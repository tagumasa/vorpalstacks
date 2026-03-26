// This is a generated file - do not edit.
//
// Generated from apigateway.proto.

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

@$core.Deprecated('Use accessAssociationSourceTypeDescriptor instead')
const AccessAssociationSourceType$json = {
  '1': 'AccessAssociationSourceType',
  '2': [
    {'1': 'ACCESS_ASSOCIATION_SOURCE_TYPE_VPCE', '2': 0},
  ],
};

/// Descriptor for `AccessAssociationSourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List accessAssociationSourceTypeDescriptor =
    $convert.base64Decode(
        'ChtBY2Nlc3NBc3NvY2lhdGlvblNvdXJjZVR5cGUSJwojQUNDRVNTX0FTU09DSUFUSU9OX1NPVV'
        'JDRV9UWVBFX1ZQQ0UQAA==');

@$core.Deprecated('Use apiKeySourceTypeDescriptor instead')
const ApiKeySourceType$json = {
  '1': 'ApiKeySourceType',
  '2': [
    {'1': 'API_KEY_SOURCE_TYPE_AUTHORIZER', '2': 0},
    {'1': 'API_KEY_SOURCE_TYPE_HEADER', '2': 1},
  ],
};

/// Descriptor for `ApiKeySourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List apiKeySourceTypeDescriptor = $convert.base64Decode(
    'ChBBcGlLZXlTb3VyY2VUeXBlEiIKHkFQSV9LRVlfU09VUkNFX1RZUEVfQVVUSE9SSVpFUhAAEh'
    '4KGkFQSV9LRVlfU09VUkNFX1RZUEVfSEVBREVSEAE=');

@$core.Deprecated('Use apiKeysFormatDescriptor instead')
const ApiKeysFormat$json = {
  '1': 'ApiKeysFormat',
  '2': [
    {'1': 'API_KEYS_FORMAT_CSV', '2': 0},
  ],
};

/// Descriptor for `ApiKeysFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List apiKeysFormatDescriptor = $convert
    .base64Decode('Cg1BcGlLZXlzRm9ybWF0EhcKE0FQSV9LRVlTX0ZPUk1BVF9DU1YQAA==');

@$core.Deprecated('Use apiStatusDescriptor instead')
const ApiStatus$json = {
  '1': 'ApiStatus',
  '2': [
    {'1': 'API_STATUS_UPDATING', '2': 0},
    {'1': 'API_STATUS_PENDING', '2': 1},
    {'1': 'API_STATUS_AVAILABLE', '2': 2},
    {'1': 'API_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `ApiStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List apiStatusDescriptor = $convert.base64Decode(
    'CglBcGlTdGF0dXMSFwoTQVBJX1NUQVRVU19VUERBVElORxAAEhYKEkFQSV9TVEFUVVNfUEVORE'
    'lORxABEhgKFEFQSV9TVEFUVVNfQVZBSUxBQkxFEAISFQoRQVBJX1NUQVRVU19GQUlMRUQQAw==');

@$core.Deprecated('Use authorizerTypeDescriptor instead')
const AuthorizerType$json = {
  '1': 'AuthorizerType',
  '2': [
    {'1': 'AUTHORIZER_TYPE_REQUEST', '2': 0},
    {'1': 'AUTHORIZER_TYPE_COGNITO_USER_POOLS', '2': 1},
    {'1': 'AUTHORIZER_TYPE_TOKEN', '2': 2},
  ],
};

/// Descriptor for `AuthorizerType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List authorizerTypeDescriptor = $convert.base64Decode(
    'Cg5BdXRob3JpemVyVHlwZRIbChdBVVRIT1JJWkVSX1RZUEVfUkVRVUVTVBAAEiYKIkFVVEhPUk'
    'laRVJfVFlQRV9DT0dOSVRPX1VTRVJfUE9PTFMQARIZChVBVVRIT1JJWkVSX1RZUEVfVE9LRU4Q'
    'Ag==');

@$core.Deprecated('Use cacheClusterSizeDescriptor instead')
const CacheClusterSize$json = {
  '1': 'CacheClusterSize',
  '2': [
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_237_GB', '2': 0},
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_13_POINT_5_GB', '2': 1},
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_118_GB', '2': 2},
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_0_POINT_5_GB', '2': 3},
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_1_POINT_6_GB', '2': 4},
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_6_POINT_1_GB', '2': 5},
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_28_POINT_4_GB', '2': 6},
    {'1': 'CACHE_CLUSTER_SIZE_SIZE_58_POINT_2_GB', '2': 7},
  ],
};

/// Descriptor for `CacheClusterSize`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cacheClusterSizeDescriptor = $convert.base64Decode(
    'ChBDYWNoZUNsdXN0ZXJTaXplEiIKHkNBQ0hFX0NMVVNURVJfU0laRV9TSVpFXzIzN19HQhAAEi'
    'kKJUNBQ0hFX0NMVVNURVJfU0laRV9TSVpFXzEzX1BPSU5UXzVfR0IQARIiCh5DQUNIRV9DTFVT'
    'VEVSX1NJWkVfU0laRV8xMThfR0IQAhIoCiRDQUNIRV9DTFVTVEVSX1NJWkVfU0laRV8wX1BPSU'
    '5UXzVfR0IQAxIoCiRDQUNIRV9DTFVTVEVSX1NJWkVfU0laRV8xX1BPSU5UXzZfR0IQBBIoCiRD'
    'QUNIRV9DTFVTVEVSX1NJWkVfU0laRV82X1BPSU5UXzFfR0IQBRIpCiVDQUNIRV9DTFVTVEVSX1'
    'NJWkVfU0laRV8yOF9QT0lOVF80X0dCEAYSKQolQ0FDSEVfQ0xVU1RFUl9TSVpFX1NJWkVfNThf'
    'UE9JTlRfMl9HQhAH');

@$core.Deprecated('Use cacheClusterStatusDescriptor instead')
const CacheClusterStatus$json = {
  '1': 'CacheClusterStatus',
  '2': [
    {'1': 'CACHE_CLUSTER_STATUS_FLUSH_IN_PROGRESS', '2': 0},
    {'1': 'CACHE_CLUSTER_STATUS_AVAILABLE', '2': 1},
    {'1': 'CACHE_CLUSTER_STATUS_CREATE_IN_PROGRESS', '2': 2},
    {'1': 'CACHE_CLUSTER_STATUS_NOT_AVAILABLE', '2': 3},
    {'1': 'CACHE_CLUSTER_STATUS_DELETE_IN_PROGRESS', '2': 4},
  ],
};

/// Descriptor for `CacheClusterStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cacheClusterStatusDescriptor = $convert.base64Decode(
    'ChJDYWNoZUNsdXN0ZXJTdGF0dXMSKgomQ0FDSEVfQ0xVU1RFUl9TVEFUVVNfRkxVU0hfSU5fUF'
    'JPR1JFU1MQABIiCh5DQUNIRV9DTFVTVEVSX1NUQVRVU19BVkFJTEFCTEUQARIrCidDQUNIRV9D'
    'TFVTVEVSX1NUQVRVU19DUkVBVEVfSU5fUFJPR1JFU1MQAhImCiJDQUNIRV9DTFVTVEVSX1NUQV'
    'RVU19OT1RfQVZBSUxBQkxFEAMSKwonQ0FDSEVfQ0xVU1RFUl9TVEFUVVNfREVMRVRFX0lOX1BS'
    'T0dSRVNTEAQ=');

@$core.Deprecated('Use connectionTypeDescriptor instead')
const ConnectionType$json = {
  '1': 'ConnectionType',
  '2': [
    {'1': 'CONNECTION_TYPE_INTERNET', '2': 0},
    {'1': 'CONNECTION_TYPE_VPC_LINK', '2': 1},
  ],
};

/// Descriptor for `ConnectionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionTypeDescriptor = $convert.base64Decode(
    'Cg5Db25uZWN0aW9uVHlwZRIcChhDT05ORUNUSU9OX1RZUEVfSU5URVJORVQQABIcChhDT05ORU'
    'NUSU9OX1RZUEVfVlBDX0xJTksQAQ==');

@$core.Deprecated('Use contentHandlingStrategyDescriptor instead')
const ContentHandlingStrategy$json = {
  '1': 'ContentHandlingStrategy',
  '2': [
    {'1': 'CONTENT_HANDLING_STRATEGY_CONVERT_TO_BINARY', '2': 0},
    {'1': 'CONTENT_HANDLING_STRATEGY_CONVERT_TO_TEXT', '2': 1},
  ],
};

/// Descriptor for `ContentHandlingStrategy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contentHandlingStrategyDescriptor = $convert.base64Decode(
    'ChdDb250ZW50SGFuZGxpbmdTdHJhdGVneRIvCitDT05URU5UX0hBTkRMSU5HX1NUUkFURUdZX0'
    'NPTlZFUlRfVE9fQklOQVJZEAASLQopQ09OVEVOVF9IQU5ETElOR19TVFJBVEVHWV9DT05WRVJU'
    'X1RPX1RFWFQQAQ==');

@$core.Deprecated('Use documentationPartTypeDescriptor instead')
const DocumentationPartType$json = {
  '1': 'DocumentationPartType',
  '2': [
    {'1': 'DOCUMENTATION_PART_TYPE_AUTHORIZER', '2': 0},
    {'1': 'DOCUMENTATION_PART_TYPE_RESPONSE_HEADER', '2': 1},
    {'1': 'DOCUMENTATION_PART_TYPE_RESOURCE', '2': 2},
    {'1': 'DOCUMENTATION_PART_TYPE_PATH_PARAMETER', '2': 3},
    {'1': 'DOCUMENTATION_PART_TYPE_REQUEST_BODY', '2': 4},
    {'1': 'DOCUMENTATION_PART_TYPE_METHOD', '2': 5},
    {'1': 'DOCUMENTATION_PART_TYPE_RESPONSE_BODY', '2': 6},
    {'1': 'DOCUMENTATION_PART_TYPE_QUERY_PARAMETER', '2': 7},
    {'1': 'DOCUMENTATION_PART_TYPE_RESPONSE', '2': 8},
    {'1': 'DOCUMENTATION_PART_TYPE_MODEL', '2': 9},
    {'1': 'DOCUMENTATION_PART_TYPE_REQUEST_HEADER', '2': 10},
    {'1': 'DOCUMENTATION_PART_TYPE_API', '2': 11},
  ],
};

/// Descriptor for `DocumentationPartType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentationPartTypeDescriptor = $convert.base64Decode(
    'ChVEb2N1bWVudGF0aW9uUGFydFR5cGUSJgoiRE9DVU1FTlRBVElPTl9QQVJUX1RZUEVfQVVUSE'
    '9SSVpFUhAAEisKJ0RPQ1VNRU5UQVRJT05fUEFSVF9UWVBFX1JFU1BPTlNFX0hFQURFUhABEiQK'
    'IERPQ1VNRU5UQVRJT05fUEFSVF9UWVBFX1JFU09VUkNFEAISKgomRE9DVU1FTlRBVElPTl9QQV'
    'JUX1RZUEVfUEFUSF9QQVJBTUVURVIQAxIoCiRET0NVTUVOVEFUSU9OX1BBUlRfVFlQRV9SRVFV'
    'RVNUX0JPRFkQBBIiCh5ET0NVTUVOVEFUSU9OX1BBUlRfVFlQRV9NRVRIT0QQBRIpCiVET0NVTU'
    'VOVEFUSU9OX1BBUlRfVFlQRV9SRVNQT05TRV9CT0RZEAYSKwonRE9DVU1FTlRBVElPTl9QQVJU'
    'X1RZUEVfUVVFUllfUEFSQU1FVEVSEAcSJAogRE9DVU1FTlRBVElPTl9QQVJUX1RZUEVfUkVTUE'
    '9OU0UQCBIhCh1ET0NVTUVOVEFUSU9OX1BBUlRfVFlQRV9NT0RFTBAJEioKJkRPQ1VNRU5UQVRJ'
    'T05fUEFSVF9UWVBFX1JFUVVFU1RfSEVBREVSEAoSHwobRE9DVU1FTlRBVElPTl9QQVJUX1RZUE'
    'VfQVBJEAs=');

@$core.Deprecated('Use domainNameStatusDescriptor instead')
const DomainNameStatus$json = {
  '1': 'DomainNameStatus',
  '2': [
    {'1': 'DOMAIN_NAME_STATUS_UPDATING', '2': 0},
    {'1': 'DOMAIN_NAME_STATUS_PENDING', '2': 1},
    {'1': 'DOMAIN_NAME_STATUS_PENDING_CERTIFICATE_REIMPORT', '2': 2},
    {'1': 'DOMAIN_NAME_STATUS_AVAILABLE', '2': 3},
    {'1': 'DOMAIN_NAME_STATUS_PENDING_OWNERSHIP_VERIFICATION', '2': 4},
    {'1': 'DOMAIN_NAME_STATUS_FAILED', '2': 5},
  ],
};

/// Descriptor for `DomainNameStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List domainNameStatusDescriptor = $convert.base64Decode(
    'ChBEb21haW5OYW1lU3RhdHVzEh8KG0RPTUFJTl9OQU1FX1NUQVRVU19VUERBVElORxAAEh4KGk'
    'RPTUFJTl9OQU1FX1NUQVRVU19QRU5ESU5HEAESMwovRE9NQUlOX05BTUVfU1RBVFVTX1BFTkRJ'
    'TkdfQ0VSVElGSUNBVEVfUkVJTVBPUlQQAhIgChxET01BSU5fTkFNRV9TVEFUVVNfQVZBSUxBQk'
    'xFEAMSNQoxRE9NQUlOX05BTUVfU1RBVFVTX1BFTkRJTkdfT1dORVJTSElQX1ZFUklGSUNBVElP'
    'ThAEEh0KGURPTUFJTl9OQU1FX1NUQVRVU19GQUlMRUQQBQ==');

@$core.Deprecated('Use endpointAccessModeDescriptor instead')
const EndpointAccessMode$json = {
  '1': 'EndpointAccessMode',
  '2': [
    {'1': 'ENDPOINT_ACCESS_MODE_STRICT', '2': 0},
    {'1': 'ENDPOINT_ACCESS_MODE_BASIC', '2': 1},
  ],
};

/// Descriptor for `EndpointAccessMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List endpointAccessModeDescriptor = $convert.base64Decode(
    'ChJFbmRwb2ludEFjY2Vzc01vZGUSHwobRU5EUE9JTlRfQUNDRVNTX01PREVfU1RSSUNUEAASHg'
    'oaRU5EUE9JTlRfQUNDRVNTX01PREVfQkFTSUMQAQ==');

@$core.Deprecated('Use endpointTypeDescriptor instead')
const EndpointType$json = {
  '1': 'EndpointType',
  '2': [
    {'1': 'ENDPOINT_TYPE_PRIVATE', '2': 0},
    {'1': 'ENDPOINT_TYPE_REGIONAL', '2': 1},
    {'1': 'ENDPOINT_TYPE_EDGE', '2': 2},
  ],
};

/// Descriptor for `EndpointType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List endpointTypeDescriptor = $convert.base64Decode(
    'CgxFbmRwb2ludFR5cGUSGQoVRU5EUE9JTlRfVFlQRV9QUklWQVRFEAASGgoWRU5EUE9JTlRfVF'
    'lQRV9SRUdJT05BTBABEhYKEkVORFBPSU5UX1RZUEVfRURHRRAC');

@$core.Deprecated('Use gatewayResponseTypeDescriptor instead')
const GatewayResponseType$json = {
  '1': 'GatewayResponseType',
  '2': [
    {'1': 'GATEWAY_RESPONSE_TYPE_WAF_FILTERED', '2': 0},
    {'1': 'GATEWAY_RESPONSE_TYPE_DEFAULT_4XX', '2': 1},
    {'1': 'GATEWAY_RESPONSE_TYPE_UNSUPPORTED_MEDIA_TYPE', '2': 2},
    {'1': 'GATEWAY_RESPONSE_TYPE_REQUEST_TOO_LARGE', '2': 3},
    {'1': 'GATEWAY_RESPONSE_TYPE_ACCESS_DENIED', '2': 4},
    {'1': 'GATEWAY_RESPONSE_TYPE_MISSING_AUTHENTICATION_TOKEN', '2': 5},
    {'1': 'GATEWAY_RESPONSE_TYPE_API_CONFIGURATION_ERROR', '2': 6},
    {'1': 'GATEWAY_RESPONSE_TYPE_THROTTLED', '2': 7},
    {'1': 'GATEWAY_RESPONSE_TYPE_INTEGRATION_FAILURE', '2': 8},
    {'1': 'GATEWAY_RESPONSE_TYPE_RESOURCE_NOT_FOUND', '2': 9},
    {'1': 'GATEWAY_RESPONSE_TYPE_QUOTA_EXCEEDED', '2': 10},
    {'1': 'GATEWAY_RESPONSE_TYPE_AUTHORIZER_FAILURE', '2': 11},
    {'1': 'GATEWAY_RESPONSE_TYPE_INVALID_SIGNATURE', '2': 12},
    {'1': 'GATEWAY_RESPONSE_TYPE_INTEGRATION_TIMEOUT', '2': 13},
    {'1': 'GATEWAY_RESPONSE_TYPE_UNAUTHORIZED', '2': 14},
    {'1': 'GATEWAY_RESPONSE_TYPE_EXPIRED_TOKEN', '2': 15},
    {'1': 'GATEWAY_RESPONSE_TYPE_INVALID_API_KEY', '2': 16},
    {'1': 'GATEWAY_RESPONSE_TYPE_BAD_REQUEST_BODY', '2': 17},
    {'1': 'GATEWAY_RESPONSE_TYPE_BAD_REQUEST_PARAMETERS', '2': 18},
    {'1': 'GATEWAY_RESPONSE_TYPE_DEFAULT_5XX', '2': 19},
    {'1': 'GATEWAY_RESPONSE_TYPE_AUTHORIZER_CONFIGURATION_ERROR', '2': 20},
  ],
};

/// Descriptor for `GatewayResponseType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List gatewayResponseTypeDescriptor = $convert.base64Decode(
    'ChNHYXRld2F5UmVzcG9uc2VUeXBlEiYKIkdBVEVXQVlfUkVTUE9OU0VfVFlQRV9XQUZfRklMVE'
    'VSRUQQABIlCiFHQVRFV0FZX1JFU1BPTlNFX1RZUEVfREVGQVVMVF80WFgQARIwCixHQVRFV0FZ'
    'X1JFU1BPTlNFX1RZUEVfVU5TVVBQT1JURURfTUVESUFfVFlQRRACEisKJ0dBVEVXQVlfUkVTUE'
    '9OU0VfVFlQRV9SRVFVRVNUX1RPT19MQVJHRRADEicKI0dBVEVXQVlfUkVTUE9OU0VfVFlQRV9B'
    'Q0NFU1NfREVOSUVEEAQSNgoyR0FURVdBWV9SRVNQT05TRV9UWVBFX01JU1NJTkdfQVVUSEVOVE'
    'lDQVRJT05fVE9LRU4QBRIxCi1HQVRFV0FZX1JFU1BPTlNFX1RZUEVfQVBJX0NPTkZJR1VSQVRJ'
    'T05fRVJST1IQBhIjCh9HQVRFV0FZX1JFU1BPTlNFX1RZUEVfVEhST1RUTEVEEAcSLQopR0FURV'
    'dBWV9SRVNQT05TRV9UWVBFX0lOVEVHUkFUSU9OX0ZBSUxVUkUQCBIsCihHQVRFV0FZX1JFU1BP'
    'TlNFX1RZUEVfUkVTT1VSQ0VfTk9UX0ZPVU5EEAkSKAokR0FURVdBWV9SRVNQT05TRV9UWVBFX1'
    'FVT1RBX0VYQ0VFREVEEAoSLAooR0FURVdBWV9SRVNQT05TRV9UWVBFX0FVVEhPUklaRVJfRkFJ'
    'TFVSRRALEisKJ0dBVEVXQVlfUkVTUE9OU0VfVFlQRV9JTlZBTElEX1NJR05BVFVSRRAMEi0KKU'
    'dBVEVXQVlfUkVTUE9OU0VfVFlQRV9JTlRFR1JBVElPTl9USU1FT1VUEA0SJgoiR0FURVdBWV9S'
    'RVNQT05TRV9UWVBFX1VOQVVUSE9SSVpFRBAOEicKI0dBVEVXQVlfUkVTUE9OU0VfVFlQRV9FWF'
    'BJUkVEX1RPS0VOEA8SKQolR0FURVdBWV9SRVNQT05TRV9UWVBFX0lOVkFMSURfQVBJX0tFWRAQ'
    'EioKJkdBVEVXQVlfUkVTUE9OU0VfVFlQRV9CQURfUkVRVUVTVF9CT0RZEBESMAosR0FURVdBWV'
    '9SRVNQT05TRV9UWVBFX0JBRF9SRVFVRVNUX1BBUkFNRVRFUlMQEhIlCiFHQVRFV0FZX1JFU1BP'
    'TlNFX1RZUEVfREVGQVVMVF81WFgQExI4CjRHQVRFV0FZX1JFU1BPTlNFX1RZUEVfQVVUSE9SSV'
    'pFUl9DT05GSUdVUkFUSU9OX0VSUk9SEBQ=');

@$core.Deprecated('Use integrationTypeDescriptor instead')
const IntegrationType$json = {
  '1': 'IntegrationType',
  '2': [
    {'1': 'INTEGRATION_TYPE_HTTP', '2': 0},
    {'1': 'INTEGRATION_TYPE_AWS_PROXY', '2': 1},
    {'1': 'INTEGRATION_TYPE_MOCK', '2': 2},
    {'1': 'INTEGRATION_TYPE_AWS', '2': 3},
    {'1': 'INTEGRATION_TYPE_HTTP_PROXY', '2': 4},
  ],
};

/// Descriptor for `IntegrationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List integrationTypeDescriptor = $convert.base64Decode(
    'Cg9JbnRlZ3JhdGlvblR5cGUSGQoVSU5URUdSQVRJT05fVFlQRV9IVFRQEAASHgoaSU5URUdSQV'
    'RJT05fVFlQRV9BV1NfUFJPWFkQARIZChVJTlRFR1JBVElPTl9UWVBFX01PQ0sQAhIYChRJTlRF'
    'R1JBVElPTl9UWVBFX0FXUxADEh8KG0lOVEVHUkFUSU9OX1RZUEVfSFRUUF9QUk9YWRAE');

@$core.Deprecated('Use ipAddressTypeDescriptor instead')
const IpAddressType$json = {
  '1': 'IpAddressType',
  '2': [
    {'1': 'IP_ADDRESS_TYPE_DUALSTACK', '2': 0},
    {'1': 'IP_ADDRESS_TYPE_IPV4', '2': 1},
  ],
};

/// Descriptor for `IpAddressType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List ipAddressTypeDescriptor = $convert.base64Decode(
    'Cg1JcEFkZHJlc3NUeXBlEh0KGUlQX0FERFJFU1NfVFlQRV9EVUFMU1RBQ0sQABIYChRJUF9BRE'
    'RSRVNTX1RZUEVfSVBWNBAB');

@$core.Deprecated('Use locationStatusTypeDescriptor instead')
const LocationStatusType$json = {
  '1': 'LocationStatusType',
  '2': [
    {'1': 'LOCATION_STATUS_TYPE_DOCUMENTED', '2': 0},
    {'1': 'LOCATION_STATUS_TYPE_UNDOCUMENTED', '2': 1},
  ],
};

/// Descriptor for `LocationStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List locationStatusTypeDescriptor = $convert.base64Decode(
    'ChJMb2NhdGlvblN0YXR1c1R5cGUSIwofTE9DQVRJT05fU1RBVFVTX1RZUEVfRE9DVU1FTlRFRB'
    'AAEiUKIUxPQ0FUSU9OX1NUQVRVU19UWVBFX1VORE9DVU1FTlRFRBAB');

@$core.Deprecated('Use opDescriptor instead')
const Op$json = {
  '1': 'Op',
  '2': [
    {'1': 'OP_ADD', '2': 0},
    {'1': 'OP_REMOVE', '2': 1},
    {'1': 'OP_MOVE', '2': 2},
    {'1': 'OP_TEST', '2': 3},
    {'1': 'OP_REPLACE', '2': 4},
    {'1': 'OP_COPY', '2': 5},
  ],
};

/// Descriptor for `Op`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opDescriptor = $convert.base64Decode(
    'CgJPcBIKCgZPUF9BREQQABINCglPUF9SRU1PVkUQARILCgdPUF9NT1ZFEAISCwoHT1BfVEVTVB'
    'ADEg4KCk9QX1JFUExBQ0UQBBILCgdPUF9DT1BZEAU=');

@$core.Deprecated('Use putModeDescriptor instead')
const PutMode$json = {
  '1': 'PutMode',
  '2': [
    {'1': 'PUT_MODE_OVERWRITE', '2': 0},
    {'1': 'PUT_MODE_MERGE', '2': 1},
  ],
};

/// Descriptor for `PutMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List putModeDescriptor = $convert.base64Decode(
    'CgdQdXRNb2RlEhYKElBVVF9NT0RFX09WRVJXUklURRAAEhIKDlBVVF9NT0RFX01FUkdFEAE=');

@$core.Deprecated('Use quotaPeriodTypeDescriptor instead')
const QuotaPeriodType$json = {
  '1': 'QuotaPeriodType',
  '2': [
    {'1': 'QUOTA_PERIOD_TYPE_WEEK', '2': 0},
    {'1': 'QUOTA_PERIOD_TYPE_DAY', '2': 1},
    {'1': 'QUOTA_PERIOD_TYPE_MONTH', '2': 2},
  ],
};

/// Descriptor for `QuotaPeriodType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List quotaPeriodTypeDescriptor = $convert.base64Decode(
    'Cg9RdW90YVBlcmlvZFR5cGUSGgoWUVVPVEFfUEVSSU9EX1RZUEVfV0VFSxAAEhkKFVFVT1RBX1'
    'BFUklPRF9UWVBFX0RBWRABEhsKF1FVT1RBX1BFUklPRF9UWVBFX01PTlRIEAI=');

@$core.Deprecated('Use resourceOwnerDescriptor instead')
const ResourceOwner$json = {
  '1': 'ResourceOwner',
  '2': [
    {'1': 'RESOURCE_OWNER_OTHER_ACCOUNTS', '2': 0},
    {'1': 'RESOURCE_OWNER_SELF', '2': 1},
  ],
};

/// Descriptor for `ResourceOwner`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resourceOwnerDescriptor = $convert.base64Decode(
    'Cg1SZXNvdXJjZU93bmVyEiEKHVJFU09VUkNFX09XTkVSX09USEVSX0FDQ09VTlRTEAASFwoTUk'
    'VTT1VSQ0VfT1dORVJfU0VMRhAB');

@$core.Deprecated('Use responseTransferModeDescriptor instead')
const ResponseTransferMode$json = {
  '1': 'ResponseTransferMode',
  '2': [
    {'1': 'RESPONSE_TRANSFER_MODE_STREAM', '2': 0},
    {'1': 'RESPONSE_TRANSFER_MODE_BUFFERED', '2': 1},
  ],
};

/// Descriptor for `ResponseTransferMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List responseTransferModeDescriptor = $convert.base64Decode(
    'ChRSZXNwb25zZVRyYW5zZmVyTW9kZRIhCh1SRVNQT05TRV9UUkFOU0ZFUl9NT0RFX1NUUkVBTR'
    'AAEiMKH1JFU1BPTlNFX1RSQU5TRkVSX01PREVfQlVGRkVSRUQQAQ==');

@$core.Deprecated('Use routingModeDescriptor instead')
const RoutingMode$json = {
  '1': 'RoutingMode',
  '2': [
    {'1': 'ROUTING_MODE_ROUTING_RULE_THEN_BASE_PATH_MAPPING', '2': 0},
    {'1': 'ROUTING_MODE_BASE_PATH_MAPPING_ONLY', '2': 1},
    {'1': 'ROUTING_MODE_ROUTING_RULE_ONLY', '2': 2},
  ],
};

/// Descriptor for `RoutingMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List routingModeDescriptor = $convert.base64Decode(
    'CgtSb3V0aW5nTW9kZRI0CjBST1VUSU5HX01PREVfUk9VVElOR19SVUxFX1RIRU5fQkFTRV9QQV'
    'RIX01BUFBJTkcQABInCiNST1VUSU5HX01PREVfQkFTRV9QQVRIX01BUFBJTkdfT05MWRABEiIK'
    'HlJPVVRJTkdfTU9ERV9ST1VUSU5HX1JVTEVfT05MWRAC');

@$core.Deprecated('Use securityPolicyDescriptor instead')
const SecurityPolicy$json = {
  '1': 'SecurityPolicy',
  '2': [
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PQ_2025_09', '2': 0},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS13_2025_EDGE', '2': 1},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PQ_2025_09', '2': 2},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS12_2018_EDGE', '2': 3},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS12_PFS_2025_EDGE', '2': 4},
    {'1': 'SECURITY_POLICY_TLS_1_2', '2': 5},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_2021_06', '2': 6},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_2025_09', '2': 7},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PFS_PQ_2025_09', '2': 8},
    {'1': 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_FIPS_2025_09', '2': 9},
    {'1': 'SECURITY_POLICY_TLS_1_0', '2': 10},
  ],
};

/// Descriptor for `SecurityPolicy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List securityPolicyDescriptor = $convert.base64Decode(
    'Cg5TZWN1cml0eVBvbGljeRI3CjNTRUNVUklUWV9QT0xJQ1lfU0VDVVJJVFlQT0xJQ1lfVExTMT'
    'NfMV8yX1BRXzIwMjVfMDkQABIyCi5TRUNVUklUWV9QT0xJQ1lfU0VDVVJJVFlQT0xJQ1lfVExT'
    'MTNfMjAyNV9FREdFEAESPAo4U0VDVVJJVFlfUE9MSUNZX1NFQ1VSSVRZUE9MSUNZX1RMUzEzXz'
    'FfMl9GSVBTX1BRXzIwMjVfMDkQAhIyCi5TRUNVUklUWV9QT0xJQ1lfU0VDVVJJVFlQT0xJQ1lf'
    'VExTMTJfMjAxOF9FREdFEAMSNgoyU0VDVVJJVFlfUE9MSUNZX1NFQ1VSSVRZUE9MSUNZX1RMUz'
    'EyX1BGU18yMDI1X0VER0UQBBIbChdTRUNVUklUWV9QT0xJQ1lfVExTXzFfMhAFEjQKMFNFQ1VS'
    'SVRZX1BPTElDWV9TRUNVUklUWVBPTElDWV9UTFMxM18xXzJfMjAyMV8wNhAGEjQKMFNFQ1VSSV'
    'RZX1BPTElDWV9TRUNVUklUWVBPTElDWV9UTFMxM18xXzNfMjAyNV8wORAHEjsKN1NFQ1VSSVRZ'
    'X1BPTElDWV9TRUNVUklUWVBPTElDWV9UTFMxM18xXzJfUEZTX1BRXzIwMjVfMDkQCBI5CjVTRU'
    'NVUklUWV9QT0xJQ1lfU0VDVVJJVFlQT0xJQ1lfVExTMTNfMV8zX0ZJUFNfMjAyNV8wORAJEhsK'
    'F1NFQ1VSSVRZX1BPTElDWV9UTFNfMV8wEAo=');

@$core
    .Deprecated('Use unauthorizedCacheControlHeaderStrategyDescriptor instead')
const UnauthorizedCacheControlHeaderStrategy$json = {
  '1': 'UnauthorizedCacheControlHeaderStrategy',
  '2': [
    {
      '1':
          'UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITH_RESPONSE_HEADER',
      '2': 0
    },
    {'1': 'UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_FAIL_WITH_403', '2': 1},
    {
      '1':
          'UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITHOUT_RESPONSE_HEADER',
      '2': 2
    },
  ],
};

/// Descriptor for `UnauthorizedCacheControlHeaderStrategy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List unauthorizedCacheControlHeaderStrategyDescriptor =
    $convert.base64Decode(
        'CiZVbmF1dGhvcml6ZWRDYWNoZUNvbnRyb2xIZWFkZXJTdHJhdGVneRJLCkdVTkFVVEhPUklaRU'
        'RfQ0FDSEVfQ09OVFJPTF9IRUFERVJfU1RSQVRFR1lfU1VDQ0VFRF9XSVRIX1JFU1BPTlNFX0hF'
        'QURFUhAAEjwKOFVOQVVUSE9SSVpFRF9DQUNIRV9DT05UUk9MX0hFQURFUl9TVFJBVEVHWV9GQU'
        'lMX1dJVEhfNDAzEAESTgpKVU5BVVRIT1JJWkVEX0NBQ0hFX0NPTlRST0xfSEVBREVSX1NUUkFU'
        'RUdZX1NVQ0NFRURfV0lUSE9VVF9SRVNQT05TRV9IRUFERVIQAg==');

@$core.Deprecated('Use vpcLinkStatusDescriptor instead')
const VpcLinkStatus$json = {
  '1': 'VpcLinkStatus',
  '2': [
    {'1': 'VPC_LINK_STATUS_PENDING', '2': 0},
    {'1': 'VPC_LINK_STATUS_AVAILABLE', '2': 1},
    {'1': 'VPC_LINK_STATUS_DELETING', '2': 2},
    {'1': 'VPC_LINK_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `VpcLinkStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List vpcLinkStatusDescriptor = $convert.base64Decode(
    'Cg1WcGNMaW5rU3RhdHVzEhsKF1ZQQ19MSU5LX1NUQVRVU19QRU5ESU5HEAASHQoZVlBDX0xJTk'
    'tfU1RBVFVTX0FWQUlMQUJMRRABEhwKGFZQQ19MSU5LX1NUQVRVU19ERUxFVElORxACEhoKFlZQ'
    'Q19MSU5LX1NUQVRVU19GQUlMRUQQAw==');

@$core.Deprecated('Use accessLogSettingsDescriptor instead')
const AccessLogSettings$json = {
  '1': 'AccessLogSettings',
  '2': [
    {
      '1': 'destinationarn',
      '3': 427601315,
      '4': 1,
      '5': 9,
      '10': 'destinationarn'
    },
    {'1': 'format', '3': 429753683, '4': 1, '5': 9, '10': 'format'},
  ],
};

/// Descriptor for `AccessLogSettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessLogSettingsDescriptor = $convert.base64Decode(
    'ChFBY2Nlc3NMb2dTZXR0aW5ncxIqCg5kZXN0aW5hdGlvbmFybhij2/LLASABKAlSDmRlc3Rpbm'
    'F0aW9uYXJuEhoKBmZvcm1hdBjTivbMASABKAlSBmZvcm1hdA==');

@$core.Deprecated('Use accountDescriptor instead')
const Account$json = {
  '1': 'Account',
  '2': [
    {
      '1': 'apikeyversion',
      '3': 148957539,
      '4': 1,
      '5': 9,
      '10': 'apikeyversion'
    },
    {
      '1': 'cloudwatchrolearn',
      '3': 144858329,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchrolearn'
    },
    {'1': 'features', '3': 528712651, '4': 3, '5': 9, '10': 'features'},
    {
      '1': 'throttlesettings',
      '3': 165000097,
      '4': 1,
      '5': 11,
      '6': '.apigateway.ThrottleSettings',
      '10': 'throttlesettings'
    },
  ],
};

/// Descriptor for `Account`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountDescriptor = $convert.base64Decode(
    'CgdBY2NvdW50EicKDWFwaWtleXZlcnNpb24Y49KDRyABKAlSDWFwaWtleXZlcnNpb24SLwoRY2'
    'xvdWR3YXRjaHJvbGVhcm4Y2bmJRSABKAlSEWNsb3Vkd2F0Y2hyb2xlYXJuEh4KCGZlYXR1cmVz'
    'GMuHjvwBIAMoCVIIZmVhdHVyZXMSSwoQdGhyb3R0bGVzZXR0aW5ncxih59ZOIAEoCzIcLmFwaW'
    'dhdGV3YXkuVGhyb3R0bGVTZXR0aW5nc1IQdGhyb3R0bGVzZXR0aW5ncw==');

@$core.Deprecated('Use apiKeyDescriptor instead')
const ApiKey$json = {
  '1': 'ApiKey',
  '2': [
    {'1': 'createddate', '3': 53061200, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'customerid', '3': 227830269, '4': 1, '5': 9, '10': 'customerid'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'enabled', '3': 49525663, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lastupdateddate',
      '3': 448453361,
      '4': 1,
      '5': 9,
      '10': 'lastupdateddate'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'stagekeys', '3': 287991830, '4': 3, '5': 9, '10': 'stagekeys'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ApiKey.TagsEntry',
      '10': 'tags'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
  '3': [ApiKey_TagsEntry$json],
};

@$core.Deprecated('Use apiKeyDescriptor instead')
const ApiKey_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ApiKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List apiKeyDescriptor = $convert.base64Decode(
    'CgZBcGlLZXkSIwoLY3JlYXRlZGRhdGUY0MymGSABKAlSC2NyZWF0ZWRkYXRlEiEKCmN1c3RvbW'
    'VyaWQY/dPRbCABKAlSCmN1c3RvbWVyaWQSJAoLZGVzY3JpcHRpb24Y6va8owEgASgJUgtkZXNj'
    'cmlwdGlvbhIbCgdlbmFibGVkGJ/nzhcgASgIUgdlbmFibGVkEhIKAmlkGOHV4bkBIAEoCVICaW'
    'QSLAoPbGFzdHVwZGF0ZWRkYXRlGPG169UBIAEoCVIPbGFzdHVwZGF0ZWRkYXRlEhUKBG5hbWUY'
    '5/vmaSABKAlSBG5hbWUSIAoJc3RhZ2VrZXlzGJbQqYkBIAMoCVIJc3RhZ2VrZXlzEjQKBHRhZ3'
    'MYodfboAEgAygLMhwuYXBpZ2F0ZXdheS5BcGlLZXkuVGFnc0VudHJ5UgR0YWdzEhcKBXZhbHVl'
    'GMun+xIgASgJUgV2YWx1ZRo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdW'
    'UYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use apiKeyIdsDescriptor instead')
const ApiKeyIds$json = {
  '1': 'ApiKeyIds',
  '2': [
    {'1': 'ids', '3': 13616490, '4': 3, '5': 9, '10': 'ids'},
    {'1': 'warnings', '3': 185617301, '4': 3, '5': 9, '10': 'warnings'},
  ],
};

/// Descriptor for `ApiKeyIds`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List apiKeyIdsDescriptor = $convert.base64Decode(
    'CglBcGlLZXlJZHMSEwoDaWRzGOqKvwYgAygJUgNpZHMSHQoId2FybmluZ3MYlZfBWCADKAlSCH'
    'dhcm5pbmdz');

@$core.Deprecated('Use apiKeysDescriptor instead')
const ApiKeys$json = {
  '1': 'ApiKeys',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ApiKey',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'warnings', '3': 185617301, '4': 3, '5': 9, '10': 'warnings'},
  ],
};

/// Descriptor for `ApiKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List apiKeysDescriptor = $convert.base64Decode(
    'CgdBcGlLZXlzEiwKBWl0ZW1zGJDn5NMBIAMoCzISLmFwaWdhdGV3YXkuQXBpS2V5UgVpdGVtcx'
    'IeCghwb3NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9uEh0KCHdhcm5pbmdzGJWXwVggAygJUgh3'
    'YXJuaW5ncw==');

@$core.Deprecated('Use apiStageDescriptor instead')
const ApiStage$json = {
  '1': 'ApiStage',
  '2': [
    {'1': 'apiid', '3': 113380971, '4': 1, '5': 9, '10': 'apiid'},
    {'1': 'stage', '3': 140155438, '4': 1, '5': 9, '10': 'stage'},
    {
      '1': 'throttle',
      '3': 395260638,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ApiStage.ThrottleEntry',
      '10': 'throttle'
    },
  ],
  '3': [ApiStage_ThrottleEntry$json],
};

@$core.Deprecated('Use apiStageDescriptor instead')
const ApiStage_ThrottleEntry$json = {
  '1': 'ThrottleEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.apigateway.ThrottleSettings',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ApiStage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List apiStageDescriptor = $convert.base64Decode(
    'CghBcGlTdGFnZRIXCgVhcGlpZBjrnIg2IAEoCVIFYXBpaWQSFwoFc3RhZ2UYrrTqQiABKAlSBX'
    'N0YWdlEkIKCHRocm90dGxlGN7lvLwBIAMoCzIiLmFwaWdhdGV3YXkuQXBpU3RhZ2UuVGhyb3R0'
    'bGVFbnRyeVIIdGhyb3R0bGUaWQoNVGhyb3R0bGVFbnRyeRIQCgNrZXkYASABKAlSA2tleRIyCg'
    'V2YWx1ZRgCIAEoCzIcLmFwaWdhdGV3YXkuVGhyb3R0bGVTZXR0aW5nc1IFdmFsdWU6AjgB');

@$core.Deprecated('Use authorizerDescriptor instead')
const Authorizer$json = {
  '1': 'Authorizer',
  '2': [
    {'1': 'authtype', '3': 162773848, '4': 1, '5': 9, '10': 'authtype'},
    {
      '1': 'authorizercredentials',
      '3': 233575233,
      '4': 1,
      '5': 9,
      '10': 'authorizercredentials'
    },
    {
      '1': 'authorizerresultttlinseconds',
      '3': 135440208,
      '4': 1,
      '5': 5,
      '10': 'authorizerresultttlinseconds'
    },
    {
      '1': 'authorizeruri',
      '3': 525146137,
      '4': 1,
      '5': 9,
      '10': 'authorizeruri'
    },
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'identitysource',
      '3': 285615231,
      '4': 1,
      '5': 9,
      '10': 'identitysource'
    },
    {
      '1': 'identityvalidationexpression',
      '3': 227211199,
      '4': 1,
      '5': 9,
      '10': 'identityvalidationexpression'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'providerarns', '3': 301486689, '4': 3, '5': 9, '10': 'providerarns'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.apigateway.AuthorizerType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `Authorizer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authorizerDescriptor = $convert.base64Decode(
    'CgpBdXRob3JpemVyEh0KCGF1dGh0eXBlGNj2zk0gASgJUghhdXRodHlwZRI3ChVhdXRob3Jpem'
    'VyY3JlZGVudGlhbHMYwaawbyABKAlSFWF1dGhvcml6ZXJjcmVkZW50aWFscxJFChxhdXRob3Jp'
    'emVycmVzdWx0dHRsaW5zZWNvbmRzGNDOykAgASgFUhxhdXRob3JpemVycmVzdWx0dHRsaW5zZW'
    'NvbmRzEigKDWF1dGhvcml6ZXJ1cmkYmbC0+gEgASgJUg1hdXRob3JpemVydXJpEhIKAmlkGOHV'
    '4bkBIAEoCVICaWQSKgoOaWRlbnRpdHlzb3VyY2UY/8iYiAEgASgJUg5pZGVudGl0eXNvdXJjZR'
    'JFChxpZGVudGl0eXZhbGlkYXRpb25leHByZXNzaW9uGL/vq2wgASgJUhxpZGVudGl0eXZhbGlk'
    'YXRpb25leHByZXNzaW9uEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSJgoMcHJvdmlkZXJhcm5zGO'
    'Gk4Y8BIAMoCVIMcHJvdmlkZXJhcm5zEjIKBHR5cGUYzuKfiQEgASgOMhouYXBpZ2F0ZXdheS5B'
    'dXRob3JpemVyVHlwZVIEdHlwZQ==');

@$core.Deprecated('Use authorizersDescriptor instead')
const Authorizers$json = {
  '1': 'Authorizers',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Authorizer',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `Authorizers`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authorizersDescriptor = $convert.base64Decode(
    'CgtBdXRob3JpemVycxIwCgVpdGVtcxiQ5+TTASADKAsyFi5hcGlnYXRld2F5LkF1dGhvcml6ZX'
    'JSBWl0ZW1zEh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24=');

@$core.Deprecated('Use badRequestExceptionDescriptor instead')
const BadRequestException$json = {
  '1': 'BadRequestException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BadRequestException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List badRequestExceptionDescriptor =
    $convert.base64Decode(
        'ChNCYWRSZXF1ZXN0RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use basePathMappingDescriptor instead')
const BasePathMapping$json = {
  '1': 'BasePathMapping',
  '2': [
    {'1': 'basepath', '3': 267528880, '4': 1, '5': 9, '10': 'basepath'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stage', '3': 140155438, '4': 1, '5': 9, '10': 'stage'},
  ],
};

/// Descriptor for `BasePathMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List basePathMappingDescriptor = $convert.base64Decode(
    'Cg9CYXNlUGF0aE1hcHBpbmcSHQoIYmFzZXBhdGgYsNXIfyABKAlSCGJhc2VwYXRoEiAKCXJlc3'
    'RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZBIXCgVzdGFnZRiutOpCIAEoCVIFc3RhZ2U=');

@$core.Deprecated('Use basePathMappingsDescriptor instead')
const BasePathMappings$json = {
  '1': 'BasePathMappings',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.BasePathMapping',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `BasePathMappings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List basePathMappingsDescriptor = $convert.base64Decode(
    'ChBCYXNlUGF0aE1hcHBpbmdzEjUKBWl0ZW1zGJDn5NMBIAMoCzIbLmFwaWdhdGV3YXkuQmFzZV'
    'BhdGhNYXBwaW5nUgVpdGVtcxIeCghwb3NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use canarySettingsDescriptor instead')
const CanarySettings$json = {
  '1': 'CanarySettings',
  '2': [
    {'1': 'deploymentid', '3': 439369188, '4': 1, '5': 9, '10': 'deploymentid'},
    {
      '1': 'percenttraffic',
      '3': 147177864,
      '4': 1,
      '5': 1,
      '10': 'percenttraffic'
    },
    {
      '1': 'stagevariableoverrides',
      '3': 221124259,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CanarySettings.StagevariableoverridesEntry',
      '10': 'stagevariableoverrides'
    },
    {
      '1': 'usestagecache',
      '3': 179841697,
      '4': 1,
      '5': 8,
      '10': 'usestagecache'
    },
  ],
  '3': [CanarySettings_StagevariableoverridesEntry$json],
};

@$core.Deprecated('Use canarySettingsDescriptor instead')
const CanarySettings_StagevariableoverridesEntry$json = {
  '1': 'StagevariableoverridesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CanarySettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List canarySettingsDescriptor = $convert.base64Decode(
    'Cg5DYW5hcnlTZXR0aW5ncxImCgxkZXBsb3ltZW50aWQY5PvA0QEgASgJUgxkZXBsb3ltZW50aW'
    'QSKQoOcGVyY2VudHRyYWZmaWMYiIOXRiABKAFSDnBlcmNlbnR0cmFmZmljEnEKFnN0YWdldmFy'
    'aWFibGVvdmVycmlkZXMYo624aSADKAsyNi5hcGlnYXRld2F5LkNhbmFyeVNldHRpbmdzLlN0YW'
    'dldmFyaWFibGVvdmVycmlkZXNFbnRyeVIWc3RhZ2V2YXJpYWJsZW92ZXJyaWRlcxInCg11c2Vz'
    'dGFnZWNhY2hlGKHV4FUgASgIUg11c2VzdGFnZWNhY2hlGkkKG1N0YWdldmFyaWFibGVvdmVycm'
    'lkZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use clientCertificateDescriptor instead')
const ClientCertificate$json = {
  '1': 'ClientCertificate',
  '2': [
    {
      '1': 'clientcertificateid',
      '3': 276222909,
      '4': 1,
      '5': 9,
      '10': 'clientcertificateid'
    },
    {'1': 'createddate', '3': 53061200, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'expirationdate',
      '3': 391033245,
      '4': 1,
      '5': 9,
      '10': 'expirationdate'
    },
    {
      '1': 'pemencodedcertificate',
      '3': 449484817,
      '4': 1,
      '5': 9,
      '10': 'pemencodedcertificate'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ClientCertificate.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [ClientCertificate_TagsEntry$json],
};

@$core.Deprecated('Use clientCertificateDescriptor instead')
const ClientCertificate_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ClientCertificate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clientCertificateDescriptor = $convert.base64Decode(
    'ChFDbGllbnRDZXJ0aWZpY2F0ZRI0ChNjbGllbnRjZXJ0aWZpY2F0ZWlkGL2n24MBIAEoCVITY2'
    'xpZW50Y2VydGlmaWNhdGVpZBIjCgtjcmVhdGVkZGF0ZRjQzKYZIAEoCVILY3JlYXRlZGRhdGUS'
    'JAoLZGVzY3JpcHRpb24Y6va8owEgASgJUgtkZXNjcmlwdGlvbhIqCg5leHBpcmF0aW9uZGF0ZR'
    'id47q6ASABKAlSDmV4cGlyYXRpb25kYXRlEjgKFXBlbWVuY29kZWRjZXJ0aWZpY2F0ZRiRsKrW'
    'ASABKAlSFXBlbWVuY29kZWRjZXJ0aWZpY2F0ZRI/CgR0YWdzGKHX26ABIAMoCzInLmFwaWdhdG'
    'V3YXkuQ2xpZW50Q2VydGlmaWNhdGUuVGFnc0VudHJ5UgR0YWdzGjcKCVRhZ3NFbnRyeRIQCgNr'
    'ZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use clientCertificatesDescriptor instead')
const ClientCertificates$json = {
  '1': 'ClientCertificates',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ClientCertificate',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `ClientCertificates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clientCertificatesDescriptor = $convert.base64Decode(
    'ChJDbGllbnRDZXJ0aWZpY2F0ZXMSNwoFaXRlbXMYkOfk0wEgAygLMh0uYXBpZ2F0ZXdheS5DbG'
    'llbnRDZXJ0aWZpY2F0ZVIFaXRlbXMSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbg==');

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

@$core.Deprecated('Use createApiKeyRequestDescriptor instead')
const CreateApiKeyRequest$json = {
  '1': 'CreateApiKeyRequest',
  '2': [
    {'1': 'customerid', '3': 227830269, '4': 1, '5': 9, '10': 'customerid'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'enabled', '3': 49525663, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'generatedistinctid',
      '3': 99833588,
      '4': 1,
      '5': 8,
      '10': 'generatedistinctid'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'stagekeys',
      '3': 287991830,
      '4': 3,
      '5': 11,
      '6': '.apigateway.StageKey',
      '10': 'stagekeys'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateApiKeyRequest.TagsEntry',
      '10': 'tags'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
  '3': [CreateApiKeyRequest_TagsEntry$json],
};

@$core.Deprecated('Use createApiKeyRequestDescriptor instead')
const CreateApiKeyRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateApiKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createApiKeyRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVBcGlLZXlSZXF1ZXN0EiEKCmN1c3RvbWVyaWQY/dPRbCABKAlSCmN1c3RvbWVyaW'
    'QSJAoLZGVzY3JpcHRpb24Y6va8owEgASgJUgtkZXNjcmlwdGlvbhIbCgdlbmFibGVkGJ/nzhcg'
    'ASgIUgdlbmFibGVkEjEKEmdlbmVyYXRlZGlzdGluY3RpZBj0rc0vIAEoCFISZ2VuZXJhdGVkaX'
    'N0aW5jdGlkEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSNgoJc3RhZ2VrZXlzGJbQqYkBIAMoCzIU'
    'LmFwaWdhdGV3YXkuU3RhZ2VLZXlSCXN0YWdla2V5cxJBCgR0YWdzGKHX26ABIAMoCzIpLmFwaW'
    'dhdGV3YXkuQ3JlYXRlQXBpS2V5UmVxdWVzdC5UYWdzRW50cnlSBHRhZ3MSFwoFdmFsdWUYy6f7'
    'EiABKAlSBXZhbHVlGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIA'
    'EoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use createAuthorizerRequestDescriptor instead')
const CreateAuthorizerRequest$json = {
  '1': 'CreateAuthorizerRequest',
  '2': [
    {'1': 'authtype', '3': 162773848, '4': 1, '5': 9, '10': 'authtype'},
    {
      '1': 'authorizercredentials',
      '3': 233575233,
      '4': 1,
      '5': 9,
      '10': 'authorizercredentials'
    },
    {
      '1': 'authorizerresultttlinseconds',
      '3': 135440208,
      '4': 1,
      '5': 5,
      '10': 'authorizerresultttlinseconds'
    },
    {
      '1': 'authorizeruri',
      '3': 525146137,
      '4': 1,
      '5': 9,
      '10': 'authorizeruri'
    },
    {
      '1': 'identitysource',
      '3': 285615231,
      '4': 1,
      '5': 9,
      '10': 'identitysource'
    },
    {
      '1': 'identityvalidationexpression',
      '3': 227211199,
      '4': 1,
      '5': 9,
      '10': 'identityvalidationexpression'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'providerarns', '3': 301486689, '4': 3, '5': 9, '10': 'providerarns'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.apigateway.AuthorizerType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `CreateAuthorizerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAuthorizerRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVBdXRob3JpemVyUmVxdWVzdBIdCghhdXRodHlwZRjY9s5NIAEoCVIIYXV0aHR5cG'
    'USNwoVYXV0aG9yaXplcmNyZWRlbnRpYWxzGMGmsG8gASgJUhVhdXRob3JpemVyY3JlZGVudGlh'
    'bHMSRQocYXV0aG9yaXplcnJlc3VsdHR0bGluc2Vjb25kcxjQzspAIAEoBVIcYXV0aG9yaXplcn'
    'Jlc3VsdHR0bGluc2Vjb25kcxIoCg1hdXRob3JpemVydXJpGJmwtPoBIAEoCVINYXV0aG9yaXpl'
    'cnVyaRIqCg5pZGVudGl0eXNvdXJjZRj/yJiIASABKAlSDmlkZW50aXR5c291cmNlEkUKHGlkZW'
    '50aXR5dmFsaWRhdGlvbmV4cHJlc3Npb24Yv++rbCABKAlSHGlkZW50aXR5dmFsaWRhdGlvbmV4'
    'cHJlc3Npb24SFQoEbmFtZRjn++ZpIAEoCVIEbmFtZRImCgxwcm92aWRlcmFybnMY4aThjwEgAy'
    'gJUgxwcm92aWRlcmFybnMSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkEjIKBHR5'
    'cGUYzuKfiQEgASgOMhouYXBpZ2F0ZXdheS5BdXRob3JpemVyVHlwZVIEdHlwZQ==');

@$core.Deprecated('Use createBasePathMappingRequestDescriptor instead')
const CreateBasePathMappingRequest$json = {
  '1': 'CreateBasePathMappingRequest',
  '2': [
    {'1': 'basepath', '3': 267528880, '4': 1, '5': 9, '10': 'basepath'},
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stage', '3': 140155438, '4': 1, '5': 9, '10': 'stage'},
  ],
};

/// Descriptor for `CreateBasePathMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBasePathMappingRequestDescriptor = $convert.base64Decode(
    'ChxDcmVhdGVCYXNlUGF0aE1hcHBpbmdSZXF1ZXN0Eh0KCGJhc2VwYXRoGLDVyH8gASgJUghiYX'
    'NlcGF0aBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG9tYWlubmFtZRImCgxkb21haW5uYW1l'
    'aWQYqPycjgEgASgJUgxkb21haW5uYW1laWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdG'
    'FwaWlkEhcKBXN0YWdlGK606kIgASgJUgVzdGFnZQ==');

@$core.Deprecated('Use createDeploymentRequestDescriptor instead')
const CreateDeploymentRequest$json = {
  '1': 'CreateDeploymentRequest',
  '2': [
    {
      '1': 'cacheclusterenabled',
      '3': 63967991,
      '4': 1,
      '5': 8,
      '10': 'cacheclusterenabled'
    },
    {
      '1': 'cacheclustersize',
      '3': 232189861,
      '4': 1,
      '5': 14,
      '6': '.apigateway.CacheClusterSize',
      '10': 'cacheclustersize'
    },
    {
      '1': 'canarysettings',
      '3': 285544261,
      '4': 1,
      '5': 11,
      '6': '.apigateway.DeploymentCanarySettings',
      '10': 'canarysettings'
    },
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'stagedescription',
      '3': 496169986,
      '4': 1,
      '5': 9,
      '10': 'stagedescription'
    },
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
    {
      '1': 'tracingenabled',
      '3': 390995731,
      '4': 1,
      '5': 8,
      '10': 'tracingenabled'
    },
    {
      '1': 'variables',
      '3': 162226883,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateDeploymentRequest.VariablesEntry',
      '10': 'variables'
    },
  ],
  '3': [CreateDeploymentRequest_VariablesEntry$json],
};

@$core.Deprecated('Use createDeploymentRequestDescriptor instead')
const CreateDeploymentRequest_VariablesEntry$json = {
  '1': 'VariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDeploymentRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVEZXBsb3ltZW50UmVxdWVzdBIzChNjYWNoZWNsdXN0ZXJlbmFibGVkGPelwB4gAS'
    'gIUhNjYWNoZWNsdXN0ZXJlbmFibGVkEksKEGNhY2hlY2x1c3RlcnNpemUYpd/bbiABKA4yHC5h'
    'cGlnYXRld2F5LkNhY2hlQ2x1c3RlclNpemVSEGNhY2hlY2x1c3RlcnNpemUSUAoOY2FuYXJ5c2'
    'V0dGluZ3MYxZ6UiAEgASgLMiQuYXBpZ2F0ZXdheS5EZXBsb3ltZW50Q2FuYXJ5U2V0dGluZ3NS'
    'DmNhbmFyeXNldHRpbmdzEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SIA'
    'oJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkEi4KEHN0YWdlZGVzY3JpcHRpb24YgujL'
    '7AEgASgJUhBzdGFnZWRlc2NyaXB0aW9uEh8KCXN0YWdlbmFtZRiP3McEIAEoCVIJc3RhZ2VuYW'
    '1lEioKDnRyYWNpbmdlbmFibGVkGJO+uLoBIAEoCFIOdHJhY2luZ2VuYWJsZWQSUwoJdmFyaWFi'
    'bGVzGMPFrU0gAygLMjIuYXBpZ2F0ZXdheS5DcmVhdGVEZXBsb3ltZW50UmVxdWVzdC5WYXJpYW'
    'JsZXNFbnRyeVIJdmFyaWFibGVzGjwKDlZhcmlhYmxlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5'
    'EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createDocumentationPartRequestDescriptor instead')
const CreateDocumentationPartRequest$json = {
  '1': 'CreateDocumentationPartRequest',
  '2': [
    {
      '1': 'location',
      '3': 200649127,
      '4': 1,
      '5': 11,
      '6': '.apigateway.DocumentationPartLocation',
      '10': 'location'
    },
    {'1': 'properties', '3': 299789533, '4': 1, '5': 9, '10': 'properties'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `CreateDocumentationPartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDocumentationPartRequestDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVEb2N1bWVudGF0aW9uUGFydFJlcXVlc3QSRAoIbG9jYXRpb24Yp9PWXyABKAsyJS'
        '5hcGlnYXRld2F5LkRvY3VtZW50YXRpb25QYXJ0TG9jYXRpb25SCGxvY2F0aW9uEiIKCnByb3Bl'
        'cnRpZXMY3dn5jgEgASgJUgpwcm9wZXJ0aWVzEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3'
        'RhcGlpZA==');

@$core.Deprecated('Use createDocumentationVersionRequestDescriptor instead')
const CreateDocumentationVersionRequest$json = {
  '1': 'CreateDocumentationVersionRequest',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'documentationversion',
      '3': 167009804,
      '4': 1,
      '5': 9,
      '10': 'documentationversion'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
};

/// Descriptor for `CreateDocumentationVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDocumentationVersionRequestDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVEb2N1bWVudGF0aW9uVmVyc2lvblJlcXVlc3QSJAoLZGVzY3JpcHRpb24Y6va8ow'
        'EgASgJUgtkZXNjcmlwdGlvbhI1ChRkb2N1bWVudGF0aW9udmVyc2lvbhiMvNFPIAEoCVIUZG9j'
        'dW1lbnRhdGlvbnZlcnNpb24SIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkEh8KCX'
        'N0YWdlbmFtZRiP3McEIAEoCVIJc3RhZ2VuYW1l');

@$core.Deprecated(
    'Use createDomainNameAccessAssociationRequestDescriptor instead')
const CreateDomainNameAccessAssociationRequest$json = {
  '1': 'CreateDomainNameAccessAssociationRequest',
  '2': [
    {
      '1': 'accessassociationsource',
      '3': 328257828,
      '4': 1,
      '5': 9,
      '10': 'accessassociationsource'
    },
    {
      '1': 'accessassociationsourcetype',
      '3': 176397628,
      '4': 1,
      '5': 14,
      '6': '.apigateway.AccessAssociationSourceType',
      '10': 'accessassociationsourcetype'
    },
    {
      '1': 'domainnamearn',
      '3': 244019094,
      '4': 1,
      '5': 9,
      '10': 'domainnamearn'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateDomainNameAccessAssociationRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [CreateDomainNameAccessAssociationRequest_TagsEntry$json],
};

@$core.Deprecated(
    'Use createDomainNameAccessAssociationRequestDescriptor instead')
const CreateDomainNameAccessAssociationRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateDomainNameAccessAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDomainNameAccessAssociationRequestDescriptor = $convert.base64Decode(
    'CihDcmVhdGVEb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb25SZXF1ZXN0EjwKF2FjY2Vzc2Fzc2'
    '9jaWF0aW9uc291cmNlGKSiw5wBIAEoCVIXYWNjZXNzYXNzb2NpYXRpb25zb3VyY2USbAobYWNj'
    'ZXNzYXNzb2NpYXRpb25zb3VyY2V0eXBlGLy6jlQgASgOMicuYXBpZ2F0ZXdheS5BY2Nlc3NBc3'
    'NvY2lhdGlvblNvdXJjZVR5cGVSG2FjY2Vzc2Fzc29jaWF0aW9uc291cmNldHlwZRInCg1kb21h'
    'aW5uYW1lYXJuGJbfrXQgASgJUg1kb21haW5uYW1lYXJuElYKBHRhZ3MYodfboAEgAygLMj4uYX'
    'BpZ2F0ZXdheS5DcmVhdGVEb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb25SZXF1ZXN0LlRhZ3NF'
    'bnRyeVIEdGFncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKA'
    'lSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createDomainNameRequestDescriptor instead')
const CreateDomainNameRequest$json = {
  '1': 'CreateDomainNameRequest',
  '2': [
    {
      '1': 'certificatearn',
      '3': 425831704,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'certificatebody',
      '3': 183933821,
      '4': 1,
      '5': 9,
      '10': 'certificatebody'
    },
    {
      '1': 'certificatechain',
      '3': 50547450,
      '4': 1,
      '5': 9,
      '10': 'certificatechain'
    },
    {
      '1': 'certificatename',
      '3': 140276948,
      '4': 1,
      '5': 9,
      '10': 'certificatename'
    },
    {
      '1': 'certificateprivatekey',
      '3': 277544067,
      '4': 1,
      '5': 9,
      '10': 'certificateprivatekey'
    },
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'endpointaccessmode',
      '3': 356705630,
      '4': 1,
      '5': 14,
      '6': '.apigateway.EndpointAccessMode',
      '10': 'endpointaccessmode'
    },
    {
      '1': 'endpointconfiguration',
      '3': 487543735,
      '4': 1,
      '5': 11,
      '6': '.apigateway.EndpointConfiguration',
      '10': 'endpointconfiguration'
    },
    {
      '1': 'mutualtlsauthentication',
      '3': 99462043,
      '4': 1,
      '5': 11,
      '6': '.apigateway.MutualTlsAuthenticationInput',
      '10': 'mutualtlsauthentication'
    },
    {
      '1': 'ownershipverificationcertificatearn',
      '3': 100550548,
      '4': 1,
      '5': 9,
      '10': 'ownershipverificationcertificatearn'
    },
    {'1': 'policy', '3': 247528064, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'regionalcertificatearn',
      '3': 137066579,
      '4': 1,
      '5': 9,
      '10': 'regionalcertificatearn'
    },
    {
      '1': 'regionalcertificatename',
      '3': 431716941,
      '4': 1,
      '5': 9,
      '10': 'regionalcertificatename'
    },
    {
      '1': 'routingmode',
      '3': 506342119,
      '4': 1,
      '5': 14,
      '6': '.apigateway.RoutingMode',
      '10': 'routingmode'
    },
    {
      '1': 'securitypolicy',
      '3': 491792990,
      '4': 1,
      '5': 14,
      '6': '.apigateway.SecurityPolicy',
      '10': 'securitypolicy'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateDomainNameRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [CreateDomainNameRequest_TagsEntry$json],
};

@$core.Deprecated('Use createDomainNameRequestDescriptor instead')
const CreateDomainNameRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateDomainNameRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDomainNameRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVEb21haW5OYW1lUmVxdWVzdBIqCg5jZXJ0aWZpY2F0ZWFybhiY2obLASABKAlSDm'
    'NlcnRpZmljYXRlYXJuEisKD2NlcnRpZmljYXRlYm9keRj9ttpXIAEoCVIPY2VydGlmaWNhdGVi'
    'b2R5Ei0KEGNlcnRpZmljYXRlY2hhaW4Y+pWNGCABKAlSEGNlcnRpZmljYXRlY2hhaW4SKwoPY2'
    'VydGlmaWNhdGVuYW1lGNTp8UIgASgJUg9jZXJ0aWZpY2F0ZW5hbWUSOAoVY2VydGlmaWNhdGVw'
    'cml2YXRla2V5GIP5q4QBIAEoCVIVY2VydGlmaWNhdGVwcml2YXRla2V5EiIKCmRvbWFpbm5hbW'
    'UYi9OPugEgASgJUgpkb21haW5uYW1lElIKEmVuZHBvaW50YWNjZXNzbW9kZRjeyouqASABKA4y'
    'Hi5hcGlnYXRld2F5LkVuZHBvaW50QWNjZXNzTW9kZVISZW5kcG9pbnRhY2Nlc3Ntb2RlElsKFW'
    'VuZHBvaW50Y29uZmlndXJhdGlvbhi3p73oASABKAsyIS5hcGlnYXRld2F5LkVuZHBvaW50Q29u'
    'ZmlndXJhdGlvblIVZW5kcG9pbnRjb25maWd1cmF0aW9uEmUKF211dHVhbHRsc2F1dGhlbnRpY2'
    'F0aW9uGJvXti8gASgLMiguYXBpZ2F0ZXdheS5NdXR1YWxUbHNBdXRoZW50aWNhdGlvbklucHV0'
    'UhdtdXR1YWx0bHNhdXRoZW50aWNhdGlvbhJTCiNvd25lcnNoaXB2ZXJpZmljYXRpb25jZXJ0aW'
    'ZpY2F0ZWFybhiUj/kvIAEoCVIjb3duZXJzaGlwdmVyaWZpY2F0aW9uY2VydGlmaWNhdGVhcm4S'
    'GQoGcG9saWN5GID1g3YgASgJUgZwb2xpY3kSOQoWcmVnaW9uYWxjZXJ0aWZpY2F0ZWFybhjT8K'
    '1BIAEoCVIWcmVnaW9uYWxjZXJ0aWZpY2F0ZWFybhI8ChdyZWdpb25hbGNlcnRpZmljYXRlbmFt'
    'ZRjN9O3NASABKAlSF3JlZ2lvbmFsY2VydGlmaWNhdGVuYW1lEj0KC3JvdXRpbmdtb2RlGOfVuP'
    'EBIAEoDjIXLmFwaWdhdGV3YXkuUm91dGluZ01vZGVSC3JvdXRpbmdtb2RlEkYKDnNlY3VyaXR5'
    'cG9saWN5GN7UwOoBIAEoDjIaLmFwaWdhdGV3YXkuU2VjdXJpdHlQb2xpY3lSDnNlY3VyaXR5cG'
    '9saWN5EkUKBHRhZ3MYodfboAEgAygLMi0uYXBpZ2F0ZXdheS5DcmVhdGVEb21haW5OYW1lUmVx'
    'dWVzdC5UYWdzRW50cnlSBHRhZ3MaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBX'
    'ZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createModelRequestDescriptor instead')
const CreateModelRequest$json = {
  '1': 'CreateModelRequest',
  '2': [
    {'1': 'contenttype', '3': 281764659, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'schema', '3': 310182711, '4': 1, '5': 9, '10': 'schema'},
  ],
};

/// Descriptor for `CreateModelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createModelRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVNb2RlbFJlcXVlc3QSJAoLY29udGVudHR5cGUYs8athgEgASgJUgtjb250ZW50dH'
    'lwZRIkCgtkZXNjcmlwdGlvbhjq9ryjASABKAlSC2Rlc2NyaXB0aW9uEhUKBG5hbWUY5/vmaSAB'
    'KAlSBG5hbWUSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkEhoKBnNjaGVtYRi3hv'
    'STASABKAlSBnNjaGVtYQ==');

@$core.Deprecated('Use createRequestValidatorRequestDescriptor instead')
const CreateRequestValidatorRequest$json = {
  '1': 'CreateRequestValidatorRequest',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'validaterequestbody',
      '3': 397505841,
      '4': 1,
      '5': 8,
      '10': 'validaterequestbody'
    },
    {
      '1': 'validaterequestparameters',
      '3': 464035801,
      '4': 1,
      '5': 8,
      '10': 'validaterequestparameters'
    },
  ],
};

/// Descriptor for `CreateRequestValidatorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRequestValidatorRequestDescriptor = $convert.base64Decode(
    'Ch1DcmVhdGVSZXF1ZXN0VmFsaWRhdG9yUmVxdWVzdBIVCgRuYW1lGOf75mkgASgJUgRuYW1lEi'
    'AKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZBI0ChN2YWxpZGF0ZXJlcXVlc3Rib2R5'
    'GLHqxb0BIAEoCFITdmFsaWRhdGVyZXF1ZXN0Ym9keRJAChl2YWxpZGF0ZXJlcXVlc3RwYXJhbW'
    'V0ZXJzGNm/ot0BIAEoCFIZdmFsaWRhdGVyZXF1ZXN0cGFyYW1ldGVycw==');

@$core.Deprecated('Use createResourceRequestDescriptor instead')
const CreateResourceRequest$json = {
  '1': 'CreateResourceRequest',
  '2': [
    {'1': 'parentid', '3': 106050857, '4': 1, '5': 9, '10': 'parentid'},
    {'1': 'pathpart', '3': 487915984, '4': 1, '5': 9, '10': 'pathpart'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `CreateResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResourceRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVSZXNvdXJjZVJlcXVlc3QSHQoIcGFyZW50aWQYqerIMiABKAlSCHBhcmVudGlkEh'
    '4KCHBhdGhwYXJ0GNCD1OgBIAEoCVIIcGF0aHBhcnQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJ'
    'cmVzdGFwaWlk');

@$core.Deprecated('Use createRestApiRequestDescriptor instead')
const CreateRestApiRequest$json = {
  '1': 'CreateRestApiRequest',
  '2': [
    {
      '1': 'apikeysource',
      '3': 108531220,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ApiKeySourceType',
      '10': 'apikeysource'
    },
    {
      '1': 'binarymediatypes',
      '3': 406416146,
      '4': 3,
      '5': 9,
      '10': 'binarymediatypes'
    },
    {'1': 'clonefrom', '3': 263376551, '4': 1, '5': 9, '10': 'clonefrom'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'disableexecuteapiendpoint',
      '3': 148140696,
      '4': 1,
      '5': 8,
      '10': 'disableexecuteapiendpoint'
    },
    {
      '1': 'endpointaccessmode',
      '3': 356705630,
      '4': 1,
      '5': 14,
      '6': '.apigateway.EndpointAccessMode',
      '10': 'endpointaccessmode'
    },
    {
      '1': 'endpointconfiguration',
      '3': 487543735,
      '4': 1,
      '5': 11,
      '6': '.apigateway.EndpointConfiguration',
      '10': 'endpointconfiguration'
    },
    {
      '1': 'minimumcompressionsize',
      '3': 254902719,
      '4': 1,
      '5': 5,
      '10': 'minimumcompressionsize'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'policy', '3': 247528064, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'securitypolicy',
      '3': 491792990,
      '4': 1,
      '5': 14,
      '6': '.apigateway.SecurityPolicy',
      '10': 'securitypolicy'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateRestApiRequest.TagsEntry',
      '10': 'tags'
    },
    {'1': 'version', '3': 108113560, '4': 1, '5': 9, '10': 'version'},
  ],
  '3': [CreateRestApiRequest_TagsEntry$json],
};

@$core.Deprecated('Use createRestApiRequestDescriptor instead')
const CreateRestApiRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateRestApiRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRestApiRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVSZXN0QXBpUmVxdWVzdBJDCgxhcGlrZXlzb3VyY2UYlJzgMyABKA4yHC5hcGlnYX'
    'Rld2F5LkFwaUtleVNvdXJjZVR5cGVSDGFwaWtleXNvdXJjZRIuChBiaW5hcnltZWRpYXR5cGVz'
    'GJLW5cEBIAMoCVIQYmluYXJ5bWVkaWF0eXBlcxIfCgljbG9uZWZyb20Yp53LfSABKAlSCWNsb2'
    '5lZnJvbRIkCgtkZXNjcmlwdGlvbhjq9ryjASABKAlSC2Rlc2NyaXB0aW9uEj8KGWRpc2FibGVl'
    'eGVjdXRlYXBpZW5kcG9pbnQYmOXRRiABKAhSGWRpc2FibGVleGVjdXRlYXBpZW5kcG9pbnQSUg'
    'oSZW5kcG9pbnRhY2Nlc3Ntb2RlGN7Ki6oBIAEoDjIeLmFwaWdhdGV3YXkuRW5kcG9pbnRBY2Nl'
    'c3NNb2RlUhJlbmRwb2ludGFjY2Vzc21vZGUSWwoVZW5kcG9pbnRjb25maWd1cmF0aW9uGLenve'
    'gBIAEoCzIhLmFwaWdhdGV3YXkuRW5kcG9pbnRDb25maWd1cmF0aW9uUhVlbmRwb2ludGNvbmZp'
    'Z3VyYXRpb24SOQoWbWluaW11bWNvbXByZXNzaW9uc2l6ZRi/g8Z5IAEoBVIWbWluaW11bWNvbX'
    'ByZXNzaW9uc2l6ZRIVCgRuYW1lGOf75mkgASgJUgRuYW1lEhkKBnBvbGljeRiA9YN2IAEoCVIG'
    'cG9saWN5EkYKDnNlY3VyaXR5cG9saWN5GN7UwOoBIAEoDjIaLmFwaWdhdGV3YXkuU2VjdXJpdH'
    'lQb2xpY3lSDnNlY3VyaXR5cG9saWN5EkIKBHRhZ3MYodfboAEgAygLMiouYXBpZ2F0ZXdheS5D'
    'cmVhdGVSZXN0QXBpUmVxdWVzdC5UYWdzRW50cnlSBHRhZ3MSGwoHdmVyc2lvbhiY3cYzIAEoCV'
    'IHdmVyc2lvbho3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createStageRequestDescriptor instead')
const CreateStageRequest$json = {
  '1': 'CreateStageRequest',
  '2': [
    {
      '1': 'cacheclusterenabled',
      '3': 63967991,
      '4': 1,
      '5': 8,
      '10': 'cacheclusterenabled'
    },
    {
      '1': 'cacheclustersize',
      '3': 232189861,
      '4': 1,
      '5': 14,
      '6': '.apigateway.CacheClusterSize',
      '10': 'cacheclustersize'
    },
    {
      '1': 'canarysettings',
      '3': 285544261,
      '4': 1,
      '5': 11,
      '6': '.apigateway.CanarySettings',
      '10': 'canarysettings'
    },
    {'1': 'deploymentid', '3': 439369188, '4': 1, '5': 9, '10': 'deploymentid'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'documentationversion',
      '3': 167009804,
      '4': 1,
      '5': 9,
      '10': 'documentationversion'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateStageRequest.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'tracingenabled',
      '3': 390995731,
      '4': 1,
      '5': 8,
      '10': 'tracingenabled'
    },
    {
      '1': 'variables',
      '3': 162226883,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateStageRequest.VariablesEntry',
      '10': 'variables'
    },
  ],
  '3': [
    CreateStageRequest_TagsEntry$json,
    CreateStageRequest_VariablesEntry$json
  ],
};

@$core.Deprecated('Use createStageRequestDescriptor instead')
const CreateStageRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use createStageRequestDescriptor instead')
const CreateStageRequest_VariablesEntry$json = {
  '1': 'VariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateStageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStageRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVTdGFnZVJlcXVlc3QSMwoTY2FjaGVjbHVzdGVyZW5hYmxlZBj3pcAeIAEoCFITY2'
    'FjaGVjbHVzdGVyZW5hYmxlZBJLChBjYWNoZWNsdXN0ZXJzaXplGKXf224gASgOMhwuYXBpZ2F0'
    'ZXdheS5DYWNoZUNsdXN0ZXJTaXplUhBjYWNoZWNsdXN0ZXJzaXplEkYKDmNhbmFyeXNldHRpbm'
    'dzGMWelIgBIAEoCzIaLmFwaWdhdGV3YXkuQ2FuYXJ5U2V0dGluZ3NSDmNhbmFyeXNldHRpbmdz'
    'EiYKDGRlcGxveW1lbnRpZBjk+8DRASABKAlSDGRlcGxveW1lbnRpZBIkCgtkZXNjcmlwdGlvbh'
    'jq9ryjASABKAlSC2Rlc2NyaXB0aW9uEjUKFGRvY3VtZW50YXRpb252ZXJzaW9uGIy80U8gASgJ'
    'UhRkb2N1bWVudGF0aW9udmVyc2lvbhIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaW'
    'QSHwoJc3RhZ2VuYW1lGI/cxwQgASgJUglzdGFnZW5hbWUSQAoEdGFncxih19ugASADKAsyKC5h'
    'cGlnYXRld2F5LkNyZWF0ZVN0YWdlUmVxdWVzdC5UYWdzRW50cnlSBHRhZ3MSKgoOdHJhY2luZ2'
    'VuYWJsZWQYk764ugEgASgIUg50cmFjaW5nZW5hYmxlZBJOCgl2YXJpYWJsZXMYw8WtTSADKAsy'
    'LS5hcGlnYXRld2F5LkNyZWF0ZVN0YWdlUmVxdWVzdC5WYXJpYWJsZXNFbnRyeVIJdmFyaWFibG'
    'VzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6'
    'AjgBGjwKDlZhcmlhYmxlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUg'
    'V2YWx1ZToCOAE=');

@$core.Deprecated('Use createUsagePlanKeyRequestDescriptor instead')
const CreateUsagePlanKeyRequest$json = {
  '1': 'CreateUsagePlanKeyRequest',
  '2': [
    {'1': 'keyid', '3': 479913282, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'keytype', '3': 107017349, '4': 1, '5': 9, '10': 'keytype'},
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `CreateUsagePlanKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUsagePlanKeyRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVVc2FnZVBsYW5LZXlSZXF1ZXN0EhgKBWtleWlkGMLK6+QBIAEoCVIFa2V5aWQSGw'
    'oHa2V5dHlwZRiF6YMzIAEoCVIHa2V5dHlwZRIkCgt1c2FnZXBsYW5pZBjX8OXyASABKAlSC3Vz'
    'YWdlcGxhbmlk');

@$core.Deprecated('Use createUsagePlanRequestDescriptor instead')
const CreateUsagePlanRequest$json = {
  '1': 'CreateUsagePlanRequest',
  '2': [
    {
      '1': 'apistages',
      '3': 64558449,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ApiStage',
      '10': 'apistages'
    },
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'quota',
      '3': 243824012,
      '4': 1,
      '5': 11,
      '6': '.apigateway.QuotaSettings',
      '10': 'quota'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateUsagePlanRequest.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'throttle',
      '3': 395260638,
      '4': 1,
      '5': 11,
      '6': '.apigateway.ThrottleSettings',
      '10': 'throttle'
    },
  ],
  '3': [CreateUsagePlanRequest_TagsEntry$json],
};

@$core.Deprecated('Use createUsagePlanRequestDescriptor instead')
const CreateUsagePlanRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateUsagePlanRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUsagePlanRequestDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVVc2FnZVBsYW5SZXF1ZXN0EjUKCWFwaXN0YWdlcxjxquQeIAMoCzIULmFwaWdhdG'
    'V3YXkuQXBpU3RhZ2VSCWFwaXN0YWdlcxIkCgtkZXNjcmlwdGlvbhjq9ryjASABKAlSC2Rlc2Ny'
    'aXB0aW9uEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSMgoFcXVvdGEYjOuhdCABKAsyGS5hcGlnYX'
    'Rld2F5LlF1b3RhU2V0dGluZ3NSBXF1b3RhEkQKBHRhZ3MYodfboAEgAygLMiwuYXBpZ2F0ZXdh'
    'eS5DcmVhdGVVc2FnZVBsYW5SZXF1ZXN0LlRhZ3NFbnRyeVIEdGFncxI8Cgh0aHJvdHRsZRje5b'
    'y8ASABKAsyHC5hcGlnYXRld2F5LlRocm90dGxlU2V0dGluZ3NSCHRocm90dGxlGjcKCVRhZ3NF'
    'bnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use createVpcLinkRequestDescriptor instead')
const CreateVpcLinkRequest$json = {
  '1': 'CreateVpcLinkRequest',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.CreateVpcLinkRequest.TagsEntry',
      '10': 'tags'
    },
    {'1': 'targetarns', '3': 46319317, '4': 3, '5': 9, '10': 'targetarns'},
  ],
  '3': [CreateVpcLinkRequest_TagsEntry$json],
};

@$core.Deprecated('Use createVpcLinkRequestDescriptor instead')
const CreateVpcLinkRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateVpcLinkRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createVpcLinkRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVWcGNMaW5rUmVxdWVzdBIkCgtkZXNjcmlwdGlvbhjq9ryjASABKAlSC2Rlc2NyaX'
    'B0aW9uEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSQgoEdGFncxih19ugASADKAsyKi5hcGlnYXRl'
    'd2F5LkNyZWF0ZVZwY0xpbmtSZXF1ZXN0LlRhZ3NFbnRyeVIEdGFncxIhCgp0YXJnZXRhcm5zGN'
    'WNixYgAygJUgp0YXJnZXRhcm5zGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2'
    'YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use deleteApiKeyRequestDescriptor instead')
const DeleteApiKeyRequest$json = {
  '1': 'DeleteApiKeyRequest',
  '2': [
    {'1': 'apikey', '3': 490490655, '4': 1, '5': 9, '10': 'apikey'},
  ],
};

/// Descriptor for `DeleteApiKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteApiKeyRequestDescriptor =
    $convert.base64Decode(
        'ChNEZWxldGVBcGlLZXlSZXF1ZXN0EhoKBmFwaWtleRiflvHpASABKAlSBmFwaWtleQ==');

@$core.Deprecated('Use deleteAuthorizerRequestDescriptor instead')
const DeleteAuthorizerRequest$json = {
  '1': 'DeleteAuthorizerRequest',
  '2': [
    {'1': 'authorizerid', '3': 111773148, '4': 1, '5': 9, '10': 'authorizerid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteAuthorizerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAuthorizerRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVBdXRob3JpemVyUmVxdWVzdBIlCgxhdXRob3JpemVyaWQY3IumNSABKAlSDGF1dG'
        'hvcml6ZXJpZBIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use deleteBasePathMappingRequestDescriptor instead')
const DeleteBasePathMappingRequest$json = {
  '1': 'DeleteBasePathMappingRequest',
  '2': [
    {'1': 'basepath', '3': 267528880, '4': 1, '5': 9, '10': 'basepath'},
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
  ],
};

/// Descriptor for `DeleteBasePathMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBasePathMappingRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVCYXNlUGF0aE1hcHBpbmdSZXF1ZXN0Eh0KCGJhc2VwYXRoGLDVyH8gASgJUghiYX'
        'NlcGF0aBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG9tYWlubmFtZRImCgxkb21haW5uYW1l'
        'aWQYqPycjgEgASgJUgxkb21haW5uYW1laWQ=');

@$core.Deprecated('Use deleteClientCertificateRequestDescriptor instead')
const DeleteClientCertificateRequest$json = {
  '1': 'DeleteClientCertificateRequest',
  '2': [
    {
      '1': 'clientcertificateid',
      '3': 276222909,
      '4': 1,
      '5': 9,
      '10': 'clientcertificateid'
    },
  ],
};

/// Descriptor for `DeleteClientCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteClientCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVDbGllbnRDZXJ0aWZpY2F0ZVJlcXVlc3QSNAoTY2xpZW50Y2VydGlmaWNhdGVpZB'
        'i9p9uDASABKAlSE2NsaWVudGNlcnRpZmljYXRlaWQ=');

@$core.Deprecated('Use deleteDeploymentRequestDescriptor instead')
const DeleteDeploymentRequest$json = {
  '1': 'DeleteDeploymentRequest',
  '2': [
    {'1': 'deploymentid', '3': 439369188, '4': 1, '5': 9, '10': 'deploymentid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDeploymentRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVEZXBsb3ltZW50UmVxdWVzdBImCgxkZXBsb3ltZW50aWQY5PvA0QEgASgJUgxkZX'
        'Bsb3ltZW50aWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use deleteDocumentationPartRequestDescriptor instead')
const DeleteDocumentationPartRequest$json = {
  '1': 'DeleteDocumentationPartRequest',
  '2': [
    {
      '1': 'documentationpartid',
      '3': 286552774,
      '4': 1,
      '5': 9,
      '10': 'documentationpartid'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteDocumentationPartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDocumentationPartRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVEb2N1bWVudGF0aW9uUGFydFJlcXVlc3QSNAoTZG9jdW1lbnRhdGlvbnBhcnRpZB'
        'jG5dGIASABKAlSE2RvY3VtZW50YXRpb25wYXJ0aWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJ'
        'cmVzdGFwaWlk');

@$core.Deprecated('Use deleteDocumentationVersionRequestDescriptor instead')
const DeleteDocumentationVersionRequest$json = {
  '1': 'DeleteDocumentationVersionRequest',
  '2': [
    {
      '1': 'documentationversion',
      '3': 167009804,
      '4': 1,
      '5': 9,
      '10': 'documentationversion'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteDocumentationVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDocumentationVersionRequestDescriptor =
    $convert.base64Decode(
        'CiFEZWxldGVEb2N1bWVudGF0aW9uVmVyc2lvblJlcXVlc3QSNQoUZG9jdW1lbnRhdGlvbnZlcn'
        'Npb24YjLzRTyABKAlSFGRvY3VtZW50YXRpb252ZXJzaW9uEiAKCXJlc3RhcGlpZBiZpIG3ASAB'
        'KAlSCXJlc3RhcGlpZA==');

@$core.Deprecated(
    'Use deleteDomainNameAccessAssociationRequestDescriptor instead')
const DeleteDomainNameAccessAssociationRequest$json = {
  '1': 'DeleteDomainNameAccessAssociationRequest',
  '2': [
    {
      '1': 'domainnameaccessassociationarn',
      '3': 281017927,
      '4': 1,
      '5': 9,
      '10': 'domainnameaccessassociationarn'
    },
  ],
};

/// Descriptor for `DeleteDomainNameAccessAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDomainNameAccessAssociationRequestDescriptor =
    $convert.base64Decode(
        'CihEZWxldGVEb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb25SZXF1ZXN0EkoKHmRvbWFpbm5hbW'
        'VhY2Nlc3Nhc3NvY2lhdGlvbmFybhjH/P+FASABKAlSHmRvbWFpbm5hbWVhY2Nlc3Nhc3NvY2lh'
        'dGlvbmFybg==');

@$core.Deprecated('Use deleteDomainNameRequestDescriptor instead')
const DeleteDomainNameRequest$json = {
  '1': 'DeleteDomainNameRequest',
  '2': [
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
  ],
};

/// Descriptor for `DeleteDomainNameRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDomainNameRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVEb21haW5OYW1lUmVxdWVzdBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG9tYW'
        'lubmFtZRImCgxkb21haW5uYW1laWQYqPycjgEgASgJUgxkb21haW5uYW1laWQ=');

@$core.Deprecated('Use deleteGatewayResponseRequestDescriptor instead')
const DeleteGatewayResponseRequest$json = {
  '1': 'DeleteGatewayResponseRequest',
  '2': [
    {
      '1': 'responsetype',
      '3': 377935935,
      '4': 1,
      '5': 14,
      '6': '.apigateway.GatewayResponseType',
      '10': 'responsetype'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteGatewayResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteGatewayResponseRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVHYXRld2F5UmVzcG9uc2VSZXF1ZXN0EkcKDHJlc3BvbnNldHlwZRi/sJu0ASABKA'
        '4yHy5hcGlnYXRld2F5LkdhdGV3YXlSZXNwb25zZVR5cGVSDHJlc3BvbnNldHlwZRIgCglyZXN0'
        'YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use deleteIntegrationRequestDescriptor instead')
const DeleteIntegrationRequest$json = {
  '1': 'DeleteIntegrationRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIntegrationRequestDescriptor = $convert.base64Decode(
    'ChhEZWxldGVJbnRlZ3JhdGlvblJlcXVlc3QSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaHR0cG'
    '1ldGhvZBIiCgpyZXNvdXJjZWlkGLG9iZgBIAEoCVIKcmVzb3VyY2VpZBIgCglyZXN0YXBpaWQY'
    'maSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use deleteIntegrationResponseRequestDescriptor instead')
const DeleteIntegrationResponseRequest$json = {
  '1': 'DeleteIntegrationResponseRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
};

/// Descriptor for `DeleteIntegrationResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIntegrationResponseRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVJbnRlZ3JhdGlvblJlc3BvbnNlUmVxdWVzdBIhCgpodHRwbWV0aG9kGPHz+zYgAS'
        'gJUgpodHRwbWV0aG9kEiIKCnJlc291cmNlaWQYsb2JmAEgASgJUgpyZXNvdXJjZWlkEiAKCXJl'
        'c3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZBIiCgpzdGF0dXNjb2RlGJ+B344BIAEoCVIKc3'
        'RhdHVzY29kZQ==');

@$core.Deprecated('Use deleteMethodRequestDescriptor instead')
const DeleteMethodRequest$json = {
  '1': 'DeleteMethodRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteMethodRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMethodRequestDescriptor = $convert.base64Decode(
    'ChNEZWxldGVNZXRob2RSZXF1ZXN0EiEKCmh0dHBtZXRob2QY8fP7NiABKAlSCmh0dHBtZXRob2'
    'QSIgoKcmVzb3VyY2VpZBixvYmYASABKAlSCnJlc291cmNlaWQSIAoJcmVzdGFwaWlkGJmkgbcB'
    'IAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use deleteMethodResponseRequestDescriptor instead')
const DeleteMethodResponseRequest$json = {
  '1': 'DeleteMethodResponseRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
};

/// Descriptor for `DeleteMethodResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMethodResponseRequestDescriptor = $convert.base64Decode(
    'ChtEZWxldGVNZXRob2RSZXNwb25zZVJlcXVlc3QSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaH'
    'R0cG1ldGhvZBIiCgpyZXNvdXJjZWlkGLG9iZgBIAEoCVIKcmVzb3VyY2VpZBIgCglyZXN0YXBp'
    'aWQYmaSBtwEgASgJUglyZXN0YXBpaWQSIgoKc3RhdHVzY29kZRifgd+OASABKAlSCnN0YXR1c2'
    'NvZGU=');

@$core.Deprecated('Use deleteModelRequestDescriptor instead')
const DeleteModelRequest$json = {
  '1': 'DeleteModelRequest',
  '2': [
    {'1': 'modelname', '3': 176835330, '4': 1, '5': 9, '10': 'modelname'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteModelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteModelRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVNb2RlbFJlcXVlc3QSHwoJbW9kZWxuYW1lGIKWqVQgASgJUgltb2RlbG5hbWUSIA'
    'oJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use deleteRequestValidatorRequestDescriptor instead')
const DeleteRequestValidatorRequest$json = {
  '1': 'DeleteRequestValidatorRequest',
  '2': [
    {
      '1': 'requestvalidatorid',
      '3': 517546134,
      '4': 1,
      '5': 9,
      '10': 'requestvalidatorid'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteRequestValidatorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRequestValidatorRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVSZXF1ZXN0VmFsaWRhdG9yUmVxdWVzdBIyChJyZXF1ZXN0dmFsaWRhdG9yaWQYls'
        'Hk9gEgASgJUhJyZXF1ZXN0dmFsaWRhdG9yaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVz'
        'dGFwaWlk');

@$core.Deprecated('Use deleteResourceRequestDescriptor instead')
const DeleteResourceRequest$json = {
  '1': 'DeleteResourceRequest',
  '2': [
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourceRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVSZXNvdXJjZVJlcXVlc3QSIgoKcmVzb3VyY2VpZBixvYmYASABKAlSCnJlc291cm'
    'NlaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use deleteRestApiRequestDescriptor instead')
const DeleteRestApiRequest$json = {
  '1': 'DeleteRestApiRequest',
  '2': [
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `DeleteRestApiRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRestApiRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVSZXN0QXBpUmVxdWVzdBIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaW'
    'Q=');

@$core.Deprecated('Use deleteStageRequestDescriptor instead')
const DeleteStageRequest$json = {
  '1': 'DeleteStageRequest',
  '2': [
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
};

/// Descriptor for `DeleteStageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStageRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVTdGFnZVJlcXVlc3QSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkEh'
    '8KCXN0YWdlbmFtZRiP3McEIAEoCVIJc3RhZ2VuYW1l');

@$core.Deprecated('Use deleteUsagePlanKeyRequestDescriptor instead')
const DeleteUsagePlanKeyRequest$json = {
  '1': 'DeleteUsagePlanKeyRequest',
  '2': [
    {'1': 'keyid', '3': 479913282, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `DeleteUsagePlanKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUsagePlanKeyRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVVc2FnZVBsYW5LZXlSZXF1ZXN0EhgKBWtleWlkGMLK6+QBIAEoCVIFa2V5aWQSJA'
        'oLdXNhZ2VwbGFuaWQY1/Dl8gEgASgJUgt1c2FnZXBsYW5pZA==');

@$core.Deprecated('Use deleteUsagePlanRequestDescriptor instead')
const DeleteUsagePlanRequest$json = {
  '1': 'DeleteUsagePlanRequest',
  '2': [
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `DeleteUsagePlanRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUsagePlanRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVVc2FnZVBsYW5SZXF1ZXN0EiQKC3VzYWdlcGxhbmlkGNfw5fIBIAEoCVILdXNhZ2'
        'VwbGFuaWQ=');

@$core.Deprecated('Use deleteVpcLinkRequestDescriptor instead')
const DeleteVpcLinkRequest$json = {
  '1': 'DeleteVpcLinkRequest',
  '2': [
    {'1': 'vpclinkid', '3': 27515846, '4': 1, '5': 9, '10': 'vpclinkid'},
  ],
};

/// Descriptor for `DeleteVpcLinkRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteVpcLinkRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVWcGNMaW5rUmVxdWVzdBIfCgl2cGNsaW5raWQYxrePDSABKAlSCXZwY2xpbmtpZA'
    '==');

@$core.Deprecated('Use deploymentDescriptor instead')
const Deployment$json = {
  '1': 'Deployment',
  '2': [
    {
      '1': 'apisummary',
      '3': 159675170,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Deployment.ApisummaryEntry',
      '10': 'apisummary'
    },
    {'1': 'createddate', '3': 53061200, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
  ],
  '3': [Deployment_ApisummaryEntry$json],
};

@$core.Deprecated('Use deploymentDescriptor instead')
const Deployment_ApisummaryEntry$json = {
  '1': 'ApisummaryEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Deployment`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deploymentDescriptor = $convert.base64Decode(
    'CgpEZXBsb3ltZW50EkkKCmFwaXN1bW1hcnkYouaRTCADKAsyJi5hcGlnYXRld2F5LkRlcGxveW'
    '1lbnQuQXBpc3VtbWFyeUVudHJ5UgphcGlzdW1tYXJ5EiMKC2NyZWF0ZWRkYXRlGNDMphkgASgJ'
    'UgtjcmVhdGVkZGF0ZRIkCgtkZXNjcmlwdGlvbhjq9ryjASABKAlSC2Rlc2NyaXB0aW9uEhIKAm'
    'lkGOHV4bkBIAEoCVICaWQaPQoPQXBpc3VtbWFyeUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQK'
    'BXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use deploymentCanarySettingsDescriptor instead')
const DeploymentCanarySettings$json = {
  '1': 'DeploymentCanarySettings',
  '2': [
    {
      '1': 'percenttraffic',
      '3': 147177864,
      '4': 1,
      '5': 1,
      '10': 'percenttraffic'
    },
    {
      '1': 'stagevariableoverrides',
      '3': 221124259,
      '4': 3,
      '5': 11,
      '6': '.apigateway.DeploymentCanarySettings.StagevariableoverridesEntry',
      '10': 'stagevariableoverrides'
    },
    {
      '1': 'usestagecache',
      '3': 179841697,
      '4': 1,
      '5': 8,
      '10': 'usestagecache'
    },
  ],
  '3': [DeploymentCanarySettings_StagevariableoverridesEntry$json],
};

@$core.Deprecated('Use deploymentCanarySettingsDescriptor instead')
const DeploymentCanarySettings_StagevariableoverridesEntry$json = {
  '1': 'StagevariableoverridesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DeploymentCanarySettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deploymentCanarySettingsDescriptor = $convert.base64Decode(
    'ChhEZXBsb3ltZW50Q2FuYXJ5U2V0dGluZ3MSKQoOcGVyY2VudHRyYWZmaWMYiIOXRiABKAFSDn'
    'BlcmNlbnR0cmFmZmljEnsKFnN0YWdldmFyaWFibGVvdmVycmlkZXMYo624aSADKAsyQC5hcGln'
    'YXRld2F5LkRlcGxveW1lbnRDYW5hcnlTZXR0aW5ncy5TdGFnZXZhcmlhYmxlb3ZlcnJpZGVzRW'
    '50cnlSFnN0YWdldmFyaWFibGVvdmVycmlkZXMSJwoNdXNlc3RhZ2VjYWNoZRih1eBVIAEoCFIN'
    'dXNlc3RhZ2VjYWNoZRpJChtTdGFnZXZhcmlhYmxlb3ZlcnJpZGVzRW50cnkSEAoDa2V5GAEgAS'
    'gJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use deploymentsDescriptor instead')
const Deployments$json = {
  '1': 'Deployments',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Deployment',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `Deployments`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deploymentsDescriptor = $convert.base64Decode(
    'CgtEZXBsb3ltZW50cxIwCgVpdGVtcxiQ5+TTASADKAsyFi5hcGlnYXRld2F5LkRlcGxveW1lbn'
    'RSBWl0ZW1zEh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24=');

@$core.Deprecated('Use documentationPartDescriptor instead')
const DocumentationPart$json = {
  '1': 'DocumentationPart',
  '2': [
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'location',
      '3': 200649127,
      '4': 1,
      '5': 11,
      '6': '.apigateway.DocumentationPartLocation',
      '10': 'location'
    },
    {'1': 'properties', '3': 299789533, '4': 1, '5': 9, '10': 'properties'},
  ],
};

/// Descriptor for `DocumentationPart`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentationPartDescriptor = $convert.base64Decode(
    'ChFEb2N1bWVudGF0aW9uUGFydBISCgJpZBjh1eG5ASABKAlSAmlkEkQKCGxvY2F0aW9uGKfT1l'
    '8gASgLMiUuYXBpZ2F0ZXdheS5Eb2N1bWVudGF0aW9uUGFydExvY2F0aW9uUghsb2NhdGlvbhIi'
    'Cgpwcm9wZXJ0aWVzGN3Z+Y4BIAEoCVIKcHJvcGVydGllcw==');

@$core.Deprecated('Use documentationPartIdsDescriptor instead')
const DocumentationPartIds$json = {
  '1': 'DocumentationPartIds',
  '2': [
    {'1': 'ids', '3': 13616490, '4': 3, '5': 9, '10': 'ids'},
    {'1': 'warnings', '3': 185617301, '4': 3, '5': 9, '10': 'warnings'},
  ],
};

/// Descriptor for `DocumentationPartIds`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentationPartIdsDescriptor = $convert.base64Decode(
    'ChREb2N1bWVudGF0aW9uUGFydElkcxITCgNpZHMY6oq/BiADKAlSA2lkcxIdCgh3YXJuaW5ncx'
    'iVl8FYIAMoCVIId2FybmluZ3M=');

@$core.Deprecated('Use documentationPartLocationDescriptor instead')
const DocumentationPartLocation$json = {
  '1': 'DocumentationPartLocation',
  '2': [
    {'1': 'method', '3': 189134641, '4': 1, '5': 9, '10': 'method'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'path', '3': 75975991, '4': 1, '5': 9, '10': 'path'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.apigateway.DocumentationPartType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `DocumentationPartLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentationPartLocationDescriptor = $convert.base64Decode(
    'ChlEb2N1bWVudGF0aW9uUGFydExvY2F0aW9uEhkKBm1ldGhvZBix7pdaIAEoCVIGbWV0aG9kEh'
    'UKBG5hbWUY5/vmaSABKAlSBG5hbWUSFQoEcGF0aBi3mp0kIAEoCVIEcGF0aBIiCgpzdGF0dXNj'
    'b2RlGJ+B344BIAEoCVIKc3RhdHVzY29kZRI5CgR0eXBlGM7in4kBIAEoDjIhLmFwaWdhdGV3YX'
    'kuRG9jdW1lbnRhdGlvblBhcnRUeXBlUgR0eXBl');

@$core.Deprecated('Use documentationPartsDescriptor instead')
const DocumentationParts$json = {
  '1': 'DocumentationParts',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.DocumentationPart',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `DocumentationParts`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentationPartsDescriptor = $convert.base64Decode(
    'ChJEb2N1bWVudGF0aW9uUGFydHMSNwoFaXRlbXMYkOfk0wEgAygLMh0uYXBpZ2F0ZXdheS5Eb2'
    'N1bWVudGF0aW9uUGFydFIFaXRlbXMSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbg==');

@$core.Deprecated('Use documentationVersionDescriptor instead')
const DocumentationVersion$json = {
  '1': 'DocumentationVersion',
  '2': [
    {'1': 'createddate', '3': 53061200, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'version', '3': 108113560, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `DocumentationVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentationVersionDescriptor = $convert.base64Decode(
    'ChREb2N1bWVudGF0aW9uVmVyc2lvbhIjCgtjcmVhdGVkZGF0ZRjQzKYZIAEoCVILY3JlYXRlZG'
    'RhdGUSJAoLZGVzY3JpcHRpb24Y6va8owEgASgJUgtkZXNjcmlwdGlvbhIbCgd2ZXJzaW9uGJjd'
    'xjMgASgJUgd2ZXJzaW9u');

@$core.Deprecated('Use documentationVersionsDescriptor instead')
const DocumentationVersions$json = {
  '1': 'DocumentationVersions',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.DocumentationVersion',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `DocumentationVersions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentationVersionsDescriptor = $convert.base64Decode(
    'ChVEb2N1bWVudGF0aW9uVmVyc2lvbnMSOgoFaXRlbXMYkOfk0wEgAygLMiAuYXBpZ2F0ZXdheS'
    '5Eb2N1bWVudGF0aW9uVmVyc2lvblIFaXRlbXMSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3Np'
    'dGlvbg==');

@$core.Deprecated('Use domainNameDescriptor instead')
const DomainName$json = {
  '1': 'DomainName',
  '2': [
    {
      '1': 'certificatearn',
      '3': 425831704,
      '4': 1,
      '5': 9,
      '10': 'certificatearn'
    },
    {
      '1': 'certificatename',
      '3': 140276948,
      '4': 1,
      '5': 9,
      '10': 'certificatename'
    },
    {
      '1': 'certificateuploaddate',
      '3': 504466814,
      '4': 1,
      '5': 9,
      '10': 'certificateuploaddate'
    },
    {
      '1': 'distributiondomainname',
      '3': 266229213,
      '4': 1,
      '5': 9,
      '10': 'distributiondomainname'
    },
    {
      '1': 'distributionhostedzoneid',
      '3': 185867208,
      '4': 1,
      '5': 9,
      '10': 'distributionhostedzoneid'
    },
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {
      '1': 'domainnamearn',
      '3': 244019094,
      '4': 1,
      '5': 9,
      '10': 'domainnamearn'
    },
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
    {
      '1': 'domainnamestatus',
      '3': 275080629,
      '4': 1,
      '5': 14,
      '6': '.apigateway.DomainNameStatus',
      '10': 'domainnamestatus'
    },
    {
      '1': 'domainnamestatusmessage',
      '3': 335045044,
      '4': 1,
      '5': 9,
      '10': 'domainnamestatusmessage'
    },
    {
      '1': 'endpointaccessmode',
      '3': 356705630,
      '4': 1,
      '5': 14,
      '6': '.apigateway.EndpointAccessMode',
      '10': 'endpointaccessmode'
    },
    {
      '1': 'endpointconfiguration',
      '3': 487543735,
      '4': 1,
      '5': 11,
      '6': '.apigateway.EndpointConfiguration',
      '10': 'endpointconfiguration'
    },
    {
      '1': 'managementpolicy',
      '3': 90103737,
      '4': 1,
      '5': 9,
      '10': 'managementpolicy'
    },
    {
      '1': 'mutualtlsauthentication',
      '3': 99462043,
      '4': 1,
      '5': 11,
      '6': '.apigateway.MutualTlsAuthentication',
      '10': 'mutualtlsauthentication'
    },
    {
      '1': 'ownershipverificationcertificatearn',
      '3': 100550548,
      '4': 1,
      '5': 9,
      '10': 'ownershipverificationcertificatearn'
    },
    {'1': 'policy', '3': 247528064, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'regionalcertificatearn',
      '3': 137066579,
      '4': 1,
      '5': 9,
      '10': 'regionalcertificatearn'
    },
    {
      '1': 'regionalcertificatename',
      '3': 431716941,
      '4': 1,
      '5': 9,
      '10': 'regionalcertificatename'
    },
    {
      '1': 'regionaldomainname',
      '3': 198256560,
      '4': 1,
      '5': 9,
      '10': 'regionaldomainname'
    },
    {
      '1': 'regionalhostedzoneid',
      '3': 467276949,
      '4': 1,
      '5': 9,
      '10': 'regionalhostedzoneid'
    },
    {
      '1': 'routingmode',
      '3': 506342119,
      '4': 1,
      '5': 14,
      '6': '.apigateway.RoutingMode',
      '10': 'routingmode'
    },
    {
      '1': 'securitypolicy',
      '3': 491792990,
      '4': 1,
      '5': 14,
      '6': '.apigateway.SecurityPolicy',
      '10': 'securitypolicy'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.DomainName.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [DomainName_TagsEntry$json],
};

@$core.Deprecated('Use domainNameDescriptor instead')
const DomainName_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DomainName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainNameDescriptor = $convert.base64Decode(
    'CgpEb21haW5OYW1lEioKDmNlcnRpZmljYXRlYXJuGJjahssBIAEoCVIOY2VydGlmaWNhdGVhcm'
    '4SKwoPY2VydGlmaWNhdGVuYW1lGNTp8UIgASgJUg9jZXJ0aWZpY2F0ZW5hbWUSOAoVY2VydGlm'
    'aWNhdGV1cGxvYWRkYXRlGP6axvABIAEoCVIVY2VydGlmaWNhdGV1cGxvYWRkYXRlEjkKFmRpc3'
    'RyaWJ1dGlvbmRvbWFpbm5hbWUY3av5fiABKAlSFmRpc3RyaWJ1dGlvbmRvbWFpbm5hbWUSPQoY'
    'ZGlzdHJpYnV0aW9uaG9zdGVkem9uZWlkGMi30FggASgJUhhkaXN0cmlidXRpb25ob3N0ZWR6b2'
    '5laWQSIgoKZG9tYWlubmFtZRiL04+6ASABKAlSCmRvbWFpbm5hbWUSJwoNZG9tYWlubmFtZWFy'
    'bhiW3610IAEoCVINZG9tYWlubmFtZWFybhImCgxkb21haW5uYW1laWQYqPycjgEgASgJUgxkb2'
    '1haW5uYW1laWQSTAoQZG9tYWlubmFtZXN0YXR1cxi1y5WDASABKA4yHC5hcGlnYXRld2F5LkRv'
    'bWFpbk5hbWVTdGF0dXNSEGRvbWFpbm5hbWVzdGF0dXMSPAoXZG9tYWlubmFtZXN0YXR1c21lc3'
    'NhZ2UYtMPhnwEgASgJUhdkb21haW5uYW1lc3RhdHVzbWVzc2FnZRJSChJlbmRwb2ludGFjY2Vz'
    'c21vZGUY3sqLqgEgASgOMh4uYXBpZ2F0ZXdheS5FbmRwb2ludEFjY2Vzc01vZGVSEmVuZHBvaW'
    '50YWNjZXNzbW9kZRJbChVlbmRwb2ludGNvbmZpZ3VyYXRpb24Yt6e96AEgASgLMiEuYXBpZ2F0'
    'ZXdheS5FbmRwb2ludENvbmZpZ3VyYXRpb25SFWVuZHBvaW50Y29uZmlndXJhdGlvbhItChBtYW'
    '5hZ2VtZW50cG9saWN5GLm/+yogASgJUhBtYW5hZ2VtZW50cG9saWN5EmAKF211dHVhbHRsc2F1'
    'dGhlbnRpY2F0aW9uGJvXti8gASgLMiMuYXBpZ2F0ZXdheS5NdXR1YWxUbHNBdXRoZW50aWNhdG'
    'lvblIXbXV0dWFsdGxzYXV0aGVudGljYXRpb24SUwojb3duZXJzaGlwdmVyaWZpY2F0aW9uY2Vy'
    'dGlmaWNhdGVhcm4YlI/5LyABKAlSI293bmVyc2hpcHZlcmlmaWNhdGlvbmNlcnRpZmljYXRlYX'
    'JuEhkKBnBvbGljeRiA9YN2IAEoCVIGcG9saWN5EjkKFnJlZ2lvbmFsY2VydGlmaWNhdGVhcm4Y'
    '0/CtQSABKAlSFnJlZ2lvbmFsY2VydGlmaWNhdGVhcm4SPAoXcmVnaW9uYWxjZXJ0aWZpY2F0ZW'
    '5hbWUYzfTtzQEgASgJUhdyZWdpb25hbGNlcnRpZmljYXRlbmFtZRIxChJyZWdpb25hbGRvbWFp'
    'bm5hbWUYsM/EXiABKAlSEnJlZ2lvbmFsZG9tYWlubmFtZRI2ChRyZWdpb25hbGhvc3RlZHpvbm'
    'VpZBiVqejeASABKAlSFHJlZ2lvbmFsaG9zdGVkem9uZWlkEj0KC3JvdXRpbmdtb2RlGOfVuPEB'
    'IAEoDjIXLmFwaWdhdGV3YXkuUm91dGluZ01vZGVSC3JvdXRpbmdtb2RlEkYKDnNlY3VyaXR5cG'
    '9saWN5GN7UwOoBIAEoDjIaLmFwaWdhdGV3YXkuU2VjdXJpdHlQb2xpY3lSDnNlY3VyaXR5cG9s'
    'aWN5EjgKBHRhZ3MYodfboAEgAygLMiAuYXBpZ2F0ZXdheS5Eb21haW5OYW1lLlRhZ3NFbnRyeV'
    'IEdGFncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZh'
    'bHVlOgI4AQ==');

@$core.Deprecated('Use domainNameAccessAssociationDescriptor instead')
const DomainNameAccessAssociation$json = {
  '1': 'DomainNameAccessAssociation',
  '2': [
    {
      '1': 'accessassociationsource',
      '3': 328257828,
      '4': 1,
      '5': 9,
      '10': 'accessassociationsource'
    },
    {
      '1': 'accessassociationsourcetype',
      '3': 176397628,
      '4': 1,
      '5': 14,
      '6': '.apigateway.AccessAssociationSourceType',
      '10': 'accessassociationsourcetype'
    },
    {
      '1': 'domainnameaccessassociationarn',
      '3': 281017927,
      '4': 1,
      '5': 9,
      '10': 'domainnameaccessassociationarn'
    },
    {
      '1': 'domainnamearn',
      '3': 244019094,
      '4': 1,
      '5': 9,
      '10': 'domainnamearn'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.DomainNameAccessAssociation.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [DomainNameAccessAssociation_TagsEntry$json],
};

@$core.Deprecated('Use domainNameAccessAssociationDescriptor instead')
const DomainNameAccessAssociation_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DomainNameAccessAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainNameAccessAssociationDescriptor = $convert.base64Decode(
    'ChtEb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb24SPAoXYWNjZXNzYXNzb2NpYXRpb25zb3VyY2'
    'UYpKLDnAEgASgJUhdhY2Nlc3Nhc3NvY2lhdGlvbnNvdXJjZRJsChthY2Nlc3Nhc3NvY2lhdGlv'
    'bnNvdXJjZXR5cGUYvLqOVCABKA4yJy5hcGlnYXRld2F5LkFjY2Vzc0Fzc29jaWF0aW9uU291cm'
    'NlVHlwZVIbYWNjZXNzYXNzb2NpYXRpb25zb3VyY2V0eXBlEkoKHmRvbWFpbm5hbWVhY2Nlc3Nh'
    'c3NvY2lhdGlvbmFybhjH/P+FASABKAlSHmRvbWFpbm5hbWVhY2Nlc3Nhc3NvY2lhdGlvbmFybh'
    'InCg1kb21haW5uYW1lYXJuGJbfrXQgASgJUg1kb21haW5uYW1lYXJuEkkKBHRhZ3MYodfboAEg'
    'AygLMjEuYXBpZ2F0ZXdheS5Eb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb24uVGFnc0VudHJ5Ug'
    'R0YWdzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFs'
    'dWU6AjgB');

@$core.Deprecated('Use domainNameAccessAssociationsDescriptor instead')
const DomainNameAccessAssociations$json = {
  '1': 'DomainNameAccessAssociations',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.DomainNameAccessAssociation',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `DomainNameAccessAssociations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainNameAccessAssociationsDescriptor =
    $convert.base64Decode(
        'ChxEb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb25zEkEKBWl0ZW1zGJDn5NMBIAMoCzInLmFwaW'
        'dhdGV3YXkuRG9tYWluTmFtZUFjY2Vzc0Fzc29jaWF0aW9uUgVpdGVtcxIeCghwb3NpdGlvbhiL'
        'nL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use domainNamesDescriptor instead')
const DomainNames$json = {
  '1': 'DomainNames',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.DomainName',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `DomainNames`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainNamesDescriptor = $convert.base64Decode(
    'CgtEb21haW5OYW1lcxIwCgVpdGVtcxiQ5+TTASADKAsyFi5hcGlnYXRld2F5LkRvbWFpbk5hbW'
    'VSBWl0ZW1zEh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24=');

@$core.Deprecated('Use endpointConfigurationDescriptor instead')
const EndpointConfiguration$json = {
  '1': 'EndpointConfiguration',
  '2': [
    {
      '1': 'ipaddresstype',
      '3': 393863813,
      '4': 1,
      '5': 14,
      '6': '.apigateway.IpAddressType',
      '10': 'ipaddresstype'
    },
    {
      '1': 'types',
      '3': 534824091,
      '4': 3,
      '5': 14,
      '6': '.apigateway.EndpointType',
      '10': 'types'
    },
    {
      '1': 'vpcendpointids',
      '3': 135085374,
      '4': 3,
      '5': 9,
      '10': 'vpcendpointids'
    },
  ],
};

/// Descriptor for `EndpointConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endpointConfigurationDescriptor = $convert.base64Decode(
    'ChVFbmRwb2ludENvbmZpZ3VyYXRpb24SQwoNaXBhZGRyZXNzdHlwZRiFxee7ASABKA4yGS5hcG'
    'lnYXRld2F5LklwQWRkcmVzc1R5cGVSDWlwYWRkcmVzc3R5cGUSMgoFdHlwZXMYm4mD/wEgAygO'
    'MhguYXBpZ2F0ZXdheS5FbmRwb2ludFR5cGVSBXR5cGVzEikKDnZwY2VuZHBvaW50aWRzGL76tE'
    'AgAygJUg52cGNlbmRwb2ludGlkcw==');

@$core.Deprecated('Use exportResponseDescriptor instead')
const ExportResponse$json = {
  '1': 'ExportResponse',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'contentdisposition',
      '3': 375146466,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {'1': 'contenttype', '3': 281764659, '4': 1, '5': 9, '10': 'contenttype'},
  ],
};

/// Descriptor for `ExportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportResponseDescriptor = $convert.base64Decode(
    'Cg5FeHBvcnRSZXNwb25zZRIcCgRib2R5GPbyqd0BIAEoDEIEiLUYAVIEYm9keRIyChJjb250ZW'
    '50ZGlzcG9zaXRpb24Y4o/xsgEgASgJUhJjb250ZW50ZGlzcG9zaXRpb24SJAoLY29udGVudHR5'
    'cGUYs8athgEgASgJUgtjb250ZW50dHlwZQ==');

@$core.Deprecated('Use flushStageAuthorizersCacheRequestDescriptor instead')
const FlushStageAuthorizersCacheRequest$json = {
  '1': 'FlushStageAuthorizersCacheRequest',
  '2': [
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
};

/// Descriptor for `FlushStageAuthorizersCacheRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List flushStageAuthorizersCacheRequestDescriptor =
    $convert.base64Decode(
        'CiFGbHVzaFN0YWdlQXV0aG9yaXplcnNDYWNoZVJlcXVlc3QSIAoJcmVzdGFwaWlkGJmkgbcBIA'
        'EoCVIJcmVzdGFwaWlkEh8KCXN0YWdlbmFtZRiP3McEIAEoCVIJc3RhZ2VuYW1l');

@$core.Deprecated('Use flushStageCacheRequestDescriptor instead')
const FlushStageCacheRequest$json = {
  '1': 'FlushStageCacheRequest',
  '2': [
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
};

/// Descriptor for `FlushStageCacheRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List flushStageCacheRequestDescriptor =
    $convert.base64Decode(
        'ChZGbHVzaFN0YWdlQ2FjaGVSZXF1ZXN0EiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcG'
        'lpZBIfCglzdGFnZW5hbWUYj9zHBCABKAlSCXN0YWdlbmFtZQ==');

@$core.Deprecated('Use gatewayResponseDescriptor instead')
const GatewayResponse$json = {
  '1': 'GatewayResponse',
  '2': [
    {
      '1': 'defaultresponse',
      '3': 164690914,
      '4': 1,
      '5': 8,
      '10': 'defaultresponse'
    },
    {
      '1': 'responseparameters',
      '3': 64271839,
      '4': 3,
      '5': 11,
      '6': '.apigateway.GatewayResponse.ResponseparametersEntry',
      '10': 'responseparameters'
    },
    {
      '1': 'responsetemplates',
      '3': 107376570,
      '4': 3,
      '5': 11,
      '6': '.apigateway.GatewayResponse.ResponsetemplatesEntry',
      '10': 'responsetemplates'
    },
    {
      '1': 'responsetype',
      '3': 377935935,
      '4': 1,
      '5': 14,
      '6': '.apigateway.GatewayResponseType',
      '10': 'responsetype'
    },
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
  '3': [
    GatewayResponse_ResponseparametersEntry$json,
    GatewayResponse_ResponsetemplatesEntry$json
  ],
};

@$core.Deprecated('Use gatewayResponseDescriptor instead')
const GatewayResponse_ResponseparametersEntry$json = {
  '1': 'ResponseparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use gatewayResponseDescriptor instead')
const GatewayResponse_ResponsetemplatesEntry$json = {
  '1': 'ResponsetemplatesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GatewayResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List gatewayResponseDescriptor = $convert.base64Decode(
    'Cg9HYXRld2F5UmVzcG9uc2USKwoPZGVmYXVsdHJlc3BvbnNlGOL3w04gASgIUg9kZWZhdWx0cm'
    'VzcG9uc2USZgoScmVzcG9uc2VwYXJhbWV0ZXJzGN/r0h4gAygLMjMuYXBpZ2F0ZXdheS5HYXRl'
    'd2F5UmVzcG9uc2UuUmVzcG9uc2VwYXJhbWV0ZXJzRW50cnlSEnJlc3BvbnNlcGFyYW1ldGVycx'
    'JjChFyZXNwb25zZXRlbXBsYXRlcxi635kzIAMoCzIyLmFwaWdhdGV3YXkuR2F0ZXdheVJlc3Bv'
    'bnNlLlJlc3BvbnNldGVtcGxhdGVzRW50cnlSEXJlc3BvbnNldGVtcGxhdGVzEkcKDHJlc3Bvbn'
    'NldHlwZRi/sJu0ASABKA4yHy5hcGlnYXRld2F5LkdhdGV3YXlSZXNwb25zZVR5cGVSDHJlc3Bv'
    'bnNldHlwZRIiCgpzdGF0dXNjb2RlGJ+B344BIAEoCVIKc3RhdHVzY29kZRpFChdSZXNwb25zZX'
    'BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6'
    'AjgBGkQKFlJlc3BvbnNldGVtcGxhdGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdW'
    'UYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use gatewayResponsesDescriptor instead')
const GatewayResponses$json = {
  '1': 'GatewayResponses',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.GatewayResponse',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GatewayResponses`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List gatewayResponsesDescriptor = $convert.base64Decode(
    'ChBHYXRld2F5UmVzcG9uc2VzEjUKBWl0ZW1zGJDn5NMBIAMoCzIbLmFwaWdhdGV3YXkuR2F0ZX'
    'dheVJlc3BvbnNlUgVpdGVtcxIeCghwb3NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use generateClientCertificateRequestDescriptor instead')
const GenerateClientCertificateRequest$json = {
  '1': 'GenerateClientCertificateRequest',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.GenerateClientCertificateRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [GenerateClientCertificateRequest_TagsEntry$json],
};

@$core.Deprecated('Use generateClientCertificateRequestDescriptor instead')
const GenerateClientCertificateRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GenerateClientCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateClientCertificateRequestDescriptor =
    $convert.base64Decode(
        'CiBHZW5lcmF0ZUNsaWVudENlcnRpZmljYXRlUmVxdWVzdBIkCgtkZXNjcmlwdGlvbhjq9ryjAS'
        'ABKAlSC2Rlc2NyaXB0aW9uEk4KBHRhZ3MYodfboAEgAygLMjYuYXBpZ2F0ZXdheS5HZW5lcmF0'
        'ZUNsaWVudENlcnRpZmljYXRlUmVxdWVzdC5UYWdzRW50cnlSBHRhZ3MaNwoJVGFnc0VudHJ5Eh'
        'AKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getAccountRequestDescriptor instead')
const GetAccountRequest$json = {
  '1': 'GetAccountRequest',
};

/// Descriptor for `GetAccountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountRequestDescriptor =
    $convert.base64Decode('ChFHZXRBY2NvdW50UmVxdWVzdA==');

@$core.Deprecated('Use getApiKeyRequestDescriptor instead')
const GetApiKeyRequest$json = {
  '1': 'GetApiKeyRequest',
  '2': [
    {'1': 'apikey', '3': 490490655, '4': 1, '5': 9, '10': 'apikey'},
    {'1': 'includevalue', '3': 462860701, '4': 1, '5': 8, '10': 'includevalue'},
  ],
};

/// Descriptor for `GetApiKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getApiKeyRequestDescriptor = $convert.base64Decode(
    'ChBHZXRBcGlLZXlSZXF1ZXN0EhoKBmFwaWtleRiflvHpASABKAlSBmFwaWtleRImCgxpbmNsdW'
    'RldmFsdWUYnePa3AEgASgIUgxpbmNsdWRldmFsdWU=');

@$core.Deprecated('Use getApiKeysRequestDescriptor instead')
const GetApiKeysRequest$json = {
  '1': 'GetApiKeysRequest',
  '2': [
    {'1': 'customerid', '3': 227830269, '4': 1, '5': 9, '10': 'customerid'},
    {
      '1': 'includevalues',
      '3': 490347326,
      '4': 1,
      '5': 8,
      '10': 'includevalues'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'namequery', '3': 51018795, '4': 1, '5': 9, '10': 'namequery'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GetApiKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getApiKeysRequestDescriptor = $convert.base64Decode(
    'ChFHZXRBcGlLZXlzUmVxdWVzdBIhCgpjdXN0b21lcmlkGP3T0WwgASgJUgpjdXN0b21lcmlkEi'
    'gKDWluY2x1ZGV2YWx1ZXMYvrbo6QEgASgIUg1pbmNsdWRldmFsdWVzEhgKBWxpbWl0GLWy65YB'
    'IAEoBVIFbGltaXQSHwoJbmFtZXF1ZXJ5GKv4qRggASgJUgluYW1lcXVlcnkSHgoIcG9zaXRpb2'
    '4Yi5y9mgEgASgJUghwb3NpdGlvbg==');

@$core.Deprecated('Use getAuthorizerRequestDescriptor instead')
const GetAuthorizerRequest$json = {
  '1': 'GetAuthorizerRequest',
  '2': [
    {'1': 'authorizerid', '3': 111773148, '4': 1, '5': 9, '10': 'authorizerid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetAuthorizerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAuthorizerRequestDescriptor = $convert.base64Decode(
    'ChRHZXRBdXRob3JpemVyUmVxdWVzdBIlCgxhdXRob3JpemVyaWQY3IumNSABKAlSDGF1dGhvcm'
    'l6ZXJpZBIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use getAuthorizersRequestDescriptor instead')
const GetAuthorizersRequest$json = {
  '1': 'GetAuthorizersRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetAuthorizersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAuthorizersRequestDescriptor = $convert.base64Decode(
    'ChVHZXRBdXRob3JpemVyc1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIeCghwb3'
    'NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9uEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3Rh'
    'cGlpZA==');

@$core.Deprecated('Use getBasePathMappingRequestDescriptor instead')
const GetBasePathMappingRequest$json = {
  '1': 'GetBasePathMappingRequest',
  '2': [
    {'1': 'basepath', '3': 267528880, '4': 1, '5': 9, '10': 'basepath'},
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
  ],
};

/// Descriptor for `GetBasePathMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBasePathMappingRequestDescriptor = $convert.base64Decode(
    'ChlHZXRCYXNlUGF0aE1hcHBpbmdSZXF1ZXN0Eh0KCGJhc2VwYXRoGLDVyH8gASgJUghiYXNlcG'
    'F0aBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG9tYWlubmFtZRImCgxkb21haW5uYW1laWQY'
    'qPycjgEgASgJUgxkb21haW5uYW1laWQ=');

@$core.Deprecated('Use getBasePathMappingsRequestDescriptor instead')
const GetBasePathMappingsRequest$json = {
  '1': 'GetBasePathMappingsRequest',
  '2': [
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GetBasePathMappingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBasePathMappingsRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRCYXNlUGF0aE1hcHBpbmdzUmVxdWVzdBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG'
        '9tYWlubmFtZRImCgxkb21haW5uYW1laWQYqPycjgEgASgJUgxkb21haW5uYW1laWQSGAoFbGlt'
        'aXQYtbLrlgEgASgFUgVsaW1pdBIeCghwb3NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use getClientCertificateRequestDescriptor instead')
const GetClientCertificateRequest$json = {
  '1': 'GetClientCertificateRequest',
  '2': [
    {
      '1': 'clientcertificateid',
      '3': 276222909,
      '4': 1,
      '5': 9,
      '10': 'clientcertificateid'
    },
  ],
};

/// Descriptor for `GetClientCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getClientCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRDbGllbnRDZXJ0aWZpY2F0ZVJlcXVlc3QSNAoTY2xpZW50Y2VydGlmaWNhdGVpZBi9p9'
        'uDASABKAlSE2NsaWVudGNlcnRpZmljYXRlaWQ=');

@$core.Deprecated('Use getClientCertificatesRequestDescriptor instead')
const GetClientCertificatesRequest$json = {
  '1': 'GetClientCertificatesRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GetClientCertificatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getClientCertificatesRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRDbGllbnRDZXJ0aWZpY2F0ZXNSZXF1ZXN0EhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaX'
        'QSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbg==');

@$core.Deprecated('Use getDeploymentRequestDescriptor instead')
const GetDeploymentRequest$json = {
  '1': 'GetDeploymentRequest',
  '2': [
    {'1': 'deploymentid', '3': 439369188, '4': 1, '5': 9, '10': 'deploymentid'},
    {'1': 'embed', '3': 136029775, '4': 3, '5': 9, '10': 'embed'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeploymentRequestDescriptor = $convert.base64Decode(
    'ChRHZXREZXBsb3ltZW50UmVxdWVzdBImCgxkZXBsb3ltZW50aWQY5PvA0QEgASgJUgxkZXBsb3'
    'ltZW50aWQSFwoFZW1iZWQYz8zuQCADKAlSBWVtYmVkEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlS'
    'CXJlc3RhcGlpZA==');

@$core.Deprecated('Use getDeploymentsRequestDescriptor instead')
const GetDeploymentsRequest$json = {
  '1': 'GetDeploymentsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetDeploymentsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeploymentsRequestDescriptor = $convert.base64Decode(
    'ChVHZXREZXBsb3ltZW50c1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIeCghwb3'
    'NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9uEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3Rh'
    'cGlpZA==');

@$core.Deprecated('Use getDocumentationPartRequestDescriptor instead')
const GetDocumentationPartRequest$json = {
  '1': 'GetDocumentationPartRequest',
  '2': [
    {
      '1': 'documentationpartid',
      '3': 286552774,
      '4': 1,
      '5': 9,
      '10': 'documentationpartid'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetDocumentationPartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDocumentationPartRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXREb2N1bWVudGF0aW9uUGFydFJlcXVlc3QSNAoTZG9jdW1lbnRhdGlvbnBhcnRpZBjG5d'
        'GIASABKAlSE2RvY3VtZW50YXRpb25wYXJ0aWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVz'
        'dGFwaWlk');

@$core.Deprecated('Use getDocumentationPartsRequestDescriptor instead')
const GetDocumentationPartsRequest$json = {
  '1': 'GetDocumentationPartsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'locationstatus',
      '3': 532215305,
      '4': 1,
      '5': 14,
      '6': '.apigateway.LocationStatusType',
      '10': 'locationstatus'
    },
    {'1': 'namequery', '3': 51018795, '4': 1, '5': 9, '10': 'namequery'},
    {'1': 'path', '3': 75975991, '4': 1, '5': 9, '10': 'path'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.apigateway.DocumentationPartType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `GetDocumentationPartsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDocumentationPartsRequestDescriptor = $convert.base64Decode(
    'ChxHZXREb2N1bWVudGF0aW9uUGFydHNSZXF1ZXN0EhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaX'
    'QSSgoObG9jYXRpb25zdGF0dXMYiezj/QEgASgOMh4uYXBpZ2F0ZXdheS5Mb2NhdGlvblN0YXR1'
    'c1R5cGVSDmxvY2F0aW9uc3RhdHVzEh8KCW5hbWVxdWVyeRir+KkYIAEoCVIJbmFtZXF1ZXJ5Eh'
    'UKBHBhdGgYt5qdJCABKAlSBHBhdGgSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbhIg'
    'CglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQSOQoEdHlwZRjO4p+JASABKA4yIS5hcG'
    'lnYXRld2F5LkRvY3VtZW50YXRpb25QYXJ0VHlwZVIEdHlwZQ==');

@$core.Deprecated('Use getDocumentationVersionRequestDescriptor instead')
const GetDocumentationVersionRequest$json = {
  '1': 'GetDocumentationVersionRequest',
  '2': [
    {
      '1': 'documentationversion',
      '3': 167009804,
      '4': 1,
      '5': 9,
      '10': 'documentationversion'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetDocumentationVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDocumentationVersionRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXREb2N1bWVudGF0aW9uVmVyc2lvblJlcXVlc3QSNQoUZG9jdW1lbnRhdGlvbnZlcnNpb2'
        '4YjLzRTyABKAlSFGRvY3VtZW50YXRpb252ZXJzaW9uEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlS'
        'CXJlc3RhcGlpZA==');

@$core.Deprecated('Use getDocumentationVersionsRequestDescriptor instead')
const GetDocumentationVersionsRequest$json = {
  '1': 'GetDocumentationVersionsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetDocumentationVersionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDocumentationVersionsRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXREb2N1bWVudGF0aW9uVmVyc2lvbnNSZXF1ZXN0EhgKBWxpbWl0GLWy65YBIAEoBVIFbG'
        'ltaXQSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbhIgCglyZXN0YXBpaWQYmaSBtwEg'
        'ASgJUglyZXN0YXBpaWQ=');

@$core
    .Deprecated('Use getDomainNameAccessAssociationsRequestDescriptor instead')
const GetDomainNameAccessAssociationsRequest$json = {
  '1': 'GetDomainNameAccessAssociationsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {
      '1': 'resourceowner',
      '3': 259175301,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ResourceOwner',
      '10': 'resourceowner'
    },
  ],
};

/// Descriptor for `GetDomainNameAccessAssociationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDomainNameAccessAssociationsRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXREb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb25zUmVxdWVzdBIYCgVsaW1pdBi1suuWAS'
        'ABKAVSBWxpbWl0Eh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24SQgoNcmVzb3VyY2Vv'
        'd25lchiF58p7IAEoDjIZLmFwaWdhdGV3YXkuUmVzb3VyY2VPd25lclINcmVzb3VyY2Vvd25lcg'
        '==');

@$core.Deprecated('Use getDomainNameRequestDescriptor instead')
const GetDomainNameRequest$json = {
  '1': 'GetDomainNameRequest',
  '2': [
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
  ],
};

/// Descriptor for `GetDomainNameRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDomainNameRequestDescriptor = $convert.base64Decode(
    'ChRHZXREb21haW5OYW1lUmVxdWVzdBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG9tYWlubm'
    'FtZRImCgxkb21haW5uYW1laWQYqPycjgEgASgJUgxkb21haW5uYW1laWQ=');

@$core.Deprecated('Use getDomainNamesRequestDescriptor instead')
const GetDomainNamesRequest$json = {
  '1': 'GetDomainNamesRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {
      '1': 'resourceowner',
      '3': 259175301,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ResourceOwner',
      '10': 'resourceowner'
    },
  ],
};

/// Descriptor for `GetDomainNamesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDomainNamesRequestDescriptor = $convert.base64Decode(
    'ChVHZXREb21haW5OYW1lc1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIeCghwb3'
    'NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9uEkIKDXJlc291cmNlb3duZXIYhefKeyABKA4yGS5h'
    'cGlnYXRld2F5LlJlc291cmNlT3duZXJSDXJlc291cmNlb3duZXI=');

@$core.Deprecated('Use getExportRequestDescriptor instead')
const GetExportRequest$json = {
  '1': 'GetExportRequest',
  '2': [
    {'1': 'accepts', '3': 192079791, '4': 1, '5': 9, '10': 'accepts'},
    {'1': 'exporttype', '3': 243495788, '4': 1, '5': 9, '10': 'exporttype'},
    {
      '1': 'parameters',
      '3': 145043162,
      '4': 3,
      '5': 11,
      '6': '.apigateway.GetExportRequest.ParametersEntry',
      '10': 'parameters'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
  '3': [GetExportRequest_ParametersEntry$json],
};

@$core.Deprecated('Use getExportRequestDescriptor instead')
const GetExportRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetExportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getExportRequestDescriptor = $convert.base64Decode(
    'ChBHZXRFeHBvcnRSZXF1ZXN0EhsKB2FjY2VwdHMYr8/LWyABKAlSB2FjY2VwdHMSIQoKZXhwb3'
    'J0dHlwZRjs5o10IAEoCVIKZXhwb3J0dHlwZRJPCgpwYXJhbWV0ZXJzGNrdlEUgAygLMiwuYXBp'
    'Z2F0ZXdheS5HZXRFeHBvcnRSZXF1ZXN0LlBhcmFtZXRlcnNFbnRyeVIKcGFyYW1ldGVycxIgCg'
    'lyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQSHwoJc3RhZ2VuYW1lGI/cxwQgASgJUglz'
    'dGFnZW5hbWUaPQoPUGFyYW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGA'
    'IgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getGatewayResponseRequestDescriptor instead')
const GetGatewayResponseRequest$json = {
  '1': 'GetGatewayResponseRequest',
  '2': [
    {
      '1': 'responsetype',
      '3': 377935935,
      '4': 1,
      '5': 14,
      '6': '.apigateway.GatewayResponseType',
      '10': 'responsetype'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetGatewayResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGatewayResponseRequestDescriptor = $convert.base64Decode(
    'ChlHZXRHYXRld2F5UmVzcG9uc2VSZXF1ZXN0EkcKDHJlc3BvbnNldHlwZRi/sJu0ASABKA4yHy'
    '5hcGlnYXRld2F5LkdhdGV3YXlSZXNwb25zZVR5cGVSDHJlc3BvbnNldHlwZRIgCglyZXN0YXBp'
    'aWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use getGatewayResponsesRequestDescriptor instead')
const GetGatewayResponsesRequest$json = {
  '1': 'GetGatewayResponsesRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetGatewayResponsesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGatewayResponsesRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRHYXRld2F5UmVzcG9uc2VzUmVxdWVzdBIYCgVsaW1pdBi1suuWASABKAVSBWxpbWl0Eh'
        '4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24SIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJ'
        'cmVzdGFwaWlk');

@$core.Deprecated('Use getIntegrationRequestDescriptor instead')
const GetIntegrationRequest$json = {
  '1': 'GetIntegrationRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIntegrationRequestDescriptor = $convert.base64Decode(
    'ChVHZXRJbnRlZ3JhdGlvblJlcXVlc3QSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaHR0cG1ldG'
    'hvZBIiCgpyZXNvdXJjZWlkGLG9iZgBIAEoCVIKcmVzb3VyY2VpZBIgCglyZXN0YXBpaWQYmaSB'
    'twEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use getIntegrationResponseRequestDescriptor instead')
const GetIntegrationResponseRequest$json = {
  '1': 'GetIntegrationResponseRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
};

/// Descriptor for `GetIntegrationResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIntegrationResponseRequestDescriptor = $convert.base64Decode(
    'Ch1HZXRJbnRlZ3JhdGlvblJlc3BvbnNlUmVxdWVzdBIhCgpodHRwbWV0aG9kGPHz+zYgASgJUg'
    'podHRwbWV0aG9kEiIKCnJlc291cmNlaWQYsb2JmAEgASgJUgpyZXNvdXJjZWlkEiAKCXJlc3Rh'
    'cGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZBIiCgpzdGF0dXNjb2RlGJ+B344BIAEoCVIKc3RhdH'
    'VzY29kZQ==');

@$core.Deprecated('Use getMethodRequestDescriptor instead')
const GetMethodRequest$json = {
  '1': 'GetMethodRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetMethodRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMethodRequestDescriptor = $convert.base64Decode(
    'ChBHZXRNZXRob2RSZXF1ZXN0EiEKCmh0dHBtZXRob2QY8fP7NiABKAlSCmh0dHBtZXRob2QSIg'
    'oKcmVzb3VyY2VpZBixvYmYASABKAlSCnJlc291cmNlaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEo'
    'CVIJcmVzdGFwaWlk');

@$core.Deprecated('Use getMethodResponseRequestDescriptor instead')
const GetMethodResponseRequest$json = {
  '1': 'GetMethodResponseRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
};

/// Descriptor for `GetMethodResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMethodResponseRequestDescriptor = $convert.base64Decode(
    'ChhHZXRNZXRob2RSZXNwb25zZVJlcXVlc3QSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaHR0cG'
    '1ldGhvZBIiCgpyZXNvdXJjZWlkGLG9iZgBIAEoCVIKcmVzb3VyY2VpZBIgCglyZXN0YXBpaWQY'
    'maSBtwEgASgJUglyZXN0YXBpaWQSIgoKc3RhdHVzY29kZRifgd+OASABKAlSCnN0YXR1c2NvZG'
    'U=');

@$core.Deprecated('Use getModelRequestDescriptor instead')
const GetModelRequest$json = {
  '1': 'GetModelRequest',
  '2': [
    {'1': 'flatten', '3': 100024266, '4': 1, '5': 8, '10': 'flatten'},
    {'1': 'modelname', '3': 176835330, '4': 1, '5': 9, '10': 'modelname'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetModelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getModelRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRNb2RlbFJlcXVlc3QSGwoHZmxhdHRlbhjK/9gvIAEoCFIHZmxhdHRlbhIfCgltb2RlbG'
    '5hbWUYgpapVCABKAlSCW1vZGVsbmFtZRIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBp'
    'aWQ=');

@$core.Deprecated('Use getModelTemplateRequestDescriptor instead')
const GetModelTemplateRequest$json = {
  '1': 'GetModelTemplateRequest',
  '2': [
    {'1': 'modelname', '3': 176835330, '4': 1, '5': 9, '10': 'modelname'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetModelTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getModelTemplateRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRNb2RlbFRlbXBsYXRlUmVxdWVzdBIfCgltb2RlbG5hbWUYgpapVCABKAlSCW1vZGVsbm'
        'FtZRIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use getModelsRequestDescriptor instead')
const GetModelsRequest$json = {
  '1': 'GetModelsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetModelsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getModelsRequestDescriptor = $convert.base64Decode(
    'ChBHZXRNb2RlbHNSZXF1ZXN0EhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSHgoIcG9zaXRpb2'
    '4Yi5y9mgEgASgJUghwb3NpdGlvbhIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use getRequestValidatorRequestDescriptor instead')
const GetRequestValidatorRequest$json = {
  '1': 'GetRequestValidatorRequest',
  '2': [
    {
      '1': 'requestvalidatorid',
      '3': 517546134,
      '4': 1,
      '5': 9,
      '10': 'requestvalidatorid'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetRequestValidatorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRequestValidatorRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRSZXF1ZXN0VmFsaWRhdG9yUmVxdWVzdBIyChJyZXF1ZXN0dmFsaWRhdG9yaWQYlsHk9g'
        'EgASgJUhJyZXF1ZXN0dmFsaWRhdG9yaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFw'
        'aWlk');

@$core.Deprecated('Use getRequestValidatorsRequestDescriptor instead')
const GetRequestValidatorsRequest$json = {
  '1': 'GetRequestValidatorsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetRequestValidatorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRequestValidatorsRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRSZXF1ZXN0VmFsaWRhdG9yc1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdB'
        'IeCghwb3NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9uEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlS'
        'CXJlc3RhcGlpZA==');

@$core.Deprecated('Use getResourceRequestDescriptor instead')
const GetResourceRequest$json = {
  '1': 'GetResourceRequest',
  '2': [
    {'1': 'embed', '3': 136029775, '4': 3, '5': 9, '10': 'embed'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourceRequestDescriptor = $convert.base64Decode(
    'ChJHZXRSZXNvdXJjZVJlcXVlc3QSFwoFZW1iZWQYz8zuQCADKAlSBWVtYmVkEiIKCnJlc291cm'
    'NlaWQYsb2JmAEgASgJUgpyZXNvdXJjZWlkEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3Rh'
    'cGlpZA==');

@$core.Deprecated('Use getResourcesRequestDescriptor instead')
const GetResourcesRequest$json = {
  '1': 'GetResourcesRequest',
  '2': [
    {'1': 'embed', '3': 136029775, '4': 3, '5': 9, '10': 'embed'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetResourcesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcesRequestDescriptor = $convert.base64Decode(
    'ChNHZXRSZXNvdXJjZXNSZXF1ZXN0EhcKBWVtYmVkGM/M7kAgAygJUgVlbWJlZBIYCgVsaW1pdB'
    'i1suuWASABKAVSBWxpbWl0Eh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24SIAoJcmVz'
    'dGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use getRestApiRequestDescriptor instead')
const GetRestApiRequest$json = {
  '1': 'GetRestApiRequest',
  '2': [
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetRestApiRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRestApiRequestDescriptor = $convert.base64Decode(
    'ChFHZXRSZXN0QXBpUmVxdWVzdBIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use getRestApisRequestDescriptor instead')
const GetRestApisRequest$json = {
  '1': 'GetRestApisRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GetRestApisRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRestApisRequestDescriptor = $convert.base64Decode(
    'ChJHZXRSZXN0QXBpc1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIeCghwb3NpdG'
    'lvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use getSdkRequestDescriptor instead')
const GetSdkRequest$json = {
  '1': 'GetSdkRequest',
  '2': [
    {
      '1': 'parameters',
      '3': 145043162,
      '4': 3,
      '5': 11,
      '6': '.apigateway.GetSdkRequest.ParametersEntry',
      '10': 'parameters'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'sdktype', '3': 146213956, '4': 1, '5': 9, '10': 'sdktype'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
  '3': [GetSdkRequest_ParametersEntry$json],
};

@$core.Deprecated('Use getSdkRequestDescriptor instead')
const GetSdkRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetSdkRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSdkRequestDescriptor = $convert.base64Decode(
    'Cg1HZXRTZGtSZXF1ZXN0EkwKCnBhcmFtZXRlcnMY2t2URSADKAsyKS5hcGlnYXRld2F5LkdldF'
    'Nka1JlcXVlc3QuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJzEiAKCXJlc3RhcGlpZBiZpIG3'
    'ASABKAlSCXJlc3RhcGlpZBIbCgdzZGt0eXBlGMSY3EUgASgJUgdzZGt0eXBlEh8KCXN0YWdlbm'
    'FtZRiP3McEIAEoCVIJc3RhZ2VuYW1lGj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlS'
    'A2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use getSdkTypeRequestDescriptor instead')
const GetSdkTypeRequest$json = {
  '1': 'GetSdkTypeRequest',
  '2': [
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetSdkTypeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSdkTypeRequestDescriptor = $convert
    .base64Decode('ChFHZXRTZGtUeXBlUmVxdWVzdBISCgJpZBjh1eG5ASABKAlSAmlk');

@$core.Deprecated('Use getSdkTypesRequestDescriptor instead')
const GetSdkTypesRequest$json = {
  '1': 'GetSdkTypesRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GetSdkTypesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSdkTypesRequestDescriptor = $convert.base64Decode(
    'ChJHZXRTZGtUeXBlc1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIeCghwb3NpdG'
    'lvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use getStageRequestDescriptor instead')
const GetStageRequest$json = {
  '1': 'GetStageRequest',
  '2': [
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
};

/// Descriptor for `GetStageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getStageRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRTdGFnZVJlcXVlc3QSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkEh8KCX'
    'N0YWdlbmFtZRiP3McEIAEoCVIJc3RhZ2VuYW1l');

@$core.Deprecated('Use getStagesRequestDescriptor instead')
const GetStagesRequest$json = {
  '1': 'GetStagesRequest',
  '2': [
    {'1': 'deploymentid', '3': 439369188, '4': 1, '5': 9, '10': 'deploymentid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `GetStagesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getStagesRequestDescriptor = $convert.base64Decode(
    'ChBHZXRTdGFnZXNSZXF1ZXN0EiYKDGRlcGxveW1lbnRpZBjk+8DRASABKAlSDGRlcGxveW1lbn'
    'RpZBIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use getTagsRequestDescriptor instead')
const GetTagsRequest$json = {
  '1': 'GetTagsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTagsRequestDescriptor = $convert.base64Decode(
    'Cg5HZXRUYWdzUmVxdWVzdBIYCgVsaW1pdBi1suuWASABKAVSBWxpbWl0Eh4KCHBvc2l0aW9uGI'
    'ucvZoBIAEoCVIIcG9zaXRpb24SIwoLcmVzb3VyY2Vhcm4YzcyqICABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use getUsagePlanKeyRequestDescriptor instead')
const GetUsagePlanKeyRequest$json = {
  '1': 'GetUsagePlanKeyRequest',
  '2': [
    {'1': 'keyid', '3': 479913282, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `GetUsagePlanKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUsagePlanKeyRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRVc2FnZVBsYW5LZXlSZXF1ZXN0EhgKBWtleWlkGMLK6+QBIAEoCVIFa2V5aWQSJAoLdX'
        'NhZ2VwbGFuaWQY1/Dl8gEgASgJUgt1c2FnZXBsYW5pZA==');

@$core.Deprecated('Use getUsagePlanKeysRequestDescriptor instead')
const GetUsagePlanKeysRequest$json = {
  '1': 'GetUsagePlanKeysRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'namequery', '3': 51018795, '4': 1, '5': 9, '10': 'namequery'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `GetUsagePlanKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUsagePlanKeysRequestDescriptor = $convert.base64Decode(
    'ChdHZXRVc2FnZVBsYW5LZXlzUmVxdWVzdBIYCgVsaW1pdBi1suuWASABKAVSBWxpbWl0Eh8KCW'
    '5hbWVxdWVyeRir+KkYIAEoCVIJbmFtZXF1ZXJ5Eh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9z'
    'aXRpb24SJAoLdXNhZ2VwbGFuaWQY1/Dl8gEgASgJUgt1c2FnZXBsYW5pZA==');

@$core.Deprecated('Use getUsagePlanRequestDescriptor instead')
const GetUsagePlanRequest$json = {
  '1': 'GetUsagePlanRequest',
  '2': [
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `GetUsagePlanRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUsagePlanRequestDescriptor = $convert.base64Decode(
    'ChNHZXRVc2FnZVBsYW5SZXF1ZXN0EiQKC3VzYWdlcGxhbmlkGNfw5fIBIAEoCVILdXNhZ2VwbG'
    'FuaWQ=');

@$core.Deprecated('Use getUsagePlansRequestDescriptor instead')
const GetUsagePlansRequest$json = {
  '1': 'GetUsagePlansRequest',
  '2': [
    {'1': 'keyid', '3': 479913282, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GetUsagePlansRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUsagePlansRequestDescriptor = $convert.base64Decode(
    'ChRHZXRVc2FnZVBsYW5zUmVxdWVzdBIYCgVrZXlpZBjCyuvkASABKAlSBWtleWlkEhgKBWxpbW'
    'l0GLWy65YBIAEoBVIFbGltaXQSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbg==');

@$core.Deprecated('Use getUsageRequestDescriptor instead')
const GetUsageRequest$json = {
  '1': 'GetUsageRequest',
  '2': [
    {'1': 'enddate', '3': 384831215, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'keyid', '3': 479913282, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `GetUsageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUsageRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRVc2FnZVJlcXVlc3QSHAoHZW5kZGF0ZRjvncC3ASABKAlSB2VuZGRhdGUSGAoFa2V5aW'
    'QYwsrr5AEgASgJUgVrZXlpZBIYCgVsaW1pdBi1suuWASABKAVSBWxpbWl0Eh4KCHBvc2l0aW9u'
    'GIucvZoBIAEoCVIIcG9zaXRpb24SIAoJc3RhcnRkYXRlGJyO/K0BIAEoCVIJc3RhcnRkYXRlEi'
    'QKC3VzYWdlcGxhbmlkGNfw5fIBIAEoCVILdXNhZ2VwbGFuaWQ=');

@$core.Deprecated('Use getVpcLinkRequestDescriptor instead')
const GetVpcLinkRequest$json = {
  '1': 'GetVpcLinkRequest',
  '2': [
    {'1': 'vpclinkid', '3': 27515846, '4': 1, '5': 9, '10': 'vpclinkid'},
  ],
};

/// Descriptor for `GetVpcLinkRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getVpcLinkRequestDescriptor = $convert.base64Decode(
    'ChFHZXRWcGNMaW5rUmVxdWVzdBIfCgl2cGNsaW5raWQYxrePDSABKAlSCXZwY2xpbmtpZA==');

@$core.Deprecated('Use getVpcLinksRequestDescriptor instead')
const GetVpcLinksRequest$json = {
  '1': 'GetVpcLinksRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `GetVpcLinksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getVpcLinksRequestDescriptor = $convert.base64Decode(
    'ChJHZXRWcGNMaW5rc1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIeCghwb3NpdG'
    'lvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use importApiKeysRequestDescriptor instead')
const ImportApiKeysRequest$json = {
  '1': 'ImportApiKeysRequest',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'failonwarnings',
      '3': 434363958,
      '4': 1,
      '5': 8,
      '10': 'failonwarnings'
    },
    {
      '1': 'format',
      '3': 429753683,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ApiKeysFormat',
      '10': 'format'
    },
  ],
};

/// Descriptor for `ImportApiKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importApiKeysRequestDescriptor = $convert.base64Decode(
    'ChRJbXBvcnRBcGlLZXlzUmVxdWVzdBIcCgRib2R5GPbyqd0BIAEoDEIEiLUYAVIEYm9keRIqCg'
    '5mYWlsb253YXJuaW5ncxi2vI/PASABKAhSDmZhaWxvbndhcm5pbmdzEjUKBmZvcm1hdBjTivbM'
    'ASABKA4yGS5hcGlnYXRld2F5LkFwaUtleXNGb3JtYXRSBmZvcm1hdA==');

@$core.Deprecated('Use importDocumentationPartsRequestDescriptor instead')
const ImportDocumentationPartsRequest$json = {
  '1': 'ImportDocumentationPartsRequest',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'failonwarnings',
      '3': 434363958,
      '4': 1,
      '5': 8,
      '10': 'failonwarnings'
    },
    {
      '1': 'mode',
      '3': 208592915,
      '4': 1,
      '5': 14,
      '6': '.apigateway.PutMode',
      '10': 'mode'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `ImportDocumentationPartsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importDocumentationPartsRequestDescriptor =
    $convert.base64Decode(
        'Ch9JbXBvcnREb2N1bWVudGF0aW9uUGFydHNSZXF1ZXN0EhwKBGJvZHkY9vKp3QEgASgMQgSItR'
        'gBUgRib2R5EioKDmZhaWxvbndhcm5pbmdzGLa8j88BIAEoCFIOZmFpbG9ud2FybmluZ3MSKgoE'
        'bW9kZRiTwLtjIAEoDjITLmFwaWdhdGV3YXkuUHV0TW9kZVIEbW9kZRIgCglyZXN0YXBpaWQYma'
        'SBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use importRestApiRequestDescriptor instead')
const ImportRestApiRequest$json = {
  '1': 'ImportRestApiRequest',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'failonwarnings',
      '3': 434363958,
      '4': 1,
      '5': 8,
      '10': 'failonwarnings'
    },
    {
      '1': 'parameters',
      '3': 145043162,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ImportRestApiRequest.ParametersEntry',
      '10': 'parameters'
    },
  ],
  '3': [ImportRestApiRequest_ParametersEntry$json],
};

@$core.Deprecated('Use importRestApiRequestDescriptor instead')
const ImportRestApiRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ImportRestApiRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importRestApiRequestDescriptor = $convert.base64Decode(
    'ChRJbXBvcnRSZXN0QXBpUmVxdWVzdBIcCgRib2R5GPbyqd0BIAEoDEIEiLUYAVIEYm9keRIqCg'
    '5mYWlsb253YXJuaW5ncxi2vI/PASABKAhSDmZhaWxvbndhcm5pbmdzElMKCnBhcmFtZXRlcnMY'
    '2t2URSADKAsyMC5hcGlnYXRld2F5LkltcG9ydFJlc3RBcGlSZXF1ZXN0LlBhcmFtZXRlcnNFbn'
    'RyeVIKcGFyYW1ldGVycxo9Cg9QYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoF'
    'dmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use integrationDescriptor instead')
const Integration$json = {
  '1': 'Integration',
  '2': [
    {
      '1': 'cachekeyparameters',
      '3': 481441313,
      '4': 3,
      '5': 9,
      '10': 'cachekeyparameters'
    },
    {
      '1': 'cachenamespace',
      '3': 85102753,
      '4': 1,
      '5': 9,
      '10': 'cachenamespace'
    },
    {'1': 'connectionid', '3': 450027965, '4': 1, '5': 9, '10': 'connectionid'},
    {
      '1': 'connectiontype',
      '3': 336253170,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ConnectionType',
      '10': 'connectiontype'
    },
    {
      '1': 'contenthandling',
      '3': 533182832,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ContentHandlingStrategy',
      '10': 'contenthandling'
    },
    {'1': 'credentials', '3': 150838226, '4': 1, '5': 9, '10': 'credentials'},
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'integrationresponses',
      '3': 386580464,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Integration.IntegrationresponsesEntry',
      '10': 'integrationresponses'
    },
    {
      '1': 'integrationtarget',
      '3': 17646705,
      '4': 1,
      '5': 9,
      '10': 'integrationtarget'
    },
    {
      '1': 'passthroughbehavior',
      '3': 310796908,
      '4': 1,
      '5': 9,
      '10': 'passthroughbehavior'
    },
    {
      '1': 'requestparameters',
      '3': 523499939,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Integration.RequestparametersEntry',
      '10': 'requestparameters'
    },
    {
      '1': 'requesttemplates',
      '3': 333512166,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Integration.RequesttemplatesEntry',
      '10': 'requesttemplates'
    },
    {
      '1': 'responsetransfermode',
      '3': 458910787,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ResponseTransferMode',
      '10': 'responsetransfermode'
    },
    {
      '1': 'timeoutinmillis',
      '3': 378229126,
      '4': 1,
      '5': 5,
      '10': 'timeoutinmillis'
    },
    {
      '1': 'tlsconfig',
      '3': 108946693,
      '4': 1,
      '5': 11,
      '6': '.apigateway.TlsConfig',
      '10': 'tlsconfig'
    },
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.apigateway.IntegrationType',
      '10': 'type'
    },
    {'1': 'uri', '3': 395269118, '4': 1, '5': 9, '10': 'uri'},
  ],
  '3': [
    Integration_IntegrationresponsesEntry$json,
    Integration_RequestparametersEntry$json,
    Integration_RequesttemplatesEntry$json
  ],
};

@$core.Deprecated('Use integrationDescriptor instead')
const Integration_IntegrationresponsesEntry$json = {
  '1': 'IntegrationresponsesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.apigateway.IntegrationResponse',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use integrationDescriptor instead')
const Integration_RequestparametersEntry$json = {
  '1': 'RequestparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use integrationDescriptor instead')
const Integration_RequesttemplatesEntry$json = {
  '1': 'RequesttemplatesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Integration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List integrationDescriptor = $convert.base64Decode(
    'CgtJbnRlZ3JhdGlvbhIyChJjYWNoZWtleXBhcmFtZXRlcnMYoezI5QEgAygJUhJjYWNoZWtleX'
    'BhcmFtZXRlcnMSKQoOY2FjaGVuYW1lc3BhY2UYoaHKKCABKAlSDmNhY2hlbmFtZXNwYWNlEiYK'
    'DGNvbm5lY3Rpb25pZBi9w8vWASABKAlSDGNvbm5lY3Rpb25pZBJGCg5jb25uZWN0aW9udHlwZR'
    'jyoaugASABKA4yGi5hcGlnYXRld2F5LkNvbm5lY3Rpb25UeXBlUg5jb25uZWN0aW9udHlwZRJR'
    'Cg9jb250ZW50aGFuZGxpbmcY8PKe/gEgASgOMiMuYXBpZ2F0ZXdheS5Db250ZW50SGFuZGxpbm'
    'dTdHJhdGVneVIPY29udGVudGhhbmRsaW5nEiMKC2NyZWRlbnRpYWxzGNK39kcgASgJUgtjcmVk'
    'ZW50aWFscxIhCgpodHRwbWV0aG9kGPHz+zYgASgJUgpodHRwbWV0aG9kEmkKFGludGVncmF0aW'
    '9ucmVzcG9uc2VzGPD/qrgBIAMoCzIxLmFwaWdhdGV3YXkuSW50ZWdyYXRpb24uSW50ZWdyYXRp'
    'b25yZXNwb25zZXNFbnRyeVIUaW50ZWdyYXRpb25yZXNwb25zZXMSLwoRaW50ZWdyYXRpb250YX'
    'JnZXQY8Yi1CCABKAlSEWludGVncmF0aW9udGFyZ2V0EjQKE3Bhc3N0aHJvdWdoYmVoYXZpb3IY'
    '7MSZlAEgASgJUhNwYXNzdGhyb3VnaGJlaGF2aW9yEmAKEXJlcXVlc3RwYXJhbWV0ZXJzGKPzz/'
    'kBIAMoCzIuLmFwaWdhdGV3YXkuSW50ZWdyYXRpb24uUmVxdWVzdHBhcmFtZXRlcnNFbnRyeVIR'
    'cmVxdWVzdHBhcmFtZXRlcnMSXQoQcmVxdWVzdHRlbXBsYXRlcxjm+4OfASADKAsyLS5hcGlnYX'
    'Rld2F5LkludGVncmF0aW9uLlJlcXVlc3R0ZW1wbGF0ZXNFbnRyeVIQcmVxdWVzdHRlbXBsYXRl'
    'cxJYChRyZXNwb25zZXRyYW5zZmVybW9kZRjD2OnaASABKA4yIC5hcGlnYXRld2F5LlJlc3Bvbn'
    'NlVHJhbnNmZXJNb2RlUhRyZXNwb25zZXRyYW5zZmVybW9kZRIsCg90aW1lb3V0aW5taWxsaXMY'
    'hqOttAEgASgFUg90aW1lb3V0aW5taWxsaXMSNgoJdGxzY29uZmlnGIXK+TMgASgLMhUuYXBpZ2'
    'F0ZXdheS5UbHNDb25maWdSCXRsc2NvbmZpZxIzCgR0eXBlGM7in4kBIAEoDjIbLmFwaWdhdGV3'
    'YXkuSW50ZWdyYXRpb25UeXBlUgR0eXBlEhQKA3VyaRj+p728ASABKAlSA3VyaRpoChlJbnRlZ3'
    'JhdGlvbnJlc3BvbnNlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EjUKBXZhbHVlGAIgASgLMh8u'
    'YXBpZ2F0ZXdheS5JbnRlZ3JhdGlvblJlc3BvbnNlUgV2YWx1ZToCOAEaRAoWUmVxdWVzdHBhcm'
    'FtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB'
    'GkMKFVJlcXVlc3R0ZW1wbGF0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIA'
    'EoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use integrationResponseDescriptor instead')
const IntegrationResponse$json = {
  '1': 'IntegrationResponse',
  '2': [
    {
      '1': 'contenthandling',
      '3': 533182832,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ContentHandlingStrategy',
      '10': 'contenthandling'
    },
    {
      '1': 'responseparameters',
      '3': 64271839,
      '4': 3,
      '5': 11,
      '6': '.apigateway.IntegrationResponse.ResponseparametersEntry',
      '10': 'responseparameters'
    },
    {
      '1': 'responsetemplates',
      '3': 107376570,
      '4': 3,
      '5': 11,
      '6': '.apigateway.IntegrationResponse.ResponsetemplatesEntry',
      '10': 'responsetemplates'
    },
    {
      '1': 'selectionpattern',
      '3': 470634042,
      '4': 1,
      '5': 9,
      '10': 'selectionpattern'
    },
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
  '3': [
    IntegrationResponse_ResponseparametersEntry$json,
    IntegrationResponse_ResponsetemplatesEntry$json
  ],
};

@$core.Deprecated('Use integrationResponseDescriptor instead')
const IntegrationResponse_ResponseparametersEntry$json = {
  '1': 'ResponseparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use integrationResponseDescriptor instead')
const IntegrationResponse_ResponsetemplatesEntry$json = {
  '1': 'ResponsetemplatesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `IntegrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List integrationResponseDescriptor = $convert.base64Decode(
    'ChNJbnRlZ3JhdGlvblJlc3BvbnNlElEKD2NvbnRlbnRoYW5kbGluZxjw8p7+ASABKA4yIy5hcG'
    'lnYXRld2F5LkNvbnRlbnRIYW5kbGluZ1N0cmF0ZWd5Ug9jb250ZW50aGFuZGxpbmcSagoScmVz'
    'cG9uc2VwYXJhbWV0ZXJzGN/r0h4gAygLMjcuYXBpZ2F0ZXdheS5JbnRlZ3JhdGlvblJlc3Bvbn'
    'NlLlJlc3BvbnNlcGFyYW1ldGVyc0VudHJ5UhJyZXNwb25zZXBhcmFtZXRlcnMSZwoRcmVzcG9u'
    'c2V0ZW1wbGF0ZXMYut+ZMyADKAsyNi5hcGlnYXRld2F5LkludGVncmF0aW9uUmVzcG9uc2UuUm'
    'VzcG9uc2V0ZW1wbGF0ZXNFbnRyeVIRcmVzcG9uc2V0ZW1wbGF0ZXMSLgoQc2VsZWN0aW9ucGF0'
    'dGVybhi6nLXgASABKAlSEHNlbGVjdGlvbnBhdHRlcm4SIgoKc3RhdHVzY29kZRifgd+OASABKA'
    'lSCnN0YXR1c2NvZGUaRQoXUmVzcG9uc2VwYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNr'
    'ZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpEChZSZXNwb25zZXRlbXBsYXRlc0VudHJ5Eh'
    'AKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use limitExceededExceptionDescriptor instead')
const LimitExceededException$json = {
  '1': 'LimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'retryafterseconds',
      '3': 436039555,
      '4': 1,
      '5': 9,
      '10': 'retryafterseconds'
    },
  ],
};

/// Descriptor for `LimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChZMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2USMA'
        'oRcmV0cnlhZnRlcnNlY29uZHMYg9/1zwEgASgJUhFyZXRyeWFmdGVyc2Vjb25kcw==');

@$core.Deprecated('Use methodDescriptor instead')
const Method$json = {
  '1': 'Method',
  '2': [
    {
      '1': 'apikeyrequired',
      '3': 435360152,
      '4': 1,
      '5': 8,
      '10': 'apikeyrequired'
    },
    {
      '1': 'authorizationscopes',
      '3': 423149932,
      '4': 3,
      '5': 9,
      '10': 'authorizationscopes'
    },
    {
      '1': 'authorizationtype',
      '3': 387986911,
      '4': 1,
      '5': 9,
      '10': 'authorizationtype'
    },
    {'1': 'authorizerid', '3': 111773148, '4': 1, '5': 9, '10': 'authorizerid'},
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'methodintegration',
      '3': 518245059,
      '4': 1,
      '5': 11,
      '6': '.apigateway.Integration',
      '10': 'methodintegration'
    },
    {
      '1': 'methodresponses',
      '3': 231818421,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Method.MethodresponsesEntry',
      '10': 'methodresponses'
    },
    {
      '1': 'operationname',
      '3': 178909574,
      '4': 1,
      '5': 9,
      '10': 'operationname'
    },
    {
      '1': 'requestmodels',
      '3': 397252853,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Method.RequestmodelsEntry',
      '10': 'requestmodels'
    },
    {
      '1': 'requestparameters',
      '3': 523499939,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Method.RequestparametersEntry',
      '10': 'requestparameters'
    },
    {
      '1': 'requestvalidatorid',
      '3': 517546134,
      '4': 1,
      '5': 9,
      '10': 'requestvalidatorid'
    },
  ],
  '3': [
    Method_MethodresponsesEntry$json,
    Method_RequestmodelsEntry$json,
    Method_RequestparametersEntry$json
  ],
};

@$core.Deprecated('Use methodDescriptor instead')
const Method_MethodresponsesEntry$json = {
  '1': 'MethodresponsesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.apigateway.MethodResponse',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use methodDescriptor instead')
const Method_RequestmodelsEntry$json = {
  '1': 'RequestmodelsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use methodDescriptor instead')
const Method_RequestparametersEntry$json = {
  '1': 'RequestparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 8, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Method`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List methodDescriptor = $convert.base64Decode(
    'CgZNZXRob2QSKgoOYXBpa2V5cmVxdWlyZWQYmKPMzwEgASgIUg5hcGlrZXlyZXF1aXJlZBI0Ch'
    'NhdXRob3JpemF0aW9uc2NvcGVzGOyC48kBIAMoCVITYXV0aG9yaXphdGlvbnNjb3BlcxIwChFh'
    'dXRob3JpemF0aW9udHlwZRjf64C5ASABKAlSEWF1dGhvcml6YXRpb250eXBlEiUKDGF1dGhvcm'
    'l6ZXJpZBjci6Y1IAEoCVIMYXV0aG9yaXplcmlkEiEKCmh0dHBtZXRob2QY8fP7NiABKAlSCmh0'
    'dHBtZXRob2QSSQoRbWV0aG9kaW50ZWdyYXRpb24Yw5WP9wEgASgLMhcuYXBpZ2F0ZXdheS5Jbn'
    'RlZ3JhdGlvblIRbWV0aG9kaW50ZWdyYXRpb24SVAoPbWV0aG9kcmVzcG9uc2VzGLWJxW4gAygL'
    'MicuYXBpZ2F0ZXdheS5NZXRob2QuTWV0aG9kcmVzcG9uc2VzRW50cnlSD21ldGhvZHJlc3Bvbn'
    'NlcxInCg1vcGVyYXRpb25uYW1lGIbjp1UgASgJUg1vcGVyYXRpb25uYW1lEk8KDXJlcXVlc3Rt'
    'b2RlbHMY9bG2vQEgAygLMiUuYXBpZ2F0ZXdheS5NZXRob2QuUmVxdWVzdG1vZGVsc0VudHJ5Ug'
    '1yZXF1ZXN0bW9kZWxzElsKEXJlcXVlc3RwYXJhbWV0ZXJzGKPzz/kBIAMoCzIpLmFwaWdhdGV3'
    'YXkuTWV0aG9kLlJlcXVlc3RwYXJhbWV0ZXJzRW50cnlSEXJlcXVlc3RwYXJhbWV0ZXJzEjIKEn'
    'JlcXVlc3R2YWxpZGF0b3JpZBiWweT2ASABKAlSEnJlcXVlc3R2YWxpZGF0b3JpZBpeChRNZXRo'
    'b2RyZXNwb25zZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIwCgV2YWx1ZRgCIAEoCzIaLmFwaW'
    'dhdGV3YXkuTWV0aG9kUmVzcG9uc2VSBXZhbHVlOgI4ARpAChJSZXF1ZXN0bW9kZWxzRW50cnkS'
    'EAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpEChZSZXF1ZXN0cG'
    'FyYW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgIUgV2YWx1ZToC'
    'OAE=');

@$core.Deprecated('Use methodResponseDescriptor instead')
const MethodResponse$json = {
  '1': 'MethodResponse',
  '2': [
    {
      '1': 'responsemodels',
      '3': 356574313,
      '4': 3,
      '5': 11,
      '6': '.apigateway.MethodResponse.ResponsemodelsEntry',
      '10': 'responsemodels'
    },
    {
      '1': 'responseparameters',
      '3': 64271839,
      '4': 3,
      '5': 11,
      '6': '.apigateway.MethodResponse.ResponseparametersEntry',
      '10': 'responseparameters'
    },
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
  '3': [
    MethodResponse_ResponsemodelsEntry$json,
    MethodResponse_ResponseparametersEntry$json
  ],
};

@$core.Deprecated('Use methodResponseDescriptor instead')
const MethodResponse_ResponsemodelsEntry$json = {
  '1': 'ResponsemodelsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use methodResponseDescriptor instead')
const MethodResponse_ResponseparametersEntry$json = {
  '1': 'ResponseparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 8, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `MethodResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List methodResponseDescriptor = $convert.base64Decode(
    'Cg5NZXRob2RSZXNwb25zZRJaCg5yZXNwb25zZW1vZGVscxjpyIOqASADKAsyLi5hcGlnYXRld2'
    'F5Lk1ldGhvZFJlc3BvbnNlLlJlc3BvbnNlbW9kZWxzRW50cnlSDnJlc3BvbnNlbW9kZWxzEmUK'
    'EnJlc3BvbnNlcGFyYW1ldGVycxjf69IeIAMoCzIyLmFwaWdhdGV3YXkuTWV0aG9kUmVzcG9uc2'
    'UuUmVzcG9uc2VwYXJhbWV0ZXJzRW50cnlSEnJlc3BvbnNlcGFyYW1ldGVycxIiCgpzdGF0dXNj'
    'b2RlGJ+B344BIAEoCVIKc3RhdHVzY29kZRpBChNSZXNwb25zZW1vZGVsc0VudHJ5EhAKA2tleR'
    'gBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaRQoXUmVzcG9uc2VwYXJhbWV0'
    'ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAhSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use methodSettingDescriptor instead')
const MethodSetting$json = {
  '1': 'MethodSetting',
  '2': [
    {
      '1': 'cachedataencrypted',
      '3': 481075060,
      '4': 1,
      '5': 8,
      '10': 'cachedataencrypted'
    },
    {
      '1': 'cachettlinseconds',
      '3': 79996982,
      '4': 1,
      '5': 5,
      '10': 'cachettlinseconds'
    },
    {
      '1': 'cachingenabled',
      '3': 489524028,
      '4': 1,
      '5': 8,
      '10': 'cachingenabled'
    },
    {
      '1': 'datatraceenabled',
      '3': 363519852,
      '4': 1,
      '5': 8,
      '10': 'datatraceenabled'
    },
    {'1': 'logginglevel', '3': 59396637, '4': 1, '5': 9, '10': 'logginglevel'},
    {
      '1': 'metricsenabled',
      '3': 142460292,
      '4': 1,
      '5': 8,
      '10': 'metricsenabled'
    },
    {
      '1': 'requireauthorizationforcachecontrol',
      '3': 529394912,
      '4': 1,
      '5': 8,
      '10': 'requireauthorizationforcachecontrol'
    },
    {
      '1': 'throttlingburstlimit',
      '3': 402901688,
      '4': 1,
      '5': 5,
      '10': 'throttlingburstlimit'
    },
    {
      '1': 'throttlingratelimit',
      '3': 371718088,
      '4': 1,
      '5': 1,
      '10': 'throttlingratelimit'
    },
    {
      '1': 'unauthorizedcachecontrolheaderstrategy',
      '3': 476741277,
      '4': 1,
      '5': 14,
      '6': '.apigateway.UnauthorizedCacheControlHeaderStrategy',
      '10': 'unauthorizedcachecontrolheaderstrategy'
    },
  ],
};

/// Descriptor for `MethodSetting`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List methodSettingDescriptor = $convert.base64Decode(
    'Cg1NZXRob2RTZXR0aW5nEjIKEmNhY2hlZGF0YWVuY3J5cHRlZBj0vrLlASABKAhSEmNhY2hlZG'
    'F0YWVuY3J5cHRlZBIvChFjYWNoZXR0bGluc2Vjb25kcxi20JImIAEoBVIRY2FjaGV0dGxpbnNl'
    'Y29uZHMSKgoOY2FjaGluZ2VuYWJsZWQYvJa26QEgASgIUg5jYWNoaW5nZW5hYmxlZBIuChBkYX'
    'RhdHJhY2VlbmFibGVkGOy+q60BIAEoCFIQZGF0YXRyYWNlZW5hYmxlZBIlCgxsb2dnaW5nbGV2'
    'ZWwYnaSpHCABKAlSDGxvZ2dpbmdsZXZlbBIpCg5tZXRyaWNzZW5hYmxlZBiEi/dDIAEoCFIObW'
    'V0cmljc2VuYWJsZWQSVAojcmVxdWlyZWF1dGhvcml6YXRpb25mb3JjYWNoZWNvbnRyb2wY4Nm3'
    '/AEgASgIUiNyZXF1aXJlYXV0aG9yaXphdGlvbmZvcmNhY2hlY29udHJvbBI2ChR0aHJvdHRsaW'
    '5nYnVyc3RsaW1pdBi4lY/AASABKAVSFHRocm90dGxpbmdidXJzdGxpbWl0EjQKE3Rocm90dGxp'
    'bmdyYXRlbGltaXQYyO+fsQEgASgBUhN0aHJvdHRsaW5ncmF0ZWxpbWl0Eo4BCiZ1bmF1dGhvcm'
    'l6ZWRjYWNoZWNvbnRyb2xoZWFkZXJzdHJhdGVneRid/anjASABKA4yMi5hcGlnYXRld2F5LlVu'
    'YXV0aG9yaXplZENhY2hlQ29udHJvbEhlYWRlclN0cmF0ZWd5UiZ1bmF1dGhvcml6ZWRjYWNoZW'
    'NvbnRyb2xoZWFkZXJzdHJhdGVneQ==');

@$core.Deprecated('Use methodSnapshotDescriptor instead')
const MethodSnapshot$json = {
  '1': 'MethodSnapshot',
  '2': [
    {
      '1': 'apikeyrequired',
      '3': 435360152,
      '4': 1,
      '5': 8,
      '10': 'apikeyrequired'
    },
    {
      '1': 'authorizationtype',
      '3': 387986911,
      '4': 1,
      '5': 9,
      '10': 'authorizationtype'
    },
  ],
};

/// Descriptor for `MethodSnapshot`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List methodSnapshotDescriptor = $convert.base64Decode(
    'Cg5NZXRob2RTbmFwc2hvdBIqCg5hcGlrZXlyZXF1aXJlZBiYo8zPASABKAhSDmFwaWtleXJlcX'
    'VpcmVkEjAKEWF1dGhvcml6YXRpb250eXBlGN/rgLkBIAEoCVIRYXV0aG9yaXphdGlvbnR5cGU=');

@$core.Deprecated('Use modelDescriptor instead')
const Model$json = {
  '1': 'Model',
  '2': [
    {'1': 'contenttype', '3': 281764659, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'schema', '3': 310182711, '4': 1, '5': 9, '10': 'schema'},
  ],
};

/// Descriptor for `Model`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelDescriptor = $convert.base64Decode(
    'CgVNb2RlbBIkCgtjb250ZW50dHlwZRizxq2GASABKAlSC2NvbnRlbnR0eXBlEiQKC2Rlc2NyaX'
    'B0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SEgoCaWQY4dXhuQEgASgJUgJpZBIVCgRuYW1l'
    'GOf75mkgASgJUgRuYW1lEhoKBnNjaGVtYRi3hvSTASABKAlSBnNjaGVtYQ==');

@$core.Deprecated('Use modelsDescriptor instead')
const Models$json = {
  '1': 'Models',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Model',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `Models`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelsDescriptor = $convert.base64Decode(
    'CgZNb2RlbHMSKwoFaXRlbXMYkOfk0wEgAygLMhEuYXBpZ2F0ZXdheS5Nb2RlbFIFaXRlbXMSHg'
    'oIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbg==');

@$core.Deprecated('Use mutualTlsAuthenticationDescriptor instead')
const MutualTlsAuthentication$json = {
  '1': 'MutualTlsAuthentication',
  '2': [
    {
      '1': 'truststoreuri',
      '3': 120246545,
      '4': 1,
      '5': 9,
      '10': 'truststoreuri'
    },
    {
      '1': 'truststoreversion',
      '3': 119080291,
      '4': 1,
      '5': 9,
      '10': 'truststoreversion'
    },
    {
      '1': 'truststorewarnings',
      '3': 420536820,
      '4': 3,
      '5': 9,
      '10': 'truststorewarnings'
    },
  ],
};

/// Descriptor for `MutualTlsAuthentication`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mutualTlsAuthenticationDescriptor = $convert.base64Decode(
    'ChdNdXR1YWxUbHNBdXRoZW50aWNhdGlvbhInCg10cnVzdHN0b3JldXJpGJGiqzkgASgJUg10cn'
    'VzdHN0b3JldXJpEi8KEXRydXN0c3RvcmV2ZXJzaW9uGOOK5DggASgJUhF0cnVzdHN0b3JldmVy'
    'c2lvbhIyChJ0cnVzdHN0b3Jld2FybmluZ3MY9MPDyAEgAygJUhJ0cnVzdHN0b3Jld2FybmluZ3'
    'M=');

@$core.Deprecated('Use mutualTlsAuthenticationInputDescriptor instead')
const MutualTlsAuthenticationInput$json = {
  '1': 'MutualTlsAuthenticationInput',
  '2': [
    {
      '1': 'truststoreuri',
      '3': 120246545,
      '4': 1,
      '5': 9,
      '10': 'truststoreuri'
    },
    {
      '1': 'truststoreversion',
      '3': 119080291,
      '4': 1,
      '5': 9,
      '10': 'truststoreversion'
    },
  ],
};

/// Descriptor for `MutualTlsAuthenticationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mutualTlsAuthenticationInputDescriptor =
    $convert.base64Decode(
        'ChxNdXR1YWxUbHNBdXRoZW50aWNhdGlvbklucHV0EicKDXRydXN0c3RvcmV1cmkYkaKrOSABKA'
        'lSDXRydXN0c3RvcmV1cmkSLwoRdHJ1c3RzdG9yZXZlcnNpb24Y44rkOCABKAlSEXRydXN0c3Rv'
        'cmV2ZXJzaW9u');

@$core.Deprecated('Use notFoundExceptionDescriptor instead')
const NotFoundException$json = {
  '1': 'NotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notFoundExceptionDescriptor = $convert.base64Decode(
    'ChFOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use patchOperationDescriptor instead')
const PatchOperation$json = {
  '1': 'PatchOperation',
  '2': [
    {'1': 'from', '3': 365789302, '4': 1, '5': 9, '10': 'from'},
    {
      '1': 'op',
      '3': 523513003,
      '4': 1,
      '5': 14,
      '6': '.apigateway.Op',
      '10': 'op'
    },
    {'1': 'path', '3': 75975991, '4': 1, '5': 9, '10': 'path'},
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `PatchOperation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchOperationDescriptor = $convert.base64Decode(
    'Cg5QYXRjaE9wZXJhdGlvbhIWCgRmcm9tGPaAtq4BIAEoCVIEZnJvbRIiCgJvcBir2dD5ASABKA'
    '4yDi5hcGlnYXRld2F5Lk9wUgJvcBIVCgRwYXRoGLeanSQgASgJUgRwYXRoEhcKBXZhbHVlGMun'
    '+xIgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use putGatewayResponseRequestDescriptor instead')
const PutGatewayResponseRequest$json = {
  '1': 'PutGatewayResponseRequest',
  '2': [
    {
      '1': 'responseparameters',
      '3': 64271839,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutGatewayResponseRequest.ResponseparametersEntry',
      '10': 'responseparameters'
    },
    {
      '1': 'responsetemplates',
      '3': 107376570,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutGatewayResponseRequest.ResponsetemplatesEntry',
      '10': 'responsetemplates'
    },
    {
      '1': 'responsetype',
      '3': 377935935,
      '4': 1,
      '5': 14,
      '6': '.apigateway.GatewayResponseType',
      '10': 'responsetype'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
  '3': [
    PutGatewayResponseRequest_ResponseparametersEntry$json,
    PutGatewayResponseRequest_ResponsetemplatesEntry$json
  ],
};

@$core.Deprecated('Use putGatewayResponseRequestDescriptor instead')
const PutGatewayResponseRequest_ResponseparametersEntry$json = {
  '1': 'ResponseparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putGatewayResponseRequestDescriptor instead')
const PutGatewayResponseRequest_ResponsetemplatesEntry$json = {
  '1': 'ResponsetemplatesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutGatewayResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putGatewayResponseRequestDescriptor = $convert.base64Decode(
    'ChlQdXRHYXRld2F5UmVzcG9uc2VSZXF1ZXN0EnAKEnJlc3BvbnNlcGFyYW1ldGVycxjf69IeIA'
    'MoCzI9LmFwaWdhdGV3YXkuUHV0R2F0ZXdheVJlc3BvbnNlUmVxdWVzdC5SZXNwb25zZXBhcmFt'
    'ZXRlcnNFbnRyeVIScmVzcG9uc2VwYXJhbWV0ZXJzEm0KEXJlc3BvbnNldGVtcGxhdGVzGLrfmT'
    'MgAygLMjwuYXBpZ2F0ZXdheS5QdXRHYXRld2F5UmVzcG9uc2VSZXF1ZXN0LlJlc3BvbnNldGVt'
    'cGxhdGVzRW50cnlSEXJlc3BvbnNldGVtcGxhdGVzEkcKDHJlc3BvbnNldHlwZRi/sJu0ASABKA'
    '4yHy5hcGlnYXRld2F5LkdhdGV3YXlSZXNwb25zZVR5cGVSDHJlc3BvbnNldHlwZRIgCglyZXN0'
    'YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQSIgoKc3RhdHVzY29kZRifgd+OASABKAlSCnN0YX'
    'R1c2NvZGUaRQoXUmVzcG9uc2VwYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoF'
    'dmFsdWUYAiABKAlSBXZhbHVlOgI4ARpEChZSZXNwb25zZXRlbXBsYXRlc0VudHJ5EhAKA2tleR'
    'gBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use putIntegrationRequestDescriptor instead')
const PutIntegrationRequest$json = {
  '1': 'PutIntegrationRequest',
  '2': [
    {
      '1': 'cachekeyparameters',
      '3': 481441313,
      '4': 3,
      '5': 9,
      '10': 'cachekeyparameters'
    },
    {
      '1': 'cachenamespace',
      '3': 85102753,
      '4': 1,
      '5': 9,
      '10': 'cachenamespace'
    },
    {'1': 'connectionid', '3': 450027965, '4': 1, '5': 9, '10': 'connectionid'},
    {
      '1': 'connectiontype',
      '3': 336253170,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ConnectionType',
      '10': 'connectiontype'
    },
    {
      '1': 'contenthandling',
      '3': 533182832,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ContentHandlingStrategy',
      '10': 'contenthandling'
    },
    {'1': 'credentials', '3': 150838226, '4': 1, '5': 9, '10': 'credentials'},
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'integrationhttpmethod',
      '3': 355314073,
      '4': 1,
      '5': 9,
      '10': 'integrationhttpmethod'
    },
    {
      '1': 'integrationtarget',
      '3': 17646705,
      '4': 1,
      '5': 9,
      '10': 'integrationtarget'
    },
    {
      '1': 'passthroughbehavior',
      '3': 310796908,
      '4': 1,
      '5': 9,
      '10': 'passthroughbehavior'
    },
    {
      '1': 'requestparameters',
      '3': 523499939,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutIntegrationRequest.RequestparametersEntry',
      '10': 'requestparameters'
    },
    {
      '1': 'requesttemplates',
      '3': 333512166,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutIntegrationRequest.RequesttemplatesEntry',
      '10': 'requesttemplates'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'responsetransfermode',
      '3': 458910787,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ResponseTransferMode',
      '10': 'responsetransfermode'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'timeoutinmillis',
      '3': 378229126,
      '4': 1,
      '5': 5,
      '10': 'timeoutinmillis'
    },
    {
      '1': 'tlsconfig',
      '3': 108946693,
      '4': 1,
      '5': 11,
      '6': '.apigateway.TlsConfig',
      '10': 'tlsconfig'
    },
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.apigateway.IntegrationType',
      '10': 'type'
    },
    {'1': 'uri', '3': 395269118, '4': 1, '5': 9, '10': 'uri'},
  ],
  '3': [
    PutIntegrationRequest_RequestparametersEntry$json,
    PutIntegrationRequest_RequesttemplatesEntry$json
  ],
};

@$core.Deprecated('Use putIntegrationRequestDescriptor instead')
const PutIntegrationRequest_RequestparametersEntry$json = {
  '1': 'RequestparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putIntegrationRequestDescriptor instead')
const PutIntegrationRequest_RequesttemplatesEntry$json = {
  '1': 'RequesttemplatesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putIntegrationRequestDescriptor = $convert.base64Decode(
    'ChVQdXRJbnRlZ3JhdGlvblJlcXVlc3QSMgoSY2FjaGVrZXlwYXJhbWV0ZXJzGKHsyOUBIAMoCV'
    'ISY2FjaGVrZXlwYXJhbWV0ZXJzEikKDmNhY2hlbmFtZXNwYWNlGKGhyiggASgJUg5jYWNoZW5h'
    'bWVzcGFjZRImCgxjb25uZWN0aW9uaWQYvcPL1gEgASgJUgxjb25uZWN0aW9uaWQSRgoOY29ubm'
    'VjdGlvbnR5cGUY8qGroAEgASgOMhouYXBpZ2F0ZXdheS5Db25uZWN0aW9uVHlwZVIOY29ubmVj'
    'dGlvbnR5cGUSUQoPY29udGVudGhhbmRsaW5nGPDynv4BIAEoDjIjLmFwaWdhdGV3YXkuQ29udG'
    'VudEhhbmRsaW5nU3RyYXRlZ3lSD2NvbnRlbnRoYW5kbGluZxIjCgtjcmVkZW50aWFscxjSt/ZH'
    'IAEoCVILY3JlZGVudGlhbHMSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaHR0cG1ldGhvZBI4Ch'
    'VpbnRlZ3JhdGlvbmh0dHBtZXRob2QYmdO2qQEgASgJUhVpbnRlZ3JhdGlvbmh0dHBtZXRob2QS'
    'LwoRaW50ZWdyYXRpb250YXJnZXQY8Yi1CCABKAlSEWludGVncmF0aW9udGFyZ2V0EjQKE3Bhc3'
    'N0aHJvdWdoYmVoYXZpb3IY7MSZlAEgASgJUhNwYXNzdGhyb3VnaGJlaGF2aW9yEmoKEXJlcXVl'
    'c3RwYXJhbWV0ZXJzGKPzz/kBIAMoCzI4LmFwaWdhdGV3YXkuUHV0SW50ZWdyYXRpb25SZXF1ZX'
    'N0LlJlcXVlc3RwYXJhbWV0ZXJzRW50cnlSEXJlcXVlc3RwYXJhbWV0ZXJzEmcKEHJlcXVlc3R0'
    'ZW1wbGF0ZXMY5vuDnwEgAygLMjcuYXBpZ2F0ZXdheS5QdXRJbnRlZ3JhdGlvblJlcXVlc3QuUm'
    'VxdWVzdHRlbXBsYXRlc0VudHJ5UhByZXF1ZXN0dGVtcGxhdGVzEiIKCnJlc291cmNlaWQYsb2J'
    'mAEgASgJUgpyZXNvdXJjZWlkElgKFHJlc3BvbnNldHJhbnNmZXJtb2RlGMPY6doBIAEoDjIgLm'
    'FwaWdhdGV3YXkuUmVzcG9uc2VUcmFuc2Zlck1vZGVSFHJlc3BvbnNldHJhbnNmZXJtb2RlEiAK'
    'CXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZBIsCg90aW1lb3V0aW5taWxsaXMYhqOttA'
    'EgASgFUg90aW1lb3V0aW5taWxsaXMSNgoJdGxzY29uZmlnGIXK+TMgASgLMhUuYXBpZ2F0ZXdh'
    'eS5UbHNDb25maWdSCXRsc2NvbmZpZxIzCgR0eXBlGM7in4kBIAEoDjIbLmFwaWdhdGV3YXkuSW'
    '50ZWdyYXRpb25UeXBlUgR0eXBlEhQKA3VyaRj+p728ASABKAlSA3VyaRpEChZSZXF1ZXN0cGFy'
    'YW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOA'
    'EaQwoVUmVxdWVzdHRlbXBsYXRlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIg'
    'ASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use putIntegrationResponseRequestDescriptor instead')
const PutIntegrationResponseRequest$json = {
  '1': 'PutIntegrationResponseRequest',
  '2': [
    {
      '1': 'contenthandling',
      '3': 533182832,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ContentHandlingStrategy',
      '10': 'contenthandling'
    },
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'responseparameters',
      '3': 64271839,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutIntegrationResponseRequest.ResponseparametersEntry',
      '10': 'responseparameters'
    },
    {
      '1': 'responsetemplates',
      '3': 107376570,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutIntegrationResponseRequest.ResponsetemplatesEntry',
      '10': 'responsetemplates'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'selectionpattern',
      '3': 470634042,
      '4': 1,
      '5': 9,
      '10': 'selectionpattern'
    },
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
  '3': [
    PutIntegrationResponseRequest_ResponseparametersEntry$json,
    PutIntegrationResponseRequest_ResponsetemplatesEntry$json
  ],
};

@$core.Deprecated('Use putIntegrationResponseRequestDescriptor instead')
const PutIntegrationResponseRequest_ResponseparametersEntry$json = {
  '1': 'ResponseparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putIntegrationResponseRequestDescriptor instead')
const PutIntegrationResponseRequest_ResponsetemplatesEntry$json = {
  '1': 'ResponsetemplatesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutIntegrationResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putIntegrationResponseRequestDescriptor = $convert.base64Decode(
    'Ch1QdXRJbnRlZ3JhdGlvblJlc3BvbnNlUmVxdWVzdBJRCg9jb250ZW50aGFuZGxpbmcY8PKe/g'
    'EgASgOMiMuYXBpZ2F0ZXdheS5Db250ZW50SGFuZGxpbmdTdHJhdGVneVIPY29udGVudGhhbmRs'
    'aW5nEiEKCmh0dHBtZXRob2QY8fP7NiABKAlSCmh0dHBtZXRob2QSIgoKcmVzb3VyY2VpZBixvY'
    'mYASABKAlSCnJlc291cmNlaWQSdAoScmVzcG9uc2VwYXJhbWV0ZXJzGN/r0h4gAygLMkEuYXBp'
    'Z2F0ZXdheS5QdXRJbnRlZ3JhdGlvblJlc3BvbnNlUmVxdWVzdC5SZXNwb25zZXBhcmFtZXRlcn'
    'NFbnRyeVIScmVzcG9uc2VwYXJhbWV0ZXJzEnEKEXJlc3BvbnNldGVtcGxhdGVzGLrfmTMgAygL'
    'MkAuYXBpZ2F0ZXdheS5QdXRJbnRlZ3JhdGlvblJlc3BvbnNlUmVxdWVzdC5SZXNwb25zZXRlbX'
    'BsYXRlc0VudHJ5UhFyZXNwb25zZXRlbXBsYXRlcxIgCglyZXN0YXBpaWQYmaSBtwEgASgJUgly'
    'ZXN0YXBpaWQSLgoQc2VsZWN0aW9ucGF0dGVybhi6nLXgASABKAlSEHNlbGVjdGlvbnBhdHRlcm'
    '4SIgoKc3RhdHVzY29kZRifgd+OASABKAlSCnN0YXR1c2NvZGUaRQoXUmVzcG9uc2VwYXJhbWV0'
    'ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpECh'
    'ZSZXNwb25zZXRlbXBsYXRlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJ'
    'UgV2YWx1ZToCOAE=');

@$core.Deprecated('Use putMethodRequestDescriptor instead')
const PutMethodRequest$json = {
  '1': 'PutMethodRequest',
  '2': [
    {
      '1': 'apikeyrequired',
      '3': 435360152,
      '4': 1,
      '5': 8,
      '10': 'apikeyrequired'
    },
    {
      '1': 'authorizationscopes',
      '3': 423149932,
      '4': 3,
      '5': 9,
      '10': 'authorizationscopes'
    },
    {
      '1': 'authorizationtype',
      '3': 387986911,
      '4': 1,
      '5': 9,
      '10': 'authorizationtype'
    },
    {'1': 'authorizerid', '3': 111773148, '4': 1, '5': 9, '10': 'authorizerid'},
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'operationname',
      '3': 178909574,
      '4': 1,
      '5': 9,
      '10': 'operationname'
    },
    {
      '1': 'requestmodels',
      '3': 397252853,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutMethodRequest.RequestmodelsEntry',
      '10': 'requestmodels'
    },
    {
      '1': 'requestparameters',
      '3': 523499939,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutMethodRequest.RequestparametersEntry',
      '10': 'requestparameters'
    },
    {
      '1': 'requestvalidatorid',
      '3': 517546134,
      '4': 1,
      '5': 9,
      '10': 'requestvalidatorid'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
  '3': [
    PutMethodRequest_RequestmodelsEntry$json,
    PutMethodRequest_RequestparametersEntry$json
  ],
};

@$core.Deprecated('Use putMethodRequestDescriptor instead')
const PutMethodRequest_RequestmodelsEntry$json = {
  '1': 'RequestmodelsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putMethodRequestDescriptor instead')
const PutMethodRequest_RequestparametersEntry$json = {
  '1': 'RequestparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 8, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutMethodRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putMethodRequestDescriptor = $convert.base64Decode(
    'ChBQdXRNZXRob2RSZXF1ZXN0EioKDmFwaWtleXJlcXVpcmVkGJijzM8BIAEoCFIOYXBpa2V5cm'
    'VxdWlyZWQSNAoTYXV0aG9yaXphdGlvbnNjb3BlcxjsguPJASADKAlSE2F1dGhvcml6YXRpb25z'
    'Y29wZXMSMAoRYXV0aG9yaXphdGlvbnR5cGUY3+uAuQEgASgJUhFhdXRob3JpemF0aW9udHlwZR'
    'IlCgxhdXRob3JpemVyaWQY3IumNSABKAlSDGF1dGhvcml6ZXJpZBIhCgpodHRwbWV0aG9kGPHz'
    '+zYgASgJUgpodHRwbWV0aG9kEicKDW9wZXJhdGlvbm5hbWUYhuOnVSABKAlSDW9wZXJhdGlvbm'
    '5hbWUSWQoNcmVxdWVzdG1vZGVscxj1sba9ASADKAsyLy5hcGlnYXRld2F5LlB1dE1ldGhvZFJl'
    'cXVlc3QuUmVxdWVzdG1vZGVsc0VudHJ5Ug1yZXF1ZXN0bW9kZWxzEmUKEXJlcXVlc3RwYXJhbW'
    'V0ZXJzGKPzz/kBIAMoCzIzLmFwaWdhdGV3YXkuUHV0TWV0aG9kUmVxdWVzdC5SZXF1ZXN0cGFy'
    'YW1ldGVyc0VudHJ5UhFyZXF1ZXN0cGFyYW1ldGVycxIyChJyZXF1ZXN0dmFsaWRhdG9yaWQYls'
    'Hk9gEgASgJUhJyZXF1ZXN0dmFsaWRhdG9yaWQSIgoKcmVzb3VyY2VpZBixvYmYASABKAlSCnJl'
    'c291cmNlaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkGkAKElJlcXVlc3Rtb2'
    'RlbHNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGkQK'
    'FlJlcXVlc3RwYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKA'
    'hSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use putMethodResponseRequestDescriptor instead')
const PutMethodResponseRequest$json = {
  '1': 'PutMethodResponseRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'responsemodels',
      '3': 356574313,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutMethodResponseRequest.ResponsemodelsEntry',
      '10': 'responsemodels'
    },
    {
      '1': 'responseparameters',
      '3': 64271839,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutMethodResponseRequest.ResponseparametersEntry',
      '10': 'responseparameters'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
  '3': [
    PutMethodResponseRequest_ResponsemodelsEntry$json,
    PutMethodResponseRequest_ResponseparametersEntry$json
  ],
};

@$core.Deprecated('Use putMethodResponseRequestDescriptor instead')
const PutMethodResponseRequest_ResponsemodelsEntry$json = {
  '1': 'ResponsemodelsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putMethodResponseRequestDescriptor instead')
const PutMethodResponseRequest_ResponseparametersEntry$json = {
  '1': 'ResponseparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 8, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutMethodResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putMethodResponseRequestDescriptor = $convert.base64Decode(
    'ChhQdXRNZXRob2RSZXNwb25zZVJlcXVlc3QSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaHR0cG'
    '1ldGhvZBIiCgpyZXNvdXJjZWlkGLG9iZgBIAEoCVIKcmVzb3VyY2VpZBJkCg5yZXNwb25zZW1v'
    'ZGVscxjpyIOqASADKAsyOC5hcGlnYXRld2F5LlB1dE1ldGhvZFJlc3BvbnNlUmVxdWVzdC5SZX'
    'Nwb25zZW1vZGVsc0VudHJ5Ug5yZXNwb25zZW1vZGVscxJvChJyZXNwb25zZXBhcmFtZXRlcnMY'
    '3+vSHiADKAsyPC5hcGlnYXRld2F5LlB1dE1ldGhvZFJlc3BvbnNlUmVxdWVzdC5SZXNwb25zZX'
    'BhcmFtZXRlcnNFbnRyeVIScmVzcG9uc2VwYXJhbWV0ZXJzEiAKCXJlc3RhcGlpZBiZpIG3ASAB'
    'KAlSCXJlc3RhcGlpZBIiCgpzdGF0dXNjb2RlGJ+B344BIAEoCVIKc3RhdHVzY29kZRpBChNSZX'
    'Nwb25zZW1vZGVsc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1'
    'ZToCOAEaRQoXUmVzcG9uc2VwYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdm'
    'FsdWUYAiABKAhSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use putRestApiRequestDescriptor instead')
const PutRestApiRequest$json = {
  '1': 'PutRestApiRequest',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'failonwarnings',
      '3': 434363958,
      '4': 1,
      '5': 8,
      '10': 'failonwarnings'
    },
    {
      '1': 'mode',
      '3': 208592915,
      '4': 1,
      '5': 14,
      '6': '.apigateway.PutMode',
      '10': 'mode'
    },
    {
      '1': 'parameters',
      '3': 145043162,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PutRestApiRequest.ParametersEntry',
      '10': 'parameters'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
  '3': [PutRestApiRequest_ParametersEntry$json],
};

@$core.Deprecated('Use putRestApiRequestDescriptor instead')
const PutRestApiRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutRestApiRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRestApiRequestDescriptor = $convert.base64Decode(
    'ChFQdXRSZXN0QXBpUmVxdWVzdBIcCgRib2R5GPbyqd0BIAEoDEIEiLUYAVIEYm9keRIqCg5mYW'
    'lsb253YXJuaW5ncxi2vI/PASABKAhSDmZhaWxvbndhcm5pbmdzEioKBG1vZGUYk8C7YyABKA4y'
    'Ey5hcGlnYXRld2F5LlB1dE1vZGVSBG1vZGUSUAoKcGFyYW1ldGVycxja3ZRFIAMoCzItLmFwaW'
    'dhdGV3YXkuUHV0UmVzdEFwaVJlcXVlc3QuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJzEiAK'
    'CXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZBo9Cg9QYXJhbWV0ZXJzRW50cnkSEAoDa2'
    'V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use quotaSettingsDescriptor instead')
const QuotaSettings$json = {
  '1': 'QuotaSettings',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 348705739, '4': 1, '5': 5, '10': 'offset'},
    {
      '1': 'period',
      '3': 432621317,
      '4': 1,
      '5': 14,
      '6': '.apigateway.QuotaPeriodType',
      '10': 'period'
    },
  ],
};

/// Descriptor for `QuotaSettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List quotaSettingsDescriptor = $convert.base64Decode(
    'Cg1RdW90YVNldHRpbmdzEhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSGgoGb2Zmc2V0GMuno6'
    'YBIAEoBVIGb2Zmc2V0EjcKBnBlcmlvZBiFjqXOASABKA4yGy5hcGlnYXRld2F5LlF1b3RhUGVy'
    'aW9kVHlwZVIGcGVyaW9k');

@$core.Deprecated(
    'Use rejectDomainNameAccessAssociationRequestDescriptor instead')
const RejectDomainNameAccessAssociationRequest$json = {
  '1': 'RejectDomainNameAccessAssociationRequest',
  '2': [
    {
      '1': 'domainnameaccessassociationarn',
      '3': 281017927,
      '4': 1,
      '5': 9,
      '10': 'domainnameaccessassociationarn'
    },
    {
      '1': 'domainnamearn',
      '3': 244019094,
      '4': 1,
      '5': 9,
      '10': 'domainnamearn'
    },
  ],
};

/// Descriptor for `RejectDomainNameAccessAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rejectDomainNameAccessAssociationRequestDescriptor =
    $convert.base64Decode(
        'CihSZWplY3REb21haW5OYW1lQWNjZXNzQXNzb2NpYXRpb25SZXF1ZXN0EkoKHmRvbWFpbm5hbW'
        'VhY2Nlc3Nhc3NvY2lhdGlvbmFybhjH/P+FASABKAlSHmRvbWFpbm5hbWVhY2Nlc3Nhc3NvY2lh'
        'dGlvbmFybhInCg1kb21haW5uYW1lYXJuGJbfrXQgASgJUg1kb21haW5uYW1lYXJu');

@$core.Deprecated('Use requestValidatorDescriptor instead')
const RequestValidator$json = {
  '1': 'RequestValidator',
  '2': [
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'validaterequestbody',
      '3': 397505841,
      '4': 1,
      '5': 8,
      '10': 'validaterequestbody'
    },
    {
      '1': 'validaterequestparameters',
      '3': 464035801,
      '4': 1,
      '5': 8,
      '10': 'validaterequestparameters'
    },
  ],
};

/// Descriptor for `RequestValidator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestValidatorDescriptor = $convert.base64Decode(
    'ChBSZXF1ZXN0VmFsaWRhdG9yEhIKAmlkGOHV4bkBIAEoCVICaWQSFQoEbmFtZRjn++ZpIAEoCV'
    'IEbmFtZRI0ChN2YWxpZGF0ZXJlcXVlc3Rib2R5GLHqxb0BIAEoCFITdmFsaWRhdGVyZXF1ZXN0'
    'Ym9keRJAChl2YWxpZGF0ZXJlcXVlc3RwYXJhbWV0ZXJzGNm/ot0BIAEoCFIZdmFsaWRhdGVyZX'
    'F1ZXN0cGFyYW1ldGVycw==');

@$core.Deprecated('Use requestValidatorsDescriptor instead')
const RequestValidators$json = {
  '1': 'RequestValidators',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.RequestValidator',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `RequestValidators`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestValidatorsDescriptor = $convert.base64Decode(
    'ChFSZXF1ZXN0VmFsaWRhdG9ycxI2CgVpdGVtcxiQ5+TTASADKAsyHC5hcGlnYXRld2F5LlJlcX'
    'Vlc3RWYWxpZGF0b3JSBWl0ZW1zEh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24=');

@$core.Deprecated('Use resourceDescriptor instead')
const Resource$json = {
  '1': 'Resource',
  '2': [
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'parentid', '3': 106050857, '4': 1, '5': 9, '10': 'parentid'},
    {'1': 'path', '3': 75975991, '4': 1, '5': 9, '10': 'path'},
    {'1': 'pathpart', '3': 487915984, '4': 1, '5': 9, '10': 'pathpart'},
    {
      '1': 'resourcemethods',
      '3': 307700458,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Resource.ResourcemethodsEntry',
      '10': 'resourcemethods'
    },
  ],
  '3': [Resource_ResourcemethodsEntry$json],
};

@$core.Deprecated('Use resourceDescriptor instead')
const Resource_ResourcemethodsEntry$json = {
  '1': 'ResourcemethodsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.apigateway.Method',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `Resource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDescriptor = $convert.base64Decode(
    'CghSZXNvdXJjZRISCgJpZBjh1eG5ASABKAlSAmlkEh0KCHBhcmVudGlkGKnqyDIgASgJUghwYX'
    'JlbnRpZBIVCgRwYXRoGLeanSQgASgJUgRwYXRoEh4KCHBhdGhwYXJ0GNCD1OgBIAEoCVIIcGF0'
    'aHBhcnQSVwoPcmVzb3VyY2VtZXRob2RzGOrF3JIBIAMoCzIpLmFwaWdhdGV3YXkuUmVzb3VyY2'
    'UuUmVzb3VyY2VtZXRob2RzRW50cnlSD3Jlc291cmNlbWV0aG9kcxpWChRSZXNvdXJjZW1ldGhv'
    'ZHNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIoCgV2YWx1ZRgCIAEoCzISLmFwaWdhdGV3YXkuTW'
    'V0aG9kUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use resourcesDescriptor instead')
const Resources$json = {
  '1': 'Resources',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Resource',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `Resources`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcesDescriptor = $convert.base64Decode(
    'CglSZXNvdXJjZXMSLgoFaXRlbXMYkOfk0wEgAygLMhQuYXBpZ2F0ZXdheS5SZXNvdXJjZVIFaX'
    'RlbXMSHgoIcG9zaXRpb24Yi5y9mgEgASgJUghwb3NpdGlvbg==');

@$core.Deprecated('Use restApiDescriptor instead')
const RestApi$json = {
  '1': 'RestApi',
  '2': [
    {
      '1': 'apikeysource',
      '3': 108531220,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ApiKeySourceType',
      '10': 'apikeysource'
    },
    {
      '1': 'apistatus',
      '3': 200568018,
      '4': 1,
      '5': 14,
      '6': '.apigateway.ApiStatus',
      '10': 'apistatus'
    },
    {
      '1': 'apistatusmessage',
      '3': 353995209,
      '4': 1,
      '5': 9,
      '10': 'apistatusmessage'
    },
    {
      '1': 'binarymediatypes',
      '3': 406416146,
      '4': 3,
      '5': 9,
      '10': 'binarymediatypes'
    },
    {'1': 'createddate', '3': 53061200, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'disableexecuteapiendpoint',
      '3': 148140696,
      '4': 1,
      '5': 8,
      '10': 'disableexecuteapiendpoint'
    },
    {
      '1': 'endpointaccessmode',
      '3': 356705630,
      '4': 1,
      '5': 14,
      '6': '.apigateway.EndpointAccessMode',
      '10': 'endpointaccessmode'
    },
    {
      '1': 'endpointconfiguration',
      '3': 487543735,
      '4': 1,
      '5': 11,
      '6': '.apigateway.EndpointConfiguration',
      '10': 'endpointconfiguration'
    },
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'minimumcompressionsize',
      '3': 254902719,
      '4': 1,
      '5': 5,
      '10': 'minimumcompressionsize'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'policy', '3': 247528064, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'rootresourceid',
      '3': 360157585,
      '4': 1,
      '5': 9,
      '10': 'rootresourceid'
    },
    {
      '1': 'securitypolicy',
      '3': 491792990,
      '4': 1,
      '5': 14,
      '6': '.apigateway.SecurityPolicy',
      '10': 'securitypolicy'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.RestApi.TagsEntry',
      '10': 'tags'
    },
    {'1': 'version', '3': 108113560, '4': 1, '5': 9, '10': 'version'},
    {'1': 'warnings', '3': 185617301, '4': 3, '5': 9, '10': 'warnings'},
  ],
  '3': [RestApi_TagsEntry$json],
};

@$core.Deprecated('Use restApiDescriptor instead')
const RestApi_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `RestApi`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restApiDescriptor = $convert.base64Decode(
    'CgdSZXN0QXBpEkMKDGFwaWtleXNvdXJjZRiUnOAzIAEoDjIcLmFwaWdhdGV3YXkuQXBpS2V5U2'
    '91cmNlVHlwZVIMYXBpa2V5c291cmNlEjYKCWFwaXN0YXR1cxjS2dFfIAEoDjIVLmFwaWdhdGV3'
    'YXkuQXBpU3RhdHVzUglhcGlzdGF0dXMSLgoQYXBpc3RhdHVzbWVzc2FnZRjJk+aoASABKAlSEG'
    'FwaXN0YXR1c21lc3NhZ2USLgoQYmluYXJ5bWVkaWF0eXBlcxiS1uXBASADKAlSEGJpbmFyeW1l'
    'ZGlhdHlwZXMSIwoLY3JlYXRlZGRhdGUY0MymGSABKAlSC2NyZWF0ZWRkYXRlEiQKC2Rlc2NyaX'
    'B0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SPwoZZGlzYWJsZWV4ZWN1dGVhcGllbmRwb2lu'
    'dBiY5dFGIAEoCFIZZGlzYWJsZWV4ZWN1dGVhcGllbmRwb2ludBJSChJlbmRwb2ludGFjY2Vzc2'
    '1vZGUY3sqLqgEgASgOMh4uYXBpZ2F0ZXdheS5FbmRwb2ludEFjY2Vzc01vZGVSEmVuZHBvaW50'
    'YWNjZXNzbW9kZRJbChVlbmRwb2ludGNvbmZpZ3VyYXRpb24Yt6e96AEgASgLMiEuYXBpZ2F0ZX'
    'dheS5FbmRwb2ludENvbmZpZ3VyYXRpb25SFWVuZHBvaW50Y29uZmlndXJhdGlvbhISCgJpZBjh'
    '1eG5ASABKAlSAmlkEjkKFm1pbmltdW1jb21wcmVzc2lvbnNpemUYv4PGeSABKAVSFm1pbmltdW'
    '1jb21wcmVzc2lvbnNpemUSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZRIZCgZwb2xpY3kYgPWDdiAB'
    'KAlSBnBvbGljeRIqCg5yb290cmVzb3VyY2VpZBiRo96rASABKAlSDnJvb3RyZXNvdXJjZWlkEk'
    'YKDnNlY3VyaXR5cG9saWN5GN7UwOoBIAEoDjIaLmFwaWdhdGV3YXkuU2VjdXJpdHlQb2xpY3lS'
    'DnNlY3VyaXR5cG9saWN5EjUKBHRhZ3MYodfboAEgAygLMh0uYXBpZ2F0ZXdheS5SZXN0QXBpLl'
    'RhZ3NFbnRyeVIEdGFncxIbCgd2ZXJzaW9uGJjdxjMgASgJUgd2ZXJzaW9uEh0KCHdhcm5pbmdz'
    'GJWXwVggAygJUgh3YXJuaW5ncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdm'
    'FsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use restApisDescriptor instead')
const RestApis$json = {
  '1': 'RestApis',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.RestApi',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `RestApis`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restApisDescriptor = $convert.base64Decode(
    'CghSZXN0QXBpcxItCgVpdGVtcxiQ5+TTASADKAsyEy5hcGlnYXRld2F5LlJlc3RBcGlSBWl0ZW'
    '1zEh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24=');

@$core.Deprecated('Use sdkConfigurationPropertyDescriptor instead')
const SdkConfigurationProperty$json = {
  '1': 'SdkConfigurationProperty',
  '2': [
    {'1': 'defaultvalue', '3': 403858624, '4': 1, '5': 9, '10': 'defaultvalue'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'friendlyname', '3': 39239102, '4': 1, '5': 9, '10': 'friendlyname'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'required', '3': 76318241, '4': 1, '5': 8, '10': 'required'},
  ],
};

/// Descriptor for `SdkConfigurationProperty`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sdkConfigurationPropertyDescriptor = $convert.base64Decode(
    'ChhTZGtDb25maWd1cmF0aW9uUHJvcGVydHkSJgoMZGVmYXVsdHZhbHVlGMDJycABIAEoCVIMZG'
    'VmYXVsdHZhbHVlEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SJQoMZnJp'
    'ZW5kbHluYW1lGL772hIgASgJUgxmcmllbmRseW5hbWUSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZR'
    'IdCghyZXF1aXJlZBihjLIkIAEoCFIIcmVxdWlyZWQ=');

@$core.Deprecated('Use sdkResponseDescriptor instead')
const SdkResponse$json = {
  '1': 'SdkResponse',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'contentdisposition',
      '3': 375146466,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {'1': 'contenttype', '3': 281764659, '4': 1, '5': 9, '10': 'contenttype'},
  ],
};

/// Descriptor for `SdkResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sdkResponseDescriptor = $convert.base64Decode(
    'CgtTZGtSZXNwb25zZRIcCgRib2R5GPbyqd0BIAEoDEIEiLUYAVIEYm9keRIyChJjb250ZW50ZG'
    'lzcG9zaXRpb24Y4o/xsgEgASgJUhJjb250ZW50ZGlzcG9zaXRpb24SJAoLY29udGVudHR5cGUY'
    's8athgEgASgJUgtjb250ZW50dHlwZQ==');

@$core.Deprecated('Use sdkTypeDescriptor instead')
const SdkType$json = {
  '1': 'SdkType',
  '2': [
    {
      '1': 'configurationproperties',
      '3': 113650241,
      '4': 3,
      '5': 11,
      '6': '.apigateway.SdkConfigurationProperty',
      '10': 'configurationproperties'
    },
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'friendlyname', '3': 39239102, '4': 1, '5': 9, '10': 'friendlyname'},
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `SdkType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sdkTypeDescriptor = $convert.base64Decode(
    'CgdTZGtUeXBlEmEKF2NvbmZpZ3VyYXRpb25wcm9wZXJ0aWVzGMHUmDYgAygLMiQuYXBpZ2F0ZX'
    'dheS5TZGtDb25maWd1cmF0aW9uUHJvcGVydHlSF2NvbmZpZ3VyYXRpb25wcm9wZXJ0aWVzEiQK'
    'C2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SJQoMZnJpZW5kbHluYW1lGL772h'
    'IgASgJUgxmcmllbmRseW5hbWUSEgoCaWQY4dXhuQEgASgJUgJpZA==');

@$core.Deprecated('Use sdkTypesDescriptor instead')
const SdkTypes$json = {
  '1': 'SdkTypes',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.SdkType',
      '10': 'items'
    },
  ],
};

/// Descriptor for `SdkTypes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sdkTypesDescriptor = $convert.base64Decode(
    'CghTZGtUeXBlcxItCgVpdGVtcxiQ5+TTASADKAsyEy5hcGlnYXRld2F5LlNka1R5cGVSBWl0ZW'
    '1z');

@$core.Deprecated('Use serviceUnavailableExceptionDescriptor instead')
const ServiceUnavailableException$json = {
  '1': 'ServiceUnavailableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'retryafterseconds',
      '3': 436039555,
      '4': 1,
      '5': 9,
      '10': 'retryafterseconds'
    },
  ],
};

/// Descriptor for `ServiceUnavailableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceUnavailableExceptionDescriptor =
    $convert.base64Decode(
        'ChtTZXJ2aWNlVW5hdmFpbGFibGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZRIwChFyZXRyeWFmdGVyc2Vjb25kcxiD3/XPASABKAlSEXJldHJ5YWZ0ZXJzZWNvbmRz');

@$core.Deprecated('Use stageDescriptor instead')
const Stage$json = {
  '1': 'Stage',
  '2': [
    {
      '1': 'accesslogsettings',
      '3': 76216807,
      '4': 1,
      '5': 11,
      '6': '.apigateway.AccessLogSettings',
      '10': 'accesslogsettings'
    },
    {
      '1': 'cacheclusterenabled',
      '3': 63967991,
      '4': 1,
      '5': 8,
      '10': 'cacheclusterenabled'
    },
    {
      '1': 'cacheclustersize',
      '3': 232189861,
      '4': 1,
      '5': 14,
      '6': '.apigateway.CacheClusterSize',
      '10': 'cacheclustersize'
    },
    {
      '1': 'cacheclusterstatus',
      '3': 385293784,
      '4': 1,
      '5': 14,
      '6': '.apigateway.CacheClusterStatus',
      '10': 'cacheclusterstatus'
    },
    {
      '1': 'canarysettings',
      '3': 285544261,
      '4': 1,
      '5': 11,
      '6': '.apigateway.CanarySettings',
      '10': 'canarysettings'
    },
    {
      '1': 'clientcertificateid',
      '3': 276222909,
      '4': 1,
      '5': 9,
      '10': 'clientcertificateid'
    },
    {'1': 'createddate', '3': 53061200, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'deploymentid', '3': 439369188, '4': 1, '5': 9, '10': 'deploymentid'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'documentationversion',
      '3': 167009804,
      '4': 1,
      '5': 9,
      '10': 'documentationversion'
    },
    {
      '1': 'lastupdateddate',
      '3': 448453361,
      '4': 1,
      '5': 9,
      '10': 'lastupdateddate'
    },
    {
      '1': 'methodsettings',
      '3': 30387838,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Stage.MethodsettingsEntry',
      '10': 'methodsettings'
    },
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Stage.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'tracingenabled',
      '3': 390995731,
      '4': 1,
      '5': 8,
      '10': 'tracingenabled'
    },
    {
      '1': 'variables',
      '3': 162226883,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Stage.VariablesEntry',
      '10': 'variables'
    },
    {'1': 'webaclarn', '3': 243701763, '4': 1, '5': 9, '10': 'webaclarn'},
  ],
  '3': [
    Stage_MethodsettingsEntry$json,
    Stage_TagsEntry$json,
    Stage_VariablesEntry$json
  ],
};

@$core.Deprecated('Use stageDescriptor instead')
const Stage_MethodsettingsEntry$json = {
  '1': 'MethodsettingsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.apigateway.MethodSetting',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use stageDescriptor instead')
const Stage_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use stageDescriptor instead')
const Stage_VariablesEntry$json = {
  '1': 'VariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Stage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stageDescriptor = $convert.base64Decode(
    'CgVTdGFnZRJOChFhY2Nlc3Nsb2dzZXR0aW5ncxjn86skIAEoCzIdLmFwaWdhdGV3YXkuQWNjZX'
    'NzTG9nU2V0dGluZ3NSEWFjY2Vzc2xvZ3NldHRpbmdzEjMKE2NhY2hlY2x1c3RlcmVuYWJsZWQY'
    '96XAHiABKAhSE2NhY2hlY2x1c3RlcmVuYWJsZWQSSwoQY2FjaGVjbHVzdGVyc2l6ZRil39tuIA'
    'EoDjIcLmFwaWdhdGV3YXkuQ2FjaGVDbHVzdGVyU2l6ZVIQY2FjaGVjbHVzdGVyc2l6ZRJSChJj'
    'YWNoZWNsdXN0ZXJzdGF0dXMY2LvctwEgASgOMh4uYXBpZ2F0ZXdheS5DYWNoZUNsdXN0ZXJTdG'
    'F0dXNSEmNhY2hlY2x1c3RlcnN0YXR1cxJGCg5jYW5hcnlzZXR0aW5ncxjFnpSIASABKAsyGi5h'
    'cGlnYXRld2F5LkNhbmFyeVNldHRpbmdzUg5jYW5hcnlzZXR0aW5ncxI0ChNjbGllbnRjZXJ0aW'
    'ZpY2F0ZWlkGL2n24MBIAEoCVITY2xpZW50Y2VydGlmaWNhdGVpZBIjCgtjcmVhdGVkZGF0ZRjQ'
    'zKYZIAEoCVILY3JlYXRlZGRhdGUSJgoMZGVwbG95bWVudGlkGOT7wNEBIAEoCVIMZGVwbG95bW'
    'VudGlkEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SNQoUZG9jdW1lbnRh'
    'dGlvbnZlcnNpb24YjLzRTyABKAlSFGRvY3VtZW50YXRpb252ZXJzaW9uEiwKD2xhc3R1cGRhdG'
    'VkZGF0ZRjxtevVASABKAlSD2xhc3R1cGRhdGVkZGF0ZRJQCg5tZXRob2RzZXR0aW5ncxj+3L4O'
    'IAMoCzIlLmFwaWdhdGV3YXkuU3RhZ2UuTWV0aG9kc2V0dGluZ3NFbnRyeVIObWV0aG9kc2V0dG'
    'luZ3MSHwoJc3RhZ2VuYW1lGI/cxwQgASgJUglzdGFnZW5hbWUSMwoEdGFncxih19ugASADKAsy'
    'Gy5hcGlnYXRld2F5LlN0YWdlLlRhZ3NFbnRyeVIEdGFncxIqCg50cmFjaW5nZW5hYmxlZBiTvr'
    'i6ASABKAhSDnRyYWNpbmdlbmFibGVkEkEKCXZhcmlhYmxlcxjDxa1NIAMoCzIgLmFwaWdhdGV3'
    'YXkuU3RhZ2UuVmFyaWFibGVzRW50cnlSCXZhcmlhYmxlcxIfCgl3ZWJhY2xhcm4Yg7CadCABKA'
    'lSCXdlYmFjbGFybhpcChNNZXRob2RzZXR0aW5nc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei8K'
    'BXZhbHVlGAIgASgLMhkuYXBpZ2F0ZXdheS5NZXRob2RTZXR0aW5nUgV2YWx1ZToCOAEaNwoJVG'
    'Fnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaPAoO'
    'VmFyaWFibGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOg'
    'I4AQ==');

@$core.Deprecated('Use stageKeyDescriptor instead')
const StageKey$json = {
  '1': 'StageKey',
  '2': [
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
};

/// Descriptor for `StageKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stageKeyDescriptor = $convert.base64Decode(
    'CghTdGFnZUtleRIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQSHwoJc3RhZ2VuYW'
    '1lGI/cxwQgASgJUglzdGFnZW5hbWU=');

@$core.Deprecated('Use stagesDescriptor instead')
const Stages$json = {
  '1': 'Stages',
  '2': [
    {
      '1': 'item',
      '3': 523776999,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Stage',
      '10': 'item'
    },
  ],
};

/// Descriptor for `Stages`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stagesDescriptor = $convert.base64Decode(
    'CgZTdGFnZXMSKQoEaXRlbRjn5+D5ASADKAsyES5hcGlnYXRld2F5LlN0YWdlUgRpdGVt');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TagResourceRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [TagResourceRequest_TagsEntry$json],
};

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSIwoLcmVzb3VyY2Vhcm4YzcyqICABKAlSC3Jlc291cmNlYX'
    'JuEkAKBHRhZ3MYodfboAEgAygLMiguYXBpZ2F0ZXdheS5UYWdSZXNvdXJjZVJlcXVlc3QuVGFn'
    'c0VudHJ5UgR0YWdzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIA'
    'EoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use tagsDescriptor instead')
const Tags$json = {
  '1': 'Tags',
  '2': [
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Tags.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [Tags_TagsEntry$json],
};

@$core.Deprecated('Use tagsDescriptor instead')
const Tags_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Tags`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagsDescriptor = $convert.base64Decode(
    'CgRUYWdzEjIKBHRhZ3MYodfboAEgAygLMhouYXBpZ2F0ZXdheS5UYWdzLlRhZ3NFbnRyeVIEdG'
    'Fncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVl'
    'OgI4AQ==');

@$core.Deprecated('Use templateDescriptor instead')
const Template$json = {
  '1': 'Template',
  '2': [
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Template`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List templateDescriptor =
    $convert.base64Decode('CghUZW1wbGF0ZRIXCgV2YWx1ZRjLp/sSIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use testInvokeAuthorizerRequestDescriptor instead')
const TestInvokeAuthorizerRequest$json = {
  '1': 'TestInvokeAuthorizerRequest',
  '2': [
    {
      '1': 'additionalcontext',
      '3': 471949152,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeAuthorizerRequest.AdditionalcontextEntry',
      '10': 'additionalcontext'
    },
    {'1': 'authorizerid', '3': 111773148, '4': 1, '5': 9, '10': 'authorizerid'},
    {'1': 'body', '3': 464157046, '4': 1, '5': 9, '10': 'body'},
    {
      '1': 'headers',
      '3': 375773674,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeAuthorizerRequest.HeadersEntry',
      '10': 'headers'
    },
    {
      '1': 'multivalueheaders',
      '3': 142421420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeAuthorizerRequest.MultivalueheadersEntry',
      '10': 'multivalueheaders'
    },
    {
      '1': 'pathwithquerystring',
      '3': 159950736,
      '4': 1,
      '5': 9,
      '10': 'pathwithquerystring'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'stagevariables',
      '3': 208524923,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeAuthorizerRequest.StagevariablesEntry',
      '10': 'stagevariables'
    },
  ],
  '3': [
    TestInvokeAuthorizerRequest_AdditionalcontextEntry$json,
    TestInvokeAuthorizerRequest_HeadersEntry$json,
    TestInvokeAuthorizerRequest_MultivalueheadersEntry$json,
    TestInvokeAuthorizerRequest_StagevariablesEntry$json
  ],
};

@$core.Deprecated('Use testInvokeAuthorizerRequestDescriptor instead')
const TestInvokeAuthorizerRequest_AdditionalcontextEntry$json = {
  '1': 'AdditionalcontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use testInvokeAuthorizerRequestDescriptor instead')
const TestInvokeAuthorizerRequest_HeadersEntry$json = {
  '1': 'HeadersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use testInvokeAuthorizerRequestDescriptor instead')
const TestInvokeAuthorizerRequest_MultivalueheadersEntry$json = {
  '1': 'MultivalueheadersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use testInvokeAuthorizerRequestDescriptor instead')
const TestInvokeAuthorizerRequest_StagevariablesEntry$json = {
  '1': 'StagevariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TestInvokeAuthorizerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testInvokeAuthorizerRequestDescriptor = $convert.base64Decode(
    'ChtUZXN0SW52b2tlQXV0aG9yaXplclJlcXVlc3QScAoRYWRkaXRpb25hbGNvbnRleHQY4L6F4Q'
    'EgAygLMj4uYXBpZ2F0ZXdheS5UZXN0SW52b2tlQXV0aG9yaXplclJlcXVlc3QuQWRkaXRpb25h'
    'bGNvbnRleHRFbnRyeVIRYWRkaXRpb25hbGNvbnRleHQSJQoMYXV0aG9yaXplcmlkGNyLpjUgAS'
    'gJUgxhdXRob3JpemVyaWQSFgoEYm9keRj28qndASABKAlSBGJvZHkSUgoHaGVhZGVycxjqs5ez'
    'ASADKAsyNC5hcGlnYXRld2F5LlRlc3RJbnZva2VBdXRob3JpemVyUmVxdWVzdC5IZWFkZXJzRW'
    '50cnlSB2hlYWRlcnMSbwoRbXVsdGl2YWx1ZWhlYWRlcnMYrNv0QyADKAsyPi5hcGlnYXRld2F5'
    'LlRlc3RJbnZva2VBdXRob3JpemVyUmVxdWVzdC5NdWx0aXZhbHVlaGVhZGVyc0VudHJ5UhFtdW'
    'x0aXZhbHVlaGVhZGVycxIzChNwYXRod2l0aHF1ZXJ5c3RyaW5nGJDPokwgASgJUhNwYXRod2l0'
    'aHF1ZXJ5c3RyaW5nEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZBJmCg5zdGFnZX'
    'ZhcmlhYmxlcxj7rLdjIAMoCzI7LmFwaWdhdGV3YXkuVGVzdEludm9rZUF1dGhvcml6ZXJSZXF1'
    'ZXN0LlN0YWdldmFyaWFibGVzRW50cnlSDnN0YWdldmFyaWFibGVzGkQKFkFkZGl0aW9uYWxjb2'
    '50ZXh0RW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARo6'
    'CgxIZWFkZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOg'
    'I4ARpEChZNdWx0aXZhbHVlaGVhZGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVl'
    'GAIgASgJUgV2YWx1ZToCOAEaQQoTU3RhZ2V2YXJpYWJsZXNFbnRyeRIQCgNrZXkYASABKAlSA2'
    'tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use testInvokeAuthorizerResponseDescriptor instead')
const TestInvokeAuthorizerResponse$json = {
  '1': 'TestInvokeAuthorizerResponse',
  '2': [
    {
      '1': 'authorization',
      '3': 288774079,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeAuthorizerResponse.AuthorizationEntry',
      '10': 'authorization'
    },
    {
      '1': 'claims',
      '3': 479124501,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeAuthorizerResponse.ClaimsEntry',
      '10': 'claims'
    },
    {'1': 'clientstatus', '3': 35642913, '4': 1, '5': 5, '10': 'clientstatus'},
    {'1': 'latency', '3': 318473050, '4': 1, '5': 3, '10': 'latency'},
    {'1': 'log', '3': 525422930, '4': 1, '5': 9, '10': 'log'},
    {'1': 'policy', '3': 247528064, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'principalid', '3': 350710285, '4': 1, '5': 9, '10': 'principalid'},
  ],
  '3': [
    TestInvokeAuthorizerResponse_AuthorizationEntry$json,
    TestInvokeAuthorizerResponse_ClaimsEntry$json
  ],
};

@$core.Deprecated('Use testInvokeAuthorizerResponseDescriptor instead')
const TestInvokeAuthorizerResponse_AuthorizationEntry$json = {
  '1': 'AuthorizationEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use testInvokeAuthorizerResponseDescriptor instead')
const TestInvokeAuthorizerResponse_ClaimsEntry$json = {
  '1': 'ClaimsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TestInvokeAuthorizerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testInvokeAuthorizerResponseDescriptor = $convert.base64Decode(
    'ChxUZXN0SW52b2tlQXV0aG9yaXplclJlc3BvbnNlEmUKDWF1dGhvcml6YXRpb24Yv6/ZiQEgAy'
    'gLMjsuYXBpZ2F0ZXdheS5UZXN0SW52b2tlQXV0aG9yaXplclJlc3BvbnNlLkF1dGhvcml6YXRp'
    'b25FbnRyeVINYXV0aG9yaXphdGlvbhJQCgZjbGFpbXMYlbi75AEgAygLMjQuYXBpZ2F0ZXdheS'
    '5UZXN0SW52b2tlQXV0aG9yaXplclJlc3BvbnNlLkNsYWltc0VudHJ5UgZjbGFpbXMSJQoMY2xp'
    'ZW50c3RhdHVzGKG8/xAgASgFUgxjbGllbnRzdGF0dXMSHAoHbGF0ZW5jeRjahu6XASABKANSB2'
    'xhdGVuY3kSFAoDbG9nGNKixfoBIAEoCVIDbG9nEhkKBnBvbGljeRiA9YN2IAEoCVIGcG9saWN5'
    'EiQKC3ByaW5jaXBhbGlkGI3UnacBIAEoCVILcHJpbmNpcGFsaWQaQAoSQXV0aG9yaXphdGlvbk'
    'VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaOQoLQ2xh'
    'aW1zRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use testInvokeMethodRequestDescriptor instead')
const TestInvokeMethodRequest$json = {
  '1': 'TestInvokeMethodRequest',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 9, '10': 'body'},
    {
      '1': 'clientcertificateid',
      '3': 276222909,
      '4': 1,
      '5': 9,
      '10': 'clientcertificateid'
    },
    {
      '1': 'headers',
      '3': 375773674,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeMethodRequest.HeadersEntry',
      '10': 'headers'
    },
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'multivalueheaders',
      '3': 142421420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeMethodRequest.MultivalueheadersEntry',
      '10': 'multivalueheaders'
    },
    {
      '1': 'pathwithquerystring',
      '3': 159950736,
      '4': 1,
      '5': 9,
      '10': 'pathwithquerystring'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {
      '1': 'stagevariables',
      '3': 208524923,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeMethodRequest.StagevariablesEntry',
      '10': 'stagevariables'
    },
  ],
  '3': [
    TestInvokeMethodRequest_HeadersEntry$json,
    TestInvokeMethodRequest_MultivalueheadersEntry$json,
    TestInvokeMethodRequest_StagevariablesEntry$json
  ],
};

@$core.Deprecated('Use testInvokeMethodRequestDescriptor instead')
const TestInvokeMethodRequest_HeadersEntry$json = {
  '1': 'HeadersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use testInvokeMethodRequestDescriptor instead')
const TestInvokeMethodRequest_MultivalueheadersEntry$json = {
  '1': 'MultivalueheadersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use testInvokeMethodRequestDescriptor instead')
const TestInvokeMethodRequest_StagevariablesEntry$json = {
  '1': 'StagevariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TestInvokeMethodRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testInvokeMethodRequestDescriptor = $convert.base64Decode(
    'ChdUZXN0SW52b2tlTWV0aG9kUmVxdWVzdBIWCgRib2R5GPbyqd0BIAEoCVIEYm9keRI0ChNjbG'
    'llbnRjZXJ0aWZpY2F0ZWlkGL2n24MBIAEoCVITY2xpZW50Y2VydGlmaWNhdGVpZBJOCgdoZWFk'
    'ZXJzGOqzl7MBIAMoCzIwLmFwaWdhdGV3YXkuVGVzdEludm9rZU1ldGhvZFJlcXVlc3QuSGVhZG'
    'Vyc0VudHJ5UgdoZWFkZXJzEiEKCmh0dHBtZXRob2QY8fP7NiABKAlSCmh0dHBtZXRob2QSawoR'
    'bXVsdGl2YWx1ZWhlYWRlcnMYrNv0QyADKAsyOi5hcGlnYXRld2F5LlRlc3RJbnZva2VNZXRob2'
    'RSZXF1ZXN0Lk11bHRpdmFsdWVoZWFkZXJzRW50cnlSEW11bHRpdmFsdWVoZWFkZXJzEjMKE3Bh'
    'dGh3aXRocXVlcnlzdHJpbmcYkM+iTCABKAlSE3BhdGh3aXRocXVlcnlzdHJpbmcSIgoKcmVzb3'
    'VyY2VpZBixvYmYASABKAlSCnJlc291cmNlaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVz'
    'dGFwaWlkEmIKDnN0YWdldmFyaWFibGVzGPust2MgAygLMjcuYXBpZ2F0ZXdheS5UZXN0SW52b2'
    'tlTWV0aG9kUmVxdWVzdC5TdGFnZXZhcmlhYmxlc0VudHJ5Ug5zdGFnZXZhcmlhYmxlcxo6CgxI'
    'ZWFkZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AR'
    'pEChZNdWx0aXZhbHVlaGVhZGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIg'
    'ASgJUgV2YWx1ZToCOAEaQQoTU3RhZ2V2YXJpYWJsZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleR'
    'IUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use testInvokeMethodResponseDescriptor instead')
const TestInvokeMethodResponse$json = {
  '1': 'TestInvokeMethodResponse',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 9, '10': 'body'},
    {
      '1': 'headers',
      '3': 375773674,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeMethodResponse.HeadersEntry',
      '10': 'headers'
    },
    {'1': 'latency', '3': 318473050, '4': 1, '5': 3, '10': 'latency'},
    {'1': 'log', '3': 525422930, '4': 1, '5': 9, '10': 'log'},
    {
      '1': 'multivalueheaders',
      '3': 142421420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.TestInvokeMethodResponse.MultivalueheadersEntry',
      '10': 'multivalueheaders'
    },
    {'1': 'status', '3': 441153520, '4': 1, '5': 5, '10': 'status'},
  ],
  '3': [
    TestInvokeMethodResponse_HeadersEntry$json,
    TestInvokeMethodResponse_MultivalueheadersEntry$json
  ],
};

@$core.Deprecated('Use testInvokeMethodResponseDescriptor instead')
const TestInvokeMethodResponse_HeadersEntry$json = {
  '1': 'HeadersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use testInvokeMethodResponseDescriptor instead')
const TestInvokeMethodResponse_MultivalueheadersEntry$json = {
  '1': 'MultivalueheadersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TestInvokeMethodResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testInvokeMethodResponseDescriptor = $convert.base64Decode(
    'ChhUZXN0SW52b2tlTWV0aG9kUmVzcG9uc2USFgoEYm9keRj28qndASABKAlSBGJvZHkSTwoHaG'
    'VhZGVycxjqs5ezASADKAsyMS5hcGlnYXRld2F5LlRlc3RJbnZva2VNZXRob2RSZXNwb25zZS5I'
    'ZWFkZXJzRW50cnlSB2hlYWRlcnMSHAoHbGF0ZW5jeRjahu6XASABKANSB2xhdGVuY3kSFAoDbG'
    '9nGNKixfoBIAEoCVIDbG9nEmwKEW11bHRpdmFsdWVoZWFkZXJzGKzb9EMgAygLMjsuYXBpZ2F0'
    'ZXdheS5UZXN0SW52b2tlTWV0aG9kUmVzcG9uc2UuTXVsdGl2YWx1ZWhlYWRlcnNFbnRyeVIRbX'
    'VsdGl2YWx1ZWhlYWRlcnMSGgoGc3RhdHVzGPDvrdIBIAEoBVIGc3RhdHVzGjoKDEhlYWRlcnNF'
    'bnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGkQKFk11bH'
    'RpdmFsdWVoZWFkZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZh'
    'bHVlOgI4AQ==');

@$core.Deprecated('Use throttleSettingsDescriptor instead')
const ThrottleSettings$json = {
  '1': 'ThrottleSettings',
  '2': [
    {'1': 'burstlimit', '3': 37855041, '4': 1, '5': 5, '10': 'burstlimit'},
    {'1': 'ratelimit', '3': 505789539, '4': 1, '5': 1, '10': 'ratelimit'},
  ],
};

/// Descriptor for `ThrottleSettings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttleSettingsDescriptor = $convert.base64Decode(
    'ChBUaHJvdHRsZVNldHRpbmdzEiEKCmJ1cnN0bGltaXQYwb6GEiABKAVSCmJ1cnN0bGltaXQSIA'
    'oJcmF0ZWxpbWl0GOP4lvEBIAEoAVIJcmF0ZWxpbWl0');

@$core.Deprecated('Use tlsConfigDescriptor instead')
const TlsConfig$json = {
  '1': 'TlsConfig',
  '2': [
    {
      '1': 'insecureskipverification',
      '3': 78100228,
      '4': 1,
      '5': 8,
      '10': 'insecureskipverification'
    },
  ],
};

/// Descriptor for `TlsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tlsConfigDescriptor = $convert.base64Decode(
    'CglUbHNDb25maWcSPQoYaW5zZWN1cmVza2lwdmVyaWZpY2F0aW9uGITuniUgASgIUhhpbnNlY3'
    'VyZXNraXB2ZXJpZmljYXRpb24=');

@$core.Deprecated('Use tooManyRequestsExceptionDescriptor instead')
const TooManyRequestsException$json = {
  '1': 'TooManyRequestsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'retryafterseconds',
      '3': 436039555,
      '4': 1,
      '5': 9,
      '10': 'retryafterseconds'
    },
  ],
};

/// Descriptor for `TooManyRequestsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyRequestsExceptionDescriptor =
    $convert.base64Decode(
        'ChhUb29NYW55UmVxdWVzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZR'
        'IwChFyZXRyeWFmdGVyc2Vjb25kcxiD3/XPASABKAlSEXJldHJ5YWZ0ZXJzZWNvbmRz');

@$core.Deprecated('Use unauthorizedExceptionDescriptor instead')
const UnauthorizedException$json = {
  '1': 'UnauthorizedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnauthorizedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unauthorizedExceptionDescriptor = $convert.base64Decode(
    'ChVVbmF1dGhvcml6ZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 78811036, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIjCgtyZXNvdXJjZWFybhjNzKogIAEoCVILcmVzb3VyY2'
    'Vhcm4SGwoHdGFna2V5cxicn8olIAMoCVIHdGFna2V5cw==');

@$core.Deprecated('Use updateAccountRequestDescriptor instead')
const UpdateAccountRequest$json = {
  '1': 'UpdateAccountRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
  ],
};

/// Descriptor for `UpdateAccountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAccountRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVBY2NvdW50UmVxdWVzdBJHCg9wYXRjaG9wZXJhdGlvbnMYrPySYCADKAsyGi5hcG'
    'lnYXRld2F5LlBhdGNoT3BlcmF0aW9uUg9wYXRjaG9wZXJhdGlvbnM=');

@$core.Deprecated('Use updateApiKeyRequestDescriptor instead')
const UpdateApiKeyRequest$json = {
  '1': 'UpdateApiKeyRequest',
  '2': [
    {'1': 'apikey', '3': 490490655, '4': 1, '5': 9, '10': 'apikey'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
  ],
};

/// Descriptor for `UpdateApiKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateApiKeyRequestDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVBcGlLZXlSZXF1ZXN0EhoKBmFwaWtleRiflvHpASABKAlSBmFwaWtleRJHCg9wYX'
    'RjaG9wZXJhdGlvbnMYrPySYCADKAsyGi5hcGlnYXRld2F5LlBhdGNoT3BlcmF0aW9uUg9wYXRj'
    'aG9wZXJhdGlvbnM=');

@$core.Deprecated('Use updateAuthorizerRequestDescriptor instead')
const UpdateAuthorizerRequest$json = {
  '1': 'UpdateAuthorizerRequest',
  '2': [
    {'1': 'authorizerid', '3': 111773148, '4': 1, '5': 9, '10': 'authorizerid'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateAuthorizerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAuthorizerRequestDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVBdXRob3JpemVyUmVxdWVzdBIlCgxhdXRob3JpemVyaWQY3IumNSABKAlSDGF1dG'
    'hvcml6ZXJpZBJHCg9wYXRjaG9wZXJhdGlvbnMYrPySYCADKAsyGi5hcGlnYXRld2F5LlBhdGNo'
    'T3BlcmF0aW9uUg9wYXRjaG9wZXJhdGlvbnMSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdG'
    'FwaWlk');

@$core.Deprecated('Use updateBasePathMappingRequestDescriptor instead')
const UpdateBasePathMappingRequest$json = {
  '1': 'UpdateBasePathMappingRequest',
  '2': [
    {'1': 'basepath', '3': 267528880, '4': 1, '5': 9, '10': 'basepath'},
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
  ],
};

/// Descriptor for `UpdateBasePathMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateBasePathMappingRequestDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVCYXNlUGF0aE1hcHBpbmdSZXF1ZXN0Eh0KCGJhc2VwYXRoGLDVyH8gASgJUghiYX'
    'NlcGF0aBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG9tYWlubmFtZRImCgxkb21haW5uYW1l'
    'aWQYqPycjgEgASgJUgxkb21haW5uYW1laWQSRwoPcGF0Y2hvcGVyYXRpb25zGKz8kmAgAygLMh'
    'ouYXBpZ2F0ZXdheS5QYXRjaE9wZXJhdGlvblIPcGF0Y2hvcGVyYXRpb25z');

@$core.Deprecated('Use updateClientCertificateRequestDescriptor instead')
const UpdateClientCertificateRequest$json = {
  '1': 'UpdateClientCertificateRequest',
  '2': [
    {
      '1': 'clientcertificateid',
      '3': 276222909,
      '4': 1,
      '5': 9,
      '10': 'clientcertificateid'
    },
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
  ],
};

/// Descriptor for `UpdateClientCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateClientCertificateRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVDbGllbnRDZXJ0aWZpY2F0ZVJlcXVlc3QSNAoTY2xpZW50Y2VydGlmaWNhdGVpZB'
        'i9p9uDASABKAlSE2NsaWVudGNlcnRpZmljYXRlaWQSRwoPcGF0Y2hvcGVyYXRpb25zGKz8kmAg'
        'AygLMhouYXBpZ2F0ZXdheS5QYXRjaE9wZXJhdGlvblIPcGF0Y2hvcGVyYXRpb25z');

@$core.Deprecated('Use updateDeploymentRequestDescriptor instead')
const UpdateDeploymentRequest$json = {
  '1': 'UpdateDeploymentRequest',
  '2': [
    {'1': 'deploymentid', '3': 439369188, '4': 1, '5': 9, '10': 'deploymentid'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateDeploymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDeploymentRequestDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVEZXBsb3ltZW50UmVxdWVzdBImCgxkZXBsb3ltZW50aWQY5PvA0QEgASgJUgxkZX'
    'Bsb3ltZW50aWQSRwoPcGF0Y2hvcGVyYXRpb25zGKz8kmAgAygLMhouYXBpZ2F0ZXdheS5QYXRj'
    'aE9wZXJhdGlvblIPcGF0Y2hvcGVyYXRpb25zEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3'
    'RhcGlpZA==');

@$core.Deprecated('Use updateDocumentationPartRequestDescriptor instead')
const UpdateDocumentationPartRequest$json = {
  '1': 'UpdateDocumentationPartRequest',
  '2': [
    {
      '1': 'documentationpartid',
      '3': 286552774,
      '4': 1,
      '5': 9,
      '10': 'documentationpartid'
    },
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateDocumentationPartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentationPartRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVEb2N1bWVudGF0aW9uUGFydFJlcXVlc3QSNAoTZG9jdW1lbnRhdGlvbnBhcnRpZB'
        'jG5dGIASABKAlSE2RvY3VtZW50YXRpb25wYXJ0aWQSRwoPcGF0Y2hvcGVyYXRpb25zGKz8kmAg'
        'AygLMhouYXBpZ2F0ZXdheS5QYXRjaE9wZXJhdGlvblIPcGF0Y2hvcGVyYXRpb25zEiAKCXJlc3'
        'RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZA==');

@$core.Deprecated('Use updateDocumentationVersionRequestDescriptor instead')
const UpdateDocumentationVersionRequest$json = {
  '1': 'UpdateDocumentationVersionRequest',
  '2': [
    {
      '1': 'documentationversion',
      '3': 167009804,
      '4': 1,
      '5': 9,
      '10': 'documentationversion'
    },
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateDocumentationVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentationVersionRequestDescriptor =
    $convert.base64Decode(
        'CiFVcGRhdGVEb2N1bWVudGF0aW9uVmVyc2lvblJlcXVlc3QSNQoUZG9jdW1lbnRhdGlvbnZlcn'
        'Npb24YjLzRTyABKAlSFGRvY3VtZW50YXRpb252ZXJzaW9uEkcKD3BhdGNob3BlcmF0aW9ucxis'
        '/JJgIAMoCzIaLmFwaWdhdGV3YXkuUGF0Y2hPcGVyYXRpb25SD3BhdGNob3BlcmF0aW9ucxIgCg'
        'lyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQ=');

@$core.Deprecated('Use updateDomainNameRequestDescriptor instead')
const UpdateDomainNameRequest$json = {
  '1': 'UpdateDomainNameRequest',
  '2': [
    {'1': 'domainname', '3': 390326667, '4': 1, '5': 9, '10': 'domainname'},
    {'1': 'domainnameid', '3': 298270248, '4': 1, '5': 9, '10': 'domainnameid'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
  ],
};

/// Descriptor for `UpdateDomainNameRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDomainNameRequestDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVEb21haW5OYW1lUmVxdWVzdBIiCgpkb21haW5uYW1lGIvTj7oBIAEoCVIKZG9tYW'
    'lubmFtZRImCgxkb21haW5uYW1laWQYqPycjgEgASgJUgxkb21haW5uYW1laWQSRwoPcGF0Y2hv'
    'cGVyYXRpb25zGKz8kmAgAygLMhouYXBpZ2F0ZXdheS5QYXRjaE9wZXJhdGlvblIPcGF0Y2hvcG'
    'VyYXRpb25z');

@$core.Deprecated('Use updateGatewayResponseRequestDescriptor instead')
const UpdateGatewayResponseRequest$json = {
  '1': 'UpdateGatewayResponseRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {
      '1': 'responsetype',
      '3': 377935935,
      '4': 1,
      '5': 14,
      '6': '.apigateway.GatewayResponseType',
      '10': 'responsetype'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateGatewayResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGatewayResponseRequestDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVHYXRld2F5UmVzcG9uc2VSZXF1ZXN0EkcKD3BhdGNob3BlcmF0aW9ucxis/JJgIA'
    'MoCzIaLmFwaWdhdGV3YXkuUGF0Y2hPcGVyYXRpb25SD3BhdGNob3BlcmF0aW9ucxJHCgxyZXNw'
    'b25zZXR5cGUYv7CbtAEgASgOMh8uYXBpZ2F0ZXdheS5HYXRld2F5UmVzcG9uc2VUeXBlUgxyZX'
    'Nwb25zZXR5cGUSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use updateIntegrationRequestDescriptor instead')
const UpdateIntegrationRequest$json = {
  '1': 'UpdateIntegrationRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateIntegrationRequestDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVJbnRlZ3JhdGlvblJlcXVlc3QSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaHR0cG'
    '1ldGhvZBJHCg9wYXRjaG9wZXJhdGlvbnMYrPySYCADKAsyGi5hcGlnYXRld2F5LlBhdGNoT3Bl'
    'cmF0aW9uUg9wYXRjaG9wZXJhdGlvbnMSIgoKcmVzb3VyY2VpZBixvYmYASABKAlSCnJlc291cm'
    'NlaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use updateIntegrationResponseRequestDescriptor instead')
const UpdateIntegrationResponseRequest$json = {
  '1': 'UpdateIntegrationResponseRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
};

/// Descriptor for `UpdateIntegrationResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateIntegrationResponseRequestDescriptor = $convert.base64Decode(
    'CiBVcGRhdGVJbnRlZ3JhdGlvblJlc3BvbnNlUmVxdWVzdBIhCgpodHRwbWV0aG9kGPHz+zYgAS'
    'gJUgpodHRwbWV0aG9kEkcKD3BhdGNob3BlcmF0aW9ucxis/JJgIAMoCzIaLmFwaWdhdGV3YXku'
    'UGF0Y2hPcGVyYXRpb25SD3BhdGNob3BlcmF0aW9ucxIiCgpyZXNvdXJjZWlkGLG9iZgBIAEoCV'
    'IKcmVzb3VyY2VpZBIgCglyZXN0YXBpaWQYmaSBtwEgASgJUglyZXN0YXBpaWQSIgoKc3RhdHVz'
    'Y29kZRifgd+OASABKAlSCnN0YXR1c2NvZGU=');

@$core.Deprecated('Use updateMethodRequestDescriptor instead')
const UpdateMethodRequest$json = {
  '1': 'UpdateMethodRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateMethodRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMethodRequestDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVNZXRob2RSZXF1ZXN0EiEKCmh0dHBtZXRob2QY8fP7NiABKAlSCmh0dHBtZXRob2'
    'QSRwoPcGF0Y2hvcGVyYXRpb25zGKz8kmAgAygLMhouYXBpZ2F0ZXdheS5QYXRjaE9wZXJhdGlv'
    'blIPcGF0Y2hvcGVyYXRpb25zEiIKCnJlc291cmNlaWQYsb2JmAEgASgJUgpyZXNvdXJjZWlkEi'
    'AKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZA==');

@$core.Deprecated('Use updateMethodResponseRequestDescriptor instead')
const UpdateMethodResponseRequest$json = {
  '1': 'UpdateMethodResponseRequest',
  '2': [
    {'1': 'httpmethod', '3': 115276273, '4': 1, '5': 9, '10': 'httpmethod'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
  ],
};

/// Descriptor for `UpdateMethodResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMethodResponseRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVNZXRob2RSZXNwb25zZVJlcXVlc3QSIQoKaHR0cG1ldGhvZBjx8/s2IAEoCVIKaH'
    'R0cG1ldGhvZBJHCg9wYXRjaG9wZXJhdGlvbnMYrPySYCADKAsyGi5hcGlnYXRld2F5LlBhdGNo'
    'T3BlcmF0aW9uUg9wYXRjaG9wZXJhdGlvbnMSIgoKcmVzb3VyY2VpZBixvYmYASABKAlSCnJlc2'
    '91cmNlaWQSIAoJcmVzdGFwaWlkGJmkgbcBIAEoCVIJcmVzdGFwaWlkEiIKCnN0YXR1c2NvZGUY'
    'n4HfjgEgASgJUgpzdGF0dXNjb2Rl');

@$core.Deprecated('Use updateModelRequestDescriptor instead')
const UpdateModelRequest$json = {
  '1': 'UpdateModelRequest',
  '2': [
    {'1': 'modelname', '3': 176835330, '4': 1, '5': 9, '10': 'modelname'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateModelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateModelRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVNb2RlbFJlcXVlc3QSHwoJbW9kZWxuYW1lGIKWqVQgASgJUgltb2RlbG5hbWUSRw'
    'oPcGF0Y2hvcGVyYXRpb25zGKz8kmAgAygLMhouYXBpZ2F0ZXdheS5QYXRjaE9wZXJhdGlvblIP'
    'cGF0Y2hvcGVyYXRpb25zEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZA==');

@$core.Deprecated('Use updateRequestValidatorRequestDescriptor instead')
const UpdateRequestValidatorRequest$json = {
  '1': 'UpdateRequestValidatorRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {
      '1': 'requestvalidatorid',
      '3': 517546134,
      '4': 1,
      '5': 9,
      '10': 'requestvalidatorid'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateRequestValidatorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRequestValidatorRequestDescriptor = $convert.base64Decode(
    'Ch1VcGRhdGVSZXF1ZXN0VmFsaWRhdG9yUmVxdWVzdBJHCg9wYXRjaG9wZXJhdGlvbnMYrPySYC'
    'ADKAsyGi5hcGlnYXRld2F5LlBhdGNoT3BlcmF0aW9uUg9wYXRjaG9wZXJhdGlvbnMSMgoScmVx'
    'dWVzdHZhbGlkYXRvcmlkGJbB5PYBIAEoCVIScmVxdWVzdHZhbGlkYXRvcmlkEiAKCXJlc3RhcG'
    'lpZBiZpIG3ASABKAlSCXJlc3RhcGlpZA==');

@$core.Deprecated('Use updateResourceRequestDescriptor instead')
const UpdateResourceRequest$json = {
  '1': 'UpdateResourceRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'resourceid', '3': 318922417, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResourceRequestDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVSZXNvdXJjZVJlcXVlc3QSRwoPcGF0Y2hvcGVyYXRpb25zGKz8kmAgAygLMhouYX'
    'BpZ2F0ZXdheS5QYXRjaE9wZXJhdGlvblIPcGF0Y2hvcGVyYXRpb25zEiIKCnJlc291cmNlaWQY'
    'sb2JmAEgASgJUgpyZXNvdXJjZWlkEiAKCXJlc3RhcGlpZBiZpIG3ASABKAlSCXJlc3RhcGlpZA'
    '==');

@$core.Deprecated('Use updateRestApiRequestDescriptor instead')
const UpdateRestApiRequest$json = {
  '1': 'UpdateRestApiRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
  ],
};

/// Descriptor for `UpdateRestApiRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateRestApiRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVSZXN0QXBpUmVxdWVzdBJHCg9wYXRjaG9wZXJhdGlvbnMYrPySYCADKAsyGi5hcG'
    'lnYXRld2F5LlBhdGNoT3BlcmF0aW9uUg9wYXRjaG9wZXJhdGlvbnMSIAoJcmVzdGFwaWlkGJmk'
    'gbcBIAEoCVIJcmVzdGFwaWlk');

@$core.Deprecated('Use updateStageRequestDescriptor instead')
const UpdateStageRequest$json = {
  '1': 'UpdateStageRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'restapiid', '3': 383799833, '4': 1, '5': 9, '10': 'restapiid'},
    {'1': 'stagename', '3': 9563663, '4': 1, '5': 9, '10': 'stagename'},
  ],
};

/// Descriptor for `UpdateStageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStageRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVTdGFnZVJlcXVlc3QSRwoPcGF0Y2hvcGVyYXRpb25zGKz8kmAgAygLMhouYXBpZ2'
    'F0ZXdheS5QYXRjaE9wZXJhdGlvblIPcGF0Y2hvcGVyYXRpb25zEiAKCXJlc3RhcGlpZBiZpIG3'
    'ASABKAlSCXJlc3RhcGlpZBIfCglzdGFnZW5hbWUYj9zHBCABKAlSCXN0YWdlbmFtZQ==');

@$core.Deprecated('Use updateUsagePlanRequestDescriptor instead')
const UpdateUsagePlanRequest$json = {
  '1': 'UpdateUsagePlanRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `UpdateUsagePlanRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUsagePlanRequestDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVVc2FnZVBsYW5SZXF1ZXN0EkcKD3BhdGNob3BlcmF0aW9ucxis/JJgIAMoCzIaLm'
    'FwaWdhdGV3YXkuUGF0Y2hPcGVyYXRpb25SD3BhdGNob3BlcmF0aW9ucxIkCgt1c2FnZXBsYW5p'
    'ZBjX8OXyASABKAlSC3VzYWdlcGxhbmlk');

@$core.Deprecated('Use updateUsageRequestDescriptor instead')
const UpdateUsageRequest$json = {
  '1': 'UpdateUsageRequest',
  '2': [
    {'1': 'keyid', '3': 479913282, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
};

/// Descriptor for `UpdateUsageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUsageRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVVc2FnZVJlcXVlc3QSGAoFa2V5aWQYwsrr5AEgASgJUgVrZXlpZBJHCg9wYXRjaG'
    '9wZXJhdGlvbnMYrPySYCADKAsyGi5hcGlnYXRld2F5LlBhdGNoT3BlcmF0aW9uUg9wYXRjaG9w'
    'ZXJhdGlvbnMSJAoLdXNhZ2VwbGFuaWQY1/Dl8gEgASgJUgt1c2FnZXBsYW5pZA==');

@$core.Deprecated('Use updateVpcLinkRequestDescriptor instead')
const UpdateVpcLinkRequest$json = {
  '1': 'UpdateVpcLinkRequest',
  '2': [
    {
      '1': 'patchoperations',
      '3': 201637420,
      '4': 3,
      '5': 11,
      '6': '.apigateway.PatchOperation',
      '10': 'patchoperations'
    },
    {'1': 'vpclinkid', '3': 27515846, '4': 1, '5': 9, '10': 'vpclinkid'},
  ],
};

/// Descriptor for `UpdateVpcLinkRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateVpcLinkRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVWcGNMaW5rUmVxdWVzdBJHCg9wYXRjaG9wZXJhdGlvbnMYrPySYCADKAsyGi5hcG'
    'lnYXRld2F5LlBhdGNoT3BlcmF0aW9uUg9wYXRjaG9wZXJhdGlvbnMSHwoJdnBjbGlua2lkGMa3'
    'jw0gASgJUgl2cGNsaW5raWQ=');

@$core.Deprecated('Use usageDescriptor instead')
const Usage$json = {
  '1': 'Usage',
  '2': [
    {'1': 'enddate', '3': 384831215, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.Usage.ItemsEntry',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'usageplanid', '3': 509179991, '4': 1, '5': 9, '10': 'usageplanid'},
  ],
  '3': [Usage_ItemsEntry$json],
};

@$core.Deprecated('Use usageDescriptor instead')
const Usage_ItemsEntry$json = {
  '1': 'ItemsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Usage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List usageDescriptor = $convert.base64Decode(
    'CgVVc2FnZRIcCgdlbmRkYXRlGO+dwLcBIAEoCVIHZW5kZGF0ZRI2CgVpdGVtcxiQ5+TTASADKA'
    'syHC5hcGlnYXRld2F5LlVzYWdlLkl0ZW1zRW50cnlSBWl0ZW1zEh4KCHBvc2l0aW9uGIucvZoB'
    'IAEoCVIIcG9zaXRpb24SIAoJc3RhcnRkYXRlGJyO/K0BIAEoCVIJc3RhcnRkYXRlEiQKC3VzYW'
    'dlcGxhbmlkGNfw5fIBIAEoCVILdXNhZ2VwbGFuaWQaOAoKSXRlbXNFbnRyeRIQCgNrZXkYASAB'
    'KAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use usagePlanDescriptor instead')
const UsagePlan$json = {
  '1': 'UsagePlan',
  '2': [
    {
      '1': 'apistages',
      '3': 64558449,
      '4': 3,
      '5': 11,
      '6': '.apigateway.ApiStage',
      '10': 'apistages'
    },
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'productcode', '3': 533381226, '4': 1, '5': 9, '10': 'productcode'},
    {
      '1': 'quota',
      '3': 243824012,
      '4': 1,
      '5': 11,
      '6': '.apigateway.QuotaSettings',
      '10': 'quota'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.UsagePlan.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'throttle',
      '3': 395260638,
      '4': 1,
      '5': 11,
      '6': '.apigateway.ThrottleSettings',
      '10': 'throttle'
    },
  ],
  '3': [UsagePlan_TagsEntry$json],
};

@$core.Deprecated('Use usagePlanDescriptor instead')
const UsagePlan_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UsagePlan`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List usagePlanDescriptor = $convert.base64Decode(
    'CglVc2FnZVBsYW4SNQoJYXBpc3RhZ2VzGPGq5B4gAygLMhQuYXBpZ2F0ZXdheS5BcGlTdGFnZV'
    'IJYXBpc3RhZ2VzEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SEgoCaWQY'
    '4dXhuQEgASgJUgJpZBIVCgRuYW1lGOf75mkgASgJUgRuYW1lEiQKC3Byb2R1Y3Rjb2RlGOqAq/'
    '4BIAEoCVILcHJvZHVjdGNvZGUSMgoFcXVvdGEYjOuhdCABKAsyGS5hcGlnYXRld2F5LlF1b3Rh'
    'U2V0dGluZ3NSBXF1b3RhEjcKBHRhZ3MYodfboAEgAygLMh8uYXBpZ2F0ZXdheS5Vc2FnZVBsYW'
    '4uVGFnc0VudHJ5UgR0YWdzEjwKCHRocm90dGxlGN7lvLwBIAEoCzIcLmFwaWdhdGV3YXkuVGhy'
    'b3R0bGVTZXR0aW5nc1IIdGhyb3R0bGUaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5Eh'
    'QKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use usagePlanKeyDescriptor instead')
const UsagePlanKey$json = {
  '1': 'UsagePlanKey',
  '2': [
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'type', '3': 287830350, '4': 1, '5': 9, '10': 'type'},
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `UsagePlanKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List usagePlanKeyDescriptor = $convert.base64Decode(
    'CgxVc2FnZVBsYW5LZXkSEgoCaWQY4dXhuQEgASgJUgJpZBIVCgRuYW1lGOf75mkgASgJUgRuYW'
    '1lEhYKBHR5cGUYzuKfiQEgASgJUgR0eXBlEhcKBXZhbHVlGMun+xIgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use usagePlanKeysDescriptor instead')
const UsagePlanKeys$json = {
  '1': 'UsagePlanKeys',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.UsagePlanKey',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `UsagePlanKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List usagePlanKeysDescriptor = $convert.base64Decode(
    'Cg1Vc2FnZVBsYW5LZXlzEjIKBWl0ZW1zGJDn5NMBIAMoCzIYLmFwaWdhdGV3YXkuVXNhZ2VQbG'
    'FuS2V5UgVpdGVtcxIeCghwb3NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use usagePlansDescriptor instead')
const UsagePlans$json = {
  '1': 'UsagePlans',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.UsagePlan',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `UsagePlans`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List usagePlansDescriptor = $convert.base64Decode(
    'CgpVc2FnZVBsYW5zEi8KBWl0ZW1zGJDn5NMBIAMoCzIVLmFwaWdhdGV3YXkuVXNhZ2VQbGFuUg'
    'VpdGVtcxIeCghwb3NpdGlvbhiLnL2aASABKAlSCHBvc2l0aW9u');

@$core.Deprecated('Use vpcLinkDescriptor instead')
const VpcLink$json = {
  '1': 'VpcLink',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.apigateway.VpcLinkStatus',
      '10': 'status'
    },
    {
      '1': 'statusmessage',
      '3': 474462255,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.apigateway.VpcLink.TagsEntry',
      '10': 'tags'
    },
    {'1': 'targetarns', '3': 46319317, '4': 3, '5': 9, '10': 'targetarns'},
  ],
  '3': [VpcLink_TagsEntry$json],
};

@$core.Deprecated('Use vpcLinkDescriptor instead')
const VpcLink_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `VpcLink`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcLinkDescriptor = $convert.base64Decode(
    'CgdWcGNMaW5rEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SEgoCaWQY4d'
    'XhuQEgASgJUgJpZBIVCgRuYW1lGOf75mkgASgJUgRuYW1lEjUKBnN0YXR1cxjw763SASABKA4y'
    'GS5hcGlnYXRld2F5LlZwY0xpbmtTdGF0dXNSBnN0YXR1cxIoCg1zdGF0dXNtZXNzYWdlGK/wnu'
    'IBIAEoCVINc3RhdHVzbWVzc2FnZRI1CgR0YWdzGKHX26ABIAMoCzIdLmFwaWdhdGV3YXkuVnBj'
    'TGluay5UYWdzRW50cnlSBHRhZ3MSIQoKdGFyZ2V0YXJucxjVjYsWIAMoCVIKdGFyZ2V0YXJucx'
    'o3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4'
    'AQ==');

@$core.Deprecated('Use vpcLinksDescriptor instead')
const VpcLinks$json = {
  '1': 'VpcLinks',
  '2': [
    {
      '1': 'items',
      '3': 444150672,
      '4': 3,
      '5': 11,
      '6': '.apigateway.VpcLink',
      '10': 'items'
    },
    {'1': 'position', '3': 323964427, '4': 1, '5': 9, '10': 'position'},
  ],
};

/// Descriptor for `VpcLinks`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcLinksDescriptor = $convert.base64Decode(
    'CghWcGNMaW5rcxItCgVpdGVtcxiQ5+TTASADKAsyEy5hcGlnYXRld2F5LlZwY0xpbmtSBWl0ZW'
    '1zEh4KCHBvc2l0aW9uGIucvZoBIAEoCVIIcG9zaXRpb24=');
