// This is a generated file - do not edit.
//
// Generated from kinesis.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'kinesis.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'kinesis.pbenum.dart';

class AccessDeniedException extends $pb.GeneratedMessage {
  factory AccessDeniedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AccessDeniedException._();

  factory AccessDeniedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccessDeniedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccessDeniedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccessDeniedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccessDeniedException copyWith(
          void Function(AccessDeniedException) updates) =>
      super.copyWith((message) => updates(message as AccessDeniedException))
          as AccessDeniedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccessDeniedException create() => AccessDeniedException._();
  @$core.override
  AccessDeniedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccessDeniedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccessDeniedException>(create);
  static AccessDeniedException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class AddTagsToStreamInput extends $pb.GeneratedMessage {
  factory AddTagsToStreamInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (tags != null) result.tags.addEntries(tags);
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  AddTagsToStreamInput._();

  factory AddTagsToStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddTagsToStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddTagsToStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'AddTagsToStreamInput.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kinesis'))
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsToStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsToStreamInput copyWith(void Function(AddTagsToStreamInput) updates) =>
      super.copyWith((message) => updates(message as AddTagsToStreamInput))
          as AddTagsToStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddTagsToStreamInput create() => AddTagsToStreamInput._();
  @$core.override
  AddTagsToStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddTagsToStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddTagsToStreamInput>(create);
  static AddTagsToStreamInput? _defaultInstance;

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(0);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class ChildShard extends $pb.GeneratedMessage {
  factory ChildShard({
    HashKeyRange? hashkeyrange,
    $core.String? shardid,
    $core.Iterable<$core.String>? parentshards,
  }) {
    final result = create();
    if (hashkeyrange != null) result.hashkeyrange = hashkeyrange;
    if (shardid != null) result.shardid = shardid;
    if (parentshards != null) result.parentshards.addAll(parentshards);
    return result;
  }

  ChildShard._();

  factory ChildShard.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChildShard.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChildShard',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<HashKeyRange>(981486, _omitFieldNames ? '' : 'hashkeyrange',
        subBuilder: HashKeyRange.create)
    ..aOS(66410951, _omitFieldNames ? '' : 'shardid')
    ..pPS(152621201, _omitFieldNames ? '' : 'parentshards')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChildShard clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChildShard copyWith(void Function(ChildShard) updates) =>
      super.copyWith((message) => updates(message as ChildShard)) as ChildShard;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChildShard create() => ChildShard._();
  @$core.override
  ChildShard createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChildShard getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChildShard>(create);
  static ChildShard? _defaultInstance;

  @$pb.TagNumber(981486)
  HashKeyRange get hashkeyrange => $_getN(0);
  @$pb.TagNumber(981486)
  set hashkeyrange(HashKeyRange value) => $_setField(981486, value);
  @$pb.TagNumber(981486)
  $core.bool hasHashkeyrange() => $_has(0);
  @$pb.TagNumber(981486)
  void clearHashkeyrange() => $_clearField(981486);
  @$pb.TagNumber(981486)
  HashKeyRange ensureHashkeyrange() => $_ensure(0);

  @$pb.TagNumber(66410951)
  $core.String get shardid => $_getSZ(1);
  @$pb.TagNumber(66410951)
  set shardid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(66410951)
  $core.bool hasShardid() => $_has(1);
  @$pb.TagNumber(66410951)
  void clearShardid() => $_clearField(66410951);

  @$pb.TagNumber(152621201)
  $pb.PbList<$core.String> get parentshards => $_getList(2);
}

class Consumer extends $pb.GeneratedMessage {
  factory Consumer({
    ConsumerStatus? consumerstatus,
    $core.String? consumerarn,
    $core.String? consumername,
    $core.String? consumercreationtimestamp,
  }) {
    final result = create();
    if (consumerstatus != null) result.consumerstatus = consumerstatus;
    if (consumerarn != null) result.consumerarn = consumerarn;
    if (consumername != null) result.consumername = consumername;
    if (consumercreationtimestamp != null)
      result.consumercreationtimestamp = consumercreationtimestamp;
    return result;
  }

  Consumer._();

  factory Consumer.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Consumer.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Consumer',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<ConsumerStatus>(33471404, _omitFieldNames ? '' : 'consumerstatus',
        enumValues: ConsumerStatus.values)
    ..aOS(41107441, _omitFieldNames ? '' : 'consumerarn')
    ..aOS(70979235, _omitFieldNames ? '' : 'consumername')
    ..aOS(356226265, _omitFieldNames ? '' : 'consumercreationtimestamp')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Consumer clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Consumer copyWith(void Function(Consumer) updates) =>
      super.copyWith((message) => updates(message as Consumer)) as Consumer;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Consumer create() => Consumer._();
  @$core.override
  Consumer createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Consumer getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Consumer>(create);
  static Consumer? _defaultInstance;

  @$pb.TagNumber(33471404)
  ConsumerStatus get consumerstatus => $_getN(0);
  @$pb.TagNumber(33471404)
  set consumerstatus(ConsumerStatus value) => $_setField(33471404, value);
  @$pb.TagNumber(33471404)
  $core.bool hasConsumerstatus() => $_has(0);
  @$pb.TagNumber(33471404)
  void clearConsumerstatus() => $_clearField(33471404);

  @$pb.TagNumber(41107441)
  $core.String get consumerarn => $_getSZ(1);
  @$pb.TagNumber(41107441)
  set consumerarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(41107441)
  $core.bool hasConsumerarn() => $_has(1);
  @$pb.TagNumber(41107441)
  void clearConsumerarn() => $_clearField(41107441);

  @$pb.TagNumber(70979235)
  $core.String get consumername => $_getSZ(2);
  @$pb.TagNumber(70979235)
  set consumername($core.String value) => $_setString(2, value);
  @$pb.TagNumber(70979235)
  $core.bool hasConsumername() => $_has(2);
  @$pb.TagNumber(70979235)
  void clearConsumername() => $_clearField(70979235);

  @$pb.TagNumber(356226265)
  $core.String get consumercreationtimestamp => $_getSZ(3);
  @$pb.TagNumber(356226265)
  set consumercreationtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(356226265)
  $core.bool hasConsumercreationtimestamp() => $_has(3);
  @$pb.TagNumber(356226265)
  void clearConsumercreationtimestamp() => $_clearField(356226265);
}

class ConsumerDescription extends $pb.GeneratedMessage {
  factory ConsumerDescription({
    ConsumerStatus? consumerstatus,
    $core.String? consumerarn,
    $core.String? consumername,
    $core.String? consumercreationtimestamp,
    $core.String? streamarn,
  }) {
    final result = create();
    if (consumerstatus != null) result.consumerstatus = consumerstatus;
    if (consumerarn != null) result.consumerarn = consumerarn;
    if (consumername != null) result.consumername = consumername;
    if (consumercreationtimestamp != null)
      result.consumercreationtimestamp = consumercreationtimestamp;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  ConsumerDescription._();

  factory ConsumerDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConsumerDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConsumerDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<ConsumerStatus>(33471404, _omitFieldNames ? '' : 'consumerstatus',
        enumValues: ConsumerStatus.values)
    ..aOS(41107441, _omitFieldNames ? '' : 'consumerarn')
    ..aOS(70979235, _omitFieldNames ? '' : 'consumername')
    ..aOS(356226265, _omitFieldNames ? '' : 'consumercreationtimestamp')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConsumerDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConsumerDescription copyWith(void Function(ConsumerDescription) updates) =>
      super.copyWith((message) => updates(message as ConsumerDescription))
          as ConsumerDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConsumerDescription create() => ConsumerDescription._();
  @$core.override
  ConsumerDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConsumerDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConsumerDescription>(create);
  static ConsumerDescription? _defaultInstance;

  @$pb.TagNumber(33471404)
  ConsumerStatus get consumerstatus => $_getN(0);
  @$pb.TagNumber(33471404)
  set consumerstatus(ConsumerStatus value) => $_setField(33471404, value);
  @$pb.TagNumber(33471404)
  $core.bool hasConsumerstatus() => $_has(0);
  @$pb.TagNumber(33471404)
  void clearConsumerstatus() => $_clearField(33471404);

  @$pb.TagNumber(41107441)
  $core.String get consumerarn => $_getSZ(1);
  @$pb.TagNumber(41107441)
  set consumerarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(41107441)
  $core.bool hasConsumerarn() => $_has(1);
  @$pb.TagNumber(41107441)
  void clearConsumerarn() => $_clearField(41107441);

  @$pb.TagNumber(70979235)
  $core.String get consumername => $_getSZ(2);
  @$pb.TagNumber(70979235)
  set consumername($core.String value) => $_setString(2, value);
  @$pb.TagNumber(70979235)
  $core.bool hasConsumername() => $_has(2);
  @$pb.TagNumber(70979235)
  void clearConsumername() => $_clearField(70979235);

  @$pb.TagNumber(356226265)
  $core.String get consumercreationtimestamp => $_getSZ(3);
  @$pb.TagNumber(356226265)
  set consumercreationtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(356226265)
  $core.bool hasConsumercreationtimestamp() => $_has(3);
  @$pb.TagNumber(356226265)
  void clearConsumercreationtimestamp() => $_clearField(356226265);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class CreateStreamInput extends $pb.GeneratedMessage {
  factory CreateStreamInput({
    StreamModeDetails? streammodedetails,
    $core.int? shardcount,
    $core.int? maxrecordsizeinkib,
    $core.int? warmthroughputmibps,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? streamname,
  }) {
    final result = create();
    if (streammodedetails != null) result.streammodedetails = streammodedetails;
    if (shardcount != null) result.shardcount = shardcount;
    if (maxrecordsizeinkib != null)
      result.maxrecordsizeinkib = maxrecordsizeinkib;
    if (warmthroughputmibps != null)
      result.warmthroughputmibps = warmthroughputmibps;
    if (tags != null) result.tags.addEntries(tags);
    if (streamname != null) result.streamname = streamname;
    return result;
  }

  CreateStreamInput._();

  factory CreateStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<StreamModeDetails>(
        12139665, _omitFieldNames ? '' : 'streammodedetails',
        subBuilder: StreamModeDetails.create)
    ..aI(92831547, _omitFieldNames ? '' : 'shardcount')
    ..aI(197267253, _omitFieldNames ? '' : 'maxrecordsizeinkib')
    ..aI(259219704, _omitFieldNames ? '' : 'warmthroughputmibps')
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateStreamInput.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kinesis'))
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStreamInput copyWith(void Function(CreateStreamInput) updates) =>
      super.copyWith((message) => updates(message as CreateStreamInput))
          as CreateStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateStreamInput create() => CreateStreamInput._();
  @$core.override
  CreateStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateStreamInput>(create);
  static CreateStreamInput? _defaultInstance;

  @$pb.TagNumber(12139665)
  StreamModeDetails get streammodedetails => $_getN(0);
  @$pb.TagNumber(12139665)
  set streammodedetails(StreamModeDetails value) => $_setField(12139665, value);
  @$pb.TagNumber(12139665)
  $core.bool hasStreammodedetails() => $_has(0);
  @$pb.TagNumber(12139665)
  void clearStreammodedetails() => $_clearField(12139665);
  @$pb.TagNumber(12139665)
  StreamModeDetails ensureStreammodedetails() => $_ensure(0);

  @$pb.TagNumber(92831547)
  $core.int get shardcount => $_getIZ(1);
  @$pb.TagNumber(92831547)
  set shardcount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(92831547)
  $core.bool hasShardcount() => $_has(1);
  @$pb.TagNumber(92831547)
  void clearShardcount() => $_clearField(92831547);

  @$pb.TagNumber(197267253)
  $core.int get maxrecordsizeinkib => $_getIZ(2);
  @$pb.TagNumber(197267253)
  set maxrecordsizeinkib($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(197267253)
  $core.bool hasMaxrecordsizeinkib() => $_has(2);
  @$pb.TagNumber(197267253)
  void clearMaxrecordsizeinkib() => $_clearField(197267253);

  @$pb.TagNumber(259219704)
  $core.int get warmthroughputmibps => $_getIZ(3);
  @$pb.TagNumber(259219704)
  set warmthroughputmibps($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(259219704)
  $core.bool hasWarmthroughputmibps() => $_has(3);
  @$pb.TagNumber(259219704)
  void clearWarmthroughputmibps() => $_clearField(259219704);

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(4);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(5);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(5);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);
}

class DecreaseStreamRetentionPeriodInput extends $pb.GeneratedMessage {
  factory DecreaseStreamRetentionPeriodInput({
    $core.int? retentionperiodhours,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (retentionperiodhours != null)
      result.retentionperiodhours = retentionperiodhours;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  DecreaseStreamRetentionPeriodInput._();

  factory DecreaseStreamRetentionPeriodInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecreaseStreamRetentionPeriodInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecreaseStreamRetentionPeriodInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(396381944, _omitFieldNames ? '' : 'retentionperiodhours')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecreaseStreamRetentionPeriodInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecreaseStreamRetentionPeriodInput copyWith(
          void Function(DecreaseStreamRetentionPeriodInput) updates) =>
      super.copyWith((message) =>
              updates(message as DecreaseStreamRetentionPeriodInput))
          as DecreaseStreamRetentionPeriodInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecreaseStreamRetentionPeriodInput create() =>
      DecreaseStreamRetentionPeriodInput._();
  @$core.override
  DecreaseStreamRetentionPeriodInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecreaseStreamRetentionPeriodInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecreaseStreamRetentionPeriodInput>(
          create);
  static DecreaseStreamRetentionPeriodInput? _defaultInstance;

  @$pb.TagNumber(396381944)
  $core.int get retentionperiodhours => $_getIZ(0);
  @$pb.TagNumber(396381944)
  set retentionperiodhours($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(396381944)
  $core.bool hasRetentionperiodhours() => $_has(0);
  @$pb.TagNumber(396381944)
  void clearRetentionperiodhours() => $_clearField(396381944);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class DeleteResourcePolicyInput extends $pb.GeneratedMessage {
  factory DeleteResourcePolicyInput({
    $core.String? resourcearn,
    $core.String? streamid,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (streamid != null) result.streamid = streamid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);
}

class DeleteStreamInput extends $pb.GeneratedMessage {
  factory DeleteStreamInput({
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
    $core.bool? enforceconsumerdeletion,
  }) {
    final result = create();
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    if (enforceconsumerdeletion != null)
      result.enforceconsumerdeletion = enforceconsumerdeletion;
    return result;
  }

  DeleteStreamInput._();

  factory DeleteStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..aOB(513184424, _omitFieldNames ? '' : 'enforceconsumerdeletion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStreamInput copyWith(void Function(DeleteStreamInput) updates) =>
      super.copyWith((message) => updates(message as DeleteStreamInput))
          as DeleteStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStreamInput create() => DeleteStreamInput._();
  @$core.override
  DeleteStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStreamInput>(create);
  static DeleteStreamInput? _defaultInstance;

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(0);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(0);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(1);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(1);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(2);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(2);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);

  @$pb.TagNumber(513184424)
  $core.bool get enforceconsumerdeletion => $_getBF(3);
  @$pb.TagNumber(513184424)
  set enforceconsumerdeletion($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(513184424)
  $core.bool hasEnforceconsumerdeletion() => $_has(3);
  @$pb.TagNumber(513184424)
  void clearEnforceconsumerdeletion() => $_clearField(513184424);
}

class DeregisterStreamConsumerInput extends $pb.GeneratedMessage {
  factory DeregisterStreamConsumerInput({
    $core.String? consumerarn,
    $core.String? consumername,
    $core.String? streamid,
    $core.String? streamarn,
  }) {
    final result = create();
    if (consumerarn != null) result.consumerarn = consumerarn;
    if (consumername != null) result.consumername = consumername;
    if (streamid != null) result.streamid = streamid;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  DeregisterStreamConsumerInput._();

  factory DeregisterStreamConsumerInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeregisterStreamConsumerInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeregisterStreamConsumerInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(41107441, _omitFieldNames ? '' : 'consumerarn')
    ..aOS(70979235, _omitFieldNames ? '' : 'consumername')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeregisterStreamConsumerInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeregisterStreamConsumerInput copyWith(
          void Function(DeregisterStreamConsumerInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeregisterStreamConsumerInput))
          as DeregisterStreamConsumerInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeregisterStreamConsumerInput create() =>
      DeregisterStreamConsumerInput._();
  @$core.override
  DeregisterStreamConsumerInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeregisterStreamConsumerInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeregisterStreamConsumerInput>(create);
  static DeregisterStreamConsumerInput? _defaultInstance;

  @$pb.TagNumber(41107441)
  $core.String get consumerarn => $_getSZ(0);
  @$pb.TagNumber(41107441)
  set consumerarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41107441)
  $core.bool hasConsumerarn() => $_has(0);
  @$pb.TagNumber(41107441)
  void clearConsumerarn() => $_clearField(41107441);

  @$pb.TagNumber(70979235)
  $core.String get consumername => $_getSZ(1);
  @$pb.TagNumber(70979235)
  set consumername($core.String value) => $_setString(1, value);
  @$pb.TagNumber(70979235)
  $core.bool hasConsumername() => $_has(1);
  @$pb.TagNumber(70979235)
  void clearConsumername() => $_clearField(70979235);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class DescribeAccountSettingsInput extends $pb.GeneratedMessage {
  factory DescribeAccountSettingsInput() => create();

  DescribeAccountSettingsInput._();

  factory DescribeAccountSettingsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAccountSettingsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAccountSettingsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsInput copyWith(
          void Function(DescribeAccountSettingsInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAccountSettingsInput))
          as DescribeAccountSettingsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsInput create() =>
      DescribeAccountSettingsInput._();
  @$core.override
  DescribeAccountSettingsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAccountSettingsInput>(create);
  static DescribeAccountSettingsInput? _defaultInstance;
}

class DescribeAccountSettingsOutput extends $pb.GeneratedMessage {
  factory DescribeAccountSettingsOutput({
    MinimumThroughputBillingCommitmentOutput?
        minimumthroughputbillingcommitment,
  }) {
    final result = create();
    if (minimumthroughputbillingcommitment != null)
      result.minimumthroughputbillingcommitment =
          minimumthroughputbillingcommitment;
    return result;
  }

  DescribeAccountSettingsOutput._();

  factory DescribeAccountSettingsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAccountSettingsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAccountSettingsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<MinimumThroughputBillingCommitmentOutput>(
        335326962, _omitFieldNames ? '' : 'minimumthroughputbillingcommitment',
        subBuilder: MinimumThroughputBillingCommitmentOutput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsOutput copyWith(
          void Function(DescribeAccountSettingsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAccountSettingsOutput))
          as DescribeAccountSettingsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsOutput create() =>
      DescribeAccountSettingsOutput._();
  @$core.override
  DescribeAccountSettingsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAccountSettingsOutput>(create);
  static DescribeAccountSettingsOutput? _defaultInstance;

  @$pb.TagNumber(335326962)
  MinimumThroughputBillingCommitmentOutput
      get minimumthroughputbillingcommitment => $_getN(0);
  @$pb.TagNumber(335326962)
  set minimumthroughputbillingcommitment(
          MinimumThroughputBillingCommitmentOutput value) =>
      $_setField(335326962, value);
  @$pb.TagNumber(335326962)
  $core.bool hasMinimumthroughputbillingcommitment() => $_has(0);
  @$pb.TagNumber(335326962)
  void clearMinimumthroughputbillingcommitment() => $_clearField(335326962);
  @$pb.TagNumber(335326962)
  MinimumThroughputBillingCommitmentOutput
      ensureMinimumthroughputbillingcommitment() => $_ensure(0);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
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
    $core.int? ondemandstreamcountlimit,
    $core.int? ondemandstreamcount,
    $core.int? openshardcount,
    $core.int? shardlimit,
  }) {
    final result = create();
    if (ondemandstreamcountlimit != null)
      result.ondemandstreamcountlimit = ondemandstreamcountlimit;
    if (ondemandstreamcount != null)
      result.ondemandstreamcount = ondemandstreamcount;
    if (openshardcount != null) result.openshardcount = openshardcount;
    if (shardlimit != null) result.shardlimit = shardlimit;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(38712458, _omitFieldNames ? '' : 'ondemandstreamcountlimit')
    ..aI(386875707, _omitFieldNames ? '' : 'ondemandstreamcount')
    ..aI(476409287, _omitFieldNames ? '' : 'openshardcount')
    ..aI(515722839, _omitFieldNames ? '' : 'shardlimit')
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

  @$pb.TagNumber(38712458)
  $core.int get ondemandstreamcountlimit => $_getIZ(0);
  @$pb.TagNumber(38712458)
  set ondemandstreamcountlimit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(38712458)
  $core.bool hasOndemandstreamcountlimit() => $_has(0);
  @$pb.TagNumber(38712458)
  void clearOndemandstreamcountlimit() => $_clearField(38712458);

  @$pb.TagNumber(386875707)
  $core.int get ondemandstreamcount => $_getIZ(1);
  @$pb.TagNumber(386875707)
  set ondemandstreamcount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(386875707)
  $core.bool hasOndemandstreamcount() => $_has(1);
  @$pb.TagNumber(386875707)
  void clearOndemandstreamcount() => $_clearField(386875707);

  @$pb.TagNumber(476409287)
  $core.int get openshardcount => $_getIZ(2);
  @$pb.TagNumber(476409287)
  set openshardcount($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(476409287)
  $core.bool hasOpenshardcount() => $_has(2);
  @$pb.TagNumber(476409287)
  void clearOpenshardcount() => $_clearField(476409287);

  @$pb.TagNumber(515722839)
  $core.int get shardlimit => $_getIZ(3);
  @$pb.TagNumber(515722839)
  set shardlimit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(515722839)
  $core.bool hasShardlimit() => $_has(3);
  @$pb.TagNumber(515722839)
  void clearShardlimit() => $_clearField(515722839);
}

class DescribeStreamConsumerInput extends $pb.GeneratedMessage {
  factory DescribeStreamConsumerInput({
    $core.String? consumerarn,
    $core.String? consumername,
    $core.String? streamid,
    $core.String? streamarn,
  }) {
    final result = create();
    if (consumerarn != null) result.consumerarn = consumerarn;
    if (consumername != null) result.consumername = consumername;
    if (streamid != null) result.streamid = streamid;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  DescribeStreamConsumerInput._();

  factory DescribeStreamConsumerInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStreamConsumerInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStreamConsumerInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(41107441, _omitFieldNames ? '' : 'consumerarn')
    ..aOS(70979235, _omitFieldNames ? '' : 'consumername')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamConsumerInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamConsumerInput copyWith(
          void Function(DescribeStreamConsumerInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeStreamConsumerInput))
          as DescribeStreamConsumerInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStreamConsumerInput create() =>
      DescribeStreamConsumerInput._();
  @$core.override
  DescribeStreamConsumerInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStreamConsumerInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStreamConsumerInput>(create);
  static DescribeStreamConsumerInput? _defaultInstance;

  @$pb.TagNumber(41107441)
  $core.String get consumerarn => $_getSZ(0);
  @$pb.TagNumber(41107441)
  set consumerarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41107441)
  $core.bool hasConsumerarn() => $_has(0);
  @$pb.TagNumber(41107441)
  void clearConsumerarn() => $_clearField(41107441);

  @$pb.TagNumber(70979235)
  $core.String get consumername => $_getSZ(1);
  @$pb.TagNumber(70979235)
  set consumername($core.String value) => $_setString(1, value);
  @$pb.TagNumber(70979235)
  $core.bool hasConsumername() => $_has(1);
  @$pb.TagNumber(70979235)
  void clearConsumername() => $_clearField(70979235);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class DescribeStreamConsumerOutput extends $pb.GeneratedMessage {
  factory DescribeStreamConsumerOutput({
    ConsumerDescription? consumerdescription,
  }) {
    final result = create();
    if (consumerdescription != null)
      result.consumerdescription = consumerdescription;
    return result;
  }

  DescribeStreamConsumerOutput._();

  factory DescribeStreamConsumerOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStreamConsumerOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStreamConsumerOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<ConsumerDescription>(
        524123998, _omitFieldNames ? '' : 'consumerdescription',
        subBuilder: ConsumerDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamConsumerOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamConsumerOutput copyWith(
          void Function(DescribeStreamConsumerOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeStreamConsumerOutput))
          as DescribeStreamConsumerOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStreamConsumerOutput create() =>
      DescribeStreamConsumerOutput._();
  @$core.override
  DescribeStreamConsumerOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStreamConsumerOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStreamConsumerOutput>(create);
  static DescribeStreamConsumerOutput? _defaultInstance;

  @$pb.TagNumber(524123998)
  ConsumerDescription get consumerdescription => $_getN(0);
  @$pb.TagNumber(524123998)
  set consumerdescription(ConsumerDescription value) =>
      $_setField(524123998, value);
  @$pb.TagNumber(524123998)
  $core.bool hasConsumerdescription() => $_has(0);
  @$pb.TagNumber(524123998)
  void clearConsumerdescription() => $_clearField(524123998);
  @$pb.TagNumber(524123998)
  ConsumerDescription ensureConsumerdescription() => $_ensure(0);
}

class DescribeStreamInput extends $pb.GeneratedMessage {
  factory DescribeStreamInput({
    $core.String? exclusivestartshardid,
    $core.int? limit,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (exclusivestartshardid != null)
      result.exclusivestartshardid = exclusivestartshardid;
    if (limit != null) result.limit = limit;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  DescribeStreamInput._();

  factory DescribeStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(44771587, _omitFieldNames ? '' : 'exclusivestartshardid')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamInput copyWith(void Function(DescribeStreamInput) updates) =>
      super.copyWith((message) => updates(message as DescribeStreamInput))
          as DescribeStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStreamInput create() => DescribeStreamInput._();
  @$core.override
  DescribeStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStreamInput>(create);
  static DescribeStreamInput? _defaultInstance;

  @$pb.TagNumber(44771587)
  $core.String get exclusivestartshardid => $_getSZ(0);
  @$pb.TagNumber(44771587)
  set exclusivestartshardid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(44771587)
  $core.bool hasExclusivestartshardid() => $_has(0);
  @$pb.TagNumber(44771587)
  void clearExclusivestartshardid() => $_clearField(44771587);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(3);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(3);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class DescribeStreamOutput extends $pb.GeneratedMessage {
  factory DescribeStreamOutput({
    StreamDescription? streamdescription,
  }) {
    final result = create();
    if (streamdescription != null) result.streamdescription = streamdescription;
    return result;
  }

  DescribeStreamOutput._();

  factory DescribeStreamOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStreamOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStreamOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<StreamDescription>(
        363737034, _omitFieldNames ? '' : 'streamdescription',
        subBuilder: StreamDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamOutput copyWith(void Function(DescribeStreamOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeStreamOutput))
          as DescribeStreamOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStreamOutput create() => DescribeStreamOutput._();
  @$core.override
  DescribeStreamOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStreamOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStreamOutput>(create);
  static DescribeStreamOutput? _defaultInstance;

  @$pb.TagNumber(363737034)
  StreamDescription get streamdescription => $_getN(0);
  @$pb.TagNumber(363737034)
  set streamdescription(StreamDescription value) =>
      $_setField(363737034, value);
  @$pb.TagNumber(363737034)
  $core.bool hasStreamdescription() => $_has(0);
  @$pb.TagNumber(363737034)
  void clearStreamdescription() => $_clearField(363737034);
  @$pb.TagNumber(363737034)
  StreamDescription ensureStreamdescription() => $_ensure(0);
}

class DescribeStreamSummaryInput extends $pb.GeneratedMessage {
  factory DescribeStreamSummaryInput({
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  DescribeStreamSummaryInput._();

  factory DescribeStreamSummaryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStreamSummaryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStreamSummaryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamSummaryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamSummaryInput copyWith(
          void Function(DescribeStreamSummaryInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeStreamSummaryInput))
          as DescribeStreamSummaryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStreamSummaryInput create() => DescribeStreamSummaryInput._();
  @$core.override
  DescribeStreamSummaryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStreamSummaryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStreamSummaryInput>(create);
  static DescribeStreamSummaryInput? _defaultInstance;

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(0);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(0);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(1);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(1);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(2);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(2);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class DescribeStreamSummaryOutput extends $pb.GeneratedMessage {
  factory DescribeStreamSummaryOutput({
    StreamDescriptionSummary? streamdescriptionsummary,
  }) {
    final result = create();
    if (streamdescriptionsummary != null)
      result.streamdescriptionsummary = streamdescriptionsummary;
    return result;
  }

  DescribeStreamSummaryOutput._();

  factory DescribeStreamSummaryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStreamSummaryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStreamSummaryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<StreamDescriptionSummary>(
        478899272, _omitFieldNames ? '' : 'streamdescriptionsummary',
        subBuilder: StreamDescriptionSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamSummaryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStreamSummaryOutput copyWith(
          void Function(DescribeStreamSummaryOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeStreamSummaryOutput))
          as DescribeStreamSummaryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStreamSummaryOutput create() =>
      DescribeStreamSummaryOutput._();
  @$core.override
  DescribeStreamSummaryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStreamSummaryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStreamSummaryOutput>(create);
  static DescribeStreamSummaryOutput? _defaultInstance;

  @$pb.TagNumber(478899272)
  StreamDescriptionSummary get streamdescriptionsummary => $_getN(0);
  @$pb.TagNumber(478899272)
  set streamdescriptionsummary(StreamDescriptionSummary value) =>
      $_setField(478899272, value);
  @$pb.TagNumber(478899272)
  $core.bool hasStreamdescriptionsummary() => $_has(0);
  @$pb.TagNumber(478899272)
  void clearStreamdescriptionsummary() => $_clearField(478899272);
  @$pb.TagNumber(478899272)
  StreamDescriptionSummary ensureStreamdescriptionsummary() => $_ensure(0);
}

class DisableEnhancedMonitoringInput extends $pb.GeneratedMessage {
  factory DisableEnhancedMonitoringInput({
    $core.Iterable<MetricsName>? shardlevelmetrics,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (shardlevelmetrics != null)
      result.shardlevelmetrics.addAll(shardlevelmetrics);
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  DisableEnhancedMonitoringInput._();

  factory DisableEnhancedMonitoringInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableEnhancedMonitoringInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableEnhancedMonitoringInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pc<MetricsName>(406711021, _omitFieldNames ? '' : 'shardlevelmetrics',
        $pb.PbFieldType.KE,
        valueOf: MetricsName.valueOf,
        enumValues: MetricsName.values,
        defaultEnumValue: MetricsName.METRICS_NAME_INCOMING_RECORDS)
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableEnhancedMonitoringInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableEnhancedMonitoringInput copyWith(
          void Function(DisableEnhancedMonitoringInput) updates) =>
      super.copyWith(
              (message) => updates(message as DisableEnhancedMonitoringInput))
          as DisableEnhancedMonitoringInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableEnhancedMonitoringInput create() =>
      DisableEnhancedMonitoringInput._();
  @$core.override
  DisableEnhancedMonitoringInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableEnhancedMonitoringInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableEnhancedMonitoringInput>(create);
  static DisableEnhancedMonitoringInput? _defaultInstance;

  @$pb.TagNumber(406711021)
  $pb.PbList<MetricsName> get shardlevelmetrics => $_getList(0);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class EnableEnhancedMonitoringInput extends $pb.GeneratedMessage {
  factory EnableEnhancedMonitoringInput({
    $core.Iterable<MetricsName>? shardlevelmetrics,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (shardlevelmetrics != null)
      result.shardlevelmetrics.addAll(shardlevelmetrics);
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  EnableEnhancedMonitoringInput._();

  factory EnableEnhancedMonitoringInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableEnhancedMonitoringInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableEnhancedMonitoringInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pc<MetricsName>(406711021, _omitFieldNames ? '' : 'shardlevelmetrics',
        $pb.PbFieldType.KE,
        valueOf: MetricsName.valueOf,
        enumValues: MetricsName.values,
        defaultEnumValue: MetricsName.METRICS_NAME_INCOMING_RECORDS)
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableEnhancedMonitoringInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableEnhancedMonitoringInput copyWith(
          void Function(EnableEnhancedMonitoringInput) updates) =>
      super.copyWith(
              (message) => updates(message as EnableEnhancedMonitoringInput))
          as EnableEnhancedMonitoringInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableEnhancedMonitoringInput create() =>
      EnableEnhancedMonitoringInput._();
  @$core.override
  EnableEnhancedMonitoringInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableEnhancedMonitoringInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableEnhancedMonitoringInput>(create);
  static EnableEnhancedMonitoringInput? _defaultInstance;

  @$pb.TagNumber(406711021)
  $pb.PbList<MetricsName> get shardlevelmetrics => $_getList(0);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class EnhancedMetrics extends $pb.GeneratedMessage {
  factory EnhancedMetrics({
    $core.Iterable<MetricsName>? shardlevelmetrics,
  }) {
    final result = create();
    if (shardlevelmetrics != null)
      result.shardlevelmetrics.addAll(shardlevelmetrics);
    return result;
  }

  EnhancedMetrics._();

  factory EnhancedMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnhancedMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnhancedMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pc<MetricsName>(406711021, _omitFieldNames ? '' : 'shardlevelmetrics',
        $pb.PbFieldType.KE,
        valueOf: MetricsName.valueOf,
        enumValues: MetricsName.values,
        defaultEnumValue: MetricsName.METRICS_NAME_INCOMING_RECORDS)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnhancedMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnhancedMetrics copyWith(void Function(EnhancedMetrics) updates) =>
      super.copyWith((message) => updates(message as EnhancedMetrics))
          as EnhancedMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnhancedMetrics create() => EnhancedMetrics._();
  @$core.override
  EnhancedMetrics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnhancedMetrics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnhancedMetrics>(create);
  static EnhancedMetrics? _defaultInstance;

  @$pb.TagNumber(406711021)
  $pb.PbList<MetricsName> get shardlevelmetrics => $_getList(0);
}

class EnhancedMonitoringOutput extends $pb.GeneratedMessage {
  factory EnhancedMonitoringOutput({
    $core.Iterable<MetricsName>? desiredshardlevelmetrics,
    $core.Iterable<MetricsName>? currentshardlevelmetrics,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (desiredshardlevelmetrics != null)
      result.desiredshardlevelmetrics.addAll(desiredshardlevelmetrics);
    if (currentshardlevelmetrics != null)
      result.currentshardlevelmetrics.addAll(currentshardlevelmetrics);
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  EnhancedMonitoringOutput._();

  factory EnhancedMonitoringOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnhancedMonitoringOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnhancedMonitoringOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pc<MetricsName>(231120285,
        _omitFieldNames ? '' : 'desiredshardlevelmetrics', $pb.PbFieldType.KE,
        valueOf: MetricsName.valueOf,
        enumValues: MetricsName.values,
        defaultEnumValue: MetricsName.METRICS_NAME_INCOMING_RECORDS)
    ..pc<MetricsName>(453144018,
        _omitFieldNames ? '' : 'currentshardlevelmetrics', $pb.PbFieldType.KE,
        valueOf: MetricsName.valueOf,
        enumValues: MetricsName.values,
        defaultEnumValue: MetricsName.METRICS_NAME_INCOMING_RECORDS)
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnhancedMonitoringOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnhancedMonitoringOutput copyWith(
          void Function(EnhancedMonitoringOutput) updates) =>
      super.copyWith((message) => updates(message as EnhancedMonitoringOutput))
          as EnhancedMonitoringOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnhancedMonitoringOutput create() => EnhancedMonitoringOutput._();
  @$core.override
  EnhancedMonitoringOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnhancedMonitoringOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnhancedMonitoringOutput>(create);
  static EnhancedMonitoringOutput? _defaultInstance;

  @$pb.TagNumber(231120285)
  $pb.PbList<MetricsName> get desiredshardlevelmetrics => $_getList(0);

  @$pb.TagNumber(453144018)
  $pb.PbList<MetricsName> get currentshardlevelmetrics => $_getList(1);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class ExpiredIteratorException extends $pb.GeneratedMessage {
  factory ExpiredIteratorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExpiredIteratorException._();

  factory ExpiredIteratorException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExpiredIteratorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExpiredIteratorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredIteratorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredIteratorException copyWith(
          void Function(ExpiredIteratorException) updates) =>
      super.copyWith((message) => updates(message as ExpiredIteratorException))
          as ExpiredIteratorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExpiredIteratorException create() => ExpiredIteratorException._();
  @$core.override
  ExpiredIteratorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExpiredIteratorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExpiredIteratorException>(create);
  static ExpiredIteratorException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExpiredNextTokenException extends $pb.GeneratedMessage {
  factory ExpiredNextTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExpiredNextTokenException._();

  factory ExpiredNextTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExpiredNextTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExpiredNextTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredNextTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredNextTokenException copyWith(
          void Function(ExpiredNextTokenException) updates) =>
      super.copyWith((message) => updates(message as ExpiredNextTokenException))
          as ExpiredNextTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExpiredNextTokenException create() => ExpiredNextTokenException._();
  @$core.override
  ExpiredNextTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExpiredNextTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExpiredNextTokenException>(create);
  static ExpiredNextTokenException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class GetRecordsInput extends $pb.GeneratedMessage {
  factory GetRecordsInput({
    $core.String? sharditerator,
    $core.int? limit,
    $core.String? streamid,
    $core.String? streamarn,
  }) {
    final result = create();
    if (sharditerator != null) result.sharditerator = sharditerator;
    if (limit != null) result.limit = limit;
    if (streamid != null) result.streamid = streamid;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  GetRecordsInput._();

  factory GetRecordsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRecordsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRecordsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(379619650, _omitFieldNames ? '' : 'sharditerator')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRecordsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRecordsInput copyWith(void Function(GetRecordsInput) updates) =>
      super.copyWith((message) => updates(message as GetRecordsInput))
          as GetRecordsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRecordsInput create() => GetRecordsInput._();
  @$core.override
  GetRecordsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRecordsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRecordsInput>(create);
  static GetRecordsInput? _defaultInstance;

  @$pb.TagNumber(379619650)
  $core.String get sharditerator => $_getSZ(0);
  @$pb.TagNumber(379619650)
  set sharditerator($core.String value) => $_setString(0, value);
  @$pb.TagNumber(379619650)
  $core.bool hasSharditerator() => $_has(0);
  @$pb.TagNumber(379619650)
  void clearSharditerator() => $_clearField(379619650);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class GetRecordsOutput extends $pb.GeneratedMessage {
  factory GetRecordsOutput({
    $core.Iterable<ChildShard>? childshards,
    $core.Iterable<Record>? records,
    $core.String? nextsharditerator,
    $fixnum.Int64? millisbehindlatest,
  }) {
    final result = create();
    if (childshards != null) result.childshards.addAll(childshards);
    if (records != null) result.records.addAll(records);
    if (nextsharditerator != null) result.nextsharditerator = nextsharditerator;
    if (millisbehindlatest != null)
      result.millisbehindlatest = millisbehindlatest;
    return result;
  }

  GetRecordsOutput._();

  factory GetRecordsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRecordsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRecordsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pPM<ChildShard>(348740657, _omitFieldNames ? '' : 'childshards',
        subBuilder: ChildShard.create)
    ..pPM<Record>(423557454, _omitFieldNames ? '' : 'records',
        subBuilder: Record.create)
    ..aOS(442470571, _omitFieldNames ? '' : 'nextsharditerator')
    ..aInt64(456422057, _omitFieldNames ? '' : 'millisbehindlatest')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRecordsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRecordsOutput copyWith(void Function(GetRecordsOutput) updates) =>
      super.copyWith((message) => updates(message as GetRecordsOutput))
          as GetRecordsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRecordsOutput create() => GetRecordsOutput._();
  @$core.override
  GetRecordsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRecordsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRecordsOutput>(create);
  static GetRecordsOutput? _defaultInstance;

  @$pb.TagNumber(348740657)
  $pb.PbList<ChildShard> get childshards => $_getList(0);

  @$pb.TagNumber(423557454)
  $pb.PbList<Record> get records => $_getList(1);

  @$pb.TagNumber(442470571)
  $core.String get nextsharditerator => $_getSZ(2);
  @$pb.TagNumber(442470571)
  set nextsharditerator($core.String value) => $_setString(2, value);
  @$pb.TagNumber(442470571)
  $core.bool hasNextsharditerator() => $_has(2);
  @$pb.TagNumber(442470571)
  void clearNextsharditerator() => $_clearField(442470571);

  @$pb.TagNumber(456422057)
  $fixnum.Int64 get millisbehindlatest => $_getI64(3);
  @$pb.TagNumber(456422057)
  set millisbehindlatest($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(456422057)
  $core.bool hasMillisbehindlatest() => $_has(3);
  @$pb.TagNumber(456422057)
  void clearMillisbehindlatest() => $_clearField(456422057);
}

class GetResourcePolicyInput extends $pb.GeneratedMessage {
  factory GetResourcePolicyInput({
    $core.String? resourcearn,
    $core.String? streamid,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (streamid != null) result.streamid = streamid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);
}

class GetResourcePolicyOutput extends $pb.GeneratedMessage {
  factory GetResourcePolicyOutput({
    $core.String? policy,
  }) {
    final result = create();
    if (policy != null) result.policy = policy;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
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
}

class GetShardIteratorInput extends $pb.GeneratedMessage {
  factory GetShardIteratorInput({
    $core.String? shardid,
    $core.String? startingsequencenumber,
    $core.String? timestamp,
    ShardIteratorType? sharditeratortype,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (shardid != null) result.shardid = shardid;
    if (startingsequencenumber != null)
      result.startingsequencenumber = startingsequencenumber;
    if (timestamp != null) result.timestamp = timestamp;
    if (sharditeratortype != null) result.sharditeratortype = sharditeratortype;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  GetShardIteratorInput._();

  factory GetShardIteratorInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetShardIteratorInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetShardIteratorInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(66410951, _omitFieldNames ? '' : 'shardid')
    ..aOS(88770150, _omitFieldNames ? '' : 'startingsequencenumber')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aE<ShardIteratorType>(
        229371818, _omitFieldNames ? '' : 'sharditeratortype',
        enumValues: ShardIteratorType.values)
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetShardIteratorInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetShardIteratorInput copyWith(
          void Function(GetShardIteratorInput) updates) =>
      super.copyWith((message) => updates(message as GetShardIteratorInput))
          as GetShardIteratorInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetShardIteratorInput create() => GetShardIteratorInput._();
  @$core.override
  GetShardIteratorInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetShardIteratorInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetShardIteratorInput>(create);
  static GetShardIteratorInput? _defaultInstance;

  @$pb.TagNumber(66410951)
  $core.String get shardid => $_getSZ(0);
  @$pb.TagNumber(66410951)
  set shardid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(66410951)
  $core.bool hasShardid() => $_has(0);
  @$pb.TagNumber(66410951)
  void clearShardid() => $_clearField(66410951);

  @$pb.TagNumber(88770150)
  $core.String get startingsequencenumber => $_getSZ(1);
  @$pb.TagNumber(88770150)
  set startingsequencenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(88770150)
  $core.bool hasStartingsequencenumber() => $_has(1);
  @$pb.TagNumber(88770150)
  void clearStartingsequencenumber() => $_clearField(88770150);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(2);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(2);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(229371818)
  ShardIteratorType get sharditeratortype => $_getN(3);
  @$pb.TagNumber(229371818)
  set sharditeratortype(ShardIteratorType value) =>
      $_setField(229371818, value);
  @$pb.TagNumber(229371818)
  $core.bool hasSharditeratortype() => $_has(3);
  @$pb.TagNumber(229371818)
  void clearSharditeratortype() => $_clearField(229371818);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(4);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(4);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(5);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(5);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(6);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(6);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class GetShardIteratorOutput extends $pb.GeneratedMessage {
  factory GetShardIteratorOutput({
    $core.String? sharditerator,
  }) {
    final result = create();
    if (sharditerator != null) result.sharditerator = sharditerator;
    return result;
  }

  GetShardIteratorOutput._();

  factory GetShardIteratorOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetShardIteratorOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetShardIteratorOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(379619650, _omitFieldNames ? '' : 'sharditerator')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetShardIteratorOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetShardIteratorOutput copyWith(
          void Function(GetShardIteratorOutput) updates) =>
      super.copyWith((message) => updates(message as GetShardIteratorOutput))
          as GetShardIteratorOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetShardIteratorOutput create() => GetShardIteratorOutput._();
  @$core.override
  GetShardIteratorOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetShardIteratorOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetShardIteratorOutput>(create);
  static GetShardIteratorOutput? _defaultInstance;

  @$pb.TagNumber(379619650)
  $core.String get sharditerator => $_getSZ(0);
  @$pb.TagNumber(379619650)
  set sharditerator($core.String value) => $_setString(0, value);
  @$pb.TagNumber(379619650)
  $core.bool hasSharditerator() => $_has(0);
  @$pb.TagNumber(379619650)
  void clearSharditerator() => $_clearField(379619650);
}

class HashKeyRange extends $pb.GeneratedMessage {
  factory HashKeyRange({
    $core.String? startinghashkey,
    $core.String? endinghashkey,
  }) {
    final result = create();
    if (startinghashkey != null) result.startinghashkey = startinghashkey;
    if (endinghashkey != null) result.endinghashkey = endinghashkey;
    return result;
  }

  HashKeyRange._();

  factory HashKeyRange.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HashKeyRange.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HashKeyRange',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(202551085, _omitFieldNames ? '' : 'startinghashkey')
    ..aOS(399409906, _omitFieldNames ? '' : 'endinghashkey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HashKeyRange clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HashKeyRange copyWith(void Function(HashKeyRange) updates) =>
      super.copyWith((message) => updates(message as HashKeyRange))
          as HashKeyRange;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HashKeyRange create() => HashKeyRange._();
  @$core.override
  HashKeyRange createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HashKeyRange getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HashKeyRange>(create);
  static HashKeyRange? _defaultInstance;

  @$pb.TagNumber(202551085)
  $core.String get startinghashkey => $_getSZ(0);
  @$pb.TagNumber(202551085)
  set startinghashkey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(202551085)
  $core.bool hasStartinghashkey() => $_has(0);
  @$pb.TagNumber(202551085)
  void clearStartinghashkey() => $_clearField(202551085);

  @$pb.TagNumber(399409906)
  $core.String get endinghashkey => $_getSZ(1);
  @$pb.TagNumber(399409906)
  set endinghashkey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(399409906)
  $core.bool hasEndinghashkey() => $_has(1);
  @$pb.TagNumber(399409906)
  void clearEndinghashkey() => $_clearField(399409906);
}

class IncreaseStreamRetentionPeriodInput extends $pb.GeneratedMessage {
  factory IncreaseStreamRetentionPeriodInput({
    $core.int? retentionperiodhours,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (retentionperiodhours != null)
      result.retentionperiodhours = retentionperiodhours;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  IncreaseStreamRetentionPeriodInput._();

  factory IncreaseStreamRetentionPeriodInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IncreaseStreamRetentionPeriodInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IncreaseStreamRetentionPeriodInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(396381944, _omitFieldNames ? '' : 'retentionperiodhours')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncreaseStreamRetentionPeriodInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncreaseStreamRetentionPeriodInput copyWith(
          void Function(IncreaseStreamRetentionPeriodInput) updates) =>
      super.copyWith((message) =>
              updates(message as IncreaseStreamRetentionPeriodInput))
          as IncreaseStreamRetentionPeriodInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IncreaseStreamRetentionPeriodInput create() =>
      IncreaseStreamRetentionPeriodInput._();
  @$core.override
  IncreaseStreamRetentionPeriodInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IncreaseStreamRetentionPeriodInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IncreaseStreamRetentionPeriodInput>(
          create);
  static IncreaseStreamRetentionPeriodInput? _defaultInstance;

  @$pb.TagNumber(396381944)
  $core.int get retentionperiodhours => $_getIZ(0);
  @$pb.TagNumber(396381944)
  set retentionperiodhours($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(396381944)
  $core.bool hasRetentionperiodhours() => $_has(0);
  @$pb.TagNumber(396381944)
  void clearRetentionperiodhours() => $_clearField(396381944);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class InternalFailureException extends $pb.GeneratedMessage {
  factory InternalFailureException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalFailureException._();

  factory InternalFailureException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalFailureException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalFailureException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalFailureException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalFailureException copyWith(
          void Function(InternalFailureException) updates) =>
      super.copyWith((message) => updates(message as InternalFailureException))
          as InternalFailureException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalFailureException create() => InternalFailureException._();
  @$core.override
  InternalFailureException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalFailureException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalFailureException>(create);
  static InternalFailureException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidArgumentException extends $pb.GeneratedMessage {
  factory InvalidArgumentException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidArgumentException._();

  factory InvalidArgumentException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidArgumentException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidArgumentException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArgumentException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArgumentException copyWith(
          void Function(InvalidArgumentException) updates) =>
      super.copyWith((message) => updates(message as InvalidArgumentException))
          as InvalidArgumentException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidArgumentException create() => InvalidArgumentException._();
  @$core.override
  InvalidArgumentException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidArgumentException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidArgumentException>(create);
  static InvalidArgumentException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSAccessDeniedException extends $pb.GeneratedMessage {
  factory KMSAccessDeniedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSAccessDeniedException._();

  factory KMSAccessDeniedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSAccessDeniedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSAccessDeniedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSAccessDeniedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSAccessDeniedException copyWith(
          void Function(KMSAccessDeniedException) updates) =>
      super.copyWith((message) => updates(message as KMSAccessDeniedException))
          as KMSAccessDeniedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSAccessDeniedException create() => KMSAccessDeniedException._();
  @$core.override
  KMSAccessDeniedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSAccessDeniedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSAccessDeniedException>(create);
  static KMSAccessDeniedException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSDisabledException extends $pb.GeneratedMessage {
  factory KMSDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSDisabledException._();

  factory KMSDisabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSDisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSDisabledException copyWith(void Function(KMSDisabledException) updates) =>
      super.copyWith((message) => updates(message as KMSDisabledException))
          as KMSDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSDisabledException create() => KMSDisabledException._();
  @$core.override
  KMSDisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSDisabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSDisabledException>(create);
  static KMSDisabledException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
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

class KMSNotFoundException extends $pb.GeneratedMessage {
  factory KMSNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSNotFoundException._();

  factory KMSNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSNotFoundException copyWith(void Function(KMSNotFoundException) updates) =>
      super.copyWith((message) => updates(message as KMSNotFoundException))
          as KMSNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSNotFoundException create() => KMSNotFoundException._();
  @$core.override
  KMSNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSNotFoundException>(create);
  static KMSNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSOptInRequired extends $pb.GeneratedMessage {
  factory KMSOptInRequired({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSOptInRequired._();

  factory KMSOptInRequired.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSOptInRequired.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSOptInRequired',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSOptInRequired clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSOptInRequired copyWith(void Function(KMSOptInRequired) updates) =>
      super.copyWith((message) => updates(message as KMSOptInRequired))
          as KMSOptInRequired;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSOptInRequired create() => KMSOptInRequired._();
  @$core.override
  KMSOptInRequired createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSOptInRequired getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSOptInRequired>(create);
  static KMSOptInRequired? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSThrottlingException extends $pb.GeneratedMessage {
  factory KMSThrottlingException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSThrottlingException._();

  factory KMSThrottlingException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSThrottlingException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSThrottlingException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSThrottlingException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSThrottlingException copyWith(
          void Function(KMSThrottlingException) updates) =>
      super.copyWith((message) => updates(message as KMSThrottlingException))
          as KMSThrottlingException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSThrottlingException create() => KMSThrottlingException._();
  @$core.override
  KMSThrottlingException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSThrottlingException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSThrottlingException>(create);
  static KMSThrottlingException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
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

class ListShardsInput extends $pb.GeneratedMessage {
  factory ListShardsInput({
    $core.String? exclusivestartshardid,
    $core.String? nexttoken,
    $core.String? streamcreationtimestamp,
    ShardFilter? shardfilter,
    $core.int? maxresults,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (exclusivestartshardid != null)
      result.exclusivestartshardid = exclusivestartshardid;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (streamcreationtimestamp != null)
      result.streamcreationtimestamp = streamcreationtimestamp;
    if (shardfilter != null) result.shardfilter = shardfilter;
    if (maxresults != null) result.maxresults = maxresults;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  ListShardsInput._();

  factory ListShardsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListShardsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListShardsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(44771587, _omitFieldNames ? '' : 'exclusivestartshardid')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(224951013, _omitFieldNames ? '' : 'streamcreationtimestamp')
    ..aOM<ShardFilter>(230254710, _omitFieldNames ? '' : 'shardfilter',
        subBuilder: ShardFilter.create)
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListShardsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListShardsInput copyWith(void Function(ListShardsInput) updates) =>
      super.copyWith((message) => updates(message as ListShardsInput))
          as ListShardsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListShardsInput create() => ListShardsInput._();
  @$core.override
  ListShardsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListShardsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListShardsInput>(create);
  static ListShardsInput? _defaultInstance;

  @$pb.TagNumber(44771587)
  $core.String get exclusivestartshardid => $_getSZ(0);
  @$pb.TagNumber(44771587)
  set exclusivestartshardid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(44771587)
  $core.bool hasExclusivestartshardid() => $_has(0);
  @$pb.TagNumber(44771587)
  void clearExclusivestartshardid() => $_clearField(44771587);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(224951013)
  $core.String get streamcreationtimestamp => $_getSZ(2);
  @$pb.TagNumber(224951013)
  set streamcreationtimestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(224951013)
  $core.bool hasStreamcreationtimestamp() => $_has(2);
  @$pb.TagNumber(224951013)
  void clearStreamcreationtimestamp() => $_clearField(224951013);

  @$pb.TagNumber(230254710)
  ShardFilter get shardfilter => $_getN(3);
  @$pb.TagNumber(230254710)
  set shardfilter(ShardFilter value) => $_setField(230254710, value);
  @$pb.TagNumber(230254710)
  $core.bool hasShardfilter() => $_has(3);
  @$pb.TagNumber(230254710)
  void clearShardfilter() => $_clearField(230254710);
  @$pb.TagNumber(230254710)
  ShardFilter ensureShardfilter() => $_ensure(3);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(4);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(4);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(5);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(5);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(6);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(6);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(7);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(7);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class ListShardsOutput extends $pb.GeneratedMessage {
  factory ListShardsOutput({
    $core.String? nexttoken,
    $core.Iterable<Shard>? shards,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (shards != null) result.shards.addAll(shards);
    return result;
  }

  ListShardsOutput._();

  factory ListShardsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListShardsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListShardsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Shard>(437117641, _omitFieldNames ? '' : 'shards',
        subBuilder: Shard.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListShardsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListShardsOutput copyWith(void Function(ListShardsOutput) updates) =>
      super.copyWith((message) => updates(message as ListShardsOutput))
          as ListShardsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListShardsOutput create() => ListShardsOutput._();
  @$core.override
  ListShardsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListShardsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListShardsOutput>(create);
  static ListShardsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(437117641)
  $pb.PbList<Shard> get shards => $_getList(1);
}

class ListStreamConsumersInput extends $pb.GeneratedMessage {
  factory ListStreamConsumersInput({
    $core.String? nexttoken,
    $core.String? streamcreationtimestamp,
    $core.int? maxresults,
    $core.String? streamid,
    $core.String? streamarn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (streamcreationtimestamp != null)
      result.streamcreationtimestamp = streamcreationtimestamp;
    if (maxresults != null) result.maxresults = maxresults;
    if (streamid != null) result.streamid = streamid;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  ListStreamConsumersInput._();

  factory ListStreamConsumersInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStreamConsumersInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStreamConsumersInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(224951013, _omitFieldNames ? '' : 'streamcreationtimestamp')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamConsumersInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamConsumersInput copyWith(
          void Function(ListStreamConsumersInput) updates) =>
      super.copyWith((message) => updates(message as ListStreamConsumersInput))
          as ListStreamConsumersInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStreamConsumersInput create() => ListStreamConsumersInput._();
  @$core.override
  ListStreamConsumersInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStreamConsumersInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStreamConsumersInput>(create);
  static ListStreamConsumersInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(224951013)
  $core.String get streamcreationtimestamp => $_getSZ(1);
  @$pb.TagNumber(224951013)
  set streamcreationtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(224951013)
  $core.bool hasStreamcreationtimestamp() => $_has(1);
  @$pb.TagNumber(224951013)
  void clearStreamcreationtimestamp() => $_clearField(224951013);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(3);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(3);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class ListStreamConsumersOutput extends $pb.GeneratedMessage {
  factory ListStreamConsumersOutput({
    $core.String? nexttoken,
    $core.Iterable<Consumer>? consumers,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (consumers != null) result.consumers.addAll(consumers);
    return result;
  }

  ListStreamConsumersOutput._();

  factory ListStreamConsumersOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStreamConsumersOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStreamConsumersOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Consumer>(348153455, _omitFieldNames ? '' : 'consumers',
        subBuilder: Consumer.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamConsumersOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamConsumersOutput copyWith(
          void Function(ListStreamConsumersOutput) updates) =>
      super.copyWith((message) => updates(message as ListStreamConsumersOutput))
          as ListStreamConsumersOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStreamConsumersOutput create() => ListStreamConsumersOutput._();
  @$core.override
  ListStreamConsumersOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStreamConsumersOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStreamConsumersOutput>(create);
  static ListStreamConsumersOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(348153455)
  $pb.PbList<Consumer> get consumers => $_getList(1);
}

class ListStreamsInput extends $pb.GeneratedMessage {
  factory ListStreamsInput({
    $core.String? nexttoken,
    $core.int? limit,
    $core.String? exclusivestartstreamname,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (limit != null) result.limit = limit;
    if (exclusivestartstreamname != null)
      result.exclusivestartstreamname = exclusivestartstreamname;
    return result;
  }

  ListStreamsInput._();

  factory ListStreamsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStreamsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStreamsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(431108955, _omitFieldNames ? '' : 'exclusivestartstreamname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamsInput copyWith(void Function(ListStreamsInput) updates) =>
      super.copyWith((message) => updates(message as ListStreamsInput))
          as ListStreamsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStreamsInput create() => ListStreamsInput._();
  @$core.override
  ListStreamsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStreamsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStreamsInput>(create);
  static ListStreamsInput? _defaultInstance;

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

  @$pb.TagNumber(431108955)
  $core.String get exclusivestartstreamname => $_getSZ(2);
  @$pb.TagNumber(431108955)
  set exclusivestartstreamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(431108955)
  $core.bool hasExclusivestartstreamname() => $_has(2);
  @$pb.TagNumber(431108955)
  void clearExclusivestartstreamname() => $_clearField(431108955);
}

class ListStreamsOutput extends $pb.GeneratedMessage {
  factory ListStreamsOutput({
    $core.bool? hasmorestreams,
    $core.String? nexttoken,
    $core.Iterable<StreamSummary>? streamsummaries,
    $core.Iterable<$core.String>? streamnames,
  }) {
    final result = create();
    if (hasmorestreams != null) result.hasmorestreams = hasmorestreams;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (streamsummaries != null) result.streamsummaries.addAll(streamsummaries);
    if (streamnames != null) result.streamnames.addAll(streamnames);
    return result;
  }

  ListStreamsOutput._();

  factory ListStreamsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStreamsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStreamsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOB(181206864, _omitFieldNames ? '' : 'hasmorestreams')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<StreamSummary>(334316244, _omitFieldNames ? '' : 'streamsummaries',
        subBuilder: StreamSummary.create)
    ..pPS(530210288, _omitFieldNames ? '' : 'streamnames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStreamsOutput copyWith(void Function(ListStreamsOutput) updates) =>
      super.copyWith((message) => updates(message as ListStreamsOutput))
          as ListStreamsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStreamsOutput create() => ListStreamsOutput._();
  @$core.override
  ListStreamsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStreamsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStreamsOutput>(create);
  static ListStreamsOutput? _defaultInstance;

  @$pb.TagNumber(181206864)
  $core.bool get hasmorestreams => $_getBF(0);
  @$pb.TagNumber(181206864)
  set hasmorestreams($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(181206864)
  $core.bool hasHasmorestreams() => $_has(0);
  @$pb.TagNumber(181206864)
  void clearHasmorestreams() => $_clearField(181206864);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(334316244)
  $pb.PbList<StreamSummary> get streamsummaries => $_getList(2);

  @$pb.TagNumber(530210288)
  $pb.PbList<$core.String> get streamnames => $_getList(3);
}

class ListTagsForResourceInput extends $pb.GeneratedMessage {
  factory ListTagsForResourceInput({
    $core.String? resourcearn,
    $core.String? streamid,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (streamid != null) result.streamid = streamid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
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

class ListTagsForStreamInput extends $pb.GeneratedMessage {
  factory ListTagsForStreamInput({
    $core.int? limit,
    $core.String? streamid,
    $core.String? exclusivestarttagkey,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (streamid != null) result.streamid = streamid;
    if (exclusivestarttagkey != null)
      result.exclusivestarttagkey = exclusivestarttagkey;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  ListTagsForStreamInput._();

  factory ListTagsForStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470465287, _omitFieldNames ? '' : 'exclusivestarttagkey')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForStreamInput copyWith(
          void Function(ListTagsForStreamInput) updates) =>
      super.copyWith((message) => updates(message as ListTagsForStreamInput))
          as ListTagsForStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForStreamInput create() => ListTagsForStreamInput._();
  @$core.override
  ListTagsForStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForStreamInput>(create);
  static ListTagsForStreamInput? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470465287)
  $core.String get exclusivestarttagkey => $_getSZ(2);
  @$pb.TagNumber(470465287)
  set exclusivestarttagkey($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470465287)
  $core.bool hasExclusivestarttagkey() => $_has(2);
  @$pb.TagNumber(470465287)
  void clearExclusivestarttagkey() => $_clearField(470465287);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(3);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(3);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class ListTagsForStreamOutput extends $pb.GeneratedMessage {
  factory ListTagsForStreamOutput({
    $core.Iterable<Tag>? tags,
    $core.bool? hasmoretags,
  }) {
    final result = create();
    if (tags != null) result.tags.addAll(tags);
    if (hasmoretags != null) result.hasmoretags = hasmoretags;
    return result;
  }

  ListTagsForStreamOutput._();

  factory ListTagsForStreamOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForStreamOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForStreamOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOB(397963248, _omitFieldNames ? '' : 'hasmoretags')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForStreamOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForStreamOutput copyWith(
          void Function(ListTagsForStreamOutput) updates) =>
      super.copyWith((message) => updates(message as ListTagsForStreamOutput))
          as ListTagsForStreamOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForStreamOutput create() => ListTagsForStreamOutput._();
  @$core.override
  ListTagsForStreamOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForStreamOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForStreamOutput>(create);
  static ListTagsForStreamOutput? _defaultInstance;

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(0);

  @$pb.TagNumber(397963248)
  $core.bool get hasmoretags => $_getBF(1);
  @$pb.TagNumber(397963248)
  set hasmoretags($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(397963248)
  $core.bool hasHasmoretags() => $_has(1);
  @$pb.TagNumber(397963248)
  void clearHasmoretags() => $_clearField(397963248);
}

class MergeShardsInput extends $pb.GeneratedMessage {
  factory MergeShardsInput({
    $core.String? streamid,
    $core.String? adjacentshardtomerge,
    $core.String? streamname,
    $core.String? shardtomerge,
    $core.String? streamarn,
  }) {
    final result = create();
    if (streamid != null) result.streamid = streamid;
    if (adjacentshardtomerge != null)
      result.adjacentshardtomerge = adjacentshardtomerge;
    if (streamname != null) result.streamname = streamname;
    if (shardtomerge != null) result.shardtomerge = shardtomerge;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  MergeShardsInput._();

  factory MergeShardsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MergeShardsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MergeShardsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(435883605, _omitFieldNames ? '' : 'adjacentshardtomerge')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(473343671, _omitFieldNames ? '' : 'shardtomerge')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeShardsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeShardsInput copyWith(void Function(MergeShardsInput) updates) =>
      super.copyWith((message) => updates(message as MergeShardsInput))
          as MergeShardsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MergeShardsInput create() => MergeShardsInput._();
  @$core.override
  MergeShardsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MergeShardsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MergeShardsInput>(create);
  static MergeShardsInput? _defaultInstance;

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(0);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(0);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(435883605)
  $core.String get adjacentshardtomerge => $_getSZ(1);
  @$pb.TagNumber(435883605)
  set adjacentshardtomerge($core.String value) => $_setString(1, value);
  @$pb.TagNumber(435883605)
  $core.bool hasAdjacentshardtomerge() => $_has(1);
  @$pb.TagNumber(435883605)
  void clearAdjacentshardtomerge() => $_clearField(435883605);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(473343671)
  $core.String get shardtomerge => $_getSZ(3);
  @$pb.TagNumber(473343671)
  set shardtomerge($core.String value) => $_setString(3, value);
  @$pb.TagNumber(473343671)
  $core.bool hasShardtomerge() => $_has(3);
  @$pb.TagNumber(473343671)
  void clearShardtomerge() => $_clearField(473343671);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class MinimumThroughputBillingCommitmentInput extends $pb.GeneratedMessage {
  factory MinimumThroughputBillingCommitmentInput({
    MinimumThroughputBillingCommitmentInputStatus? status,
  }) {
    final result = create();
    if (status != null) result.status = status;
    return result;
  }

  MinimumThroughputBillingCommitmentInput._();

  factory MinimumThroughputBillingCommitmentInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MinimumThroughputBillingCommitmentInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MinimumThroughputBillingCommitmentInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<MinimumThroughputBillingCommitmentInputStatus>(
        6222352, _omitFieldNames ? '' : 'status',
        enumValues: MinimumThroughputBillingCommitmentInputStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MinimumThroughputBillingCommitmentInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MinimumThroughputBillingCommitmentInput copyWith(
          void Function(MinimumThroughputBillingCommitmentInput) updates) =>
      super.copyWith((message) =>
              updates(message as MinimumThroughputBillingCommitmentInput))
          as MinimumThroughputBillingCommitmentInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MinimumThroughputBillingCommitmentInput create() =>
      MinimumThroughputBillingCommitmentInput._();
  @$core.override
  MinimumThroughputBillingCommitmentInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MinimumThroughputBillingCommitmentInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          MinimumThroughputBillingCommitmentInput>(create);
  static MinimumThroughputBillingCommitmentInput? _defaultInstance;

  @$pb.TagNumber(6222352)
  MinimumThroughputBillingCommitmentInputStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(MinimumThroughputBillingCommitmentInputStatus value) =>
      $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
}

class MinimumThroughputBillingCommitmentOutput extends $pb.GeneratedMessage {
  factory MinimumThroughputBillingCommitmentOutput({
    MinimumThroughputBillingCommitmentOutputStatus? status,
    $core.String? startedat,
    $core.String? earliestallowedendat,
    $core.String? endedat,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (startedat != null) result.startedat = startedat;
    if (earliestallowedendat != null)
      result.earliestallowedendat = earliestallowedendat;
    if (endedat != null) result.endedat = endedat;
    return result;
  }

  MinimumThroughputBillingCommitmentOutput._();

  factory MinimumThroughputBillingCommitmentOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MinimumThroughputBillingCommitmentOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MinimumThroughputBillingCommitmentOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<MinimumThroughputBillingCommitmentOutputStatus>(
        6222352, _omitFieldNames ? '' : 'status',
        enumValues: MinimumThroughputBillingCommitmentOutputStatus.values)
    ..aOS(77629404, _omitFieldNames ? '' : 'startedat')
    ..aOS(77927885, _omitFieldNames ? '' : 'earliestallowedendat')
    ..aOS(104122351, _omitFieldNames ? '' : 'endedat')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MinimumThroughputBillingCommitmentOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MinimumThroughputBillingCommitmentOutput copyWith(
          void Function(MinimumThroughputBillingCommitmentOutput) updates) =>
      super.copyWith((message) =>
              updates(message as MinimumThroughputBillingCommitmentOutput))
          as MinimumThroughputBillingCommitmentOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MinimumThroughputBillingCommitmentOutput create() =>
      MinimumThroughputBillingCommitmentOutput._();
  @$core.override
  MinimumThroughputBillingCommitmentOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MinimumThroughputBillingCommitmentOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          MinimumThroughputBillingCommitmentOutput>(create);
  static MinimumThroughputBillingCommitmentOutput? _defaultInstance;

  @$pb.TagNumber(6222352)
  MinimumThroughputBillingCommitmentOutputStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(MinimumThroughputBillingCommitmentOutputStatus value) =>
      $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(77629404)
  $core.String get startedat => $_getSZ(1);
  @$pb.TagNumber(77629404)
  set startedat($core.String value) => $_setString(1, value);
  @$pb.TagNumber(77629404)
  $core.bool hasStartedat() => $_has(1);
  @$pb.TagNumber(77629404)
  void clearStartedat() => $_clearField(77629404);

  @$pb.TagNumber(77927885)
  $core.String get earliestallowedendat => $_getSZ(2);
  @$pb.TagNumber(77927885)
  set earliestallowedendat($core.String value) => $_setString(2, value);
  @$pb.TagNumber(77927885)
  $core.bool hasEarliestallowedendat() => $_has(2);
  @$pb.TagNumber(77927885)
  void clearEarliestallowedendat() => $_clearField(77927885);

  @$pb.TagNumber(104122351)
  $core.String get endedat => $_getSZ(3);
  @$pb.TagNumber(104122351)
  set endedat($core.String value) => $_setString(3, value);
  @$pb.TagNumber(104122351)
  $core.bool hasEndedat() => $_has(3);
  @$pb.TagNumber(104122351)
  void clearEndedat() => $_clearField(104122351);
}

class ProvisionedThroughputExceededException extends $pb.GeneratedMessage {
  factory ProvisionedThroughputExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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
}

class PutRecordInput extends $pb.GeneratedMessage {
  factory PutRecordInput({
    $core.String? sequencenumberforordering,
    $core.String? partitionkey,
    $core.String? streamid,
    $core.String? explicithashkey,
    $core.String? streamname,
    $core.String? streamarn,
    $core.List<$core.int>? data,
  }) {
    final result = create();
    if (sequencenumberforordering != null)
      result.sequencenumberforordering = sequencenumberforordering;
    if (partitionkey != null) result.partitionkey = partitionkey;
    if (streamid != null) result.streamid = streamid;
    if (explicithashkey != null) result.explicithashkey = explicithashkey;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    if (data != null) result.data = data;
    return result;
  }

  PutRecordInput._();

  factory PutRecordInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRecordInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRecordInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(73727749, _omitFieldNames ? '' : 'sequencenumberforordering')
    ..aOS(379379617, _omitFieldNames ? '' : 'partitionkey')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(438141821, _omitFieldNames ? '' : 'explicithashkey')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..a<$core.List<$core.int>>(
        525498822, _omitFieldNames ? '' : 'data', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordInput copyWith(void Function(PutRecordInput) updates) =>
      super.copyWith((message) => updates(message as PutRecordInput))
          as PutRecordInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRecordInput create() => PutRecordInput._();
  @$core.override
  PutRecordInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRecordInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRecordInput>(create);
  static PutRecordInput? _defaultInstance;

  @$pb.TagNumber(73727749)
  $core.String get sequencenumberforordering => $_getSZ(0);
  @$pb.TagNumber(73727749)
  set sequencenumberforordering($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73727749)
  $core.bool hasSequencenumberforordering() => $_has(0);
  @$pb.TagNumber(73727749)
  void clearSequencenumberforordering() => $_clearField(73727749);

  @$pb.TagNumber(379379617)
  $core.String get partitionkey => $_getSZ(1);
  @$pb.TagNumber(379379617)
  set partitionkey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(379379617)
  $core.bool hasPartitionkey() => $_has(1);
  @$pb.TagNumber(379379617)
  void clearPartitionkey() => $_clearField(379379617);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(438141821)
  $core.String get explicithashkey => $_getSZ(3);
  @$pb.TagNumber(438141821)
  set explicithashkey($core.String value) => $_setString(3, value);
  @$pb.TagNumber(438141821)
  $core.bool hasExplicithashkey() => $_has(3);
  @$pb.TagNumber(438141821)
  void clearExplicithashkey() => $_clearField(438141821);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(4);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(4);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(5);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(5);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);

  @$pb.TagNumber(525498822)
  $core.List<$core.int> get data => $_getN(6);
  @$pb.TagNumber(525498822)
  set data($core.List<$core.int> value) => $_setBytes(6, value);
  @$pb.TagNumber(525498822)
  $core.bool hasData() => $_has(6);
  @$pb.TagNumber(525498822)
  void clearData() => $_clearField(525498822);
}

class PutRecordOutput extends $pb.GeneratedMessage {
  factory PutRecordOutput({
    $core.String? shardid,
    $core.String? sequencenumber,
    EncryptionType? encryptiontype,
  }) {
    final result = create();
    if (shardid != null) result.shardid = shardid;
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (encryptiontype != null) result.encryptiontype = encryptiontype;
    return result;
  }

  PutRecordOutput._();

  factory PutRecordOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRecordOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRecordOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(66410951, _omitFieldNames ? '' : 'shardid')
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aE<EncryptionType>(264007605, _omitFieldNames ? '' : 'encryptiontype',
        enumValues: EncryptionType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordOutput copyWith(void Function(PutRecordOutput) updates) =>
      super.copyWith((message) => updates(message as PutRecordOutput))
          as PutRecordOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRecordOutput create() => PutRecordOutput._();
  @$core.override
  PutRecordOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRecordOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRecordOutput>(create);
  static PutRecordOutput? _defaultInstance;

  @$pb.TagNumber(66410951)
  $core.String get shardid => $_getSZ(0);
  @$pb.TagNumber(66410951)
  set shardid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(66410951)
  $core.bool hasShardid() => $_has(0);
  @$pb.TagNumber(66410951)
  void clearShardid() => $_clearField(66410951);

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(1);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(1);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(264007605)
  EncryptionType get encryptiontype => $_getN(2);
  @$pb.TagNumber(264007605)
  set encryptiontype(EncryptionType value) => $_setField(264007605, value);
  @$pb.TagNumber(264007605)
  $core.bool hasEncryptiontype() => $_has(2);
  @$pb.TagNumber(264007605)
  void clearEncryptiontype() => $_clearField(264007605);
}

class PutRecordsInput extends $pb.GeneratedMessage {
  factory PutRecordsInput({
    $core.String? streamid,
    $core.Iterable<PutRecordsRequestEntry>? records,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (streamid != null) result.streamid = streamid;
    if (records != null) result.records.addAll(records);
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  PutRecordsInput._();

  factory PutRecordsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRecordsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRecordsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..pPM<PutRecordsRequestEntry>(423557454, _omitFieldNames ? '' : 'records',
        subBuilder: PutRecordsRequestEntry.create)
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsInput copyWith(void Function(PutRecordsInput) updates) =>
      super.copyWith((message) => updates(message as PutRecordsInput))
          as PutRecordsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRecordsInput create() => PutRecordsInput._();
  @$core.override
  PutRecordsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRecordsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRecordsInput>(create);
  static PutRecordsInput? _defaultInstance;

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(0);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(0);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(423557454)
  $pb.PbList<PutRecordsRequestEntry> get records => $_getList(1);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class PutRecordsOutput extends $pb.GeneratedMessage {
  factory PutRecordsOutput({
    $core.int? failedrecordcount,
    EncryptionType? encryptiontype,
    $core.Iterable<PutRecordsResultEntry>? records,
  }) {
    final result = create();
    if (failedrecordcount != null) result.failedrecordcount = failedrecordcount;
    if (encryptiontype != null) result.encryptiontype = encryptiontype;
    if (records != null) result.records.addAll(records);
    return result;
  }

  PutRecordsOutput._();

  factory PutRecordsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRecordsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRecordsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(89270457, _omitFieldNames ? '' : 'failedrecordcount')
    ..aE<EncryptionType>(264007605, _omitFieldNames ? '' : 'encryptiontype',
        enumValues: EncryptionType.values)
    ..pPM<PutRecordsResultEntry>(423557454, _omitFieldNames ? '' : 'records',
        subBuilder: PutRecordsResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsOutput copyWith(void Function(PutRecordsOutput) updates) =>
      super.copyWith((message) => updates(message as PutRecordsOutput))
          as PutRecordsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRecordsOutput create() => PutRecordsOutput._();
  @$core.override
  PutRecordsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRecordsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRecordsOutput>(create);
  static PutRecordsOutput? _defaultInstance;

  @$pb.TagNumber(89270457)
  $core.int get failedrecordcount => $_getIZ(0);
  @$pb.TagNumber(89270457)
  set failedrecordcount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(89270457)
  $core.bool hasFailedrecordcount() => $_has(0);
  @$pb.TagNumber(89270457)
  void clearFailedrecordcount() => $_clearField(89270457);

  @$pb.TagNumber(264007605)
  EncryptionType get encryptiontype => $_getN(1);
  @$pb.TagNumber(264007605)
  set encryptiontype(EncryptionType value) => $_setField(264007605, value);
  @$pb.TagNumber(264007605)
  $core.bool hasEncryptiontype() => $_has(1);
  @$pb.TagNumber(264007605)
  void clearEncryptiontype() => $_clearField(264007605);

  @$pb.TagNumber(423557454)
  $pb.PbList<PutRecordsResultEntry> get records => $_getList(2);
}

class PutRecordsRequestEntry extends $pb.GeneratedMessage {
  factory PutRecordsRequestEntry({
    $core.String? partitionkey,
    $core.String? explicithashkey,
    $core.List<$core.int>? data,
  }) {
    final result = create();
    if (partitionkey != null) result.partitionkey = partitionkey;
    if (explicithashkey != null) result.explicithashkey = explicithashkey;
    if (data != null) result.data = data;
    return result;
  }

  PutRecordsRequestEntry._();

  factory PutRecordsRequestEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRecordsRequestEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRecordsRequestEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(379379617, _omitFieldNames ? '' : 'partitionkey')
    ..aOS(438141821, _omitFieldNames ? '' : 'explicithashkey')
    ..a<$core.List<$core.int>>(
        525498822, _omitFieldNames ? '' : 'data', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsRequestEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsRequestEntry copyWith(
          void Function(PutRecordsRequestEntry) updates) =>
      super.copyWith((message) => updates(message as PutRecordsRequestEntry))
          as PutRecordsRequestEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRecordsRequestEntry create() => PutRecordsRequestEntry._();
  @$core.override
  PutRecordsRequestEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRecordsRequestEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRecordsRequestEntry>(create);
  static PutRecordsRequestEntry? _defaultInstance;

  @$pb.TagNumber(379379617)
  $core.String get partitionkey => $_getSZ(0);
  @$pb.TagNumber(379379617)
  set partitionkey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(379379617)
  $core.bool hasPartitionkey() => $_has(0);
  @$pb.TagNumber(379379617)
  void clearPartitionkey() => $_clearField(379379617);

  @$pb.TagNumber(438141821)
  $core.String get explicithashkey => $_getSZ(1);
  @$pb.TagNumber(438141821)
  set explicithashkey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(438141821)
  $core.bool hasExplicithashkey() => $_has(1);
  @$pb.TagNumber(438141821)
  void clearExplicithashkey() => $_clearField(438141821);

  @$pb.TagNumber(525498822)
  $core.List<$core.int> get data => $_getN(2);
  @$pb.TagNumber(525498822)
  set data($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(525498822)
  $core.bool hasData() => $_has(2);
  @$pb.TagNumber(525498822)
  void clearData() => $_clearField(525498822);
}

class PutRecordsResultEntry extends $pb.GeneratedMessage {
  factory PutRecordsResultEntry({
    $core.String? errorcode,
    $core.String? shardid,
    $core.String? sequencenumber,
    $core.String? errormessage,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (shardid != null) result.shardid = shardid;
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  PutRecordsResultEntry._();

  factory PutRecordsResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRecordsResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRecordsResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(66410951, _omitFieldNames ? '' : 'shardid')
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRecordsResultEntry copyWith(
          void Function(PutRecordsResultEntry) updates) =>
      super.copyWith((message) => updates(message as PutRecordsResultEntry))
          as PutRecordsResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRecordsResultEntry create() => PutRecordsResultEntry._();
  @$core.override
  PutRecordsResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRecordsResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRecordsResultEntry>(create);
  static PutRecordsResultEntry? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(66410951)
  $core.String get shardid => $_getSZ(1);
  @$pb.TagNumber(66410951)
  set shardid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(66410951)
  $core.bool hasShardid() => $_has(1);
  @$pb.TagNumber(66410951)
  void clearShardid() => $_clearField(66410951);

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(2);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(2, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(2);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(3);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(3, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(3);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class PutResourcePolicyInput extends $pb.GeneratedMessage {
  factory PutResourcePolicyInput({
    $core.String? resourcearn,
    $core.String? streamid,
    $core.String? policy,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (streamid != null) result.streamid = streamid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(2);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(2);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class Record extends $pb.GeneratedMessage {
  factory Record({
    $core.String? approximatearrivaltimestamp,
    $core.String? sequencenumber,
    EncryptionType? encryptiontype,
    $core.String? partitionkey,
    $core.List<$core.int>? data,
  }) {
    final result = create();
    if (approximatearrivaltimestamp != null)
      result.approximatearrivaltimestamp = approximatearrivaltimestamp;
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (encryptiontype != null) result.encryptiontype = encryptiontype;
    if (partitionkey != null) result.partitionkey = partitionkey;
    if (data != null) result.data = data;
    return result;
  }

  Record._();

  factory Record.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Record.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Record',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(95039887, _omitFieldNames ? '' : 'approximatearrivaltimestamp')
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aE<EncryptionType>(264007605, _omitFieldNames ? '' : 'encryptiontype',
        enumValues: EncryptionType.values)
    ..aOS(379379617, _omitFieldNames ? '' : 'partitionkey')
    ..a<$core.List<$core.int>>(
        525498822, _omitFieldNames ? '' : 'data', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Record clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Record copyWith(void Function(Record) updates) =>
      super.copyWith((message) => updates(message as Record)) as Record;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Record create() => Record._();
  @$core.override
  Record createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Record getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Record>(create);
  static Record? _defaultInstance;

  @$pb.TagNumber(95039887)
  $core.String get approximatearrivaltimestamp => $_getSZ(0);
  @$pb.TagNumber(95039887)
  set approximatearrivaltimestamp($core.String value) => $_setString(0, value);
  @$pb.TagNumber(95039887)
  $core.bool hasApproximatearrivaltimestamp() => $_has(0);
  @$pb.TagNumber(95039887)
  void clearApproximatearrivaltimestamp() => $_clearField(95039887);

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(1);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(1);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(264007605)
  EncryptionType get encryptiontype => $_getN(2);
  @$pb.TagNumber(264007605)
  set encryptiontype(EncryptionType value) => $_setField(264007605, value);
  @$pb.TagNumber(264007605)
  $core.bool hasEncryptiontype() => $_has(2);
  @$pb.TagNumber(264007605)
  void clearEncryptiontype() => $_clearField(264007605);

  @$pb.TagNumber(379379617)
  $core.String get partitionkey => $_getSZ(3);
  @$pb.TagNumber(379379617)
  set partitionkey($core.String value) => $_setString(3, value);
  @$pb.TagNumber(379379617)
  $core.bool hasPartitionkey() => $_has(3);
  @$pb.TagNumber(379379617)
  void clearPartitionkey() => $_clearField(379379617);

  @$pb.TagNumber(525498822)
  $core.List<$core.int> get data => $_getN(4);
  @$pb.TagNumber(525498822)
  set data($core.List<$core.int> value) => $_setBytes(4, value);
  @$pb.TagNumber(525498822)
  $core.bool hasData() => $_has(4);
  @$pb.TagNumber(525498822)
  void clearData() => $_clearField(525498822);
}

class RegisterStreamConsumerInput extends $pb.GeneratedMessage {
  factory RegisterStreamConsumerInput({
    $core.String? consumername,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? streamid,
    $core.String? streamarn,
  }) {
    final result = create();
    if (consumername != null) result.consumername = consumername;
    if (tags != null) result.tags.addEntries(tags);
    if (streamid != null) result.streamid = streamid;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  RegisterStreamConsumerInput._();

  factory RegisterStreamConsumerInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegisterStreamConsumerInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegisterStreamConsumerInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(70979235, _omitFieldNames ? '' : 'consumername')
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'RegisterStreamConsumerInput.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kinesis'))
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterStreamConsumerInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterStreamConsumerInput copyWith(
          void Function(RegisterStreamConsumerInput) updates) =>
      super.copyWith(
              (message) => updates(message as RegisterStreamConsumerInput))
          as RegisterStreamConsumerInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegisterStreamConsumerInput create() =>
      RegisterStreamConsumerInput._();
  @$core.override
  RegisterStreamConsumerInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegisterStreamConsumerInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegisterStreamConsumerInput>(create);
  static RegisterStreamConsumerInput? _defaultInstance;

  @$pb.TagNumber(70979235)
  $core.String get consumername => $_getSZ(0);
  @$pb.TagNumber(70979235)
  set consumername($core.String value) => $_setString(0, value);
  @$pb.TagNumber(70979235)
  $core.bool hasConsumername() => $_has(0);
  @$pb.TagNumber(70979235)
  void clearConsumername() => $_clearField(70979235);

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(1);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class RegisterStreamConsumerOutput extends $pb.GeneratedMessage {
  factory RegisterStreamConsumerOutput({
    Consumer? consumer,
  }) {
    final result = create();
    if (consumer != null) result.consumer = consumer;
    return result;
  }

  RegisterStreamConsumerOutput._();

  factory RegisterStreamConsumerOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegisterStreamConsumerOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegisterStreamConsumerOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<Consumer>(482032874, _omitFieldNames ? '' : 'consumer',
        subBuilder: Consumer.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterStreamConsumerOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterStreamConsumerOutput copyWith(
          void Function(RegisterStreamConsumerOutput) updates) =>
      super.copyWith(
              (message) => updates(message as RegisterStreamConsumerOutput))
          as RegisterStreamConsumerOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegisterStreamConsumerOutput create() =>
      RegisterStreamConsumerOutput._();
  @$core.override
  RegisterStreamConsumerOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegisterStreamConsumerOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegisterStreamConsumerOutput>(create);
  static RegisterStreamConsumerOutput? _defaultInstance;

  @$pb.TagNumber(482032874)
  Consumer get consumer => $_getN(0);
  @$pb.TagNumber(482032874)
  set consumer(Consumer value) => $_setField(482032874, value);
  @$pb.TagNumber(482032874)
  $core.bool hasConsumer() => $_has(0);
  @$pb.TagNumber(482032874)
  void clearConsumer() => $_clearField(482032874);
  @$pb.TagNumber(482032874)
  Consumer ensureConsumer() => $_ensure(0);
}

class RemoveTagsFromStreamInput extends $pb.GeneratedMessage {
  factory RemoveTagsFromStreamInput({
    $core.Iterable<$core.String>? tagkeys,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  RemoveTagsFromStreamInput._();

  factory RemoveTagsFromStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveTagsFromStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveTagsFromStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsFromStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsFromStreamInput copyWith(
          void Function(RemoveTagsFromStreamInput) updates) =>
      super.copyWith((message) => updates(message as RemoveTagsFromStreamInput))
          as RemoveTagsFromStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveTagsFromStreamInput create() => RemoveTagsFromStreamInput._();
  @$core.override
  RemoveTagsFromStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveTagsFromStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveTagsFromStreamInput>(create);
  static RemoveTagsFromStreamInput? _defaultInstance;

  @$pb.TagNumber(320659964)
  $pb.PbList<$core.String> get tagkeys => $_getList(0);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
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

class SequenceNumberRange extends $pb.GeneratedMessage {
  factory SequenceNumberRange({
    $core.String? startingsequencenumber,
    $core.String? endingsequencenumber,
  }) {
    final result = create();
    if (startingsequencenumber != null)
      result.startingsequencenumber = startingsequencenumber;
    if (endingsequencenumber != null)
      result.endingsequencenumber = endingsequencenumber;
    return result;
  }

  SequenceNumberRange._();

  factory SequenceNumberRange.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SequenceNumberRange.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SequenceNumberRange',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(88770150, _omitFieldNames ? '' : 'startingsequencenumber')
    ..aOS(328882583, _omitFieldNames ? '' : 'endingsequencenumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SequenceNumberRange clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SequenceNumberRange copyWith(void Function(SequenceNumberRange) updates) =>
      super.copyWith((message) => updates(message as SequenceNumberRange))
          as SequenceNumberRange;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SequenceNumberRange create() => SequenceNumberRange._();
  @$core.override
  SequenceNumberRange createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SequenceNumberRange getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SequenceNumberRange>(create);
  static SequenceNumberRange? _defaultInstance;

  @$pb.TagNumber(88770150)
  $core.String get startingsequencenumber => $_getSZ(0);
  @$pb.TagNumber(88770150)
  set startingsequencenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88770150)
  $core.bool hasStartingsequencenumber() => $_has(0);
  @$pb.TagNumber(88770150)
  void clearStartingsequencenumber() => $_clearField(88770150);

  @$pb.TagNumber(328882583)
  $core.String get endingsequencenumber => $_getSZ(1);
  @$pb.TagNumber(328882583)
  set endingsequencenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(328882583)
  $core.bool hasEndingsequencenumber() => $_has(1);
  @$pb.TagNumber(328882583)
  void clearEndingsequencenumber() => $_clearField(328882583);
}

class Shard extends $pb.GeneratedMessage {
  factory Shard({
    HashKeyRange? hashkeyrange,
    $core.String? shardid,
    SequenceNumberRange? sequencenumberrange,
    $core.String? adjacentparentshardid,
    $core.String? parentshardid,
  }) {
    final result = create();
    if (hashkeyrange != null) result.hashkeyrange = hashkeyrange;
    if (shardid != null) result.shardid = shardid;
    if (sequencenumberrange != null)
      result.sequencenumberrange = sequencenumberrange;
    if (adjacentparentshardid != null)
      result.adjacentparentshardid = adjacentparentshardid;
    if (parentshardid != null) result.parentshardid = parentshardid;
    return result;
  }

  Shard._();

  factory Shard.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Shard.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Shard',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<HashKeyRange>(981486, _omitFieldNames ? '' : 'hashkeyrange',
        subBuilder: HashKeyRange.create)
    ..aOS(66410951, _omitFieldNames ? '' : 'shardid')
    ..aOM<SequenceNumberRange>(
        86578567, _omitFieldNames ? '' : 'sequencenumberrange',
        subBuilder: SequenceNumberRange.create)
    ..aOS(310461245, _omitFieldNames ? '' : 'adjacentparentshardid')
    ..aOS(431117103, _omitFieldNames ? '' : 'parentshardid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Shard clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Shard copyWith(void Function(Shard) updates) =>
      super.copyWith((message) => updates(message as Shard)) as Shard;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Shard create() => Shard._();
  @$core.override
  Shard createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Shard getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Shard>(create);
  static Shard? _defaultInstance;

  @$pb.TagNumber(981486)
  HashKeyRange get hashkeyrange => $_getN(0);
  @$pb.TagNumber(981486)
  set hashkeyrange(HashKeyRange value) => $_setField(981486, value);
  @$pb.TagNumber(981486)
  $core.bool hasHashkeyrange() => $_has(0);
  @$pb.TagNumber(981486)
  void clearHashkeyrange() => $_clearField(981486);
  @$pb.TagNumber(981486)
  HashKeyRange ensureHashkeyrange() => $_ensure(0);

  @$pb.TagNumber(66410951)
  $core.String get shardid => $_getSZ(1);
  @$pb.TagNumber(66410951)
  set shardid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(66410951)
  $core.bool hasShardid() => $_has(1);
  @$pb.TagNumber(66410951)
  void clearShardid() => $_clearField(66410951);

  @$pb.TagNumber(86578567)
  SequenceNumberRange get sequencenumberrange => $_getN(2);
  @$pb.TagNumber(86578567)
  set sequencenumberrange(SequenceNumberRange value) =>
      $_setField(86578567, value);
  @$pb.TagNumber(86578567)
  $core.bool hasSequencenumberrange() => $_has(2);
  @$pb.TagNumber(86578567)
  void clearSequencenumberrange() => $_clearField(86578567);
  @$pb.TagNumber(86578567)
  SequenceNumberRange ensureSequencenumberrange() => $_ensure(2);

  @$pb.TagNumber(310461245)
  $core.String get adjacentparentshardid => $_getSZ(3);
  @$pb.TagNumber(310461245)
  set adjacentparentshardid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(310461245)
  $core.bool hasAdjacentparentshardid() => $_has(3);
  @$pb.TagNumber(310461245)
  void clearAdjacentparentshardid() => $_clearField(310461245);

  @$pb.TagNumber(431117103)
  $core.String get parentshardid => $_getSZ(4);
  @$pb.TagNumber(431117103)
  set parentshardid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(431117103)
  $core.bool hasParentshardid() => $_has(4);
  @$pb.TagNumber(431117103)
  void clearParentshardid() => $_clearField(431117103);
}

class ShardFilter extends $pb.GeneratedMessage {
  factory ShardFilter({
    $core.String? shardid,
    $core.String? timestamp,
    ShardFilterType? type,
  }) {
    final result = create();
    if (shardid != null) result.shardid = shardid;
    if (timestamp != null) result.timestamp = timestamp;
    if (type != null) result.type = type;
    return result;
  }

  ShardFilter._();

  factory ShardFilter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ShardFilter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ShardFilter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(66410951, _omitFieldNames ? '' : 'shardid')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aE<ShardFilterType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: ShardFilterType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ShardFilter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ShardFilter copyWith(void Function(ShardFilter) updates) =>
      super.copyWith((message) => updates(message as ShardFilter))
          as ShardFilter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ShardFilter create() => ShardFilter._();
  @$core.override
  ShardFilter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ShardFilter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ShardFilter>(create);
  static ShardFilter? _defaultInstance;

  @$pb.TagNumber(66410951)
  $core.String get shardid => $_getSZ(0);
  @$pb.TagNumber(66410951)
  set shardid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(66410951)
  $core.bool hasShardid() => $_has(0);
  @$pb.TagNumber(66410951)
  void clearShardid() => $_clearField(66410951);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(1);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(1);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(290836590)
  ShardFilterType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(ShardFilterType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class SplitShardInput extends $pb.GeneratedMessage {
  factory SplitShardInput({
    $core.String? shardtosplit,
    $core.String? newstartinghashkey,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (shardtosplit != null) result.shardtosplit = shardtosplit;
    if (newstartinghashkey != null)
      result.newstartinghashkey = newstartinghashkey;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  SplitShardInput._();

  factory SplitShardInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SplitShardInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SplitShardInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(124434671, _omitFieldNames ? '' : 'shardtosplit')
    ..aOS(351052649, _omitFieldNames ? '' : 'newstartinghashkey')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SplitShardInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SplitShardInput copyWith(void Function(SplitShardInput) updates) =>
      super.copyWith((message) => updates(message as SplitShardInput))
          as SplitShardInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SplitShardInput create() => SplitShardInput._();
  @$core.override
  SplitShardInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SplitShardInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SplitShardInput>(create);
  static SplitShardInput? _defaultInstance;

  @$pb.TagNumber(124434671)
  $core.String get shardtosplit => $_getSZ(0);
  @$pb.TagNumber(124434671)
  set shardtosplit($core.String value) => $_setString(0, value);
  @$pb.TagNumber(124434671)
  $core.bool hasShardtosplit() => $_has(0);
  @$pb.TagNumber(124434671)
  void clearShardtosplit() => $_clearField(124434671);

  @$pb.TagNumber(351052649)
  $core.String get newstartinghashkey => $_getSZ(1);
  @$pb.TagNumber(351052649)
  set newstartinghashkey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(351052649)
  $core.bool hasNewstartinghashkey() => $_has(1);
  @$pb.TagNumber(351052649)
  void clearNewstartinghashkey() => $_clearField(351052649);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(3);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(3);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class StartStreamEncryptionInput extends $pb.GeneratedMessage {
  factory StartStreamEncryptionInput({
    EncryptionType? encryptiontype,
    $core.String? keyid,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (encryptiontype != null) result.encryptiontype = encryptiontype;
    if (keyid != null) result.keyid = keyid;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  StartStreamEncryptionInput._();

  factory StartStreamEncryptionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartStreamEncryptionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartStreamEncryptionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<EncryptionType>(264007605, _omitFieldNames ? '' : 'encryptiontype',
        enumValues: EncryptionType.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartStreamEncryptionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartStreamEncryptionInput copyWith(
          void Function(StartStreamEncryptionInput) updates) =>
      super.copyWith(
              (message) => updates(message as StartStreamEncryptionInput))
          as StartStreamEncryptionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartStreamEncryptionInput create() => StartStreamEncryptionInput._();
  @$core.override
  StartStreamEncryptionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartStreamEncryptionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartStreamEncryptionInput>(create);
  static StartStreamEncryptionInput? _defaultInstance;

  @$pb.TagNumber(264007605)
  EncryptionType get encryptiontype => $_getN(0);
  @$pb.TagNumber(264007605)
  set encryptiontype(EncryptionType value) => $_setField(264007605, value);
  @$pb.TagNumber(264007605)
  $core.bool hasEncryptiontype() => $_has(0);
  @$pb.TagNumber(264007605)
  void clearEncryptiontype() => $_clearField(264007605);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(3);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(3);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class StartingPosition extends $pb.GeneratedMessage {
  factory StartingPosition({
    $core.String? sequencenumber,
    $core.String? timestamp,
    ShardIteratorType? type,
  }) {
    final result = create();
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (timestamp != null) result.timestamp = timestamp;
    if (type != null) result.type = type;
    return result;
  }

  StartingPosition._();

  factory StartingPosition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartingPosition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartingPosition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aE<ShardIteratorType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: ShardIteratorType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartingPosition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartingPosition copyWith(void Function(StartingPosition) updates) =>
      super.copyWith((message) => updates(message as StartingPosition))
          as StartingPosition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartingPosition create() => StartingPosition._();
  @$core.override
  StartingPosition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartingPosition getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartingPosition>(create);
  static StartingPosition? _defaultInstance;

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(0);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(0);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(1);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(1);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(290836590)
  ShardIteratorType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(ShardIteratorType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class StopStreamEncryptionInput extends $pb.GeneratedMessage {
  factory StopStreamEncryptionInput({
    EncryptionType? encryptiontype,
    $core.String? keyid,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (encryptiontype != null) result.encryptiontype = encryptiontype;
    if (keyid != null) result.keyid = keyid;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  StopStreamEncryptionInput._();

  factory StopStreamEncryptionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopStreamEncryptionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopStreamEncryptionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<EncryptionType>(264007605, _omitFieldNames ? '' : 'encryptiontype',
        enumValues: EncryptionType.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopStreamEncryptionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopStreamEncryptionInput copyWith(
          void Function(StopStreamEncryptionInput) updates) =>
      super.copyWith((message) => updates(message as StopStreamEncryptionInput))
          as StopStreamEncryptionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopStreamEncryptionInput create() => StopStreamEncryptionInput._();
  @$core.override
  StopStreamEncryptionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopStreamEncryptionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopStreamEncryptionInput>(create);
  static StopStreamEncryptionInput? _defaultInstance;

  @$pb.TagNumber(264007605)
  EncryptionType get encryptiontype => $_getN(0);
  @$pb.TagNumber(264007605)
  set encryptiontype(EncryptionType value) => $_setField(264007605, value);
  @$pb.TagNumber(264007605)
  $core.bool hasEncryptiontype() => $_has(0);
  @$pb.TagNumber(264007605)
  void clearEncryptiontype() => $_clearField(264007605);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(3);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(3);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class StreamDescription extends $pb.GeneratedMessage {
  factory StreamDescription({
    $core.bool? hasmoreshards,
    StreamModeDetails? streammodedetails,
    $core.String? streamcreationtimestamp,
    StreamStatus? streamstatus,
    EncryptionType? encryptiontype,
    $core.String? keyid,
    $core.int? retentionperiodhours,
    $core.Iterable<Shard>? shards,
    $core.Iterable<EnhancedMetrics>? enhancedmonitoring,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (hasmoreshards != null) result.hasmoreshards = hasmoreshards;
    if (streammodedetails != null) result.streammodedetails = streammodedetails;
    if (streamcreationtimestamp != null)
      result.streamcreationtimestamp = streamcreationtimestamp;
    if (streamstatus != null) result.streamstatus = streamstatus;
    if (encryptiontype != null) result.encryptiontype = encryptiontype;
    if (keyid != null) result.keyid = keyid;
    if (retentionperiodhours != null)
      result.retentionperiodhours = retentionperiodhours;
    if (shards != null) result.shards.addAll(shards);
    if (enhancedmonitoring != null)
      result.enhancedmonitoring.addAll(enhancedmonitoring);
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  StreamDescription._();

  factory StreamDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StreamDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOB(10836604, _omitFieldNames ? '' : 'hasmoreshards')
    ..aOM<StreamModeDetails>(
        12139665, _omitFieldNames ? '' : 'streammodedetails',
        subBuilder: StreamModeDetails.create)
    ..aOS(224951013, _omitFieldNames ? '' : 'streamcreationtimestamp')
    ..aE<StreamStatus>(245792976, _omitFieldNames ? '' : 'streamstatus',
        enumValues: StreamStatus.values)
    ..aE<EncryptionType>(264007605, _omitFieldNames ? '' : 'encryptiontype',
        enumValues: EncryptionType.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(396381944, _omitFieldNames ? '' : 'retentionperiodhours')
    ..pPM<Shard>(437117641, _omitFieldNames ? '' : 'shards',
        subBuilder: Shard.create)
    ..pPM<EnhancedMetrics>(
        452259826, _omitFieldNames ? '' : 'enhancedmonitoring',
        subBuilder: EnhancedMetrics.create)
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDescription copyWith(void Function(StreamDescription) updates) =>
      super.copyWith((message) => updates(message as StreamDescription))
          as StreamDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StreamDescription create() => StreamDescription._();
  @$core.override
  StreamDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StreamDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamDescription>(create);
  static StreamDescription? _defaultInstance;

  @$pb.TagNumber(10836604)
  $core.bool get hasmoreshards => $_getBF(0);
  @$pb.TagNumber(10836604)
  set hasmoreshards($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(10836604)
  $core.bool hasHasmoreshards() => $_has(0);
  @$pb.TagNumber(10836604)
  void clearHasmoreshards() => $_clearField(10836604);

  @$pb.TagNumber(12139665)
  StreamModeDetails get streammodedetails => $_getN(1);
  @$pb.TagNumber(12139665)
  set streammodedetails(StreamModeDetails value) => $_setField(12139665, value);
  @$pb.TagNumber(12139665)
  $core.bool hasStreammodedetails() => $_has(1);
  @$pb.TagNumber(12139665)
  void clearStreammodedetails() => $_clearField(12139665);
  @$pb.TagNumber(12139665)
  StreamModeDetails ensureStreammodedetails() => $_ensure(1);

  @$pb.TagNumber(224951013)
  $core.String get streamcreationtimestamp => $_getSZ(2);
  @$pb.TagNumber(224951013)
  set streamcreationtimestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(224951013)
  $core.bool hasStreamcreationtimestamp() => $_has(2);
  @$pb.TagNumber(224951013)
  void clearStreamcreationtimestamp() => $_clearField(224951013);

  @$pb.TagNumber(245792976)
  StreamStatus get streamstatus => $_getN(3);
  @$pb.TagNumber(245792976)
  set streamstatus(StreamStatus value) => $_setField(245792976, value);
  @$pb.TagNumber(245792976)
  $core.bool hasStreamstatus() => $_has(3);
  @$pb.TagNumber(245792976)
  void clearStreamstatus() => $_clearField(245792976);

  @$pb.TagNumber(264007605)
  EncryptionType get encryptiontype => $_getN(4);
  @$pb.TagNumber(264007605)
  set encryptiontype(EncryptionType value) => $_setField(264007605, value);
  @$pb.TagNumber(264007605)
  $core.bool hasEncryptiontype() => $_has(4);
  @$pb.TagNumber(264007605)
  void clearEncryptiontype() => $_clearField(264007605);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(5);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(5);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(396381944)
  $core.int get retentionperiodhours => $_getIZ(6);
  @$pb.TagNumber(396381944)
  set retentionperiodhours($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(396381944)
  $core.bool hasRetentionperiodhours() => $_has(6);
  @$pb.TagNumber(396381944)
  void clearRetentionperiodhours() => $_clearField(396381944);

  @$pb.TagNumber(437117641)
  $pb.PbList<Shard> get shards => $_getList(7);

  @$pb.TagNumber(452259826)
  $pb.PbList<EnhancedMetrics> get enhancedmonitoring => $_getList(8);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(9);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(9, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(9);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(10);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(10, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(10);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class StreamDescriptionSummary extends $pb.GeneratedMessage {
  factory StreamDescriptionSummary({
    StreamModeDetails? streammodedetails,
    $core.int? maxrecordsizeinkib,
    $core.String? streamcreationtimestamp,
    StreamStatus? streamstatus,
    EncryptionType? encryptiontype,
    $core.String? keyid,
    WarmThroughputObject? warmthroughput,
    $core.int? retentionperiodhours,
    $core.String? streamid,
    $core.int? consumercount,
    $core.Iterable<EnhancedMetrics>? enhancedmonitoring,
    $core.String? streamname,
    $core.int? openshardcount,
    $core.String? streamarn,
  }) {
    final result = create();
    if (streammodedetails != null) result.streammodedetails = streammodedetails;
    if (maxrecordsizeinkib != null)
      result.maxrecordsizeinkib = maxrecordsizeinkib;
    if (streamcreationtimestamp != null)
      result.streamcreationtimestamp = streamcreationtimestamp;
    if (streamstatus != null) result.streamstatus = streamstatus;
    if (encryptiontype != null) result.encryptiontype = encryptiontype;
    if (keyid != null) result.keyid = keyid;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (retentionperiodhours != null)
      result.retentionperiodhours = retentionperiodhours;
    if (streamid != null) result.streamid = streamid;
    if (consumercount != null) result.consumercount = consumercount;
    if (enhancedmonitoring != null)
      result.enhancedmonitoring.addAll(enhancedmonitoring);
    if (streamname != null) result.streamname = streamname;
    if (openshardcount != null) result.openshardcount = openshardcount;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  StreamDescriptionSummary._();

  factory StreamDescriptionSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StreamDescriptionSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamDescriptionSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<StreamModeDetails>(
        12139665, _omitFieldNames ? '' : 'streammodedetails',
        subBuilder: StreamModeDetails.create)
    ..aI(197267253, _omitFieldNames ? '' : 'maxrecordsizeinkib')
    ..aOS(224951013, _omitFieldNames ? '' : 'streamcreationtimestamp')
    ..aE<StreamStatus>(245792976, _omitFieldNames ? '' : 'streamstatus',
        enumValues: StreamStatus.values)
    ..aE<EncryptionType>(264007605, _omitFieldNames ? '' : 'encryptiontype',
        enumValues: EncryptionType.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOM<WarmThroughputObject>(
        290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: WarmThroughputObject.create)
    ..aI(396381944, _omitFieldNames ? '' : 'retentionperiodhours')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aI(448084721, _omitFieldNames ? '' : 'consumercount')
    ..pPM<EnhancedMetrics>(
        452259826, _omitFieldNames ? '' : 'enhancedmonitoring',
        subBuilder: EnhancedMetrics.create)
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aI(476409287, _omitFieldNames ? '' : 'openshardcount')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDescriptionSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamDescriptionSummary copyWith(
          void Function(StreamDescriptionSummary) updates) =>
      super.copyWith((message) => updates(message as StreamDescriptionSummary))
          as StreamDescriptionSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StreamDescriptionSummary create() => StreamDescriptionSummary._();
  @$core.override
  StreamDescriptionSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StreamDescriptionSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamDescriptionSummary>(create);
  static StreamDescriptionSummary? _defaultInstance;

  @$pb.TagNumber(12139665)
  StreamModeDetails get streammodedetails => $_getN(0);
  @$pb.TagNumber(12139665)
  set streammodedetails(StreamModeDetails value) => $_setField(12139665, value);
  @$pb.TagNumber(12139665)
  $core.bool hasStreammodedetails() => $_has(0);
  @$pb.TagNumber(12139665)
  void clearStreammodedetails() => $_clearField(12139665);
  @$pb.TagNumber(12139665)
  StreamModeDetails ensureStreammodedetails() => $_ensure(0);

  @$pb.TagNumber(197267253)
  $core.int get maxrecordsizeinkib => $_getIZ(1);
  @$pb.TagNumber(197267253)
  set maxrecordsizeinkib($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(197267253)
  $core.bool hasMaxrecordsizeinkib() => $_has(1);
  @$pb.TagNumber(197267253)
  void clearMaxrecordsizeinkib() => $_clearField(197267253);

  @$pb.TagNumber(224951013)
  $core.String get streamcreationtimestamp => $_getSZ(2);
  @$pb.TagNumber(224951013)
  set streamcreationtimestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(224951013)
  $core.bool hasStreamcreationtimestamp() => $_has(2);
  @$pb.TagNumber(224951013)
  void clearStreamcreationtimestamp() => $_clearField(224951013);

  @$pb.TagNumber(245792976)
  StreamStatus get streamstatus => $_getN(3);
  @$pb.TagNumber(245792976)
  set streamstatus(StreamStatus value) => $_setField(245792976, value);
  @$pb.TagNumber(245792976)
  $core.bool hasStreamstatus() => $_has(3);
  @$pb.TagNumber(245792976)
  void clearStreamstatus() => $_clearField(245792976);

  @$pb.TagNumber(264007605)
  EncryptionType get encryptiontype => $_getN(4);
  @$pb.TagNumber(264007605)
  set encryptiontype(EncryptionType value) => $_setField(264007605, value);
  @$pb.TagNumber(264007605)
  $core.bool hasEncryptiontype() => $_has(4);
  @$pb.TagNumber(264007605)
  void clearEncryptiontype() => $_clearField(264007605);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(5);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(5);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(290598659)
  WarmThroughputObject get warmthroughput => $_getN(6);
  @$pb.TagNumber(290598659)
  set warmthroughput(WarmThroughputObject value) =>
      $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(6);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  WarmThroughputObject ensureWarmthroughput() => $_ensure(6);

  @$pb.TagNumber(396381944)
  $core.int get retentionperiodhours => $_getIZ(7);
  @$pb.TagNumber(396381944)
  set retentionperiodhours($core.int value) => $_setSignedInt32(7, value);
  @$pb.TagNumber(396381944)
  $core.bool hasRetentionperiodhours() => $_has(7);
  @$pb.TagNumber(396381944)
  void clearRetentionperiodhours() => $_clearField(396381944);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(8);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(8);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(448084721)
  $core.int get consumercount => $_getIZ(9);
  @$pb.TagNumber(448084721)
  set consumercount($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(448084721)
  $core.bool hasConsumercount() => $_has(9);
  @$pb.TagNumber(448084721)
  void clearConsumercount() => $_clearField(448084721);

  @$pb.TagNumber(452259826)
  $pb.PbList<EnhancedMetrics> get enhancedmonitoring => $_getList(10);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(11);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(11, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(11);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(476409287)
  $core.int get openshardcount => $_getIZ(12);
  @$pb.TagNumber(476409287)
  set openshardcount($core.int value) => $_setSignedInt32(12, value);
  @$pb.TagNumber(476409287)
  $core.bool hasOpenshardcount() => $_has(12);
  @$pb.TagNumber(476409287)
  void clearOpenshardcount() => $_clearField(476409287);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(13);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(13, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(13);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class StreamModeDetails extends $pb.GeneratedMessage {
  factory StreamModeDetails({
    StreamMode? streammode,
  }) {
    final result = create();
    if (streammode != null) result.streammode = streammode;
    return result;
  }

  StreamModeDetails._();

  factory StreamModeDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StreamModeDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamModeDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<StreamMode>(457304819, _omitFieldNames ? '' : 'streammode',
        enumValues: StreamMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamModeDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamModeDetails copyWith(void Function(StreamModeDetails) updates) =>
      super.copyWith((message) => updates(message as StreamModeDetails))
          as StreamModeDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StreamModeDetails create() => StreamModeDetails._();
  @$core.override
  StreamModeDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StreamModeDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamModeDetails>(create);
  static StreamModeDetails? _defaultInstance;

  @$pb.TagNumber(457304819)
  StreamMode get streammode => $_getN(0);
  @$pb.TagNumber(457304819)
  set streammode(StreamMode value) => $_setField(457304819, value);
  @$pb.TagNumber(457304819)
  $core.bool hasStreammode() => $_has(0);
  @$pb.TagNumber(457304819)
  void clearStreammode() => $_clearField(457304819);
}

class StreamSummary extends $pb.GeneratedMessage {
  factory StreamSummary({
    StreamModeDetails? streammodedetails,
    $core.String? streamcreationtimestamp,
    StreamStatus? streamstatus,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (streammodedetails != null) result.streammodedetails = streammodedetails;
    if (streamcreationtimestamp != null)
      result.streamcreationtimestamp = streamcreationtimestamp;
    if (streamstatus != null) result.streamstatus = streamstatus;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  StreamSummary._();

  factory StreamSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StreamSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<StreamModeDetails>(
        12139665, _omitFieldNames ? '' : 'streammodedetails',
        subBuilder: StreamModeDetails.create)
    ..aOS(224951013, _omitFieldNames ? '' : 'streamcreationtimestamp')
    ..aE<StreamStatus>(245792976, _omitFieldNames ? '' : 'streamstatus',
        enumValues: StreamStatus.values)
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSummary copyWith(void Function(StreamSummary) updates) =>
      super.copyWith((message) => updates(message as StreamSummary))
          as StreamSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StreamSummary create() => StreamSummary._();
  @$core.override
  StreamSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StreamSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamSummary>(create);
  static StreamSummary? _defaultInstance;

  @$pb.TagNumber(12139665)
  StreamModeDetails get streammodedetails => $_getN(0);
  @$pb.TagNumber(12139665)
  set streammodedetails(StreamModeDetails value) => $_setField(12139665, value);
  @$pb.TagNumber(12139665)
  $core.bool hasStreammodedetails() => $_has(0);
  @$pb.TagNumber(12139665)
  void clearStreammodedetails() => $_clearField(12139665);
  @$pb.TagNumber(12139665)
  StreamModeDetails ensureStreammodedetails() => $_ensure(0);

  @$pb.TagNumber(224951013)
  $core.String get streamcreationtimestamp => $_getSZ(1);
  @$pb.TagNumber(224951013)
  set streamcreationtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(224951013)
  $core.bool hasStreamcreationtimestamp() => $_has(1);
  @$pb.TagNumber(224951013)
  void clearStreamcreationtimestamp() => $_clearField(224951013);

  @$pb.TagNumber(245792976)
  StreamStatus get streamstatus => $_getN(2);
  @$pb.TagNumber(245792976)
  set streamstatus(StreamStatus value) => $_setField(245792976, value);
  @$pb.TagNumber(245792976)
  $core.bool hasStreamstatus() => $_has(2);
  @$pb.TagNumber(245792976)
  void clearStreamstatus() => $_clearField(245792976);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(3);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(3);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class SubscribeToShardEvent extends $pb.GeneratedMessage {
  factory SubscribeToShardEvent({
    $core.String? continuationsequencenumber,
    $core.Iterable<ChildShard>? childshards,
    $core.Iterable<Record>? records,
    $fixnum.Int64? millisbehindlatest,
  }) {
    final result = create();
    if (continuationsequencenumber != null)
      result.continuationsequencenumber = continuationsequencenumber;
    if (childshards != null) result.childshards.addAll(childshards);
    if (records != null) result.records.addAll(records);
    if (millisbehindlatest != null)
      result.millisbehindlatest = millisbehindlatest;
    return result;
  }

  SubscribeToShardEvent._();

  factory SubscribeToShardEvent.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscribeToShardEvent.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeToShardEvent',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(143046235, _omitFieldNames ? '' : 'continuationsequencenumber')
    ..pPM<ChildShard>(348740657, _omitFieldNames ? '' : 'childshards',
        subBuilder: ChildShard.create)
    ..pPM<Record>(423557454, _omitFieldNames ? '' : 'records',
        subBuilder: Record.create)
    ..aInt64(456422057, _omitFieldNames ? '' : 'millisbehindlatest')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardEvent clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardEvent copyWith(
          void Function(SubscribeToShardEvent) updates) =>
      super.copyWith((message) => updates(message as SubscribeToShardEvent))
          as SubscribeToShardEvent;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscribeToShardEvent create() => SubscribeToShardEvent._();
  @$core.override
  SubscribeToShardEvent createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscribeToShardEvent getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeToShardEvent>(create);
  static SubscribeToShardEvent? _defaultInstance;

  @$pb.TagNumber(143046235)
  $core.String get continuationsequencenumber => $_getSZ(0);
  @$pb.TagNumber(143046235)
  set continuationsequencenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(143046235)
  $core.bool hasContinuationsequencenumber() => $_has(0);
  @$pb.TagNumber(143046235)
  void clearContinuationsequencenumber() => $_clearField(143046235);

  @$pb.TagNumber(348740657)
  $pb.PbList<ChildShard> get childshards => $_getList(1);

  @$pb.TagNumber(423557454)
  $pb.PbList<Record> get records => $_getList(2);

  @$pb.TagNumber(456422057)
  $fixnum.Int64 get millisbehindlatest => $_getI64(3);
  @$pb.TagNumber(456422057)
  set millisbehindlatest($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(456422057)
  $core.bool hasMillisbehindlatest() => $_has(3);
  @$pb.TagNumber(456422057)
  void clearMillisbehindlatest() => $_clearField(456422057);
}

class SubscribeToShardEventStream extends $pb.GeneratedMessage {
  factory SubscribeToShardEventStream({
    KMSOptInRequired? kmsoptinrequired,
    KMSThrottlingException? kmsthrottlingexception,
    KMSAccessDeniedException? kmsaccessdeniedexception,
    KMSInvalidStateException? kmsinvalidstateexception,
    KMSNotFoundException? kmsnotfoundexception,
    SubscribeToShardEvent? subscribetoshardevent,
    KMSDisabledException? kmsdisabledexception,
    ResourceInUseException? resourceinuseexception,
    ResourceNotFoundException? resourcenotfoundexception,
    InternalFailureException? internalfailureexception,
  }) {
    final result = create();
    if (kmsoptinrequired != null) result.kmsoptinrequired = kmsoptinrequired;
    if (kmsthrottlingexception != null)
      result.kmsthrottlingexception = kmsthrottlingexception;
    if (kmsaccessdeniedexception != null)
      result.kmsaccessdeniedexception = kmsaccessdeniedexception;
    if (kmsinvalidstateexception != null)
      result.kmsinvalidstateexception = kmsinvalidstateexception;
    if (kmsnotfoundexception != null)
      result.kmsnotfoundexception = kmsnotfoundexception;
    if (subscribetoshardevent != null)
      result.subscribetoshardevent = subscribetoshardevent;
    if (kmsdisabledexception != null)
      result.kmsdisabledexception = kmsdisabledexception;
    if (resourceinuseexception != null)
      result.resourceinuseexception = resourceinuseexception;
    if (resourcenotfoundexception != null)
      result.resourcenotfoundexception = resourcenotfoundexception;
    if (internalfailureexception != null)
      result.internalfailureexception = internalfailureexception;
    return result;
  }

  SubscribeToShardEventStream._();

  factory SubscribeToShardEventStream.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscribeToShardEventStream.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeToShardEventStream',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<KMSOptInRequired>(72448154, _omitFieldNames ? '' : 'kmsoptinrequired',
        subBuilder: KMSOptInRequired.create)
    ..aOM<KMSThrottlingException>(
        109345891, _omitFieldNames ? '' : 'kmsthrottlingexception',
        subBuilder: KMSThrottlingException.create)
    ..aOM<KMSAccessDeniedException>(
        110824201, _omitFieldNames ? '' : 'kmsaccessdeniedexception',
        subBuilder: KMSAccessDeniedException.create)
    ..aOM<KMSInvalidStateException>(
        183514916, _omitFieldNames ? '' : 'kmsinvalidstateexception',
        subBuilder: KMSInvalidStateException.create)
    ..aOM<KMSNotFoundException>(
        251973779, _omitFieldNames ? '' : 'kmsnotfoundexception',
        subBuilder: KMSNotFoundException.create)
    ..aOM<SubscribeToShardEvent>(
        309751871, _omitFieldNames ? '' : 'subscribetoshardevent',
        subBuilder: SubscribeToShardEvent.create)
    ..aOM<KMSDisabledException>(
        352104756, _omitFieldNames ? '' : 'kmsdisabledexception',
        subBuilder: KMSDisabledException.create)
    ..aOM<ResourceInUseException>(
        360549483, _omitFieldNames ? '' : 'resourceinuseexception',
        subBuilder: ResourceInUseException.create)
    ..aOM<ResourceNotFoundException>(
        396489184, _omitFieldNames ? '' : 'resourcenotfoundexception',
        subBuilder: ResourceNotFoundException.create)
    ..aOM<InternalFailureException>(
        445206748, _omitFieldNames ? '' : 'internalfailureexception',
        subBuilder: InternalFailureException.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardEventStream clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardEventStream copyWith(
          void Function(SubscribeToShardEventStream) updates) =>
      super.copyWith(
              (message) => updates(message as SubscribeToShardEventStream))
          as SubscribeToShardEventStream;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscribeToShardEventStream create() =>
      SubscribeToShardEventStream._();
  @$core.override
  SubscribeToShardEventStream createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscribeToShardEventStream getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeToShardEventStream>(create);
  static SubscribeToShardEventStream? _defaultInstance;

  @$pb.TagNumber(72448154)
  KMSOptInRequired get kmsoptinrequired => $_getN(0);
  @$pb.TagNumber(72448154)
  set kmsoptinrequired(KMSOptInRequired value) => $_setField(72448154, value);
  @$pb.TagNumber(72448154)
  $core.bool hasKmsoptinrequired() => $_has(0);
  @$pb.TagNumber(72448154)
  void clearKmsoptinrequired() => $_clearField(72448154);
  @$pb.TagNumber(72448154)
  KMSOptInRequired ensureKmsoptinrequired() => $_ensure(0);

  @$pb.TagNumber(109345891)
  KMSThrottlingException get kmsthrottlingexception => $_getN(1);
  @$pb.TagNumber(109345891)
  set kmsthrottlingexception(KMSThrottlingException value) =>
      $_setField(109345891, value);
  @$pb.TagNumber(109345891)
  $core.bool hasKmsthrottlingexception() => $_has(1);
  @$pb.TagNumber(109345891)
  void clearKmsthrottlingexception() => $_clearField(109345891);
  @$pb.TagNumber(109345891)
  KMSThrottlingException ensureKmsthrottlingexception() => $_ensure(1);

  @$pb.TagNumber(110824201)
  KMSAccessDeniedException get kmsaccessdeniedexception => $_getN(2);
  @$pb.TagNumber(110824201)
  set kmsaccessdeniedexception(KMSAccessDeniedException value) =>
      $_setField(110824201, value);
  @$pb.TagNumber(110824201)
  $core.bool hasKmsaccessdeniedexception() => $_has(2);
  @$pb.TagNumber(110824201)
  void clearKmsaccessdeniedexception() => $_clearField(110824201);
  @$pb.TagNumber(110824201)
  KMSAccessDeniedException ensureKmsaccessdeniedexception() => $_ensure(2);

  @$pb.TagNumber(183514916)
  KMSInvalidStateException get kmsinvalidstateexception => $_getN(3);
  @$pb.TagNumber(183514916)
  set kmsinvalidstateexception(KMSInvalidStateException value) =>
      $_setField(183514916, value);
  @$pb.TagNumber(183514916)
  $core.bool hasKmsinvalidstateexception() => $_has(3);
  @$pb.TagNumber(183514916)
  void clearKmsinvalidstateexception() => $_clearField(183514916);
  @$pb.TagNumber(183514916)
  KMSInvalidStateException ensureKmsinvalidstateexception() => $_ensure(3);

  @$pb.TagNumber(251973779)
  KMSNotFoundException get kmsnotfoundexception => $_getN(4);
  @$pb.TagNumber(251973779)
  set kmsnotfoundexception(KMSNotFoundException value) =>
      $_setField(251973779, value);
  @$pb.TagNumber(251973779)
  $core.bool hasKmsnotfoundexception() => $_has(4);
  @$pb.TagNumber(251973779)
  void clearKmsnotfoundexception() => $_clearField(251973779);
  @$pb.TagNumber(251973779)
  KMSNotFoundException ensureKmsnotfoundexception() => $_ensure(4);

  @$pb.TagNumber(309751871)
  SubscribeToShardEvent get subscribetoshardevent => $_getN(5);
  @$pb.TagNumber(309751871)
  set subscribetoshardevent(SubscribeToShardEvent value) =>
      $_setField(309751871, value);
  @$pb.TagNumber(309751871)
  $core.bool hasSubscribetoshardevent() => $_has(5);
  @$pb.TagNumber(309751871)
  void clearSubscribetoshardevent() => $_clearField(309751871);
  @$pb.TagNumber(309751871)
  SubscribeToShardEvent ensureSubscribetoshardevent() => $_ensure(5);

  @$pb.TagNumber(352104756)
  KMSDisabledException get kmsdisabledexception => $_getN(6);
  @$pb.TagNumber(352104756)
  set kmsdisabledexception(KMSDisabledException value) =>
      $_setField(352104756, value);
  @$pb.TagNumber(352104756)
  $core.bool hasKmsdisabledexception() => $_has(6);
  @$pb.TagNumber(352104756)
  void clearKmsdisabledexception() => $_clearField(352104756);
  @$pb.TagNumber(352104756)
  KMSDisabledException ensureKmsdisabledexception() => $_ensure(6);

  @$pb.TagNumber(360549483)
  ResourceInUseException get resourceinuseexception => $_getN(7);
  @$pb.TagNumber(360549483)
  set resourceinuseexception(ResourceInUseException value) =>
      $_setField(360549483, value);
  @$pb.TagNumber(360549483)
  $core.bool hasResourceinuseexception() => $_has(7);
  @$pb.TagNumber(360549483)
  void clearResourceinuseexception() => $_clearField(360549483);
  @$pb.TagNumber(360549483)
  ResourceInUseException ensureResourceinuseexception() => $_ensure(7);

  @$pb.TagNumber(396489184)
  ResourceNotFoundException get resourcenotfoundexception => $_getN(8);
  @$pb.TagNumber(396489184)
  set resourcenotfoundexception(ResourceNotFoundException value) =>
      $_setField(396489184, value);
  @$pb.TagNumber(396489184)
  $core.bool hasResourcenotfoundexception() => $_has(8);
  @$pb.TagNumber(396489184)
  void clearResourcenotfoundexception() => $_clearField(396489184);
  @$pb.TagNumber(396489184)
  ResourceNotFoundException ensureResourcenotfoundexception() => $_ensure(8);

  @$pb.TagNumber(445206748)
  InternalFailureException get internalfailureexception => $_getN(9);
  @$pb.TagNumber(445206748)
  set internalfailureexception(InternalFailureException value) =>
      $_setField(445206748, value);
  @$pb.TagNumber(445206748)
  $core.bool hasInternalfailureexception() => $_has(9);
  @$pb.TagNumber(445206748)
  void clearInternalfailureexception() => $_clearField(445206748);
  @$pb.TagNumber(445206748)
  InternalFailureException ensureInternalfailureexception() => $_ensure(9);
}

class SubscribeToShardInput extends $pb.GeneratedMessage {
  factory SubscribeToShardInput({
    $core.String? consumerarn,
    $core.String? shardid,
    $core.String? streamid,
    StartingPosition? startingposition,
  }) {
    final result = create();
    if (consumerarn != null) result.consumerarn = consumerarn;
    if (shardid != null) result.shardid = shardid;
    if (streamid != null) result.streamid = streamid;
    if (startingposition != null) result.startingposition = startingposition;
    return result;
  }

  SubscribeToShardInput._();

  factory SubscribeToShardInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscribeToShardInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeToShardInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(41107441, _omitFieldNames ? '' : 'consumerarn')
    ..aOS(66410951, _omitFieldNames ? '' : 'shardid')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOM<StartingPosition>(
        428771919, _omitFieldNames ? '' : 'startingposition',
        subBuilder: StartingPosition.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardInput copyWith(
          void Function(SubscribeToShardInput) updates) =>
      super.copyWith((message) => updates(message as SubscribeToShardInput))
          as SubscribeToShardInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscribeToShardInput create() => SubscribeToShardInput._();
  @$core.override
  SubscribeToShardInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscribeToShardInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeToShardInput>(create);
  static SubscribeToShardInput? _defaultInstance;

  @$pb.TagNumber(41107441)
  $core.String get consumerarn => $_getSZ(0);
  @$pb.TagNumber(41107441)
  set consumerarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41107441)
  $core.bool hasConsumerarn() => $_has(0);
  @$pb.TagNumber(41107441)
  void clearConsumerarn() => $_clearField(41107441);

  @$pb.TagNumber(66410951)
  $core.String get shardid => $_getSZ(1);
  @$pb.TagNumber(66410951)
  set shardid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(66410951)
  $core.bool hasShardid() => $_has(1);
  @$pb.TagNumber(66410951)
  void clearShardid() => $_clearField(66410951);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(428771919)
  StartingPosition get startingposition => $_getN(3);
  @$pb.TagNumber(428771919)
  set startingposition(StartingPosition value) => $_setField(428771919, value);
  @$pb.TagNumber(428771919)
  $core.bool hasStartingposition() => $_has(3);
  @$pb.TagNumber(428771919)
  void clearStartingposition() => $_clearField(428771919);
  @$pb.TagNumber(428771919)
  StartingPosition ensureStartingposition() => $_ensure(3);
}

class SubscribeToShardOutput extends $pb.GeneratedMessage {
  factory SubscribeToShardOutput({
    SubscribeToShardEventStream? eventstream,
  }) {
    final result = create();
    if (eventstream != null) result.eventstream = eventstream;
    return result;
  }

  SubscribeToShardOutput._();

  factory SubscribeToShardOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscribeToShardOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeToShardOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<SubscribeToShardEventStream>(
        26857888, _omitFieldNames ? '' : 'eventstream',
        subBuilder: SubscribeToShardEventStream.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeToShardOutput copyWith(
          void Function(SubscribeToShardOutput) updates) =>
      super.copyWith((message) => updates(message as SubscribeToShardOutput))
          as SubscribeToShardOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscribeToShardOutput create() => SubscribeToShardOutput._();
  @$core.override
  SubscribeToShardOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscribeToShardOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeToShardOutput>(create);
  static SubscribeToShardOutput? _defaultInstance;

  @$pb.TagNumber(26857888)
  SubscribeToShardEventStream get eventstream => $_getN(0);
  @$pb.TagNumber(26857888)
  set eventstream(SubscribeToShardEventStream value) =>
      $_setField(26857888, value);
  @$pb.TagNumber(26857888)
  $core.bool hasEventstream() => $_has(0);
  @$pb.TagNumber(26857888)
  void clearEventstream() => $_clearField(26857888);
  @$pb.TagNumber(26857888)
  SubscribeToShardEventStream ensureEventstream() => $_ensure(0);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
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
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? streamid,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tags != null) result.tags.addEntries(tags);
    if (streamid != null) result.streamid = streamid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'TagResourceInput.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kinesis'))
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(1);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);
}

class UntagResourceInput extends $pb.GeneratedMessage {
  factory UntagResourceInput({
    $core.Iterable<$core.String>? tagkeys,
    $core.String? resourcearn,
    $core.String? streamid,
  }) {
    final result = create();
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (streamid != null) result.streamid = streamid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);
}

class UpdateAccountSettingsInput extends $pb.GeneratedMessage {
  factory UpdateAccountSettingsInput({
    MinimumThroughputBillingCommitmentInput? minimumthroughputbillingcommitment,
  }) {
    final result = create();
    if (minimumthroughputbillingcommitment != null)
      result.minimumthroughputbillingcommitment =
          minimumthroughputbillingcommitment;
    return result;
  }

  UpdateAccountSettingsInput._();

  factory UpdateAccountSettingsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateAccountSettingsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAccountSettingsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<MinimumThroughputBillingCommitmentInput>(
        335326962, _omitFieldNames ? '' : 'minimumthroughputbillingcommitment',
        subBuilder: MinimumThroughputBillingCommitmentInput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsInput copyWith(
          void Function(UpdateAccountSettingsInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateAccountSettingsInput))
          as UpdateAccountSettingsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsInput create() => UpdateAccountSettingsInput._();
  @$core.override
  UpdateAccountSettingsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAccountSettingsInput>(create);
  static UpdateAccountSettingsInput? _defaultInstance;

  @$pb.TagNumber(335326962)
  MinimumThroughputBillingCommitmentInput
      get minimumthroughputbillingcommitment => $_getN(0);
  @$pb.TagNumber(335326962)
  set minimumthroughputbillingcommitment(
          MinimumThroughputBillingCommitmentInput value) =>
      $_setField(335326962, value);
  @$pb.TagNumber(335326962)
  $core.bool hasMinimumthroughputbillingcommitment() => $_has(0);
  @$pb.TagNumber(335326962)
  void clearMinimumthroughputbillingcommitment() => $_clearField(335326962);
  @$pb.TagNumber(335326962)
  MinimumThroughputBillingCommitmentInput
      ensureMinimumthroughputbillingcommitment() => $_ensure(0);
}

class UpdateAccountSettingsOutput extends $pb.GeneratedMessage {
  factory UpdateAccountSettingsOutput({
    MinimumThroughputBillingCommitmentOutput?
        minimumthroughputbillingcommitment,
  }) {
    final result = create();
    if (minimumthroughputbillingcommitment != null)
      result.minimumthroughputbillingcommitment =
          minimumthroughputbillingcommitment;
    return result;
  }

  UpdateAccountSettingsOutput._();

  factory UpdateAccountSettingsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateAccountSettingsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAccountSettingsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<MinimumThroughputBillingCommitmentOutput>(
        335326962, _omitFieldNames ? '' : 'minimumthroughputbillingcommitment',
        subBuilder: MinimumThroughputBillingCommitmentOutput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsOutput copyWith(
          void Function(UpdateAccountSettingsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateAccountSettingsOutput))
          as UpdateAccountSettingsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsOutput create() =>
      UpdateAccountSettingsOutput._();
  @$core.override
  UpdateAccountSettingsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAccountSettingsOutput>(create);
  static UpdateAccountSettingsOutput? _defaultInstance;

  @$pb.TagNumber(335326962)
  MinimumThroughputBillingCommitmentOutput
      get minimumthroughputbillingcommitment => $_getN(0);
  @$pb.TagNumber(335326962)
  set minimumthroughputbillingcommitment(
          MinimumThroughputBillingCommitmentOutput value) =>
      $_setField(335326962, value);
  @$pb.TagNumber(335326962)
  $core.bool hasMinimumthroughputbillingcommitment() => $_has(0);
  @$pb.TagNumber(335326962)
  void clearMinimumthroughputbillingcommitment() => $_clearField(335326962);
  @$pb.TagNumber(335326962)
  MinimumThroughputBillingCommitmentOutput
      ensureMinimumthroughputbillingcommitment() => $_ensure(0);
}

class UpdateMaxRecordSizeInput extends $pb.GeneratedMessage {
  factory UpdateMaxRecordSizeInput({
    $core.int? maxrecordsizeinkib,
    $core.String? streamid,
    $core.String? streamarn,
  }) {
    final result = create();
    if (maxrecordsizeinkib != null)
      result.maxrecordsizeinkib = maxrecordsizeinkib;
    if (streamid != null) result.streamid = streamid;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateMaxRecordSizeInput._();

  factory UpdateMaxRecordSizeInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateMaxRecordSizeInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateMaxRecordSizeInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(197267253, _omitFieldNames ? '' : 'maxrecordsizeinkib')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMaxRecordSizeInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMaxRecordSizeInput copyWith(
          void Function(UpdateMaxRecordSizeInput) updates) =>
      super.copyWith((message) => updates(message as UpdateMaxRecordSizeInput))
          as UpdateMaxRecordSizeInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateMaxRecordSizeInput create() => UpdateMaxRecordSizeInput._();
  @$core.override
  UpdateMaxRecordSizeInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateMaxRecordSizeInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateMaxRecordSizeInput>(create);
  static UpdateMaxRecordSizeInput? _defaultInstance;

  @$pb.TagNumber(197267253)
  $core.int get maxrecordsizeinkib => $_getIZ(0);
  @$pb.TagNumber(197267253)
  set maxrecordsizeinkib($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(197267253)
  $core.bool hasMaxrecordsizeinkib() => $_has(0);
  @$pb.TagNumber(197267253)
  void clearMaxrecordsizeinkib() => $_clearField(197267253);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(2);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(2);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class UpdateShardCountInput extends $pb.GeneratedMessage {
  factory UpdateShardCountInput({
    ScalingType? scalingtype,
    $core.int? targetshardcount,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (scalingtype != null) result.scalingtype = scalingtype;
    if (targetshardcount != null) result.targetshardcount = targetshardcount;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateShardCountInput._();

  factory UpdateShardCountInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateShardCountInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateShardCountInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aE<ScalingType>(34064531, _omitFieldNames ? '' : 'scalingtype',
        enumValues: ScalingType.values)
    ..aI(361168816, _omitFieldNames ? '' : 'targetshardcount')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateShardCountInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateShardCountInput copyWith(
          void Function(UpdateShardCountInput) updates) =>
      super.copyWith((message) => updates(message as UpdateShardCountInput))
          as UpdateShardCountInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateShardCountInput create() => UpdateShardCountInput._();
  @$core.override
  UpdateShardCountInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateShardCountInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateShardCountInput>(create);
  static UpdateShardCountInput? _defaultInstance;

  @$pb.TagNumber(34064531)
  ScalingType get scalingtype => $_getN(0);
  @$pb.TagNumber(34064531)
  set scalingtype(ScalingType value) => $_setField(34064531, value);
  @$pb.TagNumber(34064531)
  $core.bool hasScalingtype() => $_has(0);
  @$pb.TagNumber(34064531)
  void clearScalingtype() => $_clearField(34064531);

  @$pb.TagNumber(361168816)
  $core.int get targetshardcount => $_getIZ(1);
  @$pb.TagNumber(361168816)
  set targetshardcount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(361168816)
  $core.bool hasTargetshardcount() => $_has(1);
  @$pb.TagNumber(361168816)
  void clearTargetshardcount() => $_clearField(361168816);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(3);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(3);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(4);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(4);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class UpdateShardCountOutput extends $pb.GeneratedMessage {
  factory UpdateShardCountOutput({
    $core.int? currentshardcount,
    $core.int? targetshardcount,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (currentshardcount != null) result.currentshardcount = currentshardcount;
    if (targetshardcount != null) result.targetshardcount = targetshardcount;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateShardCountOutput._();

  factory UpdateShardCountOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateShardCountOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateShardCountOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(21258174, _omitFieldNames ? '' : 'currentshardcount')
    ..aI(361168816, _omitFieldNames ? '' : 'targetshardcount')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateShardCountOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateShardCountOutput copyWith(
          void Function(UpdateShardCountOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateShardCountOutput))
          as UpdateShardCountOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateShardCountOutput create() => UpdateShardCountOutput._();
  @$core.override
  UpdateShardCountOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateShardCountOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateShardCountOutput>(create);
  static UpdateShardCountOutput? _defaultInstance;

  @$pb.TagNumber(21258174)
  $core.int get currentshardcount => $_getIZ(0);
  @$pb.TagNumber(21258174)
  set currentshardcount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(21258174)
  $core.bool hasCurrentshardcount() => $_has(0);
  @$pb.TagNumber(21258174)
  void clearCurrentshardcount() => $_clearField(21258174);

  @$pb.TagNumber(361168816)
  $core.int get targetshardcount => $_getIZ(1);
  @$pb.TagNumber(361168816)
  set targetshardcount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(361168816)
  $core.bool hasTargetshardcount() => $_has(1);
  @$pb.TagNumber(361168816)
  void clearTargetshardcount() => $_clearField(361168816);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class UpdateStreamModeInput extends $pb.GeneratedMessage {
  factory UpdateStreamModeInput({
    StreamModeDetails? streammodedetails,
    $core.int? warmthroughputmibps,
    $core.String? streamid,
    $core.String? streamarn,
  }) {
    final result = create();
    if (streammodedetails != null) result.streammodedetails = streammodedetails;
    if (warmthroughputmibps != null)
      result.warmthroughputmibps = warmthroughputmibps;
    if (streamid != null) result.streamid = streamid;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateStreamModeInput._();

  factory UpdateStreamModeInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStreamModeInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStreamModeInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<StreamModeDetails>(
        12139665, _omitFieldNames ? '' : 'streammodedetails',
        subBuilder: StreamModeDetails.create)
    ..aI(259219704, _omitFieldNames ? '' : 'warmthroughputmibps')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStreamModeInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStreamModeInput copyWith(
          void Function(UpdateStreamModeInput) updates) =>
      super.copyWith((message) => updates(message as UpdateStreamModeInput))
          as UpdateStreamModeInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStreamModeInput create() => UpdateStreamModeInput._();
  @$core.override
  UpdateStreamModeInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStreamModeInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStreamModeInput>(create);
  static UpdateStreamModeInput? _defaultInstance;

  @$pb.TagNumber(12139665)
  StreamModeDetails get streammodedetails => $_getN(0);
  @$pb.TagNumber(12139665)
  set streammodedetails(StreamModeDetails value) => $_setField(12139665, value);
  @$pb.TagNumber(12139665)
  $core.bool hasStreammodedetails() => $_has(0);
  @$pb.TagNumber(12139665)
  void clearStreammodedetails() => $_clearField(12139665);
  @$pb.TagNumber(12139665)
  StreamModeDetails ensureStreammodedetails() => $_ensure(0);

  @$pb.TagNumber(259219704)
  $core.int get warmthroughputmibps => $_getIZ(1);
  @$pb.TagNumber(259219704)
  set warmthroughputmibps($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(259219704)
  $core.bool hasWarmthroughputmibps() => $_has(1);
  @$pb.TagNumber(259219704)
  void clearWarmthroughputmibps() => $_clearField(259219704);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(2);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(2);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class UpdateStreamWarmThroughputInput extends $pb.GeneratedMessage {
  factory UpdateStreamWarmThroughputInput({
    $core.int? warmthroughputmibps,
    $core.String? streamid,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (warmthroughputmibps != null)
      result.warmthroughputmibps = warmthroughputmibps;
    if (streamid != null) result.streamid = streamid;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateStreamWarmThroughputInput._();

  factory UpdateStreamWarmThroughputInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStreamWarmThroughputInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStreamWarmThroughputInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(259219704, _omitFieldNames ? '' : 'warmthroughputmibps')
    ..aOS(415266497, _omitFieldNames ? '' : 'streamid')
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStreamWarmThroughputInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStreamWarmThroughputInput copyWith(
          void Function(UpdateStreamWarmThroughputInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateStreamWarmThroughputInput))
          as UpdateStreamWarmThroughputInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStreamWarmThroughputInput create() =>
      UpdateStreamWarmThroughputInput._();
  @$core.override
  UpdateStreamWarmThroughputInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStreamWarmThroughputInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStreamWarmThroughputInput>(
          create);
  static UpdateStreamWarmThroughputInput? _defaultInstance;

  @$pb.TagNumber(259219704)
  $core.int get warmthroughputmibps => $_getIZ(0);
  @$pb.TagNumber(259219704)
  set warmthroughputmibps($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(259219704)
  $core.bool hasWarmthroughputmibps() => $_has(0);
  @$pb.TagNumber(259219704)
  void clearWarmthroughputmibps() => $_clearField(259219704);

  @$pb.TagNumber(415266497)
  $core.String get streamid => $_getSZ(1);
  @$pb.TagNumber(415266497)
  set streamid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(415266497)
  $core.bool hasStreamid() => $_has(1);
  @$pb.TagNumber(415266497)
  void clearStreamid() => $_clearField(415266497);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(2);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(2);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
}

class UpdateStreamWarmThroughputOutput extends $pb.GeneratedMessage {
  factory UpdateStreamWarmThroughputOutput({
    WarmThroughputObject? warmthroughput,
    $core.String? streamname,
    $core.String? streamarn,
  }) {
    final result = create();
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (streamname != null) result.streamname = streamname;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateStreamWarmThroughputOutput._();

  factory UpdateStreamWarmThroughputOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStreamWarmThroughputOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStreamWarmThroughputOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOM<WarmThroughputObject>(
        290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: WarmThroughputObject.create)
    ..aOS(470703047, _omitFieldNames ? '' : 'streamname')
    ..aOS(508213725, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStreamWarmThroughputOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStreamWarmThroughputOutput copyWith(
          void Function(UpdateStreamWarmThroughputOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateStreamWarmThroughputOutput))
          as UpdateStreamWarmThroughputOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStreamWarmThroughputOutput create() =>
      UpdateStreamWarmThroughputOutput._();
  @$core.override
  UpdateStreamWarmThroughputOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStreamWarmThroughputOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStreamWarmThroughputOutput>(
          create);
  static UpdateStreamWarmThroughputOutput? _defaultInstance;

  @$pb.TagNumber(290598659)
  WarmThroughputObject get warmthroughput => $_getN(0);
  @$pb.TagNumber(290598659)
  set warmthroughput(WarmThroughputObject value) =>
      $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(0);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  WarmThroughputObject ensureWarmthroughput() => $_ensure(0);

  @$pb.TagNumber(470703047)
  $core.String get streamname => $_getSZ(1);
  @$pb.TagNumber(470703047)
  set streamname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(470703047)
  $core.bool hasStreamname() => $_has(1);
  @$pb.TagNumber(470703047)
  void clearStreamname() => $_clearField(470703047);

  @$pb.TagNumber(508213725)
  $core.String get streamarn => $_getSZ(2);
  @$pb.TagNumber(508213725)
  set streamarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(508213725)
  $core.bool hasStreamarn() => $_has(2);
  @$pb.TagNumber(508213725)
  void clearStreamarn() => $_clearField(508213725);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WarmThroughputObject extends $pb.GeneratedMessage {
  factory WarmThroughputObject({
    $core.int? currentmibps,
    $core.int? targetmibps,
  }) {
    final result = create();
    if (currentmibps != null) result.currentmibps = currentmibps;
    if (targetmibps != null) result.targetmibps = targetmibps;
    return result;
  }

  WarmThroughputObject._();

  factory WarmThroughputObject.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WarmThroughputObject.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WarmThroughputObject',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kinesis'),
      createEmptyInstance: create)
    ..aI(153746976, _omitFieldNames ? '' : 'currentmibps')
    ..aI(535445626, _omitFieldNames ? '' : 'targetmibps')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WarmThroughputObject clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WarmThroughputObject copyWith(void Function(WarmThroughputObject) updates) =>
      super.copyWith((message) => updates(message as WarmThroughputObject))
          as WarmThroughputObject;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WarmThroughputObject create() => WarmThroughputObject._();
  @$core.override
  WarmThroughputObject createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WarmThroughputObject getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WarmThroughputObject>(create);
  static WarmThroughputObject? _defaultInstance;

  @$pb.TagNumber(153746976)
  $core.int get currentmibps => $_getIZ(0);
  @$pb.TagNumber(153746976)
  set currentmibps($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(153746976)
  $core.bool hasCurrentmibps() => $_has(0);
  @$pb.TagNumber(153746976)
  void clearCurrentmibps() => $_clearField(153746976);

  @$pb.TagNumber(535445626)
  $core.int get targetmibps => $_getIZ(1);
  @$pb.TagNumber(535445626)
  set targetmibps($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(535445626)
  $core.bool hasTargetmibps() => $_has(1);
  @$pb.TagNumber(535445626)
  void clearTargetmibps() => $_clearField(535445626);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
