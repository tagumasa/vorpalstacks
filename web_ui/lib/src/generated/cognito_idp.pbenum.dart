// This is a generated file - do not edit.
//
// Generated from cognito_idp.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AccountTakeoverEventActionType extends $pb.ProtobufEnum {
  static const AccountTakeoverEventActionType
      ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_BLOCK =
      AccountTakeoverEventActionType._(
          0, _omitEnumNames ? '' : 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_BLOCK');
  static const AccountTakeoverEventActionType
      ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_REQUIRED =
      AccountTakeoverEventActionType._(
          1,
          _omitEnumNames
              ? ''
              : 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_REQUIRED');
  static const AccountTakeoverEventActionType
      ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_NO_ACTION =
      AccountTakeoverEventActionType._(2,
          _omitEnumNames ? '' : 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_NO_ACTION');
  static const AccountTakeoverEventActionType
      ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_IF_CONFIGURED =
      AccountTakeoverEventActionType._(
          3,
          _omitEnumNames
              ? ''
              : 'ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_IF_CONFIGURED');

  static const $core.List<AccountTakeoverEventActionType> values =
      <AccountTakeoverEventActionType>[
    ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_BLOCK,
    ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_REQUIRED,
    ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_NO_ACTION,
    ACCOUNT_TAKEOVER_EVENT_ACTION_TYPE_MFA_IF_CONFIGURED,
  ];

  static final $core.List<AccountTakeoverEventActionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static AccountTakeoverEventActionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AccountTakeoverEventActionType._(super.value, super.name);
}

class AdvancedSecurityEnabledModeType extends $pb.ProtobufEnum {
  static const AdvancedSecurityEnabledModeType
      ADVANCED_SECURITY_ENABLED_MODE_TYPE_AUDIT =
      AdvancedSecurityEnabledModeType._(
          0, _omitEnumNames ? '' : 'ADVANCED_SECURITY_ENABLED_MODE_TYPE_AUDIT');
  static const AdvancedSecurityEnabledModeType
      ADVANCED_SECURITY_ENABLED_MODE_TYPE_ENFORCED =
      AdvancedSecurityEnabledModeType._(1,
          _omitEnumNames ? '' : 'ADVANCED_SECURITY_ENABLED_MODE_TYPE_ENFORCED');

  static const $core.List<AdvancedSecurityEnabledModeType> values =
      <AdvancedSecurityEnabledModeType>[
    ADVANCED_SECURITY_ENABLED_MODE_TYPE_AUDIT,
    ADVANCED_SECURITY_ENABLED_MODE_TYPE_ENFORCED,
  ];

  static final $core.List<AdvancedSecurityEnabledModeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AdvancedSecurityEnabledModeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AdvancedSecurityEnabledModeType._(super.value, super.name);
}

class AdvancedSecurityModeType extends $pb.ProtobufEnum {
  static const AdvancedSecurityModeType ADVANCED_SECURITY_MODE_TYPE_AUDIT =
      AdvancedSecurityModeType._(
          0, _omitEnumNames ? '' : 'ADVANCED_SECURITY_MODE_TYPE_AUDIT');
  static const AdvancedSecurityModeType ADVANCED_SECURITY_MODE_TYPE_ENFORCED =
      AdvancedSecurityModeType._(
          1, _omitEnumNames ? '' : 'ADVANCED_SECURITY_MODE_TYPE_ENFORCED');
  static const AdvancedSecurityModeType ADVANCED_SECURITY_MODE_TYPE_OFF =
      AdvancedSecurityModeType._(
          2, _omitEnumNames ? '' : 'ADVANCED_SECURITY_MODE_TYPE_OFF');

  static const $core.List<AdvancedSecurityModeType> values =
      <AdvancedSecurityModeType>[
    ADVANCED_SECURITY_MODE_TYPE_AUDIT,
    ADVANCED_SECURITY_MODE_TYPE_ENFORCED,
    ADVANCED_SECURITY_MODE_TYPE_OFF,
  ];

  static final $core.List<AdvancedSecurityModeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AdvancedSecurityModeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AdvancedSecurityModeType._(super.value, super.name);
}

class AliasAttributeType extends $pb.ProtobufEnum {
  static const AliasAttributeType ALIAS_ATTRIBUTE_TYPE_PHONE_NUMBER =
      AliasAttributeType._(
          0, _omitEnumNames ? '' : 'ALIAS_ATTRIBUTE_TYPE_PHONE_NUMBER');
  static const AliasAttributeType ALIAS_ATTRIBUTE_TYPE_EMAIL =
      AliasAttributeType._(
          1, _omitEnumNames ? '' : 'ALIAS_ATTRIBUTE_TYPE_EMAIL');
  static const AliasAttributeType ALIAS_ATTRIBUTE_TYPE_PREFERRED_USERNAME =
      AliasAttributeType._(
          2, _omitEnumNames ? '' : 'ALIAS_ATTRIBUTE_TYPE_PREFERRED_USERNAME');

  static const $core.List<AliasAttributeType> values = <AliasAttributeType>[
    ALIAS_ATTRIBUTE_TYPE_PHONE_NUMBER,
    ALIAS_ATTRIBUTE_TYPE_EMAIL,
    ALIAS_ATTRIBUTE_TYPE_PREFERRED_USERNAME,
  ];

  static final $core.List<AliasAttributeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AliasAttributeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AliasAttributeType._(super.value, super.name);
}

class AssetCategoryType extends $pb.ProtobufEnum {
  static const AssetCategoryType ASSET_CATEGORY_TYPE_FORM_LOGO =
      AssetCategoryType._(
          0, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_FORM_LOGO');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_FAVICON_ICO =
      AssetCategoryType._(
          1, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_FAVICON_ICO');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_PASSWORD_GRAPHIC =
      AssetCategoryType._(
          2, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_PASSWORD_GRAPHIC');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_PAGE_BACKGROUND =
      AssetCategoryType._(
          3, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_PAGE_BACKGROUND');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_PAGE_HEADER_BACKGROUND =
      AssetCategoryType._(4,
          _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_PAGE_HEADER_BACKGROUND');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_SMS_GRAPHIC =
      AssetCategoryType._(
          5, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_SMS_GRAPHIC');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_PASSKEY_GRAPHIC =
      AssetCategoryType._(
          6, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_PASSKEY_GRAPHIC');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_PAGE_FOOTER_BACKGROUND =
      AssetCategoryType._(7,
          _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_PAGE_FOOTER_BACKGROUND');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_EMAIL_GRAPHIC =
      AssetCategoryType._(
          8, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_EMAIL_GRAPHIC');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_PAGE_HEADER_LOGO =
      AssetCategoryType._(
          9, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_PAGE_HEADER_LOGO');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_FORM_BACKGROUND =
      AssetCategoryType._(
          10, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_FORM_BACKGROUND');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_AUTH_APP_GRAPHIC =
      AssetCategoryType._(
          11, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_AUTH_APP_GRAPHIC');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_IDP_BUTTON_ICON =
      AssetCategoryType._(
          12, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_IDP_BUTTON_ICON');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_FAVICON_SVG =
      AssetCategoryType._(
          13, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_FAVICON_SVG');
  static const AssetCategoryType ASSET_CATEGORY_TYPE_PAGE_FOOTER_LOGO =
      AssetCategoryType._(
          14, _omitEnumNames ? '' : 'ASSET_CATEGORY_TYPE_PAGE_FOOTER_LOGO');

  static const $core.List<AssetCategoryType> values = <AssetCategoryType>[
    ASSET_CATEGORY_TYPE_FORM_LOGO,
    ASSET_CATEGORY_TYPE_FAVICON_ICO,
    ASSET_CATEGORY_TYPE_PASSWORD_GRAPHIC,
    ASSET_CATEGORY_TYPE_PAGE_BACKGROUND,
    ASSET_CATEGORY_TYPE_PAGE_HEADER_BACKGROUND,
    ASSET_CATEGORY_TYPE_SMS_GRAPHIC,
    ASSET_CATEGORY_TYPE_PASSKEY_GRAPHIC,
    ASSET_CATEGORY_TYPE_PAGE_FOOTER_BACKGROUND,
    ASSET_CATEGORY_TYPE_EMAIL_GRAPHIC,
    ASSET_CATEGORY_TYPE_PAGE_HEADER_LOGO,
    ASSET_CATEGORY_TYPE_FORM_BACKGROUND,
    ASSET_CATEGORY_TYPE_AUTH_APP_GRAPHIC,
    ASSET_CATEGORY_TYPE_IDP_BUTTON_ICON,
    ASSET_CATEGORY_TYPE_FAVICON_SVG,
    ASSET_CATEGORY_TYPE_PAGE_FOOTER_LOGO,
  ];

  static final $core.List<AssetCategoryType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 14);
  static AssetCategoryType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssetCategoryType._(super.value, super.name);
}

class AssetExtensionType extends $pb.ProtobufEnum {
  static const AssetExtensionType ASSET_EXTENSION_TYPE_ICO =
      AssetExtensionType._(0, _omitEnumNames ? '' : 'ASSET_EXTENSION_TYPE_ICO');
  static const AssetExtensionType ASSET_EXTENSION_TYPE_SVG =
      AssetExtensionType._(1, _omitEnumNames ? '' : 'ASSET_EXTENSION_TYPE_SVG');
  static const AssetExtensionType ASSET_EXTENSION_TYPE_JPEG =
      AssetExtensionType._(
          2, _omitEnumNames ? '' : 'ASSET_EXTENSION_TYPE_JPEG');
  static const AssetExtensionType ASSET_EXTENSION_TYPE_PNG =
      AssetExtensionType._(3, _omitEnumNames ? '' : 'ASSET_EXTENSION_TYPE_PNG');
  static const AssetExtensionType ASSET_EXTENSION_TYPE_WEBP =
      AssetExtensionType._(
          4, _omitEnumNames ? '' : 'ASSET_EXTENSION_TYPE_WEBP');

  static const $core.List<AssetExtensionType> values = <AssetExtensionType>[
    ASSET_EXTENSION_TYPE_ICO,
    ASSET_EXTENSION_TYPE_SVG,
    ASSET_EXTENSION_TYPE_JPEG,
    ASSET_EXTENSION_TYPE_PNG,
    ASSET_EXTENSION_TYPE_WEBP,
  ];

  static final $core.List<AssetExtensionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static AssetExtensionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssetExtensionType._(super.value, super.name);
}

class AttributeDataType extends $pb.ProtobufEnum {
  static const AttributeDataType ATTRIBUTE_DATA_TYPE_DATETIME =
      AttributeDataType._(
          0, _omitEnumNames ? '' : 'ATTRIBUTE_DATA_TYPE_DATETIME');
  static const AttributeDataType ATTRIBUTE_DATA_TYPE_STRING =
      AttributeDataType._(
          1, _omitEnumNames ? '' : 'ATTRIBUTE_DATA_TYPE_STRING');
  static const AttributeDataType ATTRIBUTE_DATA_TYPE_NUMBER =
      AttributeDataType._(
          2, _omitEnumNames ? '' : 'ATTRIBUTE_DATA_TYPE_NUMBER');
  static const AttributeDataType ATTRIBUTE_DATA_TYPE_BOOLEAN =
      AttributeDataType._(
          3, _omitEnumNames ? '' : 'ATTRIBUTE_DATA_TYPE_BOOLEAN');

  static const $core.List<AttributeDataType> values = <AttributeDataType>[
    ATTRIBUTE_DATA_TYPE_DATETIME,
    ATTRIBUTE_DATA_TYPE_STRING,
    ATTRIBUTE_DATA_TYPE_NUMBER,
    ATTRIBUTE_DATA_TYPE_BOOLEAN,
  ];

  static final $core.List<AttributeDataType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static AttributeDataType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AttributeDataType._(super.value, super.name);
}

class AuthFactorType extends $pb.ProtobufEnum {
  static const AuthFactorType AUTH_FACTOR_TYPE_EMAIL_OTP =
      AuthFactorType._(0, _omitEnumNames ? '' : 'AUTH_FACTOR_TYPE_EMAIL_OTP');
  static const AuthFactorType AUTH_FACTOR_TYPE_SMS_OTP =
      AuthFactorType._(1, _omitEnumNames ? '' : 'AUTH_FACTOR_TYPE_SMS_OTP');
  static const AuthFactorType AUTH_FACTOR_TYPE_WEB_AUTHN =
      AuthFactorType._(2, _omitEnumNames ? '' : 'AUTH_FACTOR_TYPE_WEB_AUTHN');
  static const AuthFactorType AUTH_FACTOR_TYPE_PASSWORD =
      AuthFactorType._(3, _omitEnumNames ? '' : 'AUTH_FACTOR_TYPE_PASSWORD');

  static const $core.List<AuthFactorType> values = <AuthFactorType>[
    AUTH_FACTOR_TYPE_EMAIL_OTP,
    AUTH_FACTOR_TYPE_SMS_OTP,
    AUTH_FACTOR_TYPE_WEB_AUTHN,
    AUTH_FACTOR_TYPE_PASSWORD,
  ];

  static final $core.List<AuthFactorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static AuthFactorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AuthFactorType._(super.value, super.name);
}

class AuthFlowType extends $pb.ProtobufEnum {
  static const AuthFlowType AUTH_FLOW_TYPE_ADMIN_NO_SRP_AUTH = AuthFlowType._(
      0, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_ADMIN_NO_SRP_AUTH');
  static const AuthFlowType AUTH_FLOW_TYPE_ADMIN_USER_PASSWORD_AUTH =
      AuthFlowType._(
          1, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_ADMIN_USER_PASSWORD_AUTH');
  static const AuthFlowType AUTH_FLOW_TYPE_USER_AUTH =
      AuthFlowType._(2, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_USER_AUTH');
  static const AuthFlowType AUTH_FLOW_TYPE_USER_SRP_AUTH =
      AuthFlowType._(3, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_USER_SRP_AUTH');
  static const AuthFlowType AUTH_FLOW_TYPE_USER_PASSWORD_AUTH = AuthFlowType._(
      4, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_USER_PASSWORD_AUTH');
  static const AuthFlowType AUTH_FLOW_TYPE_CUSTOM_AUTH =
      AuthFlowType._(5, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_CUSTOM_AUTH');
  static const AuthFlowType AUTH_FLOW_TYPE_REFRESH_TOKEN_AUTH = AuthFlowType._(
      6, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_REFRESH_TOKEN_AUTH');
  static const AuthFlowType AUTH_FLOW_TYPE_REFRESH_TOKEN =
      AuthFlowType._(7, _omitEnumNames ? '' : 'AUTH_FLOW_TYPE_REFRESH_TOKEN');

  static const $core.List<AuthFlowType> values = <AuthFlowType>[
    AUTH_FLOW_TYPE_ADMIN_NO_SRP_AUTH,
    AUTH_FLOW_TYPE_ADMIN_USER_PASSWORD_AUTH,
    AUTH_FLOW_TYPE_USER_AUTH,
    AUTH_FLOW_TYPE_USER_SRP_AUTH,
    AUTH_FLOW_TYPE_USER_PASSWORD_AUTH,
    AUTH_FLOW_TYPE_CUSTOM_AUTH,
    AUTH_FLOW_TYPE_REFRESH_TOKEN_AUTH,
    AUTH_FLOW_TYPE_REFRESH_TOKEN,
  ];

  static final $core.List<AuthFlowType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static AuthFlowType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AuthFlowType._(super.value, super.name);
}

class ChallengeName extends $pb.ProtobufEnum {
  static const ChallengeName CHALLENGE_NAME_PASSWORD =
      ChallengeName._(0, _omitEnumNames ? '' : 'CHALLENGE_NAME_PASSWORD');
  static const ChallengeName CHALLENGE_NAME_MFA =
      ChallengeName._(1, _omitEnumNames ? '' : 'CHALLENGE_NAME_MFA');

  static const $core.List<ChallengeName> values = <ChallengeName>[
    CHALLENGE_NAME_PASSWORD,
    CHALLENGE_NAME_MFA,
  ];

  static final $core.List<ChallengeName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ChallengeName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChallengeName._(super.value, super.name);
}

class ChallengeNameType extends $pb.ProtobufEnum {
  static const ChallengeNameType CHALLENGE_NAME_TYPE_EMAIL_OTP =
      ChallengeNameType._(
          0, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_EMAIL_OTP');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_DEVICE_PASSWORD_VERIFIER =
      ChallengeNameType._(1,
          _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_DEVICE_PASSWORD_VERIFIER');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_ADMIN_NO_SRP_AUTH =
      ChallengeNameType._(
          2, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_ADMIN_NO_SRP_AUTH');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_DEVICE_SRP_AUTH =
      ChallengeNameType._(
          3, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_DEVICE_SRP_AUTH');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_SMS_OTP =
      ChallengeNameType._(
          4, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_SMS_OTP');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_PASSWORD_SRP =
      ChallengeNameType._(
          5, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_PASSWORD_SRP');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_SELECT_CHALLENGE =
      ChallengeNameType._(
          6, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_SELECT_CHALLENGE');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_PASSWORD_VERIFIER =
      ChallengeNameType._(
          7, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_PASSWORD_VERIFIER');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_NEW_PASSWORD_REQUIRED =
      ChallengeNameType._(
          8, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_NEW_PASSWORD_REQUIRED');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_WEB_AUTHN =
      ChallengeNameType._(
          9, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_WEB_AUTHN');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_CUSTOM_CHALLENGE =
      ChallengeNameType._(
          10, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_CUSTOM_CHALLENGE');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_SOFTWARE_TOKEN_MFA =
      ChallengeNameType._(
          11, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_SOFTWARE_TOKEN_MFA');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_SELECT_MFA_TYPE =
      ChallengeNameType._(
          12, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_SELECT_MFA_TYPE');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_MFA_SETUP =
      ChallengeNameType._(
          13, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_MFA_SETUP');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_PASSWORD =
      ChallengeNameType._(
          14, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_PASSWORD');
  static const ChallengeNameType CHALLENGE_NAME_TYPE_SMS_MFA =
      ChallengeNameType._(
          15, _omitEnumNames ? '' : 'CHALLENGE_NAME_TYPE_SMS_MFA');

  static const $core.List<ChallengeNameType> values = <ChallengeNameType>[
    CHALLENGE_NAME_TYPE_EMAIL_OTP,
    CHALLENGE_NAME_TYPE_DEVICE_PASSWORD_VERIFIER,
    CHALLENGE_NAME_TYPE_ADMIN_NO_SRP_AUTH,
    CHALLENGE_NAME_TYPE_DEVICE_SRP_AUTH,
    CHALLENGE_NAME_TYPE_SMS_OTP,
    CHALLENGE_NAME_TYPE_PASSWORD_SRP,
    CHALLENGE_NAME_TYPE_SELECT_CHALLENGE,
    CHALLENGE_NAME_TYPE_PASSWORD_VERIFIER,
    CHALLENGE_NAME_TYPE_NEW_PASSWORD_REQUIRED,
    CHALLENGE_NAME_TYPE_WEB_AUTHN,
    CHALLENGE_NAME_TYPE_CUSTOM_CHALLENGE,
    CHALLENGE_NAME_TYPE_SOFTWARE_TOKEN_MFA,
    CHALLENGE_NAME_TYPE_SELECT_MFA_TYPE,
    CHALLENGE_NAME_TYPE_MFA_SETUP,
    CHALLENGE_NAME_TYPE_PASSWORD,
    CHALLENGE_NAME_TYPE_SMS_MFA,
  ];

  static final $core.List<ChallengeNameType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 15);
  static ChallengeNameType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChallengeNameType._(super.value, super.name);
}

class ChallengeResponse extends $pb.ProtobufEnum {
  static const ChallengeResponse CHALLENGE_RESPONSE_FAILURE =
      ChallengeResponse._(
          0, _omitEnumNames ? '' : 'CHALLENGE_RESPONSE_FAILURE');
  static const ChallengeResponse CHALLENGE_RESPONSE_SUCCESS =
      ChallengeResponse._(
          1, _omitEnumNames ? '' : 'CHALLENGE_RESPONSE_SUCCESS');

  static const $core.List<ChallengeResponse> values = <ChallengeResponse>[
    CHALLENGE_RESPONSE_FAILURE,
    CHALLENGE_RESPONSE_SUCCESS,
  ];

  static final $core.List<ChallengeResponse?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ChallengeResponse? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChallengeResponse._(super.value, super.name);
}

class ColorSchemeModeType extends $pb.ProtobufEnum {
  static const ColorSchemeModeType COLOR_SCHEME_MODE_TYPE_LIGHT =
      ColorSchemeModeType._(
          0, _omitEnumNames ? '' : 'COLOR_SCHEME_MODE_TYPE_LIGHT');
  static const ColorSchemeModeType COLOR_SCHEME_MODE_TYPE_DYNAMIC =
      ColorSchemeModeType._(
          1, _omitEnumNames ? '' : 'COLOR_SCHEME_MODE_TYPE_DYNAMIC');
  static const ColorSchemeModeType COLOR_SCHEME_MODE_TYPE_DARK =
      ColorSchemeModeType._(
          2, _omitEnumNames ? '' : 'COLOR_SCHEME_MODE_TYPE_DARK');

  static const $core.List<ColorSchemeModeType> values = <ColorSchemeModeType>[
    COLOR_SCHEME_MODE_TYPE_LIGHT,
    COLOR_SCHEME_MODE_TYPE_DYNAMIC,
    COLOR_SCHEME_MODE_TYPE_DARK,
  ];

  static final $core.List<ColorSchemeModeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ColorSchemeModeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ColorSchemeModeType._(super.value, super.name);
}

class CompromisedCredentialsEventActionType extends $pb.ProtobufEnum {
  static const CompromisedCredentialsEventActionType
      COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_BLOCK =
      CompromisedCredentialsEventActionType._(
          0,
          _omitEnumNames
              ? ''
              : 'COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_BLOCK');
  static const CompromisedCredentialsEventActionType
      COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_NO_ACTION =
      CompromisedCredentialsEventActionType._(
          1,
          _omitEnumNames
              ? ''
              : 'COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_NO_ACTION');

  static const $core.List<CompromisedCredentialsEventActionType> values =
      <CompromisedCredentialsEventActionType>[
    COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_BLOCK,
    COMPROMISED_CREDENTIALS_EVENT_ACTION_TYPE_NO_ACTION,
  ];

  static final $core.List<CompromisedCredentialsEventActionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CompromisedCredentialsEventActionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CompromisedCredentialsEventActionType._(super.value, super.name);
}

class CustomEmailSenderLambdaVersionType extends $pb.ProtobufEnum {
  static const CustomEmailSenderLambdaVersionType
      CUSTOM_EMAIL_SENDER_LAMBDA_VERSION_TYPE_V1_0 =
      CustomEmailSenderLambdaVersionType._(0,
          _omitEnumNames ? '' : 'CUSTOM_EMAIL_SENDER_LAMBDA_VERSION_TYPE_V1_0');

  static const $core.List<CustomEmailSenderLambdaVersionType> values =
      <CustomEmailSenderLambdaVersionType>[
    CUSTOM_EMAIL_SENDER_LAMBDA_VERSION_TYPE_V1_0,
  ];

  static final $core.List<CustomEmailSenderLambdaVersionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static CustomEmailSenderLambdaVersionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CustomEmailSenderLambdaVersionType._(super.value, super.name);
}

class CustomSMSSenderLambdaVersionType extends $pb.ProtobufEnum {
  static const CustomSMSSenderLambdaVersionType
      CUSTOM_S_M_S_SENDER_LAMBDA_VERSION_TYPE_V1_0 =
      CustomSMSSenderLambdaVersionType._(0,
          _omitEnumNames ? '' : 'CUSTOM_S_M_S_SENDER_LAMBDA_VERSION_TYPE_V1_0');

  static const $core.List<CustomSMSSenderLambdaVersionType> values =
      <CustomSMSSenderLambdaVersionType>[
    CUSTOM_S_M_S_SENDER_LAMBDA_VERSION_TYPE_V1_0,
  ];

  static final $core.List<CustomSMSSenderLambdaVersionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static CustomSMSSenderLambdaVersionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CustomSMSSenderLambdaVersionType._(super.value, super.name);
}

class DefaultEmailOptionType extends $pb.ProtobufEnum {
  static const DefaultEmailOptionType
      DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_LINK = DefaultEmailOptionType._(0,
          _omitEnumNames ? '' : 'DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_LINK');
  static const DefaultEmailOptionType
      DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_CODE = DefaultEmailOptionType._(1,
          _omitEnumNames ? '' : 'DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_CODE');

  static const $core.List<DefaultEmailOptionType> values =
      <DefaultEmailOptionType>[
    DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_LINK,
    DEFAULT_EMAIL_OPTION_TYPE_CONFIRM_WITH_CODE,
  ];

  static final $core.List<DefaultEmailOptionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DefaultEmailOptionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DefaultEmailOptionType._(super.value, super.name);
}

class DeletionProtectionType extends $pb.ProtobufEnum {
  static const DeletionProtectionType DELETION_PROTECTION_TYPE_ACTIVE =
      DeletionProtectionType._(
          0, _omitEnumNames ? '' : 'DELETION_PROTECTION_TYPE_ACTIVE');
  static const DeletionProtectionType DELETION_PROTECTION_TYPE_INACTIVE =
      DeletionProtectionType._(
          1, _omitEnumNames ? '' : 'DELETION_PROTECTION_TYPE_INACTIVE');

  static const $core.List<DeletionProtectionType> values =
      <DeletionProtectionType>[
    DELETION_PROTECTION_TYPE_ACTIVE,
    DELETION_PROTECTION_TYPE_INACTIVE,
  ];

  static final $core.List<DeletionProtectionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DeletionProtectionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeletionProtectionType._(super.value, super.name);
}

class DeliveryMediumType extends $pb.ProtobufEnum {
  static const DeliveryMediumType DELIVERY_MEDIUM_TYPE_EMAIL =
      DeliveryMediumType._(
          0, _omitEnumNames ? '' : 'DELIVERY_MEDIUM_TYPE_EMAIL');
  static const DeliveryMediumType DELIVERY_MEDIUM_TYPE_SMS =
      DeliveryMediumType._(1, _omitEnumNames ? '' : 'DELIVERY_MEDIUM_TYPE_SMS');

  static const $core.List<DeliveryMediumType> values = <DeliveryMediumType>[
    DELIVERY_MEDIUM_TYPE_EMAIL,
    DELIVERY_MEDIUM_TYPE_SMS,
  ];

  static final $core.List<DeliveryMediumType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DeliveryMediumType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeliveryMediumType._(super.value, super.name);
}

class DeviceRememberedStatusType extends $pb.ProtobufEnum {
  static const DeviceRememberedStatusType
      DEVICE_REMEMBERED_STATUS_TYPE_NOT_REMEMBERED =
      DeviceRememberedStatusType._(0,
          _omitEnumNames ? '' : 'DEVICE_REMEMBERED_STATUS_TYPE_NOT_REMEMBERED');
  static const DeviceRememberedStatusType
      DEVICE_REMEMBERED_STATUS_TYPE_REMEMBERED = DeviceRememberedStatusType._(
          1, _omitEnumNames ? '' : 'DEVICE_REMEMBERED_STATUS_TYPE_REMEMBERED');

  static const $core.List<DeviceRememberedStatusType> values =
      <DeviceRememberedStatusType>[
    DEVICE_REMEMBERED_STATUS_TYPE_NOT_REMEMBERED,
    DEVICE_REMEMBERED_STATUS_TYPE_REMEMBERED,
  ];

  static final $core.List<DeviceRememberedStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DeviceRememberedStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeviceRememberedStatusType._(super.value, super.name);
}

class DomainStatusType extends $pb.ProtobufEnum {
  static const DomainStatusType DOMAIN_STATUS_TYPE_UPDATING =
      DomainStatusType._(
          0, _omitEnumNames ? '' : 'DOMAIN_STATUS_TYPE_UPDATING');
  static const DomainStatusType DOMAIN_STATUS_TYPE_ACTIVE =
      DomainStatusType._(1, _omitEnumNames ? '' : 'DOMAIN_STATUS_TYPE_ACTIVE');
  static const DomainStatusType DOMAIN_STATUS_TYPE_DELETING =
      DomainStatusType._(
          2, _omitEnumNames ? '' : 'DOMAIN_STATUS_TYPE_DELETING');
  static const DomainStatusType DOMAIN_STATUS_TYPE_CREATING =
      DomainStatusType._(
          3, _omitEnumNames ? '' : 'DOMAIN_STATUS_TYPE_CREATING');
  static const DomainStatusType DOMAIN_STATUS_TYPE_FAILED =
      DomainStatusType._(4, _omitEnumNames ? '' : 'DOMAIN_STATUS_TYPE_FAILED');

  static const $core.List<DomainStatusType> values = <DomainStatusType>[
    DOMAIN_STATUS_TYPE_UPDATING,
    DOMAIN_STATUS_TYPE_ACTIVE,
    DOMAIN_STATUS_TYPE_DELETING,
    DOMAIN_STATUS_TYPE_CREATING,
    DOMAIN_STATUS_TYPE_FAILED,
  ];

  static final $core.List<DomainStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static DomainStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DomainStatusType._(super.value, super.name);
}

class EmailSendingAccountType extends $pb.ProtobufEnum {
  static const EmailSendingAccountType
      EMAIL_SENDING_ACCOUNT_TYPE_COGNITO_DEFAULT = EmailSendingAccountType._(0,
          _omitEnumNames ? '' : 'EMAIL_SENDING_ACCOUNT_TYPE_COGNITO_DEFAULT');
  static const EmailSendingAccountType EMAIL_SENDING_ACCOUNT_TYPE_DEVELOPER =
      EmailSendingAccountType._(
          1, _omitEnumNames ? '' : 'EMAIL_SENDING_ACCOUNT_TYPE_DEVELOPER');

  static const $core.List<EmailSendingAccountType> values =
      <EmailSendingAccountType>[
    EMAIL_SENDING_ACCOUNT_TYPE_COGNITO_DEFAULT,
    EMAIL_SENDING_ACCOUNT_TYPE_DEVELOPER,
  ];

  static final $core.List<EmailSendingAccountType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static EmailSendingAccountType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EmailSendingAccountType._(super.value, super.name);
}

class EventFilterType extends $pb.ProtobufEnum {
  static const EventFilterType EVENT_FILTER_TYPE_PASSWORD_CHANGE =
      EventFilterType._(
          0, _omitEnumNames ? '' : 'EVENT_FILTER_TYPE_PASSWORD_CHANGE');
  static const EventFilterType EVENT_FILTER_TYPE_SIGN_IN =
      EventFilterType._(1, _omitEnumNames ? '' : 'EVENT_FILTER_TYPE_SIGN_IN');
  static const EventFilterType EVENT_FILTER_TYPE_SIGN_UP =
      EventFilterType._(2, _omitEnumNames ? '' : 'EVENT_FILTER_TYPE_SIGN_UP');

  static const $core.List<EventFilterType> values = <EventFilterType>[
    EVENT_FILTER_TYPE_PASSWORD_CHANGE,
    EVENT_FILTER_TYPE_SIGN_IN,
    EVENT_FILTER_TYPE_SIGN_UP,
  ];

  static final $core.List<EventFilterType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EventFilterType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventFilterType._(super.value, super.name);
}

class EventResponseType extends $pb.ProtobufEnum {
  static const EventResponseType EVENT_RESPONSE_TYPE_FAIL =
      EventResponseType._(0, _omitEnumNames ? '' : 'EVENT_RESPONSE_TYPE_FAIL');
  static const EventResponseType EVENT_RESPONSE_TYPE_PASS =
      EventResponseType._(1, _omitEnumNames ? '' : 'EVENT_RESPONSE_TYPE_PASS');
  static const EventResponseType EVENT_RESPONSE_TYPE_INPROGRESS =
      EventResponseType._(
          2, _omitEnumNames ? '' : 'EVENT_RESPONSE_TYPE_INPROGRESS');

  static const $core.List<EventResponseType> values = <EventResponseType>[
    EVENT_RESPONSE_TYPE_FAIL,
    EVENT_RESPONSE_TYPE_PASS,
    EVENT_RESPONSE_TYPE_INPROGRESS,
  ];

  static final $core.List<EventResponseType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EventResponseType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventResponseType._(super.value, super.name);
}

class EventSourceName extends $pb.ProtobufEnum {
  static const EventSourceName EVENT_SOURCE_NAME_USER_AUTH_EVENTS =
      EventSourceName._(
          0, _omitEnumNames ? '' : 'EVENT_SOURCE_NAME_USER_AUTH_EVENTS');
  static const EventSourceName EVENT_SOURCE_NAME_USER_NOTIFICATION =
      EventSourceName._(
          1, _omitEnumNames ? '' : 'EVENT_SOURCE_NAME_USER_NOTIFICATION');

  static const $core.List<EventSourceName> values = <EventSourceName>[
    EVENT_SOURCE_NAME_USER_AUTH_EVENTS,
    EVENT_SOURCE_NAME_USER_NOTIFICATION,
  ];

  static final $core.List<EventSourceName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static EventSourceName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventSourceName._(super.value, super.name);
}

class EventType extends $pb.ProtobufEnum {
  static const EventType EVENT_TYPE_SIGNIN =
      EventType._(0, _omitEnumNames ? '' : 'EVENT_TYPE_SIGNIN');
  static const EventType EVENT_TYPE_RESENDCODE =
      EventType._(1, _omitEnumNames ? '' : 'EVENT_TYPE_RESENDCODE');
  static const EventType EVENT_TYPE_PASSWORDCHANGE =
      EventType._(2, _omitEnumNames ? '' : 'EVENT_TYPE_PASSWORDCHANGE');
  static const EventType EVENT_TYPE_SIGNUP =
      EventType._(3, _omitEnumNames ? '' : 'EVENT_TYPE_SIGNUP');
  static const EventType EVENT_TYPE_FORGOTPASSWORD =
      EventType._(4, _omitEnumNames ? '' : 'EVENT_TYPE_FORGOTPASSWORD');

  static const $core.List<EventType> values = <EventType>[
    EVENT_TYPE_SIGNIN,
    EVENT_TYPE_RESENDCODE,
    EVENT_TYPE_PASSWORDCHANGE,
    EVENT_TYPE_SIGNUP,
    EVENT_TYPE_FORGOTPASSWORD,
  ];

  static final $core.List<EventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static EventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventType._(super.value, super.name);
}

class ExplicitAuthFlowsType extends $pb.ProtobufEnum {
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_CUSTOM_AUTH_FLOW_ONLY = ExplicitAuthFlowsType._(
          0,
          _omitEnumNames
              ? ''
              : 'EXPLICIT_AUTH_FLOWS_TYPE_CUSTOM_AUTH_FLOW_ONLY');
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_SRP_AUTH = ExplicitAuthFlowsType._(1,
          _omitEnumNames ? '' : 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_SRP_AUTH');
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_CUSTOM_AUTH = ExplicitAuthFlowsType._(2,
          _omitEnumNames ? '' : 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_CUSTOM_AUTH');
  static const ExplicitAuthFlowsType EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_AUTH =
      ExplicitAuthFlowsType._(
          3, _omitEnumNames ? '' : 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_AUTH');
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_ADMIN_NO_SRP_AUTH = ExplicitAuthFlowsType._(4,
          _omitEnumNames ? '' : 'EXPLICIT_AUTH_FLOWS_TYPE_ADMIN_NO_SRP_AUTH');
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_USER_PASSWORD_AUTH = ExplicitAuthFlowsType._(5,
          _omitEnumNames ? '' : 'EXPLICIT_AUTH_FLOWS_TYPE_USER_PASSWORD_AUTH');
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_ADMIN_USER_PASSWORD_AUTH =
      ExplicitAuthFlowsType._(
          6,
          _omitEnumNames
              ? ''
              : 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_ADMIN_USER_PASSWORD_AUTH');
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_PASSWORD_AUTH =
      ExplicitAuthFlowsType._(
          7,
          _omitEnumNames
              ? ''
              : 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_PASSWORD_AUTH');
  static const ExplicitAuthFlowsType
      EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_REFRESH_TOKEN_AUTH =
      ExplicitAuthFlowsType._(
          8,
          _omitEnumNames
              ? ''
              : 'EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_REFRESH_TOKEN_AUTH');

  static const $core.List<ExplicitAuthFlowsType> values =
      <ExplicitAuthFlowsType>[
    EXPLICIT_AUTH_FLOWS_TYPE_CUSTOM_AUTH_FLOW_ONLY,
    EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_SRP_AUTH,
    EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_CUSTOM_AUTH,
    EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_AUTH,
    EXPLICIT_AUTH_FLOWS_TYPE_ADMIN_NO_SRP_AUTH,
    EXPLICIT_AUTH_FLOWS_TYPE_USER_PASSWORD_AUTH,
    EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_ADMIN_USER_PASSWORD_AUTH,
    EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_USER_PASSWORD_AUTH,
    EXPLICIT_AUTH_FLOWS_TYPE_ALLOW_REFRESH_TOKEN_AUTH,
  ];

  static final $core.List<ExplicitAuthFlowsType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static ExplicitAuthFlowsType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExplicitAuthFlowsType._(super.value, super.name);
}

class FeatureType extends $pb.ProtobufEnum {
  static const FeatureType FEATURE_TYPE_DISABLED =
      FeatureType._(0, _omitEnumNames ? '' : 'FEATURE_TYPE_DISABLED');
  static const FeatureType FEATURE_TYPE_ENABLED =
      FeatureType._(1, _omitEnumNames ? '' : 'FEATURE_TYPE_ENABLED');

  static const $core.List<FeatureType> values = <FeatureType>[
    FEATURE_TYPE_DISABLED,
    FEATURE_TYPE_ENABLED,
  ];

  static final $core.List<FeatureType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FeatureType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FeatureType._(super.value, super.name);
}

class FeedbackValueType extends $pb.ProtobufEnum {
  static const FeedbackValueType FEEDBACK_VALUE_TYPE_INVALID =
      FeedbackValueType._(
          0, _omitEnumNames ? '' : 'FEEDBACK_VALUE_TYPE_INVALID');
  static const FeedbackValueType FEEDBACK_VALUE_TYPE_VALID =
      FeedbackValueType._(1, _omitEnumNames ? '' : 'FEEDBACK_VALUE_TYPE_VALID');

  static const $core.List<FeedbackValueType> values = <FeedbackValueType>[
    FEEDBACK_VALUE_TYPE_INVALID,
    FEEDBACK_VALUE_TYPE_VALID,
  ];

  static final $core.List<FeedbackValueType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FeedbackValueType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FeedbackValueType._(super.value, super.name);
}

class IdentityProviderTypeType extends $pb.ProtobufEnum {
  static const IdentityProviderTypeType
      IDENTITY_PROVIDER_TYPE_TYPE_LOGINWITHAMAZON = IdentityProviderTypeType._(
          0,
          _omitEnumNames ? '' : 'IDENTITY_PROVIDER_TYPE_TYPE_LOGINWITHAMAZON');
  static const IdentityProviderTypeType IDENTITY_PROVIDER_TYPE_TYPE_GOOGLE =
      IdentityProviderTypeType._(
          1, _omitEnumNames ? '' : 'IDENTITY_PROVIDER_TYPE_TYPE_GOOGLE');
  static const IdentityProviderTypeType IDENTITY_PROVIDER_TYPE_TYPE_FACEBOOK =
      IdentityProviderTypeType._(
          2, _omitEnumNames ? '' : 'IDENTITY_PROVIDER_TYPE_TYPE_FACEBOOK');
  static const IdentityProviderTypeType IDENTITY_PROVIDER_TYPE_TYPE_SAML =
      IdentityProviderTypeType._(
          3, _omitEnumNames ? '' : 'IDENTITY_PROVIDER_TYPE_TYPE_SAML');
  static const IdentityProviderTypeType
      IDENTITY_PROVIDER_TYPE_TYPE_SIGNINWITHAPPLE = IdentityProviderTypeType._(
          4,
          _omitEnumNames ? '' : 'IDENTITY_PROVIDER_TYPE_TYPE_SIGNINWITHAPPLE');
  static const IdentityProviderTypeType IDENTITY_PROVIDER_TYPE_TYPE_OIDC =
      IdentityProviderTypeType._(
          5, _omitEnumNames ? '' : 'IDENTITY_PROVIDER_TYPE_TYPE_OIDC');

  static const $core.List<IdentityProviderTypeType> values =
      <IdentityProviderTypeType>[
    IDENTITY_PROVIDER_TYPE_TYPE_LOGINWITHAMAZON,
    IDENTITY_PROVIDER_TYPE_TYPE_GOOGLE,
    IDENTITY_PROVIDER_TYPE_TYPE_FACEBOOK,
    IDENTITY_PROVIDER_TYPE_TYPE_SAML,
    IDENTITY_PROVIDER_TYPE_TYPE_SIGNINWITHAPPLE,
    IDENTITY_PROVIDER_TYPE_TYPE_OIDC,
  ];

  static final $core.List<IdentityProviderTypeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static IdentityProviderTypeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IdentityProviderTypeType._(super.value, super.name);
}

class InboundFederationLambdaVersionType extends $pb.ProtobufEnum {
  static const InboundFederationLambdaVersionType
      INBOUND_FEDERATION_LAMBDA_VERSION_TYPE_V1_0 =
      InboundFederationLambdaVersionType._(0,
          _omitEnumNames ? '' : 'INBOUND_FEDERATION_LAMBDA_VERSION_TYPE_V1_0');

  static const $core.List<InboundFederationLambdaVersionType> values =
      <InboundFederationLambdaVersionType>[
    INBOUND_FEDERATION_LAMBDA_VERSION_TYPE_V1_0,
  ];

  static final $core.List<InboundFederationLambdaVersionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static InboundFederationLambdaVersionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InboundFederationLambdaVersionType._(super.value, super.name);
}

class LogLevel extends $pb.ProtobufEnum {
  static const LogLevel LOG_LEVEL_INFO =
      LogLevel._(0, _omitEnumNames ? '' : 'LOG_LEVEL_INFO');
  static const LogLevel LOG_LEVEL_ERROR =
      LogLevel._(1, _omitEnumNames ? '' : 'LOG_LEVEL_ERROR');

  static const $core.List<LogLevel> values = <LogLevel>[
    LOG_LEVEL_INFO,
    LOG_LEVEL_ERROR,
  ];

  static final $core.List<LogLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static LogLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogLevel._(super.value, super.name);
}

class MessageActionType extends $pb.ProtobufEnum {
  static const MessageActionType MESSAGE_ACTION_TYPE_SUPPRESS =
      MessageActionType._(
          0, _omitEnumNames ? '' : 'MESSAGE_ACTION_TYPE_SUPPRESS');
  static const MessageActionType MESSAGE_ACTION_TYPE_RESEND =
      MessageActionType._(
          1, _omitEnumNames ? '' : 'MESSAGE_ACTION_TYPE_RESEND');

  static const $core.List<MessageActionType> values = <MessageActionType>[
    MESSAGE_ACTION_TYPE_SUPPRESS,
    MESSAGE_ACTION_TYPE_RESEND,
  ];

  static final $core.List<MessageActionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MessageActionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MessageActionType._(super.value, super.name);
}

class OAuthFlowType extends $pb.ProtobufEnum {
  static const OAuthFlowType O_AUTH_FLOW_TYPE_CLIENT_CREDENTIALS =
      OAuthFlowType._(
          0, _omitEnumNames ? '' : 'O_AUTH_FLOW_TYPE_CLIENT_CREDENTIALS');
  static const OAuthFlowType O_AUTH_FLOW_TYPE_IMPLICIT =
      OAuthFlowType._(1, _omitEnumNames ? '' : 'O_AUTH_FLOW_TYPE_IMPLICIT');
  static const OAuthFlowType O_AUTH_FLOW_TYPE_CODE =
      OAuthFlowType._(2, _omitEnumNames ? '' : 'O_AUTH_FLOW_TYPE_CODE');

  static const $core.List<OAuthFlowType> values = <OAuthFlowType>[
    O_AUTH_FLOW_TYPE_CLIENT_CREDENTIALS,
    O_AUTH_FLOW_TYPE_IMPLICIT,
    O_AUTH_FLOW_TYPE_CODE,
  ];

  static final $core.List<OAuthFlowType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static OAuthFlowType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OAuthFlowType._(super.value, super.name);
}

class PreTokenGenerationLambdaVersionType extends $pb.ProtobufEnum {
  static const PreTokenGenerationLambdaVersionType
      PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V1_0 =
      PreTokenGenerationLambdaVersionType._(
          0,
          _omitEnumNames
              ? ''
              : 'PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V1_0');
  static const PreTokenGenerationLambdaVersionType
      PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V3_0 =
      PreTokenGenerationLambdaVersionType._(
          1,
          _omitEnumNames
              ? ''
              : 'PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V3_0');
  static const PreTokenGenerationLambdaVersionType
      PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V2_0 =
      PreTokenGenerationLambdaVersionType._(
          2,
          _omitEnumNames
              ? ''
              : 'PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V2_0');

  static const $core.List<PreTokenGenerationLambdaVersionType> values =
      <PreTokenGenerationLambdaVersionType>[
    PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V1_0,
    PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V3_0,
    PRE_TOKEN_GENERATION_LAMBDA_VERSION_TYPE_V2_0,
  ];

  static final $core.List<PreTokenGenerationLambdaVersionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static PreTokenGenerationLambdaVersionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PreTokenGenerationLambdaVersionType._(super.value, super.name);
}

class PreventUserExistenceErrorTypes extends $pb.ProtobufEnum {
  static const PreventUserExistenceErrorTypes
      PREVENT_USER_EXISTENCE_ERROR_TYPES_LEGACY =
      PreventUserExistenceErrorTypes._(
          0, _omitEnumNames ? '' : 'PREVENT_USER_EXISTENCE_ERROR_TYPES_LEGACY');
  static const PreventUserExistenceErrorTypes
      PREVENT_USER_EXISTENCE_ERROR_TYPES_ENABLED =
      PreventUserExistenceErrorTypes._(1,
          _omitEnumNames ? '' : 'PREVENT_USER_EXISTENCE_ERROR_TYPES_ENABLED');

  static const $core.List<PreventUserExistenceErrorTypes> values =
      <PreventUserExistenceErrorTypes>[
    PREVENT_USER_EXISTENCE_ERROR_TYPES_LEGACY,
    PREVENT_USER_EXISTENCE_ERROR_TYPES_ENABLED,
  ];

  static final $core.List<PreventUserExistenceErrorTypes?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PreventUserExistenceErrorTypes? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PreventUserExistenceErrorTypes._(super.value, super.name);
}

class RecoveryOptionNameType extends $pb.ProtobufEnum {
  static const RecoveryOptionNameType RECOVERY_OPTION_NAME_TYPE_ADMIN_ONLY =
      RecoveryOptionNameType._(
          0, _omitEnumNames ? '' : 'RECOVERY_OPTION_NAME_TYPE_ADMIN_ONLY');
  static const RecoveryOptionNameType RECOVERY_OPTION_NAME_TYPE_VERIFIED_EMAIL =
      RecoveryOptionNameType._(
          1, _omitEnumNames ? '' : 'RECOVERY_OPTION_NAME_TYPE_VERIFIED_EMAIL');
  static const RecoveryOptionNameType
      RECOVERY_OPTION_NAME_TYPE_VERIFIED_PHONE_NUMBER =
      RecoveryOptionNameType._(
          2,
          _omitEnumNames
              ? ''
              : 'RECOVERY_OPTION_NAME_TYPE_VERIFIED_PHONE_NUMBER');

  static const $core.List<RecoveryOptionNameType> values =
      <RecoveryOptionNameType>[
    RECOVERY_OPTION_NAME_TYPE_ADMIN_ONLY,
    RECOVERY_OPTION_NAME_TYPE_VERIFIED_EMAIL,
    RECOVERY_OPTION_NAME_TYPE_VERIFIED_PHONE_NUMBER,
  ];

  static final $core.List<RecoveryOptionNameType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static RecoveryOptionNameType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecoveryOptionNameType._(super.value, super.name);
}

class RiskDecisionType extends $pb.ProtobufEnum {
  static const RiskDecisionType RISK_DECISION_TYPE_ACCOUNTTAKEOVER =
      RiskDecisionType._(
          0, _omitEnumNames ? '' : 'RISK_DECISION_TYPE_ACCOUNTTAKEOVER');
  static const RiskDecisionType RISK_DECISION_TYPE_NORISK =
      RiskDecisionType._(1, _omitEnumNames ? '' : 'RISK_DECISION_TYPE_NORISK');
  static const RiskDecisionType RISK_DECISION_TYPE_BLOCK =
      RiskDecisionType._(2, _omitEnumNames ? '' : 'RISK_DECISION_TYPE_BLOCK');

  static const $core.List<RiskDecisionType> values = <RiskDecisionType>[
    RISK_DECISION_TYPE_ACCOUNTTAKEOVER,
    RISK_DECISION_TYPE_NORISK,
    RISK_DECISION_TYPE_BLOCK,
  ];

  static final $core.List<RiskDecisionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static RiskDecisionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RiskDecisionType._(super.value, super.name);
}

class RiskLevelType extends $pb.ProtobufEnum {
  static const RiskLevelType RISK_LEVEL_TYPE_MEDIUM =
      RiskLevelType._(0, _omitEnumNames ? '' : 'RISK_LEVEL_TYPE_MEDIUM');
  static const RiskLevelType RISK_LEVEL_TYPE_LOW =
      RiskLevelType._(1, _omitEnumNames ? '' : 'RISK_LEVEL_TYPE_LOW');
  static const RiskLevelType RISK_LEVEL_TYPE_HIGH =
      RiskLevelType._(2, _omitEnumNames ? '' : 'RISK_LEVEL_TYPE_HIGH');

  static const $core.List<RiskLevelType> values = <RiskLevelType>[
    RISK_LEVEL_TYPE_MEDIUM,
    RISK_LEVEL_TYPE_LOW,
    RISK_LEVEL_TYPE_HIGH,
  ];

  static final $core.List<RiskLevelType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static RiskLevelType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RiskLevelType._(super.value, super.name);
}

class StatusType extends $pb.ProtobufEnum {
  static const StatusType STATUS_TYPE_DISABLED =
      StatusType._(0, _omitEnumNames ? '' : 'STATUS_TYPE_DISABLED');
  static const StatusType STATUS_TYPE_ENABLED =
      StatusType._(1, _omitEnumNames ? '' : 'STATUS_TYPE_ENABLED');

  static const $core.List<StatusType> values = <StatusType>[
    STATUS_TYPE_DISABLED,
    STATUS_TYPE_ENABLED,
  ];

  static final $core.List<StatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static StatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StatusType._(super.value, super.name);
}

class TermsEnforcementType extends $pb.ProtobufEnum {
  static const TermsEnforcementType TERMS_ENFORCEMENT_TYPE_NONE =
      TermsEnforcementType._(
          0, _omitEnumNames ? '' : 'TERMS_ENFORCEMENT_TYPE_NONE');

  static const $core.List<TermsEnforcementType> values = <TermsEnforcementType>[
    TERMS_ENFORCEMENT_TYPE_NONE,
  ];

  static final $core.List<TermsEnforcementType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static TermsEnforcementType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TermsEnforcementType._(super.value, super.name);
}

class TermsSourceType extends $pb.ProtobufEnum {
  static const TermsSourceType TERMS_SOURCE_TYPE_LINK =
      TermsSourceType._(0, _omitEnumNames ? '' : 'TERMS_SOURCE_TYPE_LINK');

  static const $core.List<TermsSourceType> values = <TermsSourceType>[
    TERMS_SOURCE_TYPE_LINK,
  ];

  static final $core.List<TermsSourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static TermsSourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TermsSourceType._(super.value, super.name);
}

class TimeUnitsType extends $pb.ProtobufEnum {
  static const TimeUnitsType TIME_UNITS_TYPE_MINUTES =
      TimeUnitsType._(0, _omitEnumNames ? '' : 'TIME_UNITS_TYPE_MINUTES');
  static const TimeUnitsType TIME_UNITS_TYPE_SECONDS =
      TimeUnitsType._(1, _omitEnumNames ? '' : 'TIME_UNITS_TYPE_SECONDS');
  static const TimeUnitsType TIME_UNITS_TYPE_DAYS =
      TimeUnitsType._(2, _omitEnumNames ? '' : 'TIME_UNITS_TYPE_DAYS');
  static const TimeUnitsType TIME_UNITS_TYPE_HOURS =
      TimeUnitsType._(3, _omitEnumNames ? '' : 'TIME_UNITS_TYPE_HOURS');

  static const $core.List<TimeUnitsType> values = <TimeUnitsType>[
    TIME_UNITS_TYPE_MINUTES,
    TIME_UNITS_TYPE_SECONDS,
    TIME_UNITS_TYPE_DAYS,
    TIME_UNITS_TYPE_HOURS,
  ];

  static final $core.List<TimeUnitsType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static TimeUnitsType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TimeUnitsType._(super.value, super.name);
}

class UserImportJobStatusType extends $pb.ProtobufEnum {
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_FAILED =
      UserImportJobStatusType._(
          0, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_FAILED');
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_STOPPED =
      UserImportJobStatusType._(
          1, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_STOPPED');
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_EXPIRED =
      UserImportJobStatusType._(
          2, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_EXPIRED');
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_SUCCEEDED =
      UserImportJobStatusType._(
          3, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_SUCCEEDED');
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_STOPPING =
      UserImportJobStatusType._(
          4, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_STOPPING');
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_CREATED =
      UserImportJobStatusType._(
          5, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_CREATED');
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_INPROGRESS =
      UserImportJobStatusType._(
          6, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_INPROGRESS');
  static const UserImportJobStatusType USER_IMPORT_JOB_STATUS_TYPE_PENDING =
      UserImportJobStatusType._(
          7, _omitEnumNames ? '' : 'USER_IMPORT_JOB_STATUS_TYPE_PENDING');

  static const $core.List<UserImportJobStatusType> values =
      <UserImportJobStatusType>[
    USER_IMPORT_JOB_STATUS_TYPE_FAILED,
    USER_IMPORT_JOB_STATUS_TYPE_STOPPED,
    USER_IMPORT_JOB_STATUS_TYPE_EXPIRED,
    USER_IMPORT_JOB_STATUS_TYPE_SUCCEEDED,
    USER_IMPORT_JOB_STATUS_TYPE_STOPPING,
    USER_IMPORT_JOB_STATUS_TYPE_CREATED,
    USER_IMPORT_JOB_STATUS_TYPE_INPROGRESS,
    USER_IMPORT_JOB_STATUS_TYPE_PENDING,
  ];

  static final $core.List<UserImportJobStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static UserImportJobStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UserImportJobStatusType._(super.value, super.name);
}

class UserPoolMfaType extends $pb.ProtobufEnum {
  static const UserPoolMfaType USER_POOL_MFA_TYPE_OPTIONAL =
      UserPoolMfaType._(0, _omitEnumNames ? '' : 'USER_POOL_MFA_TYPE_OPTIONAL');
  static const UserPoolMfaType USER_POOL_MFA_TYPE_OFF =
      UserPoolMfaType._(1, _omitEnumNames ? '' : 'USER_POOL_MFA_TYPE_OFF');
  static const UserPoolMfaType USER_POOL_MFA_TYPE_ON =
      UserPoolMfaType._(2, _omitEnumNames ? '' : 'USER_POOL_MFA_TYPE_ON');

  static const $core.List<UserPoolMfaType> values = <UserPoolMfaType>[
    USER_POOL_MFA_TYPE_OPTIONAL,
    USER_POOL_MFA_TYPE_OFF,
    USER_POOL_MFA_TYPE_ON,
  ];

  static final $core.List<UserPoolMfaType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static UserPoolMfaType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UserPoolMfaType._(super.value, super.name);
}

class UserPoolTierType extends $pb.ProtobufEnum {
  static const UserPoolTierType USER_POOL_TIER_TYPE_PLUS =
      UserPoolTierType._(0, _omitEnumNames ? '' : 'USER_POOL_TIER_TYPE_PLUS');
  static const UserPoolTierType USER_POOL_TIER_TYPE_ESSENTIALS =
      UserPoolTierType._(
          1, _omitEnumNames ? '' : 'USER_POOL_TIER_TYPE_ESSENTIALS');
  static const UserPoolTierType USER_POOL_TIER_TYPE_LITE =
      UserPoolTierType._(2, _omitEnumNames ? '' : 'USER_POOL_TIER_TYPE_LITE');

  static const $core.List<UserPoolTierType> values = <UserPoolTierType>[
    USER_POOL_TIER_TYPE_PLUS,
    USER_POOL_TIER_TYPE_ESSENTIALS,
    USER_POOL_TIER_TYPE_LITE,
  ];

  static final $core.List<UserPoolTierType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static UserPoolTierType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UserPoolTierType._(super.value, super.name);
}

class UserStatusType extends $pb.ProtobufEnum {
  static const UserStatusType USER_STATUS_TYPE_EXTERNAL_PROVIDER =
      UserStatusType._(
          0, _omitEnumNames ? '' : 'USER_STATUS_TYPE_EXTERNAL_PROVIDER');
  static const UserStatusType USER_STATUS_TYPE_CONFIRMED =
      UserStatusType._(1, _omitEnumNames ? '' : 'USER_STATUS_TYPE_CONFIRMED');
  static const UserStatusType USER_STATUS_TYPE_UNKNOWN =
      UserStatusType._(2, _omitEnumNames ? '' : 'USER_STATUS_TYPE_UNKNOWN');
  static const UserStatusType USER_STATUS_TYPE_UNCONFIRMED =
      UserStatusType._(3, _omitEnumNames ? '' : 'USER_STATUS_TYPE_UNCONFIRMED');
  static const UserStatusType USER_STATUS_TYPE_RESET_REQUIRED =
      UserStatusType._(
          4, _omitEnumNames ? '' : 'USER_STATUS_TYPE_RESET_REQUIRED');
  static const UserStatusType USER_STATUS_TYPE_COMPROMISED =
      UserStatusType._(5, _omitEnumNames ? '' : 'USER_STATUS_TYPE_COMPROMISED');
  static const UserStatusType USER_STATUS_TYPE_FORCE_CHANGE_PASSWORD =
      UserStatusType._(
          6, _omitEnumNames ? '' : 'USER_STATUS_TYPE_FORCE_CHANGE_PASSWORD');
  static const UserStatusType USER_STATUS_TYPE_ARCHIVED =
      UserStatusType._(7, _omitEnumNames ? '' : 'USER_STATUS_TYPE_ARCHIVED');

  static const $core.List<UserStatusType> values = <UserStatusType>[
    USER_STATUS_TYPE_EXTERNAL_PROVIDER,
    USER_STATUS_TYPE_CONFIRMED,
    USER_STATUS_TYPE_UNKNOWN,
    USER_STATUS_TYPE_UNCONFIRMED,
    USER_STATUS_TYPE_RESET_REQUIRED,
    USER_STATUS_TYPE_COMPROMISED,
    USER_STATUS_TYPE_FORCE_CHANGE_PASSWORD,
    USER_STATUS_TYPE_ARCHIVED,
  ];

  static final $core.List<UserStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static UserStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UserStatusType._(super.value, super.name);
}

class UserVerificationType extends $pb.ProtobufEnum {
  static const UserVerificationType USER_VERIFICATION_TYPE_PREFERRED =
      UserVerificationType._(
          0, _omitEnumNames ? '' : 'USER_VERIFICATION_TYPE_PREFERRED');
  static const UserVerificationType USER_VERIFICATION_TYPE_REQUIRED =
      UserVerificationType._(
          1, _omitEnumNames ? '' : 'USER_VERIFICATION_TYPE_REQUIRED');

  static const $core.List<UserVerificationType> values = <UserVerificationType>[
    USER_VERIFICATION_TYPE_PREFERRED,
    USER_VERIFICATION_TYPE_REQUIRED,
  ];

  static final $core.List<UserVerificationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static UserVerificationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UserVerificationType._(super.value, super.name);
}

class UsernameAttributeType extends $pb.ProtobufEnum {
  static const UsernameAttributeType USERNAME_ATTRIBUTE_TYPE_PHONE_NUMBER =
      UsernameAttributeType._(
          0, _omitEnumNames ? '' : 'USERNAME_ATTRIBUTE_TYPE_PHONE_NUMBER');
  static const UsernameAttributeType USERNAME_ATTRIBUTE_TYPE_EMAIL =
      UsernameAttributeType._(
          1, _omitEnumNames ? '' : 'USERNAME_ATTRIBUTE_TYPE_EMAIL');

  static const $core.List<UsernameAttributeType> values =
      <UsernameAttributeType>[
    USERNAME_ATTRIBUTE_TYPE_PHONE_NUMBER,
    USERNAME_ATTRIBUTE_TYPE_EMAIL,
  ];

  static final $core.List<UsernameAttributeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static UsernameAttributeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UsernameAttributeType._(super.value, super.name);
}

class VerifiedAttributeType extends $pb.ProtobufEnum {
  static const VerifiedAttributeType VERIFIED_ATTRIBUTE_TYPE_PHONE_NUMBER =
      VerifiedAttributeType._(
          0, _omitEnumNames ? '' : 'VERIFIED_ATTRIBUTE_TYPE_PHONE_NUMBER');
  static const VerifiedAttributeType VERIFIED_ATTRIBUTE_TYPE_EMAIL =
      VerifiedAttributeType._(
          1, _omitEnumNames ? '' : 'VERIFIED_ATTRIBUTE_TYPE_EMAIL');

  static const $core.List<VerifiedAttributeType> values =
      <VerifiedAttributeType>[
    VERIFIED_ATTRIBUTE_TYPE_PHONE_NUMBER,
    VERIFIED_ATTRIBUTE_TYPE_EMAIL,
  ];

  static final $core.List<VerifiedAttributeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static VerifiedAttributeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VerifiedAttributeType._(super.value, super.name);
}

class VerifySoftwareTokenResponseType extends $pb.ProtobufEnum {
  static const VerifySoftwareTokenResponseType
      VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_SUCCESS =
      VerifySoftwareTokenResponseType._(0,
          _omitEnumNames ? '' : 'VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_SUCCESS');
  static const VerifySoftwareTokenResponseType
      VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_ERROR =
      VerifySoftwareTokenResponseType._(
          1, _omitEnumNames ? '' : 'VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_ERROR');

  static const $core.List<VerifySoftwareTokenResponseType> values =
      <VerifySoftwareTokenResponseType>[
    VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_SUCCESS,
    VERIFY_SOFTWARE_TOKEN_RESPONSE_TYPE_ERROR,
  ];

  static final $core.List<VerifySoftwareTokenResponseType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static VerifySoftwareTokenResponseType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VerifySoftwareTokenResponseType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
