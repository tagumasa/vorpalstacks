// This is a generated file - do not edit.
//
// Generated from athena.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'athena.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'athena.pbenum.dart';

class AclConfiguration extends $pb.GeneratedMessage {
  factory AclConfiguration({
    S3AclOption? s3acloption,
  }) {
    final result = create();
    if (s3acloption != null) result.s3acloption = s3acloption;
    return result;
  }

  AclConfiguration._();

  factory AclConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AclConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AclConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<S3AclOption>(97327001, _omitFieldNames ? '' : 's3acloption',
        enumValues: S3AclOption.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AclConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AclConfiguration copyWith(void Function(AclConfiguration) updates) =>
      super.copyWith((message) => updates(message as AclConfiguration))
          as AclConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AclConfiguration create() => AclConfiguration._();
  @$core.override
  AclConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AclConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AclConfiguration>(create);
  static AclConfiguration? _defaultInstance;

  @$pb.TagNumber(97327001)
  S3AclOption get s3acloption => $_getN(0);
  @$pb.TagNumber(97327001)
  set s3acloption(S3AclOption value) => $_setField(97327001, value);
  @$pb.TagNumber(97327001)
  $core.bool hasS3acloption() => $_has(0);
  @$pb.TagNumber(97327001)
  void clearS3acloption() => $_clearField(97327001);
}

class ApplicationDPUSizes extends $pb.GeneratedMessage {
  factory ApplicationDPUSizes({
    $core.Iterable<$core.int>? supporteddpusizes,
    $core.String? applicationruntimeid,
  }) {
    final result = create();
    if (supporteddpusizes != null)
      result.supporteddpusizes.addAll(supporteddpusizes);
    if (applicationruntimeid != null)
      result.applicationruntimeid = applicationruntimeid;
    return result;
  }

  ApplicationDPUSizes._();

  factory ApplicationDPUSizes.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApplicationDPUSizes.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApplicationDPUSizes',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..p<$core.int>(227874103, _omitFieldNames ? '' : 'supporteddpusizes',
        $pb.PbFieldType.K3)
    ..aOS(300478599, _omitFieldNames ? '' : 'applicationruntimeid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplicationDPUSizes clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplicationDPUSizes copyWith(void Function(ApplicationDPUSizes) updates) =>
      super.copyWith((message) => updates(message as ApplicationDPUSizes))
          as ApplicationDPUSizes;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApplicationDPUSizes create() => ApplicationDPUSizes._();
  @$core.override
  ApplicationDPUSizes createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApplicationDPUSizes getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ApplicationDPUSizes>(create);
  static ApplicationDPUSizes? _defaultInstance;

  @$pb.TagNumber(227874103)
  $pb.PbList<$core.int> get supporteddpusizes => $_getList(0);

  @$pb.TagNumber(300478599)
  $core.String get applicationruntimeid => $_getSZ(1);
  @$pb.TagNumber(300478599)
  set applicationruntimeid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(300478599)
  $core.bool hasApplicationruntimeid() => $_has(1);
  @$pb.TagNumber(300478599)
  void clearApplicationruntimeid() => $_clearField(300478599);
}

class AthenaError extends $pb.GeneratedMessage {
  factory AthenaError({
    $core.bool? retryable,
    $core.int? errorcategory,
    $core.int? errortype,
    $core.String? errormessage,
  }) {
    final result = create();
    if (retryable != null) result.retryable = retryable;
    if (errorcategory != null) result.errorcategory = errorcategory;
    if (errortype != null) result.errortype = errortype;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  AthenaError._();

  factory AthenaError.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AthenaError.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AthenaError',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOB(83386186, _omitFieldNames ? '' : 'retryable')
    ..aI(315958414, _omitFieldNames ? '' : 'errorcategory')
    ..aI(398848954, _omitFieldNames ? '' : 'errortype')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AthenaError clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AthenaError copyWith(void Function(AthenaError) updates) =>
      super.copyWith((message) => updates(message as AthenaError))
          as AthenaError;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AthenaError create() => AthenaError._();
  @$core.override
  AthenaError createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AthenaError getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AthenaError>(create);
  static AthenaError? _defaultInstance;

  @$pb.TagNumber(83386186)
  $core.bool get retryable => $_getBF(0);
  @$pb.TagNumber(83386186)
  set retryable($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(83386186)
  $core.bool hasRetryable() => $_has(0);
  @$pb.TagNumber(83386186)
  void clearRetryable() => $_clearField(83386186);

  @$pb.TagNumber(315958414)
  $core.int get errorcategory => $_getIZ(1);
  @$pb.TagNumber(315958414)
  set errorcategory($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(315958414)
  $core.bool hasErrorcategory() => $_has(1);
  @$pb.TagNumber(315958414)
  void clearErrorcategory() => $_clearField(315958414);

  @$pb.TagNumber(398848954)
  $core.int get errortype => $_getIZ(2);
  @$pb.TagNumber(398848954)
  set errortype($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(398848954)
  $core.bool hasErrortype() => $_has(2);
  @$pb.TagNumber(398848954)
  void clearErrortype() => $_clearField(398848954);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(3);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(3, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(3);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class BatchGetNamedQueryInput extends $pb.GeneratedMessage {
  factory BatchGetNamedQueryInput({
    $core.Iterable<$core.String>? namedqueryids,
  }) {
    final result = create();
    if (namedqueryids != null) result.namedqueryids.addAll(namedqueryids);
    return result;
  }

  BatchGetNamedQueryInput._();

  factory BatchGetNamedQueryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetNamedQueryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetNamedQueryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPS(6092797, _omitFieldNames ? '' : 'namedqueryids')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetNamedQueryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetNamedQueryInput copyWith(
          void Function(BatchGetNamedQueryInput) updates) =>
      super.copyWith((message) => updates(message as BatchGetNamedQueryInput))
          as BatchGetNamedQueryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetNamedQueryInput create() => BatchGetNamedQueryInput._();
  @$core.override
  BatchGetNamedQueryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetNamedQueryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetNamedQueryInput>(create);
  static BatchGetNamedQueryInput? _defaultInstance;

  @$pb.TagNumber(6092797)
  $pb.PbList<$core.String> get namedqueryids => $_getList(0);
}

class BatchGetNamedQueryOutput extends $pb.GeneratedMessage {
  factory BatchGetNamedQueryOutput({
    $core.Iterable<NamedQuery>? namedqueries,
    $core.Iterable<UnprocessedNamedQueryId>? unprocessednamedqueryids,
  }) {
    final result = create();
    if (namedqueries != null) result.namedqueries.addAll(namedqueries);
    if (unprocessednamedqueryids != null)
      result.unprocessednamedqueryids.addAll(unprocessednamedqueryids);
    return result;
  }

  BatchGetNamedQueryOutput._();

  factory BatchGetNamedQueryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetNamedQueryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetNamedQueryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<NamedQuery>(384985627, _omitFieldNames ? '' : 'namedqueries',
        subBuilder: NamedQuery.create)
    ..pPM<UnprocessedNamedQueryId>(
        529659590, _omitFieldNames ? '' : 'unprocessednamedqueryids',
        subBuilder: UnprocessedNamedQueryId.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetNamedQueryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetNamedQueryOutput copyWith(
          void Function(BatchGetNamedQueryOutput) updates) =>
      super.copyWith((message) => updates(message as BatchGetNamedQueryOutput))
          as BatchGetNamedQueryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetNamedQueryOutput create() => BatchGetNamedQueryOutput._();
  @$core.override
  BatchGetNamedQueryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetNamedQueryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetNamedQueryOutput>(create);
  static BatchGetNamedQueryOutput? _defaultInstance;

  @$pb.TagNumber(384985627)
  $pb.PbList<NamedQuery> get namedqueries => $_getList(0);

  @$pb.TagNumber(529659590)
  $pb.PbList<UnprocessedNamedQueryId> get unprocessednamedqueryids =>
      $_getList(1);
}

class BatchGetPreparedStatementInput extends $pb.GeneratedMessage {
  factory BatchGetPreparedStatementInput({
    $core.Iterable<$core.String>? preparedstatementnames,
    $core.String? workgroup,
  }) {
    final result = create();
    if (preparedstatementnames != null)
      result.preparedstatementnames.addAll(preparedstatementnames);
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  BatchGetPreparedStatementInput._();

  factory BatchGetPreparedStatementInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetPreparedStatementInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetPreparedStatementInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPS(352693824, _omitFieldNames ? '' : 'preparedstatementnames')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetPreparedStatementInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetPreparedStatementInput copyWith(
          void Function(BatchGetPreparedStatementInput) updates) =>
      super.copyWith(
              (message) => updates(message as BatchGetPreparedStatementInput))
          as BatchGetPreparedStatementInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetPreparedStatementInput create() =>
      BatchGetPreparedStatementInput._();
  @$core.override
  BatchGetPreparedStatementInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetPreparedStatementInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetPreparedStatementInput>(create);
  static BatchGetPreparedStatementInput? _defaultInstance;

  @$pb.TagNumber(352693824)
  $pb.PbList<$core.String> get preparedstatementnames => $_getList(0);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(1);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(1);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class BatchGetPreparedStatementOutput extends $pb.GeneratedMessage {
  factory BatchGetPreparedStatementOutput({
    $core.Iterable<UnprocessedPreparedStatementName>?
        unprocessedpreparedstatementnames,
    $core.Iterable<PreparedStatement>? preparedstatements,
  }) {
    final result = create();
    if (unprocessedpreparedstatementnames != null)
      result.unprocessedpreparedstatementnames
          .addAll(unprocessedpreparedstatementnames);
    if (preparedstatements != null)
      result.preparedstatements.addAll(preparedstatements);
    return result;
  }

  BatchGetPreparedStatementOutput._();

  factory BatchGetPreparedStatementOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetPreparedStatementOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetPreparedStatementOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<UnprocessedPreparedStatementName>(
        526426761, _omitFieldNames ? '' : 'unprocessedpreparedstatementnames',
        subBuilder: UnprocessedPreparedStatementName.create)
    ..pPM<PreparedStatement>(
        526923667, _omitFieldNames ? '' : 'preparedstatements',
        subBuilder: PreparedStatement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetPreparedStatementOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetPreparedStatementOutput copyWith(
          void Function(BatchGetPreparedStatementOutput) updates) =>
      super.copyWith(
              (message) => updates(message as BatchGetPreparedStatementOutput))
          as BatchGetPreparedStatementOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetPreparedStatementOutput create() =>
      BatchGetPreparedStatementOutput._();
  @$core.override
  BatchGetPreparedStatementOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetPreparedStatementOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetPreparedStatementOutput>(
          create);
  static BatchGetPreparedStatementOutput? _defaultInstance;

  @$pb.TagNumber(526426761)
  $pb.PbList<UnprocessedPreparedStatementName>
      get unprocessedpreparedstatementnames => $_getList(0);

  @$pb.TagNumber(526923667)
  $pb.PbList<PreparedStatement> get preparedstatements => $_getList(1);
}

class BatchGetQueryExecutionInput extends $pb.GeneratedMessage {
  factory BatchGetQueryExecutionInput({
    $core.Iterable<$core.String>? queryexecutionids,
  }) {
    final result = create();
    if (queryexecutionids != null)
      result.queryexecutionids.addAll(queryexecutionids);
    return result;
  }

  BatchGetQueryExecutionInput._();

  factory BatchGetQueryExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetQueryExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetQueryExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPS(493941192, _omitFieldNames ? '' : 'queryexecutionids')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetQueryExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetQueryExecutionInput copyWith(
          void Function(BatchGetQueryExecutionInput) updates) =>
      super.copyWith(
              (message) => updates(message as BatchGetQueryExecutionInput))
          as BatchGetQueryExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetQueryExecutionInput create() =>
      BatchGetQueryExecutionInput._();
  @$core.override
  BatchGetQueryExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetQueryExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetQueryExecutionInput>(create);
  static BatchGetQueryExecutionInput? _defaultInstance;

  @$pb.TagNumber(493941192)
  $pb.PbList<$core.String> get queryexecutionids => $_getList(0);
}

class BatchGetQueryExecutionOutput extends $pb.GeneratedMessage {
  factory BatchGetQueryExecutionOutput({
    $core.Iterable<UnprocessedQueryExecutionId>? unprocessedqueryexecutionids,
    $core.Iterable<QueryExecution>? queryexecutions,
  }) {
    final result = create();
    if (unprocessedqueryexecutionids != null)
      result.unprocessedqueryexecutionids.addAll(unprocessedqueryexecutionids);
    if (queryexecutions != null) result.queryexecutions.addAll(queryexecutions);
    return result;
  }

  BatchGetQueryExecutionOutput._();

  factory BatchGetQueryExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetQueryExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetQueryExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<UnprocessedQueryExecutionId>(
        85288387, _omitFieldNames ? '' : 'unprocessedqueryexecutionids',
        subBuilder: UnprocessedQueryExecutionId.create)
    ..pPM<QueryExecution>(91405361, _omitFieldNames ? '' : 'queryexecutions',
        subBuilder: QueryExecution.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetQueryExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetQueryExecutionOutput copyWith(
          void Function(BatchGetQueryExecutionOutput) updates) =>
      super.copyWith(
              (message) => updates(message as BatchGetQueryExecutionOutput))
          as BatchGetQueryExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetQueryExecutionOutput create() =>
      BatchGetQueryExecutionOutput._();
  @$core.override
  BatchGetQueryExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetQueryExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetQueryExecutionOutput>(create);
  static BatchGetQueryExecutionOutput? _defaultInstance;

  @$pb.TagNumber(85288387)
  $pb.PbList<UnprocessedQueryExecutionId> get unprocessedqueryexecutionids =>
      $_getList(0);

  @$pb.TagNumber(91405361)
  $pb.PbList<QueryExecution> get queryexecutions => $_getList(1);
}

class CalculationConfiguration extends $pb.GeneratedMessage {
  factory CalculationConfiguration({
    $core.String? codeblock,
  }) {
    final result = create();
    if (codeblock != null) result.codeblock = codeblock;
    return result;
  }

  CalculationConfiguration._();

  factory CalculationConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CalculationConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CalculationConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23945838, _omitFieldNames ? '' : 'codeblock')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationConfiguration copyWith(
          void Function(CalculationConfiguration) updates) =>
      super.copyWith((message) => updates(message as CalculationConfiguration))
          as CalculationConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CalculationConfiguration create() => CalculationConfiguration._();
  @$core.override
  CalculationConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CalculationConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CalculationConfiguration>(create);
  static CalculationConfiguration? _defaultInstance;

  @$pb.TagNumber(23945838)
  $core.String get codeblock => $_getSZ(0);
  @$pb.TagNumber(23945838)
  set codeblock($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23945838)
  $core.bool hasCodeblock() => $_has(0);
  @$pb.TagNumber(23945838)
  void clearCodeblock() => $_clearField(23945838);
}

class CalculationResult extends $pb.GeneratedMessage {
  factory CalculationResult({
    $core.String? stdouts3uri,
    $core.String? results3uri,
    $core.String? stderrors3uri,
    $core.String? resulttype,
  }) {
    final result = create();
    if (stdouts3uri != null) result.stdouts3uri = stdouts3uri;
    if (results3uri != null) result.results3uri = results3uri;
    if (stderrors3uri != null) result.stderrors3uri = stderrors3uri;
    if (resulttype != null) result.resulttype = resulttype;
    return result;
  }

  CalculationResult._();

  factory CalculationResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CalculationResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CalculationResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(100477551, _omitFieldNames ? '' : 'stdouts3uri')
    ..aOS(120920277, _omitFieldNames ? '' : 'results3uri')
    ..aOS(227945541, _omitFieldNames ? '' : 'stderrors3uri')
    ..aOS(261078637, _omitFieldNames ? '' : 'resulttype')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationResult copyWith(void Function(CalculationResult) updates) =>
      super.copyWith((message) => updates(message as CalculationResult))
          as CalculationResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CalculationResult create() => CalculationResult._();
  @$core.override
  CalculationResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CalculationResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CalculationResult>(create);
  static CalculationResult? _defaultInstance;

  @$pb.TagNumber(100477551)
  $core.String get stdouts3uri => $_getSZ(0);
  @$pb.TagNumber(100477551)
  set stdouts3uri($core.String value) => $_setString(0, value);
  @$pb.TagNumber(100477551)
  $core.bool hasStdouts3uri() => $_has(0);
  @$pb.TagNumber(100477551)
  void clearStdouts3uri() => $_clearField(100477551);

  @$pb.TagNumber(120920277)
  $core.String get results3uri => $_getSZ(1);
  @$pb.TagNumber(120920277)
  set results3uri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(120920277)
  $core.bool hasResults3uri() => $_has(1);
  @$pb.TagNumber(120920277)
  void clearResults3uri() => $_clearField(120920277);

  @$pb.TagNumber(227945541)
  $core.String get stderrors3uri => $_getSZ(2);
  @$pb.TagNumber(227945541)
  set stderrors3uri($core.String value) => $_setString(2, value);
  @$pb.TagNumber(227945541)
  $core.bool hasStderrors3uri() => $_has(2);
  @$pb.TagNumber(227945541)
  void clearStderrors3uri() => $_clearField(227945541);

  @$pb.TagNumber(261078637)
  $core.String get resulttype => $_getSZ(3);
  @$pb.TagNumber(261078637)
  set resulttype($core.String value) => $_setString(3, value);
  @$pb.TagNumber(261078637)
  $core.bool hasResulttype() => $_has(3);
  @$pb.TagNumber(261078637)
  void clearResulttype() => $_clearField(261078637);
}

class CalculationStatistics extends $pb.GeneratedMessage {
  factory CalculationStatistics({
    $fixnum.Int64? dpuexecutioninmillis,
    $core.String? progress,
  }) {
    final result = create();
    if (dpuexecutioninmillis != null)
      result.dpuexecutioninmillis = dpuexecutioninmillis;
    if (progress != null) result.progress = progress;
    return result;
  }

  CalculationStatistics._();

  factory CalculationStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CalculationStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CalculationStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(174857936, _omitFieldNames ? '' : 'dpuexecutioninmillis')
    ..aOS(439787879, _omitFieldNames ? '' : 'progress')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationStatistics copyWith(
          void Function(CalculationStatistics) updates) =>
      super.copyWith((message) => updates(message as CalculationStatistics))
          as CalculationStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CalculationStatistics create() => CalculationStatistics._();
  @$core.override
  CalculationStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CalculationStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CalculationStatistics>(create);
  static CalculationStatistics? _defaultInstance;

  @$pb.TagNumber(174857936)
  $fixnum.Int64 get dpuexecutioninmillis => $_getI64(0);
  @$pb.TagNumber(174857936)
  set dpuexecutioninmillis($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(174857936)
  $core.bool hasDpuexecutioninmillis() => $_has(0);
  @$pb.TagNumber(174857936)
  void clearDpuexecutioninmillis() => $_clearField(174857936);

  @$pb.TagNumber(439787879)
  $core.String get progress => $_getSZ(1);
  @$pb.TagNumber(439787879)
  set progress($core.String value) => $_setString(1, value);
  @$pb.TagNumber(439787879)
  $core.bool hasProgress() => $_has(1);
  @$pb.TagNumber(439787879)
  void clearProgress() => $_clearField(439787879);
}

class CalculationStatus extends $pb.GeneratedMessage {
  factory CalculationStatus({
    $core.String? completiondatetime,
    $core.String? statechangereason,
    $core.String? submissiondatetime,
    CalculationExecutionState? state,
  }) {
    final result = create();
    if (completiondatetime != null)
      result.completiondatetime = completiondatetime;
    if (statechangereason != null) result.statechangereason = statechangereason;
    if (submissiondatetime != null)
      result.submissiondatetime = submissiondatetime;
    if (state != null) result.state = state;
    return result;
  }

  CalculationStatus._();

  factory CalculationStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CalculationStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CalculationStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(175822779, _omitFieldNames ? '' : 'completiondatetime')
    ..aOS(228940439, _omitFieldNames ? '' : 'statechangereason')
    ..aOS(449650437, _omitFieldNames ? '' : 'submissiondatetime')
    ..aE<CalculationExecutionState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: CalculationExecutionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationStatus copyWith(void Function(CalculationStatus) updates) =>
      super.copyWith((message) => updates(message as CalculationStatus))
          as CalculationStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CalculationStatus create() => CalculationStatus._();
  @$core.override
  CalculationStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CalculationStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CalculationStatus>(create);
  static CalculationStatus? _defaultInstance;

  @$pb.TagNumber(175822779)
  $core.String get completiondatetime => $_getSZ(0);
  @$pb.TagNumber(175822779)
  set completiondatetime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(175822779)
  $core.bool hasCompletiondatetime() => $_has(0);
  @$pb.TagNumber(175822779)
  void clearCompletiondatetime() => $_clearField(175822779);

  @$pb.TagNumber(228940439)
  $core.String get statechangereason => $_getSZ(1);
  @$pb.TagNumber(228940439)
  set statechangereason($core.String value) => $_setString(1, value);
  @$pb.TagNumber(228940439)
  $core.bool hasStatechangereason() => $_has(1);
  @$pb.TagNumber(228940439)
  void clearStatechangereason() => $_clearField(228940439);

  @$pb.TagNumber(449650437)
  $core.String get submissiondatetime => $_getSZ(2);
  @$pb.TagNumber(449650437)
  set submissiondatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(449650437)
  $core.bool hasSubmissiondatetime() => $_has(2);
  @$pb.TagNumber(449650437)
  void clearSubmissiondatetime() => $_clearField(449650437);

  @$pb.TagNumber(502047895)
  CalculationExecutionState get state => $_getN(3);
  @$pb.TagNumber(502047895)
  set state(CalculationExecutionState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(3);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class CalculationSummary extends $pb.GeneratedMessage {
  factory CalculationSummary({
    CalculationStatus? status,
    $core.String? calculationexecutionid,
    $core.String? description,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (calculationexecutionid != null)
      result.calculationexecutionid = calculationexecutionid;
    if (description != null) result.description = description;
    return result;
  }

  CalculationSummary._();

  factory CalculationSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CalculationSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CalculationSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<CalculationStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: CalculationStatus.create)
    ..aOS(80028050, _omitFieldNames ? '' : 'calculationexecutionid')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculationSummary copyWith(void Function(CalculationSummary) updates) =>
      super.copyWith((message) => updates(message as CalculationSummary))
          as CalculationSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CalculationSummary create() => CalculationSummary._();
  @$core.override
  CalculationSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CalculationSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CalculationSummary>(create);
  static CalculationSummary? _defaultInstance;

  @$pb.TagNumber(6222352)
  CalculationStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(CalculationStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  CalculationStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(80028050)
  $core.String get calculationexecutionid => $_getSZ(1);
  @$pb.TagNumber(80028050)
  set calculationexecutionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(80028050)
  $core.bool hasCalculationexecutionid() => $_has(1);
  @$pb.TagNumber(80028050)
  void clearCalculationexecutionid() => $_clearField(80028050);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);
}

class CancelCapacityReservationInput extends $pb.GeneratedMessage {
  factory CancelCapacityReservationInput({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  CancelCapacityReservationInput._();

  factory CancelCapacityReservationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelCapacityReservationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelCapacityReservationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelCapacityReservationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelCapacityReservationInput copyWith(
          void Function(CancelCapacityReservationInput) updates) =>
      super.copyWith(
              (message) => updates(message as CancelCapacityReservationInput))
          as CancelCapacityReservationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelCapacityReservationInput create() =>
      CancelCapacityReservationInput._();
  @$core.override
  CancelCapacityReservationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelCapacityReservationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelCapacityReservationInput>(create);
  static CancelCapacityReservationInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CancelCapacityReservationOutput extends $pb.GeneratedMessage {
  factory CancelCapacityReservationOutput() => create();

  CancelCapacityReservationOutput._();

  factory CancelCapacityReservationOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelCapacityReservationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelCapacityReservationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelCapacityReservationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelCapacityReservationOutput copyWith(
          void Function(CancelCapacityReservationOutput) updates) =>
      super.copyWith(
              (message) => updates(message as CancelCapacityReservationOutput))
          as CancelCapacityReservationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelCapacityReservationOutput create() =>
      CancelCapacityReservationOutput._();
  @$core.override
  CancelCapacityReservationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelCapacityReservationOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelCapacityReservationOutput>(
          create);
  static CancelCapacityReservationOutput? _defaultInstance;
}

class CapacityAllocation extends $pb.GeneratedMessage {
  factory CapacityAllocation({
    CapacityAllocationStatus? status,
    $core.String? statusmessage,
    $core.String? requestcompletiontime,
    $core.String? requesttime,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (statusmessage != null) result.statusmessage = statusmessage;
    if (requestcompletiontime != null)
      result.requestcompletiontime = requestcompletiontime;
    if (requesttime != null) result.requesttime = requesttime;
    return result;
  }

  CapacityAllocation._();

  factory CapacityAllocation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CapacityAllocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CapacityAllocation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<CapacityAllocationStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: CapacityAllocationStatus.values)
    ..aOS(72590095, _omitFieldNames ? '' : 'statusmessage')
    ..aOS(324037826, _omitFieldNames ? '' : 'requestcompletiontime')
    ..aOS(507812148, _omitFieldNames ? '' : 'requesttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityAllocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityAllocation copyWith(void Function(CapacityAllocation) updates) =>
      super.copyWith((message) => updates(message as CapacityAllocation))
          as CapacityAllocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CapacityAllocation create() => CapacityAllocation._();
  @$core.override
  CapacityAllocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CapacityAllocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CapacityAllocation>(create);
  static CapacityAllocation? _defaultInstance;

  @$pb.TagNumber(6222352)
  CapacityAllocationStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(CapacityAllocationStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(72590095)
  $core.String get statusmessage => $_getSZ(1);
  @$pb.TagNumber(72590095)
  set statusmessage($core.String value) => $_setString(1, value);
  @$pb.TagNumber(72590095)
  $core.bool hasStatusmessage() => $_has(1);
  @$pb.TagNumber(72590095)
  void clearStatusmessage() => $_clearField(72590095);

  @$pb.TagNumber(324037826)
  $core.String get requestcompletiontime => $_getSZ(2);
  @$pb.TagNumber(324037826)
  set requestcompletiontime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(324037826)
  $core.bool hasRequestcompletiontime() => $_has(2);
  @$pb.TagNumber(324037826)
  void clearRequestcompletiontime() => $_clearField(324037826);

  @$pb.TagNumber(507812148)
  $core.String get requesttime => $_getSZ(3);
  @$pb.TagNumber(507812148)
  set requesttime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(507812148)
  $core.bool hasRequesttime() => $_has(3);
  @$pb.TagNumber(507812148)
  void clearRequesttime() => $_clearField(507812148);
}

class CapacityAssignment extends $pb.GeneratedMessage {
  factory CapacityAssignment({
    $core.Iterable<$core.String>? workgroupnames,
  }) {
    final result = create();
    if (workgroupnames != null) result.workgroupnames.addAll(workgroupnames);
    return result;
  }

  CapacityAssignment._();

  factory CapacityAssignment.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CapacityAssignment.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CapacityAssignment',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPS(486049026, _omitFieldNames ? '' : 'workgroupnames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityAssignment clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityAssignment copyWith(void Function(CapacityAssignment) updates) =>
      super.copyWith((message) => updates(message as CapacityAssignment))
          as CapacityAssignment;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CapacityAssignment create() => CapacityAssignment._();
  @$core.override
  CapacityAssignment createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CapacityAssignment getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CapacityAssignment>(create);
  static CapacityAssignment? _defaultInstance;

  @$pb.TagNumber(486049026)
  $pb.PbList<$core.String> get workgroupnames => $_getList(0);
}

class CapacityAssignmentConfiguration extends $pb.GeneratedMessage {
  factory CapacityAssignmentConfiguration({
    $core.String? capacityreservationname,
    $core.Iterable<CapacityAssignment>? capacityassignments,
  }) {
    final result = create();
    if (capacityreservationname != null)
      result.capacityreservationname = capacityreservationname;
    if (capacityassignments != null)
      result.capacityassignments.addAll(capacityassignments);
    return result;
  }

  CapacityAssignmentConfiguration._();

  factory CapacityAssignmentConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CapacityAssignmentConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CapacityAssignmentConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(327567687, _omitFieldNames ? '' : 'capacityreservationname')
    ..pPM<CapacityAssignment>(
        345772294, _omitFieldNames ? '' : 'capacityassignments',
        subBuilder: CapacityAssignment.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityAssignmentConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityAssignmentConfiguration copyWith(
          void Function(CapacityAssignmentConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as CapacityAssignmentConfiguration))
          as CapacityAssignmentConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CapacityAssignmentConfiguration create() =>
      CapacityAssignmentConfiguration._();
  @$core.override
  CapacityAssignmentConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CapacityAssignmentConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CapacityAssignmentConfiguration>(
          create);
  static CapacityAssignmentConfiguration? _defaultInstance;

  @$pb.TagNumber(327567687)
  $core.String get capacityreservationname => $_getSZ(0);
  @$pb.TagNumber(327567687)
  set capacityreservationname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(327567687)
  $core.bool hasCapacityreservationname() => $_has(0);
  @$pb.TagNumber(327567687)
  void clearCapacityreservationname() => $_clearField(327567687);

  @$pb.TagNumber(345772294)
  $pb.PbList<CapacityAssignment> get capacityassignments => $_getList(1);
}

class CapacityReservation extends $pb.GeneratedMessage {
  factory CapacityReservation({
    CapacityReservationStatus? status,
    $core.String? creationtime,
    $core.int? allocateddpus,
    $core.String? name,
    CapacityAllocation? lastallocation,
    $core.int? targetdpus,
    $core.String? lastsuccessfulallocationtime,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (creationtime != null) result.creationtime = creationtime;
    if (allocateddpus != null) result.allocateddpus = allocateddpus;
    if (name != null) result.name = name;
    if (lastallocation != null) result.lastallocation = lastallocation;
    if (targetdpus != null) result.targetdpus = targetdpus;
    if (lastsuccessfulallocationtime != null)
      result.lastsuccessfulallocationtime = lastsuccessfulallocationtime;
    return result;
  }

  CapacityReservation._();

  factory CapacityReservation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CapacityReservation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CapacityReservation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<CapacityReservationStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: CapacityReservationStatus.values)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aI(252958879, _omitFieldNames ? '' : 'allocateddpus')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<CapacityAllocation>(
        274654476, _omitFieldNames ? '' : 'lastallocation',
        subBuilder: CapacityAllocation.create)
    ..aI(367520745, _omitFieldNames ? '' : 'targetdpus')
    ..aOS(383368641, _omitFieldNames ? '' : 'lastsuccessfulallocationtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityReservation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityReservation copyWith(void Function(CapacityReservation) updates) =>
      super.copyWith((message) => updates(message as CapacityReservation))
          as CapacityReservation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CapacityReservation create() => CapacityReservation._();
  @$core.override
  CapacityReservation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CapacityReservation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CapacityReservation>(create);
  static CapacityReservation? _defaultInstance;

  @$pb.TagNumber(6222352)
  CapacityReservationStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(CapacityReservationStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(252958879)
  $core.int get allocateddpus => $_getIZ(2);
  @$pb.TagNumber(252958879)
  set allocateddpus($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(252958879)
  $core.bool hasAllocateddpus() => $_has(2);
  @$pb.TagNumber(252958879)
  void clearAllocateddpus() => $_clearField(252958879);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(274654476)
  CapacityAllocation get lastallocation => $_getN(4);
  @$pb.TagNumber(274654476)
  set lastallocation(CapacityAllocation value) => $_setField(274654476, value);
  @$pb.TagNumber(274654476)
  $core.bool hasLastallocation() => $_has(4);
  @$pb.TagNumber(274654476)
  void clearLastallocation() => $_clearField(274654476);
  @$pb.TagNumber(274654476)
  CapacityAllocation ensureLastallocation() => $_ensure(4);

  @$pb.TagNumber(367520745)
  $core.int get targetdpus => $_getIZ(5);
  @$pb.TagNumber(367520745)
  set targetdpus($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(367520745)
  $core.bool hasTargetdpus() => $_has(5);
  @$pb.TagNumber(367520745)
  void clearTargetdpus() => $_clearField(367520745);

  @$pb.TagNumber(383368641)
  $core.String get lastsuccessfulallocationtime => $_getSZ(6);
  @$pb.TagNumber(383368641)
  set lastsuccessfulallocationtime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(383368641)
  $core.bool hasLastsuccessfulallocationtime() => $_has(6);
  @$pb.TagNumber(383368641)
  void clearLastsuccessfulallocationtime() => $_clearField(383368641);
}

class Classification extends $pb.GeneratedMessage {
  factory Classification({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? properties,
    $core.String? name,
  }) {
    final result = create();
    if (properties != null) result.properties.addEntries(properties);
    if (name != null) result.name = name;
    return result;
  }

  Classification._();

  factory Classification.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Classification.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Classification',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        29886973, _omitFieldNames ? '' : 'properties',
        entryClassName: 'Classification.PropertiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Classification clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Classification copyWith(void Function(Classification) updates) =>
      super.copyWith((message) => updates(message as Classification))
          as Classification;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Classification create() => Classification._();
  @$core.override
  Classification createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Classification getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Classification>(create);
  static Classification? _defaultInstance;

  @$pb.TagNumber(29886973)
  $pb.PbMap<$core.String, $core.String> get properties => $_getMap(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CloudWatchLoggingConfiguration extends $pb.GeneratedMessage {
  factory CloudWatchLoggingConfiguration({
    $core.String? loggroup,
    $core.String? logstreamnameprefix,
    $core.bool? enabled,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? logtypes,
  }) {
    final result = create();
    if (loggroup != null) result.loggroup = loggroup;
    if (logstreamnameprefix != null)
      result.logstreamnameprefix = logstreamnameprefix;
    if (enabled != null) result.enabled = enabled;
    if (logtypes != null) result.logtypes.addEntries(logtypes);
    return result;
  }

  CloudWatchLoggingConfiguration._();

  factory CloudWatchLoggingConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudWatchLoggingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudWatchLoggingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(148580073, _omitFieldNames ? '' : 'loggroup')
    ..aOS(437213671, _omitFieldNames ? '' : 'logstreamnameprefix')
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..m<$core.String, $core.String>(
        491693055, _omitFieldNames ? '' : 'logtypes',
        entryClassName: 'CloudWatchLoggingConfiguration.LogtypesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchLoggingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchLoggingConfiguration copyWith(
          void Function(CloudWatchLoggingConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as CloudWatchLoggingConfiguration))
          as CloudWatchLoggingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudWatchLoggingConfiguration create() =>
      CloudWatchLoggingConfiguration._();
  @$core.override
  CloudWatchLoggingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudWatchLoggingConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudWatchLoggingConfiguration>(create);
  static CloudWatchLoggingConfiguration? _defaultInstance;

  @$pb.TagNumber(148580073)
  $core.String get loggroup => $_getSZ(0);
  @$pb.TagNumber(148580073)
  set loggroup($core.String value) => $_setString(0, value);
  @$pb.TagNumber(148580073)
  $core.bool hasLoggroup() => $_has(0);
  @$pb.TagNumber(148580073)
  void clearLoggroup() => $_clearField(148580073);

  @$pb.TagNumber(437213671)
  $core.String get logstreamnameprefix => $_getSZ(1);
  @$pb.TagNumber(437213671)
  set logstreamnameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(437213671)
  $core.bool hasLogstreamnameprefix() => $_has(1);
  @$pb.TagNumber(437213671)
  void clearLogstreamnameprefix() => $_clearField(437213671);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(2);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(2);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);

  @$pb.TagNumber(491693055)
  $pb.PbMap<$core.String, $core.String> get logtypes => $_getMap(3);
}

class Column extends $pb.GeneratedMessage {
  factory Column({
    $core.String? name,
    $core.String? type,
    $core.String? comment,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (comment != null) result.comment = comment;
    return result;
  }

  Column._();

  factory Column.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Column.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Column',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Column clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Column copyWith(void Function(Column) updates) =>
      super.copyWith((message) => updates(message as Column)) as Column;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Column create() => Column._();
  @$core.override
  Column createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Column getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Column>(create);
  static Column? _defaultInstance;

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

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(2);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(2, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(2);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);
}

class ColumnInfo extends $pb.GeneratedMessage {
  factory ColumnInfo({
    $core.int? precision,
    $core.int? scale,
    $core.bool? casesensitive,
    $core.String? name,
    $core.String? tablename,
    $core.String? type,
    ColumnNullable? nullable,
    $core.String? schemaname,
    $core.String? label,
    $core.String? catalogname,
  }) {
    final result = create();
    if (precision != null) result.precision = precision;
    if (scale != null) result.scale = scale;
    if (casesensitive != null) result.casesensitive = casesensitive;
    if (name != null) result.name = name;
    if (tablename != null) result.tablename = tablename;
    if (type != null) result.type = type;
    if (nullable != null) result.nullable = nullable;
    if (schemaname != null) result.schemaname = schemaname;
    if (label != null) result.label = label;
    if (catalogname != null) result.catalogname = catalogname;
    return result;
  }

  ColumnInfo._();

  factory ColumnInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ColumnInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ColumnInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aI(110022584, _omitFieldNames ? '' : 'precision')
    ..aI(139628050, _omitFieldNames ? '' : 'scale')
    ..aOB(258546956, _omitFieldNames ? '' : 'casesensitive')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..aE<ColumnNullable>(373261405, _omitFieldNames ? '' : 'nullable',
        enumValues: ColumnNullable.values)
    ..aOS(443785942, _omitFieldNames ? '' : 'schemaname')
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..aOS(518825212, _omitFieldNames ? '' : 'catalogname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ColumnInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ColumnInfo copyWith(void Function(ColumnInfo) updates) =>
      super.copyWith((message) => updates(message as ColumnInfo)) as ColumnInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ColumnInfo create() => ColumnInfo._();
  @$core.override
  ColumnInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ColumnInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ColumnInfo>(create);
  static ColumnInfo? _defaultInstance;

  @$pb.TagNumber(110022584)
  $core.int get precision => $_getIZ(0);
  @$pb.TagNumber(110022584)
  set precision($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(110022584)
  $core.bool hasPrecision() => $_has(0);
  @$pb.TagNumber(110022584)
  void clearPrecision() => $_clearField(110022584);

  @$pb.TagNumber(139628050)
  $core.int get scale => $_getIZ(1);
  @$pb.TagNumber(139628050)
  set scale($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(139628050)
  $core.bool hasScale() => $_has(1);
  @$pb.TagNumber(139628050)
  void clearScale() => $_clearField(139628050);

  @$pb.TagNumber(258546956)
  $core.bool get casesensitive => $_getBF(2);
  @$pb.TagNumber(258546956)
  set casesensitive($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(258546956)
  $core.bool hasCasesensitive() => $_has(2);
  @$pb.TagNumber(258546956)
  void clearCasesensitive() => $_clearField(258546956);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(4);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(4);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(5);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(5, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(5);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(373261405)
  ColumnNullable get nullable => $_getN(6);
  @$pb.TagNumber(373261405)
  set nullable(ColumnNullable value) => $_setField(373261405, value);
  @$pb.TagNumber(373261405)
  $core.bool hasNullable() => $_has(6);
  @$pb.TagNumber(373261405)
  void clearNullable() => $_clearField(373261405);

  @$pb.TagNumber(443785942)
  $core.String get schemaname => $_getSZ(7);
  @$pb.TagNumber(443785942)
  set schemaname($core.String value) => $_setString(7, value);
  @$pb.TagNumber(443785942)
  $core.bool hasSchemaname() => $_has(7);
  @$pb.TagNumber(443785942)
  void clearSchemaname() => $_clearField(443785942);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(8);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(8, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(8);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);

  @$pb.TagNumber(518825212)
  $core.String get catalogname => $_getSZ(9);
  @$pb.TagNumber(518825212)
  set catalogname($core.String value) => $_setString(9, value);
  @$pb.TagNumber(518825212)
  $core.bool hasCatalogname() => $_has(9);
  @$pb.TagNumber(518825212)
  void clearCatalogname() => $_clearField(518825212);
}

class CreateCapacityReservationInput extends $pb.GeneratedMessage {
  factory CreateCapacityReservationInput({
    $core.String? name,
    $core.int? targetdpus,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (targetdpus != null) result.targetdpus = targetdpus;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateCapacityReservationInput._();

  factory CreateCapacityReservationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateCapacityReservationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateCapacityReservationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aI(367520745, _omitFieldNames ? '' : 'targetdpus')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCapacityReservationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCapacityReservationInput copyWith(
          void Function(CreateCapacityReservationInput) updates) =>
      super.copyWith(
              (message) => updates(message as CreateCapacityReservationInput))
          as CreateCapacityReservationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateCapacityReservationInput create() =>
      CreateCapacityReservationInput._();
  @$core.override
  CreateCapacityReservationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateCapacityReservationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateCapacityReservationInput>(create);
  static CreateCapacityReservationInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(367520745)
  $core.int get targetdpus => $_getIZ(1);
  @$pb.TagNumber(367520745)
  set targetdpus($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(367520745)
  $core.bool hasTargetdpus() => $_has(1);
  @$pb.TagNumber(367520745)
  void clearTargetdpus() => $_clearField(367520745);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(2);
}

class CreateCapacityReservationOutput extends $pb.GeneratedMessage {
  factory CreateCapacityReservationOutput() => create();

  CreateCapacityReservationOutput._();

  factory CreateCapacityReservationOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateCapacityReservationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateCapacityReservationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCapacityReservationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCapacityReservationOutput copyWith(
          void Function(CreateCapacityReservationOutput) updates) =>
      super.copyWith(
              (message) => updates(message as CreateCapacityReservationOutput))
          as CreateCapacityReservationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateCapacityReservationOutput create() =>
      CreateCapacityReservationOutput._();
  @$core.override
  CreateCapacityReservationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateCapacityReservationOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateCapacityReservationOutput>(
          create);
  static CreateCapacityReservationOutput? _defaultInstance;
}

class CreateDataCatalogInput extends $pb.GeneratedMessage {
  factory CreateDataCatalogInput({
    $core.String? description,
    $core.String? name,
    DataCatalogType? type,
    $core.Iterable<Tag>? tags,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (tags != null) result.tags.addAll(tags);
    if (parameters != null) result.parameters.addEntries(parameters);
    return result;
  }

  CreateDataCatalogInput._();

  factory CreateDataCatalogInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDataCatalogInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDataCatalogInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<DataCatalogType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DataCatalogType.values)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..m<$core.String, $core.String>(
        494900218, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'CreateDataCatalogInput.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDataCatalogInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDataCatalogInput copyWith(
          void Function(CreateDataCatalogInput) updates) =>
      super.copyWith((message) => updates(message as CreateDataCatalogInput))
          as CreateDataCatalogInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDataCatalogInput create() => CreateDataCatalogInput._();
  @$core.override
  CreateDataCatalogInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDataCatalogInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDataCatalogInput>(create);
  static CreateDataCatalogInput? _defaultInstance;

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

  @$pb.TagNumber(290836590)
  DataCatalogType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(DataCatalogType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(3);

  @$pb.TagNumber(494900218)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(4);
}

class CreateDataCatalogOutput extends $pb.GeneratedMessage {
  factory CreateDataCatalogOutput({
    DataCatalog? datacatalog,
  }) {
    final result = create();
    if (datacatalog != null) result.datacatalog = datacatalog;
    return result;
  }

  CreateDataCatalogOutput._();

  factory CreateDataCatalogOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDataCatalogOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDataCatalogOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<DataCatalog>(209694649, _omitFieldNames ? '' : 'datacatalog',
        subBuilder: DataCatalog.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDataCatalogOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDataCatalogOutput copyWith(
          void Function(CreateDataCatalogOutput) updates) =>
      super.copyWith((message) => updates(message as CreateDataCatalogOutput))
          as CreateDataCatalogOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDataCatalogOutput create() => CreateDataCatalogOutput._();
  @$core.override
  CreateDataCatalogOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDataCatalogOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDataCatalogOutput>(create);
  static CreateDataCatalogOutput? _defaultInstance;

  @$pb.TagNumber(209694649)
  DataCatalog get datacatalog => $_getN(0);
  @$pb.TagNumber(209694649)
  set datacatalog(DataCatalog value) => $_setField(209694649, value);
  @$pb.TagNumber(209694649)
  $core.bool hasDatacatalog() => $_has(0);
  @$pb.TagNumber(209694649)
  void clearDatacatalog() => $_clearField(209694649);
  @$pb.TagNumber(209694649)
  DataCatalog ensureDatacatalog() => $_ensure(0);
}

class CreateNamedQueryInput extends $pb.GeneratedMessage {
  factory CreateNamedQueryInput({
    $core.String? description,
    $core.String? name,
    $core.String? database,
    $core.String? querystring,
    $core.String? clientrequesttoken,
    $core.String? workgroup,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (database != null) result.database = database;
    if (querystring != null) result.querystring = querystring;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  CreateNamedQueryInput._();

  factory CreateNamedQueryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateNamedQueryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateNamedQueryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(278147289, _omitFieldNames ? '' : 'database')
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNamedQueryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNamedQueryInput copyWith(
          void Function(CreateNamedQueryInput) updates) =>
      super.copyWith((message) => updates(message as CreateNamedQueryInput))
          as CreateNamedQueryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateNamedQueryInput create() => CreateNamedQueryInput._();
  @$core.override
  CreateNamedQueryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateNamedQueryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateNamedQueryInput>(create);
  static CreateNamedQueryInput? _defaultInstance;

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

  @$pb.TagNumber(278147289)
  $core.String get database => $_getSZ(2);
  @$pb.TagNumber(278147289)
  set database($core.String value) => $_setString(2, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(2);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(3);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(3, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(3);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(4);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(4);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(5);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(5, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(5);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class CreateNamedQueryOutput extends $pb.GeneratedMessage {
  factory CreateNamedQueryOutput({
    $core.String? namedqueryid,
  }) {
    final result = create();
    if (namedqueryid != null) result.namedqueryid = namedqueryid;
    return result;
  }

  CreateNamedQueryOutput._();

  factory CreateNamedQueryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateNamedQueryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateNamedQueryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(330896872, _omitFieldNames ? '' : 'namedqueryid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNamedQueryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNamedQueryOutput copyWith(
          void Function(CreateNamedQueryOutput) updates) =>
      super.copyWith((message) => updates(message as CreateNamedQueryOutput))
          as CreateNamedQueryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateNamedQueryOutput create() => CreateNamedQueryOutput._();
  @$core.override
  CreateNamedQueryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateNamedQueryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateNamedQueryOutput>(create);
  static CreateNamedQueryOutput? _defaultInstance;

  @$pb.TagNumber(330896872)
  $core.String get namedqueryid => $_getSZ(0);
  @$pb.TagNumber(330896872)
  set namedqueryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(330896872)
  $core.bool hasNamedqueryid() => $_has(0);
  @$pb.TagNumber(330896872)
  void clearNamedqueryid() => $_clearField(330896872);
}

class CreateNotebookInput extends $pb.GeneratedMessage {
  factory CreateNotebookInput({
    $core.String? name,
    $core.String? clientrequesttoken,
    $core.String? workgroup,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  CreateNotebookInput._();

  factory CreateNotebookInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateNotebookInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateNotebookInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNotebookInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNotebookInput copyWith(void Function(CreateNotebookInput) updates) =>
      super.copyWith((message) => updates(message as CreateNotebookInput))
          as CreateNotebookInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateNotebookInput create() => CreateNotebookInput._();
  @$core.override
  CreateNotebookInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateNotebookInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateNotebookInput>(create);
  static CreateNotebookInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(1);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(1);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(2);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(2);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class CreateNotebookOutput extends $pb.GeneratedMessage {
  factory CreateNotebookOutput({
    $core.String? notebookid,
  }) {
    final result = create();
    if (notebookid != null) result.notebookid = notebookid;
    return result;
  }

  CreateNotebookOutput._();

  factory CreateNotebookOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateNotebookOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateNotebookOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNotebookOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateNotebookOutput copyWith(void Function(CreateNotebookOutput) updates) =>
      super.copyWith((message) => updates(message as CreateNotebookOutput))
          as CreateNotebookOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateNotebookOutput create() => CreateNotebookOutput._();
  @$core.override
  CreateNotebookOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateNotebookOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateNotebookOutput>(create);
  static CreateNotebookOutput? _defaultInstance;

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(0);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(0);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);
}

class CreatePreparedStatementInput extends $pb.GeneratedMessage {
  factory CreatePreparedStatementInput({
    $core.String? statementname,
    $core.String? description,
    $core.String? querystatement,
    $core.String? workgroup,
  }) {
    final result = create();
    if (statementname != null) result.statementname = statementname;
    if (description != null) result.description = description;
    if (querystatement != null) result.querystatement = querystatement;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  CreatePreparedStatementInput._();

  factory CreatePreparedStatementInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePreparedStatementInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePreparedStatementInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(340852217, _omitFieldNames ? '' : 'querystatement')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePreparedStatementInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePreparedStatementInput copyWith(
          void Function(CreatePreparedStatementInput) updates) =>
      super.copyWith(
              (message) => updates(message as CreatePreparedStatementInput))
          as CreatePreparedStatementInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePreparedStatementInput create() =>
      CreatePreparedStatementInput._();
  @$core.override
  CreatePreparedStatementInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePreparedStatementInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePreparedStatementInput>(create);
  static CreatePreparedStatementInput? _defaultInstance;

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(0);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(0);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(340852217)
  $core.String get querystatement => $_getSZ(2);
  @$pb.TagNumber(340852217)
  set querystatement($core.String value) => $_setString(2, value);
  @$pb.TagNumber(340852217)
  $core.bool hasQuerystatement() => $_has(2);
  @$pb.TagNumber(340852217)
  void clearQuerystatement() => $_clearField(340852217);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(3);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(3, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(3);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class CreatePreparedStatementOutput extends $pb.GeneratedMessage {
  factory CreatePreparedStatementOutput() => create();

  CreatePreparedStatementOutput._();

  factory CreatePreparedStatementOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePreparedStatementOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePreparedStatementOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePreparedStatementOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePreparedStatementOutput copyWith(
          void Function(CreatePreparedStatementOutput) updates) =>
      super.copyWith(
              (message) => updates(message as CreatePreparedStatementOutput))
          as CreatePreparedStatementOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePreparedStatementOutput create() =>
      CreatePreparedStatementOutput._();
  @$core.override
  CreatePreparedStatementOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePreparedStatementOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePreparedStatementOutput>(create);
  static CreatePreparedStatementOutput? _defaultInstance;
}

class CreatePresignedNotebookUrlRequest extends $pb.GeneratedMessage {
  factory CreatePresignedNotebookUrlRequest({
    $core.String? sessionid,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    return result;
  }

  CreatePresignedNotebookUrlRequest._();

  factory CreatePresignedNotebookUrlRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePresignedNotebookUrlRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePresignedNotebookUrlRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePresignedNotebookUrlRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePresignedNotebookUrlRequest copyWith(
          void Function(CreatePresignedNotebookUrlRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreatePresignedNotebookUrlRequest))
          as CreatePresignedNotebookUrlRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePresignedNotebookUrlRequest create() =>
      CreatePresignedNotebookUrlRequest._();
  @$core.override
  CreatePresignedNotebookUrlRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePresignedNotebookUrlRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePresignedNotebookUrlRequest>(
          create);
  static CreatePresignedNotebookUrlRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);
}

class CreatePresignedNotebookUrlResponse extends $pb.GeneratedMessage {
  factory CreatePresignedNotebookUrlResponse({
    $core.String? authtoken,
    $fixnum.Int64? authtokenexpirationtime,
    $core.String? notebookurl,
  }) {
    final result = create();
    if (authtoken != null) result.authtoken = authtoken;
    if (authtokenexpirationtime != null)
      result.authtokenexpirationtime = authtokenexpirationtime;
    if (notebookurl != null) result.notebookurl = notebookurl;
    return result;
  }

  CreatePresignedNotebookUrlResponse._();

  factory CreatePresignedNotebookUrlResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePresignedNotebookUrlResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePresignedNotebookUrlResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(349771493, _omitFieldNames ? '' : 'authtoken')
    ..aInt64(406065689, _omitFieldNames ? '' : 'authtokenexpirationtime')
    ..aOS(454352434, _omitFieldNames ? '' : 'notebookurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePresignedNotebookUrlResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePresignedNotebookUrlResponse copyWith(
          void Function(CreatePresignedNotebookUrlResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreatePresignedNotebookUrlResponse))
          as CreatePresignedNotebookUrlResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePresignedNotebookUrlResponse create() =>
      CreatePresignedNotebookUrlResponse._();
  @$core.override
  CreatePresignedNotebookUrlResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePresignedNotebookUrlResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePresignedNotebookUrlResponse>(
          create);
  static CreatePresignedNotebookUrlResponse? _defaultInstance;

  @$pb.TagNumber(349771493)
  $core.String get authtoken => $_getSZ(0);
  @$pb.TagNumber(349771493)
  set authtoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(349771493)
  $core.bool hasAuthtoken() => $_has(0);
  @$pb.TagNumber(349771493)
  void clearAuthtoken() => $_clearField(349771493);

  @$pb.TagNumber(406065689)
  $fixnum.Int64 get authtokenexpirationtime => $_getI64(1);
  @$pb.TagNumber(406065689)
  set authtokenexpirationtime($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(406065689)
  $core.bool hasAuthtokenexpirationtime() => $_has(1);
  @$pb.TagNumber(406065689)
  void clearAuthtokenexpirationtime() => $_clearField(406065689);

  @$pb.TagNumber(454352434)
  $core.String get notebookurl => $_getSZ(2);
  @$pb.TagNumber(454352434)
  set notebookurl($core.String value) => $_setString(2, value);
  @$pb.TagNumber(454352434)
  $core.bool hasNotebookurl() => $_has(2);
  @$pb.TagNumber(454352434)
  void clearNotebookurl() => $_clearField(454352434);
}

class CreateWorkGroupInput extends $pb.GeneratedMessage {
  factory CreateWorkGroupInput({
    $core.String? description,
    $core.String? name,
    $core.Iterable<Tag>? tags,
    WorkGroupConfiguration? configuration,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    if (configuration != null) result.configuration = configuration;
    return result;
  }

  CreateWorkGroupInput._();

  factory CreateWorkGroupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWorkGroupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWorkGroupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOM<WorkGroupConfiguration>(
        442426458, _omitFieldNames ? '' : 'configuration',
        subBuilder: WorkGroupConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWorkGroupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWorkGroupInput copyWith(void Function(CreateWorkGroupInput) updates) =>
      super.copyWith((message) => updates(message as CreateWorkGroupInput))
          as CreateWorkGroupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWorkGroupInput create() => CreateWorkGroupInput._();
  @$core.override
  CreateWorkGroupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWorkGroupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWorkGroupInput>(create);
  static CreateWorkGroupInput? _defaultInstance;

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

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(2);

  @$pb.TagNumber(442426458)
  WorkGroupConfiguration get configuration => $_getN(3);
  @$pb.TagNumber(442426458)
  set configuration(WorkGroupConfiguration value) =>
      $_setField(442426458, value);
  @$pb.TagNumber(442426458)
  $core.bool hasConfiguration() => $_has(3);
  @$pb.TagNumber(442426458)
  void clearConfiguration() => $_clearField(442426458);
  @$pb.TagNumber(442426458)
  WorkGroupConfiguration ensureConfiguration() => $_ensure(3);
}

class CreateWorkGroupOutput extends $pb.GeneratedMessage {
  factory CreateWorkGroupOutput() => create();

  CreateWorkGroupOutput._();

  factory CreateWorkGroupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWorkGroupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWorkGroupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWorkGroupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWorkGroupOutput copyWith(
          void Function(CreateWorkGroupOutput) updates) =>
      super.copyWith((message) => updates(message as CreateWorkGroupOutput))
          as CreateWorkGroupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWorkGroupOutput create() => CreateWorkGroupOutput._();
  @$core.override
  CreateWorkGroupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWorkGroupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWorkGroupOutput>(create);
  static CreateWorkGroupOutput? _defaultInstance;
}

class CustomerContentEncryptionConfiguration extends $pb.GeneratedMessage {
  factory CustomerContentEncryptionConfiguration({
    $core.String? kmskey,
  }) {
    final result = create();
    if (kmskey != null) result.kmskey = kmskey;
    return result;
  }

  CustomerContentEncryptionConfiguration._();

  factory CustomerContentEncryptionConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomerContentEncryptionConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomerContentEncryptionConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(114561194, _omitFieldNames ? '' : 'kmskey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomerContentEncryptionConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomerContentEncryptionConfiguration copyWith(
          void Function(CustomerContentEncryptionConfiguration) updates) =>
      super.copyWith((message) =>
              updates(message as CustomerContentEncryptionConfiguration))
          as CustomerContentEncryptionConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomerContentEncryptionConfiguration create() =>
      CustomerContentEncryptionConfiguration._();
  @$core.override
  CustomerContentEncryptionConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomerContentEncryptionConfiguration getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CustomerContentEncryptionConfiguration>(create);
  static CustomerContentEncryptionConfiguration? _defaultInstance;

  @$pb.TagNumber(114561194)
  $core.String get kmskey => $_getSZ(0);
  @$pb.TagNumber(114561194)
  set kmskey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(114561194)
  $core.bool hasKmskey() => $_has(0);
  @$pb.TagNumber(114561194)
  void clearKmskey() => $_clearField(114561194);
}

class DataCatalog extends $pb.GeneratedMessage {
  factory DataCatalog({
    DataCatalogStatus? status,
    $core.String? description,
    $core.String? name,
    DataCatalogType? type,
    $core.String? error,
    ConnectionType? connectiontype,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (error != null) result.error = error;
    if (connectiontype != null) result.connectiontype = connectiontype;
    if (parameters != null) result.parameters.addEntries(parameters);
    return result;
  }

  DataCatalog._();

  factory DataCatalog.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataCatalog.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataCatalog',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<DataCatalogStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: DataCatalogStatus.values)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<DataCatalogType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DataCatalogType.values)
    ..aOS(328047858, _omitFieldNames ? '' : 'error')
    ..aE<ConnectionType>(489507282, _omitFieldNames ? '' : 'connectiontype',
        enumValues: ConnectionType.values)
    ..m<$core.String, $core.String>(
        494900218, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'DataCatalog.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataCatalog clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataCatalog copyWith(void Function(DataCatalog) updates) =>
      super.copyWith((message) => updates(message as DataCatalog))
          as DataCatalog;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataCatalog create() => DataCatalog._();
  @$core.override
  DataCatalog createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataCatalog getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataCatalog>(create);
  static DataCatalog? _defaultInstance;

  @$pb.TagNumber(6222352)
  DataCatalogStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(DataCatalogStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

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

  @$pb.TagNumber(290836590)
  DataCatalogType get type => $_getN(3);
  @$pb.TagNumber(290836590)
  set type(DataCatalogType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(3);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(328047858)
  $core.String get error => $_getSZ(4);
  @$pb.TagNumber(328047858)
  set error($core.String value) => $_setString(4, value);
  @$pb.TagNumber(328047858)
  $core.bool hasError() => $_has(4);
  @$pb.TagNumber(328047858)
  void clearError() => $_clearField(328047858);

  @$pb.TagNumber(489507282)
  ConnectionType get connectiontype => $_getN(5);
  @$pb.TagNumber(489507282)
  set connectiontype(ConnectionType value) => $_setField(489507282, value);
  @$pb.TagNumber(489507282)
  $core.bool hasConnectiontype() => $_has(5);
  @$pb.TagNumber(489507282)
  void clearConnectiontype() => $_clearField(489507282);

  @$pb.TagNumber(494900218)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(6);
}

class DataCatalogSummary extends $pb.GeneratedMessage {
  factory DataCatalogSummary({
    DataCatalogStatus? status,
    DataCatalogType? type,
    $core.String? error,
    ConnectionType? connectiontype,
    $core.String? catalogname,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (type != null) result.type = type;
    if (error != null) result.error = error;
    if (connectiontype != null) result.connectiontype = connectiontype;
    if (catalogname != null) result.catalogname = catalogname;
    return result;
  }

  DataCatalogSummary._();

  factory DataCatalogSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataCatalogSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataCatalogSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<DataCatalogStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: DataCatalogStatus.values)
    ..aE<DataCatalogType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DataCatalogType.values)
    ..aOS(328047858, _omitFieldNames ? '' : 'error')
    ..aE<ConnectionType>(489507282, _omitFieldNames ? '' : 'connectiontype',
        enumValues: ConnectionType.values)
    ..aOS(518825212, _omitFieldNames ? '' : 'catalogname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataCatalogSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataCatalogSummary copyWith(void Function(DataCatalogSummary) updates) =>
      super.copyWith((message) => updates(message as DataCatalogSummary))
          as DataCatalogSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataCatalogSummary create() => DataCatalogSummary._();
  @$core.override
  DataCatalogSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataCatalogSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataCatalogSummary>(create);
  static DataCatalogSummary? _defaultInstance;

  @$pb.TagNumber(6222352)
  DataCatalogStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(DataCatalogStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(290836590)
  DataCatalogType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(DataCatalogType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(328047858)
  $core.String get error => $_getSZ(2);
  @$pb.TagNumber(328047858)
  set error($core.String value) => $_setString(2, value);
  @$pb.TagNumber(328047858)
  $core.bool hasError() => $_has(2);
  @$pb.TagNumber(328047858)
  void clearError() => $_clearField(328047858);

  @$pb.TagNumber(489507282)
  ConnectionType get connectiontype => $_getN(3);
  @$pb.TagNumber(489507282)
  set connectiontype(ConnectionType value) => $_setField(489507282, value);
  @$pb.TagNumber(489507282)
  $core.bool hasConnectiontype() => $_has(3);
  @$pb.TagNumber(489507282)
  void clearConnectiontype() => $_clearField(489507282);

  @$pb.TagNumber(518825212)
  $core.String get catalogname => $_getSZ(4);
  @$pb.TagNumber(518825212)
  set catalogname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(518825212)
  $core.bool hasCatalogname() => $_has(4);
  @$pb.TagNumber(518825212)
  void clearCatalogname() => $_clearField(518825212);
}

class Database extends $pb.GeneratedMessage {
  factory Database({
    $core.String? description,
    $core.String? name,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (parameters != null) result.parameters.addEntries(parameters);
    return result;
  }

  Database._();

  factory Database.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Database.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Database',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..m<$core.String, $core.String>(
        494900218, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'Database.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Database clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Database copyWith(void Function(Database) updates) =>
      super.copyWith((message) => updates(message as Database)) as Database;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Database create() => Database._();
  @$core.override
  Database createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Database getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Database>(create);
  static Database? _defaultInstance;

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

  @$pb.TagNumber(494900218)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(2);
}

class Datum extends $pb.GeneratedMessage {
  factory Datum({
    $core.String? varcharvalue,
  }) {
    final result = create();
    if (varcharvalue != null) result.varcharvalue = varcharvalue;
    return result;
  }

  Datum._();

  factory Datum.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Datum.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Datum',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(286740796, _omitFieldNames ? '' : 'varcharvalue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Datum clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Datum copyWith(void Function(Datum) updates) =>
      super.copyWith((message) => updates(message as Datum)) as Datum;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Datum create() => Datum._();
  @$core.override
  Datum createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Datum getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Datum>(create);
  static Datum? _defaultInstance;

  @$pb.TagNumber(286740796)
  $core.String get varcharvalue => $_getSZ(0);
  @$pb.TagNumber(286740796)
  set varcharvalue($core.String value) => $_setString(0, value);
  @$pb.TagNumber(286740796)
  $core.bool hasVarcharvalue() => $_has(0);
  @$pb.TagNumber(286740796)
  void clearVarcharvalue() => $_clearField(286740796);
}

class DeleteCapacityReservationInput extends $pb.GeneratedMessage {
  factory DeleteCapacityReservationInput({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeleteCapacityReservationInput._();

  factory DeleteCapacityReservationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteCapacityReservationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteCapacityReservationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCapacityReservationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCapacityReservationInput copyWith(
          void Function(DeleteCapacityReservationInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteCapacityReservationInput))
          as DeleteCapacityReservationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteCapacityReservationInput create() =>
      DeleteCapacityReservationInput._();
  @$core.override
  DeleteCapacityReservationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteCapacityReservationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteCapacityReservationInput>(create);
  static DeleteCapacityReservationInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeleteCapacityReservationOutput extends $pb.GeneratedMessage {
  factory DeleteCapacityReservationOutput() => create();

  DeleteCapacityReservationOutput._();

  factory DeleteCapacityReservationOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteCapacityReservationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteCapacityReservationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCapacityReservationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCapacityReservationOutput copyWith(
          void Function(DeleteCapacityReservationOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteCapacityReservationOutput))
          as DeleteCapacityReservationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteCapacityReservationOutput create() =>
      DeleteCapacityReservationOutput._();
  @$core.override
  DeleteCapacityReservationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteCapacityReservationOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteCapacityReservationOutput>(
          create);
  static DeleteCapacityReservationOutput? _defaultInstance;
}

class DeleteDataCatalogInput extends $pb.GeneratedMessage {
  factory DeleteDataCatalogInput({
    $core.String? name,
    $core.bool? deletecatalogonly,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (deletecatalogonly != null) result.deletecatalogonly = deletecatalogonly;
    return result;
  }

  DeleteDataCatalogInput._();

  factory DeleteDataCatalogInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDataCatalogInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDataCatalogInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOB(310631780, _omitFieldNames ? '' : 'deletecatalogonly')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDataCatalogInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDataCatalogInput copyWith(
          void Function(DeleteDataCatalogInput) updates) =>
      super.copyWith((message) => updates(message as DeleteDataCatalogInput))
          as DeleteDataCatalogInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDataCatalogInput create() => DeleteDataCatalogInput._();
  @$core.override
  DeleteDataCatalogInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDataCatalogInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDataCatalogInput>(create);
  static DeleteDataCatalogInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(310631780)
  $core.bool get deletecatalogonly => $_getBF(1);
  @$pb.TagNumber(310631780)
  set deletecatalogonly($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(310631780)
  $core.bool hasDeletecatalogonly() => $_has(1);
  @$pb.TagNumber(310631780)
  void clearDeletecatalogonly() => $_clearField(310631780);
}

class DeleteDataCatalogOutput extends $pb.GeneratedMessage {
  factory DeleteDataCatalogOutput({
    DataCatalog? datacatalog,
  }) {
    final result = create();
    if (datacatalog != null) result.datacatalog = datacatalog;
    return result;
  }

  DeleteDataCatalogOutput._();

  factory DeleteDataCatalogOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDataCatalogOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDataCatalogOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<DataCatalog>(209694649, _omitFieldNames ? '' : 'datacatalog',
        subBuilder: DataCatalog.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDataCatalogOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDataCatalogOutput copyWith(
          void Function(DeleteDataCatalogOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteDataCatalogOutput))
          as DeleteDataCatalogOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDataCatalogOutput create() => DeleteDataCatalogOutput._();
  @$core.override
  DeleteDataCatalogOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDataCatalogOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDataCatalogOutput>(create);
  static DeleteDataCatalogOutput? _defaultInstance;

  @$pb.TagNumber(209694649)
  DataCatalog get datacatalog => $_getN(0);
  @$pb.TagNumber(209694649)
  set datacatalog(DataCatalog value) => $_setField(209694649, value);
  @$pb.TagNumber(209694649)
  $core.bool hasDatacatalog() => $_has(0);
  @$pb.TagNumber(209694649)
  void clearDatacatalog() => $_clearField(209694649);
  @$pb.TagNumber(209694649)
  DataCatalog ensureDatacatalog() => $_ensure(0);
}

class DeleteNamedQueryInput extends $pb.GeneratedMessage {
  factory DeleteNamedQueryInput({
    $core.String? namedqueryid,
  }) {
    final result = create();
    if (namedqueryid != null) result.namedqueryid = namedqueryid;
    return result;
  }

  DeleteNamedQueryInput._();

  factory DeleteNamedQueryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteNamedQueryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteNamedQueryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(330896872, _omitFieldNames ? '' : 'namedqueryid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNamedQueryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNamedQueryInput copyWith(
          void Function(DeleteNamedQueryInput) updates) =>
      super.copyWith((message) => updates(message as DeleteNamedQueryInput))
          as DeleteNamedQueryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteNamedQueryInput create() => DeleteNamedQueryInput._();
  @$core.override
  DeleteNamedQueryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteNamedQueryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteNamedQueryInput>(create);
  static DeleteNamedQueryInput? _defaultInstance;

  @$pb.TagNumber(330896872)
  $core.String get namedqueryid => $_getSZ(0);
  @$pb.TagNumber(330896872)
  set namedqueryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(330896872)
  $core.bool hasNamedqueryid() => $_has(0);
  @$pb.TagNumber(330896872)
  void clearNamedqueryid() => $_clearField(330896872);
}

class DeleteNamedQueryOutput extends $pb.GeneratedMessage {
  factory DeleteNamedQueryOutput() => create();

  DeleteNamedQueryOutput._();

  factory DeleteNamedQueryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteNamedQueryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteNamedQueryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNamedQueryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNamedQueryOutput copyWith(
          void Function(DeleteNamedQueryOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteNamedQueryOutput))
          as DeleteNamedQueryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteNamedQueryOutput create() => DeleteNamedQueryOutput._();
  @$core.override
  DeleteNamedQueryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteNamedQueryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteNamedQueryOutput>(create);
  static DeleteNamedQueryOutput? _defaultInstance;
}

class DeleteNotebookInput extends $pb.GeneratedMessage {
  factory DeleteNotebookInput({
    $core.String? notebookid,
  }) {
    final result = create();
    if (notebookid != null) result.notebookid = notebookid;
    return result;
  }

  DeleteNotebookInput._();

  factory DeleteNotebookInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteNotebookInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteNotebookInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNotebookInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNotebookInput copyWith(void Function(DeleteNotebookInput) updates) =>
      super.copyWith((message) => updates(message as DeleteNotebookInput))
          as DeleteNotebookInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteNotebookInput create() => DeleteNotebookInput._();
  @$core.override
  DeleteNotebookInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteNotebookInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteNotebookInput>(create);
  static DeleteNotebookInput? _defaultInstance;

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(0);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(0);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);
}

class DeleteNotebookOutput extends $pb.GeneratedMessage {
  factory DeleteNotebookOutput() => create();

  DeleteNotebookOutput._();

  factory DeleteNotebookOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteNotebookOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteNotebookOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNotebookOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteNotebookOutput copyWith(void Function(DeleteNotebookOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteNotebookOutput))
          as DeleteNotebookOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteNotebookOutput create() => DeleteNotebookOutput._();
  @$core.override
  DeleteNotebookOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteNotebookOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteNotebookOutput>(create);
  static DeleteNotebookOutput? _defaultInstance;
}

class DeletePreparedStatementInput extends $pb.GeneratedMessage {
  factory DeletePreparedStatementInput({
    $core.String? statementname,
    $core.String? workgroup,
  }) {
    final result = create();
    if (statementname != null) result.statementname = statementname;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  DeletePreparedStatementInput._();

  factory DeletePreparedStatementInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeletePreparedStatementInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeletePreparedStatementInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePreparedStatementInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePreparedStatementInput copyWith(
          void Function(DeletePreparedStatementInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeletePreparedStatementInput))
          as DeletePreparedStatementInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeletePreparedStatementInput create() =>
      DeletePreparedStatementInput._();
  @$core.override
  DeletePreparedStatementInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeletePreparedStatementInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeletePreparedStatementInput>(create);
  static DeletePreparedStatementInput? _defaultInstance;

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(0);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(0);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(1);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(1);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class DeletePreparedStatementOutput extends $pb.GeneratedMessage {
  factory DeletePreparedStatementOutput() => create();

  DeletePreparedStatementOutput._();

  factory DeletePreparedStatementOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeletePreparedStatementOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeletePreparedStatementOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePreparedStatementOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePreparedStatementOutput copyWith(
          void Function(DeletePreparedStatementOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DeletePreparedStatementOutput))
          as DeletePreparedStatementOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeletePreparedStatementOutput create() =>
      DeletePreparedStatementOutput._();
  @$core.override
  DeletePreparedStatementOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeletePreparedStatementOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeletePreparedStatementOutput>(create);
  static DeletePreparedStatementOutput? _defaultInstance;
}

class DeleteWorkGroupInput extends $pb.GeneratedMessage {
  factory DeleteWorkGroupInput({
    $core.bool? recursivedeleteoption,
    $core.String? workgroup,
  }) {
    final result = create();
    if (recursivedeleteoption != null)
      result.recursivedeleteoption = recursivedeleteoption;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  DeleteWorkGroupInput._();

  factory DeleteWorkGroupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteWorkGroupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteWorkGroupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOB(398167566, _omitFieldNames ? '' : 'recursivedeleteoption')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWorkGroupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWorkGroupInput copyWith(void Function(DeleteWorkGroupInput) updates) =>
      super.copyWith((message) => updates(message as DeleteWorkGroupInput))
          as DeleteWorkGroupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteWorkGroupInput create() => DeleteWorkGroupInput._();
  @$core.override
  DeleteWorkGroupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteWorkGroupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteWorkGroupInput>(create);
  static DeleteWorkGroupInput? _defaultInstance;

  @$pb.TagNumber(398167566)
  $core.bool get recursivedeleteoption => $_getBF(0);
  @$pb.TagNumber(398167566)
  set recursivedeleteoption($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(398167566)
  $core.bool hasRecursivedeleteoption() => $_has(0);
  @$pb.TagNumber(398167566)
  void clearRecursivedeleteoption() => $_clearField(398167566);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(1);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(1);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class DeleteWorkGroupOutput extends $pb.GeneratedMessage {
  factory DeleteWorkGroupOutput() => create();

  DeleteWorkGroupOutput._();

  factory DeleteWorkGroupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteWorkGroupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteWorkGroupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWorkGroupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWorkGroupOutput copyWith(
          void Function(DeleteWorkGroupOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteWorkGroupOutput))
          as DeleteWorkGroupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteWorkGroupOutput create() => DeleteWorkGroupOutput._();
  @$core.override
  DeleteWorkGroupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteWorkGroupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteWorkGroupOutput>(create);
  static DeleteWorkGroupOutput? _defaultInstance;
}

class EncryptionConfiguration extends $pb.GeneratedMessage {
  factory EncryptionConfiguration({
    $core.String? kmskey,
    EncryptionOption? encryptionoption,
  }) {
    final result = create();
    if (kmskey != null) result.kmskey = kmskey;
    if (encryptionoption != null) result.encryptionoption = encryptionoption;
    return result;
  }

  EncryptionConfiguration._();

  factory EncryptionConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EncryptionConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EncryptionConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(114561194, _omitFieldNames ? '' : 'kmskey')
    ..aE<EncryptionOption>(160833062, _omitFieldNames ? '' : 'encryptionoption',
        enumValues: EncryptionOption.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptionConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptionConfiguration copyWith(
          void Function(EncryptionConfiguration) updates) =>
      super.copyWith((message) => updates(message as EncryptionConfiguration))
          as EncryptionConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EncryptionConfiguration create() => EncryptionConfiguration._();
  @$core.override
  EncryptionConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EncryptionConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EncryptionConfiguration>(create);
  static EncryptionConfiguration? _defaultInstance;

  @$pb.TagNumber(114561194)
  $core.String get kmskey => $_getSZ(0);
  @$pb.TagNumber(114561194)
  set kmskey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(114561194)
  $core.bool hasKmskey() => $_has(0);
  @$pb.TagNumber(114561194)
  void clearKmskey() => $_clearField(114561194);

  @$pb.TagNumber(160833062)
  EncryptionOption get encryptionoption => $_getN(1);
  @$pb.TagNumber(160833062)
  set encryptionoption(EncryptionOption value) => $_setField(160833062, value);
  @$pb.TagNumber(160833062)
  $core.bool hasEncryptionoption() => $_has(1);
  @$pb.TagNumber(160833062)
  void clearEncryptionoption() => $_clearField(160833062);
}

class EngineConfiguration extends $pb.GeneratedMessage {
  factory EngineConfiguration({
    $core.int? defaultexecutordpusize,
    $core.Iterable<Classification>? classifications,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? sparkproperties,
    $core.int? coordinatordpusize,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        additionalconfigs,
    $core.int? maxconcurrentdpus,
  }) {
    final result = create();
    if (defaultexecutordpusize != null)
      result.defaultexecutordpusize = defaultexecutordpusize;
    if (classifications != null) result.classifications.addAll(classifications);
    if (sparkproperties != null)
      result.sparkproperties.addEntries(sparkproperties);
    if (coordinatordpusize != null)
      result.coordinatordpusize = coordinatordpusize;
    if (additionalconfigs != null)
      result.additionalconfigs.addEntries(additionalconfigs);
    if (maxconcurrentdpus != null) result.maxconcurrentdpus = maxconcurrentdpus;
    return result;
  }

  EngineConfiguration._();

  factory EngineConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EngineConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EngineConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aI(42210020, _omitFieldNames ? '' : 'defaultexecutordpusize')
    ..pPM<Classification>(129400407, _omitFieldNames ? '' : 'classifications',
        subBuilder: Classification.create)
    ..m<$core.String, $core.String>(
        157244376, _omitFieldNames ? '' : 'sparkproperties',
        entryClassName: 'EngineConfiguration.SparkpropertiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..aI(158247866, _omitFieldNames ? '' : 'coordinatordpusize')
    ..m<$core.String, $core.String>(
        341218716, _omitFieldNames ? '' : 'additionalconfigs',
        entryClassName: 'EngineConfiguration.AdditionalconfigsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..aI(349976345, _omitFieldNames ? '' : 'maxconcurrentdpus')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EngineConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EngineConfiguration copyWith(void Function(EngineConfiguration) updates) =>
      super.copyWith((message) => updates(message as EngineConfiguration))
          as EngineConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EngineConfiguration create() => EngineConfiguration._();
  @$core.override
  EngineConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EngineConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EngineConfiguration>(create);
  static EngineConfiguration? _defaultInstance;

  @$pb.TagNumber(42210020)
  $core.int get defaultexecutordpusize => $_getIZ(0);
  @$pb.TagNumber(42210020)
  set defaultexecutordpusize($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(42210020)
  $core.bool hasDefaultexecutordpusize() => $_has(0);
  @$pb.TagNumber(42210020)
  void clearDefaultexecutordpusize() => $_clearField(42210020);

  @$pb.TagNumber(129400407)
  $pb.PbList<Classification> get classifications => $_getList(1);

  @$pb.TagNumber(157244376)
  $pb.PbMap<$core.String, $core.String> get sparkproperties => $_getMap(2);

  @$pb.TagNumber(158247866)
  $core.int get coordinatordpusize => $_getIZ(3);
  @$pb.TagNumber(158247866)
  set coordinatordpusize($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(158247866)
  $core.bool hasCoordinatordpusize() => $_has(3);
  @$pb.TagNumber(158247866)
  void clearCoordinatordpusize() => $_clearField(158247866);

  @$pb.TagNumber(341218716)
  $pb.PbMap<$core.String, $core.String> get additionalconfigs => $_getMap(4);

  @$pb.TagNumber(349976345)
  $core.int get maxconcurrentdpus => $_getIZ(5);
  @$pb.TagNumber(349976345)
  set maxconcurrentdpus($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(349976345)
  $core.bool hasMaxconcurrentdpus() => $_has(5);
  @$pb.TagNumber(349976345)
  void clearMaxconcurrentdpus() => $_clearField(349976345);
}

class EngineVersion extends $pb.GeneratedMessage {
  factory EngineVersion({
    $core.String? selectedengineversion,
    $core.String? effectiveengineversion,
  }) {
    final result = create();
    if (selectedengineversion != null)
      result.selectedengineversion = selectedengineversion;
    if (effectiveengineversion != null)
      result.effectiveengineversion = effectiveengineversion;
    return result;
  }

  EngineVersion._();

  factory EngineVersion.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EngineVersion.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EngineVersion',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(199609827, _omitFieldNames ? '' : 'selectedengineversion')
    ..aOS(382949365, _omitFieldNames ? '' : 'effectiveengineversion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EngineVersion clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EngineVersion copyWith(void Function(EngineVersion) updates) =>
      super.copyWith((message) => updates(message as EngineVersion))
          as EngineVersion;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EngineVersion create() => EngineVersion._();
  @$core.override
  EngineVersion createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EngineVersion getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EngineVersion>(create);
  static EngineVersion? _defaultInstance;

  @$pb.TagNumber(199609827)
  $core.String get selectedengineversion => $_getSZ(0);
  @$pb.TagNumber(199609827)
  set selectedengineversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(199609827)
  $core.bool hasSelectedengineversion() => $_has(0);
  @$pb.TagNumber(199609827)
  void clearSelectedengineversion() => $_clearField(199609827);

  @$pb.TagNumber(382949365)
  $core.String get effectiveengineversion => $_getSZ(1);
  @$pb.TagNumber(382949365)
  set effectiveengineversion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(382949365)
  $core.bool hasEffectiveengineversion() => $_has(1);
  @$pb.TagNumber(382949365)
  void clearEffectiveengineversion() => $_clearField(382949365);
}

class ExecutorsSummary extends $pb.GeneratedMessage {
  factory ExecutorsSummary({
    $fixnum.Int64? startdatetime,
    ExecutorState? executorstate,
    ExecutorType? executortype,
    $fixnum.Int64? executorsize,
    $fixnum.Int64? terminationdatetime,
    $core.String? executorid,
  }) {
    final result = create();
    if (startdatetime != null) result.startdatetime = startdatetime;
    if (executorstate != null) result.executorstate = executorstate;
    if (executortype != null) result.executortype = executortype;
    if (executorsize != null) result.executorsize = executorsize;
    if (terminationdatetime != null)
      result.terminationdatetime = terminationdatetime;
    if (executorid != null) result.executorid = executorid;
    return result;
  }

  ExecutorsSummary._();

  factory ExecutorsSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutorsSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutorsSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(88518355, _omitFieldNames ? '' : 'startdatetime')
    ..aE<ExecutorState>(187122914, _omitFieldNames ? '' : 'executorstate',
        enumValues: ExecutorState.values)
    ..aE<ExecutorType>(213715221, _omitFieldNames ? '' : 'executortype',
        enumValues: ExecutorType.values)
    ..aInt64(220298710, _omitFieldNames ? '' : 'executorsize')
    ..aInt64(253850405, _omitFieldNames ? '' : 'terminationdatetime')
    ..aOS(307566450, _omitFieldNames ? '' : 'executorid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutorsSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutorsSummary copyWith(void Function(ExecutorsSummary) updates) =>
      super.copyWith((message) => updates(message as ExecutorsSummary))
          as ExecutorsSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutorsSummary create() => ExecutorsSummary._();
  @$core.override
  ExecutorsSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutorsSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutorsSummary>(create);
  static ExecutorsSummary? _defaultInstance;

  @$pb.TagNumber(88518355)
  $fixnum.Int64 get startdatetime => $_getI64(0);
  @$pb.TagNumber(88518355)
  set startdatetime($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(88518355)
  $core.bool hasStartdatetime() => $_has(0);
  @$pb.TagNumber(88518355)
  void clearStartdatetime() => $_clearField(88518355);

  @$pb.TagNumber(187122914)
  ExecutorState get executorstate => $_getN(1);
  @$pb.TagNumber(187122914)
  set executorstate(ExecutorState value) => $_setField(187122914, value);
  @$pb.TagNumber(187122914)
  $core.bool hasExecutorstate() => $_has(1);
  @$pb.TagNumber(187122914)
  void clearExecutorstate() => $_clearField(187122914);

  @$pb.TagNumber(213715221)
  ExecutorType get executortype => $_getN(2);
  @$pb.TagNumber(213715221)
  set executortype(ExecutorType value) => $_setField(213715221, value);
  @$pb.TagNumber(213715221)
  $core.bool hasExecutortype() => $_has(2);
  @$pb.TagNumber(213715221)
  void clearExecutortype() => $_clearField(213715221);

  @$pb.TagNumber(220298710)
  $fixnum.Int64 get executorsize => $_getI64(3);
  @$pb.TagNumber(220298710)
  set executorsize($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(220298710)
  $core.bool hasExecutorsize() => $_has(3);
  @$pb.TagNumber(220298710)
  void clearExecutorsize() => $_clearField(220298710);

  @$pb.TagNumber(253850405)
  $fixnum.Int64 get terminationdatetime => $_getI64(4);
  @$pb.TagNumber(253850405)
  set terminationdatetime($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(253850405)
  $core.bool hasTerminationdatetime() => $_has(4);
  @$pb.TagNumber(253850405)
  void clearTerminationdatetime() => $_clearField(253850405);

  @$pb.TagNumber(307566450)
  $core.String get executorid => $_getSZ(5);
  @$pb.TagNumber(307566450)
  set executorid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(307566450)
  $core.bool hasExecutorid() => $_has(5);
  @$pb.TagNumber(307566450)
  void clearExecutorid() => $_clearField(307566450);
}

class ExportNotebookInput extends $pb.GeneratedMessage {
  factory ExportNotebookInput({
    $core.String? notebookid,
  }) {
    final result = create();
    if (notebookid != null) result.notebookid = notebookid;
    return result;
  }

  ExportNotebookInput._();

  factory ExportNotebookInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportNotebookInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportNotebookInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportNotebookInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportNotebookInput copyWith(void Function(ExportNotebookInput) updates) =>
      super.copyWith((message) => updates(message as ExportNotebookInput))
          as ExportNotebookInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportNotebookInput create() => ExportNotebookInput._();
  @$core.override
  ExportNotebookInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportNotebookInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportNotebookInput>(create);
  static ExportNotebookInput? _defaultInstance;

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(0);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(0);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);
}

class ExportNotebookOutput extends $pb.GeneratedMessage {
  factory ExportNotebookOutput({
    $core.String? payload,
    NotebookMetadata? notebookmetadata,
  }) {
    final result = create();
    if (payload != null) result.payload = payload;
    if (notebookmetadata != null) result.notebookmetadata = notebookmetadata;
    return result;
  }

  ExportNotebookOutput._();

  factory ExportNotebookOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportNotebookOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportNotebookOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(6526790, _omitFieldNames ? '' : 'payload')
    ..aOM<NotebookMetadata>(77536390, _omitFieldNames ? '' : 'notebookmetadata',
        subBuilder: NotebookMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportNotebookOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportNotebookOutput copyWith(void Function(ExportNotebookOutput) updates) =>
      super.copyWith((message) => updates(message as ExportNotebookOutput))
          as ExportNotebookOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportNotebookOutput create() => ExportNotebookOutput._();
  @$core.override
  ExportNotebookOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportNotebookOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportNotebookOutput>(create);
  static ExportNotebookOutput? _defaultInstance;

  @$pb.TagNumber(6526790)
  $core.String get payload => $_getSZ(0);
  @$pb.TagNumber(6526790)
  set payload($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6526790)
  $core.bool hasPayload() => $_has(0);
  @$pb.TagNumber(6526790)
  void clearPayload() => $_clearField(6526790);

  @$pb.TagNumber(77536390)
  NotebookMetadata get notebookmetadata => $_getN(1);
  @$pb.TagNumber(77536390)
  set notebookmetadata(NotebookMetadata value) => $_setField(77536390, value);
  @$pb.TagNumber(77536390)
  $core.bool hasNotebookmetadata() => $_has(1);
  @$pb.TagNumber(77536390)
  void clearNotebookmetadata() => $_clearField(77536390);
  @$pb.TagNumber(77536390)
  NotebookMetadata ensureNotebookmetadata() => $_ensure(1);
}

class FilterDefinition extends $pb.GeneratedMessage {
  factory FilterDefinition({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  FilterDefinition._();

  factory FilterDefinition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilterDefinition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilterDefinition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilterDefinition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilterDefinition copyWith(void Function(FilterDefinition) updates) =>
      super.copyWith((message) => updates(message as FilterDefinition))
          as FilterDefinition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilterDefinition create() => FilterDefinition._();
  @$core.override
  FilterDefinition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilterDefinition getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilterDefinition>(create);
  static FilterDefinition? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class GetCalculationExecutionCodeRequest extends $pb.GeneratedMessage {
  factory GetCalculationExecutionCodeRequest({
    $core.String? calculationexecutionid,
  }) {
    final result = create();
    if (calculationexecutionid != null)
      result.calculationexecutionid = calculationexecutionid;
    return result;
  }

  GetCalculationExecutionCodeRequest._();

  factory GetCalculationExecutionCodeRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCalculationExecutionCodeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCalculationExecutionCodeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(80028050, _omitFieldNames ? '' : 'calculationexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionCodeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionCodeRequest copyWith(
          void Function(GetCalculationExecutionCodeRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetCalculationExecutionCodeRequest))
          as GetCalculationExecutionCodeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionCodeRequest create() =>
      GetCalculationExecutionCodeRequest._();
  @$core.override
  GetCalculationExecutionCodeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionCodeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCalculationExecutionCodeRequest>(
          create);
  static GetCalculationExecutionCodeRequest? _defaultInstance;

  @$pb.TagNumber(80028050)
  $core.String get calculationexecutionid => $_getSZ(0);
  @$pb.TagNumber(80028050)
  set calculationexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(80028050)
  $core.bool hasCalculationexecutionid() => $_has(0);
  @$pb.TagNumber(80028050)
  void clearCalculationexecutionid() => $_clearField(80028050);
}

class GetCalculationExecutionCodeResponse extends $pb.GeneratedMessage {
  factory GetCalculationExecutionCodeResponse({
    $core.String? codeblock,
  }) {
    final result = create();
    if (codeblock != null) result.codeblock = codeblock;
    return result;
  }

  GetCalculationExecutionCodeResponse._();

  factory GetCalculationExecutionCodeResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCalculationExecutionCodeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCalculationExecutionCodeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23945838, _omitFieldNames ? '' : 'codeblock')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionCodeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionCodeResponse copyWith(
          void Function(GetCalculationExecutionCodeResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetCalculationExecutionCodeResponse))
          as GetCalculationExecutionCodeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionCodeResponse create() =>
      GetCalculationExecutionCodeResponse._();
  @$core.override
  GetCalculationExecutionCodeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionCodeResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetCalculationExecutionCodeResponse>(create);
  static GetCalculationExecutionCodeResponse? _defaultInstance;

  @$pb.TagNumber(23945838)
  $core.String get codeblock => $_getSZ(0);
  @$pb.TagNumber(23945838)
  set codeblock($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23945838)
  $core.bool hasCodeblock() => $_has(0);
  @$pb.TagNumber(23945838)
  void clearCodeblock() => $_clearField(23945838);
}

class GetCalculationExecutionRequest extends $pb.GeneratedMessage {
  factory GetCalculationExecutionRequest({
    $core.String? calculationexecutionid,
  }) {
    final result = create();
    if (calculationexecutionid != null)
      result.calculationexecutionid = calculationexecutionid;
    return result;
  }

  GetCalculationExecutionRequest._();

  factory GetCalculationExecutionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCalculationExecutionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCalculationExecutionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(80028050, _omitFieldNames ? '' : 'calculationexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionRequest copyWith(
          void Function(GetCalculationExecutionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetCalculationExecutionRequest))
          as GetCalculationExecutionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionRequest create() =>
      GetCalculationExecutionRequest._();
  @$core.override
  GetCalculationExecutionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCalculationExecutionRequest>(create);
  static GetCalculationExecutionRequest? _defaultInstance;

  @$pb.TagNumber(80028050)
  $core.String get calculationexecutionid => $_getSZ(0);
  @$pb.TagNumber(80028050)
  set calculationexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(80028050)
  $core.bool hasCalculationexecutionid() => $_has(0);
  @$pb.TagNumber(80028050)
  void clearCalculationexecutionid() => $_clearField(80028050);
}

class GetCalculationExecutionResponse extends $pb.GeneratedMessage {
  factory GetCalculationExecutionResponse({
    CalculationStatus? status,
    $core.String? sessionid,
    $core.String? calculationexecutionid,
    $core.String? description,
    CalculationResult? result,
    $core.String? workingdirectory,
    CalculationStatistics? statistics,
  }) {
    final result$ = create();
    if (status != null) result$.status = status;
    if (sessionid != null) result$.sessionid = sessionid;
    if (calculationexecutionid != null)
      result$.calculationexecutionid = calculationexecutionid;
    if (description != null) result$.description = description;
    if (result != null) result$.result = result;
    if (workingdirectory != null) result$.workingdirectory = workingdirectory;
    if (statistics != null) result$.statistics = statistics;
    return result$;
  }

  GetCalculationExecutionResponse._();

  factory GetCalculationExecutionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCalculationExecutionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCalculationExecutionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<CalculationStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: CalculationStatus.create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aOS(80028050, _omitFieldNames ? '' : 'calculationexecutionid')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOM<CalculationResult>(273346629, _omitFieldNames ? '' : 'result',
        subBuilder: CalculationResult.create)
    ..aOS(478970252, _omitFieldNames ? '' : 'workingdirectory')
    ..aOM<CalculationStatistics>(510636075, _omitFieldNames ? '' : 'statistics',
        subBuilder: CalculationStatistics.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionResponse copyWith(
          void Function(GetCalculationExecutionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetCalculationExecutionResponse))
          as GetCalculationExecutionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionResponse create() =>
      GetCalculationExecutionResponse._();
  @$core.override
  GetCalculationExecutionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCalculationExecutionResponse>(
          create);
  static GetCalculationExecutionResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  CalculationStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(CalculationStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  CalculationStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(1);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(1);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(80028050)
  $core.String get calculationexecutionid => $_getSZ(2);
  @$pb.TagNumber(80028050)
  set calculationexecutionid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(80028050)
  $core.bool hasCalculationexecutionid() => $_has(2);
  @$pb.TagNumber(80028050)
  void clearCalculationexecutionid() => $_clearField(80028050);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(273346629)
  CalculationResult get result => $_getN(4);
  @$pb.TagNumber(273346629)
  set result(CalculationResult value) => $_setField(273346629, value);
  @$pb.TagNumber(273346629)
  $core.bool hasResult() => $_has(4);
  @$pb.TagNumber(273346629)
  void clearResult() => $_clearField(273346629);
  @$pb.TagNumber(273346629)
  CalculationResult ensureResult() => $_ensure(4);

  @$pb.TagNumber(478970252)
  $core.String get workingdirectory => $_getSZ(5);
  @$pb.TagNumber(478970252)
  set workingdirectory($core.String value) => $_setString(5, value);
  @$pb.TagNumber(478970252)
  $core.bool hasWorkingdirectory() => $_has(5);
  @$pb.TagNumber(478970252)
  void clearWorkingdirectory() => $_clearField(478970252);

  @$pb.TagNumber(510636075)
  CalculationStatistics get statistics => $_getN(6);
  @$pb.TagNumber(510636075)
  set statistics(CalculationStatistics value) => $_setField(510636075, value);
  @$pb.TagNumber(510636075)
  $core.bool hasStatistics() => $_has(6);
  @$pb.TagNumber(510636075)
  void clearStatistics() => $_clearField(510636075);
  @$pb.TagNumber(510636075)
  CalculationStatistics ensureStatistics() => $_ensure(6);
}

class GetCalculationExecutionStatusRequest extends $pb.GeneratedMessage {
  factory GetCalculationExecutionStatusRequest({
    $core.String? calculationexecutionid,
  }) {
    final result = create();
    if (calculationexecutionid != null)
      result.calculationexecutionid = calculationexecutionid;
    return result;
  }

  GetCalculationExecutionStatusRequest._();

  factory GetCalculationExecutionStatusRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCalculationExecutionStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCalculationExecutionStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(80028050, _omitFieldNames ? '' : 'calculationexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionStatusRequest copyWith(
          void Function(GetCalculationExecutionStatusRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetCalculationExecutionStatusRequest))
          as GetCalculationExecutionStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionStatusRequest create() =>
      GetCalculationExecutionStatusRequest._();
  @$core.override
  GetCalculationExecutionStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionStatusRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetCalculationExecutionStatusRequest>(create);
  static GetCalculationExecutionStatusRequest? _defaultInstance;

  @$pb.TagNumber(80028050)
  $core.String get calculationexecutionid => $_getSZ(0);
  @$pb.TagNumber(80028050)
  set calculationexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(80028050)
  $core.bool hasCalculationexecutionid() => $_has(0);
  @$pb.TagNumber(80028050)
  void clearCalculationexecutionid() => $_clearField(80028050);
}

class GetCalculationExecutionStatusResponse extends $pb.GeneratedMessage {
  factory GetCalculationExecutionStatusResponse({
    CalculationStatus? status,
    CalculationStatistics? statistics,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (statistics != null) result.statistics = statistics;
    return result;
  }

  GetCalculationExecutionStatusResponse._();

  factory GetCalculationExecutionStatusResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCalculationExecutionStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCalculationExecutionStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<CalculationStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: CalculationStatus.create)
    ..aOM<CalculationStatistics>(510636075, _omitFieldNames ? '' : 'statistics',
        subBuilder: CalculationStatistics.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCalculationExecutionStatusResponse copyWith(
          void Function(GetCalculationExecutionStatusResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetCalculationExecutionStatusResponse))
          as GetCalculationExecutionStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionStatusResponse create() =>
      GetCalculationExecutionStatusResponse._();
  @$core.override
  GetCalculationExecutionStatusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCalculationExecutionStatusResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetCalculationExecutionStatusResponse>(create);
  static GetCalculationExecutionStatusResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  CalculationStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(CalculationStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  CalculationStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(510636075)
  CalculationStatistics get statistics => $_getN(1);
  @$pb.TagNumber(510636075)
  set statistics(CalculationStatistics value) => $_setField(510636075, value);
  @$pb.TagNumber(510636075)
  $core.bool hasStatistics() => $_has(1);
  @$pb.TagNumber(510636075)
  void clearStatistics() => $_clearField(510636075);
  @$pb.TagNumber(510636075)
  CalculationStatistics ensureStatistics() => $_ensure(1);
}

class GetCapacityAssignmentConfigurationInput extends $pb.GeneratedMessage {
  factory GetCapacityAssignmentConfigurationInput({
    $core.String? capacityreservationname,
  }) {
    final result = create();
    if (capacityreservationname != null)
      result.capacityreservationname = capacityreservationname;
    return result;
  }

  GetCapacityAssignmentConfigurationInput._();

  factory GetCapacityAssignmentConfigurationInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCapacityAssignmentConfigurationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCapacityAssignmentConfigurationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(327567687, _omitFieldNames ? '' : 'capacityreservationname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityAssignmentConfigurationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityAssignmentConfigurationInput copyWith(
          void Function(GetCapacityAssignmentConfigurationInput) updates) =>
      super.copyWith((message) =>
              updates(message as GetCapacityAssignmentConfigurationInput))
          as GetCapacityAssignmentConfigurationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCapacityAssignmentConfigurationInput create() =>
      GetCapacityAssignmentConfigurationInput._();
  @$core.override
  GetCapacityAssignmentConfigurationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCapacityAssignmentConfigurationInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetCapacityAssignmentConfigurationInput>(create);
  static GetCapacityAssignmentConfigurationInput? _defaultInstance;

  @$pb.TagNumber(327567687)
  $core.String get capacityreservationname => $_getSZ(0);
  @$pb.TagNumber(327567687)
  set capacityreservationname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(327567687)
  $core.bool hasCapacityreservationname() => $_has(0);
  @$pb.TagNumber(327567687)
  void clearCapacityreservationname() => $_clearField(327567687);
}

class GetCapacityAssignmentConfigurationOutput extends $pb.GeneratedMessage {
  factory GetCapacityAssignmentConfigurationOutput({
    CapacityAssignmentConfiguration? capacityassignmentconfiguration,
  }) {
    final result = create();
    if (capacityassignmentconfiguration != null)
      result.capacityassignmentconfiguration = capacityassignmentconfiguration;
    return result;
  }

  GetCapacityAssignmentConfigurationOutput._();

  factory GetCapacityAssignmentConfigurationOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCapacityAssignmentConfigurationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCapacityAssignmentConfigurationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<CapacityAssignmentConfiguration>(
        9812735, _omitFieldNames ? '' : 'capacityassignmentconfiguration',
        subBuilder: CapacityAssignmentConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityAssignmentConfigurationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityAssignmentConfigurationOutput copyWith(
          void Function(GetCapacityAssignmentConfigurationOutput) updates) =>
      super.copyWith((message) =>
              updates(message as GetCapacityAssignmentConfigurationOutput))
          as GetCapacityAssignmentConfigurationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCapacityAssignmentConfigurationOutput create() =>
      GetCapacityAssignmentConfigurationOutput._();
  @$core.override
  GetCapacityAssignmentConfigurationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCapacityAssignmentConfigurationOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetCapacityAssignmentConfigurationOutput>(create);
  static GetCapacityAssignmentConfigurationOutput? _defaultInstance;

  @$pb.TagNumber(9812735)
  CapacityAssignmentConfiguration get capacityassignmentconfiguration =>
      $_getN(0);
  @$pb.TagNumber(9812735)
  set capacityassignmentconfiguration(CapacityAssignmentConfiguration value) =>
      $_setField(9812735, value);
  @$pb.TagNumber(9812735)
  $core.bool hasCapacityassignmentconfiguration() => $_has(0);
  @$pb.TagNumber(9812735)
  void clearCapacityassignmentconfiguration() => $_clearField(9812735);
  @$pb.TagNumber(9812735)
  CapacityAssignmentConfiguration ensureCapacityassignmentconfiguration() =>
      $_ensure(0);
}

class GetCapacityReservationInput extends $pb.GeneratedMessage {
  factory GetCapacityReservationInput({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  GetCapacityReservationInput._();

  factory GetCapacityReservationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCapacityReservationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCapacityReservationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityReservationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityReservationInput copyWith(
          void Function(GetCapacityReservationInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetCapacityReservationInput))
          as GetCapacityReservationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCapacityReservationInput create() =>
      GetCapacityReservationInput._();
  @$core.override
  GetCapacityReservationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCapacityReservationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCapacityReservationInput>(create);
  static GetCapacityReservationInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class GetCapacityReservationOutput extends $pb.GeneratedMessage {
  factory GetCapacityReservationOutput({
    CapacityReservation? capacityreservation,
  }) {
    final result = create();
    if (capacityreservation != null)
      result.capacityreservation = capacityreservation;
    return result;
  }

  GetCapacityReservationOutput._();

  factory GetCapacityReservationOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCapacityReservationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCapacityReservationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<CapacityReservation>(
        451869958, _omitFieldNames ? '' : 'capacityreservation',
        subBuilder: CapacityReservation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityReservationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCapacityReservationOutput copyWith(
          void Function(GetCapacityReservationOutput) updates) =>
      super.copyWith(
              (message) => updates(message as GetCapacityReservationOutput))
          as GetCapacityReservationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCapacityReservationOutput create() =>
      GetCapacityReservationOutput._();
  @$core.override
  GetCapacityReservationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCapacityReservationOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCapacityReservationOutput>(create);
  static GetCapacityReservationOutput? _defaultInstance;

  @$pb.TagNumber(451869958)
  CapacityReservation get capacityreservation => $_getN(0);
  @$pb.TagNumber(451869958)
  set capacityreservation(CapacityReservation value) =>
      $_setField(451869958, value);
  @$pb.TagNumber(451869958)
  $core.bool hasCapacityreservation() => $_has(0);
  @$pb.TagNumber(451869958)
  void clearCapacityreservation() => $_clearField(451869958);
  @$pb.TagNumber(451869958)
  CapacityReservation ensureCapacityreservation() => $_ensure(0);
}

class GetDataCatalogInput extends $pb.GeneratedMessage {
  factory GetDataCatalogInput({
    $core.String? name,
    $core.String? workgroup,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  GetDataCatalogInput._();

  factory GetDataCatalogInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDataCatalogInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDataCatalogInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataCatalogInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataCatalogInput copyWith(void Function(GetDataCatalogInput) updates) =>
      super.copyWith((message) => updates(message as GetDataCatalogInput))
          as GetDataCatalogInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDataCatalogInput create() => GetDataCatalogInput._();
  @$core.override
  GetDataCatalogInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDataCatalogInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDataCatalogInput>(create);
  static GetDataCatalogInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(1);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(1);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class GetDataCatalogOutput extends $pb.GeneratedMessage {
  factory GetDataCatalogOutput({
    DataCatalog? datacatalog,
  }) {
    final result = create();
    if (datacatalog != null) result.datacatalog = datacatalog;
    return result;
  }

  GetDataCatalogOutput._();

  factory GetDataCatalogOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDataCatalogOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDataCatalogOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<DataCatalog>(209694649, _omitFieldNames ? '' : 'datacatalog',
        subBuilder: DataCatalog.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataCatalogOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataCatalogOutput copyWith(void Function(GetDataCatalogOutput) updates) =>
      super.copyWith((message) => updates(message as GetDataCatalogOutput))
          as GetDataCatalogOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDataCatalogOutput create() => GetDataCatalogOutput._();
  @$core.override
  GetDataCatalogOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDataCatalogOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDataCatalogOutput>(create);
  static GetDataCatalogOutput? _defaultInstance;

  @$pb.TagNumber(209694649)
  DataCatalog get datacatalog => $_getN(0);
  @$pb.TagNumber(209694649)
  set datacatalog(DataCatalog value) => $_setField(209694649, value);
  @$pb.TagNumber(209694649)
  $core.bool hasDatacatalog() => $_has(0);
  @$pb.TagNumber(209694649)
  void clearDatacatalog() => $_clearField(209694649);
  @$pb.TagNumber(209694649)
  DataCatalog ensureDatacatalog() => $_ensure(0);
}

class GetDatabaseInput extends $pb.GeneratedMessage {
  factory GetDatabaseInput({
    $core.String? databasename,
    $core.String? workgroup,
    $core.String? catalogname,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (workgroup != null) result.workgroup = workgroup;
    if (catalogname != null) result.catalogname = catalogname;
    return result;
  }

  GetDatabaseInput._();

  factory GetDatabaseInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDatabaseInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDatabaseInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..aOS(518825212, _omitFieldNames ? '' : 'catalogname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDatabaseInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDatabaseInput copyWith(void Function(GetDatabaseInput) updates) =>
      super.copyWith((message) => updates(message as GetDatabaseInput))
          as GetDatabaseInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDatabaseInput create() => GetDatabaseInput._();
  @$core.override
  GetDatabaseInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDatabaseInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDatabaseInput>(create);
  static GetDatabaseInput? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(1);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(1);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(518825212)
  $core.String get catalogname => $_getSZ(2);
  @$pb.TagNumber(518825212)
  set catalogname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518825212)
  $core.bool hasCatalogname() => $_has(2);
  @$pb.TagNumber(518825212)
  void clearCatalogname() => $_clearField(518825212);
}

class GetDatabaseOutput extends $pb.GeneratedMessage {
  factory GetDatabaseOutput({
    Database? database,
  }) {
    final result = create();
    if (database != null) result.database = database;
    return result;
  }

  GetDatabaseOutput._();

  factory GetDatabaseOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDatabaseOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDatabaseOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<Database>(278147289, _omitFieldNames ? '' : 'database',
        subBuilder: Database.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDatabaseOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDatabaseOutput copyWith(void Function(GetDatabaseOutput) updates) =>
      super.copyWith((message) => updates(message as GetDatabaseOutput))
          as GetDatabaseOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDatabaseOutput create() => GetDatabaseOutput._();
  @$core.override
  GetDatabaseOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDatabaseOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDatabaseOutput>(create);
  static GetDatabaseOutput? _defaultInstance;

  @$pb.TagNumber(278147289)
  Database get database => $_getN(0);
  @$pb.TagNumber(278147289)
  set database(Database value) => $_setField(278147289, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(0);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);
  @$pb.TagNumber(278147289)
  Database ensureDatabase() => $_ensure(0);
}

class GetNamedQueryInput extends $pb.GeneratedMessage {
  factory GetNamedQueryInput({
    $core.String? namedqueryid,
  }) {
    final result = create();
    if (namedqueryid != null) result.namedqueryid = namedqueryid;
    return result;
  }

  GetNamedQueryInput._();

  factory GetNamedQueryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNamedQueryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNamedQueryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(330896872, _omitFieldNames ? '' : 'namedqueryid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNamedQueryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNamedQueryInput copyWith(void Function(GetNamedQueryInput) updates) =>
      super.copyWith((message) => updates(message as GetNamedQueryInput))
          as GetNamedQueryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNamedQueryInput create() => GetNamedQueryInput._();
  @$core.override
  GetNamedQueryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNamedQueryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNamedQueryInput>(create);
  static GetNamedQueryInput? _defaultInstance;

  @$pb.TagNumber(330896872)
  $core.String get namedqueryid => $_getSZ(0);
  @$pb.TagNumber(330896872)
  set namedqueryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(330896872)
  $core.bool hasNamedqueryid() => $_has(0);
  @$pb.TagNumber(330896872)
  void clearNamedqueryid() => $_clearField(330896872);
}

class GetNamedQueryOutput extends $pb.GeneratedMessage {
  factory GetNamedQueryOutput({
    NamedQuery? namedquery,
  }) {
    final result = create();
    if (namedquery != null) result.namedquery = namedquery;
    return result;
  }

  GetNamedQueryOutput._();

  factory GetNamedQueryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNamedQueryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNamedQueryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<NamedQuery>(170671691, _omitFieldNames ? '' : 'namedquery',
        subBuilder: NamedQuery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNamedQueryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNamedQueryOutput copyWith(void Function(GetNamedQueryOutput) updates) =>
      super.copyWith((message) => updates(message as GetNamedQueryOutput))
          as GetNamedQueryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNamedQueryOutput create() => GetNamedQueryOutput._();
  @$core.override
  GetNamedQueryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNamedQueryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNamedQueryOutput>(create);
  static GetNamedQueryOutput? _defaultInstance;

  @$pb.TagNumber(170671691)
  NamedQuery get namedquery => $_getN(0);
  @$pb.TagNumber(170671691)
  set namedquery(NamedQuery value) => $_setField(170671691, value);
  @$pb.TagNumber(170671691)
  $core.bool hasNamedquery() => $_has(0);
  @$pb.TagNumber(170671691)
  void clearNamedquery() => $_clearField(170671691);
  @$pb.TagNumber(170671691)
  NamedQuery ensureNamedquery() => $_ensure(0);
}

class GetNotebookMetadataInput extends $pb.GeneratedMessage {
  factory GetNotebookMetadataInput({
    $core.String? notebookid,
  }) {
    final result = create();
    if (notebookid != null) result.notebookid = notebookid;
    return result;
  }

  GetNotebookMetadataInput._();

  factory GetNotebookMetadataInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNotebookMetadataInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNotebookMetadataInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNotebookMetadataInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNotebookMetadataInput copyWith(
          void Function(GetNotebookMetadataInput) updates) =>
      super.copyWith((message) => updates(message as GetNotebookMetadataInput))
          as GetNotebookMetadataInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNotebookMetadataInput create() => GetNotebookMetadataInput._();
  @$core.override
  GetNotebookMetadataInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNotebookMetadataInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNotebookMetadataInput>(create);
  static GetNotebookMetadataInput? _defaultInstance;

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(0);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(0);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);
}

class GetNotebookMetadataOutput extends $pb.GeneratedMessage {
  factory GetNotebookMetadataOutput({
    NotebookMetadata? notebookmetadata,
  }) {
    final result = create();
    if (notebookmetadata != null) result.notebookmetadata = notebookmetadata;
    return result;
  }

  GetNotebookMetadataOutput._();

  factory GetNotebookMetadataOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNotebookMetadataOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNotebookMetadataOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<NotebookMetadata>(77536390, _omitFieldNames ? '' : 'notebookmetadata',
        subBuilder: NotebookMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNotebookMetadataOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNotebookMetadataOutput copyWith(
          void Function(GetNotebookMetadataOutput) updates) =>
      super.copyWith((message) => updates(message as GetNotebookMetadataOutput))
          as GetNotebookMetadataOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNotebookMetadataOutput create() => GetNotebookMetadataOutput._();
  @$core.override
  GetNotebookMetadataOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNotebookMetadataOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNotebookMetadataOutput>(create);
  static GetNotebookMetadataOutput? _defaultInstance;

  @$pb.TagNumber(77536390)
  NotebookMetadata get notebookmetadata => $_getN(0);
  @$pb.TagNumber(77536390)
  set notebookmetadata(NotebookMetadata value) => $_setField(77536390, value);
  @$pb.TagNumber(77536390)
  $core.bool hasNotebookmetadata() => $_has(0);
  @$pb.TagNumber(77536390)
  void clearNotebookmetadata() => $_clearField(77536390);
  @$pb.TagNumber(77536390)
  NotebookMetadata ensureNotebookmetadata() => $_ensure(0);
}

class GetPreparedStatementInput extends $pb.GeneratedMessage {
  factory GetPreparedStatementInput({
    $core.String? statementname,
    $core.String? workgroup,
  }) {
    final result = create();
    if (statementname != null) result.statementname = statementname;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  GetPreparedStatementInput._();

  factory GetPreparedStatementInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPreparedStatementInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPreparedStatementInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPreparedStatementInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPreparedStatementInput copyWith(
          void Function(GetPreparedStatementInput) updates) =>
      super.copyWith((message) => updates(message as GetPreparedStatementInput))
          as GetPreparedStatementInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPreparedStatementInput create() => GetPreparedStatementInput._();
  @$core.override
  GetPreparedStatementInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPreparedStatementInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPreparedStatementInput>(create);
  static GetPreparedStatementInput? _defaultInstance;

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(0);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(0);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(1);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(1);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class GetPreparedStatementOutput extends $pb.GeneratedMessage {
  factory GetPreparedStatementOutput({
    PreparedStatement? preparedstatement,
  }) {
    final result = create();
    if (preparedstatement != null) result.preparedstatement = preparedstatement;
    return result;
  }

  GetPreparedStatementOutput._();

  factory GetPreparedStatementOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPreparedStatementOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPreparedStatementOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<PreparedStatement>(
        164916502, _omitFieldNames ? '' : 'preparedstatement',
        subBuilder: PreparedStatement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPreparedStatementOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPreparedStatementOutput copyWith(
          void Function(GetPreparedStatementOutput) updates) =>
      super.copyWith(
              (message) => updates(message as GetPreparedStatementOutput))
          as GetPreparedStatementOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPreparedStatementOutput create() => GetPreparedStatementOutput._();
  @$core.override
  GetPreparedStatementOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPreparedStatementOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPreparedStatementOutput>(create);
  static GetPreparedStatementOutput? _defaultInstance;

  @$pb.TagNumber(164916502)
  PreparedStatement get preparedstatement => $_getN(0);
  @$pb.TagNumber(164916502)
  set preparedstatement(PreparedStatement value) =>
      $_setField(164916502, value);
  @$pb.TagNumber(164916502)
  $core.bool hasPreparedstatement() => $_has(0);
  @$pb.TagNumber(164916502)
  void clearPreparedstatement() => $_clearField(164916502);
  @$pb.TagNumber(164916502)
  PreparedStatement ensurePreparedstatement() => $_ensure(0);
}

class GetQueryExecutionInput extends $pb.GeneratedMessage {
  factory GetQueryExecutionInput({
    $core.String? queryexecutionid,
  }) {
    final result = create();
    if (queryexecutionid != null) result.queryexecutionid = queryexecutionid;
    return result;
  }

  GetQueryExecutionInput._();

  factory GetQueryExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(467615503, _omitFieldNames ? '' : 'queryexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryExecutionInput copyWith(
          void Function(GetQueryExecutionInput) updates) =>
      super.copyWith((message) => updates(message as GetQueryExecutionInput))
          as GetQueryExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryExecutionInput create() => GetQueryExecutionInput._();
  @$core.override
  GetQueryExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryExecutionInput>(create);
  static GetQueryExecutionInput? _defaultInstance;

  @$pb.TagNumber(467615503)
  $core.String get queryexecutionid => $_getSZ(0);
  @$pb.TagNumber(467615503)
  set queryexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(467615503)
  $core.bool hasQueryexecutionid() => $_has(0);
  @$pb.TagNumber(467615503)
  void clearQueryexecutionid() => $_clearField(467615503);
}

class GetQueryExecutionOutput extends $pb.GeneratedMessage {
  factory GetQueryExecutionOutput({
    QueryExecution? queryexecution,
  }) {
    final result = create();
    if (queryexecution != null) result.queryexecution = queryexecution;
    return result;
  }

  GetQueryExecutionOutput._();

  factory GetQueryExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<QueryExecution>(62173540, _omitFieldNames ? '' : 'queryexecution',
        subBuilder: QueryExecution.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryExecutionOutput copyWith(
          void Function(GetQueryExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as GetQueryExecutionOutput))
          as GetQueryExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryExecutionOutput create() => GetQueryExecutionOutput._();
  @$core.override
  GetQueryExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryExecutionOutput>(create);
  static GetQueryExecutionOutput? _defaultInstance;

  @$pb.TagNumber(62173540)
  QueryExecution get queryexecution => $_getN(0);
  @$pb.TagNumber(62173540)
  set queryexecution(QueryExecution value) => $_setField(62173540, value);
  @$pb.TagNumber(62173540)
  $core.bool hasQueryexecution() => $_has(0);
  @$pb.TagNumber(62173540)
  void clearQueryexecution() => $_clearField(62173540);
  @$pb.TagNumber(62173540)
  QueryExecution ensureQueryexecution() => $_ensure(0);
}

class GetQueryResultsInput extends $pb.GeneratedMessage {
  factory GetQueryResultsInput({
    QueryResultType? queryresulttype,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? queryexecutionid,
  }) {
    final result = create();
    if (queryresulttype != null) result.queryresulttype = queryresulttype;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (queryexecutionid != null) result.queryexecutionid = queryexecutionid;
    return result;
  }

  GetQueryResultsInput._();

  factory GetQueryResultsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryResultsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryResultsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<QueryResultType>(3724043, _omitFieldNames ? '' : 'queryresulttype',
        enumValues: QueryResultType.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(467615503, _omitFieldNames ? '' : 'queryexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsInput copyWith(void Function(GetQueryResultsInput) updates) =>
      super.copyWith((message) => updates(message as GetQueryResultsInput))
          as GetQueryResultsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryResultsInput create() => GetQueryResultsInput._();
  @$core.override
  GetQueryResultsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryResultsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryResultsInput>(create);
  static GetQueryResultsInput? _defaultInstance;

  @$pb.TagNumber(3724043)
  QueryResultType get queryresulttype => $_getN(0);
  @$pb.TagNumber(3724043)
  set queryresulttype(QueryResultType value) => $_setField(3724043, value);
  @$pb.TagNumber(3724043)
  $core.bool hasQueryresulttype() => $_has(0);
  @$pb.TagNumber(3724043)
  void clearQueryresulttype() => $_clearField(3724043);

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

  @$pb.TagNumber(467615503)
  $core.String get queryexecutionid => $_getSZ(3);
  @$pb.TagNumber(467615503)
  set queryexecutionid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(467615503)
  $core.bool hasQueryexecutionid() => $_has(3);
  @$pb.TagNumber(467615503)
  void clearQueryexecutionid() => $_clearField(467615503);
}

class GetQueryResultsOutput extends $pb.GeneratedMessage {
  factory GetQueryResultsOutput({
    ResultSet? resultset,
    $core.String? nexttoken,
    $fixnum.Int64? updatecount,
  }) {
    final result = create();
    if (resultset != null) result.resultset = resultset;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (updatecount != null) result.updatecount = updatecount;
    return result;
  }

  GetQueryResultsOutput._();

  factory GetQueryResultsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryResultsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryResultsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<ResultSet>(146844365, _omitFieldNames ? '' : 'resultset',
        subBuilder: ResultSet.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aInt64(312282732, _omitFieldNames ? '' : 'updatecount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryResultsOutput copyWith(
          void Function(GetQueryResultsOutput) updates) =>
      super.copyWith((message) => updates(message as GetQueryResultsOutput))
          as GetQueryResultsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryResultsOutput create() => GetQueryResultsOutput._();
  @$core.override
  GetQueryResultsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryResultsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryResultsOutput>(create);
  static GetQueryResultsOutput? _defaultInstance;

  @$pb.TagNumber(146844365)
  ResultSet get resultset => $_getN(0);
  @$pb.TagNumber(146844365)
  set resultset(ResultSet value) => $_setField(146844365, value);
  @$pb.TagNumber(146844365)
  $core.bool hasResultset() => $_has(0);
  @$pb.TagNumber(146844365)
  void clearResultset() => $_clearField(146844365);
  @$pb.TagNumber(146844365)
  ResultSet ensureResultset() => $_ensure(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(312282732)
  $fixnum.Int64 get updatecount => $_getI64(2);
  @$pb.TagNumber(312282732)
  set updatecount($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(312282732)
  $core.bool hasUpdatecount() => $_has(2);
  @$pb.TagNumber(312282732)
  void clearUpdatecount() => $_clearField(312282732);
}

class GetQueryRuntimeStatisticsInput extends $pb.GeneratedMessage {
  factory GetQueryRuntimeStatisticsInput({
    $core.String? queryexecutionid,
  }) {
    final result = create();
    if (queryexecutionid != null) result.queryexecutionid = queryexecutionid;
    return result;
  }

  GetQueryRuntimeStatisticsInput._();

  factory GetQueryRuntimeStatisticsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryRuntimeStatisticsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryRuntimeStatisticsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(467615503, _omitFieldNames ? '' : 'queryexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryRuntimeStatisticsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryRuntimeStatisticsInput copyWith(
          void Function(GetQueryRuntimeStatisticsInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetQueryRuntimeStatisticsInput))
          as GetQueryRuntimeStatisticsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryRuntimeStatisticsInput create() =>
      GetQueryRuntimeStatisticsInput._();
  @$core.override
  GetQueryRuntimeStatisticsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryRuntimeStatisticsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryRuntimeStatisticsInput>(create);
  static GetQueryRuntimeStatisticsInput? _defaultInstance;

  @$pb.TagNumber(467615503)
  $core.String get queryexecutionid => $_getSZ(0);
  @$pb.TagNumber(467615503)
  set queryexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(467615503)
  $core.bool hasQueryexecutionid() => $_has(0);
  @$pb.TagNumber(467615503)
  void clearQueryexecutionid() => $_clearField(467615503);
}

class GetQueryRuntimeStatisticsOutput extends $pb.GeneratedMessage {
  factory GetQueryRuntimeStatisticsOutput({
    QueryRuntimeStatistics? queryruntimestatistics,
  }) {
    final result = create();
    if (queryruntimestatistics != null)
      result.queryruntimestatistics = queryruntimestatistics;
    return result;
  }

  GetQueryRuntimeStatisticsOutput._();

  factory GetQueryRuntimeStatisticsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryRuntimeStatisticsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryRuntimeStatisticsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<QueryRuntimeStatistics>(
        157051307, _omitFieldNames ? '' : 'queryruntimestatistics',
        subBuilder: QueryRuntimeStatistics.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryRuntimeStatisticsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryRuntimeStatisticsOutput copyWith(
          void Function(GetQueryRuntimeStatisticsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as GetQueryRuntimeStatisticsOutput))
          as GetQueryRuntimeStatisticsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryRuntimeStatisticsOutput create() =>
      GetQueryRuntimeStatisticsOutput._();
  @$core.override
  GetQueryRuntimeStatisticsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryRuntimeStatisticsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryRuntimeStatisticsOutput>(
          create);
  static GetQueryRuntimeStatisticsOutput? _defaultInstance;

  @$pb.TagNumber(157051307)
  QueryRuntimeStatistics get queryruntimestatistics => $_getN(0);
  @$pb.TagNumber(157051307)
  set queryruntimestatistics(QueryRuntimeStatistics value) =>
      $_setField(157051307, value);
  @$pb.TagNumber(157051307)
  $core.bool hasQueryruntimestatistics() => $_has(0);
  @$pb.TagNumber(157051307)
  void clearQueryruntimestatistics() => $_clearField(157051307);
  @$pb.TagNumber(157051307)
  QueryRuntimeStatistics ensureQueryruntimestatistics() => $_ensure(0);
}

class GetResourceDashboardRequest extends $pb.GeneratedMessage {
  factory GetResourceDashboardRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  GetResourceDashboardRequest._();

  factory GetResourceDashboardRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourceDashboardRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourceDashboardRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourceDashboardRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourceDashboardRequest copyWith(
          void Function(GetResourceDashboardRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetResourceDashboardRequest))
          as GetResourceDashboardRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourceDashboardRequest create() =>
      GetResourceDashboardRequest._();
  @$core.override
  GetResourceDashboardRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourceDashboardRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourceDashboardRequest>(create);
  static GetResourceDashboardRequest? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class GetResourceDashboardResponse extends $pb.GeneratedMessage {
  factory GetResourceDashboardResponse({
    $core.String? url,
  }) {
    final result = create();
    if (url != null) result.url = url;
    return result;
  }

  GetResourceDashboardResponse._();

  factory GetResourceDashboardResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourceDashboardResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourceDashboardResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(354018239, _omitFieldNames ? '' : 'url')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourceDashboardResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourceDashboardResponse copyWith(
          void Function(GetResourceDashboardResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetResourceDashboardResponse))
          as GetResourceDashboardResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourceDashboardResponse create() =>
      GetResourceDashboardResponse._();
  @$core.override
  GetResourceDashboardResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourceDashboardResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourceDashboardResponse>(create);
  static GetResourceDashboardResponse? _defaultInstance;

  @$pb.TagNumber(354018239)
  $core.String get url => $_getSZ(0);
  @$pb.TagNumber(354018239)
  set url($core.String value) => $_setString(0, value);
  @$pb.TagNumber(354018239)
  $core.bool hasUrl() => $_has(0);
  @$pb.TagNumber(354018239)
  void clearUrl() => $_clearField(354018239);
}

class GetSessionEndpointRequest extends $pb.GeneratedMessage {
  factory GetSessionEndpointRequest({
    $core.String? sessionid,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    return result;
  }

  GetSessionEndpointRequest._();

  factory GetSessionEndpointRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionEndpointRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionEndpointRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionEndpointRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionEndpointRequest copyWith(
          void Function(GetSessionEndpointRequest) updates) =>
      super.copyWith((message) => updates(message as GetSessionEndpointRequest))
          as GetSessionEndpointRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionEndpointRequest create() => GetSessionEndpointRequest._();
  @$core.override
  GetSessionEndpointRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionEndpointRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionEndpointRequest>(create);
  static GetSessionEndpointRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);
}

class GetSessionEndpointResponse extends $pb.GeneratedMessage {
  factory GetSessionEndpointResponse({
    $core.String? endpointurl,
    $core.String? authtoken,
    $core.String? authtokenexpirationtime,
  }) {
    final result = create();
    if (endpointurl != null) result.endpointurl = endpointurl;
    if (authtoken != null) result.authtoken = authtoken;
    if (authtokenexpirationtime != null)
      result.authtokenexpirationtime = authtokenexpirationtime;
    return result;
  }

  GetSessionEndpointResponse._();

  factory GetSessionEndpointResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionEndpointResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionEndpointResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(31787414, _omitFieldNames ? '' : 'endpointurl')
    ..aOS(349771493, _omitFieldNames ? '' : 'authtoken')
    ..aOS(406065689, _omitFieldNames ? '' : 'authtokenexpirationtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionEndpointResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionEndpointResponse copyWith(
          void Function(GetSessionEndpointResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetSessionEndpointResponse))
          as GetSessionEndpointResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionEndpointResponse create() => GetSessionEndpointResponse._();
  @$core.override
  GetSessionEndpointResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionEndpointResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionEndpointResponse>(create);
  static GetSessionEndpointResponse? _defaultInstance;

  @$pb.TagNumber(31787414)
  $core.String get endpointurl => $_getSZ(0);
  @$pb.TagNumber(31787414)
  set endpointurl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(31787414)
  $core.bool hasEndpointurl() => $_has(0);
  @$pb.TagNumber(31787414)
  void clearEndpointurl() => $_clearField(31787414);

  @$pb.TagNumber(349771493)
  $core.String get authtoken => $_getSZ(1);
  @$pb.TagNumber(349771493)
  set authtoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(349771493)
  $core.bool hasAuthtoken() => $_has(1);
  @$pb.TagNumber(349771493)
  void clearAuthtoken() => $_clearField(349771493);

  @$pb.TagNumber(406065689)
  $core.String get authtokenexpirationtime => $_getSZ(2);
  @$pb.TagNumber(406065689)
  set authtokenexpirationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(406065689)
  $core.bool hasAuthtokenexpirationtime() => $_has(2);
  @$pb.TagNumber(406065689)
  void clearAuthtokenexpirationtime() => $_clearField(406065689);
}

class GetSessionRequest extends $pb.GeneratedMessage {
  factory GetSessionRequest({
    $core.String? sessionid,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    return result;
  }

  GetSessionRequest._();

  factory GetSessionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionRequest copyWith(void Function(GetSessionRequest) updates) =>
      super.copyWith((message) => updates(message as GetSessionRequest))
          as GetSessionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionRequest create() => GetSessionRequest._();
  @$core.override
  GetSessionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionRequest>(create);
  static GetSessionRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);
}

class GetSessionResponse extends $pb.GeneratedMessage {
  factory GetSessionResponse({
    SessionStatus? status,
    $core.String? sessionid,
    $core.String? engineversion,
    $core.String? description,
    SessionConfiguration? sessionconfiguration,
    EngineConfiguration? engineconfiguration,
    MonitoringConfiguration? monitoringconfiguration,
    $core.String? workgroup,
    SessionStatistics? statistics,
    $core.String? notebookversion,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (sessionid != null) result.sessionid = sessionid;
    if (engineversion != null) result.engineversion = engineversion;
    if (description != null) result.description = description;
    if (sessionconfiguration != null)
      result.sessionconfiguration = sessionconfiguration;
    if (engineconfiguration != null)
      result.engineconfiguration = engineconfiguration;
    if (monitoringconfiguration != null)
      result.monitoringconfiguration = monitoringconfiguration;
    if (workgroup != null) result.workgroup = workgroup;
    if (statistics != null) result.statistics = statistics;
    if (notebookversion != null) result.notebookversion = notebookversion;
    return result;
  }

  GetSessionResponse._();

  factory GetSessionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<SessionStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: SessionStatus.create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aOS(44953462, _omitFieldNames ? '' : 'engineversion')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOM<SessionConfiguration>(
        211592776, _omitFieldNames ? '' : 'sessionconfiguration',
        subBuilder: SessionConfiguration.create)
    ..aOM<EngineConfiguration>(
        341629412, _omitFieldNames ? '' : 'engineconfiguration',
        subBuilder: EngineConfiguration.create)
    ..aOM<MonitoringConfiguration>(
        364891928, _omitFieldNames ? '' : 'monitoringconfiguration',
        subBuilder: MonitoringConfiguration.create)
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..aOM<SessionStatistics>(510636075, _omitFieldNames ? '' : 'statistics',
        subBuilder: SessionStatistics.create)
    ..aOS(528689837, _omitFieldNames ? '' : 'notebookversion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionResponse copyWith(void Function(GetSessionResponse) updates) =>
      super.copyWith((message) => updates(message as GetSessionResponse))
          as GetSessionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionResponse create() => GetSessionResponse._();
  @$core.override
  GetSessionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionResponse>(create);
  static GetSessionResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  SessionStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(SessionStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  SessionStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(1);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(1);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(44953462)
  $core.String get engineversion => $_getSZ(2);
  @$pb.TagNumber(44953462)
  set engineversion($core.String value) => $_setString(2, value);
  @$pb.TagNumber(44953462)
  $core.bool hasEngineversion() => $_has(2);
  @$pb.TagNumber(44953462)
  void clearEngineversion() => $_clearField(44953462);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(211592776)
  SessionConfiguration get sessionconfiguration => $_getN(4);
  @$pb.TagNumber(211592776)
  set sessionconfiguration(SessionConfiguration value) =>
      $_setField(211592776, value);
  @$pb.TagNumber(211592776)
  $core.bool hasSessionconfiguration() => $_has(4);
  @$pb.TagNumber(211592776)
  void clearSessionconfiguration() => $_clearField(211592776);
  @$pb.TagNumber(211592776)
  SessionConfiguration ensureSessionconfiguration() => $_ensure(4);

  @$pb.TagNumber(341629412)
  EngineConfiguration get engineconfiguration => $_getN(5);
  @$pb.TagNumber(341629412)
  set engineconfiguration(EngineConfiguration value) =>
      $_setField(341629412, value);
  @$pb.TagNumber(341629412)
  $core.bool hasEngineconfiguration() => $_has(5);
  @$pb.TagNumber(341629412)
  void clearEngineconfiguration() => $_clearField(341629412);
  @$pb.TagNumber(341629412)
  EngineConfiguration ensureEngineconfiguration() => $_ensure(5);

  @$pb.TagNumber(364891928)
  MonitoringConfiguration get monitoringconfiguration => $_getN(6);
  @$pb.TagNumber(364891928)
  set monitoringconfiguration(MonitoringConfiguration value) =>
      $_setField(364891928, value);
  @$pb.TagNumber(364891928)
  $core.bool hasMonitoringconfiguration() => $_has(6);
  @$pb.TagNumber(364891928)
  void clearMonitoringconfiguration() => $_clearField(364891928);
  @$pb.TagNumber(364891928)
  MonitoringConfiguration ensureMonitoringconfiguration() => $_ensure(6);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(7);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(7, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(7);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(510636075)
  SessionStatistics get statistics => $_getN(8);
  @$pb.TagNumber(510636075)
  set statistics(SessionStatistics value) => $_setField(510636075, value);
  @$pb.TagNumber(510636075)
  $core.bool hasStatistics() => $_has(8);
  @$pb.TagNumber(510636075)
  void clearStatistics() => $_clearField(510636075);
  @$pb.TagNumber(510636075)
  SessionStatistics ensureStatistics() => $_ensure(8);

  @$pb.TagNumber(528689837)
  $core.String get notebookversion => $_getSZ(9);
  @$pb.TagNumber(528689837)
  set notebookversion($core.String value) => $_setString(9, value);
  @$pb.TagNumber(528689837)
  $core.bool hasNotebookversion() => $_has(9);
  @$pb.TagNumber(528689837)
  void clearNotebookversion() => $_clearField(528689837);
}

class GetSessionStatusRequest extends $pb.GeneratedMessage {
  factory GetSessionStatusRequest({
    $core.String? sessionid,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    return result;
  }

  GetSessionStatusRequest._();

  factory GetSessionStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionStatusRequest copyWith(
          void Function(GetSessionStatusRequest) updates) =>
      super.copyWith((message) => updates(message as GetSessionStatusRequest))
          as GetSessionStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionStatusRequest create() => GetSessionStatusRequest._();
  @$core.override
  GetSessionStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionStatusRequest>(create);
  static GetSessionStatusRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);
}

class GetSessionStatusResponse extends $pb.GeneratedMessage {
  factory GetSessionStatusResponse({
    SessionStatus? status,
    $core.String? sessionid,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (sessionid != null) result.sessionid = sessionid;
    return result;
  }

  GetSessionStatusResponse._();

  factory GetSessionStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<SessionStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: SessionStatus.create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionStatusResponse copyWith(
          void Function(GetSessionStatusResponse) updates) =>
      super.copyWith((message) => updates(message as GetSessionStatusResponse))
          as GetSessionStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionStatusResponse create() => GetSessionStatusResponse._();
  @$core.override
  GetSessionStatusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionStatusResponse>(create);
  static GetSessionStatusResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  SessionStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(SessionStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  SessionStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(1);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(1);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);
}

class GetTableMetadataInput extends $pb.GeneratedMessage {
  factory GetTableMetadataInput({
    $core.String? databasename,
    $core.String? tablename,
    $core.String? workgroup,
    $core.String? catalogname,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (tablename != null) result.tablename = tablename;
    if (workgroup != null) result.workgroup = workgroup;
    if (catalogname != null) result.catalogname = catalogname;
    return result;
  }

  GetTableMetadataInput._();

  factory GetTableMetadataInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTableMetadataInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTableMetadataInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..aOS(518825212, _omitFieldNames ? '' : 'catalogname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTableMetadataInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTableMetadataInput copyWith(
          void Function(GetTableMetadataInput) updates) =>
      super.copyWith((message) => updates(message as GetTableMetadataInput))
          as GetTableMetadataInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTableMetadataInput create() => GetTableMetadataInput._();
  @$core.override
  GetTableMetadataInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTableMetadataInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTableMetadataInput>(create);
  static GetTableMetadataInput? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(2);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(2);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(518825212)
  $core.String get catalogname => $_getSZ(3);
  @$pb.TagNumber(518825212)
  set catalogname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(518825212)
  $core.bool hasCatalogname() => $_has(3);
  @$pb.TagNumber(518825212)
  void clearCatalogname() => $_clearField(518825212);
}

class GetTableMetadataOutput extends $pb.GeneratedMessage {
  factory GetTableMetadataOutput({
    TableMetadata? tablemetadata,
  }) {
    final result = create();
    if (tablemetadata != null) result.tablemetadata = tablemetadata;
    return result;
  }

  GetTableMetadataOutput._();

  factory GetTableMetadataOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTableMetadataOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTableMetadataOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<TableMetadata>(393071395, _omitFieldNames ? '' : 'tablemetadata',
        subBuilder: TableMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTableMetadataOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTableMetadataOutput copyWith(
          void Function(GetTableMetadataOutput) updates) =>
      super.copyWith((message) => updates(message as GetTableMetadataOutput))
          as GetTableMetadataOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTableMetadataOutput create() => GetTableMetadataOutput._();
  @$core.override
  GetTableMetadataOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTableMetadataOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTableMetadataOutput>(create);
  static GetTableMetadataOutput? _defaultInstance;

  @$pb.TagNumber(393071395)
  TableMetadata get tablemetadata => $_getN(0);
  @$pb.TagNumber(393071395)
  set tablemetadata(TableMetadata value) => $_setField(393071395, value);
  @$pb.TagNumber(393071395)
  $core.bool hasTablemetadata() => $_has(0);
  @$pb.TagNumber(393071395)
  void clearTablemetadata() => $_clearField(393071395);
  @$pb.TagNumber(393071395)
  TableMetadata ensureTablemetadata() => $_ensure(0);
}

class GetWorkGroupInput extends $pb.GeneratedMessage {
  factory GetWorkGroupInput({
    $core.String? workgroup,
  }) {
    final result = create();
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  GetWorkGroupInput._();

  factory GetWorkGroupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWorkGroupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWorkGroupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWorkGroupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWorkGroupInput copyWith(void Function(GetWorkGroupInput) updates) =>
      super.copyWith((message) => updates(message as GetWorkGroupInput))
          as GetWorkGroupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWorkGroupInput create() => GetWorkGroupInput._();
  @$core.override
  GetWorkGroupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWorkGroupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWorkGroupInput>(create);
  static GetWorkGroupInput? _defaultInstance;

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(0);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(0, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(0);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class GetWorkGroupOutput extends $pb.GeneratedMessage {
  factory GetWorkGroupOutput({
    WorkGroup? workgroup,
  }) {
    final result = create();
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  GetWorkGroupOutput._();

  factory GetWorkGroupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWorkGroupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWorkGroupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<WorkGroup>(505960068, _omitFieldNames ? '' : 'workgroup',
        subBuilder: WorkGroup.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWorkGroupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWorkGroupOutput copyWith(void Function(GetWorkGroupOutput) updates) =>
      super.copyWith((message) => updates(message as GetWorkGroupOutput))
          as GetWorkGroupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWorkGroupOutput create() => GetWorkGroupOutput._();
  @$core.override
  GetWorkGroupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWorkGroupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWorkGroupOutput>(create);
  static GetWorkGroupOutput? _defaultInstance;

  @$pb.TagNumber(505960068)
  WorkGroup get workgroup => $_getN(0);
  @$pb.TagNumber(505960068)
  set workgroup(WorkGroup value) => $_setField(505960068, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(0);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
  @$pb.TagNumber(505960068)
  WorkGroup ensureWorkgroup() => $_ensure(0);
}

class IdentityCenterConfiguration extends $pb.GeneratedMessage {
  factory IdentityCenterConfiguration({
    $core.bool? enableidentitycenter,
    $core.String? identitycenterinstancearn,
  }) {
    final result = create();
    if (enableidentitycenter != null)
      result.enableidentitycenter = enableidentitycenter;
    if (identitycenterinstancearn != null)
      result.identitycenterinstancearn = identitycenterinstancearn;
    return result;
  }

  IdentityCenterConfiguration._();

  factory IdentityCenterConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IdentityCenterConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IdentityCenterConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOB(418515238, _omitFieldNames ? '' : 'enableidentitycenter')
    ..aOS(469575873, _omitFieldNames ? '' : 'identitycenterinstancearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityCenterConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityCenterConfiguration copyWith(
          void Function(IdentityCenterConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as IdentityCenterConfiguration))
          as IdentityCenterConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IdentityCenterConfiguration create() =>
      IdentityCenterConfiguration._();
  @$core.override
  IdentityCenterConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IdentityCenterConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IdentityCenterConfiguration>(create);
  static IdentityCenterConfiguration? _defaultInstance;

  @$pb.TagNumber(418515238)
  $core.bool get enableidentitycenter => $_getBF(0);
  @$pb.TagNumber(418515238)
  set enableidentitycenter($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(418515238)
  $core.bool hasEnableidentitycenter() => $_has(0);
  @$pb.TagNumber(418515238)
  void clearEnableidentitycenter() => $_clearField(418515238);

  @$pb.TagNumber(469575873)
  $core.String get identitycenterinstancearn => $_getSZ(1);
  @$pb.TagNumber(469575873)
  set identitycenterinstancearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(469575873)
  $core.bool hasIdentitycenterinstancearn() => $_has(1);
  @$pb.TagNumber(469575873)
  void clearIdentitycenterinstancearn() => $_clearField(469575873);
}

class ImportNotebookInput extends $pb.GeneratedMessage {
  factory ImportNotebookInput({
    $core.String? payload,
    $core.String? name,
    NotebookType? type,
    $core.String? clientrequesttoken,
    $core.String? notebooks3locationuri,
    $core.String? workgroup,
  }) {
    final result = create();
    if (payload != null) result.payload = payload;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (notebooks3locationuri != null)
      result.notebooks3locationuri = notebooks3locationuri;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  ImportNotebookInput._();

  factory ImportNotebookInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportNotebookInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportNotebookInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(6526790, _omitFieldNames ? '' : 'payload')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<NotebookType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: NotebookType.values)
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..aOS(460333982, _omitFieldNames ? '' : 'notebooks3locationuri')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportNotebookInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportNotebookInput copyWith(void Function(ImportNotebookInput) updates) =>
      super.copyWith((message) => updates(message as ImportNotebookInput))
          as ImportNotebookInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportNotebookInput create() => ImportNotebookInput._();
  @$core.override
  ImportNotebookInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportNotebookInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportNotebookInput>(create);
  static ImportNotebookInput? _defaultInstance;

  @$pb.TagNumber(6526790)
  $core.String get payload => $_getSZ(0);
  @$pb.TagNumber(6526790)
  set payload($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6526790)
  $core.bool hasPayload() => $_has(0);
  @$pb.TagNumber(6526790)
  void clearPayload() => $_clearField(6526790);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  NotebookType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(NotebookType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(3);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(3);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(460333982)
  $core.String get notebooks3locationuri => $_getSZ(4);
  @$pb.TagNumber(460333982)
  set notebooks3locationuri($core.String value) => $_setString(4, value);
  @$pb.TagNumber(460333982)
  $core.bool hasNotebooks3locationuri() => $_has(4);
  @$pb.TagNumber(460333982)
  void clearNotebooks3locationuri() => $_clearField(460333982);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(5);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(5, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(5);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class ImportNotebookOutput extends $pb.GeneratedMessage {
  factory ImportNotebookOutput({
    $core.String? notebookid,
  }) {
    final result = create();
    if (notebookid != null) result.notebookid = notebookid;
    return result;
  }

  ImportNotebookOutput._();

  factory ImportNotebookOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportNotebookOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportNotebookOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportNotebookOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportNotebookOutput copyWith(void Function(ImportNotebookOutput) updates) =>
      super.copyWith((message) => updates(message as ImportNotebookOutput))
          as ImportNotebookOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportNotebookOutput create() => ImportNotebookOutput._();
  @$core.override
  ImportNotebookOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportNotebookOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportNotebookOutput>(create);
  static ImportNotebookOutput? _defaultInstance;

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(0);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(0);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
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

class InvalidRequestException extends $pb.GeneratedMessage {
  factory InvalidRequestException({
    $core.String? message,
    $core.String? athenaerrorcode,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (athenaerrorcode != null) result.athenaerrorcode = athenaerrorcode;
    return result;
  }

  InvalidRequestException._();

  factory InvalidRequestException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidRequestException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidRequestException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(335153222, _omitFieldNames ? '' : 'athenaerrorcode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidRequestException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidRequestException copyWith(
          void Function(InvalidRequestException) updates) =>
      super.copyWith((message) => updates(message as InvalidRequestException))
          as InvalidRequestException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidRequestException create() => InvalidRequestException._();
  @$core.override
  InvalidRequestException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidRequestException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidRequestException>(create);
  static InvalidRequestException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(335153222)
  $core.String get athenaerrorcode => $_getSZ(1);
  @$pb.TagNumber(335153222)
  set athenaerrorcode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(335153222)
  $core.bool hasAthenaerrorcode() => $_has(1);
  @$pb.TagNumber(335153222)
  void clearAthenaerrorcode() => $_clearField(335153222);
}

class ListApplicationDPUSizesInput extends $pb.GeneratedMessage {
  factory ListApplicationDPUSizesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListApplicationDPUSizesInput._();

  factory ListApplicationDPUSizesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListApplicationDPUSizesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListApplicationDPUSizesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApplicationDPUSizesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApplicationDPUSizesInput copyWith(
          void Function(ListApplicationDPUSizesInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListApplicationDPUSizesInput))
          as ListApplicationDPUSizesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListApplicationDPUSizesInput create() =>
      ListApplicationDPUSizesInput._();
  @$core.override
  ListApplicationDPUSizesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListApplicationDPUSizesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListApplicationDPUSizesInput>(create);
  static ListApplicationDPUSizesInput? _defaultInstance;

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

class ListApplicationDPUSizesOutput extends $pb.GeneratedMessage {
  factory ListApplicationDPUSizesOutput({
    $core.Iterable<ApplicationDPUSizes>? applicationdpusizes,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (applicationdpusizes != null)
      result.applicationdpusizes.addAll(applicationdpusizes);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListApplicationDPUSizesOutput._();

  factory ListApplicationDPUSizesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListApplicationDPUSizesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListApplicationDPUSizesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<ApplicationDPUSizes>(
        60251851, _omitFieldNames ? '' : 'applicationdpusizes',
        subBuilder: ApplicationDPUSizes.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApplicationDPUSizesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApplicationDPUSizesOutput copyWith(
          void Function(ListApplicationDPUSizesOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListApplicationDPUSizesOutput))
          as ListApplicationDPUSizesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListApplicationDPUSizesOutput create() =>
      ListApplicationDPUSizesOutput._();
  @$core.override
  ListApplicationDPUSizesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListApplicationDPUSizesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListApplicationDPUSizesOutput>(create);
  static ListApplicationDPUSizesOutput? _defaultInstance;

  @$pb.TagNumber(60251851)
  $pb.PbList<ApplicationDPUSizes> get applicationdpusizes => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListCalculationExecutionsRequest extends $pb.GeneratedMessage {
  factory ListCalculationExecutionsRequest({
    $core.String? sessionid,
    CalculationExecutionState? statefilter,
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    if (statefilter != null) result.statefilter = statefilter;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListCalculationExecutionsRequest._();

  factory ListCalculationExecutionsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCalculationExecutionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCalculationExecutionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aE<CalculationExecutionState>(
        184693297, _omitFieldNames ? '' : 'statefilter',
        enumValues: CalculationExecutionState.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCalculationExecutionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCalculationExecutionsRequest copyWith(
          void Function(ListCalculationExecutionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListCalculationExecutionsRequest))
          as ListCalculationExecutionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCalculationExecutionsRequest create() =>
      ListCalculationExecutionsRequest._();
  @$core.override
  ListCalculationExecutionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCalculationExecutionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCalculationExecutionsRequest>(
          create);
  static ListCalculationExecutionsRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(184693297)
  CalculationExecutionState get statefilter => $_getN(1);
  @$pb.TagNumber(184693297)
  set statefilter(CalculationExecutionState value) =>
      $_setField(184693297, value);
  @$pb.TagNumber(184693297)
  $core.bool hasStatefilter() => $_has(1);
  @$pb.TagNumber(184693297)
  void clearStatefilter() => $_clearField(184693297);

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
}

class ListCalculationExecutionsResponse extends $pb.GeneratedMessage {
  factory ListCalculationExecutionsResponse({
    $core.String? nexttoken,
    $core.Iterable<CalculationSummary>? calculations,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (calculations != null) result.calculations.addAll(calculations);
    return result;
  }

  ListCalculationExecutionsResponse._();

  factory ListCalculationExecutionsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCalculationExecutionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCalculationExecutionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<CalculationSummary>(335028832, _omitFieldNames ? '' : 'calculations',
        subBuilder: CalculationSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCalculationExecutionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCalculationExecutionsResponse copyWith(
          void Function(ListCalculationExecutionsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListCalculationExecutionsResponse))
          as ListCalculationExecutionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCalculationExecutionsResponse create() =>
      ListCalculationExecutionsResponse._();
  @$core.override
  ListCalculationExecutionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCalculationExecutionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCalculationExecutionsResponse>(
          create);
  static ListCalculationExecutionsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(335028832)
  $pb.PbList<CalculationSummary> get calculations => $_getList(1);
}

class ListCapacityReservationsInput extends $pb.GeneratedMessage {
  factory ListCapacityReservationsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListCapacityReservationsInput._();

  factory ListCapacityReservationsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCapacityReservationsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCapacityReservationsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCapacityReservationsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCapacityReservationsInput copyWith(
          void Function(ListCapacityReservationsInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListCapacityReservationsInput))
          as ListCapacityReservationsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCapacityReservationsInput create() =>
      ListCapacityReservationsInput._();
  @$core.override
  ListCapacityReservationsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCapacityReservationsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCapacityReservationsInput>(create);
  static ListCapacityReservationsInput? _defaultInstance;

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

class ListCapacityReservationsOutput extends $pb.GeneratedMessage {
  factory ListCapacityReservationsOutput({
    $core.String? nexttoken,
    $core.Iterable<CapacityReservation>? capacityreservations,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (capacityreservations != null)
      result.capacityreservations.addAll(capacityreservations);
    return result;
  }

  ListCapacityReservationsOutput._();

  factory ListCapacityReservationsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCapacityReservationsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCapacityReservationsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<CapacityReservation>(
        473497795, _omitFieldNames ? '' : 'capacityreservations',
        subBuilder: CapacityReservation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCapacityReservationsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCapacityReservationsOutput copyWith(
          void Function(ListCapacityReservationsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListCapacityReservationsOutput))
          as ListCapacityReservationsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCapacityReservationsOutput create() =>
      ListCapacityReservationsOutput._();
  @$core.override
  ListCapacityReservationsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCapacityReservationsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCapacityReservationsOutput>(create);
  static ListCapacityReservationsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(473497795)
  $pb.PbList<CapacityReservation> get capacityreservations => $_getList(1);
}

class ListDataCatalogsInput extends $pb.GeneratedMessage {
  factory ListDataCatalogsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  ListDataCatalogsInput._();

  factory ListDataCatalogsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDataCatalogsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDataCatalogsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDataCatalogsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDataCatalogsInput copyWith(
          void Function(ListDataCatalogsInput) updates) =>
      super.copyWith((message) => updates(message as ListDataCatalogsInput))
          as ListDataCatalogsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDataCatalogsInput create() => ListDataCatalogsInput._();
  @$core.override
  ListDataCatalogsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDataCatalogsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDataCatalogsInput>(create);
  static ListDataCatalogsInput? _defaultInstance;

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(2);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(2);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class ListDataCatalogsOutput extends $pb.GeneratedMessage {
  factory ListDataCatalogsOutput({
    $core.String? nexttoken,
    $core.Iterable<DataCatalogSummary>? datacatalogssummary,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (datacatalogssummary != null)
      result.datacatalogssummary.addAll(datacatalogssummary);
    return result;
  }

  ListDataCatalogsOutput._();

  factory ListDataCatalogsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDataCatalogsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDataCatalogsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<DataCatalogSummary>(
        499180880, _omitFieldNames ? '' : 'datacatalogssummary',
        subBuilder: DataCatalogSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDataCatalogsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDataCatalogsOutput copyWith(
          void Function(ListDataCatalogsOutput) updates) =>
      super.copyWith((message) => updates(message as ListDataCatalogsOutput))
          as ListDataCatalogsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDataCatalogsOutput create() => ListDataCatalogsOutput._();
  @$core.override
  ListDataCatalogsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDataCatalogsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDataCatalogsOutput>(create);
  static ListDataCatalogsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(499180880)
  $pb.PbList<DataCatalogSummary> get datacatalogssummary => $_getList(1);
}

class ListDatabasesInput extends $pb.GeneratedMessage {
  factory ListDatabasesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
    $core.String? catalogname,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    if (catalogname != null) result.catalogname = catalogname;
    return result;
  }

  ListDatabasesInput._();

  factory ListDatabasesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDatabasesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDatabasesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..aOS(518825212, _omitFieldNames ? '' : 'catalogname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesInput copyWith(void Function(ListDatabasesInput) updates) =>
      super.copyWith((message) => updates(message as ListDatabasesInput))
          as ListDatabasesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDatabasesInput create() => ListDatabasesInput._();
  @$core.override
  ListDatabasesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDatabasesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDatabasesInput>(create);
  static ListDatabasesInput? _defaultInstance;

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(2);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(2);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(518825212)
  $core.String get catalogname => $_getSZ(3);
  @$pb.TagNumber(518825212)
  set catalogname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(518825212)
  $core.bool hasCatalogname() => $_has(3);
  @$pb.TagNumber(518825212)
  void clearCatalogname() => $_clearField(518825212);
}

class ListDatabasesOutput extends $pb.GeneratedMessage {
  factory ListDatabasesOutput({
    $core.String? nexttoken,
    $core.Iterable<Database>? databaselist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (databaselist != null) result.databaselist.addAll(databaselist);
    return result;
  }

  ListDatabasesOutput._();

  factory ListDatabasesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDatabasesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDatabasesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Database>(393720897, _omitFieldNames ? '' : 'databaselist',
        subBuilder: Database.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesOutput copyWith(void Function(ListDatabasesOutput) updates) =>
      super.copyWith((message) => updates(message as ListDatabasesOutput))
          as ListDatabasesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDatabasesOutput create() => ListDatabasesOutput._();
  @$core.override
  ListDatabasesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDatabasesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDatabasesOutput>(create);
  static ListDatabasesOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(393720897)
  $pb.PbList<Database> get databaselist => $_getList(1);
}

class ListEngineVersionsInput extends $pb.GeneratedMessage {
  factory ListEngineVersionsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListEngineVersionsInput._();

  factory ListEngineVersionsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEngineVersionsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEngineVersionsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEngineVersionsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEngineVersionsInput copyWith(
          void Function(ListEngineVersionsInput) updates) =>
      super.copyWith((message) => updates(message as ListEngineVersionsInput))
          as ListEngineVersionsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEngineVersionsInput create() => ListEngineVersionsInput._();
  @$core.override
  ListEngineVersionsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEngineVersionsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEngineVersionsInput>(create);
  static ListEngineVersionsInput? _defaultInstance;

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

class ListEngineVersionsOutput extends $pb.GeneratedMessage {
  factory ListEngineVersionsOutput({
    $core.String? nexttoken,
    $core.Iterable<EngineVersion>? engineversions,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (engineversions != null) result.engineversions.addAll(engineversions);
    return result;
  }

  ListEngineVersionsOutput._();

  factory ListEngineVersionsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEngineVersionsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEngineVersionsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<EngineVersion>(500123251, _omitFieldNames ? '' : 'engineversions',
        subBuilder: EngineVersion.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEngineVersionsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEngineVersionsOutput copyWith(
          void Function(ListEngineVersionsOutput) updates) =>
      super.copyWith((message) => updates(message as ListEngineVersionsOutput))
          as ListEngineVersionsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEngineVersionsOutput create() => ListEngineVersionsOutput._();
  @$core.override
  ListEngineVersionsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEngineVersionsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEngineVersionsOutput>(create);
  static ListEngineVersionsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(500123251)
  $pb.PbList<EngineVersion> get engineversions => $_getList(1);
}

class ListExecutorsRequest extends $pb.GeneratedMessage {
  factory ListExecutorsRequest({
    $core.String? sessionid,
    ExecutorState? executorstatefilter,
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    if (executorstatefilter != null)
      result.executorstatefilter = executorstatefilter;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListExecutorsRequest._();

  factory ListExecutorsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListExecutorsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListExecutorsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aE<ExecutorState>(208448852, _omitFieldNames ? '' : 'executorstatefilter',
        enumValues: ExecutorState.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutorsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutorsRequest copyWith(void Function(ListExecutorsRequest) updates) =>
      super.copyWith((message) => updates(message as ListExecutorsRequest))
          as ListExecutorsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListExecutorsRequest create() => ListExecutorsRequest._();
  @$core.override
  ListExecutorsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListExecutorsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListExecutorsRequest>(create);
  static ListExecutorsRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(208448852)
  ExecutorState get executorstatefilter => $_getN(1);
  @$pb.TagNumber(208448852)
  set executorstatefilter(ExecutorState value) => $_setField(208448852, value);
  @$pb.TagNumber(208448852)
  $core.bool hasExecutorstatefilter() => $_has(1);
  @$pb.TagNumber(208448852)
  void clearExecutorstatefilter() => $_clearField(208448852);

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
}

class ListExecutorsResponse extends $pb.GeneratedMessage {
  factory ListExecutorsResponse({
    $core.String? sessionid,
    $core.String? nexttoken,
    $core.Iterable<ExecutorsSummary>? executorssummary,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (executorssummary != null)
      result.executorssummary.addAll(executorssummary);
    return result;
  }

  ListExecutorsResponse._();

  factory ListExecutorsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListExecutorsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListExecutorsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ExecutorsSummary>(
        450697516, _omitFieldNames ? '' : 'executorssummary',
        subBuilder: ExecutorsSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutorsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutorsResponse copyWith(
          void Function(ListExecutorsResponse) updates) =>
      super.copyWith((message) => updates(message as ListExecutorsResponse))
          as ListExecutorsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListExecutorsResponse create() => ListExecutorsResponse._();
  @$core.override
  ListExecutorsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListExecutorsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListExecutorsResponse>(create);
  static ListExecutorsResponse? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(450697516)
  $pb.PbList<ExecutorsSummary> get executorssummary => $_getList(2);
}

class ListNamedQueriesInput extends $pb.GeneratedMessage {
  factory ListNamedQueriesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  ListNamedQueriesInput._();

  factory ListNamedQueriesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListNamedQueriesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListNamedQueriesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNamedQueriesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNamedQueriesInput copyWith(
          void Function(ListNamedQueriesInput) updates) =>
      super.copyWith((message) => updates(message as ListNamedQueriesInput))
          as ListNamedQueriesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListNamedQueriesInput create() => ListNamedQueriesInput._();
  @$core.override
  ListNamedQueriesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListNamedQueriesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListNamedQueriesInput>(create);
  static ListNamedQueriesInput? _defaultInstance;

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(2);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(2);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class ListNamedQueriesOutput extends $pb.GeneratedMessage {
  factory ListNamedQueriesOutput({
    $core.Iterable<$core.String>? namedqueryids,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (namedqueryids != null) result.namedqueryids.addAll(namedqueryids);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListNamedQueriesOutput._();

  factory ListNamedQueriesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListNamedQueriesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListNamedQueriesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPS(6092797, _omitFieldNames ? '' : 'namedqueryids')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNamedQueriesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNamedQueriesOutput copyWith(
          void Function(ListNamedQueriesOutput) updates) =>
      super.copyWith((message) => updates(message as ListNamedQueriesOutput))
          as ListNamedQueriesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListNamedQueriesOutput create() => ListNamedQueriesOutput._();
  @$core.override
  ListNamedQueriesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListNamedQueriesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListNamedQueriesOutput>(create);
  static ListNamedQueriesOutput? _defaultInstance;

  @$pb.TagNumber(6092797)
  $pb.PbList<$core.String> get namedqueryids => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListNotebookMetadataInput extends $pb.GeneratedMessage {
  factory ListNotebookMetadataInput({
    FilterDefinition? filters,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
  }) {
    final result = create();
    if (filters != null) result.filters = filters;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  ListNotebookMetadataInput._();

  factory ListNotebookMetadataInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListNotebookMetadataInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListNotebookMetadataInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<FilterDefinition>(188393197, _omitFieldNames ? '' : 'filters',
        subBuilder: FilterDefinition.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookMetadataInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookMetadataInput copyWith(
          void Function(ListNotebookMetadataInput) updates) =>
      super.copyWith((message) => updates(message as ListNotebookMetadataInput))
          as ListNotebookMetadataInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListNotebookMetadataInput create() => ListNotebookMetadataInput._();
  @$core.override
  ListNotebookMetadataInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListNotebookMetadataInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListNotebookMetadataInput>(create);
  static ListNotebookMetadataInput? _defaultInstance;

  @$pb.TagNumber(188393197)
  FilterDefinition get filters => $_getN(0);
  @$pb.TagNumber(188393197)
  set filters(FilterDefinition value) => $_setField(188393197, value);
  @$pb.TagNumber(188393197)
  $core.bool hasFilters() => $_has(0);
  @$pb.TagNumber(188393197)
  void clearFilters() => $_clearField(188393197);
  @$pb.TagNumber(188393197)
  FilterDefinition ensureFilters() => $_ensure(0);

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(3);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(3, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(3);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class ListNotebookMetadataOutput extends $pb.GeneratedMessage {
  factory ListNotebookMetadataOutput({
    $core.String? nexttoken,
    $core.Iterable<NotebookMetadata>? notebookmetadatalist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (notebookmetadatalist != null)
      result.notebookmetadatalist.addAll(notebookmetadatalist);
    return result;
  }

  ListNotebookMetadataOutput._();

  factory ListNotebookMetadataOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListNotebookMetadataOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListNotebookMetadataOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<NotebookMetadata>(
        319238242, _omitFieldNames ? '' : 'notebookmetadatalist',
        subBuilder: NotebookMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookMetadataOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookMetadataOutput copyWith(
          void Function(ListNotebookMetadataOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListNotebookMetadataOutput))
          as ListNotebookMetadataOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListNotebookMetadataOutput create() => ListNotebookMetadataOutput._();
  @$core.override
  ListNotebookMetadataOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListNotebookMetadataOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListNotebookMetadataOutput>(create);
  static ListNotebookMetadataOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(319238242)
  $pb.PbList<NotebookMetadata> get notebookmetadatalist => $_getList(1);
}

class ListNotebookSessionsRequest extends $pb.GeneratedMessage {
  factory ListNotebookSessionsRequest({
    $core.String? notebookid,
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (notebookid != null) result.notebookid = notebookid;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListNotebookSessionsRequest._();

  factory ListNotebookSessionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListNotebookSessionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListNotebookSessionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookSessionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookSessionsRequest copyWith(
          void Function(ListNotebookSessionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListNotebookSessionsRequest))
          as ListNotebookSessionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListNotebookSessionsRequest create() =>
      ListNotebookSessionsRequest._();
  @$core.override
  ListNotebookSessionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListNotebookSessionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListNotebookSessionsRequest>(create);
  static ListNotebookSessionsRequest? _defaultInstance;

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(0);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(0);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);

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
}

class ListNotebookSessionsResponse extends $pb.GeneratedMessage {
  factory ListNotebookSessionsResponse({
    $core.String? nexttoken,
    $core.Iterable<NotebookSessionSummary>? notebooksessionslist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (notebooksessionslist != null)
      result.notebooksessionslist.addAll(notebooksessionslist);
    return result;
  }

  ListNotebookSessionsResponse._();

  factory ListNotebookSessionsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListNotebookSessionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListNotebookSessionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<NotebookSessionSummary>(
        247821802, _omitFieldNames ? '' : 'notebooksessionslist',
        subBuilder: NotebookSessionSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookSessionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListNotebookSessionsResponse copyWith(
          void Function(ListNotebookSessionsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListNotebookSessionsResponse))
          as ListNotebookSessionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListNotebookSessionsResponse create() =>
      ListNotebookSessionsResponse._();
  @$core.override
  ListNotebookSessionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListNotebookSessionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListNotebookSessionsResponse>(create);
  static ListNotebookSessionsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(247821802)
  $pb.PbList<NotebookSessionSummary> get notebooksessionslist => $_getList(1);
}

class ListPreparedStatementsInput extends $pb.GeneratedMessage {
  factory ListPreparedStatementsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  ListPreparedStatementsInput._();

  factory ListPreparedStatementsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPreparedStatementsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPreparedStatementsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPreparedStatementsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPreparedStatementsInput copyWith(
          void Function(ListPreparedStatementsInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListPreparedStatementsInput))
          as ListPreparedStatementsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPreparedStatementsInput create() =>
      ListPreparedStatementsInput._();
  @$core.override
  ListPreparedStatementsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPreparedStatementsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPreparedStatementsInput>(create);
  static ListPreparedStatementsInput? _defaultInstance;

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(2);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(2);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class ListPreparedStatementsOutput extends $pb.GeneratedMessage {
  factory ListPreparedStatementsOutput({
    $core.String? nexttoken,
    $core.Iterable<PreparedStatementSummary>? preparedstatements,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (preparedstatements != null)
      result.preparedstatements.addAll(preparedstatements);
    return result;
  }

  ListPreparedStatementsOutput._();

  factory ListPreparedStatementsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPreparedStatementsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPreparedStatementsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<PreparedStatementSummary>(
        526923667, _omitFieldNames ? '' : 'preparedstatements',
        subBuilder: PreparedStatementSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPreparedStatementsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPreparedStatementsOutput copyWith(
          void Function(ListPreparedStatementsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListPreparedStatementsOutput))
          as ListPreparedStatementsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPreparedStatementsOutput create() =>
      ListPreparedStatementsOutput._();
  @$core.override
  ListPreparedStatementsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPreparedStatementsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPreparedStatementsOutput>(create);
  static ListPreparedStatementsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(526923667)
  $pb.PbList<PreparedStatementSummary> get preparedstatements => $_getList(1);
}

class ListQueryExecutionsInput extends $pb.GeneratedMessage {
  factory ListQueryExecutionsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  ListQueryExecutionsInput._();

  factory ListQueryExecutionsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueryExecutionsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueryExecutionsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryExecutionsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryExecutionsInput copyWith(
          void Function(ListQueryExecutionsInput) updates) =>
      super.copyWith((message) => updates(message as ListQueryExecutionsInput))
          as ListQueryExecutionsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueryExecutionsInput create() => ListQueryExecutionsInput._();
  @$core.override
  ListQueryExecutionsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueryExecutionsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueryExecutionsInput>(create);
  static ListQueryExecutionsInput? _defaultInstance;

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(2);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(2);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class ListQueryExecutionsOutput extends $pb.GeneratedMessage {
  factory ListQueryExecutionsOutput({
    $core.String? nexttoken,
    $core.Iterable<$core.String>? queryexecutionids,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (queryexecutionids != null)
      result.queryexecutionids.addAll(queryexecutionids);
    return result;
  }

  ListQueryExecutionsOutput._();

  factory ListQueryExecutionsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueryExecutionsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueryExecutionsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPS(493941192, _omitFieldNames ? '' : 'queryexecutionids')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryExecutionsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryExecutionsOutput copyWith(
          void Function(ListQueryExecutionsOutput) updates) =>
      super.copyWith((message) => updates(message as ListQueryExecutionsOutput))
          as ListQueryExecutionsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueryExecutionsOutput create() => ListQueryExecutionsOutput._();
  @$core.override
  ListQueryExecutionsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueryExecutionsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueryExecutionsOutput>(create);
  static ListQueryExecutionsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(493941192)
  $pb.PbList<$core.String> get queryexecutionids => $_getList(1);
}

class ListSessionsRequest extends $pb.GeneratedMessage {
  factory ListSessionsRequest({
    SessionState? statefilter,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
  }) {
    final result = create();
    if (statefilter != null) result.statefilter = statefilter;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  ListSessionsRequest._();

  factory ListSessionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSessionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSessionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<SessionState>(184693297, _omitFieldNames ? '' : 'statefilter',
        enumValues: SessionState.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSessionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSessionsRequest copyWith(void Function(ListSessionsRequest) updates) =>
      super.copyWith((message) => updates(message as ListSessionsRequest))
          as ListSessionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSessionsRequest create() => ListSessionsRequest._();
  @$core.override
  ListSessionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSessionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSessionsRequest>(create);
  static ListSessionsRequest? _defaultInstance;

  @$pb.TagNumber(184693297)
  SessionState get statefilter => $_getN(0);
  @$pb.TagNumber(184693297)
  set statefilter(SessionState value) => $_setField(184693297, value);
  @$pb.TagNumber(184693297)
  $core.bool hasStatefilter() => $_has(0);
  @$pb.TagNumber(184693297)
  void clearStatefilter() => $_clearField(184693297);

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(3);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(3, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(3);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class ListSessionsResponse extends $pb.GeneratedMessage {
  factory ListSessionsResponse({
    $core.String? nexttoken,
    $core.Iterable<SessionSummary>? sessions,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (sessions != null) result.sessions.addAll(sessions);
    return result;
  }

  ListSessionsResponse._();

  factory ListSessionsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSessionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSessionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<SessionSummary>(379226861, _omitFieldNames ? '' : 'sessions',
        subBuilder: SessionSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSessionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSessionsResponse copyWith(void Function(ListSessionsResponse) updates) =>
      super.copyWith((message) => updates(message as ListSessionsResponse))
          as ListSessionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSessionsResponse create() => ListSessionsResponse._();
  @$core.override
  ListSessionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSessionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSessionsResponse>(create);
  static ListSessionsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(379226861)
  $pb.PbList<SessionSummary> get sessions => $_getList(1);
}

class ListTableMetadataInput extends $pb.GeneratedMessage {
  factory ListTableMetadataInput({
    $core.String? databasename,
    $core.String? expression,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? workgroup,
    $core.String? catalogname,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (expression != null) result.expression = expression;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (workgroup != null) result.workgroup = workgroup;
    if (catalogname != null) result.catalogname = catalogname;
    return result;
  }

  ListTableMetadataInput._();

  factory ListTableMetadataInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTableMetadataInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTableMetadataInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(193051916, _omitFieldNames ? '' : 'expression')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..aOS(518825212, _omitFieldNames ? '' : 'catalogname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTableMetadataInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTableMetadataInput copyWith(
          void Function(ListTableMetadataInput) updates) =>
      super.copyWith((message) => updates(message as ListTableMetadataInput))
          as ListTableMetadataInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTableMetadataInput create() => ListTableMetadataInput._();
  @$core.override
  ListTableMetadataInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTableMetadataInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTableMetadataInput>(create);
  static ListTableMetadataInput? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(193051916)
  $core.String get expression => $_getSZ(1);
  @$pb.TagNumber(193051916)
  set expression($core.String value) => $_setString(1, value);
  @$pb.TagNumber(193051916)
  $core.bool hasExpression() => $_has(1);
  @$pb.TagNumber(193051916)
  void clearExpression() => $_clearField(193051916);

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

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(4);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(4, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(4);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(518825212)
  $core.String get catalogname => $_getSZ(5);
  @$pb.TagNumber(518825212)
  set catalogname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(518825212)
  $core.bool hasCatalogname() => $_has(5);
  @$pb.TagNumber(518825212)
  void clearCatalogname() => $_clearField(518825212);
}

class ListTableMetadataOutput extends $pb.GeneratedMessage {
  factory ListTableMetadataOutput({
    $core.String? nexttoken,
    $core.Iterable<TableMetadata>? tablemetadatalist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (tablemetadatalist != null)
      result.tablemetadatalist.addAll(tablemetadatalist);
    return result;
  }

  ListTableMetadataOutput._();

  factory ListTableMetadataOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTableMetadataOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTableMetadataOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<TableMetadata>(532092415, _omitFieldNames ? '' : 'tablemetadatalist',
        subBuilder: TableMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTableMetadataOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTableMetadataOutput copyWith(
          void Function(ListTableMetadataOutput) updates) =>
      super.copyWith((message) => updates(message as ListTableMetadataOutput))
          as ListTableMetadataOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTableMetadataOutput create() => ListTableMetadataOutput._();
  @$core.override
  ListTableMetadataOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTableMetadataOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTableMetadataOutput>(create);
  static ListTableMetadataOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(532092415)
  $pb.PbList<TableMetadata> get tablemetadatalist => $_getList(1);
}

class ListTagsForResourceInput extends $pb.GeneratedMessage {
  factory ListTagsForResourceInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(2);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(2);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class ListTagsForResourceOutput extends $pb.GeneratedMessage {
  factory ListTagsForResourceOutput({
    $core.String? nexttoken,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
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

class ListWorkGroupsInput extends $pb.GeneratedMessage {
  factory ListWorkGroupsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListWorkGroupsInput._();

  factory ListWorkGroupsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListWorkGroupsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListWorkGroupsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWorkGroupsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWorkGroupsInput copyWith(void Function(ListWorkGroupsInput) updates) =>
      super.copyWith((message) => updates(message as ListWorkGroupsInput))
          as ListWorkGroupsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWorkGroupsInput create() => ListWorkGroupsInput._();
  @$core.override
  ListWorkGroupsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListWorkGroupsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListWorkGroupsInput>(create);
  static ListWorkGroupsInput? _defaultInstance;

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

class ListWorkGroupsOutput extends $pb.GeneratedMessage {
  factory ListWorkGroupsOutput({
    $core.Iterable<WorkGroupSummary>? workgroups,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (workgroups != null) result.workgroups.addAll(workgroups);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListWorkGroupsOutput._();

  factory ListWorkGroupsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListWorkGroupsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListWorkGroupsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<WorkGroupSummary>(159439825, _omitFieldNames ? '' : 'workgroups',
        subBuilder: WorkGroupSummary.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWorkGroupsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWorkGroupsOutput copyWith(void Function(ListWorkGroupsOutput) updates) =>
      super.copyWith((message) => updates(message as ListWorkGroupsOutput))
          as ListWorkGroupsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWorkGroupsOutput create() => ListWorkGroupsOutput._();
  @$core.override
  ListWorkGroupsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListWorkGroupsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListWorkGroupsOutput>(create);
  static ListWorkGroupsOutput? _defaultInstance;

  @$pb.TagNumber(159439825)
  $pb.PbList<WorkGroupSummary> get workgroups => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ManagedLoggingConfiguration extends $pb.GeneratedMessage {
  factory ManagedLoggingConfiguration({
    $core.String? kmskey,
    $core.bool? enabled,
  }) {
    final result = create();
    if (kmskey != null) result.kmskey = kmskey;
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  ManagedLoggingConfiguration._();

  factory ManagedLoggingConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedLoggingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedLoggingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(114561194, _omitFieldNames ? '' : 'kmskey')
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedLoggingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedLoggingConfiguration copyWith(
          void Function(ManagedLoggingConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as ManagedLoggingConfiguration))
          as ManagedLoggingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedLoggingConfiguration create() =>
      ManagedLoggingConfiguration._();
  @$core.override
  ManagedLoggingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedLoggingConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedLoggingConfiguration>(create);
  static ManagedLoggingConfiguration? _defaultInstance;

  @$pb.TagNumber(114561194)
  $core.String get kmskey => $_getSZ(0);
  @$pb.TagNumber(114561194)
  set kmskey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(114561194)
  $core.bool hasKmskey() => $_has(0);
  @$pb.TagNumber(114561194)
  void clearKmskey() => $_clearField(114561194);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);
}

class ManagedQueryResultsConfiguration extends $pb.GeneratedMessage {
  factory ManagedQueryResultsConfiguration({
    ManagedQueryResultsEncryptionConfiguration? encryptionconfiguration,
    $core.bool? enabled,
  }) {
    final result = create();
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  ManagedQueryResultsConfiguration._();

  factory ManagedQueryResultsConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedQueryResultsConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedQueryResultsConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<ManagedQueryResultsEncryptionConfiguration>(
        225764215, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: ManagedQueryResultsEncryptionConfiguration.create)
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedQueryResultsConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedQueryResultsConfiguration copyWith(
          void Function(ManagedQueryResultsConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as ManagedQueryResultsConfiguration))
          as ManagedQueryResultsConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedQueryResultsConfiguration create() =>
      ManagedQueryResultsConfiguration._();
  @$core.override
  ManagedQueryResultsConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedQueryResultsConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedQueryResultsConfiguration>(
          create);
  static ManagedQueryResultsConfiguration? _defaultInstance;

  @$pb.TagNumber(225764215)
  ManagedQueryResultsEncryptionConfiguration get encryptionconfiguration =>
      $_getN(0);
  @$pb.TagNumber(225764215)
  set encryptionconfiguration(
          ManagedQueryResultsEncryptionConfiguration value) =>
      $_setField(225764215, value);
  @$pb.TagNumber(225764215)
  $core.bool hasEncryptionconfiguration() => $_has(0);
  @$pb.TagNumber(225764215)
  void clearEncryptionconfiguration() => $_clearField(225764215);
  @$pb.TagNumber(225764215)
  ManagedQueryResultsEncryptionConfiguration ensureEncryptionconfiguration() =>
      $_ensure(0);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);
}

class ManagedQueryResultsConfigurationUpdates extends $pb.GeneratedMessage {
  factory ManagedQueryResultsConfigurationUpdates({
    ManagedQueryResultsEncryptionConfiguration? encryptionconfiguration,
    $core.bool? removeencryptionconfiguration,
    $core.bool? enabled,
  }) {
    final result = create();
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (removeencryptionconfiguration != null)
      result.removeencryptionconfiguration = removeencryptionconfiguration;
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  ManagedQueryResultsConfigurationUpdates._();

  factory ManagedQueryResultsConfigurationUpdates.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedQueryResultsConfigurationUpdates.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedQueryResultsConfigurationUpdates',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<ManagedQueryResultsEncryptionConfiguration>(
        225764215, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: ManagedQueryResultsEncryptionConfiguration.create)
    ..aOB(294004335, _omitFieldNames ? '' : 'removeencryptionconfiguration')
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedQueryResultsConfigurationUpdates clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedQueryResultsConfigurationUpdates copyWith(
          void Function(ManagedQueryResultsConfigurationUpdates) updates) =>
      super.copyWith((message) =>
              updates(message as ManagedQueryResultsConfigurationUpdates))
          as ManagedQueryResultsConfigurationUpdates;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedQueryResultsConfigurationUpdates create() =>
      ManagedQueryResultsConfigurationUpdates._();
  @$core.override
  ManagedQueryResultsConfigurationUpdates createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedQueryResultsConfigurationUpdates getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ManagedQueryResultsConfigurationUpdates>(create);
  static ManagedQueryResultsConfigurationUpdates? _defaultInstance;

  @$pb.TagNumber(225764215)
  ManagedQueryResultsEncryptionConfiguration get encryptionconfiguration =>
      $_getN(0);
  @$pb.TagNumber(225764215)
  set encryptionconfiguration(
          ManagedQueryResultsEncryptionConfiguration value) =>
      $_setField(225764215, value);
  @$pb.TagNumber(225764215)
  $core.bool hasEncryptionconfiguration() => $_has(0);
  @$pb.TagNumber(225764215)
  void clearEncryptionconfiguration() => $_clearField(225764215);
  @$pb.TagNumber(225764215)
  ManagedQueryResultsEncryptionConfiguration ensureEncryptionconfiguration() =>
      $_ensure(0);

  @$pb.TagNumber(294004335)
  $core.bool get removeencryptionconfiguration => $_getBF(1);
  @$pb.TagNumber(294004335)
  set removeencryptionconfiguration($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(294004335)
  $core.bool hasRemoveencryptionconfiguration() => $_has(1);
  @$pb.TagNumber(294004335)
  void clearRemoveencryptionconfiguration() => $_clearField(294004335);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(2);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(2);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);
}

class ManagedQueryResultsEncryptionConfiguration extends $pb.GeneratedMessage {
  factory ManagedQueryResultsEncryptionConfiguration({
    $core.String? kmskey,
  }) {
    final result = create();
    if (kmskey != null) result.kmskey = kmskey;
    return result;
  }

  ManagedQueryResultsEncryptionConfiguration._();

  factory ManagedQueryResultsEncryptionConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedQueryResultsEncryptionConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedQueryResultsEncryptionConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(114561194, _omitFieldNames ? '' : 'kmskey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedQueryResultsEncryptionConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedQueryResultsEncryptionConfiguration copyWith(
          void Function(ManagedQueryResultsEncryptionConfiguration) updates) =>
      super.copyWith((message) =>
              updates(message as ManagedQueryResultsEncryptionConfiguration))
          as ManagedQueryResultsEncryptionConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedQueryResultsEncryptionConfiguration create() =>
      ManagedQueryResultsEncryptionConfiguration._();
  @$core.override
  ManagedQueryResultsEncryptionConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedQueryResultsEncryptionConfiguration getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ManagedQueryResultsEncryptionConfiguration>(create);
  static ManagedQueryResultsEncryptionConfiguration? _defaultInstance;

  @$pb.TagNumber(114561194)
  $core.String get kmskey => $_getSZ(0);
  @$pb.TagNumber(114561194)
  set kmskey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(114561194)
  $core.bool hasKmskey() => $_has(0);
  @$pb.TagNumber(114561194)
  void clearKmskey() => $_clearField(114561194);
}

class MetadataException extends $pb.GeneratedMessage {
  factory MetadataException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  MetadataException._();

  factory MetadataException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetadataException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetadataException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetadataException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetadataException copyWith(void Function(MetadataException) updates) =>
      super.copyWith((message) => updates(message as MetadataException))
          as MetadataException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetadataException create() => MetadataException._();
  @$core.override
  MetadataException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetadataException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetadataException>(create);
  static MetadataException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class MonitoringConfiguration extends $pb.GeneratedMessage {
  factory MonitoringConfiguration({
    S3LoggingConfiguration? s3loggingconfiguration,
    CloudWatchLoggingConfiguration? cloudwatchloggingconfiguration,
    ManagedLoggingConfiguration? managedloggingconfiguration,
  }) {
    final result = create();
    if (s3loggingconfiguration != null)
      result.s3loggingconfiguration = s3loggingconfiguration;
    if (cloudwatchloggingconfiguration != null)
      result.cloudwatchloggingconfiguration = cloudwatchloggingconfiguration;
    if (managedloggingconfiguration != null)
      result.managedloggingconfiguration = managedloggingconfiguration;
    return result;
  }

  MonitoringConfiguration._();

  factory MonitoringConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MonitoringConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MonitoringConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<S3LoggingConfiguration>(
        15691207, _omitFieldNames ? '' : 's3loggingconfiguration',
        subBuilder: S3LoggingConfiguration.create)
    ..aOM<CloudWatchLoggingConfiguration>(
        60618733, _omitFieldNames ? '' : 'cloudwatchloggingconfiguration',
        subBuilder: CloudWatchLoggingConfiguration.create)
    ..aOM<ManagedLoggingConfiguration>(
        270845728, _omitFieldNames ? '' : 'managedloggingconfiguration',
        subBuilder: ManagedLoggingConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MonitoringConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MonitoringConfiguration copyWith(
          void Function(MonitoringConfiguration) updates) =>
      super.copyWith((message) => updates(message as MonitoringConfiguration))
          as MonitoringConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MonitoringConfiguration create() => MonitoringConfiguration._();
  @$core.override
  MonitoringConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MonitoringConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MonitoringConfiguration>(create);
  static MonitoringConfiguration? _defaultInstance;

  @$pb.TagNumber(15691207)
  S3LoggingConfiguration get s3loggingconfiguration => $_getN(0);
  @$pb.TagNumber(15691207)
  set s3loggingconfiguration(S3LoggingConfiguration value) =>
      $_setField(15691207, value);
  @$pb.TagNumber(15691207)
  $core.bool hasS3loggingconfiguration() => $_has(0);
  @$pb.TagNumber(15691207)
  void clearS3loggingconfiguration() => $_clearField(15691207);
  @$pb.TagNumber(15691207)
  S3LoggingConfiguration ensureS3loggingconfiguration() => $_ensure(0);

  @$pb.TagNumber(60618733)
  CloudWatchLoggingConfiguration get cloudwatchloggingconfiguration =>
      $_getN(1);
  @$pb.TagNumber(60618733)
  set cloudwatchloggingconfiguration(CloudWatchLoggingConfiguration value) =>
      $_setField(60618733, value);
  @$pb.TagNumber(60618733)
  $core.bool hasCloudwatchloggingconfiguration() => $_has(1);
  @$pb.TagNumber(60618733)
  void clearCloudwatchloggingconfiguration() => $_clearField(60618733);
  @$pb.TagNumber(60618733)
  CloudWatchLoggingConfiguration ensureCloudwatchloggingconfiguration() =>
      $_ensure(1);

  @$pb.TagNumber(270845728)
  ManagedLoggingConfiguration get managedloggingconfiguration => $_getN(2);
  @$pb.TagNumber(270845728)
  set managedloggingconfiguration(ManagedLoggingConfiguration value) =>
      $_setField(270845728, value);
  @$pb.TagNumber(270845728)
  $core.bool hasManagedloggingconfiguration() => $_has(2);
  @$pb.TagNumber(270845728)
  void clearManagedloggingconfiguration() => $_clearField(270845728);
  @$pb.TagNumber(270845728)
  ManagedLoggingConfiguration ensureManagedloggingconfiguration() =>
      $_ensure(2);
}

class NamedQuery extends $pb.GeneratedMessage {
  factory NamedQuery({
    $core.String? description,
    $core.String? name,
    $core.String? database,
    $core.String? namedqueryid,
    $core.String? querystring,
    $core.String? workgroup,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (database != null) result.database = database;
    if (namedqueryid != null) result.namedqueryid = namedqueryid;
    if (querystring != null) result.querystring = querystring;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  NamedQuery._();

  factory NamedQuery.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NamedQuery.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NamedQuery',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(278147289, _omitFieldNames ? '' : 'database')
    ..aOS(330896872, _omitFieldNames ? '' : 'namedqueryid')
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NamedQuery clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NamedQuery copyWith(void Function(NamedQuery) updates) =>
      super.copyWith((message) => updates(message as NamedQuery)) as NamedQuery;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NamedQuery create() => NamedQuery._();
  @$core.override
  NamedQuery createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NamedQuery getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NamedQuery>(create);
  static NamedQuery? _defaultInstance;

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

  @$pb.TagNumber(278147289)
  $core.String get database => $_getSZ(2);
  @$pb.TagNumber(278147289)
  set database($core.String value) => $_setString(2, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(2);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);

  @$pb.TagNumber(330896872)
  $core.String get namedqueryid => $_getSZ(3);
  @$pb.TagNumber(330896872)
  set namedqueryid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(330896872)
  $core.bool hasNamedqueryid() => $_has(3);
  @$pb.TagNumber(330896872)
  void clearNamedqueryid() => $_clearField(330896872);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(4);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(4, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(4);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(5);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(5, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(5);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class NotebookMetadata extends $pb.GeneratedMessage {
  factory NotebookMetadata({
    $core.String? creationtime,
    $core.String? notebookid,
    $core.String? lastmodifiedtime,
    $core.String? name,
    NotebookType? type,
    $core.String? workgroup,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (notebookid != null) result.notebookid = notebookid;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  NotebookMetadata._();

  factory NotebookMetadata.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotebookMetadata.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotebookMetadata',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<NotebookType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: NotebookType.values)
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotebookMetadata clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotebookMetadata copyWith(void Function(NotebookMetadata) updates) =>
      super.copyWith((message) => updates(message as NotebookMetadata))
          as NotebookMetadata;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotebookMetadata create() => NotebookMetadata._();
  @$core.override
  NotebookMetadata createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotebookMetadata getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NotebookMetadata>(create);
  static NotebookMetadata? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(1);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(1);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(2);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(2);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  NotebookType get type => $_getN(4);
  @$pb.TagNumber(290836590)
  set type(NotebookType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(4);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(5);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(5, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(5);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class NotebookSessionSummary extends $pb.GeneratedMessage {
  factory NotebookSessionSummary({
    $core.String? sessionid,
    $core.String? creationtime,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    if (creationtime != null) result.creationtime = creationtime;
    return result;
  }

  NotebookSessionSummary._();

  factory NotebookSessionSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotebookSessionSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotebookSessionSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotebookSessionSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotebookSessionSummary copyWith(
          void Function(NotebookSessionSummary) updates) =>
      super.copyWith((message) => updates(message as NotebookSessionSummary))
          as NotebookSessionSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotebookSessionSummary create() => NotebookSessionSummary._();
  @$core.override
  NotebookSessionSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotebookSessionSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NotebookSessionSummary>(create);
  static NotebookSessionSummary? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);
}

class PreparedStatement extends $pb.GeneratedMessage {
  factory PreparedStatement({
    $core.String? statementname,
    $core.String? description,
    $core.String? lastmodifiedtime,
    $core.String? workgroupname,
    $core.String? querystatement,
  }) {
    final result = create();
    if (statementname != null) result.statementname = statementname;
    if (description != null) result.description = description;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (workgroupname != null) result.workgroupname = workgroupname;
    if (querystatement != null) result.querystatement = querystatement;
    return result;
  }

  PreparedStatement._();

  factory PreparedStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PreparedStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PreparedStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(258526185, _omitFieldNames ? '' : 'workgroupname')
    ..aOS(340852217, _omitFieldNames ? '' : 'querystatement')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PreparedStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PreparedStatement copyWith(void Function(PreparedStatement) updates) =>
      super.copyWith((message) => updates(message as PreparedStatement))
          as PreparedStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PreparedStatement create() => PreparedStatement._();
  @$core.override
  PreparedStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PreparedStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PreparedStatement>(create);
  static PreparedStatement? _defaultInstance;

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(0);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(0);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(2);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(2);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(258526185)
  $core.String get workgroupname => $_getSZ(3);
  @$pb.TagNumber(258526185)
  set workgroupname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(258526185)
  $core.bool hasWorkgroupname() => $_has(3);
  @$pb.TagNumber(258526185)
  void clearWorkgroupname() => $_clearField(258526185);

  @$pb.TagNumber(340852217)
  $core.String get querystatement => $_getSZ(4);
  @$pb.TagNumber(340852217)
  set querystatement($core.String value) => $_setString(4, value);
  @$pb.TagNumber(340852217)
  $core.bool hasQuerystatement() => $_has(4);
  @$pb.TagNumber(340852217)
  void clearQuerystatement() => $_clearField(340852217);
}

class PreparedStatementSummary extends $pb.GeneratedMessage {
  factory PreparedStatementSummary({
    $core.String? statementname,
    $core.String? lastmodifiedtime,
  }) {
    final result = create();
    if (statementname != null) result.statementname = statementname;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    return result;
  }

  PreparedStatementSummary._();

  factory PreparedStatementSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PreparedStatementSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PreparedStatementSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PreparedStatementSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PreparedStatementSummary copyWith(
          void Function(PreparedStatementSummary) updates) =>
      super.copyWith((message) => updates(message as PreparedStatementSummary))
          as PreparedStatementSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PreparedStatementSummary create() => PreparedStatementSummary._();
  @$core.override
  PreparedStatementSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PreparedStatementSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PreparedStatementSummary>(create);
  static PreparedStatementSummary? _defaultInstance;

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(0);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(0);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(1);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(1);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);
}

class PutCapacityAssignmentConfigurationInput extends $pb.GeneratedMessage {
  factory PutCapacityAssignmentConfigurationInput({
    $core.String? capacityreservationname,
    $core.Iterable<CapacityAssignment>? capacityassignments,
  }) {
    final result = create();
    if (capacityreservationname != null)
      result.capacityreservationname = capacityreservationname;
    if (capacityassignments != null)
      result.capacityassignments.addAll(capacityassignments);
    return result;
  }

  PutCapacityAssignmentConfigurationInput._();

  factory PutCapacityAssignmentConfigurationInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutCapacityAssignmentConfigurationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutCapacityAssignmentConfigurationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(327567687, _omitFieldNames ? '' : 'capacityreservationname')
    ..pPM<CapacityAssignment>(
        345772294, _omitFieldNames ? '' : 'capacityassignments',
        subBuilder: CapacityAssignment.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutCapacityAssignmentConfigurationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutCapacityAssignmentConfigurationInput copyWith(
          void Function(PutCapacityAssignmentConfigurationInput) updates) =>
      super.copyWith((message) =>
              updates(message as PutCapacityAssignmentConfigurationInput))
          as PutCapacityAssignmentConfigurationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutCapacityAssignmentConfigurationInput create() =>
      PutCapacityAssignmentConfigurationInput._();
  @$core.override
  PutCapacityAssignmentConfigurationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutCapacityAssignmentConfigurationInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          PutCapacityAssignmentConfigurationInput>(create);
  static PutCapacityAssignmentConfigurationInput? _defaultInstance;

  @$pb.TagNumber(327567687)
  $core.String get capacityreservationname => $_getSZ(0);
  @$pb.TagNumber(327567687)
  set capacityreservationname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(327567687)
  $core.bool hasCapacityreservationname() => $_has(0);
  @$pb.TagNumber(327567687)
  void clearCapacityreservationname() => $_clearField(327567687);

  @$pb.TagNumber(345772294)
  $pb.PbList<CapacityAssignment> get capacityassignments => $_getList(1);
}

class PutCapacityAssignmentConfigurationOutput extends $pb.GeneratedMessage {
  factory PutCapacityAssignmentConfigurationOutput() => create();

  PutCapacityAssignmentConfigurationOutput._();

  factory PutCapacityAssignmentConfigurationOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutCapacityAssignmentConfigurationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutCapacityAssignmentConfigurationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutCapacityAssignmentConfigurationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutCapacityAssignmentConfigurationOutput copyWith(
          void Function(PutCapacityAssignmentConfigurationOutput) updates) =>
      super.copyWith((message) =>
              updates(message as PutCapacityAssignmentConfigurationOutput))
          as PutCapacityAssignmentConfigurationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutCapacityAssignmentConfigurationOutput create() =>
      PutCapacityAssignmentConfigurationOutput._();
  @$core.override
  PutCapacityAssignmentConfigurationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutCapacityAssignmentConfigurationOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          PutCapacityAssignmentConfigurationOutput>(create);
  static PutCapacityAssignmentConfigurationOutput? _defaultInstance;
}

class QueryExecution extends $pb.GeneratedMessage {
  factory QueryExecution({
    QueryExecutionStatus? status,
    EngineVersion? engineversion,
    QueryExecutionContext? queryexecutioncontext,
    ManagedQueryResultsConfiguration? managedqueryresultsconfiguration,
    ResultConfiguration? resultconfiguration,
    QueryResultsS3AccessGrantsConfiguration?
        queryresultss3accessgrantsconfiguration,
    StatementType? statementtype,
    $core.String? queryexecutionid,
    ResultReuseConfiguration? resultreuseconfiguration,
    $core.String? substatementtype,
    $core.String? workgroup,
    QueryExecutionStatistics? statistics,
    $core.String? query,
    $core.Iterable<$core.String>? executionparameters,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (engineversion != null) result.engineversion = engineversion;
    if (queryexecutioncontext != null)
      result.queryexecutioncontext = queryexecutioncontext;
    if (managedqueryresultsconfiguration != null)
      result.managedqueryresultsconfiguration =
          managedqueryresultsconfiguration;
    if (resultconfiguration != null)
      result.resultconfiguration = resultconfiguration;
    if (queryresultss3accessgrantsconfiguration != null)
      result.queryresultss3accessgrantsconfiguration =
          queryresultss3accessgrantsconfiguration;
    if (statementtype != null) result.statementtype = statementtype;
    if (queryexecutionid != null) result.queryexecutionid = queryexecutionid;
    if (resultreuseconfiguration != null)
      result.resultreuseconfiguration = resultreuseconfiguration;
    if (substatementtype != null) result.substatementtype = substatementtype;
    if (workgroup != null) result.workgroup = workgroup;
    if (statistics != null) result.statistics = statistics;
    if (query != null) result.query = query;
    if (executionparameters != null)
      result.executionparameters.addAll(executionparameters);
    return result;
  }

  QueryExecution._();

  factory QueryExecution.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryExecution.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryExecution',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<QueryExecutionStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: QueryExecutionStatus.create)
    ..aOM<EngineVersion>(44953462, _omitFieldNames ? '' : 'engineversion',
        subBuilder: EngineVersion.create)
    ..aOM<QueryExecutionContext>(
        139302379, _omitFieldNames ? '' : 'queryexecutioncontext',
        subBuilder: QueryExecutionContext.create)
    ..aOM<ManagedQueryResultsConfiguration>(
        159215683, _omitFieldNames ? '' : 'managedqueryresultsconfiguration',
        subBuilder: ManagedQueryResultsConfiguration.create)
    ..aOM<ResultConfiguration>(
        183031503, _omitFieldNames ? '' : 'resultconfiguration',
        subBuilder: ResultConfiguration.create)
    ..aOM<QueryResultsS3AccessGrantsConfiguration>(264403639,
        _omitFieldNames ? '' : 'queryresultss3accessgrantsconfiguration',
        subBuilder: QueryResultsS3AccessGrantsConfiguration.create)
    ..aE<StatementType>(286463655, _omitFieldNames ? '' : 'statementtype',
        enumValues: StatementType.values)
    ..aOS(467615503, _omitFieldNames ? '' : 'queryexecutionid')
    ..aOM<ResultReuseConfiguration>(
        482796543, _omitFieldNames ? '' : 'resultreuseconfiguration',
        subBuilder: ResultReuseConfiguration.create)
    ..aOS(504257527, _omitFieldNames ? '' : 'substatementtype')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..aOM<QueryExecutionStatistics>(
        510636075, _omitFieldNames ? '' : 'statistics',
        subBuilder: QueryExecutionStatistics.create)
    ..aOS(512354180, _omitFieldNames ? '' : 'query')
    ..pPS(527591242, _omitFieldNames ? '' : 'executionparameters')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecution clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecution copyWith(void Function(QueryExecution) updates) =>
      super.copyWith((message) => updates(message as QueryExecution))
          as QueryExecution;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryExecution create() => QueryExecution._();
  @$core.override
  QueryExecution createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryExecution getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryExecution>(create);
  static QueryExecution? _defaultInstance;

  @$pb.TagNumber(6222352)
  QueryExecutionStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(QueryExecutionStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  QueryExecutionStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(44953462)
  EngineVersion get engineversion => $_getN(1);
  @$pb.TagNumber(44953462)
  set engineversion(EngineVersion value) => $_setField(44953462, value);
  @$pb.TagNumber(44953462)
  $core.bool hasEngineversion() => $_has(1);
  @$pb.TagNumber(44953462)
  void clearEngineversion() => $_clearField(44953462);
  @$pb.TagNumber(44953462)
  EngineVersion ensureEngineversion() => $_ensure(1);

  @$pb.TagNumber(139302379)
  QueryExecutionContext get queryexecutioncontext => $_getN(2);
  @$pb.TagNumber(139302379)
  set queryexecutioncontext(QueryExecutionContext value) =>
      $_setField(139302379, value);
  @$pb.TagNumber(139302379)
  $core.bool hasQueryexecutioncontext() => $_has(2);
  @$pb.TagNumber(139302379)
  void clearQueryexecutioncontext() => $_clearField(139302379);
  @$pb.TagNumber(139302379)
  QueryExecutionContext ensureQueryexecutioncontext() => $_ensure(2);

  @$pb.TagNumber(159215683)
  ManagedQueryResultsConfiguration get managedqueryresultsconfiguration =>
      $_getN(3);
  @$pb.TagNumber(159215683)
  set managedqueryresultsconfiguration(
          ManagedQueryResultsConfiguration value) =>
      $_setField(159215683, value);
  @$pb.TagNumber(159215683)
  $core.bool hasManagedqueryresultsconfiguration() => $_has(3);
  @$pb.TagNumber(159215683)
  void clearManagedqueryresultsconfiguration() => $_clearField(159215683);
  @$pb.TagNumber(159215683)
  ManagedQueryResultsConfiguration ensureManagedqueryresultsconfiguration() =>
      $_ensure(3);

  @$pb.TagNumber(183031503)
  ResultConfiguration get resultconfiguration => $_getN(4);
  @$pb.TagNumber(183031503)
  set resultconfiguration(ResultConfiguration value) =>
      $_setField(183031503, value);
  @$pb.TagNumber(183031503)
  $core.bool hasResultconfiguration() => $_has(4);
  @$pb.TagNumber(183031503)
  void clearResultconfiguration() => $_clearField(183031503);
  @$pb.TagNumber(183031503)
  ResultConfiguration ensureResultconfiguration() => $_ensure(4);

  @$pb.TagNumber(264403639)
  QueryResultsS3AccessGrantsConfiguration
      get queryresultss3accessgrantsconfiguration => $_getN(5);
  @$pb.TagNumber(264403639)
  set queryresultss3accessgrantsconfiguration(
          QueryResultsS3AccessGrantsConfiguration value) =>
      $_setField(264403639, value);
  @$pb.TagNumber(264403639)
  $core.bool hasQueryresultss3accessgrantsconfiguration() => $_has(5);
  @$pb.TagNumber(264403639)
  void clearQueryresultss3accessgrantsconfiguration() =>
      $_clearField(264403639);
  @$pb.TagNumber(264403639)
  QueryResultsS3AccessGrantsConfiguration
      ensureQueryresultss3accessgrantsconfiguration() => $_ensure(5);

  @$pb.TagNumber(286463655)
  StatementType get statementtype => $_getN(6);
  @$pb.TagNumber(286463655)
  set statementtype(StatementType value) => $_setField(286463655, value);
  @$pb.TagNumber(286463655)
  $core.bool hasStatementtype() => $_has(6);
  @$pb.TagNumber(286463655)
  void clearStatementtype() => $_clearField(286463655);

  @$pb.TagNumber(467615503)
  $core.String get queryexecutionid => $_getSZ(7);
  @$pb.TagNumber(467615503)
  set queryexecutionid($core.String value) => $_setString(7, value);
  @$pb.TagNumber(467615503)
  $core.bool hasQueryexecutionid() => $_has(7);
  @$pb.TagNumber(467615503)
  void clearQueryexecutionid() => $_clearField(467615503);

  @$pb.TagNumber(482796543)
  ResultReuseConfiguration get resultreuseconfiguration => $_getN(8);
  @$pb.TagNumber(482796543)
  set resultreuseconfiguration(ResultReuseConfiguration value) =>
      $_setField(482796543, value);
  @$pb.TagNumber(482796543)
  $core.bool hasResultreuseconfiguration() => $_has(8);
  @$pb.TagNumber(482796543)
  void clearResultreuseconfiguration() => $_clearField(482796543);
  @$pb.TagNumber(482796543)
  ResultReuseConfiguration ensureResultreuseconfiguration() => $_ensure(8);

  @$pb.TagNumber(504257527)
  $core.String get substatementtype => $_getSZ(9);
  @$pb.TagNumber(504257527)
  set substatementtype($core.String value) => $_setString(9, value);
  @$pb.TagNumber(504257527)
  $core.bool hasSubstatementtype() => $_has(9);
  @$pb.TagNumber(504257527)
  void clearSubstatementtype() => $_clearField(504257527);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(10);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(10, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(10);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(510636075)
  QueryExecutionStatistics get statistics => $_getN(11);
  @$pb.TagNumber(510636075)
  set statistics(QueryExecutionStatistics value) =>
      $_setField(510636075, value);
  @$pb.TagNumber(510636075)
  $core.bool hasStatistics() => $_has(11);
  @$pb.TagNumber(510636075)
  void clearStatistics() => $_clearField(510636075);
  @$pb.TagNumber(510636075)
  QueryExecutionStatistics ensureStatistics() => $_ensure(11);

  @$pb.TagNumber(512354180)
  $core.String get query => $_getSZ(12);
  @$pb.TagNumber(512354180)
  set query($core.String value) => $_setString(12, value);
  @$pb.TagNumber(512354180)
  $core.bool hasQuery() => $_has(12);
  @$pb.TagNumber(512354180)
  void clearQuery() => $_clearField(512354180);

  @$pb.TagNumber(527591242)
  $pb.PbList<$core.String> get executionparameters => $_getList(13);
}

class QueryExecutionContext extends $pb.GeneratedMessage {
  factory QueryExecutionContext({
    $core.String? catalog,
    $core.String? database,
  }) {
    final result = create();
    if (catalog != null) result.catalog = catalog;
    if (database != null) result.database = database;
    return result;
  }

  QueryExecutionContext._();

  factory QueryExecutionContext.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryExecutionContext.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryExecutionContext',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(111213497, _omitFieldNames ? '' : 'catalog')
    ..aOS(278147289, _omitFieldNames ? '' : 'database')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionContext clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionContext copyWith(
          void Function(QueryExecutionContext) updates) =>
      super.copyWith((message) => updates(message as QueryExecutionContext))
          as QueryExecutionContext;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryExecutionContext create() => QueryExecutionContext._();
  @$core.override
  QueryExecutionContext createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryExecutionContext getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryExecutionContext>(create);
  static QueryExecutionContext? _defaultInstance;

  @$pb.TagNumber(111213497)
  $core.String get catalog => $_getSZ(0);
  @$pb.TagNumber(111213497)
  set catalog($core.String value) => $_setString(0, value);
  @$pb.TagNumber(111213497)
  $core.bool hasCatalog() => $_has(0);
  @$pb.TagNumber(111213497)
  void clearCatalog() => $_clearField(111213497);

  @$pb.TagNumber(278147289)
  $core.String get database => $_getSZ(1);
  @$pb.TagNumber(278147289)
  set database($core.String value) => $_setString(1, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(1);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);
}

class QueryExecutionStatistics extends $pb.GeneratedMessage {
  factory QueryExecutionStatistics({
    $fixnum.Int64? serviceprocessingtimeinmillis,
    $fixnum.Int64? datascannedinbytes,
    $fixnum.Int64? servicepreprocessingtimeinmillis,
    $fixnum.Int64? queryqueuetimeinmillis,
    $fixnum.Int64? engineexecutiontimeinmillis,
    $fixnum.Int64? totalexecutiontimeinmillis,
    ResultReuseInformation? resultreuseinformation,
    $core.double? dpucount,
    $core.String? datamanifestlocation,
    $fixnum.Int64? queryplanningtimeinmillis,
  }) {
    final result = create();
    if (serviceprocessingtimeinmillis != null)
      result.serviceprocessingtimeinmillis = serviceprocessingtimeinmillis;
    if (datascannedinbytes != null)
      result.datascannedinbytes = datascannedinbytes;
    if (servicepreprocessingtimeinmillis != null)
      result.servicepreprocessingtimeinmillis =
          servicepreprocessingtimeinmillis;
    if (queryqueuetimeinmillis != null)
      result.queryqueuetimeinmillis = queryqueuetimeinmillis;
    if (engineexecutiontimeinmillis != null)
      result.engineexecutiontimeinmillis = engineexecutiontimeinmillis;
    if (totalexecutiontimeinmillis != null)
      result.totalexecutiontimeinmillis = totalexecutiontimeinmillis;
    if (resultreuseinformation != null)
      result.resultreuseinformation = resultreuseinformation;
    if (dpucount != null) result.dpucount = dpucount;
    if (datamanifestlocation != null)
      result.datamanifestlocation = datamanifestlocation;
    if (queryplanningtimeinmillis != null)
      result.queryplanningtimeinmillis = queryplanningtimeinmillis;
    return result;
  }

  QueryExecutionStatistics._();

  factory QueryExecutionStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryExecutionStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryExecutionStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(81556436, _omitFieldNames ? '' : 'serviceprocessingtimeinmillis')
    ..aInt64(125540604, _omitFieldNames ? '' : 'datascannedinbytes')
    ..aInt64(
        185806899, _omitFieldNames ? '' : 'servicepreprocessingtimeinmillis')
    ..aInt64(221814191, _omitFieldNames ? '' : 'queryqueuetimeinmillis')
    ..aInt64(228680528, _omitFieldNames ? '' : 'engineexecutiontimeinmillis')
    ..aInt64(251249322, _omitFieldNames ? '' : 'totalexecutiontimeinmillis')
    ..aOM<ResultReuseInformation>(
        261773827, _omitFieldNames ? '' : 'resultreuseinformation',
        subBuilder: ResultReuseInformation.create)
    ..aD(296377724, _omitFieldNames ? '' : 'dpucount')
    ..aOS(329970068, _omitFieldNames ? '' : 'datamanifestlocation')
    ..aInt64(482562295, _omitFieldNames ? '' : 'queryplanningtimeinmillis')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionStatistics copyWith(
          void Function(QueryExecutionStatistics) updates) =>
      super.copyWith((message) => updates(message as QueryExecutionStatistics))
          as QueryExecutionStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryExecutionStatistics create() => QueryExecutionStatistics._();
  @$core.override
  QueryExecutionStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryExecutionStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryExecutionStatistics>(create);
  static QueryExecutionStatistics? _defaultInstance;

  @$pb.TagNumber(81556436)
  $fixnum.Int64 get serviceprocessingtimeinmillis => $_getI64(0);
  @$pb.TagNumber(81556436)
  set serviceprocessingtimeinmillis($fixnum.Int64 value) =>
      $_setInt64(0, value);
  @$pb.TagNumber(81556436)
  $core.bool hasServiceprocessingtimeinmillis() => $_has(0);
  @$pb.TagNumber(81556436)
  void clearServiceprocessingtimeinmillis() => $_clearField(81556436);

  @$pb.TagNumber(125540604)
  $fixnum.Int64 get datascannedinbytes => $_getI64(1);
  @$pb.TagNumber(125540604)
  set datascannedinbytes($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(125540604)
  $core.bool hasDatascannedinbytes() => $_has(1);
  @$pb.TagNumber(125540604)
  void clearDatascannedinbytes() => $_clearField(125540604);

  @$pb.TagNumber(185806899)
  $fixnum.Int64 get servicepreprocessingtimeinmillis => $_getI64(2);
  @$pb.TagNumber(185806899)
  set servicepreprocessingtimeinmillis($fixnum.Int64 value) =>
      $_setInt64(2, value);
  @$pb.TagNumber(185806899)
  $core.bool hasServicepreprocessingtimeinmillis() => $_has(2);
  @$pb.TagNumber(185806899)
  void clearServicepreprocessingtimeinmillis() => $_clearField(185806899);

  @$pb.TagNumber(221814191)
  $fixnum.Int64 get queryqueuetimeinmillis => $_getI64(3);
  @$pb.TagNumber(221814191)
  set queryqueuetimeinmillis($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(221814191)
  $core.bool hasQueryqueuetimeinmillis() => $_has(3);
  @$pb.TagNumber(221814191)
  void clearQueryqueuetimeinmillis() => $_clearField(221814191);

  @$pb.TagNumber(228680528)
  $fixnum.Int64 get engineexecutiontimeinmillis => $_getI64(4);
  @$pb.TagNumber(228680528)
  set engineexecutiontimeinmillis($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(228680528)
  $core.bool hasEngineexecutiontimeinmillis() => $_has(4);
  @$pb.TagNumber(228680528)
  void clearEngineexecutiontimeinmillis() => $_clearField(228680528);

  @$pb.TagNumber(251249322)
  $fixnum.Int64 get totalexecutiontimeinmillis => $_getI64(5);
  @$pb.TagNumber(251249322)
  set totalexecutiontimeinmillis($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(251249322)
  $core.bool hasTotalexecutiontimeinmillis() => $_has(5);
  @$pb.TagNumber(251249322)
  void clearTotalexecutiontimeinmillis() => $_clearField(251249322);

  @$pb.TagNumber(261773827)
  ResultReuseInformation get resultreuseinformation => $_getN(6);
  @$pb.TagNumber(261773827)
  set resultreuseinformation(ResultReuseInformation value) =>
      $_setField(261773827, value);
  @$pb.TagNumber(261773827)
  $core.bool hasResultreuseinformation() => $_has(6);
  @$pb.TagNumber(261773827)
  void clearResultreuseinformation() => $_clearField(261773827);
  @$pb.TagNumber(261773827)
  ResultReuseInformation ensureResultreuseinformation() => $_ensure(6);

  @$pb.TagNumber(296377724)
  $core.double get dpucount => $_getN(7);
  @$pb.TagNumber(296377724)
  set dpucount($core.double value) => $_setDouble(7, value);
  @$pb.TagNumber(296377724)
  $core.bool hasDpucount() => $_has(7);
  @$pb.TagNumber(296377724)
  void clearDpucount() => $_clearField(296377724);

  @$pb.TagNumber(329970068)
  $core.String get datamanifestlocation => $_getSZ(8);
  @$pb.TagNumber(329970068)
  set datamanifestlocation($core.String value) => $_setString(8, value);
  @$pb.TagNumber(329970068)
  $core.bool hasDatamanifestlocation() => $_has(8);
  @$pb.TagNumber(329970068)
  void clearDatamanifestlocation() => $_clearField(329970068);

  @$pb.TagNumber(482562295)
  $fixnum.Int64 get queryplanningtimeinmillis => $_getI64(9);
  @$pb.TagNumber(482562295)
  set queryplanningtimeinmillis($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(482562295)
  $core.bool hasQueryplanningtimeinmillis() => $_has(9);
  @$pb.TagNumber(482562295)
  void clearQueryplanningtimeinmillis() => $_clearField(482562295);
}

class QueryExecutionStatus extends $pb.GeneratedMessage {
  factory QueryExecutionStatus({
    AthenaError? athenaerror,
    $core.String? completiondatetime,
    $core.String? statechangereason,
    $core.String? submissiondatetime,
    QueryExecutionState? state,
  }) {
    final result = create();
    if (athenaerror != null) result.athenaerror = athenaerror;
    if (completiondatetime != null)
      result.completiondatetime = completiondatetime;
    if (statechangereason != null) result.statechangereason = statechangereason;
    if (submissiondatetime != null)
      result.submissiondatetime = submissiondatetime;
    if (state != null) result.state = state;
    return result;
  }

  QueryExecutionStatus._();

  factory QueryExecutionStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryExecutionStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryExecutionStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<AthenaError>(82066385, _omitFieldNames ? '' : 'athenaerror',
        subBuilder: AthenaError.create)
    ..aOS(175822779, _omitFieldNames ? '' : 'completiondatetime')
    ..aOS(228940439, _omitFieldNames ? '' : 'statechangereason')
    ..aOS(449650437, _omitFieldNames ? '' : 'submissiondatetime')
    ..aE<QueryExecutionState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: QueryExecutionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionStatus copyWith(void Function(QueryExecutionStatus) updates) =>
      super.copyWith((message) => updates(message as QueryExecutionStatus))
          as QueryExecutionStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryExecutionStatus create() => QueryExecutionStatus._();
  @$core.override
  QueryExecutionStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryExecutionStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryExecutionStatus>(create);
  static QueryExecutionStatus? _defaultInstance;

  @$pb.TagNumber(82066385)
  AthenaError get athenaerror => $_getN(0);
  @$pb.TagNumber(82066385)
  set athenaerror(AthenaError value) => $_setField(82066385, value);
  @$pb.TagNumber(82066385)
  $core.bool hasAthenaerror() => $_has(0);
  @$pb.TagNumber(82066385)
  void clearAthenaerror() => $_clearField(82066385);
  @$pb.TagNumber(82066385)
  AthenaError ensureAthenaerror() => $_ensure(0);

  @$pb.TagNumber(175822779)
  $core.String get completiondatetime => $_getSZ(1);
  @$pb.TagNumber(175822779)
  set completiondatetime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(175822779)
  $core.bool hasCompletiondatetime() => $_has(1);
  @$pb.TagNumber(175822779)
  void clearCompletiondatetime() => $_clearField(175822779);

  @$pb.TagNumber(228940439)
  $core.String get statechangereason => $_getSZ(2);
  @$pb.TagNumber(228940439)
  set statechangereason($core.String value) => $_setString(2, value);
  @$pb.TagNumber(228940439)
  $core.bool hasStatechangereason() => $_has(2);
  @$pb.TagNumber(228940439)
  void clearStatechangereason() => $_clearField(228940439);

  @$pb.TagNumber(449650437)
  $core.String get submissiondatetime => $_getSZ(3);
  @$pb.TagNumber(449650437)
  set submissiondatetime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(449650437)
  $core.bool hasSubmissiondatetime() => $_has(3);
  @$pb.TagNumber(449650437)
  void clearSubmissiondatetime() => $_clearField(449650437);

  @$pb.TagNumber(502047895)
  QueryExecutionState get state => $_getN(4);
  @$pb.TagNumber(502047895)
  set state(QueryExecutionState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(4);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class QueryResultsS3AccessGrantsConfiguration extends $pb.GeneratedMessage {
  factory QueryResultsS3AccessGrantsConfiguration({
    $core.bool? enables3accessgrants,
    AuthenticationType? authenticationtype,
    $core.bool? createuserlevelprefix,
  }) {
    final result = create();
    if (enables3accessgrants != null)
      result.enables3accessgrants = enables3accessgrants;
    if (authenticationtype != null)
      result.authenticationtype = authenticationtype;
    if (createuserlevelprefix != null)
      result.createuserlevelprefix = createuserlevelprefix;
    return result;
  }

  QueryResultsS3AccessGrantsConfiguration._();

  factory QueryResultsS3AccessGrantsConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryResultsS3AccessGrantsConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryResultsS3AccessGrantsConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOB(96318870, _omitFieldNames ? '' : 'enables3accessgrants')
    ..aE<AuthenticationType>(
        277200874, _omitFieldNames ? '' : 'authenticationtype',
        enumValues: AuthenticationType.values)
    ..aOB(493230469, _omitFieldNames ? '' : 'createuserlevelprefix')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryResultsS3AccessGrantsConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryResultsS3AccessGrantsConfiguration copyWith(
          void Function(QueryResultsS3AccessGrantsConfiguration) updates) =>
      super.copyWith((message) =>
              updates(message as QueryResultsS3AccessGrantsConfiguration))
          as QueryResultsS3AccessGrantsConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryResultsS3AccessGrantsConfiguration create() =>
      QueryResultsS3AccessGrantsConfiguration._();
  @$core.override
  QueryResultsS3AccessGrantsConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryResultsS3AccessGrantsConfiguration getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          QueryResultsS3AccessGrantsConfiguration>(create);
  static QueryResultsS3AccessGrantsConfiguration? _defaultInstance;

  @$pb.TagNumber(96318870)
  $core.bool get enables3accessgrants => $_getBF(0);
  @$pb.TagNumber(96318870)
  set enables3accessgrants($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(96318870)
  $core.bool hasEnables3accessgrants() => $_has(0);
  @$pb.TagNumber(96318870)
  void clearEnables3accessgrants() => $_clearField(96318870);

  @$pb.TagNumber(277200874)
  AuthenticationType get authenticationtype => $_getN(1);
  @$pb.TagNumber(277200874)
  set authenticationtype(AuthenticationType value) =>
      $_setField(277200874, value);
  @$pb.TagNumber(277200874)
  $core.bool hasAuthenticationtype() => $_has(1);
  @$pb.TagNumber(277200874)
  void clearAuthenticationtype() => $_clearField(277200874);

  @$pb.TagNumber(493230469)
  $core.bool get createuserlevelprefix => $_getBF(2);
  @$pb.TagNumber(493230469)
  set createuserlevelprefix($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(493230469)
  $core.bool hasCreateuserlevelprefix() => $_has(2);
  @$pb.TagNumber(493230469)
  void clearCreateuserlevelprefix() => $_clearField(493230469);
}

class QueryRuntimeStatistics extends $pb.GeneratedMessage {
  factory QueryRuntimeStatistics({
    QueryStage? outputstage,
    QueryRuntimeStatisticsRows? rows,
    QueryRuntimeStatisticsTimeline? timeline,
  }) {
    final result = create();
    if (outputstage != null) result.outputstage = outputstage;
    if (rows != null) result.rows = rows;
    if (timeline != null) result.timeline = timeline;
    return result;
  }

  QueryRuntimeStatistics._();

  factory QueryRuntimeStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryRuntimeStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryRuntimeStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<QueryStage>(61483955, _omitFieldNames ? '' : 'outputstage',
        subBuilder: QueryStage.create)
    ..aOM<QueryRuntimeStatisticsRows>(174428857, _omitFieldNames ? '' : 'rows',
        subBuilder: QueryRuntimeStatisticsRows.create)
    ..aOM<QueryRuntimeStatisticsTimeline>(
        282559479, _omitFieldNames ? '' : 'timeline',
        subBuilder: QueryRuntimeStatisticsTimeline.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRuntimeStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRuntimeStatistics copyWith(
          void Function(QueryRuntimeStatistics) updates) =>
      super.copyWith((message) => updates(message as QueryRuntimeStatistics))
          as QueryRuntimeStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryRuntimeStatistics create() => QueryRuntimeStatistics._();
  @$core.override
  QueryRuntimeStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryRuntimeStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryRuntimeStatistics>(create);
  static QueryRuntimeStatistics? _defaultInstance;

  @$pb.TagNumber(61483955)
  QueryStage get outputstage => $_getN(0);
  @$pb.TagNumber(61483955)
  set outputstage(QueryStage value) => $_setField(61483955, value);
  @$pb.TagNumber(61483955)
  $core.bool hasOutputstage() => $_has(0);
  @$pb.TagNumber(61483955)
  void clearOutputstage() => $_clearField(61483955);
  @$pb.TagNumber(61483955)
  QueryStage ensureOutputstage() => $_ensure(0);

  @$pb.TagNumber(174428857)
  QueryRuntimeStatisticsRows get rows => $_getN(1);
  @$pb.TagNumber(174428857)
  set rows(QueryRuntimeStatisticsRows value) => $_setField(174428857, value);
  @$pb.TagNumber(174428857)
  $core.bool hasRows() => $_has(1);
  @$pb.TagNumber(174428857)
  void clearRows() => $_clearField(174428857);
  @$pb.TagNumber(174428857)
  QueryRuntimeStatisticsRows ensureRows() => $_ensure(1);

  @$pb.TagNumber(282559479)
  QueryRuntimeStatisticsTimeline get timeline => $_getN(2);
  @$pb.TagNumber(282559479)
  set timeline(QueryRuntimeStatisticsTimeline value) =>
      $_setField(282559479, value);
  @$pb.TagNumber(282559479)
  $core.bool hasTimeline() => $_has(2);
  @$pb.TagNumber(282559479)
  void clearTimeline() => $_clearField(282559479);
  @$pb.TagNumber(282559479)
  QueryRuntimeStatisticsTimeline ensureTimeline() => $_ensure(2);
}

class QueryRuntimeStatisticsRows extends $pb.GeneratedMessage {
  factory QueryRuntimeStatisticsRows({
    $fixnum.Int64? inputrows,
    $fixnum.Int64? outputrows,
    $fixnum.Int64? outputbytes,
    $fixnum.Int64? inputbytes,
  }) {
    final result = create();
    if (inputrows != null) result.inputrows = inputrows;
    if (outputrows != null) result.outputrows = outputrows;
    if (outputbytes != null) result.outputbytes = outputbytes;
    if (inputbytes != null) result.inputbytes = inputbytes;
    return result;
  }

  QueryRuntimeStatisticsRows._();

  factory QueryRuntimeStatisticsRows.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryRuntimeStatisticsRows.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryRuntimeStatisticsRows',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(30551127, _omitFieldNames ? '' : 'inputrows')
    ..aInt64(138873322, _omitFieldNames ? '' : 'outputrows')
    ..aInt64(318971400, _omitFieldNames ? '' : 'outputbytes')
    ..aInt64(431684783, _omitFieldNames ? '' : 'inputbytes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRuntimeStatisticsRows clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRuntimeStatisticsRows copyWith(
          void Function(QueryRuntimeStatisticsRows) updates) =>
      super.copyWith(
              (message) => updates(message as QueryRuntimeStatisticsRows))
          as QueryRuntimeStatisticsRows;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryRuntimeStatisticsRows create() => QueryRuntimeStatisticsRows._();
  @$core.override
  QueryRuntimeStatisticsRows createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryRuntimeStatisticsRows getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryRuntimeStatisticsRows>(create);
  static QueryRuntimeStatisticsRows? _defaultInstance;

  @$pb.TagNumber(30551127)
  $fixnum.Int64 get inputrows => $_getI64(0);
  @$pb.TagNumber(30551127)
  set inputrows($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(30551127)
  $core.bool hasInputrows() => $_has(0);
  @$pb.TagNumber(30551127)
  void clearInputrows() => $_clearField(30551127);

  @$pb.TagNumber(138873322)
  $fixnum.Int64 get outputrows => $_getI64(1);
  @$pb.TagNumber(138873322)
  set outputrows($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(138873322)
  $core.bool hasOutputrows() => $_has(1);
  @$pb.TagNumber(138873322)
  void clearOutputrows() => $_clearField(138873322);

  @$pb.TagNumber(318971400)
  $fixnum.Int64 get outputbytes => $_getI64(2);
  @$pb.TagNumber(318971400)
  set outputbytes($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(318971400)
  $core.bool hasOutputbytes() => $_has(2);
  @$pb.TagNumber(318971400)
  void clearOutputbytes() => $_clearField(318971400);

  @$pb.TagNumber(431684783)
  $fixnum.Int64 get inputbytes => $_getI64(3);
  @$pb.TagNumber(431684783)
  set inputbytes($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(431684783)
  $core.bool hasInputbytes() => $_has(3);
  @$pb.TagNumber(431684783)
  void clearInputbytes() => $_clearField(431684783);
}

class QueryRuntimeStatisticsTimeline extends $pb.GeneratedMessage {
  factory QueryRuntimeStatisticsTimeline({
    $fixnum.Int64? serviceprocessingtimeinmillis,
    $fixnum.Int64? servicepreprocessingtimeinmillis,
    $fixnum.Int64? queryqueuetimeinmillis,
    $fixnum.Int64? engineexecutiontimeinmillis,
    $fixnum.Int64? totalexecutiontimeinmillis,
    $fixnum.Int64? queryplanningtimeinmillis,
  }) {
    final result = create();
    if (serviceprocessingtimeinmillis != null)
      result.serviceprocessingtimeinmillis = serviceprocessingtimeinmillis;
    if (servicepreprocessingtimeinmillis != null)
      result.servicepreprocessingtimeinmillis =
          servicepreprocessingtimeinmillis;
    if (queryqueuetimeinmillis != null)
      result.queryqueuetimeinmillis = queryqueuetimeinmillis;
    if (engineexecutiontimeinmillis != null)
      result.engineexecutiontimeinmillis = engineexecutiontimeinmillis;
    if (totalexecutiontimeinmillis != null)
      result.totalexecutiontimeinmillis = totalexecutiontimeinmillis;
    if (queryplanningtimeinmillis != null)
      result.queryplanningtimeinmillis = queryplanningtimeinmillis;
    return result;
  }

  QueryRuntimeStatisticsTimeline._();

  factory QueryRuntimeStatisticsTimeline.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryRuntimeStatisticsTimeline.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryRuntimeStatisticsTimeline',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(81556436, _omitFieldNames ? '' : 'serviceprocessingtimeinmillis')
    ..aInt64(
        185806899, _omitFieldNames ? '' : 'servicepreprocessingtimeinmillis')
    ..aInt64(221814191, _omitFieldNames ? '' : 'queryqueuetimeinmillis')
    ..aInt64(228680528, _omitFieldNames ? '' : 'engineexecutiontimeinmillis')
    ..aInt64(251249322, _omitFieldNames ? '' : 'totalexecutiontimeinmillis')
    ..aInt64(482562295, _omitFieldNames ? '' : 'queryplanningtimeinmillis')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRuntimeStatisticsTimeline clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRuntimeStatisticsTimeline copyWith(
          void Function(QueryRuntimeStatisticsTimeline) updates) =>
      super.copyWith(
              (message) => updates(message as QueryRuntimeStatisticsTimeline))
          as QueryRuntimeStatisticsTimeline;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryRuntimeStatisticsTimeline create() =>
      QueryRuntimeStatisticsTimeline._();
  @$core.override
  QueryRuntimeStatisticsTimeline createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryRuntimeStatisticsTimeline getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryRuntimeStatisticsTimeline>(create);
  static QueryRuntimeStatisticsTimeline? _defaultInstance;

  @$pb.TagNumber(81556436)
  $fixnum.Int64 get serviceprocessingtimeinmillis => $_getI64(0);
  @$pb.TagNumber(81556436)
  set serviceprocessingtimeinmillis($fixnum.Int64 value) =>
      $_setInt64(0, value);
  @$pb.TagNumber(81556436)
  $core.bool hasServiceprocessingtimeinmillis() => $_has(0);
  @$pb.TagNumber(81556436)
  void clearServiceprocessingtimeinmillis() => $_clearField(81556436);

  @$pb.TagNumber(185806899)
  $fixnum.Int64 get servicepreprocessingtimeinmillis => $_getI64(1);
  @$pb.TagNumber(185806899)
  set servicepreprocessingtimeinmillis($fixnum.Int64 value) =>
      $_setInt64(1, value);
  @$pb.TagNumber(185806899)
  $core.bool hasServicepreprocessingtimeinmillis() => $_has(1);
  @$pb.TagNumber(185806899)
  void clearServicepreprocessingtimeinmillis() => $_clearField(185806899);

  @$pb.TagNumber(221814191)
  $fixnum.Int64 get queryqueuetimeinmillis => $_getI64(2);
  @$pb.TagNumber(221814191)
  set queryqueuetimeinmillis($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(221814191)
  $core.bool hasQueryqueuetimeinmillis() => $_has(2);
  @$pb.TagNumber(221814191)
  void clearQueryqueuetimeinmillis() => $_clearField(221814191);

  @$pb.TagNumber(228680528)
  $fixnum.Int64 get engineexecutiontimeinmillis => $_getI64(3);
  @$pb.TagNumber(228680528)
  set engineexecutiontimeinmillis($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(228680528)
  $core.bool hasEngineexecutiontimeinmillis() => $_has(3);
  @$pb.TagNumber(228680528)
  void clearEngineexecutiontimeinmillis() => $_clearField(228680528);

  @$pb.TagNumber(251249322)
  $fixnum.Int64 get totalexecutiontimeinmillis => $_getI64(4);
  @$pb.TagNumber(251249322)
  set totalexecutiontimeinmillis($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(251249322)
  $core.bool hasTotalexecutiontimeinmillis() => $_has(4);
  @$pb.TagNumber(251249322)
  void clearTotalexecutiontimeinmillis() => $_clearField(251249322);

  @$pb.TagNumber(482562295)
  $fixnum.Int64 get queryplanningtimeinmillis => $_getI64(5);
  @$pb.TagNumber(482562295)
  set queryplanningtimeinmillis($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(482562295)
  $core.bool hasQueryplanningtimeinmillis() => $_has(5);
  @$pb.TagNumber(482562295)
  void clearQueryplanningtimeinmillis() => $_clearField(482562295);
}

class QueryStage extends $pb.GeneratedMessage {
  factory QueryStage({
    $fixnum.Int64? inputrows,
    QueryStagePlanNode? querystageplan,
    $core.Iterable<QueryStage>? substages,
    $fixnum.Int64? outputrows,
    $fixnum.Int64? outputbytes,
    $fixnum.Int64? stageid,
    $fixnum.Int64? executiontime,
    $fixnum.Int64? inputbytes,
    $core.String? state,
  }) {
    final result = create();
    if (inputrows != null) result.inputrows = inputrows;
    if (querystageplan != null) result.querystageplan = querystageplan;
    if (substages != null) result.substages.addAll(substages);
    if (outputrows != null) result.outputrows = outputrows;
    if (outputbytes != null) result.outputbytes = outputbytes;
    if (stageid != null) result.stageid = stageid;
    if (executiontime != null) result.executiontime = executiontime;
    if (inputbytes != null) result.inputbytes = inputbytes;
    if (state != null) result.state = state;
    return result;
  }

  QueryStage._();

  factory QueryStage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryStage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryStage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(30551127, _omitFieldNames ? '' : 'inputrows')
    ..aOM<QueryStagePlanNode>(36313545, _omitFieldNames ? '' : 'querystageplan',
        subBuilder: QueryStagePlanNode.create)
    ..pPM<QueryStage>(114732779, _omitFieldNames ? '' : 'substages',
        subBuilder: QueryStage.create)
    ..aInt64(138873322, _omitFieldNames ? '' : 'outputrows')
    ..aInt64(318971400, _omitFieldNames ? '' : 'outputbytes')
    ..aInt64(328165497, _omitFieldNames ? '' : 'stageid')
    ..aInt64(379716053, _omitFieldNames ? '' : 'executiontime')
    ..aInt64(431684783, _omitFieldNames ? '' : 'inputbytes')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStage copyWith(void Function(QueryStage) updates) =>
      super.copyWith((message) => updates(message as QueryStage)) as QueryStage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryStage create() => QueryStage._();
  @$core.override
  QueryStage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryStage getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryStage>(create);
  static QueryStage? _defaultInstance;

  @$pb.TagNumber(30551127)
  $fixnum.Int64 get inputrows => $_getI64(0);
  @$pb.TagNumber(30551127)
  set inputrows($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(30551127)
  $core.bool hasInputrows() => $_has(0);
  @$pb.TagNumber(30551127)
  void clearInputrows() => $_clearField(30551127);

  @$pb.TagNumber(36313545)
  QueryStagePlanNode get querystageplan => $_getN(1);
  @$pb.TagNumber(36313545)
  set querystageplan(QueryStagePlanNode value) => $_setField(36313545, value);
  @$pb.TagNumber(36313545)
  $core.bool hasQuerystageplan() => $_has(1);
  @$pb.TagNumber(36313545)
  void clearQuerystageplan() => $_clearField(36313545);
  @$pb.TagNumber(36313545)
  QueryStagePlanNode ensureQuerystageplan() => $_ensure(1);

  @$pb.TagNumber(114732779)
  $pb.PbList<QueryStage> get substages => $_getList(2);

  @$pb.TagNumber(138873322)
  $fixnum.Int64 get outputrows => $_getI64(3);
  @$pb.TagNumber(138873322)
  set outputrows($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(138873322)
  $core.bool hasOutputrows() => $_has(3);
  @$pb.TagNumber(138873322)
  void clearOutputrows() => $_clearField(138873322);

  @$pb.TagNumber(318971400)
  $fixnum.Int64 get outputbytes => $_getI64(4);
  @$pb.TagNumber(318971400)
  set outputbytes($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(318971400)
  $core.bool hasOutputbytes() => $_has(4);
  @$pb.TagNumber(318971400)
  void clearOutputbytes() => $_clearField(318971400);

  @$pb.TagNumber(328165497)
  $fixnum.Int64 get stageid => $_getI64(5);
  @$pb.TagNumber(328165497)
  set stageid($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(328165497)
  $core.bool hasStageid() => $_has(5);
  @$pb.TagNumber(328165497)
  void clearStageid() => $_clearField(328165497);

  @$pb.TagNumber(379716053)
  $fixnum.Int64 get executiontime => $_getI64(6);
  @$pb.TagNumber(379716053)
  set executiontime($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(379716053)
  $core.bool hasExecutiontime() => $_has(6);
  @$pb.TagNumber(379716053)
  void clearExecutiontime() => $_clearField(379716053);

  @$pb.TagNumber(431684783)
  $fixnum.Int64 get inputbytes => $_getI64(7);
  @$pb.TagNumber(431684783)
  set inputbytes($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(431684783)
  $core.bool hasInputbytes() => $_has(7);
  @$pb.TagNumber(431684783)
  void clearInputbytes() => $_clearField(431684783);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(8);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(8, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(8);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class QueryStagePlanNode extends $pb.GeneratedMessage {
  factory QueryStagePlanNode({
    $core.String? identifier,
    $core.Iterable<$core.String>? remotesources,
    $core.Iterable<QueryStagePlanNode>? children,
    $core.String? name,
  }) {
    final result = create();
    if (identifier != null) result.identifier = identifier;
    if (remotesources != null) result.remotesources.addAll(remotesources);
    if (children != null) result.children.addAll(children);
    if (name != null) result.name = name;
    return result;
  }

  QueryStagePlanNode._();

  factory QueryStagePlanNode.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryStagePlanNode.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryStagePlanNode',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(41865311, _omitFieldNames ? '' : 'identifier')
    ..pPS(143743504, _omitFieldNames ? '' : 'remotesources')
    ..pPM<QueryStagePlanNode>(188567027, _omitFieldNames ? '' : 'children',
        subBuilder: QueryStagePlanNode.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStagePlanNode clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStagePlanNode copyWith(void Function(QueryStagePlanNode) updates) =>
      super.copyWith((message) => updates(message as QueryStagePlanNode))
          as QueryStagePlanNode;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryStagePlanNode create() => QueryStagePlanNode._();
  @$core.override
  QueryStagePlanNode createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryStagePlanNode getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryStagePlanNode>(create);
  static QueryStagePlanNode? _defaultInstance;

  @$pb.TagNumber(41865311)
  $core.String get identifier => $_getSZ(0);
  @$pb.TagNumber(41865311)
  set identifier($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41865311)
  $core.bool hasIdentifier() => $_has(0);
  @$pb.TagNumber(41865311)
  void clearIdentifier() => $_clearField(41865311);

  @$pb.TagNumber(143743504)
  $pb.PbList<$core.String> get remotesources => $_getList(1);

  @$pb.TagNumber(188567027)
  $pb.PbList<QueryStagePlanNode> get children => $_getList(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ResourceNotFoundException extends $pb.GeneratedMessage {
  factory ResourceNotFoundException({
    $core.String? message,
    $core.String? resourcename,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (resourcename != null) result.resourcename = resourcename;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(269834071, _omitFieldNames ? '' : 'resourcename')
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

  @$pb.TagNumber(269834071)
  $core.String get resourcename => $_getSZ(1);
  @$pb.TagNumber(269834071)
  set resourcename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(269834071)
  $core.bool hasResourcename() => $_has(1);
  @$pb.TagNumber(269834071)
  void clearResourcename() => $_clearField(269834071);
}

class ResultConfiguration extends $pb.GeneratedMessage {
  factory ResultConfiguration({
    AclConfiguration? aclconfiguration,
    $core.String? outputlocation,
    $core.String? expectedbucketowner,
    EncryptionConfiguration? encryptionconfiguration,
  }) {
    final result = create();
    if (aclconfiguration != null) result.aclconfiguration = aclconfiguration;
    if (outputlocation != null) result.outputlocation = outputlocation;
    if (expectedbucketowner != null)
      result.expectedbucketowner = expectedbucketowner;
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    return result;
  }

  ResultConfiguration._();

  factory ResultConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResultConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResultConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<AclConfiguration>(9179900, _omitFieldNames ? '' : 'aclconfiguration',
        subBuilder: AclConfiguration.create)
    ..aOS(67991028, _omitFieldNames ? '' : 'outputlocation')
    ..aOS(132066983, _omitFieldNames ? '' : 'expectedbucketowner')
    ..aOM<EncryptionConfiguration>(
        225764215, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultConfiguration copyWith(void Function(ResultConfiguration) updates) =>
      super.copyWith((message) => updates(message as ResultConfiguration))
          as ResultConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResultConfiguration create() => ResultConfiguration._();
  @$core.override
  ResultConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResultConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResultConfiguration>(create);
  static ResultConfiguration? _defaultInstance;

  @$pb.TagNumber(9179900)
  AclConfiguration get aclconfiguration => $_getN(0);
  @$pb.TagNumber(9179900)
  set aclconfiguration(AclConfiguration value) => $_setField(9179900, value);
  @$pb.TagNumber(9179900)
  $core.bool hasAclconfiguration() => $_has(0);
  @$pb.TagNumber(9179900)
  void clearAclconfiguration() => $_clearField(9179900);
  @$pb.TagNumber(9179900)
  AclConfiguration ensureAclconfiguration() => $_ensure(0);

  @$pb.TagNumber(67991028)
  $core.String get outputlocation => $_getSZ(1);
  @$pb.TagNumber(67991028)
  set outputlocation($core.String value) => $_setString(1, value);
  @$pb.TagNumber(67991028)
  $core.bool hasOutputlocation() => $_has(1);
  @$pb.TagNumber(67991028)
  void clearOutputlocation() => $_clearField(67991028);

  @$pb.TagNumber(132066983)
  $core.String get expectedbucketowner => $_getSZ(2);
  @$pb.TagNumber(132066983)
  set expectedbucketowner($core.String value) => $_setString(2, value);
  @$pb.TagNumber(132066983)
  $core.bool hasExpectedbucketowner() => $_has(2);
  @$pb.TagNumber(132066983)
  void clearExpectedbucketowner() => $_clearField(132066983);

  @$pb.TagNumber(225764215)
  EncryptionConfiguration get encryptionconfiguration => $_getN(3);
  @$pb.TagNumber(225764215)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(225764215, value);
  @$pb.TagNumber(225764215)
  $core.bool hasEncryptionconfiguration() => $_has(3);
  @$pb.TagNumber(225764215)
  void clearEncryptionconfiguration() => $_clearField(225764215);
  @$pb.TagNumber(225764215)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(3);
}

class ResultConfigurationUpdates extends $pb.GeneratedMessage {
  factory ResultConfigurationUpdates({
    AclConfiguration? aclconfiguration,
    $core.String? outputlocation,
    $core.String? expectedbucketowner,
    $core.bool? removeoutputlocation,
    $core.bool? removeexpectedbucketowner,
    EncryptionConfiguration? encryptionconfiguration,
    $core.bool? removeencryptionconfiguration,
    $core.bool? removeaclconfiguration,
  }) {
    final result = create();
    if (aclconfiguration != null) result.aclconfiguration = aclconfiguration;
    if (outputlocation != null) result.outputlocation = outputlocation;
    if (expectedbucketowner != null)
      result.expectedbucketowner = expectedbucketowner;
    if (removeoutputlocation != null)
      result.removeoutputlocation = removeoutputlocation;
    if (removeexpectedbucketowner != null)
      result.removeexpectedbucketowner = removeexpectedbucketowner;
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (removeencryptionconfiguration != null)
      result.removeencryptionconfiguration = removeencryptionconfiguration;
    if (removeaclconfiguration != null)
      result.removeaclconfiguration = removeaclconfiguration;
    return result;
  }

  ResultConfigurationUpdates._();

  factory ResultConfigurationUpdates.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResultConfigurationUpdates.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResultConfigurationUpdates',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<AclConfiguration>(9179900, _omitFieldNames ? '' : 'aclconfiguration',
        subBuilder: AclConfiguration.create)
    ..aOS(67991028, _omitFieldNames ? '' : 'outputlocation')
    ..aOS(132066983, _omitFieldNames ? '' : 'expectedbucketowner')
    ..aOB(197075308, _omitFieldNames ? '' : 'removeoutputlocation')
    ..aOB(198850607, _omitFieldNames ? '' : 'removeexpectedbucketowner')
    ..aOM<EncryptionConfiguration>(
        225764215, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOB(294004335, _omitFieldNames ? '' : 'removeencryptionconfiguration')
    ..aOB(378505796, _omitFieldNames ? '' : 'removeaclconfiguration')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultConfigurationUpdates clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultConfigurationUpdates copyWith(
          void Function(ResultConfigurationUpdates) updates) =>
      super.copyWith(
              (message) => updates(message as ResultConfigurationUpdates))
          as ResultConfigurationUpdates;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResultConfigurationUpdates create() => ResultConfigurationUpdates._();
  @$core.override
  ResultConfigurationUpdates createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResultConfigurationUpdates getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResultConfigurationUpdates>(create);
  static ResultConfigurationUpdates? _defaultInstance;

  @$pb.TagNumber(9179900)
  AclConfiguration get aclconfiguration => $_getN(0);
  @$pb.TagNumber(9179900)
  set aclconfiguration(AclConfiguration value) => $_setField(9179900, value);
  @$pb.TagNumber(9179900)
  $core.bool hasAclconfiguration() => $_has(0);
  @$pb.TagNumber(9179900)
  void clearAclconfiguration() => $_clearField(9179900);
  @$pb.TagNumber(9179900)
  AclConfiguration ensureAclconfiguration() => $_ensure(0);

  @$pb.TagNumber(67991028)
  $core.String get outputlocation => $_getSZ(1);
  @$pb.TagNumber(67991028)
  set outputlocation($core.String value) => $_setString(1, value);
  @$pb.TagNumber(67991028)
  $core.bool hasOutputlocation() => $_has(1);
  @$pb.TagNumber(67991028)
  void clearOutputlocation() => $_clearField(67991028);

  @$pb.TagNumber(132066983)
  $core.String get expectedbucketowner => $_getSZ(2);
  @$pb.TagNumber(132066983)
  set expectedbucketowner($core.String value) => $_setString(2, value);
  @$pb.TagNumber(132066983)
  $core.bool hasExpectedbucketowner() => $_has(2);
  @$pb.TagNumber(132066983)
  void clearExpectedbucketowner() => $_clearField(132066983);

  @$pb.TagNumber(197075308)
  $core.bool get removeoutputlocation => $_getBF(3);
  @$pb.TagNumber(197075308)
  set removeoutputlocation($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(197075308)
  $core.bool hasRemoveoutputlocation() => $_has(3);
  @$pb.TagNumber(197075308)
  void clearRemoveoutputlocation() => $_clearField(197075308);

  @$pb.TagNumber(198850607)
  $core.bool get removeexpectedbucketowner => $_getBF(4);
  @$pb.TagNumber(198850607)
  set removeexpectedbucketowner($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(198850607)
  $core.bool hasRemoveexpectedbucketowner() => $_has(4);
  @$pb.TagNumber(198850607)
  void clearRemoveexpectedbucketowner() => $_clearField(198850607);

  @$pb.TagNumber(225764215)
  EncryptionConfiguration get encryptionconfiguration => $_getN(5);
  @$pb.TagNumber(225764215)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(225764215, value);
  @$pb.TagNumber(225764215)
  $core.bool hasEncryptionconfiguration() => $_has(5);
  @$pb.TagNumber(225764215)
  void clearEncryptionconfiguration() => $_clearField(225764215);
  @$pb.TagNumber(225764215)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(5);

  @$pb.TagNumber(294004335)
  $core.bool get removeencryptionconfiguration => $_getBF(6);
  @$pb.TagNumber(294004335)
  set removeencryptionconfiguration($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(294004335)
  $core.bool hasRemoveencryptionconfiguration() => $_has(6);
  @$pb.TagNumber(294004335)
  void clearRemoveencryptionconfiguration() => $_clearField(294004335);

  @$pb.TagNumber(378505796)
  $core.bool get removeaclconfiguration => $_getBF(7);
  @$pb.TagNumber(378505796)
  set removeaclconfiguration($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(378505796)
  $core.bool hasRemoveaclconfiguration() => $_has(7);
  @$pb.TagNumber(378505796)
  void clearRemoveaclconfiguration() => $_clearField(378505796);
}

class ResultReuseByAgeConfiguration extends $pb.GeneratedMessage {
  factory ResultReuseByAgeConfiguration({
    $core.int? maxageinminutes,
    $core.bool? enabled,
  }) {
    final result = create();
    if (maxageinminutes != null) result.maxageinminutes = maxageinminutes;
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  ResultReuseByAgeConfiguration._();

  factory ResultReuseByAgeConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResultReuseByAgeConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResultReuseByAgeConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aI(211566877, _omitFieldNames ? '' : 'maxageinminutes')
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultReuseByAgeConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultReuseByAgeConfiguration copyWith(
          void Function(ResultReuseByAgeConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as ResultReuseByAgeConfiguration))
          as ResultReuseByAgeConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResultReuseByAgeConfiguration create() =>
      ResultReuseByAgeConfiguration._();
  @$core.override
  ResultReuseByAgeConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResultReuseByAgeConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResultReuseByAgeConfiguration>(create);
  static ResultReuseByAgeConfiguration? _defaultInstance;

  @$pb.TagNumber(211566877)
  $core.int get maxageinminutes => $_getIZ(0);
  @$pb.TagNumber(211566877)
  set maxageinminutes($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(211566877)
  $core.bool hasMaxageinminutes() => $_has(0);
  @$pb.TagNumber(211566877)
  void clearMaxageinminutes() => $_clearField(211566877);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);
}

class ResultReuseConfiguration extends $pb.GeneratedMessage {
  factory ResultReuseConfiguration({
    ResultReuseByAgeConfiguration? resultreusebyageconfiguration,
  }) {
    final result = create();
    if (resultreusebyageconfiguration != null)
      result.resultreusebyageconfiguration = resultreusebyageconfiguration;
    return result;
  }

  ResultReuseConfiguration._();

  factory ResultReuseConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResultReuseConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResultReuseConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<ResultReuseByAgeConfiguration>(
        13481299, _omitFieldNames ? '' : 'resultreusebyageconfiguration',
        subBuilder: ResultReuseByAgeConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultReuseConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultReuseConfiguration copyWith(
          void Function(ResultReuseConfiguration) updates) =>
      super.copyWith((message) => updates(message as ResultReuseConfiguration))
          as ResultReuseConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResultReuseConfiguration create() => ResultReuseConfiguration._();
  @$core.override
  ResultReuseConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResultReuseConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResultReuseConfiguration>(create);
  static ResultReuseConfiguration? _defaultInstance;

  @$pb.TagNumber(13481299)
  ResultReuseByAgeConfiguration get resultreusebyageconfiguration => $_getN(0);
  @$pb.TagNumber(13481299)
  set resultreusebyageconfiguration(ResultReuseByAgeConfiguration value) =>
      $_setField(13481299, value);
  @$pb.TagNumber(13481299)
  $core.bool hasResultreusebyageconfiguration() => $_has(0);
  @$pb.TagNumber(13481299)
  void clearResultreusebyageconfiguration() => $_clearField(13481299);
  @$pb.TagNumber(13481299)
  ResultReuseByAgeConfiguration ensureResultreusebyageconfiguration() =>
      $_ensure(0);
}

class ResultReuseInformation extends $pb.GeneratedMessage {
  factory ResultReuseInformation({
    $core.bool? reusedpreviousresult,
  }) {
    final result = create();
    if (reusedpreviousresult != null)
      result.reusedpreviousresult = reusedpreviousresult;
    return result;
  }

  ResultReuseInformation._();

  factory ResultReuseInformation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResultReuseInformation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResultReuseInformation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOB(448068450, _omitFieldNames ? '' : 'reusedpreviousresult')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultReuseInformation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultReuseInformation copyWith(
          void Function(ResultReuseInformation) updates) =>
      super.copyWith((message) => updates(message as ResultReuseInformation))
          as ResultReuseInformation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResultReuseInformation create() => ResultReuseInformation._();
  @$core.override
  ResultReuseInformation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResultReuseInformation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResultReuseInformation>(create);
  static ResultReuseInformation? _defaultInstance;

  @$pb.TagNumber(448068450)
  $core.bool get reusedpreviousresult => $_getBF(0);
  @$pb.TagNumber(448068450)
  set reusedpreviousresult($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(448068450)
  $core.bool hasReusedpreviousresult() => $_has(0);
  @$pb.TagNumber(448068450)
  void clearReusedpreviousresult() => $_clearField(448068450);
}

class ResultSet extends $pb.GeneratedMessage {
  factory ResultSet({
    $core.Iterable<Row>? rows,
    ResultSetMetadata? resultsetmetadata,
  }) {
    final result = create();
    if (rows != null) result.rows.addAll(rows);
    if (resultsetmetadata != null) result.resultsetmetadata = resultsetmetadata;
    return result;
  }

  ResultSet._();

  factory ResultSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResultSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResultSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<Row>(174428857, _omitFieldNames ? '' : 'rows', subBuilder: Row.create)
    ..aOM<ResultSetMetadata>(
        317769226, _omitFieldNames ? '' : 'resultsetmetadata',
        subBuilder: ResultSetMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultSet copyWith(void Function(ResultSet) updates) =>
      super.copyWith((message) => updates(message as ResultSet)) as ResultSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResultSet create() => ResultSet._();
  @$core.override
  ResultSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResultSet getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ResultSet>(create);
  static ResultSet? _defaultInstance;

  @$pb.TagNumber(174428857)
  $pb.PbList<Row> get rows => $_getList(0);

  @$pb.TagNumber(317769226)
  ResultSetMetadata get resultsetmetadata => $_getN(1);
  @$pb.TagNumber(317769226)
  set resultsetmetadata(ResultSetMetadata value) =>
      $_setField(317769226, value);
  @$pb.TagNumber(317769226)
  $core.bool hasResultsetmetadata() => $_has(1);
  @$pb.TagNumber(317769226)
  void clearResultsetmetadata() => $_clearField(317769226);
  @$pb.TagNumber(317769226)
  ResultSetMetadata ensureResultsetmetadata() => $_ensure(1);
}

class ResultSetMetadata extends $pb.GeneratedMessage {
  factory ResultSetMetadata({
    $core.Iterable<ColumnInfo>? columninfo,
  }) {
    final result = create();
    if (columninfo != null) result.columninfo.addAll(columninfo);
    return result;
  }

  ResultSetMetadata._();

  factory ResultSetMetadata.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResultSetMetadata.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResultSetMetadata',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<ColumnInfo>(364742404, _omitFieldNames ? '' : 'columninfo',
        subBuilder: ColumnInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultSetMetadata clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResultSetMetadata copyWith(void Function(ResultSetMetadata) updates) =>
      super.copyWith((message) => updates(message as ResultSetMetadata))
          as ResultSetMetadata;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResultSetMetadata create() => ResultSetMetadata._();
  @$core.override
  ResultSetMetadata createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResultSetMetadata getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResultSetMetadata>(create);
  static ResultSetMetadata? _defaultInstance;

  @$pb.TagNumber(364742404)
  $pb.PbList<ColumnInfo> get columninfo => $_getList(0);
}

class Row extends $pb.GeneratedMessage {
  factory Row({
    $core.Iterable<Datum>? data,
  }) {
    final result = create();
    if (data != null) result.data.addAll(data);
    return result;
  }

  Row._();

  factory Row.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Row.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Row',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<Datum>(525498822, _omitFieldNames ? '' : 'data',
        subBuilder: Datum.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Row clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Row copyWith(void Function(Row) updates) =>
      super.copyWith((message) => updates(message as Row)) as Row;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Row create() => Row._();
  @$core.override
  Row createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Row getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Row>(create);
  static Row? _defaultInstance;

  @$pb.TagNumber(525498822)
  $pb.PbList<Datum> get data => $_getList(0);
}

class S3LoggingConfiguration extends $pb.GeneratedMessage {
  factory S3LoggingConfiguration({
    $core.String? kmskey,
    $core.String? loglocation,
    $core.bool? enabled,
  }) {
    final result = create();
    if (kmskey != null) result.kmskey = kmskey;
    if (loglocation != null) result.loglocation = loglocation;
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  S3LoggingConfiguration._();

  factory S3LoggingConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory S3LoggingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'S3LoggingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(114561194, _omitFieldNames ? '' : 'kmskey')
    ..aOS(192028619, _omitFieldNames ? '' : 'loglocation')
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3LoggingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3LoggingConfiguration copyWith(
          void Function(S3LoggingConfiguration) updates) =>
      super.copyWith((message) => updates(message as S3LoggingConfiguration))
          as S3LoggingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static S3LoggingConfiguration create() => S3LoggingConfiguration._();
  @$core.override
  S3LoggingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static S3LoggingConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<S3LoggingConfiguration>(create);
  static S3LoggingConfiguration? _defaultInstance;

  @$pb.TagNumber(114561194)
  $core.String get kmskey => $_getSZ(0);
  @$pb.TagNumber(114561194)
  set kmskey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(114561194)
  $core.bool hasKmskey() => $_has(0);
  @$pb.TagNumber(114561194)
  void clearKmskey() => $_clearField(114561194);

  @$pb.TagNumber(192028619)
  $core.String get loglocation => $_getSZ(1);
  @$pb.TagNumber(192028619)
  set loglocation($core.String value) => $_setString(1, value);
  @$pb.TagNumber(192028619)
  $core.bool hasLoglocation() => $_has(1);
  @$pb.TagNumber(192028619)
  void clearLoglocation() => $_clearField(192028619);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(2);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(2);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);
}

class SessionAlreadyExistsException extends $pb.GeneratedMessage {
  factory SessionAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  SessionAlreadyExistsException._();

  factory SessionAlreadyExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SessionAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SessionAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionAlreadyExistsException copyWith(
          void Function(SessionAlreadyExistsException) updates) =>
      super.copyWith(
              (message) => updates(message as SessionAlreadyExistsException))
          as SessionAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SessionAlreadyExistsException create() =>
      SessionAlreadyExistsException._();
  @$core.override
  SessionAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SessionAlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SessionAlreadyExistsException>(create);
  static SessionAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class SessionConfiguration extends $pb.GeneratedMessage {
  factory SessionConfiguration({
    $fixnum.Int64? idletimeoutseconds,
    EncryptionConfiguration? encryptionconfiguration,
    $core.String? executionrole,
    $core.String? workingdirectory,
    $core.int? sessionidletimeoutinminutes,
  }) {
    final result = create();
    if (idletimeoutseconds != null)
      result.idletimeoutseconds = idletimeoutseconds;
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (executionrole != null) result.executionrole = executionrole;
    if (workingdirectory != null) result.workingdirectory = workingdirectory;
    if (sessionidletimeoutinminutes != null)
      result.sessionidletimeoutinminutes = sessionidletimeoutinminutes;
    return result;
  }

  SessionConfiguration._();

  factory SessionConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SessionConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SessionConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(61279260, _omitFieldNames ? '' : 'idletimeoutseconds')
    ..aOM<EncryptionConfiguration>(
        225764215, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOS(253307658, _omitFieldNames ? '' : 'executionrole')
    ..aOS(478970252, _omitFieldNames ? '' : 'workingdirectory')
    ..aI(515304989, _omitFieldNames ? '' : 'sessionidletimeoutinminutes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionConfiguration copyWith(void Function(SessionConfiguration) updates) =>
      super.copyWith((message) => updates(message as SessionConfiguration))
          as SessionConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SessionConfiguration create() => SessionConfiguration._();
  @$core.override
  SessionConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SessionConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SessionConfiguration>(create);
  static SessionConfiguration? _defaultInstance;

  @$pb.TagNumber(61279260)
  $fixnum.Int64 get idletimeoutseconds => $_getI64(0);
  @$pb.TagNumber(61279260)
  set idletimeoutseconds($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(61279260)
  $core.bool hasIdletimeoutseconds() => $_has(0);
  @$pb.TagNumber(61279260)
  void clearIdletimeoutseconds() => $_clearField(61279260);

  @$pb.TagNumber(225764215)
  EncryptionConfiguration get encryptionconfiguration => $_getN(1);
  @$pb.TagNumber(225764215)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(225764215, value);
  @$pb.TagNumber(225764215)
  $core.bool hasEncryptionconfiguration() => $_has(1);
  @$pb.TagNumber(225764215)
  void clearEncryptionconfiguration() => $_clearField(225764215);
  @$pb.TagNumber(225764215)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(1);

  @$pb.TagNumber(253307658)
  $core.String get executionrole => $_getSZ(2);
  @$pb.TagNumber(253307658)
  set executionrole($core.String value) => $_setString(2, value);
  @$pb.TagNumber(253307658)
  $core.bool hasExecutionrole() => $_has(2);
  @$pb.TagNumber(253307658)
  void clearExecutionrole() => $_clearField(253307658);

  @$pb.TagNumber(478970252)
  $core.String get workingdirectory => $_getSZ(3);
  @$pb.TagNumber(478970252)
  set workingdirectory($core.String value) => $_setString(3, value);
  @$pb.TagNumber(478970252)
  $core.bool hasWorkingdirectory() => $_has(3);
  @$pb.TagNumber(478970252)
  void clearWorkingdirectory() => $_clearField(478970252);

  @$pb.TagNumber(515304989)
  $core.int get sessionidletimeoutinminutes => $_getIZ(4);
  @$pb.TagNumber(515304989)
  set sessionidletimeoutinminutes($core.int value) =>
      $_setSignedInt32(4, value);
  @$pb.TagNumber(515304989)
  $core.bool hasSessionidletimeoutinminutes() => $_has(4);
  @$pb.TagNumber(515304989)
  void clearSessionidletimeoutinminutes() => $_clearField(515304989);
}

class SessionStatistics extends $pb.GeneratedMessage {
  factory SessionStatistics({
    $fixnum.Int64? dpuexecutioninmillis,
  }) {
    final result = create();
    if (dpuexecutioninmillis != null)
      result.dpuexecutioninmillis = dpuexecutioninmillis;
    return result;
  }

  SessionStatistics._();

  factory SessionStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SessionStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SessionStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aInt64(174857936, _omitFieldNames ? '' : 'dpuexecutioninmillis')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionStatistics copyWith(void Function(SessionStatistics) updates) =>
      super.copyWith((message) => updates(message as SessionStatistics))
          as SessionStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SessionStatistics create() => SessionStatistics._();
  @$core.override
  SessionStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SessionStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SessionStatistics>(create);
  static SessionStatistics? _defaultInstance;

  @$pb.TagNumber(174857936)
  $fixnum.Int64 get dpuexecutioninmillis => $_getI64(0);
  @$pb.TagNumber(174857936)
  set dpuexecutioninmillis($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(174857936)
  $core.bool hasDpuexecutioninmillis() => $_has(0);
  @$pb.TagNumber(174857936)
  void clearDpuexecutioninmillis() => $_clearField(174857936);
}

class SessionStatus extends $pb.GeneratedMessage {
  factory SessionStatus({
    $core.String? startdatetime,
    $core.String? statechangereason,
    $core.String? enddatetime,
    $core.String? lastmodifieddatetime,
    $core.String? idlesincedatetime,
    SessionState? state,
  }) {
    final result = create();
    if (startdatetime != null) result.startdatetime = startdatetime;
    if (statechangereason != null) result.statechangereason = statechangereason;
    if (enddatetime != null) result.enddatetime = enddatetime;
    if (lastmodifieddatetime != null)
      result.lastmodifieddatetime = lastmodifieddatetime;
    if (idlesincedatetime != null) result.idlesincedatetime = idlesincedatetime;
    if (state != null) result.state = state;
    return result;
  }

  SessionStatus._();

  factory SessionStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SessionStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SessionStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(88518355, _omitFieldNames ? '' : 'startdatetime')
    ..aOS(228940439, _omitFieldNames ? '' : 'statechangereason')
    ..aOS(372000488, _omitFieldNames ? '' : 'enddatetime')
    ..aOS(406490164, _omitFieldNames ? '' : 'lastmodifieddatetime')
    ..aOS(454742375, _omitFieldNames ? '' : 'idlesincedatetime')
    ..aE<SessionState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: SessionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionStatus copyWith(void Function(SessionStatus) updates) =>
      super.copyWith((message) => updates(message as SessionStatus))
          as SessionStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SessionStatus create() => SessionStatus._();
  @$core.override
  SessionStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SessionStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SessionStatus>(create);
  static SessionStatus? _defaultInstance;

  @$pb.TagNumber(88518355)
  $core.String get startdatetime => $_getSZ(0);
  @$pb.TagNumber(88518355)
  set startdatetime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88518355)
  $core.bool hasStartdatetime() => $_has(0);
  @$pb.TagNumber(88518355)
  void clearStartdatetime() => $_clearField(88518355);

  @$pb.TagNumber(228940439)
  $core.String get statechangereason => $_getSZ(1);
  @$pb.TagNumber(228940439)
  set statechangereason($core.String value) => $_setString(1, value);
  @$pb.TagNumber(228940439)
  $core.bool hasStatechangereason() => $_has(1);
  @$pb.TagNumber(228940439)
  void clearStatechangereason() => $_clearField(228940439);

  @$pb.TagNumber(372000488)
  $core.String get enddatetime => $_getSZ(2);
  @$pb.TagNumber(372000488)
  set enddatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(372000488)
  $core.bool hasEnddatetime() => $_has(2);
  @$pb.TagNumber(372000488)
  void clearEnddatetime() => $_clearField(372000488);

  @$pb.TagNumber(406490164)
  $core.String get lastmodifieddatetime => $_getSZ(3);
  @$pb.TagNumber(406490164)
  set lastmodifieddatetime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(406490164)
  $core.bool hasLastmodifieddatetime() => $_has(3);
  @$pb.TagNumber(406490164)
  void clearLastmodifieddatetime() => $_clearField(406490164);

  @$pb.TagNumber(454742375)
  $core.String get idlesincedatetime => $_getSZ(4);
  @$pb.TagNumber(454742375)
  set idlesincedatetime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(454742375)
  $core.bool hasIdlesincedatetime() => $_has(4);
  @$pb.TagNumber(454742375)
  void clearIdlesincedatetime() => $_clearField(454742375);

  @$pb.TagNumber(502047895)
  SessionState get state => $_getN(5);
  @$pb.TagNumber(502047895)
  set state(SessionState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class SessionSummary extends $pb.GeneratedMessage {
  factory SessionSummary({
    SessionStatus? status,
    $core.String? sessionid,
    EngineVersion? engineversion,
    $core.String? description,
    $core.String? notebookversion,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (sessionid != null) result.sessionid = sessionid;
    if (engineversion != null) result.engineversion = engineversion;
    if (description != null) result.description = description;
    if (notebookversion != null) result.notebookversion = notebookversion;
    return result;
  }

  SessionSummary._();

  factory SessionSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SessionSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SessionSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<SessionStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: SessionStatus.create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aOM<EngineVersion>(44953462, _omitFieldNames ? '' : 'engineversion',
        subBuilder: EngineVersion.create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(528689837, _omitFieldNames ? '' : 'notebookversion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionSummary copyWith(void Function(SessionSummary) updates) =>
      super.copyWith((message) => updates(message as SessionSummary))
          as SessionSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SessionSummary create() => SessionSummary._();
  @$core.override
  SessionSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SessionSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SessionSummary>(create);
  static SessionSummary? _defaultInstance;

  @$pb.TagNumber(6222352)
  SessionStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(SessionStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  SessionStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(1);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(1);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(44953462)
  EngineVersion get engineversion => $_getN(2);
  @$pb.TagNumber(44953462)
  set engineversion(EngineVersion value) => $_setField(44953462, value);
  @$pb.TagNumber(44953462)
  $core.bool hasEngineversion() => $_has(2);
  @$pb.TagNumber(44953462)
  void clearEngineversion() => $_clearField(44953462);
  @$pb.TagNumber(44953462)
  EngineVersion ensureEngineversion() => $_ensure(2);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(528689837)
  $core.String get notebookversion => $_getSZ(4);
  @$pb.TagNumber(528689837)
  set notebookversion($core.String value) => $_setString(4, value);
  @$pb.TagNumber(528689837)
  $core.bool hasNotebookversion() => $_has(4);
  @$pb.TagNumber(528689837)
  void clearNotebookversion() => $_clearField(528689837);
}

class StartCalculationExecutionRequest extends $pb.GeneratedMessage {
  factory StartCalculationExecutionRequest({
    $core.String? sessionid,
    $core.String? codeblock,
    $core.String? description,
    $core.String? clientrequesttoken,
    CalculationConfiguration? calculationconfiguration,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    if (codeblock != null) result.codeblock = codeblock;
    if (description != null) result.description = description;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (calculationconfiguration != null)
      result.calculationconfiguration = calculationconfiguration;
    return result;
  }

  StartCalculationExecutionRequest._();

  factory StartCalculationExecutionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartCalculationExecutionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartCalculationExecutionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aOS(23945838, _omitFieldNames ? '' : 'codeblock')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..aOM<CalculationConfiguration>(
        459676797, _omitFieldNames ? '' : 'calculationconfiguration',
        subBuilder: CalculationConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartCalculationExecutionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartCalculationExecutionRequest copyWith(
          void Function(StartCalculationExecutionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StartCalculationExecutionRequest))
          as StartCalculationExecutionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartCalculationExecutionRequest create() =>
      StartCalculationExecutionRequest._();
  @$core.override
  StartCalculationExecutionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartCalculationExecutionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartCalculationExecutionRequest>(
          create);
  static StartCalculationExecutionRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(23945838)
  $core.String get codeblock => $_getSZ(1);
  @$pb.TagNumber(23945838)
  set codeblock($core.String value) => $_setString(1, value);
  @$pb.TagNumber(23945838)
  $core.bool hasCodeblock() => $_has(1);
  @$pb.TagNumber(23945838)
  void clearCodeblock() => $_clearField(23945838);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(3);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(3);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(459676797)
  CalculationConfiguration get calculationconfiguration => $_getN(4);
  @$pb.TagNumber(459676797)
  set calculationconfiguration(CalculationConfiguration value) =>
      $_setField(459676797, value);
  @$pb.TagNumber(459676797)
  $core.bool hasCalculationconfiguration() => $_has(4);
  @$pb.TagNumber(459676797)
  void clearCalculationconfiguration() => $_clearField(459676797);
  @$pb.TagNumber(459676797)
  CalculationConfiguration ensureCalculationconfiguration() => $_ensure(4);
}

class StartCalculationExecutionResponse extends $pb.GeneratedMessage {
  factory StartCalculationExecutionResponse({
    $core.String? calculationexecutionid,
    CalculationExecutionState? state,
  }) {
    final result = create();
    if (calculationexecutionid != null)
      result.calculationexecutionid = calculationexecutionid;
    if (state != null) result.state = state;
    return result;
  }

  StartCalculationExecutionResponse._();

  factory StartCalculationExecutionResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartCalculationExecutionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartCalculationExecutionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(80028050, _omitFieldNames ? '' : 'calculationexecutionid')
    ..aE<CalculationExecutionState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: CalculationExecutionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartCalculationExecutionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartCalculationExecutionResponse copyWith(
          void Function(StartCalculationExecutionResponse) updates) =>
      super.copyWith((message) =>
              updates(message as StartCalculationExecutionResponse))
          as StartCalculationExecutionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartCalculationExecutionResponse create() =>
      StartCalculationExecutionResponse._();
  @$core.override
  StartCalculationExecutionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartCalculationExecutionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartCalculationExecutionResponse>(
          create);
  static StartCalculationExecutionResponse? _defaultInstance;

  @$pb.TagNumber(80028050)
  $core.String get calculationexecutionid => $_getSZ(0);
  @$pb.TagNumber(80028050)
  set calculationexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(80028050)
  $core.bool hasCalculationexecutionid() => $_has(0);
  @$pb.TagNumber(80028050)
  void clearCalculationexecutionid() => $_clearField(80028050);

  @$pb.TagNumber(502047895)
  CalculationExecutionState get state => $_getN(1);
  @$pb.TagNumber(502047895)
  set state(CalculationExecutionState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(1);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class StartQueryExecutionInput extends $pb.GeneratedMessage {
  factory StartQueryExecutionInput({
    QueryExecutionContext? queryexecutioncontext,
    ResultConfiguration? resultconfiguration,
    EngineConfiguration? engineconfiguration,
    $core.String? querystring,
    $core.String? clientrequesttoken,
    ResultReuseConfiguration? resultreuseconfiguration,
    $core.String? workgroup,
    $core.Iterable<$core.String>? executionparameters,
  }) {
    final result = create();
    if (queryexecutioncontext != null)
      result.queryexecutioncontext = queryexecutioncontext;
    if (resultconfiguration != null)
      result.resultconfiguration = resultconfiguration;
    if (engineconfiguration != null)
      result.engineconfiguration = engineconfiguration;
    if (querystring != null) result.querystring = querystring;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (resultreuseconfiguration != null)
      result.resultreuseconfiguration = resultreuseconfiguration;
    if (workgroup != null) result.workgroup = workgroup;
    if (executionparameters != null)
      result.executionparameters.addAll(executionparameters);
    return result;
  }

  StartQueryExecutionInput._();

  factory StartQueryExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartQueryExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartQueryExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<QueryExecutionContext>(
        139302379, _omitFieldNames ? '' : 'queryexecutioncontext',
        subBuilder: QueryExecutionContext.create)
    ..aOM<ResultConfiguration>(
        183031503, _omitFieldNames ? '' : 'resultconfiguration',
        subBuilder: ResultConfiguration.create)
    ..aOM<EngineConfiguration>(
        341629412, _omitFieldNames ? '' : 'engineconfiguration',
        subBuilder: EngineConfiguration.create)
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..aOM<ResultReuseConfiguration>(
        482796543, _omitFieldNames ? '' : 'resultreuseconfiguration',
        subBuilder: ResultReuseConfiguration.create)
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..pPS(527591242, _omitFieldNames ? '' : 'executionparameters')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryExecutionInput copyWith(
          void Function(StartQueryExecutionInput) updates) =>
      super.copyWith((message) => updates(message as StartQueryExecutionInput))
          as StartQueryExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartQueryExecutionInput create() => StartQueryExecutionInput._();
  @$core.override
  StartQueryExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartQueryExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartQueryExecutionInput>(create);
  static StartQueryExecutionInput? _defaultInstance;

  @$pb.TagNumber(139302379)
  QueryExecutionContext get queryexecutioncontext => $_getN(0);
  @$pb.TagNumber(139302379)
  set queryexecutioncontext(QueryExecutionContext value) =>
      $_setField(139302379, value);
  @$pb.TagNumber(139302379)
  $core.bool hasQueryexecutioncontext() => $_has(0);
  @$pb.TagNumber(139302379)
  void clearQueryexecutioncontext() => $_clearField(139302379);
  @$pb.TagNumber(139302379)
  QueryExecutionContext ensureQueryexecutioncontext() => $_ensure(0);

  @$pb.TagNumber(183031503)
  ResultConfiguration get resultconfiguration => $_getN(1);
  @$pb.TagNumber(183031503)
  set resultconfiguration(ResultConfiguration value) =>
      $_setField(183031503, value);
  @$pb.TagNumber(183031503)
  $core.bool hasResultconfiguration() => $_has(1);
  @$pb.TagNumber(183031503)
  void clearResultconfiguration() => $_clearField(183031503);
  @$pb.TagNumber(183031503)
  ResultConfiguration ensureResultconfiguration() => $_ensure(1);

  @$pb.TagNumber(341629412)
  EngineConfiguration get engineconfiguration => $_getN(2);
  @$pb.TagNumber(341629412)
  set engineconfiguration(EngineConfiguration value) =>
      $_setField(341629412, value);
  @$pb.TagNumber(341629412)
  $core.bool hasEngineconfiguration() => $_has(2);
  @$pb.TagNumber(341629412)
  void clearEngineconfiguration() => $_clearField(341629412);
  @$pb.TagNumber(341629412)
  EngineConfiguration ensureEngineconfiguration() => $_ensure(2);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(3);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(3, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(3);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(4);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(4);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(482796543)
  ResultReuseConfiguration get resultreuseconfiguration => $_getN(5);
  @$pb.TagNumber(482796543)
  set resultreuseconfiguration(ResultReuseConfiguration value) =>
      $_setField(482796543, value);
  @$pb.TagNumber(482796543)
  $core.bool hasResultreuseconfiguration() => $_has(5);
  @$pb.TagNumber(482796543)
  void clearResultreuseconfiguration() => $_clearField(482796543);
  @$pb.TagNumber(482796543)
  ResultReuseConfiguration ensureResultreuseconfiguration() => $_ensure(5);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(6);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(6, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(6);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(527591242)
  $pb.PbList<$core.String> get executionparameters => $_getList(7);
}

class StartQueryExecutionOutput extends $pb.GeneratedMessage {
  factory StartQueryExecutionOutput({
    $core.String? queryexecutionid,
  }) {
    final result = create();
    if (queryexecutionid != null) result.queryexecutionid = queryexecutionid;
    return result;
  }

  StartQueryExecutionOutput._();

  factory StartQueryExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartQueryExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartQueryExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(467615503, _omitFieldNames ? '' : 'queryexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartQueryExecutionOutput copyWith(
          void Function(StartQueryExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as StartQueryExecutionOutput))
          as StartQueryExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartQueryExecutionOutput create() => StartQueryExecutionOutput._();
  @$core.override
  StartQueryExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartQueryExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartQueryExecutionOutput>(create);
  static StartQueryExecutionOutput? _defaultInstance;

  @$pb.TagNumber(467615503)
  $core.String get queryexecutionid => $_getSZ(0);
  @$pb.TagNumber(467615503)
  set queryexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(467615503)
  $core.bool hasQueryexecutionid() => $_has(0);
  @$pb.TagNumber(467615503)
  void clearQueryexecutionid() => $_clearField(467615503);
}

class StartSessionRequest extends $pb.GeneratedMessage {
  factory StartSessionRequest({
    $core.String? description,
    $core.String? executionrole,
    EngineConfiguration? engineconfiguration,
    MonitoringConfiguration? monitoringconfiguration,
    $core.Iterable<Tag>? tags,
    $core.bool? copyworkgrouptags,
    $core.String? clientrequesttoken,
    $core.String? workgroup,
    $core.int? sessionidletimeoutinminutes,
    $core.String? notebookversion,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (executionrole != null) result.executionrole = executionrole;
    if (engineconfiguration != null)
      result.engineconfiguration = engineconfiguration;
    if (monitoringconfiguration != null)
      result.monitoringconfiguration = monitoringconfiguration;
    if (tags != null) result.tags.addAll(tags);
    if (copyworkgrouptags != null) result.copyworkgrouptags = copyworkgrouptags;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (workgroup != null) result.workgroup = workgroup;
    if (sessionidletimeoutinminutes != null)
      result.sessionidletimeoutinminutes = sessionidletimeoutinminutes;
    if (notebookversion != null) result.notebookversion = notebookversion;
    return result;
  }

  StartSessionRequest._();

  factory StartSessionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartSessionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartSessionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(253307658, _omitFieldNames ? '' : 'executionrole')
    ..aOM<EngineConfiguration>(
        341629412, _omitFieldNames ? '' : 'engineconfiguration',
        subBuilder: EngineConfiguration.create)
    ..aOM<MonitoringConfiguration>(
        364891928, _omitFieldNames ? '' : 'monitoringconfiguration',
        subBuilder: MonitoringConfiguration.create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOB(429355506, _omitFieldNames ? '' : 'copyworkgrouptags')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..aI(515304989, _omitFieldNames ? '' : 'sessionidletimeoutinminutes')
    ..aOS(528689837, _omitFieldNames ? '' : 'notebookversion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSessionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSessionRequest copyWith(void Function(StartSessionRequest) updates) =>
      super.copyWith((message) => updates(message as StartSessionRequest))
          as StartSessionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartSessionRequest create() => StartSessionRequest._();
  @$core.override
  StartSessionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartSessionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartSessionRequest>(create);
  static StartSessionRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(253307658)
  $core.String get executionrole => $_getSZ(1);
  @$pb.TagNumber(253307658)
  set executionrole($core.String value) => $_setString(1, value);
  @$pb.TagNumber(253307658)
  $core.bool hasExecutionrole() => $_has(1);
  @$pb.TagNumber(253307658)
  void clearExecutionrole() => $_clearField(253307658);

  @$pb.TagNumber(341629412)
  EngineConfiguration get engineconfiguration => $_getN(2);
  @$pb.TagNumber(341629412)
  set engineconfiguration(EngineConfiguration value) =>
      $_setField(341629412, value);
  @$pb.TagNumber(341629412)
  $core.bool hasEngineconfiguration() => $_has(2);
  @$pb.TagNumber(341629412)
  void clearEngineconfiguration() => $_clearField(341629412);
  @$pb.TagNumber(341629412)
  EngineConfiguration ensureEngineconfiguration() => $_ensure(2);

  @$pb.TagNumber(364891928)
  MonitoringConfiguration get monitoringconfiguration => $_getN(3);
  @$pb.TagNumber(364891928)
  set monitoringconfiguration(MonitoringConfiguration value) =>
      $_setField(364891928, value);
  @$pb.TagNumber(364891928)
  $core.bool hasMonitoringconfiguration() => $_has(3);
  @$pb.TagNumber(364891928)
  void clearMonitoringconfiguration() => $_clearField(364891928);
  @$pb.TagNumber(364891928)
  MonitoringConfiguration ensureMonitoringconfiguration() => $_ensure(3);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(4);

  @$pb.TagNumber(429355506)
  $core.bool get copyworkgrouptags => $_getBF(5);
  @$pb.TagNumber(429355506)
  set copyworkgrouptags($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(429355506)
  $core.bool hasCopyworkgrouptags() => $_has(5);
  @$pb.TagNumber(429355506)
  void clearCopyworkgrouptags() => $_clearField(429355506);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(6);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(6, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(6);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(7);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(7, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(7);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);

  @$pb.TagNumber(515304989)
  $core.int get sessionidletimeoutinminutes => $_getIZ(8);
  @$pb.TagNumber(515304989)
  set sessionidletimeoutinminutes($core.int value) =>
      $_setSignedInt32(8, value);
  @$pb.TagNumber(515304989)
  $core.bool hasSessionidletimeoutinminutes() => $_has(8);
  @$pb.TagNumber(515304989)
  void clearSessionidletimeoutinminutes() => $_clearField(515304989);

  @$pb.TagNumber(528689837)
  $core.String get notebookversion => $_getSZ(9);
  @$pb.TagNumber(528689837)
  set notebookversion($core.String value) => $_setString(9, value);
  @$pb.TagNumber(528689837)
  $core.bool hasNotebookversion() => $_has(9);
  @$pb.TagNumber(528689837)
  void clearNotebookversion() => $_clearField(528689837);
}

class StartSessionResponse extends $pb.GeneratedMessage {
  factory StartSessionResponse({
    $core.String? sessionid,
    SessionState? state,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    if (state != null) result.state = state;
    return result;
  }

  StartSessionResponse._();

  factory StartSessionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartSessionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartSessionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aE<SessionState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: SessionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSessionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSessionResponse copyWith(void Function(StartSessionResponse) updates) =>
      super.copyWith((message) => updates(message as StartSessionResponse))
          as StartSessionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartSessionResponse create() => StartSessionResponse._();
  @$core.override
  StartSessionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartSessionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartSessionResponse>(create);
  static StartSessionResponse? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(502047895)
  SessionState get state => $_getN(1);
  @$pb.TagNumber(502047895)
  set state(SessionState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(1);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class StopCalculationExecutionRequest extends $pb.GeneratedMessage {
  factory StopCalculationExecutionRequest({
    $core.String? calculationexecutionid,
  }) {
    final result = create();
    if (calculationexecutionid != null)
      result.calculationexecutionid = calculationexecutionid;
    return result;
  }

  StopCalculationExecutionRequest._();

  factory StopCalculationExecutionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopCalculationExecutionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopCalculationExecutionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(80028050, _omitFieldNames ? '' : 'calculationexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopCalculationExecutionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopCalculationExecutionRequest copyWith(
          void Function(StopCalculationExecutionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StopCalculationExecutionRequest))
          as StopCalculationExecutionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopCalculationExecutionRequest create() =>
      StopCalculationExecutionRequest._();
  @$core.override
  StopCalculationExecutionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopCalculationExecutionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopCalculationExecutionRequest>(
          create);
  static StopCalculationExecutionRequest? _defaultInstance;

  @$pb.TagNumber(80028050)
  $core.String get calculationexecutionid => $_getSZ(0);
  @$pb.TagNumber(80028050)
  set calculationexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(80028050)
  $core.bool hasCalculationexecutionid() => $_has(0);
  @$pb.TagNumber(80028050)
  void clearCalculationexecutionid() => $_clearField(80028050);
}

class StopCalculationExecutionResponse extends $pb.GeneratedMessage {
  factory StopCalculationExecutionResponse({
    CalculationExecutionState? state,
  }) {
    final result = create();
    if (state != null) result.state = state;
    return result;
  }

  StopCalculationExecutionResponse._();

  factory StopCalculationExecutionResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopCalculationExecutionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopCalculationExecutionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<CalculationExecutionState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: CalculationExecutionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopCalculationExecutionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopCalculationExecutionResponse copyWith(
          void Function(StopCalculationExecutionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as StopCalculationExecutionResponse))
          as StopCalculationExecutionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopCalculationExecutionResponse create() =>
      StopCalculationExecutionResponse._();
  @$core.override
  StopCalculationExecutionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopCalculationExecutionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopCalculationExecutionResponse>(
          create);
  static StopCalculationExecutionResponse? _defaultInstance;

  @$pb.TagNumber(502047895)
  CalculationExecutionState get state => $_getN(0);
  @$pb.TagNumber(502047895)
  set state(CalculationExecutionState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(0);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class StopQueryExecutionInput extends $pb.GeneratedMessage {
  factory StopQueryExecutionInput({
    $core.String? queryexecutionid,
  }) {
    final result = create();
    if (queryexecutionid != null) result.queryexecutionid = queryexecutionid;
    return result;
  }

  StopQueryExecutionInput._();

  factory StopQueryExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopQueryExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopQueryExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(467615503, _omitFieldNames ? '' : 'queryexecutionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopQueryExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopQueryExecutionInput copyWith(
          void Function(StopQueryExecutionInput) updates) =>
      super.copyWith((message) => updates(message as StopQueryExecutionInput))
          as StopQueryExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopQueryExecutionInput create() => StopQueryExecutionInput._();
  @$core.override
  StopQueryExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopQueryExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopQueryExecutionInput>(create);
  static StopQueryExecutionInput? _defaultInstance;

  @$pb.TagNumber(467615503)
  $core.String get queryexecutionid => $_getSZ(0);
  @$pb.TagNumber(467615503)
  set queryexecutionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(467615503)
  $core.bool hasQueryexecutionid() => $_has(0);
  @$pb.TagNumber(467615503)
  void clearQueryexecutionid() => $_clearField(467615503);
}

class StopQueryExecutionOutput extends $pb.GeneratedMessage {
  factory StopQueryExecutionOutput() => create();

  StopQueryExecutionOutput._();

  factory StopQueryExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopQueryExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopQueryExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopQueryExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopQueryExecutionOutput copyWith(
          void Function(StopQueryExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as StopQueryExecutionOutput))
          as StopQueryExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopQueryExecutionOutput create() => StopQueryExecutionOutput._();
  @$core.override
  StopQueryExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopQueryExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopQueryExecutionOutput>(create);
  static StopQueryExecutionOutput? _defaultInstance;
}

class TableMetadata extends $pb.GeneratedMessage {
  factory TableMetadata({
    $core.Iterable<Column>? columns,
    $core.Iterable<Column>? partitionkeys,
    $core.String? name,
    $core.String? tabletype,
    $core.String? createtime,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
    $core.String? lastaccesstime,
  }) {
    final result = create();
    if (columns != null) result.columns.addAll(columns);
    if (partitionkeys != null) result.partitionkeys.addAll(partitionkeys);
    if (name != null) result.name = name;
    if (tabletype != null) result.tabletype = tabletype;
    if (createtime != null) result.createtime = createtime;
    if (parameters != null) result.parameters.addEntries(parameters);
    if (lastaccesstime != null) result.lastaccesstime = lastaccesstime;
    return result;
  }

  TableMetadata._();

  factory TableMetadata.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableMetadata.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableMetadata',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPM<Column>(169177053, _omitFieldNames ? '' : 'columns',
        subBuilder: Column.create)
    ..pPM<Column>(200562986, _omitFieldNames ? '' : 'partitionkeys',
        subBuilder: Column.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(476171176, _omitFieldNames ? '' : 'tabletype')
    ..aOS(490895933, _omitFieldNames ? '' : 'createtime')
    ..m<$core.String, $core.String>(
        494900218, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'TableMetadata.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..aOS(516574551, _omitFieldNames ? '' : 'lastaccesstime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableMetadata clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableMetadata copyWith(void Function(TableMetadata) updates) =>
      super.copyWith((message) => updates(message as TableMetadata))
          as TableMetadata;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableMetadata create() => TableMetadata._();
  @$core.override
  TableMetadata createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableMetadata getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableMetadata>(create);
  static TableMetadata? _defaultInstance;

  @$pb.TagNumber(169177053)
  $pb.PbList<Column> get columns => $_getList(0);

  @$pb.TagNumber(200562986)
  $pb.PbList<Column> get partitionkeys => $_getList(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(476171176)
  $core.String get tabletype => $_getSZ(3);
  @$pb.TagNumber(476171176)
  set tabletype($core.String value) => $_setString(3, value);
  @$pb.TagNumber(476171176)
  $core.bool hasTabletype() => $_has(3);
  @$pb.TagNumber(476171176)
  void clearTabletype() => $_clearField(476171176);

  @$pb.TagNumber(490895933)
  $core.String get createtime => $_getSZ(4);
  @$pb.TagNumber(490895933)
  set createtime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(490895933)
  $core.bool hasCreatetime() => $_has(4);
  @$pb.TagNumber(490895933)
  void clearCreatetime() => $_clearField(490895933);

  @$pb.TagNumber(494900218)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(5);

  @$pb.TagNumber(516574551)
  $core.String get lastaccesstime => $_getSZ(6);
  @$pb.TagNumber(516574551)
  set lastaccesstime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(516574551)
  $core.bool hasLastaccesstime() => $_has(6);
  @$pb.TagNumber(516574551)
  void clearLastaccesstime() => $_clearField(516574551);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
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

class TerminateSessionRequest extends $pb.GeneratedMessage {
  factory TerminateSessionRequest({
    $core.String? sessionid,
  }) {
    final result = create();
    if (sessionid != null) result.sessionid = sessionid;
    return result;
  }

  TerminateSessionRequest._();

  factory TerminateSessionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TerminateSessionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TerminateSessionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TerminateSessionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TerminateSessionRequest copyWith(
          void Function(TerminateSessionRequest) updates) =>
      super.copyWith((message) => updates(message as TerminateSessionRequest))
          as TerminateSessionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TerminateSessionRequest create() => TerminateSessionRequest._();
  @$core.override
  TerminateSessionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TerminateSessionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TerminateSessionRequest>(create);
  static TerminateSessionRequest? _defaultInstance;

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(0);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(0);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);
}

class TerminateSessionResponse extends $pb.GeneratedMessage {
  factory TerminateSessionResponse({
    SessionState? state,
  }) {
    final result = create();
    if (state != null) result.state = state;
    return result;
  }

  TerminateSessionResponse._();

  factory TerminateSessionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TerminateSessionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TerminateSessionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<SessionState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: SessionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TerminateSessionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TerminateSessionResponse copyWith(
          void Function(TerminateSessionResponse) updates) =>
      super.copyWith((message) => updates(message as TerminateSessionResponse))
          as TerminateSessionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TerminateSessionResponse create() => TerminateSessionResponse._();
  @$core.override
  TerminateSessionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TerminateSessionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TerminateSessionResponse>(create);
  static TerminateSessionResponse? _defaultInstance;

  @$pb.TagNumber(502047895)
  SessionState get state => $_getN(0);
  @$pb.TagNumber(502047895)
  set state(SessionState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(0);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class TooManyRequestsException extends $pb.GeneratedMessage {
  factory TooManyRequestsException({
    ThrottleReason? reason,
    $core.String? message,
  }) {
    final result = create();
    if (reason != null) result.reason = reason;
    if (message != null) result.message = message;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aE<ThrottleReason>(20005178, _omitFieldNames ? '' : 'reason',
        enumValues: ThrottleReason.values)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(20005178)
  ThrottleReason get reason => $_getN(0);
  @$pb.TagNumber(20005178)
  set reason(ThrottleReason value) => $_setField(20005178, value);
  @$pb.TagNumber(20005178)
  $core.bool hasReason() => $_has(0);
  @$pb.TagNumber(20005178)
  void clearReason() => $_clearField(20005178);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class UnprocessedNamedQueryId extends $pb.GeneratedMessage {
  factory UnprocessedNamedQueryId({
    $core.String? errorcode,
    $core.String? namedqueryid,
    $core.String? errormessage,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (namedqueryid != null) result.namedqueryid = namedqueryid;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  UnprocessedNamedQueryId._();

  factory UnprocessedNamedQueryId.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnprocessedNamedQueryId.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnprocessedNamedQueryId',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(330896872, _omitFieldNames ? '' : 'namedqueryid')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedNamedQueryId clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedNamedQueryId copyWith(
          void Function(UnprocessedNamedQueryId) updates) =>
      super.copyWith((message) => updates(message as UnprocessedNamedQueryId))
          as UnprocessedNamedQueryId;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnprocessedNamedQueryId create() => UnprocessedNamedQueryId._();
  @$core.override
  UnprocessedNamedQueryId createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnprocessedNamedQueryId getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnprocessedNamedQueryId>(create);
  static UnprocessedNamedQueryId? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(330896872)
  $core.String get namedqueryid => $_getSZ(1);
  @$pb.TagNumber(330896872)
  set namedqueryid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(330896872)
  $core.bool hasNamedqueryid() => $_has(1);
  @$pb.TagNumber(330896872)
  void clearNamedqueryid() => $_clearField(330896872);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(2);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(2);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class UnprocessedPreparedStatementName extends $pb.GeneratedMessage {
  factory UnprocessedPreparedStatementName({
    $core.String? statementname,
    $core.String? errorcode,
    $core.String? errormessage,
  }) {
    final result = create();
    if (statementname != null) result.statementname = statementname;
    if (errorcode != null) result.errorcode = errorcode;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  UnprocessedPreparedStatementName._();

  factory UnprocessedPreparedStatementName.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnprocessedPreparedStatementName.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnprocessedPreparedStatementName',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedPreparedStatementName clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedPreparedStatementName copyWith(
          void Function(UnprocessedPreparedStatementName) updates) =>
      super.copyWith(
              (message) => updates(message as UnprocessedPreparedStatementName))
          as UnprocessedPreparedStatementName;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnprocessedPreparedStatementName create() =>
      UnprocessedPreparedStatementName._();
  @$core.override
  UnprocessedPreparedStatementName createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnprocessedPreparedStatementName getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnprocessedPreparedStatementName>(
          create);
  static UnprocessedPreparedStatementName? _defaultInstance;

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(0);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(0);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(1);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(1);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(2);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(2);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class UnprocessedQueryExecutionId extends $pb.GeneratedMessage {
  factory UnprocessedQueryExecutionId({
    $core.String? errorcode,
    $core.String? queryexecutionid,
    $core.String? errormessage,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (queryexecutionid != null) result.queryexecutionid = queryexecutionid;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  UnprocessedQueryExecutionId._();

  factory UnprocessedQueryExecutionId.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnprocessedQueryExecutionId.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnprocessedQueryExecutionId',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(467615503, _omitFieldNames ? '' : 'queryexecutionid')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedQueryExecutionId clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedQueryExecutionId copyWith(
          void Function(UnprocessedQueryExecutionId) updates) =>
      super.copyWith(
              (message) => updates(message as UnprocessedQueryExecutionId))
          as UnprocessedQueryExecutionId;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnprocessedQueryExecutionId create() =>
      UnprocessedQueryExecutionId._();
  @$core.override
  UnprocessedQueryExecutionId createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnprocessedQueryExecutionId getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnprocessedQueryExecutionId>(create);
  static UnprocessedQueryExecutionId? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(467615503)
  $core.String get queryexecutionid => $_getSZ(1);
  @$pb.TagNumber(467615503)
  set queryexecutionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(467615503)
  $core.bool hasQueryexecutionid() => $_has(1);
  @$pb.TagNumber(467615503)
  void clearQueryexecutionid() => $_clearField(467615503);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(2);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(2);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
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

class UpdateCapacityReservationInput extends $pb.GeneratedMessage {
  factory UpdateCapacityReservationInput({
    $core.String? name,
    $core.int? targetdpus,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (targetdpus != null) result.targetdpus = targetdpus;
    return result;
  }

  UpdateCapacityReservationInput._();

  factory UpdateCapacityReservationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateCapacityReservationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateCapacityReservationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aI(367520745, _omitFieldNames ? '' : 'targetdpus')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCapacityReservationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCapacityReservationInput copyWith(
          void Function(UpdateCapacityReservationInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateCapacityReservationInput))
          as UpdateCapacityReservationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateCapacityReservationInput create() =>
      UpdateCapacityReservationInput._();
  @$core.override
  UpdateCapacityReservationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateCapacityReservationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateCapacityReservationInput>(create);
  static UpdateCapacityReservationInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(367520745)
  $core.int get targetdpus => $_getIZ(1);
  @$pb.TagNumber(367520745)
  set targetdpus($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(367520745)
  $core.bool hasTargetdpus() => $_has(1);
  @$pb.TagNumber(367520745)
  void clearTargetdpus() => $_clearField(367520745);
}

class UpdateCapacityReservationOutput extends $pb.GeneratedMessage {
  factory UpdateCapacityReservationOutput() => create();

  UpdateCapacityReservationOutput._();

  factory UpdateCapacityReservationOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateCapacityReservationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateCapacityReservationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCapacityReservationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCapacityReservationOutput copyWith(
          void Function(UpdateCapacityReservationOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateCapacityReservationOutput))
          as UpdateCapacityReservationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateCapacityReservationOutput create() =>
      UpdateCapacityReservationOutput._();
  @$core.override
  UpdateCapacityReservationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateCapacityReservationOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateCapacityReservationOutput>(
          create);
  static UpdateCapacityReservationOutput? _defaultInstance;
}

class UpdateDataCatalogInput extends $pb.GeneratedMessage {
  factory UpdateDataCatalogInput({
    $core.String? description,
    $core.String? name,
    DataCatalogType? type,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (parameters != null) result.parameters.addEntries(parameters);
    return result;
  }

  UpdateDataCatalogInput._();

  factory UpdateDataCatalogInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDataCatalogInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDataCatalogInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<DataCatalogType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: DataCatalogType.values)
    ..m<$core.String, $core.String>(
        494900218, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'UpdateDataCatalogInput.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('athena'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDataCatalogInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDataCatalogInput copyWith(
          void Function(UpdateDataCatalogInput) updates) =>
      super.copyWith((message) => updates(message as UpdateDataCatalogInput))
          as UpdateDataCatalogInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDataCatalogInput create() => UpdateDataCatalogInput._();
  @$core.override
  UpdateDataCatalogInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDataCatalogInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDataCatalogInput>(create);
  static UpdateDataCatalogInput? _defaultInstance;

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

  @$pb.TagNumber(290836590)
  DataCatalogType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(DataCatalogType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(494900218)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(3);
}

class UpdateDataCatalogOutput extends $pb.GeneratedMessage {
  factory UpdateDataCatalogOutput() => create();

  UpdateDataCatalogOutput._();

  factory UpdateDataCatalogOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDataCatalogOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDataCatalogOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDataCatalogOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDataCatalogOutput copyWith(
          void Function(UpdateDataCatalogOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateDataCatalogOutput))
          as UpdateDataCatalogOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDataCatalogOutput create() => UpdateDataCatalogOutput._();
  @$core.override
  UpdateDataCatalogOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDataCatalogOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDataCatalogOutput>(create);
  static UpdateDataCatalogOutput? _defaultInstance;
}

class UpdateNamedQueryInput extends $pb.GeneratedMessage {
  factory UpdateNamedQueryInput({
    $core.String? description,
    $core.String? name,
    $core.String? namedqueryid,
    $core.String? querystring,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (namedqueryid != null) result.namedqueryid = namedqueryid;
    if (querystring != null) result.querystring = querystring;
    return result;
  }

  UpdateNamedQueryInput._();

  factory UpdateNamedQueryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateNamedQueryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateNamedQueryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(330896872, _omitFieldNames ? '' : 'namedqueryid')
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNamedQueryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNamedQueryInput copyWith(
          void Function(UpdateNamedQueryInput) updates) =>
      super.copyWith((message) => updates(message as UpdateNamedQueryInput))
          as UpdateNamedQueryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateNamedQueryInput create() => UpdateNamedQueryInput._();
  @$core.override
  UpdateNamedQueryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateNamedQueryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateNamedQueryInput>(create);
  static UpdateNamedQueryInput? _defaultInstance;

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

  @$pb.TagNumber(330896872)
  $core.String get namedqueryid => $_getSZ(2);
  @$pb.TagNumber(330896872)
  set namedqueryid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(330896872)
  $core.bool hasNamedqueryid() => $_has(2);
  @$pb.TagNumber(330896872)
  void clearNamedqueryid() => $_clearField(330896872);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(3);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(3, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(3);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);
}

class UpdateNamedQueryOutput extends $pb.GeneratedMessage {
  factory UpdateNamedQueryOutput() => create();

  UpdateNamedQueryOutput._();

  factory UpdateNamedQueryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateNamedQueryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateNamedQueryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNamedQueryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNamedQueryOutput copyWith(
          void Function(UpdateNamedQueryOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateNamedQueryOutput))
          as UpdateNamedQueryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateNamedQueryOutput create() => UpdateNamedQueryOutput._();
  @$core.override
  UpdateNamedQueryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateNamedQueryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateNamedQueryOutput>(create);
  static UpdateNamedQueryOutput? _defaultInstance;
}

class UpdateNotebookInput extends $pb.GeneratedMessage {
  factory UpdateNotebookInput({
    $core.String? payload,
    $core.String? sessionid,
    $core.String? notebookid,
    NotebookType? type,
    $core.String? clientrequesttoken,
  }) {
    final result = create();
    if (payload != null) result.payload = payload;
    if (sessionid != null) result.sessionid = sessionid;
    if (notebookid != null) result.notebookid = notebookid;
    if (type != null) result.type = type;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    return result;
  }

  UpdateNotebookInput._();

  factory UpdateNotebookInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateNotebookInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateNotebookInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(6526790, _omitFieldNames ? '' : 'payload')
    ..aOS(20529723, _omitFieldNames ? '' : 'sessionid')
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..aE<NotebookType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: NotebookType.values)
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookInput copyWith(void Function(UpdateNotebookInput) updates) =>
      super.copyWith((message) => updates(message as UpdateNotebookInput))
          as UpdateNotebookInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateNotebookInput create() => UpdateNotebookInput._();
  @$core.override
  UpdateNotebookInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateNotebookInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateNotebookInput>(create);
  static UpdateNotebookInput? _defaultInstance;

  @$pb.TagNumber(6526790)
  $core.String get payload => $_getSZ(0);
  @$pb.TagNumber(6526790)
  set payload($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6526790)
  $core.bool hasPayload() => $_has(0);
  @$pb.TagNumber(6526790)
  void clearPayload() => $_clearField(6526790);

  @$pb.TagNumber(20529723)
  $core.String get sessionid => $_getSZ(1);
  @$pb.TagNumber(20529723)
  set sessionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(20529723)
  $core.bool hasSessionid() => $_has(1);
  @$pb.TagNumber(20529723)
  void clearSessionid() => $_clearField(20529723);

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(2);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(2);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);

  @$pb.TagNumber(290836590)
  NotebookType get type => $_getN(3);
  @$pb.TagNumber(290836590)
  set type(NotebookType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(3);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(4);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(4);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);
}

class UpdateNotebookMetadataInput extends $pb.GeneratedMessage {
  factory UpdateNotebookMetadataInput({
    $core.String? notebookid,
    $core.String? name,
    $core.String? clientrequesttoken,
  }) {
    final result = create();
    if (notebookid != null) result.notebookid = notebookid;
    if (name != null) result.name = name;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    return result;
  }

  UpdateNotebookMetadataInput._();

  factory UpdateNotebookMetadataInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateNotebookMetadataInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateNotebookMetadataInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(157637214, _omitFieldNames ? '' : 'notebookid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookMetadataInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookMetadataInput copyWith(
          void Function(UpdateNotebookMetadataInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateNotebookMetadataInput))
          as UpdateNotebookMetadataInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateNotebookMetadataInput create() =>
      UpdateNotebookMetadataInput._();
  @$core.override
  UpdateNotebookMetadataInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateNotebookMetadataInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateNotebookMetadataInput>(create);
  static UpdateNotebookMetadataInput? _defaultInstance;

  @$pb.TagNumber(157637214)
  $core.String get notebookid => $_getSZ(0);
  @$pb.TagNumber(157637214)
  set notebookid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157637214)
  $core.bool hasNotebookid() => $_has(0);
  @$pb.TagNumber(157637214)
  void clearNotebookid() => $_clearField(157637214);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(2);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(2);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);
}

class UpdateNotebookMetadataOutput extends $pb.GeneratedMessage {
  factory UpdateNotebookMetadataOutput() => create();

  UpdateNotebookMetadataOutput._();

  factory UpdateNotebookMetadataOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateNotebookMetadataOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateNotebookMetadataOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookMetadataOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookMetadataOutput copyWith(
          void Function(UpdateNotebookMetadataOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateNotebookMetadataOutput))
          as UpdateNotebookMetadataOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateNotebookMetadataOutput create() =>
      UpdateNotebookMetadataOutput._();
  @$core.override
  UpdateNotebookMetadataOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateNotebookMetadataOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateNotebookMetadataOutput>(create);
  static UpdateNotebookMetadataOutput? _defaultInstance;
}

class UpdateNotebookOutput extends $pb.GeneratedMessage {
  factory UpdateNotebookOutput() => create();

  UpdateNotebookOutput._();

  factory UpdateNotebookOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateNotebookOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateNotebookOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateNotebookOutput copyWith(void Function(UpdateNotebookOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateNotebookOutput))
          as UpdateNotebookOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateNotebookOutput create() => UpdateNotebookOutput._();
  @$core.override
  UpdateNotebookOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateNotebookOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateNotebookOutput>(create);
  static UpdateNotebookOutput? _defaultInstance;
}

class UpdatePreparedStatementInput extends $pb.GeneratedMessage {
  factory UpdatePreparedStatementInput({
    $core.String? statementname,
    $core.String? description,
    $core.String? querystatement,
    $core.String? workgroup,
  }) {
    final result = create();
    if (statementname != null) result.statementname = statementname;
    if (description != null) result.description = description;
    if (querystatement != null) result.querystatement = querystatement;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  UpdatePreparedStatementInput._();

  factory UpdatePreparedStatementInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdatePreparedStatementInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdatePreparedStatementInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(340852217, _omitFieldNames ? '' : 'querystatement')
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdatePreparedStatementInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdatePreparedStatementInput copyWith(
          void Function(UpdatePreparedStatementInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdatePreparedStatementInput))
          as UpdatePreparedStatementInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdatePreparedStatementInput create() =>
      UpdatePreparedStatementInput._();
  @$core.override
  UpdatePreparedStatementInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdatePreparedStatementInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdatePreparedStatementInput>(create);
  static UpdatePreparedStatementInput? _defaultInstance;

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(0);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(0);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(340852217)
  $core.String get querystatement => $_getSZ(2);
  @$pb.TagNumber(340852217)
  set querystatement($core.String value) => $_setString(2, value);
  @$pb.TagNumber(340852217)
  $core.bool hasQuerystatement() => $_has(2);
  @$pb.TagNumber(340852217)
  void clearQuerystatement() => $_clearField(340852217);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(3);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(3, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(3);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class UpdatePreparedStatementOutput extends $pb.GeneratedMessage {
  factory UpdatePreparedStatementOutput() => create();

  UpdatePreparedStatementOutput._();

  factory UpdatePreparedStatementOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdatePreparedStatementOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdatePreparedStatementOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdatePreparedStatementOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdatePreparedStatementOutput copyWith(
          void Function(UpdatePreparedStatementOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdatePreparedStatementOutput))
          as UpdatePreparedStatementOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdatePreparedStatementOutput create() =>
      UpdatePreparedStatementOutput._();
  @$core.override
  UpdatePreparedStatementOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdatePreparedStatementOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdatePreparedStatementOutput>(create);
  static UpdatePreparedStatementOutput? _defaultInstance;
}

class UpdateWorkGroupInput extends $pb.GeneratedMessage {
  factory UpdateWorkGroupInput({
    $core.String? description,
    WorkGroupConfigurationUpdates? configurationupdates,
    WorkGroupState? state,
    $core.String? workgroup,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (configurationupdates != null)
      result.configurationupdates = configurationupdates;
    if (state != null) result.state = state;
    if (workgroup != null) result.workgroup = workgroup;
    return result;
  }

  UpdateWorkGroupInput._();

  factory UpdateWorkGroupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateWorkGroupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateWorkGroupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOM<WorkGroupConfigurationUpdates>(
        133706738, _omitFieldNames ? '' : 'configurationupdates',
        subBuilder: WorkGroupConfigurationUpdates.create)
    ..aE<WorkGroupState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: WorkGroupState.values)
    ..aOS(505960068, _omitFieldNames ? '' : 'workgroup')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWorkGroupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWorkGroupInput copyWith(void Function(UpdateWorkGroupInput) updates) =>
      super.copyWith((message) => updates(message as UpdateWorkGroupInput))
          as UpdateWorkGroupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateWorkGroupInput create() => UpdateWorkGroupInput._();
  @$core.override
  UpdateWorkGroupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateWorkGroupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateWorkGroupInput>(create);
  static UpdateWorkGroupInput? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(133706738)
  WorkGroupConfigurationUpdates get configurationupdates => $_getN(1);
  @$pb.TagNumber(133706738)
  set configurationupdates(WorkGroupConfigurationUpdates value) =>
      $_setField(133706738, value);
  @$pb.TagNumber(133706738)
  $core.bool hasConfigurationupdates() => $_has(1);
  @$pb.TagNumber(133706738)
  void clearConfigurationupdates() => $_clearField(133706738);
  @$pb.TagNumber(133706738)
  WorkGroupConfigurationUpdates ensureConfigurationupdates() => $_ensure(1);

  @$pb.TagNumber(502047895)
  WorkGroupState get state => $_getN(2);
  @$pb.TagNumber(502047895)
  set state(WorkGroupState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(2);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(505960068)
  $core.String get workgroup => $_getSZ(3);
  @$pb.TagNumber(505960068)
  set workgroup($core.String value) => $_setString(3, value);
  @$pb.TagNumber(505960068)
  $core.bool hasWorkgroup() => $_has(3);
  @$pb.TagNumber(505960068)
  void clearWorkgroup() => $_clearField(505960068);
}

class UpdateWorkGroupOutput extends $pb.GeneratedMessage {
  factory UpdateWorkGroupOutput() => create();

  UpdateWorkGroupOutput._();

  factory UpdateWorkGroupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateWorkGroupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateWorkGroupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWorkGroupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWorkGroupOutput copyWith(
          void Function(UpdateWorkGroupOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateWorkGroupOutput))
          as UpdateWorkGroupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateWorkGroupOutput create() => UpdateWorkGroupOutput._();
  @$core.override
  UpdateWorkGroupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateWorkGroupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateWorkGroupOutput>(create);
  static UpdateWorkGroupOutput? _defaultInstance;
}

class WorkGroup extends $pb.GeneratedMessage {
  factory WorkGroup({
    $core.String? creationtime,
    $core.String? description,
    $core.String? name,
    $core.String? identitycenterapplicationarn,
    WorkGroupConfiguration? configuration,
    WorkGroupState? state,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (identitycenterapplicationarn != null)
      result.identitycenterapplicationarn = identitycenterapplicationarn;
    if (configuration != null) result.configuration = configuration;
    if (state != null) result.state = state;
    return result;
  }

  WorkGroup._();

  factory WorkGroup.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WorkGroup.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WorkGroup',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338293532, _omitFieldNames ? '' : 'identitycenterapplicationarn')
    ..aOM<WorkGroupConfiguration>(
        442426458, _omitFieldNames ? '' : 'configuration',
        subBuilder: WorkGroupConfiguration.create)
    ..aE<WorkGroupState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: WorkGroupState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroup clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroup copyWith(void Function(WorkGroup) updates) =>
      super.copyWith((message) => updates(message as WorkGroup)) as WorkGroup;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WorkGroup create() => WorkGroup._();
  @$core.override
  WorkGroup createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WorkGroup getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<WorkGroup>(create);
  static WorkGroup? _defaultInstance;

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

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338293532)
  $core.String get identitycenterapplicationarn => $_getSZ(3);
  @$pb.TagNumber(338293532)
  set identitycenterapplicationarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(338293532)
  $core.bool hasIdentitycenterapplicationarn() => $_has(3);
  @$pb.TagNumber(338293532)
  void clearIdentitycenterapplicationarn() => $_clearField(338293532);

  @$pb.TagNumber(442426458)
  WorkGroupConfiguration get configuration => $_getN(4);
  @$pb.TagNumber(442426458)
  set configuration(WorkGroupConfiguration value) =>
      $_setField(442426458, value);
  @$pb.TagNumber(442426458)
  $core.bool hasConfiguration() => $_has(4);
  @$pb.TagNumber(442426458)
  void clearConfiguration() => $_clearField(442426458);
  @$pb.TagNumber(442426458)
  WorkGroupConfiguration ensureConfiguration() => $_ensure(4);

  @$pb.TagNumber(502047895)
  WorkGroupState get state => $_getN(5);
  @$pb.TagNumber(502047895)
  set state(WorkGroupState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class WorkGroupConfiguration extends $pb.GeneratedMessage {
  factory WorkGroupConfiguration({
    EngineVersion? engineversion,
    $core.bool? enforceworkgroupconfiguration,
    ManagedQueryResultsConfiguration? managedqueryresultsconfiguration,
    CustomerContentEncryptionConfiguration?
        customercontentencryptionconfiguration,
    ResultConfiguration? resultconfiguration,
    IdentityCenterConfiguration? identitycenterconfiguration,
    $core.bool? enableminimumencryptionconfiguration,
    $core.String? executionrole,
    QueryResultsS3AccessGrantsConfiguration?
        queryresultss3accessgrantsconfiguration,
    $fixnum.Int64? bytesscannedcutoffperquery,
    EngineConfiguration? engineconfiguration,
    MonitoringConfiguration? monitoringconfiguration,
    $core.String? additionalconfiguration,
    $core.bool? requesterpaysenabled,
    $core.bool? publishcloudwatchmetricsenabled,
  }) {
    final result = create();
    if (engineversion != null) result.engineversion = engineversion;
    if (enforceworkgroupconfiguration != null)
      result.enforceworkgroupconfiguration = enforceworkgroupconfiguration;
    if (managedqueryresultsconfiguration != null)
      result.managedqueryresultsconfiguration =
          managedqueryresultsconfiguration;
    if (customercontentencryptionconfiguration != null)
      result.customercontentencryptionconfiguration =
          customercontentencryptionconfiguration;
    if (resultconfiguration != null)
      result.resultconfiguration = resultconfiguration;
    if (identitycenterconfiguration != null)
      result.identitycenterconfiguration = identitycenterconfiguration;
    if (enableminimumencryptionconfiguration != null)
      result.enableminimumencryptionconfiguration =
          enableminimumencryptionconfiguration;
    if (executionrole != null) result.executionrole = executionrole;
    if (queryresultss3accessgrantsconfiguration != null)
      result.queryresultss3accessgrantsconfiguration =
          queryresultss3accessgrantsconfiguration;
    if (bytesscannedcutoffperquery != null)
      result.bytesscannedcutoffperquery = bytesscannedcutoffperquery;
    if (engineconfiguration != null)
      result.engineconfiguration = engineconfiguration;
    if (monitoringconfiguration != null)
      result.monitoringconfiguration = monitoringconfiguration;
    if (additionalconfiguration != null)
      result.additionalconfiguration = additionalconfiguration;
    if (requesterpaysenabled != null)
      result.requesterpaysenabled = requesterpaysenabled;
    if (publishcloudwatchmetricsenabled != null)
      result.publishcloudwatchmetricsenabled = publishcloudwatchmetricsenabled;
    return result;
  }

  WorkGroupConfiguration._();

  factory WorkGroupConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WorkGroupConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WorkGroupConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<EngineVersion>(44953462, _omitFieldNames ? '' : 'engineversion',
        subBuilder: EngineVersion.create)
    ..aOB(152602624, _omitFieldNames ? '' : 'enforceworkgroupconfiguration')
    ..aOM<ManagedQueryResultsConfiguration>(
        159215683, _omitFieldNames ? '' : 'managedqueryresultsconfiguration',
        subBuilder: ManagedQueryResultsConfiguration.create)
    ..aOM<CustomerContentEncryptionConfiguration>(165213900,
        _omitFieldNames ? '' : 'customercontentencryptionconfiguration',
        subBuilder: CustomerContentEncryptionConfiguration.create)
    ..aOM<ResultConfiguration>(
        183031503, _omitFieldNames ? '' : 'resultconfiguration',
        subBuilder: ResultConfiguration.create)
    ..aOM<IdentityCenterConfiguration>(
        236974917, _omitFieldNames ? '' : 'identitycenterconfiguration',
        subBuilder: IdentityCenterConfiguration.create)
    ..aOB(238637616,
        _omitFieldNames ? '' : 'enableminimumencryptionconfiguration')
    ..aOS(253307658, _omitFieldNames ? '' : 'executionrole')
    ..aOM<QueryResultsS3AccessGrantsConfiguration>(264403639,
        _omitFieldNames ? '' : 'queryresultss3accessgrantsconfiguration',
        subBuilder: QueryResultsS3AccessGrantsConfiguration.create)
    ..aInt64(265761289, _omitFieldNames ? '' : 'bytesscannedcutoffperquery')
    ..aOM<EngineConfiguration>(
        341629412, _omitFieldNames ? '' : 'engineconfiguration',
        subBuilder: EngineConfiguration.create)
    ..aOM<MonitoringConfiguration>(
        364891928, _omitFieldNames ? '' : 'monitoringconfiguration',
        subBuilder: MonitoringConfiguration.create)
    ..aOS(389584375, _omitFieldNames ? '' : 'additionalconfiguration')
    ..aOB(448832780, _omitFieldNames ? '' : 'requesterpaysenabled')
    ..aOB(493112579, _omitFieldNames ? '' : 'publishcloudwatchmetricsenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroupConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroupConfiguration copyWith(
          void Function(WorkGroupConfiguration) updates) =>
      super.copyWith((message) => updates(message as WorkGroupConfiguration))
          as WorkGroupConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WorkGroupConfiguration create() => WorkGroupConfiguration._();
  @$core.override
  WorkGroupConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WorkGroupConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WorkGroupConfiguration>(create);
  static WorkGroupConfiguration? _defaultInstance;

  @$pb.TagNumber(44953462)
  EngineVersion get engineversion => $_getN(0);
  @$pb.TagNumber(44953462)
  set engineversion(EngineVersion value) => $_setField(44953462, value);
  @$pb.TagNumber(44953462)
  $core.bool hasEngineversion() => $_has(0);
  @$pb.TagNumber(44953462)
  void clearEngineversion() => $_clearField(44953462);
  @$pb.TagNumber(44953462)
  EngineVersion ensureEngineversion() => $_ensure(0);

  @$pb.TagNumber(152602624)
  $core.bool get enforceworkgroupconfiguration => $_getBF(1);
  @$pb.TagNumber(152602624)
  set enforceworkgroupconfiguration($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(152602624)
  $core.bool hasEnforceworkgroupconfiguration() => $_has(1);
  @$pb.TagNumber(152602624)
  void clearEnforceworkgroupconfiguration() => $_clearField(152602624);

  @$pb.TagNumber(159215683)
  ManagedQueryResultsConfiguration get managedqueryresultsconfiguration =>
      $_getN(2);
  @$pb.TagNumber(159215683)
  set managedqueryresultsconfiguration(
          ManagedQueryResultsConfiguration value) =>
      $_setField(159215683, value);
  @$pb.TagNumber(159215683)
  $core.bool hasManagedqueryresultsconfiguration() => $_has(2);
  @$pb.TagNumber(159215683)
  void clearManagedqueryresultsconfiguration() => $_clearField(159215683);
  @$pb.TagNumber(159215683)
  ManagedQueryResultsConfiguration ensureManagedqueryresultsconfiguration() =>
      $_ensure(2);

  @$pb.TagNumber(165213900)
  CustomerContentEncryptionConfiguration
      get customercontentencryptionconfiguration => $_getN(3);
  @$pb.TagNumber(165213900)
  set customercontentencryptionconfiguration(
          CustomerContentEncryptionConfiguration value) =>
      $_setField(165213900, value);
  @$pb.TagNumber(165213900)
  $core.bool hasCustomercontentencryptionconfiguration() => $_has(3);
  @$pb.TagNumber(165213900)
  void clearCustomercontentencryptionconfiguration() => $_clearField(165213900);
  @$pb.TagNumber(165213900)
  CustomerContentEncryptionConfiguration
      ensureCustomercontentencryptionconfiguration() => $_ensure(3);

  @$pb.TagNumber(183031503)
  ResultConfiguration get resultconfiguration => $_getN(4);
  @$pb.TagNumber(183031503)
  set resultconfiguration(ResultConfiguration value) =>
      $_setField(183031503, value);
  @$pb.TagNumber(183031503)
  $core.bool hasResultconfiguration() => $_has(4);
  @$pb.TagNumber(183031503)
  void clearResultconfiguration() => $_clearField(183031503);
  @$pb.TagNumber(183031503)
  ResultConfiguration ensureResultconfiguration() => $_ensure(4);

  @$pb.TagNumber(236974917)
  IdentityCenterConfiguration get identitycenterconfiguration => $_getN(5);
  @$pb.TagNumber(236974917)
  set identitycenterconfiguration(IdentityCenterConfiguration value) =>
      $_setField(236974917, value);
  @$pb.TagNumber(236974917)
  $core.bool hasIdentitycenterconfiguration() => $_has(5);
  @$pb.TagNumber(236974917)
  void clearIdentitycenterconfiguration() => $_clearField(236974917);
  @$pb.TagNumber(236974917)
  IdentityCenterConfiguration ensureIdentitycenterconfiguration() =>
      $_ensure(5);

  @$pb.TagNumber(238637616)
  $core.bool get enableminimumencryptionconfiguration => $_getBF(6);
  @$pb.TagNumber(238637616)
  set enableminimumencryptionconfiguration($core.bool value) =>
      $_setBool(6, value);
  @$pb.TagNumber(238637616)
  $core.bool hasEnableminimumencryptionconfiguration() => $_has(6);
  @$pb.TagNumber(238637616)
  void clearEnableminimumencryptionconfiguration() => $_clearField(238637616);

  @$pb.TagNumber(253307658)
  $core.String get executionrole => $_getSZ(7);
  @$pb.TagNumber(253307658)
  set executionrole($core.String value) => $_setString(7, value);
  @$pb.TagNumber(253307658)
  $core.bool hasExecutionrole() => $_has(7);
  @$pb.TagNumber(253307658)
  void clearExecutionrole() => $_clearField(253307658);

  @$pb.TagNumber(264403639)
  QueryResultsS3AccessGrantsConfiguration
      get queryresultss3accessgrantsconfiguration => $_getN(8);
  @$pb.TagNumber(264403639)
  set queryresultss3accessgrantsconfiguration(
          QueryResultsS3AccessGrantsConfiguration value) =>
      $_setField(264403639, value);
  @$pb.TagNumber(264403639)
  $core.bool hasQueryresultss3accessgrantsconfiguration() => $_has(8);
  @$pb.TagNumber(264403639)
  void clearQueryresultss3accessgrantsconfiguration() =>
      $_clearField(264403639);
  @$pb.TagNumber(264403639)
  QueryResultsS3AccessGrantsConfiguration
      ensureQueryresultss3accessgrantsconfiguration() => $_ensure(8);

  @$pb.TagNumber(265761289)
  $fixnum.Int64 get bytesscannedcutoffperquery => $_getI64(9);
  @$pb.TagNumber(265761289)
  set bytesscannedcutoffperquery($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(265761289)
  $core.bool hasBytesscannedcutoffperquery() => $_has(9);
  @$pb.TagNumber(265761289)
  void clearBytesscannedcutoffperquery() => $_clearField(265761289);

  @$pb.TagNumber(341629412)
  EngineConfiguration get engineconfiguration => $_getN(10);
  @$pb.TagNumber(341629412)
  set engineconfiguration(EngineConfiguration value) =>
      $_setField(341629412, value);
  @$pb.TagNumber(341629412)
  $core.bool hasEngineconfiguration() => $_has(10);
  @$pb.TagNumber(341629412)
  void clearEngineconfiguration() => $_clearField(341629412);
  @$pb.TagNumber(341629412)
  EngineConfiguration ensureEngineconfiguration() => $_ensure(10);

  @$pb.TagNumber(364891928)
  MonitoringConfiguration get monitoringconfiguration => $_getN(11);
  @$pb.TagNumber(364891928)
  set monitoringconfiguration(MonitoringConfiguration value) =>
      $_setField(364891928, value);
  @$pb.TagNumber(364891928)
  $core.bool hasMonitoringconfiguration() => $_has(11);
  @$pb.TagNumber(364891928)
  void clearMonitoringconfiguration() => $_clearField(364891928);
  @$pb.TagNumber(364891928)
  MonitoringConfiguration ensureMonitoringconfiguration() => $_ensure(11);

  @$pb.TagNumber(389584375)
  $core.String get additionalconfiguration => $_getSZ(12);
  @$pb.TagNumber(389584375)
  set additionalconfiguration($core.String value) => $_setString(12, value);
  @$pb.TagNumber(389584375)
  $core.bool hasAdditionalconfiguration() => $_has(12);
  @$pb.TagNumber(389584375)
  void clearAdditionalconfiguration() => $_clearField(389584375);

  @$pb.TagNumber(448832780)
  $core.bool get requesterpaysenabled => $_getBF(13);
  @$pb.TagNumber(448832780)
  set requesterpaysenabled($core.bool value) => $_setBool(13, value);
  @$pb.TagNumber(448832780)
  $core.bool hasRequesterpaysenabled() => $_has(13);
  @$pb.TagNumber(448832780)
  void clearRequesterpaysenabled() => $_clearField(448832780);

  @$pb.TagNumber(493112579)
  $core.bool get publishcloudwatchmetricsenabled => $_getBF(14);
  @$pb.TagNumber(493112579)
  set publishcloudwatchmetricsenabled($core.bool value) => $_setBool(14, value);
  @$pb.TagNumber(493112579)
  $core.bool hasPublishcloudwatchmetricsenabled() => $_has(14);
  @$pb.TagNumber(493112579)
  void clearPublishcloudwatchmetricsenabled() => $_clearField(493112579);
}

class WorkGroupConfigurationUpdates extends $pb.GeneratedMessage {
  factory WorkGroupConfigurationUpdates({
    EngineVersion? engineversion,
    ResultConfigurationUpdates? resultconfigurationupdates,
    $core.bool? removebytesscannedcutoffperquery,
    $core.bool? enforceworkgroupconfiguration,
    CustomerContentEncryptionConfiguration?
        customercontentencryptionconfiguration,
    ManagedQueryResultsConfigurationUpdates?
        managedqueryresultsconfigurationupdates,
    $core.bool? enableminimumencryptionconfiguration,
    $core.String? executionrole,
    QueryResultsS3AccessGrantsConfiguration?
        queryresultss3accessgrantsconfiguration,
    $fixnum.Int64? bytesscannedcutoffperquery,
    EngineConfiguration? engineconfiguration,
    MonitoringConfiguration? monitoringconfiguration,
    $core.String? additionalconfiguration,
    $core.bool? removecustomercontentencryptionconfiguration,
    $core.bool? requesterpaysenabled,
    $core.bool? publishcloudwatchmetricsenabled,
  }) {
    final result = create();
    if (engineversion != null) result.engineversion = engineversion;
    if (resultconfigurationupdates != null)
      result.resultconfigurationupdates = resultconfigurationupdates;
    if (removebytesscannedcutoffperquery != null)
      result.removebytesscannedcutoffperquery =
          removebytesscannedcutoffperquery;
    if (enforceworkgroupconfiguration != null)
      result.enforceworkgroupconfiguration = enforceworkgroupconfiguration;
    if (customercontentencryptionconfiguration != null)
      result.customercontentencryptionconfiguration =
          customercontentencryptionconfiguration;
    if (managedqueryresultsconfigurationupdates != null)
      result.managedqueryresultsconfigurationupdates =
          managedqueryresultsconfigurationupdates;
    if (enableminimumencryptionconfiguration != null)
      result.enableminimumencryptionconfiguration =
          enableminimumencryptionconfiguration;
    if (executionrole != null) result.executionrole = executionrole;
    if (queryresultss3accessgrantsconfiguration != null)
      result.queryresultss3accessgrantsconfiguration =
          queryresultss3accessgrantsconfiguration;
    if (bytesscannedcutoffperquery != null)
      result.bytesscannedcutoffperquery = bytesscannedcutoffperquery;
    if (engineconfiguration != null)
      result.engineconfiguration = engineconfiguration;
    if (monitoringconfiguration != null)
      result.monitoringconfiguration = monitoringconfiguration;
    if (additionalconfiguration != null)
      result.additionalconfiguration = additionalconfiguration;
    if (removecustomercontentencryptionconfiguration != null)
      result.removecustomercontentencryptionconfiguration =
          removecustomercontentencryptionconfiguration;
    if (requesterpaysenabled != null)
      result.requesterpaysenabled = requesterpaysenabled;
    if (publishcloudwatchmetricsenabled != null)
      result.publishcloudwatchmetricsenabled = publishcloudwatchmetricsenabled;
    return result;
  }

  WorkGroupConfigurationUpdates._();

  factory WorkGroupConfigurationUpdates.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WorkGroupConfigurationUpdates.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WorkGroupConfigurationUpdates',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<EngineVersion>(44953462, _omitFieldNames ? '' : 'engineversion',
        subBuilder: EngineVersion.create)
    ..aOM<ResultConfigurationUpdates>(
        97856045, _omitFieldNames ? '' : 'resultconfigurationupdates',
        subBuilder: ResultConfigurationUpdates.create)
    ..aOB(130456497, _omitFieldNames ? '' : 'removebytesscannedcutoffperquery')
    ..aOB(152602624, _omitFieldNames ? '' : 'enforceworkgroupconfiguration')
    ..aOM<CustomerContentEncryptionConfiguration>(165213900,
        _omitFieldNames ? '' : 'customercontentencryptionconfiguration',
        subBuilder: CustomerContentEncryptionConfiguration.create)
    ..aOM<ManagedQueryResultsConfigurationUpdates>(211421145,
        _omitFieldNames ? '' : 'managedqueryresultsconfigurationupdates',
        subBuilder: ManagedQueryResultsConfigurationUpdates.create)
    ..aOB(238637616,
        _omitFieldNames ? '' : 'enableminimumencryptionconfiguration')
    ..aOS(253307658, _omitFieldNames ? '' : 'executionrole')
    ..aOM<QueryResultsS3AccessGrantsConfiguration>(264403639,
        _omitFieldNames ? '' : 'queryresultss3accessgrantsconfiguration',
        subBuilder: QueryResultsS3AccessGrantsConfiguration.create)
    ..aInt64(265761289, _omitFieldNames ? '' : 'bytesscannedcutoffperquery')
    ..aOM<EngineConfiguration>(
        341629412, _omitFieldNames ? '' : 'engineconfiguration',
        subBuilder: EngineConfiguration.create)
    ..aOM<MonitoringConfiguration>(
        364891928, _omitFieldNames ? '' : 'monitoringconfiguration',
        subBuilder: MonitoringConfiguration.create)
    ..aOS(389584375, _omitFieldNames ? '' : 'additionalconfiguration')
    ..aOB(418402884,
        _omitFieldNames ? '' : 'removecustomercontentencryptionconfiguration')
    ..aOB(448832780, _omitFieldNames ? '' : 'requesterpaysenabled')
    ..aOB(493112579, _omitFieldNames ? '' : 'publishcloudwatchmetricsenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroupConfigurationUpdates clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroupConfigurationUpdates copyWith(
          void Function(WorkGroupConfigurationUpdates) updates) =>
      super.copyWith(
              (message) => updates(message as WorkGroupConfigurationUpdates))
          as WorkGroupConfigurationUpdates;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WorkGroupConfigurationUpdates create() =>
      WorkGroupConfigurationUpdates._();
  @$core.override
  WorkGroupConfigurationUpdates createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WorkGroupConfigurationUpdates getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WorkGroupConfigurationUpdates>(create);
  static WorkGroupConfigurationUpdates? _defaultInstance;

  @$pb.TagNumber(44953462)
  EngineVersion get engineversion => $_getN(0);
  @$pb.TagNumber(44953462)
  set engineversion(EngineVersion value) => $_setField(44953462, value);
  @$pb.TagNumber(44953462)
  $core.bool hasEngineversion() => $_has(0);
  @$pb.TagNumber(44953462)
  void clearEngineversion() => $_clearField(44953462);
  @$pb.TagNumber(44953462)
  EngineVersion ensureEngineversion() => $_ensure(0);

  @$pb.TagNumber(97856045)
  ResultConfigurationUpdates get resultconfigurationupdates => $_getN(1);
  @$pb.TagNumber(97856045)
  set resultconfigurationupdates(ResultConfigurationUpdates value) =>
      $_setField(97856045, value);
  @$pb.TagNumber(97856045)
  $core.bool hasResultconfigurationupdates() => $_has(1);
  @$pb.TagNumber(97856045)
  void clearResultconfigurationupdates() => $_clearField(97856045);
  @$pb.TagNumber(97856045)
  ResultConfigurationUpdates ensureResultconfigurationupdates() => $_ensure(1);

  @$pb.TagNumber(130456497)
  $core.bool get removebytesscannedcutoffperquery => $_getBF(2);
  @$pb.TagNumber(130456497)
  set removebytesscannedcutoffperquery($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(130456497)
  $core.bool hasRemovebytesscannedcutoffperquery() => $_has(2);
  @$pb.TagNumber(130456497)
  void clearRemovebytesscannedcutoffperquery() => $_clearField(130456497);

  @$pb.TagNumber(152602624)
  $core.bool get enforceworkgroupconfiguration => $_getBF(3);
  @$pb.TagNumber(152602624)
  set enforceworkgroupconfiguration($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(152602624)
  $core.bool hasEnforceworkgroupconfiguration() => $_has(3);
  @$pb.TagNumber(152602624)
  void clearEnforceworkgroupconfiguration() => $_clearField(152602624);

  @$pb.TagNumber(165213900)
  CustomerContentEncryptionConfiguration
      get customercontentencryptionconfiguration => $_getN(4);
  @$pb.TagNumber(165213900)
  set customercontentencryptionconfiguration(
          CustomerContentEncryptionConfiguration value) =>
      $_setField(165213900, value);
  @$pb.TagNumber(165213900)
  $core.bool hasCustomercontentencryptionconfiguration() => $_has(4);
  @$pb.TagNumber(165213900)
  void clearCustomercontentencryptionconfiguration() => $_clearField(165213900);
  @$pb.TagNumber(165213900)
  CustomerContentEncryptionConfiguration
      ensureCustomercontentencryptionconfiguration() => $_ensure(4);

  @$pb.TagNumber(211421145)
  ManagedQueryResultsConfigurationUpdates
      get managedqueryresultsconfigurationupdates => $_getN(5);
  @$pb.TagNumber(211421145)
  set managedqueryresultsconfigurationupdates(
          ManagedQueryResultsConfigurationUpdates value) =>
      $_setField(211421145, value);
  @$pb.TagNumber(211421145)
  $core.bool hasManagedqueryresultsconfigurationupdates() => $_has(5);
  @$pb.TagNumber(211421145)
  void clearManagedqueryresultsconfigurationupdates() =>
      $_clearField(211421145);
  @$pb.TagNumber(211421145)
  ManagedQueryResultsConfigurationUpdates
      ensureManagedqueryresultsconfigurationupdates() => $_ensure(5);

  @$pb.TagNumber(238637616)
  $core.bool get enableminimumencryptionconfiguration => $_getBF(6);
  @$pb.TagNumber(238637616)
  set enableminimumencryptionconfiguration($core.bool value) =>
      $_setBool(6, value);
  @$pb.TagNumber(238637616)
  $core.bool hasEnableminimumencryptionconfiguration() => $_has(6);
  @$pb.TagNumber(238637616)
  void clearEnableminimumencryptionconfiguration() => $_clearField(238637616);

  @$pb.TagNumber(253307658)
  $core.String get executionrole => $_getSZ(7);
  @$pb.TagNumber(253307658)
  set executionrole($core.String value) => $_setString(7, value);
  @$pb.TagNumber(253307658)
  $core.bool hasExecutionrole() => $_has(7);
  @$pb.TagNumber(253307658)
  void clearExecutionrole() => $_clearField(253307658);

  @$pb.TagNumber(264403639)
  QueryResultsS3AccessGrantsConfiguration
      get queryresultss3accessgrantsconfiguration => $_getN(8);
  @$pb.TagNumber(264403639)
  set queryresultss3accessgrantsconfiguration(
          QueryResultsS3AccessGrantsConfiguration value) =>
      $_setField(264403639, value);
  @$pb.TagNumber(264403639)
  $core.bool hasQueryresultss3accessgrantsconfiguration() => $_has(8);
  @$pb.TagNumber(264403639)
  void clearQueryresultss3accessgrantsconfiguration() =>
      $_clearField(264403639);
  @$pb.TagNumber(264403639)
  QueryResultsS3AccessGrantsConfiguration
      ensureQueryresultss3accessgrantsconfiguration() => $_ensure(8);

  @$pb.TagNumber(265761289)
  $fixnum.Int64 get bytesscannedcutoffperquery => $_getI64(9);
  @$pb.TagNumber(265761289)
  set bytesscannedcutoffperquery($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(265761289)
  $core.bool hasBytesscannedcutoffperquery() => $_has(9);
  @$pb.TagNumber(265761289)
  void clearBytesscannedcutoffperquery() => $_clearField(265761289);

  @$pb.TagNumber(341629412)
  EngineConfiguration get engineconfiguration => $_getN(10);
  @$pb.TagNumber(341629412)
  set engineconfiguration(EngineConfiguration value) =>
      $_setField(341629412, value);
  @$pb.TagNumber(341629412)
  $core.bool hasEngineconfiguration() => $_has(10);
  @$pb.TagNumber(341629412)
  void clearEngineconfiguration() => $_clearField(341629412);
  @$pb.TagNumber(341629412)
  EngineConfiguration ensureEngineconfiguration() => $_ensure(10);

  @$pb.TagNumber(364891928)
  MonitoringConfiguration get monitoringconfiguration => $_getN(11);
  @$pb.TagNumber(364891928)
  set monitoringconfiguration(MonitoringConfiguration value) =>
      $_setField(364891928, value);
  @$pb.TagNumber(364891928)
  $core.bool hasMonitoringconfiguration() => $_has(11);
  @$pb.TagNumber(364891928)
  void clearMonitoringconfiguration() => $_clearField(364891928);
  @$pb.TagNumber(364891928)
  MonitoringConfiguration ensureMonitoringconfiguration() => $_ensure(11);

  @$pb.TagNumber(389584375)
  $core.String get additionalconfiguration => $_getSZ(12);
  @$pb.TagNumber(389584375)
  set additionalconfiguration($core.String value) => $_setString(12, value);
  @$pb.TagNumber(389584375)
  $core.bool hasAdditionalconfiguration() => $_has(12);
  @$pb.TagNumber(389584375)
  void clearAdditionalconfiguration() => $_clearField(389584375);

  @$pb.TagNumber(418402884)
  $core.bool get removecustomercontentencryptionconfiguration => $_getBF(13);
  @$pb.TagNumber(418402884)
  set removecustomercontentencryptionconfiguration($core.bool value) =>
      $_setBool(13, value);
  @$pb.TagNumber(418402884)
  $core.bool hasRemovecustomercontentencryptionconfiguration() => $_has(13);
  @$pb.TagNumber(418402884)
  void clearRemovecustomercontentencryptionconfiguration() =>
      $_clearField(418402884);

  @$pb.TagNumber(448832780)
  $core.bool get requesterpaysenabled => $_getBF(14);
  @$pb.TagNumber(448832780)
  set requesterpaysenabled($core.bool value) => $_setBool(14, value);
  @$pb.TagNumber(448832780)
  $core.bool hasRequesterpaysenabled() => $_has(14);
  @$pb.TagNumber(448832780)
  void clearRequesterpaysenabled() => $_clearField(448832780);

  @$pb.TagNumber(493112579)
  $core.bool get publishcloudwatchmetricsenabled => $_getBF(15);
  @$pb.TagNumber(493112579)
  set publishcloudwatchmetricsenabled($core.bool value) => $_setBool(15, value);
  @$pb.TagNumber(493112579)
  $core.bool hasPublishcloudwatchmetricsenabled() => $_has(15);
  @$pb.TagNumber(493112579)
  void clearPublishcloudwatchmetricsenabled() => $_clearField(493112579);
}

class WorkGroupSummary extends $pb.GeneratedMessage {
  factory WorkGroupSummary({
    EngineVersion? engineversion,
    $core.String? creationtime,
    $core.String? description,
    $core.String? name,
    $core.String? identitycenterapplicationarn,
    WorkGroupState? state,
  }) {
    final result = create();
    if (engineversion != null) result.engineversion = engineversion;
    if (creationtime != null) result.creationtime = creationtime;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (identitycenterapplicationarn != null)
      result.identitycenterapplicationarn = identitycenterapplicationarn;
    if (state != null) result.state = state;
    return result;
  }

  WorkGroupSummary._();

  factory WorkGroupSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WorkGroupSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WorkGroupSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'athena'),
      createEmptyInstance: create)
    ..aOM<EngineVersion>(44953462, _omitFieldNames ? '' : 'engineversion',
        subBuilder: EngineVersion.create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338293532, _omitFieldNames ? '' : 'identitycenterapplicationarn')
    ..aE<WorkGroupState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: WorkGroupState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroupSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WorkGroupSummary copyWith(void Function(WorkGroupSummary) updates) =>
      super.copyWith((message) => updates(message as WorkGroupSummary))
          as WorkGroupSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WorkGroupSummary create() => WorkGroupSummary._();
  @$core.override
  WorkGroupSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WorkGroupSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WorkGroupSummary>(create);
  static WorkGroupSummary? _defaultInstance;

  @$pb.TagNumber(44953462)
  EngineVersion get engineversion => $_getN(0);
  @$pb.TagNumber(44953462)
  set engineversion(EngineVersion value) => $_setField(44953462, value);
  @$pb.TagNumber(44953462)
  $core.bool hasEngineversion() => $_has(0);
  @$pb.TagNumber(44953462)
  void clearEngineversion() => $_clearField(44953462);
  @$pb.TagNumber(44953462)
  EngineVersion ensureEngineversion() => $_ensure(0);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

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

  @$pb.TagNumber(338293532)
  $core.String get identitycenterapplicationarn => $_getSZ(4);
  @$pb.TagNumber(338293532)
  set identitycenterapplicationarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(338293532)
  $core.bool hasIdentitycenterapplicationarn() => $_has(4);
  @$pb.TagNumber(338293532)
  void clearIdentitycenterapplicationarn() => $_clearField(338293532);

  @$pb.TagNumber(502047895)
  WorkGroupState get state => $_getN(5);
  @$pb.TagNumber(502047895)
  set state(WorkGroupState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
