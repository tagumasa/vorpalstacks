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

class CertificateExport extends $pb.ProtobufEnum {
  static const CertificateExport CERTIFICATE_EXPORT_DISABLED =
      CertificateExport._(
          0, _omitEnumNames ? '' : 'CERTIFICATE_EXPORT_DISABLED');
  static const CertificateExport CERTIFICATE_EXPORT_ENABLED =
      CertificateExport._(
          1, _omitEnumNames ? '' : 'CERTIFICATE_EXPORT_ENABLED');

  static const $core.List<CertificateExport> values = <CertificateExport>[
    CERTIFICATE_EXPORT_DISABLED,
    CERTIFICATE_EXPORT_ENABLED,
  ];

  static final $core.List<CertificateExport?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CertificateExport? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CertificateExport._(super.value, super.name);
}

class CertificateManagedBy extends $pb.ProtobufEnum {
  static const CertificateManagedBy CERTIFICATE_MANAGED_BY_CLOUDFRONT =
      CertificateManagedBy._(
          0, _omitEnumNames ? '' : 'CERTIFICATE_MANAGED_BY_CLOUDFRONT');

  static const $core.List<CertificateManagedBy> values = <CertificateManagedBy>[
    CERTIFICATE_MANAGED_BY_CLOUDFRONT,
  ];

  static final $core.List<CertificateManagedBy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static CertificateManagedBy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CertificateManagedBy._(super.value, super.name);
}

class CertificateStatus extends $pb.ProtobufEnum {
  static const CertificateStatus CERTIFICATE_STATUS_REVOKED =
      CertificateStatus._(
          0, _omitEnumNames ? '' : 'CERTIFICATE_STATUS_REVOKED');
  static const CertificateStatus CERTIFICATE_STATUS_ISSUED =
      CertificateStatus._(1, _omitEnumNames ? '' : 'CERTIFICATE_STATUS_ISSUED');
  static const CertificateStatus CERTIFICATE_STATUS_PENDING_VALIDATION =
      CertificateStatus._(
          2, _omitEnumNames ? '' : 'CERTIFICATE_STATUS_PENDING_VALIDATION');
  static const CertificateStatus CERTIFICATE_STATUS_VALIDATION_TIMED_OUT =
      CertificateStatus._(
          3, _omitEnumNames ? '' : 'CERTIFICATE_STATUS_VALIDATION_TIMED_OUT');
  static const CertificateStatus CERTIFICATE_STATUS_EXPIRED =
      CertificateStatus._(
          4, _omitEnumNames ? '' : 'CERTIFICATE_STATUS_EXPIRED');
  static const CertificateStatus CERTIFICATE_STATUS_INACTIVE =
      CertificateStatus._(
          5, _omitEnumNames ? '' : 'CERTIFICATE_STATUS_INACTIVE');
  static const CertificateStatus CERTIFICATE_STATUS_FAILED =
      CertificateStatus._(6, _omitEnumNames ? '' : 'CERTIFICATE_STATUS_FAILED');

  static const $core.List<CertificateStatus> values = <CertificateStatus>[
    CERTIFICATE_STATUS_REVOKED,
    CERTIFICATE_STATUS_ISSUED,
    CERTIFICATE_STATUS_PENDING_VALIDATION,
    CERTIFICATE_STATUS_VALIDATION_TIMED_OUT,
    CERTIFICATE_STATUS_EXPIRED,
    CERTIFICATE_STATUS_INACTIVE,
    CERTIFICATE_STATUS_FAILED,
  ];

  static final $core.List<CertificateStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static CertificateStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CertificateStatus._(super.value, super.name);
}

class CertificateTransparencyLoggingPreference extends $pb.ProtobufEnum {
  static const CertificateTransparencyLoggingPreference
      CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_DISABLED =
      CertificateTransparencyLoggingPreference._(
          0,
          _omitEnumNames
              ? ''
              : 'CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_DISABLED');
  static const CertificateTransparencyLoggingPreference
      CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_ENABLED =
      CertificateTransparencyLoggingPreference._(
          1,
          _omitEnumNames
              ? ''
              : 'CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_ENABLED');

  static const $core.List<CertificateTransparencyLoggingPreference> values =
      <CertificateTransparencyLoggingPreference>[
    CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_DISABLED,
    CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_ENABLED,
  ];

  static final $core.List<CertificateTransparencyLoggingPreference?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CertificateTransparencyLoggingPreference? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CertificateTransparencyLoggingPreference._(super.value, super.name);
}

class CertificateType extends $pb.ProtobufEnum {
  static const CertificateType CERTIFICATE_TYPE_IMPORTED =
      CertificateType._(0, _omitEnumNames ? '' : 'CERTIFICATE_TYPE_IMPORTED');
  static const CertificateType CERTIFICATE_TYPE_PRIVATE =
      CertificateType._(1, _omitEnumNames ? '' : 'CERTIFICATE_TYPE_PRIVATE');
  static const CertificateType CERTIFICATE_TYPE_AMAZON_ISSUED =
      CertificateType._(
          2, _omitEnumNames ? '' : 'CERTIFICATE_TYPE_AMAZON_ISSUED');

  static const $core.List<CertificateType> values = <CertificateType>[
    CERTIFICATE_TYPE_IMPORTED,
    CERTIFICATE_TYPE_PRIVATE,
    CERTIFICATE_TYPE_AMAZON_ISSUED,
  ];

  static final $core.List<CertificateType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static CertificateType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CertificateType._(super.value, super.name);
}

class DomainStatus extends $pb.ProtobufEnum {
  static const DomainStatus DOMAIN_STATUS_PENDING_VALIDATION = DomainStatus._(
      0, _omitEnumNames ? '' : 'DOMAIN_STATUS_PENDING_VALIDATION');
  static const DomainStatus DOMAIN_STATUS_SUCCESS =
      DomainStatus._(1, _omitEnumNames ? '' : 'DOMAIN_STATUS_SUCCESS');
  static const DomainStatus DOMAIN_STATUS_FAILED =
      DomainStatus._(2, _omitEnumNames ? '' : 'DOMAIN_STATUS_FAILED');

  static const $core.List<DomainStatus> values = <DomainStatus>[
    DOMAIN_STATUS_PENDING_VALIDATION,
    DOMAIN_STATUS_SUCCESS,
    DOMAIN_STATUS_FAILED,
  ];

  static final $core.List<DomainStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static DomainStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DomainStatus._(super.value, super.name);
}

class ExtendedKeyUsageName extends $pb.ProtobufEnum {
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_ANY =
      ExtendedKeyUsageName._(
          0, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_ANY');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_IPSEC_END_SYSTEM =
      ExtendedKeyUsageName._(
          1, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_IPSEC_END_SYSTEM');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_IPSEC_USER =
      ExtendedKeyUsageName._(
          2, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_IPSEC_USER');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_CUSTOM =
      ExtendedKeyUsageName._(
          3, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_CUSTOM');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_NONE =
      ExtendedKeyUsageName._(
          4, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_NONE');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_EMAIL_PROTECTION =
      ExtendedKeyUsageName._(
          5, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_EMAIL_PROTECTION');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_TIME_STAMPING =
      ExtendedKeyUsageName._(
          6, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_TIME_STAMPING');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_CODE_SIGNING =
      ExtendedKeyUsageName._(
          7, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_CODE_SIGNING');
  static const ExtendedKeyUsageName
      EXTENDED_KEY_USAGE_NAME_TLS_WEB_CLIENT_AUTHENTICATION =
      ExtendedKeyUsageName._(
          8,
          _omitEnumNames
              ? ''
              : 'EXTENDED_KEY_USAGE_NAME_TLS_WEB_CLIENT_AUTHENTICATION');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_IPSEC_TUNNEL =
      ExtendedKeyUsageName._(
          9, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_IPSEC_TUNNEL');
  static const ExtendedKeyUsageName EXTENDED_KEY_USAGE_NAME_OCSP_SIGNING =
      ExtendedKeyUsageName._(
          10, _omitEnumNames ? '' : 'EXTENDED_KEY_USAGE_NAME_OCSP_SIGNING');
  static const ExtendedKeyUsageName
      EXTENDED_KEY_USAGE_NAME_TLS_WEB_SERVER_AUTHENTICATION =
      ExtendedKeyUsageName._(
          11,
          _omitEnumNames
              ? ''
              : 'EXTENDED_KEY_USAGE_NAME_TLS_WEB_SERVER_AUTHENTICATION');

  static const $core.List<ExtendedKeyUsageName> values = <ExtendedKeyUsageName>[
    EXTENDED_KEY_USAGE_NAME_ANY,
    EXTENDED_KEY_USAGE_NAME_IPSEC_END_SYSTEM,
    EXTENDED_KEY_USAGE_NAME_IPSEC_USER,
    EXTENDED_KEY_USAGE_NAME_CUSTOM,
    EXTENDED_KEY_USAGE_NAME_NONE,
    EXTENDED_KEY_USAGE_NAME_EMAIL_PROTECTION,
    EXTENDED_KEY_USAGE_NAME_TIME_STAMPING,
    EXTENDED_KEY_USAGE_NAME_CODE_SIGNING,
    EXTENDED_KEY_USAGE_NAME_TLS_WEB_CLIENT_AUTHENTICATION,
    EXTENDED_KEY_USAGE_NAME_IPSEC_TUNNEL,
    EXTENDED_KEY_USAGE_NAME_OCSP_SIGNING,
    EXTENDED_KEY_USAGE_NAME_TLS_WEB_SERVER_AUTHENTICATION,
  ];

  static final $core.List<ExtendedKeyUsageName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 11);
  static ExtendedKeyUsageName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExtendedKeyUsageName._(super.value, super.name);
}

class FailureReason extends $pb.ProtobufEnum {
  static const FailureReason FAILURE_REASON_PCA_INVALID_ARGS = FailureReason._(
      0, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_INVALID_ARGS');
  static const FailureReason FAILURE_REASON_PCA_INVALID_STATE = FailureReason._(
      1, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_INVALID_STATE');
  static const FailureReason FAILURE_REASON_OTHER =
      FailureReason._(2, _omitEnumNames ? '' : 'FAILURE_REASON_OTHER');
  static const FailureReason FAILURE_REASON_CAA_ERROR =
      FailureReason._(3, _omitEnumNames ? '' : 'FAILURE_REASON_CAA_ERROR');
  static const FailureReason FAILURE_REASON_ADDITIONAL_VERIFICATION_REQUIRED =
      FailureReason._(
          4,
          _omitEnumNames
              ? ''
              : 'FAILURE_REASON_ADDITIONAL_VERIFICATION_REQUIRED');
  static const FailureReason FAILURE_REASON_PCA_NAME_CONSTRAINTS_VALIDATION =
      FailureReason._(
          5,
          _omitEnumNames
              ? ''
              : 'FAILURE_REASON_PCA_NAME_CONSTRAINTS_VALIDATION');
  static const FailureReason FAILURE_REASON_INVALID_PUBLIC_DOMAIN =
      FailureReason._(
          6, _omitEnumNames ? '' : 'FAILURE_REASON_INVALID_PUBLIC_DOMAIN');
  static const FailureReason FAILURE_REASON_DOMAIN_NOT_ALLOWED =
      FailureReason._(
          7, _omitEnumNames ? '' : 'FAILURE_REASON_DOMAIN_NOT_ALLOWED');
  static const FailureReason FAILURE_REASON_NO_AVAILABLE_CONTACTS =
      FailureReason._(
          8, _omitEnumNames ? '' : 'FAILURE_REASON_NO_AVAILABLE_CONTACTS');
  static const FailureReason FAILURE_REASON_PCA_INVALID_DURATION =
      FailureReason._(
          9, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_INVALID_DURATION');
  static const FailureReason FAILURE_REASON_PCA_REQUEST_FAILED =
      FailureReason._(
          10, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_REQUEST_FAILED');
  static const FailureReason FAILURE_REASON_DOMAIN_VALIDATION_DENIED =
      FailureReason._(
          11, _omitEnumNames ? '' : 'FAILURE_REASON_DOMAIN_VALIDATION_DENIED');
  static const FailureReason FAILURE_REASON_PCA_INVALID_ARN = FailureReason._(
      12, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_INVALID_ARN');
  static const FailureReason FAILURE_REASON_PCA_RESOURCE_NOT_FOUND =
      FailureReason._(
          13, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_RESOURCE_NOT_FOUND');
  static const FailureReason FAILURE_REASON_PCA_ACCESS_DENIED = FailureReason._(
      14, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_ACCESS_DENIED');
  static const FailureReason FAILURE_REASON_PCA_LIMIT_EXCEEDED =
      FailureReason._(
          15, _omitEnumNames ? '' : 'FAILURE_REASON_PCA_LIMIT_EXCEEDED');
  static const FailureReason FAILURE_REASON_SLR_NOT_FOUND =
      FailureReason._(16, _omitEnumNames ? '' : 'FAILURE_REASON_SLR_NOT_FOUND');

  static const $core.List<FailureReason> values = <FailureReason>[
    FAILURE_REASON_PCA_INVALID_ARGS,
    FAILURE_REASON_PCA_INVALID_STATE,
    FAILURE_REASON_OTHER,
    FAILURE_REASON_CAA_ERROR,
    FAILURE_REASON_ADDITIONAL_VERIFICATION_REQUIRED,
    FAILURE_REASON_PCA_NAME_CONSTRAINTS_VALIDATION,
    FAILURE_REASON_INVALID_PUBLIC_DOMAIN,
    FAILURE_REASON_DOMAIN_NOT_ALLOWED,
    FAILURE_REASON_NO_AVAILABLE_CONTACTS,
    FAILURE_REASON_PCA_INVALID_DURATION,
    FAILURE_REASON_PCA_REQUEST_FAILED,
    FAILURE_REASON_DOMAIN_VALIDATION_DENIED,
    FAILURE_REASON_PCA_INVALID_ARN,
    FAILURE_REASON_PCA_RESOURCE_NOT_FOUND,
    FAILURE_REASON_PCA_ACCESS_DENIED,
    FAILURE_REASON_PCA_LIMIT_EXCEEDED,
    FAILURE_REASON_SLR_NOT_FOUND,
  ];

  static final $core.List<FailureReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 16);
  static FailureReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FailureReason._(super.value, super.name);
}

class KeyAlgorithm extends $pb.ProtobufEnum {
  static const KeyAlgorithm KEY_ALGORITHM_EC_PRIME256V1 =
      KeyAlgorithm._(0, _omitEnumNames ? '' : 'KEY_ALGORITHM_EC_PRIME256V1');
  static const KeyAlgorithm KEY_ALGORITHM_RSA_2048 =
      KeyAlgorithm._(1, _omitEnumNames ? '' : 'KEY_ALGORITHM_RSA_2048');
  static const KeyAlgorithm KEY_ALGORITHM_RSA_3072 =
      KeyAlgorithm._(2, _omitEnumNames ? '' : 'KEY_ALGORITHM_RSA_3072');
  static const KeyAlgorithm KEY_ALGORITHM_EC_SECP521R1 =
      KeyAlgorithm._(3, _omitEnumNames ? '' : 'KEY_ALGORITHM_EC_SECP521R1');
  static const KeyAlgorithm KEY_ALGORITHM_RSA_1024 =
      KeyAlgorithm._(4, _omitEnumNames ? '' : 'KEY_ALGORITHM_RSA_1024');
  static const KeyAlgorithm KEY_ALGORITHM_RSA_4096 =
      KeyAlgorithm._(5, _omitEnumNames ? '' : 'KEY_ALGORITHM_RSA_4096');
  static const KeyAlgorithm KEY_ALGORITHM_EC_SECP384R1 =
      KeyAlgorithm._(6, _omitEnumNames ? '' : 'KEY_ALGORITHM_EC_SECP384R1');

  static const $core.List<KeyAlgorithm> values = <KeyAlgorithm>[
    KEY_ALGORITHM_EC_PRIME256V1,
    KEY_ALGORITHM_RSA_2048,
    KEY_ALGORITHM_RSA_3072,
    KEY_ALGORITHM_EC_SECP521R1,
    KEY_ALGORITHM_RSA_1024,
    KEY_ALGORITHM_RSA_4096,
    KEY_ALGORITHM_EC_SECP384R1,
  ];

  static final $core.List<KeyAlgorithm?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static KeyAlgorithm? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyAlgorithm._(super.value, super.name);
}

class KeyUsageName extends $pb.ProtobufEnum {
  static const KeyUsageName KEY_USAGE_NAME_ANY =
      KeyUsageName._(0, _omitEnumNames ? '' : 'KEY_USAGE_NAME_ANY');
  static const KeyUsageName KEY_USAGE_NAME_CERTIFICATE_SIGNING = KeyUsageName._(
      1, _omitEnumNames ? '' : 'KEY_USAGE_NAME_CERTIFICATE_SIGNING');
  static const KeyUsageName KEY_USAGE_NAME_CUSTOM =
      KeyUsageName._(2, _omitEnumNames ? '' : 'KEY_USAGE_NAME_CUSTOM');
  static const KeyUsageName KEY_USAGE_NAME_DATA_ENCIPHERMENT = KeyUsageName._(
      3, _omitEnumNames ? '' : 'KEY_USAGE_NAME_DATA_ENCIPHERMENT');
  static const KeyUsageName KEY_USAGE_NAME_CRL_SIGNING =
      KeyUsageName._(4, _omitEnumNames ? '' : 'KEY_USAGE_NAME_CRL_SIGNING');
  static const KeyUsageName KEY_USAGE_NAME_KEY_AGREEMENT =
      KeyUsageName._(5, _omitEnumNames ? '' : 'KEY_USAGE_NAME_KEY_AGREEMENT');
  static const KeyUsageName KEY_USAGE_NAME_DECIPHER_ONLY =
      KeyUsageName._(6, _omitEnumNames ? '' : 'KEY_USAGE_NAME_DECIPHER_ONLY');
  static const KeyUsageName KEY_USAGE_NAME_ENCHIPER_ONLY =
      KeyUsageName._(7, _omitEnumNames ? '' : 'KEY_USAGE_NAME_ENCHIPER_ONLY');
  static const KeyUsageName KEY_USAGE_NAME_NON_REPUDATION =
      KeyUsageName._(8, _omitEnumNames ? '' : 'KEY_USAGE_NAME_NON_REPUDATION');
  static const KeyUsageName KEY_USAGE_NAME_KEY_ENCIPHERMENT = KeyUsageName._(
      9, _omitEnumNames ? '' : 'KEY_USAGE_NAME_KEY_ENCIPHERMENT');
  static const KeyUsageName KEY_USAGE_NAME_DIGITAL_SIGNATURE = KeyUsageName._(
      10, _omitEnumNames ? '' : 'KEY_USAGE_NAME_DIGITAL_SIGNATURE');

  static const $core.List<KeyUsageName> values = <KeyUsageName>[
    KEY_USAGE_NAME_ANY,
    KEY_USAGE_NAME_CERTIFICATE_SIGNING,
    KEY_USAGE_NAME_CUSTOM,
    KEY_USAGE_NAME_DATA_ENCIPHERMENT,
    KEY_USAGE_NAME_CRL_SIGNING,
    KEY_USAGE_NAME_KEY_AGREEMENT,
    KEY_USAGE_NAME_DECIPHER_ONLY,
    KEY_USAGE_NAME_ENCHIPER_ONLY,
    KEY_USAGE_NAME_NON_REPUDATION,
    KEY_USAGE_NAME_KEY_ENCIPHERMENT,
    KEY_USAGE_NAME_DIGITAL_SIGNATURE,
  ];

  static final $core.List<KeyUsageName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 10);
  static KeyUsageName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyUsageName._(super.value, super.name);
}

class RecordType extends $pb.ProtobufEnum {
  static const RecordType RECORD_TYPE_CNAME =
      RecordType._(0, _omitEnumNames ? '' : 'RECORD_TYPE_CNAME');

  static const $core.List<RecordType> values = <RecordType>[
    RECORD_TYPE_CNAME,
  ];

  static final $core.List<RecordType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static RecordType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecordType._(super.value, super.name);
}

class RenewalEligibility extends $pb.ProtobufEnum {
  static const RenewalEligibility RENEWAL_ELIGIBILITY_ELIGIBLE =
      RenewalEligibility._(
          0, _omitEnumNames ? '' : 'RENEWAL_ELIGIBILITY_ELIGIBLE');
  static const RenewalEligibility RENEWAL_ELIGIBILITY_INELIGIBLE =
      RenewalEligibility._(
          1, _omitEnumNames ? '' : 'RENEWAL_ELIGIBILITY_INELIGIBLE');

  static const $core.List<RenewalEligibility> values = <RenewalEligibility>[
    RENEWAL_ELIGIBILITY_ELIGIBLE,
    RENEWAL_ELIGIBILITY_INELIGIBLE,
  ];

  static final $core.List<RenewalEligibility?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RenewalEligibility? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RenewalEligibility._(super.value, super.name);
}

class RenewalStatus extends $pb.ProtobufEnum {
  static const RenewalStatus RENEWAL_STATUS_PENDING_AUTO_RENEWAL =
      RenewalStatus._(
          0, _omitEnumNames ? '' : 'RENEWAL_STATUS_PENDING_AUTO_RENEWAL');
  static const RenewalStatus RENEWAL_STATUS_PENDING_VALIDATION =
      RenewalStatus._(
          1, _omitEnumNames ? '' : 'RENEWAL_STATUS_PENDING_VALIDATION');
  static const RenewalStatus RENEWAL_STATUS_SUCCESS =
      RenewalStatus._(2, _omitEnumNames ? '' : 'RENEWAL_STATUS_SUCCESS');
  static const RenewalStatus RENEWAL_STATUS_FAILED =
      RenewalStatus._(3, _omitEnumNames ? '' : 'RENEWAL_STATUS_FAILED');

  static const $core.List<RenewalStatus> values = <RenewalStatus>[
    RENEWAL_STATUS_PENDING_AUTO_RENEWAL,
    RENEWAL_STATUS_PENDING_VALIDATION,
    RENEWAL_STATUS_SUCCESS,
    RENEWAL_STATUS_FAILED,
  ];

  static final $core.List<RenewalStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static RenewalStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RenewalStatus._(super.value, super.name);
}

class RevocationReason extends $pb.ProtobufEnum {
  static const RevocationReason REVOCATION_REASON_AFFILIATION_CHANGED =
      RevocationReason._(
          0, _omitEnumNames ? '' : 'REVOCATION_REASON_AFFILIATION_CHANGED');
  static const RevocationReason REVOCATION_REASON_SUPERCEDED =
      RevocationReason._(
          1, _omitEnumNames ? '' : 'REVOCATION_REASON_SUPERCEDED');
  static const RevocationReason REVOCATION_REASON_UNSPECIFIED =
      RevocationReason._(
          2, _omitEnumNames ? '' : 'REVOCATION_REASON_UNSPECIFIED');
  static const RevocationReason REVOCATION_REASON_CERTIFICATE_HOLD =
      RevocationReason._(
          3, _omitEnumNames ? '' : 'REVOCATION_REASON_CERTIFICATE_HOLD');
  static const RevocationReason REVOCATION_REASON_REMOVE_FROM_CRL =
      RevocationReason._(
          4, _omitEnumNames ? '' : 'REVOCATION_REASON_REMOVE_FROM_CRL');
  static const RevocationReason REVOCATION_REASON_CESSATION_OF_OPERATION =
      RevocationReason._(
          5, _omitEnumNames ? '' : 'REVOCATION_REASON_CESSATION_OF_OPERATION');
  static const RevocationReason REVOCATION_REASON_PRIVILEGE_WITHDRAWN =
      RevocationReason._(
          6, _omitEnumNames ? '' : 'REVOCATION_REASON_PRIVILEGE_WITHDRAWN');
  static const RevocationReason REVOCATION_REASON_A_A_COMPROMISE =
      RevocationReason._(
          7, _omitEnumNames ? '' : 'REVOCATION_REASON_A_A_COMPROMISE');
  static const RevocationReason REVOCATION_REASON_SUPERSEDED =
      RevocationReason._(
          8, _omitEnumNames ? '' : 'REVOCATION_REASON_SUPERSEDED');
  static const RevocationReason REVOCATION_REASON_CA_COMPROMISE =
      RevocationReason._(
          9, _omitEnumNames ? '' : 'REVOCATION_REASON_CA_COMPROMISE');
  static const RevocationReason REVOCATION_REASON_KEY_COMPROMISE =
      RevocationReason._(
          10, _omitEnumNames ? '' : 'REVOCATION_REASON_KEY_COMPROMISE');

  static const $core.List<RevocationReason> values = <RevocationReason>[
    REVOCATION_REASON_AFFILIATION_CHANGED,
    REVOCATION_REASON_SUPERCEDED,
    REVOCATION_REASON_UNSPECIFIED,
    REVOCATION_REASON_CERTIFICATE_HOLD,
    REVOCATION_REASON_REMOVE_FROM_CRL,
    REVOCATION_REASON_CESSATION_OF_OPERATION,
    REVOCATION_REASON_PRIVILEGE_WITHDRAWN,
    REVOCATION_REASON_A_A_COMPROMISE,
    REVOCATION_REASON_SUPERSEDED,
    REVOCATION_REASON_CA_COMPROMISE,
    REVOCATION_REASON_KEY_COMPROMISE,
  ];

  static final $core.List<RevocationReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 10);
  static RevocationReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RevocationReason._(super.value, super.name);
}

class SortBy extends $pb.ProtobufEnum {
  static const SortBy SORT_BY_CREATED_AT =
      SortBy._(0, _omitEnumNames ? '' : 'SORT_BY_CREATED_AT');

  static const $core.List<SortBy> values = <SortBy>[
    SORT_BY_CREATED_AT,
  ];

  static final $core.List<SortBy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static SortBy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SortBy._(super.value, super.name);
}

class SortOrder extends $pb.ProtobufEnum {
  static const SortOrder SORT_ORDER_ASCENDING =
      SortOrder._(0, _omitEnumNames ? '' : 'SORT_ORDER_ASCENDING');
  static const SortOrder SORT_ORDER_DESCENDING =
      SortOrder._(1, _omitEnumNames ? '' : 'SORT_ORDER_DESCENDING');

  static const $core.List<SortOrder> values = <SortOrder>[
    SORT_ORDER_ASCENDING,
    SORT_ORDER_DESCENDING,
  ];

  static final $core.List<SortOrder?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SortOrder? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SortOrder._(super.value, super.name);
}

class ValidationMethod extends $pb.ProtobufEnum {
  static const ValidationMethod VALIDATION_METHOD_HTTP =
      ValidationMethod._(0, _omitEnumNames ? '' : 'VALIDATION_METHOD_HTTP');
  static const ValidationMethod VALIDATION_METHOD_EMAIL =
      ValidationMethod._(1, _omitEnumNames ? '' : 'VALIDATION_METHOD_EMAIL');
  static const ValidationMethod VALIDATION_METHOD_DNS =
      ValidationMethod._(2, _omitEnumNames ? '' : 'VALIDATION_METHOD_DNS');

  static const $core.List<ValidationMethod> values = <ValidationMethod>[
    VALIDATION_METHOD_HTTP,
    VALIDATION_METHOD_EMAIL,
    VALIDATION_METHOD_DNS,
  ];

  static final $core.List<ValidationMethod?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ValidationMethod? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ValidationMethod._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
