// This is a generated file - do not edit.
//
// Generated from acm.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use certificateExportDescriptor instead')
const CertificateExport$json = {
  '1': 'CertificateExport',
  '2': [
    {'1': 'CERTIFICATE_EXPORT_DISABLED', '2': 0},
    {'1': 'CERTIFICATE_EXPORT_ENABLED', '2': 1},
  ],
};

/// Descriptor for `CertificateExport`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List certificateExportDescriptor = $convert.base64Decode(
    'ChFDZXJ0aWZpY2F0ZUV4cG9ydBIfChtDRVJUSUZJQ0FURV9FWFBPUlRfRElTQUJMRUQQABIeCh'
    'pDRVJUSUZJQ0FURV9FWFBPUlRfRU5BQkxFRBAB');

@$core.Deprecated('Use certificateManagedByDescriptor instead')
const CertificateManagedBy$json = {
  '1': 'CertificateManagedBy',
  '2': [
    {'1': 'CERTIFICATE_MANAGED_BY_CLOUDFRONT', '2': 0},
  ],
};

/// Descriptor for `CertificateManagedBy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List certificateManagedByDescriptor = $convert.base64Decode(
    'ChRDZXJ0aWZpY2F0ZU1hbmFnZWRCeRIlCiFDRVJUSUZJQ0FURV9NQU5BR0VEX0JZX0NMT1VERl'
    'JPTlQQAA==');

@$core.Deprecated('Use certificateStatusDescriptor instead')
const CertificateStatus$json = {
  '1': 'CertificateStatus',
  '2': [
    {'1': 'CERTIFICATE_STATUS_REVOKED', '2': 0},
    {'1': 'CERTIFICATE_STATUS_ISSUED', '2': 1},
    {'1': 'CERTIFICATE_STATUS_PENDING_VALIDATION', '2': 2},
    {'1': 'CERTIFICATE_STATUS_VALIDATION_TIMED_OUT', '2': 3},
    {'1': 'CERTIFICATE_STATUS_EXPIRED', '2': 4},
    {'1': 'CERTIFICATE_STATUS_INACTIVE', '2': 5},
    {'1': 'CERTIFICATE_STATUS_FAILED', '2': 6},
  ],
};

/// Descriptor for `CertificateStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List certificateStatusDescriptor = $convert.base64Decode(
    'ChFDZXJ0aWZpY2F0ZVN0YXR1cxIeChpDRVJUSUZJQ0FURV9TVEFUVVNfUkVWT0tFRBAAEh0KGU'
    'NFUlRJRklDQVRFX1NUQVRVU19JU1NVRUQQARIpCiVDRVJUSUZJQ0FURV9TVEFUVVNfUEVORElO'
    'R19WQUxJREFUSU9OEAISKwonQ0VSVElGSUNBVEVfU1RBVFVTX1ZBTElEQVRJT05fVElNRURfT1'
    'VUEAMSHgoaQ0VSVElGSUNBVEVfU1RBVFVTX0VYUElSRUQQBBIfChtDRVJUSUZJQ0FURV9TVEFU'
    'VVNfSU5BQ1RJVkUQBRIdChlDRVJUSUZJQ0FURV9TVEFUVVNfRkFJTEVEEAY=');

@$core.Deprecated(
    'Use certificateTransparencyLoggingPreferenceDescriptor instead')
const CertificateTransparencyLoggingPreference$json = {
  '1': 'CertificateTransparencyLoggingPreference',
  '2': [
    {'1': 'CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_DISABLED', '2': 0},
    {'1': 'CERTIFICATE_TRANSPARENCY_LOGGING_PREFERENCE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `CertificateTransparencyLoggingPreference`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List certificateTransparencyLoggingPreferenceDescriptor =
    $convert.base64Decode(
        'CihDZXJ0aWZpY2F0ZVRyYW5zcGFyZW5jeUxvZ2dpbmdQcmVmZXJlbmNlEjgKNENFUlRJRklDQV'
        'RFX1RSQU5TUEFSRU5DWV9MT0dHSU5HX1BSRUZFUkVOQ0VfRElTQUJMRUQQABI3CjNDRVJUSUZJ'
        'Q0FURV9UUkFOU1BBUkVOQ1lfTE9HR0lOR19QUkVGRVJFTkNFX0VOQUJMRUQQAQ==');

@$core.Deprecated('Use certificateTypeDescriptor instead')
const CertificateType$json = {
  '1': 'CertificateType',
  '2': [
    {'1': 'CERTIFICATE_TYPE_IMPORTED', '2': 0},
    {'1': 'CERTIFICATE_TYPE_PRIVATE', '2': 1},
    {'1': 'CERTIFICATE_TYPE_AMAZON_ISSUED', '2': 2},
  ],
};

/// Descriptor for `CertificateType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List certificateTypeDescriptor = $convert.base64Decode(
    'Cg9DZXJ0aWZpY2F0ZVR5cGUSHQoZQ0VSVElGSUNBVEVfVFlQRV9JTVBPUlRFRBAAEhwKGENFUl'
    'RJRklDQVRFX1RZUEVfUFJJVkFURRABEiIKHkNFUlRJRklDQVRFX1RZUEVfQU1BWk9OX0lTU1VF'
    'RBAC');

@$core.Deprecated('Use domainStatusDescriptor instead')
const DomainStatus$json = {
  '1': 'DomainStatus',
  '2': [
    {'1': 'DOMAIN_STATUS_PENDING_VALIDATION', '2': 0},
    {'1': 'DOMAIN_STATUS_SUCCESS', '2': 1},
    {'1': 'DOMAIN_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `DomainStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List domainStatusDescriptor = $convert.base64Decode(
    'CgxEb21haW5TdGF0dXMSJAogRE9NQUlOX1NUQVRVU19QRU5ESU5HX1ZBTElEQVRJT04QABIZCh'
    'VET01BSU5fU1RBVFVTX1NVQ0NFU1MQARIYChRET01BSU5fU1RBVFVTX0ZBSUxFRBAC');

@$core.Deprecated('Use extendedKeyUsageNameDescriptor instead')
const ExtendedKeyUsageName$json = {
  '1': 'ExtendedKeyUsageName',
  '2': [
    {'1': 'EXTENDED_KEY_USAGE_NAME_ANY', '2': 0},
    {'1': 'EXTENDED_KEY_USAGE_NAME_IPSEC_END_SYSTEM', '2': 1},
    {'1': 'EXTENDED_KEY_USAGE_NAME_IPSEC_USER', '2': 2},
    {'1': 'EXTENDED_KEY_USAGE_NAME_CUSTOM', '2': 3},
    {'1': 'EXTENDED_KEY_USAGE_NAME_NONE', '2': 4},
    {'1': 'EXTENDED_KEY_USAGE_NAME_EMAIL_PROTECTION', '2': 5},
    {'1': 'EXTENDED_KEY_USAGE_NAME_TIME_STAMPING', '2': 6},
    {'1': 'EXTENDED_KEY_USAGE_NAME_CODE_SIGNING', '2': 7},
    {'1': 'EXTENDED_KEY_USAGE_NAME_TLS_WEB_CLIENT_AUTHENTICATION', '2': 8},
    {'1': 'EXTENDED_KEY_USAGE_NAME_IPSEC_TUNNEL', '2': 9},
    {'1': 'EXTENDED_KEY_USAGE_NAME_OCSP_SIGNING', '2': 10},
    {'1': 'EXTENDED_KEY_USAGE_NAME_TLS_WEB_SERVER_AUTHENTICATION', '2': 11},
  ],
};

/// Descriptor for `ExtendedKeyUsageName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List extendedKeyUsageNameDescriptor = $convert.base64Decode(
    'ChRFeHRlbmRlZEtleVVzYWdlTmFtZRIfChtFWFRFTkRFRF9LRVlfVVNBR0VfTkFNRV9BTlkQAB'
    'IsCihFWFRFTkRFRF9LRVlfVVNBR0VfTkFNRV9JUFNFQ19FTkRfU1lTVEVNEAESJgoiRVhURU5E'
    'RURfS0VZX1VTQUdFX05BTUVfSVBTRUNfVVNFUhACEiIKHkVYVEVOREVEX0tFWV9VU0FHRV9OQU'
    '1FX0NVU1RPTRADEiAKHEVYVEVOREVEX0tFWV9VU0FHRV9OQU1FX05PTkUQBBIsCihFWFRFTkRF'
    'RF9LRVlfVVNBR0VfTkFNRV9FTUFJTF9QUk9URUNUSU9OEAUSKQolRVhURU5ERURfS0VZX1VTQU'
    'dFX05BTUVfVElNRV9TVEFNUElORxAGEigKJEVYVEVOREVEX0tFWV9VU0FHRV9OQU1FX0NPREVf'
    'U0lHTklORxAHEjkKNUVYVEVOREVEX0tFWV9VU0FHRV9OQU1FX1RMU19XRUJfQ0xJRU5UX0FVVE'
    'hFTlRJQ0FUSU9OEAgSKAokRVhURU5ERURfS0VZX1VTQUdFX05BTUVfSVBTRUNfVFVOTkVMEAkS'
    'KAokRVhURU5ERURfS0VZX1VTQUdFX05BTUVfT0NTUF9TSUdOSU5HEAoSOQo1RVhURU5ERURfS0'
    'VZX1VTQUdFX05BTUVfVExTX1dFQl9TRVJWRVJfQVVUSEVOVElDQVRJT04QCw==');

@$core.Deprecated('Use failureReasonDescriptor instead')
const FailureReason$json = {
  '1': 'FailureReason',
  '2': [
    {'1': 'FAILURE_REASON_PCA_INVALID_ARGS', '2': 0},
    {'1': 'FAILURE_REASON_PCA_INVALID_STATE', '2': 1},
    {'1': 'FAILURE_REASON_OTHER', '2': 2},
    {'1': 'FAILURE_REASON_CAA_ERROR', '2': 3},
    {'1': 'FAILURE_REASON_ADDITIONAL_VERIFICATION_REQUIRED', '2': 4},
    {'1': 'FAILURE_REASON_PCA_NAME_CONSTRAINTS_VALIDATION', '2': 5},
    {'1': 'FAILURE_REASON_INVALID_PUBLIC_DOMAIN', '2': 6},
    {'1': 'FAILURE_REASON_DOMAIN_NOT_ALLOWED', '2': 7},
    {'1': 'FAILURE_REASON_NO_AVAILABLE_CONTACTS', '2': 8},
    {'1': 'FAILURE_REASON_PCA_INVALID_DURATION', '2': 9},
    {'1': 'FAILURE_REASON_PCA_REQUEST_FAILED', '2': 10},
    {'1': 'FAILURE_REASON_DOMAIN_VALIDATION_DENIED', '2': 11},
    {'1': 'FAILURE_REASON_PCA_INVALID_ARN', '2': 12},
    {'1': 'FAILURE_REASON_PCA_RESOURCE_NOT_FOUND', '2': 13},
    {'1': 'FAILURE_REASON_PCA_ACCESS_DENIED', '2': 14},
    {'1': 'FAILURE_REASON_PCA_LIMIT_EXCEEDED', '2': 15},
    {'1': 'FAILURE_REASON_SLR_NOT_FOUND', '2': 16},
  ],
};

/// Descriptor for `FailureReason`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List failureReasonDescriptor = $convert.base64Decode(
    'Cg1GYWlsdXJlUmVhc29uEiMKH0ZBSUxVUkVfUkVBU09OX1BDQV9JTlZBTElEX0FSR1MQABIkCi'
    'BGQUlMVVJFX1JFQVNPTl9QQ0FfSU5WQUxJRF9TVEFURRABEhgKFEZBSUxVUkVfUkVBU09OX09U'
    'SEVSEAISHAoYRkFJTFVSRV9SRUFTT05fQ0FBX0VSUk9SEAMSMwovRkFJTFVSRV9SRUFTT05fQU'
    'RESVRJT05BTF9WRVJJRklDQVRJT05fUkVRVUlSRUQQBBIyCi5GQUlMVVJFX1JFQVNPTl9QQ0Ff'
    'TkFNRV9DT05TVFJBSU5UU19WQUxJREFUSU9OEAUSKAokRkFJTFVSRV9SRUFTT05fSU5WQUxJRF'
    '9QVUJMSUNfRE9NQUlOEAYSJQohRkFJTFVSRV9SRUFTT05fRE9NQUlOX05PVF9BTExPV0VEEAcS'
    'KAokRkFJTFVSRV9SRUFTT05fTk9fQVZBSUxBQkxFX0NPTlRBQ1RTEAgSJwojRkFJTFVSRV9SRU'
    'FTT05fUENBX0lOVkFMSURfRFVSQVRJT04QCRIlCiFGQUlMVVJFX1JFQVNPTl9QQ0FfUkVRVUVT'
    'VF9GQUlMRUQQChIrCidGQUlMVVJFX1JFQVNPTl9ET01BSU5fVkFMSURBVElPTl9ERU5JRUQQCx'
    'IiCh5GQUlMVVJFX1JFQVNPTl9QQ0FfSU5WQUxJRF9BUk4QDBIpCiVGQUlMVVJFX1JFQVNPTl9Q'
    'Q0FfUkVTT1VSQ0VfTk9UX0ZPVU5EEA0SJAogRkFJTFVSRV9SRUFTT05fUENBX0FDQ0VTU19ERU'
    '5JRUQQDhIlCiFGQUlMVVJFX1JFQVNPTl9QQ0FfTElNSVRfRVhDRUVERUQQDxIgChxGQUlMVVJF'
    'X1JFQVNPTl9TTFJfTk9UX0ZPVU5EEBA=');

@$core.Deprecated('Use keyAlgorithmDescriptor instead')
const KeyAlgorithm$json = {
  '1': 'KeyAlgorithm',
  '2': [
    {'1': 'KEY_ALGORITHM_EC_PRIME256V1', '2': 0},
    {'1': 'KEY_ALGORITHM_RSA_2048', '2': 1},
    {'1': 'KEY_ALGORITHM_RSA_3072', '2': 2},
    {'1': 'KEY_ALGORITHM_EC_SECP521R1', '2': 3},
    {'1': 'KEY_ALGORITHM_RSA_1024', '2': 4},
    {'1': 'KEY_ALGORITHM_RSA_4096', '2': 5},
    {'1': 'KEY_ALGORITHM_EC_SECP384R1', '2': 6},
  ],
};

/// Descriptor for `KeyAlgorithm`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyAlgorithmDescriptor = $convert.base64Decode(
    'CgxLZXlBbGdvcml0aG0SHwobS0VZX0FMR09SSVRITV9FQ19QUklNRTI1NlYxEAASGgoWS0VZX0'
    'FMR09SSVRITV9SU0FfMjA0OBABEhoKFktFWV9BTEdPUklUSE1fUlNBXzMwNzIQAhIeChpLRVlf'
    'QUxHT1JJVEhNX0VDX1NFQ1A1MjFSMRADEhoKFktFWV9BTEdPUklUSE1fUlNBXzEwMjQQBBIaCh'
    'ZLRVlfQUxHT1JJVEhNX1JTQV80MDk2EAUSHgoaS0VZX0FMR09SSVRITV9FQ19TRUNQMzg0UjEQ'
    'Bg==');

@$core.Deprecated('Use keyUsageNameDescriptor instead')
const KeyUsageName$json = {
  '1': 'KeyUsageName',
  '2': [
    {'1': 'KEY_USAGE_NAME_ANY', '2': 0},
    {'1': 'KEY_USAGE_NAME_CERTIFICATE_SIGNING', '2': 1},
    {'1': 'KEY_USAGE_NAME_CUSTOM', '2': 2},
    {'1': 'KEY_USAGE_NAME_DATA_ENCIPHERMENT', '2': 3},
    {'1': 'KEY_USAGE_NAME_CRL_SIGNING', '2': 4},
    {'1': 'KEY_USAGE_NAME_KEY_AGREEMENT', '2': 5},
    {'1': 'KEY_USAGE_NAME_DECIPHER_ONLY', '2': 6},
    {'1': 'KEY_USAGE_NAME_ENCHIPER_ONLY', '2': 7},
    {'1': 'KEY_USAGE_NAME_NON_REPUDATION', '2': 8},
    {'1': 'KEY_USAGE_NAME_KEY_ENCIPHERMENT', '2': 9},
    {'1': 'KEY_USAGE_NAME_DIGITAL_SIGNATURE', '2': 10},
  ],
};

/// Descriptor for `KeyUsageName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyUsageNameDescriptor = $convert.base64Decode(
    'CgxLZXlVc2FnZU5hbWUSFgoSS0VZX1VTQUdFX05BTUVfQU5ZEAASJgoiS0VZX1VTQUdFX05BTU'
    'VfQ0VSVElGSUNBVEVfU0lHTklORxABEhkKFUtFWV9VU0FHRV9OQU1FX0NVU1RPTRACEiQKIEtF'
    'WV9VU0FHRV9OQU1FX0RBVEFfRU5DSVBIRVJNRU5UEAMSHgoaS0VZX1VTQUdFX05BTUVfQ1JMX1'
    'NJR05JTkcQBBIgChxLRVlfVVNBR0VfTkFNRV9LRVlfQUdSRUVNRU5UEAUSIAocS0VZX1VTQUdF'
    'X05BTUVfREVDSVBIRVJfT05MWRAGEiAKHEtFWV9VU0FHRV9OQU1FX0VOQ0hJUEVSX09OTFkQBx'
    'IhCh1LRVlfVVNBR0VfTkFNRV9OT05fUkVQVURBVElPThAIEiMKH0tFWV9VU0FHRV9OQU1FX0tF'
    'WV9FTkNJUEhFUk1FTlQQCRIkCiBLRVlfVVNBR0VfTkFNRV9ESUdJVEFMX1NJR05BVFVSRRAK');

@$core.Deprecated('Use recordTypeDescriptor instead')
const RecordType$json = {
  '1': 'RecordType',
  '2': [
    {'1': 'RECORD_TYPE_CNAME', '2': 0},
  ],
};

/// Descriptor for `RecordType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recordTypeDescriptor =
    $convert.base64Decode('CgpSZWNvcmRUeXBlEhUKEVJFQ09SRF9UWVBFX0NOQU1FEAA=');

@$core.Deprecated('Use renewalEligibilityDescriptor instead')
const RenewalEligibility$json = {
  '1': 'RenewalEligibility',
  '2': [
    {'1': 'RENEWAL_ELIGIBILITY_ELIGIBLE', '2': 0},
    {'1': 'RENEWAL_ELIGIBILITY_INELIGIBLE', '2': 1},
  ],
};

/// Descriptor for `RenewalEligibility`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List renewalEligibilityDescriptor = $convert.base64Decode(
    'ChJSZW5ld2FsRWxpZ2liaWxpdHkSIAocUkVORVdBTF9FTElHSUJJTElUWV9FTElHSUJMRRAAEi'
    'IKHlJFTkVXQUxfRUxJR0lCSUxJVFlfSU5FTElHSUJMRRAB');

@$core.Deprecated('Use renewalStatusDescriptor instead')
const RenewalStatus$json = {
  '1': 'RenewalStatus',
  '2': [
    {'1': 'RENEWAL_STATUS_PENDING_AUTO_RENEWAL', '2': 0},
    {'1': 'RENEWAL_STATUS_PENDING_VALIDATION', '2': 1},
    {'1': 'RENEWAL_STATUS_SUCCESS', '2': 2},
    {'1': 'RENEWAL_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `RenewalStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List renewalStatusDescriptor = $convert.base64Decode(
    'Cg1SZW5ld2FsU3RhdHVzEicKI1JFTkVXQUxfU1RBVFVTX1BFTkRJTkdfQVVUT19SRU5FV0FMEA'
    'ASJQohUkVORVdBTF9TVEFUVVNfUEVORElOR19WQUxJREFUSU9OEAESGgoWUkVORVdBTF9TVEFU'
    'VVNfU1VDQ0VTUxACEhkKFVJFTkVXQUxfU1RBVFVTX0ZBSUxFRBAD');

@$core.Deprecated('Use revocationReasonDescriptor instead')
const RevocationReason$json = {
  '1': 'RevocationReason',
  '2': [
    {'1': 'REVOCATION_REASON_AFFILIATION_CHANGED', '2': 0},
    {'1': 'REVOCATION_REASON_SUPERCEDED', '2': 1},
    {'1': 'REVOCATION_REASON_UNSPECIFIED', '2': 2},
    {'1': 'REVOCATION_REASON_CERTIFICATE_HOLD', '2': 3},
    {'1': 'REVOCATION_REASON_REMOVE_FROM_CRL', '2': 4},
    {'1': 'REVOCATION_REASON_CESSATION_OF_OPERATION', '2': 5},
    {'1': 'REVOCATION_REASON_PRIVILEGE_WITHDRAWN', '2': 6},
    {'1': 'REVOCATION_REASON_A_A_COMPROMISE', '2': 7},
    {'1': 'REVOCATION_REASON_SUPERSEDED', '2': 8},
    {'1': 'REVOCATION_REASON_CA_COMPROMISE', '2': 9},
    {'1': 'REVOCATION_REASON_KEY_COMPROMISE', '2': 10},
  ],
};

/// Descriptor for `RevocationReason`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List revocationReasonDescriptor = $convert.base64Decode(
    'ChBSZXZvY2F0aW9uUmVhc29uEikKJVJFVk9DQVRJT05fUkVBU09OX0FGRklMSUFUSU9OX0NIQU'
    '5HRUQQABIgChxSRVZPQ0FUSU9OX1JFQVNPTl9TVVBFUkNFREVEEAESIQodUkVWT0NBVElPTl9S'
    'RUFTT05fVU5TUEVDSUZJRUQQAhImCiJSRVZPQ0FUSU9OX1JFQVNPTl9DRVJUSUZJQ0FURV9IT0'
    'xEEAMSJQohUkVWT0NBVElPTl9SRUFTT05fUkVNT1ZFX0ZST01fQ1JMEAQSLAooUkVWT0NBVElP'
    'Tl9SRUFTT05fQ0VTU0FUSU9OX09GX09QRVJBVElPThAFEikKJVJFVk9DQVRJT05fUkVBU09OX1'
    'BSSVZJTEVHRV9XSVRIRFJBV04QBhIkCiBSRVZPQ0FUSU9OX1JFQVNPTl9BX0FfQ09NUFJPTUlT'
    'RRAHEiAKHFJFVk9DQVRJT05fUkVBU09OX1NVUEVSU0VERUQQCBIjCh9SRVZPQ0FUSU9OX1JFQV'
    'NPTl9DQV9DT01QUk9NSVNFEAkSJAogUkVWT0NBVElPTl9SRUFTT05fS0VZX0NPTVBST01JU0UQ'
    'Cg==');

@$core.Deprecated('Use sortByDescriptor instead')
const SortBy$json = {
  '1': 'SortBy',
  '2': [
    {'1': 'SORT_BY_CREATED_AT', '2': 0},
  ],
};

/// Descriptor for `SortBy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sortByDescriptor =
    $convert.base64Decode('CgZTb3J0QnkSFgoSU09SVF9CWV9DUkVBVEVEX0FUEAA=');

@$core.Deprecated('Use sortOrderDescriptor instead')
const SortOrder$json = {
  '1': 'SortOrder',
  '2': [
    {'1': 'SORT_ORDER_ASCENDING', '2': 0},
    {'1': 'SORT_ORDER_DESCENDING', '2': 1},
  ],
};

/// Descriptor for `SortOrder`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sortOrderDescriptor = $convert.base64Decode(
    'CglTb3J0T3JkZXISGAoUU09SVF9PUkRFUl9BU0NFTkRJTkcQABIZChVTT1JUX09SREVSX0RFU0'
    'NFTkRJTkcQAQ==');

@$core.Deprecated('Use validationMethodDescriptor instead')
const ValidationMethod$json = {
  '1': 'ValidationMethod',
  '2': [
    {'1': 'VALIDATION_METHOD_HTTP', '2': 0},
    {'1': 'VALIDATION_METHOD_EMAIL', '2': 1},
    {'1': 'VALIDATION_METHOD_DNS', '2': 2},
  ],
};

/// Descriptor for `ValidationMethod`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List validationMethodDescriptor = $convert.base64Decode(
    'ChBWYWxpZGF0aW9uTWV0aG9kEhoKFlZBTElEQVRJT05fTUVUSE9EX0hUVFAQABIbChdWQUxJRE'
    'FUSU9OX01FVEhPRF9FTUFJTBABEhkKFVZBTElEQVRJT05fTUVUSE9EX0ROUxAC');

@$core.Deprecated('Use accessDeniedExceptionDescriptor instead')
const AccessDeniedException$json = {
  '1': 'AccessDeniedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedExceptionDescriptor = $convert.base64Decode(
    'ChVBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use addTagsToCertificateRequestDescriptor instead')
const AddTagsToCertificateRequest$json = {
  '1': 'AddTagsToCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.acm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `AddTagsToCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addTagsToCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChtBZGRUYWdzVG9DZXJ0aWZpY2F0ZVJlcXVlc3QSKQoOY2VydGlmaWNhdGVhcm4Y+MqZLCABKA'
        'lSDmNlcnRpZmljYXRlYXJuEiAKBHRhZ3MYwcH2tQEgAygLMgguYWNtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use certificateDetailDescriptor instead')
const CertificateDetail$json = {
  '1': 'CertificateDetail',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'certificateauthorityarn',
      '3': 266069181,
      '4': 1,
      '5': 9,
      '10': 'certificateauthorityarn'
    },
    {'1': 'createdat', '3': 258192751, '4': 1, '5': 9, '10': 'createdat'},
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'domainvalidationoptions',
      '3': 335573705,
      '4': 3,
      '5': 11,
      '6': '.acm.DomainValidation',
      '10': 'domainvalidationoptions'
    },
    {
      '1': 'extendedkeyusages',
      '3': 531267688,
      '4': 3,
      '5': 11,
      '6': '.acm.ExtendedKeyUsage',
      '10': 'extendedkeyusages'
    },
    {
      '1': 'failurereason',
      '3': 232322142,
      '4': 1,
      '5': 14,
      '6': '.acm.FailureReason',
      '10': 'failurereason'
    },
    {'1': 'importedat', '3': 348649225, '4': 1, '5': 9, '10': 'importedat'},
    {'1': 'inuseby', '3': 330307273, '4': 3, '5': 9, '10': 'inuseby'},
    {'1': 'issuedat', '3': 17449786, '4': 1, '5': 9, '10': 'issuedat'},
    {'1': 'issuer', '3': 528708823, '4': 1, '5': 9, '10': 'issuer'},
    {
      '1': 'keyalgorithm',
      '3': 452317818,
      '4': 1,
      '5': 14,
      '6': '.acm.KeyAlgorithm',
      '10': 'keyalgorithm'
    },
    {
      '1': 'keyusages',
      '3': 345433681,
      '4': 3,
      '5': 11,
      '6': '.acm.KeyUsage',
      '10': 'keyusages'
    },
    {
      '1': 'managedby',
      '3': 455511232,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateManagedBy',
      '10': 'managedby'
    },
    {'1': 'notafter', '3': 287678033, '4': 1, '5': 9, '10': 'notafter'},
    {'1': 'notbefore', '3': 459074038, '4': 1, '5': 9, '10': 'notbefore'},
    {
      '1': 'options',
      '3': 356388166,
      '4': 1,
      '5': 11,
      '6': '.acm.CertificateOptions',
      '10': 'options'
    },
    {
      '1': 'renewaleligibility',
      '3': 172871849,
      '4': 1,
      '5': 14,
      '6': '.acm.RenewalEligibility',
      '10': 'renewaleligibility'
    },
    {
      '1': 'renewalsummary',
      '3': 125255166,
      '4': 1,
      '5': 11,
      '6': '.acm.RenewalSummary',
      '10': 'renewalsummary'
    },
    {
      '1': 'revocationreason',
      '3': 331836358,
      '4': 1,
      '5': 14,
      '6': '.acm.RevocationReason',
      '10': 'revocationreason'
    },
    {'1': 'revokedat', '3': 63251417, '4': 1, '5': 9, '10': 'revokedat'},
    {'1': 'serial', '3': 143954586, '4': 1, '5': 9, '10': 'serial'},
    {
      '1': 'signaturealgorithm',
      '3': 476410739,
      '4': 1,
      '5': 9,
      '10': 'signaturealgorithm'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateStatus',
      '10': 'status'
    },
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
    {
      '1': 'subjectalternativenames',
      '3': 109998119,
      '4': 3,
      '5': 9,
      '10': 'subjectalternativenames'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `CertificateDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateDetailDescriptor = $convert.base64Decode(
    'ChFDZXJ0aWZpY2F0ZURldGFpbBIpCg5jZXJ0aWZpY2F0ZWFybhj4ypksIAEoCVIOY2VydGlmaW'
    'NhdGVhcm4SOwoXY2VydGlmaWNhdGVhdXRob3JpdHlhcm4YvcnvfiABKAlSF2NlcnRpZmljYXRl'
    'YXV0aG9yaXR5YXJuEh8KCWNyZWF0ZWRhdBjv6o57IAEoCVIJY3JlYXRlZGF0EiEKCmRvbWFpbm'
    '5hbWUY6834XCABKAlSCmRvbWFpbm5hbWUSUwoXZG9tYWludmFsaWRhdGlvbm9wdGlvbnMYyeWB'
    'oAEgAygLMhUuYWNtLkRvbWFpblZhbGlkYXRpb25SF2RvbWFpbnZhbGlkYXRpb25vcHRpb25zEk'
    'cKEWV4dGVuZGVka2V5dXNhZ2VzGOiAqv0BIAMoCzIVLmFjbS5FeHRlbmRlZEtleVVzYWdlUhFl'
    'eHRlbmRlZGtleXVzYWdlcxI7Cg1mYWlsdXJlcmVhc29uGN7o424gASgOMhIuYWNtLkZhaWx1cm'
    'VSZWFzb25SDWZhaWx1cmVyZWFzb24SIgoKaW1wb3J0ZWRhdBiJ7p+mASABKAlSCmltcG9ydGVk'
    'YXQSHAoHaW51c2VieRjJrcCdASADKAlSB2ludXNlYnkSHQoIaXNzdWVkYXQYuoapCCABKAlSCG'
    'lzc3VlZGF0EhoKBmlzc3VlchjX6Y38ASABKAlSBmlzc3VlchI5CgxrZXlhbGdvcml0aG0Y+qTX'
    '1wEgASgOMhEuYWNtLktleUFsZ29yaXRobVIMa2V5YWxnb3JpdGhtEi8KCWtleXVzYWdlcxjRzN'
    'ukASADKAsyDS5hY20uS2V5VXNhZ2VSCWtleXVzYWdlcxI7CgltYW5hZ2VkYnkYwJma2QEgASgO'
    'MhkuYWNtLkNlcnRpZmljYXRlTWFuYWdlZEJ5UgltYW5hZ2VkYnkSHgoIbm90YWZ0ZXIY0byWiQ'
    'EgASgJUghub3RhZnRlchIgCglub3RiZWZvcmUY9tPz2gEgASgJUglub3RiZWZvcmUSNQoHb3B0'
    'aW9ucxjGmvipASABKAsyFy5hY20uQ2VydGlmaWNhdGVPcHRpb25zUgdvcHRpb25zEkoKEnJlbm'
    'V3YWxlbGlnaWJpbGl0eRipobdSIAEoDjIXLmFjbS5SZW5ld2FsRWxpZ2liaWxpdHlSEnJlbmV3'
    'YWxlbGlnaWJpbGl0eRI+Cg5yZW5ld2Fsc3VtbWFyeRj++9w7IAEoCzITLmFjbS5SZW5ld2FsU3'
    'VtbWFyeVIOcmVuZXdhbHN1bW1hcnkSRQoQcmV2b2NhdGlvbnJlYXNvbhjG152eASABKA4yFS5h'
    'Y20uUmV2b2NhdGlvblJlYXNvblIQcmV2b2NhdGlvbnJlYXNvbhIfCglyZXZva2VkYXQY2ceUHi'
    'ABKAlSCXJldm9rZWRhdBIZCgZzZXJpYWwYmqXSRCABKAlSBnNlcmlhbBIyChJzaWduYXR1cmVh'
    'bGdvcml0aG0Y8+aV4wEgASgJUhJzaWduYXR1cmVhbGdvcml0aG0SMQoGc3RhdHVzGJDk+wIgAS'
    'gOMhYuYWNtLkNlcnRpZmljYXRlU3RhdHVzUgZzdGF0dXMSGwoHc3ViamVjdBjwyeQDIAEoCVIH'
    'c3ViamVjdBI7ChdzdWJqZWN0YWx0ZXJuYXRpdmVuYW1lcxin4Lk0IAMoCVIXc3ViamVjdGFsdG'
    'VybmF0aXZlbmFtZXMSLAoEdHlwZRjuoNeKASABKA4yFC5hY20uQ2VydGlmaWNhdGVUeXBlUgR0'
    'eXBl');

@$core.Deprecated('Use certificateOptionsDescriptor instead')
const CertificateOptions$json = {
  '1': 'CertificateOptions',
  '2': [
    {
      '1': 'certificatetransparencyloggingpreference',
      '3': 414636075,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateTransparencyLoggingPreference',
      '10': 'certificatetransparencyloggingpreference'
    },
    {
      '1': 'export',
      '3': 140724692,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateExport',
      '10': 'export'
    },
  ],
};

/// Descriptor for `CertificateOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateOptionsDescriptor = $convert.base64Decode(
    'ChJDZXJ0aWZpY2F0ZU9wdGlvbnMSjQEKKGNlcnRpZmljYXRldHJhbnNwYXJlbmN5bG9nZ2luZ3'
    'ByZWZlcmVuY2UYq7DbxQEgASgOMi0uYWNtLkNlcnRpZmljYXRlVHJhbnNwYXJlbmN5TG9nZ2lu'
    'Z1ByZWZlcmVuY2VSKGNlcnRpZmljYXRldHJhbnNwYXJlbmN5bG9nZ2luZ3ByZWZlcmVuY2USMQ'
    'oGZXhwb3J0GNSTjUMgASgOMhYuYWNtLkNlcnRpZmljYXRlRXhwb3J0UgZleHBvcnQ=');

@$core.Deprecated('Use certificateSummaryDescriptor instead')
const CertificateSummary$json = {
  '1': 'CertificateSummary',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {'1': 'createdat', '3': 258192751, '4': 1, '5': 9, '10': 'createdat'},
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'exportoption',
      '3': 19500687,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateExport',
      '10': 'exportoption'
    },
    {'1': 'exported', '3': 491164947, '4': 1, '5': 8, '10': 'exported'},
    {
      '1': 'extendedkeyusages',
      '3': 531267688,
      '4': 3,
      '5': 14,
      '6': '.acm.ExtendedKeyUsageName',
      '10': 'extendedkeyusages'
    },
    {
      '1': 'hasadditionalsubjectalternativenames',
      '3': 389669028,
      '4': 1,
      '5': 8,
      '10': 'hasadditionalsubjectalternativenames'
    },
    {'1': 'importedat', '3': 348649225, '4': 1, '5': 9, '10': 'importedat'},
    {'1': 'inuse', '3': 398346234, '4': 1, '5': 8, '10': 'inuse'},
    {'1': 'issuedat', '3': 17449786, '4': 1, '5': 9, '10': 'issuedat'},
    {
      '1': 'keyalgorithm',
      '3': 452317818,
      '4': 1,
      '5': 14,
      '6': '.acm.KeyAlgorithm',
      '10': 'keyalgorithm'
    },
    {
      '1': 'keyusages',
      '3': 345433681,
      '4': 3,
      '5': 14,
      '6': '.acm.KeyUsageName',
      '10': 'keyusages'
    },
    {
      '1': 'managedby',
      '3': 455511232,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateManagedBy',
      '10': 'managedby'
    },
    {'1': 'notafter', '3': 287678033, '4': 1, '5': 9, '10': 'notafter'},
    {'1': 'notbefore', '3': 459074038, '4': 1, '5': 9, '10': 'notbefore'},
    {
      '1': 'renewaleligibility',
      '3': 172871849,
      '4': 1,
      '5': 14,
      '6': '.acm.RenewalEligibility',
      '10': 'renewaleligibility'
    },
    {'1': 'revokedat', '3': 63251417, '4': 1, '5': 9, '10': 'revokedat'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateStatus',
      '10': 'status'
    },
    {
      '1': 'subjectalternativenamesummaries',
      '3': 248841448,
      '4': 3,
      '5': 9,
      '10': 'subjectalternativenamesummaries'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `CertificateSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateSummaryDescriptor = $convert.base64Decode(
    'ChJDZXJ0aWZpY2F0ZVN1bW1hcnkSKQoOY2VydGlmaWNhdGVhcm4Y+MqZLCABKAlSDmNlcnRpZm'
    'ljYXRlYXJuEh8KCWNyZWF0ZWRhdBjv6o57IAEoCVIJY3JlYXRlZGF0EiEKCmRvbWFpbm5hbWUY'
    '6834XCABKAlSCmRvbWFpbm5hbWUSPQoMZXhwb3J0b3B0aW9uGI+dpgkgASgOMhYuYWNtLkNlcn'
    'RpZmljYXRlRXhwb3J0UgxleHBvcnRvcHRpb24SHgoIZXhwb3J0ZWQYk6qa6gEgASgIUghleHBv'
    'cnRlZBJLChFleHRlbmRlZGtleXVzYWdlcxjogKr9ASADKA4yGS5hY20uRXh0ZW5kZWRLZXlVc2'
    'FnZU5hbWVSEWV4dGVuZGVka2V5dXNhZ2VzElYKJGhhc2FkZGl0aW9uYWxzdWJqZWN0YWx0ZXJu'
    'YXRpdmVuYW1lcxikwee5ASABKAhSJGhhc2FkZGl0aW9uYWxzdWJqZWN0YWx0ZXJuYXRpdmVuYW'
    '1lcxIiCgppbXBvcnRlZGF0GInun6YBIAEoCVIKaW1wb3J0ZWRhdBIYCgVpbnVzZRj6j/m9ASAB'
    'KAhSBWludXNlEh0KCGlzc3VlZGF0GLqGqQggASgJUghpc3N1ZWRhdBI5CgxrZXlhbGdvcml0aG'
    '0Y+qTX1wEgASgOMhEuYWNtLktleUFsZ29yaXRobVIMa2V5YWxnb3JpdGhtEjMKCWtleXVzYWdl'
    'cxjRzNukASADKA4yES5hY20uS2V5VXNhZ2VOYW1lUglrZXl1c2FnZXMSOwoJbWFuYWdlZGJ5GM'
    'CZmtkBIAEoDjIZLmFjbS5DZXJ0aWZpY2F0ZU1hbmFnZWRCeVIJbWFuYWdlZGJ5Eh4KCG5vdGFm'
    'dGVyGNG8lokBIAEoCVIIbm90YWZ0ZXISIAoJbm90YmVmb3JlGPbT89oBIAEoCVIJbm90YmVmb3'
    'JlEkoKEnJlbmV3YWxlbGlnaWJpbGl0eRipobdSIAEoDjIXLmFjbS5SZW5ld2FsRWxpZ2liaWxp'
    'dHlSEnJlbmV3YWxlbGlnaWJpbGl0eRIfCglyZXZva2VkYXQY2ceUHiABKAlSCXJldm9rZWRhdB'
    'IxCgZzdGF0dXMYkOT7AiABKA4yFi5hY20uQ2VydGlmaWNhdGVTdGF0dXNSBnN0YXR1cxJLCh9z'
    'dWJqZWN0YWx0ZXJuYXRpdmVuYW1lc3VtbWFyaWVzGOiJ1HYgAygJUh9zdWJqZWN0YWx0ZXJuYX'
    'RpdmVuYW1lc3VtbWFyaWVzEiwKBHR5cGUY7qDXigEgASgOMhQuYWNtLkNlcnRpZmljYXRlVHlw'
    'ZVIEdHlwZQ==');

@$core.Deprecated('Use conflictExceptionDescriptor instead')
const ConflictException$json = {
  '1': 'ConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictExceptionDescriptor = $convert.base64Decode(
    'ChFDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use deleteCertificateRequestDescriptor instead')
const DeleteCertificateRequest$json = {
  '1': 'DeleteCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `DeleteCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVDZXJ0aWZpY2F0ZVJlcXVlc3QSKQoOY2VydGlmaWNhdGVhcm4Y+MqZLCABKAlSDm'
        'NlcnRpZmljYXRlYXJu');

@$core.Deprecated('Use describeCertificateRequestDescriptor instead')
const DescribeCertificateRequest$json = {
  '1': 'DescribeCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `DescribeCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZUNlcnRpZmljYXRlUmVxdWVzdBIpCg5jZXJ0aWZpY2F0ZWFybhj4ypksIAEoCV'
        'IOY2VydGlmaWNhdGVhcm4=');

@$core.Deprecated('Use describeCertificateResponseDescriptor instead')
const DescribeCertificateResponse$json = {
  '1': 'DescribeCertificateResponse',
  '2': [
    {
      '1': 'certificate',
      '3': 198060817,
      '4': 1,
      '5': 11,
      '6': '.acm.CertificateDetail',
      '10': 'certificate'
    },
  ],
};

/// Descriptor for `DescribeCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeCertificateResponseDescriptor =
    $convert.base64Decode(
        'ChtEZXNjcmliZUNlcnRpZmljYXRlUmVzcG9uc2USOwoLY2VydGlmaWNhdGUYkda4XiABKAsyFi'
        '5hY20uQ2VydGlmaWNhdGVEZXRhaWxSC2NlcnRpZmljYXRl');

@$core.Deprecated('Use domainValidationDescriptor instead')
const DomainValidation$json = {
  '1': 'DomainValidation',
  '2': [
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'httpredirect',
      '3': 375515090,
      '4': 1,
      '5': 11,
      '6': '.acm.HttpRedirect',
      '10': 'httpredirect'
    },
    {
      '1': 'resourcerecord',
      '3': 207153213,
      '4': 1,
      '5': 11,
      '6': '.acm.ResourceRecord',
      '10': 'resourcerecord'
    },
    {
      '1': 'validationdomain',
      '3': 161300799,
      '4': 1,
      '5': 9,
      '10': 'validationdomain'
    },
    {
      '1': 'validationemails',
      '3': 375615664,
      '4': 3,
      '5': 9,
      '10': 'validationemails'
    },
    {
      '1': 'validationmethod',
      '3': 58092520,
      '4': 1,
      '5': 14,
      '6': '.acm.ValidationMethod',
      '10': 'validationmethod'
    },
    {
      '1': 'validationstatus',
      '3': 426907749,
      '4': 1,
      '5': 14,
      '6': '.acm.DomainStatus',
      '10': 'validationstatus'
    },
  ],
};

/// Descriptor for `DomainValidation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainValidationDescriptor = $convert.base64Decode(
    'ChBEb21haW5WYWxpZGF0aW9uEiEKCmRvbWFpbm5hbWUY6834XCABKAlSCmRvbWFpbm5hbWUSOQ'
    'oMaHR0cHJlZGlyZWN0GNLPh7MBIAEoCzIRLmFjbS5IdHRwUmVkaXJlY3RSDGh0dHByZWRpcmVj'
    'dBI+Cg5yZXNvdXJjZXJlY29yZBi90ONiIAEoCzITLmFjbS5SZXNvdXJjZVJlY29yZFIOcmVzb3'
    'VyY2VyZWNvcmQSLQoQdmFsaWRhdGlvbmRvbWFpbhi/gvVMIAEoCVIQdmFsaWRhdGlvbmRvbWFp'
    'bhIuChB2YWxpZGF0aW9uZW1haWxzGLDhjbMBIAMoCVIQdmFsaWRhdGlvbmVtYWlscxJEChB2YW'
    'xpZGF0aW9ubWV0aG9kGOjX2RsgASgOMhUuYWNtLlZhbGlkYXRpb25NZXRob2RSEHZhbGlkYXRp'
    'b25tZXRob2QSQQoQdmFsaWRhdGlvbnN0YXR1cxjlsMjLASABKA4yES5hY20uRG9tYWluU3RhdH'
    'VzUhB2YWxpZGF0aW9uc3RhdHVz');

@$core.Deprecated('Use domainValidationOptionDescriptor instead')
const DomainValidationOption$json = {
  '1': 'DomainValidationOption',
  '2': [
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'validationdomain',
      '3': 161300799,
      '4': 1,
      '5': 9,
      '10': 'validationdomain'
    },
  ],
};

/// Descriptor for `DomainValidationOption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainValidationOptionDescriptor = $convert.base64Decode(
    'ChZEb21haW5WYWxpZGF0aW9uT3B0aW9uEiEKCmRvbWFpbm5hbWUY6834XCABKAlSCmRvbWFpbm'
    '5hbWUSLQoQdmFsaWRhdGlvbmRvbWFpbhi/gvVMIAEoCVIQdmFsaWRhdGlvbmRvbWFpbg==');

@$core.Deprecated('Use expiryEventsConfigurationDescriptor instead')
const ExpiryEventsConfiguration$json = {
  '1': 'ExpiryEventsConfiguration',
  '2': [
    {
      '1': 'daysbeforeexpiry',
      '3': 292642703,
      '4': 1,
      '5': 5,
      '10': 'daysbeforeexpiry'
    },
  ],
};

/// Descriptor for `ExpiryEventsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expiryEventsConfigurationDescriptor =
    $convert.base64Decode(
        'ChlFeHBpcnlFdmVudHNDb25maWd1cmF0aW9uEi4KEGRheXNiZWZvcmVleHBpcnkYj7/FiwEgAS'
        'gFUhBkYXlzYmVmb3JlZXhwaXJ5');

@$core.Deprecated('Use exportCertificateRequestDescriptor instead')
const ExportCertificateRequest$json = {
  '1': 'ExportCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {'1': 'passphrase', '3': 421146054, '4': 1, '5': 12, '10': 'passphrase'},
  ],
};

/// Descriptor for `ExportCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChhFeHBvcnRDZXJ0aWZpY2F0ZVJlcXVlc3QSKQoOY2VydGlmaWNhdGVhcm4Y+MqZLCABKAlSDm'
        'NlcnRpZmljYXRlYXJuEiIKCnBhc3NwaHJhc2UYxtvoyAEgASgMUgpwYXNzcGhyYXNl');

@$core.Deprecated('Use exportCertificateResponseDescriptor instead')
const ExportCertificateResponse$json = {
  '1': 'ExportCertificateResponse',
  '2': [
    {'1': 'certificate', '3': 198060817, '4': 1, '5': 9, '10': 'certificate'},
    {
      '1': 'certificatechain',
      '3': 369292378,
      '4': 1,
      '5': 9,
      '10': 'certificatechain'
    },
    {'1': 'privatekey', '3': 173312700, '4': 1, '5': 9, '10': 'privatekey'},
  ],
};

/// Descriptor for `ExportCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportCertificateResponseDescriptor = $convert.base64Decode(
    'ChlFeHBvcnRDZXJ0aWZpY2F0ZVJlc3BvbnNlEiMKC2NlcnRpZmljYXRlGJHWuF4gASgJUgtjZX'
    'J0aWZpY2F0ZRIuChBjZXJ0aWZpY2F0ZWNoYWluGNroi7ABIAEoCVIQY2VydGlmaWNhdGVjaGFp'
    'bhIhCgpwcml2YXRla2V5GLyV0lIgASgJUgpwcml2YXRla2V5');

@$core.Deprecated('Use extendedKeyUsageDescriptor instead')
const ExtendedKeyUsage$json = {
  '1': 'ExtendedKeyUsage',
  '2': [
    {
      '1': 'name',
      '3': 266367751,
      '4': 1,
      '5': 14,
      '6': '.acm.ExtendedKeyUsageName',
      '10': 'name'
    },
    {'1': 'oid', '3': 504527812, '4': 1, '5': 9, '10': 'oid'},
  ],
};

/// Descriptor for `ExtendedKeyUsage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List extendedKeyUsageDescriptor = $convert.base64Decode(
    'ChBFeHRlbmRlZEtleVVzYWdlEjAKBG5hbWUYh+aBfyABKA4yGS5hY20uRXh0ZW5kZWRLZXlVc2'
    'FnZU5hbWVSBG5hbWUSFAoDb2lkGMT3yfABIAEoCVIDb2lk');

@$core.Deprecated('Use filtersDescriptor instead')
const Filters$json = {
  '1': 'Filters',
  '2': [
    {
      '1': 'exportoption',
      '3': 442051119,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateExport',
      '10': 'exportoption'
    },
    {
      '1': 'extendedkeyusage',
      '3': 420785679,
      '4': 3,
      '5': 14,
      '6': '.acm.ExtendedKeyUsageName',
      '10': 'extendedkeyusage'
    },
    {
      '1': 'keytypes',
      '3': 27369686,
      '4': 3,
      '5': 14,
      '6': '.acm.KeyAlgorithm',
      '10': 'keytypes'
    },
    {
      '1': 'keyusage',
      '3': 502227108,
      '4': 3,
      '5': 14,
      '6': '.acm.KeyUsageName',
      '10': 'keyusage'
    },
    {
      '1': 'managedby',
      '3': 401567328,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateManagedBy',
      '10': 'managedby'
    },
  ],
};

/// Descriptor for `Filters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filtersDescriptor = $convert.base64Decode(
    'CgdGaWx0ZXJzEj4KDGV4cG9ydG9wdGlvbhiv1OTSASABKA4yFi5hY20uQ2VydGlmaWNhdGVFeH'
    'BvcnRSDGV4cG9ydG9wdGlvbhJJChBleHRlbmRlZGtleXVzYWdlGI/c0sgBIAMoDjIZLmFjbS5F'
    'eHRlbmRlZEtleVVzYWdlTmFtZVIQZXh0ZW5kZWRrZXl1c2FnZRIwCghrZXl0eXBlcxjWwYYNIA'
    'MoDjIRLmFjbS5LZXlBbGdvcml0aG1SCGtleXR5cGVzEjEKCGtleXVzYWdlGKTBve8BIAMoDjIR'
    'LmFjbS5LZXlVc2FnZU5hbWVSCGtleXVzYWdlEjsKCW1hbmFnZWRieRjg3L2/ASABKA4yGS5hY2'
    '0uQ2VydGlmaWNhdGVNYW5hZ2VkQnlSCW1hbmFnZWRieQ==');

@$core.Deprecated('Use getAccountConfigurationResponseDescriptor instead')
const GetAccountConfigurationResponse$json = {
  '1': 'GetAccountConfigurationResponse',
  '2': [
    {
      '1': 'expiryevents',
      '3': 358763356,
      '4': 1,
      '5': 11,
      '6': '.acm.ExpiryEventsConfiguration',
      '10': 'expiryevents'
    },
  ],
};

/// Descriptor for `GetAccountConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountConfigurationResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXRBY2NvdW50Q29uZmlndXJhdGlvblJlc3BvbnNlEkYKDGV4cGlyeWV2ZW50cxjclomrAS'
        'ABKAsyHi5hY20uRXhwaXJ5RXZlbnRzQ29uZmlndXJhdGlvblIMZXhwaXJ5ZXZlbnRz');

@$core.Deprecated('Use getCertificateRequestDescriptor instead')
const GetCertificateRequest$json = {
  '1': 'GetCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `GetCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCertificateRequestDescriptor = $convert.base64Decode(
    'ChVHZXRDZXJ0aWZpY2F0ZVJlcXVlc3QSKQoOY2VydGlmaWNhdGVhcm4Y+MqZLCABKAlSDmNlcn'
    'RpZmljYXRlYXJu');

@$core.Deprecated('Use getCertificateResponseDescriptor instead')
const GetCertificateResponse$json = {
  '1': 'GetCertificateResponse',
  '2': [
    {'1': 'certificate', '3': 198060817, '4': 1, '5': 9, '10': 'certificate'},
    {
      '1': 'certificatechain',
      '3': 369292378,
      '4': 1,
      '5': 9,
      '10': 'certificatechain'
    },
  ],
};

/// Descriptor for `GetCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCertificateResponseDescriptor = $convert.base64Decode(
    'ChZHZXRDZXJ0aWZpY2F0ZVJlc3BvbnNlEiMKC2NlcnRpZmljYXRlGJHWuF4gASgJUgtjZXJ0aW'
    'ZpY2F0ZRIuChBjZXJ0aWZpY2F0ZWNoYWluGNroi7ABIAEoCVIQY2VydGlmaWNhdGVjaGFpbg==');

@$core.Deprecated('Use httpRedirectDescriptor instead')
const HttpRedirect$json = {
  '1': 'HttpRedirect',
  '2': [
    {'1': 'redirectfrom', '3': 52853514, '4': 1, '5': 9, '10': 'redirectfrom'},
    {'1': 'redirectto', '3': 472243857, '4': 1, '5': 9, '10': 'redirectto'},
  ],
};

/// Descriptor for `HttpRedirect`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List httpRedirectDescriptor = $convert.base64Decode(
    'CgxIdHRwUmVkaXJlY3QSJQoMcmVkaXJlY3Rmcm9tGIr2mRkgASgJUgxyZWRpcmVjdGZyb20SIg'
    'oKcmVkaXJlY3R0bxiRvZfhASABKAlSCnJlZGlyZWN0dG8=');

@$core.Deprecated('Use importCertificateRequestDescriptor instead')
const ImportCertificateRequest$json = {
  '1': 'ImportCertificateRequest',
  '2': [
    {'1': 'certificate', '3': 198060817, '4': 1, '5': 12, '10': 'certificate'},
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'certificatechain',
      '3': 369292378,
      '4': 1,
      '5': 12,
      '10': 'certificatechain'
    },
    {'1': 'privatekey', '3': 173312700, '4': 1, '5': 12, '10': 'privatekey'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.acm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ImportCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importCertificateRequestDescriptor = $convert.base64Decode(
    'ChhJbXBvcnRDZXJ0aWZpY2F0ZVJlcXVlc3QSIwoLY2VydGlmaWNhdGUYkda4XiABKAxSC2Nlcn'
    'RpZmljYXRlEikKDmNlcnRpZmljYXRlYXJuGPjKmSwgASgJUg5jZXJ0aWZpY2F0ZWFybhIuChBj'
    'ZXJ0aWZpY2F0ZWNoYWluGNroi7ABIAEoDFIQY2VydGlmaWNhdGVjaGFpbhIhCgpwcml2YXRla2'
    'V5GLyV0lIgASgMUgpwcml2YXRla2V5EiAKBHRhZ3MYwcH2tQEgAygLMgguYWNtLlRhZ1IEdGFn'
    'cw==');

@$core.Deprecated('Use importCertificateResponseDescriptor instead')
const ImportCertificateResponse$json = {
  '1': 'ImportCertificateResponse',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `ImportCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importCertificateResponseDescriptor =
    $convert.base64Decode(
        'ChlJbXBvcnRDZXJ0aWZpY2F0ZVJlc3BvbnNlEikKDmNlcnRpZmljYXRlYXJuGPjKmSwgASgJUg'
        '5jZXJ0aWZpY2F0ZWFybg==');

@$core.Deprecated('Use invalidArgsExceptionDescriptor instead')
const InvalidArgsException$json = {
  '1': 'InvalidArgsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidArgsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidArgsExceptionDescriptor =
    $convert.base64Decode(
        'ChRJbnZhbGlkQXJnc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidArnExceptionDescriptor instead')
const InvalidArnException$json = {
  '1': 'InvalidArnException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidArnException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidArnExceptionDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkQXJuRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core
    .Deprecated('Use invalidDomainValidationOptionsExceptionDescriptor instead')
const InvalidDomainValidationOptionsException$json = {
  '1': 'InvalidDomainValidationOptionsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDomainValidationOptionsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDomainValidationOptionsExceptionDescriptor =
    $convert.base64Decode(
        'CidJbnZhbGlkRG9tYWluVmFsaWRhdGlvbk9wdGlvbnNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkc'
        'gnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidParameterExceptionDescriptor instead')
const InvalidParameterException$json = {
  '1': 'InvalidParameterException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkUGFyYW1ldGVyRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidStateExceptionDescriptor instead')
const InvalidStateException$json = {
  '1': 'InvalidStateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidStateExceptionDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkU3RhdGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidTagExceptionDescriptor instead')
const InvalidTagException$json = {
  '1': 'InvalidTagException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTagException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTagExceptionDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkVGFnRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use keyUsageDescriptor instead')
const KeyUsage$json = {
  '1': 'KeyUsage',
  '2': [
    {
      '1': 'name',
      '3': 266367751,
      '4': 1,
      '5': 14,
      '6': '.acm.KeyUsageName',
      '10': 'name'
    },
  ],
};

/// Descriptor for `KeyUsage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyUsageDescriptor = $convert.base64Decode(
    'CghLZXlVc2FnZRIoCgRuYW1lGIfmgX8gASgOMhEuYWNtLktleVVzYWdlTmFtZVIEbmFtZQ==');

@$core.Deprecated('Use limitExceededExceptionDescriptor instead')
const LimitExceededException$json = {
  '1': 'LimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChZMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use listCertificatesRequestDescriptor instead')
const ListCertificatesRequest$json = {
  '1': 'ListCertificatesRequest',
  '2': [
    {
      '1': 'certificatestatuses',
      '3': 203375987,
      '4': 3,
      '5': 14,
      '6': '.acm.CertificateStatus',
      '10': 'certificatestatuses'
    },
    {
      '1': 'includes',
      '3': 162404509,
      '4': 1,
      '5': 11,
      '6': '.acm.Filters',
      '10': 'includes'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'sortby',
      '3': 186052369,
      '4': 1,
      '5': 14,
      '6': '.acm.SortBy',
      '10': 'sortby'
    },
    {
      '1': 'sortorder',
      '3': 274231684,
      '4': 1,
      '5': 14,
      '6': '.acm.SortOrder',
      '10': 'sortorder'
    },
  ],
};

/// Descriptor for `ListCertificatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCertificatesRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0Q2VydGlmaWNhdGVzUmVxdWVzdBJLChNjZXJ0aWZpY2F0ZXN0YXR1c2VzGPOK/WAgAy'
    'gOMhYuYWNtLkNlcnRpZmljYXRlU3RhdHVzUhNjZXJ0aWZpY2F0ZXN0YXR1c2VzEisKCGluY2x1'
    'ZGVzGJ2xuE0gASgLMgwuYWNtLkZpbHRlcnNSCGluY2x1ZGVzEh4KCG1heGl0ZW1zGJTW2vEBIA'
    'EoBVIIbWF4aXRlbXMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SJgoGc29ydGJ5'
    'GJHe21ggASgOMgsuYWNtLlNvcnRCeVIGc29ydGJ5EjAKCXNvcnRvcmRlchiE4+GCASABKA4yDi'
    '5hY20uU29ydE9yZGVyUglzb3J0b3JkZXI=');

@$core.Deprecated('Use listCertificatesResponseDescriptor instead')
const ListCertificatesResponse$json = {
  '1': 'ListCertificatesResponse',
  '2': [
    {
      '1': 'certificatesummarylist',
      '3': 487936529,
      '4': 3,
      '5': 11,
      '6': '.acm.CertificateSummary',
      '10': 'certificatesummarylist'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCertificatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCertificatesResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0Q2VydGlmaWNhdGVzUmVzcG9uc2USUwoWY2VydGlmaWNhdGVzdW1tYXJ5bGlzdBiRpN'
    'XoASADKAsyFy5hY20uQ2VydGlmaWNhdGVTdW1tYXJ5UhZjZXJ0aWZpY2F0ZXN1bW1hcnlsaXN0'
    'Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listTagsForCertificateRequestDescriptor instead')
const ListTagsForCertificateRequest$json = {
  '1': 'ListTagsForCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `ListTagsForCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0VGFnc0ZvckNlcnRpZmljYXRlUmVxdWVzdBIpCg5jZXJ0aWZpY2F0ZWFybhj4ypksIA'
        'EoCVIOY2VydGlmaWNhdGVhcm4=');

@$core.Deprecated('Use listTagsForCertificateResponseDescriptor instead')
const ListTagsForCertificateResponse$json = {
  '1': 'ListTagsForCertificateResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.acm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForCertificateResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0VGFnc0ZvckNlcnRpZmljYXRlUmVzcG9uc2USIAoEdGFncxjBwfa1ASADKAsyCC5hY2'
        '0uVGFnUgR0YWdz');

@$core.Deprecated('Use putAccountConfigurationRequestDescriptor instead')
const PutAccountConfigurationRequest$json = {
  '1': 'PutAccountConfigurationRequest',
  '2': [
    {
      '1': 'expiryevents',
      '3': 358763356,
      '4': 1,
      '5': 11,
      '6': '.acm.ExpiryEventsConfiguration',
      '10': 'expiryevents'
    },
    {
      '1': 'idempotencytoken',
      '3': 56833648,
      '4': 1,
      '5': 9,
      '10': 'idempotencytoken'
    },
  ],
};

/// Descriptor for `PutAccountConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountConfigurationRequestDescriptor =
    $convert.base64Decode(
        'Ch5QdXRBY2NvdW50Q29uZmlndXJhdGlvblJlcXVlc3QSRgoMZXhwaXJ5ZXZlbnRzGNyWiasBIA'
        'EoCzIeLmFjbS5FeHBpcnlFdmVudHNDb25maWd1cmF0aW9uUgxleHBpcnlldmVudHMSLQoQaWRl'
        'bXBvdGVuY3l0b2tlbhjw7IwbIAEoCVIQaWRlbXBvdGVuY3l0b2tlbg==');

@$core.Deprecated('Use removeTagsFromCertificateRequestDescriptor instead')
const RemoveTagsFromCertificateRequest$json = {
  '1': 'RemoveTagsFromCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.acm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `RemoveTagsFromCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTagsFromCertificateRequestDescriptor =
    $convert.base64Decode(
        'CiBSZW1vdmVUYWdzRnJvbUNlcnRpZmljYXRlUmVxdWVzdBIpCg5jZXJ0aWZpY2F0ZWFybhj4yp'
        'ksIAEoCVIOY2VydGlmaWNhdGVhcm4SIAoEdGFncxjBwfa1ASADKAsyCC5hY20uVGFnUgR0YWdz');

@$core.Deprecated('Use renewCertificateRequestDescriptor instead')
const RenewCertificateRequest$json = {
  '1': 'RenewCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `RenewCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renewCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChdSZW5ld0NlcnRpZmljYXRlUmVxdWVzdBIpCg5jZXJ0aWZpY2F0ZWFybhj4ypksIAEoCVIOY2'
        'VydGlmaWNhdGVhcm4=');

@$core.Deprecated('Use renewalSummaryDescriptor instead')
const RenewalSummary$json = {
  '1': 'RenewalSummary',
  '2': [
    {
      '1': 'domainvalidationoptions',
      '3': 335573705,
      '4': 3,
      '5': 11,
      '6': '.acm.DomainValidation',
      '10': 'domainvalidationoptions'
    },
    {
      '1': 'renewalstatus',
      '3': 277232086,
      '4': 1,
      '5': 14,
      '6': '.acm.RenewalStatus',
      '10': 'renewalstatus'
    },
    {
      '1': 'renewalstatusreason',
      '3': 206499082,
      '4': 1,
      '5': 14,
      '6': '.acm.FailureReason',
      '10': 'renewalstatusreason'
    },
    {'1': 'updatedat', '3': 213581206, '4': 1, '5': 9, '10': 'updatedat'},
  ],
};

/// Descriptor for `RenewalSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renewalSummaryDescriptor = $convert.base64Decode(
    'Cg5SZW5ld2FsU3VtbWFyeRJTChdkb21haW52YWxpZGF0aW9ub3B0aW9ucxjJ5YGgASADKAsyFS'
    '5hY20uRG9tYWluVmFsaWRhdGlvblIXZG9tYWludmFsaWRhdGlvbm9wdGlvbnMSPAoNcmVuZXdh'
    'bHN0YXR1cxjW85iEASABKA4yEi5hY20uUmVuZXdhbFN0YXR1c1INcmVuZXdhbHN0YXR1cxJHCh'
    'NyZW5ld2Fsc3RhdHVzcmVhc29uGIrau2IgASgOMhIuYWNtLkZhaWx1cmVSZWFzb25SE3JlbmV3'
    'YWxzdGF0dXNyZWFzb24SHwoJdXBkYXRlZGF0GJb762UgASgJUgl1cGRhdGVkYXQ=');

@$core.Deprecated('Use requestCertificateRequestDescriptor instead')
const RequestCertificateRequest$json = {
  '1': 'RequestCertificateRequest',
  '2': [
    {
      '1': 'certificateauthorityarn',
      '3': 266069181,
      '4': 1,
      '5': 9,
      '10': 'certificateauthorityarn'
    },
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'domainvalidationoptions',
      '3': 335573705,
      '4': 3,
      '5': 11,
      '6': '.acm.DomainValidationOption',
      '10': 'domainvalidationoptions'
    },
    {
      '1': 'idempotencytoken',
      '3': 56833648,
      '4': 1,
      '5': 9,
      '10': 'idempotencytoken'
    },
    {
      '1': 'keyalgorithm',
      '3': 452317818,
      '4': 1,
      '5': 14,
      '6': '.acm.KeyAlgorithm',
      '10': 'keyalgorithm'
    },
    {
      '1': 'managedby',
      '3': 455511232,
      '4': 1,
      '5': 14,
      '6': '.acm.CertificateManagedBy',
      '10': 'managedby'
    },
    {
      '1': 'options',
      '3': 356388166,
      '4': 1,
      '5': 11,
      '6': '.acm.CertificateOptions',
      '10': 'options'
    },
    {
      '1': 'subjectalternativenames',
      '3': 109998119,
      '4': 3,
      '5': 9,
      '10': 'subjectalternativenames'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.acm.Tag',
      '10': 'tags'
    },
    {
      '1': 'validationmethod',
      '3': 58092520,
      '4': 1,
      '5': 14,
      '6': '.acm.ValidationMethod',
      '10': 'validationmethod'
    },
  ],
};

/// Descriptor for `RequestCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestCertificateRequestDescriptor = $convert.base64Decode(
    'ChlSZXF1ZXN0Q2VydGlmaWNhdGVSZXF1ZXN0EjsKF2NlcnRpZmljYXRlYXV0aG9yaXR5YXJuGL'
    '3J734gASgJUhdjZXJ0aWZpY2F0ZWF1dGhvcml0eWFybhIhCgpkb21haW5uYW1lGOvN+FwgASgJ'
    'Ugpkb21haW5uYW1lElkKF2RvbWFpbnZhbGlkYXRpb25vcHRpb25zGMnlgaABIAMoCzIbLmFjbS'
    '5Eb21haW5WYWxpZGF0aW9uT3B0aW9uUhdkb21haW52YWxpZGF0aW9ub3B0aW9ucxItChBpZGVt'
    'cG90ZW5jeXRva2VuGPDsjBsgASgJUhBpZGVtcG90ZW5jeXRva2VuEjkKDGtleWFsZ29yaXRobR'
    'j6pNfXASABKA4yES5hY20uS2V5QWxnb3JpdGhtUgxrZXlhbGdvcml0aG0SOwoJbWFuYWdlZGJ5'
    'GMCZmtkBIAEoDjIZLmFjbS5DZXJ0aWZpY2F0ZU1hbmFnZWRCeVIJbWFuYWdlZGJ5EjUKB29wdG'
    'lvbnMYxpr4qQEgASgLMhcuYWNtLkNlcnRpZmljYXRlT3B0aW9uc1IHb3B0aW9ucxI7ChdzdWJq'
    'ZWN0YWx0ZXJuYXRpdmVuYW1lcxin4Lk0IAMoCVIXc3ViamVjdGFsdGVybmF0aXZlbmFtZXMSIA'
    'oEdGFncxjBwfa1ASADKAsyCC5hY20uVGFnUgR0YWdzEkQKEHZhbGlkYXRpb25tZXRob2QY6NfZ'
    'GyABKA4yFS5hY20uVmFsaWRhdGlvbk1ldGhvZFIQdmFsaWRhdGlvbm1ldGhvZA==');

@$core.Deprecated('Use requestCertificateResponseDescriptor instead')
const RequestCertificateResponse$json = {
  '1': 'RequestCertificateResponse',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `RequestCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestCertificateResponseDescriptor =
    $convert.base64Decode(
        'ChpSZXF1ZXN0Q2VydGlmaWNhdGVSZXNwb25zZRIpCg5jZXJ0aWZpY2F0ZWFybhj4ypksIAEoCV'
        'IOY2VydGlmaWNhdGVhcm4=');

@$core.Deprecated('Use requestInProgressExceptionDescriptor instead')
const RequestInProgressException$json = {
  '1': 'RequestInProgressException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RequestInProgressException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestInProgressExceptionDescriptor =
    $convert.base64Decode(
        'ChpSZXF1ZXN0SW5Qcm9ncmVzc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use resendValidationEmailRequestDescriptor instead')
const ResendValidationEmailRequest$json = {
  '1': 'ResendValidationEmailRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'validationdomain',
      '3': 161300799,
      '4': 1,
      '5': 9,
      '10': 'validationdomain'
    },
  ],
};

/// Descriptor for `ResendValidationEmailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resendValidationEmailRequestDescriptor =
    $convert.base64Decode(
        'ChxSZXNlbmRWYWxpZGF0aW9uRW1haWxSZXF1ZXN0EikKDmNlcnRpZmljYXRlYXJuGPjKmSwgAS'
        'gJUg5jZXJ0aWZpY2F0ZWFybhIaCgZkb21haW4YkpLy8AEgASgJUgZkb21haW4SLQoQdmFsaWRh'
        'dGlvbmRvbWFpbhi/gvVMIAEoCVIQdmFsaWRhdGlvbmRvbWFpbg==');

@$core.Deprecated('Use resourceInUseExceptionDescriptor instead')
const ResourceInUseException$json = {
  '1': 'ResourceInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChZSZXNvdXJjZUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use resourceRecordDescriptor instead')
const ResourceRecord$json = {
  '1': 'ResourceRecord',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.acm.RecordType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ResourceRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceRecordDescriptor = $convert.base64Decode(
    'Cg5SZXNvdXJjZVJlY29yZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEicKBHR5cGUY7qDXigEgAS'
    'gOMg8uYWNtLlJlY29yZFR5cGVSBHR5cGUSGAoFdmFsdWUY6/KfigEgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use revokeCertificateRequestDescriptor instead')
const RevokeCertificateRequest$json = {
  '1': 'RevokeCertificateRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'revocationreason',
      '3': 331836358,
      '4': 1,
      '5': 14,
      '6': '.acm.RevocationReason',
      '10': 'revocationreason'
    },
  ],
};

/// Descriptor for `RevokeCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeCertificateRequestDescriptor = $convert.base64Decode(
    'ChhSZXZva2VDZXJ0aWZpY2F0ZVJlcXVlc3QSKQoOY2VydGlmaWNhdGVhcm4Y+MqZLCABKAlSDm'
    'NlcnRpZmljYXRlYXJuEkUKEHJldm9jYXRpb25yZWFzb24YxtedngEgASgOMhUuYWNtLlJldm9j'
    'YXRpb25SZWFzb25SEHJldm9jYXRpb25yZWFzb24=');

@$core.Deprecated('Use revokeCertificateResponseDescriptor instead')
const RevokeCertificateResponse$json = {
  '1': 'RevokeCertificateResponse',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
  ],
};

/// Descriptor for `RevokeCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeCertificateResponseDescriptor =
    $convert.base64Decode(
        'ChlSZXZva2VDZXJ0aWZpY2F0ZVJlc3BvbnNlEikKDmNlcnRpZmljYXRlYXJuGPjKmSwgASgJUg'
        '5jZXJ0aWZpY2F0ZWFybg==');

@$core.Deprecated('Use tagDescriptor instead')
const Tag$json = {
  '1': 'Tag',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Tag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagDescriptor = $convert.base64Decode(
    'CgNUYWcSEwoDa2V5GI2S62ggASgJUgNrZXkSGAoFdmFsdWUY6/KfigEgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use tagPolicyExceptionDescriptor instead')
const TagPolicyException$json = {
  '1': 'TagPolicyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TagPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagPolicyExceptionDescriptor =
    $convert.base64Decode(
        'ChJUYWdQb2xpY3lFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use throttlingExceptionDescriptor instead')
const ThrottlingException$json = {
  '1': 'ThrottlingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttlingExceptionDescriptor =
    $convert.base64Decode(
        'ChNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyTagsExceptionDescriptor instead')
const TooManyTagsException$json = {
  '1': 'TooManyTagsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyTagsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyTagsExceptionDescriptor =
    $convert.base64Decode(
        'ChRUb29NYW55VGFnc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use updateCertificateOptionsRequestDescriptor instead')
const UpdateCertificateOptionsRequest$json = {
  '1': 'UpdateCertificateOptionsRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'options',
      '3': 356388166,
      '4': 1,
      '5': 11,
      '6': '.acm.CertificateOptions',
      '10': 'options'
    },
  ],
};

/// Descriptor for `UpdateCertificateOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCertificateOptionsRequestDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVDZXJ0aWZpY2F0ZU9wdGlvbnNSZXF1ZXN0EikKDmNlcnRpZmljYXRlYXJuGPjKmS'
        'wgASgJUg5jZXJ0aWZpY2F0ZWFybhI1CgdvcHRpb25zGMaa+KkBIAEoCzIXLmFjbS5DZXJ0aWZp'
        'Y2F0ZU9wdGlvbnNSB29wdGlvbnM=');

@$core.Deprecated('Use validationExceptionDescriptor instead')
const ValidationException$json = {
  '1': 'ValidationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationExceptionDescriptor =
    $convert.base64Decode(
        'ChNWYWxpZGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');
