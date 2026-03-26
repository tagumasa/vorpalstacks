// This is a generated file - do not edit.
//
// Generated from cognito_idp.proto.

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

@$core.Deprecated('Use accountTakeoverEventActionTypeDescriptor instead')
const AccountTakeoverEventActionType$json = {
  '1': 'AccountTakeoverEventActionType',
  '2': [
    {'1': 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_BLOCK', '2': 0},
    {'1': 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_REQUIRED', '2': 1},
    {'1': 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_NO_ACTION', '2': 2},
    {'1': 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_IF_CONFIGURED', '2': 3},
  ],
};

/// Descriptor for `AccountTakeoverEventActionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List accountTakeoverEventActionTypeDescriptor = $convert.base64Decode(
    'Ch5BY2NvdW50VGFrZW92ZXJFdmVudEFjdGlvblR5cGUSLAooQUNDT1VOVF9UQUtFT1ZFUl9FVk'
    'VOVF9BQ1RJT05fVFlQRV9CTE9DSxAAEjMKL0FDQ09VTlRfVEFLRU9WRVJfRVZFTlRfQUNUSU9O'
    'X1RZUEVfTUZBX1JFUVVJUkVEEAESMAosQUNDT1VOVF9UQUtFT1ZFUl9FVkVOVF9BQ1RJT05fVF'
    'lQRV9OT19BQ1RJT04QAhI4CjRBQ0NPVU5UX1RBS0VPVkVSX0VWRU5UX0FDVElPTl9UWVBFX01G'
    'QV9JRl9DT05GSUdVUkVEEAM=');

@$core.Deprecated('Use advancedSecurityEnabledModeTypeDescriptor instead')
const AdvancedSecurityEnabledModeType$json = {
  '1': 'AdvancedSecurityEnabledModeType',
  '2': [
    {'1': 'ADVANCED_SECURITY_ENABLED_MODE_TYPE_AUDIT', '2': 0},
    {'1': 'ADVANCED_SECURITY_ENABLED_MODE_TYPE_ENFORCED', '2': 1},
  ],
};

/// Descriptor for `AdvancedSecurityEnabledModeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List advancedSecurityEnabledModeTypeDescriptor =
    $convert.base64Decode(
        'Ch9BZHZhbmNlZFNlY3VyaXR5RW5hYmxlZE1vZGVUeXBlEi0KKUFEVkFOQ0VEX1NFQ1VSSVRZX0'
        'VOQUJMRURfTU9ERV9UWVBFX0FVRElUEAASMAosQURWQU5DRURfU0VDVVJJVFlfRU5BQkxFRF9N'
        'T0RFX1RZUEVfRU5GT1JDRUQQAQ==');

@$core.Deprecated('Use advancedSecurityModeTypeDescriptor instead')
const AdvancedSecurityModeType$json = {
  '1': 'AdvancedSecurityModeType',
  '2': [
    {'1': 'ADVANCED_SECURITY_MODE_TYPE_AUDIT', '2': 0},
    {'1': 'ADVANCED_SECURITY_MODE_TYPE_ENFORCED', '2': 1},
    {'1': 'ADVANCED_SECURITY_MODE_TYPE_OFF', '2': 2},
  ],
};

/// Descriptor for `AdvancedSecurityModeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List advancedSecurityModeTypeDescriptor = $convert.base64Decode(
    'ChhBZHZhbmNlZFNlY3VyaXR5TW9kZVR5cGUSJQohQURWQU5DRURfU0VDVVJJVFlfTU9ERV9UWV'
    'BFX0FVRElUEAASKAokQURWQU5DRURfU0VDVVJJVFlfTU9ERV9UWVBFX0VORk9SQ0VEEAESIwof'
    'QURWQU5DRURfU0VDVVJJVFlfTU9ERV9UWVBFX09GRhAC');

@$core.Deprecated('Use aliasAttributeTypeDescriptor instead')
const AliasAttributeType$json = {
  '1': 'AliasAttributeType',
  '2': [
    {'1': 'ALIAS_ATTRIBUTE_TYPE_PHONE_NUMBER', '2': 0},
    {'1': 'ALIAS_ATTRIBUTE_TYPE_EMAIL', '2': 1},
    {'1': 'ALIAS_ATTRIBUTE_TYPE_PREFERRED_USERNAME', '2': 2},
  ],
};

/// Descriptor for `AliasAttributeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List aliasAttributeTypeDescriptor = $convert.base64Decode(
    'ChJBbGlhc0F0dHJpYnV0ZVR5cGUSJQohQUxJQVNfQVRUUklCVVRFX1RZUEVfUEhPTkVfTlVNQk'
    'VSEAASHgoaQUxJQVNfQVRUUklCVVRFX1RZUEVfRU1BSUwQARIrCidBTElBU19BVFRSSUJVVEVf'
    'VFlQRV9QUkVGRVJSRURfVVNFUk5BTUUQAg==');

@$core.Deprecated('Use assetCategoryTypeDescriptor instead')
const AssetCategoryType$json = {
  '1': 'AssetCategoryType',
  '2': [
    {'1': 'ASSET_CATEGORY_TYPE_FORM_LOGO', '2': 0},
    {'1': 'ASSET_CATEGORY_TYPE_FAVICON_ICO', '2': 1},
    {'1': 'ASSET_CATEGORY_TYPE_PASSWORD_GRAPHIC', '2': 2},
    {'1': 'ASSET_CATEGORY_TYPE_PAGE_BACKGROUND', '2': 3},
    {'1': 'ASSET_CATEGORY_TYPE_PAGE_HEADER_BACKGROUND', '2': 4},
    {'1': 'ASSET_CATEGORY_TYPE_SMS_GRAPHIC', '2': 5},
    {'1': 'ASSET_CATEGORY_TYPE_PASSKEY_GRAPHIC', '2': 6},
    {'1': 'ASSET_CATEGORY_TYPE_PAGE_FOOTER_BACKGROUND', '2': 7},
    {'1': 'ASSET_CATEGORY_TYPE_EMAIL_GRAPHIC', '2': 8},
    {'1': 'ASSET_CATEGORY_TYPE_PAGE_HEADER_LOGO', '2': 9},
    {'1': 'ASSET_CATEGORY_TYPE_FORM_BACKGROUND', '2': 10},
    {'1': 'ASSET_CATEGORY_TYPE_AUTH_APP_GRAPHIC', '2': 11},
    {'1': 'ASSET_CATEGORY_TYPE_IDP_BUTTON_ICON', '2': 12},
    {'1': 'ASSET_CATEGORY_TYPE_FAVICON_SVG', '2': 13},
    {'1': 'ASSET_CATEGORY_TYPE_PAGE_FOOTER_LOGO', '2': 14},
  ],
};

/// Descriptor for `AssetCategoryType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List assetCategoryTypeDescriptor = $convert.base64Decode(
    'ChFBc3NldENhdGVnb3J5VHlwZRIhCh1BU1NFVF9DQVRFR09SWV9UWVBFX0ZPUk1fTE9HTxAAEi'
    'MKH0FTU0VUX0NBVEVHT1JZX1RZUEVfRkFWSUNPTl9JQ08QARIoCiRBU1NFVF9DQVRFR09SWV9U'
    'WVBFX1BBU1NXT1JEX0dSQVBISUMQAhInCiNBU1NFVF9DQVRFR09SWV9UWVBFX1BBR0VfQkFDS0'
    'dST1VORBADEi4KKkFTU0VUX0NBVEVHT1JZX1RZUEVfUEFHRV9IRUFERVJfQkFDS0dST1VORBAE'
    'EiMKH0FTU0VUX0NBVEVHT1JZX1RZUEVfU01TX0dSQVBISUMQBRInCiNBU1NFVF9DQVRFR09SWV'
    '9UWVBFX1BBU1NLRVlfR1JBUEhJQxAGEi4KKkFTU0VUX0NBVEVHT1JZX1RZUEVfUEFHRV9GT09U'
    'RVJfQkFDS0dST1VORBAHEiUKIUFTU0VUX0NBVEVHT1JZX1RZUEVfRU1BSUxfR1JBUEhJQxAIEi'
    'gKJEFTU0VUX0NBVEVHT1JZX1RZUEVfUEFHRV9IRUFERVJfTE9HTxAJEicKI0FTU0VUX0NBVEVH'
    'T1JZX1RZUEVfRk9STV9CQUNLR1JPVU5EEAoSKAokQVNTRVRfQ0FURUdPUllfVFlQRV9BVVRIX0'
    'FQUF9HUkFQSElDEAsSJwojQVNTRVRfQ0FURUdPUllfVFlQRV9JRFBfQlVUVE9OX0lDT04QDBIj'
    'Ch9BU1NFVF9DQVRFR09SWV9UWVBFX0ZBVklDT05fU1ZHEA0SKAokQVNTRVRfQ0FURUdPUllfVF'
    'lQRV9QQUdFX0ZPT1RFUl9MT0dPEA4=');

@$core.Deprecated('Use assetExtensionTypeDescriptor instead')
const AssetExtensionType$json = {
  '1': 'AssetExtensionType',
  '2': [
    {'1': 'ASSET_EXTENSION_TYPE_ICO', '2': 0},
    {'1': 'ASSET_EXTENSION_TYPE_SVG', '2': 1},
    {'1': 'ASSET_EXTENSION_TYPE_JPEG', '2': 2},
    {'1': 'ASSET_EXTENSION_TYPE_PNG', '2': 3},
    {'1': 'ASSET_EXTENSION_TYPE_WEBP', '2': 4},
  ],
};

/// Descriptor for `AssetExtensionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List assetExtensionTypeDescriptor = $convert.base64Decode(
    'ChJBc3NldEV4dGVuc2lvblR5cGUSHAoYQVNTRVRfRVhURU5TSU9OX1RZUEVfSUNPEAASHAoYQV'
    'NTRVRfRVhURU5TSU9OX1RZUEVfU1ZHEAESHQoZQVNTRVRfRVhURU5TSU9OX1RZUEVfSlBFRxAC'
    'EhwKGEFTU0VUX0VYVEVOU0lPTl9UWVBFX1BORxADEh0KGUFTU0VUX0VYVEVOU0lPTl9UWVBFX1'
    'dFQlAQBA==');

@$core.Deprecated('Use attributeDataTypeDescriptor instead')
const AttributeDataType$json = {
  '1': 'AttributeDataType',
  '2': [
    {'1': 'ATTRIBUTE_DATA_TYPE_DATETIME', '2': 0},
    {'1': 'ATTRIBUTE_DATA_TYPE_STRING', '2': 1},
    {'1': 'ATTRIBUTE_DATA_TYPE_NUMBER', '2': 2},
    {'1': 'ATTRIBUTE_DATA_TYPE_BOOLEAN', '2': 3},
  ],
};

/// Descriptor for `AttributeDataType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List attributeDataTypeDescriptor = $convert.base64Decode(
    'ChFBdHRyaWJ1dGVEYXRhVHlwZRIgChxBVFRSSUJVVEVfREFUQV9UWVBFX0RBVEVUSU1FEAASHg'
    'oaQVRUUklCVVRFX0RBVEFfVFlQRV9TVFJJTkcQARIeChpBVFRSSUJVVEVfREFUQV9UWVBFX05V'
    'TUJFUhACEh8KG0FUVFJJQlVURV9EQVRBX1RZUEVfQk9PTEVBThAD');

@$core.Deprecated('Use authFactorTypeDescriptor instead')
const AuthFactorType$json = {
  '1': 'AuthFactorType',
  '2': [
    {'1': 'AUTH_FACTOR_TYPE_EMAIL_OTP', '2': 0},
    {'1': 'AUTH_FACTOR_TYPE_SMS_OTP', '2': 1},
    {'1': 'AUTH_FACTOR_TYPE_WEB_AUTHN', '2': 2},
    {'1': 'AUTH_FACTOR_TYPE_PASSWORD', '2': 3},
  ],
};

/// Descriptor for `AuthFactorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List authFactorTypeDescriptor = $convert.base64Decode(
    'Cg5BdXRoRmFjdG9yVHlwZRIeChpBVVRIX0ZBQ1RPUl9UWVBFX0VNQUlMX09UUBAAEhwKGEFVVE'
    'hfRkFDVE9SX1RZUEVfU01TX09UUBABEh4KGkFVVEhfRkFDVE9SX1RZUEVfV0VCX0FVVEhOEAIS'
    'HQoZQVVUSF9GQUNUT1JfVFlQRV9QQVNTV09SRBAD');

@$core.Deprecated('Use authFlowTypeDescriptor instead')
const AuthFlowType$json = {
  '1': 'AuthFlowType',
  '2': [
    {'1': 'AUTH_FLOW_TYPE_ADMIN_NO_SRP_AUTH', '2': 0},
    {'1': 'AUTH_FLOW_TYPE_ADMIN_USER_PASSWORD_AUTH', '2': 1},
    {'1': 'AUTH_FLOW_TYPE_USER_AUTH', '2': 2},
    {'1': 'AUTH_FLOW_TYPE_USER_SRP_AUTH', '2': 3},
    {'1': 'AUTH_FLOW_TYPE_USER_PASSWORD_AUTH', '2': 4},
    {'1': 'AUTH_FLOW_TYPE_CUSTOM_AUTH', '2': 5},
    {'1': 'AUTH_FLOW_TYPE_REFRESH_TOKEN_AUTH', '2': 6},
    {'1': 'AUTH_FLOW_TYPE_REFRESH_TOKEN', '2': 7},
  ],
};

/// Descriptor for `AuthFlowType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List authFlowTypeDescriptor = $convert.base64Decode(
    'CgxBdXRoRmxvd1R5cGUSJAogQVVUSF9GTE9XX1RZUEVfQURNSU5fTk9fU1JQX0FVVEgQABIrCi'
    'dBVVRIX0ZMT1dfVFlQRV9BRE1JTl9VU0VSX1BBU1NXT1JEX0FVVEgQARIcChhBVVRIX0ZMT1df'
    'VFlQRV9VU0VSX0FVVEgQAhIgChxBVVRIX0ZMT1dfVFlQRV9VU0VSX1NSUF9BVVRIEAMSJQohQV'
    'VUSF9GTE9XX1RZUEVfVVNFUl9QQVNTV09SRF9BVVRIEAQSHgoaQVVUSF9GTE9XX1RZUEVfQ1VT'
    'VE9NX0FVVEgQBRIlCiFBVVRIX0ZMT1dfVFlQRV9SRUZSRVNIX1RPS0VOX0FVVEgQBhIgChxBVV'
    'RIX0ZMT1dfVFlQRV9SRUZSRVNIX1RPS0VOEAc=');

@$core.Deprecated('Use challengeNameDescriptor instead')
const ChallengeName$json = {
  '1': 'ChallengeName',
  '2': [
    {'1': 'CHALLENGE_NAME_PASSWORD', '2': 0},
    {'1': 'CHALLENGE_NAME_MFA', '2': 1},
  ],
};

/// Descriptor for `ChallengeName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List challengeNameDescriptor = $convert.base64Decode(
    'Cg1DaGFsbGVuZ2VOYW1lEhsKF0NIQUxMRU5HRV9OQU1FX1BBU1NXT1JEEAASFgoSQ0hBTExFTk'
    'dFX05BTUVfTUZBEAE=');

@$core.Deprecated('Use challengeNameTypeDescriptor instead')
const ChallengeNameType$json = {
  '1': 'ChallengeNameType',
  '2': [
    {'1': 'CHALLENGE_NAME_TYPE_EMAIL_OTP', '2': 0},
    {'1': 'CHALLENGE_NAME_TYPE_DEVICE_PASSWORD_VERIFIER', '2': 1},
    {'1': 'CHALLENGE_NAME_TYPE_ADMIN_NO_SRP_AUTH', '2': 2},
    {'1': 'CHALLENGE_NAME_TYPE_DEVICE_SRP_AUTH', '2': 3},
    {'1': 'CHALLENGE_NAME_TYPE_SMS_OTP', '2': 4},
    {'1': 'CHALLENGE_NAME_TYPE_PASSWORD_SRP', '2': 5},
    {'1': 'CHALLENGE_NAME_TYPE_SELECT_CHALLENGE', '2': 6},
    {'1': 'CHALLENGE_NAME_TYPE_PASSWORD_VERIFIER', '2': 7},
    {'1': 'CHALLENGE_NAME_TYPE_NEW_PASSWORD_REQUIRED', '2': 8},
    {'1': 'CHALLENGE_NAME_TYPE_WEB_AUTHN', '2': 9},
    {'1': 'CHALLENGE_NAME_TYPE_CUSTOM_CHALLENGE', '2': 10},
    {'1': 'CHALLENGE_NAME_TYPE_SOFTWARE_TOKEN_MFA', '2': 11},
    {'1': 'CHALLENGE_NAME_TYPE_SELECT_MFA_TYPE', '2': 12},
    {'1': 'CHALLENGE_NAME_TYPE_MFA_SETUP', '2': 13},
    {'1': 'CHALLENGE_NAME_TYPE_PASSWORD', '2': 14},
    {'1': 'CHALLENGE_NAME_TYPE_SMS_MFA', '2': 15},
  ],
};

/// Descriptor for `ChallengeNameType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List challengeNameTypeDescriptor = $convert.base64Decode(
    'ChFDaGFsbGVuZ2VOYW1lVHlwZRIhCh1DSEFMTEVOR0VfTkFNRV9UWVBFX0VNQUlMX09UUBAAEj'
    'AKLENIQUxMRU5HRV9OQU1FX1RZUEVfREVWSUNFX1BBU1NXT1JEX1ZFUklGSUVSEAESKQolQ0hB'
    'TExFTkdFX05BTUVfVFlQRV9BRE1JTl9OT19TUlBfQVVUSBACEicKI0NIQUxMRU5HRV9OQU1FX1'
    'RZUEVfREVWSUNFX1NSUF9BVVRIEAMSHwobQ0hBTExFTkdFX05BTUVfVFlQRV9TTVNfT1RQEAQS'
    'JAogQ0hBTExFTkdFX05BTUVfVFlQRV9QQVNTV09SRF9TUlAQBRIoCiRDSEFMTEVOR0VfTkFNRV'
    '9UWVBFX1NFTEVDVF9DSEFMTEVOR0UQBhIpCiVDSEFMTEVOR0VfTkFNRV9UWVBFX1BBU1NXT1JE'
    'X1ZFUklGSUVSEAcSLQopQ0hBTExFTkdFX05BTUVfVFlQRV9ORVdfUEFTU1dPUkRfUkVRVUlSRU'
    'QQCBIhCh1DSEFMTEVOR0VfTkFNRV9UWVBFX1dFQl9BVVRIThAJEigKJENIQUxMRU5HRV9OQU1F'
    'X1RZUEVfQ1VTVE9NX0NIQUxMRU5HRRAKEioKJkNIQUxMRU5HRV9OQU1FX1RZUEVfU09GVFdBUk'
    'VfVE9LRU5fTUZBEAsSJwojQ0hBTExFTkdFX05BTUVfVFlQRV9TRUxFQ1RfTUZBX1RZUEUQDBIh'
    'Ch1DSEFMTEVOR0VfTkFNRV9UWVBFX01GQV9TRVRVUBANEiAKHENIQUxMRU5HRV9OQU1FX1RZUE'
    'VfUEFTU1dPUkQQDhIfChtDSEFMTEVOR0VfTkFNRV9UWVBFX1NNU19NRkEQDw==');

@$core.Deprecated('Use challengeResponseDescriptor instead')
const ChallengeResponse$json = {
  '1': 'ChallengeResponse',
  '2': [
    {'1': 'CHALLENGE_RESPONSE_FAILURE', '2': 0},
    {'1': 'CHALLENGE_RESPONSE_SUCCESS', '2': 1},
  ],
};

/// Descriptor for `ChallengeResponse`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List challengeResponseDescriptor = $convert.base64Decode(
    'ChFDaGFsbGVuZ2VSZXNwb25zZRIeChpDSEFMTEVOR0VfUkVTUE9OU0VfRkFJTFVSRRAAEh4KGk'
    'NIQUxMRU5HRV9SRVNQT05TRV9TVUNDRVNTEAE=');

@$core.Deprecated('Use colorSchemeModeTypeDescriptor instead')
const ColorSchemeModeType$json = {
  '1': 'ColorSchemeModeType',
  '2': [
    {'1': 'COLOR_SCHEME_MODE_TYPE_LIGHT', '2': 0},
    {'1': 'COLOR_SCHEME_MODE_TYPE_DYNAMIC', '2': 1},
    {'1': 'COLOR_SCHEME_MODE_TYPE_DARK', '2': 2},
  ],
};

/// Descriptor for `ColorSchemeModeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List colorSchemeModeTypeDescriptor = $convert.base64Decode(
    'ChNDb2xvclNjaGVtZU1vZGVUeXBlEiAKHENPTE9SX1NDSEVNRV9NT0RFX1RZUEVfTElHSFQQAB'
    'IiCh5DT0xPUl9TQ0hFTUVfTU9ERV9UWVBFX0RZTkFNSUMQARIfChtDT0xPUl9TQ0hFTUVfTU9E'
    'RV9UWVBFX0RBUksQAg==');

@$core.Deprecated('Use compromisedCredentialsEventActionTypeDescriptor instead')
const CompromisedCredentialsEventActionType$json = {
  '1': 'CompromisedCredentialsEventActionType',
  '2': [
    {'1': 'COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_BLOCK', '2': 0},
    {'1': 'COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_NO_ACTION', '2': 1},
  ],
};

/// Descriptor for `CompromisedCredentialsEventActionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List compromisedCredentialsEventActionTypeDescriptor =
    $convert.base64Decode(
        'CiVDb21wcm9taXNlZENyZWRlbnRpYWxzRXZlbnRBY3Rpb25UeXBlEjMKL0NPTVBST01JU0VEX0'
        'NSRURFTlRJQUxTX0VWRU5UX0FDVElPTl9UWVBFX0JMT0NLEAASNwozQ09NUFJPTUlTRURfQ1JF'
        'REVOVElBTFNfRVZFTlRfQUNUSU9OX1RZUEVfTk9fQUNUSU9OEAE=');

@$core.Deprecated('Use customEmailSenderLambdaVersionTypeDescriptor instead')
const CustomEmailSenderLambdaVersionType$json = {
  '1': 'CustomEmailSenderLambdaVersionType',
  '2': [
    {'1': 'CUSTOM_EMAIL_SENDER_LAMBDA_VERSION_TYPE_V1_0', '2': 0},
  ],
};

/// Descriptor for `CustomEmailSenderLambdaVersionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List customEmailSenderLambdaVersionTypeDescriptor =
    $convert.base64Decode(
        'CiJDdXN0b21FbWFpbFNlbmRlckxhbWJkYVZlcnNpb25UeXBlEjAKLENVU1RPTV9FTUFJTF9TRU'
        '5ERVJfTEFNQkRBX1ZFUlNJT05fVFlQRV9WMV8wEAA=');

@$core.Deprecated('Use customSMSSenderLambdaVersionTypeDescriptor instead')
const CustomSMSSenderLambdaVersionType$json = {
  '1': 'CustomSMSSenderLambdaVersionType',
  '2': [
    {'1': 'CUSTOM_S_M_S_SENDER_LAMBDA_VERSION_TYPE_V1_0', '2': 0},
  ],
};

/// Descriptor for `CustomSMSSenderLambdaVersionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List customSMSSenderLambdaVersionTypeDescriptor =
    $convert.base64Decode(
        'CiBDdXN0b21TTVNTZW5kZXJMYW1iZGFWZXJzaW9uVHlwZRIwCixDVVNUT01fU19NX1NfU0VORE'
        'VSX0xBTUJEQV9WRVJTSU9OX1RZUEVfVjFfMBAA');

@$core.Deprecated('Use defaultEmailOptionTypeDescriptor instead')
const DefaultEmailOptionType$json = {
  '1': 'DefaultEmailOptionType',
  '2': [
    {'1': 'DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_LINK', '2': 0},
    {'1': 'DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_CODE', '2': 1},
  ],
};

/// Descriptor for `DefaultEmailOptionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List defaultEmailOptionTypeDescriptor = $convert.base64Decode(
    'ChZEZWZhdWx0RW1haWxPcHRpb25UeXBlEi8KK0RFRkFVTFRfRU1BSUxfT1BUSU9OX1RZUEVfQ0'
    '9ORklSTV9XSVRIX0xJTksQABIvCitERUZBVUxUX0VNQUlMX09QVElPTl9UWVBFX0NPTkZJUk1f'
    'V0lUSF9DT0RFEAE=');

@$core.Deprecated('Use deletionProtectionTypeDescriptor instead')
const DeletionProtectionType$json = {
  '1': 'DeletionProtectionType',
  '2': [
    {'1': 'DELETION_PROTECTION_TYPE_ACTIVE', '2': 0},
    {'1': 'DELETION_PROTECTION_TYPE_INACTIVE', '2': 1},
  ],
};

/// Descriptor for `DeletionProtectionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deletionProtectionTypeDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGlvblByb3RlY3Rpb25UeXBlEiMKH0RFTEVUSU9OX1BST1RFQ1RJT05fVFlQRV9BQ1'
        'RJVkUQABIlCiFERUxFVElPTl9QUk9URUNUSU9OX1RZUEVfSU5BQ1RJVkUQAQ==');

@$core.Deprecated('Use deliveryMediumTypeDescriptor instead')
const DeliveryMediumType$json = {
  '1': 'DeliveryMediumType',
  '2': [
    {'1': 'DELIVERY_MEDIUM_TYPE_EMAIL', '2': 0},
    {'1': 'DELIVERY_MEDIUM_TYPE_SMS', '2': 1},
  ],
};

/// Descriptor for `DeliveryMediumType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deliveryMediumTypeDescriptor = $convert.base64Decode(
    'ChJEZWxpdmVyeU1lZGl1bVR5cGUSHgoaREVMSVZFUllfTUVESVVNX1RZUEVfRU1BSUwQABIcCh'
    'hERUxJVkVSWV9NRURJVU1fVFlQRV9TTVMQAQ==');

@$core.Deprecated('Use deviceRememberedStatusTypeDescriptor instead')
const DeviceRememberedStatusType$json = {
  '1': 'DeviceRememberedStatusType',
  '2': [
    {'1': 'DEVICE_REMEMBERED_STATUS_TYPE_NOT_REMEMBERED', '2': 0},
    {'1': 'DEVICE_REMEMBERED_STATUS_TYPE_REMEMBERED', '2': 1},
  ],
};

/// Descriptor for `DeviceRememberedStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deviceRememberedStatusTypeDescriptor =
    $convert.base64Decode(
        'ChpEZXZpY2VSZW1lbWJlcmVkU3RhdHVzVHlwZRIwCixERVZJQ0VfUkVNRU1CRVJFRF9TVEFUVV'
        'NfVFlQRV9OT1RfUkVNRU1CRVJFRBAAEiwKKERFVklDRV9SRU1FTUJFUkVEX1NUQVRVU19UWVBF'
        'X1JFTUVNQkVSRUQQAQ==');

@$core.Deprecated('Use domainStatusTypeDescriptor instead')
const DomainStatusType$json = {
  '1': 'DomainStatusType',
  '2': [
    {'1': 'DOMAIN_STATUS_TYPE_UPDATING', '2': 0},
    {'1': 'DOMAIN_STATUS_TYPE_ACTIVE', '2': 1},
    {'1': 'DOMAIN_STATUS_TYPE_DELETING', '2': 2},
    {'1': 'DOMAIN_STATUS_TYPE_CREATING', '2': 3},
    {'1': 'DOMAIN_STATUS_TYPE_FAILED', '2': 4},
  ],
};

/// Descriptor for `DomainStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List domainStatusTypeDescriptor = $convert.base64Decode(
    'ChBEb21haW5TdGF0dXNUeXBlEh8KG0RPTUFJTl9TVEFUVVNfVFlQRV9VUERBVElORxAAEh0KGU'
    'RPTUFJTl9TVEFUVVNfVFlQRV9BQ1RJVkUQARIfChtET01BSU5fU1RBVFVTX1RZUEVfREVMRVRJ'
    'TkcQAhIfChtET01BSU5fU1RBVFVTX1RZUEVfQ1JFQVRJTkcQAxIdChlET01BSU5fU1RBVFVTX1'
    'RZUEVfRkFJTEVEEAQ=');

@$core.Deprecated('Use emailSendingAccountTypeDescriptor instead')
const EmailSendingAccountType$json = {
  '1': 'EmailSendingAccountType',
  '2': [
    {'1': 'EMAIL_SENDING_ACCOUNT_TYPE_COGNITO_DEFAULT', '2': 0},
    {'1': 'EMAIL_SENDING_ACCOUNT_TYPE_DEVELOPER', '2': 1},
  ],
};

/// Descriptor for `EmailSendingAccountType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List emailSendingAccountTypeDescriptor = $convert.base64Decode(
    'ChdFbWFpbFNlbmRpbmdBY2NvdW50VHlwZRIuCipFTUFJTF9TRU5ESU5HX0FDQ09VTlRfVFlQRV'
    '9DT0dOSVRPX0RFRkFVTFQQABIoCiRFTUFJTF9TRU5ESU5HX0FDQ09VTlRfVFlQRV9ERVZFTE9Q'
    'RVIQAQ==');

@$core.Deprecated('Use eventFilterTypeDescriptor instead')
const EventFilterType$json = {
  '1': 'EventFilterType',
  '2': [
    {'1': 'EVENT_FILTER_TYPE_PASSWORD_CHANGE', '2': 0},
    {'1': 'EVENT_FILTER_TYPE_SIGN_IN', '2': 1},
    {'1': 'EVENT_FILTER_TYPE_SIGN_UP', '2': 2},
  ],
};

/// Descriptor for `EventFilterType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventFilterTypeDescriptor = $convert.base64Decode(
    'Cg9FdmVudEZpbHRlclR5cGUSJQohRVZFTlRfRklMVEVSX1RZUEVfUEFTU1dPUkRfQ0hBTkdFEA'
    'ASHQoZRVZFTlRfRklMVEVSX1RZUEVfU0lHTl9JThABEh0KGUVWRU5UX0ZJTFRFUl9UWVBFX1NJ'
    'R05fVVAQAg==');

@$core.Deprecated('Use eventResponseTypeDescriptor instead')
const EventResponseType$json = {
  '1': 'EventResponseType',
  '2': [
    {'1': 'EVENT_RESPONSE_TYPE_FAIL', '2': 0},
    {'1': 'EVENT_RESPONSE_TYPE_PASS', '2': 1},
    {'1': 'EVENT_RESPONSE_TYPE_INPROGRESS', '2': 2},
  ],
};

/// Descriptor for `EventResponseType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventResponseTypeDescriptor = $convert.base64Decode(
    'ChFFdmVudFJlc3BvbnNlVHlwZRIcChhFVkVOVF9SRVNQT05TRV9UWVBFX0ZBSUwQABIcChhFVk'
    'VOVF9SRVNQT05TRV9UWVBFX1BBU1MQARIiCh5FVkVOVF9SRVNQT05TRV9UWVBFX0lOUFJPR1JF'
    'U1MQAg==');

@$core.Deprecated('Use eventSourceNameDescriptor instead')
const EventSourceName$json = {
  '1': 'EventSourceName',
  '2': [
    {'1': 'EVENT_SOURCE_NAME_USER_AUTH_EVENTS', '2': 0},
    {'1': 'EVENT_SOURCE_NAME_USER_NOTIFICATION', '2': 1},
  ],
};

/// Descriptor for `EventSourceName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventSourceNameDescriptor = $convert.base64Decode(
    'Cg9FdmVudFNvdXJjZU5hbWUSJgoiRVZFTlRfU09VUkNFX05BTUVfVVNFUl9BVVRIX0VWRU5UUx'
    'AAEicKI0VWRU5UX1NPVVJDRV9OQU1FX1VTRVJfTk9USUZJQ0FUSU9OEAE=');

@$core.Deprecated('Use eventTypeDescriptor instead')
const EventType$json = {
  '1': 'EventType',
  '2': [
    {'1': 'EVENT_TYPE_SIGNIN', '2': 0},
    {'1': 'EVENT_TYPE_RESENDCODE', '2': 1},
    {'1': 'EVENT_TYPE_PASSWORDCHANGE', '2': 2},
    {'1': 'EVENT_TYPE_SIGNUP', '2': 3},
    {'1': 'EVENT_TYPE_FORGOTPASSWORD', '2': 4},
  ],
};

/// Descriptor for `EventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventTypeDescriptor = $convert.base64Decode(
    'CglFdmVudFR5cGUSFQoRRVZFTlRfVFlQRV9TSUdOSU4QABIZChVFVkVOVF9UWVBFX1JFU0VORE'
    'NPREUQARIdChlFVkVOVF9UWVBFX1BBU1NXT1JEQ0hBTkdFEAISFQoRRVZFTlRfVFlQRV9TSUdO'
    'VVAQAxIdChlFVkVOVF9UWVBFX0ZPUkdPVFBBU1NXT1JEEAQ=');

@$core.Deprecated('Use explicitAuthFlowsTypeDescriptor instead')
const ExplicitAuthFlowsType$json = {
  '1': 'ExplicitAuthFlowsType',
  '2': [
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_CUSTOM_AUTH_FLOW_ONLY', '2': 0},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_SRP_AUTH', '2': 1},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_CUSTOM_AUTH', '2': 2},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_AUTH', '2': 3},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_ADMIN_NO_SRP_AUTH', '2': 4},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_USER_PASSWORD_AUTH', '2': 5},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_ADMIN_USER_PASSWORD_AUTH', '2': 6},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_PASSWORD_AUTH', '2': 7},
    {'1': 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_REFRESH_TOKEN_AUTH', '2': 8},
  ],
};

/// Descriptor for `ExplicitAuthFlowsType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List explicitAuthFlowsTypeDescriptor = $convert.base64Decode(
    'ChVFeHBsaWNpdEF1dGhGbG93c1R5cGUSMgouRVhQTElDSVRfQVVUSF9GTE9XU19UWVBFX0NVU1'
    'RPTV9BVVRIX0ZMT1dfT05MWRAAEjAKLEVYUExJQ0lUX0FVVEhfRkxPV1NfVFlQRV9BTExPV19V'
    'U0VSX1NSUF9BVVRIEAESLgoqRVhQTElDSVRfQVVUSF9GTE9XU19UWVBFX0FMTE9XX0NVU1RPTV'
    '9BVVRIEAISLAooRVhQTElDSVRfQVVUSF9GTE9XU19UWVBFX0FMTE9XX1VTRVJfQVVUSBADEi4K'
    'KkVYUExJQ0lUX0FVVEhfRkxPV1NfVFlQRV9BRE1JTl9OT19TUlBfQVVUSBAEEi8KK0VYUExJQ0'
    'lUX0FVVEhfRkxPV1NfVFlQRV9VU0VSX1BBU1NXT1JEX0FVVEgQBRI7CjdFWFBMSUNJVF9BVVRI'
    'X0ZMT1dTX1RZUEVfQUxMT1dfQURNSU5fVVNFUl9QQVNTV09SRF9BVVRIEAYSNQoxRVhQTElDSV'
    'RfQVVUSF9GTE9XU19UWVBFX0FMTE9XX1VTRVJfUEFTU1dPUkRfQVVUSBAHEjUKMUVYUExJQ0lU'
    'X0FVVEhfRkxPV1NfVFlQRV9BTExPV19SRUZSRVNIX1RPS0VOX0FVVEgQCA==');

@$core.Deprecated('Use featureTypeDescriptor instead')
const FeatureType$json = {
  '1': 'FeatureType',
  '2': [
    {'1': 'FEATURE_TYPE_DISABLED', '2': 0},
    {'1': 'FEATURE_TYPE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `FeatureType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List featureTypeDescriptor = $convert.base64Decode(
    'CgtGZWF0dXJlVHlwZRIZChVGRUFUVVJFX1RZUEVfRElTQUJMRUQQABIYChRGRUFUVVJFX1RZUE'
    'VfRU5BQkxFRBAB');

@$core.Deprecated('Use feedbackValueTypeDescriptor instead')
const FeedbackValueType$json = {
  '1': 'FeedbackValueType',
  '2': [
    {'1': 'FEEDBACK_VALUE_TYPE_INVALID', '2': 0},
    {'1': 'FEEDBACK_VALUE_TYPE_VALID', '2': 1},
  ],
};

/// Descriptor for `FeedbackValueType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List feedbackValueTypeDescriptor = $convert.base64Decode(
    'ChFGZWVkYmFja1ZhbHVlVHlwZRIfChtGRUVEQkFDS19WQUxVRV9UWVBFX0lOVkFMSUQQABIdCh'
    'lGRUVEQkFDS19WQUxVRV9UWVBFX1ZBTElEEAE=');

@$core.Deprecated('Use identityProviderTypeTypeDescriptor instead')
const IdentityProviderTypeType$json = {
  '1': 'IdentityProviderTypeType',
  '2': [
    {'1': 'IDENTITY_PROVIDER_TYPE_TYPE_LOGINWITHAMAZON', '2': 0},
    {'1': 'IDENTITY_PROVIDER_TYPE_TYPE_GOOGLE', '2': 1},
    {'1': 'IDENTITY_PROVIDER_TYPE_TYPE_FACEBOOK', '2': 2},
    {'1': 'IDENTITY_PROVIDER_TYPE_TYPE_SAML', '2': 3},
    {'1': 'IDENTITY_PROVIDER_TYPE_TYPE_SIGNINWITHAPPLE', '2': 4},
    {'1': 'IDENTITY_PROVIDER_TYPE_TYPE_OIDC', '2': 5},
  ],
};

/// Descriptor for `IdentityProviderTypeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List identityProviderTypeTypeDescriptor = $convert.base64Decode(
    'ChhJZGVudGl0eVByb3ZpZGVyVHlwZVR5cGUSLworSURFTlRJVFlfUFJPVklERVJfVFlQRV9UWV'
    'BFX0xPR0lOV0lUSEFNQVpPThAAEiYKIklERU5USVRZX1BST1ZJREVSX1RZUEVfVFlQRV9HT09H'
    'TEUQARIoCiRJREVOVElUWV9QUk9WSURFUl9UWVBFX1RZUEVfRkFDRUJPT0sQAhIkCiBJREVOVE'
    'lUWV9QUk9WSURFUl9UWVBFX1RZUEVfU0FNTBADEi8KK0lERU5USVRZX1BST1ZJREVSX1RZUEVf'
    'VFlQRV9TSUdOSU5XSVRIQVBQTEUQBBIkCiBJREVOVElUWV9QUk9WSURFUl9UWVBFX1RZUEVfT0'
    'lEQxAF');

@$core.Deprecated('Use inboundFederationLambdaVersionTypeDescriptor instead')
const InboundFederationLambdaVersionType$json = {
  '1': 'InboundFederationLambdaVersionType',
  '2': [
    {'1': 'INBOUND_FEDERATION_LAMBDA_VERSION_TYPE_V1_0', '2': 0},
  ],
};

/// Descriptor for `InboundFederationLambdaVersionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inboundFederationLambdaVersionTypeDescriptor =
    $convert.base64Decode(
        'CiJJbmJvdW5kRmVkZXJhdGlvbkxhbWJkYVZlcnNpb25UeXBlEi8KK0lOQk9VTkRfRkVERVJBVE'
        'lPTl9MQU1CREFfVkVSU0lPTl9UWVBFX1YxXzAQAA==');

@$core.Deprecated('Use logLevelDescriptor instead')
const LogLevel$json = {
  '1': 'LogLevel',
  '2': [
    {'1': 'LOG_LEVEL_INFO', '2': 0},
    {'1': 'LOG_LEVEL_ERROR', '2': 1},
  ],
};

/// Descriptor for `LogLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List logLevelDescriptor = $convert.base64Decode(
    'CghMb2dMZXZlbBISCg5MT0dfTEVWRUxfSU5GTxAAEhMKD0xPR19MRVZFTF9FUlJPUhAB');

@$core.Deprecated('Use messageActionTypeDescriptor instead')
const MessageActionType$json = {
  '1': 'MessageActionType',
  '2': [
    {'1': 'MESSAGE_ACTION_TYPE_SUPPRESS', '2': 0},
    {'1': 'MESSAGE_ACTION_TYPE_RESEND', '2': 1},
  ],
};

/// Descriptor for `MessageActionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List messageActionTypeDescriptor = $convert.base64Decode(
    'ChFNZXNzYWdlQWN0aW9uVHlwZRIgChxNRVNTQUdFX0FDVElPTl9UWVBFX1NVUFBSRVNTEAASHg'
    'oaTUVTU0FHRV9BQ1RJT05fVFlQRV9SRVNFTkQQAQ==');

@$core.Deprecated('Use oAuthFlowTypeDescriptor instead')
const OAuthFlowType$json = {
  '1': 'OAuthFlowType',
  '2': [
    {'1': 'O_AUTH_FLOW_TYPE_CLIENT_CREDENTIALS', '2': 0},
    {'1': 'O_AUTH_FLOW_TYPE_IMPLICIT', '2': 1},
    {'1': 'O_AUTH_FLOW_TYPE_CODE', '2': 2},
  ],
};

/// Descriptor for `OAuthFlowType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List oAuthFlowTypeDescriptor = $convert.base64Decode(
    'Cg1PQXV0aEZsb3dUeXBlEicKI09fQVVUSF9GTE9XX1RZUEVfQ0xJRU5UX0NSRURFTlRJQUxTEA'
    'ASHQoZT19BVVRIX0ZMT1dfVFlQRV9JTVBMSUNJVBABEhkKFU9fQVVUSF9GTE9XX1RZUEVfQ09E'
    'RRAC');

@$core.Deprecated('Use preTokenGenerationLambdaVersionTypeDescriptor instead')
const PreTokenGenerationLambdaVersionType$json = {
  '1': 'PreTokenGenerationLambdaVersionType',
  '2': [
    {'1': 'PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V1_0', '2': 0},
    {'1': 'PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V3_0', '2': 1},
    {'1': 'PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V2_0', '2': 2},
  ],
};

/// Descriptor for `PreTokenGenerationLambdaVersionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List preTokenGenerationLambdaVersionTypeDescriptor =
    $convert.base64Decode(
        'CiNQcmVUb2tlbkdlbmVyYXRpb25MYW1iZGFWZXJzaW9uVHlwZRIxCi1QUkVfVE9LRU5fR0VORV'
        'JBVElPTl9MQU1CREFfVkVSU0lPTl9UWVBFX1YxXzAQABIxCi1QUkVfVE9LRU5fR0VORVJBVElP'
        'Tl9MQU1CREFfVkVSU0lPTl9UWVBFX1YzXzAQARIxCi1QUkVfVE9LRU5fR0VORVJBVElPTl9MQU'
        '1CREFfVkVSU0lPTl9UWVBFX1YyXzAQAg==');

@$core.Deprecated('Use preventUserExistenceErrorTypesDescriptor instead')
const PreventUserExistenceErrorTypes$json = {
  '1': 'PreventUserExistenceErrorTypes',
  '2': [
    {'1': 'PREVENT_USER_EXISTENCE_ERROR_TYPES_LEGACY', '2': 0},
    {'1': 'PREVENT_USER_EXISTENCE_ERROR_TYPES_ENABLED', '2': 1},
  ],
};

/// Descriptor for `PreventUserExistenceErrorTypes`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List preventUserExistenceErrorTypesDescriptor =
    $convert.base64Decode(
        'Ch5QcmV2ZW50VXNlckV4aXN0ZW5jZUVycm9yVHlwZXMSLQopUFJFVkVOVF9VU0VSX0VYSVNURU'
        '5DRV9FUlJPUl9UWVBFU19MRUdBQ1kQABIuCipQUkVWRU5UX1VTRVJfRVhJU1RFTkNFX0VSUk9S'
        'X1RZUEVTX0VOQUJMRUQQAQ==');

@$core.Deprecated('Use recoveryOptionNameTypeDescriptor instead')
const RecoveryOptionNameType$json = {
  '1': 'RecoveryOptionNameType',
  '2': [
    {'1': 'RECOVERY_OPTION_NAME_TYPE_ADMIN_ONLY', '2': 0},
    {'1': 'RECOVERY_OPTION_NAME_TYPE_VERIFIED_EMAIL', '2': 1},
    {'1': 'RECOVERY_OPTION_NAME_TYPE_VERIFIED_PHONE_NUMBER', '2': 2},
  ],
};

/// Descriptor for `RecoveryOptionNameType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recoveryOptionNameTypeDescriptor = $convert.base64Decode(
    'ChZSZWNvdmVyeU9wdGlvbk5hbWVUeXBlEigKJFJFQ09WRVJZX09QVElPTl9OQU1FX1RZUEVfQU'
    'RNSU5fT05MWRAAEiwKKFJFQ09WRVJZX09QVElPTl9OQU1FX1RZUEVfVkVSSUZJRURfRU1BSUwQ'
    'ARIzCi9SRUNPVkVSWV9PUFRJT05fTkFNRV9UWVBFX1ZFUklGSUVEX1BIT05FX05VTUJFUhAC');

@$core.Deprecated('Use riskDecisionTypeDescriptor instead')
const RiskDecisionType$json = {
  '1': 'RiskDecisionType',
  '2': [
    {'1': 'RISK_DECISION_TYPE_ACCOUNTTAKEOVER', '2': 0},
    {'1': 'RISK_DECISION_TYPE_NORISK', '2': 1},
    {'1': 'RISK_DECISION_TYPE_BLOCK', '2': 2},
  ],
};

/// Descriptor for `RiskDecisionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List riskDecisionTypeDescriptor = $convert.base64Decode(
    'ChBSaXNrRGVjaXNpb25UeXBlEiYKIlJJU0tfREVDSVNJT05fVFlQRV9BQ0NPVU5UVEFLRU9WRV'
    'IQABIdChlSSVNLX0RFQ0lTSU9OX1RZUEVfTk9SSVNLEAESHAoYUklTS19ERUNJU0lPTl9UWVBF'
    'X0JMT0NLEAI=');

@$core.Deprecated('Use riskLevelTypeDescriptor instead')
const RiskLevelType$json = {
  '1': 'RiskLevelType',
  '2': [
    {'1': 'RISK_LEVEL_TYPE_MEDIUM', '2': 0},
    {'1': 'RISK_LEVEL_TYPE_LOW', '2': 1},
    {'1': 'RISK_LEVEL_TYPE_HIGH', '2': 2},
  ],
};

/// Descriptor for `RiskLevelType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List riskLevelTypeDescriptor = $convert.base64Decode(
    'Cg1SaXNrTGV2ZWxUeXBlEhoKFlJJU0tfTEVWRUxfVFlQRV9NRURJVU0QABIXChNSSVNLX0xFVk'
    'VMX1RZUEVfTE9XEAESGAoUUklTS19MRVZFTF9UWVBFX0hJR0gQAg==');

@$core.Deprecated('Use statusTypeDescriptor instead')
const StatusType$json = {
  '1': 'StatusType',
  '2': [
    {'1': 'STATUS_TYPE_DISABLED', '2': 0},
    {'1': 'STATUS_TYPE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `StatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statusTypeDescriptor = $convert.base64Decode(
    'CgpTdGF0dXNUeXBlEhgKFFNUQVRVU19UWVBFX0RJU0FCTEVEEAASFwoTU1RBVFVTX1RZUEVfRU'
    '5BQkxFRBAB');

@$core.Deprecated('Use termsEnforcementTypeDescriptor instead')
const TermsEnforcementType$json = {
  '1': 'TermsEnforcementType',
  '2': [
    {'1': 'TERMS_ENFORCEMENT_TYPE_NONE', '2': 0},
  ],
};

/// Descriptor for `TermsEnforcementType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List termsEnforcementTypeDescriptor = $convert.base64Decode(
    'ChRUZXJtc0VuZm9yY2VtZW50VHlwZRIfChtURVJNU19FTkZPUkNFTUVOVF9UWVBFX05PTkUQAA'
    '==');

@$core.Deprecated('Use termsSourceTypeDescriptor instead')
const TermsSourceType$json = {
  '1': 'TermsSourceType',
  '2': [
    {'1': 'TERMS_SOURCE_TYPE_LINK', '2': 0},
  ],
};

/// Descriptor for `TermsSourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List termsSourceTypeDescriptor = $convert.base64Decode(
    'Cg9UZXJtc1NvdXJjZVR5cGUSGgoWVEVSTVNfU09VUkNFX1RZUEVfTElOSxAA');

@$core.Deprecated('Use timeUnitsTypeDescriptor instead')
const TimeUnitsType$json = {
  '1': 'TimeUnitsType',
  '2': [
    {'1': 'TIME_UNITS_TYPE_MINUTES', '2': 0},
    {'1': 'TIME_UNITS_TYPE_SECONDS', '2': 1},
    {'1': 'TIME_UNITS_TYPE_DAYS', '2': 2},
    {'1': 'TIME_UNITS_TYPE_HOURS', '2': 3},
  ],
};

/// Descriptor for `TimeUnitsType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List timeUnitsTypeDescriptor = $convert.base64Decode(
    'Cg1UaW1lVW5pdHNUeXBlEhsKF1RJTUVfVU5JVFNfVFlQRV9NSU5VVEVTEAASGwoXVElNRV9VTk'
    'lUU19UWVBFX1NFQ09ORFMQARIYChRUSU1FX1VOSVRTX1RZUEVfREFZUxACEhkKFVRJTUVfVU5J'
    'VFNfVFlQRV9IT1VSUxAD');

@$core.Deprecated('Use userImportJobStatusTypeDescriptor instead')
const UserImportJobStatusType$json = {
  '1': 'UserImportJobStatusType',
  '2': [
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_FAILED', '2': 0},
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_STOPPED', '2': 1},
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_EXPIRED', '2': 2},
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_SUCCEEDED', '2': 3},
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_STOPPING', '2': 4},
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_CREATED', '2': 5},
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_INPROGRESS', '2': 6},
    {'1': 'USER_IMPORT_JOB_STATUS_TYPE_PENDING', '2': 7},
  ],
};

/// Descriptor for `UserImportJobStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List userImportJobStatusTypeDescriptor = $convert.base64Decode(
    'ChdVc2VySW1wb3J0Sm9iU3RhdHVzVHlwZRImCiJVU0VSX0lNUE9SVF9KT0JfU1RBVFVTX1RZUE'
    'VfRkFJTEVEEAASJwojVVNFUl9JTVBPUlRfSk9CX1NUQVRVU19UWVBFX1NUT1BQRUQQARInCiNV'
    'U0VSX0lNUE9SVF9KT0JfU1RBVFVTX1RZUEVfRVhQSVJFRBACEikKJVVTRVJfSU1QT1JUX0pPQl'
    '9TVEFUVVNfVFlQRV9TVUNDRUVERUQQAxIoCiRVU0VSX0lNUE9SVF9KT0JfU1RBVFVTX1RZUEVf'
    'U1RPUFBJTkcQBBInCiNVU0VSX0lNUE9SVF9KT0JfU1RBVFVTX1RZUEVfQ1JFQVRFRBAFEioKJl'
    'VTRVJfSU1QT1JUX0pPQl9TVEFUVVNfVFlQRV9JTlBST0dSRVNTEAYSJwojVVNFUl9JTVBPUlRf'
    'Sk9CX1NUQVRVU19UWVBFX1BFTkRJTkcQBw==');

@$core.Deprecated('Use userPoolMfaTypeDescriptor instead')
const UserPoolMfaType$json = {
  '1': 'UserPoolMfaType',
  '2': [
    {'1': 'USER_POOL_MFA_TYPE_OPTIONAL', '2': 0},
    {'1': 'USER_POOL_MFA_TYPE_OFF', '2': 1},
    {'1': 'USER_POOL_MFA_TYPE_ON', '2': 2},
  ],
};

/// Descriptor for `UserPoolMfaType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List userPoolMfaTypeDescriptor = $convert.base64Decode(
    'Cg9Vc2VyUG9vbE1mYVR5cGUSHwobVVNFUl9QT09MX01GQV9UWVBFX09QVElPTkFMEAASGgoWVV'
    'NFUl9QT09MX01GQV9UWVBFX09GRhABEhkKFVVTRVJfUE9PTF9NRkFfVFlQRV9PThAC');

@$core.Deprecated('Use userPoolTierTypeDescriptor instead')
const UserPoolTierType$json = {
  '1': 'UserPoolTierType',
  '2': [
    {'1': 'USER_POOL_TIER_TYPE_PLUS', '2': 0},
    {'1': 'USER_POOL_TIER_TYPE_ESSENTIALS', '2': 1},
    {'1': 'USER_POOL_TIER_TYPE_LITE', '2': 2},
  ],
};

/// Descriptor for `UserPoolTierType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List userPoolTierTypeDescriptor = $convert.base64Decode(
    'ChBVc2VyUG9vbFRpZXJUeXBlEhwKGFVTRVJfUE9PTF9USUVSX1RZUEVfUExVUxAAEiIKHlVTRV'
    'JfUE9PTF9USUVSX1RZUEVfRVNTRU5USUFMUxABEhwKGFVTRVJfUE9PTF9USUVSX1RZUEVfTElU'
    'RRAC');

@$core.Deprecated('Use userStatusTypeDescriptor instead')
const UserStatusType$json = {
  '1': 'UserStatusType',
  '2': [
    {'1': 'USER_STATUS_TYPE_EXTERNAL_PROVIDER', '2': 0},
    {'1': 'USER_STATUS_TYPE_CONFIRMED', '2': 1},
    {'1': 'USER_STATUS_TYPE_UNKNOWN', '2': 2},
    {'1': 'USER_STATUS_TYPE_UNCONFIRMED', '2': 3},
    {'1': 'USER_STATUS_TYPE_RESET_REQUIRED', '2': 4},
    {'1': 'USER_STATUS_TYPE_COMPROMISED', '2': 5},
    {'1': 'USER_STATUS_TYPE_FORCE_CHANGE_PASSWORD', '2': 6},
    {'1': 'USER_STATUS_TYPE_ARCHIVED', '2': 7},
  ],
};

/// Descriptor for `UserStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List userStatusTypeDescriptor = $convert.base64Decode(
    'Cg5Vc2VyU3RhdHVzVHlwZRImCiJVU0VSX1NUQVRVU19UWVBFX0VYVEVSTkFMX1BST1ZJREVSEA'
    'ASHgoaVVNFUl9TVEFUVVNfVFlQRV9DT05GSVJNRUQQARIcChhVU0VSX1NUQVRVU19UWVBFX1VO'
    'S05PV04QAhIgChxVU0VSX1NUQVRVU19UWVBFX1VOQ09ORklSTUVEEAMSIwofVVNFUl9TVEFUVV'
    'NfVFlQRV9SRVNFVF9SRVFVSVJFRBAEEiAKHFVTRVJfU1RBVFVTX1RZUEVfQ09NUFJPTUlTRUQQ'
    'BRIqCiZVU0VSX1NUQVRVU19UWVBFX0ZPUkNFX0NIQU5HRV9QQVNTV09SRBAGEh0KGVVTRVJfU1'
    'RBVFVTX1RZUEVfQVJDSElWRUQQBw==');

@$core.Deprecated('Use userVerificationTypeDescriptor instead')
const UserVerificationType$json = {
  '1': 'UserVerificationType',
  '2': [
    {'1': 'USER_VERIFICATION_TYPE_PREFERRED', '2': 0},
    {'1': 'USER_VERIFICATION_TYPE_REQUIRED', '2': 1},
  ],
};

/// Descriptor for `UserVerificationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List userVerificationTypeDescriptor = $convert.base64Decode(
    'ChRVc2VyVmVyaWZpY2F0aW9uVHlwZRIkCiBVU0VSX1ZFUklGSUNBVElPTl9UWVBFX1BSRUZFUl'
    'JFRBAAEiMKH1VTRVJfVkVSSUZJQ0FUSU9OX1RZUEVfUkVRVUlSRUQQAQ==');

@$core.Deprecated('Use usernameAttributeTypeDescriptor instead')
const UsernameAttributeType$json = {
  '1': 'UsernameAttributeType',
  '2': [
    {'1': 'USERNAME_ATTRIBUTE_TYPE_PHONE_NUMBER', '2': 0},
    {'1': 'USERNAME_ATTRIBUTE_TYPE_EMAIL', '2': 1},
  ],
};

/// Descriptor for `UsernameAttributeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List usernameAttributeTypeDescriptor = $convert.base64Decode(
    'ChVVc2VybmFtZUF0dHJpYnV0ZVR5cGUSKAokVVNFUk5BTUVfQVRUUklCVVRFX1RZUEVfUEhPTk'
    'VfTlVNQkVSEAASIQodVVNFUk5BTUVfQVRUUklCVVRFX1RZUEVfRU1BSUwQAQ==');

@$core.Deprecated('Use verifiedAttributeTypeDescriptor instead')
const VerifiedAttributeType$json = {
  '1': 'VerifiedAttributeType',
  '2': [
    {'1': 'VERIFIED_ATTRIBUTE_TYPE_PHONE_NUMBER', '2': 0},
    {'1': 'VERIFIED_ATTRIBUTE_TYPE_EMAIL', '2': 1},
  ],
};

/// Descriptor for `VerifiedAttributeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List verifiedAttributeTypeDescriptor = $convert.base64Decode(
    'ChVWZXJpZmllZEF0dHJpYnV0ZVR5cGUSKAokVkVSSUZJRURfQVRUUklCVVRFX1RZUEVfUEhPTk'
    'VfTlVNQkVSEAASIQodVkVSSUZJRURfQVRUUklCVVRFX1RZUEVfRU1BSUwQAQ==');

@$core.Deprecated('Use verifySoftwareTokenResponseTypeDescriptor instead')
const VerifySoftwareTokenResponseType$json = {
  '1': 'VerifySoftwareTokenResponseType',
  '2': [
    {'1': 'VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_SUCCESS', '2': 0},
    {'1': 'VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_ERROR', '2': 1},
  ],
};

/// Descriptor for `VerifySoftwareTokenResponseType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List verifySoftwareTokenResponseTypeDescriptor =
    $convert.base64Decode(
        'Ch9WZXJpZnlTb2Z0d2FyZVRva2VuUmVzcG9uc2VUeXBlEi8KK1ZFUklGWV9TT0ZUV0FSRV9UT0'
        'tFTl9SRVNQT05TRV9UWVBFX1NVQ0NFU1MQABItCilWRVJJRllfU09GVFdBUkVfVE9LRU5fUkVT'
        'UE9OU0VfVFlQRV9FUlJPUhAB');

@$core.Deprecated('Use accessDeniedExceptionDescriptor instead')
const AccessDeniedException$json = {
  '1': 'AccessDeniedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedExceptionDescriptor = $convert.base64Decode(
    'ChVBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use accountRecoverySettingTypeDescriptor instead')
const AccountRecoverySettingType$json = {
  '1': 'AccountRecoverySettingType',
  '2': [
    {
      '1': 'recoverymechanisms',
      '3': 294786489,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.RecoveryOptionType',
      '10': 'recoverymechanisms'
    },
  ],
};

/// Descriptor for `AccountRecoverySettingType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountRecoverySettingTypeDescriptor =
    $convert.base64Decode(
        'ChpBY2NvdW50UmVjb3ZlcnlTZXR0aW5nVHlwZRJTChJyZWNvdmVyeW1lY2hhbmlzbXMYuavIjA'
        'EgAygLMh8uY29nbml0b19pZHAuUmVjb3ZlcnlPcHRpb25UeXBlUhJyZWNvdmVyeW1lY2hhbmlz'
        'bXM=');

@$core.Deprecated('Use accountTakeoverActionTypeDescriptor instead')
const AccountTakeoverActionType$json = {
  '1': 'AccountTakeoverActionType',
  '2': [
    {
      '1': 'eventaction',
      '3': 85257574,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AccountTakeoverEventActionType',
      '10': 'eventaction'
    },
    {'1': 'notify', '3': 314575197, '4': 1, '5': 8, '10': 'notify'},
  ],
};

/// Descriptor for `AccountTakeoverActionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountTakeoverActionTypeDescriptor = $convert.base64Decode(
    'ChlBY2NvdW50VGFrZW92ZXJBY3Rpb25UeXBlElAKC2V2ZW50YWN0aW9uGOba0yggASgOMisuY2'
    '9nbml0b19pZHAuQWNjb3VudFRha2VvdmVyRXZlbnRBY3Rpb25UeXBlUgtldmVudGFjdGlvbhIa'
    'CgZub3RpZnkY3ZKAlgEgASgIUgZub3RpZnk=');

@$core.Deprecated('Use accountTakeoverActionsTypeDescriptor instead')
const AccountTakeoverActionsType$json = {
  '1': 'AccountTakeoverActionsType',
  '2': [
    {
      '1': 'highaction',
      '3': 234016520,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountTakeoverActionType',
      '10': 'highaction'
    },
    {
      '1': 'lowaction',
      '3': 281569908,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountTakeoverActionType',
      '10': 'lowaction'
    },
    {
      '1': 'mediumaction',
      '3': 391798941,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountTakeoverActionType',
      '10': 'mediumaction'
    },
  ],
};

/// Descriptor for `AccountTakeoverActionsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountTakeoverActionsTypeDescriptor = $convert.base64Decode(
    'ChpBY2NvdW50VGFrZW92ZXJBY3Rpb25zVHlwZRJJCgpoaWdoYWN0aW9uGIiey28gASgLMiYuY2'
    '9nbml0b19pZHAuQWNjb3VudFRha2VvdmVyQWN0aW9uVHlwZVIKaGlnaGFjdGlvbhJICglsb3dh'
    'Y3Rpb24Y9NShhgEgASgLMiYuY29nbml0b19pZHAuQWNjb3VudFRha2VvdmVyQWN0aW9uVHlwZV'
    'IJbG93YWN0aW9uEk4KDG1lZGl1bWFjdGlvbhidwem6ASABKAsyJi5jb2duaXRvX2lkcC5BY2Nv'
    'dW50VGFrZW92ZXJBY3Rpb25UeXBlUgxtZWRpdW1hY3Rpb24=');

@$core.Deprecated('Use accountTakeoverRiskConfigurationTypeDescriptor instead')
const AccountTakeoverRiskConfigurationType$json = {
  '1': 'AccountTakeoverRiskConfigurationType',
  '2': [
    {
      '1': 'actions',
      '3': 106935557,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountTakeoverActionsType',
      '10': 'actions'
    },
    {
      '1': 'notifyconfiguration',
      '3': 203106647,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.NotifyConfigurationType',
      '10': 'notifyconfiguration'
    },
  ],
};

/// Descriptor for `AccountTakeoverRiskConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountTakeoverRiskConfigurationTypeDescriptor =
    $convert.base64Decode(
        'CiRBY2NvdW50VGFrZW92ZXJSaXNrQ29uZmlndXJhdGlvblR5cGUSRAoHYWN0aW9ucxiF6v4yIA'
        'EoCzInLmNvZ25pdG9faWRwLkFjY291bnRUYWtlb3ZlckFjdGlvbnNUeXBlUgdhY3Rpb25zElkK'
        'E25vdGlmeWNvbmZpZ3VyYXRpb24Y19LsYCABKAsyJC5jb2duaXRvX2lkcC5Ob3RpZnlDb25maW'
        'd1cmF0aW9uVHlwZVITbm90aWZ5Y29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use addCustomAttributesRequestDescriptor instead')
const AddCustomAttributesRequest$json = {
  '1': 'AddCustomAttributesRequest',
  '2': [
    {
      '1': 'customattributes',
      '3': 423725708,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.SchemaAttributeType',
      '10': 'customattributes'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `AddCustomAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addCustomAttributesRequestDescriptor =
    $convert.base64Decode(
        'ChpBZGRDdXN0b21BdHRyaWJ1dGVzUmVxdWVzdBJQChBjdXN0b21hdHRyaWJ1dGVzGIyVhsoBIA'
        'MoCzIgLmNvZ25pdG9faWRwLlNjaGVtYUF0dHJpYnV0ZVR5cGVSEGN1c3RvbWF0dHJpYnV0ZXMS'
        'IgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use addCustomAttributesResponseDescriptor instead')
const AddCustomAttributesResponse$json = {
  '1': 'AddCustomAttributesResponse',
};

/// Descriptor for `AddCustomAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addCustomAttributesResponseDescriptor =
    $convert.base64Decode('ChtBZGRDdXN0b21BdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated('Use addUserPoolClientSecretRequestDescriptor instead')
const AddUserPoolClientSecretRequest$json = {
  '1': 'AddUserPoolClientSecretRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'clientsecret', '3': 500734711, '4': 1, '5': 9, '10': 'clientsecret'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `AddUserPoolClientSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addUserPoolClientSecretRequestDescriptor =
    $convert.base64Decode(
        'Ch5BZGRVc2VyUG9vbENsaWVudFNlY3JldFJlcXVlc3QSHgoIY2xpZW50aWQYpOiG1gEgASgJUg'
        'hjbGllbnRpZBImCgxjbGllbnRzZWNyZXQY97Xi7gEgASgJUgxjbGllbnRzZWNyZXQSIgoKdXNl'
        'cnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use addUserPoolClientSecretResponseDescriptor instead')
const AddUserPoolClientSecretResponse$json = {
  '1': 'AddUserPoolClientSecretResponse',
  '2': [
    {
      '1': 'clientsecretdescriptor',
      '3': 428044828,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ClientSecretDescriptorType',
      '10': 'clientsecretdescriptor'
    },
  ],
};

/// Descriptor for `AddUserPoolClientSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addUserPoolClientSecretResponseDescriptor =
    $convert.base64Decode(
        'Ch9BZGRVc2VyUG9vbENsaWVudFNlY3JldFJlc3BvbnNlEmMKFmNsaWVudHNlY3JldGRlc2NyaX'
        'B0b3IYnOSNzAEgASgLMicuY29nbml0b19pZHAuQ2xpZW50U2VjcmV0RGVzY3JpcHRvclR5cGVS'
        'FmNsaWVudHNlY3JldGRlc2NyaXB0b3I=');

@$core.Deprecated('Use adminAddUserToGroupRequestDescriptor instead')
const AdminAddUserToGroupRequest$json = {
  '1': 'AdminAddUserToGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminAddUserToGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminAddUserToGroupRequestDescriptor =
    $convert.base64Decode(
        'ChpBZG1pbkFkZFVzZXJUb0dyb3VwUmVxdWVzdBIgCglncm91cG5hbWUYyMqgqgEgASgJUglncm'
        '91cG5hbWUSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQSHgoIdXNlcm5hbWUY'
        '2qmj4AEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use adminConfirmSignUpRequestDescriptor instead')
const AdminConfirmSignUpRequest$json = {
  '1': 'AdminConfirmSignUpRequest',
  '2': [
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AdminConfirmSignUpRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [AdminConfirmSignUpRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use adminConfirmSignUpRequestDescriptor instead')
const AdminConfirmSignUpRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminConfirmSignUpRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminConfirmSignUpRequestDescriptor = $convert.base64Decode(
    'ChlBZG1pbkNvbmZpcm1TaWduVXBSZXF1ZXN0EmUKDmNsaWVudG1ldGFkYXRhGMyv/2EgAygLMj'
    'ouY29nbml0b19pZHAuQWRtaW5Db25maXJtU2lnblVwUmVxdWVzdC5DbGllbnRtZXRhZGF0YUVu'
    'dHJ5Ug5jbGllbnRtZXRhZGF0YRIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZB'
    'IeCgh1c2VybmFtZRjaqaPgASABKAlSCHVzZXJuYW1lGkEKE0NsaWVudG1ldGFkYXRhRW50cnkS'
    'EAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use adminConfirmSignUpResponseDescriptor instead')
const AdminConfirmSignUpResponse$json = {
  '1': 'AdminConfirmSignUpResponse',
};

/// Descriptor for `AdminConfirmSignUpResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminConfirmSignUpResponseDescriptor =
    $convert.base64Decode('ChpBZG1pbkNvbmZpcm1TaWduVXBSZXNwb25zZQ==');

@$core.Deprecated('Use adminCreateUserConfigTypeDescriptor instead')
const AdminCreateUserConfigType$json = {
  '1': 'AdminCreateUserConfigType',
  '2': [
    {
      '1': 'allowadmincreateuseronly',
      '3': 399154521,
      '4': 1,
      '5': 8,
      '10': 'allowadmincreateuseronly'
    },
    {
      '1': 'invitemessagetemplate',
      '3': 162590438,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.MessageTemplateType',
      '10': 'invitemessagetemplate'
    },
    {
      '1': 'unusedaccountvaliditydays',
      '3': 354236064,
      '4': 1,
      '5': 5,
      '10': 'unusedaccountvaliditydays'
    },
  ],
};

/// Descriptor for `AdminCreateUserConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminCreateUserConfigTypeDescriptor = $convert.base64Decode(
    'ChlBZG1pbkNyZWF0ZVVzZXJDb25maWdUeXBlEj4KGGFsbG93YWRtaW5jcmVhdGV1c2Vyb25seR'
    'jZuqq+ASABKAhSGGFsbG93YWRtaW5jcmVhdGV1c2Vyb25seRJZChVpbnZpdGVtZXNzYWdldGVt'
    'cGxhdGUY5t3DTSABKAsyIC5jb2duaXRvX2lkcC5NZXNzYWdlVGVtcGxhdGVUeXBlUhVpbnZpdG'
    'VtZXNzYWdldGVtcGxhdGUSQAoZdW51c2VkYWNjb3VudHZhbGlkaXR5ZGF5cxig7fSoASABKAVS'
    'GXVudXNlZGFjY291bnR2YWxpZGl0eWRheXM=');

@$core.Deprecated('Use adminCreateUserRequestDescriptor instead')
const AdminCreateUserRequest$json = {
  '1': 'AdminCreateUserRequest',
  '2': [
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AdminCreateUserRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {
      '1': 'desireddeliverymediums',
      '3': 82642144,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.DeliveryMediumType',
      '10': 'desireddeliverymediums'
    },
    {
      '1': 'forcealiascreation',
      '3': 449762314,
      '4': 1,
      '5': 8,
      '10': 'forcealiascreation'
    },
    {
      '1': 'messageaction',
      '3': 233944615,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.MessageActionType',
      '10': 'messageaction'
    },
    {
      '1': 'temporarypassword',
      '3': 274827278,
      '4': 1,
      '5': 9,
      '10': 'temporarypassword'
    },
    {
      '1': 'userattributes',
      '3': 194667064,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'userattributes'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
    {
      '1': 'validationdata',
      '3': 235118411,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'validationdata'
    },
  ],
  '3': [AdminCreateUserRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use adminCreateUserRequestDescriptor instead')
const AdminCreateUserRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminCreateUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminCreateUserRequestDescriptor = $convert.base64Decode(
    'ChZBZG1pbkNyZWF0ZVVzZXJSZXF1ZXN0EmIKDmNsaWVudG1ldGFkYXRhGMyv/2EgAygLMjcuY2'
    '9nbml0b19pZHAuQWRtaW5DcmVhdGVVc2VyUmVxdWVzdC5DbGllbnRtZXRhZGF0YUVudHJ5Ug5j'
    'bGllbnRtZXRhZGF0YRJaChZkZXNpcmVkZGVsaXZlcnltZWRpdW1zGOCJtCcgAygOMh8uY29nbm'
    'l0b19pZHAuRGVsaXZlcnlNZWRpdW1UeXBlUhZkZXNpcmVkZGVsaXZlcnltZWRpdW1zEjIKEmZv'
    'cmNlYWxpYXNjcmVhdGlvbhiKqLvWASABKAhSEmZvcmNlYWxpYXNjcmVhdGlvbhJHCg1tZXNzYW'
    'dlYWN0aW9uGKfsxm8gASgOMh4uY29nbml0b19pZHAuTWVzc2FnZUFjdGlvblR5cGVSDW1lc3Nh'
    'Z2VhY3Rpb24SMAoRdGVtcG9yYXJ5cGFzc3dvcmQYjpCGgwEgASgJUhF0ZW1wb3JhcnlwYXNzd2'
    '9yZBJFCg51c2VyYXR0cmlidXRlcxi4xOlcIAMoCzIaLmNvZ25pdG9faWRwLkF0dHJpYnV0ZVR5'
    'cGVSDnVzZXJhdHRyaWJ1dGVzEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlkEh'
    '4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWUSRQoOdmFsaWRhdGlvbmRhdGEYy76OcCAD'
    'KAsyGi5jb2duaXRvX2lkcC5BdHRyaWJ1dGVUeXBlUg52YWxpZGF0aW9uZGF0YRpBChNDbGllbn'
    'RtZXRhZGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToC'
    'OAE=');

@$core.Deprecated('Use adminCreateUserResponseDescriptor instead')
const AdminCreateUserResponse$json = {
  '1': 'AdminCreateUserResponse',
  '2': [
    {
      '1': 'user',
      '3': 10894867,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserType',
      '10': 'user'
    },
  ],
};

/// Descriptor for `AdminCreateUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminCreateUserResponseDescriptor =
    $convert.base64Decode(
        'ChdBZG1pbkNyZWF0ZVVzZXJSZXNwb25zZRIsCgR1c2VyGJP8mAUgASgLMhUuY29nbml0b19pZH'
        'AuVXNlclR5cGVSBHVzZXI=');

@$core.Deprecated('Use adminDeleteUserAttributesRequestDescriptor instead')
const AdminDeleteUserAttributesRequest$json = {
  '1': 'AdminDeleteUserAttributesRequest',
  '2': [
    {
      '1': 'userattributenames',
      '3': 132104459,
      '4': 3,
      '5': 9,
      '10': 'userattributenames'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminDeleteUserAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminDeleteUserAttributesRequestDescriptor =
    $convert.base64Decode(
        'CiBBZG1pbkRlbGV0ZVVzZXJBdHRyaWJ1dGVzUmVxdWVzdBIxChJ1c2VyYXR0cmlidXRlbmFtZX'
        'MYi4L/PiADKAlSEnVzZXJhdHRyaWJ1dGVuYW1lcxIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIK'
        'dXNlcnBvb2xpZBIeCgh1c2VybmFtZRjaqaPgASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use adminDeleteUserAttributesResponseDescriptor instead')
const AdminDeleteUserAttributesResponse$json = {
  '1': 'AdminDeleteUserAttributesResponse',
};

/// Descriptor for `AdminDeleteUserAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminDeleteUserAttributesResponseDescriptor =
    $convert.base64Decode('CiFBZG1pbkRlbGV0ZVVzZXJBdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated('Use adminDeleteUserRequestDescriptor instead')
const AdminDeleteUserRequest$json = {
  '1': 'AdminDeleteUserRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminDeleteUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminDeleteUserRequestDescriptor =
    $convert.base64Decode(
        'ChZBZG1pbkRlbGV0ZVVzZXJSZXF1ZXN0EiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG'
        '9vbGlkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminDisableProviderForUserRequestDescriptor instead')
const AdminDisableProviderForUserRequest$json = {
  '1': 'AdminDisableProviderForUserRequest',
  '2': [
    {
      '1': 'user',
      '3': 10894867,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ProviderUserIdentifierType',
      '10': 'user'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `AdminDisableProviderForUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminDisableProviderForUserRequestDescriptor =
    $convert.base64Decode(
        'CiJBZG1pbkRpc2FibGVQcm92aWRlckZvclVzZXJSZXF1ZXN0Ej4KBHVzZXIYk/yYBSABKAsyJy'
        '5jb2duaXRvX2lkcC5Qcm92aWRlclVzZXJJZGVudGlmaWVyVHlwZVIEdXNlchIiCgp1c2VycG9v'
        'bGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use adminDisableProviderForUserResponseDescriptor instead')
const AdminDisableProviderForUserResponse$json = {
  '1': 'AdminDisableProviderForUserResponse',
};

/// Descriptor for `AdminDisableProviderForUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminDisableProviderForUserResponseDescriptor =
    $convert
        .base64Decode('CiNBZG1pbkRpc2FibGVQcm92aWRlckZvclVzZXJSZXNwb25zZQ==');

@$core.Deprecated('Use adminDisableUserRequestDescriptor instead')
const AdminDisableUserRequest$json = {
  '1': 'AdminDisableUserRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminDisableUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminDisableUserRequestDescriptor =
    $convert.base64Decode(
        'ChdBZG1pbkRpc2FibGVVc2VyUmVxdWVzdBIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcn'
        'Bvb2xpZBIeCgh1c2VybmFtZRjaqaPgASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use adminDisableUserResponseDescriptor instead')
const AdminDisableUserResponse$json = {
  '1': 'AdminDisableUserResponse',
};

/// Descriptor for `AdminDisableUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminDisableUserResponseDescriptor =
    $convert.base64Decode('ChhBZG1pbkRpc2FibGVVc2VyUmVzcG9uc2U=');

@$core.Deprecated('Use adminEnableUserRequestDescriptor instead')
const AdminEnableUserRequest$json = {
  '1': 'AdminEnableUserRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminEnableUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminEnableUserRequestDescriptor =
    $convert.base64Decode(
        'ChZBZG1pbkVuYWJsZVVzZXJSZXF1ZXN0EiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG'
        '9vbGlkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminEnableUserResponseDescriptor instead')
const AdminEnableUserResponse$json = {
  '1': 'AdminEnableUserResponse',
};

/// Descriptor for `AdminEnableUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminEnableUserResponseDescriptor =
    $convert.base64Decode('ChdBZG1pbkVuYWJsZVVzZXJSZXNwb25zZQ==');

@$core.Deprecated('Use adminForgetDeviceRequestDescriptor instead')
const AdminForgetDeviceRequest$json = {
  '1': 'AdminForgetDeviceRequest',
  '2': [
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminForgetDeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminForgetDeviceRequestDescriptor = $convert.base64Decode(
    'ChhBZG1pbkZvcmdldERldmljZVJlcXVlc3QSIAoJZGV2aWNla2V5GKvkyLYBIAEoCVIJZGV2aW'
    'Nla2V5EiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlkEh4KCHVzZXJuYW1lGNqp'
    'o+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminGetDeviceRequestDescriptor instead')
const AdminGetDeviceRequest$json = {
  '1': 'AdminGetDeviceRequest',
  '2': [
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminGetDeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminGetDeviceRequestDescriptor = $convert.base64Decode(
    'ChVBZG1pbkdldERldmljZVJlcXVlc3QSIAoJZGV2aWNla2V5GKvkyLYBIAEoCVIJZGV2aWNla2'
    'V5EiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlkEh4KCHVzZXJuYW1lGNqpo+AB'
    'IAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminGetDeviceResponseDescriptor instead')
const AdminGetDeviceResponse$json = {
  '1': 'AdminGetDeviceResponse',
  '2': [
    {
      '1': 'device',
      '3': 499889172,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.DeviceType',
      '10': 'device'
    },
  ],
};

/// Descriptor for `AdminGetDeviceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminGetDeviceResponseDescriptor =
    $convert.base64Decode(
        'ChZBZG1pbkdldERldmljZVJlc3BvbnNlEjMKBmRldmljZRiU6K7uASABKAsyFy5jb2duaXRvX2'
        'lkcC5EZXZpY2VUeXBlUgZkZXZpY2U=');

@$core.Deprecated('Use adminGetUserRequestDescriptor instead')
const AdminGetUserRequest$json = {
  '1': 'AdminGetUserRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminGetUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminGetUserRequestDescriptor = $convert.base64Decode(
    'ChNBZG1pbkdldFVzZXJSZXF1ZXN0EiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbG'
    'lkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminGetUserResponseDescriptor instead')
const AdminGetUserResponse$json = {
  '1': 'AdminGetUserResponse',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'mfaoptions',
      '3': 501540826,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.MFAOptionType',
      '10': 'mfaoptions'
    },
    {
      '1': 'preferredmfasetting',
      '3': 111271921,
      '4': 1,
      '5': 9,
      '10': 'preferredmfasetting'
    },
    {
      '1': 'userattributes',
      '3': 194667064,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'userattributes'
    },
    {
      '1': 'usercreatedate',
      '3': 73013267,
      '4': 1,
      '5': 9,
      '10': 'usercreatedate'
    },
    {
      '1': 'userlastmodifieddate',
      '3': 80916802,
      '4': 1,
      '5': 9,
      '10': 'userlastmodifieddate'
    },
    {
      '1': 'usermfasettinglist',
      '3': 230885,
      '4': 3,
      '5': 9,
      '10': 'usermfasettinglist'
    },
    {
      '1': 'userstatus',
      '3': 189848701,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserStatusType',
      '10': 'userstatus'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminGetUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminGetUserResponseDescriptor = $convert.base64Decode(
    'ChRBZG1pbkdldFVzZXJSZXNwb25zZRIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBI+Cg'
    'ptZmFvcHRpb25zGNrPk+8BIAMoCzIaLmNvZ25pdG9faWRwLk1GQU9wdGlvblR5cGVSCm1mYW9w'
    'dGlvbnMSMwoTcHJlZmVycmVkbWZhc2V0dGluZxjxv4c1IAEoCVITcHJlZmVycmVkbWZhc2V0dG'
    'luZxJFCg51c2VyYXR0cmlidXRlcxi4xOlcIAMoCzIaLmNvZ25pdG9faWRwLkF0dHJpYnV0ZVR5'
    'cGVSDnVzZXJhdHRyaWJ1dGVzEikKDnVzZXJjcmVhdGVkYXRlGJOw6CIgASgJUg51c2VyY3JlYX'
    'RlZGF0ZRI1ChR1c2VybGFzdG1vZGlmaWVkZGF0ZRjC4somIAEoCVIUdXNlcmxhc3Rtb2RpZmll'
    'ZGRhdGUSMAoSdXNlcm1mYXNldHRpbmdsaXN0GOWLDiADKAlSEnVzZXJtZmFzZXR0aW5nbGlzdB'
    'I+Cgp1c2Vyc3RhdHVzGP24w1ogASgOMhsuY29nbml0b19pZHAuVXNlclN0YXR1c1R5cGVSCnVz'
    'ZXJzdGF0dXMSHgoIdXNlcm5hbWUY2qmj4AEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use adminInitiateAuthRequestDescriptor instead')
const AdminInitiateAuthRequest$json = {
  '1': 'AdminInitiateAuthRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {
      '1': 'authflow',
      '3': 143701420,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AuthFlowType',
      '10': 'authflow'
    },
    {
      '1': 'authparameters',
      '3': 258276552,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AdminInitiateAuthRequest.AuthparametersEntry',
      '10': 'authparameters'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AdminInitiateAuthRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {
      '1': 'contextdata',
      '3': 235700261,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ContextDataType',
      '10': 'contextdata'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [
    AdminInitiateAuthRequest_AuthparametersEntry$json,
    AdminInitiateAuthRequest_ClientmetadataEntry$json
  ],
};

@$core.Deprecated('Use adminInitiateAuthRequestDescriptor instead')
const AdminInitiateAuthRequest_AuthparametersEntry$json = {
  '1': 'AuthparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use adminInitiateAuthRequestDescriptor instead')
const AdminInitiateAuthRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminInitiateAuthRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminInitiateAuthRequestDescriptor = $convert.base64Decode(
    'ChhBZG1pbkluaXRpYXRlQXV0aFJlcXVlc3QSUwoRYW5hbHl0aWNzbWV0YWRhdGEY94XBPCABKA'
    'syIi5jb2duaXRvX2lkcC5BbmFseXRpY3NNZXRhZGF0YVR5cGVSEWFuYWx5dGljc21ldGFkYXRh'
    'EjgKCGF1dGhmbG93GKzrwkQgASgOMhkuY29nbml0b19pZHAuQXV0aEZsb3dUeXBlUghhdXRoZm'
    'xvdxJkCg5hdXRocGFyYW1ldGVycxjI+ZN7IAMoCzI5LmNvZ25pdG9faWRwLkFkbWluSW5pdGlh'
    'dGVBdXRoUmVxdWVzdC5BdXRocGFyYW1ldGVyc0VudHJ5Ug5hdXRocGFyYW1ldGVycxIeCghjbG'
    'llbnRpZBik6IbWASABKAlSCGNsaWVudGlkEmQKDmNsaWVudG1ldGFkYXRhGMyv/2EgAygLMjku'
    'Y29nbml0b19pZHAuQWRtaW5Jbml0aWF0ZUF1dGhSZXF1ZXN0LkNsaWVudG1ldGFkYXRhRW50cn'
    'lSDmNsaWVudG1ldGFkYXRhEkEKC2NvbnRleHRkYXRhGKWAsnAgASgLMhwuY29nbml0b19pZHAu'
    'Q29udGV4dERhdGFUeXBlUgtjb250ZXh0ZGF0YRIbCgdzZXNzaW9uGJiZowIgASgJUgdzZXNzaW'
    '9uEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlkGkEKE0F1dGhwYXJhbWV0ZXJz'
    'RW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpBChNDbG'
    'llbnRtZXRhZGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1'
    'ZToCOAE=');

@$core.Deprecated('Use adminInitiateAuthResponseDescriptor instead')
const AdminInitiateAuthResponse$json = {
  '1': 'AdminInitiateAuthResponse',
  '2': [
    {
      '1': 'authenticationresult',
      '3': 519327313,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AuthenticationResultType',
      '10': 'authenticationresult'
    },
    {
      '1': 'availablechallenges',
      '3': 275725719,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'availablechallenges'
    },
    {
      '1': 'challengename',
      '3': 170761310,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'challengename'
    },
    {
      '1': 'challengeparameters',
      '3': 79757023,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AdminInitiateAuthResponse.ChallengeparametersEntry',
      '10': 'challengeparameters'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
  ],
  '3': [AdminInitiateAuthResponse_ChallengeparametersEntry$json],
};

@$core.Deprecated('Use adminInitiateAuthResponseDescriptor instead')
const AdminInitiateAuthResponse_ChallengeparametersEntry$json = {
  '1': 'ChallengeparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminInitiateAuthResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminInitiateAuthResponseDescriptor = $convert.base64Decode(
    'ChlBZG1pbkluaXRpYXRlQXV0aFJlc3BvbnNlEl0KFGF1dGhlbnRpY2F0aW9ucmVzdWx0GNGc0f'
    'cBIAEoCzIlLmNvZ25pdG9faWRwLkF1dGhlbnRpY2F0aW9uUmVzdWx0VHlwZVIUYXV0aGVudGlj'
    'YXRpb25yZXN1bHQSVAoTYXZhaWxhYmxlY2hhbGxlbmdlcxiX+7yDASADKA4yHi5jb2duaXRvX2'
    'lkcC5DaGFsbGVuZ2VOYW1lVHlwZVITYXZhaWxhYmxlY2hhbGxlbmdlcxJHCg1jaGFsbGVuZ2Vu'
    'YW1lGN64tlEgASgOMh4uY29nbml0b19pZHAuQ2hhbGxlbmdlTmFtZVR5cGVSDWNoYWxsZW5nZW'
    '5hbWUSdAoTY2hhbGxlbmdlcGFyYW1ldGVycxjf/YMmIAMoCzI/LmNvZ25pdG9faWRwLkFkbWlu'
    'SW5pdGlhdGVBdXRoUmVzcG9uc2UuQ2hhbGxlbmdlcGFyYW1ldGVyc0VudHJ5UhNjaGFsbGVuZ2'
    'VwYXJhbWV0ZXJzEhsKB3Nlc3Npb24YmJmjAiABKAlSB3Nlc3Npb24aRgoYQ2hhbGxlbmdlcGFy'
    'YW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOA'
    'E=');

@$core.Deprecated('Use adminLinkProviderForUserRequestDescriptor instead')
const AdminLinkProviderForUserRequest$json = {
  '1': 'AdminLinkProviderForUserRequest',
  '2': [
    {
      '1': 'destinationuser',
      '3': 457045229,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ProviderUserIdentifierType',
      '10': 'destinationuser'
    },
    {
      '1': 'sourceuser',
      '3': 207754844,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ProviderUserIdentifierType',
      '10': 'sourceuser'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `AdminLinkProviderForUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminLinkProviderForUserRequestDescriptor = $convert.base64Decode(
    'Ch9BZG1pbkxpbmtQcm92aWRlckZvclVzZXJSZXF1ZXN0ElUKD2Rlc3RpbmF0aW9udXNlchjt6f'
    'fZASABKAsyJy5jb2duaXRvX2lkcC5Qcm92aWRlclVzZXJJZGVudGlmaWVyVHlwZVIPZGVzdGlu'
    'YXRpb251c2VyEkoKCnNvdXJjZXVzZXIY3KyIYyABKAsyJy5jb2duaXRvX2lkcC5Qcm92aWRlcl'
    'VzZXJJZGVudGlmaWVyVHlwZVIKc291cmNldXNlchIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIK'
    'dXNlcnBvb2xpZA==');

@$core.Deprecated('Use adminLinkProviderForUserResponseDescriptor instead')
const AdminLinkProviderForUserResponse$json = {
  '1': 'AdminLinkProviderForUserResponse',
};

/// Descriptor for `AdminLinkProviderForUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminLinkProviderForUserResponseDescriptor =
    $convert.base64Decode('CiBBZG1pbkxpbmtQcm92aWRlckZvclVzZXJSZXNwb25zZQ==');

@$core.Deprecated('Use adminListDevicesRequestDescriptor instead')
const AdminListDevicesRequest$json = {
  '1': 'AdminListDevicesRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminListDevicesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListDevicesRequestDescriptor = $convert.base64Decode(
    'ChdBZG1pbkxpc3REZXZpY2VzUmVxdWVzdBIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0EiwKD3'
    'BhZ2luYXRpb250b2tlbhj/isCtASABKAlSD3BhZ2luYXRpb250b2tlbhIiCgp1c2VycG9vbGlk'
    'GP7Gi50BIAEoCVIKdXNlcnBvb2xpZBIeCgh1c2VybmFtZRjaqaPgASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use adminListDevicesResponseDescriptor instead')
const AdminListDevicesResponse$json = {
  '1': 'AdminListDevicesResponse',
  '2': [
    {
      '1': 'devices',
      '3': 128774945,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.DeviceType',
      '10': 'devices'
    },
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
  ],
};

/// Descriptor for `AdminListDevicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListDevicesResponseDescriptor = $convert.base64Decode(
    'ChhBZG1pbkxpc3REZXZpY2VzUmVzcG9uc2USNAoHZGV2aWNlcxih5rM9IAMoCzIXLmNvZ25pdG'
    '9faWRwLkRldmljZVR5cGVSB2RldmljZXMSLAoPcGFnaW5hdGlvbnRva2VuGP+KwK0BIAEoCVIP'
    'cGFnaW5hdGlvbnRva2Vu');

@$core.Deprecated('Use adminListGroupsForUserRequestDescriptor instead')
const AdminListGroupsForUserRequest$json = {
  '1': 'AdminListGroupsForUserRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminListGroupsForUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListGroupsForUserRequestDescriptor =
    $convert.base64Decode(
        'Ch1BZG1pbkxpc3RHcm91cHNGb3JVc2VyUmVxdWVzdBIYCgVsaW1pdBjVldnEASABKAVSBWxpbW'
        'l0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiIKCnVzZXJwb29saWQY/saLnQEg'
        'ASgJUgp1c2VycG9vbGlkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminListGroupsForUserResponseDescriptor instead')
const AdminListGroupsForUserResponse$json = {
  '1': 'AdminListGroupsForUserResponse',
  '2': [
    {
      '1': 'groups',
      '3': 360662414,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.GroupType',
      '10': 'groups'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `AdminListGroupsForUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListGroupsForUserResponseDescriptor =
    $convert.base64Decode(
        'Ch5BZG1pbkxpc3RHcm91cHNGb3JVc2VyUmVzcG9uc2USMgoGZ3JvdXBzGI6L/asBIAMoCzIWLm'
        'NvZ25pdG9faWRwLkdyb3VwVHlwZVIGZ3JvdXBzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4'
        'dHRva2Vu');

@$core.Deprecated('Use adminListUserAuthEventsRequestDescriptor instead')
const AdminListUserAuthEventsRequest$json = {
  '1': 'AdminListUserAuthEventsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminListUserAuthEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListUserAuthEventsRequestDescriptor =
    $convert.base64Decode(
        'Ch5BZG1pbkxpc3RVc2VyQXV0aEV2ZW50c1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKA'
        'VSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIgoKdXNlcnBv'
        'b2xpZBj+xoudASABKAlSCnVzZXJwb29saWQSHgoIdXNlcm5hbWUY2qmj4AEgASgJUgh1c2Vybm'
        'FtZQ==');

@$core.Deprecated('Use adminListUserAuthEventsResponseDescriptor instead')
const AdminListUserAuthEventsResponse$json = {
  '1': 'AdminListUserAuthEventsResponse',
  '2': [
    {
      '1': 'authevents',
      '3': 470980699,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AuthEventType',
      '10': 'authevents'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `AdminListUserAuthEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListUserAuthEventsResponseDescriptor =
    $convert.base64Decode(
        'Ch9BZG1pbkxpc3RVc2VyQXV0aEV2ZW50c1Jlc3BvbnNlEj4KCmF1dGhldmVudHMY27DK4AEgAy'
        'gLMhouY29nbml0b19pZHAuQXV0aEV2ZW50VHlwZVIKYXV0aGV2ZW50cxIfCgluZXh0dG9rZW4Y'
        '/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use adminRemoveUserFromGroupRequestDescriptor instead')
const AdminRemoveUserFromGroupRequest$json = {
  '1': 'AdminRemoveUserFromGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminRemoveUserFromGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminRemoveUserFromGroupRequestDescriptor =
    $convert.base64Decode(
        'Ch9BZG1pblJlbW92ZVVzZXJGcm9tR3JvdXBSZXF1ZXN0EiAKCWdyb3VwbmFtZRjIyqCqASABKA'
        'lSCWdyb3VwbmFtZRIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZBIeCgh1c2Vy'
        'bmFtZRjaqaPgASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use adminResetUserPasswordRequestDescriptor instead')
const AdminResetUserPasswordRequest$json = {
  '1': 'AdminResetUserPasswordRequest',
  '2': [
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AdminResetUserPasswordRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [AdminResetUserPasswordRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use adminResetUserPasswordRequestDescriptor instead')
const AdminResetUserPasswordRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminResetUserPasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminResetUserPasswordRequestDescriptor = $convert.base64Decode(
    'Ch1BZG1pblJlc2V0VXNlclBhc3N3b3JkUmVxdWVzdBJpCg5jbGllbnRtZXRhZGF0YRjMr/9hIA'
    'MoCzI+LmNvZ25pdG9faWRwLkFkbWluUmVzZXRVc2VyUGFzc3dvcmRSZXF1ZXN0LkNsaWVudG1l'
    'dGFkYXRhRW50cnlSDmNsaWVudG1ldGFkYXRhEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2'
    'VycG9vbGlkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWUaQQoTQ2xpZW50bWV0YWRh'
    'dGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use adminResetUserPasswordResponseDescriptor instead')
const AdminResetUserPasswordResponse$json = {
  '1': 'AdminResetUserPasswordResponse',
};

/// Descriptor for `AdminResetUserPasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminResetUserPasswordResponseDescriptor =
    $convert.base64Decode('Ch5BZG1pblJlc2V0VXNlclBhc3N3b3JkUmVzcG9uc2U=');

@$core.Deprecated('Use adminRespondToAuthChallengeRequestDescriptor instead')
const AdminRespondToAuthChallengeRequest$json = {
  '1': 'AdminRespondToAuthChallengeRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {
      '1': 'challengename',
      '3': 170761310,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'challengename'
    },
    {
      '1': 'challengeresponses',
      '3': 60230615,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_idp.AdminRespondToAuthChallengeRequest.ChallengeresponsesEntry',
      '10': 'challengeresponses'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_idp.AdminRespondToAuthChallengeRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {
      '1': 'contextdata',
      '3': 235700261,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ContextDataType',
      '10': 'contextdata'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [
    AdminRespondToAuthChallengeRequest_ChallengeresponsesEntry$json,
    AdminRespondToAuthChallengeRequest_ClientmetadataEntry$json
  ],
};

@$core.Deprecated('Use adminRespondToAuthChallengeRequestDescriptor instead')
const AdminRespondToAuthChallengeRequest_ChallengeresponsesEntry$json = {
  '1': 'ChallengeresponsesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use adminRespondToAuthChallengeRequestDescriptor instead')
const AdminRespondToAuthChallengeRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminRespondToAuthChallengeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminRespondToAuthChallengeRequestDescriptor = $convert.base64Decode(
    'CiJBZG1pblJlc3BvbmRUb0F1dGhDaGFsbGVuZ2VSZXF1ZXN0ElMKEWFuYWx5dGljc21ldGFkYX'
    'RhGPeFwTwgASgLMiIuY29nbml0b19pZHAuQW5hbHl0aWNzTWV0YWRhdGFUeXBlUhFhbmFseXRp'
    'Y3NtZXRhZGF0YRJHCg1jaGFsbGVuZ2VuYW1lGN64tlEgASgOMh4uY29nbml0b19pZHAuQ2hhbG'
    'xlbmdlTmFtZVR5cGVSDWNoYWxsZW5nZW5hbWUSegoSY2hhbGxlbmdlcmVzcG9uc2VzGNeX3Bwg'
    'AygLMkcuY29nbml0b19pZHAuQWRtaW5SZXNwb25kVG9BdXRoQ2hhbGxlbmdlUmVxdWVzdC5DaG'
    'FsbGVuZ2VyZXNwb25zZXNFbnRyeVISY2hhbGxlbmdlcmVzcG9uc2VzEh4KCGNsaWVudGlkGKTo'
    'htYBIAEoCVIIY2xpZW50aWQSbgoOY2xpZW50bWV0YWRhdGEYzK//YSADKAsyQy5jb2duaXRvX2'
    'lkcC5BZG1pblJlc3BvbmRUb0F1dGhDaGFsbGVuZ2VSZXF1ZXN0LkNsaWVudG1ldGFkYXRhRW50'
    'cnlSDmNsaWVudG1ldGFkYXRhEkEKC2NvbnRleHRkYXRhGKWAsnAgASgLMhwuY29nbml0b19pZH'
    'AuQ29udGV4dERhdGFUeXBlUgtjb250ZXh0ZGF0YRIbCgdzZXNzaW9uGJiZowIgASgJUgdzZXNz'
    'aW9uEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlkGkUKF0NoYWxsZW5nZXJlc3'
    'BvbnNlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEa'
    'QQoTQ2xpZW50bWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCV'
    'IFdmFsdWU6AjgB');

@$core.Deprecated('Use adminRespondToAuthChallengeResponseDescriptor instead')
const AdminRespondToAuthChallengeResponse$json = {
  '1': 'AdminRespondToAuthChallengeResponse',
  '2': [
    {
      '1': 'authenticationresult',
      '3': 519327313,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AuthenticationResultType',
      '10': 'authenticationresult'
    },
    {
      '1': 'challengename',
      '3': 170761310,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'challengename'
    },
    {
      '1': 'challengeparameters',
      '3': 79757023,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_idp.AdminRespondToAuthChallengeResponse.ChallengeparametersEntry',
      '10': 'challengeparameters'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
  ],
  '3': [AdminRespondToAuthChallengeResponse_ChallengeparametersEntry$json],
};

@$core.Deprecated('Use adminRespondToAuthChallengeResponseDescriptor instead')
const AdminRespondToAuthChallengeResponse_ChallengeparametersEntry$json = {
  '1': 'ChallengeparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminRespondToAuthChallengeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminRespondToAuthChallengeResponseDescriptor = $convert.base64Decode(
    'CiNBZG1pblJlc3BvbmRUb0F1dGhDaGFsbGVuZ2VSZXNwb25zZRJdChRhdXRoZW50aWNhdGlvbn'
    'Jlc3VsdBjRnNH3ASABKAsyJS5jb2duaXRvX2lkcC5BdXRoZW50aWNhdGlvblJlc3VsdFR5cGVS'
    'FGF1dGhlbnRpY2F0aW9ucmVzdWx0EkcKDWNoYWxsZW5nZW5hbWUY3ri2USABKA4yHi5jb2duaX'
    'RvX2lkcC5DaGFsbGVuZ2VOYW1lVHlwZVINY2hhbGxlbmdlbmFtZRJ+ChNjaGFsbGVuZ2VwYXJh'
    'bWV0ZXJzGN/9gyYgAygLMkkuY29nbml0b19pZHAuQWRtaW5SZXNwb25kVG9BdXRoQ2hhbGxlbm'
    'dlUmVzcG9uc2UuQ2hhbGxlbmdlcGFyYW1ldGVyc0VudHJ5UhNjaGFsbGVuZ2VwYXJhbWV0ZXJz'
    'EhsKB3Nlc3Npb24YmJmjAiABKAlSB3Nlc3Npb24aRgoYQ2hhbGxlbmdlcGFyYW1ldGVyc0VudH'
    'J5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use adminSetUserMFAPreferenceRequestDescriptor instead')
const AdminSetUserMFAPreferenceRequest$json = {
  '1': 'AdminSetUserMFAPreferenceRequest',
  '2': [
    {
      '1': 'emailmfasettings',
      '3': 306207991,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailMfaSettingsType',
      '10': 'emailmfasettings'
    },
    {
      '1': 'smsmfasettings',
      '3': 220659382,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SMSMfaSettingsType',
      '10': 'smsmfasettings'
    },
    {
      '1': 'softwaretokenmfasettings',
      '3': 22910957,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SoftwareTokenMfaSettingsType',
      '10': 'softwaretokenmfasettings'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminSetUserMFAPreferenceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminSetUserMFAPreferenceRequestDescriptor = $convert.base64Decode(
    'CiBBZG1pblNldFVzZXJNRkFQcmVmZXJlbmNlUmVxdWVzdBJRChBlbWFpbG1mYXNldHRpbmdzGP'
    'e5gZIBIAEoCzIhLmNvZ25pdG9faWRwLkVtYWlsTWZhU2V0dGluZ3NUeXBlUhBlbWFpbG1mYXNl'
    'dHRpbmdzEkoKDnNtc21mYXNldHRpbmdzGLb9m2kgASgLMh8uY29nbml0b19pZHAuU01TTWZhU2'
    'V0dGluZ3NUeXBlUg5zbXNtZmFzZXR0aW5ncxJoChhzb2Z0d2FyZXRva2VubWZhc2V0dGluZ3MY'
    '7a/2CiABKAsyKS5jb2duaXRvX2lkcC5Tb2Z0d2FyZVRva2VuTWZhU2V0dGluZ3NUeXBlUhhzb2'
    'Z0d2FyZXRva2VubWZhc2V0dGluZ3MSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29s'
    'aWQSHgoIdXNlcm5hbWUY2qmj4AEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use adminSetUserMFAPreferenceResponseDescriptor instead')
const AdminSetUserMFAPreferenceResponse$json = {
  '1': 'AdminSetUserMFAPreferenceResponse',
};

/// Descriptor for `AdminSetUserMFAPreferenceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminSetUserMFAPreferenceResponseDescriptor =
    $convert.base64Decode('CiFBZG1pblNldFVzZXJNRkFQcmVmZXJlbmNlUmVzcG9uc2U=');

@$core.Deprecated('Use adminSetUserPasswordRequestDescriptor instead')
const AdminSetUserPasswordRequest$json = {
  '1': 'AdminSetUserPasswordRequest',
  '2': [
    {'1': 'password', '3': 214108217, '4': 1, '5': 9, '10': 'password'},
    {'1': 'permanent', '3': 141702118, '4': 1, '5': 8, '10': 'permanent'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminSetUserPasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminSetUserPasswordRequestDescriptor =
    $convert.base64Decode(
        'ChtBZG1pblNldFVzZXJQYXNzd29yZFJlcXVlc3QSHQoIcGFzc3dvcmQYuZCMZiABKAlSCHBhc3'
        'N3b3JkEh8KCXBlcm1hbmVudBjm58hDIAEoCFIJcGVybWFuZW50EiIKCnVzZXJwb29saWQY/saL'
        'nQEgASgJUgp1c2VycG9vbGlkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminSetUserPasswordResponseDescriptor instead')
const AdminSetUserPasswordResponse$json = {
  '1': 'AdminSetUserPasswordResponse',
};

/// Descriptor for `AdminSetUserPasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminSetUserPasswordResponseDescriptor =
    $convert.base64Decode('ChxBZG1pblNldFVzZXJQYXNzd29yZFJlc3BvbnNl');

@$core.Deprecated('Use adminSetUserSettingsRequestDescriptor instead')
const AdminSetUserSettingsRequest$json = {
  '1': 'AdminSetUserSettingsRequest',
  '2': [
    {
      '1': 'mfaoptions',
      '3': 501540826,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.MFAOptionType',
      '10': 'mfaoptions'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminSetUserSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminSetUserSettingsRequestDescriptor =
    $convert.base64Decode(
        'ChtBZG1pblNldFVzZXJTZXR0aW5nc1JlcXVlc3QSPgoKbWZhb3B0aW9ucxjaz5PvASADKAsyGi'
        '5jb2duaXRvX2lkcC5NRkFPcHRpb25UeXBlUgptZmFvcHRpb25zEiIKCnVzZXJwb29saWQY/saL'
        'nQEgASgJUgp1c2VycG9vbGlkEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use adminSetUserSettingsResponseDescriptor instead')
const AdminSetUserSettingsResponse$json = {
  '1': 'AdminSetUserSettingsResponse',
};

/// Descriptor for `AdminSetUserSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminSetUserSettingsResponseDescriptor =
    $convert.base64Decode('ChxBZG1pblNldFVzZXJTZXR0aW5nc1Jlc3BvbnNl');

@$core.Deprecated('Use adminUpdateAuthEventFeedbackRequestDescriptor instead')
const AdminUpdateAuthEventFeedbackRequest$json = {
  '1': 'AdminUpdateAuthEventFeedbackRequest',
  '2': [
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 9, '10': 'eventid'},
    {
      '1': 'feedbackvalue',
      '3': 259452876,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.FeedbackValueType',
      '10': 'feedbackvalue'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminUpdateAuthEventFeedbackRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUpdateAuthEventFeedbackRequestDescriptor =
    $convert.base64Decode(
        'CiNBZG1pblVwZGF0ZUF1dGhFdmVudEZlZWRiYWNrUmVxdWVzdBIcCgdldmVudGlkGNOW3bMBIA'
        'EoCVIHZXZlbnRpZBJHCg1mZWVkYmFja3ZhbHVlGMzf23sgASgOMh4uY29nbml0b19pZHAuRmVl'
        'ZGJhY2tWYWx1ZVR5cGVSDWZlZWRiYWNrdmFsdWUSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCn'
        'VzZXJwb29saWQSHgoIdXNlcm5hbWUY2qmj4AEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use adminUpdateAuthEventFeedbackResponseDescriptor instead')
const AdminUpdateAuthEventFeedbackResponse$json = {
  '1': 'AdminUpdateAuthEventFeedbackResponse',
};

/// Descriptor for `AdminUpdateAuthEventFeedbackResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUpdateAuthEventFeedbackResponseDescriptor =
    $convert
        .base64Decode('CiRBZG1pblVwZGF0ZUF1dGhFdmVudEZlZWRiYWNrUmVzcG9uc2U=');

@$core.Deprecated('Use adminUpdateDeviceStatusRequestDescriptor instead')
const AdminUpdateDeviceStatusRequest$json = {
  '1': 'AdminUpdateDeviceStatusRequest',
  '2': [
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
    {
      '1': 'devicerememberedstatus',
      '3': 111455992,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DeviceRememberedStatusType',
      '10': 'devicerememberedstatus'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminUpdateDeviceStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUpdateDeviceStatusRequestDescriptor = $convert.base64Decode(
    'Ch5BZG1pblVwZGF0ZURldmljZVN0YXR1c1JlcXVlc3QSIAoJZGV2aWNla2V5GKvkyLYBIAEoCV'
    'IJZGV2aWNla2V5EmIKFmRldmljZXJlbWVtYmVyZWRzdGF0dXMY+N2SNSABKA4yJy5jb2duaXRv'
    'X2lkcC5EZXZpY2VSZW1lbWJlcmVkU3RhdHVzVHlwZVIWZGV2aWNlcmVtZW1iZXJlZHN0YXR1cx'
    'IiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZBIeCgh1c2VybmFtZRjaqaPgASAB'
    'KAlSCHVzZXJuYW1l');

@$core.Deprecated('Use adminUpdateDeviceStatusResponseDescriptor instead')
const AdminUpdateDeviceStatusResponse$json = {
  '1': 'AdminUpdateDeviceStatusResponse',
};

/// Descriptor for `AdminUpdateDeviceStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUpdateDeviceStatusResponseDescriptor =
    $convert.base64Decode('Ch9BZG1pblVwZGF0ZURldmljZVN0YXR1c1Jlc3BvbnNl');

@$core.Deprecated('Use adminUpdateUserAttributesRequestDescriptor instead')
const AdminUpdateUserAttributesRequest$json = {
  '1': 'AdminUpdateUserAttributesRequest',
  '2': [
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AdminUpdateUserAttributesRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {
      '1': 'userattributes',
      '3': 194667064,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'userattributes'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [AdminUpdateUserAttributesRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use adminUpdateUserAttributesRequestDescriptor instead')
const AdminUpdateUserAttributesRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AdminUpdateUserAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUpdateUserAttributesRequestDescriptor = $convert.base64Decode(
    'CiBBZG1pblVwZGF0ZVVzZXJBdHRyaWJ1dGVzUmVxdWVzdBJsCg5jbGllbnRtZXRhZGF0YRjMr/'
    '9hIAMoCzJBLmNvZ25pdG9faWRwLkFkbWluVXBkYXRlVXNlckF0dHJpYnV0ZXNSZXF1ZXN0LkNs'
    'aWVudG1ldGFkYXRhRW50cnlSDmNsaWVudG1ldGFkYXRhEkUKDnVzZXJhdHRyaWJ1dGVzGLjE6V'
    'wgAygLMhouY29nbml0b19pZHAuQXR0cmlidXRlVHlwZVIOdXNlcmF0dHJpYnV0ZXMSIgoKdXNl'
    'cnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQSHgoIdXNlcm5hbWUY2qmj4AEgASgJUgh1c2'
    'VybmFtZRpBChNDbGllbnRtZXRhZGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVl'
    'GAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use adminUpdateUserAttributesResponseDescriptor instead')
const AdminUpdateUserAttributesResponse$json = {
  '1': 'AdminUpdateUserAttributesResponse',
};

/// Descriptor for `AdminUpdateUserAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUpdateUserAttributesResponseDescriptor =
    $convert.base64Decode('CiFBZG1pblVwZGF0ZVVzZXJBdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated('Use adminUserGlobalSignOutRequestDescriptor instead')
const AdminUserGlobalSignOutRequest$json = {
  '1': 'AdminUserGlobalSignOutRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `AdminUserGlobalSignOutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUserGlobalSignOutRequestDescriptor =
    $convert.base64Decode(
        'Ch1BZG1pblVzZXJHbG9iYWxTaWduT3V0UmVxdWVzdBIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCV'
        'IKdXNlcnBvb2xpZBIeCgh1c2VybmFtZRjaqaPgASABKAlSCHVzZXJuYW1l');

@$core.Deprecated('Use adminUserGlobalSignOutResponseDescriptor instead')
const AdminUserGlobalSignOutResponse$json = {
  '1': 'AdminUserGlobalSignOutResponse',
};

/// Descriptor for `AdminUserGlobalSignOutResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUserGlobalSignOutResponseDescriptor =
    $convert.base64Decode('Ch5BZG1pblVzZXJHbG9iYWxTaWduT3V0UmVzcG9uc2U=');

@$core.Deprecated('Use advancedSecurityAdditionalFlowsTypeDescriptor instead')
const AdvancedSecurityAdditionalFlowsType$json = {
  '1': 'AdvancedSecurityAdditionalFlowsType',
  '2': [
    {
      '1': 'customauthmode',
      '3': 236623048,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AdvancedSecurityEnabledModeType',
      '10': 'customauthmode'
    },
  ],
};

/// Descriptor for `AdvancedSecurityAdditionalFlowsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List advancedSecurityAdditionalFlowsTypeDescriptor =
    $convert.base64Decode(
        'CiNBZHZhbmNlZFNlY3VyaXR5QWRkaXRpb25hbEZsb3dzVHlwZRJXCg5jdXN0b21hdXRobW9kZR'
        'jIqepwIAEoDjIsLmNvZ25pdG9faWRwLkFkdmFuY2VkU2VjdXJpdHlFbmFibGVkTW9kZVR5cGVS'
        'DmN1c3RvbWF1dGhtb2Rl');

@$core.Deprecated('Use aliasExistsExceptionDescriptor instead')
const AliasExistsException$json = {
  '1': 'AliasExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AliasExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aliasExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChRBbGlhc0V4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use analyticsConfigurationTypeDescriptor instead')
const AnalyticsConfigurationType$json = {
  '1': 'AnalyticsConfigurationType',
  '2': [
    {
      '1': 'applicationarn',
      '3': 524185341,
      '4': 1,
      '5': 9,
      '10': 'applicationarn'
    },
    {
      '1': 'applicationid',
      '3': 132869857,
      '4': 1,
      '5': 9,
      '10': 'applicationid'
    },
    {'1': 'externalid', '3': 271401992, '4': 1, '5': 9, '10': 'externalid'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'userdatashared',
      '3': 370377352,
      '4': 1,
      '5': 8,
      '10': 'userdatashared'
    },
  ],
};

/// Descriptor for `AnalyticsConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyticsConfigurationTypeDescriptor = $convert.base64Decode(
    'ChpBbmFseXRpY3NDb25maWd1cmF0aW9uVHlwZRIqCg5hcHBsaWNhdGlvbmFybhj93fn5ASABKA'
    'lSDmFwcGxpY2F0aW9uYXJuEicKDWFwcGxpY2F0aW9uaWQY4d2tPyABKAlSDWFwcGxpY2F0aW9u'
    'aWQSIgoKZXh0ZXJuYWxpZBiIiLWBASABKAlSCmV4dGVybmFsaWQSHAoHcm9sZWFybhiB+OeZAS'
    'ABKAlSB3JvbGVhcm4SKgoOdXNlcmRhdGFzaGFyZWQYiIXOsAEgASgIUg51c2VyZGF0YXNoYXJl'
    'ZA==');

@$core.Deprecated('Use analyticsMetadataTypeDescriptor instead')
const AnalyticsMetadataType$json = {
  '1': 'AnalyticsMetadataType',
  '2': [
    {
      '1': 'analyticsendpointid',
      '3': 341236308,
      '4': 1,
      '5': 9,
      '10': 'analyticsendpointid'
    },
  ],
};

/// Descriptor for `AnalyticsMetadataType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyticsMetadataTypeDescriptor = $convert.base64Decode(
    'ChVBbmFseXRpY3NNZXRhZGF0YVR5cGUSNAoTYW5hbHl0aWNzZW5kcG9pbnRpZBjUtNuiASABKA'
    'lSE2FuYWx5dGljc2VuZHBvaW50aWQ=');

@$core.Deprecated('Use assetTypeDescriptor instead')
const AssetType$json = {
  '1': 'AssetType',
  '2': [
    {'1': 'bytes', '3': 397169317, '4': 1, '5': 12, '10': 'bytes'},
    {
      '1': 'category',
      '3': 263447954,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AssetCategoryType',
      '10': 'category'
    },
    {
      '1': 'colormode',
      '3': 183765676,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ColorSchemeModeType',
      '10': 'colormode'
    },
    {
      '1': 'extension',
      '3': 381907515,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AssetExtensionType',
      '10': 'extension'
    },
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
  ],
};

/// Descriptor for `AssetType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assetTypeDescriptor = $convert.base64Decode(
    'CglBc3NldFR5cGUSGAoFYnl0ZXMYpaWxvQEgASgMUgVieXRlcxI9CghjYXRlZ29yeRiSy899IA'
    'EoDjIeLmNvZ25pdG9faWRwLkFzc2V0Q2F0ZWdvcnlUeXBlUghjYXRlZ29yeRJBCgljb2xvcm1v'
    'ZGUYrJXQVyABKA4yIC5jb2duaXRvX2lkcC5Db2xvclNjaGVtZU1vZGVUeXBlUgljb2xvcm1vZG'
    'USQQoJZXh0ZW5zaW9uGLvkjbYBIAEoDjIfLmNvZ25pdG9faWRwLkFzc2V0RXh0ZW5zaW9uVHlw'
    'ZVIJZXh0ZW5zaW9uEiIKCnJlc291cmNlaWQYkbrx+gEgASgJUgpyZXNvdXJjZWlk');

@$core.Deprecated('Use associateSoftwareTokenRequestDescriptor instead')
const AssociateSoftwareTokenRequest$json = {
  '1': 'AssociateSoftwareTokenRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
  ],
};

/// Descriptor for `AssociateSoftwareTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateSoftwareTokenRequestDescriptor =
    $convert.base64Decode(
        'Ch1Bc3NvY2lhdGVTb2Z0d2FyZVRva2VuUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIAEoCV'
        'ILYWNjZXNzdG9rZW4SGwoHc2Vzc2lvbhiYmaMCIAEoCVIHc2Vzc2lvbg==');

@$core.Deprecated('Use associateSoftwareTokenResponseDescriptor instead')
const AssociateSoftwareTokenResponse$json = {
  '1': 'AssociateSoftwareTokenResponse',
  '2': [
    {'1': 'secretcode', '3': 252366169, '4': 1, '5': 9, '10': 'secretcode'},
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
  ],
};

/// Descriptor for `AssociateSoftwareTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateSoftwareTokenResponseDescriptor =
    $convert.base64Decode(
        'Ch5Bc3NvY2lhdGVTb2Z0d2FyZVRva2VuUmVzcG9uc2USIQoKc2VjcmV0Y29kZRjZmqt4IAEoCV'
        'IKc2VjcmV0Y29kZRIbCgdzZXNzaW9uGJiZowIgASgJUgdzZXNzaW9u');

@$core.Deprecated('Use attributeTypeDescriptor instead')
const AttributeType$json = {
  '1': 'AttributeType',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `AttributeType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attributeTypeDescriptor = $convert.base64Decode(
    'Cg1BdHRyaWJ1dGVUeXBlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGAoFdmFsdWUY6/KfigEgAS'
    'gJUgV2YWx1ZQ==');

@$core.Deprecated('Use authEventTypeDescriptor instead')
const AuthEventType$json = {
  '1': 'AuthEventType',
  '2': [
    {
      '1': 'challengeresponses',
      '3': 60230615,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ChallengeResponseType',
      '10': 'challengeresponses'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'eventcontextdata',
      '3': 432307527,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EventContextDataType',
      '10': 'eventcontextdata'
    },
    {
      '1': 'eventfeedback',
      '3': 156992389,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EventFeedbackType',
      '10': 'eventfeedback'
    },
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 9, '10': 'eventid'},
    {
      '1': 'eventresponse',
      '3': 198651437,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.EventResponseType',
      '10': 'eventresponse'
    },
    {
      '1': 'eventrisk',
      '3': 8230945,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EventRiskType',
      '10': 'eventrisk'
    },
    {
      '1': 'eventtype',
      '3': 468897896,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.EventType',
      '10': 'eventtype'
    },
  ],
};

/// Descriptor for `AuthEventType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authEventTypeDescriptor = $convert.base64Decode(
    'Cg1BdXRoRXZlbnRUeXBlElUKEmNoYWxsZW5nZXJlc3BvbnNlcxjXl9wcIAMoCzIiLmNvZ25pdG'
    '9faWRwLkNoYWxsZW5nZVJlc3BvbnNlVHlwZVISY2hhbGxlbmdlcmVzcG9uc2VzEiYKDGNyZWF0'
    'aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9uZGF0ZRJRChBldmVudGNvbnRleHRkYXRhGMf6kc'
    '4BIAEoCzIhLmNvZ25pdG9faWRwLkV2ZW50Q29udGV4dERhdGFUeXBlUhBldmVudGNvbnRleHRk'
    'YXRhEkcKDWV2ZW50ZmVlZGJhY2sYhYfuSiABKAsyHi5jb2duaXRvX2lkcC5FdmVudEZlZWRiYW'
    'NrVHlwZVINZXZlbnRmZWVkYmFjaxIcCgdldmVudGlkGNOW3bMBIAEoCVIHZXZlbnRpZBJHCg1l'
    'dmVudHJlc3BvbnNlGK3c3F4gASgOMh4uY29nbml0b19pZHAuRXZlbnRSZXNwb25zZVR5cGVSDW'
    'V2ZW50cmVzcG9uc2USOwoJZXZlbnRyaXNrGKGw9gMgASgLMhouY29nbml0b19pZHAuRXZlbnRS'
    'aXNrVHlwZVIJZXZlbnRyaXNrEjgKCWV2ZW50dHlwZRjooMvfASABKA4yFi5jb2duaXRvX2lkcC'
    '5FdmVudFR5cGVSCWV2ZW50dHlwZQ==');

@$core.Deprecated('Use authenticationResultTypeDescriptor instead')
const AuthenticationResultType$json = {
  '1': 'AuthenticationResultType',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'expiresin', '3': 69497873, '4': 1, '5': 5, '10': 'expiresin'},
    {'1': 'idtoken', '3': 228470, '4': 1, '5': 9, '10': 'idtoken'},
    {
      '1': 'newdevicemetadata',
      '3': 383914787,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.NewDeviceMetadataType',
      '10': 'newdevicemetadata'
    },
    {'1': 'refreshtoken', '3': 253777778, '4': 1, '5': 9, '10': 'refreshtoken'},
    {'1': 'tokentype', '3': 231695523, '4': 1, '5': 9, '10': 'tokentype'},
  ],
};

/// Descriptor for `AuthenticationResultType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authenticationResultTypeDescriptor = $convert.base64Decode(
    'ChhBdXRoZW50aWNhdGlvblJlc3VsdFR5cGUSIwoLYWNjZXNzdG9rZW4YibyQRiABKAlSC2FjY2'
    'Vzc3Rva2VuEh8KCWV4cGlyZXNpbhiR6JEhIAEoBVIJZXhwaXJlc2luEhoKB2lkdG9rZW4Y9vgN'
    'IAEoCVIHaWR0b2tlbhJUChFuZXdkZXZpY2VtZXRhZGF0YRijpoi3ASABKAsyIi5jb2duaXRvX2'
    'lkcC5OZXdEZXZpY2VNZXRhZGF0YVR5cGVSEW5ld2RldmljZW1ldGFkYXRhEiUKDHJlZnJlc2h0'
    'b2tlbhjyroF5IAEoCVIMcmVmcmVzaHRva2VuEh8KCXRva2VudHlwZRijyb1uIAEoCVIJdG9rZW'
    '50eXBl');

@$core.Deprecated('Use challengeResponseTypeDescriptor instead')
const ChallengeResponseType$json = {
  '1': 'ChallengeResponseType',
  '2': [
    {
      '1': 'challengename',
      '3': 170761310,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeName',
      '10': 'challengename'
    },
    {
      '1': 'challengeresponse',
      '3': 268501730,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeResponse',
      '10': 'challengeresponse'
    },
  ],
};

/// Descriptor for `ChallengeResponseType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List challengeResponseTypeDescriptor = $convert.base64Decode(
    'ChVDaGFsbGVuZ2VSZXNwb25zZVR5cGUSQwoNY2hhbGxlbmdlbmFtZRjeuLZRIAEoDjIaLmNvZ2'
    '5pdG9faWRwLkNoYWxsZW5nZU5hbWVSDWNoYWxsZW5nZW5hbWUSUAoRY2hhbGxlbmdlcmVzcG9u'
    'c2UY4oWEgAEgASgOMh4uY29nbml0b19pZHAuQ2hhbGxlbmdlUmVzcG9uc2VSEWNoYWxsZW5nZX'
    'Jlc3BvbnNl');

@$core.Deprecated('Use changePasswordRequestDescriptor instead')
const ChangePasswordRequest$json = {
  '1': 'ChangePasswordRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'previouspassword',
      '3': 303785026,
      '4': 1,
      '5': 9,
      '10': 'previouspassword'
    },
    {
      '1': 'proposedpassword',
      '3': 234890803,
      '4': 1,
      '5': 9,
      '10': 'proposedpassword'
    },
  ],
};

/// Descriptor for `ChangePasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changePasswordRequestDescriptor = $convert.base64Decode(
    'ChVDaGFuZ2VQYXNzd29yZFJlcXVlc3QSIwoLYWNjZXNzdG9rZW4YibyQRiABKAlSC2FjY2Vzc3'
    'Rva2VuEi4KEHByZXZpb3VzcGFzc3dvcmQYwsjtkAEgASgJUhBwcmV2aW91c3Bhc3N3b3JkEi0K'
    'EHByb3Bvc2VkcGFzc3dvcmQYs8yAcCABKAlSEHByb3Bvc2VkcGFzc3dvcmQ=');

@$core.Deprecated('Use changePasswordResponseDescriptor instead')
const ChangePasswordResponse$json = {
  '1': 'ChangePasswordResponse',
};

/// Descriptor for `ChangePasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changePasswordResponseDescriptor =
    $convert.base64Decode('ChZDaGFuZ2VQYXNzd29yZFJlc3BvbnNl');

@$core.Deprecated('Use clientSecretDescriptorTypeDescriptor instead')
const ClientSecretDescriptorType$json = {
  '1': 'ClientSecretDescriptorType',
  '2': [
    {
      '1': 'clientsecretcreatedate',
      '3': 345458871,
      '4': 1,
      '5': 9,
      '10': 'clientsecretcreatedate'
    },
    {
      '1': 'clientsecretid',
      '3': 51685996,
      '4': 1,
      '5': 9,
      '10': 'clientsecretid'
    },
    {
      '1': 'clientsecretvalue',
      '3': 475370712,
      '4': 1,
      '5': 9,
      '10': 'clientsecretvalue'
    },
  ],
};

/// Descriptor for `ClientSecretDescriptorType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clientSecretDescriptorTypeDescriptor = $convert.base64Decode(
    'ChpDbGllbnRTZWNyZXREZXNjcmlwdG9yVHlwZRI6ChZjbGllbnRzZWNyZXRjcmVhdGVkYXRlGL'
    'eR3aQBIAEoCVIWY2xpZW50c2VjcmV0Y3JlYXRlZGF0ZRIpCg5jbGllbnRzZWNyZXRpZBjs1NIY'
    'IAEoCVIOY2xpZW50c2VjcmV0aWQSMAoRY2xpZW50c2VjcmV0dmFsdWUY2KnW4gEgASgJUhFjbG'
    'llbnRzZWNyZXR2YWx1ZQ==');

@$core.Deprecated('Use cloudWatchLogsConfigurationTypeDescriptor instead')
const CloudWatchLogsConfigurationType$json = {
  '1': 'CloudWatchLogsConfigurationType',
  '2': [
    {'1': 'loggrouparn', '3': 220362512, '4': 1, '5': 9, '10': 'loggrouparn'},
  ],
};

/// Descriptor for `CloudWatchLogsConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchLogsConfigurationTypeDescriptor =
    $convert.base64Decode(
        'Ch9DbG91ZFdhdGNoTG9nc0NvbmZpZ3VyYXRpb25UeXBlEiMKC2xvZ2dyb3VwYXJuGJDuiWkgAS'
        'gJUgtsb2dncm91cGFybg==');

@$core.Deprecated('Use codeDeliveryDetailsTypeDescriptor instead')
const CodeDeliveryDetailsType$json = {
  '1': 'CodeDeliveryDetailsType',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'deliverymedium',
      '3': 140498471,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DeliveryMediumType',
      '10': 'deliverymedium'
    },
    {'1': 'destination', '3': 457443680, '4': 1, '5': 9, '10': 'destination'},
  ],
};

/// Descriptor for `CodeDeliveryDetailsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeDeliveryDetailsTypeDescriptor = $convert.base64Decode(
    'ChdDb2RlRGVsaXZlcnlEZXRhaWxzVHlwZRIoCg1hdHRyaWJ1dGVuYW1lGK2VmKgBIAEoCVINYX'
    'R0cmlidXRlbmFtZRJKCg5kZWxpdmVyeW1lZGl1bRinrP9CIAEoDjIfLmNvZ25pdG9faWRwLkRl'
    'bGl2ZXJ5TWVkaXVtVHlwZVIOZGVsaXZlcnltZWRpdW0SJAoLZGVzdGluYXRpb24Y4JKQ2gEgAS'
    'gJUgtkZXN0aW5hdGlvbg==');

@$core.Deprecated('Use codeDeliveryFailureExceptionDescriptor instead')
const CodeDeliveryFailureException$json = {
  '1': 'CodeDeliveryFailureException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CodeDeliveryFailureException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeDeliveryFailureExceptionDescriptor =
    $convert.base64Decode(
        'ChxDb2RlRGVsaXZlcnlGYWlsdXJlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use codeMismatchExceptionDescriptor instead')
const CodeMismatchException$json = {
  '1': 'CodeMismatchException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CodeMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeMismatchExceptionDescriptor = $convert.base64Decode(
    'ChVDb2RlTWlzbWF0Y2hFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use completeWebAuthnRegistrationRequestDescriptor instead')
const CompleteWebAuthnRegistrationRequest$json = {
  '1': 'CompleteWebAuthnRegistrationRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'credential', '3': 265594649, '4': 1, '5': 9, '10': 'credential'},
  ],
};

/// Descriptor for `CompleteWebAuthnRegistrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completeWebAuthnRegistrationRequestDescriptor =
    $convert.base64Decode(
        'CiNDb21wbGV0ZVdlYkF1dGhuUmVnaXN0cmF0aW9uUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJ'
        'BGIAEoCVILYWNjZXNzdG9rZW4SIQoKY3JlZGVudGlhbBiZztJ+IAEoCVIKY3JlZGVudGlhbA==');

@$core.Deprecated('Use completeWebAuthnRegistrationResponseDescriptor instead')
const CompleteWebAuthnRegistrationResponse$json = {
  '1': 'CompleteWebAuthnRegistrationResponse',
};

/// Descriptor for `CompleteWebAuthnRegistrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completeWebAuthnRegistrationResponseDescriptor =
    $convert
        .base64Decode('CiRDb21wbGV0ZVdlYkF1dGhuUmVnaXN0cmF0aW9uUmVzcG9uc2U=');

@$core.Deprecated('Use compromisedCredentialsActionsTypeDescriptor instead')
const CompromisedCredentialsActionsType$json = {
  '1': 'CompromisedCredentialsActionsType',
  '2': [
    {
      '1': 'eventaction',
      '3': 85257574,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.CompromisedCredentialsEventActionType',
      '10': 'eventaction'
    },
  ],
};

/// Descriptor for `CompromisedCredentialsActionsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List compromisedCredentialsActionsTypeDescriptor =
    $convert.base64Decode(
        'CiFDb21wcm9taXNlZENyZWRlbnRpYWxzQWN0aW9uc1R5cGUSVwoLZXZlbnRhY3Rpb24Y5trTKC'
        'ABKA4yMi5jb2duaXRvX2lkcC5Db21wcm9taXNlZENyZWRlbnRpYWxzRXZlbnRBY3Rpb25UeXBl'
        'UgtldmVudGFjdGlvbg==');

@$core.Deprecated(
    'Use compromisedCredentialsRiskConfigurationTypeDescriptor instead')
const CompromisedCredentialsRiskConfigurationType$json = {
  '1': 'CompromisedCredentialsRiskConfigurationType',
  '2': [
    {
      '1': 'actions',
      '3': 106935557,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CompromisedCredentialsActionsType',
      '10': 'actions'
    },
    {
      '1': 'eventfilter',
      '3': 285931770,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.EventFilterType',
      '10': 'eventfilter'
    },
  ],
};

/// Descriptor for `CompromisedCredentialsRiskConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    compromisedCredentialsRiskConfigurationTypeDescriptor =
    $convert.base64Decode(
        'CitDb21wcm9taXNlZENyZWRlbnRpYWxzUmlza0NvbmZpZ3VyYXRpb25UeXBlEksKB2FjdGlvbn'
        'MYher+MiABKAsyLi5jb2duaXRvX2lkcC5Db21wcm9taXNlZENyZWRlbnRpYWxzQWN0aW9uc1R5'
        'cGVSB2FjdGlvbnMSQgoLZXZlbnRmaWx0ZXIY+vGriAEgAygOMhwuY29nbml0b19pZHAuRXZlbn'
        'RGaWx0ZXJUeXBlUgtldmVudGZpbHRlcg==');

@$core.Deprecated('Use concurrentModificationExceptionDescriptor instead')
const ConcurrentModificationException$json = {
  '1': 'ConcurrentModificationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConcurrentModificationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrentModificationExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Db25jdXJyZW50TW9kaWZpY2F0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use confirmDeviceRequestDescriptor instead')
const ConfirmDeviceRequest$json = {
  '1': 'ConfirmDeviceRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
    {'1': 'devicename', '3': 514720633, '4': 1, '5': 9, '10': 'devicename'},
    {
      '1': 'devicesecretverifierconfig',
      '3': 310432876,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.DeviceSecretVerifierConfigType',
      '10': 'devicesecretverifierconfig'
    },
  ],
};

/// Descriptor for `ConfirmDeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmDeviceRequestDescriptor = $convert.base64Decode(
    'ChRDb25maXJtRGV2aWNlUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIAEoCVILYWNjZXNzdG'
    '9rZW4SIAoJZGV2aWNla2V5GKvkyLYBIAEoCVIJZGV2aWNla2V5EiIKCmRldmljZW5hbWUY+Ya4'
    '9QEgASgJUgpkZXZpY2VuYW1lEm8KGmRldmljZXNlY3JldHZlcmlmaWVyY29uZmlnGOyog5QBIA'
    'EoCzIrLmNvZ25pdG9faWRwLkRldmljZVNlY3JldFZlcmlmaWVyQ29uZmlnVHlwZVIaZGV2aWNl'
    'c2VjcmV0dmVyaWZpZXJjb25maWc=');

@$core.Deprecated('Use confirmDeviceResponseDescriptor instead')
const ConfirmDeviceResponse$json = {
  '1': 'ConfirmDeviceResponse',
  '2': [
    {
      '1': 'userconfirmationnecessary',
      '3': 445943147,
      '4': 1,
      '5': 8,
      '10': 'userconfirmationnecessary'
    },
  ],
};

/// Descriptor for `ConfirmDeviceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmDeviceResponseDescriptor = $convert.base64Decode(
    'ChVDb25maXJtRGV2aWNlUmVzcG9uc2USQAoZdXNlcmNvbmZpcm1hdGlvbm5lY2Vzc2FyeRjrmt'
    'LUASABKAhSGXVzZXJjb25maXJtYXRpb25uZWNlc3Nhcnk=');

@$core.Deprecated('Use confirmForgotPasswordRequestDescriptor instead')
const ConfirmForgotPasswordRequest$json = {
  '1': 'ConfirmForgotPasswordRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ConfirmForgotPasswordRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {
      '1': 'confirmationcode',
      '3': 109915362,
      '4': 1,
      '5': 9,
      '10': 'confirmationcode'
    },
    {'1': 'password', '3': 214108217, '4': 1, '5': 9, '10': 'password'},
    {'1': 'secrethash', '3': 321025630, '4': 1, '5': 9, '10': 'secrethash'},
    {
      '1': 'usercontextdata',
      '3': 169951134,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserContextDataType',
      '10': 'usercontextdata'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [ConfirmForgotPasswordRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use confirmForgotPasswordRequestDescriptor instead')
const ConfirmForgotPasswordRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ConfirmForgotPasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmForgotPasswordRequestDescriptor = $convert.base64Decode(
    'ChxDb25maXJtRm9yZ290UGFzc3dvcmRSZXF1ZXN0ElMKEWFuYWx5dGljc21ldGFkYXRhGPeFwT'
    'wgASgLMiIuY29nbml0b19pZHAuQW5hbHl0aWNzTWV0YWRhdGFUeXBlUhFhbmFseXRpY3NtZXRh'
    'ZGF0YRIeCghjbGllbnRpZBik6IbWASABKAlSCGNsaWVudGlkEmgKDmNsaWVudG1ldGFkYXRhGM'
    'yv/2EgAygLMj0uY29nbml0b19pZHAuQ29uZmlybUZvcmdvdFBhc3N3b3JkUmVxdWVzdC5DbGll'
    'bnRtZXRhZGF0YUVudHJ5Ug5jbGllbnRtZXRhZGF0YRItChBjb25maXJtYXRpb25jb2RlGOLZtD'
    'QgASgJUhBjb25maXJtYXRpb25jb2RlEh0KCHBhc3N3b3JkGLmQjGYgASgJUghwYXNzd29yZBIi'
    'CgpzZWNyZXRoYXNoGN7siZkBIAEoCVIKc2VjcmV0aGFzaBJNCg91c2VyY29udGV4dGRhdGEYnv'
    '+EUSABKAsyIC5jb2duaXRvX2lkcC5Vc2VyQ29udGV4dERhdGFUeXBlUg91c2VyY29udGV4dGRh'
    'dGESHgoIdXNlcm5hbWUY2qmj4AEgASgJUgh1c2VybmFtZRpBChNDbGllbnRtZXRhZGF0YUVudH'
    'J5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use confirmForgotPasswordResponseDescriptor instead')
const ConfirmForgotPasswordResponse$json = {
  '1': 'ConfirmForgotPasswordResponse',
};

/// Descriptor for `ConfirmForgotPasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmForgotPasswordResponseDescriptor =
    $convert.base64Decode('Ch1Db25maXJtRm9yZ290UGFzc3dvcmRSZXNwb25zZQ==');

@$core.Deprecated('Use confirmSignUpRequestDescriptor instead')
const ConfirmSignUpRequest$json = {
  '1': 'ConfirmSignUpRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ConfirmSignUpRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {
      '1': 'confirmationcode',
      '3': 109915362,
      '4': 1,
      '5': 9,
      '10': 'confirmationcode'
    },
    {
      '1': 'forcealiascreation',
      '3': 449762314,
      '4': 1,
      '5': 8,
      '10': 'forcealiascreation'
    },
    {'1': 'secrethash', '3': 321025630, '4': 1, '5': 9, '10': 'secrethash'},
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {
      '1': 'usercontextdata',
      '3': 169951134,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserContextDataType',
      '10': 'usercontextdata'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [ConfirmSignUpRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use confirmSignUpRequestDescriptor instead')
const ConfirmSignUpRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ConfirmSignUpRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmSignUpRequestDescriptor = $convert.base64Decode(
    'ChRDb25maXJtU2lnblVwUmVxdWVzdBJTChFhbmFseXRpY3NtZXRhZGF0YRj3hcE8IAEoCzIiLm'
    'NvZ25pdG9faWRwLkFuYWx5dGljc01ldGFkYXRhVHlwZVIRYW5hbHl0aWNzbWV0YWRhdGESHgoI'
    'Y2xpZW50aWQYpOiG1gEgASgJUghjbGllbnRpZBJgCg5jbGllbnRtZXRhZGF0YRjMr/9hIAMoCz'
    'I1LmNvZ25pdG9faWRwLkNvbmZpcm1TaWduVXBSZXF1ZXN0LkNsaWVudG1ldGFkYXRhRW50cnlS'
    'DmNsaWVudG1ldGFkYXRhEi0KEGNvbmZpcm1hdGlvbmNvZGUY4tm0NCABKAlSEGNvbmZpcm1hdG'
    'lvbmNvZGUSMgoSZm9yY2VhbGlhc2NyZWF0aW9uGIqou9YBIAEoCFISZm9yY2VhbGlhc2NyZWF0'
    'aW9uEiIKCnNlY3JldGhhc2gY3uyJmQEgASgJUgpzZWNyZXRoYXNoEhsKB3Nlc3Npb24YmJmjAi'
    'ABKAlSB3Nlc3Npb24STQoPdXNlcmNvbnRleHRkYXRhGJ7/hFEgASgLMiAuY29nbml0b19pZHAu'
    'VXNlckNvbnRleHREYXRhVHlwZVIPdXNlcmNvbnRleHRkYXRhEh4KCHVzZXJuYW1lGNqpo+ABIA'
    'EoCVIIdXNlcm5hbWUaQQoTQ2xpZW50bWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIU'
    'CgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use confirmSignUpResponseDescriptor instead')
const ConfirmSignUpResponse$json = {
  '1': 'ConfirmSignUpResponse',
  '2': [
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
  ],
};

/// Descriptor for `ConfirmSignUpResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmSignUpResponseDescriptor = $convert.base64Decode(
    'ChVDb25maXJtU2lnblVwUmVzcG9uc2USGwoHc2Vzc2lvbhiYmaMCIAEoCVIHc2Vzc2lvbg==');

@$core.Deprecated('Use contextDataTypeDescriptor instead')
const ContextDataType$json = {
  '1': 'ContextDataType',
  '2': [
    {'1': 'encodeddata', '3': 153276170, '4': 1, '5': 9, '10': 'encodeddata'},
    {
      '1': 'httpheaders',
      '3': 214337162,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.HttpHeader',
      '10': 'httpheaders'
    },
    {'1': 'ipaddress', '3': 1800397, '4': 1, '5': 9, '10': 'ipaddress'},
    {'1': 'servername', '3': 47211162, '4': 1, '5': 9, '10': 'servername'},
    {'1': 'serverpath', '3': 200845298, '4': 1, '5': 9, '10': 'serverpath'},
  ],
};

/// Descriptor for `ContextDataType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextDataTypeDescriptor = $convert.base64Decode(
    'Cg9Db250ZXh0RGF0YVR5cGUSIwoLZW5jb2RlZGRhdGEYip6LSSABKAlSC2VuY29kZWRkYXRhEj'
    'wKC2h0dHBoZWFkZXJzGIqNmmYgAygLMhcuY29nbml0b19pZHAuSHR0cEhlYWRlclILaHR0cGhl'
    'YWRlcnMSHgoJaXBhZGRyZXNzGM3xbSABKAlSCWlwYWRkcmVzcxIhCgpzZXJ2ZXJuYW1lGJrFwR'
    'YgASgJUgpzZXJ2ZXJuYW1lEiEKCnNlcnZlcnBhdGgY8s/iXyABKAlSCnNlcnZlcnBhdGg=');

@$core.Deprecated('Use createGroupRequestDescriptor instead')
const CreateGroupRequest$json = {
  '1': 'CreateGroupRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'precedence', '3': 206704142, '4': 1, '5': 5, '10': 'precedence'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `CreateGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGroupRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVHcm91cFJlcXVlc3QSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW'
    '9uEiAKCWdyb3VwbmFtZRjIyqCqASABKAlSCWdyb3VwbmFtZRIhCgpwcmVjZWRlbmNlGI6cyGIg'
    'ASgFUgpwcmVjZWRlbmNlEhwKB3JvbGVhcm4YgfjnmQEgASgJUgdyb2xlYXJuEiIKCnVzZXJwb2'
    '9saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use createGroupResponseDescriptor instead')
const CreateGroupResponse$json = {
  '1': 'CreateGroupResponse',
  '2': [
    {
      '1': 'group',
      '3': 91525165,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.GroupType',
      '10': 'group'
    },
  ],
};

/// Descriptor for `CreateGroupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGroupResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVHcm91cFJlc3BvbnNlEi8KBWdyb3VwGK2g0isgASgLMhYuY29nbml0b19pZHAuR3'
    'JvdXBUeXBlUgVncm91cA==');

@$core.Deprecated('Use createIdentityProviderRequestDescriptor instead')
const CreateIdentityProviderRequest$json = {
  '1': 'CreateIdentityProviderRequest',
  '2': [
    {
      '1': 'attributemapping',
      '3': 116923092,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.CreateIdentityProviderRequest.AttributemappingEntry',
      '10': 'attributemapping'
    },
    {
      '1': 'idpidentifiers',
      '3': 205051409,
      '4': 3,
      '5': 9,
      '10': 'idpidentifiers'
    },
    {
      '1': 'providerdetails',
      '3': 476397115,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.CreateIdentityProviderRequest.ProviderdetailsEntry',
      '10': 'providerdetails'
    },
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
    {
      '1': 'providertype',
      '3': 337296789,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.IdentityProviderTypeType',
      '10': 'providertype'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [
    CreateIdentityProviderRequest_AttributemappingEntry$json,
    CreateIdentityProviderRequest_ProviderdetailsEntry$json
  ],
};

@$core.Deprecated('Use createIdentityProviderRequestDescriptor instead')
const CreateIdentityProviderRequest_AttributemappingEntry$json = {
  '1': 'AttributemappingEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use createIdentityProviderRequestDescriptor instead')
const CreateIdentityProviderRequest_ProviderdetailsEntry$json = {
  '1': 'ProviderdetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateIdentityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createIdentityProviderRequestDescriptor = $convert.base64Decode(
    'Ch1DcmVhdGVJZGVudGl0eVByb3ZpZGVyUmVxdWVzdBJvChBhdHRyaWJ1dGVtYXBwaW5nGNS14D'
    'cgAygLMkAuY29nbml0b19pZHAuQ3JlYXRlSWRlbnRpdHlQcm92aWRlclJlcXVlc3QuQXR0cmli'
    'dXRlbWFwcGluZ0VudHJ5UhBhdHRyaWJ1dGVtYXBwaW5nEikKDmlkcGlkZW50aWZpZXJzGJGs42'
    'EgAygJUg5pZHBpZGVudGlmaWVycxJtCg9wcm92aWRlcmRldGFpbHMYu/yU4wEgAygLMj8uY29n'
    'bml0b19pZHAuQ3JlYXRlSWRlbnRpdHlQcm92aWRlclJlcXVlc3QuUHJvdmlkZXJkZXRhaWxzRW'
    '50cnlSD3Byb3ZpZGVyZGV0YWlscxImCgxwcm92aWRlcm5hbWUY+KGo5wEgASgJUgxwcm92aWRl'
    'cm5hbWUSTQoMcHJvdmlkZXJ0eXBlGJX76qABIAEoDjIlLmNvZ25pdG9faWRwLklkZW50aXR5UH'
    'JvdmlkZXJUeXBlVHlwZVIMcHJvdmlkZXJ0eXBlEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1'
    'c2VycG9vbGlkGkMKFUF0dHJpYnV0ZW1hcHBpbmdFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCg'
    'V2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGkIKFFByb3ZpZGVyZGV0YWlsc0VudHJ5EhAKA2tleRgB'
    'IAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createIdentityProviderResponseDescriptor instead')
const CreateIdentityProviderResponse$json = {
  '1': 'CreateIdentityProviderResponse',
  '2': [
    {
      '1': 'identityprovider',
      '3': 450450187,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.IdentityProviderType',
      '10': 'identityprovider'
    },
  ],
};

/// Descriptor for `CreateIdentityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createIdentityProviderResponseDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVJZGVudGl0eVByb3ZpZGVyUmVzcG9uc2USUQoQaWRlbnRpdHlwcm92aWRlchiLpu'
        'XWASABKAsyIS5jb2duaXRvX2lkcC5JZGVudGl0eVByb3ZpZGVyVHlwZVIQaWRlbnRpdHlwcm92'
        'aWRlcg==');

@$core.Deprecated('Use createManagedLoginBrandingRequestDescriptor instead')
const CreateManagedLoginBrandingRequest$json = {
  '1': 'CreateManagedLoginBrandingRequest',
  '2': [
    {
      '1': 'assets',
      '3': 141150041,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AssetType',
      '10': 'assets'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'settings', '3': 184911657, '4': 1, '5': 9, '10': 'settings'},
    {
      '1': 'usecognitoprovidedvalues',
      '3': 408662835,
      '4': 1,
      '5': 8,
      '10': 'usecognitoprovidedvalues'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `CreateManagedLoginBrandingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createManagedLoginBrandingRequestDescriptor = $convert.base64Decode(
    'CiFDcmVhdGVNYW5hZ2VkTG9naW5CcmFuZGluZ1JlcXVlc3QSMQoGYXNzZXRzGNmOp0MgAygLMh'
    'YuY29nbml0b19pZHAuQXNzZXRUeXBlUgZhc3NldHMSHgoIY2xpZW50aWQYpOiG1gEgASgJUghj'
    'bGllbnRpZBIdCghzZXR0aW5ncxipjpZYIAEoCVIIc2V0dGluZ3MSPgoYdXNlY29nbml0b3Byb3'
    'ZpZGVkdmFsdWVzGLPm7sIBIAEoCFIYdXNlY29nbml0b3Byb3ZpZGVkdmFsdWVzEiIKCnVzZXJw'
    'b29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use createManagedLoginBrandingResponseDescriptor instead')
const CreateManagedLoginBrandingResponse$json = {
  '1': 'CreateManagedLoginBrandingResponse',
  '2': [
    {
      '1': 'managedloginbranding',
      '3': 338791109,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ManagedLoginBrandingType',
      '10': 'managedloginbranding'
    },
  ],
};

/// Descriptor for `CreateManagedLoginBrandingResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createManagedLoginBrandingResponseDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVNYW5hZ2VkTG9naW5CcmFuZGluZ1Jlc3BvbnNlEl0KFG1hbmFnZWRsb2dpbmJyYW'
        '5kaW5nGMWVxqEBIAEoCzIlLmNvZ25pdG9faWRwLk1hbmFnZWRMb2dpbkJyYW5kaW5nVHlwZVIU'
        'bWFuYWdlZGxvZ2luYnJhbmRpbmc=');

@$core.Deprecated('Use createResourceServerRequestDescriptor instead')
const CreateResourceServerRequest$json = {
  '1': 'CreateResourceServerRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'scopes',
      '3': 464684393,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ResourceServerScopeType',
      '10': 'scopes'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `CreateResourceServerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResourceServerRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVSZXNvdXJjZVNlcnZlclJlcXVlc3QSIQoKaWRlbnRpZmllchjfoPsTIAEoCVIKaW'
    'RlbnRpZmllchIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkAKBnNjb3BlcxjpisrdASADKAsyJC5j'
    'b2duaXRvX2lkcC5SZXNvdXJjZVNlcnZlclNjb3BlVHlwZVIGc2NvcGVzEiIKCnVzZXJwb29saW'
    'QY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use createResourceServerResponseDescriptor instead')
const CreateResourceServerResponse$json = {
  '1': 'CreateResourceServerResponse',
  '2': [
    {
      '1': 'resourceserver',
      '3': 506282051,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ResourceServerType',
      '10': 'resourceserver'
    },
  ],
};

/// Descriptor for `CreateResourceServerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResourceServerResponseDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVSZXNvdXJjZVNlcnZlclJlc3BvbnNlEksKDnJlc291cmNlc2VydmVyGMOAtfEBIA'
        'EoCzIfLmNvZ25pdG9faWRwLlJlc291cmNlU2VydmVyVHlwZVIOcmVzb3VyY2VzZXJ2ZXI=');

@$core.Deprecated('Use createTermsRequestDescriptor instead')
const CreateTermsRequest$json = {
  '1': 'CreateTermsRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'enforcement',
      '3': 412213242,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TermsEnforcementType',
      '10': 'enforcement'
    },
    {
      '1': 'links',
      '3': 302123151,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.CreateTermsRequest.LinksEntry',
      '10': 'links'
    },
    {'1': 'termsname', '3': 303051560, '4': 1, '5': 9, '10': 'termsname'},
    {
      '1': 'termssource',
      '3': 122689594,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TermsSourceType',
      '10': 'termssource'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [CreateTermsRequest_LinksEntry$json],
};

@$core.Deprecated('Use createTermsRequestDescriptor instead')
const CreateTermsRequest_LinksEntry$json = {
  '1': 'LinksEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateTermsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTermsRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVUZXJtc1JlcXVlc3QSHgoIY2xpZW50aWQYpOiG1gEgASgJUghjbGllbnRpZBJHCg'
    'tlbmZvcmNlbWVudBj6v8fEASABKA4yIS5jb2duaXRvX2lkcC5UZXJtc0VuZm9yY2VtZW50VHlw'
    'ZVILZW5mb3JjZW1lbnQSRAoFbGlua3MYj5GIkAEgAygLMiouY29nbml0b19pZHAuQ3JlYXRlVG'
    'VybXNSZXF1ZXN0LkxpbmtzRW50cnlSBWxpbmtzEiAKCXRlcm1zbmFtZRio5sCQASABKAlSCXRl'
    'cm1zbmFtZRJBCgt0ZXJtc3NvdXJjZRi6sMA6IAEoDjIcLmNvZ25pdG9faWRwLlRlcm1zU291cm'
    'NlVHlwZVILdGVybXNzb3VyY2USIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQa'
    'OAoKTGlua3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6Aj'
    'gB');

@$core.Deprecated('Use createTermsResponseDescriptor instead')
const CreateTermsResponse$json = {
  '1': 'CreateTermsResponse',
  '2': [
    {
      '1': 'terms',
      '3': 339062221,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.TermsType',
      '10': 'terms'
    },
  ],
};

/// Descriptor for `CreateTermsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTermsResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVUZXJtc1Jlc3BvbnNlEjAKBXRlcm1zGM3b1qEBIAEoCzIWLmNvZ25pdG9faWRwLl'
    'Rlcm1zVHlwZVIFdGVybXM=');

@$core.Deprecated('Use createUserImportJobRequestDescriptor instead')
const CreateUserImportJobRequest$json = {
  '1': 'CreateUserImportJobRequest',
  '2': [
    {
      '1': 'cloudwatchlogsrolearn',
      '3': 55454690,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsrolearn'
    },
    {'1': 'jobname', '3': 498531160, '4': 1, '5': 9, '10': 'jobname'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `CreateUserImportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserImportJobRequestDescriptor =
    $convert.base64Decode(
        'ChpDcmVhdGVVc2VySW1wb3J0Sm9iUmVxdWVzdBI3ChVjbG91ZHdhdGNobG9nc3JvbGVhcm4Y4t'
        'e4GiABKAlSFWNsb3Vkd2F0Y2hsb2dzcm9sZWFybhIcCgdqb2JuYW1lGNj22+0BIAEoCVIHam9i'
        'bmFtZRIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use createUserImportJobResponseDescriptor instead')
const CreateUserImportJobResponse$json = {
  '1': 'CreateUserImportJobResponse',
  '2': [
    {
      '1': 'userimportjob',
      '3': 473682999,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserImportJobType',
      '10': 'userimportjob'
    },
  ],
};

/// Descriptor for `CreateUserImportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserImportJobResponseDescriptor =
    $convert.base64Decode(
        'ChtDcmVhdGVVc2VySW1wb3J0Sm9iUmVzcG9uc2USSAoNdXNlcmltcG9ydGpvYhi3qO/hASABKA'
        'syHi5jb2duaXRvX2lkcC5Vc2VySW1wb3J0Sm9iVHlwZVINdXNlcmltcG9ydGpvYg==');

@$core.Deprecated('Use createUserPoolClientRequestDescriptor instead')
const CreateUserPoolClientRequest$json = {
  '1': 'CreateUserPoolClientRequest',
  '2': [
    {
      '1': 'accesstokenvalidity',
      '3': 260874267,
      '4': 1,
      '5': 5,
      '10': 'accesstokenvalidity'
    },
    {
      '1': 'allowedoauthflows',
      '3': 268290584,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.OAuthFlowType',
      '10': 'allowedoauthflows'
    },
    {
      '1': 'allowedoauthflowsuserpoolclient',
      '3': 520095610,
      '4': 1,
      '5': 8,
      '10': 'allowedoauthflowsuserpoolclient'
    },
    {
      '1': 'allowedoauthscopes',
      '3': 39385504,
      '4': 3,
      '5': 9,
      '10': 'allowedoauthscopes'
    },
    {
      '1': 'analyticsconfiguration',
      '3': 229750388,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsConfigurationType',
      '10': 'analyticsconfiguration'
    },
    {
      '1': 'authsessionvalidity',
      '3': 223873468,
      '4': 1,
      '5': 5,
      '10': 'authsessionvalidity'
    },
    {'1': 'callbackurls', '3': 227703885, '4': 3, '5': 9, '10': 'callbackurls'},
    {'1': 'clientname', '3': 340245630, '4': 1, '5': 9, '10': 'clientname'},
    {'1': 'clientsecret', '3': 500734711, '4': 1, '5': 9, '10': 'clientsecret'},
    {
      '1': 'defaultredirecturi',
      '3': 311293253,
      '4': 1,
      '5': 9,
      '10': 'defaultredirecturi'
    },
    {
      '1': 'enablepropagateadditionalusercontextdata',
      '3': 201651031,
      '4': 1,
      '5': 8,
      '10': 'enablepropagateadditionalusercontextdata'
    },
    {
      '1': 'enabletokenrevocation',
      '3': 178186392,
      '4': 1,
      '5': 8,
      '10': 'enabletokenrevocation'
    },
    {
      '1': 'explicitauthflows',
      '3': 277179621,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.ExplicitAuthFlowsType',
      '10': 'explicitauthflows'
    },
    {
      '1': 'generatesecret',
      '3': 233116579,
      '4': 1,
      '5': 8,
      '10': 'generatesecret'
    },
    {
      '1': 'idtokenvalidity',
      '3': 312934952,
      '4': 1,
      '5': 5,
      '10': 'idtokenvalidity'
    },
    {'1': 'logouturls', '3': 468187518, '4': 3, '5': 9, '10': 'logouturls'},
    {
      '1': 'preventuserexistenceerrors',
      '3': 188235606,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.PreventUserExistenceErrorTypes',
      '10': 'preventuserexistenceerrors'
    },
    {
      '1': 'readattributes',
      '3': 334413205,
      '4': 3,
      '5': 9,
      '10': 'readattributes'
    },
    {
      '1': 'refreshtokenrotation',
      '3': 199284564,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.RefreshTokenRotationType',
      '10': 'refreshtokenrotation'
    },
    {
      '1': 'refreshtokenvalidity',
      '3': 303433364,
      '4': 1,
      '5': 5,
      '10': 'refreshtokenvalidity'
    },
    {
      '1': 'supportedidentityproviders',
      '3': 439564368,
      '4': 3,
      '5': 9,
      '10': 'supportedidentityproviders'
    },
    {
      '1': 'tokenvalidityunits',
      '3': 2056664,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.TokenValidityUnitsType',
      '10': 'tokenvalidityunits'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {
      '1': 'writeattributes',
      '3': 440236318,
      '4': 3,
      '5': 9,
      '10': 'writeattributes'
    },
  ],
};

/// Descriptor for `CreateUserPoolClientRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserPoolClientRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVVc2VyUG9vbENsaWVudFJlcXVlc3QSMwoTYWNjZXNzdG9rZW52YWxpZGl0eRibwL'
    'J8IAEoBVITYWNjZXNzdG9rZW52YWxpZGl0eRJLChFhbGxvd2Vkb2F1dGhmbG93cxiYlPd/IAMo'
    'DjIaLmNvZ25pdG9faWRwLk9BdXRoRmxvd1R5cGVSEWFsbG93ZWRvYXV0aGZsb3dzEkwKH2FsbG'
    '93ZWRvYXV0aGZsb3dzdXNlcnBvb2xjbGllbnQY+o6A+AEgASgIUh9hbGxvd2Vkb2F1dGhmbG93'
    'c3VzZXJwb29sY2xpZW50EjEKEmFsbG93ZWRvYXV0aHNjb3Blcxig8+MSIAMoCVISYWxsb3dlZG'
    '9hdXRoc2NvcGVzEmIKFmFuYWx5dGljc2NvbmZpZ3VyYXRpb24Y9OzGbSABKAsyJy5jb2duaXRv'
    'X2lkcC5BbmFseXRpY3NDb25maWd1cmF0aW9uVHlwZVIWYW5hbHl0aWNzY29uZmlndXJhdGlvbh'
    'IzChNhdXRoc2Vzc2lvbnZhbGlkaXR5GLyT4GogASgFUhNhdXRoc2Vzc2lvbnZhbGlkaXR5EiUK'
    'DGNhbGxiYWNrdXJscxjN+MlsIAMoCVIMY2FsbGJhY2t1cmxzEiIKCmNsaWVudG5hbWUY/vieog'
    'EgASgJUgpjbGllbnRuYW1lEiYKDGNsaWVudHNlY3JldBj3teLuASABKAlSDGNsaWVudHNlY3Jl'
    'dBIyChJkZWZhdWx0cmVkaXJlY3R1cmkYxeq3lAEgASgJUhJkZWZhdWx0cmVkaXJlY3R1cmkSXQ'
    'ooZW5hYmxlcHJvcGFnYXRlYWRkaXRpb25hbHVzZXJjb250ZXh0ZGF0YRjX5pNgIAEoCFIoZW5h'
    'YmxlcHJvcGFnYXRlYWRkaXRpb25hbHVzZXJjb250ZXh0ZGF0YRI3ChVlbmFibGV0b2tlbnJldm'
    '9jYXRpb24YmNH7VCABKAhSFWVuYWJsZXRva2VucmV2b2NhdGlvbhJUChFleHBsaWNpdGF1dGhm'
    'bG93cxjl2ZWEASADKA4yIi5jb2duaXRvX2lkcC5FeHBsaWNpdEF1dGhGbG93c1R5cGVSEWV4cG'
    'xpY2l0YXV0aGZsb3dzEikKDmdlbmVyYXRlc2VjcmV0GKOnlG8gASgIUg5nZW5lcmF0ZXNlY3Jl'
    'dBIsCg9pZHRva2VudmFsaWRpdHkYqISclQEgASgFUg9pZHRva2VudmFsaWRpdHkSIgoKbG9nb3'
    'V0dXJscxj+8p/fASADKAlSCmxvZ291dHVybHMSbgoacHJldmVudHVzZXJleGlzdGVuY2VlcnJv'
    'cnMY1v7gWSABKA4yKy5jb2duaXRvX2lkcC5QcmV2ZW50VXNlckV4aXN0ZW5jZUVycm9yVHlwZX'
    'NSGnByZXZlbnR1c2VyZXhpc3RlbmNlZXJyb3JzEioKDnJlYWRhdHRyaWJ1dGVzGJX7up8BIAMo'
    'CVIOcmVhZGF0dHJpYnV0ZXMSXAoUcmVmcmVzaHRva2Vucm90YXRpb24Y1K6DXyABKAsyJS5jb2'
    'duaXRvX2lkcC5SZWZyZXNoVG9rZW5Sb3RhdGlvblR5cGVSFHJlZnJlc2h0b2tlbnJvdGF0aW9u'
    'EjYKFHJlZnJlc2h0b2tlbnZhbGlkaXR5GJSN2JABIAEoBVIUcmVmcmVzaHRva2VudmFsaWRpdH'
    'kSQgoac3VwcG9ydGVkaWRlbnRpdHlwcm92aWRlcnMY0PDM0QEgAygJUhpzdXBwb3J0ZWRpZGVu'
    'dGl0eXByb3ZpZGVycxJVChJ0b2tlbnZhbGlkaXR5dW5pdHMY2MN9IAEoCzIjLmNvZ25pdG9faW'
    'RwLlRva2VuVmFsaWRpdHlVbml0c1R5cGVSEnRva2VudmFsaWRpdHl1bml0cxIiCgp1c2VycG9v'
    'bGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZBIsCg93cml0ZWF0dHJpYnV0ZXMYnvL10QEgAygJUg'
    '93cml0ZWF0dHJpYnV0ZXM=');

@$core.Deprecated('Use createUserPoolClientResponseDescriptor instead')
const CreateUserPoolClientResponse$json = {
  '1': 'CreateUserPoolClientResponse',
  '2': [
    {
      '1': 'userpoolclient',
      '3': 138497904,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolClientType',
      '10': 'userpoolclient'
    },
  ],
};

/// Descriptor for `CreateUserPoolClientResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserPoolClientResponseDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVVc2VyUG9vbENsaWVudFJlc3BvbnNlEkoKDnVzZXJwb29sY2xpZW50GPCehUIgAS'
        'gLMh8uY29nbml0b19pZHAuVXNlclBvb2xDbGllbnRUeXBlUg51c2VycG9vbGNsaWVudA==');

@$core.Deprecated('Use createUserPoolDomainRequestDescriptor instead')
const CreateUserPoolDomainRequest$json = {
  '1': 'CreateUserPoolDomainRequest',
  '2': [
    {
      '1': 'customdomainconfig',
      '3': 472479421,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CustomDomainConfigType',
      '10': 'customdomainconfig'
    },
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'managedloginversion',
      '3': 479901038,
      '4': 1,
      '5': 5,
      '10': 'managedloginversion'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `CreateUserPoolDomainRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserPoolDomainRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVVc2VyUG9vbERvbWFpblJlcXVlc3QSVwoSY3VzdG9tZG9tYWluY29uZmlnGL3tpe'
    'EBIAEoCzIjLmNvZ25pdG9faWRwLkN1c3RvbURvbWFpbkNvbmZpZ1R5cGVSEmN1c3RvbWRvbWFp'
    'bmNvbmZpZxIaCgZkb21haW4YkpLy8AEgASgJUgZkb21haW4SNAoTbWFuYWdlZGxvZ2ludmVyc2'
    'lvbhju6urkASABKAVSE21hbmFnZWRsb2dpbnZlcnNpb24SIgoKdXNlcnBvb2xpZBj+xoudASAB'
    'KAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use createUserPoolDomainResponseDescriptor instead')
const CreateUserPoolDomainResponse$json = {
  '1': 'CreateUserPoolDomainResponse',
  '2': [
    {
      '1': 'cloudfrontdomain',
      '3': 436051942,
      '4': 1,
      '5': 9,
      '10': 'cloudfrontdomain'
    },
    {
      '1': 'managedloginversion',
      '3': 479901038,
      '4': 1,
      '5': 5,
      '10': 'managedloginversion'
    },
  ],
};

/// Descriptor for `CreateUserPoolDomainResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserPoolDomainResponseDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVVc2VyUG9vbERvbWFpblJlc3BvbnNlEi4KEGNsb3VkZnJvbnRkb21haW4Y5r/2zw'
        'EgASgJUhBjbG91ZGZyb250ZG9tYWluEjQKE21hbmFnZWRsb2dpbnZlcnNpb24Y7urq5AEgASgF'
        'UhNtYW5hZ2VkbG9naW52ZXJzaW9u');

@$core.Deprecated('Use createUserPoolRequestDescriptor instead')
const CreateUserPoolRequest$json = {
  '1': 'CreateUserPoolRequest',
  '2': [
    {
      '1': 'accountrecoverysetting',
      '3': 219232186,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountRecoverySettingType',
      '10': 'accountrecoverysetting'
    },
    {
      '1': 'admincreateuserconfig',
      '3': 364968418,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AdminCreateUserConfigType',
      '10': 'admincreateuserconfig'
    },
    {
      '1': 'aliasattributes',
      '3': 189876251,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.AliasAttributeType',
      '10': 'aliasattributes'
    },
    {
      '1': 'autoverifiedattributes',
      '3': 467729812,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.VerifiedAttributeType',
      '10': 'autoverifiedattributes'
    },
    {
      '1': 'deletionprotection',
      '3': 504781905,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DeletionProtectionType',
      '10': 'deletionprotection'
    },
    {
      '1': 'deviceconfiguration',
      '3': 512944140,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.DeviceConfigurationType',
      '10': 'deviceconfiguration'
    },
    {
      '1': 'emailconfiguration',
      '3': 528317976,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailConfigurationType',
      '10': 'emailconfiguration'
    },
    {
      '1': 'emailverificationmessage',
      '3': 172634664,
      '4': 1,
      '5': 9,
      '10': 'emailverificationmessage'
    },
    {
      '1': 'emailverificationsubject',
      '3': 224298169,
      '4': 1,
      '5': 9,
      '10': 'emailverificationsubject'
    },
    {
      '1': 'lambdaconfig',
      '3': 291837797,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.LambdaConfigType',
      '10': 'lambdaconfig'
    },
    {
      '1': 'mfaconfiguration',
      '3': 259020766,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolMfaType',
      '10': 'mfaconfiguration'
    },
    {
      '1': 'policies',
      '3': 40015384,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolPolicyType',
      '10': 'policies'
    },
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
    {
      '1': 'schema',
      '3': 412122455,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.SchemaAttributeType',
      '10': 'schema'
    },
    {
      '1': 'smsauthenticationmessage',
      '3': 356104990,
      '4': 1,
      '5': 9,
      '10': 'smsauthenticationmessage'
    },
    {
      '1': 'smsconfiguration',
      '3': 10321849,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SmsConfigurationType',
      '10': 'smsconfiguration'
    },
    {
      '1': 'smsverificationmessage',
      '3': 497665917,
      '4': 1,
      '5': 9,
      '10': 'smsverificationmessage'
    },
    {
      '1': 'userattributeupdatesettings',
      '3': 319670235,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserAttributeUpdateSettingsType',
      '10': 'userattributeupdatesettings'
    },
    {
      '1': 'userpooladdons',
      '3': 296941112,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolAddOnsType',
      '10': 'userpooladdons'
    },
    {
      '1': 'userpooltags',
      '3': 341705322,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.CreateUserPoolRequest.UserpooltagsEntry',
      '10': 'userpooltags'
    },
    {
      '1': 'userpooltier',
      '3': 80461029,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolTierType',
      '10': 'userpooltier'
    },
    {
      '1': 'usernameattributes',
      '3': 196392641,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.UsernameAttributeType',
      '10': 'usernameattributes'
    },
    {
      '1': 'usernameconfiguration',
      '3': 15447334,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UsernameConfigurationType',
      '10': 'usernameconfiguration'
    },
    {
      '1': 'verificationmessagetemplate',
      '3': 502836004,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.VerificationMessageTemplateType',
      '10': 'verificationmessagetemplate'
    },
  ],
  '3': [CreateUserPoolRequest_UserpooltagsEntry$json],
};

@$core.Deprecated('Use createUserPoolRequestDescriptor instead')
const CreateUserPoolRequest_UserpooltagsEntry$json = {
  '1': 'UserpooltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateUserPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserPoolRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVVc2VyUG9vbFJlcXVlc3QSYgoWYWNjb3VudHJlY292ZXJ5c2V0dGluZxi678RoIA'
    'EoCzInLmNvZ25pdG9faWRwLkFjY291bnRSZWNvdmVyeVNldHRpbmdUeXBlUhZhY2NvdW50cmVj'
    'b3ZlcnlzZXR0aW5nEmAKFWFkbWluY3JlYXRldXNlcmNvbmZpZxji84OuASABKAsyJi5jb2duaX'
    'RvX2lkcC5BZG1pbkNyZWF0ZVVzZXJDb25maWdUeXBlUhVhZG1pbmNyZWF0ZXVzZXJjb25maWcS'
    'TAoPYWxpYXNhdHRyaWJ1dGVzGJuQxVogAygOMh8uY29nbml0b19pZHAuQWxpYXNBdHRyaWJ1dG'
    'VUeXBlUg9hbGlhc2F0dHJpYnV0ZXMSXgoWYXV0b3ZlcmlmaWVkYXR0cmlidXRlcxiU+4PfASAD'
    'KA4yIi5jb2duaXRvX2lkcC5WZXJpZmllZEF0dHJpYnV0ZVR5cGVSFmF1dG92ZXJpZmllZGF0dH'
    'JpYnV0ZXMSVwoSZGVsZXRpb25wcm90ZWN0aW9uGNG42fABIAEoDjIjLmNvZ25pdG9faWRwLkRl'
    'bGV0aW9uUHJvdGVjdGlvblR5cGVSEmRlbGV0aW9ucHJvdGVjdGlvbhJaChNkZXZpY2Vjb25maW'
    'd1cmF0aW9uGIzQy/QBIAEoCzIkLmNvZ25pdG9faWRwLkRldmljZUNvbmZpZ3VyYXRpb25UeXBl'
    'UhNkZXZpY2Vjb25maWd1cmF0aW9uElcKEmVtYWlsY29uZmlndXJhdGlvbhiY/PX7ASABKAsyIy'
    '5jb2duaXRvX2lkcC5FbWFpbENvbmZpZ3VyYXRpb25UeXBlUhJlbWFpbGNvbmZpZ3VyYXRpb24S'
    'PQoYZW1haWx2ZXJpZmljYXRpb25tZXNzYWdlGKjkqFIgASgJUhhlbWFpbHZlcmlmaWNhdGlvbm'
    '1lc3NhZ2USPQoYZW1haWx2ZXJpZmljYXRpb25zdWJqZWN0GLmJ+mogASgJUhhlbWFpbHZlcmlm'
    'aWNhdGlvbnN1YmplY3QSRQoMbGFtYmRhY29uZmlnGOWulIsBIAEoCzIdLmNvZ25pdG9faWRwLk'
    'xhbWJkYUNvbmZpZ1R5cGVSDGxhbWJkYWNvbmZpZxJLChBtZmFjb25maWd1cmF0aW9uGN6vwXsg'
    'ASgOMhwuY29nbml0b19pZHAuVXNlclBvb2xNZmFUeXBlUhBtZmFjb25maWd1cmF0aW9uEj4KCH'
    'BvbGljaWVzGJisihMgASgLMh8uY29nbml0b19pZHAuVXNlclBvb2xQb2xpY3lUeXBlUghwb2xp'
    'Y2llcxIdCghwb29sbmFtZRjJjYUnIAEoCVIIcG9vbG5hbWUSPAoGc2NoZW1hGNf6wcQBIAMoCz'
    'IgLmNvZ25pdG9faWRwLlNjaGVtYUF0dHJpYnV0ZVR5cGVSBnNjaGVtYRI+ChhzbXNhdXRoZW50'
    'aWNhdGlvbm1lc3NhZ2UYnvbmqQEgASgJUhhzbXNhdXRoZW50aWNhdGlvbm1lc3NhZ2USUAoQc2'
    '1zY29uZmlndXJhdGlvbhi5//UEIAEoCzIhLmNvZ25pdG9faWRwLlNtc0NvbmZpZ3VyYXRpb25U'
    'eXBlUhBzbXNjb25maWd1cmF0aW9uEjoKFnNtc3ZlcmlmaWNhdGlvbm1lc3NhZ2UY/Y6n7QEgAS'
    'gJUhZzbXN2ZXJpZmljYXRpb25tZXNzYWdlEnIKG3VzZXJhdHRyaWJ1dGV1cGRhdGVzZXR0aW5n'
    'cxjbj7eYASABKAsyLC5jb2duaXRvX2lkcC5Vc2VyQXR0cmlidXRlVXBkYXRlU2V0dGluZ3NUeX'
    'BlUht1c2VyYXR0cmlidXRldXBkYXRlc2V0dGluZ3MSSwoOdXNlcnBvb2xhZGRvbnMYuOzLjQEg'
    'ASgLMh8uY29nbml0b19pZHAuVXNlclBvb2xBZGRPbnNUeXBlUg51c2VycG9vbGFkZG9ucxJcCg'
    'x1c2VycG9vbHRhZ3MY6oT4ogEgAygLMjQuY29nbml0b19pZHAuQ3JlYXRlVXNlclBvb2xSZXF1'
    'ZXN0LlVzZXJwb29sdGFnc0VudHJ5Ugx1c2VycG9vbHRhZ3MSRAoMdXNlcnBvb2x0aWVyGOX5ri'
    'YgASgOMh0uY29nbml0b19pZHAuVXNlclBvb2xUaWVyVHlwZVIMdXNlcnBvb2x0aWVyElUKEnVz'
    'ZXJuYW1lYXR0cmlidXRlcxjB7dJdIAMoDjIiLmNvZ25pdG9faWRwLlVzZXJuYW1lQXR0cmlidX'
    'RlVHlwZVISdXNlcm5hbWVhdHRyaWJ1dGVzEl8KFXVzZXJuYW1lY29uZmlndXJhdGlvbhim6q4H'
    'IAEoCzImLmNvZ25pdG9faWRwLlVzZXJuYW1lQ29uZmlndXJhdGlvblR5cGVSFXVzZXJuYW1lY2'
    '9uZmlndXJhdGlvbhJyCht2ZXJpZmljYXRpb25tZXNzYWdldGVtcGxhdGUYpNbi7wEgASgLMiwu'
    'Y29nbml0b19pZHAuVmVyaWZpY2F0aW9uTWVzc2FnZVRlbXBsYXRlVHlwZVIbdmVyaWZpY2F0aW'
    '9ubWVzc2FnZXRlbXBsYXRlGj8KEVVzZXJwb29sdGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5'
    'EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createUserPoolResponseDescriptor instead')
const CreateUserPoolResponse$json = {
  '1': 'CreateUserPoolResponse',
  '2': [
    {
      '1': 'userpool',
      '3': 404993697,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolType',
      '10': 'userpool'
    },
  ],
};

/// Descriptor for `CreateUserPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createUserPoolResponseDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVVc2VyUG9vbFJlc3BvbnNlEjkKCHVzZXJwb29sGKHtjsEBIAEoCzIZLmNvZ25pdG'
        '9faWRwLlVzZXJQb29sVHlwZVIIdXNlcnBvb2w=');

@$core.Deprecated('Use customDomainConfigTypeDescriptor instead')
const CustomDomainConfigType$json = {
  '1': 'CustomDomainConfigType',
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

/// Descriptor for `CustomDomainConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customDomainConfigTypeDescriptor =
    $convert.base64Decode(
        'ChZDdXN0b21Eb21haW5Db25maWdUeXBlEikKDmNlcnRpZmljYXRlYXJuGPjKmSwgASgJUg5jZX'
        'J0aWZpY2F0ZWFybg==');

@$core.Deprecated('Use customEmailLambdaVersionConfigTypeDescriptor instead')
const CustomEmailLambdaVersionConfigType$json = {
  '1': 'CustomEmailLambdaVersionConfigType',
  '2': [
    {'1': 'lambdaarn', '3': 196258582, '4': 1, '5': 9, '10': 'lambdaarn'},
    {
      '1': 'lambdaversion',
      '3': 429236631,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.CustomEmailSenderLambdaVersionType',
      '10': 'lambdaversion'
    },
  ],
};

/// Descriptor for `CustomEmailLambdaVersionConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customEmailLambdaVersionConfigTypeDescriptor =
    $convert.base64Decode(
        'CiJDdXN0b21FbWFpbExhbWJkYVZlcnNpb25Db25maWdUeXBlEh8KCWxhbWJkYWFybhiW1spdIA'
        'EoCVIJbGFtYmRhYXJuElkKDWxhbWJkYXZlcnNpb24Yl8PWzAEgASgOMi8uY29nbml0b19pZHAu'
        'Q3VzdG9tRW1haWxTZW5kZXJMYW1iZGFWZXJzaW9uVHlwZVINbGFtYmRhdmVyc2lvbg==');

@$core.Deprecated('Use customSMSLambdaVersionConfigTypeDescriptor instead')
const CustomSMSLambdaVersionConfigType$json = {
  '1': 'CustomSMSLambdaVersionConfigType',
  '2': [
    {'1': 'lambdaarn', '3': 196258582, '4': 1, '5': 9, '10': 'lambdaarn'},
    {
      '1': 'lambdaversion',
      '3': 429236631,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.CustomSMSSenderLambdaVersionType',
      '10': 'lambdaversion'
    },
  ],
};

/// Descriptor for `CustomSMSLambdaVersionConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customSMSLambdaVersionConfigTypeDescriptor =
    $convert.base64Decode(
        'CiBDdXN0b21TTVNMYW1iZGFWZXJzaW9uQ29uZmlnVHlwZRIfCglsYW1iZGFhcm4YltbKXSABKA'
        'lSCWxhbWJkYWFybhJXCg1sYW1iZGF2ZXJzaW9uGJfD1swBIAEoDjItLmNvZ25pdG9faWRwLkN1'
        'c3RvbVNNU1NlbmRlckxhbWJkYVZlcnNpb25UeXBlUg1sYW1iZGF2ZXJzaW9u');

@$core.Deprecated('Use deleteGroupRequestDescriptor instead')
const DeleteGroupRequest$json = {
  '1': 'DeleteGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteGroupRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVHcm91cFJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1lEi'
    'IKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use deleteIdentityProviderRequestDescriptor instead')
const DeleteIdentityProviderRequest$json = {
  '1': 'DeleteIdentityProviderRequest',
  '2': [
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteIdentityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIdentityProviderRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVJZGVudGl0eVByb3ZpZGVyUmVxdWVzdBImCgxwcm92aWRlcm5hbWUY+KGo5wEgAS'
        'gJUgxwcm92aWRlcm5hbWUSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use deleteManagedLoginBrandingRequestDescriptor instead')
const DeleteManagedLoginBrandingRequest$json = {
  '1': 'DeleteManagedLoginBrandingRequest',
  '2': [
    {
      '1': 'managedloginbrandingid',
      '3': 57829482,
      '4': 1,
      '5': 9,
      '10': 'managedloginbrandingid'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteManagedLoginBrandingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteManagedLoginBrandingRequestDescriptor =
    $convert.base64Decode(
        'CiFEZWxldGVNYW5hZ2VkTG9naW5CcmFuZGluZ1JlcXVlc3QSOQoWbWFuYWdlZGxvZ2luYnJhbm'
        'RpbmdpZBjq0MkbIAEoCVIWbWFuYWdlZGxvZ2luYnJhbmRpbmdpZBIiCgp1c2VycG9vbGlkGP7G'
        'i50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use deleteResourceServerRequestDescriptor instead')
const DeleteResourceServerRequest$json = {
  '1': 'DeleteResourceServerRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteResourceServerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourceServerRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVSZXNvdXJjZVNlcnZlclJlcXVlc3QSIQoKaWRlbnRpZmllchjfoPsTIAEoCVIKaW'
        'RlbnRpZmllchIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use deleteTermsRequestDescriptor instead')
const DeleteTermsRequest$json = {
  '1': 'DeleteTermsRequest',
  '2': [
    {'1': 'termsid', '3': 331306210, '4': 1, '5': 9, '10': 'termsid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteTermsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTermsRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVUZXJtc1JlcXVlc3QSHAoHdGVybXNpZBjiqf2dASABKAlSB3Rlcm1zaWQSIgoKdX'
    'NlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use deleteUserAttributesRequestDescriptor instead')
const DeleteUserAttributesRequest$json = {
  '1': 'DeleteUserAttributesRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'userattributenames',
      '3': 132104459,
      '4': 3,
      '5': 9,
      '10': 'userattributenames'
    },
  ],
};

/// Descriptor for `DeleteUserAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserAttributesRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVVc2VyQXR0cmlidXRlc1JlcXVlc3QSIwoLYWNjZXNzdG9rZW4YibyQRiABKAlSC2'
        'FjY2Vzc3Rva2VuEjEKEnVzZXJhdHRyaWJ1dGVuYW1lcxiLgv8+IAMoCVISdXNlcmF0dHJpYnV0'
        'ZW5hbWVz');

@$core.Deprecated('Use deleteUserAttributesResponseDescriptor instead')
const DeleteUserAttributesResponse$json = {
  '1': 'DeleteUserAttributesResponse',
};

/// Descriptor for `DeleteUserAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserAttributesResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVVc2VyQXR0cmlidXRlc1Jlc3BvbnNl');

@$core.Deprecated('Use deleteUserPoolClientRequestDescriptor instead')
const DeleteUserPoolClientRequest$json = {
  '1': 'DeleteUserPoolClientRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteUserPoolClientRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPoolClientRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVVc2VyUG9vbENsaWVudFJlcXVlc3QSHgoIY2xpZW50aWQYpOiG1gEgASgJUghjbG'
        'llbnRpZBIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use deleteUserPoolClientSecretRequestDescriptor instead')
const DeleteUserPoolClientSecretRequest$json = {
  '1': 'DeleteUserPoolClientSecretRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientsecretid',
      '3': 51685996,
      '4': 1,
      '5': 9,
      '10': 'clientsecretid'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteUserPoolClientSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPoolClientSecretRequestDescriptor =
    $convert.base64Decode(
        'CiFEZWxldGVVc2VyUG9vbENsaWVudFNlY3JldFJlcXVlc3QSHgoIY2xpZW50aWQYpOiG1gEgAS'
        'gJUghjbGllbnRpZBIpCg5jbGllbnRzZWNyZXRpZBjs1NIYIAEoCVIOY2xpZW50c2VjcmV0aWQS'
        'IgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use deleteUserPoolClientSecretResponseDescriptor instead')
const DeleteUserPoolClientSecretResponse$json = {
  '1': 'DeleteUserPoolClientSecretResponse',
};

/// Descriptor for `DeleteUserPoolClientSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPoolClientSecretResponseDescriptor =
    $convert.base64Decode('CiJEZWxldGVVc2VyUG9vbENsaWVudFNlY3JldFJlc3BvbnNl');

@$core.Deprecated('Use deleteUserPoolDomainRequestDescriptor instead')
const DeleteUserPoolDomainRequest$json = {
  '1': 'DeleteUserPoolDomainRequest',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteUserPoolDomainRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPoolDomainRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVVc2VyUG9vbERvbWFpblJlcXVlc3QSGgoGZG9tYWluGJKS8vABIAEoCVIGZG9tYW'
        'luEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use deleteUserPoolDomainResponseDescriptor instead')
const DeleteUserPoolDomainResponse$json = {
  '1': 'DeleteUserPoolDomainResponse',
};

/// Descriptor for `DeleteUserPoolDomainResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPoolDomainResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVVc2VyUG9vbERvbWFpblJlc3BvbnNl');

@$core.Deprecated('Use deleteUserPoolRequestDescriptor instead')
const DeleteUserPoolRequest$json = {
  '1': 'DeleteUserPoolRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DeleteUserPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserPoolRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVVc2VyUG9vbFJlcXVlc3QSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb2'
    '9saWQ=');

@$core.Deprecated('Use deleteUserRequestDescriptor instead')
const DeleteUserRequest$json = {
  '1': 'DeleteUserRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
  ],
};

/// Descriptor for `DeleteUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteUserRequestDescriptor = $convert.base64Decode(
    'ChFEZWxldGVVc2VyUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIAEoCVILYWNjZXNzdG9rZW'
    '4=');

@$core.Deprecated('Use deleteWebAuthnCredentialRequestDescriptor instead')
const DeleteWebAuthnCredentialRequest$json = {
  '1': 'DeleteWebAuthnCredentialRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'credentialid', '3': 244832166, '4': 1, '5': 9, '10': 'credentialid'},
  ],
};

/// Descriptor for `DeleteWebAuthnCredentialRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteWebAuthnCredentialRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVXZWJBdXRobkNyZWRlbnRpYWxSZXF1ZXN0EiMKC2FjY2Vzc3Rva2VuGIm8kEYgAS'
        'gJUgthY2Nlc3N0b2tlbhIlCgxjcmVkZW50aWFsaWQYpq/fdCABKAlSDGNyZWRlbnRpYWxpZA==');

@$core.Deprecated('Use deleteWebAuthnCredentialResponseDescriptor instead')
const DeleteWebAuthnCredentialResponse$json = {
  '1': 'DeleteWebAuthnCredentialResponse',
};

/// Descriptor for `DeleteWebAuthnCredentialResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteWebAuthnCredentialResponseDescriptor =
    $convert.base64Decode('CiBEZWxldGVXZWJBdXRobkNyZWRlbnRpYWxSZXNwb25zZQ==');

@$core.Deprecated('Use describeIdentityProviderRequestDescriptor instead')
const DescribeIdentityProviderRequest$json = {
  '1': 'DescribeIdentityProviderRequest',
  '2': [
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeIdentityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeIdentityProviderRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZUlkZW50aXR5UHJvdmlkZXJSZXF1ZXN0EiYKDHByb3ZpZGVybmFtZRj4oajnAS'
        'ABKAlSDHByb3ZpZGVybmFtZRIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use describeIdentityProviderResponseDescriptor instead')
const DescribeIdentityProviderResponse$json = {
  '1': 'DescribeIdentityProviderResponse',
  '2': [
    {
      '1': 'identityprovider',
      '3': 450450187,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.IdentityProviderType',
      '10': 'identityprovider'
    },
  ],
};

/// Descriptor for `DescribeIdentityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeIdentityProviderResponseDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZUlkZW50aXR5UHJvdmlkZXJSZXNwb25zZRJRChBpZGVudGl0eXByb3ZpZGVyGI'
        'um5dYBIAEoCzIhLmNvZ25pdG9faWRwLklkZW50aXR5UHJvdmlkZXJUeXBlUhBpZGVudGl0eXBy'
        'b3ZpZGVy');

@$core.Deprecated(
    'Use describeManagedLoginBrandingByClientRequestDescriptor instead')
const DescribeManagedLoginBrandingByClientRequest$json = {
  '1': 'DescribeManagedLoginBrandingByClientRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'returnmergedresources',
      '3': 496753301,
      '4': 1,
      '5': 8,
      '10': 'returnmergedresources'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeManagedLoginBrandingByClientRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeManagedLoginBrandingByClientRequestDescriptor =
    $convert.base64Decode(
        'CitEZXNjcmliZU1hbmFnZWRMb2dpbkJyYW5kaW5nQnlDbGllbnRSZXF1ZXN0Eh4KCGNsaWVudG'
        'lkGKTohtYBIAEoCVIIY2xpZW50aWQSOAoVcmV0dXJubWVyZ2VkcmVzb3VyY2VzGJW17+wBIAEo'
        'CFIVcmV0dXJubWVyZ2VkcmVzb3VyY2VzEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG'
        '9vbGlk');

@$core.Deprecated(
    'Use describeManagedLoginBrandingByClientResponseDescriptor instead')
const DescribeManagedLoginBrandingByClientResponse$json = {
  '1': 'DescribeManagedLoginBrandingByClientResponse',
  '2': [
    {
      '1': 'managedloginbranding',
      '3': 338791109,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ManagedLoginBrandingType',
      '10': 'managedloginbranding'
    },
  ],
};

/// Descriptor for `DescribeManagedLoginBrandingByClientResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeManagedLoginBrandingByClientResponseDescriptor =
    $convert.base64Decode(
        'CixEZXNjcmliZU1hbmFnZWRMb2dpbkJyYW5kaW5nQnlDbGllbnRSZXNwb25zZRJdChRtYW5hZ2'
        'VkbG9naW5icmFuZGluZxjFlcahASABKAsyJS5jb2duaXRvX2lkcC5NYW5hZ2VkTG9naW5CcmFu'
        'ZGluZ1R5cGVSFG1hbmFnZWRsb2dpbmJyYW5kaW5n');

@$core.Deprecated('Use describeManagedLoginBrandingRequestDescriptor instead')
const DescribeManagedLoginBrandingRequest$json = {
  '1': 'DescribeManagedLoginBrandingRequest',
  '2': [
    {
      '1': 'managedloginbrandingid',
      '3': 57829482,
      '4': 1,
      '5': 9,
      '10': 'managedloginbrandingid'
    },
    {
      '1': 'returnmergedresources',
      '3': 496753301,
      '4': 1,
      '5': 8,
      '10': 'returnmergedresources'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeManagedLoginBrandingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeManagedLoginBrandingRequestDescriptor =
    $convert.base64Decode(
        'CiNEZXNjcmliZU1hbmFnZWRMb2dpbkJyYW5kaW5nUmVxdWVzdBI5ChZtYW5hZ2VkbG9naW5icm'
        'FuZGluZ2lkGOrQyRsgASgJUhZtYW5hZ2VkbG9naW5icmFuZGluZ2lkEjgKFXJldHVybm1lcmdl'
        'ZHJlc291cmNlcxiVte/sASABKAhSFXJldHVybm1lcmdlZHJlc291cmNlcxIiCgp1c2VycG9vbG'
        'lkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use describeManagedLoginBrandingResponseDescriptor instead')
const DescribeManagedLoginBrandingResponse$json = {
  '1': 'DescribeManagedLoginBrandingResponse',
  '2': [
    {
      '1': 'managedloginbranding',
      '3': 338791109,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ManagedLoginBrandingType',
      '10': 'managedloginbranding'
    },
  ],
};

/// Descriptor for `DescribeManagedLoginBrandingResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeManagedLoginBrandingResponseDescriptor =
    $convert.base64Decode(
        'CiREZXNjcmliZU1hbmFnZWRMb2dpbkJyYW5kaW5nUmVzcG9uc2USXQoUbWFuYWdlZGxvZ2luYn'
        'JhbmRpbmcYxZXGoQEgASgLMiUuY29nbml0b19pZHAuTWFuYWdlZExvZ2luQnJhbmRpbmdUeXBl'
        'UhRtYW5hZ2VkbG9naW5icmFuZGluZw==');

@$core.Deprecated('Use describeResourceServerRequestDescriptor instead')
const DescribeResourceServerRequest$json = {
  '1': 'DescribeResourceServerRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeResourceServerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeResourceServerRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZVJlc291cmNlU2VydmVyUmVxdWVzdBIhCgppZGVudGlmaWVyGN+g+xMgASgJUg'
        'ppZGVudGlmaWVyEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use describeResourceServerResponseDescriptor instead')
const DescribeResourceServerResponse$json = {
  '1': 'DescribeResourceServerResponse',
  '2': [
    {
      '1': 'resourceserver',
      '3': 506282051,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ResourceServerType',
      '10': 'resourceserver'
    },
  ],
};

/// Descriptor for `DescribeResourceServerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeResourceServerResponseDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZVJlc291cmNlU2VydmVyUmVzcG9uc2USSwoOcmVzb3VyY2VzZXJ2ZXIYw4C18Q'
        'EgASgLMh8uY29nbml0b19pZHAuUmVzb3VyY2VTZXJ2ZXJUeXBlUg5yZXNvdXJjZXNlcnZlcg==');

@$core.Deprecated('Use describeRiskConfigurationRequestDescriptor instead')
const DescribeRiskConfigurationRequest$json = {
  '1': 'DescribeRiskConfigurationRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeRiskConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeRiskConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZVJpc2tDb25maWd1cmF0aW9uUmVxdWVzdBIeCghjbGllbnRpZBik6IbWASABKA'
        'lSCGNsaWVudGlkEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use describeRiskConfigurationResponseDescriptor instead')
const DescribeRiskConfigurationResponse$json = {
  '1': 'DescribeRiskConfigurationResponse',
  '2': [
    {
      '1': 'riskconfiguration',
      '3': 374237465,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.RiskConfigurationType',
      '10': 'riskconfiguration'
    },
  ],
};

/// Descriptor for `DescribeRiskConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeRiskConfigurationResponseDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZVJpc2tDb25maWd1cmF0aW9uUmVzcG9uc2USVAoRcmlza2NvbmZpZ3VyYXRpb2'
        '4YmdK5sgEgASgLMiIuY29nbml0b19pZHAuUmlza0NvbmZpZ3VyYXRpb25UeXBlUhFyaXNrY29u'
        'ZmlndXJhdGlvbg==');

@$core.Deprecated('Use describeTermsRequestDescriptor instead')
const DescribeTermsRequest$json = {
  '1': 'DescribeTermsRequest',
  '2': [
    {'1': 'termsid', '3': 331306210, '4': 1, '5': 9, '10': 'termsid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeTermsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTermsRequestDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZVRlcm1zUmVxdWVzdBIcCgd0ZXJtc2lkGOKp/Z0BIAEoCVIHdGVybXNpZBIiCg'
    'p1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use describeTermsResponseDescriptor instead')
const DescribeTermsResponse$json = {
  '1': 'DescribeTermsResponse',
  '2': [
    {
      '1': 'terms',
      '3': 339062221,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.TermsType',
      '10': 'terms'
    },
  ],
};

/// Descriptor for `DescribeTermsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTermsResponseDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZVRlcm1zUmVzcG9uc2USMAoFdGVybXMYzdvWoQEgASgLMhYuY29nbml0b19pZH'
    'AuVGVybXNUeXBlUgV0ZXJtcw==');

@$core.Deprecated('Use describeUserImportJobRequestDescriptor instead')
const DescribeUserImportJobRequest$json = {
  '1': 'DescribeUserImportJobRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeUserImportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserImportJobRequestDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZVVzZXJJbXBvcnRKb2JSZXF1ZXN0EhcKBWpvYmlkGNLU3TMgASgJUgVqb2JpZB'
        'IiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use describeUserImportJobResponseDescriptor instead')
const DescribeUserImportJobResponse$json = {
  '1': 'DescribeUserImportJobResponse',
  '2': [
    {
      '1': 'userimportjob',
      '3': 473682999,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserImportJobType',
      '10': 'userimportjob'
    },
  ],
};

/// Descriptor for `DescribeUserImportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserImportJobResponseDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZVVzZXJJbXBvcnRKb2JSZXNwb25zZRJICg11c2VyaW1wb3J0am9iGLeo7+EBIA'
        'EoCzIeLmNvZ25pdG9faWRwLlVzZXJJbXBvcnRKb2JUeXBlUg11c2VyaW1wb3J0am9i');

@$core.Deprecated('Use describeUserPoolClientRequestDescriptor instead')
const DescribeUserPoolClientRequest$json = {
  '1': 'DescribeUserPoolClientRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeUserPoolClientRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserPoolClientRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZVVzZXJQb29sQ2xpZW50UmVxdWVzdBIeCghjbGllbnRpZBik6IbWASABKAlSCG'
        'NsaWVudGlkEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use describeUserPoolClientResponseDescriptor instead')
const DescribeUserPoolClientResponse$json = {
  '1': 'DescribeUserPoolClientResponse',
  '2': [
    {
      '1': 'userpoolclient',
      '3': 138497904,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolClientType',
      '10': 'userpoolclient'
    },
  ],
};

/// Descriptor for `DescribeUserPoolClientResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserPoolClientResponseDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZVVzZXJQb29sQ2xpZW50UmVzcG9uc2USSgoOdXNlcnBvb2xjbGllbnQY8J6FQi'
        'ABKAsyHy5jb2duaXRvX2lkcC5Vc2VyUG9vbENsaWVudFR5cGVSDnVzZXJwb29sY2xpZW50');

@$core.Deprecated('Use describeUserPoolDomainRequestDescriptor instead')
const DescribeUserPoolDomainRequest$json = {
  '1': 'DescribeUserPoolDomainRequest',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
  ],
};

/// Descriptor for `DescribeUserPoolDomainRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserPoolDomainRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZVVzZXJQb29sRG9tYWluUmVxdWVzdBIaCgZkb21haW4YkpLy8AEgASgJUgZkb2'
        '1haW4=');

@$core.Deprecated('Use describeUserPoolDomainResponseDescriptor instead')
const DescribeUserPoolDomainResponse$json = {
  '1': 'DescribeUserPoolDomainResponse',
  '2': [
    {
      '1': 'domaindescription',
      '3': 401253238,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.DomainDescriptionType',
      '10': 'domaindescription'
    },
  ],
};

/// Descriptor for `DescribeUserPoolDomainResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserPoolDomainResponseDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZVVzZXJQb29sRG9tYWluUmVzcG9uc2USVAoRZG9tYWluZGVzY3JpcHRpb24Y9s'
        'aqvwEgASgLMiIuY29nbml0b19pZHAuRG9tYWluRGVzY3JpcHRpb25UeXBlUhFkb21haW5kZXNj'
        'cmlwdGlvbg==');

@$core.Deprecated('Use describeUserPoolRequestDescriptor instead')
const DescribeUserPoolRequest$json = {
  '1': 'DescribeUserPoolRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `DescribeUserPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserPoolRequestDescriptor =
    $convert.base64Decode(
        'ChdEZXNjcmliZVVzZXJQb29sUmVxdWVzdBIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcn'
        'Bvb2xpZA==');

@$core.Deprecated('Use describeUserPoolResponseDescriptor instead')
const DescribeUserPoolResponse$json = {
  '1': 'DescribeUserPoolResponse',
  '2': [
    {
      '1': 'userpool',
      '3': 404993697,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolType',
      '10': 'userpool'
    },
  ],
};

/// Descriptor for `DescribeUserPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeUserPoolResponseDescriptor =
    $convert.base64Decode(
        'ChhEZXNjcmliZVVzZXJQb29sUmVzcG9uc2USOQoIdXNlcnBvb2wYoe2OwQEgASgLMhkuY29nbm'
        'l0b19pZHAuVXNlclBvb2xUeXBlUgh1c2VycG9vbA==');

@$core.Deprecated('Use deviceConfigurationTypeDescriptor instead')
const DeviceConfigurationType$json = {
  '1': 'DeviceConfigurationType',
  '2': [
    {
      '1': 'challengerequiredonnewdevice',
      '3': 291794359,
      '4': 1,
      '5': 8,
      '10': 'challengerequiredonnewdevice'
    },
    {
      '1': 'deviceonlyrememberedonuserprompt',
      '3': 199534264,
      '4': 1,
      '5': 8,
      '10': 'deviceonlyrememberedonuserprompt'
    },
  ],
};

/// Descriptor for `DeviceConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceConfigurationTypeDescriptor = $convert.base64Decode(
    'ChdEZXZpY2VDb25maWd1cmF0aW9uVHlwZRJGChxjaGFsbGVuZ2VyZXF1aXJlZG9ubmV3ZGV2aW'
    'NlGLfbkYsBIAEoCFIcY2hhbGxlbmdlcmVxdWlyZWRvbm5ld2RldmljZRJNCiBkZXZpY2Vvbmx5'
    'cmVtZW1iZXJlZG9udXNlcnByb21wdBi4zZJfIAEoCFIgZGV2aWNlb25seXJlbWVtYmVyZWRvbn'
    'VzZXJwcm9tcHQ=');

@$core.Deprecated('Use deviceKeyExistsExceptionDescriptor instead')
const DeviceKeyExistsException$json = {
  '1': 'DeviceKeyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DeviceKeyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceKeyExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChhEZXZpY2VLZXlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use deviceSecretVerifierConfigTypeDescriptor instead')
const DeviceSecretVerifierConfigType$json = {
  '1': 'DeviceSecretVerifierConfigType',
  '2': [
    {
      '1': 'passwordverifier',
      '3': 5891351,
      '4': 1,
      '5': 9,
      '10': 'passwordverifier'
    },
    {'1': 'salt', '3': 37508770, '4': 1, '5': 9, '10': 'salt'},
  ],
};

/// Descriptor for `DeviceSecretVerifierConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceSecretVerifierConfigTypeDescriptor =
    $convert.base64Decode(
        'Ch5EZXZpY2VTZWNyZXRWZXJpZmllckNvbmZpZ1R5cGUSLQoQcGFzc3dvcmR2ZXJpZmllchiXyu'
        'cCIAEoCVIQcGFzc3dvcmR2ZXJpZmllchIVCgRzYWx0GKKt8REgASgJUgRzYWx0');

@$core.Deprecated('Use deviceTypeDescriptor instead')
const DeviceType$json = {
  '1': 'DeviceType',
  '2': [
    {
      '1': 'deviceattributes',
      '3': 204381535,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'deviceattributes'
    },
    {
      '1': 'devicecreatedate',
      '3': 489197016,
      '4': 1,
      '5': 9,
      '10': 'devicecreatedate'
    },
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
    {
      '1': 'devicelastauthenticateddate',
      '3': 63931991,
      '4': 1,
      '5': 9,
      '10': 'devicelastauthenticateddate'
    },
    {
      '1': 'devicelastmodifieddate',
      '3': 508534725,
      '4': 1,
      '5': 9,
      '10': 'devicelastmodifieddate'
    },
  ],
};

/// Descriptor for `DeviceType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceTypeDescriptor = $convert.base64Decode(
    'CgpEZXZpY2VUeXBlEkkKEGRldmljZWF0dHJpYnV0ZXMY37q6YSADKAsyGi5jb2duaXRvX2lkcC'
    '5BdHRyaWJ1dGVUeXBlUhBkZXZpY2VhdHRyaWJ1dGVzEi4KEGRldmljZWNyZWF0ZWRhdGUY2Jui'
    '6QEgASgJUhBkZXZpY2VjcmVhdGVkYXRlEiAKCWRldmljZWtleRir5Mi2ASABKAlSCWRldmljZW'
    'tleRJDChtkZXZpY2VsYXN0YXV0aGVudGljYXRlZGRhdGUY14y+HiABKAlSG2RldmljZWxhc3Rh'
    'dXRoZW50aWNhdGVkZGF0ZRI6ChZkZXZpY2VsYXN0bW9kaWZpZWRkYXRlGMW/vvIBIAEoCVIWZG'
    'V2aWNlbGFzdG1vZGlmaWVkZGF0ZQ==');

@$core.Deprecated('Use domainDescriptionTypeDescriptor instead')
const DomainDescriptionType$json = {
  '1': 'DomainDescriptionType',
  '2': [
    {'1': 'awsaccountid', '3': 370093421, '4': 1, '5': 9, '10': 'awsaccountid'},
    {
      '1': 'cloudfrontdistribution',
      '3': 171268400,
      '4': 1,
      '5': 9,
      '10': 'cloudfrontdistribution'
    },
    {
      '1': 'customdomainconfig',
      '3': 472479421,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CustomDomainConfigType',
      '10': 'customdomainconfig'
    },
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'managedloginversion',
      '3': 479901038,
      '4': 1,
      '5': 5,
      '10': 'managedloginversion'
    },
    {'1': 's3bucket', '3': 114031434, '4': 1, '5': 9, '10': 's3bucket'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DomainStatusType',
      '10': 'status'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `DomainDescriptionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainDescriptionTypeDescriptor = $convert.base64Decode(
    'ChVEb21haW5EZXNjcmlwdGlvblR5cGUSJgoMYXdzYWNjb3VudGlkGO3avLABIAEoCVIMYXdzYW'
    'Njb3VudGlkEjkKFmNsb3VkZnJvbnRkaXN0cmlidXRpb24YsLLVUSABKAlSFmNsb3VkZnJvbnRk'
    'aXN0cmlidXRpb24SVwoSY3VzdG9tZG9tYWluY29uZmlnGL3tpeEBIAEoCzIjLmNvZ25pdG9faW'
    'RwLkN1c3RvbURvbWFpbkNvbmZpZ1R5cGVSEmN1c3RvbWRvbWFpbmNvbmZpZxIaCgZkb21haW4Y'
    'kpLy8AEgASgJUgZkb21haW4SNAoTbWFuYWdlZGxvZ2ludmVyc2lvbhju6urkASABKAVSE21hbm'
    'FnZWRsb2dpbnZlcnNpb24SHQoIczNidWNrZXQYyvavNiABKAlSCHMzYnVja2V0EjgKBnN0YXR1'
    'cxiQ5PsCIAEoDjIdLmNvZ25pdG9faWRwLkRvbWFpblN0YXR1c1R5cGVSBnN0YXR1cxIiCgp1c2'
    'VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZBIcCgd2ZXJzaW9uGLiqt+4BIAEoCVIHdmVy'
    'c2lvbg==');

@$core.Deprecated('Use duplicateProviderExceptionDescriptor instead')
const DuplicateProviderException$json = {
  '1': 'DuplicateProviderException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DuplicateProviderException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List duplicateProviderExceptionDescriptor =
    $convert.base64Decode(
        'ChpEdXBsaWNhdGVQcm92aWRlckV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use emailConfigurationTypeDescriptor instead')
const EmailConfigurationType$json = {
  '1': 'EmailConfigurationType',
  '2': [
    {
      '1': 'configurationset',
      '3': 42178544,
      '4': 1,
      '5': 9,
      '10': 'configurationset'
    },
    {
      '1': 'emailsendingaccount',
      '3': 284286811,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.EmailSendingAccountType',
      '10': 'emailsendingaccount'
    },
    {'1': 'from', '3': 410269078, '4': 1, '5': 9, '10': 'from'},
    {
      '1': 'replytoemailaddress',
      '3': 166144599,
      '4': 1,
      '5': 9,
      '10': 'replytoemailaddress'
    },
    {'1': 'sourcearn', '3': 439903072, '4': 1, '5': 9, '10': 'sourcearn'},
  ],
};

/// Descriptor for `EmailConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailConfigurationTypeDescriptor = $convert.base64Decode(
    'ChZFbWFpbENvbmZpZ3VyYXRpb25UeXBlEi0KEGNvbmZpZ3VyYXRpb25zZXQY8K+OFCABKAlSEG'
    'NvbmZpZ3VyYXRpb25zZXQSWgoTZW1haWxzZW5kaW5nYWNjb3VudBjbvseHASABKA4yJC5jb2du'
    'aXRvX2lkcC5FbWFpbFNlbmRpbmdBY2NvdW50VHlwZVITZW1haWxzZW5kaW5nYWNjb3VudBIWCg'
    'Rmcm9tGJbr0MMBIAEoCVIEZnJvbRIzChNyZXBseXRvZW1haWxhZGRyZXNzGNfUnE8gASgJUhNy'
    'ZXBseXRvZW1haWxhZGRyZXNzEiAKCXNvdXJjZWFybhjgxuHRASABKAlSCXNvdXJjZWFybg==');

@$core.Deprecated('Use emailMfaConfigTypeDescriptor instead')
const EmailMfaConfigType$json = {
  '1': 'EmailMfaConfigType',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
  ],
};

/// Descriptor for `EmailMfaConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailMfaConfigTypeDescriptor = $convert.base64Decode(
    'ChJFbWFpbE1mYUNvbmZpZ1R5cGUSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZRIbCgdzdW'
    'JqZWN0GPDJ5AMgASgJUgdzdWJqZWN0');

@$core.Deprecated('Use emailMfaSettingsTypeDescriptor instead')
const EmailMfaSettingsType$json = {
  '1': 'EmailMfaSettingsType',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'preferredmfa', '3': 195248469, '4': 1, '5': 8, '10': 'preferredmfa'},
  ],
};

/// Descriptor for `EmailMfaSettingsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailMfaSettingsTypeDescriptor = $convert.base64Decode(
    'ChRFbWFpbE1mYVNldHRpbmdzVHlwZRIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIlCg'
    'xwcmVmZXJyZWRtZmEY1YKNXSABKAhSDHByZWZlcnJlZG1mYQ==');

@$core.Deprecated('Use enableSoftwareTokenMFAExceptionDescriptor instead')
const EnableSoftwareTokenMFAException$json = {
  '1': 'EnableSoftwareTokenMFAException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EnableSoftwareTokenMFAException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableSoftwareTokenMFAExceptionDescriptor =
    $convert.base64Decode(
        'Ch9FbmFibGVTb2Z0d2FyZVRva2VuTUZBRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use eventContextDataTypeDescriptor instead')
const EventContextDataType$json = {
  '1': 'EventContextDataType',
  '2': [
    {'1': 'city', '3': 275461731, '4': 1, '5': 9, '10': 'city'},
    {'1': 'country', '3': 83164786, '4': 1, '5': 9, '10': 'country'},
    {'1': 'devicename', '3': 514720633, '4': 1, '5': 9, '10': 'devicename'},
    {'1': 'ipaddress', '3': 1800397, '4': 1, '5': 9, '10': 'ipaddress'},
    {'1': 'timezone', '3': 246302531, '4': 1, '5': 9, '10': 'timezone'},
  ],
};

/// Descriptor for `EventContextDataType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventContextDataTypeDescriptor = $convert.base64Decode(
    'ChRFdmVudENvbnRleHREYXRhVHlwZRIWCgRjaXR5GOPsrIMBIAEoCVIEY2l0eRIbCgdjb3VudH'
    'J5GPL80ycgASgJUgdjb3VudHJ5EiIKCmRldmljZW5hbWUY+Ya49QEgASgJUgpkZXZpY2VuYW1l'
    'Eh4KCWlwYWRkcmVzcxjN8W0gASgJUglpcGFkZHJlc3MSHQoIdGltZXpvbmUYw465dSABKAlSCH'
    'RpbWV6b25l');

@$core.Deprecated('Use eventFeedbackTypeDescriptor instead')
const EventFeedbackType$json = {
  '1': 'EventFeedbackType',
  '2': [
    {'1': 'feedbackdate', '3': 39988103, '4': 1, '5': 9, '10': 'feedbackdate'},
    {
      '1': 'feedbackvalue',
      '3': 259452876,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.FeedbackValueType',
      '10': 'feedbackvalue'
    },
    {'1': 'provider', '3': 363366621, '4': 1, '5': 9, '10': 'provider'},
  ],
};

/// Descriptor for `EventFeedbackType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventFeedbackTypeDescriptor = $convert.base64Decode(
    'ChFFdmVudEZlZWRiYWNrVHlwZRIlCgxmZWVkYmFja2RhdGUYh9eIEyABKAlSDGZlZWRiYWNrZG'
    'F0ZRJHCg1mZWVkYmFja3ZhbHVlGMzf23sgASgOMh4uY29nbml0b19pZHAuRmVlZGJhY2tWYWx1'
    'ZVR5cGVSDWZlZWRiYWNrdmFsdWUSHgoIcHJvdmlkZXIY3ZGirQEgASgJUghwcm92aWRlcg==');

@$core.Deprecated('Use eventRiskTypeDescriptor instead')
const EventRiskType$json = {
  '1': 'EventRiskType',
  '2': [
    {
      '1': 'compromisedcredentialsdetected',
      '3': 251099316,
      '4': 1,
      '5': 8,
      '10': 'compromisedcredentialsdetected'
    },
    {
      '1': 'riskdecision',
      '3': 476551633,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.RiskDecisionType',
      '10': 'riskdecision'
    },
    {
      '1': 'risklevel',
      '3': 212098729,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.RiskLevelType',
      '10': 'risklevel'
    },
  ],
};

/// Descriptor for `EventRiskType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventRiskTypeDescriptor = $convert.base64Decode(
    'Cg1FdmVudFJpc2tUeXBlEkkKHmNvbXByb21pc2VkY3JlZGVudGlhbHNkZXRlY3RlZBi08d13IA'
    'EoCFIeY29tcHJvbWlzZWRjcmVkZW50aWFsc2RldGVjdGVkEkUKDHJpc2tkZWNpc2lvbhjRs57j'
    'ASABKA4yHS5jb2duaXRvX2lkcC5SaXNrRGVjaXNpb25UeXBlUgxyaXNrZGVjaXNpb24SOwoJcm'
    'lza2xldmVsGKm9kWUgASgOMhouY29nbml0b19pZHAuUmlza0xldmVsVHlwZVIJcmlza2xldmVs');

@$core.Deprecated('Use expiredCodeExceptionDescriptor instead')
const ExpiredCodeException$json = {
  '1': 'ExpiredCodeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExpiredCodeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expiredCodeExceptionDescriptor =
    $convert.base64Decode(
        'ChRFeHBpcmVkQ29kZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use featureUnavailableInTierExceptionDescriptor instead')
const FeatureUnavailableInTierException$json = {
  '1': 'FeatureUnavailableInTierException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FeatureUnavailableInTierException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List featureUnavailableInTierExceptionDescriptor =
    $convert.base64Decode(
        'CiFGZWF0dXJlVW5hdmFpbGFibGVJblRpZXJFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use firehoseConfigurationTypeDescriptor instead')
const FirehoseConfigurationType$json = {
  '1': 'FirehoseConfigurationType',
  '2': [
    {'1': 'streamarn', '3': 513423709, '4': 1, '5': 9, '10': 'streamarn'},
  ],
};

/// Descriptor for `FirehoseConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List firehoseConfigurationTypeDescriptor =
    $convert.base64Decode(
        'ChlGaXJlaG9zZUNvbmZpZ3VyYXRpb25UeXBlEiAKCXN0cmVhbWFybhjd8uj0ASABKAlSCXN0cm'
        'VhbWFybg==');

@$core.Deprecated('Use forbiddenExceptionDescriptor instead')
const ForbiddenException$json = {
  '1': 'ForbiddenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ForbiddenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List forbiddenExceptionDescriptor =
    $convert.base64Decode(
        'ChJGb3JiaWRkZW5FeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use forgetDeviceRequestDescriptor instead')
const ForgetDeviceRequest$json = {
  '1': 'ForgetDeviceRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
  ],
};

/// Descriptor for `ForgetDeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List forgetDeviceRequestDescriptor = $convert.base64Decode(
    'ChNGb3JnZXREZXZpY2VSZXF1ZXN0EiMKC2FjY2Vzc3Rva2VuGIm8kEYgASgJUgthY2Nlc3N0b2'
    'tlbhIgCglkZXZpY2VrZXkYq+TItgEgASgJUglkZXZpY2VrZXk=');

@$core.Deprecated('Use forgotPasswordRequestDescriptor instead')
const ForgotPasswordRequest$json = {
  '1': 'ForgotPasswordRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ForgotPasswordRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'secrethash', '3': 321025630, '4': 1, '5': 9, '10': 'secrethash'},
    {
      '1': 'usercontextdata',
      '3': 169951134,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserContextDataType',
      '10': 'usercontextdata'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [ForgotPasswordRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use forgotPasswordRequestDescriptor instead')
const ForgotPasswordRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ForgotPasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List forgotPasswordRequestDescriptor = $convert.base64Decode(
    'ChVGb3Jnb3RQYXNzd29yZFJlcXVlc3QSUwoRYW5hbHl0aWNzbWV0YWRhdGEY94XBPCABKAsyIi'
    '5jb2duaXRvX2lkcC5BbmFseXRpY3NNZXRhZGF0YVR5cGVSEWFuYWx5dGljc21ldGFkYXRhEh4K'
    'CGNsaWVudGlkGKTohtYBIAEoCVIIY2xpZW50aWQSYQoOY2xpZW50bWV0YWRhdGEYzK//YSADKA'
    'syNi5jb2duaXRvX2lkcC5Gb3Jnb3RQYXNzd29yZFJlcXVlc3QuQ2xpZW50bWV0YWRhdGFFbnRy'
    'eVIOY2xpZW50bWV0YWRhdGESIgoKc2VjcmV0aGFzaBje7ImZASABKAlSCnNlY3JldGhhc2gSTQ'
    'oPdXNlcmNvbnRleHRkYXRhGJ7/hFEgASgLMiAuY29nbml0b19pZHAuVXNlckNvbnRleHREYXRh'
    'VHlwZVIPdXNlcmNvbnRleHRkYXRhEh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWUaQQ'
    'oTQ2xpZW50bWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIF'
    'dmFsdWU6AjgB');

@$core.Deprecated('Use forgotPasswordResponseDescriptor instead')
const ForgotPasswordResponse$json = {
  '1': 'ForgotPasswordResponse',
  '2': [
    {
      '1': 'codedeliverydetails',
      '3': 423272803,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CodeDeliveryDetailsType',
      '10': 'codedeliverydetails'
    },
  ],
};

/// Descriptor for `ForgotPasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List forgotPasswordResponseDescriptor = $convert.base64Decode(
    'ChZGb3Jnb3RQYXNzd29yZFJlc3BvbnNlEloKE2NvZGVkZWxpdmVyeWRldGFpbHMY48LqyQEgAS'
    'gLMiQuY29nbml0b19pZHAuQ29kZURlbGl2ZXJ5RGV0YWlsc1R5cGVSE2NvZGVkZWxpdmVyeWRl'
    'dGFpbHM=');

@$core.Deprecated('Use getCSVHeaderRequestDescriptor instead')
const GetCSVHeaderRequest$json = {
  '1': 'GetCSVHeaderRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetCSVHeaderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCSVHeaderRequestDescriptor = $convert.base64Decode(
    'ChNHZXRDU1ZIZWFkZXJSZXF1ZXN0EiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbG'
    'lk');

@$core.Deprecated('Use getCSVHeaderResponseDescriptor instead')
const GetCSVHeaderResponse$json = {
  '1': 'GetCSVHeaderResponse',
  '2': [
    {'1': 'csvheader', '3': 294305815, '4': 3, '5': 9, '10': 'csvheader'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetCSVHeaderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCSVHeaderResponseDescriptor = $convert.base64Decode(
    'ChRHZXRDU1ZIZWFkZXJSZXNwb25zZRIgCgljc3ZoZWFkZXIYl4CrjAEgAygJUgljc3ZoZWFkZX'
    'ISIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use getDeviceRequestDescriptor instead')
const GetDeviceRequest$json = {
  '1': 'GetDeviceRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
  ],
};

/// Descriptor for `GetDeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeviceRequestDescriptor = $convert.base64Decode(
    'ChBHZXREZXZpY2VSZXF1ZXN0EiMKC2FjY2Vzc3Rva2VuGIm8kEYgASgJUgthY2Nlc3N0b2tlbh'
    'IgCglkZXZpY2VrZXkYq+TItgEgASgJUglkZXZpY2VrZXk=');

@$core.Deprecated('Use getDeviceResponseDescriptor instead')
const GetDeviceResponse$json = {
  '1': 'GetDeviceResponse',
  '2': [
    {
      '1': 'device',
      '3': 499889172,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.DeviceType',
      '10': 'device'
    },
  ],
};

/// Descriptor for `GetDeviceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeviceResponseDescriptor = $convert.base64Decode(
    'ChFHZXREZXZpY2VSZXNwb25zZRIzCgZkZXZpY2UYlOiu7gEgASgLMhcuY29nbml0b19pZHAuRG'
    'V2aWNlVHlwZVIGZGV2aWNl');

@$core.Deprecated('Use getGroupRequestDescriptor instead')
const GetGroupRequest$json = {
  '1': 'GetGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGroupRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRHcm91cFJlcXVlc3QSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1lEiIKCn'
    'VzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use getGroupResponseDescriptor instead')
const GetGroupResponse$json = {
  '1': 'GetGroupResponse',
  '2': [
    {
      '1': 'group',
      '3': 91525165,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.GroupType',
      '10': 'group'
    },
  ],
};

/// Descriptor for `GetGroupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getGroupResponseDescriptor = $convert.base64Decode(
    'ChBHZXRHcm91cFJlc3BvbnNlEi8KBWdyb3VwGK2g0isgASgLMhYuY29nbml0b19pZHAuR3JvdX'
    'BUeXBlUgVncm91cA==');

@$core
    .Deprecated('Use getIdentityProviderByIdentifierRequestDescriptor instead')
const GetIdentityProviderByIdentifierRequest$json = {
  '1': 'GetIdentityProviderByIdentifierRequest',
  '2': [
    {
      '1': 'idpidentifier',
      '3': 17161156,
      '4': 1,
      '5': 9,
      '10': 'idpidentifier'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetIdentityProviderByIdentifierRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIdentityProviderByIdentifierRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXRJZGVudGl0eVByb3ZpZGVyQnlJZGVudGlmaWVyUmVxdWVzdBInCg1pZHBpZGVudGlmaW'
        'VyGMS3lwggASgJUg1pZHBpZGVudGlmaWVyEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2Vy'
        'cG9vbGlk');

@$core
    .Deprecated('Use getIdentityProviderByIdentifierResponseDescriptor instead')
const GetIdentityProviderByIdentifierResponse$json = {
  '1': 'GetIdentityProviderByIdentifierResponse',
  '2': [
    {
      '1': 'identityprovider',
      '3': 450450187,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.IdentityProviderType',
      '10': 'identityprovider'
    },
  ],
};

/// Descriptor for `GetIdentityProviderByIdentifierResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIdentityProviderByIdentifierResponseDescriptor =
    $convert.base64Decode(
        'CidHZXRJZGVudGl0eVByb3ZpZGVyQnlJZGVudGlmaWVyUmVzcG9uc2USUQoQaWRlbnRpdHlwcm'
        '92aWRlchiLpuXWASABKAsyIS5jb2duaXRvX2lkcC5JZGVudGl0eVByb3ZpZGVyVHlwZVIQaWRl'
        'bnRpdHlwcm92aWRlcg==');

@$core.Deprecated('Use getLogDeliveryConfigurationRequestDescriptor instead')
const GetLogDeliveryConfigurationRequest$json = {
  '1': 'GetLogDeliveryConfigurationRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetLogDeliveryConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogDeliveryConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiJHZXRMb2dEZWxpdmVyeUNvbmZpZ3VyYXRpb25SZXF1ZXN0EiIKCnVzZXJwb29saWQY/saLnQ'
        'EgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use getLogDeliveryConfigurationResponseDescriptor instead')
const GetLogDeliveryConfigurationResponse$json = {
  '1': 'GetLogDeliveryConfigurationResponse',
  '2': [
    {
      '1': 'logdeliveryconfiguration',
      '3': 332495878,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.LogDeliveryConfigurationType',
      '10': 'logdeliveryconfiguration'
    },
  ],
};

/// Descriptor for `GetLogDeliveryConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogDeliveryConfigurationResponseDescriptor =
    $convert.base64Decode(
        'CiNHZXRMb2dEZWxpdmVyeUNvbmZpZ3VyYXRpb25SZXNwb25zZRJpChhsb2dkZWxpdmVyeWNvbm'
        'ZpZ3VyYXRpb24YhvjFngEgASgLMikuY29nbml0b19pZHAuTG9nRGVsaXZlcnlDb25maWd1cmF0'
        'aW9uVHlwZVIYbG9nZGVsaXZlcnljb25maWd1cmF0aW9u');

@$core.Deprecated('Use getSigningCertificateRequestDescriptor instead')
const GetSigningCertificateRequest$json = {
  '1': 'GetSigningCertificateRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetSigningCertificateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSigningCertificateRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRTaWduaW5nQ2VydGlmaWNhdGVSZXF1ZXN0EiIKCnVzZXJwb29saWQY/saLnQEgASgJUg'
        'p1c2VycG9vbGlk');

@$core.Deprecated('Use getSigningCertificateResponseDescriptor instead')
const GetSigningCertificateResponse$json = {
  '1': 'GetSigningCertificateResponse',
  '2': [
    {'1': 'certificate', '3': 198060817, '4': 1, '5': 9, '10': 'certificate'},
  ],
};

/// Descriptor for `GetSigningCertificateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSigningCertificateResponseDescriptor =
    $convert.base64Decode(
        'Ch1HZXRTaWduaW5nQ2VydGlmaWNhdGVSZXNwb25zZRIjCgtjZXJ0aWZpY2F0ZRiR1rheIAEoCV'
        'ILY2VydGlmaWNhdGU=');

@$core.Deprecated('Use getTokensFromRefreshTokenRequestDescriptor instead')
const GetTokensFromRefreshTokenRequest$json = {
  '1': 'GetTokensFromRefreshTokenRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.GetTokensFromRefreshTokenRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'clientsecret', '3': 500734711, '4': 1, '5': 9, '10': 'clientsecret'},
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
    {'1': 'refreshtoken', '3': 253777778, '4': 1, '5': 9, '10': 'refreshtoken'},
  ],
  '3': [GetTokensFromRefreshTokenRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use getTokensFromRefreshTokenRequestDescriptor instead')
const GetTokensFromRefreshTokenRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetTokensFromRefreshTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTokensFromRefreshTokenRequestDescriptor = $convert.base64Decode(
    'CiBHZXRUb2tlbnNGcm9tUmVmcmVzaFRva2VuUmVxdWVzdBIeCghjbGllbnRpZBik6IbWASABKA'
    'lSCGNsaWVudGlkEmwKDmNsaWVudG1ldGFkYXRhGMyv/2EgAygLMkEuY29nbml0b19pZHAuR2V0'
    'VG9rZW5zRnJvbVJlZnJlc2hUb2tlblJlcXVlc3QuQ2xpZW50bWV0YWRhdGFFbnRyeVIOY2xpZW'
    '50bWV0YWRhdGESJgoMY2xpZW50c2VjcmV0GPe14u4BIAEoCVIMY2xpZW50c2VjcmV0EiAKCWRl'
    'dmljZWtleRir5Mi2ASABKAlSCWRldmljZWtleRIlCgxyZWZyZXNodG9rZW4Y8q6BeSABKAlSDH'
    'JlZnJlc2h0b2tlbhpBChNDbGllbnRtZXRhZGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQK'
    'BXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getTokensFromRefreshTokenResponseDescriptor instead')
const GetTokensFromRefreshTokenResponse$json = {
  '1': 'GetTokensFromRefreshTokenResponse',
  '2': [
    {
      '1': 'authenticationresult',
      '3': 519327313,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AuthenticationResultType',
      '10': 'authenticationresult'
    },
  ],
};

/// Descriptor for `GetTokensFromRefreshTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTokensFromRefreshTokenResponseDescriptor =
    $convert.base64Decode(
        'CiFHZXRUb2tlbnNGcm9tUmVmcmVzaFRva2VuUmVzcG9uc2USXQoUYXV0aGVudGljYXRpb25yZX'
        'N1bHQY0ZzR9wEgASgLMiUuY29nbml0b19pZHAuQXV0aGVudGljYXRpb25SZXN1bHRUeXBlUhRh'
        'dXRoZW50aWNhdGlvbnJlc3VsdA==');

@$core.Deprecated('Use getUICustomizationRequestDescriptor instead')
const GetUICustomizationRequest$json = {
  '1': 'GetUICustomizationRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetUICustomizationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUICustomizationRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRVSUN1c3RvbWl6YXRpb25SZXF1ZXN0Eh4KCGNsaWVudGlkGKTohtYBIAEoCVIIY2xpZW'
        '50aWQSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use getUICustomizationResponseDescriptor instead')
const GetUICustomizationResponse$json = {
  '1': 'GetUICustomizationResponse',
  '2': [
    {
      '1': 'uicustomization',
      '3': 28966753,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UICustomizationType',
      '10': 'uicustomization'
    },
  ],
};

/// Descriptor for `GetUICustomizationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUICustomizationResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRVSUN1c3RvbWl6YXRpb25SZXNwb25zZRJNCg91aWN1c3RvbWl6YXRpb24Y4f7nDSABKA'
        'syIC5jb2duaXRvX2lkcC5VSUN1c3RvbWl6YXRpb25UeXBlUg91aWN1c3RvbWl6YXRpb24=');

@$core
    .Deprecated('Use getUserAttributeVerificationCodeRequestDescriptor instead')
const GetUserAttributeVerificationCodeRequest$json = {
  '1': 'GetUserAttributeVerificationCodeRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_idp.GetUserAttributeVerificationCodeRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
  ],
  '3': [GetUserAttributeVerificationCodeRequest_ClientmetadataEntry$json],
};

@$core
    .Deprecated('Use getUserAttributeVerificationCodeRequestDescriptor instead')
const GetUserAttributeVerificationCodeRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetUserAttributeVerificationCodeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserAttributeVerificationCodeRequestDescriptor =
    $convert.base64Decode(
        'CidHZXRVc2VyQXR0cmlidXRlVmVyaWZpY2F0aW9uQ29kZVJlcXVlc3QSIwoLYWNjZXNzdG9rZW'
        '4YibyQRiABKAlSC2FjY2Vzc3Rva2VuEigKDWF0dHJpYnV0ZW5hbWUYrZWYqAEgASgJUg1hdHRy'
        'aWJ1dGVuYW1lEnMKDmNsaWVudG1ldGFkYXRhGMyv/2EgAygLMkguY29nbml0b19pZHAuR2V0VX'
        'NlckF0dHJpYnV0ZVZlcmlmaWNhdGlvbkNvZGVSZXF1ZXN0LkNsaWVudG1ldGFkYXRhRW50cnlS'
        'DmNsaWVudG1ldGFkYXRhGkEKE0NsaWVudG1ldGFkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZX'
        'kSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated(
    'Use getUserAttributeVerificationCodeResponseDescriptor instead')
const GetUserAttributeVerificationCodeResponse$json = {
  '1': 'GetUserAttributeVerificationCodeResponse',
  '2': [
    {
      '1': 'codedeliverydetails',
      '3': 423272803,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CodeDeliveryDetailsType',
      '10': 'codedeliverydetails'
    },
  ],
};

/// Descriptor for `GetUserAttributeVerificationCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserAttributeVerificationCodeResponseDescriptor =
    $convert.base64Decode(
        'CihHZXRVc2VyQXR0cmlidXRlVmVyaWZpY2F0aW9uQ29kZVJlc3BvbnNlEloKE2NvZGVkZWxpdm'
        'VyeWRldGFpbHMY48LqyQEgASgLMiQuY29nbml0b19pZHAuQ29kZURlbGl2ZXJ5RGV0YWlsc1R5'
        'cGVSE2NvZGVkZWxpdmVyeWRldGFpbHM=');

@$core.Deprecated('Use getUserAuthFactorsRequestDescriptor instead')
const GetUserAuthFactorsRequest$json = {
  '1': 'GetUserAuthFactorsRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
  ],
};

/// Descriptor for `GetUserAuthFactorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserAuthFactorsRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRVc2VyQXV0aEZhY3RvcnNSZXF1ZXN0EiMKC2FjY2Vzc3Rva2VuGIm8kEYgASgJUgthY2'
        'Nlc3N0b2tlbg==');

@$core.Deprecated('Use getUserAuthFactorsResponseDescriptor instead')
const GetUserAuthFactorsResponse$json = {
  '1': 'GetUserAuthFactorsResponse',
  '2': [
    {
      '1': 'configureduserauthfactors',
      '3': 207621521,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.AuthFactorType',
      '10': 'configureduserauthfactors'
    },
    {
      '1': 'preferredmfasetting',
      '3': 111271921,
      '4': 1,
      '5': 9,
      '10': 'preferredmfasetting'
    },
    {
      '1': 'usermfasettinglist',
      '3': 230885,
      '4': 3,
      '5': 9,
      '10': 'usermfasettinglist'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetUserAuthFactorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserAuthFactorsResponseDescriptor = $convert.base64Decode(
    'ChpHZXRVc2VyQXV0aEZhY3RvcnNSZXNwb25zZRJcChljb25maWd1cmVkdXNlcmF1dGhmYWN0b3'
    'JzGJGbgGMgAygOMhsuY29nbml0b19pZHAuQXV0aEZhY3RvclR5cGVSGWNvbmZpZ3VyZWR1c2Vy'
    'YXV0aGZhY3RvcnMSMwoTcHJlZmVycmVkbWZhc2V0dGluZxjxv4c1IAEoCVITcHJlZmVycmVkbW'
    'Zhc2V0dGluZxIwChJ1c2VybWZhc2V0dGluZ2xpc3QY5YsOIAMoCVISdXNlcm1mYXNldHRpbmds'
    'aXN0Eh4KCHVzZXJuYW1lGNqpo+ABIAEoCVIIdXNlcm5hbWU=');

@$core.Deprecated('Use getUserPoolMfaConfigRequestDescriptor instead')
const GetUserPoolMfaConfigRequest$json = {
  '1': 'GetUserPoolMfaConfigRequest',
  '2': [
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GetUserPoolMfaConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserPoolMfaConfigRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRVc2VyUG9vbE1mYUNvbmZpZ1JlcXVlc3QSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCn'
        'VzZXJwb29saWQ=');

@$core.Deprecated('Use getUserPoolMfaConfigResponseDescriptor instead')
const GetUserPoolMfaConfigResponse$json = {
  '1': 'GetUserPoolMfaConfigResponse',
  '2': [
    {
      '1': 'emailmfaconfiguration',
      '3': 482754548,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailMfaConfigType',
      '10': 'emailmfaconfiguration'
    },
    {
      '1': 'mfaconfiguration',
      '3': 259020766,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolMfaType',
      '10': 'mfaconfiguration'
    },
    {
      '1': 'smsmfaconfiguration',
      '3': 153073099,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SmsMfaConfigType',
      '10': 'smsmfaconfiguration'
    },
    {
      '1': 'softwaretokenmfaconfiguration',
      '3': 502085950,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SoftwareTokenMfaConfigType',
      '10': 'softwaretokenmfaconfiguration'
    },
    {
      '1': 'webauthnconfiguration',
      '3': 506289104,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.WebAuthnConfigurationType',
      '10': 'webauthnconfiguration'
    },
  ],
};

/// Descriptor for `GetUserPoolMfaConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserPoolMfaConfigResponseDescriptor = $convert.base64Decode(
    'ChxHZXRVc2VyUG9vbE1mYUNvbmZpZ1Jlc3BvbnNlElkKFWVtYWlsbWZhY29uZmlndXJhdGlvbh'
    'j0/5jmASABKAsyHy5jb2duaXRvX2lkcC5FbWFpbE1mYUNvbmZpZ1R5cGVSFWVtYWlsbWZhY29u'
    'ZmlndXJhdGlvbhJLChBtZmFjb25maWd1cmF0aW9uGN6vwXsgASgOMhwuY29nbml0b19pZHAuVX'
    'NlclBvb2xNZmFUeXBlUhBtZmFjb25maWd1cmF0aW9uElIKE3Ntc21mYWNvbmZpZ3VyYXRpb24Y'
    'y+v+SCABKAsyHS5jb2duaXRvX2lkcC5TbXNNZmFDb25maWdUeXBlUhNzbXNtZmFjb25maWd1cm'
    'F0aW9uEnEKHXNvZnR3YXJldG9rZW5tZmFjb25maWd1cmF0aW9uGL7ytO8BIAEoCzInLmNvZ25p'
    'dG9faWRwLlNvZnR3YXJlVG9rZW5NZmFDb25maWdUeXBlUh1zb2Z0d2FyZXRva2VubWZhY29uZm'
    'lndXJhdGlvbhJgChV3ZWJhdXRobmNvbmZpZ3VyYXRpb24Y0Le18QEgASgLMiYuY29nbml0b19p'
    'ZHAuV2ViQXV0aG5Db25maWd1cmF0aW9uVHlwZVIVd2ViYXV0aG5jb25maWd1cmF0aW9u');

@$core.Deprecated('Use getUserRequestDescriptor instead')
const GetUserRequest$json = {
  '1': 'GetUserRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
  ],
};

/// Descriptor for `GetUserRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserRequestDescriptor = $convert.base64Decode(
    'Cg5HZXRVc2VyUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIAEoCVILYWNjZXNzdG9rZW4=');

@$core.Deprecated('Use getUserResponseDescriptor instead')
const GetUserResponse$json = {
  '1': 'GetUserResponse',
  '2': [
    {
      '1': 'mfaoptions',
      '3': 501540826,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.MFAOptionType',
      '10': 'mfaoptions'
    },
    {
      '1': 'preferredmfasetting',
      '3': 111271921,
      '4': 1,
      '5': 9,
      '10': 'preferredmfasetting'
    },
    {
      '1': 'userattributes',
      '3': 194667064,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'userattributes'
    },
    {
      '1': 'usermfasettinglist',
      '3': 230885,
      '4': 3,
      '5': 9,
      '10': 'usermfasettinglist'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `GetUserResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUserResponseDescriptor = $convert.base64Decode(
    'Cg9HZXRVc2VyUmVzcG9uc2USPgoKbWZhb3B0aW9ucxjaz5PvASADKAsyGi5jb2duaXRvX2lkcC'
    '5NRkFPcHRpb25UeXBlUgptZmFvcHRpb25zEjMKE3ByZWZlcnJlZG1mYXNldHRpbmcY8b+HNSAB'
    'KAlSE3ByZWZlcnJlZG1mYXNldHRpbmcSRQoOdXNlcmF0dHJpYnV0ZXMYuMTpXCADKAsyGi5jb2'
    'duaXRvX2lkcC5BdHRyaWJ1dGVUeXBlUg51c2VyYXR0cmlidXRlcxIwChJ1c2VybWZhc2V0dGlu'
    'Z2xpc3QY5YsOIAMoCVISdXNlcm1mYXNldHRpbmdsaXN0Eh4KCHVzZXJuYW1lGNqpo+ABIAEoCV'
    'IIdXNlcm5hbWU=');

@$core.Deprecated('Use globalSignOutRequestDescriptor instead')
const GlobalSignOutRequest$json = {
  '1': 'GlobalSignOutRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
  ],
};

/// Descriptor for `GlobalSignOutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalSignOutRequestDescriptor = $convert.base64Decode(
    'ChRHbG9iYWxTaWduT3V0UmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIAEoCVILYWNjZXNzdG'
    '9rZW4=');

@$core.Deprecated('Use globalSignOutResponseDescriptor instead')
const GlobalSignOutResponse$json = {
  '1': 'GlobalSignOutResponse',
};

/// Descriptor for `GlobalSignOutResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalSignOutResponseDescriptor =
    $convert.base64Decode('ChVHbG9iYWxTaWduT3V0UmVzcG9uc2U=');

@$core.Deprecated('Use groupExistsExceptionDescriptor instead')
const GroupExistsException$json = {
  '1': 'GroupExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `GroupExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List groupExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChRHcm91cEV4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use groupTypeDescriptor instead')
const GroupType$json = {
  '1': 'GroupType',
  '2': [
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'precedence', '3': 206704142, '4': 1, '5': 5, '10': 'precedence'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `GroupType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List groupTypeDescriptor = $convert.base64Decode(
    'CglHcm91cFR5cGUSJgoMY3JlYXRpb25kYXRlGOHYt4kBIAEoCVIMY3JlYXRpb25kYXRlEiMKC2'
    'Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIgCglncm91cG5hbWUYyMqgqgEgASgJ'
    'Uglncm91cG5hbWUSLQoQbGFzdG1vZGlmaWVkZGF0ZRjTiMgLIAEoCVIQbGFzdG1vZGlmaWVkZG'
    'F0ZRIhCgpwcmVjZWRlbmNlGI6cyGIgASgFUgpwcmVjZWRlbmNlEhwKB3JvbGVhcm4YgfjnmQEg'
    'ASgJUgdyb2xlYXJuEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use httpHeaderDescriptor instead')
const HttpHeader$json = {
  '1': 'HttpHeader',
  '2': [
    {'1': 'headername', '3': 409221572, '4': 1, '5': 9, '10': 'headername'},
    {'1': 'headervalue', '3': 13503826, '4': 1, '5': 9, '10': 'headervalue'},
  ],
};

/// Descriptor for `HttpHeader`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List httpHeaderDescriptor = $convert.base64Decode(
    'CgpIdHRwSGVhZGVyEiIKCmhlYWRlcm5hbWUYxPOQwwEgASgJUgpoZWFkZXJuYW1lEiMKC2hlYW'
    'RlcnZhbHVlGNKauAYgASgJUgtoZWFkZXJ2YWx1ZQ==');

@$core.Deprecated('Use identityProviderTypeDescriptor instead')
const IdentityProviderType$json = {
  '1': 'IdentityProviderType',
  '2': [
    {
      '1': 'attributemapping',
      '3': 116923092,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.IdentityProviderType.AttributemappingEntry',
      '10': 'attributemapping'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'idpidentifiers',
      '3': 205051409,
      '4': 3,
      '5': 9,
      '10': 'idpidentifiers'
    },
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'providerdetails',
      '3': 476397115,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.IdentityProviderType.ProviderdetailsEntry',
      '10': 'providerdetails'
    },
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
    {
      '1': 'providertype',
      '3': 337296789,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.IdentityProviderTypeType',
      '10': 'providertype'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [
    IdentityProviderType_AttributemappingEntry$json,
    IdentityProviderType_ProviderdetailsEntry$json
  ],
};

@$core.Deprecated('Use identityProviderTypeDescriptor instead')
const IdentityProviderType_AttributemappingEntry$json = {
  '1': 'AttributemappingEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use identityProviderTypeDescriptor instead')
const IdentityProviderType_ProviderdetailsEntry$json = {
  '1': 'ProviderdetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `IdentityProviderType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List identityProviderTypeDescriptor = $convert.base64Decode(
    'ChRJZGVudGl0eVByb3ZpZGVyVHlwZRJmChBhdHRyaWJ1dGVtYXBwaW5nGNS14DcgAygLMjcuY2'
    '9nbml0b19pZHAuSWRlbnRpdHlQcm92aWRlclR5cGUuQXR0cmlidXRlbWFwcGluZ0VudHJ5UhBh'
    'dHRyaWJ1dGVtYXBwaW5nEiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9uZGF0ZR'
    'IpCg5pZHBpZGVudGlmaWVycxiRrONhIAMoCVIOaWRwaWRlbnRpZmllcnMSLQoQbGFzdG1vZGlm'
    'aWVkZGF0ZRjTiMgLIAEoCVIQbGFzdG1vZGlmaWVkZGF0ZRJkCg9wcm92aWRlcmRldGFpbHMYu/'
    'yU4wEgAygLMjYuY29nbml0b19pZHAuSWRlbnRpdHlQcm92aWRlclR5cGUuUHJvdmlkZXJkZXRh'
    'aWxzRW50cnlSD3Byb3ZpZGVyZGV0YWlscxImCgxwcm92aWRlcm5hbWUY+KGo5wEgASgJUgxwcm'
    '92aWRlcm5hbWUSTQoMcHJvdmlkZXJ0eXBlGJX76qABIAEoDjIlLmNvZ25pdG9faWRwLklkZW50'
    'aXR5UHJvdmlkZXJUeXBlVHlwZVIMcHJvdmlkZXJ0eXBlEiIKCnVzZXJwb29saWQY/saLnQEgAS'
    'gJUgp1c2VycG9vbGlkGkMKFUF0dHJpYnV0ZW1hcHBpbmdFbnRyeRIQCgNrZXkYASABKAlSA2tl'
    'eRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGkIKFFByb3ZpZGVyZGV0YWlsc0VudHJ5EhAKA2'
    'tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use inboundFederationLambdaTypeDescriptor instead')
const InboundFederationLambdaType$json = {
  '1': 'InboundFederationLambdaType',
  '2': [
    {'1': 'lambdaarn', '3': 196258582, '4': 1, '5': 9, '10': 'lambdaarn'},
    {
      '1': 'lambdaversion',
      '3': 429236631,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.InboundFederationLambdaVersionType',
      '10': 'lambdaversion'
    },
  ],
};

/// Descriptor for `InboundFederationLambdaType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inboundFederationLambdaTypeDescriptor =
    $convert.base64Decode(
        'ChtJbmJvdW5kRmVkZXJhdGlvbkxhbWJkYVR5cGUSHwoJbGFtYmRhYXJuGJbWyl0gASgJUglsYW'
        '1iZGFhcm4SWQoNbGFtYmRhdmVyc2lvbhiXw9bMASABKA4yLy5jb2duaXRvX2lkcC5JbmJvdW5k'
        'RmVkZXJhdGlvbkxhbWJkYVZlcnNpb25UeXBlUg1sYW1iZGF2ZXJzaW9u');

@$core.Deprecated('Use initiateAuthRequestDescriptor instead')
const InitiateAuthRequest$json = {
  '1': 'InitiateAuthRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {
      '1': 'authflow',
      '3': 143701420,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AuthFlowType',
      '10': 'authflow'
    },
    {
      '1': 'authparameters',
      '3': 258276552,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.InitiateAuthRequest.AuthparametersEntry',
      '10': 'authparameters'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.InitiateAuthRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {
      '1': 'usercontextdata',
      '3': 169951134,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserContextDataType',
      '10': 'usercontextdata'
    },
  ],
  '3': [
    InitiateAuthRequest_AuthparametersEntry$json,
    InitiateAuthRequest_ClientmetadataEntry$json
  ],
};

@$core.Deprecated('Use initiateAuthRequestDescriptor instead')
const InitiateAuthRequest_AuthparametersEntry$json = {
  '1': 'AuthparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use initiateAuthRequestDescriptor instead')
const InitiateAuthRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `InitiateAuthRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List initiateAuthRequestDescriptor = $convert.base64Decode(
    'ChNJbml0aWF0ZUF1dGhSZXF1ZXN0ElMKEWFuYWx5dGljc21ldGFkYXRhGPeFwTwgASgLMiIuY2'
    '9nbml0b19pZHAuQW5hbHl0aWNzTWV0YWRhdGFUeXBlUhFhbmFseXRpY3NtZXRhZGF0YRI4Cghh'
    'dXRoZmxvdxis68JEIAEoDjIZLmNvZ25pdG9faWRwLkF1dGhGbG93VHlwZVIIYXV0aGZsb3cSXw'
    'oOYXV0aHBhcmFtZXRlcnMYyPmTeyADKAsyNC5jb2duaXRvX2lkcC5Jbml0aWF0ZUF1dGhSZXF1'
    'ZXN0LkF1dGhwYXJhbWV0ZXJzRW50cnlSDmF1dGhwYXJhbWV0ZXJzEh4KCGNsaWVudGlkGKToht'
    'YBIAEoCVIIY2xpZW50aWQSXwoOY2xpZW50bWV0YWRhdGEYzK//YSADKAsyNC5jb2duaXRvX2lk'
    'cC5Jbml0aWF0ZUF1dGhSZXF1ZXN0LkNsaWVudG1ldGFkYXRhRW50cnlSDmNsaWVudG1ldGFkYX'
    'RhEhsKB3Nlc3Npb24YmJmjAiABKAlSB3Nlc3Npb24STQoPdXNlcmNvbnRleHRkYXRhGJ7/hFEg'
    'ASgLMiAuY29nbml0b19pZHAuVXNlckNvbnRleHREYXRhVHlwZVIPdXNlcmNvbnRleHRkYXRhGk'
    'EKE0F1dGhwYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4ARpBChNDbGllbnRtZXRhZGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBX'
    'ZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use initiateAuthResponseDescriptor instead')
const InitiateAuthResponse$json = {
  '1': 'InitiateAuthResponse',
  '2': [
    {
      '1': 'authenticationresult',
      '3': 519327313,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AuthenticationResultType',
      '10': 'authenticationresult'
    },
    {
      '1': 'availablechallenges',
      '3': 275725719,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'availablechallenges'
    },
    {
      '1': 'challengename',
      '3': 170761310,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'challengename'
    },
    {
      '1': 'challengeparameters',
      '3': 79757023,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.InitiateAuthResponse.ChallengeparametersEntry',
      '10': 'challengeparameters'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
  ],
  '3': [InitiateAuthResponse_ChallengeparametersEntry$json],
};

@$core.Deprecated('Use initiateAuthResponseDescriptor instead')
const InitiateAuthResponse_ChallengeparametersEntry$json = {
  '1': 'ChallengeparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `InitiateAuthResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List initiateAuthResponseDescriptor = $convert.base64Decode(
    'ChRJbml0aWF0ZUF1dGhSZXNwb25zZRJdChRhdXRoZW50aWNhdGlvbnJlc3VsdBjRnNH3ASABKA'
    'syJS5jb2duaXRvX2lkcC5BdXRoZW50aWNhdGlvblJlc3VsdFR5cGVSFGF1dGhlbnRpY2F0aW9u'
    'cmVzdWx0ElQKE2F2YWlsYWJsZWNoYWxsZW5nZXMYl/u8gwEgAygOMh4uY29nbml0b19pZHAuQ2'
    'hhbGxlbmdlTmFtZVR5cGVSE2F2YWlsYWJsZWNoYWxsZW5nZXMSRwoNY2hhbGxlbmdlbmFtZRje'
    'uLZRIAEoDjIeLmNvZ25pdG9faWRwLkNoYWxsZW5nZU5hbWVUeXBlUg1jaGFsbGVuZ2VuYW1lEm'
    '8KE2NoYWxsZW5nZXBhcmFtZXRlcnMY3/2DJiADKAsyOi5jb2duaXRvX2lkcC5Jbml0aWF0ZUF1'
    'dGhSZXNwb25zZS5DaGFsbGVuZ2VwYXJhbWV0ZXJzRW50cnlSE2NoYWxsZW5nZXBhcmFtZXRlcn'
    'MSGwoHc2Vzc2lvbhiYmaMCIAEoCVIHc2Vzc2lvbhpGChhDaGFsbGVuZ2VwYXJhbWV0ZXJzRW50'
    'cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use internalErrorExceptionDescriptor instead')
const InternalErrorException$json = {
  '1': 'InternalErrorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalErrorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalErrorExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbnRlcm5hbEVycm9yRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use internalServerExceptionDescriptor instead')
const InternalServerException$json = {
  '1': 'InternalServerException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServerException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServerExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnRlcm5hbFNlcnZlckV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidEmailRoleAccessPolicyExceptionDescriptor instead')
const InvalidEmailRoleAccessPolicyException$json = {
  '1': 'InvalidEmailRoleAccessPolicyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEmailRoleAccessPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEmailRoleAccessPolicyExceptionDescriptor =
    $convert.base64Decode(
        'CiVJbnZhbGlkRW1haWxSb2xlQWNjZXNzUG9saWN5RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidLambdaResponseExceptionDescriptor instead')
const InvalidLambdaResponseException$json = {
  '1': 'InvalidLambdaResponseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidLambdaResponseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidLambdaResponseExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkTGFtYmRhUmVzcG9uc2VFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidOAuthFlowExceptionDescriptor instead')
const InvalidOAuthFlowException$json = {
  '1': 'InvalidOAuthFlowException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOAuthFlowException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOAuthFlowExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkT0F1dGhGbG93RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidParameterExceptionDescriptor instead')
const InvalidParameterException$json = {
  '1': 'InvalidParameterException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'reasoncode', '3': 207873105, '4': 1, '5': 9, '10': 'reasoncode'},
  ],
};

/// Descriptor for `InvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkUGFyYW1ldGVyRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'USIQoKcmVhc29uY29kZRjRyI9jIAEoCVIKcmVhc29uY29kZQ==');

@$core.Deprecated('Use invalidPasswordExceptionDescriptor instead')
const InvalidPasswordException$json = {
  '1': 'InvalidPasswordException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidPasswordException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidPasswordExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkUGFzc3dvcmRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidSmsRoleAccessPolicyExceptionDescriptor instead')
const InvalidSmsRoleAccessPolicyException$json = {
  '1': 'InvalidSmsRoleAccessPolicyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSmsRoleAccessPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSmsRoleAccessPolicyExceptionDescriptor =
    $convert.base64Decode(
        'CiNJbnZhbGlkU21zUm9sZUFjY2Vzc1BvbGljeUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use invalidSmsRoleTrustRelationshipExceptionDescriptor instead')
const InvalidSmsRoleTrustRelationshipException$json = {
  '1': 'InvalidSmsRoleTrustRelationshipException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSmsRoleTrustRelationshipException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSmsRoleTrustRelationshipExceptionDescriptor =
    $convert.base64Decode(
        'CihJbnZhbGlkU21zUm9sZVRydXN0UmVsYXRpb25zaGlwRXhjZXB0aW9uEhsKB21lc3NhZ2UY5Z'
        'HIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidUserPoolConfigurationExceptionDescriptor instead')
const InvalidUserPoolConfigurationException$json = {
  '1': 'InvalidUserPoolConfigurationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidUserPoolConfigurationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidUserPoolConfigurationExceptionDescriptor =
    $convert.base64Decode(
        'CiVJbnZhbGlkVXNlclBvb2xDb25maWd1cmF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use lambdaConfigTypeDescriptor instead')
const LambdaConfigType$json = {
  '1': 'LambdaConfigType',
  '2': [
    {
      '1': 'createauthchallenge',
      '3': 254734041,
      '4': 1,
      '5': 9,
      '10': 'createauthchallenge'
    },
    {
      '1': 'customemailsender',
      '3': 319793392,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CustomEmailLambdaVersionConfigType',
      '10': 'customemailsender'
    },
    {
      '1': 'custommessage',
      '3': 184293966,
      '4': 1,
      '5': 9,
      '10': 'custommessage'
    },
    {
      '1': 'customsmssender',
      '3': 394119899,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CustomSMSLambdaVersionConfigType',
      '10': 'customsmssender'
    },
    {
      '1': 'defineauthchallenge',
      '3': 220544290,
      '4': 1,
      '5': 9,
      '10': 'defineauthchallenge'
    },
    {
      '1': 'inboundfederation',
      '3': 87415260,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.InboundFederationLambdaType',
      '10': 'inboundfederation'
    },
    {'1': 'kmskeyid', '3': 13250477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'postauthentication',
      '3': 199890284,
      '4': 1,
      '5': 9,
      '10': 'postauthentication'
    },
    {
      '1': 'postconfirmation',
      '3': 178069195,
      '4': 1,
      '5': 9,
      '10': 'postconfirmation'
    },
    {
      '1': 'preauthentication',
      '3': 374491163,
      '4': 1,
      '5': 9,
      '10': 'preauthentication'
    },
    {'1': 'presignup', '3': 249108711, '4': 1, '5': 9, '10': 'presignup'},
    {
      '1': 'pretokengeneration',
      '3': 285445042,
      '4': 1,
      '5': 9,
      '10': 'pretokengeneration'
    },
    {
      '1': 'pretokengenerationconfig',
      '3': 197795476,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.PreTokenGenerationVersionConfigType',
      '10': 'pretokengenerationconfig'
    },
    {
      '1': 'usermigration',
      '3': 31196067,
      '4': 1,
      '5': 9,
      '10': 'usermigration'
    },
    {
      '1': 'verifyauthchallengeresponse',
      '3': 152402305,
      '4': 1,
      '5': 9,
      '10': 'verifyauthchallengeresponse'
    },
  ],
};

/// Descriptor for `LambdaConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaConfigTypeDescriptor = $convert.base64Decode(
    'ChBMYW1iZGFDb25maWdUeXBlEjMKE2NyZWF0ZWF1dGhjaGFsbGVuZ2UY2d27eSABKAlSE2NyZW'
    'F0ZWF1dGhjaGFsbGVuZ2USYQoRY3VzdG9tZW1haWxzZW5kZXIY8NG+mAEgASgLMi8uY29nbml0'
    'b19pZHAuQ3VzdG9tRW1haWxMYW1iZGFWZXJzaW9uQ29uZmlnVHlwZVIRY3VzdG9tZW1haWxzZW'
    '5kZXISJwoNY3VzdG9tbWVzc2FnZRjOtPBXIAEoCVINY3VzdG9tbWVzc2FnZRJbCg9jdXN0b21z'
    'bXNzZW5kZXIY25X3uwEgASgLMi0uY29nbml0b19pZHAuQ3VzdG9tU01TTGFtYmRhVmVyc2lvbk'
    'NvbmZpZ1R5cGVSD2N1c3RvbXNtc3NlbmRlchIzChNkZWZpbmVhdXRoY2hhbGxlbmdlGKL6lGkg'
    'ASgJUhNkZWZpbmVhdXRoY2hhbGxlbmdlElkKEWluYm91bmRmZWRlcmF0aW9uGNyz1ykgASgLMi'
    'guY29nbml0b19pZHAuSW5ib3VuZEZlZGVyYXRpb25MYW1iZGFUeXBlUhFpbmJvdW5kZmVkZXJh'
    'dGlvbhIdCghrbXNrZXlpZBit36gGIAEoCVIIa21za2V5aWQSMQoScG9zdGF1dGhlbnRpY2F0aW'
    '9uGOyqqF8gASgJUhJwb3N0YXV0aGVudGljYXRpb24SLQoQcG9zdGNvbmZpcm1hdGlvbhjLvfRU'
    'IAEoCVIQcG9zdGNvbmZpcm1hdGlvbhIwChFwcmVhdXRoZW50aWNhdGlvbhibkMmyASABKAlSEX'
    'ByZWF1dGhlbnRpY2F0aW9uEh8KCXByZXNpZ251cBjnseR2IAEoCVIJcHJlc2lnbnVwEjIKEnBy'
    'ZXRva2VuZ2VuZXJhdGlvbhiyl46IASABKAlSEnByZXRva2VuZ2VuZXJhdGlvbhJvChhwcmV0b2'
    'tlbmdlbmVyYXRpb25jb25maWcYlL2oXiABKAsyMC5jb2duaXRvX2lkcC5QcmVUb2tlbkdlbmVy'
    'YXRpb25WZXJzaW9uQ29uZmlnVHlwZVIYcHJldG9rZW5nZW5lcmF0aW9uY29uZmlnEicKDXVzZX'
    'JtaWdyYXRpb24Yo4fwDiABKAlSDXVzZXJtaWdyYXRpb24SQwobdmVyaWZ5YXV0aGNoYWxsZW5n'
    'ZXJlc3BvbnNlGIHz1UggASgJUht2ZXJpZnlhdXRoY2hhbGxlbmdlcmVzcG9uc2U=');

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

@$core.Deprecated('Use listDevicesRequestDescriptor instead')
const ListDevicesRequest$json = {
  '1': 'ListDevicesRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
  ],
};

/// Descriptor for `ListDevicesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDevicesRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0RGV2aWNlc1JlcXVlc3QSIwoLYWNjZXNzdG9rZW4YibyQRiABKAlSC2FjY2Vzc3Rva2'
    'VuEhgKBWxpbWl0GNWV2cQBIAEoBVIFbGltaXQSLAoPcGFnaW5hdGlvbnRva2VuGP+KwK0BIAEo'
    'CVIPcGFnaW5hdGlvbnRva2Vu');

@$core.Deprecated('Use listDevicesResponseDescriptor instead')
const ListDevicesResponse$json = {
  '1': 'ListDevicesResponse',
  '2': [
    {
      '1': 'devices',
      '3': 128774945,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.DeviceType',
      '10': 'devices'
    },
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
  ],
};

/// Descriptor for `ListDevicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDevicesResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0RGV2aWNlc1Jlc3BvbnNlEjQKB2RldmljZXMYoeazPSADKAsyFy5jb2duaXRvX2lkcC'
    '5EZXZpY2VUeXBlUgdkZXZpY2VzEiwKD3BhZ2luYXRpb250b2tlbhj/isCtASABKAlSD3BhZ2lu'
    'YXRpb250b2tlbg==');

@$core.Deprecated('Use listGroupsRequestDescriptor instead')
const ListGroupsRequest$json = {
  '1': 'ListGroupsRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListGroupsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupsRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0R3JvdXBzUmVxdWVzdBIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0Eh8KCW5leHR0b2'
    'tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9v'
    'bGlk');

@$core.Deprecated('Use listGroupsResponseDescriptor instead')
const ListGroupsResponse$json = {
  '1': 'ListGroupsResponse',
  '2': [
    {
      '1': 'groups',
      '3': 360662414,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.GroupType',
      '10': 'groups'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListGroupsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGroupsResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0R3JvdXBzUmVzcG9uc2USMgoGZ3JvdXBzGI6L/asBIAMoCzIWLmNvZ25pdG9faWRwLk'
    'dyb3VwVHlwZVIGZ3JvdXBzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listIdentityProvidersRequestDescriptor instead')
const ListIdentityProvidersRequest$json = {
  '1': 'ListIdentityProvidersRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListIdentityProvidersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIdentityProvidersRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0SWRlbnRpdHlQcm92aWRlcnNSZXF1ZXN0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUg'
        'ptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiIKCnVzZXJwb29s'
        'aWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use listIdentityProvidersResponseDescriptor instead')
const ListIdentityProvidersResponse$json = {
  '1': 'ListIdentityProvidersResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'providers',
      '3': 122628990,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ProviderDescription',
      '10': 'providers'
    },
  ],
};

/// Descriptor for `ListIdentityProvidersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIdentityProvidersResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0SWRlbnRpdHlQcm92aWRlcnNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW'
        '5leHR0b2tlbhJBCglwcm92aWRlcnMY/ta8OiADKAsyIC5jb2duaXRvX2lkcC5Qcm92aWRlckRl'
        'c2NyaXB0aW9uUglwcm92aWRlcnM=');

@$core.Deprecated('Use listResourceServersRequestDescriptor instead')
const ListResourceServersRequest$json = {
  '1': 'ListResourceServersRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListResourceServersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceServersRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0UmVzb3VyY2VTZXJ2ZXJzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
        'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIiCgp1c2VycG9vbGlk'
        'GP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use listResourceServersResponseDescriptor instead')
const ListResourceServersResponse$json = {
  '1': 'ListResourceServersResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourceservers',
      '3': 305925412,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ResourceServerType',
      '10': 'resourceservers'
    },
  ],
};

/// Descriptor for `ListResourceServersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceServersResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0UmVzb3VyY2VTZXJ2ZXJzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4STQoPcmVzb3VyY2VzZXJ2ZXJzGKSa8JEBIAMoCzIfLmNvZ25pdG9faWRwLlJlc291'
        'cmNlU2VydmVyVHlwZVIPcmVzb3VyY2VzZXJ2ZXJz');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3'
        'Jlc291cmNlYXJu');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ListTagsForResourceResponse.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [ListTagsForResourceResponse_TagsEntry$json],
};

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USSgoEdGFncxjBwfa1ASADKAsyMi5jb2duaX'
        'RvX2lkcC5MaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2UuVGFnc0VudHJ5UgR0YWdzGjcKCVRh'
        'Z3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use listTermsRequestDescriptor instead')
const ListTermsRequest$json = {
  '1': 'ListTermsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListTermsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTermsRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0VGVybXNSZXF1ZXN0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh'
    '8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiIKCnVzZXJwb29saWQY/saLnQEgASgJ'
    'Ugp1c2VycG9vbGlk');

@$core.Deprecated('Use listTermsResponseDescriptor instead')
const ListTermsResponse$json = {
  '1': 'ListTermsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'terms',
      '3': 339062221,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.TermsDescriptionType',
      '10': 'terms'
    },
  ],
};

/// Descriptor for `ListTermsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTermsResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0VGVybXNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhI7Cg'
    'V0ZXJtcxjN29ahASADKAsyIS5jb2duaXRvX2lkcC5UZXJtc0Rlc2NyaXB0aW9uVHlwZVIFdGVy'
    'bXM=');

@$core.Deprecated('Use listUserImportJobsRequestDescriptor instead')
const ListUserImportJobsRequest$json = {
  '1': 'ListUserImportJobsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListUserImportJobsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserImportJobsRequestDescriptor = $convert.base64Decode(
    'ChlMaXN0VXNlckltcG9ydEpvYnNSZXF1ZXN0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYX'
    'hyZXN1bHRzEiwKD3BhZ2luYXRpb250b2tlbhj/isCtASABKAlSD3BhZ2luYXRpb250b2tlbhIi'
    'Cgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use listUserImportJobsResponseDescriptor instead')
const ListUserImportJobsResponse$json = {
  '1': 'ListUserImportJobsResponse',
  '2': [
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
    {
      '1': 'userimportjobs',
      '3': 388966560,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UserImportJobType',
      '10': 'userimportjobs'
    },
  ],
};

/// Descriptor for `ListUserImportJobsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserImportJobsResponseDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VXNlckltcG9ydEpvYnNSZXNwb25zZRIsCg9wYWdpbmF0aW9udG9rZW4Y/4rArQEgAS'
        'gJUg9wYWdpbmF0aW9udG9rZW4SSgoOdXNlcmltcG9ydGpvYnMYoNG8uQEgAygLMh4uY29nbml0'
        'b19pZHAuVXNlckltcG9ydEpvYlR5cGVSDnVzZXJpbXBvcnRqb2Jz');

@$core.Deprecated('Use listUserPoolClientSecretsRequestDescriptor instead')
const ListUserPoolClientSecretsRequest$json = {
  '1': 'ListUserPoolClientSecretsRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListUserPoolClientSecretsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoolClientSecretsRequestDescriptor =
    $convert.base64Decode(
        'CiBMaXN0VXNlclBvb2xDbGllbnRTZWNyZXRzUmVxdWVzdBIeCghjbGllbnRpZBik6IbWASABKA'
        'lSCGNsaWVudGlkEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiIKCnVzZXJwb29s'
        'aWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use listUserPoolClientSecretsResponseDescriptor instead')
const ListUserPoolClientSecretsResponse$json = {
  '1': 'ListUserPoolClientSecretsResponse',
  '2': [
    {
      '1': 'clientsecrets',
      '3': 16465760,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ClientSecretDescriptorType',
      '10': 'clientsecrets'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListUserPoolClientSecretsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoolClientSecretsResponseDescriptor =
    $convert.base64Decode(
        'CiFMaXN0VXNlclBvb2xDbGllbnRTZWNyZXRzUmVzcG9uc2USUAoNY2xpZW50c2VjcmV0cxjg/u'
        'wHIAMoCzInLmNvZ25pdG9faWRwLkNsaWVudFNlY3JldERlc2NyaXB0b3JUeXBlUg1jbGllbnRz'
        'ZWNyZXRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listUserPoolClientsRequestDescriptor instead')
const ListUserPoolClientsRequest$json = {
  '1': 'ListUserPoolClientsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListUserPoolClientsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoolClientsRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VXNlclBvb2xDbGllbnRzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
        'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIiCgp1c2VycG9vbGlk'
        'GP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use listUserPoolClientsResponseDescriptor instead')
const ListUserPoolClientsResponse$json = {
  '1': 'ListUserPoolClientsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'userpoolclients',
      '3': 449808661,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UserPoolClientDescription',
      '10': 'userpoolclients'
    },
  ],
};

/// Descriptor for `ListUserPoolClientsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoolClientsResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VXNlclBvb2xDbGllbnRzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SVAoPdXNlcnBvb2xjbGllbnRzGJWSvtYBIAMoCzImLmNvZ25pdG9faWRwLlVzZXJQ'
        'b29sQ2xpZW50RGVzY3JpcHRpb25SD3VzZXJwb29sY2xpZW50cw==');

@$core.Deprecated('Use listUserPoolsRequestDescriptor instead')
const ListUserPoolsRequest$json = {
  '1': 'ListUserPoolsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListUserPoolsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoolsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0VXNlclBvb2xzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdW'
    'x0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listUserPoolsResponseDescriptor instead')
const ListUserPoolsResponse$json = {
  '1': 'ListUserPoolsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'userpools',
      '3': 322489898,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UserPoolDescriptionType',
      '10': 'userpools'
    },
  ],
};

/// Descriptor for `ListUserPoolsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUserPoolsResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0VXNlclBvb2xzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
    '4SRgoJdXNlcnBvb2xzGKqc45kBIAMoCzIkLmNvZ25pdG9faWRwLlVzZXJQb29sRGVzY3JpcHRp'
    'b25UeXBlUgl1c2VycG9vbHM=');

@$core.Deprecated('Use listUsersInGroupRequestDescriptor instead')
const ListUsersInGroupRequest$json = {
  '1': 'ListUsersInGroupRequest',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListUsersInGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUsersInGroupRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0VXNlcnNJbkdyb3VwUmVxdWVzdBIgCglncm91cG5hbWUYyMqgqgEgASgJUglncm91cG'
    '5hbWUSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5l'
    'eHR0b2tlbhIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use listUsersInGroupResponseDescriptor instead')
const ListUsersInGroupResponse$json = {
  '1': 'ListUsersInGroupResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'users',
      '3': 112472756,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UserType',
      '10': 'users'
    },
  ],
};

/// Descriptor for `ListUsersInGroupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUsersInGroupResponseDescriptor =
    $convert.base64Decode(
        'ChhMaXN0VXNlcnNJbkdyb3VwUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
        '9rZW4SLgoFdXNlcnMYtOXQNSADKAsyFS5jb2duaXRvX2lkcC5Vc2VyVHlwZVIFdXNlcnM=');

@$core.Deprecated('Use listUsersRequestDescriptor instead')
const ListUsersRequest$json = {
  '1': 'ListUsersRequest',
  '2': [
    {
      '1': 'attributestoget',
      '3': 311382592,
      '4': 3,
      '5': 9,
      '10': 'attributestoget'
    },
    {'1': 'filter', '3': 346669208, '4': 1, '5': 9, '10': 'filter'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ListUsersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUsersRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0VXNlcnNSZXF1ZXN0EiwKD2F0dHJpYnV0ZXN0b2dldBjApL2UASADKAlSD2F0dHJpYn'
    'V0ZXN0b2dldBIaCgZmaWx0ZXIYmIGnpQEgASgJUgZmaWx0ZXISGAoFbGltaXQY1ZXZxAEgASgF'
    'UgVsaW1pdBIsCg9wYWdpbmF0aW9udG9rZW4Y/4rArQEgASgJUg9wYWdpbmF0aW9udG9rZW4SIg'
    'oKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use listUsersResponseDescriptor instead')
const ListUsersResponse$json = {
  '1': 'ListUsersResponse',
  '2': [
    {
      '1': 'paginationtoken',
      '3': 363857279,
      '4': 1,
      '5': 9,
      '10': 'paginationtoken'
    },
    {
      '1': 'users',
      '3': 112472756,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UserType',
      '10': 'users'
    },
  ],
};

/// Descriptor for `ListUsersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listUsersResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0VXNlcnNSZXNwb25zZRIsCg9wYWdpbmF0aW9udG9rZW4Y/4rArQEgASgJUg9wYWdpbm'
    'F0aW9udG9rZW4SLgoFdXNlcnMYtOXQNSADKAsyFS5jb2duaXRvX2lkcC5Vc2VyVHlwZVIFdXNl'
    'cnM=');

@$core.Deprecated('Use listWebAuthnCredentialsRequestDescriptor instead')
const ListWebAuthnCredentialsRequest$json = {
  '1': 'ListWebAuthnCredentialsRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListWebAuthnCredentialsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWebAuthnCredentialsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0V2ViQXV0aG5DcmVkZW50aWFsc1JlcXVlc3QSIwoLYWNjZXNzdG9rZW4YibyQRiABKA'
        'lSC2FjY2Vzc3Rva2VuEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh8KCW5l'
        'eHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listWebAuthnCredentialsResponseDescriptor instead')
const ListWebAuthnCredentialsResponse$json = {
  '1': 'ListWebAuthnCredentialsResponse',
  '2': [
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.WebAuthnCredentialDescription',
      '10': 'credentials'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListWebAuthnCredentialsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWebAuthnCredentialsResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0V2ViQXV0aG5DcmVkZW50aWFsc1Jlc3BvbnNlElAKC2NyZWRlbnRpYWxzGPKajrYBIA'
        'MoCzIqLmNvZ25pdG9faWRwLldlYkF1dGhuQ3JlZGVudGlhbERlc2NyaXB0aW9uUgtjcmVkZW50'
        'aWFscxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use logConfigurationTypeDescriptor instead')
const LogConfigurationType$json = {
  '1': 'LogConfigurationType',
  '2': [
    {
      '1': 'cloudwatchlogsconfiguration',
      '3': 300899477,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CloudWatchLogsConfigurationType',
      '10': 'cloudwatchlogsconfiguration'
    },
    {
      '1': 'eventsource',
      '3': 37841339,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.EventSourceName',
      '10': 'eventsource'
    },
    {
      '1': 'firehoseconfiguration',
      '3': 141015491,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.FirehoseConfigurationType',
      '10': 'firehoseconfiguration'
    },
    {
      '1': 'loglevel',
      '3': 171121018,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.LogLevel',
      '10': 'loglevel'
    },
    {
      '1': 's3configuration',
      '3': 27828476,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.S3ConfigurationType',
      '10': 's3configuration'
    },
  ],
};

/// Descriptor for `LogConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logConfigurationTypeDescriptor = $convert.base64Decode(
    'ChRMb2dDb25maWd1cmF0aW9uVHlwZRJyChtjbG91ZHdhdGNobG9nc2NvbmZpZ3VyYXRpb24Ylb'
    'm9jwEgASgLMiwuY29nbml0b19pZHAuQ2xvdWRXYXRjaExvZ3NDb25maWd1cmF0aW9uVHlwZVIb'
    'Y2xvdWR3YXRjaGxvZ3Njb25maWd1cmF0aW9uEkEKC2V2ZW50c291cmNlGLvThRIgASgOMhwuY2'
    '9nbml0b19pZHAuRXZlbnRTb3VyY2VOYW1lUgtldmVudHNvdXJjZRJfChVmaXJlaG9zZWNvbmZp'
    'Z3VyYXRpb24Yw/OeQyABKAsyJi5jb2duaXRvX2lkcC5GaXJlaG9zZUNvbmZpZ3VyYXRpb25UeX'
    'BlUhVmaXJlaG9zZWNvbmZpZ3VyYXRpb24SNAoIbG9nbGV2ZWwY+rLMUSABKA4yFS5jb2duaXRv'
    'X2lkcC5Mb2dMZXZlbFIIbG9nbGV2ZWwSTQoPczNjb25maWd1cmF0aW9uGPzBog0gASgLMiAuY2'
    '9nbml0b19pZHAuUzNDb25maWd1cmF0aW9uVHlwZVIPczNjb25maWd1cmF0aW9u');

@$core.Deprecated('Use logDeliveryConfigurationTypeDescriptor instead')
const LogDeliveryConfigurationType$json = {
  '1': 'LogDeliveryConfigurationType',
  '2': [
    {
      '1': 'logconfigurations',
      '3': 89180603,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.LogConfigurationType',
      '10': 'logconfigurations'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `LogDeliveryConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logDeliveryConfigurationTypeDescriptor =
    $convert.base64Decode(
        'ChxMb2dEZWxpdmVyeUNvbmZpZ3VyYXRpb25UeXBlElIKEWxvZ2NvbmZpZ3VyYXRpb25zGLuTwy'
        'ogAygLMiEuY29nbml0b19pZHAuTG9nQ29uZmlndXJhdGlvblR5cGVSEWxvZ2NvbmZpZ3VyYXRp'
        'b25zEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use mFAMethodNotFoundExceptionDescriptor instead')
const MFAMethodNotFoundException$json = {
  '1': 'MFAMethodNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MFAMethodNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mFAMethodNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChpNRkFNZXRob2ROb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use mFAOptionTypeDescriptor instead')
const MFAOptionType$json = {
  '1': 'MFAOptionType',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'deliverymedium',
      '3': 140498471,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DeliveryMediumType',
      '10': 'deliverymedium'
    },
  ],
};

/// Descriptor for `MFAOptionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mFAOptionTypeDescriptor = $convert.base64Decode(
    'Cg1NRkFPcHRpb25UeXBlEigKDWF0dHJpYnV0ZW5hbWUYrZWYqAEgASgJUg1hdHRyaWJ1dGVuYW'
    '1lEkoKDmRlbGl2ZXJ5bWVkaXVtGKes/0IgASgOMh8uY29nbml0b19pZHAuRGVsaXZlcnlNZWRp'
    'dW1UeXBlUg5kZWxpdmVyeW1lZGl1bQ==');

@$core.Deprecated('Use managedLoginBrandingExistsExceptionDescriptor instead')
const ManagedLoginBrandingExistsException$json = {
  '1': 'ManagedLoginBrandingExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ManagedLoginBrandingExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedLoginBrandingExistsExceptionDescriptor =
    $convert.base64Decode(
        'CiNNYW5hZ2VkTG9naW5CcmFuZGluZ0V4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use managedLoginBrandingTypeDescriptor instead')
const ManagedLoginBrandingType$json = {
  '1': 'ManagedLoginBrandingType',
  '2': [
    {
      '1': 'assets',
      '3': 141150041,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AssetType',
      '10': 'assets'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'managedloginbrandingid',
      '3': 57829482,
      '4': 1,
      '5': 9,
      '10': 'managedloginbrandingid'
    },
    {'1': 'settings', '3': 184911657, '4': 1, '5': 9, '10': 'settings'},
    {
      '1': 'usecognitoprovidedvalues',
      '3': 408662835,
      '4': 1,
      '5': 8,
      '10': 'usecognitoprovidedvalues'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ManagedLoginBrandingType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedLoginBrandingTypeDescriptor = $convert.base64Decode(
    'ChhNYW5hZ2VkTG9naW5CcmFuZGluZ1R5cGUSMQoGYXNzZXRzGNmOp0MgAygLMhYuY29nbml0b1'
    '9pZHAuQXNzZXRUeXBlUgZhc3NldHMSJgoMY3JlYXRpb25kYXRlGOHYt4kBIAEoCVIMY3JlYXRp'
    'b25kYXRlEi0KEGxhc3Rtb2RpZmllZGRhdGUY04jICyABKAlSEGxhc3Rtb2RpZmllZGRhdGUSOQ'
    'oWbWFuYWdlZGxvZ2luYnJhbmRpbmdpZBjq0MkbIAEoCVIWbWFuYWdlZGxvZ2luYnJhbmRpbmdp'
    'ZBIdCghzZXR0aW5ncxipjpZYIAEoCVIIc2V0dGluZ3MSPgoYdXNlY29nbml0b3Byb3ZpZGVkdm'
    'FsdWVzGLPm7sIBIAEoCFIYdXNlY29nbml0b3Byb3ZpZGVkdmFsdWVzEiIKCnVzZXJwb29saWQY'
    '/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use messageTemplateTypeDescriptor instead')
const MessageTemplateType$json = {
  '1': 'MessageTemplateType',
  '2': [
    {'1': 'emailmessage', '3': 247992871, '4': 1, '5': 9, '10': 'emailmessage'},
    {'1': 'emailsubject', '3': 215040870, '4': 1, '5': 9, '10': 'emailsubject'},
    {'1': 'smsmessage', '3': 441530194, '4': 1, '5': 9, '10': 'smsmessage'},
  ],
};

/// Descriptor for `MessageTemplateType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageTemplateTypeDescriptor = $convert.base64Decode(
    'ChNNZXNzYWdlVGVtcGxhdGVUeXBlEiUKDGVtYWlsbWVzc2FnZRinpKB2IAEoCVIMZW1haWxtZX'
    'NzYWdlEiUKDGVtYWlsc3ViamVjdBjmhsVmIAEoCVIMZW1haWxzdWJqZWN0EiIKCnNtc21lc3Nh'
    'Z2UY0u7E0gEgASgJUgpzbXNtZXNzYWdl');

@$core.Deprecated('Use newDeviceMetadataTypeDescriptor instead')
const NewDeviceMetadataType$json = {
  '1': 'NewDeviceMetadataType',
  '2': [
    {
      '1': 'devicegroupkey',
      '3': 11731234,
      '4': 1,
      '5': 9,
      '10': 'devicegroupkey'
    },
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
  ],
};

/// Descriptor for `NewDeviceMetadataType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List newDeviceMetadataTypeDescriptor = $convert.base64Decode(
    'ChVOZXdEZXZpY2VNZXRhZGF0YVR5cGUSKQoOZGV2aWNlZ3JvdXBrZXkYooLMBSABKAlSDmRldm'
    'ljZWdyb3Vwa2V5EiAKCWRldmljZWtleRir5Mi2ASABKAlSCWRldmljZWtleQ==');

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

@$core.Deprecated('Use notifyConfigurationTypeDescriptor instead')
const NotifyConfigurationType$json = {
  '1': 'NotifyConfigurationType',
  '2': [
    {
      '1': 'blockemail',
      '3': 413505327,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.NotifyEmailType',
      '10': 'blockemail'
    },
    {'1': 'from', '3': 410269078, '4': 1, '5': 9, '10': 'from'},
    {
      '1': 'mfaemail',
      '3': 330686644,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.NotifyEmailType',
      '10': 'mfaemail'
    },
    {
      '1': 'noactionemail',
      '3': 203016289,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.NotifyEmailType',
      '10': 'noactionemail'
    },
    {'1': 'replyto', '3': 82405747, '4': 1, '5': 9, '10': 'replyto'},
    {'1': 'sourcearn', '3': 439903072, '4': 1, '5': 9, '10': 'sourcearn'},
  ],
};

/// Descriptor for `NotifyConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notifyConfigurationTypeDescriptor = $convert.base64Decode(
    'ChdOb3RpZnlDb25maWd1cmF0aW9uVHlwZRJACgpibG9ja2VtYWlsGK+ulsUBIAEoCzIcLmNvZ2'
    '5pdG9faWRwLk5vdGlmeUVtYWlsVHlwZVIKYmxvY2tlbWFpbBIWCgRmcm9tGJbr0MMBIAEoCVIE'
    'ZnJvbRI8CghtZmFlbWFpbBi0wdedASABKAsyHC5jb2duaXRvX2lkcC5Ob3RpZnlFbWFpbFR5cG'
    'VSCG1mYWVtYWlsEkUKDW5vYWN0aW9uZW1haWwY4ZDnYCABKAsyHC5jb2duaXRvX2lkcC5Ob3Rp'
    'ZnlFbWFpbFR5cGVSDW5vYWN0aW9uZW1haWwSGwoHcmVwbHl0bxjz0qUnIAEoCVIHcmVwbHl0bx'
    'IgCglzb3VyY2Vhcm4Y4Mbh0QEgASgJUglzb3VyY2Vhcm4=');

@$core.Deprecated('Use notifyEmailTypeDescriptor instead')
const NotifyEmailType$json = {
  '1': 'NotifyEmailType',
  '2': [
    {'1': 'htmlbody', '3': 148096253, '4': 1, '5': 9, '10': 'htmlbody'},
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
    {'1': 'textbody', '3': 91008147, '4': 1, '5': 9, '10': 'textbody'},
  ],
};

/// Descriptor for `NotifyEmailType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notifyEmailTypeDescriptor = $convert.base64Decode(
    'Cg9Ob3RpZnlFbWFpbFR5cGUSHQoIaHRtbGJvZHkY/YnPRiABKAlSCGh0bWxib2R5EhsKB3N1Ym'
    'plY3QY8MnkAyABKAlSB3N1YmplY3QSHQoIdGV4dGJvZHkYk9myKyABKAlSCHRleHRib2R5');

@$core.Deprecated('Use numberAttributeConstraintsTypeDescriptor instead')
const NumberAttributeConstraintsType$json = {
  '1': 'NumberAttributeConstraintsType',
  '2': [
    {'1': 'maxvalue', '3': 480618407, '4': 1, '5': 9, '10': 'maxvalue'},
    {'1': 'minvalue', '3': 454592205, '4': 1, '5': 9, '10': 'minvalue'},
  ],
};

/// Descriptor for `NumberAttributeConstraintsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List numberAttributeConstraintsTypeDescriptor =
    $convert.base64Decode(
        'Ch5OdW1iZXJBdHRyaWJ1dGVDb25zdHJhaW50c1R5cGUSHgoIbWF4dmFsdWUYp8+W5QEgASgJUg'
        'htYXh2YWx1ZRIeCghtaW52YWx1ZRjNjeLYASABKAlSCG1pbnZhbHVl');

@$core
    .Deprecated('Use passwordHistoryPolicyViolationExceptionDescriptor instead')
const PasswordHistoryPolicyViolationException$json = {
  '1': 'PasswordHistoryPolicyViolationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PasswordHistoryPolicyViolationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List passwordHistoryPolicyViolationExceptionDescriptor =
    $convert.base64Decode(
        'CidQYXNzd29yZEhpc3RvcnlQb2xpY3lWaW9sYXRpb25FeGNlcHRpb24SGwoHbWVzc2FnZRjlkc'
        'gnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use passwordPolicyTypeDescriptor instead')
const PasswordPolicyType$json = {
  '1': 'PasswordPolicyType',
  '2': [
    {
      '1': 'minimumlength',
      '3': 417878840,
      '4': 1,
      '5': 5,
      '10': 'minimumlength'
    },
    {
      '1': 'passwordhistorysize',
      '3': 413355954,
      '4': 1,
      '5': 5,
      '10': 'passwordhistorysize'
    },
    {
      '1': 'requirelowercase',
      '3': 14334134,
      '4': 1,
      '5': 8,
      '10': 'requirelowercase'
    },
    {
      '1': 'requirenumbers',
      '3': 385325811,
      '4': 1,
      '5': 8,
      '10': 'requirenumbers'
    },
    {
      '1': 'requiresymbols',
      '3': 98968906,
      '4': 1,
      '5': 8,
      '10': 'requiresymbols'
    },
    {
      '1': 'requireuppercase',
      '3': 209568997,
      '4': 1,
      '5': 8,
      '10': 'requireuppercase'
    },
    {
      '1': 'temporarypasswordvaliditydays',
      '3': 379562265,
      '4': 1,
      '5': 5,
      '10': 'temporarypasswordvaliditydays'
    },
  ],
};

/// Descriptor for `PasswordPolicyType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List passwordPolicyTypeDescriptor = $convert.base64Decode(
    'ChJQYXNzd29yZFBvbGljeVR5cGUSKAoNbWluaW11bWxlbmd0aBi4pqHHASABKAVSDW1pbmltdW'
    '1sZW5ndGgSNAoTcGFzc3dvcmRoaXN0b3J5c2l6ZRiyn43FASABKAVSE3Bhc3N3b3JkaGlzdG9y'
    'eXNpemUSLQoQcmVxdWlyZWxvd2VyY2FzZRi28eoGIAEoCFIQcmVxdWlyZWxvd2VyY2FzZRIqCg'
    '5yZXF1aXJlbnVtYmVycxjztd63ASABKAhSDnJlcXVpcmVudW1iZXJzEikKDnJlcXVpcmVzeW1i'
    'b2xzGMrKmC8gASgIUg5yZXF1aXJlc3ltYm9scxItChByZXF1aXJldXBwZXJjYXNlGOWJ92MgAS'
    'gIUhByZXF1aXJldXBwZXJjYXNlEkgKHXRlbXBvcmFyeXBhc3N3b3JkdmFsaWRpdHlkYXlzGJnS'
    '/rQBIAEoBVIddGVtcG9yYXJ5cGFzc3dvcmR2YWxpZGl0eWRheXM=');

@$core.Deprecated('Use passwordResetRequiredExceptionDescriptor instead')
const PasswordResetRequiredException$json = {
  '1': 'PasswordResetRequiredException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PasswordResetRequiredException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List passwordResetRequiredExceptionDescriptor =
    $convert.base64Decode(
        'Ch5QYXNzd29yZFJlc2V0UmVxdWlyZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use preTokenGenerationVersionConfigTypeDescriptor instead')
const PreTokenGenerationVersionConfigType$json = {
  '1': 'PreTokenGenerationVersionConfigType',
  '2': [
    {'1': 'lambdaarn', '3': 196258582, '4': 1, '5': 9, '10': 'lambdaarn'},
    {
      '1': 'lambdaversion',
      '3': 429236631,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.PreTokenGenerationLambdaVersionType',
      '10': 'lambdaversion'
    },
  ],
};

/// Descriptor for `PreTokenGenerationVersionConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List preTokenGenerationVersionConfigTypeDescriptor =
    $convert.base64Decode(
        'CiNQcmVUb2tlbkdlbmVyYXRpb25WZXJzaW9uQ29uZmlnVHlwZRIfCglsYW1iZGFhcm4YltbKXS'
        'ABKAlSCWxhbWJkYWFybhJaCg1sYW1iZGF2ZXJzaW9uGJfD1swBIAEoDjIwLmNvZ25pdG9faWRw'
        'LlByZVRva2VuR2VuZXJhdGlvbkxhbWJkYVZlcnNpb25UeXBlUg1sYW1iZGF2ZXJzaW9u');

@$core.Deprecated('Use preconditionNotMetExceptionDescriptor instead')
const PreconditionNotMetException$json = {
  '1': 'PreconditionNotMetException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PreconditionNotMetException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List preconditionNotMetExceptionDescriptor =
    $convert.base64Decode(
        'ChtQcmVjb25kaXRpb25Ob3RNZXRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use providerDescriptionDescriptor instead')
const ProviderDescription$json = {
  '1': 'ProviderDescription',
  '2': [
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
    {
      '1': 'providertype',
      '3': 337296789,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.IdentityProviderTypeType',
      '10': 'providertype'
    },
  ],
};

/// Descriptor for `ProviderDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List providerDescriptionDescriptor = $convert.base64Decode(
    'ChNQcm92aWRlckRlc2NyaXB0aW9uEiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW'
    '9uZGF0ZRItChBsYXN0bW9kaWZpZWRkYXRlGNOIyAsgASgJUhBsYXN0bW9kaWZpZWRkYXRlEiYK'
    'DHByb3ZpZGVybmFtZRj4oajnASABKAlSDHByb3ZpZGVybmFtZRJNCgxwcm92aWRlcnR5cGUYlf'
    'vqoAEgASgOMiUuY29nbml0b19pZHAuSWRlbnRpdHlQcm92aWRlclR5cGVUeXBlUgxwcm92aWRl'
    'cnR5cGU=');

@$core.Deprecated('Use providerUserIdentifierTypeDescriptor instead')
const ProviderUserIdentifierType$json = {
  '1': 'ProviderUserIdentifierType',
  '2': [
    {
      '1': 'providerattributename',
      '3': 17050040,
      '4': 1,
      '5': 9,
      '10': 'providerattributename'
    },
    {
      '1': 'providerattributevalue',
      '3': 37044150,
      '4': 1,
      '5': 9,
      '10': 'providerattributevalue'
    },
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
  ],
};

/// Descriptor for `ProviderUserIdentifierType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List providerUserIdentifierTypeDescriptor = $convert.base64Decode(
    'ChpQcm92aWRlclVzZXJJZGVudGlmaWVyVHlwZRI3ChVwcm92aWRlcmF0dHJpYnV0ZW5hbWUYuN'
    'OQCCABKAlSFXByb3ZpZGVyYXR0cmlidXRlbmFtZRI5ChZwcm92aWRlcmF0dHJpYnV0ZXZhbHVl'
    'GLb/1BEgASgJUhZwcm92aWRlcmF0dHJpYnV0ZXZhbHVlEiYKDHByb3ZpZGVybmFtZRj4oajnAS'
    'ABKAlSDHByb3ZpZGVybmFtZQ==');

@$core.Deprecated('Use recoveryOptionTypeDescriptor instead')
const RecoveryOptionType$json = {
  '1': 'RecoveryOptionType',
  '2': [
    {
      '1': 'name',
      '3': 266367751,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.RecoveryOptionNameType',
      '10': 'name'
    },
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
  ],
};

/// Descriptor for `RecoveryOptionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recoveryOptionTypeDescriptor = $convert.base64Decode(
    'ChJSZWNvdmVyeU9wdGlvblR5cGUSOgoEbmFtZRiH5oF/IAEoDjIjLmNvZ25pdG9faWRwLlJlY2'
    '92ZXJ5T3B0aW9uTmFtZVR5cGVSBG5hbWUSHQoIcHJpb3JpdHkYqr62NCABKAVSCHByaW9yaXR5');

@$core.Deprecated('Use refreshTokenReuseExceptionDescriptor instead')
const RefreshTokenReuseException$json = {
  '1': 'RefreshTokenReuseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RefreshTokenReuseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refreshTokenReuseExceptionDescriptor =
    $convert.base64Decode(
        'ChpSZWZyZXNoVG9rZW5SZXVzZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use refreshTokenRotationTypeDescriptor instead')
const RefreshTokenRotationType$json = {
  '1': 'RefreshTokenRotationType',
  '2': [
    {
      '1': 'feature',
      '3': 512819934,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.FeatureType',
      '10': 'feature'
    },
    {
      '1': 'retrygraceperiodseconds',
      '3': 498715956,
      '4': 1,
      '5': 5,
      '10': 'retrygraceperiodseconds'
    },
  ],
};

/// Descriptor for `RefreshTokenRotationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refreshTokenRotationTypeDescriptor = $convert.base64Decode(
    'ChhSZWZyZXNoVG9rZW5Sb3RhdGlvblR5cGUSNgoHZmVhdHVyZRjehcT0ASABKA4yGC5jb2duaX'
    'RvX2lkcC5GZWF0dXJlVHlwZVIHZmVhdHVyZRI8ChdyZXRyeWdyYWNlcGVyaW9kc2Vjb25kcxi0'
    'muftASABKAVSF3JldHJ5Z3JhY2VwZXJpb2RzZWNvbmRz');

@$core.Deprecated('Use resendConfirmationCodeRequestDescriptor instead')
const ResendConfirmationCodeRequest$json = {
  '1': 'ResendConfirmationCodeRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ResendConfirmationCodeRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'secrethash', '3': 321025630, '4': 1, '5': 9, '10': 'secrethash'},
    {
      '1': 'usercontextdata',
      '3': 169951134,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserContextDataType',
      '10': 'usercontextdata'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
  '3': [ResendConfirmationCodeRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use resendConfirmationCodeRequestDescriptor instead')
const ResendConfirmationCodeRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ResendConfirmationCodeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resendConfirmationCodeRequestDescriptor = $convert.base64Decode(
    'Ch1SZXNlbmRDb25maXJtYXRpb25Db2RlUmVxdWVzdBJTChFhbmFseXRpY3NtZXRhZGF0YRj3hc'
    'E8IAEoCzIiLmNvZ25pdG9faWRwLkFuYWx5dGljc01ldGFkYXRhVHlwZVIRYW5hbHl0aWNzbWV0'
    'YWRhdGESHgoIY2xpZW50aWQYpOiG1gEgASgJUghjbGllbnRpZBJpCg5jbGllbnRtZXRhZGF0YR'
    'jMr/9hIAMoCzI+LmNvZ25pdG9faWRwLlJlc2VuZENvbmZpcm1hdGlvbkNvZGVSZXF1ZXN0LkNs'
    'aWVudG1ldGFkYXRhRW50cnlSDmNsaWVudG1ldGFkYXRhEiIKCnNlY3JldGhhc2gY3uyJmQEgAS'
    'gJUgpzZWNyZXRoYXNoEk0KD3VzZXJjb250ZXh0ZGF0YRie/4RRIAEoCzIgLmNvZ25pdG9faWRw'
    'LlVzZXJDb250ZXh0RGF0YVR5cGVSD3VzZXJjb250ZXh0ZGF0YRIeCgh1c2VybmFtZRjaqaPgAS'
    'ABKAlSCHVzZXJuYW1lGkEKE0NsaWVudG1ldGFkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZXkS'
    'FAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use resendConfirmationCodeResponseDescriptor instead')
const ResendConfirmationCodeResponse$json = {
  '1': 'ResendConfirmationCodeResponse',
  '2': [
    {
      '1': 'codedeliverydetails',
      '3': 423272803,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CodeDeliveryDetailsType',
      '10': 'codedeliverydetails'
    },
  ],
};

/// Descriptor for `ResendConfirmationCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resendConfirmationCodeResponseDescriptor =
    $convert.base64Decode(
        'Ch5SZXNlbmRDb25maXJtYXRpb25Db2RlUmVzcG9uc2USWgoTY29kZWRlbGl2ZXJ5ZGV0YWlscx'
        'jjwurJASABKAsyJC5jb2duaXRvX2lkcC5Db2RlRGVsaXZlcnlEZXRhaWxzVHlwZVITY29kZWRl'
        'bGl2ZXJ5ZGV0YWlscw==');

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

@$core.Deprecated('Use resourceServerScopeTypeDescriptor instead')
const ResourceServerScopeType$json = {
  '1': 'ResourceServerScopeType',
  '2': [
    {
      '1': 'scopedescription',
      '3': 131415120,
      '4': 1,
      '5': 9,
      '10': 'scopedescription'
    },
    {'1': 'scopename', '3': 58366161, '4': 1, '5': 9, '10': 'scopename'},
  ],
};

/// Descriptor for `ResourceServerScopeType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceServerScopeTypeDescriptor =
    $convert.base64Decode(
        'ChdSZXNvdXJjZVNlcnZlclNjb3BlVHlwZRItChBzY29wZWRlc2NyaXB0aW9uGND41D4gASgJUh'
        'BzY29wZWRlc2NyaXB0aW9uEh8KCXNjb3BlbmFtZRjRseobIAEoCVIJc2NvcGVuYW1l');

@$core.Deprecated('Use resourceServerTypeDescriptor instead')
const ResourceServerType$json = {
  '1': 'ResourceServerType',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'scopes',
      '3': 464684393,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ResourceServerScopeType',
      '10': 'scopes'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `ResourceServerType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceServerTypeDescriptor = $convert.base64Decode(
    'ChJSZXNvdXJjZVNlcnZlclR5cGUSIQoKaWRlbnRpZmllchjfoPsTIAEoCVIKaWRlbnRpZmllch'
    'IVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkAKBnNjb3BlcxjpisrdASADKAsyJC5jb2duaXRvX2lk'
    'cC5SZXNvdXJjZVNlcnZlclNjb3BlVHlwZVIGc2NvcGVzEiIKCnVzZXJwb29saWQY/saLnQEgAS'
    'gJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use respondToAuthChallengeRequestDescriptor instead')
const RespondToAuthChallengeRequest$json = {
  '1': 'RespondToAuthChallengeRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {
      '1': 'challengename',
      '3': 170761310,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'challengename'
    },
    {
      '1': 'challengeresponses',
      '3': 60230615,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.RespondToAuthChallengeRequest.ChallengeresponsesEntry',
      '10': 'challengeresponses'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.RespondToAuthChallengeRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {
      '1': 'usercontextdata',
      '3': 169951134,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserContextDataType',
      '10': 'usercontextdata'
    },
  ],
  '3': [
    RespondToAuthChallengeRequest_ChallengeresponsesEntry$json,
    RespondToAuthChallengeRequest_ClientmetadataEntry$json
  ],
};

@$core.Deprecated('Use respondToAuthChallengeRequestDescriptor instead')
const RespondToAuthChallengeRequest_ChallengeresponsesEntry$json = {
  '1': 'ChallengeresponsesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use respondToAuthChallengeRequestDescriptor instead')
const RespondToAuthChallengeRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `RespondToAuthChallengeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List respondToAuthChallengeRequestDescriptor = $convert.base64Decode(
    'Ch1SZXNwb25kVG9BdXRoQ2hhbGxlbmdlUmVxdWVzdBJTChFhbmFseXRpY3NtZXRhZGF0YRj3hc'
    'E8IAEoCzIiLmNvZ25pdG9faWRwLkFuYWx5dGljc01ldGFkYXRhVHlwZVIRYW5hbHl0aWNzbWV0'
    'YWRhdGESRwoNY2hhbGxlbmdlbmFtZRjeuLZRIAEoDjIeLmNvZ25pdG9faWRwLkNoYWxsZW5nZU'
    '5hbWVUeXBlUg1jaGFsbGVuZ2VuYW1lEnUKEmNoYWxsZW5nZXJlc3BvbnNlcxjXl9wcIAMoCzJC'
    'LmNvZ25pdG9faWRwLlJlc3BvbmRUb0F1dGhDaGFsbGVuZ2VSZXF1ZXN0LkNoYWxsZW5nZXJlc3'
    'BvbnNlc0VudHJ5UhJjaGFsbGVuZ2VyZXNwb25zZXMSHgoIY2xpZW50aWQYpOiG1gEgASgJUghj'
    'bGllbnRpZBJpCg5jbGllbnRtZXRhZGF0YRjMr/9hIAMoCzI+LmNvZ25pdG9faWRwLlJlc3Bvbm'
    'RUb0F1dGhDaGFsbGVuZ2VSZXF1ZXN0LkNsaWVudG1ldGFkYXRhRW50cnlSDmNsaWVudG1ldGFk'
    'YXRhEhsKB3Nlc3Npb24YmJmjAiABKAlSB3Nlc3Npb24STQoPdXNlcmNvbnRleHRkYXRhGJ7/hF'
    'EgASgLMiAuY29nbml0b19pZHAuVXNlckNvbnRleHREYXRhVHlwZVIPdXNlcmNvbnRleHRkYXRh'
    'GkUKF0NoYWxsZW5nZXJlc3BvbnNlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGA'
    'IgASgJUgV2YWx1ZToCOAEaQQoTQ2xpZW50bWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tl'
    'eRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use respondToAuthChallengeResponseDescriptor instead')
const RespondToAuthChallengeResponse$json = {
  '1': 'RespondToAuthChallengeResponse',
  '2': [
    {
      '1': 'authenticationresult',
      '3': 519327313,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AuthenticationResultType',
      '10': 'authenticationresult'
    },
    {
      '1': 'challengename',
      '3': 170761310,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.ChallengeNameType',
      '10': 'challengename'
    },
    {
      '1': 'challengeparameters',
      '3': 79757023,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_idp.RespondToAuthChallengeResponse.ChallengeparametersEntry',
      '10': 'challengeparameters'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
  ],
  '3': [RespondToAuthChallengeResponse_ChallengeparametersEntry$json],
};

@$core.Deprecated('Use respondToAuthChallengeResponseDescriptor instead')
const RespondToAuthChallengeResponse_ChallengeparametersEntry$json = {
  '1': 'ChallengeparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `RespondToAuthChallengeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List respondToAuthChallengeResponseDescriptor = $convert.base64Decode(
    'Ch5SZXNwb25kVG9BdXRoQ2hhbGxlbmdlUmVzcG9uc2USXQoUYXV0aGVudGljYXRpb25yZXN1bH'
    'QY0ZzR9wEgASgLMiUuY29nbml0b19pZHAuQXV0aGVudGljYXRpb25SZXN1bHRUeXBlUhRhdXRo'
    'ZW50aWNhdGlvbnJlc3VsdBJHCg1jaGFsbGVuZ2VuYW1lGN64tlEgASgOMh4uY29nbml0b19pZH'
    'AuQ2hhbGxlbmdlTmFtZVR5cGVSDWNoYWxsZW5nZW5hbWUSeQoTY2hhbGxlbmdlcGFyYW1ldGVy'
    'cxjf/YMmIAMoCzJELmNvZ25pdG9faWRwLlJlc3BvbmRUb0F1dGhDaGFsbGVuZ2VSZXNwb25zZS'
    '5DaGFsbGVuZ2VwYXJhbWV0ZXJzRW50cnlSE2NoYWxsZW5nZXBhcmFtZXRlcnMSGwoHc2Vzc2lv'
    'bhiYmaMCIAEoCVIHc2Vzc2lvbhpGChhDaGFsbGVuZ2VwYXJhbWV0ZXJzRW50cnkSEAoDa2V5GA'
    'EgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use revokeTokenRequestDescriptor instead')
const RevokeTokenRequest$json = {
  '1': 'RevokeTokenRequest',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'clientsecret', '3': 500734711, '4': 1, '5': 9, '10': 'clientsecret'},
    {'1': 'token', '3': 439704531, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `RevokeTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeTokenRequestDescriptor = $convert.base64Decode(
    'ChJSZXZva2VUb2tlblJlcXVlc3QSHgoIY2xpZW50aWQYpOiG1gEgASgJUghjbGllbnRpZBImCg'
    'xjbGllbnRzZWNyZXQY97Xi7gEgASgJUgxjbGllbnRzZWNyZXQSGAoFdG9rZW4Y07fV0QEgASgJ'
    'UgV0b2tlbg==');

@$core.Deprecated('Use revokeTokenResponseDescriptor instead')
const RevokeTokenResponse$json = {
  '1': 'RevokeTokenResponse',
};

/// Descriptor for `RevokeTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeTokenResponseDescriptor =
    $convert.base64Decode('ChNSZXZva2VUb2tlblJlc3BvbnNl');

@$core.Deprecated('Use riskConfigurationTypeDescriptor instead')
const RiskConfigurationType$json = {
  '1': 'RiskConfigurationType',
  '2': [
    {
      '1': 'accounttakeoverriskconfiguration',
      '3': 184926217,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountTakeoverRiskConfigurationType',
      '10': 'accounttakeoverriskconfiguration'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'compromisedcredentialsriskconfiguration',
      '3': 85540233,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CompromisedCredentialsRiskConfigurationType',
      '10': 'compromisedcredentialsriskconfiguration'
    },
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'riskexceptionconfiguration',
      '3': 2938748,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.RiskExceptionConfigurationType',
      '10': 'riskexceptionconfiguration'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `RiskConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List riskConfigurationTypeDescriptor = $convert.base64Decode(
    'ChVSaXNrQ29uZmlndXJhdGlvblR5cGUSgAEKIGFjY291bnR0YWtlb3ZlcnJpc2tjb25maWd1cm'
    'F0aW9uGImAl1ggASgLMjEuY29nbml0b19pZHAuQWNjb3VudFRha2VvdmVyUmlza0NvbmZpZ3Vy'
    'YXRpb25UeXBlUiBhY2NvdW50dGFrZW92ZXJyaXNrY29uZmlndXJhdGlvbhIeCghjbGllbnRpZB'
    'ik6IbWASABKAlSCGNsaWVudGlkEpUBCidjb21wcm9taXNlZGNyZWRlbnRpYWxzcmlza2NvbmZp'
    'Z3VyYXRpb24YifvkKCABKAsyOC5jb2duaXRvX2lkcC5Db21wcm9taXNlZENyZWRlbnRpYWxzUm'
    'lza0NvbmZpZ3VyYXRpb25UeXBlUidjb21wcm9taXNlZGNyZWRlbnRpYWxzcmlza2NvbmZpZ3Vy'
    'YXRpb24SLQoQbGFzdG1vZGlmaWVkZGF0ZRjTiMgLIAEoCVIQbGFzdG1vZGlmaWVkZGF0ZRJuCh'
    'pyaXNrZXhjZXB0aW9uY29uZmlndXJhdGlvbhj8rrMBIAEoCzIrLmNvZ25pdG9faWRwLlJpc2tF'
    'eGNlcHRpb25Db25maWd1cmF0aW9uVHlwZVIacmlza2V4Y2VwdGlvbmNvbmZpZ3VyYXRpb24SIg'
    'oKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use riskExceptionConfigurationTypeDescriptor instead')
const RiskExceptionConfigurationType$json = {
  '1': 'RiskExceptionConfigurationType',
  '2': [
    {
      '1': 'blockediprangelist',
      '3': 31179024,
      '4': 3,
      '5': 9,
      '10': 'blockediprangelist'
    },
    {
      '1': 'skippediprangelist',
      '3': 114395508,
      '4': 3,
      '5': 9,
      '10': 'skippediprangelist'
    },
  ],
};

/// Descriptor for `RiskExceptionConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List riskExceptionConfigurationTypeDescriptor =
    $convert.base64Decode(
        'Ch5SaXNrRXhjZXB0aW9uQ29uZmlndXJhdGlvblR5cGUSMQoSYmxvY2tlZGlwcmFuZ2VsaXN0GJ'
        'CC7w4gAygJUhJibG9ja2VkaXByYW5nZWxpc3QSMQoSc2tpcHBlZGlwcmFuZ2VsaXN0GPSSxjYg'
        'AygJUhJza2lwcGVkaXByYW5nZWxpc3Q=');

@$core.Deprecated('Use s3ConfigurationTypeDescriptor instead')
const S3ConfigurationType$json = {
  '1': 'S3ConfigurationType',
  '2': [
    {'1': 'bucketarn', '3': 255683899, '4': 1, '5': 9, '10': 'bucketarn'},
  ],
};

/// Descriptor for `S3ConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3ConfigurationTypeDescriptor = $convert.base64Decode(
    'ChNTM0NvbmZpZ3VyYXRpb25UeXBlEh8KCWJ1Y2tldGFybhi72vV5IAEoCVIJYnVja2V0YXJu');

@$core.Deprecated('Use sMSMfaSettingsTypeDescriptor instead')
const SMSMfaSettingsType$json = {
  '1': 'SMSMfaSettingsType',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'preferredmfa', '3': 195248469, '4': 1, '5': 8, '10': 'preferredmfa'},
  ],
};

/// Descriptor for `SMSMfaSettingsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sMSMfaSettingsTypeDescriptor = $convert.base64Decode(
    'ChJTTVNNZmFTZXR0aW5nc1R5cGUSHAoHZW5hYmxlZBi/yJvkASABKAhSB2VuYWJsZWQSJQoMcH'
    'JlZmVycmVkbWZhGNWCjV0gASgIUgxwcmVmZXJyZWRtZmE=');

@$core.Deprecated('Use schemaAttributeTypeDescriptor instead')
const SchemaAttributeType$json = {
  '1': 'SchemaAttributeType',
  '2': [
    {
      '1': 'attributedatatype',
      '3': 400166480,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AttributeDataType',
      '10': 'attributedatatype'
    },
    {
      '1': 'developeronlyattribute',
      '3': 195850806,
      '4': 1,
      '5': 8,
      '10': 'developeronlyattribute'
    },
    {'1': 'mutable', '3': 149670730, '4': 1, '5': 8, '10': 'mutable'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'numberattributeconstraints',
      '3': 381084525,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.NumberAttributeConstraintsType',
      '10': 'numberattributeconstraints'
    },
    {'1': 'required', '3': 300200513, '4': 1, '5': 8, '10': 'required'},
    {
      '1': 'stringattributeconstraints',
      '3': 270228629,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.StringAttributeConstraintsType',
      '10': 'stringattributeconstraints'
    },
  ],
};

/// Descriptor for `SchemaAttributeType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List schemaAttributeTypeDescriptor = $convert.base64Decode(
    'ChNTY2hlbWFBdHRyaWJ1dGVUeXBlElAKEWF0dHJpYnV0ZWRhdGF0eXBlGNCc6L4BIAEoDjIeLm'
    'NvZ25pdG9faWRwLkF0dHJpYnV0ZURhdGFUeXBlUhFhdHRyaWJ1dGVkYXRhdHlwZRI5ChZkZXZl'
    'bG9wZXJvbmx5YXR0cmlidXRlGLbksV0gASgIUhZkZXZlbG9wZXJvbmx5YXR0cmlidXRlEhsKB2'
    '11dGFibGUYypavRyABKAhSB211dGFibGUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRJvChpudW1i'
    'ZXJhdHRyaWJ1dGVjb25zdHJhaW50cxjtxtu1ASABKAsyKy5jb2duaXRvX2lkcC5OdW1iZXJBdH'
    'RyaWJ1dGVDb25zdHJhaW50c1R5cGVSGm51bWJlcmF0dHJpYnV0ZWNvbnN0cmFpbnRzEh4KCHJl'
    'cXVpcmVkGMHkko8BIAEoCFIIcmVxdWlyZWQSbwoac3RyaW5nYXR0cmlidXRlY29uc3RyYWludH'
    'MYlbntgAEgASgLMisuY29nbml0b19pZHAuU3RyaW5nQXR0cmlidXRlQ29uc3RyYWludHNUeXBl'
    'UhpzdHJpbmdhdHRyaWJ1dGVjb25zdHJhaW50cw==');

@$core.Deprecated('Use scopeDoesNotExistExceptionDescriptor instead')
const ScopeDoesNotExistException$json = {
  '1': 'ScopeDoesNotExistException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ScopeDoesNotExistException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scopeDoesNotExistExceptionDescriptor =
    $convert.base64Decode(
        'ChpTY29wZURvZXNOb3RFeGlzdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use setLogDeliveryConfigurationRequestDescriptor instead')
const SetLogDeliveryConfigurationRequest$json = {
  '1': 'SetLogDeliveryConfigurationRequest',
  '2': [
    {
      '1': 'logconfigurations',
      '3': 89180603,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.LogConfigurationType',
      '10': 'logconfigurations'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `SetLogDeliveryConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setLogDeliveryConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiJTZXRMb2dEZWxpdmVyeUNvbmZpZ3VyYXRpb25SZXF1ZXN0ElIKEWxvZ2NvbmZpZ3VyYXRpb2'
        '5zGLuTwyogAygLMiEuY29nbml0b19pZHAuTG9nQ29uZmlndXJhdGlvblR5cGVSEWxvZ2NvbmZp'
        'Z3VyYXRpb25zEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use setLogDeliveryConfigurationResponseDescriptor instead')
const SetLogDeliveryConfigurationResponse$json = {
  '1': 'SetLogDeliveryConfigurationResponse',
  '2': [
    {
      '1': 'logdeliveryconfiguration',
      '3': 332495878,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.LogDeliveryConfigurationType',
      '10': 'logdeliveryconfiguration'
    },
  ],
};

/// Descriptor for `SetLogDeliveryConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setLogDeliveryConfigurationResponseDescriptor =
    $convert.base64Decode(
        'CiNTZXRMb2dEZWxpdmVyeUNvbmZpZ3VyYXRpb25SZXNwb25zZRJpChhsb2dkZWxpdmVyeWNvbm'
        'ZpZ3VyYXRpb24YhvjFngEgASgLMikuY29nbml0b19pZHAuTG9nRGVsaXZlcnlDb25maWd1cmF0'
        'aW9uVHlwZVIYbG9nZGVsaXZlcnljb25maWd1cmF0aW9u');

@$core.Deprecated('Use setRiskConfigurationRequestDescriptor instead')
const SetRiskConfigurationRequest$json = {
  '1': 'SetRiskConfigurationRequest',
  '2': [
    {
      '1': 'accounttakeoverriskconfiguration',
      '3': 184926217,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountTakeoverRiskConfigurationType',
      '10': 'accounttakeoverriskconfiguration'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'compromisedcredentialsriskconfiguration',
      '3': 85540233,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CompromisedCredentialsRiskConfigurationType',
      '10': 'compromisedcredentialsriskconfiguration'
    },
    {
      '1': 'riskexceptionconfiguration',
      '3': 2938748,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.RiskExceptionConfigurationType',
      '10': 'riskexceptionconfiguration'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `SetRiskConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setRiskConfigurationRequestDescriptor = $convert.base64Decode(
    'ChtTZXRSaXNrQ29uZmlndXJhdGlvblJlcXVlc3QSgAEKIGFjY291bnR0YWtlb3ZlcnJpc2tjb2'
    '5maWd1cmF0aW9uGImAl1ggASgLMjEuY29nbml0b19pZHAuQWNjb3VudFRha2VvdmVyUmlza0Nv'
    'bmZpZ3VyYXRpb25UeXBlUiBhY2NvdW50dGFrZW92ZXJyaXNrY29uZmlndXJhdGlvbhIeCghjbG'
    'llbnRpZBik6IbWASABKAlSCGNsaWVudGlkEpUBCidjb21wcm9taXNlZGNyZWRlbnRpYWxzcmlz'
    'a2NvbmZpZ3VyYXRpb24YifvkKCABKAsyOC5jb2duaXRvX2lkcC5Db21wcm9taXNlZENyZWRlbn'
    'RpYWxzUmlza0NvbmZpZ3VyYXRpb25UeXBlUidjb21wcm9taXNlZGNyZWRlbnRpYWxzcmlza2Nv'
    'bmZpZ3VyYXRpb24Sbgoacmlza2V4Y2VwdGlvbmNvbmZpZ3VyYXRpb24Y/K6zASABKAsyKy5jb2'
    'duaXRvX2lkcC5SaXNrRXhjZXB0aW9uQ29uZmlndXJhdGlvblR5cGVSGnJpc2tleGNlcHRpb25j'
    'b25maWd1cmF0aW9uEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use setRiskConfigurationResponseDescriptor instead')
const SetRiskConfigurationResponse$json = {
  '1': 'SetRiskConfigurationResponse',
  '2': [
    {
      '1': 'riskconfiguration',
      '3': 374237465,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.RiskConfigurationType',
      '10': 'riskconfiguration'
    },
  ],
};

/// Descriptor for `SetRiskConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setRiskConfigurationResponseDescriptor =
    $convert.base64Decode(
        'ChxTZXRSaXNrQ29uZmlndXJhdGlvblJlc3BvbnNlElQKEXJpc2tjb25maWd1cmF0aW9uGJnSub'
        'IBIAEoCzIiLmNvZ25pdG9faWRwLlJpc2tDb25maWd1cmF0aW9uVHlwZVIRcmlza2NvbmZpZ3Vy'
        'YXRpb24=');

@$core.Deprecated('Use setUICustomizationRequestDescriptor instead')
const SetUICustomizationRequest$json = {
  '1': 'SetUICustomizationRequest',
  '2': [
    {'1': 'css', '3': 371300961, '4': 1, '5': 9, '10': 'css'},
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'imagefile', '3': 372468861, '4': 1, '5': 12, '10': 'imagefile'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `SetUICustomizationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUICustomizationRequestDescriptor = $convert.base64Decode(
    'ChlTZXRVSUN1c3RvbWl6YXRpb25SZXF1ZXN0EhQKA2NzcxjhtIaxASABKAlSA2NzcxIeCghjbG'
    'llbnRpZBik6IbWASABKAlSCGNsaWVudGlkEiAKCWltYWdlZmlsZRj92M2xASABKAxSCWltYWdl'
    'ZmlsZRIiCgp1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use setUICustomizationResponseDescriptor instead')
const SetUICustomizationResponse$json = {
  '1': 'SetUICustomizationResponse',
  '2': [
    {
      '1': 'uicustomization',
      '3': 28966753,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UICustomizationType',
      '10': 'uicustomization'
    },
  ],
};

/// Descriptor for `SetUICustomizationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUICustomizationResponseDescriptor =
    $convert.base64Decode(
        'ChpTZXRVSUN1c3RvbWl6YXRpb25SZXNwb25zZRJNCg91aWN1c3RvbWl6YXRpb24Y4f7nDSABKA'
        'syIC5jb2duaXRvX2lkcC5VSUN1c3RvbWl6YXRpb25UeXBlUg91aWN1c3RvbWl6YXRpb24=');

@$core.Deprecated('Use setUserMFAPreferenceRequestDescriptor instead')
const SetUserMFAPreferenceRequest$json = {
  '1': 'SetUserMFAPreferenceRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'emailmfasettings',
      '3': 306207991,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailMfaSettingsType',
      '10': 'emailmfasettings'
    },
    {
      '1': 'smsmfasettings',
      '3': 220659382,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SMSMfaSettingsType',
      '10': 'smsmfasettings'
    },
    {
      '1': 'softwaretokenmfasettings',
      '3': 22910957,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SoftwareTokenMfaSettingsType',
      '10': 'softwaretokenmfasettings'
    },
  ],
};

/// Descriptor for `SetUserMFAPreferenceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUserMFAPreferenceRequestDescriptor = $convert.base64Decode(
    'ChtTZXRVc2VyTUZBUHJlZmVyZW5jZVJlcXVlc3QSIwoLYWNjZXNzdG9rZW4YibyQRiABKAlSC2'
    'FjY2Vzc3Rva2VuElEKEGVtYWlsbWZhc2V0dGluZ3MY97mBkgEgASgLMiEuY29nbml0b19pZHAu'
    'RW1haWxNZmFTZXR0aW5nc1R5cGVSEGVtYWlsbWZhc2V0dGluZ3MSSgoOc21zbWZhc2V0dGluZ3'
    'MYtv2baSABKAsyHy5jb2duaXRvX2lkcC5TTVNNZmFTZXR0aW5nc1R5cGVSDnNtc21mYXNldHRp'
    'bmdzEmgKGHNvZnR3YXJldG9rZW5tZmFzZXR0aW5ncxjtr/YKIAEoCzIpLmNvZ25pdG9faWRwLl'
    'NvZnR3YXJlVG9rZW5NZmFTZXR0aW5nc1R5cGVSGHNvZnR3YXJldG9rZW5tZmFzZXR0aW5ncw==');

@$core.Deprecated('Use setUserMFAPreferenceResponseDescriptor instead')
const SetUserMFAPreferenceResponse$json = {
  '1': 'SetUserMFAPreferenceResponse',
};

/// Descriptor for `SetUserMFAPreferenceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUserMFAPreferenceResponseDescriptor =
    $convert.base64Decode('ChxTZXRVc2VyTUZBUHJlZmVyZW5jZVJlc3BvbnNl');

@$core.Deprecated('Use setUserPoolMfaConfigRequestDescriptor instead')
const SetUserPoolMfaConfigRequest$json = {
  '1': 'SetUserPoolMfaConfigRequest',
  '2': [
    {
      '1': 'emailmfaconfiguration',
      '3': 482754548,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailMfaConfigType',
      '10': 'emailmfaconfiguration'
    },
    {
      '1': 'mfaconfiguration',
      '3': 259020766,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolMfaType',
      '10': 'mfaconfiguration'
    },
    {
      '1': 'smsmfaconfiguration',
      '3': 153073099,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SmsMfaConfigType',
      '10': 'smsmfaconfiguration'
    },
    {
      '1': 'softwaretokenmfaconfiguration',
      '3': 502085950,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SoftwareTokenMfaConfigType',
      '10': 'softwaretokenmfaconfiguration'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {
      '1': 'webauthnconfiguration',
      '3': 506289104,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.WebAuthnConfigurationType',
      '10': 'webauthnconfiguration'
    },
  ],
};

/// Descriptor for `SetUserPoolMfaConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUserPoolMfaConfigRequestDescriptor = $convert.base64Decode(
    'ChtTZXRVc2VyUG9vbE1mYUNvbmZpZ1JlcXVlc3QSWQoVZW1haWxtZmFjb25maWd1cmF0aW9uGP'
    'T/mOYBIAEoCzIfLmNvZ25pdG9faWRwLkVtYWlsTWZhQ29uZmlnVHlwZVIVZW1haWxtZmFjb25m'
    'aWd1cmF0aW9uEksKEG1mYWNvbmZpZ3VyYXRpb24Y3q/BeyABKA4yHC5jb2duaXRvX2lkcC5Vc2'
    'VyUG9vbE1mYVR5cGVSEG1mYWNvbmZpZ3VyYXRpb24SUgoTc21zbWZhY29uZmlndXJhdGlvbhjL'
    '6/5IIAEoCzIdLmNvZ25pdG9faWRwLlNtc01mYUNvbmZpZ1R5cGVSE3Ntc21mYWNvbmZpZ3VyYX'
    'Rpb24ScQodc29mdHdhcmV0b2tlbm1mYWNvbmZpZ3VyYXRpb24YvvK07wEgASgLMicuY29nbml0'
    'b19pZHAuU29mdHdhcmVUb2tlbk1mYUNvbmZpZ1R5cGVSHXNvZnR3YXJldG9rZW5tZmFjb25maW'
    'd1cmF0aW9uEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlkEmAKFXdlYmF1dGhu'
    'Y29uZmlndXJhdGlvbhjQt7XxASABKAsyJi5jb2duaXRvX2lkcC5XZWJBdXRobkNvbmZpZ3VyYX'
    'Rpb25UeXBlUhV3ZWJhdXRobmNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use setUserPoolMfaConfigResponseDescriptor instead')
const SetUserPoolMfaConfigResponse$json = {
  '1': 'SetUserPoolMfaConfigResponse',
  '2': [
    {
      '1': 'emailmfaconfiguration',
      '3': 482754548,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailMfaConfigType',
      '10': 'emailmfaconfiguration'
    },
    {
      '1': 'mfaconfiguration',
      '3': 259020766,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolMfaType',
      '10': 'mfaconfiguration'
    },
    {
      '1': 'smsmfaconfiguration',
      '3': 153073099,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SmsMfaConfigType',
      '10': 'smsmfaconfiguration'
    },
    {
      '1': 'softwaretokenmfaconfiguration',
      '3': 502085950,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SoftwareTokenMfaConfigType',
      '10': 'softwaretokenmfaconfiguration'
    },
    {
      '1': 'webauthnconfiguration',
      '3': 506289104,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.WebAuthnConfigurationType',
      '10': 'webauthnconfiguration'
    },
  ],
};

/// Descriptor for `SetUserPoolMfaConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUserPoolMfaConfigResponseDescriptor = $convert.base64Decode(
    'ChxTZXRVc2VyUG9vbE1mYUNvbmZpZ1Jlc3BvbnNlElkKFWVtYWlsbWZhY29uZmlndXJhdGlvbh'
    'j0/5jmASABKAsyHy5jb2duaXRvX2lkcC5FbWFpbE1mYUNvbmZpZ1R5cGVSFWVtYWlsbWZhY29u'
    'ZmlndXJhdGlvbhJLChBtZmFjb25maWd1cmF0aW9uGN6vwXsgASgOMhwuY29nbml0b19pZHAuVX'
    'NlclBvb2xNZmFUeXBlUhBtZmFjb25maWd1cmF0aW9uElIKE3Ntc21mYWNvbmZpZ3VyYXRpb24Y'
    'y+v+SCABKAsyHS5jb2duaXRvX2lkcC5TbXNNZmFDb25maWdUeXBlUhNzbXNtZmFjb25maWd1cm'
    'F0aW9uEnEKHXNvZnR3YXJldG9rZW5tZmFjb25maWd1cmF0aW9uGL7ytO8BIAEoCzInLmNvZ25p'
    'dG9faWRwLlNvZnR3YXJlVG9rZW5NZmFDb25maWdUeXBlUh1zb2Z0d2FyZXRva2VubWZhY29uZm'
    'lndXJhdGlvbhJgChV3ZWJhdXRobmNvbmZpZ3VyYXRpb24Y0Le18QEgASgLMiYuY29nbml0b19p'
    'ZHAuV2ViQXV0aG5Db25maWd1cmF0aW9uVHlwZVIVd2ViYXV0aG5jb25maWd1cmF0aW9u');

@$core.Deprecated('Use setUserSettingsRequestDescriptor instead')
const SetUserSettingsRequest$json = {
  '1': 'SetUserSettingsRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'mfaoptions',
      '3': 501540826,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.MFAOptionType',
      '10': 'mfaoptions'
    },
  ],
};

/// Descriptor for `SetUserSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUserSettingsRequestDescriptor = $convert.base64Decode(
    'ChZTZXRVc2VyU2V0dGluZ3NSZXF1ZXN0EiMKC2FjY2Vzc3Rva2VuGIm8kEYgASgJUgthY2Nlc3'
    'N0b2tlbhI+CgptZmFvcHRpb25zGNrPk+8BIAMoCzIaLmNvZ25pdG9faWRwLk1GQU9wdGlvblR5'
    'cGVSCm1mYW9wdGlvbnM=');

@$core.Deprecated('Use setUserSettingsResponseDescriptor instead')
const SetUserSettingsResponse$json = {
  '1': 'SetUserSettingsResponse',
};

/// Descriptor for `SetUserSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setUserSettingsResponseDescriptor =
    $convert.base64Decode('ChdTZXRVc2VyU2V0dGluZ3NSZXNwb25zZQ==');

@$core.Deprecated('Use signInPolicyTypeDescriptor instead')
const SignInPolicyType$json = {
  '1': 'SignInPolicyType',
  '2': [
    {
      '1': 'allowedfirstauthfactors',
      '3': 315187180,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.AuthFactorType',
      '10': 'allowedfirstauthfactors'
    },
  ],
};

/// Descriptor for `SignInPolicyType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signInPolicyTypeDescriptor = $convert.base64Decode(
    'ChBTaWduSW5Qb2xpY3lUeXBlElkKF2FsbG93ZWRmaXJzdGF1dGhmYWN0b3JzGOy/pZYBIAMoDj'
    'IbLmNvZ25pdG9faWRwLkF1dGhGYWN0b3JUeXBlUhdhbGxvd2VkZmlyc3RhdXRoZmFjdG9ycw==');

@$core.Deprecated('Use signUpRequestDescriptor instead')
const SignUpRequest$json = {
  '1': 'SignUpRequest',
  '2': [
    {
      '1': 'analyticsmetadata',
      '3': 126894839,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsMetadataType',
      '10': 'analyticsmetadata'
    },
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.SignUpRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {'1': 'password', '3': 214108217, '4': 1, '5': 9, '10': 'password'},
    {'1': 'secrethash', '3': 321025630, '4': 1, '5': 9, '10': 'secrethash'},
    {
      '1': 'userattributes',
      '3': 194667064,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'userattributes'
    },
    {
      '1': 'usercontextdata',
      '3': 169951134,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserContextDataType',
      '10': 'usercontextdata'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
    {
      '1': 'validationdata',
      '3': 235118411,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'validationdata'
    },
  ],
  '3': [SignUpRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use signUpRequestDescriptor instead')
const SignUpRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SignUpRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signUpRequestDescriptor = $convert.base64Decode(
    'Cg1TaWduVXBSZXF1ZXN0ElMKEWFuYWx5dGljc21ldGFkYXRhGPeFwTwgASgLMiIuY29nbml0b1'
    '9pZHAuQW5hbHl0aWNzTWV0YWRhdGFUeXBlUhFhbmFseXRpY3NtZXRhZGF0YRIeCghjbGllbnRp'
    'ZBik6IbWASABKAlSCGNsaWVudGlkElkKDmNsaWVudG1ldGFkYXRhGMyv/2EgAygLMi4uY29nbm'
    'l0b19pZHAuU2lnblVwUmVxdWVzdC5DbGllbnRtZXRhZGF0YUVudHJ5Ug5jbGllbnRtZXRhZGF0'
    'YRIdCghwYXNzd29yZBi5kIxmIAEoCVIIcGFzc3dvcmQSIgoKc2VjcmV0aGFzaBje7ImZASABKA'
    'lSCnNlY3JldGhhc2gSRQoOdXNlcmF0dHJpYnV0ZXMYuMTpXCADKAsyGi5jb2duaXRvX2lkcC5B'
    'dHRyaWJ1dGVUeXBlUg51c2VyYXR0cmlidXRlcxJNCg91c2VyY29udGV4dGRhdGEYnv+EUSABKA'
    'syIC5jb2duaXRvX2lkcC5Vc2VyQ29udGV4dERhdGFUeXBlUg91c2VyY29udGV4dGRhdGESHgoI'
    'dXNlcm5hbWUY2qmj4AEgASgJUgh1c2VybmFtZRJFCg52YWxpZGF0aW9uZGF0YRjLvo5wIAMoCz'
    'IaLmNvZ25pdG9faWRwLkF0dHJpYnV0ZVR5cGVSDnZhbGlkYXRpb25kYXRhGkEKE0NsaWVudG1l'
    'dGFkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ'
    '==');

@$core.Deprecated('Use signUpResponseDescriptor instead')
const SignUpResponse$json = {
  '1': 'SignUpResponse',
  '2': [
    {
      '1': 'codedeliverydetails',
      '3': 423272803,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CodeDeliveryDetailsType',
      '10': 'codedeliverydetails'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {
      '1': 'userconfirmed',
      '3': 200708940,
      '4': 1,
      '5': 8,
      '10': 'userconfirmed'
    },
    {'1': 'usersub', '3': 171046417, '4': 1, '5': 9, '10': 'usersub'},
  ],
};

/// Descriptor for `SignUpResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signUpResponseDescriptor = $convert.base64Decode(
    'Cg5TaWduVXBSZXNwb25zZRJaChNjb2RlZGVsaXZlcnlkZXRhaWxzGOPC6skBIAEoCzIkLmNvZ2'
    '5pdG9faWRwLkNvZGVEZWxpdmVyeURldGFpbHNUeXBlUhNjb2RlZGVsaXZlcnlkZXRhaWxzEhsK'
    'B3Nlc3Npb24YmJmjAiABKAlSB3Nlc3Npb24SJwoNdXNlcmNvbmZpcm1lZBjMptpfIAEoCFINdX'
    'NlcmNvbmZpcm1lZBIbCgd1c2Vyc3ViGJHsx1EgASgJUgd1c2Vyc3Vi');

@$core.Deprecated('Use smsConfigurationTypeDescriptor instead')
const SmsConfigurationType$json = {
  '1': 'SmsConfigurationType',
  '2': [
    {'1': 'externalid', '3': 271401992, '4': 1, '5': 9, '10': 'externalid'},
    {'1': 'snscallerarn', '3': 128123070, '4': 1, '5': 9, '10': 'snscallerarn'},
    {'1': 'snsregion', '3': 413686894, '4': 1, '5': 9, '10': 'snsregion'},
  ],
};

/// Descriptor for `SmsConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List smsConfigurationTypeDescriptor = $convert.base64Decode(
    'ChRTbXNDb25maWd1cmF0aW9uVHlwZRIiCgpleHRlcm5hbGlkGIiItYEBIAEoCVIKZXh0ZXJuYW'
    'xpZBIlCgxzbnNjYWxsZXJhcm4YvoGMPSABKAlSDHNuc2NhbGxlcmFybhIgCglzbnNyZWdpb24Y'
    '7rihxQEgASgJUglzbnNyZWdpb24=');

@$core.Deprecated('Use smsMfaConfigTypeDescriptor instead')
const SmsMfaConfigType$json = {
  '1': 'SmsMfaConfigType',
  '2': [
    {
      '1': 'smsauthenticationmessage',
      '3': 356104990,
      '4': 1,
      '5': 9,
      '10': 'smsauthenticationmessage'
    },
    {
      '1': 'smsconfiguration',
      '3': 10321849,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SmsConfigurationType',
      '10': 'smsconfiguration'
    },
  ],
};

/// Descriptor for `SmsMfaConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List smsMfaConfigTypeDescriptor = $convert.base64Decode(
    'ChBTbXNNZmFDb25maWdUeXBlEj4KGHNtc2F1dGhlbnRpY2F0aW9ubWVzc2FnZRie9uapASABKA'
    'lSGHNtc2F1dGhlbnRpY2F0aW9ubWVzc2FnZRJQChBzbXNjb25maWd1cmF0aW9uGLn/9QQgASgL'
    'MiEuY29nbml0b19pZHAuU21zQ29uZmlndXJhdGlvblR5cGVSEHNtc2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use softwareTokenMFANotFoundExceptionDescriptor instead')
const SoftwareTokenMFANotFoundException$json = {
  '1': 'SoftwareTokenMFANotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SoftwareTokenMFANotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List softwareTokenMFANotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CiFTb2Z0d2FyZVRva2VuTUZBTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use softwareTokenMfaConfigTypeDescriptor instead')
const SoftwareTokenMfaConfigType$json = {
  '1': 'SoftwareTokenMfaConfigType',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
  ],
};

/// Descriptor for `SoftwareTokenMfaConfigType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List softwareTokenMfaConfigTypeDescriptor =
    $convert.base64Decode(
        'ChpTb2Z0d2FyZVRva2VuTWZhQ29uZmlnVHlwZRIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYm'
        'xlZA==');

@$core.Deprecated('Use softwareTokenMfaSettingsTypeDescriptor instead')
const SoftwareTokenMfaSettingsType$json = {
  '1': 'SoftwareTokenMfaSettingsType',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'preferredmfa', '3': 195248469, '4': 1, '5': 8, '10': 'preferredmfa'},
  ],
};

/// Descriptor for `SoftwareTokenMfaSettingsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List softwareTokenMfaSettingsTypeDescriptor =
    $convert.base64Decode(
        'ChxTb2Z0d2FyZVRva2VuTWZhU2V0dGluZ3NUeXBlEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbm'
        'FibGVkEiUKDHByZWZlcnJlZG1mYRjVgo1dIAEoCFIMcHJlZmVycmVkbWZh');

@$core.Deprecated('Use startUserImportJobRequestDescriptor instead')
const StartUserImportJobRequest$json = {
  '1': 'StartUserImportJobRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `StartUserImportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startUserImportJobRequestDescriptor =
    $convert.base64Decode(
        'ChlTdGFydFVzZXJJbXBvcnRKb2JSZXF1ZXN0EhcKBWpvYmlkGNLU3TMgASgJUgVqb2JpZBIiCg'
        'p1c2VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZA==');

@$core.Deprecated('Use startUserImportJobResponseDescriptor instead')
const StartUserImportJobResponse$json = {
  '1': 'StartUserImportJobResponse',
  '2': [
    {
      '1': 'userimportjob',
      '3': 473682999,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserImportJobType',
      '10': 'userimportjob'
    },
  ],
};

/// Descriptor for `StartUserImportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startUserImportJobResponseDescriptor =
    $convert.base64Decode(
        'ChpTdGFydFVzZXJJbXBvcnRKb2JSZXNwb25zZRJICg11c2VyaW1wb3J0am9iGLeo7+EBIAEoCz'
        'IeLmNvZ25pdG9faWRwLlVzZXJJbXBvcnRKb2JUeXBlUg11c2VyaW1wb3J0am9i');

@$core.Deprecated('Use startWebAuthnRegistrationRequestDescriptor instead')
const StartWebAuthnRegistrationRequest$json = {
  '1': 'StartWebAuthnRegistrationRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
  ],
};

/// Descriptor for `StartWebAuthnRegistrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startWebAuthnRegistrationRequestDescriptor =
    $convert.base64Decode(
        'CiBTdGFydFdlYkF1dGhuUmVnaXN0cmF0aW9uUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIA'
        'EoCVILYWNjZXNzdG9rZW4=');

@$core.Deprecated('Use startWebAuthnRegistrationResponseDescriptor instead')
const StartWebAuthnRegistrationResponse$json = {
  '1': 'StartWebAuthnRegistrationResponse',
  '2': [
    {
      '1': 'credentialcreationoptions',
      '3': 66218090,
      '4': 1,
      '5': 9,
      '10': 'credentialcreationoptions'
    },
  ],
};

/// Descriptor for `StartWebAuthnRegistrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startWebAuthnRegistrationResponseDescriptor =
    $convert.base64Decode(
        'CiFTdGFydFdlYkF1dGhuUmVnaXN0cmF0aW9uUmVzcG9uc2USPwoZY3JlZGVudGlhbGNyZWF0aW'
        '9ub3B0aW9ucxjq0MkfIAEoCVIZY3JlZGVudGlhbGNyZWF0aW9ub3B0aW9ucw==');

@$core.Deprecated('Use stopUserImportJobRequestDescriptor instead')
const StopUserImportJobRequest$json = {
  '1': 'StopUserImportJobRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `StopUserImportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopUserImportJobRequestDescriptor =
    $convert.base64Decode(
        'ChhTdG9wVXNlckltcG9ydEpvYlJlcXVlc3QSFwoFam9iaWQY0tTdMyABKAlSBWpvYmlkEiIKCn'
        'VzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use stopUserImportJobResponseDescriptor instead')
const StopUserImportJobResponse$json = {
  '1': 'StopUserImportJobResponse',
  '2': [
    {
      '1': 'userimportjob',
      '3': 473682999,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserImportJobType',
      '10': 'userimportjob'
    },
  ],
};

/// Descriptor for `StopUserImportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopUserImportJobResponseDescriptor =
    $convert.base64Decode(
        'ChlTdG9wVXNlckltcG9ydEpvYlJlc3BvbnNlEkgKDXVzZXJpbXBvcnRqb2IYt6jv4QEgASgLMh'
        '4uY29nbml0b19pZHAuVXNlckltcG9ydEpvYlR5cGVSDXVzZXJpbXBvcnRqb2I=');

@$core.Deprecated('Use stringAttributeConstraintsTypeDescriptor instead')
const StringAttributeConstraintsType$json = {
  '1': 'StringAttributeConstraintsType',
  '2': [
    {'1': 'maxlength', '3': 490672418, '4': 1, '5': 9, '10': 'maxlength'},
    {'1': 'minlength', '3': 46267224, '4': 1, '5': 9, '10': 'minlength'},
  ],
};

/// Descriptor for `StringAttributeConstraintsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stringAttributeConstraintsTypeDescriptor =
    $convert.base64Decode(
        'Ch5TdHJpbmdBdHRyaWJ1dGVDb25zdHJhaW50c1R5cGUSIAoJbWF4bGVuZ3RoGKKi/OkBIAEoCV'
        'IJbWF4bGVuZ3RoEh8KCW1pbmxlbmd0aBjY9ocWIAEoCVIJbWlubGVuZ3Ro');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.TagResourceRequest.TagsEntry',
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
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJjZW'
    'FybhJBCgR0YWdzGMHB9rUBIAMoCzIpLmNvZ25pdG9faWRwLlRhZ1Jlc291cmNlUmVxdWVzdC5U'
    'YWdzRW50cnlSBHRhZ3MaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGA'
    'IgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use tagResourceResponseDescriptor instead')
const TagResourceResponse$json = {
  '1': 'TagResourceResponse',
};

/// Descriptor for `TagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceResponseDescriptor =
    $convert.base64Decode('ChNUYWdSZXNvdXJjZVJlc3BvbnNl');

@$core.Deprecated('Use termsDescriptionTypeDescriptor instead')
const TermsDescriptionType$json = {
  '1': 'TermsDescriptionType',
  '2': [
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'enforcement',
      '3': 412213242,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TermsEnforcementType',
      '10': 'enforcement'
    },
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'termsid', '3': 331306210, '4': 1, '5': 9, '10': 'termsid'},
    {'1': 'termsname', '3': 303051560, '4': 1, '5': 9, '10': 'termsname'},
  ],
};

/// Descriptor for `TermsDescriptionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List termsDescriptionTypeDescriptor = $convert.base64Decode(
    'ChRUZXJtc0Rlc2NyaXB0aW9uVHlwZRImCgxjcmVhdGlvbmRhdGUY4di3iQEgASgJUgxjcmVhdG'
    'lvbmRhdGUSRwoLZW5mb3JjZW1lbnQY+r/HxAEgASgOMiEuY29nbml0b19pZHAuVGVybXNFbmZv'
    'cmNlbWVudFR5cGVSC2VuZm9yY2VtZW50Ei0KEGxhc3Rtb2RpZmllZGRhdGUY04jICyABKAlSEG'
    'xhc3Rtb2RpZmllZGRhdGUSHAoHdGVybXNpZBjiqf2dASABKAlSB3Rlcm1zaWQSIAoJdGVybXNu'
    'YW1lGKjmwJABIAEoCVIJdGVybXNuYW1l');

@$core.Deprecated('Use termsExistsExceptionDescriptor instead')
const TermsExistsException$json = {
  '1': 'TermsExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TermsExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List termsExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChRUZXJtc0V4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use termsTypeDescriptor instead')
const TermsType$json = {
  '1': 'TermsType',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'enforcement',
      '3': 412213242,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TermsEnforcementType',
      '10': 'enforcement'
    },
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'links',
      '3': 302123151,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.TermsType.LinksEntry',
      '10': 'links'
    },
    {'1': 'termsid', '3': 331306210, '4': 1, '5': 9, '10': 'termsid'},
    {'1': 'termsname', '3': 303051560, '4': 1, '5': 9, '10': 'termsname'},
    {
      '1': 'termssource',
      '3': 122689594,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TermsSourceType',
      '10': 'termssource'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [TermsType_LinksEntry$json],
};

@$core.Deprecated('Use termsTypeDescriptor instead')
const TermsType_LinksEntry$json = {
  '1': 'LinksEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TermsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List termsTypeDescriptor = $convert.base64Decode(
    'CglUZXJtc1R5cGUSHgoIY2xpZW50aWQYpOiG1gEgASgJUghjbGllbnRpZBImCgxjcmVhdGlvbm'
    'RhdGUY4di3iQEgASgJUgxjcmVhdGlvbmRhdGUSRwoLZW5mb3JjZW1lbnQY+r/HxAEgASgOMiEu'
    'Y29nbml0b19pZHAuVGVybXNFbmZvcmNlbWVudFR5cGVSC2VuZm9yY2VtZW50Ei0KEGxhc3Rtb2'
    'RpZmllZGRhdGUY04jICyABKAlSEGxhc3Rtb2RpZmllZGRhdGUSOwoFbGlua3MYj5GIkAEgAygL'
    'MiEuY29nbml0b19pZHAuVGVybXNUeXBlLkxpbmtzRW50cnlSBWxpbmtzEhwKB3Rlcm1zaWQY4q'
    'n9nQEgASgJUgd0ZXJtc2lkEiAKCXRlcm1zbmFtZRio5sCQASABKAlSCXRlcm1zbmFtZRJBCgt0'
    'ZXJtc3NvdXJjZRi6sMA6IAEoDjIcLmNvZ25pdG9faWRwLlRlcm1zU291cmNlVHlwZVILdGVybX'
    'Nzb3VyY2USIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQaOAoKTGlua3NFbnRy'
    'eRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use tierChangeNotAllowedExceptionDescriptor instead')
const TierChangeNotAllowedException$json = {
  '1': 'TierChangeNotAllowedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TierChangeNotAllowedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tierChangeNotAllowedExceptionDescriptor =
    $convert.base64Decode(
        'Ch1UaWVyQ2hhbmdlTm90QWxsb3dlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use tokenValidityUnitsTypeDescriptor instead')
const TokenValidityUnitsType$json = {
  '1': 'TokenValidityUnitsType',
  '2': [
    {
      '1': 'accesstoken',
      '3': 147070473,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TimeUnitsType',
      '10': 'accesstoken'
    },
    {
      '1': 'idtoken',
      '3': 228470,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TimeUnitsType',
      '10': 'idtoken'
    },
    {
      '1': 'refreshtoken',
      '3': 253777778,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TimeUnitsType',
      '10': 'refreshtoken'
    },
  ],
};

/// Descriptor for `TokenValidityUnitsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tokenValidityUnitsTypeDescriptor = $convert.base64Decode(
    'ChZUb2tlblZhbGlkaXR5VW5pdHNUeXBlEj8KC2FjY2Vzc3Rva2VuGIm8kEYgASgOMhouY29nbm'
    'l0b19pZHAuVGltZVVuaXRzVHlwZVILYWNjZXNzdG9rZW4SNgoHaWR0b2tlbhj2+A0gASgOMhou'
    'Y29nbml0b19pZHAuVGltZVVuaXRzVHlwZVIHaWR0b2tlbhJBCgxyZWZyZXNodG9rZW4Y8q6BeS'
    'ABKA4yGi5jb2duaXRvX2lkcC5UaW1lVW5pdHNUeXBlUgxyZWZyZXNodG9rZW4=');

@$core.Deprecated('Use tooManyFailedAttemptsExceptionDescriptor instead')
const TooManyFailedAttemptsException$json = {
  '1': 'TooManyFailedAttemptsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyFailedAttemptsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyFailedAttemptsExceptionDescriptor =
    $convert.base64Decode(
        'Ch5Ub29NYW55RmFpbGVkQXR0ZW1wdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use tooManyRequestsExceptionDescriptor instead')
const TooManyRequestsException$json = {
  '1': 'TooManyRequestsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyRequestsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyRequestsExceptionDescriptor =
    $convert.base64Decode(
        'ChhUb29NYW55UmVxdWVzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use uICustomizationTypeDescriptor instead')
const UICustomizationType$json = {
  '1': 'UICustomizationType',
  '2': [
    {'1': 'css', '3': 371300961, '4': 1, '5': 9, '10': 'css'},
    {'1': 'cssversion', '3': 307918093, '4': 1, '5': 9, '10': 'cssversion'},
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'imageurl', '3': 361905604, '4': 1, '5': 9, '10': 'imageurl'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `UICustomizationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uICustomizationTypeDescriptor = $convert.base64Decode(
    'ChNVSUN1c3RvbWl6YXRpb25UeXBlEhQKA2NzcxjhtIaxASABKAlSA2NzcxIiCgpjc3N2ZXJzaW'
    '9uGI3q6ZIBIAEoCVIKY3NzdmVyc2lvbhIeCghjbGllbnRpZBik6IbWASABKAlSCGNsaWVudGlk'
    'EiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9uZGF0ZRIeCghpbWFnZXVybBjE+8'
    'isASABKAlSCGltYWdldXJsEi0KEGxhc3Rtb2RpZmllZGRhdGUY04jICyABKAlSEGxhc3Rtb2Rp'
    'ZmllZGRhdGUSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQ=');

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

@$core.Deprecated('Use unexpectedLambdaExceptionDescriptor instead')
const UnexpectedLambdaException$json = {
  '1': 'UnexpectedLambdaException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnexpectedLambdaException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unexpectedLambdaExceptionDescriptor =
    $convert.base64Decode(
        'ChlVbmV4cGVjdGVkTGFtYmRhRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use unsupportedIdentityProviderExceptionDescriptor instead')
const UnsupportedIdentityProviderException$json = {
  '1': 'UnsupportedIdentityProviderException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedIdentityProviderException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedIdentityProviderExceptionDescriptor =
    $convert.base64Decode(
        'CiRVbnN1cHBvcnRlZElkZW50aXR5UHJvdmlkZXJFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use unsupportedOperationExceptionDescriptor instead')
const UnsupportedOperationException$json = {
  '1': 'UnsupportedOperationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedOperationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedOperationExceptionDescriptor =
    $convert.base64Decode(
        'Ch1VbnN1cHBvcnRlZE9wZXJhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use unsupportedTokenTypeExceptionDescriptor instead')
const UnsupportedTokenTypeException$json = {
  '1': 'UnsupportedTokenTypeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedTokenTypeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedTokenTypeExceptionDescriptor =
    $convert.base64Decode(
        'Ch1VbnN1cHBvcnRlZFRva2VuVHlwZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use unsupportedUserStateExceptionDescriptor instead')
const UnsupportedUserStateException$json = {
  '1': 'UnsupportedUserStateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedUserStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedUserStateExceptionDescriptor =
    $convert.base64Decode(
        'Ch1VbnN1cHBvcnRlZFVzZXJTdGF0ZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc291cm'
    'NlYXJuEhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagResourceResponseDescriptor instead')
const UntagResourceResponse$json = {
  '1': 'UntagResourceResponse',
};

/// Descriptor for `UntagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceResponseDescriptor =
    $convert.base64Decode('ChVVbnRhZ1Jlc291cmNlUmVzcG9uc2U=');

@$core.Deprecated('Use updateAuthEventFeedbackRequestDescriptor instead')
const UpdateAuthEventFeedbackRequest$json = {
  '1': 'UpdateAuthEventFeedbackRequest',
  '2': [
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 9, '10': 'eventid'},
    {
      '1': 'feedbacktoken',
      '3': 253886296,
      '4': 1,
      '5': 9,
      '10': 'feedbacktoken'
    },
    {
      '1': 'feedbackvalue',
      '3': 259452876,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.FeedbackValueType',
      '10': 'feedbackvalue'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UpdateAuthEventFeedbackRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAuthEventFeedbackRequestDescriptor = $convert.base64Decode(
    'Ch5VcGRhdGVBdXRoRXZlbnRGZWVkYmFja1JlcXVlc3QSHAoHZXZlbnRpZBjTlt2zASABKAlSB2'
    'V2ZW50aWQSJwoNZmVlZGJhY2t0b2tlbhjY/od5IAEoCVINZmVlZGJhY2t0b2tlbhJHCg1mZWVk'
    'YmFja3ZhbHVlGMzf23sgASgOMh4uY29nbml0b19pZHAuRmVlZGJhY2tWYWx1ZVR5cGVSDWZlZW'
    'RiYWNrdmFsdWUSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQSHgoIdXNlcm5h'
    'bWUY2qmj4AEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use updateAuthEventFeedbackResponseDescriptor instead')
const UpdateAuthEventFeedbackResponse$json = {
  '1': 'UpdateAuthEventFeedbackResponse',
};

/// Descriptor for `UpdateAuthEventFeedbackResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAuthEventFeedbackResponseDescriptor =
    $convert.base64Decode('Ch9VcGRhdGVBdXRoRXZlbnRGZWVkYmFja1Jlc3BvbnNl');

@$core.Deprecated('Use updateDeviceStatusRequestDescriptor instead')
const UpdateDeviceStatusRequest$json = {
  '1': 'UpdateDeviceStatusRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {'1': 'devicekey', '3': 382874155, '4': 1, '5': 9, '10': 'devicekey'},
    {
      '1': 'devicerememberedstatus',
      '3': 111455992,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DeviceRememberedStatusType',
      '10': 'devicerememberedstatus'
    },
  ],
};

/// Descriptor for `UpdateDeviceStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDeviceStatusRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVEZXZpY2VTdGF0dXNSZXF1ZXN0EiMKC2FjY2Vzc3Rva2VuGIm8kEYgASgJUgthY2'
    'Nlc3N0b2tlbhIgCglkZXZpY2VrZXkYq+TItgEgASgJUglkZXZpY2VrZXkSYgoWZGV2aWNlcmVt'
    'ZW1iZXJlZHN0YXR1cxj43ZI1IAEoDjInLmNvZ25pdG9faWRwLkRldmljZVJlbWVtYmVyZWRTdG'
    'F0dXNUeXBlUhZkZXZpY2VyZW1lbWJlcmVkc3RhdHVz');

@$core.Deprecated('Use updateDeviceStatusResponseDescriptor instead')
const UpdateDeviceStatusResponse$json = {
  '1': 'UpdateDeviceStatusResponse',
};

/// Descriptor for `UpdateDeviceStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDeviceStatusResponseDescriptor =
    $convert.base64Decode('ChpVcGRhdGVEZXZpY2VTdGF0dXNSZXNwb25zZQ==');

@$core.Deprecated('Use updateGroupRequestDescriptor instead')
const UpdateGroupRequest$json = {
  '1': 'UpdateGroupRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'precedence', '3': 206704142, '4': 1, '5': 5, '10': 'precedence'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `UpdateGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGroupRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVHcm91cFJlcXVlc3QSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW'
    '9uEiAKCWdyb3VwbmFtZRjIyqCqASABKAlSCWdyb3VwbmFtZRIhCgpwcmVjZWRlbmNlGI6cyGIg'
    'ASgFUgpwcmVjZWRlbmNlEhwKB3JvbGVhcm4YgfjnmQEgASgJUgdyb2xlYXJuEiIKCnVzZXJwb2'
    '9saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use updateGroupResponseDescriptor instead')
const UpdateGroupResponse$json = {
  '1': 'UpdateGroupResponse',
  '2': [
    {
      '1': 'group',
      '3': 91525165,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.GroupType',
      '10': 'group'
    },
  ],
};

/// Descriptor for `UpdateGroupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGroupResponseDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVHcm91cFJlc3BvbnNlEi8KBWdyb3VwGK2g0isgASgLMhYuY29nbml0b19pZHAuR3'
    'JvdXBUeXBlUgVncm91cA==');

@$core.Deprecated('Use updateIdentityProviderRequestDescriptor instead')
const UpdateIdentityProviderRequest$json = {
  '1': 'UpdateIdentityProviderRequest',
  '2': [
    {
      '1': 'attributemapping',
      '3': 116923092,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UpdateIdentityProviderRequest.AttributemappingEntry',
      '10': 'attributemapping'
    },
    {
      '1': 'idpidentifiers',
      '3': 205051409,
      '4': 3,
      '5': 9,
      '10': 'idpidentifiers'
    },
    {
      '1': 'providerdetails',
      '3': 476397115,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UpdateIdentityProviderRequest.ProviderdetailsEntry',
      '10': 'providerdetails'
    },
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [
    UpdateIdentityProviderRequest_AttributemappingEntry$json,
    UpdateIdentityProviderRequest_ProviderdetailsEntry$json
  ],
};

@$core.Deprecated('Use updateIdentityProviderRequestDescriptor instead')
const UpdateIdentityProviderRequest_AttributemappingEntry$json = {
  '1': 'AttributemappingEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateIdentityProviderRequestDescriptor instead')
const UpdateIdentityProviderRequest_ProviderdetailsEntry$json = {
  '1': 'ProviderdetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateIdentityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateIdentityProviderRequestDescriptor = $convert.base64Decode(
    'Ch1VcGRhdGVJZGVudGl0eVByb3ZpZGVyUmVxdWVzdBJvChBhdHRyaWJ1dGVtYXBwaW5nGNS14D'
    'cgAygLMkAuY29nbml0b19pZHAuVXBkYXRlSWRlbnRpdHlQcm92aWRlclJlcXVlc3QuQXR0cmli'
    'dXRlbWFwcGluZ0VudHJ5UhBhdHRyaWJ1dGVtYXBwaW5nEikKDmlkcGlkZW50aWZpZXJzGJGs42'
    'EgAygJUg5pZHBpZGVudGlmaWVycxJtCg9wcm92aWRlcmRldGFpbHMYu/yU4wEgAygLMj8uY29n'
    'bml0b19pZHAuVXBkYXRlSWRlbnRpdHlQcm92aWRlclJlcXVlc3QuUHJvdmlkZXJkZXRhaWxzRW'
    '50cnlSD3Byb3ZpZGVyZGV0YWlscxImCgxwcm92aWRlcm5hbWUY+KGo5wEgASgJUgxwcm92aWRl'
    'cm5hbWUSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQaQwoVQXR0cmlidXRlbW'
    'FwcGluZ0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEa'
    'QgoUUHJvdmlkZXJkZXRhaWxzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKA'
    'lSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use updateIdentityProviderResponseDescriptor instead')
const UpdateIdentityProviderResponse$json = {
  '1': 'UpdateIdentityProviderResponse',
  '2': [
    {
      '1': 'identityprovider',
      '3': 450450187,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.IdentityProviderType',
      '10': 'identityprovider'
    },
  ],
};

/// Descriptor for `UpdateIdentityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateIdentityProviderResponseDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVJZGVudGl0eVByb3ZpZGVyUmVzcG9uc2USUQoQaWRlbnRpdHlwcm92aWRlchiLpu'
        'XWASABKAsyIS5jb2duaXRvX2lkcC5JZGVudGl0eVByb3ZpZGVyVHlwZVIQaWRlbnRpdHlwcm92'
        'aWRlcg==');

@$core.Deprecated('Use updateManagedLoginBrandingRequestDescriptor instead')
const UpdateManagedLoginBrandingRequest$json = {
  '1': 'UpdateManagedLoginBrandingRequest',
  '2': [
    {
      '1': 'assets',
      '3': 141150041,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AssetType',
      '10': 'assets'
    },
    {
      '1': 'managedloginbrandingid',
      '3': 57829482,
      '4': 1,
      '5': 9,
      '10': 'managedloginbrandingid'
    },
    {'1': 'settings', '3': 184911657, '4': 1, '5': 9, '10': 'settings'},
    {
      '1': 'usecognitoprovidedvalues',
      '3': 408662835,
      '4': 1,
      '5': 8,
      '10': 'usecognitoprovidedvalues'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `UpdateManagedLoginBrandingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateManagedLoginBrandingRequestDescriptor = $convert.base64Decode(
    'CiFVcGRhdGVNYW5hZ2VkTG9naW5CcmFuZGluZ1JlcXVlc3QSMQoGYXNzZXRzGNmOp0MgAygLMh'
    'YuY29nbml0b19pZHAuQXNzZXRUeXBlUgZhc3NldHMSOQoWbWFuYWdlZGxvZ2luYnJhbmRpbmdp'
    'ZBjq0MkbIAEoCVIWbWFuYWdlZGxvZ2luYnJhbmRpbmdpZBIdCghzZXR0aW5ncxipjpZYIAEoCV'
    'IIc2V0dGluZ3MSPgoYdXNlY29nbml0b3Byb3ZpZGVkdmFsdWVzGLPm7sIBIAEoCFIYdXNlY29n'
    'bml0b3Byb3ZpZGVkdmFsdWVzEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use updateManagedLoginBrandingResponseDescriptor instead')
const UpdateManagedLoginBrandingResponse$json = {
  '1': 'UpdateManagedLoginBrandingResponse',
  '2': [
    {
      '1': 'managedloginbranding',
      '3': 338791109,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ManagedLoginBrandingType',
      '10': 'managedloginbranding'
    },
  ],
};

/// Descriptor for `UpdateManagedLoginBrandingResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateManagedLoginBrandingResponseDescriptor =
    $convert.base64Decode(
        'CiJVcGRhdGVNYW5hZ2VkTG9naW5CcmFuZGluZ1Jlc3BvbnNlEl0KFG1hbmFnZWRsb2dpbmJyYW'
        '5kaW5nGMWVxqEBIAEoCzIlLmNvZ25pdG9faWRwLk1hbmFnZWRMb2dpbkJyYW5kaW5nVHlwZVIU'
        'bWFuYWdlZGxvZ2luYnJhbmRpbmc=');

@$core.Deprecated('Use updateResourceServerRequestDescriptor instead')
const UpdateResourceServerRequest$json = {
  '1': 'UpdateResourceServerRequest',
  '2': [
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'scopes',
      '3': 464684393,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.ResourceServerScopeType',
      '10': 'scopes'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `UpdateResourceServerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResourceServerRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVSZXNvdXJjZVNlcnZlclJlcXVlc3QSIQoKaWRlbnRpZmllchjfoPsTIAEoCVIKaW'
    'RlbnRpZmllchIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkAKBnNjb3BlcxjpisrdASADKAsyJC5j'
    'b2duaXRvX2lkcC5SZXNvdXJjZVNlcnZlclNjb3BlVHlwZVIGc2NvcGVzEiIKCnVzZXJwb29saW'
    'QY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use updateResourceServerResponseDescriptor instead')
const UpdateResourceServerResponse$json = {
  '1': 'UpdateResourceServerResponse',
  '2': [
    {
      '1': 'resourceserver',
      '3': 506282051,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.ResourceServerType',
      '10': 'resourceserver'
    },
  ],
};

/// Descriptor for `UpdateResourceServerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResourceServerResponseDescriptor =
    $convert.base64Decode(
        'ChxVcGRhdGVSZXNvdXJjZVNlcnZlclJlc3BvbnNlEksKDnJlc291cmNlc2VydmVyGMOAtfEBIA'
        'EoCzIfLmNvZ25pdG9faWRwLlJlc291cmNlU2VydmVyVHlwZVIOcmVzb3VyY2VzZXJ2ZXI=');

@$core.Deprecated('Use updateTermsRequestDescriptor instead')
const UpdateTermsRequest$json = {
  '1': 'UpdateTermsRequest',
  '2': [
    {
      '1': 'enforcement',
      '3': 412213242,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TermsEnforcementType',
      '10': 'enforcement'
    },
    {
      '1': 'links',
      '3': 302123151,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UpdateTermsRequest.LinksEntry',
      '10': 'links'
    },
    {'1': 'termsid', '3': 331306210, '4': 1, '5': 9, '10': 'termsid'},
    {'1': 'termsname', '3': 303051560, '4': 1, '5': 9, '10': 'termsname'},
    {
      '1': 'termssource',
      '3': 122689594,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.TermsSourceType',
      '10': 'termssource'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
  '3': [UpdateTermsRequest_LinksEntry$json],
};

@$core.Deprecated('Use updateTermsRequestDescriptor instead')
const UpdateTermsRequest_LinksEntry$json = {
  '1': 'LinksEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateTermsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTermsRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVUZXJtc1JlcXVlc3QSRwoLZW5mb3JjZW1lbnQY+r/HxAEgASgOMiEuY29nbml0b1'
    '9pZHAuVGVybXNFbmZvcmNlbWVudFR5cGVSC2VuZm9yY2VtZW50EkQKBWxpbmtzGI+RiJABIAMo'
    'CzIqLmNvZ25pdG9faWRwLlVwZGF0ZVRlcm1zUmVxdWVzdC5MaW5rc0VudHJ5UgVsaW5rcxIcCg'
    'd0ZXJtc2lkGOKp/Z0BIAEoCVIHdGVybXNpZBIgCgl0ZXJtc25hbWUYqObAkAEgASgJUgl0ZXJt'
    'c25hbWUSQQoLdGVybXNzb3VyY2UYurDAOiABKA4yHC5jb2duaXRvX2lkcC5UZXJtc1NvdXJjZV'
    'R5cGVSC3Rlcm1zc291cmNlEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlkGjgK'
    'CkxpbmtzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ'
    '==');

@$core.Deprecated('Use updateTermsResponseDescriptor instead')
const UpdateTermsResponse$json = {
  '1': 'UpdateTermsResponse',
  '2': [
    {
      '1': 'terms',
      '3': 339062221,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.TermsType',
      '10': 'terms'
    },
  ],
};

/// Descriptor for `UpdateTermsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTermsResponseDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVUZXJtc1Jlc3BvbnNlEjAKBXRlcm1zGM3b1qEBIAEoCzIWLmNvZ25pdG9faWRwLl'
    'Rlcm1zVHlwZVIFdGVybXM=');

@$core.Deprecated('Use updateUserAttributesRequestDescriptor instead')
const UpdateUserAttributesRequest$json = {
  '1': 'UpdateUserAttributesRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'clientmetadata',
      '3': 205510604,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UpdateUserAttributesRequest.ClientmetadataEntry',
      '10': 'clientmetadata'
    },
    {
      '1': 'userattributes',
      '3': 194667064,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'userattributes'
    },
  ],
  '3': [UpdateUserAttributesRequest_ClientmetadataEntry$json],
};

@$core.Deprecated('Use updateUserAttributesRequestDescriptor instead')
const UpdateUserAttributesRequest_ClientmetadataEntry$json = {
  '1': 'ClientmetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateUserAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserAttributesRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVVc2VyQXR0cmlidXRlc1JlcXVlc3QSIwoLYWNjZXNzdG9rZW4YibyQRiABKAlSC2'
    'FjY2Vzc3Rva2VuEmcKDmNsaWVudG1ldGFkYXRhGMyv/2EgAygLMjwuY29nbml0b19pZHAuVXBk'
    'YXRlVXNlckF0dHJpYnV0ZXNSZXF1ZXN0LkNsaWVudG1ldGFkYXRhRW50cnlSDmNsaWVudG1ldG'
    'FkYXRhEkUKDnVzZXJhdHRyaWJ1dGVzGLjE6VwgAygLMhouY29nbml0b19pZHAuQXR0cmlidXRl'
    'VHlwZVIOdXNlcmF0dHJpYnV0ZXMaQQoTQ2xpZW50bWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKA'
    'lSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use updateUserAttributesResponseDescriptor instead')
const UpdateUserAttributesResponse$json = {
  '1': 'UpdateUserAttributesResponse',
  '2': [
    {
      '1': 'codedeliverydetailslist',
      '3': 298732863,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.CodeDeliveryDetailsType',
      '10': 'codedeliverydetailslist'
    },
  ],
};

/// Descriptor for `UpdateUserAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserAttributesResponseDescriptor =
    $convert.base64Decode(
        'ChxVcGRhdGVVc2VyQXR0cmlidXRlc1Jlc3BvbnNlEmIKF2NvZGVkZWxpdmVyeWRldGFpbHNsaX'
        'N0GL+auY4BIAMoCzIkLmNvZ25pdG9faWRwLkNvZGVEZWxpdmVyeURldGFpbHNUeXBlUhdjb2Rl'
        'ZGVsaXZlcnlkZXRhaWxzbGlzdA==');

@$core.Deprecated('Use updateUserPoolClientRequestDescriptor instead')
const UpdateUserPoolClientRequest$json = {
  '1': 'UpdateUserPoolClientRequest',
  '2': [
    {
      '1': 'accesstokenvalidity',
      '3': 260874267,
      '4': 1,
      '5': 5,
      '10': 'accesstokenvalidity'
    },
    {
      '1': 'allowedoauthflows',
      '3': 268290584,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.OAuthFlowType',
      '10': 'allowedoauthflows'
    },
    {
      '1': 'allowedoauthflowsuserpoolclient',
      '3': 520095610,
      '4': 1,
      '5': 8,
      '10': 'allowedoauthflowsuserpoolclient'
    },
    {
      '1': 'allowedoauthscopes',
      '3': 39385504,
      '4': 3,
      '5': 9,
      '10': 'allowedoauthscopes'
    },
    {
      '1': 'analyticsconfiguration',
      '3': 229750388,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsConfigurationType',
      '10': 'analyticsconfiguration'
    },
    {
      '1': 'authsessionvalidity',
      '3': 223873468,
      '4': 1,
      '5': 5,
      '10': 'authsessionvalidity'
    },
    {'1': 'callbackurls', '3': 227703885, '4': 3, '5': 9, '10': 'callbackurls'},
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'clientname', '3': 340245630, '4': 1, '5': 9, '10': 'clientname'},
    {
      '1': 'defaultredirecturi',
      '3': 311293253,
      '4': 1,
      '5': 9,
      '10': 'defaultredirecturi'
    },
    {
      '1': 'enablepropagateadditionalusercontextdata',
      '3': 201651031,
      '4': 1,
      '5': 8,
      '10': 'enablepropagateadditionalusercontextdata'
    },
    {
      '1': 'enabletokenrevocation',
      '3': 178186392,
      '4': 1,
      '5': 8,
      '10': 'enabletokenrevocation'
    },
    {
      '1': 'explicitauthflows',
      '3': 277179621,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.ExplicitAuthFlowsType',
      '10': 'explicitauthflows'
    },
    {
      '1': 'idtokenvalidity',
      '3': 312934952,
      '4': 1,
      '5': 5,
      '10': 'idtokenvalidity'
    },
    {'1': 'logouturls', '3': 468187518, '4': 3, '5': 9, '10': 'logouturls'},
    {
      '1': 'preventuserexistenceerrors',
      '3': 188235606,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.PreventUserExistenceErrorTypes',
      '10': 'preventuserexistenceerrors'
    },
    {
      '1': 'readattributes',
      '3': 334413205,
      '4': 3,
      '5': 9,
      '10': 'readattributes'
    },
    {
      '1': 'refreshtokenrotation',
      '3': 199284564,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.RefreshTokenRotationType',
      '10': 'refreshtokenrotation'
    },
    {
      '1': 'refreshtokenvalidity',
      '3': 303433364,
      '4': 1,
      '5': 5,
      '10': 'refreshtokenvalidity'
    },
    {
      '1': 'supportedidentityproviders',
      '3': 439564368,
      '4': 3,
      '5': 9,
      '10': 'supportedidentityproviders'
    },
    {
      '1': 'tokenvalidityunits',
      '3': 2056664,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.TokenValidityUnitsType',
      '10': 'tokenvalidityunits'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {
      '1': 'writeattributes',
      '3': 440236318,
      '4': 3,
      '5': 9,
      '10': 'writeattributes'
    },
  ],
};

/// Descriptor for `UpdateUserPoolClientRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserPoolClientRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVVc2VyUG9vbENsaWVudFJlcXVlc3QSMwoTYWNjZXNzdG9rZW52YWxpZGl0eRibwL'
    'J8IAEoBVITYWNjZXNzdG9rZW52YWxpZGl0eRJLChFhbGxvd2Vkb2F1dGhmbG93cxiYlPd/IAMo'
    'DjIaLmNvZ25pdG9faWRwLk9BdXRoRmxvd1R5cGVSEWFsbG93ZWRvYXV0aGZsb3dzEkwKH2FsbG'
    '93ZWRvYXV0aGZsb3dzdXNlcnBvb2xjbGllbnQY+o6A+AEgASgIUh9hbGxvd2Vkb2F1dGhmbG93'
    'c3VzZXJwb29sY2xpZW50EjEKEmFsbG93ZWRvYXV0aHNjb3Blcxig8+MSIAMoCVISYWxsb3dlZG'
    '9hdXRoc2NvcGVzEmIKFmFuYWx5dGljc2NvbmZpZ3VyYXRpb24Y9OzGbSABKAsyJy5jb2duaXRv'
    'X2lkcC5BbmFseXRpY3NDb25maWd1cmF0aW9uVHlwZVIWYW5hbHl0aWNzY29uZmlndXJhdGlvbh'
    'IzChNhdXRoc2Vzc2lvbnZhbGlkaXR5GLyT4GogASgFUhNhdXRoc2Vzc2lvbnZhbGlkaXR5EiUK'
    'DGNhbGxiYWNrdXJscxjN+MlsIAMoCVIMY2FsbGJhY2t1cmxzEh4KCGNsaWVudGlkGKTohtYBIA'
    'EoCVIIY2xpZW50aWQSIgoKY2xpZW50bmFtZRj++J6iASABKAlSCmNsaWVudG5hbWUSMgoSZGVm'
    'YXVsdHJlZGlyZWN0dXJpGMXqt5QBIAEoCVISZGVmYXVsdHJlZGlyZWN0dXJpEl0KKGVuYWJsZX'
    'Byb3BhZ2F0ZWFkZGl0aW9uYWx1c2VyY29udGV4dGRhdGEY1+aTYCABKAhSKGVuYWJsZXByb3Bh'
    'Z2F0ZWFkZGl0aW9uYWx1c2VyY29udGV4dGRhdGESNwoVZW5hYmxldG9rZW5yZXZvY2F0aW9uGJ'
    'jR+1QgASgIUhVlbmFibGV0b2tlbnJldm9jYXRpb24SVAoRZXhwbGljaXRhdXRoZmxvd3MY5dmV'
    'hAEgAygOMiIuY29nbml0b19pZHAuRXhwbGljaXRBdXRoRmxvd3NUeXBlUhFleHBsaWNpdGF1dG'
    'hmbG93cxIsCg9pZHRva2VudmFsaWRpdHkYqISclQEgASgFUg9pZHRva2VudmFsaWRpdHkSIgoK'
    'bG9nb3V0dXJscxj+8p/fASADKAlSCmxvZ291dHVybHMSbgoacHJldmVudHVzZXJleGlzdGVuY2'
    'VlcnJvcnMY1v7gWSABKA4yKy5jb2duaXRvX2lkcC5QcmV2ZW50VXNlckV4aXN0ZW5jZUVycm9y'
    'VHlwZXNSGnByZXZlbnR1c2VyZXhpc3RlbmNlZXJyb3JzEioKDnJlYWRhdHRyaWJ1dGVzGJX7up'
    '8BIAMoCVIOcmVhZGF0dHJpYnV0ZXMSXAoUcmVmcmVzaHRva2Vucm90YXRpb24Y1K6DXyABKAsy'
    'JS5jb2duaXRvX2lkcC5SZWZyZXNoVG9rZW5Sb3RhdGlvblR5cGVSFHJlZnJlc2h0b2tlbnJvdG'
    'F0aW9uEjYKFHJlZnJlc2h0b2tlbnZhbGlkaXR5GJSN2JABIAEoBVIUcmVmcmVzaHRva2VudmFs'
    'aWRpdHkSQgoac3VwcG9ydGVkaWRlbnRpdHlwcm92aWRlcnMY0PDM0QEgAygJUhpzdXBwb3J0ZW'
    'RpZGVudGl0eXByb3ZpZGVycxJVChJ0b2tlbnZhbGlkaXR5dW5pdHMY2MN9IAEoCzIjLmNvZ25p'
    'dG9faWRwLlRva2VuVmFsaWRpdHlVbml0c1R5cGVSEnRva2VudmFsaWRpdHl1bml0cxIiCgp1c2'
    'VycG9vbGlkGP7Gi50BIAEoCVIKdXNlcnBvb2xpZBIsCg93cml0ZWF0dHJpYnV0ZXMYnvL10QEg'
    'AygJUg93cml0ZWF0dHJpYnV0ZXM=');

@$core.Deprecated('Use updateUserPoolClientResponseDescriptor instead')
const UpdateUserPoolClientResponse$json = {
  '1': 'UpdateUserPoolClientResponse',
  '2': [
    {
      '1': 'userpoolclient',
      '3': 138497904,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolClientType',
      '10': 'userpoolclient'
    },
  ],
};

/// Descriptor for `UpdateUserPoolClientResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserPoolClientResponseDescriptor =
    $convert.base64Decode(
        'ChxVcGRhdGVVc2VyUG9vbENsaWVudFJlc3BvbnNlEkoKDnVzZXJwb29sY2xpZW50GPCehUIgAS'
        'gLMh8uY29nbml0b19pZHAuVXNlclBvb2xDbGllbnRUeXBlUg51c2VycG9vbGNsaWVudA==');

@$core.Deprecated('Use updateUserPoolDomainRequestDescriptor instead')
const UpdateUserPoolDomainRequest$json = {
  '1': 'UpdateUserPoolDomainRequest',
  '2': [
    {
      '1': 'customdomainconfig',
      '3': 472479421,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.CustomDomainConfigType',
      '10': 'customdomainconfig'
    },
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'managedloginversion',
      '3': 479901038,
      '4': 1,
      '5': 5,
      '10': 'managedloginversion'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `UpdateUserPoolDomainRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserPoolDomainRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVVc2VyUG9vbERvbWFpblJlcXVlc3QSVwoSY3VzdG9tZG9tYWluY29uZmlnGL3tpe'
    'EBIAEoCzIjLmNvZ25pdG9faWRwLkN1c3RvbURvbWFpbkNvbmZpZ1R5cGVSEmN1c3RvbWRvbWFp'
    'bmNvbmZpZxIaCgZkb21haW4YkpLy8AEgASgJUgZkb21haW4SNAoTbWFuYWdlZGxvZ2ludmVyc2'
    'lvbhju6urkASABKAVSE21hbmFnZWRsb2dpbnZlcnNpb24SIgoKdXNlcnBvb2xpZBj+xoudASAB'
    'KAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use updateUserPoolDomainResponseDescriptor instead')
const UpdateUserPoolDomainResponse$json = {
  '1': 'UpdateUserPoolDomainResponse',
  '2': [
    {
      '1': 'cloudfrontdomain',
      '3': 436051942,
      '4': 1,
      '5': 9,
      '10': 'cloudfrontdomain'
    },
    {
      '1': 'managedloginversion',
      '3': 479901038,
      '4': 1,
      '5': 5,
      '10': 'managedloginversion'
    },
  ],
};

/// Descriptor for `UpdateUserPoolDomainResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserPoolDomainResponseDescriptor =
    $convert.base64Decode(
        'ChxVcGRhdGVVc2VyUG9vbERvbWFpblJlc3BvbnNlEi4KEGNsb3VkZnJvbnRkb21haW4Y5r/2zw'
        'EgASgJUhBjbG91ZGZyb250ZG9tYWluEjQKE21hbmFnZWRsb2dpbnZlcnNpb24Y7urq5AEgASgF'
        'UhNtYW5hZ2VkbG9naW52ZXJzaW9u');

@$core.Deprecated('Use updateUserPoolRequestDescriptor instead')
const UpdateUserPoolRequest$json = {
  '1': 'UpdateUserPoolRequest',
  '2': [
    {
      '1': 'accountrecoverysetting',
      '3': 219232186,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountRecoverySettingType',
      '10': 'accountrecoverysetting'
    },
    {
      '1': 'admincreateuserconfig',
      '3': 364968418,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AdminCreateUserConfigType',
      '10': 'admincreateuserconfig'
    },
    {
      '1': 'autoverifiedattributes',
      '3': 467729812,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.VerifiedAttributeType',
      '10': 'autoverifiedattributes'
    },
    {
      '1': 'deletionprotection',
      '3': 504781905,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DeletionProtectionType',
      '10': 'deletionprotection'
    },
    {
      '1': 'deviceconfiguration',
      '3': 512944140,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.DeviceConfigurationType',
      '10': 'deviceconfiguration'
    },
    {
      '1': 'emailconfiguration',
      '3': 528317976,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailConfigurationType',
      '10': 'emailconfiguration'
    },
    {
      '1': 'emailverificationmessage',
      '3': 172634664,
      '4': 1,
      '5': 9,
      '10': 'emailverificationmessage'
    },
    {
      '1': 'emailverificationsubject',
      '3': 224298169,
      '4': 1,
      '5': 9,
      '10': 'emailverificationsubject'
    },
    {
      '1': 'lambdaconfig',
      '3': 291837797,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.LambdaConfigType',
      '10': 'lambdaconfig'
    },
    {
      '1': 'mfaconfiguration',
      '3': 259020766,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolMfaType',
      '10': 'mfaconfiguration'
    },
    {
      '1': 'policies',
      '3': 40015384,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolPolicyType',
      '10': 'policies'
    },
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
    {
      '1': 'smsauthenticationmessage',
      '3': 356104990,
      '4': 1,
      '5': 9,
      '10': 'smsauthenticationmessage'
    },
    {
      '1': 'smsconfiguration',
      '3': 10321849,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SmsConfigurationType',
      '10': 'smsconfiguration'
    },
    {
      '1': 'smsverificationmessage',
      '3': 497665917,
      '4': 1,
      '5': 9,
      '10': 'smsverificationmessage'
    },
    {
      '1': 'userattributeupdatesettings',
      '3': 319670235,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserAttributeUpdateSettingsType',
      '10': 'userattributeupdatesettings'
    },
    {
      '1': 'userpooladdons',
      '3': 296941112,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolAddOnsType',
      '10': 'userpooladdons'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {
      '1': 'userpooltags',
      '3': 341705322,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UpdateUserPoolRequest.UserpooltagsEntry',
      '10': 'userpooltags'
    },
    {
      '1': 'userpooltier',
      '3': 80461029,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolTierType',
      '10': 'userpooltier'
    },
    {
      '1': 'verificationmessagetemplate',
      '3': 502836004,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.VerificationMessageTemplateType',
      '10': 'verificationmessagetemplate'
    },
  ],
  '3': [UpdateUserPoolRequest_UserpooltagsEntry$json],
};

@$core.Deprecated('Use updateUserPoolRequestDescriptor instead')
const UpdateUserPoolRequest_UserpooltagsEntry$json = {
  '1': 'UserpooltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateUserPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserPoolRequestDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVVc2VyUG9vbFJlcXVlc3QSYgoWYWNjb3VudHJlY292ZXJ5c2V0dGluZxi678RoIA'
    'EoCzInLmNvZ25pdG9faWRwLkFjY291bnRSZWNvdmVyeVNldHRpbmdUeXBlUhZhY2NvdW50cmVj'
    'b3ZlcnlzZXR0aW5nEmAKFWFkbWluY3JlYXRldXNlcmNvbmZpZxji84OuASABKAsyJi5jb2duaX'
    'RvX2lkcC5BZG1pbkNyZWF0ZVVzZXJDb25maWdUeXBlUhVhZG1pbmNyZWF0ZXVzZXJjb25maWcS'
    'XgoWYXV0b3ZlcmlmaWVkYXR0cmlidXRlcxiU+4PfASADKA4yIi5jb2duaXRvX2lkcC5WZXJpZm'
    'llZEF0dHJpYnV0ZVR5cGVSFmF1dG92ZXJpZmllZGF0dHJpYnV0ZXMSVwoSZGVsZXRpb25wcm90'
    'ZWN0aW9uGNG42fABIAEoDjIjLmNvZ25pdG9faWRwLkRlbGV0aW9uUHJvdGVjdGlvblR5cGVSEm'
    'RlbGV0aW9ucHJvdGVjdGlvbhJaChNkZXZpY2Vjb25maWd1cmF0aW9uGIzQy/QBIAEoCzIkLmNv'
    'Z25pdG9faWRwLkRldmljZUNvbmZpZ3VyYXRpb25UeXBlUhNkZXZpY2Vjb25maWd1cmF0aW9uEl'
    'cKEmVtYWlsY29uZmlndXJhdGlvbhiY/PX7ASABKAsyIy5jb2duaXRvX2lkcC5FbWFpbENvbmZp'
    'Z3VyYXRpb25UeXBlUhJlbWFpbGNvbmZpZ3VyYXRpb24SPQoYZW1haWx2ZXJpZmljYXRpb25tZX'
    'NzYWdlGKjkqFIgASgJUhhlbWFpbHZlcmlmaWNhdGlvbm1lc3NhZ2USPQoYZW1haWx2ZXJpZmlj'
    'YXRpb25zdWJqZWN0GLmJ+mogASgJUhhlbWFpbHZlcmlmaWNhdGlvbnN1YmplY3QSRQoMbGFtYm'
    'RhY29uZmlnGOWulIsBIAEoCzIdLmNvZ25pdG9faWRwLkxhbWJkYUNvbmZpZ1R5cGVSDGxhbWJk'
    'YWNvbmZpZxJLChBtZmFjb25maWd1cmF0aW9uGN6vwXsgASgOMhwuY29nbml0b19pZHAuVXNlcl'
    'Bvb2xNZmFUeXBlUhBtZmFjb25maWd1cmF0aW9uEj4KCHBvbGljaWVzGJisihMgASgLMh8uY29n'
    'bml0b19pZHAuVXNlclBvb2xQb2xpY3lUeXBlUghwb2xpY2llcxIdCghwb29sbmFtZRjJjYUnIA'
    'EoCVIIcG9vbG5hbWUSPgoYc21zYXV0aGVudGljYXRpb25tZXNzYWdlGJ725qkBIAEoCVIYc21z'
    'YXV0aGVudGljYXRpb25tZXNzYWdlElAKEHNtc2NvbmZpZ3VyYXRpb24Yuf/1BCABKAsyIS5jb2'
    'duaXRvX2lkcC5TbXNDb25maWd1cmF0aW9uVHlwZVIQc21zY29uZmlndXJhdGlvbhI6ChZzbXN2'
    'ZXJpZmljYXRpb25tZXNzYWdlGP2Op+0BIAEoCVIWc21zdmVyaWZpY2F0aW9ubWVzc2FnZRJyCh'
    't1c2VyYXR0cmlidXRldXBkYXRlc2V0dGluZ3MY24+3mAEgASgLMiwuY29nbml0b19pZHAuVXNl'
    'ckF0dHJpYnV0ZVVwZGF0ZVNldHRpbmdzVHlwZVIbdXNlcmF0dHJpYnV0ZXVwZGF0ZXNldHRpbm'
    'dzEksKDnVzZXJwb29sYWRkb25zGLjsy40BIAEoCzIfLmNvZ25pdG9faWRwLlVzZXJQb29sQWRk'
    'T25zVHlwZVIOdXNlcnBvb2xhZGRvbnMSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb2'
    '9saWQSXAoMdXNlcnBvb2x0YWdzGOqE+KIBIAMoCzI0LmNvZ25pdG9faWRwLlVwZGF0ZVVzZXJQ'
    'b29sUmVxdWVzdC5Vc2VycG9vbHRhZ3NFbnRyeVIMdXNlcnBvb2x0YWdzEkQKDHVzZXJwb29sdG'
    'llchjl+a4mIAEoDjIdLmNvZ25pdG9faWRwLlVzZXJQb29sVGllclR5cGVSDHVzZXJwb29sdGll'
    'chJyCht2ZXJpZmljYXRpb25tZXNzYWdldGVtcGxhdGUYpNbi7wEgASgLMiwuY29nbml0b19pZH'
    'AuVmVyaWZpY2F0aW9uTWVzc2FnZVRlbXBsYXRlVHlwZVIbdmVyaWZpY2F0aW9ubWVzc2FnZXRl'
    'bXBsYXRlGj8KEVVzZXJwb29sdGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGA'
    'IgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use updateUserPoolResponseDescriptor instead')
const UpdateUserPoolResponse$json = {
  '1': 'UpdateUserPoolResponse',
};

/// Descriptor for `UpdateUserPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateUserPoolResponseDescriptor =
    $convert.base64Decode('ChZVcGRhdGVVc2VyUG9vbFJlc3BvbnNl');

@$core.Deprecated('Use userAttributeUpdateSettingsTypeDescriptor instead')
const UserAttributeUpdateSettingsType$json = {
  '1': 'UserAttributeUpdateSettingsType',
  '2': [
    {
      '1': 'attributesrequireverificationbeforeupdate',
      '3': 71064847,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.VerifiedAttributeType',
      '10': 'attributesrequireverificationbeforeupdate'
    },
  ],
};

/// Descriptor for `UserAttributeUpdateSettingsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userAttributeUpdateSettingsTypeDescriptor =
    $convert.base64Decode(
        'Ch9Vc2VyQXR0cmlidXRlVXBkYXRlU2V0dGluZ3NUeXBlEoMBCilhdHRyaWJ1dGVzcmVxdWlyZX'
        'ZlcmlmaWNhdGlvbmJlZm9yZXVwZGF0ZRiPuvEhIAMoDjIiLmNvZ25pdG9faWRwLlZlcmlmaWVk'
        'QXR0cmlidXRlVHlwZVIpYXR0cmlidXRlc3JlcXVpcmV2ZXJpZmljYXRpb25iZWZvcmV1cGRhdG'
        'U=');

@$core.Deprecated('Use userContextDataTypeDescriptor instead')
const UserContextDataType$json = {
  '1': 'UserContextDataType',
  '2': [
    {'1': 'encodeddata', '3': 153276170, '4': 1, '5': 9, '10': 'encodeddata'},
    {'1': 'ipaddress', '3': 1800397, '4': 1, '5': 9, '10': 'ipaddress'},
  ],
};

/// Descriptor for `UserContextDataType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userContextDataTypeDescriptor = $convert.base64Decode(
    'ChNVc2VyQ29udGV4dERhdGFUeXBlEiMKC2VuY29kZWRkYXRhGIqei0kgASgJUgtlbmNvZGVkZG'
    'F0YRIeCglpcGFkZHJlc3MYzfFtIAEoCVIJaXBhZGRyZXNz');

@$core.Deprecated('Use userImportInProgressExceptionDescriptor instead')
const UserImportInProgressException$json = {
  '1': 'UserImportInProgressException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UserImportInProgressException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userImportInProgressExceptionDescriptor =
    $convert.base64Decode(
        'Ch1Vc2VySW1wb3J0SW5Qcm9ncmVzc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use userImportJobTypeDescriptor instead')
const UserImportJobType$json = {
  '1': 'UserImportJobType',
  '2': [
    {
      '1': 'cloudwatchlogsrolearn',
      '3': 55454690,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsrolearn'
    },
    {
      '1': 'completiondate',
      '3': 130397444,
      '4': 1,
      '5': 9,
      '10': 'completiondate'
    },
    {
      '1': 'completionmessage',
      '3': 254331463,
      '4': 1,
      '5': 9,
      '10': 'completionmessage'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'failedusers', '3': 109599085, '4': 1, '5': 3, '10': 'failedusers'},
    {
      '1': 'importedusers',
      '3': 165762452,
      '4': 1,
      '5': 3,
      '10': 'importedusers'
    },
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {'1': 'jobname', '3': 498531160, '4': 1, '5': 9, '10': 'jobname'},
    {'1': 'presignedurl', '3': 334334652, '4': 1, '5': 9, '10': 'presignedurl'},
    {'1': 'skippedusers', '3': 31320534, '4': 1, '5': 3, '10': 'skippedusers'},
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserImportJobStatusType',
      '10': 'status'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `UserImportJobType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userImportJobTypeDescriptor = $convert.base64Decode(
    'ChFVc2VySW1wb3J0Sm9iVHlwZRI3ChVjbG91ZHdhdGNobG9nc3JvbGVhcm4Y4te4GiABKAlSFW'
    'Nsb3Vkd2F0Y2hsb2dzcm9sZWFybhIpCg5jb21wbGV0aW9uZGF0ZRiE6pY+IAEoCVIOY29tcGxl'
    'dGlvbmRhdGUSLwoRY29tcGxldGlvbm1lc3NhZ2UYx5SjeSABKAlSEWNvbXBsZXRpb25tZXNzYW'
    'dlEiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9uZGF0ZRIjCgtmYWlsZWR1c2Vy'
    'cxjtsqE0IAEoA1ILZmFpbGVkdXNlcnMSJwoNaW1wb3J0ZWR1c2VycxiUq4VPIAEoA1INaW1wb3'
    'J0ZWR1c2VycxIXCgVqb2JpZBjS1N0zIAEoCVIFam9iaWQSHAoHam9ibmFtZRjY9tvtASABKAlS'
    'B2pvYm5hbWUSJgoMcHJlc2lnbmVkdXJsGLyVtp8BIAEoCVIMcHJlc2lnbmVkdXJsEiUKDHNraX'
    'BwZWR1c2VycxjW0/cOIAEoA1IMc2tpcHBlZHVzZXJzEiAKCXN0YXJ0ZGF0ZRj8+KDUASABKAlS'
    'CXN0YXJ0ZGF0ZRI/CgZzdGF0dXMYkOT7AiABKA4yJC5jb2duaXRvX2lkcC5Vc2VySW1wb3J0Sm'
    '9iU3RhdHVzVHlwZVIGc3RhdHVzEiIKCnVzZXJwb29saWQY/saLnQEgASgJUgp1c2VycG9vbGlk');

@$core.Deprecated('Use userLambdaValidationExceptionDescriptor instead')
const UserLambdaValidationException$json = {
  '1': 'UserLambdaValidationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UserLambdaValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userLambdaValidationExceptionDescriptor =
    $convert.base64Decode(
        'Ch1Vc2VyTGFtYmRhVmFsaWRhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use userNotConfirmedExceptionDescriptor instead')
const UserNotConfirmedException$json = {
  '1': 'UserNotConfirmedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UserNotConfirmedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userNotConfirmedExceptionDescriptor =
    $convert.base64Decode(
        'ChlVc2VyTm90Q29uZmlybWVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use userNotFoundExceptionDescriptor instead')
const UserNotFoundException$json = {
  '1': 'UserNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UserNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userNotFoundExceptionDescriptor = $convert.base64Decode(
    'ChVVc2VyTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use userPoolAddOnNotEnabledExceptionDescriptor instead')
const UserPoolAddOnNotEnabledException$json = {
  '1': 'UserPoolAddOnNotEnabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UserPoolAddOnNotEnabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolAddOnNotEnabledExceptionDescriptor =
    $convert.base64Decode(
        'CiBVc2VyUG9vbEFkZE9uTm90RW5hYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use userPoolAddOnsTypeDescriptor instead')
const UserPoolAddOnsType$json = {
  '1': 'UserPoolAddOnsType',
  '2': [
    {
      '1': 'advancedsecurityadditionalflows',
      '3': 92064958,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AdvancedSecurityAdditionalFlowsType',
      '10': 'advancedsecurityadditionalflows'
    },
    {
      '1': 'advancedsecuritymode',
      '3': 160610251,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.AdvancedSecurityModeType',
      '10': 'advancedsecuritymode'
    },
  ],
};

/// Descriptor for `UserPoolAddOnsType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolAddOnsTypeDescriptor = $convert.base64Decode(
    'ChJVc2VyUG9vbEFkZE9uc1R5cGUSfQofYWR2YW5jZWRzZWN1cml0eWFkZGl0aW9uYWxmbG93cx'
    'i+mfMrIAEoCzIwLmNvZ25pdG9faWRwLkFkdmFuY2VkU2VjdXJpdHlBZGRpdGlvbmFsRmxvd3NU'
    'eXBlUh9hZHZhbmNlZHNlY3VyaXR5YWRkaXRpb25hbGZsb3dzElwKFGFkdmFuY2Vkc2VjdXJpdH'
    'ltb2RlGMvvykwgASgOMiUuY29nbml0b19pZHAuQWR2YW5jZWRTZWN1cml0eU1vZGVUeXBlUhRh'
    'ZHZhbmNlZHNlY3VyaXR5bW9kZQ==');

@$core.Deprecated('Use userPoolClientDescriptionDescriptor instead')
const UserPoolClientDescription$json = {
  '1': 'UserPoolClientDescription',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'clientname', '3': 340245630, '4': 1, '5': 9, '10': 'clientname'},
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
  ],
};

/// Descriptor for `UserPoolClientDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolClientDescriptionDescriptor = $convert.base64Decode(
    'ChlVc2VyUG9vbENsaWVudERlc2NyaXB0aW9uEh4KCGNsaWVudGlkGKTohtYBIAEoCVIIY2xpZW'
    '50aWQSIgoKY2xpZW50bmFtZRj++J6iASABKAlSCmNsaWVudG5hbWUSIgoKdXNlcnBvb2xpZBj+'
    'xoudASABKAlSCnVzZXJwb29saWQ=');

@$core.Deprecated('Use userPoolClientTypeDescriptor instead')
const UserPoolClientType$json = {
  '1': 'UserPoolClientType',
  '2': [
    {
      '1': 'accesstokenvalidity',
      '3': 260874267,
      '4': 1,
      '5': 5,
      '10': 'accesstokenvalidity'
    },
    {
      '1': 'allowedoauthflows',
      '3': 268290584,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.OAuthFlowType',
      '10': 'allowedoauthflows'
    },
    {
      '1': 'allowedoauthflowsuserpoolclient',
      '3': 520095610,
      '4': 1,
      '5': 8,
      '10': 'allowedoauthflowsuserpoolclient'
    },
    {
      '1': 'allowedoauthscopes',
      '3': 39385504,
      '4': 3,
      '5': 9,
      '10': 'allowedoauthscopes'
    },
    {
      '1': 'analyticsconfiguration',
      '3': 229750388,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AnalyticsConfigurationType',
      '10': 'analyticsconfiguration'
    },
    {
      '1': 'authsessionvalidity',
      '3': 223873468,
      '4': 1,
      '5': 5,
      '10': 'authsessionvalidity'
    },
    {'1': 'callbackurls', '3': 227703885, '4': 3, '5': 9, '10': 'callbackurls'},
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'clientname', '3': 340245630, '4': 1, '5': 9, '10': 'clientname'},
    {'1': 'clientsecret', '3': 500734711, '4': 1, '5': 9, '10': 'clientsecret'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'defaultredirecturi',
      '3': 311293253,
      '4': 1,
      '5': 9,
      '10': 'defaultredirecturi'
    },
    {
      '1': 'enablepropagateadditionalusercontextdata',
      '3': 201651031,
      '4': 1,
      '5': 8,
      '10': 'enablepropagateadditionalusercontextdata'
    },
    {
      '1': 'enabletokenrevocation',
      '3': 178186392,
      '4': 1,
      '5': 8,
      '10': 'enabletokenrevocation'
    },
    {
      '1': 'explicitauthflows',
      '3': 277179621,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.ExplicitAuthFlowsType',
      '10': 'explicitauthflows'
    },
    {
      '1': 'idtokenvalidity',
      '3': 312934952,
      '4': 1,
      '5': 5,
      '10': 'idtokenvalidity'
    },
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'logouturls', '3': 468187518, '4': 3, '5': 9, '10': 'logouturls'},
    {
      '1': 'preventuserexistenceerrors',
      '3': 188235606,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.PreventUserExistenceErrorTypes',
      '10': 'preventuserexistenceerrors'
    },
    {
      '1': 'readattributes',
      '3': 334413205,
      '4': 3,
      '5': 9,
      '10': 'readattributes'
    },
    {
      '1': 'refreshtokenrotation',
      '3': 199284564,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.RefreshTokenRotationType',
      '10': 'refreshtokenrotation'
    },
    {
      '1': 'refreshtokenvalidity',
      '3': 303433364,
      '4': 1,
      '5': 5,
      '10': 'refreshtokenvalidity'
    },
    {
      '1': 'supportedidentityproviders',
      '3': 439564368,
      '4': 3,
      '5': 9,
      '10': 'supportedidentityproviders'
    },
    {
      '1': 'tokenvalidityunits',
      '3': 2056664,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.TokenValidityUnitsType',
      '10': 'tokenvalidityunits'
    },
    {'1': 'userpoolid', '3': 329442174, '4': 1, '5': 9, '10': 'userpoolid'},
    {
      '1': 'writeattributes',
      '3': 440236318,
      '4': 3,
      '5': 9,
      '10': 'writeattributes'
    },
  ],
};

/// Descriptor for `UserPoolClientType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolClientTypeDescriptor = $convert.base64Decode(
    'ChJVc2VyUG9vbENsaWVudFR5cGUSMwoTYWNjZXNzdG9rZW52YWxpZGl0eRibwLJ8IAEoBVITYW'
    'NjZXNzdG9rZW52YWxpZGl0eRJLChFhbGxvd2Vkb2F1dGhmbG93cxiYlPd/IAMoDjIaLmNvZ25p'
    'dG9faWRwLk9BdXRoRmxvd1R5cGVSEWFsbG93ZWRvYXV0aGZsb3dzEkwKH2FsbG93ZWRvYXV0aG'
    'Zsb3dzdXNlcnBvb2xjbGllbnQY+o6A+AEgASgIUh9hbGxvd2Vkb2F1dGhmbG93c3VzZXJwb29s'
    'Y2xpZW50EjEKEmFsbG93ZWRvYXV0aHNjb3Blcxig8+MSIAMoCVISYWxsb3dlZG9hdXRoc2NvcG'
    'VzEmIKFmFuYWx5dGljc2NvbmZpZ3VyYXRpb24Y9OzGbSABKAsyJy5jb2duaXRvX2lkcC5BbmFs'
    'eXRpY3NDb25maWd1cmF0aW9uVHlwZVIWYW5hbHl0aWNzY29uZmlndXJhdGlvbhIzChNhdXRoc2'
    'Vzc2lvbnZhbGlkaXR5GLyT4GogASgFUhNhdXRoc2Vzc2lvbnZhbGlkaXR5EiUKDGNhbGxiYWNr'
    'dXJscxjN+MlsIAMoCVIMY2FsbGJhY2t1cmxzEh4KCGNsaWVudGlkGKTohtYBIAEoCVIIY2xpZW'
    '50aWQSIgoKY2xpZW50bmFtZRj++J6iASABKAlSCmNsaWVudG5hbWUSJgoMY2xpZW50c2VjcmV0'
    'GPe14u4BIAEoCVIMY2xpZW50c2VjcmV0EiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZW'
    'F0aW9uZGF0ZRIyChJkZWZhdWx0cmVkaXJlY3R1cmkYxeq3lAEgASgJUhJkZWZhdWx0cmVkaXJl'
    'Y3R1cmkSXQooZW5hYmxlcHJvcGFnYXRlYWRkaXRpb25hbHVzZXJjb250ZXh0ZGF0YRjX5pNgIA'
    'EoCFIoZW5hYmxlcHJvcGFnYXRlYWRkaXRpb25hbHVzZXJjb250ZXh0ZGF0YRI3ChVlbmFibGV0'
    'b2tlbnJldm9jYXRpb24YmNH7VCABKAhSFWVuYWJsZXRva2VucmV2b2NhdGlvbhJUChFleHBsaW'
    'NpdGF1dGhmbG93cxjl2ZWEASADKA4yIi5jb2duaXRvX2lkcC5FeHBsaWNpdEF1dGhGbG93c1R5'
    'cGVSEWV4cGxpY2l0YXV0aGZsb3dzEiwKD2lkdG9rZW52YWxpZGl0eRiohJyVASABKAVSD2lkdG'
    '9rZW52YWxpZGl0eRItChBsYXN0bW9kaWZpZWRkYXRlGNOIyAsgASgJUhBsYXN0bW9kaWZpZWRk'
    'YXRlEiIKCmxvZ291dHVybHMY/vKf3wEgAygJUgpsb2dvdXR1cmxzEm4KGnByZXZlbnR1c2VyZX'
    'hpc3RlbmNlZXJyb3JzGNb+4FkgASgOMisuY29nbml0b19pZHAuUHJldmVudFVzZXJFeGlzdGVu'
    'Y2VFcnJvclR5cGVzUhpwcmV2ZW50dXNlcmV4aXN0ZW5jZWVycm9ycxIqCg5yZWFkYXR0cmlidX'
    'RlcxiV+7qfASADKAlSDnJlYWRhdHRyaWJ1dGVzElwKFHJlZnJlc2h0b2tlbnJvdGF0aW9uGNSu'
    'g18gASgLMiUuY29nbml0b19pZHAuUmVmcmVzaFRva2VuUm90YXRpb25UeXBlUhRyZWZyZXNodG'
    '9rZW5yb3RhdGlvbhI2ChRyZWZyZXNodG9rZW52YWxpZGl0eRiUjdiQASABKAVSFHJlZnJlc2h0'
    'b2tlbnZhbGlkaXR5EkIKGnN1cHBvcnRlZGlkZW50aXR5cHJvdmlkZXJzGNDwzNEBIAMoCVIac3'
    'VwcG9ydGVkaWRlbnRpdHlwcm92aWRlcnMSVQoSdG9rZW52YWxpZGl0eXVuaXRzGNjDfSABKAsy'
    'Iy5jb2duaXRvX2lkcC5Ub2tlblZhbGlkaXR5VW5pdHNUeXBlUhJ0b2tlbnZhbGlkaXR5dW5pdH'
    'MSIgoKdXNlcnBvb2xpZBj+xoudASABKAlSCnVzZXJwb29saWQSLAoPd3JpdGVhdHRyaWJ1dGVz'
    'GJ7y9dEBIAMoCVIPd3JpdGVhdHRyaWJ1dGVz');

@$core.Deprecated('Use userPoolDescriptionTypeDescriptor instead')
const UserPoolDescriptionType$json = {
  '1': 'UserPoolDescriptionType',
  '2': [
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lambdaconfig',
      '3': 291837797,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.LambdaConfigType',
      '10': 'lambdaconfig'
    },
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.StatusType',
      '10': 'status'
    },
  ],
};

/// Descriptor for `UserPoolDescriptionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolDescriptionTypeDescriptor = $convert.base64Decode(
    'ChdVc2VyUG9vbERlc2NyaXB0aW9uVHlwZRImCgxjcmVhdGlvbmRhdGUY4di3iQEgASgJUgxjcm'
    'VhdGlvbmRhdGUSEgoCaWQYgfKitwEgASgJUgJpZBJFCgxsYW1iZGFjb25maWcY5a6UiwEgASgL'
    'Mh0uY29nbml0b19pZHAuTGFtYmRhQ29uZmlnVHlwZVIMbGFtYmRhY29uZmlnEi0KEGxhc3Rtb2'
    'RpZmllZGRhdGUY04jICyABKAlSEGxhc3Rtb2RpZmllZGRhdGUSFQoEbmFtZRiH5oF/IAEoCVIE'
    'bmFtZRIyCgZzdGF0dXMYkOT7AiABKA4yFy5jb2duaXRvX2lkcC5TdGF0dXNUeXBlUgZzdGF0dX'
    'M=');

@$core.Deprecated('Use userPoolPolicyTypeDescriptor instead')
const UserPoolPolicyType$json = {
  '1': 'UserPoolPolicyType',
  '2': [
    {
      '1': 'passwordpolicy',
      '3': 235700227,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.PasswordPolicyType',
      '10': 'passwordpolicy'
    },
    {
      '1': 'signinpolicy',
      '3': 386634174,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SignInPolicyType',
      '10': 'signinpolicy'
    },
  ],
};

/// Descriptor for `UserPoolPolicyType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolPolicyTypeDescriptor = $convert.base64Decode(
    'ChJVc2VyUG9vbFBvbGljeVR5cGUSSgoOcGFzc3dvcmRwb2xpY3kYg4CycCABKAsyHy5jb2duaX'
    'RvX2lkcC5QYXNzd29yZFBvbGljeVR5cGVSDnBhc3N3b3JkcG9saWN5EkUKDHNpZ25pbnBvbGlj'
    'eRi+o664ASABKAsyHS5jb2duaXRvX2lkcC5TaWduSW5Qb2xpY3lUeXBlUgxzaWduaW5wb2xpY3'
    'k=');

@$core.Deprecated('Use userPoolTaggingExceptionDescriptor instead')
const UserPoolTaggingException$json = {
  '1': 'UserPoolTaggingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UserPoolTaggingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolTaggingExceptionDescriptor =
    $convert.base64Decode(
        'ChhVc2VyUG9vbFRhZ2dpbmdFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use userPoolTypeDescriptor instead')
const UserPoolType$json = {
  '1': 'UserPoolType',
  '2': [
    {
      '1': 'accountrecoverysetting',
      '3': 219232186,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AccountRecoverySettingType',
      '10': 'accountrecoverysetting'
    },
    {
      '1': 'admincreateuserconfig',
      '3': 364968418,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.AdminCreateUserConfigType',
      '10': 'admincreateuserconfig'
    },
    {
      '1': 'aliasattributes',
      '3': 189876251,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.AliasAttributeType',
      '10': 'aliasattributes'
    },
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'autoverifiedattributes',
      '3': 467729812,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.VerifiedAttributeType',
      '10': 'autoverifiedattributes'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'customdomain', '3': 459317555, '4': 1, '5': 9, '10': 'customdomain'},
    {
      '1': 'deletionprotection',
      '3': 504781905,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DeletionProtectionType',
      '10': 'deletionprotection'
    },
    {
      '1': 'deviceconfiguration',
      '3': 512944140,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.DeviceConfigurationType',
      '10': 'deviceconfiguration'
    },
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'emailconfiguration',
      '3': 528317976,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.EmailConfigurationType',
      '10': 'emailconfiguration'
    },
    {
      '1': 'emailconfigurationfailure',
      '3': 518209776,
      '4': 1,
      '5': 9,
      '10': 'emailconfigurationfailure'
    },
    {
      '1': 'emailverificationmessage',
      '3': 172634664,
      '4': 1,
      '5': 9,
      '10': 'emailverificationmessage'
    },
    {
      '1': 'emailverificationsubject',
      '3': 224298169,
      '4': 1,
      '5': 9,
      '10': 'emailverificationsubject'
    },
    {
      '1': 'estimatednumberofusers',
      '3': 20814484,
      '4': 1,
      '5': 5,
      '10': 'estimatednumberofusers'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lambdaconfig',
      '3': 291837797,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.LambdaConfigType',
      '10': 'lambdaconfig'
    },
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'mfaconfiguration',
      '3': 259020766,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolMfaType',
      '10': 'mfaconfiguration'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'policies',
      '3': 40015384,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolPolicyType',
      '10': 'policies'
    },
    {
      '1': 'schemaattributes',
      '3': 417986564,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.SchemaAttributeType',
      '10': 'schemaattributes'
    },
    {
      '1': 'smsauthenticationmessage',
      '3': 356104990,
      '4': 1,
      '5': 9,
      '10': 'smsauthenticationmessage'
    },
    {
      '1': 'smsconfiguration',
      '3': 10321849,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.SmsConfigurationType',
      '10': 'smsconfiguration'
    },
    {
      '1': 'smsconfigurationfailure',
      '3': 525270907,
      '4': 1,
      '5': 9,
      '10': 'smsconfigurationfailure'
    },
    {
      '1': 'smsverificationmessage',
      '3': 497665917,
      '4': 1,
      '5': 9,
      '10': 'smsverificationmessage'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.StatusType',
      '10': 'status'
    },
    {
      '1': 'userattributeupdatesettings',
      '3': 319670235,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserAttributeUpdateSettingsType',
      '10': 'userattributeupdatesettings'
    },
    {
      '1': 'userpooladdons',
      '3': 296941112,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UserPoolAddOnsType',
      '10': 'userpooladdons'
    },
    {
      '1': 'userpooltags',
      '3': 341705322,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.UserPoolType.UserpooltagsEntry',
      '10': 'userpooltags'
    },
    {
      '1': 'userpooltier',
      '3': 80461029,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserPoolTierType',
      '10': 'userpooltier'
    },
    {
      '1': 'usernameattributes',
      '3': 196392641,
      '4': 3,
      '5': 14,
      '6': '.cognito_idp.UsernameAttributeType',
      '10': 'usernameattributes'
    },
    {
      '1': 'usernameconfiguration',
      '3': 15447334,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.UsernameConfigurationType',
      '10': 'usernameconfiguration'
    },
    {
      '1': 'verificationmessagetemplate',
      '3': 502836004,
      '4': 1,
      '5': 11,
      '6': '.cognito_idp.VerificationMessageTemplateType',
      '10': 'verificationmessagetemplate'
    },
  ],
  '3': [UserPoolType_UserpooltagsEntry$json],
};

@$core.Deprecated('Use userPoolTypeDescriptor instead')
const UserPoolType_UserpooltagsEntry$json = {
  '1': 'UserpooltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UserPoolType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userPoolTypeDescriptor = $convert.base64Decode(
    'CgxVc2VyUG9vbFR5cGUSYgoWYWNjb3VudHJlY292ZXJ5c2V0dGluZxi678RoIAEoCzInLmNvZ2'
    '5pdG9faWRwLkFjY291bnRSZWNvdmVyeVNldHRpbmdUeXBlUhZhY2NvdW50cmVjb3ZlcnlzZXR0'
    'aW5nEmAKFWFkbWluY3JlYXRldXNlcmNvbmZpZxji84OuASABKAsyJi5jb2duaXRvX2lkcC5BZG'
    '1pbkNyZWF0ZVVzZXJDb25maWdUeXBlUhVhZG1pbmNyZWF0ZXVzZXJjb25maWcSTAoPYWxpYXNh'
    'dHRyaWJ1dGVzGJuQxVogAygOMh8uY29nbml0b19pZHAuQWxpYXNBdHRyaWJ1dGVUeXBlUg9hbG'
    'lhc2F0dHJpYnV0ZXMSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEl4KFmF1dG92ZXJpZmllZGF0dHJp'
    'YnV0ZXMYlPuD3wEgAygOMiIuY29nbml0b19pZHAuVmVyaWZpZWRBdHRyaWJ1dGVUeXBlUhZhdX'
    'RvdmVyaWZpZWRhdHRyaWJ1dGVzEiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9u'
    'ZGF0ZRImCgxjdXN0b21kb21haW4Ys8KC2wEgASgJUgxjdXN0b21kb21haW4SVwoSZGVsZXRpb2'
    '5wcm90ZWN0aW9uGNG42fABIAEoDjIjLmNvZ25pdG9faWRwLkRlbGV0aW9uUHJvdGVjdGlvblR5'
    'cGVSEmRlbGV0aW9ucHJvdGVjdGlvbhJaChNkZXZpY2Vjb25maWd1cmF0aW9uGIzQy/QBIAEoCz'
    'IkLmNvZ25pdG9faWRwLkRldmljZUNvbmZpZ3VyYXRpb25UeXBlUhNkZXZpY2Vjb25maWd1cmF0'
    'aW9uEhoKBmRvbWFpbhiSkvLwASABKAlSBmRvbWFpbhJXChJlbWFpbGNvbmZpZ3VyYXRpb24YmP'
    'z1+wEgASgLMiMuY29nbml0b19pZHAuRW1haWxDb25maWd1cmF0aW9uVHlwZVISZW1haWxjb25m'
    'aWd1cmF0aW9uEkAKGWVtYWlsY29uZmlndXJhdGlvbmZhaWx1cmUY8IGN9wEgASgJUhllbWFpbG'
    'NvbmZpZ3VyYXRpb25mYWlsdXJlEj0KGGVtYWlsdmVyaWZpY2F0aW9ubWVzc2FnZRio5KhSIAEo'
    'CVIYZW1haWx2ZXJpZmljYXRpb25tZXNzYWdlEj0KGGVtYWlsdmVyaWZpY2F0aW9uc3ViamVjdB'
    'i5ifpqIAEoCVIYZW1haWx2ZXJpZmljYXRpb25zdWJqZWN0EjkKFmVzdGltYXRlZG51bWJlcm9m'
    'dXNlcnMYlLX2CSABKAVSFmVzdGltYXRlZG51bWJlcm9mdXNlcnMSEgoCaWQYgfKitwEgASgJUg'
    'JpZBJFCgxsYW1iZGFjb25maWcY5a6UiwEgASgLMh0uY29nbml0b19pZHAuTGFtYmRhQ29uZmln'
    'VHlwZVIMbGFtYmRhY29uZmlnEi0KEGxhc3Rtb2RpZmllZGRhdGUY04jICyABKAlSEGxhc3Rtb2'
    'RpZmllZGRhdGUSSwoQbWZhY29uZmlndXJhdGlvbhjer8F7IAEoDjIcLmNvZ25pdG9faWRwLlVz'
    'ZXJQb29sTWZhVHlwZVIQbWZhY29uZmlndXJhdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEj'
    '4KCHBvbGljaWVzGJisihMgASgLMh8uY29nbml0b19pZHAuVXNlclBvb2xQb2xpY3lUeXBlUghw'
    'b2xpY2llcxJQChBzY2hlbWFhdHRyaWJ1dGVzGITwp8cBIAMoCzIgLmNvZ25pdG9faWRwLlNjaG'
    'VtYUF0dHJpYnV0ZVR5cGVSEHNjaGVtYWF0dHJpYnV0ZXMSPgoYc21zYXV0aGVudGljYXRpb25t'
    'ZXNzYWdlGJ725qkBIAEoCVIYc21zYXV0aGVudGljYXRpb25tZXNzYWdlElAKEHNtc2NvbmZpZ3'
    'VyYXRpb24Yuf/1BCABKAsyIS5jb2duaXRvX2lkcC5TbXNDb25maWd1cmF0aW9uVHlwZVIQc21z'
    'Y29uZmlndXJhdGlvbhI8ChdzbXNjb25maWd1cmF0aW9uZmFpbHVyZRj7/rv6ASABKAlSF3Ntc2'
    'NvbmZpZ3VyYXRpb25mYWlsdXJlEjoKFnNtc3ZlcmlmaWNhdGlvbm1lc3NhZ2UY/Y6n7QEgASgJ'
    'UhZzbXN2ZXJpZmljYXRpb25tZXNzYWdlEjIKBnN0YXR1cxiQ5PsCIAEoDjIXLmNvZ25pdG9faW'
    'RwLlN0YXR1c1R5cGVSBnN0YXR1cxJyCht1c2VyYXR0cmlidXRldXBkYXRlc2V0dGluZ3MY24+3'
    'mAEgASgLMiwuY29nbml0b19pZHAuVXNlckF0dHJpYnV0ZVVwZGF0ZVNldHRpbmdzVHlwZVIbdX'
    'NlcmF0dHJpYnV0ZXVwZGF0ZXNldHRpbmdzEksKDnVzZXJwb29sYWRkb25zGLjsy40BIAEoCzIf'
    'LmNvZ25pdG9faWRwLlVzZXJQb29sQWRkT25zVHlwZVIOdXNlcnBvb2xhZGRvbnMSUwoMdXNlcn'
    'Bvb2x0YWdzGOqE+KIBIAMoCzIrLmNvZ25pdG9faWRwLlVzZXJQb29sVHlwZS5Vc2VycG9vbHRh'
    'Z3NFbnRyeVIMdXNlcnBvb2x0YWdzEkQKDHVzZXJwb29sdGllchjl+a4mIAEoDjIdLmNvZ25pdG'
    '9faWRwLlVzZXJQb29sVGllclR5cGVSDHVzZXJwb29sdGllchJVChJ1c2VybmFtZWF0dHJpYnV0'
    'ZXMYwe3SXSADKA4yIi5jb2duaXRvX2lkcC5Vc2VybmFtZUF0dHJpYnV0ZVR5cGVSEnVzZXJuYW'
    '1lYXR0cmlidXRlcxJfChV1c2VybmFtZWNvbmZpZ3VyYXRpb24YpuquByABKAsyJi5jb2duaXRv'
    'X2lkcC5Vc2VybmFtZUNvbmZpZ3VyYXRpb25UeXBlUhV1c2VybmFtZWNvbmZpZ3VyYXRpb24Scg'
    'obdmVyaWZpY2F0aW9ubWVzc2FnZXRlbXBsYXRlGKTW4u8BIAEoCzIsLmNvZ25pdG9faWRwLlZl'
    'cmlmaWNhdGlvbk1lc3NhZ2VUZW1wbGF0ZVR5cGVSG3ZlcmlmaWNhdGlvbm1lc3NhZ2V0ZW1wbG'
    'F0ZRo/ChFVc2VycG9vbHRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEo'
    'CVIFdmFsdWU6AjgB');

@$core.Deprecated('Use userTypeDescriptor instead')
const UserType$json = {
  '1': 'UserType',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.AttributeType',
      '10': 'attributes'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'mfaoptions',
      '3': 501540826,
      '4': 3,
      '5': 11,
      '6': '.cognito_idp.MFAOptionType',
      '10': 'mfaoptions'
    },
    {
      '1': 'usercreatedate',
      '3': 73013267,
      '4': 1,
      '5': 9,
      '10': 'usercreatedate'
    },
    {
      '1': 'userlastmodifieddate',
      '3': 80916802,
      '4': 1,
      '5': 9,
      '10': 'userlastmodifieddate'
    },
    {
      '1': 'userstatus',
      '3': 189848701,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserStatusType',
      '10': 'userstatus'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `UserType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userTypeDescriptor = $convert.base64Decode(
    'CghVc2VyVHlwZRI9CgphdHRyaWJ1dGVzGLWp+2MgAygLMhouY29nbml0b19pZHAuQXR0cmlidX'
    'RlVHlwZVIKYXR0cmlidXRlcxIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBI+CgptZmFv'
    'cHRpb25zGNrPk+8BIAMoCzIaLmNvZ25pdG9faWRwLk1GQU9wdGlvblR5cGVSCm1mYW9wdGlvbn'
    'MSKQoOdXNlcmNyZWF0ZWRhdGUYk7DoIiABKAlSDnVzZXJjcmVhdGVkYXRlEjUKFHVzZXJsYXN0'
    'bW9kaWZpZWRkYXRlGMLiyiYgASgJUhR1c2VybGFzdG1vZGlmaWVkZGF0ZRI+Cgp1c2Vyc3RhdH'
    'VzGP24w1ogASgOMhsuY29nbml0b19pZHAuVXNlclN0YXR1c1R5cGVSCnVzZXJzdGF0dXMSHgoI'
    'dXNlcm5hbWUY2qmj4AEgASgJUgh1c2VybmFtZQ==');

@$core.Deprecated('Use usernameConfigurationTypeDescriptor instead')
const UsernameConfigurationType$json = {
  '1': 'UsernameConfigurationType',
  '2': [
    {
      '1': 'casesensitive',
      '3': 258546956,
      '4': 1,
      '5': 8,
      '10': 'casesensitive'
    },
  ],
};

/// Descriptor for `UsernameConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List usernameConfigurationTypeDescriptor =
    $convert.base64Decode(
        'ChlVc2VybmFtZUNvbmZpZ3VyYXRpb25UeXBlEicKDWNhc2VzZW5zaXRpdmUYjLqkeyABKAhSDW'
        'Nhc2VzZW5zaXRpdmU=');

@$core.Deprecated('Use usernameExistsExceptionDescriptor instead')
const UsernameExistsException$json = {
  '1': 'UsernameExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UsernameExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List usernameExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChdVc2VybmFtZUV4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use verificationMessageTemplateTypeDescriptor instead')
const VerificationMessageTemplateType$json = {
  '1': 'VerificationMessageTemplateType',
  '2': [
    {
      '1': 'defaultemailoption',
      '3': 2801728,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.DefaultEmailOptionType',
      '10': 'defaultemailoption'
    },
    {'1': 'emailmessage', '3': 247992871, '4': 1, '5': 9, '10': 'emailmessage'},
    {
      '1': 'emailmessagebylink',
      '3': 444701508,
      '4': 1,
      '5': 9,
      '10': 'emailmessagebylink'
    },
    {'1': 'emailsubject', '3': 215040870, '4': 1, '5': 9, '10': 'emailsubject'},
    {
      '1': 'emailsubjectbylink',
      '3': 351344357,
      '4': 1,
      '5': 9,
      '10': 'emailsubjectbylink'
    },
    {'1': 'smsmessage', '3': 359006226, '4': 1, '5': 9, '10': 'smsmessage'},
  ],
};

/// Descriptor for `VerificationMessageTemplateType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verificationMessageTemplateTypeDescriptor = $convert.base64Decode(
    'Ch9WZXJpZmljYXRpb25NZXNzYWdlVGVtcGxhdGVUeXBlElYKEmRlZmF1bHRlbWFpbG9wdGlvbh'
    'jAgKsBIAEoDjIjLmNvZ25pdG9faWRwLkRlZmF1bHRFbWFpbE9wdGlvblR5cGVSEmRlZmF1bHRl'
    'bWFpbG9wdGlvbhIlCgxlbWFpbG1lc3NhZ2UYp6SgdiABKAlSDGVtYWlsbWVzc2FnZRIyChJlbW'
    'FpbG1lc3NhZ2VieWxpbmsYxLaG1AEgASgJUhJlbWFpbG1lc3NhZ2VieWxpbmsSJQoMZW1haWxz'
    'dWJqZWN0GOaGxWYgASgJUgxlbWFpbHN1YmplY3QSMgoSZW1haWxzdWJqZWN0YnlsaW5rGOWtxK'
    'cBIAEoCVISZW1haWxzdWJqZWN0YnlsaW5rEiIKCnNtc21lc3NhZ2UYkoCYqwEgASgJUgpzbXNt'
    'ZXNzYWdl');

@$core.Deprecated('Use verifySoftwareTokenRequestDescriptor instead')
const VerifySoftwareTokenRequest$json = {
  '1': 'VerifySoftwareTokenRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'friendlydevicename',
      '3': 89691332,
      '4': 1,
      '5': 9,
      '10': 'friendlydevicename'
    },
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {'1': 'usercode', '3': 349164576, '4': 1, '5': 9, '10': 'usercode'},
  ],
};

/// Descriptor for `VerifySoftwareTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifySoftwareTokenRequestDescriptor = $convert.base64Decode(
    'ChpWZXJpZnlTb2Z0d2FyZVRva2VuUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIAEoCVILYW'
    'NjZXNzdG9rZW4SMQoSZnJpZW5kbHlkZXZpY2VuYW1lGMSp4iogASgJUhJmcmllbmRseWRldmlj'
    'ZW5hbWUSGwoHc2Vzc2lvbhiYmaMCIAEoCVIHc2Vzc2lvbhIeCgh1c2VyY29kZRigqL+mASABKA'
    'lSCHVzZXJjb2Rl');

@$core.Deprecated('Use verifySoftwareTokenResponseDescriptor instead')
const VerifySoftwareTokenResponse$json = {
  '1': 'VerifySoftwareTokenResponse',
  '2': [
    {'1': 'session', '3': 4770968, '4': 1, '5': 9, '10': 'session'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.VerifySoftwareTokenResponseType',
      '10': 'status'
    },
  ],
};

/// Descriptor for `VerifySoftwareTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifySoftwareTokenResponseDescriptor =
    $convert.base64Decode(
        'ChtWZXJpZnlTb2Z0d2FyZVRva2VuUmVzcG9uc2USGwoHc2Vzc2lvbhiYmaMCIAEoCVIHc2Vzc2'
        'lvbhJHCgZzdGF0dXMYkOT7AiABKA4yLC5jb2duaXRvX2lkcC5WZXJpZnlTb2Z0d2FyZVRva2Vu'
        'UmVzcG9uc2VUeXBlUgZzdGF0dXM=');

@$core.Deprecated('Use verifyUserAttributeRequestDescriptor instead')
const VerifyUserAttributeRequest$json = {
  '1': 'VerifyUserAttributeRequest',
  '2': [
    {'1': 'accesstoken', '3': 147070473, '4': 1, '5': 9, '10': 'accesstoken'},
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {'1': 'code', '3': 425572629, '4': 1, '5': 9, '10': 'code'},
  ],
};

/// Descriptor for `VerifyUserAttributeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyUserAttributeRequestDescriptor =
    $convert.base64Decode(
        'ChpWZXJpZnlVc2VyQXR0cmlidXRlUmVxdWVzdBIjCgthY2Nlc3N0b2tlbhiJvJBGIAEoCVILYW'
        'NjZXNzdG9rZW4SKAoNYXR0cmlidXRlbmFtZRitlZioASABKAlSDWF0dHJpYnV0ZW5hbWUSFgoE'
        'Y29kZRiV8vbKASABKAlSBGNvZGU=');

@$core.Deprecated('Use verifyUserAttributeResponseDescriptor instead')
const VerifyUserAttributeResponse$json = {
  '1': 'VerifyUserAttributeResponse',
};

/// Descriptor for `VerifyUserAttributeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyUserAttributeResponseDescriptor =
    $convert.base64Decode('ChtWZXJpZnlVc2VyQXR0cmlidXRlUmVzcG9uc2U=');

@$core.Deprecated('Use webAuthnChallengeNotFoundExceptionDescriptor instead')
const WebAuthnChallengeNotFoundException$json = {
  '1': 'WebAuthnChallengeNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `WebAuthnChallengeNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnChallengeNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CiJXZWJBdXRobkNoYWxsZW5nZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use webAuthnClientMismatchExceptionDescriptor instead')
const WebAuthnClientMismatchException$json = {
  '1': 'WebAuthnClientMismatchException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `WebAuthnClientMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnClientMismatchExceptionDescriptor =
    $convert.base64Decode(
        'Ch9XZWJBdXRobkNsaWVudE1pc21hdGNoRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use webAuthnConfigurationMissingExceptionDescriptor instead')
const WebAuthnConfigurationMissingException$json = {
  '1': 'WebAuthnConfigurationMissingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `WebAuthnConfigurationMissingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnConfigurationMissingExceptionDescriptor =
    $convert.base64Decode(
        'CiVXZWJBdXRobkNvbmZpZ3VyYXRpb25NaXNzaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use webAuthnConfigurationTypeDescriptor instead')
const WebAuthnConfigurationType$json = {
  '1': 'WebAuthnConfigurationType',
  '2': [
    {
      '1': 'relyingpartyid',
      '3': 211618527,
      '4': 1,
      '5': 9,
      '10': 'relyingpartyid'
    },
    {
      '1': 'userverification',
      '3': 263409814,
      '4': 1,
      '5': 14,
      '6': '.cognito_idp.UserVerificationType',
      '10': 'userverification'
    },
  ],
};

/// Descriptor for `WebAuthnConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnConfigurationTypeDescriptor = $convert.base64Decode(
    'ChlXZWJBdXRobkNvbmZpZ3VyYXRpb25UeXBlEikKDnJlbHlpbmdwYXJ0eWlkGN+V9GQgASgJUg'
    '5yZWx5aW5ncGFydHlpZBJQChB1c2VydmVyaWZpY2F0aW9uGJahzX0gASgOMiEuY29nbml0b19p'
    'ZHAuVXNlclZlcmlmaWNhdGlvblR5cGVSEHVzZXJ2ZXJpZmljYXRpb24=');

@$core.Deprecated('Use webAuthnCredentialDescriptionDescriptor instead')
const WebAuthnCredentialDescription$json = {
  '1': 'WebAuthnCredentialDescription',
  '2': [
    {
      '1': 'authenticatorattachment',
      '3': 61147702,
      '4': 1,
      '5': 9,
      '10': 'authenticatorattachment'
    },
    {
      '1': 'authenticatortransports',
      '3': 152292059,
      '4': 3,
      '5': 9,
      '10': 'authenticatortransports'
    },
    {'1': 'createdat', '3': 258192751, '4': 1, '5': 9, '10': 'createdat'},
    {'1': 'credentialid', '3': 244832166, '4': 1, '5': 9, '10': 'credentialid'},
    {
      '1': 'friendlycredentialname',
      '3': 27218241,
      '4': 1,
      '5': 9,
      '10': 'friendlycredentialname'
    },
    {
      '1': 'relyingpartyid',
      '3': 211618527,
      '4': 1,
      '5': 9,
      '10': 'relyingpartyid'
    },
  ],
};

/// Descriptor for `WebAuthnCredentialDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnCredentialDescriptionDescriptor = $convert.base64Decode(
    'Ch1XZWJBdXRobkNyZWRlbnRpYWxEZXNjcmlwdGlvbhI7ChdhdXRoZW50aWNhdG9yYXR0YWNobW'
    'VudBi2lJQdIAEoCVIXYXV0aGVudGljYXRvcmF0dGFjaG1lbnQSOwoXYXV0aGVudGljYXRvcnRy'
    'YW5zcG9ydHMY25XPSCADKAlSF2F1dGhlbnRpY2F0b3J0cmFuc3BvcnRzEh8KCWNyZWF0ZWRhdB'
    'jv6o57IAEoCVIJY3JlYXRlZGF0EiUKDGNyZWRlbnRpYWxpZBimr990IAEoCVIMY3JlZGVudGlh'
    'bGlkEjkKFmZyaWVuZGx5Y3JlZGVudGlhbG5hbWUYwaL9DCABKAlSFmZyaWVuZGx5Y3JlZGVudG'
    'lhbG5hbWUSKQoOcmVseWluZ3BhcnR5aWQY35X0ZCABKAlSDnJlbHlpbmdwYXJ0eWlk');

@$core
    .Deprecated('Use webAuthnCredentialNotSupportedExceptionDescriptor instead')
const WebAuthnCredentialNotSupportedException$json = {
  '1': 'WebAuthnCredentialNotSupportedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `WebAuthnCredentialNotSupportedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnCredentialNotSupportedExceptionDescriptor =
    $convert.base64Decode(
        'CidXZWJBdXRobkNyZWRlbnRpYWxOb3RTdXBwb3J0ZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkc'
        'gnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use webAuthnNotEnabledExceptionDescriptor instead')
const WebAuthnNotEnabledException$json = {
  '1': 'WebAuthnNotEnabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `WebAuthnNotEnabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnNotEnabledExceptionDescriptor =
    $convert.base64Decode(
        'ChtXZWJBdXRobk5vdEVuYWJsZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use webAuthnOriginNotAllowedExceptionDescriptor instead')
const WebAuthnOriginNotAllowedException$json = {
  '1': 'WebAuthnOriginNotAllowedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `WebAuthnOriginNotAllowedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnOriginNotAllowedExceptionDescriptor =
    $convert.base64Decode(
        'CiFXZWJBdXRobk9yaWdpbk5vdEFsbG93ZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use webAuthnRelyingPartyMismatchExceptionDescriptor instead')
const WebAuthnRelyingPartyMismatchException$json = {
  '1': 'WebAuthnRelyingPartyMismatchException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `WebAuthnRelyingPartyMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List webAuthnRelyingPartyMismatchExceptionDescriptor =
    $convert.base64Decode(
        'CiVXZWJBdXRoblJlbHlpbmdQYXJ0eU1pc21hdGNoRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');
