// This is a generated file - do not edit.
//
// Generated from route53.proto.

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

@$core.Deprecated('Use acceleratedRecoveryStatusDescriptor instead')
const AcceleratedRecoveryStatus$json = {
  '1': 'AcceleratedRecoveryStatus',
  '2': [
    {'1': 'ACCELERATED_RECOVERY_STATUS_DISABLED', '2': 0},
    {'1': 'ACCELERATED_RECOVERY_STATUS_DISABLE_FAILED', '2': 1},
    {'1': 'ACCELERATED_RECOVERY_STATUS_DISABLING_HOSTED_ZONE_LOCKED', '2': 2},
    {'1': 'ACCELERATED_RECOVERY_STATUS_ENABLING', '2': 3},
    {'1': 'ACCELERATED_RECOVERY_STATUS_ENABLE_FAILED', '2': 4},
    {'1': 'ACCELERATED_RECOVERY_STATUS_ENABLING_HOSTED_ZONE_LOCKED', '2': 5},
    {'1': 'ACCELERATED_RECOVERY_STATUS_DISABLING', '2': 6},
    {'1': 'ACCELERATED_RECOVERY_STATUS_ENABLED', '2': 7},
  ],
};

/// Descriptor for `AcceleratedRecoveryStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List acceleratedRecoveryStatusDescriptor = $convert.base64Decode(
    'ChlBY2NlbGVyYXRlZFJlY292ZXJ5U3RhdHVzEigKJEFDQ0VMRVJBVEVEX1JFQ09WRVJZX1NUQV'
    'RVU19ESVNBQkxFRBAAEi4KKkFDQ0VMRVJBVEVEX1JFQ09WRVJZX1NUQVRVU19ESVNBQkxFX0ZB'
    'SUxFRBABEjwKOEFDQ0VMRVJBVEVEX1JFQ09WRVJZX1NUQVRVU19ESVNBQkxJTkdfSE9TVEVEX1'
    'pPTkVfTE9DS0VEEAISKAokQUNDRUxFUkFURURfUkVDT1ZFUllfU1RBVFVTX0VOQUJMSU5HEAMS'
    'LQopQUNDRUxFUkFURURfUkVDT1ZFUllfU1RBVFVTX0VOQUJMRV9GQUlMRUQQBBI7CjdBQ0NFTE'
    'VSQVRFRF9SRUNPVkVSWV9TVEFUVVNfRU5BQkxJTkdfSE9TVEVEX1pPTkVfTE9DS0VEEAUSKQol'
    'QUNDRUxFUkFURURfUkVDT1ZFUllfU1RBVFVTX0RJU0FCTElORxAGEicKI0FDQ0VMRVJBVEVEX1'
    'JFQ09WRVJZX1NUQVRVU19FTkFCTEVEEAc=');

@$core.Deprecated('Use accountLimitTypeDescriptor instead')
const AccountLimitType$json = {
  '1': 'AccountLimitType',
  '2': [
    {'1': 'ACCOUNT_LIMIT_TYPE_MAX_REUSABLE_DELEGATION_SETS_BY_OWNER', '2': 0},
    {'1': 'ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICY_INSTANCES_BY_OWNER', '2': 1},
    {'1': 'ACCOUNT_LIMIT_TYPE_MAX_HEALTH_CHECKS_BY_OWNER', '2': 2},
    {'1': 'ACCOUNT_LIMIT_TYPE_MAX_HOSTED_ZONES_BY_OWNER', '2': 3},
    {'1': 'ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICIES_BY_OWNER', '2': 4},
  ],
};

/// Descriptor for `AccountLimitType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List accountLimitTypeDescriptor = $convert.base64Decode(
    'ChBBY2NvdW50TGltaXRUeXBlEjwKOEFDQ09VTlRfTElNSVRfVFlQRV9NQVhfUkVVU0FCTEVfRE'
    'VMRUdBVElPTl9TRVRTX0JZX09XTkVSEAASPAo4QUNDT1VOVF9MSU1JVF9UWVBFX01BWF9UUkFG'
    'RklDX1BPTElDWV9JTlNUQU5DRVNfQllfT1dORVIQARIxCi1BQ0NPVU5UX0xJTUlUX1RZUEVfTU'
    'FYX0hFQUxUSF9DSEVDS1NfQllfT1dORVIQAhIwCixBQ0NPVU5UX0xJTUlUX1RZUEVfTUFYX0hP'
    'U1RFRF9aT05FU19CWV9PV05FUhADEjQKMEFDQ09VTlRfTElNSVRfVFlQRV9NQVhfVFJBRkZJQ1'
    '9QT0xJQ0lFU19CWV9PV05FUhAE');

@$core.Deprecated('Use changeActionDescriptor instead')
const ChangeAction$json = {
  '1': 'ChangeAction',
  '2': [
    {'1': 'CHANGE_ACTION_CREATE', '2': 0},
    {'1': 'CHANGE_ACTION_DELETE', '2': 1},
    {'1': 'CHANGE_ACTION_UPSERT', '2': 2},
  ],
};

/// Descriptor for `ChangeAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List changeActionDescriptor = $convert.base64Decode(
    'CgxDaGFuZ2VBY3Rpb24SGAoUQ0hBTkdFX0FDVElPTl9DUkVBVEUQABIYChRDSEFOR0VfQUNUSU'
    '9OX0RFTEVURRABEhgKFENIQU5HRV9BQ1RJT05fVVBTRVJUEAI=');

@$core.Deprecated('Use changeStatusDescriptor instead')
const ChangeStatus$json = {
  '1': 'ChangeStatus',
  '2': [
    {'1': 'CHANGE_STATUS_PENDING', '2': 0},
    {'1': 'CHANGE_STATUS_INSYNC', '2': 1},
  ],
};

/// Descriptor for `ChangeStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List changeStatusDescriptor = $convert.base64Decode(
    'CgxDaGFuZ2VTdGF0dXMSGQoVQ0hBTkdFX1NUQVRVU19QRU5ESU5HEAASGAoUQ0hBTkdFX1NUQV'
    'RVU19JTlNZTkMQAQ==');

@$core.Deprecated('Use cidrCollectionChangeActionDescriptor instead')
const CidrCollectionChangeAction$json = {
  '1': 'CidrCollectionChangeAction',
  '2': [
    {'1': 'CIDR_COLLECTION_CHANGE_ACTION_DELETE_IF_EXISTS', '2': 0},
    {'1': 'CIDR_COLLECTION_CHANGE_ACTION_PUT', '2': 1},
  ],
};

/// Descriptor for `CidrCollectionChangeAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cidrCollectionChangeActionDescriptor =
    $convert.base64Decode(
        'ChpDaWRyQ29sbGVjdGlvbkNoYW5nZUFjdGlvbhIyCi5DSURSX0NPTExFQ1RJT05fQ0hBTkdFX0'
        'FDVElPTl9ERUxFVEVfSUZfRVhJU1RTEAASJQohQ0lEUl9DT0xMRUNUSU9OX0NIQU5HRV9BQ1RJ'
        'T05fUFVUEAE=');

@$core.Deprecated('Use cloudWatchRegionDescriptor instead')
const CloudWatchRegion$json = {
  '1': 'CloudWatchRegion',
  '2': [
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTHEAST_3', '2': 0},
    {'1': 'CLOUD_WATCH_REGION_US_ISOF_SOUTH_1', '2': 1},
    {'1': 'CLOUD_WATCH_REGION_US_ISO_EAST_1', '2': 2},
    {'1': 'CLOUD_WATCH_REGION_EU_WEST_3', '2': 3},
    {'1': 'CLOUD_WATCH_REGION_US_ISO_WEST_1', '2': 4},
    {'1': 'CLOUD_WATCH_REGION_US_GOV_EAST_1', '2': 5},
    {'1': 'CLOUD_WATCH_REGION_AP_NORTHEAST_1', '2': 6},
    {'1': 'CLOUD_WATCH_REGION_IL_CENTRAL_1', '2': 7},
    {'1': 'CLOUD_WATCH_REGION_EU_NORTH_1', '2': 8},
    {'1': 'CLOUD_WATCH_REGION_CN_NORTH_1', '2': 9},
    {'1': 'CLOUD_WATCH_REGION_MX_CENTRAL_1', '2': 10},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTHEAST_5', '2': 11},
    {'1': 'CLOUD_WATCH_REGION_EUSC_DE_EAST_1', '2': 12},
    {'1': 'CLOUD_WATCH_REGION_US_ISOF_EAST_1', '2': 13},
    {'1': 'CLOUD_WATCH_REGION_AP_NORTHEAST_2', '2': 14},
    {'1': 'CLOUD_WATCH_REGION_US_WEST_1', '2': 15},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTH_1', '2': 16},
    {'1': 'CLOUD_WATCH_REGION_US_EAST_1', '2': 17},
    {'1': 'CLOUD_WATCH_REGION_US_ISOB_WEST_1', '2': 18},
    {'1': 'CLOUD_WATCH_REGION_CA_WEST_1', '2': 19},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTHEAST_2', '2': 20},
    {'1': 'CLOUD_WATCH_REGION_EU_SOUTH_1', '2': 21},
    {'1': 'CLOUD_WATCH_REGION_AP_EAST_2', '2': 22},
    {'1': 'CLOUD_WATCH_REGION_EU_WEST_2', '2': 23},
    {'1': 'CLOUD_WATCH_REGION_US_WEST_2', '2': 24},
    {'1': 'CLOUD_WATCH_REGION_ME_CENTRAL_1', '2': 25},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTHEAST_7', '2': 26},
    {'1': 'CLOUD_WATCH_REGION_US_ISOB_EAST_1', '2': 27},
    {'1': 'CLOUD_WATCH_REGION_CN_NORTHWEST_1', '2': 28},
    {'1': 'CLOUD_WATCH_REGION_EU_SOUTH_2', '2': 29},
    {'1': 'CLOUD_WATCH_REGION_AP_EAST_1', '2': 30},
    {'1': 'CLOUD_WATCH_REGION_EU_CENTRAL_1', '2': 31},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTHEAST_4', '2': 32},
    {'1': 'CLOUD_WATCH_REGION_CA_CENTRAL_1', '2': 33},
    {'1': 'CLOUD_WATCH_REGION_SA_EAST_1', '2': 34},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTH_2', '2': 35},
    {'1': 'CLOUD_WATCH_REGION_US_GOV_WEST_1', '2': 36},
    {'1': 'CLOUD_WATCH_REGION_EU_CENTRAL_2', '2': 37},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTHEAST_1', '2': 38},
    {'1': 'CLOUD_WATCH_REGION_AF_SOUTH_1', '2': 39},
    {'1': 'CLOUD_WATCH_REGION_ME_SOUTH_1', '2': 40},
    {'1': 'CLOUD_WATCH_REGION_EU_WEST_1', '2': 41},
    {'1': 'CLOUD_WATCH_REGION_EU_ISOE_WEST_1', '2': 42},
    {'1': 'CLOUD_WATCH_REGION_AP_SOUTHEAST_6', '2': 43},
    {'1': 'CLOUD_WATCH_REGION_AP_NORTHEAST_3', '2': 44},
    {'1': 'CLOUD_WATCH_REGION_US_EAST_2', '2': 45},
  ],
};

/// Descriptor for `CloudWatchRegion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cloudWatchRegionDescriptor = $convert.base64Decode(
    'ChBDbG91ZFdhdGNoUmVnaW9uEiUKIUNMT1VEX1dBVENIX1JFR0lPTl9BUF9TT1VUSEVBU1RfMx'
    'AAEiYKIkNMT1VEX1dBVENIX1JFR0lPTl9VU19JU09GX1NPVVRIXzEQARIkCiBDTE9VRF9XQVRD'
    'SF9SRUdJT05fVVNfSVNPX0VBU1RfMRACEiAKHENMT1VEX1dBVENIX1JFR0lPTl9FVV9XRVNUXz'
    'MQAxIkCiBDTE9VRF9XQVRDSF9SRUdJT05fVVNfSVNPX1dFU1RfMRAEEiQKIENMT1VEX1dBVENI'
    'X1JFR0lPTl9VU19HT1ZfRUFTVF8xEAUSJQohQ0xPVURfV0FUQ0hfUkVHSU9OX0FQX05PUlRIRU'
    'FTVF8xEAYSIwofQ0xPVURfV0FUQ0hfUkVHSU9OX0lMX0NFTlRSQUxfMRAHEiEKHUNMT1VEX1dB'
    'VENIX1JFR0lPTl9FVV9OT1JUSF8xEAgSIQodQ0xPVURfV0FUQ0hfUkVHSU9OX0NOX05PUlRIXz'
    'EQCRIjCh9DTE9VRF9XQVRDSF9SRUdJT05fTVhfQ0VOVFJBTF8xEAoSJQohQ0xPVURfV0FUQ0hf'
    'UkVHSU9OX0FQX1NPVVRIRUFTVF81EAsSJQohQ0xPVURfV0FUQ0hfUkVHSU9OX0VVU0NfREVfRU'
    'FTVF8xEAwSJQohQ0xPVURfV0FUQ0hfUkVHSU9OX1VTX0lTT0ZfRUFTVF8xEA0SJQohQ0xPVURf'
    'V0FUQ0hfUkVHSU9OX0FQX05PUlRIRUFTVF8yEA4SIAocQ0xPVURfV0FUQ0hfUkVHSU9OX1VTX1'
    'dFU1RfMRAPEiEKHUNMT1VEX1dBVENIX1JFR0lPTl9BUF9TT1VUSF8xEBASIAocQ0xPVURfV0FU'
    'Q0hfUkVHSU9OX1VTX0VBU1RfMRAREiUKIUNMT1VEX1dBVENIX1JFR0lPTl9VU19JU09CX1dFU1'
    'RfMRASEiAKHENMT1VEX1dBVENIX1JFR0lPTl9DQV9XRVNUXzEQExIlCiFDTE9VRF9XQVRDSF9S'
    'RUdJT05fQVBfU09VVEhFQVNUXzIQFBIhCh1DTE9VRF9XQVRDSF9SRUdJT05fRVVfU09VVEhfMR'
    'AVEiAKHENMT1VEX1dBVENIX1JFR0lPTl9BUF9FQVNUXzIQFhIgChxDTE9VRF9XQVRDSF9SRUdJ'
    'T05fRVVfV0VTVF8yEBcSIAocQ0xPVURfV0FUQ0hfUkVHSU9OX1VTX1dFU1RfMhAYEiMKH0NMT1'
    'VEX1dBVENIX1JFR0lPTl9NRV9DRU5UUkFMXzEQGRIlCiFDTE9VRF9XQVRDSF9SRUdJT05fQVBf'
    'U09VVEhFQVNUXzcQGhIlCiFDTE9VRF9XQVRDSF9SRUdJT05fVVNfSVNPQl9FQVNUXzEQGxIlCi'
    'FDTE9VRF9XQVRDSF9SRUdJT05fQ05fTk9SVEhXRVNUXzEQHBIhCh1DTE9VRF9XQVRDSF9SRUdJ'
    'T05fRVVfU09VVEhfMhAdEiAKHENMT1VEX1dBVENIX1JFR0lPTl9BUF9FQVNUXzEQHhIjCh9DTE'
    '9VRF9XQVRDSF9SRUdJT05fRVVfQ0VOVFJBTF8xEB8SJQohQ0xPVURfV0FUQ0hfUkVHSU9OX0FQ'
    'X1NPVVRIRUFTVF80ECASIwofQ0xPVURfV0FUQ0hfUkVHSU9OX0NBX0NFTlRSQUxfMRAhEiAKHE'
    'NMT1VEX1dBVENIX1JFR0lPTl9TQV9FQVNUXzEQIhIhCh1DTE9VRF9XQVRDSF9SRUdJT05fQVBf'
    'U09VVEhfMhAjEiQKIENMT1VEX1dBVENIX1JFR0lPTl9VU19HT1ZfV0VTVF8xECQSIwofQ0xPVU'
    'RfV0FUQ0hfUkVHSU9OX0VVX0NFTlRSQUxfMhAlEiUKIUNMT1VEX1dBVENIX1JFR0lPTl9BUF9T'
    'T1VUSEVBU1RfMRAmEiEKHUNMT1VEX1dBVENIX1JFR0lPTl9BRl9TT1VUSF8xECcSIQodQ0xPVU'
    'RfV0FUQ0hfUkVHSU9OX01FX1NPVVRIXzEQKBIgChxDTE9VRF9XQVRDSF9SRUdJT05fRVVfV0VT'
    'VF8xECkSJQohQ0xPVURfV0FUQ0hfUkVHSU9OX0VVX0lTT0VfV0VTVF8xECoSJQohQ0xPVURfV0'
    'FUQ0hfUkVHSU9OX0FQX1NPVVRIRUFTVF82ECsSJQohQ0xPVURfV0FUQ0hfUkVHSU9OX0FQX05P'
    'UlRIRUFTVF8zECwSIAocQ0xPVURfV0FUQ0hfUkVHSU9OX1VTX0VBU1RfMhAt');

@$core.Deprecated('Use comparisonOperatorDescriptor instead')
const ComparisonOperator$json = {
  '1': 'ComparisonOperator',
  '2': [
    {'1': 'COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD', '2': 0},
    {'1': 'COMPARISON_OPERATOR_GREATERTHANTHRESHOLD', '2': 1},
    {'1': 'COMPARISON_OPERATOR_LESSTHANTHRESHOLD', '2': 2},
    {'1': 'COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD', '2': 3},
  ],
};

/// Descriptor for `ComparisonOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List comparisonOperatorDescriptor = $convert.base64Decode(
    'ChJDb21wYXJpc29uT3BlcmF0b3ISNQoxQ09NUEFSSVNPTl9PUEVSQVRPUl9HUkVBVEVSVEhBTk'
    '9SRVFVQUxUT1RIUkVTSE9MRBAAEiwKKENPTVBBUklTT05fT1BFUkFUT1JfR1JFQVRFUlRIQU5U'
    'SFJFU0hPTEQQARIpCiVDT01QQVJJU09OX09QRVJBVE9SX0xFU1NUSEFOVEhSRVNIT0xEEAISMg'
    'ouQ09NUEFSSVNPTl9PUEVSQVRPUl9MRVNTVEhBTk9SRVFVQUxUT1RIUkVTSE9MRBAD');

@$core.Deprecated('Use healthCheckRegionDescriptor instead')
const HealthCheckRegion$json = {
  '1': 'HealthCheckRegion',
  '2': [
    {'1': 'HEALTH_CHECK_REGION_AP_NORTHEAST_1', '2': 0},
    {'1': 'HEALTH_CHECK_REGION_US_WEST_1', '2': 1},
    {'1': 'HEALTH_CHECK_REGION_US_EAST_1', '2': 2},
    {'1': 'HEALTH_CHECK_REGION_AP_SOUTHEAST_2', '2': 3},
    {'1': 'HEALTH_CHECK_REGION_US_WEST_2', '2': 4},
    {'1': 'HEALTH_CHECK_REGION_SA_EAST_1', '2': 5},
    {'1': 'HEALTH_CHECK_REGION_AP_SOUTHEAST_1', '2': 6},
    {'1': 'HEALTH_CHECK_REGION_EU_WEST_1', '2': 7},
  ],
};

/// Descriptor for `HealthCheckRegion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List healthCheckRegionDescriptor = $convert.base64Decode(
    'ChFIZWFsdGhDaGVja1JlZ2lvbhImCiJIRUFMVEhfQ0hFQ0tfUkVHSU9OX0FQX05PUlRIRUFTVF'
    '8xEAASIQodSEVBTFRIX0NIRUNLX1JFR0lPTl9VU19XRVNUXzEQARIhCh1IRUFMVEhfQ0hFQ0tf'
    'UkVHSU9OX1VTX0VBU1RfMRACEiYKIkhFQUxUSF9DSEVDS19SRUdJT05fQVBfU09VVEhFQVNUXz'
    'IQAxIhCh1IRUFMVEhfQ0hFQ0tfUkVHSU9OX1VTX1dFU1RfMhAEEiEKHUhFQUxUSF9DSEVDS19S'
    'RUdJT05fU0FfRUFTVF8xEAUSJgoiSEVBTFRIX0NIRUNLX1JFR0lPTl9BUF9TT1VUSEVBU1RfMR'
    'AGEiEKHUhFQUxUSF9DSEVDS19SRUdJT05fRVVfV0VTVF8xEAc=');

@$core.Deprecated('Use healthCheckTypeDescriptor instead')
const HealthCheckType$json = {
  '1': 'HealthCheckType',
  '2': [
    {'1': 'HEALTH_CHECK_TYPE_HTTP', '2': 0},
    {'1': 'HEALTH_CHECK_TYPE_TCP', '2': 1},
    {'1': 'HEALTH_CHECK_TYPE_RECOVERY_CONTROL', '2': 2},
    {'1': 'HEALTH_CHECK_TYPE_CALCULATED', '2': 3},
    {'1': 'HEALTH_CHECK_TYPE_CLOUDWATCH_METRIC', '2': 4},
    {'1': 'HEALTH_CHECK_TYPE_HTTPS_STR_MATCH', '2': 5},
    {'1': 'HEALTH_CHECK_TYPE_HTTP_STR_MATCH', '2': 6},
    {'1': 'HEALTH_CHECK_TYPE_HTTPS', '2': 7},
  ],
};

/// Descriptor for `HealthCheckType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List healthCheckTypeDescriptor = $convert.base64Decode(
    'Cg9IZWFsdGhDaGVja1R5cGUSGgoWSEVBTFRIX0NIRUNLX1RZUEVfSFRUUBAAEhkKFUhFQUxUSF'
    '9DSEVDS19UWVBFX1RDUBABEiYKIkhFQUxUSF9DSEVDS19UWVBFX1JFQ09WRVJZX0NPTlRST0wQ'
    'AhIgChxIRUFMVEhfQ0hFQ0tfVFlQRV9DQUxDVUxBVEVEEAMSJwojSEVBTFRIX0NIRUNLX1RZUE'
    'VfQ0xPVURXQVRDSF9NRVRSSUMQBBIlCiFIRUFMVEhfQ0hFQ0tfVFlQRV9IVFRQU19TVFJfTUFU'
    'Q0gQBRIkCiBIRUFMVEhfQ0hFQ0tfVFlQRV9IVFRQX1NUUl9NQVRDSBAGEhsKF0hFQUxUSF9DSE'
    'VDS19UWVBFX0hUVFBTEAc=');

@$core.Deprecated('Use hostedZoneLimitTypeDescriptor instead')
const HostedZoneLimitType$json = {
  '1': 'HostedZoneLimitType',
  '2': [
    {'1': 'HOSTED_ZONE_LIMIT_TYPE_MAX_RRSETS_BY_ZONE', '2': 0},
    {'1': 'HOSTED_ZONE_LIMIT_TYPE_MAX_VPCS_ASSOCIATED_BY_ZONE', '2': 1},
  ],
};

/// Descriptor for `HostedZoneLimitType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List hostedZoneLimitTypeDescriptor = $convert.base64Decode(
    'ChNIb3N0ZWRab25lTGltaXRUeXBlEi0KKUhPU1RFRF9aT05FX0xJTUlUX1RZUEVfTUFYX1JSU0'
    'VUU19CWV9aT05FEAASNgoySE9TVEVEX1pPTkVfTElNSVRfVFlQRV9NQVhfVlBDU19BU1NPQ0lB'
    'VEVEX0JZX1pPTkUQAQ==');

@$core.Deprecated('Use hostedZoneTypeDescriptor instead')
const HostedZoneType$json = {
  '1': 'HostedZoneType',
  '2': [
    {'1': 'HOSTED_ZONE_TYPE_PRIVATE_HOSTED_ZONE', '2': 0},
  ],
};

/// Descriptor for `HostedZoneType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List hostedZoneTypeDescriptor = $convert.base64Decode(
    'Cg5Ib3N0ZWRab25lVHlwZRIoCiRIT1NURURfWk9ORV9UWVBFX1BSSVZBVEVfSE9TVEVEX1pPTk'
    'UQAA==');

@$core.Deprecated('Use insufficientDataHealthStatusDescriptor instead')
const InsufficientDataHealthStatus$json = {
  '1': 'InsufficientDataHealthStatus',
  '2': [
    {'1': 'INSUFFICIENT_DATA_HEALTH_STATUS_HEALTHY', '2': 0},
    {'1': 'INSUFFICIENT_DATA_HEALTH_STATUS_LASTKNOWNSTATUS', '2': 1},
    {'1': 'INSUFFICIENT_DATA_HEALTH_STATUS_UNHEALTHY', '2': 2},
  ],
};

/// Descriptor for `InsufficientDataHealthStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List insufficientDataHealthStatusDescriptor = $convert.base64Decode(
    'ChxJbnN1ZmZpY2llbnREYXRhSGVhbHRoU3RhdHVzEisKJ0lOU1VGRklDSUVOVF9EQVRBX0hFQU'
    'xUSF9TVEFUVVNfSEVBTFRIWRAAEjMKL0lOU1VGRklDSUVOVF9EQVRBX0hFQUxUSF9TVEFUVVNf'
    'TEFTVEtOT1dOU1RBVFVTEAESLQopSU5TVUZGSUNJRU5UX0RBVEFfSEVBTFRIX1NUQVRVU19VTk'
    'hFQUxUSFkQAg==');

@$core.Deprecated('Use rRTypeDescriptor instead')
const RRType$json = {
  '1': 'RRType',
  '2': [
    {'1': 'R_R_TYPE_PTR', '2': 0},
    {'1': 'R_R_TYPE_CAA', '2': 1},
    {'1': 'R_R_TYPE_SOA', '2': 2},
    {'1': 'R_R_TYPE_SRV', '2': 3},
    {'1': 'R_R_TYPE_CNAME', '2': 4},
    {'1': 'R_R_TYPE_TLSA', '2': 5},
    {'1': 'R_R_TYPE_A', '2': 6},
    {'1': 'R_R_TYPE_MX', '2': 7},
    {'1': 'R_R_TYPE_SSHFP', '2': 8},
    {'1': 'R_R_TYPE_DS', '2': 9},
    {'1': 'R_R_TYPE_AAAA', '2': 10},
    {'1': 'R_R_TYPE_NAPTR', '2': 11},
    {'1': 'R_R_TYPE_NS', '2': 12},
    {'1': 'R_R_TYPE_HTTPS', '2': 13},
    {'1': 'R_R_TYPE_SPF', '2': 14},
    {'1': 'R_R_TYPE_SVCB', '2': 15},
    {'1': 'R_R_TYPE_TXT', '2': 16},
  ],
};

/// Descriptor for `RRType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List rRTypeDescriptor = $convert.base64Decode(
    'CgZSUlR5cGUSEAoMUl9SX1RZUEVfUFRSEAASEAoMUl9SX1RZUEVfQ0FBEAESEAoMUl9SX1RZUE'
    'VfU09BEAISEAoMUl9SX1RZUEVfU1JWEAMSEgoOUl9SX1RZUEVfQ05BTUUQBBIRCg1SX1JfVFlQ'
    'RV9UTFNBEAUSDgoKUl9SX1RZUEVfQRAGEg8KC1JfUl9UWVBFX01YEAcSEgoOUl9SX1RZUEVfU1'
    'NIRlAQCBIPCgtSX1JfVFlQRV9EUxAJEhEKDVJfUl9UWVBFX0FBQUEQChISCg5SX1JfVFlQRV9O'
    'QVBUUhALEg8KC1JfUl9UWVBFX05TEAwSEgoOUl9SX1RZUEVfSFRUUFMQDRIQCgxSX1JfVFlQRV'
    '9TUEYQDhIRCg1SX1JfVFlQRV9TVkNCEA8SEAoMUl9SX1RZUEVfVFhUEBA=');

@$core.Deprecated('Use resettableElementNameDescriptor instead')
const ResettableElementName$json = {
  '1': 'ResettableElementName',
  '2': [
    {'1': 'RESETTABLE_ELEMENT_NAME_REGIONS', '2': 0},
    {'1': 'RESETTABLE_ELEMENT_NAME_CHILDHEALTHCHECKS', '2': 1},
    {'1': 'RESETTABLE_ELEMENT_NAME_FULLYQUALIFIEDDOMAINNAME', '2': 2},
    {'1': 'RESETTABLE_ELEMENT_NAME_RESOURCEPATH', '2': 3},
  ],
};

/// Descriptor for `ResettableElementName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resettableElementNameDescriptor = $convert.base64Decode(
    'ChVSZXNldHRhYmxlRWxlbWVudE5hbWUSIwofUkVTRVRUQUJMRV9FTEVNRU5UX05BTUVfUkVHSU'
    '9OUxAAEi0KKVJFU0VUVEFCTEVfRUxFTUVOVF9OQU1FX0NISUxESEVBTFRIQ0hFQ0tTEAESNAow'
    'UkVTRVRUQUJMRV9FTEVNRU5UX05BTUVfRlVMTFlRVUFMSUZJRURET01BSU5OQU1FEAISKAokUk'
    'VTRVRUQUJMRV9FTEVNRU5UX05BTUVfUkVTT1VSQ0VQQVRIEAM=');

@$core.Deprecated('Use resourceRecordSetFailoverDescriptor instead')
const ResourceRecordSetFailover$json = {
  '1': 'ResourceRecordSetFailover',
  '2': [
    {'1': 'RESOURCE_RECORD_SET_FAILOVER_SECONDARY', '2': 0},
    {'1': 'RESOURCE_RECORD_SET_FAILOVER_PRIMARY', '2': 1},
  ],
};

/// Descriptor for `ResourceRecordSetFailover`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resourceRecordSetFailoverDescriptor = $convert.base64Decode(
    'ChlSZXNvdXJjZVJlY29yZFNldEZhaWxvdmVyEioKJlJFU09VUkNFX1JFQ09SRF9TRVRfRkFJTE'
    '9WRVJfU0VDT05EQVJZEAASKAokUkVTT1VSQ0VfUkVDT1JEX1NFVF9GQUlMT1ZFUl9QUklNQVJZ'
    'EAE=');

@$core.Deprecated('Use resourceRecordSetRegionDescriptor instead')
const ResourceRecordSetRegion$json = {
  '1': 'ResourceRecordSetRegion',
  '2': [
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_3', '2': 0},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_WEST_3', '2': 1},
    {'1': 'RESOURCE_RECORD_SET_REGION_US_GOV_EAST_1', '2': 2},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_1', '2': 3},
    {'1': 'RESOURCE_RECORD_SET_REGION_IL_CENTRAL_1', '2': 4},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_NORTH_1', '2': 5},
    {'1': 'RESOURCE_RECORD_SET_REGION_CN_NORTH_1', '2': 6},
    {'1': 'RESOURCE_RECORD_SET_REGION_MX_CENTRAL_1', '2': 7},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_5', '2': 8},
    {'1': 'RESOURCE_RECORD_SET_REGION_EUSC_DE_EAST_1', '2': 9},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_2', '2': 10},
    {'1': 'RESOURCE_RECORD_SET_REGION_US_WEST_1', '2': 11},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTH_1', '2': 12},
    {'1': 'RESOURCE_RECORD_SET_REGION_US_EAST_1', '2': 13},
    {'1': 'RESOURCE_RECORD_SET_REGION_CA_WEST_1', '2': 14},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_2', '2': 15},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_SOUTH_1', '2': 16},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_EAST_2', '2': 17},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_WEST_2', '2': 18},
    {'1': 'RESOURCE_RECORD_SET_REGION_US_WEST_2', '2': 19},
    {'1': 'RESOURCE_RECORD_SET_REGION_ME_CENTRAL_1', '2': 20},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_7', '2': 21},
    {'1': 'RESOURCE_RECORD_SET_REGION_CN_NORTHWEST_1', '2': 22},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_SOUTH_2', '2': 23},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_EAST_1', '2': 24},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_CENTRAL_1', '2': 25},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_4', '2': 26},
    {'1': 'RESOURCE_RECORD_SET_REGION_CA_CENTRAL_1', '2': 27},
    {'1': 'RESOURCE_RECORD_SET_REGION_SA_EAST_1', '2': 28},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTH_2', '2': 29},
    {'1': 'RESOURCE_RECORD_SET_REGION_US_GOV_WEST_1', '2': 30},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_CENTRAL_2', '2': 31},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_1', '2': 32},
    {'1': 'RESOURCE_RECORD_SET_REGION_AF_SOUTH_1', '2': 33},
    {'1': 'RESOURCE_RECORD_SET_REGION_ME_SOUTH_1', '2': 34},
    {'1': 'RESOURCE_RECORD_SET_REGION_EU_WEST_1', '2': 35},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_6', '2': 36},
    {'1': 'RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_3', '2': 37},
    {'1': 'RESOURCE_RECORD_SET_REGION_US_EAST_2', '2': 38},
  ],
};

/// Descriptor for `ResourceRecordSetRegion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resourceRecordSetRegionDescriptor = $convert.base64Decode(
    'ChdSZXNvdXJjZVJlY29yZFNldFJlZ2lvbhItCilSRVNPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl'
    '9BUF9TT1VUSEVBU1RfMxAAEigKJFJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX0VVX1dFU1Rf'
    'MxABEiwKKFJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX1VTX0dPVl9FQVNUXzEQAhItCilSRV'
    'NPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl9BUF9OT1JUSEVBU1RfMRADEisKJ1JFU09VUkNFX1JF'
    'Q09SRF9TRVRfUkVHSU9OX0lMX0NFTlRSQUxfMRAEEikKJVJFU09VUkNFX1JFQ09SRF9TRVRfUk'
    'VHSU9OX0VVX05PUlRIXzEQBRIpCiVSRVNPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl9DTl9OT1JU'
    'SF8xEAYSKwonUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fTVhfQ0VOVFJBTF8xEAcSLQopUk'
    'VTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fQVBfU09VVEhFQVNUXzUQCBItCilSRVNPVVJDRV9S'
    'RUNPUkRfU0VUX1JFR0lPTl9FVVNDX0RFX0VBU1RfMRAJEi0KKVJFU09VUkNFX1JFQ09SRF9TRV'
    'RfUkVHSU9OX0FQX05PUlRIRUFTVF8yEAoSKAokUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05f'
    'VVNfV0VTVF8xEAsSKQolUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fQVBfU09VVEhfMRAMEi'
    'gKJFJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX1VTX0VBU1RfMRANEigKJFJFU09VUkNFX1JF'
    'Q09SRF9TRVRfUkVHSU9OX0NBX1dFU1RfMRAOEi0KKVJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU'
    '9OX0FQX1NPVVRIRUFTVF8yEA8SKQolUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fRVVfU09V'
    'VEhfMRAQEigKJFJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX0FQX0VBU1RfMhAREigKJFJFU0'
    '9VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX0VVX1dFU1RfMhASEigKJFJFU09VUkNFX1JFQ09SRF9T'
    'RVRfUkVHSU9OX1VTX1dFU1RfMhATEisKJ1JFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX01FX0'
    'NFTlRSQUxfMRAUEi0KKVJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX0FQX1NPVVRIRUFTVF83'
    'EBUSLQopUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fQ05fTk9SVEhXRVNUXzEQFhIpCiVSRV'
    'NPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl9FVV9TT1VUSF8yEBcSKAokUkVTT1VSQ0VfUkVDT1JE'
    'X1NFVF9SRUdJT05fQVBfRUFTVF8xEBgSKwonUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fRV'
    'VfQ0VOVFJBTF8xEBkSLQopUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fQVBfU09VVEhFQVNU'
    'XzQQGhIrCidSRVNPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl9DQV9DRU5UUkFMXzEQGxIoCiRSRV'
    'NPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl9TQV9FQVNUXzEQHBIpCiVSRVNPVVJDRV9SRUNPUkRf'
    'U0VUX1JFR0lPTl9BUF9TT1VUSF8yEB0SLAooUkVTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fVV'
    'NfR09WX1dFU1RfMRAeEisKJ1JFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX0VVX0NFTlRSQUxf'
    'MhAfEi0KKVJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX0FQX1NPVVRIRUFTVF8xECASKQolUk'
    'VTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fQUZfU09VVEhfMRAhEikKJVJFU09VUkNFX1JFQ09S'
    'RF9TRVRfUkVHSU9OX01FX1NPVVRIXzEQIhIoCiRSRVNPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl'
    '9FVV9XRVNUXzEQIxItCilSRVNPVVJDRV9SRUNPUkRfU0VUX1JFR0lPTl9BUF9TT1VUSEVBU1Rf'
    'NhAkEi0KKVJFU09VUkNFX1JFQ09SRF9TRVRfUkVHSU9OX0FQX05PUlRIRUFTVF8zECUSKAokUk'
    'VTT1VSQ0VfUkVDT1JEX1NFVF9SRUdJT05fVVNfRUFTVF8yECY=');

@$core.Deprecated('Use reusableDelegationSetLimitTypeDescriptor instead')
const ReusableDelegationSetLimitType$json = {
  '1': 'ReusableDelegationSetLimitType',
  '2': [
    {
      '1':
          'REUSABLE_DELEGATION_SET_LIMIT_TYPE_MAX_ZONES_BY_REUSABLE_DELEGATION_SET',
      '2': 0
    },
  ],
};

/// Descriptor for `ReusableDelegationSetLimitType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reusableDelegationSetLimitTypeDescriptor =
    $convert.base64Decode(
        'Ch5SZXVzYWJsZURlbGVnYXRpb25TZXRMaW1pdFR5cGUSSwpHUkVVU0FCTEVfREVMRUdBVElPTl'
        '9TRVRfTElNSVRfVFlQRV9NQVhfWk9ORVNfQllfUkVVU0FCTEVfREVMRUdBVElPTl9TRVQQAA==');

@$core.Deprecated('Use statisticDescriptor instead')
const Statistic$json = {
  '1': 'Statistic',
  '2': [
    {'1': 'STATISTIC_SUM', '2': 0},
    {'1': 'STATISTIC_SAMPLECOUNT', '2': 1},
    {'1': 'STATISTIC_AVERAGE', '2': 2},
    {'1': 'STATISTIC_MAXIMUM', '2': 3},
    {'1': 'STATISTIC_MINIMUM', '2': 4},
  ],
};

/// Descriptor for `Statistic`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statisticDescriptor = $convert.base64Decode(
    'CglTdGF0aXN0aWMSEQoNU1RBVElTVElDX1NVTRAAEhkKFVNUQVRJU1RJQ19TQU1QTEVDT1VOVB'
    'ABEhUKEVNUQVRJU1RJQ19BVkVSQUdFEAISFQoRU1RBVElTVElDX01BWElNVU0QAxIVChFTVEFU'
    'SVNUSUNfTUlOSU1VTRAE');

@$core.Deprecated('Use tagResourceTypeDescriptor instead')
const TagResourceType$json = {
  '1': 'TagResourceType',
  '2': [
    {'1': 'TAG_RESOURCE_TYPE_HEALTHCHECK', '2': 0},
    {'1': 'TAG_RESOURCE_TYPE_HOSTEDZONE', '2': 1},
  ],
};

/// Descriptor for `TagResourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tagResourceTypeDescriptor = $convert.base64Decode(
    'Cg9UYWdSZXNvdXJjZVR5cGUSIQodVEFHX1JFU09VUkNFX1RZUEVfSEVBTFRIQ0hFQ0sQABIgCh'
    'xUQUdfUkVTT1VSQ0VfVFlQRV9IT1NURURaT05FEAE=');

@$core.Deprecated('Use vPCRegionDescriptor instead')
const VPCRegion$json = {
  '1': 'VPCRegion',
  '2': [
    {'1': 'V_P_C_REGION_AP_SOUTHEAST_3', '2': 0},
    {'1': 'V_P_C_REGION_US_ISOF_SOUTH_1', '2': 1},
    {'1': 'V_P_C_REGION_US_ISO_EAST_1', '2': 2},
    {'1': 'V_P_C_REGION_EU_WEST_3', '2': 3},
    {'1': 'V_P_C_REGION_US_ISO_WEST_1', '2': 4},
    {'1': 'V_P_C_REGION_US_GOV_EAST_1', '2': 5},
    {'1': 'V_P_C_REGION_AP_NORTHEAST_1', '2': 6},
    {'1': 'V_P_C_REGION_IL_CENTRAL_1', '2': 7},
    {'1': 'V_P_C_REGION_EU_NORTH_1', '2': 8},
    {'1': 'V_P_C_REGION_CN_NORTH_1', '2': 9},
    {'1': 'V_P_C_REGION_MX_CENTRAL_1', '2': 10},
    {'1': 'V_P_C_REGION_AP_SOUTHEAST_5', '2': 11},
    {'1': 'V_P_C_REGION_EUSC_DE_EAST_1', '2': 12},
    {'1': 'V_P_C_REGION_US_ISOF_EAST_1', '2': 13},
    {'1': 'V_P_C_REGION_AP_NORTHEAST_2', '2': 14},
    {'1': 'V_P_C_REGION_US_WEST_1', '2': 15},
    {'1': 'V_P_C_REGION_AP_SOUTH_1', '2': 16},
    {'1': 'V_P_C_REGION_US_EAST_1', '2': 17},
    {'1': 'V_P_C_REGION_US_ISOB_WEST_1', '2': 18},
    {'1': 'V_P_C_REGION_CA_WEST_1', '2': 19},
    {'1': 'V_P_C_REGION_AP_SOUTHEAST_2', '2': 20},
    {'1': 'V_P_C_REGION_EU_SOUTH_1', '2': 21},
    {'1': 'V_P_C_REGION_AP_EAST_2', '2': 22},
    {'1': 'V_P_C_REGION_EU_WEST_2', '2': 23},
    {'1': 'V_P_C_REGION_US_WEST_2', '2': 24},
    {'1': 'V_P_C_REGION_ME_CENTRAL_1', '2': 25},
    {'1': 'V_P_C_REGION_AP_SOUTHEAST_7', '2': 26},
    {'1': 'V_P_C_REGION_US_ISOB_EAST_1', '2': 27},
    {'1': 'V_P_C_REGION_CN_NORTHWEST_1', '2': 28},
    {'1': 'V_P_C_REGION_EU_SOUTH_2', '2': 29},
    {'1': 'V_P_C_REGION_AP_EAST_1', '2': 30},
    {'1': 'V_P_C_REGION_EU_CENTRAL_1', '2': 31},
    {'1': 'V_P_C_REGION_AP_SOUTHEAST_4', '2': 32},
    {'1': 'V_P_C_REGION_CA_CENTRAL_1', '2': 33},
    {'1': 'V_P_C_REGION_SA_EAST_1', '2': 34},
    {'1': 'V_P_C_REGION_AP_SOUTH_2', '2': 35},
    {'1': 'V_P_C_REGION_US_GOV_WEST_1', '2': 36},
    {'1': 'V_P_C_REGION_EU_CENTRAL_2', '2': 37},
    {'1': 'V_P_C_REGION_AP_SOUTHEAST_1', '2': 38},
    {'1': 'V_P_C_REGION_AF_SOUTH_1', '2': 39},
    {'1': 'V_P_C_REGION_ME_SOUTH_1', '2': 40},
    {'1': 'V_P_C_REGION_EU_WEST_1', '2': 41},
    {'1': 'V_P_C_REGION_EU_ISOE_WEST_1', '2': 42},
    {'1': 'V_P_C_REGION_AP_SOUTHEAST_6', '2': 43},
    {'1': 'V_P_C_REGION_AP_NORTHEAST_3', '2': 44},
    {'1': 'V_P_C_REGION_US_EAST_2', '2': 45},
  ],
};

/// Descriptor for `VPCRegion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List vPCRegionDescriptor = $convert.base64Decode(
    'CglWUENSZWdpb24SHwobVl9QX0NfUkVHSU9OX0FQX1NPVVRIRUFTVF8zEAASIAocVl9QX0NfUk'
    'VHSU9OX1VTX0lTT0ZfU09VVEhfMRABEh4KGlZfUF9DX1JFR0lPTl9VU19JU09fRUFTVF8xEAIS'
    'GgoWVl9QX0NfUkVHSU9OX0VVX1dFU1RfMxADEh4KGlZfUF9DX1JFR0lPTl9VU19JU09fV0VTVF'
    '8xEAQSHgoaVl9QX0NfUkVHSU9OX1VTX0dPVl9FQVNUXzEQBRIfChtWX1BfQ19SRUdJT05fQVBf'
    'Tk9SVEhFQVNUXzEQBhIdChlWX1BfQ19SRUdJT05fSUxfQ0VOVFJBTF8xEAcSGwoXVl9QX0NfUk'
    'VHSU9OX0VVX05PUlRIXzEQCBIbChdWX1BfQ19SRUdJT05fQ05fTk9SVEhfMRAJEh0KGVZfUF9D'
    'X1JFR0lPTl9NWF9DRU5UUkFMXzEQChIfChtWX1BfQ19SRUdJT05fQVBfU09VVEhFQVNUXzUQCx'
    'IfChtWX1BfQ19SRUdJT05fRVVTQ19ERV9FQVNUXzEQDBIfChtWX1BfQ19SRUdJT05fVVNfSVNP'
    'Rl9FQVNUXzEQDRIfChtWX1BfQ19SRUdJT05fQVBfTk9SVEhFQVNUXzIQDhIaChZWX1BfQ19SRU'
    'dJT05fVVNfV0VTVF8xEA8SGwoXVl9QX0NfUkVHSU9OX0FQX1NPVVRIXzEQEBIaChZWX1BfQ19S'
    'RUdJT05fVVNfRUFTVF8xEBESHwobVl9QX0NfUkVHSU9OX1VTX0lTT0JfV0VTVF8xEBISGgoWVl'
    '9QX0NfUkVHSU9OX0NBX1dFU1RfMRATEh8KG1ZfUF9DX1JFR0lPTl9BUF9TT1VUSEVBU1RfMhAU'
    'EhsKF1ZfUF9DX1JFR0lPTl9FVV9TT1VUSF8xEBUSGgoWVl9QX0NfUkVHSU9OX0FQX0VBU1RfMh'
    'AWEhoKFlZfUF9DX1JFR0lPTl9FVV9XRVNUXzIQFxIaChZWX1BfQ19SRUdJT05fVVNfV0VTVF8y'
    'EBgSHQoZVl9QX0NfUkVHSU9OX01FX0NFTlRSQUxfMRAZEh8KG1ZfUF9DX1JFR0lPTl9BUF9TT1'
    'VUSEVBU1RfNxAaEh8KG1ZfUF9DX1JFR0lPTl9VU19JU09CX0VBU1RfMRAbEh8KG1ZfUF9DX1JF'
    'R0lPTl9DTl9OT1JUSFdFU1RfMRAcEhsKF1ZfUF9DX1JFR0lPTl9FVV9TT1VUSF8yEB0SGgoWVl'
    '9QX0NfUkVHSU9OX0FQX0VBU1RfMRAeEh0KGVZfUF9DX1JFR0lPTl9FVV9DRU5UUkFMXzEQHxIf'
    'ChtWX1BfQ19SRUdJT05fQVBfU09VVEhFQVNUXzQQIBIdChlWX1BfQ19SRUdJT05fQ0FfQ0VOVF'
    'JBTF8xECESGgoWVl9QX0NfUkVHSU9OX1NBX0VBU1RfMRAiEhsKF1ZfUF9DX1JFR0lPTl9BUF9T'
    'T1VUSF8yECMSHgoaVl9QX0NfUkVHSU9OX1VTX0dPVl9XRVNUXzEQJBIdChlWX1BfQ19SRUdJT0'
    '5fRVVfQ0VOVFJBTF8yECUSHwobVl9QX0NfUkVHSU9OX0FQX1NPVVRIRUFTVF8xECYSGwoXVl9Q'
    'X0NfUkVHSU9OX0FGX1NPVVRIXzEQJxIbChdWX1BfQ19SRUdJT05fTUVfU09VVEhfMRAoEhoKFl'
    'ZfUF9DX1JFR0lPTl9FVV9XRVNUXzEQKRIfChtWX1BfQ19SRUdJT05fRVVfSVNPRV9XRVNUXzEQ'
    'KhIfChtWX1BfQ19SRUdJT05fQVBfU09VVEhFQVNUXzYQKxIfChtWX1BfQ19SRUdJT05fQVBfTk'
    '9SVEhFQVNUXzMQLBIaChZWX1BfQ19SRUdJT05fVVNfRUFTVF8yEC0=');

@$core.Deprecated('Use accountLimitDescriptor instead')
const AccountLimit$json = {
  '1': 'AccountLimit',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.AccountLimitType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 3, '10': 'value'},
  ],
};

/// Descriptor for `AccountLimit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountLimitDescriptor = $convert.base64Decode(
    'CgxBY2NvdW50TGltaXQSMQoEdHlwZRjuoNeKASABKA4yGS5yb3V0ZTUzLkFjY291bnRMaW1pdF'
    'R5cGVSBHR5cGUSGAoFdmFsdWUY6/KfigEgASgDUgV2YWx1ZQ==');

@$core.Deprecated('Use activateKeySigningKeyRequestDescriptor instead')
const ActivateKeySigningKeyRequest$json = {
  '1': 'ActivateKeySigningKeyRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `ActivateKeySigningKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activateKeySigningKeyRequestDescriptor =
    $convert.base64Decode(
        'ChxBY3RpdmF0ZUtleVNpZ25pbmdLZXlSZXF1ZXN0EiYKDGhvc3RlZHpvbmVpZBj+zp6lASABKA'
        'lSDGhvc3RlZHpvbmVpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use activateKeySigningKeyResponseDescriptor instead')
const ActivateKeySigningKeyResponse$json = {
  '1': 'ActivateKeySigningKeyResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `ActivateKeySigningKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activateKeySigningKeyResponseDescriptor =
    $convert.base64Decode(
        'Ch1BY3RpdmF0ZUtleVNpZ25pbmdLZXlSZXNwb25zZRI3CgpjaGFuZ2VpbmZvGO61i9ABIAEoCz'
        'ITLnJvdXRlNTMuQ2hhbmdlSW5mb1IKY2hhbmdlaW5mbw==');

@$core.Deprecated('Use alarmIdentifierDescriptor instead')
const AlarmIdentifier$json = {
  '1': 'AlarmIdentifier',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'region',
      '3': 154040478,
      '4': 1,
      '5': 14,
      '6': '.route53.CloudWatchRegion',
      '10': 'region'
    },
  ],
};

/// Descriptor for `AlarmIdentifier`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alarmIdentifierDescriptor = $convert.base64Decode(
    'Cg9BbGFybUlkZW50aWZpZXISFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRI0CgZyZWdpb24YnvG5SS'
    'ABKA4yGS5yb3V0ZTUzLkNsb3VkV2F0Y2hSZWdpb25SBnJlZ2lvbg==');

@$core.Deprecated('Use aliasTargetDescriptor instead')
const AliasTarget$json = {
  '1': 'AliasTarget',
  '2': [
    {'1': 'dnsname', '3': 171901432, '4': 1, '5': 9, '10': 'dnsname'},
    {
      '1': 'evaluatetargethealth',
      '3': 409666542,
      '4': 1,
      '5': 8,
      '10': 'evaluatetargethealth'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `AliasTarget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aliasTargetDescriptor = $convert.base64Decode(
    'CgtBbGlhc1RhcmdldBIbCgdkbnNuYW1lGPiD/FEgASgJUgdkbnNuYW1lEjYKFGV2YWx1YXRldG'
    'FyZ2V0aGVhbHRoGO6HrMMBIAEoCFIUZXZhbHVhdGV0YXJnZXRoZWFsdGgSJgoMaG9zdGVkem9u'
    'ZWlkGP7OnqUBIAEoCVIMaG9zdGVkem9uZWlk');

@$core.Deprecated('Use associateVPCWithHostedZoneRequestDescriptor instead')
const AssociateVPCWithHostedZoneRequest$json = {
  '1': 'AssociateVPCWithHostedZoneRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'vpc',
      '3': 506158953,
      '4': 1,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpc'
    },
  ],
};

/// Descriptor for `AssociateVPCWithHostedZoneRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateVPCWithHostedZoneRequestDescriptor =
    $convert.base64Decode(
        'CiFBc3NvY2lhdGVWUENXaXRoSG9zdGVkWm9uZVJlcXVlc3QSHAoHY29tbWVudBj/v77CASABKA'
        'lSB2NvbW1lbnQSJgoMaG9zdGVkem9uZWlkGP7OnqUBIAEoCVIMaG9zdGVkem9uZWlkEiIKA3Zw'
        'Yxjpvq3xASABKAsyDC5yb3V0ZTUzLlZQQ1IDdnBj');

@$core.Deprecated('Use associateVPCWithHostedZoneResponseDescriptor instead')
const AssociateVPCWithHostedZoneResponse$json = {
  '1': 'AssociateVPCWithHostedZoneResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `AssociateVPCWithHostedZoneResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateVPCWithHostedZoneResponseDescriptor =
    $convert.base64Decode(
        'CiJBc3NvY2lhdGVWUENXaXRoSG9zdGVkWm9uZVJlc3BvbnNlEjcKCmNoYW5nZWluZm8Y7rWL0A'
        'EgASgLMhMucm91dGU1My5DaGFuZ2VJbmZvUgpjaGFuZ2VpbmZv');

@$core.Deprecated('Use changeDescriptor instead')
const Change$json = {
  '1': 'Change',
  '2': [
    {
      '1': 'action',
      '3': 175614240,
      '4': 1,
      '5': 14,
      '6': '.route53.ChangeAction',
      '10': 'action'
    },
    {
      '1': 'resourcerecordset',
      '3': 344937781,
      '4': 1,
      '5': 11,
      '6': '.route53.ResourceRecordSet',
      '10': 'resourcerecordset'
    },
  ],
};

/// Descriptor for `Change`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeDescriptor = $convert.base64Decode(
    'CgZDaGFuZ2USMAoGYWN0aW9uGKDS3lMgASgOMhUucm91dGU1My5DaGFuZ2VBY3Rpb25SBmFjdG'
    'lvbhJMChFyZXNvdXJjZXJlY29yZHNldBi1qr2kASABKAsyGi5yb3V0ZTUzLlJlc291cmNlUmVj'
    'b3JkU2V0UhFyZXNvdXJjZXJlY29yZHNldA==');

@$core.Deprecated('Use changeBatchDescriptor instead')
const ChangeBatch$json = {
  '1': 'ChangeBatch',
  '2': [
    {
      '1': 'changes',
      '3': 516230891,
      '4': 3,
      '5': 11,
      '6': '.route53.Change',
      '10': 'changes'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
  ],
};

/// Descriptor for `ChangeBatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeBatchDescriptor = $convert.base64Decode(
    'CgtDaGFuZ2VCYXRjaBItCgdjaGFuZ2VzGOudlPYBIAMoCzIPLnJvdXRlNTMuQ2hhbmdlUgdjaG'
    'FuZ2VzEhwKB2NvbW1lbnQY/7++wgEgASgJUgdjb21tZW50');

@$core.Deprecated('Use changeCidrCollectionRequestDescriptor instead')
const ChangeCidrCollectionRequest$json = {
  '1': 'ChangeCidrCollectionRequest',
  '2': [
    {
      '1': 'changes',
      '3': 516230891,
      '4': 3,
      '5': 11,
      '6': '.route53.CidrCollectionChange',
      '10': 'changes'
    },
    {
      '1': 'collectionversion',
      '3': 95829948,
      '4': 1,
      '5': 3,
      '10': 'collectionversion'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `ChangeCidrCollectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeCidrCollectionRequestDescriptor =
    $convert.base64Decode(
        'ChtDaGFuZ2VDaWRyQ29sbGVjdGlvblJlcXVlc3QSOwoHY2hhbmdlcxjrnZT2ASADKAsyHS5yb3'
        'V0ZTUzLkNpZHJDb2xsZWN0aW9uQ2hhbmdlUgdjaGFuZ2VzEi8KEWNvbGxlY3Rpb252ZXJzaW9u'
        'GLz/2C0gASgDUhFjb2xsZWN0aW9udmVyc2lvbhISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use changeCidrCollectionResponseDescriptor instead')
const ChangeCidrCollectionResponse$json = {
  '1': 'ChangeCidrCollectionResponse',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `ChangeCidrCollectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeCidrCollectionResponseDescriptor =
    $convert.base64Decode(
        'ChxDaGFuZ2VDaWRyQ29sbGVjdGlvblJlc3BvbnNlEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use changeInfoDescriptor instead')
const ChangeInfo$json = {
  '1': 'ChangeInfo',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.route53.ChangeStatus',
      '10': 'status'
    },
    {'1': 'submittedat', '3': 343958936, '4': 1, '5': 9, '10': 'submittedat'},
  ],
};

/// Descriptor for `ChangeInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeInfoDescriptor = $convert.base64Decode(
    'CgpDaGFuZ2VJbmZvEhwKB2NvbW1lbnQY/7++wgEgASgJUgdjb21tZW50EhIKAmlkGIHyorcBIA'
    'EoCVICaWQSMAoGc3RhdHVzGJDk+wIgASgOMhUucm91dGU1My5DaGFuZ2VTdGF0dXNSBnN0YXR1'
    'cxIkCgtzdWJtaXR0ZWRhdBiYy4GkASABKAlSC3N1Ym1pdHRlZGF0');

@$core.Deprecated('Use changeResourceRecordSetsRequestDescriptor instead')
const ChangeResourceRecordSetsRequest$json = {
  '1': 'ChangeResourceRecordSetsRequest',
  '2': [
    {
      '1': 'changebatch',
      '3': 23600436,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeBatch',
      '10': 'changebatch'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `ChangeResourceRecordSetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeResourceRecordSetsRequestDescriptor =
    $convert.base64Decode(
        'Ch9DaGFuZ2VSZXNvdXJjZVJlY29yZFNldHNSZXF1ZXN0EjkKC2NoYW5nZWJhdGNoGLS6oAsgAS'
        'gLMhQucm91dGU1My5DaGFuZ2VCYXRjaFILY2hhbmdlYmF0Y2gSJgoMaG9zdGVkem9uZWlkGP7O'
        'nqUBIAEoCVIMaG9zdGVkem9uZWlk');

@$core.Deprecated('Use changeResourceRecordSetsResponseDescriptor instead')
const ChangeResourceRecordSetsResponse$json = {
  '1': 'ChangeResourceRecordSetsResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `ChangeResourceRecordSetsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeResourceRecordSetsResponseDescriptor =
    $convert.base64Decode(
        'CiBDaGFuZ2VSZXNvdXJjZVJlY29yZFNldHNSZXNwb25zZRI3CgpjaGFuZ2VpbmZvGO61i9ABIA'
        'EoCzITLnJvdXRlNTMuQ2hhbmdlSW5mb1IKY2hhbmdlaW5mbw==');

@$core.Deprecated('Use changeTagsForResourceRequestDescriptor instead')
const ChangeTagsForResourceRequest$json = {
  '1': 'ChangeTagsForResourceRequest',
  '2': [
    {
      '1': 'addtags',
      '3': 503857054,
      '4': 3,
      '5': 11,
      '6': '.route53.Tag',
      '10': 'addtags'
    },
    {
      '1': 'removetagkeys',
      '3': 128174116,
      '4': 3,
      '5': 9,
      '10': 'removetagkeys'
    },
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.route53.TagResourceType',
      '10': 'resourcetype'
    },
  ],
};

/// Descriptor for `ChangeTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeTagsForResourceRequestDescriptor = $convert.base64Decode(
    'ChxDaGFuZ2VUYWdzRm9yUmVzb3VyY2VSZXF1ZXN0EioKB2FkZHRhZ3MYnv+g8AEgAygLMgwucm'
    '91dGU1My5UYWdSB2FkZHRhZ3MSJwoNcmVtb3ZldGFna2V5cxikkI89IAMoCVINcmVtb3ZldGFn'
    'a2V5cxIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcmVzb3VyY2VpZBJACgxyZXNvdXJjZXR5cG'
    'UY3r7YjwEgASgOMhgucm91dGU1My5UYWdSZXNvdXJjZVR5cGVSDHJlc291cmNldHlwZQ==');

@$core.Deprecated('Use changeTagsForResourceResponseDescriptor instead')
const ChangeTagsForResourceResponse$json = {
  '1': 'ChangeTagsForResourceResponse',
};

/// Descriptor for `ChangeTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeTagsForResourceResponseDescriptor =
    $convert.base64Decode('Ch1DaGFuZ2VUYWdzRm9yUmVzb3VyY2VSZXNwb25zZQ==');

@$core.Deprecated('Use cidrBlockInUseExceptionDescriptor instead')
const CidrBlockInUseException$json = {
  '1': 'CidrBlockInUseException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CidrBlockInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrBlockInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChdDaWRyQmxvY2tJblVzZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use cidrBlockSummaryDescriptor instead')
const CidrBlockSummary$json = {
  '1': 'CidrBlockSummary',
  '2': [
    {'1': 'cidrblock', '3': 56494185, '4': 1, '5': 9, '10': 'cidrblock'},
    {'1': 'locationname', '3': 158186566, '4': 1, '5': 9, '10': 'locationname'},
  ],
};

/// Descriptor for `CidrBlockSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrBlockSummaryDescriptor = $convert.base64Decode(
    'ChBDaWRyQmxvY2tTdW1tYXJ5Eh8KCWNpZHJibG9jaxjpkPgaIAEoCVIJY2lkcmJsb2NrEiUKDG'
    'xvY2F0aW9ubmFtZRjG+LZLIAEoCVIMbG9jYXRpb25uYW1l');

@$core.Deprecated('Use cidrCollectionDescriptor instead')
const CidrCollection$json = {
  '1': 'CidrCollection',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `CidrCollection`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrCollectionDescriptor = $convert.base64Decode(
    'Cg5DaWRyQ29sbGVjdGlvbhIUCgNhcm4YnZvtvwEgASgJUgNhcm4SEgoCaWQYgfKitwEgASgJUg'
    'JpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhwKB3ZlcnNpb24YuKq37gEgASgDUgd2ZXJzaW9u');

@$core.Deprecated('Use cidrCollectionAlreadyExistsExceptionDescriptor instead')
const CidrCollectionAlreadyExistsException$json = {
  '1': 'CidrCollectionAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CidrCollectionAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrCollectionAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'CiRDaWRyQ29sbGVjdGlvbkFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use cidrCollectionChangeDescriptor instead')
const CidrCollectionChange$json = {
  '1': 'CidrCollectionChange',
  '2': [
    {
      '1': 'action',
      '3': 175614240,
      '4': 1,
      '5': 14,
      '6': '.route53.CidrCollectionChangeAction',
      '10': 'action'
    },
    {'1': 'cidrlist', '3': 288964704, '4': 3, '5': 9, '10': 'cidrlist'},
    {'1': 'locationname', '3': 158186566, '4': 1, '5': 9, '10': 'locationname'},
  ],
};

/// Descriptor for `CidrCollectionChange`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrCollectionChangeDescriptor = $convert.base64Decode(
    'ChRDaWRyQ29sbGVjdGlvbkNoYW5nZRI+CgZhY3Rpb24YoNLeUyABKA4yIy5yb3V0ZTUzLkNpZH'
    'JDb2xsZWN0aW9uQ2hhbmdlQWN0aW9uUgZhY3Rpb24SHgoIY2lkcmxpc3QY4IDliQEgAygJUghj'
    'aWRybGlzdBIlCgxsb2NhdGlvbm5hbWUYxvi2SyABKAlSDGxvY2F0aW9ubmFtZQ==');

@$core.Deprecated('Use cidrCollectionInUseExceptionDescriptor instead')
const CidrCollectionInUseException$json = {
  '1': 'CidrCollectionInUseException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CidrCollectionInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrCollectionInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChxDaWRyQ29sbGVjdGlvbkluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core
    .Deprecated('Use cidrCollectionVersionMismatchExceptionDescriptor instead')
const CidrCollectionVersionMismatchException$json = {
  '1': 'CidrCollectionVersionMismatchException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CidrCollectionVersionMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrCollectionVersionMismatchExceptionDescriptor =
    $convert.base64Decode(
        'CiZDaWRyQ29sbGVjdGlvblZlcnNpb25NaXNtYXRjaEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3'
        'AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use cidrRoutingConfigDescriptor instead')
const CidrRoutingConfig$json = {
  '1': 'CidrRoutingConfig',
  '2': [
    {'1': 'collectionid', '3': 128052453, '4': 1, '5': 9, '10': 'collectionid'},
    {'1': 'locationname', '3': 158186566, '4': 1, '5': 9, '10': 'locationname'},
  ],
};

/// Descriptor for `CidrRoutingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cidrRoutingConfigDescriptor = $convert.base64Decode(
    'ChFDaWRyUm91dGluZ0NvbmZpZxIlCgxjb2xsZWN0aW9uaWQY5dmHPSABKAlSDGNvbGxlY3Rpb2'
    '5pZBIlCgxsb2NhdGlvbm5hbWUYxvi2SyABKAlSDGxvY2F0aW9ubmFtZQ==');

@$core.Deprecated('Use cloudWatchAlarmConfigurationDescriptor instead')
const CloudWatchAlarmConfiguration$json = {
  '1': 'CloudWatchAlarmConfiguration',
  '2': [
    {
      '1': 'comparisonoperator',
      '3': 7059861,
      '4': 1,
      '5': 14,
      '6': '.route53.ComparisonOperator',
      '10': 'comparisonoperator'
    },
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.route53.Dimension',
      '10': 'dimensions'
    },
    {
      '1': 'evaluationperiods',
      '3': 214794856,
      '4': 1,
      '5': 5,
      '10': 'evaluationperiods'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {
      '1': 'statistic',
      '3': 67293470,
      '4': 1,
      '5': 14,
      '6': '.route53.Statistic',
      '10': 'statistic'
    },
    {'1': 'threshold', '3': 394884505, '4': 1, '5': 1, '10': 'threshold'},
  ],
};

/// Descriptor for `CloudWatchAlarmConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchAlarmConfigurationDescriptor = $convert.base64Decode(
    'ChxDbG91ZFdhdGNoQWxhcm1Db25maWd1cmF0aW9uEk4KEmNvbXBhcmlzb25vcGVyYXRvchiV86'
    '4DIAEoDjIbLnJvdXRlNTMuQ29tcGFyaXNvbk9wZXJhdG9yUhJjb21wYXJpc29ub3BlcmF0b3IS'
    'NgoKZGltZW5zaW9ucxjRm9/cASADKAsyEi5yb3V0ZTUzLkRpbWVuc2lvblIKZGltZW5zaW9ucx'
    'IvChFldmFsdWF0aW9ucGVyaW9kcxjohLZmIAEoBVIRZXZhbHVhdGlvbnBlcmlvZHMSIQoKbWV0'
    'cmljbmFtZRj7vtoyIAEoCVIKbWV0cmljbmFtZRIgCgluYW1lc3BhY2UYwYS5qQEgASgJUgluYW'
    '1lc3BhY2USGQoGcGVyaW9kGKWIkjkgASgFUgZwZXJpb2QSMwoJc3RhdGlzdGljGJ6iiyAgASgO'
    'MhIucm91dGU1My5TdGF0aXN0aWNSCXN0YXRpc3RpYxIgCgl0aHJlc2hvbGQYmeulvAEgASgBUg'
    'l0aHJlc2hvbGQ=');

@$core.Deprecated('Use collectionSummaryDescriptor instead')
const CollectionSummary$json = {
  '1': 'CollectionSummary',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `CollectionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List collectionSummaryDescriptor = $convert.base64Decode(
    'ChFDb2xsZWN0aW9uU3VtbWFyeRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SEgoCaWQYgfKitwEgAS'
    'gJUgJpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhwKB3ZlcnNpb24YuKq37gEgASgDUgd2ZXJz'
    'aW9u');

@$core.Deprecated('Use concurrentModificationDescriptor instead')
const ConcurrentModification$json = {
  '1': 'ConcurrentModification',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConcurrentModification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrentModificationDescriptor =
    $convert.base64Decode(
        'ChZDb25jdXJyZW50TW9kaWZpY2F0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use conflictingDomainExistsDescriptor instead')
const ConflictingDomainExists$json = {
  '1': 'ConflictingDomainExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictingDomainExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictingDomainExistsDescriptor =
    $convert.base64Decode(
        'ChdDb25mbGljdGluZ0RvbWFpbkV4aXN0cxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use conflictingTypesDescriptor instead')
const ConflictingTypes$json = {
  '1': 'ConflictingTypes',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictingTypes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictingTypesDescriptor = $convert.base64Decode(
    'ChBDb25mbGljdGluZ1R5cGVzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use coordinatesDescriptor instead')
const Coordinates$json = {
  '1': 'Coordinates',
  '2': [
    {'1': 'latitude', '3': 226378086, '4': 1, '5': 9, '10': 'latitude'},
    {'1': 'longitude', '3': 119104751, '4': 1, '5': 9, '10': 'longitude'},
  ],
};

/// Descriptor for `Coordinates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List coordinatesDescriptor = $convert.base64Decode(
    'CgtDb29yZGluYXRlcxIdCghsYXRpdHVkZRjmgvlrIAEoCVIIbGF0aXR1ZGUSHwoJbG9uZ2l0dW'
    'RlGO/J5TggASgJUglsb25naXR1ZGU=');

@$core.Deprecated('Use createCidrCollectionRequestDescriptor instead')
const CreateCidrCollectionRequest$json = {
  '1': 'CreateCidrCollectionRequest',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreateCidrCollectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCidrCollectionRequestDescriptor =
    $convert.base64Decode(
        'ChtDcmVhdGVDaWRyQ29sbGVjdGlvblJlcXVlc3QSKwoPY2FsbGVycmVmZXJlbmNlGJiZjUggAS'
        'gJUg9jYWxsZXJyZWZlcmVuY2USFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use createCidrCollectionResponseDescriptor instead')
const CreateCidrCollectionResponse$json = {
  '1': 'CreateCidrCollectionResponse',
  '2': [
    {
      '1': 'collection',
      '3': 114272434,
      '4': 1,
      '5': 11,
      '6': '.route53.CidrCollection',
      '10': 'collection'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateCidrCollectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCidrCollectionResponseDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVDaWRyQ29sbGVjdGlvblJlc3BvbnNlEjoKCmNvbGxlY3Rpb24YstG+NiABKAsyFy'
        '5yb3V0ZTUzLkNpZHJDb2xsZWN0aW9uUgpjb2xsZWN0aW9uEh4KCGxvY2F0aW9uGMebgt4BIAEo'
        'CVIIbG9jYXRpb24=');

@$core.Deprecated('Use createHealthCheckRequestDescriptor instead')
const CreateHealthCheckRequest$json = {
  '1': 'CreateHealthCheckRequest',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {
      '1': 'healthcheckconfig',
      '3': 83111204,
      '4': 1,
      '5': 11,
      '6': '.route53.HealthCheckConfig',
      '10': 'healthcheckconfig'
    },
  ],
};

/// Descriptor for `CreateHealthCheckRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createHealthCheckRequestDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVIZWFsdGhDaGVja1JlcXVlc3QSKwoPY2FsbGVycmVmZXJlbmNlGJiZjUggASgJUg'
    '9jYWxsZXJyZWZlcmVuY2USSwoRaGVhbHRoY2hlY2tjb25maWcYpNrQJyABKAsyGi5yb3V0ZTUz'
    'LkhlYWx0aENoZWNrQ29uZmlnUhFoZWFsdGhjaGVja2NvbmZpZw==');

@$core.Deprecated('Use createHealthCheckResponseDescriptor instead')
const CreateHealthCheckResponse$json = {
  '1': 'CreateHealthCheckResponse',
  '2': [
    {
      '1': 'healthcheck',
      '3': 377540610,
      '4': 1,
      '5': 11,
      '6': '.route53.HealthCheck',
      '10': 'healthcheck'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateHealthCheckResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createHealthCheckResponseDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVIZWFsdGhDaGVja1Jlc3BvbnNlEjoKC2hlYWx0aGNoZWNrGIKgg7QBIAEoCzIULn'
    'JvdXRlNTMuSGVhbHRoQ2hlY2tSC2hlYWx0aGNoZWNrEh4KCGxvY2F0aW9uGMebgt4BIAEoCVII'
    'bG9jYXRpb24=');

@$core.Deprecated('Use createHostedZoneRequestDescriptor instead')
const CreateHostedZoneRequest$json = {
  '1': 'CreateHostedZoneRequest',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {
      '1': 'delegationsetid',
      '3': 307328801,
      '4': 1,
      '5': 9,
      '10': 'delegationsetid'
    },
    {
      '1': 'hostedzoneconfig',
      '3': 881519,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZoneConfig',
      '10': 'hostedzoneconfig'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'vpc',
      '3': 506158953,
      '4': 1,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpc'
    },
  ],
};

/// Descriptor for `CreateHostedZoneRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createHostedZoneRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVIb3N0ZWRab25lUmVxdWVzdBIrCg9jYWxsZXJyZWZlcmVuY2UYmJmNSCABKAlSD2'
    'NhbGxlcnJlZmVyZW5jZRIsCg9kZWxlZ2F0aW9uc2V0aWQYoe7FkgEgASgJUg9kZWxlZ2F0aW9u'
    'c2V0aWQSRwoQaG9zdGVkem9uZWNvbmZpZxjv5jUgASgLMhkucm91dGU1My5Ib3N0ZWRab25lQ2'
    '9uZmlnUhBob3N0ZWR6b25lY29uZmlnEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSIgoDdnBjGOm+'
    'rfEBIAEoCzIMLnJvdXRlNTMuVlBDUgN2cGM=');

@$core.Deprecated('Use createHostedZoneResponseDescriptor instead')
const CreateHostedZoneResponse$json = {
  '1': 'CreateHostedZoneResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
    {
      '1': 'delegationset',
      '3': 190771750,
      '4': 1,
      '5': 11,
      '6': '.route53.DelegationSet',
      '10': 'delegationset'
    },
    {
      '1': 'hostedzone',
      '3': 465689249,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZone',
      '10': 'hostedzone'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'vpc',
      '3': 506158953,
      '4': 1,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpc'
    },
  ],
};

/// Descriptor for `CreateHostedZoneResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createHostedZoneResponseDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVIb3N0ZWRab25lUmVzcG9uc2USNwoKY2hhbmdlaW5mbxjutYvQASABKAsyEy5yb3'
    'V0ZTUzLkNoYW5nZUluZm9SCmNoYW5nZWluZm8SPwoNZGVsZWdhdGlvbnNldBim5PtaIAEoCzIW'
    'LnJvdXRlNTMuRGVsZWdhdGlvblNldFINZGVsZWdhdGlvbnNldBI3Cgpob3N0ZWR6b25lGKG1h9'
    '4BIAEoCzITLnJvdXRlNTMuSG9zdGVkWm9uZVIKaG9zdGVkem9uZRIeCghsb2NhdGlvbhjHm4Le'
    'ASABKAlSCGxvY2F0aW9uEiIKA3ZwYxjpvq3xASABKAsyDC5yb3V0ZTUzLlZQQ1IDdnBj');

@$core.Deprecated('Use createKeySigningKeyRequestDescriptor instead')
const CreateKeySigningKeyRequest$json = {
  '1': 'CreateKeySigningKeyRequest',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'keymanagementservicearn',
      '3': 360687398,
      '4': 1,
      '5': 9,
      '10': 'keymanagementservicearn'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `CreateKeySigningKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeySigningKeyRequestDescriptor = $convert.base64Decode(
    'ChpDcmVhdGVLZXlTaWduaW5nS2V5UmVxdWVzdBIrCg9jYWxsZXJyZWZlcmVuY2UYmJmNSCABKA'
    'lSD2NhbGxlcnJlZmVyZW5jZRImCgxob3N0ZWR6b25laWQY/s6epQEgASgJUgxob3N0ZWR6b25l'
    'aWQSPAoXa2V5bWFuYWdlbWVudHNlcnZpY2Vhcm4Yps7+qwEgASgJUhdrZXltYW5hZ2VtZW50c2'
    'VydmljZWFybhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhkKBnN0YXR1cxiQ5PsCIAEoCVIGc3Rh'
    'dHVz');

@$core.Deprecated('Use createKeySigningKeyResponseDescriptor instead')
const CreateKeySigningKeyResponse$json = {
  '1': 'CreateKeySigningKeyResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
    {
      '1': 'keysigningkey',
      '3': 91098059,
      '4': 1,
      '5': 11,
      '6': '.route53.KeySigningKey',
      '10': 'keysigningkey'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateKeySigningKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeySigningKeyResponseDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVLZXlTaWduaW5nS2V5UmVzcG9uc2USNwoKY2hhbmdlaW5mbxjutYvQASABKAsyEy'
    '5yb3V0ZTUzLkNoYW5nZUluZm9SCmNoYW5nZWluZm8SPwoNa2V5c2lnbmluZ2tleRjLl7grIAEo'
    'CzIWLnJvdXRlNTMuS2V5U2lnbmluZ0tleVINa2V5c2lnbmluZ2tleRIeCghsb2NhdGlvbhjHm4'
    'LeASABKAlSCGxvY2F0aW9u');

@$core.Deprecated('Use createQueryLoggingConfigRequestDescriptor instead')
const CreateQueryLoggingConfigRequest$json = {
  '1': 'CreateQueryLoggingConfigRequest',
  '2': [
    {
      '1': 'cloudwatchlogsloggrouparn',
      '3': 519613355,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsloggrouparn'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `CreateQueryLoggingConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createQueryLoggingConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch9DcmVhdGVRdWVyeUxvZ2dpbmdDb25maWdSZXF1ZXN0EkAKGWNsb3Vkd2F0Y2hsb2dzbG9nZ3'
        'JvdXBhcm4Yq9fi9wEgASgJUhljbG91ZHdhdGNobG9nc2xvZ2dyb3VwYXJuEiYKDGhvc3RlZHpv'
        'bmVpZBj+zp6lASABKAlSDGhvc3RlZHpvbmVpZA==');

@$core.Deprecated('Use createQueryLoggingConfigResponseDescriptor instead')
const CreateQueryLoggingConfigResponse$json = {
  '1': 'CreateQueryLoggingConfigResponse',
  '2': [
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'queryloggingconfig',
      '3': 492751675,
      '4': 1,
      '5': 11,
      '6': '.route53.QueryLoggingConfig',
      '10': 'queryloggingconfig'
    },
  ],
};

/// Descriptor for `CreateQueryLoggingConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createQueryLoggingConfigResponseDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVRdWVyeUxvZ2dpbmdDb25maWdSZXNwb25zZRIeCghsb2NhdGlvbhjHm4LeASABKA'
        'lSCGxvY2F0aW9uEk8KEnF1ZXJ5bG9nZ2luZ2NvbmZpZxi7lvvqASABKAsyGy5yb3V0ZTUzLlF1'
        'ZXJ5TG9nZ2luZ0NvbmZpZ1IScXVlcnlsb2dnaW5nY29uZmln');

@$core.Deprecated('Use createReusableDelegationSetRequestDescriptor instead')
const CreateReusableDelegationSetRequest$json = {
  '1': 'CreateReusableDelegationSetRequest',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `CreateReusableDelegationSetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createReusableDelegationSetRequestDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVSZXVzYWJsZURlbGVnYXRpb25TZXRSZXF1ZXN0EisKD2NhbGxlcnJlZmVyZW5jZR'
        'iYmY1IIAEoCVIPY2FsbGVycmVmZXJlbmNlEiYKDGhvc3RlZHpvbmVpZBj+zp6lASABKAlSDGhv'
        'c3RlZHpvbmVpZA==');

@$core.Deprecated('Use createReusableDelegationSetResponseDescriptor instead')
const CreateReusableDelegationSetResponse$json = {
  '1': 'CreateReusableDelegationSetResponse',
  '2': [
    {
      '1': 'delegationset',
      '3': 190771750,
      '4': 1,
      '5': 11,
      '6': '.route53.DelegationSet',
      '10': 'delegationset'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateReusableDelegationSetResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createReusableDelegationSetResponseDescriptor =
    $convert.base64Decode(
        'CiNDcmVhdGVSZXVzYWJsZURlbGVnYXRpb25TZXRSZXNwb25zZRI/Cg1kZWxlZ2F0aW9uc2V0GK'
        'bk+1ogASgLMhYucm91dGU1My5EZWxlZ2F0aW9uU2V0Ug1kZWxlZ2F0aW9uc2V0Eh4KCGxvY2F0'
        'aW9uGMebgt4BIAEoCVIIbG9jYXRpb24=');

@$core.Deprecated('Use createTrafficPolicyInstanceRequestDescriptor instead')
const CreateTrafficPolicyInstanceRequest$json = {
  '1': 'CreateTrafficPolicyInstanceRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'ttl', '3': 526904700, '4': 1, '5': 3, '10': 'ttl'},
    {
      '1': 'trafficpolicyid',
      '3': 40235222,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyid'
    },
    {
      '1': 'trafficpolicyversion',
      '3': 479078485,
      '4': 1,
      '5': 5,
      '10': 'trafficpolicyversion'
    },
  ],
};

/// Descriptor for `CreateTrafficPolicyInstanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrafficPolicyInstanceRequestDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVUcmFmZmljUG9saWN5SW5zdGFuY2VSZXF1ZXN0EiYKDGhvc3RlZHpvbmVpZBj+zp'
        '6lASABKAlSDGhvc3RlZHpvbmVpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhQKA3R0bBj82p/7'
        'ASABKANSA3R0bBIrCg90cmFmZmljcG9saWN5aWQY1uGXEyABKAlSD3RyYWZmaWNwb2xpY3lpZB'
        'I2ChR0cmFmZmljcG9saWN5dmVyc2lvbhjV0LjkASABKAVSFHRyYWZmaWNwb2xpY3l2ZXJzaW9u');

@$core.Deprecated('Use createTrafficPolicyInstanceResponseDescriptor instead')
const CreateTrafficPolicyInstanceResponse$json = {
  '1': 'CreateTrafficPolicyInstanceResponse',
  '2': [
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'trafficpolicyinstance',
      '3': 205651476,
      '4': 1,
      '5': 11,
      '6': '.route53.TrafficPolicyInstance',
      '10': 'trafficpolicyinstance'
    },
  ],
};

/// Descriptor for `CreateTrafficPolicyInstanceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrafficPolicyInstanceResponseDescriptor =
    $convert.base64Decode(
        'CiNDcmVhdGVUcmFmZmljUG9saWN5SW5zdGFuY2VSZXNwb25zZRIeCghsb2NhdGlvbhjHm4LeAS'
        'ABKAlSCGxvY2F0aW9uElcKFXRyYWZmaWNwb2xpY3lpbnN0YW5jZRiU/IdiIAEoCzIeLnJvdXRl'
        'NTMuVHJhZmZpY1BvbGljeUluc3RhbmNlUhV0cmFmZmljcG9saWN5aW5zdGFuY2U=');

@$core.Deprecated('Use createTrafficPolicyRequestDescriptor instead')
const CreateTrafficPolicyRequest$json = {
  '1': 'CreateTrafficPolicyRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'document', '3': 407108341, '4': 1, '5': 9, '10': 'document'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreateTrafficPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrafficPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChpDcmVhdGVUcmFmZmljUG9saWN5UmVxdWVzdBIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbW'
        'VudBIeCghkb2N1bWVudBj19Y/CASABKAlSCGRvY3VtZW50EhUKBG5hbWUYh+aBfyABKAlSBG5h'
        'bWU=');

@$core.Deprecated('Use createTrafficPolicyResponseDescriptor instead')
const CreateTrafficPolicyResponse$json = {
  '1': 'CreateTrafficPolicyResponse',
  '2': [
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'trafficpolicy',
      '3': 154595657,
      '4': 1,
      '5': 11,
      '6': '.route53.TrafficPolicy',
      '10': 'trafficpolicy'
    },
  ],
};

/// Descriptor for `CreateTrafficPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrafficPolicyResponseDescriptor =
    $convert.base64Decode(
        'ChtDcmVhdGVUcmFmZmljUG9saWN5UmVzcG9uc2USHgoIbG9jYXRpb24Yx5uC3gEgASgJUghsb2'
        'NhdGlvbhI/Cg10cmFmZmljcG9saWN5GMni20kgASgLMhYucm91dGU1My5UcmFmZmljUG9saWN5'
        'Ug10cmFmZmljcG9saWN5');

@$core.Deprecated('Use createTrafficPolicyVersionRequestDescriptor instead')
const CreateTrafficPolicyVersionRequest$json = {
  '1': 'CreateTrafficPolicyVersionRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'document', '3': 407108341, '4': 1, '5': 9, '10': 'document'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `CreateTrafficPolicyVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrafficPolicyVersionRequestDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVUcmFmZmljUG9saWN5VmVyc2lvblJlcXVlc3QSHAoHY29tbWVudBj/v77CASABKA'
        'lSB2NvbW1lbnQSHgoIZG9jdW1lbnQY9fWPwgEgASgJUghkb2N1bWVudBISCgJpZBiB8qK3ASAB'
        'KAlSAmlk');

@$core.Deprecated('Use createTrafficPolicyVersionResponseDescriptor instead')
const CreateTrafficPolicyVersionResponse$json = {
  '1': 'CreateTrafficPolicyVersionResponse',
  '2': [
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'trafficpolicy',
      '3': 154595657,
      '4': 1,
      '5': 11,
      '6': '.route53.TrafficPolicy',
      '10': 'trafficpolicy'
    },
  ],
};

/// Descriptor for `CreateTrafficPolicyVersionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrafficPolicyVersionResponseDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVUcmFmZmljUG9saWN5VmVyc2lvblJlc3BvbnNlEh4KCGxvY2F0aW9uGMebgt4BIA'
        'EoCVIIbG9jYXRpb24SPwoNdHJhZmZpY3BvbGljeRjJ4ttJIAEoCzIWLnJvdXRlNTMuVHJhZmZp'
        'Y1BvbGljeVINdHJhZmZpY3BvbGljeQ==');

@$core.Deprecated(
    'Use createVPCAssociationAuthorizationRequestDescriptor instead')
const CreateVPCAssociationAuthorizationRequest$json = {
  '1': 'CreateVPCAssociationAuthorizationRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'vpc',
      '3': 506158953,
      '4': 1,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpc'
    },
  ],
};

/// Descriptor for `CreateVPCAssociationAuthorizationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createVPCAssociationAuthorizationRequestDescriptor =
    $convert.base64Decode(
        'CihDcmVhdGVWUENBc3NvY2lhdGlvbkF1dGhvcml6YXRpb25SZXF1ZXN0EiYKDGhvc3RlZHpvbm'
        'VpZBj+zp6lASABKAlSDGhvc3RlZHpvbmVpZBIiCgN2cGMY6b6t8QEgASgLMgwucm91dGU1My5W'
        'UENSA3ZwYw==');

@$core.Deprecated(
    'Use createVPCAssociationAuthorizationResponseDescriptor instead')
const CreateVPCAssociationAuthorizationResponse$json = {
  '1': 'CreateVPCAssociationAuthorizationResponse',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'vpc',
      '3': 506158953,
      '4': 1,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpc'
    },
  ],
};

/// Descriptor for `CreateVPCAssociationAuthorizationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createVPCAssociationAuthorizationResponseDescriptor = $convert.base64Decode(
        'CilDcmVhdGVWUENBc3NvY2lhdGlvbkF1dGhvcml6YXRpb25SZXNwb25zZRImCgxob3N0ZWR6b2'
        '5laWQY/s6epQEgASgJUgxob3N0ZWR6b25laWQSIgoDdnBjGOm+rfEBIAEoCzIMLnJvdXRlNTMu'
        'VlBDUgN2cGM=');

@$core.Deprecated('Use dNSSECNotFoundDescriptor instead')
const DNSSECNotFound$json = {
  '1': 'DNSSECNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DNSSECNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dNSSECNotFoundDescriptor = $convert.base64Decode(
    'Cg5ETlNTRUNOb3RGb3VuZBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use dNSSECStatusDescriptor instead')
const DNSSECStatus$json = {
  '1': 'DNSSECStatus',
  '2': [
    {
      '1': 'servesignature',
      '3': 466239599,
      '4': 1,
      '5': 9,
      '10': 'servesignature'
    },
    {
      '1': 'statusmessage',
      '3': 72590095,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
  ],
};

/// Descriptor for `DNSSECStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dNSSECStatusDescriptor = $convert.base64Decode(
    'CgxETlNTRUNTdGF0dXMSKgoOc2VydmVzaWduYXR1cmUY74Cp3gEgASgJUg5zZXJ2ZXNpZ25hdH'
    'VyZRInCg1zdGF0dXNtZXNzYWdlGI/GziIgASgJUg1zdGF0dXNtZXNzYWdl');

@$core.Deprecated('Use deactivateKeySigningKeyRequestDescriptor instead')
const DeactivateKeySigningKeyRequest$json = {
  '1': 'DeactivateKeySigningKeyRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeactivateKeySigningKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deactivateKeySigningKeyRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWFjdGl2YXRlS2V5U2lnbmluZ0tleVJlcXVlc3QSJgoMaG9zdGVkem9uZWlkGP7OnqUBIA'
        'EoCVIMaG9zdGVkem9uZWlkEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deactivateKeySigningKeyResponseDescriptor instead')
const DeactivateKeySigningKeyResponse$json = {
  '1': 'DeactivateKeySigningKeyResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `DeactivateKeySigningKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deactivateKeySigningKeyResponseDescriptor =
    $convert.base64Decode(
        'Ch9EZWFjdGl2YXRlS2V5U2lnbmluZ0tleVJlc3BvbnNlEjcKCmNoYW5nZWluZm8Y7rWL0AEgAS'
        'gLMhMucm91dGU1My5DaGFuZ2VJbmZvUgpjaGFuZ2VpbmZv');

@$core.Deprecated('Use delegationSetDescriptor instead')
const DelegationSet$json = {
  '1': 'DelegationSet',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'nameservers', '3': 340971511, '4': 3, '5': 9, '10': 'nameservers'},
  ],
};

/// Descriptor for `DelegationSet`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationSetDescriptor = $convert.base64Decode(
    'Cg1EZWxlZ2F0aW9uU2V0EisKD2NhbGxlcnJlZmVyZW5jZRiYmY1IIAEoCVIPY2FsbGVycmVmZX'
    'JlbmNlEhIKAmlkGIHyorcBIAEoCVICaWQSJAoLbmFtZXNlcnZlcnMY95/LogEgAygJUgtuYW1l'
    'c2VydmVycw==');

@$core.Deprecated('Use delegationSetAlreadyCreatedDescriptor instead')
const DelegationSetAlreadyCreated$json = {
  '1': 'DelegationSetAlreadyCreated',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DelegationSetAlreadyCreated`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationSetAlreadyCreatedDescriptor =
    $convert.base64Decode(
        'ChtEZWxlZ2F0aW9uU2V0QWxyZWFkeUNyZWF0ZWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use delegationSetAlreadyReusableDescriptor instead')
const DelegationSetAlreadyReusable$json = {
  '1': 'DelegationSetAlreadyReusable',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DelegationSetAlreadyReusable`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationSetAlreadyReusableDescriptor =
    $convert.base64Decode(
        'ChxEZWxlZ2F0aW9uU2V0QWxyZWFkeVJldXNhYmxlEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use delegationSetInUseDescriptor instead')
const DelegationSetInUse$json = {
  '1': 'DelegationSetInUse',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DelegationSetInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationSetInUseDescriptor =
    $convert.base64Decode(
        'ChJEZWxlZ2F0aW9uU2V0SW5Vc2USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use delegationSetNotAvailableDescriptor instead')
const DelegationSetNotAvailable$json = {
  '1': 'DelegationSetNotAvailable',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DelegationSetNotAvailable`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationSetNotAvailableDescriptor =
    $convert.base64Decode(
        'ChlEZWxlZ2F0aW9uU2V0Tm90QXZhaWxhYmxlEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use delegationSetNotReusableDescriptor instead')
const DelegationSetNotReusable$json = {
  '1': 'DelegationSetNotReusable',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DelegationSetNotReusable`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List delegationSetNotReusableDescriptor =
    $convert.base64Decode(
        'ChhEZWxlZ2F0aW9uU2V0Tm90UmV1c2FibGUSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use deleteCidrCollectionRequestDescriptor instead')
const DeleteCidrCollectionRequest$json = {
  '1': 'DeleteCidrCollectionRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteCidrCollectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCidrCollectionRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVDaWRyQ29sbGVjdGlvblJlcXVlc3QSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use deleteCidrCollectionResponseDescriptor instead')
const DeleteCidrCollectionResponse$json = {
  '1': 'DeleteCidrCollectionResponse',
};

/// Descriptor for `DeleteCidrCollectionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCidrCollectionResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVDaWRyQ29sbGVjdGlvblJlc3BvbnNl');

@$core.Deprecated('Use deleteHealthCheckRequestDescriptor instead')
const DeleteHealthCheckRequest$json = {
  '1': 'DeleteHealthCheckRequest',
  '2': [
    {
      '1': 'healthcheckid',
      '3': 312971637,
      '4': 1,
      '5': 9,
      '10': 'healthcheckid'
    },
  ],
};

/// Descriptor for `DeleteHealthCheckRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteHealthCheckRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVIZWFsdGhDaGVja1JlcXVlc3QSKAoNaGVhbHRoY2hlY2tpZBj1op6VASABKAlSDW'
        'hlYWx0aGNoZWNraWQ=');

@$core.Deprecated('Use deleteHealthCheckResponseDescriptor instead')
const DeleteHealthCheckResponse$json = {
  '1': 'DeleteHealthCheckResponse',
};

/// Descriptor for `DeleteHealthCheckResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteHealthCheckResponseDescriptor =
    $convert.base64Decode('ChlEZWxldGVIZWFsdGhDaGVja1Jlc3BvbnNl');

@$core.Deprecated('Use deleteHostedZoneRequestDescriptor instead')
const DeleteHostedZoneRequest$json = {
  '1': 'DeleteHostedZoneRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteHostedZoneRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteHostedZoneRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVIb3N0ZWRab25lUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use deleteHostedZoneResponseDescriptor instead')
const DeleteHostedZoneResponse$json = {
  '1': 'DeleteHostedZoneResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `DeleteHostedZoneResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteHostedZoneResponseDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVIb3N0ZWRab25lUmVzcG9uc2USNwoKY2hhbmdlaW5mbxjutYvQASABKAsyEy5yb3'
        'V0ZTUzLkNoYW5nZUluZm9SCmNoYW5nZWluZm8=');

@$core.Deprecated('Use deleteKeySigningKeyRequestDescriptor instead')
const DeleteKeySigningKeyRequest$json = {
  '1': 'DeleteKeySigningKeyRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteKeySigningKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteKeySigningKeyRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVLZXlTaWduaW5nS2V5UmVxdWVzdBImCgxob3N0ZWR6b25laWQY/s6epQEgASgJUg'
        'xob3N0ZWR6b25laWQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use deleteKeySigningKeyResponseDescriptor instead')
const DeleteKeySigningKeyResponse$json = {
  '1': 'DeleteKeySigningKeyResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `DeleteKeySigningKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteKeySigningKeyResponseDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVLZXlTaWduaW5nS2V5UmVzcG9uc2USNwoKY2hhbmdlaW5mbxjutYvQASABKAsyEy'
        '5yb3V0ZTUzLkNoYW5nZUluZm9SCmNoYW5nZWluZm8=');

@$core.Deprecated('Use deleteQueryLoggingConfigRequestDescriptor instead')
const DeleteQueryLoggingConfigRequest$json = {
  '1': 'DeleteQueryLoggingConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteQueryLoggingConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteQueryLoggingConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVRdWVyeUxvZ2dpbmdDb25maWdSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use deleteQueryLoggingConfigResponseDescriptor instead')
const DeleteQueryLoggingConfigResponse$json = {
  '1': 'DeleteQueryLoggingConfigResponse',
};

/// Descriptor for `DeleteQueryLoggingConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteQueryLoggingConfigResponseDescriptor =
    $convert.base64Decode('CiBEZWxldGVRdWVyeUxvZ2dpbmdDb25maWdSZXNwb25zZQ==');

@$core.Deprecated('Use deleteReusableDelegationSetRequestDescriptor instead')
const DeleteReusableDelegationSetRequest$json = {
  '1': 'DeleteReusableDelegationSetRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteReusableDelegationSetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteReusableDelegationSetRequestDescriptor =
    $convert.base64Decode(
        'CiJEZWxldGVSZXVzYWJsZURlbGVnYXRpb25TZXRSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'Q=');

@$core.Deprecated('Use deleteReusableDelegationSetResponseDescriptor instead')
const DeleteReusableDelegationSetResponse$json = {
  '1': 'DeleteReusableDelegationSetResponse',
};

/// Descriptor for `DeleteReusableDelegationSetResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteReusableDelegationSetResponseDescriptor =
    $convert
        .base64Decode('CiNEZWxldGVSZXVzYWJsZURlbGVnYXRpb25TZXRSZXNwb25zZQ==');

@$core.Deprecated('Use deleteTrafficPolicyInstanceRequestDescriptor instead')
const DeleteTrafficPolicyInstanceRequest$json = {
  '1': 'DeleteTrafficPolicyInstanceRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteTrafficPolicyInstanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTrafficPolicyInstanceRequestDescriptor =
    $convert.base64Decode(
        'CiJEZWxldGVUcmFmZmljUG9saWN5SW5zdGFuY2VSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'Q=');

@$core.Deprecated('Use deleteTrafficPolicyInstanceResponseDescriptor instead')
const DeleteTrafficPolicyInstanceResponse$json = {
  '1': 'DeleteTrafficPolicyInstanceResponse',
};

/// Descriptor for `DeleteTrafficPolicyInstanceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTrafficPolicyInstanceResponseDescriptor =
    $convert
        .base64Decode('CiNEZWxldGVUcmFmZmljUG9saWN5SW5zdGFuY2VSZXNwb25zZQ==');

@$core.Deprecated('Use deleteTrafficPolicyRequestDescriptor instead')
const DeleteTrafficPolicyRequest$json = {
  '1': 'DeleteTrafficPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 5, '10': 'version'},
  ],
};

/// Descriptor for `DeleteTrafficPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTrafficPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVUcmFmZmljUG9saWN5UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEhwKB3Zlcn'
        'Npb24YuKq37gEgASgFUgd2ZXJzaW9u');

@$core.Deprecated('Use deleteTrafficPolicyResponseDescriptor instead')
const DeleteTrafficPolicyResponse$json = {
  '1': 'DeleteTrafficPolicyResponse',
};

/// Descriptor for `DeleteTrafficPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTrafficPolicyResponseDescriptor =
    $convert.base64Decode('ChtEZWxldGVUcmFmZmljUG9saWN5UmVzcG9uc2U=');

@$core.Deprecated(
    'Use deleteVPCAssociationAuthorizationRequestDescriptor instead')
const DeleteVPCAssociationAuthorizationRequest$json = {
  '1': 'DeleteVPCAssociationAuthorizationRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'vpc',
      '3': 506158953,
      '4': 1,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpc'
    },
  ],
};

/// Descriptor for `DeleteVPCAssociationAuthorizationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteVPCAssociationAuthorizationRequestDescriptor =
    $convert.base64Decode(
        'CihEZWxldGVWUENBc3NvY2lhdGlvbkF1dGhvcml6YXRpb25SZXF1ZXN0EiYKDGhvc3RlZHpvbm'
        'VpZBj+zp6lASABKAlSDGhvc3RlZHpvbmVpZBIiCgN2cGMY6b6t8QEgASgLMgwucm91dGU1My5W'
        'UENSA3ZwYw==');

@$core.Deprecated(
    'Use deleteVPCAssociationAuthorizationResponseDescriptor instead')
const DeleteVPCAssociationAuthorizationResponse$json = {
  '1': 'DeleteVPCAssociationAuthorizationResponse',
};

/// Descriptor for `DeleteVPCAssociationAuthorizationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteVPCAssociationAuthorizationResponseDescriptor = $convert.base64Decode(
        'CilEZWxldGVWUENBc3NvY2lhdGlvbkF1dGhvcml6YXRpb25SZXNwb25zZQ==');

@$core.Deprecated('Use dimensionDescriptor instead')
const Dimension$json = {
  '1': 'Dimension',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Dimension`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dimensionDescriptor = $convert.base64Decode(
    'CglEaW1lbnNpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIYCgV2YWx1ZRjr8p+KASABKAlSBX'
    'ZhbHVl');

@$core.Deprecated('Use disableHostedZoneDNSSECRequestDescriptor instead')
const DisableHostedZoneDNSSECRequest$json = {
  '1': 'DisableHostedZoneDNSSECRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `DisableHostedZoneDNSSECRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableHostedZoneDNSSECRequestDescriptor =
    $convert.base64Decode(
        'Ch5EaXNhYmxlSG9zdGVkWm9uZUROU1NFQ1JlcXVlc3QSJgoMaG9zdGVkem9uZWlkGP7OnqUBIA'
        'EoCVIMaG9zdGVkem9uZWlk');

@$core.Deprecated('Use disableHostedZoneDNSSECResponseDescriptor instead')
const DisableHostedZoneDNSSECResponse$json = {
  '1': 'DisableHostedZoneDNSSECResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `DisableHostedZoneDNSSECResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableHostedZoneDNSSECResponseDescriptor =
    $convert.base64Decode(
        'Ch9EaXNhYmxlSG9zdGVkWm9uZUROU1NFQ1Jlc3BvbnNlEjcKCmNoYW5nZWluZm8Y7rWL0AEgAS'
        'gLMhMucm91dGU1My5DaGFuZ2VJbmZvUgpjaGFuZ2VpbmZv');

@$core.Deprecated('Use disassociateVPCFromHostedZoneRequestDescriptor instead')
const DisassociateVPCFromHostedZoneRequest$json = {
  '1': 'DisassociateVPCFromHostedZoneRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'vpc',
      '3': 506158953,
      '4': 1,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpc'
    },
  ],
};

/// Descriptor for `DisassociateVPCFromHostedZoneRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disassociateVPCFromHostedZoneRequestDescriptor =
    $convert.base64Decode(
        'CiREaXNhc3NvY2lhdGVWUENGcm9tSG9zdGVkWm9uZVJlcXVlc3QSHAoHY29tbWVudBj/v77CAS'
        'ABKAlSB2NvbW1lbnQSJgoMaG9zdGVkem9uZWlkGP7OnqUBIAEoCVIMaG9zdGVkem9uZWlkEiIK'
        'A3ZwYxjpvq3xASABKAsyDC5yb3V0ZTUzLlZQQ1IDdnBj');

@$core.Deprecated('Use disassociateVPCFromHostedZoneResponseDescriptor instead')
const DisassociateVPCFromHostedZoneResponse$json = {
  '1': 'DisassociateVPCFromHostedZoneResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `DisassociateVPCFromHostedZoneResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disassociateVPCFromHostedZoneResponseDescriptor =
    $convert.base64Decode(
        'CiVEaXNhc3NvY2lhdGVWUENGcm9tSG9zdGVkWm9uZVJlc3BvbnNlEjcKCmNoYW5nZWluZm8Y7r'
        'WL0AEgASgLMhMucm91dGU1My5DaGFuZ2VJbmZvUgpjaGFuZ2VpbmZv');

@$core.Deprecated('Use enableHostedZoneDNSSECRequestDescriptor instead')
const EnableHostedZoneDNSSECRequest$json = {
  '1': 'EnableHostedZoneDNSSECRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `EnableHostedZoneDNSSECRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableHostedZoneDNSSECRequestDescriptor =
    $convert.base64Decode(
        'Ch1FbmFibGVIb3N0ZWRab25lRE5TU0VDUmVxdWVzdBImCgxob3N0ZWR6b25laWQY/s6epQEgAS'
        'gJUgxob3N0ZWR6b25laWQ=');

@$core.Deprecated('Use enableHostedZoneDNSSECResponseDescriptor instead')
const EnableHostedZoneDNSSECResponse$json = {
  '1': 'EnableHostedZoneDNSSECResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `EnableHostedZoneDNSSECResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableHostedZoneDNSSECResponseDescriptor =
    $convert.base64Decode(
        'Ch5FbmFibGVIb3N0ZWRab25lRE5TU0VDUmVzcG9uc2USNwoKY2hhbmdlaW5mbxjutYvQASABKA'
        'syEy5yb3V0ZTUzLkNoYW5nZUluZm9SCmNoYW5nZWluZm8=');

@$core.Deprecated('Use geoLocationDescriptor instead')
const GeoLocation$json = {
  '1': 'GeoLocation',
  '2': [
    {
      '1': 'continentcode',
      '3': 180194347,
      '4': 1,
      '5': 9,
      '10': 'continentcode'
    },
    {'1': 'countrycode', '3': 485287065, '4': 1, '5': 9, '10': 'countrycode'},
    {
      '1': 'subdivisioncode',
      '3': 444328268,
      '4': 1,
      '5': 9,
      '10': 'subdivisioncode'
    },
  ],
};

/// Descriptor for `GeoLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List geoLocationDescriptor = $convert.base64Decode(
    'CgtHZW9Mb2NhdGlvbhInCg1jb250aW5lbnRjb2RlGKuY9lUgASgJUg1jb250aW5lbnRjb2RlEi'
    'QKC2NvdW50cnljb2RlGJnJs+cBIAEoCVILY291bnRyeWNvZGUSLAoPc3ViZGl2aXNpb25jb2Rl'
    'GMzS79MBIAEoCVIPc3ViZGl2aXNpb25jb2Rl');

@$core.Deprecated('Use geoLocationDetailsDescriptor instead')
const GeoLocationDetails$json = {
  '1': 'GeoLocationDetails',
  '2': [
    {
      '1': 'continentcode',
      '3': 180194347,
      '4': 1,
      '5': 9,
      '10': 'continentcode'
    },
    {
      '1': 'continentname',
      '3': 203176953,
      '4': 1,
      '5': 9,
      '10': 'continentname'
    },
    {'1': 'countrycode', '3': 485287065, '4': 1, '5': 9, '10': 'countrycode'},
    {'1': 'countryname', '3': 16070667, '4': 1, '5': 9, '10': 'countryname'},
    {
      '1': 'subdivisioncode',
      '3': 444328268,
      '4': 1,
      '5': 9,
      '10': 'subdivisioncode'
    },
    {
      '1': 'subdivisionname',
      '3': 412782294,
      '4': 1,
      '5': 9,
      '10': 'subdivisionname'
    },
  ],
};

/// Descriptor for `GeoLocationDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List geoLocationDetailsDescriptor = $convert.base64Decode(
    'ChJHZW9Mb2NhdGlvbkRldGFpbHMSJwoNY29udGluZW50Y29kZRirmPZVIAEoCVINY29udGluZW'
    '50Y29kZRInCg1jb250aW5lbnRuYW1lGPn38GAgASgJUg1jb250aW5lbnRuYW1lEiQKC2NvdW50'
    'cnljb2RlGJnJs+cBIAEoCVILY291bnRyeWNvZGUSIwoLY291bnRyeW5hbWUYi/DUByABKAlSC2'
    'NvdW50cnluYW1lEiwKD3N1YmRpdmlzaW9uY29kZRjM0u/TASABKAlSD3N1YmRpdmlzaW9uY29k'
    'ZRIsCg9zdWJkaXZpc2lvbm5hbWUY1p3qxAEgASgJUg9zdWJkaXZpc2lvbm5hbWU=');

@$core.Deprecated('Use geoProximityLocationDescriptor instead')
const GeoProximityLocation$json = {
  '1': 'GeoProximityLocation',
  '2': [
    {'1': 'awsregion', '3': 245430451, '4': 1, '5': 9, '10': 'awsregion'},
    {'1': 'bias', '3': 60849893, '4': 1, '5': 5, '10': 'bias'},
    {
      '1': 'coordinates',
      '3': 231283401,
      '4': 1,
      '5': 11,
      '6': '.route53.Coordinates',
      '10': 'coordinates'
    },
    {
      '1': 'localzonegroup',
      '3': 60150354,
      '4': 1,
      '5': 9,
      '10': 'localzonegroup'
    },
  ],
};

/// Descriptor for `GeoProximityLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List geoProximityLocationDescriptor = $convert.base64Decode(
    'ChRHZW9Qcm94aW1pdHlMb2NhdGlvbhIfCglhd3NyZWdpb24Ys/GDdSABKAlSCWF3c3JlZ2lvbh'
    'IVCgRiaWFzGOX9gR0gASgFUgRiaWFzEjkKC2Nvb3JkaW5hdGVzGMm1pG4gASgLMhQucm91dGU1'
    'My5Db29yZGluYXRlc1ILY29vcmRpbmF0ZXMSKQoObG9jYWx6b25lZ3JvdXAY0qTXHCABKAlSDm'
    'xvY2Fsem9uZWdyb3Vw');

@$core.Deprecated('Use getAccountLimitRequestDescriptor instead')
const GetAccountLimitRequest$json = {
  '1': 'GetAccountLimitRequest',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.AccountLimitType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `GetAccountLimitRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountLimitRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRBY2NvdW50TGltaXRSZXF1ZXN0EjEKBHR5cGUY7qDXigEgASgOMhkucm91dGU1My5BY2'
        'NvdW50TGltaXRUeXBlUgR0eXBl');

@$core.Deprecated('Use getAccountLimitResponseDescriptor instead')
const GetAccountLimitResponse$json = {
  '1': 'GetAccountLimitResponse',
  '2': [
    {'1': 'count', '3': 31963285, '4': 1, '5': 3, '10': 'count'},
    {
      '1': 'limit',
      '3': 412502741,
      '4': 1,
      '5': 11,
      '6': '.route53.AccountLimit',
      '10': 'limit'
    },
  ],
};

/// Descriptor for `GetAccountLimitResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountLimitResponseDescriptor =
    $convert.base64Decode(
        'ChdHZXRBY2NvdW50TGltaXRSZXNwb25zZRIXCgVjb3VudBiV8Z4PIAEoA1IFY291bnQSLwoFbG'
        'ltaXQY1ZXZxAEgASgLMhUucm91dGU1My5BY2NvdW50TGltaXRSBWxpbWl0');

@$core.Deprecated('Use getChangeRequestDescriptor instead')
const GetChangeRequest$json = {
  '1': 'GetChangeRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetChangeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChangeRequestDescriptor = $convert
    .base64Decode('ChBHZXRDaGFuZ2VSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getChangeResponseDescriptor instead')
const GetChangeResponse$json = {
  '1': 'GetChangeResponse',
  '2': [
    {
      '1': 'changeinfo',
      '3': 436394734,
      '4': 1,
      '5': 11,
      '6': '.route53.ChangeInfo',
      '10': 'changeinfo'
    },
  ],
};

/// Descriptor for `GetChangeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChangeResponseDescriptor = $convert.base64Decode(
    'ChFHZXRDaGFuZ2VSZXNwb25zZRI3CgpjaGFuZ2VpbmZvGO61i9ABIAEoCzITLnJvdXRlNTMuQ2'
    'hhbmdlSW5mb1IKY2hhbmdlaW5mbw==');

@$core.Deprecated('Use getCheckerIpRangesRequestDescriptor instead')
const GetCheckerIpRangesRequest$json = {
  '1': 'GetCheckerIpRangesRequest',
};

/// Descriptor for `GetCheckerIpRangesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCheckerIpRangesRequestDescriptor =
    $convert.base64Decode('ChlHZXRDaGVja2VySXBSYW5nZXNSZXF1ZXN0');

@$core.Deprecated('Use getCheckerIpRangesResponseDescriptor instead')
const GetCheckerIpRangesResponse$json = {
  '1': 'GetCheckerIpRangesResponse',
  '2': [
    {
      '1': 'checkeripranges',
      '3': 376098860,
      '4': 3,
      '5': 9,
      '10': 'checkeripranges'
    },
  ],
};

/// Descriptor for `GetCheckerIpRangesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCheckerIpRangesResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRDaGVja2VySXBSYW5nZXNSZXNwb25zZRIsCg9jaGVja2VyaXByYW5nZXMYrKCrswEgAy'
        'gJUg9jaGVja2VyaXByYW5nZXM=');

@$core.Deprecated('Use getDNSSECRequestDescriptor instead')
const GetDNSSECRequest$json = {
  '1': 'GetDNSSECRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `GetDNSSECRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDNSSECRequestDescriptor = $convert.base64Decode(
    'ChBHZXRETlNTRUNSZXF1ZXN0EiYKDGhvc3RlZHpvbmVpZBj+zp6lASABKAlSDGhvc3RlZHpvbm'
    'VpZA==');

@$core.Deprecated('Use getDNSSECResponseDescriptor instead')
const GetDNSSECResponse$json = {
  '1': 'GetDNSSECResponse',
  '2': [
    {
      '1': 'keysigningkeys',
      '3': 87847996,
      '4': 3,
      '5': 11,
      '6': '.route53.KeySigningKey',
      '10': 'keysigningkeys'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.route53.DNSSECStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetDNSSECResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDNSSECResponseDescriptor = $convert.base64Decode(
    'ChFHZXRETlNTRUNSZXNwb25zZRJBCg5rZXlzaWduaW5na2V5cxi86PEpIAMoCzIWLnJvdXRlNT'
    'MuS2V5U2lnbmluZ0tleVIOa2V5c2lnbmluZ2tleXMSMAoGc3RhdHVzGJDk+wIgASgLMhUucm91'
    'dGU1My5ETlNTRUNTdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use getGeoLocationRequestDescriptor instead')
const GetGeoLocationRequest$json = {
  '1': 'GetGeoLocationRequest',
  '2': [
    {
      '1': 'continentcode',
      '3': 180194347,
      '4': 1,
      '5': 9,
      '10': 'continentcode'
    },
    {'1': 'countrycode', '3': 485287065, '4': 1, '5': 9, '10': 'countrycode'},
    {
      '1': 'subdivisioncode',
      '3': 444328268,
      '4': 1,
      '5': 9,
      '10': 'subdivisioncode'
    },
  ],
};

/// Descriptor for `GetGeoLocationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGeoLocationRequestDescriptor = $convert.base64Decode(
    'ChVHZXRHZW9Mb2NhdGlvblJlcXVlc3QSJwoNY29udGluZW50Y29kZRirmPZVIAEoCVINY29udG'
    'luZW50Y29kZRIkCgtjb3VudHJ5Y29kZRiZybPnASABKAlSC2NvdW50cnljb2RlEiwKD3N1YmRp'
    'dmlzaW9uY29kZRjM0u/TASABKAlSD3N1YmRpdmlzaW9uY29kZQ==');

@$core.Deprecated('Use getGeoLocationResponseDescriptor instead')
const GetGeoLocationResponse$json = {
  '1': 'GetGeoLocationResponse',
  '2': [
    {
      '1': 'geolocationdetails',
      '3': 70816602,
      '4': 1,
      '5': 11,
      '6': '.route53.GeoLocationDetails',
      '10': 'geolocationdetails'
    },
  ],
};

/// Descriptor for `GetGeoLocationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGeoLocationResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRHZW9Mb2NhdGlvblJlc3BvbnNlEk4KEmdlb2xvY2F0aW9uZGV0YWlscxjapuIhIAEoCz'
        'IbLnJvdXRlNTMuR2VvTG9jYXRpb25EZXRhaWxzUhJnZW9sb2NhdGlvbmRldGFpbHM=');

@$core.Deprecated('Use getHealthCheckCountRequestDescriptor instead')
const GetHealthCheckCountRequest$json = {
  '1': 'GetHealthCheckCountRequest',
};

/// Descriptor for `GetHealthCheckCountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckCountRequestDescriptor =
    $convert.base64Decode('ChpHZXRIZWFsdGhDaGVja0NvdW50UmVxdWVzdA==');

@$core.Deprecated('Use getHealthCheckCountResponseDescriptor instead')
const GetHealthCheckCountResponse$json = {
  '1': 'GetHealthCheckCountResponse',
  '2': [
    {
      '1': 'healthcheckcount',
      '3': 332677785,
      '4': 1,
      '5': 3,
      '10': 'healthcheckcount'
    },
  ],
};

/// Descriptor for `GetHealthCheckCountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckCountResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRIZWFsdGhDaGVja0NvdW50UmVzcG9uc2USLgoQaGVhbHRoY2hlY2tjb3VudBiZhdGeAS'
        'ABKANSEGhlYWx0aGNoZWNrY291bnQ=');

@$core
    .Deprecated('Use getHealthCheckLastFailureReasonRequestDescriptor instead')
const GetHealthCheckLastFailureReasonRequest$json = {
  '1': 'GetHealthCheckLastFailureReasonRequest',
  '2': [
    {
      '1': 'healthcheckid',
      '3': 312971637,
      '4': 1,
      '5': 9,
      '10': 'healthcheckid'
    },
  ],
};

/// Descriptor for `GetHealthCheckLastFailureReasonRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckLastFailureReasonRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXRIZWFsdGhDaGVja0xhc3RGYWlsdXJlUmVhc29uUmVxdWVzdBIoCg1oZWFsdGhjaGVja2'
        'lkGPWinpUBIAEoCVINaGVhbHRoY2hlY2tpZA==');

@$core
    .Deprecated('Use getHealthCheckLastFailureReasonResponseDescriptor instead')
const GetHealthCheckLastFailureReasonResponse$json = {
  '1': 'GetHealthCheckLastFailureReasonResponse',
  '2': [
    {
      '1': 'healthcheckobservations',
      '3': 391890103,
      '4': 3,
      '5': 11,
      '6': '.route53.HealthCheckObservation',
      '10': 'healthcheckobservations'
    },
  ],
};

/// Descriptor for `GetHealthCheckLastFailureReasonResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckLastFailureReasonResponseDescriptor =
    $convert.base64Decode(
        'CidHZXRIZWFsdGhDaGVja0xhc3RGYWlsdXJlUmVhc29uUmVzcG9uc2USXQoXaGVhbHRoY2hlY2'
        'tvYnNlcnZhdGlvbnMYt4nvugEgAygLMh8ucm91dGU1My5IZWFsdGhDaGVja09ic2VydmF0aW9u'
        'UhdoZWFsdGhjaGVja29ic2VydmF0aW9ucw==');

@$core.Deprecated('Use getHealthCheckRequestDescriptor instead')
const GetHealthCheckRequest$json = {
  '1': 'GetHealthCheckRequest',
  '2': [
    {
      '1': 'healthcheckid',
      '3': 312971637,
      '4': 1,
      '5': 9,
      '10': 'healthcheckid'
    },
  ],
};

/// Descriptor for `GetHealthCheckRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckRequestDescriptor = $convert.base64Decode(
    'ChVHZXRIZWFsdGhDaGVja1JlcXVlc3QSKAoNaGVhbHRoY2hlY2tpZBj1op6VASABKAlSDWhlYW'
    'x0aGNoZWNraWQ=');

@$core.Deprecated('Use getHealthCheckResponseDescriptor instead')
const GetHealthCheckResponse$json = {
  '1': 'GetHealthCheckResponse',
  '2': [
    {
      '1': 'healthcheck',
      '3': 377540610,
      '4': 1,
      '5': 11,
      '6': '.route53.HealthCheck',
      '10': 'healthcheck'
    },
  ],
};

/// Descriptor for `GetHealthCheckResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRIZWFsdGhDaGVja1Jlc3BvbnNlEjoKC2hlYWx0aGNoZWNrGIKgg7QBIAEoCzIULnJvdX'
        'RlNTMuSGVhbHRoQ2hlY2tSC2hlYWx0aGNoZWNr');

@$core.Deprecated('Use getHealthCheckStatusRequestDescriptor instead')
const GetHealthCheckStatusRequest$json = {
  '1': 'GetHealthCheckStatusRequest',
  '2': [
    {
      '1': 'healthcheckid',
      '3': 312971637,
      '4': 1,
      '5': 9,
      '10': 'healthcheckid'
    },
  ],
};

/// Descriptor for `GetHealthCheckStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckStatusRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRIZWFsdGhDaGVja1N0YXR1c1JlcXVlc3QSKAoNaGVhbHRoY2hlY2tpZBj1op6VASABKA'
        'lSDWhlYWx0aGNoZWNraWQ=');

@$core.Deprecated('Use getHealthCheckStatusResponseDescriptor instead')
const GetHealthCheckStatusResponse$json = {
  '1': 'GetHealthCheckStatusResponse',
  '2': [
    {
      '1': 'healthcheckobservations',
      '3': 391890103,
      '4': 3,
      '5': 11,
      '6': '.route53.HealthCheckObservation',
      '10': 'healthcheckobservations'
    },
  ],
};

/// Descriptor for `GetHealthCheckStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthCheckStatusResponseDescriptor =
    $convert.base64Decode(
        'ChxHZXRIZWFsdGhDaGVja1N0YXR1c1Jlc3BvbnNlEl0KF2hlYWx0aGNoZWNrb2JzZXJ2YXRpb2'
        '5zGLeJ77oBIAMoCzIfLnJvdXRlNTMuSGVhbHRoQ2hlY2tPYnNlcnZhdGlvblIXaGVhbHRoY2hl'
        'Y2tvYnNlcnZhdGlvbnM=');

@$core.Deprecated('Use getHostedZoneCountRequestDescriptor instead')
const GetHostedZoneCountRequest$json = {
  '1': 'GetHostedZoneCountRequest',
};

/// Descriptor for `GetHostedZoneCountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHostedZoneCountRequestDescriptor =
    $convert.base64Decode('ChlHZXRIb3N0ZWRab25lQ291bnRSZXF1ZXN0');

@$core.Deprecated('Use getHostedZoneCountResponseDescriptor instead')
const GetHostedZoneCountResponse$json = {
  '1': 'GetHostedZoneCountResponse',
  '2': [
    {
      '1': 'hostedzonecount',
      '3': 93723536,
      '4': 1,
      '5': 3,
      '10': 'hostedzonecount'
    },
  ],
};

/// Descriptor for `GetHostedZoneCountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHostedZoneCountResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRIb3N0ZWRab25lQ291bnRSZXNwb25zZRIrCg9ob3N0ZWR6b25lY291bnQYkLfYLCABKA'
        'NSD2hvc3RlZHpvbmVjb3VudA==');

@$core.Deprecated('Use getHostedZoneLimitRequestDescriptor instead')
const GetHostedZoneLimitRequest$json = {
  '1': 'GetHostedZoneLimitRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.HostedZoneLimitType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `GetHostedZoneLimitRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHostedZoneLimitRequestDescriptor = $convert.base64Decode(
    'ChlHZXRIb3N0ZWRab25lTGltaXRSZXF1ZXN0EiYKDGhvc3RlZHpvbmVpZBj+zp6lASABKAlSDG'
    'hvc3RlZHpvbmVpZBI0CgR0eXBlGO6g14oBIAEoDjIcLnJvdXRlNTMuSG9zdGVkWm9uZUxpbWl0'
    'VHlwZVIEdHlwZQ==');

@$core.Deprecated('Use getHostedZoneLimitResponseDescriptor instead')
const GetHostedZoneLimitResponse$json = {
  '1': 'GetHostedZoneLimitResponse',
  '2': [
    {'1': 'count', '3': 31963285, '4': 1, '5': 3, '10': 'count'},
    {
      '1': 'limit',
      '3': 412502741,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZoneLimit',
      '10': 'limit'
    },
  ],
};

/// Descriptor for `GetHostedZoneLimitResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHostedZoneLimitResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRIb3N0ZWRab25lTGltaXRSZXNwb25zZRIXCgVjb3VudBiV8Z4PIAEoA1IFY291bnQSMg'
        'oFbGltaXQY1ZXZxAEgASgLMhgucm91dGU1My5Ib3N0ZWRab25lTGltaXRSBWxpbWl0');

@$core.Deprecated('Use getHostedZoneRequestDescriptor instead')
const GetHostedZoneRequest$json = {
  '1': 'GetHostedZoneRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetHostedZoneRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHostedZoneRequestDescriptor = $convert
    .base64Decode('ChRHZXRIb3N0ZWRab25lUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use getHostedZoneResponseDescriptor instead')
const GetHostedZoneResponse$json = {
  '1': 'GetHostedZoneResponse',
  '2': [
    {
      '1': 'delegationset',
      '3': 190771750,
      '4': 1,
      '5': 11,
      '6': '.route53.DelegationSet',
      '10': 'delegationset'
    },
    {
      '1': 'hostedzone',
      '3': 465689249,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZone',
      '10': 'hostedzone'
    },
    {
      '1': 'vpcs',
      '3': 424064898,
      '4': 3,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpcs'
    },
  ],
};

/// Descriptor for `GetHostedZoneResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHostedZoneResponseDescriptor = $convert.base64Decode(
    'ChVHZXRIb3N0ZWRab25lUmVzcG9uc2USPwoNZGVsZWdhdGlvbnNldBim5PtaIAEoCzIWLnJvdX'
    'RlNTMuRGVsZWdhdGlvblNldFINZGVsZWdhdGlvbnNldBI3Cgpob3N0ZWR6b25lGKG1h94BIAEo'
    'CzITLnJvdXRlNTMuSG9zdGVkWm9uZVIKaG9zdGVkem9uZRIkCgR2cGNzGILvmsoBIAMoCzIMLn'
    'JvdXRlNTMuVlBDUgR2cGNz');

@$core.Deprecated('Use getQueryLoggingConfigRequestDescriptor instead')
const GetQueryLoggingConfigRequest$json = {
  '1': 'GetQueryLoggingConfigRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetQueryLoggingConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryLoggingConfigRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRRdWVyeUxvZ2dpbmdDb25maWdSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getQueryLoggingConfigResponseDescriptor instead')
const GetQueryLoggingConfigResponse$json = {
  '1': 'GetQueryLoggingConfigResponse',
  '2': [
    {
      '1': 'queryloggingconfig',
      '3': 492751675,
      '4': 1,
      '5': 11,
      '6': '.route53.QueryLoggingConfig',
      '10': 'queryloggingconfig'
    },
  ],
};

/// Descriptor for `GetQueryLoggingConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryLoggingConfigResponseDescriptor =
    $convert.base64Decode(
        'Ch1HZXRRdWVyeUxvZ2dpbmdDb25maWdSZXNwb25zZRJPChJxdWVyeWxvZ2dpbmdjb25maWcYu5'
        'b76gEgASgLMhsucm91dGU1My5RdWVyeUxvZ2dpbmdDb25maWdSEnF1ZXJ5bG9nZ2luZ2NvbmZp'
        'Zw==');

@$core.Deprecated('Use getReusableDelegationSetLimitRequestDescriptor instead')
const GetReusableDelegationSetLimitRequest$json = {
  '1': 'GetReusableDelegationSetLimitRequest',
  '2': [
    {
      '1': 'delegationsetid',
      '3': 307328801,
      '4': 1,
      '5': 9,
      '10': 'delegationsetid'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.ReusableDelegationSetLimitType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `GetReusableDelegationSetLimitRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReusableDelegationSetLimitRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRSZXVzYWJsZURlbGVnYXRpb25TZXRMaW1pdFJlcXVlc3QSLAoPZGVsZWdhdGlvbnNldG'
        'lkGKHuxZIBIAEoCVIPZGVsZWdhdGlvbnNldGlkEj8KBHR5cGUY7qDXigEgASgOMicucm91dGU1'
        'My5SZXVzYWJsZURlbGVnYXRpb25TZXRMaW1pdFR5cGVSBHR5cGU=');

@$core.Deprecated('Use getReusableDelegationSetLimitResponseDescriptor instead')
const GetReusableDelegationSetLimitResponse$json = {
  '1': 'GetReusableDelegationSetLimitResponse',
  '2': [
    {'1': 'count', '3': 31963285, '4': 1, '5': 3, '10': 'count'},
    {
      '1': 'limit',
      '3': 412502741,
      '4': 1,
      '5': 11,
      '6': '.route53.ReusableDelegationSetLimit',
      '10': 'limit'
    },
  ],
};

/// Descriptor for `GetReusableDelegationSetLimitResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReusableDelegationSetLimitResponseDescriptor =
    $convert.base64Decode(
        'CiVHZXRSZXVzYWJsZURlbGVnYXRpb25TZXRMaW1pdFJlc3BvbnNlEhcKBWNvdW50GJXxng8gAS'
        'gDUgVjb3VudBI9CgVsaW1pdBjVldnEASABKAsyIy5yb3V0ZTUzLlJldXNhYmxlRGVsZWdhdGlv'
        'blNldExpbWl0UgVsaW1pdA==');

@$core.Deprecated('Use getReusableDelegationSetRequestDescriptor instead')
const GetReusableDelegationSetRequest$json = {
  '1': 'GetReusableDelegationSetRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetReusableDelegationSetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReusableDelegationSetRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRSZXVzYWJsZURlbGVnYXRpb25TZXRSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getReusableDelegationSetResponseDescriptor instead')
const GetReusableDelegationSetResponse$json = {
  '1': 'GetReusableDelegationSetResponse',
  '2': [
    {
      '1': 'delegationset',
      '3': 190771750,
      '4': 1,
      '5': 11,
      '6': '.route53.DelegationSet',
      '10': 'delegationset'
    },
  ],
};

/// Descriptor for `GetReusableDelegationSetResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReusableDelegationSetResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXRSZXVzYWJsZURlbGVnYXRpb25TZXRSZXNwb25zZRI/Cg1kZWxlZ2F0aW9uc2V0GKbk+1'
        'ogASgLMhYucm91dGU1My5EZWxlZ2F0aW9uU2V0Ug1kZWxlZ2F0aW9uc2V0');

@$core.Deprecated('Use getTrafficPolicyInstanceCountRequestDescriptor instead')
const GetTrafficPolicyInstanceCountRequest$json = {
  '1': 'GetTrafficPolicyInstanceCountRequest',
};

/// Descriptor for `GetTrafficPolicyInstanceCountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrafficPolicyInstanceCountRequestDescriptor =
    $convert
        .base64Decode('CiRHZXRUcmFmZmljUG9saWN5SW5zdGFuY2VDb3VudFJlcXVlc3Q=');

@$core.Deprecated('Use getTrafficPolicyInstanceCountResponseDescriptor instead')
const GetTrafficPolicyInstanceCountResponse$json = {
  '1': 'GetTrafficPolicyInstanceCountResponse',
  '2': [
    {
      '1': 'trafficpolicyinstancecount',
      '3': 190494003,
      '4': 1,
      '5': 5,
      '10': 'trafficpolicyinstancecount'
    },
  ],
};

/// Descriptor for `GetTrafficPolicyInstanceCountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrafficPolicyInstanceCountResponseDescriptor =
    $convert.base64Decode(
        'CiVHZXRUcmFmZmljUG9saWN5SW5zdGFuY2VDb3VudFJlc3BvbnNlEkEKGnRyYWZmaWNwb2xpY3'
        'lpbnN0YW5jZWNvdW50GLPq6logASgFUhp0cmFmZmljcG9saWN5aW5zdGFuY2Vjb3VudA==');

@$core.Deprecated('Use getTrafficPolicyInstanceRequestDescriptor instead')
const GetTrafficPolicyInstanceRequest$json = {
  '1': 'GetTrafficPolicyInstanceRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetTrafficPolicyInstanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrafficPolicyInstanceRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRUcmFmZmljUG9saWN5SW5zdGFuY2VSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getTrafficPolicyInstanceResponseDescriptor instead')
const GetTrafficPolicyInstanceResponse$json = {
  '1': 'GetTrafficPolicyInstanceResponse',
  '2': [
    {
      '1': 'trafficpolicyinstance',
      '3': 205651476,
      '4': 1,
      '5': 11,
      '6': '.route53.TrafficPolicyInstance',
      '10': 'trafficpolicyinstance'
    },
  ],
};

/// Descriptor for `GetTrafficPolicyInstanceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrafficPolicyInstanceResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXRUcmFmZmljUG9saWN5SW5zdGFuY2VSZXNwb25zZRJXChV0cmFmZmljcG9saWN5aW5zdG'
        'FuY2UYlPyHYiABKAsyHi5yb3V0ZTUzLlRyYWZmaWNQb2xpY3lJbnN0YW5jZVIVdHJhZmZpY3Bv'
        'bGljeWluc3RhbmNl');

@$core.Deprecated('Use getTrafficPolicyRequestDescriptor instead')
const GetTrafficPolicyRequest$json = {
  '1': 'GetTrafficPolicyRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 5, '10': 'version'},
  ],
};

/// Descriptor for `GetTrafficPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrafficPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRUcmFmZmljUG9saWN5UmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEhwKB3ZlcnNpb2'
        '4YuKq37gEgASgFUgd2ZXJzaW9u');

@$core.Deprecated('Use getTrafficPolicyResponseDescriptor instead')
const GetTrafficPolicyResponse$json = {
  '1': 'GetTrafficPolicyResponse',
  '2': [
    {
      '1': 'trafficpolicy',
      '3': 154595657,
      '4': 1,
      '5': 11,
      '6': '.route53.TrafficPolicy',
      '10': 'trafficpolicy'
    },
  ],
};

/// Descriptor for `GetTrafficPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrafficPolicyResponseDescriptor =
    $convert.base64Decode(
        'ChhHZXRUcmFmZmljUG9saWN5UmVzcG9uc2USPwoNdHJhZmZpY3BvbGljeRjJ4ttJIAEoCzIWLn'
        'JvdXRlNTMuVHJhZmZpY1BvbGljeVINdHJhZmZpY3BvbGljeQ==');

@$core.Deprecated('Use healthCheckDescriptor instead')
const HealthCheck$json = {
  '1': 'HealthCheck',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {
      '1': 'cloudwatchalarmconfiguration',
      '3': 358915841,
      '4': 1,
      '5': 11,
      '6': '.route53.CloudWatchAlarmConfiguration',
      '10': 'cloudwatchalarmconfiguration'
    },
    {
      '1': 'healthcheckconfig',
      '3': 83111204,
      '4': 1,
      '5': 11,
      '6': '.route53.HealthCheckConfig',
      '10': 'healthcheckconfig'
    },
    {
      '1': 'healthcheckversion',
      '3': 89568396,
      '4': 1,
      '5': 3,
      '10': 'healthcheckversion'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'linkedservice',
      '3': 438614164,
      '4': 1,
      '5': 11,
      '6': '.route53.LinkedService',
      '10': 'linkedservice'
    },
  ],
};

/// Descriptor for `HealthCheck`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthCheckDescriptor = $convert.base64Decode(
    'CgtIZWFsdGhDaGVjaxIrCg9jYWxsZXJyZWZlcmVuY2UYmJmNSCABKAlSD2NhbGxlcnJlZmVyZW'
    '5jZRJtChxjbG91ZHdhdGNoYWxhcm1jb25maWd1cmF0aW9uGIG+kqsBIAEoCzIlLnJvdXRlNTMu'
    'Q2xvdWRXYXRjaEFsYXJtQ29uZmlndXJhdGlvblIcY2xvdWR3YXRjaGFsYXJtY29uZmlndXJhdG'
    'lvbhJLChFoZWFsdGhjaGVja2NvbmZpZxik2tAnIAEoCzIaLnJvdXRlNTMuSGVhbHRoQ2hlY2tD'
    'b25maWdSEWhlYWx0aGNoZWNrY29uZmlnEjEKEmhlYWx0aGNoZWNrdmVyc2lvbhiM6doqIAEoA1'
    'ISaGVhbHRoY2hlY2t2ZXJzaW9uEhIKAmlkGIHyorcBIAEoCVICaWQSQAoNbGlua2Vkc2Vydmlj'
    'ZRiU8ZLRASABKAsyFi5yb3V0ZTUzLkxpbmtlZFNlcnZpY2VSDWxpbmtlZHNlcnZpY2U=');

@$core.Deprecated('Use healthCheckAlreadyExistsDescriptor instead')
const HealthCheckAlreadyExists$json = {
  '1': 'HealthCheckAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HealthCheckAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthCheckAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChhIZWFsdGhDaGVja0FscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use healthCheckConfigDescriptor instead')
const HealthCheckConfig$json = {
  '1': 'HealthCheckConfig',
  '2': [
    {
      '1': 'alarmidentifier',
      '3': 536124346,
      '4': 1,
      '5': 11,
      '6': '.route53.AlarmIdentifier',
      '10': 'alarmidentifier'
    },
    {
      '1': 'childhealthchecks',
      '3': 485535935,
      '4': 3,
      '5': 9,
      '10': 'childhealthchecks'
    },
    {'1': 'disabled', '3': 533633318, '4': 1, '5': 8, '10': 'disabled'},
    {'1': 'enablesni', '3': 70122887, '4': 1, '5': 8, '10': 'enablesni'},
    {
      '1': 'failurethreshold',
      '3': 176846565,
      '4': 1,
      '5': 5,
      '10': 'failurethreshold'
    },
    {
      '1': 'fullyqualifieddomainname',
      '3': 459321509,
      '4': 1,
      '5': 9,
      '10': 'fullyqualifieddomainname'
    },
    {
      '1': 'healththreshold',
      '3': 215873163,
      '4': 1,
      '5': 5,
      '10': 'healththreshold'
    },
    {'1': 'ipaddress', '3': 169333741, '4': 1, '5': 9, '10': 'ipaddress'},
    {
      '1': 'insufficientdatahealthstatus',
      '3': 493115723,
      '4': 1,
      '5': 14,
      '6': '.route53.InsufficientDataHealthStatus',
      '10': 'insufficientdatahealthstatus'
    },
    {'1': 'inverted', '3': 55175513, '4': 1, '5': 8, '10': 'inverted'},
    {
      '1': 'measurelatency',
      '3': 87136848,
      '4': 1,
      '5': 8,
      '10': 'measurelatency'
    },
    {'1': 'port', '3': 46480583, '4': 1, '5': 5, '10': 'port'},
    {
      '1': 'regions',
      '3': 36200107,
      '4': 3,
      '5': 14,
      '6': '.route53.HealthCheckRegion',
      '10': 'regions'
    },
    {
      '1': 'requestinterval',
      '3': 350673112,
      '4': 1,
      '5': 5,
      '10': 'requestinterval'
    },
    {'1': 'resourcepath', '3': 117584551, '4': 1, '5': 9, '10': 'resourcepath'},
    {
      '1': 'routingcontrolarn',
      '3': 206883790,
      '4': 1,
      '5': 9,
      '10': 'routingcontrolarn'
    },
    {'1': 'searchstring', '3': 318687365, '4': 1, '5': 9, '10': 'searchstring'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.HealthCheckType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `HealthCheckConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthCheckConfigDescriptor = $convert.base64Decode(
    'ChFIZWFsdGhDaGVja0NvbmZpZxJGCg9hbGFybWlkZW50aWZpZXIYurfS/wEgASgLMhgucm91dG'
    'U1My5BbGFybUlkZW50aWZpZXJSD2FsYXJtaWRlbnRpZmllchIwChFjaGlsZGhlYWx0aGNoZWNr'
    'cxi/4cLnASADKAlSEWNoaWxkaGVhbHRoY2hlY2tzEh4KCGRpc2FibGVkGKayuv4BIAEoCFIIZG'
    'lzYWJsZWQSHwoJZW5hYmxlc25pGIf7tyEgASgIUgllbmFibGVzbmkSLQoQZmFpbHVyZXRocmVz'
    'aG9sZBjl7alUIAEoBVIQZmFpbHVyZXRocmVzaG9sZBI+ChhmdWxseXF1YWxpZmllZGRvbWFpbm'
    '5hbWUYpeGC2wEgASgJUhhmdWxseXF1YWxpZmllZGRvbWFpbm5hbWUSKwoPaGVhbHRodGhyZXNo'
    'b2xkGIvt92YgASgFUg9oZWFsdGh0aHJlc2hvbGQSHwoJaXBhZGRyZXNzGO2n31AgASgJUglpcG'
    'FkZHJlc3MSbQocaW5zdWZmaWNpZW50ZGF0YWhlYWx0aHN0YXR1cxjLspHrASABKA4yJS5yb3V0'
    'ZTUzLkluc3VmZmljaWVudERhdGFIZWFsdGhTdGF0dXNSHGluc3VmZmljaWVudGRhdGFoZWFsdG'
    'hzdGF0dXMSHQoIaW52ZXJ0ZWQY2dKnGiABKAhSCGludmVydGVkEikKDm1lYXN1cmVsYXRlbmN5'
    'GNC0xikgASgIUg5tZWFzdXJlbGF0ZW5jeRIVCgRwb3J0GMf5lBYgASgFUgRwb3J0EjcKB3JlZ2'
    'lvbnMYq72hESADKA4yGi5yb3V0ZTUzLkhlYWx0aENoZWNrUmVnaW9uUgdyZWdpb25zEiwKD3Jl'
    'cXVlc3RpbnRlcnZhbBjYsZunASABKAVSD3JlcXVlc3RpbnRlcnZhbBIlCgxyZXNvdXJjZXBhdG'
    'gYp+WIOCABKAlSDHJlc291cmNlcGF0aBIvChFyb3V0aW5nY29udHJvbGFybhjOl9NiIAEoCVIR'
    'cm91dGluZ2NvbnRyb2xhcm4SJgoMc2VhcmNoc3RyaW5nGIWR+5cBIAEoCVIMc2VhcmNoc3RyaW'
    '5nEjAKBHR5cGUY7qDXigEgASgOMhgucm91dGU1My5IZWFsdGhDaGVja1R5cGVSBHR5cGU=');

@$core.Deprecated('Use healthCheckInUseDescriptor instead')
const HealthCheckInUse$json = {
  '1': 'HealthCheckInUse',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HealthCheckInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthCheckInUseDescriptor = $convert.base64Decode(
    'ChBIZWFsdGhDaGVja0luVXNlEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use healthCheckObservationDescriptor instead')
const HealthCheckObservation$json = {
  '1': 'HealthCheckObservation',
  '2': [
    {'1': 'ipaddress', '3': 169333741, '4': 1, '5': 9, '10': 'ipaddress'},
    {
      '1': 'region',
      '3': 154040478,
      '4': 1,
      '5': 14,
      '6': '.route53.HealthCheckRegion',
      '10': 'region'
    },
    {
      '1': 'statusreport',
      '3': 27958834,
      '4': 1,
      '5': 11,
      '6': '.route53.StatusReport',
      '10': 'statusreport'
    },
  ],
};

/// Descriptor for `HealthCheckObservation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthCheckObservationDescriptor = $convert.base64Decode(
    'ChZIZWFsdGhDaGVja09ic2VydmF0aW9uEh8KCWlwYWRkcmVzcxjtp99QIAEoCVIJaXBhZGRyZX'
    'NzEjUKBnJlZ2lvbhie8blJIAEoDjIaLnJvdXRlNTMuSGVhbHRoQ2hlY2tSZWdpb25SBnJlZ2lv'
    'bhI8CgxzdGF0dXNyZXBvcnQYsryqDSABKAsyFS5yb3V0ZTUzLlN0YXR1c1JlcG9ydFIMc3RhdH'
    'VzcmVwb3J0');

@$core.Deprecated('Use healthCheckVersionMismatchDescriptor instead')
const HealthCheckVersionMismatch$json = {
  '1': 'HealthCheckVersionMismatch',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HealthCheckVersionMismatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthCheckVersionMismatchDescriptor =
    $convert.base64Decode(
        'ChpIZWFsdGhDaGVja1ZlcnNpb25NaXNtYXRjaBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use hostedZoneDescriptor instead')
const HostedZone$json = {
  '1': 'HostedZone',
  '2': [
    {
      '1': 'callerreference',
      '3': 151211160,
      '4': 1,
      '5': 9,
      '10': 'callerreference'
    },
    {
      '1': 'config',
      '3': 169009384,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZoneConfig',
      '10': 'config'
    },
    {
      '1': 'features',
      '3': 205993963,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZoneFeatures',
      '10': 'features'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'linkedservice',
      '3': 438614164,
      '4': 1,
      '5': 11,
      '6': '.route53.LinkedService',
      '10': 'linkedservice'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'resourcerecordsetcount',
      '3': 334391116,
      '4': 1,
      '5': 3,
      '10': 'resourcerecordsetcount'
    },
  ],
};

/// Descriptor for `HostedZone`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneDescriptor = $convert.base64Decode(
    'CgpIb3N0ZWRab25lEisKD2NhbGxlcnJlZmVyZW5jZRiYmY1IIAEoCVIPY2FsbGVycmVmZXJlbm'
    'NlEjQKBmNvbmZpZxjowctQIAEoCzIZLnJvdXRlNTMuSG9zdGVkWm9uZUNvbmZpZ1IGY29uZmln'
    'EjoKCGZlYXR1cmVzGOvvnGIgASgLMhsucm91dGU1My5Ib3N0ZWRab25lRmVhdHVyZXNSCGZlYX'
    'R1cmVzEhIKAmlkGIHyorcBIAEoCVICaWQSQAoNbGlua2Vkc2VydmljZRiU8ZLRASABKAsyFi5y'
    'b3V0ZTUzLkxpbmtlZFNlcnZpY2VSDWxpbmtlZHNlcnZpY2USFQoEbmFtZRiH5oF/IAEoCVIEbm'
    'FtZRI6ChZyZXNvdXJjZXJlY29yZHNldGNvdW50GMzOuZ8BIAEoA1IWcmVzb3VyY2VyZWNvcmRz'
    'ZXRjb3VudA==');

@$core.Deprecated('Use hostedZoneAlreadyExistsDescriptor instead')
const HostedZoneAlreadyExists$json = {
  '1': 'HostedZoneAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HostedZoneAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChdIb3N0ZWRab25lQWxyZWFkeUV4aXN0cxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use hostedZoneConfigDescriptor instead')
const HostedZoneConfig$json = {
  '1': 'HostedZoneConfig',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'privatezone', '3': 506703867, '4': 1, '5': 8, '10': 'privatezone'},
  ],
};

/// Descriptor for `HostedZoneConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneConfigDescriptor = $convert.base64Decode(
    'ChBIb3N0ZWRab25lQ29uZmlnEhwKB2NvbW1lbnQY/7++wgEgASgJUgdjb21tZW50EiQKC3ByaX'
    'ZhdGV6b25lGPvfzvEBIAEoCFILcHJpdmF0ZXpvbmU=');

@$core.Deprecated('Use hostedZoneFailureReasonsDescriptor instead')
const HostedZoneFailureReasons$json = {
  '1': 'HostedZoneFailureReasons',
  '2': [
    {
      '1': 'acceleratedrecovery',
      '3': 356087092,
      '4': 1,
      '5': 9,
      '10': 'acceleratedrecovery'
    },
  ],
};

/// Descriptor for `HostedZoneFailureReasons`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneFailureReasonsDescriptor =
    $convert.base64Decode(
        'ChhIb3N0ZWRab25lRmFpbHVyZVJlYXNvbnMSNAoTYWNjZWxlcmF0ZWRyZWNvdmVyeRi06uWpAS'
        'ABKAlSE2FjY2VsZXJhdGVkcmVjb3Zlcnk=');

@$core.Deprecated('Use hostedZoneFeaturesDescriptor instead')
const HostedZoneFeatures$json = {
  '1': 'HostedZoneFeatures',
  '2': [
    {
      '1': 'acceleratedrecoverystatus',
      '3': 507534246,
      '4': 1,
      '5': 14,
      '6': '.route53.AcceleratedRecoveryStatus',
      '10': 'acceleratedrecoverystatus'
    },
    {
      '1': 'failurereasons',
      '3': 445146219,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZoneFailureReasons',
      '10': 'failurereasons'
    },
  ],
};

/// Descriptor for `HostedZoneFeatures`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneFeaturesDescriptor = $convert.base64Decode(
    'ChJIb3N0ZWRab25lRmVhdHVyZXMSZAoZYWNjZWxlcmF0ZWRyZWNvdmVyeXN0YXR1cximt4HyAS'
    'ABKA4yIi5yb3V0ZTUzLkFjY2VsZXJhdGVkUmVjb3ZlcnlTdGF0dXNSGWFjY2VsZXJhdGVkcmVj'
    'b3ZlcnlzdGF0dXMSTQoOZmFpbHVyZXJlYXNvbnMY68ih1AEgASgLMiEucm91dGU1My5Ib3N0ZW'
    'Rab25lRmFpbHVyZVJlYXNvbnNSDmZhaWx1cmVyZWFzb25z');

@$core.Deprecated('Use hostedZoneLimitDescriptor instead')
const HostedZoneLimit$json = {
  '1': 'HostedZoneLimit',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.HostedZoneLimitType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 3, '10': 'value'},
  ],
};

/// Descriptor for `HostedZoneLimit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneLimitDescriptor = $convert.base64Decode(
    'Cg9Ib3N0ZWRab25lTGltaXQSNAoEdHlwZRjuoNeKASABKA4yHC5yb3V0ZTUzLkhvc3RlZFpvbm'
    'VMaW1pdFR5cGVSBHR5cGUSGAoFdmFsdWUY6/KfigEgASgDUgV2YWx1ZQ==');

@$core.Deprecated('Use hostedZoneNotEmptyDescriptor instead')
const HostedZoneNotEmpty$json = {
  '1': 'HostedZoneNotEmpty',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HostedZoneNotEmpty`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneNotEmptyDescriptor =
    $convert.base64Decode(
        'ChJIb3N0ZWRab25lTm90RW1wdHkSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use hostedZoneNotFoundDescriptor instead')
const HostedZoneNotFound$json = {
  '1': 'HostedZoneNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HostedZoneNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneNotFoundDescriptor =
    $convert.base64Decode(
        'ChJIb3N0ZWRab25lTm90Rm91bmQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use hostedZoneNotPrivateDescriptor instead')
const HostedZoneNotPrivate$json = {
  '1': 'HostedZoneNotPrivate',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HostedZoneNotPrivate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneNotPrivateDescriptor =
    $convert.base64Decode(
        'ChRIb3N0ZWRab25lTm90UHJpdmF0ZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use hostedZoneOwnerDescriptor instead')
const HostedZoneOwner$json = {
  '1': 'HostedZoneOwner',
  '2': [
    {
      '1': 'owningaccount',
      '3': 339968073,
      '4': 1,
      '5': 9,
      '10': 'owningaccount'
    },
    {
      '1': 'owningservice',
      '3': 462487817,
      '4': 1,
      '5': 9,
      '10': 'owningservice'
    },
  ],
};

/// Descriptor for `HostedZoneOwner`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneOwnerDescriptor = $convert.base64Decode(
    'Cg9Ib3N0ZWRab25lT3duZXISKAoNb3duaW5nYWNjb3VudBjJgI6iASABKAlSDW93bmluZ2FjY2'
    '91bnQSKAoNb3duaW5nc2VydmljZRiJgsTcASABKAlSDW93bmluZ3NlcnZpY2U=');

@$core.Deprecated('Use hostedZonePartiallyDelegatedDescriptor instead')
const HostedZonePartiallyDelegated$json = {
  '1': 'HostedZonePartiallyDelegated',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HostedZonePartiallyDelegated`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZonePartiallyDelegatedDescriptor =
    $convert.base64Decode(
        'ChxIb3N0ZWRab25lUGFydGlhbGx5RGVsZWdhdGVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use hostedZoneSummaryDescriptor instead')
const HostedZoneSummary$json = {
  '1': 'HostedZoneSummary',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZoneOwner',
      '10': 'owner'
    },
  ],
};

/// Descriptor for `HostedZoneSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hostedZoneSummaryDescriptor = $convert.base64Decode(
    'ChFIb3N0ZWRab25lU3VtbWFyeRImCgxob3N0ZWR6b25laWQY/s6epQEgASgJUgxob3N0ZWR6b2'
    '5laWQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIyCgVvd25lchj1/IrZASABKAsyGC5yb3V0ZTUz'
    'Lkhvc3RlZFpvbmVPd25lclIFb3duZXI=');

@$core.Deprecated('Use incompatibleVersionDescriptor instead')
const IncompatibleVersion$json = {
  '1': 'IncompatibleVersion',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IncompatibleVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List incompatibleVersionDescriptor =
    $convert.base64Decode(
        'ChNJbmNvbXBhdGlibGVWZXJzaW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use insufficientCloudWatchLogsResourcePolicyDescriptor instead')
const InsufficientCloudWatchLogsResourcePolicy$json = {
  '1': 'InsufficientCloudWatchLogsResourcePolicy',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InsufficientCloudWatchLogsResourcePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insufficientCloudWatchLogsResourcePolicyDescriptor =
    $convert.base64Decode(
        'CihJbnN1ZmZpY2llbnRDbG91ZFdhdGNoTG9nc1Jlc291cmNlUG9saWN5EhsKB21lc3NhZ2UY5Z'
        'HIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidArgumentDescriptor instead')
const InvalidArgument$json = {
  '1': 'InvalidArgument',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidArgument`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidArgumentDescriptor = $convert.base64Decode(
    'Cg9JbnZhbGlkQXJndW1lbnQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidChangeBatchDescriptor instead')
const InvalidChangeBatch$json = {
  '1': 'InvalidChangeBatch',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'messages', '3': 230838, '4': 3, '5': 9, '10': 'messages'},
  ],
};

/// Descriptor for `InvalidChangeBatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidChangeBatchDescriptor = $convert.base64Decode(
    'ChJJbnZhbGlkQ2hhbmdlQmF0Y2gSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZRIcCghtZX'
    'NzYWdlcxi2iw4gAygJUghtZXNzYWdlcw==');

@$core.Deprecated('Use invalidDomainNameDescriptor instead')
const InvalidDomainName$json = {
  '1': 'InvalidDomainName',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDomainName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDomainNameDescriptor = $convert.base64Decode(
    'ChFJbnZhbGlkRG9tYWluTmFtZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidInputDescriptor instead')
const InvalidInput$json = {
  '1': 'InvalidInput',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInputDescriptor = $convert.base64Decode(
    'CgxJbnZhbGlkSW5wdXQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidKMSArnDescriptor instead')
const InvalidKMSArn$json = {
  '1': 'InvalidKMSArn',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidKMSArn`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidKMSArnDescriptor = $convert.base64Decode(
    'Cg1JbnZhbGlkS01TQXJuEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidKeySigningKeyNameDescriptor instead')
const InvalidKeySigningKeyName$json = {
  '1': 'InvalidKeySigningKeyName',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidKeySigningKeyName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidKeySigningKeyNameDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkS2V5U2lnbmluZ0tleU5hbWUSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidKeySigningKeyStatusDescriptor instead')
const InvalidKeySigningKeyStatus$json = {
  '1': 'InvalidKeySigningKeyStatus',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidKeySigningKeyStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidKeySigningKeyStatusDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkS2V5U2lnbmluZ0tleVN0YXR1cxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidPaginationTokenDescriptor instead')
const InvalidPaginationToken$json = {
  '1': 'InvalidPaginationToken',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidPaginationToken`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidPaginationTokenDescriptor =
    $convert.base64Decode(
        'ChZJbnZhbGlkUGFnaW5hdGlvblRva2VuEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidSigningStatusDescriptor instead')
const InvalidSigningStatus$json = {
  '1': 'InvalidSigningStatus',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSigningStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSigningStatusDescriptor =
    $convert.base64Decode(
        'ChRJbnZhbGlkU2lnbmluZ1N0YXR1cxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidTrafficPolicyDocumentDescriptor instead')
const InvalidTrafficPolicyDocument$json = {
  '1': 'InvalidTrafficPolicyDocument',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTrafficPolicyDocument`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTrafficPolicyDocumentDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkVHJhZmZpY1BvbGljeURvY3VtZW50EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidVPCIdDescriptor instead')
const InvalidVPCId$json = {
  '1': 'InvalidVPCId',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidVPCId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidVPCIdDescriptor = $convert.base64Decode(
    'CgxJbnZhbGlkVlBDSWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use keySigningKeyDescriptor instead')
const KeySigningKey$json = {
  '1': 'KeySigningKey',
  '2': [
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'dnskeyrecord', '3': 33128395, '4': 1, '5': 9, '10': 'dnskeyrecord'},
    {'1': 'dsrecord', '3': 123290342, '4': 1, '5': 9, '10': 'dsrecord'},
    {
      '1': 'digestalgorithmmnemonic',
      '3': 209907099,
      '4': 1,
      '5': 9,
      '10': 'digestalgorithmmnemonic'
    },
    {
      '1': 'digestalgorithmtype',
      '3': 255175667,
      '4': 1,
      '5': 5,
      '10': 'digestalgorithmtype'
    },
    {'1': 'digestvalue', '3': 277868387, '4': 1, '5': 9, '10': 'digestvalue'},
    {'1': 'flag', '3': 504820984, '4': 1, '5': 5, '10': 'flag'},
    {'1': 'keytag', '3': 245478853, '4': 1, '5': 5, '10': 'keytag'},
    {'1': 'kmsarn', '3': 205584974, '4': 1, '5': 9, '10': 'kmsarn'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'publickey', '3': 167335776, '4': 1, '5': 9, '10': 'publickey'},
    {
      '1': 'signingalgorithmmnemonic',
      '3': 379080734,
      '4': 1,
      '5': 9,
      '10': 'signingalgorithmmnemonic'
    },
    {
      '1': 'signingalgorithmtype',
      '3': 407768874,
      '4': 1,
      '5': 5,
      '10': 'signingalgorithmtype'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'statusmessage',
      '3': 72590095,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
  ],
};

/// Descriptor for `KeySigningKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keySigningKeyDescriptor = $convert.base64Decode(
    'Cg1LZXlTaWduaW5nS2V5EiQKC2NyZWF0ZWRkYXRlGLCw58YBIAEoCVILY3JlYXRlZGRhdGUSJQ'
    'oMZG5za2V5cmVjb3JkGMv/5Q8gASgJUgxkbnNrZXlyZWNvcmQSHQoIZHNyZWNvcmQY5oXlOiAB'
    'KAlSCGRzcmVjb3JkEjsKF2RpZ2VzdGFsZ29yaXRobW1uZW1vbmljGJvbi2QgASgJUhdkaWdlc3'
    'RhbGdvcml0aG1tbmVtb25pYxIzChNkaWdlc3RhbGdvcml0aG10eXBlGPPX1nkgASgFUhNkaWdl'
    'c3RhbGdvcml0aG10eXBlEiQKC2RpZ2VzdHZhbHVlGOPev4QBIAEoCVILZGlnZXN0dmFsdWUSFg'
    'oEZmxhZxj46dvwASABKAVSBGZsYWcSGQoGa2V5dGFnGMXrhnUgASgFUgZrZXl0YWcSGQoGa21z'
    'YXJuGM70g2IgASgJUgZrbXNhcm4SLQoQbGFzdG1vZGlmaWVkZGF0ZRjTiMgLIAEoCVIQbGFzdG'
    '1vZGlmaWVkZGF0ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEh8KCXB1YmxpY2tleRjgruVPIAEo'
    'CVIJcHVibGlja2V5Ej4KGHNpZ25pbmdhbGdvcml0aG1tbmVtb25pYxieoOG0ASABKAlSGHNpZ2'
    '5pbmdhbGdvcml0aG1tbmVtb25pYxI2ChRzaWduaW5nYWxnb3JpdGhtdHlwZRiqnrjCASABKAVS'
    'FHNpZ25pbmdhbGdvcml0aG10eXBlEhkKBnN0YXR1cxiQ5PsCIAEoCVIGc3RhdHVzEicKDXN0YX'
    'R1c21lc3NhZ2UYj8bOIiABKAlSDXN0YXR1c21lc3NhZ2U=');

@$core.Deprecated('Use keySigningKeyAlreadyExistsDescriptor instead')
const KeySigningKeyAlreadyExists$json = {
  '1': 'KeySigningKeyAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KeySigningKeyAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keySigningKeyAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChpLZXlTaWduaW5nS2V5QWxyZWFkeUV4aXN0cxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use keySigningKeyInParentDSRecordDescriptor instead')
const KeySigningKeyInParentDSRecord$json = {
  '1': 'KeySigningKeyInParentDSRecord',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KeySigningKeyInParentDSRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keySigningKeyInParentDSRecordDescriptor =
    $convert.base64Decode(
        'Ch1LZXlTaWduaW5nS2V5SW5QYXJlbnREU1JlY29yZBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use keySigningKeyInUseDescriptor instead')
const KeySigningKeyInUse$json = {
  '1': 'KeySigningKeyInUse',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KeySigningKeyInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keySigningKeyInUseDescriptor =
    $convert.base64Decode(
        'ChJLZXlTaWduaW5nS2V5SW5Vc2USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use keySigningKeyWithActiveStatusNotFoundDescriptor instead')
const KeySigningKeyWithActiveStatusNotFound$json = {
  '1': 'KeySigningKeyWithActiveStatusNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KeySigningKeyWithActiveStatusNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keySigningKeyWithActiveStatusNotFoundDescriptor =
    $convert.base64Decode(
        'CiVLZXlTaWduaW5nS2V5V2l0aEFjdGl2ZVN0YXR1c05vdEZvdW5kEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use lastVPCAssociationDescriptor instead')
const LastVPCAssociation$json = {
  '1': 'LastVPCAssociation',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LastVPCAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lastVPCAssociationDescriptor =
    $convert.base64Decode(
        'ChJMYXN0VlBDQXNzb2NpYXRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use limitsExceededDescriptor instead')
const LimitsExceeded$json = {
  '1': 'LimitsExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LimitsExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitsExceededDescriptor = $convert.base64Decode(
    'Cg5MaW1pdHNFeGNlZWRlZBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use linkedServiceDescriptor instead')
const LinkedService$json = {
  '1': 'LinkedService',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'serviceprincipal',
      '3': 146694383,
      '4': 1,
      '5': 9,
      '10': 'serviceprincipal'
    },
  ],
};

/// Descriptor for `LinkedService`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List linkedServiceDescriptor = $convert.base64Decode(
    'Cg1MaW5rZWRTZXJ2aWNlEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhItCh'
    'BzZXJ2aWNlcHJpbmNpcGFsGO/B+UUgASgJUhBzZXJ2aWNlcHJpbmNpcGFs');

@$core.Deprecated('Use listCidrBlocksRequestDescriptor instead')
const ListCidrBlocksRequest$json = {
  '1': 'ListCidrBlocksRequest',
  '2': [
    {'1': 'collectionid', '3': 128052453, '4': 1, '5': 9, '10': 'collectionid'},
    {'1': 'locationname', '3': 158186566, '4': 1, '5': 9, '10': 'locationname'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 9, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCidrBlocksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCidrBlocksRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0Q2lkckJsb2Nrc1JlcXVlc3QSJQoMY29sbGVjdGlvbmlkGOXZhz0gASgJUgxjb2xsZW'
    'N0aW9uaWQSJQoMbG9jYXRpb25uYW1lGMb4tksgASgJUgxsb2NhdGlvbm5hbWUSIgoKbWF4cmVz'
    'dWx0cxiyqJuDASABKAlSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
    '9rZW4=');

@$core.Deprecated('Use listCidrBlocksResponseDescriptor instead')
const ListCidrBlocksResponse$json = {
  '1': 'ListCidrBlocksResponse',
  '2': [
    {
      '1': 'cidrblocks',
      '3': 134660738,
      '4': 3,
      '5': 11,
      '6': '.route53.CidrBlockSummary',
      '10': 'cidrblocks'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCidrBlocksResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCidrBlocksResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0Q2lkckJsb2Nrc1Jlc3BvbnNlEjwKCmNpZHJibG9ja3MYgoWbQCADKAsyGS5yb3V0ZT'
    'UzLkNpZHJCbG9ja1N1bW1hcnlSCmNpZHJibG9ja3MSHwoJbmV4dHRva2VuGP6EumcgASgJUglu'
    'ZXh0dG9rZW4=');

@$core.Deprecated('Use listCidrCollectionsRequestDescriptor instead')
const ListCidrCollectionsRequest$json = {
  '1': 'ListCidrCollectionsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 9, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCidrCollectionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCidrCollectionsRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0Q2lkckNvbGxlY3Rpb25zUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoCVIKbW'
        'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listCidrCollectionsResponseDescriptor instead')
const ListCidrCollectionsResponse$json = {
  '1': 'ListCidrCollectionsResponse',
  '2': [
    {
      '1': 'cidrcollections',
      '3': 502378541,
      '4': 3,
      '5': 11,
      '6': '.route53.CollectionSummary',
      '10': 'cidrcollections'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCidrCollectionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCidrCollectionsResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0Q2lkckNvbGxlY3Rpb25zUmVzcG9uc2USSAoPY2lkcmNvbGxlY3Rpb25zGK3gxu8BIA'
        'MoCzIaLnJvdXRlNTMuQ29sbGVjdGlvblN1bW1hcnlSD2NpZHJjb2xsZWN0aW9ucxIfCgluZXh0'
        'dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listCidrLocationsRequestDescriptor instead')
const ListCidrLocationsRequest$json = {
  '1': 'ListCidrLocationsRequest',
  '2': [
    {'1': 'collectionid', '3': 128052453, '4': 1, '5': 9, '10': 'collectionid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 9, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCidrLocationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCidrLocationsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0Q2lkckxvY2F0aW9uc1JlcXVlc3QSJQoMY29sbGVjdGlvbmlkGOXZhz0gASgJUgxjb2'
    'xsZWN0aW9uaWQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAlSCm1heHJlc3VsdHMSHwoJbmV4dHRv'
    'a2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listCidrLocationsResponseDescriptor instead')
const ListCidrLocationsResponse$json = {
  '1': 'ListCidrLocationsResponse',
  '2': [
    {
      '1': 'cidrlocations',
      '3': 480481962,
      '4': 3,
      '5': 11,
      '6': '.route53.LocationSummary',
      '10': 'cidrlocations'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCidrLocationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCidrLocationsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0Q2lkckxvY2F0aW9uc1Jlc3BvbnNlEkIKDWNpZHJsb2NhdGlvbnMYqqWO5QEgAygLMh'
    'gucm91dGU1My5Mb2NhdGlvblN1bW1hcnlSDWNpZHJsb2NhdGlvbnMSHwoJbmV4dHRva2VuGP6E'
    'umcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listGeoLocationsRequestDescriptor instead')
const ListGeoLocationsRequest$json = {
  '1': 'ListGeoLocationsRequest',
  '2': [
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'startcontinentcode',
      '3': 272774541,
      '4': 1,
      '5': 9,
      '10': 'startcontinentcode'
    },
    {
      '1': 'startcountrycode',
      '3': 182995791,
      '4': 1,
      '5': 9,
      '10': 'startcountrycode'
    },
    {
      '1': 'startsubdivisioncode',
      '3': 206846290,
      '4': 1,
      '5': 9,
      '10': 'startsubdivisioncode'
    },
  ],
};

/// Descriptor for `ListGeoLocationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGeoLocationsRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0R2VvTG9jYXRpb25zUmVxdWVzdBIeCghtYXhpdGVtcxiU1trxASABKAlSCG1heGl0ZW'
    '1zEjIKEnN0YXJ0Y29udGluZW50Y29kZRiN64iCASABKAlSEnN0YXJ0Y29udGluZW50Y29kZRIt'
    'ChBzdGFydGNvdW50cnljb2RlGM+WoVcgASgJUhBzdGFydGNvdW50cnljb2RlEjUKFHN0YXJ0c3'
    'ViZGl2aXNpb25jb2RlGNLy0GIgASgJUhRzdGFydHN1YmRpdmlzaW9uY29kZQ==');

@$core.Deprecated('Use listGeoLocationsResponseDescriptor instead')
const ListGeoLocationsResponse$json = {
  '1': 'ListGeoLocationsResponse',
  '2': [
    {
      '1': 'geolocationdetailslist',
      '3': 448447430,
      '4': 3,
      '5': 11,
      '6': '.route53.GeoLocationDetails',
      '10': 'geolocationdetailslist'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'nextcontinentcode',
      '3': 14234018,
      '4': 1,
      '5': 9,
      '10': 'nextcontinentcode'
    },
    {
      '1': 'nextcountrycode',
      '3': 147948496,
      '4': 1,
      '5': 9,
      '10': 'nextcountrycode'
    },
    {
      '1': 'nextsubdivisioncode',
      '3': 385449253,
      '4': 1,
      '5': 9,
      '10': 'nextsubdivisioncode'
    },
  ],
};

/// Descriptor for `ListGeoLocationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGeoLocationsResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0R2VvTG9jYXRpb25zUmVzcG9uc2USVwoWZ2VvbG9jYXRpb25kZXRhaWxzbGlzdBjGh+'
    'vVASADKAsyGy5yb3V0ZTUzLkdlb0xvY2F0aW9uRGV0YWlsc1IWZ2VvbG9jYXRpb25kZXRhaWxz'
    'bGlzdBIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSHgoIbWF4aXRlbXMYlN'
    'ba8QEgASgJUghtYXhpdGVtcxIvChFuZXh0Y29udGluZW50Y29kZRii4+QGIAEoCVIRbmV4dGNv'
    'bnRpbmVudGNvZGUSKwoPbmV4dGNvdW50cnljb2RlGNCHxkYgASgJUg9uZXh0Y291bnRyeWNvZG'
    'USNAoTbmV4dHN1YmRpdmlzaW9uY29kZRil+uW3ASABKAlSE25leHRzdWJkaXZpc2lvbmNvZGU=');

@$core.Deprecated('Use listHealthChecksRequestDescriptor instead')
const ListHealthChecksRequest$json = {
  '1': 'ListHealthChecksRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListHealthChecksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHealthChecksRequestDescriptor =
    $convert.base64Decode(
        'ChdMaXN0SGVhbHRoQ2hlY2tzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCg'
        'htYXhpdGVtcxiU1trxASABKAlSCG1heGl0ZW1z');

@$core.Deprecated('Use listHealthChecksResponseDescriptor instead')
const ListHealthChecksResponse$json = {
  '1': 'ListHealthChecksResponse',
  '2': [
    {
      '1': 'healthchecks',
      '3': 516432759,
      '4': 3,
      '5': 11,
      '6': '.route53.HealthCheck',
      '10': 'healthchecks'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListHealthChecksResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHealthChecksResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0SGVhbHRoQ2hlY2tzUmVzcG9uc2USPAoMaGVhbHRoY2hlY2tzGPfGoPYBIAMoCzIULn'
    'JvdXRlNTMuSGVhbHRoQ2hlY2tSDGhlYWx0aGNoZWNrcxIjCgtpc3RydW5jYXRlZBjan7hzIAEo'
    'CFILaXN0cnVuY2F0ZWQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlN'
    'ba8QEgASgJUghtYXhpdGVtcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use listHostedZonesByNameRequestDescriptor instead')
const ListHostedZonesByNameRequest$json = {
  '1': 'ListHostedZonesByNameRequest',
  '2': [
    {'1': 'dnsname', '3': 171901432, '4': 1, '5': 9, '10': 'dnsname'},
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListHostedZonesByNameRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHostedZonesByNameRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0SG9zdGVkWm9uZXNCeU5hbWVSZXF1ZXN0EhsKB2Ruc25hbWUY+IP8USABKAlSB2Ruc2'
        '5hbWUSJgoMaG9zdGVkem9uZWlkGP7OnqUBIAEoCVIMaG9zdGVkem9uZWlkEh4KCG1heGl0ZW1z'
        'GJTW2vEBIAEoCVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listHostedZonesByNameResponseDescriptor instead')
const ListHostedZonesByNameResponse$json = {
  '1': 'ListHostedZonesByNameResponse',
  '2': [
    {'1': 'dnsname', '3': 171901432, '4': 1, '5': 9, '10': 'dnsname'},
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {
      '1': 'hostedzones',
      '3': 86735402,
      '4': 3,
      '5': 11,
      '6': '.route53.HostedZone',
      '10': 'hostedzones'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {'1': 'nextdnsname', '3': 488331797, '4': 1, '5': 9, '10': 'nextdnsname'},
    {
      '1': 'nexthostedzoneid',
      '3': 162450165,
      '4': 1,
      '5': 9,
      '10': 'nexthostedzoneid'
    },
  ],
};

/// Descriptor for `ListHostedZonesByNameResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHostedZonesByNameResponseDescriptor = $convert.base64Decode(
    'Ch1MaXN0SG9zdGVkWm9uZXNCeU5hbWVSZXNwb25zZRIbCgdkbnNuYW1lGPiD/FEgASgJUgdkbn'
    'NuYW1lEiYKDGhvc3RlZHpvbmVpZBj+zp6lASABKAlSDGhvc3RlZHpvbmVpZBI4Cgtob3N0ZWR6'
    'b25lcxiq9K0pIAMoCzITLnJvdXRlNTMuSG9zdGVkWm9uZVILaG9zdGVkem9uZXMSIwoLaXN0cn'
    'VuY2F0ZWQY2p+4cyABKAhSC2lzdHJ1bmNhdGVkEh4KCG1heGl0ZW1zGJTW2vEBIAEoCVIIbWF4'
    'aXRlbXMSJAoLbmV4dGRuc25hbWUYlbTt6AEgASgJUgtuZXh0ZG5zbmFtZRItChBuZXh0aG9zdG'
    'Vkem9uZWlkGPWVu00gASgJUhBuZXh0aG9zdGVkem9uZWlk');

@$core.Deprecated('Use listHostedZonesByVPCRequestDescriptor instead')
const ListHostedZonesByVPCRequest$json = {
  '1': 'ListHostedZonesByVPCRequest',
  '2': [
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'vpcid', '3': 325135798, '4': 1, '5': 9, '10': 'vpcid'},
    {
      '1': 'vpcregion',
      '3': 474180765,
      '4': 1,
      '5': 14,
      '6': '.route53.VPCRegion',
      '10': 'vpcregion'
    },
  ],
};

/// Descriptor for `ListHostedZonesByVPCRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHostedZonesByVPCRequestDescriptor = $convert.base64Decode(
    'ChtMaXN0SG9zdGVkWm9uZXNCeVZQQ1JlcXVlc3QSHgoIbWF4aXRlbXMYlNba8QEgASgJUghtYX'
    'hpdGVtcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIYCgV2cGNpZBi224SbASAB'
    'KAlSBXZwY2lkEjQKCXZwY3JlZ2lvbhid2Y3iASABKA4yEi5yb3V0ZTUzLlZQQ1JlZ2lvblIJdn'
    'BjcmVnaW9u');

@$core.Deprecated('Use listHostedZonesByVPCResponseDescriptor instead')
const ListHostedZonesByVPCResponse$json = {
  '1': 'ListHostedZonesByVPCResponse',
  '2': [
    {
      '1': 'hostedzonesummaries',
      '3': 111631021,
      '4': 3,
      '5': 11,
      '6': '.route53.HostedZoneSummary',
      '10': 'hostedzonesummaries'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListHostedZonesByVPCResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHostedZonesByVPCResponseDescriptor = $convert.base64Decode(
    'ChxMaXN0SG9zdGVkWm9uZXNCeVZQQ1Jlc3BvbnNlEk8KE2hvc3RlZHpvbmVzdW1tYXJpZXMYrb'
    'WdNSADKAsyGi5yb3V0ZTUzLkhvc3RlZFpvbmVTdW1tYXJ5UhNob3N0ZWR6b25lc3VtbWFyaWVz'
    'Eh4KCG1heGl0ZW1zGJTW2vEBIAEoCVIIbWF4aXRlbXMSHwoJbmV4dHRva2VuGP6EumcgASgJUg'
    'luZXh0dG9rZW4=');

@$core.Deprecated('Use listHostedZonesRequestDescriptor instead')
const ListHostedZonesRequest$json = {
  '1': 'ListHostedZonesRequest',
  '2': [
    {
      '1': 'delegationsetid',
      '3': 307328801,
      '4': 1,
      '5': 9,
      '10': 'delegationsetid'
    },
    {
      '1': 'hostedzonetype',
      '3': 409319401,
      '4': 1,
      '5': 14,
      '6': '.route53.HostedZoneType',
      '10': 'hostedzonetype'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListHostedZonesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHostedZonesRequestDescriptor = $convert.base64Decode(
    'ChZMaXN0SG9zdGVkWm9uZXNSZXF1ZXN0EiwKD2RlbGVnYXRpb25zZXRpZBih7sWSASABKAlSD2'
    'RlbGVnYXRpb25zZXRpZBJDCg5ob3N0ZWR6b25ldHlwZRjp75bDASABKA4yFy5yb3V0ZTUzLkhv'
    'c3RlZFpvbmVUeXBlUg5ob3N0ZWR6b25ldHlwZRIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlch'
    'IeCghtYXhpdGVtcxiU1trxASABKAlSCG1heGl0ZW1z');

@$core.Deprecated('Use listHostedZonesResponseDescriptor instead')
const ListHostedZonesResponse$json = {
  '1': 'ListHostedZonesResponse',
  '2': [
    {
      '1': 'hostedzones',
      '3': 86735402,
      '4': 3,
      '5': 11,
      '6': '.route53.HostedZone',
      '10': 'hostedzones'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListHostedZonesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHostedZonesResponseDescriptor = $convert.base64Decode(
    'ChdMaXN0SG9zdGVkWm9uZXNSZXNwb25zZRI4Cgtob3N0ZWR6b25lcxiq9K0pIAMoCzITLnJvdX'
    'RlNTMuSG9zdGVkWm9uZVILaG9zdGVkem9uZXMSIwoLaXN0cnVuY2F0ZWQY2p+4cyABKAhSC2lz'
    'dHJ1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIA'
    'EoCVIIbWF4aXRlbXMSIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leHRtYXJrZXI=');

@$core.Deprecated('Use listQueryLoggingConfigsRequestDescriptor instead')
const ListQueryLoggingConfigsRequest$json = {
  '1': 'ListQueryLoggingConfigsRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 9, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListQueryLoggingConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueryLoggingConfigsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0UXVlcnlMb2dnaW5nQ29uZmlnc1JlcXVlc3QSJgoMaG9zdGVkem9uZWlkGP7OnqUBIA'
        'EoCVIMaG9zdGVkem9uZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgJUgptYXhyZXN1bHRzEh8K'
        'CW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listQueryLoggingConfigsResponseDescriptor instead')
const ListQueryLoggingConfigsResponse$json = {
  '1': 'ListQueryLoggingConfigsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queryloggingconfigs',
      '3': 87688172,
      '4': 3,
      '5': 11,
      '6': '.route53.QueryLoggingConfig',
      '10': 'queryloggingconfigs'
    },
  ],
};

/// Descriptor for `ListQueryLoggingConfigsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueryLoggingConfigsResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0UXVlcnlMb2dnaW5nQ29uZmlnc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCV'
        'IJbmV4dHRva2VuElAKE3F1ZXJ5bG9nZ2luZ2NvbmZpZ3MY7IfoKSADKAsyGy5yb3V0ZTUzLlF1'
        'ZXJ5TG9nZ2luZ0NvbmZpZ1ITcXVlcnlsb2dnaW5nY29uZmlncw==');

@$core.Deprecated('Use listResourceRecordSetsRequestDescriptor instead')
const ListResourceRecordSetsRequest$json = {
  '1': 'ListResourceRecordSetsRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'startrecordidentifier',
      '3': 518502950,
      '4': 1,
      '5': 9,
      '10': 'startrecordidentifier'
    },
    {
      '1': 'startrecordname',
      '3': 145299062,
      '4': 1,
      '5': 9,
      '10': 'startrecordname'
    },
    {
      '1': 'startrecordtype',
      '3': 408714791,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'startrecordtype'
    },
  ],
};

/// Descriptor for `ListResourceRecordSetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceRecordSetsRequestDescriptor = $convert.base64Decode(
    'Ch1MaXN0UmVzb3VyY2VSZWNvcmRTZXRzUmVxdWVzdBImCgxob3N0ZWR6b25laWQY/s6epQEgAS'
    'gJUgxob3N0ZWR6b25laWQSHgoIbWF4aXRlbXMYlNba8QEgASgJUghtYXhpdGVtcxI4ChVzdGFy'
    'dHJlY29yZGlkZW50aWZpZXIYpvSe9wEgASgJUhVzdGFydHJlY29yZGlkZW50aWZpZXISKwoPc3'
    'RhcnRyZWNvcmRuYW1lGPaspEUgASgJUg9zdGFydHJlY29yZG5hbWUSPQoPc3RhcnRyZWNvcmR0'
    'eXBlGKf88cIBIAEoDjIPLnJvdXRlNTMuUlJUeXBlUg9zdGFydHJlY29yZHR5cGU=');

@$core.Deprecated('Use listResourceRecordSetsResponseDescriptor instead')
const ListResourceRecordSetsResponse$json = {
  '1': 'ListResourceRecordSetsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'nextrecordidentifier',
      '3': 424069527,
      '4': 1,
      '5': 9,
      '10': 'nextrecordidentifier'
    },
    {
      '1': 'nextrecordname',
      '3': 131258783,
      '4': 1,
      '5': 9,
      '10': 'nextrecordname'
    },
    {
      '1': 'nextrecordtype',
      '3': 97817846,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'nextrecordtype'
    },
    {
      '1': 'resourcerecordsets',
      '3': 77807302,
      '4': 3,
      '5': 11,
      '6': '.route53.ResourceRecordSet',
      '10': 'resourcerecordsets'
    },
  ],
};

/// Descriptor for `ListResourceRecordSetsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceRecordSetsResponseDescriptor = $convert.base64Decode(
    'Ch5MaXN0UmVzb3VyY2VSZWNvcmRTZXRzUmVzcG9uc2USIwoLaXN0cnVuY2F0ZWQY2p+4cyABKA'
    'hSC2lzdHJ1bmNhdGVkEh4KCG1heGl0ZW1zGJTW2vEBIAEoCVIIbWF4aXRlbXMSNgoUbmV4dHJl'
    'Y29yZGlkZW50aWZpZXIYl5ObygEgASgJUhRuZXh0cmVjb3JkaWRlbnRpZmllchIpCg5uZXh0cm'
    'Vjb3JkbmFtZRifs8s+IAEoCVIObmV4dHJlY29yZG5hbWUSOgoObmV4dHJlY29yZHR5cGUY9qnS'
    'LiABKA4yDy5yb3V0ZTUzLlJSVHlwZVIObmV4dHJlY29yZHR5cGUSTQoScmVzb3VyY2VyZWNvcm'
    'RzZXRzGMb9jCUgAygLMhoucm91dGU1My5SZXNvdXJjZVJlY29yZFNldFIScmVzb3VyY2VyZWNv'
    'cmRzZXRz');

@$core.Deprecated('Use listReusableDelegationSetsRequestDescriptor instead')
const ListReusableDelegationSetsRequest$json = {
  '1': 'ListReusableDelegationSetsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListReusableDelegationSetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReusableDelegationSetsRequestDescriptor =
    $convert.base64Decode(
        'CiFMaXN0UmV1c2FibGVEZWxlZ2F0aW9uU2V0c1JlcXVlc3QSGQoGbWFya2VyGLjdzSogASgJUg'
        'ZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgJUghtYXhpdGVtcw==');

@$core.Deprecated('Use listReusableDelegationSetsResponseDescriptor instead')
const ListReusableDelegationSetsResponse$json = {
  '1': 'ListReusableDelegationSetsResponse',
  '2': [
    {
      '1': 'delegationsets',
      '3': 477592931,
      '4': 3,
      '5': 11,
      '6': '.route53.DelegationSet',
      '10': 'delegationsets'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListReusableDelegationSetsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReusableDelegationSetsResponseDescriptor =
    $convert.base64Decode(
        'CiJMaXN0UmV1c2FibGVEZWxlZ2F0aW9uU2V0c1Jlc3BvbnNlEkIKDmRlbGVnYXRpb25zZXRzGO'
        'P63eMBIAMoCzIWLnJvdXRlNTMuRGVsZWdhdGlvblNldFIOZGVsZWdhdGlvbnNldHMSIwoLaXN0'
        'cnVuY2F0ZWQY2p+4cyABKAhSC2lzdHJ1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2'
        'VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoCVIIbWF4aXRlbXMSIgoKbmV4dG1hcmtlchijga79ASAB'
        'KAlSCm5leHRtYXJrZXI=');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.route53.TagResourceType',
      '10': 'resourcetype'
    },
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcm'
        'Vzb3VyY2VpZBJACgxyZXNvdXJjZXR5cGUY3r7YjwEgASgOMhgucm91dGU1My5UYWdSZXNvdXJj'
        'ZVR5cGVSDHJlc291cmNldHlwZQ==');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {
      '1': 'resourcetagset',
      '3': 249268514,
      '4': 1,
      '5': 11,
      '6': '.route53.ResourceTagSet',
      '10': 'resourcetagset'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USQgoOcmVzb3VyY2V0YWdzZXQYopLudiABKA'
        'syFy5yb3V0ZTUzLlJlc291cmNlVGFnU2V0Ug5yZXNvdXJjZXRhZ3NldA==');

@$core.Deprecated('Use listTagsForResourcesRequestDescriptor instead')
const ListTagsForResourcesRequest$json = {
  '1': 'ListTagsForResourcesRequest',
  '2': [
    {'1': 'resourceids', '3': 23528154, '4': 3, '5': 9, '10': 'resourceids'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.route53.TagResourceType',
      '10': 'resourcetype'
    },
  ],
};

/// Descriptor for `ListTagsForResourcesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourcesRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlc1JlcXVlc3QSIwoLcmVzb3VyY2VpZHMY2oWcCyADKAlSC3'
        'Jlc291cmNlaWRzEkAKDHJlc291cmNldHlwZRjevtiPASABKA4yGC5yb3V0ZTUzLlRhZ1Jlc291'
        'cmNlVHlwZVIMcmVzb3VyY2V0eXBl');

@$core.Deprecated('Use listTagsForResourcesResponseDescriptor instead')
const ListTagsForResourcesResponse$json = {
  '1': 'ListTagsForResourcesResponse',
  '2': [
    {
      '1': 'resourcetagsets',
      '3': 362359831,
      '4': 3,
      '5': 11,
      '6': '.route53.ResourceTagSet',
      '10': 'resourcetagsets'
    },
  ],
};

/// Descriptor for `ListTagsForResourcesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourcesResponseDescriptor =
    $convert.base64Decode(
        'ChxMaXN0VGFnc0ZvclJlc291cmNlc1Jlc3BvbnNlEkUKD3Jlc291cmNldGFnc2V0cxiX2OSsAS'
        'ADKAsyFy5yb3V0ZTUzLlJlc291cmNlVGFnU2V0Ug9yZXNvdXJjZXRhZ3NldHM=');

@$core.Deprecated('Use listTrafficPoliciesRequestDescriptor instead')
const ListTrafficPoliciesRequest$json = {
  '1': 'ListTrafficPoliciesRequest',
  '2': [
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyidmarker',
      '3': 426883336,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyidmarker'
    },
  ],
};

/// Descriptor for `ListTrafficPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPoliciesRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VHJhZmZpY1BvbGljaWVzUmVxdWVzdBIeCghtYXhpdGVtcxiU1trxASABKAlSCG1heG'
        'l0ZW1zEjgKFXRyYWZmaWNwb2xpY3lpZG1hcmtlchiI8sbLASABKAlSFXRyYWZmaWNwb2xpY3lp'
        'ZG1hcmtlcg==');

@$core.Deprecated('Use listTrafficPoliciesResponseDescriptor instead')
const ListTrafficPoliciesResponse$json = {
  '1': 'ListTrafficPoliciesResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyidmarker',
      '3': 426883336,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyidmarker'
    },
    {
      '1': 'trafficpolicysummaries',
      '3': 206945717,
      '4': 3,
      '5': 11,
      '6': '.route53.TrafficPolicySummary',
      '10': 'trafficpolicysummaries'
    },
  ],
};

/// Descriptor for `ListTrafficPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPoliciesResponseDescriptor = $convert.base64Decode(
    'ChtMaXN0VHJhZmZpY1BvbGljaWVzUmVzcG9uc2USIwoLaXN0cnVuY2F0ZWQY2p+4cyABKAhSC2'
    'lzdHJ1bmNhdGVkEh4KCG1heGl0ZW1zGJTW2vEBIAEoCVIIbWF4aXRlbXMSOAoVdHJhZmZpY3Bv'
    'bGljeWlkbWFya2VyGIjyxssBIAEoCVIVdHJhZmZpY3BvbGljeWlkbWFya2VyElgKFnRyYWZmaW'
    'Nwb2xpY3lzdW1tYXJpZXMYtfvWYiADKAsyHS5yb3V0ZTUzLlRyYWZmaWNQb2xpY3lTdW1tYXJ5'
    'UhZ0cmFmZmljcG9saWN5c3VtbWFyaWVz');

@$core.Deprecated(
    'Use listTrafficPolicyInstancesByHostedZoneRequestDescriptor instead')
const ListTrafficPolicyInstancesByHostedZoneRequest$json = {
  '1': 'ListTrafficPolicyInstancesByHostedZoneRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyinstancenamemarker',
      '3': 136215379,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyinstancenamemarker'
    },
    {
      '1': 'trafficpolicyinstancetypemarker',
      '3': 30935978,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'trafficpolicyinstancetypemarker'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyInstancesByHostedZoneRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listTrafficPolicyInstancesByHostedZoneRequestDescriptor =
    $convert.base64Decode(
        'Ci1MaXN0VHJhZmZpY1BvbGljeUluc3RhbmNlc0J5SG9zdGVkWm9uZVJlcXVlc3QSJgoMaG9zdG'
        'Vkem9uZWlkGP7OnqUBIAEoCVIMaG9zdGVkem9uZWlkEh4KCG1heGl0ZW1zGJTW2vEBIAEoCVII'
        'bWF4aXRlbXMSSwofdHJhZmZpY3BvbGljeWluc3RhbmNlbmFtZW1hcmtlchjT9vlAIAEoCVIfdH'
        'JhZmZpY3BvbGljeWluc3RhbmNlbmFtZW1hcmtlchJcCh90cmFmZmljcG9saWN5aW5zdGFuY2V0'
        'eXBlbWFya2VyGKqX4A4gASgOMg8ucm91dGU1My5SUlR5cGVSH3RyYWZmaWNwb2xpY3lpbnN0YW'
        '5jZXR5cGVtYXJrZXI=');

@$core.Deprecated(
    'Use listTrafficPolicyInstancesByHostedZoneResponseDescriptor instead')
const ListTrafficPolicyInstancesByHostedZoneResponse$json = {
  '1': 'ListTrafficPolicyInstancesByHostedZoneResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyinstancenamemarker',
      '3': 136215379,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyinstancenamemarker'
    },
    {
      '1': 'trafficpolicyinstancetypemarker',
      '3': 30935978,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'trafficpolicyinstancetypemarker'
    },
    {
      '1': 'trafficpolicyinstances',
      '3': 199455009,
      '4': 3,
      '5': 11,
      '6': '.route53.TrafficPolicyInstance',
      '10': 'trafficpolicyinstances'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyInstancesByHostedZoneResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listTrafficPolicyInstancesByHostedZoneResponseDescriptor =
    $convert.base64Decode(
        'Ci5MaXN0VHJhZmZpY1BvbGljeUluc3RhbmNlc0J5SG9zdGVkWm9uZVJlc3BvbnNlEiMKC2lzdH'
        'J1bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBIeCghtYXhpdGVtcxiU1trxASABKAlSCG1h'
        'eGl0ZW1zEksKH3RyYWZmaWNwb2xpY3lpbnN0YW5jZW5hbWVtYXJrZXIY0/b5QCABKAlSH3RyYW'
        'ZmaWNwb2xpY3lpbnN0YW5jZW5hbWVtYXJrZXISXAofdHJhZmZpY3BvbGljeWluc3RhbmNldHlw'
        'ZW1hcmtlchiql+AOIAEoDjIPLnJvdXRlNTMuUlJUeXBlUh90cmFmZmljcG9saWN5aW5zdGFuY2'
        'V0eXBlbWFya2VyElkKFnRyYWZmaWNwb2xpY3lpbnN0YW5jZXMYoeKNXyADKAsyHi5yb3V0ZTUz'
        'LlRyYWZmaWNQb2xpY3lJbnN0YW5jZVIWdHJhZmZpY3BvbGljeWluc3RhbmNlcw==');

@$core.Deprecated(
    'Use listTrafficPolicyInstancesByPolicyRequestDescriptor instead')
const ListTrafficPolicyInstancesByPolicyRequest$json = {
  '1': 'ListTrafficPolicyInstancesByPolicyRequest',
  '2': [
    {
      '1': 'hostedzoneidmarker',
      '3': 475055952,
      '4': 1,
      '5': 9,
      '10': 'hostedzoneidmarker'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyid',
      '3': 40235222,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyid'
    },
    {
      '1': 'trafficpolicyinstancenamemarker',
      '3': 136215379,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyinstancenamemarker'
    },
    {
      '1': 'trafficpolicyinstancetypemarker',
      '3': 30935978,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'trafficpolicyinstancetypemarker'
    },
    {
      '1': 'trafficpolicyversion',
      '3': 479078485,
      '4': 1,
      '5': 5,
      '10': 'trafficpolicyversion'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyInstancesByPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPolicyInstancesByPolicyRequestDescriptor = $convert.base64Decode(
    'CilMaXN0VHJhZmZpY1BvbGljeUluc3RhbmNlc0J5UG9saWN5UmVxdWVzdBIyChJob3N0ZWR6b2'
    '5laWRtYXJrZXIY0I7D4gEgASgJUhJob3N0ZWR6b25laWRtYXJrZXISHgoIbWF4aXRlbXMYlNba'
    '8QEgASgJUghtYXhpdGVtcxIrCg90cmFmZmljcG9saWN5aWQY1uGXEyABKAlSD3RyYWZmaWNwb2'
    'xpY3lpZBJLCh90cmFmZmljcG9saWN5aW5zdGFuY2VuYW1lbWFya2VyGNP2+UAgASgJUh90cmFm'
    'ZmljcG9saWN5aW5zdGFuY2VuYW1lbWFya2VyElwKH3RyYWZmaWNwb2xpY3lpbnN0YW5jZXR5cG'
    'VtYXJrZXIYqpfgDiABKA4yDy5yb3V0ZTUzLlJSVHlwZVIfdHJhZmZpY3BvbGljeWluc3RhbmNl'
    'dHlwZW1hcmtlchI2ChR0cmFmZmljcG9saWN5dmVyc2lvbhjV0LjkASABKAVSFHRyYWZmaWNwb2'
    'xpY3l2ZXJzaW9u');

@$core.Deprecated(
    'Use listTrafficPolicyInstancesByPolicyResponseDescriptor instead')
const ListTrafficPolicyInstancesByPolicyResponse$json = {
  '1': 'ListTrafficPolicyInstancesByPolicyResponse',
  '2': [
    {
      '1': 'hostedzoneidmarker',
      '3': 475055952,
      '4': 1,
      '5': 9,
      '10': 'hostedzoneidmarker'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyinstancenamemarker',
      '3': 136215379,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyinstancenamemarker'
    },
    {
      '1': 'trafficpolicyinstancetypemarker',
      '3': 30935978,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'trafficpolicyinstancetypemarker'
    },
    {
      '1': 'trafficpolicyinstances',
      '3': 199455009,
      '4': 3,
      '5': 11,
      '6': '.route53.TrafficPolicyInstance',
      '10': 'trafficpolicyinstances'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyInstancesByPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPolicyInstancesByPolicyResponseDescriptor = $convert.base64Decode(
    'CipMaXN0VHJhZmZpY1BvbGljeUluc3RhbmNlc0J5UG9saWN5UmVzcG9uc2USMgoSaG9zdGVkem'
    '9uZWlkbWFya2VyGNCOw+IBIAEoCVISaG9zdGVkem9uZWlkbWFya2VyEiMKC2lzdHJ1bmNhdGVk'
    'GNqfuHMgASgIUgtpc3RydW5jYXRlZBIeCghtYXhpdGVtcxiU1trxASABKAlSCG1heGl0ZW1zEk'
    'sKH3RyYWZmaWNwb2xpY3lpbnN0YW5jZW5hbWVtYXJrZXIY0/b5QCABKAlSH3RyYWZmaWNwb2xp'
    'Y3lpbnN0YW5jZW5hbWVtYXJrZXISXAofdHJhZmZpY3BvbGljeWluc3RhbmNldHlwZW1hcmtlch'
    'iql+AOIAEoDjIPLnJvdXRlNTMuUlJUeXBlUh90cmFmZmljcG9saWN5aW5zdGFuY2V0eXBlbWFy'
    'a2VyElkKFnRyYWZmaWNwb2xpY3lpbnN0YW5jZXMYoeKNXyADKAsyHi5yb3V0ZTUzLlRyYWZmaW'
    'NQb2xpY3lJbnN0YW5jZVIWdHJhZmZpY3BvbGljeWluc3RhbmNlcw==');

@$core.Deprecated('Use listTrafficPolicyInstancesRequestDescriptor instead')
const ListTrafficPolicyInstancesRequest$json = {
  '1': 'ListTrafficPolicyInstancesRequest',
  '2': [
    {
      '1': 'hostedzoneidmarker',
      '3': 475055952,
      '4': 1,
      '5': 9,
      '10': 'hostedzoneidmarker'
    },
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyinstancenamemarker',
      '3': 136215379,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyinstancenamemarker'
    },
    {
      '1': 'trafficpolicyinstancetypemarker',
      '3': 30935978,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'trafficpolicyinstancetypemarker'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyInstancesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPolicyInstancesRequestDescriptor = $convert.base64Decode(
    'CiFMaXN0VHJhZmZpY1BvbGljeUluc3RhbmNlc1JlcXVlc3QSMgoSaG9zdGVkem9uZWlkbWFya2'
    'VyGNCOw+IBIAEoCVISaG9zdGVkem9uZWlkbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoCVII'
    'bWF4aXRlbXMSSwofdHJhZmZpY3BvbGljeWluc3RhbmNlbmFtZW1hcmtlchjT9vlAIAEoCVIfdH'
    'JhZmZpY3BvbGljeWluc3RhbmNlbmFtZW1hcmtlchJcCh90cmFmZmljcG9saWN5aW5zdGFuY2V0'
    'eXBlbWFya2VyGKqX4A4gASgOMg8ucm91dGU1My5SUlR5cGVSH3RyYWZmaWNwb2xpY3lpbnN0YW'
    '5jZXR5cGVtYXJrZXI=');

@$core.Deprecated('Use listTrafficPolicyInstancesResponseDescriptor instead')
const ListTrafficPolicyInstancesResponse$json = {
  '1': 'ListTrafficPolicyInstancesResponse',
  '2': [
    {
      '1': 'hostedzoneidmarker',
      '3': 475055952,
      '4': 1,
      '5': 9,
      '10': 'hostedzoneidmarker'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyinstancenamemarker',
      '3': 136215379,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyinstancenamemarker'
    },
    {
      '1': 'trafficpolicyinstancetypemarker',
      '3': 30935978,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'trafficpolicyinstancetypemarker'
    },
    {
      '1': 'trafficpolicyinstances',
      '3': 199455009,
      '4': 3,
      '5': 11,
      '6': '.route53.TrafficPolicyInstance',
      '10': 'trafficpolicyinstances'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyInstancesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPolicyInstancesResponseDescriptor = $convert.base64Decode(
    'CiJMaXN0VHJhZmZpY1BvbGljeUluc3RhbmNlc1Jlc3BvbnNlEjIKEmhvc3RlZHpvbmVpZG1hcm'
    'tlchjQjsPiASABKAlSEmhvc3RlZHpvbmVpZG1hcmtlchIjCgtpc3RydW5jYXRlZBjan7hzIAEo'
    'CFILaXN0cnVuY2F0ZWQSHgoIbWF4aXRlbXMYlNba8QEgASgJUghtYXhpdGVtcxJLCh90cmFmZm'
    'ljcG9saWN5aW5zdGFuY2VuYW1lbWFya2VyGNP2+UAgASgJUh90cmFmZmljcG9saWN5aW5zdGFu'
    'Y2VuYW1lbWFya2VyElwKH3RyYWZmaWNwb2xpY3lpbnN0YW5jZXR5cGVtYXJrZXIYqpfgDiABKA'
    '4yDy5yb3V0ZTUzLlJSVHlwZVIfdHJhZmZpY3BvbGljeWluc3RhbmNldHlwZW1hcmtlchJZChZ0'
    'cmFmZmljcG9saWN5aW5zdGFuY2VzGKHijV8gAygLMh4ucm91dGU1My5UcmFmZmljUG9saWN5SW'
    '5zdGFuY2VSFnRyYWZmaWNwb2xpY3lpbnN0YW5jZXM=');

@$core.Deprecated('Use listTrafficPolicyVersionsRequestDescriptor instead')
const ListTrafficPolicyVersionsRequest$json = {
  '1': 'ListTrafficPolicyVersionsRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicyversionmarker',
      '3': 228559295,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyversionmarker'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyVersionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPolicyVersionsRequestDescriptor =
    $convert.base64Decode(
        'CiBMaXN0VHJhZmZpY1BvbGljeVZlcnNpb25zUmVxdWVzdBISCgJpZBiB8qK3ASABKAlSAmlkEh'
        '4KCG1heGl0ZW1zGJTW2vEBIAEoCVIIbWF4aXRlbXMSQQoadHJhZmZpY3BvbGljeXZlcnNpb25t'
        'YXJrZXIYv5P+bCABKAlSGnRyYWZmaWNwb2xpY3l2ZXJzaW9ubWFya2Vy');

@$core.Deprecated('Use listTrafficPolicyVersionsResponseDescriptor instead')
const ListTrafficPolicyVersionsResponse$json = {
  '1': 'ListTrafficPolicyVersionsResponse',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 9, '10': 'maxitems'},
    {
      '1': 'trafficpolicies',
      '3': 111427421,
      '4': 3,
      '5': 11,
      '6': '.route53.TrafficPolicy',
      '10': 'trafficpolicies'
    },
    {
      '1': 'trafficpolicyversionmarker',
      '3': 228559295,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyversionmarker'
    },
  ],
};

/// Descriptor for `ListTrafficPolicyVersionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrafficPolicyVersionsResponseDescriptor = $convert.base64Decode(
    'CiFMaXN0VHJhZmZpY1BvbGljeVZlcnNpb25zUmVzcG9uc2USIwoLaXN0cnVuY2F0ZWQY2p+4cy'
    'ABKAhSC2lzdHJ1bmNhdGVkEh4KCG1heGl0ZW1zGJTW2vEBIAEoCVIIbWF4aXRlbXMSQwoPdHJh'
    'ZmZpY3BvbGljaWVzGN3+kDUgAygLMhYucm91dGU1My5UcmFmZmljUG9saWN5Ug90cmFmZmljcG'
    '9saWNpZXMSQQoadHJhZmZpY3BvbGljeXZlcnNpb25tYXJrZXIYv5P+bCABKAlSGnRyYWZmaWNw'
    'b2xpY3l2ZXJzaW9ubWFya2Vy');

@$core
    .Deprecated('Use listVPCAssociationAuthorizationsRequestDescriptor instead')
const ListVPCAssociationAuthorizationsRequest$json = {
  '1': 'ListVPCAssociationAuthorizationsRequest',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 9, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListVPCAssociationAuthorizationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVPCAssociationAuthorizationsRequestDescriptor =
    $convert.base64Decode(
        'CidMaXN0VlBDQXNzb2NpYXRpb25BdXRob3JpemF0aW9uc1JlcXVlc3QSJgoMaG9zdGVkem9uZW'
        'lkGP7OnqUBIAEoCVIMaG9zdGVkem9uZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgJUgptYXhy'
        'ZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated(
    'Use listVPCAssociationAuthorizationsResponseDescriptor instead')
const ListVPCAssociationAuthorizationsResponse$json = {
  '1': 'ListVPCAssociationAuthorizationsResponse',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'vpcs',
      '3': 424064898,
      '4': 3,
      '5': 11,
      '6': '.route53.VPC',
      '10': 'vpcs'
    },
  ],
};

/// Descriptor for `ListVPCAssociationAuthorizationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVPCAssociationAuthorizationsResponseDescriptor =
    $convert.base64Decode(
        'CihMaXN0VlBDQXNzb2NpYXRpb25BdXRob3JpemF0aW9uc1Jlc3BvbnNlEiYKDGhvc3RlZHpvbm'
        'VpZBj+zp6lASABKAlSDGhvc3RlZHpvbmVpZBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0'
        'b2tlbhIkCgR2cGNzGILvmsoBIAMoCzIMLnJvdXRlNTMuVlBDUgR2cGNz');

@$core.Deprecated('Use locationSummaryDescriptor instead')
const LocationSummary$json = {
  '1': 'LocationSummary',
  '2': [
    {'1': 'locationname', '3': 158186566, '4': 1, '5': 9, '10': 'locationname'},
  ],
};

/// Descriptor for `LocationSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locationSummaryDescriptor = $convert.base64Decode(
    'Cg9Mb2NhdGlvblN1bW1hcnkSJQoMbG9jYXRpb25uYW1lGMb4tksgASgJUgxsb2NhdGlvbm5hbW'
    'U=');

@$core.Deprecated('Use noSuchChangeDescriptor instead')
const NoSuchChange$json = {
  '1': 'NoSuchChange',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchChange`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchChangeDescriptor = $convert.base64Decode(
    'CgxOb1N1Y2hDaGFuZ2USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use noSuchCidrCollectionExceptionDescriptor instead')
const NoSuchCidrCollectionException$json = {
  '1': 'NoSuchCidrCollectionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchCidrCollectionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchCidrCollectionExceptionDescriptor =
    $convert.base64Decode(
        'Ch1Ob1N1Y2hDaWRyQ29sbGVjdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use noSuchCidrLocationExceptionDescriptor instead')
const NoSuchCidrLocationException$json = {
  '1': 'NoSuchCidrLocationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchCidrLocationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchCidrLocationExceptionDescriptor =
    $convert.base64Decode(
        'ChtOb1N1Y2hDaWRyTG9jYXRpb25FeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use noSuchCloudWatchLogsLogGroupDescriptor instead')
const NoSuchCloudWatchLogsLogGroup$json = {
  '1': 'NoSuchCloudWatchLogsLogGroup',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchCloudWatchLogsLogGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchCloudWatchLogsLogGroupDescriptor =
    $convert.base64Decode(
        'ChxOb1N1Y2hDbG91ZFdhdGNoTG9nc0xvZ0dyb3VwEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use noSuchDelegationSetDescriptor instead')
const NoSuchDelegationSet$json = {
  '1': 'NoSuchDelegationSet',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchDelegationSet`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchDelegationSetDescriptor =
    $convert.base64Decode(
        'ChNOb1N1Y2hEZWxlZ2F0aW9uU2V0EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use noSuchGeoLocationDescriptor instead')
const NoSuchGeoLocation$json = {
  '1': 'NoSuchGeoLocation',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchGeoLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchGeoLocationDescriptor = $convert.base64Decode(
    'ChFOb1N1Y2hHZW9Mb2NhdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use noSuchHealthCheckDescriptor instead')
const NoSuchHealthCheck$json = {
  '1': 'NoSuchHealthCheck',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchHealthCheck`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchHealthCheckDescriptor = $convert.base64Decode(
    'ChFOb1N1Y2hIZWFsdGhDaGVjaxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use noSuchHostedZoneDescriptor instead')
const NoSuchHostedZone$json = {
  '1': 'NoSuchHostedZone',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchHostedZone`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchHostedZoneDescriptor = $convert.base64Decode(
    'ChBOb1N1Y2hIb3N0ZWRab25lEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use noSuchKeySigningKeyDescriptor instead')
const NoSuchKeySigningKey$json = {
  '1': 'NoSuchKeySigningKey',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchKeySigningKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchKeySigningKeyDescriptor =
    $convert.base64Decode(
        'ChNOb1N1Y2hLZXlTaWduaW5nS2V5EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use noSuchQueryLoggingConfigDescriptor instead')
const NoSuchQueryLoggingConfig$json = {
  '1': 'NoSuchQueryLoggingConfig',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchQueryLoggingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchQueryLoggingConfigDescriptor =
    $convert.base64Decode(
        'ChhOb1N1Y2hRdWVyeUxvZ2dpbmdDb25maWcSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use noSuchTrafficPolicyDescriptor instead')
const NoSuchTrafficPolicy$json = {
  '1': 'NoSuchTrafficPolicy',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchTrafficPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchTrafficPolicyDescriptor =
    $convert.base64Decode(
        'ChNOb1N1Y2hUcmFmZmljUG9saWN5EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use noSuchTrafficPolicyInstanceDescriptor instead')
const NoSuchTrafficPolicyInstance$json = {
  '1': 'NoSuchTrafficPolicyInstance',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoSuchTrafficPolicyInstance`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchTrafficPolicyInstanceDescriptor =
    $convert.base64Decode(
        'ChtOb1N1Y2hUcmFmZmljUG9saWN5SW5zdGFuY2USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use notAuthorizedExceptionDescriptor instead')
const NotAuthorizedException$json = {
  '1': 'NotAuthorizedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NotAuthorizedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notAuthorizedExceptionDescriptor =
    $convert.base64Decode(
        'ChZOb3RBdXRob3JpemVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use priorRequestNotCompleteDescriptor instead')
const PriorRequestNotComplete$json = {
  '1': 'PriorRequestNotComplete',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PriorRequestNotComplete`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List priorRequestNotCompleteDescriptor =
    $convert.base64Decode(
        'ChdQcmlvclJlcXVlc3ROb3RDb21wbGV0ZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use publicZoneVPCAssociationDescriptor instead')
const PublicZoneVPCAssociation$json = {
  '1': 'PublicZoneVPCAssociation',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PublicZoneVPCAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicZoneVPCAssociationDescriptor =
    $convert.base64Decode(
        'ChhQdWJsaWNab25lVlBDQXNzb2NpYXRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use queryLoggingConfigDescriptor instead')
const QueryLoggingConfig$json = {
  '1': 'QueryLoggingConfig',
  '2': [
    {
      '1': 'cloudwatchlogsloggrouparn',
      '3': 519613355,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsloggrouparn'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `QueryLoggingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryLoggingConfigDescriptor = $convert.base64Decode(
    'ChJRdWVyeUxvZ2dpbmdDb25maWcSQAoZY2xvdWR3YXRjaGxvZ3Nsb2dncm91cGFybhir1+L3AS'
    'ABKAlSGWNsb3Vkd2F0Y2hsb2dzbG9nZ3JvdXBhcm4SJgoMaG9zdGVkem9uZWlkGP7OnqUBIAEo'
    'CVIMaG9zdGVkem9uZWlkEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use queryLoggingConfigAlreadyExistsDescriptor instead')
const QueryLoggingConfigAlreadyExists$json = {
  '1': 'QueryLoggingConfigAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueryLoggingConfigAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryLoggingConfigAlreadyExistsDescriptor =
    $convert.base64Decode(
        'Ch9RdWVyeUxvZ2dpbmdDb25maWdBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use resourceRecordDescriptor instead')
const ResourceRecord$json = {
  '1': 'ResourceRecord',
  '2': [
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ResourceRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceRecordDescriptor = $convert
    .base64Decode('Cg5SZXNvdXJjZVJlY29yZBIYCgV2YWx1ZRjr8p+KASABKAlSBXZhbHVl');

@$core.Deprecated('Use resourceRecordSetDescriptor instead')
const ResourceRecordSet$json = {
  '1': 'ResourceRecordSet',
  '2': [
    {
      '1': 'aliastarget',
      '3': 317299867,
      '4': 1,
      '5': 11,
      '6': '.route53.AliasTarget',
      '10': 'aliastarget'
    },
    {
      '1': 'cidrroutingconfig',
      '3': 357245246,
      '4': 1,
      '5': 11,
      '6': '.route53.CidrRoutingConfig',
      '10': 'cidrroutingconfig'
    },
    {
      '1': 'failover',
      '3': 26793064,
      '4': 1,
      '5': 14,
      '6': '.route53.ResourceRecordSetFailover',
      '10': 'failover'
    },
    {
      '1': 'geolocation',
      '3': 267973346,
      '4': 1,
      '5': 11,
      '6': '.route53.GeoLocation',
      '10': 'geolocation'
    },
    {
      '1': 'geoproximitylocation',
      '3': 94319785,
      '4': 1,
      '5': 11,
      '6': '.route53.GeoProximityLocation',
      '10': 'geoproximitylocation'
    },
    {
      '1': 'healthcheckid',
      '3': 312971637,
      '4': 1,
      '5': 9,
      '10': 'healthcheckid'
    },
    {
      '1': 'multivalueanswer',
      '3': 424105486,
      '4': 1,
      '5': 8,
      '10': 'multivalueanswer'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'region',
      '3': 154040478,
      '4': 1,
      '5': 14,
      '6': '.route53.ResourceRecordSetRegion',
      '10': 'region'
    },
    {
      '1': 'resourcerecords',
      '3': 519418974,
      '4': 3,
      '5': 11,
      '6': '.route53.ResourceRecord',
      '10': 'resourcerecords'
    },
    {
      '1': 'setidentifier',
      '3': 201408985,
      '4': 1,
      '5': 9,
      '10': 'setidentifier'
    },
    {'1': 'ttl', '3': 526904700, '4': 1, '5': 3, '10': 'ttl'},
    {
      '1': 'trafficpolicyinstanceid',
      '3': 251421439,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyinstanceid'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'type'
    },
    {'1': 'weight', '3': 422581466, '4': 1, '5': 3, '10': 'weight'},
  ],
};

/// Descriptor for `ResourceRecordSet`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceRecordSetDescriptor = $convert.base64Decode(
    'ChFSZXNvdXJjZVJlY29yZFNldBI6CgthbGlhc3RhcmdldBibuaaXASABKAsyFC5yb3V0ZTUzLk'
    'FsaWFzVGFyZ2V0UgthbGlhc3RhcmdldBJMChFjaWRycm91dGluZ2NvbmZpZxi+wqyqASABKAsy'
    'Gi5yb3V0ZTUzLkNpZHJSb3V0aW5nQ29uZmlnUhFjaWRycm91dGluZ2NvbmZpZxJBCghmYWlsb3'
    'ZlchjoqOMMIAEoDjIiLnJvdXRlNTMuUmVzb3VyY2VSZWNvcmRTZXRGYWlsb3ZlclIIZmFpbG92'
    'ZXISOQoLZ2VvbG9jYXRpb24Y4uXjfyABKAsyFC5yb3V0ZTUzLkdlb0xvY2F0aW9uUgtnZW9sb2'
    'NhdGlvbhJUChRnZW9wcm94aW1pdHlsb2NhdGlvbhip6fwsIAEoCzIdLnJvdXRlNTMuR2VvUHJv'
    'eGltaXR5TG9jYXRpb25SFGdlb3Byb3hpbWl0eWxvY2F0aW9uEigKDWhlYWx0aGNoZWNraWQY9a'
    'KelQEgASgJUg1oZWFsdGhjaGVja2lkEi4KEG11bHRpdmFsdWVhbnN3ZXIYjqydygEgASgIUhBt'
    'dWx0aXZhbHVlYW5zd2VyEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSOwoGcmVnaW9uGJ7xuUkgAS'
    'gOMiAucm91dGU1My5SZXNvdXJjZVJlY29yZFNldFJlZ2lvblIGcmVnaW9uEkUKD3Jlc291cmNl'
    'cmVjb3Jkcxje6Nb3ASADKAsyFy5yb3V0ZTUzLlJlc291cmNlUmVjb3JkUg9yZXNvdXJjZXJlY2'
    '9yZHMSJwoNc2V0aWRlbnRpZmllchjZg4VgIAEoCVINc2V0aWRlbnRpZmllchIUCgN0dGwY/Nqf'
    '+wEgASgDUgN0dGwSOwoXdHJhZmZpY3BvbGljeWluc3RhbmNlaWQY/8XxdyABKAlSF3RyYWZmaW'
    'Nwb2xpY3lpbnN0YW5jZWlkEicKBHR5cGUY7qDXigEgASgOMg8ucm91dGU1My5SUlR5cGVSBHR5'
    'cGUSGgoGd2VpZ2h0GNqpwMkBIAEoA1IGd2VpZ2h0');

@$core.Deprecated('Use resourceTagSetDescriptor instead')
const ResourceTagSet$json = {
  '1': 'ResourceTagSet',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.route53.TagResourceType',
      '10': 'resourcetype'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.route53.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ResourceTagSet`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceTagSetDescriptor = $convert.base64Decode(
    'Cg5SZXNvdXJjZVRhZ1NldBIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcmVzb3VyY2VpZBJACg'
    'xyZXNvdXJjZXR5cGUY3r7YjwEgASgOMhgucm91dGU1My5UYWdSZXNvdXJjZVR5cGVSDHJlc291'
    'cmNldHlwZRIkCgR0YWdzGMHB9rUBIAMoCzIMLnJvdXRlNTMuVGFnUgR0YWdz');

@$core.Deprecated('Use reusableDelegationSetLimitDescriptor instead')
const ReusableDelegationSetLimit$json = {
  '1': 'ReusableDelegationSetLimit',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.ReusableDelegationSetLimitType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 3, '10': 'value'},
  ],
};

/// Descriptor for `ReusableDelegationSetLimit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reusableDelegationSetLimitDescriptor =
    $convert.base64Decode(
        'ChpSZXVzYWJsZURlbGVnYXRpb25TZXRMaW1pdBI/CgR0eXBlGO6g14oBIAEoDjInLnJvdXRlNT'
        'MuUmV1c2FibGVEZWxlZ2F0aW9uU2V0TGltaXRUeXBlUgR0eXBlEhgKBXZhbHVlGOvyn4oBIAEo'
        'A1IFdmFsdWU=');

@$core.Deprecated('Use statusReportDescriptor instead')
const StatusReport$json = {
  '1': 'StatusReport',
  '2': [
    {'1': 'checkedtime', '3': 152550560, '4': 1, '5': 9, '10': 'checkedtime'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `StatusReport`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statusReportDescriptor = $convert.base64Decode(
    'CgxTdGF0dXNSZXBvcnQSIwoLY2hlY2tlZHRpbWUYoPneSCABKAlSC2NoZWNrZWR0aW1lEhkKBn'
    'N0YXR1cxiQ5PsCIAEoCVIGc3RhdHVz');

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

@$core.Deprecated('Use testDNSAnswerRequestDescriptor instead')
const TestDNSAnswerRequest$json = {
  '1': 'TestDNSAnswerRequest',
  '2': [
    {
      '1': 'edns0clientsubnetip',
      '3': 506730999,
      '4': 1,
      '5': 9,
      '10': 'edns0clientsubnetip'
    },
    {
      '1': 'edns0clientsubnetmask',
      '3': 381322620,
      '4': 1,
      '5': 9,
      '10': 'edns0clientsubnetmask'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'recordname', '3': 204939016, '4': 1, '5': 9, '10': 'recordname'},
    {
      '1': 'recordtype',
      '3': 441248261,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'recordtype'
    },
    {'1': 'resolverip', '3': 95000907, '4': 1, '5': 9, '10': 'resolverip'},
  ],
};

/// Descriptor for `TestDNSAnswerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testDNSAnswerRequestDescriptor = $convert.base64Decode(
    'ChRUZXN0RE5TQW5zd2VyUmVxdWVzdBI0ChNlZG5zMGNsaWVudHN1Ym5ldGlwGPez0PEBIAEoCV'
    'ITZWRuczBjbGllbnRzdWJuZXRpcBI4ChVlZG5zMGNsaWVudHN1Ym5ldG1hc2sY/IrqtQEgASgJ'
    'UhVlZG5zMGNsaWVudHN1Ym5ldG1hc2sSJgoMaG9zdGVkem9uZWlkGP7OnqUBIAEoCVIMaG9zdG'
    'Vkem9uZWlkEiEKCnJlY29yZG5hbWUYiL7cYSABKAlSCnJlY29yZG5hbWUSMwoKcmVjb3JkdHlw'
    'ZRiF1LPSASABKA4yDy5yb3V0ZTUzLlJSVHlwZVIKcmVjb3JkdHlwZRIhCgpyZXNvbHZlcmlwGM'
    'uypi0gASgJUgpyZXNvbHZlcmlw');

@$core.Deprecated('Use testDNSAnswerResponseDescriptor instead')
const TestDNSAnswerResponse$json = {
  '1': 'TestDNSAnswerResponse',
  '2': [
    {'1': 'nameserver', '3': 474180450, '4': 1, '5': 9, '10': 'nameserver'},
    {'1': 'protocol', '3': 173534166, '4': 1, '5': 9, '10': 'protocol'},
    {'1': 'recorddata', '3': 527821973, '4': 3, '5': 9, '10': 'recorddata'},
    {'1': 'recordname', '3': 204939016, '4': 1, '5': 9, '10': 'recordname'},
    {
      '1': 'recordtype',
      '3': 441248261,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'recordtype'
    },
    {'1': 'responsecode', '3': 447553700, '4': 1, '5': 9, '10': 'responsecode'},
  ],
};

/// Descriptor for `TestDNSAnswerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testDNSAnswerResponseDescriptor = $convert.base64Decode(
    'ChVUZXN0RE5TQW5zd2VyUmVzcG9uc2USIgoKbmFtZXNlcnZlchji1o3iASABKAlSCm5hbWVzZX'
    'J2ZXISHQoIcHJvdG9jb2wY1tffUiABKAlSCHByb3RvY29sEiIKCnJlY29yZGRhdGEYldnX+wEg'
    'AygJUgpyZWNvcmRkYXRhEiEKCnJlY29yZG5hbWUYiL7cYSABKAlSCnJlY29yZG5hbWUSMwoKcm'
    'Vjb3JkdHlwZRiF1LPSASABKA4yDy5yb3V0ZTUzLlJSVHlwZVIKcmVjb3JkdHlwZRImCgxyZXNw'
    'b25zZWNvZGUYpMG01QEgASgJUgxyZXNwb25zZWNvZGU=');

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

@$core.Deprecated('Use tooManyHealthChecksDescriptor instead')
const TooManyHealthChecks$json = {
  '1': 'TooManyHealthChecks',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyHealthChecks`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyHealthChecksDescriptor =
    $convert.base64Decode(
        'ChNUb29NYW55SGVhbHRoQ2hlY2tzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyHostedZonesDescriptor instead')
const TooManyHostedZones$json = {
  '1': 'TooManyHostedZones',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyHostedZones`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyHostedZonesDescriptor =
    $convert.base64Decode(
        'ChJUb29NYW55SG9zdGVkWm9uZXMSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use tooManyKeySigningKeysDescriptor instead')
const TooManyKeySigningKeys$json = {
  '1': 'TooManyKeySigningKeys',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyKeySigningKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyKeySigningKeysDescriptor = $convert.base64Decode(
    'ChVUb29NYW55S2V5U2lnbmluZ0tleXMSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use tooManyTrafficPoliciesDescriptor instead')
const TooManyTrafficPolicies$json = {
  '1': 'TooManyTrafficPolicies',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyTrafficPolicies`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyTrafficPoliciesDescriptor =
    $convert.base64Decode(
        'ChZUb29NYW55VHJhZmZpY1BvbGljaWVzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tooManyTrafficPolicyInstancesDescriptor instead')
const TooManyTrafficPolicyInstances$json = {
  '1': 'TooManyTrafficPolicyInstances',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyTrafficPolicyInstances`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyTrafficPolicyInstancesDescriptor =
    $convert.base64Decode(
        'Ch1Ub29NYW55VHJhZmZpY1BvbGljeUluc3RhbmNlcxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated(
    'Use tooManyTrafficPolicyVersionsForCurrentPolicyDescriptor instead')
const TooManyTrafficPolicyVersionsForCurrentPolicy$json = {
  '1': 'TooManyTrafficPolicyVersionsForCurrentPolicy',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyTrafficPolicyVersionsForCurrentPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    tooManyTrafficPolicyVersionsForCurrentPolicyDescriptor =
    $convert.base64Decode(
        'CixUb29NYW55VHJhZmZpY1BvbGljeVZlcnNpb25zRm9yQ3VycmVudFBvbGljeRIbCgdtZXNzYW'
        'dlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use tooManyVPCAssociationAuthorizationsDescriptor instead')
const TooManyVPCAssociationAuthorizations$json = {
  '1': 'TooManyVPCAssociationAuthorizations',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyVPCAssociationAuthorizations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyVPCAssociationAuthorizationsDescriptor =
    $convert.base64Decode(
        'CiNUb29NYW55VlBDQXNzb2NpYXRpb25BdXRob3JpemF0aW9ucxIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use trafficPolicyDescriptor instead')
const TrafficPolicy$json = {
  '1': 'TrafficPolicy',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'document', '3': 407108341, '4': 1, '5': 9, '10': 'document'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'type'
    },
    {'1': 'version', '3': 500028728, '4': 1, '5': 5, '10': 'version'},
  ],
};

/// Descriptor for `TrafficPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trafficPolicyDescriptor = $convert.base64Decode(
    'Cg1UcmFmZmljUG9saWN5EhwKB2NvbW1lbnQY/7++wgEgASgJUgdjb21tZW50Eh4KCGRvY3VtZW'
    '50GPX1j8IBIAEoCVIIZG9jdW1lbnQSEgoCaWQYgfKitwEgASgJUgJpZBIVCgRuYW1lGIfmgX8g'
    'ASgJUgRuYW1lEicKBHR5cGUY7qDXigEgASgOMg8ucm91dGU1My5SUlR5cGVSBHR5cGUSHAoHdm'
    'Vyc2lvbhi4qrfuASABKAVSB3ZlcnNpb24=');

@$core.Deprecated('Use trafficPolicyAlreadyExistsDescriptor instead')
const TrafficPolicyAlreadyExists$json = {
  '1': 'TrafficPolicyAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrafficPolicyAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trafficPolicyAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChpUcmFmZmljUG9saWN5QWxyZWFkeUV4aXN0cxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use trafficPolicyInUseDescriptor instead')
const TrafficPolicyInUse$json = {
  '1': 'TrafficPolicyInUse',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrafficPolicyInUse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trafficPolicyInUseDescriptor =
    $convert.base64Decode(
        'ChJUcmFmZmljUG9saWN5SW5Vc2USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use trafficPolicyInstanceDescriptor instead')
const TrafficPolicyInstance$json = {
  '1': 'TrafficPolicyInstance',
  '2': [
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {'1': 'ttl', '3': 526904700, '4': 1, '5': 3, '10': 'ttl'},
    {
      '1': 'trafficpolicyid',
      '3': 40235222,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyid'
    },
    {
      '1': 'trafficpolicytype',
      '3': 214345537,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'trafficpolicytype'
    },
    {
      '1': 'trafficpolicyversion',
      '3': 479078485,
      '4': 1,
      '5': 5,
      '10': 'trafficpolicyversion'
    },
  ],
};

/// Descriptor for `TrafficPolicyInstance`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trafficPolicyInstanceDescriptor = $convert.base64Decode(
    'ChVUcmFmZmljUG9saWN5SW5zdGFuY2USJgoMaG9zdGVkem9uZWlkGP7OnqUBIAEoCVIMaG9zdG'
    'Vkem9uZWlkEhIKAmlkGIHyorcBIAEoCVICaWQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2Fn'
    'ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhgKBXN0YXRlGJfJsu8BIAEoCVIFc3RhdGUSFAoDdH'
    'RsGPzan/sBIAEoA1IDdHRsEisKD3RyYWZmaWNwb2xpY3lpZBjW4ZcTIAEoCVIPdHJhZmZpY3Bv'
    'bGljeWlkEkAKEXRyYWZmaWNwb2xpY3l0eXBlGMHOmmYgASgOMg8ucm91dGU1My5SUlR5cGVSEX'
    'RyYWZmaWNwb2xpY3l0eXBlEjYKFHRyYWZmaWNwb2xpY3l2ZXJzaW9uGNXQuOQBIAEoBVIUdHJh'
    'ZmZpY3BvbGljeXZlcnNpb24=');

@$core.Deprecated('Use trafficPolicyInstanceAlreadyExistsDescriptor instead')
const TrafficPolicyInstanceAlreadyExists$json = {
  '1': 'TrafficPolicyInstanceAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrafficPolicyInstanceAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trafficPolicyInstanceAlreadyExistsDescriptor =
    $convert.base64Decode(
        'CiJUcmFmZmljUG9saWN5SW5zdGFuY2VBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use trafficPolicySummaryDescriptor instead')
const TrafficPolicySummary$json = {
  '1': 'TrafficPolicySummary',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'latestversion',
      '3': 424864587,
      '4': 1,
      '5': 5,
      '10': 'latestversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'trafficpolicycount',
      '3': 157833448,
      '4': 1,
      '5': 5,
      '10': 'trafficpolicycount'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.route53.RRType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `TrafficPolicySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trafficPolicySummaryDescriptor = $convert.base64Decode(
    'ChRUcmFmZmljUG9saWN5U3VtbWFyeRISCgJpZBiB8qK3ASABKAlSAmlkEigKDWxhdGVzdHZlcn'
    'Npb24Yy9bLygEgASgFUg1sYXRlc3R2ZXJzaW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSMQoS'
    'dHJhZmZpY3BvbGljeWNvdW50GOixoUsgASgFUhJ0cmFmZmljcG9saWN5Y291bnQSJwoEdHlwZR'
    'juoNeKASABKA4yDy5yb3V0ZTUzLlJSVHlwZVIEdHlwZQ==');

@$core.Deprecated('Use updateHealthCheckRequestDescriptor instead')
const UpdateHealthCheckRequest$json = {
  '1': 'UpdateHealthCheckRequest',
  '2': [
    {
      '1': 'alarmidentifier',
      '3': 536124346,
      '4': 1,
      '5': 11,
      '6': '.route53.AlarmIdentifier',
      '10': 'alarmidentifier'
    },
    {
      '1': 'childhealthchecks',
      '3': 485535935,
      '4': 3,
      '5': 9,
      '10': 'childhealthchecks'
    },
    {'1': 'disabled', '3': 533633318, '4': 1, '5': 8, '10': 'disabled'},
    {'1': 'enablesni', '3': 70122887, '4': 1, '5': 8, '10': 'enablesni'},
    {
      '1': 'failurethreshold',
      '3': 176846565,
      '4': 1,
      '5': 5,
      '10': 'failurethreshold'
    },
    {
      '1': 'fullyqualifieddomainname',
      '3': 459321509,
      '4': 1,
      '5': 9,
      '10': 'fullyqualifieddomainname'
    },
    {
      '1': 'healthcheckid',
      '3': 312971637,
      '4': 1,
      '5': 9,
      '10': 'healthcheckid'
    },
    {
      '1': 'healthcheckversion',
      '3': 89568396,
      '4': 1,
      '5': 3,
      '10': 'healthcheckversion'
    },
    {
      '1': 'healththreshold',
      '3': 215873163,
      '4': 1,
      '5': 5,
      '10': 'healththreshold'
    },
    {'1': 'ipaddress', '3': 169333741, '4': 1, '5': 9, '10': 'ipaddress'},
    {
      '1': 'insufficientdatahealthstatus',
      '3': 493115723,
      '4': 1,
      '5': 14,
      '6': '.route53.InsufficientDataHealthStatus',
      '10': 'insufficientdatahealthstatus'
    },
    {'1': 'inverted', '3': 55175513, '4': 1, '5': 8, '10': 'inverted'},
    {'1': 'port', '3': 46480583, '4': 1, '5': 5, '10': 'port'},
    {
      '1': 'regions',
      '3': 36200107,
      '4': 3,
      '5': 14,
      '6': '.route53.HealthCheckRegion',
      '10': 'regions'
    },
    {
      '1': 'resetelements',
      '3': 16543458,
      '4': 3,
      '5': 14,
      '6': '.route53.ResettableElementName',
      '10': 'resetelements'
    },
    {'1': 'resourcepath', '3': 117584551, '4': 1, '5': 9, '10': 'resourcepath'},
    {'1': 'searchstring', '3': 318687365, '4': 1, '5': 9, '10': 'searchstring'},
  ],
};

/// Descriptor for `UpdateHealthCheckRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateHealthCheckRequestDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVIZWFsdGhDaGVja1JlcXVlc3QSRgoPYWxhcm1pZGVudGlmaWVyGLq30v8BIAEoCz'
    'IYLnJvdXRlNTMuQWxhcm1JZGVudGlmaWVyUg9hbGFybWlkZW50aWZpZXISMAoRY2hpbGRoZWFs'
    'dGhjaGVja3MYv+HC5wEgAygJUhFjaGlsZGhlYWx0aGNoZWNrcxIeCghkaXNhYmxlZBimsrr+AS'
    'ABKAhSCGRpc2FibGVkEh8KCWVuYWJsZXNuaRiH+7chIAEoCFIJZW5hYmxlc25pEi0KEGZhaWx1'
    'cmV0aHJlc2hvbGQY5e2pVCABKAVSEGZhaWx1cmV0aHJlc2hvbGQSPgoYZnVsbHlxdWFsaWZpZW'
    'Rkb21haW5uYW1lGKXhgtsBIAEoCVIYZnVsbHlxdWFsaWZpZWRkb21haW5uYW1lEigKDWhlYWx0'
    'aGNoZWNraWQY9aKelQEgASgJUg1oZWFsdGhjaGVja2lkEjEKEmhlYWx0aGNoZWNrdmVyc2lvbh'
    'iM6doqIAEoA1ISaGVhbHRoY2hlY2t2ZXJzaW9uEisKD2hlYWx0aHRocmVzaG9sZBiL7fdmIAEo'
    'BVIPaGVhbHRodGhyZXNob2xkEh8KCWlwYWRkcmVzcxjtp99QIAEoCVIJaXBhZGRyZXNzEm0KHG'
    'luc3VmZmljaWVudGRhdGFoZWFsdGhzdGF0dXMYy7KR6wEgASgOMiUucm91dGU1My5JbnN1ZmZp'
    'Y2llbnREYXRhSGVhbHRoU3RhdHVzUhxpbnN1ZmZpY2llbnRkYXRhaGVhbHRoc3RhdHVzEh0KCG'
    'ludmVydGVkGNnSpxogASgIUghpbnZlcnRlZBIVCgRwb3J0GMf5lBYgASgFUgRwb3J0EjcKB3Jl'
    'Z2lvbnMYq72hESADKA4yGi5yb3V0ZTUzLkhlYWx0aENoZWNrUmVnaW9uUgdyZWdpb25zEkcKDX'
    'Jlc2V0ZWxlbWVudHMY4t3xByADKA4yHi5yb3V0ZTUzLlJlc2V0dGFibGVFbGVtZW50TmFtZVIN'
    'cmVzZXRlbGVtZW50cxIlCgxyZXNvdXJjZXBhdGgYp+WIOCABKAlSDHJlc291cmNlcGF0aBImCg'
    'xzZWFyY2hzdHJpbmcYhZH7lwEgASgJUgxzZWFyY2hzdHJpbmc=');

@$core.Deprecated('Use updateHealthCheckResponseDescriptor instead')
const UpdateHealthCheckResponse$json = {
  '1': 'UpdateHealthCheckResponse',
  '2': [
    {
      '1': 'healthcheck',
      '3': 377540610,
      '4': 1,
      '5': 11,
      '6': '.route53.HealthCheck',
      '10': 'healthcheck'
    },
  ],
};

/// Descriptor for `UpdateHealthCheckResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateHealthCheckResponseDescriptor =
    $convert.base64Decode(
        'ChlVcGRhdGVIZWFsdGhDaGVja1Jlc3BvbnNlEjoKC2hlYWx0aGNoZWNrGIKgg7QBIAEoCzIULn'
        'JvdXRlNTMuSGVhbHRoQ2hlY2tSC2hlYWx0aGNoZWNr');

@$core.Deprecated('Use updateHostedZoneCommentRequestDescriptor instead')
const UpdateHostedZoneCommentRequest$json = {
  '1': 'UpdateHostedZoneCommentRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `UpdateHostedZoneCommentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateHostedZoneCommentRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVIb3N0ZWRab25lQ29tbWVudFJlcXVlc3QSHAoHY29tbWVudBj/v77CASABKAlSB2'
        'NvbW1lbnQSEgoCaWQYgfKitwEgASgJUgJpZA==');

@$core.Deprecated('Use updateHostedZoneCommentResponseDescriptor instead')
const UpdateHostedZoneCommentResponse$json = {
  '1': 'UpdateHostedZoneCommentResponse',
  '2': [
    {
      '1': 'hostedzone',
      '3': 465689249,
      '4': 1,
      '5': 11,
      '6': '.route53.HostedZone',
      '10': 'hostedzone'
    },
  ],
};

/// Descriptor for `UpdateHostedZoneCommentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateHostedZoneCommentResponseDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVIb3N0ZWRab25lQ29tbWVudFJlc3BvbnNlEjcKCmhvc3RlZHpvbmUYobWH3gEgAS'
        'gLMhMucm91dGU1My5Ib3N0ZWRab25lUgpob3N0ZWR6b25l');

@$core.Deprecated('Use updateHostedZoneFeaturesRequestDescriptor instead')
const UpdateHostedZoneFeaturesRequest$json = {
  '1': 'UpdateHostedZoneFeaturesRequest',
  '2': [
    {
      '1': 'enableacceleratedrecovery',
      '3': 482177307,
      '4': 1,
      '5': 8,
      '10': 'enableacceleratedrecovery'
    },
    {'1': 'hostedzoneid', '3': 346531710, '4': 1, '5': 9, '10': 'hostedzoneid'},
  ],
};

/// Descriptor for `UpdateHostedZoneFeaturesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateHostedZoneFeaturesRequestDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVIb3N0ZWRab25lRmVhdHVyZXNSZXF1ZXN0EkAKGWVuYWJsZWFjY2VsZXJhdGVkcm'
        'Vjb3ZlcnkYm+L15QEgASgIUhllbmFibGVhY2NlbGVyYXRlZHJlY292ZXJ5EiYKDGhvc3RlZHpv'
        'bmVpZBj+zp6lASABKAlSDGhvc3RlZHpvbmVpZA==');

@$core.Deprecated('Use updateHostedZoneFeaturesResponseDescriptor instead')
const UpdateHostedZoneFeaturesResponse$json = {
  '1': 'UpdateHostedZoneFeaturesResponse',
};

/// Descriptor for `UpdateHostedZoneFeaturesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateHostedZoneFeaturesResponseDescriptor =
    $convert.base64Decode('CiBVcGRhdGVIb3N0ZWRab25lRmVhdHVyZXNSZXNwb25zZQ==');

@$core.Deprecated('Use updateTrafficPolicyCommentRequestDescriptor instead')
const UpdateTrafficPolicyCommentRequest$json = {
  '1': 'UpdateTrafficPolicyCommentRequest',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 5, '10': 'version'},
  ],
};

/// Descriptor for `UpdateTrafficPolicyCommentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrafficPolicyCommentRequestDescriptor =
    $convert.base64Decode(
        'CiFVcGRhdGVUcmFmZmljUG9saWN5Q29tbWVudFJlcXVlc3QSHAoHY29tbWVudBj/v77CASABKA'
        'lSB2NvbW1lbnQSEgoCaWQYgfKitwEgASgJUgJpZBIcCgd2ZXJzaW9uGLiqt+4BIAEoBVIHdmVy'
        'c2lvbg==');

@$core.Deprecated('Use updateTrafficPolicyCommentResponseDescriptor instead')
const UpdateTrafficPolicyCommentResponse$json = {
  '1': 'UpdateTrafficPolicyCommentResponse',
  '2': [
    {
      '1': 'trafficpolicy',
      '3': 154595657,
      '4': 1,
      '5': 11,
      '6': '.route53.TrafficPolicy',
      '10': 'trafficpolicy'
    },
  ],
};

/// Descriptor for `UpdateTrafficPolicyCommentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrafficPolicyCommentResponseDescriptor =
    $convert.base64Decode(
        'CiJVcGRhdGVUcmFmZmljUG9saWN5Q29tbWVudFJlc3BvbnNlEj8KDXRyYWZmaWNwb2xpY3kYye'
        'LbSSABKAsyFi5yb3V0ZTUzLlRyYWZmaWNQb2xpY3lSDXRyYWZmaWNwb2xpY3k=');

@$core.Deprecated('Use updateTrafficPolicyInstanceRequestDescriptor instead')
const UpdateTrafficPolicyInstanceRequest$json = {
  '1': 'UpdateTrafficPolicyInstanceRequest',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'ttl', '3': 526904700, '4': 1, '5': 3, '10': 'ttl'},
    {
      '1': 'trafficpolicyid',
      '3': 40235222,
      '4': 1,
      '5': 9,
      '10': 'trafficpolicyid'
    },
    {
      '1': 'trafficpolicyversion',
      '3': 479078485,
      '4': 1,
      '5': 5,
      '10': 'trafficpolicyversion'
    },
  ],
};

/// Descriptor for `UpdateTrafficPolicyInstanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrafficPolicyInstanceRequestDescriptor =
    $convert.base64Decode(
        'CiJVcGRhdGVUcmFmZmljUG9saWN5SW5zdGFuY2VSZXF1ZXN0EhIKAmlkGIHyorcBIAEoCVICaW'
        'QSFAoDdHRsGPzan/sBIAEoA1IDdHRsEisKD3RyYWZmaWNwb2xpY3lpZBjW4ZcTIAEoCVIPdHJh'
        'ZmZpY3BvbGljeWlkEjYKFHRyYWZmaWNwb2xpY3l2ZXJzaW9uGNXQuOQBIAEoBVIUdHJhZmZpY3'
        'BvbGljeXZlcnNpb24=');

@$core.Deprecated('Use updateTrafficPolicyInstanceResponseDescriptor instead')
const UpdateTrafficPolicyInstanceResponse$json = {
  '1': 'UpdateTrafficPolicyInstanceResponse',
  '2': [
    {
      '1': 'trafficpolicyinstance',
      '3': 205651476,
      '4': 1,
      '5': 11,
      '6': '.route53.TrafficPolicyInstance',
      '10': 'trafficpolicyinstance'
    },
  ],
};

/// Descriptor for `UpdateTrafficPolicyInstanceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrafficPolicyInstanceResponseDescriptor =
    $convert.base64Decode(
        'CiNVcGRhdGVUcmFmZmljUG9saWN5SW5zdGFuY2VSZXNwb25zZRJXChV0cmFmZmljcG9saWN5aW'
        '5zdGFuY2UYlPyHYiABKAsyHi5yb3V0ZTUzLlRyYWZmaWNQb2xpY3lJbnN0YW5jZVIVdHJhZmZp'
        'Y3BvbGljeWluc3RhbmNl');

@$core.Deprecated('Use vPCDescriptor instead')
const VPC$json = {
  '1': 'VPC',
  '2': [
    {'1': 'vpcid', '3': 325135798, '4': 1, '5': 9, '10': 'vpcid'},
    {
      '1': 'vpcregion',
      '3': 474180765,
      '4': 1,
      '5': 14,
      '6': '.route53.VPCRegion',
      '10': 'vpcregion'
    },
  ],
};

/// Descriptor for `VPC`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vPCDescriptor = $convert.base64Decode(
    'CgNWUEMSGAoFdnBjaWQYttuEmwEgASgJUgV2cGNpZBI0Cgl2cGNyZWdpb24YndmN4gEgASgOMh'
    'Iucm91dGU1My5WUENSZWdpb25SCXZwY3JlZ2lvbg==');

@$core.Deprecated('Use vPCAssociationAuthorizationNotFoundDescriptor instead')
const VPCAssociationAuthorizationNotFound$json = {
  '1': 'VPCAssociationAuthorizationNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `VPCAssociationAuthorizationNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vPCAssociationAuthorizationNotFoundDescriptor =
    $convert.base64Decode(
        'CiNWUENBc3NvY2lhdGlvbkF1dGhvcml6YXRpb25Ob3RGb3VuZBIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use vPCAssociationNotFoundDescriptor instead')
const VPCAssociationNotFound$json = {
  '1': 'VPCAssociationNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `VPCAssociationNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vPCAssociationNotFoundDescriptor =
    $convert.base64Decode(
        'ChZWUENBc3NvY2lhdGlvbk5vdEZvdW5kEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');
