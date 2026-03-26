// This is a generated file - do not edit.
//
// Generated from cloudfront.proto.

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

@$core.Deprecated('Use cachePolicyCookieBehaviorDescriptor instead')
const CachePolicyCookieBehavior$json = {
  '1': 'CachePolicyCookieBehavior',
  '2': [
    {'1': 'CACHE_POLICY_COOKIE_BEHAVIOR_ALL', '2': 0},
    {'1': 'CACHE_POLICY_COOKIE_BEHAVIOR_WHITELIST', '2': 1},
    {'1': 'CACHE_POLICY_COOKIE_BEHAVIOR_NONE', '2': 2},
    {'1': 'CACHE_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT', '2': 3},
  ],
};

/// Descriptor for `CachePolicyCookieBehavior`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cachePolicyCookieBehaviorDescriptor = $convert.base64Decode(
    'ChlDYWNoZVBvbGljeUNvb2tpZUJlaGF2aW9yEiQKIENBQ0hFX1BPTElDWV9DT09LSUVfQkVIQV'
    'ZJT1JfQUxMEAASKgomQ0FDSEVfUE9MSUNZX0NPT0tJRV9CRUhBVklPUl9XSElURUxJU1QQARIl'
    'CiFDQUNIRV9QT0xJQ1lfQ09PS0lFX0JFSEFWSU9SX05PTkUQAhIqCiZDQUNIRV9QT0xJQ1lfQ0'
    '9PS0lFX0JFSEFWSU9SX0FMTEVYQ0VQVBAD');

@$core.Deprecated('Use cachePolicyHeaderBehaviorDescriptor instead')
const CachePolicyHeaderBehavior$json = {
  '1': 'CachePolicyHeaderBehavior',
  '2': [
    {'1': 'CACHE_POLICY_HEADER_BEHAVIOR_WHITELIST', '2': 0},
    {'1': 'CACHE_POLICY_HEADER_BEHAVIOR_NONE', '2': 1},
  ],
};

/// Descriptor for `CachePolicyHeaderBehavior`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cachePolicyHeaderBehaviorDescriptor = $convert.base64Decode(
    'ChlDYWNoZVBvbGljeUhlYWRlckJlaGF2aW9yEioKJkNBQ0hFX1BPTElDWV9IRUFERVJfQkVIQV'
    'ZJT1JfV0hJVEVMSVNUEAASJQohQ0FDSEVfUE9MSUNZX0hFQURFUl9CRUhBVklPUl9OT05FEAE=');

@$core.Deprecated('Use cachePolicyQueryStringBehaviorDescriptor instead')
const CachePolicyQueryStringBehavior$json = {
  '1': 'CachePolicyQueryStringBehavior',
  '2': [
    {'1': 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALL', '2': 0},
    {'1': 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST', '2': 1},
    {'1': 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_NONE', '2': 2},
    {'1': 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT', '2': 3},
  ],
};

/// Descriptor for `CachePolicyQueryStringBehavior`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cachePolicyQueryStringBehaviorDescriptor = $convert.base64Decode(
    'Ch5DYWNoZVBvbGljeVF1ZXJ5U3RyaW5nQmVoYXZpb3ISKgomQ0FDSEVfUE9MSUNZX1FVRVJZX1'
    'NUUklOR19CRUhBVklPUl9BTEwQABIwCixDQUNIRV9QT0xJQ1lfUVVFUllfU1RSSU5HX0JFSEFW'
    'SU9SX1dISVRFTElTVBABEisKJ0NBQ0hFX1BPTElDWV9RVUVSWV9TVFJJTkdfQkVIQVZJT1JfTk'
    '9ORRACEjAKLENBQ0hFX1BPTElDWV9RVUVSWV9TVFJJTkdfQkVIQVZJT1JfQUxMRVhDRVBUEAM=');

@$core.Deprecated('Use cachePolicyTypeDescriptor instead')
const CachePolicyType$json = {
  '1': 'CachePolicyType',
  '2': [
    {'1': 'CACHE_POLICY_TYPE_MANAGED', '2': 0},
    {'1': 'CACHE_POLICY_TYPE_CUSTOM', '2': 1},
  ],
};

/// Descriptor for `CachePolicyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cachePolicyTypeDescriptor = $convert.base64Decode(
    'Cg9DYWNoZVBvbGljeVR5cGUSHQoZQ0FDSEVfUE9MSUNZX1RZUEVfTUFOQUdFRBAAEhwKGENBQ0'
    'hFX1BPTElDWV9UWVBFX0NVU1RPTRAB');

@$core.Deprecated('Use certificateSourceDescriptor instead')
const CertificateSource$json = {
  '1': 'CertificateSource',
  '2': [
    {'1': 'CERTIFICATE_SOURCE_IAM', '2': 0},
    {'1': 'CERTIFICATE_SOURCE_ACM', '2': 1},
    {'1': 'CERTIFICATE_SOURCE_CLOUDFRONT', '2': 2},
  ],
};

/// Descriptor for `CertificateSource`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List certificateSourceDescriptor = $convert.base64Decode(
    'ChFDZXJ0aWZpY2F0ZVNvdXJjZRIaChZDRVJUSUZJQ0FURV9TT1VSQ0VfSUFNEAASGgoWQ0VSVE'
    'lGSUNBVEVfU09VUkNFX0FDTRABEiEKHUNFUlRJRklDQVRFX1NPVVJDRV9DTE9VREZST05UEAI=');

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

@$core.Deprecated('Use connectionModeDescriptor instead')
const ConnectionMode$json = {
  '1': 'ConnectionMode',
  '2': [
    {'1': 'CONNECTION_MODE_DIRECT', '2': 0},
    {'1': 'CONNECTION_MODE_TENANTONLY', '2': 1},
  ],
};

/// Descriptor for `ConnectionMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionModeDescriptor = $convert.base64Decode(
    'Cg5Db25uZWN0aW9uTW9kZRIaChZDT05ORUNUSU9OX01PREVfRElSRUNUEAASHgoaQ09OTkVDVE'
    'lPTl9NT0RFX1RFTkFOVE9OTFkQAQ==');

@$core.Deprecated('Use continuousDeploymentPolicyTypeDescriptor instead')
const ContinuousDeploymentPolicyType$json = {
  '1': 'ContinuousDeploymentPolicyType',
  '2': [
    {'1': 'CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEHEADER', '2': 0},
    {'1': 'CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEWEIGHT', '2': 1},
  ],
};

/// Descriptor for `ContinuousDeploymentPolicyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List continuousDeploymentPolicyTypeDescriptor =
    $convert.base64Decode(
        'Ch5Db250aW51b3VzRGVwbG95bWVudFBvbGljeVR5cGUSMgouQ09OVElOVU9VU19ERVBMT1lNRU'
        '5UX1BPTElDWV9UWVBFX1NJTkdMRUhFQURFUhAAEjIKLkNPTlRJTlVPVVNfREVQTE9ZTUVOVF9Q'
        'T0xJQ1lfVFlQRV9TSU5HTEVXRUlHSFQQAQ==');

@$core.Deprecated('Use customizationActionTypeDescriptor instead')
const CustomizationActionType$json = {
  '1': 'CustomizationActionType',
  '2': [
    {'1': 'CUSTOMIZATION_ACTION_TYPE_OVERRIDE', '2': 0},
    {'1': 'CUSTOMIZATION_ACTION_TYPE_DISABLE', '2': 1},
  ],
};

/// Descriptor for `CustomizationActionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List customizationActionTypeDescriptor =
    $convert.base64Decode(
        'ChdDdXN0b21pemF0aW9uQWN0aW9uVHlwZRImCiJDVVNUT01JWkFUSU9OX0FDVElPTl9UWVBFX0'
        '9WRVJSSURFEAASJQohQ1VTVE9NSVpBVElPTl9BQ1RJT05fVFlQRV9ESVNBQkxFEAE=');

@$core.Deprecated('Use distributionResourceTypeDescriptor instead')
const DistributionResourceType$json = {
  '1': 'DistributionResourceType',
  '2': [
    {'1': 'DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTION', '2': 0},
    {'1': 'DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTIONTENANT', '2': 1},
  ],
};

/// Descriptor for `DistributionResourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List distributionResourceTypeDescriptor = $convert.base64Decode(
    'ChhEaXN0cmlidXRpb25SZXNvdXJjZVR5cGUSKwonRElTVFJJQlVUSU9OX1JFU09VUkNFX1RZUE'
    'VfRElTVFJJQlVUSU9OEAASMQotRElTVFJJQlVUSU9OX1JFU09VUkNFX1RZUEVfRElTVFJJQlVU'
    'SU9OVEVOQU5UEAE=');

@$core.Deprecated('Use dnsConfigurationStatusDescriptor instead')
const DnsConfigurationStatus$json = {
  '1': 'DnsConfigurationStatus',
  '2': [
    {'1': 'DNS_CONFIGURATION_STATUS_VALID', '2': 0},
    {'1': 'DNS_CONFIGURATION_STATUS_INVALID', '2': 1},
    {'1': 'DNS_CONFIGURATION_STATUS_UNKNOWN', '2': 2},
  ],
};

/// Descriptor for `DnsConfigurationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dnsConfigurationStatusDescriptor = $convert.base64Decode(
    'ChZEbnNDb25maWd1cmF0aW9uU3RhdHVzEiIKHkROU19DT05GSUdVUkFUSU9OX1NUQVRVU19WQU'
    'xJRBAAEiQKIEROU19DT05GSUdVUkFUSU9OX1NUQVRVU19JTlZBTElEEAESJAogRE5TX0NPTkZJ'
    'R1VSQVRJT05fU1RBVFVTX1VOS05PV04QAg==');

@$core.Deprecated('Use domainStatusDescriptor instead')
const DomainStatus$json = {
  '1': 'DomainStatus',
  '2': [
    {'1': 'DOMAIN_STATUS_ACTIVE', '2': 0},
    {'1': 'DOMAIN_STATUS_INACTIVE', '2': 1},
  ],
};

/// Descriptor for `DomainStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List domainStatusDescriptor = $convert.base64Decode(
    'CgxEb21haW5TdGF0dXMSGAoURE9NQUlOX1NUQVRVU19BQ1RJVkUQABIaChZET01BSU5fU1RBVF'
    'VTX0lOQUNUSVZFEAE=');

@$core.Deprecated('Use eventTypeDescriptor instead')
const EventType$json = {
  '1': 'EventType',
  '2': [
    {'1': 'EVENT_TYPE_VIEWER_REQUEST', '2': 0},
    {'1': 'EVENT_TYPE_ORIGIN_REQUEST', '2': 1},
    {'1': 'EVENT_TYPE_VIEWER_RESPONSE', '2': 2},
    {'1': 'EVENT_TYPE_ORIGIN_RESPONSE', '2': 3},
  ],
};

/// Descriptor for `EventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventTypeDescriptor = $convert.base64Decode(
    'CglFdmVudFR5cGUSHQoZRVZFTlRfVFlQRV9WSUVXRVJfUkVRVUVTVBAAEh0KGUVWRU5UX1RZUE'
    'VfT1JJR0lOX1JFUVVFU1QQARIeChpFVkVOVF9UWVBFX1ZJRVdFUl9SRVNQT05TRRACEh4KGkVW'
    'RU5UX1RZUEVfT1JJR0lOX1JFU1BPTlNFEAM=');

@$core.Deprecated('Use formatDescriptor instead')
const Format$json = {
  '1': 'Format',
  '2': [
    {'1': 'FORMAT_URLENCODED', '2': 0},
  ],
};

/// Descriptor for `Format`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List formatDescriptor =
    $convert.base64Decode('CgZGb3JtYXQSFQoRRk9STUFUX1VSTEVOQ09ERUQQAA==');

@$core.Deprecated('Use frameOptionsListDescriptor instead')
const FrameOptionsList$json = {
  '1': 'FrameOptionsList',
  '2': [
    {'1': 'FRAME_OPTIONS_LIST_SAMEORIGIN', '2': 0},
    {'1': 'FRAME_OPTIONS_LIST_DENY', '2': 1},
  ],
};

/// Descriptor for `FrameOptionsList`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List frameOptionsListDescriptor = $convert.base64Decode(
    'ChBGcmFtZU9wdGlvbnNMaXN0EiEKHUZSQU1FX09QVElPTlNfTElTVF9TQU1FT1JJR0lOEAASGw'
    'oXRlJBTUVfT1BUSU9OU19MSVNUX0RFTlkQAQ==');

@$core.Deprecated('Use functionRuntimeDescriptor instead')
const FunctionRuntime$json = {
  '1': 'FunctionRuntime',
  '2': [
    {'1': 'FUNCTION_RUNTIME_CLOUDFRONT_JS_1_0', '2': 0},
    {'1': 'FUNCTION_RUNTIME_CLOUDFRONT_JS_2_0', '2': 1},
  ],
};

/// Descriptor for `FunctionRuntime`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List functionRuntimeDescriptor = $convert.base64Decode(
    'Cg9GdW5jdGlvblJ1bnRpbWUSJgoiRlVOQ1RJT05fUlVOVElNRV9DTE9VREZST05UX0pTXzFfMB'
    'AAEiYKIkZVTkNUSU9OX1JVTlRJTUVfQ0xPVURGUk9OVF9KU18yXzAQAQ==');

@$core.Deprecated('Use functionStageDescriptor instead')
const FunctionStage$json = {
  '1': 'FunctionStage',
  '2': [
    {'1': 'FUNCTION_STAGE_DEVELOPMENT', '2': 0},
    {'1': 'FUNCTION_STAGE_LIVE', '2': 1},
  ],
};

/// Descriptor for `FunctionStage`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List functionStageDescriptor = $convert.base64Decode(
    'Cg1GdW5jdGlvblN0YWdlEh4KGkZVTkNUSU9OX1NUQUdFX0RFVkVMT1BNRU5UEAASFwoTRlVOQ1'
    'RJT05fU1RBR0VfTElWRRAB');

@$core.Deprecated('Use geoRestrictionTypeDescriptor instead')
const GeoRestrictionType$json = {
  '1': 'GeoRestrictionType',
  '2': [
    {'1': 'GEO_RESTRICTION_TYPE_BLACKLIST', '2': 0},
    {'1': 'GEO_RESTRICTION_TYPE_WHITELIST', '2': 1},
    {'1': 'GEO_RESTRICTION_TYPE_NONE', '2': 2},
  ],
};

/// Descriptor for `GeoRestrictionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List geoRestrictionTypeDescriptor = $convert.base64Decode(
    'ChJHZW9SZXN0cmljdGlvblR5cGUSIgoeR0VPX1JFU1RSSUNUSU9OX1RZUEVfQkxBQ0tMSVNUEA'
    'ASIgoeR0VPX1JFU1RSSUNUSU9OX1RZUEVfV0hJVEVMSVNUEAESHQoZR0VPX1JFU1RSSUNUSU9O'
    'X1RZUEVfTk9ORRAC');

@$core.Deprecated('Use httpVersionDescriptor instead')
const HttpVersion$json = {
  '1': 'HttpVersion',
  '2': [
    {'1': 'HTTP_VERSION_HTTP2', '2': 0},
    {'1': 'HTTP_VERSION_HTTP3', '2': 1},
    {'1': 'HTTP_VERSION_HTTP1_1', '2': 2},
    {'1': 'HTTP_VERSION_HTTP2AND3', '2': 3},
  ],
};

/// Descriptor for `HttpVersion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List httpVersionDescriptor = $convert.base64Decode(
    'CgtIdHRwVmVyc2lvbhIWChJIVFRQX1ZFUlNJT05fSFRUUDIQABIWChJIVFRQX1ZFUlNJT05fSF'
    'RUUDMQARIYChRIVFRQX1ZFUlNJT05fSFRUUDFfMRACEhoKFkhUVFBfVkVSU0lPTl9IVFRQMkFO'
    'RDMQAw==');

@$core.Deprecated('Use iCPRecordalStatusDescriptor instead')
const ICPRecordalStatus$json = {
  '1': 'ICPRecordalStatus',
  '2': [
    {'1': 'I_C_P_RECORDAL_STATUS_PENDING', '2': 0},
    {'1': 'I_C_P_RECORDAL_STATUS_SUSPENDED', '2': 1},
    {'1': 'I_C_P_RECORDAL_STATUS_APPROVED', '2': 2},
  ],
};

/// Descriptor for `ICPRecordalStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List iCPRecordalStatusDescriptor = $convert.base64Decode(
    'ChFJQ1BSZWNvcmRhbFN0YXR1cxIhCh1JX0NfUF9SRUNPUkRBTF9TVEFUVVNfUEVORElORxAAEi'
    'MKH0lfQ19QX1JFQ09SREFMX1NUQVRVU19TVVNQRU5ERUQQARIiCh5JX0NfUF9SRUNPUkRBTF9T'
    'VEFUVVNfQVBQUk9WRUQQAg==');

@$core.Deprecated('Use importSourceTypeDescriptor instead')
const ImportSourceType$json = {
  '1': 'ImportSourceType',
  '2': [
    {'1': 'IMPORT_SOURCE_TYPE_S3', '2': 0},
  ],
};

/// Descriptor for `ImportSourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importSourceTypeDescriptor = $convert.base64Decode(
    'ChBJbXBvcnRTb3VyY2VUeXBlEhkKFUlNUE9SVF9TT1VSQ0VfVFlQRV9TMxAA');

@$core.Deprecated('Use ipAddressTypeDescriptor instead')
const IpAddressType$json = {
  '1': 'IpAddressType',
  '2': [
    {'1': 'IP_ADDRESS_TYPE_IPV6', '2': 0},
    {'1': 'IP_ADDRESS_TYPE_DUALSTACK', '2': 1},
    {'1': 'IP_ADDRESS_TYPE_IPV4', '2': 2},
  ],
};

/// Descriptor for `IpAddressType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List ipAddressTypeDescriptor = $convert.base64Decode(
    'Cg1JcEFkZHJlc3NUeXBlEhgKFElQX0FERFJFU1NfVFlQRV9JUFY2EAASHQoZSVBfQUREUkVTU1'
    '9UWVBFX0RVQUxTVEFDSxABEhgKFElQX0FERFJFU1NfVFlQRV9JUFY0EAI=');

@$core.Deprecated('Use ipamCidrStatusDescriptor instead')
const IpamCidrStatus$json = {
  '1': 'IpamCidrStatus',
  '2': [
    {'1': 'IPAM_CIDR_STATUS_PROVISIONED', '2': 0},
    {'1': 'IPAM_CIDR_STATUS_FAILEDDEPROVISION', '2': 1},
    {'1': 'IPAM_CIDR_STATUS_FAILEDPROVISION', '2': 2},
    {'1': 'IPAM_CIDR_STATUS_ADVERTISED', '2': 3},
    {'1': 'IPAM_CIDR_STATUS_DEPROVISIONED', '2': 4},
    {'1': 'IPAM_CIDR_STATUS_WITHDRAWN', '2': 5},
    {'1': 'IPAM_CIDR_STATUS_PROVISIONING', '2': 6},
    {'1': 'IPAM_CIDR_STATUS_FAILEDADVERTISE', '2': 7},
    {'1': 'IPAM_CIDR_STATUS_FAILEDWITHDRAW', '2': 8},
    {'1': 'IPAM_CIDR_STATUS_DEPROVISIONING', '2': 9},
    {'1': 'IPAM_CIDR_STATUS_ADVERTISING', '2': 10},
    {'1': 'IPAM_CIDR_STATUS_WITHDRAWING', '2': 11},
  ],
};

/// Descriptor for `IpamCidrStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List ipamCidrStatusDescriptor = $convert.base64Decode(
    'Cg5JcGFtQ2lkclN0YXR1cxIgChxJUEFNX0NJRFJfU1RBVFVTX1BST1ZJU0lPTkVEEAASJgoiSV'
    'BBTV9DSURSX1NUQVRVU19GQUlMRURERVBST1ZJU0lPThABEiQKIElQQU1fQ0lEUl9TVEFUVVNf'
    'RkFJTEVEUFJPVklTSU9OEAISHwobSVBBTV9DSURSX1NUQVRVU19BRFZFUlRJU0VEEAMSIgoeSV'
    'BBTV9DSURSX1NUQVRVU19ERVBST1ZJU0lPTkVEEAQSHgoaSVBBTV9DSURSX1NUQVRVU19XSVRI'
    'RFJBV04QBRIhCh1JUEFNX0NJRFJfU1RBVFVTX1BST1ZJU0lPTklORxAGEiQKIElQQU1fQ0lEUl'
    '9TVEFUVVNfRkFJTEVEQURWRVJUSVNFEAcSIwofSVBBTV9DSURSX1NUQVRVU19GQUlMRURXSVRI'
    'RFJBVxAIEiMKH0lQQU1fQ0lEUl9TVEFUVVNfREVQUk9WSVNJT05JTkcQCRIgChxJUEFNX0NJRF'
    'JfU1RBVFVTX0FEVkVSVElTSU5HEAoSIAocSVBBTV9DSURSX1NUQVRVU19XSVRIRFJBV0lORxAL');

@$core.Deprecated('Use itemSelectionDescriptor instead')
const ItemSelection$json = {
  '1': 'ItemSelection',
  '2': [
    {'1': 'ITEM_SELECTION_ALL', '2': 0},
    {'1': 'ITEM_SELECTION_WHITELIST', '2': 1},
    {'1': 'ITEM_SELECTION_NONE', '2': 2},
  ],
};

/// Descriptor for `ItemSelection`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List itemSelectionDescriptor = $convert.base64Decode(
    'Cg1JdGVtU2VsZWN0aW9uEhYKEklURU1fU0VMRUNUSU9OX0FMTBAAEhwKGElURU1fU0VMRUNUSU'
    '9OX1dISVRFTElTVBABEhcKE0lURU1fU0VMRUNUSU9OX05PTkUQAg==');

@$core.Deprecated('Use managedCertificateStatusDescriptor instead')
const ManagedCertificateStatus$json = {
  '1': 'ManagedCertificateStatus',
  '2': [
    {'1': 'MANAGED_CERTIFICATE_STATUS_PENDINGVALIDATION', '2': 0},
    {'1': 'MANAGED_CERTIFICATE_STATUS_REVOKED', '2': 1},
    {'1': 'MANAGED_CERTIFICATE_STATUS_VALIDATIONTIMEDOUT', '2': 2},
    {'1': 'MANAGED_CERTIFICATE_STATUS_FAILED', '2': 3},
    {'1': 'MANAGED_CERTIFICATE_STATUS_EXPIRED', '2': 4},
    {'1': 'MANAGED_CERTIFICATE_STATUS_ISSUED', '2': 5},
    {'1': 'MANAGED_CERTIFICATE_STATUS_INACTIVE', '2': 6},
  ],
};

/// Descriptor for `ManagedCertificateStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List managedCertificateStatusDescriptor = $convert.base64Decode(
    'ChhNYW5hZ2VkQ2VydGlmaWNhdGVTdGF0dXMSMAosTUFOQUdFRF9DRVJUSUZJQ0FURV9TVEFUVV'
    'NfUEVORElOR1ZBTElEQVRJT04QABImCiJNQU5BR0VEX0NFUlRJRklDQVRFX1NUQVRVU19SRVZP'
    'S0VEEAESMQotTUFOQUdFRF9DRVJUSUZJQ0FURV9TVEFUVVNfVkFMSURBVElPTlRJTUVET1VUEA'
    'ISJQohTUFOQUdFRF9DRVJUSUZJQ0FURV9TVEFUVVNfRkFJTEVEEAMSJgoiTUFOQUdFRF9DRVJU'
    'SUZJQ0FURV9TVEFUVVNfRVhQSVJFRBAEEiUKIU1BTkFHRURfQ0VSVElGSUNBVEVfU1RBVFVTX0'
    'lTU1VFRBAFEicKI01BTkFHRURfQ0VSVElGSUNBVEVfU1RBVFVTX0lOQUNUSVZFEAY=');

@$core.Deprecated('Use methodDescriptor instead')
const Method$json = {
  '1': 'Method',
  '2': [
    {'1': 'METHOD_OPTIONS', '2': 0},
    {'1': 'METHOD_PATCH', '2': 1},
    {'1': 'METHOD_POST', '2': 2},
    {'1': 'METHOD_HEAD', '2': 3},
    {'1': 'METHOD_GET', '2': 4},
    {'1': 'METHOD_DELETE', '2': 5},
    {'1': 'METHOD_PUT', '2': 6},
  ],
};

/// Descriptor for `Method`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List methodDescriptor = $convert.base64Decode(
    'CgZNZXRob2QSEgoOTUVUSE9EX09QVElPTlMQABIQCgxNRVRIT0RfUEFUQ0gQARIPCgtNRVRIT0'
    'RfUE9TVBACEg8KC01FVEhPRF9IRUFEEAMSDgoKTUVUSE9EX0dFVBAEEhEKDU1FVEhPRF9ERUxF'
    'VEUQBRIOCgpNRVRIT0RfUFVUEAY=');

@$core.Deprecated('Use minimumProtocolVersionDescriptor instead')
const MinimumProtocolVersion$json = {
  '1': 'MinimumProtocolVersion',
  '2': [
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1_1_2016', '2': 0},
    {'1': 'MINIMUM_PROTOCOL_VERSION_SSLV3', '2': 1},
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1_2016', '2': 2},
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2019', '2': 3},
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1_3_2025', '2': 4},
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2021', '2': 5},
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1', '2': 6},
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2018', '2': 7},
    {'1': 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2025', '2': 8},
  ],
};

/// Descriptor for `MinimumProtocolVersion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List minimumProtocolVersionDescriptor = $convert.base64Decode(
    'ChZNaW5pbXVtUHJvdG9jb2xWZXJzaW9uEikKJU1JTklNVU1fUFJPVE9DT0xfVkVSU0lPTl9UTF'
    'NWMV8xXzIwMTYQABIiCh5NSU5JTVVNX1BST1RPQ09MX1ZFUlNJT05fU1NMVjMQARInCiNNSU5J'
    'TVVNX1BST1RPQ09MX1ZFUlNJT05fVExTVjFfMjAxNhACEikKJU1JTklNVU1fUFJPVE9DT0xfVk'
    'VSU0lPTl9UTFNWMV8yXzIwMTkQAxIpCiVNSU5JTVVNX1BST1RPQ09MX1ZFUlNJT05fVExTVjFf'
    'M18yMDI1EAQSKQolTUlOSU1VTV9QUk9UT0NPTF9WRVJTSU9OX1RMU1YxXzJfMjAyMRAFEiIKHk'
    '1JTklNVU1fUFJPVE9DT0xfVkVSU0lPTl9UTFNWMRAGEikKJU1JTklNVU1fUFJPVE9DT0xfVkVS'
    'U0lPTl9UTFNWMV8yXzIwMTgQBxIpCiVNSU5JTVVNX1BST1RPQ09MX1ZFUlNJT05fVExTVjFfMl'
    '8yMDI1EAg=');

@$core.Deprecated('Use originAccessControlOriginTypesDescriptor instead')
const OriginAccessControlOriginTypes$json = {
  '1': 'OriginAccessControlOriginTypes',
  '2': [
    {'1': 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_S3', '2': 0},
    {'1': 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_LAMBDA', '2': 1},
    {'1': 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIASTORE', '2': 2},
    {'1': 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIAPACKAGEV2', '2': 3},
  ],
};

/// Descriptor for `OriginAccessControlOriginTypes`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originAccessControlOriginTypesDescriptor = $convert.base64Decode(
    'Ch5PcmlnaW5BY2Nlc3NDb250cm9sT3JpZ2luVHlwZXMSKQolT1JJR0lOX0FDQ0VTU19DT05UUk'
    '9MX09SSUdJTl9UWVBFU19TMxAAEi0KKU9SSUdJTl9BQ0NFU1NfQ09OVFJPTF9PUklHSU5fVFlQ'
    'RVNfTEFNQkRBEAESMQotT1JJR0lOX0FDQ0VTU19DT05UUk9MX09SSUdJTl9UWVBFU19NRURJQV'
    'NUT1JFEAISNQoxT1JJR0lOX0FDQ0VTU19DT05UUk9MX09SSUdJTl9UWVBFU19NRURJQVBBQ0tB'
    'R0VWMhAD');

@$core.Deprecated('Use originAccessControlSigningBehaviorsDescriptor instead')
const OriginAccessControlSigningBehaviors$json = {
  '1': 'OriginAccessControlSigningBehaviors',
  '2': [
    {'1': 'ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_ALWAYS', '2': 0},
    {'1': 'ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NO_OVERRIDE', '2': 1},
    {'1': 'ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NEVER', '2': 2},
  ],
};

/// Descriptor for `OriginAccessControlSigningBehaviors`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originAccessControlSigningBehaviorsDescriptor =
    $convert.base64Decode(
        'CiNPcmlnaW5BY2Nlc3NDb250cm9sU2lnbmluZ0JlaGF2aW9ycxIyCi5PUklHSU5fQUNDRVNTX0'
        'NPTlRST0xfU0lHTklOR19CRUhBVklPUlNfQUxXQVlTEAASNwozT1JJR0lOX0FDQ0VTU19DT05U'
        'Uk9MX1NJR05JTkdfQkVIQVZJT1JTX05PX09WRVJSSURFEAESMQotT1JJR0lOX0FDQ0VTU19DT0'
        '5UUk9MX1NJR05JTkdfQkVIQVZJT1JTX05FVkVSEAI=');

@$core.Deprecated('Use originAccessControlSigningProtocolsDescriptor instead')
const OriginAccessControlSigningProtocols$json = {
  '1': 'OriginAccessControlSigningProtocols',
  '2': [
    {'1': 'ORIGIN_ACCESS_CONTROL_SIGNING_PROTOCOLS_SIGV4', '2': 0},
  ],
};

/// Descriptor for `OriginAccessControlSigningProtocols`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originAccessControlSigningProtocolsDescriptor =
    $convert.base64Decode(
        'CiNPcmlnaW5BY2Nlc3NDb250cm9sU2lnbmluZ1Byb3RvY29scxIxCi1PUklHSU5fQUNDRVNTX0'
        'NPTlRST0xfU0lHTklOR19QUk9UT0NPTFNfU0lHVjQQAA==');

@$core.Deprecated('Use originGroupSelectionCriteriaDescriptor instead')
const OriginGroupSelectionCriteria$json = {
  '1': 'OriginGroupSelectionCriteria',
  '2': [
    {'1': 'ORIGIN_GROUP_SELECTION_CRITERIA_MEDIAQUALITYBASED', '2': 0},
    {'1': 'ORIGIN_GROUP_SELECTION_CRITERIA_DEFAULT', '2': 1},
  ],
};

/// Descriptor for `OriginGroupSelectionCriteria`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originGroupSelectionCriteriaDescriptor =
    $convert.base64Decode(
        'ChxPcmlnaW5Hcm91cFNlbGVjdGlvbkNyaXRlcmlhEjUKMU9SSUdJTl9HUk9VUF9TRUxFQ1RJT0'
        '5fQ1JJVEVSSUFfTUVESUFRVUFMSVRZQkFTRUQQABIrCidPUklHSU5fR1JPVVBfU0VMRUNUSU9O'
        'X0NSSVRFUklBX0RFRkFVTFQQAQ==');

@$core.Deprecated('Use originProtocolPolicyDescriptor instead')
const OriginProtocolPolicy$json = {
  '1': 'OriginProtocolPolicy',
  '2': [
    {'1': 'ORIGIN_PROTOCOL_POLICY_HTTP_ONLY', '2': 0},
    {'1': 'ORIGIN_PROTOCOL_POLICY_MATCH_VIEWER', '2': 1},
    {'1': 'ORIGIN_PROTOCOL_POLICY_HTTPS_ONLY', '2': 2},
  ],
};

/// Descriptor for `OriginProtocolPolicy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originProtocolPolicyDescriptor = $convert.base64Decode(
    'ChRPcmlnaW5Qcm90b2NvbFBvbGljeRIkCiBPUklHSU5fUFJPVE9DT0xfUE9MSUNZX0hUVFBfT0'
    '5MWRAAEicKI09SSUdJTl9QUk9UT0NPTF9QT0xJQ1lfTUFUQ0hfVklFV0VSEAESJQohT1JJR0lO'
    'X1BST1RPQ09MX1BPTElDWV9IVFRQU19PTkxZEAI=');

@$core.Deprecated('Use originRequestPolicyCookieBehaviorDescriptor instead')
const OriginRequestPolicyCookieBehavior$json = {
  '1': 'OriginRequestPolicyCookieBehavior',
  '2': [
    {'1': 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALL', '2': 0},
    {'1': 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_WHITELIST', '2': 1},
    {'1': 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_NONE', '2': 2},
    {'1': 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT', '2': 3},
  ],
};

/// Descriptor for `OriginRequestPolicyCookieBehavior`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originRequestPolicyCookieBehaviorDescriptor = $convert.base64Decode(
    'CiFPcmlnaW5SZXF1ZXN0UG9saWN5Q29va2llQmVoYXZpb3ISLQopT1JJR0lOX1JFUVVFU1RfUE'
    '9MSUNZX0NPT0tJRV9CRUhBVklPUl9BTEwQABIzCi9PUklHSU5fUkVRVUVTVF9QT0xJQ1lfQ09P'
    'S0lFX0JFSEFWSU9SX1dISVRFTElTVBABEi4KKk9SSUdJTl9SRVFVRVNUX1BPTElDWV9DT09LSU'
    'VfQkVIQVZJT1JfTk9ORRACEjMKL09SSUdJTl9SRVFVRVNUX1BPTElDWV9DT09LSUVfQkVIQVZJ'
    'T1JfQUxMRVhDRVBUEAM=');

@$core.Deprecated('Use originRequestPolicyHeaderBehaviorDescriptor instead')
const OriginRequestPolicyHeaderBehavior$json = {
  '1': 'OriginRequestPolicyHeaderBehavior',
  '2': [
    {
      '1':
          'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWERANDWHITELISTCLOUDFRONT',
      '2': 0
    },
    {'1': 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_WHITELIST', '2': 1},
    {'1': 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_NONE', '2': 2},
    {'1': 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWER', '2': 3},
    {'1': 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLEXCEPT', '2': 4},
  ],
};

/// Descriptor for `OriginRequestPolicyHeaderBehavior`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originRequestPolicyHeaderBehaviorDescriptor = $convert.base64Decode(
    'CiFPcmlnaW5SZXF1ZXN0UG9saWN5SGVhZGVyQmVoYXZpb3ISSQpFT1JJR0lOX1JFUVVFU1RfUE'
    '9MSUNZX0hFQURFUl9CRUhBVklPUl9BTExWSUVXRVJBTkRXSElURUxJU1RDTE9VREZST05UEAAS'
    'MwovT1JJR0lOX1JFUVVFU1RfUE9MSUNZX0hFQURFUl9CRUhBVklPUl9XSElURUxJU1QQARIuCi'
    'pPUklHSU5fUkVRVUVTVF9QT0xJQ1lfSEVBREVSX0JFSEFWSU9SX05PTkUQAhIzCi9PUklHSU5f'
    'UkVRVUVTVF9QT0xJQ1lfSEVBREVSX0JFSEFWSU9SX0FMTFZJRVdFUhADEjMKL09SSUdJTl9SRV'
    'FVRVNUX1BPTElDWV9IRUFERVJfQkVIQVZJT1JfQUxMRVhDRVBUEAQ=');

@$core
    .Deprecated('Use originRequestPolicyQueryStringBehaviorDescriptor instead')
const OriginRequestPolicyQueryStringBehavior$json = {
  '1': 'OriginRequestPolicyQueryStringBehavior',
  '2': [
    {'1': 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALL', '2': 0},
    {'1': 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST', '2': 1},
    {'1': 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_NONE', '2': 2},
    {'1': 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT', '2': 3},
  ],
};

/// Descriptor for `OriginRequestPolicyQueryStringBehavior`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originRequestPolicyQueryStringBehaviorDescriptor =
    $convert.base64Decode(
        'CiZPcmlnaW5SZXF1ZXN0UG9saWN5UXVlcnlTdHJpbmdCZWhhdmlvchIzCi9PUklHSU5fUkVRVU'
        'VTVF9QT0xJQ1lfUVVFUllfU1RSSU5HX0JFSEFWSU9SX0FMTBAAEjkKNU9SSUdJTl9SRVFVRVNU'
        'X1BPTElDWV9RVUVSWV9TVFJJTkdfQkVIQVZJT1JfV0hJVEVMSVNUEAESNAowT1JJR0lOX1JFUV'
        'VFU1RfUE9MSUNZX1FVRVJZX1NUUklOR19CRUhBVklPUl9OT05FEAISOQo1T1JJR0lOX1JFUVVF'
        'U1RfUE9MSUNZX1FVRVJZX1NUUklOR19CRUhBVklPUl9BTExFWENFUFQQAw==');

@$core.Deprecated('Use originRequestPolicyTypeDescriptor instead')
const OriginRequestPolicyType$json = {
  '1': 'OriginRequestPolicyType',
  '2': [
    {'1': 'ORIGIN_REQUEST_POLICY_TYPE_MANAGED', '2': 0},
    {'1': 'ORIGIN_REQUEST_POLICY_TYPE_CUSTOM', '2': 1},
  ],
};

/// Descriptor for `OriginRequestPolicyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originRequestPolicyTypeDescriptor =
    $convert.base64Decode(
        'ChdPcmlnaW5SZXF1ZXN0UG9saWN5VHlwZRImCiJPUklHSU5fUkVRVUVTVF9QT0xJQ1lfVFlQRV'
        '9NQU5BR0VEEAASJQohT1JJR0lOX1JFUVVFU1RfUE9MSUNZX1RZUEVfQ1VTVE9NEAE=');

@$core.Deprecated('Use priceClassDescriptor instead')
const PriceClass$json = {
  '1': 'PriceClass',
  '2': [
    {'1': 'PRICE_CLASS_NONE', '2': 0},
    {'1': 'PRICE_CLASS_PRICECLASS_200', '2': 1},
    {'1': 'PRICE_CLASS_PRICECLASS_100', '2': 2},
    {'1': 'PRICE_CLASS_PRICECLASS_ALL', '2': 3},
  ],
};

/// Descriptor for `PriceClass`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List priceClassDescriptor = $convert.base64Decode(
    'CgpQcmljZUNsYXNzEhQKEFBSSUNFX0NMQVNTX05PTkUQABIeChpQUklDRV9DTEFTU19QUklDRU'
    'NMQVNTXzIwMBABEh4KGlBSSUNFX0NMQVNTX1BSSUNFQ0xBU1NfMTAwEAISHgoaUFJJQ0VfQ0xB'
    'U1NfUFJJQ0VDTEFTU19BTEwQAw==');

@$core.Deprecated('Use realtimeMetricsSubscriptionStatusDescriptor instead')
const RealtimeMetricsSubscriptionStatus$json = {
  '1': 'RealtimeMetricsSubscriptionStatus',
  '2': [
    {'1': 'REALTIME_METRICS_SUBSCRIPTION_STATUS_DISABLED', '2': 0},
    {'1': 'REALTIME_METRICS_SUBSCRIPTION_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `RealtimeMetricsSubscriptionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List realtimeMetricsSubscriptionStatusDescriptor =
    $convert.base64Decode(
        'CiFSZWFsdGltZU1ldHJpY3NTdWJzY3JpcHRpb25TdGF0dXMSMQotUkVBTFRJTUVfTUVUUklDU1'
        '9TVUJTQ1JJUFRJT05fU1RBVFVTX0RJU0FCTEVEEAASMAosUkVBTFRJTUVfTUVUUklDU19TVUJT'
        'Q1JJUFRJT05fU1RBVFVTX0VOQUJMRUQQAQ==');

@$core.Deprecated('Use referrerPolicyListDescriptor instead')
const ReferrerPolicyList$json = {
  '1': 'ReferrerPolicyList',
  '2': [
    {'1': 'REFERRER_POLICY_LIST_UNSAFE_URL', '2': 0},
    {'1': 'REFERRER_POLICY_LIST_STRICT_ORIGIN_WHEN_CROSS_ORIGIN', '2': 1},
    {'1': 'REFERRER_POLICY_LIST_NO_REFERRER', '2': 2},
    {'1': 'REFERRER_POLICY_LIST_NO_REFERRER_WHEN_DOWNGRADE', '2': 3},
    {'1': 'REFERRER_POLICY_LIST_ORIGIN', '2': 4},
    {'1': 'REFERRER_POLICY_LIST_STRICT_ORIGIN', '2': 5},
    {'1': 'REFERRER_POLICY_LIST_ORIGIN_WHEN_CROSS_ORIGIN', '2': 6},
    {'1': 'REFERRER_POLICY_LIST_SAME_ORIGIN', '2': 7},
  ],
};

/// Descriptor for `ReferrerPolicyList`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List referrerPolicyListDescriptor = $convert.base64Decode(
    'ChJSZWZlcnJlclBvbGljeUxpc3QSIwofUkVGRVJSRVJfUE9MSUNZX0xJU1RfVU5TQUZFX1VSTB'
    'AAEjgKNFJFRkVSUkVSX1BPTElDWV9MSVNUX1NUUklDVF9PUklHSU5fV0hFTl9DUk9TU19PUklH'
    'SU4QARIkCiBSRUZFUlJFUl9QT0xJQ1lfTElTVF9OT19SRUZFUlJFUhACEjMKL1JFRkVSUkVSX1'
    'BPTElDWV9MSVNUX05PX1JFRkVSUkVSX1dIRU5fRE9XTkdSQURFEAMSHwobUkVGRVJSRVJfUE9M'
    'SUNZX0xJU1RfT1JJR0lOEAQSJgoiUkVGRVJSRVJfUE9MSUNZX0xJU1RfU1RSSUNUX09SSUdJTh'
    'AFEjEKLVJFRkVSUkVSX1BPTElDWV9MSVNUX09SSUdJTl9XSEVOX0NST1NTX09SSUdJThAGEiQK'
    'IFJFRkVSUkVSX1BPTElDWV9MSVNUX1NBTUVfT1JJR0lOEAc=');

@$core.Deprecated(
    'Use responseHeadersPolicyAccessControlAllowMethodsValuesDescriptor instead')
const ResponseHeadersPolicyAccessControlAllowMethodsValues$json = {
  '1': 'ResponseHeadersPolicyAccessControlAllowMethodsValues',
  '2': [
    {
      '1':
          'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_OPTIONS',
      '2': 0
    },
    {
      '1': 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PATCH',
      '2': 1
    },
    {
      '1': 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_POST',
      '2': 2
    },
    {
      '1': 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_HEAD',
      '2': 3
    },
    {
      '1': 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_GET',
      '2': 4
    },
    {
      '1': 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_DELETE',
      '2': 5
    },
    {
      '1': 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_ALL',
      '2': 6
    },
    {
      '1': 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PUT',
      '2': 7
    },
  ],
};

/// Descriptor for `ResponseHeadersPolicyAccessControlAllowMethodsValues`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyAccessControlAllowMethodsValuesDescriptor = $convert.base64Decode(
    'CjRSZXNwb25zZUhlYWRlcnNQb2xpY3lBY2Nlc3NDb250cm9sQWxsb3dNZXRob2RzVmFsdWVzEk'
    'cKQ1JFU1BPTlNFX0hFQURFUlNfUE9MSUNZX0FDQ0VTU19DT05UUk9MX0FMTE9XX01FVEhPRFNf'
    'VkFMVUVTX09QVElPTlMQABJFCkFSRVNQT05TRV9IRUFERVJTX1BPTElDWV9BQ0NFU1NfQ09OVF'
    'JPTF9BTExPV19NRVRIT0RTX1ZBTFVFU19QQVRDSBABEkQKQFJFU1BPTlNFX0hFQURFUlNfUE9M'
    'SUNZX0FDQ0VTU19DT05UUk9MX0FMTE9XX01FVEhPRFNfVkFMVUVTX1BPU1QQAhJECkBSRVNQT0'
    '5TRV9IRUFERVJTX1BPTElDWV9BQ0NFU1NfQ09OVFJPTF9BTExPV19NRVRIT0RTX1ZBTFVFU19I'
    'RUFEEAMSQwo/UkVTUE9OU0VfSEVBREVSU19QT0xJQ1lfQUNDRVNTX0NPTlRST0xfQUxMT1dfTU'
    'VUSE9EU19WQUxVRVNfR0VUEAQSRgpCUkVTUE9OU0VfSEVBREVSU19QT0xJQ1lfQUNDRVNTX0NP'
    'TlRST0xfQUxMT1dfTUVUSE9EU19WQUxVRVNfREVMRVRFEAUSQwo/UkVTUE9OU0VfSEVBREVSU1'
    '9QT0xJQ1lfQUNDRVNTX0NPTlRST0xfQUxMT1dfTUVUSE9EU19WQUxVRVNfQUxMEAYSQwo/UkVT'
    'UE9OU0VfSEVBREVSU19QT0xJQ1lfQUNDRVNTX0NPTlRST0xfQUxMT1dfTUVUSE9EU19WQUxVRV'
    'NfUFVUEAc=');

@$core.Deprecated('Use responseHeadersPolicyTypeDescriptor instead')
const ResponseHeadersPolicyType$json = {
  '1': 'ResponseHeadersPolicyType',
  '2': [
    {'1': 'RESPONSE_HEADERS_POLICY_TYPE_MANAGED', '2': 0},
    {'1': 'RESPONSE_HEADERS_POLICY_TYPE_CUSTOM', '2': 1},
  ],
};

/// Descriptor for `ResponseHeadersPolicyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyTypeDescriptor = $convert.base64Decode(
    'ChlSZXNwb25zZUhlYWRlcnNQb2xpY3lUeXBlEigKJFJFU1BPTlNFX0hFQURFUlNfUE9MSUNZX1'
    'RZUEVfTUFOQUdFRBAAEicKI1JFU1BPTlNFX0hFQURFUlNfUE9MSUNZX1RZUEVfQ1VTVE9NEAE=');

@$core.Deprecated('Use sSLSupportMethodDescriptor instead')
const SSLSupportMethod$json = {
  '1': 'SSLSupportMethod',
  '2': [
    {'1': 'S_S_L_SUPPORT_METHOD_STATIC_IP', '2': 0},
    {'1': 'S_S_L_SUPPORT_METHOD_VIP', '2': 1},
    {'1': 'S_S_L_SUPPORT_METHOD_SNI_ONLY', '2': 2},
  ],
};

/// Descriptor for `SSLSupportMethod`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sSLSupportMethodDescriptor = $convert.base64Decode(
    'ChBTU0xTdXBwb3J0TWV0aG9kEiIKHlNfU19MX1NVUFBPUlRfTUVUSE9EX1NUQVRJQ19JUBAAEh'
    'wKGFNfU19MX1NVUFBPUlRfTUVUSE9EX1ZJUBABEiEKHVNfU19MX1NVUFBPUlRfTUVUSE9EX1NO'
    'SV9PTkxZEAI=');

@$core.Deprecated('Use sslProtocolDescriptor instead')
const SslProtocol$json = {
  '1': 'SslProtocol',
  '2': [
    {'1': 'SSL_PROTOCOL_SSLV3', '2': 0},
    {'1': 'SSL_PROTOCOL_TLSV1_1', '2': 1},
    {'1': 'SSL_PROTOCOL_TLSV1', '2': 2},
    {'1': 'SSL_PROTOCOL_TLSV1_2', '2': 3},
  ],
};

/// Descriptor for `SslProtocol`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sslProtocolDescriptor = $convert.base64Decode(
    'CgtTc2xQcm90b2NvbBIWChJTU0xfUFJPVE9DT0xfU1NMVjMQABIYChRTU0xfUFJPVE9DT0xfVE'
    'xTVjFfMRABEhYKElNTTF9QUk9UT0NPTF9UTFNWMRACEhgKFFNTTF9QUk9UT0NPTF9UTFNWMV8y'
    'EAM=');

@$core.Deprecated('Use trustStoreStatusDescriptor instead')
const TrustStoreStatus$json = {
  '1': 'TrustStoreStatus',
  '2': [
    {'1': 'TRUST_STORE_STATUS_ACTIVE', '2': 0},
    {'1': 'TRUST_STORE_STATUS_FAILED', '2': 1},
    {'1': 'TRUST_STORE_STATUS_PENDING', '2': 2},
  ],
};

/// Descriptor for `TrustStoreStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List trustStoreStatusDescriptor = $convert.base64Decode(
    'ChBUcnVzdFN0b3JlU3RhdHVzEh0KGVRSVVNUX1NUT1JFX1NUQVRVU19BQ1RJVkUQABIdChlUUl'
    'VTVF9TVE9SRV9TVEFUVVNfRkFJTEVEEAESHgoaVFJVU1RfU1RPUkVfU1RBVFVTX1BFTkRJTkcQ'
    'Ag==');

@$core.Deprecated('Use validationTokenHostDescriptor instead')
const ValidationTokenHost$json = {
  '1': 'ValidationTokenHost',
  '2': [
    {'1': 'VALIDATION_TOKEN_HOST_SELFHOSTED', '2': 0},
    {'1': 'VALIDATION_TOKEN_HOST_CLOUDFRONT', '2': 1},
  ],
};

/// Descriptor for `ValidationTokenHost`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List validationTokenHostDescriptor = $convert.base64Decode(
    'ChNWYWxpZGF0aW9uVG9rZW5Ib3N0EiQKIFZBTElEQVRJT05fVE9LRU5fSE9TVF9TRUxGSE9TVE'
    'VEEAASJAogVkFMSURBVElPTl9UT0tFTl9IT1NUX0NMT1VERlJPTlQQAQ==');

@$core.Deprecated('Use viewerMtlsModeDescriptor instead')
const ViewerMtlsMode$json = {
  '1': 'ViewerMtlsMode',
  '2': [
    {'1': 'VIEWER_MTLS_MODE_OPTIONAL', '2': 0},
    {'1': 'VIEWER_MTLS_MODE_REQUIRED', '2': 1},
  ],
};

/// Descriptor for `ViewerMtlsMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List viewerMtlsModeDescriptor = $convert.base64Decode(
    'Cg5WaWV3ZXJNdGxzTW9kZRIdChlWSUVXRVJfTVRMU19NT0RFX09QVElPTkFMEAASHQoZVklFV0'
    'VSX01UTFNfTU9ERV9SRVFVSVJFRBAB');

@$core.Deprecated('Use viewerProtocolPolicyDescriptor instead')
const ViewerProtocolPolicy$json = {
  '1': 'ViewerProtocolPolicy',
  '2': [
    {'1': 'VIEWER_PROTOCOL_POLICY_ALLOW_ALL', '2': 0},
    {'1': 'VIEWER_PROTOCOL_POLICY_HTTPS_ONLY', '2': 1},
    {'1': 'VIEWER_PROTOCOL_POLICY_REDIRECT_TO_HTTPS', '2': 2},
  ],
};

/// Descriptor for `ViewerProtocolPolicy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List viewerProtocolPolicyDescriptor = $convert.base64Decode(
    'ChRWaWV3ZXJQcm90b2NvbFBvbGljeRIkCiBWSUVXRVJfUFJPVE9DT0xfUE9MSUNZX0FMTE9XX0'
    'FMTBAAEiUKIVZJRVdFUl9QUk9UT0NPTF9QT0xJQ1lfSFRUUFNfT05MWRABEiwKKFZJRVdFUl9Q'
    'Uk9UT0NPTF9QT0xJQ1lfUkVESVJFQ1RfVE9fSFRUUFMQAg==');

@$core.Deprecated('Use accessDeniedDescriptor instead')
const AccessDenied$json = {
  '1': 'AccessDenied',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccessDenied`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedDescriptor = $convert.base64Decode(
    'CgxBY2Nlc3NEZW5pZWQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use activeTrustedKeyGroupsDescriptor instead')
const ActiveTrustedKeyGroups$json = {
  '1': 'ActiveTrustedKeyGroups',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.KGKeyPairIds',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ActiveTrustedKeyGroups`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activeTrustedKeyGroupsDescriptor = $convert.base64Decode(
    'ChZBY3RpdmVUcnVzdGVkS2V5R3JvdXBzEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVkEj'
    'EKBWl0ZW1zGLDw2AEgAygLMhguY2xvdWRmcm9udC5LR0tleVBhaXJJZHNSBWl0ZW1zEh0KCHF1'
    'YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use activeTrustedSignersDescriptor instead')
const ActiveTrustedSigners$json = {
  '1': 'ActiveTrustedSigners',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.Signer',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ActiveTrustedSigners`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activeTrustedSignersDescriptor = $convert.base64Decode(
    'ChRBY3RpdmVUcnVzdGVkU2lnbmVycxIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIrCg'
    'VpdGVtcxiw8NgBIAMoCzISLmNsb3VkZnJvbnQuU2lnbmVyUgVpdGVtcxIdCghxdWFudGl0eRj5'
    '5dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use aliasICPRecordalDescriptor instead')
const AliasICPRecordal$json = {
  '1': 'AliasICPRecordal',
  '2': [
    {'1': 'cname', '3': 420235386, '4': 1, '5': 9, '10': 'cname'},
    {
      '1': 'icprecordalstatus',
      '3': 230990942,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ICPRecordalStatus',
      '10': 'icprecordalstatus'
    },
  ],
};

/// Descriptor for `AliasICPRecordal`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aliasICPRecordalDescriptor = $convert.base64Decode(
    'ChBBbGlhc0lDUFJlY29yZGFsEhgKBWNuYW1lGPqQscgBIAEoCVIFY25hbWUSTgoRaWNwcmVjb3'
    'JkYWxzdGF0dXMY3siSbiABKA4yHS5jbG91ZGZyb250LklDUFJlY29yZGFsU3RhdHVzUhFpY3By'
    'ZWNvcmRhbHN0YXR1cw==');

@$core.Deprecated('Use aliasesDescriptor instead')
const Aliases$json = {
  '1': 'Aliases',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `Aliases`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aliasesDescriptor = $convert.base64Decode(
    'CgdBbGlhc2VzEhcKBWl0ZW1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBV'
    'IIcXVhbnRpdHk=');

@$core.Deprecated('Use allowedMethodsDescriptor instead')
const AllowedMethods$json = {
  '1': 'AllowedMethods',
  '2': [
    {
      '1': 'cachedmethods',
      '3': 62836516,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachedMethods',
      '10': 'cachedmethods'
    },
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 14,
      '6': '.cloudfront.Method',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `AllowedMethods`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List allowedMethodsDescriptor = $convert.base64Decode(
    'Cg5BbGxvd2VkTWV0aG9kcxJCCg1jYWNoZWRtZXRob2RzGKSe+x0gASgLMhkuY2xvdWRmcm9udC'
    '5DYWNoZWRNZXRob2RzUg1jYWNoZWRtZXRob2RzEisKBWl0ZW1zGLDw2AEgAygOMhIuY2xvdWRm'
    'cm9udC5NZXRob2RSBWl0ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use anycastIpListDescriptor instead')
const AnycastIpList$json = {
  '1': 'AnycastIpList',
  '2': [
    {'1': 'anycastips', '3': 113332317, '4': 3, '5': 9, '10': 'anycastips'},
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'ipaddresstype',
      '3': 459110693,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.IpAddressType',
      '10': 'ipaddresstype'
    },
    {'1': 'ipcount', '3': 475138532, '4': 1, '5': 5, '10': 'ipcount'},
    {
      '1': 'ipamconfig',
      '3': 310812553,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.IpamConfig',
      '10': 'ipamconfig'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `AnycastIpList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anycastIpListDescriptor = $convert.base64Decode(
    'Cg1BbnljYXN0SXBMaXN0EiEKCmFueWNhc3RpcHMY3aCFNiADKAlSCmFueWNhc3RpcHMSFAoDYX'
    'JuGJ2b7b8BIAEoCVIDYXJuEhIKAmlkGIHyorcBIAEoCVICaWQSQwoNaXBhZGRyZXNzdHlwZRil'
    '8vXaASABKA4yGS5jbG91ZGZyb250LklwQWRkcmVzc1R5cGVSDWlwYWRkcmVzc3R5cGUSHAoHaX'
    'Bjb3VudBjkk8jiASABKAVSB2lwY291bnQSOgoKaXBhbWNvbmZpZxiJv5qUASABKAsyFi5jbG91'
    'ZGZyb250LklwYW1Db25maWdSCmlwYW1jb25maWcSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIA'
    'EoCVIQbGFzdG1vZGlmaWVkdGltZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhkKBnN0YXR1cxiQ'
    '5PsCIAEoCVIGc3RhdHVz');

@$core.Deprecated('Use anycastIpListCollectionDescriptor instead')
const AnycastIpListCollection$json = {
  '1': 'AnycastIpListCollection',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.AnycastIpListSummary',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `AnycastIpListCollection`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anycastIpListCollectionDescriptor = $convert.base64Decode(
    'ChdBbnljYXN0SXBMaXN0Q29sbGVjdGlvbhIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cn'
    'VuY2F0ZWQSOQoFaXRlbXMYsPDYASADKAsyIC5jbG91ZGZyb250LkFueWNhc3RJcExpc3RTdW1t'
    'YXJ5UgVpdGVtcxIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trxAS'
    'ABKAVSCG1heGl0ZW1zEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2VyEh0KCHF1'
    'YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use anycastIpListSummaryDescriptor instead')
const AnycastIpListSummary$json = {
  '1': 'AnycastIpListSummary',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'ipaddresstype',
      '3': 459110693,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.IpAddressType',
      '10': 'ipaddresstype'
    },
    {'1': 'ipcount', '3': 475138532, '4': 1, '5': 5, '10': 'ipcount'},
    {
      '1': 'ipamconfig',
      '3': 310812553,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.IpamConfig',
      '10': 'ipamconfig'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `AnycastIpListSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anycastIpListSummaryDescriptor = $convert.base64Decode(
    'ChRBbnljYXN0SXBMaXN0U3VtbWFyeRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SFgoEZXRhZxiB37'
    'OVASABKAlSBGV0YWcSEgoCaWQYgfKitwEgASgJUgJpZBJDCg1pcGFkZHJlc3N0eXBlGKXy9doB'
    'IAEoDjIZLmNsb3VkZnJvbnQuSXBBZGRyZXNzVHlwZVINaXBhZGRyZXNzdHlwZRIcCgdpcGNvdW'
    '50GOSTyOIBIAEoBVIHaXBjb3VudBI6CgppcGFtY29uZmlnGIm/mpQBIAEoCzIWLmNsb3VkZnJv'
    'bnQuSXBhbUNvbmZpZ1IKaXBhbWNvbmZpZxItChBsYXN0bW9kaWZpZWR0aW1lGOCC/HAgASgJUh'
    'BsYXN0bW9kaWZpZWR0aW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGQoGc3RhdHVzGJDk+wIg'
    'ASgJUgZzdGF0dXM=');

@$core.Deprecated('Use associateAliasRequestDescriptor instead')
const AssociateAliasRequest$json = {
  '1': 'AssociateAliasRequest',
  '2': [
    {'1': 'alias', '3': 48362232, '4': 1, '5': 9, '10': 'alias'},
    {
      '1': 'targetdistributionid',
      '3': 250750584,
      '4': 1,
      '5': 9,
      '10': 'targetdistributionid'
    },
  ],
};

/// Descriptor for `AssociateAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateAliasRequestDescriptor = $convert.base64Decode(
    'ChVBc3NvY2lhdGVBbGlhc1JlcXVlc3QSFwoFYWxpYXMY+OWHFyABKAlSBWFsaWFzEjUKFHRhcm'
    'dldGRpc3RyaWJ1dGlvbmlkGPjMyHcgASgJUhR0YXJnZXRkaXN0cmlidXRpb25pZA==');

@$core.Deprecated(
    'Use associateDistributionTenantWebACLRequestDescriptor instead')
const AssociateDistributionTenantWebACLRequest$json = {
  '1': 'AssociateDistributionTenantWebACLRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'webaclarn', '3': 82506659, '4': 1, '5': 9, '10': 'webaclarn'},
  ],
};

/// Descriptor for `AssociateDistributionTenantWebACLRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateDistributionTenantWebACLRequestDescriptor =
    $convert.base64Decode(
        'CihBc3NvY2lhdGVEaXN0cmlidXRpb25UZW5hbnRXZWJBQ0xSZXF1ZXN0EhIKAmlkGIHyorcBIA'
        'EoCVICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaBIfCgl3ZWJhY2xhcm4Yo+erJyAB'
        'KAlSCXdlYmFjbGFybg==');

@$core
    .Deprecated('Use associateDistributionTenantWebACLResultDescriptor instead')
const AssociateDistributionTenantWebACLResult$json = {
  '1': 'AssociateDistributionTenantWebACLResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'webaclarn', '3': 82506659, '4': 1, '5': 9, '10': 'webaclarn'},
  ],
};

/// Descriptor for `AssociateDistributionTenantWebACLResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateDistributionTenantWebACLResultDescriptor =
    $convert.base64Decode(
        'CidBc3NvY2lhdGVEaXN0cmlidXRpb25UZW5hbnRXZWJBQ0xSZXN1bHQSFgoEZXRhZxiB37OVAS'
        'ABKAlSBGV0YWcSEgoCaWQYgfKitwEgASgJUgJpZBIfCgl3ZWJhY2xhcm4Yo+erJyABKAlSCXdl'
        'YmFjbGFybg==');

@$core.Deprecated('Use associateDistributionWebACLRequestDescriptor instead')
const AssociateDistributionWebACLRequest$json = {
  '1': 'AssociateDistributionWebACLRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'webaclarn', '3': 82506659, '4': 1, '5': 9, '10': 'webaclarn'},
  ],
};

/// Descriptor for `AssociateDistributionWebACLRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateDistributionWebACLRequestDescriptor =
    $convert.base64Decode(
        'CiJBc3NvY2lhdGVEaXN0cmlidXRpb25XZWJBQ0xSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'QSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaBIfCgl3ZWJhY2xhcm4Yo+erJyABKAlSCXdl'
        'YmFjbGFybg==');

@$core.Deprecated('Use associateDistributionWebACLResultDescriptor instead')
const AssociateDistributionWebACLResult$json = {
  '1': 'AssociateDistributionWebACLResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'webaclarn', '3': 82506659, '4': 1, '5': 9, '10': 'webaclarn'},
  ],
};

/// Descriptor for `AssociateDistributionWebACLResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateDistributionWebACLResultDescriptor =
    $convert.base64Decode(
        'CiFBc3NvY2lhdGVEaXN0cmlidXRpb25XZWJBQ0xSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBG'
        'V0YWcSEgoCaWQYgfKitwEgASgJUgJpZBIfCgl3ZWJhY2xhcm4Yo+erJyABKAlSCXdlYmFjbGFy'
        'bg==');

@$core.Deprecated('Use batchTooLargeDescriptor instead')
const BatchTooLarge$json = {
  '1': 'BatchTooLarge',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BatchTooLarge`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchTooLargeDescriptor = $convert.base64Decode(
    'Cg1CYXRjaFRvb0xhcmdlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use cNAMEAlreadyExistsDescriptor instead')
const CNAMEAlreadyExists$json = {
  '1': 'CNAMEAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CNAMEAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cNAMEAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChJDTkFNRUFscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use caCertificatesBundleS3LocationDescriptor instead')
const CaCertificatesBundleS3Location$json = {
  '1': 'CaCertificatesBundleS3Location',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `CaCertificatesBundleS3Location`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List caCertificatesBundleS3LocationDescriptor =
    $convert.base64Decode(
        'Ch5DYUNlcnRpZmljYXRlc0J1bmRsZVMzTG9jYXRpb24SGQoGYnVja2V0GNjquBogASgJUgZidW'
        'NrZXQSEwoDa2V5GI2S62ggASgJUgNrZXkSGQoGcmVnaW9uGJ7xuUkgASgJUgZyZWdpb24SHAoH'
        'dmVyc2lvbhi4qrfuASABKAlSB3ZlcnNpb24=');

@$core.Deprecated('Use caCertificatesBundleSourceDescriptor instead')
const CaCertificatesBundleSource$json = {
  '1': 'CaCertificatesBundleSource',
  '2': [
    {
      '1': 'cacertificatesbundles3location',
      '3': 113649925,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CaCertificatesBundleS3Location',
      '10': 'cacertificatesbundles3location'
    },
  ],
};

/// Descriptor for `CaCertificatesBundleSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List caCertificatesBundleSourceDescriptor =
    $convert.base64Decode(
        'ChpDYUNlcnRpZmljYXRlc0J1bmRsZVNvdXJjZRJ1Ch5jYWNlcnRpZmljYXRlc2J1bmRsZXMzbG'
        '9jYXRpb24YhdKYNiABKAsyKi5jbG91ZGZyb250LkNhQ2VydGlmaWNhdGVzQnVuZGxlUzNMb2Nh'
        'dGlvblIeY2FjZXJ0aWZpY2F0ZXNidW5kbGVzM2xvY2F0aW9u');

@$core.Deprecated('Use cacheBehaviorDescriptor instead')
const CacheBehavior$json = {
  '1': 'CacheBehavior',
  '2': [
    {
      '1': 'allowedmethods',
      '3': 56383476,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.AllowedMethods',
      '10': 'allowedmethods'
    },
    {
      '1': 'cachepolicyid',
      '3': 431434163,
      '4': 1,
      '5': 9,
      '10': 'cachepolicyid'
    },
    {'1': 'compress', '3': 235468462, '4': 1, '5': 8, '10': 'compress'},
    {'1': 'defaultttl', '3': 391646391, '4': 1, '5': 3, '10': 'defaultttl'},
    {
      '1': 'fieldlevelencryptionid',
      '3': 450714616,
      '4': 1,
      '5': 9,
      '10': 'fieldlevelencryptionid'
    },
    {
      '1': 'forwardedvalues',
      '3': 34815362,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ForwardedValues',
      '10': 'forwardedvalues'
    },
    {
      '1': 'functionassociations',
      '3': 457445650,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionAssociations',
      '10': 'functionassociations'
    },
    {
      '1': 'grpcconfig',
      '3': 406090728,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.GrpcConfig',
      '10': 'grpcconfig'
    },
    {
      '1': 'lambdafunctionassociations',
      '3': 46888655,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.LambdaFunctionAssociations',
      '10': 'lambdafunctionassociations'
    },
    {'1': 'maxttl', '3': 451484784, '4': 1, '5': 3, '10': 'maxttl'},
    {'1': 'minttl', '3': 420784162, '4': 1, '5': 3, '10': 'minttl'},
    {
      '1': 'originrequestpolicyid',
      '3': 298538616,
      '4': 1,
      '5': 9,
      '10': 'originrequestpolicyid'
    },
    {'1': 'pathpattern', '3': 266478053, '4': 1, '5': 9, '10': 'pathpattern'},
    {
      '1': 'realtimelogconfigarn',
      '3': 152963408,
      '4': 1,
      '5': 9,
      '10': 'realtimelogconfigarn'
    },
    {
      '1': 'responseheaderspolicyid',
      '3': 244029524,
      '4': 1,
      '5': 9,
      '10': 'responseheaderspolicyid'
    },
    {
      '1': 'smoothstreaming',
      '3': 92667114,
      '4': 1,
      '5': 8,
      '10': 'smoothstreaming'
    },
    {
      '1': 'targetoriginid',
      '3': 381807144,
      '4': 1,
      '5': 9,
      '10': 'targetoriginid'
    },
    {
      '1': 'trustedkeygroups',
      '3': 436720164,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustedKeyGroups',
      '10': 'trustedkeygroups'
    },
    {
      '1': 'trustedsigners',
      '3': 82003176,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustedSigners',
      '10': 'trustedsigners'
    },
    {
      '1': 'viewerprotocolpolicy',
      '3': 107534830,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ViewerProtocolPolicy',
      '10': 'viewerprotocolpolicy'
    },
  ],
};

/// Descriptor for `CacheBehavior`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cacheBehaviorDescriptor = $convert.base64Decode(
    'Cg1DYWNoZUJlaGF2aW9yEkUKDmFsbG93ZWRtZXRob2RzGPSv8RogASgLMhouY2xvdWRmcm9udC'
    '5BbGxvd2VkTWV0aG9kc1IOYWxsb3dlZG1ldGhvZHMSKAoNY2FjaGVwb2xpY3lpZBiz09zNASAB'
    'KAlSDWNhY2hlcG9saWN5aWQSHQoIY29tcHJlc3MYru2jcCABKAhSCGNvbXByZXNzEiIKCmRlZm'
    'F1bHR0dGwYt5ngugEgASgDUgpkZWZhdWx0dHRsEjoKFmZpZWxkbGV2ZWxlbmNyeXB0aW9uaWQY'
    '+Lf11gEgASgJUhZmaWVsZGxldmVsZW5jcnlwdGlvbmlkEkgKD2ZvcndhcmRlZHZhbHVlcxiC+8'
    'wQIAEoCzIbLmNsb3VkZnJvbnQuRm9yd2FyZGVkVmFsdWVzUg9mb3J3YXJkZWR2YWx1ZXMSWAoU'
    'ZnVuY3Rpb25hc3NvY2lhdGlvbnMYkqKQ2gEgASgLMiAuY2xvdWRmcm9udC5GdW5jdGlvbkFzc2'
    '9jaWF0aW9uc1IUZnVuY3Rpb25hc3NvY2lhdGlvbnMSOgoKZ3JwY2NvbmZpZxjo59HBASABKAsy'
    'Fi5jbG91ZGZyb250LkdycGNDb25maWdSCmdycGNjb25maWcSaQoabGFtYmRhZnVuY3Rpb25hc3'
    'NvY2lhdGlvbnMYz+2tFiABKAsyJi5jbG91ZGZyb250LkxhbWJkYUZ1bmN0aW9uQXNzb2NpYXRp'
    'b25zUhpsYW1iZGFmdW5jdGlvbmFzc29jaWF0aW9ucxIaCgZtYXh0dGwY8Lik1wEgASgDUgZtYX'
    'h0dGwSGgoGbWludHRsGKLQ0sgBIAEoA1IGbWludHRsEjgKFW9yaWdpbnJlcXVlc3Rwb2xpY3lp'
    'ZBj4rK2OASABKAlSFW9yaWdpbnJlcXVlc3Rwb2xpY3lpZBIjCgtwYXRocGF0dGVybhjlw4h/IA'
    'EoCVILcGF0aHBhdHRlcm4SNQoUcmVhbHRpbWVsb2djb25maWdhcm4Y0JL4SCABKAlSFHJlYWx0'
    'aW1lbG9nY29uZmlnYXJuEjsKF3Jlc3BvbnNlaGVhZGVyc3BvbGljeWlkGNSwrnQgASgJUhdyZX'
    'Nwb25zZWhlYWRlcnNwb2xpY3lpZBIrCg9zbW9vdGhzdHJlYW1pbmcY6vmXLCABKAhSD3Ntb290'
    'aHN0cmVhbWluZxIqCg50YXJnZXRvcmlnaW5pZBio1Ie2ASABKAlSDnRhcmdldG9yaWdpbmlkEk'
    'wKEHRydXN0ZWRrZXlncm91cHMYpKSf0AEgASgLMhwuY2xvdWRmcm9udC5UcnVzdGVkS2V5R3Jv'
    'dXBzUhB0cnVzdGVka2V5Z3JvdXBzEkUKDnRydXN0ZWRzaWduZXJzGOiJjScgASgLMhouY2xvdW'
    'Rmcm9udC5UcnVzdGVkU2lnbmVyc1IOdHJ1c3RlZHNpZ25lcnMSVwoUdmlld2VycHJvdG9jb2xw'
    'b2xpY3kY7rOjMyABKA4yIC5jbG91ZGZyb250LlZpZXdlclByb3RvY29sUG9saWN5UhR2aWV3ZX'
    'Jwcm90b2NvbHBvbGljeQ==');

@$core.Deprecated('Use cacheBehaviorsDescriptor instead')
const CacheBehaviors$json = {
  '1': 'CacheBehaviors',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.CacheBehavior',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `CacheBehaviors`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cacheBehaviorsDescriptor = $convert.base64Decode(
    'Cg5DYWNoZUJlaGF2aW9ycxIyCgVpdGVtcxiw8NgBIAMoCzIZLmNsb3VkZnJvbnQuQ2FjaGVCZW'
    'hhdmlvclIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use cachePolicyDescriptor instead')
const CachePolicy$json = {
  '1': 'CachePolicy',
  '2': [
    {
      '1': 'cachepolicyconfig',
      '3': 407094126,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyConfig',
      '10': 'cachepolicyconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `CachePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyDescriptor = $convert.base64Decode(
    'CgtDYWNoZVBvbGljeRJPChFjYWNoZXBvbGljeWNvbmZpZxjuho/CASABKAsyHS5jbG91ZGZyb2'
    '50LkNhY2hlUG9saWN5Q29uZmlnUhFjYWNoZXBvbGljeWNvbmZpZxISCgJpZBiB8qK3ASABKAlS'
    'AmlkEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWU=');

@$core.Deprecated('Use cachePolicyAlreadyExistsDescriptor instead')
const CachePolicyAlreadyExists$json = {
  '1': 'CachePolicyAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CachePolicyAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChhDYWNoZVBvbGljeUFscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use cachePolicyConfigDescriptor instead')
const CachePolicyConfig$json = {
  '1': 'CachePolicyConfig',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'defaultttl', '3': 391646391, '4': 1, '5': 3, '10': 'defaultttl'},
    {'1': 'maxttl', '3': 451484784, '4': 1, '5': 3, '10': 'maxttl'},
    {'1': 'minttl', '3': 420784162, '4': 1, '5': 3, '10': 'minttl'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parametersincachekeyandforwardedtoorigin',
      '3': 298859768,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ParametersInCacheKeyAndForwardedToOrigin',
      '10': 'parametersincachekeyandforwardedtoorigin'
    },
  ],
};

/// Descriptor for `CachePolicyConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyConfigDescriptor = $convert.base64Decode(
    'ChFDYWNoZVBvbGljeUNvbmZpZxIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbWVudBIiCgpkZW'
    'ZhdWx0dHRsGLeZ4LoBIAEoA1IKZGVmYXVsdHR0bBIaCgZtYXh0dGwY8Lik1wEgASgDUgZtYXh0'
    'dGwSGgoGbWludHRsGKLQ0sgBIAEoA1IGbWludHRsEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSlA'
    'EKKHBhcmFtZXRlcnNpbmNhY2hla2V5YW5kZm9yd2FyZGVkdG9vcmlnaW4Y+PnAjgEgASgLMjQu'
    'Y2xvdWRmcm9udC5QYXJhbWV0ZXJzSW5DYWNoZUtleUFuZEZvcndhcmRlZFRvT3JpZ2luUihwYX'
    'JhbWV0ZXJzaW5jYWNoZWtleWFuZGZvcndhcmRlZHRvb3JpZ2lu');

@$core.Deprecated('Use cachePolicyCookiesConfigDescriptor instead')
const CachePolicyCookiesConfig$json = {
  '1': 'CachePolicyCookiesConfig',
  '2': [
    {
      '1': 'cookiebehavior',
      '3': 86223930,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CachePolicyCookieBehavior',
      '10': 'cookiebehavior'
    },
    {
      '1': 'cookies',
      '3': 418634949,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CookieNames',
      '10': 'cookies'
    },
  ],
};

/// Descriptor for `CachePolicyCookiesConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyCookiesConfigDescriptor = $convert.base64Decode(
    'ChhDYWNoZVBvbGljeUNvb2tpZXNDb25maWcSUAoOY29va2llYmVoYXZpb3IYutiOKSABKA4yJS'
    '5jbG91ZGZyb250LkNhY2hlUG9saWN5Q29va2llQmVoYXZpb3JSDmNvb2tpZWJlaGF2aW9yEjUK'
    'B2Nvb2tpZXMYxbnPxwEgASgLMhcuY2xvdWRmcm9udC5Db29raWVOYW1lc1IHY29va2llcw==');

@$core.Deprecated('Use cachePolicyHeadersConfigDescriptor instead')
const CachePolicyHeadersConfig$json = {
  '1': 'CachePolicyHeadersConfig',
  '2': [
    {
      '1': 'headerbehavior',
      '3': 348515467,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CachePolicyHeaderBehavior',
      '10': 'headerbehavior'
    },
    {
      '1': 'headers',
      '3': 323967370,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Headers',
      '10': 'headers'
    },
  ],
};

/// Descriptor for `CachePolicyHeadersConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyHeadersConfigDescriptor = $convert.base64Decode(
    'ChhDYWNoZVBvbGljeUhlYWRlcnNDb25maWcSUQoOaGVhZGVyYmVoYXZpb3IYi9mXpgEgASgOMi'
    'UuY2xvdWRmcm9udC5DYWNoZVBvbGljeUhlYWRlckJlaGF2aW9yUg5oZWFkZXJiZWhhdmlvchIx'
    'CgdoZWFkZXJzGIqzvZoBIAEoCzITLmNsb3VkZnJvbnQuSGVhZGVyc1IHaGVhZGVycw==');

@$core.Deprecated('Use cachePolicyInUseDescriptor instead')
const CachePolicyInUse$json = {
  '1': 'CachePolicyInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CachePolicyInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyInUseDescriptor = $convert.base64Decode(
    'ChBDYWNoZVBvbGljeUluVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use cachePolicyListDescriptor instead')
const CachePolicyList$json = {
  '1': 'CachePolicyList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.CachePolicySummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `CachePolicyList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyListDescriptor = $convert.base64Decode(
    'Cg9DYWNoZVBvbGljeUxpc3QSNwoFaXRlbXMYsPDYASADKAsyHi5jbG91ZGZyb250LkNhY2hlUG'
    '9saWN5U3VtbWFyeVIFaXRlbXMSHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpu'
    'ZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj55dxfIAEoBVIIcX'
    'VhbnRpdHk=');

@$core.Deprecated('Use cachePolicyQueryStringsConfigDescriptor instead')
const CachePolicyQueryStringsConfig$json = {
  '1': 'CachePolicyQueryStringsConfig',
  '2': [
    {
      '1': 'querystringbehavior',
      '3': 235804641,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CachePolicyQueryStringBehavior',
      '10': 'querystringbehavior'
    },
    {
      '1': 'querystrings',
      '3': 478781456,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.QueryStringNames',
      '10': 'querystrings'
    },
  ],
};

/// Descriptor for `CachePolicyQueryStringsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicyQueryStringsConfigDescriptor = $convert.base64Decode(
    'Ch1DYWNoZVBvbGljeVF1ZXJ5U3RyaW5nc0NvbmZpZxJfChNxdWVyeXN0cmluZ2JlaGF2aW9yGO'
    'GvuHAgASgOMiouY2xvdWRmcm9udC5DYWNoZVBvbGljeVF1ZXJ5U3RyaW5nQmVoYXZpb3JSE3F1'
    'ZXJ5c3RyaW5nYmVoYXZpb3ISRAoMcXVlcnlzdHJpbmdzGJDApuQBIAEoCzIcLmNsb3VkZnJvbn'
    'QuUXVlcnlTdHJpbmdOYW1lc1IMcXVlcnlzdHJpbmdz');

@$core.Deprecated('Use cachePolicySummaryDescriptor instead')
const CachePolicySummary$json = {
  '1': 'CachePolicySummary',
  '2': [
    {
      '1': 'cachepolicy',
      '3': 439848032,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicy',
      '10': 'cachepolicy'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CachePolicyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `CachePolicySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachePolicySummaryDescriptor = $convert.base64Decode(
    'ChJDYWNoZVBvbGljeVN1bW1hcnkSPQoLY2FjaGVwb2xpY3kY4Jje0QEgASgLMhcuY2xvdWRmcm'
    '9udC5DYWNoZVBvbGljeVILY2FjaGVwb2xpY3kSMwoEdHlwZRjuoNeKASABKA4yGy5jbG91ZGZy'
    'b250LkNhY2hlUG9saWN5VHlwZVIEdHlwZQ==');

@$core.Deprecated('Use cachedMethodsDescriptor instead')
const CachedMethods$json = {
  '1': 'CachedMethods',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 14,
      '6': '.cloudfront.Method',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `CachedMethods`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cachedMethodsDescriptor = $convert.base64Decode(
    'Cg1DYWNoZWRNZXRob2RzEisKBWl0ZW1zGLDw2AEgAygOMhIuY2xvdWRmcm9udC5NZXRob2RSBW'
    'l0ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use cannotChangeImmutablePublicKeyFieldsDescriptor instead')
const CannotChangeImmutablePublicKeyFields$json = {
  '1': 'CannotChangeImmutablePublicKeyFields',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CannotChangeImmutablePublicKeyFields`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cannotChangeImmutablePublicKeyFieldsDescriptor =
    $convert.base64Decode(
        'CiRDYW5ub3RDaGFuZ2VJbW11dGFibGVQdWJsaWNLZXlGaWVsZHMSGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use cannotDeleteEntityWhileInUseDescriptor instead')
const CannotDeleteEntityWhileInUse$json = {
  '1': 'CannotDeleteEntityWhileInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CannotDeleteEntityWhileInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cannotDeleteEntityWhileInUseDescriptor =
    $convert.base64Decode(
        'ChxDYW5ub3REZWxldGVFbnRpdHlXaGlsZUluVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use cannotUpdateEntityWhileInUseDescriptor instead')
const CannotUpdateEntityWhileInUse$json = {
  '1': 'CannotUpdateEntityWhileInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CannotUpdateEntityWhileInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cannotUpdateEntityWhileInUseDescriptor =
    $convert.base64Decode(
        'ChxDYW5ub3RVcGRhdGVFbnRpdHlXaGlsZUluVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use certificateDescriptor instead')
const Certificate$json = {
  '1': 'Certificate',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `Certificate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateDescriptor =
    $convert.base64Decode('CgtDZXJ0aWZpY2F0ZRIUCgNhcm4YnZvtvwEgASgJUgNhcm4=');

@$core.Deprecated('Use cloudFrontOriginAccessIdentityDescriptor instead')
const CloudFrontOriginAccessIdentity$json = {
  '1': 'CloudFrontOriginAccessIdentity',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentityconfig',
      '3': 111945038,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentityConfig',
      '10': 'cloudfrontoriginaccessidentityconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 's3canonicaluserid',
      '3': 428876102,
      '4': 1,
      '5': 9,
      '10': 's3canonicaluserid'
    },
  ],
};

/// Descriptor for `CloudFrontOriginAccessIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudFrontOriginAccessIdentityDescriptor = $convert.base64Decode(
    'Ch5DbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHkShwEKJGNsb3VkZnJvbnRvcmlnaW5hY2'
    'Nlc3NpZGVudGl0eWNvbmZpZxjOyrA1IAEoCzIwLmNsb3VkZnJvbnQuQ2xvdWRGcm9udE9yaWdp'
    'bkFjY2Vzc0lkZW50aXR5Q29uZmlnUiRjbG91ZGZyb250b3JpZ2luYWNjZXNzaWRlbnRpdHljb2'
    '5maWcSEgoCaWQYgfKitwEgASgJUgJpZBIwChFzM2Nhbm9uaWNhbHVzZXJpZBjGwsDMASABKAlS'
    'EXMzY2Fub25pY2FsdXNlcmlk');

@$core.Deprecated(
    'Use cloudFrontOriginAccessIdentityAlreadyExistsDescriptor instead')
const CloudFrontOriginAccessIdentityAlreadyExists$json = {
  '1': 'CloudFrontOriginAccessIdentityAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudFrontOriginAccessIdentityAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    cloudFrontOriginAccessIdentityAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CitDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2'
        'UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use cloudFrontOriginAccessIdentityConfigDescriptor instead')
const CloudFrontOriginAccessIdentityConfig$json = {
  '1': 'CloudFrontOriginAccessIdentityConfig',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
  ],
};

/// Descriptor for `CloudFrontOriginAccessIdentityConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudFrontOriginAccessIdentityConfigDescriptor =
    $convert.base64Decode(
        'CiRDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlDb25maWcSKwoPY2FsbGVycmVmZXJlbm'
        'NlGJiZjUggASgJUg9jYWxsZXJyZWZlcmVuY2USHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1l'
        'bnQ=');

@$core.Deprecated('Use cloudFrontOriginAccessIdentityInUseDescriptor instead')
const CloudFrontOriginAccessIdentityInUse$json = {
  '1': 'CloudFrontOriginAccessIdentityInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudFrontOriginAccessIdentityInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudFrontOriginAccessIdentityInUseDescriptor =
    $convert.base64Decode(
        'CiNDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlJblVzZRIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use cloudFrontOriginAccessIdentityListDescriptor instead')
const CloudFrontOriginAccessIdentityList$json = {
  '1': 'CloudFrontOriginAccessIdentityList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentitySummary',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `CloudFrontOriginAccessIdentityList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudFrontOriginAccessIdentityListDescriptor = $convert.base64Decode(
    'CiJDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlMaXN0EiMKC2lzdHJ1bmNhdGVkGNqfuH'
    'MgASgIUgtpc3RydW5jYXRlZBJKCgVpdGVtcxiw8NgBIAMoCzIxLmNsb3VkZnJvbnQuQ2xvdWRG'
    'cm9udE9yaWdpbkFjY2Vzc0lkZW50aXR5U3VtbWFyeVIFaXRlbXMSGQoGbWFya2VyGLjdzSogAS'
    'gJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpuZXh0bWFya2Vy'
    'GKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use cloudFrontOriginAccessIdentitySummaryDescriptor instead')
const CloudFrontOriginAccessIdentitySummary$json = {
  '1': 'CloudFrontOriginAccessIdentitySummary',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 's3canonicaluserid',
      '3': 428876102,
      '4': 1,
      '5': 9,
      '10': 's3canonicaluserid'
    },
  ],
};

/// Descriptor for `CloudFrontOriginAccessIdentitySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudFrontOriginAccessIdentitySummaryDescriptor =
    $convert.base64Decode(
        'CiVDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlTdW1tYXJ5EhwKB2NvbW1lbnQY/7++wg'
        'EgASgJUgdjb21tZW50EhIKAmlkGIHyorcBIAEoCVICaWQSMAoRczNjYW5vbmljYWx1c2VyaWQY'
        'xsLAzAEgASgJUhFzM2Nhbm9uaWNhbHVzZXJpZA==');

@$core.Deprecated('Use conflictingAliasDescriptor instead')
const ConflictingAlias$json = {
  '1': 'ConflictingAlias',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {'1': 'alias', '3': 48362232, '4': 1, '5': 9, '10': 'alias'},
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
  ],
};

/// Descriptor for `ConflictingAlias`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictingAliasDescriptor = $convert.base64Decode(
    'ChBDb25mbGljdGluZ0FsaWFzEh8KCWFjY291bnRpZBjSwbkfIAEoCVIJYWNjb3VudGlkEhcKBW'
    'FsaWFzGPjlhxcgASgJUgVhbGlhcxIpCg5kaXN0cmlidXRpb25pZBjnsftDIAEoCVIOZGlzdHJp'
    'YnV0aW9uaWQ=');

@$core.Deprecated('Use conflictingAliasesListDescriptor instead')
const ConflictingAliasesList$json = {
  '1': 'ConflictingAliasesList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ConflictingAlias',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ConflictingAliasesList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictingAliasesListDescriptor = $convert.base64Decode(
    'ChZDb25mbGljdGluZ0FsaWFzZXNMaXN0EjUKBWl0ZW1zGLDw2AEgAygLMhwuY2xvdWRmcm9udC'
    '5Db25mbGljdGluZ0FsaWFzUgVpdGVtcxIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z'
    'EiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2VyEh0KCHF1YW50aXR5GPnl3F8gAS'
    'gFUghxdWFudGl0eQ==');

@$core.Deprecated('Use connectionFunctionAssociationDescriptor instead')
const ConnectionFunctionAssociation$json = {
  '1': 'ConnectionFunctionAssociation',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `ConnectionFunctionAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionFunctionAssociationDescriptor =
    $convert.base64Decode(
        'Ch1Db25uZWN0aW9uRnVuY3Rpb25Bc3NvY2lhdGlvbhISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use connectionFunctionSummaryDescriptor instead')
const ConnectionFunctionSummary$json = {
  '1': 'ConnectionFunctionSummary',
  '2': [
    {
      '1': 'connectionfunctionarn',
      '3': 69152757,
      '4': 1,
      '5': 9,
      '10': 'connectionfunctionarn'
    },
    {
      '1': 'connectionfunctionconfig',
      '3': 349557744,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionConfig',
      '10': 'connectionfunctionconfig'
    },
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `ConnectionFunctionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionFunctionSummaryDescriptor = $convert.base64Decode(
    'ChlDb25uZWN0aW9uRnVuY3Rpb25TdW1tYXJ5EjcKFWNvbm5lY3Rpb25mdW5jdGlvbmFybhj13/'
    'wgIAEoCVIVY29ubmVjdGlvbmZ1bmN0aW9uYXJuEloKGGNvbm5lY3Rpb25mdW5jdGlvbmNvbmZp'
    'Zxjwp9emASABKAsyGi5jbG91ZGZyb250LkZ1bmN0aW9uQ29uZmlnUhhjb25uZWN0aW9uZnVuY3'
    'Rpb25jb25maWcSIwoLY3JlYXRlZHRpbWUY8+vzOSABKAlSC2NyZWF0ZWR0aW1lEhIKAmlkGIHy'
    'orcBIAEoCVICaWQSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdG'
    'ltZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjIKBXN0YWdlGM6X2HAgASgOMhkuY2xvdWRmcm9u'
    'dC5GdW5jdGlvblN0YWdlUgVzdGFnZRIZCgZzdGF0dXMYkOT7AiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use connectionFunctionTestResultDescriptor instead')
const ConnectionFunctionTestResult$json = {
  '1': 'ConnectionFunctionTestResult',
  '2': [
    {
      '1': 'computeutilization',
      '3': 247332359,
      '4': 1,
      '5': 9,
      '10': 'computeutilization'
    },
    {
      '1': 'connectionfunctionerrormessage',
      '3': 360436321,
      '4': 1,
      '5': 9,
      '10': 'connectionfunctionerrormessage'
    },
    {
      '1': 'connectionfunctionexecutionlogs',
      '3': 194700351,
      '4': 3,
      '5': 9,
      '10': 'connectionfunctionexecutionlogs'
    },
    {
      '1': 'connectionfunctionoutput',
      '3': 235201677,
      '4': 1,
      '5': 9,
      '10': 'connectionfunctionoutput'
    },
    {
      '1': 'connectionfunctionsummary',
      '3': 62528396,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionSummary',
      '10': 'connectionfunctionsummary'
    },
  ],
};

/// Descriptor for `ConnectionFunctionTestResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionFunctionTestResultDescriptor = $convert.base64Decode(
    'ChxDb25uZWN0aW9uRnVuY3Rpb25UZXN0UmVzdWx0EjEKEmNvbXB1dGV1dGlsaXphdGlvbhiH/P'
    'd1IAEoCVISY29tcHV0ZXV0aWxpemF0aW9uEkoKHmNvbm5lY3Rpb25mdW5jdGlvbmVycm9ybWVz'
    'c2FnZRjhpO+rASABKAlSHmNvbm5lY3Rpb25mdW5jdGlvbmVycm9ybWVzc2FnZRJLCh9jb25uZW'
    'N0aW9uZnVuY3Rpb25leGVjdXRpb25sb2dzGL/I61wgAygJUh9jb25uZWN0aW9uZnVuY3Rpb25l'
    'eGVjdXRpb25sb2dzEj0KGGNvbm5lY3Rpb25mdW5jdGlvbm91dHB1dBiNyZNwIAEoCVIYY29ubm'
    'VjdGlvbmZ1bmN0aW9ub3V0cHV0EmYKGWNvbm5lY3Rpb25mdW5jdGlvbnN1bW1hcnkYjLfoHSAB'
    'KAsyJS5jbG91ZGZyb250LkNvbm5lY3Rpb25GdW5jdGlvblN1bW1hcnlSGWNvbm5lY3Rpb25mdW'
    '5jdGlvbnN1bW1hcnk=');

@$core.Deprecated('Use connectionGroupDescriptor instead')
const ConnectionGroup$json = {
  '1': 'ConnectionGroup',
  '2': [
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ipv6enabled', '3': 324073322, '4': 1, '5': 8, '10': 'ipv6enabled'},
    {'1': 'isdefault', '3': 101743631, '4': 1, '5': 8, '10': 'isdefault'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'routingendpoint',
      '3': 409864193,
      '4': 1,
      '5': 9,
      '10': 'routingendpoint'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ConnectionGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionGroupDescriptor = $convert.base64Decode(
    'Cg9Db25uZWN0aW9uR3JvdXASLAoPYW55Y2FzdGlwbGlzdGlkGIOs+M0BIAEoCVIPYW55Y2FzdG'
    'lwbGlzdGlkEhQKA2Fybhidm+2/ASABKAlSA2FybhIjCgtjcmVhdGVkdGltZRjz6/M5IAEoCVIL'
    'Y3JlYXRlZHRpbWUSHAoHZW5hYmxlZBi/yJvkASABKAhSB2VuYWJsZWQSEgoCaWQYgfKitwEgAS'
    'gJUgJpZBIkCgtpcHY2ZW5hYmxlZBjq7sOaASABKAhSC2lwdjZlbmFibGVkEh8KCWlzZGVmYXVs'
    'dBiP+MEwIAEoCFIJaXNkZWZhdWx0Ei0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3'
    'Rtb2RpZmllZHRpbWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIsCg9yb3V0aW5nZW5kcG9pbnQY'
    'gZC4wwEgASgJUg9yb3V0aW5nZW5kcG9pbnQSGQoGc3RhdHVzGJDk+wIgASgJUgZzdGF0dXMSKA'
    'oEdGFncxjBwfa1ASABKAsyEC5jbG91ZGZyb250LlRhZ3NSBHRhZ3M=');

@$core.Deprecated('Use connectionGroupAssociationFilterDescriptor instead')
const ConnectionGroupAssociationFilter$json = {
  '1': 'ConnectionGroupAssociationFilter',
  '2': [
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
  ],
};

/// Descriptor for `ConnectionGroupAssociationFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionGroupAssociationFilterDescriptor =
    $convert.base64Decode(
        'CiBDb25uZWN0aW9uR3JvdXBBc3NvY2lhdGlvbkZpbHRlchIsCg9hbnljYXN0aXBsaXN0aWQYg6'
        'z4zQEgASgJUg9hbnljYXN0aXBsaXN0aWQ=');

@$core.Deprecated('Use connectionGroupSummaryDescriptor instead')
const ConnectionGroupSummary$json = {
  '1': 'ConnectionGroupSummary',
  '2': [
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'isdefault', '3': 101743631, '4': 1, '5': 8, '10': 'isdefault'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'routingendpoint',
      '3': 409864193,
      '4': 1,
      '5': 9,
      '10': 'routingendpoint'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `ConnectionGroupSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionGroupSummaryDescriptor = $convert.base64Decode(
    'ChZDb25uZWN0aW9uR3JvdXBTdW1tYXJ5EiwKD2FueWNhc3RpcGxpc3RpZBiDrPjNASABKAlSD2'
    'FueWNhc3RpcGxpc3RpZBIUCgNhcm4YnZvtvwEgASgJUgNhcm4SIwoLY3JlYXRlZHRpbWUY8+vz'
    'OSABKAlSC2NyZWF0ZWR0aW1lEhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEhwKB2VuYWJsZWQYv8'
    'ib5AEgASgIUgdlbmFibGVkEhIKAmlkGIHyorcBIAEoCVICaWQSHwoJaXNkZWZhdWx0GI/4wTAg'
    'ASgIUglpc2RlZmF1bHQSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaW'
    'VkdGltZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiwKD3JvdXRpbmdlbmRwb2ludBiBkLjDASAB'
    'KAlSD3JvdXRpbmdlbmRwb2ludBIZCgZzdGF0dXMYkOT7AiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use contentTypeProfileDescriptor instead')
const ContentTypeProfile$json = {
  '1': 'ContentTypeProfile',
  '2': [
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {
      '1': 'format',
      '3': 531693427,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.Format',
      '10': 'format'
    },
    {'1': 'profileid', '3': 407138548, '4': 1, '5': 9, '10': 'profileid'},
  ],
};

/// Descriptor for `ContentTypeProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contentTypeProfileDescriptor = $convert.base64Decode(
    'ChJDb250ZW50VHlwZVByb2ZpbGUSJAoLY29udGVudHR5cGUYk9XongEgASgJUgtjb250ZW50dH'
    'lwZRIuCgZmb3JtYXQY8/7D/QEgASgOMhIuY2xvdWRmcm9udC5Gb3JtYXRSBmZvcm1hdBIgCglw'
    'cm9maWxlaWQY9OGRwgEgASgJUglwcm9maWxlaWQ=');

@$core.Deprecated('Use contentTypeProfileConfigDescriptor instead')
const ContentTypeProfileConfig$json = {
  '1': 'ContentTypeProfileConfig',
  '2': [
    {
      '1': 'contenttypeprofiles',
      '3': 110010517,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContentTypeProfiles',
      '10': 'contenttypeprofiles'
    },
    {
      '1': 'forwardwhencontenttypeisunknown',
      '3': 472006304,
      '4': 1,
      '5': 8,
      '10': 'forwardwhencontenttypeisunknown'
    },
  ],
};

/// Descriptor for `ContentTypeProfileConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contentTypeProfileConfigDescriptor = $convert.base64Decode(
    'ChhDb250ZW50VHlwZVByb2ZpbGVDb25maWcSVAoTY29udGVudHR5cGVwcm9maWxlcxiVwbo0IA'
    'EoCzIfLmNsb3VkZnJvbnQuQ29udGVudFR5cGVQcm9maWxlc1ITY29udGVudHR5cGVwcm9maWxl'
    'cxJMCh9mb3J3YXJkd2hlbmNvbnRlbnR0eXBlaXN1bmtub3duGKD9iOEBIAEoCFIfZm9yd2FyZH'
    'doZW5jb250ZW50dHlwZWlzdW5rbm93bg==');

@$core.Deprecated('Use contentTypeProfilesDescriptor instead')
const ContentTypeProfiles$json = {
  '1': 'ContentTypeProfiles',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ContentTypeProfile',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ContentTypeProfiles`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contentTypeProfilesDescriptor = $convert.base64Decode(
    'ChNDb250ZW50VHlwZVByb2ZpbGVzEjcKBWl0ZW1zGLDw2AEgAygLMh4uY2xvdWRmcm9udC5Db2'
    '50ZW50VHlwZVByb2ZpbGVSBWl0ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use continuousDeploymentPolicyDescriptor instead')
const ContinuousDeploymentPolicy$json = {
  '1': 'ContinuousDeploymentPolicy',
  '2': [
    {
      '1': 'continuousdeploymentpolicyconfig',
      '3': 161949042,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicyConfig',
      '10': 'continuousdeploymentpolicyconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `ContinuousDeploymentPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentPolicyDescriptor = $convert.base64Decode(
    'ChpDb250aW51b3VzRGVwbG95bWVudFBvbGljeRJ7CiBjb250aW51b3VzZGVwbG95bWVudHBvbG'
    'ljeWNvbmZpZxjyypxNIAEoCzIsLmNsb3VkZnJvbnQuQ29udGludW91c0RlcGxveW1lbnRQb2xp'
    'Y3lDb25maWdSIGNvbnRpbnVvdXNkZXBsb3ltZW50cG9saWN5Y29uZmlnEhIKAmlkGIHyorcBIA'
    'EoCVICaWQSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZQ==');

@$core
    .Deprecated('Use continuousDeploymentPolicyAlreadyExistsDescriptor instead')
const ContinuousDeploymentPolicyAlreadyExists$json = {
  '1': 'ContinuousDeploymentPolicyAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ContinuousDeploymentPolicyAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentPolicyAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CidDb250aW51b3VzRGVwbG95bWVudFBvbGljeUFscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use continuousDeploymentPolicyConfigDescriptor instead')
const ContinuousDeploymentPolicyConfig$json = {
  '1': 'ContinuousDeploymentPolicyConfig',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'stagingdistributiondnsnames',
      '3': 313322754,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StagingDistributionDnsNames',
      '10': 'stagingdistributiondnsnames'
    },
    {
      '1': 'trafficconfig',
      '3': 483116577,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrafficConfig',
      '10': 'trafficconfig'
    },
  ],
};

/// Descriptor for `ContinuousDeploymentPolicyConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentPolicyConfigDescriptor = $convert.base64Decode(
    'CiBDb250aW51b3VzRGVwbG95bWVudFBvbGljeUNvbmZpZxIcCgdlbmFibGVkGL/Im+QBIAEoCF'
    'IHZW5hYmxlZBJtChtzdGFnaW5nZGlzdHJpYnV0aW9uZG5zbmFtZXMYgtqzlQEgASgLMicuY2xv'
    'dWRmcm9udC5TdGFnaW5nRGlzdHJpYnV0aW9uRG5zTmFtZXNSG3N0YWdpbmdkaXN0cmlidXRpb2'
    '5kbnNuYW1lcxJDCg10cmFmZmljY29uZmlnGKGMr+YBIAEoCzIZLmNsb3VkZnJvbnQuVHJhZmZp'
    'Y0NvbmZpZ1INdHJhZmZpY2NvbmZpZw==');

@$core.Deprecated('Use continuousDeploymentPolicyInUseDescriptor instead')
const ContinuousDeploymentPolicyInUse$json = {
  '1': 'ContinuousDeploymentPolicyInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ContinuousDeploymentPolicyInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentPolicyInUseDescriptor =
    $convert.base64Decode(
        'Ch9Db250aW51b3VzRGVwbG95bWVudFBvbGljeUluVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use continuousDeploymentPolicyListDescriptor instead')
const ContinuousDeploymentPolicyList$json = {
  '1': 'ContinuousDeploymentPolicyList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicySummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ContinuousDeploymentPolicyList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentPolicyListDescriptor =
    $convert.base64Decode(
        'Ch5Db250aW51b3VzRGVwbG95bWVudFBvbGljeUxpc3QSRgoFaXRlbXMYsPDYASADKAsyLS5jbG'
        '91ZGZyb250LkNvbnRpbnVvdXNEZXBsb3ltZW50UG9saWN5U3VtbWFyeVIFaXRlbXMSHgoIbWF4'
        'aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG'
        '1hcmtlchIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use continuousDeploymentPolicySummaryDescriptor instead')
const ContinuousDeploymentPolicySummary$json = {
  '1': 'ContinuousDeploymentPolicySummary',
  '2': [
    {
      '1': 'continuousdeploymentpolicy',
      '3': 36616788,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicy',
      '10': 'continuousdeploymentpolicy'
    },
  ],
};

/// Descriptor for `ContinuousDeploymentPolicySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentPolicySummaryDescriptor =
    $convert.base64Decode(
        'CiFDb250aW51b3VzRGVwbG95bWVudFBvbGljeVN1bW1hcnkSaQoaY29udGludW91c2RlcGxveW'
        '1lbnRwb2xpY3kY1PS6ESABKAsyJi5jbG91ZGZyb250LkNvbnRpbnVvdXNEZXBsb3ltZW50UG9s'
        'aWN5Uhpjb250aW51b3VzZGVwbG95bWVudHBvbGljeQ==');

@$core
    .Deprecated('Use continuousDeploymentSingleHeaderConfigDescriptor instead')
const ContinuousDeploymentSingleHeaderConfig$json = {
  '1': 'ContinuousDeploymentSingleHeaderConfig',
  '2': [
    {'1': 'header', '3': 290429313, '4': 1, '5': 9, '10': 'header'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ContinuousDeploymentSingleHeaderConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentSingleHeaderConfigDescriptor =
    $convert.base64Decode(
        'CiZDb250aW51b3VzRGVwbG95bWVudFNpbmdsZUhlYWRlckNvbmZpZxIaCgZoZWFkZXIYgbO+ig'
        'EgASgJUgZoZWFkZXISGAoFdmFsdWUY6/KfigEgASgJUgV2YWx1ZQ==');

@$core
    .Deprecated('Use continuousDeploymentSingleWeightConfigDescriptor instead')
const ContinuousDeploymentSingleWeightConfig$json = {
  '1': 'ContinuousDeploymentSingleWeightConfig',
  '2': [
    {
      '1': 'sessionstickinessconfig',
      '3': 870690,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.SessionStickinessConfig',
      '10': 'sessionstickinessconfig'
    },
    {'1': 'weight', '3': 422581466, '4': 1, '5': 2, '10': 'weight'},
  ],
};

/// Descriptor for `ContinuousDeploymentSingleWeightConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousDeploymentSingleWeightConfigDescriptor =
    $convert.base64Decode(
        'CiZDb250aW51b3VzRGVwbG95bWVudFNpbmdsZVdlaWdodENvbmZpZxJfChdzZXNzaW9uc3RpY2'
        'tpbmVzc2NvbmZpZxiikjUgASgLMiMuY2xvdWRmcm9udC5TZXNzaW9uU3RpY2tpbmVzc0NvbmZp'
        'Z1IXc2Vzc2lvbnN0aWNraW5lc3Njb25maWcSGgoGd2VpZ2h0GNqpwMkBIAEoAlIGd2VpZ2h0');

@$core.Deprecated('Use cookieNamesDescriptor instead')
const CookieNames$json = {
  '1': 'CookieNames',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `CookieNames`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cookieNamesDescriptor = $convert.base64Decode(
    'CgtDb29raWVOYW1lcxIXCgVpdGVtcxiw8NgBIAMoCVIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXy'
    'ABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use cookiePreferenceDescriptor instead')
const CookiePreference$json = {
  '1': 'CookiePreference',
  '2': [
    {
      '1': 'forward',
      '3': 84444091,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ItemSelection',
      '10': 'forward'
    },
    {
      '1': 'whitelistednames',
      '3': 172336098,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CookieNames',
      '10': 'whitelistednames'
    },
  ],
};

/// Descriptor for `CookiePreference`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cookiePreferenceDescriptor = $convert.base64Decode(
    'ChBDb29raWVQcmVmZXJlbmNlEjYKB2ZvcndhcmQYu4eiKCABKA4yGS5jbG91ZGZyb250Lkl0ZW'
    '1TZWxlY3Rpb25SB2ZvcndhcmQSRgoQd2hpdGVsaXN0ZWRuYW1lcxjix5ZSIAEoCzIXLmNsb3Vk'
    'ZnJvbnQuQ29va2llTmFtZXNSEHdoaXRlbGlzdGVkbmFtZXM=');

@$core.Deprecated('Use copyDistributionRequestDescriptor instead')
const CopyDistributionRequest$json = {
  '1': 'CopyDistributionRequest',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'primarydistributionid',
      '3': 430692915,
      '4': 1,
      '5': 9,
      '10': 'primarydistributionid'
    },
    {'1': 'staging', '3': 193058759, '4': 1, '5': 8, '10': 'staging'},
  ],
};

/// Descriptor for `CopyDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyDistributionRequestDescriptor = $convert.base64Decode(
    'ChdDb3B5RGlzdHJpYnV0aW9uUmVxdWVzdBIrCg9jYWxsZXJyZWZlcmVuY2UYmJmNSCABKAlSD2'
    'NhbGxlcnJlZmVyZW5jZRIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIbCgdpZm1hdGNo'
    'GNCWtywgASgJUgdpZm1hdGNoEjgKFXByaW1hcnlkaXN0cmlidXRpb25pZBiztK/NASABKAlSFX'
    'ByaW1hcnlkaXN0cmlidXRpb25pZBIbCgdzdGFnaW5nGMevh1wgASgIUgdzdGFnaW5n');

@$core.Deprecated('Use copyDistributionResultDescriptor instead')
const CopyDistributionResult$json = {
  '1': 'CopyDistributionResult',
  '2': [
    {
      '1': 'distribution',
      '3': 105183308,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Distribution',
      '8': {},
      '10': 'distribution'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CopyDistributionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyDistributionResultDescriptor = $convert.base64Decode(
    'ChZDb3B5RGlzdHJpYnV0aW9uUmVzdWx0EkUKDGRpc3RyaWJ1dGlvbhjM8JMyIAEoCzIYLmNsb3'
    'VkZnJvbnQuRGlzdHJpYnV0aW9uQgSItRgBUgxkaXN0cmlidXRpb24SFgoEZXRhZxiB37OVASAB'
    'KAlSBGV0YWcSHgoIbG9jYXRpb24Yx5uC3gEgASgJUghsb2NhdGlvbg==');

@$core.Deprecated('Use createAnycastIpListRequestDescriptor instead')
const CreateAnycastIpListRequest$json = {
  '1': 'CreateAnycastIpListRequest',
  '2': [
    {
      '1': 'ipaddresstype',
      '3': 459110693,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.IpAddressType',
      '10': 'ipaddresstype'
    },
    {'1': 'ipcount', '3': 475138532, '4': 1, '5': 5, '10': 'ipcount'},
    {
      '1': 'ipamcidrconfigs',
      '3': 314539184,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.IpamCidrConfig',
      '10': 'ipamcidrconfigs'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateAnycastIpListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAnycastIpListRequestDescriptor = $convert.base64Decode(
    'ChpDcmVhdGVBbnljYXN0SXBMaXN0UmVxdWVzdBJDCg1pcGFkZHJlc3N0eXBlGKXy9doBIAEoDj'
    'IZLmNsb3VkZnJvbnQuSXBBZGRyZXNzVHlwZVINaXBhZGRyZXNzdHlwZRIcCgdpcGNvdW50GOST'
    'yOIBIAEoBVIHaXBjb3VudBJICg9pcGFtY2lkcmNvbmZpZ3MYsPn9lQEgAygLMhouY2xvdWRmcm'
    '9udC5JcGFtQ2lkckNvbmZpZ1IPaXBhbWNpZHJjb25maWdzEhUKBG5hbWUYh+aBfyABKAlSBG5h'
    'bWUSKAoEdGFncxjBwfa1ASABKAsyEC5jbG91ZGZyb250LlRhZ3NSBHRhZ3M=');

@$core.Deprecated('Use createAnycastIpListResultDescriptor instead')
const CreateAnycastIpListResult$json = {
  '1': 'CreateAnycastIpListResult',
  '2': [
    {
      '1': 'anycastiplist',
      '3': 190550768,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.AnycastIpList',
      '8': {},
      '10': 'anycastiplist'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `CreateAnycastIpListResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAnycastIpListResultDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVBbnljYXN0SXBMaXN0UmVzdWx0EkgKDWFueWNhc3RpcGxpc3QY8KXuWiABKAsyGS'
    '5jbG91ZGZyb250LkFueWNhc3RJcExpc3RCBIi1GAFSDWFueWNhc3RpcGxpc3QSFgoEZXRhZxiB'
    '37OVASABKAlSBGV0YWc=');

@$core.Deprecated('Use createCachePolicyRequestDescriptor instead')
const CreateCachePolicyRequest$json = {
  '1': 'CreateCachePolicyRequest',
  '2': [
    {
      '1': 'cachepolicyconfig',
      '3': 407094126,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyConfig',
      '8': {},
      '10': 'cachepolicyconfig'
    },
  ],
};

/// Descriptor for `CreateCachePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCachePolicyRequestDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVDYWNoZVBvbGljeVJlcXVlc3QSVQoRY2FjaGVwb2xpY3ljb25maWcY7oaPwgEgAS'
    'gLMh0uY2xvdWRmcm9udC5DYWNoZVBvbGljeUNvbmZpZ0IEiLUYAVIRY2FjaGVwb2xpY3ljb25m'
    'aWc=');

@$core.Deprecated('Use createCachePolicyResultDescriptor instead')
const CreateCachePolicyResult$json = {
  '1': 'CreateCachePolicyResult',
  '2': [
    {
      '1': 'cachepolicy',
      '3': 439848032,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicy',
      '8': {},
      '10': 'cachepolicy'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateCachePolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCachePolicyResultDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVDYWNoZVBvbGljeVJlc3VsdBJDCgtjYWNoZXBvbGljeRjgmN7RASABKAsyFy5jbG'
    '91ZGZyb250LkNhY2hlUG9saWN5QgSItRgBUgtjYWNoZXBvbGljeRIWCgRldGFnGIHfs5UBIAEo'
    'CVIEZXRhZxIeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated(
    'Use createCloudFrontOriginAccessIdentityRequestDescriptor instead')
const CreateCloudFrontOriginAccessIdentityRequest$json = {
  '1': 'CreateCloudFrontOriginAccessIdentityRequest',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentityconfig',
      '3': 111945038,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentityConfig',
      '8': {},
      '10': 'cloudfrontoriginaccessidentityconfig'
    },
  ],
};

/// Descriptor for `CreateCloudFrontOriginAccessIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createCloudFrontOriginAccessIdentityRequestDescriptor =
    $convert.base64Decode(
        'CitDcmVhdGVDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlSZXF1ZXN0Eo0BCiRjbG91ZG'
        'Zyb250b3JpZ2luYWNjZXNzaWRlbnRpdHljb25maWcYzsqwNSABKAsyMC5jbG91ZGZyb250LkNs'
        'b3VkRnJvbnRPcmlnaW5BY2Nlc3NJZGVudGl0eUNvbmZpZ0IEiLUYAVIkY2xvdWRmcm9udG9yaW'
        'dpbmFjY2Vzc2lkZW50aXR5Y29uZmln');

@$core.Deprecated(
    'Use createCloudFrontOriginAccessIdentityResultDescriptor instead')
const CreateCloudFrontOriginAccessIdentityResult$json = {
  '1': 'CreateCloudFrontOriginAccessIdentityResult',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentity',
      '3': 109497984,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentity',
      '8': {},
      '10': 'cloudfrontoriginaccessidentity'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateCloudFrontOriginAccessIdentityResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createCloudFrontOriginAccessIdentityResultDescriptor =
    $convert.base64Decode(
        'CipDcmVhdGVDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlSZXN1bHQSewoeY2xvdWRmcm'
        '9udG9yaWdpbmFjY2Vzc2lkZW50aXR5GICdmzQgASgLMiouY2xvdWRmcm9udC5DbG91ZEZyb250'
        'T3JpZ2luQWNjZXNzSWRlbnRpdHlCBIi1GAFSHmNsb3VkZnJvbnRvcmlnaW5hY2Nlc3NpZGVudG'
        'l0eRIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxIeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0'
        'aW9u');

@$core.Deprecated('Use createConnectionFunctionRequestDescriptor instead')
const CreateConnectionFunctionRequest$json = {
  '1': 'CreateConnectionFunctionRequest',
  '2': [
    {
      '1': 'connectionfunctioncode',
      '3': 502949501,
      '4': 1,
      '5': 12,
      '10': 'connectionfunctioncode'
    },
    {
      '1': 'connectionfunctionconfig',
      '3': 349557744,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionConfig',
      '10': 'connectionfunctionconfig'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionFunctionRequestDescriptor = $convert.base64Decode(
    'Ch9DcmVhdGVDb25uZWN0aW9uRnVuY3Rpb25SZXF1ZXN0EjoKFmNvbm5lY3Rpb25mdW5jdGlvbm'
    'NvZGUY/czp7wEgASgMUhZjb25uZWN0aW9uZnVuY3Rpb25jb2RlEloKGGNvbm5lY3Rpb25mdW5j'
    'dGlvbmNvbmZpZxjwp9emASABKAsyGi5jbG91ZGZyb250LkZ1bmN0aW9uQ29uZmlnUhhjb25uZW'
    'N0aW9uZnVuY3Rpb25jb25maWcSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIoCgR0YWdzGMHB9rUB'
    'IAEoCzIQLmNsb3VkZnJvbnQuVGFnc1IEdGFncw==');

@$core.Deprecated('Use createConnectionFunctionResultDescriptor instead')
const CreateConnectionFunctionResult$json = {
  '1': 'CreateConnectionFunctionResult',
  '2': [
    {
      '1': 'connectionfunctionsummary',
      '3': 62528396,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionSummary',
      '8': {},
      '10': 'connectionfunctionsummary'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateConnectionFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionFunctionResultDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVDb25uZWN0aW9uRnVuY3Rpb25SZXN1bHQSbAoZY29ubmVjdGlvbmZ1bmN0aW9uc3'
        'VtbWFyeRiMt+gdIAEoCzIlLmNsb3VkZnJvbnQuQ29ubmVjdGlvbkZ1bmN0aW9uU3VtbWFyeUIE'
        'iLUYAVIZY29ubmVjdGlvbmZ1bmN0aW9uc3VtbWFyeRIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZx'
        'IeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated('Use createConnectionGroupRequestDescriptor instead')
const CreateConnectionGroupRequest$json = {
  '1': 'CreateConnectionGroupRequest',
  '2': [
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'ipv6enabled', '3': 324073322, '4': 1, '5': 8, '10': 'ipv6enabled'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateConnectionGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionGroupRequestDescriptor = $convert.base64Decode(
    'ChxDcmVhdGVDb25uZWN0aW9uR3JvdXBSZXF1ZXN0EiwKD2FueWNhc3RpcGxpc3RpZBiDrPjNAS'
    'ABKAlSD2FueWNhc3RpcGxpc3RpZBIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIkCgtp'
    'cHY2ZW5hYmxlZBjq7sOaASABKAhSC2lwdjZlbmFibGVkEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'USKAoEdGFncxjBwfa1ASABKAsyEC5jbG91ZGZyb250LlRhZ3NSBHRhZ3M=');

@$core.Deprecated('Use createConnectionGroupResultDescriptor instead')
const CreateConnectionGroupResult$json = {
  '1': 'CreateConnectionGroupResult',
  '2': [
    {
      '1': 'connectiongroup',
      '3': 517217105,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionGroup',
      '8': {},
      '10': 'connectiongroup'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `CreateConnectionGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConnectionGroupResultDescriptor =
    $convert.base64Decode(
        'ChtDcmVhdGVDb25uZWN0aW9uR3JvdXBSZXN1bHQSTwoPY29ubmVjdGlvbmdyb3VwGNG20PYBIA'
        'EoCzIbLmNsb3VkZnJvbnQuQ29ubmVjdGlvbkdyb3VwQgSItRgBUg9jb25uZWN0aW9uZ3JvdXAS'
        'FgoEZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core
    .Deprecated('Use createContinuousDeploymentPolicyRequestDescriptor instead')
const CreateContinuousDeploymentPolicyRequest$json = {
  '1': 'CreateContinuousDeploymentPolicyRequest',
  '2': [
    {
      '1': 'continuousdeploymentpolicyconfig',
      '3': 161949042,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicyConfig',
      '8': {},
      '10': 'continuousdeploymentpolicyconfig'
    },
  ],
};

/// Descriptor for `CreateContinuousDeploymentPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContinuousDeploymentPolicyRequestDescriptor =
    $convert.base64Decode(
        'CidDcmVhdGVDb250aW51b3VzRGVwbG95bWVudFBvbGljeVJlcXVlc3QSgQEKIGNvbnRpbnVvdX'
        'NkZXBsb3ltZW50cG9saWN5Y29uZmlnGPLKnE0gASgLMiwuY2xvdWRmcm9udC5Db250aW51b3Vz'
        'RGVwbG95bWVudFBvbGljeUNvbmZpZ0IEiLUYAVIgY29udGludW91c2RlcGxveW1lbnRwb2xpY3'
        'ljb25maWc=');

@$core
    .Deprecated('Use createContinuousDeploymentPolicyResultDescriptor instead')
const CreateContinuousDeploymentPolicyResult$json = {
  '1': 'CreateContinuousDeploymentPolicyResult',
  '2': [
    {
      '1': 'continuousdeploymentpolicy',
      '3': 36616788,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicy',
      '8': {},
      '10': 'continuousdeploymentpolicy'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateContinuousDeploymentPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContinuousDeploymentPolicyResultDescriptor =
    $convert.base64Decode(
        'CiZDcmVhdGVDb250aW51b3VzRGVwbG95bWVudFBvbGljeVJlc3VsdBJvChpjb250aW51b3VzZG'
        'VwbG95bWVudHBvbGljeRjU9LoRIAEoCzImLmNsb3VkZnJvbnQuQ29udGludW91c0RlcGxveW1l'
        'bnRQb2xpY3lCBIi1GAFSGmNvbnRpbnVvdXNkZXBsb3ltZW50cG9saWN5EhYKBGV0YWcYgd+zlQ'
        'EgASgJUgRldGFnEh4KCGxvY2F0aW9uGMebgt4BIAEoCVIIbG9jYXRpb24=');

@$core.Deprecated('Use createDistributionRequestDescriptor instead')
const CreateDistributionRequest$json = {
  '1': 'CreateDistributionRequest',
  '2': [
    {
      '1': 'distributionconfig',
      '3': 528940762,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionConfig',
      '8': {},
      '10': 'distributionconfig'
    },
  ],
};

/// Descriptor for `CreateDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDistributionRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVEaXN0cmlidXRpb25SZXF1ZXN0ElgKEmRpc3RyaWJ1dGlvbmNvbmZpZxja/Zv8AS'
    'ABKAsyHi5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkNvbmZpZ0IEiLUYAVISZGlzdHJpYnV0aW9u'
    'Y29uZmln');

@$core.Deprecated('Use createDistributionResultDescriptor instead')
const CreateDistributionResult$json = {
  '1': 'CreateDistributionResult',
  '2': [
    {
      '1': 'distribution',
      '3': 105183308,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Distribution',
      '8': {},
      '10': 'distribution'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateDistributionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDistributionResultDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVEaXN0cmlidXRpb25SZXN1bHQSRQoMZGlzdHJpYnV0aW9uGMzwkzIgASgLMhguY2'
    'xvdWRmcm9udC5EaXN0cmlidXRpb25CBIi1GAFSDGRpc3RyaWJ1dGlvbhIWCgRldGFnGIHfs5UB'
    'IAEoCVIEZXRhZxIeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated('Use createDistributionTenantRequestDescriptor instead')
const CreateDistributionTenantRequest$json = {
  '1': 'CreateDistributionTenantRequest',
  '2': [
    {
      '1': 'connectiongroupid',
      '3': 169532206,
      '4': 1,
      '5': 9,
      '10': 'connectiongroupid'
    },
    {
      '1': 'customizations',
      '3': 70755200,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Customizations',
      '10': 'customizations'
    },
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'domains',
      '3': 149701959,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DomainItem',
      '10': 'domains'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'managedcertificaterequest',
      '3': 310259927,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ManagedCertificateRequest',
      '10': 'managedcertificaterequest'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.Parameter',
      '10': 'parameters'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateDistributionTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDistributionTenantRequestDescriptor = $convert.base64Decode(
    'Ch9DcmVhdGVEaXN0cmlidXRpb25UZW5hbnRSZXF1ZXN0Ei8KEWNvbm5lY3Rpb25ncm91cGlkGK'
    '6261AgASgJUhFjb25uZWN0aW9uZ3JvdXBpZBJFCg5jdXN0b21pemF0aW9ucxiAx94hIAEoCzIa'
    'LmNsb3VkZnJvbnQuQ3VzdG9taXphdGlvbnNSDmN1c3RvbWl6YXRpb25zEikKDmRpc3RyaWJ1dG'
    'lvbmlkGOex+0MgASgJUg5kaXN0cmlidXRpb25pZBIzCgdkb21haW5zGMeKsUcgAygLMhYuY2xv'
    'dWRmcm9udC5Eb21haW5JdGVtUgdkb21haW5zEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibG'
    'VkEmcKGW1hbmFnZWRjZXJ0aWZpY2F0ZXJlcXVlc3QY1+H4kwEgASgLMiUuY2xvdWRmcm9udC5N'
    'YW5hZ2VkQ2VydGlmaWNhdGVSZXF1ZXN0UhltYW5hZ2VkY2VydGlmaWNhdGVyZXF1ZXN0EhUKBG'
    '5hbWUYh+aBfyABKAlSBG5hbWUSOQoKcGFyYW1ldGVycxj6p/7rASADKAsyFS5jbG91ZGZyb250'
    'LlBhcmFtZXRlclIKcGFyYW1ldGVycxIoCgR0YWdzGMHB9rUBIAEoCzIQLmNsb3VkZnJvbnQuVG'
    'Fnc1IEdGFncw==');

@$core.Deprecated('Use createDistributionTenantResultDescriptor instead')
const CreateDistributionTenantResult$json = {
  '1': 'CreateDistributionTenantResult',
  '2': [
    {
      '1': 'distributiontenant',
      '3': 510856916,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionTenant',
      '8': {},
      '10': 'distributiontenant'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `CreateDistributionTenantResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDistributionTenantResultDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVEaXN0cmlidXRpb25UZW5hbnRSZXN1bHQSWAoSZGlzdHJpYnV0aW9udGVuYW50GN'
        'SdzPMBIAEoCzIeLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uVGVuYW50QgSItRgBUhJkaXN0cmli'
        'dXRpb250ZW5hbnQSFgoEZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core.Deprecated('Use createDistributionWithTagsRequestDescriptor instead')
const CreateDistributionWithTagsRequest$json = {
  '1': 'CreateDistributionWithTagsRequest',
  '2': [
    {
      '1': 'distributionconfigwithtags',
      '3': 195182985,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionConfigWithTags',
      '8': {},
      '10': 'distributionconfigwithtags'
    },
  ],
};

/// Descriptor for `CreateDistributionWithTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDistributionWithTagsRequestDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVEaXN0cmlidXRpb25XaXRoVGFnc1JlcXVlc3QSbwoaZGlzdHJpYnV0aW9uY29uZm'
        'lnd2l0aHRhZ3MYiYOJXSABKAsyJi5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkNvbmZpZ1dpdGhU'
        'YWdzQgSItRgBUhpkaXN0cmlidXRpb25jb25maWd3aXRodGFncw==');

@$core.Deprecated('Use createDistributionWithTagsResultDescriptor instead')
const CreateDistributionWithTagsResult$json = {
  '1': 'CreateDistributionWithTagsResult',
  '2': [
    {
      '1': 'distribution',
      '3': 105183308,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Distribution',
      '8': {},
      '10': 'distribution'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateDistributionWithTagsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDistributionWithTagsResultDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVEaXN0cmlidXRpb25XaXRoVGFnc1Jlc3VsdBJFCgxkaXN0cmlidXRpb24YzPCTMi'
        'ABKAsyGC5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkIEiLUYAVIMZGlzdHJpYnV0aW9uEhYKBGV0'
        'YWcYgd+zlQEgASgJUgRldGFnEh4KCGxvY2F0aW9uGMebgt4BIAEoCVIIbG9jYXRpb24=');

@$core
    .Deprecated('Use createFieldLevelEncryptionConfigRequestDescriptor instead')
const CreateFieldLevelEncryptionConfigRequest$json = {
  '1': 'CreateFieldLevelEncryptionConfigRequest',
  '2': [
    {
      '1': 'fieldlevelencryptionconfig',
      '3': 499294709,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionConfig',
      '8': {},
      '10': 'fieldlevelencryptionconfig'
    },
  ],
};

/// Descriptor for `CreateFieldLevelEncryptionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFieldLevelEncryptionConfigRequestDescriptor =
    $convert.base64Decode(
        'CidDcmVhdGVGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ1JlcXVlc3QScAoaZmllbGRsZXZlbG'
        'VuY3J5cHRpb25jb25maWcY9cOK7gEgASgLMiYuY2xvdWRmcm9udC5GaWVsZExldmVsRW5jcnlw'
        'dGlvbkNvbmZpZ0IEiLUYAVIaZmllbGRsZXZlbGVuY3J5cHRpb25jb25maWc=');

@$core
    .Deprecated('Use createFieldLevelEncryptionConfigResultDescriptor instead')
const CreateFieldLevelEncryptionConfigResult$json = {
  '1': 'CreateFieldLevelEncryptionConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryption',
      '3': 473382747,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryption',
      '8': {},
      '10': 'fieldlevelencryption'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateFieldLevelEncryptionConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFieldLevelEncryptionConfigResultDescriptor =
    $convert.base64Decode(
        'CiZDcmVhdGVGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ1Jlc3VsdBIWCgRldGFnGIHfs5UBIA'
        'EoCVIEZXRhZxJeChRmaWVsZGxldmVsZW5jcnlwdGlvbhjb/tzhASABKAsyIC5jbG91ZGZyb250'
        'LkZpZWxkTGV2ZWxFbmNyeXB0aW9uQgSItRgBUhRmaWVsZGxldmVsZW5jcnlwdGlvbhIeCghsb2'
        'NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated(
    'Use createFieldLevelEncryptionProfileRequestDescriptor instead')
const CreateFieldLevelEncryptionProfileRequest$json = {
  '1': 'CreateFieldLevelEncryptionProfileRequest',
  '2': [
    {
      '1': 'fieldlevelencryptionprofileconfig',
      '3': 199371734,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfileConfig',
      '8': {},
      '10': 'fieldlevelencryptionprofileconfig'
    },
  ],
};

/// Descriptor for `CreateFieldLevelEncryptionProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFieldLevelEncryptionProfileRequestDescriptor =
    $convert.base64Decode(
        'CihDcmVhdGVGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVSZXF1ZXN0EoQBCiFmaWVsZGxldm'
        'VsZW5jcnlwdGlvbnByb2ZpbGVjb25maWcY1teIXyABKAsyLS5jbG91ZGZyb250LkZpZWxkTGV2'
        'ZWxFbmNyeXB0aW9uUHJvZmlsZUNvbmZpZ0IEiLUYAVIhZmllbGRsZXZlbGVuY3J5cHRpb25wcm'
        '9maWxlY29uZmln');

@$core
    .Deprecated('Use createFieldLevelEncryptionProfileResultDescriptor instead')
const CreateFieldLevelEncryptionProfileResult$json = {
  '1': 'CreateFieldLevelEncryptionProfileResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryptionprofile',
      '3': 344546136,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfile',
      '8': {},
      '10': 'fieldlevelencryptionprofile'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateFieldLevelEncryptionProfileResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFieldLevelEncryptionProfileResultDescriptor =
    $convert.base64Decode(
        'CidDcmVhdGVGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVSZXN1bHQSFgoEZXRhZxiB37OVAS'
        'ABKAlSBGV0YWcScwobZmllbGRsZXZlbGVuY3J5cHRpb25wcm9maWxlGNi2paQBIAEoCzInLmNs'
        'b3VkZnJvbnQuRmllbGRMZXZlbEVuY3J5cHRpb25Qcm9maWxlQgSItRgBUhtmaWVsZGxldmVsZW'
        '5jcnlwdGlvbnByb2ZpbGUSHgoIbG9jYXRpb24Yx5uC3gEgASgJUghsb2NhdGlvbg==');

@$core.Deprecated('Use createFunctionRequestDescriptor instead')
const CreateFunctionRequest$json = {
  '1': 'CreateFunctionRequest',
  '2': [
    {
      '1': 'functioncode',
      '3': 405947809,
      '4': 1,
      '5': 12,
      '10': 'functioncode'
    },
    {
      '1': 'functionconfig',
      '3': 116111484,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionConfig',
      '10': 'functionconfig'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreateFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFunctionRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVGdW5jdGlvblJlcXVlc3QSJgoMZnVuY3Rpb25jb2RlGKGLycEBIAEoDFIMZnVuY3'
    'Rpb25jb2RlEkUKDmZ1bmN0aW9uY29uZmlnGPzwrjcgASgLMhouY2xvdWRmcm9udC5GdW5jdGlv'
    'bkNvbmZpZ1IOZnVuY3Rpb25jb25maWcSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use createFunctionResultDescriptor instead')
const CreateFunctionResult$json = {
  '1': 'CreateFunctionResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'functionsummary',
      '3': 523316264,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionSummary',
      '8': {},
      '10': 'functionsummary'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFunctionResultDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVGdW5jdGlvblJlc3VsdBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxJPCg9mdW5jdG'
    'lvbnN1bW1hcnkYqNjE+QEgASgLMhsuY2xvdWRmcm9udC5GdW5jdGlvblN1bW1hcnlCBIi1GAFS'
    'D2Z1bmN0aW9uc3VtbWFyeRIeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated(
    'Use createInvalidationForDistributionTenantRequestDescriptor instead')
const CreateInvalidationForDistributionTenantRequest$json = {
  '1': 'CreateInvalidationForDistributionTenantRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'invalidationbatch',
      '3': 458631988,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.InvalidationBatch',
      '8': {},
      '10': 'invalidationbatch'
    },
  ],
};

/// Descriptor for `CreateInvalidationForDistributionTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createInvalidationForDistributionTenantRequestDescriptor =
    $convert.base64Decode(
        'Ci5DcmVhdGVJbnZhbGlkYXRpb25Gb3JEaXN0cmlidXRpb25UZW5hbnRSZXF1ZXN0EhIKAmlkGI'
        'HyorcBIAEoCVICaWQSVQoRaW52YWxpZGF0aW9uYmF0Y2gYtNbY2gEgASgLMh0uY2xvdWRmcm9u'
        'dC5JbnZhbGlkYXRpb25CYXRjaEIEiLUYAVIRaW52YWxpZGF0aW9uYmF0Y2g=');

@$core.Deprecated(
    'Use createInvalidationForDistributionTenantResultDescriptor instead')
const CreateInvalidationForDistributionTenantResult$json = {
  '1': 'CreateInvalidationForDistributionTenantResult',
  '2': [
    {
      '1': 'invalidation',
      '3': 77924830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Invalidation',
      '8': {},
      '10': 'invalidation'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateInvalidationForDistributionTenantResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createInvalidationForDistributionTenantResultDescriptor =
    $convert.base64Decode(
        'Ci1DcmVhdGVJbnZhbGlkYXRpb25Gb3JEaXN0cmlidXRpb25UZW5hbnRSZXN1bHQSRQoMaW52YW'
        'xpZGF0aW9uGN6TlCUgASgLMhguY2xvdWRmcm9udC5JbnZhbGlkYXRpb25CBIi1GAFSDGludmFs'
        'aWRhdGlvbhIeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated('Use createInvalidationRequestDescriptor instead')
const CreateInvalidationRequest$json = {
  '1': 'CreateInvalidationRequest',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'invalidationbatch',
      '3': 458631988,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.InvalidationBatch',
      '8': {},
      '10': 'invalidationbatch'
    },
  ],
};

/// Descriptor for `CreateInvalidationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createInvalidationRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVJbnZhbGlkYXRpb25SZXF1ZXN0EikKDmRpc3RyaWJ1dGlvbmlkGOex+0MgASgJUg'
    '5kaXN0cmlidXRpb25pZBJVChFpbnZhbGlkYXRpb25iYXRjaBi01tjaASABKAsyHS5jbG91ZGZy'
    'b250LkludmFsaWRhdGlvbkJhdGNoQgSItRgBUhFpbnZhbGlkYXRpb25iYXRjaA==');

@$core.Deprecated('Use createInvalidationResultDescriptor instead')
const CreateInvalidationResult$json = {
  '1': 'CreateInvalidationResult',
  '2': [
    {
      '1': 'invalidation',
      '3': 77924830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Invalidation',
      '8': {},
      '10': 'invalidation'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateInvalidationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createInvalidationResultDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVJbnZhbGlkYXRpb25SZXN1bHQSRQoMaW52YWxpZGF0aW9uGN6TlCUgASgLMhguY2'
    'xvdWRmcm9udC5JbnZhbGlkYXRpb25CBIi1GAFSDGludmFsaWRhdGlvbhIeCghsb2NhdGlvbhjH'
    'm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated('Use createKeyGroupRequestDescriptor instead')
const CreateKeyGroupRequest$json = {
  '1': 'CreateKeyGroupRequest',
  '2': [
    {
      '1': 'keygroupconfig',
      '3': 143012494,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroupConfig',
      '8': {},
      '10': 'keygroupconfig'
    },
  ],
};

/// Descriptor for `CreateKeyGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeyGroupRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVLZXlHcm91cFJlcXVlc3QSSwoOa2V5Z3JvdXBjb25maWcYjuWYRCABKAsyGi5jbG'
    '91ZGZyb250LktleUdyb3VwQ29uZmlnQgSItRgBUg5rZXlncm91cGNvbmZpZw==');

@$core.Deprecated('Use createKeyGroupResultDescriptor instead')
const CreateKeyGroupResult$json = {
  '1': 'CreateKeyGroupResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'keygroup',
      '3': 518748096,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroup',
      '8': {},
      '10': 'keygroup'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateKeyGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeyGroupResultDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVLZXlHcm91cFJlc3VsdBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxI6CghrZXlncm'
    '91cBjA7633ASABKAsyFC5jbG91ZGZyb250LktleUdyb3VwQgSItRgBUghrZXlncm91cBIeCghs'
    'b2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated('Use createKeyValueStoreRequestDescriptor instead')
const CreateKeyValueStoreRequest$json = {
  '1': 'CreateKeyValueStoreRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'importsource',
      '3': 41128754,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ImportSource',
      '10': 'importsource'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreateKeyValueStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeyValueStoreRequestDescriptor =
    $convert.base64Decode(
        'ChpDcmVhdGVLZXlWYWx1ZVN0b3JlUmVxdWVzdBIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbW'
        'VudBI/CgxpbXBvcnRzb3VyY2UYsqbOEyABKAsyGC5jbG91ZGZyb250LkltcG9ydFNvdXJjZVIM'
        'aW1wb3J0c291cmNlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use createKeyValueStoreResultDescriptor instead')
const CreateKeyValueStoreResult$json = {
  '1': 'CreateKeyValueStoreResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'keyvaluestore',
      '3': 151113103,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyValueStore',
      '8': {},
      '10': 'keyvaluestore'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateKeyValueStoreResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeyValueStoreResultDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVLZXlWYWx1ZVN0b3JlUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEkgKDW'
    'tleXZhbHVlc3RvcmUYj5uHSCABKAsyGS5jbG91ZGZyb250LktleVZhbHVlU3RvcmVCBIi1GAFS'
    'DWtleXZhbHVlc3RvcmUSHgoIbG9jYXRpb24Yx5uC3gEgASgJUghsb2NhdGlvbg==');

@$core.Deprecated('Use createMonitoringSubscriptionRequestDescriptor instead')
const CreateMonitoringSubscriptionRequest$json = {
  '1': 'CreateMonitoringSubscriptionRequest',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'monitoringsubscription',
      '3': 456621195,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.MonitoringSubscription',
      '8': {},
      '10': 'monitoringsubscription'
    },
  ],
};

/// Descriptor for `CreateMonitoringSubscriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMonitoringSubscriptionRequestDescriptor =
    $convert.base64Decode(
        'CiNDcmVhdGVNb25pdG9yaW5nU3Vic2NyaXB0aW9uUmVxdWVzdBIpCg5kaXN0cmlidXRpb25pZB'
        'jnsftDIAEoCVIOZGlzdHJpYnV0aW9uaWQSZAoWbW9uaXRvcmluZ3N1YnNjcmlwdGlvbhiL+d3Z'
        'ASABKAsyIi5jbG91ZGZyb250Lk1vbml0b3JpbmdTdWJzY3JpcHRpb25CBIi1GAFSFm1vbml0b3'
        'JpbmdzdWJzY3JpcHRpb24=');

@$core.Deprecated('Use createMonitoringSubscriptionResultDescriptor instead')
const CreateMonitoringSubscriptionResult$json = {
  '1': 'CreateMonitoringSubscriptionResult',
  '2': [
    {
      '1': 'monitoringsubscription',
      '3': 456621195,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.MonitoringSubscription',
      '8': {},
      '10': 'monitoringsubscription'
    },
  ],
};

/// Descriptor for `CreateMonitoringSubscriptionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMonitoringSubscriptionResultDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVNb25pdG9yaW5nU3Vic2NyaXB0aW9uUmVzdWx0EmQKFm1vbml0b3JpbmdzdWJzY3'
        'JpcHRpb24Yi/nd2QEgASgLMiIuY2xvdWRmcm9udC5Nb25pdG9yaW5nU3Vic2NyaXB0aW9uQgSI'
        'tRgBUhZtb25pdG9yaW5nc3Vic2NyaXB0aW9u');

@$core.Deprecated('Use createOriginAccessControlRequestDescriptor instead')
const CreateOriginAccessControlRequest$json = {
  '1': 'CreateOriginAccessControlRequest',
  '2': [
    {
      '1': 'originaccesscontrolconfig',
      '3': 143834977,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControlConfig',
      '8': {},
      '10': 'originaccesscontrolconfig'
    },
  ],
};

/// Descriptor for `CreateOriginAccessControlRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOriginAccessControlRequestDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVPcmlnaW5BY2Nlc3NDb250cm9sUmVxdWVzdBJsChlvcmlnaW5hY2Nlc3Njb250cm'
        '9sY29uZmlnGOH+ykQgASgLMiUuY2xvdWRmcm9udC5PcmlnaW5BY2Nlc3NDb250cm9sQ29uZmln'
        'QgSItRgBUhlvcmlnaW5hY2Nlc3Njb250cm9sY29uZmln');

@$core.Deprecated('Use createOriginAccessControlResultDescriptor instead')
const CreateOriginAccessControlResult$json = {
  '1': 'CreateOriginAccessControlResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'originaccesscontrol',
      '3': 238302375,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControl',
      '8': {},
      '10': 'originaccesscontrol'
    },
  ],
};

/// Descriptor for `CreateOriginAccessControlResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOriginAccessControlResultDescriptor =
    $convert.base64Decode(
        'Ch9DcmVhdGVPcmlnaW5BY2Nlc3NDb250cm9sUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldG'
        'FnEh4KCGxvY2F0aW9uGMebgt4BIAEoCVIIbG9jYXRpb24SWgoTb3JpZ2luYWNjZXNzY29udHJv'
        'bBin6dBxIAEoCzIfLmNsb3VkZnJvbnQuT3JpZ2luQWNjZXNzQ29udHJvbEIEiLUYAVITb3JpZ2'
        'luYWNjZXNzY29udHJvbA==');

@$core.Deprecated('Use createOriginRequestPolicyRequestDescriptor instead')
const CreateOriginRequestPolicyRequest$json = {
  '1': 'CreateOriginRequestPolicyRequest',
  '2': [
    {
      '1': 'originrequestpolicyconfig',
      '3': 37078133,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyConfig',
      '8': {},
      '10': 'originrequestpolicyconfig'
    },
  ],
};

/// Descriptor for `CreateOriginRequestPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOriginRequestPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVPcmlnaW5SZXF1ZXN0UG9saWN5UmVxdWVzdBJsChlvcmlnaW5yZXF1ZXN0cG9saW'
        'N5Y29uZmlnGPWI1xEgASgLMiUuY2xvdWRmcm9udC5PcmlnaW5SZXF1ZXN0UG9saWN5Q29uZmln'
        'QgSItRgBUhlvcmlnaW5yZXF1ZXN0cG9saWN5Y29uZmln');

@$core.Deprecated('Use createOriginRequestPolicyResultDescriptor instead')
const CreateOriginRequestPolicyResult$json = {
  '1': 'CreateOriginRequestPolicyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'originrequestpolicy',
      '3': 386733531,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicy',
      '8': {},
      '10': 'originrequestpolicy'
    },
  ],
};

/// Descriptor for `CreateOriginRequestPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOriginRequestPolicyResultDescriptor =
    $convert.base64Decode(
        'Ch9DcmVhdGVPcmlnaW5SZXF1ZXN0UG9saWN5UmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldG'
        'FnEh4KCGxvY2F0aW9uGMebgt4BIAEoCVIIbG9jYXRpb24SWwoTb3JpZ2lucmVxdWVzdHBvbGlj'
        'eRjbq7S4ASABKAsyHy5jbG91ZGZyb250Lk9yaWdpblJlcXVlc3RQb2xpY3lCBIi1GAFSE29yaW'
        'dpbnJlcXVlc3Rwb2xpY3k=');

@$core.Deprecated('Use createPublicKeyRequestDescriptor instead')
const CreatePublicKeyRequest$json = {
  '1': 'CreatePublicKeyRequest',
  '2': [
    {
      '1': 'publickeyconfig',
      '3': 228537966,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKeyConfig',
      '8': {},
      '10': 'publickeyconfig'
    },
  ],
};

/// Descriptor for `CreatePublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPublicKeyRequestDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVQdWJsaWNLZXlSZXF1ZXN0Ek4KD3B1YmxpY2tleWNvbmZpZxju7PxsIAEoCzIbLm'
        'Nsb3VkZnJvbnQuUHVibGljS2V5Q29uZmlnQgSItRgBUg9wdWJsaWNrZXljb25maWc=');

@$core.Deprecated('Use createPublicKeyResultDescriptor instead')
const CreatePublicKeyResult$json = {
  '1': 'CreatePublicKeyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'publickey',
      '3': 167335776,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKey',
      '8': {},
      '10': 'publickey'
    },
  ],
};

/// Descriptor for `CreatePublicKeyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPublicKeyResultDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVQdWJsaWNLZXlSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSHgoIbG9jYX'
    'Rpb24Yx5uC3gEgASgJUghsb2NhdGlvbhI8CglwdWJsaWNrZXkY4K7lTyABKAsyFS5jbG91ZGZy'
    'b250LlB1YmxpY0tleUIEiLUYAVIJcHVibGlja2V5');

@$core.Deprecated('Use createRealtimeLogConfigRequestDescriptor instead')
const CreateRealtimeLogConfigRequest$json = {
  '1': 'CreateRealtimeLogConfigRequest',
  '2': [
    {
      '1': 'endpoints',
      '3': 436023390,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.EndPoint',
      '10': 'endpoints'
    },
    {'1': 'fields', '3': 319339933, '4': 3, '5': 9, '10': 'fields'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'samplingrate', '3': 272929747, '4': 1, '5': 3, '10': 'samplingrate'},
  ],
};

/// Descriptor for `CreateRealtimeLogConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRealtimeLogConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVSZWFsdGltZUxvZ0NvbmZpZ1JlcXVlc3QSNgoJZW5kcG9pbnRzGN7g9M8BIAMoCz'
        'IULmNsb3VkZnJvbnQuRW5kUG9pbnRSCWVuZHBvaW50cxIaCgZmaWVsZHMYnfuimAEgAygJUgZm'
        'aWVsZHMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRImCgxzYW1wbGluZ3JhdGUY06eSggEgASgDUg'
        'xzYW1wbGluZ3JhdGU=');

@$core.Deprecated('Use createRealtimeLogConfigResultDescriptor instead')
const CreateRealtimeLogConfigResult$json = {
  '1': 'CreateRealtimeLogConfigResult',
  '2': [
    {
      '1': 'realtimelogconfig',
      '3': 95917609,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.RealtimeLogConfig',
      '10': 'realtimelogconfig'
    },
  ],
};

/// Descriptor for `CreateRealtimeLogConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRealtimeLogConfigResultDescriptor =
    $convert.base64Decode(
        'Ch1DcmVhdGVSZWFsdGltZUxvZ0NvbmZpZ1Jlc3VsdBJOChFyZWFsdGltZWxvZ2NvbmZpZxiprN'
        '4tIAEoCzIdLmNsb3VkZnJvbnQuUmVhbHRpbWVMb2dDb25maWdSEXJlYWx0aW1lbG9nY29uZmln');

@$core.Deprecated('Use createResponseHeadersPolicyRequestDescriptor instead')
const CreateResponseHeadersPolicyRequest$json = {
  '1': 'CreateResponseHeadersPolicyRequest',
  '2': [
    {
      '1': 'responseheaderspolicyconfig',
      '3': 159056825,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyConfig',
      '8': {},
      '10': 'responseheaderspolicyconfig'
    },
  ],
};

/// Descriptor for `CreateResponseHeadersPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResponseHeadersPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVSZXNwb25zZUhlYWRlcnNQb2xpY3lSZXF1ZXN0EnIKG3Jlc3BvbnNlaGVhZGVyc3'
        'BvbGljeWNvbmZpZxi5h+xLIAEoCzInLmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9saWN5'
        'Q29uZmlnQgSItRgBUhtyZXNwb25zZWhlYWRlcnNwb2xpY3ljb25maWc=');

@$core.Deprecated('Use createResponseHeadersPolicyResultDescriptor instead')
const CreateResponseHeadersPolicyResult$json = {
  '1': 'CreateResponseHeadersPolicyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'responseheaderspolicy',
      '3': 418204719,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicy',
      '8': {},
      '10': 'responseheaderspolicy'
    },
  ],
};

/// Descriptor for `CreateResponseHeadersPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResponseHeadersPolicyResultDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVSZXNwb25zZUhlYWRlcnNQb2xpY3lSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBG'
        'V0YWcSHgoIbG9jYXRpb24Yx5uC3gEgASgJUghsb2NhdGlvbhJhChVyZXNwb25zZWhlYWRlcnNw'
        'b2xpY3kYr5i1xwEgASgLMiEuY2xvdWRmcm9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lCBIi1GA'
        'FSFXJlc3BvbnNlaGVhZGVyc3BvbGljeQ==');

@$core.Deprecated('Use createStreamingDistributionRequestDescriptor instead')
const CreateStreamingDistributionRequest$json = {
  '1': 'CreateStreamingDistributionRequest',
  '2': [
    {
      '1': 'streamingdistributionconfig',
      '3': 291115944,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionConfig',
      '8': {},
      '10': 'streamingdistributionconfig'
    },
  ],
};

/// Descriptor for `CreateStreamingDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStreamingDistributionRequestDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVTdHJlYW1pbmdEaXN0cmlidXRpb25SZXF1ZXN0EnMKG3N0cmVhbWluZ2Rpc3RyaW'
        'J1dGlvbmNvbmZpZxiop+iKASABKAsyJy5jbG91ZGZyb250LlN0cmVhbWluZ0Rpc3RyaWJ1dGlv'
        'bkNvbmZpZ0IEiLUYAVIbc3RyZWFtaW5nZGlzdHJpYnV0aW9uY29uZmln');

@$core.Deprecated('Use createStreamingDistributionResultDescriptor instead')
const CreateStreamingDistributionResult$json = {
  '1': 'CreateStreamingDistributionResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'streamingdistribution',
      '3': 294813830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistribution',
      '8': {},
      '10': 'streamingdistribution'
    },
  ],
};

/// Descriptor for `CreateStreamingDistributionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStreamingDistributionResultDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVTdHJlYW1pbmdEaXN0cmlidXRpb25SZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBG'
        'V0YWcSHgoIbG9jYXRpb24Yx5uC3gEgASgJUghsb2NhdGlvbhJhChVzdHJlYW1pbmdkaXN0cmli'
        'dXRpb24YhoHKjAEgASgLMiEuY2xvdWRmcm9udC5TdHJlYW1pbmdEaXN0cmlidXRpb25CBIi1GA'
        'FSFXN0cmVhbWluZ2Rpc3RyaWJ1dGlvbg==');

@$core.Deprecated(
    'Use createStreamingDistributionWithTagsRequestDescriptor instead')
const CreateStreamingDistributionWithTagsRequest$json = {
  '1': 'CreateStreamingDistributionWithTagsRequest',
  '2': [
    {
      '1': 'streamingdistributionconfigwithtags',
      '3': 510934755,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionConfigWithTags',
      '8': {},
      '10': 'streamingdistributionconfigwithtags'
    },
  ],
};

/// Descriptor for `CreateStreamingDistributionWithTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createStreamingDistributionWithTagsRequestDescriptor =
    $convert.base64Decode(
        'CipDcmVhdGVTdHJlYW1pbmdEaXN0cmlidXRpb25XaXRoVGFnc1JlcXVlc3QSiwEKI3N0cmVhbW'
        'luZ2Rpc3RyaWJ1dGlvbmNvbmZpZ3dpdGh0YWdzGOP90PMBIAEoCzIvLmNsb3VkZnJvbnQuU3Ry'
        'ZWFtaW5nRGlzdHJpYnV0aW9uQ29uZmlnV2l0aFRhZ3NCBIi1GAFSI3N0cmVhbWluZ2Rpc3RyaW'
        'J1dGlvbmNvbmZpZ3dpdGh0YWdz');

@$core.Deprecated(
    'Use createStreamingDistributionWithTagsResultDescriptor instead')
const CreateStreamingDistributionWithTagsResult$json = {
  '1': 'CreateStreamingDistributionWithTagsResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'streamingdistribution',
      '3': 294813830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistribution',
      '8': {},
      '10': 'streamingdistribution'
    },
  ],
};

/// Descriptor for `CreateStreamingDistributionWithTagsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createStreamingDistributionWithTagsResultDescriptor = $convert.base64Decode(
        'CilDcmVhdGVTdHJlYW1pbmdEaXN0cmlidXRpb25XaXRoVGFnc1Jlc3VsdBIWCgRldGFnGIHfs5'
        'UBIAEoCVIEZXRhZxIeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9uEmEKFXN0cmVhbWlu'
        'Z2Rpc3RyaWJ1dGlvbhiGgcqMASABKAsyIS5jbG91ZGZyb250LlN0cmVhbWluZ0Rpc3RyaWJ1dG'
        'lvbkIEiLUYAVIVc3RyZWFtaW5nZGlzdHJpYnV0aW9u');

@$core.Deprecated('Use createTrustStoreRequestDescriptor instead')
const CreateTrustStoreRequest$json = {
  '1': 'CreateTrustStoreRequest',
  '2': [
    {
      '1': 'cacertificatesbundlesource',
      '3': 328448937,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CaCertificatesBundleSource',
      '10': 'cacertificatesbundlesource'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateTrustStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrustStoreRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVUcnVzdFN0b3JlUmVxdWVzdBJqChpjYWNlcnRpZmljYXRlc2J1bmRsZXNvdXJjZR'
    'ip986cASABKAsyJi5jbG91ZGZyb250LkNhQ2VydGlmaWNhdGVzQnVuZGxlU291cmNlUhpjYWNl'
    'cnRpZmljYXRlc2J1bmRsZXNvdXJjZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEigKBHRhZ3MYwc'
    'H2tQEgASgLMhAuY2xvdWRmcm9udC5UYWdzUgR0YWdz');

@$core.Deprecated('Use createTrustStoreResultDescriptor instead')
const CreateTrustStoreResult$json = {
  '1': 'CreateTrustStoreResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'truststore',
      '3': 224815327,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustStore',
      '8': {},
      '10': 'truststore'
    },
  ],
};

/// Descriptor for `CreateTrustStoreResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrustStoreResultDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVUcnVzdFN0b3JlUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEj8KCnRydX'
    'N0c3RvcmUY39GZayABKAsyFi5jbG91ZGZyb250LlRydXN0U3RvcmVCBIi1GAFSCnRydXN0c3Rv'
    'cmU=');

@$core.Deprecated('Use createVpcOriginRequestDescriptor instead')
const CreateVpcOriginRequest$json = {
  '1': 'CreateVpcOriginRequest',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
    {
      '1': 'vpcoriginendpointconfig',
      '3': 45660030,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOriginEndpointConfig',
      '10': 'vpcoriginendpointconfig'
    },
  ],
};

/// Descriptor for `CreateVpcOriginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createVpcOriginRequestDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVWcGNPcmlnaW5SZXF1ZXN0EigKBHRhZ3MYwcH2tQEgASgLMhAuY2xvdWRmcm9udC'
    '5UYWdzUgR0YWdzEmAKF3ZwY29yaWdpbmVuZHBvaW50Y29uZmlnGP7u4hUgASgLMiMuY2xvdWRm'
    'cm9udC5WcGNPcmlnaW5FbmRwb2ludENvbmZpZ1IXdnBjb3JpZ2luZW5kcG9pbnRjb25maWc=');

@$core.Deprecated('Use createVpcOriginResultDescriptor instead')
const CreateVpcOriginResult$json = {
  '1': 'CreateVpcOriginResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'vpcorigin',
      '3': 159181387,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOrigin',
      '8': {},
      '10': 'vpcorigin'
    },
  ],
};

/// Descriptor for `CreateVpcOriginResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createVpcOriginResultDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVWcGNPcmlnaW5SZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSHgoIbG9jYX'
    'Rpb24Yx5uC3gEgASgJUghsb2NhdGlvbhI8Cgl2cGNvcmlnaW4Yy9TzSyABKAsyFS5jbG91ZGZy'
    'b250LlZwY09yaWdpbkIEiLUYAVIJdnBjb3JpZ2lu');

@$core.Deprecated('Use customErrorResponseDescriptor instead')
const CustomErrorResponse$json = {
  '1': 'CustomErrorResponse',
  '2': [
    {
      '1': 'errorcachingminttl',
      '3': 111450431,
      '4': 1,
      '5': 3,
      '10': 'errorcachingminttl'
    },
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 5, '10': 'errorcode'},
    {'1': 'responsecode', '3': 447553700, '4': 1, '5': 9, '10': 'responsecode'},
    {
      '1': 'responsepagepath',
      '3': 523062007,
      '4': 1,
      '5': 9,
      '10': 'responsepagepath'
    },
  ],
};

/// Descriptor for `CustomErrorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customErrorResponseDescriptor = $convert.base64Decode(
    'ChNDdXN0b21FcnJvclJlc3BvbnNlEjEKEmVycm9yY2FjaGluZ21pbnR0bBi/spI1IAEoA1ISZX'
    'Jyb3JjYWNoaW5nbWludHRsEh8KCWVycm9yY29kZRiZ1sMQIAEoBVIJZXJyb3Jjb2RlEiYKDHJl'
    'c3BvbnNlY29kZRikwbTVASABKAlSDHJlc3BvbnNlY29kZRIuChByZXNwb25zZXBhZ2VwYXRoGP'
    'eVtfkBIAEoCVIQcmVzcG9uc2VwYWdlcGF0aA==');

@$core.Deprecated('Use customErrorResponsesDescriptor instead')
const CustomErrorResponses$json = {
  '1': 'CustomErrorResponses',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.CustomErrorResponse',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `CustomErrorResponses`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customErrorResponsesDescriptor = $convert.base64Decode(
    'ChRDdXN0b21FcnJvclJlc3BvbnNlcxI4CgVpdGVtcxiw8NgBIAMoCzIfLmNsb3VkZnJvbnQuQ3'
    'VzdG9tRXJyb3JSZXNwb25zZVIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use customHeadersDescriptor instead')
const CustomHeaders$json = {
  '1': 'CustomHeaders',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.OriginCustomHeader',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `CustomHeaders`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customHeadersDescriptor = $convert.base64Decode(
    'Cg1DdXN0b21IZWFkZXJzEjcKBWl0ZW1zGLDw2AEgAygLMh4uY2xvdWRmcm9udC5PcmlnaW5DdX'
    'N0b21IZWFkZXJSBWl0ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use customOriginConfigDescriptor instead')
const CustomOriginConfig$json = {
  '1': 'CustomOriginConfig',
  '2': [
    {'1': 'httpport', '3': 375976295, '4': 1, '5': 5, '10': 'httpport'},
    {'1': 'httpsport', '3': 153043918, '4': 1, '5': 5, '10': 'httpsport'},
    {
      '1': 'ipaddresstype',
      '3': 459110693,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.IpAddressType',
      '10': 'ipaddresstype'
    },
    {
      '1': 'originkeepalivetimeout',
      '3': 214128603,
      '4': 1,
      '5': 5,
      '10': 'originkeepalivetimeout'
    },
    {
      '1': 'originmtlsconfig',
      '3': 268477222,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginMtlsConfig',
      '10': 'originmtlsconfig'
    },
    {
      '1': 'originprotocolpolicy',
      '3': 234303342,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginProtocolPolicy',
      '10': 'originprotocolpolicy'
    },
    {
      '1': 'originreadtimeout',
      '3': 387717023,
      '4': 1,
      '5': 5,
      '10': 'originreadtimeout'
    },
    {
      '1': 'originsslprotocols',
      '3': 403853493,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginSslProtocols',
      '10': 'originsslprotocols'
    },
  ],
};

/// Descriptor for `CustomOriginConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customOriginConfigDescriptor = $convert.base64Decode(
    'ChJDdXN0b21PcmlnaW5Db25maWcSHgoIaHR0cHBvcnQY5+KjswEgASgFUghodHRwcG9ydBIfCg'
    'lodHRwc3BvcnQYzof9SCABKAVSCWh0dHBzcG9ydBJDCg1pcGFkZHJlc3N0eXBlGKXy9doBIAEo'
    'DjIZLmNsb3VkZnJvbnQuSXBBZGRyZXNzVHlwZVINaXBhZGRyZXNzdHlwZRI5ChZvcmlnaW5rZW'
    'VwYWxpdmV0aW1lb3V0GNuvjWYgASgFUhZvcmlnaW5rZWVwYWxpdmV0aW1lb3V0EkwKEG9yaWdp'
    'bm10bHNjb25maWcYpsaCgAEgASgLMhwuY2xvdWRmcm9udC5PcmlnaW5NdGxzQ29uZmlnUhBvcm'
    'lnaW5tdGxzY29uZmlnElcKFG9yaWdpbnByb3RvY29scG9saWN5GO7e3G8gASgOMiAuY2xvdWRm'
    'cm9udC5PcmlnaW5Qcm90b2NvbFBvbGljeVIUb3JpZ2lucHJvdG9jb2xwb2xpY3kSMAoRb3JpZ2'
    'lucmVhZHRpbWVvdXQYn6/wuAEgASgFUhFvcmlnaW5yZWFkdGltZW91dBJSChJvcmlnaW5zc2xw'
    'cm90b2NvbHMYtaHJwAEgASgLMh4uY2xvdWRmcm9udC5PcmlnaW5Tc2xQcm90b2NvbHNSEm9yaW'
    'dpbnNzbHByb3RvY29scw==');

@$core.Deprecated('Use customizationsDescriptor instead')
const Customizations$json = {
  '1': 'Customizations',
  '2': [
    {
      '1': 'certificate',
      '3': 198060817,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Certificate',
      '10': 'certificate'
    },
    {
      '1': 'georestrictions',
      '3': 309444102,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.GeoRestrictionCustomization',
      '10': 'georestrictions'
    },
    {
      '1': 'webacl',
      '3': 348583488,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.WebAclCustomization',
      '10': 'webacl'
    },
  ],
};

/// Descriptor for `Customizations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customizationsDescriptor = $convert.base64Decode(
    'Cg5DdXN0b21pemF0aW9ucxI8CgtjZXJ0aWZpY2F0ZRiR1rheIAEoCzIXLmNsb3VkZnJvbnQuQ2'
    'VydGlmaWNhdGVSC2NlcnRpZmljYXRlElUKD2dlb3Jlc3RyaWN0aW9ucxiG/MaTASABKAsyJy5j'
    'bG91ZGZyb250Lkdlb1Jlc3RyaWN0aW9uQ3VzdG9taXphdGlvblIPZ2VvcmVzdHJpY3Rpb25zEj'
    'sKBndlYmFjbBjA7JumASABKAsyHy5jbG91ZGZyb250LldlYkFjbEN1c3RvbWl6YXRpb25SBndl'
    'YmFjbA==');

@$core.Deprecated('Use defaultCacheBehaviorDescriptor instead')
const DefaultCacheBehavior$json = {
  '1': 'DefaultCacheBehavior',
  '2': [
    {
      '1': 'allowedmethods',
      '3': 56383476,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.AllowedMethods',
      '10': 'allowedmethods'
    },
    {
      '1': 'cachepolicyid',
      '3': 431434163,
      '4': 1,
      '5': 9,
      '10': 'cachepolicyid'
    },
    {'1': 'compress', '3': 235468462, '4': 1, '5': 8, '10': 'compress'},
    {'1': 'defaultttl', '3': 391646391, '4': 1, '5': 3, '10': 'defaultttl'},
    {
      '1': 'fieldlevelencryptionid',
      '3': 450714616,
      '4': 1,
      '5': 9,
      '10': 'fieldlevelencryptionid'
    },
    {
      '1': 'forwardedvalues',
      '3': 34815362,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ForwardedValues',
      '10': 'forwardedvalues'
    },
    {
      '1': 'functionassociations',
      '3': 457445650,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionAssociations',
      '10': 'functionassociations'
    },
    {
      '1': 'grpcconfig',
      '3': 406090728,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.GrpcConfig',
      '10': 'grpcconfig'
    },
    {
      '1': 'lambdafunctionassociations',
      '3': 46888655,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.LambdaFunctionAssociations',
      '10': 'lambdafunctionassociations'
    },
    {'1': 'maxttl', '3': 451484784, '4': 1, '5': 3, '10': 'maxttl'},
    {'1': 'minttl', '3': 420784162, '4': 1, '5': 3, '10': 'minttl'},
    {
      '1': 'originrequestpolicyid',
      '3': 298538616,
      '4': 1,
      '5': 9,
      '10': 'originrequestpolicyid'
    },
    {
      '1': 'realtimelogconfigarn',
      '3': 152963408,
      '4': 1,
      '5': 9,
      '10': 'realtimelogconfigarn'
    },
    {
      '1': 'responseheaderspolicyid',
      '3': 244029524,
      '4': 1,
      '5': 9,
      '10': 'responseheaderspolicyid'
    },
    {
      '1': 'smoothstreaming',
      '3': 92667114,
      '4': 1,
      '5': 8,
      '10': 'smoothstreaming'
    },
    {
      '1': 'targetoriginid',
      '3': 381807144,
      '4': 1,
      '5': 9,
      '10': 'targetoriginid'
    },
    {
      '1': 'trustedkeygroups',
      '3': 436720164,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustedKeyGroups',
      '10': 'trustedkeygroups'
    },
    {
      '1': 'trustedsigners',
      '3': 82003176,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustedSigners',
      '10': 'trustedsigners'
    },
    {
      '1': 'viewerprotocolpolicy',
      '3': 107534830,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ViewerProtocolPolicy',
      '10': 'viewerprotocolpolicy'
    },
  ],
};

/// Descriptor for `DefaultCacheBehavior`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List defaultCacheBehaviorDescriptor = $convert.base64Decode(
    'ChREZWZhdWx0Q2FjaGVCZWhhdmlvchJFCg5hbGxvd2VkbWV0aG9kcxj0r/EaIAEoCzIaLmNsb3'
    'VkZnJvbnQuQWxsb3dlZE1ldGhvZHNSDmFsbG93ZWRtZXRob2RzEigKDWNhY2hlcG9saWN5aWQY'
    's9PczQEgASgJUg1jYWNoZXBvbGljeWlkEh0KCGNvbXByZXNzGK7to3AgASgIUghjb21wcmVzcx'
    'IiCgpkZWZhdWx0dHRsGLeZ4LoBIAEoA1IKZGVmYXVsdHR0bBI6ChZmaWVsZGxldmVsZW5jcnlw'
    'dGlvbmlkGPi39dYBIAEoCVIWZmllbGRsZXZlbGVuY3J5cHRpb25pZBJICg9mb3J3YXJkZWR2YW'
    'x1ZXMYgvvMECABKAsyGy5jbG91ZGZyb250LkZvcndhcmRlZFZhbHVlc1IPZm9yd2FyZGVkdmFs'
    'dWVzElgKFGZ1bmN0aW9uYXNzb2NpYXRpb25zGJKikNoBIAEoCzIgLmNsb3VkZnJvbnQuRnVuY3'
    'Rpb25Bc3NvY2lhdGlvbnNSFGZ1bmN0aW9uYXNzb2NpYXRpb25zEjoKCmdycGNjb25maWcY6OfR'
    'wQEgASgLMhYuY2xvdWRmcm9udC5HcnBjQ29uZmlnUgpncnBjY29uZmlnEmkKGmxhbWJkYWZ1bm'
    'N0aW9uYXNzb2NpYXRpb25zGM/trRYgASgLMiYuY2xvdWRmcm9udC5MYW1iZGFGdW5jdGlvbkFz'
    'c29jaWF0aW9uc1IabGFtYmRhZnVuY3Rpb25hc3NvY2lhdGlvbnMSGgoGbWF4dHRsGPC4pNcBIA'
    'EoA1IGbWF4dHRsEhoKBm1pbnR0bBii0NLIASABKANSBm1pbnR0bBI4ChVvcmlnaW5yZXF1ZXN0'
    'cG9saWN5aWQY+KytjgEgASgJUhVvcmlnaW5yZXF1ZXN0cG9saWN5aWQSNQoUcmVhbHRpbWVsb2'
    'djb25maWdhcm4Y0JL4SCABKAlSFHJlYWx0aW1lbG9nY29uZmlnYXJuEjsKF3Jlc3BvbnNlaGVh'
    'ZGVyc3BvbGljeWlkGNSwrnQgASgJUhdyZXNwb25zZWhlYWRlcnNwb2xpY3lpZBIrCg9zbW9vdG'
    'hzdHJlYW1pbmcY6vmXLCABKAhSD3Ntb290aHN0cmVhbWluZxIqCg50YXJnZXRvcmlnaW5pZBio'
    '1Ie2ASABKAlSDnRhcmdldG9yaWdpbmlkEkwKEHRydXN0ZWRrZXlncm91cHMYpKSf0AEgASgLMh'
    'wuY2xvdWRmcm9udC5UcnVzdGVkS2V5R3JvdXBzUhB0cnVzdGVka2V5Z3JvdXBzEkUKDnRydXN0'
    'ZWRzaWduZXJzGOiJjScgASgLMhouY2xvdWRmcm9udC5UcnVzdGVkU2lnbmVyc1IOdHJ1c3RlZH'
    'NpZ25lcnMSVwoUdmlld2VycHJvdG9jb2xwb2xpY3kY7rOjMyABKA4yIC5jbG91ZGZyb250LlZp'
    'ZXdlclByb3RvY29sUG9saWN5UhR2aWV3ZXJwcm90b2NvbHBvbGljeQ==');

@$core.Deprecated('Use deleteAnycastIpListRequestDescriptor instead')
const DeleteAnycastIpListRequest$json = {
  '1': 'DeleteAnycastIpListRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteAnycastIpListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAnycastIpListRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVBbnljYXN0SXBMaXN0UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEhsKB2lmbW'
        'F0Y2gY0Ja3LCABKAlSB2lmbWF0Y2g=');

@$core.Deprecated('Use deleteCachePolicyRequestDescriptor instead')
const DeleteCachePolicyRequest$json = {
  '1': 'DeleteCachePolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteCachePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCachePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVDYWNoZVBvbGljeVJlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZBIbCgdpZm1hdG'
        'NoGNCWtywgASgJUgdpZm1hdGNo');

@$core.Deprecated(
    'Use deleteCloudFrontOriginAccessIdentityRequestDescriptor instead')
const DeleteCloudFrontOriginAccessIdentityRequest$json = {
  '1': 'DeleteCloudFrontOriginAccessIdentityRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteCloudFrontOriginAccessIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteCloudFrontOriginAccessIdentityRequestDescriptor =
    $convert.base64Decode(
        'CitEZWxldGVDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlSZXF1ZXN0EhIKAmlkGIHyor'
        'cBIAEoCVICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteConnectionFunctionRequestDescriptor instead')
const DeleteConnectionFunctionRequest$json = {
  '1': 'DeleteConnectionFunctionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteConnectionFunctionRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVDb25uZWN0aW9uRnVuY3Rpb25SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGw'
        'oHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteConnectionGroupRequestDescriptor instead')
const DeleteConnectionGroupRequest$json = {
  '1': 'DeleteConnectionGroupRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteConnectionGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteConnectionGroupRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVDb25uZWN0aW9uR3JvdXBSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaW'
        'ZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core
    .Deprecated('Use deleteContinuousDeploymentPolicyRequestDescriptor instead')
const DeleteContinuousDeploymentPolicyRequest$json = {
  '1': 'DeleteContinuousDeploymentPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteContinuousDeploymentPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteContinuousDeploymentPolicyRequestDescriptor =
    $convert.base64Decode(
        'CidEZWxldGVDb250aW51b3VzRGVwbG95bWVudFBvbGljeVJlcXVlc3QSEgoCaWQYgfKitwEgAS'
        'gJUgJpZBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNo');

@$core.Deprecated('Use deleteDistributionRequestDescriptor instead')
const DeleteDistributionRequest$json = {
  '1': 'DeleteDistributionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDistributionRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVEaXN0cmlidXRpb25SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYX'
        'RjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteDistributionTenantRequestDescriptor instead')
const DeleteDistributionTenantRequest$json = {
  '1': 'DeleteDistributionTenantRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteDistributionTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDistributionTenantRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVEaXN0cmlidXRpb25UZW5hbnRSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGw'
        'oHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core
    .Deprecated('Use deleteFieldLevelEncryptionConfigRequestDescriptor instead')
const DeleteFieldLevelEncryptionConfigRequest$json = {
  '1': 'DeleteFieldLevelEncryptionConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteFieldLevelEncryptionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFieldLevelEncryptionConfigRequestDescriptor =
    $convert.base64Decode(
        'CidEZWxldGVGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ1JlcXVlc3QSEgoCaWQYgfKitwEgAS'
        'gJUgJpZBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNo');

@$core.Deprecated(
    'Use deleteFieldLevelEncryptionProfileRequestDescriptor instead')
const DeleteFieldLevelEncryptionProfileRequest$json = {
  '1': 'DeleteFieldLevelEncryptionProfileRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteFieldLevelEncryptionProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFieldLevelEncryptionProfileRequestDescriptor =
    $convert.base64Decode(
        'CihEZWxldGVGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVSZXF1ZXN0EhIKAmlkGIHyorcBIA'
        'EoCVICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteFunctionRequestDescriptor instead')
const DeleteFunctionRequest$json = {
  '1': 'DeleteFunctionRequest',
  '2': [
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFunctionRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVGdW5jdGlvblJlcXVlc3QSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaBIVCg'
    'RuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use deleteKeyGroupRequestDescriptor instead')
const DeleteKeyGroupRequest$json = {
  '1': 'DeleteKeyGroupRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteKeyGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteKeyGroupRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVLZXlHcm91cFJlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZBIbCgdpZm1hdGNoGN'
    'CWtywgASgJUgdpZm1hdGNo');

@$core.Deprecated('Use deleteKeyValueStoreRequestDescriptor instead')
const DeleteKeyValueStoreRequest$json = {
  '1': 'DeleteKeyValueStoreRequest',
  '2': [
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteKeyValueStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteKeyValueStoreRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVLZXlWYWx1ZVN0b3JlUmVxdWVzdBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdG'
        'NoEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deleteMonitoringSubscriptionRequestDescriptor instead')
const DeleteMonitoringSubscriptionRequest$json = {
  '1': 'DeleteMonitoringSubscriptionRequest',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
  ],
};

/// Descriptor for `DeleteMonitoringSubscriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMonitoringSubscriptionRequestDescriptor =
    $convert.base64Decode(
        'CiNEZWxldGVNb25pdG9yaW5nU3Vic2NyaXB0aW9uUmVxdWVzdBIpCg5kaXN0cmlidXRpb25pZB'
        'jnsftDIAEoCVIOZGlzdHJpYnV0aW9uaWQ=');

@$core.Deprecated('Use deleteMonitoringSubscriptionResultDescriptor instead')
const DeleteMonitoringSubscriptionResult$json = {
  '1': 'DeleteMonitoringSubscriptionResult',
};

/// Descriptor for `DeleteMonitoringSubscriptionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMonitoringSubscriptionResultDescriptor =
    $convert.base64Decode('CiJEZWxldGVNb25pdG9yaW5nU3Vic2NyaXB0aW9uUmVzdWx0');

@$core.Deprecated('Use deleteOriginAccessControlRequestDescriptor instead')
const DeleteOriginAccessControlRequest$json = {
  '1': 'DeleteOriginAccessControlRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteOriginAccessControlRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteOriginAccessControlRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVPcmlnaW5BY2Nlc3NDb250cm9sUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEh'
        'sKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2g=');

@$core.Deprecated('Use deleteOriginRequestPolicyRequestDescriptor instead')
const DeleteOriginRequestPolicyRequest$json = {
  '1': 'DeleteOriginRequestPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteOriginRequestPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteOriginRequestPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVPcmlnaW5SZXF1ZXN0UG9saWN5UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEh'
        'sKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2g=');

@$core.Deprecated('Use deletePublicKeyRequestDescriptor instead')
const DeletePublicKeyRequest$json = {
  '1': 'DeletePublicKeyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeletePublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePublicKeyRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVQdWJsaWNLZXlSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaB'
        'jQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteRealtimeLogConfigRequestDescriptor instead')
const DeleteRealtimeLogConfigRequest$json = {
  '1': 'DeleteRealtimeLogConfigRequest',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteRealtimeLogConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRealtimeLogConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVSZWFsdGltZUxvZ0NvbmZpZ1JlcXVlc3QSFAoDYXJuGJ2cr70BIAEoCVIDYXJuEh'
        'UKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deleteResourcePolicyRequestDescriptor instead')
const DeleteResourcePolicyRequest$json = {
  '1': 'DeleteResourcePolicyRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `DeleteResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVSZXNvdXJjZVBvbGljeVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUg'
        'tyZXNvdXJjZWFybg==');

@$core.Deprecated('Use deleteResponseHeadersPolicyRequestDescriptor instead')
const DeleteResponseHeadersPolicyRequest$json = {
  '1': 'DeleteResponseHeadersPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteResponseHeadersPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResponseHeadersPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiJEZWxldGVSZXNwb25zZUhlYWRlcnNQb2xpY3lSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'QSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteStreamingDistributionRequestDescriptor instead')
const DeleteStreamingDistributionRequest$json = {
  '1': 'DeleteStreamingDistributionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteStreamingDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStreamingDistributionRequestDescriptor =
    $convert.base64Decode(
        'CiJEZWxldGVTdHJlYW1pbmdEaXN0cmlidXRpb25SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'QSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteTrustStoreRequestDescriptor instead')
const DeleteTrustStoreRequest$json = {
  '1': 'DeleteTrustStoreRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteTrustStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTrustStoreRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVUcnVzdFN0b3JlUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEhsKB2lmbWF0Y2'
        'gY0Ja3LCABKAlSB2lmbWF0Y2g=');

@$core.Deprecated('Use deleteVpcOriginRequestDescriptor instead')
const DeleteVpcOriginRequest$json = {
  '1': 'DeleteVpcOriginRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DeleteVpcOriginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteVpcOriginRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVWcGNPcmlnaW5SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaB'
        'jQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use deleteVpcOriginResultDescriptor instead')
const DeleteVpcOriginResult$json = {
  '1': 'DeleteVpcOriginResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'vpcorigin',
      '3': 159181387,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOrigin',
      '8': {},
      '10': 'vpcorigin'
    },
  ],
};

/// Descriptor for `DeleteVpcOriginResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteVpcOriginResultDescriptor = $convert.base64Decode(
    'ChVEZWxldGVWcGNPcmlnaW5SZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSPAoJdnBjb3'
    'JpZ2luGMvU80sgASgLMhUuY2xvdWRmcm9udC5WcGNPcmlnaW5CBIi1GAFSCXZwY29yaWdpbg==');

@$core.Deprecated('Use describeConnectionFunctionRequestDescriptor instead')
const DescribeConnectionFunctionRequest$json = {
  '1': 'DescribeConnectionFunctionRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `DescribeConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeConnectionFunctionRequestDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZUNvbm5lY3Rpb25GdW5jdGlvblJlcXVlc3QSIQoKaWRlbnRpZmllchjfoPsTIA'
        'EoCVIKaWRlbnRpZmllchIyCgVzdGFnZRjOl9hwIAEoDjIZLmNsb3VkZnJvbnQuRnVuY3Rpb25T'
        'dGFnZVIFc3RhZ2U=');

@$core.Deprecated('Use describeConnectionFunctionResultDescriptor instead')
const DescribeConnectionFunctionResult$json = {
  '1': 'DescribeConnectionFunctionResult',
  '2': [
    {
      '1': 'connectionfunctionsummary',
      '3': 62528396,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionSummary',
      '8': {},
      '10': 'connectionfunctionsummary'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `DescribeConnectionFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeConnectionFunctionResultDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZUNvbm5lY3Rpb25GdW5jdGlvblJlc3VsdBJsChljb25uZWN0aW9uZnVuY3Rpb2'
        '5zdW1tYXJ5GIy36B0gASgLMiUuY2xvdWRmcm9udC5Db25uZWN0aW9uRnVuY3Rpb25TdW1tYXJ5'
        'QgSItRgBUhljb25uZWN0aW9uZnVuY3Rpb25zdW1tYXJ5EhYKBGV0YWcYgd+zlQEgASgJUgRldG'
        'Fn');

@$core.Deprecated('Use describeFunctionRequestDescriptor instead')
const DescribeFunctionRequest$json = {
  '1': 'DescribeFunctionRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `DescribeFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeFunctionRequestDescriptor =
    $convert.base64Decode(
        'ChdEZXNjcmliZUZ1bmN0aW9uUmVxdWVzdBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjIKBXN0YW'
        'dlGM6X2HAgASgOMhkuY2xvdWRmcm9udC5GdW5jdGlvblN0YWdlUgVzdGFnZQ==');

@$core.Deprecated('Use describeFunctionResultDescriptor instead')
const DescribeFunctionResult$json = {
  '1': 'DescribeFunctionResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'functionsummary',
      '3': 523316264,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionSummary',
      '8': {},
      '10': 'functionsummary'
    },
  ],
};

/// Descriptor for `DescribeFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeFunctionResultDescriptor = $convert.base64Decode(
    'ChZEZXNjcmliZUZ1bmN0aW9uUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEk8KD2Z1bm'
    'N0aW9uc3VtbWFyeRio2MT5ASABKAsyGy5jbG91ZGZyb250LkZ1bmN0aW9uU3VtbWFyeUIEiLUY'
    'AVIPZnVuY3Rpb25zdW1tYXJ5');

@$core.Deprecated('Use describeKeyValueStoreRequestDescriptor instead')
const DescribeKeyValueStoreRequest$json = {
  '1': 'DescribeKeyValueStoreRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeKeyValueStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeKeyValueStoreRequestDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZUtleVZhbHVlU3RvcmVSZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use describeKeyValueStoreResultDescriptor instead')
const DescribeKeyValueStoreResult$json = {
  '1': 'DescribeKeyValueStoreResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'keyvaluestore',
      '3': 151113103,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyValueStore',
      '8': {},
      '10': 'keyvaluestore'
    },
  ],
};

/// Descriptor for `DescribeKeyValueStoreResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeKeyValueStoreResultDescriptor =
    $convert.base64Decode(
        'ChtEZXNjcmliZUtleVZhbHVlU3RvcmVSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSSA'
        'oNa2V5dmFsdWVzdG9yZRiPm4dIIAEoCzIZLmNsb3VkZnJvbnQuS2V5VmFsdWVTdG9yZUIEiLUY'
        'AVINa2V5dmFsdWVzdG9yZQ==');

@$core.Deprecated(
    'Use disassociateDistributionTenantWebACLRequestDescriptor instead')
const DisassociateDistributionTenantWebACLRequest$json = {
  '1': 'DisassociateDistributionTenantWebACLRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DisassociateDistributionTenantWebACLRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    disassociateDistributionTenantWebACLRequestDescriptor =
    $convert.base64Decode(
        'CitEaXNhc3NvY2lhdGVEaXN0cmlidXRpb25UZW5hbnRXZWJBQ0xSZXF1ZXN0EhIKAmlkGIHyor'
        'cBIAEoCVICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated(
    'Use disassociateDistributionTenantWebACLResultDescriptor instead')
const DisassociateDistributionTenantWebACLResult$json = {
  '1': 'DisassociateDistributionTenantWebACLResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DisassociateDistributionTenantWebACLResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    disassociateDistributionTenantWebACLResultDescriptor =
    $convert.base64Decode(
        'CipEaXNhc3NvY2lhdGVEaXN0cmlidXRpb25UZW5hbnRXZWJBQ0xSZXN1bHQSFgoEZXRhZxiB37'
        'OVASABKAlSBGV0YWcSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use disassociateDistributionWebACLRequestDescriptor instead')
const DisassociateDistributionWebACLRequest$json = {
  '1': 'DisassociateDistributionWebACLRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `DisassociateDistributionWebACLRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disassociateDistributionWebACLRequestDescriptor =
    $convert.base64Decode(
        'CiVEaXNhc3NvY2lhdGVEaXN0cmlidXRpb25XZWJBQ0xSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCV'
        'ICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated('Use disassociateDistributionWebACLResultDescriptor instead')
const DisassociateDistributionWebACLResult$json = {
  '1': 'DisassociateDistributionWebACLResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DisassociateDistributionWebACLResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disassociateDistributionWebACLResultDescriptor =
    $convert.base64Decode(
        'CiREaXNhc3NvY2lhdGVEaXN0cmlidXRpb25XZWJBQ0xSZXN1bHQSFgoEZXRhZxiB37OVASABKA'
        'lSBGV0YWcSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use distributionDescriptor instead')
const Distribution$json = {
  '1': 'Distribution',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'activetrustedkeygroups',
      '3': 204442122,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ActiveTrustedKeyGroups',
      '10': 'activetrustedkeygroups'
    },
    {
      '1': 'activetrustedsigners',
      '3': 356944374,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ActiveTrustedSigners',
      '10': 'activetrustedsigners'
    },
    {
      '1': 'aliasicprecordals',
      '3': 14607563,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.AliasICPRecordal',
      '10': 'aliasicprecordals'
    },
    {
      '1': 'distributionconfig',
      '3': 528940762,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionConfig',
      '10': 'distributionconfig'
    },
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'inprogressinvalidationbatches',
      '3': 143904528,
      '4': 1,
      '5': 5,
      '10': 'inprogressinvalidationbatches'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `Distribution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionDescriptor = $convert.base64Decode(
    'CgxEaXN0cmlidXRpb24SFAoDYXJuGJ2cr70BIAEoCVIDYXJuEl0KFmFjdGl2ZXRydXN0ZWRrZX'
    'lncm91cHMYipS+YSABKAsyIi5jbG91ZGZyb250LkFjdGl2ZVRydXN0ZWRLZXlHcm91cHNSFmFj'
    'dGl2ZXRydXN0ZWRrZXlncm91cHMSWAoUYWN0aXZldHJ1c3RlZHNpZ25lcnMY9pOaqgEgASgLMi'
    'AuY2xvdWRmcm9udC5BY3RpdmVUcnVzdGVkU2lnbmVyc1IUYWN0aXZldHJ1c3RlZHNpZ25lcnMS'
    'TQoRYWxpYXNpY3ByZWNvcmRhbHMYy8n7BiADKAsyHC5jbG91ZGZyb250LkFsaWFzSUNQUmVjb3'
    'JkYWxSEWFsaWFzaWNwcmVjb3JkYWxzElIKEmRpc3RyaWJ1dGlvbmNvbmZpZxja/Zv8ASABKAsy'
    'Hi5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkNvbmZpZ1ISZGlzdHJpYnV0aW9uY29uZmlnEiEKCm'
    'RvbWFpbm5hbWUY6834XCABKAlSCmRvbWFpbm5hbWUSEgoCaWQYgfKitwEgASgJUgJpZBJHCh1p'
    'bnByb2dyZXNzaW52YWxpZGF0aW9uYmF0Y2hlcxiQns9EIAEoBVIdaW5wcm9ncmVzc2ludmFsaW'
    'RhdGlvbmJhdGNoZXMSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVk'
    'dGltZRIZCgZzdGF0dXMYkOT7AiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use distributionAlreadyExistsDescriptor instead')
const DistributionAlreadyExists$json = {
  '1': 'DistributionAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DistributionAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChlEaXN0cmlidXRpb25BbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use distributionConfigDescriptor instead')
const DistributionConfig$json = {
  '1': 'DistributionConfig',
  '2': [
    {
      '1': 'aliases',
      '3': 476693696,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Aliases',
      '10': 'aliases'
    },
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
    {
      '1': 'cachebehaviors',
      '3': 70698397,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CacheBehaviors',
      '10': 'cachebehaviors'
    },
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'connectionfunctionassociation',
      '3': 253879893,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionAssociation',
      '10': 'connectionfunctionassociation'
    },
    {
      '1': 'connectionmode',
      '3': 82068023,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ConnectionMode',
      '10': 'connectionmode'
    },
    {
      '1': 'continuousdeploymentpolicyid',
      '3': 370023231,
      '4': 1,
      '5': 9,
      '10': 'continuousdeploymentpolicyid'
    },
    {
      '1': 'customerrorresponses',
      '3': 306704557,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CustomErrorResponses',
      '10': 'customerrorresponses'
    },
    {
      '1': 'defaultcachebehavior',
      '3': 346164111,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DefaultCacheBehavior',
      '10': 'defaultcachebehavior'
    },
    {
      '1': 'defaultrootobject',
      '3': 58521698,
      '4': 1,
      '5': 9,
      '10': 'defaultrootobject'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'httpversion',
      '3': 390398520,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.HttpVersion',
      '10': 'httpversion'
    },
    {
      '1': 'isipv6enabled',
      '3': 526764378,
      '4': 1,
      '5': 8,
      '10': 'isipv6enabled'
    },
    {
      '1': 'logging',
      '3': 65655615,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.LoggingConfig',
      '10': 'logging'
    },
    {
      '1': 'origingroups',
      '3': 397149984,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginGroups',
      '10': 'origingroups'
    },
    {
      '1': 'origins',
      '3': 277391381,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Origins',
      '10': 'origins'
    },
    {
      '1': 'priceclass',
      '3': 488692315,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.PriceClass',
      '10': 'priceclass'
    },
    {
      '1': 'restrictions',
      '3': 521105179,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Restrictions',
      '10': 'restrictions'
    },
    {'1': 'staging', '3': 193058759, '4': 1, '5': 8, '10': 'staging'},
    {
      '1': 'tenantconfig',
      '3': 477582312,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TenantConfig',
      '10': 'tenantconfig'
    },
    {
      '1': 'viewercertificate',
      '3': 214782047,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ViewerCertificate',
      '10': 'viewercertificate'
    },
    {
      '1': 'viewermtlsconfig',
      '3': 76850598,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ViewerMtlsConfig',
      '10': 'viewermtlsconfig'
    },
    {'1': 'webaclid', '3': 161274579, '4': 1, '5': 9, '10': 'webaclid'},
  ],
};

/// Descriptor for `DistributionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionConfigDescriptor = $convert.base64Decode(
    'ChJEaXN0cmlidXRpb25Db25maWcSMQoHYWxpYXNlcxjAiafjASABKAsyEy5jbG91ZGZyb250Lk'
    'FsaWFzZXNSB2FsaWFzZXMSLAoPYW55Y2FzdGlwbGlzdGlkGIOs+M0BIAEoCVIPYW55Y2FzdGlw'
    'bGlzdGlkEkUKDmNhY2hlYmVoYXZpb3JzGJ2L2yEgASgLMhouY2xvdWRmcm9udC5DYWNoZUJlaG'
    'F2aW9yc1IOY2FjaGViZWhhdmlvcnMSKwoPY2FsbGVycmVmZXJlbmNlGJiZjUggASgJUg9jYWxs'
    'ZXJyZWZlcmVuY2USHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1lbnQScgodY29ubmVjdGlvbm'
    'Z1bmN0aW9uYXNzb2NpYXRpb24Y1cyHeSABKAsyKS5jbG91ZGZyb250LkNvbm5lY3Rpb25GdW5j'
    'dGlvbkFzc29jaWF0aW9uUh1jb25uZWN0aW9uZnVuY3Rpb25hc3NvY2lhdGlvbhJFCg5jb25uZW'
    'N0aW9ubW9kZRi3hJEnIAEoDjIaLmNsb3VkZnJvbnQuQ29ubmVjdGlvbk1vZGVSDmNvbm5lY3Rp'
    'b25tb2RlEkYKHGNvbnRpbnVvdXNkZXBsb3ltZW50cG9saWN5aWQYv7a4sAEgASgJUhxjb250aW'
    '51b3VzZGVwbG95bWVudHBvbGljeWlkElgKFGN1c3RvbWVycm9ycmVzcG9uc2VzGK3hn5IBIAEo'
    'CzIgLmNsb3VkZnJvbnQuQ3VzdG9tRXJyb3JSZXNwb25zZXNSFGN1c3RvbWVycm9ycmVzcG9uc2'
    'VzElgKFGRlZmF1bHRjYWNoZWJlaGF2aW9yGI+XiKUBIAEoCzIgLmNsb3VkZnJvbnQuRGVmYXVs'
    'dENhY2hlQmVoYXZpb3JSFGRlZmF1bHRjYWNoZWJlaGF2aW9yEi8KEWRlZmF1bHRyb290b2JqZW'
    'N0GOLw8xsgASgJUhFkZWZhdWx0cm9vdG9iamVjdBIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5h'
    'YmxlZBI9CgtodHRwdmVyc2lvbhi4hJS6ASABKA4yFy5jbG91ZGZyb250Lkh0dHBWZXJzaW9uUg'
    'todHRwdmVyc2lvbhIoCg1pc2lwdjZlbmFibGVkGNqSl/sBIAEoCFINaXNpcHY2ZW5hYmxlZBI2'
    'Cgdsb2dnaW5nGL+mpx8gASgLMhkuY2xvdWRmcm9udC5Mb2dnaW5nQ29uZmlnUgdsb2dnaW5nEk'
    'AKDG9yaWdpbmdyb3VwcxigjrC9ASABKAsyGC5jbG91ZGZyb250Lk9yaWdpbkdyb3Vwc1IMb3Jp'
    'Z2luZ3JvdXBzEjEKB29yaWdpbnMYldCihAEgASgLMhMuY2xvdWRmcm9udC5PcmlnaW5zUgdvcm'
    'lnaW5zEjoKCnByaWNlY2xhc3MY27SD6QEgASgOMhYuY2xvdWRmcm9udC5QcmljZUNsYXNzUgpw'
    'cmljZWNsYXNzEkAKDHJlc3RyaWN0aW9ucxib3r34ASABKAsyGC5jbG91ZGZyb250LlJlc3RyaW'
    'N0aW9uc1IMcmVzdHJpY3Rpb25zEhsKB3N0YWdpbmcYx6+HXCABKAhSB3N0YWdpbmcSQAoMdGVu'
    'YW50Y29uZmlnGOin3eMBIAEoCzIYLmNsb3VkZnJvbnQuVGVuYW50Q29uZmlnUgx0ZW5hbnRjb2'
    '5maWcSTgoRdmlld2VyY2VydGlmaWNhdGUY36C1ZiABKAsyHS5jbG91ZGZyb250LlZpZXdlckNl'
    'cnRpZmljYXRlUhF2aWV3ZXJjZXJ0aWZpY2F0ZRJLChB2aWV3ZXJtdGxzY29uZmlnGKbL0iQgAS'
    'gLMhwuY2xvdWRmcm9udC5WaWV3ZXJNdGxzQ29uZmlnUhB2aWV3ZXJtdGxzY29uZmlnEh0KCHdl'
    'YmFjbGlkGNO180wgASgJUgh3ZWJhY2xpZA==');

@$core.Deprecated('Use distributionConfigWithTagsDescriptor instead')
const DistributionConfigWithTags$json = {
  '1': 'DistributionConfigWithTags',
  '2': [
    {
      '1': 'distributionconfig',
      '3': 528940762,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionConfig',
      '10': 'distributionconfig'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `DistributionConfigWithTags`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionConfigWithTagsDescriptor =
    $convert.base64Decode(
        'ChpEaXN0cmlidXRpb25Db25maWdXaXRoVGFncxJSChJkaXN0cmlidXRpb25jb25maWcY2v2b/A'
        'EgASgLMh4uY2xvdWRmcm9udC5EaXN0cmlidXRpb25Db25maWdSEmRpc3RyaWJ1dGlvbmNvbmZp'
        'ZxIoCgR0YWdzGMHB9rUBIAEoCzIQLmNsb3VkZnJvbnQuVGFnc1IEdGFncw==');

@$core.Deprecated('Use distributionIdListDescriptor instead')
const DistributionIdList$json = {
  '1': 'DistributionIdList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `DistributionIdList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionIdListDescriptor = $convert.base64Decode(
    'ChJEaXN0cmlidXRpb25JZExpc3QSIwoLaXN0cnVuY2F0ZWQY2p+4cyABKAhSC2lzdHJ1bmNhdG'
    'VkEhcKBWl0ZW1zGLDw2AEgAygJUgVpdGVtcxIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIe'
    'CghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUg'
    'puZXh0bWFya2VyEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use distributionIdOwnerDescriptor instead')
const DistributionIdOwner$json = {
  '1': 'DistributionIdOwner',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'owneraccountid',
      '3': 369721751,
      '4': 1,
      '5': 9,
      '10': 'owneraccountid'
    },
  ],
};

/// Descriptor for `DistributionIdOwner`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionIdOwnerDescriptor = $convert.base64Decode(
    'ChNEaXN0cmlidXRpb25JZE93bmVyEikKDmRpc3RyaWJ1dGlvbmlkGOex+0MgASgJUg5kaXN0cm'
    'lidXRpb25pZBIqCg5vd25lcmFjY291bnRpZBiXg6awASABKAlSDm93bmVyYWNjb3VudGlk');

@$core.Deprecated('Use distributionIdOwnerListDescriptor instead')
const DistributionIdOwnerList$json = {
  '1': 'DistributionIdOwnerList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DistributionIdOwner',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `DistributionIdOwnerList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionIdOwnerListDescriptor = $convert.base64Decode(
    'ChdEaXN0cmlidXRpb25JZE93bmVyTGlzdBIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cn'
    'VuY2F0ZWQSOAoFaXRlbXMYsPDYASADKAsyHy5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbklkT3du'
    'ZXJSBWl0ZW1zEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIA'
    'EoBVIIbWF4aXRlbXMSIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leHRtYXJrZXISHQoIcXVh'
    'bnRpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use distributionListDescriptor instead')
const DistributionList$json = {
  '1': 'DistributionList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DistributionSummary',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `DistributionList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionListDescriptor = $convert.base64Decode(
    'ChBEaXN0cmlidXRpb25MaXN0EiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZB'
    'I4CgVpdGVtcxiw8NgBIAMoCzIfLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uU3VtbWFyeVIFaXRl'
    'bXMSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYX'
    'hpdGVtcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj5'
    '5dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use distributionNotDisabledDescriptor instead')
const DistributionNotDisabled$json = {
  '1': 'DistributionNotDisabled',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DistributionNotDisabled`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionNotDisabledDescriptor =
    $convert.base64Decode(
        'ChdEaXN0cmlidXRpb25Ob3REaXNhYmxlZBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use distributionResourceIdDescriptor instead')
const DistributionResourceId$json = {
  '1': 'DistributionResourceId',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'distributiontenantid',
      '3': 123323327,
      '4': 1,
      '5': 9,
      '10': 'distributiontenantid'
    },
  ],
};

/// Descriptor for `DistributionResourceId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionResourceIdDescriptor = $convert.base64Decode(
    'ChZEaXN0cmlidXRpb25SZXNvdXJjZUlkEikKDmRpc3RyaWJ1dGlvbmlkGOex+0MgASgJUg5kaX'
    'N0cmlidXRpb25pZBI1ChRkaXN0cmlidXRpb250ZW5hbnRpZBi/h+c6IAEoCVIUZGlzdHJpYnV0'
    'aW9udGVuYW50aWQ=');

@$core.Deprecated('Use distributionSummaryDescriptor instead')
const DistributionSummary$json = {
  '1': 'DistributionSummary',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'aliasicprecordals',
      '3': 14607563,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.AliasICPRecordal',
      '10': 'aliasicprecordals'
    },
    {
      '1': 'aliases',
      '3': 476693696,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Aliases',
      '10': 'aliases'
    },
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
    {
      '1': 'cachebehaviors',
      '3': 70698397,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CacheBehaviors',
      '10': 'cachebehaviors'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'connectionfunctionassociation',
      '3': 253879893,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionAssociation',
      '10': 'connectionfunctionassociation'
    },
    {
      '1': 'connectionmode',
      '3': 82068023,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ConnectionMode',
      '10': 'connectionmode'
    },
    {
      '1': 'customerrorresponses',
      '3': 306704557,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CustomErrorResponses',
      '10': 'customerrorresponses'
    },
    {
      '1': 'defaultcachebehavior',
      '3': 346164111,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DefaultCacheBehavior',
      '10': 'defaultcachebehavior'
    },
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'httpversion',
      '3': 390398520,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.HttpVersion',
      '10': 'httpversion'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'isipv6enabled',
      '3': 526764378,
      '4': 1,
      '5': 8,
      '10': 'isipv6enabled'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'origingroups',
      '3': 397149984,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginGroups',
      '10': 'origingroups'
    },
    {
      '1': 'origins',
      '3': 277391381,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Origins',
      '10': 'origins'
    },
    {
      '1': 'priceclass',
      '3': 488692315,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.PriceClass',
      '10': 'priceclass'
    },
    {
      '1': 'restrictions',
      '3': 521105179,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Restrictions',
      '10': 'restrictions'
    },
    {'1': 'staging', '3': 193058759, '4': 1, '5': 8, '10': 'staging'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'viewercertificate',
      '3': 214782047,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ViewerCertificate',
      '10': 'viewercertificate'
    },
    {
      '1': 'viewermtlsconfig',
      '3': 76850598,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ViewerMtlsConfig',
      '10': 'viewermtlsconfig'
    },
    {'1': 'webaclid', '3': 161274579, '4': 1, '5': 9, '10': 'webaclid'},
  ],
};

/// Descriptor for `DistributionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionSummaryDescriptor = $convert.base64Decode(
    'ChNEaXN0cmlidXRpb25TdW1tYXJ5EhQKA2FybhidnK+9ASABKAlSA2FybhJNChFhbGlhc2ljcH'
    'JlY29yZGFscxjLyfsGIAMoCzIcLmNsb3VkZnJvbnQuQWxpYXNJQ1BSZWNvcmRhbFIRYWxpYXNp'
    'Y3ByZWNvcmRhbHMSMQoHYWxpYXNlcxjAiafjASABKAsyEy5jbG91ZGZyb250LkFsaWFzZXNSB2'
    'FsaWFzZXMSLAoPYW55Y2FzdGlwbGlzdGlkGIOs+M0BIAEoCVIPYW55Y2FzdGlwbGlzdGlkEkUK'
    'DmNhY2hlYmVoYXZpb3JzGJ2L2yEgASgLMhouY2xvdWRmcm9udC5DYWNoZUJlaGF2aW9yc1IOY2'
    'FjaGViZWhhdmlvcnMSHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1lbnQScgodY29ubmVjdGlv'
    'bmZ1bmN0aW9uYXNzb2NpYXRpb24Y1cyHeSABKAsyKS5jbG91ZGZyb250LkNvbm5lY3Rpb25GdW'
    '5jdGlvbkFzc29jaWF0aW9uUh1jb25uZWN0aW9uZnVuY3Rpb25hc3NvY2lhdGlvbhJFCg5jb25u'
    'ZWN0aW9ubW9kZRi3hJEnIAEoDjIaLmNsb3VkZnJvbnQuQ29ubmVjdGlvbk1vZGVSDmNvbm5lY3'
    'Rpb25tb2RlElgKFGN1c3RvbWVycm9ycmVzcG9uc2VzGK3hn5IBIAEoCzIgLmNsb3VkZnJvbnQu'
    'Q3VzdG9tRXJyb3JSZXNwb25zZXNSFGN1c3RvbWVycm9ycmVzcG9uc2VzElgKFGRlZmF1bHRjYW'
    'NoZWJlaGF2aW9yGI+XiKUBIAEoCzIgLmNsb3VkZnJvbnQuRGVmYXVsdENhY2hlQmVoYXZpb3JS'
    'FGRlZmF1bHRjYWNoZWJlaGF2aW9yEiEKCmRvbWFpbm5hbWUY6834XCABKAlSCmRvbWFpbm5hbW'
    'USFgoEZXRhZxiB37OVASABKAlSBGV0YWcSHAoHZW5hYmxlZBi/yJvkASABKAhSB2VuYWJsZWQS'
    'PQoLaHR0cHZlcnNpb24YuISUugEgASgOMhcuY2xvdWRmcm9udC5IdHRwVmVyc2lvblILaHR0cH'
    'ZlcnNpb24SEgoCaWQYgfKitwEgASgJUgJpZBIoCg1pc2lwdjZlbmFibGVkGNqSl/sBIAEoCFIN'
    'aXNpcHY2ZW5hYmxlZBItChBsYXN0bW9kaWZpZWR0aW1lGOCC/HAgASgJUhBsYXN0bW9kaWZpZW'
    'R0aW1lEkAKDG9yaWdpbmdyb3VwcxigjrC9ASABKAsyGC5jbG91ZGZyb250Lk9yaWdpbkdyb3Vw'
    'c1IMb3JpZ2luZ3JvdXBzEjEKB29yaWdpbnMYldCihAEgASgLMhMuY2xvdWRmcm9udC5PcmlnaW'
    '5zUgdvcmlnaW5zEjoKCnByaWNlY2xhc3MY27SD6QEgASgOMhYuY2xvdWRmcm9udC5QcmljZUNs'
    'YXNzUgpwcmljZWNsYXNzEkAKDHJlc3RyaWN0aW9ucxib3r34ASABKAsyGC5jbG91ZGZyb250Ll'
    'Jlc3RyaWN0aW9uc1IMcmVzdHJpY3Rpb25zEhsKB3N0YWdpbmcYx6+HXCABKAhSB3N0YWdpbmcS'
    'GQoGc3RhdHVzGJDk+wIgASgJUgZzdGF0dXMSTgoRdmlld2VyY2VydGlmaWNhdGUY36C1ZiABKA'
    'syHS5jbG91ZGZyb250LlZpZXdlckNlcnRpZmljYXRlUhF2aWV3ZXJjZXJ0aWZpY2F0ZRJLChB2'
    'aWV3ZXJtdGxzY29uZmlnGKbL0iQgASgLMhwuY2xvdWRmcm9udC5WaWV3ZXJNdGxzQ29uZmlnUh'
    'B2aWV3ZXJtdGxzY29uZmlnEh0KCHdlYmFjbGlkGNO180wgASgJUgh3ZWJhY2xpZA==');

@$core.Deprecated('Use distributionTenantDescriptor instead')
const DistributionTenant$json = {
  '1': 'DistributionTenant',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'connectiongroupid',
      '3': 169532206,
      '4': 1,
      '5': 9,
      '10': 'connectiongroupid'
    },
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {
      '1': 'customizations',
      '3': 70755200,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Customizations',
      '10': 'customizations'
    },
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'domains',
      '3': 149701959,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DomainResult',
      '10': 'domains'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.Parameter',
      '10': 'parameters'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `DistributionTenant`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionTenantDescriptor = $convert.base64Decode(
    'ChJEaXN0cmlidXRpb25UZW5hbnQSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEi8KEWNvbm5lY3Rpb2'
    '5ncm91cGlkGK6261AgASgJUhFjb25uZWN0aW9uZ3JvdXBpZBIjCgtjcmVhdGVkdGltZRjz6/M5'
    'IAEoCVILY3JlYXRlZHRpbWUSRQoOY3VzdG9taXphdGlvbnMYgMfeISABKAsyGi5jbG91ZGZyb2'
    '50LkN1c3RvbWl6YXRpb25zUg5jdXN0b21pemF0aW9ucxIpCg5kaXN0cmlidXRpb25pZBjnsftD'
    'IAEoCVIOZGlzdHJpYnV0aW9uaWQSNQoHZG9tYWlucxjHirFHIAMoCzIYLmNsb3VkZnJvbnQuRG'
    '9tYWluUmVzdWx0Ugdkb21haW5zEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVkEhIKAmlk'
    'GIHyorcBIAEoCVICaWQSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaW'
    'VkdGltZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjkKCnBhcmFtZXRlcnMY+qf+6wEgAygLMhUu'
    'Y2xvdWRmcm9udC5QYXJhbWV0ZXJSCnBhcmFtZXRlcnMSGQoGc3RhdHVzGJDk+wIgASgJUgZzdG'
    'F0dXMSKAoEdGFncxjBwfa1ASABKAsyEC5jbG91ZGZyb250LlRhZ3NSBHRhZ3M=');

@$core.Deprecated('Use distributionTenantAssociationFilterDescriptor instead')
const DistributionTenantAssociationFilter$json = {
  '1': 'DistributionTenantAssociationFilter',
  '2': [
    {
      '1': 'connectiongroupid',
      '3': 169532206,
      '4': 1,
      '5': 9,
      '10': 'connectiongroupid'
    },
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
  ],
};

/// Descriptor for `DistributionTenantAssociationFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionTenantAssociationFilterDescriptor =
    $convert.base64Decode(
        'CiNEaXN0cmlidXRpb25UZW5hbnRBc3NvY2lhdGlvbkZpbHRlchIvChFjb25uZWN0aW9uZ3JvdX'
        'BpZBiututQIAEoCVIRY29ubmVjdGlvbmdyb3VwaWQSKQoOZGlzdHJpYnV0aW9uaWQY57H7QyAB'
        'KAlSDmRpc3RyaWJ1dGlvbmlk');

@$core.Deprecated('Use distributionTenantSummaryDescriptor instead')
const DistributionTenantSummary$json = {
  '1': 'DistributionTenantSummary',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'connectiongroupid',
      '3': 169532206,
      '4': 1,
      '5': 9,
      '10': 'connectiongroupid'
    },
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {
      '1': 'customizations',
      '3': 70755200,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Customizations',
      '10': 'customizations'
    },
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'domains',
      '3': 149701959,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DomainResult',
      '10': 'domains'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `DistributionTenantSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distributionTenantSummaryDescriptor = $convert.base64Decode(
    'ChlEaXN0cmlidXRpb25UZW5hbnRTdW1tYXJ5EhQKA2Fybhidm+2/ASABKAlSA2FybhIvChFjb2'
    '5uZWN0aW9uZ3JvdXBpZBiututQIAEoCVIRY29ubmVjdGlvbmdyb3VwaWQSIwoLY3JlYXRlZHRp'
    'bWUY8+vzOSABKAlSC2NyZWF0ZWR0aW1lEkUKDmN1c3RvbWl6YXRpb25zGIDH3iEgASgLMhouY2'
    'xvdWRmcm9udC5DdXN0b21pemF0aW9uc1IOY3VzdG9taXphdGlvbnMSKQoOZGlzdHJpYnV0aW9u'
    'aWQY57H7QyABKAlSDmRpc3RyaWJ1dGlvbmlkEjUKB2RvbWFpbnMYx4qxRyADKAsyGC5jbG91ZG'
    'Zyb250LkRvbWFpblJlc3VsdFIHZG9tYWlucxIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxIcCgdl'
    'bmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBISCgJpZBiB8qK3ASABKAlSAmlkEi0KEGxhc3Rtb2'
    'RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWUSFQoEbmFtZRiH5oF/IAEoCVIE'
    'bmFtZRIZCgZzdGF0dXMYkOT7AiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use dnsConfigurationDescriptor instead')
const DnsConfiguration$json = {
  '1': 'DnsConfiguration',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.DnsConfigurationStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `DnsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dnsConfigurationDescriptor = $convert.base64Decode(
    'ChBEbnNDb25maWd1cmF0aW9uEhoKBmRvbWFpbhiSkvLwASABKAlSBmRvbWFpbhIZCgZyZWFzb2'
    '4YuoLFCSABKAlSBnJlYXNvbhI9CgZzdGF0dXMYkOT7AiABKA4yIi5jbG91ZGZyb250LkRuc0Nv'
    'bmZpZ3VyYXRpb25TdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use domainConflictDescriptor instead')
const DomainConflict$json = {
  '1': 'DomainConflict',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.DistributionResourceType',
      '10': 'resourcetype'
    },
  ],
};

/// Descriptor for `DomainConflict`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainConflictDescriptor = $convert.base64Decode(
    'Cg5Eb21haW5Db25mbGljdBIfCglhY2NvdW50aWQY0sG5HyABKAlSCWFjY291bnRpZBIaCgZkb2'
    '1haW4YkpLy8AEgASgJUgZkb21haW4SIgoKcmVzb3VyY2VpZBiRuvH6ASABKAlSCnJlc291cmNl'
    'aWQSTAoMcmVzb3VyY2V0eXBlGN6+2I8BIAEoDjIkLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uUm'
    'Vzb3VyY2VUeXBlUgxyZXNvdXJjZXR5cGU=');

@$core.Deprecated('Use domainItemDescriptor instead')
const DomainItem$json = {
  '1': 'DomainItem',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
  ],
};

/// Descriptor for `DomainItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainItemDescriptor = $convert
    .base64Decode('CgpEb21haW5JdGVtEhoKBmRvbWFpbhiSkvLwASABKAlSBmRvbWFpbg==');

@$core.Deprecated('Use domainResultDescriptor instead')
const DomainResult$json = {
  '1': 'DomainResult',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.DomainStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `DomainResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainResultDescriptor = $convert.base64Decode(
    'CgxEb21haW5SZXN1bHQSGgoGZG9tYWluGJKS8vABIAEoCVIGZG9tYWluEjMKBnN0YXR1cxiQ5P'
    'sCIAEoDjIYLmNsb3VkZnJvbnQuRG9tYWluU3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use encryptionEntitiesDescriptor instead')
const EncryptionEntities$json = {
  '1': 'EncryptionEntities',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.EncryptionEntity',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `EncryptionEntities`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionEntitiesDescriptor = $convert.base64Decode(
    'ChJFbmNyeXB0aW9uRW50aXRpZXMSNQoFaXRlbXMYsPDYASADKAsyHC5jbG91ZGZyb250LkVuY3'
    'J5cHRpb25FbnRpdHlSBWl0ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use encryptionEntityDescriptor instead')
const EncryptionEntity$json = {
  '1': 'EncryptionEntity',
  '2': [
    {
      '1': 'fieldpatterns',
      '3': 462160669,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldPatterns',
      '10': 'fieldpatterns'
    },
    {'1': 'providerid', '3': 509712370, '4': 1, '5': 9, '10': 'providerid'},
    {'1': 'publickeyid', '3': 15822003, '4': 1, '5': 9, '10': 'publickeyid'},
  ],
};

/// Descriptor for `EncryptionEntity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionEntityDescriptor = $convert.base64Decode(
    'ChBFbmNyeXB0aW9uRW50aXR5EkMKDWZpZWxkcGF0dGVybnMYnYaw3AEgASgLMhkuY2xvdWRmcm'
    '9udC5GaWVsZFBhdHRlcm5zUg1maWVsZHBhdHRlcm5zEiIKCnByb3ZpZGVyaWQY8q+G8wEgASgJ'
    'Ugpwcm92aWRlcmlkEiMKC3B1YmxpY2tleWlkGLPZxQcgASgJUgtwdWJsaWNrZXlpZA==');

@$core.Deprecated('Use endPointDescriptor instead')
const EndPoint$json = {
  '1': 'EndPoint',
  '2': [
    {
      '1': 'kinesisstreamconfig',
      '3': 104737892,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KinesisStreamConfig',
      '10': 'kinesisstreamconfig'
    },
    {'1': 'streamtype', '3': 536645422, '4': 1, '5': 9, '10': 'streamtype'},
  ],
};

/// Descriptor for `EndPoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endPointDescriptor = $convert.base64Decode(
    'CghFbmRQb2ludBJUChNraW5lc2lzc3RyZWFtY29uZmlnGOTY+DEgASgLMh8uY2xvdWRmcm9udC'
    '5LaW5lc2lzU3RyZWFtQ29uZmlnUhNraW5lc2lzc3RyZWFtY29uZmlnEiIKCnN0cmVhbXR5cGUY'
    'rp7y/wEgASgJUgpzdHJlYW10eXBl');

@$core.Deprecated('Use entityAlreadyExistsDescriptor instead')
const EntityAlreadyExists$json = {
  '1': 'EntityAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EntityAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChNFbnRpdHlBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use entityLimitExceededDescriptor instead')
const EntityLimitExceeded$json = {
  '1': 'EntityLimitExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EntityLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityLimitExceededDescriptor =
    $convert.base64Decode(
        'ChNFbnRpdHlMaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use entityNotFoundDescriptor instead')
const EntityNotFound$json = {
  '1': 'EntityNotFound',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EntityNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityNotFoundDescriptor = $convert.base64Decode(
    'Cg5FbnRpdHlOb3RGb3VuZBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use entitySizeLimitExceededDescriptor instead')
const EntitySizeLimitExceeded$json = {
  '1': 'EntitySizeLimitExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EntitySizeLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entitySizeLimitExceededDescriptor =
    $convert.base64Decode(
        'ChdFbnRpdHlTaXplTGltaXRFeGNlZWRlZBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use fieldLevelEncryptionDescriptor instead')
const FieldLevelEncryption$json = {
  '1': 'FieldLevelEncryption',
  '2': [
    {
      '1': 'fieldlevelencryptionconfig',
      '3': 499294709,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionConfig',
      '10': 'fieldlevelencryptionconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `FieldLevelEncryption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionDescriptor = $convert.base64Decode(
    'ChRGaWVsZExldmVsRW5jcnlwdGlvbhJqChpmaWVsZGxldmVsZW5jcnlwdGlvbmNvbmZpZxj1w4'
    'ruASABKAsyJi5jbG91ZGZyb250LkZpZWxkTGV2ZWxFbmNyeXB0aW9uQ29uZmlnUhpmaWVsZGxl'
    'dmVsZW5jcnlwdGlvbmNvbmZpZxISCgJpZBiB8qK3ASABKAlSAmlkEi0KEGxhc3Rtb2RpZmllZH'
    'RpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWU=');

@$core.Deprecated('Use fieldLevelEncryptionConfigDescriptor instead')
const FieldLevelEncryptionConfig$json = {
  '1': 'FieldLevelEncryptionConfig',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'contenttypeprofileconfig',
      '3': 361713406,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContentTypeProfileConfig',
      '10': 'contenttypeprofileconfig'
    },
    {
      '1': 'queryargprofileconfig',
      '3': 351336125,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.QueryArgProfileConfig',
      '10': 'queryargprofileconfig'
    },
  ],
};

/// Descriptor for `FieldLevelEncryptionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionConfigDescriptor = $convert.base64Decode(
    'ChpGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZxIrCg9jYWxsZXJyZWZlcmVuY2UYmJmNSCABKA'
    'lSD2NhbGxlcnJlZmVyZW5jZRIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbWVudBJkChhjb250'
    'ZW50dHlwZXByb2ZpbGVjb25maWcY/p29rAEgASgLMiQuY2xvdWRmcm9udC5Db250ZW50VHlwZV'
    'Byb2ZpbGVDb25maWdSGGNvbnRlbnR0eXBlcHJvZmlsZWNvbmZpZxJbChVxdWVyeWFyZ3Byb2Zp'
    'bGVjb25maWcYve3DpwEgASgLMiEuY2xvdWRmcm9udC5RdWVyeUFyZ1Byb2ZpbGVDb25maWdSFX'
    'F1ZXJ5YXJncHJvZmlsZWNvbmZpZw==');

@$core
    .Deprecated('Use fieldLevelEncryptionConfigAlreadyExistsDescriptor instead')
const FieldLevelEncryptionConfigAlreadyExists$json = {
  '1': 'FieldLevelEncryptionConfigAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FieldLevelEncryptionConfigAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionConfigAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CidGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ0FscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use fieldLevelEncryptionConfigInUseDescriptor instead')
const FieldLevelEncryptionConfigInUse$json = {
  '1': 'FieldLevelEncryptionConfigInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FieldLevelEncryptionConfigInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionConfigInUseDescriptor =
    $convert.base64Decode(
        'Ch9GaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ0luVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use fieldLevelEncryptionListDescriptor instead')
const FieldLevelEncryptionList$json = {
  '1': 'FieldLevelEncryptionList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionSummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `FieldLevelEncryptionList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionListDescriptor = $convert.base64Decode(
    'ChhGaWVsZExldmVsRW5jcnlwdGlvbkxpc3QSQAoFaXRlbXMYsPDYASADKAsyJy5jbG91ZGZyb2'
    '50LkZpZWxkTGV2ZWxFbmNyeXB0aW9uU3VtbWFyeVIFaXRlbXMSHgoIbWF4aXRlbXMYlNba8QEg'
    'ASgFUghtYXhpdGVtcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdW'
    'FudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use fieldLevelEncryptionProfileDescriptor instead')
const FieldLevelEncryptionProfile$json = {
  '1': 'FieldLevelEncryptionProfile',
  '2': [
    {
      '1': 'fieldlevelencryptionprofileconfig',
      '3': 199371734,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfileConfig',
      '10': 'fieldlevelencryptionprofileconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `FieldLevelEncryptionProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionProfileDescriptor = $convert.base64Decode(
    'ChtGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGUSfgohZmllbGRsZXZlbGVuY3J5cHRpb25wcm'
    '9maWxlY29uZmlnGNbXiF8gASgLMi0uY2xvdWRmcm9udC5GaWVsZExldmVsRW5jcnlwdGlvblBy'
    'b2ZpbGVDb25maWdSIWZpZWxkbGV2ZWxlbmNyeXB0aW9ucHJvZmlsZWNvbmZpZxISCgJpZBiB8q'
    'K3ASABKAlSAmlkEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRp'
    'bWU=');

@$core.Deprecated(
    'Use fieldLevelEncryptionProfileAlreadyExistsDescriptor instead')
const FieldLevelEncryptionProfileAlreadyExists$json = {
  '1': 'FieldLevelEncryptionProfileAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FieldLevelEncryptionProfileAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionProfileAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CihGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use fieldLevelEncryptionProfileConfigDescriptor instead')
const FieldLevelEncryptionProfileConfig$json = {
  '1': 'FieldLevelEncryptionProfileConfig',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'encryptionentities',
      '3': 400717088,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.EncryptionEntities',
      '10': 'encryptionentities'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `FieldLevelEncryptionProfileConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionProfileConfigDescriptor =
    $convert.base64Decode(
        'CiFGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVDb25maWcSKwoPY2FsbGVycmVmZXJlbmNlGJ'
        'iZjUggASgJUg9jYWxsZXJyZWZlcmVuY2USHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1lbnQS'
        'UgoSZW5jcnlwdGlvbmVudGl0aWVzGKDqib8BIAEoCzIeLmNsb3VkZnJvbnQuRW5jcnlwdGlvbk'
        'VudGl0aWVzUhJlbmNyeXB0aW9uZW50aXRpZXMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use fieldLevelEncryptionProfileInUseDescriptor instead')
const FieldLevelEncryptionProfileInUse$json = {
  '1': 'FieldLevelEncryptionProfileInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FieldLevelEncryptionProfileInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionProfileInUseDescriptor =
    $convert.base64Decode(
        'CiBGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVJblVzZRIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use fieldLevelEncryptionProfileListDescriptor instead')
const FieldLevelEncryptionProfileList$json = {
  '1': 'FieldLevelEncryptionProfileList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfileSummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `FieldLevelEncryptionProfileList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionProfileListDescriptor =
    $convert.base64Decode(
        'Ch9GaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVMaXN0EkcKBWl0ZW1zGLDw2AEgAygLMi4uY2'
        'xvdWRmcm9udC5GaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVTdW1tYXJ5UgVpdGVtcxIeCght'
        'YXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZX'
        'h0bWFya2VyEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core
    .Deprecated('Use fieldLevelEncryptionProfileSizeExceededDescriptor instead')
const FieldLevelEncryptionProfileSizeExceeded$json = {
  '1': 'FieldLevelEncryptionProfileSizeExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FieldLevelEncryptionProfileSizeExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionProfileSizeExceededDescriptor =
    $convert.base64Decode(
        'CidGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVTaXplRXhjZWVkZWQSGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use fieldLevelEncryptionProfileSummaryDescriptor instead')
const FieldLevelEncryptionProfileSummary$json = {
  '1': 'FieldLevelEncryptionProfileSummary',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'encryptionentities',
      '3': 400717088,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.EncryptionEntities',
      '10': 'encryptionentities'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `FieldLevelEncryptionProfileSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionProfileSummaryDescriptor =
    $convert.base64Decode(
        'CiJGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVTdW1tYXJ5EhwKB2NvbW1lbnQY/7++wgEgAS'
        'gJUgdjb21tZW50ElIKEmVuY3J5cHRpb25lbnRpdGllcxig6om/ASABKAsyHi5jbG91ZGZyb250'
        'LkVuY3J5cHRpb25FbnRpdGllc1ISZW5jcnlwdGlvbmVudGl0aWVzEhIKAmlkGIHyorcBIAEoCV'
        'ICaWQSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZRIVCgRu'
        'YW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use fieldLevelEncryptionSummaryDescriptor instead')
const FieldLevelEncryptionSummary$json = {
  '1': 'FieldLevelEncryptionSummary',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'contenttypeprofileconfig',
      '3': 361713406,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContentTypeProfileConfig',
      '10': 'contenttypeprofileconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'queryargprofileconfig',
      '3': 351336125,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.QueryArgProfileConfig',
      '10': 'queryargprofileconfig'
    },
  ],
};

/// Descriptor for `FieldLevelEncryptionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldLevelEncryptionSummaryDescriptor = $convert.base64Decode(
    'ChtGaWVsZExldmVsRW5jcnlwdGlvblN1bW1hcnkSHAoHY29tbWVudBj/v77CASABKAlSB2NvbW'
    '1lbnQSZAoYY29udGVudHR5cGVwcm9maWxlY29uZmlnGP6dvawBIAEoCzIkLmNsb3VkZnJvbnQu'
    'Q29udGVudFR5cGVQcm9maWxlQ29uZmlnUhhjb250ZW50dHlwZXByb2ZpbGVjb25maWcSEgoCaW'
    'QYgfKitwEgASgJUgJpZBItChBsYXN0bW9kaWZpZWR0aW1lGOCC/HAgASgJUhBsYXN0bW9kaWZp'
    'ZWR0aW1lElsKFXF1ZXJ5YXJncHJvZmlsZWNvbmZpZxi97cOnASABKAsyIS5jbG91ZGZyb250Ll'
    'F1ZXJ5QXJnUHJvZmlsZUNvbmZpZ1IVcXVlcnlhcmdwcm9maWxlY29uZmln');

@$core.Deprecated('Use fieldPatternsDescriptor instead')
const FieldPatterns$json = {
  '1': 'FieldPatterns',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `FieldPatterns`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldPatternsDescriptor = $convert.base64Decode(
    'Cg1GaWVsZFBhdHRlcm5zEhcKBWl0ZW1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eRj55d'
    'xfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use forwardedValuesDescriptor instead')
const ForwardedValues$json = {
  '1': 'ForwardedValues',
  '2': [
    {
      '1': 'cookies',
      '3': 418634949,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CookiePreference',
      '10': 'cookies'
    },
    {
      '1': 'headers',
      '3': 323967370,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Headers',
      '10': 'headers'
    },
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 8, '10': 'querystring'},
    {
      '1': 'querystringcachekeys',
      '3': 510632337,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.QueryStringCacheKeys',
      '10': 'querystringcachekeys'
    },
  ],
};

/// Descriptor for `ForwardedValues`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List forwardedValuesDescriptor = $convert.base64Decode(
    'Cg9Gb3J3YXJkZWRWYWx1ZXMSOgoHY29va2llcxjFuc/HASABKAsyHC5jbG91ZGZyb250LkNvb2'
    'tpZVByZWZlcmVuY2VSB2Nvb2tpZXMSMQoHaGVhZGVycxiKs72aASABKAsyEy5jbG91ZGZyb250'
    'LkhlYWRlcnNSB2hlYWRlcnMSJAoLcXVlcnlzdHJpbmcY58rvzwEgASgIUgtxdWVyeXN0cmluZx'
    'JYChRxdWVyeXN0cmluZ2NhY2hla2V5cxiRw77zASABKAsyIC5jbG91ZGZyb250LlF1ZXJ5U3Ry'
    'aW5nQ2FjaGVLZXlzUhRxdWVyeXN0cmluZ2NhY2hla2V5cw==');

@$core.Deprecated('Use functionAlreadyExistsDescriptor instead')
const FunctionAlreadyExists$json = {
  '1': 'FunctionAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FunctionAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionAlreadyExistsDescriptor = $convert.base64Decode(
    'ChVGdW5jdGlvbkFscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use functionAssociationDescriptor instead')
const FunctionAssociation$json = {
  '1': 'FunctionAssociation',
  '2': [
    {
      '1': 'eventtype',
      '3': 468897896,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.EventType',
      '10': 'eventtype'
    },
    {'1': 'functionarn', '3': 387130481, '4': 1, '5': 9, '10': 'functionarn'},
  ],
};

/// Descriptor for `FunctionAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionAssociationDescriptor = $convert.base64Decode(
    'ChNGdW5jdGlvbkFzc29jaWF0aW9uEjcKCWV2ZW50dHlwZRjooMvfASABKA4yFS5jbG91ZGZyb2'
    '50LkV2ZW50VHlwZVIJZXZlbnR0eXBlEiQKC2Z1bmN0aW9uYXJuGPHIzLgBIAEoCVILZnVuY3Rp'
    'b25hcm4=');

@$core.Deprecated('Use functionAssociationsDescriptor instead')
const FunctionAssociations$json = {
  '1': 'FunctionAssociations',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.FunctionAssociation',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `FunctionAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionAssociationsDescriptor = $convert.base64Decode(
    'ChRGdW5jdGlvbkFzc29jaWF0aW9ucxI4CgVpdGVtcxiw8NgBIAMoCzIfLmNsb3VkZnJvbnQuRn'
    'VuY3Rpb25Bc3NvY2lhdGlvblIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use functionConfigDescriptor instead')
const FunctionConfig$json = {
  '1': 'FunctionConfig',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'keyvaluestoreassociations',
      '3': 383569407,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyValueStoreAssociations',
      '10': 'keyvaluestoreassociations'
    },
    {
      '1': 'runtime',
      '3': 359311308,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionRuntime',
      '10': 'runtime'
    },
  ],
};

/// Descriptor for `FunctionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionConfigDescriptor = $convert.base64Decode(
    'Cg5GdW5jdGlvbkNvbmZpZxIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbWVudBJnChlrZXl2YW'
    'x1ZXN0b3JlYXNzb2NpYXRpb25zGP+b87YBIAEoCzIlLmNsb3VkZnJvbnQuS2V5VmFsdWVTdG9y'
    'ZUFzc29jaWF0aW9uc1IZa2V5dmFsdWVzdG9yZWFzc29jaWF0aW9ucxI5CgdydW50aW1lGMzPqq'
    'sBIAEoDjIbLmNsb3VkZnJvbnQuRnVuY3Rpb25SdW50aW1lUgdydW50aW1l');

@$core.Deprecated('Use functionInUseDescriptor instead')
const FunctionInUse$json = {
  '1': 'FunctionInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FunctionInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionInUseDescriptor = $convert.base64Decode(
    'Cg1GdW5jdGlvbkluVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use functionListDescriptor instead')
const FunctionList$json = {
  '1': 'FunctionList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.FunctionSummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `FunctionList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionListDescriptor = $convert.base64Decode(
    'CgxGdW5jdGlvbkxpc3QSNAoFaXRlbXMYsPDYASADKAsyGy5jbG91ZGZyb250LkZ1bmN0aW9uU3'
    'VtbWFyeVIFaXRlbXMSHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpuZXh0bWFy'
    'a2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdH'
    'k=');

@$core.Deprecated('Use functionMetadataDescriptor instead')
const FunctionMetadata$json = {
  '1': 'FunctionMetadata',
  '2': [
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'functionarn', '3': 387130481, '4': 1, '5': 9, '10': 'functionarn'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `FunctionMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionMetadataDescriptor = $convert.base64Decode(
    'ChBGdW5jdGlvbk1ldGFkYXRhEiMKC2NyZWF0ZWR0aW1lGPPr8zkgASgJUgtjcmVhdGVkdGltZR'
    'IkCgtmdW5jdGlvbmFybhjxyMy4ASABKAlSC2Z1bmN0aW9uYXJuEi0KEGxhc3Rtb2RpZmllZHRp'
    'bWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWUSMgoFc3RhZ2UYzpfYcCABKA4yGS5jbG91ZG'
    'Zyb250LkZ1bmN0aW9uU3RhZ2VSBXN0YWdl');

@$core.Deprecated('Use functionSizeLimitExceededDescriptor instead')
const FunctionSizeLimitExceeded$json = {
  '1': 'FunctionSizeLimitExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FunctionSizeLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionSizeLimitExceededDescriptor =
    $convert.base64Decode(
        'ChlGdW5jdGlvblNpemVMaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use functionSummaryDescriptor instead')
const FunctionSummary$json = {
  '1': 'FunctionSummary',
  '2': [
    {
      '1': 'functionconfig',
      '3': 116111484,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionConfig',
      '10': 'functionconfig'
    },
    {
      '1': 'functionmetadata',
      '3': 503503117,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionMetadata',
      '10': 'functionmetadata'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `FunctionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionSummaryDescriptor = $convert.base64Decode(
    'Cg9GdW5jdGlvblN1bW1hcnkSRQoOZnVuY3Rpb25jb25maWcY/PCuNyABKAsyGi5jbG91ZGZyb2'
    '50LkZ1bmN0aW9uQ29uZmlnUg5mdW5jdGlvbmNvbmZpZxJMChBmdW5jdGlvbm1ldGFkYXRhGI2y'
    'i/ABIAEoCzIcLmNsb3VkZnJvbnQuRnVuY3Rpb25NZXRhZGF0YVIQZnVuY3Rpb25tZXRhZGF0YR'
    'IVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhkKBnN0YXR1cxiQ5PsCIAEoCVIGc3RhdHVz');

@$core.Deprecated('Use geoRestrictionDescriptor instead')
const GeoRestriction$json = {
  '1': 'GeoRestriction',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
    {
      '1': 'restrictiontype',
      '3': 266540262,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.GeoRestrictionType',
      '10': 'restrictiontype'
    },
  ],
};

/// Descriptor for `GeoRestriction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List geoRestrictionDescriptor = $convert.base64Decode(
    'Cg5HZW9SZXN0cmljdGlvbhIXCgVpdGVtcxiw8NgBIAMoCVIFaXRlbXMSHQoIcXVhbnRpdHkY+e'
    'XcXyABKAVSCHF1YW50aXR5EksKD3Jlc3RyaWN0aW9udHlwZRjmqYx/IAEoDjIeLmNsb3VkZnJv'
    'bnQuR2VvUmVzdHJpY3Rpb25UeXBlUg9yZXN0cmljdGlvbnR5cGU=');

@$core.Deprecated('Use geoRestrictionCustomizationDescriptor instead')
const GeoRestrictionCustomization$json = {
  '1': 'GeoRestrictionCustomization',
  '2': [
    {'1': 'locations', '3': 85922800, '4': 3, '5': 9, '10': 'locations'},
    {
      '1': 'restrictiontype',
      '3': 266540262,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.GeoRestrictionType',
      '10': 'restrictiontype'
    },
  ],
};

/// Descriptor for `GeoRestrictionCustomization`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List geoRestrictionCustomizationDescriptor =
    $convert.base64Decode(
        'ChtHZW9SZXN0cmljdGlvbkN1c3RvbWl6YXRpb24SHwoJbG9jYXRpb25zGPCn/CggAygJUglsb2'
        'NhdGlvbnMSSwoPcmVzdHJpY3Rpb250eXBlGOapjH8gASgOMh4uY2xvdWRmcm9udC5HZW9SZXN0'
        'cmljdGlvblR5cGVSD3Jlc3RyaWN0aW9udHlwZQ==');

@$core.Deprecated('Use getAnycastIpListRequestDescriptor instead')
const GetAnycastIpListRequest$json = {
  '1': 'GetAnycastIpListRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetAnycastIpListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAnycastIpListRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRBbnljYXN0SXBMaXN0UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use getAnycastIpListResultDescriptor instead')
const GetAnycastIpListResult$json = {
  '1': 'GetAnycastIpListResult',
  '2': [
    {
      '1': 'anycastiplist',
      '3': 190550768,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.AnycastIpList',
      '8': {},
      '10': 'anycastiplist'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetAnycastIpListResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAnycastIpListResultDescriptor = $convert.base64Decode(
    'ChZHZXRBbnljYXN0SXBMaXN0UmVzdWx0EkgKDWFueWNhc3RpcGxpc3QY8KXuWiABKAsyGS5jbG'
    '91ZGZyb250LkFueWNhc3RJcExpc3RCBIi1GAFSDWFueWNhc3RpcGxpc3QSFgoEZXRhZxiB37OV'
    'ASABKAlSBGV0YWc=');

@$core.Deprecated('Use getCachePolicyConfigRequestDescriptor instead')
const GetCachePolicyConfigRequest$json = {
  '1': 'GetCachePolicyConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetCachePolicyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCachePolicyConfigRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRDYWNoZVBvbGljeUNvbmZpZ1JlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use getCachePolicyConfigResultDescriptor instead')
const GetCachePolicyConfigResult$json = {
  '1': 'GetCachePolicyConfigResult',
  '2': [
    {
      '1': 'cachepolicyconfig',
      '3': 407094126,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyConfig',
      '8': {},
      '10': 'cachepolicyconfig'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetCachePolicyConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCachePolicyConfigResultDescriptor =
    $convert.base64Decode(
        'ChpHZXRDYWNoZVBvbGljeUNvbmZpZ1Jlc3VsdBJVChFjYWNoZXBvbGljeWNvbmZpZxjuho/CAS'
        'ABKAsyHS5jbG91ZGZyb250LkNhY2hlUG9saWN5Q29uZmlnQgSItRgBUhFjYWNoZXBvbGljeWNv'
        'bmZpZxIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw==');

@$core.Deprecated('Use getCachePolicyRequestDescriptor instead')
const GetCachePolicyRequest$json = {
  '1': 'GetCachePolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetCachePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCachePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChVHZXRDYWNoZVBvbGljeVJlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use getCachePolicyResultDescriptor instead')
const GetCachePolicyResult$json = {
  '1': 'GetCachePolicyResult',
  '2': [
    {
      '1': 'cachepolicy',
      '3': 439848032,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicy',
      '8': {},
      '10': 'cachepolicy'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetCachePolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCachePolicyResultDescriptor = $convert.base64Decode(
    'ChRHZXRDYWNoZVBvbGljeVJlc3VsdBJDCgtjYWNoZXBvbGljeRjgmN7RASABKAsyFy5jbG91ZG'
    'Zyb250LkNhY2hlUG9saWN5QgSItRgBUgtjYWNoZXBvbGljeRIWCgRldGFnGIHfs5UBIAEoCVIE'
    'ZXRhZw==');

@$core.Deprecated(
    'Use getCloudFrontOriginAccessIdentityConfigRequestDescriptor instead')
const GetCloudFrontOriginAccessIdentityConfigRequest$json = {
  '1': 'GetCloudFrontOriginAccessIdentityConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetCloudFrontOriginAccessIdentityConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getCloudFrontOriginAccessIdentityConfigRequestDescriptor =
    $convert.base64Decode(
        'Ci5HZXRDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlDb25maWdSZXF1ZXN0EhIKAmlkGI'
        'HyorcBIAEoCVICaWQ=');

@$core.Deprecated(
    'Use getCloudFrontOriginAccessIdentityConfigResultDescriptor instead')
const GetCloudFrontOriginAccessIdentityConfigResult$json = {
  '1': 'GetCloudFrontOriginAccessIdentityConfigResult',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentityconfig',
      '3': 111945038,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentityConfig',
      '8': {},
      '10': 'cloudfrontoriginaccessidentityconfig'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetCloudFrontOriginAccessIdentityConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getCloudFrontOriginAccessIdentityConfigResultDescriptor =
    $convert.base64Decode(
        'Ci1HZXRDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlDb25maWdSZXN1bHQSjQEKJGNsb3'
        'VkZnJvbnRvcmlnaW5hY2Nlc3NpZGVudGl0eWNvbmZpZxjOyrA1IAEoCzIwLmNsb3VkZnJvbnQu'
        'Q2xvdWRGcm9udE9yaWdpbkFjY2Vzc0lkZW50aXR5Q29uZmlnQgSItRgBUiRjbG91ZGZyb250b3'
        'JpZ2luYWNjZXNzaWRlbnRpdHljb25maWcSFgoEZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core.Deprecated(
    'Use getCloudFrontOriginAccessIdentityRequestDescriptor instead')
const GetCloudFrontOriginAccessIdentityRequest$json = {
  '1': 'GetCloudFrontOriginAccessIdentityRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetCloudFrontOriginAccessIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCloudFrontOriginAccessIdentityRequestDescriptor =
    $convert.base64Decode(
        'CihHZXRDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlSZXF1ZXN0EhIKAmlkGIHyorcBIA'
        'EoCVICaWQ=');

@$core
    .Deprecated('Use getCloudFrontOriginAccessIdentityResultDescriptor instead')
const GetCloudFrontOriginAccessIdentityResult$json = {
  '1': 'GetCloudFrontOriginAccessIdentityResult',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentity',
      '3': 109497984,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentity',
      '8': {},
      '10': 'cloudfrontoriginaccessidentity'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetCloudFrontOriginAccessIdentityResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCloudFrontOriginAccessIdentityResultDescriptor =
    $convert.base64Decode(
        'CidHZXRDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlSZXN1bHQSewoeY2xvdWRmcm9udG'
        '9yaWdpbmFjY2Vzc2lkZW50aXR5GICdmzQgASgLMiouY2xvdWRmcm9udC5DbG91ZEZyb250T3Jp'
        'Z2luQWNjZXNzSWRlbnRpdHlCBIi1GAFSHmNsb3VkZnJvbnRvcmlnaW5hY2Nlc3NpZGVudGl0eR'
        'IWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw==');

@$core.Deprecated('Use getConnectionFunctionRequestDescriptor instead')
const GetConnectionFunctionRequest$json = {
  '1': 'GetConnectionFunctionRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `GetConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConnectionFunctionRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRDb25uZWN0aW9uRnVuY3Rpb25SZXF1ZXN0EiEKCmlkZW50aWZpZXIY36D7EyABKAlSCm'
        'lkZW50aWZpZXISMgoFc3RhZ2UYzpfYcCABKA4yGS5jbG91ZGZyb250LkZ1bmN0aW9uU3RhZ2VS'
        'BXN0YWdl');

@$core.Deprecated('Use getConnectionFunctionResultDescriptor instead')
const GetConnectionFunctionResult$json = {
  '1': 'GetConnectionFunctionResult',
  '2': [
    {
      '1': 'connectionfunctioncode',
      '3': 502949501,
      '4': 1,
      '5': 12,
      '8': {},
      '10': 'connectionfunctioncode'
    },
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetConnectionFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConnectionFunctionResultDescriptor =
    $convert.base64Decode(
        'ChtHZXRDb25uZWN0aW9uRnVuY3Rpb25SZXN1bHQSQAoWY29ubmVjdGlvbmZ1bmN0aW9uY29kZR'
        'j9zOnvASABKAxCBIi1GAFSFmNvbm5lY3Rpb25mdW5jdGlvbmNvZGUSJAoLY29udGVudHR5cGUY'
        'k9XongEgASgJUgtjb250ZW50dHlwZRIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw==');

@$core.Deprecated(
    'Use getConnectionGroupByRoutingEndpointRequestDescriptor instead')
const GetConnectionGroupByRoutingEndpointRequest$json = {
  '1': 'GetConnectionGroupByRoutingEndpointRequest',
  '2': [
    {
      '1': 'routingendpoint',
      '3': 409864193,
      '4': 1,
      '5': 9,
      '10': 'routingendpoint'
    },
  ],
};

/// Descriptor for `GetConnectionGroupByRoutingEndpointRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getConnectionGroupByRoutingEndpointRequestDescriptor =
    $convert.base64Decode(
        'CipHZXRDb25uZWN0aW9uR3JvdXBCeVJvdXRpbmdFbmRwb2ludFJlcXVlc3QSLAoPcm91dGluZ2'
        'VuZHBvaW50GIGQuMMBIAEoCVIPcm91dGluZ2VuZHBvaW50');

@$core.Deprecated(
    'Use getConnectionGroupByRoutingEndpointResultDescriptor instead')
const GetConnectionGroupByRoutingEndpointResult$json = {
  '1': 'GetConnectionGroupByRoutingEndpointResult',
  '2': [
    {
      '1': 'connectiongroup',
      '3': 517217105,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionGroup',
      '8': {},
      '10': 'connectiongroup'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetConnectionGroupByRoutingEndpointResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getConnectionGroupByRoutingEndpointResultDescriptor = $convert.base64Decode(
        'CilHZXRDb25uZWN0aW9uR3JvdXBCeVJvdXRpbmdFbmRwb2ludFJlc3VsdBJPCg9jb25uZWN0aW'
        '9uZ3JvdXAY0bbQ9gEgASgLMhsuY2xvdWRmcm9udC5Db25uZWN0aW9uR3JvdXBCBIi1GAFSD2Nv'
        'bm5lY3Rpb25ncm91cBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw==');

@$core.Deprecated('Use getConnectionGroupRequestDescriptor instead')
const GetConnectionGroupRequest$json = {
  '1': 'GetConnectionGroupRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `GetConnectionGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConnectionGroupRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRDb25uZWN0aW9uR3JvdXBSZXF1ZXN0EiEKCmlkZW50aWZpZXIY36D7EyABKAlSCmlkZW'
        '50aWZpZXI=');

@$core.Deprecated('Use getConnectionGroupResultDescriptor instead')
const GetConnectionGroupResult$json = {
  '1': 'GetConnectionGroupResult',
  '2': [
    {
      '1': 'connectiongroup',
      '3': 517217105,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionGroup',
      '8': {},
      '10': 'connectiongroup'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetConnectionGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConnectionGroupResultDescriptor = $convert.base64Decode(
    'ChhHZXRDb25uZWN0aW9uR3JvdXBSZXN1bHQSTwoPY29ubmVjdGlvbmdyb3VwGNG20PYBIAEoCz'
    'IbLmNsb3VkZnJvbnQuQ29ubmVjdGlvbkdyb3VwQgSItRgBUg9jb25uZWN0aW9uZ3JvdXASFgoE'
    'ZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core.Deprecated(
    'Use getContinuousDeploymentPolicyConfigRequestDescriptor instead')
const GetContinuousDeploymentPolicyConfigRequest$json = {
  '1': 'GetContinuousDeploymentPolicyConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetContinuousDeploymentPolicyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getContinuousDeploymentPolicyConfigRequestDescriptor =
    $convert.base64Decode(
        'CipHZXRDb250aW51b3VzRGVwbG95bWVudFBvbGljeUNvbmZpZ1JlcXVlc3QSEgoCaWQYgfKitw'
        'EgASgJUgJpZA==');

@$core.Deprecated(
    'Use getContinuousDeploymentPolicyConfigResultDescriptor instead')
const GetContinuousDeploymentPolicyConfigResult$json = {
  '1': 'GetContinuousDeploymentPolicyConfigResult',
  '2': [
    {
      '1': 'continuousdeploymentpolicyconfig',
      '3': 161949042,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicyConfig',
      '8': {},
      '10': 'continuousdeploymentpolicyconfig'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetContinuousDeploymentPolicyConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getContinuousDeploymentPolicyConfigResultDescriptor = $convert.base64Decode(
        'CilHZXRDb250aW51b3VzRGVwbG95bWVudFBvbGljeUNvbmZpZ1Jlc3VsdBKBAQogY29udGludW'
        '91c2RlcGxveW1lbnRwb2xpY3ljb25maWcY8sqcTSABKAsyLC5jbG91ZGZyb250LkNvbnRpbnVv'
        'dXNEZXBsb3ltZW50UG9saWN5Q29uZmlnQgSItRgBUiBjb250aW51b3VzZGVwbG95bWVudHBvbG'
        'ljeWNvbmZpZxIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw==');

@$core.Deprecated('Use getContinuousDeploymentPolicyRequestDescriptor instead')
const GetContinuousDeploymentPolicyRequest$json = {
  '1': 'GetContinuousDeploymentPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetContinuousDeploymentPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContinuousDeploymentPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRDb250aW51b3VzRGVwbG95bWVudFBvbGljeVJlcXVlc3QSEgoCaWQYgfKitwEgASgJUg'
        'JpZA==');

@$core.Deprecated('Use getContinuousDeploymentPolicyResultDescriptor instead')
const GetContinuousDeploymentPolicyResult$json = {
  '1': 'GetContinuousDeploymentPolicyResult',
  '2': [
    {
      '1': 'continuousdeploymentpolicy',
      '3': 36616788,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicy',
      '8': {},
      '10': 'continuousdeploymentpolicy'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetContinuousDeploymentPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContinuousDeploymentPolicyResultDescriptor =
    $convert.base64Decode(
        'CiNHZXRDb250aW51b3VzRGVwbG95bWVudFBvbGljeVJlc3VsdBJvChpjb250aW51b3VzZGVwbG'
        '95bWVudHBvbGljeRjU9LoRIAEoCzImLmNsb3VkZnJvbnQuQ29udGludW91c0RlcGxveW1lbnRQ'
        'b2xpY3lCBIi1GAFSGmNvbnRpbnVvdXNkZXBsb3ltZW50cG9saWN5EhYKBGV0YWcYgd+zlQEgAS'
        'gJUgRldGFn');

@$core.Deprecated('Use getDistributionConfigRequestDescriptor instead')
const GetDistributionConfigRequest$json = {
  '1': 'GetDistributionConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetDistributionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionConfigRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXREaXN0cmlidXRpb25Db25maWdSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getDistributionConfigResultDescriptor instead')
const GetDistributionConfigResult$json = {
  '1': 'GetDistributionConfigResult',
  '2': [
    {
      '1': 'distributionconfig',
      '3': 528940762,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionConfig',
      '8': {},
      '10': 'distributionconfig'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetDistributionConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionConfigResultDescriptor =
    $convert.base64Decode(
        'ChtHZXREaXN0cmlidXRpb25Db25maWdSZXN1bHQSWAoSZGlzdHJpYnV0aW9uY29uZmlnGNr9m/'
        'wBIAEoCzIeLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uQ29uZmlnQgSItRgBUhJkaXN0cmlidXRp'
        'b25jb25maWcSFgoEZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core.Deprecated('Use getDistributionRequestDescriptor instead')
const GetDistributionRequest$json = {
  '1': 'GetDistributionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXREaXN0cmlidXRpb25SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getDistributionResultDescriptor instead')
const GetDistributionResult$json = {
  '1': 'GetDistributionResult',
  '2': [
    {
      '1': 'distribution',
      '3': 105183308,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Distribution',
      '8': {},
      '10': 'distribution'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetDistributionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionResultDescriptor = $convert.base64Decode(
    'ChVHZXREaXN0cmlidXRpb25SZXN1bHQSRQoMZGlzdHJpYnV0aW9uGMzwkzIgASgLMhguY2xvdW'
    'Rmcm9udC5EaXN0cmlidXRpb25CBIi1GAFSDGRpc3RyaWJ1dGlvbhIWCgRldGFnGIHfs5UBIAEo'
    'CVIEZXRhZw==');

@$core.Deprecated('Use getDistributionTenantByDomainRequestDescriptor instead')
const GetDistributionTenantByDomainRequest$json = {
  '1': 'GetDistributionTenantByDomainRequest',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
  ],
};

/// Descriptor for `GetDistributionTenantByDomainRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionTenantByDomainRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXREaXN0cmlidXRpb25UZW5hbnRCeURvbWFpblJlcXVlc3QSGgoGZG9tYWluGJKS8vABIA'
        'EoCVIGZG9tYWlu');

@$core.Deprecated('Use getDistributionTenantByDomainResultDescriptor instead')
const GetDistributionTenantByDomainResult$json = {
  '1': 'GetDistributionTenantByDomainResult',
  '2': [
    {
      '1': 'distributiontenant',
      '3': 510856916,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionTenant',
      '8': {},
      '10': 'distributiontenant'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetDistributionTenantByDomainResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionTenantByDomainResultDescriptor =
    $convert.base64Decode(
        'CiNHZXREaXN0cmlidXRpb25UZW5hbnRCeURvbWFpblJlc3VsdBJYChJkaXN0cmlidXRpb250ZW'
        '5hbnQY1J3M8wEgASgLMh4uY2xvdWRmcm9udC5EaXN0cmlidXRpb25UZW5hbnRCBIi1GAFSEmRp'
        'c3RyaWJ1dGlvbnRlbmFudBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw==');

@$core.Deprecated('Use getDistributionTenantRequestDescriptor instead')
const GetDistributionTenantRequest$json = {
  '1': 'GetDistributionTenantRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `GetDistributionTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionTenantRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXREaXN0cmlidXRpb25UZW5hbnRSZXF1ZXN0EiEKCmlkZW50aWZpZXIY36D7EyABKAlSCm'
        'lkZW50aWZpZXI=');

@$core.Deprecated('Use getDistributionTenantResultDescriptor instead')
const GetDistributionTenantResult$json = {
  '1': 'GetDistributionTenantResult',
  '2': [
    {
      '1': 'distributiontenant',
      '3': 510856916,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionTenant',
      '8': {},
      '10': 'distributiontenant'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `GetDistributionTenantResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDistributionTenantResultDescriptor =
    $convert.base64Decode(
        'ChtHZXREaXN0cmlidXRpb25UZW5hbnRSZXN1bHQSWAoSZGlzdHJpYnV0aW9udGVuYW50GNSdzP'
        'MBIAEoCzIeLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uVGVuYW50QgSItRgBUhJkaXN0cmlidXRp'
        'b250ZW5hbnQSFgoEZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core.Deprecated('Use getFieldLevelEncryptionConfigRequestDescriptor instead')
const GetFieldLevelEncryptionConfigRequest$json = {
  '1': 'GetFieldLevelEncryptionConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetFieldLevelEncryptionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFieldLevelEncryptionConfigRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ1JlcXVlc3QSEgoCaWQYgfKitwEgASgJUg'
        'JpZA==');

@$core.Deprecated('Use getFieldLevelEncryptionConfigResultDescriptor instead')
const GetFieldLevelEncryptionConfigResult$json = {
  '1': 'GetFieldLevelEncryptionConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryptionconfig',
      '3': 499294709,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionConfig',
      '8': {},
      '10': 'fieldlevelencryptionconfig'
    },
  ],
};

/// Descriptor for `GetFieldLevelEncryptionConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFieldLevelEncryptionConfigResultDescriptor =
    $convert.base64Decode(
        'CiNHZXRGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ1Jlc3VsdBIWCgRldGFnGIHfs5UBIAEoCV'
        'IEZXRhZxJwChpmaWVsZGxldmVsZW5jcnlwdGlvbmNvbmZpZxj1w4ruASABKAsyJi5jbG91ZGZy'
        'b250LkZpZWxkTGV2ZWxFbmNyeXB0aW9uQ29uZmlnQgSItRgBUhpmaWVsZGxldmVsZW5jcnlwdG'
        'lvbmNvbmZpZw==');

@$core.Deprecated(
    'Use getFieldLevelEncryptionProfileConfigRequestDescriptor instead')
const GetFieldLevelEncryptionProfileConfigRequest$json = {
  '1': 'GetFieldLevelEncryptionProfileConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetFieldLevelEncryptionProfileConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getFieldLevelEncryptionProfileConfigRequestDescriptor =
    $convert.base64Decode(
        'CitHZXRGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVDb25maWdSZXF1ZXN0EhIKAmlkGIHyor'
        'cBIAEoCVICaWQ=');

@$core.Deprecated(
    'Use getFieldLevelEncryptionProfileConfigResultDescriptor instead')
const GetFieldLevelEncryptionProfileConfigResult$json = {
  '1': 'GetFieldLevelEncryptionProfileConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryptionprofileconfig',
      '3': 199371734,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfileConfig',
      '8': {},
      '10': 'fieldlevelencryptionprofileconfig'
    },
  ],
};

/// Descriptor for `GetFieldLevelEncryptionProfileConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getFieldLevelEncryptionProfileConfigResultDescriptor =
    $convert.base64Decode(
        'CipHZXRGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVDb25maWdSZXN1bHQSFgoEZXRhZxiB37'
        'OVASABKAlSBGV0YWcShAEKIWZpZWxkbGV2ZWxlbmNyeXB0aW9ucHJvZmlsZWNvbmZpZxjW14hf'
        'IAEoCzItLmNsb3VkZnJvbnQuRmllbGRMZXZlbEVuY3J5cHRpb25Qcm9maWxlQ29uZmlnQgSItR'
        'gBUiFmaWVsZGxldmVsZW5jcnlwdGlvbnByb2ZpbGVjb25maWc=');

@$core.Deprecated('Use getFieldLevelEncryptionProfileRequestDescriptor instead')
const GetFieldLevelEncryptionProfileRequest$json = {
  '1': 'GetFieldLevelEncryptionProfileRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetFieldLevelEncryptionProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFieldLevelEncryptionProfileRequestDescriptor =
    $convert.base64Decode(
        'CiVHZXRGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCV'
        'ICaWQ=');

@$core.Deprecated('Use getFieldLevelEncryptionProfileResultDescriptor instead')
const GetFieldLevelEncryptionProfileResult$json = {
  '1': 'GetFieldLevelEncryptionProfileResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryptionprofile',
      '3': 344546136,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfile',
      '8': {},
      '10': 'fieldlevelencryptionprofile'
    },
  ],
};

/// Descriptor for `GetFieldLevelEncryptionProfileResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFieldLevelEncryptionProfileResultDescriptor =
    $convert.base64Decode(
        'CiRHZXRGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVSZXN1bHQSFgoEZXRhZxiB37OVASABKA'
        'lSBGV0YWcScwobZmllbGRsZXZlbGVuY3J5cHRpb25wcm9maWxlGNi2paQBIAEoCzInLmNsb3Vk'
        'ZnJvbnQuRmllbGRMZXZlbEVuY3J5cHRpb25Qcm9maWxlQgSItRgBUhtmaWVsZGxldmVsZW5jcn'
        'lwdGlvbnByb2ZpbGU=');

@$core.Deprecated('Use getFieldLevelEncryptionRequestDescriptor instead')
const GetFieldLevelEncryptionRequest$json = {
  '1': 'GetFieldLevelEncryptionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetFieldLevelEncryptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFieldLevelEncryptionRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXRGaWVsZExldmVsRW5jcnlwdGlvblJlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use getFieldLevelEncryptionResultDescriptor instead')
const GetFieldLevelEncryptionResult$json = {
  '1': 'GetFieldLevelEncryptionResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryption',
      '3': 473382747,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryption',
      '8': {},
      '10': 'fieldlevelencryption'
    },
  ],
};

/// Descriptor for `GetFieldLevelEncryptionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFieldLevelEncryptionResultDescriptor =
    $convert.base64Decode(
        'Ch1HZXRGaWVsZExldmVsRW5jcnlwdGlvblJlc3VsdBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZx'
        'JeChRmaWVsZGxldmVsZW5jcnlwdGlvbhjb/tzhASABKAsyIC5jbG91ZGZyb250LkZpZWxkTGV2'
        'ZWxFbmNyeXB0aW9uQgSItRgBUhRmaWVsZGxldmVsZW5jcnlwdGlvbg==');

@$core.Deprecated('Use getFunctionRequestDescriptor instead')
const GetFunctionRequest$json = {
  '1': 'GetFunctionRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `GetFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionRequestDescriptor = $convert.base64Decode(
    'ChJHZXRGdW5jdGlvblJlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIyCgVzdGFnZRjOl9'
    'hwIAEoDjIZLmNsb3VkZnJvbnQuRnVuY3Rpb25TdGFnZVIFc3RhZ2U=');

@$core.Deprecated('Use getFunctionResultDescriptor instead')
const GetFunctionResult$json = {
  '1': 'GetFunctionResult',
  '2': [
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'functioncode',
      '3': 405947809,
      '4': 1,
      '5': 12,
      '8': {},
      '10': 'functioncode'
    },
  ],
};

/// Descriptor for `GetFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionResultDescriptor = $convert.base64Decode(
    'ChFHZXRGdW5jdGlvblJlc3VsdBIkCgtjb250ZW50dHlwZRiT1eieASABKAlSC2NvbnRlbnR0eX'
    'BlEhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEiwKDGZ1bmN0aW9uY29kZRihi8nBASABKAxCBIi1'
    'GAFSDGZ1bmN0aW9uY29kZQ==');

@$core.Deprecated(
    'Use getInvalidationForDistributionTenantRequestDescriptor instead')
const GetInvalidationForDistributionTenantRequest$json = {
  '1': 'GetInvalidationForDistributionTenantRequest',
  '2': [
    {
      '1': 'distributiontenantid',
      '3': 123323327,
      '4': 1,
      '5': 9,
      '10': 'distributiontenantid'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetInvalidationForDistributionTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getInvalidationForDistributionTenantRequestDescriptor =
    $convert.base64Decode(
        'CitHZXRJbnZhbGlkYXRpb25Gb3JEaXN0cmlidXRpb25UZW5hbnRSZXF1ZXN0EjUKFGRpc3RyaW'
        'J1dGlvbnRlbmFudGlkGL+H5zogASgJUhRkaXN0cmlidXRpb250ZW5hbnRpZBISCgJpZBiB8qK3'
        'ASABKAlSAmlk');

@$core.Deprecated(
    'Use getInvalidationForDistributionTenantResultDescriptor instead')
const GetInvalidationForDistributionTenantResult$json = {
  '1': 'GetInvalidationForDistributionTenantResult',
  '2': [
    {
      '1': 'invalidation',
      '3': 77924830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Invalidation',
      '8': {},
      '10': 'invalidation'
    },
  ],
};

/// Descriptor for `GetInvalidationForDistributionTenantResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getInvalidationForDistributionTenantResultDescriptor =
    $convert.base64Decode(
        'CipHZXRJbnZhbGlkYXRpb25Gb3JEaXN0cmlidXRpb25UZW5hbnRSZXN1bHQSRQoMaW52YWxpZG'
        'F0aW9uGN6TlCUgASgLMhguY2xvdWRmcm9udC5JbnZhbGlkYXRpb25CBIi1GAFSDGludmFsaWRh'
        'dGlvbg==');

@$core.Deprecated('Use getInvalidationRequestDescriptor instead')
const GetInvalidationRequest$json = {
  '1': 'GetInvalidationRequest',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetInvalidationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInvalidationRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRJbnZhbGlkYXRpb25SZXF1ZXN0EikKDmRpc3RyaWJ1dGlvbmlkGOex+0MgASgJUg5kaX'
        'N0cmlidXRpb25pZBISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use getInvalidationResultDescriptor instead')
const GetInvalidationResult$json = {
  '1': 'GetInvalidationResult',
  '2': [
    {
      '1': 'invalidation',
      '3': 77924830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Invalidation',
      '8': {},
      '10': 'invalidation'
    },
  ],
};

/// Descriptor for `GetInvalidationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInvalidationResultDescriptor = $convert.base64Decode(
    'ChVHZXRJbnZhbGlkYXRpb25SZXN1bHQSRQoMaW52YWxpZGF0aW9uGN6TlCUgASgLMhguY2xvdW'
    'Rmcm9udC5JbnZhbGlkYXRpb25CBIi1GAFSDGludmFsaWRhdGlvbg==');

@$core.Deprecated('Use getKeyGroupConfigRequestDescriptor instead')
const GetKeyGroupConfigRequest$json = {
  '1': 'GetKeyGroupConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetKeyGroupConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyGroupConfigRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRLZXlHcm91cENvbmZpZ1JlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use getKeyGroupConfigResultDescriptor instead')
const GetKeyGroupConfigResult$json = {
  '1': 'GetKeyGroupConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'keygroupconfig',
      '3': 143012494,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroupConfig',
      '8': {},
      '10': 'keygroupconfig'
    },
  ],
};

/// Descriptor for `GetKeyGroupConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyGroupConfigResultDescriptor = $convert.base64Decode(
    'ChdHZXRLZXlHcm91cENvbmZpZ1Jlc3VsdBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxJLCg5rZX'
    'lncm91cGNvbmZpZxiO5ZhEIAEoCzIaLmNsb3VkZnJvbnQuS2V5R3JvdXBDb25maWdCBIi1GAFS'
    'DmtleWdyb3VwY29uZmln');

@$core.Deprecated('Use getKeyGroupRequestDescriptor instead')
const GetKeyGroupRequest$json = {
  '1': 'GetKeyGroupRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetKeyGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyGroupRequestDescriptor = $convert
    .base64Decode('ChJHZXRLZXlHcm91cFJlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use getKeyGroupResultDescriptor instead')
const GetKeyGroupResult$json = {
  '1': 'GetKeyGroupResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'keygroup',
      '3': 518748096,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroup',
      '8': {},
      '10': 'keygroup'
    },
  ],
};

/// Descriptor for `GetKeyGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyGroupResultDescriptor = $convert.base64Decode(
    'ChFHZXRLZXlHcm91cFJlc3VsdBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxI6CghrZXlncm91cB'
    'jA7633ASABKAsyFC5jbG91ZGZyb250LktleUdyb3VwQgSItRgBUghrZXlncm91cA==');

@$core.Deprecated('Use getManagedCertificateDetailsRequestDescriptor instead')
const GetManagedCertificateDetailsRequest$json = {
  '1': 'GetManagedCertificateDetailsRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `GetManagedCertificateDetailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getManagedCertificateDetailsRequestDescriptor =
    $convert.base64Decode(
        'CiNHZXRNYW5hZ2VkQ2VydGlmaWNhdGVEZXRhaWxzUmVxdWVzdBIhCgppZGVudGlmaWVyGN+g+x'
        'MgASgJUgppZGVudGlmaWVy');

@$core.Deprecated('Use getManagedCertificateDetailsResultDescriptor instead')
const GetManagedCertificateDetailsResult$json = {
  '1': 'GetManagedCertificateDetailsResult',
  '2': [
    {
      '1': 'managedcertificatedetails',
      '3': 517987386,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ManagedCertificateDetails',
      '8': {},
      '10': 'managedcertificatedetails'
    },
  ],
};

/// Descriptor for `GetManagedCertificateDetailsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getManagedCertificateDetailsResultDescriptor =
    $convert.base64Decode(
        'CiJHZXRNYW5hZ2VkQ2VydGlmaWNhdGVEZXRhaWxzUmVzdWx0Em0KGW1hbmFnZWRjZXJ0aWZpY2'
        'F0ZWRldGFpbHMYurj/9gEgASgLMiUuY2xvdWRmcm9udC5NYW5hZ2VkQ2VydGlmaWNhdGVEZXRh'
        'aWxzQgSItRgBUhltYW5hZ2VkY2VydGlmaWNhdGVkZXRhaWxz');

@$core.Deprecated('Use getMonitoringSubscriptionRequestDescriptor instead')
const GetMonitoringSubscriptionRequest$json = {
  '1': 'GetMonitoringSubscriptionRequest',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
  ],
};

/// Descriptor for `GetMonitoringSubscriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMonitoringSubscriptionRequestDescriptor =
    $convert.base64Decode(
        'CiBHZXRNb25pdG9yaW5nU3Vic2NyaXB0aW9uUmVxdWVzdBIpCg5kaXN0cmlidXRpb25pZBjnsf'
        'tDIAEoCVIOZGlzdHJpYnV0aW9uaWQ=');

@$core.Deprecated('Use getMonitoringSubscriptionResultDescriptor instead')
const GetMonitoringSubscriptionResult$json = {
  '1': 'GetMonitoringSubscriptionResult',
  '2': [
    {
      '1': 'monitoringsubscription',
      '3': 456621195,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.MonitoringSubscription',
      '8': {},
      '10': 'monitoringsubscription'
    },
  ],
};

/// Descriptor for `GetMonitoringSubscriptionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMonitoringSubscriptionResultDescriptor =
    $convert.base64Decode(
        'Ch9HZXRNb25pdG9yaW5nU3Vic2NyaXB0aW9uUmVzdWx0EmQKFm1vbml0b3JpbmdzdWJzY3JpcH'
        'Rpb24Yi/nd2QEgASgLMiIuY2xvdWRmcm9udC5Nb25pdG9yaW5nU3Vic2NyaXB0aW9uQgSItRgB'
        'UhZtb25pdG9yaW5nc3Vic2NyaXB0aW9u');

@$core.Deprecated('Use getOriginAccessControlConfigRequestDescriptor instead')
const GetOriginAccessControlConfigRequest$json = {
  '1': 'GetOriginAccessControlConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetOriginAccessControlConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginAccessControlConfigRequestDescriptor =
    $convert.base64Decode(
        'CiNHZXRPcmlnaW5BY2Nlc3NDb250cm9sQ29uZmlnUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAm'
        'lk');

@$core.Deprecated('Use getOriginAccessControlConfigResultDescriptor instead')
const GetOriginAccessControlConfigResult$json = {
  '1': 'GetOriginAccessControlConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'originaccesscontrolconfig',
      '3': 143834977,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControlConfig',
      '8': {},
      '10': 'originaccesscontrolconfig'
    },
  ],
};

/// Descriptor for `GetOriginAccessControlConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginAccessControlConfigResultDescriptor =
    $convert.base64Decode(
        'CiJHZXRPcmlnaW5BY2Nlc3NDb250cm9sQ29uZmlnUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUg'
        'RldGFnEmwKGW9yaWdpbmFjY2Vzc2NvbnRyb2xjb25maWcY4f7KRCABKAsyJS5jbG91ZGZyb250'
        'Lk9yaWdpbkFjY2Vzc0NvbnRyb2xDb25maWdCBIi1GAFSGW9yaWdpbmFjY2Vzc2NvbnRyb2xjb2'
        '5maWc=');

@$core.Deprecated('Use getOriginAccessControlRequestDescriptor instead')
const GetOriginAccessControlRequest$json = {
  '1': 'GetOriginAccessControlRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetOriginAccessControlRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginAccessControlRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRPcmlnaW5BY2Nlc3NDb250cm9sUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use getOriginAccessControlResultDescriptor instead')
const GetOriginAccessControlResult$json = {
  '1': 'GetOriginAccessControlResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'originaccesscontrol',
      '3': 238302375,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControl',
      '8': {},
      '10': 'originaccesscontrol'
    },
  ],
};

/// Descriptor for `GetOriginAccessControlResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginAccessControlResultDescriptor =
    $convert.base64Decode(
        'ChxHZXRPcmlnaW5BY2Nlc3NDb250cm9sUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEl'
        'oKE29yaWdpbmFjY2Vzc2NvbnRyb2wYp+nQcSABKAsyHy5jbG91ZGZyb250Lk9yaWdpbkFjY2Vz'
        'c0NvbnRyb2xCBIi1GAFSE29yaWdpbmFjY2Vzc2NvbnRyb2w=');

@$core.Deprecated('Use getOriginRequestPolicyConfigRequestDescriptor instead')
const GetOriginRequestPolicyConfigRequest$json = {
  '1': 'GetOriginRequestPolicyConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetOriginRequestPolicyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginRequestPolicyConfigRequestDescriptor =
    $convert.base64Decode(
        'CiNHZXRPcmlnaW5SZXF1ZXN0UG9saWN5Q29uZmlnUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAm'
        'lk');

@$core.Deprecated('Use getOriginRequestPolicyConfigResultDescriptor instead')
const GetOriginRequestPolicyConfigResult$json = {
  '1': 'GetOriginRequestPolicyConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'originrequestpolicyconfig',
      '3': 37078133,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyConfig',
      '8': {},
      '10': 'originrequestpolicyconfig'
    },
  ],
};

/// Descriptor for `GetOriginRequestPolicyConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginRequestPolicyConfigResultDescriptor =
    $convert.base64Decode(
        'CiJHZXRPcmlnaW5SZXF1ZXN0UG9saWN5Q29uZmlnUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUg'
        'RldGFnEmwKGW9yaWdpbnJlcXVlc3Rwb2xpY3ljb25maWcY9YjXESABKAsyJS5jbG91ZGZyb250'
        'Lk9yaWdpblJlcXVlc3RQb2xpY3lDb25maWdCBIi1GAFSGW9yaWdpbnJlcXVlc3Rwb2xpY3ljb2'
        '5maWc=');

@$core.Deprecated('Use getOriginRequestPolicyRequestDescriptor instead')
const GetOriginRequestPolicyRequest$json = {
  '1': 'GetOriginRequestPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetOriginRequestPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginRequestPolicyRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRPcmlnaW5SZXF1ZXN0UG9saWN5UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use getOriginRequestPolicyResultDescriptor instead')
const GetOriginRequestPolicyResult$json = {
  '1': 'GetOriginRequestPolicyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'originrequestpolicy',
      '3': 386733531,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicy',
      '8': {},
      '10': 'originrequestpolicy'
    },
  ],
};

/// Descriptor for `GetOriginRequestPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOriginRequestPolicyResultDescriptor =
    $convert.base64Decode(
        'ChxHZXRPcmlnaW5SZXF1ZXN0UG9saWN5UmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEl'
        'sKE29yaWdpbnJlcXVlc3Rwb2xpY3kY26u0uAEgASgLMh8uY2xvdWRmcm9udC5PcmlnaW5SZXF1'
        'ZXN0UG9saWN5QgSItRgBUhNvcmlnaW5yZXF1ZXN0cG9saWN5');

@$core.Deprecated('Use getPublicKeyConfigRequestDescriptor instead')
const GetPublicKeyConfigRequest$json = {
  '1': 'GetPublicKeyConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetPublicKeyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyConfigRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRQdWJsaWNLZXlDb25maWdSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getPublicKeyConfigResultDescriptor instead')
const GetPublicKeyConfigResult$json = {
  '1': 'GetPublicKeyConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'publickeyconfig',
      '3': 228537966,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKeyConfig',
      '8': {},
      '10': 'publickeyconfig'
    },
  ],
};

/// Descriptor for `GetPublicKeyConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyConfigResultDescriptor = $convert.base64Decode(
    'ChhHZXRQdWJsaWNLZXlDb25maWdSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSTgoPcH'
    'VibGlja2V5Y29uZmlnGO7s/GwgASgLMhsuY2xvdWRmcm9udC5QdWJsaWNLZXlDb25maWdCBIi1'
    'GAFSD3B1YmxpY2tleWNvbmZpZw==');

@$core.Deprecated('Use getPublicKeyRequestDescriptor instead')
const GetPublicKeyRequest$json = {
  '1': 'GetPublicKeyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyRequestDescriptor = $convert
    .base64Decode('ChNHZXRQdWJsaWNLZXlSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getPublicKeyResultDescriptor instead')
const GetPublicKeyResult$json = {
  '1': 'GetPublicKeyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'publickey',
      '3': 167335776,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKey',
      '8': {},
      '10': 'publickey'
    },
  ],
};

/// Descriptor for `GetPublicKeyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyResultDescriptor = $convert.base64Decode(
    'ChJHZXRQdWJsaWNLZXlSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSPAoJcHVibGlja2'
    'V5GOCu5U8gASgLMhUuY2xvdWRmcm9udC5QdWJsaWNLZXlCBIi1GAFSCXB1YmxpY2tleQ==');

@$core.Deprecated('Use getRealtimeLogConfigRequestDescriptor instead')
const GetRealtimeLogConfigRequest$json = {
  '1': 'GetRealtimeLogConfigRequest',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetRealtimeLogConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRealtimeLogConfigRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRSZWFsdGltZUxvZ0NvbmZpZ1JlcXVlc3QSFAoDYXJuGJ2cr70BIAEoCVIDYXJuEhUKBG'
        '5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use getRealtimeLogConfigResultDescriptor instead')
const GetRealtimeLogConfigResult$json = {
  '1': 'GetRealtimeLogConfigResult',
  '2': [
    {
      '1': 'realtimelogconfig',
      '3': 95917609,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.RealtimeLogConfig',
      '10': 'realtimelogconfig'
    },
  ],
};

/// Descriptor for `GetRealtimeLogConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRealtimeLogConfigResultDescriptor =
    $convert.base64Decode(
        'ChpHZXRSZWFsdGltZUxvZ0NvbmZpZ1Jlc3VsdBJOChFyZWFsdGltZWxvZ2NvbmZpZxiprN4tIA'
        'EoCzIdLmNsb3VkZnJvbnQuUmVhbHRpbWVMb2dDb25maWdSEXJlYWx0aW1lbG9nY29uZmln');

@$core.Deprecated('Use getResourcePolicyRequestDescriptor instead')
const GetResourcePolicyRequest$json = {
  '1': 'GetResourcePolicyRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZX'
        'NvdXJjZWFybg==');

@$core.Deprecated('Use getResourcePolicyResultDescriptor instead')
const GetResourcePolicyResult$json = {
  '1': 'GetResourcePolicyResult',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetResourcePolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyResultDescriptor = $convert.base64Decode(
    'ChdHZXRSZXNvdXJjZVBvbGljeVJlc3VsdBIpCg5wb2xpY3lkb2N1bWVudBjLrsFxIAEoCVIOcG'
    '9saWN5ZG9jdW1lbnQSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJjZWFybg==');

@$core.Deprecated('Use getResponseHeadersPolicyConfigRequestDescriptor instead')
const GetResponseHeadersPolicyConfigRequest$json = {
  '1': 'GetResponseHeadersPolicyConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetResponseHeadersPolicyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResponseHeadersPolicyConfigRequestDescriptor =
    $convert.base64Decode(
        'CiVHZXRSZXNwb25zZUhlYWRlcnNQb2xpY3lDb25maWdSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCV'
        'ICaWQ=');

@$core.Deprecated('Use getResponseHeadersPolicyConfigResultDescriptor instead')
const GetResponseHeadersPolicyConfigResult$json = {
  '1': 'GetResponseHeadersPolicyConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'responseheaderspolicyconfig',
      '3': 159056825,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyConfig',
      '8': {},
      '10': 'responseheaderspolicyconfig'
    },
  ],
};

/// Descriptor for `GetResponseHeadersPolicyConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResponseHeadersPolicyConfigResultDescriptor =
    $convert.base64Decode(
        'CiRHZXRSZXNwb25zZUhlYWRlcnNQb2xpY3lDb25maWdSZXN1bHQSFgoEZXRhZxiB37OVASABKA'
        'lSBGV0YWcScgobcmVzcG9uc2VoZWFkZXJzcG9saWN5Y29uZmlnGLmH7EsgASgLMicuY2xvdWRm'
        'cm9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lDb25maWdCBIi1GAFSG3Jlc3BvbnNlaGVhZGVyc3'
        'BvbGljeWNvbmZpZw==');

@$core.Deprecated('Use getResponseHeadersPolicyRequestDescriptor instead')
const GetResponseHeadersPolicyRequest$json = {
  '1': 'GetResponseHeadersPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetResponseHeadersPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResponseHeadersPolicyRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRSZXNwb25zZUhlYWRlcnNQb2xpY3lSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getResponseHeadersPolicyResultDescriptor instead')
const GetResponseHeadersPolicyResult$json = {
  '1': 'GetResponseHeadersPolicyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'responseheaderspolicy',
      '3': 418204719,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicy',
      '8': {},
      '10': 'responseheaderspolicy'
    },
  ],
};

/// Descriptor for `GetResponseHeadersPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResponseHeadersPolicyResultDescriptor =
    $convert.base64Decode(
        'Ch5HZXRSZXNwb25zZUhlYWRlcnNQb2xpY3lSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YW'
        'cSYQoVcmVzcG9uc2VoZWFkZXJzcG9saWN5GK+YtccBIAEoCzIhLmNsb3VkZnJvbnQuUmVzcG9u'
        'c2VIZWFkZXJzUG9saWN5QgSItRgBUhVyZXNwb25zZWhlYWRlcnNwb2xpY3k=');

@$core.Deprecated('Use getStreamingDistributionConfigRequestDescriptor instead')
const GetStreamingDistributionConfigRequest$json = {
  '1': 'GetStreamingDistributionConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetStreamingDistributionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getStreamingDistributionConfigRequestDescriptor =
    $convert.base64Decode(
        'CiVHZXRTdHJlYW1pbmdEaXN0cmlidXRpb25Db25maWdSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCV'
        'ICaWQ=');

@$core.Deprecated('Use getStreamingDistributionConfigResultDescriptor instead')
const GetStreamingDistributionConfigResult$json = {
  '1': 'GetStreamingDistributionConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'streamingdistributionconfig',
      '3': 291115944,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionConfig',
      '8': {},
      '10': 'streamingdistributionconfig'
    },
  ],
};

/// Descriptor for `GetStreamingDistributionConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getStreamingDistributionConfigResultDescriptor =
    $convert.base64Decode(
        'CiRHZXRTdHJlYW1pbmdEaXN0cmlidXRpb25Db25maWdSZXN1bHQSFgoEZXRhZxiB37OVASABKA'
        'lSBGV0YWcScwobc3RyZWFtaW5nZGlzdHJpYnV0aW9uY29uZmlnGKin6IoBIAEoCzInLmNsb3Vk'
        'ZnJvbnQuU3RyZWFtaW5nRGlzdHJpYnV0aW9uQ29uZmlnQgSItRgBUhtzdHJlYW1pbmdkaXN0cm'
        'lidXRpb25jb25maWc=');

@$core.Deprecated('Use getStreamingDistributionRequestDescriptor instead')
const GetStreamingDistributionRequest$json = {
  '1': 'GetStreamingDistributionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetStreamingDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getStreamingDistributionRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRTdHJlYW1pbmdEaXN0cmlidXRpb25SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getStreamingDistributionResultDescriptor instead')
const GetStreamingDistributionResult$json = {
  '1': 'GetStreamingDistributionResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'streamingdistribution',
      '3': 294813830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistribution',
      '8': {},
      '10': 'streamingdistribution'
    },
  ],
};

/// Descriptor for `GetStreamingDistributionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getStreamingDistributionResultDescriptor =
    $convert.base64Decode(
        'Ch5HZXRTdHJlYW1pbmdEaXN0cmlidXRpb25SZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YW'
        'cSYQoVc3RyZWFtaW5nZGlzdHJpYnV0aW9uGIaByowBIAEoCzIhLmNsb3VkZnJvbnQuU3RyZWFt'
        'aW5nRGlzdHJpYnV0aW9uQgSItRgBUhVzdHJlYW1pbmdkaXN0cmlidXRpb24=');

@$core.Deprecated('Use getTrustStoreRequestDescriptor instead')
const GetTrustStoreRequest$json = {
  '1': 'GetTrustStoreRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `GetTrustStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrustStoreRequestDescriptor = $convert.base64Decode(
    'ChRHZXRUcnVzdFN0b3JlUmVxdWVzdBIhCgppZGVudGlmaWVyGN+g+xMgASgJUgppZGVudGlmaW'
    'Vy');

@$core.Deprecated('Use getTrustStoreResultDescriptor instead')
const GetTrustStoreResult$json = {
  '1': 'GetTrustStoreResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'truststore',
      '3': 224815327,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustStore',
      '8': {},
      '10': 'truststore'
    },
  ],
};

/// Descriptor for `GetTrustStoreResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrustStoreResultDescriptor = $convert.base64Decode(
    'ChNHZXRUcnVzdFN0b3JlUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEj8KCnRydXN0c3'
    'RvcmUY39GZayABKAsyFi5jbG91ZGZyb250LlRydXN0U3RvcmVCBIi1GAFSCnRydXN0c3RvcmU=');

@$core.Deprecated('Use getVpcOriginRequestDescriptor instead')
const GetVpcOriginRequest$json = {
  '1': 'GetVpcOriginRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetVpcOriginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getVpcOriginRequestDescriptor = $convert
    .base64Decode('ChNHZXRWcGNPcmlnaW5SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getVpcOriginResultDescriptor instead')
const GetVpcOriginResult$json = {
  '1': 'GetVpcOriginResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'vpcorigin',
      '3': 159181387,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOrigin',
      '8': {},
      '10': 'vpcorigin'
    },
  ],
};

/// Descriptor for `GetVpcOriginResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getVpcOriginResultDescriptor = $convert.base64Decode(
    'ChJHZXRWcGNPcmlnaW5SZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSPAoJdnBjb3JpZ2'
    'luGMvU80sgASgLMhUuY2xvdWRmcm9udC5WcGNPcmlnaW5CBIi1GAFSCXZwY29yaWdpbg==');

@$core.Deprecated('Use grpcConfigDescriptor instead')
const GrpcConfig$json = {
  '1': 'GrpcConfig',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
  ],
};

/// Descriptor for `GrpcConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grpcConfigDescriptor = $convert
    .base64Decode('CgpHcnBjQ29uZmlnEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVk');

@$core.Deprecated('Use headersDescriptor instead')
const Headers$json = {
  '1': 'Headers',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `Headers`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List headersDescriptor = $convert.base64Decode(
    'CgdIZWFkZXJzEhcKBWl0ZW1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBV'
    'IIcXVhbnRpdHk=');

@$core.Deprecated('Use illegalDeleteDescriptor instead')
const IllegalDelete$json = {
  '1': 'IllegalDelete',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IllegalDelete`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List illegalDeleteDescriptor = $convert.base64Decode(
    'Cg1JbGxlZ2FsRGVsZXRlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use illegalFieldLevelEncryptionConfigAssociationWithCacheBehaviorDescriptor instead')
const IllegalFieldLevelEncryptionConfigAssociationWithCacheBehavior$json = {
  '1': 'IllegalFieldLevelEncryptionConfigAssociationWithCacheBehavior',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IllegalFieldLevelEncryptionConfigAssociationWithCacheBehavior`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    illegalFieldLevelEncryptionConfigAssociationWithCacheBehaviorDescriptor =
    $convert.base64Decode(
        'Cj1JbGxlZ2FsRmllbGRMZXZlbEVuY3J5cHRpb25Db25maWdBc3NvY2lhdGlvbldpdGhDYWNoZU'
        'JlaGF2aW9yEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use illegalOriginAccessConfigurationDescriptor instead')
const IllegalOriginAccessConfiguration$json = {
  '1': 'IllegalOriginAccessConfiguration',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IllegalOriginAccessConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List illegalOriginAccessConfigurationDescriptor =
    $convert.base64Decode(
        'CiBJbGxlZ2FsT3JpZ2luQWNjZXNzQ29uZmlndXJhdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use illegalUpdateDescriptor instead')
const IllegalUpdate$json = {
  '1': 'IllegalUpdate',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IllegalUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List illegalUpdateDescriptor = $convert.base64Decode(
    'Cg1JbGxlZ2FsVXBkYXRlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use importSourceDescriptor instead')
const ImportSource$json = {
  '1': 'ImportSource',
  '2': [
    {'1': 'sourcearn', '3': 445113056, '4': 1, '5': 9, '10': 'sourcearn'},
    {
      '1': 'sourcetype',
      '3': 195731217,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ImportSourceType',
      '10': 'sourcetype'
    },
  ],
};

/// Descriptor for `ImportSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importSourceDescriptor = $convert.base64Decode(
    'CgxJbXBvcnRTb3VyY2USIAoJc291cmNlYXJuGODFn9QBIAEoCVIJc291cmNlYXJuEj8KCnNvdX'
    'JjZXR5cGUYkb6qXSABKA4yHC5jbG91ZGZyb250LkltcG9ydFNvdXJjZVR5cGVSCnNvdXJjZXR5'
    'cGU=');

@$core.Deprecated('Use inconsistentQuantitiesDescriptor instead')
const InconsistentQuantities$json = {
  '1': 'InconsistentQuantities',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InconsistentQuantities`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inconsistentQuantitiesDescriptor =
    $convert.base64Decode(
        'ChZJbmNvbnNpc3RlbnRRdWFudGl0aWVzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidArgumentDescriptor instead')
const InvalidArgument$json = {
  '1': 'InvalidArgument',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidArgument`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidArgumentDescriptor = $convert.base64Decode(
    'Cg9JbnZhbGlkQXJndW1lbnQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidAssociationDescriptor instead')
const InvalidAssociation$json = {
  '1': 'InvalidAssociation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAssociationDescriptor =
    $convert.base64Decode(
        'ChJJbnZhbGlkQXNzb2NpYXRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidDefaultRootObjectDescriptor instead')
const InvalidDefaultRootObject$json = {
  '1': 'InvalidDefaultRootObject',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDefaultRootObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDefaultRootObjectDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkRGVmYXVsdFJvb3RPYmplY3QSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core
    .Deprecated('Use invalidDomainNameForOriginAccessControlDescriptor instead')
const InvalidDomainNameForOriginAccessControl$json = {
  '1': 'InvalidDomainNameForOriginAccessControl',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDomainNameForOriginAccessControl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDomainNameForOriginAccessControlDescriptor =
    $convert.base64Decode(
        'CidJbnZhbGlkRG9tYWluTmFtZUZvck9yaWdpbkFjY2Vzc0NvbnRyb2wSGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidErrorCodeDescriptor instead')
const InvalidErrorCode$json = {
  '1': 'InvalidErrorCode',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidErrorCode`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidErrorCodeDescriptor = $convert.base64Decode(
    'ChBJbnZhbGlkRXJyb3JDb2RlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidForwardCookiesDescriptor instead')
const InvalidForwardCookies$json = {
  '1': 'InvalidForwardCookies',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidForwardCookies`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidForwardCookiesDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkRm9yd2FyZENvb2tpZXMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidFunctionAssociationDescriptor instead')
const InvalidFunctionAssociation$json = {
  '1': 'InvalidFunctionAssociation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidFunctionAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidFunctionAssociationDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkRnVuY3Rpb25Bc3NvY2lhdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidGeoRestrictionParameterDescriptor instead')
const InvalidGeoRestrictionParameter$json = {
  '1': 'InvalidGeoRestrictionParameter',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidGeoRestrictionParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidGeoRestrictionParameterDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkR2VvUmVzdHJpY3Rpb25QYXJhbWV0ZXISGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidHeadersForS3OriginDescriptor instead')
const InvalidHeadersForS3Origin$json = {
  '1': 'InvalidHeadersForS3Origin',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidHeadersForS3Origin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidHeadersForS3OriginDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkSGVhZGVyc0ZvclMzT3JpZ2luEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidIfMatchVersionDescriptor instead')
const InvalidIfMatchVersion$json = {
  '1': 'InvalidIfMatchVersion',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidIfMatchVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidIfMatchVersionDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkSWZNYXRjaFZlcnNpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidLambdaFunctionAssociationDescriptor instead')
const InvalidLambdaFunctionAssociation$json = {
  '1': 'InvalidLambdaFunctionAssociation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidLambdaFunctionAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidLambdaFunctionAssociationDescriptor =
    $convert.base64Decode(
        'CiBJbnZhbGlkTGFtYmRhRnVuY3Rpb25Bc3NvY2lhdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use invalidLocationCodeDescriptor instead')
const InvalidLocationCode$json = {
  '1': 'InvalidLocationCode',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidLocationCode`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidLocationCodeDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkTG9jYXRpb25Db2RlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidMinimumProtocolVersionDescriptor instead')
const InvalidMinimumProtocolVersion$json = {
  '1': 'InvalidMinimumProtocolVersion',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidMinimumProtocolVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidMinimumProtocolVersionDescriptor =
    $convert.base64Decode(
        'Ch1JbnZhbGlkTWluaW11bVByb3RvY29sVmVyc2lvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use invalidOriginDescriptor instead')
const InvalidOrigin$json = {
  '1': 'InvalidOrigin',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOrigin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOriginDescriptor = $convert.base64Decode(
    'Cg1JbnZhbGlkT3JpZ2luEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidOriginAccessControlDescriptor instead')
const InvalidOriginAccessControl$json = {
  '1': 'InvalidOriginAccessControl',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOriginAccessControl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOriginAccessControlDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkT3JpZ2luQWNjZXNzQ29udHJvbBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidOriginAccessIdentityDescriptor instead')
const InvalidOriginAccessIdentity$json = {
  '1': 'InvalidOriginAccessIdentity',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOriginAccessIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOriginAccessIdentityDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkT3JpZ2luQWNjZXNzSWRlbnRpdHkSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use invalidOriginKeepaliveTimeoutDescriptor instead')
const InvalidOriginKeepaliveTimeout$json = {
  '1': 'InvalidOriginKeepaliveTimeout',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOriginKeepaliveTimeout`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOriginKeepaliveTimeoutDescriptor =
    $convert.base64Decode(
        'Ch1JbnZhbGlkT3JpZ2luS2VlcGFsaXZlVGltZW91dBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use invalidOriginReadTimeoutDescriptor instead')
const InvalidOriginReadTimeout$json = {
  '1': 'InvalidOriginReadTimeout',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOriginReadTimeout`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOriginReadTimeoutDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkT3JpZ2luUmVhZFRpbWVvdXQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidProtocolSettingsDescriptor instead')
const InvalidProtocolSettings$json = {
  '1': 'InvalidProtocolSettings',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidProtocolSettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidProtocolSettingsDescriptor =
    $convert.base64Decode(
        'ChdJbnZhbGlkUHJvdG9jb2xTZXR0aW5ncxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidQueryStringParametersDescriptor instead')
const InvalidQueryStringParameters$json = {
  '1': 'InvalidQueryStringParameters',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidQueryStringParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidQueryStringParametersDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkUXVlcnlTdHJpbmdQYXJhbWV0ZXJzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidRelativePathDescriptor instead')
const InvalidRelativePath$json = {
  '1': 'InvalidRelativePath',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidRelativePath`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRelativePathDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkUmVsYXRpdmVQYXRoEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidRequiredProtocolDescriptor instead')
const InvalidRequiredProtocol$json = {
  '1': 'InvalidRequiredProtocol',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidRequiredProtocol`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRequiredProtocolDescriptor =
    $convert.base64Decode(
        'ChdJbnZhbGlkUmVxdWlyZWRQcm90b2NvbBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidResponseCodeDescriptor instead')
const InvalidResponseCode$json = {
  '1': 'InvalidResponseCode',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidResponseCode`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidResponseCodeDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkUmVzcG9uc2VDb2RlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidTTLOrderDescriptor instead')
const InvalidTTLOrder$json = {
  '1': 'InvalidTTLOrder',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTTLOrder`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTTLOrderDescriptor = $convert.base64Decode(
    'Cg9JbnZhbGlkVFRMT3JkZXISGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidTaggingDescriptor instead')
const InvalidTagging$json = {
  '1': 'InvalidTagging',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTagging`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTaggingDescriptor = $convert.base64Decode(
    'Cg5JbnZhbGlkVGFnZ2luZxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidViewerCertificateDescriptor instead')
const InvalidViewerCertificate$json = {
  '1': 'InvalidViewerCertificate',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidViewerCertificate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidViewerCertificateDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkVmlld2VyQ2VydGlmaWNhdGUSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidWebACLIdDescriptor instead')
const InvalidWebACLId$json = {
  '1': 'InvalidWebACLId',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidWebACLId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidWebACLIdDescriptor = $convert.base64Decode(
    'Cg9JbnZhbGlkV2ViQUNMSWQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidationDescriptor instead')
const Invalidation$json = {
  '1': 'Invalidation',
  '2': [
    {'1': 'createtime', '3': 490895933, '4': 1, '5': 9, '10': 'createtime'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'invalidationbatch',
      '3': 458631988,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.InvalidationBatch',
      '10': 'invalidationbatch'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `Invalidation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidationDescriptor = $convert.base64Decode(
    'CgxJbnZhbGlkYXRpb24SIgoKY3JlYXRldGltZRi99InqASABKAlSCmNyZWF0ZXRpbWUSEgoCaW'
    'QYgfKitwEgASgJUgJpZBJPChFpbnZhbGlkYXRpb25iYXRjaBi01tjaASABKAsyHS5jbG91ZGZy'
    'b250LkludmFsaWRhdGlvbkJhdGNoUhFpbnZhbGlkYXRpb25iYXRjaBIZCgZzdGF0dXMYkOT7Ai'
    'ABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use invalidationBatchDescriptor instead')
const InvalidationBatch$json = {
  '1': 'InvalidationBatch',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {
      '1': 'paths',
      '3': 402204224,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Paths',
      '10': 'paths'
    },
  ],
};

/// Descriptor for `InvalidationBatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidationBatchDescriptor = $convert.base64Decode(
    'ChFJbnZhbGlkYXRpb25CYXRjaBIrCg9jYWxsZXJyZWZlcmVuY2UYmJmNSCABKAlSD2NhbGxlcn'
    'JlZmVyZW5jZRIrCgVwYXRocxjAzOS/ASABKAsyES5jbG91ZGZyb250LlBhdGhzUgVwYXRocw==');

@$core.Deprecated('Use invalidationListDescriptor instead')
const InvalidationList$json = {
  '1': 'InvalidationList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.InvalidationSummary',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `InvalidationList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidationListDescriptor = $convert.base64Decode(
    'ChBJbnZhbGlkYXRpb25MaXN0EiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZB'
    'I4CgVpdGVtcxiw8NgBIAMoCzIfLmNsb3VkZnJvbnQuSW52YWxpZGF0aW9uU3VtbWFyeVIFaXRl'
    'bXMSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYX'
    'hpdGVtcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj5'
    '5dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use invalidationSummaryDescriptor instead')
const InvalidationSummary$json = {
  '1': 'InvalidationSummary',
  '2': [
    {'1': 'createtime', '3': 490895933, '4': 1, '5': 9, '10': 'createtime'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `InvalidationSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidationSummaryDescriptor = $convert.base64Decode(
    'ChNJbnZhbGlkYXRpb25TdW1tYXJ5EiIKCmNyZWF0ZXRpbWUYvfSJ6gEgASgJUgpjcmVhdGV0aW'
    '1lEhIKAmlkGIHyorcBIAEoCVICaWQSGQoGc3RhdHVzGJDk+wIgASgJUgZzdGF0dXM=');

@$core.Deprecated('Use ipamCidrConfigDescriptor instead')
const IpamCidrConfig$json = {
  '1': 'IpamCidrConfig',
  '2': [
    {'1': 'anycastip', '3': 343152584, '4': 1, '5': 9, '10': 'anycastip'},
    {'1': 'cidr', '3': 390242040, '4': 1, '5': 9, '10': 'cidr'},
    {'1': 'ipampoolarn', '3': 116678828, '4': 1, '5': 9, '10': 'ipampoolarn'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.IpamCidrStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `IpamCidrConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ipamCidrConfigDescriptor = $convert.base64Decode(
    'Cg5JcGFtQ2lkckNvbmZpZxIgCglhbnljYXN0aXAYyK/QowEgASgJUglhbnljYXN0aXASFgoEY2'
    'lkchj4vYq6ASABKAlSBGNpZHISIwoLaXBhbXBvb2xhcm4YrMHRNyABKAlSC2lwYW1wb29sYXJu'
    'EjUKBnN0YXR1cxiQ5PsCIAEoDjIaLmNsb3VkZnJvbnQuSXBhbUNpZHJTdGF0dXNSBnN0YXR1cw'
    '==');

@$core.Deprecated('Use ipamConfigDescriptor instead')
const IpamConfig$json = {
  '1': 'IpamConfig',
  '2': [
    {
      '1': 'ipamcidrconfigs',
      '3': 314539184,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.IpamCidrConfig',
      '10': 'ipamcidrconfigs'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `IpamConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ipamConfigDescriptor = $convert.base64Decode(
    'CgpJcGFtQ29uZmlnEkgKD2lwYW1jaWRyY29uZmlncxiw+f2VASADKAsyGi5jbG91ZGZyb250Lk'
    'lwYW1DaWRyQ29uZmlnUg9pcGFtY2lkcmNvbmZpZ3MSHQoIcXVhbnRpdHkY+eXcXyABKAVSCHF1'
    'YW50aXR5');

@$core.Deprecated('Use kGKeyPairIdsDescriptor instead')
const KGKeyPairIds$json = {
  '1': 'KGKeyPairIds',
  '2': [
    {'1': 'keygroupid', '3': 497763283, '4': 1, '5': 9, '10': 'keygroupid'},
    {
      '1': 'keypairids',
      '3': 5804971,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyPairIds',
      '10': 'keypairids'
    },
  ],
};

/// Descriptor for `KGKeyPairIds`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kGKeyPairIdsDescriptor = $convert.base64Decode(
    'CgxLR0tleVBhaXJJZHMSIgoKa2V5Z3JvdXBpZBjTh63tASABKAlSCmtleWdyb3VwaWQSOQoKa2'
    'V5cGFpcmlkcxirp+ICIAEoCzIWLmNsb3VkZnJvbnQuS2V5UGFpcklkc1IKa2V5cGFpcmlkcw==');

@$core.Deprecated('Use keyGroupDescriptor instead')
const KeyGroup$json = {
  '1': 'KeyGroup',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'keygroupconfig',
      '3': 143012494,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroupConfig',
      '10': 'keygroupconfig'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `KeyGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyGroupDescriptor = $convert.base64Decode(
    'CghLZXlHcm91cBISCgJpZBiB8qK3ASABKAlSAmlkEkUKDmtleWdyb3VwY29uZmlnGI7lmEQgAS'
    'gLMhouY2xvdWRmcm9udC5LZXlHcm91cENvbmZpZ1IOa2V5Z3JvdXBjb25maWcSLQoQbGFzdG1v'
    'ZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZQ==');

@$core.Deprecated('Use keyGroupAlreadyExistsDescriptor instead')
const KeyGroupAlreadyExists$json = {
  '1': 'KeyGroupAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KeyGroupAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyGroupAlreadyExistsDescriptor = $convert.base64Decode(
    'ChVLZXlHcm91cEFscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use keyGroupConfigDescriptor instead')
const KeyGroupConfig$json = {
  '1': 'KeyGroupConfig',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `KeyGroupConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyGroupConfigDescriptor = $convert.base64Decode(
    'Cg5LZXlHcm91cENvbmZpZxIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbWVudBIXCgVpdGVtcx'
    'iw8NgBIAMoCVIFaXRlbXMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use keyGroupListDescriptor instead')
const KeyGroupList$json = {
  '1': 'KeyGroupList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.KeyGroupSummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `KeyGroupList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyGroupListDescriptor = $convert.base64Decode(
    'CgxLZXlHcm91cExpc3QSNAoFaXRlbXMYsPDYASADKAsyGy5jbG91ZGZyb250LktleUdyb3VwU3'
    'VtbWFyeVIFaXRlbXMSHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpuZXh0bWFy'
    'a2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdH'
    'k=');

@$core.Deprecated('Use keyGroupSummaryDescriptor instead')
const KeyGroupSummary$json = {
  '1': 'KeyGroupSummary',
  '2': [
    {
      '1': 'keygroup',
      '3': 518748096,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroup',
      '10': 'keygroup'
    },
  ],
};

/// Descriptor for `KeyGroupSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyGroupSummaryDescriptor = $convert.base64Decode(
    'Cg9LZXlHcm91cFN1bW1hcnkSNAoIa2V5Z3JvdXAYwO+t9wEgASgLMhQuY2xvdWRmcm9udC5LZX'
    'lHcm91cFIIa2V5Z3JvdXA=');

@$core.Deprecated('Use keyPairIdsDescriptor instead')
const KeyPairIds$json = {
  '1': 'KeyPairIds',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `KeyPairIds`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyPairIdsDescriptor = $convert.base64Decode(
    'CgpLZXlQYWlySWRzEhcKBWl0ZW1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIA'
    'EoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use keyValueStoreDescriptor instead')
const KeyValueStore$json = {
  '1': 'KeyValueStore',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `KeyValueStore`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyValueStoreDescriptor = $convert.base64Decode(
    'Cg1LZXlWYWx1ZVN0b3JlEhQKA2FybhidnK+9ASABKAlSA2FybhIcCgdjb21tZW50GP+/vsIBIA'
    'EoCVIHY29tbWVudBISCgJpZBiB8qK3ASABKAlSAmlkEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8'
    'cCABKAlSEGxhc3Rtb2RpZmllZHRpbWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIZCgZzdGF0dX'
    'MYkOT7AiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use keyValueStoreAssociationDescriptor instead')
const KeyValueStoreAssociation$json = {
  '1': 'KeyValueStoreAssociation',
  '2': [
    {
      '1': 'keyvaluestorearn',
      '3': 296227242,
      '4': 1,
      '5': 9,
      '10': 'keyvaluestorearn'
    },
  ],
};

/// Descriptor for `KeyValueStoreAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyValueStoreAssociationDescriptor =
    $convert.base64Decode(
        'ChhLZXlWYWx1ZVN0b3JlQXNzb2NpYXRpb24SLgoQa2V5dmFsdWVzdG9yZWFybhiqo6CNASABKA'
        'lSEGtleXZhbHVlc3RvcmVhcm4=');

@$core.Deprecated('Use keyValueStoreAssociationsDescriptor instead')
const KeyValueStoreAssociations$json = {
  '1': 'KeyValueStoreAssociations',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.KeyValueStoreAssociation',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `KeyValueStoreAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyValueStoreAssociationsDescriptor = $convert.base64Decode(
    'ChlLZXlWYWx1ZVN0b3JlQXNzb2NpYXRpb25zEj0KBWl0ZW1zGLDw2AEgAygLMiQuY2xvdWRmcm'
    '9udC5LZXlWYWx1ZVN0b3JlQXNzb2NpYXRpb25SBWl0ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgF'
    'UghxdWFudGl0eQ==');

@$core.Deprecated('Use keyValueStoreListDescriptor instead')
const KeyValueStoreList$json = {
  '1': 'KeyValueStoreList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.KeyValueStore',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `KeyValueStoreList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyValueStoreListDescriptor = $convert.base64Decode(
    'ChFLZXlWYWx1ZVN0b3JlTGlzdBIyCgVpdGVtcxiw8NgBIAMoCzIZLmNsb3VkZnJvbnQuS2V5Vm'
    'FsdWVTdG9yZVIFaXRlbXMSHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpuZXh0'
    'bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbn'
    'RpdHk=');

@$core.Deprecated('Use kinesisStreamConfigDescriptor instead')
const KinesisStreamConfig$json = {
  '1': 'KinesisStreamConfig',
  '2': [
    {'1': 'rolearn', '3': 327777153, '4': 1, '5': 9, '10': 'rolearn'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
  ],
};

/// Descriptor for `KinesisStreamConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kinesisStreamConfigDescriptor = $convert.base64Decode(
    'ChNLaW5lc2lzU3RyZWFtQ29uZmlnEhwKB3JvbGVhcm4YgfelnAEgASgJUgdyb2xlYXJuEiAKCX'
    'N0cmVhbWFybhjd86ryASABKAlSCXN0cmVhbWFybg==');

@$core.Deprecated('Use lambdaFunctionAssociationDescriptor instead')
const LambdaFunctionAssociation$json = {
  '1': 'LambdaFunctionAssociation',
  '2': [
    {
      '1': 'eventtype',
      '3': 468897896,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.EventType',
      '10': 'eventtype'
    },
    {'1': 'includebody', '3': 83342912, '4': 1, '5': 8, '10': 'includebody'},
    {
      '1': 'lambdafunctionarn',
      '3': 527901498,
      '4': 1,
      '5': 9,
      '10': 'lambdafunctionarn'
    },
  ],
};

/// Descriptor for `LambdaFunctionAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionAssociationDescriptor = $convert.base64Decode(
    'ChlMYW1iZGFGdW5jdGlvbkFzc29jaWF0aW9uEjcKCWV2ZW50dHlwZRjooMvfASABKA4yFS5jbG'
    '91ZGZyb250LkV2ZW50VHlwZVIJZXZlbnR0eXBlEiMKC2luY2x1ZGVib2R5GMDs3icgASgIUgtp'
    'bmNsdWRlYm9keRIwChFsYW1iZGFmdW5jdGlvbmFybhi6xtz7ASABKAlSEWxhbWJkYWZ1bmN0aW'
    '9uYXJu');

@$core.Deprecated('Use lambdaFunctionAssociationsDescriptor instead')
const LambdaFunctionAssociations$json = {
  '1': 'LambdaFunctionAssociations',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.LambdaFunctionAssociation',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `LambdaFunctionAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionAssociationsDescriptor =
    $convert.base64Decode(
        'ChpMYW1iZGFGdW5jdGlvbkFzc29jaWF0aW9ucxI+CgVpdGVtcxiw8NgBIAMoCzIlLmNsb3VkZn'
        'JvbnQuTGFtYmRhRnVuY3Rpb25Bc3NvY2lhdGlvblIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXyAB'
        'KAVSCHF1YW50aXR5');

@$core.Deprecated('Use listAnycastIpListsRequestDescriptor instead')
const ListAnycastIpListsRequest$json = {
  '1': 'ListAnycastIpListsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListAnycastIpListsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAnycastIpListsRequestDescriptor =
    $convert.base64Decode(
        'ChlMaXN0QW55Y2FzdElwTGlzdHNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh'
        '4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listAnycastIpListsResultDescriptor instead')
const ListAnycastIpListsResult$json = {
  '1': 'ListAnycastIpListsResult',
  '2': [
    {
      '1': 'anycastiplists',
      '3': 489147285,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.AnycastIpListCollection',
      '8': {},
      '10': 'anycastiplists'
    },
  ],
};

/// Descriptor for `ListAnycastIpListsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAnycastIpListsResultDescriptor = $convert.base64Decode(
    'ChhMaXN0QW55Y2FzdElwTGlzdHNSZXN1bHQSVQoOYW55Y2FzdGlwbGlzdHMYlZef6QEgASgLMi'
    'MuY2xvdWRmcm9udC5BbnljYXN0SXBMaXN0Q29sbGVjdGlvbkIEiLUYAVIOYW55Y2FzdGlwbGlz'
    'dHM=');

@$core.Deprecated('Use listCachePoliciesRequestDescriptor instead')
const ListCachePoliciesRequest$json = {
  '1': 'ListCachePoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CachePolicyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ListCachePoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCachePoliciesRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0Q2FjaGVQb2xpY2llc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHg'
    'oIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIzCgR0eXBlGO6g14oBIAEoDjIbLmNsb3Vk'
    'ZnJvbnQuQ2FjaGVQb2xpY3lUeXBlUgR0eXBl');

@$core.Deprecated('Use listCachePoliciesResultDescriptor instead')
const ListCachePoliciesResult$json = {
  '1': 'ListCachePoliciesResult',
  '2': [
    {
      '1': 'cachepolicylist',
      '3': 116073544,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyList',
      '8': {},
      '10': 'cachepolicylist'
    },
  ],
};

/// Descriptor for `ListCachePoliciesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCachePoliciesResultDescriptor =
    $convert.base64Decode(
        'ChdMaXN0Q2FjaGVQb2xpY2llc1Jlc3VsdBJOCg9jYWNoZXBvbGljeWxpc3QYyMisNyABKAsyGy'
        '5jbG91ZGZyb250LkNhY2hlUG9saWN5TGlzdEIEiLUYAVIPY2FjaGVwb2xpY3lsaXN0');

@$core.Deprecated(
    'Use listCloudFrontOriginAccessIdentitiesRequestDescriptor instead')
const ListCloudFrontOriginAccessIdentitiesRequest$json = {
  '1': 'ListCloudFrontOriginAccessIdentitiesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListCloudFrontOriginAccessIdentitiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listCloudFrontOriginAccessIdentitiesRequestDescriptor =
    $convert.base64Decode(
        'CitMaXN0Q2xvdWRGcm9udE9yaWdpbkFjY2Vzc0lkZW50aXRpZXNSZXF1ZXN0EhkKBm1hcmtlch'
        'i43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated(
    'Use listCloudFrontOriginAccessIdentitiesResultDescriptor instead')
const ListCloudFrontOriginAccessIdentitiesResult$json = {
  '1': 'ListCloudFrontOriginAccessIdentitiesResult',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentitylist',
      '3': 247741992,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentityList',
      '8': {},
      '10': 'cloudfrontoriginaccessidentitylist'
    },
  ],
};

/// Descriptor for `ListCloudFrontOriginAccessIdentitiesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listCloudFrontOriginAccessIdentitiesResultDescriptor =
    $convert.base64Decode(
        'CipMaXN0Q2xvdWRGcm9udE9yaWdpbkFjY2Vzc0lkZW50aXRpZXNSZXN1bHQShwEKImNsb3VkZn'
        'JvbnRvcmlnaW5hY2Nlc3NpZGVudGl0eWxpc3QYqPyQdiABKAsyLi5jbG91ZGZyb250LkNsb3Vk'
        'RnJvbnRPcmlnaW5BY2Nlc3NJZGVudGl0eUxpc3RCBIi1GAFSImNsb3VkZnJvbnRvcmlnaW5hY2'
        'Nlc3NpZGVudGl0eWxpc3Q=');

@$core.Deprecated('Use listConflictingAliasesRequestDescriptor instead')
const ListConflictingAliasesRequest$json = {
  '1': 'ListConflictingAliasesRequest',
  '2': [
    {'1': 'alias', '3': 48362232, '4': 1, '5': 9, '10': 'alias'},
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListConflictingAliasesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConflictingAliasesRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q29uZmxpY3RpbmdBbGlhc2VzUmVxdWVzdBIXCgVhbGlhcxj45YcXIAEoCVIFYWxpYX'
        'MSKQoOZGlzdHJpYnV0aW9uaWQY57H7QyABKAlSDmRpc3RyaWJ1dGlvbmlkEhkKBm1hcmtlchi4'
        '3c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listConflictingAliasesResultDescriptor instead')
const ListConflictingAliasesResult$json = {
  '1': 'ListConflictingAliasesResult',
  '2': [
    {
      '1': 'conflictingaliaseslist',
      '3': 132174420,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConflictingAliasesList',
      '8': {},
      '10': 'conflictingaliaseslist'
    },
  ],
};

/// Descriptor for `ListConflictingAliasesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConflictingAliasesResultDescriptor =
    $convert.base64Decode(
        'ChxMaXN0Q29uZmxpY3RpbmdBbGlhc2VzUmVzdWx0EmMKFmNvbmZsaWN0aW5nYWxpYXNlc2xpc3'
        'QY1KSDPyABKAsyIi5jbG91ZGZyb250LkNvbmZsaWN0aW5nQWxpYXNlc0xpc3RCBIi1GAFSFmNv'
        'bmZsaWN0aW5nYWxpYXNlc2xpc3Q=');

@$core.Deprecated('Use listConnectionFunctionsRequestDescriptor instead')
const ListConnectionFunctionsRequest$json = {
  '1': 'ListConnectionFunctionsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `ListConnectionFunctionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConnectionFunctionsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0Q29ubmVjdGlvbkZ1bmN0aW9uc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYX'
        'JrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIyCgVzdGFnZRjOl9hwIAEoDjIZ'
        'LmNsb3VkZnJvbnQuRnVuY3Rpb25TdGFnZVIFc3RhZ2U=');

@$core.Deprecated('Use listConnectionFunctionsResultDescriptor instead')
const ListConnectionFunctionsResult$json = {
  '1': 'ListConnectionFunctionsResult',
  '2': [
    {
      '1': 'connectionfunctions',
      '3': 424487851,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionSummary',
      '10': 'connectionfunctions'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListConnectionFunctionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConnectionFunctionsResultDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q29ubmVjdGlvbkZ1bmN0aW9uc1Jlc3VsdBJbChNjb25uZWN0aW9uZnVuY3Rpb25zGK'
        'vXtMoBIAMoCzIlLmNsb3VkZnJvbnQuQ29ubmVjdGlvbkZ1bmN0aW9uU3VtbWFyeVITY29ubmVj'
        'dGlvbmZ1bmN0aW9ucxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use listConnectionGroupsRequestDescriptor instead')
const ListConnectionGroupsRequest$json = {
  '1': 'ListConnectionGroupsRequest',
  '2': [
    {
      '1': 'associationfilter',
      '3': 309904703,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionGroupAssociationFilter',
      '10': 'associationfilter'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListConnectionGroupsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConnectionGroupsRequestDescriptor = $convert.base64Decode(
    'ChtMaXN0Q29ubmVjdGlvbkdyb3Vwc1JlcXVlc3QSXgoRYXNzb2NpYXRpb25maWx0ZXIYv4rjkw'
    'EgASgLMiwuY2xvdWRmcm9udC5Db25uZWN0aW9uR3JvdXBBc3NvY2lhdGlvbkZpbHRlclIRYXNz'
    'b2NpYXRpb25maWx0ZXISGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlN'
    'ba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listConnectionGroupsResultDescriptor instead')
const ListConnectionGroupsResult$json = {
  '1': 'ListConnectionGroupsResult',
  '2': [
    {
      '1': 'connectiongroups',
      '3': 182892570,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ConnectionGroupSummary',
      '10': 'connectiongroups'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListConnectionGroupsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConnectionGroupsResultDescriptor =
    $convert.base64Decode(
        'ChpMaXN0Q29ubmVjdGlvbkdyb3Vwc1Jlc3VsdBJRChBjb25uZWN0aW9uZ3JvdXBzGJrwmlcgAy'
        'gLMiIuY2xvdWRmcm9udC5Db25uZWN0aW9uR3JvdXBTdW1tYXJ5UhBjb25uZWN0aW9uZ3JvdXBz'
        'EiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2Vy');

@$core
    .Deprecated('Use listContinuousDeploymentPoliciesRequestDescriptor instead')
const ListContinuousDeploymentPoliciesRequest$json = {
  '1': 'ListContinuousDeploymentPoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListContinuousDeploymentPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContinuousDeploymentPoliciesRequestDescriptor =
    $convert.base64Decode(
        'CidMaXN0Q29udGludW91c0RlcGxveW1lbnRQb2xpY2llc1JlcXVlc3QSGQoGbWFya2VyGLjdzS'
        'ogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core
    .Deprecated('Use listContinuousDeploymentPoliciesResultDescriptor instead')
const ListContinuousDeploymentPoliciesResult$json = {
  '1': 'ListContinuousDeploymentPoliciesResult',
  '2': [
    {
      '1': 'continuousdeploymentpolicylist',
      '3': 499667708,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicyList',
      '8': {},
      '10': 'continuousdeploymentpolicylist'
    },
  ],
};

/// Descriptor for `ListContinuousDeploymentPoliciesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContinuousDeploymentPoliciesResultDescriptor =
    $convert.base64Decode(
        'CiZMaXN0Q29udGludW91c0RlcGxveW1lbnRQb2xpY2llc1Jlc3VsdBJ8Ch5jb250aW51b3VzZG'
        'VwbG95bWVudHBvbGljeWxpc3QY/KWh7gEgASgLMiouY2xvdWRmcm9udC5Db250aW51b3VzRGVw'
        'bG95bWVudFBvbGljeUxpc3RCBIi1GAFSHmNvbnRpbnVvdXNkZXBsb3ltZW50cG9saWN5bGlzdA'
        '==');

@$core.Deprecated(
    'Use listDistributionTenantsByCustomizationRequestDescriptor instead')
const ListDistributionTenantsByCustomizationRequest$json = {
  '1': 'ListDistributionTenantsByCustomizationRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'webaclarn', '3': 82506659, '4': 1, '5': 9, '10': 'webaclarn'},
  ],
};

/// Descriptor for `ListDistributionTenantsByCustomizationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionTenantsByCustomizationRequestDescriptor =
    $convert.base64Decode(
        'Ci1MaXN0RGlzdHJpYnV0aW9uVGVuYW50c0J5Q3VzdG9taXphdGlvblJlcXVlc3QSKQoOY2VydG'
        'lmaWNhdGVhcm4Y+MqZLCABKAlSDmNlcnRpZmljYXRlYXJuEhkKBm1hcmtlchi43c0qIAEoCVIG'
        'bWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSHwoJd2ViYWNsYXJuGKPnqy'
        'cgASgJUgl3ZWJhY2xhcm4=');

@$core.Deprecated(
    'Use listDistributionTenantsByCustomizationResultDescriptor instead')
const ListDistributionTenantsByCustomizationResult$json = {
  '1': 'ListDistributionTenantsByCustomizationResult',
  '2': [
    {
      '1': 'distributiontenantlist',
      '3': 354250620,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DistributionTenantSummary',
      '10': 'distributiontenantlist'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListDistributionTenantsByCustomizationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionTenantsByCustomizationResultDescriptor =
    $convert.base64Decode(
        'CixMaXN0RGlzdHJpYnV0aW9uVGVuYW50c0J5Q3VzdG9taXphdGlvblJlc3VsdBJhChZkaXN0cm'
        'lidXRpb250ZW5hbnRsaXN0GPze9agBIAMoCzIlLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uVGVu'
        'YW50U3VtbWFyeVIWZGlzdHJpYnV0aW9udGVuYW50bGlzdBIiCgpuZXh0bWFya2VyGKOBrv0BIA'
        'EoCVIKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use listDistributionTenantsRequestDescriptor instead')
const ListDistributionTenantsRequest$json = {
  '1': 'ListDistributionTenantsRequest',
  '2': [
    {
      '1': 'associationfilter',
      '3': 309904703,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionTenantAssociationFilter',
      '10': 'associationfilter'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDistributionTenantsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionTenantsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0RGlzdHJpYnV0aW9uVGVuYW50c1JlcXVlc3QSYQoRYXNzb2NpYXRpb25maWx0ZXIYv4'
        'rjkwEgASgLMi8uY2xvdWRmcm9udC5EaXN0cmlidXRpb25UZW5hbnRBc3NvY2lhdGlvbkZpbHRl'
        'clIRYXNzb2NpYXRpb25maWx0ZXISGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aX'
        'RlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listDistributionTenantsResultDescriptor instead')
const ListDistributionTenantsResult$json = {
  '1': 'ListDistributionTenantsResult',
  '2': [
    {
      '1': 'distributiontenantlist',
      '3': 354250620,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DistributionTenantSummary',
      '10': 'distributiontenantlist'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListDistributionTenantsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionTenantsResultDescriptor = $convert.base64Decode(
    'Ch1MaXN0RGlzdHJpYnV0aW9uVGVuYW50c1Jlc3VsdBJhChZkaXN0cmlidXRpb250ZW5hbnRsaX'
    'N0GPze9agBIAMoCzIlLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uVGVuYW50U3VtbWFyeVIWZGlz'
    'dHJpYnV0aW9udGVuYW50bGlzdBIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlcg'
    '==');

@$core.Deprecated(
    'Use listDistributionsByAnycastIpListIdRequestDescriptor instead')
const ListDistributionsByAnycastIpListIdRequest$json = {
  '1': 'ListDistributionsByAnycastIpListIdRequest',
  '2': [
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDistributionsByAnycastIpListIdRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByAnycastIpListIdRequestDescriptor = $convert.base64Decode(
        'CilMaXN0RGlzdHJpYnV0aW9uc0J5QW55Y2FzdElwTGlzdElkUmVxdWVzdBIsCg9hbnljYXN0aX'
        'BsaXN0aWQYg6z4zQEgASgJUg9hbnljYXN0aXBsaXN0aWQSGQoGbWFya2VyGLjdzSogASgJUgZt'
        'YXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated(
    'Use listDistributionsByAnycastIpListIdResultDescriptor instead')
const ListDistributionsByAnycastIpListIdResult$json = {
  '1': 'ListDistributionsByAnycastIpListIdResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByAnycastIpListIdResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByAnycastIpListIdResultDescriptor =
    $convert.base64Decode(
        'CihMaXN0RGlzdHJpYnV0aW9uc0J5QW55Y2FzdElwTGlzdElkUmVzdWx0ElEKEGRpc3RyaWJ1dG'
        'lvbmxpc3QY9P+DFiABKAsyHC5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkxpc3RCBIi1GAFSEGRp'
        'c3RyaWJ1dGlvbmxpc3Q=');

@$core
    .Deprecated('Use listDistributionsByCachePolicyIdRequestDescriptor instead')
const ListDistributionsByCachePolicyIdRequest$json = {
  '1': 'ListDistributionsByCachePolicyIdRequest',
  '2': [
    {
      '1': 'cachepolicyid',
      '3': 431434163,
      '4': 1,
      '5': 9,
      '10': 'cachepolicyid'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDistributionsByCachePolicyIdRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByCachePolicyIdRequestDescriptor =
    $convert.base64Decode(
        'CidMaXN0RGlzdHJpYnV0aW9uc0J5Q2FjaGVQb2xpY3lJZFJlcXVlc3QSKAoNY2FjaGVwb2xpY3'
        'lpZBiz09zNASABKAlSDWNhY2hlcG9saWN5aWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXIS'
        'HgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core
    .Deprecated('Use listDistributionsByCachePolicyIdResultDescriptor instead')
const ListDistributionsByCachePolicyIdResult$json = {
  '1': 'ListDistributionsByCachePolicyIdResult',
  '2': [
    {
      '1': 'distributionidlist',
      '3': 481677267,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionIdList',
      '8': {},
      '10': 'distributionidlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByCachePolicyIdResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByCachePolicyIdResultDescriptor =
    $convert.base64Decode(
        'CiZMaXN0RGlzdHJpYnV0aW9uc0J5Q2FjaGVQb2xpY3lJZFJlc3VsdBJYChJkaXN0cmlidXRpb2'
        '5pZGxpc3QY05/X5QEgASgLMh4uY2xvdWRmcm9udC5EaXN0cmlidXRpb25JZExpc3RCBIi1GAFS'
        'EmRpc3RyaWJ1dGlvbmlkbGlzdA==');

@$core.Deprecated(
    'Use listDistributionsByConnectionFunctionRequestDescriptor instead')
const ListDistributionsByConnectionFunctionRequest$json = {
  '1': 'ListDistributionsByConnectionFunctionRequest',
  '2': [
    {
      '1': 'connectionfunctionidentifier',
      '3': 241101239,
      '4': 1,
      '5': 9,
      '10': 'connectionfunctionidentifier'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDistributionsByConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByConnectionFunctionRequestDescriptor =
    $convert.base64Decode(
        'CixMaXN0RGlzdHJpYnV0aW9uc0J5Q29ubmVjdGlvbkZ1bmN0aW9uUmVxdWVzdBJFChxjb25uZW'
        'N0aW9uZnVuY3Rpb25pZGVudGlmaWVyGLfT+3IgASgJUhxjb25uZWN0aW9uZnVuY3Rpb25pZGVu'
        'dGlmaWVyEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBV'
        'IIbWF4aXRlbXM=');

@$core.Deprecated(
    'Use listDistributionsByConnectionFunctionResultDescriptor instead')
const ListDistributionsByConnectionFunctionResult$json = {
  '1': 'ListDistributionsByConnectionFunctionResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByConnectionFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByConnectionFunctionResultDescriptor =
    $convert.base64Decode(
        'CitMaXN0RGlzdHJpYnV0aW9uc0J5Q29ubmVjdGlvbkZ1bmN0aW9uUmVzdWx0ElEKEGRpc3RyaW'
        'J1dGlvbmxpc3QY9P+DFiABKAsyHC5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkxpc3RCBIi1GAFS'
        'EGRpc3RyaWJ1dGlvbmxpc3Q=');

@$core.Deprecated(
    'Use listDistributionsByConnectionModeRequestDescriptor instead')
const ListDistributionsByConnectionModeRequest$json = {
  '1': 'ListDistributionsByConnectionModeRequest',
  '2': [
    {
      '1': 'connectionmode',
      '3': 82068023,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ConnectionMode',
      '10': 'connectionmode'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDistributionsByConnectionModeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByConnectionModeRequestDescriptor =
    $convert.base64Decode(
        'CihMaXN0RGlzdHJpYnV0aW9uc0J5Q29ubmVjdGlvbk1vZGVSZXF1ZXN0EkUKDmNvbm5lY3Rpb2'
        '5tb2RlGLeEkScgASgOMhouY2xvdWRmcm9udC5Db25uZWN0aW9uTW9kZVIOY29ubmVjdGlvbm1v'
        'ZGUSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYX'
        'hpdGVtcw==');

@$core
    .Deprecated('Use listDistributionsByConnectionModeResultDescriptor instead')
const ListDistributionsByConnectionModeResult$json = {
  '1': 'ListDistributionsByConnectionModeResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByConnectionModeResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByConnectionModeResultDescriptor =
    $convert.base64Decode(
        'CidMaXN0RGlzdHJpYnV0aW9uc0J5Q29ubmVjdGlvbk1vZGVSZXN1bHQSUQoQZGlzdHJpYnV0aW'
        '9ubGlzdBj0/4MWIAEoCzIcLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uTGlzdEIEiLUYAVIQZGlz'
        'dHJpYnV0aW9ubGlzdA==');

@$core.Deprecated('Use listDistributionsByKeyGroupRequestDescriptor instead')
const ListDistributionsByKeyGroupRequest$json = {
  '1': 'ListDistributionsByKeyGroupRequest',
  '2': [
    {'1': 'keygroupid', '3': 497763283, '4': 1, '5': 9, '10': 'keygroupid'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDistributionsByKeyGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByKeyGroupRequestDescriptor =
    $convert.base64Decode(
        'CiJMaXN0RGlzdHJpYnV0aW9uc0J5S2V5R3JvdXBSZXF1ZXN0EiIKCmtleWdyb3VwaWQY04et7Q'
        'EgASgJUgprZXlncm91cGlkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1z'
        'GJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listDistributionsByKeyGroupResultDescriptor instead')
const ListDistributionsByKeyGroupResult$json = {
  '1': 'ListDistributionsByKeyGroupResult',
  '2': [
    {
      '1': 'distributionidlist',
      '3': 481677267,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionIdList',
      '8': {},
      '10': 'distributionidlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByKeyGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByKeyGroupResultDescriptor =
    $convert.base64Decode(
        'CiFMaXN0RGlzdHJpYnV0aW9uc0J5S2V5R3JvdXBSZXN1bHQSWAoSZGlzdHJpYnV0aW9uaWRsaX'
        'N0GNOf1+UBIAEoCzIeLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uSWRMaXN0QgSItRgBUhJkaXN0'
        'cmlidXRpb25pZGxpc3Q=');

@$core.Deprecated(
    'Use listDistributionsByOriginRequestPolicyIdRequestDescriptor instead')
const ListDistributionsByOriginRequestPolicyIdRequest$json = {
  '1': 'ListDistributionsByOriginRequestPolicyIdRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'originrequestpolicyid',
      '3': 298538616,
      '4': 1,
      '5': 9,
      '10': 'originrequestpolicyid'
    },
  ],
};

/// Descriptor for `ListDistributionsByOriginRequestPolicyIdRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByOriginRequestPolicyIdRequestDescriptor =
    $convert.base64Decode(
        'Ci9MaXN0RGlzdHJpYnV0aW9uc0J5T3JpZ2luUmVxdWVzdFBvbGljeUlkUmVxdWVzdBIZCgZtYX'
        'JrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEjgK'
        'FW9yaWdpbnJlcXVlc3Rwb2xpY3lpZBj4rK2OASABKAlSFW9yaWdpbnJlcXVlc3Rwb2xpY3lpZA'
        '==');

@$core.Deprecated(
    'Use listDistributionsByOriginRequestPolicyIdResultDescriptor instead')
const ListDistributionsByOriginRequestPolicyIdResult$json = {
  '1': 'ListDistributionsByOriginRequestPolicyIdResult',
  '2': [
    {
      '1': 'distributionidlist',
      '3': 481677267,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionIdList',
      '8': {},
      '10': 'distributionidlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByOriginRequestPolicyIdResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByOriginRequestPolicyIdResultDescriptor =
    $convert.base64Decode(
        'Ci5MaXN0RGlzdHJpYnV0aW9uc0J5T3JpZ2luUmVxdWVzdFBvbGljeUlkUmVzdWx0ElgKEmRpc3'
        'RyaWJ1dGlvbmlkbGlzdBjTn9flASABKAsyHi5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbklkTGlz'
        'dEIEiLUYAVISZGlzdHJpYnV0aW9uaWRsaXN0');

@$core
    .Deprecated('Use listDistributionsByOwnedResourceRequestDescriptor instead')
const ListDistributionsByOwnedResourceRequest$json = {
  '1': 'ListDistributionsByOwnedResourceRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListDistributionsByOwnedResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByOwnedResourceRequestDescriptor =
    $convert.base64Decode(
        'CidMaXN0RGlzdHJpYnV0aW9uc0J5T3duZWRSZXNvdXJjZVJlcXVlc3QSGQoGbWFya2VyGLjdzS'
        'ogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxIkCgtyZXNvdXJj'
        'ZWFybhit+NmtASABKAlSC3Jlc291cmNlYXJu');

@$core
    .Deprecated('Use listDistributionsByOwnedResourceResultDescriptor instead')
const ListDistributionsByOwnedResourceResult$json = {
  '1': 'ListDistributionsByOwnedResourceResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionIdOwnerList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByOwnedResourceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByOwnedResourceResultDescriptor =
    $convert.base64Decode(
        'CiZMaXN0RGlzdHJpYnV0aW9uc0J5T3duZWRSZXNvdXJjZVJlc3VsdBJYChBkaXN0cmlidXRpb2'
        '5saXN0GPT/gxYgASgLMiMuY2xvdWRmcm9udC5EaXN0cmlidXRpb25JZE93bmVyTGlzdEIEiLUY'
        'AVIQZGlzdHJpYnV0aW9ubGlzdA==');

@$core.Deprecated(
    'Use listDistributionsByRealtimeLogConfigRequestDescriptor instead')
const ListDistributionsByRealtimeLogConfigRequest$json = {
  '1': 'ListDistributionsByRealtimeLogConfigRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'realtimelogconfigarn',
      '3': 152963408,
      '4': 1,
      '5': 9,
      '10': 'realtimelogconfigarn'
    },
    {
      '1': 'realtimelogconfigname',
      '3': 503770540,
      '4': 1,
      '5': 9,
      '10': 'realtimelogconfigname'
    },
  ],
};

/// Descriptor for `ListDistributionsByRealtimeLogConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByRealtimeLogConfigRequestDescriptor =
    $convert.base64Decode(
        'CitMaXN0RGlzdHJpYnV0aW9uc0J5UmVhbHRpbWVMb2dDb25maWdSZXF1ZXN0EhkKBm1hcmtlch'
        'i43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSNQoUcmVh'
        'bHRpbWVsb2djb25maWdhcm4Y0JL4SCABKAlSFHJlYWx0aW1lbG9nY29uZmlnYXJuEjgKFXJlYW'
        'x0aW1lbG9nY29uZmlnbmFtZRis25vwASABKAlSFXJlYWx0aW1lbG9nY29uZmlnbmFtZQ==');

@$core.Deprecated(
    'Use listDistributionsByRealtimeLogConfigResultDescriptor instead')
const ListDistributionsByRealtimeLogConfigResult$json = {
  '1': 'ListDistributionsByRealtimeLogConfigResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByRealtimeLogConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByRealtimeLogConfigResultDescriptor =
    $convert.base64Decode(
        'CipMaXN0RGlzdHJpYnV0aW9uc0J5UmVhbHRpbWVMb2dDb25maWdSZXN1bHQSUQoQZGlzdHJpYn'
        'V0aW9ubGlzdBj0/4MWIAEoCzIcLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uTGlzdEIEiLUYAVIQ'
        'ZGlzdHJpYnV0aW9ubGlzdA==');

@$core.Deprecated(
    'Use listDistributionsByResponseHeadersPolicyIdRequestDescriptor instead')
const ListDistributionsByResponseHeadersPolicyIdRequest$json = {
  '1': 'ListDistributionsByResponseHeadersPolicyIdRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'responseheaderspolicyid',
      '3': 244029524,
      '4': 1,
      '5': 9,
      '10': 'responseheaderspolicyid'
    },
  ],
};

/// Descriptor for `ListDistributionsByResponseHeadersPolicyIdRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByResponseHeadersPolicyIdRequestDescriptor =
    $convert.base64Decode(
        'CjFMaXN0RGlzdHJpYnV0aW9uc0J5UmVzcG9uc2VIZWFkZXJzUG9saWN5SWRSZXF1ZXN0EhkKBm'
        '1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMS'
        'OwoXcmVzcG9uc2VoZWFkZXJzcG9saWN5aWQY1LCudCABKAlSF3Jlc3BvbnNlaGVhZGVyc3BvbG'
        'ljeWlk');

@$core.Deprecated(
    'Use listDistributionsByResponseHeadersPolicyIdResultDescriptor instead')
const ListDistributionsByResponseHeadersPolicyIdResult$json = {
  '1': 'ListDistributionsByResponseHeadersPolicyIdResult',
  '2': [
    {
      '1': 'distributionidlist',
      '3': 481677267,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionIdList',
      '8': {},
      '10': 'distributionidlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByResponseHeadersPolicyIdResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDistributionsByResponseHeadersPolicyIdResultDescriptor =
    $convert.base64Decode(
        'CjBMaXN0RGlzdHJpYnV0aW9uc0J5UmVzcG9uc2VIZWFkZXJzUG9saWN5SWRSZXN1bHQSWAoSZG'
        'lzdHJpYnV0aW9uaWRsaXN0GNOf1+UBIAEoCzIeLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uSWRM'
        'aXN0QgSItRgBUhJkaXN0cmlidXRpb25pZGxpc3Q=');

@$core.Deprecated('Use listDistributionsByTrustStoreRequestDescriptor instead')
const ListDistributionsByTrustStoreRequest$json = {
  '1': 'ListDistributionsByTrustStoreRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'truststoreidentifier',
      '3': 327436686,
      '4': 1,
      '5': 9,
      '10': 'truststoreidentifier'
    },
  ],
};

/// Descriptor for `ListDistributionsByTrustStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByTrustStoreRequestDescriptor =
    $convert.base64Decode(
        'CiRMaXN0RGlzdHJpYnV0aW9uc0J5VHJ1c3RTdG9yZVJlcXVlc3QSGQoGbWFya2VyGLjdzSogAS'
        'gJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcxI2ChR0cnVzdHN0b3Jl'
        'aWRlbnRpZmllchiOk5GcASABKAlSFHRydXN0c3RvcmVpZGVudGlmaWVy');

@$core.Deprecated('Use listDistributionsByTrustStoreResultDescriptor instead')
const ListDistributionsByTrustStoreResult$json = {
  '1': 'ListDistributionsByTrustStoreResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByTrustStoreResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByTrustStoreResultDescriptor =
    $convert.base64Decode(
        'CiNMaXN0RGlzdHJpYnV0aW9uc0J5VHJ1c3RTdG9yZVJlc3VsdBJRChBkaXN0cmlidXRpb25saX'
        'N0GPT/gxYgASgLMhwuY2xvdWRmcm9udC5EaXN0cmlidXRpb25MaXN0QgSItRgBUhBkaXN0cmli'
        'dXRpb25saXN0');

@$core.Deprecated('Use listDistributionsByVpcOriginIdRequestDescriptor instead')
const ListDistributionsByVpcOriginIdRequest$json = {
  '1': 'ListDistributionsByVpcOriginIdRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'vpcoriginid', '3': 365404648, '4': 1, '5': 9, '10': 'vpcoriginid'},
  ],
};

/// Descriptor for `ListDistributionsByVpcOriginIdRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByVpcOriginIdRequestDescriptor =
    $convert.base64Decode(
        'CiVMaXN0RGlzdHJpYnV0aW9uc0J5VnBjT3JpZ2luSWRSZXF1ZXN0EhkKBm1hcmtlchi43c0qIA'
        'EoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSJAoLdnBjb3JpZ2lu'
        'aWQY6MOergEgASgJUgt2cGNvcmlnaW5pZA==');

@$core.Deprecated('Use listDistributionsByVpcOriginIdResultDescriptor instead')
const ListDistributionsByVpcOriginIdResult$json = {
  '1': 'ListDistributionsByVpcOriginIdResult',
  '2': [
    {
      '1': 'distributionidlist',
      '3': 481677267,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionIdList',
      '8': {},
      '10': 'distributionidlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByVpcOriginIdResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByVpcOriginIdResultDescriptor =
    $convert.base64Decode(
        'CiRMaXN0RGlzdHJpYnV0aW9uc0J5VnBjT3JpZ2luSWRSZXN1bHQSWAoSZGlzdHJpYnV0aW9uaW'
        'RsaXN0GNOf1+UBIAEoCzIeLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uSWRMaXN0QgSItRgBUhJk'
        'aXN0cmlidXRpb25pZGxpc3Q=');

@$core.Deprecated('Use listDistributionsByWebACLIdRequestDescriptor instead')
const ListDistributionsByWebACLIdRequest$json = {
  '1': 'ListDistributionsByWebACLIdRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'webaclid', '3': 161274579, '4': 1, '5': 9, '10': 'webaclid'},
  ],
};

/// Descriptor for `ListDistributionsByWebACLIdRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByWebACLIdRequestDescriptor =
    $convert.base64Decode(
        'CiJMaXN0RGlzdHJpYnV0aW9uc0J5V2ViQUNMSWRSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCV'
        'IGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSHQoId2ViYWNsaWQY07Xz'
        'TCABKAlSCHdlYmFjbGlk');

@$core.Deprecated('Use listDistributionsByWebACLIdResultDescriptor instead')
const ListDistributionsByWebACLIdResult$json = {
  '1': 'ListDistributionsByWebACLIdResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsByWebACLIdResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsByWebACLIdResultDescriptor =
    $convert.base64Decode(
        'CiFMaXN0RGlzdHJpYnV0aW9uc0J5V2ViQUNMSWRSZXN1bHQSUQoQZGlzdHJpYnV0aW9ubGlzdB'
        'j0/4MWIAEoCzIcLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uTGlzdEIEiLUYAVIQZGlzdHJpYnV0'
        'aW9ubGlzdA==');

@$core.Deprecated('Use listDistributionsRequestDescriptor instead')
const ListDistributionsRequest$json = {
  '1': 'ListDistributionsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDistributionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsRequestDescriptor =
    $convert.base64Decode(
        'ChhMaXN0RGlzdHJpYnV0aW9uc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHg'
        'oIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listDistributionsResultDescriptor instead')
const ListDistributionsResult$json = {
  '1': 'ListDistributionsResult',
  '2': [
    {
      '1': 'distributionlist',
      '3': 46202868,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionList',
      '8': {},
      '10': 'distributionlist'
    },
  ],
};

/// Descriptor for `ListDistributionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDistributionsResultDescriptor = $convert.base64Decode(
    'ChdMaXN0RGlzdHJpYnV0aW9uc1Jlc3VsdBJRChBkaXN0cmlidXRpb25saXN0GPT/gxYgASgLMh'
    'wuY2xvdWRmcm9udC5EaXN0cmlidXRpb25MaXN0QgSItRgBUhBkaXN0cmlidXRpb25saXN0');

@$core.Deprecated('Use listDomainConflictsRequestDescriptor instead')
const ListDomainConflictsRequest$json = {
  '1': 'ListDomainConflictsRequest',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'domaincontrolvalidationresource',
      '3': 405541194,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionResourceId',
      '10': 'domaincontrolvalidationresource'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListDomainConflictsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDomainConflictsRequestDescriptor = $convert.base64Decode(
    'ChpMaXN0RG9tYWluQ29uZmxpY3RzUmVxdWVzdBIaCgZkb21haW4YkpLy8AEgASgJUgZkb21haW'
    '4ScAofZG9tYWluY29udHJvbHZhbGlkYXRpb25yZXNvdXJjZRjKorDBASABKAsyIi5jbG91ZGZy'
    'b250LkRpc3RyaWJ1dGlvblJlc291cmNlSWRSH2RvbWFpbmNvbnRyb2x2YWxpZGF0aW9ucmVzb3'
    'VyY2USGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUght'
    'YXhpdGVtcw==');

@$core.Deprecated('Use listDomainConflictsResultDescriptor instead')
const ListDomainConflictsResult$json = {
  '1': 'ListDomainConflictsResult',
  '2': [
    {
      '1': 'domainconflicts',
      '3': 324386609,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DomainConflict',
      '10': 'domainconflicts'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListDomainConflictsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDomainConflictsResultDescriptor = $convert.base64Decode(
    'ChlMaXN0RG9tYWluQ29uZmxpY3RzUmVzdWx0EkgKD2RvbWFpbmNvbmZsaWN0cxix/taaASADKA'
    'syGi5jbG91ZGZyb250LkRvbWFpbkNvbmZsaWN0Ug9kb21haW5jb25mbGljdHMSIgoKbmV4dG1h'
    'cmtlchijga79ASABKAlSCm5leHRtYXJrZXI=');

@$core
    .Deprecated('Use listFieldLevelEncryptionConfigsRequestDescriptor instead')
const ListFieldLevelEncryptionConfigsRequest$json = {
  '1': 'ListFieldLevelEncryptionConfigsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListFieldLevelEncryptionConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFieldLevelEncryptionConfigsRequestDescriptor =
    $convert.base64Decode(
        'CiZMaXN0RmllbGRMZXZlbEVuY3J5cHRpb25Db25maWdzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKi'
        'ABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated('Use listFieldLevelEncryptionConfigsResultDescriptor instead')
const ListFieldLevelEncryptionConfigsResult$json = {
  '1': 'ListFieldLevelEncryptionConfigsResult',
  '2': [
    {
      '1': 'fieldlevelencryptionlist',
      '3': 364178743,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionList',
      '8': {},
      '10': 'fieldlevelencryptionlist'
    },
  ],
};

/// Descriptor for `ListFieldLevelEncryptionConfigsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFieldLevelEncryptionConfigsResultDescriptor =
    $convert.base64Decode(
        'CiVMaXN0RmllbGRMZXZlbEVuY3J5cHRpb25Db25maWdzUmVzdWx0EmoKGGZpZWxkbGV2ZWxlbm'
        'NyeXB0aW9ubGlzdBi32tOtASABKAsyJC5jbG91ZGZyb250LkZpZWxkTGV2ZWxFbmNyeXB0aW9u'
        'TGlzdEIEiLUYAVIYZmllbGRsZXZlbGVuY3J5cHRpb25saXN0');

@$core
    .Deprecated('Use listFieldLevelEncryptionProfilesRequestDescriptor instead')
const ListFieldLevelEncryptionProfilesRequest$json = {
  '1': 'ListFieldLevelEncryptionProfilesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListFieldLevelEncryptionProfilesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFieldLevelEncryptionProfilesRequestDescriptor =
    $convert.base64Decode(
        'CidMaXN0RmllbGRMZXZlbEVuY3J5cHRpb25Qcm9maWxlc1JlcXVlc3QSGQoGbWFya2VyGLjdzS'
        'ogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core
    .Deprecated('Use listFieldLevelEncryptionProfilesResultDescriptor instead')
const ListFieldLevelEncryptionProfilesResult$json = {
  '1': 'ListFieldLevelEncryptionProfilesResult',
  '2': [
    {
      '1': 'fieldlevelencryptionprofilelist',
      '3': 405827328,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfileList',
      '8': {},
      '10': 'fieldlevelencryptionprofilelist'
    },
  ],
};

/// Descriptor for `ListFieldLevelEncryptionProfilesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFieldLevelEncryptionProfilesResultDescriptor =
    $convert.base64Decode(
        'CiZMaXN0RmllbGRMZXZlbEVuY3J5cHRpb25Qcm9maWxlc1Jlc3VsdBJ/Ch9maWVsZGxldmVsZW'
        '5jcnlwdGlvbnByb2ZpbGVsaXN0GIDewcEBIAEoCzIrLmNsb3VkZnJvbnQuRmllbGRMZXZlbEVu'
        'Y3J5cHRpb25Qcm9maWxlTGlzdEIEiLUYAVIfZmllbGRsZXZlbGVuY3J5cHRpb25wcm9maWxlbG'
        'lzdA==');

@$core.Deprecated('Use listFunctionsRequestDescriptor instead')
const ListFunctionsRequest$json = {
  '1': 'ListFunctionsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `ListFunctionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0RnVuY3Rpb25zUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYX'
    'hpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEjIKBXN0YWdlGM6X2HAgASgOMhkuY2xvdWRmcm9u'
    'dC5GdW5jdGlvblN0YWdlUgVzdGFnZQ==');

@$core.Deprecated('Use listFunctionsResultDescriptor instead')
const ListFunctionsResult$json = {
  '1': 'ListFunctionsResult',
  '2': [
    {
      '1': 'functionlist',
      '3': 73245974,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionList',
      '8': {},
      '10': 'functionlist'
    },
  ],
};

/// Descriptor for `ListFunctionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionsResultDescriptor = $convert.base64Decode(
    'ChNMaXN0RnVuY3Rpb25zUmVzdWx0EkUKDGZ1bmN0aW9ubGlzdBiWyvYiIAEoCzIYLmNsb3VkZn'
    'JvbnQuRnVuY3Rpb25MaXN0QgSItRgBUgxmdW5jdGlvbmxpc3Q=');

@$core.Deprecated(
    'Use listInvalidationsForDistributionTenantRequestDescriptor instead')
const ListInvalidationsForDistributionTenantRequest$json = {
  '1': 'ListInvalidationsForDistributionTenantRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListInvalidationsForDistributionTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listInvalidationsForDistributionTenantRequestDescriptor =
    $convert.base64Decode(
        'Ci1MaXN0SW52YWxpZGF0aW9uc0ZvckRpc3RyaWJ1dGlvblRlbmFudFJlcXVlc3QSEgoCaWQYgf'
        'KitwEgASgJUgJpZBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trx'
        'ASABKAVSCG1heGl0ZW1z');

@$core.Deprecated(
    'Use listInvalidationsForDistributionTenantResultDescriptor instead')
const ListInvalidationsForDistributionTenantResult$json = {
  '1': 'ListInvalidationsForDistributionTenantResult',
  '2': [
    {
      '1': 'invalidationlist',
      '3': 58264922,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.InvalidationList',
      '8': {},
      '10': 'invalidationlist'
    },
  ],
};

/// Descriptor for `ListInvalidationsForDistributionTenantResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listInvalidationsForDistributionTenantResultDescriptor =
    $convert.base64Decode(
        'CixMaXN0SW52YWxpZGF0aW9uc0ZvckRpc3RyaWJ1dGlvblRlbmFudFJlc3VsdBJRChBpbnZhbG'
        'lkYXRpb25saXN0GNqa5BsgASgLMhwuY2xvdWRmcm9udC5JbnZhbGlkYXRpb25MaXN0QgSItRgB'
        'UhBpbnZhbGlkYXRpb25saXN0');

@$core.Deprecated('Use listInvalidationsRequestDescriptor instead')
const ListInvalidationsRequest$json = {
  '1': 'ListInvalidationsRequest',
  '2': [
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListInvalidationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInvalidationsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0SW52YWxpZGF0aW9uc1JlcXVlc3QSKQoOZGlzdHJpYnV0aW9uaWQY57H7QyABKAlSDm'
    'Rpc3RyaWJ1dGlvbmlkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW'
    '2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listInvalidationsResultDescriptor instead')
const ListInvalidationsResult$json = {
  '1': 'ListInvalidationsResult',
  '2': [
    {
      '1': 'invalidationlist',
      '3': 58264922,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.InvalidationList',
      '8': {},
      '10': 'invalidationlist'
    },
  ],
};

/// Descriptor for `ListInvalidationsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInvalidationsResultDescriptor = $convert.base64Decode(
    'ChdMaXN0SW52YWxpZGF0aW9uc1Jlc3VsdBJRChBpbnZhbGlkYXRpb25saXN0GNqa5BsgASgLMh'
    'wuY2xvdWRmcm9udC5JbnZhbGlkYXRpb25MaXN0QgSItRgBUhBpbnZhbGlkYXRpb25saXN0');

@$core.Deprecated('Use listKeyGroupsRequestDescriptor instead')
const ListKeyGroupsRequest$json = {
  '1': 'ListKeyGroupsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListKeyGroupsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyGroupsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0S2V5R3JvdXBzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYX'
    'hpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated('Use listKeyGroupsResultDescriptor instead')
const ListKeyGroupsResult$json = {
  '1': 'ListKeyGroupsResult',
  '2': [
    {
      '1': 'keygrouplist',
      '3': 349267816,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroupList',
      '8': {},
      '10': 'keygrouplist'
    },
  ],
};

/// Descriptor for `ListKeyGroupsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyGroupsResultDescriptor = $convert.base64Decode(
    'ChNMaXN0S2V5R3JvdXBzUmVzdWx0EkYKDGtleWdyb3VwbGlzdBjozsWmASABKAsyGC5jbG91ZG'
    'Zyb250LktleUdyb3VwTGlzdEIEiLUYAVIMa2V5Z3JvdXBsaXN0');

@$core.Deprecated('Use listKeyValueStoresRequestDescriptor instead')
const ListKeyValueStoresRequest$json = {
  '1': 'ListKeyValueStoresRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `ListKeyValueStoresRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyValueStoresRequestDescriptor = $convert.base64Decode(
    'ChlMaXN0S2V5VmFsdWVTdG9yZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh'
    '4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSGQoGc3RhdHVzGJDk+wIgASgJUgZzdGF0'
    'dXM=');

@$core.Deprecated('Use listKeyValueStoresResultDescriptor instead')
const ListKeyValueStoresResult$json = {
  '1': 'ListKeyValueStoresResult',
  '2': [
    {
      '1': 'keyvaluestorelist',
      '3': 41649787,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyValueStoreList',
      '8': {},
      '10': 'keyvaluestorelist'
    },
  ],
};

/// Descriptor for `ListKeyValueStoresResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyValueStoresResultDescriptor = $convert.base64Decode(
    'ChhMaXN0S2V5VmFsdWVTdG9yZXNSZXN1bHQSVAoRa2V5dmFsdWVzdG9yZWxpc3QY+4zuEyABKA'
    'syHS5jbG91ZGZyb250LktleVZhbHVlU3RvcmVMaXN0QgSItRgBUhFrZXl2YWx1ZXN0b3JlbGlz'
    'dA==');

@$core.Deprecated('Use listOriginAccessControlsRequestDescriptor instead')
const ListOriginAccessControlsRequest$json = {
  '1': 'ListOriginAccessControlsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListOriginAccessControlsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOriginAccessControlsRequestDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0T3JpZ2luQWNjZXNzQ29udHJvbHNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbW'
        'Fya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listOriginAccessControlsResultDescriptor instead')
const ListOriginAccessControlsResult$json = {
  '1': 'ListOriginAccessControlsResult',
  '2': [
    {
      '1': 'originaccesscontrollist',
      '3': 140052627,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControlList',
      '8': {},
      '10': 'originaccesscontrollist'
    },
  ],
};

/// Descriptor for `ListOriginAccessControlsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOriginAccessControlsResultDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0T3JpZ2luQWNjZXNzQ29udHJvbHNSZXN1bHQSZgoXb3JpZ2luYWNjZXNzY29udHJvbG'
        'xpc3QYk5HkQiABKAsyIy5jbG91ZGZyb250Lk9yaWdpbkFjY2Vzc0NvbnRyb2xMaXN0QgSItRgB'
        'UhdvcmlnaW5hY2Nlc3Njb250cm9sbGlzdA==');

@$core.Deprecated('Use listOriginRequestPoliciesRequestDescriptor instead')
const ListOriginRequestPoliciesRequest$json = {
  '1': 'ListOriginRequestPoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginRequestPolicyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ListOriginRequestPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOriginRequestPoliciesRequestDescriptor =
    $convert.base64Decode(
        'CiBMaXN0T3JpZ2luUmVxdWVzdFBvbGljaWVzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm'
        '1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEjsKBHR5cGUY7qDXigEgASgO'
        'MiMuY2xvdWRmcm9udC5PcmlnaW5SZXF1ZXN0UG9saWN5VHlwZVIEdHlwZQ==');

@$core.Deprecated('Use listOriginRequestPoliciesResultDescriptor instead')
const ListOriginRequestPoliciesResult$json = {
  '1': 'ListOriginRequestPoliciesResult',
  '2': [
    {
      '1': 'originrequestpolicylist',
      '3': 464628663,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyList',
      '8': {},
      '10': 'originrequestpolicylist'
    },
  ],
};

/// Descriptor for `ListOriginRequestPoliciesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOriginRequestPoliciesResultDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0T3JpZ2luUmVxdWVzdFBvbGljaWVzUmVzdWx0EmcKF29yaWdpbnJlcXVlc3Rwb2xpY3'
        'lsaXN0GLfXxt0BIAEoCzIjLmNsb3VkZnJvbnQuT3JpZ2luUmVxdWVzdFBvbGljeUxpc3RCBIi1'
        'GAFSF29yaWdpbnJlcXVlc3Rwb2xpY3lsaXN0');

@$core.Deprecated('Use listPublicKeysRequestDescriptor instead')
const ListPublicKeysRequest$json = {
  '1': 'ListPublicKeysRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListPublicKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPublicKeysRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0UHVibGljS2V5c1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbW'
    'F4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listPublicKeysResultDescriptor instead')
const ListPublicKeysResult$json = {
  '1': 'ListPublicKeysResult',
  '2': [
    {
      '1': 'publickeylist',
      '3': 129933128,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKeyList',
      '8': {},
      '10': 'publickeylist'
    },
  ],
};

/// Descriptor for `ListPublicKeysResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPublicKeysResultDescriptor = $convert.base64Decode(
    'ChRMaXN0UHVibGljS2V5c1Jlc3VsdBJICg1wdWJsaWNrZXlsaXN0GMi++j0gASgLMhkuY2xvdW'
    'Rmcm9udC5QdWJsaWNLZXlMaXN0QgSItRgBUg1wdWJsaWNrZXlsaXN0');

@$core.Deprecated('Use listRealtimeLogConfigsRequestDescriptor instead')
const ListRealtimeLogConfigsRequest$json = {
  '1': 'ListRealtimeLogConfigsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListRealtimeLogConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRealtimeLogConfigsRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0UmVhbHRpbWVMb2dDb25maWdzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcm'
        'tlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated('Use listRealtimeLogConfigsResultDescriptor instead')
const ListRealtimeLogConfigsResult$json = {
  '1': 'ListRealtimeLogConfigsResult',
  '2': [
    {
      '1': 'realtimelogconfigs',
      '3': 453095746,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.RealtimeLogConfigs',
      '8': {},
      '10': 'realtimelogconfigs'
    },
  ],
};

/// Descriptor for `ListRealtimeLogConfigsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRealtimeLogConfigsResultDescriptor =
    $convert.base64Decode(
        'ChxMaXN0UmVhbHRpbWVMb2dDb25maWdzUmVzdWx0ElgKEnJlYWx0aW1lbG9nY29uZmlncxjC4o'
        'bYASABKAsyHi5jbG91ZGZyb250LlJlYWx0aW1lTG9nQ29uZmlnc0IEiLUYAVIScmVhbHRpbWVs'
        'b2djb25maWdz');

@$core.Deprecated('Use listResponseHeadersPoliciesRequestDescriptor instead')
const ListResponseHeadersPoliciesRequest$json = {
  '1': 'ListResponseHeadersPoliciesRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ResponseHeadersPolicyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ListResponseHeadersPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResponseHeadersPoliciesRequestDescriptor =
    $convert.base64Decode(
        'CiJMaXN0UmVzcG9uc2VIZWFkZXJzUG9saWNpZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCV'
        'IGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSPQoEdHlwZRjuoNeKASAB'
        'KA4yJS5jbG91ZGZyb250LlJlc3BvbnNlSGVhZGVyc1BvbGljeVR5cGVSBHR5cGU=');

@$core.Deprecated('Use listResponseHeadersPoliciesResultDescriptor instead')
const ListResponseHeadersPoliciesResult$json = {
  '1': 'ListResponseHeadersPoliciesResult',
  '2': [
    {
      '1': 'responseheaderspolicylist',
      '3': 10074075,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyList',
      '8': {},
      '10': 'responseheaderspolicylist'
    },
  ],
};

/// Descriptor for `ListResponseHeadersPoliciesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResponseHeadersPoliciesResultDescriptor =
    $convert.base64Decode(
        'CiFMaXN0UmVzcG9uc2VIZWFkZXJzUG9saWNpZXNSZXN1bHQSbAoZcmVzcG9uc2VoZWFkZXJzcG'
        '9saWN5bGlzdBjb7+YEIAEoCzIlLmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9saWN5TGlz'
        'dEIEiLUYAVIZcmVzcG9uc2VoZWFkZXJzcG9saWN5bGlzdA==');

@$core.Deprecated('Use listStreamingDistributionsRequestDescriptor instead')
const ListStreamingDistributionsRequest$json = {
  '1': 'ListStreamingDistributionsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListStreamingDistributionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStreamingDistributionsRequestDescriptor =
    $convert.base64Decode(
        'CiFMaXN0U3RyZWFtaW5nRGlzdHJpYnV0aW9uc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUg'
        'ZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listStreamingDistributionsResultDescriptor instead')
const ListStreamingDistributionsResult$json = {
  '1': 'ListStreamingDistributionsResult',
  '2': [
    {
      '1': 'streamingdistributionlist',
      '3': 80532578,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionList',
      '8': {},
      '10': 'streamingdistributionlist'
    },
  ],
};

/// Descriptor for `ListStreamingDistributionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStreamingDistributionsResultDescriptor =
    $convert.base64Decode(
        'CiBMaXN0U3RyZWFtaW5nRGlzdHJpYnV0aW9uc1Jlc3VsdBJsChlzdHJlYW1pbmdkaXN0cmlidX'
        'Rpb25saXN0GOKosyYgASgLMiUuY2xvdWRmcm9udC5TdHJlYW1pbmdEaXN0cmlidXRpb25MaXN0'
        'QgSItRgBUhlzdHJlYW1pbmdkaXN0cmlidXRpb25saXN0');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resource', '3': 61921302, '4': 1, '5': 9, '10': 'resource'},
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIdCghyZXNvdXJjZRiWsMMdIAEoCVIIcmVzb3'
        'VyY2U=');

@$core.Deprecated('Use listTagsForResourceResultDescriptor instead')
const ListTagsForResourceResult$json = {
  '1': 'ListTagsForResourceResult',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '8': {},
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResultDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFnc0ZvclJlc291cmNlUmVzdWx0Ei4KBHRhZ3MYwcH2tQEgASgLMhAuY2xvdWRmcm'
        '9udC5UYWdzQgSItRgBUgR0YWdz');

@$core.Deprecated('Use listTrustStoresRequestDescriptor instead')
const ListTrustStoresRequest$json = {
  '1': 'ListTrustStoresRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListTrustStoresRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrustStoresRequestDescriptor =
    $convert.base64Decode(
        'ChZMaXN0VHJ1c3RTdG9yZXNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG'
        '1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listTrustStoresResultDescriptor instead')
const ListTrustStoresResult$json = {
  '1': 'ListTrustStoresResult',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {
      '1': 'truststorelist',
      '3': 73369611,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.TrustStoreSummary',
      '10': 'truststorelist'
    },
  ],
};

/// Descriptor for `ListTrustStoresResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrustStoresResultDescriptor = $convert.base64Decode(
    'ChVMaXN0VHJ1c3RTdG9yZXNSZXN1bHQSIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leHRtYX'
    'JrZXISSAoOdHJ1c3RzdG9yZWxpc3QYi5D+IiADKAsyHS5jbG91ZGZyb250LlRydXN0U3RvcmVT'
    'dW1tYXJ5Ug50cnVzdHN0b3JlbGlzdA==');

@$core.Deprecated('Use listVpcOriginsRequestDescriptor instead')
const ListVpcOriginsRequest$json = {
  '1': 'ListVpcOriginsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListVpcOriginsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVpcOriginsRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0VnBjT3JpZ2luc1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbW'
    'F4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listVpcOriginsResultDescriptor instead')
const ListVpcOriginsResult$json = {
  '1': 'ListVpcOriginsResult',
  '2': [
    {
      '1': 'vpcoriginlist',
      '3': 177956775,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOriginList',
      '8': {},
      '10': 'vpcoriginlist'
    },
  ],
};

/// Descriptor for `ListVpcOriginsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVpcOriginsResultDescriptor = $convert.base64Decode(
    'ChRMaXN0VnBjT3JpZ2luc1Jlc3VsdBJICg12cGNvcmlnaW5saXN0GKfP7VQgASgLMhkuY2xvdW'
    'Rmcm9udC5WcGNPcmlnaW5MaXN0QgSItRgBUg12cGNvcmlnaW5saXN0');

@$core.Deprecated('Use loggingConfigDescriptor instead')
const LoggingConfig$json = {
  '1': 'LoggingConfig',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'includecookies',
      '3': 385372575,
      '4': 1,
      '5': 8,
      '10': 'includecookies'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `LoggingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loggingConfigDescriptor = $convert.base64Decode(
    'Cg1Mb2dnaW5nQ29uZmlnEhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EhwKB2VuYWJsZWQYv8'
    'ib5AEgASgIUgdlbmFibGVkEioKDmluY2x1ZGVjb29raWVzGJ+j4bcBIAEoCFIOaW5jbHVkZWNv'
    'b2tpZXMSGgoGcHJlZml4GOqz04IBIAEoCVIGcHJlZml4');

@$core.Deprecated('Use managedCertificateDetailsDescriptor instead')
const ManagedCertificateDetails$json = {
  '1': 'ManagedCertificateDetails',
  '2': [
    {
      '1': 'certificatearn',
      '3': 92693880,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'certificatestatus',
      '3': 508638899,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ManagedCertificateStatus',
      '10': 'certificatestatus'
    },
    {
      '1': 'validationtokendetails',
      '3': 370214280,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ValidationTokenDetail',
      '10': 'validationtokendetails'
    },
    {
      '1': 'validationtokenhost',
      '3': 521935290,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ValidationTokenHost',
      '10': 'validationtokenhost'
    },
  ],
};

/// Descriptor for `ManagedCertificateDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedCertificateDetailsDescriptor = $convert.base64Decode(
    'ChlNYW5hZ2VkQ2VydGlmaWNhdGVEZXRhaWxzEikKDmNlcnRpZmljYXRlYXJuGPjKmSwgASgJUg'
    '5jZXJ0aWZpY2F0ZWFybhJWChFjZXJ0aWZpY2F0ZXN0YXR1cxiz7cTyASABKA4yJC5jbG91ZGZy'
    'b250Lk1hbmFnZWRDZXJ0aWZpY2F0ZVN0YXR1c1IRY2VydGlmaWNhdGVzdGF0dXMSXQoWdmFsaW'
    'RhdGlvbnRva2VuZGV0YWlscxiIi8SwASADKAsyIS5jbG91ZGZyb250LlZhbGlkYXRpb25Ub2tl'
    'bkRldGFpbFIWdmFsaWRhdGlvbnRva2VuZGV0YWlscxJVChN2YWxpZGF0aW9udG9rZW5ob3N0GL'
    'qz8PgBIAEoDjIfLmNsb3VkZnJvbnQuVmFsaWRhdGlvblRva2VuSG9zdFITdmFsaWRhdGlvbnRv'
    'a2VuaG9zdA==');

@$core.Deprecated('Use managedCertificateRequestDescriptor instead')
const ManagedCertificateRequest$json = {
  '1': 'ManagedCertificateRequest',
  '2': [
    {
      '1': 'certificatetransparencyloggingpreference',
      '3': 414636075,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CertificateTransparencyLoggingPreference',
      '10': 'certificatetransparencyloggingpreference'
    },
    {
      '1': 'primarydomainname',
      '3': 229268599,
      '4': 1,
      '5': 9,
      '10': 'primarydomainname'
    },
    {
      '1': 'validationtokenhost',
      '3': 521935290,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ValidationTokenHost',
      '10': 'validationtokenhost'
    },
  ],
};

/// Descriptor for `ManagedCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedCertificateRequestDescriptor = $convert.base64Decode(
    'ChlNYW5hZ2VkQ2VydGlmaWNhdGVSZXF1ZXN0EpQBCihjZXJ0aWZpY2F0ZXRyYW5zcGFyZW5jeW'
    'xvZ2dpbmdwcmVmZXJlbmNlGKuw28UBIAEoDjI0LmNsb3VkZnJvbnQuQ2VydGlmaWNhdGVUcmFu'
    'c3BhcmVuY3lMb2dnaW5nUHJlZmVyZW5jZVIoY2VydGlmaWNhdGV0cmFuc3BhcmVuY3lsb2dnaW'
    '5ncHJlZmVyZW5jZRIvChFwcmltYXJ5ZG9tYWlubmFtZRj3uKltIAEoCVIRcHJpbWFyeWRvbWFp'
    'bm5hbWUSVQoTdmFsaWRhdGlvbnRva2VuaG9zdBi6s/D4ASABKA4yHy5jbG91ZGZyb250LlZhbG'
    'lkYXRpb25Ub2tlbkhvc3RSE3ZhbGlkYXRpb250b2tlbmhvc3Q=');

@$core.Deprecated('Use missingBodyDescriptor instead')
const MissingBody$json = {
  '1': 'MissingBody',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MissingBody`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List missingBodyDescriptor = $convert
    .base64Decode('CgtNaXNzaW5nQm9keRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use monitoringSubscriptionDescriptor instead')
const MonitoringSubscription$json = {
  '1': 'MonitoringSubscription',
  '2': [
    {
      '1': 'realtimemetricssubscriptionconfig',
      '3': 5061935,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.RealtimeMetricsSubscriptionConfig',
      '10': 'realtimemetricssubscriptionconfig'
    },
  ],
};

/// Descriptor for `MonitoringSubscription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List monitoringSubscriptionDescriptor = $convert.base64Decode(
    'ChZNb25pdG9yaW5nU3Vic2NyaXB0aW9uEn4KIXJlYWx0aW1lbWV0cmljc3N1YnNjcmlwdGlvbm'
    'NvbmZpZxiv+rQCIAEoCzItLmNsb3VkZnJvbnQuUmVhbHRpbWVNZXRyaWNzU3Vic2NyaXB0aW9u'
    'Q29uZmlnUiFyZWFsdGltZW1ldHJpY3NzdWJzY3JpcHRpb25jb25maWc=');

@$core.Deprecated('Use monitoringSubscriptionAlreadyExistsDescriptor instead')
const MonitoringSubscriptionAlreadyExists$json = {
  '1': 'MonitoringSubscriptionAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MonitoringSubscriptionAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List monitoringSubscriptionAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CiNNb25pdG9yaW5nU3Vic2NyaXB0aW9uQWxyZWFkeUV4aXN0cxIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use noSuchCachePolicyDescriptor instead')
const NoSuchCachePolicy$json = {
  '1': 'NoSuchCachePolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchCachePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchCachePolicyDescriptor = $convert.base64Decode(
    'ChFOb1N1Y2hDYWNoZVBvbGljeRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use noSuchCloudFrontOriginAccessIdentityDescriptor instead')
const NoSuchCloudFrontOriginAccessIdentity$json = {
  '1': 'NoSuchCloudFrontOriginAccessIdentity',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchCloudFrontOriginAccessIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchCloudFrontOriginAccessIdentityDescriptor =
    $convert.base64Decode(
        'CiROb1N1Y2hDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHkSGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use noSuchContinuousDeploymentPolicyDescriptor instead')
const NoSuchContinuousDeploymentPolicy$json = {
  '1': 'NoSuchContinuousDeploymentPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchContinuousDeploymentPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchContinuousDeploymentPolicyDescriptor =
    $convert.base64Decode(
        'CiBOb1N1Y2hDb250aW51b3VzRGVwbG95bWVudFBvbGljeRIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use noSuchDistributionDescriptor instead')
const NoSuchDistribution$json = {
  '1': 'NoSuchDistribution',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchDistribution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchDistributionDescriptor =
    $convert.base64Decode(
        'ChJOb1N1Y2hEaXN0cmlidXRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use noSuchFieldLevelEncryptionConfigDescriptor instead')
const NoSuchFieldLevelEncryptionConfig$json = {
  '1': 'NoSuchFieldLevelEncryptionConfig',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchFieldLevelEncryptionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchFieldLevelEncryptionConfigDescriptor =
    $convert.base64Decode(
        'CiBOb1N1Y2hGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZxIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use noSuchFieldLevelEncryptionProfileDescriptor instead')
const NoSuchFieldLevelEncryptionProfile$json = {
  '1': 'NoSuchFieldLevelEncryptionProfile',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchFieldLevelEncryptionProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchFieldLevelEncryptionProfileDescriptor =
    $convert.base64Decode(
        'CiFOb1N1Y2hGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGUSGwoHbWVzc2FnZRiFs7twIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use noSuchFunctionExistsDescriptor instead')
const NoSuchFunctionExists$json = {
  '1': 'NoSuchFunctionExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchFunctionExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchFunctionExistsDescriptor =
    $convert.base64Decode(
        'ChROb1N1Y2hGdW5jdGlvbkV4aXN0cxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use noSuchInvalidationDescriptor instead')
const NoSuchInvalidation$json = {
  '1': 'NoSuchInvalidation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchInvalidation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchInvalidationDescriptor =
    $convert.base64Decode(
        'ChJOb1N1Y2hJbnZhbGlkYXRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use noSuchMonitoringSubscriptionDescriptor instead')
const NoSuchMonitoringSubscription$json = {
  '1': 'NoSuchMonitoringSubscription',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchMonitoringSubscription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchMonitoringSubscriptionDescriptor =
    $convert.base64Decode(
        'ChxOb1N1Y2hNb25pdG9yaW5nU3Vic2NyaXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use noSuchOriginDescriptor instead')
const NoSuchOrigin$json = {
  '1': 'NoSuchOrigin',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchOrigin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchOriginDescriptor = $convert.base64Decode(
    'CgxOb1N1Y2hPcmlnaW4SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use noSuchOriginAccessControlDescriptor instead')
const NoSuchOriginAccessControl$json = {
  '1': 'NoSuchOriginAccessControl',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchOriginAccessControl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchOriginAccessControlDescriptor =
    $convert.base64Decode(
        'ChlOb1N1Y2hPcmlnaW5BY2Nlc3NDb250cm9sEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use noSuchOriginRequestPolicyDescriptor instead')
const NoSuchOriginRequestPolicy$json = {
  '1': 'NoSuchOriginRequestPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchOriginRequestPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchOriginRequestPolicyDescriptor =
    $convert.base64Decode(
        'ChlOb1N1Y2hPcmlnaW5SZXF1ZXN0UG9saWN5EhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use noSuchPublicKeyDescriptor instead')
const NoSuchPublicKey$json = {
  '1': 'NoSuchPublicKey',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchPublicKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchPublicKeyDescriptor = $convert.base64Decode(
    'Cg9Ob1N1Y2hQdWJsaWNLZXkSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use noSuchRealtimeLogConfigDescriptor instead')
const NoSuchRealtimeLogConfig$json = {
  '1': 'NoSuchRealtimeLogConfig',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchRealtimeLogConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchRealtimeLogConfigDescriptor =
    $convert.base64Decode(
        'ChdOb1N1Y2hSZWFsdGltZUxvZ0NvbmZpZxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use noSuchResourceDescriptor instead')
const NoSuchResource$json = {
  '1': 'NoSuchResource',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchResource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchResourceDescriptor = $convert.base64Decode(
    'Cg5Ob1N1Y2hSZXNvdXJjZRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use noSuchResponseHeadersPolicyDescriptor instead')
const NoSuchResponseHeadersPolicy$json = {
  '1': 'NoSuchResponseHeadersPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchResponseHeadersPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchResponseHeadersPolicyDescriptor =
    $convert.base64Decode(
        'ChtOb1N1Y2hSZXNwb25zZUhlYWRlcnNQb2xpY3kSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use noSuchStreamingDistributionDescriptor instead')
const NoSuchStreamingDistribution$json = {
  '1': 'NoSuchStreamingDistribution',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchStreamingDistribution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchStreamingDistributionDescriptor =
    $convert.base64Decode(
        'ChtOb1N1Y2hTdHJlYW1pbmdEaXN0cmlidXRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use originDescriptor instead')
const Origin$json = {
  '1': 'Origin',
  '2': [
    {
      '1': 'connectionattempts',
      '3': 130111316,
      '4': 1,
      '5': 5,
      '10': 'connectionattempts'
    },
    {
      '1': 'connectiontimeout',
      '3': 285150581,
      '4': 1,
      '5': 5,
      '10': 'connectiontimeout'
    },
    {
      '1': 'customheaders',
      '3': 112565669,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CustomHeaders',
      '10': 'customheaders'
    },
    {
      '1': 'customoriginconfig',
      '3': 97020983,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CustomOriginConfig',
      '10': 'customoriginconfig'
    },
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'originaccesscontrolid',
      '3': 344049692,
      '4': 1,
      '5': 9,
      '10': 'originaccesscontrolid'
    },
    {'1': 'originpath', '3': 66064073, '4': 1, '5': 9, '10': 'originpath'},
    {
      '1': 'originshield',
      '3': 482702287,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginShield',
      '10': 'originshield'
    },
    {
      '1': 'responsecompletiontimeout',
      '3': 197954616,
      '4': 1,
      '5': 5,
      '10': 'responsecompletiontimeout'
    },
    {
      '1': 's3originconfig',
      '3': 78562208,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.S3OriginConfig',
      '10': 's3originconfig'
    },
    {
      '1': 'vpcoriginconfig',
      '3': 107511973,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOriginConfig',
      '10': 'vpcoriginconfig'
    },
  ],
};

/// Descriptor for `Origin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originDescriptor = $convert.base64Decode(
    'CgZPcmlnaW4SMQoSY29ubmVjdGlvbmF0dGVtcHRzGNSuhT4gASgFUhJjb25uZWN0aW9uYXR0ZW'
    '1wdHMSMAoRY29ubmVjdGlvbnRpbWVvdXQY9Zr8hwEgASgFUhFjb25uZWN0aW9udGltZW91dBJC'
    'Cg1jdXN0b21oZWFkZXJzGKW71jUgASgLMhkuY2xvdWRmcm9udC5DdXN0b21IZWFkZXJzUg1jdX'
    'N0b21oZWFkZXJzElEKEmN1c3RvbW9yaWdpbmNvbmZpZxi32KEuIAEoCzIeLmNsb3VkZnJvbnQu'
    'Q3VzdG9tT3JpZ2luQ29uZmlnUhJjdXN0b21vcmlnaW5jb25maWcSIQoKZG9tYWlubmFtZRjrzf'
    'hcIAEoCVIKZG9tYWlubmFtZRISCgJpZBiB8qK3ASABKAlSAmlkEjgKFW9yaWdpbmFjY2Vzc2Nv'
    'bnRyb2xpZBickIekASABKAlSFW9yaWdpbmFjY2Vzc2NvbnRyb2xpZBIhCgpvcmlnaW5wYXRoGM'
    'mdwB8gASgJUgpvcmlnaW5wYXRoEkAKDG9yaWdpbnNoaWVsZBjP55XmASABKAsyGC5jbG91ZGZy'
    'b250Lk9yaWdpblNoaWVsZFIMb3JpZ2luc2hpZWxkEj8KGXJlc3BvbnNlY29tcGxldGlvbnRpbW'
    'VvdXQYuJiyXiABKAVSGXJlc3BvbnNlY29tcGxldGlvbnRpbWVvdXQSRQoOczNvcmlnaW5jb25m'
    'aWcYoIe7JSABKAsyGi5jbG91ZGZyb250LlMzT3JpZ2luQ29uZmlnUg5zM29yaWdpbmNvbmZpZx'
    'JICg92cGNvcmlnaW5jb25maWcYpYGiMyABKAsyGy5jbG91ZGZyb250LlZwY09yaWdpbkNvbmZp'
    'Z1IPdnBjb3JpZ2luY29uZmln');

@$core.Deprecated('Use originAccessControlDescriptor instead')
const OriginAccessControl$json = {
  '1': 'OriginAccessControl',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'originaccesscontrolconfig',
      '3': 143834977,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControlConfig',
      '10': 'originaccesscontrolconfig'
    },
  ],
};

/// Descriptor for `OriginAccessControl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originAccessControlDescriptor = $convert.base64Decode(
    'ChNPcmlnaW5BY2Nlc3NDb250cm9sEhIKAmlkGIHyorcBIAEoCVICaWQSZgoZb3JpZ2luYWNjZX'
    'NzY29udHJvbGNvbmZpZxjh/spEIAEoCzIlLmNsb3VkZnJvbnQuT3JpZ2luQWNjZXNzQ29udHJv'
    'bENvbmZpZ1IZb3JpZ2luYWNjZXNzY29udHJvbGNvbmZpZw==');

@$core.Deprecated('Use originAccessControlAlreadyExistsDescriptor instead')
const OriginAccessControlAlreadyExists$json = {
  '1': 'OriginAccessControlAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OriginAccessControlAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originAccessControlAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CiBPcmlnaW5BY2Nlc3NDb250cm9sQWxyZWFkeUV4aXN0cxIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use originAccessControlConfigDescriptor instead')
const OriginAccessControlConfig$json = {
  '1': 'OriginAccessControlConfig',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'originaccesscontrolorigintype',
      '3': 496410329,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginAccessControlOriginTypes',
      '10': 'originaccesscontrolorigintype'
    },
    {
      '1': 'signingbehavior',
      '3': 533461047,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginAccessControlSigningBehaviors',
      '10': 'signingbehavior'
    },
    {
      '1': 'signingprotocol',
      '3': 341451169,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginAccessControlSigningProtocols',
      '10': 'signingprotocol'
    },
  ],
};

/// Descriptor for `OriginAccessControlConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originAccessControlConfigDescriptor = $convert.base64Decode(
    'ChlPcmlnaW5BY2Nlc3NDb250cm9sQ29uZmlnEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZX'
    'NjcmlwdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEnQKHW9yaWdpbmFjY2Vzc2NvbnRyb2xv'
    'cmlnaW50eXBlGNm92uwBIAEoDjIqLmNsb3VkZnJvbnQuT3JpZ2luQWNjZXNzQ29udHJvbE9yaW'
    'dpblR5cGVzUh1vcmlnaW5hY2Nlc3Njb250cm9sb3JpZ2ludHlwZRJdCg9zaWduaW5nYmVoYXZp'
    'b3IYt/Cv/gEgASgOMi8uY2xvdWRmcm9udC5PcmlnaW5BY2Nlc3NDb250cm9sU2lnbmluZ0JlaG'
    'F2aW9yc1IPc2lnbmluZ2JlaGF2aW9yEl0KD3NpZ25pbmdwcm90b2NvbBihw+iiASABKA4yLy5j'
    'bG91ZGZyb250Lk9yaWdpbkFjY2Vzc0NvbnRyb2xTaWduaW5nUHJvdG9jb2xzUg9zaWduaW5ncH'
    'JvdG9jb2w=');

@$core.Deprecated('Use originAccessControlInUseDescriptor instead')
const OriginAccessControlInUse$json = {
  '1': 'OriginAccessControlInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OriginAccessControlInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originAccessControlInUseDescriptor =
    $convert.base64Decode(
        'ChhPcmlnaW5BY2Nlc3NDb250cm9sSW5Vc2USGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use originAccessControlListDescriptor instead')
const OriginAccessControlList$json = {
  '1': 'OriginAccessControlList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.OriginAccessControlSummary',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `OriginAccessControlList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originAccessControlListDescriptor = $convert.base64Decode(
    'ChdPcmlnaW5BY2Nlc3NDb250cm9sTGlzdBIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cn'
    'VuY2F0ZWQSPwoFaXRlbXMYsPDYASADKAsyJi5jbG91ZGZyb250Lk9yaWdpbkFjY2Vzc0NvbnRy'
    'b2xTdW1tYXJ5UgVpdGVtcxIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcx'
    'iU1trxASABKAVSCG1heGl0ZW1zEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2Vy'
    'Eh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use originAccessControlSummaryDescriptor instead')
const OriginAccessControlSummary$json = {
  '1': 'OriginAccessControlSummary',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'originaccesscontrolorigintype',
      '3': 496410329,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginAccessControlOriginTypes',
      '10': 'originaccesscontrolorigintype'
    },
    {
      '1': 'signingbehavior',
      '3': 533461047,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginAccessControlSigningBehaviors',
      '10': 'signingbehavior'
    },
    {
      '1': 'signingprotocol',
      '3': 341451169,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginAccessControlSigningProtocols',
      '10': 'signingprotocol'
    },
  ],
};

/// Descriptor for `OriginAccessControlSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originAccessControlSummaryDescriptor = $convert.base64Decode(
    'ChpPcmlnaW5BY2Nlc3NDb250cm9sU3VtbWFyeRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZG'
    'VzY3JpcHRpb24SEgoCaWQYgfKitwEgASgJUgJpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEnQK'
    'HW9yaWdpbmFjY2Vzc2NvbnRyb2xvcmlnaW50eXBlGNm92uwBIAEoDjIqLmNsb3VkZnJvbnQuT3'
    'JpZ2luQWNjZXNzQ29udHJvbE9yaWdpblR5cGVzUh1vcmlnaW5hY2Nlc3Njb250cm9sb3JpZ2lu'
    'dHlwZRJdCg9zaWduaW5nYmVoYXZpb3IYt/Cv/gEgASgOMi8uY2xvdWRmcm9udC5PcmlnaW5BY2'
    'Nlc3NDb250cm9sU2lnbmluZ0JlaGF2aW9yc1IPc2lnbmluZ2JlaGF2aW9yEl0KD3NpZ25pbmdw'
    'cm90b2NvbBihw+iiASABKA4yLy5jbG91ZGZyb250Lk9yaWdpbkFjY2Vzc0NvbnRyb2xTaWduaW'
    '5nUHJvdG9jb2xzUg9zaWduaW5ncHJvdG9jb2w=');

@$core.Deprecated('Use originCustomHeaderDescriptor instead')
const OriginCustomHeader$json = {
  '1': 'OriginCustomHeader',
  '2': [
    {'1': 'headername', '3': 235244324, '4': 1, '5': 9, '10': 'headername'},
    {'1': 'headervalue', '3': 25328434, '4': 1, '5': 9, '10': 'headervalue'},
  ],
};

/// Descriptor for `OriginCustomHeader`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originCustomHeaderDescriptor = $convert.base64Decode(
    'ChJPcmlnaW5DdXN0b21IZWFkZXISIQoKaGVhZGVybmFtZRiklpZwIAEoCVIKaGVhZGVybmFtZR'
    'IjCgtoZWFkZXJ2YWx1ZRiy9okMIAEoCVILaGVhZGVydmFsdWU=');

@$core.Deprecated('Use originGroupDescriptor instead')
const OriginGroup$json = {
  '1': 'OriginGroup',
  '2': [
    {
      '1': 'failovercriteria',
      '3': 474175659,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginGroupFailoverCriteria',
      '10': 'failovercriteria'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'members',
      '3': 222074545,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginGroupMembers',
      '10': 'members'
    },
    {
      '1': 'selectioncriteria',
      '3': 3089005,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginGroupSelectionCriteria',
      '10': 'selectioncriteria'
    },
  ],
};

/// Descriptor for `OriginGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originGroupDescriptor = $convert.base64Decode(
    'CgtPcmlnaW5Hcm91cBJXChBmYWlsb3ZlcmNyaXRlcmlhGKuxjeIBIAEoCzInLmNsb3VkZnJvbn'
    'QuT3JpZ2luR3JvdXBGYWlsb3ZlckNyaXRlcmlhUhBmYWlsb3ZlcmNyaXRlcmlhEhIKAmlkGIHy'
    'orcBIAEoCVICaWQSOwoHbWVtYmVycxixrfJpIAEoCzIeLmNsb3VkZnJvbnQuT3JpZ2luR3JvdX'
    'BNZW1iZXJzUgdtZW1iZXJzElkKEXNlbGVjdGlvbmNyaXRlcmlhGO3EvAEgASgOMiguY2xvdWRm'
    'cm9udC5PcmlnaW5Hcm91cFNlbGVjdGlvbkNyaXRlcmlhUhFzZWxlY3Rpb25jcml0ZXJpYQ==');

@$core.Deprecated('Use originGroupFailoverCriteriaDescriptor instead')
const OriginGroupFailoverCriteria$json = {
  '1': 'OriginGroupFailoverCriteria',
  '2': [
    {
      '1': 'statuscodes',
      '3': 255295480,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StatusCodes',
      '10': 'statuscodes'
    },
  ],
};

/// Descriptor for `OriginGroupFailoverCriteria`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originGroupFailoverCriteriaDescriptor =
    $convert.base64Decode(
        'ChtPcmlnaW5Hcm91cEZhaWxvdmVyQ3JpdGVyaWESPAoLc3RhdHVzY29kZXMY+P/deSABKAsyFy'
        '5jbG91ZGZyb250LlN0YXR1c0NvZGVzUgtzdGF0dXNjb2Rlcw==');

@$core.Deprecated('Use originGroupMemberDescriptor instead')
const OriginGroupMember$json = {
  '1': 'OriginGroupMember',
  '2': [
    {'1': 'originid', '3': 456733571, '4': 1, '5': 9, '10': 'originid'},
  ],
};

/// Descriptor for `OriginGroupMember`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originGroupMemberDescriptor = $convert.base64Decode(
    'ChFPcmlnaW5Hcm91cE1lbWJlchIeCghvcmlnaW5pZBiD5+TZASABKAlSCG9yaWdpbmlk');

@$core.Deprecated('Use originGroupMembersDescriptor instead')
const OriginGroupMembers$json = {
  '1': 'OriginGroupMembers',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.OriginGroupMember',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `OriginGroupMembers`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originGroupMembersDescriptor = $convert.base64Decode(
    'ChJPcmlnaW5Hcm91cE1lbWJlcnMSNgoFaXRlbXMYsPDYASADKAsyHS5jbG91ZGZyb250Lk9yaW'
    'dpbkdyb3VwTWVtYmVyUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use originGroupsDescriptor instead')
const OriginGroups$json = {
  '1': 'OriginGroups',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.OriginGroup',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `OriginGroups`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originGroupsDescriptor = $convert.base64Decode(
    'CgxPcmlnaW5Hcm91cHMSMAoFaXRlbXMYsPDYASADKAsyFy5jbG91ZGZyb250Lk9yaWdpbkdyb3'
    'VwUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use originMtlsConfigDescriptor instead')
const OriginMtlsConfig$json = {
  '1': 'OriginMtlsConfig',
  '2': [
    {
      '1': 'clientcertificatearn',
      '3': 78787841,
      '4': 1,
      '5': 9,
      '10': 'clientcertificatearn'
    },
  ],
};

/// Descriptor for `OriginMtlsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originMtlsConfigDescriptor = $convert.base64Decode(
    'ChBPcmlnaW5NdGxzQ29uZmlnEjUKFGNsaWVudGNlcnRpZmljYXRlYXJuGIHqyCUgASgJUhRjbG'
    'llbnRjZXJ0aWZpY2F0ZWFybg==');

@$core.Deprecated('Use originRequestPolicyDescriptor instead')
const OriginRequestPolicy$json = {
  '1': 'OriginRequestPolicy',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'originrequestpolicyconfig',
      '3': 37078133,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyConfig',
      '10': 'originrequestpolicyconfig'
    },
  ],
};

/// Descriptor for `OriginRequestPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyDescriptor = $convert.base64Decode(
    'ChNPcmlnaW5SZXF1ZXN0UG9saWN5EhIKAmlkGIHyorcBIAEoCVICaWQSLQoQbGFzdG1vZGlmaW'
    'VkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZRJmChlvcmlnaW5yZXF1ZXN0cG9saWN5'
    'Y29uZmlnGPWI1xEgASgLMiUuY2xvdWRmcm9udC5PcmlnaW5SZXF1ZXN0UG9saWN5Q29uZmlnUh'
    'lvcmlnaW5yZXF1ZXN0cG9saWN5Y29uZmln');

@$core.Deprecated('Use originRequestPolicyAlreadyExistsDescriptor instead')
const OriginRequestPolicyAlreadyExists$json = {
  '1': 'OriginRequestPolicyAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OriginRequestPolicyAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CiBPcmlnaW5SZXF1ZXN0UG9saWN5QWxyZWFkeUV4aXN0cxIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use originRequestPolicyConfigDescriptor instead')
const OriginRequestPolicyConfig$json = {
  '1': 'OriginRequestPolicyConfig',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'cookiesconfig',
      '3': 243885539,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyCookiesConfig',
      '10': 'cookiesconfig'
    },
    {
      '1': 'headersconfig',
      '3': 368880604,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyHeadersConfig',
      '10': 'headersconfig'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'querystringsconfig',
      '3': 129166430,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyQueryStringsConfig',
      '10': 'querystringsconfig'
    },
  ],
};

/// Descriptor for `OriginRequestPolicyConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyConfigDescriptor = $convert.base64Decode(
    'ChlPcmlnaW5SZXF1ZXN0UG9saWN5Q29uZmlnEhwKB2NvbW1lbnQY/7++wgEgASgJUgdjb21tZW'
    '50ElUKDWNvb2tpZXNjb25maWcY48uldCABKAsyLC5jbG91ZGZyb250Lk9yaWdpblJlcXVlc3RQ'
    'b2xpY3lDb29raWVzQ29uZmlnUg1jb29raWVzY29uZmlnElYKDWhlYWRlcnNjb25maWcY3Nfyrw'
    'EgASgLMiwuY2xvdWRmcm9udC5PcmlnaW5SZXF1ZXN0UG9saWN5SGVhZGVyc0NvbmZpZ1INaGVh'
    'ZGVyc2NvbmZpZxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEmQKEnF1ZXJ5c3RyaW5nc2NvbmZpZx'
    'je2Ms9IAEoCzIxLmNsb3VkZnJvbnQuT3JpZ2luUmVxdWVzdFBvbGljeVF1ZXJ5U3RyaW5nc0Nv'
    'bmZpZ1IScXVlcnlzdHJpbmdzY29uZmln');

@$core.Deprecated('Use originRequestPolicyCookiesConfigDescriptor instead')
const OriginRequestPolicyCookiesConfig$json = {
  '1': 'OriginRequestPolicyCookiesConfig',
  '2': [
    {
      '1': 'cookiebehavior',
      '3': 86223930,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginRequestPolicyCookieBehavior',
      '10': 'cookiebehavior'
    },
    {
      '1': 'cookies',
      '3': 418634949,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CookieNames',
      '10': 'cookies'
    },
  ],
};

/// Descriptor for `OriginRequestPolicyCookiesConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyCookiesConfigDescriptor =
    $convert.base64Decode(
        'CiBPcmlnaW5SZXF1ZXN0UG9saWN5Q29va2llc0NvbmZpZxJYCg5jb29raWViZWhhdmlvchi62I'
        '4pIAEoDjItLmNsb3VkZnJvbnQuT3JpZ2luUmVxdWVzdFBvbGljeUNvb2tpZUJlaGF2aW9yUg5j'
        'b29raWViZWhhdmlvchI1Cgdjb29raWVzGMW5z8cBIAEoCzIXLmNsb3VkZnJvbnQuQ29va2llTm'
        'FtZXNSB2Nvb2tpZXM=');

@$core.Deprecated('Use originRequestPolicyHeadersConfigDescriptor instead')
const OriginRequestPolicyHeadersConfig$json = {
  '1': 'OriginRequestPolicyHeadersConfig',
  '2': [
    {
      '1': 'headerbehavior',
      '3': 348515467,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginRequestPolicyHeaderBehavior',
      '10': 'headerbehavior'
    },
    {
      '1': 'headers',
      '3': 323967370,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Headers',
      '10': 'headers'
    },
  ],
};

/// Descriptor for `OriginRequestPolicyHeadersConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyHeadersConfigDescriptor =
    $convert.base64Decode(
        'CiBPcmlnaW5SZXF1ZXN0UG9saWN5SGVhZGVyc0NvbmZpZxJZCg5oZWFkZXJiZWhhdmlvchiL2Z'
        'emASABKA4yLS5jbG91ZGZyb250Lk9yaWdpblJlcXVlc3RQb2xpY3lIZWFkZXJCZWhhdmlvclIO'
        'aGVhZGVyYmVoYXZpb3ISMQoHaGVhZGVycxiKs72aASABKAsyEy5jbG91ZGZyb250LkhlYWRlcn'
        'NSB2hlYWRlcnM=');

@$core.Deprecated('Use originRequestPolicyInUseDescriptor instead')
const OriginRequestPolicyInUse$json = {
  '1': 'OriginRequestPolicyInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OriginRequestPolicyInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyInUseDescriptor =
    $convert.base64Decode(
        'ChhPcmlnaW5SZXF1ZXN0UG9saWN5SW5Vc2USGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use originRequestPolicyListDescriptor instead')
const OriginRequestPolicyList$json = {
  '1': 'OriginRequestPolicyList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicySummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `OriginRequestPolicyList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyListDescriptor = $convert.base64Decode(
    'ChdPcmlnaW5SZXF1ZXN0UG9saWN5TGlzdBI/CgVpdGVtcxiw8NgBIAMoCzImLmNsb3VkZnJvbn'
    'QuT3JpZ2luUmVxdWVzdFBvbGljeVN1bW1hcnlSBWl0ZW1zEh4KCG1heGl0ZW1zGJTW2vEBIAEo'
    'BVIIbWF4aXRlbXMSIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leHRtYXJrZXISHQoIcXVhbn'
    'RpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use originRequestPolicyQueryStringsConfigDescriptor instead')
const OriginRequestPolicyQueryStringsConfig$json = {
  '1': 'OriginRequestPolicyQueryStringsConfig',
  '2': [
    {
      '1': 'querystringbehavior',
      '3': 235804641,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginRequestPolicyQueryStringBehavior',
      '10': 'querystringbehavior'
    },
    {
      '1': 'querystrings',
      '3': 478781456,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.QueryStringNames',
      '10': 'querystrings'
    },
  ],
};

/// Descriptor for `OriginRequestPolicyQueryStringsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicyQueryStringsConfigDescriptor =
    $convert.base64Decode(
        'CiVPcmlnaW5SZXF1ZXN0UG9saWN5UXVlcnlTdHJpbmdzQ29uZmlnEmcKE3F1ZXJ5c3RyaW5nYm'
        'VoYXZpb3IY4a+4cCABKA4yMi5jbG91ZGZyb250Lk9yaWdpblJlcXVlc3RQb2xpY3lRdWVyeVN0'
        'cmluZ0JlaGF2aW9yUhNxdWVyeXN0cmluZ2JlaGF2aW9yEkQKDHF1ZXJ5c3RyaW5ncxiQwKbkAS'
        'ABKAsyHC5jbG91ZGZyb250LlF1ZXJ5U3RyaW5nTmFtZXNSDHF1ZXJ5c3RyaW5ncw==');

@$core.Deprecated('Use originRequestPolicySummaryDescriptor instead')
const OriginRequestPolicySummary$json = {
  '1': 'OriginRequestPolicySummary',
  '2': [
    {
      '1': 'originrequestpolicy',
      '3': 386733531,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicy',
      '10': 'originrequestpolicy'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginRequestPolicyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `OriginRequestPolicySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originRequestPolicySummaryDescriptor = $convert.base64Decode(
    'ChpPcmlnaW5SZXF1ZXN0UG9saWN5U3VtbWFyeRJVChNvcmlnaW5yZXF1ZXN0cG9saWN5GNurtL'
    'gBIAEoCzIfLmNsb3VkZnJvbnQuT3JpZ2luUmVxdWVzdFBvbGljeVITb3JpZ2lucmVxdWVzdHBv'
    'bGljeRI7CgR0eXBlGO6g14oBIAEoDjIjLmNsb3VkZnJvbnQuT3JpZ2luUmVxdWVzdFBvbGljeV'
    'R5cGVSBHR5cGU=');

@$core.Deprecated('Use originShieldDescriptor instead')
const OriginShield$json = {
  '1': 'OriginShield',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'originshieldregion',
      '3': 448983227,
      '4': 1,
      '5': 9,
      '10': 'originshieldregion'
    },
  ],
};

/// Descriptor for `OriginShield`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originShieldDescriptor = $convert.base64Decode(
    'CgxPcmlnaW5TaGllbGQSHAoHZW5hYmxlZBi/yJvkASABKAhSB2VuYWJsZWQSMgoSb3JpZ2luc2'
    'hpZWxkcmVnaW9uGLvhi9YBIAEoCVISb3JpZ2luc2hpZWxkcmVnaW9u');

@$core.Deprecated('Use originSslProtocolsDescriptor instead')
const OriginSslProtocols$json = {
  '1': 'OriginSslProtocols',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 14,
      '6': '.cloudfront.SslProtocol',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `OriginSslProtocols`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originSslProtocolsDescriptor = $convert.base64Decode(
    'ChJPcmlnaW5Tc2xQcm90b2NvbHMSMAoFaXRlbXMYsPDYASADKA4yFy5jbG91ZGZyb250LlNzbF'
    'Byb3RvY29sUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use originsDescriptor instead')
const Origins$json = {
  '1': 'Origins',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.Origin',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `Origins`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List originsDescriptor = $convert.base64Decode(
    'CgdPcmlnaW5zEisKBWl0ZW1zGLDw2AEgAygLMhIuY2xvdWRmcm9udC5PcmlnaW5SBWl0ZW1zEh'
    '0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use parameterDescriptor instead')
const Parameter$json = {
  '1': 'Parameter',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Parameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterDescriptor = $convert.base64Decode(
    'CglQYXJhbWV0ZXISFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIYCgV2YWx1ZRjr8p+KASABKAlSBX'
    'ZhbHVl');

@$core.Deprecated('Use parameterDefinitionDescriptor instead')
const ParameterDefinition$json = {
  '1': 'ParameterDefinition',
  '2': [
    {
      '1': 'definition',
      '3': 222998209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ParameterDefinitionSchema',
      '10': 'definition'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `ParameterDefinition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterDefinitionDescriptor = $convert.base64Decode(
    'ChNQYXJhbWV0ZXJEZWZpbml0aW9uEkgKCmRlZmluaXRpb24Ywd2qaiABKAsyJS5jbG91ZGZyb2'
    '50LlBhcmFtZXRlckRlZmluaXRpb25TY2hlbWFSCmRlZmluaXRpb24SFQoEbmFtZRiH5oF/IAEo'
    'CVIEbmFtZQ==');

@$core.Deprecated('Use parameterDefinitionSchemaDescriptor instead')
const ParameterDefinitionSchema$json = {
  '1': 'ParameterDefinitionSchema',
  '2': [
    {
      '1': 'stringschema',
      '3': 444400660,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StringSchemaConfig',
      '10': 'stringschema'
    },
  ],
};

/// Descriptor for `ParameterDefinitionSchema`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterDefinitionSchemaDescriptor =
    $convert.base64Decode(
        'ChlQYXJhbWV0ZXJEZWZpbml0aW9uU2NoZW1hEkYKDHN0cmluZ3NjaGVtYRiUiPTTASABKAsyHi'
        '5jbG91ZGZyb250LlN0cmluZ1NjaGVtYUNvbmZpZ1IMc3RyaW5nc2NoZW1h');

@$core.Deprecated(
    'Use parametersInCacheKeyAndForwardedToOriginDescriptor instead')
const ParametersInCacheKeyAndForwardedToOrigin$json = {
  '1': 'ParametersInCacheKeyAndForwardedToOrigin',
  '2': [
    {
      '1': 'cookiesconfig',
      '3': 243885539,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyCookiesConfig',
      '10': 'cookiesconfig'
    },
    {
      '1': 'enableacceptencodingbrotli',
      '3': 105219116,
      '4': 1,
      '5': 8,
      '10': 'enableacceptencodingbrotli'
    },
    {
      '1': 'enableacceptencodinggzip',
      '3': 150977880,
      '4': 1,
      '5': 8,
      '10': 'enableacceptencodinggzip'
    },
    {
      '1': 'headersconfig',
      '3': 368880604,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyHeadersConfig',
      '10': 'headersconfig'
    },
    {
      '1': 'querystringsconfig',
      '3': 129166430,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyQueryStringsConfig',
      '10': 'querystringsconfig'
    },
  ],
};

/// Descriptor for `ParametersInCacheKeyAndForwardedToOrigin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parametersInCacheKeyAndForwardedToOriginDescriptor = $convert.base64Decode(
    'CihQYXJhbWV0ZXJzSW5DYWNoZUtleUFuZEZvcndhcmRlZFRvT3JpZ2luEk0KDWNvb2tpZXNjb2'
    '5maWcY48uldCABKAsyJC5jbG91ZGZyb250LkNhY2hlUG9saWN5Q29va2llc0NvbmZpZ1INY29v'
    'a2llc2NvbmZpZxJBChplbmFibGVhY2NlcHRlbmNvZGluZ2Jyb3RsaRisiJYyIAEoCFIaZW5hYm'
    'xlYWNjZXB0ZW5jb2Rpbmdicm90bGkSPQoYZW5hYmxlYWNjZXB0ZW5jb2RpbmdnemlwGNj6/kcg'
    'ASgIUhhlbmFibGVhY2NlcHRlbmNvZGluZ2d6aXASTgoNaGVhZGVyc2NvbmZpZxjc1/KvASABKA'
    'syJC5jbG91ZGZyb250LkNhY2hlUG9saWN5SGVhZGVyc0NvbmZpZ1INaGVhZGVyc2NvbmZpZxJc'
    'ChJxdWVyeXN0cmluZ3Njb25maWcY3tjLPSABKAsyKS5jbG91ZGZyb250LkNhY2hlUG9saWN5UX'
    'VlcnlTdHJpbmdzQ29uZmlnUhJxdWVyeXN0cmluZ3Njb25maWc=');

@$core.Deprecated('Use pathsDescriptor instead')
const Paths$json = {
  '1': 'Paths',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `Paths`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pathsDescriptor = $convert.base64Decode(
    'CgVQYXRocxIXCgVpdGVtcxiw8NgBIAMoCVIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXyABKAVSCH'
    'F1YW50aXR5');

@$core.Deprecated('Use preconditionFailedDescriptor instead')
const PreconditionFailed$json = {
  '1': 'PreconditionFailed',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PreconditionFailed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List preconditionFailedDescriptor =
    $convert.base64Decode(
        'ChJQcmVjb25kaXRpb25GYWlsZWQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use publicKeyDescriptor instead')
const PublicKey$json = {
  '1': 'PublicKey',
  '2': [
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'publickeyconfig',
      '3': 228537966,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKeyConfig',
      '10': 'publickeyconfig'
    },
  ],
};

/// Descriptor for `PublicKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyDescriptor = $convert.base64Decode(
    'CglQdWJsaWNLZXkSIwoLY3JlYXRlZHRpbWUY8+vzOSABKAlSC2NyZWF0ZWR0aW1lEhIKAmlkGI'
    'HyorcBIAEoCVICaWQSSAoPcHVibGlja2V5Y29uZmlnGO7s/GwgASgLMhsuY2xvdWRmcm9udC5Q'
    'dWJsaWNLZXlDb25maWdSD3B1YmxpY2tleWNvbmZpZw==');

@$core.Deprecated('Use publicKeyAlreadyExistsDescriptor instead')
const PublicKeyAlreadyExists$json = {
  '1': 'PublicKeyAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PublicKeyAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChZQdWJsaWNLZXlBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use publicKeyConfigDescriptor instead')
const PublicKeyConfig$json = {
  '1': 'PublicKeyConfig',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'encodedkey', '3': 13297, '4': 1, '5': 9, '10': 'encodedkey'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `PublicKeyConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyConfigDescriptor = $convert.base64Decode(
    'Cg9QdWJsaWNLZXlDb25maWcSKwoPY2FsbGVycmVmZXJlbmNlGJiZjUggASgJUg9jYWxsZXJyZW'
    'ZlcmVuY2USHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1lbnQSHwoKZW5jb2RlZGtleRjxZyAB'
    'KAlSCmVuY29kZWRrZXkSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use publicKeyInUseDescriptor instead')
const PublicKeyInUse$json = {
  '1': 'PublicKeyInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PublicKeyInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyInUseDescriptor = $convert.base64Decode(
    'Cg5QdWJsaWNLZXlJblVzZRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use publicKeyListDescriptor instead')
const PublicKeyList$json = {
  '1': 'PublicKeyList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.PublicKeySummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `PublicKeyList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyListDescriptor = $convert.base64Decode(
    'Cg1QdWJsaWNLZXlMaXN0EjUKBWl0ZW1zGLDw2AEgAygLMhwuY2xvdWRmcm9udC5QdWJsaWNLZX'
    'lTdW1tYXJ5UgVpdGVtcxIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1zEiIKCm5leHRt'
    'YXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2VyEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudG'
    'l0eQ==');

@$core.Deprecated('Use publicKeySummaryDescriptor instead')
const PublicKeySummary$json = {
  '1': 'PublicKeySummary',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'encodedkey', '3': 13297, '4': 1, '5': 9, '10': 'encodedkey'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `PublicKeySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeySummaryDescriptor = $convert.base64Decode(
    'ChBQdWJsaWNLZXlTdW1tYXJ5EhwKB2NvbW1lbnQY/7++wgEgASgJUgdjb21tZW50EiMKC2NyZW'
    'F0ZWR0aW1lGPPr8zkgASgJUgtjcmVhdGVkdGltZRIfCgplbmNvZGVka2V5GPFnIAEoCVIKZW5j'
    'b2RlZGtleRISCgJpZBiB8qK3ASABKAlSAmlkEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use publishConnectionFunctionRequestDescriptor instead')
const PublishConnectionFunctionRequest$json = {
  '1': 'PublishConnectionFunctionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `PublishConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishConnectionFunctionRequestDescriptor =
    $convert.base64Decode(
        'CiBQdWJsaXNoQ29ubmVjdGlvbkZ1bmN0aW9uUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEh'
        'sKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2g=');

@$core.Deprecated('Use publishConnectionFunctionResultDescriptor instead')
const PublishConnectionFunctionResult$json = {
  '1': 'PublishConnectionFunctionResult',
  '2': [
    {
      '1': 'connectionfunctionsummary',
      '3': 62528396,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionSummary',
      '8': {},
      '10': 'connectionfunctionsummary'
    },
  ],
};

/// Descriptor for `PublishConnectionFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishConnectionFunctionResultDescriptor =
    $convert.base64Decode(
        'Ch9QdWJsaXNoQ29ubmVjdGlvbkZ1bmN0aW9uUmVzdWx0EmwKGWNvbm5lY3Rpb25mdW5jdGlvbn'
        'N1bW1hcnkYjLfoHSABKAsyJS5jbG91ZGZyb250LkNvbm5lY3Rpb25GdW5jdGlvblN1bW1hcnlC'
        'BIi1GAFSGWNvbm5lY3Rpb25mdW5jdGlvbnN1bW1hcnk=');

@$core.Deprecated('Use publishFunctionRequestDescriptor instead')
const PublishFunctionRequest$json = {
  '1': 'PublishFunctionRequest',
  '2': [
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `PublishFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishFunctionRequestDescriptor =
    $convert.base64Decode(
        'ChZQdWJsaXNoRnVuY3Rpb25SZXF1ZXN0EhsKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2gSFQ'
        'oEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use publishFunctionResultDescriptor instead')
const PublishFunctionResult$json = {
  '1': 'PublishFunctionResult',
  '2': [
    {
      '1': 'functionsummary',
      '3': 523316264,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionSummary',
      '8': {},
      '10': 'functionsummary'
    },
  ],
};

/// Descriptor for `PublishFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishFunctionResultDescriptor = $convert.base64Decode(
    'ChVQdWJsaXNoRnVuY3Rpb25SZXN1bHQSTwoPZnVuY3Rpb25zdW1tYXJ5GKjYxPkBIAEoCzIbLm'
    'Nsb3VkZnJvbnQuRnVuY3Rpb25TdW1tYXJ5QgSItRgBUg9mdW5jdGlvbnN1bW1hcnk=');

@$core.Deprecated('Use putResourcePolicyRequestDescriptor instead')
const PutResourcePolicyRequest$json = {
  '1': 'PutResourcePolicyRequest',
  '2': [
    {
      '1': 'policydocument',
      '3': 238049099,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `PutResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhQdXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSKQoOcG9saWN5ZG9jdW1lbnQYy67BcSABKAlSDn'
        'BvbGljeWRvY3VtZW50EiQKC3Jlc291cmNlYXJuGK342a0BIAEoCVILcmVzb3VyY2Vhcm4=');

@$core.Deprecated('Use putResourcePolicyResultDescriptor instead')
const PutResourcePolicyResult$json = {
  '1': 'PutResourcePolicyResult',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `PutResourcePolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyResultDescriptor =
    $convert.base64Decode(
        'ChdQdXRSZXNvdXJjZVBvbGljeVJlc3VsdBIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc2'
        '91cmNlYXJu');

@$core.Deprecated('Use queryArgProfileDescriptor instead')
const QueryArgProfile$json = {
  '1': 'QueryArgProfile',
  '2': [
    {'1': 'profileid', '3': 407138548, '4': 1, '5': 9, '10': 'profileid'},
    {'1': 'queryarg', '3': 99859842, '4': 1, '5': 9, '10': 'queryarg'},
  ],
};

/// Descriptor for `QueryArgProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryArgProfileDescriptor = $convert.base64Decode(
    'Cg9RdWVyeUFyZ1Byb2ZpbGUSIAoJcHJvZmlsZWlkGPThkcIBIAEoCVIJcHJvZmlsZWlkEh0KCH'
    'F1ZXJ5YXJnGIL7zi8gASgJUghxdWVyeWFyZw==');

@$core.Deprecated('Use queryArgProfileConfigDescriptor instead')
const QueryArgProfileConfig$json = {
  '1': 'QueryArgProfileConfig',
  '2': [
    {
      '1': 'forwardwhenqueryargprofileisunknown',
      '3': 43564528,
      '4': 1,
      '5': 8,
      '10': 'forwardwhenqueryargprofileisunknown'
    },
    {
      '1': 'queryargprofiles',
      '3': 155212116,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.QueryArgProfiles',
      '10': 'queryargprofiles'
    },
  ],
};

/// Descriptor for `QueryArgProfileConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryArgProfileConfigDescriptor = $convert.base64Decode(
    'ChVRdWVyeUFyZ1Byb2ZpbGVDb25maWcSUwojZm9yd2FyZHdoZW5xdWVyeWFyZ3Byb2ZpbGVpc3'
    'Vua25vd24Y8PviFCABKAhSI2ZvcndhcmR3aGVucXVlcnlhcmdwcm9maWxlaXN1bmtub3duEksK'
    'EHF1ZXJ5YXJncHJvZmlsZXMY1LKBSiABKAsyHC5jbG91ZGZyb250LlF1ZXJ5QXJnUHJvZmlsZX'
    'NSEHF1ZXJ5YXJncHJvZmlsZXM=');

@$core.Deprecated('Use queryArgProfileEmptyDescriptor instead')
const QueryArgProfileEmpty$json = {
  '1': 'QueryArgProfileEmpty',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueryArgProfileEmpty`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryArgProfileEmptyDescriptor =
    $convert.base64Decode(
        'ChRRdWVyeUFyZ1Byb2ZpbGVFbXB0eRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use queryArgProfilesDescriptor instead')
const QueryArgProfiles$json = {
  '1': 'QueryArgProfiles',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.QueryArgProfile',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `QueryArgProfiles`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryArgProfilesDescriptor = $convert.base64Decode(
    'ChBRdWVyeUFyZ1Byb2ZpbGVzEjQKBWl0ZW1zGLDw2AEgAygLMhsuY2xvdWRmcm9udC5RdWVyeU'
    'FyZ1Byb2ZpbGVSBWl0ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use queryStringCacheKeysDescriptor instead')
const QueryStringCacheKeys$json = {
  '1': 'QueryStringCacheKeys',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `QueryStringCacheKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStringCacheKeysDescriptor = $convert.base64Decode(
    'ChRRdWVyeVN0cmluZ0NhY2hlS2V5cxIXCgVpdGVtcxiw8NgBIAMoCVIFaXRlbXMSHQoIcXVhbn'
    'RpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use queryStringNamesDescriptor instead')
const QueryStringNames$json = {
  '1': 'QueryStringNames',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `QueryStringNames`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStringNamesDescriptor = $convert.base64Decode(
    'ChBRdWVyeVN0cmluZ05hbWVzEhcKBWl0ZW1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eR'
    'j55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use realtimeLogConfigDescriptor instead')
const RealtimeLogConfig$json = {
  '1': 'RealtimeLogConfig',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'endpoints',
      '3': 436023390,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.EndPoint',
      '10': 'endpoints'
    },
    {'1': 'fields', '3': 319339933, '4': 3, '5': 9, '10': 'fields'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'samplingrate', '3': 272929747, '4': 1, '5': 3, '10': 'samplingrate'},
  ],
};

/// Descriptor for `RealtimeLogConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List realtimeLogConfigDescriptor = $convert.base64Decode(
    'ChFSZWFsdGltZUxvZ0NvbmZpZxIUCgNhcm4YnZyvvQEgASgJUgNhcm4SNgoJZW5kcG9pbnRzGN'
    '7g9M8BIAMoCzIULmNsb3VkZnJvbnQuRW5kUG9pbnRSCWVuZHBvaW50cxIaCgZmaWVsZHMYnfui'
    'mAEgAygJUgZmaWVsZHMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRImCgxzYW1wbGluZ3JhdGUY06'
    'eSggEgASgDUgxzYW1wbGluZ3JhdGU=');

@$core.Deprecated('Use realtimeLogConfigAlreadyExistsDescriptor instead')
const RealtimeLogConfigAlreadyExists$json = {
  '1': 'RealtimeLogConfigAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RealtimeLogConfigAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List realtimeLogConfigAlreadyExistsDescriptor =
    $convert.base64Decode(
        'Ch5SZWFsdGltZUxvZ0NvbmZpZ0FscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use realtimeLogConfigInUseDescriptor instead')
const RealtimeLogConfigInUse$json = {
  '1': 'RealtimeLogConfigInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RealtimeLogConfigInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List realtimeLogConfigInUseDescriptor =
    $convert.base64Decode(
        'ChZSZWFsdGltZUxvZ0NvbmZpZ0luVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use realtimeLogConfigOwnerMismatchDescriptor instead')
const RealtimeLogConfigOwnerMismatch$json = {
  '1': 'RealtimeLogConfigOwnerMismatch',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RealtimeLogConfigOwnerMismatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List realtimeLogConfigOwnerMismatchDescriptor =
    $convert.base64Decode(
        'Ch5SZWFsdGltZUxvZ0NvbmZpZ093bmVyTWlzbWF0Y2gSGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use realtimeLogConfigsDescriptor instead')
const RealtimeLogConfigs$json = {
  '1': 'RealtimeLogConfigs',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.RealtimeLogConfig',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `RealtimeLogConfigs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List realtimeLogConfigsDescriptor = $convert.base64Decode(
    'ChJSZWFsdGltZUxvZ0NvbmZpZ3MSIwoLaXN0cnVuY2F0ZWQY2p+4cyABKAhSC2lzdHJ1bmNhdG'
    'VkEjYKBWl0ZW1zGLDw2AEgAygLMh0uY2xvdWRmcm9udC5SZWFsdGltZUxvZ0NvbmZpZ1IFaXRl'
    'bXMSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYX'
    'hpdGVtcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use realtimeMetricsSubscriptionConfigDescriptor instead')
const RealtimeMetricsSubscriptionConfig$json = {
  '1': 'RealtimeMetricsSubscriptionConfig',
  '2': [
    {
      '1': 'realtimemetricssubscriptionstatus',
      '3': 158225507,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.RealtimeMetricsSubscriptionStatus',
      '10': 'realtimemetricssubscriptionstatus'
    },
  ],
};

/// Descriptor for `RealtimeMetricsSubscriptionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List realtimeMetricsSubscriptionConfigDescriptor =
    $convert.base64Decode(
        'CiFSZWFsdGltZU1ldHJpY3NTdWJzY3JpcHRpb25Db25maWcSfgohcmVhbHRpbWVtZXRyaWNzc3'
        'Vic2NyaXB0aW9uc3RhdHVzGOOouUsgASgOMi0uY2xvdWRmcm9udC5SZWFsdGltZU1ldHJpY3NT'
        'dWJzY3JpcHRpb25TdGF0dXNSIXJlYWx0aW1lbWV0cmljc3N1YnNjcmlwdGlvbnN0YXR1cw==');

@$core.Deprecated('Use resourceInUseDescriptor instead')
const ResourceInUse$json = {
  '1': 'ResourceInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceInUseDescriptor = $convert.base64Decode(
    'Cg1SZXNvdXJjZUluVXNlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use resourceNotDisabledDescriptor instead')
const ResourceNotDisabled$json = {
  '1': 'ResourceNotDisabled',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotDisabled`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotDisabledDescriptor =
    $convert.base64Decode(
        'ChNSZXNvdXJjZU5vdERpc2FibGVkEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use responseHeadersPolicyDescriptor instead')
const ResponseHeadersPolicy$json = {
  '1': 'ResponseHeadersPolicy',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'responseheaderspolicyconfig',
      '3': 159056825,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyConfig',
      '10': 'responseheaderspolicyconfig'
    },
  ],
};

/// Descriptor for `ResponseHeadersPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyDescriptor = $convert.base64Decode(
    'ChVSZXNwb25zZUhlYWRlcnNQb2xpY3kSEgoCaWQYgfKitwEgASgJUgJpZBItChBsYXN0bW9kaW'
    'ZpZWR0aW1lGOCC/HAgASgJUhBsYXN0bW9kaWZpZWR0aW1lEmwKG3Jlc3BvbnNlaGVhZGVyc3Bv'
    'bGljeWNvbmZpZxi5h+xLIAEoCzInLmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9saWN5Q2'
    '9uZmlnUhtyZXNwb25zZWhlYWRlcnNwb2xpY3ljb25maWc=');

@$core.Deprecated(
    'Use responseHeadersPolicyAccessControlAllowHeadersDescriptor instead')
const ResponseHeadersPolicyAccessControlAllowHeaders$json = {
  '1': 'ResponseHeadersPolicyAccessControlAllowHeaders',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyAccessControlAllowHeaders`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    responseHeadersPolicyAccessControlAllowHeadersDescriptor =
    $convert.base64Decode(
        'Ci5SZXNwb25zZUhlYWRlcnNQb2xpY3lBY2Nlc3NDb250cm9sQWxsb3dIZWFkZXJzEhcKBWl0ZW'
        '1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated(
    'Use responseHeadersPolicyAccessControlAllowMethodsDescriptor instead')
const ResponseHeadersPolicyAccessControlAllowMethods$json = {
  '1': 'ResponseHeadersPolicyAccessControlAllowMethods',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 14,
      '6': '.cloudfront.ResponseHeadersPolicyAccessControlAllowMethodsValues',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyAccessControlAllowMethods`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    responseHeadersPolicyAccessControlAllowMethodsDescriptor =
    $convert.base64Decode(
        'Ci5SZXNwb25zZUhlYWRlcnNQb2xpY3lBY2Nlc3NDb250cm9sQWxsb3dNZXRob2RzElkKBWl0ZW'
        '1zGLDw2AEgAygOMkAuY2xvdWRmcm9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lBY2Nlc3NDb250'
        'cm9sQWxsb3dNZXRob2RzVmFsdWVzUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbn'
        'RpdHk=');

@$core.Deprecated(
    'Use responseHeadersPolicyAccessControlAllowOriginsDescriptor instead')
const ResponseHeadersPolicyAccessControlAllowOrigins$json = {
  '1': 'ResponseHeadersPolicyAccessControlAllowOrigins',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyAccessControlAllowOrigins`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    responseHeadersPolicyAccessControlAllowOriginsDescriptor =
    $convert.base64Decode(
        'Ci5SZXNwb25zZUhlYWRlcnNQb2xpY3lBY2Nlc3NDb250cm9sQWxsb3dPcmlnaW5zEhcKBWl0ZW'
        '1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated(
    'Use responseHeadersPolicyAccessControlExposeHeadersDescriptor instead')
const ResponseHeadersPolicyAccessControlExposeHeaders$json = {
  '1': 'ResponseHeadersPolicyAccessControlExposeHeaders',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyAccessControlExposeHeaders`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    responseHeadersPolicyAccessControlExposeHeadersDescriptor =
    $convert.base64Decode(
        'Ci9SZXNwb25zZUhlYWRlcnNQb2xpY3lBY2Nlc3NDb250cm9sRXhwb3NlSGVhZGVycxIXCgVpdG'
        'Vtcxiw8NgBIAMoCVIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use responseHeadersPolicyAlreadyExistsDescriptor instead')
const ResponseHeadersPolicyAlreadyExists$json = {
  '1': 'ResponseHeadersPolicyAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CiJSZXNwb25zZUhlYWRlcnNQb2xpY3lBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UYhbO7cCABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use responseHeadersPolicyConfigDescriptor instead')
const ResponseHeadersPolicyConfig$json = {
  '1': 'ResponseHeadersPolicyConfig',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'corsconfig',
      '3': 238574547,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyCorsConfig',
      '10': 'corsconfig'
    },
    {
      '1': 'customheadersconfig',
      '3': 121935235,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyCustomHeadersConfig',
      '10': 'customheadersconfig'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'removeheadersconfig',
      '3': 61840100,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyRemoveHeadersConfig',
      '10': 'removeheadersconfig'
    },
    {
      '1': 'securityheadersconfig',
      '3': 412944602,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicySecurityHeadersConfig',
      '10': 'securityheadersconfig'
    },
    {
      '1': 'servertimingheadersconfig',
      '3': 74432993,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyServerTimingHeadersConfig',
      '10': 'servertimingheadersconfig'
    },
  ],
};

/// Descriptor for `ResponseHeadersPolicyConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyConfigDescriptor = $convert.base64Decode(
    'ChtSZXNwb25zZUhlYWRlcnNQb2xpY3lDb25maWcSHAoHY29tbWVudBj/v77CASABKAlSB2NvbW'
    '1lbnQSTgoKY29yc2NvbmZpZxjTt+FxIAEoCzIrLmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJz'
    'UG9saWN5Q29yc0NvbmZpZ1IKY29yc2NvbmZpZxJpChNjdXN0b21oZWFkZXJzY29uZmlnGIOrkj'
    'ogASgLMjQuY2xvdWRmcm9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lDdXN0b21IZWFkZXJzQ29u'
    'ZmlnUhNjdXN0b21oZWFkZXJzY29uZmlnEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSaQoTcmVtb3'
    'ZlaGVhZGVyc2NvbmZpZxjktb4dIAEoCzI0LmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9s'
    'aWN5UmVtb3ZlSGVhZGVyc0NvbmZpZ1ITcmVtb3ZlaGVhZGVyc2NvbmZpZxJwChVzZWN1cml0eW'
    'hlYWRlcnNjb25maWcY2pH0xAEgASgLMjYuY2xvdWRmcm9udC5SZXNwb25zZUhlYWRlcnNQb2xp'
    'Y3lTZWN1cml0eUhlYWRlcnNDb25maWdSFXNlY3VyaXR5aGVhZGVyc2NvbmZpZxJ7ChlzZXJ2ZX'
    'J0aW1pbmdoZWFkZXJzY29uZmlnGOGDvyMgASgLMjouY2xvdWRmcm9udC5SZXNwb25zZUhlYWRl'
    'cnNQb2xpY3lTZXJ2ZXJUaW1pbmdIZWFkZXJzQ29uZmlnUhlzZXJ2ZXJ0aW1pbmdoZWFkZXJzY2'
    '9uZmln');

@$core.Deprecated(
    'Use responseHeadersPolicyContentSecurityPolicyDescriptor instead')
const ResponseHeadersPolicyContentSecurityPolicy$json = {
  '1': 'ResponseHeadersPolicyContentSecurityPolicy',
  '2': [
    {
      '1': 'contentsecuritypolicy',
      '3': 531505511,
      '4': 1,
      '5': 9,
      '10': 'contentsecuritypolicy'
    },
    {'1': 'override', '3': 139255902, '4': 1, '5': 8, '10': 'override'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyContentSecurityPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    responseHeadersPolicyContentSecurityPolicyDescriptor =
    $convert.base64Decode(
        'CipSZXNwb25zZUhlYWRlcnNQb2xpY3lDb250ZW50U2VjdXJpdHlQb2xpY3kSOAoVY29udGVudH'
        'NlY3VyaXR5cG9saWN5GOfCuP0BIAEoCVIVY29udGVudHNlY3VyaXR5cG9saWN5Eh0KCG92ZXJy'
        'aWRlGN7As0IgASgIUghvdmVycmlkZQ==');

@$core
    .Deprecated('Use responseHeadersPolicyContentTypeOptionsDescriptor instead')
const ResponseHeadersPolicyContentTypeOptions$json = {
  '1': 'ResponseHeadersPolicyContentTypeOptions',
  '2': [
    {'1': 'override', '3': 139255902, '4': 1, '5': 8, '10': 'override'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyContentTypeOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyContentTypeOptionsDescriptor =
    $convert.base64Decode(
        'CidSZXNwb25zZUhlYWRlcnNQb2xpY3lDb250ZW50VHlwZU9wdGlvbnMSHQoIb3ZlcnJpZGUY3s'
        'CzQiABKAhSCG92ZXJyaWRl');

@$core.Deprecated('Use responseHeadersPolicyCorsConfigDescriptor instead')
const ResponseHeadersPolicyCorsConfig$json = {
  '1': 'ResponseHeadersPolicyCorsConfig',
  '2': [
    {
      '1': 'accesscontrolallowcredentials',
      '3': 387134634,
      '4': 1,
      '5': 8,
      '10': 'accesscontrolallowcredentials'
    },
    {
      '1': 'accesscontrolallowheaders',
      '3': 392007202,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyAccessControlAllowHeaders',
      '10': 'accesscontrolallowheaders'
    },
    {
      '1': 'accesscontrolallowmethods',
      '3': 73974706,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyAccessControlAllowMethods',
      '10': 'accesscontrolallowmethods'
    },
    {
      '1': 'accesscontrolalloworigins',
      '3': 454589053,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyAccessControlAllowOrigins',
      '10': 'accesscontrolalloworigins'
    },
    {
      '1': 'accesscontrolexposeheaders',
      '3': 39057771,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyAccessControlExposeHeaders',
      '10': 'accesscontrolexposeheaders'
    },
    {
      '1': 'accesscontrolmaxagesec',
      '3': 528754691,
      '4': 1,
      '5': 5,
      '10': 'accesscontrolmaxagesec'
    },
    {
      '1': 'originoverride',
      '3': 297680576,
      '4': 1,
      '5': 8,
      '10': 'originoverride'
    },
  ],
};

/// Descriptor for `ResponseHeadersPolicyCorsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyCorsConfigDescriptor = $convert.base64Decode(
    'Ch9SZXNwb25zZUhlYWRlcnNQb2xpY3lDb3JzQ29uZmlnEkgKHWFjY2Vzc2NvbnRyb2xhbGxvd2'
    'NyZWRlbnRpYWxzGKrpzLgBIAEoCFIdYWNjZXNzY29udHJvbGFsbG93Y3JlZGVudGlhbHMSfAoZ'
    'YWNjZXNzY29udHJvbGFsbG93aGVhZGVycxiinPa6ASABKAsyOi5jbG91ZGZyb250LlJlc3Bvbn'
    'NlSGVhZGVyc1BvbGljeUFjY2Vzc0NvbnRyb2xBbGxvd0hlYWRlcnNSGWFjY2Vzc2NvbnRyb2xh'
    'bGxvd2hlYWRlcnMSewoZYWNjZXNzY29udHJvbGFsbG93bWV0aG9kcxiyh6MjIAEoCzI6LmNsb3'
    'VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9saWN5QWNjZXNzQ29udHJvbEFsbG93TWV0aG9kc1IZ'
    'YWNjZXNzY29udHJvbGFsbG93bWV0aG9kcxJ8ChlhY2Nlc3Njb250cm9sYWxsb3dvcmlnaW5zGP'
    '304dgBIAEoCzI6LmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9saWN5QWNjZXNzQ29udHJv'
    'bEFsbG93T3JpZ2luc1IZYWNjZXNzY29udHJvbGFsbG93b3JpZ2lucxJ+ChphY2Nlc3Njb250cm'
    '9sZXhwb3NlaGVhZGVycxjr8s8SIAEoCzI7LmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9s'
    'aWN5QWNjZXNzQ29udHJvbEV4cG9zZUhlYWRlcnNSGmFjY2Vzc2NvbnRyb2xleHBvc2VoZWFkZX'
    'JzEjoKFmFjY2Vzc2NvbnRyb2xtYXhhZ2VzZWMYg9CQ/AEgASgFUhZhY2Nlc3Njb250cm9sbWF4'
    'YWdlc2VjEioKDm9yaWdpbm92ZXJyaWRlGMD9+I0BIAEoCFIOb3JpZ2lub3ZlcnJpZGU=');

@$core.Deprecated('Use responseHeadersPolicyCustomHeaderDescriptor instead')
const ResponseHeadersPolicyCustomHeader$json = {
  '1': 'ResponseHeadersPolicyCustomHeader',
  '2': [
    {'1': 'header', '3': 290429313, '4': 1, '5': 9, '10': 'header'},
    {'1': 'override', '3': 139255902, '4': 1, '5': 8, '10': 'override'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyCustomHeader`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyCustomHeaderDescriptor =
    $convert.base64Decode(
        'CiFSZXNwb25zZUhlYWRlcnNQb2xpY3lDdXN0b21IZWFkZXISGgoGaGVhZGVyGIGzvooBIAEoCV'
        'IGaGVhZGVyEh0KCG92ZXJyaWRlGN7As0IgASgIUghvdmVycmlkZRIYCgV2YWx1ZRjr8p+KASAB'
        'KAlSBXZhbHVl');

@$core.Deprecated(
    'Use responseHeadersPolicyCustomHeadersConfigDescriptor instead')
const ResponseHeadersPolicyCustomHeadersConfig$json = {
  '1': 'ResponseHeadersPolicyCustomHeadersConfig',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyCustomHeader',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyCustomHeadersConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyCustomHeadersConfigDescriptor =
    $convert.base64Decode(
        'CihSZXNwb25zZUhlYWRlcnNQb2xpY3lDdXN0b21IZWFkZXJzQ29uZmlnEkYKBWl0ZW1zGLDw2A'
        'EgAygLMi0uY2xvdWRmcm9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lDdXN0b21IZWFkZXJSBWl0'
        'ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use responseHeadersPolicyFrameOptionsDescriptor instead')
const ResponseHeadersPolicyFrameOptions$json = {
  '1': 'ResponseHeadersPolicyFrameOptions',
  '2': [
    {
      '1': 'frameoption',
      '3': 494597996,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FrameOptionsList',
      '10': 'frameoption'
    },
    {'1': 'override', '3': 139255902, '4': 1, '5': 8, '10': 'override'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyFrameOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyFrameOptionsDescriptor =
    $convert.base64Decode(
        'CiFSZXNwb25zZUhlYWRlcnNQb2xpY3lGcmFtZU9wdGlvbnMSQgoLZnJhbWVvcHRpb24Y7O7r6w'
        'EgASgOMhwuY2xvdWRmcm9udC5GcmFtZU9wdGlvbnNMaXN0UgtmcmFtZW9wdGlvbhIdCghvdmVy'
        'cmlkZRjewLNCIAEoCFIIb3ZlcnJpZGU=');

@$core.Deprecated('Use responseHeadersPolicyInUseDescriptor instead')
const ResponseHeadersPolicyInUse$json = {
  '1': 'ResponseHeadersPolicyInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyInUseDescriptor =
    $convert.base64Decode(
        'ChpSZXNwb25zZUhlYWRlcnNQb2xpY3lJblVzZRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use responseHeadersPolicyListDescriptor instead')
const ResponseHeadersPolicyList$json = {
  '1': 'ResponseHeadersPolicyList',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicySummary',
      '10': 'items'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyListDescriptor = $convert.base64Decode(
    'ChlSZXNwb25zZUhlYWRlcnNQb2xpY3lMaXN0EkEKBWl0ZW1zGLDw2AEgAygLMiguY2xvdWRmcm'
    '9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lTdW1tYXJ5UgVpdGVtcxIeCghtYXhpdGVtcxiU1trx'
    'ASABKAVSCG1heGl0ZW1zEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2VyEh0KCH'
    'F1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use responseHeadersPolicyReferrerPolicyDescriptor instead')
const ResponseHeadersPolicyReferrerPolicy$json = {
  '1': 'ResponseHeadersPolicyReferrerPolicy',
  '2': [
    {'1': 'override', '3': 139255902, '4': 1, '5': 8, '10': 'override'},
    {
      '1': 'referrerpolicy',
      '3': 413520665,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ReferrerPolicyList',
      '10': 'referrerpolicy'
    },
  ],
};

/// Descriptor for `ResponseHeadersPolicyReferrerPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyReferrerPolicyDescriptor =
    $convert.base64Decode(
        'CiNSZXNwb25zZUhlYWRlcnNQb2xpY3lSZWZlcnJlclBvbGljeRIdCghvdmVycmlkZRjewLNCIA'
        'EoCFIIb3ZlcnJpZGUSSgoOcmVmZXJyZXJwb2xpY3kYmaaXxQEgASgOMh4uY2xvdWRmcm9udC5S'
        'ZWZlcnJlclBvbGljeUxpc3RSDnJlZmVycmVycG9saWN5');

@$core.Deprecated('Use responseHeadersPolicyRemoveHeaderDescriptor instead')
const ResponseHeadersPolicyRemoveHeader$json = {
  '1': 'ResponseHeadersPolicyRemoveHeader',
  '2': [
    {'1': 'header', '3': 290429313, '4': 1, '5': 9, '10': 'header'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyRemoveHeader`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyRemoveHeaderDescriptor =
    $convert.base64Decode(
        'CiFSZXNwb25zZUhlYWRlcnNQb2xpY3lSZW1vdmVIZWFkZXISGgoGaGVhZGVyGIGzvooBIAEoCV'
        'IGaGVhZGVy');

@$core.Deprecated(
    'Use responseHeadersPolicyRemoveHeadersConfigDescriptor instead')
const ResponseHeadersPolicyRemoveHeadersConfig$json = {
  '1': 'ResponseHeadersPolicyRemoveHeadersConfig',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyRemoveHeader',
      '10': 'items'
    },
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyRemoveHeadersConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyRemoveHeadersConfigDescriptor =
    $convert.base64Decode(
        'CihSZXNwb25zZUhlYWRlcnNQb2xpY3lSZW1vdmVIZWFkZXJzQ29uZmlnEkYKBWl0ZW1zGLDw2A'
        'EgAygLMi0uY2xvdWRmcm9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lSZW1vdmVIZWFkZXJSBWl0'
        'ZW1zEh0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated(
    'Use responseHeadersPolicySecurityHeadersConfigDescriptor instead')
const ResponseHeadersPolicySecurityHeadersConfig$json = {
  '1': 'ResponseHeadersPolicySecurityHeadersConfig',
  '2': [
    {
      '1': 'contentsecuritypolicy',
      '3': 531505511,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyContentSecurityPolicy',
      '10': 'contentsecuritypolicy'
    },
    {
      '1': 'contenttypeoptions',
      '3': 401984701,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyContentTypeOptions',
      '10': 'contenttypeoptions'
    },
    {
      '1': 'frameoptions',
      '3': 9632457,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyFrameOptions',
      '10': 'frameoptions'
    },
    {
      '1': 'referrerpolicy',
      '3': 413520665,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyReferrerPolicy',
      '10': 'referrerpolicy'
    },
    {
      '1': 'stricttransportsecurity',
      '3': 262335824,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyStrictTransportSecurity',
      '10': 'stricttransportsecurity'
    },
    {
      '1': 'xssprotection',
      '3': 12055347,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyXSSProtection',
      '10': 'xssprotection'
    },
  ],
};

/// Descriptor for `ResponseHeadersPolicySecurityHeadersConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicySecurityHeadersConfigDescriptor = $convert.base64Decode(
    'CipSZXNwb25zZUhlYWRlcnNQb2xpY3lTZWN1cml0eUhlYWRlcnNDb25maWcScAoVY29udGVudH'
    'NlY3VyaXR5cG9saWN5GOfCuP0BIAEoCzI2LmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9s'
    'aWN5Q29udGVudFNlY3VyaXR5UG9saWN5UhVjb250ZW50c2VjdXJpdHlwb2xpY3kSZwoSY29udG'
    'VudHR5cGVvcHRpb25zGL2Z178BIAEoCzIzLmNsb3VkZnJvbnQuUmVzcG9uc2VIZWFkZXJzUG9s'
    'aWN5Q29udGVudFR5cGVPcHRpb25zUhJjb250ZW50dHlwZW9wdGlvbnMSVAoMZnJhbWVvcHRpb2'
    '5zGMn1ywQgASgLMi0uY2xvdWRmcm9udC5SZXNwb25zZUhlYWRlcnNQb2xpY3lGcmFtZU9wdGlv'
    'bnNSDGZyYW1lb3B0aW9ucxJbCg5yZWZlcnJlcnBvbGljeRiZppfFASABKAsyLy5jbG91ZGZyb2'
    '50LlJlc3BvbnNlSGVhZGVyc1BvbGljeVJlZmVycmVyUG9saWN5Ug5yZWZlcnJlcnBvbGljeRJ1'
    'ChdzdHJpY3R0cmFuc3BvcnRzZWN1cml0eRjQ2ot9IAEoCzI4LmNsb3VkZnJvbnQuUmVzcG9uc2'
    'VIZWFkZXJzUG9saWN5U3RyaWN0VHJhbnNwb3J0U2VjdXJpdHlSF3N0cmljdHRyYW5zcG9ydHNl'
    'Y3VyaXR5ElcKDXhzc3Byb3RlY3Rpb24Ys+bfBSABKAsyLi5jbG91ZGZyb250LlJlc3BvbnNlSG'
    'VhZGVyc1BvbGljeVhTU1Byb3RlY3Rpb25SDXhzc3Byb3RlY3Rpb24=');

@$core.Deprecated(
    'Use responseHeadersPolicyServerTimingHeadersConfigDescriptor instead')
const ResponseHeadersPolicyServerTimingHeadersConfig$json = {
  '1': 'ResponseHeadersPolicyServerTimingHeadersConfig',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'samplingrate', '3': 272929747, '4': 1, '5': 1, '10': 'samplingrate'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyServerTimingHeadersConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    responseHeadersPolicyServerTimingHeadersConfigDescriptor =
    $convert.base64Decode(
        'Ci5SZXNwb25zZUhlYWRlcnNQb2xpY3lTZXJ2ZXJUaW1pbmdIZWFkZXJzQ29uZmlnEhwKB2VuYW'
        'JsZWQYv8ib5AEgASgIUgdlbmFibGVkEiYKDHNhbXBsaW5ncmF0ZRjTp5KCASABKAFSDHNhbXBs'
        'aW5ncmF0ZQ==');

@$core.Deprecated(
    'Use responseHeadersPolicyStrictTransportSecurityDescriptor instead')
const ResponseHeadersPolicyStrictTransportSecurity$json = {
  '1': 'ResponseHeadersPolicyStrictTransportSecurity',
  '2': [
    {
      '1': 'accesscontrolmaxagesec',
      '3': 528754691,
      '4': 1,
      '5': 5,
      '10': 'accesscontrolmaxagesec'
    },
    {
      '1': 'includesubdomains',
      '3': 233672837,
      '4': 1,
      '5': 8,
      '10': 'includesubdomains'
    },
    {'1': 'override', '3': 139255902, '4': 1, '5': 8, '10': 'override'},
    {'1': 'preload', '3': 99847847, '4': 1, '5': 8, '10': 'preload'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyStrictTransportSecurity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    responseHeadersPolicyStrictTransportSecurityDescriptor =
    $convert.base64Decode(
        'CixSZXNwb25zZUhlYWRlcnNQb2xpY3lTdHJpY3RUcmFuc3BvcnRTZWN1cml0eRI6ChZhY2Nlc3'
        'Njb250cm9sbWF4YWdlc2VjGIPQkPwBIAEoBVIWYWNjZXNzY29udHJvbG1heGFnZXNlYxIvChFp'
        'bmNsdWRlc3ViZG9tYWlucxiFobZvIAEoCFIRaW5jbHVkZXN1YmRvbWFpbnMSHQoIb3ZlcnJpZG'
        'UY3sCzQiABKAhSCG92ZXJyaWRlEhsKB3ByZWxvYWQYp53OLyABKAhSB3ByZWxvYWQ=');

@$core.Deprecated('Use responseHeadersPolicySummaryDescriptor instead')
const ResponseHeadersPolicySummary$json = {
  '1': 'ResponseHeadersPolicySummary',
  '2': [
    {
      '1': 'responseheaderspolicy',
      '3': 418204719,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicy',
      '10': 'responseheaderspolicy'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ResponseHeadersPolicyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ResponseHeadersPolicySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicySummaryDescriptor = $convert.base64Decode(
    'ChxSZXNwb25zZUhlYWRlcnNQb2xpY3lTdW1tYXJ5ElsKFXJlc3BvbnNlaGVhZGVyc3BvbGljeR'
    'ivmLXHASABKAsyIS5jbG91ZGZyb250LlJlc3BvbnNlSGVhZGVyc1BvbGljeVIVcmVzcG9uc2Vo'
    'ZWFkZXJzcG9saWN5Ej0KBHR5cGUY7qDXigEgASgOMiUuY2xvdWRmcm9udC5SZXNwb25zZUhlYW'
    'RlcnNQb2xpY3lUeXBlUgR0eXBl');

@$core.Deprecated('Use responseHeadersPolicyXSSProtectionDescriptor instead')
const ResponseHeadersPolicyXSSProtection$json = {
  '1': 'ResponseHeadersPolicyXSSProtection',
  '2': [
    {'1': 'modeblock', '3': 238813784, '4': 1, '5': 8, '10': 'modeblock'},
    {'1': 'override', '3': 139255902, '4': 1, '5': 8, '10': 'override'},
    {'1': 'protection', '3': 411601691, '4': 1, '5': 8, '10': 'protection'},
    {'1': 'reporturi', '3': 176037092, '4': 1, '5': 9, '10': 'reporturi'},
  ],
};

/// Descriptor for `ResponseHeadersPolicyXSSProtection`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseHeadersPolicyXSSProtectionDescriptor =
    $convert.base64Decode(
        'CiJSZXNwb25zZUhlYWRlcnNQb2xpY3lYU1NQcm90ZWN0aW9uEh8KCW1vZGVibG9jaxjYhPBxIA'
        'EoCFIJbW9kZWJsb2NrEh0KCG92ZXJyaWRlGN7As0IgASgIUghvdmVycmlkZRIiCgpwcm90ZWN0'
        'aW9uGJuWosQBIAEoCFIKcHJvdGVjdGlvbhIfCglyZXBvcnR1cmkY5Ln4UyABKAlSCXJlcG9ydH'
        'VyaQ==');

@$core.Deprecated('Use restrictionsDescriptor instead')
const Restrictions$json = {
  '1': 'Restrictions',
  '2': [
    {
      '1': 'georestriction',
      '3': 358834549,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.GeoRestriction',
      '10': 'georestriction'
    },
  ],
};

/// Descriptor for `Restrictions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restrictionsDescriptor = $convert.base64Decode(
    'CgxSZXN0cmljdGlvbnMSRgoOZ2VvcmVzdHJpY3Rpb24Y9cKNqwEgASgLMhouY2xvdWRmcm9udC'
    '5HZW9SZXN0cmljdGlvblIOZ2VvcmVzdHJpY3Rpb24=');

@$core.Deprecated('Use s3OriginDescriptor instead')
const S3Origin$json = {
  '1': 'S3Origin',
  '2': [
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'originaccessidentity',
      '3': 226897844,
      '4': 1,
      '5': 9,
      '10': 'originaccessidentity'
    },
  ],
};

/// Descriptor for `S3Origin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3OriginDescriptor = $convert.base64Decode(
    'CghTM09yaWdpbhIhCgpkb21haW5uYW1lGOvN+FwgASgJUgpkb21haW5uYW1lEjUKFG9yaWdpbm'
    'FjY2Vzc2lkZW50aXR5GLTfmGwgASgJUhRvcmlnaW5hY2Nlc3NpZGVudGl0eQ==');

@$core.Deprecated('Use s3OriginConfigDescriptor instead')
const S3OriginConfig$json = {
  '1': 'S3OriginConfig',
  '2': [
    {
      '1': 'originaccessidentity',
      '3': 226897844,
      '4': 1,
      '5': 9,
      '10': 'originaccessidentity'
    },
    {
      '1': 'originreadtimeout',
      '3': 387717023,
      '4': 1,
      '5': 5,
      '10': 'originreadtimeout'
    },
  ],
};

/// Descriptor for `S3OriginConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3OriginConfigDescriptor = $convert.base64Decode(
    'Cg5TM09yaWdpbkNvbmZpZxI1ChRvcmlnaW5hY2Nlc3NpZGVudGl0eRi035hsIAEoCVIUb3JpZ2'
    'luYWNjZXNzaWRlbnRpdHkSMAoRb3JpZ2lucmVhZHRpbWVvdXQYn6/wuAEgASgFUhFvcmlnaW5y'
    'ZWFkdGltZW91dA==');

@$core.Deprecated('Use sessionStickinessConfigDescriptor instead')
const SessionStickinessConfig$json = {
  '1': 'SessionStickinessConfig',
  '2': [
    {'1': 'idlettl', '3': 399608574, '4': 1, '5': 5, '10': 'idlettl'},
    {'1': 'maximumttl', '3': 447302168, '4': 1, '5': 5, '10': 'maximumttl'},
  ],
};

/// Descriptor for `SessionStickinessConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionStickinessConfigDescriptor =
    $convert.base64Decode(
        'ChdTZXNzaW9uU3RpY2tpbmVzc0NvbmZpZxIcCgdpZGxldHRsGP6Vxr4BIAEoBVIHaWRsZXR0bB'
        'IiCgptYXhpbXVtdHRsGJiUpdUBIAEoBVIKbWF4aW11bXR0bA==');

@$core.Deprecated('Use signerDescriptor instead')
const Signer$json = {
  '1': 'Signer',
  '2': [
    {
      '1': 'awsaccountnumber',
      '3': 397632853,
      '4': 1,
      '5': 9,
      '10': 'awsaccountnumber'
    },
    {
      '1': 'keypairids',
      '3': 5804971,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyPairIds',
      '10': 'keypairids'
    },
  ],
};

/// Descriptor for `Signer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signerDescriptor = $convert.base64Decode(
    'CgZTaWduZXISLgoQYXdzYWNjb3VudG51bWJlchjVys29ASABKAlSEGF3c2FjY291bnRudW1iZX'
    'ISOQoKa2V5cGFpcmlkcxirp+ICIAEoCzIWLmNsb3VkZnJvbnQuS2V5UGFpcklkc1IKa2V5cGFp'
    'cmlkcw==');

@$core.Deprecated('Use stagingDistributionDnsNamesDescriptor instead')
const StagingDistributionDnsNames$json = {
  '1': 'StagingDistributionDnsNames',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `StagingDistributionDnsNames`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stagingDistributionDnsNamesDescriptor =
    $convert.base64Decode(
        'ChtTdGFnaW5nRGlzdHJpYnV0aW9uRG5zTmFtZXMSFwoFaXRlbXMYsPDYASADKAlSBWl0ZW1zEh'
        '0KCHF1YW50aXR5GPnl3F8gASgFUghxdWFudGl0eQ==');

@$core.Deprecated('Use stagingDistributionInUseDescriptor instead')
const StagingDistributionInUse$json = {
  '1': 'StagingDistributionInUse',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StagingDistributionInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stagingDistributionInUseDescriptor =
    $convert.base64Decode(
        'ChhTdGFnaW5nRGlzdHJpYnV0aW9uSW5Vc2USGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use statusCodesDescriptor instead')
const StatusCodes$json = {
  '1': 'StatusCodes',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 5, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `StatusCodes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statusCodesDescriptor = $convert.base64Decode(
    'CgtTdGF0dXNDb2RlcxIXCgVpdGVtcxiw8NgBIAMoBVIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXy'
    'ABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use streamingDistributionDescriptor instead')
const StreamingDistribution$json = {
  '1': 'StreamingDistribution',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'activetrustedsigners',
      '3': 356944374,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ActiveTrustedSigners',
      '10': 'activetrustedsigners'
    },
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'streamingdistributionconfig',
      '3': 291115944,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionConfig',
      '10': 'streamingdistributionconfig'
    },
  ],
};

/// Descriptor for `StreamingDistribution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingDistributionDescriptor = $convert.base64Decode(
    'ChVTdHJlYW1pbmdEaXN0cmlidXRpb24SFAoDYXJuGJ2cr70BIAEoCVIDYXJuElgKFGFjdGl2ZX'
    'RydXN0ZWRzaWduZXJzGPaTmqoBIAEoCzIgLmNsb3VkZnJvbnQuQWN0aXZlVHJ1c3RlZFNpZ25l'
    'cnNSFGFjdGl2ZXRydXN0ZWRzaWduZXJzEiEKCmRvbWFpbm5hbWUY6834XCABKAlSCmRvbWFpbm'
    '5hbWUSEgoCaWQYgfKitwEgASgJUgJpZBItChBsYXN0bW9kaWZpZWR0aW1lGOCC/HAgASgJUhBs'
    'YXN0bW9kaWZpZWR0aW1lEhkKBnN0YXR1cxiQ5PsCIAEoCVIGc3RhdHVzEm0KG3N0cmVhbWluZ2'
    'Rpc3RyaWJ1dGlvbmNvbmZpZxiop+iKASABKAsyJy5jbG91ZGZyb250LlN0cmVhbWluZ0Rpc3Ry'
    'aWJ1dGlvbkNvbmZpZ1Ibc3RyZWFtaW5nZGlzdHJpYnV0aW9uY29uZmln');

@$core.Deprecated('Use streamingDistributionAlreadyExistsDescriptor instead')
const StreamingDistributionAlreadyExists$json = {
  '1': 'StreamingDistributionAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StreamingDistributionAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingDistributionAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CiJTdHJlYW1pbmdEaXN0cmlidXRpb25BbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UYhbO7cCABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use streamingDistributionConfigDescriptor instead')
const StreamingDistributionConfig$json = {
  '1': 'StreamingDistributionConfig',
  '2': [
    {
      '1': 'aliases',
      '3': 476693696,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Aliases',
      '10': 'aliases'
    },
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'logging',
      '3': 65655615,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingLoggingConfig',
      '10': 'logging'
    },
    {
      '1': 'priceclass',
      '3': 488692315,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.PriceClass',
      '10': 'priceclass'
    },
    {
      '1': 's3origin',
      '3': 197642510,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.S3Origin',
      '10': 's3origin'
    },
    {
      '1': 'trustedsigners',
      '3': 82003176,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustedSigners',
      '10': 'trustedsigners'
    },
  ],
};

/// Descriptor for `StreamingDistributionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingDistributionConfigDescriptor = $convert.base64Decode(
    'ChtTdHJlYW1pbmdEaXN0cmlidXRpb25Db25maWcSMQoHYWxpYXNlcxjAiafjASABKAsyEy5jbG'
    '91ZGZyb250LkFsaWFzZXNSB2FsaWFzZXMSKwoPY2FsbGVycmVmZXJlbmNlGJiZjUggASgJUg9j'
    'YWxsZXJyZWZlcmVuY2USHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1lbnQSHAoHZW5hYmxlZB'
    'i/yJvkASABKAhSB2VuYWJsZWQSPwoHbG9nZ2luZxi/pqcfIAEoCzIiLmNsb3VkZnJvbnQuU3Ry'
    'ZWFtaW5nTG9nZ2luZ0NvbmZpZ1IHbG9nZ2luZxI6CgpwcmljZWNsYXNzGNu0g+kBIAEoDjIWLm'
    'Nsb3VkZnJvbnQuUHJpY2VDbGFzc1IKcHJpY2VjbGFzcxIzCghzM29yaWdpbhiOkp9eIAEoCzIU'
    'LmNsb3VkZnJvbnQuUzNPcmlnaW5SCHMzb3JpZ2luEkUKDnRydXN0ZWRzaWduZXJzGOiJjScgAS'
    'gLMhouY2xvdWRmcm9udC5UcnVzdGVkU2lnbmVyc1IOdHJ1c3RlZHNpZ25lcnM=');

@$core.Deprecated('Use streamingDistributionConfigWithTagsDescriptor instead')
const StreamingDistributionConfigWithTags$json = {
  '1': 'StreamingDistributionConfigWithTags',
  '2': [
    {
      '1': 'streamingdistributionconfig',
      '3': 291115944,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionConfig',
      '10': 'streamingdistributionconfig'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `StreamingDistributionConfigWithTags`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingDistributionConfigWithTagsDescriptor =
    $convert.base64Decode(
        'CiNTdHJlYW1pbmdEaXN0cmlidXRpb25Db25maWdXaXRoVGFncxJtChtzdHJlYW1pbmdkaXN0cm'
        'lidXRpb25jb25maWcYqKfoigEgASgLMicuY2xvdWRmcm9udC5TdHJlYW1pbmdEaXN0cmlidXRp'
        'b25Db25maWdSG3N0cmVhbWluZ2Rpc3RyaWJ1dGlvbmNvbmZpZxIoCgR0YWdzGMHB9rUBIAEoCz'
        'IQLmNsb3VkZnJvbnQuVGFnc1IEdGFncw==');

@$core.Deprecated('Use streamingDistributionListDescriptor instead')
const StreamingDistributionList$json = {
  '1': 'StreamingDistributionList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionSummary',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `StreamingDistributionList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingDistributionListDescriptor = $convert.base64Decode(
    'ChlTdHJlYW1pbmdEaXN0cmlidXRpb25MaXN0EiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3'
    'RydW5jYXRlZBJBCgVpdGVtcxiw8NgBIAMoCzIoLmNsb3VkZnJvbnQuU3RyZWFtaW5nRGlzdHJp'
    'YnV0aW9uU3VtbWFyeVIFaXRlbXMSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aX'
    'RlbXMYlNba8QEgASgFUghtYXhpdGVtcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1h'
    'cmtlchIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use streamingDistributionNotDisabledDescriptor instead')
const StreamingDistributionNotDisabled$json = {
  '1': 'StreamingDistributionNotDisabled',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StreamingDistributionNotDisabled`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingDistributionNotDisabledDescriptor =
    $convert.base64Decode(
        'CiBTdHJlYW1pbmdEaXN0cmlidXRpb25Ob3REaXNhYmxlZBIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use streamingDistributionSummaryDescriptor instead')
const StreamingDistributionSummary$json = {
  '1': 'StreamingDistributionSummary',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'aliases',
      '3': 476693696,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Aliases',
      '10': 'aliases'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'domainname', '3': 194914027, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'priceclass',
      '3': 488692315,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.PriceClass',
      '10': 'priceclass'
    },
    {
      '1': 's3origin',
      '3': 197642510,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.S3Origin',
      '10': 's3origin'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'trustedsigners',
      '3': 82003176,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustedSigners',
      '10': 'trustedsigners'
    },
  ],
};

/// Descriptor for `StreamingDistributionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingDistributionSummaryDescriptor = $convert.base64Decode(
    'ChxTdHJlYW1pbmdEaXN0cmlidXRpb25TdW1tYXJ5EhQKA2FybhidnK+9ASABKAlSA2FybhIxCg'
    'dhbGlhc2VzGMCJp+MBIAEoCzITLmNsb3VkZnJvbnQuQWxpYXNlc1IHYWxpYXNlcxIcCgdjb21t'
    'ZW50GP+/vsIBIAEoCVIHY29tbWVudBIhCgpkb21haW5uYW1lGOvN+FwgASgJUgpkb21haW5uYW'
    '1lEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVkEhIKAmlkGIHyorcBIAEoCVICaWQSLQoQ'
    'bGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZRI6CgpwcmljZWNsYX'
    'NzGNu0g+kBIAEoDjIWLmNsb3VkZnJvbnQuUHJpY2VDbGFzc1IKcHJpY2VjbGFzcxIzCghzM29y'
    'aWdpbhiOkp9eIAEoCzIULmNsb3VkZnJvbnQuUzNPcmlnaW5SCHMzb3JpZ2luEhkKBnN0YXR1cx'
    'iQ5PsCIAEoCVIGc3RhdHVzEkUKDnRydXN0ZWRzaWduZXJzGOiJjScgASgLMhouY2xvdWRmcm9u'
    'dC5UcnVzdGVkU2lnbmVyc1IOdHJ1c3RlZHNpZ25lcnM=');

@$core.Deprecated('Use streamingLoggingConfigDescriptor instead')
const StreamingLoggingConfig$json = {
  '1': 'StreamingLoggingConfig',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `StreamingLoggingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamingLoggingConfigDescriptor = $convert.base64Decode(
    'ChZTdHJlYW1pbmdMb2dnaW5nQ29uZmlnEhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EhwKB2'
    'VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVkEhoKBnByZWZpeBjqs9OCASABKAlSBnByZWZpeA==');

@$core.Deprecated('Use stringSchemaConfigDescriptor instead')
const StringSchemaConfig$json = {
  '1': 'StringSchemaConfig',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'defaultvalue', '3': 218709920, '4': 1, '5': 9, '10': 'defaultvalue'},
    {'1': 'required', '3': 300200513, '4': 1, '5': 8, '10': 'required'},
  ],
};

/// Descriptor for `StringSchemaConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stringSchemaConfigDescriptor = $convert.base64Decode(
    'ChJTdHJpbmdTY2hlbWFDb25maWcSHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1lbnQSJQoMZG'
    'VmYXVsdHZhbHVlGKD/pGggASgJUgxkZWZhdWx0dmFsdWUSHgoIcmVxdWlyZWQYweSSjwEgASgI'
    'UghyZXF1aXJlZA==');

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

@$core.Deprecated('Use tagKeysDescriptor instead')
const TagKeys$json = {
  '1': 'TagKeys',
  '2': [
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
  ],
};

/// Descriptor for `TagKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagKeysDescriptor =
    $convert.base64Decode('CgdUYWdLZXlzEhcKBWl0ZW1zGLDw2AEgAygJUgVpdGVtcw==');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resource', '3': 61921302, '4': 1, '5': 9, '10': 'resource'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Tags',
      '8': {},
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSHQoIcmVzb3VyY2UYlrDDHSABKAlSCHJlc291cmNlEi4KBH'
    'RhZ3MYwcH2tQEgASgLMhAuY2xvdWRmcm9udC5UYWdzQgSItRgBUgR0YWdz');

@$core.Deprecated('Use tagsDescriptor instead')
const Tags$json = {
  '1': 'Tags',
  '2': [
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.Tag',
      '10': 'items'
    },
  ],
};

/// Descriptor for `Tags`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagsDescriptor = $convert.base64Decode(
    'CgRUYWdzEigKBWl0ZW1zGLDw2AEgAygLMg8uY2xvdWRmcm9udC5UYWdSBWl0ZW1z');

@$core.Deprecated('Use tenantConfigDescriptor instead')
const TenantConfig$json = {
  '1': 'TenantConfig',
  '2': [
    {
      '1': 'parameterdefinitions',
      '3': 504985607,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.ParameterDefinition',
      '10': 'parameterdefinitions'
    },
  ],
};

/// Descriptor for `TenantConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tenantConfigDescriptor = $convert.base64Decode(
    'CgxUZW5hbnRDb25maWcSVwoUcGFyYW1ldGVyZGVmaW5pdGlvbnMYh/Dl8AEgAygLMh8uY2xvdW'
    'Rmcm9udC5QYXJhbWV0ZXJEZWZpbml0aW9uUhRwYXJhbWV0ZXJkZWZpbml0aW9ucw==');

@$core.Deprecated('Use testConnectionFunctionRequestDescriptor instead')
const TestConnectionFunctionRequest$json = {
  '1': 'TestConnectionFunctionRequest',
  '2': [
    {
      '1': 'connectionobject',
      '3': 262521015,
      '4': 1,
      '5': 12,
      '10': 'connectionobject'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `TestConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testConnectionFunctionRequestDescriptor = $convert.base64Decode(
    'Ch1UZXN0Q29ubmVjdGlvbkZ1bmN0aW9uUmVxdWVzdBItChBjb25uZWN0aW9ub2JqZWN0GLeBl3'
    '0gASgMUhBjb25uZWN0aW9ub2JqZWN0EhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaBjQ'
    'lrcsIAEoCVIHaWZtYXRjaBIyCgVzdGFnZRjOl9hwIAEoDjIZLmNsb3VkZnJvbnQuRnVuY3Rpb2'
    '5TdGFnZVIFc3RhZ2U=');

@$core.Deprecated('Use testConnectionFunctionResultDescriptor instead')
const TestConnectionFunctionResult$json = {
  '1': 'TestConnectionFunctionResult',
  '2': [
    {
      '1': 'connectionfunctiontestresult',
      '3': 37901549,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionTestResult',
      '8': {},
      '10': 'connectionfunctiontestresult'
    },
  ],
};

/// Descriptor for `TestConnectionFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testConnectionFunctionResultDescriptor =
    $convert.base64Decode(
        'ChxUZXN0Q29ubmVjdGlvbkZ1bmN0aW9uUmVzdWx0EnUKHGNvbm5lY3Rpb25mdW5jdGlvbnRlc3'
        'RyZXN1bHQY7amJEiABKAsyKC5jbG91ZGZyb250LkNvbm5lY3Rpb25GdW5jdGlvblRlc3RSZXN1'
        'bHRCBIi1GAFSHGNvbm5lY3Rpb25mdW5jdGlvbnRlc3RyZXN1bHQ=');

@$core.Deprecated('Use testFunctionFailedDescriptor instead')
const TestFunctionFailed$json = {
  '1': 'TestFunctionFailed',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TestFunctionFailed`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testFunctionFailedDescriptor =
    $convert.base64Decode(
        'ChJUZXN0RnVuY3Rpb25GYWlsZWQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use testFunctionRequestDescriptor instead')
const TestFunctionRequest$json = {
  '1': 'TestFunctionRequest',
  '2': [
    {'1': 'eventobject', '3': 59339613, '4': 1, '5': 12, '10': 'eventobject'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'stage',
      '3': 236325838,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.FunctionStage',
      '10': 'stage'
    },
  ],
};

/// Descriptor for `TestFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testFunctionRequestDescriptor = $convert.base64Decode(
    'ChNUZXN0RnVuY3Rpb25SZXF1ZXN0EiMKC2V2ZW50b2JqZWN0GN3mpRwgASgMUgtldmVudG9iam'
    'VjdBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNoEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUS'
    'MgoFc3RhZ2UYzpfYcCABKA4yGS5jbG91ZGZyb250LkZ1bmN0aW9uU3RhZ2VSBXN0YWdl');

@$core.Deprecated('Use testFunctionResultDescriptor instead')
const TestFunctionResult$json = {
  '1': 'TestFunctionResult',
  '2': [
    {
      '1': 'testresult',
      '3': 361341317,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TestResult',
      '8': {},
      '10': 'testresult'
    },
  ],
};

/// Descriptor for `TestFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testFunctionResultDescriptor = $convert.base64Decode(
    'ChJUZXN0RnVuY3Rpb25SZXN1bHQSQAoKdGVzdHJlc3VsdBiFw6asASABKAsyFi5jbG91ZGZyb2'
    '50LlRlc3RSZXN1bHRCBIi1GAFSCnRlc3RyZXN1bHQ=');

@$core.Deprecated('Use testResultDescriptor instead')
const TestResult$json = {
  '1': 'TestResult',
  '2': [
    {
      '1': 'computeutilization',
      '3': 247332359,
      '4': 1,
      '5': 9,
      '10': 'computeutilization'
    },
    {
      '1': 'functionerrormessage',
      '3': 424938733,
      '4': 1,
      '5': 9,
      '10': 'functionerrormessage'
    },
    {
      '1': 'functionexecutionlogs',
      '3': 344853419,
      '4': 3,
      '5': 9,
      '10': 'functionexecutionlogs'
    },
    {
      '1': 'functionoutput',
      '3': 171161345,
      '4': 1,
      '5': 9,
      '10': 'functionoutput'
    },
    {
      '1': 'functionsummary',
      '3': 523316264,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionSummary',
      '10': 'functionsummary'
    },
  ],
};

/// Descriptor for `TestResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testResultDescriptor = $convert.base64Decode(
    'CgpUZXN0UmVzdWx0EjEKEmNvbXB1dGV1dGlsaXphdGlvbhiH/Pd1IAEoCVISY29tcHV0ZXV0aW'
    'xpemF0aW9uEjYKFGZ1bmN0aW9uZXJyb3JtZXNzYWdlGO2Z0MoBIAEoCVIUZnVuY3Rpb25lcnJv'
    'cm1lc3NhZ2USOAoVZnVuY3Rpb25leGVjdXRpb25sb2dzGKuXuKQBIAMoCVIVZnVuY3Rpb25leG'
    'VjdXRpb25sb2dzEikKDmZ1bmN0aW9ub3V0cHV0GIHuzlEgASgJUg5mdW5jdGlvbm91dHB1dBJJ'
    'Cg9mdW5jdGlvbnN1bW1hcnkYqNjE+QEgASgLMhsuY2xvdWRmcm9udC5GdW5jdGlvblN1bW1hcn'
    'lSD2Z1bmN0aW9uc3VtbWFyeQ==');

@$core.Deprecated('Use tooLongCSPInResponseHeadersPolicyDescriptor instead')
const TooLongCSPInResponseHeadersPolicy$json = {
  '1': 'TooLongCSPInResponseHeadersPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooLongCSPInResponseHeadersPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooLongCSPInResponseHeadersPolicyDescriptor =
    $convert.base64Decode(
        'CiFUb29Mb25nQ1NQSW5SZXNwb25zZUhlYWRlcnNQb2xpY3kSGwoHbWVzc2FnZRiFs7twIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use tooManyCacheBehaviorsDescriptor instead')
const TooManyCacheBehaviors$json = {
  '1': 'TooManyCacheBehaviors',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCacheBehaviors`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyCacheBehaviorsDescriptor = $convert.base64Decode(
    'ChVUb29NYW55Q2FjaGVCZWhhdmlvcnMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use tooManyCachePoliciesDescriptor instead')
const TooManyCachePolicies$json = {
  '1': 'TooManyCachePolicies',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCachePolicies`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyCachePoliciesDescriptor =
    $convert.base64Decode(
        'ChRUb29NYW55Q2FjaGVQb2xpY2llcxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use tooManyCertificatesDescriptor instead')
const TooManyCertificates$json = {
  '1': 'TooManyCertificates',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCertificates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyCertificatesDescriptor =
    $convert.base64Decode(
        'ChNUb29NYW55Q2VydGlmaWNhdGVzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core
    .Deprecated('Use tooManyCloudFrontOriginAccessIdentitiesDescriptor instead')
const TooManyCloudFrontOriginAccessIdentities$json = {
  '1': 'TooManyCloudFrontOriginAccessIdentities',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCloudFrontOriginAccessIdentities`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyCloudFrontOriginAccessIdentitiesDescriptor =
    $convert.base64Decode(
        'CidUb29NYW55Q2xvdWRGcm9udE9yaWdpbkFjY2Vzc0lkZW50aXRpZXMSGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use tooManyContinuousDeploymentPoliciesDescriptor instead')
const TooManyContinuousDeploymentPolicies$json = {
  '1': 'TooManyContinuousDeploymentPolicies',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyContinuousDeploymentPolicies`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyContinuousDeploymentPoliciesDescriptor =
    $convert.base64Decode(
        'CiNUb29NYW55Q29udGludW91c0RlcGxveW1lbnRQb2xpY2llcxIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use tooManyCookieNamesInWhiteListDescriptor instead')
const TooManyCookieNamesInWhiteList$json = {
  '1': 'TooManyCookieNamesInWhiteList',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCookieNamesInWhiteList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyCookieNamesInWhiteListDescriptor =
    $convert.base64Decode(
        'Ch1Ub29NYW55Q29va2llTmFtZXNJbldoaXRlTGlzdBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use tooManyCookiesInCachePolicyDescriptor instead')
const TooManyCookiesInCachePolicy$json = {
  '1': 'TooManyCookiesInCachePolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCookiesInCachePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyCookiesInCachePolicyDescriptor =
    $convert.base64Decode(
        'ChtUb29NYW55Q29va2llc0luQ2FjaGVQb2xpY3kSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use tooManyCookiesInOriginRequestPolicyDescriptor instead')
const TooManyCookiesInOriginRequestPolicy$json = {
  '1': 'TooManyCookiesInOriginRequestPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCookiesInOriginRequestPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyCookiesInOriginRequestPolicyDescriptor =
    $convert.base64Decode(
        'CiNUb29NYW55Q29va2llc0luT3JpZ2luUmVxdWVzdFBvbGljeRIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use tooManyCustomHeadersInResponseHeadersPolicyDescriptor instead')
const TooManyCustomHeadersInResponseHeadersPolicy$json = {
  '1': 'TooManyCustomHeadersInResponseHeadersPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyCustomHeadersInResponseHeadersPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyCustomHeadersInResponseHeadersPolicyDescriptor =
    $convert.base64Decode(
        'CitUb29NYW55Q3VzdG9tSGVhZGVyc0luUmVzcG9uc2VIZWFkZXJzUG9saWN5EhsKB21lc3NhZ2'
        'UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyDistributionCNAMEsDescriptor instead')
const TooManyDistributionCNAMEs$json = {
  '1': 'TooManyDistributionCNAMEs',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionCNAMEs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyDistributionCNAMEsDescriptor =
    $convert.base64Decode(
        'ChlUb29NYW55RGlzdHJpYnV0aW9uQ05BTUVzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use tooManyDistributionsDescriptor instead')
const TooManyDistributions$json = {
  '1': 'TooManyDistributions',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyDistributionsDescriptor =
    $convert.base64Decode(
        'ChRUb29NYW55RGlzdHJpYnV0aW9ucxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use tooManyDistributionsAssociatedToCachePolicyDescriptor instead')
const TooManyDistributionsAssociatedToCachePolicy$json = {
  '1': 'TooManyDistributionsAssociatedToCachePolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsAssociatedToCachePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsAssociatedToCachePolicyDescriptor =
    $convert.base64Decode(
        'CitUb29NYW55RGlzdHJpYnV0aW9uc0Fzc29jaWF0ZWRUb0NhY2hlUG9saWN5EhsKB21lc3NhZ2'
        'UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use tooManyDistributionsAssociatedToFieldLevelEncryptionConfigDescriptor instead')
const TooManyDistributionsAssociatedToFieldLevelEncryptionConfig$json = {
  '1': 'TooManyDistributionsAssociatedToFieldLevelEncryptionConfig',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsAssociatedToFieldLevelEncryptionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsAssociatedToFieldLevelEncryptionConfigDescriptor =
    $convert.base64Decode(
        'CjpUb29NYW55RGlzdHJpYnV0aW9uc0Fzc29jaWF0ZWRUb0ZpZWxkTGV2ZWxFbmNyeXB0aW9uQ2'
        '9uZmlnEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use tooManyDistributionsAssociatedToKeyGroupDescriptor instead')
const TooManyDistributionsAssociatedToKeyGroup$json = {
  '1': 'TooManyDistributionsAssociatedToKeyGroup',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsAssociatedToKeyGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyDistributionsAssociatedToKeyGroupDescriptor =
    $convert.base64Decode(
        'CihUb29NYW55RGlzdHJpYnV0aW9uc0Fzc29jaWF0ZWRUb0tleUdyb3VwEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use tooManyDistributionsAssociatedToOriginAccessControlDescriptor instead')
const TooManyDistributionsAssociatedToOriginAccessControl$json = {
  '1': 'TooManyDistributionsAssociatedToOriginAccessControl',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsAssociatedToOriginAccessControl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsAssociatedToOriginAccessControlDescriptor =
    $convert.base64Decode(
        'CjNUb29NYW55RGlzdHJpYnV0aW9uc0Fzc29jaWF0ZWRUb09yaWdpbkFjY2Vzc0NvbnRyb2wSGw'
        'oHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated(
    'Use tooManyDistributionsAssociatedToOriginRequestPolicyDescriptor instead')
const TooManyDistributionsAssociatedToOriginRequestPolicy$json = {
  '1': 'TooManyDistributionsAssociatedToOriginRequestPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsAssociatedToOriginRequestPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsAssociatedToOriginRequestPolicyDescriptor =
    $convert.base64Decode(
        'CjNUb29NYW55RGlzdHJpYnV0aW9uc0Fzc29jaWF0ZWRUb09yaWdpblJlcXVlc3RQb2xpY3kSGw'
        'oHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated(
    'Use tooManyDistributionsAssociatedToResponseHeadersPolicyDescriptor instead')
const TooManyDistributionsAssociatedToResponseHeadersPolicy$json = {
  '1': 'TooManyDistributionsAssociatedToResponseHeadersPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsAssociatedToResponseHeadersPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsAssociatedToResponseHeadersPolicyDescriptor =
    $convert.base64Decode(
        'CjVUb29NYW55RGlzdHJpYnV0aW9uc0Fzc29jaWF0ZWRUb1Jlc3BvbnNlSGVhZGVyc1BvbGljeR'
        'IbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use tooManyDistributionsWithFunctionAssociationsDescriptor instead')
const TooManyDistributionsWithFunctionAssociations$json = {
  '1': 'TooManyDistributionsWithFunctionAssociations',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsWithFunctionAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsWithFunctionAssociationsDescriptor =
    $convert.base64Decode(
        'CixUb29NYW55RGlzdHJpYnV0aW9uc1dpdGhGdW5jdGlvbkFzc29jaWF0aW9ucxIbCgdtZXNzYW'
        'dlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use tooManyDistributionsWithLambdaAssociationsDescriptor instead')
const TooManyDistributionsWithLambdaAssociations$json = {
  '1': 'TooManyDistributionsWithLambdaAssociations',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsWithLambdaAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsWithLambdaAssociationsDescriptor =
    $convert.base64Decode(
        'CipUb29NYW55RGlzdHJpYnV0aW9uc1dpdGhMYW1iZGFBc3NvY2lhdGlvbnMSGwoHbWVzc2FnZR'
        'iFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated(
    'Use tooManyDistributionsWithSingleFunctionARNDescriptor instead')
const TooManyDistributionsWithSingleFunctionARN$json = {
  '1': 'TooManyDistributionsWithSingleFunctionARN',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyDistributionsWithSingleFunctionARN`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyDistributionsWithSingleFunctionARNDescriptor = $convert.base64Decode(
        'CilUb29NYW55RGlzdHJpYnV0aW9uc1dpdGhTaW5nbGVGdW5jdGlvbkFSThIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use tooManyFieldLevelEncryptionConfigsDescriptor instead')
const TooManyFieldLevelEncryptionConfigs$json = {
  '1': 'TooManyFieldLevelEncryptionConfigs',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFieldLevelEncryptionConfigs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyFieldLevelEncryptionConfigsDescriptor =
    $convert.base64Decode(
        'CiJUb29NYW55RmllbGRMZXZlbEVuY3J5cHRpb25Db25maWdzEhsKB21lc3NhZ2UYhbO7cCABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use tooManyFieldLevelEncryptionContentTypeProfilesDescriptor instead')
const TooManyFieldLevelEncryptionContentTypeProfiles$json = {
  '1': 'TooManyFieldLevelEncryptionContentTypeProfiles',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFieldLevelEncryptionContentTypeProfiles`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyFieldLevelEncryptionContentTypeProfilesDescriptor =
    $convert.base64Decode(
        'Ci5Ub29NYW55RmllbGRMZXZlbEVuY3J5cHRpb25Db250ZW50VHlwZVByb2ZpbGVzEhsKB21lc3'
        'NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use tooManyFieldLevelEncryptionEncryptionEntitiesDescriptor instead')
const TooManyFieldLevelEncryptionEncryptionEntities$json = {
  '1': 'TooManyFieldLevelEncryptionEncryptionEntities',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFieldLevelEncryptionEncryptionEntities`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyFieldLevelEncryptionEncryptionEntitiesDescriptor =
    $convert.base64Decode(
        'Ci1Ub29NYW55RmllbGRMZXZlbEVuY3J5cHRpb25FbmNyeXB0aW9uRW50aXRpZXMSGwoHbWVzc2'
        'FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated(
    'Use tooManyFieldLevelEncryptionFieldPatternsDescriptor instead')
const TooManyFieldLevelEncryptionFieldPatterns$json = {
  '1': 'TooManyFieldLevelEncryptionFieldPatterns',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFieldLevelEncryptionFieldPatterns`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyFieldLevelEncryptionFieldPatternsDescriptor =
    $convert.base64Decode(
        'CihUb29NYW55RmllbGRMZXZlbEVuY3J5cHRpb25GaWVsZFBhdHRlcm5zEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyFieldLevelEncryptionProfilesDescriptor instead')
const TooManyFieldLevelEncryptionProfiles$json = {
  '1': 'TooManyFieldLevelEncryptionProfiles',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFieldLevelEncryptionProfiles`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyFieldLevelEncryptionProfilesDescriptor =
    $convert.base64Decode(
        'CiNUb29NYW55RmllbGRMZXZlbEVuY3J5cHRpb25Qcm9maWxlcxIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use tooManyFieldLevelEncryptionQueryArgProfilesDescriptor instead')
const TooManyFieldLevelEncryptionQueryArgProfiles$json = {
  '1': 'TooManyFieldLevelEncryptionQueryArgProfiles',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFieldLevelEncryptionQueryArgProfiles`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyFieldLevelEncryptionQueryArgProfilesDescriptor =
    $convert.base64Decode(
        'CitUb29NYW55RmllbGRMZXZlbEVuY3J5cHRpb25RdWVyeUFyZ1Byb2ZpbGVzEhsKB21lc3NhZ2'
        'UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyFunctionAssociationsDescriptor instead')
const TooManyFunctionAssociations$json = {
  '1': 'TooManyFunctionAssociations',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFunctionAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyFunctionAssociationsDescriptor =
    $convert.base64Decode(
        'ChtUb29NYW55RnVuY3Rpb25Bc3NvY2lhdGlvbnMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use tooManyFunctionsDescriptor instead')
const TooManyFunctions$json = {
  '1': 'TooManyFunctions',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFunctions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyFunctionsDescriptor = $convert.base64Decode(
    'ChBUb29NYW55RnVuY3Rpb25zEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyHeadersInCachePolicyDescriptor instead')
const TooManyHeadersInCachePolicy$json = {
  '1': 'TooManyHeadersInCachePolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyHeadersInCachePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyHeadersInCachePolicyDescriptor =
    $convert.base64Decode(
        'ChtUb29NYW55SGVhZGVyc0luQ2FjaGVQb2xpY3kSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use tooManyHeadersInForwardedValuesDescriptor instead')
const TooManyHeadersInForwardedValues$json = {
  '1': 'TooManyHeadersInForwardedValues',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyHeadersInForwardedValues`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyHeadersInForwardedValuesDescriptor =
    $convert.base64Decode(
        'Ch9Ub29NYW55SGVhZGVyc0luRm9yd2FyZGVkVmFsdWVzEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use tooManyHeadersInOriginRequestPolicyDescriptor instead')
const TooManyHeadersInOriginRequestPolicy$json = {
  '1': 'TooManyHeadersInOriginRequestPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyHeadersInOriginRequestPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyHeadersInOriginRequestPolicyDescriptor =
    $convert.base64Decode(
        'CiNUb29NYW55SGVhZGVyc0luT3JpZ2luUmVxdWVzdFBvbGljeRIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use tooManyInvalidationsInProgressDescriptor instead')
const TooManyInvalidationsInProgress$json = {
  '1': 'TooManyInvalidationsInProgress',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyInvalidationsInProgress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyInvalidationsInProgressDescriptor =
    $convert.base64Decode(
        'Ch5Ub29NYW55SW52YWxpZGF0aW9uc0luUHJvZ3Jlc3MSGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use tooManyKeyGroupsDescriptor instead')
const TooManyKeyGroups$json = {
  '1': 'TooManyKeyGroups',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyKeyGroups`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyKeyGroupsDescriptor = $convert.base64Decode(
    'ChBUb29NYW55S2V5R3JvdXBzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use tooManyKeyGroupsAssociatedToDistributionDescriptor instead')
const TooManyKeyGroupsAssociatedToDistribution$json = {
  '1': 'TooManyKeyGroupsAssociatedToDistribution',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyKeyGroupsAssociatedToDistribution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyKeyGroupsAssociatedToDistributionDescriptor =
    $convert.base64Decode(
        'CihUb29NYW55S2V5R3JvdXBzQXNzb2NpYXRlZFRvRGlzdHJpYnV0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyLambdaFunctionAssociationsDescriptor instead')
const TooManyLambdaFunctionAssociations$json = {
  '1': 'TooManyLambdaFunctionAssociations',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyLambdaFunctionAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyLambdaFunctionAssociationsDescriptor =
    $convert.base64Decode(
        'CiFUb29NYW55TGFtYmRhRnVuY3Rpb25Bc3NvY2lhdGlvbnMSGwoHbWVzc2FnZRiFs7twIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use tooManyOriginAccessControlsDescriptor instead')
const TooManyOriginAccessControls$json = {
  '1': 'TooManyOriginAccessControls',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyOriginAccessControls`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyOriginAccessControlsDescriptor =
    $convert.base64Decode(
        'ChtUb29NYW55T3JpZ2luQWNjZXNzQ29udHJvbHMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use tooManyOriginCustomHeadersDescriptor instead')
const TooManyOriginCustomHeaders$json = {
  '1': 'TooManyOriginCustomHeaders',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyOriginCustomHeaders`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyOriginCustomHeadersDescriptor =
    $convert.base64Decode(
        'ChpUb29NYW55T3JpZ2luQ3VzdG9tSGVhZGVycxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use tooManyOriginGroupsPerDistributionDescriptor instead')
const TooManyOriginGroupsPerDistribution$json = {
  '1': 'TooManyOriginGroupsPerDistribution',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyOriginGroupsPerDistribution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyOriginGroupsPerDistributionDescriptor =
    $convert.base64Decode(
        'CiJUb29NYW55T3JpZ2luR3JvdXBzUGVyRGlzdHJpYnV0aW9uEhsKB21lc3NhZ2UYhbO7cCABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyOriginRequestPoliciesDescriptor instead')
const TooManyOriginRequestPolicies$json = {
  '1': 'TooManyOriginRequestPolicies',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyOriginRequestPolicies`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyOriginRequestPoliciesDescriptor =
    $convert.base64Decode(
        'ChxUb29NYW55T3JpZ2luUmVxdWVzdFBvbGljaWVzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use tooManyOriginsDescriptor instead')
const TooManyOrigins$json = {
  '1': 'TooManyOrigins',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyOrigins`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyOriginsDescriptor = $convert.base64Decode(
    'Cg5Ub29NYW55T3JpZ2lucxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use tooManyPublicKeysDescriptor instead')
const TooManyPublicKeys$json = {
  '1': 'TooManyPublicKeys',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyPublicKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyPublicKeysDescriptor = $convert.base64Decode(
    'ChFUb29NYW55UHVibGljS2V5cxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use tooManyPublicKeysInKeyGroupDescriptor instead')
const TooManyPublicKeysInKeyGroup$json = {
  '1': 'TooManyPublicKeysInKeyGroup',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyPublicKeysInKeyGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyPublicKeysInKeyGroupDescriptor =
    $convert.base64Decode(
        'ChtUb29NYW55UHVibGljS2V5c0luS2V5R3JvdXASGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use tooManyQueryStringParametersDescriptor instead')
const TooManyQueryStringParameters$json = {
  '1': 'TooManyQueryStringParameters',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyQueryStringParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyQueryStringParametersDescriptor =
    $convert.base64Decode(
        'ChxUb29NYW55UXVlcnlTdHJpbmdQYXJhbWV0ZXJzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use tooManyQueryStringsInCachePolicyDescriptor instead')
const TooManyQueryStringsInCachePolicy$json = {
  '1': 'TooManyQueryStringsInCachePolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyQueryStringsInCachePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyQueryStringsInCachePolicyDescriptor =
    $convert.base64Decode(
        'CiBUb29NYW55UXVlcnlTdHJpbmdzSW5DYWNoZVBvbGljeRIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated(
    'Use tooManyQueryStringsInOriginRequestPolicyDescriptor instead')
const TooManyQueryStringsInOriginRequestPolicy$json = {
  '1': 'TooManyQueryStringsInOriginRequestPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyQueryStringsInOriginRequestPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyQueryStringsInOriginRequestPolicyDescriptor =
    $convert.base64Decode(
        'CihUb29NYW55UXVlcnlTdHJpbmdzSW5PcmlnaW5SZXF1ZXN0UG9saWN5EhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyRealtimeLogConfigsDescriptor instead')
const TooManyRealtimeLogConfigs$json = {
  '1': 'TooManyRealtimeLogConfigs',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyRealtimeLogConfigs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyRealtimeLogConfigsDescriptor =
    $convert.base64Decode(
        'ChlUb29NYW55UmVhbHRpbWVMb2dDb25maWdzEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated(
    'Use tooManyRemoveHeadersInResponseHeadersPolicyDescriptor instead')
const TooManyRemoveHeadersInResponseHeadersPolicy$json = {
  '1': 'TooManyRemoveHeadersInResponseHeadersPolicy',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyRemoveHeadersInResponseHeadersPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyRemoveHeadersInResponseHeadersPolicyDescriptor =
    $convert.base64Decode(
        'CitUb29NYW55UmVtb3ZlSGVhZGVyc0luUmVzcG9uc2VIZWFkZXJzUG9saWN5EhsKB21lc3NhZ2'
        'UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyResponseHeadersPoliciesDescriptor instead')
const TooManyResponseHeadersPolicies$json = {
  '1': 'TooManyResponseHeadersPolicies',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyResponseHeadersPolicies`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyResponseHeadersPoliciesDescriptor =
    $convert.base64Decode(
        'Ch5Ub29NYW55UmVzcG9uc2VIZWFkZXJzUG9saWNpZXMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use tooManyStreamingDistributionCNAMEsDescriptor instead')
const TooManyStreamingDistributionCNAMEs$json = {
  '1': 'TooManyStreamingDistributionCNAMEs',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyStreamingDistributionCNAMEs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyStreamingDistributionCNAMEsDescriptor =
    $convert.base64Decode(
        'CiJUb29NYW55U3RyZWFtaW5nRGlzdHJpYnV0aW9uQ05BTUVzEhsKB21lc3NhZ2UYhbO7cCABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyStreamingDistributionsDescriptor instead')
const TooManyStreamingDistributions$json = {
  '1': 'TooManyStreamingDistributions',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyStreamingDistributions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyStreamingDistributionsDescriptor =
    $convert.base64Decode(
        'Ch1Ub29NYW55U3RyZWFtaW5nRGlzdHJpYnV0aW9ucxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use tooManyTrustedSignersDescriptor instead')
const TooManyTrustedSigners$json = {
  '1': 'TooManyTrustedSigners',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyTrustedSigners`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyTrustedSignersDescriptor = $convert.base64Decode(
    'ChVUb29NYW55VHJ1c3RlZFNpZ25lcnMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use trafficConfigDescriptor instead')
const TrafficConfig$json = {
  '1': 'TrafficConfig',
  '2': [
    {
      '1': 'singleheaderconfig',
      '3': 307731851,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentSingleHeaderConfig',
      '10': 'singleheaderconfig'
    },
    {
      '1': 'singleweightconfig',
      '3': 428983664,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentSingleWeightConfig',
      '10': 'singleweightconfig'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ContinuousDeploymentPolicyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `TrafficConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trafficConfigDescriptor = $convert.base64Decode(
    'Cg1UcmFmZmljQ29uZmlnEmYKEnNpbmdsZWhlYWRlcmNvbmZpZxiLu96SASABKAsyMi5jbG91ZG'
    'Zyb250LkNvbnRpbnVvdXNEZXBsb3ltZW50U2luZ2xlSGVhZGVyQ29uZmlnUhJzaW5nbGVoZWFk'
    'ZXJjb25maWcSZgoSc2luZ2xld2VpZ2h0Y29uZmlnGPCKx8wBIAEoCzIyLmNsb3VkZnJvbnQuQ2'
    '9udGludW91c0RlcGxveW1lbnRTaW5nbGVXZWlnaHRDb25maWdSEnNpbmdsZXdlaWdodGNvbmZp'
    'ZxJCCgR0eXBlGO6g14oBIAEoDjIqLmNsb3VkZnJvbnQuQ29udGludW91c0RlcGxveW1lbnRQb2'
    'xpY3lUeXBlUgR0eXBl');

@$core.Deprecated('Use trustStoreDescriptor instead')
const TrustStore$json = {
  '1': 'TrustStore',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'numberofcacertificates',
      '3': 203958970,
      '4': 1,
      '5': 5,
      '10': 'numberofcacertificates'
    },
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.TrustStoreStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `TrustStore`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trustStoreDescriptor = $convert.base64Decode(
    'CgpUcnVzdFN0b3JlEhQKA2Fybhidm+2/ASABKAlSA2FybhISCgJpZBiB8qK3ASABKAlSAmlkEi'
    '0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWUSFQoEbmFtZRiH'
    '5oF/IAEoCVIEbmFtZRI5ChZudW1iZXJvZmNhY2VydGlmaWNhdGVzGLrVoGEgASgFUhZudW1iZX'
    'JvZmNhY2VydGlmaWNhdGVzEhkKBnJlYXNvbhi6gsUJIAEoCVIGcmVhc29uEjcKBnN0YXR1cxiQ'
    '5PsCIAEoDjIcLmNsb3VkZnJvbnQuVHJ1c3RTdG9yZVN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use trustStoreConfigDescriptor instead')
const TrustStoreConfig$json = {
  '1': 'TrustStoreConfig',
  '2': [
    {
      '1': 'advertisetruststorecanames',
      '3': 199979436,
      '4': 1,
      '5': 8,
      '10': 'advertisetruststorecanames'
    },
    {
      '1': 'ignorecertificateexpiry',
      '3': 64719090,
      '4': 1,
      '5': 8,
      '10': 'ignorecertificateexpiry'
    },
    {'1': 'truststoreid', '3': 67569668, '4': 1, '5': 9, '10': 'truststoreid'},
  ],
};

/// Descriptor for `TrustStoreConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trustStoreConfigDescriptor = $convert.base64Decode(
    'ChBUcnVzdFN0b3JlQ29uZmlnEkEKGmFkdmVydGlzZXRydXN0c3RvcmVjYW5hbWVzGKzjrV8gAS'
    'gIUhphZHZlcnRpc2V0cnVzdHN0b3JlY2FuYW1lcxI7ChdpZ25vcmVjZXJ0aWZpY2F0ZWV4cGly'
    'eRjyke4eIAEoCFIXaWdub3JlY2VydGlmaWNhdGVleHBpcnkSJQoMdHJ1c3RzdG9yZWlkGISQnC'
    'AgASgJUgx0cnVzdHN0b3JlaWQ=');

@$core.Deprecated('Use trustStoreSummaryDescriptor instead')
const TrustStoreSummary$json = {
  '1': 'TrustStoreSummary',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'numberofcacertificates',
      '3': 203958970,
      '4': 1,
      '5': 5,
      '10': 'numberofcacertificates'
    },
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.TrustStoreStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `TrustStoreSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trustStoreSummaryDescriptor = $convert.base64Decode(
    'ChFUcnVzdFN0b3JlU3VtbWFyeRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SFgoEZXRhZxiB37OVAS'
    'ABKAlSBGV0YWcSEgoCaWQYgfKitwEgASgJUgJpZBItChBsYXN0bW9kaWZpZWR0aW1lGOCC/HAg'
    'ASgJUhBsYXN0bW9kaWZpZWR0aW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSOQoWbnVtYmVyb2'
    'ZjYWNlcnRpZmljYXRlcxi61aBhIAEoBVIWbnVtYmVyb2ZjYWNlcnRpZmljYXRlcxIZCgZyZWFz'
    'b24YuoLFCSABKAlSBnJlYXNvbhI3CgZzdGF0dXMYkOT7AiABKA4yHC5jbG91ZGZyb250LlRydX'
    'N0U3RvcmVTdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use trustedKeyGroupDoesNotExistDescriptor instead')
const TrustedKeyGroupDoesNotExist$json = {
  '1': 'TrustedKeyGroupDoesNotExist',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrustedKeyGroupDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trustedKeyGroupDoesNotExistDescriptor =
    $convert.base64Decode(
        'ChtUcnVzdGVkS2V5R3JvdXBEb2VzTm90RXhpc3QSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use trustedKeyGroupsDescriptor instead')
const TrustedKeyGroups$json = {
  '1': 'TrustedKeyGroups',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `TrustedKeyGroups`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trustedKeyGroupsDescriptor = $convert.base64Decode(
    'ChBUcnVzdGVkS2V5R3JvdXBzEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVkEhcKBWl0ZW'
    '1zGLDw2AEgAygJUgVpdGVtcxIdCghxdWFudGl0eRj55dxfIAEoBVIIcXVhbnRpdHk=');

@$core.Deprecated('Use trustedSignerDoesNotExistDescriptor instead')
const TrustedSignerDoesNotExist$json = {
  '1': 'TrustedSignerDoesNotExist',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrustedSignerDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trustedSignerDoesNotExistDescriptor =
    $convert.base64Decode(
        'ChlUcnVzdGVkU2lnbmVyRG9lc05vdEV4aXN0EhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use trustedSignersDescriptor instead')
const TrustedSigners$json = {
  '1': 'TrustedSigners',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'items', '3': 3553328, '4': 3, '5': 9, '10': 'items'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `TrustedSigners`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trustedSignersDescriptor = $convert.base64Decode(
    'Cg5UcnVzdGVkU2lnbmVycxIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIXCgVpdGVtcx'
    'iw8NgBIAMoCVIFaXRlbXMSHQoIcXVhbnRpdHkY+eXcXyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use unsupportedOperationDescriptor instead')
const UnsupportedOperation$json = {
  '1': 'UnsupportedOperation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedOperation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedOperationDescriptor =
    $convert.base64Decode(
        'ChRVbnN1cHBvcnRlZE9wZXJhdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resource', '3': 61921302, '4': 1, '5': 9, '10': 'resource'},
    {
      '1': 'tagkeys',
      '3': 320659964,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TagKeys',
      '8': {},
      '10': 'tagkeys'
    },
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIdCghyZXNvdXJjZRiWsMMdIAEoCVIIcmVzb3VyY2USNw'
    'oHdGFna2V5cxj8w/OYASABKAsyEy5jbG91ZGZyb250LlRhZ0tleXNCBIi1GAFSB3RhZ2tleXM=');

@$core.Deprecated('Use updateAnycastIpListRequestDescriptor instead')
const UpdateAnycastIpListRequest$json = {
  '1': 'UpdateAnycastIpListRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'ipaddresstype',
      '3': 459110693,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.IpAddressType',
      '10': 'ipaddresstype'
    },
  ],
};

/// Descriptor for `UpdateAnycastIpListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAnycastIpListRequestDescriptor =
    $convert.base64Decode(
        'ChpVcGRhdGVBbnljYXN0SXBMaXN0UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEhsKB2lmbW'
        'F0Y2gY0Ja3LCABKAlSB2lmbWF0Y2gSQwoNaXBhZGRyZXNzdHlwZRil8vXaASABKA4yGS5jbG91'
        'ZGZyb250LklwQWRkcmVzc1R5cGVSDWlwYWRkcmVzc3R5cGU=');

@$core.Deprecated('Use updateAnycastIpListResultDescriptor instead')
const UpdateAnycastIpListResult$json = {
  '1': 'UpdateAnycastIpListResult',
  '2': [
    {
      '1': 'anycastiplist',
      '3': 190550768,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.AnycastIpList',
      '8': {},
      '10': 'anycastiplist'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateAnycastIpListResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAnycastIpListResultDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVBbnljYXN0SXBMaXN0UmVzdWx0EkgKDWFueWNhc3RpcGxpc3QY8KXuWiABKAsyGS'
    '5jbG91ZGZyb250LkFueWNhc3RJcExpc3RCBIi1GAFSDWFueWNhc3RpcGxpc3QSFgoEZXRhZxiB'
    '37OVASABKAlSBGV0YWc=');

@$core.Deprecated('Use updateCachePolicyRequestDescriptor instead')
const UpdateCachePolicyRequest$json = {
  '1': 'UpdateCachePolicyRequest',
  '2': [
    {
      '1': 'cachepolicyconfig',
      '3': 407094126,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicyConfig',
      '8': {},
      '10': 'cachepolicyconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateCachePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCachePolicyRequestDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVDYWNoZVBvbGljeVJlcXVlc3QSVQoRY2FjaGVwb2xpY3ljb25maWcY7oaPwgEgAS'
    'gLMh0uY2xvdWRmcm9udC5DYWNoZVBvbGljeUNvbmZpZ0IEiLUYAVIRY2FjaGVwb2xpY3ljb25m'
    'aWcSEgoCaWQYgfKitwEgASgJUgJpZBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNo');

@$core.Deprecated('Use updateCachePolicyResultDescriptor instead')
const UpdateCachePolicyResult$json = {
  '1': 'UpdateCachePolicyResult',
  '2': [
    {
      '1': 'cachepolicy',
      '3': 439848032,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CachePolicy',
      '8': {},
      '10': 'cachepolicy'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateCachePolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCachePolicyResultDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVDYWNoZVBvbGljeVJlc3VsdBJDCgtjYWNoZXBvbGljeRjgmN7RASABKAsyFy5jbG'
    '91ZGZyb250LkNhY2hlUG9saWN5QgSItRgBUgtjYWNoZXBvbGljeRIWCgRldGFnGIHfs5UBIAEo'
    'CVIEZXRhZw==');

@$core.Deprecated(
    'Use updateCloudFrontOriginAccessIdentityRequestDescriptor instead')
const UpdateCloudFrontOriginAccessIdentityRequest$json = {
  '1': 'UpdateCloudFrontOriginAccessIdentityRequest',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentityconfig',
      '3': 111945038,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentityConfig',
      '8': {},
      '10': 'cloudfrontoriginaccessidentityconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateCloudFrontOriginAccessIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateCloudFrontOriginAccessIdentityRequestDescriptor =
    $convert.base64Decode(
        'CitVcGRhdGVDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlSZXF1ZXN0Eo0BCiRjbG91ZG'
        'Zyb250b3JpZ2luYWNjZXNzaWRlbnRpdHljb25maWcYzsqwNSABKAsyMC5jbG91ZGZyb250LkNs'
        'b3VkRnJvbnRPcmlnaW5BY2Nlc3NJZGVudGl0eUNvbmZpZ0IEiLUYAVIkY2xvdWRmcm9udG9yaW'
        'dpbmFjY2Vzc2lkZW50aXR5Y29uZmlnEhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaBjQ'
        'lrcsIAEoCVIHaWZtYXRjaA==');

@$core.Deprecated(
    'Use updateCloudFrontOriginAccessIdentityResultDescriptor instead')
const UpdateCloudFrontOriginAccessIdentityResult$json = {
  '1': 'UpdateCloudFrontOriginAccessIdentityResult',
  '2': [
    {
      '1': 'cloudfrontoriginaccessidentity',
      '3': 109497984,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CloudFrontOriginAccessIdentity',
      '8': {},
      '10': 'cloudfrontoriginaccessidentity'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateCloudFrontOriginAccessIdentityResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateCloudFrontOriginAccessIdentityResultDescriptor =
    $convert.base64Decode(
        'CipVcGRhdGVDbG91ZEZyb250T3JpZ2luQWNjZXNzSWRlbnRpdHlSZXN1bHQSewoeY2xvdWRmcm'
        '9udG9yaWdpbmFjY2Vzc2lkZW50aXR5GICdmzQgASgLMiouY2xvdWRmcm9udC5DbG91ZEZyb250'
        'T3JpZ2luQWNjZXNzSWRlbnRpdHlCBIi1GAFSHmNsb3VkZnJvbnRvcmlnaW5hY2Nlc3NpZGVudG'
        'l0eRIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw==');

@$core.Deprecated('Use updateConnectionFunctionRequestDescriptor instead')
const UpdateConnectionFunctionRequest$json = {
  '1': 'UpdateConnectionFunctionRequest',
  '2': [
    {
      '1': 'connectionfunctioncode',
      '3': 502949501,
      '4': 1,
      '5': 12,
      '10': 'connectionfunctioncode'
    },
    {
      '1': 'connectionfunctionconfig',
      '3': 349557744,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionConfig',
      '10': 'connectionfunctionconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateConnectionFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionFunctionRequestDescriptor = $convert.base64Decode(
    'Ch9VcGRhdGVDb25uZWN0aW9uRnVuY3Rpb25SZXF1ZXN0EjoKFmNvbm5lY3Rpb25mdW5jdGlvbm'
    'NvZGUY/czp7wEgASgMUhZjb25uZWN0aW9uZnVuY3Rpb25jb2RlEloKGGNvbm5lY3Rpb25mdW5j'
    'dGlvbmNvbmZpZxjwp9emASABKAsyGi5jbG91ZGZyb250LkZ1bmN0aW9uQ29uZmlnUhhjb25uZW'
    'N0aW9uZnVuY3Rpb25jb25maWcSEgoCaWQYgfKitwEgASgJUgJpZBIbCgdpZm1hdGNoGNCWtywg'
    'ASgJUgdpZm1hdGNo');

@$core.Deprecated('Use updateConnectionFunctionResultDescriptor instead')
const UpdateConnectionFunctionResult$json = {
  '1': 'UpdateConnectionFunctionResult',
  '2': [
    {
      '1': 'connectionfunctionsummary',
      '3': 62528396,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionFunctionSummary',
      '8': {},
      '10': 'connectionfunctionsummary'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateConnectionFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionFunctionResultDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVDb25uZWN0aW9uRnVuY3Rpb25SZXN1bHQSbAoZY29ubmVjdGlvbmZ1bmN0aW9uc3'
        'VtbWFyeRiMt+gdIAEoCzIlLmNsb3VkZnJvbnQuQ29ubmVjdGlvbkZ1bmN0aW9uU3VtbWFyeUIE'
        'iLUYAVIZY29ubmVjdGlvbmZ1bmN0aW9uc3VtbWFyeRIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZw'
        '==');

@$core.Deprecated('Use updateConnectionGroupRequestDescriptor instead')
const UpdateConnectionGroupRequest$json = {
  '1': 'UpdateConnectionGroupRequest',
  '2': [
    {
      '1': 'anycastiplistid',
      '3': 431887875,
      '4': 1,
      '5': 9,
      '10': 'anycastiplistid'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'ipv6enabled', '3': 324073322, '4': 1, '5': 8, '10': 'ipv6enabled'},
  ],
};

/// Descriptor for `UpdateConnectionGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionGroupRequestDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVDb25uZWN0aW9uR3JvdXBSZXF1ZXN0EiwKD2FueWNhc3RpcGxpc3RpZBiDrPjNAS'
    'ABKAlSD2FueWNhc3RpcGxpc3RpZBIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBISCgJp'
    'ZBiB8qK3ASABKAlSAmlkEhsKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2gSJAoLaXB2NmVuYW'
    'JsZWQY6u7DmgEgASgIUgtpcHY2ZW5hYmxlZA==');

@$core.Deprecated('Use updateConnectionGroupResultDescriptor instead')
const UpdateConnectionGroupResult$json = {
  '1': 'UpdateConnectionGroupResult',
  '2': [
    {
      '1': 'connectiongroup',
      '3': 517217105,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ConnectionGroup',
      '8': {},
      '10': 'connectiongroup'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateConnectionGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConnectionGroupResultDescriptor =
    $convert.base64Decode(
        'ChtVcGRhdGVDb25uZWN0aW9uR3JvdXBSZXN1bHQSTwoPY29ubmVjdGlvbmdyb3VwGNG20PYBIA'
        'EoCzIbLmNsb3VkZnJvbnQuQ29ubmVjdGlvbkdyb3VwQgSItRgBUg9jb25uZWN0aW9uZ3JvdXAS'
        'FgoEZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core
    .Deprecated('Use updateContinuousDeploymentPolicyRequestDescriptor instead')
const UpdateContinuousDeploymentPolicyRequest$json = {
  '1': 'UpdateContinuousDeploymentPolicyRequest',
  '2': [
    {
      '1': 'continuousdeploymentpolicyconfig',
      '3': 161949042,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicyConfig',
      '8': {},
      '10': 'continuousdeploymentpolicyconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateContinuousDeploymentPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContinuousDeploymentPolicyRequestDescriptor =
    $convert.base64Decode(
        'CidVcGRhdGVDb250aW51b3VzRGVwbG95bWVudFBvbGljeVJlcXVlc3QSgQEKIGNvbnRpbnVvdX'
        'NkZXBsb3ltZW50cG9saWN5Y29uZmlnGPLKnE0gASgLMiwuY2xvdWRmcm9udC5Db250aW51b3Vz'
        'RGVwbG95bWVudFBvbGljeUNvbmZpZ0IEiLUYAVIgY29udGludW91c2RlcGxveW1lbnRwb2xpY3'
        'ljb25maWcSEgoCaWQYgfKitwEgASgJUgJpZBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNo');

@$core
    .Deprecated('Use updateContinuousDeploymentPolicyResultDescriptor instead')
const UpdateContinuousDeploymentPolicyResult$json = {
  '1': 'UpdateContinuousDeploymentPolicyResult',
  '2': [
    {
      '1': 'continuousdeploymentpolicy',
      '3': 36616788,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ContinuousDeploymentPolicy',
      '8': {},
      '10': 'continuousdeploymentpolicy'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateContinuousDeploymentPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContinuousDeploymentPolicyResultDescriptor =
    $convert.base64Decode(
        'CiZVcGRhdGVDb250aW51b3VzRGVwbG95bWVudFBvbGljeVJlc3VsdBJvChpjb250aW51b3VzZG'
        'VwbG95bWVudHBvbGljeRjU9LoRIAEoCzImLmNsb3VkZnJvbnQuQ29udGludW91c0RlcGxveW1l'
        'bnRQb2xpY3lCBIi1GAFSGmNvbnRpbnVvdXNkZXBsb3ltZW50cG9saWN5EhYKBGV0YWcYgd+zlQ'
        'EgASgJUgRldGFn');

@$core.Deprecated('Use updateDistributionRequestDescriptor instead')
const UpdateDistributionRequest$json = {
  '1': 'UpdateDistributionRequest',
  '2': [
    {
      '1': 'distributionconfig',
      '3': 528940762,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionConfig',
      '8': {},
      '10': 'distributionconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDistributionRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVEaXN0cmlidXRpb25SZXF1ZXN0ElgKEmRpc3RyaWJ1dGlvbmNvbmZpZxja/Zv8AS'
    'ABKAsyHi5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkNvbmZpZ0IEiLUYAVISZGlzdHJpYnV0aW9u'
    'Y29uZmlnEhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaA'
    '==');

@$core.Deprecated('Use updateDistributionResultDescriptor instead')
const UpdateDistributionResult$json = {
  '1': 'UpdateDistributionResult',
  '2': [
    {
      '1': 'distribution',
      '3': 105183308,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Distribution',
      '8': {},
      '10': 'distribution'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateDistributionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDistributionResultDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVEaXN0cmlidXRpb25SZXN1bHQSRQoMZGlzdHJpYnV0aW9uGMzwkzIgASgLMhguY2'
    'xvdWRmcm9udC5EaXN0cmlidXRpb25CBIi1GAFSDGRpc3RyaWJ1dGlvbhIWCgRldGFnGIHfs5UB'
    'IAEoCVIEZXRhZw==');

@$core.Deprecated('Use updateDistributionTenantRequestDescriptor instead')
const UpdateDistributionTenantRequest$json = {
  '1': 'UpdateDistributionTenantRequest',
  '2': [
    {
      '1': 'connectiongroupid',
      '3': 169532206,
      '4': 1,
      '5': 9,
      '10': 'connectiongroupid'
    },
    {
      '1': 'customizations',
      '3': 70755200,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Customizations',
      '10': 'customizations'
    },
    {
      '1': 'distributionid',
      '3': 142530791,
      '4': 1,
      '5': 9,
      '10': 'distributionid'
    },
    {
      '1': 'domains',
      '3': 149701959,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DomainItem',
      '10': 'domains'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'managedcertificaterequest',
      '3': 310259927,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ManagedCertificateRequest',
      '10': 'managedcertificaterequest'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.Parameter',
      '10': 'parameters'
    },
  ],
};

/// Descriptor for `UpdateDistributionTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDistributionTenantRequestDescriptor = $convert.base64Decode(
    'Ch9VcGRhdGVEaXN0cmlidXRpb25UZW5hbnRSZXF1ZXN0Ei8KEWNvbm5lY3Rpb25ncm91cGlkGK'
    '6261AgASgJUhFjb25uZWN0aW9uZ3JvdXBpZBJFCg5jdXN0b21pemF0aW9ucxiAx94hIAEoCzIa'
    'LmNsb3VkZnJvbnQuQ3VzdG9taXphdGlvbnNSDmN1c3RvbWl6YXRpb25zEikKDmRpc3RyaWJ1dG'
    'lvbmlkGOex+0MgASgJUg5kaXN0cmlidXRpb25pZBIzCgdkb21haW5zGMeKsUcgAygLMhYuY2xv'
    'dWRmcm9udC5Eb21haW5JdGVtUgdkb21haW5zEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibG'
    'VkEhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaBJnChlt'
    'YW5hZ2VkY2VydGlmaWNhdGVyZXF1ZXN0GNfh+JMBIAEoCzIlLmNsb3VkZnJvbnQuTWFuYWdlZE'
    'NlcnRpZmljYXRlUmVxdWVzdFIZbWFuYWdlZGNlcnRpZmljYXRlcmVxdWVzdBI5CgpwYXJhbWV0'
    'ZXJzGPqn/usBIAMoCzIVLmNsb3VkZnJvbnQuUGFyYW1ldGVyUgpwYXJhbWV0ZXJz');

@$core.Deprecated('Use updateDistributionTenantResultDescriptor instead')
const UpdateDistributionTenantResult$json = {
  '1': 'UpdateDistributionTenantResult',
  '2': [
    {
      '1': 'distributiontenant',
      '3': 510856916,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionTenant',
      '8': {},
      '10': 'distributiontenant'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateDistributionTenantResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDistributionTenantResultDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVEaXN0cmlidXRpb25UZW5hbnRSZXN1bHQSWAoSZGlzdHJpYnV0aW9udGVuYW50GN'
        'SdzPMBIAEoCzIeLmNsb3VkZnJvbnQuRGlzdHJpYnV0aW9uVGVuYW50QgSItRgBUhJkaXN0cmli'
        'dXRpb250ZW5hbnQSFgoEZXRhZxiB37OVASABKAlSBGV0YWc=');

@$core.Deprecated(
    'Use updateDistributionWithStagingConfigRequestDescriptor instead')
const UpdateDistributionWithStagingConfigRequest$json = {
  '1': 'UpdateDistributionWithStagingConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'stagingdistributionid',
      '3': 120990786,
      '4': 1,
      '5': 9,
      '10': 'stagingdistributionid'
    },
  ],
};

/// Descriptor for `UpdateDistributionWithStagingConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateDistributionWithStagingConfigRequestDescriptor =
    $convert.base64Decode(
        'CipVcGRhdGVEaXN0cmlidXRpb25XaXRoU3RhZ2luZ0NvbmZpZ1JlcXVlc3QSEgoCaWQYgfKitw'
        'EgASgJUgJpZBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNoEjcKFXN0YWdpbmdkaXN0cmli'
        'dXRpb25pZBjC2Ng5IAEoCVIVc3RhZ2luZ2Rpc3RyaWJ1dGlvbmlk');

@$core.Deprecated(
    'Use updateDistributionWithStagingConfigResultDescriptor instead')
const UpdateDistributionWithStagingConfigResult$json = {
  '1': 'UpdateDistributionWithStagingConfigResult',
  '2': [
    {
      '1': 'distribution',
      '3': 105183308,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.Distribution',
      '8': {},
      '10': 'distribution'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
  ],
};

/// Descriptor for `UpdateDistributionWithStagingConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateDistributionWithStagingConfigResultDescriptor = $convert.base64Decode(
        'CilVcGRhdGVEaXN0cmlidXRpb25XaXRoU3RhZ2luZ0NvbmZpZ1Jlc3VsdBJFCgxkaXN0cmlidX'
        'Rpb24YzPCTMiABKAsyGC5jbG91ZGZyb250LkRpc3RyaWJ1dGlvbkIEiLUYAVIMZGlzdHJpYnV0'
        'aW9uEhYKBGV0YWcYgd+zlQEgASgJUgRldGFn');

@$core.Deprecated('Use updateDomainAssociationRequestDescriptor instead')
const UpdateDomainAssociationRequest$json = {
  '1': 'UpdateDomainAssociationRequest',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'targetresource',
      '3': 523474061,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.DistributionResourceId',
      '10': 'targetresource'
    },
  ],
};

/// Descriptor for `UpdateDomainAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDomainAssociationRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVEb21haW5Bc3NvY2lhdGlvblJlcXVlc3QSGgoGZG9tYWluGJKS8vABIAEoCVIGZG'
        '9tYWluEhsKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2gSTgoOdGFyZ2V0cmVzb3VyY2UYjanO'
        '+QEgASgLMiIuY2xvdWRmcm9udC5EaXN0cmlidXRpb25SZXNvdXJjZUlkUg50YXJnZXRyZXNvdX'
        'JjZQ==');

@$core.Deprecated('Use updateDomainAssociationResultDescriptor instead')
const UpdateDomainAssociationResult$json = {
  '1': 'UpdateDomainAssociationResult',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
  ],
};

/// Descriptor for `UpdateDomainAssociationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDomainAssociationResultDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVEb21haW5Bc3NvY2lhdGlvblJlc3VsdBIaCgZkb21haW4YkpLy8AEgASgJUgZkb2'
        '1haW4SFgoEZXRhZxiB37OVASABKAlSBGV0YWcSIgoKcmVzb3VyY2VpZBiRuvH6ASABKAlSCnJl'
        'c291cmNlaWQ=');

@$core
    .Deprecated('Use updateFieldLevelEncryptionConfigRequestDescriptor instead')
const UpdateFieldLevelEncryptionConfigRequest$json = {
  '1': 'UpdateFieldLevelEncryptionConfigRequest',
  '2': [
    {
      '1': 'fieldlevelencryptionconfig',
      '3': 499294709,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionConfig',
      '8': {},
      '10': 'fieldlevelencryptionconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateFieldLevelEncryptionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFieldLevelEncryptionConfigRequestDescriptor =
    $convert.base64Decode(
        'CidVcGRhdGVGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ1JlcXVlc3QScAoaZmllbGRsZXZlbG'
        'VuY3J5cHRpb25jb25maWcY9cOK7gEgASgLMiYuY2xvdWRmcm9udC5GaWVsZExldmVsRW5jcnlw'
        'dGlvbkNvbmZpZ0IEiLUYAVIaZmllbGRsZXZlbGVuY3J5cHRpb25jb25maWcSEgoCaWQYgfKitw'
        'EgASgJUgJpZBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNo');

@$core
    .Deprecated('Use updateFieldLevelEncryptionConfigResultDescriptor instead')
const UpdateFieldLevelEncryptionConfigResult$json = {
  '1': 'UpdateFieldLevelEncryptionConfigResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryption',
      '3': 473382747,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryption',
      '8': {},
      '10': 'fieldlevelencryption'
    },
  ],
};

/// Descriptor for `UpdateFieldLevelEncryptionConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFieldLevelEncryptionConfigResultDescriptor =
    $convert.base64Decode(
        'CiZVcGRhdGVGaWVsZExldmVsRW5jcnlwdGlvbkNvbmZpZ1Jlc3VsdBIWCgRldGFnGIHfs5UBIA'
        'EoCVIEZXRhZxJeChRmaWVsZGxldmVsZW5jcnlwdGlvbhjb/tzhASABKAsyIC5jbG91ZGZyb250'
        'LkZpZWxkTGV2ZWxFbmNyeXB0aW9uQgSItRgBUhRmaWVsZGxldmVsZW5jcnlwdGlvbg==');

@$core.Deprecated(
    'Use updateFieldLevelEncryptionProfileRequestDescriptor instead')
const UpdateFieldLevelEncryptionProfileRequest$json = {
  '1': 'UpdateFieldLevelEncryptionProfileRequest',
  '2': [
    {
      '1': 'fieldlevelencryptionprofileconfig',
      '3': 199371734,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfileConfig',
      '8': {},
      '10': 'fieldlevelencryptionprofileconfig'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateFieldLevelEncryptionProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFieldLevelEncryptionProfileRequestDescriptor =
    $convert.base64Decode(
        'CihVcGRhdGVGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVSZXF1ZXN0EoQBCiFmaWVsZGxldm'
        'VsZW5jcnlwdGlvbnByb2ZpbGVjb25maWcY1teIXyABKAsyLS5jbG91ZGZyb250LkZpZWxkTGV2'
        'ZWxFbmNyeXB0aW9uUHJvZmlsZUNvbmZpZ0IEiLUYAVIhZmllbGRsZXZlbGVuY3J5cHRpb25wcm'
        '9maWxlY29uZmlnEhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZt'
        'YXRjaA==');

@$core
    .Deprecated('Use updateFieldLevelEncryptionProfileResultDescriptor instead')
const UpdateFieldLevelEncryptionProfileResult$json = {
  '1': 'UpdateFieldLevelEncryptionProfileResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'fieldlevelencryptionprofile',
      '3': 344546136,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FieldLevelEncryptionProfile',
      '8': {},
      '10': 'fieldlevelencryptionprofile'
    },
  ],
};

/// Descriptor for `UpdateFieldLevelEncryptionProfileResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFieldLevelEncryptionProfileResultDescriptor =
    $convert.base64Decode(
        'CidVcGRhdGVGaWVsZExldmVsRW5jcnlwdGlvblByb2ZpbGVSZXN1bHQSFgoEZXRhZxiB37OVAS'
        'ABKAlSBGV0YWcScwobZmllbGRsZXZlbGVuY3J5cHRpb25wcm9maWxlGNi2paQBIAEoCzInLmNs'
        'b3VkZnJvbnQuRmllbGRMZXZlbEVuY3J5cHRpb25Qcm9maWxlQgSItRgBUhtmaWVsZGxldmVsZW'
        '5jcnlwdGlvbnByb2ZpbGU=');

@$core.Deprecated('Use updateFunctionRequestDescriptor instead')
const UpdateFunctionRequest$json = {
  '1': 'UpdateFunctionRequest',
  '2': [
    {
      '1': 'functioncode',
      '3': 405947809,
      '4': 1,
      '5': 12,
      '10': 'functioncode'
    },
    {
      '1': 'functionconfig',
      '3': 116111484,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionConfig',
      '10': 'functionconfig'
    },
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFunctionRequestDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVGdW5jdGlvblJlcXVlc3QSJgoMZnVuY3Rpb25jb2RlGKGLycEBIAEoDFIMZnVuY3'
    'Rpb25jb2RlEkUKDmZ1bmN0aW9uY29uZmlnGPzwrjcgASgLMhouY2xvdWRmcm9udC5GdW5jdGlv'
    'bkNvbmZpZ1IOZnVuY3Rpb25jb25maWcSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaBIVCg'
    'RuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use updateFunctionResultDescriptor instead')
const UpdateFunctionResult$json = {
  '1': 'UpdateFunctionResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'functionsummary',
      '3': 523316264,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.FunctionSummary',
      '8': {},
      '10': 'functionsummary'
    },
  ],
};

/// Descriptor for `UpdateFunctionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFunctionResultDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVGdW5jdGlvblJlc3VsdBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxJPCg9mdW5jdG'
    'lvbnN1bW1hcnkYqNjE+QEgASgLMhsuY2xvdWRmcm9udC5GdW5jdGlvblN1bW1hcnlCBIi1GAFS'
    'D2Z1bmN0aW9uc3VtbWFyeQ==');

@$core.Deprecated('Use updateKeyGroupRequestDescriptor instead')
const UpdateKeyGroupRequest$json = {
  '1': 'UpdateKeyGroupRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'keygroupconfig',
      '3': 143012494,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroupConfig',
      '8': {},
      '10': 'keygroupconfig'
    },
  ],
};

/// Descriptor for `UpdateKeyGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKeyGroupRequestDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVLZXlHcm91cFJlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZBIbCgdpZm1hdGNoGN'
    'CWtywgASgJUgdpZm1hdGNoEksKDmtleWdyb3VwY29uZmlnGI7lmEQgASgLMhouY2xvdWRmcm9u'
    'dC5LZXlHcm91cENvbmZpZ0IEiLUYAVIOa2V5Z3JvdXBjb25maWc=');

@$core.Deprecated('Use updateKeyGroupResultDescriptor instead')
const UpdateKeyGroupResult$json = {
  '1': 'UpdateKeyGroupResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'keygroup',
      '3': 518748096,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyGroup',
      '8': {},
      '10': 'keygroup'
    },
  ],
};

/// Descriptor for `UpdateKeyGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKeyGroupResultDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVLZXlHcm91cFJlc3VsdBIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxI6CghrZXlncm'
    '91cBjA7633ASABKAsyFC5jbG91ZGZyb250LktleUdyb3VwQgSItRgBUghrZXlncm91cA==');

@$core.Deprecated('Use updateKeyValueStoreRequestDescriptor instead')
const UpdateKeyValueStoreRequest$json = {
  '1': 'UpdateKeyValueStoreRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateKeyValueStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKeyValueStoreRequestDescriptor =
    $convert.base64Decode(
        'ChpVcGRhdGVLZXlWYWx1ZVN0b3JlUmVxdWVzdBIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbW'
        'VudBIbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNoEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use updateKeyValueStoreResultDescriptor instead')
const UpdateKeyValueStoreResult$json = {
  '1': 'UpdateKeyValueStoreResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'keyvaluestore',
      '3': 151113103,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.KeyValueStore',
      '8': {},
      '10': 'keyvaluestore'
    },
  ],
};

/// Descriptor for `UpdateKeyValueStoreResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKeyValueStoreResultDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVLZXlWYWx1ZVN0b3JlUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEkgKDW'
    'tleXZhbHVlc3RvcmUYj5uHSCABKAsyGS5jbG91ZGZyb250LktleVZhbHVlU3RvcmVCBIi1GAFS'
    'DWtleXZhbHVlc3RvcmU=');

@$core.Deprecated('Use updateOriginAccessControlRequestDescriptor instead')
const UpdateOriginAccessControlRequest$json = {
  '1': 'UpdateOriginAccessControlRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'originaccesscontrolconfig',
      '3': 143834977,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControlConfig',
      '8': {},
      '10': 'originaccesscontrolconfig'
    },
  ],
};

/// Descriptor for `UpdateOriginAccessControlRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOriginAccessControlRequestDescriptor =
    $convert.base64Decode(
        'CiBVcGRhdGVPcmlnaW5BY2Nlc3NDb250cm9sUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEh'
        'sKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2gSbAoZb3JpZ2luYWNjZXNzY29udHJvbGNvbmZp'
        'Zxjh/spEIAEoCzIlLmNsb3VkZnJvbnQuT3JpZ2luQWNjZXNzQ29udHJvbENvbmZpZ0IEiLUYAV'
        'IZb3JpZ2luYWNjZXNzY29udHJvbGNvbmZpZw==');

@$core.Deprecated('Use updateOriginAccessControlResultDescriptor instead')
const UpdateOriginAccessControlResult$json = {
  '1': 'UpdateOriginAccessControlResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'originaccesscontrol',
      '3': 238302375,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginAccessControl',
      '8': {},
      '10': 'originaccesscontrol'
    },
  ],
};

/// Descriptor for `UpdateOriginAccessControlResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOriginAccessControlResultDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVPcmlnaW5BY2Nlc3NDb250cm9sUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldG'
        'FnEloKE29yaWdpbmFjY2Vzc2NvbnRyb2wYp+nQcSABKAsyHy5jbG91ZGZyb250Lk9yaWdpbkFj'
        'Y2Vzc0NvbnRyb2xCBIi1GAFSE29yaWdpbmFjY2Vzc2NvbnRyb2w=');

@$core.Deprecated('Use updateOriginRequestPolicyRequestDescriptor instead')
const UpdateOriginRequestPolicyRequest$json = {
  '1': 'UpdateOriginRequestPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'originrequestpolicyconfig',
      '3': 37078133,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicyConfig',
      '8': {},
      '10': 'originrequestpolicyconfig'
    },
  ],
};

/// Descriptor for `UpdateOriginRequestPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOriginRequestPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiBVcGRhdGVPcmlnaW5SZXF1ZXN0UG9saWN5UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEh'
        'sKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2gSbAoZb3JpZ2lucmVxdWVzdHBvbGljeWNvbmZp'
        'Zxj1iNcRIAEoCzIlLmNsb3VkZnJvbnQuT3JpZ2luUmVxdWVzdFBvbGljeUNvbmZpZ0IEiLUYAV'
        'IZb3JpZ2lucmVxdWVzdHBvbGljeWNvbmZpZw==');

@$core.Deprecated('Use updateOriginRequestPolicyResultDescriptor instead')
const UpdateOriginRequestPolicyResult$json = {
  '1': 'UpdateOriginRequestPolicyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'originrequestpolicy',
      '3': 386733531,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginRequestPolicy',
      '8': {},
      '10': 'originrequestpolicy'
    },
  ],
};

/// Descriptor for `UpdateOriginRequestPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOriginRequestPolicyResultDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVPcmlnaW5SZXF1ZXN0UG9saWN5UmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldG'
        'FnElsKE29yaWdpbnJlcXVlc3Rwb2xpY3kY26u0uAEgASgLMh8uY2xvdWRmcm9udC5PcmlnaW5S'
        'ZXF1ZXN0UG9saWN5QgSItRgBUhNvcmlnaW5yZXF1ZXN0cG9saWN5');

@$core.Deprecated('Use updatePublicKeyRequestDescriptor instead')
const UpdatePublicKeyRequest$json = {
  '1': 'UpdatePublicKeyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'publickeyconfig',
      '3': 228537966,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKeyConfig',
      '8': {},
      '10': 'publickeyconfig'
    },
  ],
};

/// Descriptor for `UpdatePublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePublicKeyRequestDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVQdWJsaWNLZXlSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaB'
    'jQlrcsIAEoCVIHaWZtYXRjaBJOCg9wdWJsaWNrZXljb25maWcY7uz8bCABKAsyGy5jbG91ZGZy'
    'b250LlB1YmxpY0tleUNvbmZpZ0IEiLUYAVIPcHVibGlja2V5Y29uZmln');

@$core.Deprecated('Use updatePublicKeyResultDescriptor instead')
const UpdatePublicKeyResult$json = {
  '1': 'UpdatePublicKeyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'publickey',
      '3': 167335776,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.PublicKey',
      '8': {},
      '10': 'publickey'
    },
  ],
};

/// Descriptor for `UpdatePublicKeyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePublicKeyResultDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVQdWJsaWNLZXlSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSPAoJcHVibG'
    'lja2V5GOCu5U8gASgLMhUuY2xvdWRmcm9udC5QdWJsaWNLZXlCBIi1GAFSCXB1YmxpY2tleQ==');

@$core.Deprecated('Use updateRealtimeLogConfigRequestDescriptor instead')
const UpdateRealtimeLogConfigRequest$json = {
  '1': 'UpdateRealtimeLogConfigRequest',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'endpoints',
      '3': 436023390,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.EndPoint',
      '10': 'endpoints'
    },
    {'1': 'fields', '3': 319339933, '4': 3, '5': 9, '10': 'fields'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'samplingrate', '3': 272929747, '4': 1, '5': 3, '10': 'samplingrate'},
  ],
};

/// Descriptor for `UpdateRealtimeLogConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRealtimeLogConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVSZWFsdGltZUxvZ0NvbmZpZ1JlcXVlc3QSFAoDYXJuGJ2cr70BIAEoCVIDYXJuEj'
        'YKCWVuZHBvaW50cxje4PTPASADKAsyFC5jbG91ZGZyb250LkVuZFBvaW50UgllbmRwb2ludHMS'
        'GgoGZmllbGRzGJ37opgBIAMoCVIGZmllbGRzEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJgoMc2'
        'FtcGxpbmdyYXRlGNOnkoIBIAEoA1IMc2FtcGxpbmdyYXRl');

@$core.Deprecated('Use updateRealtimeLogConfigResultDescriptor instead')
const UpdateRealtimeLogConfigResult$json = {
  '1': 'UpdateRealtimeLogConfigResult',
  '2': [
    {
      '1': 'realtimelogconfig',
      '3': 95917609,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.RealtimeLogConfig',
      '10': 'realtimelogconfig'
    },
  ],
};

/// Descriptor for `UpdateRealtimeLogConfigResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRealtimeLogConfigResultDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVSZWFsdGltZUxvZ0NvbmZpZ1Jlc3VsdBJOChFyZWFsdGltZWxvZ2NvbmZpZxiprN'
        '4tIAEoCzIdLmNsb3VkZnJvbnQuUmVhbHRpbWVMb2dDb25maWdSEXJlYWx0aW1lbG9nY29uZmln');

@$core.Deprecated('Use updateResponseHeadersPolicyRequestDescriptor instead')
const UpdateResponseHeadersPolicyRequest$json = {
  '1': 'UpdateResponseHeadersPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'responseheaderspolicyconfig',
      '3': 159056825,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicyConfig',
      '8': {},
      '10': 'responseheaderspolicyconfig'
    },
  ],
};

/// Descriptor for `UpdateResponseHeadersPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResponseHeadersPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiJVcGRhdGVSZXNwb25zZUhlYWRlcnNQb2xpY3lSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'QSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaBJyChtyZXNwb25zZWhlYWRlcnNwb2xpY3lj'
        'b25maWcYuYfsSyABKAsyJy5jbG91ZGZyb250LlJlc3BvbnNlSGVhZGVyc1BvbGljeUNvbmZpZ0'
        'IEiLUYAVIbcmVzcG9uc2VoZWFkZXJzcG9saWN5Y29uZmln');

@$core.Deprecated('Use updateResponseHeadersPolicyResultDescriptor instead')
const UpdateResponseHeadersPolicyResult$json = {
  '1': 'UpdateResponseHeadersPolicyResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'responseheaderspolicy',
      '3': 418204719,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.ResponseHeadersPolicy',
      '8': {},
      '10': 'responseheaderspolicy'
    },
  ],
};

/// Descriptor for `UpdateResponseHeadersPolicyResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResponseHeadersPolicyResultDescriptor =
    $convert.base64Decode(
        'CiFVcGRhdGVSZXNwb25zZUhlYWRlcnNQb2xpY3lSZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBG'
        'V0YWcSYQoVcmVzcG9uc2VoZWFkZXJzcG9saWN5GK+YtccBIAEoCzIhLmNsb3VkZnJvbnQuUmVz'
        'cG9uc2VIZWFkZXJzUG9saWN5QgSItRgBUhVyZXNwb25zZWhlYWRlcnNwb2xpY3k=');

@$core.Deprecated('Use updateStreamingDistributionRequestDescriptor instead')
const UpdateStreamingDistributionRequest$json = {
  '1': 'UpdateStreamingDistributionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'streamingdistributionconfig',
      '3': 291115944,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistributionConfig',
      '8': {},
      '10': 'streamingdistributionconfig'
    },
  ],
};

/// Descriptor for `UpdateStreamingDistributionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStreamingDistributionRequestDescriptor =
    $convert.base64Decode(
        'CiJVcGRhdGVTdHJlYW1pbmdEaXN0cmlidXRpb25SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'QSGwoHaWZtYXRjaBjQlrcsIAEoCVIHaWZtYXRjaBJzChtzdHJlYW1pbmdkaXN0cmlidXRpb25j'
        'b25maWcYqKfoigEgASgLMicuY2xvdWRmcm9udC5TdHJlYW1pbmdEaXN0cmlidXRpb25Db25maW'
        'dCBIi1GAFSG3N0cmVhbWluZ2Rpc3RyaWJ1dGlvbmNvbmZpZw==');

@$core.Deprecated('Use updateStreamingDistributionResultDescriptor instead')
const UpdateStreamingDistributionResult$json = {
  '1': 'UpdateStreamingDistributionResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'streamingdistribution',
      '3': 294813830,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.StreamingDistribution',
      '8': {},
      '10': 'streamingdistribution'
    },
  ],
};

/// Descriptor for `UpdateStreamingDistributionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStreamingDistributionResultDescriptor =
    $convert.base64Decode(
        'CiFVcGRhdGVTdHJlYW1pbmdEaXN0cmlidXRpb25SZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBG'
        'V0YWcSYQoVc3RyZWFtaW5nZGlzdHJpYnV0aW9uGIaByowBIAEoCzIhLmNsb3VkZnJvbnQuU3Ry'
        'ZWFtaW5nRGlzdHJpYnV0aW9uQgSItRgBUhVzdHJlYW1pbmdkaXN0cmlidXRpb24=');

@$core.Deprecated('Use updateTrustStoreRequestDescriptor instead')
const UpdateTrustStoreRequest$json = {
  '1': 'UpdateTrustStoreRequest',
  '2': [
    {
      '1': 'cacertificatesbundlesource',
      '3': 328448937,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.CaCertificatesBundleSource',
      '8': {},
      '10': 'cacertificatesbundlesource'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
  ],
};

/// Descriptor for `UpdateTrustStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrustStoreRequestDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVUcnVzdFN0b3JlUmVxdWVzdBJwChpjYWNlcnRpZmljYXRlc2J1bmRsZXNvdXJjZR'
    'ip986cASABKAsyJi5jbG91ZGZyb250LkNhQ2VydGlmaWNhdGVzQnVuZGxlU291cmNlQgSItRgB'
    'UhpjYWNlcnRpZmljYXRlc2J1bmRsZXNvdXJjZRISCgJpZBiB8qK3ASABKAlSAmlkEhsKB2lmbW'
    'F0Y2gY0Ja3LCABKAlSB2lmbWF0Y2g=');

@$core.Deprecated('Use updateTrustStoreResultDescriptor instead')
const UpdateTrustStoreResult$json = {
  '1': 'UpdateTrustStoreResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'truststore',
      '3': 224815327,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustStore',
      '8': {},
      '10': 'truststore'
    },
  ],
};

/// Descriptor for `UpdateTrustStoreResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrustStoreResultDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVUcnVzdFN0b3JlUmVzdWx0EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEj8KCnRydX'
    'N0c3RvcmUY39GZayABKAsyFi5jbG91ZGZyb250LlRydXN0U3RvcmVCBIi1GAFSCnRydXN0c3Rv'
    'cmU=');

@$core.Deprecated('Use updateVpcOriginRequestDescriptor instead')
const UpdateVpcOriginRequest$json = {
  '1': 'UpdateVpcOriginRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'vpcoriginendpointconfig',
      '3': 45660030,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOriginEndpointConfig',
      '8': {},
      '10': 'vpcoriginendpointconfig'
    },
  ],
};

/// Descriptor for `UpdateVpcOriginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateVpcOriginRequestDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVWcGNPcmlnaW5SZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQSGwoHaWZtYXRjaB'
    'jQlrcsIAEoCVIHaWZtYXRjaBJmChd2cGNvcmlnaW5lbmRwb2ludGNvbmZpZxj+7uIVIAEoCzIj'
    'LmNsb3VkZnJvbnQuVnBjT3JpZ2luRW5kcG9pbnRDb25maWdCBIi1GAFSF3ZwY29yaWdpbmVuZH'
    'BvaW50Y29uZmln');

@$core.Deprecated('Use updateVpcOriginResultDescriptor instead')
const UpdateVpcOriginResult$json = {
  '1': 'UpdateVpcOriginResult',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'vpcorigin',
      '3': 159181387,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOrigin',
      '8': {},
      '10': 'vpcorigin'
    },
  ],
};

/// Descriptor for `UpdateVpcOriginResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateVpcOriginResultDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVWcGNPcmlnaW5SZXN1bHQSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSPAoJdnBjb3'
    'JpZ2luGMvU80sgASgLMhUuY2xvdWRmcm9udC5WcGNPcmlnaW5CBIi1GAFSCXZwY29yaWdpbg==');

@$core.Deprecated('Use validationTokenDetailDescriptor instead')
const ValidationTokenDetail$json = {
  '1': 'ValidationTokenDetail',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'redirectfrom', '3': 52853514, '4': 1, '5': 9, '10': 'redirectfrom'},
    {'1': 'redirectto', '3': 472243857, '4': 1, '5': 9, '10': 'redirectto'},
  ],
};

/// Descriptor for `ValidationTokenDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationTokenDetailDescriptor = $convert.base64Decode(
    'ChVWYWxpZGF0aW9uVG9rZW5EZXRhaWwSGgoGZG9tYWluGJKS8vABIAEoCVIGZG9tYWluEiUKDH'
    'JlZGlyZWN0ZnJvbRiK9pkZIAEoCVIMcmVkaXJlY3Rmcm9tEiIKCnJlZGlyZWN0dG8Ykb2X4QEg'
    'ASgJUgpyZWRpcmVjdHRv');

@$core.Deprecated('Use verifyDnsConfigurationRequestDescriptor instead')
const VerifyDnsConfigurationRequest$json = {
  '1': 'VerifyDnsConfigurationRequest',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `VerifyDnsConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyDnsConfigurationRequestDescriptor =
    $convert.base64Decode(
        'Ch1WZXJpZnlEbnNDb25maWd1cmF0aW9uUmVxdWVzdBIaCgZkb21haW4YkpLy8AEgASgJUgZkb2'
        '1haW4SIQoKaWRlbnRpZmllchjfoPsTIAEoCVIKaWRlbnRpZmllcg==');

@$core.Deprecated('Use verifyDnsConfigurationResultDescriptor instead')
const VerifyDnsConfigurationResult$json = {
  '1': 'VerifyDnsConfigurationResult',
  '2': [
    {
      '1': 'dnsconfigurationlist',
      '3': 355373827,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.DnsConfiguration',
      '10': 'dnsconfigurationlist'
    },
  ],
};

/// Descriptor for `VerifyDnsConfigurationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyDnsConfigurationResultDescriptor =
    $convert.base64Decode(
        'ChxWZXJpZnlEbnNDb25maWd1cmF0aW9uUmVzdWx0ElQKFGRuc2NvbmZpZ3VyYXRpb25saXN0GI'
        'OmuqkBIAMoCzIcLmNsb3VkZnJvbnQuRG5zQ29uZmlndXJhdGlvblIUZG5zY29uZmlndXJhdGlv'
        'bmxpc3Q=');

@$core.Deprecated('Use viewerCertificateDescriptor instead')
const ViewerCertificate$json = {
  '1': 'ViewerCertificate',
  '2': [
    {
      '1': 'acmcertificatearn',
      '3': 294529483,
      '4': 1,
      '5': 9,
      '10': 'acmcertificatearn'
    },
    {'1': 'certificate', '3': 198060817, '4': 1, '5': 9, '10': 'certificate'},
    {
      '1': 'certificatesource',
      '3': 63765974,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CertificateSource',
      '10': 'certificatesource'
    },
    {
      '1': 'cloudfrontdefaultcertificate',
      '3': 529306422,
      '4': 1,
      '5': 8,
      '10': 'cloudfrontdefaultcertificate'
    },
    {
      '1': 'iamcertificateid',
      '3': 92232821,
      '4': 1,
      '5': 9,
      '10': 'iamcertificateid'
    },
    {
      '1': 'minimumprotocolversion',
      '3': 367416622,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.MinimumProtocolVersion',
      '10': 'minimumprotocolversion'
    },
    {
      '1': 'sslsupportmethod',
      '3': 441557986,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.SSLSupportMethod',
      '10': 'sslsupportmethod'
    },
  ],
};

/// Descriptor for `ViewerCertificate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List viewerCertificateDescriptor = $convert.base64Decode(
    'ChFWaWV3ZXJDZXJ0aWZpY2F0ZRIwChFhY21jZXJ0aWZpY2F0ZWFybhjL07iMASABKAlSEWFjbW'
    'NlcnRpZmljYXRlYXJuEiMKC2NlcnRpZmljYXRlGJHWuF4gASgJUgtjZXJ0aWZpY2F0ZRJOChFj'
    'ZXJ0aWZpY2F0ZXNvdXJjZRjW+7MeIAEoDjIdLmNsb3VkZnJvbnQuQ2VydGlmaWNhdGVTb3VyY2'
    'VSEWNlcnRpZmljYXRlc291cmNlEkYKHGNsb3VkZnJvbnRkZWZhdWx0Y2VydGlmaWNhdGUYtqay'
    '/AEgASgIUhxjbG91ZGZyb250ZGVmYXVsdGNlcnRpZmljYXRlEi0KEGlhbWNlcnRpZmljYXRlaW'
    'QY9bj9KyABKAlSEGlhbWNlcnRpZmljYXRlaWQSXgoWbWluaW11bXByb3RvY29sdmVyc2lvbhiu'
    'qpmvASABKA4yIi5jbG91ZGZyb250Lk1pbmltdW1Qcm90b2NvbFZlcnNpb25SFm1pbmltdW1wcm'
    '90b2NvbHZlcnNpb24STAoQc3Nsc3VwcG9ydG1ldGhvZBjix8bSASABKA4yHC5jbG91ZGZyb250'
    'LlNTTFN1cHBvcnRNZXRob2RSEHNzbHN1cHBvcnRtZXRob2Q=');

@$core.Deprecated('Use viewerMtlsConfigDescriptor instead')
const ViewerMtlsConfig$json = {
  '1': 'ViewerMtlsConfig',
  '2': [
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.ViewerMtlsMode',
      '10': 'mode'
    },
    {
      '1': 'truststoreconfig',
      '3': 397230249,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.TrustStoreConfig',
      '10': 'truststoreconfig'
    },
  ],
};

/// Descriptor for `ViewerMtlsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List viewerMtlsConfigDescriptor = $convert.base64Decode(
    'ChBWaWV3ZXJNdGxzQ29uZmlnEjIKBG1vZGUYs+65mgEgASgOMhouY2xvdWRmcm9udC5WaWV3ZX'
    'JNdGxzTW9kZVIEbW9kZRJMChB0cnVzdHN0b3JlY29uZmlnGKmBtb0BIAEoCzIcLmNsb3VkZnJv'
    'bnQuVHJ1c3RTdG9yZUNvbmZpZ1IQdHJ1c3RzdG9yZWNvbmZpZw==');

@$core.Deprecated('Use vpcOriginDescriptor instead')
const VpcOrigin$json = {
  '1': 'VpcOrigin',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'vpcoriginendpointconfig',
      '3': 45660030,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.VpcOriginEndpointConfig',
      '10': 'vpcoriginendpointconfig'
    },
  ],
};

/// Descriptor for `VpcOrigin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcOriginDescriptor = $convert.base64Decode(
    'CglWcGNPcmlnaW4SHwoJYWNjb3VudGlkGNLBuR8gASgJUglhY2NvdW50aWQSFAoDYXJuGJ2b7b'
    '8BIAEoCVIDYXJuEiMKC2NyZWF0ZWR0aW1lGPPr8zkgASgJUgtjcmVhdGVkdGltZRISCgJpZBiB'
    '8qK3ASABKAlSAmlkEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZH'
    'RpbWUSGQoGc3RhdHVzGJDk+wIgASgJUgZzdGF0dXMSYAoXdnBjb3JpZ2luZW5kcG9pbnRjb25m'
    'aWcY/u7iFSABKAsyIy5jbG91ZGZyb250LlZwY09yaWdpbkVuZHBvaW50Q29uZmlnUhd2cGNvcm'
    'lnaW5lbmRwb2ludGNvbmZpZw==');

@$core.Deprecated('Use vpcOriginConfigDescriptor instead')
const VpcOriginConfig$json = {
  '1': 'VpcOriginConfig',
  '2': [
    {
      '1': 'originkeepalivetimeout',
      '3': 214128603,
      '4': 1,
      '5': 5,
      '10': 'originkeepalivetimeout'
    },
    {
      '1': 'originreadtimeout',
      '3': 387717023,
      '4': 1,
      '5': 5,
      '10': 'originreadtimeout'
    },
    {
      '1': 'owneraccountid',
      '3': 369721751,
      '4': 1,
      '5': 9,
      '10': 'owneraccountid'
    },
    {'1': 'vpcoriginid', '3': 365404648, '4': 1, '5': 9, '10': 'vpcoriginid'},
  ],
};

/// Descriptor for `VpcOriginConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcOriginConfigDescriptor = $convert.base64Decode(
    'Cg9WcGNPcmlnaW5Db25maWcSOQoWb3JpZ2lua2VlcGFsaXZldGltZW91dBjbr41mIAEoBVIWb3'
    'JpZ2lua2VlcGFsaXZldGltZW91dBIwChFvcmlnaW5yZWFkdGltZW91dBifr/C4ASABKAVSEW9y'
    'aWdpbnJlYWR0aW1lb3V0EioKDm93bmVyYWNjb3VudGlkGJeDprABIAEoCVIOb3duZXJhY2NvdW'
    '50aWQSJAoLdnBjb3JpZ2luaWQY6MOergEgASgJUgt2cGNvcmlnaW5pZA==');

@$core.Deprecated('Use vpcOriginEndpointConfigDescriptor instead')
const VpcOriginEndpointConfig$json = {
  '1': 'VpcOriginEndpointConfig',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'httpport', '3': 375976295, '4': 1, '5': 5, '10': 'httpport'},
    {'1': 'httpsport', '3': 153043918, '4': 1, '5': 5, '10': 'httpsport'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'originprotocolpolicy',
      '3': 234303342,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.OriginProtocolPolicy',
      '10': 'originprotocolpolicy'
    },
    {
      '1': 'originsslprotocols',
      '3': 403853493,
      '4': 1,
      '5': 11,
      '6': '.cloudfront.OriginSslProtocols',
      '10': 'originsslprotocols'
    },
  ],
};

/// Descriptor for `VpcOriginEndpointConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcOriginEndpointConfigDescriptor = $convert.base64Decode(
    'ChdWcGNPcmlnaW5FbmRwb2ludENvbmZpZxIUCgNhcm4YnZvtvwEgASgJUgNhcm4SHgoIaHR0cH'
    'BvcnQY5+KjswEgASgFUghodHRwcG9ydBIfCglodHRwc3BvcnQYzof9SCABKAVSCWh0dHBzcG9y'
    'dBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lElcKFG9yaWdpbnByb3RvY29scG9saWN5GO7e3G8gAS'
    'gOMiAuY2xvdWRmcm9udC5PcmlnaW5Qcm90b2NvbFBvbGljeVIUb3JpZ2lucHJvdG9jb2xwb2xp'
    'Y3kSUgoSb3JpZ2luc3NscHJvdG9jb2xzGLWhycABIAEoCzIeLmNsb3VkZnJvbnQuT3JpZ2luU3'
    'NsUHJvdG9jb2xzUhJvcmlnaW5zc2xwcm90b2NvbHM=');

@$core.Deprecated('Use vpcOriginListDescriptor instead')
const VpcOriginList$json = {
  '1': 'VpcOriginList',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.cloudfront.VpcOriginSummary',
      '10': 'items'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'quantity', '3': 200749817, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `VpcOriginList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcOriginListDescriptor = $convert.base64Decode(
    'Cg1WcGNPcmlnaW5MaXN0EiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBI1Cg'
    'VpdGVtcxiw8NgBIAMoCzIcLmNsb3VkZnJvbnQuVnBjT3JpZ2luU3VtbWFyeVIFaXRlbXMSGQoG'
    'bWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUghtYXhpdGVtcx'
    'IiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIdCghxdWFudGl0eRj55dxfIAEo'
    'BVIIcXVhbnRpdHk=');

@$core.Deprecated('Use vpcOriginSummaryDescriptor instead')
const VpcOriginSummary$json = {
  '1': 'VpcOriginSummary',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'originendpointarn',
      '3': 52504830,
      '4': 1,
      '5': 9,
      '10': 'originendpointarn'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `VpcOriginSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcOriginSummaryDescriptor = $convert.base64Decode(
    'ChBWcGNPcmlnaW5TdW1tYXJ5Eh8KCWFjY291bnRpZBjSwbkfIAEoCVIJYWNjb3VudGlkEhQKA2'
    'Fybhidm+2/ASABKAlSA2FybhIjCgtjcmVhdGVkdGltZRjz6/M5IAEoCVILY3JlYXRlZHRpbWUS'
    'EgoCaWQYgfKitwEgASgJUgJpZBItChBsYXN0bW9kaWZpZWR0aW1lGOCC/HAgASgJUhBsYXN0bW'
    '9kaWZpZWR0aW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSLwoRb3JpZ2luZW5kcG9pbnRhcm4Y'
    '/tGEGSABKAlSEW9yaWdpbmVuZHBvaW50YXJuEhkKBnN0YXR1cxiQ5PsCIAEoCVIGc3RhdHVz');

@$core.Deprecated('Use webAclCustomizationDescriptor instead')
const WebAclCustomization$json = {
  '1': 'WebAclCustomization',
  '2': [
    {
      '1': 'action',
      '3': 175614240,
      '4': 1,
      '5': 14,
      '6': '.cloudfront.CustomizationActionType',
      '10': 'action'
    },
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `WebAclCustomization`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAclCustomizationDescriptor = $convert.base64Decode(
    'ChNXZWJBY2xDdXN0b21pemF0aW9uEj4KBmFjdGlvbhig0t5TIAEoDjIjLmNsb3VkZnJvbnQuQ3'
    'VzdG9taXphdGlvbkFjdGlvblR5cGVSBmFjdGlvbhIUCgNhcm4YnZvtvwEgASgJUgNhcm4=');
