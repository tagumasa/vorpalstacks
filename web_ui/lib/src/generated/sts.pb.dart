// This is a generated file - do not edit.
//
// Generated from sts.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class AssumeRoleRequest extends $pb.GeneratedMessage {
  factory AssumeRoleRequest({
    $core.Iterable<PolicyDescriptorType>? policyarns,
    $core.Iterable<ProvidedContext>? providedcontexts,
    $core.String? externalid,
    $core.String? tokencode,
    $core.String? rolesessionname,
    $core.String? rolearn,
    $core.Iterable<Tag>? tags,
    $core.String? serialnumber,
    $core.int? durationseconds,
    $core.Iterable<$core.String>? transitivetagkeys,
    $core.String? sourceidentity,
    $core.String? policy,
  }) {
    final result = create();
    if (policyarns != null) result.policyarns.addAll(policyarns);
    if (providedcontexts != null)
      result.providedcontexts.addAll(providedcontexts);
    if (externalid != null) result.externalid = externalid;
    if (tokencode != null) result.tokencode = tokencode;
    if (rolesessionname != null) result.rolesessionname = rolesessionname;
    if (rolearn != null) result.rolearn = rolearn;
    if (tags != null) result.tags.addAll(tags);
    if (serialnumber != null) result.serialnumber = serialnumber;
    if (durationseconds != null) result.durationseconds = durationseconds;
    if (transitivetagkeys != null)
      result.transitivetagkeys.addAll(transitivetagkeys);
    if (sourceidentity != null) result.sourceidentity = sourceidentity;
    if (policy != null) result.policy = policy;
    return result;
  }

  AssumeRoleRequest._();

  factory AssumeRoleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRoleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRoleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..pPM<PolicyDescriptorType>(183785508, _omitFieldNames ? '' : 'policyarns',
        subBuilder: PolicyDescriptorType.create)
    ..pPM<ProvidedContext>(228510151, _omitFieldNames ? '' : 'providedcontexts',
        subBuilder: ProvidedContext.create)
    ..aOS(271401992, _omitFieldNames ? '' : 'externalid')
    ..aOS(300671456, _omitFieldNames ? '' : 'tokencode')
    ..aOS(315098849, _omitFieldNames ? '' : 'rolesessionname')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(418274661, _omitFieldNames ? '' : 'serialnumber')
    ..aI(451873635, _omitFieldNames ? '' : 'durationseconds')
    ..pPS(452608727, _omitFieldNames ? '' : 'transitivetagkeys')
    ..aOS(466635355, _omitFieldNames ? '' : 'sourceidentity')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleRequest copyWith(void Function(AssumeRoleRequest) updates) =>
      super.copyWith((message) => updates(message as AssumeRoleRequest))
          as AssumeRoleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRoleRequest create() => AssumeRoleRequest._();
  @$core.override
  AssumeRoleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRoleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRoleRequest>(create);
  static AssumeRoleRequest? _defaultInstance;

  @$pb.TagNumber(183785508)
  $pb.PbList<PolicyDescriptorType> get policyarns => $_getList(0);

  @$pb.TagNumber(228510151)
  $pb.PbList<ProvidedContext> get providedcontexts => $_getList(1);

  @$pb.TagNumber(271401992)
  $core.String get externalid => $_getSZ(2);
  @$pb.TagNumber(271401992)
  set externalid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(271401992)
  $core.bool hasExternalid() => $_has(2);
  @$pb.TagNumber(271401992)
  void clearExternalid() => $_clearField(271401992);

  @$pb.TagNumber(300671456)
  $core.String get tokencode => $_getSZ(3);
  @$pb.TagNumber(300671456)
  set tokencode($core.String value) => $_setString(3, value);
  @$pb.TagNumber(300671456)
  $core.bool hasTokencode() => $_has(3);
  @$pb.TagNumber(300671456)
  void clearTokencode() => $_clearField(300671456);

  @$pb.TagNumber(315098849)
  $core.String get rolesessionname => $_getSZ(4);
  @$pb.TagNumber(315098849)
  set rolesessionname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(315098849)
  $core.bool hasRolesessionname() => $_has(4);
  @$pb.TagNumber(315098849)
  void clearRolesessionname() => $_clearField(315098849);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(5);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(5);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(6);

  @$pb.TagNumber(418274661)
  $core.String get serialnumber => $_getSZ(7);
  @$pb.TagNumber(418274661)
  set serialnumber($core.String value) => $_setString(7, value);
  @$pb.TagNumber(418274661)
  $core.bool hasSerialnumber() => $_has(7);
  @$pb.TagNumber(418274661)
  void clearSerialnumber() => $_clearField(418274661);

  @$pb.TagNumber(451873635)
  $core.int get durationseconds => $_getIZ(8);
  @$pb.TagNumber(451873635)
  set durationseconds($core.int value) => $_setSignedInt32(8, value);
  @$pb.TagNumber(451873635)
  $core.bool hasDurationseconds() => $_has(8);
  @$pb.TagNumber(451873635)
  void clearDurationseconds() => $_clearField(451873635);

  @$pb.TagNumber(452608727)
  $pb.PbList<$core.String> get transitivetagkeys => $_getList(9);

  @$pb.TagNumber(466635355)
  $core.String get sourceidentity => $_getSZ(10);
  @$pb.TagNumber(466635355)
  set sourceidentity($core.String value) => $_setString(10, value);
  @$pb.TagNumber(466635355)
  $core.bool hasSourceidentity() => $_has(10);
  @$pb.TagNumber(466635355)
  void clearSourceidentity() => $_clearField(466635355);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(11);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(11, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(11);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class AssumeRoleResponse extends $pb.GeneratedMessage {
  factory AssumeRoleResponse({
    AssumedRoleUser? assumedroleuser,
    Credentials? credentials,
    $core.String? sourceidentity,
    $core.int? packedpolicysize,
  }) {
    final result = create();
    if (assumedroleuser != null) result.assumedroleuser = assumedroleuser;
    if (credentials != null) result.credentials = credentials;
    if (sourceidentity != null) result.sourceidentity = sourceidentity;
    if (packedpolicysize != null) result.packedpolicysize = packedpolicysize;
    return result;
  }

  AssumeRoleResponse._();

  factory AssumeRoleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRoleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRoleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOM<AssumedRoleUser>(314673579, _omitFieldNames ? '' : 'assumedroleuser',
        subBuilder: AssumedRoleUser.create)
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..aOS(466635355, _omitFieldNames ? '' : 'sourceidentity')
    ..aI(511234267, _omitFieldNames ? '' : 'packedpolicysize')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleResponse copyWith(void Function(AssumeRoleResponse) updates) =>
      super.copyWith((message) => updates(message as AssumeRoleResponse))
          as AssumeRoleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRoleResponse create() => AssumeRoleResponse._();
  @$core.override
  AssumeRoleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRoleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRoleResponse>(create);
  static AssumeRoleResponse? _defaultInstance;

  @$pb.TagNumber(314673579)
  AssumedRoleUser get assumedroleuser => $_getN(0);
  @$pb.TagNumber(314673579)
  set assumedroleuser(AssumedRoleUser value) => $_setField(314673579, value);
  @$pb.TagNumber(314673579)
  $core.bool hasAssumedroleuser() => $_has(0);
  @$pb.TagNumber(314673579)
  void clearAssumedroleuser() => $_clearField(314673579);
  @$pb.TagNumber(314673579)
  AssumedRoleUser ensureAssumedroleuser() => $_ensure(0);

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(1);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(1);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(1);

  @$pb.TagNumber(466635355)
  $core.String get sourceidentity => $_getSZ(2);
  @$pb.TagNumber(466635355)
  set sourceidentity($core.String value) => $_setString(2, value);
  @$pb.TagNumber(466635355)
  $core.bool hasSourceidentity() => $_has(2);
  @$pb.TagNumber(466635355)
  void clearSourceidentity() => $_clearField(466635355);

  @$pb.TagNumber(511234267)
  $core.int get packedpolicysize => $_getIZ(3);
  @$pb.TagNumber(511234267)
  set packedpolicysize($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(511234267)
  $core.bool hasPackedpolicysize() => $_has(3);
  @$pb.TagNumber(511234267)
  void clearPackedpolicysize() => $_clearField(511234267);
}

class AssumeRoleWithSAMLRequest extends $pb.GeneratedMessage {
  factory AssumeRoleWithSAMLRequest({
    $core.String? principalarn,
    $core.Iterable<PolicyDescriptorType>? policyarns,
    $core.String? samlassertion,
    $core.String? rolearn,
    $core.int? durationseconds,
    $core.String? policy,
  }) {
    final result = create();
    if (principalarn != null) result.principalarn = principalarn;
    if (policyarns != null) result.policyarns.addAll(policyarns);
    if (samlassertion != null) result.samlassertion = samlassertion;
    if (rolearn != null) result.rolearn = rolearn;
    if (durationseconds != null) result.durationseconds = durationseconds;
    if (policy != null) result.policy = policy;
    return result;
  }

  AssumeRoleWithSAMLRequest._();

  factory AssumeRoleWithSAMLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRoleWithSAMLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRoleWithSAMLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(93469969, _omitFieldNames ? '' : 'principalarn')
    ..pPM<PolicyDescriptorType>(183785508, _omitFieldNames ? '' : 'policyarns',
        subBuilder: PolicyDescriptorType.create)
    ..aOS(202921933, _omitFieldNames ? '' : 'samlassertion')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aI(451873635, _omitFieldNames ? '' : 'durationseconds')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithSAMLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithSAMLRequest copyWith(
          void Function(AssumeRoleWithSAMLRequest) updates) =>
      super.copyWith((message) => updates(message as AssumeRoleWithSAMLRequest))
          as AssumeRoleWithSAMLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithSAMLRequest create() => AssumeRoleWithSAMLRequest._();
  @$core.override
  AssumeRoleWithSAMLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithSAMLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRoleWithSAMLRequest>(create);
  static AssumeRoleWithSAMLRequest? _defaultInstance;

  @$pb.TagNumber(93469969)
  $core.String get principalarn => $_getSZ(0);
  @$pb.TagNumber(93469969)
  set principalarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(93469969)
  $core.bool hasPrincipalarn() => $_has(0);
  @$pb.TagNumber(93469969)
  void clearPrincipalarn() => $_clearField(93469969);

  @$pb.TagNumber(183785508)
  $pb.PbList<PolicyDescriptorType> get policyarns => $_getList(1);

  @$pb.TagNumber(202921933)
  $core.String get samlassertion => $_getSZ(2);
  @$pb.TagNumber(202921933)
  set samlassertion($core.String value) => $_setString(2, value);
  @$pb.TagNumber(202921933)
  $core.bool hasSamlassertion() => $_has(2);
  @$pb.TagNumber(202921933)
  void clearSamlassertion() => $_clearField(202921933);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(3);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(3);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(451873635)
  $core.int get durationseconds => $_getIZ(4);
  @$pb.TagNumber(451873635)
  set durationseconds($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(451873635)
  $core.bool hasDurationseconds() => $_has(4);
  @$pb.TagNumber(451873635)
  void clearDurationseconds() => $_clearField(451873635);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(5);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(5, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(5);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class AssumeRoleWithSAMLResponse extends $pb.GeneratedMessage {
  factory AssumeRoleWithSAMLResponse({
    $core.String? subject,
    $core.String? subjecttype,
    $core.String? audience,
    AssumedRoleUser? assumedroleuser,
    Credentials? credentials,
    $core.String? sourceidentity,
    $core.int? packedpolicysize,
    $core.String? namequalifier,
    $core.String? issuer,
  }) {
    final result = create();
    if (subject != null) result.subject = subject;
    if (subjecttype != null) result.subjecttype = subjecttype;
    if (audience != null) result.audience = audience;
    if (assumedroleuser != null) result.assumedroleuser = assumedroleuser;
    if (credentials != null) result.credentials = credentials;
    if (sourceidentity != null) result.sourceidentity = sourceidentity;
    if (packedpolicysize != null) result.packedpolicysize = packedpolicysize;
    if (namequalifier != null) result.namequalifier = namequalifier;
    if (issuer != null) result.issuer = issuer;
    return result;
  }

  AssumeRoleWithSAMLResponse._();

  factory AssumeRoleWithSAMLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRoleWithSAMLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRoleWithSAMLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(7939312, _omitFieldNames ? '' : 'subject')
    ..aOS(222881976, _omitFieldNames ? '' : 'subjecttype')
    ..aOS(284892548, _omitFieldNames ? '' : 'audience')
    ..aOM<AssumedRoleUser>(314673579, _omitFieldNames ? '' : 'assumedroleuser',
        subBuilder: AssumedRoleUser.create)
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..aOS(466635355, _omitFieldNames ? '' : 'sourceidentity')
    ..aI(511234267, _omitFieldNames ? '' : 'packedpolicysize')
    ..aOS(521907559, _omitFieldNames ? '' : 'namequalifier')
    ..aOS(528708823, _omitFieldNames ? '' : 'issuer')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithSAMLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithSAMLResponse copyWith(
          void Function(AssumeRoleWithSAMLResponse) updates) =>
      super.copyWith(
              (message) => updates(message as AssumeRoleWithSAMLResponse))
          as AssumeRoleWithSAMLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithSAMLResponse create() => AssumeRoleWithSAMLResponse._();
  @$core.override
  AssumeRoleWithSAMLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithSAMLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRoleWithSAMLResponse>(create);
  static AssumeRoleWithSAMLResponse? _defaultInstance;

  @$pb.TagNumber(7939312)
  $core.String get subject => $_getSZ(0);
  @$pb.TagNumber(7939312)
  set subject($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7939312)
  $core.bool hasSubject() => $_has(0);
  @$pb.TagNumber(7939312)
  void clearSubject() => $_clearField(7939312);

  @$pb.TagNumber(222881976)
  $core.String get subjecttype => $_getSZ(1);
  @$pb.TagNumber(222881976)
  set subjecttype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(222881976)
  $core.bool hasSubjecttype() => $_has(1);
  @$pb.TagNumber(222881976)
  void clearSubjecttype() => $_clearField(222881976);

  @$pb.TagNumber(284892548)
  $core.String get audience => $_getSZ(2);
  @$pb.TagNumber(284892548)
  set audience($core.String value) => $_setString(2, value);
  @$pb.TagNumber(284892548)
  $core.bool hasAudience() => $_has(2);
  @$pb.TagNumber(284892548)
  void clearAudience() => $_clearField(284892548);

  @$pb.TagNumber(314673579)
  AssumedRoleUser get assumedroleuser => $_getN(3);
  @$pb.TagNumber(314673579)
  set assumedroleuser(AssumedRoleUser value) => $_setField(314673579, value);
  @$pb.TagNumber(314673579)
  $core.bool hasAssumedroleuser() => $_has(3);
  @$pb.TagNumber(314673579)
  void clearAssumedroleuser() => $_clearField(314673579);
  @$pb.TagNumber(314673579)
  AssumedRoleUser ensureAssumedroleuser() => $_ensure(3);

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(4);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(4);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(4);

  @$pb.TagNumber(466635355)
  $core.String get sourceidentity => $_getSZ(5);
  @$pb.TagNumber(466635355)
  set sourceidentity($core.String value) => $_setString(5, value);
  @$pb.TagNumber(466635355)
  $core.bool hasSourceidentity() => $_has(5);
  @$pb.TagNumber(466635355)
  void clearSourceidentity() => $_clearField(466635355);

  @$pb.TagNumber(511234267)
  $core.int get packedpolicysize => $_getIZ(6);
  @$pb.TagNumber(511234267)
  set packedpolicysize($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(511234267)
  $core.bool hasPackedpolicysize() => $_has(6);
  @$pb.TagNumber(511234267)
  void clearPackedpolicysize() => $_clearField(511234267);

  @$pb.TagNumber(521907559)
  $core.String get namequalifier => $_getSZ(7);
  @$pb.TagNumber(521907559)
  set namequalifier($core.String value) => $_setString(7, value);
  @$pb.TagNumber(521907559)
  $core.bool hasNamequalifier() => $_has(7);
  @$pb.TagNumber(521907559)
  void clearNamequalifier() => $_clearField(521907559);

  @$pb.TagNumber(528708823)
  $core.String get issuer => $_getSZ(8);
  @$pb.TagNumber(528708823)
  set issuer($core.String value) => $_setString(8, value);
  @$pb.TagNumber(528708823)
  $core.bool hasIssuer() => $_has(8);
  @$pb.TagNumber(528708823)
  void clearIssuer() => $_clearField(528708823);
}

class AssumeRoleWithWebIdentityRequest extends $pb.GeneratedMessage {
  factory AssumeRoleWithWebIdentityRequest({
    $core.Iterable<PolicyDescriptorType>? policyarns,
    $core.String? webidentitytoken,
    $core.String? rolesessionname,
    $core.String? rolearn,
    $core.int? durationseconds,
    $core.String? policy,
    $core.String? providerid,
  }) {
    final result = create();
    if (policyarns != null) result.policyarns.addAll(policyarns);
    if (webidentitytoken != null) result.webidentitytoken = webidentitytoken;
    if (rolesessionname != null) result.rolesessionname = rolesessionname;
    if (rolearn != null) result.rolearn = rolearn;
    if (durationseconds != null) result.durationseconds = durationseconds;
    if (policy != null) result.policy = policy;
    if (providerid != null) result.providerid = providerid;
    return result;
  }

  AssumeRoleWithWebIdentityRequest._();

  factory AssumeRoleWithWebIdentityRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRoleWithWebIdentityRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRoleWithWebIdentityRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..pPM<PolicyDescriptorType>(183785508, _omitFieldNames ? '' : 'policyarns',
        subBuilder: PolicyDescriptorType.create)
    ..aOS(234014869, _omitFieldNames ? '' : 'webidentitytoken')
    ..aOS(315098849, _omitFieldNames ? '' : 'rolesessionname')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aI(451873635, _omitFieldNames ? '' : 'durationseconds')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..aOS(509712370, _omitFieldNames ? '' : 'providerid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithWebIdentityRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithWebIdentityRequest copyWith(
          void Function(AssumeRoleWithWebIdentityRequest) updates) =>
      super.copyWith(
              (message) => updates(message as AssumeRoleWithWebIdentityRequest))
          as AssumeRoleWithWebIdentityRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithWebIdentityRequest create() =>
      AssumeRoleWithWebIdentityRequest._();
  @$core.override
  AssumeRoleWithWebIdentityRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithWebIdentityRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRoleWithWebIdentityRequest>(
          create);
  static AssumeRoleWithWebIdentityRequest? _defaultInstance;

  @$pb.TagNumber(183785508)
  $pb.PbList<PolicyDescriptorType> get policyarns => $_getList(0);

  @$pb.TagNumber(234014869)
  $core.String get webidentitytoken => $_getSZ(1);
  @$pb.TagNumber(234014869)
  set webidentitytoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(234014869)
  $core.bool hasWebidentitytoken() => $_has(1);
  @$pb.TagNumber(234014869)
  void clearWebidentitytoken() => $_clearField(234014869);

  @$pb.TagNumber(315098849)
  $core.String get rolesessionname => $_getSZ(2);
  @$pb.TagNumber(315098849)
  set rolesessionname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(315098849)
  $core.bool hasRolesessionname() => $_has(2);
  @$pb.TagNumber(315098849)
  void clearRolesessionname() => $_clearField(315098849);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(3);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(3);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(451873635)
  $core.int get durationseconds => $_getIZ(4);
  @$pb.TagNumber(451873635)
  set durationseconds($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(451873635)
  $core.bool hasDurationseconds() => $_has(4);
  @$pb.TagNumber(451873635)
  void clearDurationseconds() => $_clearField(451873635);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(5);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(5, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(5);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);

  @$pb.TagNumber(509712370)
  $core.String get providerid => $_getSZ(6);
  @$pb.TagNumber(509712370)
  set providerid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(509712370)
  $core.bool hasProviderid() => $_has(6);
  @$pb.TagNumber(509712370)
  void clearProviderid() => $_clearField(509712370);
}

class AssumeRoleWithWebIdentityResponse extends $pb.GeneratedMessage {
  factory AssumeRoleWithWebIdentityResponse({
    $core.String? subjectfromwebidentitytoken,
    $core.String? audience,
    AssumedRoleUser? assumedroleuser,
    $core.String? provider,
    Credentials? credentials,
    $core.String? sourceidentity,
    $core.int? packedpolicysize,
  }) {
    final result = create();
    if (subjectfromwebidentitytoken != null)
      result.subjectfromwebidentitytoken = subjectfromwebidentitytoken;
    if (audience != null) result.audience = audience;
    if (assumedroleuser != null) result.assumedroleuser = assumedroleuser;
    if (provider != null) result.provider = provider;
    if (credentials != null) result.credentials = credentials;
    if (sourceidentity != null) result.sourceidentity = sourceidentity;
    if (packedpolicysize != null) result.packedpolicysize = packedpolicysize;
    return result;
  }

  AssumeRoleWithWebIdentityResponse._();

  factory AssumeRoleWithWebIdentityResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRoleWithWebIdentityResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRoleWithWebIdentityResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(96354739, _omitFieldNames ? '' : 'subjectfromwebidentitytoken')
    ..aOS(284892548, _omitFieldNames ? '' : 'audience')
    ..aOM<AssumedRoleUser>(314673579, _omitFieldNames ? '' : 'assumedroleuser',
        subBuilder: AssumedRoleUser.create)
    ..aOS(363366621, _omitFieldNames ? '' : 'provider')
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..aOS(466635355, _omitFieldNames ? '' : 'sourceidentity')
    ..aI(511234267, _omitFieldNames ? '' : 'packedpolicysize')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithWebIdentityResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRoleWithWebIdentityResponse copyWith(
          void Function(AssumeRoleWithWebIdentityResponse) updates) =>
      super.copyWith((message) =>
              updates(message as AssumeRoleWithWebIdentityResponse))
          as AssumeRoleWithWebIdentityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithWebIdentityResponse create() =>
      AssumeRoleWithWebIdentityResponse._();
  @$core.override
  AssumeRoleWithWebIdentityResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRoleWithWebIdentityResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRoleWithWebIdentityResponse>(
          create);
  static AssumeRoleWithWebIdentityResponse? _defaultInstance;

  @$pb.TagNumber(96354739)
  $core.String get subjectfromwebidentitytoken => $_getSZ(0);
  @$pb.TagNumber(96354739)
  set subjectfromwebidentitytoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(96354739)
  $core.bool hasSubjectfromwebidentitytoken() => $_has(0);
  @$pb.TagNumber(96354739)
  void clearSubjectfromwebidentitytoken() => $_clearField(96354739);

  @$pb.TagNumber(284892548)
  $core.String get audience => $_getSZ(1);
  @$pb.TagNumber(284892548)
  set audience($core.String value) => $_setString(1, value);
  @$pb.TagNumber(284892548)
  $core.bool hasAudience() => $_has(1);
  @$pb.TagNumber(284892548)
  void clearAudience() => $_clearField(284892548);

  @$pb.TagNumber(314673579)
  AssumedRoleUser get assumedroleuser => $_getN(2);
  @$pb.TagNumber(314673579)
  set assumedroleuser(AssumedRoleUser value) => $_setField(314673579, value);
  @$pb.TagNumber(314673579)
  $core.bool hasAssumedroleuser() => $_has(2);
  @$pb.TagNumber(314673579)
  void clearAssumedroleuser() => $_clearField(314673579);
  @$pb.TagNumber(314673579)
  AssumedRoleUser ensureAssumedroleuser() => $_ensure(2);

  @$pb.TagNumber(363366621)
  $core.String get provider => $_getSZ(3);
  @$pb.TagNumber(363366621)
  set provider($core.String value) => $_setString(3, value);
  @$pb.TagNumber(363366621)
  $core.bool hasProvider() => $_has(3);
  @$pb.TagNumber(363366621)
  void clearProvider() => $_clearField(363366621);

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(4);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(4);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(4);

  @$pb.TagNumber(466635355)
  $core.String get sourceidentity => $_getSZ(5);
  @$pb.TagNumber(466635355)
  set sourceidentity($core.String value) => $_setString(5, value);
  @$pb.TagNumber(466635355)
  $core.bool hasSourceidentity() => $_has(5);
  @$pb.TagNumber(466635355)
  void clearSourceidentity() => $_clearField(466635355);

  @$pb.TagNumber(511234267)
  $core.int get packedpolicysize => $_getIZ(6);
  @$pb.TagNumber(511234267)
  set packedpolicysize($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(511234267)
  $core.bool hasPackedpolicysize() => $_has(6);
  @$pb.TagNumber(511234267)
  void clearPackedpolicysize() => $_clearField(511234267);
}

class AssumeRootRequest extends $pb.GeneratedMessage {
  factory AssumeRootRequest({
    $core.String? targetprincipal,
    PolicyDescriptorType? taskpolicyarn,
    $core.int? durationseconds,
  }) {
    final result = create();
    if (targetprincipal != null) result.targetprincipal = targetprincipal;
    if (taskpolicyarn != null) result.taskpolicyarn = taskpolicyarn;
    if (durationseconds != null) result.durationseconds = durationseconds;
    return result;
  }

  AssumeRootRequest._();

  factory AssumeRootRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRootRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRootRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(215916667, _omitFieldNames ? '' : 'targetprincipal')
    ..aOM<PolicyDescriptorType>(
        332406370, _omitFieldNames ? '' : 'taskpolicyarn',
        subBuilder: PolicyDescriptorType.create)
    ..aI(451873635, _omitFieldNames ? '' : 'durationseconds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRootRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRootRequest copyWith(void Function(AssumeRootRequest) updates) =>
      super.copyWith((message) => updates(message as AssumeRootRequest))
          as AssumeRootRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRootRequest create() => AssumeRootRequest._();
  @$core.override
  AssumeRootRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRootRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRootRequest>(create);
  static AssumeRootRequest? _defaultInstance;

  @$pb.TagNumber(215916667)
  $core.String get targetprincipal => $_getSZ(0);
  @$pb.TagNumber(215916667)
  set targetprincipal($core.String value) => $_setString(0, value);
  @$pb.TagNumber(215916667)
  $core.bool hasTargetprincipal() => $_has(0);
  @$pb.TagNumber(215916667)
  void clearTargetprincipal() => $_clearField(215916667);

  @$pb.TagNumber(332406370)
  PolicyDescriptorType get taskpolicyarn => $_getN(1);
  @$pb.TagNumber(332406370)
  set taskpolicyarn(PolicyDescriptorType value) => $_setField(332406370, value);
  @$pb.TagNumber(332406370)
  $core.bool hasTaskpolicyarn() => $_has(1);
  @$pb.TagNumber(332406370)
  void clearTaskpolicyarn() => $_clearField(332406370);
  @$pb.TagNumber(332406370)
  PolicyDescriptorType ensureTaskpolicyarn() => $_ensure(1);

  @$pb.TagNumber(451873635)
  $core.int get durationseconds => $_getIZ(2);
  @$pb.TagNumber(451873635)
  set durationseconds($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(451873635)
  $core.bool hasDurationseconds() => $_has(2);
  @$pb.TagNumber(451873635)
  void clearDurationseconds() => $_clearField(451873635);
}

class AssumeRootResponse extends $pb.GeneratedMessage {
  factory AssumeRootResponse({
    Credentials? credentials,
    $core.String? sourceidentity,
  }) {
    final result = create();
    if (credentials != null) result.credentials = credentials;
    if (sourceidentity != null) result.sourceidentity = sourceidentity;
    return result;
  }

  AssumeRootResponse._();

  factory AssumeRootResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumeRootResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumeRootResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..aOS(466635355, _omitFieldNames ? '' : 'sourceidentity')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRootResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumeRootResponse copyWith(void Function(AssumeRootResponse) updates) =>
      super.copyWith((message) => updates(message as AssumeRootResponse))
          as AssumeRootResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumeRootResponse create() => AssumeRootResponse._();
  @$core.override
  AssumeRootResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumeRootResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumeRootResponse>(create);
  static AssumeRootResponse? _defaultInstance;

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(0);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(0);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(0);

  @$pb.TagNumber(466635355)
  $core.String get sourceidentity => $_getSZ(1);
  @$pb.TagNumber(466635355)
  set sourceidentity($core.String value) => $_setString(1, value);
  @$pb.TagNumber(466635355)
  $core.bool hasSourceidentity() => $_has(1);
  @$pb.TagNumber(466635355)
  void clearSourceidentity() => $_clearField(466635355);
}

class AssumedRoleUser extends $pb.GeneratedMessage {
  factory AssumedRoleUser({
    $core.String? assumedroleid,
    $core.String? arn,
  }) {
    final result = create();
    if (assumedroleid != null) result.assumedroleid = assumedroleid;
    if (arn != null) result.arn = arn;
    return result;
  }

  AssumedRoleUser._();

  factory AssumedRoleUser.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssumedRoleUser.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssumedRoleUser',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(244783337, _omitFieldNames ? '' : 'assumedroleid')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumedRoleUser clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssumedRoleUser copyWith(void Function(AssumedRoleUser) updates) =>
      super.copyWith((message) => updates(message as AssumedRoleUser))
          as AssumedRoleUser;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssumedRoleUser create() => AssumedRoleUser._();
  @$core.override
  AssumedRoleUser createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssumedRoleUser getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssumedRoleUser>(create);
  static AssumedRoleUser? _defaultInstance;

  @$pb.TagNumber(244783337)
  $core.String get assumedroleid => $_getSZ(0);
  @$pb.TagNumber(244783337)
  set assumedroleid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(244783337)
  $core.bool hasAssumedroleid() => $_has(0);
  @$pb.TagNumber(244783337)
  void clearAssumedroleid() => $_clearField(244783337);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class Credentials extends $pb.GeneratedMessage {
  factory Credentials({
    $core.String? secretaccesskey,
    $core.String? sessiontoken,
    $core.String? expiration,
    $core.String? accesskeyid,
  }) {
    final result = create();
    if (secretaccesskey != null) result.secretaccesskey = secretaccesskey;
    if (sessiontoken != null) result.sessiontoken = sessiontoken;
    if (expiration != null) result.expiration = expiration;
    if (accesskeyid != null) result.accesskeyid = accesskeyid;
    return result;
  }

  Credentials._();

  factory Credentials.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Credentials.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Credentials',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(172288487, _omitFieldNames ? '' : 'secretaccesskey')
    ..aOS(211161069, _omitFieldNames ? '' : 'sessiontoken')
    ..aOS(245879945, _omitFieldNames ? '' : 'expiration')
    ..aOS(453893024, _omitFieldNames ? '' : 'accesskeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Credentials clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Credentials copyWith(void Function(Credentials) updates) =>
      super.copyWith((message) => updates(message as Credentials))
          as Credentials;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Credentials create() => Credentials._();
  @$core.override
  Credentials createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Credentials getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Credentials>(create);
  static Credentials? _defaultInstance;

  @$pb.TagNumber(172288487)
  $core.String get secretaccesskey => $_getSZ(0);
  @$pb.TagNumber(172288487)
  set secretaccesskey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(172288487)
  $core.bool hasSecretaccesskey() => $_has(0);
  @$pb.TagNumber(172288487)
  void clearSecretaccesskey() => $_clearField(172288487);

  @$pb.TagNumber(211161069)
  $core.String get sessiontoken => $_getSZ(1);
  @$pb.TagNumber(211161069)
  set sessiontoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(211161069)
  $core.bool hasSessiontoken() => $_has(1);
  @$pb.TagNumber(211161069)
  void clearSessiontoken() => $_clearField(211161069);

  @$pb.TagNumber(245879945)
  $core.String get expiration => $_getSZ(2);
  @$pb.TagNumber(245879945)
  set expiration($core.String value) => $_setString(2, value);
  @$pb.TagNumber(245879945)
  $core.bool hasExpiration() => $_has(2);
  @$pb.TagNumber(245879945)
  void clearExpiration() => $_clearField(245879945);

  @$pb.TagNumber(453893024)
  $core.String get accesskeyid => $_getSZ(3);
  @$pb.TagNumber(453893024)
  set accesskeyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(453893024)
  $core.bool hasAccesskeyid() => $_has(3);
  @$pb.TagNumber(453893024)
  void clearAccesskeyid() => $_clearField(453893024);
}

class DecodeAuthorizationMessageRequest extends $pb.GeneratedMessage {
  factory DecodeAuthorizationMessageRequest({
    $core.String? encodedmessage,
  }) {
    final result = create();
    if (encodedmessage != null) result.encodedmessage = encodedmessage;
    return result;
  }

  DecodeAuthorizationMessageRequest._();

  factory DecodeAuthorizationMessageRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecodeAuthorizationMessageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecodeAuthorizationMessageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(176632049, _omitFieldNames ? '' : 'encodedmessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeAuthorizationMessageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeAuthorizationMessageRequest copyWith(
          void Function(DecodeAuthorizationMessageRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DecodeAuthorizationMessageRequest))
          as DecodeAuthorizationMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecodeAuthorizationMessageRequest create() =>
      DecodeAuthorizationMessageRequest._();
  @$core.override
  DecodeAuthorizationMessageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecodeAuthorizationMessageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecodeAuthorizationMessageRequest>(
          create);
  static DecodeAuthorizationMessageRequest? _defaultInstance;

  @$pb.TagNumber(176632049)
  $core.String get encodedmessage => $_getSZ(0);
  @$pb.TagNumber(176632049)
  set encodedmessage($core.String value) => $_setString(0, value);
  @$pb.TagNumber(176632049)
  $core.bool hasEncodedmessage() => $_has(0);
  @$pb.TagNumber(176632049)
  void clearEncodedmessage() => $_clearField(176632049);
}

class DecodeAuthorizationMessageResponse extends $pb.GeneratedMessage {
  factory DecodeAuthorizationMessageResponse({
    $core.String? decodedmessage,
  }) {
    final result = create();
    if (decodedmessage != null) result.decodedmessage = decodedmessage;
    return result;
  }

  DecodeAuthorizationMessageResponse._();

  factory DecodeAuthorizationMessageResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecodeAuthorizationMessageResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecodeAuthorizationMessageResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(475373641, _omitFieldNames ? '' : 'decodedmessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeAuthorizationMessageResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeAuthorizationMessageResponse copyWith(
          void Function(DecodeAuthorizationMessageResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DecodeAuthorizationMessageResponse))
          as DecodeAuthorizationMessageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecodeAuthorizationMessageResponse create() =>
      DecodeAuthorizationMessageResponse._();
  @$core.override
  DecodeAuthorizationMessageResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecodeAuthorizationMessageResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecodeAuthorizationMessageResponse>(
          create);
  static DecodeAuthorizationMessageResponse? _defaultInstance;

  @$pb.TagNumber(475373641)
  $core.String get decodedmessage => $_getSZ(0);
  @$pb.TagNumber(475373641)
  set decodedmessage($core.String value) => $_setString(0, value);
  @$pb.TagNumber(475373641)
  $core.bool hasDecodedmessage() => $_has(0);
  @$pb.TagNumber(475373641)
  void clearDecodedmessage() => $_clearField(475373641);
}

class ExpiredTokenException extends $pb.GeneratedMessage {
  factory ExpiredTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExpiredTokenException._();

  factory ExpiredTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExpiredTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExpiredTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredTokenException copyWith(
          void Function(ExpiredTokenException) updates) =>
      super.copyWith((message) => updates(message as ExpiredTokenException))
          as ExpiredTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExpiredTokenException create() => ExpiredTokenException._();
  @$core.override
  ExpiredTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExpiredTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExpiredTokenException>(create);
  static ExpiredTokenException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExpiredTradeInTokenException extends $pb.GeneratedMessage {
  factory ExpiredTradeInTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExpiredTradeInTokenException._();

  factory ExpiredTradeInTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExpiredTradeInTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExpiredTradeInTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredTradeInTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredTradeInTokenException copyWith(
          void Function(ExpiredTradeInTokenException) updates) =>
      super.copyWith(
              (message) => updates(message as ExpiredTradeInTokenException))
          as ExpiredTradeInTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExpiredTradeInTokenException create() =>
      ExpiredTradeInTokenException._();
  @$core.override
  ExpiredTradeInTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExpiredTradeInTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExpiredTradeInTokenException>(create);
  static ExpiredTradeInTokenException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class FederatedUser extends $pb.GeneratedMessage {
  factory FederatedUser({
    $core.String? federateduserid,
    $core.String? arn,
  }) {
    final result = create();
    if (federateduserid != null) result.federateduserid = federateduserid;
    if (arn != null) result.arn = arn;
    return result;
  }

  FederatedUser._();

  factory FederatedUser.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FederatedUser.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FederatedUser',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(366276542, _omitFieldNames ? '' : 'federateduserid')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FederatedUser clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FederatedUser copyWith(void Function(FederatedUser) updates) =>
      super.copyWith((message) => updates(message as FederatedUser))
          as FederatedUser;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FederatedUser create() => FederatedUser._();
  @$core.override
  FederatedUser createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FederatedUser getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FederatedUser>(create);
  static FederatedUser? _defaultInstance;

  @$pb.TagNumber(366276542)
  $core.String get federateduserid => $_getSZ(0);
  @$pb.TagNumber(366276542)
  set federateduserid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(366276542)
  $core.bool hasFederateduserid() => $_has(0);
  @$pb.TagNumber(366276542)
  void clearFederateduserid() => $_clearField(366276542);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class GetAccessKeyInfoRequest extends $pb.GeneratedMessage {
  factory GetAccessKeyInfoRequest({
    $core.String? accesskeyid,
  }) {
    final result = create();
    if (accesskeyid != null) result.accesskeyid = accesskeyid;
    return result;
  }

  GetAccessKeyInfoRequest._();

  factory GetAccessKeyInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccessKeyInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccessKeyInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(453893024, _omitFieldNames ? '' : 'accesskeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccessKeyInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccessKeyInfoRequest copyWith(
          void Function(GetAccessKeyInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetAccessKeyInfoRequest))
          as GetAccessKeyInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccessKeyInfoRequest create() => GetAccessKeyInfoRequest._();
  @$core.override
  GetAccessKeyInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccessKeyInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccessKeyInfoRequest>(create);
  static GetAccessKeyInfoRequest? _defaultInstance;

  @$pb.TagNumber(453893024)
  $core.String get accesskeyid => $_getSZ(0);
  @$pb.TagNumber(453893024)
  set accesskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(453893024)
  $core.bool hasAccesskeyid() => $_has(0);
  @$pb.TagNumber(453893024)
  void clearAccesskeyid() => $_clearField(453893024);
}

class GetAccessKeyInfoResponse extends $pb.GeneratedMessage {
  factory GetAccessKeyInfoResponse({
    $core.String? account,
  }) {
    final result = create();
    if (account != null) result.account = account;
    return result;
  }

  GetAccessKeyInfoResponse._();

  factory GetAccessKeyInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccessKeyInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccessKeyInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(435725053, _omitFieldNames ? '' : 'account')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccessKeyInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccessKeyInfoResponse copyWith(
          void Function(GetAccessKeyInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetAccessKeyInfoResponse))
          as GetAccessKeyInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccessKeyInfoResponse create() => GetAccessKeyInfoResponse._();
  @$core.override
  GetAccessKeyInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccessKeyInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccessKeyInfoResponse>(create);
  static GetAccessKeyInfoResponse? _defaultInstance;

  @$pb.TagNumber(435725053)
  $core.String get account => $_getSZ(0);
  @$pb.TagNumber(435725053)
  set account($core.String value) => $_setString(0, value);
  @$pb.TagNumber(435725053)
  $core.bool hasAccount() => $_has(0);
  @$pb.TagNumber(435725053)
  void clearAccount() => $_clearField(435725053);
}

class GetCallerIdentityRequest extends $pb.GeneratedMessage {
  factory GetCallerIdentityRequest() => create();

  GetCallerIdentityRequest._();

  factory GetCallerIdentityRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCallerIdentityRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCallerIdentityRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCallerIdentityRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCallerIdentityRequest copyWith(
          void Function(GetCallerIdentityRequest) updates) =>
      super.copyWith((message) => updates(message as GetCallerIdentityRequest))
          as GetCallerIdentityRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCallerIdentityRequest create() => GetCallerIdentityRequest._();
  @$core.override
  GetCallerIdentityRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCallerIdentityRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCallerIdentityRequest>(create);
  static GetCallerIdentityRequest? _defaultInstance;
}

class GetCallerIdentityResponse extends $pb.GeneratedMessage {
  factory GetCallerIdentityResponse({
    $core.String? userid,
    $core.String? arn,
    $core.String? account,
  }) {
    final result = create();
    if (userid != null) result.userid = userid;
    if (arn != null) result.arn = arn;
    if (account != null) result.account = account;
    return result;
  }

  GetCallerIdentityResponse._();

  factory GetCallerIdentityResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCallerIdentityResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCallerIdentityResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(10274112, _omitFieldNames ? '' : 'userid')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(435725053, _omitFieldNames ? '' : 'account')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCallerIdentityResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCallerIdentityResponse copyWith(
          void Function(GetCallerIdentityResponse) updates) =>
      super.copyWith((message) => updates(message as GetCallerIdentityResponse))
          as GetCallerIdentityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCallerIdentityResponse create() => GetCallerIdentityResponse._();
  @$core.override
  GetCallerIdentityResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCallerIdentityResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCallerIdentityResponse>(create);
  static GetCallerIdentityResponse? _defaultInstance;

  @$pb.TagNumber(10274112)
  $core.String get userid => $_getSZ(0);
  @$pb.TagNumber(10274112)
  set userid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(10274112)
  $core.bool hasUserid() => $_has(0);
  @$pb.TagNumber(10274112)
  void clearUserid() => $_clearField(10274112);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(435725053)
  $core.String get account => $_getSZ(2);
  @$pb.TagNumber(435725053)
  set account($core.String value) => $_setString(2, value);
  @$pb.TagNumber(435725053)
  $core.bool hasAccount() => $_has(2);
  @$pb.TagNumber(435725053)
  void clearAccount() => $_clearField(435725053);
}

class GetDelegatedAccessTokenRequest extends $pb.GeneratedMessage {
  factory GetDelegatedAccessTokenRequest({
    $core.String? tradeintoken,
  }) {
    final result = create();
    if (tradeintoken != null) result.tradeintoken = tradeintoken;
    return result;
  }

  GetDelegatedAccessTokenRequest._();

  factory GetDelegatedAccessTokenRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDelegatedAccessTokenRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDelegatedAccessTokenRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(10008816, _omitFieldNames ? '' : 'tradeintoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDelegatedAccessTokenRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDelegatedAccessTokenRequest copyWith(
          void Function(GetDelegatedAccessTokenRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetDelegatedAccessTokenRequest))
          as GetDelegatedAccessTokenRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDelegatedAccessTokenRequest create() =>
      GetDelegatedAccessTokenRequest._();
  @$core.override
  GetDelegatedAccessTokenRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDelegatedAccessTokenRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDelegatedAccessTokenRequest>(create);
  static GetDelegatedAccessTokenRequest? _defaultInstance;

  @$pb.TagNumber(10008816)
  $core.String get tradeintoken => $_getSZ(0);
  @$pb.TagNumber(10008816)
  set tradeintoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(10008816)
  $core.bool hasTradeintoken() => $_has(0);
  @$pb.TagNumber(10008816)
  void clearTradeintoken() => $_clearField(10008816);
}

class GetDelegatedAccessTokenResponse extends $pb.GeneratedMessage {
  factory GetDelegatedAccessTokenResponse({
    $core.String? assumedprincipal,
    Credentials? credentials,
    $core.int? packedpolicysize,
  }) {
    final result = create();
    if (assumedprincipal != null) result.assumedprincipal = assumedprincipal;
    if (credentials != null) result.credentials = credentials;
    if (packedpolicysize != null) result.packedpolicysize = packedpolicysize;
    return result;
  }

  GetDelegatedAccessTokenResponse._();

  factory GetDelegatedAccessTokenResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDelegatedAccessTokenResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDelegatedAccessTokenResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(359093742, _omitFieldNames ? '' : 'assumedprincipal')
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..aI(511234267, _omitFieldNames ? '' : 'packedpolicysize')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDelegatedAccessTokenResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDelegatedAccessTokenResponse copyWith(
          void Function(GetDelegatedAccessTokenResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetDelegatedAccessTokenResponse))
          as GetDelegatedAccessTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDelegatedAccessTokenResponse create() =>
      GetDelegatedAccessTokenResponse._();
  @$core.override
  GetDelegatedAccessTokenResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDelegatedAccessTokenResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDelegatedAccessTokenResponse>(
          create);
  static GetDelegatedAccessTokenResponse? _defaultInstance;

  @$pb.TagNumber(359093742)
  $core.String get assumedprincipal => $_getSZ(0);
  @$pb.TagNumber(359093742)
  set assumedprincipal($core.String value) => $_setString(0, value);
  @$pb.TagNumber(359093742)
  $core.bool hasAssumedprincipal() => $_has(0);
  @$pb.TagNumber(359093742)
  void clearAssumedprincipal() => $_clearField(359093742);

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(1);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(1);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(1);

  @$pb.TagNumber(511234267)
  $core.int get packedpolicysize => $_getIZ(2);
  @$pb.TagNumber(511234267)
  set packedpolicysize($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(511234267)
  $core.bool hasPackedpolicysize() => $_has(2);
  @$pb.TagNumber(511234267)
  void clearPackedpolicysize() => $_clearField(511234267);
}

class GetFederationTokenRequest extends $pb.GeneratedMessage {
  factory GetFederationTokenRequest({
    $core.Iterable<PolicyDescriptorType>? policyarns,
    $core.String? name,
    $core.Iterable<Tag>? tags,
    $core.int? durationseconds,
    $core.String? policy,
  }) {
    final result = create();
    if (policyarns != null) result.policyarns.addAll(policyarns);
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    if (durationseconds != null) result.durationseconds = durationseconds;
    if (policy != null) result.policy = policy;
    return result;
  }

  GetFederationTokenRequest._();

  factory GetFederationTokenRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetFederationTokenRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetFederationTokenRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..pPM<PolicyDescriptorType>(183785508, _omitFieldNames ? '' : 'policyarns',
        subBuilder: PolicyDescriptorType.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aI(451873635, _omitFieldNames ? '' : 'durationseconds')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetFederationTokenRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetFederationTokenRequest copyWith(
          void Function(GetFederationTokenRequest) updates) =>
      super.copyWith((message) => updates(message as GetFederationTokenRequest))
          as GetFederationTokenRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetFederationTokenRequest create() => GetFederationTokenRequest._();
  @$core.override
  GetFederationTokenRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetFederationTokenRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetFederationTokenRequest>(create);
  static GetFederationTokenRequest? _defaultInstance;

  @$pb.TagNumber(183785508)
  $pb.PbList<PolicyDescriptorType> get policyarns => $_getList(0);

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

  @$pb.TagNumber(451873635)
  $core.int get durationseconds => $_getIZ(3);
  @$pb.TagNumber(451873635)
  set durationseconds($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(451873635)
  $core.bool hasDurationseconds() => $_has(3);
  @$pb.TagNumber(451873635)
  void clearDurationseconds() => $_clearField(451873635);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(4);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(4, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(4);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class GetFederationTokenResponse extends $pb.GeneratedMessage {
  factory GetFederationTokenResponse({
    FederatedUser? federateduser,
    Credentials? credentials,
    $core.int? packedpolicysize,
  }) {
    final result = create();
    if (federateduser != null) result.federateduser = federateduser;
    if (credentials != null) result.credentials = credentials;
    if (packedpolicysize != null) result.packedpolicysize = packedpolicysize;
    return result;
  }

  GetFederationTokenResponse._();

  factory GetFederationTokenResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetFederationTokenResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetFederationTokenResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOM<FederatedUser>(328666081, _omitFieldNames ? '' : 'federateduser',
        subBuilder: FederatedUser.create)
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..aI(511234267, _omitFieldNames ? '' : 'packedpolicysize')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetFederationTokenResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetFederationTokenResponse copyWith(
          void Function(GetFederationTokenResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetFederationTokenResponse))
          as GetFederationTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetFederationTokenResponse create() => GetFederationTokenResponse._();
  @$core.override
  GetFederationTokenResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetFederationTokenResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetFederationTokenResponse>(create);
  static GetFederationTokenResponse? _defaultInstance;

  @$pb.TagNumber(328666081)
  FederatedUser get federateduser => $_getN(0);
  @$pb.TagNumber(328666081)
  set federateduser(FederatedUser value) => $_setField(328666081, value);
  @$pb.TagNumber(328666081)
  $core.bool hasFederateduser() => $_has(0);
  @$pb.TagNumber(328666081)
  void clearFederateduser() => $_clearField(328666081);
  @$pb.TagNumber(328666081)
  FederatedUser ensureFederateduser() => $_ensure(0);

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(1);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(1);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(1);

  @$pb.TagNumber(511234267)
  $core.int get packedpolicysize => $_getIZ(2);
  @$pb.TagNumber(511234267)
  set packedpolicysize($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(511234267)
  $core.bool hasPackedpolicysize() => $_has(2);
  @$pb.TagNumber(511234267)
  void clearPackedpolicysize() => $_clearField(511234267);
}

class GetSessionTokenRequest extends $pb.GeneratedMessage {
  factory GetSessionTokenRequest({
    $core.String? tokencode,
    $core.String? serialnumber,
    $core.int? durationseconds,
  }) {
    final result = create();
    if (tokencode != null) result.tokencode = tokencode;
    if (serialnumber != null) result.serialnumber = serialnumber;
    if (durationseconds != null) result.durationseconds = durationseconds;
    return result;
  }

  GetSessionTokenRequest._();

  factory GetSessionTokenRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionTokenRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionTokenRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(300671456, _omitFieldNames ? '' : 'tokencode')
    ..aOS(418274661, _omitFieldNames ? '' : 'serialnumber')
    ..aI(451873635, _omitFieldNames ? '' : 'durationseconds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionTokenRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionTokenRequest copyWith(
          void Function(GetSessionTokenRequest) updates) =>
      super.copyWith((message) => updates(message as GetSessionTokenRequest))
          as GetSessionTokenRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionTokenRequest create() => GetSessionTokenRequest._();
  @$core.override
  GetSessionTokenRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionTokenRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionTokenRequest>(create);
  static GetSessionTokenRequest? _defaultInstance;

  @$pb.TagNumber(300671456)
  $core.String get tokencode => $_getSZ(0);
  @$pb.TagNumber(300671456)
  set tokencode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(300671456)
  $core.bool hasTokencode() => $_has(0);
  @$pb.TagNumber(300671456)
  void clearTokencode() => $_clearField(300671456);

  @$pb.TagNumber(418274661)
  $core.String get serialnumber => $_getSZ(1);
  @$pb.TagNumber(418274661)
  set serialnumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(418274661)
  $core.bool hasSerialnumber() => $_has(1);
  @$pb.TagNumber(418274661)
  void clearSerialnumber() => $_clearField(418274661);

  @$pb.TagNumber(451873635)
  $core.int get durationseconds => $_getIZ(2);
  @$pb.TagNumber(451873635)
  set durationseconds($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(451873635)
  $core.bool hasDurationseconds() => $_has(2);
  @$pb.TagNumber(451873635)
  void clearDurationseconds() => $_clearField(451873635);
}

class GetSessionTokenResponse extends $pb.GeneratedMessage {
  factory GetSessionTokenResponse({
    Credentials? credentials,
  }) {
    final result = create();
    if (credentials != null) result.credentials = credentials;
    return result;
  }

  GetSessionTokenResponse._();

  factory GetSessionTokenResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSessionTokenResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSessionTokenResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionTokenResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSessionTokenResponse copyWith(
          void Function(GetSessionTokenResponse) updates) =>
      super.copyWith((message) => updates(message as GetSessionTokenResponse))
          as GetSessionTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSessionTokenResponse create() => GetSessionTokenResponse._();
  @$core.override
  GetSessionTokenResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSessionTokenResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSessionTokenResponse>(create);
  static GetSessionTokenResponse? _defaultInstance;

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(0);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(0);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(0);
}

class GetWebIdentityTokenRequest extends $pb.GeneratedMessage {
  factory GetWebIdentityTokenRequest({
    $core.Iterable<$core.String>? audience,
    $core.Iterable<Tag>? tags,
    $core.int? durationseconds,
    $core.String? signingalgorithm,
  }) {
    final result = create();
    if (audience != null) result.audience.addAll(audience);
    if (tags != null) result.tags.addAll(tags);
    if (durationseconds != null) result.durationseconds = durationseconds;
    if (signingalgorithm != null) result.signingalgorithm = signingalgorithm;
    return result;
  }

  GetWebIdentityTokenRequest._();

  factory GetWebIdentityTokenRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWebIdentityTokenRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWebIdentityTokenRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..pPS(284892548, _omitFieldNames ? '' : 'audience')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aI(451873635, _omitFieldNames ? '' : 'durationseconds')
    ..aOS(488091842, _omitFieldNames ? '' : 'signingalgorithm')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebIdentityTokenRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebIdentityTokenRequest copyWith(
          void Function(GetWebIdentityTokenRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetWebIdentityTokenRequest))
          as GetWebIdentityTokenRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWebIdentityTokenRequest create() => GetWebIdentityTokenRequest._();
  @$core.override
  GetWebIdentityTokenRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWebIdentityTokenRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWebIdentityTokenRequest>(create);
  static GetWebIdentityTokenRequest? _defaultInstance;

  @$pb.TagNumber(284892548)
  $pb.PbList<$core.String> get audience => $_getList(0);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);

  @$pb.TagNumber(451873635)
  $core.int get durationseconds => $_getIZ(2);
  @$pb.TagNumber(451873635)
  set durationseconds($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(451873635)
  $core.bool hasDurationseconds() => $_has(2);
  @$pb.TagNumber(451873635)
  void clearDurationseconds() => $_clearField(451873635);

  @$pb.TagNumber(488091842)
  $core.String get signingalgorithm => $_getSZ(3);
  @$pb.TagNumber(488091842)
  set signingalgorithm($core.String value) => $_setString(3, value);
  @$pb.TagNumber(488091842)
  $core.bool hasSigningalgorithm() => $_has(3);
  @$pb.TagNumber(488091842)
  void clearSigningalgorithm() => $_clearField(488091842);
}

class GetWebIdentityTokenResponse extends $pb.GeneratedMessage {
  factory GetWebIdentityTokenResponse({
    $core.String? webidentitytoken,
    $core.String? expiration,
  }) {
    final result = create();
    if (webidentitytoken != null) result.webidentitytoken = webidentitytoken;
    if (expiration != null) result.expiration = expiration;
    return result;
  }

  GetWebIdentityTokenResponse._();

  factory GetWebIdentityTokenResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWebIdentityTokenResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWebIdentityTokenResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(234014869, _omitFieldNames ? '' : 'webidentitytoken')
    ..aOS(245879945, _omitFieldNames ? '' : 'expiration')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebIdentityTokenResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebIdentityTokenResponse copyWith(
          void Function(GetWebIdentityTokenResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetWebIdentityTokenResponse))
          as GetWebIdentityTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWebIdentityTokenResponse create() =>
      GetWebIdentityTokenResponse._();
  @$core.override
  GetWebIdentityTokenResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWebIdentityTokenResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWebIdentityTokenResponse>(create);
  static GetWebIdentityTokenResponse? _defaultInstance;

  @$pb.TagNumber(234014869)
  $core.String get webidentitytoken => $_getSZ(0);
  @$pb.TagNumber(234014869)
  set webidentitytoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234014869)
  $core.bool hasWebidentitytoken() => $_has(0);
  @$pb.TagNumber(234014869)
  void clearWebidentitytoken() => $_clearField(234014869);

  @$pb.TagNumber(245879945)
  $core.String get expiration => $_getSZ(1);
  @$pb.TagNumber(245879945)
  set expiration($core.String value) => $_setString(1, value);
  @$pb.TagNumber(245879945)
  $core.bool hasExpiration() => $_has(1);
  @$pb.TagNumber(245879945)
  void clearExpiration() => $_clearField(245879945);
}

class IDPCommunicationErrorException extends $pb.GeneratedMessage {
  factory IDPCommunicationErrorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IDPCommunicationErrorException._();

  factory IDPCommunicationErrorException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IDPCommunicationErrorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IDPCommunicationErrorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IDPCommunicationErrorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IDPCommunicationErrorException copyWith(
          void Function(IDPCommunicationErrorException) updates) =>
      super.copyWith(
              (message) => updates(message as IDPCommunicationErrorException))
          as IDPCommunicationErrorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IDPCommunicationErrorException create() =>
      IDPCommunicationErrorException._();
  @$core.override
  IDPCommunicationErrorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IDPCommunicationErrorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IDPCommunicationErrorException>(create);
  static IDPCommunicationErrorException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class IDPRejectedClaimException extends $pb.GeneratedMessage {
  factory IDPRejectedClaimException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IDPRejectedClaimException._();

  factory IDPRejectedClaimException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IDPRejectedClaimException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IDPRejectedClaimException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IDPRejectedClaimException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IDPRejectedClaimException copyWith(
          void Function(IDPRejectedClaimException) updates) =>
      super.copyWith((message) => updates(message as IDPRejectedClaimException))
          as IDPRejectedClaimException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IDPRejectedClaimException create() => IDPRejectedClaimException._();
  @$core.override
  IDPRejectedClaimException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IDPRejectedClaimException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IDPRejectedClaimException>(create);
  static IDPRejectedClaimException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidAuthorizationMessageException extends $pb.GeneratedMessage {
  factory InvalidAuthorizationMessageException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidAuthorizationMessageException._();

  factory InvalidAuthorizationMessageException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidAuthorizationMessageException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidAuthorizationMessageException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAuthorizationMessageException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAuthorizationMessageException copyWith(
          void Function(InvalidAuthorizationMessageException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidAuthorizationMessageException))
          as InvalidAuthorizationMessageException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidAuthorizationMessageException create() =>
      InvalidAuthorizationMessageException._();
  @$core.override
  InvalidAuthorizationMessageException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidAuthorizationMessageException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidAuthorizationMessageException>(create);
  static InvalidAuthorizationMessageException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidIdentityTokenException extends $pb.GeneratedMessage {
  factory InvalidIdentityTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidIdentityTokenException._();

  factory InvalidIdentityTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidIdentityTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidIdentityTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidIdentityTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidIdentityTokenException copyWith(
          void Function(InvalidIdentityTokenException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidIdentityTokenException))
          as InvalidIdentityTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidIdentityTokenException create() =>
      InvalidIdentityTokenException._();
  @$core.override
  InvalidIdentityTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidIdentityTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidIdentityTokenException>(create);
  static InvalidIdentityTokenException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class JWTPayloadSizeExceededException extends $pb.GeneratedMessage {
  factory JWTPayloadSizeExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  JWTPayloadSizeExceededException._();

  factory JWTPayloadSizeExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory JWTPayloadSizeExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'JWTPayloadSizeExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JWTPayloadSizeExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JWTPayloadSizeExceededException copyWith(
          void Function(JWTPayloadSizeExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as JWTPayloadSizeExceededException))
          as JWTPayloadSizeExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static JWTPayloadSizeExceededException create() =>
      JWTPayloadSizeExceededException._();
  @$core.override
  JWTPayloadSizeExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static JWTPayloadSizeExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<JWTPayloadSizeExceededException>(
          create);
  static JWTPayloadSizeExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class MalformedPolicyDocumentException extends $pb.GeneratedMessage {
  factory MalformedPolicyDocumentException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  MalformedPolicyDocumentException._();

  factory MalformedPolicyDocumentException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MalformedPolicyDocumentException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MalformedPolicyDocumentException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MalformedPolicyDocumentException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MalformedPolicyDocumentException copyWith(
          void Function(MalformedPolicyDocumentException) updates) =>
      super.copyWith(
              (message) => updates(message as MalformedPolicyDocumentException))
          as MalformedPolicyDocumentException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MalformedPolicyDocumentException create() =>
      MalformedPolicyDocumentException._();
  @$core.override
  MalformedPolicyDocumentException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MalformedPolicyDocumentException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MalformedPolicyDocumentException>(
          create);
  static MalformedPolicyDocumentException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class OutboundWebIdentityFederationDisabledException
    extends $pb.GeneratedMessage {
  factory OutboundWebIdentityFederationDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  OutboundWebIdentityFederationDisabledException._();

  factory OutboundWebIdentityFederationDisabledException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OutboundWebIdentityFederationDisabledException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OutboundWebIdentityFederationDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OutboundWebIdentityFederationDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OutboundWebIdentityFederationDisabledException copyWith(
          void Function(OutboundWebIdentityFederationDisabledException)
              updates) =>
      super.copyWith((message) => updates(
              message as OutboundWebIdentityFederationDisabledException))
          as OutboundWebIdentityFederationDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OutboundWebIdentityFederationDisabledException create() =>
      OutboundWebIdentityFederationDisabledException._();
  @$core.override
  OutboundWebIdentityFederationDisabledException createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static OutboundWebIdentityFederationDisabledException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          OutboundWebIdentityFederationDisabledException>(create);
  static OutboundWebIdentityFederationDisabledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PackedPolicyTooLargeException extends $pb.GeneratedMessage {
  factory PackedPolicyTooLargeException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PackedPolicyTooLargeException._();

  factory PackedPolicyTooLargeException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PackedPolicyTooLargeException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PackedPolicyTooLargeException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PackedPolicyTooLargeException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PackedPolicyTooLargeException copyWith(
          void Function(PackedPolicyTooLargeException) updates) =>
      super.copyWith(
              (message) => updates(message as PackedPolicyTooLargeException))
          as PackedPolicyTooLargeException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PackedPolicyTooLargeException create() =>
      PackedPolicyTooLargeException._();
  @$core.override
  PackedPolicyTooLargeException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PackedPolicyTooLargeException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PackedPolicyTooLargeException>(create);
  static PackedPolicyTooLargeException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PolicyDescriptorType extends $pb.GeneratedMessage {
  factory PolicyDescriptorType({
    $core.String? arn,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    return result;
  }

  PolicyDescriptorType._();

  factory PolicyDescriptorType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PolicyDescriptorType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PolicyDescriptorType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(359604989, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PolicyDescriptorType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PolicyDescriptorType copyWith(void Function(PolicyDescriptorType) updates) =>
      super.copyWith((message) => updates(message as PolicyDescriptorType))
          as PolicyDescriptorType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PolicyDescriptorType create() => PolicyDescriptorType._();
  @$core.override
  PolicyDescriptorType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PolicyDescriptorType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PolicyDescriptorType>(create);
  static PolicyDescriptorType? _defaultInstance;

  @$pb.TagNumber(359604989)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(359604989)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(359604989)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(359604989)
  void clearArn() => $_clearField(359604989);
}

class ProvidedContext extends $pb.GeneratedMessage {
  factory ProvidedContext({
    $core.String? contextassertion,
    $core.String? providerarn,
  }) {
    final result = create();
    if (contextassertion != null) result.contextassertion = contextassertion;
    if (providerarn != null) result.providerarn = providerarn;
    return result;
  }

  ProvidedContext._();

  factory ProvidedContext.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProvidedContext.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProvidedContext',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(351907089, _omitFieldNames ? '' : 'contextassertion')
    ..aOS(426083188, _omitFieldNames ? '' : 'providerarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvidedContext clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvidedContext copyWith(void Function(ProvidedContext) updates) =>
      super.copyWith((message) => updates(message as ProvidedContext))
          as ProvidedContext;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProvidedContext create() => ProvidedContext._();
  @$core.override
  ProvidedContext createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProvidedContext getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProvidedContext>(create);
  static ProvidedContext? _defaultInstance;

  @$pb.TagNumber(351907089)
  $core.String get contextassertion => $_getSZ(0);
  @$pb.TagNumber(351907089)
  set contextassertion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(351907089)
  $core.bool hasContextassertion() => $_has(0);
  @$pb.TagNumber(351907089)
  void clearContextassertion() => $_clearField(351907089);

  @$pb.TagNumber(426083188)
  $core.String get providerarn => $_getSZ(1);
  @$pb.TagNumber(426083188)
  set providerarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(426083188)
  $core.bool hasProviderarn() => $_has(1);
  @$pb.TagNumber(426083188)
  void clearProviderarn() => $_clearField(426083188);
}

class RegionDisabledException extends $pb.GeneratedMessage {
  factory RegionDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  RegionDisabledException._();

  factory RegionDisabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegionDisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegionDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegionDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegionDisabledException copyWith(
          void Function(RegionDisabledException) updates) =>
      super.copyWith((message) => updates(message as RegionDisabledException))
          as RegionDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegionDisabledException create() => RegionDisabledException._();
  @$core.override
  RegionDisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegionDisabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegionDisabledException>(create);
  static RegionDisabledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class SessionDurationEscalationException extends $pb.GeneratedMessage {
  factory SessionDurationEscalationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  SessionDurationEscalationException._();

  factory SessionDurationEscalationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SessionDurationEscalationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SessionDurationEscalationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionDurationEscalationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SessionDurationEscalationException copyWith(
          void Function(SessionDurationEscalationException) updates) =>
      super.copyWith((message) =>
              updates(message as SessionDurationEscalationException))
          as SessionDurationEscalationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SessionDurationEscalationException create() =>
      SessionDurationEscalationException._();
  @$core.override
  SessionDurationEscalationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SessionDurationEscalationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SessionDurationEscalationException>(
          create);
  static SessionDurationEscalationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sts'),
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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
