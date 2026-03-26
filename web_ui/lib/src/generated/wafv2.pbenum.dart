// This is a generated file - do not edit.
//
// Generated from wafv2.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ActionValue extends $pb.ProtobufEnum {
  static const ActionValue ACTION_VALUE_CHALLENGE =
      ActionValue._(0, _omitEnumNames ? '' : 'ACTION_VALUE_CHALLENGE');
  static const ActionValue ACTION_VALUE_COUNT =
      ActionValue._(1, _omitEnumNames ? '' : 'ACTION_VALUE_COUNT');
  static const ActionValue ACTION_VALUE_BLOCK =
      ActionValue._(2, _omitEnumNames ? '' : 'ACTION_VALUE_BLOCK');
  static const ActionValue ACTION_VALUE_ALLOW =
      ActionValue._(3, _omitEnumNames ? '' : 'ACTION_VALUE_ALLOW');
  static const ActionValue ACTION_VALUE_EXCLUDED_AS_COUNT =
      ActionValue._(4, _omitEnumNames ? '' : 'ACTION_VALUE_EXCLUDED_AS_COUNT');
  static const ActionValue ACTION_VALUE_CAPTCHA =
      ActionValue._(5, _omitEnumNames ? '' : 'ACTION_VALUE_CAPTCHA');

  static const $core.List<ActionValue> values = <ActionValue>[
    ACTION_VALUE_CHALLENGE,
    ACTION_VALUE_COUNT,
    ACTION_VALUE_BLOCK,
    ACTION_VALUE_ALLOW,
    ACTION_VALUE_EXCLUDED_AS_COUNT,
    ACTION_VALUE_CAPTCHA,
  ];

  static final $core.List<ActionValue?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ActionValue? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ActionValue._(super.value, super.name);
}

class AssociatedResourceType extends $pb.ProtobufEnum {
  static const AssociatedResourceType
      ASSOCIATED_RESOURCE_TYPE_COGNITO_USER_POOL = AssociatedResourceType._(0,
          _omitEnumNames ? '' : 'ASSOCIATED_RESOURCE_TYPE_COGNITO_USER_POOL');
  static const AssociatedResourceType
      ASSOCIATED_RESOURCE_TYPE_VERIFIED_ACCESS_INSTANCE =
      AssociatedResourceType._(
          1,
          _omitEnumNames
              ? ''
              : 'ASSOCIATED_RESOURCE_TYPE_VERIFIED_ACCESS_INSTANCE');
  static const AssociatedResourceType ASSOCIATED_RESOURCE_TYPE_CLOUDFRONT =
      AssociatedResourceType._(
          2, _omitEnumNames ? '' : 'ASSOCIATED_RESOURCE_TYPE_CLOUDFRONT');
  static const AssociatedResourceType ASSOCIATED_RESOURCE_TYPE_API_GATEWAY =
      AssociatedResourceType._(
          3, _omitEnumNames ? '' : 'ASSOCIATED_RESOURCE_TYPE_API_GATEWAY');
  static const AssociatedResourceType
      ASSOCIATED_RESOURCE_TYPE_APP_RUNNER_SERVICE = AssociatedResourceType._(4,
          _omitEnumNames ? '' : 'ASSOCIATED_RESOURCE_TYPE_APP_RUNNER_SERVICE');

  static const $core.List<AssociatedResourceType> values =
      <AssociatedResourceType>[
    ASSOCIATED_RESOURCE_TYPE_COGNITO_USER_POOL,
    ASSOCIATED_RESOURCE_TYPE_VERIFIED_ACCESS_INSTANCE,
    ASSOCIATED_RESOURCE_TYPE_CLOUDFRONT,
    ASSOCIATED_RESOURCE_TYPE_API_GATEWAY,
    ASSOCIATED_RESOURCE_TYPE_APP_RUNNER_SERVICE,
  ];

  static final $core.List<AssociatedResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static AssociatedResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociatedResourceType._(super.value, super.name);
}

class BodyParsingFallbackBehavior extends $pb.ProtobufEnum {
  static const BodyParsingFallbackBehavior
      BODY_PARSING_FALLBACK_BEHAVIOR_EVALUATE_AS_STRING =
      BodyParsingFallbackBehavior._(
          0,
          _omitEnumNames
              ? ''
              : 'BODY_PARSING_FALLBACK_BEHAVIOR_EVALUATE_AS_STRING');
  static const BodyParsingFallbackBehavior
      BODY_PARSING_FALLBACK_BEHAVIOR_MATCH = BodyParsingFallbackBehavior._(
          1, _omitEnumNames ? '' : 'BODY_PARSING_FALLBACK_BEHAVIOR_MATCH');
  static const BodyParsingFallbackBehavior
      BODY_PARSING_FALLBACK_BEHAVIOR_NO_MATCH = BodyParsingFallbackBehavior._(
          2, _omitEnumNames ? '' : 'BODY_PARSING_FALLBACK_BEHAVIOR_NO_MATCH');

  static const $core.List<BodyParsingFallbackBehavior> values =
      <BodyParsingFallbackBehavior>[
    BODY_PARSING_FALLBACK_BEHAVIOR_EVALUATE_AS_STRING,
    BODY_PARSING_FALLBACK_BEHAVIOR_MATCH,
    BODY_PARSING_FALLBACK_BEHAVIOR_NO_MATCH,
  ];

  static final $core.List<BodyParsingFallbackBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static BodyParsingFallbackBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BodyParsingFallbackBehavior._(super.value, super.name);
}

class ComparisonOperator extends $pb.ProtobufEnum {
  static const ComparisonOperator COMPARISON_OPERATOR_GE =
      ComparisonOperator._(0, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_GE');
  static const ComparisonOperator COMPARISON_OPERATOR_EQ =
      ComparisonOperator._(1, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_EQ');
  static const ComparisonOperator COMPARISON_OPERATOR_LT =
      ComparisonOperator._(2, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_LT');
  static const ComparisonOperator COMPARISON_OPERATOR_GT =
      ComparisonOperator._(3, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_GT');
  static const ComparisonOperator COMPARISON_OPERATOR_LE =
      ComparisonOperator._(4, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_LE');
  static const ComparisonOperator COMPARISON_OPERATOR_NE =
      ComparisonOperator._(5, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_NE');

  static const $core.List<ComparisonOperator> values = <ComparisonOperator>[
    COMPARISON_OPERATOR_GE,
    COMPARISON_OPERATOR_EQ,
    COMPARISON_OPERATOR_LT,
    COMPARISON_OPERATOR_GT,
    COMPARISON_OPERATOR_LE,
    COMPARISON_OPERATOR_NE,
  ];

  static final $core.List<ComparisonOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ComparisonOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComparisonOperator._(super.value, super.name);
}

class CountryCode extends $pb.ProtobufEnum {
  static const CountryCode COUNTRY_CODE_US =
      CountryCode._(0, _omitEnumNames ? '' : 'COUNTRY_CODE_US');
  static const CountryCode COUNTRY_CODE_YE =
      CountryCode._(1, _omitEnumNames ? '' : 'COUNTRY_CODE_YE');
  static const CountryCode COUNTRY_CODE_TO =
      CountryCode._(2, _omitEnumNames ? '' : 'COUNTRY_CODE_TO');
  static const CountryCode COUNTRY_CODE_SI =
      CountryCode._(3, _omitEnumNames ? '' : 'COUNTRY_CODE_SI');
  static const CountryCode COUNTRY_CODE_MA =
      CountryCode._(4, _omitEnumNames ? '' : 'COUNTRY_CODE_MA');
  static const CountryCode COUNTRY_CODE_BY =
      CountryCode._(5, _omitEnumNames ? '' : 'COUNTRY_CODE_BY');
  static const CountryCode COUNTRY_CODE_PK =
      CountryCode._(6, _omitEnumNames ? '' : 'COUNTRY_CODE_PK');
  static const CountryCode COUNTRY_CODE_BB =
      CountryCode._(7, _omitEnumNames ? '' : 'COUNTRY_CODE_BB');
  static const CountryCode COUNTRY_CODE_CO =
      CountryCode._(8, _omitEnumNames ? '' : 'COUNTRY_CODE_CO');
  static const CountryCode COUNTRY_CODE_CZ =
      CountryCode._(9, _omitEnumNames ? '' : 'COUNTRY_CODE_CZ');
  static const CountryCode COUNTRY_CODE_FI =
      CountryCode._(10, _omitEnumNames ? '' : 'COUNTRY_CODE_FI');
  static const CountryCode COUNTRY_CODE_VA =
      CountryCode._(11, _omitEnumNames ? '' : 'COUNTRY_CODE_VA');
  static const CountryCode COUNTRY_CODE_GU =
      CountryCode._(12, _omitEnumNames ? '' : 'COUNTRY_CODE_GU');
  static const CountryCode COUNTRY_CODE_AU =
      CountryCode._(13, _omitEnumNames ? '' : 'COUNTRY_CODE_AU');
  static const CountryCode COUNTRY_CODE_TZ =
      CountryCode._(14, _omitEnumNames ? '' : 'COUNTRY_CODE_TZ');
  static const CountryCode COUNTRY_CODE_TG =
      CountryCode._(15, _omitEnumNames ? '' : 'COUNTRY_CODE_TG');
  static const CountryCode COUNTRY_CODE_IT =
      CountryCode._(16, _omitEnumNames ? '' : 'COUNTRY_CODE_IT');
  static const CountryCode COUNTRY_CODE_LS =
      CountryCode._(17, _omitEnumNames ? '' : 'COUNTRY_CODE_LS');
  static const CountryCode COUNTRY_CODE_RU =
      CountryCode._(18, _omitEnumNames ? '' : 'COUNTRY_CODE_RU');
  static const CountryCode COUNTRY_CODE_SZ =
      CountryCode._(19, _omitEnumNames ? '' : 'COUNTRY_CODE_SZ');
  static const CountryCode COUNTRY_CODE_SA =
      CountryCode._(20, _omitEnumNames ? '' : 'COUNTRY_CODE_SA');
  static const CountryCode COUNTRY_CODE_NO =
      CountryCode._(21, _omitEnumNames ? '' : 'COUNTRY_CODE_NO');
  static const CountryCode COUNTRY_CODE_BQ =
      CountryCode._(22, _omitEnumNames ? '' : 'COUNTRY_CODE_BQ');
  static const CountryCode COUNTRY_CODE_CG =
      CountryCode._(23, _omitEnumNames ? '' : 'COUNTRY_CODE_CG');
  static const CountryCode COUNTRY_CODE_UM =
      CountryCode._(24, _omitEnumNames ? '' : 'COUNTRY_CODE_UM');
  static const CountryCode COUNTRY_CODE_CR =
      CountryCode._(25, _omitEnumNames ? '' : 'COUNTRY_CODE_CR');
  static const CountryCode COUNTRY_CODE_VI =
      CountryCode._(26, _omitEnumNames ? '' : 'COUNTRY_CODE_VI');
  static const CountryCode COUNTRY_CODE_GM =
      CountryCode._(27, _omitEnumNames ? '' : 'COUNTRY_CODE_GM');
  static const CountryCode COUNTRY_CODE_DJ =
      CountryCode._(28, _omitEnumNames ? '' : 'COUNTRY_CODE_DJ');
  static const CountryCode COUNTRY_CODE_TR =
      CountryCode._(29, _omitEnumNames ? '' : 'COUNTRY_CODE_TR');
  static const CountryCode COUNTRY_CODE_KN =
      CountryCode._(30, _omitEnumNames ? '' : 'COUNTRY_CODE_KN');
  static const CountryCode COUNTRY_CODE_EG =
      CountryCode._(31, _omitEnumNames ? '' : 'COUNTRY_CODE_EG');
  static const CountryCode COUNTRY_CODE_HM =
      CountryCode._(32, _omitEnumNames ? '' : 'COUNTRY_CODE_HM');
  static const CountryCode COUNTRY_CODE_HR =
      CountryCode._(33, _omitEnumNames ? '' : 'COUNTRY_CODE_HR');
  static const CountryCode COUNTRY_CODE_IL =
      CountryCode._(34, _omitEnumNames ? '' : 'COUNTRY_CODE_IL');
  static const CountryCode COUNTRY_CODE_SR =
      CountryCode._(35, _omitEnumNames ? '' : 'COUNTRY_CODE_SR');
  static const CountryCode COUNTRY_CODE_MF =
      CountryCode._(36, _omitEnumNames ? '' : 'COUNTRY_CODE_MF');
  static const CountryCode COUNTRY_CODE_NG =
      CountryCode._(37, _omitEnumNames ? '' : 'COUNTRY_CODE_NG');
  static const CountryCode COUNTRY_CODE_MQ =
      CountryCode._(38, _omitEnumNames ? '' : 'COUNTRY_CODE_MQ');
  static const CountryCode COUNTRY_CODE_NP =
      CountryCode._(39, _omitEnumNames ? '' : 'COUNTRY_CODE_NP');
  static const CountryCode COUNTRY_CODE_PF =
      CountryCode._(40, _omitEnumNames ? '' : 'COUNTRY_CODE_PF');
  static const CountryCode COUNTRY_CODE_BI =
      CountryCode._(41, _omitEnumNames ? '' : 'COUNTRY_CODE_BI');
  static const CountryCode COUNTRY_CODE_GE =
      CountryCode._(42, _omitEnumNames ? '' : 'COUNTRY_CODE_GE');
  static const CountryCode COUNTRY_CODE_AR =
      CountryCode._(43, _omitEnumNames ? '' : 'COUNTRY_CODE_AR');
  static const CountryCode COUNTRY_CODE_AE =
      CountryCode._(44, _omitEnumNames ? '' : 'COUNTRY_CODE_AE');
  static const CountryCode COUNTRY_CODE_TJ =
      CountryCode._(45, _omitEnumNames ? '' : 'COUNTRY_CODE_TJ');
  static const CountryCode COUNTRY_CODE_ET =
      CountryCode._(46, _omitEnumNames ? '' : 'COUNTRY_CODE_ET');
  static const CountryCode COUNTRY_CODE_RO =
      CountryCode._(47, _omitEnumNames ? '' : 'COUNTRY_CODE_RO');
  static const CountryCode COUNTRY_CODE_IQ =
      CountryCode._(48, _omitEnumNames ? '' : 'COUNTRY_CODE_IQ');
  static const CountryCode COUNTRY_CODE_LV =
      CountryCode._(49, _omitEnumNames ? '' : 'COUNTRY_CODE_LV');
  static const CountryCode COUNTRY_CODE_ID =
      CountryCode._(50, _omitEnumNames ? '' : 'COUNTRY_CODE_ID');
  static const CountryCode COUNTRY_CODE_LC =
      CountryCode._(51, _omitEnumNames ? '' : 'COUNTRY_CODE_LC');
  static const CountryCode COUNTRY_CODE_SJ =
      CountryCode._(52, _omitEnumNames ? '' : 'COUNTRY_CODE_SJ');
  static const CountryCode COUNTRY_CODE_MN =
      CountryCode._(53, _omitEnumNames ? '' : 'COUNTRY_CODE_MN');
  static const CountryCode COUNTRY_CODE_PY =
      CountryCode._(54, _omitEnumNames ? '' : 'COUNTRY_CODE_PY');
  static const CountryCode COUNTRY_CODE_MY =
      CountryCode._(55, _omitEnumNames ? '' : 'COUNTRY_CODE_MY');
  static const CountryCode COUNTRY_CODE_PN =
      CountryCode._(56, _omitEnumNames ? '' : 'COUNTRY_CODE_PN');
  static const CountryCode COUNTRY_CODE_BA =
      CountryCode._(57, _omitEnumNames ? '' : 'COUNTRY_CODE_BA');
  static const CountryCode COUNTRY_CODE_CH =
      CountryCode._(58, _omitEnumNames ? '' : 'COUNTRY_CODE_CH');
  static const CountryCode COUNTRY_CODE_CW =
      CountryCode._(59, _omitEnumNames ? '' : 'COUNTRY_CODE_CW');
  static const CountryCode COUNTRY_CODE_AZ =
      CountryCode._(60, _omitEnumNames ? '' : 'COUNTRY_CODE_AZ');
  static const CountryCode COUNTRY_CODE_DM =
      CountryCode._(61, _omitEnumNames ? '' : 'COUNTRY_CODE_DM');
  static const CountryCode COUNTRY_CODE_AM =
      CountryCode._(62, _omitEnumNames ? '' : 'COUNTRY_CODE_AM');
  static const CountryCode COUNTRY_CODE_DZ =
      CountryCode._(63, _omitEnumNames ? '' : 'COUNTRY_CODE_DZ');
  static const CountryCode COUNTRY_CODE_LK =
      CountryCode._(64, _omitEnumNames ? '' : 'COUNTRY_CODE_LK');
  static const CountryCode COUNTRY_CODE_SB =
      CountryCode._(65, _omitEnumNames ? '' : 'COUNTRY_CODE_SB');
  static const CountryCode COUNTRY_CODE_MV =
      CountryCode._(66, _omitEnumNames ? '' : 'COUNTRY_CODE_MV');
  static const CountryCode COUNTRY_CODE_PA =
      CountryCode._(67, _omitEnumNames ? '' : 'COUNTRY_CODE_PA');
  static const CountryCode COUNTRY_CODE_BT =
      CountryCode._(68, _omitEnumNames ? '' : 'COUNTRY_CODE_BT');
  static const CountryCode COUNTRY_CODE_GN =
      CountryCode._(69, _omitEnumNames ? '' : 'COUNTRY_CODE_GN');
  static const CountryCode COUNTRY_CODE_DE =
      CountryCode._(70, _omitEnumNames ? '' : 'COUNTRY_CODE_DE');
  static const CountryCode COUNTRY_CODE_TM =
      CountryCode._(71, _omitEnumNames ? '' : 'COUNTRY_CODE_TM');
  static const CountryCode COUNTRY_CODE_HU =
      CountryCode._(72, _omitEnumNames ? '' : 'COUNTRY_CODE_HU');
  static const CountryCode COUNTRY_CODE_SO =
      CountryCode._(73, _omitEnumNames ? '' : 'COUNTRY_CODE_SO');
  static const CountryCode COUNTRY_CODE_MC =
      CountryCode._(74, _omitEnumNames ? '' : 'COUNTRY_CODE_MC');
  static const CountryCode COUNTRY_CODE_PT =
      CountryCode._(75, _omitEnumNames ? '' : 'COUNTRY_CODE_PT');
  static const CountryCode COUNTRY_CODE_BL =
      CountryCode._(76, _omitEnumNames ? '' : 'COUNTRY_CODE_BL');
  static const CountryCode COUNTRY_CODE_CM =
      CountryCode._(77, _omitEnumNames ? '' : 'COUNTRY_CODE_CM');
  static const CountryCode COUNTRY_CODE_CX =
      CountryCode._(78, _omitEnumNames ? '' : 'COUNTRY_CODE_CX');
  static const CountryCode COUNTRY_CODE_FK =
      CountryCode._(79, _omitEnumNames ? '' : 'COUNTRY_CODE_FK');
  static const CountryCode COUNTRY_CODE_VC =
      CountryCode._(80, _omitEnumNames ? '' : 'COUNTRY_CODE_VC');
  static const CountryCode COUNTRY_CODE_GS =
      CountryCode._(81, _omitEnumNames ? '' : 'COUNTRY_CODE_GS');
  static const CountryCode COUNTRY_CODE_GF =
      CountryCode._(82, _omitEnumNames ? '' : 'COUNTRY_CODE_GF');
  static const CountryCode COUNTRY_CODE_AW =
      CountryCode._(83, _omitEnumNames ? '' : 'COUNTRY_CODE_AW');
  static const CountryCode COUNTRY_CODE_RW =
      CountryCode._(84, _omitEnumNames ? '' : 'COUNTRY_CODE_RW');
  static const CountryCode COUNTRY_CODE_SX =
      CountryCode._(85, _omitEnumNames ? '' : 'COUNTRY_CODE_SX');
  static const CountryCode COUNTRY_CODE_SG =
      CountryCode._(86, _omitEnumNames ? '' : 'COUNTRY_CODE_SG');
  static const CountryCode COUNTRY_CODE_MK =
      CountryCode._(87, _omitEnumNames ? '' : 'COUNTRY_CODE_MK');
  static const CountryCode COUNTRY_CODE_NZ =
      CountryCode._(88, _omitEnumNames ? '' : 'COUNTRY_CODE_NZ');
  static const CountryCode COUNTRY_CODE_BS =
      CountryCode._(89, _omitEnumNames ? '' : 'COUNTRY_CODE_BS');
  static const CountryCode COUNTRY_CODE_BD =
      CountryCode._(90, _omitEnumNames ? '' : 'COUNTRY_CODE_BD');
  static const CountryCode COUNTRY_CODE_UZ =
      CountryCode._(91, _omitEnumNames ? '' : 'COUNTRY_CODE_UZ');
  static const CountryCode COUNTRY_CODE_JE =
      CountryCode._(92, _omitEnumNames ? '' : 'COUNTRY_CODE_JE');
  static const CountryCode COUNTRY_CODE_WF =
      CountryCode._(93, _omitEnumNames ? '' : 'COUNTRY_CODE_WF');
  static const CountryCode COUNTRY_CODE_HK =
      CountryCode._(94, _omitEnumNames ? '' : 'COUNTRY_CODE_HK');
  static const CountryCode COUNTRY_CODE_IN =
      CountryCode._(95, _omitEnumNames ? '' : 'COUNTRY_CODE_IN');
  static const CountryCode COUNTRY_CODE_LY =
      CountryCode._(96, _omitEnumNames ? '' : 'COUNTRY_CODE_LY');
  static const CountryCode COUNTRY_CODE_NI =
      CountryCode._(97, _omitEnumNames ? '' : 'COUNTRY_CODE_NI');
  static const CountryCode COUNTRY_CODE_MS =
      CountryCode._(98, _omitEnumNames ? '' : 'COUNTRY_CODE_MS');
  static const CountryCode COUNTRY_CODE_NR =
      CountryCode._(99, _omitEnumNames ? '' : 'COUNTRY_CODE_NR');
  static const CountryCode COUNTRY_CODE_UG =
      CountryCode._(100, _omitEnumNames ? '' : 'COUNTRY_CODE_UG');
  static const CountryCode COUNTRY_CODE_JM =
      CountryCode._(101, _omitEnumNames ? '' : 'COUNTRY_CODE_JM');
  static const CountryCode COUNTRY_CODE_AG =
      CountryCode._(102, _omitEnumNames ? '' : 'COUNTRY_CODE_AG');
  static const CountryCode COUNTRY_CODE_KY =
      CountryCode._(103, _omitEnumNames ? '' : 'COUNTRY_CODE_KY');
  static const CountryCode COUNTRY_CODE_TH =
      CountryCode._(104, _omitEnumNames ? '' : 'COUNTRY_CODE_TH');
  static const CountryCode COUNTRY_CODE_IS =
      CountryCode._(105, _omitEnumNames ? '' : 'COUNTRY_CODE_IS');
  static const CountryCode COUNTRY_CODE_LT =
      CountryCode._(106, _omitEnumNames ? '' : 'COUNTRY_CODE_LT');
  static const CountryCode COUNTRY_CODE_LA =
      CountryCode._(107, _omitEnumNames ? '' : 'COUNTRY_CODE_LA');
  static const CountryCode COUNTRY_CODE_SH =
      CountryCode._(108, _omitEnumNames ? '' : 'COUNTRY_CODE_SH');
  static const CountryCode COUNTRY_CODE_NA =
      CountryCode._(109, _omitEnumNames ? '' : 'COUNTRY_CODE_NA');
  static const CountryCode COUNTRY_CODE_PW =
      CountryCode._(110, _omitEnumNames ? '' : 'COUNTRY_CODE_PW');
  static const CountryCode COUNTRY_CODE_PL =
      CountryCode._(111, _omitEnumNames ? '' : 'COUNTRY_CODE_PL');
  static const CountryCode COUNTRY_CODE_CN =
      CountryCode._(112, _omitEnumNames ? '' : 'COUNTRY_CODE_CN');
  static const CountryCode COUNTRY_CODE_VU =
      CountryCode._(113, _omitEnumNames ? '' : 'COUNTRY_CODE_VU');
  static const CountryCode COUNTRY_CODE_CU =
      CountryCode._(114, _omitEnumNames ? '' : 'COUNTRY_CODE_CU');
  static const CountryCode COUNTRY_CODE_VN =
      CountryCode._(115, _omitEnumNames ? '' : 'COUNTRY_CODE_VN');
  static const CountryCode COUNTRY_CODE_GT =
      CountryCode._(116, _omitEnumNames ? '' : 'COUNTRY_CODE_GT');
  static const CountryCode COUNTRY_CODE_AT =
      CountryCode._(117, _omitEnumNames ? '' : 'COUNTRY_CODE_AT');
  static const CountryCode COUNTRY_CODE_ZM =
      CountryCode._(118, _omitEnumNames ? '' : 'COUNTRY_CODE_ZM');
  static const CountryCode COUNTRY_CODE_AO =
      CountryCode._(119, _omitEnumNames ? '' : 'COUNTRY_CODE_AO');
  static const CountryCode COUNTRY_CODE_LI =
      CountryCode._(120, _omitEnumNames ? '' : 'COUNTRY_CODE_LI');
  static const CountryCode COUNTRY_CODE_NL =
      CountryCode._(121, _omitEnumNames ? '' : 'COUNTRY_CODE_NL');
  static const CountryCode COUNTRY_CODE_MH =
      CountryCode._(122, _omitEnumNames ? '' : 'COUNTRY_CODE_MH');
  static const CountryCode COUNTRY_CODE_BV =
      CountryCode._(123, _omitEnumNames ? '' : 'COUNTRY_CODE_BV');
  static const CountryCode COUNTRY_CODE_CF =
      CountryCode._(124, _omitEnumNames ? '' : 'COUNTRY_CODE_CF');
  static const CountryCode COUNTRY_CODE_GY =
      CountryCode._(125, _omitEnumNames ? '' : 'COUNTRY_CODE_GY');
  static const CountryCode COUNTRY_CODE_GL =
      CountryCode._(126, _omitEnumNames ? '' : 'COUNTRY_CODE_GL');
  static const CountryCode COUNTRY_CODE_DK =
      CountryCode._(127, _omitEnumNames ? '' : 'COUNTRY_CODE_DK');
  static const CountryCode COUNTRY_CODE_KR =
      CountryCode._(128, _omitEnumNames ? '' : 'COUNTRY_CODE_KR');
  static const CountryCode COUNTRY_CODE_YT =
      CountryCode._(129, _omitEnumNames ? '' : 'COUNTRY_CODE_YT');
  static const CountryCode COUNTRY_CODE_KI =
      CountryCode._(130, _omitEnumNames ? '' : 'COUNTRY_CODE_KI');
  static const CountryCode COUNTRY_CODE_HN =
      CountryCode._(131, _omitEnumNames ? '' : 'COUNTRY_CODE_HN');
  static const CountryCode COUNTRY_CODE_SM =
      CountryCode._(132, _omitEnumNames ? '' : 'COUNTRY_CODE_SM');
  static const CountryCode COUNTRY_CODE_ME =
      CountryCode._(133, _omitEnumNames ? '' : 'COUNTRY_CODE_ME');
  static const CountryCode COUNTRY_CODE_PR =
      CountryCode._(134, _omitEnumNames ? '' : 'COUNTRY_CODE_PR');
  static const CountryCode COUNTRY_CODE_MP =
      CountryCode._(135, _omitEnumNames ? '' : 'COUNTRY_CODE_MP');
  static const CountryCode COUNTRY_CODE_PG =
      CountryCode._(136, _omitEnumNames ? '' : 'COUNTRY_CODE_PG');
  static const CountryCode COUNTRY_CODE_BN =
      CountryCode._(137, _omitEnumNames ? '' : 'COUNTRY_CODE_BN');
  static const CountryCode COUNTRY_CODE_FM =
      CountryCode._(138, _omitEnumNames ? '' : 'COUNTRY_CODE_FM');
  static const CountryCode COUNTRY_CODE_VE =
      CountryCode._(139, _omitEnumNames ? '' : 'COUNTRY_CODE_VE');
  static const CountryCode COUNTRY_CODE_GQ =
      CountryCode._(140, _omitEnumNames ? '' : 'COUNTRY_CODE_GQ');
  static const CountryCode COUNTRY_CODE_GD =
      CountryCode._(141, _omitEnumNames ? '' : 'COUNTRY_CODE_GD');
  static const CountryCode COUNTRY_CODE_ZW =
      CountryCode._(142, _omitEnumNames ? '' : 'COUNTRY_CODE_ZW');
  static const CountryCode COUNTRY_CODE_AQ =
      CountryCode._(143, _omitEnumNames ? '' : 'COUNTRY_CODE_AQ');
  static const CountryCode COUNTRY_CODE_JP =
      CountryCode._(144, _omitEnumNames ? '' : 'COUNTRY_CODE_JP');
  static const CountryCode COUNTRY_CODE_WS =
      CountryCode._(145, _omitEnumNames ? '' : 'COUNTRY_CODE_WS');
  static const CountryCode COUNTRY_CODE_AD =
      CountryCode._(146, _omitEnumNames ? '' : 'COUNTRY_CODE_AD');
  static const CountryCode COUNTRY_CODE_KZ =
      CountryCode._(147, _omitEnumNames ? '' : 'COUNTRY_CODE_KZ');
  static const CountryCode COUNTRY_CODE_TK =
      CountryCode._(148, _omitEnumNames ? '' : 'COUNTRY_CODE_TK');
  static const CountryCode COUNTRY_CODE_ES =
      CountryCode._(149, _omitEnumNames ? '' : 'COUNTRY_CODE_ES');
  static const CountryCode COUNTRY_CODE_SE =
      CountryCode._(150, _omitEnumNames ? '' : 'COUNTRY_CODE_SE');
  static const CountryCode COUNTRY_CODE_MM =
      CountryCode._(151, _omitEnumNames ? '' : 'COUNTRY_CODE_MM');
  static const CountryCode COUNTRY_CODE_MX =
      CountryCode._(152, _omitEnumNames ? '' : 'COUNTRY_CODE_MX');
  static const CountryCode COUNTRY_CODE_BF =
      CountryCode._(153, _omitEnumNames ? '' : 'COUNTRY_CODE_BF');
  static const CountryCode COUNTRY_CODE_CK =
      CountryCode._(154, _omitEnumNames ? '' : 'COUNTRY_CODE_CK');
  static const CountryCode COUNTRY_CODE_CV =
      CountryCode._(155, _omitEnumNames ? '' : 'COUNTRY_CODE_CV');
  static const CountryCode COUNTRY_CODE_XK =
      CountryCode._(156, _omitEnumNames ? '' : 'COUNTRY_CODE_XK');
  static const CountryCode COUNTRY_CODE_GI =
      CountryCode._(157, _omitEnumNames ? '' : 'COUNTRY_CODE_GI');
  static const CountryCode COUNTRY_CODE_KW =
      CountryCode._(158, _omitEnumNames ? '' : 'COUNTRY_CODE_KW');
  static const CountryCode COUNTRY_CODE_TV =
      CountryCode._(159, _omitEnumNames ? '' : 'COUNTRY_CODE_TV');
  static const CountryCode COUNTRY_CODE_AL =
      CountryCode._(160, _omitEnumNames ? '' : 'COUNTRY_CODE_AL');
  static const CountryCode COUNTRY_CODE_TC =
      CountryCode._(161, _omitEnumNames ? '' : 'COUNTRY_CODE_TC');
  static const CountryCode COUNTRY_CODE_SV =
      CountryCode._(162, _omitEnumNames ? '' : 'COUNTRY_CODE_SV');
  static const CountryCode COUNTRY_CODE_MU =
      CountryCode._(163, _omitEnumNames ? '' : 'COUNTRY_CODE_MU');
  static const CountryCode COUNTRY_CODE_CC =
      CountryCode._(164, _omitEnumNames ? '' : 'COUNTRY_CODE_CC');
  static const CountryCode COUNTRY_CODE_GA =
      CountryCode._(165, _omitEnumNames ? '' : 'COUNTRY_CODE_GA');
  static const CountryCode COUNTRY_CODE_JO =
      CountryCode._(166, _omitEnumNames ? '' : 'COUNTRY_CODE_JO');
  static const CountryCode COUNTRY_CODE_TN =
      CountryCode._(167, _omitEnumNames ? '' : 'COUNTRY_CODE_TN');
  static const CountryCode COUNTRY_CODE_EC =
      CountryCode._(168, _omitEnumNames ? '' : 'COUNTRY_CODE_EC');
  static const CountryCode COUNTRY_CODE_SN =
      CountryCode._(169, _omitEnumNames ? '' : 'COUNTRY_CODE_SN');
  static const CountryCode COUNTRY_CODE_NC =
      CountryCode._(170, _omitEnumNames ? '' : 'COUNTRY_CODE_NC');
  static const CountryCode COUNTRY_CODE_BM =
      CountryCode._(171, _omitEnumNames ? '' : 'COUNTRY_CODE_BM');
  static const CountryCode COUNTRY_CODE_CL =
      CountryCode._(172, _omitEnumNames ? '' : 'COUNTRY_CODE_CL');
  static const CountryCode COUNTRY_CODE_UA =
      CountryCode._(173, _omitEnumNames ? '' : 'COUNTRY_CODE_UA');
  static const CountryCode COUNTRY_CODE_GR =
      CountryCode._(174, _omitEnumNames ? '' : 'COUNTRY_CODE_GR');
  static const CountryCode COUNTRY_CODE_AI =
      CountryCode._(175, _omitEnumNames ? '' : 'COUNTRY_CODE_AI');
  static const CountryCode COUNTRY_CODE_KG =
      CountryCode._(176, _omitEnumNames ? '' : 'COUNTRY_CODE_KG');
  static const CountryCode COUNTRY_CODE_TF =
      CountryCode._(177, _omitEnumNames ? '' : 'COUNTRY_CODE_TF');
  static const CountryCode COUNTRY_CODE_LR =
      CountryCode._(178, _omitEnumNames ? '' : 'COUNTRY_CODE_LR');
  static const CountryCode COUNTRY_CODE_BE =
      CountryCode._(179, _omitEnumNames ? '' : 'COUNTRY_CODE_BE');
  static const CountryCode COUNTRY_CODE_CD =
      CountryCode._(180, _omitEnumNames ? '' : 'COUNTRY_CODE_CD');
  static const CountryCode COUNTRY_CODE_UY =
      CountryCode._(181, _omitEnumNames ? '' : 'COUNTRY_CODE_UY');
  static const CountryCode COUNTRY_CODE_KP =
      CountryCode._(182, _omitEnumNames ? '' : 'COUNTRY_CODE_KP');
  static const CountryCode COUNTRY_CODE_EH =
      CountryCode._(183, _omitEnumNames ? '' : 'COUNTRY_CODE_EH');
  static const CountryCode COUNTRY_CODE_IM =
      CountryCode._(184, _omitEnumNames ? '' : 'COUNTRY_CODE_IM');
  static const CountryCode COUNTRY_CODE_SS =
      CountryCode._(185, _omitEnumNames ? '' : 'COUNTRY_CODE_SS');
  static const CountryCode COUNTRY_CODE_MG =
      CountryCode._(186, _omitEnumNames ? '' : 'COUNTRY_CODE_MG');
  static const CountryCode COUNTRY_CODE_NF =
      CountryCode._(187, _omitEnumNames ? '' : 'COUNTRY_CODE_NF');
  static const CountryCode COUNTRY_CODE_MR =
      CountryCode._(188, _omitEnumNames ? '' : 'COUNTRY_CODE_MR');
  static const CountryCode COUNTRY_CODE_PE =
      CountryCode._(189, _omitEnumNames ? '' : 'COUNTRY_CODE_PE');
  static const CountryCode COUNTRY_CODE_BH =
      CountryCode._(190, _omitEnumNames ? '' : 'COUNTRY_CODE_BH');
  static const CountryCode COUNTRY_CODE_FO =
      CountryCode._(191, _omitEnumNames ? '' : 'COUNTRY_CODE_FO');
  static const CountryCode COUNTRY_CODE_VG =
      CountryCode._(192, _omitEnumNames ? '' : 'COUNTRY_CODE_VG');
  static const CountryCode COUNTRY_CODE_GB =
      CountryCode._(193, _omitEnumNames ? '' : 'COUNTRY_CODE_GB');
  static const CountryCode COUNTRY_CODE_AS =
      CountryCode._(194, _omitEnumNames ? '' : 'COUNTRY_CODE_AS');
  static const CountryCode COUNTRY_CODE_AF =
      CountryCode._(195, _omitEnumNames ? '' : 'COUNTRY_CODE_AF');
  static const CountryCode COUNTRY_CODE_IR =
      CountryCode._(196, _omitEnumNames ? '' : 'COUNTRY_CODE_IR');
  static const CountryCode COUNTRY_CODE_LU =
      CountryCode._(197, _omitEnumNames ? '' : 'COUNTRY_CODE_LU');
  static const CountryCode COUNTRY_CODE_RS =
      CountryCode._(198, _omitEnumNames ? '' : 'COUNTRY_CODE_RS');
  static const CountryCode COUNTRY_CODE_IE =
      CountryCode._(199, _omitEnumNames ? '' : 'COUNTRY_CODE_IE');
  static const CountryCode COUNTRY_CODE_LB =
      CountryCode._(200, _omitEnumNames ? '' : 'COUNTRY_CODE_LB');
  static const CountryCode COUNTRY_CODE_SK =
      CountryCode._(201, _omitEnumNames ? '' : 'COUNTRY_CODE_SK');
  static const CountryCode COUNTRY_CODE_MO =
      CountryCode._(202, _omitEnumNames ? '' : 'COUNTRY_CODE_MO');
  static const CountryCode COUNTRY_CODE_MZ =
      CountryCode._(203, _omitEnumNames ? '' : 'COUNTRY_CODE_MZ');
  static const CountryCode COUNTRY_CODE_PM =
      CountryCode._(204, _omitEnumNames ? '' : 'COUNTRY_CODE_PM');
  static const CountryCode COUNTRY_CODE_CI =
      CountryCode._(205, _omitEnumNames ? '' : 'COUNTRY_CODE_CI');
  static const CountryCode COUNTRY_CODE_GW =
      CountryCode._(206, _omitEnumNames ? '' : 'COUNTRY_CODE_GW');
  static const CountryCode COUNTRY_CODE_TT =
      CountryCode._(207, _omitEnumNames ? '' : 'COUNTRY_CODE_TT');
  static const CountryCode COUNTRY_CODE_OM =
      CountryCode._(208, _omitEnumNames ? '' : 'COUNTRY_CODE_OM');
  static const CountryCode COUNTRY_CODE_ST =
      CountryCode._(209, _omitEnumNames ? '' : 'COUNTRY_CODE_ST');
  static const CountryCode COUNTRY_CODE_SC =
      CountryCode._(210, _omitEnumNames ? '' : 'COUNTRY_CODE_SC');
  static const CountryCode COUNTRY_CODE_MW =
      CountryCode._(211, _omitEnumNames ? '' : 'COUNTRY_CODE_MW');
  static const CountryCode COUNTRY_CODE_BW =
      CountryCode._(212, _omitEnumNames ? '' : 'COUNTRY_CODE_BW');
  static const CountryCode COUNTRY_CODE_CA =
      CountryCode._(213, _omitEnumNames ? '' : 'COUNTRY_CODE_CA');
  static const CountryCode COUNTRY_CODE_ZA =
      CountryCode._(214, _omitEnumNames ? '' : 'COUNTRY_CODE_ZA');
  static const CountryCode COUNTRY_CODE_TL =
      CountryCode._(215, _omitEnumNames ? '' : 'COUNTRY_CODE_TL');
  static const CountryCode COUNTRY_CODE_KH =
      CountryCode._(216, _omitEnumNames ? '' : 'COUNTRY_CODE_KH');
  static const CountryCode COUNTRY_CODE_EE =
      CountryCode._(217, _omitEnumNames ? '' : 'COUNTRY_CODE_EE');
  static const CountryCode COUNTRY_CODE_HT =
      CountryCode._(218, _omitEnumNames ? '' : 'COUNTRY_CODE_HT');
  static const CountryCode COUNTRY_CODE_SL =
      CountryCode._(219, _omitEnumNames ? '' : 'COUNTRY_CODE_SL');
  static const CountryCode COUNTRY_CODE_MD =
      CountryCode._(220, _omitEnumNames ? '' : 'COUNTRY_CODE_MD');
  static const CountryCode COUNTRY_CODE_NE =
      CountryCode._(221, _omitEnumNames ? '' : 'COUNTRY_CODE_NE');
  static const CountryCode COUNTRY_CODE_PS =
      CountryCode._(222, _omitEnumNames ? '' : 'COUNTRY_CODE_PS');
  static const CountryCode COUNTRY_CODE_BZ =
      CountryCode._(223, _omitEnumNames ? '' : 'COUNTRY_CODE_BZ');
  static const CountryCode COUNTRY_CODE_PH =
      CountryCode._(224, _omitEnumNames ? '' : 'COUNTRY_CODE_PH');
  static const CountryCode COUNTRY_CODE_BO =
      CountryCode._(225, _omitEnumNames ? '' : 'COUNTRY_CODE_BO');
  static const CountryCode COUNTRY_CODE_QA =
      CountryCode._(226, _omitEnumNames ? '' : 'COUNTRY_CODE_QA');
  static const CountryCode COUNTRY_CODE_CY =
      CountryCode._(227, _omitEnumNames ? '' : 'COUNTRY_CODE_CY');
  static const CountryCode COUNTRY_CODE_FJ =
      CountryCode._(228, _omitEnumNames ? '' : 'COUNTRY_CODE_FJ');
  static const CountryCode COUNTRY_CODE_GP =
      CountryCode._(229, _omitEnumNames ? '' : 'COUNTRY_CODE_GP');
  static const CountryCode COUNTRY_CODE_GG =
      CountryCode._(230, _omitEnumNames ? '' : 'COUNTRY_CODE_GG');
  static const CountryCode COUNTRY_CODE_KE =
      CountryCode._(231, _omitEnumNames ? '' : 'COUNTRY_CODE_KE');
  static const CountryCode COUNTRY_CODE_TD =
      CountryCode._(232, _omitEnumNames ? '' : 'COUNTRY_CODE_TD');
  static const CountryCode COUNTRY_CODE_ER =
      CountryCode._(233, _omitEnumNames ? '' : 'COUNTRY_CODE_ER');
  static const CountryCode COUNTRY_CODE_SY =
      CountryCode._(234, _omitEnumNames ? '' : 'COUNTRY_CODE_SY');
  static const CountryCode COUNTRY_CODE_SD =
      CountryCode._(235, _omitEnumNames ? '' : 'COUNTRY_CODE_SD');
  static const CountryCode COUNTRY_CODE_ML =
      CountryCode._(236, _omitEnumNames ? '' : 'COUNTRY_CODE_ML');
  static const CountryCode COUNTRY_CODE_BR =
      CountryCode._(237, _omitEnumNames ? '' : 'COUNTRY_CODE_BR');
  static const CountryCode COUNTRY_CODE_BG =
      CountryCode._(238, _omitEnumNames ? '' : 'COUNTRY_CODE_BG');
  static const CountryCode COUNTRY_CODE_FR =
      CountryCode._(239, _omitEnumNames ? '' : 'COUNTRY_CODE_FR');
  static const CountryCode COUNTRY_CODE_GH =
      CountryCode._(240, _omitEnumNames ? '' : 'COUNTRY_CODE_GH');
  static const CountryCode COUNTRY_CODE_AX =
      CountryCode._(241, _omitEnumNames ? '' : 'COUNTRY_CODE_AX');
  static const CountryCode COUNTRY_CODE_DO =
      CountryCode._(242, _omitEnumNames ? '' : 'COUNTRY_CODE_DO');
  static const CountryCode COUNTRY_CODE_TW =
      CountryCode._(243, _omitEnumNames ? '' : 'COUNTRY_CODE_TW');
  static const CountryCode COUNTRY_CODE_KM =
      CountryCode._(244, _omitEnumNames ? '' : 'COUNTRY_CODE_KM');
  static const CountryCode COUNTRY_CODE_RE =
      CountryCode._(245, _omitEnumNames ? '' : 'COUNTRY_CODE_RE');
  static const CountryCode COUNTRY_CODE_IO =
      CountryCode._(246, _omitEnumNames ? '' : 'COUNTRY_CODE_IO');
  static const CountryCode COUNTRY_CODE_MT =
      CountryCode._(247, _omitEnumNames ? '' : 'COUNTRY_CODE_MT');
  static const CountryCode COUNTRY_CODE_NU =
      CountryCode._(248, _omitEnumNames ? '' : 'COUNTRY_CODE_NU');
  static const CountryCode COUNTRY_CODE_BJ =
      CountryCode._(249, _omitEnumNames ? '' : 'COUNTRY_CODE_BJ');

  static const $core.List<CountryCode> values = <CountryCode>[
    COUNTRY_CODE_US,
    COUNTRY_CODE_YE,
    COUNTRY_CODE_TO,
    COUNTRY_CODE_SI,
    COUNTRY_CODE_MA,
    COUNTRY_CODE_BY,
    COUNTRY_CODE_PK,
    COUNTRY_CODE_BB,
    COUNTRY_CODE_CO,
    COUNTRY_CODE_CZ,
    COUNTRY_CODE_FI,
    COUNTRY_CODE_VA,
    COUNTRY_CODE_GU,
    COUNTRY_CODE_AU,
    COUNTRY_CODE_TZ,
    COUNTRY_CODE_TG,
    COUNTRY_CODE_IT,
    COUNTRY_CODE_LS,
    COUNTRY_CODE_RU,
    COUNTRY_CODE_SZ,
    COUNTRY_CODE_SA,
    COUNTRY_CODE_NO,
    COUNTRY_CODE_BQ,
    COUNTRY_CODE_CG,
    COUNTRY_CODE_UM,
    COUNTRY_CODE_CR,
    COUNTRY_CODE_VI,
    COUNTRY_CODE_GM,
    COUNTRY_CODE_DJ,
    COUNTRY_CODE_TR,
    COUNTRY_CODE_KN,
    COUNTRY_CODE_EG,
    COUNTRY_CODE_HM,
    COUNTRY_CODE_HR,
    COUNTRY_CODE_IL,
    COUNTRY_CODE_SR,
    COUNTRY_CODE_MF,
    COUNTRY_CODE_NG,
    COUNTRY_CODE_MQ,
    COUNTRY_CODE_NP,
    COUNTRY_CODE_PF,
    COUNTRY_CODE_BI,
    COUNTRY_CODE_GE,
    COUNTRY_CODE_AR,
    COUNTRY_CODE_AE,
    COUNTRY_CODE_TJ,
    COUNTRY_CODE_ET,
    COUNTRY_CODE_RO,
    COUNTRY_CODE_IQ,
    COUNTRY_CODE_LV,
    COUNTRY_CODE_ID,
    COUNTRY_CODE_LC,
    COUNTRY_CODE_SJ,
    COUNTRY_CODE_MN,
    COUNTRY_CODE_PY,
    COUNTRY_CODE_MY,
    COUNTRY_CODE_PN,
    COUNTRY_CODE_BA,
    COUNTRY_CODE_CH,
    COUNTRY_CODE_CW,
    COUNTRY_CODE_AZ,
    COUNTRY_CODE_DM,
    COUNTRY_CODE_AM,
    COUNTRY_CODE_DZ,
    COUNTRY_CODE_LK,
    COUNTRY_CODE_SB,
    COUNTRY_CODE_MV,
    COUNTRY_CODE_PA,
    COUNTRY_CODE_BT,
    COUNTRY_CODE_GN,
    COUNTRY_CODE_DE,
    COUNTRY_CODE_TM,
    COUNTRY_CODE_HU,
    COUNTRY_CODE_SO,
    COUNTRY_CODE_MC,
    COUNTRY_CODE_PT,
    COUNTRY_CODE_BL,
    COUNTRY_CODE_CM,
    COUNTRY_CODE_CX,
    COUNTRY_CODE_FK,
    COUNTRY_CODE_VC,
    COUNTRY_CODE_GS,
    COUNTRY_CODE_GF,
    COUNTRY_CODE_AW,
    COUNTRY_CODE_RW,
    COUNTRY_CODE_SX,
    COUNTRY_CODE_SG,
    COUNTRY_CODE_MK,
    COUNTRY_CODE_NZ,
    COUNTRY_CODE_BS,
    COUNTRY_CODE_BD,
    COUNTRY_CODE_UZ,
    COUNTRY_CODE_JE,
    COUNTRY_CODE_WF,
    COUNTRY_CODE_HK,
    COUNTRY_CODE_IN,
    COUNTRY_CODE_LY,
    COUNTRY_CODE_NI,
    COUNTRY_CODE_MS,
    COUNTRY_CODE_NR,
    COUNTRY_CODE_UG,
    COUNTRY_CODE_JM,
    COUNTRY_CODE_AG,
    COUNTRY_CODE_KY,
    COUNTRY_CODE_TH,
    COUNTRY_CODE_IS,
    COUNTRY_CODE_LT,
    COUNTRY_CODE_LA,
    COUNTRY_CODE_SH,
    COUNTRY_CODE_NA,
    COUNTRY_CODE_PW,
    COUNTRY_CODE_PL,
    COUNTRY_CODE_CN,
    COUNTRY_CODE_VU,
    COUNTRY_CODE_CU,
    COUNTRY_CODE_VN,
    COUNTRY_CODE_GT,
    COUNTRY_CODE_AT,
    COUNTRY_CODE_ZM,
    COUNTRY_CODE_AO,
    COUNTRY_CODE_LI,
    COUNTRY_CODE_NL,
    COUNTRY_CODE_MH,
    COUNTRY_CODE_BV,
    COUNTRY_CODE_CF,
    COUNTRY_CODE_GY,
    COUNTRY_CODE_GL,
    COUNTRY_CODE_DK,
    COUNTRY_CODE_KR,
    COUNTRY_CODE_YT,
    COUNTRY_CODE_KI,
    COUNTRY_CODE_HN,
    COUNTRY_CODE_SM,
    COUNTRY_CODE_ME,
    COUNTRY_CODE_PR,
    COUNTRY_CODE_MP,
    COUNTRY_CODE_PG,
    COUNTRY_CODE_BN,
    COUNTRY_CODE_FM,
    COUNTRY_CODE_VE,
    COUNTRY_CODE_GQ,
    COUNTRY_CODE_GD,
    COUNTRY_CODE_ZW,
    COUNTRY_CODE_AQ,
    COUNTRY_CODE_JP,
    COUNTRY_CODE_WS,
    COUNTRY_CODE_AD,
    COUNTRY_CODE_KZ,
    COUNTRY_CODE_TK,
    COUNTRY_CODE_ES,
    COUNTRY_CODE_SE,
    COUNTRY_CODE_MM,
    COUNTRY_CODE_MX,
    COUNTRY_CODE_BF,
    COUNTRY_CODE_CK,
    COUNTRY_CODE_CV,
    COUNTRY_CODE_XK,
    COUNTRY_CODE_GI,
    COUNTRY_CODE_KW,
    COUNTRY_CODE_TV,
    COUNTRY_CODE_AL,
    COUNTRY_CODE_TC,
    COUNTRY_CODE_SV,
    COUNTRY_CODE_MU,
    COUNTRY_CODE_CC,
    COUNTRY_CODE_GA,
    COUNTRY_CODE_JO,
    COUNTRY_CODE_TN,
    COUNTRY_CODE_EC,
    COUNTRY_CODE_SN,
    COUNTRY_CODE_NC,
    COUNTRY_CODE_BM,
    COUNTRY_CODE_CL,
    COUNTRY_CODE_UA,
    COUNTRY_CODE_GR,
    COUNTRY_CODE_AI,
    COUNTRY_CODE_KG,
    COUNTRY_CODE_TF,
    COUNTRY_CODE_LR,
    COUNTRY_CODE_BE,
    COUNTRY_CODE_CD,
    COUNTRY_CODE_UY,
    COUNTRY_CODE_KP,
    COUNTRY_CODE_EH,
    COUNTRY_CODE_IM,
    COUNTRY_CODE_SS,
    COUNTRY_CODE_MG,
    COUNTRY_CODE_NF,
    COUNTRY_CODE_MR,
    COUNTRY_CODE_PE,
    COUNTRY_CODE_BH,
    COUNTRY_CODE_FO,
    COUNTRY_CODE_VG,
    COUNTRY_CODE_GB,
    COUNTRY_CODE_AS,
    COUNTRY_CODE_AF,
    COUNTRY_CODE_IR,
    COUNTRY_CODE_LU,
    COUNTRY_CODE_RS,
    COUNTRY_CODE_IE,
    COUNTRY_CODE_LB,
    COUNTRY_CODE_SK,
    COUNTRY_CODE_MO,
    COUNTRY_CODE_MZ,
    COUNTRY_CODE_PM,
    COUNTRY_CODE_CI,
    COUNTRY_CODE_GW,
    COUNTRY_CODE_TT,
    COUNTRY_CODE_OM,
    COUNTRY_CODE_ST,
    COUNTRY_CODE_SC,
    COUNTRY_CODE_MW,
    COUNTRY_CODE_BW,
    COUNTRY_CODE_CA,
    COUNTRY_CODE_ZA,
    COUNTRY_CODE_TL,
    COUNTRY_CODE_KH,
    COUNTRY_CODE_EE,
    COUNTRY_CODE_HT,
    COUNTRY_CODE_SL,
    COUNTRY_CODE_MD,
    COUNTRY_CODE_NE,
    COUNTRY_CODE_PS,
    COUNTRY_CODE_BZ,
    COUNTRY_CODE_PH,
    COUNTRY_CODE_BO,
    COUNTRY_CODE_QA,
    COUNTRY_CODE_CY,
    COUNTRY_CODE_FJ,
    COUNTRY_CODE_GP,
    COUNTRY_CODE_GG,
    COUNTRY_CODE_KE,
    COUNTRY_CODE_TD,
    COUNTRY_CODE_ER,
    COUNTRY_CODE_SY,
    COUNTRY_CODE_SD,
    COUNTRY_CODE_ML,
    COUNTRY_CODE_BR,
    COUNTRY_CODE_BG,
    COUNTRY_CODE_FR,
    COUNTRY_CODE_GH,
    COUNTRY_CODE_AX,
    COUNTRY_CODE_DO,
    COUNTRY_CODE_TW,
    COUNTRY_CODE_KM,
    COUNTRY_CODE_RE,
    COUNTRY_CODE_IO,
    COUNTRY_CODE_MT,
    COUNTRY_CODE_NU,
    COUNTRY_CODE_BJ,
  ];

  static final $core.List<CountryCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 249);
  static CountryCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CountryCode._(super.value, super.name);
}

class DataProtectionAction extends $pb.ProtobufEnum {
  static const DataProtectionAction DATA_PROTECTION_ACTION_HASH =
      DataProtectionAction._(
          0, _omitEnumNames ? '' : 'DATA_PROTECTION_ACTION_HASH');
  static const DataProtectionAction DATA_PROTECTION_ACTION_SUBSTITUTION =
      DataProtectionAction._(
          1, _omitEnumNames ? '' : 'DATA_PROTECTION_ACTION_SUBSTITUTION');

  static const $core.List<DataProtectionAction> values = <DataProtectionAction>[
    DATA_PROTECTION_ACTION_HASH,
    DATA_PROTECTION_ACTION_SUBSTITUTION,
  ];

  static final $core.List<DataProtectionAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DataProtectionAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataProtectionAction._(super.value, super.name);
}

class FailureReason extends $pb.ProtobufEnum {
  static const FailureReason FAILURE_REASON_TOKEN_MISSING =
      FailureReason._(0, _omitEnumNames ? '' : 'FAILURE_REASON_TOKEN_MISSING');
  static const FailureReason FAILURE_REASON_TOKEN_INVALID =
      FailureReason._(1, _omitEnumNames ? '' : 'FAILURE_REASON_TOKEN_INVALID');
  static const FailureReason FAILURE_REASON_TOKEN_DOMAIN_MISMATCH =
      FailureReason._(
          2, _omitEnumNames ? '' : 'FAILURE_REASON_TOKEN_DOMAIN_MISMATCH');
  static const FailureReason FAILURE_REASON_TOKEN_EXPIRED =
      FailureReason._(3, _omitEnumNames ? '' : 'FAILURE_REASON_TOKEN_EXPIRED');

  static const $core.List<FailureReason> values = <FailureReason>[
    FAILURE_REASON_TOKEN_MISSING,
    FAILURE_REASON_TOKEN_INVALID,
    FAILURE_REASON_TOKEN_DOMAIN_MISMATCH,
    FAILURE_REASON_TOKEN_EXPIRED,
  ];

  static final $core.List<FailureReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static FailureReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FailureReason._(super.value, super.name);
}

class FallbackBehavior extends $pb.ProtobufEnum {
  static const FallbackBehavior FALLBACK_BEHAVIOR_MATCH =
      FallbackBehavior._(0, _omitEnumNames ? '' : 'FALLBACK_BEHAVIOR_MATCH');
  static const FallbackBehavior FALLBACK_BEHAVIOR_NO_MATCH =
      FallbackBehavior._(1, _omitEnumNames ? '' : 'FALLBACK_BEHAVIOR_NO_MATCH');

  static const $core.List<FallbackBehavior> values = <FallbackBehavior>[
    FALLBACK_BEHAVIOR_MATCH,
    FALLBACK_BEHAVIOR_NO_MATCH,
  ];

  static final $core.List<FallbackBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FallbackBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FallbackBehavior._(super.value, super.name);
}

class FieldToProtectType extends $pb.ProtobufEnum {
  static const FieldToProtectType FIELD_TO_PROTECT_TYPE_SINGLE_HEADER =
      FieldToProtectType._(
          0, _omitEnumNames ? '' : 'FIELD_TO_PROTECT_TYPE_SINGLE_HEADER');
  static const FieldToProtectType FIELD_TO_PROTECT_TYPE_QUERY_STRING =
      FieldToProtectType._(
          1, _omitEnumNames ? '' : 'FIELD_TO_PROTECT_TYPE_QUERY_STRING');
  static const FieldToProtectType FIELD_TO_PROTECT_TYPE_BODY =
      FieldToProtectType._(
          2, _omitEnumNames ? '' : 'FIELD_TO_PROTECT_TYPE_BODY');
  static const FieldToProtectType FIELD_TO_PROTECT_TYPE_SINGLE_COOKIE =
      FieldToProtectType._(
          3, _omitEnumNames ? '' : 'FIELD_TO_PROTECT_TYPE_SINGLE_COOKIE');
  static const FieldToProtectType FIELD_TO_PROTECT_TYPE_SINGLE_QUERY_ARGUMENT =
      FieldToProtectType._(4,
          _omitEnumNames ? '' : 'FIELD_TO_PROTECT_TYPE_SINGLE_QUERY_ARGUMENT');

  static const $core.List<FieldToProtectType> values = <FieldToProtectType>[
    FIELD_TO_PROTECT_TYPE_SINGLE_HEADER,
    FIELD_TO_PROTECT_TYPE_QUERY_STRING,
    FIELD_TO_PROTECT_TYPE_BODY,
    FIELD_TO_PROTECT_TYPE_SINGLE_COOKIE,
    FIELD_TO_PROTECT_TYPE_SINGLE_QUERY_ARGUMENT,
  ];

  static final $core.List<FieldToProtectType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static FieldToProtectType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FieldToProtectType._(super.value, super.name);
}

class FilterBehavior extends $pb.ProtobufEnum {
  static const FilterBehavior FILTER_BEHAVIOR_DROP =
      FilterBehavior._(0, _omitEnumNames ? '' : 'FILTER_BEHAVIOR_DROP');
  static const FilterBehavior FILTER_BEHAVIOR_KEEP =
      FilterBehavior._(1, _omitEnumNames ? '' : 'FILTER_BEHAVIOR_KEEP');

  static const $core.List<FilterBehavior> values = <FilterBehavior>[
    FILTER_BEHAVIOR_DROP,
    FILTER_BEHAVIOR_KEEP,
  ];

  static final $core.List<FilterBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FilterBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FilterBehavior._(super.value, super.name);
}

class FilterRequirement extends $pb.ProtobufEnum {
  static const FilterRequirement FILTER_REQUIREMENT_MEETS_ANY =
      FilterRequirement._(
          0, _omitEnumNames ? '' : 'FILTER_REQUIREMENT_MEETS_ANY');
  static const FilterRequirement FILTER_REQUIREMENT_MEETS_ALL =
      FilterRequirement._(
          1, _omitEnumNames ? '' : 'FILTER_REQUIREMENT_MEETS_ALL');

  static const $core.List<FilterRequirement> values = <FilterRequirement>[
    FILTER_REQUIREMENT_MEETS_ANY,
    FILTER_REQUIREMENT_MEETS_ALL,
  ];

  static final $core.List<FilterRequirement?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FilterRequirement? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FilterRequirement._(super.value, super.name);
}

class ForwardedIPPosition extends $pb.ProtobufEnum {
  static const ForwardedIPPosition FORWARDED_I_P_POSITION_LAST =
      ForwardedIPPosition._(
          0, _omitEnumNames ? '' : 'FORWARDED_I_P_POSITION_LAST');
  static const ForwardedIPPosition FORWARDED_I_P_POSITION_ANY =
      ForwardedIPPosition._(
          1, _omitEnumNames ? '' : 'FORWARDED_I_P_POSITION_ANY');
  static const ForwardedIPPosition FORWARDED_I_P_POSITION_FIRST =
      ForwardedIPPosition._(
          2, _omitEnumNames ? '' : 'FORWARDED_I_P_POSITION_FIRST');

  static const $core.List<ForwardedIPPosition> values = <ForwardedIPPosition>[
    FORWARDED_I_P_POSITION_LAST,
    FORWARDED_I_P_POSITION_ANY,
    FORWARDED_I_P_POSITION_FIRST,
  ];

  static final $core.List<ForwardedIPPosition?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ForwardedIPPosition? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ForwardedIPPosition._(super.value, super.name);
}

class IPAddressVersion extends $pb.ProtobufEnum {
  static const IPAddressVersion I_P_ADDRESS_VERSION_IPV6 =
      IPAddressVersion._(0, _omitEnumNames ? '' : 'I_P_ADDRESS_VERSION_IPV6');
  static const IPAddressVersion I_P_ADDRESS_VERSION_IPV4 =
      IPAddressVersion._(1, _omitEnumNames ? '' : 'I_P_ADDRESS_VERSION_IPV4');

  static const $core.List<IPAddressVersion> values = <IPAddressVersion>[
    I_P_ADDRESS_VERSION_IPV6,
    I_P_ADDRESS_VERSION_IPV4,
  ];

  static final $core.List<IPAddressVersion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IPAddressVersion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IPAddressVersion._(super.value, super.name);
}

class InspectionLevel extends $pb.ProtobufEnum {
  static const InspectionLevel INSPECTION_LEVEL_TARGETED =
      InspectionLevel._(0, _omitEnumNames ? '' : 'INSPECTION_LEVEL_TARGETED');
  static const InspectionLevel INSPECTION_LEVEL_COMMON =
      InspectionLevel._(1, _omitEnumNames ? '' : 'INSPECTION_LEVEL_COMMON');

  static const $core.List<InspectionLevel> values = <InspectionLevel>[
    INSPECTION_LEVEL_TARGETED,
    INSPECTION_LEVEL_COMMON,
  ];

  static final $core.List<InspectionLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InspectionLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InspectionLevel._(super.value, super.name);
}

class JsonMatchScope extends $pb.ProtobufEnum {
  static const JsonMatchScope JSON_MATCH_SCOPE_VALUE =
      JsonMatchScope._(0, _omitEnumNames ? '' : 'JSON_MATCH_SCOPE_VALUE');
  static const JsonMatchScope JSON_MATCH_SCOPE_ALL =
      JsonMatchScope._(1, _omitEnumNames ? '' : 'JSON_MATCH_SCOPE_ALL');
  static const JsonMatchScope JSON_MATCH_SCOPE_KEY =
      JsonMatchScope._(2, _omitEnumNames ? '' : 'JSON_MATCH_SCOPE_KEY');

  static const $core.List<JsonMatchScope> values = <JsonMatchScope>[
    JSON_MATCH_SCOPE_VALUE,
    JSON_MATCH_SCOPE_ALL,
    JSON_MATCH_SCOPE_KEY,
  ];

  static final $core.List<JsonMatchScope?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static JsonMatchScope? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const JsonMatchScope._(super.value, super.name);
}

class LabelMatchScope extends $pb.ProtobufEnum {
  static const LabelMatchScope LABEL_MATCH_SCOPE_LABEL =
      LabelMatchScope._(0, _omitEnumNames ? '' : 'LABEL_MATCH_SCOPE_LABEL');
  static const LabelMatchScope LABEL_MATCH_SCOPE_NAMESPACE =
      LabelMatchScope._(1, _omitEnumNames ? '' : 'LABEL_MATCH_SCOPE_NAMESPACE');

  static const $core.List<LabelMatchScope> values = <LabelMatchScope>[
    LABEL_MATCH_SCOPE_LABEL,
    LABEL_MATCH_SCOPE_NAMESPACE,
  ];

  static final $core.List<LabelMatchScope?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static LabelMatchScope? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LabelMatchScope._(super.value, super.name);
}

class LogScope extends $pb.ProtobufEnum {
  static const LogScope LOG_SCOPE_CUSTOMER =
      LogScope._(0, _omitEnumNames ? '' : 'LOG_SCOPE_CUSTOMER');
  static const LogScope LOG_SCOPE_SECURITY_LAKE =
      LogScope._(1, _omitEnumNames ? '' : 'LOG_SCOPE_SECURITY_LAKE');
  static const LogScope LOG_SCOPE_CLOUDWATCH_TELEMETRY_RULE_MANAGED =
      LogScope._(2,
          _omitEnumNames ? '' : 'LOG_SCOPE_CLOUDWATCH_TELEMETRY_RULE_MANAGED');

  static const $core.List<LogScope> values = <LogScope>[
    LOG_SCOPE_CUSTOMER,
    LOG_SCOPE_SECURITY_LAKE,
    LOG_SCOPE_CLOUDWATCH_TELEMETRY_RULE_MANAGED,
  ];

  static final $core.List<LogScope?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static LogScope? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogScope._(super.value, super.name);
}

class LogType extends $pb.ProtobufEnum {
  static const LogType LOG_TYPE_WAF_LOGS =
      LogType._(0, _omitEnumNames ? '' : 'LOG_TYPE_WAF_LOGS');

  static const $core.List<LogType> values = <LogType>[
    LOG_TYPE_WAF_LOGS,
  ];

  static final $core.List<LogType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static LogType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogType._(super.value, super.name);
}

class LowReputationMode extends $pb.ProtobufEnum {
  static const LowReputationMode LOW_REPUTATION_MODE_ALWAYS_ON =
      LowReputationMode._(
          0, _omitEnumNames ? '' : 'LOW_REPUTATION_MODE_ALWAYS_ON');
  static const LowReputationMode LOW_REPUTATION_MODE_ACTIVE_UNDER_DDOS =
      LowReputationMode._(
          1, _omitEnumNames ? '' : 'LOW_REPUTATION_MODE_ACTIVE_UNDER_DDOS');

  static const $core.List<LowReputationMode> values = <LowReputationMode>[
    LOW_REPUTATION_MODE_ALWAYS_ON,
    LOW_REPUTATION_MODE_ACTIVE_UNDER_DDOS,
  ];

  static final $core.List<LowReputationMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static LowReputationMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LowReputationMode._(super.value, super.name);
}

class MapMatchScope extends $pb.ProtobufEnum {
  static const MapMatchScope MAP_MATCH_SCOPE_VALUE =
      MapMatchScope._(0, _omitEnumNames ? '' : 'MAP_MATCH_SCOPE_VALUE');
  static const MapMatchScope MAP_MATCH_SCOPE_ALL =
      MapMatchScope._(1, _omitEnumNames ? '' : 'MAP_MATCH_SCOPE_ALL');
  static const MapMatchScope MAP_MATCH_SCOPE_KEY =
      MapMatchScope._(2, _omitEnumNames ? '' : 'MAP_MATCH_SCOPE_KEY');

  static const $core.List<MapMatchScope> values = <MapMatchScope>[
    MAP_MATCH_SCOPE_VALUE,
    MAP_MATCH_SCOPE_ALL,
    MAP_MATCH_SCOPE_KEY,
  ];

  static final $core.List<MapMatchScope?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static MapMatchScope? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MapMatchScope._(super.value, super.name);
}

class OversizeHandling extends $pb.ProtobufEnum {
  static const OversizeHandling OVERSIZE_HANDLING_MATCH =
      OversizeHandling._(0, _omitEnumNames ? '' : 'OVERSIZE_HANDLING_MATCH');
  static const OversizeHandling OVERSIZE_HANDLING_CONTINUE =
      OversizeHandling._(1, _omitEnumNames ? '' : 'OVERSIZE_HANDLING_CONTINUE');
  static const OversizeHandling OVERSIZE_HANDLING_NO_MATCH =
      OversizeHandling._(2, _omitEnumNames ? '' : 'OVERSIZE_HANDLING_NO_MATCH');

  static const $core.List<OversizeHandling> values = <OversizeHandling>[
    OVERSIZE_HANDLING_MATCH,
    OVERSIZE_HANDLING_CONTINUE,
    OVERSIZE_HANDLING_NO_MATCH,
  ];

  static final $core.List<OversizeHandling?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static OversizeHandling? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OversizeHandling._(super.value, super.name);
}

class ParameterExceptionField extends $pb.ProtobufEnum {
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_COOKIE_MATCH_PATTERN =
      ParameterExceptionField._(
          0,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_COOKIE_MATCH_PATTERN');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_IP_ADDRESS_VERSION = ParameterExceptionField._(
          1,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_IP_ADDRESS_VERSION');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_NOT_STATEMENT =
      ParameterExceptionField._(
          2, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_NOT_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_FORWARDED_IP_CONFIG = ParameterExceptionField._(
          3,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_FORWARDED_IP_CONFIG');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_LOW_REPUTATION_MODE = ParameterExceptionField._(
          4,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_LOW_REPUTATION_MODE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_DATA_PROTECTION_CONFIG =
      ParameterExceptionField._(
          5,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_DATA_PROTECTION_CONFIG');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_CHANGE_PROPAGATION_STATUS =
      ParameterExceptionField._(
          6,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_CHANGE_PROPAGATION_STATUS');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_EXPIRE_TIMESTAMP = ParameterExceptionField._(7,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_EXPIRE_TIMESTAMP');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_ASSOCIABLE_RESOURCE = ParameterExceptionField._(
          8,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_ASSOCIABLE_RESOURCE');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_SINGLE_HEADER =
      ParameterExceptionField._(
          9, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_SINGLE_HEADER');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_ACP_RULE_SET_RESPONSE_INSPECTION =
      ParameterExceptionField._(
          10,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_ACP_RULE_SET_RESPONSE_INSPECTION');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_IP_SET_FORWARDED_IP_CONFIG =
      ParameterExceptionField._(
          11,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_IP_SET_FORWARDED_IP_CONFIG');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_EXCLUDED_RULE =
      ParameterExceptionField._(
          12, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_EXCLUDED_RULE');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_IP_ADDRESS =
      ParameterExceptionField._(
          13, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_IP_ADDRESS');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_METRIC_NAME =
      ParameterExceptionField._(
          14, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_METRIC_NAME');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_FILTER_CONDITION = ParameterExceptionField._(15,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_FILTER_CONDITION');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_OVERRIDE_ACTION = ParameterExceptionField._(16,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_OVERRIDE_ACTION');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_WEB_ACL =
      ParameterExceptionField._(
          17, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_WEB_ACL');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_TEXT_TRANSFORMATION = ParameterExceptionField._(
          18,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_TEXT_TRANSFORMATION');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_JSON_MATCH_PATTERN = ParameterExceptionField._(
          19,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_JSON_MATCH_PATTERN');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_BODY_PARSING_FALLBACK_BEHAVIOR =
      ParameterExceptionField._(
          20,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_BODY_PARSING_FALLBACK_BEHAVIOR');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_PAYLOAD_TYPE =
      ParameterExceptionField._(
          21, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_PAYLOAD_TYPE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_CHALLENGE_CONFIG = ParameterExceptionField._(22,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_CHALLENGE_CONFIG');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_HEADER_MATCH_PATTERN =
      ParameterExceptionField._(
          23,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_HEADER_MATCH_PATTERN');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_HEADER_NAME =
      ParameterExceptionField._(
          24, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_HEADER_NAME');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_AND_STATEMENT =
      ParameterExceptionField._(
          25, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_AND_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_SET = ParameterExceptionField._(26,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_SET');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_MAP_MATCH_SCOPE = ParameterExceptionField._(27,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_MAP_MATCH_SCOPE');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_OR_STATEMENT =
      ParameterExceptionField._(
          28, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_OR_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_RESPONSE_CONTENT_TYPE =
      ParameterExceptionField._(
          29,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_RESPONSE_CONTENT_TYPE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_RATE_BASED_STATEMENT =
      ParameterExceptionField._(
          30,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_RATE_BASED_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_LOGGING_FILTER = ParameterExceptionField._(
          31, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_LOGGING_FILTER');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_SQLI_MATCH_STATEMENT =
      ParameterExceptionField._(
          32,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_SQLI_MATCH_STATEMENT');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_STATEMENT =
      ParameterExceptionField._(
          33, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_IP_SET_REFERENCE_STATEMENT =
      ParameterExceptionField._(
          34,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_IP_SET_REFERENCE_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_FALLBACK_BEHAVIOR = ParameterExceptionField._(
          35,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_FALLBACK_BEHAVIOR');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_LABEL_MATCH_STATEMENT =
      ParameterExceptionField._(
          36,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_LABEL_MATCH_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_SET_STATEMENT =
      ParameterExceptionField._(
          37,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_SET_STATEMENT');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RULE =
      ParameterExceptionField._(
          38, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RULE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_XSS_MATCH_STATEMENT = ParameterExceptionField._(
          39,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_XSS_MATCH_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_JSON_MATCH_SCOPE = ParameterExceptionField._(40,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_JSON_MATCH_SCOPE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_GROUP_CONFIG =
      ParameterExceptionField._(
          41,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_GROUP_CONFIG');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_SCOPE_DOWN =
      ParameterExceptionField._(
          42, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_SCOPE_DOWN');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_REGEX_PATTERN_REFERENCE_STATEMENT =
      ParameterExceptionField._(
          43,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_REGEX_PATTERN_REFERENCE_STATEMENT');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_TAGS =
      ParameterExceptionField._(
          44, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_TAGS');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_SIZE_CONSTRAINT_STATEMENT =
      ParameterExceptionField._(
          45,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_SIZE_CONSTRAINT_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_GEO_MATCH_STATEMENT = ParameterExceptionField._(
          46,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_GEO_MATCH_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_CUSTOM_RESPONSE = ParameterExceptionField._(47,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_CUSTOM_RESPONSE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_SINGLE_QUERY_ARGUMENT =
      ParameterExceptionField._(
          48,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_SINGLE_QUERY_ARGUMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_ATP_RULE_SET_RESPONSE_INSPECTION =
      ParameterExceptionField._(
          49,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_ATP_RULE_SET_RESPONSE_INSPECTION');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_ASSOCIATED_RESOURCE_TYPE =
      ParameterExceptionField._(
          50,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_ASSOCIATED_RESOURCE_TYPE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_FIREWALL_MANAGER_STATEMENT =
      ParameterExceptionField._(
          51,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_FIREWALL_MANAGER_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_CUSTOM_RESPONSE_BODY =
      ParameterExceptionField._(
          52,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_CUSTOM_RESPONSE_BODY');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_SCOPE_VALUE =
      ParameterExceptionField._(
          53, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_SCOPE_VALUE');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_POSITION =
      ParameterExceptionField._(
          54, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_POSITION');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_OVERSIZE_HANDLING = ParameterExceptionField._(
          55,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_OVERSIZE_HANDLING');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_LOG_DESTINATION = ParameterExceptionField._(56,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_LOG_DESTINATION');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_IP_SET =
      ParameterExceptionField._(
          57, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_IP_SET');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_REGEX_PATTERN_SET = ParameterExceptionField._(
          58,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_REGEX_PATTERN_SET');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RULE_GROUP =
      ParameterExceptionField._(
          59, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RULE_GROUP');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RESOURCE_TYPE =
      ParameterExceptionField._(
          60, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RESOURCE_TYPE');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RESOURCE_ARN =
      ParameterExceptionField._(
          61, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RESOURCE_ARN');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_DEFAULT_ACTION = ParameterExceptionField._(
          62, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_DEFAULT_ACTION');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_STATEMENT =
      ParameterExceptionField._(
          63,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_STATEMENT');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_CUSTOM_REQUEST_HANDLING =
      ParameterExceptionField._(
          64,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_CUSTOM_REQUEST_HANDLING');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_TAG_KEYS =
      ParameterExceptionField._(
          65, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_TAG_KEYS');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_FIELD_TO_MATCH = ParameterExceptionField._(
          66, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_FIELD_TO_MATCH');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RULE_ACTION =
      ParameterExceptionField._(
          67, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RULE_ACTION');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_CUSTOM_KEYS =
      ParameterExceptionField._(
          68, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_CUSTOM_KEYS');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_ENTITY_LIMIT =
      ParameterExceptionField._(
          69, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_ENTITY_LIMIT');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_TOKEN_DOMAIN =
      ParameterExceptionField._(
          70, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_TOKEN_DOMAIN');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_RULE_GROUP_REFERENCE_STATEMENT =
      ParameterExceptionField._(
          71,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_RULE_GROUP_REFERENCE_STATEMENT');

  static const $core.List<ParameterExceptionField> values =
      <ParameterExceptionField>[
    PARAMETER_EXCEPTION_FIELD_COOKIE_MATCH_PATTERN,
    PARAMETER_EXCEPTION_FIELD_IP_ADDRESS_VERSION,
    PARAMETER_EXCEPTION_FIELD_NOT_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_FORWARDED_IP_CONFIG,
    PARAMETER_EXCEPTION_FIELD_LOW_REPUTATION_MODE,
    PARAMETER_EXCEPTION_FIELD_DATA_PROTECTION_CONFIG,
    PARAMETER_EXCEPTION_FIELD_CHANGE_PROPAGATION_STATUS,
    PARAMETER_EXCEPTION_FIELD_EXPIRE_TIMESTAMP,
    PARAMETER_EXCEPTION_FIELD_ASSOCIABLE_RESOURCE,
    PARAMETER_EXCEPTION_FIELD_SINGLE_HEADER,
    PARAMETER_EXCEPTION_FIELD_ACP_RULE_SET_RESPONSE_INSPECTION,
    PARAMETER_EXCEPTION_FIELD_IP_SET_FORWARDED_IP_CONFIG,
    PARAMETER_EXCEPTION_FIELD_EXCLUDED_RULE,
    PARAMETER_EXCEPTION_FIELD_IP_ADDRESS,
    PARAMETER_EXCEPTION_FIELD_METRIC_NAME,
    PARAMETER_EXCEPTION_FIELD_FILTER_CONDITION,
    PARAMETER_EXCEPTION_FIELD_OVERRIDE_ACTION,
    PARAMETER_EXCEPTION_FIELD_WEB_ACL,
    PARAMETER_EXCEPTION_FIELD_TEXT_TRANSFORMATION,
    PARAMETER_EXCEPTION_FIELD_JSON_MATCH_PATTERN,
    PARAMETER_EXCEPTION_FIELD_BODY_PARSING_FALLBACK_BEHAVIOR,
    PARAMETER_EXCEPTION_FIELD_PAYLOAD_TYPE,
    PARAMETER_EXCEPTION_FIELD_CHALLENGE_CONFIG,
    PARAMETER_EXCEPTION_FIELD_HEADER_MATCH_PATTERN,
    PARAMETER_EXCEPTION_FIELD_HEADER_NAME,
    PARAMETER_EXCEPTION_FIELD_AND_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_SET,
    PARAMETER_EXCEPTION_FIELD_MAP_MATCH_SCOPE,
    PARAMETER_EXCEPTION_FIELD_OR_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_RESPONSE_CONTENT_TYPE,
    PARAMETER_EXCEPTION_FIELD_RATE_BASED_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_LOGGING_FILTER,
    PARAMETER_EXCEPTION_FIELD_SQLI_MATCH_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_IP_SET_REFERENCE_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_FALLBACK_BEHAVIOR,
    PARAMETER_EXCEPTION_FIELD_LABEL_MATCH_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_SET_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_RULE,
    PARAMETER_EXCEPTION_FIELD_XSS_MATCH_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_JSON_MATCH_SCOPE,
    PARAMETER_EXCEPTION_FIELD_MANAGED_RULE_GROUP_CONFIG,
    PARAMETER_EXCEPTION_FIELD_SCOPE_DOWN,
    PARAMETER_EXCEPTION_FIELD_REGEX_PATTERN_REFERENCE_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_TAGS,
    PARAMETER_EXCEPTION_FIELD_SIZE_CONSTRAINT_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_GEO_MATCH_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_CUSTOM_RESPONSE,
    PARAMETER_EXCEPTION_FIELD_SINGLE_QUERY_ARGUMENT,
    PARAMETER_EXCEPTION_FIELD_ATP_RULE_SET_RESPONSE_INSPECTION,
    PARAMETER_EXCEPTION_FIELD_ASSOCIATED_RESOURCE_TYPE,
    PARAMETER_EXCEPTION_FIELD_FIREWALL_MANAGER_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_CUSTOM_RESPONSE_BODY,
    PARAMETER_EXCEPTION_FIELD_SCOPE_VALUE,
    PARAMETER_EXCEPTION_FIELD_POSITION,
    PARAMETER_EXCEPTION_FIELD_OVERSIZE_HANDLING,
    PARAMETER_EXCEPTION_FIELD_LOG_DESTINATION,
    PARAMETER_EXCEPTION_FIELD_IP_SET,
    PARAMETER_EXCEPTION_FIELD_REGEX_PATTERN_SET,
    PARAMETER_EXCEPTION_FIELD_RULE_GROUP,
    PARAMETER_EXCEPTION_FIELD_RESOURCE_TYPE,
    PARAMETER_EXCEPTION_FIELD_RESOURCE_ARN,
    PARAMETER_EXCEPTION_FIELD_DEFAULT_ACTION,
    PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_STATEMENT,
    PARAMETER_EXCEPTION_FIELD_CUSTOM_REQUEST_HANDLING,
    PARAMETER_EXCEPTION_FIELD_TAG_KEYS,
    PARAMETER_EXCEPTION_FIELD_FIELD_TO_MATCH,
    PARAMETER_EXCEPTION_FIELD_RULE_ACTION,
    PARAMETER_EXCEPTION_FIELD_CUSTOM_KEYS,
    PARAMETER_EXCEPTION_FIELD_ENTITY_LIMIT,
    PARAMETER_EXCEPTION_FIELD_TOKEN_DOMAIN,
    PARAMETER_EXCEPTION_FIELD_RULE_GROUP_REFERENCE_STATEMENT,
  ];

  static final $core.List<ParameterExceptionField?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 71);
  static ParameterExceptionField? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ParameterExceptionField._(super.value, super.name);
}

class PayloadType extends $pb.ProtobufEnum {
  static const PayloadType PAYLOAD_TYPE_JSON =
      PayloadType._(0, _omitEnumNames ? '' : 'PAYLOAD_TYPE_JSON');
  static const PayloadType PAYLOAD_TYPE_FORM_ENCODED =
      PayloadType._(1, _omitEnumNames ? '' : 'PAYLOAD_TYPE_FORM_ENCODED');

  static const $core.List<PayloadType> values = <PayloadType>[
    PAYLOAD_TYPE_JSON,
    PAYLOAD_TYPE_FORM_ENCODED,
  ];

  static final $core.List<PayloadType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PayloadType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PayloadType._(super.value, super.name);
}

class Platform extends $pb.ProtobufEnum {
  static const Platform PLATFORM_ANDROID =
      Platform._(0, _omitEnumNames ? '' : 'PLATFORM_ANDROID');
  static const Platform PLATFORM_IOS =
      Platform._(1, _omitEnumNames ? '' : 'PLATFORM_IOS');

  static const $core.List<Platform> values = <Platform>[
    PLATFORM_ANDROID,
    PLATFORM_IOS,
  ];

  static final $core.List<Platform?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static Platform? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Platform._(super.value, super.name);
}

class PositionalConstraint extends $pb.ProtobufEnum {
  static const PositionalConstraint POSITIONAL_CONSTRAINT_EXACTLY =
      PositionalConstraint._(
          0, _omitEnumNames ? '' : 'POSITIONAL_CONSTRAINT_EXACTLY');
  static const PositionalConstraint POSITIONAL_CONSTRAINT_CONTAINS =
      PositionalConstraint._(
          1, _omitEnumNames ? '' : 'POSITIONAL_CONSTRAINT_CONTAINS');
  static const PositionalConstraint POSITIONAL_CONSTRAINT_STARTS_WITH =
      PositionalConstraint._(
          2, _omitEnumNames ? '' : 'POSITIONAL_CONSTRAINT_STARTS_WITH');
  static const PositionalConstraint POSITIONAL_CONSTRAINT_ENDS_WITH =
      PositionalConstraint._(
          3, _omitEnumNames ? '' : 'POSITIONAL_CONSTRAINT_ENDS_WITH');
  static const PositionalConstraint POSITIONAL_CONSTRAINT_CONTAINS_WORD =
      PositionalConstraint._(
          4, _omitEnumNames ? '' : 'POSITIONAL_CONSTRAINT_CONTAINS_WORD');

  static const $core.List<PositionalConstraint> values = <PositionalConstraint>[
    POSITIONAL_CONSTRAINT_EXACTLY,
    POSITIONAL_CONSTRAINT_CONTAINS,
    POSITIONAL_CONSTRAINT_STARTS_WITH,
    POSITIONAL_CONSTRAINT_ENDS_WITH,
    POSITIONAL_CONSTRAINT_CONTAINS_WORD,
  ];

  static final $core.List<PositionalConstraint?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static PositionalConstraint? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PositionalConstraint._(super.value, super.name);
}

class RateBasedStatementAggregateKeyType extends $pb.ProtobufEnum {
  static const RateBasedStatementAggregateKeyType
      RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_CONSTANT =
      RateBasedStatementAggregateKeyType._(
          0,
          _omitEnumNames
              ? ''
              : 'RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_CONSTANT');
  static const RateBasedStatementAggregateKeyType
      RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_IP =
      RateBasedStatementAggregateKeyType._(1,
          _omitEnumNames ? '' : 'RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_IP');
  static const RateBasedStatementAggregateKeyType
      RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_FORWARDED_IP =
      RateBasedStatementAggregateKeyType._(
          2,
          _omitEnumNames
              ? ''
              : 'RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_FORWARDED_IP');
  static const RateBasedStatementAggregateKeyType
      RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_CUSTOM_KEYS =
      RateBasedStatementAggregateKeyType._(
          3,
          _omitEnumNames
              ? ''
              : 'RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_CUSTOM_KEYS');

  static const $core.List<RateBasedStatementAggregateKeyType> values =
      <RateBasedStatementAggregateKeyType>[
    RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_CONSTANT,
    RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_IP,
    RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_FORWARDED_IP,
    RATE_BASED_STATEMENT_AGGREGATE_KEY_TYPE_CUSTOM_KEYS,
  ];

  static final $core.List<RateBasedStatementAggregateKeyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static RateBasedStatementAggregateKeyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RateBasedStatementAggregateKeyType._(super.value, super.name);
}

class ResourceType extends $pb.ProtobufEnum {
  static const ResourceType RESOURCE_TYPE_VERIFIED_ACCESS_INSTANCE =
      ResourceType._(
          0, _omitEnumNames ? '' : 'RESOURCE_TYPE_VERIFIED_ACCESS_INSTANCE');
  static const ResourceType RESOURCE_TYPE_COGNITIO_USER_POOL = ResourceType._(
      1, _omitEnumNames ? '' : 'RESOURCE_TYPE_COGNITIO_USER_POOL');
  static const ResourceType RESOURCE_TYPE_APPLICATION_LOAD_BALANCER =
      ResourceType._(
          2, _omitEnumNames ? '' : 'RESOURCE_TYPE_APPLICATION_LOAD_BALANCER');
  static const ResourceType RESOURCE_TYPE_APPSYNC =
      ResourceType._(3, _omitEnumNames ? '' : 'RESOURCE_TYPE_APPSYNC');
  static const ResourceType RESOURCE_TYPE_AMPLIFY =
      ResourceType._(4, _omitEnumNames ? '' : 'RESOURCE_TYPE_AMPLIFY');
  static const ResourceType RESOURCE_TYPE_API_GATEWAY =
      ResourceType._(5, _omitEnumNames ? '' : 'RESOURCE_TYPE_API_GATEWAY');
  static const ResourceType RESOURCE_TYPE_APP_RUNNER_SERVICE = ResourceType._(
      6, _omitEnumNames ? '' : 'RESOURCE_TYPE_APP_RUNNER_SERVICE');

  static const $core.List<ResourceType> values = <ResourceType>[
    RESOURCE_TYPE_VERIFIED_ACCESS_INSTANCE,
    RESOURCE_TYPE_COGNITIO_USER_POOL,
    RESOURCE_TYPE_APPLICATION_LOAD_BALANCER,
    RESOURCE_TYPE_APPSYNC,
    RESOURCE_TYPE_AMPLIFY,
    RESOURCE_TYPE_API_GATEWAY,
    RESOURCE_TYPE_APP_RUNNER_SERVICE,
  ];

  static final $core.List<ResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static ResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceType._(super.value, super.name);
}

class ResponseContentType extends $pb.ProtobufEnum {
  static const ResponseContentType RESPONSE_CONTENT_TYPE_APPLICATION_JSON =
      ResponseContentType._(
          0, _omitEnumNames ? '' : 'RESPONSE_CONTENT_TYPE_APPLICATION_JSON');
  static const ResponseContentType RESPONSE_CONTENT_TYPE_TEXT_PLAIN =
      ResponseContentType._(
          1, _omitEnumNames ? '' : 'RESPONSE_CONTENT_TYPE_TEXT_PLAIN');
  static const ResponseContentType RESPONSE_CONTENT_TYPE_TEXT_HTML =
      ResponseContentType._(
          2, _omitEnumNames ? '' : 'RESPONSE_CONTENT_TYPE_TEXT_HTML');

  static const $core.List<ResponseContentType> values = <ResponseContentType>[
    RESPONSE_CONTENT_TYPE_APPLICATION_JSON,
    RESPONSE_CONTENT_TYPE_TEXT_PLAIN,
    RESPONSE_CONTENT_TYPE_TEXT_HTML,
  ];

  static final $core.List<ResponseContentType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ResponseContentType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResponseContentType._(super.value, super.name);
}

class Scope extends $pb.ProtobufEnum {
  static const Scope SCOPE_CLOUDFRONT =
      Scope._(0, _omitEnumNames ? '' : 'SCOPE_CLOUDFRONT');
  static const Scope SCOPE_REGIONAL =
      Scope._(1, _omitEnumNames ? '' : 'SCOPE_REGIONAL');

  static const $core.List<Scope> values = <Scope>[
    SCOPE_CLOUDFRONT,
    SCOPE_REGIONAL,
  ];

  static final $core.List<Scope?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static Scope? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Scope._(super.value, super.name);
}

class SensitivityLevel extends $pb.ProtobufEnum {
  static const SensitivityLevel SENSITIVITY_LEVEL_LOW =
      SensitivityLevel._(0, _omitEnumNames ? '' : 'SENSITIVITY_LEVEL_LOW');
  static const SensitivityLevel SENSITIVITY_LEVEL_HIGH =
      SensitivityLevel._(1, _omitEnumNames ? '' : 'SENSITIVITY_LEVEL_HIGH');

  static const $core.List<SensitivityLevel> values = <SensitivityLevel>[
    SENSITIVITY_LEVEL_LOW,
    SENSITIVITY_LEVEL_HIGH,
  ];

  static final $core.List<SensitivityLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SensitivityLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SensitivityLevel._(super.value, super.name);
}

class SensitivityToAct extends $pb.ProtobufEnum {
  static const SensitivityToAct SENSITIVITY_TO_ACT_MEDIUM =
      SensitivityToAct._(0, _omitEnumNames ? '' : 'SENSITIVITY_TO_ACT_MEDIUM');
  static const SensitivityToAct SENSITIVITY_TO_ACT_LOW =
      SensitivityToAct._(1, _omitEnumNames ? '' : 'SENSITIVITY_TO_ACT_LOW');
  static const SensitivityToAct SENSITIVITY_TO_ACT_HIGH =
      SensitivityToAct._(2, _omitEnumNames ? '' : 'SENSITIVITY_TO_ACT_HIGH');

  static const $core.List<SensitivityToAct> values = <SensitivityToAct>[
    SENSITIVITY_TO_ACT_MEDIUM,
    SENSITIVITY_TO_ACT_LOW,
    SENSITIVITY_TO_ACT_HIGH,
  ];

  static final $core.List<SensitivityToAct?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SensitivityToAct? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SensitivityToAct._(super.value, super.name);
}

class SizeInspectionLimit extends $pb.ProtobufEnum {
  static const SizeInspectionLimit SIZE_INSPECTION_LIMIT_KB_48 =
      SizeInspectionLimit._(
          0, _omitEnumNames ? '' : 'SIZE_INSPECTION_LIMIT_KB_48');
  static const SizeInspectionLimit SIZE_INSPECTION_LIMIT_KB_64 =
      SizeInspectionLimit._(
          1, _omitEnumNames ? '' : 'SIZE_INSPECTION_LIMIT_KB_64');
  static const SizeInspectionLimit SIZE_INSPECTION_LIMIT_KB_16 =
      SizeInspectionLimit._(
          2, _omitEnumNames ? '' : 'SIZE_INSPECTION_LIMIT_KB_16');
  static const SizeInspectionLimit SIZE_INSPECTION_LIMIT_KB_32 =
      SizeInspectionLimit._(
          3, _omitEnumNames ? '' : 'SIZE_INSPECTION_LIMIT_KB_32');

  static const $core.List<SizeInspectionLimit> values = <SizeInspectionLimit>[
    SIZE_INSPECTION_LIMIT_KB_48,
    SIZE_INSPECTION_LIMIT_KB_64,
    SIZE_INSPECTION_LIMIT_KB_16,
    SIZE_INSPECTION_LIMIT_KB_32,
  ];

  static final $core.List<SizeInspectionLimit?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static SizeInspectionLimit? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SizeInspectionLimit._(super.value, super.name);
}

class TextTransformationType extends $pb.ProtobufEnum {
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_SQL_HEX_DECODE =
      TextTransformationType._(
          0, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_SQL_HEX_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_URL_DECODE =
      TextTransformationType._(
          1, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_URL_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_HEX_DECODE =
      TextTransformationType._(
          2, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_HEX_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_NORMALIZE_PATH =
      TextTransformationType._(
          3, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_NORMALIZE_PATH');
  static const TextTransformationType
      TEXT_TRANSFORMATION_TYPE_COMPRESS_WHITE_SPACE = TextTransformationType._(
          4,
          _omitEnumNames
              ? ''
              : 'TEXT_TRANSFORMATION_TYPE_COMPRESS_WHITE_SPACE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_NONE =
      TextTransformationType._(
          5, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_NONE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_BASE64_DECODE =
      TextTransformationType._(
          6, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_BASE64_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_URL_DECODE_UNI =
      TextTransformationType._(
          7, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_URL_DECODE_UNI');
  static const TextTransformationType
      TEXT_TRANSFORMATION_TYPE_HTML_ENTITY_DECODE = TextTransformationType._(8,
          _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_HTML_ENTITY_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_JS_DECODE =
      TextTransformationType._(
          9, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_JS_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_MD5 =
      TextTransformationType._(
          10, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_MD5');
  static const TextTransformationType
      TEXT_TRANSFORMATION_TYPE_NORMALIZE_PATH_WIN = TextTransformationType._(11,
          _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_NORMALIZE_PATH_WIN');
  static const TextTransformationType
      TEXT_TRANSFORMATION_TYPE_REPLACE_COMMENTS = TextTransformationType._(12,
          _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_REPLACE_COMMENTS');
  static const TextTransformationType
      TEXT_TRANSFORMATION_TYPE_BASE64_DECODE_EXT = TextTransformationType._(13,
          _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_BASE64_DECODE_EXT');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_UTF8_TO_UNICODE =
      TextTransformationType._(
          14, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_UTF8_TO_UNICODE');
  static const TextTransformationType
      TEXT_TRANSFORMATION_TYPE_ESCAPE_SEQ_DECODE = TextTransformationType._(15,
          _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_ESCAPE_SEQ_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_CSS_DECODE =
      TextTransformationType._(
          16, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_CSS_DECODE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_REPLACE_NULLS =
      TextTransformationType._(
          17, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_REPLACE_NULLS');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_LOWERCASE =
      TextTransformationType._(
          18, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_LOWERCASE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_CMD_LINE =
      TextTransformationType._(
          19, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_CMD_LINE');
  static const TextTransformationType TEXT_TRANSFORMATION_TYPE_REMOVE_NULLS =
      TextTransformationType._(
          20, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_TYPE_REMOVE_NULLS');

  static const $core.List<TextTransformationType> values =
      <TextTransformationType>[
    TEXT_TRANSFORMATION_TYPE_SQL_HEX_DECODE,
    TEXT_TRANSFORMATION_TYPE_URL_DECODE,
    TEXT_TRANSFORMATION_TYPE_HEX_DECODE,
    TEXT_TRANSFORMATION_TYPE_NORMALIZE_PATH,
    TEXT_TRANSFORMATION_TYPE_COMPRESS_WHITE_SPACE,
    TEXT_TRANSFORMATION_TYPE_NONE,
    TEXT_TRANSFORMATION_TYPE_BASE64_DECODE,
    TEXT_TRANSFORMATION_TYPE_URL_DECODE_UNI,
    TEXT_TRANSFORMATION_TYPE_HTML_ENTITY_DECODE,
    TEXT_TRANSFORMATION_TYPE_JS_DECODE,
    TEXT_TRANSFORMATION_TYPE_MD5,
    TEXT_TRANSFORMATION_TYPE_NORMALIZE_PATH_WIN,
    TEXT_TRANSFORMATION_TYPE_REPLACE_COMMENTS,
    TEXT_TRANSFORMATION_TYPE_BASE64_DECODE_EXT,
    TEXT_TRANSFORMATION_TYPE_UTF8_TO_UNICODE,
    TEXT_TRANSFORMATION_TYPE_ESCAPE_SEQ_DECODE,
    TEXT_TRANSFORMATION_TYPE_CSS_DECODE,
    TEXT_TRANSFORMATION_TYPE_REPLACE_NULLS,
    TEXT_TRANSFORMATION_TYPE_LOWERCASE,
    TEXT_TRANSFORMATION_TYPE_CMD_LINE,
    TEXT_TRANSFORMATION_TYPE_REMOVE_NULLS,
  ];

  static final $core.List<TextTransformationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 20);
  static TextTransformationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TextTransformationType._(super.value, super.name);
}

class UsageOfAction extends $pb.ProtobufEnum {
  static const UsageOfAction USAGE_OF_ACTION_DISABLED =
      UsageOfAction._(0, _omitEnumNames ? '' : 'USAGE_OF_ACTION_DISABLED');
  static const UsageOfAction USAGE_OF_ACTION_ENABLED =
      UsageOfAction._(1, _omitEnumNames ? '' : 'USAGE_OF_ACTION_ENABLED');

  static const $core.List<UsageOfAction> values = <UsageOfAction>[
    USAGE_OF_ACTION_DISABLED,
    USAGE_OF_ACTION_ENABLED,
  ];

  static final $core.List<UsageOfAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static UsageOfAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UsageOfAction._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
