// This is a generated file - do not edit.
//
// Generated from apigateway.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AccessAssociationSourceType extends $pb.ProtobufEnum {
  static const AccessAssociationSourceType ACCESS_ASSOCIATION_SOURCE_TYPE_VPCE =
      AccessAssociationSourceType._(
          0, _omitEnumNames ? '' : 'ACCESS_ASSOCIATION_SOURCE_TYPE_VPCE');

  static const $core.List<AccessAssociationSourceType> values =
      <AccessAssociationSourceType>[
    ACCESS_ASSOCIATION_SOURCE_TYPE_VPCE,
  ];

  static final $core.List<AccessAssociationSourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static AccessAssociationSourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AccessAssociationSourceType._(super.value, super.name);
}

class ApiKeySourceType extends $pb.ProtobufEnum {
  static const ApiKeySourceType API_KEY_SOURCE_TYPE_AUTHORIZER =
      ApiKeySourceType._(
          0, _omitEnumNames ? '' : 'API_KEY_SOURCE_TYPE_AUTHORIZER');
  static const ApiKeySourceType API_KEY_SOURCE_TYPE_HEADER =
      ApiKeySourceType._(1, _omitEnumNames ? '' : 'API_KEY_SOURCE_TYPE_HEADER');

  static const $core.List<ApiKeySourceType> values = <ApiKeySourceType>[
    API_KEY_SOURCE_TYPE_AUTHORIZER,
    API_KEY_SOURCE_TYPE_HEADER,
  ];

  static final $core.List<ApiKeySourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ApiKeySourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ApiKeySourceType._(super.value, super.name);
}

class ApiKeysFormat extends $pb.ProtobufEnum {
  static const ApiKeysFormat API_KEYS_FORMAT_CSV =
      ApiKeysFormat._(0, _omitEnumNames ? '' : 'API_KEYS_FORMAT_CSV');

  static const $core.List<ApiKeysFormat> values = <ApiKeysFormat>[
    API_KEYS_FORMAT_CSV,
  ];

  static final $core.List<ApiKeysFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ApiKeysFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ApiKeysFormat._(super.value, super.name);
}

class ApiStatus extends $pb.ProtobufEnum {
  static const ApiStatus API_STATUS_UPDATING =
      ApiStatus._(0, _omitEnumNames ? '' : 'API_STATUS_UPDATING');
  static const ApiStatus API_STATUS_PENDING =
      ApiStatus._(1, _omitEnumNames ? '' : 'API_STATUS_PENDING');
  static const ApiStatus API_STATUS_AVAILABLE =
      ApiStatus._(2, _omitEnumNames ? '' : 'API_STATUS_AVAILABLE');
  static const ApiStatus API_STATUS_FAILED =
      ApiStatus._(3, _omitEnumNames ? '' : 'API_STATUS_FAILED');

  static const $core.List<ApiStatus> values = <ApiStatus>[
    API_STATUS_UPDATING,
    API_STATUS_PENDING,
    API_STATUS_AVAILABLE,
    API_STATUS_FAILED,
  ];

  static final $core.List<ApiStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ApiStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ApiStatus._(super.value, super.name);
}

class AuthorizerType extends $pb.ProtobufEnum {
  static const AuthorizerType AUTHORIZER_TYPE_REQUEST =
      AuthorizerType._(0, _omitEnumNames ? '' : 'AUTHORIZER_TYPE_REQUEST');
  static const AuthorizerType AUTHORIZER_TYPE_COGNITO_USER_POOLS =
      AuthorizerType._(
          1, _omitEnumNames ? '' : 'AUTHORIZER_TYPE_COGNITO_USER_POOLS');
  static const AuthorizerType AUTHORIZER_TYPE_TOKEN =
      AuthorizerType._(2, _omitEnumNames ? '' : 'AUTHORIZER_TYPE_TOKEN');

  static const $core.List<AuthorizerType> values = <AuthorizerType>[
    AUTHORIZER_TYPE_REQUEST,
    AUTHORIZER_TYPE_COGNITO_USER_POOLS,
    AUTHORIZER_TYPE_TOKEN,
  ];

  static final $core.List<AuthorizerType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AuthorizerType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AuthorizerType._(super.value, super.name);
}

class CacheClusterSize extends $pb.ProtobufEnum {
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_237_GB =
      CacheClusterSize._(
          0, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_237_GB');
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_13_POINT_5_GB =
      CacheClusterSize._(
          1, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_13_POINT_5_GB');
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_118_GB =
      CacheClusterSize._(
          2, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_118_GB');
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_0_POINT_5_GB =
      CacheClusterSize._(
          3, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_0_POINT_5_GB');
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_1_POINT_6_GB =
      CacheClusterSize._(
          4, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_1_POINT_6_GB');
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_6_POINT_1_GB =
      CacheClusterSize._(
          5, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_6_POINT_1_GB');
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_28_POINT_4_GB =
      CacheClusterSize._(
          6, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_28_POINT_4_GB');
  static const CacheClusterSize CACHE_CLUSTER_SIZE_SIZE_58_POINT_2_GB =
      CacheClusterSize._(
          7, _omitEnumNames ? '' : 'CACHE_CLUSTER_SIZE_SIZE_58_POINT_2_GB');

  static const $core.List<CacheClusterSize> values = <CacheClusterSize>[
    CACHE_CLUSTER_SIZE_SIZE_237_GB,
    CACHE_CLUSTER_SIZE_SIZE_13_POINT_5_GB,
    CACHE_CLUSTER_SIZE_SIZE_118_GB,
    CACHE_CLUSTER_SIZE_SIZE_0_POINT_5_GB,
    CACHE_CLUSTER_SIZE_SIZE_1_POINT_6_GB,
    CACHE_CLUSTER_SIZE_SIZE_6_POINT_1_GB,
    CACHE_CLUSTER_SIZE_SIZE_28_POINT_4_GB,
    CACHE_CLUSTER_SIZE_SIZE_58_POINT_2_GB,
  ];

  static final $core.List<CacheClusterSize?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static CacheClusterSize? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CacheClusterSize._(super.value, super.name);
}

class CacheClusterStatus extends $pb.ProtobufEnum {
  static const CacheClusterStatus CACHE_CLUSTER_STATUS_FLUSH_IN_PROGRESS =
      CacheClusterStatus._(
          0, _omitEnumNames ? '' : 'CACHE_CLUSTER_STATUS_FLUSH_IN_PROGRESS');
  static const CacheClusterStatus CACHE_CLUSTER_STATUS_AVAILABLE =
      CacheClusterStatus._(
          1, _omitEnumNames ? '' : 'CACHE_CLUSTER_STATUS_AVAILABLE');
  static const CacheClusterStatus CACHE_CLUSTER_STATUS_CREATE_IN_PROGRESS =
      CacheClusterStatus._(
          2, _omitEnumNames ? '' : 'CACHE_CLUSTER_STATUS_CREATE_IN_PROGRESS');
  static const CacheClusterStatus CACHE_CLUSTER_STATUS_NOT_AVAILABLE =
      CacheClusterStatus._(
          3, _omitEnumNames ? '' : 'CACHE_CLUSTER_STATUS_NOT_AVAILABLE');
  static const CacheClusterStatus CACHE_CLUSTER_STATUS_DELETE_IN_PROGRESS =
      CacheClusterStatus._(
          4, _omitEnumNames ? '' : 'CACHE_CLUSTER_STATUS_DELETE_IN_PROGRESS');

  static const $core.List<CacheClusterStatus> values = <CacheClusterStatus>[
    CACHE_CLUSTER_STATUS_FLUSH_IN_PROGRESS,
    CACHE_CLUSTER_STATUS_AVAILABLE,
    CACHE_CLUSTER_STATUS_CREATE_IN_PROGRESS,
    CACHE_CLUSTER_STATUS_NOT_AVAILABLE,
    CACHE_CLUSTER_STATUS_DELETE_IN_PROGRESS,
  ];

  static final $core.List<CacheClusterStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static CacheClusterStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CacheClusterStatus._(super.value, super.name);
}

class ConnectionType extends $pb.ProtobufEnum {
  static const ConnectionType CONNECTION_TYPE_INTERNET =
      ConnectionType._(0, _omitEnumNames ? '' : 'CONNECTION_TYPE_INTERNET');
  static const ConnectionType CONNECTION_TYPE_VPC_LINK =
      ConnectionType._(1, _omitEnumNames ? '' : 'CONNECTION_TYPE_VPC_LINK');

  static const $core.List<ConnectionType> values = <ConnectionType>[
    CONNECTION_TYPE_INTERNET,
    CONNECTION_TYPE_VPC_LINK,
  ];

  static final $core.List<ConnectionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ConnectionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionType._(super.value, super.name);
}

class ContentHandlingStrategy extends $pb.ProtobufEnum {
  static const ContentHandlingStrategy
      CONTENT_HANDLING_STRATEGY_CONVERT_TO_BINARY = ContentHandlingStrategy._(0,
          _omitEnumNames ? '' : 'CONTENT_HANDLING_STRATEGY_CONVERT_TO_BINARY');
  static const ContentHandlingStrategy
      CONTENT_HANDLING_STRATEGY_CONVERT_TO_TEXT = ContentHandlingStrategy._(
          1, _omitEnumNames ? '' : 'CONTENT_HANDLING_STRATEGY_CONVERT_TO_TEXT');

  static const $core.List<ContentHandlingStrategy> values =
      <ContentHandlingStrategy>[
    CONTENT_HANDLING_STRATEGY_CONVERT_TO_BINARY,
    CONTENT_HANDLING_STRATEGY_CONVERT_TO_TEXT,
  ];

  static final $core.List<ContentHandlingStrategy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ContentHandlingStrategy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContentHandlingStrategy._(super.value, super.name);
}

class DocumentationPartType extends $pb.ProtobufEnum {
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_AUTHORIZER =
      DocumentationPartType._(
          0, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_AUTHORIZER');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_RESPONSE_HEADER =
      DocumentationPartType._(
          1, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_RESPONSE_HEADER');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_RESOURCE =
      DocumentationPartType._(
          2, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_RESOURCE');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_PATH_PARAMETER =
      DocumentationPartType._(
          3, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_PATH_PARAMETER');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_REQUEST_BODY =
      DocumentationPartType._(
          4, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_REQUEST_BODY');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_METHOD =
      DocumentationPartType._(
          5, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_METHOD');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_RESPONSE_BODY =
      DocumentationPartType._(
          6, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_RESPONSE_BODY');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_QUERY_PARAMETER =
      DocumentationPartType._(
          7, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_QUERY_PARAMETER');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_RESPONSE =
      DocumentationPartType._(
          8, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_RESPONSE');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_MODEL =
      DocumentationPartType._(
          9, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_MODEL');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_REQUEST_HEADER =
      DocumentationPartType._(
          10, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_REQUEST_HEADER');
  static const DocumentationPartType DOCUMENTATION_PART_TYPE_API =
      DocumentationPartType._(
          11, _omitEnumNames ? '' : 'DOCUMENTATION_PART_TYPE_API');

  static const $core.List<DocumentationPartType> values =
      <DocumentationPartType>[
    DOCUMENTATION_PART_TYPE_AUTHORIZER,
    DOCUMENTATION_PART_TYPE_RESPONSE_HEADER,
    DOCUMENTATION_PART_TYPE_RESOURCE,
    DOCUMENTATION_PART_TYPE_PATH_PARAMETER,
    DOCUMENTATION_PART_TYPE_REQUEST_BODY,
    DOCUMENTATION_PART_TYPE_METHOD,
    DOCUMENTATION_PART_TYPE_RESPONSE_BODY,
    DOCUMENTATION_PART_TYPE_QUERY_PARAMETER,
    DOCUMENTATION_PART_TYPE_RESPONSE,
    DOCUMENTATION_PART_TYPE_MODEL,
    DOCUMENTATION_PART_TYPE_REQUEST_HEADER,
    DOCUMENTATION_PART_TYPE_API,
  ];

  static final $core.List<DocumentationPartType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 11);
  static DocumentationPartType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentationPartType._(super.value, super.name);
}

class DomainNameStatus extends $pb.ProtobufEnum {
  static const DomainNameStatus DOMAIN_NAME_STATUS_UPDATING =
      DomainNameStatus._(
          0, _omitEnumNames ? '' : 'DOMAIN_NAME_STATUS_UPDATING');
  static const DomainNameStatus DOMAIN_NAME_STATUS_PENDING =
      DomainNameStatus._(1, _omitEnumNames ? '' : 'DOMAIN_NAME_STATUS_PENDING');
  static const DomainNameStatus
      DOMAIN_NAME_STATUS_PENDING_CERTIFICATE_REIMPORT = DomainNameStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'DOMAIN_NAME_STATUS_PENDING_CERTIFICATE_REIMPORT');
  static const DomainNameStatus DOMAIN_NAME_STATUS_AVAILABLE =
      DomainNameStatus._(
          3, _omitEnumNames ? '' : 'DOMAIN_NAME_STATUS_AVAILABLE');
  static const DomainNameStatus
      DOMAIN_NAME_STATUS_PENDING_OWNERSHIP_VERIFICATION = DomainNameStatus._(
          4,
          _omitEnumNames
              ? ''
              : 'DOMAIN_NAME_STATUS_PENDING_OWNERSHIP_VERIFICATION');
  static const DomainNameStatus DOMAIN_NAME_STATUS_FAILED =
      DomainNameStatus._(5, _omitEnumNames ? '' : 'DOMAIN_NAME_STATUS_FAILED');

  static const $core.List<DomainNameStatus> values = <DomainNameStatus>[
    DOMAIN_NAME_STATUS_UPDATING,
    DOMAIN_NAME_STATUS_PENDING,
    DOMAIN_NAME_STATUS_PENDING_CERTIFICATE_REIMPORT,
    DOMAIN_NAME_STATUS_AVAILABLE,
    DOMAIN_NAME_STATUS_PENDING_OWNERSHIP_VERIFICATION,
    DOMAIN_NAME_STATUS_FAILED,
  ];

  static final $core.List<DomainNameStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static DomainNameStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DomainNameStatus._(super.value, super.name);
}

class EndpointAccessMode extends $pb.ProtobufEnum {
  static const EndpointAccessMode ENDPOINT_ACCESS_MODE_STRICT =
      EndpointAccessMode._(
          0, _omitEnumNames ? '' : 'ENDPOINT_ACCESS_MODE_STRICT');
  static const EndpointAccessMode ENDPOINT_ACCESS_MODE_BASIC =
      EndpointAccessMode._(
          1, _omitEnumNames ? '' : 'ENDPOINT_ACCESS_MODE_BASIC');

  static const $core.List<EndpointAccessMode> values = <EndpointAccessMode>[
    ENDPOINT_ACCESS_MODE_STRICT,
    ENDPOINT_ACCESS_MODE_BASIC,
  ];

  static final $core.List<EndpointAccessMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static EndpointAccessMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EndpointAccessMode._(super.value, super.name);
}

class EndpointType extends $pb.ProtobufEnum {
  static const EndpointType ENDPOINT_TYPE_PRIVATE =
      EndpointType._(0, _omitEnumNames ? '' : 'ENDPOINT_TYPE_PRIVATE');
  static const EndpointType ENDPOINT_TYPE_REGIONAL =
      EndpointType._(1, _omitEnumNames ? '' : 'ENDPOINT_TYPE_REGIONAL');
  static const EndpointType ENDPOINT_TYPE_EDGE =
      EndpointType._(2, _omitEnumNames ? '' : 'ENDPOINT_TYPE_EDGE');

  static const $core.List<EndpointType> values = <EndpointType>[
    ENDPOINT_TYPE_PRIVATE,
    ENDPOINT_TYPE_REGIONAL,
    ENDPOINT_TYPE_EDGE,
  ];

  static final $core.List<EndpointType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EndpointType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EndpointType._(super.value, super.name);
}

class GatewayResponseType extends $pb.ProtobufEnum {
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_WAF_FILTERED =
      GatewayResponseType._(
          0, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_WAF_FILTERED');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_DEFAULT_4XX =
      GatewayResponseType._(
          1, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_DEFAULT_4XX');
  static const GatewayResponseType
      GATEWAY_RESPONSE_TYPE_UNSUPPORTED_MEDIA_TYPE = GatewayResponseType._(2,
          _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_UNSUPPORTED_MEDIA_TYPE');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_REQUEST_TOO_LARGE =
      GatewayResponseType._(
          3, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_REQUEST_TOO_LARGE');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_ACCESS_DENIED =
      GatewayResponseType._(
          4, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_ACCESS_DENIED');
  static const GatewayResponseType
      GATEWAY_RESPONSE_TYPE_MISSING_AUTHENTICATION_TOKEN =
      GatewayResponseType._(
          5,
          _omitEnumNames
              ? ''
              : 'GATEWAY_RESPONSE_TYPE_MISSING_AUTHENTICATION_TOKEN');
  static const GatewayResponseType
      GATEWAY_RESPONSE_TYPE_API_CONFIGURATION_ERROR = GatewayResponseType._(
          6,
          _omitEnumNames
              ? ''
              : 'GATEWAY_RESPONSE_TYPE_API_CONFIGURATION_ERROR');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_THROTTLED =
      GatewayResponseType._(
          7, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_THROTTLED');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_INTEGRATION_FAILURE =
      GatewayResponseType._(
          8, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_INTEGRATION_FAILURE');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_RESOURCE_NOT_FOUND =
      GatewayResponseType._(
          9, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_RESOURCE_NOT_FOUND');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_QUOTA_EXCEEDED =
      GatewayResponseType._(
          10, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_QUOTA_EXCEEDED');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_AUTHORIZER_FAILURE =
      GatewayResponseType._(
          11, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_AUTHORIZER_FAILURE');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_INVALID_SIGNATURE =
      GatewayResponseType._(
          12, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_INVALID_SIGNATURE');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_INTEGRATION_TIMEOUT =
      GatewayResponseType._(13,
          _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_INTEGRATION_TIMEOUT');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_UNAUTHORIZED =
      GatewayResponseType._(
          14, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_UNAUTHORIZED');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_EXPIRED_TOKEN =
      GatewayResponseType._(
          15, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_EXPIRED_TOKEN');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_INVALID_API_KEY =
      GatewayResponseType._(
          16, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_INVALID_API_KEY');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_BAD_REQUEST_BODY =
      GatewayResponseType._(
          17, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_BAD_REQUEST_BODY');
  static const GatewayResponseType
      GATEWAY_RESPONSE_TYPE_BAD_REQUEST_PARAMETERS = GatewayResponseType._(18,
          _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_BAD_REQUEST_PARAMETERS');
  static const GatewayResponseType GATEWAY_RESPONSE_TYPE_DEFAULT_5XX =
      GatewayResponseType._(
          19, _omitEnumNames ? '' : 'GATEWAY_RESPONSE_TYPE_DEFAULT_5XX');
  static const GatewayResponseType
      GATEWAY_RESPONSE_TYPE_AUTHORIZER_CONFIGURATION_ERROR =
      GatewayResponseType._(
          20,
          _omitEnumNames
              ? ''
              : 'GATEWAY_RESPONSE_TYPE_AUTHORIZER_CONFIGURATION_ERROR');

  static const $core.List<GatewayResponseType> values = <GatewayResponseType>[
    GATEWAY_RESPONSE_TYPE_WAF_FILTERED,
    GATEWAY_RESPONSE_TYPE_DEFAULT_4XX,
    GATEWAY_RESPONSE_TYPE_UNSUPPORTED_MEDIA_TYPE,
    GATEWAY_RESPONSE_TYPE_REQUEST_TOO_LARGE,
    GATEWAY_RESPONSE_TYPE_ACCESS_DENIED,
    GATEWAY_RESPONSE_TYPE_MISSING_AUTHENTICATION_TOKEN,
    GATEWAY_RESPONSE_TYPE_API_CONFIGURATION_ERROR,
    GATEWAY_RESPONSE_TYPE_THROTTLED,
    GATEWAY_RESPONSE_TYPE_INTEGRATION_FAILURE,
    GATEWAY_RESPONSE_TYPE_RESOURCE_NOT_FOUND,
    GATEWAY_RESPONSE_TYPE_QUOTA_EXCEEDED,
    GATEWAY_RESPONSE_TYPE_AUTHORIZER_FAILURE,
    GATEWAY_RESPONSE_TYPE_INVALID_SIGNATURE,
    GATEWAY_RESPONSE_TYPE_INTEGRATION_TIMEOUT,
    GATEWAY_RESPONSE_TYPE_UNAUTHORIZED,
    GATEWAY_RESPONSE_TYPE_EXPIRED_TOKEN,
    GATEWAY_RESPONSE_TYPE_INVALID_API_KEY,
    GATEWAY_RESPONSE_TYPE_BAD_REQUEST_BODY,
    GATEWAY_RESPONSE_TYPE_BAD_REQUEST_PARAMETERS,
    GATEWAY_RESPONSE_TYPE_DEFAULT_5XX,
    GATEWAY_RESPONSE_TYPE_AUTHORIZER_CONFIGURATION_ERROR,
  ];

  static final $core.List<GatewayResponseType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 20);
  static GatewayResponseType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GatewayResponseType._(super.value, super.name);
}

class IntegrationType extends $pb.ProtobufEnum {
  static const IntegrationType INTEGRATION_TYPE_HTTP =
      IntegrationType._(0, _omitEnumNames ? '' : 'INTEGRATION_TYPE_HTTP');
  static const IntegrationType INTEGRATION_TYPE_AWS_PROXY =
      IntegrationType._(1, _omitEnumNames ? '' : 'INTEGRATION_TYPE_AWS_PROXY');
  static const IntegrationType INTEGRATION_TYPE_MOCK =
      IntegrationType._(2, _omitEnumNames ? '' : 'INTEGRATION_TYPE_MOCK');
  static const IntegrationType INTEGRATION_TYPE_AWS =
      IntegrationType._(3, _omitEnumNames ? '' : 'INTEGRATION_TYPE_AWS');
  static const IntegrationType INTEGRATION_TYPE_HTTP_PROXY =
      IntegrationType._(4, _omitEnumNames ? '' : 'INTEGRATION_TYPE_HTTP_PROXY');

  static const $core.List<IntegrationType> values = <IntegrationType>[
    INTEGRATION_TYPE_HTTP,
    INTEGRATION_TYPE_AWS_PROXY,
    INTEGRATION_TYPE_MOCK,
    INTEGRATION_TYPE_AWS,
    INTEGRATION_TYPE_HTTP_PROXY,
  ];

  static final $core.List<IntegrationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static IntegrationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IntegrationType._(super.value, super.name);
}

class IpAddressType extends $pb.ProtobufEnum {
  static const IpAddressType IP_ADDRESS_TYPE_DUALSTACK =
      IpAddressType._(0, _omitEnumNames ? '' : 'IP_ADDRESS_TYPE_DUALSTACK');
  static const IpAddressType IP_ADDRESS_TYPE_IPV4 =
      IpAddressType._(1, _omitEnumNames ? '' : 'IP_ADDRESS_TYPE_IPV4');

  static const $core.List<IpAddressType> values = <IpAddressType>[
    IP_ADDRESS_TYPE_DUALSTACK,
    IP_ADDRESS_TYPE_IPV4,
  ];

  static final $core.List<IpAddressType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IpAddressType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IpAddressType._(super.value, super.name);
}

class LocationStatusType extends $pb.ProtobufEnum {
  static const LocationStatusType LOCATION_STATUS_TYPE_DOCUMENTED =
      LocationStatusType._(
          0, _omitEnumNames ? '' : 'LOCATION_STATUS_TYPE_DOCUMENTED');
  static const LocationStatusType LOCATION_STATUS_TYPE_UNDOCUMENTED =
      LocationStatusType._(
          1, _omitEnumNames ? '' : 'LOCATION_STATUS_TYPE_UNDOCUMENTED');

  static const $core.List<LocationStatusType> values = <LocationStatusType>[
    LOCATION_STATUS_TYPE_DOCUMENTED,
    LOCATION_STATUS_TYPE_UNDOCUMENTED,
  ];

  static final $core.List<LocationStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static LocationStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LocationStatusType._(super.value, super.name);
}

class Op extends $pb.ProtobufEnum {
  static const Op OP_ADD = Op._(0, _omitEnumNames ? '' : 'OP_ADD');
  static const Op OP_REMOVE = Op._(1, _omitEnumNames ? '' : 'OP_REMOVE');
  static const Op OP_MOVE = Op._(2, _omitEnumNames ? '' : 'OP_MOVE');
  static const Op OP_TEST = Op._(3, _omitEnumNames ? '' : 'OP_TEST');
  static const Op OP_REPLACE = Op._(4, _omitEnumNames ? '' : 'OP_REPLACE');
  static const Op OP_COPY = Op._(5, _omitEnumNames ? '' : 'OP_COPY');

  static const $core.List<Op> values = <Op>[
    OP_ADD,
    OP_REMOVE,
    OP_MOVE,
    OP_TEST,
    OP_REPLACE,
    OP_COPY,
  ];

  static final $core.List<Op?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static Op? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Op._(super.value, super.name);
}

class PutMode extends $pb.ProtobufEnum {
  static const PutMode PUT_MODE_OVERWRITE =
      PutMode._(0, _omitEnumNames ? '' : 'PUT_MODE_OVERWRITE');
  static const PutMode PUT_MODE_MERGE =
      PutMode._(1, _omitEnumNames ? '' : 'PUT_MODE_MERGE');

  static const $core.List<PutMode> values = <PutMode>[
    PUT_MODE_OVERWRITE,
    PUT_MODE_MERGE,
  ];

  static final $core.List<PutMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PutMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PutMode._(super.value, super.name);
}

class QuotaPeriodType extends $pb.ProtobufEnum {
  static const QuotaPeriodType QUOTA_PERIOD_TYPE_WEEK =
      QuotaPeriodType._(0, _omitEnumNames ? '' : 'QUOTA_PERIOD_TYPE_WEEK');
  static const QuotaPeriodType QUOTA_PERIOD_TYPE_DAY =
      QuotaPeriodType._(1, _omitEnumNames ? '' : 'QUOTA_PERIOD_TYPE_DAY');
  static const QuotaPeriodType QUOTA_PERIOD_TYPE_MONTH =
      QuotaPeriodType._(2, _omitEnumNames ? '' : 'QUOTA_PERIOD_TYPE_MONTH');

  static const $core.List<QuotaPeriodType> values = <QuotaPeriodType>[
    QUOTA_PERIOD_TYPE_WEEK,
    QUOTA_PERIOD_TYPE_DAY,
    QUOTA_PERIOD_TYPE_MONTH,
  ];

  static final $core.List<QuotaPeriodType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static QuotaPeriodType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QuotaPeriodType._(super.value, super.name);
}

class ResourceOwner extends $pb.ProtobufEnum {
  static const ResourceOwner RESOURCE_OWNER_OTHER_ACCOUNTS =
      ResourceOwner._(0, _omitEnumNames ? '' : 'RESOURCE_OWNER_OTHER_ACCOUNTS');
  static const ResourceOwner RESOURCE_OWNER_SELF =
      ResourceOwner._(1, _omitEnumNames ? '' : 'RESOURCE_OWNER_SELF');

  static const $core.List<ResourceOwner> values = <ResourceOwner>[
    RESOURCE_OWNER_OTHER_ACCOUNTS,
    RESOURCE_OWNER_SELF,
  ];

  static final $core.List<ResourceOwner?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ResourceOwner? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceOwner._(super.value, super.name);
}

class ResponseTransferMode extends $pb.ProtobufEnum {
  static const ResponseTransferMode RESPONSE_TRANSFER_MODE_STREAM =
      ResponseTransferMode._(
          0, _omitEnumNames ? '' : 'RESPONSE_TRANSFER_MODE_STREAM');
  static const ResponseTransferMode RESPONSE_TRANSFER_MODE_BUFFERED =
      ResponseTransferMode._(
          1, _omitEnumNames ? '' : 'RESPONSE_TRANSFER_MODE_BUFFERED');

  static const $core.List<ResponseTransferMode> values = <ResponseTransferMode>[
    RESPONSE_TRANSFER_MODE_STREAM,
    RESPONSE_TRANSFER_MODE_BUFFERED,
  ];

  static final $core.List<ResponseTransferMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ResponseTransferMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResponseTransferMode._(super.value, super.name);
}

class RoutingMode extends $pb.ProtobufEnum {
  static const RoutingMode ROUTING_MODE_ROUTING_RULE_THEN_BASE_PATH_MAPPING =
      RoutingMode._(
          0,
          _omitEnumNames
              ? ''
              : 'ROUTING_MODE_ROUTING_RULE_THEN_BASE_PATH_MAPPING');
  static const RoutingMode ROUTING_MODE_BASE_PATH_MAPPING_ONLY = RoutingMode._(
      1, _omitEnumNames ? '' : 'ROUTING_MODE_BASE_PATH_MAPPING_ONLY');
  static const RoutingMode ROUTING_MODE_ROUTING_RULE_ONLY =
      RoutingMode._(2, _omitEnumNames ? '' : 'ROUTING_MODE_ROUTING_RULE_ONLY');

  static const $core.List<RoutingMode> values = <RoutingMode>[
    ROUTING_MODE_ROUTING_RULE_THEN_BASE_PATH_MAPPING,
    ROUTING_MODE_BASE_PATH_MAPPING_ONLY,
    ROUTING_MODE_ROUTING_RULE_ONLY,
  ];

  static final $core.List<RoutingMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static RoutingMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RoutingMode._(super.value, super.name);
}

class SecurityPolicy extends $pb.ProtobufEnum {
  static const SecurityPolicy
      SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PQ_2025_09 = SecurityPolicy._(
          0,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PQ_2025_09');
  static const SecurityPolicy SECURITY_POLICY_SECURITYPOLICY_TLS13_2025_EDGE =
      SecurityPolicy._(
          1,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS13_2025_EDGE');
  static const SecurityPolicy
      SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PQ_2025_09 =
      SecurityPolicy._(
          2,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PQ_2025_09');
  static const SecurityPolicy SECURITY_POLICY_SECURITYPOLICY_TLS12_2018_EDGE =
      SecurityPolicy._(
          3,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS12_2018_EDGE');
  static const SecurityPolicy
      SECURITY_POLICY_SECURITYPOLICY_TLS12_PFS_2025_EDGE = SecurityPolicy._(
          4,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS12_PFS_2025_EDGE');
  static const SecurityPolicy SECURITY_POLICY_TLS_1_2 =
      SecurityPolicy._(5, _omitEnumNames ? '' : 'SECURITY_POLICY_TLS_1_2');
  static const SecurityPolicy SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_2021_06 =
      SecurityPolicy._(
          6,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_2021_06');
  static const SecurityPolicy SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_2025_09 =
      SecurityPolicy._(
          7,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_2025_09');
  static const SecurityPolicy
      SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PFS_PQ_2025_09 =
      SecurityPolicy._(
          8,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PFS_PQ_2025_09');
  static const SecurityPolicy
      SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_FIPS_2025_09 = SecurityPolicy._(
          9,
          _omitEnumNames
              ? ''
              : 'SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_FIPS_2025_09');
  static const SecurityPolicy SECURITY_POLICY_TLS_1_0 =
      SecurityPolicy._(10, _omitEnumNames ? '' : 'SECURITY_POLICY_TLS_1_0');

  static const $core.List<SecurityPolicy> values = <SecurityPolicy>[
    SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PQ_2025_09,
    SECURITY_POLICY_SECURITYPOLICY_TLS13_2025_EDGE,
    SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PQ_2025_09,
    SECURITY_POLICY_SECURITYPOLICY_TLS12_2018_EDGE,
    SECURITY_POLICY_SECURITYPOLICY_TLS12_PFS_2025_EDGE,
    SECURITY_POLICY_TLS_1_2,
    SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_2021_06,
    SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_2025_09,
    SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PFS_PQ_2025_09,
    SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_FIPS_2025_09,
    SECURITY_POLICY_TLS_1_0,
  ];

  static final $core.List<SecurityPolicy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 10);
  static SecurityPolicy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SecurityPolicy._(super.value, super.name);
}

class UnauthorizedCacheControlHeaderStrategy extends $pb.ProtobufEnum {
  static const UnauthorizedCacheControlHeaderStrategy
      UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITH_RESPONSE_HEADER =
      UnauthorizedCacheControlHeaderStrategy._(
          0,
          _omitEnumNames
              ? ''
              : 'UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITH_RESPONSE_HEADER');
  static const UnauthorizedCacheControlHeaderStrategy
      UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_FAIL_WITH_403 =
      UnauthorizedCacheControlHeaderStrategy._(
          1,
          _omitEnumNames
              ? ''
              : 'UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_FAIL_WITH_403');
  static const UnauthorizedCacheControlHeaderStrategy
      UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITHOUT_RESPONSE_HEADER =
      UnauthorizedCacheControlHeaderStrategy._(
          2,
          _omitEnumNames
              ? ''
              : 'UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITHOUT_RESPONSE_HEADER');

  static const $core.List<UnauthorizedCacheControlHeaderStrategy> values =
      <UnauthorizedCacheControlHeaderStrategy>[
    UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITH_RESPONSE_HEADER,
    UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_FAIL_WITH_403,
    UNAUTHORIZED_CACHE_CONTROL_HEADER_STRATEGY_SUCCEED_WITHOUT_RESPONSE_HEADER,
  ];

  static final $core.List<UnauthorizedCacheControlHeaderStrategy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static UnauthorizedCacheControlHeaderStrategy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UnauthorizedCacheControlHeaderStrategy._(super.value, super.name);
}

class VpcLinkStatus extends $pb.ProtobufEnum {
  static const VpcLinkStatus VPC_LINK_STATUS_PENDING =
      VpcLinkStatus._(0, _omitEnumNames ? '' : 'VPC_LINK_STATUS_PENDING');
  static const VpcLinkStatus VPC_LINK_STATUS_AVAILABLE =
      VpcLinkStatus._(1, _omitEnumNames ? '' : 'VPC_LINK_STATUS_AVAILABLE');
  static const VpcLinkStatus VPC_LINK_STATUS_DELETING =
      VpcLinkStatus._(2, _omitEnumNames ? '' : 'VPC_LINK_STATUS_DELETING');
  static const VpcLinkStatus VPC_LINK_STATUS_FAILED =
      VpcLinkStatus._(3, _omitEnumNames ? '' : 'VPC_LINK_STATUS_FAILED');

  static const $core.List<VpcLinkStatus> values = <VpcLinkStatus>[
    VPC_LINK_STATUS_PENDING,
    VPC_LINK_STATUS_AVAILABLE,
    VPC_LINK_STATUS_DELETING,
    VPC_LINK_STATUS_FAILED,
  ];

  static final $core.List<VpcLinkStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static VpcLinkStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VpcLinkStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
