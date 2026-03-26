// This is a generated file - do not edit.
//
// Generated from cloudfront.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class CachePolicyCookieBehavior extends $pb.ProtobufEnum {
  static const CachePolicyCookieBehavior CACHE_POLICY_COOKIE_BEHAVIOR_ALL =
      CachePolicyCookieBehavior._(
          0, _omitEnumNames ? '' : 'CACHE_POLICY_COOKIE_BEHAVIOR_ALL');
  static const CachePolicyCookieBehavior
      CACHE_POLICY_COOKIE_BEHAVIOR_WHITELIST = CachePolicyCookieBehavior._(
          1, _omitEnumNames ? '' : 'CACHE_POLICY_COOKIE_BEHAVIOR_WHITELIST');
  static const CachePolicyCookieBehavior CACHE_POLICY_COOKIE_BEHAVIOR_NONE =
      CachePolicyCookieBehavior._(
          2, _omitEnumNames ? '' : 'CACHE_POLICY_COOKIE_BEHAVIOR_NONE');
  static const CachePolicyCookieBehavior
      CACHE_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT = CachePolicyCookieBehavior._(
          3, _omitEnumNames ? '' : 'CACHE_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT');

  static const $core.List<CachePolicyCookieBehavior> values =
      <CachePolicyCookieBehavior>[
    CACHE_POLICY_COOKIE_BEHAVIOR_ALL,
    CACHE_POLICY_COOKIE_BEHAVIOR_WHITELIST,
    CACHE_POLICY_COOKIE_BEHAVIOR_NONE,
    CACHE_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT,
  ];

  static final $core.List<CachePolicyCookieBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static CachePolicyCookieBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CachePolicyCookieBehavior._(super.value, super.name);
}

class CachePolicyHeaderBehavior extends $pb.ProtobufEnum {
  static const CachePolicyHeaderBehavior
      CACHE_POLICY_HEADER_BEHAVIOR_WHITELIST = CachePolicyHeaderBehavior._(
          0, _omitEnumNames ? '' : 'CACHE_POLICY_HEADER_BEHAVIOR_WHITELIST');
  static const CachePolicyHeaderBehavior CACHE_POLICY_HEADER_BEHAVIOR_NONE =
      CachePolicyHeaderBehavior._(
          1, _omitEnumNames ? '' : 'CACHE_POLICY_HEADER_BEHAVIOR_NONE');

  static const $core.List<CachePolicyHeaderBehavior> values =
      <CachePolicyHeaderBehavior>[
    CACHE_POLICY_HEADER_BEHAVIOR_WHITELIST,
    CACHE_POLICY_HEADER_BEHAVIOR_NONE,
  ];

  static final $core.List<CachePolicyHeaderBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CachePolicyHeaderBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CachePolicyHeaderBehavior._(super.value, super.name);
}

class CachePolicyQueryStringBehavior extends $pb.ProtobufEnum {
  static const CachePolicyQueryStringBehavior
      CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALL = CachePolicyQueryStringBehavior._(
          0, _omitEnumNames ? '' : 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALL');
  static const CachePolicyQueryStringBehavior
      CACHE_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST =
      CachePolicyQueryStringBehavior._(1,
          _omitEnumNames ? '' : 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST');
  static const CachePolicyQueryStringBehavior
      CACHE_POLICY_QUERY_STRING_BEHAVIOR_NONE =
      CachePolicyQueryStringBehavior._(
          2, _omitEnumNames ? '' : 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_NONE');
  static const CachePolicyQueryStringBehavior
      CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT =
      CachePolicyQueryStringBehavior._(3,
          _omitEnumNames ? '' : 'CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT');

  static const $core.List<CachePolicyQueryStringBehavior> values =
      <CachePolicyQueryStringBehavior>[
    CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALL,
    CACHE_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST,
    CACHE_POLICY_QUERY_STRING_BEHAVIOR_NONE,
    CACHE_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT,
  ];

  static final $core.List<CachePolicyQueryStringBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static CachePolicyQueryStringBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CachePolicyQueryStringBehavior._(super.value, super.name);
}

class CachePolicyType extends $pb.ProtobufEnum {
  static const CachePolicyType CACHE_POLICY_TYPE_MANAGED =
      CachePolicyType._(0, _omitEnumNames ? '' : 'CACHE_POLICY_TYPE_MANAGED');
  static const CachePolicyType CACHE_POLICY_TYPE_CUSTOM =
      CachePolicyType._(1, _omitEnumNames ? '' : 'CACHE_POLICY_TYPE_CUSTOM');

  static const $core.List<CachePolicyType> values = <CachePolicyType>[
    CACHE_POLICY_TYPE_MANAGED,
    CACHE_POLICY_TYPE_CUSTOM,
  ];

  static final $core.List<CachePolicyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CachePolicyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CachePolicyType._(super.value, super.name);
}

class CertificateSource extends $pb.ProtobufEnum {
  static const CertificateSource CERTIFICATE_SOURCE_IAM =
      CertificateSource._(0, _omitEnumNames ? '' : 'CERTIFICATE_SOURCE_IAM');
  static const CertificateSource CERTIFICATE_SOURCE_ACM =
      CertificateSource._(1, _omitEnumNames ? '' : 'CERTIFICATE_SOURCE_ACM');
  static const CertificateSource CERTIFICATE_SOURCE_CLOUDFRONT =
      CertificateSource._(
          2, _omitEnumNames ? '' : 'CERTIFICATE_SOURCE_CLOUDFRONT');

  static const $core.List<CertificateSource> values = <CertificateSource>[
    CERTIFICATE_SOURCE_IAM,
    CERTIFICATE_SOURCE_ACM,
    CERTIFICATE_SOURCE_CLOUDFRONT,
  ];

  static final $core.List<CertificateSource?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static CertificateSource? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CertificateSource._(super.value, super.name);
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

class ConnectionMode extends $pb.ProtobufEnum {
  static const ConnectionMode CONNECTION_MODE_DIRECT =
      ConnectionMode._(0, _omitEnumNames ? '' : 'CONNECTION_MODE_DIRECT');
  static const ConnectionMode CONNECTION_MODE_TENANTONLY =
      ConnectionMode._(1, _omitEnumNames ? '' : 'CONNECTION_MODE_TENANTONLY');

  static const $core.List<ConnectionMode> values = <ConnectionMode>[
    CONNECTION_MODE_DIRECT,
    CONNECTION_MODE_TENANTONLY,
  ];

  static final $core.List<ConnectionMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ConnectionMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionMode._(super.value, super.name);
}

class ContinuousDeploymentPolicyType extends $pb.ProtobufEnum {
  static const ContinuousDeploymentPolicyType
      CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEHEADER =
      ContinuousDeploymentPolicyType._(
          0,
          _omitEnumNames
              ? ''
              : 'CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEHEADER');
  static const ContinuousDeploymentPolicyType
      CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEWEIGHT =
      ContinuousDeploymentPolicyType._(
          1,
          _omitEnumNames
              ? ''
              : 'CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEWEIGHT');

  static const $core.List<ContinuousDeploymentPolicyType> values =
      <ContinuousDeploymentPolicyType>[
    CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEHEADER,
    CONTINUOUS_DEPLOYMENT_POLICY_TYPE_SINGLEWEIGHT,
  ];

  static final $core.List<ContinuousDeploymentPolicyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ContinuousDeploymentPolicyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContinuousDeploymentPolicyType._(super.value, super.name);
}

class CustomizationActionType extends $pb.ProtobufEnum {
  static const CustomizationActionType CUSTOMIZATION_ACTION_TYPE_OVERRIDE =
      CustomizationActionType._(
          0, _omitEnumNames ? '' : 'CUSTOMIZATION_ACTION_TYPE_OVERRIDE');
  static const CustomizationActionType CUSTOMIZATION_ACTION_TYPE_DISABLE =
      CustomizationActionType._(
          1, _omitEnumNames ? '' : 'CUSTOMIZATION_ACTION_TYPE_DISABLE');

  static const $core.List<CustomizationActionType> values =
      <CustomizationActionType>[
    CUSTOMIZATION_ACTION_TYPE_OVERRIDE,
    CUSTOMIZATION_ACTION_TYPE_DISABLE,
  ];

  static final $core.List<CustomizationActionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CustomizationActionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CustomizationActionType._(super.value, super.name);
}

class DistributionResourceType extends $pb.ProtobufEnum {
  static const DistributionResourceType
      DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTION = DistributionResourceType._(
          0, _omitEnumNames ? '' : 'DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTION');
  static const DistributionResourceType
      DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTIONTENANT =
      DistributionResourceType._(
          1,
          _omitEnumNames
              ? ''
              : 'DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTIONTENANT');

  static const $core.List<DistributionResourceType> values =
      <DistributionResourceType>[
    DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTION,
    DISTRIBUTION_RESOURCE_TYPE_DISTRIBUTIONTENANT,
  ];

  static final $core.List<DistributionResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DistributionResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DistributionResourceType._(super.value, super.name);
}

class DnsConfigurationStatus extends $pb.ProtobufEnum {
  static const DnsConfigurationStatus DNS_CONFIGURATION_STATUS_VALID =
      DnsConfigurationStatus._(
          0, _omitEnumNames ? '' : 'DNS_CONFIGURATION_STATUS_VALID');
  static const DnsConfigurationStatus DNS_CONFIGURATION_STATUS_INVALID =
      DnsConfigurationStatus._(
          1, _omitEnumNames ? '' : 'DNS_CONFIGURATION_STATUS_INVALID');
  static const DnsConfigurationStatus DNS_CONFIGURATION_STATUS_UNKNOWN =
      DnsConfigurationStatus._(
          2, _omitEnumNames ? '' : 'DNS_CONFIGURATION_STATUS_UNKNOWN');

  static const $core.List<DnsConfigurationStatus> values =
      <DnsConfigurationStatus>[
    DNS_CONFIGURATION_STATUS_VALID,
    DNS_CONFIGURATION_STATUS_INVALID,
    DNS_CONFIGURATION_STATUS_UNKNOWN,
  ];

  static final $core.List<DnsConfigurationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static DnsConfigurationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DnsConfigurationStatus._(super.value, super.name);
}

class DomainStatus extends $pb.ProtobufEnum {
  static const DomainStatus DOMAIN_STATUS_ACTIVE =
      DomainStatus._(0, _omitEnumNames ? '' : 'DOMAIN_STATUS_ACTIVE');
  static const DomainStatus DOMAIN_STATUS_INACTIVE =
      DomainStatus._(1, _omitEnumNames ? '' : 'DOMAIN_STATUS_INACTIVE');

  static const $core.List<DomainStatus> values = <DomainStatus>[
    DOMAIN_STATUS_ACTIVE,
    DOMAIN_STATUS_INACTIVE,
  ];

  static final $core.List<DomainStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DomainStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DomainStatus._(super.value, super.name);
}

class EventType extends $pb.ProtobufEnum {
  static const EventType EVENT_TYPE_VIEWER_REQUEST =
      EventType._(0, _omitEnumNames ? '' : 'EVENT_TYPE_VIEWER_REQUEST');
  static const EventType EVENT_TYPE_ORIGIN_REQUEST =
      EventType._(1, _omitEnumNames ? '' : 'EVENT_TYPE_ORIGIN_REQUEST');
  static const EventType EVENT_TYPE_VIEWER_RESPONSE =
      EventType._(2, _omitEnumNames ? '' : 'EVENT_TYPE_VIEWER_RESPONSE');
  static const EventType EVENT_TYPE_ORIGIN_RESPONSE =
      EventType._(3, _omitEnumNames ? '' : 'EVENT_TYPE_ORIGIN_RESPONSE');

  static const $core.List<EventType> values = <EventType>[
    EVENT_TYPE_VIEWER_REQUEST,
    EVENT_TYPE_ORIGIN_REQUEST,
    EVENT_TYPE_VIEWER_RESPONSE,
    EVENT_TYPE_ORIGIN_RESPONSE,
  ];

  static final $core.List<EventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static EventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventType._(super.value, super.name);
}

class Format extends $pb.ProtobufEnum {
  static const Format FORMAT_URLENCODED =
      Format._(0, _omitEnumNames ? '' : 'FORMAT_URLENCODED');

  static const $core.List<Format> values = <Format>[
    FORMAT_URLENCODED,
  ];

  static final $core.List<Format?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static Format? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Format._(super.value, super.name);
}

class FrameOptionsList extends $pb.ProtobufEnum {
  static const FrameOptionsList FRAME_OPTIONS_LIST_SAMEORIGIN =
      FrameOptionsList._(
          0, _omitEnumNames ? '' : 'FRAME_OPTIONS_LIST_SAMEORIGIN');
  static const FrameOptionsList FRAME_OPTIONS_LIST_DENY =
      FrameOptionsList._(1, _omitEnumNames ? '' : 'FRAME_OPTIONS_LIST_DENY');

  static const $core.List<FrameOptionsList> values = <FrameOptionsList>[
    FRAME_OPTIONS_LIST_SAMEORIGIN,
    FRAME_OPTIONS_LIST_DENY,
  ];

  static final $core.List<FrameOptionsList?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FrameOptionsList? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FrameOptionsList._(super.value, super.name);
}

class FunctionRuntime extends $pb.ProtobufEnum {
  static const FunctionRuntime FUNCTION_RUNTIME_CLOUDFRONT_JS_1_0 =
      FunctionRuntime._(
          0, _omitEnumNames ? '' : 'FUNCTION_RUNTIME_CLOUDFRONT_JS_1_0');
  static const FunctionRuntime FUNCTION_RUNTIME_CLOUDFRONT_JS_2_0 =
      FunctionRuntime._(
          1, _omitEnumNames ? '' : 'FUNCTION_RUNTIME_CLOUDFRONT_JS_2_0');

  static const $core.List<FunctionRuntime> values = <FunctionRuntime>[
    FUNCTION_RUNTIME_CLOUDFRONT_JS_1_0,
    FUNCTION_RUNTIME_CLOUDFRONT_JS_2_0,
  ];

  static final $core.List<FunctionRuntime?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FunctionRuntime? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FunctionRuntime._(super.value, super.name);
}

class FunctionStage extends $pb.ProtobufEnum {
  static const FunctionStage FUNCTION_STAGE_DEVELOPMENT =
      FunctionStage._(0, _omitEnumNames ? '' : 'FUNCTION_STAGE_DEVELOPMENT');
  static const FunctionStage FUNCTION_STAGE_LIVE =
      FunctionStage._(1, _omitEnumNames ? '' : 'FUNCTION_STAGE_LIVE');

  static const $core.List<FunctionStage> values = <FunctionStage>[
    FUNCTION_STAGE_DEVELOPMENT,
    FUNCTION_STAGE_LIVE,
  ];

  static final $core.List<FunctionStage?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FunctionStage? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FunctionStage._(super.value, super.name);
}

class GeoRestrictionType extends $pb.ProtobufEnum {
  static const GeoRestrictionType GEO_RESTRICTION_TYPE_BLACKLIST =
      GeoRestrictionType._(
          0, _omitEnumNames ? '' : 'GEO_RESTRICTION_TYPE_BLACKLIST');
  static const GeoRestrictionType GEO_RESTRICTION_TYPE_WHITELIST =
      GeoRestrictionType._(
          1, _omitEnumNames ? '' : 'GEO_RESTRICTION_TYPE_WHITELIST');
  static const GeoRestrictionType GEO_RESTRICTION_TYPE_NONE =
      GeoRestrictionType._(
          2, _omitEnumNames ? '' : 'GEO_RESTRICTION_TYPE_NONE');

  static const $core.List<GeoRestrictionType> values = <GeoRestrictionType>[
    GEO_RESTRICTION_TYPE_BLACKLIST,
    GEO_RESTRICTION_TYPE_WHITELIST,
    GEO_RESTRICTION_TYPE_NONE,
  ];

  static final $core.List<GeoRestrictionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static GeoRestrictionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GeoRestrictionType._(super.value, super.name);
}

class HttpVersion extends $pb.ProtobufEnum {
  static const HttpVersion HTTP_VERSION_HTTP2 =
      HttpVersion._(0, _omitEnumNames ? '' : 'HTTP_VERSION_HTTP2');
  static const HttpVersion HTTP_VERSION_HTTP3 =
      HttpVersion._(1, _omitEnumNames ? '' : 'HTTP_VERSION_HTTP3');
  static const HttpVersion HTTP_VERSION_HTTP1_1 =
      HttpVersion._(2, _omitEnumNames ? '' : 'HTTP_VERSION_HTTP1_1');
  static const HttpVersion HTTP_VERSION_HTTP2AND3 =
      HttpVersion._(3, _omitEnumNames ? '' : 'HTTP_VERSION_HTTP2AND3');

  static const $core.List<HttpVersion> values = <HttpVersion>[
    HTTP_VERSION_HTTP2,
    HTTP_VERSION_HTTP3,
    HTTP_VERSION_HTTP1_1,
    HTTP_VERSION_HTTP2AND3,
  ];

  static final $core.List<HttpVersion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static HttpVersion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HttpVersion._(super.value, super.name);
}

class ICPRecordalStatus extends $pb.ProtobufEnum {
  static const ICPRecordalStatus I_C_P_RECORDAL_STATUS_PENDING =
      ICPRecordalStatus._(
          0, _omitEnumNames ? '' : 'I_C_P_RECORDAL_STATUS_PENDING');
  static const ICPRecordalStatus I_C_P_RECORDAL_STATUS_SUSPENDED =
      ICPRecordalStatus._(
          1, _omitEnumNames ? '' : 'I_C_P_RECORDAL_STATUS_SUSPENDED');
  static const ICPRecordalStatus I_C_P_RECORDAL_STATUS_APPROVED =
      ICPRecordalStatus._(
          2, _omitEnumNames ? '' : 'I_C_P_RECORDAL_STATUS_APPROVED');

  static const $core.List<ICPRecordalStatus> values = <ICPRecordalStatus>[
    I_C_P_RECORDAL_STATUS_PENDING,
    I_C_P_RECORDAL_STATUS_SUSPENDED,
    I_C_P_RECORDAL_STATUS_APPROVED,
  ];

  static final $core.List<ICPRecordalStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ICPRecordalStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ICPRecordalStatus._(super.value, super.name);
}

class ImportSourceType extends $pb.ProtobufEnum {
  static const ImportSourceType IMPORT_SOURCE_TYPE_S3 =
      ImportSourceType._(0, _omitEnumNames ? '' : 'IMPORT_SOURCE_TYPE_S3');

  static const $core.List<ImportSourceType> values = <ImportSourceType>[
    IMPORT_SOURCE_TYPE_S3,
  ];

  static final $core.List<ImportSourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ImportSourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportSourceType._(super.value, super.name);
}

class IpAddressType extends $pb.ProtobufEnum {
  static const IpAddressType IP_ADDRESS_TYPE_IPV6 =
      IpAddressType._(0, _omitEnumNames ? '' : 'IP_ADDRESS_TYPE_IPV6');
  static const IpAddressType IP_ADDRESS_TYPE_DUALSTACK =
      IpAddressType._(1, _omitEnumNames ? '' : 'IP_ADDRESS_TYPE_DUALSTACK');
  static const IpAddressType IP_ADDRESS_TYPE_IPV4 =
      IpAddressType._(2, _omitEnumNames ? '' : 'IP_ADDRESS_TYPE_IPV4');

  static const $core.List<IpAddressType> values = <IpAddressType>[
    IP_ADDRESS_TYPE_IPV6,
    IP_ADDRESS_TYPE_DUALSTACK,
    IP_ADDRESS_TYPE_IPV4,
  ];

  static final $core.List<IpAddressType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static IpAddressType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IpAddressType._(super.value, super.name);
}

class IpamCidrStatus extends $pb.ProtobufEnum {
  static const IpamCidrStatus IPAM_CIDR_STATUS_PROVISIONED =
      IpamCidrStatus._(0, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_PROVISIONED');
  static const IpamCidrStatus IPAM_CIDR_STATUS_FAILEDDEPROVISION =
      IpamCidrStatus._(
          1, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_FAILEDDEPROVISION');
  static const IpamCidrStatus IPAM_CIDR_STATUS_FAILEDPROVISION =
      IpamCidrStatus._(
          2, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_FAILEDPROVISION');
  static const IpamCidrStatus IPAM_CIDR_STATUS_ADVERTISED =
      IpamCidrStatus._(3, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_ADVERTISED');
  static const IpamCidrStatus IPAM_CIDR_STATUS_DEPROVISIONED = IpamCidrStatus._(
      4, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_DEPROVISIONED');
  static const IpamCidrStatus IPAM_CIDR_STATUS_WITHDRAWN =
      IpamCidrStatus._(5, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_WITHDRAWN');
  static const IpamCidrStatus IPAM_CIDR_STATUS_PROVISIONING = IpamCidrStatus._(
      6, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_PROVISIONING');
  static const IpamCidrStatus IPAM_CIDR_STATUS_FAILEDADVERTISE =
      IpamCidrStatus._(
          7, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_FAILEDADVERTISE');
  static const IpamCidrStatus IPAM_CIDR_STATUS_FAILEDWITHDRAW =
      IpamCidrStatus._(
          8, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_FAILEDWITHDRAW');
  static const IpamCidrStatus IPAM_CIDR_STATUS_DEPROVISIONING =
      IpamCidrStatus._(
          9, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_DEPROVISIONING');
  static const IpamCidrStatus IPAM_CIDR_STATUS_ADVERTISING = IpamCidrStatus._(
      10, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_ADVERTISING');
  static const IpamCidrStatus IPAM_CIDR_STATUS_WITHDRAWING = IpamCidrStatus._(
      11, _omitEnumNames ? '' : 'IPAM_CIDR_STATUS_WITHDRAWING');

  static const $core.List<IpamCidrStatus> values = <IpamCidrStatus>[
    IPAM_CIDR_STATUS_PROVISIONED,
    IPAM_CIDR_STATUS_FAILEDDEPROVISION,
    IPAM_CIDR_STATUS_FAILEDPROVISION,
    IPAM_CIDR_STATUS_ADVERTISED,
    IPAM_CIDR_STATUS_DEPROVISIONED,
    IPAM_CIDR_STATUS_WITHDRAWN,
    IPAM_CIDR_STATUS_PROVISIONING,
    IPAM_CIDR_STATUS_FAILEDADVERTISE,
    IPAM_CIDR_STATUS_FAILEDWITHDRAW,
    IPAM_CIDR_STATUS_DEPROVISIONING,
    IPAM_CIDR_STATUS_ADVERTISING,
    IPAM_CIDR_STATUS_WITHDRAWING,
  ];

  static final $core.List<IpamCidrStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 11);
  static IpamCidrStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IpamCidrStatus._(super.value, super.name);
}

class ItemSelection extends $pb.ProtobufEnum {
  static const ItemSelection ITEM_SELECTION_ALL =
      ItemSelection._(0, _omitEnumNames ? '' : 'ITEM_SELECTION_ALL');
  static const ItemSelection ITEM_SELECTION_WHITELIST =
      ItemSelection._(1, _omitEnumNames ? '' : 'ITEM_SELECTION_WHITELIST');
  static const ItemSelection ITEM_SELECTION_NONE =
      ItemSelection._(2, _omitEnumNames ? '' : 'ITEM_SELECTION_NONE');

  static const $core.List<ItemSelection> values = <ItemSelection>[
    ITEM_SELECTION_ALL,
    ITEM_SELECTION_WHITELIST,
    ITEM_SELECTION_NONE,
  ];

  static final $core.List<ItemSelection?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ItemSelection? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ItemSelection._(super.value, super.name);
}

class ManagedCertificateStatus extends $pb.ProtobufEnum {
  static const ManagedCertificateStatus
      MANAGED_CERTIFICATE_STATUS_PENDINGVALIDATION = ManagedCertificateStatus._(
          0,
          _omitEnumNames ? '' : 'MANAGED_CERTIFICATE_STATUS_PENDINGVALIDATION');
  static const ManagedCertificateStatus MANAGED_CERTIFICATE_STATUS_REVOKED =
      ManagedCertificateStatus._(
          1, _omitEnumNames ? '' : 'MANAGED_CERTIFICATE_STATUS_REVOKED');
  static const ManagedCertificateStatus
      MANAGED_CERTIFICATE_STATUS_VALIDATIONTIMEDOUT =
      ManagedCertificateStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'MANAGED_CERTIFICATE_STATUS_VALIDATIONTIMEDOUT');
  static const ManagedCertificateStatus MANAGED_CERTIFICATE_STATUS_FAILED =
      ManagedCertificateStatus._(
          3, _omitEnumNames ? '' : 'MANAGED_CERTIFICATE_STATUS_FAILED');
  static const ManagedCertificateStatus MANAGED_CERTIFICATE_STATUS_EXPIRED =
      ManagedCertificateStatus._(
          4, _omitEnumNames ? '' : 'MANAGED_CERTIFICATE_STATUS_EXPIRED');
  static const ManagedCertificateStatus MANAGED_CERTIFICATE_STATUS_ISSUED =
      ManagedCertificateStatus._(
          5, _omitEnumNames ? '' : 'MANAGED_CERTIFICATE_STATUS_ISSUED');
  static const ManagedCertificateStatus MANAGED_CERTIFICATE_STATUS_INACTIVE =
      ManagedCertificateStatus._(
          6, _omitEnumNames ? '' : 'MANAGED_CERTIFICATE_STATUS_INACTIVE');

  static const $core.List<ManagedCertificateStatus> values =
      <ManagedCertificateStatus>[
    MANAGED_CERTIFICATE_STATUS_PENDINGVALIDATION,
    MANAGED_CERTIFICATE_STATUS_REVOKED,
    MANAGED_CERTIFICATE_STATUS_VALIDATIONTIMEDOUT,
    MANAGED_CERTIFICATE_STATUS_FAILED,
    MANAGED_CERTIFICATE_STATUS_EXPIRED,
    MANAGED_CERTIFICATE_STATUS_ISSUED,
    MANAGED_CERTIFICATE_STATUS_INACTIVE,
  ];

  static final $core.List<ManagedCertificateStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static ManagedCertificateStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ManagedCertificateStatus._(super.value, super.name);
}

class Method extends $pb.ProtobufEnum {
  static const Method METHOD_OPTIONS =
      Method._(0, _omitEnumNames ? '' : 'METHOD_OPTIONS');
  static const Method METHOD_PATCH =
      Method._(1, _omitEnumNames ? '' : 'METHOD_PATCH');
  static const Method METHOD_POST =
      Method._(2, _omitEnumNames ? '' : 'METHOD_POST');
  static const Method METHOD_HEAD =
      Method._(3, _omitEnumNames ? '' : 'METHOD_HEAD');
  static const Method METHOD_GET =
      Method._(4, _omitEnumNames ? '' : 'METHOD_GET');
  static const Method METHOD_DELETE =
      Method._(5, _omitEnumNames ? '' : 'METHOD_DELETE');
  static const Method METHOD_PUT =
      Method._(6, _omitEnumNames ? '' : 'METHOD_PUT');

  static const $core.List<Method> values = <Method>[
    METHOD_OPTIONS,
    METHOD_PATCH,
    METHOD_POST,
    METHOD_HEAD,
    METHOD_GET,
    METHOD_DELETE,
    METHOD_PUT,
  ];

  static final $core.List<Method?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static Method? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Method._(super.value, super.name);
}

class MinimumProtocolVersion extends $pb.ProtobufEnum {
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1_1_2016 =
      MinimumProtocolVersion._(
          0, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1_1_2016');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_SSLV3 =
      MinimumProtocolVersion._(
          1, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_SSLV3');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1_2016 =
      MinimumProtocolVersion._(
          2, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1_2016');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1_2_2019 =
      MinimumProtocolVersion._(
          3, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2019');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1_3_2025 =
      MinimumProtocolVersion._(
          4, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1_3_2025');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1_2_2021 =
      MinimumProtocolVersion._(
          5, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2021');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1 =
      MinimumProtocolVersion._(
          6, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1_2_2018 =
      MinimumProtocolVersion._(
          7, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2018');
  static const MinimumProtocolVersion MINIMUM_PROTOCOL_VERSION_TLSV1_2_2025 =
      MinimumProtocolVersion._(
          8, _omitEnumNames ? '' : 'MINIMUM_PROTOCOL_VERSION_TLSV1_2_2025');

  static const $core.List<MinimumProtocolVersion> values =
      <MinimumProtocolVersion>[
    MINIMUM_PROTOCOL_VERSION_TLSV1_1_2016,
    MINIMUM_PROTOCOL_VERSION_SSLV3,
    MINIMUM_PROTOCOL_VERSION_TLSV1_2016,
    MINIMUM_PROTOCOL_VERSION_TLSV1_2_2019,
    MINIMUM_PROTOCOL_VERSION_TLSV1_3_2025,
    MINIMUM_PROTOCOL_VERSION_TLSV1_2_2021,
    MINIMUM_PROTOCOL_VERSION_TLSV1,
    MINIMUM_PROTOCOL_VERSION_TLSV1_2_2018,
    MINIMUM_PROTOCOL_VERSION_TLSV1_2_2025,
  ];

  static final $core.List<MinimumProtocolVersion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static MinimumProtocolVersion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MinimumProtocolVersion._(super.value, super.name);
}

class OriginAccessControlOriginTypes extends $pb.ProtobufEnum {
  static const OriginAccessControlOriginTypes
      ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_S3 = OriginAccessControlOriginTypes._(
          0, _omitEnumNames ? '' : 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_S3');
  static const OriginAccessControlOriginTypes
      ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_LAMBDA =
      OriginAccessControlOriginTypes._(
          1, _omitEnumNames ? '' : 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_LAMBDA');
  static const OriginAccessControlOriginTypes
      ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIASTORE =
      OriginAccessControlOriginTypes._(
          2,
          _omitEnumNames
              ? ''
              : 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIASTORE');
  static const OriginAccessControlOriginTypes
      ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIAPACKAGEV2 =
      OriginAccessControlOriginTypes._(
          3,
          _omitEnumNames
              ? ''
              : 'ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIAPACKAGEV2');

  static const $core.List<OriginAccessControlOriginTypes> values =
      <OriginAccessControlOriginTypes>[
    ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_S3,
    ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_LAMBDA,
    ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIASTORE,
    ORIGIN_ACCESS_CONTROL_ORIGIN_TYPES_MEDIAPACKAGEV2,
  ];

  static final $core.List<OriginAccessControlOriginTypes?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static OriginAccessControlOriginTypes? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginAccessControlOriginTypes._(super.value, super.name);
}

class OriginAccessControlSigningBehaviors extends $pb.ProtobufEnum {
  static const OriginAccessControlSigningBehaviors
      ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_ALWAYS =
      OriginAccessControlSigningBehaviors._(
          0,
          _omitEnumNames
              ? ''
              : 'ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_ALWAYS');
  static const OriginAccessControlSigningBehaviors
      ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NO_OVERRIDE =
      OriginAccessControlSigningBehaviors._(
          1,
          _omitEnumNames
              ? ''
              : 'ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NO_OVERRIDE');
  static const OriginAccessControlSigningBehaviors
      ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NEVER =
      OriginAccessControlSigningBehaviors._(
          2,
          _omitEnumNames
              ? ''
              : 'ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NEVER');

  static const $core.List<OriginAccessControlSigningBehaviors> values =
      <OriginAccessControlSigningBehaviors>[
    ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_ALWAYS,
    ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NO_OVERRIDE,
    ORIGIN_ACCESS_CONTROL_SIGNING_BEHAVIORS_NEVER,
  ];

  static final $core.List<OriginAccessControlSigningBehaviors?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static OriginAccessControlSigningBehaviors? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginAccessControlSigningBehaviors._(super.value, super.name);
}

class OriginAccessControlSigningProtocols extends $pb.ProtobufEnum {
  static const OriginAccessControlSigningProtocols
      ORIGIN_ACCESS_CONTROL_SIGNING_PROTOCOLS_SIGV4 =
      OriginAccessControlSigningProtocols._(
          0,
          _omitEnumNames
              ? ''
              : 'ORIGIN_ACCESS_CONTROL_SIGNING_PROTOCOLS_SIGV4');

  static const $core.List<OriginAccessControlSigningProtocols> values =
      <OriginAccessControlSigningProtocols>[
    ORIGIN_ACCESS_CONTROL_SIGNING_PROTOCOLS_SIGV4,
  ];

  static final $core.List<OriginAccessControlSigningProtocols?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static OriginAccessControlSigningProtocols? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginAccessControlSigningProtocols._(super.value, super.name);
}

class OriginGroupSelectionCriteria extends $pb.ProtobufEnum {
  static const OriginGroupSelectionCriteria
      ORIGIN_GROUP_SELECTION_CRITERIA_MEDIAQUALITYBASED =
      OriginGroupSelectionCriteria._(
          0,
          _omitEnumNames
              ? ''
              : 'ORIGIN_GROUP_SELECTION_CRITERIA_MEDIAQUALITYBASED');
  static const OriginGroupSelectionCriteria
      ORIGIN_GROUP_SELECTION_CRITERIA_DEFAULT = OriginGroupSelectionCriteria._(
          1, _omitEnumNames ? '' : 'ORIGIN_GROUP_SELECTION_CRITERIA_DEFAULT');

  static const $core.List<OriginGroupSelectionCriteria> values =
      <OriginGroupSelectionCriteria>[
    ORIGIN_GROUP_SELECTION_CRITERIA_MEDIAQUALITYBASED,
    ORIGIN_GROUP_SELECTION_CRITERIA_DEFAULT,
  ];

  static final $core.List<OriginGroupSelectionCriteria?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static OriginGroupSelectionCriteria? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginGroupSelectionCriteria._(super.value, super.name);
}

class OriginProtocolPolicy extends $pb.ProtobufEnum {
  static const OriginProtocolPolicy ORIGIN_PROTOCOL_POLICY_HTTP_ONLY =
      OriginProtocolPolicy._(
          0, _omitEnumNames ? '' : 'ORIGIN_PROTOCOL_POLICY_HTTP_ONLY');
  static const OriginProtocolPolicy ORIGIN_PROTOCOL_POLICY_MATCH_VIEWER =
      OriginProtocolPolicy._(
          1, _omitEnumNames ? '' : 'ORIGIN_PROTOCOL_POLICY_MATCH_VIEWER');
  static const OriginProtocolPolicy ORIGIN_PROTOCOL_POLICY_HTTPS_ONLY =
      OriginProtocolPolicy._(
          2, _omitEnumNames ? '' : 'ORIGIN_PROTOCOL_POLICY_HTTPS_ONLY');

  static const $core.List<OriginProtocolPolicy> values = <OriginProtocolPolicy>[
    ORIGIN_PROTOCOL_POLICY_HTTP_ONLY,
    ORIGIN_PROTOCOL_POLICY_MATCH_VIEWER,
    ORIGIN_PROTOCOL_POLICY_HTTPS_ONLY,
  ];

  static final $core.List<OriginProtocolPolicy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static OriginProtocolPolicy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginProtocolPolicy._(super.value, super.name);
}

class OriginRequestPolicyCookieBehavior extends $pb.ProtobufEnum {
  static const OriginRequestPolicyCookieBehavior
      ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALL =
      OriginRequestPolicyCookieBehavior._(
          0, _omitEnumNames ? '' : 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALL');
  static const OriginRequestPolicyCookieBehavior
      ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_WHITELIST =
      OriginRequestPolicyCookieBehavior._(
          1,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_WHITELIST');
  static const OriginRequestPolicyCookieBehavior
      ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_NONE =
      OriginRequestPolicyCookieBehavior._(2,
          _omitEnumNames ? '' : 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_NONE');
  static const OriginRequestPolicyCookieBehavior
      ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT =
      OriginRequestPolicyCookieBehavior._(
          3,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT');

  static const $core.List<OriginRequestPolicyCookieBehavior> values =
      <OriginRequestPolicyCookieBehavior>[
    ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALL,
    ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_WHITELIST,
    ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_NONE,
    ORIGIN_REQUEST_POLICY_COOKIE_BEHAVIOR_ALLEXCEPT,
  ];

  static final $core.List<OriginRequestPolicyCookieBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static OriginRequestPolicyCookieBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginRequestPolicyCookieBehavior._(super.value, super.name);
}

class OriginRequestPolicyHeaderBehavior extends $pb.ProtobufEnum {
  static const OriginRequestPolicyHeaderBehavior
      ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWERANDWHITELISTCLOUDFRONT =
      OriginRequestPolicyHeaderBehavior._(
          0,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWERANDWHITELISTCLOUDFRONT');
  static const OriginRequestPolicyHeaderBehavior
      ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_WHITELIST =
      OriginRequestPolicyHeaderBehavior._(
          1,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_WHITELIST');
  static const OriginRequestPolicyHeaderBehavior
      ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_NONE =
      OriginRequestPolicyHeaderBehavior._(2,
          _omitEnumNames ? '' : 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_NONE');
  static const OriginRequestPolicyHeaderBehavior
      ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWER =
      OriginRequestPolicyHeaderBehavior._(
          3,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWER');
  static const OriginRequestPolicyHeaderBehavior
      ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLEXCEPT =
      OriginRequestPolicyHeaderBehavior._(
          4,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLEXCEPT');

  static const $core.List<OriginRequestPolicyHeaderBehavior> values =
      <OriginRequestPolicyHeaderBehavior>[
    ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWERANDWHITELISTCLOUDFRONT,
    ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_WHITELIST,
    ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_NONE,
    ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLVIEWER,
    ORIGIN_REQUEST_POLICY_HEADER_BEHAVIOR_ALLEXCEPT,
  ];

  static final $core.List<OriginRequestPolicyHeaderBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static OriginRequestPolicyHeaderBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginRequestPolicyHeaderBehavior._(super.value, super.name);
}

class OriginRequestPolicyQueryStringBehavior extends $pb.ProtobufEnum {
  static const OriginRequestPolicyQueryStringBehavior
      ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALL =
      OriginRequestPolicyQueryStringBehavior._(
          0,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALL');
  static const OriginRequestPolicyQueryStringBehavior
      ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST =
      OriginRequestPolicyQueryStringBehavior._(
          1,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST');
  static const OriginRequestPolicyQueryStringBehavior
      ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_NONE =
      OriginRequestPolicyQueryStringBehavior._(
          2,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_NONE');
  static const OriginRequestPolicyQueryStringBehavior
      ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT =
      OriginRequestPolicyQueryStringBehavior._(
          3,
          _omitEnumNames
              ? ''
              : 'ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT');

  static const $core.List<OriginRequestPolicyQueryStringBehavior> values =
      <OriginRequestPolicyQueryStringBehavior>[
    ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALL,
    ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_WHITELIST,
    ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_NONE,
    ORIGIN_REQUEST_POLICY_QUERY_STRING_BEHAVIOR_ALLEXCEPT,
  ];

  static final $core.List<OriginRequestPolicyQueryStringBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static OriginRequestPolicyQueryStringBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginRequestPolicyQueryStringBehavior._(super.value, super.name);
}

class OriginRequestPolicyType extends $pb.ProtobufEnum {
  static const OriginRequestPolicyType ORIGIN_REQUEST_POLICY_TYPE_MANAGED =
      OriginRequestPolicyType._(
          0, _omitEnumNames ? '' : 'ORIGIN_REQUEST_POLICY_TYPE_MANAGED');
  static const OriginRequestPolicyType ORIGIN_REQUEST_POLICY_TYPE_CUSTOM =
      OriginRequestPolicyType._(
          1, _omitEnumNames ? '' : 'ORIGIN_REQUEST_POLICY_TYPE_CUSTOM');

  static const $core.List<OriginRequestPolicyType> values =
      <OriginRequestPolicyType>[
    ORIGIN_REQUEST_POLICY_TYPE_MANAGED,
    ORIGIN_REQUEST_POLICY_TYPE_CUSTOM,
  ];

  static final $core.List<OriginRequestPolicyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static OriginRequestPolicyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginRequestPolicyType._(super.value, super.name);
}

class PriceClass extends $pb.ProtobufEnum {
  static const PriceClass PRICE_CLASS_NONE =
      PriceClass._(0, _omitEnumNames ? '' : 'PRICE_CLASS_NONE');
  static const PriceClass PRICE_CLASS_PRICECLASS_200 =
      PriceClass._(1, _omitEnumNames ? '' : 'PRICE_CLASS_PRICECLASS_200');
  static const PriceClass PRICE_CLASS_PRICECLASS_100 =
      PriceClass._(2, _omitEnumNames ? '' : 'PRICE_CLASS_PRICECLASS_100');
  static const PriceClass PRICE_CLASS_PRICECLASS_ALL =
      PriceClass._(3, _omitEnumNames ? '' : 'PRICE_CLASS_PRICECLASS_ALL');

  static const $core.List<PriceClass> values = <PriceClass>[
    PRICE_CLASS_NONE,
    PRICE_CLASS_PRICECLASS_200,
    PRICE_CLASS_PRICECLASS_100,
    PRICE_CLASS_PRICECLASS_ALL,
  ];

  static final $core.List<PriceClass?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static PriceClass? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PriceClass._(super.value, super.name);
}

class RealtimeMetricsSubscriptionStatus extends $pb.ProtobufEnum {
  static const RealtimeMetricsSubscriptionStatus
      REALTIME_METRICS_SUBSCRIPTION_STATUS_DISABLED =
      RealtimeMetricsSubscriptionStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'REALTIME_METRICS_SUBSCRIPTION_STATUS_DISABLED');
  static const RealtimeMetricsSubscriptionStatus
      REALTIME_METRICS_SUBSCRIPTION_STATUS_ENABLED =
      RealtimeMetricsSubscriptionStatus._(1,
          _omitEnumNames ? '' : 'REALTIME_METRICS_SUBSCRIPTION_STATUS_ENABLED');

  static const $core.List<RealtimeMetricsSubscriptionStatus> values =
      <RealtimeMetricsSubscriptionStatus>[
    REALTIME_METRICS_SUBSCRIPTION_STATUS_DISABLED,
    REALTIME_METRICS_SUBSCRIPTION_STATUS_ENABLED,
  ];

  static final $core.List<RealtimeMetricsSubscriptionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RealtimeMetricsSubscriptionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RealtimeMetricsSubscriptionStatus._(super.value, super.name);
}

class ReferrerPolicyList extends $pb.ProtobufEnum {
  static const ReferrerPolicyList REFERRER_POLICY_LIST_UNSAFE_URL =
      ReferrerPolicyList._(
          0, _omitEnumNames ? '' : 'REFERRER_POLICY_LIST_UNSAFE_URL');
  static const ReferrerPolicyList
      REFERRER_POLICY_LIST_STRICT_ORIGIN_WHEN_CROSS_ORIGIN =
      ReferrerPolicyList._(
          1,
          _omitEnumNames
              ? ''
              : 'REFERRER_POLICY_LIST_STRICT_ORIGIN_WHEN_CROSS_ORIGIN');
  static const ReferrerPolicyList REFERRER_POLICY_LIST_NO_REFERRER =
      ReferrerPolicyList._(
          2, _omitEnumNames ? '' : 'REFERRER_POLICY_LIST_NO_REFERRER');
  static const ReferrerPolicyList
      REFERRER_POLICY_LIST_NO_REFERRER_WHEN_DOWNGRADE = ReferrerPolicyList._(
          3,
          _omitEnumNames
              ? ''
              : 'REFERRER_POLICY_LIST_NO_REFERRER_WHEN_DOWNGRADE');
  static const ReferrerPolicyList REFERRER_POLICY_LIST_ORIGIN =
      ReferrerPolicyList._(
          4, _omitEnumNames ? '' : 'REFERRER_POLICY_LIST_ORIGIN');
  static const ReferrerPolicyList REFERRER_POLICY_LIST_STRICT_ORIGIN =
      ReferrerPolicyList._(
          5, _omitEnumNames ? '' : 'REFERRER_POLICY_LIST_STRICT_ORIGIN');
  static const ReferrerPolicyList
      REFERRER_POLICY_LIST_ORIGIN_WHEN_CROSS_ORIGIN = ReferrerPolicyList._(
          6,
          _omitEnumNames
              ? ''
              : 'REFERRER_POLICY_LIST_ORIGIN_WHEN_CROSS_ORIGIN');
  static const ReferrerPolicyList REFERRER_POLICY_LIST_SAME_ORIGIN =
      ReferrerPolicyList._(
          7, _omitEnumNames ? '' : 'REFERRER_POLICY_LIST_SAME_ORIGIN');

  static const $core.List<ReferrerPolicyList> values = <ReferrerPolicyList>[
    REFERRER_POLICY_LIST_UNSAFE_URL,
    REFERRER_POLICY_LIST_STRICT_ORIGIN_WHEN_CROSS_ORIGIN,
    REFERRER_POLICY_LIST_NO_REFERRER,
    REFERRER_POLICY_LIST_NO_REFERRER_WHEN_DOWNGRADE,
    REFERRER_POLICY_LIST_ORIGIN,
    REFERRER_POLICY_LIST_STRICT_ORIGIN,
    REFERRER_POLICY_LIST_ORIGIN_WHEN_CROSS_ORIGIN,
    REFERRER_POLICY_LIST_SAME_ORIGIN,
  ];

  static final $core.List<ReferrerPolicyList?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static ReferrerPolicyList? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReferrerPolicyList._(super.value, super.name);
}

class ResponseHeadersPolicyAccessControlAllowMethodsValues
    extends $pb.ProtobufEnum {
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_OPTIONS =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          0,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_OPTIONS');
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PATCH =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          1,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PATCH');
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_POST =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          2,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_POST');
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_HEAD =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          3,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_HEAD');
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_GET =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          4,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_GET');
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_DELETE =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          5,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_DELETE');
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_ALL =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          6,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_ALL');
  static const ResponseHeadersPolicyAccessControlAllowMethodsValues
      RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PUT =
      ResponseHeadersPolicyAccessControlAllowMethodsValues._(
          7,
          _omitEnumNames
              ? ''
              : 'RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PUT');

  static const $core.List<ResponseHeadersPolicyAccessControlAllowMethodsValues>
      values = <ResponseHeadersPolicyAccessControlAllowMethodsValues>[
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_OPTIONS,
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PATCH,
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_POST,
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_HEAD,
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_GET,
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_DELETE,
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_ALL,
    RESPONSE_HEADERS_POLICY_ACCESS_CONTROL_ALLOW_METHODS_VALUES_PUT,
  ];

  static final $core.List<ResponseHeadersPolicyAccessControlAllowMethodsValues?>
      _byValue = $pb.ProtobufEnum.$_initByValueList(values, 7);
  static ResponseHeadersPolicyAccessControlAllowMethodsValues? valueOf(
          $core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResponseHeadersPolicyAccessControlAllowMethodsValues._(
      super.value, super.name);
}

class ResponseHeadersPolicyType extends $pb.ProtobufEnum {
  static const ResponseHeadersPolicyType RESPONSE_HEADERS_POLICY_TYPE_MANAGED =
      ResponseHeadersPolicyType._(
          0, _omitEnumNames ? '' : 'RESPONSE_HEADERS_POLICY_TYPE_MANAGED');
  static const ResponseHeadersPolicyType RESPONSE_HEADERS_POLICY_TYPE_CUSTOM =
      ResponseHeadersPolicyType._(
          1, _omitEnumNames ? '' : 'RESPONSE_HEADERS_POLICY_TYPE_CUSTOM');

  static const $core.List<ResponseHeadersPolicyType> values =
      <ResponseHeadersPolicyType>[
    RESPONSE_HEADERS_POLICY_TYPE_MANAGED,
    RESPONSE_HEADERS_POLICY_TYPE_CUSTOM,
  ];

  static final $core.List<ResponseHeadersPolicyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ResponseHeadersPolicyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResponseHeadersPolicyType._(super.value, super.name);
}

class SSLSupportMethod extends $pb.ProtobufEnum {
  static const SSLSupportMethod S_S_L_SUPPORT_METHOD_STATIC_IP =
      SSLSupportMethod._(
          0, _omitEnumNames ? '' : 'S_S_L_SUPPORT_METHOD_STATIC_IP');
  static const SSLSupportMethod S_S_L_SUPPORT_METHOD_VIP =
      SSLSupportMethod._(1, _omitEnumNames ? '' : 'S_S_L_SUPPORT_METHOD_VIP');
  static const SSLSupportMethod S_S_L_SUPPORT_METHOD_SNI_ONLY =
      SSLSupportMethod._(
          2, _omitEnumNames ? '' : 'S_S_L_SUPPORT_METHOD_SNI_ONLY');

  static const $core.List<SSLSupportMethod> values = <SSLSupportMethod>[
    S_S_L_SUPPORT_METHOD_STATIC_IP,
    S_S_L_SUPPORT_METHOD_VIP,
    S_S_L_SUPPORT_METHOD_SNI_ONLY,
  ];

  static final $core.List<SSLSupportMethod?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SSLSupportMethod? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SSLSupportMethod._(super.value, super.name);
}

class SslProtocol extends $pb.ProtobufEnum {
  static const SslProtocol SSL_PROTOCOL_SSLV3 =
      SslProtocol._(0, _omitEnumNames ? '' : 'SSL_PROTOCOL_SSLV3');
  static const SslProtocol SSL_PROTOCOL_TLSV1_1 =
      SslProtocol._(1, _omitEnumNames ? '' : 'SSL_PROTOCOL_TLSV1_1');
  static const SslProtocol SSL_PROTOCOL_TLSV1 =
      SslProtocol._(2, _omitEnumNames ? '' : 'SSL_PROTOCOL_TLSV1');
  static const SslProtocol SSL_PROTOCOL_TLSV1_2 =
      SslProtocol._(3, _omitEnumNames ? '' : 'SSL_PROTOCOL_TLSV1_2');

  static const $core.List<SslProtocol> values = <SslProtocol>[
    SSL_PROTOCOL_SSLV3,
    SSL_PROTOCOL_TLSV1_1,
    SSL_PROTOCOL_TLSV1,
    SSL_PROTOCOL_TLSV1_2,
  ];

  static final $core.List<SslProtocol?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static SslProtocol? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SslProtocol._(super.value, super.name);
}

class TrustStoreStatus extends $pb.ProtobufEnum {
  static const TrustStoreStatus TRUST_STORE_STATUS_ACTIVE =
      TrustStoreStatus._(0, _omitEnumNames ? '' : 'TRUST_STORE_STATUS_ACTIVE');
  static const TrustStoreStatus TRUST_STORE_STATUS_FAILED =
      TrustStoreStatus._(1, _omitEnumNames ? '' : 'TRUST_STORE_STATUS_FAILED');
  static const TrustStoreStatus TRUST_STORE_STATUS_PENDING =
      TrustStoreStatus._(2, _omitEnumNames ? '' : 'TRUST_STORE_STATUS_PENDING');

  static const $core.List<TrustStoreStatus> values = <TrustStoreStatus>[
    TRUST_STORE_STATUS_ACTIVE,
    TRUST_STORE_STATUS_FAILED,
    TRUST_STORE_STATUS_PENDING,
  ];

  static final $core.List<TrustStoreStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static TrustStoreStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TrustStoreStatus._(super.value, super.name);
}

class ValidationTokenHost extends $pb.ProtobufEnum {
  static const ValidationTokenHost VALIDATION_TOKEN_HOST_SELFHOSTED =
      ValidationTokenHost._(
          0, _omitEnumNames ? '' : 'VALIDATION_TOKEN_HOST_SELFHOSTED');
  static const ValidationTokenHost VALIDATION_TOKEN_HOST_CLOUDFRONT =
      ValidationTokenHost._(
          1, _omitEnumNames ? '' : 'VALIDATION_TOKEN_HOST_CLOUDFRONT');

  static const $core.List<ValidationTokenHost> values = <ValidationTokenHost>[
    VALIDATION_TOKEN_HOST_SELFHOSTED,
    VALIDATION_TOKEN_HOST_CLOUDFRONT,
  ];

  static final $core.List<ValidationTokenHost?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ValidationTokenHost? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ValidationTokenHost._(super.value, super.name);
}

class ViewerMtlsMode extends $pb.ProtobufEnum {
  static const ViewerMtlsMode VIEWER_MTLS_MODE_OPTIONAL =
      ViewerMtlsMode._(0, _omitEnumNames ? '' : 'VIEWER_MTLS_MODE_OPTIONAL');
  static const ViewerMtlsMode VIEWER_MTLS_MODE_REQUIRED =
      ViewerMtlsMode._(1, _omitEnumNames ? '' : 'VIEWER_MTLS_MODE_REQUIRED');

  static const $core.List<ViewerMtlsMode> values = <ViewerMtlsMode>[
    VIEWER_MTLS_MODE_OPTIONAL,
    VIEWER_MTLS_MODE_REQUIRED,
  ];

  static final $core.List<ViewerMtlsMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ViewerMtlsMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ViewerMtlsMode._(super.value, super.name);
}

class ViewerProtocolPolicy extends $pb.ProtobufEnum {
  static const ViewerProtocolPolicy VIEWER_PROTOCOL_POLICY_ALLOW_ALL =
      ViewerProtocolPolicy._(
          0, _omitEnumNames ? '' : 'VIEWER_PROTOCOL_POLICY_ALLOW_ALL');
  static const ViewerProtocolPolicy VIEWER_PROTOCOL_POLICY_HTTPS_ONLY =
      ViewerProtocolPolicy._(
          1, _omitEnumNames ? '' : 'VIEWER_PROTOCOL_POLICY_HTTPS_ONLY');
  static const ViewerProtocolPolicy VIEWER_PROTOCOL_POLICY_REDIRECT_TO_HTTPS =
      ViewerProtocolPolicy._(
          2, _omitEnumNames ? '' : 'VIEWER_PROTOCOL_POLICY_REDIRECT_TO_HTTPS');

  static const $core.List<ViewerProtocolPolicy> values = <ViewerProtocolPolicy>[
    VIEWER_PROTOCOL_POLICY_ALLOW_ALL,
    VIEWER_PROTOCOL_POLICY_HTTPS_ONLY,
    VIEWER_PROTOCOL_POLICY_REDIRECT_TO_HTTPS,
  ];

  static final $core.List<ViewerProtocolPolicy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ViewerProtocolPolicy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ViewerProtocolPolicy._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
