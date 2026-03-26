// This is a generated file - do not edit.
//
// Generated from sns.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class LanguageCodeString extends $pb.ProtobufEnum {
  static const LanguageCodeString LANGUAGE_CODE_STRING_ZH_CN =
      LanguageCodeString._(
          0, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_ZH_CN');
  static const LanguageCodeString LANGUAGE_CODE_STRING_EN_GB =
      LanguageCodeString._(
          1, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_EN_GB');
  static const LanguageCodeString LANGUAGE_CODE_STRING_IT_IT =
      LanguageCodeString._(
          2, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_IT_IT');
  static const LanguageCodeString LANGUAGE_CODE_STRING_JP_JP =
      LanguageCodeString._(
          3, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_JP_JP');
  static const LanguageCodeString LANGUAGE_CODE_STRING_ES_ES =
      LanguageCodeString._(
          4, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_ES_ES');
  static const LanguageCodeString LANGUAGE_CODE_STRING_DE_DE =
      LanguageCodeString._(
          5, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_DE_DE');
  static const LanguageCodeString LANGUAGE_CODE_STRING_KR_KR =
      LanguageCodeString._(
          6, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_KR_KR');
  static const LanguageCodeString LANGUAGE_CODE_STRING_FR_FR =
      LanguageCodeString._(
          7, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_FR_FR');
  static const LanguageCodeString LANGUAGE_CODE_STRING_EN_US =
      LanguageCodeString._(
          8, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_EN_US');
  static const LanguageCodeString LANGUAGE_CODE_STRING_PT_BR =
      LanguageCodeString._(
          9, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_PT_BR');
  static const LanguageCodeString LANGUAGE_CODE_STRING_ZH_TW =
      LanguageCodeString._(
          10, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_ZH_TW');
  static const LanguageCodeString LANGUAGE_CODE_STRING_ES_419 =
      LanguageCodeString._(
          11, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_ES_419');
  static const LanguageCodeString LANGUAGE_CODE_STRING_FR_CA =
      LanguageCodeString._(
          12, _omitEnumNames ? '' : 'LANGUAGE_CODE_STRING_FR_CA');

  static const $core.List<LanguageCodeString> values = <LanguageCodeString>[
    LANGUAGE_CODE_STRING_ZH_CN,
    LANGUAGE_CODE_STRING_EN_GB,
    LANGUAGE_CODE_STRING_IT_IT,
    LANGUAGE_CODE_STRING_JP_JP,
    LANGUAGE_CODE_STRING_ES_ES,
    LANGUAGE_CODE_STRING_DE_DE,
    LANGUAGE_CODE_STRING_KR_KR,
    LANGUAGE_CODE_STRING_FR_FR,
    LANGUAGE_CODE_STRING_EN_US,
    LANGUAGE_CODE_STRING_PT_BR,
    LANGUAGE_CODE_STRING_ZH_TW,
    LANGUAGE_CODE_STRING_ES_419,
    LANGUAGE_CODE_STRING_FR_CA,
  ];

  static final $core.List<LanguageCodeString?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 12);
  static LanguageCodeString? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LanguageCodeString._(super.value, super.name);
}

class NumberCapability extends $pb.ProtobufEnum {
  static const NumberCapability NUMBER_CAPABILITY_SMS =
      NumberCapability._(0, _omitEnumNames ? '' : 'NUMBER_CAPABILITY_SMS');
  static const NumberCapability NUMBER_CAPABILITY_VOICE =
      NumberCapability._(1, _omitEnumNames ? '' : 'NUMBER_CAPABILITY_VOICE');
  static const NumberCapability NUMBER_CAPABILITY_MMS =
      NumberCapability._(2, _omitEnumNames ? '' : 'NUMBER_CAPABILITY_MMS');

  static const $core.List<NumberCapability> values = <NumberCapability>[
    NUMBER_CAPABILITY_SMS,
    NUMBER_CAPABILITY_VOICE,
    NUMBER_CAPABILITY_MMS,
  ];

  static final $core.List<NumberCapability?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static NumberCapability? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NumberCapability._(super.value, super.name);
}

class RouteType extends $pb.ProtobufEnum {
  static const RouteType ROUTE_TYPE_TRANSACTIONAL =
      RouteType._(0, _omitEnumNames ? '' : 'ROUTE_TYPE_TRANSACTIONAL');
  static const RouteType ROUTE_TYPE_PREMIUM =
      RouteType._(1, _omitEnumNames ? '' : 'ROUTE_TYPE_PREMIUM');
  static const RouteType ROUTE_TYPE_PROMOTIONAL =
      RouteType._(2, _omitEnumNames ? '' : 'ROUTE_TYPE_PROMOTIONAL');

  static const $core.List<RouteType> values = <RouteType>[
    ROUTE_TYPE_TRANSACTIONAL,
    ROUTE_TYPE_PREMIUM,
    ROUTE_TYPE_PROMOTIONAL,
  ];

  static final $core.List<RouteType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static RouteType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RouteType._(super.value, super.name);
}

class SMSSandboxPhoneNumberVerificationStatus extends $pb.ProtobufEnum {
  static const SMSSandboxPhoneNumberVerificationStatus
      S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_VERIFIED =
      SMSSandboxPhoneNumberVerificationStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_VERIFIED');
  static const SMSSandboxPhoneNumberVerificationStatus
      S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_PENDING =
      SMSSandboxPhoneNumberVerificationStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_PENDING');

  static const $core.List<SMSSandboxPhoneNumberVerificationStatus> values =
      <SMSSandboxPhoneNumberVerificationStatus>[
    S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_VERIFIED,
    S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_PENDING,
  ];

  static final $core.List<SMSSandboxPhoneNumberVerificationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SMSSandboxPhoneNumberVerificationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SMSSandboxPhoneNumberVerificationStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
