// This is a generated file - do not edit.
//
// Generated from waf.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ChangeAction extends $pb.ProtobufEnum {
  static const ChangeAction CHANGE_ACTION_INSERT =
      ChangeAction._(0, _omitEnumNames ? '' : 'CHANGE_ACTION_INSERT');
  static const ChangeAction CHANGE_ACTION_DELETE =
      ChangeAction._(1, _omitEnumNames ? '' : 'CHANGE_ACTION_DELETE');

  static const $core.List<ChangeAction> values = <ChangeAction>[
    CHANGE_ACTION_INSERT,
    CHANGE_ACTION_DELETE,
  ];

  static final $core.List<ChangeAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ChangeAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChangeAction._(super.value, super.name);
}

class ChangeTokenStatus extends $pb.ProtobufEnum {
  static const ChangeTokenStatus CHANGE_TOKEN_STATUS_PENDING =
      ChangeTokenStatus._(
          0, _omitEnumNames ? '' : 'CHANGE_TOKEN_STATUS_PENDING');
  static const ChangeTokenStatus CHANGE_TOKEN_STATUS_INSYNC =
      ChangeTokenStatus._(
          1, _omitEnumNames ? '' : 'CHANGE_TOKEN_STATUS_INSYNC');
  static const ChangeTokenStatus CHANGE_TOKEN_STATUS_PROVISIONED =
      ChangeTokenStatus._(
          2, _omitEnumNames ? '' : 'CHANGE_TOKEN_STATUS_PROVISIONED');

  static const $core.List<ChangeTokenStatus> values = <ChangeTokenStatus>[
    CHANGE_TOKEN_STATUS_PENDING,
    CHANGE_TOKEN_STATUS_INSYNC,
    CHANGE_TOKEN_STATUS_PROVISIONED,
  ];

  static final $core.List<ChangeTokenStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ChangeTokenStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChangeTokenStatus._(super.value, super.name);
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

class GeoMatchConstraintType extends $pb.ProtobufEnum {
  static const GeoMatchConstraintType GEO_MATCH_CONSTRAINT_TYPE_COUNTRY =
      GeoMatchConstraintType._(
          0, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_TYPE_COUNTRY');

  static const $core.List<GeoMatchConstraintType> values =
      <GeoMatchConstraintType>[
    GEO_MATCH_CONSTRAINT_TYPE_COUNTRY,
  ];

  static final $core.List<GeoMatchConstraintType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static GeoMatchConstraintType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GeoMatchConstraintType._(super.value, super.name);
}

class GeoMatchConstraintValue extends $pb.ProtobufEnum {
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_US =
      GeoMatchConstraintValue._(
          0, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_US');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_YE =
      GeoMatchConstraintValue._(
          1, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_YE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TO =
      GeoMatchConstraintValue._(
          2, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SI =
      GeoMatchConstraintValue._(
          3, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MA =
      GeoMatchConstraintValue._(
          4, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BY =
      GeoMatchConstraintValue._(
          5, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PK =
      GeoMatchConstraintValue._(
          6, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BB =
      GeoMatchConstraintValue._(
          7, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BB');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CO =
      GeoMatchConstraintValue._(
          8, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CZ =
      GeoMatchConstraintValue._(
          9, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_FI =
      GeoMatchConstraintValue._(
          10, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_FI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_VA =
      GeoMatchConstraintValue._(
          11, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_VA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GU =
      GeoMatchConstraintValue._(
          12, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AU =
      GeoMatchConstraintValue._(
          13, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TZ =
      GeoMatchConstraintValue._(
          14, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TG =
      GeoMatchConstraintValue._(
          15, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IT =
      GeoMatchConstraintValue._(
          16, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LS =
      GeoMatchConstraintValue._(
          17, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_RU =
      GeoMatchConstraintValue._(
          18, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_RU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SZ =
      GeoMatchConstraintValue._(
          19, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SA =
      GeoMatchConstraintValue._(
          20, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NO =
      GeoMatchConstraintValue._(
          21, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BQ =
      GeoMatchConstraintValue._(
          22, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BQ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CG =
      GeoMatchConstraintValue._(
          23, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_UM =
      GeoMatchConstraintValue._(
          24, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_UM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CR =
      GeoMatchConstraintValue._(
          25, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_VI =
      GeoMatchConstraintValue._(
          26, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_VI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GM =
      GeoMatchConstraintValue._(
          27, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_DJ =
      GeoMatchConstraintValue._(
          28, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_DJ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TR =
      GeoMatchConstraintValue._(
          29, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KN =
      GeoMatchConstraintValue._(
          30, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_EG =
      GeoMatchConstraintValue._(
          31, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_EG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_HM =
      GeoMatchConstraintValue._(
          32, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_HM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_HR =
      GeoMatchConstraintValue._(
          33, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_HR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IL =
      GeoMatchConstraintValue._(
          34, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SR =
      GeoMatchConstraintValue._(
          35, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MF =
      GeoMatchConstraintValue._(
          36, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NG =
      GeoMatchConstraintValue._(
          37, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MQ =
      GeoMatchConstraintValue._(
          38, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MQ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NP =
      GeoMatchConstraintValue._(
          39, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NP');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PF =
      GeoMatchConstraintValue._(
          40, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BI =
      GeoMatchConstraintValue._(
          41, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GE =
      GeoMatchConstraintValue._(
          42, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AR =
      GeoMatchConstraintValue._(
          43, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AE =
      GeoMatchConstraintValue._(
          44, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TJ =
      GeoMatchConstraintValue._(
          45, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TJ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ET =
      GeoMatchConstraintValue._(
          46, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ET');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_RO =
      GeoMatchConstraintValue._(
          47, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_RO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IQ =
      GeoMatchConstraintValue._(
          48, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IQ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LV =
      GeoMatchConstraintValue._(
          49, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LV');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ID =
      GeoMatchConstraintValue._(
          50, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ID');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LC =
      GeoMatchConstraintValue._(
          51, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SJ =
      GeoMatchConstraintValue._(
          52, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SJ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MN =
      GeoMatchConstraintValue._(
          53, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PY =
      GeoMatchConstraintValue._(
          54, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MY =
      GeoMatchConstraintValue._(
          55, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PN =
      GeoMatchConstraintValue._(
          56, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BA =
      GeoMatchConstraintValue._(
          57, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CH =
      GeoMatchConstraintValue._(
          58, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CW =
      GeoMatchConstraintValue._(
          59, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AZ =
      GeoMatchConstraintValue._(
          60, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_DM =
      GeoMatchConstraintValue._(
          61, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_DM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AM =
      GeoMatchConstraintValue._(
          62, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_DZ =
      GeoMatchConstraintValue._(
          63, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_DZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LK =
      GeoMatchConstraintValue._(
          64, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SB =
      GeoMatchConstraintValue._(
          65, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SB');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MV =
      GeoMatchConstraintValue._(
          66, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MV');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PA =
      GeoMatchConstraintValue._(
          67, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BT =
      GeoMatchConstraintValue._(
          68, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GN =
      GeoMatchConstraintValue._(
          69, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_DE =
      GeoMatchConstraintValue._(
          70, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_DE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TM =
      GeoMatchConstraintValue._(
          71, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_HU =
      GeoMatchConstraintValue._(
          72, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_HU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SO =
      GeoMatchConstraintValue._(
          73, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MC =
      GeoMatchConstraintValue._(
          74, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PT =
      GeoMatchConstraintValue._(
          75, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BL =
      GeoMatchConstraintValue._(
          76, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CM =
      GeoMatchConstraintValue._(
          77, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CX =
      GeoMatchConstraintValue._(
          78, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CX');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_FK =
      GeoMatchConstraintValue._(
          79, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_FK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_VC =
      GeoMatchConstraintValue._(
          80, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_VC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GS =
      GeoMatchConstraintValue._(
          81, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GF =
      GeoMatchConstraintValue._(
          82, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AW =
      GeoMatchConstraintValue._(
          83, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_RW =
      GeoMatchConstraintValue._(
          84, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_RW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SX =
      GeoMatchConstraintValue._(
          85, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SX');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SG =
      GeoMatchConstraintValue._(
          86, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MK =
      GeoMatchConstraintValue._(
          87, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NZ =
      GeoMatchConstraintValue._(
          88, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BS =
      GeoMatchConstraintValue._(
          89, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BD =
      GeoMatchConstraintValue._(
          90, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BD');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_UZ =
      GeoMatchConstraintValue._(
          91, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_UZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_JE =
      GeoMatchConstraintValue._(
          92, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_JE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_WF =
      GeoMatchConstraintValue._(
          93, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_WF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_HK =
      GeoMatchConstraintValue._(
          94, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_HK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IN =
      GeoMatchConstraintValue._(
          95, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LY =
      GeoMatchConstraintValue._(
          96, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NI =
      GeoMatchConstraintValue._(
          97, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MS =
      GeoMatchConstraintValue._(
          98, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NR =
      GeoMatchConstraintValue._(
          99, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_UG =
      GeoMatchConstraintValue._(
          100, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_UG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_JM =
      GeoMatchConstraintValue._(
          101, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_JM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AG =
      GeoMatchConstraintValue._(
          102, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KY =
      GeoMatchConstraintValue._(
          103, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TH =
      GeoMatchConstraintValue._(
          104, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IS =
      GeoMatchConstraintValue._(
          105, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LT =
      GeoMatchConstraintValue._(
          106, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LA =
      GeoMatchConstraintValue._(
          107, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SH =
      GeoMatchConstraintValue._(
          108, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NA =
      GeoMatchConstraintValue._(
          109, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PW =
      GeoMatchConstraintValue._(
          110, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PL =
      GeoMatchConstraintValue._(
          111, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CN =
      GeoMatchConstraintValue._(
          112, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_VU =
      GeoMatchConstraintValue._(
          113, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_VU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CU =
      GeoMatchConstraintValue._(
          114, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_VN =
      GeoMatchConstraintValue._(
          115, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_VN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GT =
      GeoMatchConstraintValue._(
          116, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AT =
      GeoMatchConstraintValue._(
          117, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ZM =
      GeoMatchConstraintValue._(
          118, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ZM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AO =
      GeoMatchConstraintValue._(
          119, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LI =
      GeoMatchConstraintValue._(
          120, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NL =
      GeoMatchConstraintValue._(
          121, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MH =
      GeoMatchConstraintValue._(
          122, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BV =
      GeoMatchConstraintValue._(
          123, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BV');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CF =
      GeoMatchConstraintValue._(
          124, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GY =
      GeoMatchConstraintValue._(
          125, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GL =
      GeoMatchConstraintValue._(
          126, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_DK =
      GeoMatchConstraintValue._(
          127, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_DK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KR =
      GeoMatchConstraintValue._(
          128, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_YT =
      GeoMatchConstraintValue._(
          129, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_YT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KI =
      GeoMatchConstraintValue._(
          130, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_HN =
      GeoMatchConstraintValue._(
          131, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_HN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SM =
      GeoMatchConstraintValue._(
          132, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ME =
      GeoMatchConstraintValue._(
          133, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ME');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PR =
      GeoMatchConstraintValue._(
          134, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MP =
      GeoMatchConstraintValue._(
          135, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MP');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PG =
      GeoMatchConstraintValue._(
          136, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BN =
      GeoMatchConstraintValue._(
          137, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_FM =
      GeoMatchConstraintValue._(
          138, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_FM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_VE =
      GeoMatchConstraintValue._(
          139, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_VE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GQ =
      GeoMatchConstraintValue._(
          140, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GQ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GD =
      GeoMatchConstraintValue._(
          141, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GD');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ZW =
      GeoMatchConstraintValue._(
          142, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ZW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AQ =
      GeoMatchConstraintValue._(
          143, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AQ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_JP =
      GeoMatchConstraintValue._(
          144, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_JP');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_WS =
      GeoMatchConstraintValue._(
          145, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_WS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AD =
      GeoMatchConstraintValue._(
          146, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AD');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KZ =
      GeoMatchConstraintValue._(
          147, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TK =
      GeoMatchConstraintValue._(
          148, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ES =
      GeoMatchConstraintValue._(
          149, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ES');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SE =
      GeoMatchConstraintValue._(
          150, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MM =
      GeoMatchConstraintValue._(
          151, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MX =
      GeoMatchConstraintValue._(
          152, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MX');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BF =
      GeoMatchConstraintValue._(
          153, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CK =
      GeoMatchConstraintValue._(
          154, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CV =
      GeoMatchConstraintValue._(
          155, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CV');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GI =
      GeoMatchConstraintValue._(
          156, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KW =
      GeoMatchConstraintValue._(
          157, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TV =
      GeoMatchConstraintValue._(
          158, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TV');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AL =
      GeoMatchConstraintValue._(
          159, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TC =
      GeoMatchConstraintValue._(
          160, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SV =
      GeoMatchConstraintValue._(
          161, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SV');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MU =
      GeoMatchConstraintValue._(
          162, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CC =
      GeoMatchConstraintValue._(
          163, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GA =
      GeoMatchConstraintValue._(
          164, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_JO =
      GeoMatchConstraintValue._(
          165, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_JO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TN =
      GeoMatchConstraintValue._(
          166, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_EC =
      GeoMatchConstraintValue._(
          167, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_EC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SN =
      GeoMatchConstraintValue._(
          168, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SN');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NC =
      GeoMatchConstraintValue._(
          169, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BM =
      GeoMatchConstraintValue._(
          170, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CL =
      GeoMatchConstraintValue._(
          171, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_UA =
      GeoMatchConstraintValue._(
          172, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_UA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GR =
      GeoMatchConstraintValue._(
          173, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AI =
      GeoMatchConstraintValue._(
          174, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KG =
      GeoMatchConstraintValue._(
          175, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TF =
      GeoMatchConstraintValue._(
          176, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LR =
      GeoMatchConstraintValue._(
          177, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BE =
      GeoMatchConstraintValue._(
          178, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CD =
      GeoMatchConstraintValue._(
          179, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CD');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_UY =
      GeoMatchConstraintValue._(
          180, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_UY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KP =
      GeoMatchConstraintValue._(
          181, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KP');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_EH =
      GeoMatchConstraintValue._(
          182, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_EH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IM =
      GeoMatchConstraintValue._(
          183, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SS =
      GeoMatchConstraintValue._(
          184, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MG =
      GeoMatchConstraintValue._(
          185, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NF =
      GeoMatchConstraintValue._(
          186, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MR =
      GeoMatchConstraintValue._(
          187, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PE =
      GeoMatchConstraintValue._(
          188, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BH =
      GeoMatchConstraintValue._(
          189, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_FO =
      GeoMatchConstraintValue._(
          190, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_FO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_VG =
      GeoMatchConstraintValue._(
          191, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_VG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GB =
      GeoMatchConstraintValue._(
          192, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GB');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AS =
      GeoMatchConstraintValue._(
          193, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AF =
      GeoMatchConstraintValue._(
          194, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AF');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IR =
      GeoMatchConstraintValue._(
          195, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LU =
      GeoMatchConstraintValue._(
          196, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_RS =
      GeoMatchConstraintValue._(
          197, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_RS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IE =
      GeoMatchConstraintValue._(
          198, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_LB =
      GeoMatchConstraintValue._(
          199, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_LB');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SK =
      GeoMatchConstraintValue._(
          200, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SK');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MO =
      GeoMatchConstraintValue._(
          201, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MZ =
      GeoMatchConstraintValue._(
          202, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PM =
      GeoMatchConstraintValue._(
          203, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CI =
      GeoMatchConstraintValue._(
          204, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CI');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GW =
      GeoMatchConstraintValue._(
          205, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TT =
      GeoMatchConstraintValue._(
          206, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_OM =
      GeoMatchConstraintValue._(
          207, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_OM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ST =
      GeoMatchConstraintValue._(
          208, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ST');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SC =
      GeoMatchConstraintValue._(
          209, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SC');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MW =
      GeoMatchConstraintValue._(
          210, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BW =
      GeoMatchConstraintValue._(
          211, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CA =
      GeoMatchConstraintValue._(
          212, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ZA =
      GeoMatchConstraintValue._(
          213, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ZA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TL =
      GeoMatchConstraintValue._(
          214, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KH =
      GeoMatchConstraintValue._(
          215, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_EE =
      GeoMatchConstraintValue._(
          216, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_EE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_HT =
      GeoMatchConstraintValue._(
          217, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_HT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SL =
      GeoMatchConstraintValue._(
          218, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SL');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MD =
      GeoMatchConstraintValue._(
          219, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MD');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NE =
      GeoMatchConstraintValue._(
          220, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PS =
      GeoMatchConstraintValue._(
          221, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PS');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BZ =
      GeoMatchConstraintValue._(
          222, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BZ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_PH =
      GeoMatchConstraintValue._(
          223, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_PH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BO =
      GeoMatchConstraintValue._(
          224, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_QA =
      GeoMatchConstraintValue._(
          225, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_QA');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_CY =
      GeoMatchConstraintValue._(
          226, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_CY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_FJ =
      GeoMatchConstraintValue._(
          227, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_FJ');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GP =
      GeoMatchConstraintValue._(
          228, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GP');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GG =
      GeoMatchConstraintValue._(
          229, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KE =
      GeoMatchConstraintValue._(
          230, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TD =
      GeoMatchConstraintValue._(
          231, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TD');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ER =
      GeoMatchConstraintValue._(
          232, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ER');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SY =
      GeoMatchConstraintValue._(
          233, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SY');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_SD =
      GeoMatchConstraintValue._(
          234, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_SD');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_ML =
      GeoMatchConstraintValue._(
          235, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_ML');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BR =
      GeoMatchConstraintValue._(
          236, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BG =
      GeoMatchConstraintValue._(
          237, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BG');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_FR =
      GeoMatchConstraintValue._(
          238, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_FR');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_GH =
      GeoMatchConstraintValue._(
          239, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_GH');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_AX =
      GeoMatchConstraintValue._(
          240, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_AX');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_DO =
      GeoMatchConstraintValue._(
          241, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_DO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_TW =
      GeoMatchConstraintValue._(
          242, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_TW');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_KM =
      GeoMatchConstraintValue._(
          243, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_KM');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_RE =
      GeoMatchConstraintValue._(
          244, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_RE');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_IO =
      GeoMatchConstraintValue._(
          245, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_IO');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_MT =
      GeoMatchConstraintValue._(
          246, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_MT');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_NU =
      GeoMatchConstraintValue._(
          247, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_NU');
  static const GeoMatchConstraintValue GEO_MATCH_CONSTRAINT_VALUE_BJ =
      GeoMatchConstraintValue._(
          248, _omitEnumNames ? '' : 'GEO_MATCH_CONSTRAINT_VALUE_BJ');

  static const $core.List<GeoMatchConstraintValue> values =
      <GeoMatchConstraintValue>[
    GEO_MATCH_CONSTRAINT_VALUE_US,
    GEO_MATCH_CONSTRAINT_VALUE_YE,
    GEO_MATCH_CONSTRAINT_VALUE_TO,
    GEO_MATCH_CONSTRAINT_VALUE_SI,
    GEO_MATCH_CONSTRAINT_VALUE_MA,
    GEO_MATCH_CONSTRAINT_VALUE_BY,
    GEO_MATCH_CONSTRAINT_VALUE_PK,
    GEO_MATCH_CONSTRAINT_VALUE_BB,
    GEO_MATCH_CONSTRAINT_VALUE_CO,
    GEO_MATCH_CONSTRAINT_VALUE_CZ,
    GEO_MATCH_CONSTRAINT_VALUE_FI,
    GEO_MATCH_CONSTRAINT_VALUE_VA,
    GEO_MATCH_CONSTRAINT_VALUE_GU,
    GEO_MATCH_CONSTRAINT_VALUE_AU,
    GEO_MATCH_CONSTRAINT_VALUE_TZ,
    GEO_MATCH_CONSTRAINT_VALUE_TG,
    GEO_MATCH_CONSTRAINT_VALUE_IT,
    GEO_MATCH_CONSTRAINT_VALUE_LS,
    GEO_MATCH_CONSTRAINT_VALUE_RU,
    GEO_MATCH_CONSTRAINT_VALUE_SZ,
    GEO_MATCH_CONSTRAINT_VALUE_SA,
    GEO_MATCH_CONSTRAINT_VALUE_NO,
    GEO_MATCH_CONSTRAINT_VALUE_BQ,
    GEO_MATCH_CONSTRAINT_VALUE_CG,
    GEO_MATCH_CONSTRAINT_VALUE_UM,
    GEO_MATCH_CONSTRAINT_VALUE_CR,
    GEO_MATCH_CONSTRAINT_VALUE_VI,
    GEO_MATCH_CONSTRAINT_VALUE_GM,
    GEO_MATCH_CONSTRAINT_VALUE_DJ,
    GEO_MATCH_CONSTRAINT_VALUE_TR,
    GEO_MATCH_CONSTRAINT_VALUE_KN,
    GEO_MATCH_CONSTRAINT_VALUE_EG,
    GEO_MATCH_CONSTRAINT_VALUE_HM,
    GEO_MATCH_CONSTRAINT_VALUE_HR,
    GEO_MATCH_CONSTRAINT_VALUE_IL,
    GEO_MATCH_CONSTRAINT_VALUE_SR,
    GEO_MATCH_CONSTRAINT_VALUE_MF,
    GEO_MATCH_CONSTRAINT_VALUE_NG,
    GEO_MATCH_CONSTRAINT_VALUE_MQ,
    GEO_MATCH_CONSTRAINT_VALUE_NP,
    GEO_MATCH_CONSTRAINT_VALUE_PF,
    GEO_MATCH_CONSTRAINT_VALUE_BI,
    GEO_MATCH_CONSTRAINT_VALUE_GE,
    GEO_MATCH_CONSTRAINT_VALUE_AR,
    GEO_MATCH_CONSTRAINT_VALUE_AE,
    GEO_MATCH_CONSTRAINT_VALUE_TJ,
    GEO_MATCH_CONSTRAINT_VALUE_ET,
    GEO_MATCH_CONSTRAINT_VALUE_RO,
    GEO_MATCH_CONSTRAINT_VALUE_IQ,
    GEO_MATCH_CONSTRAINT_VALUE_LV,
    GEO_MATCH_CONSTRAINT_VALUE_ID,
    GEO_MATCH_CONSTRAINT_VALUE_LC,
    GEO_MATCH_CONSTRAINT_VALUE_SJ,
    GEO_MATCH_CONSTRAINT_VALUE_MN,
    GEO_MATCH_CONSTRAINT_VALUE_PY,
    GEO_MATCH_CONSTRAINT_VALUE_MY,
    GEO_MATCH_CONSTRAINT_VALUE_PN,
    GEO_MATCH_CONSTRAINT_VALUE_BA,
    GEO_MATCH_CONSTRAINT_VALUE_CH,
    GEO_MATCH_CONSTRAINT_VALUE_CW,
    GEO_MATCH_CONSTRAINT_VALUE_AZ,
    GEO_MATCH_CONSTRAINT_VALUE_DM,
    GEO_MATCH_CONSTRAINT_VALUE_AM,
    GEO_MATCH_CONSTRAINT_VALUE_DZ,
    GEO_MATCH_CONSTRAINT_VALUE_LK,
    GEO_MATCH_CONSTRAINT_VALUE_SB,
    GEO_MATCH_CONSTRAINT_VALUE_MV,
    GEO_MATCH_CONSTRAINT_VALUE_PA,
    GEO_MATCH_CONSTRAINT_VALUE_BT,
    GEO_MATCH_CONSTRAINT_VALUE_GN,
    GEO_MATCH_CONSTRAINT_VALUE_DE,
    GEO_MATCH_CONSTRAINT_VALUE_TM,
    GEO_MATCH_CONSTRAINT_VALUE_HU,
    GEO_MATCH_CONSTRAINT_VALUE_SO,
    GEO_MATCH_CONSTRAINT_VALUE_MC,
    GEO_MATCH_CONSTRAINT_VALUE_PT,
    GEO_MATCH_CONSTRAINT_VALUE_BL,
    GEO_MATCH_CONSTRAINT_VALUE_CM,
    GEO_MATCH_CONSTRAINT_VALUE_CX,
    GEO_MATCH_CONSTRAINT_VALUE_FK,
    GEO_MATCH_CONSTRAINT_VALUE_VC,
    GEO_MATCH_CONSTRAINT_VALUE_GS,
    GEO_MATCH_CONSTRAINT_VALUE_GF,
    GEO_MATCH_CONSTRAINT_VALUE_AW,
    GEO_MATCH_CONSTRAINT_VALUE_RW,
    GEO_MATCH_CONSTRAINT_VALUE_SX,
    GEO_MATCH_CONSTRAINT_VALUE_SG,
    GEO_MATCH_CONSTRAINT_VALUE_MK,
    GEO_MATCH_CONSTRAINT_VALUE_NZ,
    GEO_MATCH_CONSTRAINT_VALUE_BS,
    GEO_MATCH_CONSTRAINT_VALUE_BD,
    GEO_MATCH_CONSTRAINT_VALUE_UZ,
    GEO_MATCH_CONSTRAINT_VALUE_JE,
    GEO_MATCH_CONSTRAINT_VALUE_WF,
    GEO_MATCH_CONSTRAINT_VALUE_HK,
    GEO_MATCH_CONSTRAINT_VALUE_IN,
    GEO_MATCH_CONSTRAINT_VALUE_LY,
    GEO_MATCH_CONSTRAINT_VALUE_NI,
    GEO_MATCH_CONSTRAINT_VALUE_MS,
    GEO_MATCH_CONSTRAINT_VALUE_NR,
    GEO_MATCH_CONSTRAINT_VALUE_UG,
    GEO_MATCH_CONSTRAINT_VALUE_JM,
    GEO_MATCH_CONSTRAINT_VALUE_AG,
    GEO_MATCH_CONSTRAINT_VALUE_KY,
    GEO_MATCH_CONSTRAINT_VALUE_TH,
    GEO_MATCH_CONSTRAINT_VALUE_IS,
    GEO_MATCH_CONSTRAINT_VALUE_LT,
    GEO_MATCH_CONSTRAINT_VALUE_LA,
    GEO_MATCH_CONSTRAINT_VALUE_SH,
    GEO_MATCH_CONSTRAINT_VALUE_NA,
    GEO_MATCH_CONSTRAINT_VALUE_PW,
    GEO_MATCH_CONSTRAINT_VALUE_PL,
    GEO_MATCH_CONSTRAINT_VALUE_CN,
    GEO_MATCH_CONSTRAINT_VALUE_VU,
    GEO_MATCH_CONSTRAINT_VALUE_CU,
    GEO_MATCH_CONSTRAINT_VALUE_VN,
    GEO_MATCH_CONSTRAINT_VALUE_GT,
    GEO_MATCH_CONSTRAINT_VALUE_AT,
    GEO_MATCH_CONSTRAINT_VALUE_ZM,
    GEO_MATCH_CONSTRAINT_VALUE_AO,
    GEO_MATCH_CONSTRAINT_VALUE_LI,
    GEO_MATCH_CONSTRAINT_VALUE_NL,
    GEO_MATCH_CONSTRAINT_VALUE_MH,
    GEO_MATCH_CONSTRAINT_VALUE_BV,
    GEO_MATCH_CONSTRAINT_VALUE_CF,
    GEO_MATCH_CONSTRAINT_VALUE_GY,
    GEO_MATCH_CONSTRAINT_VALUE_GL,
    GEO_MATCH_CONSTRAINT_VALUE_DK,
    GEO_MATCH_CONSTRAINT_VALUE_KR,
    GEO_MATCH_CONSTRAINT_VALUE_YT,
    GEO_MATCH_CONSTRAINT_VALUE_KI,
    GEO_MATCH_CONSTRAINT_VALUE_HN,
    GEO_MATCH_CONSTRAINT_VALUE_SM,
    GEO_MATCH_CONSTRAINT_VALUE_ME,
    GEO_MATCH_CONSTRAINT_VALUE_PR,
    GEO_MATCH_CONSTRAINT_VALUE_MP,
    GEO_MATCH_CONSTRAINT_VALUE_PG,
    GEO_MATCH_CONSTRAINT_VALUE_BN,
    GEO_MATCH_CONSTRAINT_VALUE_FM,
    GEO_MATCH_CONSTRAINT_VALUE_VE,
    GEO_MATCH_CONSTRAINT_VALUE_GQ,
    GEO_MATCH_CONSTRAINT_VALUE_GD,
    GEO_MATCH_CONSTRAINT_VALUE_ZW,
    GEO_MATCH_CONSTRAINT_VALUE_AQ,
    GEO_MATCH_CONSTRAINT_VALUE_JP,
    GEO_MATCH_CONSTRAINT_VALUE_WS,
    GEO_MATCH_CONSTRAINT_VALUE_AD,
    GEO_MATCH_CONSTRAINT_VALUE_KZ,
    GEO_MATCH_CONSTRAINT_VALUE_TK,
    GEO_MATCH_CONSTRAINT_VALUE_ES,
    GEO_MATCH_CONSTRAINT_VALUE_SE,
    GEO_MATCH_CONSTRAINT_VALUE_MM,
    GEO_MATCH_CONSTRAINT_VALUE_MX,
    GEO_MATCH_CONSTRAINT_VALUE_BF,
    GEO_MATCH_CONSTRAINT_VALUE_CK,
    GEO_MATCH_CONSTRAINT_VALUE_CV,
    GEO_MATCH_CONSTRAINT_VALUE_GI,
    GEO_MATCH_CONSTRAINT_VALUE_KW,
    GEO_MATCH_CONSTRAINT_VALUE_TV,
    GEO_MATCH_CONSTRAINT_VALUE_AL,
    GEO_MATCH_CONSTRAINT_VALUE_TC,
    GEO_MATCH_CONSTRAINT_VALUE_SV,
    GEO_MATCH_CONSTRAINT_VALUE_MU,
    GEO_MATCH_CONSTRAINT_VALUE_CC,
    GEO_MATCH_CONSTRAINT_VALUE_GA,
    GEO_MATCH_CONSTRAINT_VALUE_JO,
    GEO_MATCH_CONSTRAINT_VALUE_TN,
    GEO_MATCH_CONSTRAINT_VALUE_EC,
    GEO_MATCH_CONSTRAINT_VALUE_SN,
    GEO_MATCH_CONSTRAINT_VALUE_NC,
    GEO_MATCH_CONSTRAINT_VALUE_BM,
    GEO_MATCH_CONSTRAINT_VALUE_CL,
    GEO_MATCH_CONSTRAINT_VALUE_UA,
    GEO_MATCH_CONSTRAINT_VALUE_GR,
    GEO_MATCH_CONSTRAINT_VALUE_AI,
    GEO_MATCH_CONSTRAINT_VALUE_KG,
    GEO_MATCH_CONSTRAINT_VALUE_TF,
    GEO_MATCH_CONSTRAINT_VALUE_LR,
    GEO_MATCH_CONSTRAINT_VALUE_BE,
    GEO_MATCH_CONSTRAINT_VALUE_CD,
    GEO_MATCH_CONSTRAINT_VALUE_UY,
    GEO_MATCH_CONSTRAINT_VALUE_KP,
    GEO_MATCH_CONSTRAINT_VALUE_EH,
    GEO_MATCH_CONSTRAINT_VALUE_IM,
    GEO_MATCH_CONSTRAINT_VALUE_SS,
    GEO_MATCH_CONSTRAINT_VALUE_MG,
    GEO_MATCH_CONSTRAINT_VALUE_NF,
    GEO_MATCH_CONSTRAINT_VALUE_MR,
    GEO_MATCH_CONSTRAINT_VALUE_PE,
    GEO_MATCH_CONSTRAINT_VALUE_BH,
    GEO_MATCH_CONSTRAINT_VALUE_FO,
    GEO_MATCH_CONSTRAINT_VALUE_VG,
    GEO_MATCH_CONSTRAINT_VALUE_GB,
    GEO_MATCH_CONSTRAINT_VALUE_AS,
    GEO_MATCH_CONSTRAINT_VALUE_AF,
    GEO_MATCH_CONSTRAINT_VALUE_IR,
    GEO_MATCH_CONSTRAINT_VALUE_LU,
    GEO_MATCH_CONSTRAINT_VALUE_RS,
    GEO_MATCH_CONSTRAINT_VALUE_IE,
    GEO_MATCH_CONSTRAINT_VALUE_LB,
    GEO_MATCH_CONSTRAINT_VALUE_SK,
    GEO_MATCH_CONSTRAINT_VALUE_MO,
    GEO_MATCH_CONSTRAINT_VALUE_MZ,
    GEO_MATCH_CONSTRAINT_VALUE_PM,
    GEO_MATCH_CONSTRAINT_VALUE_CI,
    GEO_MATCH_CONSTRAINT_VALUE_GW,
    GEO_MATCH_CONSTRAINT_VALUE_TT,
    GEO_MATCH_CONSTRAINT_VALUE_OM,
    GEO_MATCH_CONSTRAINT_VALUE_ST,
    GEO_MATCH_CONSTRAINT_VALUE_SC,
    GEO_MATCH_CONSTRAINT_VALUE_MW,
    GEO_MATCH_CONSTRAINT_VALUE_BW,
    GEO_MATCH_CONSTRAINT_VALUE_CA,
    GEO_MATCH_CONSTRAINT_VALUE_ZA,
    GEO_MATCH_CONSTRAINT_VALUE_TL,
    GEO_MATCH_CONSTRAINT_VALUE_KH,
    GEO_MATCH_CONSTRAINT_VALUE_EE,
    GEO_MATCH_CONSTRAINT_VALUE_HT,
    GEO_MATCH_CONSTRAINT_VALUE_SL,
    GEO_MATCH_CONSTRAINT_VALUE_MD,
    GEO_MATCH_CONSTRAINT_VALUE_NE,
    GEO_MATCH_CONSTRAINT_VALUE_PS,
    GEO_MATCH_CONSTRAINT_VALUE_BZ,
    GEO_MATCH_CONSTRAINT_VALUE_PH,
    GEO_MATCH_CONSTRAINT_VALUE_BO,
    GEO_MATCH_CONSTRAINT_VALUE_QA,
    GEO_MATCH_CONSTRAINT_VALUE_CY,
    GEO_MATCH_CONSTRAINT_VALUE_FJ,
    GEO_MATCH_CONSTRAINT_VALUE_GP,
    GEO_MATCH_CONSTRAINT_VALUE_GG,
    GEO_MATCH_CONSTRAINT_VALUE_KE,
    GEO_MATCH_CONSTRAINT_VALUE_TD,
    GEO_MATCH_CONSTRAINT_VALUE_ER,
    GEO_MATCH_CONSTRAINT_VALUE_SY,
    GEO_MATCH_CONSTRAINT_VALUE_SD,
    GEO_MATCH_CONSTRAINT_VALUE_ML,
    GEO_MATCH_CONSTRAINT_VALUE_BR,
    GEO_MATCH_CONSTRAINT_VALUE_BG,
    GEO_MATCH_CONSTRAINT_VALUE_FR,
    GEO_MATCH_CONSTRAINT_VALUE_GH,
    GEO_MATCH_CONSTRAINT_VALUE_AX,
    GEO_MATCH_CONSTRAINT_VALUE_DO,
    GEO_MATCH_CONSTRAINT_VALUE_TW,
    GEO_MATCH_CONSTRAINT_VALUE_KM,
    GEO_MATCH_CONSTRAINT_VALUE_RE,
    GEO_MATCH_CONSTRAINT_VALUE_IO,
    GEO_MATCH_CONSTRAINT_VALUE_MT,
    GEO_MATCH_CONSTRAINT_VALUE_NU,
    GEO_MATCH_CONSTRAINT_VALUE_BJ,
  ];

  static final $core.List<GeoMatchConstraintValue?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 248);
  static GeoMatchConstraintValue? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GeoMatchConstraintValue._(super.value, super.name);
}

class IPSetDescriptorType extends $pb.ProtobufEnum {
  static const IPSetDescriptorType I_P_SET_DESCRIPTOR_TYPE_IPV6 =
      IPSetDescriptorType._(
          0, _omitEnumNames ? '' : 'I_P_SET_DESCRIPTOR_TYPE_IPV6');
  static const IPSetDescriptorType I_P_SET_DESCRIPTOR_TYPE_IPV4 =
      IPSetDescriptorType._(
          1, _omitEnumNames ? '' : 'I_P_SET_DESCRIPTOR_TYPE_IPV4');

  static const $core.List<IPSetDescriptorType> values = <IPSetDescriptorType>[
    I_P_SET_DESCRIPTOR_TYPE_IPV6,
    I_P_SET_DESCRIPTOR_TYPE_IPV4,
  ];

  static final $core.List<IPSetDescriptorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IPSetDescriptorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IPSetDescriptorType._(super.value, super.name);
}

class MatchFieldType extends $pb.ProtobufEnum {
  static const MatchFieldType MATCH_FIELD_TYPE_URI =
      MatchFieldType._(0, _omitEnumNames ? '' : 'MATCH_FIELD_TYPE_URI');
  static const MatchFieldType MATCH_FIELD_TYPE_HEADER =
      MatchFieldType._(1, _omitEnumNames ? '' : 'MATCH_FIELD_TYPE_HEADER');
  static const MatchFieldType MATCH_FIELD_TYPE_QUERY_STRING = MatchFieldType._(
      2, _omitEnumNames ? '' : 'MATCH_FIELD_TYPE_QUERY_STRING');
  static const MatchFieldType MATCH_FIELD_TYPE_BODY =
      MatchFieldType._(3, _omitEnumNames ? '' : 'MATCH_FIELD_TYPE_BODY');
  static const MatchFieldType MATCH_FIELD_TYPE_ALL_QUERY_ARGS =
      MatchFieldType._(
          4, _omitEnumNames ? '' : 'MATCH_FIELD_TYPE_ALL_QUERY_ARGS');
  static const MatchFieldType MATCH_FIELD_TYPE_METHOD =
      MatchFieldType._(5, _omitEnumNames ? '' : 'MATCH_FIELD_TYPE_METHOD');
  static const MatchFieldType MATCH_FIELD_TYPE_SINGLE_QUERY_ARG =
      MatchFieldType._(
          6, _omitEnumNames ? '' : 'MATCH_FIELD_TYPE_SINGLE_QUERY_ARG');

  static const $core.List<MatchFieldType> values = <MatchFieldType>[
    MATCH_FIELD_TYPE_URI,
    MATCH_FIELD_TYPE_HEADER,
    MATCH_FIELD_TYPE_QUERY_STRING,
    MATCH_FIELD_TYPE_BODY,
    MATCH_FIELD_TYPE_ALL_QUERY_ARGS,
    MATCH_FIELD_TYPE_METHOD,
    MATCH_FIELD_TYPE_SINGLE_QUERY_ARG,
  ];

  static final $core.List<MatchFieldType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static MatchFieldType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MatchFieldType._(super.value, super.name);
}

class MigrationErrorType extends $pb.ProtobufEnum {
  static const MigrationErrorType MIGRATION_ERROR_TYPE_S3_BUCKET_NO_PERMISSION =
      MigrationErrorType._(0,
          _omitEnumNames ? '' : 'MIGRATION_ERROR_TYPE_S3_BUCKET_NO_PERMISSION');
  static const MigrationErrorType
      MIGRATION_ERROR_TYPE_S3_BUCKET_NOT_ACCESSIBLE = MigrationErrorType._(
          1,
          _omitEnumNames
              ? ''
              : 'MIGRATION_ERROR_TYPE_S3_BUCKET_NOT_ACCESSIBLE');
  static const MigrationErrorType MIGRATION_ERROR_TYPE_ENTITY_NOT_FOUND =
      MigrationErrorType._(
          2, _omitEnumNames ? '' : 'MIGRATION_ERROR_TYPE_ENTITY_NOT_FOUND');
  static const MigrationErrorType MIGRATION_ERROR_TYPE_S3_BUCKET_NOT_FOUND =
      MigrationErrorType._(
          3, _omitEnumNames ? '' : 'MIGRATION_ERROR_TYPE_S3_BUCKET_NOT_FOUND');
  static const MigrationErrorType MIGRATION_ERROR_TYPE_S3_INTERNAL_ERROR =
      MigrationErrorType._(
          4, _omitEnumNames ? '' : 'MIGRATION_ERROR_TYPE_S3_INTERNAL_ERROR');
  static const MigrationErrorType MIGRATION_ERROR_TYPE_ENTITY_NOT_SUPPORTED =
      MigrationErrorType._(
          5, _omitEnumNames ? '' : 'MIGRATION_ERROR_TYPE_ENTITY_NOT_SUPPORTED');
  static const MigrationErrorType
      MIGRATION_ERROR_TYPE_S3_BUCKET_INVALID_REGION = MigrationErrorType._(
          6,
          _omitEnumNames
              ? ''
              : 'MIGRATION_ERROR_TYPE_S3_BUCKET_INVALID_REGION');

  static const $core.List<MigrationErrorType> values = <MigrationErrorType>[
    MIGRATION_ERROR_TYPE_S3_BUCKET_NO_PERMISSION,
    MIGRATION_ERROR_TYPE_S3_BUCKET_NOT_ACCESSIBLE,
    MIGRATION_ERROR_TYPE_ENTITY_NOT_FOUND,
    MIGRATION_ERROR_TYPE_S3_BUCKET_NOT_FOUND,
    MIGRATION_ERROR_TYPE_S3_INTERNAL_ERROR,
    MIGRATION_ERROR_TYPE_ENTITY_NOT_SUPPORTED,
    MIGRATION_ERROR_TYPE_S3_BUCKET_INVALID_REGION,
  ];

  static final $core.List<MigrationErrorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static MigrationErrorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MigrationErrorType._(super.value, super.name);
}

class ParameterExceptionField extends $pb.ProtobufEnum {
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_POSITIONAL_CONSTRAINT =
      ParameterExceptionField._(
          0,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_POSITIONAL_CONSTRAINT');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_CHANGE_ACTION =
      ParameterExceptionField._(
          1, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_CHANGE_ACTION');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_NEXT_MARKER =
      ParameterExceptionField._(
          2, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_NEXT_MARKER');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_WAF_ACTION =
      ParameterExceptionField._(
          3, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_WAF_ACTION');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_PREDICATE_TYPE = ParameterExceptionField._(
          4, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_PREDICATE_TYPE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_SIZE_CONSTRAINT_COMPARISON_OPERATOR =
      ParameterExceptionField._(
          5,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_SIZE_CONSTRAINT_COMPARISON_OPERATOR');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_IPSET_TYPE =
      ParameterExceptionField._(
          6, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_IPSET_TYPE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_SQL_INJECTION_MATCH_FIELD_TYPE =
      ParameterExceptionField._(
          7,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_SQL_INJECTION_MATCH_FIELD_TYPE');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RULE_TYPE =
      ParameterExceptionField._(
          8, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RULE_TYPE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_FIELD_TYPE =
      ParameterExceptionField._(
          9,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_FIELD_TYPE');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_TAGS =
      ParameterExceptionField._(
          10, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_TAGS');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_GEO_MATCH_LOCATION_VALUE =
      ParameterExceptionField._(
          11,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_GEO_MATCH_LOCATION_VALUE');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_WAF_OVERRIDE_ACTION = ParameterExceptionField._(
          12,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_WAF_OVERRIDE_ACTION');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RESOURCE_ARN =
      ParameterExceptionField._(
          13, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RESOURCE_ARN');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_RATE_KEY =
      ParameterExceptionField._(
          14, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_RATE_KEY');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_TEXT_TRANSFORMATION =
      ParameterExceptionField._(
          15,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_TEXT_TRANSFORMATION');
  static const ParameterExceptionField PARAMETER_EXCEPTION_FIELD_TAG_KEYS =
      ParameterExceptionField._(
          16, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_FIELD_TAG_KEYS');
  static const ParameterExceptionField
      PARAMETER_EXCEPTION_FIELD_GEO_MATCH_LOCATION_TYPE =
      ParameterExceptionField._(
          17,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_FIELD_GEO_MATCH_LOCATION_TYPE');

  static const $core.List<ParameterExceptionField> values =
      <ParameterExceptionField>[
    PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_POSITIONAL_CONSTRAINT,
    PARAMETER_EXCEPTION_FIELD_CHANGE_ACTION,
    PARAMETER_EXCEPTION_FIELD_NEXT_MARKER,
    PARAMETER_EXCEPTION_FIELD_WAF_ACTION,
    PARAMETER_EXCEPTION_FIELD_PREDICATE_TYPE,
    PARAMETER_EXCEPTION_FIELD_SIZE_CONSTRAINT_COMPARISON_OPERATOR,
    PARAMETER_EXCEPTION_FIELD_IPSET_TYPE,
    PARAMETER_EXCEPTION_FIELD_SQL_INJECTION_MATCH_FIELD_TYPE,
    PARAMETER_EXCEPTION_FIELD_RULE_TYPE,
    PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_FIELD_TYPE,
    PARAMETER_EXCEPTION_FIELD_TAGS,
    PARAMETER_EXCEPTION_FIELD_GEO_MATCH_LOCATION_VALUE,
    PARAMETER_EXCEPTION_FIELD_WAF_OVERRIDE_ACTION,
    PARAMETER_EXCEPTION_FIELD_RESOURCE_ARN,
    PARAMETER_EXCEPTION_FIELD_RATE_KEY,
    PARAMETER_EXCEPTION_FIELD_BYTE_MATCH_TEXT_TRANSFORMATION,
    PARAMETER_EXCEPTION_FIELD_TAG_KEYS,
    PARAMETER_EXCEPTION_FIELD_GEO_MATCH_LOCATION_TYPE,
  ];

  static final $core.List<ParameterExceptionField?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 17);
  static ParameterExceptionField? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ParameterExceptionField._(super.value, super.name);
}

class ParameterExceptionReason extends $pb.ProtobufEnum {
  static const ParameterExceptionReason
      PARAMETER_EXCEPTION_REASON_INVALID_OPTION = ParameterExceptionReason._(
          0, _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_REASON_INVALID_OPTION');
  static const ParameterExceptionReason
      PARAMETER_EXCEPTION_REASON_ILLEGAL_ARGUMENT = ParameterExceptionReason._(
          1,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_REASON_ILLEGAL_ARGUMENT');
  static const ParameterExceptionReason
      PARAMETER_EXCEPTION_REASON_ILLEGAL_COMBINATION =
      ParameterExceptionReason._(
          2,
          _omitEnumNames
              ? ''
              : 'PARAMETER_EXCEPTION_REASON_ILLEGAL_COMBINATION');
  static const ParameterExceptionReason
      PARAMETER_EXCEPTION_REASON_INVALID_TAG_KEY = ParameterExceptionReason._(3,
          _omitEnumNames ? '' : 'PARAMETER_EXCEPTION_REASON_INVALID_TAG_KEY');

  static const $core.List<ParameterExceptionReason> values =
      <ParameterExceptionReason>[
    PARAMETER_EXCEPTION_REASON_INVALID_OPTION,
    PARAMETER_EXCEPTION_REASON_ILLEGAL_ARGUMENT,
    PARAMETER_EXCEPTION_REASON_ILLEGAL_COMBINATION,
    PARAMETER_EXCEPTION_REASON_INVALID_TAG_KEY,
  ];

  static final $core.List<ParameterExceptionReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ParameterExceptionReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ParameterExceptionReason._(super.value, super.name);
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

class PredicateType extends $pb.ProtobufEnum {
  static const PredicateType PREDICATE_TYPE_IP_MATCH =
      PredicateType._(0, _omitEnumNames ? '' : 'PREDICATE_TYPE_IP_MATCH');
  static const PredicateType PREDICATE_TYPE_BYTE_MATCH =
      PredicateType._(1, _omitEnumNames ? '' : 'PREDICATE_TYPE_BYTE_MATCH');
  static const PredicateType PREDICATE_TYPE_SQL_INJECTION_MATCH =
      PredicateType._(
          2, _omitEnumNames ? '' : 'PREDICATE_TYPE_SQL_INJECTION_MATCH');
  static const PredicateType PREDICATE_TYPE_REGEX_MATCH =
      PredicateType._(3, _omitEnumNames ? '' : 'PREDICATE_TYPE_REGEX_MATCH');
  static const PredicateType PREDICATE_TYPE_GEO_MATCH =
      PredicateType._(4, _omitEnumNames ? '' : 'PREDICATE_TYPE_GEO_MATCH');
  static const PredicateType PREDICATE_TYPE_SIZE_CONSTRAINT = PredicateType._(
      5, _omitEnumNames ? '' : 'PREDICATE_TYPE_SIZE_CONSTRAINT');
  static const PredicateType PREDICATE_TYPE_XSS_MATCH =
      PredicateType._(6, _omitEnumNames ? '' : 'PREDICATE_TYPE_XSS_MATCH');

  static const $core.List<PredicateType> values = <PredicateType>[
    PREDICATE_TYPE_IP_MATCH,
    PREDICATE_TYPE_BYTE_MATCH,
    PREDICATE_TYPE_SQL_INJECTION_MATCH,
    PREDICATE_TYPE_REGEX_MATCH,
    PREDICATE_TYPE_GEO_MATCH,
    PREDICATE_TYPE_SIZE_CONSTRAINT,
    PREDICATE_TYPE_XSS_MATCH,
  ];

  static final $core.List<PredicateType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static PredicateType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PredicateType._(super.value, super.name);
}

class RateKey extends $pb.ProtobufEnum {
  static const RateKey RATE_KEY_IP =
      RateKey._(0, _omitEnumNames ? '' : 'RATE_KEY_IP');

  static const $core.List<RateKey> values = <RateKey>[
    RATE_KEY_IP,
  ];

  static final $core.List<RateKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static RateKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RateKey._(super.value, super.name);
}

class TextTransformation extends $pb.ProtobufEnum {
  static const TextTransformation TEXT_TRANSFORMATION_URL_DECODE =
      TextTransformation._(
          0, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_URL_DECODE');
  static const TextTransformation TEXT_TRANSFORMATION_COMPRESS_WHITE_SPACE =
      TextTransformation._(
          1, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_COMPRESS_WHITE_SPACE');
  static const TextTransformation TEXT_TRANSFORMATION_NONE =
      TextTransformation._(2, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_NONE');
  static const TextTransformation TEXT_TRANSFORMATION_HTML_ENTITY_DECODE =
      TextTransformation._(
          3, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_HTML_ENTITY_DECODE');
  static const TextTransformation TEXT_TRANSFORMATION_LOWERCASE =
      TextTransformation._(
          4, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_LOWERCASE');
  static const TextTransformation TEXT_TRANSFORMATION_CMD_LINE =
      TextTransformation._(
          5, _omitEnumNames ? '' : 'TEXT_TRANSFORMATION_CMD_LINE');

  static const $core.List<TextTransformation> values = <TextTransformation>[
    TEXT_TRANSFORMATION_URL_DECODE,
    TEXT_TRANSFORMATION_COMPRESS_WHITE_SPACE,
    TEXT_TRANSFORMATION_NONE,
    TEXT_TRANSFORMATION_HTML_ENTITY_DECODE,
    TEXT_TRANSFORMATION_LOWERCASE,
    TEXT_TRANSFORMATION_CMD_LINE,
  ];

  static final $core.List<TextTransformation?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static TextTransformation? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TextTransformation._(super.value, super.name);
}

class WafActionType extends $pb.ProtobufEnum {
  static const WafActionType WAF_ACTION_TYPE_COUNT =
      WafActionType._(0, _omitEnumNames ? '' : 'WAF_ACTION_TYPE_COUNT');
  static const WafActionType WAF_ACTION_TYPE_BLOCK =
      WafActionType._(1, _omitEnumNames ? '' : 'WAF_ACTION_TYPE_BLOCK');
  static const WafActionType WAF_ACTION_TYPE_ALLOW =
      WafActionType._(2, _omitEnumNames ? '' : 'WAF_ACTION_TYPE_ALLOW');

  static const $core.List<WafActionType> values = <WafActionType>[
    WAF_ACTION_TYPE_COUNT,
    WAF_ACTION_TYPE_BLOCK,
    WAF_ACTION_TYPE_ALLOW,
  ];

  static final $core.List<WafActionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static WafActionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WafActionType._(super.value, super.name);
}

class WafOverrideActionType extends $pb.ProtobufEnum {
  static const WafOverrideActionType WAF_OVERRIDE_ACTION_TYPE_NONE =
      WafOverrideActionType._(
          0, _omitEnumNames ? '' : 'WAF_OVERRIDE_ACTION_TYPE_NONE');
  static const WafOverrideActionType WAF_OVERRIDE_ACTION_TYPE_COUNT =
      WafOverrideActionType._(
          1, _omitEnumNames ? '' : 'WAF_OVERRIDE_ACTION_TYPE_COUNT');

  static const $core.List<WafOverrideActionType> values =
      <WafOverrideActionType>[
    WAF_OVERRIDE_ACTION_TYPE_NONE,
    WAF_OVERRIDE_ACTION_TYPE_COUNT,
  ];

  static final $core.List<WafOverrideActionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static WafOverrideActionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WafOverrideActionType._(super.value, super.name);
}

class WafRuleType extends $pb.ProtobufEnum {
  static const WafRuleType WAF_RULE_TYPE_GROUP =
      WafRuleType._(0, _omitEnumNames ? '' : 'WAF_RULE_TYPE_GROUP');
  static const WafRuleType WAF_RULE_TYPE_RATE_BASED =
      WafRuleType._(1, _omitEnumNames ? '' : 'WAF_RULE_TYPE_RATE_BASED');
  static const WafRuleType WAF_RULE_TYPE_REGULAR =
      WafRuleType._(2, _omitEnumNames ? '' : 'WAF_RULE_TYPE_REGULAR');

  static const $core.List<WafRuleType> values = <WafRuleType>[
    WAF_RULE_TYPE_GROUP,
    WAF_RULE_TYPE_RATE_BASED,
    WAF_RULE_TYPE_REGULAR,
  ];

  static final $core.List<WafRuleType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static WafRuleType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WafRuleType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
