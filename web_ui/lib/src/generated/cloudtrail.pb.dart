// This is a generated file - do not edit.
//
// Generated from cloudtrail.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'cloudtrail.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'cloudtrail.pbenum.dart';

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class AccountHasOngoingImportException extends $pb.GeneratedMessage {
  factory AccountHasOngoingImportException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AccountHasOngoingImportException._();

  factory AccountHasOngoingImportException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccountHasOngoingImportException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountHasOngoingImportException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountHasOngoingImportException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountHasOngoingImportException copyWith(
          void Function(AccountHasOngoingImportException) updates) =>
      super.copyWith(
              (message) => updates(message as AccountHasOngoingImportException))
          as AccountHasOngoingImportException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountHasOngoingImportException create() =>
      AccountHasOngoingImportException._();
  @$core.override
  AccountHasOngoingImportException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccountHasOngoingImportException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountHasOngoingImportException>(
          create);
  static AccountHasOngoingImportException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class AccountNotFoundException extends $pb.GeneratedMessage {
  factory AccountNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AccountNotFoundException._();

  factory AccountNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccountNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountNotFoundException copyWith(
          void Function(AccountNotFoundException) updates) =>
      super.copyWith((message) => updates(message as AccountNotFoundException))
          as AccountNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountNotFoundException create() => AccountNotFoundException._();
  @$core.override
  AccountNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccountNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountNotFoundException>(create);
  static AccountNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class AccountNotRegisteredException extends $pb.GeneratedMessage {
  factory AccountNotRegisteredException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AccountNotRegisteredException._();

  factory AccountNotRegisteredException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccountNotRegisteredException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountNotRegisteredException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountNotRegisteredException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountNotRegisteredException copyWith(
          void Function(AccountNotRegisteredException) updates) =>
      super.copyWith(
              (message) => updates(message as AccountNotRegisteredException))
          as AccountNotRegisteredException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountNotRegisteredException create() =>
      AccountNotRegisteredException._();
  @$core.override
  AccountNotRegisteredException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccountNotRegisteredException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountNotRegisteredException>(create);
  static AccountNotRegisteredException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class AccountRegisteredException extends $pb.GeneratedMessage {
  factory AccountRegisteredException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AccountRegisteredException._();

  factory AccountRegisteredException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccountRegisteredException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountRegisteredException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountRegisteredException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountRegisteredException copyWith(
          void Function(AccountRegisteredException) updates) =>
      super.copyWith(
              (message) => updates(message as AccountRegisteredException))
          as AccountRegisteredException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountRegisteredException create() => AccountRegisteredException._();
  @$core.override
  AccountRegisteredException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccountRegisteredException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountRegisteredException>(create);
  static AccountRegisteredException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class AddTagsRequest extends $pb.GeneratedMessage {
  factory AddTagsRequest({
    $core.Iterable<Tag>? tagslist,
    $core.String? resourceid,
  }) {
    final result = create();
    if (tagslist != null) result.tagslist.addAll(tagslist);
    if (resourceid != null) result.resourceid = resourceid;
    return result;
  }

  AddTagsRequest._();

  factory AddTagsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddTagsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddTagsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..aOS(526146833, _omitFieldNames ? '' : 'resourceid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsRequest copyWith(void Function(AddTagsRequest) updates) =>
      super.copyWith((message) => updates(message as AddTagsRequest))
          as AddTagsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddTagsRequest create() => AddTagsRequest._();
  @$core.override
  AddTagsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddTagsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddTagsRequest>(create);
  static AddTagsRequest? _defaultInstance;

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(0);

  @$pb.TagNumber(526146833)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(526146833)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(526146833)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(526146833)
  void clearResourceid() => $_clearField(526146833);
}

class AddTagsResponse extends $pb.GeneratedMessage {
  factory AddTagsResponse() => create();

  AddTagsResponse._();

  factory AddTagsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddTagsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddTagsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsResponse copyWith(void Function(AddTagsResponse) updates) =>
      super.copyWith((message) => updates(message as AddTagsResponse))
          as AddTagsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddTagsResponse create() => AddTagsResponse._();
  @$core.override
  AddTagsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddTagsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddTagsResponse>(create);
  static AddTagsResponse? _defaultInstance;
}

class AdvancedEventSelector extends $pb.GeneratedMessage {
  factory AdvancedEventSelector({
    $core.String? name,
    $core.Iterable<AdvancedFieldSelector>? fieldselectors,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (fieldselectors != null) result.fieldselectors.addAll(fieldselectors);
    return result;
  }

  AdvancedEventSelector._();

  factory AdvancedEventSelector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AdvancedEventSelector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AdvancedEventSelector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<AdvancedFieldSelector>(
        445033784, _omitFieldNames ? '' : 'fieldselectors',
        subBuilder: AdvancedFieldSelector.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AdvancedEventSelector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AdvancedEventSelector copyWith(
          void Function(AdvancedEventSelector) updates) =>
      super.copyWith((message) => updates(message as AdvancedEventSelector))
          as AdvancedEventSelector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AdvancedEventSelector create() => AdvancedEventSelector._();
  @$core.override
  AdvancedEventSelector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AdvancedEventSelector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AdvancedEventSelector>(create);
  static AdvancedEventSelector? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(445033784)
  $pb.PbList<AdvancedFieldSelector> get fieldselectors => $_getList(1);
}

class AdvancedFieldSelector extends $pb.GeneratedMessage {
  factory AdvancedFieldSelector({
    $core.Iterable<$core.String>? notendswith,
    $core.String? field_263732488,
    $core.Iterable<$core.String>? notstartswith,
    $core.Iterable<$core.String>? notequals,
    $core.Iterable<$core.String>? startswith,
    $core.Iterable<$core.String>? endswith,
    $core.Iterable<$core.String>? equals,
  }) {
    final result = create();
    if (notendswith != null) result.notendswith.addAll(notendswith);
    if (field_263732488 != null) result.field_263732488 = field_263732488;
    if (notstartswith != null) result.notstartswith.addAll(notstartswith);
    if (notequals != null) result.notequals.addAll(notequals);
    if (startswith != null) result.startswith.addAll(startswith);
    if (endswith != null) result.endswith.addAll(endswith);
    if (equals != null) result.equals.addAll(equals);
    return result;
  }

  AdvancedFieldSelector._();

  factory AdvancedFieldSelector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AdvancedFieldSelector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AdvancedFieldSelector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPS(252128485, _omitFieldNames ? '' : 'notendswith')
    ..aOS(263732488, _omitFieldNames ? '' : 'field')
    ..pPS(305277948, _omitFieldNames ? '' : 'notstartswith')
    ..pPS(424280088, _omitFieldNames ? '' : 'notequals')
    ..pPS(468546557, _omitFieldNames ? '' : 'startswith')
    ..pPS(505385124, _omitFieldNames ? '' : 'endswith')
    ..pPS(513367477, _omitFieldNames ? '' : 'equals')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AdvancedFieldSelector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AdvancedFieldSelector copyWith(
          void Function(AdvancedFieldSelector) updates) =>
      super.copyWith((message) => updates(message as AdvancedFieldSelector))
          as AdvancedFieldSelector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AdvancedFieldSelector create() => AdvancedFieldSelector._();
  @$core.override
  AdvancedFieldSelector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AdvancedFieldSelector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AdvancedFieldSelector>(create);
  static AdvancedFieldSelector? _defaultInstance;

  @$pb.TagNumber(252128485)
  $pb.PbList<$core.String> get notendswith => $_getList(0);

  @$pb.TagNumber(263732488)
  $core.String get field_263732488 => $_getSZ(1);
  @$pb.TagNumber(263732488)
  set field_263732488($core.String value) => $_setString(1, value);
  @$pb.TagNumber(263732488)
  $core.bool hasField_263732488() => $_has(1);
  @$pb.TagNumber(263732488)
  void clearField_263732488() => $_clearField(263732488);

  @$pb.TagNumber(305277948)
  $pb.PbList<$core.String> get notstartswith => $_getList(2);

  @$pb.TagNumber(424280088)
  $pb.PbList<$core.String> get notequals => $_getList(3);

  @$pb.TagNumber(468546557)
  $pb.PbList<$core.String> get startswith => $_getList(4);

  @$pb.TagNumber(505385124)
  $pb.PbList<$core.String> get endswith => $_getList(5);

  @$pb.TagNumber(513367477)
  $pb.PbList<$core.String> get equals => $_getList(6);
}

class AggregationConfiguration extends $pb.GeneratedMessage {
  factory AggregationConfiguration({
    $core.Iterable<Template>? templates,
    EventCategoryAggregation? eventcategory,
  }) {
    final result = create();
    if (templates != null) result.templates.addAll(templates);
    if (eventcategory != null) result.eventcategory = eventcategory;
    return result;
  }

  AggregationConfiguration._();

  factory AggregationConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AggregationConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AggregationConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pc<Template>(
        63634313, _omitFieldNames ? '' : 'templates', $pb.PbFieldType.KE,
        valueOf: Template.valueOf,
        enumValues: Template.values,
        defaultEnumValue: Template.TEMPLATE_API_ACTIVITY)
    ..aE<EventCategoryAggregation>(
        164668724, _omitFieldNames ? '' : 'eventcategory',
        enumValues: EventCategoryAggregation.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AggregationConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AggregationConfiguration copyWith(
          void Function(AggregationConfiguration) updates) =>
      super.copyWith((message) => updates(message as AggregationConfiguration))
          as AggregationConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AggregationConfiguration create() => AggregationConfiguration._();
  @$core.override
  AggregationConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AggregationConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AggregationConfiguration>(create);
  static AggregationConfiguration? _defaultInstance;

  @$pb.TagNumber(63634313)
  $pb.PbList<Template> get templates => $_getList(0);

  @$pb.TagNumber(164668724)
  EventCategoryAggregation get eventcategory => $_getN(1);
  @$pb.TagNumber(164668724)
  set eventcategory(EventCategoryAggregation value) =>
      $_setField(164668724, value);
  @$pb.TagNumber(164668724)
  $core.bool hasEventcategory() => $_has(1);
  @$pb.TagNumber(164668724)
  void clearEventcategory() => $_clearField(164668724);
}

class CancelQueryRequest extends $pb.GeneratedMessage {
  factory CancelQueryRequest({
    $core.String? queryid,
    $core.String? eventdatastore,
    $core.String? eventdatastoreowneraccountid,
  }) {
    final result = create();
    if (queryid != null) result.queryid = queryid;
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    return result;
  }

  CancelQueryRequest._();

  factory CancelQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelQueryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryRequest copyWith(void Function(CancelQueryRequest) updates) =>
      super.copyWith((message) => updates(message as CancelQueryRequest))
          as CancelQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelQueryRequest create() => CancelQueryRequest._();
  @$core.override
  CancelQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelQueryRequest>(create);
  static CancelQueryRequest? _defaultInstance;

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(0);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(0);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(1);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(1, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(1);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(2);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(2);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);
}

class CancelQueryResponse extends $pb.GeneratedMessage {
  factory CancelQueryResponse({
    $core.String? queryid,
    QueryStatus? querystatus,
    $core.String? eventdatastoreowneraccountid,
  }) {
    final result = create();
    if (queryid != null) result.queryid = queryid;
    if (querystatus != null) result.querystatus = querystatus;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    return result;
  }

  CancelQueryResponse._();

  factory CancelQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelQueryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..aE<QueryStatus>(367016406, _omitFieldNames ? '' : 'querystatus',
        enumValues: QueryStatus.values)
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryResponse copyWith(void Function(CancelQueryResponse) updates) =>
      super.copyWith((message) => updates(message as CancelQueryResponse))
          as CancelQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelQueryResponse create() => CancelQueryResponse._();
  @$core.override
  CancelQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelQueryResponse>(create);
  static CancelQueryResponse? _defaultInstance;

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(0);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(0);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(367016406)
  QueryStatus get querystatus => $_getN(1);
  @$pb.TagNumber(367016406)
  set querystatus(QueryStatus value) => $_setField(367016406, value);
  @$pb.TagNumber(367016406)
  $core.bool hasQuerystatus() => $_has(1);
  @$pb.TagNumber(367016406)
  void clearQuerystatus() => $_clearField(367016406);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(2);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(2);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);
}

class CannotDelegateManagementAccountException extends $pb.GeneratedMessage {
  factory CannotDelegateManagementAccountException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CannotDelegateManagementAccountException._();

  factory CannotDelegateManagementAccountException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CannotDelegateManagementAccountException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CannotDelegateManagementAccountException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CannotDelegateManagementAccountException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CannotDelegateManagementAccountException copyWith(
          void Function(CannotDelegateManagementAccountException) updates) =>
      super.copyWith((message) =>
              updates(message as CannotDelegateManagementAccountException))
          as CannotDelegateManagementAccountException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CannotDelegateManagementAccountException create() =>
      CannotDelegateManagementAccountException._();
  @$core.override
  CannotDelegateManagementAccountException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CannotDelegateManagementAccountException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CannotDelegateManagementAccountException>(create);
  static CannotDelegateManagementAccountException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class Channel extends $pb.GeneratedMessage {
  factory Channel({
    $core.String? channelarn,
    $core.String? name,
  }) {
    final result = create();
    if (channelarn != null) result.channelarn = channelarn;
    if (name != null) result.name = name;
    return result;
  }

  Channel._();

  factory Channel.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Channel.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Channel',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(99276476, _omitFieldNames ? '' : 'channelarn')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Channel clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Channel copyWith(void Function(Channel) updates) =>
      super.copyWith((message) => updates(message as Channel)) as Channel;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Channel create() => Channel._();
  @$core.override
  Channel createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Channel getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Channel>(create);
  static Channel? _defaultInstance;

  @$pb.TagNumber(99276476)
  $core.String get channelarn => $_getSZ(0);
  @$pb.TagNumber(99276476)
  set channelarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(99276476)
  $core.bool hasChannelarn() => $_has(0);
  @$pb.TagNumber(99276476)
  void clearChannelarn() => $_clearField(99276476);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ChannelARNInvalidException extends $pb.GeneratedMessage {
  factory ChannelARNInvalidException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ChannelARNInvalidException._();

  factory ChannelARNInvalidException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChannelARNInvalidException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChannelARNInvalidException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelARNInvalidException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelARNInvalidException copyWith(
          void Function(ChannelARNInvalidException) updates) =>
      super.copyWith(
              (message) => updates(message as ChannelARNInvalidException))
          as ChannelARNInvalidException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChannelARNInvalidException create() => ChannelARNInvalidException._();
  @$core.override
  ChannelARNInvalidException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChannelARNInvalidException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChannelARNInvalidException>(create);
  static ChannelARNInvalidException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ChannelAlreadyExistsException extends $pb.GeneratedMessage {
  factory ChannelAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ChannelAlreadyExistsException._();

  factory ChannelAlreadyExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChannelAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChannelAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelAlreadyExistsException copyWith(
          void Function(ChannelAlreadyExistsException) updates) =>
      super.copyWith(
              (message) => updates(message as ChannelAlreadyExistsException))
          as ChannelAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChannelAlreadyExistsException create() =>
      ChannelAlreadyExistsException._();
  @$core.override
  ChannelAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChannelAlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChannelAlreadyExistsException>(create);
  static ChannelAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ChannelExistsForEDSException extends $pb.GeneratedMessage {
  factory ChannelExistsForEDSException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ChannelExistsForEDSException._();

  factory ChannelExistsForEDSException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChannelExistsForEDSException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChannelExistsForEDSException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelExistsForEDSException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelExistsForEDSException copyWith(
          void Function(ChannelExistsForEDSException) updates) =>
      super.copyWith(
              (message) => updates(message as ChannelExistsForEDSException))
          as ChannelExistsForEDSException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChannelExistsForEDSException create() =>
      ChannelExistsForEDSException._();
  @$core.override
  ChannelExistsForEDSException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChannelExistsForEDSException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChannelExistsForEDSException>(create);
  static ChannelExistsForEDSException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ChannelMaxLimitExceededException extends $pb.GeneratedMessage {
  factory ChannelMaxLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ChannelMaxLimitExceededException._();

  factory ChannelMaxLimitExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChannelMaxLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChannelMaxLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelMaxLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelMaxLimitExceededException copyWith(
          void Function(ChannelMaxLimitExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as ChannelMaxLimitExceededException))
          as ChannelMaxLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChannelMaxLimitExceededException create() =>
      ChannelMaxLimitExceededException._();
  @$core.override
  ChannelMaxLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChannelMaxLimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChannelMaxLimitExceededException>(
          create);
  static ChannelMaxLimitExceededException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ChannelNotFoundException extends $pb.GeneratedMessage {
  factory ChannelNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ChannelNotFoundException._();

  factory ChannelNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChannelNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChannelNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChannelNotFoundException copyWith(
          void Function(ChannelNotFoundException) updates) =>
      super.copyWith((message) => updates(message as ChannelNotFoundException))
          as ChannelNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChannelNotFoundException create() => ChannelNotFoundException._();
  @$core.override
  ChannelNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChannelNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChannelNotFoundException>(create);
  static ChannelNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CloudTrailARNInvalidException extends $pb.GeneratedMessage {
  factory CloudTrailARNInvalidException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudTrailARNInvalidException._();

  factory CloudTrailARNInvalidException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudTrailARNInvalidException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudTrailARNInvalidException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudTrailARNInvalidException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudTrailARNInvalidException copyWith(
          void Function(CloudTrailARNInvalidException) updates) =>
      super.copyWith(
              (message) => updates(message as CloudTrailARNInvalidException))
          as CloudTrailARNInvalidException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudTrailARNInvalidException create() =>
      CloudTrailARNInvalidException._();
  @$core.override
  CloudTrailARNInvalidException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudTrailARNInvalidException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudTrailARNInvalidException>(create);
  static CloudTrailARNInvalidException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CloudTrailAccessNotEnabledException extends $pb.GeneratedMessage {
  factory CloudTrailAccessNotEnabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudTrailAccessNotEnabledException._();

  factory CloudTrailAccessNotEnabledException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudTrailAccessNotEnabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudTrailAccessNotEnabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudTrailAccessNotEnabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudTrailAccessNotEnabledException copyWith(
          void Function(CloudTrailAccessNotEnabledException) updates) =>
      super.copyWith((message) =>
              updates(message as CloudTrailAccessNotEnabledException))
          as CloudTrailAccessNotEnabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudTrailAccessNotEnabledException create() =>
      CloudTrailAccessNotEnabledException._();
  @$core.override
  CloudTrailAccessNotEnabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudTrailAccessNotEnabledException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CloudTrailAccessNotEnabledException>(create);
  static CloudTrailAccessNotEnabledException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CloudTrailInvalidClientTokenIdException extends $pb.GeneratedMessage {
  factory CloudTrailInvalidClientTokenIdException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudTrailInvalidClientTokenIdException._();

  factory CloudTrailInvalidClientTokenIdException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudTrailInvalidClientTokenIdException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudTrailInvalidClientTokenIdException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudTrailInvalidClientTokenIdException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudTrailInvalidClientTokenIdException copyWith(
          void Function(CloudTrailInvalidClientTokenIdException) updates) =>
      super.copyWith((message) =>
              updates(message as CloudTrailInvalidClientTokenIdException))
          as CloudTrailInvalidClientTokenIdException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudTrailInvalidClientTokenIdException create() =>
      CloudTrailInvalidClientTokenIdException._();
  @$core.override
  CloudTrailInvalidClientTokenIdException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudTrailInvalidClientTokenIdException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CloudTrailInvalidClientTokenIdException>(create);
  static CloudTrailInvalidClientTokenIdException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CloudWatchLogsDeliveryUnavailableException extends $pb.GeneratedMessage {
  factory CloudWatchLogsDeliveryUnavailableException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudWatchLogsDeliveryUnavailableException._();

  factory CloudWatchLogsDeliveryUnavailableException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudWatchLogsDeliveryUnavailableException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudWatchLogsDeliveryUnavailableException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchLogsDeliveryUnavailableException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchLogsDeliveryUnavailableException copyWith(
          void Function(CloudWatchLogsDeliveryUnavailableException) updates) =>
      super.copyWith((message) =>
              updates(message as CloudWatchLogsDeliveryUnavailableException))
          as CloudWatchLogsDeliveryUnavailableException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudWatchLogsDeliveryUnavailableException create() =>
      CloudWatchLogsDeliveryUnavailableException._();
  @$core.override
  CloudWatchLogsDeliveryUnavailableException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudWatchLogsDeliveryUnavailableException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CloudWatchLogsDeliveryUnavailableException>(create);
  static CloudWatchLogsDeliveryUnavailableException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
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

class ContextKeySelector extends $pb.GeneratedMessage {
  factory ContextKeySelector({
    Type? type,
    $core.Iterable<$core.String>? equals,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (equals != null) result.equals.addAll(equals);
    return result;
  }

  ContextKeySelector._();

  factory ContextKeySelector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ContextKeySelector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContextKeySelector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<Type>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: Type.values)
    ..pPS(513367477, _omitFieldNames ? '' : 'equals')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContextKeySelector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContextKeySelector copyWith(void Function(ContextKeySelector) updates) =>
      super.copyWith((message) => updates(message as ContextKeySelector))
          as ContextKeySelector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ContextKeySelector create() => ContextKeySelector._();
  @$core.override
  ContextKeySelector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ContextKeySelector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ContextKeySelector>(create);
  static ContextKeySelector? _defaultInstance;

  @$pb.TagNumber(290836590)
  Type get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(Type value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(513367477)
  $pb.PbList<$core.String> get equals => $_getList(1);
}

class CreateChannelRequest extends $pb.GeneratedMessage {
  factory CreateChannelRequest({
    $core.String? source,
    $core.String? name,
    $core.Iterable<Tag>? tags,
    $core.Iterable<Destination>? destinations,
  }) {
    final result = create();
    if (source != null) result.source = source;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    if (destinations != null) result.destinations.addAll(destinations);
    return result;
  }

  CreateChannelRequest._();

  factory CreateChannelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateChannelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateChannelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(31630329, _omitFieldNames ? '' : 'source')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..pPM<Destination>(404385861, _omitFieldNames ? '' : 'destinations',
        subBuilder: Destination.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateChannelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateChannelRequest copyWith(void Function(CreateChannelRequest) updates) =>
      super.copyWith((message) => updates(message as CreateChannelRequest))
          as CreateChannelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateChannelRequest create() => CreateChannelRequest._();
  @$core.override
  CreateChannelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateChannelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateChannelRequest>(create);
  static CreateChannelRequest? _defaultInstance;

  @$pb.TagNumber(31630329)
  $core.String get source => $_getSZ(0);
  @$pb.TagNumber(31630329)
  set source($core.String value) => $_setString(0, value);
  @$pb.TagNumber(31630329)
  $core.bool hasSource() => $_has(0);
  @$pb.TagNumber(31630329)
  void clearSource() => $_clearField(31630329);

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

  @$pb.TagNumber(404385861)
  $pb.PbList<Destination> get destinations => $_getList(3);
}

class CreateChannelResponse extends $pb.GeneratedMessage {
  factory CreateChannelResponse({
    $core.String? source,
    $core.String? channelarn,
    $core.String? name,
    $core.Iterable<Tag>? tags,
    $core.Iterable<Destination>? destinations,
  }) {
    final result = create();
    if (source != null) result.source = source;
    if (channelarn != null) result.channelarn = channelarn;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    if (destinations != null) result.destinations.addAll(destinations);
    return result;
  }

  CreateChannelResponse._();

  factory CreateChannelResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateChannelResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateChannelResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(31630329, _omitFieldNames ? '' : 'source')
    ..aOS(99276476, _omitFieldNames ? '' : 'channelarn')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..pPM<Destination>(404385861, _omitFieldNames ? '' : 'destinations',
        subBuilder: Destination.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateChannelResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateChannelResponse copyWith(
          void Function(CreateChannelResponse) updates) =>
      super.copyWith((message) => updates(message as CreateChannelResponse))
          as CreateChannelResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateChannelResponse create() => CreateChannelResponse._();
  @$core.override
  CreateChannelResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateChannelResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateChannelResponse>(create);
  static CreateChannelResponse? _defaultInstance;

  @$pb.TagNumber(31630329)
  $core.String get source => $_getSZ(0);
  @$pb.TagNumber(31630329)
  set source($core.String value) => $_setString(0, value);
  @$pb.TagNumber(31630329)
  $core.bool hasSource() => $_has(0);
  @$pb.TagNumber(31630329)
  void clearSource() => $_clearField(31630329);

  @$pb.TagNumber(99276476)
  $core.String get channelarn => $_getSZ(1);
  @$pb.TagNumber(99276476)
  set channelarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(99276476)
  $core.bool hasChannelarn() => $_has(1);
  @$pb.TagNumber(99276476)
  void clearChannelarn() => $_clearField(99276476);

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

  @$pb.TagNumber(404385861)
  $pb.PbList<Destination> get destinations => $_getList(4);
}

class CreateDashboardRequest extends $pb.GeneratedMessage {
  factory CreateDashboardRequest({
    RefreshSchedule? refreshschedule,
    $core.String? name,
    $core.bool? terminationprotectionenabled,
    $core.Iterable<Tag>? tagslist,
    $core.Iterable<RequestWidget>? widgets,
  }) {
    final result = create();
    if (refreshschedule != null) result.refreshschedule = refreshschedule;
    if (name != null) result.name = name;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (tagslist != null) result.tagslist.addAll(tagslist);
    if (widgets != null) result.widgets.addAll(widgets);
    return result;
  }

  CreateDashboardRequest._();

  factory CreateDashboardRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDashboardRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDashboardRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<RefreshSchedule>(261773338, _omitFieldNames ? '' : 'refreshschedule',
        subBuilder: RefreshSchedule.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..pPM<RequestWidget>(501826147, _omitFieldNames ? '' : 'widgets',
        subBuilder: RequestWidget.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDashboardRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDashboardRequest copyWith(
          void Function(CreateDashboardRequest) updates) =>
      super.copyWith((message) => updates(message as CreateDashboardRequest))
          as CreateDashboardRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDashboardRequest create() => CreateDashboardRequest._();
  @$core.override
  CreateDashboardRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDashboardRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDashboardRequest>(create);
  static CreateDashboardRequest? _defaultInstance;

  @$pb.TagNumber(261773338)
  RefreshSchedule get refreshschedule => $_getN(0);
  @$pb.TagNumber(261773338)
  set refreshschedule(RefreshSchedule value) => $_setField(261773338, value);
  @$pb.TagNumber(261773338)
  $core.bool hasRefreshschedule() => $_has(0);
  @$pb.TagNumber(261773338)
  void clearRefreshschedule() => $_clearField(261773338);
  @$pb.TagNumber(261773338)
  RefreshSchedule ensureRefreshschedule() => $_ensure(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(2);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(2);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(3);

  @$pb.TagNumber(501826147)
  $pb.PbList<RequestWidget> get widgets => $_getList(4);
}

class CreateDashboardResponse extends $pb.GeneratedMessage {
  factory CreateDashboardResponse({
    $core.String? dashboardarn,
    RefreshSchedule? refreshschedule,
    $core.String? name,
    DashboardType? type,
    $core.bool? terminationprotectionenabled,
    $core.Iterable<Tag>? tagslist,
    $core.Iterable<Widget>? widgets,
  }) {
    final result = create();
    if (dashboardarn != null) result.dashboardarn = dashboardarn;
    if (refreshschedule != null) result.refreshschedule = refreshschedule;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (tagslist != null) result.tagslist.addAll(tagslist);
    if (widgets != null) result.widgets.addAll(widgets);
    return result;
  }

  CreateDashboardResponse._();

  factory CreateDashboardResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDashboardResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDashboardResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(108051951, _omitFieldNames ? '' : 'dashboardarn')
    ..aOM<RefreshSchedule>(261773338, _omitFieldNames ? '' : 'refreshschedule',
        subBuilder: RefreshSchedule.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<DashboardType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DashboardType.values)
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..pPM<Widget>(501826147, _omitFieldNames ? '' : 'widgets',
        subBuilder: Widget.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDashboardResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDashboardResponse copyWith(
          void Function(CreateDashboardResponse) updates) =>
      super.copyWith((message) => updates(message as CreateDashboardResponse))
          as CreateDashboardResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDashboardResponse create() => CreateDashboardResponse._();
  @$core.override
  CreateDashboardResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDashboardResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDashboardResponse>(create);
  static CreateDashboardResponse? _defaultInstance;

  @$pb.TagNumber(108051951)
  $core.String get dashboardarn => $_getSZ(0);
  @$pb.TagNumber(108051951)
  set dashboardarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(108051951)
  $core.bool hasDashboardarn() => $_has(0);
  @$pb.TagNumber(108051951)
  void clearDashboardarn() => $_clearField(108051951);

  @$pb.TagNumber(261773338)
  RefreshSchedule get refreshschedule => $_getN(1);
  @$pb.TagNumber(261773338)
  set refreshschedule(RefreshSchedule value) => $_setField(261773338, value);
  @$pb.TagNumber(261773338)
  $core.bool hasRefreshschedule() => $_has(1);
  @$pb.TagNumber(261773338)
  void clearRefreshschedule() => $_clearField(261773338);
  @$pb.TagNumber(261773338)
  RefreshSchedule ensureRefreshschedule() => $_ensure(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  DashboardType get type => $_getN(3);
  @$pb.TagNumber(290836590)
  set type(DashboardType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(3);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(4);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(4);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(5);

  @$pb.TagNumber(501826147)
  $pb.PbList<Widget> get widgets => $_getList(6);
}

class CreateEventDataStoreRequest extends $pb.GeneratedMessage {
  factory CreateEventDataStoreRequest({
    $core.bool? multiregionenabled,
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? kmskeyid,
    $core.bool? startingestion,
    BillingMode? billingmode,
    $core.int? retentionperiod,
    $core.String? name,
    $core.bool? terminationprotectionenabled,
    $core.bool? organizationenabled,
    $core.Iterable<Tag>? tagslist,
  }) {
    final result = create();
    if (multiregionenabled != null)
      result.multiregionenabled = multiregionenabled;
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (startingestion != null) result.startingestion = startingestion;
    if (billingmode != null) result.billingmode = billingmode;
    if (retentionperiod != null) result.retentionperiod = retentionperiod;
    if (name != null) result.name = name;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (organizationenabled != null)
      result.organizationenabled = organizationenabled;
    if (tagslist != null) result.tagslist.addAll(tagslist);
    return result;
  }

  CreateEventDataStoreRequest._();

  factory CreateEventDataStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateEventDataStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateEventDataStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOB(20620620, _omitFieldNames ? '' : 'multiregionenabled')
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOB(168263832, _omitFieldNames ? '' : 'startingestion')
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aI(196383721, _omitFieldNames ? '' : 'retentionperiod')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOB(480171176, _omitFieldNames ? '' : 'organizationenabled')
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventDataStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventDataStoreRequest copyWith(
          void Function(CreateEventDataStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateEventDataStoreRequest))
          as CreateEventDataStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateEventDataStoreRequest create() =>
      CreateEventDataStoreRequest._();
  @$core.override
  CreateEventDataStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateEventDataStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateEventDataStoreRequest>(create);
  static CreateEventDataStoreRequest? _defaultInstance;

  @$pb.TagNumber(20620620)
  $core.bool get multiregionenabled => $_getBF(0);
  @$pb.TagNumber(20620620)
  set multiregionenabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(20620620)
  $core.bool hasMultiregionenabled() => $_has(0);
  @$pb.TagNumber(20620620)
  void clearMultiregionenabled() => $_clearField(20620620);

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(1);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(2);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(2);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(168263832)
  $core.bool get startingestion => $_getBF(3);
  @$pb.TagNumber(168263832)
  set startingestion($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(168263832)
  $core.bool hasStartingestion() => $_has(3);
  @$pb.TagNumber(168263832)
  void clearStartingestion() => $_clearField(168263832);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(4);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(4);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(196383721)
  $core.int get retentionperiod => $_getIZ(5);
  @$pb.TagNumber(196383721)
  set retentionperiod($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(196383721)
  $core.bool hasRetentionperiod() => $_has(5);
  @$pb.TagNumber(196383721)
  void clearRetentionperiod() => $_clearField(196383721);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(7);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(7);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(480171176)
  $core.bool get organizationenabled => $_getBF(8);
  @$pb.TagNumber(480171176)
  set organizationenabled($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(480171176)
  $core.bool hasOrganizationenabled() => $_has(8);
  @$pb.TagNumber(480171176)
  void clearOrganizationenabled() => $_clearField(480171176);

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(9);
}

class CreateEventDataStoreResponse extends $pb.GeneratedMessage {
  factory CreateEventDataStoreResponse({
    EventDataStoreStatus? status,
    $core.bool? multiregionenabled,
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? updatedtimestamp,
    $core.String? kmskeyid,
    BillingMode? billingmode,
    $core.int? retentionperiod,
    $core.String? name,
    $core.String? eventdatastorearn,
    $core.String? createdtimestamp,
    $core.bool? terminationprotectionenabled,
    $core.bool? organizationenabled,
    $core.Iterable<Tag>? tagslist,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (multiregionenabled != null)
      result.multiregionenabled = multiregionenabled;
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (billingmode != null) result.billingmode = billingmode;
    if (retentionperiod != null) result.retentionperiod = retentionperiod;
    if (name != null) result.name = name;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (organizationenabled != null)
      result.organizationenabled = organizationenabled;
    if (tagslist != null) result.tagslist.addAll(tagslist);
    return result;
  }

  CreateEventDataStoreResponse._();

  factory CreateEventDataStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateEventDataStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateEventDataStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<EventDataStoreStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: EventDataStoreStatus.values)
    ..aOB(20620620, _omitFieldNames ? '' : 'multiregionenabled')
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aI(196383721, _omitFieldNames ? '' : 'retentionperiod')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOB(480171176, _omitFieldNames ? '' : 'organizationenabled')
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventDataStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventDataStoreResponse copyWith(
          void Function(CreateEventDataStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateEventDataStoreResponse))
          as CreateEventDataStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateEventDataStoreResponse create() =>
      CreateEventDataStoreResponse._();
  @$core.override
  CreateEventDataStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateEventDataStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateEventDataStoreResponse>(create);
  static CreateEventDataStoreResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  EventDataStoreStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(EventDataStoreStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(20620620)
  $core.bool get multiregionenabled => $_getBF(1);
  @$pb.TagNumber(20620620)
  set multiregionenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(20620620)
  $core.bool hasMultiregionenabled() => $_has(1);
  @$pb.TagNumber(20620620)
  void clearMultiregionenabled() => $_clearField(20620620);

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(2);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(3);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(3);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(4);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(4);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(5);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(5);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(196383721)
  $core.int get retentionperiod => $_getIZ(6);
  @$pb.TagNumber(196383721)
  set retentionperiod($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(196383721)
  $core.bool hasRetentionperiod() => $_has(6);
  @$pb.TagNumber(196383721)
  void clearRetentionperiod() => $_clearField(196383721);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(8);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(8);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(9);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(9, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(9);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(10);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(10);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(480171176)
  $core.bool get organizationenabled => $_getBF(11);
  @$pb.TagNumber(480171176)
  set organizationenabled($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(480171176)
  $core.bool hasOrganizationenabled() => $_has(11);
  @$pb.TagNumber(480171176)
  void clearOrganizationenabled() => $_clearField(480171176);

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(12);
}

class CreateTrailRequest extends $pb.GeneratedMessage {
  factory CreateTrailRequest({
    $core.String? kmskeyid,
    $core.String? cloudwatchlogsrolearn,
    $core.bool? enablelogfilevalidation,
    $core.bool? isorganizationtrail,
    $core.String? s3keyprefix,
    $core.bool? includeglobalserviceevents,
    $core.String? name,
    $core.String? s3bucketname,
    $core.String? snstopicname,
    $core.bool? ismultiregiontrail,
    $core.Iterable<Tag>? tagslist,
    $core.String? cloudwatchlogsloggrouparn,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (cloudwatchlogsrolearn != null)
      result.cloudwatchlogsrolearn = cloudwatchlogsrolearn;
    if (enablelogfilevalidation != null)
      result.enablelogfilevalidation = enablelogfilevalidation;
    if (isorganizationtrail != null)
      result.isorganizationtrail = isorganizationtrail;
    if (s3keyprefix != null) result.s3keyprefix = s3keyprefix;
    if (includeglobalserviceevents != null)
      result.includeglobalserviceevents = includeglobalserviceevents;
    if (name != null) result.name = name;
    if (s3bucketname != null) result.s3bucketname = s3bucketname;
    if (snstopicname != null) result.snstopicname = snstopicname;
    if (ismultiregiontrail != null)
      result.ismultiregiontrail = ismultiregiontrail;
    if (tagslist != null) result.tagslist.addAll(tagslist);
    if (cloudwatchlogsloggrouparn != null)
      result.cloudwatchlogsloggrouparn = cloudwatchlogsloggrouparn;
    return result;
  }

  CreateTrailRequest._();

  factory CreateTrailRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrailRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrailRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(55454690, _omitFieldNames ? '' : 'cloudwatchlogsrolearn')
    ..aOB(80236930, _omitFieldNames ? '' : 'enablelogfilevalidation')
    ..aOB(145256127, _omitFieldNames ? '' : 'isorganizationtrail')
    ..aOS(206015359, _omitFieldNames ? '' : 's3keyprefix')
    ..aOB(212227451, _omitFieldNames ? '' : 'includeglobalserviceevents')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(320495427, _omitFieldNames ? '' : 's3bucketname')
    ..aOS(415454800, _omitFieldNames ? '' : 'snstopicname')
    ..aOB(468658571, _omitFieldNames ? '' : 'ismultiregiontrail')
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..aOS(519613355, _omitFieldNames ? '' : 'cloudwatchlogsloggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrailRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrailRequest copyWith(void Function(CreateTrailRequest) updates) =>
      super.copyWith((message) => updates(message as CreateTrailRequest))
          as CreateTrailRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrailRequest create() => CreateTrailRequest._();
  @$core.override
  CreateTrailRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrailRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTrailRequest>(create);
  static CreateTrailRequest? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(55454690)
  $core.String get cloudwatchlogsrolearn => $_getSZ(1);
  @$pb.TagNumber(55454690)
  set cloudwatchlogsrolearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(55454690)
  $core.bool hasCloudwatchlogsrolearn() => $_has(1);
  @$pb.TagNumber(55454690)
  void clearCloudwatchlogsrolearn() => $_clearField(55454690);

  @$pb.TagNumber(80236930)
  $core.bool get enablelogfilevalidation => $_getBF(2);
  @$pb.TagNumber(80236930)
  set enablelogfilevalidation($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(80236930)
  $core.bool hasEnablelogfilevalidation() => $_has(2);
  @$pb.TagNumber(80236930)
  void clearEnablelogfilevalidation() => $_clearField(80236930);

  @$pb.TagNumber(145256127)
  $core.bool get isorganizationtrail => $_getBF(3);
  @$pb.TagNumber(145256127)
  set isorganizationtrail($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(145256127)
  $core.bool hasIsorganizationtrail() => $_has(3);
  @$pb.TagNumber(145256127)
  void clearIsorganizationtrail() => $_clearField(145256127);

  @$pb.TagNumber(206015359)
  $core.String get s3keyprefix => $_getSZ(4);
  @$pb.TagNumber(206015359)
  set s3keyprefix($core.String value) => $_setString(4, value);
  @$pb.TagNumber(206015359)
  $core.bool hasS3keyprefix() => $_has(4);
  @$pb.TagNumber(206015359)
  void clearS3keyprefix() => $_clearField(206015359);

  @$pb.TagNumber(212227451)
  $core.bool get includeglobalserviceevents => $_getBF(5);
  @$pb.TagNumber(212227451)
  set includeglobalserviceevents($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(212227451)
  $core.bool hasIncludeglobalserviceevents() => $_has(5);
  @$pb.TagNumber(212227451)
  void clearIncludeglobalserviceevents() => $_clearField(212227451);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(320495427)
  $core.String get s3bucketname => $_getSZ(7);
  @$pb.TagNumber(320495427)
  set s3bucketname($core.String value) => $_setString(7, value);
  @$pb.TagNumber(320495427)
  $core.bool hasS3bucketname() => $_has(7);
  @$pb.TagNumber(320495427)
  void clearS3bucketname() => $_clearField(320495427);

  @$pb.TagNumber(415454800)
  $core.String get snstopicname => $_getSZ(8);
  @$pb.TagNumber(415454800)
  set snstopicname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(415454800)
  $core.bool hasSnstopicname() => $_has(8);
  @$pb.TagNumber(415454800)
  void clearSnstopicname() => $_clearField(415454800);

  @$pb.TagNumber(468658571)
  $core.bool get ismultiregiontrail => $_getBF(9);
  @$pb.TagNumber(468658571)
  set ismultiregiontrail($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(468658571)
  $core.bool hasIsmultiregiontrail() => $_has(9);
  @$pb.TagNumber(468658571)
  void clearIsmultiregiontrail() => $_clearField(468658571);

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(10);

  @$pb.TagNumber(519613355)
  $core.String get cloudwatchlogsloggrouparn => $_getSZ(11);
  @$pb.TagNumber(519613355)
  set cloudwatchlogsloggrouparn($core.String value) => $_setString(11, value);
  @$pb.TagNumber(519613355)
  $core.bool hasCloudwatchlogsloggrouparn() => $_has(11);
  @$pb.TagNumber(519613355)
  void clearCloudwatchlogsloggrouparn() => $_clearField(519613355);
}

class CreateTrailResponse extends $pb.GeneratedMessage {
  factory CreateTrailResponse({
    $core.bool? logfilevalidationenabled,
    $core.String? trailarn,
    $core.String? kmskeyid,
    $core.String? cloudwatchlogsrolearn,
    $core.bool? isorganizationtrail,
    $core.String? s3keyprefix,
    $core.bool? includeglobalserviceevents,
    $core.String? name,
    $core.String? s3bucketname,
    $core.String? snstopicarn,
    $core.String? snstopicname,
    $core.bool? ismultiregiontrail,
    $core.String? cloudwatchlogsloggrouparn,
  }) {
    final result = create();
    if (logfilevalidationenabled != null)
      result.logfilevalidationenabled = logfilevalidationenabled;
    if (trailarn != null) result.trailarn = trailarn;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (cloudwatchlogsrolearn != null)
      result.cloudwatchlogsrolearn = cloudwatchlogsrolearn;
    if (isorganizationtrail != null)
      result.isorganizationtrail = isorganizationtrail;
    if (s3keyprefix != null) result.s3keyprefix = s3keyprefix;
    if (includeglobalserviceevents != null)
      result.includeglobalserviceevents = includeglobalserviceevents;
    if (name != null) result.name = name;
    if (s3bucketname != null) result.s3bucketname = s3bucketname;
    if (snstopicarn != null) result.snstopicarn = snstopicarn;
    if (snstopicname != null) result.snstopicname = snstopicname;
    if (ismultiregiontrail != null)
      result.ismultiregiontrail = ismultiregiontrail;
    if (cloudwatchlogsloggrouparn != null)
      result.cloudwatchlogsloggrouparn = cloudwatchlogsloggrouparn;
    return result;
  }

  CreateTrailResponse._();

  factory CreateTrailResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrailResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrailResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOB(35904346, _omitFieldNames ? '' : 'logfilevalidationenabled')
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(55454690, _omitFieldNames ? '' : 'cloudwatchlogsrolearn')
    ..aOB(145256127, _omitFieldNames ? '' : 'isorganizationtrail')
    ..aOS(206015359, _omitFieldNames ? '' : 's3keyprefix')
    ..aOB(212227451, _omitFieldNames ? '' : 'includeglobalserviceevents')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(320495427, _omitFieldNames ? '' : 's3bucketname')
    ..aOS(380025580, _omitFieldNames ? '' : 'snstopicarn')
    ..aOS(415454800, _omitFieldNames ? '' : 'snstopicname')
    ..aOB(468658571, _omitFieldNames ? '' : 'ismultiregiontrail')
    ..aOS(519613355, _omitFieldNames ? '' : 'cloudwatchlogsloggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrailResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrailResponse copyWith(void Function(CreateTrailResponse) updates) =>
      super.copyWith((message) => updates(message as CreateTrailResponse))
          as CreateTrailResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrailResponse create() => CreateTrailResponse._();
  @$core.override
  CreateTrailResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrailResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTrailResponse>(create);
  static CreateTrailResponse? _defaultInstance;

  @$pb.TagNumber(35904346)
  $core.bool get logfilevalidationenabled => $_getBF(0);
  @$pb.TagNumber(35904346)
  set logfilevalidationenabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(35904346)
  $core.bool hasLogfilevalidationenabled() => $_has(0);
  @$pb.TagNumber(35904346)
  void clearLogfilevalidationenabled() => $_clearField(35904346);

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(1);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(1);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(2);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(2);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(55454690)
  $core.String get cloudwatchlogsrolearn => $_getSZ(3);
  @$pb.TagNumber(55454690)
  set cloudwatchlogsrolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(55454690)
  $core.bool hasCloudwatchlogsrolearn() => $_has(3);
  @$pb.TagNumber(55454690)
  void clearCloudwatchlogsrolearn() => $_clearField(55454690);

  @$pb.TagNumber(145256127)
  $core.bool get isorganizationtrail => $_getBF(4);
  @$pb.TagNumber(145256127)
  set isorganizationtrail($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(145256127)
  $core.bool hasIsorganizationtrail() => $_has(4);
  @$pb.TagNumber(145256127)
  void clearIsorganizationtrail() => $_clearField(145256127);

  @$pb.TagNumber(206015359)
  $core.String get s3keyprefix => $_getSZ(5);
  @$pb.TagNumber(206015359)
  set s3keyprefix($core.String value) => $_setString(5, value);
  @$pb.TagNumber(206015359)
  $core.bool hasS3keyprefix() => $_has(5);
  @$pb.TagNumber(206015359)
  void clearS3keyprefix() => $_clearField(206015359);

  @$pb.TagNumber(212227451)
  $core.bool get includeglobalserviceevents => $_getBF(6);
  @$pb.TagNumber(212227451)
  set includeglobalserviceevents($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(212227451)
  $core.bool hasIncludeglobalserviceevents() => $_has(6);
  @$pb.TagNumber(212227451)
  void clearIncludeglobalserviceevents() => $_clearField(212227451);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(320495427)
  $core.String get s3bucketname => $_getSZ(8);
  @$pb.TagNumber(320495427)
  set s3bucketname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(320495427)
  $core.bool hasS3bucketname() => $_has(8);
  @$pb.TagNumber(320495427)
  void clearS3bucketname() => $_clearField(320495427);

  @$pb.TagNumber(380025580)
  $core.String get snstopicarn => $_getSZ(9);
  @$pb.TagNumber(380025580)
  set snstopicarn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(380025580)
  $core.bool hasSnstopicarn() => $_has(9);
  @$pb.TagNumber(380025580)
  void clearSnstopicarn() => $_clearField(380025580);

  @$pb.TagNumber(415454800)
  $core.String get snstopicname => $_getSZ(10);
  @$pb.TagNumber(415454800)
  set snstopicname($core.String value) => $_setString(10, value);
  @$pb.TagNumber(415454800)
  $core.bool hasSnstopicname() => $_has(10);
  @$pb.TagNumber(415454800)
  void clearSnstopicname() => $_clearField(415454800);

  @$pb.TagNumber(468658571)
  $core.bool get ismultiregiontrail => $_getBF(11);
  @$pb.TagNumber(468658571)
  set ismultiregiontrail($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(468658571)
  $core.bool hasIsmultiregiontrail() => $_has(11);
  @$pb.TagNumber(468658571)
  void clearIsmultiregiontrail() => $_clearField(468658571);

  @$pb.TagNumber(519613355)
  $core.String get cloudwatchlogsloggrouparn => $_getSZ(12);
  @$pb.TagNumber(519613355)
  set cloudwatchlogsloggrouparn($core.String value) => $_setString(12, value);
  @$pb.TagNumber(519613355)
  $core.bool hasCloudwatchlogsloggrouparn() => $_has(12);
  @$pb.TagNumber(519613355)
  void clearCloudwatchlogsloggrouparn() => $_clearField(519613355);
}

class DashboardDetail extends $pb.GeneratedMessage {
  factory DashboardDetail({
    $core.String? dashboardarn,
    DashboardType? type,
  }) {
    final result = create();
    if (dashboardarn != null) result.dashboardarn = dashboardarn;
    if (type != null) result.type = type;
    return result;
  }

  DashboardDetail._();

  factory DashboardDetail.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DashboardDetail.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DashboardDetail',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(108051951, _omitFieldNames ? '' : 'dashboardarn')
    ..aE<DashboardType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DashboardType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardDetail clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardDetail copyWith(void Function(DashboardDetail) updates) =>
      super.copyWith((message) => updates(message as DashboardDetail))
          as DashboardDetail;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DashboardDetail create() => DashboardDetail._();
  @$core.override
  DashboardDetail createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DashboardDetail getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DashboardDetail>(create);
  static DashboardDetail? _defaultInstance;

  @$pb.TagNumber(108051951)
  $core.String get dashboardarn => $_getSZ(0);
  @$pb.TagNumber(108051951)
  set dashboardarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(108051951)
  $core.bool hasDashboardarn() => $_has(0);
  @$pb.TagNumber(108051951)
  void clearDashboardarn() => $_clearField(108051951);

  @$pb.TagNumber(290836590)
  DashboardType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(DashboardType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class DataResource extends $pb.GeneratedMessage {
  factory DataResource({
    $core.Iterable<$core.String>? values,
    $core.String? type,
  }) {
    final result = create();
    if (values != null) result.values.addAll(values);
    if (type != null) result.type = type;
    return result;
  }

  DataResource._();

  factory DataResource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataResource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataResource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPS(223158876, _omitFieldNames ? '' : 'values')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataResource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataResource copyWith(void Function(DataResource) updates) =>
      super.copyWith((message) => updates(message as DataResource))
          as DataResource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataResource create() => DataResource._();
  @$core.override
  DataResource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataResource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataResource>(create);
  static DataResource? _defaultInstance;

  @$pb.TagNumber(223158876)
  $pb.PbList<$core.String> get values => $_getList(0);

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(1);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(1, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class DelegatedAdminAccountLimitExceededException extends $pb.GeneratedMessage {
  factory DelegatedAdminAccountLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DelegatedAdminAccountLimitExceededException._();

  factory DelegatedAdminAccountLimitExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DelegatedAdminAccountLimitExceededException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DelegatedAdminAccountLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegatedAdminAccountLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegatedAdminAccountLimitExceededException copyWith(
          void Function(DelegatedAdminAccountLimitExceededException) updates) =>
      super.copyWith((message) =>
              updates(message as DelegatedAdminAccountLimitExceededException))
          as DelegatedAdminAccountLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DelegatedAdminAccountLimitExceededException create() =>
      DelegatedAdminAccountLimitExceededException._();
  @$core.override
  DelegatedAdminAccountLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DelegatedAdminAccountLimitExceededException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DelegatedAdminAccountLimitExceededException>(create);
  static DelegatedAdminAccountLimitExceededException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class DeleteChannelRequest extends $pb.GeneratedMessage {
  factory DeleteChannelRequest({
    $core.String? channel,
  }) {
    final result = create();
    if (channel != null) result.channel = channel;
    return result;
  }

  DeleteChannelRequest._();

  factory DeleteChannelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteChannelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteChannelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(90016325, _omitFieldNames ? '' : 'channel')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChannelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChannelRequest copyWith(void Function(DeleteChannelRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteChannelRequest))
          as DeleteChannelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteChannelRequest create() => DeleteChannelRequest._();
  @$core.override
  DeleteChannelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteChannelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteChannelRequest>(create);
  static DeleteChannelRequest? _defaultInstance;

  @$pb.TagNumber(90016325)
  $core.String get channel => $_getSZ(0);
  @$pb.TagNumber(90016325)
  set channel($core.String value) => $_setString(0, value);
  @$pb.TagNumber(90016325)
  $core.bool hasChannel() => $_has(0);
  @$pb.TagNumber(90016325)
  void clearChannel() => $_clearField(90016325);
}

class DeleteChannelResponse extends $pb.GeneratedMessage {
  factory DeleteChannelResponse() => create();

  DeleteChannelResponse._();

  factory DeleteChannelResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteChannelResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteChannelResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChannelResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChannelResponse copyWith(
          void Function(DeleteChannelResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteChannelResponse))
          as DeleteChannelResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteChannelResponse create() => DeleteChannelResponse._();
  @$core.override
  DeleteChannelResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteChannelResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteChannelResponse>(create);
  static DeleteChannelResponse? _defaultInstance;
}

class DeleteDashboardRequest extends $pb.GeneratedMessage {
  factory DeleteDashboardRequest({
    $core.String? dashboardid,
  }) {
    final result = create();
    if (dashboardid != null) result.dashboardid = dashboardid;
    return result;
  }

  DeleteDashboardRequest._();

  factory DeleteDashboardRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDashboardRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDashboardRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(430615703, _omitFieldNames ? '' : 'dashboardid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardRequest copyWith(
          void Function(DeleteDashboardRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteDashboardRequest))
          as DeleteDashboardRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDashboardRequest create() => DeleteDashboardRequest._();
  @$core.override
  DeleteDashboardRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDashboardRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDashboardRequest>(create);
  static DeleteDashboardRequest? _defaultInstance;

  @$pb.TagNumber(430615703)
  $core.String get dashboardid => $_getSZ(0);
  @$pb.TagNumber(430615703)
  set dashboardid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(430615703)
  $core.bool hasDashboardid() => $_has(0);
  @$pb.TagNumber(430615703)
  void clearDashboardid() => $_clearField(430615703);
}

class DeleteDashboardResponse extends $pb.GeneratedMessage {
  factory DeleteDashboardResponse() => create();

  DeleteDashboardResponse._();

  factory DeleteDashboardResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDashboardResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDashboardResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardResponse copyWith(
          void Function(DeleteDashboardResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteDashboardResponse))
          as DeleteDashboardResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDashboardResponse create() => DeleteDashboardResponse._();
  @$core.override
  DeleteDashboardResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDashboardResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDashboardResponse>(create);
  static DeleteDashboardResponse? _defaultInstance;
}

class DeleteEventDataStoreRequest extends $pb.GeneratedMessage {
  factory DeleteEventDataStoreRequest({
    $core.String? eventdatastore,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    return result;
  }

  DeleteEventDataStoreRequest._();

  factory DeleteEventDataStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteEventDataStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteEventDataStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEventDataStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEventDataStoreRequest copyWith(
          void Function(DeleteEventDataStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteEventDataStoreRequest))
          as DeleteEventDataStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteEventDataStoreRequest create() =>
      DeleteEventDataStoreRequest._();
  @$core.override
  DeleteEventDataStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteEventDataStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteEventDataStoreRequest>(create);
  static DeleteEventDataStoreRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);
}

class DeleteEventDataStoreResponse extends $pb.GeneratedMessage {
  factory DeleteEventDataStoreResponse() => create();

  DeleteEventDataStoreResponse._();

  factory DeleteEventDataStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteEventDataStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteEventDataStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEventDataStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEventDataStoreResponse copyWith(
          void Function(DeleteEventDataStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteEventDataStoreResponse))
          as DeleteEventDataStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteEventDataStoreResponse create() =>
      DeleteEventDataStoreResponse._();
  @$core.override
  DeleteEventDataStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteEventDataStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteEventDataStoreResponse>(create);
  static DeleteEventDataStoreResponse? _defaultInstance;
}

class DeleteResourcePolicyRequest extends $pb.GeneratedMessage {
  factory DeleteResourcePolicyRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class DeleteResourcePolicyResponse extends $pb.GeneratedMessage {
  factory DeleteResourcePolicyResponse() => create();

  DeleteResourcePolicyResponse._();

  factory DeleteResourcePolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteResourcePolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteResourcePolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
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
}

class DeleteTrailRequest extends $pb.GeneratedMessage {
  factory DeleteTrailRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeleteTrailRequest._();

  factory DeleteTrailRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTrailRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTrailRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrailRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrailRequest copyWith(void Function(DeleteTrailRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteTrailRequest))
          as DeleteTrailRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTrailRequest create() => DeleteTrailRequest._();
  @$core.override
  DeleteTrailRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTrailRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTrailRequest>(create);
  static DeleteTrailRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeleteTrailResponse extends $pb.GeneratedMessage {
  factory DeleteTrailResponse() => create();

  DeleteTrailResponse._();

  factory DeleteTrailResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTrailResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTrailResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrailResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrailResponse copyWith(void Function(DeleteTrailResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteTrailResponse))
          as DeleteTrailResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTrailResponse create() => DeleteTrailResponse._();
  @$core.override
  DeleteTrailResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTrailResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTrailResponse>(create);
  static DeleteTrailResponse? _defaultInstance;
}

class DeregisterOrganizationDelegatedAdminRequest extends $pb.GeneratedMessage {
  factory DeregisterOrganizationDelegatedAdminRequest({
    $core.String? delegatedadminaccountid,
  }) {
    final result = create();
    if (delegatedadminaccountid != null)
      result.delegatedadminaccountid = delegatedadminaccountid;
    return result;
  }

  DeregisterOrganizationDelegatedAdminRequest._();

  factory DeregisterOrganizationDelegatedAdminRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeregisterOrganizationDelegatedAdminRequest.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeregisterOrganizationDelegatedAdminRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(104507462, _omitFieldNames ? '' : 'delegatedadminaccountid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeregisterOrganizationDelegatedAdminRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeregisterOrganizationDelegatedAdminRequest copyWith(
          void Function(DeregisterOrganizationDelegatedAdminRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeregisterOrganizationDelegatedAdminRequest))
          as DeregisterOrganizationDelegatedAdminRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeregisterOrganizationDelegatedAdminRequest create() =>
      DeregisterOrganizationDelegatedAdminRequest._();
  @$core.override
  DeregisterOrganizationDelegatedAdminRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeregisterOrganizationDelegatedAdminRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeregisterOrganizationDelegatedAdminRequest>(create);
  static DeregisterOrganizationDelegatedAdminRequest? _defaultInstance;

  @$pb.TagNumber(104507462)
  $core.String get delegatedadminaccountid => $_getSZ(0);
  @$pb.TagNumber(104507462)
  set delegatedadminaccountid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(104507462)
  $core.bool hasDelegatedadminaccountid() => $_has(0);
  @$pb.TagNumber(104507462)
  void clearDelegatedadminaccountid() => $_clearField(104507462);
}

class DeregisterOrganizationDelegatedAdminResponse
    extends $pb.GeneratedMessage {
  factory DeregisterOrganizationDelegatedAdminResponse() => create();

  DeregisterOrganizationDelegatedAdminResponse._();

  factory DeregisterOrganizationDelegatedAdminResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeregisterOrganizationDelegatedAdminResponse.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeregisterOrganizationDelegatedAdminResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeregisterOrganizationDelegatedAdminResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeregisterOrganizationDelegatedAdminResponse copyWith(
          void Function(DeregisterOrganizationDelegatedAdminResponse)
              updates) =>
      super.copyWith((message) =>
              updates(message as DeregisterOrganizationDelegatedAdminResponse))
          as DeregisterOrganizationDelegatedAdminResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeregisterOrganizationDelegatedAdminResponse create() =>
      DeregisterOrganizationDelegatedAdminResponse._();
  @$core.override
  DeregisterOrganizationDelegatedAdminResponse createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static DeregisterOrganizationDelegatedAdminResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeregisterOrganizationDelegatedAdminResponse>(create);
  static DeregisterOrganizationDelegatedAdminResponse? _defaultInstance;
}

class DescribeQueryRequest extends $pb.GeneratedMessage {
  factory DescribeQueryRequest({
    $core.String? refreshid,
    $core.String? queryid,
    $core.String? eventdatastore,
    $core.String? eventdatastoreowneraccountid,
    $core.String? queryalias,
  }) {
    final result = create();
    if (refreshid != null) result.refreshid = refreshid;
    if (queryid != null) result.queryid = queryid;
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    if (queryalias != null) result.queryalias = queryalias;
    return result;
  }

  DescribeQueryRequest._();

  factory DescribeQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeQueryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(57407610, _omitFieldNames ? '' : 'refreshid')
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..aOS(512762270, _omitFieldNames ? '' : 'queryalias')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeQueryRequest copyWith(void Function(DescribeQueryRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeQueryRequest))
          as DescribeQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeQueryRequest create() => DescribeQueryRequest._();
  @$core.override
  DescribeQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeQueryRequest>(create);
  static DescribeQueryRequest? _defaultInstance;

  @$pb.TagNumber(57407610)
  $core.String get refreshid => $_getSZ(0);
  @$pb.TagNumber(57407610)
  set refreshid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(57407610)
  $core.bool hasRefreshid() => $_has(0);
  @$pb.TagNumber(57407610)
  void clearRefreshid() => $_clearField(57407610);

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(1);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(1);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(2);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(2, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(2);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(3);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(3);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);

  @$pb.TagNumber(512762270)
  $core.String get queryalias => $_getSZ(4);
  @$pb.TagNumber(512762270)
  set queryalias($core.String value) => $_setString(4, value);
  @$pb.TagNumber(512762270)
  $core.bool hasQueryalias() => $_has(4);
  @$pb.TagNumber(512762270)
  void clearQueryalias() => $_clearField(512762270);
}

class DescribeQueryResponse extends $pb.GeneratedMessage {
  factory DescribeQueryResponse({
    $core.String? queryid,
    $core.String? deliverys3uri,
    QueryStatisticsForDescribeQuery? querystatistics,
    $core.String? prompt,
    QueryStatus? querystatus,
    $core.String? querystring,
    $core.String? eventdatastoreowneraccountid,
    DeliveryStatus? deliverystatus,
    $core.String? errormessage,
  }) {
    final result = create();
    if (queryid != null) result.queryid = queryid;
    if (deliverys3uri != null) result.deliverys3uri = deliverys3uri;
    if (querystatistics != null) result.querystatistics = querystatistics;
    if (prompt != null) result.prompt = prompt;
    if (querystatus != null) result.querystatus = querystatus;
    if (querystring != null) result.querystring = querystring;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    if (deliverystatus != null) result.deliverystatus = deliverystatus;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  DescribeQueryResponse._();

  factory DescribeQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeQueryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..aOS(230884460, _omitFieldNames ? '' : 'deliverys3uri')
    ..aOM<QueryStatisticsForDescribeQuery>(
        260794841, _omitFieldNames ? '' : 'querystatistics',
        subBuilder: QueryStatisticsForDescribeQuery.create)
    ..aOS(263974748, _omitFieldNames ? '' : 'prompt')
    ..aE<QueryStatus>(367016406, _omitFieldNames ? '' : 'querystatus',
        enumValues: QueryStatus.values)
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..aE<DeliveryStatus>(483265504, _omitFieldNames ? '' : 'deliverystatus',
        enumValues: DeliveryStatus.values)
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeQueryResponse copyWith(
          void Function(DescribeQueryResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeQueryResponse))
          as DescribeQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeQueryResponse create() => DescribeQueryResponse._();
  @$core.override
  DescribeQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeQueryResponse>(create);
  static DescribeQueryResponse? _defaultInstance;

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(0);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(0);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(230884460)
  $core.String get deliverys3uri => $_getSZ(1);
  @$pb.TagNumber(230884460)
  set deliverys3uri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(230884460)
  $core.bool hasDeliverys3uri() => $_has(1);
  @$pb.TagNumber(230884460)
  void clearDeliverys3uri() => $_clearField(230884460);

  @$pb.TagNumber(260794841)
  QueryStatisticsForDescribeQuery get querystatistics => $_getN(2);
  @$pb.TagNumber(260794841)
  set querystatistics(QueryStatisticsForDescribeQuery value) =>
      $_setField(260794841, value);
  @$pb.TagNumber(260794841)
  $core.bool hasQuerystatistics() => $_has(2);
  @$pb.TagNumber(260794841)
  void clearQuerystatistics() => $_clearField(260794841);
  @$pb.TagNumber(260794841)
  QueryStatisticsForDescribeQuery ensureQuerystatistics() => $_ensure(2);

  @$pb.TagNumber(263974748)
  $core.String get prompt => $_getSZ(3);
  @$pb.TagNumber(263974748)
  set prompt($core.String value) => $_setString(3, value);
  @$pb.TagNumber(263974748)
  $core.bool hasPrompt() => $_has(3);
  @$pb.TagNumber(263974748)
  void clearPrompt() => $_clearField(263974748);

  @$pb.TagNumber(367016406)
  QueryStatus get querystatus => $_getN(4);
  @$pb.TagNumber(367016406)
  set querystatus(QueryStatus value) => $_setField(367016406, value);
  @$pb.TagNumber(367016406)
  $core.bool hasQuerystatus() => $_has(4);
  @$pb.TagNumber(367016406)
  void clearQuerystatus() => $_clearField(367016406);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(5);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(5, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(5);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(6);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(6);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);

  @$pb.TagNumber(483265504)
  DeliveryStatus get deliverystatus => $_getN(7);
  @$pb.TagNumber(483265504)
  set deliverystatus(DeliveryStatus value) => $_setField(483265504, value);
  @$pb.TagNumber(483265504)
  $core.bool hasDeliverystatus() => $_has(7);
  @$pb.TagNumber(483265504)
  void clearDeliverystatus() => $_clearField(483265504);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(8);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(8, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(8);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class DescribeTrailsRequest extends $pb.GeneratedMessage {
  factory DescribeTrailsRequest({
    $core.bool? includeshadowtrails,
    $core.Iterable<$core.String>? trailnamelist,
  }) {
    final result = create();
    if (includeshadowtrails != null)
      result.includeshadowtrails = includeshadowtrails;
    if (trailnamelist != null) result.trailnamelist.addAll(trailnamelist);
    return result;
  }

  DescribeTrailsRequest._();

  factory DescribeTrailsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTrailsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTrailsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOB(200547819, _omitFieldNames ? '' : 'includeshadowtrails')
    ..pPS(529562769, _omitFieldNames ? '' : 'trailnamelist')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTrailsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTrailsRequest copyWith(
          void Function(DescribeTrailsRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeTrailsRequest))
          as DescribeTrailsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTrailsRequest create() => DescribeTrailsRequest._();
  @$core.override
  DescribeTrailsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTrailsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTrailsRequest>(create);
  static DescribeTrailsRequest? _defaultInstance;

  @$pb.TagNumber(200547819)
  $core.bool get includeshadowtrails => $_getBF(0);
  @$pb.TagNumber(200547819)
  set includeshadowtrails($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(200547819)
  $core.bool hasIncludeshadowtrails() => $_has(0);
  @$pb.TagNumber(200547819)
  void clearIncludeshadowtrails() => $_clearField(200547819);

  @$pb.TagNumber(529562769)
  $pb.PbList<$core.String> get trailnamelist => $_getList(1);
}

class DescribeTrailsResponse extends $pb.GeneratedMessage {
  factory DescribeTrailsResponse({
    $core.Iterable<Trail>? traillist,
  }) {
    final result = create();
    if (traillist != null) result.traillist.addAll(traillist);
    return result;
  }

  DescribeTrailsResponse._();

  factory DescribeTrailsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTrailsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTrailsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<Trail>(508517484, _omitFieldNames ? '' : 'traillist',
        subBuilder: Trail.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTrailsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTrailsResponse copyWith(
          void Function(DescribeTrailsResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeTrailsResponse))
          as DescribeTrailsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTrailsResponse create() => DescribeTrailsResponse._();
  @$core.override
  DescribeTrailsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTrailsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTrailsResponse>(create);
  static DescribeTrailsResponse? _defaultInstance;

  @$pb.TagNumber(508517484)
  $pb.PbList<Trail> get traillist => $_getList(0);
}

class Destination extends $pb.GeneratedMessage {
  factory Destination({
    DestinationType? type,
    $core.String? location,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (location != null) result.location = location;
    return result;
  }

  Destination._();

  factory Destination.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Destination.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Destination',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<DestinationType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DestinationType.values)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Destination clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Destination copyWith(void Function(Destination) updates) =>
      super.copyWith((message) => updates(message as Destination))
          as Destination;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Destination create() => Destination._();
  @$core.override
  Destination createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Destination getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Destination>(create);
  static Destination? _defaultInstance;

  @$pb.TagNumber(290836590)
  DestinationType get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(DestinationType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class DisableFederationRequest extends $pb.GeneratedMessage {
  factory DisableFederationRequest({
    $core.String? eventdatastore,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    return result;
  }

  DisableFederationRequest._();

  factory DisableFederationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableFederationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableFederationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableFederationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableFederationRequest copyWith(
          void Function(DisableFederationRequest) updates) =>
      super.copyWith((message) => updates(message as DisableFederationRequest))
          as DisableFederationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableFederationRequest create() => DisableFederationRequest._();
  @$core.override
  DisableFederationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableFederationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableFederationRequest>(create);
  static DisableFederationRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);
}

class DisableFederationResponse extends $pb.GeneratedMessage {
  factory DisableFederationResponse({
    FederationStatus? federationstatus,
    $core.String? eventdatastorearn,
  }) {
    final result = create();
    if (federationstatus != null) result.federationstatus = federationstatus;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    return result;
  }

  DisableFederationResponse._();

  factory DisableFederationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableFederationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableFederationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<FederationStatus>(146235383, _omitFieldNames ? '' : 'federationstatus',
        enumValues: FederationStatus.values)
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableFederationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableFederationResponse copyWith(
          void Function(DisableFederationResponse) updates) =>
      super.copyWith((message) => updates(message as DisableFederationResponse))
          as DisableFederationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableFederationResponse create() => DisableFederationResponse._();
  @$core.override
  DisableFederationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableFederationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableFederationResponse>(create);
  static DisableFederationResponse? _defaultInstance;

  @$pb.TagNumber(146235383)
  FederationStatus get federationstatus => $_getN(0);
  @$pb.TagNumber(146235383)
  set federationstatus(FederationStatus value) => $_setField(146235383, value);
  @$pb.TagNumber(146235383)
  $core.bool hasFederationstatus() => $_has(0);
  @$pb.TagNumber(146235383)
  void clearFederationstatus() => $_clearField(146235383);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(1);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(1);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);
}

class EnableFederationRequest extends $pb.GeneratedMessage {
  factory EnableFederationRequest({
    $core.String? eventdatastore,
    $core.String? federationrolearn,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (federationrolearn != null) result.federationrolearn = federationrolearn;
    return result;
  }

  EnableFederationRequest._();

  factory EnableFederationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableFederationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableFederationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(504464364, _omitFieldNames ? '' : 'federationrolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableFederationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableFederationRequest copyWith(
          void Function(EnableFederationRequest) updates) =>
      super.copyWith((message) => updates(message as EnableFederationRequest))
          as EnableFederationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableFederationRequest create() => EnableFederationRequest._();
  @$core.override
  EnableFederationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableFederationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableFederationRequest>(create);
  static EnableFederationRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(504464364)
  $core.String get federationrolearn => $_getSZ(1);
  @$pb.TagNumber(504464364)
  set federationrolearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(504464364)
  $core.bool hasFederationrolearn() => $_has(1);
  @$pb.TagNumber(504464364)
  void clearFederationrolearn() => $_clearField(504464364);
}

class EnableFederationResponse extends $pb.GeneratedMessage {
  factory EnableFederationResponse({
    FederationStatus? federationstatus,
    $core.String? eventdatastorearn,
    $core.String? federationrolearn,
  }) {
    final result = create();
    if (federationstatus != null) result.federationstatus = federationstatus;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (federationrolearn != null) result.federationrolearn = federationrolearn;
    return result;
  }

  EnableFederationResponse._();

  factory EnableFederationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableFederationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableFederationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<FederationStatus>(146235383, _omitFieldNames ? '' : 'federationstatus',
        enumValues: FederationStatus.values)
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(504464364, _omitFieldNames ? '' : 'federationrolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableFederationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableFederationResponse copyWith(
          void Function(EnableFederationResponse) updates) =>
      super.copyWith((message) => updates(message as EnableFederationResponse))
          as EnableFederationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableFederationResponse create() => EnableFederationResponse._();
  @$core.override
  EnableFederationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableFederationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableFederationResponse>(create);
  static EnableFederationResponse? _defaultInstance;

  @$pb.TagNumber(146235383)
  FederationStatus get federationstatus => $_getN(0);
  @$pb.TagNumber(146235383)
  set federationstatus(FederationStatus value) => $_setField(146235383, value);
  @$pb.TagNumber(146235383)
  $core.bool hasFederationstatus() => $_has(0);
  @$pb.TagNumber(146235383)
  void clearFederationstatus() => $_clearField(146235383);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(1);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(1);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(504464364)
  $core.String get federationrolearn => $_getSZ(2);
  @$pb.TagNumber(504464364)
  set federationrolearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(504464364)
  $core.bool hasFederationrolearn() => $_has(2);
  @$pb.TagNumber(504464364)
  void clearFederationrolearn() => $_clearField(504464364);
}

class Event extends $pb.GeneratedMessage {
  factory Event({
    $core.String? eventsource,
    $core.String? readonly,
    $core.String? eventtime,
    $core.String? eventname,
    $core.String? cloudtrailevent,
    $core.Iterable<Resource>? resources,
    $core.String? eventid,
    $core.String? accesskeyid,
    $core.String? username,
  }) {
    final result = create();
    if (eventsource != null) result.eventsource = eventsource;
    if (readonly != null) result.readonly = readonly;
    if (eventtime != null) result.eventtime = eventtime;
    if (eventname != null) result.eventname = eventname;
    if (cloudtrailevent != null) result.cloudtrailevent = cloudtrailevent;
    if (resources != null) result.resources.addAll(resources);
    if (eventid != null) result.eventid = eventid;
    if (accesskeyid != null) result.accesskeyid = accesskeyid;
    if (username != null) result.username = username;
    return result;
  }

  Event._();

  factory Event.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Event.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Event',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(37841339, _omitFieldNames ? '' : 'eventsource')
    ..aOS(129039352, _omitFieldNames ? '' : 'readonly')
    ..aOS(222361647, _omitFieldNames ? '' : 'eventtime')
    ..aOS(264746781, _omitFieldNames ? '' : 'eventname')
    ..aOS(344297255, _omitFieldNames ? '' : 'cloudtrailevent')
    ..pPM<Resource>(358918291, _omitFieldNames ? '' : 'resources',
        subBuilder: Resource.create)
    ..aOS(376916819, _omitFieldNames ? '' : 'eventid')
    ..aOS(453893024, _omitFieldNames ? '' : 'accesskeyid')
    ..aOS(470340826, _omitFieldNames ? '' : 'username')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Event clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Event copyWith(void Function(Event) updates) =>
      super.copyWith((message) => updates(message as Event)) as Event;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Event create() => Event._();
  @$core.override
  Event createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Event getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Event>(create);
  static Event? _defaultInstance;

  @$pb.TagNumber(37841339)
  $core.String get eventsource => $_getSZ(0);
  @$pb.TagNumber(37841339)
  set eventsource($core.String value) => $_setString(0, value);
  @$pb.TagNumber(37841339)
  $core.bool hasEventsource() => $_has(0);
  @$pb.TagNumber(37841339)
  void clearEventsource() => $_clearField(37841339);

  @$pb.TagNumber(129039352)
  $core.String get readonly => $_getSZ(1);
  @$pb.TagNumber(129039352)
  set readonly($core.String value) => $_setString(1, value);
  @$pb.TagNumber(129039352)
  $core.bool hasReadonly() => $_has(1);
  @$pb.TagNumber(129039352)
  void clearReadonly() => $_clearField(129039352);

  @$pb.TagNumber(222361647)
  $core.String get eventtime => $_getSZ(2);
  @$pb.TagNumber(222361647)
  set eventtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(222361647)
  $core.bool hasEventtime() => $_has(2);
  @$pb.TagNumber(222361647)
  void clearEventtime() => $_clearField(222361647);

  @$pb.TagNumber(264746781)
  $core.String get eventname => $_getSZ(3);
  @$pb.TagNumber(264746781)
  set eventname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(264746781)
  $core.bool hasEventname() => $_has(3);
  @$pb.TagNumber(264746781)
  void clearEventname() => $_clearField(264746781);

  @$pb.TagNumber(344297255)
  $core.String get cloudtrailevent => $_getSZ(4);
  @$pb.TagNumber(344297255)
  set cloudtrailevent($core.String value) => $_setString(4, value);
  @$pb.TagNumber(344297255)
  $core.bool hasCloudtrailevent() => $_has(4);
  @$pb.TagNumber(344297255)
  void clearCloudtrailevent() => $_clearField(344297255);

  @$pb.TagNumber(358918291)
  $pb.PbList<Resource> get resources => $_getList(5);

  @$pb.TagNumber(376916819)
  $core.String get eventid => $_getSZ(6);
  @$pb.TagNumber(376916819)
  set eventid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(376916819)
  $core.bool hasEventid() => $_has(6);
  @$pb.TagNumber(376916819)
  void clearEventid() => $_clearField(376916819);

  @$pb.TagNumber(453893024)
  $core.String get accesskeyid => $_getSZ(7);
  @$pb.TagNumber(453893024)
  set accesskeyid($core.String value) => $_setString(7, value);
  @$pb.TagNumber(453893024)
  $core.bool hasAccesskeyid() => $_has(7);
  @$pb.TagNumber(453893024)
  void clearAccesskeyid() => $_clearField(453893024);

  @$pb.TagNumber(470340826)
  $core.String get username => $_getSZ(8);
  @$pb.TagNumber(470340826)
  set username($core.String value) => $_setString(8, value);
  @$pb.TagNumber(470340826)
  $core.bool hasUsername() => $_has(8);
  @$pb.TagNumber(470340826)
  void clearUsername() => $_clearField(470340826);
}

class EventDataStore extends $pb.GeneratedMessage {
  factory EventDataStore({
    EventDataStoreStatus? status,
    $core.bool? multiregionenabled,
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? updatedtimestamp,
    $core.int? retentionperiod,
    $core.String? name,
    $core.String? eventdatastorearn,
    $core.String? createdtimestamp,
    $core.bool? terminationprotectionenabled,
    $core.bool? organizationenabled,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (multiregionenabled != null)
      result.multiregionenabled = multiregionenabled;
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (retentionperiod != null) result.retentionperiod = retentionperiod;
    if (name != null) result.name = name;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (organizationenabled != null)
      result.organizationenabled = organizationenabled;
    return result;
  }

  EventDataStore._();

  factory EventDataStore.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStore.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStore',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<EventDataStoreStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: EventDataStoreStatus.values)
    ..aOB(20620620, _omitFieldNames ? '' : 'multiregionenabled')
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aI(196383721, _omitFieldNames ? '' : 'retentionperiod')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOB(480171176, _omitFieldNames ? '' : 'organizationenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStore clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStore copyWith(void Function(EventDataStore) updates) =>
      super.copyWith((message) => updates(message as EventDataStore))
          as EventDataStore;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStore create() => EventDataStore._();
  @$core.override
  EventDataStore createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStore getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EventDataStore>(create);
  static EventDataStore? _defaultInstance;

  @$pb.TagNumber(6222352)
  EventDataStoreStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(EventDataStoreStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(20620620)
  $core.bool get multiregionenabled => $_getBF(1);
  @$pb.TagNumber(20620620)
  set multiregionenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(20620620)
  $core.bool hasMultiregionenabled() => $_has(1);
  @$pb.TagNumber(20620620)
  void clearMultiregionenabled() => $_clearField(20620620);

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(2);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(3);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(3);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(196383721)
  $core.int get retentionperiod => $_getIZ(4);
  @$pb.TagNumber(196383721)
  set retentionperiod($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(196383721)
  $core.bool hasRetentionperiod() => $_has(4);
  @$pb.TagNumber(196383721)
  void clearRetentionperiod() => $_clearField(196383721);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(6);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(6);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(7);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(7, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(7);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(8);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(8);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(480171176)
  $core.bool get organizationenabled => $_getBF(9);
  @$pb.TagNumber(480171176)
  set organizationenabled($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(480171176)
  $core.bool hasOrganizationenabled() => $_has(9);
  @$pb.TagNumber(480171176)
  void clearOrganizationenabled() => $_clearField(480171176);
}

class EventDataStoreARNInvalidException extends $pb.GeneratedMessage {
  factory EventDataStoreARNInvalidException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EventDataStoreARNInvalidException._();

  factory EventDataStoreARNInvalidException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStoreARNInvalidException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStoreARNInvalidException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreARNInvalidException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreARNInvalidException copyWith(
          void Function(EventDataStoreARNInvalidException) updates) =>
      super.copyWith((message) =>
              updates(message as EventDataStoreARNInvalidException))
          as EventDataStoreARNInvalidException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStoreARNInvalidException create() =>
      EventDataStoreARNInvalidException._();
  @$core.override
  EventDataStoreARNInvalidException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStoreARNInvalidException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EventDataStoreARNInvalidException>(
          create);
  static EventDataStoreARNInvalidException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class EventDataStoreAlreadyExistsException extends $pb.GeneratedMessage {
  factory EventDataStoreAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EventDataStoreAlreadyExistsException._();

  factory EventDataStoreAlreadyExistsException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStoreAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStoreAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreAlreadyExistsException copyWith(
          void Function(EventDataStoreAlreadyExistsException) updates) =>
      super.copyWith((message) =>
              updates(message as EventDataStoreAlreadyExistsException))
          as EventDataStoreAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStoreAlreadyExistsException create() =>
      EventDataStoreAlreadyExistsException._();
  @$core.override
  EventDataStoreAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStoreAlreadyExistsException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          EventDataStoreAlreadyExistsException>(create);
  static EventDataStoreAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class EventDataStoreFederationEnabledException extends $pb.GeneratedMessage {
  factory EventDataStoreFederationEnabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EventDataStoreFederationEnabledException._();

  factory EventDataStoreFederationEnabledException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStoreFederationEnabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStoreFederationEnabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreFederationEnabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreFederationEnabledException copyWith(
          void Function(EventDataStoreFederationEnabledException) updates) =>
      super.copyWith((message) =>
              updates(message as EventDataStoreFederationEnabledException))
          as EventDataStoreFederationEnabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStoreFederationEnabledException create() =>
      EventDataStoreFederationEnabledException._();
  @$core.override
  EventDataStoreFederationEnabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStoreFederationEnabledException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          EventDataStoreFederationEnabledException>(create);
  static EventDataStoreFederationEnabledException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class EventDataStoreHasOngoingImportException extends $pb.GeneratedMessage {
  factory EventDataStoreHasOngoingImportException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EventDataStoreHasOngoingImportException._();

  factory EventDataStoreHasOngoingImportException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStoreHasOngoingImportException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStoreHasOngoingImportException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreHasOngoingImportException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreHasOngoingImportException copyWith(
          void Function(EventDataStoreHasOngoingImportException) updates) =>
      super.copyWith((message) =>
              updates(message as EventDataStoreHasOngoingImportException))
          as EventDataStoreHasOngoingImportException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStoreHasOngoingImportException create() =>
      EventDataStoreHasOngoingImportException._();
  @$core.override
  EventDataStoreHasOngoingImportException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStoreHasOngoingImportException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          EventDataStoreHasOngoingImportException>(create);
  static EventDataStoreHasOngoingImportException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class EventDataStoreMaxLimitExceededException extends $pb.GeneratedMessage {
  factory EventDataStoreMaxLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EventDataStoreMaxLimitExceededException._();

  factory EventDataStoreMaxLimitExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStoreMaxLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStoreMaxLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreMaxLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreMaxLimitExceededException copyWith(
          void Function(EventDataStoreMaxLimitExceededException) updates) =>
      super.copyWith((message) =>
              updates(message as EventDataStoreMaxLimitExceededException))
          as EventDataStoreMaxLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStoreMaxLimitExceededException create() =>
      EventDataStoreMaxLimitExceededException._();
  @$core.override
  EventDataStoreMaxLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStoreMaxLimitExceededException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          EventDataStoreMaxLimitExceededException>(create);
  static EventDataStoreMaxLimitExceededException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class EventDataStoreNotFoundException extends $pb.GeneratedMessage {
  factory EventDataStoreNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EventDataStoreNotFoundException._();

  factory EventDataStoreNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStoreNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStoreNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreNotFoundException copyWith(
          void Function(EventDataStoreNotFoundException) updates) =>
      super.copyWith(
              (message) => updates(message as EventDataStoreNotFoundException))
          as EventDataStoreNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStoreNotFoundException create() =>
      EventDataStoreNotFoundException._();
  @$core.override
  EventDataStoreNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStoreNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EventDataStoreNotFoundException>(
          create);
  static EventDataStoreNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class EventDataStoreTerminationProtectedException extends $pb.GeneratedMessage {
  factory EventDataStoreTerminationProtectedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EventDataStoreTerminationProtectedException._();

  factory EventDataStoreTerminationProtectedException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventDataStoreTerminationProtectedException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventDataStoreTerminationProtectedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreTerminationProtectedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventDataStoreTerminationProtectedException copyWith(
          void Function(EventDataStoreTerminationProtectedException) updates) =>
      super.copyWith((message) =>
              updates(message as EventDataStoreTerminationProtectedException))
          as EventDataStoreTerminationProtectedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventDataStoreTerminationProtectedException create() =>
      EventDataStoreTerminationProtectedException._();
  @$core.override
  EventDataStoreTerminationProtectedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventDataStoreTerminationProtectedException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          EventDataStoreTerminationProtectedException>(create);
  static EventDataStoreTerminationProtectedException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class EventSelector extends $pb.GeneratedMessage {
  factory EventSelector({
    $core.Iterable<DataResource>? dataresources,
    $core.bool? includemanagementevents,
    $core.Iterable<$core.String>? excludemanagementeventsources,
    ReadWriteType? readwritetype,
  }) {
    final result = create();
    if (dataresources != null) result.dataresources.addAll(dataresources);
    if (includemanagementevents != null)
      result.includemanagementevents = includemanagementevents;
    if (excludemanagementeventsources != null)
      result.excludemanagementeventsources
          .addAll(excludemanagementeventsources);
    if (readwritetype != null) result.readwritetype = readwritetype;
    return result;
  }

  EventSelector._();

  factory EventSelector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventSelector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventSelector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<DataResource>(126123155, _omitFieldNames ? '' : 'dataresources',
        subBuilder: DataResource.create)
    ..aOB(215824550, _omitFieldNames ? '' : 'includemanagementevents')
    ..pPS(225676985, _omitFieldNames ? '' : 'excludemanagementeventsources')
    ..aE<ReadWriteType>(296653333, _omitFieldNames ? '' : 'readwritetype',
        enumValues: ReadWriteType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventSelector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventSelector copyWith(void Function(EventSelector) updates) =>
      super.copyWith((message) => updates(message as EventSelector))
          as EventSelector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventSelector create() => EventSelector._();
  @$core.override
  EventSelector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventSelector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EventSelector>(create);
  static EventSelector? _defaultInstance;

  @$pb.TagNumber(126123155)
  $pb.PbList<DataResource> get dataresources => $_getList(0);

  @$pb.TagNumber(215824550)
  $core.bool get includemanagementevents => $_getBF(1);
  @$pb.TagNumber(215824550)
  set includemanagementevents($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(215824550)
  $core.bool hasIncludemanagementevents() => $_has(1);
  @$pb.TagNumber(215824550)
  void clearIncludemanagementevents() => $_clearField(215824550);

  @$pb.TagNumber(225676985)
  $pb.PbList<$core.String> get excludemanagementeventsources => $_getList(2);

  @$pb.TagNumber(296653333)
  ReadWriteType get readwritetype => $_getN(3);
  @$pb.TagNumber(296653333)
  set readwritetype(ReadWriteType value) => $_setField(296653333, value);
  @$pb.TagNumber(296653333)
  $core.bool hasReadwritetype() => $_has(3);
  @$pb.TagNumber(296653333)
  void clearReadwritetype() => $_clearField(296653333);
}

class GenerateQueryRequest extends $pb.GeneratedMessage {
  factory GenerateQueryRequest({
    $core.Iterable<$core.String>? eventdatastores,
    $core.String? prompt,
  }) {
    final result = create();
    if (eventdatastores != null) result.eventdatastores.addAll(eventdatastores);
    if (prompt != null) result.prompt = prompt;
    return result;
  }

  GenerateQueryRequest._();

  factory GenerateQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateQueryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPS(152154314, _omitFieldNames ? '' : 'eventdatastores')
    ..aOS(263974748, _omitFieldNames ? '' : 'prompt')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateQueryRequest copyWith(void Function(GenerateQueryRequest) updates) =>
      super.copyWith((message) => updates(message as GenerateQueryRequest))
          as GenerateQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateQueryRequest create() => GenerateQueryRequest._();
  @$core.override
  GenerateQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateQueryRequest>(create);
  static GenerateQueryRequest? _defaultInstance;

  @$pb.TagNumber(152154314)
  $pb.PbList<$core.String> get eventdatastores => $_getList(0);

  @$pb.TagNumber(263974748)
  $core.String get prompt => $_getSZ(1);
  @$pb.TagNumber(263974748)
  set prompt($core.String value) => $_setString(1, value);
  @$pb.TagNumber(263974748)
  $core.bool hasPrompt() => $_has(1);
  @$pb.TagNumber(263974748)
  void clearPrompt() => $_clearField(263974748);
}

class GenerateQueryResponse extends $pb.GeneratedMessage {
  factory GenerateQueryResponse({
    $core.String? querystatement,
    $core.String? eventdatastoreowneraccountid,
    $core.String? queryalias,
  }) {
    final result = create();
    if (querystatement != null) result.querystatement = querystatement;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    if (queryalias != null) result.queryalias = queryalias;
    return result;
  }

  GenerateQueryResponse._();

  factory GenerateQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateQueryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(340852217, _omitFieldNames ? '' : 'querystatement')
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..aOS(512762270, _omitFieldNames ? '' : 'queryalias')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateQueryResponse copyWith(
          void Function(GenerateQueryResponse) updates) =>
      super.copyWith((message) => updates(message as GenerateQueryResponse))
          as GenerateQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateQueryResponse create() => GenerateQueryResponse._();
  @$core.override
  GenerateQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateQueryResponse>(create);
  static GenerateQueryResponse? _defaultInstance;

  @$pb.TagNumber(340852217)
  $core.String get querystatement => $_getSZ(0);
  @$pb.TagNumber(340852217)
  set querystatement($core.String value) => $_setString(0, value);
  @$pb.TagNumber(340852217)
  $core.bool hasQuerystatement() => $_has(0);
  @$pb.TagNumber(340852217)
  void clearQuerystatement() => $_clearField(340852217);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(1);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(1);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);

  @$pb.TagNumber(512762270)
  $core.String get queryalias => $_getSZ(2);
  @$pb.TagNumber(512762270)
  set queryalias($core.String value) => $_setString(2, value);
  @$pb.TagNumber(512762270)
  $core.bool hasQueryalias() => $_has(2);
  @$pb.TagNumber(512762270)
  void clearQueryalias() => $_clearField(512762270);
}

class GenerateResponseException extends $pb.GeneratedMessage {
  factory GenerateResponseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  GenerateResponseException._();

  factory GenerateResponseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateResponseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateResponseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateResponseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateResponseException copyWith(
          void Function(GenerateResponseException) updates) =>
      super.copyWith((message) => updates(message as GenerateResponseException))
          as GenerateResponseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateResponseException create() => GenerateResponseException._();
  @$core.override
  GenerateResponseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateResponseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateResponseException>(create);
  static GenerateResponseException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class GetChannelRequest extends $pb.GeneratedMessage {
  factory GetChannelRequest({
    $core.String? channel,
  }) {
    final result = create();
    if (channel != null) result.channel = channel;
    return result;
  }

  GetChannelRequest._();

  factory GetChannelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChannelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChannelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(90016325, _omitFieldNames ? '' : 'channel')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChannelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChannelRequest copyWith(void Function(GetChannelRequest) updates) =>
      super.copyWith((message) => updates(message as GetChannelRequest))
          as GetChannelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChannelRequest create() => GetChannelRequest._();
  @$core.override
  GetChannelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChannelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChannelRequest>(create);
  static GetChannelRequest? _defaultInstance;

  @$pb.TagNumber(90016325)
  $core.String get channel => $_getSZ(0);
  @$pb.TagNumber(90016325)
  set channel($core.String value) => $_setString(0, value);
  @$pb.TagNumber(90016325)
  $core.bool hasChannel() => $_has(0);
  @$pb.TagNumber(90016325)
  void clearChannel() => $_clearField(90016325);
}

class GetChannelResponse extends $pb.GeneratedMessage {
  factory GetChannelResponse({
    $core.String? source,
    $core.String? channelarn,
    IngestionStatus? ingestionstatus,
    $core.String? name,
    $core.Iterable<Destination>? destinations,
    SourceConfig? sourceconfig,
  }) {
    final result = create();
    if (source != null) result.source = source;
    if (channelarn != null) result.channelarn = channelarn;
    if (ingestionstatus != null) result.ingestionstatus = ingestionstatus;
    if (name != null) result.name = name;
    if (destinations != null) result.destinations.addAll(destinations);
    if (sourceconfig != null) result.sourceconfig = sourceconfig;
    return result;
  }

  GetChannelResponse._();

  factory GetChannelResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChannelResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChannelResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(31630329, _omitFieldNames ? '' : 'source')
    ..aOS(99276476, _omitFieldNames ? '' : 'channelarn')
    ..aOM<IngestionStatus>(188310656, _omitFieldNames ? '' : 'ingestionstatus',
        subBuilder: IngestionStatus.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Destination>(404385861, _omitFieldNames ? '' : 'destinations',
        subBuilder: Destination.create)
    ..aOM<SourceConfig>(415255383, _omitFieldNames ? '' : 'sourceconfig',
        subBuilder: SourceConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChannelResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChannelResponse copyWith(void Function(GetChannelResponse) updates) =>
      super.copyWith((message) => updates(message as GetChannelResponse))
          as GetChannelResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChannelResponse create() => GetChannelResponse._();
  @$core.override
  GetChannelResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChannelResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChannelResponse>(create);
  static GetChannelResponse? _defaultInstance;

  @$pb.TagNumber(31630329)
  $core.String get source => $_getSZ(0);
  @$pb.TagNumber(31630329)
  set source($core.String value) => $_setString(0, value);
  @$pb.TagNumber(31630329)
  $core.bool hasSource() => $_has(0);
  @$pb.TagNumber(31630329)
  void clearSource() => $_clearField(31630329);

  @$pb.TagNumber(99276476)
  $core.String get channelarn => $_getSZ(1);
  @$pb.TagNumber(99276476)
  set channelarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(99276476)
  $core.bool hasChannelarn() => $_has(1);
  @$pb.TagNumber(99276476)
  void clearChannelarn() => $_clearField(99276476);

  @$pb.TagNumber(188310656)
  IngestionStatus get ingestionstatus => $_getN(2);
  @$pb.TagNumber(188310656)
  set ingestionstatus(IngestionStatus value) => $_setField(188310656, value);
  @$pb.TagNumber(188310656)
  $core.bool hasIngestionstatus() => $_has(2);
  @$pb.TagNumber(188310656)
  void clearIngestionstatus() => $_clearField(188310656);
  @$pb.TagNumber(188310656)
  IngestionStatus ensureIngestionstatus() => $_ensure(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(404385861)
  $pb.PbList<Destination> get destinations => $_getList(4);

  @$pb.TagNumber(415255383)
  SourceConfig get sourceconfig => $_getN(5);
  @$pb.TagNumber(415255383)
  set sourceconfig(SourceConfig value) => $_setField(415255383, value);
  @$pb.TagNumber(415255383)
  $core.bool hasSourceconfig() => $_has(5);
  @$pb.TagNumber(415255383)
  void clearSourceconfig() => $_clearField(415255383);
  @$pb.TagNumber(415255383)
  SourceConfig ensureSourceconfig() => $_ensure(5);
}

class GetDashboardRequest extends $pb.GeneratedMessage {
  factory GetDashboardRequest({
    $core.String? dashboardid,
  }) {
    final result = create();
    if (dashboardid != null) result.dashboardid = dashboardid;
    return result;
  }

  GetDashboardRequest._();

  factory GetDashboardRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDashboardRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDashboardRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(430615703, _omitFieldNames ? '' : 'dashboardid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardRequest copyWith(void Function(GetDashboardRequest) updates) =>
      super.copyWith((message) => updates(message as GetDashboardRequest))
          as GetDashboardRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDashboardRequest create() => GetDashboardRequest._();
  @$core.override
  GetDashboardRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDashboardRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDashboardRequest>(create);
  static GetDashboardRequest? _defaultInstance;

  @$pb.TagNumber(430615703)
  $core.String get dashboardid => $_getSZ(0);
  @$pb.TagNumber(430615703)
  set dashboardid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(430615703)
  $core.bool hasDashboardid() => $_has(0);
  @$pb.TagNumber(430615703)
  void clearDashboardid() => $_clearField(430615703);
}

class GetDashboardResponse extends $pb.GeneratedMessage {
  factory GetDashboardResponse({
    DashboardStatus? status,
    $core.String? updatedtimestamp,
    $core.String? dashboardarn,
    RefreshSchedule? refreshschedule,
    $core.String? lastrefreshid,
    DashboardType? type,
    $core.String? createdtimestamp,
    $core.bool? terminationprotectionenabled,
    $core.String? lastrefreshfailurereason,
    $core.Iterable<Widget>? widgets,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (dashboardarn != null) result.dashboardarn = dashboardarn;
    if (refreshschedule != null) result.refreshschedule = refreshschedule;
    if (lastrefreshid != null) result.lastrefreshid = lastrefreshid;
    if (type != null) result.type = type;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (lastrefreshfailurereason != null)
      result.lastrefreshfailurereason = lastrefreshfailurereason;
    if (widgets != null) result.widgets.addAll(widgets);
    return result;
  }

  GetDashboardResponse._();

  factory GetDashboardResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDashboardResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDashboardResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<DashboardStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: DashboardStatus.values)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOS(108051951, _omitFieldNames ? '' : 'dashboardarn')
    ..aOM<RefreshSchedule>(261773338, _omitFieldNames ? '' : 'refreshschedule',
        subBuilder: RefreshSchedule.create)
    ..aOS(272889110, _omitFieldNames ? '' : 'lastrefreshid')
    ..aE<DashboardType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DashboardType.values)
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOS(493889247, _omitFieldNames ? '' : 'lastrefreshfailurereason')
    ..pPM<Widget>(501826147, _omitFieldNames ? '' : 'widgets',
        subBuilder: Widget.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardResponse copyWith(void Function(GetDashboardResponse) updates) =>
      super.copyWith((message) => updates(message as GetDashboardResponse))
          as GetDashboardResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDashboardResponse create() => GetDashboardResponse._();
  @$core.override
  GetDashboardResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDashboardResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDashboardResponse>(create);
  static GetDashboardResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  DashboardStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(DashboardStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(1);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(1);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(108051951)
  $core.String get dashboardarn => $_getSZ(2);
  @$pb.TagNumber(108051951)
  set dashboardarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(108051951)
  $core.bool hasDashboardarn() => $_has(2);
  @$pb.TagNumber(108051951)
  void clearDashboardarn() => $_clearField(108051951);

  @$pb.TagNumber(261773338)
  RefreshSchedule get refreshschedule => $_getN(3);
  @$pb.TagNumber(261773338)
  set refreshschedule(RefreshSchedule value) => $_setField(261773338, value);
  @$pb.TagNumber(261773338)
  $core.bool hasRefreshschedule() => $_has(3);
  @$pb.TagNumber(261773338)
  void clearRefreshschedule() => $_clearField(261773338);
  @$pb.TagNumber(261773338)
  RefreshSchedule ensureRefreshschedule() => $_ensure(3);

  @$pb.TagNumber(272889110)
  $core.String get lastrefreshid => $_getSZ(4);
  @$pb.TagNumber(272889110)
  set lastrefreshid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(272889110)
  $core.bool hasLastrefreshid() => $_has(4);
  @$pb.TagNumber(272889110)
  void clearLastrefreshid() => $_clearField(272889110);

  @$pb.TagNumber(290836590)
  DashboardType get type => $_getN(5);
  @$pb.TagNumber(290836590)
  set type(DashboardType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(5);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(6);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(6, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(6);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(7);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(7);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(493889247)
  $core.String get lastrefreshfailurereason => $_getSZ(8);
  @$pb.TagNumber(493889247)
  set lastrefreshfailurereason($core.String value) => $_setString(8, value);
  @$pb.TagNumber(493889247)
  $core.bool hasLastrefreshfailurereason() => $_has(8);
  @$pb.TagNumber(493889247)
  void clearLastrefreshfailurereason() => $_clearField(493889247);

  @$pb.TagNumber(501826147)
  $pb.PbList<Widget> get widgets => $_getList(9);
}

class GetEventConfigurationRequest extends $pb.GeneratedMessage {
  factory GetEventConfigurationRequest({
    $core.String? eventdatastore,
    $core.String? trailname,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (trailname != null) result.trailname = trailname;
    return result;
  }

  GetEventConfigurationRequest._();

  factory GetEventConfigurationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEventConfigurationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEventConfigurationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(507774985, _omitFieldNames ? '' : 'trailname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventConfigurationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventConfigurationRequest copyWith(
          void Function(GetEventConfigurationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetEventConfigurationRequest))
          as GetEventConfigurationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEventConfigurationRequest create() =>
      GetEventConfigurationRequest._();
  @$core.override
  GetEventConfigurationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEventConfigurationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEventConfigurationRequest>(create);
  static GetEventConfigurationRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(507774985)
  $core.String get trailname => $_getSZ(1);
  @$pb.TagNumber(507774985)
  set trailname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(507774985)
  $core.bool hasTrailname() => $_has(1);
  @$pb.TagNumber(507774985)
  void clearTrailname() => $_clearField(507774985);
}

class GetEventConfigurationResponse extends $pb.GeneratedMessage {
  factory GetEventConfigurationResponse({
    $core.String? trailarn,
    $core.String? eventdatastorearn,
    $core.Iterable<ContextKeySelector>? contextkeyselectors,
    $core.Iterable<AggregationConfiguration>? aggregationconfigurations,
    MaxEventSize? maxeventsize,
  }) {
    final result = create();
    if (trailarn != null) result.trailarn = trailarn;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (contextkeyselectors != null)
      result.contextkeyselectors.addAll(contextkeyselectors);
    if (aggregationconfigurations != null)
      result.aggregationconfigurations.addAll(aggregationconfigurations);
    if (maxeventsize != null) result.maxeventsize = maxeventsize;
    return result;
  }

  GetEventConfigurationResponse._();

  factory GetEventConfigurationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEventConfigurationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEventConfigurationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..pPM<ContextKeySelector>(
        342040300, _omitFieldNames ? '' : 'contextkeyselectors',
        subBuilder: ContextKeySelector.create)
    ..pPM<AggregationConfiguration>(
        481530463, _omitFieldNames ? '' : 'aggregationconfigurations',
        subBuilder: AggregationConfiguration.create)
    ..aE<MaxEventSize>(517627763, _omitFieldNames ? '' : 'maxeventsize',
        enumValues: MaxEventSize.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventConfigurationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventConfigurationResponse copyWith(
          void Function(GetEventConfigurationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetEventConfigurationResponse))
          as GetEventConfigurationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEventConfigurationResponse create() =>
      GetEventConfigurationResponse._();
  @$core.override
  GetEventConfigurationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEventConfigurationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEventConfigurationResponse>(create);
  static GetEventConfigurationResponse? _defaultInstance;

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(0);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(0);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(1);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(1);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(342040300)
  $pb.PbList<ContextKeySelector> get contextkeyselectors => $_getList(2);

  @$pb.TagNumber(481530463)
  $pb.PbList<AggregationConfiguration> get aggregationconfigurations =>
      $_getList(3);

  @$pb.TagNumber(517627763)
  MaxEventSize get maxeventsize => $_getN(4);
  @$pb.TagNumber(517627763)
  set maxeventsize(MaxEventSize value) => $_setField(517627763, value);
  @$pb.TagNumber(517627763)
  $core.bool hasMaxeventsize() => $_has(4);
  @$pb.TagNumber(517627763)
  void clearMaxeventsize() => $_clearField(517627763);
}

class GetEventDataStoreRequest extends $pb.GeneratedMessage {
  factory GetEventDataStoreRequest({
    $core.String? eventdatastore,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    return result;
  }

  GetEventDataStoreRequest._();

  factory GetEventDataStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEventDataStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEventDataStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventDataStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventDataStoreRequest copyWith(
          void Function(GetEventDataStoreRequest) updates) =>
      super.copyWith((message) => updates(message as GetEventDataStoreRequest))
          as GetEventDataStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEventDataStoreRequest create() => GetEventDataStoreRequest._();
  @$core.override
  GetEventDataStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEventDataStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEventDataStoreRequest>(create);
  static GetEventDataStoreRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);
}

class GetEventDataStoreResponse extends $pb.GeneratedMessage {
  factory GetEventDataStoreResponse({
    EventDataStoreStatus? status,
    $core.bool? multiregionenabled,
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? updatedtimestamp,
    $core.String? kmskeyid,
    FederationStatus? federationstatus,
    BillingMode? billingmode,
    $core.int? retentionperiod,
    $core.Iterable<PartitionKey>? partitionkeys,
    $core.String? name,
    $core.String? eventdatastorearn,
    $core.String? createdtimestamp,
    $core.bool? terminationprotectionenabled,
    $core.bool? organizationenabled,
    $core.String? federationrolearn,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (multiregionenabled != null)
      result.multiregionenabled = multiregionenabled;
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (federationstatus != null) result.federationstatus = federationstatus;
    if (billingmode != null) result.billingmode = billingmode;
    if (retentionperiod != null) result.retentionperiod = retentionperiod;
    if (partitionkeys != null) result.partitionkeys.addAll(partitionkeys);
    if (name != null) result.name = name;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (organizationenabled != null)
      result.organizationenabled = organizationenabled;
    if (federationrolearn != null) result.federationrolearn = federationrolearn;
    return result;
  }

  GetEventDataStoreResponse._();

  factory GetEventDataStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEventDataStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEventDataStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<EventDataStoreStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: EventDataStoreStatus.values)
    ..aOB(20620620, _omitFieldNames ? '' : 'multiregionenabled')
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aE<FederationStatus>(146235383, _omitFieldNames ? '' : 'federationstatus',
        enumValues: FederationStatus.values)
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aI(196383721, _omitFieldNames ? '' : 'retentionperiod')
    ..pPM<PartitionKey>(200562986, _omitFieldNames ? '' : 'partitionkeys',
        subBuilder: PartitionKey.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOB(480171176, _omitFieldNames ? '' : 'organizationenabled')
    ..aOS(504464364, _omitFieldNames ? '' : 'federationrolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventDataStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventDataStoreResponse copyWith(
          void Function(GetEventDataStoreResponse) updates) =>
      super.copyWith((message) => updates(message as GetEventDataStoreResponse))
          as GetEventDataStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEventDataStoreResponse create() => GetEventDataStoreResponse._();
  @$core.override
  GetEventDataStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEventDataStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEventDataStoreResponse>(create);
  static GetEventDataStoreResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  EventDataStoreStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(EventDataStoreStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(20620620)
  $core.bool get multiregionenabled => $_getBF(1);
  @$pb.TagNumber(20620620)
  set multiregionenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(20620620)
  $core.bool hasMultiregionenabled() => $_has(1);
  @$pb.TagNumber(20620620)
  void clearMultiregionenabled() => $_clearField(20620620);

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(2);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(3);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(3);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(4);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(4);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(146235383)
  FederationStatus get federationstatus => $_getN(5);
  @$pb.TagNumber(146235383)
  set federationstatus(FederationStatus value) => $_setField(146235383, value);
  @$pb.TagNumber(146235383)
  $core.bool hasFederationstatus() => $_has(5);
  @$pb.TagNumber(146235383)
  void clearFederationstatus() => $_clearField(146235383);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(6);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(6);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(196383721)
  $core.int get retentionperiod => $_getIZ(7);
  @$pb.TagNumber(196383721)
  set retentionperiod($core.int value) => $_setSignedInt32(7, value);
  @$pb.TagNumber(196383721)
  $core.bool hasRetentionperiod() => $_has(7);
  @$pb.TagNumber(196383721)
  void clearRetentionperiod() => $_clearField(196383721);

  @$pb.TagNumber(200562986)
  $pb.PbList<PartitionKey> get partitionkeys => $_getList(8);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(9);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(9, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(9);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(10);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(10, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(10);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(11);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(11, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(11);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(12);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(12, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(12);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(480171176)
  $core.bool get organizationenabled => $_getBF(13);
  @$pb.TagNumber(480171176)
  set organizationenabled($core.bool value) => $_setBool(13, value);
  @$pb.TagNumber(480171176)
  $core.bool hasOrganizationenabled() => $_has(13);
  @$pb.TagNumber(480171176)
  void clearOrganizationenabled() => $_clearField(480171176);

  @$pb.TagNumber(504464364)
  $core.String get federationrolearn => $_getSZ(14);
  @$pb.TagNumber(504464364)
  set federationrolearn($core.String value) => $_setString(14, value);
  @$pb.TagNumber(504464364)
  $core.bool hasFederationrolearn() => $_has(14);
  @$pb.TagNumber(504464364)
  void clearFederationrolearn() => $_clearField(504464364);
}

class GetEventSelectorsRequest extends $pb.GeneratedMessage {
  factory GetEventSelectorsRequest({
    $core.String? trailname,
  }) {
    final result = create();
    if (trailname != null) result.trailname = trailname;
    return result;
  }

  GetEventSelectorsRequest._();

  factory GetEventSelectorsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEventSelectorsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEventSelectorsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(507774985, _omitFieldNames ? '' : 'trailname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventSelectorsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventSelectorsRequest copyWith(
          void Function(GetEventSelectorsRequest) updates) =>
      super.copyWith((message) => updates(message as GetEventSelectorsRequest))
          as GetEventSelectorsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEventSelectorsRequest create() => GetEventSelectorsRequest._();
  @$core.override
  GetEventSelectorsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEventSelectorsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEventSelectorsRequest>(create);
  static GetEventSelectorsRequest? _defaultInstance;

  @$pb.TagNumber(507774985)
  $core.String get trailname => $_getSZ(0);
  @$pb.TagNumber(507774985)
  set trailname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(507774985)
  $core.bool hasTrailname() => $_has(0);
  @$pb.TagNumber(507774985)
  void clearTrailname() => $_clearField(507774985);
}

class GetEventSelectorsResponse extends $pb.GeneratedMessage {
  factory GetEventSelectorsResponse({
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? trailarn,
    $core.Iterable<EventSelector>? eventselectors,
  }) {
    final result = create();
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (trailarn != null) result.trailarn = trailarn;
    if (eventselectors != null) result.eventselectors.addAll(eventselectors);
    return result;
  }

  GetEventSelectorsResponse._();

  factory GetEventSelectorsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEventSelectorsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEventSelectorsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..pPM<EventSelector>(344575328, _omitFieldNames ? '' : 'eventselectors',
        subBuilder: EventSelector.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventSelectorsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEventSelectorsResponse copyWith(
          void Function(GetEventSelectorsResponse) updates) =>
      super.copyWith((message) => updates(message as GetEventSelectorsResponse))
          as GetEventSelectorsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEventSelectorsResponse create() => GetEventSelectorsResponse._();
  @$core.override
  GetEventSelectorsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEventSelectorsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEventSelectorsResponse>(create);
  static GetEventSelectorsResponse? _defaultInstance;

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(0);

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(1);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(1);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(344575328)
  $pb.PbList<EventSelector> get eventselectors => $_getList(2);
}

class GetImportRequest extends $pb.GeneratedMessage {
  factory GetImportRequest({
    $core.String? importid,
  }) {
    final result = create();
    if (importid != null) result.importid = importid;
    return result;
  }

  GetImportRequest._();

  factory GetImportRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetImportRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetImportRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImportRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImportRequest copyWith(void Function(GetImportRequest) updates) =>
      super.copyWith((message) => updates(message as GetImportRequest))
          as GetImportRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetImportRequest create() => GetImportRequest._();
  @$core.override
  GetImportRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetImportRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetImportRequest>(create);
  static GetImportRequest? _defaultInstance;

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(0);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(0);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class GetImportResponse extends $pb.GeneratedMessage {
  factory GetImportResponse({
    ImportSource? importsource,
    $core.String? updatedtimestamp,
    ImportStatistics? importstatistics,
    $core.String? starteventtime,
    ImportStatus? importstatus,
    $core.String? endeventtime,
    $core.String? createdtimestamp,
    $core.Iterable<$core.String>? destinations,
    $core.String? importid,
  }) {
    final result = create();
    if (importsource != null) result.importsource = importsource;
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (importstatistics != null) result.importstatistics = importstatistics;
    if (starteventtime != null) result.starteventtime = starteventtime;
    if (importstatus != null) result.importstatus = importstatus;
    if (endeventtime != null) result.endeventtime = endeventtime;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (destinations != null) result.destinations.addAll(destinations);
    if (importid != null) result.importid = importid;
    return result;
  }

  GetImportResponse._();

  factory GetImportResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetImportResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetImportResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<ImportSource>(41128754, _omitFieldNames ? '' : 'importsource',
        subBuilder: ImportSource.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOM<ImportStatistics>(48175528, _omitFieldNames ? '' : 'importstatistics',
        subBuilder: ImportStatistics.create)
    ..aOS(107272573, _omitFieldNames ? '' : 'starteventtime')
    ..aE<ImportStatus>(129077631, _omitFieldNames ? '' : 'importstatus',
        enumValues: ImportStatus.values)
    ..aOS(260441984, _omitFieldNames ? '' : 'endeventtime')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..pPS(404385861, _omitFieldNames ? '' : 'destinations')
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImportResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetImportResponse copyWith(void Function(GetImportResponse) updates) =>
      super.copyWith((message) => updates(message as GetImportResponse))
          as GetImportResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetImportResponse create() => GetImportResponse._();
  @$core.override
  GetImportResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetImportResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetImportResponse>(create);
  static GetImportResponse? _defaultInstance;

  @$pb.TagNumber(41128754)
  ImportSource get importsource => $_getN(0);
  @$pb.TagNumber(41128754)
  set importsource(ImportSource value) => $_setField(41128754, value);
  @$pb.TagNumber(41128754)
  $core.bool hasImportsource() => $_has(0);
  @$pb.TagNumber(41128754)
  void clearImportsource() => $_clearField(41128754);
  @$pb.TagNumber(41128754)
  ImportSource ensureImportsource() => $_ensure(0);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(1);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(1);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(48175528)
  ImportStatistics get importstatistics => $_getN(2);
  @$pb.TagNumber(48175528)
  set importstatistics(ImportStatistics value) => $_setField(48175528, value);
  @$pb.TagNumber(48175528)
  $core.bool hasImportstatistics() => $_has(2);
  @$pb.TagNumber(48175528)
  void clearImportstatistics() => $_clearField(48175528);
  @$pb.TagNumber(48175528)
  ImportStatistics ensureImportstatistics() => $_ensure(2);

  @$pb.TagNumber(107272573)
  $core.String get starteventtime => $_getSZ(3);
  @$pb.TagNumber(107272573)
  set starteventtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(107272573)
  $core.bool hasStarteventtime() => $_has(3);
  @$pb.TagNumber(107272573)
  void clearStarteventtime() => $_clearField(107272573);

  @$pb.TagNumber(129077631)
  ImportStatus get importstatus => $_getN(4);
  @$pb.TagNumber(129077631)
  set importstatus(ImportStatus value) => $_setField(129077631, value);
  @$pb.TagNumber(129077631)
  $core.bool hasImportstatus() => $_has(4);
  @$pb.TagNumber(129077631)
  void clearImportstatus() => $_clearField(129077631);

  @$pb.TagNumber(260441984)
  $core.String get endeventtime => $_getSZ(5);
  @$pb.TagNumber(260441984)
  set endeventtime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(260441984)
  $core.bool hasEndeventtime() => $_has(5);
  @$pb.TagNumber(260441984)
  void clearEndeventtime() => $_clearField(260441984);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(6);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(6, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(6);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(404385861)
  $pb.PbList<$core.String> get destinations => $_getList(7);

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(8);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(8);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class GetInsightSelectorsRequest extends $pb.GeneratedMessage {
  factory GetInsightSelectorsRequest({
    $core.String? eventdatastore,
    $core.String? trailname,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (trailname != null) result.trailname = trailname;
    return result;
  }

  GetInsightSelectorsRequest._();

  factory GetInsightSelectorsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetInsightSelectorsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetInsightSelectorsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(507774985, _omitFieldNames ? '' : 'trailname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightSelectorsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightSelectorsRequest copyWith(
          void Function(GetInsightSelectorsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetInsightSelectorsRequest))
          as GetInsightSelectorsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetInsightSelectorsRequest create() => GetInsightSelectorsRequest._();
  @$core.override
  GetInsightSelectorsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetInsightSelectorsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetInsightSelectorsRequest>(create);
  static GetInsightSelectorsRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(507774985)
  $core.String get trailname => $_getSZ(1);
  @$pb.TagNumber(507774985)
  set trailname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(507774985)
  $core.bool hasTrailname() => $_has(1);
  @$pb.TagNumber(507774985)
  void clearTrailname() => $_clearField(507774985);
}

class GetInsightSelectorsResponse extends $pb.GeneratedMessage {
  factory GetInsightSelectorsResponse({
    $core.String? trailarn,
    $core.Iterable<InsightSelector>? insightselectors,
    $core.String? eventdatastorearn,
    $core.String? insightsdestination,
  }) {
    final result = create();
    if (trailarn != null) result.trailarn = trailarn;
    if (insightselectors != null)
      result.insightselectors.addAll(insightselectors);
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (insightsdestination != null)
      result.insightsdestination = insightsdestination;
    return result;
  }

  GetInsightSelectorsResponse._();

  factory GetInsightSelectorsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetInsightSelectorsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetInsightSelectorsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..pPM<InsightSelector>(49628748, _omitFieldNames ? '' : 'insightselectors',
        subBuilder: InsightSelector.create)
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(396776649, _omitFieldNames ? '' : 'insightsdestination')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightSelectorsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightSelectorsResponse copyWith(
          void Function(GetInsightSelectorsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetInsightSelectorsResponse))
          as GetInsightSelectorsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetInsightSelectorsResponse create() =>
      GetInsightSelectorsResponse._();
  @$core.override
  GetInsightSelectorsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetInsightSelectorsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetInsightSelectorsResponse>(create);
  static GetInsightSelectorsResponse? _defaultInstance;

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(0);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(0);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(49628748)
  $pb.PbList<InsightSelector> get insightselectors => $_getList(1);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(2);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(2);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(396776649)
  $core.String get insightsdestination => $_getSZ(3);
  @$pb.TagNumber(396776649)
  set insightsdestination($core.String value) => $_setString(3, value);
  @$pb.TagNumber(396776649)
  $core.bool hasInsightsdestination() => $_has(3);
  @$pb.TagNumber(396776649)
  void clearInsightsdestination() => $_clearField(396776649);
}

class GetQueryResultsRequest extends $pb.GeneratedMessage {
  factory GetQueryResultsRequest({
    $core.String? queryid,
    $core.int? maxqueryresults,
    $core.String? eventdatastore,
    $core.String? nexttoken,
    $core.String? eventdatastoreowneraccountid,
  }) {
    final result = create();
    if (queryid != null) result.queryid = queryid;
    if (maxqueryresults != null) result.maxqueryresults = maxqueryresults;
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    return result;
  }

  GetQueryResultsRequest._();

  factory GetQueryResultsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryResultsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryResultsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..aI(120660448, _omitFieldNames ? '' : 'maxqueryresults')
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsRequest copyWith(
          void Function(GetQueryResultsRequest) updates) =>
      super.copyWith((message) => updates(message as GetQueryResultsRequest))
          as GetQueryResultsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryResultsRequest create() => GetQueryResultsRequest._();
  @$core.override
  GetQueryResultsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryResultsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryResultsRequest>(create);
  static GetQueryResultsRequest? _defaultInstance;

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(0);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(0);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(120660448)
  $core.int get maxqueryresults => $_getIZ(1);
  @$pb.TagNumber(120660448)
  set maxqueryresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(120660448)
  $core.bool hasMaxqueryresults() => $_has(1);
  @$pb.TagNumber(120660448)
  void clearMaxqueryresults() => $_clearField(120660448);

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(2);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(2, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(2);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(3);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(3);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(4);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(4);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);
}

class GetQueryResultsResponse extends $pb.GeneratedMessage {
  factory GetQueryResultsResponse({
    $core.String? nexttoken,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? queryresultrows,
    QueryStatistics? querystatistics,
    QueryStatus? querystatus,
    $core.String? errormessage,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (queryresultrows != null)
      result.queryresultrows.addEntries(queryresultrows);
    if (querystatistics != null) result.querystatistics = querystatistics;
    if (querystatus != null) result.querystatus = querystatus;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  GetQueryResultsResponse._();

  factory GetQueryResultsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryResultsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryResultsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..m<$core.String, $core.String>(
        240264704, _omitFieldNames ? '' : 'queryresultrows',
        entryClassName: 'GetQueryResultsResponse.QueryresultrowsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cloudtrail'))
    ..aOM<QueryStatistics>(260794841, _omitFieldNames ? '' : 'querystatistics',
        subBuilder: QueryStatistics.create)
    ..aE<QueryStatus>(367016406, _omitFieldNames ? '' : 'querystatus',
        enumValues: QueryStatus.values)
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsResponse copyWith(
          void Function(GetQueryResultsResponse) updates) =>
      super.copyWith((message) => updates(message as GetQueryResultsResponse))
          as GetQueryResultsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryResultsResponse create() => GetQueryResultsResponse._();
  @$core.override
  GetQueryResultsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryResultsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryResultsResponse>(create);
  static GetQueryResultsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(240264704)
  $pb.PbMap<$core.String, $core.String> get queryresultrows => $_getMap(1);

  @$pb.TagNumber(260794841)
  QueryStatistics get querystatistics => $_getN(2);
  @$pb.TagNumber(260794841)
  set querystatistics(QueryStatistics value) => $_setField(260794841, value);
  @$pb.TagNumber(260794841)
  $core.bool hasQuerystatistics() => $_has(2);
  @$pb.TagNumber(260794841)
  void clearQuerystatistics() => $_clearField(260794841);
  @$pb.TagNumber(260794841)
  QueryStatistics ensureQuerystatistics() => $_ensure(2);

  @$pb.TagNumber(367016406)
  QueryStatus get querystatus => $_getN(3);
  @$pb.TagNumber(367016406)
  set querystatus(QueryStatus value) => $_setField(367016406, value);
  @$pb.TagNumber(367016406)
  $core.bool hasQuerystatus() => $_has(3);
  @$pb.TagNumber(367016406)
  void clearQuerystatus() => $_clearField(367016406);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(4);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(4, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(4);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class GetResourcePolicyRequest extends $pb.GeneratedMessage {
  factory GetResourcePolicyRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class GetResourcePolicyResponse extends $pb.GeneratedMessage {
  factory GetResourcePolicyResponse({
    $core.String? resourcepolicy,
    $core.String? resourcearn,
    $core.String? delegatedadminresourcepolicy,
  }) {
    final result = create();
    if (resourcepolicy != null) result.resourcepolicy = resourcepolicy;
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (delegatedadminresourcepolicy != null)
      result.delegatedadminresourcepolicy = delegatedadminresourcepolicy;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(15747632, _omitFieldNames ? '' : 'resourcepolicy')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(510966132, _omitFieldNames ? '' : 'delegatedadminresourcepolicy')
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(510966132)
  $core.String get delegatedadminresourcepolicy => $_getSZ(2);
  @$pb.TagNumber(510966132)
  set delegatedadminresourcepolicy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(510966132)
  $core.bool hasDelegatedadminresourcepolicy() => $_has(2);
  @$pb.TagNumber(510966132)
  void clearDelegatedadminresourcepolicy() => $_clearField(510966132);
}

class GetTrailRequest extends $pb.GeneratedMessage {
  factory GetTrailRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  GetTrailRequest._();

  factory GetTrailRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrailRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrailRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailRequest copyWith(void Function(GetTrailRequest) updates) =>
      super.copyWith((message) => updates(message as GetTrailRequest))
          as GetTrailRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrailRequest create() => GetTrailRequest._();
  @$core.override
  GetTrailRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrailRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrailRequest>(create);
  static GetTrailRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class GetTrailResponse extends $pb.GeneratedMessage {
  factory GetTrailResponse({
    Trail? trail,
  }) {
    final result = create();
    if (trail != null) result.trail = trail;
    return result;
  }

  GetTrailResponse._();

  factory GetTrailResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrailResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrailResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<Trail>(470969828, _omitFieldNames ? '' : 'trail',
        subBuilder: Trail.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailResponse copyWith(void Function(GetTrailResponse) updates) =>
      super.copyWith((message) => updates(message as GetTrailResponse))
          as GetTrailResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrailResponse create() => GetTrailResponse._();
  @$core.override
  GetTrailResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrailResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrailResponse>(create);
  static GetTrailResponse? _defaultInstance;

  @$pb.TagNumber(470969828)
  Trail get trail => $_getN(0);
  @$pb.TagNumber(470969828)
  set trail(Trail value) => $_setField(470969828, value);
  @$pb.TagNumber(470969828)
  $core.bool hasTrail() => $_has(0);
  @$pb.TagNumber(470969828)
  void clearTrail() => $_clearField(470969828);
  @$pb.TagNumber(470969828)
  Trail ensureTrail() => $_ensure(0);
}

class GetTrailStatusRequest extends $pb.GeneratedMessage {
  factory GetTrailStatusRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  GetTrailStatusRequest._();

  factory GetTrailStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrailStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrailStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailStatusRequest copyWith(
          void Function(GetTrailStatusRequest) updates) =>
      super.copyWith((message) => updates(message as GetTrailStatusRequest))
          as GetTrailStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrailStatusRequest create() => GetTrailStatusRequest._();
  @$core.override
  GetTrailStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrailStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrailStatusRequest>(create);
  static GetTrailStatusRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class GetTrailStatusResponse extends $pb.GeneratedMessage {
  factory GetTrailStatusResponse({
    $core.String? latestnotificationtime,
    $core.String? latestnotificationattemptsucceeded,
    $core.String? timeloggingstarted,
    $core.bool? islogging,
    $core.String? latestnotificationattempttime,
    $core.String? latestdeliveryattempttime,
    $core.String? latestdeliveryerror,
    $core.String? latestcloudwatchlogsdeliveryerror,
    $core.String? latestnotificationerror,
    $core.String? timeloggingstopped,
    $core.String? latestdeliverytime,
    $core.String? startloggingtime,
    $core.String? latestdeliveryattemptsucceeded,
    $core.String? latestdigestdeliverytime,
    $core.String? stoploggingtime,
    $core.String? latestdigestdeliveryerror,
    $core.String? latestcloudwatchlogsdeliverytime,
  }) {
    final result = create();
    if (latestnotificationtime != null)
      result.latestnotificationtime = latestnotificationtime;
    if (latestnotificationattemptsucceeded != null)
      result.latestnotificationattemptsucceeded =
          latestnotificationattemptsucceeded;
    if (timeloggingstarted != null)
      result.timeloggingstarted = timeloggingstarted;
    if (islogging != null) result.islogging = islogging;
    if (latestnotificationattempttime != null)
      result.latestnotificationattempttime = latestnotificationattempttime;
    if (latestdeliveryattempttime != null)
      result.latestdeliveryattempttime = latestdeliveryattempttime;
    if (latestdeliveryerror != null)
      result.latestdeliveryerror = latestdeliveryerror;
    if (latestcloudwatchlogsdeliveryerror != null)
      result.latestcloudwatchlogsdeliveryerror =
          latestcloudwatchlogsdeliveryerror;
    if (latestnotificationerror != null)
      result.latestnotificationerror = latestnotificationerror;
    if (timeloggingstopped != null)
      result.timeloggingstopped = timeloggingstopped;
    if (latestdeliverytime != null)
      result.latestdeliverytime = latestdeliverytime;
    if (startloggingtime != null) result.startloggingtime = startloggingtime;
    if (latestdeliveryattemptsucceeded != null)
      result.latestdeliveryattemptsucceeded = latestdeliveryattemptsucceeded;
    if (latestdigestdeliverytime != null)
      result.latestdigestdeliverytime = latestdigestdeliverytime;
    if (stoploggingtime != null) result.stoploggingtime = stoploggingtime;
    if (latestdigestdeliveryerror != null)
      result.latestdigestdeliveryerror = latestdigestdeliveryerror;
    if (latestcloudwatchlogsdeliverytime != null)
      result.latestcloudwatchlogsdeliverytime =
          latestcloudwatchlogsdeliverytime;
    return result;
  }

  GetTrailStatusResponse._();

  factory GetTrailStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrailStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrailStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(2082317, _omitFieldNames ? '' : 'latestnotificationtime')
    ..aOS(3764412, _omitFieldNames ? '' : 'latestnotificationattemptsucceeded')
    ..aOS(18088239, _omitFieldNames ? '' : 'timeloggingstarted')
    ..aOB(61005519, _omitFieldNames ? '' : 'islogging')
    ..aOS(104168880, _omitFieldNames ? '' : 'latestnotificationattempttime')
    ..aOS(109095435, _omitFieldNames ? '' : 'latestdeliveryattempttime')
    ..aOS(119953649, _omitFieldNames ? '' : 'latestdeliveryerror')
    ..aOS(238605214, _omitFieldNames ? '' : 'latestcloudwatchlogsdeliveryerror')
    ..aOS(261077338, _omitFieldNames ? '' : 'latestnotificationerror')
    ..aOS(272406495, _omitFieldNames ? '' : 'timeloggingstopped')
    ..aOS(320897248, _omitFieldNames ? '' : 'latestdeliverytime')
    ..aOS(351696170, _omitFieldNames ? '' : 'startloggingtime')
    ..aOS(432951833, _omitFieldNames ? '' : 'latestdeliveryattemptsucceeded')
    ..aOS(469527024, _omitFieldNames ? '' : 'latestdigestdeliverytime')
    ..aOS(469874232, _omitFieldNames ? '' : 'stoploggingtime')
    ..aOS(514209249, _omitFieldNames ? '' : 'latestdigestdeliveryerror')
    ..aOS(517138001, _omitFieldNames ? '' : 'latestcloudwatchlogsdeliverytime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrailStatusResponse copyWith(
          void Function(GetTrailStatusResponse) updates) =>
      super.copyWith((message) => updates(message as GetTrailStatusResponse))
          as GetTrailStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrailStatusResponse create() => GetTrailStatusResponse._();
  @$core.override
  GetTrailStatusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrailStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrailStatusResponse>(create);
  static GetTrailStatusResponse? _defaultInstance;

  @$pb.TagNumber(2082317)
  $core.String get latestnotificationtime => $_getSZ(0);
  @$pb.TagNumber(2082317)
  set latestnotificationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(2082317)
  $core.bool hasLatestnotificationtime() => $_has(0);
  @$pb.TagNumber(2082317)
  void clearLatestnotificationtime() => $_clearField(2082317);

  @$pb.TagNumber(3764412)
  $core.String get latestnotificationattemptsucceeded => $_getSZ(1);
  @$pb.TagNumber(3764412)
  set latestnotificationattemptsucceeded($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(3764412)
  $core.bool hasLatestnotificationattemptsucceeded() => $_has(1);
  @$pb.TagNumber(3764412)
  void clearLatestnotificationattemptsucceeded() => $_clearField(3764412);

  @$pb.TagNumber(18088239)
  $core.String get timeloggingstarted => $_getSZ(2);
  @$pb.TagNumber(18088239)
  set timeloggingstarted($core.String value) => $_setString(2, value);
  @$pb.TagNumber(18088239)
  $core.bool hasTimeloggingstarted() => $_has(2);
  @$pb.TagNumber(18088239)
  void clearTimeloggingstarted() => $_clearField(18088239);

  @$pb.TagNumber(61005519)
  $core.bool get islogging => $_getBF(3);
  @$pb.TagNumber(61005519)
  set islogging($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(61005519)
  $core.bool hasIslogging() => $_has(3);
  @$pb.TagNumber(61005519)
  void clearIslogging() => $_clearField(61005519);

  @$pb.TagNumber(104168880)
  $core.String get latestnotificationattempttime => $_getSZ(4);
  @$pb.TagNumber(104168880)
  set latestnotificationattempttime($core.String value) =>
      $_setString(4, value);
  @$pb.TagNumber(104168880)
  $core.bool hasLatestnotificationattempttime() => $_has(4);
  @$pb.TagNumber(104168880)
  void clearLatestnotificationattempttime() => $_clearField(104168880);

  @$pb.TagNumber(109095435)
  $core.String get latestdeliveryattempttime => $_getSZ(5);
  @$pb.TagNumber(109095435)
  set latestdeliveryattempttime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(109095435)
  $core.bool hasLatestdeliveryattempttime() => $_has(5);
  @$pb.TagNumber(109095435)
  void clearLatestdeliveryattempttime() => $_clearField(109095435);

  @$pb.TagNumber(119953649)
  $core.String get latestdeliveryerror => $_getSZ(6);
  @$pb.TagNumber(119953649)
  set latestdeliveryerror($core.String value) => $_setString(6, value);
  @$pb.TagNumber(119953649)
  $core.bool hasLatestdeliveryerror() => $_has(6);
  @$pb.TagNumber(119953649)
  void clearLatestdeliveryerror() => $_clearField(119953649);

  @$pb.TagNumber(238605214)
  $core.String get latestcloudwatchlogsdeliveryerror => $_getSZ(7);
  @$pb.TagNumber(238605214)
  set latestcloudwatchlogsdeliveryerror($core.String value) =>
      $_setString(7, value);
  @$pb.TagNumber(238605214)
  $core.bool hasLatestcloudwatchlogsdeliveryerror() => $_has(7);
  @$pb.TagNumber(238605214)
  void clearLatestcloudwatchlogsdeliveryerror() => $_clearField(238605214);

  @$pb.TagNumber(261077338)
  $core.String get latestnotificationerror => $_getSZ(8);
  @$pb.TagNumber(261077338)
  set latestnotificationerror($core.String value) => $_setString(8, value);
  @$pb.TagNumber(261077338)
  $core.bool hasLatestnotificationerror() => $_has(8);
  @$pb.TagNumber(261077338)
  void clearLatestnotificationerror() => $_clearField(261077338);

  @$pb.TagNumber(272406495)
  $core.String get timeloggingstopped => $_getSZ(9);
  @$pb.TagNumber(272406495)
  set timeloggingstopped($core.String value) => $_setString(9, value);
  @$pb.TagNumber(272406495)
  $core.bool hasTimeloggingstopped() => $_has(9);
  @$pb.TagNumber(272406495)
  void clearTimeloggingstopped() => $_clearField(272406495);

  @$pb.TagNumber(320897248)
  $core.String get latestdeliverytime => $_getSZ(10);
  @$pb.TagNumber(320897248)
  set latestdeliverytime($core.String value) => $_setString(10, value);
  @$pb.TagNumber(320897248)
  $core.bool hasLatestdeliverytime() => $_has(10);
  @$pb.TagNumber(320897248)
  void clearLatestdeliverytime() => $_clearField(320897248);

  @$pb.TagNumber(351696170)
  $core.String get startloggingtime => $_getSZ(11);
  @$pb.TagNumber(351696170)
  set startloggingtime($core.String value) => $_setString(11, value);
  @$pb.TagNumber(351696170)
  $core.bool hasStartloggingtime() => $_has(11);
  @$pb.TagNumber(351696170)
  void clearStartloggingtime() => $_clearField(351696170);

  @$pb.TagNumber(432951833)
  $core.String get latestdeliveryattemptsucceeded => $_getSZ(12);
  @$pb.TagNumber(432951833)
  set latestdeliveryattemptsucceeded($core.String value) =>
      $_setString(12, value);
  @$pb.TagNumber(432951833)
  $core.bool hasLatestdeliveryattemptsucceeded() => $_has(12);
  @$pb.TagNumber(432951833)
  void clearLatestdeliveryattemptsucceeded() => $_clearField(432951833);

  @$pb.TagNumber(469527024)
  $core.String get latestdigestdeliverytime => $_getSZ(13);
  @$pb.TagNumber(469527024)
  set latestdigestdeliverytime($core.String value) => $_setString(13, value);
  @$pb.TagNumber(469527024)
  $core.bool hasLatestdigestdeliverytime() => $_has(13);
  @$pb.TagNumber(469527024)
  void clearLatestdigestdeliverytime() => $_clearField(469527024);

  @$pb.TagNumber(469874232)
  $core.String get stoploggingtime => $_getSZ(14);
  @$pb.TagNumber(469874232)
  set stoploggingtime($core.String value) => $_setString(14, value);
  @$pb.TagNumber(469874232)
  $core.bool hasStoploggingtime() => $_has(14);
  @$pb.TagNumber(469874232)
  void clearStoploggingtime() => $_clearField(469874232);

  @$pb.TagNumber(514209249)
  $core.String get latestdigestdeliveryerror => $_getSZ(15);
  @$pb.TagNumber(514209249)
  set latestdigestdeliveryerror($core.String value) => $_setString(15, value);
  @$pb.TagNumber(514209249)
  $core.bool hasLatestdigestdeliveryerror() => $_has(15);
  @$pb.TagNumber(514209249)
  void clearLatestdigestdeliveryerror() => $_clearField(514209249);

  @$pb.TagNumber(517138001)
  $core.String get latestcloudwatchlogsdeliverytime => $_getSZ(16);
  @$pb.TagNumber(517138001)
  set latestcloudwatchlogsdeliverytime($core.String value) =>
      $_setString(16, value);
  @$pb.TagNumber(517138001)
  $core.bool hasLatestcloudwatchlogsdeliverytime() => $_has(16);
  @$pb.TagNumber(517138001)
  void clearLatestcloudwatchlogsdeliverytime() => $_clearField(517138001);
}

class ImportFailureListItem extends $pb.GeneratedMessage {
  factory ImportFailureListItem({
    ImportFailureStatus? status,
    $core.String? lastupdatedtime,
    $core.String? errortype,
    $core.String? location,
    $core.String? errormessage,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (lastupdatedtime != null) result.lastupdatedtime = lastupdatedtime;
    if (errortype != null) result.errortype = errortype;
    if (location != null) result.location = location;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  ImportFailureListItem._();

  factory ImportFailureListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportFailureListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportFailureListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<ImportFailureStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: ImportFailureStatus.values)
    ..aOS(177046166, _omitFieldNames ? '' : 'lastupdatedtime')
    ..aOS(398848954, _omitFieldNames ? '' : 'errortype')
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportFailureListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportFailureListItem copyWith(
          void Function(ImportFailureListItem) updates) =>
      super.copyWith((message) => updates(message as ImportFailureListItem))
          as ImportFailureListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportFailureListItem create() => ImportFailureListItem._();
  @$core.override
  ImportFailureListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportFailureListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportFailureListItem>(create);
  static ImportFailureListItem? _defaultInstance;

  @$pb.TagNumber(6222352)
  ImportFailureStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(ImportFailureStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(177046166)
  $core.String get lastupdatedtime => $_getSZ(1);
  @$pb.TagNumber(177046166)
  set lastupdatedtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(177046166)
  $core.bool hasLastupdatedtime() => $_has(1);
  @$pb.TagNumber(177046166)
  void clearLastupdatedtime() => $_clearField(177046166);

  @$pb.TagNumber(398848954)
  $core.String get errortype => $_getSZ(2);
  @$pb.TagNumber(398848954)
  set errortype($core.String value) => $_setString(2, value);
  @$pb.TagNumber(398848954)
  $core.bool hasErrortype() => $_has(2);
  @$pb.TagNumber(398848954)
  void clearErrortype() => $_clearField(398848954);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(3);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(3, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(3);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(4);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(4, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(4);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ImportSource extends $pb.GeneratedMessage {
  factory ImportSource({
    S3ImportSource? s3,
  }) {
    final result = create();
    if (s3 != null) result.s3 = s3;
    return result;
  }

  ImportSource._();

  factory ImportSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<S3ImportSource>(100795332, _omitFieldNames ? '' : 's3',
        subBuilder: S3ImportSource.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportSource copyWith(void Function(ImportSource) updates) =>
      super.copyWith((message) => updates(message as ImportSource))
          as ImportSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportSource create() => ImportSource._();
  @$core.override
  ImportSource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportSource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportSource>(create);
  static ImportSource? _defaultInstance;

  @$pb.TagNumber(100795332)
  S3ImportSource get s3 => $_getN(0);
  @$pb.TagNumber(100795332)
  set s3(S3ImportSource value) => $_setField(100795332, value);
  @$pb.TagNumber(100795332)
  $core.bool hasS3() => $_has(0);
  @$pb.TagNumber(100795332)
  void clearS3() => $_clearField(100795332);
  @$pb.TagNumber(100795332)
  S3ImportSource ensureS3() => $_ensure(0);
}

class ImportStatistics extends $pb.GeneratedMessage {
  factory ImportStatistics({
    $fixnum.Int64? failedentries,
    $fixnum.Int64? eventscompleted,
    $fixnum.Int64? prefixesfound,
    $fixnum.Int64? filescompleted,
    $fixnum.Int64? prefixescompleted,
  }) {
    final result = create();
    if (failedentries != null) result.failedentries = failedentries;
    if (eventscompleted != null) result.eventscompleted = eventscompleted;
    if (prefixesfound != null) result.prefixesfound = prefixesfound;
    if (filescompleted != null) result.filescompleted = filescompleted;
    if (prefixescompleted != null) result.prefixescompleted = prefixescompleted;
    return result;
  }

  ImportStatistics._();

  factory ImportStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aInt64(86416685, _omitFieldNames ? '' : 'failedentries')
    ..aInt64(159166494, _omitFieldNames ? '' : 'eventscompleted')
    ..aInt64(244228528, _omitFieldNames ? '' : 'prefixesfound')
    ..aInt64(300203834, _omitFieldNames ? '' : 'filescompleted')
    ..aInt64(433011203, _omitFieldNames ? '' : 'prefixescompleted')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportStatistics copyWith(void Function(ImportStatistics) updates) =>
      super.copyWith((message) => updates(message as ImportStatistics))
          as ImportStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportStatistics create() => ImportStatistics._();
  @$core.override
  ImportStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportStatistics>(create);
  static ImportStatistics? _defaultInstance;

  @$pb.TagNumber(86416685)
  $fixnum.Int64 get failedentries => $_getI64(0);
  @$pb.TagNumber(86416685)
  set failedentries($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(86416685)
  $core.bool hasFailedentries() => $_has(0);
  @$pb.TagNumber(86416685)
  void clearFailedentries() => $_clearField(86416685);

  @$pb.TagNumber(159166494)
  $fixnum.Int64 get eventscompleted => $_getI64(1);
  @$pb.TagNumber(159166494)
  set eventscompleted($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(159166494)
  $core.bool hasEventscompleted() => $_has(1);
  @$pb.TagNumber(159166494)
  void clearEventscompleted() => $_clearField(159166494);

  @$pb.TagNumber(244228528)
  $fixnum.Int64 get prefixesfound => $_getI64(2);
  @$pb.TagNumber(244228528)
  set prefixesfound($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(244228528)
  $core.bool hasPrefixesfound() => $_has(2);
  @$pb.TagNumber(244228528)
  void clearPrefixesfound() => $_clearField(244228528);

  @$pb.TagNumber(300203834)
  $fixnum.Int64 get filescompleted => $_getI64(3);
  @$pb.TagNumber(300203834)
  set filescompleted($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(300203834)
  $core.bool hasFilescompleted() => $_has(3);
  @$pb.TagNumber(300203834)
  void clearFilescompleted() => $_clearField(300203834);

  @$pb.TagNumber(433011203)
  $fixnum.Int64 get prefixescompleted => $_getI64(4);
  @$pb.TagNumber(433011203)
  set prefixescompleted($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(433011203)
  $core.bool hasPrefixescompleted() => $_has(4);
  @$pb.TagNumber(433011203)
  void clearPrefixescompleted() => $_clearField(433011203);
}

class ImportsListItem extends $pb.GeneratedMessage {
  factory ImportsListItem({
    $core.String? updatedtimestamp,
    ImportStatus? importstatus,
    $core.String? createdtimestamp,
    $core.Iterable<$core.String>? destinations,
    $core.String? importid,
  }) {
    final result = create();
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (importstatus != null) result.importstatus = importstatus;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (destinations != null) result.destinations.addAll(destinations);
    if (importid != null) result.importid = importid;
    return result;
  }

  ImportsListItem._();

  factory ImportsListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportsListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportsListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aE<ImportStatus>(129077631, _omitFieldNames ? '' : 'importstatus',
        enumValues: ImportStatus.values)
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..pPS(404385861, _omitFieldNames ? '' : 'destinations')
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportsListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportsListItem copyWith(void Function(ImportsListItem) updates) =>
      super.copyWith((message) => updates(message as ImportsListItem))
          as ImportsListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportsListItem create() => ImportsListItem._();
  @$core.override
  ImportsListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportsListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportsListItem>(create);
  static ImportsListItem? _defaultInstance;

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(0);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(0, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(0);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(129077631)
  ImportStatus get importstatus => $_getN(1);
  @$pb.TagNumber(129077631)
  set importstatus(ImportStatus value) => $_setField(129077631, value);
  @$pb.TagNumber(129077631)
  $core.bool hasImportstatus() => $_has(1);
  @$pb.TagNumber(129077631)
  void clearImportstatus() => $_clearField(129077631);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(2);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(2);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(404385861)
  $pb.PbList<$core.String> get destinations => $_getList(3);

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(4);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(4);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class InactiveEventDataStoreException extends $pb.GeneratedMessage {
  factory InactiveEventDataStoreException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InactiveEventDataStoreException._();

  factory InactiveEventDataStoreException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InactiveEventDataStoreException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InactiveEventDataStoreException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InactiveEventDataStoreException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InactiveEventDataStoreException copyWith(
          void Function(InactiveEventDataStoreException) updates) =>
      super.copyWith(
              (message) => updates(message as InactiveEventDataStoreException))
          as InactiveEventDataStoreException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InactiveEventDataStoreException create() =>
      InactiveEventDataStoreException._();
  @$core.override
  InactiveEventDataStoreException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InactiveEventDataStoreException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InactiveEventDataStoreException>(
          create);
  static InactiveEventDataStoreException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InactiveQueryException extends $pb.GeneratedMessage {
  factory InactiveQueryException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InactiveQueryException._();

  factory InactiveQueryException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InactiveQueryException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InactiveQueryException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InactiveQueryException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InactiveQueryException copyWith(
          void Function(InactiveQueryException) updates) =>
      super.copyWith((message) => updates(message as InactiveQueryException))
          as InactiveQueryException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InactiveQueryException create() => InactiveQueryException._();
  @$core.override
  InactiveQueryException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InactiveQueryException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InactiveQueryException>(create);
  static InactiveQueryException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class IngestionStatus extends $pb.GeneratedMessage {
  factory IngestionStatus({
    $core.String? latestingestionsuccesseventid,
    $core.String? latestingestionsuccesstime,
    $core.String? latestingestionattempteventid,
    $core.String? latestingestionattempttime,
    $core.String? latestingestionerrorcode,
  }) {
    final result = create();
    if (latestingestionsuccesseventid != null)
      result.latestingestionsuccesseventid = latestingestionsuccesseventid;
    if (latestingestionsuccesstime != null)
      result.latestingestionsuccesstime = latestingestionsuccesstime;
    if (latestingestionattempteventid != null)
      result.latestingestionattempteventid = latestingestionattempteventid;
    if (latestingestionattempttime != null)
      result.latestingestionattempttime = latestingestionattempttime;
    if (latestingestionerrorcode != null)
      result.latestingestionerrorcode = latestingestionerrorcode;
    return result;
  }

  IngestionStatus._();

  factory IngestionStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IngestionStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IngestionStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(31034217, _omitFieldNames ? '' : 'latestingestionsuccesseventid')
    ..aOS(184241615, _omitFieldNames ? '' : 'latestingestionsuccesstime')
    ..aOS(204299739, _omitFieldNames ? '' : 'latestingestionattempteventid')
    ..aOS(313195389, _omitFieldNames ? '' : 'latestingestionattempttime')
    ..aOS(415423856, _omitFieldNames ? '' : 'latestingestionerrorcode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IngestionStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IngestionStatus copyWith(void Function(IngestionStatus) updates) =>
      super.copyWith((message) => updates(message as IngestionStatus))
          as IngestionStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IngestionStatus create() => IngestionStatus._();
  @$core.override
  IngestionStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IngestionStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IngestionStatus>(create);
  static IngestionStatus? _defaultInstance;

  @$pb.TagNumber(31034217)
  $core.String get latestingestionsuccesseventid => $_getSZ(0);
  @$pb.TagNumber(31034217)
  set latestingestionsuccesseventid($core.String value) =>
      $_setString(0, value);
  @$pb.TagNumber(31034217)
  $core.bool hasLatestingestionsuccesseventid() => $_has(0);
  @$pb.TagNumber(31034217)
  void clearLatestingestionsuccesseventid() => $_clearField(31034217);

  @$pb.TagNumber(184241615)
  $core.String get latestingestionsuccesstime => $_getSZ(1);
  @$pb.TagNumber(184241615)
  set latestingestionsuccesstime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(184241615)
  $core.bool hasLatestingestionsuccesstime() => $_has(1);
  @$pb.TagNumber(184241615)
  void clearLatestingestionsuccesstime() => $_clearField(184241615);

  @$pb.TagNumber(204299739)
  $core.String get latestingestionattempteventid => $_getSZ(2);
  @$pb.TagNumber(204299739)
  set latestingestionattempteventid($core.String value) =>
      $_setString(2, value);
  @$pb.TagNumber(204299739)
  $core.bool hasLatestingestionattempteventid() => $_has(2);
  @$pb.TagNumber(204299739)
  void clearLatestingestionattempteventid() => $_clearField(204299739);

  @$pb.TagNumber(313195389)
  $core.String get latestingestionattempttime => $_getSZ(3);
  @$pb.TagNumber(313195389)
  set latestingestionattempttime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(313195389)
  $core.bool hasLatestingestionattempttime() => $_has(3);
  @$pb.TagNumber(313195389)
  void clearLatestingestionattempttime() => $_clearField(313195389);

  @$pb.TagNumber(415423856)
  $core.String get latestingestionerrorcode => $_getSZ(4);
  @$pb.TagNumber(415423856)
  set latestingestionerrorcode($core.String value) => $_setString(4, value);
  @$pb.TagNumber(415423856)
  $core.bool hasLatestingestionerrorcode() => $_has(4);
  @$pb.TagNumber(415423856)
  void clearLatestingestionerrorcode() => $_clearField(415423856);
}

class InsightNotEnabledException extends $pb.GeneratedMessage {
  factory InsightNotEnabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InsightNotEnabledException._();

  factory InsightNotEnabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsightNotEnabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsightNotEnabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightNotEnabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightNotEnabledException copyWith(
          void Function(InsightNotEnabledException) updates) =>
      super.copyWith(
              (message) => updates(message as InsightNotEnabledException))
          as InsightNotEnabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsightNotEnabledException create() => InsightNotEnabledException._();
  @$core.override
  InsightNotEnabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsightNotEnabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InsightNotEnabledException>(create);
  static InsightNotEnabledException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InsightSelector extends $pb.GeneratedMessage {
  factory InsightSelector({
    $core.Iterable<SourceEventCategory>? eventcategories,
    InsightType? insighttype,
  }) {
    final result = create();
    if (eventcategories != null) result.eventcategories.addAll(eventcategories);
    if (insighttype != null) result.insighttype = insighttype;
    return result;
  }

  InsightSelector._();

  factory InsightSelector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsightSelector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsightSelector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pc<SourceEventCategory>(
        3676820, _omitFieldNames ? '' : 'eventcategories', $pb.PbFieldType.KE,
        valueOf: SourceEventCategory.valueOf,
        enumValues: SourceEventCategory.values,
        defaultEnumValue: SourceEventCategory.SOURCE_EVENT_CATEGORY_DATA)
    ..aE<InsightType>(530375860, _omitFieldNames ? '' : 'insighttype',
        enumValues: InsightType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightSelector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightSelector copyWith(void Function(InsightSelector) updates) =>
      super.copyWith((message) => updates(message as InsightSelector))
          as InsightSelector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsightSelector create() => InsightSelector._();
  @$core.override
  InsightSelector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsightSelector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InsightSelector>(create);
  static InsightSelector? _defaultInstance;

  @$pb.TagNumber(3676820)
  $pb.PbList<SourceEventCategory> get eventcategories => $_getList(0);

  @$pb.TagNumber(530375860)
  InsightType get insighttype => $_getN(1);
  @$pb.TagNumber(530375860)
  set insighttype(InsightType value) => $_setField(530375860, value);
  @$pb.TagNumber(530375860)
  $core.bool hasInsighttype() => $_has(1);
  @$pb.TagNumber(530375860)
  void clearInsighttype() => $_clearField(530375860);
}

class InsufficientDependencyServiceAccessPermissionException
    extends $pb.GeneratedMessage {
  factory InsufficientDependencyServiceAccessPermissionException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InsufficientDependencyServiceAccessPermissionException._();

  factory InsufficientDependencyServiceAccessPermissionException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsufficientDependencyServiceAccessPermissionException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames
          ? ''
          : 'InsufficientDependencyServiceAccessPermissionException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientDependencyServiceAccessPermissionException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientDependencyServiceAccessPermissionException copyWith(
          void Function(InsufficientDependencyServiceAccessPermissionException)
              updates) =>
      super.copyWith((message) => updates(message
              as InsufficientDependencyServiceAccessPermissionException))
          as InsufficientDependencyServiceAccessPermissionException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsufficientDependencyServiceAccessPermissionException create() =>
      InsufficientDependencyServiceAccessPermissionException._();
  @$core.override
  InsufficientDependencyServiceAccessPermissionException
      createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsufficientDependencyServiceAccessPermissionException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InsufficientDependencyServiceAccessPermissionException>(create);
  static InsufficientDependencyServiceAccessPermissionException?
      _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InsufficientEncryptionPolicyException extends $pb.GeneratedMessage {
  factory InsufficientEncryptionPolicyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InsufficientEncryptionPolicyException._();

  factory InsufficientEncryptionPolicyException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsufficientEncryptionPolicyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsufficientEncryptionPolicyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientEncryptionPolicyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientEncryptionPolicyException copyWith(
          void Function(InsufficientEncryptionPolicyException) updates) =>
      super.copyWith((message) =>
              updates(message as InsufficientEncryptionPolicyException))
          as InsufficientEncryptionPolicyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsufficientEncryptionPolicyException create() =>
      InsufficientEncryptionPolicyException._();
  @$core.override
  InsufficientEncryptionPolicyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsufficientEncryptionPolicyException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InsufficientEncryptionPolicyException>(create);
  static InsufficientEncryptionPolicyException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InsufficientIAMAccessPermissionException extends $pb.GeneratedMessage {
  factory InsufficientIAMAccessPermissionException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InsufficientIAMAccessPermissionException._();

  factory InsufficientIAMAccessPermissionException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsufficientIAMAccessPermissionException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsufficientIAMAccessPermissionException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientIAMAccessPermissionException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientIAMAccessPermissionException copyWith(
          void Function(InsufficientIAMAccessPermissionException) updates) =>
      super.copyWith((message) =>
              updates(message as InsufficientIAMAccessPermissionException))
          as InsufficientIAMAccessPermissionException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsufficientIAMAccessPermissionException create() =>
      InsufficientIAMAccessPermissionException._();
  @$core.override
  InsufficientIAMAccessPermissionException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsufficientIAMAccessPermissionException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InsufficientIAMAccessPermissionException>(create);
  static InsufficientIAMAccessPermissionException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InsufficientS3BucketPolicyException extends $pb.GeneratedMessage {
  factory InsufficientS3BucketPolicyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InsufficientS3BucketPolicyException._();

  factory InsufficientS3BucketPolicyException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsufficientS3BucketPolicyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsufficientS3BucketPolicyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientS3BucketPolicyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientS3BucketPolicyException copyWith(
          void Function(InsufficientS3BucketPolicyException) updates) =>
      super.copyWith((message) =>
              updates(message as InsufficientS3BucketPolicyException))
          as InsufficientS3BucketPolicyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsufficientS3BucketPolicyException create() =>
      InsufficientS3BucketPolicyException._();
  @$core.override
  InsufficientS3BucketPolicyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsufficientS3BucketPolicyException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InsufficientS3BucketPolicyException>(create);
  static InsufficientS3BucketPolicyException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InsufficientSnsTopicPolicyException extends $pb.GeneratedMessage {
  factory InsufficientSnsTopicPolicyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InsufficientSnsTopicPolicyException._();

  factory InsufficientSnsTopicPolicyException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsufficientSnsTopicPolicyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsufficientSnsTopicPolicyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientSnsTopicPolicyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientSnsTopicPolicyException copyWith(
          void Function(InsufficientSnsTopicPolicyException) updates) =>
      super.copyWith((message) =>
              updates(message as InsufficientSnsTopicPolicyException))
          as InsufficientSnsTopicPolicyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsufficientSnsTopicPolicyException create() =>
      InsufficientSnsTopicPolicyException._();
  @$core.override
  InsufficientSnsTopicPolicyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsufficientSnsTopicPolicyException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InsufficientSnsTopicPolicyException>(create);
  static InsufficientSnsTopicPolicyException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidCloudWatchLogsLogGroupArnException extends $pb.GeneratedMessage {
  factory InvalidCloudWatchLogsLogGroupArnException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidCloudWatchLogsLogGroupArnException._();

  factory InvalidCloudWatchLogsLogGroupArnException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidCloudWatchLogsLogGroupArnException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidCloudWatchLogsLogGroupArnException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidCloudWatchLogsLogGroupArnException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidCloudWatchLogsLogGroupArnException copyWith(
          void Function(InvalidCloudWatchLogsLogGroupArnException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidCloudWatchLogsLogGroupArnException))
          as InvalidCloudWatchLogsLogGroupArnException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidCloudWatchLogsLogGroupArnException create() =>
      InvalidCloudWatchLogsLogGroupArnException._();
  @$core.override
  InvalidCloudWatchLogsLogGroupArnException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidCloudWatchLogsLogGroupArnException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidCloudWatchLogsLogGroupArnException>(create);
  static InvalidCloudWatchLogsLogGroupArnException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidCloudWatchLogsRoleArnException extends $pb.GeneratedMessage {
  factory InvalidCloudWatchLogsRoleArnException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidCloudWatchLogsRoleArnException._();

  factory InvalidCloudWatchLogsRoleArnException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidCloudWatchLogsRoleArnException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidCloudWatchLogsRoleArnException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidCloudWatchLogsRoleArnException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidCloudWatchLogsRoleArnException copyWith(
          void Function(InvalidCloudWatchLogsRoleArnException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidCloudWatchLogsRoleArnException))
          as InvalidCloudWatchLogsRoleArnException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidCloudWatchLogsRoleArnException create() =>
      InvalidCloudWatchLogsRoleArnException._();
  @$core.override
  InvalidCloudWatchLogsRoleArnException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidCloudWatchLogsRoleArnException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidCloudWatchLogsRoleArnException>(create);
  static InvalidCloudWatchLogsRoleArnException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidDateRangeException extends $pb.GeneratedMessage {
  factory InvalidDateRangeException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidDateRangeException._();

  factory InvalidDateRangeException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidDateRangeException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidDateRangeException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDateRangeException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDateRangeException copyWith(
          void Function(InvalidDateRangeException) updates) =>
      super.copyWith((message) => updates(message as InvalidDateRangeException))
          as InvalidDateRangeException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidDateRangeException create() => InvalidDateRangeException._();
  @$core.override
  InvalidDateRangeException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidDateRangeException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidDateRangeException>(create);
  static InvalidDateRangeException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidEventCategoryException extends $pb.GeneratedMessage {
  factory InvalidEventCategoryException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidEventCategoryException._();

  factory InvalidEventCategoryException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidEventCategoryException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidEventCategoryException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventCategoryException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventCategoryException copyWith(
          void Function(InvalidEventCategoryException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidEventCategoryException))
          as InvalidEventCategoryException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidEventCategoryException create() =>
      InvalidEventCategoryException._();
  @$core.override
  InvalidEventCategoryException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidEventCategoryException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidEventCategoryException>(create);
  static InvalidEventCategoryException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidEventDataStoreCategoryException extends $pb.GeneratedMessage {
  factory InvalidEventDataStoreCategoryException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidEventDataStoreCategoryException._();

  factory InvalidEventDataStoreCategoryException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidEventDataStoreCategoryException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidEventDataStoreCategoryException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventDataStoreCategoryException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventDataStoreCategoryException copyWith(
          void Function(InvalidEventDataStoreCategoryException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidEventDataStoreCategoryException))
          as InvalidEventDataStoreCategoryException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidEventDataStoreCategoryException create() =>
      InvalidEventDataStoreCategoryException._();
  @$core.override
  InvalidEventDataStoreCategoryException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidEventDataStoreCategoryException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidEventDataStoreCategoryException>(create);
  static InvalidEventDataStoreCategoryException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidEventDataStoreStatusException extends $pb.GeneratedMessage {
  factory InvalidEventDataStoreStatusException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidEventDataStoreStatusException._();

  factory InvalidEventDataStoreStatusException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidEventDataStoreStatusException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidEventDataStoreStatusException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventDataStoreStatusException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventDataStoreStatusException copyWith(
          void Function(InvalidEventDataStoreStatusException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidEventDataStoreStatusException))
          as InvalidEventDataStoreStatusException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidEventDataStoreStatusException create() =>
      InvalidEventDataStoreStatusException._();
  @$core.override
  InvalidEventDataStoreStatusException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidEventDataStoreStatusException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidEventDataStoreStatusException>(create);
  static InvalidEventDataStoreStatusException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidEventSelectorsException extends $pb.GeneratedMessage {
  factory InvalidEventSelectorsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidEventSelectorsException._();

  factory InvalidEventSelectorsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidEventSelectorsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidEventSelectorsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventSelectorsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventSelectorsException copyWith(
          void Function(InvalidEventSelectorsException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidEventSelectorsException))
          as InvalidEventSelectorsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidEventSelectorsException create() =>
      InvalidEventSelectorsException._();
  @$core.override
  InvalidEventSelectorsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidEventSelectorsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidEventSelectorsException>(create);
  static InvalidEventSelectorsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidHomeRegionException extends $pb.GeneratedMessage {
  factory InvalidHomeRegionException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidHomeRegionException._();

  factory InvalidHomeRegionException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidHomeRegionException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidHomeRegionException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidHomeRegionException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidHomeRegionException copyWith(
          void Function(InvalidHomeRegionException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidHomeRegionException))
          as InvalidHomeRegionException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidHomeRegionException create() => InvalidHomeRegionException._();
  @$core.override
  InvalidHomeRegionException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidHomeRegionException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidHomeRegionException>(create);
  static InvalidHomeRegionException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidImportSourceException extends $pb.GeneratedMessage {
  factory InvalidImportSourceException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidImportSourceException._();

  factory InvalidImportSourceException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidImportSourceException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidImportSourceException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidImportSourceException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidImportSourceException copyWith(
          void Function(InvalidImportSourceException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidImportSourceException))
          as InvalidImportSourceException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidImportSourceException create() =>
      InvalidImportSourceException._();
  @$core.override
  InvalidImportSourceException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidImportSourceException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidImportSourceException>(create);
  static InvalidImportSourceException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidInsightSelectorsException extends $pb.GeneratedMessage {
  factory InvalidInsightSelectorsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidInsightSelectorsException._();

  factory InvalidInsightSelectorsException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidInsightSelectorsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidInsightSelectorsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidInsightSelectorsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidInsightSelectorsException copyWith(
          void Function(InvalidInsightSelectorsException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidInsightSelectorsException))
          as InvalidInsightSelectorsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidInsightSelectorsException create() =>
      InvalidInsightSelectorsException._();
  @$core.override
  InvalidInsightSelectorsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidInsightSelectorsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidInsightSelectorsException>(
          create);
  static InvalidInsightSelectorsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidKmsKeyIdException extends $pb.GeneratedMessage {
  factory InvalidKmsKeyIdException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidKmsKeyIdException._();

  factory InvalidKmsKeyIdException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidKmsKeyIdException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidKmsKeyIdException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKmsKeyIdException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKmsKeyIdException copyWith(
          void Function(InvalidKmsKeyIdException) updates) =>
      super.copyWith((message) => updates(message as InvalidKmsKeyIdException))
          as InvalidKmsKeyIdException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidKmsKeyIdException create() => InvalidKmsKeyIdException._();
  @$core.override
  InvalidKmsKeyIdException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidKmsKeyIdException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidKmsKeyIdException>(create);
  static InvalidKmsKeyIdException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidLookupAttributesException extends $pb.GeneratedMessage {
  factory InvalidLookupAttributesException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidLookupAttributesException._();

  factory InvalidLookupAttributesException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidLookupAttributesException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidLookupAttributesException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidLookupAttributesException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidLookupAttributesException copyWith(
          void Function(InvalidLookupAttributesException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidLookupAttributesException))
          as InvalidLookupAttributesException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidLookupAttributesException create() =>
      InvalidLookupAttributesException._();
  @$core.override
  InvalidLookupAttributesException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidLookupAttributesException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidLookupAttributesException>(
          create);
  static InvalidLookupAttributesException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidMaxResultsException extends $pb.GeneratedMessage {
  factory InvalidMaxResultsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidMaxResultsException._();

  factory InvalidMaxResultsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidMaxResultsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidMaxResultsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidMaxResultsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidMaxResultsException copyWith(
          void Function(InvalidMaxResultsException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidMaxResultsException))
          as InvalidMaxResultsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidMaxResultsException create() => InvalidMaxResultsException._();
  @$core.override
  InvalidMaxResultsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidMaxResultsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidMaxResultsException>(create);
  static InvalidMaxResultsException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
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

class InvalidParameterCombinationException extends $pb.GeneratedMessage {
  factory InvalidParameterCombinationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidParameterCombinationException._();

  factory InvalidParameterCombinationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidParameterCombinationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidParameterCombinationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterCombinationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterCombinationException copyWith(
          void Function(InvalidParameterCombinationException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidParameterCombinationException))
          as InvalidParameterCombinationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidParameterCombinationException create() =>
      InvalidParameterCombinationException._();
  @$core.override
  InvalidParameterCombinationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidParameterCombinationException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidParameterCombinationException>(create);
  static InvalidParameterCombinationException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
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

class InvalidQueryStatementException extends $pb.GeneratedMessage {
  factory InvalidQueryStatementException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidQueryStatementException._();

  factory InvalidQueryStatementException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidQueryStatementException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidQueryStatementException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidQueryStatementException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidQueryStatementException copyWith(
          void Function(InvalidQueryStatementException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidQueryStatementException))
          as InvalidQueryStatementException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidQueryStatementException create() =>
      InvalidQueryStatementException._();
  @$core.override
  InvalidQueryStatementException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidQueryStatementException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidQueryStatementException>(create);
  static InvalidQueryStatementException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidQueryStatusException extends $pb.GeneratedMessage {
  factory InvalidQueryStatusException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidQueryStatusException._();

  factory InvalidQueryStatusException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidQueryStatusException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidQueryStatusException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidQueryStatusException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidQueryStatusException copyWith(
          void Function(InvalidQueryStatusException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidQueryStatusException))
          as InvalidQueryStatusException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidQueryStatusException create() =>
      InvalidQueryStatusException._();
  @$core.override
  InvalidQueryStatusException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidQueryStatusException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidQueryStatusException>(create);
  static InvalidQueryStatusException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidS3BucketNameException extends $pb.GeneratedMessage {
  factory InvalidS3BucketNameException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidS3BucketNameException._();

  factory InvalidS3BucketNameException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidS3BucketNameException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidS3BucketNameException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidS3BucketNameException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidS3BucketNameException copyWith(
          void Function(InvalidS3BucketNameException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidS3BucketNameException))
          as InvalidS3BucketNameException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidS3BucketNameException create() =>
      InvalidS3BucketNameException._();
  @$core.override
  InvalidS3BucketNameException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidS3BucketNameException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidS3BucketNameException>(create);
  static InvalidS3BucketNameException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidS3PrefixException extends $pb.GeneratedMessage {
  factory InvalidS3PrefixException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidS3PrefixException._();

  factory InvalidS3PrefixException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidS3PrefixException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidS3PrefixException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidS3PrefixException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidS3PrefixException copyWith(
          void Function(InvalidS3PrefixException) updates) =>
      super.copyWith((message) => updates(message as InvalidS3PrefixException))
          as InvalidS3PrefixException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidS3PrefixException create() => InvalidS3PrefixException._();
  @$core.override
  InvalidS3PrefixException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidS3PrefixException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidS3PrefixException>(create);
  static InvalidS3PrefixException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidSnsTopicNameException extends $pb.GeneratedMessage {
  factory InvalidSnsTopicNameException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidSnsTopicNameException._();

  factory InvalidSnsTopicNameException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidSnsTopicNameException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidSnsTopicNameException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSnsTopicNameException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSnsTopicNameException copyWith(
          void Function(InvalidSnsTopicNameException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidSnsTopicNameException))
          as InvalidSnsTopicNameException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidSnsTopicNameException create() =>
      InvalidSnsTopicNameException._();
  @$core.override
  InvalidSnsTopicNameException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidSnsTopicNameException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidSnsTopicNameException>(create);
  static InvalidSnsTopicNameException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidSourceException extends $pb.GeneratedMessage {
  factory InvalidSourceException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidSourceException._();

  factory InvalidSourceException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidSourceException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidSourceException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSourceException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSourceException copyWith(
          void Function(InvalidSourceException) updates) =>
      super.copyWith((message) => updates(message as InvalidSourceException))
          as InvalidSourceException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidSourceException create() => InvalidSourceException._();
  @$core.override
  InvalidSourceException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidSourceException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidSourceException>(create);
  static InvalidSourceException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidTagParameterException extends $pb.GeneratedMessage {
  factory InvalidTagParameterException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidTagParameterException._();

  factory InvalidTagParameterException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidTagParameterException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidTagParameterException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTagParameterException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTagParameterException copyWith(
          void Function(InvalidTagParameterException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidTagParameterException))
          as InvalidTagParameterException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidTagParameterException create() =>
      InvalidTagParameterException._();
  @$core.override
  InvalidTagParameterException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidTagParameterException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidTagParameterException>(create);
  static InvalidTagParameterException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidTimeRangeException extends $pb.GeneratedMessage {
  factory InvalidTimeRangeException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidTimeRangeException._();

  factory InvalidTimeRangeException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidTimeRangeException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidTimeRangeException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTimeRangeException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTimeRangeException copyWith(
          void Function(InvalidTimeRangeException) updates) =>
      super.copyWith((message) => updates(message as InvalidTimeRangeException))
          as InvalidTimeRangeException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidTimeRangeException create() => InvalidTimeRangeException._();
  @$core.override
  InvalidTimeRangeException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidTimeRangeException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidTimeRangeException>(create);
  static InvalidTimeRangeException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidTokenException extends $pb.GeneratedMessage {
  factory InvalidTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidTokenException._();

  factory InvalidTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTokenException copyWith(
          void Function(InvalidTokenException) updates) =>
      super.copyWith((message) => updates(message as InvalidTokenException))
          as InvalidTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidTokenException create() => InvalidTokenException._();
  @$core.override
  InvalidTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidTokenException>(create);
  static InvalidTokenException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidTrailNameException extends $pb.GeneratedMessage {
  factory InvalidTrailNameException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidTrailNameException._();

  factory InvalidTrailNameException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidTrailNameException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidTrailNameException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTrailNameException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTrailNameException copyWith(
          void Function(InvalidTrailNameException) updates) =>
      super.copyWith((message) => updates(message as InvalidTrailNameException))
          as InvalidTrailNameException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidTrailNameException create() => InvalidTrailNameException._();
  @$core.override
  InvalidTrailNameException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidTrailNameException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidTrailNameException>(create);
  static InvalidTrailNameException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class KmsException extends $pb.GeneratedMessage {
  factory KmsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsException._();

  factory KmsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsException copyWith(void Function(KmsException) updates) =>
      super.copyWith((message) => updates(message as KmsException))
          as KmsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsException create() => KmsException._();
  @$core.override
  KmsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsException>(create);
  static KmsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class KmsKeyDisabledException extends $pb.GeneratedMessage {
  factory KmsKeyDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsKeyDisabledException._();

  factory KmsKeyDisabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsKeyDisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsKeyDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsKeyDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsKeyDisabledException copyWith(
          void Function(KmsKeyDisabledException) updates) =>
      super.copyWith((message) => updates(message as KmsKeyDisabledException))
          as KmsKeyDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsKeyDisabledException create() => KmsKeyDisabledException._();
  @$core.override
  KmsKeyDisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsKeyDisabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsKeyDisabledException>(create);
  static KmsKeyDisabledException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class KmsKeyNotFoundException extends $pb.GeneratedMessage {
  factory KmsKeyNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsKeyNotFoundException._();

  factory KmsKeyNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsKeyNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsKeyNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsKeyNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsKeyNotFoundException copyWith(
          void Function(KmsKeyNotFoundException) updates) =>
      super.copyWith((message) => updates(message as KmsKeyNotFoundException))
          as KmsKeyNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsKeyNotFoundException create() => KmsKeyNotFoundException._();
  @$core.override
  KmsKeyNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsKeyNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsKeyNotFoundException>(create);
  static KmsKeyNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ListChannelsRequest extends $pb.GeneratedMessage {
  factory ListChannelsRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListChannelsRequest._();

  factory ListChannelsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListChannelsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListChannelsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChannelsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChannelsRequest copyWith(void Function(ListChannelsRequest) updates) =>
      super.copyWith((message) => updates(message as ListChannelsRequest))
          as ListChannelsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListChannelsRequest create() => ListChannelsRequest._();
  @$core.override
  ListChannelsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListChannelsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListChannelsRequest>(create);
  static ListChannelsRequest? _defaultInstance;

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

class ListChannelsResponse extends $pb.GeneratedMessage {
  factory ListChannelsResponse({
    $core.Iterable<Channel>? channels,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (channels != null) result.channels.addAll(channels);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListChannelsResponse._();

  factory ListChannelsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListChannelsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListChannelsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<Channel>(155227286, _omitFieldNames ? '' : 'channels',
        subBuilder: Channel.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChannelsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChannelsResponse copyWith(void Function(ListChannelsResponse) updates) =>
      super.copyWith((message) => updates(message as ListChannelsResponse))
          as ListChannelsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListChannelsResponse create() => ListChannelsResponse._();
  @$core.override
  ListChannelsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListChannelsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListChannelsResponse>(create);
  static ListChannelsResponse? _defaultInstance;

  @$pb.TagNumber(155227286)
  $pb.PbList<Channel> get channels => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListDashboardsRequest extends $pb.GeneratedMessage {
  factory ListDashboardsRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
    DashboardType? type,
    $core.String? nameprefix,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (type != null) result.type = type;
    if (nameprefix != null) result.nameprefix = nameprefix;
    return result;
  }

  ListDashboardsRequest._();

  factory ListDashboardsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDashboardsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDashboardsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aE<DashboardType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DashboardType.values)
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsRequest copyWith(
          void Function(ListDashboardsRequest) updates) =>
      super.copyWith((message) => updates(message as ListDashboardsRequest))
          as ListDashboardsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDashboardsRequest create() => ListDashboardsRequest._();
  @$core.override
  ListDashboardsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDashboardsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDashboardsRequest>(create);
  static ListDashboardsRequest? _defaultInstance;

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

  @$pb.TagNumber(290836590)
  DashboardType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(DashboardType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(3);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(3, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(3);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);
}

class ListDashboardsResponse extends $pb.GeneratedMessage {
  factory ListDashboardsResponse({
    $core.String? nexttoken,
    $core.Iterable<DashboardDetail>? dashboards,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (dashboards != null) result.dashboards.addAll(dashboards);
    return result;
  }

  ListDashboardsResponse._();

  factory ListDashboardsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDashboardsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDashboardsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<DashboardDetail>(261493913, _omitFieldNames ? '' : 'dashboards',
        subBuilder: DashboardDetail.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsResponse copyWith(
          void Function(ListDashboardsResponse) updates) =>
      super.copyWith((message) => updates(message as ListDashboardsResponse))
          as ListDashboardsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDashboardsResponse create() => ListDashboardsResponse._();
  @$core.override
  ListDashboardsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDashboardsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDashboardsResponse>(create);
  static ListDashboardsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(261493913)
  $pb.PbList<DashboardDetail> get dashboards => $_getList(1);
}

class ListEventDataStoresRequest extends $pb.GeneratedMessage {
  factory ListEventDataStoresRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListEventDataStoresRequest._();

  factory ListEventDataStoresRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEventDataStoresRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEventDataStoresRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventDataStoresRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventDataStoresRequest copyWith(
          void Function(ListEventDataStoresRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListEventDataStoresRequest))
          as ListEventDataStoresRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEventDataStoresRequest create() => ListEventDataStoresRequest._();
  @$core.override
  ListEventDataStoresRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEventDataStoresRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEventDataStoresRequest>(create);
  static ListEventDataStoresRequest? _defaultInstance;

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

class ListEventDataStoresResponse extends $pb.GeneratedMessage {
  factory ListEventDataStoresResponse({
    $core.Iterable<EventDataStore>? eventdatastores,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (eventdatastores != null) result.eventdatastores.addAll(eventdatastores);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListEventDataStoresResponse._();

  factory ListEventDataStoresResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEventDataStoresResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEventDataStoresResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<EventDataStore>(152154314, _omitFieldNames ? '' : 'eventdatastores',
        subBuilder: EventDataStore.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventDataStoresResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventDataStoresResponse copyWith(
          void Function(ListEventDataStoresResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListEventDataStoresResponse))
          as ListEventDataStoresResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEventDataStoresResponse create() =>
      ListEventDataStoresResponse._();
  @$core.override
  ListEventDataStoresResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEventDataStoresResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEventDataStoresResponse>(create);
  static ListEventDataStoresResponse? _defaultInstance;

  @$pb.TagNumber(152154314)
  $pb.PbList<EventDataStore> get eventdatastores => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListImportFailuresRequest extends $pb.GeneratedMessage {
  factory ListImportFailuresRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? importid,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (importid != null) result.importid = importid;
    return result;
  }

  ListImportFailuresRequest._();

  factory ListImportFailuresRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListImportFailuresRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListImportFailuresRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportFailuresRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportFailuresRequest copyWith(
          void Function(ListImportFailuresRequest) updates) =>
      super.copyWith((message) => updates(message as ListImportFailuresRequest))
          as ListImportFailuresRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListImportFailuresRequest create() => ListImportFailuresRequest._();
  @$core.override
  ListImportFailuresRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListImportFailuresRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListImportFailuresRequest>(create);
  static ListImportFailuresRequest? _defaultInstance;

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

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(2);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(2);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class ListImportFailuresResponse extends $pb.GeneratedMessage {
  factory ListImportFailuresResponse({
    $core.String? nexttoken,
    $core.Iterable<ImportFailureListItem>? failures,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (failures != null) result.failures.addAll(failures);
    return result;
  }

  ListImportFailuresResponse._();

  factory ListImportFailuresResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListImportFailuresResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListImportFailuresResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ImportFailureListItem>(335467271, _omitFieldNames ? '' : 'failures',
        subBuilder: ImportFailureListItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportFailuresResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportFailuresResponse copyWith(
          void Function(ListImportFailuresResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListImportFailuresResponse))
          as ListImportFailuresResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListImportFailuresResponse create() => ListImportFailuresResponse._();
  @$core.override
  ListImportFailuresResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListImportFailuresResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListImportFailuresResponse>(create);
  static ListImportFailuresResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(335467271)
  $pb.PbList<ImportFailureListItem> get failures => $_getList(1);
}

class ListImportsRequest extends $pb.GeneratedMessage {
  factory ListImportsRequest({
    ImportStatus? importstatus,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? destination,
  }) {
    final result = create();
    if (importstatus != null) result.importstatus = importstatus;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (destination != null) result.destination = destination;
    return result;
  }

  ListImportsRequest._();

  factory ListImportsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListImportsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListImportsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<ImportStatus>(129077631, _omitFieldNames ? '' : 'importstatus',
        enumValues: ImportStatus.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(457443680, _omitFieldNames ? '' : 'destination')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsRequest copyWith(void Function(ListImportsRequest) updates) =>
      super.copyWith((message) => updates(message as ListImportsRequest))
          as ListImportsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListImportsRequest create() => ListImportsRequest._();
  @$core.override
  ListImportsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListImportsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListImportsRequest>(create);
  static ListImportsRequest? _defaultInstance;

  @$pb.TagNumber(129077631)
  ImportStatus get importstatus => $_getN(0);
  @$pb.TagNumber(129077631)
  set importstatus(ImportStatus value) => $_setField(129077631, value);
  @$pb.TagNumber(129077631)
  $core.bool hasImportstatus() => $_has(0);
  @$pb.TagNumber(129077631)
  void clearImportstatus() => $_clearField(129077631);

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

  @$pb.TagNumber(457443680)
  $core.String get destination => $_getSZ(3);
  @$pb.TagNumber(457443680)
  set destination($core.String value) => $_setString(3, value);
  @$pb.TagNumber(457443680)
  $core.bool hasDestination() => $_has(3);
  @$pb.TagNumber(457443680)
  void clearDestination() => $_clearField(457443680);
}

class ListImportsResponse extends $pb.GeneratedMessage {
  factory ListImportsResponse({
    $core.String? nexttoken,
    $core.Iterable<ImportsListItem>? imports,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (imports != null) result.imports.addAll(imports);
    return result;
  }

  ListImportsResponse._();

  factory ListImportsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListImportsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListImportsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ImportsListItem>(497717894, _omitFieldNames ? '' : 'imports',
        subBuilder: ImportsListItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsResponse copyWith(void Function(ListImportsResponse) updates) =>
      super.copyWith((message) => updates(message as ListImportsResponse))
          as ListImportsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListImportsResponse create() => ListImportsResponse._();
  @$core.override
  ListImportsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListImportsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListImportsResponse>(create);
  static ListImportsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(497717894)
  $pb.PbList<ImportsListItem> get imports => $_getList(1);
}

class ListInsightsDataRequest extends $pb.GeneratedMessage {
  factory ListInsightsDataRequest({
    $core.String? endtime,
    ListInsightsDataType? datatype,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? starttime,
    $core.String? insightsource,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? dimensions,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (datatype != null) result.datatype = datatype;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (starttime != null) result.starttime = starttime;
    if (insightsource != null) result.insightsource = insightsource;
    if (dimensions != null) result.dimensions.addEntries(dimensions);
    return result;
  }

  ListInsightsDataRequest._();

  factory ListInsightsDataRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListInsightsDataRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListInsightsDataRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aE<ListInsightsDataType>(67988590, _omitFieldNames ? '' : 'datatype',
        enumValues: ListInsightsDataType.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..aOS(438381967, _omitFieldNames ? '' : 'insightsource')
    ..m<$core.String, $core.String>(
        462933457, _omitFieldNames ? '' : 'dimensions',
        entryClassName: 'ListInsightsDataRequest.DimensionsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cloudtrail'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsDataRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsDataRequest copyWith(
          void Function(ListInsightsDataRequest) updates) =>
      super.copyWith((message) => updates(message as ListInsightsDataRequest))
          as ListInsightsDataRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListInsightsDataRequest create() => ListInsightsDataRequest._();
  @$core.override
  ListInsightsDataRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListInsightsDataRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListInsightsDataRequest>(create);
  static ListInsightsDataRequest? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(67988590)
  ListInsightsDataType get datatype => $_getN(1);
  @$pb.TagNumber(67988590)
  set datatype(ListInsightsDataType value) => $_setField(67988590, value);
  @$pb.TagNumber(67988590)
  $core.bool hasDatatype() => $_has(1);
  @$pb.TagNumber(67988590)
  void clearDatatype() => $_clearField(67988590);

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

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(4);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(4);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);

  @$pb.TagNumber(438381967)
  $core.String get insightsource => $_getSZ(5);
  @$pb.TagNumber(438381967)
  set insightsource($core.String value) => $_setString(5, value);
  @$pb.TagNumber(438381967)
  $core.bool hasInsightsource() => $_has(5);
  @$pb.TagNumber(438381967)
  void clearInsightsource() => $_clearField(438381967);

  @$pb.TagNumber(462933457)
  $pb.PbMap<$core.String, $core.String> get dimensions => $_getMap(6);
}

class ListInsightsDataResponse extends $pb.GeneratedMessage {
  factory ListInsightsDataResponse({
    $core.Iterable<Event>? events,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (events != null) result.events.addAll(events);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListInsightsDataResponse._();

  factory ListInsightsDataResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListInsightsDataResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListInsightsDataResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<Event>(3416229, _omitFieldNames ? '' : 'events',
        subBuilder: Event.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsDataResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsDataResponse copyWith(
          void Function(ListInsightsDataResponse) updates) =>
      super.copyWith((message) => updates(message as ListInsightsDataResponse))
          as ListInsightsDataResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListInsightsDataResponse create() => ListInsightsDataResponse._();
  @$core.override
  ListInsightsDataResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListInsightsDataResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListInsightsDataResponse>(create);
  static ListInsightsDataResponse? _defaultInstance;

  @$pb.TagNumber(3416229)
  $pb.PbList<Event> get events => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListInsightsMetricDataRequest extends $pb.GeneratedMessage {
  factory ListInsightsMetricDataRequest({
    $core.String? errorcode,
    $core.String? eventsource,
    $core.String? endtime,
    InsightsMetricDataType? datatype,
    $core.int? period,
    $core.String? nexttoken,
    $core.String? eventname,
    $core.int? maxresults,
    $core.String? starttime,
    $core.String? trailname,
    InsightType? insighttype,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (eventsource != null) result.eventsource = eventsource;
    if (endtime != null) result.endtime = endtime;
    if (datatype != null) result.datatype = datatype;
    if (period != null) result.period = period;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (eventname != null) result.eventname = eventname;
    if (maxresults != null) result.maxresults = maxresults;
    if (starttime != null) result.starttime = starttime;
    if (trailname != null) result.trailname = trailname;
    if (insighttype != null) result.insighttype = insighttype;
    return result;
  }

  ListInsightsMetricDataRequest._();

  factory ListInsightsMetricDataRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListInsightsMetricDataRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListInsightsMetricDataRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(37841339, _omitFieldNames ? '' : 'eventsource')
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aE<InsightsMetricDataType>(67988590, _omitFieldNames ? '' : 'datatype',
        enumValues: InsightsMetricDataType.values)
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(264746781, _omitFieldNames ? '' : 'eventname')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..aOS(507774985, _omitFieldNames ? '' : 'trailname')
    ..aE<InsightType>(530375860, _omitFieldNames ? '' : 'insighttype',
        enumValues: InsightType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsMetricDataRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsMetricDataRequest copyWith(
          void Function(ListInsightsMetricDataRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListInsightsMetricDataRequest))
          as ListInsightsMetricDataRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListInsightsMetricDataRequest create() =>
      ListInsightsMetricDataRequest._();
  @$core.override
  ListInsightsMetricDataRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListInsightsMetricDataRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListInsightsMetricDataRequest>(create);
  static ListInsightsMetricDataRequest? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(37841339)
  $core.String get eventsource => $_getSZ(1);
  @$pb.TagNumber(37841339)
  set eventsource($core.String value) => $_setString(1, value);
  @$pb.TagNumber(37841339)
  $core.bool hasEventsource() => $_has(1);
  @$pb.TagNumber(37841339)
  void clearEventsource() => $_clearField(37841339);

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(2);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(2);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(67988590)
  InsightsMetricDataType get datatype => $_getN(3);
  @$pb.TagNumber(67988590)
  set datatype(InsightsMetricDataType value) => $_setField(67988590, value);
  @$pb.TagNumber(67988590)
  $core.bool hasDatatype() => $_has(3);
  @$pb.TagNumber(67988590)
  void clearDatatype() => $_clearField(67988590);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(4);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(4);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(5);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(5, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(5);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(264746781)
  $core.String get eventname => $_getSZ(6);
  @$pb.TagNumber(264746781)
  set eventname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(264746781)
  $core.bool hasEventname() => $_has(6);
  @$pb.TagNumber(264746781)
  void clearEventname() => $_clearField(264746781);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(7);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(7, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(7);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(8);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(8, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(8);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);

  @$pb.TagNumber(507774985)
  $core.String get trailname => $_getSZ(9);
  @$pb.TagNumber(507774985)
  set trailname($core.String value) => $_setString(9, value);
  @$pb.TagNumber(507774985)
  $core.bool hasTrailname() => $_has(9);
  @$pb.TagNumber(507774985)
  void clearTrailname() => $_clearField(507774985);

  @$pb.TagNumber(530375860)
  InsightType get insighttype => $_getN(10);
  @$pb.TagNumber(530375860)
  set insighttype(InsightType value) => $_setField(530375860, value);
  @$pb.TagNumber(530375860)
  $core.bool hasInsighttype() => $_has(10);
  @$pb.TagNumber(530375860)
  void clearInsighttype() => $_clearField(530375860);
}

class ListInsightsMetricDataResponse extends $pb.GeneratedMessage {
  factory ListInsightsMetricDataResponse({
    $core.String? errorcode,
    $core.String? eventsource,
    $core.String? trailarn,
    $core.Iterable<$core.String>? timestamps,
    $core.String? nexttoken,
    $core.Iterable<$core.double>? values,
    $core.String? eventname,
    InsightType? insighttype,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (eventsource != null) result.eventsource = eventsource;
    if (trailarn != null) result.trailarn = trailarn;
    if (timestamps != null) result.timestamps.addAll(timestamps);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (values != null) result.values.addAll(values);
    if (eventname != null) result.eventname = eventname;
    if (insighttype != null) result.insighttype = insighttype;
    return result;
  }

  ListInsightsMetricDataResponse._();

  factory ListInsightsMetricDataResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListInsightsMetricDataResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListInsightsMetricDataResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(37841339, _omitFieldNames ? '' : 'eventsource')
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..pPS(213534737, _omitFieldNames ? '' : 'timestamps')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..p<$core.double>(
        223158876, _omitFieldNames ? '' : 'values', $pb.PbFieldType.KD)
    ..aOS(264746781, _omitFieldNames ? '' : 'eventname')
    ..aE<InsightType>(530375860, _omitFieldNames ? '' : 'insighttype',
        enumValues: InsightType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsMetricDataResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListInsightsMetricDataResponse copyWith(
          void Function(ListInsightsMetricDataResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListInsightsMetricDataResponse))
          as ListInsightsMetricDataResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListInsightsMetricDataResponse create() =>
      ListInsightsMetricDataResponse._();
  @$core.override
  ListInsightsMetricDataResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListInsightsMetricDataResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListInsightsMetricDataResponse>(create);
  static ListInsightsMetricDataResponse? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(37841339)
  $core.String get eventsource => $_getSZ(1);
  @$pb.TagNumber(37841339)
  set eventsource($core.String value) => $_setString(1, value);
  @$pb.TagNumber(37841339)
  $core.bool hasEventsource() => $_has(1);
  @$pb.TagNumber(37841339)
  void clearEventsource() => $_clearField(37841339);

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(2);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(2);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(213534737)
  $pb.PbList<$core.String> get timestamps => $_getList(3);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(4);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(4);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(223158876)
  $pb.PbList<$core.double> get values => $_getList(5);

  @$pb.TagNumber(264746781)
  $core.String get eventname => $_getSZ(6);
  @$pb.TagNumber(264746781)
  set eventname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(264746781)
  $core.bool hasEventname() => $_has(6);
  @$pb.TagNumber(264746781)
  void clearEventname() => $_clearField(264746781);

  @$pb.TagNumber(530375860)
  InsightType get insighttype => $_getN(7);
  @$pb.TagNumber(530375860)
  set insighttype(InsightType value) => $_setField(530375860, value);
  @$pb.TagNumber(530375860)
  $core.bool hasInsighttype() => $_has(7);
  @$pb.TagNumber(530375860)
  void clearInsighttype() => $_clearField(530375860);
}

class ListPublicKeysRequest extends $pb.GeneratedMessage {
  factory ListPublicKeysRequest({
    $core.String? endtime,
    $core.String? nexttoken,
    $core.String? starttime,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (starttime != null) result.starttime = starttime;
    return result;
  }

  ListPublicKeysRequest._();

  factory ListPublicKeysRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPublicKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPublicKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPublicKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPublicKeysRequest copyWith(
          void Function(ListPublicKeysRequest) updates) =>
      super.copyWith((message) => updates(message as ListPublicKeysRequest))
          as ListPublicKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPublicKeysRequest create() => ListPublicKeysRequest._();
  @$core.override
  ListPublicKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPublicKeysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPublicKeysRequest>(create);
  static ListPublicKeysRequest? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(2);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(2);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);
}

class ListPublicKeysResponse extends $pb.GeneratedMessage {
  factory ListPublicKeysResponse({
    $core.Iterable<PublicKey>? publickeylist,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (publickeylist != null) result.publickeylist.addAll(publickeylist);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListPublicKeysResponse._();

  factory ListPublicKeysResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPublicKeysResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPublicKeysResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<PublicKey>(129933128, _omitFieldNames ? '' : 'publickeylist',
        subBuilder: PublicKey.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPublicKeysResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPublicKeysResponse copyWith(
          void Function(ListPublicKeysResponse) updates) =>
      super.copyWith((message) => updates(message as ListPublicKeysResponse))
          as ListPublicKeysResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPublicKeysResponse create() => ListPublicKeysResponse._();
  @$core.override
  ListPublicKeysResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPublicKeysResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPublicKeysResponse>(create);
  static ListPublicKeysResponse? _defaultInstance;

  @$pb.TagNumber(129933128)
  $pb.PbList<PublicKey> get publickeylist => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListQueriesRequest extends $pb.GeneratedMessage {
  factory ListQueriesRequest({
    $core.String? endtime,
    $core.String? eventdatastore,
    $core.String? nexttoken,
    $core.int? maxresults,
    QueryStatus? querystatus,
    $core.String? starttime,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (querystatus != null) result.querystatus = querystatus;
    if (starttime != null) result.starttime = starttime;
    return result;
  }

  ListQueriesRequest._();

  factory ListQueriesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueriesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueriesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aE<QueryStatus>(367016406, _omitFieldNames ? '' : 'querystatus',
        enumValues: QueryStatus.values)
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueriesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueriesRequest copyWith(void Function(ListQueriesRequest) updates) =>
      super.copyWith((message) => updates(message as ListQueriesRequest))
          as ListQueriesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueriesRequest create() => ListQueriesRequest._();
  @$core.override
  ListQueriesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueriesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueriesRequest>(create);
  static ListQueriesRequest? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(1);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(1, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(1);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

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

  @$pb.TagNumber(367016406)
  QueryStatus get querystatus => $_getN(4);
  @$pb.TagNumber(367016406)
  set querystatus(QueryStatus value) => $_setField(367016406, value);
  @$pb.TagNumber(367016406)
  $core.bool hasQuerystatus() => $_has(4);
  @$pb.TagNumber(367016406)
  void clearQuerystatus() => $_clearField(367016406);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(5);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(5);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);
}

class ListQueriesResponse extends $pb.GeneratedMessage {
  factory ListQueriesResponse({
    $core.String? nexttoken,
    $core.Iterable<Query>? queries,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (queries != null) result.queries.addAll(queries);
    return result;
  }

  ListQueriesResponse._();

  factory ListQueriesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueriesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueriesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Query>(305659620, _omitFieldNames ? '' : 'queries',
        subBuilder: Query.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueriesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueriesResponse copyWith(void Function(ListQueriesResponse) updates) =>
      super.copyWith((message) => updates(message as ListQueriesResponse))
          as ListQueriesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueriesResponse create() => ListQueriesResponse._();
  @$core.override
  ListQueriesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueriesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueriesResponse>(create);
  static ListQueriesResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(305659620)
  $pb.PbList<Query> get queries => $_getList(1);
}

class ListTagsRequest extends $pb.GeneratedMessage {
  factory ListTagsRequest({
    $core.String? nexttoken,
    $core.Iterable<$core.String>? resourceidlist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (resourceidlist != null) result.resourceidlist.addAll(resourceidlist);
    return result;
  }

  ListTagsRequest._();

  factory ListTagsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPS(272650681, _omitFieldNames ? '' : 'resourceidlist')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsRequest copyWith(void Function(ListTagsRequest) updates) =>
      super.copyWith((message) => updates(message as ListTagsRequest))
          as ListTagsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsRequest create() => ListTagsRequest._();
  @$core.override
  ListTagsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsRequest>(create);
  static ListTagsRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(272650681)
  $pb.PbList<$core.String> get resourceidlist => $_getList(1);
}

class ListTagsResponse extends $pb.GeneratedMessage {
  factory ListTagsResponse({
    $core.String? nexttoken,
    $core.Iterable<ResourceTag>? resourcetaglist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (resourcetaglist != null) result.resourcetaglist.addAll(resourcetaglist);
    return result;
  }

  ListTagsResponse._();

  factory ListTagsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ResourceTag>(301974828, _omitFieldNames ? '' : 'resourcetaglist',
        subBuilder: ResourceTag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsResponse copyWith(void Function(ListTagsResponse) updates) =>
      super.copyWith((message) => updates(message as ListTagsResponse))
          as ListTagsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsResponse create() => ListTagsResponse._();
  @$core.override
  ListTagsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsResponse>(create);
  static ListTagsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(301974828)
  $pb.PbList<ResourceTag> get resourcetaglist => $_getList(1);
}

class ListTrailsRequest extends $pb.GeneratedMessage {
  factory ListTrailsRequest({
    $core.String? nexttoken,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListTrailsRequest._();

  factory ListTrailsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrailsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrailsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrailsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrailsRequest copyWith(void Function(ListTrailsRequest) updates) =>
      super.copyWith((message) => updates(message as ListTrailsRequest))
          as ListTrailsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrailsRequest create() => ListTrailsRequest._();
  @$core.override
  ListTrailsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrailsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrailsRequest>(create);
  static ListTrailsRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListTrailsResponse extends $pb.GeneratedMessage {
  factory ListTrailsResponse({
    $core.Iterable<TrailInfo>? trails,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (trails != null) result.trails.addAll(trails);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListTrailsResponse._();

  factory ListTrailsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrailsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrailsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<TrailInfo>(16939441, _omitFieldNames ? '' : 'trails',
        subBuilder: TrailInfo.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrailsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrailsResponse copyWith(void Function(ListTrailsResponse) updates) =>
      super.copyWith((message) => updates(message as ListTrailsResponse))
          as ListTrailsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrailsResponse create() => ListTrailsResponse._();
  @$core.override
  ListTrailsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrailsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrailsResponse>(create);
  static ListTrailsResponse? _defaultInstance;

  @$pb.TagNumber(16939441)
  $pb.PbList<TrailInfo> get trails => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class LookupAttribute extends $pb.GeneratedMessage {
  factory LookupAttribute({
    $core.String? attributevalue,
    LookupAttributeKey? attributekey,
  }) {
    final result = create();
    if (attributevalue != null) result.attributevalue = attributevalue;
    if (attributekey != null) result.attributekey = attributekey;
    return result;
  }

  LookupAttribute._();

  factory LookupAttribute.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LookupAttribute.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LookupAttribute',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(96769221, _omitFieldNames ? '' : 'attributevalue')
    ..aE<LookupAttributeKey>(104472119, _omitFieldNames ? '' : 'attributekey',
        enumValues: LookupAttributeKey.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupAttribute clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupAttribute copyWith(void Function(LookupAttribute) updates) =>
      super.copyWith((message) => updates(message as LookupAttribute))
          as LookupAttribute;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LookupAttribute create() => LookupAttribute._();
  @$core.override
  LookupAttribute createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LookupAttribute getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LookupAttribute>(create);
  static LookupAttribute? _defaultInstance;

  @$pb.TagNumber(96769221)
  $core.String get attributevalue => $_getSZ(0);
  @$pb.TagNumber(96769221)
  set attributevalue($core.String value) => $_setString(0, value);
  @$pb.TagNumber(96769221)
  $core.bool hasAttributevalue() => $_has(0);
  @$pb.TagNumber(96769221)
  void clearAttributevalue() => $_clearField(96769221);

  @$pb.TagNumber(104472119)
  LookupAttributeKey get attributekey => $_getN(1);
  @$pb.TagNumber(104472119)
  set attributekey(LookupAttributeKey value) => $_setField(104472119, value);
  @$pb.TagNumber(104472119)
  $core.bool hasAttributekey() => $_has(1);
  @$pb.TagNumber(104472119)
  void clearAttributekey() => $_clearField(104472119);
}

class LookupEventsRequest extends $pb.GeneratedMessage {
  factory LookupEventsRequest({
    $core.String? endtime,
    $core.Iterable<LookupAttribute>? lookupattributes,
    EventCategory? eventcategory,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? starttime,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (lookupattributes != null)
      result.lookupattributes.addAll(lookupattributes);
    if (eventcategory != null) result.eventcategory = eventcategory;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (starttime != null) result.starttime = starttime;
    return result;
  }

  LookupEventsRequest._();

  factory LookupEventsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LookupEventsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LookupEventsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..pPM<LookupAttribute>(162162567, _omitFieldNames ? '' : 'lookupattributes',
        subBuilder: LookupAttribute.create)
    ..aE<EventCategory>(164668724, _omitFieldNames ? '' : 'eventcategory',
        enumValues: EventCategory.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupEventsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupEventsRequest copyWith(void Function(LookupEventsRequest) updates) =>
      super.copyWith((message) => updates(message as LookupEventsRequest))
          as LookupEventsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LookupEventsRequest create() => LookupEventsRequest._();
  @$core.override
  LookupEventsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LookupEventsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LookupEventsRequest>(create);
  static LookupEventsRequest? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(162162567)
  $pb.PbList<LookupAttribute> get lookupattributes => $_getList(1);

  @$pb.TagNumber(164668724)
  EventCategory get eventcategory => $_getN(2);
  @$pb.TagNumber(164668724)
  set eventcategory(EventCategory value) => $_setField(164668724, value);
  @$pb.TagNumber(164668724)
  $core.bool hasEventcategory() => $_has(2);
  @$pb.TagNumber(164668724)
  void clearEventcategory() => $_clearField(164668724);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(3);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(3);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(4);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(4);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(5);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(5);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);
}

class LookupEventsResponse extends $pb.GeneratedMessage {
  factory LookupEventsResponse({
    $core.Iterable<Event>? events,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (events != null) result.events.addAll(events);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  LookupEventsResponse._();

  factory LookupEventsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LookupEventsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LookupEventsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<Event>(3416229, _omitFieldNames ? '' : 'events',
        subBuilder: Event.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupEventsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupEventsResponse copyWith(void Function(LookupEventsResponse) updates) =>
      super.copyWith((message) => updates(message as LookupEventsResponse))
          as LookupEventsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LookupEventsResponse create() => LookupEventsResponse._();
  @$core.override
  LookupEventsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LookupEventsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LookupEventsResponse>(create);
  static LookupEventsResponse? _defaultInstance;

  @$pb.TagNumber(3416229)
  $pb.PbList<Event> get events => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class MaxConcurrentQueriesException extends $pb.GeneratedMessage {
  factory MaxConcurrentQueriesException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  MaxConcurrentQueriesException._();

  factory MaxConcurrentQueriesException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MaxConcurrentQueriesException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MaxConcurrentQueriesException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MaxConcurrentQueriesException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MaxConcurrentQueriesException copyWith(
          void Function(MaxConcurrentQueriesException) updates) =>
      super.copyWith(
              (message) => updates(message as MaxConcurrentQueriesException))
          as MaxConcurrentQueriesException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MaxConcurrentQueriesException create() =>
      MaxConcurrentQueriesException._();
  @$core.override
  MaxConcurrentQueriesException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MaxConcurrentQueriesException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MaxConcurrentQueriesException>(create);
  static MaxConcurrentQueriesException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class MaximumNumberOfTrailsExceededException extends $pb.GeneratedMessage {
  factory MaximumNumberOfTrailsExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  MaximumNumberOfTrailsExceededException._();

  factory MaximumNumberOfTrailsExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MaximumNumberOfTrailsExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MaximumNumberOfTrailsExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MaximumNumberOfTrailsExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MaximumNumberOfTrailsExceededException copyWith(
          void Function(MaximumNumberOfTrailsExceededException) updates) =>
      super.copyWith((message) =>
              updates(message as MaximumNumberOfTrailsExceededException))
          as MaximumNumberOfTrailsExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MaximumNumberOfTrailsExceededException create() =>
      MaximumNumberOfTrailsExceededException._();
  @$core.override
  MaximumNumberOfTrailsExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MaximumNumberOfTrailsExceededException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          MaximumNumberOfTrailsExceededException>(create);
  static MaximumNumberOfTrailsExceededException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class NoManagementAccountSLRExistsException extends $pb.GeneratedMessage {
  factory NoManagementAccountSLRExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoManagementAccountSLRExistsException._();

  factory NoManagementAccountSLRExistsException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoManagementAccountSLRExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoManagementAccountSLRExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoManagementAccountSLRExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoManagementAccountSLRExistsException copyWith(
          void Function(NoManagementAccountSLRExistsException) updates) =>
      super.copyWith((message) =>
              updates(message as NoManagementAccountSLRExistsException))
          as NoManagementAccountSLRExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoManagementAccountSLRExistsException create() =>
      NoManagementAccountSLRExistsException._();
  @$core.override
  NoManagementAccountSLRExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoManagementAccountSLRExistsException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          NoManagementAccountSLRExistsException>(create);
  static NoManagementAccountSLRExistsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class NotOrganizationManagementAccountException extends $pb.GeneratedMessage {
  factory NotOrganizationManagementAccountException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NotOrganizationManagementAccountException._();

  factory NotOrganizationManagementAccountException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotOrganizationManagementAccountException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotOrganizationManagementAccountException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotOrganizationManagementAccountException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotOrganizationManagementAccountException copyWith(
          void Function(NotOrganizationManagementAccountException) updates) =>
      super.copyWith((message) =>
              updates(message as NotOrganizationManagementAccountException))
          as NotOrganizationManagementAccountException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotOrganizationManagementAccountException create() =>
      NotOrganizationManagementAccountException._();
  @$core.override
  NotOrganizationManagementAccountException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotOrganizationManagementAccountException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          NotOrganizationManagementAccountException>(create);
  static NotOrganizationManagementAccountException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class NotOrganizationMasterAccountException extends $pb.GeneratedMessage {
  factory NotOrganizationMasterAccountException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NotOrganizationMasterAccountException._();

  factory NotOrganizationMasterAccountException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotOrganizationMasterAccountException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotOrganizationMasterAccountException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotOrganizationMasterAccountException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotOrganizationMasterAccountException copyWith(
          void Function(NotOrganizationMasterAccountException) updates) =>
      super.copyWith((message) =>
              updates(message as NotOrganizationMasterAccountException))
          as NotOrganizationMasterAccountException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotOrganizationMasterAccountException create() =>
      NotOrganizationMasterAccountException._();
  @$core.override
  NotOrganizationMasterAccountException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotOrganizationMasterAccountException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          NotOrganizationMasterAccountException>(create);
  static NotOrganizationMasterAccountException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class OperationNotPermittedException extends $pb.GeneratedMessage {
  factory OperationNotPermittedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  OperationNotPermittedException._();

  factory OperationNotPermittedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OperationNotPermittedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OperationNotPermittedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OperationNotPermittedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OperationNotPermittedException copyWith(
          void Function(OperationNotPermittedException) updates) =>
      super.copyWith(
              (message) => updates(message as OperationNotPermittedException))
          as OperationNotPermittedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OperationNotPermittedException create() =>
      OperationNotPermittedException._();
  @$core.override
  OperationNotPermittedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OperationNotPermittedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OperationNotPermittedException>(create);
  static OperationNotPermittedException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class OrganizationNotInAllFeaturesModeException extends $pb.GeneratedMessage {
  factory OrganizationNotInAllFeaturesModeException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  OrganizationNotInAllFeaturesModeException._();

  factory OrganizationNotInAllFeaturesModeException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OrganizationNotInAllFeaturesModeException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OrganizationNotInAllFeaturesModeException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OrganizationNotInAllFeaturesModeException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OrganizationNotInAllFeaturesModeException copyWith(
          void Function(OrganizationNotInAllFeaturesModeException) updates) =>
      super.copyWith((message) =>
              updates(message as OrganizationNotInAllFeaturesModeException))
          as OrganizationNotInAllFeaturesModeException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OrganizationNotInAllFeaturesModeException create() =>
      OrganizationNotInAllFeaturesModeException._();
  @$core.override
  OrganizationNotInAllFeaturesModeException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OrganizationNotInAllFeaturesModeException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          OrganizationNotInAllFeaturesModeException>(create);
  static OrganizationNotInAllFeaturesModeException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class OrganizationsNotInUseException extends $pb.GeneratedMessage {
  factory OrganizationsNotInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  OrganizationsNotInUseException._();

  factory OrganizationsNotInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OrganizationsNotInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OrganizationsNotInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OrganizationsNotInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OrganizationsNotInUseException copyWith(
          void Function(OrganizationsNotInUseException) updates) =>
      super.copyWith(
              (message) => updates(message as OrganizationsNotInUseException))
          as OrganizationsNotInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OrganizationsNotInUseException create() =>
      OrganizationsNotInUseException._();
  @$core.override
  OrganizationsNotInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OrganizationsNotInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OrganizationsNotInUseException>(create);
  static OrganizationsNotInUseException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class PartitionKey extends $pb.GeneratedMessage {
  factory PartitionKey({
    $core.String? name,
    $core.String? type,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    return result;
  }

  PartitionKey._();

  factory PartitionKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PartitionKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PartitionKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartitionKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartitionKey copyWith(void Function(PartitionKey) updates) =>
      super.copyWith((message) => updates(message as PartitionKey))
          as PartitionKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PartitionKey create() => PartitionKey._();
  @$core.override
  PartitionKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PartitionKey getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PartitionKey>(create);
  static PartitionKey? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(1);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(1, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class PublicKey extends $pb.GeneratedMessage {
  factory PublicKey({
    $core.String? validityendtime,
    $core.String? validitystarttime,
    $core.List<$core.int>? value,
    $core.String? fingerprint,
  }) {
    final result = create();
    if (validityendtime != null) result.validityendtime = validityendtime;
    if (validitystarttime != null) result.validitystarttime = validitystarttime;
    if (value != null) result.value = value;
    if (fingerprint != null) result.fingerprint = fingerprint;
    return result;
  }

  PublicKey._();

  factory PublicKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublicKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublicKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(192996550, _omitFieldNames ? '' : 'validityendtime')
    ..aOS(288335101, _omitFieldNames ? '' : 'validitystarttime')
    ..a<$core.List<$core.int>>(
        289929579, _omitFieldNames ? '' : 'value', $pb.PbFieldType.OY)
    ..aOS(502547484, _omitFieldNames ? '' : 'fingerprint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublicKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublicKey copyWith(void Function(PublicKey) updates) =>
      super.copyWith((message) => updates(message as PublicKey)) as PublicKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublicKey create() => PublicKey._();
  @$core.override
  PublicKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublicKey getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PublicKey>(create);
  static PublicKey? _defaultInstance;

  @$pb.TagNumber(192996550)
  $core.String get validityendtime => $_getSZ(0);
  @$pb.TagNumber(192996550)
  set validityendtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(192996550)
  $core.bool hasValidityendtime() => $_has(0);
  @$pb.TagNumber(192996550)
  void clearValidityendtime() => $_clearField(192996550);

  @$pb.TagNumber(288335101)
  $core.String get validitystarttime => $_getSZ(1);
  @$pb.TagNumber(288335101)
  set validitystarttime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(288335101)
  $core.bool hasValiditystarttime() => $_has(1);
  @$pb.TagNumber(288335101)
  void clearValiditystarttime() => $_clearField(288335101);

  @$pb.TagNumber(289929579)
  $core.List<$core.int> get value => $_getN(2);
  @$pb.TagNumber(289929579)
  set value($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(2);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(502547484)
  $core.String get fingerprint => $_getSZ(3);
  @$pb.TagNumber(502547484)
  set fingerprint($core.String value) => $_setString(3, value);
  @$pb.TagNumber(502547484)
  $core.bool hasFingerprint() => $_has(3);
  @$pb.TagNumber(502547484)
  void clearFingerprint() => $_clearField(502547484);
}

class PutEventConfigurationRequest extends $pb.GeneratedMessage {
  factory PutEventConfigurationRequest({
    $core.String? eventdatastore,
    $core.Iterable<ContextKeySelector>? contextkeyselectors,
    $core.Iterable<AggregationConfiguration>? aggregationconfigurations,
    $core.String? trailname,
    MaxEventSize? maxeventsize,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (contextkeyselectors != null)
      result.contextkeyselectors.addAll(contextkeyselectors);
    if (aggregationconfigurations != null)
      result.aggregationconfigurations.addAll(aggregationconfigurations);
    if (trailname != null) result.trailname = trailname;
    if (maxeventsize != null) result.maxeventsize = maxeventsize;
    return result;
  }

  PutEventConfigurationRequest._();

  factory PutEventConfigurationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventConfigurationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventConfigurationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..pPM<ContextKeySelector>(
        342040300, _omitFieldNames ? '' : 'contextkeyselectors',
        subBuilder: ContextKeySelector.create)
    ..pPM<AggregationConfiguration>(
        481530463, _omitFieldNames ? '' : 'aggregationconfigurations',
        subBuilder: AggregationConfiguration.create)
    ..aOS(507774985, _omitFieldNames ? '' : 'trailname')
    ..aE<MaxEventSize>(517627763, _omitFieldNames ? '' : 'maxeventsize',
        enumValues: MaxEventSize.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventConfigurationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventConfigurationRequest copyWith(
          void Function(PutEventConfigurationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PutEventConfigurationRequest))
          as PutEventConfigurationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventConfigurationRequest create() =>
      PutEventConfigurationRequest._();
  @$core.override
  PutEventConfigurationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventConfigurationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventConfigurationRequest>(create);
  static PutEventConfigurationRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(342040300)
  $pb.PbList<ContextKeySelector> get contextkeyselectors => $_getList(1);

  @$pb.TagNumber(481530463)
  $pb.PbList<AggregationConfiguration> get aggregationconfigurations =>
      $_getList(2);

  @$pb.TagNumber(507774985)
  $core.String get trailname => $_getSZ(3);
  @$pb.TagNumber(507774985)
  set trailname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(507774985)
  $core.bool hasTrailname() => $_has(3);
  @$pb.TagNumber(507774985)
  void clearTrailname() => $_clearField(507774985);

  @$pb.TagNumber(517627763)
  MaxEventSize get maxeventsize => $_getN(4);
  @$pb.TagNumber(517627763)
  set maxeventsize(MaxEventSize value) => $_setField(517627763, value);
  @$pb.TagNumber(517627763)
  $core.bool hasMaxeventsize() => $_has(4);
  @$pb.TagNumber(517627763)
  void clearMaxeventsize() => $_clearField(517627763);
}

class PutEventConfigurationResponse extends $pb.GeneratedMessage {
  factory PutEventConfigurationResponse({
    $core.String? trailarn,
    $core.String? eventdatastorearn,
    $core.Iterable<ContextKeySelector>? contextkeyselectors,
    $core.Iterable<AggregationConfiguration>? aggregationconfigurations,
    MaxEventSize? maxeventsize,
  }) {
    final result = create();
    if (trailarn != null) result.trailarn = trailarn;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (contextkeyselectors != null)
      result.contextkeyselectors.addAll(contextkeyselectors);
    if (aggregationconfigurations != null)
      result.aggregationconfigurations.addAll(aggregationconfigurations);
    if (maxeventsize != null) result.maxeventsize = maxeventsize;
    return result;
  }

  PutEventConfigurationResponse._();

  factory PutEventConfigurationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventConfigurationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventConfigurationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..pPM<ContextKeySelector>(
        342040300, _omitFieldNames ? '' : 'contextkeyselectors',
        subBuilder: ContextKeySelector.create)
    ..pPM<AggregationConfiguration>(
        481530463, _omitFieldNames ? '' : 'aggregationconfigurations',
        subBuilder: AggregationConfiguration.create)
    ..aE<MaxEventSize>(517627763, _omitFieldNames ? '' : 'maxeventsize',
        enumValues: MaxEventSize.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventConfigurationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventConfigurationResponse copyWith(
          void Function(PutEventConfigurationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PutEventConfigurationResponse))
          as PutEventConfigurationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventConfigurationResponse create() =>
      PutEventConfigurationResponse._();
  @$core.override
  PutEventConfigurationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventConfigurationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventConfigurationResponse>(create);
  static PutEventConfigurationResponse? _defaultInstance;

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(0);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(0);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(1);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(1);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(342040300)
  $pb.PbList<ContextKeySelector> get contextkeyselectors => $_getList(2);

  @$pb.TagNumber(481530463)
  $pb.PbList<AggregationConfiguration> get aggregationconfigurations =>
      $_getList(3);

  @$pb.TagNumber(517627763)
  MaxEventSize get maxeventsize => $_getN(4);
  @$pb.TagNumber(517627763)
  set maxeventsize(MaxEventSize value) => $_setField(517627763, value);
  @$pb.TagNumber(517627763)
  $core.bool hasMaxeventsize() => $_has(4);
  @$pb.TagNumber(517627763)
  void clearMaxeventsize() => $_clearField(517627763);
}

class PutEventSelectorsRequest extends $pb.GeneratedMessage {
  factory PutEventSelectorsRequest({
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.Iterable<EventSelector>? eventselectors,
    $core.String? trailname,
  }) {
    final result = create();
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (eventselectors != null) result.eventselectors.addAll(eventselectors);
    if (trailname != null) result.trailname = trailname;
    return result;
  }

  PutEventSelectorsRequest._();

  factory PutEventSelectorsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventSelectorsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventSelectorsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..pPM<EventSelector>(344575328, _omitFieldNames ? '' : 'eventselectors',
        subBuilder: EventSelector.create)
    ..aOS(507774985, _omitFieldNames ? '' : 'trailname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventSelectorsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventSelectorsRequest copyWith(
          void Function(PutEventSelectorsRequest) updates) =>
      super.copyWith((message) => updates(message as PutEventSelectorsRequest))
          as PutEventSelectorsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventSelectorsRequest create() => PutEventSelectorsRequest._();
  @$core.override
  PutEventSelectorsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventSelectorsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventSelectorsRequest>(create);
  static PutEventSelectorsRequest? _defaultInstance;

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(0);

  @$pb.TagNumber(344575328)
  $pb.PbList<EventSelector> get eventselectors => $_getList(1);

  @$pb.TagNumber(507774985)
  $core.String get trailname => $_getSZ(2);
  @$pb.TagNumber(507774985)
  set trailname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(507774985)
  $core.bool hasTrailname() => $_has(2);
  @$pb.TagNumber(507774985)
  void clearTrailname() => $_clearField(507774985);
}

class PutEventSelectorsResponse extends $pb.GeneratedMessage {
  factory PutEventSelectorsResponse({
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? trailarn,
    $core.Iterable<EventSelector>? eventselectors,
  }) {
    final result = create();
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (trailarn != null) result.trailarn = trailarn;
    if (eventselectors != null) result.eventselectors.addAll(eventselectors);
    return result;
  }

  PutEventSelectorsResponse._();

  factory PutEventSelectorsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventSelectorsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventSelectorsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..pPM<EventSelector>(344575328, _omitFieldNames ? '' : 'eventselectors',
        subBuilder: EventSelector.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventSelectorsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventSelectorsResponse copyWith(
          void Function(PutEventSelectorsResponse) updates) =>
      super.copyWith((message) => updates(message as PutEventSelectorsResponse))
          as PutEventSelectorsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventSelectorsResponse create() => PutEventSelectorsResponse._();
  @$core.override
  PutEventSelectorsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventSelectorsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventSelectorsResponse>(create);
  static PutEventSelectorsResponse? _defaultInstance;

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(0);

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(1);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(1);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(344575328)
  $pb.PbList<EventSelector> get eventselectors => $_getList(2);
}

class PutInsightSelectorsRequest extends $pb.GeneratedMessage {
  factory PutInsightSelectorsRequest({
    $core.Iterable<InsightSelector>? insightselectors,
    $core.String? eventdatastore,
    $core.String? insightsdestination,
    $core.String? trailname,
  }) {
    final result = create();
    if (insightselectors != null)
      result.insightselectors.addAll(insightselectors);
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (insightsdestination != null)
      result.insightsdestination = insightsdestination;
    if (trailname != null) result.trailname = trailname;
    return result;
  }

  PutInsightSelectorsRequest._();

  factory PutInsightSelectorsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutInsightSelectorsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutInsightSelectorsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<InsightSelector>(49628748, _omitFieldNames ? '' : 'insightselectors',
        subBuilder: InsightSelector.create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aOS(396776649, _omitFieldNames ? '' : 'insightsdestination')
    ..aOS(507774985, _omitFieldNames ? '' : 'trailname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightSelectorsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightSelectorsRequest copyWith(
          void Function(PutInsightSelectorsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PutInsightSelectorsRequest))
          as PutInsightSelectorsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutInsightSelectorsRequest create() => PutInsightSelectorsRequest._();
  @$core.override
  PutInsightSelectorsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutInsightSelectorsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutInsightSelectorsRequest>(create);
  static PutInsightSelectorsRequest? _defaultInstance;

  @$pb.TagNumber(49628748)
  $pb.PbList<InsightSelector> get insightselectors => $_getList(0);

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(1);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(1, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(1);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(396776649)
  $core.String get insightsdestination => $_getSZ(2);
  @$pb.TagNumber(396776649)
  set insightsdestination($core.String value) => $_setString(2, value);
  @$pb.TagNumber(396776649)
  $core.bool hasInsightsdestination() => $_has(2);
  @$pb.TagNumber(396776649)
  void clearInsightsdestination() => $_clearField(396776649);

  @$pb.TagNumber(507774985)
  $core.String get trailname => $_getSZ(3);
  @$pb.TagNumber(507774985)
  set trailname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(507774985)
  $core.bool hasTrailname() => $_has(3);
  @$pb.TagNumber(507774985)
  void clearTrailname() => $_clearField(507774985);
}

class PutInsightSelectorsResponse extends $pb.GeneratedMessage {
  factory PutInsightSelectorsResponse({
    $core.String? trailarn,
    $core.Iterable<InsightSelector>? insightselectors,
    $core.String? eventdatastorearn,
    $core.String? insightsdestination,
  }) {
    final result = create();
    if (trailarn != null) result.trailarn = trailarn;
    if (insightselectors != null)
      result.insightselectors.addAll(insightselectors);
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (insightsdestination != null)
      result.insightsdestination = insightsdestination;
    return result;
  }

  PutInsightSelectorsResponse._();

  factory PutInsightSelectorsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutInsightSelectorsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutInsightSelectorsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..pPM<InsightSelector>(49628748, _omitFieldNames ? '' : 'insightselectors',
        subBuilder: InsightSelector.create)
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(396776649, _omitFieldNames ? '' : 'insightsdestination')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightSelectorsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightSelectorsResponse copyWith(
          void Function(PutInsightSelectorsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PutInsightSelectorsResponse))
          as PutInsightSelectorsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutInsightSelectorsResponse create() =>
      PutInsightSelectorsResponse._();
  @$core.override
  PutInsightSelectorsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutInsightSelectorsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutInsightSelectorsResponse>(create);
  static PutInsightSelectorsResponse? _defaultInstance;

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(0);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(0);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(49628748)
  $pb.PbList<InsightSelector> get insightselectors => $_getList(1);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(2);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(2);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(396776649)
  $core.String get insightsdestination => $_getSZ(3);
  @$pb.TagNumber(396776649)
  set insightsdestination($core.String value) => $_setString(3, value);
  @$pb.TagNumber(396776649)
  $core.bool hasInsightsdestination() => $_has(3);
  @$pb.TagNumber(396776649)
  void clearInsightsdestination() => $_clearField(396776649);
}

class PutResourcePolicyRequest extends $pb.GeneratedMessage {
  factory PutResourcePolicyRequest({
    $core.String? resourcepolicy,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcepolicy != null) result.resourcepolicy = resourcepolicy;
    if (resourcearn != null) result.resourcearn = resourcearn;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(15747632, _omitFieldNames ? '' : 'resourcepolicy')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class PutResourcePolicyResponse extends $pb.GeneratedMessage {
  factory PutResourcePolicyResponse({
    $core.String? resourcepolicy,
    $core.String? resourcearn,
    $core.String? delegatedadminresourcepolicy,
  }) {
    final result = create();
    if (resourcepolicy != null) result.resourcepolicy = resourcepolicy;
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (delegatedadminresourcepolicy != null)
      result.delegatedadminresourcepolicy = delegatedadminresourcepolicy;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(15747632, _omitFieldNames ? '' : 'resourcepolicy')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(510966132, _omitFieldNames ? '' : 'delegatedadminresourcepolicy')
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

  @$pb.TagNumber(15747632)
  $core.String get resourcepolicy => $_getSZ(0);
  @$pb.TagNumber(15747632)
  set resourcepolicy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(15747632)
  $core.bool hasResourcepolicy() => $_has(0);
  @$pb.TagNumber(15747632)
  void clearResourcepolicy() => $_clearField(15747632);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(510966132)
  $core.String get delegatedadminresourcepolicy => $_getSZ(2);
  @$pb.TagNumber(510966132)
  set delegatedadminresourcepolicy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(510966132)
  $core.bool hasDelegatedadminresourcepolicy() => $_has(2);
  @$pb.TagNumber(510966132)
  void clearDelegatedadminresourcepolicy() => $_clearField(510966132);
}

class Query extends $pb.GeneratedMessage {
  factory Query({
    $core.String? creationtime,
    $core.String? queryid,
    QueryStatus? querystatus,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (queryid != null) result.queryid = queryid;
    if (querystatus != null) result.querystatus = querystatus;
    return result;
  }

  Query._();

  factory Query.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Query.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Query',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..aE<QueryStatus>(367016406, _omitFieldNames ? '' : 'querystatus',
        enumValues: QueryStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Query clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Query copyWith(void Function(Query) updates) =>
      super.copyWith((message) => updates(message as Query)) as Query;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Query create() => Query._();
  @$core.override
  Query createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Query getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Query>(create);
  static Query? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(1);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(1);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(367016406)
  QueryStatus get querystatus => $_getN(2);
  @$pb.TagNumber(367016406)
  set querystatus(QueryStatus value) => $_setField(367016406, value);
  @$pb.TagNumber(367016406)
  $core.bool hasQuerystatus() => $_has(2);
  @$pb.TagNumber(367016406)
  void clearQuerystatus() => $_clearField(367016406);
}

class QueryIdNotFoundException extends $pb.GeneratedMessage {
  factory QueryIdNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  QueryIdNotFoundException._();

  factory QueryIdNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryIdNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryIdNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryIdNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryIdNotFoundException copyWith(
          void Function(QueryIdNotFoundException) updates) =>
      super.copyWith((message) => updates(message as QueryIdNotFoundException))
          as QueryIdNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryIdNotFoundException create() => QueryIdNotFoundException._();
  @$core.override
  QueryIdNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryIdNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryIdNotFoundException>(create);
  static QueryIdNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class QueryStatistics extends $pb.GeneratedMessage {
  factory QueryStatistics({
    $core.int? totalresultscount,
    $fixnum.Int64? bytesscanned,
    $core.int? resultscount,
  }) {
    final result = create();
    if (totalresultscount != null) result.totalresultscount = totalresultscount;
    if (bytesscanned != null) result.bytesscanned = bytesscanned;
    if (resultscount != null) result.resultscount = resultscount;
    return result;
  }

  QueryStatistics._();

  factory QueryStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aI(167126221, _omitFieldNames ? '' : 'totalresultscount')
    ..aInt64(186950329, _omitFieldNames ? '' : 'bytesscanned')
    ..aI(530057669, _omitFieldNames ? '' : 'resultscount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStatistics copyWith(void Function(QueryStatistics) updates) =>
      super.copyWith((message) => updates(message as QueryStatistics))
          as QueryStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryStatistics create() => QueryStatistics._();
  @$core.override
  QueryStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryStatistics>(create);
  static QueryStatistics? _defaultInstance;

  @$pb.TagNumber(167126221)
  $core.int get totalresultscount => $_getIZ(0);
  @$pb.TagNumber(167126221)
  set totalresultscount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(167126221)
  $core.bool hasTotalresultscount() => $_has(0);
  @$pb.TagNumber(167126221)
  void clearTotalresultscount() => $_clearField(167126221);

  @$pb.TagNumber(186950329)
  $fixnum.Int64 get bytesscanned => $_getI64(1);
  @$pb.TagNumber(186950329)
  set bytesscanned($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(186950329)
  $core.bool hasBytesscanned() => $_has(1);
  @$pb.TagNumber(186950329)
  void clearBytesscanned() => $_clearField(186950329);

  @$pb.TagNumber(530057669)
  $core.int get resultscount => $_getIZ(2);
  @$pb.TagNumber(530057669)
  set resultscount($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(530057669)
  $core.bool hasResultscount() => $_has(2);
  @$pb.TagNumber(530057669)
  void clearResultscount() => $_clearField(530057669);
}

class QueryStatisticsForDescribeQuery extends $pb.GeneratedMessage {
  factory QueryStatisticsForDescribeQuery({
    $fixnum.Int64? eventsmatched,
    $core.String? creationtime,
    $fixnum.Int64? bytesscanned,
    $core.int? executiontimeinmillis,
    $fixnum.Int64? eventsscanned,
  }) {
    final result = create();
    if (eventsmatched != null) result.eventsmatched = eventsmatched;
    if (creationtime != null) result.creationtime = creationtime;
    if (bytesscanned != null) result.bytesscanned = bytesscanned;
    if (executiontimeinmillis != null)
      result.executiontimeinmillis = executiontimeinmillis;
    if (eventsscanned != null) result.eventsscanned = eventsscanned;
    return result;
  }

  QueryStatisticsForDescribeQuery._();

  factory QueryStatisticsForDescribeQuery.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryStatisticsForDescribeQuery.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryStatisticsForDescribeQuery',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aInt64(19501829, _omitFieldNames ? '' : 'eventsmatched')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aInt64(186950329, _omitFieldNames ? '' : 'bytesscanned')
    ..aI(447596178, _omitFieldNames ? '' : 'executiontimeinmillis')
    ..aInt64(475241657, _omitFieldNames ? '' : 'eventsscanned')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStatisticsForDescribeQuery clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStatisticsForDescribeQuery copyWith(
          void Function(QueryStatisticsForDescribeQuery) updates) =>
      super.copyWith(
              (message) => updates(message as QueryStatisticsForDescribeQuery))
          as QueryStatisticsForDescribeQuery;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryStatisticsForDescribeQuery create() =>
      QueryStatisticsForDescribeQuery._();
  @$core.override
  QueryStatisticsForDescribeQuery createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryStatisticsForDescribeQuery getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryStatisticsForDescribeQuery>(
          create);
  static QueryStatisticsForDescribeQuery? _defaultInstance;

  @$pb.TagNumber(19501829)
  $fixnum.Int64 get eventsmatched => $_getI64(0);
  @$pb.TagNumber(19501829)
  set eventsmatched($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(19501829)
  $core.bool hasEventsmatched() => $_has(0);
  @$pb.TagNumber(19501829)
  void clearEventsmatched() => $_clearField(19501829);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(186950329)
  $fixnum.Int64 get bytesscanned => $_getI64(2);
  @$pb.TagNumber(186950329)
  set bytesscanned($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(186950329)
  $core.bool hasBytesscanned() => $_has(2);
  @$pb.TagNumber(186950329)
  void clearBytesscanned() => $_clearField(186950329);

  @$pb.TagNumber(447596178)
  $core.int get executiontimeinmillis => $_getIZ(3);
  @$pb.TagNumber(447596178)
  set executiontimeinmillis($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(447596178)
  $core.bool hasExecutiontimeinmillis() => $_has(3);
  @$pb.TagNumber(447596178)
  void clearExecutiontimeinmillis() => $_clearField(447596178);

  @$pb.TagNumber(475241657)
  $fixnum.Int64 get eventsscanned => $_getI64(4);
  @$pb.TagNumber(475241657)
  set eventsscanned($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(475241657)
  $core.bool hasEventsscanned() => $_has(4);
  @$pb.TagNumber(475241657)
  void clearEventsscanned() => $_clearField(475241657);
}

class RefreshSchedule extends $pb.GeneratedMessage {
  factory RefreshSchedule({
    RefreshScheduleStatus? status,
    $core.String? timeofday,
    RefreshScheduleFrequency? frequency,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (timeofday != null) result.timeofday = timeofday;
    if (frequency != null) result.frequency = frequency;
    return result;
  }

  RefreshSchedule._();

  factory RefreshSchedule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RefreshSchedule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RefreshSchedule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<RefreshScheduleStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: RefreshScheduleStatus.values)
    ..aOS(77605358, _omitFieldNames ? '' : 'timeofday')
    ..aOM<RefreshScheduleFrequency>(
        227673762, _omitFieldNames ? '' : 'frequency',
        subBuilder: RefreshScheduleFrequency.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RefreshSchedule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RefreshSchedule copyWith(void Function(RefreshSchedule) updates) =>
      super.copyWith((message) => updates(message as RefreshSchedule))
          as RefreshSchedule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RefreshSchedule create() => RefreshSchedule._();
  @$core.override
  RefreshSchedule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RefreshSchedule getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RefreshSchedule>(create);
  static RefreshSchedule? _defaultInstance;

  @$pb.TagNumber(6222352)
  RefreshScheduleStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(RefreshScheduleStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(77605358)
  $core.String get timeofday => $_getSZ(1);
  @$pb.TagNumber(77605358)
  set timeofday($core.String value) => $_setString(1, value);
  @$pb.TagNumber(77605358)
  $core.bool hasTimeofday() => $_has(1);
  @$pb.TagNumber(77605358)
  void clearTimeofday() => $_clearField(77605358);

  @$pb.TagNumber(227673762)
  RefreshScheduleFrequency get frequency => $_getN(2);
  @$pb.TagNumber(227673762)
  set frequency(RefreshScheduleFrequency value) => $_setField(227673762, value);
  @$pb.TagNumber(227673762)
  $core.bool hasFrequency() => $_has(2);
  @$pb.TagNumber(227673762)
  void clearFrequency() => $_clearField(227673762);
  @$pb.TagNumber(227673762)
  RefreshScheduleFrequency ensureFrequency() => $_ensure(2);
}

class RefreshScheduleFrequency extends $pb.GeneratedMessage {
  factory RefreshScheduleFrequency({
    RefreshScheduleFrequencyUnit? unit,
    $core.int? value,
  }) {
    final result = create();
    if (unit != null) result.unit = unit;
    if (value != null) result.value = value;
    return result;
  }

  RefreshScheduleFrequency._();

  factory RefreshScheduleFrequency.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RefreshScheduleFrequency.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RefreshScheduleFrequency',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<RefreshScheduleFrequencyUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: RefreshScheduleFrequencyUnit.values)
    ..aI(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RefreshScheduleFrequency clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RefreshScheduleFrequency copyWith(
          void Function(RefreshScheduleFrequency) updates) =>
      super.copyWith((message) => updates(message as RefreshScheduleFrequency))
          as RefreshScheduleFrequency;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RefreshScheduleFrequency create() => RefreshScheduleFrequency._();
  @$core.override
  RefreshScheduleFrequency createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RefreshScheduleFrequency getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RefreshScheduleFrequency>(create);
  static RefreshScheduleFrequency? _defaultInstance;

  @$pb.TagNumber(148989480)
  RefreshScheduleFrequencyUnit get unit => $_getN(0);
  @$pb.TagNumber(148989480)
  set unit(RefreshScheduleFrequencyUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(0);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(289929579)
  $core.int get value => $_getIZ(1);
  @$pb.TagNumber(289929579)
  set value($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class RegisterOrganizationDelegatedAdminRequest extends $pb.GeneratedMessage {
  factory RegisterOrganizationDelegatedAdminRequest({
    $core.String? memberaccountid,
  }) {
    final result = create();
    if (memberaccountid != null) result.memberaccountid = memberaccountid;
    return result;
  }

  RegisterOrganizationDelegatedAdminRequest._();

  factory RegisterOrganizationDelegatedAdminRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegisterOrganizationDelegatedAdminRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegisterOrganizationDelegatedAdminRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(374644988, _omitFieldNames ? '' : 'memberaccountid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterOrganizationDelegatedAdminRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterOrganizationDelegatedAdminRequest copyWith(
          void Function(RegisterOrganizationDelegatedAdminRequest) updates) =>
      super.copyWith((message) =>
              updates(message as RegisterOrganizationDelegatedAdminRequest))
          as RegisterOrganizationDelegatedAdminRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegisterOrganizationDelegatedAdminRequest create() =>
      RegisterOrganizationDelegatedAdminRequest._();
  @$core.override
  RegisterOrganizationDelegatedAdminRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegisterOrganizationDelegatedAdminRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          RegisterOrganizationDelegatedAdminRequest>(create);
  static RegisterOrganizationDelegatedAdminRequest? _defaultInstance;

  @$pb.TagNumber(374644988)
  $core.String get memberaccountid => $_getSZ(0);
  @$pb.TagNumber(374644988)
  set memberaccountid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(374644988)
  $core.bool hasMemberaccountid() => $_has(0);
  @$pb.TagNumber(374644988)
  void clearMemberaccountid() => $_clearField(374644988);
}

class RegisterOrganizationDelegatedAdminResponse extends $pb.GeneratedMessage {
  factory RegisterOrganizationDelegatedAdminResponse() => create();

  RegisterOrganizationDelegatedAdminResponse._();

  factory RegisterOrganizationDelegatedAdminResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegisterOrganizationDelegatedAdminResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegisterOrganizationDelegatedAdminResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterOrganizationDelegatedAdminResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegisterOrganizationDelegatedAdminResponse copyWith(
          void Function(RegisterOrganizationDelegatedAdminResponse) updates) =>
      super.copyWith((message) =>
              updates(message as RegisterOrganizationDelegatedAdminResponse))
          as RegisterOrganizationDelegatedAdminResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegisterOrganizationDelegatedAdminResponse create() =>
      RegisterOrganizationDelegatedAdminResponse._();
  @$core.override
  RegisterOrganizationDelegatedAdminResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegisterOrganizationDelegatedAdminResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          RegisterOrganizationDelegatedAdminResponse>(create);
  static RegisterOrganizationDelegatedAdminResponse? _defaultInstance;
}

class RemoveTagsRequest extends $pb.GeneratedMessage {
  factory RemoveTagsRequest({
    $core.Iterable<Tag>? tagslist,
    $core.String? resourceid,
  }) {
    final result = create();
    if (tagslist != null) result.tagslist.addAll(tagslist);
    if (resourceid != null) result.resourceid = resourceid;
    return result;
  }

  RemoveTagsRequest._();

  factory RemoveTagsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveTagsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveTagsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..aOS(526146833, _omitFieldNames ? '' : 'resourceid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsRequest copyWith(void Function(RemoveTagsRequest) updates) =>
      super.copyWith((message) => updates(message as RemoveTagsRequest))
          as RemoveTagsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveTagsRequest create() => RemoveTagsRequest._();
  @$core.override
  RemoveTagsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveTagsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveTagsRequest>(create);
  static RemoveTagsRequest? _defaultInstance;

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(0);

  @$pb.TagNumber(526146833)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(526146833)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(526146833)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(526146833)
  void clearResourceid() => $_clearField(526146833);
}

class RemoveTagsResponse extends $pb.GeneratedMessage {
  factory RemoveTagsResponse() => create();

  RemoveTagsResponse._();

  factory RemoveTagsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveTagsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveTagsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsResponse copyWith(void Function(RemoveTagsResponse) updates) =>
      super.copyWith((message) => updates(message as RemoveTagsResponse))
          as RemoveTagsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveTagsResponse create() => RemoveTagsResponse._();
  @$core.override
  RemoveTagsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveTagsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveTagsResponse>(create);
  static RemoveTagsResponse? _defaultInstance;
}

class RequestWidget extends $pb.GeneratedMessage {
  factory RequestWidget({
    $core.Iterable<$core.String>? queryparameters,
    $core.String? querystatement,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? viewproperties,
  }) {
    final result = create();
    if (queryparameters != null) result.queryparameters.addAll(queryparameters);
    if (querystatement != null) result.querystatement = querystatement;
    if (viewproperties != null)
      result.viewproperties.addEntries(viewproperties);
    return result;
  }

  RequestWidget._();

  factory RequestWidget.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestWidget.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestWidget',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPS(76317660, _omitFieldNames ? '' : 'queryparameters')
    ..aOS(340852217, _omitFieldNames ? '' : 'querystatement')
    ..m<$core.String, $core.String>(
        381974398, _omitFieldNames ? '' : 'viewproperties',
        entryClassName: 'RequestWidget.ViewpropertiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cloudtrail'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestWidget clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestWidget copyWith(void Function(RequestWidget) updates) =>
      super.copyWith((message) => updates(message as RequestWidget))
          as RequestWidget;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestWidget create() => RequestWidget._();
  @$core.override
  RequestWidget createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestWidget getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestWidget>(create);
  static RequestWidget? _defaultInstance;

  @$pb.TagNumber(76317660)
  $pb.PbList<$core.String> get queryparameters => $_getList(0);

  @$pb.TagNumber(340852217)
  $core.String get querystatement => $_getSZ(1);
  @$pb.TagNumber(340852217)
  set querystatement($core.String value) => $_setString(1, value);
  @$pb.TagNumber(340852217)
  $core.bool hasQuerystatement() => $_has(1);
  @$pb.TagNumber(340852217)
  void clearQuerystatement() => $_clearField(340852217);

  @$pb.TagNumber(381974398)
  $pb.PbMap<$core.String, $core.String> get viewproperties => $_getMap(2);
}

class Resource extends $pb.GeneratedMessage {
  factory Resource({
    $core.String? resourcename,
    $core.String? resourcetype,
  }) {
    final result = create();
    if (resourcename != null) result.resourcename = resourcename;
    if (resourcetype != null) result.resourcetype = resourcetype;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(269834071, _omitFieldNames ? '' : 'resourcename')
    ..aOS(301342558, _omitFieldNames ? '' : 'resourcetype')
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

  @$pb.TagNumber(269834071)
  $core.String get resourcename => $_getSZ(0);
  @$pb.TagNumber(269834071)
  set resourcename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(269834071)
  $core.bool hasResourcename() => $_has(0);
  @$pb.TagNumber(269834071)
  void clearResourcename() => $_clearField(269834071);

  @$pb.TagNumber(301342558)
  $core.String get resourcetype => $_getSZ(1);
  @$pb.TagNumber(301342558)
  set resourcetype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(301342558)
  $core.bool hasResourcetype() => $_has(1);
  @$pb.TagNumber(301342558)
  void clearResourcetype() => $_clearField(301342558);
}

class ResourceARNNotValidException extends $pb.GeneratedMessage {
  factory ResourceARNNotValidException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceARNNotValidException._();

  factory ResourceARNNotValidException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceARNNotValidException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceARNNotValidException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceARNNotValidException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceARNNotValidException copyWith(
          void Function(ResourceARNNotValidException) updates) =>
      super.copyWith(
              (message) => updates(message as ResourceARNNotValidException))
          as ResourceARNNotValidException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceARNNotValidException create() =>
      ResourceARNNotValidException._();
  @$core.override
  ResourceARNNotValidException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceARNNotValidException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceARNNotValidException>(create);
  static ResourceARNNotValidException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
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

class ResourcePolicyNotFoundException extends $pb.GeneratedMessage {
  factory ResourcePolicyNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourcePolicyNotFoundException._();

  factory ResourcePolicyNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourcePolicyNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourcePolicyNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourcePolicyNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourcePolicyNotFoundException copyWith(
          void Function(ResourcePolicyNotFoundException) updates) =>
      super.copyWith(
              (message) => updates(message as ResourcePolicyNotFoundException))
          as ResourcePolicyNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourcePolicyNotFoundException create() =>
      ResourcePolicyNotFoundException._();
  @$core.override
  ResourcePolicyNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourcePolicyNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourcePolicyNotFoundException>(
          create);
  static ResourcePolicyNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ResourcePolicyNotValidException extends $pb.GeneratedMessage {
  factory ResourcePolicyNotValidException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourcePolicyNotValidException._();

  factory ResourcePolicyNotValidException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourcePolicyNotValidException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourcePolicyNotValidException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourcePolicyNotValidException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourcePolicyNotValidException copyWith(
          void Function(ResourcePolicyNotValidException) updates) =>
      super.copyWith(
              (message) => updates(message as ResourcePolicyNotValidException))
          as ResourcePolicyNotValidException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourcePolicyNotValidException create() =>
      ResourcePolicyNotValidException._();
  @$core.override
  ResourcePolicyNotValidException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourcePolicyNotValidException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourcePolicyNotValidException>(
          create);
  static ResourcePolicyNotValidException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ResourceTag extends $pb.GeneratedMessage {
  factory ResourceTag({
    $core.Iterable<Tag>? tagslist,
    $core.String? resourceid,
  }) {
    final result = create();
    if (tagslist != null) result.tagslist.addAll(tagslist);
    if (resourceid != null) result.resourceid = resourceid;
    return result;
  }

  ResourceTag._();

  factory ResourceTag.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceTag.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceTag',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<Tag>(497422889, _omitFieldNames ? '' : 'tagslist',
        subBuilder: Tag.create)
    ..aOS(526146833, _omitFieldNames ? '' : 'resourceid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceTag clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceTag copyWith(void Function(ResourceTag) updates) =>
      super.copyWith((message) => updates(message as ResourceTag))
          as ResourceTag;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceTag create() => ResourceTag._();
  @$core.override
  ResourceTag createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceTag getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceTag>(create);
  static ResourceTag? _defaultInstance;

  @$pb.TagNumber(497422889)
  $pb.PbList<Tag> get tagslist => $_getList(0);

  @$pb.TagNumber(526146833)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(526146833)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(526146833)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(526146833)
  void clearResourceid() => $_clearField(526146833);
}

class ResourceTypeNotSupportedException extends $pb.GeneratedMessage {
  factory ResourceTypeNotSupportedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceTypeNotSupportedException._();

  factory ResourceTypeNotSupportedException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceTypeNotSupportedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceTypeNotSupportedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceTypeNotSupportedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceTypeNotSupportedException copyWith(
          void Function(ResourceTypeNotSupportedException) updates) =>
      super.copyWith((message) =>
              updates(message as ResourceTypeNotSupportedException))
          as ResourceTypeNotSupportedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceTypeNotSupportedException create() =>
      ResourceTypeNotSupportedException._();
  @$core.override
  ResourceTypeNotSupportedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceTypeNotSupportedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceTypeNotSupportedException>(
          create);
  static ResourceTypeNotSupportedException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class RestoreEventDataStoreRequest extends $pb.GeneratedMessage {
  factory RestoreEventDataStoreRequest({
    $core.String? eventdatastore,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    return result;
  }

  RestoreEventDataStoreRequest._();

  factory RestoreEventDataStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreEventDataStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreEventDataStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreEventDataStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreEventDataStoreRequest copyWith(
          void Function(RestoreEventDataStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as RestoreEventDataStoreRequest))
          as RestoreEventDataStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreEventDataStoreRequest create() =>
      RestoreEventDataStoreRequest._();
  @$core.override
  RestoreEventDataStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreEventDataStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreEventDataStoreRequest>(create);
  static RestoreEventDataStoreRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);
}

class RestoreEventDataStoreResponse extends $pb.GeneratedMessage {
  factory RestoreEventDataStoreResponse({
    EventDataStoreStatus? status,
    $core.bool? multiregionenabled,
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? updatedtimestamp,
    $core.String? kmskeyid,
    BillingMode? billingmode,
    $core.int? retentionperiod,
    $core.String? name,
    $core.String? eventdatastorearn,
    $core.String? createdtimestamp,
    $core.bool? terminationprotectionenabled,
    $core.bool? organizationenabled,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (multiregionenabled != null)
      result.multiregionenabled = multiregionenabled;
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (billingmode != null) result.billingmode = billingmode;
    if (retentionperiod != null) result.retentionperiod = retentionperiod;
    if (name != null) result.name = name;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (organizationenabled != null)
      result.organizationenabled = organizationenabled;
    return result;
  }

  RestoreEventDataStoreResponse._();

  factory RestoreEventDataStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreEventDataStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreEventDataStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<EventDataStoreStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: EventDataStoreStatus.values)
    ..aOB(20620620, _omitFieldNames ? '' : 'multiregionenabled')
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aI(196383721, _omitFieldNames ? '' : 'retentionperiod')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOB(480171176, _omitFieldNames ? '' : 'organizationenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreEventDataStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreEventDataStoreResponse copyWith(
          void Function(RestoreEventDataStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as RestoreEventDataStoreResponse))
          as RestoreEventDataStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreEventDataStoreResponse create() =>
      RestoreEventDataStoreResponse._();
  @$core.override
  RestoreEventDataStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreEventDataStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreEventDataStoreResponse>(create);
  static RestoreEventDataStoreResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  EventDataStoreStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(EventDataStoreStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(20620620)
  $core.bool get multiregionenabled => $_getBF(1);
  @$pb.TagNumber(20620620)
  set multiregionenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(20620620)
  $core.bool hasMultiregionenabled() => $_has(1);
  @$pb.TagNumber(20620620)
  void clearMultiregionenabled() => $_clearField(20620620);

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(2);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(3);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(3);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(4);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(4);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(5);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(5);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(196383721)
  $core.int get retentionperiod => $_getIZ(6);
  @$pb.TagNumber(196383721)
  set retentionperiod($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(196383721)
  $core.bool hasRetentionperiod() => $_has(6);
  @$pb.TagNumber(196383721)
  void clearRetentionperiod() => $_clearField(196383721);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(8);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(8);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(9);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(9, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(9);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(10);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(10);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(480171176)
  $core.bool get organizationenabled => $_getBF(11);
  @$pb.TagNumber(480171176)
  set organizationenabled($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(480171176)
  $core.bool hasOrganizationenabled() => $_has(11);
  @$pb.TagNumber(480171176)
  void clearOrganizationenabled() => $_clearField(480171176);
}

class S3BucketDoesNotExistException extends $pb.GeneratedMessage {
  factory S3BucketDoesNotExistException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  S3BucketDoesNotExistException._();

  factory S3BucketDoesNotExistException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory S3BucketDoesNotExistException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'S3BucketDoesNotExistException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3BucketDoesNotExistException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3BucketDoesNotExistException copyWith(
          void Function(S3BucketDoesNotExistException) updates) =>
      super.copyWith(
              (message) => updates(message as S3BucketDoesNotExistException))
          as S3BucketDoesNotExistException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static S3BucketDoesNotExistException create() =>
      S3BucketDoesNotExistException._();
  @$core.override
  S3BucketDoesNotExistException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static S3BucketDoesNotExistException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<S3BucketDoesNotExistException>(create);
  static S3BucketDoesNotExistException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class S3ImportSource extends $pb.GeneratedMessage {
  factory S3ImportSource({
    $core.String? s3locationuri,
    $core.String? s3bucketaccessrolearn,
    $core.String? s3bucketregion,
  }) {
    final result = create();
    if (s3locationuri != null) result.s3locationuri = s3locationuri;
    if (s3bucketaccessrolearn != null)
      result.s3bucketaccessrolearn = s3bucketaccessrolearn;
    if (s3bucketregion != null) result.s3bucketregion = s3bucketregion;
    return result;
  }

  S3ImportSource._();

  factory S3ImportSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory S3ImportSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'S3ImportSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(115913071, _omitFieldNames ? '' : 's3locationuri')
    ..aOS(184258167, _omitFieldNames ? '' : 's3bucketaccessrolearn')
    ..aOS(223825074, _omitFieldNames ? '' : 's3bucketregion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3ImportSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3ImportSource copyWith(void Function(S3ImportSource) updates) =>
      super.copyWith((message) => updates(message as S3ImportSource))
          as S3ImportSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static S3ImportSource create() => S3ImportSource._();
  @$core.override
  S3ImportSource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static S3ImportSource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<S3ImportSource>(create);
  static S3ImportSource? _defaultInstance;

  @$pb.TagNumber(115913071)
  $core.String get s3locationuri => $_getSZ(0);
  @$pb.TagNumber(115913071)
  set s3locationuri($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115913071)
  $core.bool hasS3locationuri() => $_has(0);
  @$pb.TagNumber(115913071)
  void clearS3locationuri() => $_clearField(115913071);

  @$pb.TagNumber(184258167)
  $core.String get s3bucketaccessrolearn => $_getSZ(1);
  @$pb.TagNumber(184258167)
  set s3bucketaccessrolearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(184258167)
  $core.bool hasS3bucketaccessrolearn() => $_has(1);
  @$pb.TagNumber(184258167)
  void clearS3bucketaccessrolearn() => $_clearField(184258167);

  @$pb.TagNumber(223825074)
  $core.String get s3bucketregion => $_getSZ(2);
  @$pb.TagNumber(223825074)
  set s3bucketregion($core.String value) => $_setString(2, value);
  @$pb.TagNumber(223825074)
  $core.bool hasS3bucketregion() => $_has(2);
  @$pb.TagNumber(223825074)
  void clearS3bucketregion() => $_clearField(223825074);
}

class SearchSampleQueriesRequest extends $pb.GeneratedMessage {
  factory SearchSampleQueriesRequest({
    $core.String? nexttoken,
    $core.String? searchphrase,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (searchphrase != null) result.searchphrase = searchphrase;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  SearchSampleQueriesRequest._();

  factory SearchSampleQueriesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SearchSampleQueriesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchSampleQueriesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(243779077, _omitFieldNames ? '' : 'searchphrase')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchSampleQueriesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchSampleQueriesRequest copyWith(
          void Function(SearchSampleQueriesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as SearchSampleQueriesRequest))
          as SearchSampleQueriesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SearchSampleQueriesRequest create() => SearchSampleQueriesRequest._();
  @$core.override
  SearchSampleQueriesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SearchSampleQueriesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchSampleQueriesRequest>(create);
  static SearchSampleQueriesRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(243779077)
  $core.String get searchphrase => $_getSZ(1);
  @$pb.TagNumber(243779077)
  set searchphrase($core.String value) => $_setString(1, value);
  @$pb.TagNumber(243779077)
  $core.bool hasSearchphrase() => $_has(1);
  @$pb.TagNumber(243779077)
  void clearSearchphrase() => $_clearField(243779077);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class SearchSampleQueriesResponse extends $pb.GeneratedMessage {
  factory SearchSampleQueriesResponse({
    $core.Iterable<SearchSampleQueriesSearchResult>? searchresults,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (searchresults != null) result.searchresults.addAll(searchresults);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  SearchSampleQueriesResponse._();

  factory SearchSampleQueriesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SearchSampleQueriesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchSampleQueriesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<SearchSampleQueriesSearchResult>(
        29101186, _omitFieldNames ? '' : 'searchresults',
        subBuilder: SearchSampleQueriesSearchResult.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchSampleQueriesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchSampleQueriesResponse copyWith(
          void Function(SearchSampleQueriesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as SearchSampleQueriesResponse))
          as SearchSampleQueriesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SearchSampleQueriesResponse create() =>
      SearchSampleQueriesResponse._();
  @$core.override
  SearchSampleQueriesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SearchSampleQueriesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchSampleQueriesResponse>(create);
  static SearchSampleQueriesResponse? _defaultInstance;

  @$pb.TagNumber(29101186)
  $pb.PbList<SearchSampleQueriesSearchResult> get searchresults => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class SearchSampleQueriesSearchResult extends $pb.GeneratedMessage {
  factory SearchSampleQueriesSearchResult({
    $core.String? sql,
    $core.String? description,
    $core.String? name,
    $core.double? relevance,
  }) {
    final result = create();
    if (sql != null) result.sql = sql;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (relevance != null) result.relevance = relevance;
    return result;
  }

  SearchSampleQueriesSearchResult._();

  factory SearchSampleQueriesSearchResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SearchSampleQueriesSearchResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SearchSampleQueriesSearchResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(18818720, _omitFieldNames ? '' : 'sql')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aD(287814281, _omitFieldNames ? '' : 'relevance',
        fieldType: $pb.PbFieldType.OF)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchSampleQueriesSearchResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SearchSampleQueriesSearchResult copyWith(
          void Function(SearchSampleQueriesSearchResult) updates) =>
      super.copyWith(
              (message) => updates(message as SearchSampleQueriesSearchResult))
          as SearchSampleQueriesSearchResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SearchSampleQueriesSearchResult create() =>
      SearchSampleQueriesSearchResult._();
  @$core.override
  SearchSampleQueriesSearchResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SearchSampleQueriesSearchResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SearchSampleQueriesSearchResult>(
          create);
  static SearchSampleQueriesSearchResult? _defaultInstance;

  @$pb.TagNumber(18818720)
  $core.String get sql => $_getSZ(0);
  @$pb.TagNumber(18818720)
  set sql($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18818720)
  $core.bool hasSql() => $_has(0);
  @$pb.TagNumber(18818720)
  void clearSql() => $_clearField(18818720);

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

  @$pb.TagNumber(287814281)
  $core.double get relevance => $_getN(3);
  @$pb.TagNumber(287814281)
  set relevance($core.double value) => $_setFloat(3, value);
  @$pb.TagNumber(287814281)
  $core.bool hasRelevance() => $_has(3);
  @$pb.TagNumber(287814281)
  void clearRelevance() => $_clearField(287814281);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
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

class SourceConfig extends $pb.GeneratedMessage {
  factory SourceConfig({
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.bool? applytoallregions,
  }) {
    final result = create();
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (applytoallregions != null) result.applytoallregions = applytoallregions;
    return result;
  }

  SourceConfig._();

  factory SourceConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SourceConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SourceConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOB(512541819, _omitFieldNames ? '' : 'applytoallregions')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SourceConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SourceConfig copyWith(void Function(SourceConfig) updates) =>
      super.copyWith((message) => updates(message as SourceConfig))
          as SourceConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SourceConfig create() => SourceConfig._();
  @$core.override
  SourceConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SourceConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SourceConfig>(create);
  static SourceConfig? _defaultInstance;

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(0);

  @$pb.TagNumber(512541819)
  $core.bool get applytoallregions => $_getBF(1);
  @$pb.TagNumber(512541819)
  set applytoallregions($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(512541819)
  $core.bool hasApplytoallregions() => $_has(1);
  @$pb.TagNumber(512541819)
  void clearApplytoallregions() => $_clearField(512541819);
}

class StartDashboardRefreshRequest extends $pb.GeneratedMessage {
  factory StartDashboardRefreshRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        queryparametervalues,
    $core.String? dashboardid,
  }) {
    final result = create();
    if (queryparametervalues != null)
      result.queryparametervalues.addEntries(queryparametervalues);
    if (dashboardid != null) result.dashboardid = dashboardid;
    return result;
  }

  StartDashboardRefreshRequest._();

  factory StartDashboardRefreshRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartDashboardRefreshRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartDashboardRefreshRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        398732929, _omitFieldNames ? '' : 'queryparametervalues',
        entryClassName:
            'StartDashboardRefreshRequest.QueryparametervaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cloudtrail'))
    ..aOS(430615703, _omitFieldNames ? '' : 'dashboardid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartDashboardRefreshRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartDashboardRefreshRequest copyWith(
          void Function(StartDashboardRefreshRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StartDashboardRefreshRequest))
          as StartDashboardRefreshRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartDashboardRefreshRequest create() =>
      StartDashboardRefreshRequest._();
  @$core.override
  StartDashboardRefreshRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartDashboardRefreshRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartDashboardRefreshRequest>(create);
  static StartDashboardRefreshRequest? _defaultInstance;

  @$pb.TagNumber(398732929)
  $pb.PbMap<$core.String, $core.String> get queryparametervalues => $_getMap(0);

  @$pb.TagNumber(430615703)
  $core.String get dashboardid => $_getSZ(1);
  @$pb.TagNumber(430615703)
  set dashboardid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(430615703)
  $core.bool hasDashboardid() => $_has(1);
  @$pb.TagNumber(430615703)
  void clearDashboardid() => $_clearField(430615703);
}

class StartDashboardRefreshResponse extends $pb.GeneratedMessage {
  factory StartDashboardRefreshResponse({
    $core.String? refreshid,
  }) {
    final result = create();
    if (refreshid != null) result.refreshid = refreshid;
    return result;
  }

  StartDashboardRefreshResponse._();

  factory StartDashboardRefreshResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartDashboardRefreshResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartDashboardRefreshResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(57407610, _omitFieldNames ? '' : 'refreshid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartDashboardRefreshResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartDashboardRefreshResponse copyWith(
          void Function(StartDashboardRefreshResponse) updates) =>
      super.copyWith(
              (message) => updates(message as StartDashboardRefreshResponse))
          as StartDashboardRefreshResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartDashboardRefreshResponse create() =>
      StartDashboardRefreshResponse._();
  @$core.override
  StartDashboardRefreshResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartDashboardRefreshResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartDashboardRefreshResponse>(create);
  static StartDashboardRefreshResponse? _defaultInstance;

  @$pb.TagNumber(57407610)
  $core.String get refreshid => $_getSZ(0);
  @$pb.TagNumber(57407610)
  set refreshid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(57407610)
  $core.bool hasRefreshid() => $_has(0);
  @$pb.TagNumber(57407610)
  void clearRefreshid() => $_clearField(57407610);
}

class StartEventDataStoreIngestionRequest extends $pb.GeneratedMessage {
  factory StartEventDataStoreIngestionRequest({
    $core.String? eventdatastore,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    return result;
  }

  StartEventDataStoreIngestionRequest._();

  factory StartEventDataStoreIngestionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartEventDataStoreIngestionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartEventDataStoreIngestionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartEventDataStoreIngestionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartEventDataStoreIngestionRequest copyWith(
          void Function(StartEventDataStoreIngestionRequest) updates) =>
      super.copyWith((message) =>
              updates(message as StartEventDataStoreIngestionRequest))
          as StartEventDataStoreIngestionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartEventDataStoreIngestionRequest create() =>
      StartEventDataStoreIngestionRequest._();
  @$core.override
  StartEventDataStoreIngestionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartEventDataStoreIngestionRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          StartEventDataStoreIngestionRequest>(create);
  static StartEventDataStoreIngestionRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);
}

class StartEventDataStoreIngestionResponse extends $pb.GeneratedMessage {
  factory StartEventDataStoreIngestionResponse() => create();

  StartEventDataStoreIngestionResponse._();

  factory StartEventDataStoreIngestionResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartEventDataStoreIngestionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartEventDataStoreIngestionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartEventDataStoreIngestionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartEventDataStoreIngestionResponse copyWith(
          void Function(StartEventDataStoreIngestionResponse) updates) =>
      super.copyWith((message) =>
              updates(message as StartEventDataStoreIngestionResponse))
          as StartEventDataStoreIngestionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartEventDataStoreIngestionResponse create() =>
      StartEventDataStoreIngestionResponse._();
  @$core.override
  StartEventDataStoreIngestionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartEventDataStoreIngestionResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          StartEventDataStoreIngestionResponse>(create);
  static StartEventDataStoreIngestionResponse? _defaultInstance;
}

class StartImportRequest extends $pb.GeneratedMessage {
  factory StartImportRequest({
    ImportSource? importsource,
    $core.String? starteventtime,
    $core.String? endeventtime,
    $core.Iterable<$core.String>? destinations,
    $core.String? importid,
  }) {
    final result = create();
    if (importsource != null) result.importsource = importsource;
    if (starteventtime != null) result.starteventtime = starteventtime;
    if (endeventtime != null) result.endeventtime = endeventtime;
    if (destinations != null) result.destinations.addAll(destinations);
    if (importid != null) result.importid = importid;
    return result;
  }

  StartImportRequest._();

  factory StartImportRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartImportRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartImportRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<ImportSource>(41128754, _omitFieldNames ? '' : 'importsource',
        subBuilder: ImportSource.create)
    ..aOS(107272573, _omitFieldNames ? '' : 'starteventtime')
    ..aOS(260441984, _omitFieldNames ? '' : 'endeventtime')
    ..pPS(404385861, _omitFieldNames ? '' : 'destinations')
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartImportRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartImportRequest copyWith(void Function(StartImportRequest) updates) =>
      super.copyWith((message) => updates(message as StartImportRequest))
          as StartImportRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartImportRequest create() => StartImportRequest._();
  @$core.override
  StartImportRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartImportRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartImportRequest>(create);
  static StartImportRequest? _defaultInstance;

  @$pb.TagNumber(41128754)
  ImportSource get importsource => $_getN(0);
  @$pb.TagNumber(41128754)
  set importsource(ImportSource value) => $_setField(41128754, value);
  @$pb.TagNumber(41128754)
  $core.bool hasImportsource() => $_has(0);
  @$pb.TagNumber(41128754)
  void clearImportsource() => $_clearField(41128754);
  @$pb.TagNumber(41128754)
  ImportSource ensureImportsource() => $_ensure(0);

  @$pb.TagNumber(107272573)
  $core.String get starteventtime => $_getSZ(1);
  @$pb.TagNumber(107272573)
  set starteventtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(107272573)
  $core.bool hasStarteventtime() => $_has(1);
  @$pb.TagNumber(107272573)
  void clearStarteventtime() => $_clearField(107272573);

  @$pb.TagNumber(260441984)
  $core.String get endeventtime => $_getSZ(2);
  @$pb.TagNumber(260441984)
  set endeventtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(260441984)
  $core.bool hasEndeventtime() => $_has(2);
  @$pb.TagNumber(260441984)
  void clearEndeventtime() => $_clearField(260441984);

  @$pb.TagNumber(404385861)
  $pb.PbList<$core.String> get destinations => $_getList(3);

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(4);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(4);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class StartImportResponse extends $pb.GeneratedMessage {
  factory StartImportResponse({
    ImportSource? importsource,
    $core.String? updatedtimestamp,
    $core.String? starteventtime,
    ImportStatus? importstatus,
    $core.String? endeventtime,
    $core.String? createdtimestamp,
    $core.Iterable<$core.String>? destinations,
    $core.String? importid,
  }) {
    final result = create();
    if (importsource != null) result.importsource = importsource;
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (starteventtime != null) result.starteventtime = starteventtime;
    if (importstatus != null) result.importstatus = importstatus;
    if (endeventtime != null) result.endeventtime = endeventtime;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (destinations != null) result.destinations.addAll(destinations);
    if (importid != null) result.importid = importid;
    return result;
  }

  StartImportResponse._();

  factory StartImportResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartImportResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartImportResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<ImportSource>(41128754, _omitFieldNames ? '' : 'importsource',
        subBuilder: ImportSource.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOS(107272573, _omitFieldNames ? '' : 'starteventtime')
    ..aE<ImportStatus>(129077631, _omitFieldNames ? '' : 'importstatus',
        enumValues: ImportStatus.values)
    ..aOS(260441984, _omitFieldNames ? '' : 'endeventtime')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..pPS(404385861, _omitFieldNames ? '' : 'destinations')
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartImportResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartImportResponse copyWith(void Function(StartImportResponse) updates) =>
      super.copyWith((message) => updates(message as StartImportResponse))
          as StartImportResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartImportResponse create() => StartImportResponse._();
  @$core.override
  StartImportResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartImportResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartImportResponse>(create);
  static StartImportResponse? _defaultInstance;

  @$pb.TagNumber(41128754)
  ImportSource get importsource => $_getN(0);
  @$pb.TagNumber(41128754)
  set importsource(ImportSource value) => $_setField(41128754, value);
  @$pb.TagNumber(41128754)
  $core.bool hasImportsource() => $_has(0);
  @$pb.TagNumber(41128754)
  void clearImportsource() => $_clearField(41128754);
  @$pb.TagNumber(41128754)
  ImportSource ensureImportsource() => $_ensure(0);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(1);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(1);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(107272573)
  $core.String get starteventtime => $_getSZ(2);
  @$pb.TagNumber(107272573)
  set starteventtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(107272573)
  $core.bool hasStarteventtime() => $_has(2);
  @$pb.TagNumber(107272573)
  void clearStarteventtime() => $_clearField(107272573);

  @$pb.TagNumber(129077631)
  ImportStatus get importstatus => $_getN(3);
  @$pb.TagNumber(129077631)
  set importstatus(ImportStatus value) => $_setField(129077631, value);
  @$pb.TagNumber(129077631)
  $core.bool hasImportstatus() => $_has(3);
  @$pb.TagNumber(129077631)
  void clearImportstatus() => $_clearField(129077631);

  @$pb.TagNumber(260441984)
  $core.String get endeventtime => $_getSZ(4);
  @$pb.TagNumber(260441984)
  set endeventtime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(260441984)
  $core.bool hasEndeventtime() => $_has(4);
  @$pb.TagNumber(260441984)
  void clearEndeventtime() => $_clearField(260441984);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(5);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(5, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(5);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(404385861)
  $pb.PbList<$core.String> get destinations => $_getList(6);

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(7);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(7, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(7);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class StartLoggingRequest extends $pb.GeneratedMessage {
  factory StartLoggingRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  StartLoggingRequest._();

  factory StartLoggingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartLoggingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartLoggingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartLoggingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartLoggingRequest copyWith(void Function(StartLoggingRequest) updates) =>
      super.copyWith((message) => updates(message as StartLoggingRequest))
          as StartLoggingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartLoggingRequest create() => StartLoggingRequest._();
  @$core.override
  StartLoggingRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartLoggingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartLoggingRequest>(create);
  static StartLoggingRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class StartLoggingResponse extends $pb.GeneratedMessage {
  factory StartLoggingResponse() => create();

  StartLoggingResponse._();

  factory StartLoggingResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartLoggingResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartLoggingResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartLoggingResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartLoggingResponse copyWith(void Function(StartLoggingResponse) updates) =>
      super.copyWith((message) => updates(message as StartLoggingResponse))
          as StartLoggingResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartLoggingResponse create() => StartLoggingResponse._();
  @$core.override
  StartLoggingResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartLoggingResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartLoggingResponse>(create);
  static StartLoggingResponse? _defaultInstance;
}

class StartQueryRequest extends $pb.GeneratedMessage {
  factory StartQueryRequest({
    $core.Iterable<$core.String>? queryparameters,
    $core.String? deliverys3uri,
    $core.String? querystatement,
    $core.String? eventdatastoreowneraccountid,
    $core.String? queryalias,
  }) {
    final result = create();
    if (queryparameters != null) result.queryparameters.addAll(queryparameters);
    if (deliverys3uri != null) result.deliverys3uri = deliverys3uri;
    if (querystatement != null) result.querystatement = querystatement;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    if (queryalias != null) result.queryalias = queryalias;
    return result;
  }

  StartQueryRequest._();

  factory StartQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartQueryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPS(76317660, _omitFieldNames ? '' : 'queryparameters')
    ..aOS(230884460, _omitFieldNames ? '' : 'deliverys3uri')
    ..aOS(340852217, _omitFieldNames ? '' : 'querystatement')
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..aOS(512762270, _omitFieldNames ? '' : 'queryalias')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryRequest copyWith(void Function(StartQueryRequest) updates) =>
      super.copyWith((message) => updates(message as StartQueryRequest))
          as StartQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartQueryRequest create() => StartQueryRequest._();
  @$core.override
  StartQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartQueryRequest>(create);
  static StartQueryRequest? _defaultInstance;

  @$pb.TagNumber(76317660)
  $pb.PbList<$core.String> get queryparameters => $_getList(0);

  @$pb.TagNumber(230884460)
  $core.String get deliverys3uri => $_getSZ(1);
  @$pb.TagNumber(230884460)
  set deliverys3uri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(230884460)
  $core.bool hasDeliverys3uri() => $_has(1);
  @$pb.TagNumber(230884460)
  void clearDeliverys3uri() => $_clearField(230884460);

  @$pb.TagNumber(340852217)
  $core.String get querystatement => $_getSZ(2);
  @$pb.TagNumber(340852217)
  set querystatement($core.String value) => $_setString(2, value);
  @$pb.TagNumber(340852217)
  $core.bool hasQuerystatement() => $_has(2);
  @$pb.TagNumber(340852217)
  void clearQuerystatement() => $_clearField(340852217);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(3);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(3);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);

  @$pb.TagNumber(512762270)
  $core.String get queryalias => $_getSZ(4);
  @$pb.TagNumber(512762270)
  set queryalias($core.String value) => $_setString(4, value);
  @$pb.TagNumber(512762270)
  $core.bool hasQueryalias() => $_has(4);
  @$pb.TagNumber(512762270)
  void clearQueryalias() => $_clearField(512762270);
}

class StartQueryResponse extends $pb.GeneratedMessage {
  factory StartQueryResponse({
    $core.String? queryid,
    $core.String? eventdatastoreowneraccountid,
  }) {
    final result = create();
    if (queryid != null) result.queryid = queryid;
    if (eventdatastoreowneraccountid != null)
      result.eventdatastoreowneraccountid = eventdatastoreowneraccountid;
    return result;
  }

  StartQueryResponse._();

  factory StartQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartQueryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..aOS(471609008, _omitFieldNames ? '' : 'eventdatastoreowneraccountid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryResponse copyWith(void Function(StartQueryResponse) updates) =>
      super.copyWith((message) => updates(message as StartQueryResponse))
          as StartQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartQueryResponse create() => StartQueryResponse._();
  @$core.override
  StartQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartQueryResponse>(create);
  static StartQueryResponse? _defaultInstance;

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(0);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(0);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(471609008)
  $core.String get eventdatastoreowneraccountid => $_getSZ(1);
  @$pb.TagNumber(471609008)
  set eventdatastoreowneraccountid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(471609008)
  $core.bool hasEventdatastoreowneraccountid() => $_has(1);
  @$pb.TagNumber(471609008)
  void clearEventdatastoreowneraccountid() => $_clearField(471609008);
}

class StopEventDataStoreIngestionRequest extends $pb.GeneratedMessage {
  factory StopEventDataStoreIngestionRequest({
    $core.String? eventdatastore,
  }) {
    final result = create();
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    return result;
  }

  StopEventDataStoreIngestionRequest._();

  factory StopEventDataStoreIngestionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopEventDataStoreIngestionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopEventDataStoreIngestionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopEventDataStoreIngestionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopEventDataStoreIngestionRequest copyWith(
          void Function(StopEventDataStoreIngestionRequest) updates) =>
      super.copyWith((message) =>
              updates(message as StopEventDataStoreIngestionRequest))
          as StopEventDataStoreIngestionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopEventDataStoreIngestionRequest create() =>
      StopEventDataStoreIngestionRequest._();
  @$core.override
  StopEventDataStoreIngestionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopEventDataStoreIngestionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopEventDataStoreIngestionRequest>(
          create);
  static StopEventDataStoreIngestionRequest? _defaultInstance;

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(0);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(0, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(0);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);
}

class StopEventDataStoreIngestionResponse extends $pb.GeneratedMessage {
  factory StopEventDataStoreIngestionResponse() => create();

  StopEventDataStoreIngestionResponse._();

  factory StopEventDataStoreIngestionResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopEventDataStoreIngestionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopEventDataStoreIngestionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopEventDataStoreIngestionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopEventDataStoreIngestionResponse copyWith(
          void Function(StopEventDataStoreIngestionResponse) updates) =>
      super.copyWith((message) =>
              updates(message as StopEventDataStoreIngestionResponse))
          as StopEventDataStoreIngestionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopEventDataStoreIngestionResponse create() =>
      StopEventDataStoreIngestionResponse._();
  @$core.override
  StopEventDataStoreIngestionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopEventDataStoreIngestionResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          StopEventDataStoreIngestionResponse>(create);
  static StopEventDataStoreIngestionResponse? _defaultInstance;
}

class StopImportRequest extends $pb.GeneratedMessage {
  factory StopImportRequest({
    $core.String? importid,
  }) {
    final result = create();
    if (importid != null) result.importid = importid;
    return result;
  }

  StopImportRequest._();

  factory StopImportRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopImportRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopImportRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopImportRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopImportRequest copyWith(void Function(StopImportRequest) updates) =>
      super.copyWith((message) => updates(message as StopImportRequest))
          as StopImportRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopImportRequest create() => StopImportRequest._();
  @$core.override
  StopImportRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopImportRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopImportRequest>(create);
  static StopImportRequest? _defaultInstance;

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(0);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(0);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class StopImportResponse extends $pb.GeneratedMessage {
  factory StopImportResponse({
    ImportSource? importsource,
    $core.String? updatedtimestamp,
    ImportStatistics? importstatistics,
    $core.String? starteventtime,
    ImportStatus? importstatus,
    $core.String? endeventtime,
    $core.String? createdtimestamp,
    $core.Iterable<$core.String>? destinations,
    $core.String? importid,
  }) {
    final result = create();
    if (importsource != null) result.importsource = importsource;
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (importstatistics != null) result.importstatistics = importstatistics;
    if (starteventtime != null) result.starteventtime = starteventtime;
    if (importstatus != null) result.importstatus = importstatus;
    if (endeventtime != null) result.endeventtime = endeventtime;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (destinations != null) result.destinations.addAll(destinations);
    if (importid != null) result.importid = importid;
    return result;
  }

  StopImportResponse._();

  factory StopImportResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopImportResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopImportResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<ImportSource>(41128754, _omitFieldNames ? '' : 'importsource',
        subBuilder: ImportSource.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOM<ImportStatistics>(48175528, _omitFieldNames ? '' : 'importstatistics',
        subBuilder: ImportStatistics.create)
    ..aOS(107272573, _omitFieldNames ? '' : 'starteventtime')
    ..aE<ImportStatus>(129077631, _omitFieldNames ? '' : 'importstatus',
        enumValues: ImportStatus.values)
    ..aOS(260441984, _omitFieldNames ? '' : 'endeventtime')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..pPS(404385861, _omitFieldNames ? '' : 'destinations')
    ..aOS(420153946, _omitFieldNames ? '' : 'importid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopImportResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopImportResponse copyWith(void Function(StopImportResponse) updates) =>
      super.copyWith((message) => updates(message as StopImportResponse))
          as StopImportResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopImportResponse create() => StopImportResponse._();
  @$core.override
  StopImportResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopImportResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopImportResponse>(create);
  static StopImportResponse? _defaultInstance;

  @$pb.TagNumber(41128754)
  ImportSource get importsource => $_getN(0);
  @$pb.TagNumber(41128754)
  set importsource(ImportSource value) => $_setField(41128754, value);
  @$pb.TagNumber(41128754)
  $core.bool hasImportsource() => $_has(0);
  @$pb.TagNumber(41128754)
  void clearImportsource() => $_clearField(41128754);
  @$pb.TagNumber(41128754)
  ImportSource ensureImportsource() => $_ensure(0);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(1);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(1);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(48175528)
  ImportStatistics get importstatistics => $_getN(2);
  @$pb.TagNumber(48175528)
  set importstatistics(ImportStatistics value) => $_setField(48175528, value);
  @$pb.TagNumber(48175528)
  $core.bool hasImportstatistics() => $_has(2);
  @$pb.TagNumber(48175528)
  void clearImportstatistics() => $_clearField(48175528);
  @$pb.TagNumber(48175528)
  ImportStatistics ensureImportstatistics() => $_ensure(2);

  @$pb.TagNumber(107272573)
  $core.String get starteventtime => $_getSZ(3);
  @$pb.TagNumber(107272573)
  set starteventtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(107272573)
  $core.bool hasStarteventtime() => $_has(3);
  @$pb.TagNumber(107272573)
  void clearStarteventtime() => $_clearField(107272573);

  @$pb.TagNumber(129077631)
  ImportStatus get importstatus => $_getN(4);
  @$pb.TagNumber(129077631)
  set importstatus(ImportStatus value) => $_setField(129077631, value);
  @$pb.TagNumber(129077631)
  $core.bool hasImportstatus() => $_has(4);
  @$pb.TagNumber(129077631)
  void clearImportstatus() => $_clearField(129077631);

  @$pb.TagNumber(260441984)
  $core.String get endeventtime => $_getSZ(5);
  @$pb.TagNumber(260441984)
  set endeventtime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(260441984)
  $core.bool hasEndeventtime() => $_has(5);
  @$pb.TagNumber(260441984)
  void clearEndeventtime() => $_clearField(260441984);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(6);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(6, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(6);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(404385861)
  $pb.PbList<$core.String> get destinations => $_getList(7);

  @$pb.TagNumber(420153946)
  $core.String get importid => $_getSZ(8);
  @$pb.TagNumber(420153946)
  set importid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(420153946)
  $core.bool hasImportid() => $_has(8);
  @$pb.TagNumber(420153946)
  void clearImportid() => $_clearField(420153946);
}

class StopLoggingRequest extends $pb.GeneratedMessage {
  factory StopLoggingRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  StopLoggingRequest._();

  factory StopLoggingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopLoggingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopLoggingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopLoggingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopLoggingRequest copyWith(void Function(StopLoggingRequest) updates) =>
      super.copyWith((message) => updates(message as StopLoggingRequest))
          as StopLoggingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopLoggingRequest create() => StopLoggingRequest._();
  @$core.override
  StopLoggingRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopLoggingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopLoggingRequest>(create);
  static StopLoggingRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class StopLoggingResponse extends $pb.GeneratedMessage {
  factory StopLoggingResponse() => create();

  StopLoggingResponse._();

  factory StopLoggingResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopLoggingResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopLoggingResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopLoggingResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopLoggingResponse copyWith(void Function(StopLoggingResponse) updates) =>
      super.copyWith((message) => updates(message as StopLoggingResponse))
          as StopLoggingResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopLoggingResponse create() => StopLoggingResponse._();
  @$core.override
  StopLoggingResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopLoggingResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopLoggingResponse>(create);
  static StopLoggingResponse? _defaultInstance;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
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

class TagsLimitExceededException extends $pb.GeneratedMessage {
  factory TagsLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TagsLimitExceededException._();

  factory TagsLimitExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagsLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagsLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagsLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagsLimitExceededException copyWith(
          void Function(TagsLimitExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as TagsLimitExceededException))
          as TagsLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagsLimitExceededException create() => TagsLimitExceededException._();
  @$core.override
  TagsLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagsLimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagsLimitExceededException>(create);
  static TagsLimitExceededException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
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

class Trail extends $pb.GeneratedMessage {
  factory Trail({
    $core.bool? logfilevalidationenabled,
    $core.String? trailarn,
    $core.String? kmskeyid,
    $core.String? cloudwatchlogsrolearn,
    $core.bool? hascustomeventselectors,
    $core.bool? isorganizationtrail,
    $core.bool? hasinsightselectors,
    $core.String? s3keyprefix,
    $core.bool? includeglobalserviceevents,
    $core.String? name,
    $core.String? homeregion,
    $core.String? s3bucketname,
    $core.String? snstopicarn,
    $core.String? snstopicname,
    $core.bool? ismultiregiontrail,
    $core.String? cloudwatchlogsloggrouparn,
  }) {
    final result = create();
    if (logfilevalidationenabled != null)
      result.logfilevalidationenabled = logfilevalidationenabled;
    if (trailarn != null) result.trailarn = trailarn;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (cloudwatchlogsrolearn != null)
      result.cloudwatchlogsrolearn = cloudwatchlogsrolearn;
    if (hascustomeventselectors != null)
      result.hascustomeventselectors = hascustomeventselectors;
    if (isorganizationtrail != null)
      result.isorganizationtrail = isorganizationtrail;
    if (hasinsightselectors != null)
      result.hasinsightselectors = hasinsightselectors;
    if (s3keyprefix != null) result.s3keyprefix = s3keyprefix;
    if (includeglobalserviceevents != null)
      result.includeglobalserviceevents = includeglobalserviceevents;
    if (name != null) result.name = name;
    if (homeregion != null) result.homeregion = homeregion;
    if (s3bucketname != null) result.s3bucketname = s3bucketname;
    if (snstopicarn != null) result.snstopicarn = snstopicarn;
    if (snstopicname != null) result.snstopicname = snstopicname;
    if (ismultiregiontrail != null)
      result.ismultiregiontrail = ismultiregiontrail;
    if (cloudwatchlogsloggrouparn != null)
      result.cloudwatchlogsloggrouparn = cloudwatchlogsloggrouparn;
    return result;
  }

  Trail._();

  factory Trail.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Trail.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Trail',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOB(35904346, _omitFieldNames ? '' : 'logfilevalidationenabled')
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(55454690, _omitFieldNames ? '' : 'cloudwatchlogsrolearn')
    ..aOB(79311547, _omitFieldNames ? '' : 'hascustomeventselectors')
    ..aOB(145256127, _omitFieldNames ? '' : 'isorganizationtrail')
    ..aOB(146532978, _omitFieldNames ? '' : 'hasinsightselectors')
    ..aOS(206015359, _omitFieldNames ? '' : 's3keyprefix')
    ..aOB(212227451, _omitFieldNames ? '' : 'includeglobalserviceevents')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(300797531, _omitFieldNames ? '' : 'homeregion')
    ..aOS(320495427, _omitFieldNames ? '' : 's3bucketname')
    ..aOS(380025580, _omitFieldNames ? '' : 'snstopicarn')
    ..aOS(415454800, _omitFieldNames ? '' : 'snstopicname')
    ..aOB(468658571, _omitFieldNames ? '' : 'ismultiregiontrail')
    ..aOS(519613355, _omitFieldNames ? '' : 'cloudwatchlogsloggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Trail clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Trail copyWith(void Function(Trail) updates) =>
      super.copyWith((message) => updates(message as Trail)) as Trail;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Trail create() => Trail._();
  @$core.override
  Trail createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Trail getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Trail>(create);
  static Trail? _defaultInstance;

  @$pb.TagNumber(35904346)
  $core.bool get logfilevalidationenabled => $_getBF(0);
  @$pb.TagNumber(35904346)
  set logfilevalidationenabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(35904346)
  $core.bool hasLogfilevalidationenabled() => $_has(0);
  @$pb.TagNumber(35904346)
  void clearLogfilevalidationenabled() => $_clearField(35904346);

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(1);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(1);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(2);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(2);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(55454690)
  $core.String get cloudwatchlogsrolearn => $_getSZ(3);
  @$pb.TagNumber(55454690)
  set cloudwatchlogsrolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(55454690)
  $core.bool hasCloudwatchlogsrolearn() => $_has(3);
  @$pb.TagNumber(55454690)
  void clearCloudwatchlogsrolearn() => $_clearField(55454690);

  @$pb.TagNumber(79311547)
  $core.bool get hascustomeventselectors => $_getBF(4);
  @$pb.TagNumber(79311547)
  set hascustomeventselectors($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(79311547)
  $core.bool hasHascustomeventselectors() => $_has(4);
  @$pb.TagNumber(79311547)
  void clearHascustomeventselectors() => $_clearField(79311547);

  @$pb.TagNumber(145256127)
  $core.bool get isorganizationtrail => $_getBF(5);
  @$pb.TagNumber(145256127)
  set isorganizationtrail($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(145256127)
  $core.bool hasIsorganizationtrail() => $_has(5);
  @$pb.TagNumber(145256127)
  void clearIsorganizationtrail() => $_clearField(145256127);

  @$pb.TagNumber(146532978)
  $core.bool get hasinsightselectors => $_getBF(6);
  @$pb.TagNumber(146532978)
  set hasinsightselectors($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(146532978)
  $core.bool hasHasinsightselectors() => $_has(6);
  @$pb.TagNumber(146532978)
  void clearHasinsightselectors() => $_clearField(146532978);

  @$pb.TagNumber(206015359)
  $core.String get s3keyprefix => $_getSZ(7);
  @$pb.TagNumber(206015359)
  set s3keyprefix($core.String value) => $_setString(7, value);
  @$pb.TagNumber(206015359)
  $core.bool hasS3keyprefix() => $_has(7);
  @$pb.TagNumber(206015359)
  void clearS3keyprefix() => $_clearField(206015359);

  @$pb.TagNumber(212227451)
  $core.bool get includeglobalserviceevents => $_getBF(8);
  @$pb.TagNumber(212227451)
  set includeglobalserviceevents($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(212227451)
  $core.bool hasIncludeglobalserviceevents() => $_has(8);
  @$pb.TagNumber(212227451)
  void clearIncludeglobalserviceevents() => $_clearField(212227451);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(9);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(9, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(9);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(300797531)
  $core.String get homeregion => $_getSZ(10);
  @$pb.TagNumber(300797531)
  set homeregion($core.String value) => $_setString(10, value);
  @$pb.TagNumber(300797531)
  $core.bool hasHomeregion() => $_has(10);
  @$pb.TagNumber(300797531)
  void clearHomeregion() => $_clearField(300797531);

  @$pb.TagNumber(320495427)
  $core.String get s3bucketname => $_getSZ(11);
  @$pb.TagNumber(320495427)
  set s3bucketname($core.String value) => $_setString(11, value);
  @$pb.TagNumber(320495427)
  $core.bool hasS3bucketname() => $_has(11);
  @$pb.TagNumber(320495427)
  void clearS3bucketname() => $_clearField(320495427);

  @$pb.TagNumber(380025580)
  $core.String get snstopicarn => $_getSZ(12);
  @$pb.TagNumber(380025580)
  set snstopicarn($core.String value) => $_setString(12, value);
  @$pb.TagNumber(380025580)
  $core.bool hasSnstopicarn() => $_has(12);
  @$pb.TagNumber(380025580)
  void clearSnstopicarn() => $_clearField(380025580);

  @$pb.TagNumber(415454800)
  $core.String get snstopicname => $_getSZ(13);
  @$pb.TagNumber(415454800)
  set snstopicname($core.String value) => $_setString(13, value);
  @$pb.TagNumber(415454800)
  $core.bool hasSnstopicname() => $_has(13);
  @$pb.TagNumber(415454800)
  void clearSnstopicname() => $_clearField(415454800);

  @$pb.TagNumber(468658571)
  $core.bool get ismultiregiontrail => $_getBF(14);
  @$pb.TagNumber(468658571)
  set ismultiregiontrail($core.bool value) => $_setBool(14, value);
  @$pb.TagNumber(468658571)
  $core.bool hasIsmultiregiontrail() => $_has(14);
  @$pb.TagNumber(468658571)
  void clearIsmultiregiontrail() => $_clearField(468658571);

  @$pb.TagNumber(519613355)
  $core.String get cloudwatchlogsloggrouparn => $_getSZ(15);
  @$pb.TagNumber(519613355)
  set cloudwatchlogsloggrouparn($core.String value) => $_setString(15, value);
  @$pb.TagNumber(519613355)
  $core.bool hasCloudwatchlogsloggrouparn() => $_has(15);
  @$pb.TagNumber(519613355)
  void clearCloudwatchlogsloggrouparn() => $_clearField(519613355);
}

class TrailAlreadyExistsException extends $pb.GeneratedMessage {
  factory TrailAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TrailAlreadyExistsException._();

  factory TrailAlreadyExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrailAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrailAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailAlreadyExistsException copyWith(
          void Function(TrailAlreadyExistsException) updates) =>
      super.copyWith(
              (message) => updates(message as TrailAlreadyExistsException))
          as TrailAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrailAlreadyExistsException create() =>
      TrailAlreadyExistsException._();
  @$core.override
  TrailAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrailAlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrailAlreadyExistsException>(create);
  static TrailAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class TrailInfo extends $pb.GeneratedMessage {
  factory TrailInfo({
    $core.String? trailarn,
    $core.String? name,
    $core.String? homeregion,
  }) {
    final result = create();
    if (trailarn != null) result.trailarn = trailarn;
    if (name != null) result.name = name;
    if (homeregion != null) result.homeregion = homeregion;
    return result;
  }

  TrailInfo._();

  factory TrailInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrailInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrailInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(300797531, _omitFieldNames ? '' : 'homeregion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailInfo copyWith(void Function(TrailInfo) updates) =>
      super.copyWith((message) => updates(message as TrailInfo)) as TrailInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrailInfo create() => TrailInfo._();
  @$core.override
  TrailInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrailInfo getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TrailInfo>(create);
  static TrailInfo? _defaultInstance;

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(0);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(0);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(300797531)
  $core.String get homeregion => $_getSZ(2);
  @$pb.TagNumber(300797531)
  set homeregion($core.String value) => $_setString(2, value);
  @$pb.TagNumber(300797531)
  $core.bool hasHomeregion() => $_has(2);
  @$pb.TagNumber(300797531)
  void clearHomeregion() => $_clearField(300797531);
}

class TrailNotFoundException extends $pb.GeneratedMessage {
  factory TrailNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TrailNotFoundException._();

  factory TrailNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrailNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrailNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailNotFoundException copyWith(
          void Function(TrailNotFoundException) updates) =>
      super.copyWith((message) => updates(message as TrailNotFoundException))
          as TrailNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrailNotFoundException create() => TrailNotFoundException._();
  @$core.override
  TrailNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrailNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrailNotFoundException>(create);
  static TrailNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class TrailNotProvidedException extends $pb.GeneratedMessage {
  factory TrailNotProvidedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TrailNotProvidedException._();

  factory TrailNotProvidedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrailNotProvidedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrailNotProvidedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailNotProvidedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrailNotProvidedException copyWith(
          void Function(TrailNotProvidedException) updates) =>
      super.copyWith((message) => updates(message as TrailNotProvidedException))
          as TrailNotProvidedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrailNotProvidedException create() => TrailNotProvidedException._();
  @$core.override
  TrailNotProvidedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrailNotProvidedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrailNotProvidedException>(create);
  static TrailNotProvidedException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class UnsupportedOperationException extends $pb.GeneratedMessage {
  factory UnsupportedOperationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  UnsupportedOperationException._();

  factory UnsupportedOperationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnsupportedOperationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnsupportedOperationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsupportedOperationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsupportedOperationException copyWith(
          void Function(UnsupportedOperationException) updates) =>
      super.copyWith(
              (message) => updates(message as UnsupportedOperationException))
          as UnsupportedOperationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnsupportedOperationException create() =>
      UnsupportedOperationException._();
  @$core.override
  UnsupportedOperationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnsupportedOperationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnsupportedOperationException>(create);
  static UnsupportedOperationException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class UpdateChannelRequest extends $pb.GeneratedMessage {
  factory UpdateChannelRequest({
    $core.String? channel,
    $core.String? name,
    $core.Iterable<Destination>? destinations,
  }) {
    final result = create();
    if (channel != null) result.channel = channel;
    if (name != null) result.name = name;
    if (destinations != null) result.destinations.addAll(destinations);
    return result;
  }

  UpdateChannelRequest._();

  factory UpdateChannelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateChannelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateChannelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(90016325, _omitFieldNames ? '' : 'channel')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Destination>(404385861, _omitFieldNames ? '' : 'destinations',
        subBuilder: Destination.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateChannelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateChannelRequest copyWith(void Function(UpdateChannelRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateChannelRequest))
          as UpdateChannelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateChannelRequest create() => UpdateChannelRequest._();
  @$core.override
  UpdateChannelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateChannelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateChannelRequest>(create);
  static UpdateChannelRequest? _defaultInstance;

  @$pb.TagNumber(90016325)
  $core.String get channel => $_getSZ(0);
  @$pb.TagNumber(90016325)
  set channel($core.String value) => $_setString(0, value);
  @$pb.TagNumber(90016325)
  $core.bool hasChannel() => $_has(0);
  @$pb.TagNumber(90016325)
  void clearChannel() => $_clearField(90016325);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(404385861)
  $pb.PbList<Destination> get destinations => $_getList(2);
}

class UpdateChannelResponse extends $pb.GeneratedMessage {
  factory UpdateChannelResponse({
    $core.String? source,
    $core.String? channelarn,
    $core.String? name,
    $core.Iterable<Destination>? destinations,
  }) {
    final result = create();
    if (source != null) result.source = source;
    if (channelarn != null) result.channelarn = channelarn;
    if (name != null) result.name = name;
    if (destinations != null) result.destinations.addAll(destinations);
    return result;
  }

  UpdateChannelResponse._();

  factory UpdateChannelResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateChannelResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateChannelResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(31630329, _omitFieldNames ? '' : 'source')
    ..aOS(99276476, _omitFieldNames ? '' : 'channelarn')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Destination>(404385861, _omitFieldNames ? '' : 'destinations',
        subBuilder: Destination.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateChannelResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateChannelResponse copyWith(
          void Function(UpdateChannelResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateChannelResponse))
          as UpdateChannelResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateChannelResponse create() => UpdateChannelResponse._();
  @$core.override
  UpdateChannelResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateChannelResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateChannelResponse>(create);
  static UpdateChannelResponse? _defaultInstance;

  @$pb.TagNumber(31630329)
  $core.String get source => $_getSZ(0);
  @$pb.TagNumber(31630329)
  set source($core.String value) => $_setString(0, value);
  @$pb.TagNumber(31630329)
  $core.bool hasSource() => $_has(0);
  @$pb.TagNumber(31630329)
  void clearSource() => $_clearField(31630329);

  @$pb.TagNumber(99276476)
  $core.String get channelarn => $_getSZ(1);
  @$pb.TagNumber(99276476)
  set channelarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(99276476)
  $core.bool hasChannelarn() => $_has(1);
  @$pb.TagNumber(99276476)
  void clearChannelarn() => $_clearField(99276476);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(404385861)
  $pb.PbList<Destination> get destinations => $_getList(3);
}

class UpdateDashboardRequest extends $pb.GeneratedMessage {
  factory UpdateDashboardRequest({
    RefreshSchedule? refreshschedule,
    $core.bool? terminationprotectionenabled,
    $core.String? dashboardid,
    $core.Iterable<RequestWidget>? widgets,
  }) {
    final result = create();
    if (refreshschedule != null) result.refreshschedule = refreshschedule;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (dashboardid != null) result.dashboardid = dashboardid;
    if (widgets != null) result.widgets.addAll(widgets);
    return result;
  }

  UpdateDashboardRequest._();

  factory UpdateDashboardRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDashboardRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDashboardRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOM<RefreshSchedule>(261773338, _omitFieldNames ? '' : 'refreshschedule',
        subBuilder: RefreshSchedule.create)
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOS(430615703, _omitFieldNames ? '' : 'dashboardid')
    ..pPM<RequestWidget>(501826147, _omitFieldNames ? '' : 'widgets',
        subBuilder: RequestWidget.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDashboardRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDashboardRequest copyWith(
          void Function(UpdateDashboardRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateDashboardRequest))
          as UpdateDashboardRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDashboardRequest create() => UpdateDashboardRequest._();
  @$core.override
  UpdateDashboardRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDashboardRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDashboardRequest>(create);
  static UpdateDashboardRequest? _defaultInstance;

  @$pb.TagNumber(261773338)
  RefreshSchedule get refreshschedule => $_getN(0);
  @$pb.TagNumber(261773338)
  set refreshschedule(RefreshSchedule value) => $_setField(261773338, value);
  @$pb.TagNumber(261773338)
  $core.bool hasRefreshschedule() => $_has(0);
  @$pb.TagNumber(261773338)
  void clearRefreshschedule() => $_clearField(261773338);
  @$pb.TagNumber(261773338)
  RefreshSchedule ensureRefreshschedule() => $_ensure(0);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(1);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(1);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(430615703)
  $core.String get dashboardid => $_getSZ(2);
  @$pb.TagNumber(430615703)
  set dashboardid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(430615703)
  $core.bool hasDashboardid() => $_has(2);
  @$pb.TagNumber(430615703)
  void clearDashboardid() => $_clearField(430615703);

  @$pb.TagNumber(501826147)
  $pb.PbList<RequestWidget> get widgets => $_getList(3);
}

class UpdateDashboardResponse extends $pb.GeneratedMessage {
  factory UpdateDashboardResponse({
    $core.String? updatedtimestamp,
    $core.String? dashboardarn,
    RefreshSchedule? refreshschedule,
    $core.String? name,
    DashboardType? type,
    $core.String? createdtimestamp,
    $core.bool? terminationprotectionenabled,
    $core.Iterable<Widget>? widgets,
  }) {
    final result = create();
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (dashboardarn != null) result.dashboardarn = dashboardarn;
    if (refreshschedule != null) result.refreshschedule = refreshschedule;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (widgets != null) result.widgets.addAll(widgets);
    return result;
  }

  UpdateDashboardResponse._();

  factory UpdateDashboardResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDashboardResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDashboardResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOS(108051951, _omitFieldNames ? '' : 'dashboardarn')
    ..aOM<RefreshSchedule>(261773338, _omitFieldNames ? '' : 'refreshschedule',
        subBuilder: RefreshSchedule.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<DashboardType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DashboardType.values)
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..pPM<Widget>(501826147, _omitFieldNames ? '' : 'widgets',
        subBuilder: Widget.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDashboardResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDashboardResponse copyWith(
          void Function(UpdateDashboardResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateDashboardResponse))
          as UpdateDashboardResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDashboardResponse create() => UpdateDashboardResponse._();
  @$core.override
  UpdateDashboardResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDashboardResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDashboardResponse>(create);
  static UpdateDashboardResponse? _defaultInstance;

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(0);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(0, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(0);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(108051951)
  $core.String get dashboardarn => $_getSZ(1);
  @$pb.TagNumber(108051951)
  set dashboardarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(108051951)
  $core.bool hasDashboardarn() => $_has(1);
  @$pb.TagNumber(108051951)
  void clearDashboardarn() => $_clearField(108051951);

  @$pb.TagNumber(261773338)
  RefreshSchedule get refreshschedule => $_getN(2);
  @$pb.TagNumber(261773338)
  set refreshschedule(RefreshSchedule value) => $_setField(261773338, value);
  @$pb.TagNumber(261773338)
  $core.bool hasRefreshschedule() => $_has(2);
  @$pb.TagNumber(261773338)
  void clearRefreshschedule() => $_clearField(261773338);
  @$pb.TagNumber(261773338)
  RefreshSchedule ensureRefreshschedule() => $_ensure(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  DashboardType get type => $_getN(4);
  @$pb.TagNumber(290836590)
  set type(DashboardType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(4);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(5);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(5, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(5);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(6);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(6);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(501826147)
  $pb.PbList<Widget> get widgets => $_getList(7);
}

class UpdateEventDataStoreRequest extends $pb.GeneratedMessage {
  factory UpdateEventDataStoreRequest({
    $core.bool? multiregionenabled,
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? kmskeyid,
    $core.String? eventdatastore,
    BillingMode? billingmode,
    $core.int? retentionperiod,
    $core.String? name,
    $core.bool? terminationprotectionenabled,
    $core.bool? organizationenabled,
  }) {
    final result = create();
    if (multiregionenabled != null)
      result.multiregionenabled = multiregionenabled;
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (eventdatastore != null) result.eventdatastore = eventdatastore;
    if (billingmode != null) result.billingmode = billingmode;
    if (retentionperiod != null) result.retentionperiod = retentionperiod;
    if (name != null) result.name = name;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (organizationenabled != null)
      result.organizationenabled = organizationenabled;
    return result;
  }

  UpdateEventDataStoreRequest._();

  factory UpdateEventDataStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateEventDataStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateEventDataStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOB(20620620, _omitFieldNames ? '' : 'multiregionenabled')
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(136801729, _omitFieldNames ? '' : 'eventdatastore')
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aI(196383721, _omitFieldNames ? '' : 'retentionperiod')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOB(480171176, _omitFieldNames ? '' : 'organizationenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateEventDataStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateEventDataStoreRequest copyWith(
          void Function(UpdateEventDataStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateEventDataStoreRequest))
          as UpdateEventDataStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateEventDataStoreRequest create() =>
      UpdateEventDataStoreRequest._();
  @$core.override
  UpdateEventDataStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateEventDataStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateEventDataStoreRequest>(create);
  static UpdateEventDataStoreRequest? _defaultInstance;

  @$pb.TagNumber(20620620)
  $core.bool get multiregionenabled => $_getBF(0);
  @$pb.TagNumber(20620620)
  set multiregionenabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(20620620)
  $core.bool hasMultiregionenabled() => $_has(0);
  @$pb.TagNumber(20620620)
  void clearMultiregionenabled() => $_clearField(20620620);

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(1);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(2);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(2);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(136801729)
  $core.String get eventdatastore => $_getSZ(3);
  @$pb.TagNumber(136801729)
  set eventdatastore($core.String value) => $_setString(3, value);
  @$pb.TagNumber(136801729)
  $core.bool hasEventdatastore() => $_has(3);
  @$pb.TagNumber(136801729)
  void clearEventdatastore() => $_clearField(136801729);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(4);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(4);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(196383721)
  $core.int get retentionperiod => $_getIZ(5);
  @$pb.TagNumber(196383721)
  set retentionperiod($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(196383721)
  $core.bool hasRetentionperiod() => $_has(5);
  @$pb.TagNumber(196383721)
  void clearRetentionperiod() => $_clearField(196383721);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(7);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(7);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(480171176)
  $core.bool get organizationenabled => $_getBF(8);
  @$pb.TagNumber(480171176)
  set organizationenabled($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(480171176)
  $core.bool hasOrganizationenabled() => $_has(8);
  @$pb.TagNumber(480171176)
  void clearOrganizationenabled() => $_clearField(480171176);
}

class UpdateEventDataStoreResponse extends $pb.GeneratedMessage {
  factory UpdateEventDataStoreResponse({
    EventDataStoreStatus? status,
    $core.bool? multiregionenabled,
    $core.Iterable<AdvancedEventSelector>? advancedeventselectors,
    $core.String? updatedtimestamp,
    $core.String? kmskeyid,
    FederationStatus? federationstatus,
    BillingMode? billingmode,
    $core.int? retentionperiod,
    $core.String? name,
    $core.String? eventdatastorearn,
    $core.String? createdtimestamp,
    $core.bool? terminationprotectionenabled,
    $core.bool? organizationenabled,
    $core.String? federationrolearn,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (multiregionenabled != null)
      result.multiregionenabled = multiregionenabled;
    if (advancedeventselectors != null)
      result.advancedeventselectors.addAll(advancedeventselectors);
    if (updatedtimestamp != null) result.updatedtimestamp = updatedtimestamp;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (federationstatus != null) result.federationstatus = federationstatus;
    if (billingmode != null) result.billingmode = billingmode;
    if (retentionperiod != null) result.retentionperiod = retentionperiod;
    if (name != null) result.name = name;
    if (eventdatastorearn != null) result.eventdatastorearn = eventdatastorearn;
    if (createdtimestamp != null) result.createdtimestamp = createdtimestamp;
    if (terminationprotectionenabled != null)
      result.terminationprotectionenabled = terminationprotectionenabled;
    if (organizationenabled != null)
      result.organizationenabled = organizationenabled;
    if (federationrolearn != null) result.federationrolearn = federationrolearn;
    return result;
  }

  UpdateEventDataStoreResponse._();

  factory UpdateEventDataStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateEventDataStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateEventDataStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aE<EventDataStoreStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: EventDataStoreStatus.values)
    ..aOB(20620620, _omitFieldNames ? '' : 'multiregionenabled')
    ..pPM<AdvancedEventSelector>(
        36838194, _omitFieldNames ? '' : 'advancedeventselectors',
        subBuilder: AdvancedEventSelector.create)
    ..aOS(44364161, _omitFieldNames ? '' : 'updatedtimestamp')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aE<FederationStatus>(146235383, _omitFieldNames ? '' : 'federationstatus',
        enumValues: FederationStatus.values)
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aI(196383721, _omitFieldNames ? '' : 'retentionperiod')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(331732456, _omitFieldNames ? '' : 'eventdatastorearn')
    ..aOS(334753274, _omitFieldNames ? '' : 'createdtimestamp')
    ..aOB(376863196, _omitFieldNames ? '' : 'terminationprotectionenabled')
    ..aOB(480171176, _omitFieldNames ? '' : 'organizationenabled')
    ..aOS(504464364, _omitFieldNames ? '' : 'federationrolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateEventDataStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateEventDataStoreResponse copyWith(
          void Function(UpdateEventDataStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateEventDataStoreResponse))
          as UpdateEventDataStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateEventDataStoreResponse create() =>
      UpdateEventDataStoreResponse._();
  @$core.override
  UpdateEventDataStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateEventDataStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateEventDataStoreResponse>(create);
  static UpdateEventDataStoreResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  EventDataStoreStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(EventDataStoreStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(20620620)
  $core.bool get multiregionenabled => $_getBF(1);
  @$pb.TagNumber(20620620)
  set multiregionenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(20620620)
  $core.bool hasMultiregionenabled() => $_has(1);
  @$pb.TagNumber(20620620)
  void clearMultiregionenabled() => $_clearField(20620620);

  @$pb.TagNumber(36838194)
  $pb.PbList<AdvancedEventSelector> get advancedeventselectors => $_getList(2);

  @$pb.TagNumber(44364161)
  $core.String get updatedtimestamp => $_getSZ(3);
  @$pb.TagNumber(44364161)
  set updatedtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(44364161)
  $core.bool hasUpdatedtimestamp() => $_has(3);
  @$pb.TagNumber(44364161)
  void clearUpdatedtimestamp() => $_clearField(44364161);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(4);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(4);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(146235383)
  FederationStatus get federationstatus => $_getN(5);
  @$pb.TagNumber(146235383)
  set federationstatus(FederationStatus value) => $_setField(146235383, value);
  @$pb.TagNumber(146235383)
  $core.bool hasFederationstatus() => $_has(5);
  @$pb.TagNumber(146235383)
  void clearFederationstatus() => $_clearField(146235383);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(6);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(6);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(196383721)
  $core.int get retentionperiod => $_getIZ(7);
  @$pb.TagNumber(196383721)
  set retentionperiod($core.int value) => $_setSignedInt32(7, value);
  @$pb.TagNumber(196383721)
  $core.bool hasRetentionperiod() => $_has(7);
  @$pb.TagNumber(196383721)
  void clearRetentionperiod() => $_clearField(196383721);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(8);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(8, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(8);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(331732456)
  $core.String get eventdatastorearn => $_getSZ(9);
  @$pb.TagNumber(331732456)
  set eventdatastorearn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(331732456)
  $core.bool hasEventdatastorearn() => $_has(9);
  @$pb.TagNumber(331732456)
  void clearEventdatastorearn() => $_clearField(331732456);

  @$pb.TagNumber(334753274)
  $core.String get createdtimestamp => $_getSZ(10);
  @$pb.TagNumber(334753274)
  set createdtimestamp($core.String value) => $_setString(10, value);
  @$pb.TagNumber(334753274)
  $core.bool hasCreatedtimestamp() => $_has(10);
  @$pb.TagNumber(334753274)
  void clearCreatedtimestamp() => $_clearField(334753274);

  @$pb.TagNumber(376863196)
  $core.bool get terminationprotectionenabled => $_getBF(11);
  @$pb.TagNumber(376863196)
  set terminationprotectionenabled($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(376863196)
  $core.bool hasTerminationprotectionenabled() => $_has(11);
  @$pb.TagNumber(376863196)
  void clearTerminationprotectionenabled() => $_clearField(376863196);

  @$pb.TagNumber(480171176)
  $core.bool get organizationenabled => $_getBF(12);
  @$pb.TagNumber(480171176)
  set organizationenabled($core.bool value) => $_setBool(12, value);
  @$pb.TagNumber(480171176)
  $core.bool hasOrganizationenabled() => $_has(12);
  @$pb.TagNumber(480171176)
  void clearOrganizationenabled() => $_clearField(480171176);

  @$pb.TagNumber(504464364)
  $core.String get federationrolearn => $_getSZ(13);
  @$pb.TagNumber(504464364)
  set federationrolearn($core.String value) => $_setString(13, value);
  @$pb.TagNumber(504464364)
  $core.bool hasFederationrolearn() => $_has(13);
  @$pb.TagNumber(504464364)
  void clearFederationrolearn() => $_clearField(504464364);
}

class UpdateTrailRequest extends $pb.GeneratedMessage {
  factory UpdateTrailRequest({
    $core.String? kmskeyid,
    $core.String? cloudwatchlogsrolearn,
    $core.bool? enablelogfilevalidation,
    $core.bool? isorganizationtrail,
    $core.String? s3keyprefix,
    $core.bool? includeglobalserviceevents,
    $core.String? name,
    $core.String? s3bucketname,
    $core.String? snstopicname,
    $core.bool? ismultiregiontrail,
    $core.String? cloudwatchlogsloggrouparn,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (cloudwatchlogsrolearn != null)
      result.cloudwatchlogsrolearn = cloudwatchlogsrolearn;
    if (enablelogfilevalidation != null)
      result.enablelogfilevalidation = enablelogfilevalidation;
    if (isorganizationtrail != null)
      result.isorganizationtrail = isorganizationtrail;
    if (s3keyprefix != null) result.s3keyprefix = s3keyprefix;
    if (includeglobalserviceevents != null)
      result.includeglobalserviceevents = includeglobalserviceevents;
    if (name != null) result.name = name;
    if (s3bucketname != null) result.s3bucketname = s3bucketname;
    if (snstopicname != null) result.snstopicname = snstopicname;
    if (ismultiregiontrail != null)
      result.ismultiregiontrail = ismultiregiontrail;
    if (cloudwatchlogsloggrouparn != null)
      result.cloudwatchlogsloggrouparn = cloudwatchlogsloggrouparn;
    return result;
  }

  UpdateTrailRequest._();

  factory UpdateTrailRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTrailRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTrailRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(55454690, _omitFieldNames ? '' : 'cloudwatchlogsrolearn')
    ..aOB(80236930, _omitFieldNames ? '' : 'enablelogfilevalidation')
    ..aOB(145256127, _omitFieldNames ? '' : 'isorganizationtrail')
    ..aOS(206015359, _omitFieldNames ? '' : 's3keyprefix')
    ..aOB(212227451, _omitFieldNames ? '' : 'includeglobalserviceevents')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(320495427, _omitFieldNames ? '' : 's3bucketname')
    ..aOS(415454800, _omitFieldNames ? '' : 'snstopicname')
    ..aOB(468658571, _omitFieldNames ? '' : 'ismultiregiontrail')
    ..aOS(519613355, _omitFieldNames ? '' : 'cloudwatchlogsloggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrailRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrailRequest copyWith(void Function(UpdateTrailRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateTrailRequest))
          as UpdateTrailRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTrailRequest create() => UpdateTrailRequest._();
  @$core.override
  UpdateTrailRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTrailRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTrailRequest>(create);
  static UpdateTrailRequest? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(55454690)
  $core.String get cloudwatchlogsrolearn => $_getSZ(1);
  @$pb.TagNumber(55454690)
  set cloudwatchlogsrolearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(55454690)
  $core.bool hasCloudwatchlogsrolearn() => $_has(1);
  @$pb.TagNumber(55454690)
  void clearCloudwatchlogsrolearn() => $_clearField(55454690);

  @$pb.TagNumber(80236930)
  $core.bool get enablelogfilevalidation => $_getBF(2);
  @$pb.TagNumber(80236930)
  set enablelogfilevalidation($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(80236930)
  $core.bool hasEnablelogfilevalidation() => $_has(2);
  @$pb.TagNumber(80236930)
  void clearEnablelogfilevalidation() => $_clearField(80236930);

  @$pb.TagNumber(145256127)
  $core.bool get isorganizationtrail => $_getBF(3);
  @$pb.TagNumber(145256127)
  set isorganizationtrail($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(145256127)
  $core.bool hasIsorganizationtrail() => $_has(3);
  @$pb.TagNumber(145256127)
  void clearIsorganizationtrail() => $_clearField(145256127);

  @$pb.TagNumber(206015359)
  $core.String get s3keyprefix => $_getSZ(4);
  @$pb.TagNumber(206015359)
  set s3keyprefix($core.String value) => $_setString(4, value);
  @$pb.TagNumber(206015359)
  $core.bool hasS3keyprefix() => $_has(4);
  @$pb.TagNumber(206015359)
  void clearS3keyprefix() => $_clearField(206015359);

  @$pb.TagNumber(212227451)
  $core.bool get includeglobalserviceevents => $_getBF(5);
  @$pb.TagNumber(212227451)
  set includeglobalserviceevents($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(212227451)
  $core.bool hasIncludeglobalserviceevents() => $_has(5);
  @$pb.TagNumber(212227451)
  void clearIncludeglobalserviceevents() => $_clearField(212227451);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(320495427)
  $core.String get s3bucketname => $_getSZ(7);
  @$pb.TagNumber(320495427)
  set s3bucketname($core.String value) => $_setString(7, value);
  @$pb.TagNumber(320495427)
  $core.bool hasS3bucketname() => $_has(7);
  @$pb.TagNumber(320495427)
  void clearS3bucketname() => $_clearField(320495427);

  @$pb.TagNumber(415454800)
  $core.String get snstopicname => $_getSZ(8);
  @$pb.TagNumber(415454800)
  set snstopicname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(415454800)
  $core.bool hasSnstopicname() => $_has(8);
  @$pb.TagNumber(415454800)
  void clearSnstopicname() => $_clearField(415454800);

  @$pb.TagNumber(468658571)
  $core.bool get ismultiregiontrail => $_getBF(9);
  @$pb.TagNumber(468658571)
  set ismultiregiontrail($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(468658571)
  $core.bool hasIsmultiregiontrail() => $_has(9);
  @$pb.TagNumber(468658571)
  void clearIsmultiregiontrail() => $_clearField(468658571);

  @$pb.TagNumber(519613355)
  $core.String get cloudwatchlogsloggrouparn => $_getSZ(10);
  @$pb.TagNumber(519613355)
  set cloudwatchlogsloggrouparn($core.String value) => $_setString(10, value);
  @$pb.TagNumber(519613355)
  $core.bool hasCloudwatchlogsloggrouparn() => $_has(10);
  @$pb.TagNumber(519613355)
  void clearCloudwatchlogsloggrouparn() => $_clearField(519613355);
}

class UpdateTrailResponse extends $pb.GeneratedMessage {
  factory UpdateTrailResponse({
    $core.bool? logfilevalidationenabled,
    $core.String? trailarn,
    $core.String? kmskeyid,
    $core.String? cloudwatchlogsrolearn,
    $core.bool? isorganizationtrail,
    $core.String? s3keyprefix,
    $core.bool? includeglobalserviceevents,
    $core.String? name,
    $core.String? s3bucketname,
    $core.String? snstopicarn,
    $core.String? snstopicname,
    $core.bool? ismultiregiontrail,
    $core.String? cloudwatchlogsloggrouparn,
  }) {
    final result = create();
    if (logfilevalidationenabled != null)
      result.logfilevalidationenabled = logfilevalidationenabled;
    if (trailarn != null) result.trailarn = trailarn;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (cloudwatchlogsrolearn != null)
      result.cloudwatchlogsrolearn = cloudwatchlogsrolearn;
    if (isorganizationtrail != null)
      result.isorganizationtrail = isorganizationtrail;
    if (s3keyprefix != null) result.s3keyprefix = s3keyprefix;
    if (includeglobalserviceevents != null)
      result.includeglobalserviceevents = includeglobalserviceevents;
    if (name != null) result.name = name;
    if (s3bucketname != null) result.s3bucketname = s3bucketname;
    if (snstopicarn != null) result.snstopicarn = snstopicarn;
    if (snstopicname != null) result.snstopicname = snstopicname;
    if (ismultiregiontrail != null)
      result.ismultiregiontrail = ismultiregiontrail;
    if (cloudwatchlogsloggrouparn != null)
      result.cloudwatchlogsloggrouparn = cloudwatchlogsloggrouparn;
    return result;
  }

  UpdateTrailResponse._();

  factory UpdateTrailResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTrailResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTrailResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..aOB(35904346, _omitFieldNames ? '' : 'logfilevalidationenabled')
    ..aOS(39585143, _omitFieldNames ? '' : 'trailarn')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(55454690, _omitFieldNames ? '' : 'cloudwatchlogsrolearn')
    ..aOB(145256127, _omitFieldNames ? '' : 'isorganizationtrail')
    ..aOS(206015359, _omitFieldNames ? '' : 's3keyprefix')
    ..aOB(212227451, _omitFieldNames ? '' : 'includeglobalserviceevents')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(320495427, _omitFieldNames ? '' : 's3bucketname')
    ..aOS(380025580, _omitFieldNames ? '' : 'snstopicarn')
    ..aOS(415454800, _omitFieldNames ? '' : 'snstopicname')
    ..aOB(468658571, _omitFieldNames ? '' : 'ismultiregiontrail')
    ..aOS(519613355, _omitFieldNames ? '' : 'cloudwatchlogsloggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrailResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrailResponse copyWith(void Function(UpdateTrailResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateTrailResponse))
          as UpdateTrailResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTrailResponse create() => UpdateTrailResponse._();
  @$core.override
  UpdateTrailResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTrailResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTrailResponse>(create);
  static UpdateTrailResponse? _defaultInstance;

  @$pb.TagNumber(35904346)
  $core.bool get logfilevalidationenabled => $_getBF(0);
  @$pb.TagNumber(35904346)
  set logfilevalidationenabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(35904346)
  $core.bool hasLogfilevalidationenabled() => $_has(0);
  @$pb.TagNumber(35904346)
  void clearLogfilevalidationenabled() => $_clearField(35904346);

  @$pb.TagNumber(39585143)
  $core.String get trailarn => $_getSZ(1);
  @$pb.TagNumber(39585143)
  set trailarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(39585143)
  $core.bool hasTrailarn() => $_has(1);
  @$pb.TagNumber(39585143)
  void clearTrailarn() => $_clearField(39585143);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(2);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(2);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(55454690)
  $core.String get cloudwatchlogsrolearn => $_getSZ(3);
  @$pb.TagNumber(55454690)
  set cloudwatchlogsrolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(55454690)
  $core.bool hasCloudwatchlogsrolearn() => $_has(3);
  @$pb.TagNumber(55454690)
  void clearCloudwatchlogsrolearn() => $_clearField(55454690);

  @$pb.TagNumber(145256127)
  $core.bool get isorganizationtrail => $_getBF(4);
  @$pb.TagNumber(145256127)
  set isorganizationtrail($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(145256127)
  $core.bool hasIsorganizationtrail() => $_has(4);
  @$pb.TagNumber(145256127)
  void clearIsorganizationtrail() => $_clearField(145256127);

  @$pb.TagNumber(206015359)
  $core.String get s3keyprefix => $_getSZ(5);
  @$pb.TagNumber(206015359)
  set s3keyprefix($core.String value) => $_setString(5, value);
  @$pb.TagNumber(206015359)
  $core.bool hasS3keyprefix() => $_has(5);
  @$pb.TagNumber(206015359)
  void clearS3keyprefix() => $_clearField(206015359);

  @$pb.TagNumber(212227451)
  $core.bool get includeglobalserviceevents => $_getBF(6);
  @$pb.TagNumber(212227451)
  set includeglobalserviceevents($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(212227451)
  $core.bool hasIncludeglobalserviceevents() => $_has(6);
  @$pb.TagNumber(212227451)
  void clearIncludeglobalserviceevents() => $_clearField(212227451);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(320495427)
  $core.String get s3bucketname => $_getSZ(8);
  @$pb.TagNumber(320495427)
  set s3bucketname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(320495427)
  $core.bool hasS3bucketname() => $_has(8);
  @$pb.TagNumber(320495427)
  void clearS3bucketname() => $_clearField(320495427);

  @$pb.TagNumber(380025580)
  $core.String get snstopicarn => $_getSZ(9);
  @$pb.TagNumber(380025580)
  set snstopicarn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(380025580)
  $core.bool hasSnstopicarn() => $_has(9);
  @$pb.TagNumber(380025580)
  void clearSnstopicarn() => $_clearField(380025580);

  @$pb.TagNumber(415454800)
  $core.String get snstopicname => $_getSZ(10);
  @$pb.TagNumber(415454800)
  set snstopicname($core.String value) => $_setString(10, value);
  @$pb.TagNumber(415454800)
  $core.bool hasSnstopicname() => $_has(10);
  @$pb.TagNumber(415454800)
  void clearSnstopicname() => $_clearField(415454800);

  @$pb.TagNumber(468658571)
  $core.bool get ismultiregiontrail => $_getBF(11);
  @$pb.TagNumber(468658571)
  set ismultiregiontrail($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(468658571)
  $core.bool hasIsmultiregiontrail() => $_has(11);
  @$pb.TagNumber(468658571)
  void clearIsmultiregiontrail() => $_clearField(468658571);

  @$pb.TagNumber(519613355)
  $core.String get cloudwatchlogsloggrouparn => $_getSZ(12);
  @$pb.TagNumber(519613355)
  set cloudwatchlogsloggrouparn($core.String value) => $_setString(12, value);
  @$pb.TagNumber(519613355)
  $core.bool hasCloudwatchlogsloggrouparn() => $_has(12);
  @$pb.TagNumber(519613355)
  void clearCloudwatchlogsloggrouparn() => $_clearField(519613355);
}

class Widget extends $pb.GeneratedMessage {
  factory Widget({
    $core.Iterable<$core.String>? queryparameters,
    $core.String? querystatement,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? viewproperties,
    $core.String? queryalias,
  }) {
    final result = create();
    if (queryparameters != null) result.queryparameters.addAll(queryparameters);
    if (querystatement != null) result.querystatement = querystatement;
    if (viewproperties != null)
      result.viewproperties.addEntries(viewproperties);
    if (queryalias != null) result.queryalias = queryalias;
    return result;
  }

  Widget._();

  factory Widget.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Widget.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Widget',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'cloudtrail'),
      createEmptyInstance: create)
    ..pPS(76317660, _omitFieldNames ? '' : 'queryparameters')
    ..aOS(340852217, _omitFieldNames ? '' : 'querystatement')
    ..m<$core.String, $core.String>(
        381974398, _omitFieldNames ? '' : 'viewproperties',
        entryClassName: 'Widget.ViewpropertiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cloudtrail'))
    ..aOS(512762270, _omitFieldNames ? '' : 'queryalias')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Widget clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Widget copyWith(void Function(Widget) updates) =>
      super.copyWith((message) => updates(message as Widget)) as Widget;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Widget create() => Widget._();
  @$core.override
  Widget createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Widget getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Widget>(create);
  static Widget? _defaultInstance;

  @$pb.TagNumber(76317660)
  $pb.PbList<$core.String> get queryparameters => $_getList(0);

  @$pb.TagNumber(340852217)
  $core.String get querystatement => $_getSZ(1);
  @$pb.TagNumber(340852217)
  set querystatement($core.String value) => $_setString(1, value);
  @$pb.TagNumber(340852217)
  $core.bool hasQuerystatement() => $_has(1);
  @$pb.TagNumber(340852217)
  void clearQuerystatement() => $_clearField(340852217);

  @$pb.TagNumber(381974398)
  $pb.PbMap<$core.String, $core.String> get viewproperties => $_getMap(2);

  @$pb.TagNumber(512762270)
  $core.String get queryalias => $_getSZ(3);
  @$pb.TagNumber(512762270)
  set queryalias($core.String value) => $_setString(3, value);
  @$pb.TagNumber(512762270)
  $core.bool hasQueryalias() => $_has(3);
  @$pb.TagNumber(512762270)
  void clearQueryalias() => $_clearField(512762270);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
