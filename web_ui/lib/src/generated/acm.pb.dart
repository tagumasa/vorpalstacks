// This is a generated file - do not edit.
//
// Generated from acm.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'acm.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'acm.pbenum.dart';

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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

class AddTagsToCertificateRequest extends $pb.GeneratedMessage {
  factory AddTagsToCertificateRequest({
    $core.String? certificatearn,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  AddTagsToCertificateRequest._();

  factory AddTagsToCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddTagsToCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddTagsToCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsToCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddTagsToCertificateRequest copyWith(
          void Function(AddTagsToCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as AddTagsToCertificateRequest))
          as AddTagsToCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddTagsToCertificateRequest create() =>
      AddTagsToCertificateRequest._();
  @$core.override
  AddTagsToCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddTagsToCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddTagsToCertificateRequest>(create);
  static AddTagsToCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);
}

class CertificateDetail extends $pb.GeneratedMessage {
  factory CertificateDetail({
    CertificateStatus? status,
    $core.String? subject,
    $core.String? issuedat,
    $core.String? revokedat,
    $core.String? certificatearn,
    $core.Iterable<$core.String>? subjectalternativenames,
    RenewalSummary? renewalsummary,
    $core.String? serial,
    RenewalEligibility? renewaleligibility,
    $core.String? domainname,
    FailureReason? failurereason,
    $core.String? createdat,
    $core.String? certificateauthorityarn,
    $core.String? notafter,
    CertificateType? type,
    $core.Iterable<$core.String>? inuseby,
    RevocationReason? revocationreason,
    $core.Iterable<DomainValidation>? domainvalidationoptions,
    $core.Iterable<KeyUsage>? keyusages,
    $core.String? importedat,
    CertificateOptions? options,
    KeyAlgorithm? keyalgorithm,
    CertificateManagedBy? managedby,
    $core.String? notbefore,
    $core.String? signaturealgorithm,
    $core.String? issuer,
    $core.Iterable<ExtendedKeyUsage>? extendedkeyusages,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (subject != null) result.subject = subject;
    if (issuedat != null) result.issuedat = issuedat;
    if (revokedat != null) result.revokedat = revokedat;
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (subjectalternativenames != null)
      result.subjectalternativenames.addAll(subjectalternativenames);
    if (renewalsummary != null) result.renewalsummary = renewalsummary;
    if (serial != null) result.serial = serial;
    if (renewaleligibility != null)
      result.renewaleligibility = renewaleligibility;
    if (domainname != null) result.domainname = domainname;
    if (failurereason != null) result.failurereason = failurereason;
    if (createdat != null) result.createdat = createdat;
    if (certificateauthorityarn != null)
      result.certificateauthorityarn = certificateauthorityarn;
    if (notafter != null) result.notafter = notafter;
    if (type != null) result.type = type;
    if (inuseby != null) result.inuseby.addAll(inuseby);
    if (revocationreason != null) result.revocationreason = revocationreason;
    if (domainvalidationoptions != null)
      result.domainvalidationoptions.addAll(domainvalidationoptions);
    if (keyusages != null) result.keyusages.addAll(keyusages);
    if (importedat != null) result.importedat = importedat;
    if (options != null) result.options = options;
    if (keyalgorithm != null) result.keyalgorithm = keyalgorithm;
    if (managedby != null) result.managedby = managedby;
    if (notbefore != null) result.notbefore = notbefore;
    if (signaturealgorithm != null)
      result.signaturealgorithm = signaturealgorithm;
    if (issuer != null) result.issuer = issuer;
    if (extendedkeyusages != null)
      result.extendedkeyusages.addAll(extendedkeyusages);
    return result;
  }

  CertificateDetail._();

  factory CertificateDetail.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CertificateDetail.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CertificateDetail',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aE<CertificateStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: CertificateStatus.values)
    ..aOS(7939312, _omitFieldNames ? '' : 'subject')
    ..aOS(17449786, _omitFieldNames ? '' : 'issuedat')
    ..aOS(63251417, _omitFieldNames ? '' : 'revokedat')
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..pPS(109998119, _omitFieldNames ? '' : 'subjectalternativenames')
    ..aOM<RenewalSummary>(125255166, _omitFieldNames ? '' : 'renewalsummary',
        subBuilder: RenewalSummary.create)
    ..aOS(143954586, _omitFieldNames ? '' : 'serial')
    ..aE<RenewalEligibility>(
        172871849, _omitFieldNames ? '' : 'renewaleligibility',
        enumValues: RenewalEligibility.values)
    ..aOS(194914027, _omitFieldNames ? '' : 'domainname')
    ..aE<FailureReason>(232322142, _omitFieldNames ? '' : 'failurereason',
        enumValues: FailureReason.values)
    ..aOS(258192751, _omitFieldNames ? '' : 'createdat')
    ..aOS(266069181, _omitFieldNames ? '' : 'certificateauthorityarn')
    ..aOS(287678033, _omitFieldNames ? '' : 'notafter')
    ..aE<CertificateType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: CertificateType.values)
    ..pPS(330307273, _omitFieldNames ? '' : 'inuseby')
    ..aE<RevocationReason>(331836358, _omitFieldNames ? '' : 'revocationreason',
        enumValues: RevocationReason.values)
    ..pPM<DomainValidation>(
        335573705, _omitFieldNames ? '' : 'domainvalidationoptions',
        subBuilder: DomainValidation.create)
    ..pPM<KeyUsage>(345433681, _omitFieldNames ? '' : 'keyusages',
        subBuilder: KeyUsage.create)
    ..aOS(348649225, _omitFieldNames ? '' : 'importedat')
    ..aOM<CertificateOptions>(356388166, _omitFieldNames ? '' : 'options',
        subBuilder: CertificateOptions.create)
    ..aE<KeyAlgorithm>(452317818, _omitFieldNames ? '' : 'keyalgorithm',
        enumValues: KeyAlgorithm.values)
    ..aE<CertificateManagedBy>(455511232, _omitFieldNames ? '' : 'managedby',
        enumValues: CertificateManagedBy.values)
    ..aOS(459074038, _omitFieldNames ? '' : 'notbefore')
    ..aOS(476410739, _omitFieldNames ? '' : 'signaturealgorithm')
    ..aOS(528708823, _omitFieldNames ? '' : 'issuer')
    ..pPM<ExtendedKeyUsage>(
        531267688, _omitFieldNames ? '' : 'extendedkeyusages',
        subBuilder: ExtendedKeyUsage.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateDetail clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateDetail copyWith(void Function(CertificateDetail) updates) =>
      super.copyWith((message) => updates(message as CertificateDetail))
          as CertificateDetail;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CertificateDetail create() => CertificateDetail._();
  @$core.override
  CertificateDetail createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CertificateDetail getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CertificateDetail>(create);
  static CertificateDetail? _defaultInstance;

  @$pb.TagNumber(6222352)
  CertificateStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(CertificateStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(7939312)
  $core.String get subject => $_getSZ(1);
  @$pb.TagNumber(7939312)
  set subject($core.String value) => $_setString(1, value);
  @$pb.TagNumber(7939312)
  $core.bool hasSubject() => $_has(1);
  @$pb.TagNumber(7939312)
  void clearSubject() => $_clearField(7939312);

  @$pb.TagNumber(17449786)
  $core.String get issuedat => $_getSZ(2);
  @$pb.TagNumber(17449786)
  set issuedat($core.String value) => $_setString(2, value);
  @$pb.TagNumber(17449786)
  $core.bool hasIssuedat() => $_has(2);
  @$pb.TagNumber(17449786)
  void clearIssuedat() => $_clearField(17449786);

  @$pb.TagNumber(63251417)
  $core.String get revokedat => $_getSZ(3);
  @$pb.TagNumber(63251417)
  set revokedat($core.String value) => $_setString(3, value);
  @$pb.TagNumber(63251417)
  $core.bool hasRevokedat() => $_has(3);
  @$pb.TagNumber(63251417)
  void clearRevokedat() => $_clearField(63251417);

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(4);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(4);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(109998119)
  $pb.PbList<$core.String> get subjectalternativenames => $_getList(5);

  @$pb.TagNumber(125255166)
  RenewalSummary get renewalsummary => $_getN(6);
  @$pb.TagNumber(125255166)
  set renewalsummary(RenewalSummary value) => $_setField(125255166, value);
  @$pb.TagNumber(125255166)
  $core.bool hasRenewalsummary() => $_has(6);
  @$pb.TagNumber(125255166)
  void clearRenewalsummary() => $_clearField(125255166);
  @$pb.TagNumber(125255166)
  RenewalSummary ensureRenewalsummary() => $_ensure(6);

  @$pb.TagNumber(143954586)
  $core.String get serial => $_getSZ(7);
  @$pb.TagNumber(143954586)
  set serial($core.String value) => $_setString(7, value);
  @$pb.TagNumber(143954586)
  $core.bool hasSerial() => $_has(7);
  @$pb.TagNumber(143954586)
  void clearSerial() => $_clearField(143954586);

  @$pb.TagNumber(172871849)
  RenewalEligibility get renewaleligibility => $_getN(8);
  @$pb.TagNumber(172871849)
  set renewaleligibility(RenewalEligibility value) =>
      $_setField(172871849, value);
  @$pb.TagNumber(172871849)
  $core.bool hasRenewaleligibility() => $_has(8);
  @$pb.TagNumber(172871849)
  void clearRenewaleligibility() => $_clearField(172871849);

  @$pb.TagNumber(194914027)
  $core.String get domainname => $_getSZ(9);
  @$pb.TagNumber(194914027)
  set domainname($core.String value) => $_setString(9, value);
  @$pb.TagNumber(194914027)
  $core.bool hasDomainname() => $_has(9);
  @$pb.TagNumber(194914027)
  void clearDomainname() => $_clearField(194914027);

  @$pb.TagNumber(232322142)
  FailureReason get failurereason => $_getN(10);
  @$pb.TagNumber(232322142)
  set failurereason(FailureReason value) => $_setField(232322142, value);
  @$pb.TagNumber(232322142)
  $core.bool hasFailurereason() => $_has(10);
  @$pb.TagNumber(232322142)
  void clearFailurereason() => $_clearField(232322142);

  @$pb.TagNumber(258192751)
  $core.String get createdat => $_getSZ(11);
  @$pb.TagNumber(258192751)
  set createdat($core.String value) => $_setString(11, value);
  @$pb.TagNumber(258192751)
  $core.bool hasCreatedat() => $_has(11);
  @$pb.TagNumber(258192751)
  void clearCreatedat() => $_clearField(258192751);

  @$pb.TagNumber(266069181)
  $core.String get certificateauthorityarn => $_getSZ(12);
  @$pb.TagNumber(266069181)
  set certificateauthorityarn($core.String value) => $_setString(12, value);
  @$pb.TagNumber(266069181)
  $core.bool hasCertificateauthorityarn() => $_has(12);
  @$pb.TagNumber(266069181)
  void clearCertificateauthorityarn() => $_clearField(266069181);

  @$pb.TagNumber(287678033)
  $core.String get notafter => $_getSZ(13);
  @$pb.TagNumber(287678033)
  set notafter($core.String value) => $_setString(13, value);
  @$pb.TagNumber(287678033)
  $core.bool hasNotafter() => $_has(13);
  @$pb.TagNumber(287678033)
  void clearNotafter() => $_clearField(287678033);

  @$pb.TagNumber(290836590)
  CertificateType get type => $_getN(14);
  @$pb.TagNumber(290836590)
  set type(CertificateType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(14);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(330307273)
  $pb.PbList<$core.String> get inuseby => $_getList(15);

  @$pb.TagNumber(331836358)
  RevocationReason get revocationreason => $_getN(16);
  @$pb.TagNumber(331836358)
  set revocationreason(RevocationReason value) => $_setField(331836358, value);
  @$pb.TagNumber(331836358)
  $core.bool hasRevocationreason() => $_has(16);
  @$pb.TagNumber(331836358)
  void clearRevocationreason() => $_clearField(331836358);

  @$pb.TagNumber(335573705)
  $pb.PbList<DomainValidation> get domainvalidationoptions => $_getList(17);

  @$pb.TagNumber(345433681)
  $pb.PbList<KeyUsage> get keyusages => $_getList(18);

  @$pb.TagNumber(348649225)
  $core.String get importedat => $_getSZ(19);
  @$pb.TagNumber(348649225)
  set importedat($core.String value) => $_setString(19, value);
  @$pb.TagNumber(348649225)
  $core.bool hasImportedat() => $_has(19);
  @$pb.TagNumber(348649225)
  void clearImportedat() => $_clearField(348649225);

  @$pb.TagNumber(356388166)
  CertificateOptions get options => $_getN(20);
  @$pb.TagNumber(356388166)
  set options(CertificateOptions value) => $_setField(356388166, value);
  @$pb.TagNumber(356388166)
  $core.bool hasOptions() => $_has(20);
  @$pb.TagNumber(356388166)
  void clearOptions() => $_clearField(356388166);
  @$pb.TagNumber(356388166)
  CertificateOptions ensureOptions() => $_ensure(20);

  @$pb.TagNumber(452317818)
  KeyAlgorithm get keyalgorithm => $_getN(21);
  @$pb.TagNumber(452317818)
  set keyalgorithm(KeyAlgorithm value) => $_setField(452317818, value);
  @$pb.TagNumber(452317818)
  $core.bool hasKeyalgorithm() => $_has(21);
  @$pb.TagNumber(452317818)
  void clearKeyalgorithm() => $_clearField(452317818);

  @$pb.TagNumber(455511232)
  CertificateManagedBy get managedby => $_getN(22);
  @$pb.TagNumber(455511232)
  set managedby(CertificateManagedBy value) => $_setField(455511232, value);
  @$pb.TagNumber(455511232)
  $core.bool hasManagedby() => $_has(22);
  @$pb.TagNumber(455511232)
  void clearManagedby() => $_clearField(455511232);

  @$pb.TagNumber(459074038)
  $core.String get notbefore => $_getSZ(23);
  @$pb.TagNumber(459074038)
  set notbefore($core.String value) => $_setString(23, value);
  @$pb.TagNumber(459074038)
  $core.bool hasNotbefore() => $_has(23);
  @$pb.TagNumber(459074038)
  void clearNotbefore() => $_clearField(459074038);

  @$pb.TagNumber(476410739)
  $core.String get signaturealgorithm => $_getSZ(24);
  @$pb.TagNumber(476410739)
  set signaturealgorithm($core.String value) => $_setString(24, value);
  @$pb.TagNumber(476410739)
  $core.bool hasSignaturealgorithm() => $_has(24);
  @$pb.TagNumber(476410739)
  void clearSignaturealgorithm() => $_clearField(476410739);

  @$pb.TagNumber(528708823)
  $core.String get issuer => $_getSZ(25);
  @$pb.TagNumber(528708823)
  set issuer($core.String value) => $_setString(25, value);
  @$pb.TagNumber(528708823)
  $core.bool hasIssuer() => $_has(25);
  @$pb.TagNumber(528708823)
  void clearIssuer() => $_clearField(528708823);

  @$pb.TagNumber(531267688)
  $pb.PbList<ExtendedKeyUsage> get extendedkeyusages => $_getList(26);
}

class CertificateOptions extends $pb.GeneratedMessage {
  factory CertificateOptions({
    CertificateExport? export,
    CertificateTransparencyLoggingPreference?
        certificatetransparencyloggingpreference,
  }) {
    final result = create();
    if (export != null) result.export = export;
    if (certificatetransparencyloggingpreference != null)
      result.certificatetransparencyloggingpreference =
          certificatetransparencyloggingpreference;
    return result;
  }

  CertificateOptions._();

  factory CertificateOptions.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CertificateOptions.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CertificateOptions',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aE<CertificateExport>(140724692, _omitFieldNames ? '' : 'export',
        enumValues: CertificateExport.values)
    ..aE<CertificateTransparencyLoggingPreference>(414636075,
        _omitFieldNames ? '' : 'certificatetransparencyloggingpreference',
        enumValues: CertificateTransparencyLoggingPreference.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateOptions clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateOptions copyWith(void Function(CertificateOptions) updates) =>
      super.copyWith((message) => updates(message as CertificateOptions))
          as CertificateOptions;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CertificateOptions create() => CertificateOptions._();
  @$core.override
  CertificateOptions createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CertificateOptions getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CertificateOptions>(create);
  static CertificateOptions? _defaultInstance;

  @$pb.TagNumber(140724692)
  CertificateExport get export => $_getN(0);
  @$pb.TagNumber(140724692)
  set export(CertificateExport value) => $_setField(140724692, value);
  @$pb.TagNumber(140724692)
  $core.bool hasExport() => $_has(0);
  @$pb.TagNumber(140724692)
  void clearExport() => $_clearField(140724692);

  @$pb.TagNumber(414636075)
  CertificateTransparencyLoggingPreference
      get certificatetransparencyloggingpreference => $_getN(1);
  @$pb.TagNumber(414636075)
  set certificatetransparencyloggingpreference(
          CertificateTransparencyLoggingPreference value) =>
      $_setField(414636075, value);
  @$pb.TagNumber(414636075)
  $core.bool hasCertificatetransparencyloggingpreference() => $_has(1);
  @$pb.TagNumber(414636075)
  void clearCertificatetransparencyloggingpreference() =>
      $_clearField(414636075);
}

class CertificateSummary extends $pb.GeneratedMessage {
  factory CertificateSummary({
    CertificateStatus? status,
    $core.String? issuedat,
    CertificateExport? exportoption,
    $core.String? revokedat,
    $core.String? certificatearn,
    RenewalEligibility? renewaleligibility,
    $core.String? domainname,
    $core.Iterable<$core.String>? subjectalternativenamesummaries,
    $core.String? createdat,
    $core.String? notafter,
    CertificateType? type,
    $core.Iterable<KeyUsageName>? keyusages,
    $core.String? importedat,
    $core.bool? hasadditionalsubjectalternativenames,
    $core.bool? inuse,
    KeyAlgorithm? keyalgorithm,
    CertificateManagedBy? managedby,
    $core.String? notbefore,
    $core.bool? exported,
    $core.Iterable<ExtendedKeyUsageName>? extendedkeyusages,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (issuedat != null) result.issuedat = issuedat;
    if (exportoption != null) result.exportoption = exportoption;
    if (revokedat != null) result.revokedat = revokedat;
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (renewaleligibility != null)
      result.renewaleligibility = renewaleligibility;
    if (domainname != null) result.domainname = domainname;
    if (subjectalternativenamesummaries != null)
      result.subjectalternativenamesummaries
          .addAll(subjectalternativenamesummaries);
    if (createdat != null) result.createdat = createdat;
    if (notafter != null) result.notafter = notafter;
    if (type != null) result.type = type;
    if (keyusages != null) result.keyusages.addAll(keyusages);
    if (importedat != null) result.importedat = importedat;
    if (hasadditionalsubjectalternativenames != null)
      result.hasadditionalsubjectalternativenames =
          hasadditionalsubjectalternativenames;
    if (inuse != null) result.inuse = inuse;
    if (keyalgorithm != null) result.keyalgorithm = keyalgorithm;
    if (managedby != null) result.managedby = managedby;
    if (notbefore != null) result.notbefore = notbefore;
    if (exported != null) result.exported = exported;
    if (extendedkeyusages != null)
      result.extendedkeyusages.addAll(extendedkeyusages);
    return result;
  }

  CertificateSummary._();

  factory CertificateSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CertificateSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CertificateSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aE<CertificateStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: CertificateStatus.values)
    ..aOS(17449786, _omitFieldNames ? '' : 'issuedat')
    ..aE<CertificateExport>(19500687, _omitFieldNames ? '' : 'exportoption',
        enumValues: CertificateExport.values)
    ..aOS(63251417, _omitFieldNames ? '' : 'revokedat')
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..aE<RenewalEligibility>(
        172871849, _omitFieldNames ? '' : 'renewaleligibility',
        enumValues: RenewalEligibility.values)
    ..aOS(194914027, _omitFieldNames ? '' : 'domainname')
    ..pPS(248841448, _omitFieldNames ? '' : 'subjectalternativenamesummaries')
    ..aOS(258192751, _omitFieldNames ? '' : 'createdat')
    ..aOS(287678033, _omitFieldNames ? '' : 'notafter')
    ..aE<CertificateType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: CertificateType.values)
    ..pc<KeyUsageName>(
        345433681, _omitFieldNames ? '' : 'keyusages', $pb.PbFieldType.KE,
        valueOf: KeyUsageName.valueOf,
        enumValues: KeyUsageName.values,
        defaultEnumValue: KeyUsageName.KEY_USAGE_NAME_ANY)
    ..aOS(348649225, _omitFieldNames ? '' : 'importedat')
    ..aOB(389669028,
        _omitFieldNames ? '' : 'hasadditionalsubjectalternativenames')
    ..aOB(398346234, _omitFieldNames ? '' : 'inuse')
    ..aE<KeyAlgorithm>(452317818, _omitFieldNames ? '' : 'keyalgorithm',
        enumValues: KeyAlgorithm.values)
    ..aE<CertificateManagedBy>(455511232, _omitFieldNames ? '' : 'managedby',
        enumValues: CertificateManagedBy.values)
    ..aOS(459074038, _omitFieldNames ? '' : 'notbefore')
    ..aOB(491164947, _omitFieldNames ? '' : 'exported')
    ..pc<ExtendedKeyUsageName>(531267688,
        _omitFieldNames ? '' : 'extendedkeyusages', $pb.PbFieldType.KE,
        valueOf: ExtendedKeyUsageName.valueOf,
        enumValues: ExtendedKeyUsageName.values,
        defaultEnumValue: ExtendedKeyUsageName.EXTENDED_KEY_USAGE_NAME_ANY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateSummary copyWith(void Function(CertificateSummary) updates) =>
      super.copyWith((message) => updates(message as CertificateSummary))
          as CertificateSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CertificateSummary create() => CertificateSummary._();
  @$core.override
  CertificateSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CertificateSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CertificateSummary>(create);
  static CertificateSummary? _defaultInstance;

  @$pb.TagNumber(6222352)
  CertificateStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(CertificateStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(17449786)
  $core.String get issuedat => $_getSZ(1);
  @$pb.TagNumber(17449786)
  set issuedat($core.String value) => $_setString(1, value);
  @$pb.TagNumber(17449786)
  $core.bool hasIssuedat() => $_has(1);
  @$pb.TagNumber(17449786)
  void clearIssuedat() => $_clearField(17449786);

  @$pb.TagNumber(19500687)
  CertificateExport get exportoption => $_getN(2);
  @$pb.TagNumber(19500687)
  set exportoption(CertificateExport value) => $_setField(19500687, value);
  @$pb.TagNumber(19500687)
  $core.bool hasExportoption() => $_has(2);
  @$pb.TagNumber(19500687)
  void clearExportoption() => $_clearField(19500687);

  @$pb.TagNumber(63251417)
  $core.String get revokedat => $_getSZ(3);
  @$pb.TagNumber(63251417)
  set revokedat($core.String value) => $_setString(3, value);
  @$pb.TagNumber(63251417)
  $core.bool hasRevokedat() => $_has(3);
  @$pb.TagNumber(63251417)
  void clearRevokedat() => $_clearField(63251417);

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(4);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(4);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(172871849)
  RenewalEligibility get renewaleligibility => $_getN(5);
  @$pb.TagNumber(172871849)
  set renewaleligibility(RenewalEligibility value) =>
      $_setField(172871849, value);
  @$pb.TagNumber(172871849)
  $core.bool hasRenewaleligibility() => $_has(5);
  @$pb.TagNumber(172871849)
  void clearRenewaleligibility() => $_clearField(172871849);

  @$pb.TagNumber(194914027)
  $core.String get domainname => $_getSZ(6);
  @$pb.TagNumber(194914027)
  set domainname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(194914027)
  $core.bool hasDomainname() => $_has(6);
  @$pb.TagNumber(194914027)
  void clearDomainname() => $_clearField(194914027);

  @$pb.TagNumber(248841448)
  $pb.PbList<$core.String> get subjectalternativenamesummaries => $_getList(7);

  @$pb.TagNumber(258192751)
  $core.String get createdat => $_getSZ(8);
  @$pb.TagNumber(258192751)
  set createdat($core.String value) => $_setString(8, value);
  @$pb.TagNumber(258192751)
  $core.bool hasCreatedat() => $_has(8);
  @$pb.TagNumber(258192751)
  void clearCreatedat() => $_clearField(258192751);

  @$pb.TagNumber(287678033)
  $core.String get notafter => $_getSZ(9);
  @$pb.TagNumber(287678033)
  set notafter($core.String value) => $_setString(9, value);
  @$pb.TagNumber(287678033)
  $core.bool hasNotafter() => $_has(9);
  @$pb.TagNumber(287678033)
  void clearNotafter() => $_clearField(287678033);

  @$pb.TagNumber(290836590)
  CertificateType get type => $_getN(10);
  @$pb.TagNumber(290836590)
  set type(CertificateType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(10);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(345433681)
  $pb.PbList<KeyUsageName> get keyusages => $_getList(11);

  @$pb.TagNumber(348649225)
  $core.String get importedat => $_getSZ(12);
  @$pb.TagNumber(348649225)
  set importedat($core.String value) => $_setString(12, value);
  @$pb.TagNumber(348649225)
  $core.bool hasImportedat() => $_has(12);
  @$pb.TagNumber(348649225)
  void clearImportedat() => $_clearField(348649225);

  @$pb.TagNumber(389669028)
  $core.bool get hasadditionalsubjectalternativenames => $_getBF(13);
  @$pb.TagNumber(389669028)
  set hasadditionalsubjectalternativenames($core.bool value) =>
      $_setBool(13, value);
  @$pb.TagNumber(389669028)
  $core.bool hasHasadditionalsubjectalternativenames() => $_has(13);
  @$pb.TagNumber(389669028)
  void clearHasadditionalsubjectalternativenames() => $_clearField(389669028);

  @$pb.TagNumber(398346234)
  $core.bool get inuse => $_getBF(14);
  @$pb.TagNumber(398346234)
  set inuse($core.bool value) => $_setBool(14, value);
  @$pb.TagNumber(398346234)
  $core.bool hasInuse() => $_has(14);
  @$pb.TagNumber(398346234)
  void clearInuse() => $_clearField(398346234);

  @$pb.TagNumber(452317818)
  KeyAlgorithm get keyalgorithm => $_getN(15);
  @$pb.TagNumber(452317818)
  set keyalgorithm(KeyAlgorithm value) => $_setField(452317818, value);
  @$pb.TagNumber(452317818)
  $core.bool hasKeyalgorithm() => $_has(15);
  @$pb.TagNumber(452317818)
  void clearKeyalgorithm() => $_clearField(452317818);

  @$pb.TagNumber(455511232)
  CertificateManagedBy get managedby => $_getN(16);
  @$pb.TagNumber(455511232)
  set managedby(CertificateManagedBy value) => $_setField(455511232, value);
  @$pb.TagNumber(455511232)
  $core.bool hasManagedby() => $_has(16);
  @$pb.TagNumber(455511232)
  void clearManagedby() => $_clearField(455511232);

  @$pb.TagNumber(459074038)
  $core.String get notbefore => $_getSZ(17);
  @$pb.TagNumber(459074038)
  set notbefore($core.String value) => $_setString(17, value);
  @$pb.TagNumber(459074038)
  $core.bool hasNotbefore() => $_has(17);
  @$pb.TagNumber(459074038)
  void clearNotbefore() => $_clearField(459074038);

  @$pb.TagNumber(491164947)
  $core.bool get exported => $_getBF(18);
  @$pb.TagNumber(491164947)
  set exported($core.bool value) => $_setBool(18, value);
  @$pb.TagNumber(491164947)
  $core.bool hasExported() => $_has(18);
  @$pb.TagNumber(491164947)
  void clearExported() => $_clearField(491164947);

  @$pb.TagNumber(531267688)
  $pb.PbList<ExtendedKeyUsageName> get extendedkeyusages => $_getList(19);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DeleteCertificateRequest extends $pb.GeneratedMessage {
  factory DeleteCertificateRequest({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  DeleteCertificateRequest._();

  factory DeleteCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCertificateRequest copyWith(
          void Function(DeleteCertificateRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteCertificateRequest))
          as DeleteCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteCertificateRequest create() => DeleteCertificateRequest._();
  @$core.override
  DeleteCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteCertificateRequest>(create);
  static DeleteCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
}

class DescribeCertificateRequest extends $pb.GeneratedMessage {
  factory DescribeCertificateRequest({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  DescribeCertificateRequest._();

  factory DescribeCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCertificateRequest copyWith(
          void Function(DescribeCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeCertificateRequest))
          as DescribeCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeCertificateRequest create() => DescribeCertificateRequest._();
  @$core.override
  DescribeCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeCertificateRequest>(create);
  static DescribeCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
}

class DescribeCertificateResponse extends $pb.GeneratedMessage {
  factory DescribeCertificateResponse({
    CertificateDetail? certificate,
  }) {
    final result = create();
    if (certificate != null) result.certificate = certificate;
    return result;
  }

  DescribeCertificateResponse._();

  factory DescribeCertificateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeCertificateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeCertificateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOM<CertificateDetail>(198060817, _omitFieldNames ? '' : 'certificate',
        subBuilder: CertificateDetail.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCertificateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCertificateResponse copyWith(
          void Function(DescribeCertificateResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeCertificateResponse))
          as DescribeCertificateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeCertificateResponse create() =>
      DescribeCertificateResponse._();
  @$core.override
  DescribeCertificateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeCertificateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeCertificateResponse>(create);
  static DescribeCertificateResponse? _defaultInstance;

  @$pb.TagNumber(198060817)
  CertificateDetail get certificate => $_getN(0);
  @$pb.TagNumber(198060817)
  set certificate(CertificateDetail value) => $_setField(198060817, value);
  @$pb.TagNumber(198060817)
  $core.bool hasCertificate() => $_has(0);
  @$pb.TagNumber(198060817)
  void clearCertificate() => $_clearField(198060817);
  @$pb.TagNumber(198060817)
  CertificateDetail ensureCertificate() => $_ensure(0);
}

class DomainValidation extends $pb.GeneratedMessage {
  factory DomainValidation({
    ValidationMethod? validationmethod,
    $core.String? validationdomain,
    $core.String? domainname,
    ResourceRecord? resourcerecord,
    HttpRedirect? httpredirect,
    $core.Iterable<$core.String>? validationemails,
    DomainStatus? validationstatus,
  }) {
    final result = create();
    if (validationmethod != null) result.validationmethod = validationmethod;
    if (validationdomain != null) result.validationdomain = validationdomain;
    if (domainname != null) result.domainname = domainname;
    if (resourcerecord != null) result.resourcerecord = resourcerecord;
    if (httpredirect != null) result.httpredirect = httpredirect;
    if (validationemails != null)
      result.validationemails.addAll(validationemails);
    if (validationstatus != null) result.validationstatus = validationstatus;
    return result;
  }

  DomainValidation._();

  factory DomainValidation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DomainValidation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DomainValidation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aE<ValidationMethod>(58092520, _omitFieldNames ? '' : 'validationmethod',
        enumValues: ValidationMethod.values)
    ..aOS(161300799, _omitFieldNames ? '' : 'validationdomain')
    ..aOS(194914027, _omitFieldNames ? '' : 'domainname')
    ..aOM<ResourceRecord>(207153213, _omitFieldNames ? '' : 'resourcerecord',
        subBuilder: ResourceRecord.create)
    ..aOM<HttpRedirect>(375515090, _omitFieldNames ? '' : 'httpredirect',
        subBuilder: HttpRedirect.create)
    ..pPS(375615664, _omitFieldNames ? '' : 'validationemails')
    ..aE<DomainStatus>(426907749, _omitFieldNames ? '' : 'validationstatus',
        enumValues: DomainStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainValidation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainValidation copyWith(void Function(DomainValidation) updates) =>
      super.copyWith((message) => updates(message as DomainValidation))
          as DomainValidation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DomainValidation create() => DomainValidation._();
  @$core.override
  DomainValidation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DomainValidation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DomainValidation>(create);
  static DomainValidation? _defaultInstance;

  @$pb.TagNumber(58092520)
  ValidationMethod get validationmethod => $_getN(0);
  @$pb.TagNumber(58092520)
  set validationmethod(ValidationMethod value) => $_setField(58092520, value);
  @$pb.TagNumber(58092520)
  $core.bool hasValidationmethod() => $_has(0);
  @$pb.TagNumber(58092520)
  void clearValidationmethod() => $_clearField(58092520);

  @$pb.TagNumber(161300799)
  $core.String get validationdomain => $_getSZ(1);
  @$pb.TagNumber(161300799)
  set validationdomain($core.String value) => $_setString(1, value);
  @$pb.TagNumber(161300799)
  $core.bool hasValidationdomain() => $_has(1);
  @$pb.TagNumber(161300799)
  void clearValidationdomain() => $_clearField(161300799);

  @$pb.TagNumber(194914027)
  $core.String get domainname => $_getSZ(2);
  @$pb.TagNumber(194914027)
  set domainname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(194914027)
  $core.bool hasDomainname() => $_has(2);
  @$pb.TagNumber(194914027)
  void clearDomainname() => $_clearField(194914027);

  @$pb.TagNumber(207153213)
  ResourceRecord get resourcerecord => $_getN(3);
  @$pb.TagNumber(207153213)
  set resourcerecord(ResourceRecord value) => $_setField(207153213, value);
  @$pb.TagNumber(207153213)
  $core.bool hasResourcerecord() => $_has(3);
  @$pb.TagNumber(207153213)
  void clearResourcerecord() => $_clearField(207153213);
  @$pb.TagNumber(207153213)
  ResourceRecord ensureResourcerecord() => $_ensure(3);

  @$pb.TagNumber(375515090)
  HttpRedirect get httpredirect => $_getN(4);
  @$pb.TagNumber(375515090)
  set httpredirect(HttpRedirect value) => $_setField(375515090, value);
  @$pb.TagNumber(375515090)
  $core.bool hasHttpredirect() => $_has(4);
  @$pb.TagNumber(375515090)
  void clearHttpredirect() => $_clearField(375515090);
  @$pb.TagNumber(375515090)
  HttpRedirect ensureHttpredirect() => $_ensure(4);

  @$pb.TagNumber(375615664)
  $pb.PbList<$core.String> get validationemails => $_getList(5);

  @$pb.TagNumber(426907749)
  DomainStatus get validationstatus => $_getN(6);
  @$pb.TagNumber(426907749)
  set validationstatus(DomainStatus value) => $_setField(426907749, value);
  @$pb.TagNumber(426907749)
  $core.bool hasValidationstatus() => $_has(6);
  @$pb.TagNumber(426907749)
  void clearValidationstatus() => $_clearField(426907749);
}

class DomainValidationOption extends $pb.GeneratedMessage {
  factory DomainValidationOption({
    $core.String? validationdomain,
    $core.String? domainname,
  }) {
    final result = create();
    if (validationdomain != null) result.validationdomain = validationdomain;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  DomainValidationOption._();

  factory DomainValidationOption.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DomainValidationOption.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DomainValidationOption',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(161300799, _omitFieldNames ? '' : 'validationdomain')
    ..aOS(194914027, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainValidationOption clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainValidationOption copyWith(
          void Function(DomainValidationOption) updates) =>
      super.copyWith((message) => updates(message as DomainValidationOption))
          as DomainValidationOption;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DomainValidationOption create() => DomainValidationOption._();
  @$core.override
  DomainValidationOption createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DomainValidationOption getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DomainValidationOption>(create);
  static DomainValidationOption? _defaultInstance;

  @$pb.TagNumber(161300799)
  $core.String get validationdomain => $_getSZ(0);
  @$pb.TagNumber(161300799)
  set validationdomain($core.String value) => $_setString(0, value);
  @$pb.TagNumber(161300799)
  $core.bool hasValidationdomain() => $_has(0);
  @$pb.TagNumber(161300799)
  void clearValidationdomain() => $_clearField(161300799);

  @$pb.TagNumber(194914027)
  $core.String get domainname => $_getSZ(1);
  @$pb.TagNumber(194914027)
  set domainname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(194914027)
  $core.bool hasDomainname() => $_has(1);
  @$pb.TagNumber(194914027)
  void clearDomainname() => $_clearField(194914027);
}

class ExpiryEventsConfiguration extends $pb.GeneratedMessage {
  factory ExpiryEventsConfiguration({
    $core.int? daysbeforeexpiry,
  }) {
    final result = create();
    if (daysbeforeexpiry != null) result.daysbeforeexpiry = daysbeforeexpiry;
    return result;
  }

  ExpiryEventsConfiguration._();

  factory ExpiryEventsConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExpiryEventsConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExpiryEventsConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aI(292642703, _omitFieldNames ? '' : 'daysbeforeexpiry')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiryEventsConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiryEventsConfiguration copyWith(
          void Function(ExpiryEventsConfiguration) updates) =>
      super.copyWith((message) => updates(message as ExpiryEventsConfiguration))
          as ExpiryEventsConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExpiryEventsConfiguration create() => ExpiryEventsConfiguration._();
  @$core.override
  ExpiryEventsConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExpiryEventsConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExpiryEventsConfiguration>(create);
  static ExpiryEventsConfiguration? _defaultInstance;

  @$pb.TagNumber(292642703)
  $core.int get daysbeforeexpiry => $_getIZ(0);
  @$pb.TagNumber(292642703)
  set daysbeforeexpiry($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(292642703)
  $core.bool hasDaysbeforeexpiry() => $_has(0);
  @$pb.TagNumber(292642703)
  void clearDaysbeforeexpiry() => $_clearField(292642703);
}

class ExportCertificateRequest extends $pb.GeneratedMessage {
  factory ExportCertificateRequest({
    $core.String? certificatearn,
    $core.List<$core.int>? passphrase,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (passphrase != null) result.passphrase = passphrase;
    return result;
  }

  ExportCertificateRequest._();

  factory ExportCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..a<$core.List<$core.int>>(
        421146054, _omitFieldNames ? '' : 'passphrase', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportCertificateRequest copyWith(
          void Function(ExportCertificateRequest) updates) =>
      super.copyWith((message) => updates(message as ExportCertificateRequest))
          as ExportCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportCertificateRequest create() => ExportCertificateRequest._();
  @$core.override
  ExportCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportCertificateRequest>(create);
  static ExportCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(421146054)
  $core.List<$core.int> get passphrase => $_getN(1);
  @$pb.TagNumber(421146054)
  set passphrase($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(421146054)
  $core.bool hasPassphrase() => $_has(1);
  @$pb.TagNumber(421146054)
  void clearPassphrase() => $_clearField(421146054);
}

class ExportCertificateResponse extends $pb.GeneratedMessage {
  factory ExportCertificateResponse({
    $core.String? privatekey,
    $core.String? certificate,
    $core.String? certificatechain,
  }) {
    final result = create();
    if (privatekey != null) result.privatekey = privatekey;
    if (certificate != null) result.certificate = certificate;
    if (certificatechain != null) result.certificatechain = certificatechain;
    return result;
  }

  ExportCertificateResponse._();

  factory ExportCertificateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportCertificateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportCertificateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(173312700, _omitFieldNames ? '' : 'privatekey')
    ..aOS(198060817, _omitFieldNames ? '' : 'certificate')
    ..aOS(369292378, _omitFieldNames ? '' : 'certificatechain')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportCertificateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportCertificateResponse copyWith(
          void Function(ExportCertificateResponse) updates) =>
      super.copyWith((message) => updates(message as ExportCertificateResponse))
          as ExportCertificateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportCertificateResponse create() => ExportCertificateResponse._();
  @$core.override
  ExportCertificateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportCertificateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportCertificateResponse>(create);
  static ExportCertificateResponse? _defaultInstance;

  @$pb.TagNumber(173312700)
  $core.String get privatekey => $_getSZ(0);
  @$pb.TagNumber(173312700)
  set privatekey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(173312700)
  $core.bool hasPrivatekey() => $_has(0);
  @$pb.TagNumber(173312700)
  void clearPrivatekey() => $_clearField(173312700);

  @$pb.TagNumber(198060817)
  $core.String get certificate => $_getSZ(1);
  @$pb.TagNumber(198060817)
  set certificate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(198060817)
  $core.bool hasCertificate() => $_has(1);
  @$pb.TagNumber(198060817)
  void clearCertificate() => $_clearField(198060817);

  @$pb.TagNumber(369292378)
  $core.String get certificatechain => $_getSZ(2);
  @$pb.TagNumber(369292378)
  set certificatechain($core.String value) => $_setString(2, value);
  @$pb.TagNumber(369292378)
  $core.bool hasCertificatechain() => $_has(2);
  @$pb.TagNumber(369292378)
  void clearCertificatechain() => $_clearField(369292378);
}

class ExtendedKeyUsage extends $pb.GeneratedMessage {
  factory ExtendedKeyUsage({
    ExtendedKeyUsageName? name,
    $core.String? oid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (oid != null) result.oid = oid;
    return result;
  }

  ExtendedKeyUsage._();

  factory ExtendedKeyUsage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExtendedKeyUsage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExtendedKeyUsage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aE<ExtendedKeyUsageName>(266367751, _omitFieldNames ? '' : 'name',
        enumValues: ExtendedKeyUsageName.values)
    ..aOS(504527812, _omitFieldNames ? '' : 'oid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExtendedKeyUsage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExtendedKeyUsage copyWith(void Function(ExtendedKeyUsage) updates) =>
      super.copyWith((message) => updates(message as ExtendedKeyUsage))
          as ExtendedKeyUsage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExtendedKeyUsage create() => ExtendedKeyUsage._();
  @$core.override
  ExtendedKeyUsage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExtendedKeyUsage getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExtendedKeyUsage>(create);
  static ExtendedKeyUsage? _defaultInstance;

  @$pb.TagNumber(266367751)
  ExtendedKeyUsageName get name => $_getN(0);
  @$pb.TagNumber(266367751)
  set name(ExtendedKeyUsageName value) => $_setField(266367751, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(504527812)
  $core.String get oid => $_getSZ(1);
  @$pb.TagNumber(504527812)
  set oid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(504527812)
  $core.bool hasOid() => $_has(1);
  @$pb.TagNumber(504527812)
  void clearOid() => $_clearField(504527812);
}

class Filters extends $pb.GeneratedMessage {
  factory Filters({
    $core.Iterable<KeyAlgorithm>? keytypes,
    CertificateManagedBy? managedby,
    $core.Iterable<ExtendedKeyUsageName>? extendedkeyusage,
    CertificateExport? exportoption,
    $core.Iterable<KeyUsageName>? keyusage,
  }) {
    final result = create();
    if (keytypes != null) result.keytypes.addAll(keytypes);
    if (managedby != null) result.managedby = managedby;
    if (extendedkeyusage != null)
      result.extendedkeyusage.addAll(extendedkeyusage);
    if (exportoption != null) result.exportoption = exportoption;
    if (keyusage != null) result.keyusage.addAll(keyusage);
    return result;
  }

  Filters._();

  factory Filters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Filters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Filters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..pc<KeyAlgorithm>(
        27369686, _omitFieldNames ? '' : 'keytypes', $pb.PbFieldType.KE,
        valueOf: KeyAlgorithm.valueOf,
        enumValues: KeyAlgorithm.values,
        defaultEnumValue: KeyAlgorithm.KEY_ALGORITHM_EC_PRIME256V1)
    ..aE<CertificateManagedBy>(401567328, _omitFieldNames ? '' : 'managedby',
        enumValues: CertificateManagedBy.values)
    ..pc<ExtendedKeyUsageName>(420785679,
        _omitFieldNames ? '' : 'extendedkeyusage', $pb.PbFieldType.KE,
        valueOf: ExtendedKeyUsageName.valueOf,
        enumValues: ExtendedKeyUsageName.values,
        defaultEnumValue: ExtendedKeyUsageName.EXTENDED_KEY_USAGE_NAME_ANY)
    ..aE<CertificateExport>(442051119, _omitFieldNames ? '' : 'exportoption',
        enumValues: CertificateExport.values)
    ..pc<KeyUsageName>(
        502227108, _omitFieldNames ? '' : 'keyusage', $pb.PbFieldType.KE,
        valueOf: KeyUsageName.valueOf,
        enumValues: KeyUsageName.values,
        defaultEnumValue: KeyUsageName.KEY_USAGE_NAME_ANY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Filters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Filters copyWith(void Function(Filters) updates) =>
      super.copyWith((message) => updates(message as Filters)) as Filters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Filters create() => Filters._();
  @$core.override
  Filters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Filters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Filters>(create);
  static Filters? _defaultInstance;

  @$pb.TagNumber(27369686)
  $pb.PbList<KeyAlgorithm> get keytypes => $_getList(0);

  @$pb.TagNumber(401567328)
  CertificateManagedBy get managedby => $_getN(1);
  @$pb.TagNumber(401567328)
  set managedby(CertificateManagedBy value) => $_setField(401567328, value);
  @$pb.TagNumber(401567328)
  $core.bool hasManagedby() => $_has(1);
  @$pb.TagNumber(401567328)
  void clearManagedby() => $_clearField(401567328);

  @$pb.TagNumber(420785679)
  $pb.PbList<ExtendedKeyUsageName> get extendedkeyusage => $_getList(2);

  @$pb.TagNumber(442051119)
  CertificateExport get exportoption => $_getN(3);
  @$pb.TagNumber(442051119)
  set exportoption(CertificateExport value) => $_setField(442051119, value);
  @$pb.TagNumber(442051119)
  $core.bool hasExportoption() => $_has(3);
  @$pb.TagNumber(442051119)
  void clearExportoption() => $_clearField(442051119);

  @$pb.TagNumber(502227108)
  $pb.PbList<KeyUsageName> get keyusage => $_getList(4);
}

class GetAccountConfigurationResponse extends $pb.GeneratedMessage {
  factory GetAccountConfigurationResponse({
    ExpiryEventsConfiguration? expiryevents,
  }) {
    final result = create();
    if (expiryevents != null) result.expiryevents = expiryevents;
    return result;
  }

  GetAccountConfigurationResponse._();

  factory GetAccountConfigurationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccountConfigurationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccountConfigurationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOM<ExpiryEventsConfiguration>(
        358763356, _omitFieldNames ? '' : 'expiryevents',
        subBuilder: ExpiryEventsConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountConfigurationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountConfigurationResponse copyWith(
          void Function(GetAccountConfigurationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetAccountConfigurationResponse))
          as GetAccountConfigurationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountConfigurationResponse create() =>
      GetAccountConfigurationResponse._();
  @$core.override
  GetAccountConfigurationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccountConfigurationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccountConfigurationResponse>(
          create);
  static GetAccountConfigurationResponse? _defaultInstance;

  @$pb.TagNumber(358763356)
  ExpiryEventsConfiguration get expiryevents => $_getN(0);
  @$pb.TagNumber(358763356)
  set expiryevents(ExpiryEventsConfiguration value) =>
      $_setField(358763356, value);
  @$pb.TagNumber(358763356)
  $core.bool hasExpiryevents() => $_has(0);
  @$pb.TagNumber(358763356)
  void clearExpiryevents() => $_clearField(358763356);
  @$pb.TagNumber(358763356)
  ExpiryEventsConfiguration ensureExpiryevents() => $_ensure(0);
}

class GetCertificateRequest extends $pb.GeneratedMessage {
  factory GetCertificateRequest({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  GetCertificateRequest._();

  factory GetCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCertificateRequest copyWith(
          void Function(GetCertificateRequest) updates) =>
      super.copyWith((message) => updates(message as GetCertificateRequest))
          as GetCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCertificateRequest create() => GetCertificateRequest._();
  @$core.override
  GetCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCertificateRequest>(create);
  static GetCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
}

class GetCertificateResponse extends $pb.GeneratedMessage {
  factory GetCertificateResponse({
    $core.String? certificate,
    $core.String? certificatechain,
  }) {
    final result = create();
    if (certificate != null) result.certificate = certificate;
    if (certificatechain != null) result.certificatechain = certificatechain;
    return result;
  }

  GetCertificateResponse._();

  factory GetCertificateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCertificateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCertificateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(198060817, _omitFieldNames ? '' : 'certificate')
    ..aOS(369292378, _omitFieldNames ? '' : 'certificatechain')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCertificateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCertificateResponse copyWith(
          void Function(GetCertificateResponse) updates) =>
      super.copyWith((message) => updates(message as GetCertificateResponse))
          as GetCertificateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCertificateResponse create() => GetCertificateResponse._();
  @$core.override
  GetCertificateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCertificateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCertificateResponse>(create);
  static GetCertificateResponse? _defaultInstance;

  @$pb.TagNumber(198060817)
  $core.String get certificate => $_getSZ(0);
  @$pb.TagNumber(198060817)
  set certificate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(198060817)
  $core.bool hasCertificate() => $_has(0);
  @$pb.TagNumber(198060817)
  void clearCertificate() => $_clearField(198060817);

  @$pb.TagNumber(369292378)
  $core.String get certificatechain => $_getSZ(1);
  @$pb.TagNumber(369292378)
  set certificatechain($core.String value) => $_setString(1, value);
  @$pb.TagNumber(369292378)
  $core.bool hasCertificatechain() => $_has(1);
  @$pb.TagNumber(369292378)
  void clearCertificatechain() => $_clearField(369292378);
}

class HttpRedirect extends $pb.GeneratedMessage {
  factory HttpRedirect({
    $core.String? redirectfrom,
    $core.String? redirectto,
  }) {
    final result = create();
    if (redirectfrom != null) result.redirectfrom = redirectfrom;
    if (redirectto != null) result.redirectto = redirectto;
    return result;
  }

  HttpRedirect._();

  factory HttpRedirect.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HttpRedirect.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HttpRedirect',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(52853514, _omitFieldNames ? '' : 'redirectfrom')
    ..aOS(472243857, _omitFieldNames ? '' : 'redirectto')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HttpRedirect clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HttpRedirect copyWith(void Function(HttpRedirect) updates) =>
      super.copyWith((message) => updates(message as HttpRedirect))
          as HttpRedirect;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HttpRedirect create() => HttpRedirect._();
  @$core.override
  HttpRedirect createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HttpRedirect getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HttpRedirect>(create);
  static HttpRedirect? _defaultInstance;

  @$pb.TagNumber(52853514)
  $core.String get redirectfrom => $_getSZ(0);
  @$pb.TagNumber(52853514)
  set redirectfrom($core.String value) => $_setString(0, value);
  @$pb.TagNumber(52853514)
  $core.bool hasRedirectfrom() => $_has(0);
  @$pb.TagNumber(52853514)
  void clearRedirectfrom() => $_clearField(52853514);

  @$pb.TagNumber(472243857)
  $core.String get redirectto => $_getSZ(1);
  @$pb.TagNumber(472243857)
  set redirectto($core.String value) => $_setString(1, value);
  @$pb.TagNumber(472243857)
  $core.bool hasRedirectto() => $_has(1);
  @$pb.TagNumber(472243857)
  void clearRedirectto() => $_clearField(472243857);
}

class ImportCertificateRequest extends $pb.GeneratedMessage {
  factory ImportCertificateRequest({
    $core.String? certificatearn,
    $core.List<$core.int>? privatekey,
    $core.List<$core.int>? certificate,
    $core.List<$core.int>? certificatechain,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (privatekey != null) result.privatekey = privatekey;
    if (certificate != null) result.certificate = certificate;
    if (certificatechain != null) result.certificatechain = certificatechain;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  ImportCertificateRequest._();

  factory ImportCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..a<$core.List<$core.int>>(
        173312700, _omitFieldNames ? '' : 'privatekey', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(
        198060817, _omitFieldNames ? '' : 'certificate', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(369292378,
        _omitFieldNames ? '' : 'certificatechain', $pb.PbFieldType.OY)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportCertificateRequest copyWith(
          void Function(ImportCertificateRequest) updates) =>
      super.copyWith((message) => updates(message as ImportCertificateRequest))
          as ImportCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportCertificateRequest create() => ImportCertificateRequest._();
  @$core.override
  ImportCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportCertificateRequest>(create);
  static ImportCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(173312700)
  $core.List<$core.int> get privatekey => $_getN(1);
  @$pb.TagNumber(173312700)
  set privatekey($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(173312700)
  $core.bool hasPrivatekey() => $_has(1);
  @$pb.TagNumber(173312700)
  void clearPrivatekey() => $_clearField(173312700);

  @$pb.TagNumber(198060817)
  $core.List<$core.int> get certificate => $_getN(2);
  @$pb.TagNumber(198060817)
  set certificate($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(198060817)
  $core.bool hasCertificate() => $_has(2);
  @$pb.TagNumber(198060817)
  void clearCertificate() => $_clearField(198060817);

  @$pb.TagNumber(369292378)
  $core.List<$core.int> get certificatechain => $_getN(3);
  @$pb.TagNumber(369292378)
  set certificatechain($core.List<$core.int> value) => $_setBytes(3, value);
  @$pb.TagNumber(369292378)
  $core.bool hasCertificatechain() => $_has(3);
  @$pb.TagNumber(369292378)
  void clearCertificatechain() => $_clearField(369292378);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(4);
}

class ImportCertificateResponse extends $pb.GeneratedMessage {
  factory ImportCertificateResponse({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  ImportCertificateResponse._();

  factory ImportCertificateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportCertificateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportCertificateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportCertificateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportCertificateResponse copyWith(
          void Function(ImportCertificateResponse) updates) =>
      super.copyWith((message) => updates(message as ImportCertificateResponse))
          as ImportCertificateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportCertificateResponse create() => ImportCertificateResponse._();
  @$core.override
  ImportCertificateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportCertificateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportCertificateResponse>(create);
  static ImportCertificateResponse? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
}

class InvalidArgsException extends $pb.GeneratedMessage {
  factory InvalidArgsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidArgsException._();

  factory InvalidArgsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidArgsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidArgsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArgsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArgsException copyWith(void Function(InvalidArgsException) updates) =>
      super.copyWith((message) => updates(message as InvalidArgsException))
          as InvalidArgsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidArgsException create() => InvalidArgsException._();
  @$core.override
  InvalidArgsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidArgsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidArgsException>(create);
  static InvalidArgsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidArnException extends $pb.GeneratedMessage {
  factory InvalidArnException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidArnException._();

  factory InvalidArnException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidArnException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidArnException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArnException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArnException copyWith(void Function(InvalidArnException) updates) =>
      super.copyWith((message) => updates(message as InvalidArnException))
          as InvalidArnException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidArnException create() => InvalidArnException._();
  @$core.override
  InvalidArnException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidArnException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidArnException>(create);
  static InvalidArnException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidDomainValidationOptionsException extends $pb.GeneratedMessage {
  factory InvalidDomainValidationOptionsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidDomainValidationOptionsException._();

  factory InvalidDomainValidationOptionsException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidDomainValidationOptionsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidDomainValidationOptionsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDomainValidationOptionsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDomainValidationOptionsException copyWith(
          void Function(InvalidDomainValidationOptionsException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidDomainValidationOptionsException))
          as InvalidDomainValidationOptionsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidDomainValidationOptionsException create() =>
      InvalidDomainValidationOptionsException._();
  @$core.override
  InvalidDomainValidationOptionsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidDomainValidationOptionsException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidDomainValidationOptionsException>(create);
  static InvalidDomainValidationOptionsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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

class InvalidTagException extends $pb.GeneratedMessage {
  factory InvalidTagException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidTagException._();

  factory InvalidTagException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidTagException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidTagException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTagException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTagException copyWith(void Function(InvalidTagException) updates) =>
      super.copyWith((message) => updates(message as InvalidTagException))
          as InvalidTagException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidTagException create() => InvalidTagException._();
  @$core.override
  InvalidTagException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidTagException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidTagException>(create);
  static InvalidTagException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KeyUsage extends $pb.GeneratedMessage {
  factory KeyUsage({
    KeyUsageName? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  KeyUsage._();

  factory KeyUsage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeyUsage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeyUsage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aE<KeyUsageName>(266367751, _omitFieldNames ? '' : 'name',
        enumValues: KeyUsageName.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyUsage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyUsage copyWith(void Function(KeyUsage) updates) =>
      super.copyWith((message) => updates(message as KeyUsage)) as KeyUsage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeyUsage create() => KeyUsage._();
  @$core.override
  KeyUsage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeyUsage getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<KeyUsage>(create);
  static KeyUsage? _defaultInstance;

  @$pb.TagNumber(266367751)
  KeyUsageName get name => $_getN(0);
  @$pb.TagNumber(266367751)
  set name(KeyUsageName value) => $_setField(266367751, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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

class ListCertificatesRequest extends $pb.GeneratedMessage {
  factory ListCertificatesRequest({
    Filters? includes,
    SortBy? sortby,
    $core.Iterable<CertificateStatus>? certificatestatuses,
    $core.String? nexttoken,
    SortOrder? sortorder,
    $core.int? maxitems,
  }) {
    final result = create();
    if (includes != null) result.includes = includes;
    if (sortby != null) result.sortby = sortby;
    if (certificatestatuses != null)
      result.certificatestatuses.addAll(certificatestatuses);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (sortorder != null) result.sortorder = sortorder;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListCertificatesRequest._();

  factory ListCertificatesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCertificatesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCertificatesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOM<Filters>(162404509, _omitFieldNames ? '' : 'includes',
        subBuilder: Filters.create)
    ..aE<SortBy>(186052369, _omitFieldNames ? '' : 'sortby',
        enumValues: SortBy.values)
    ..pc<CertificateStatus>(203375987,
        _omitFieldNames ? '' : 'certificatestatuses', $pb.PbFieldType.KE,
        valueOf: CertificateStatus.valueOf,
        enumValues: CertificateStatus.values,
        defaultEnumValue: CertificateStatus.CERTIFICATE_STATUS_REVOKED)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aE<SortOrder>(274231684, _omitFieldNames ? '' : 'sortorder',
        enumValues: SortOrder.values)
    ..aI(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCertificatesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCertificatesRequest copyWith(
          void Function(ListCertificatesRequest) updates) =>
      super.copyWith((message) => updates(message as ListCertificatesRequest))
          as ListCertificatesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCertificatesRequest create() => ListCertificatesRequest._();
  @$core.override
  ListCertificatesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCertificatesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCertificatesRequest>(create);
  static ListCertificatesRequest? _defaultInstance;

  @$pb.TagNumber(162404509)
  Filters get includes => $_getN(0);
  @$pb.TagNumber(162404509)
  set includes(Filters value) => $_setField(162404509, value);
  @$pb.TagNumber(162404509)
  $core.bool hasIncludes() => $_has(0);
  @$pb.TagNumber(162404509)
  void clearIncludes() => $_clearField(162404509);
  @$pb.TagNumber(162404509)
  Filters ensureIncludes() => $_ensure(0);

  @$pb.TagNumber(186052369)
  SortBy get sortby => $_getN(1);
  @$pb.TagNumber(186052369)
  set sortby(SortBy value) => $_setField(186052369, value);
  @$pb.TagNumber(186052369)
  $core.bool hasSortby() => $_has(1);
  @$pb.TagNumber(186052369)
  void clearSortby() => $_clearField(186052369);

  @$pb.TagNumber(203375987)
  $pb.PbList<CertificateStatus> get certificatestatuses => $_getList(2);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(3);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(3);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(274231684)
  SortOrder get sortorder => $_getN(4);
  @$pb.TagNumber(274231684)
  set sortorder(SortOrder value) => $_setField(274231684, value);
  @$pb.TagNumber(274231684)
  $core.bool hasSortorder() => $_has(4);
  @$pb.TagNumber(274231684)
  void clearSortorder() => $_clearField(274231684);

  @$pb.TagNumber(506899220)
  $core.int get maxitems => $_getIZ(5);
  @$pb.TagNumber(506899220)
  set maxitems($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(5);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListCertificatesResponse extends $pb.GeneratedMessage {
  factory ListCertificatesResponse({
    $core.String? nexttoken,
    $core.Iterable<CertificateSummary>? certificatesummarylist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (certificatesummarylist != null)
      result.certificatesummarylist.addAll(certificatesummarylist);
    return result;
  }

  ListCertificatesResponse._();

  factory ListCertificatesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCertificatesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCertificatesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<CertificateSummary>(
        487936529, _omitFieldNames ? '' : 'certificatesummarylist',
        subBuilder: CertificateSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCertificatesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCertificatesResponse copyWith(
          void Function(ListCertificatesResponse) updates) =>
      super.copyWith((message) => updates(message as ListCertificatesResponse))
          as ListCertificatesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCertificatesResponse create() => ListCertificatesResponse._();
  @$core.override
  ListCertificatesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCertificatesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCertificatesResponse>(create);
  static ListCertificatesResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(487936529)
  $pb.PbList<CertificateSummary> get certificatesummarylist => $_getList(1);
}

class ListTagsForCertificateRequest extends $pb.GeneratedMessage {
  factory ListTagsForCertificateRequest({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  ListTagsForCertificateRequest._();

  factory ListTagsForCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForCertificateRequest copyWith(
          void Function(ListTagsForCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListTagsForCertificateRequest))
          as ListTagsForCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForCertificateRequest create() =>
      ListTagsForCertificateRequest._();
  @$core.override
  ListTagsForCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForCertificateRequest>(create);
  static ListTagsForCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
}

class ListTagsForCertificateResponse extends $pb.GeneratedMessage {
  factory ListTagsForCertificateResponse({
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  ListTagsForCertificateResponse._();

  factory ListTagsForCertificateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForCertificateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForCertificateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForCertificateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForCertificateResponse copyWith(
          void Function(ListTagsForCertificateResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListTagsForCertificateResponse))
          as ListTagsForCertificateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForCertificateResponse create() =>
      ListTagsForCertificateResponse._();
  @$core.override
  ListTagsForCertificateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForCertificateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForCertificateResponse>(create);
  static ListTagsForCertificateResponse? _defaultInstance;

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(0);
}

class PutAccountConfigurationRequest extends $pb.GeneratedMessage {
  factory PutAccountConfigurationRequest({
    $core.String? idempotencytoken,
    ExpiryEventsConfiguration? expiryevents,
  }) {
    final result = create();
    if (idempotencytoken != null) result.idempotencytoken = idempotencytoken;
    if (expiryevents != null) result.expiryevents = expiryevents;
    return result;
  }

  PutAccountConfigurationRequest._();

  factory PutAccountConfigurationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutAccountConfigurationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutAccountConfigurationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(56833648, _omitFieldNames ? '' : 'idempotencytoken')
    ..aOM<ExpiryEventsConfiguration>(
        358763356, _omitFieldNames ? '' : 'expiryevents',
        subBuilder: ExpiryEventsConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAccountConfigurationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAccountConfigurationRequest copyWith(
          void Function(PutAccountConfigurationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PutAccountConfigurationRequest))
          as PutAccountConfigurationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutAccountConfigurationRequest create() =>
      PutAccountConfigurationRequest._();
  @$core.override
  PutAccountConfigurationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutAccountConfigurationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutAccountConfigurationRequest>(create);
  static PutAccountConfigurationRequest? _defaultInstance;

  @$pb.TagNumber(56833648)
  $core.String get idempotencytoken => $_getSZ(0);
  @$pb.TagNumber(56833648)
  set idempotencytoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56833648)
  $core.bool hasIdempotencytoken() => $_has(0);
  @$pb.TagNumber(56833648)
  void clearIdempotencytoken() => $_clearField(56833648);

  @$pb.TagNumber(358763356)
  ExpiryEventsConfiguration get expiryevents => $_getN(1);
  @$pb.TagNumber(358763356)
  set expiryevents(ExpiryEventsConfiguration value) =>
      $_setField(358763356, value);
  @$pb.TagNumber(358763356)
  $core.bool hasExpiryevents() => $_has(1);
  @$pb.TagNumber(358763356)
  void clearExpiryevents() => $_clearField(358763356);
  @$pb.TagNumber(358763356)
  ExpiryEventsConfiguration ensureExpiryevents() => $_ensure(1);
}

class RemoveTagsFromCertificateRequest extends $pb.GeneratedMessage {
  factory RemoveTagsFromCertificateRequest({
    $core.String? certificatearn,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  RemoveTagsFromCertificateRequest._();

  factory RemoveTagsFromCertificateRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveTagsFromCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveTagsFromCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsFromCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTagsFromCertificateRequest copyWith(
          void Function(RemoveTagsFromCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as RemoveTagsFromCertificateRequest))
          as RemoveTagsFromCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveTagsFromCertificateRequest create() =>
      RemoveTagsFromCertificateRequest._();
  @$core.override
  RemoveTagsFromCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveTagsFromCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveTagsFromCertificateRequest>(
          create);
  static RemoveTagsFromCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);
}

class RenewCertificateRequest extends $pb.GeneratedMessage {
  factory RenewCertificateRequest({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  RenewCertificateRequest._();

  factory RenewCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RenewCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RenewCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewCertificateRequest copyWith(
          void Function(RenewCertificateRequest) updates) =>
      super.copyWith((message) => updates(message as RenewCertificateRequest))
          as RenewCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RenewCertificateRequest create() => RenewCertificateRequest._();
  @$core.override
  RenewCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RenewCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RenewCertificateRequest>(create);
  static RenewCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
}

class RenewalSummary extends $pb.GeneratedMessage {
  factory RenewalSummary({
    FailureReason? renewalstatusreason,
    $core.String? updatedat,
    RenewalStatus? renewalstatus,
    $core.Iterable<DomainValidation>? domainvalidationoptions,
  }) {
    final result = create();
    if (renewalstatusreason != null)
      result.renewalstatusreason = renewalstatusreason;
    if (updatedat != null) result.updatedat = updatedat;
    if (renewalstatus != null) result.renewalstatus = renewalstatus;
    if (domainvalidationoptions != null)
      result.domainvalidationoptions.addAll(domainvalidationoptions);
    return result;
  }

  RenewalSummary._();

  factory RenewalSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RenewalSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RenewalSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aE<FailureReason>(206499082, _omitFieldNames ? '' : 'renewalstatusreason',
        enumValues: FailureReason.values)
    ..aOS(213581206, _omitFieldNames ? '' : 'updatedat')
    ..aE<RenewalStatus>(277232086, _omitFieldNames ? '' : 'renewalstatus',
        enumValues: RenewalStatus.values)
    ..pPM<DomainValidation>(
        335573705, _omitFieldNames ? '' : 'domainvalidationoptions',
        subBuilder: DomainValidation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewalSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RenewalSummary copyWith(void Function(RenewalSummary) updates) =>
      super.copyWith((message) => updates(message as RenewalSummary))
          as RenewalSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RenewalSummary create() => RenewalSummary._();
  @$core.override
  RenewalSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RenewalSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RenewalSummary>(create);
  static RenewalSummary? _defaultInstance;

  @$pb.TagNumber(206499082)
  FailureReason get renewalstatusreason => $_getN(0);
  @$pb.TagNumber(206499082)
  set renewalstatusreason(FailureReason value) => $_setField(206499082, value);
  @$pb.TagNumber(206499082)
  $core.bool hasRenewalstatusreason() => $_has(0);
  @$pb.TagNumber(206499082)
  void clearRenewalstatusreason() => $_clearField(206499082);

  @$pb.TagNumber(213581206)
  $core.String get updatedat => $_getSZ(1);
  @$pb.TagNumber(213581206)
  set updatedat($core.String value) => $_setString(1, value);
  @$pb.TagNumber(213581206)
  $core.bool hasUpdatedat() => $_has(1);
  @$pb.TagNumber(213581206)
  void clearUpdatedat() => $_clearField(213581206);

  @$pb.TagNumber(277232086)
  RenewalStatus get renewalstatus => $_getN(2);
  @$pb.TagNumber(277232086)
  set renewalstatus(RenewalStatus value) => $_setField(277232086, value);
  @$pb.TagNumber(277232086)
  $core.bool hasRenewalstatus() => $_has(2);
  @$pb.TagNumber(277232086)
  void clearRenewalstatus() => $_clearField(277232086);

  @$pb.TagNumber(335573705)
  $pb.PbList<DomainValidation> get domainvalidationoptions => $_getList(3);
}

class RequestCertificateRequest extends $pb.GeneratedMessage {
  factory RequestCertificateRequest({
    $core.String? idempotencytoken,
    ValidationMethod? validationmethod,
    $core.Iterable<$core.String>? subjectalternativenames,
    $core.String? domainname,
    $core.String? certificateauthorityarn,
    $core.Iterable<DomainValidationOption>? domainvalidationoptions,
    CertificateOptions? options,
    $core.Iterable<Tag>? tags,
    KeyAlgorithm? keyalgorithm,
    CertificateManagedBy? managedby,
  }) {
    final result = create();
    if (idempotencytoken != null) result.idempotencytoken = idempotencytoken;
    if (validationmethod != null) result.validationmethod = validationmethod;
    if (subjectalternativenames != null)
      result.subjectalternativenames.addAll(subjectalternativenames);
    if (domainname != null) result.domainname = domainname;
    if (certificateauthorityarn != null)
      result.certificateauthorityarn = certificateauthorityarn;
    if (domainvalidationoptions != null)
      result.domainvalidationoptions.addAll(domainvalidationoptions);
    if (options != null) result.options = options;
    if (tags != null) result.tags.addAll(tags);
    if (keyalgorithm != null) result.keyalgorithm = keyalgorithm;
    if (managedby != null) result.managedby = managedby;
    return result;
  }

  RequestCertificateRequest._();

  factory RequestCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(56833648, _omitFieldNames ? '' : 'idempotencytoken')
    ..aE<ValidationMethod>(58092520, _omitFieldNames ? '' : 'validationmethod',
        enumValues: ValidationMethod.values)
    ..pPS(109998119, _omitFieldNames ? '' : 'subjectalternativenames')
    ..aOS(194914027, _omitFieldNames ? '' : 'domainname')
    ..aOS(266069181, _omitFieldNames ? '' : 'certificateauthorityarn')
    ..pPM<DomainValidationOption>(
        335573705, _omitFieldNames ? '' : 'domainvalidationoptions',
        subBuilder: DomainValidationOption.create)
    ..aOM<CertificateOptions>(356388166, _omitFieldNames ? '' : 'options',
        subBuilder: CertificateOptions.create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aE<KeyAlgorithm>(452317818, _omitFieldNames ? '' : 'keyalgorithm',
        enumValues: KeyAlgorithm.values)
    ..aE<CertificateManagedBy>(455511232, _omitFieldNames ? '' : 'managedby',
        enumValues: CertificateManagedBy.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestCertificateRequest copyWith(
          void Function(RequestCertificateRequest) updates) =>
      super.copyWith((message) => updates(message as RequestCertificateRequest))
          as RequestCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestCertificateRequest create() => RequestCertificateRequest._();
  @$core.override
  RequestCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestCertificateRequest>(create);
  static RequestCertificateRequest? _defaultInstance;

  @$pb.TagNumber(56833648)
  $core.String get idempotencytoken => $_getSZ(0);
  @$pb.TagNumber(56833648)
  set idempotencytoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56833648)
  $core.bool hasIdempotencytoken() => $_has(0);
  @$pb.TagNumber(56833648)
  void clearIdempotencytoken() => $_clearField(56833648);

  @$pb.TagNumber(58092520)
  ValidationMethod get validationmethod => $_getN(1);
  @$pb.TagNumber(58092520)
  set validationmethod(ValidationMethod value) => $_setField(58092520, value);
  @$pb.TagNumber(58092520)
  $core.bool hasValidationmethod() => $_has(1);
  @$pb.TagNumber(58092520)
  void clearValidationmethod() => $_clearField(58092520);

  @$pb.TagNumber(109998119)
  $pb.PbList<$core.String> get subjectalternativenames => $_getList(2);

  @$pb.TagNumber(194914027)
  $core.String get domainname => $_getSZ(3);
  @$pb.TagNumber(194914027)
  set domainname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(194914027)
  $core.bool hasDomainname() => $_has(3);
  @$pb.TagNumber(194914027)
  void clearDomainname() => $_clearField(194914027);

  @$pb.TagNumber(266069181)
  $core.String get certificateauthorityarn => $_getSZ(4);
  @$pb.TagNumber(266069181)
  set certificateauthorityarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266069181)
  $core.bool hasCertificateauthorityarn() => $_has(4);
  @$pb.TagNumber(266069181)
  void clearCertificateauthorityarn() => $_clearField(266069181);

  @$pb.TagNumber(335573705)
  $pb.PbList<DomainValidationOption> get domainvalidationoptions =>
      $_getList(5);

  @$pb.TagNumber(356388166)
  CertificateOptions get options => $_getN(6);
  @$pb.TagNumber(356388166)
  set options(CertificateOptions value) => $_setField(356388166, value);
  @$pb.TagNumber(356388166)
  $core.bool hasOptions() => $_has(6);
  @$pb.TagNumber(356388166)
  void clearOptions() => $_clearField(356388166);
  @$pb.TagNumber(356388166)
  CertificateOptions ensureOptions() => $_ensure(6);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(7);

  @$pb.TagNumber(452317818)
  KeyAlgorithm get keyalgorithm => $_getN(8);
  @$pb.TagNumber(452317818)
  set keyalgorithm(KeyAlgorithm value) => $_setField(452317818, value);
  @$pb.TagNumber(452317818)
  $core.bool hasKeyalgorithm() => $_has(8);
  @$pb.TagNumber(452317818)
  void clearKeyalgorithm() => $_clearField(452317818);

  @$pb.TagNumber(455511232)
  CertificateManagedBy get managedby => $_getN(9);
  @$pb.TagNumber(455511232)
  set managedby(CertificateManagedBy value) => $_setField(455511232, value);
  @$pb.TagNumber(455511232)
  $core.bool hasManagedby() => $_has(9);
  @$pb.TagNumber(455511232)
  void clearManagedby() => $_clearField(455511232);
}

class RequestCertificateResponse extends $pb.GeneratedMessage {
  factory RequestCertificateResponse({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  RequestCertificateResponse._();

  factory RequestCertificateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestCertificateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestCertificateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestCertificateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestCertificateResponse copyWith(
          void Function(RequestCertificateResponse) updates) =>
      super.copyWith(
              (message) => updates(message as RequestCertificateResponse))
          as RequestCertificateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestCertificateResponse create() => RequestCertificateResponse._();
  @$core.override
  RequestCertificateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestCertificateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestCertificateResponse>(create);
  static RequestCertificateResponse? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
}

class RequestInProgressException extends $pb.GeneratedMessage {
  factory RequestInProgressException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  RequestInProgressException._();

  factory RequestInProgressException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestInProgressException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestInProgressException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestInProgressException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestInProgressException copyWith(
          void Function(RequestInProgressException) updates) =>
      super.copyWith(
              (message) => updates(message as RequestInProgressException))
          as RequestInProgressException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestInProgressException create() => RequestInProgressException._();
  @$core.override
  RequestInProgressException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestInProgressException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestInProgressException>(create);
  static RequestInProgressException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ResendValidationEmailRequest extends $pb.GeneratedMessage {
  factory ResendValidationEmailRequest({
    $core.String? certificatearn,
    $core.String? validationdomain,
    $core.String? domain,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (validationdomain != null) result.validationdomain = validationdomain;
    if (domain != null) result.domain = domain;
    return result;
  }

  ResendValidationEmailRequest._();

  factory ResendValidationEmailRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResendValidationEmailRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResendValidationEmailRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..aOS(161300799, _omitFieldNames ? '' : 'validationdomain')
    ..aOS(505186578, _omitFieldNames ? '' : 'domain')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResendValidationEmailRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResendValidationEmailRequest copyWith(
          void Function(ResendValidationEmailRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ResendValidationEmailRequest))
          as ResendValidationEmailRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResendValidationEmailRequest create() =>
      ResendValidationEmailRequest._();
  @$core.override
  ResendValidationEmailRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResendValidationEmailRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResendValidationEmailRequest>(create);
  static ResendValidationEmailRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(161300799)
  $core.String get validationdomain => $_getSZ(1);
  @$pb.TagNumber(161300799)
  set validationdomain($core.String value) => $_setString(1, value);
  @$pb.TagNumber(161300799)
  $core.bool hasValidationdomain() => $_has(1);
  @$pb.TagNumber(161300799)
  void clearValidationdomain() => $_clearField(161300799);

  @$pb.TagNumber(505186578)
  $core.String get domain => $_getSZ(2);
  @$pb.TagNumber(505186578)
  set domain($core.String value) => $_setString(2, value);
  @$pb.TagNumber(505186578)
  $core.bool hasDomain() => $_has(2);
  @$pb.TagNumber(505186578)
  void clearDomain() => $_clearField(505186578);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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

class ResourceRecord extends $pb.GeneratedMessage {
  factory ResourceRecord({
    $core.String? name,
    $core.String? value,
    RecordType? type,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    if (type != null) result.type = type;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..aE<RecordType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: RecordType.values)
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

  @$pb.TagNumber(290836590)
  RecordType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(RecordType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class RevokeCertificateRequest extends $pb.GeneratedMessage {
  factory RevokeCertificateRequest({
    $core.String? certificatearn,
    RevocationReason? revocationreason,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (revocationreason != null) result.revocationreason = revocationreason;
    return result;
  }

  RevokeCertificateRequest._();

  factory RevokeCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RevokeCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RevokeCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..aE<RevocationReason>(331836358, _omitFieldNames ? '' : 'revocationreason',
        enumValues: RevocationReason.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RevokeCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RevokeCertificateRequest copyWith(
          void Function(RevokeCertificateRequest) updates) =>
      super.copyWith((message) => updates(message as RevokeCertificateRequest))
          as RevokeCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RevokeCertificateRequest create() => RevokeCertificateRequest._();
  @$core.override
  RevokeCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RevokeCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RevokeCertificateRequest>(create);
  static RevokeCertificateRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(331836358)
  RevocationReason get revocationreason => $_getN(1);
  @$pb.TagNumber(331836358)
  set revocationreason(RevocationReason value) => $_setField(331836358, value);
  @$pb.TagNumber(331836358)
  $core.bool hasRevocationreason() => $_has(1);
  @$pb.TagNumber(331836358)
  void clearRevocationreason() => $_clearField(331836358);
}

class RevokeCertificateResponse extends $pb.GeneratedMessage {
  factory RevokeCertificateResponse({
    $core.String? certificatearn,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    return result;
  }

  RevokeCertificateResponse._();

  factory RevokeCertificateResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RevokeCertificateResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RevokeCertificateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RevokeCertificateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RevokeCertificateResponse copyWith(
          void Function(RevokeCertificateResponse) updates) =>
      super.copyWith((message) => updates(message as RevokeCertificateResponse))
          as RevokeCertificateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RevokeCertificateResponse create() => RevokeCertificateResponse._();
  @$core.override
  RevokeCertificateResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RevokeCertificateResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RevokeCertificateResponse>(create);
  static RevokeCertificateResponse? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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

class TagPolicyException extends $pb.GeneratedMessage {
  factory TagPolicyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TagPolicyException._();

  factory TagPolicyException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagPolicyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagPolicyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagPolicyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagPolicyException copyWith(void Function(TagPolicyException) updates) =>
      super.copyWith((message) => updates(message as TagPolicyException))
          as TagPolicyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagPolicyException create() => TagPolicyException._();
  @$core.override
  TagPolicyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagPolicyException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagPolicyException>(create);
  static TagPolicyException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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

class TooManyTagsException extends $pb.GeneratedMessage {
  factory TooManyTagsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyTagsException._();

  factory TooManyTagsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyTagsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyTagsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTagsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTagsException copyWith(void Function(TooManyTagsException) updates) =>
      super.copyWith((message) => updates(message as TooManyTagsException))
          as TooManyTagsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyTagsException create() => TooManyTagsException._();
  @$core.override
  TooManyTagsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyTagsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyTagsException>(create);
  static TooManyTagsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class UpdateCertificateOptionsRequest extends $pb.GeneratedMessage {
  factory UpdateCertificateOptionsRequest({
    $core.String? certificatearn,
    CertificateOptions? options,
  }) {
    final result = create();
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (options != null) result.options = options;
    return result;
  }

  UpdateCertificateOptionsRequest._();

  factory UpdateCertificateOptionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateCertificateOptionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateCertificateOptionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
      createEmptyInstance: create)
    ..aOS(92693880, _omitFieldNames ? '' : 'certificatearn')
    ..aOM<CertificateOptions>(356388166, _omitFieldNames ? '' : 'options',
        subBuilder: CertificateOptions.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCertificateOptionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCertificateOptionsRequest copyWith(
          void Function(UpdateCertificateOptionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateCertificateOptionsRequest))
          as UpdateCertificateOptionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateCertificateOptionsRequest create() =>
      UpdateCertificateOptionsRequest._();
  @$core.override
  UpdateCertificateOptionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateCertificateOptionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateCertificateOptionsRequest>(
          create);
  static UpdateCertificateOptionsRequest? _defaultInstance;

  @$pb.TagNumber(92693880)
  $core.String get certificatearn => $_getSZ(0);
  @$pb.TagNumber(92693880)
  set certificatearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(92693880)
  $core.bool hasCertificatearn() => $_has(0);
  @$pb.TagNumber(92693880)
  void clearCertificatearn() => $_clearField(92693880);

  @$pb.TagNumber(356388166)
  CertificateOptions get options => $_getN(1);
  @$pb.TagNumber(356388166)
  set options(CertificateOptions value) => $_setField(356388166, value);
  @$pb.TagNumber(356388166)
  $core.bool hasOptions() => $_has(1);
  @$pb.TagNumber(356388166)
  void clearOptions() => $_clearField(356388166);
  @$pb.TagNumber(356388166)
  CertificateOptions ensureOptions() => $_ensure(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'acm'),
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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
