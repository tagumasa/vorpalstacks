// This is a generated file - do not edit.
//
// Generated from admin_config.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class ConfigEntry extends $pb.GeneratedMessage {
  factory ConfigEntry({
    $core.String? key,
    $core.String? value,
    $core.String? type,
    $core.String? source,
    $core.String? description,
    $core.bool? readOnly,
    $fixnum.Int64? updatedAt,
    $core.String? envVar,
    $core.String? category,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    if (type != null) result.type = type;
    if (source != null) result.source = source;
    if (description != null) result.description = description;
    if (readOnly != null) result.readOnly = readOnly;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (envVar != null) result.envVar = envVar;
    if (category != null) result.category = category;
    return result;
  }

  ConfigEntry._();

  factory ConfigEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConfigEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConfigEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'key')
    ..aOS(2, _omitFieldNames ? '' : 'value')
    ..aOS(3, _omitFieldNames ? '' : 'type')
    ..aOS(4, _omitFieldNames ? '' : 'source')
    ..aOS(5, _omitFieldNames ? '' : 'description')
    ..aOB(6, _omitFieldNames ? '' : 'readOnly')
    ..aInt64(7, _omitFieldNames ? '' : 'updatedAt')
    ..aOS(8, _omitFieldNames ? '' : 'envVar')
    ..aOS(9, _omitFieldNames ? '' : 'category')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfigEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfigEntry copyWith(void Function(ConfigEntry) updates) =>
      super.copyWith((message) => updates(message as ConfigEntry))
          as ConfigEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConfigEntry create() => ConfigEntry._();
  @$core.override
  ConfigEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConfigEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConfigEntry>(create);
  static ConfigEntry? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearKey() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(2)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(2)
  void clearValue() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get type => $_getSZ(2);
  @$pb.TagNumber(3)
  set type($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(3)
  void clearType() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get source => $_getSZ(3);
  @$pb.TagNumber(4)
  set source($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSource() => $_has(3);
  @$pb.TagNumber(4)
  void clearSource() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(5)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(5)
  void clearDescription() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.bool get readOnly => $_getBF(5);
  @$pb.TagNumber(6)
  set readOnly($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(6)
  $core.bool hasReadOnly() => $_has(5);
  @$pb.TagNumber(6)
  void clearReadOnly() => $_clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get updatedAt => $_getI64(6);
  @$pb.TagNumber(7)
  set updatedAt($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasUpdatedAt() => $_has(6);
  @$pb.TagNumber(7)
  void clearUpdatedAt() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get envVar => $_getSZ(7);
  @$pb.TagNumber(8)
  set envVar($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasEnvVar() => $_has(7);
  @$pb.TagNumber(8)
  void clearEnvVar() => $_clearField(8);

  @$pb.TagNumber(9)
  $core.String get category => $_getSZ(8);
  @$pb.TagNumber(9)
  set category($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasCategory() => $_has(8);
  @$pb.TagNumber(9)
  void clearCategory() => $_clearField(9);
}

class GetConfigRequest extends $pb.GeneratedMessage {
  factory GetConfigRequest({
    $core.String? key,
  }) {
    final result = create();
    if (key != null) result.key = key;
    return result;
  }

  GetConfigRequest._();

  factory GetConfigRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetConfigRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetConfigRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'key')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConfigRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConfigRequest copyWith(void Function(GetConfigRequest) updates) =>
      super.copyWith((message) => updates(message as GetConfigRequest))
          as GetConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetConfigRequest create() => GetConfigRequest._();
  @$core.override
  GetConfigRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetConfigRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetConfigRequest>(create);
  static GetConfigRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearKey() => $_clearField(1);
}

class GetConfigResponse extends $pb.GeneratedMessage {
  factory GetConfigResponse({
    ConfigEntry? entry,
  }) {
    final result = create();
    if (entry != null) result.entry = entry;
    return result;
  }

  GetConfigResponse._();

  factory GetConfigResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetConfigResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetConfigResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOM<ConfigEntry>(1, _omitFieldNames ? '' : 'entry',
        subBuilder: ConfigEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConfigResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConfigResponse copyWith(void Function(GetConfigResponse) updates) =>
      super.copyWith((message) => updates(message as GetConfigResponse))
          as GetConfigResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetConfigResponse create() => GetConfigResponse._();
  @$core.override
  GetConfigResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetConfigResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetConfigResponse>(create);
  static GetConfigResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ConfigEntry get entry => $_getN(0);
  @$pb.TagNumber(1)
  set entry(ConfigEntry value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasEntry() => $_has(0);
  @$pb.TagNumber(1)
  void clearEntry() => $_clearField(1);
  @$pb.TagNumber(1)
  ConfigEntry ensureEntry() => $_ensure(0);
}

class ListConfigRequest extends $pb.GeneratedMessage {
  factory ListConfigRequest({
    $core.String? category,
  }) {
    final result = create();
    if (category != null) result.category = category;
    return result;
  }

  ListConfigRequest._();

  factory ListConfigRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListConfigRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListConfigRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'category')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConfigRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConfigRequest copyWith(void Function(ListConfigRequest) updates) =>
      super.copyWith((message) => updates(message as ListConfigRequest))
          as ListConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListConfigRequest create() => ListConfigRequest._();
  @$core.override
  ListConfigRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListConfigRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListConfigRequest>(create);
  static ListConfigRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get category => $_getSZ(0);
  @$pb.TagNumber(1)
  set category($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCategory() => $_has(0);
  @$pb.TagNumber(1)
  void clearCategory() => $_clearField(1);
}

class ListConfigResponse extends $pb.GeneratedMessage {
  factory ListConfigResponse({
    $core.Iterable<ConfigEntry>? entries,
  }) {
    final result = create();
    if (entries != null) result.entries.addAll(entries);
    return result;
  }

  ListConfigResponse._();

  factory ListConfigResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListConfigResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListConfigResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..pPM<ConfigEntry>(1, _omitFieldNames ? '' : 'entries',
        subBuilder: ConfigEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConfigResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConfigResponse copyWith(void Function(ListConfigResponse) updates) =>
      super.copyWith((message) => updates(message as ListConfigResponse))
          as ListConfigResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListConfigResponse create() => ListConfigResponse._();
  @$core.override
  ListConfigResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListConfigResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListConfigResponse>(create);
  static ListConfigResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<ConfigEntry> get entries => $_getList(0);
}

class UpdateConfigRequest extends $pb.GeneratedMessage {
  factory UpdateConfigRequest({
    $core.String? key,
    $core.String? value,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    return result;
  }

  UpdateConfigRequest._();

  factory UpdateConfigRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConfigRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConfigRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'key')
    ..aOS(2, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConfigRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConfigRequest copyWith(void Function(UpdateConfigRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateConfigRequest))
          as UpdateConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConfigRequest create() => UpdateConfigRequest._();
  @$core.override
  UpdateConfigRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateConfigRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateConfigRequest>(create);
  static UpdateConfigRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearKey() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(2)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(2)
  void clearValue() => $_clearField(2);
}

class ResetConfigRequest extends $pb.GeneratedMessage {
  factory ResetConfigRequest({
    $core.String? key,
  }) {
    final result = create();
    if (key != null) result.key = key;
    return result;
  }

  ResetConfigRequest._();

  factory ResetConfigRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResetConfigRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResetConfigRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'key')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResetConfigRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResetConfigRequest copyWith(void Function(ResetConfigRequest) updates) =>
      super.copyWith((message) => updates(message as ResetConfigRequest))
          as ResetConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResetConfigRequest create() => ResetConfigRequest._();
  @$core.override
  ResetConfigRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResetConfigRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResetConfigRequest>(create);
  static ResetConfigRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearKey() => $_clearField(1);
}

class ListServicesRequest extends $pb.GeneratedMessage {
  factory ListServicesRequest() => create();

  ListServicesRequest._();

  factory ListServicesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListServicesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListServicesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListServicesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListServicesRequest copyWith(void Function(ListServicesRequest) updates) =>
      super.copyWith((message) => updates(message as ListServicesRequest))
          as ListServicesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListServicesRequest create() => ListServicesRequest._();
  @$core.override
  ListServicesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListServicesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListServicesRequest>(create);
  static ListServicesRequest? _defaultInstance;
}

class ListServicesResponse extends $pb.GeneratedMessage {
  factory ListServicesResponse({
    $core.Iterable<ServiceInfo>? services,
  }) {
    final result = create();
    if (services != null) result.services.addAll(services);
    return result;
  }

  ListServicesResponse._();

  factory ListServicesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListServicesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListServicesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..pPM<ServiceInfo>(1, _omitFieldNames ? '' : 'services',
        subBuilder: ServiceInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListServicesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListServicesResponse copyWith(void Function(ListServicesResponse) updates) =>
      super.copyWith((message) => updates(message as ListServicesResponse))
          as ListServicesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListServicesResponse create() => ListServicesResponse._();
  @$core.override
  ListServicesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListServicesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListServicesResponse>(create);
  static ListServicesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<ServiceInfo> get services => $_getList(0);
}

class ServiceInfo extends $pb.GeneratedMessage {
  factory ServiceInfo({
    $core.String? name,
    $core.bool? enabled,
    $core.int? resourceCount,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (enabled != null) result.enabled = enabled;
    if (resourceCount != null) result.resourceCount = resourceCount;
    return result;
  }

  ServiceInfo._();

  factory ServiceInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ServiceInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ServiceInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOB(2, _omitFieldNames ? '' : 'enabled')
    ..aI(3, _omitFieldNames ? '' : 'resourceCount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceInfo copyWith(void Function(ServiceInfo) updates) =>
      super.copyWith((message) => updates(message as ServiceInfo))
          as ServiceInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ServiceInfo create() => ServiceInfo._();
  @$core.override
  ServiceInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ServiceInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ServiceInfo>(create);
  static ServiceInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(2)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(2)
  void clearEnabled() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get resourceCount => $_getIZ(2);
  @$pb.TagNumber(3)
  set resourceCount($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasResourceCount() => $_has(2);
  @$pb.TagNumber(3)
  void clearResourceCount() => $_clearField(3);
}

class GetServiceStatusRequest extends $pb.GeneratedMessage {
  factory GetServiceStatusRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  GetServiceStatusRequest._();

  factory GetServiceStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetServiceStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetServiceStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetServiceStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetServiceStatusRequest copyWith(
          void Function(GetServiceStatusRequest) updates) =>
      super.copyWith((message) => updates(message as GetServiceStatusRequest))
          as GetServiceStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetServiceStatusRequest create() => GetServiceStatusRequest._();
  @$core.override
  GetServiceStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetServiceStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetServiceStatusRequest>(create);
  static GetServiceStatusRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);
}

class ServiceStatus extends $pb.GeneratedMessage {
  factory ServiceStatus({
    $core.String? name,
    $core.bool? enabled,
    $core.String? mode,
    $core.int? resourceCount,
    $fixnum.Int64? lastRequestAt,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (enabled != null) result.enabled = enabled;
    if (mode != null) result.mode = mode;
    if (resourceCount != null) result.resourceCount = resourceCount;
    if (lastRequestAt != null) result.lastRequestAt = lastRequestAt;
    return result;
  }

  ServiceStatus._();

  factory ServiceStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ServiceStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ServiceStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOB(2, _omitFieldNames ? '' : 'enabled')
    ..aOS(3, _omitFieldNames ? '' : 'mode')
    ..aI(4, _omitFieldNames ? '' : 'resourceCount')
    ..aInt64(5, _omitFieldNames ? '' : 'lastRequestAt')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceStatus copyWith(void Function(ServiceStatus) updates) =>
      super.copyWith((message) => updates(message as ServiceStatus))
          as ServiceStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ServiceStatus create() => ServiceStatus._();
  @$core.override
  ServiceStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ServiceStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ServiceStatus>(create);
  static ServiceStatus? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(2)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(2)
  void clearEnabled() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get mode => $_getSZ(2);
  @$pb.TagNumber(3)
  set mode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMode() => $_has(2);
  @$pb.TagNumber(3)
  void clearMode() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.int get resourceCount => $_getIZ(3);
  @$pb.TagNumber(4)
  set resourceCount($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasResourceCount() => $_has(3);
  @$pb.TagNumber(4)
  void clearResourceCount() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get lastRequestAt => $_getI64(4);
  @$pb.TagNumber(5)
  set lastRequestAt($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasLastRequestAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearLastRequestAt() => $_clearField(5);
}

class GetResourcePortRequest extends $pb.GeneratedMessage {
  factory GetResourcePortRequest({
    $core.String? servicePortKey,
    $core.String? resourceId,
  }) {
    final result = create();
    if (servicePortKey != null) result.servicePortKey = servicePortKey;
    if (resourceId != null) result.resourceId = resourceId;
    return result;
  }

  GetResourcePortRequest._();

  factory GetResourcePortRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourcePortRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourcePortRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'servicePortKey')
    ..aOS(2, _omitFieldNames ? '' : 'resourceId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePortRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePortRequest copyWith(
          void Function(GetResourcePortRequest) updates) =>
      super.copyWith((message) => updates(message as GetResourcePortRequest))
          as GetResourcePortRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourcePortRequest create() => GetResourcePortRequest._();
  @$core.override
  GetResourcePortRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourcePortRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourcePortRequest>(create);
  static GetResourcePortRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get servicePortKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set servicePortKey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasServicePortKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearServicePortKey() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get resourceId => $_getSZ(1);
  @$pb.TagNumber(2)
  set resourceId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasResourceId() => $_has(1);
  @$pb.TagNumber(2)
  void clearResourceId() => $_clearField(2);
}

class GetResourcePortResponse extends $pb.GeneratedMessage {
  factory GetResourcePortResponse({
    $core.int? port,
    $core.String? source,
  }) {
    final result = create();
    if (port != null) result.port = port;
    if (source != null) result.source = source;
    return result;
  }

  GetResourcePortResponse._();

  factory GetResourcePortResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourcePortResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourcePortResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'port')
    ..aOS(2, _omitFieldNames ? '' : 'source')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePortResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePortResponse copyWith(
          void Function(GetResourcePortResponse) updates) =>
      super.copyWith((message) => updates(message as GetResourcePortResponse))
          as GetResourcePortResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourcePortResponse create() => GetResourcePortResponse._();
  @$core.override
  GetResourcePortResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourcePortResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourcePortResponse>(create);
  static GetResourcePortResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get port => $_getIZ(0);
  @$pb.TagNumber(1)
  set port($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPort() => $_has(0);
  @$pb.TagNumber(1)
  void clearPort() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get source => $_getSZ(1);
  @$pb.TagNumber(2)
  set source($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSource() => $_has(1);
  @$pb.TagNumber(2)
  void clearSource() => $_clearField(2);
}

class SetResourcePortRequest extends $pb.GeneratedMessage {
  factory SetResourcePortRequest({
    $core.String? servicePortKey,
    $core.String? resourceId,
    $core.int? port,
  }) {
    final result = create();
    if (servicePortKey != null) result.servicePortKey = servicePortKey;
    if (resourceId != null) result.resourceId = resourceId;
    if (port != null) result.port = port;
    return result;
  }

  SetResourcePortRequest._();

  factory SetResourcePortRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetResourcePortRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetResourcePortRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'admin_config'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'servicePortKey')
    ..aOS(2, _omitFieldNames ? '' : 'resourceId')
    ..aI(3, _omitFieldNames ? '' : 'port')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetResourcePortRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetResourcePortRequest copyWith(
          void Function(SetResourcePortRequest) updates) =>
      super.copyWith((message) => updates(message as SetResourcePortRequest))
          as SetResourcePortRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetResourcePortRequest create() => SetResourcePortRequest._();
  @$core.override
  SetResourcePortRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetResourcePortRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetResourcePortRequest>(create);
  static SetResourcePortRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get servicePortKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set servicePortKey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasServicePortKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearServicePortKey() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get resourceId => $_getSZ(1);
  @$pb.TagNumber(2)
  set resourceId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasResourceId() => $_has(1);
  @$pb.TagNumber(2)
  void clearResourceId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get port => $_getIZ(2);
  @$pb.TagNumber(3)
  set port($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPort() => $_has(2);
  @$pb.TagNumber(3)
  void clearPort() => $_clearField(3);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
