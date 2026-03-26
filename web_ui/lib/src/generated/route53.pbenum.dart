// This is a generated file - do not edit.
//
// Generated from route53.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AcceleratedRecoveryStatus extends $pb.ProtobufEnum {
  static const AcceleratedRecoveryStatus ACCELERATED_RECOVERY_STATUS_DISABLED =
      AcceleratedRecoveryStatus._(
          0, _omitEnumNames ? '' : 'ACCELERATED_RECOVERY_STATUS_DISABLED');
  static const AcceleratedRecoveryStatus
      ACCELERATED_RECOVERY_STATUS_DISABLE_FAILED = AcceleratedRecoveryStatus._(
          1,
          _omitEnumNames ? '' : 'ACCELERATED_RECOVERY_STATUS_DISABLE_FAILED');
  static const AcceleratedRecoveryStatus
      ACCELERATED_RECOVERY_STATUS_DISABLING_HOSTED_ZONE_LOCKED =
      AcceleratedRecoveryStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'ACCELERATED_RECOVERY_STATUS_DISABLING_HOSTED_ZONE_LOCKED');
  static const AcceleratedRecoveryStatus ACCELERATED_RECOVERY_STATUS_ENABLING =
      AcceleratedRecoveryStatus._(
          3, _omitEnumNames ? '' : 'ACCELERATED_RECOVERY_STATUS_ENABLING');
  static const AcceleratedRecoveryStatus
      ACCELERATED_RECOVERY_STATUS_ENABLE_FAILED = AcceleratedRecoveryStatus._(
          4, _omitEnumNames ? '' : 'ACCELERATED_RECOVERY_STATUS_ENABLE_FAILED');
  static const AcceleratedRecoveryStatus
      ACCELERATED_RECOVERY_STATUS_ENABLING_HOSTED_ZONE_LOCKED =
      AcceleratedRecoveryStatus._(
          5,
          _omitEnumNames
              ? ''
              : 'ACCELERATED_RECOVERY_STATUS_ENABLING_HOSTED_ZONE_LOCKED');
  static const AcceleratedRecoveryStatus ACCELERATED_RECOVERY_STATUS_DISABLING =
      AcceleratedRecoveryStatus._(
          6, _omitEnumNames ? '' : 'ACCELERATED_RECOVERY_STATUS_DISABLING');
  static const AcceleratedRecoveryStatus ACCELERATED_RECOVERY_STATUS_ENABLED =
      AcceleratedRecoveryStatus._(
          7, _omitEnumNames ? '' : 'ACCELERATED_RECOVERY_STATUS_ENABLED');

  static const $core.List<AcceleratedRecoveryStatus> values =
      <AcceleratedRecoveryStatus>[
    ACCELERATED_RECOVERY_STATUS_DISABLED,
    ACCELERATED_RECOVERY_STATUS_DISABLE_FAILED,
    ACCELERATED_RECOVERY_STATUS_DISABLING_HOSTED_ZONE_LOCKED,
    ACCELERATED_RECOVERY_STATUS_ENABLING,
    ACCELERATED_RECOVERY_STATUS_ENABLE_FAILED,
    ACCELERATED_RECOVERY_STATUS_ENABLING_HOSTED_ZONE_LOCKED,
    ACCELERATED_RECOVERY_STATUS_DISABLING,
    ACCELERATED_RECOVERY_STATUS_ENABLED,
  ];

  static final $core.List<AcceleratedRecoveryStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static AcceleratedRecoveryStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AcceleratedRecoveryStatus._(super.value, super.name);
}

class AccountLimitType extends $pb.ProtobufEnum {
  static const AccountLimitType
      ACCOUNT_LIMIT_TYPE_MAX_REUSABLE_DELEGATION_SETS_BY_OWNER =
      AccountLimitType._(
          0,
          _omitEnumNames
              ? ''
              : 'ACCOUNT_LIMIT_TYPE_MAX_REUSABLE_DELEGATION_SETS_BY_OWNER');
  static const AccountLimitType
      ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICY_INSTANCES_BY_OWNER =
      AccountLimitType._(
          1,
          _omitEnumNames
              ? ''
              : 'ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICY_INSTANCES_BY_OWNER');
  static const AccountLimitType ACCOUNT_LIMIT_TYPE_MAX_HEALTH_CHECKS_BY_OWNER =
      AccountLimitType._(
          2,
          _omitEnumNames
              ? ''
              : 'ACCOUNT_LIMIT_TYPE_MAX_HEALTH_CHECKS_BY_OWNER');
  static const AccountLimitType ACCOUNT_LIMIT_TYPE_MAX_HOSTED_ZONES_BY_OWNER =
      AccountLimitType._(3,
          _omitEnumNames ? '' : 'ACCOUNT_LIMIT_TYPE_MAX_HOSTED_ZONES_BY_OWNER');
  static const AccountLimitType
      ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICIES_BY_OWNER = AccountLimitType._(
          4,
          _omitEnumNames
              ? ''
              : 'ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICIES_BY_OWNER');

  static const $core.List<AccountLimitType> values = <AccountLimitType>[
    ACCOUNT_LIMIT_TYPE_MAX_REUSABLE_DELEGATION_SETS_BY_OWNER,
    ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICY_INSTANCES_BY_OWNER,
    ACCOUNT_LIMIT_TYPE_MAX_HEALTH_CHECKS_BY_OWNER,
    ACCOUNT_LIMIT_TYPE_MAX_HOSTED_ZONES_BY_OWNER,
    ACCOUNT_LIMIT_TYPE_MAX_TRAFFIC_POLICIES_BY_OWNER,
  ];

  static final $core.List<AccountLimitType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static AccountLimitType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AccountLimitType._(super.value, super.name);
}

class ChangeAction extends $pb.ProtobufEnum {
  static const ChangeAction CHANGE_ACTION_CREATE =
      ChangeAction._(0, _omitEnumNames ? '' : 'CHANGE_ACTION_CREATE');
  static const ChangeAction CHANGE_ACTION_DELETE =
      ChangeAction._(1, _omitEnumNames ? '' : 'CHANGE_ACTION_DELETE');
  static const ChangeAction CHANGE_ACTION_UPSERT =
      ChangeAction._(2, _omitEnumNames ? '' : 'CHANGE_ACTION_UPSERT');

  static const $core.List<ChangeAction> values = <ChangeAction>[
    CHANGE_ACTION_CREATE,
    CHANGE_ACTION_DELETE,
    CHANGE_ACTION_UPSERT,
  ];

  static final $core.List<ChangeAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ChangeAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChangeAction._(super.value, super.name);
}

class ChangeStatus extends $pb.ProtobufEnum {
  static const ChangeStatus CHANGE_STATUS_PENDING =
      ChangeStatus._(0, _omitEnumNames ? '' : 'CHANGE_STATUS_PENDING');
  static const ChangeStatus CHANGE_STATUS_INSYNC =
      ChangeStatus._(1, _omitEnumNames ? '' : 'CHANGE_STATUS_INSYNC');

  static const $core.List<ChangeStatus> values = <ChangeStatus>[
    CHANGE_STATUS_PENDING,
    CHANGE_STATUS_INSYNC,
  ];

  static final $core.List<ChangeStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ChangeStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChangeStatus._(super.value, super.name);
}

class CidrCollectionChangeAction extends $pb.ProtobufEnum {
  static const CidrCollectionChangeAction
      CIDR_COLLECTION_CHANGE_ACTION_DELETE_IF_EXISTS =
      CidrCollectionChangeAction._(
          0,
          _omitEnumNames
              ? ''
              : 'CIDR_COLLECTION_CHANGE_ACTION_DELETE_IF_EXISTS');
  static const CidrCollectionChangeAction CIDR_COLLECTION_CHANGE_ACTION_PUT =
      CidrCollectionChangeAction._(
          1, _omitEnumNames ? '' : 'CIDR_COLLECTION_CHANGE_ACTION_PUT');

  static const $core.List<CidrCollectionChangeAction> values =
      <CidrCollectionChangeAction>[
    CIDR_COLLECTION_CHANGE_ACTION_DELETE_IF_EXISTS,
    CIDR_COLLECTION_CHANGE_ACTION_PUT,
  ];

  static final $core.List<CidrCollectionChangeAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CidrCollectionChangeAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CidrCollectionChangeAction._(super.value, super.name);
}

class CloudWatchRegion extends $pb.ProtobufEnum {
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTHEAST_3 =
      CloudWatchRegion._(
          0, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTHEAST_3');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_ISOF_SOUTH_1 =
      CloudWatchRegion._(
          1, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_ISOF_SOUTH_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_ISO_EAST_1 =
      CloudWatchRegion._(
          2, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_ISO_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_WEST_3 =
      CloudWatchRegion._(
          3, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_WEST_3');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_ISO_WEST_1 =
      CloudWatchRegion._(
          4, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_ISO_WEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_GOV_EAST_1 =
      CloudWatchRegion._(
          5, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_GOV_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_NORTHEAST_1 =
      CloudWatchRegion._(
          6, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_NORTHEAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_IL_CENTRAL_1 =
      CloudWatchRegion._(
          7, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_IL_CENTRAL_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_NORTH_1 =
      CloudWatchRegion._(
          8, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_NORTH_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_CN_NORTH_1 =
      CloudWatchRegion._(
          9, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_CN_NORTH_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_MX_CENTRAL_1 =
      CloudWatchRegion._(
          10, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_MX_CENTRAL_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTHEAST_5 =
      CloudWatchRegion._(
          11, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTHEAST_5');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EUSC_DE_EAST_1 =
      CloudWatchRegion._(
          12, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EUSC_DE_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_ISOF_EAST_1 =
      CloudWatchRegion._(
          13, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_ISOF_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_NORTHEAST_2 =
      CloudWatchRegion._(
          14, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_NORTHEAST_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_WEST_1 =
      CloudWatchRegion._(
          15, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_WEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTH_1 =
      CloudWatchRegion._(
          16, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTH_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_EAST_1 =
      CloudWatchRegion._(
          17, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_ISOB_WEST_1 =
      CloudWatchRegion._(
          18, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_ISOB_WEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_CA_WEST_1 =
      CloudWatchRegion._(
          19, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_CA_WEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTHEAST_2 =
      CloudWatchRegion._(
          20, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTHEAST_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_SOUTH_1 =
      CloudWatchRegion._(
          21, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_SOUTH_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_EAST_2 =
      CloudWatchRegion._(
          22, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_EAST_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_WEST_2 =
      CloudWatchRegion._(
          23, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_WEST_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_WEST_2 =
      CloudWatchRegion._(
          24, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_WEST_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_ME_CENTRAL_1 =
      CloudWatchRegion._(
          25, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_ME_CENTRAL_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTHEAST_7 =
      CloudWatchRegion._(
          26, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTHEAST_7');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_ISOB_EAST_1 =
      CloudWatchRegion._(
          27, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_ISOB_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_CN_NORTHWEST_1 =
      CloudWatchRegion._(
          28, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_CN_NORTHWEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_SOUTH_2 =
      CloudWatchRegion._(
          29, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_SOUTH_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_EAST_1 =
      CloudWatchRegion._(
          30, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_CENTRAL_1 =
      CloudWatchRegion._(
          31, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_CENTRAL_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTHEAST_4 =
      CloudWatchRegion._(
          32, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTHEAST_4');
  static const CloudWatchRegion CLOUD_WATCH_REGION_CA_CENTRAL_1 =
      CloudWatchRegion._(
          33, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_CA_CENTRAL_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_SA_EAST_1 =
      CloudWatchRegion._(
          34, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_SA_EAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTH_2 =
      CloudWatchRegion._(
          35, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTH_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_GOV_WEST_1 =
      CloudWatchRegion._(
          36, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_GOV_WEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_CENTRAL_2 =
      CloudWatchRegion._(
          37, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_CENTRAL_2');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTHEAST_1 =
      CloudWatchRegion._(
          38, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTHEAST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AF_SOUTH_1 =
      CloudWatchRegion._(
          39, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AF_SOUTH_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_ME_SOUTH_1 =
      CloudWatchRegion._(
          40, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_ME_SOUTH_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_WEST_1 =
      CloudWatchRegion._(
          41, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_WEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_EU_ISOE_WEST_1 =
      CloudWatchRegion._(
          42, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_EU_ISOE_WEST_1');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_SOUTHEAST_6 =
      CloudWatchRegion._(
          43, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_SOUTHEAST_6');
  static const CloudWatchRegion CLOUD_WATCH_REGION_AP_NORTHEAST_3 =
      CloudWatchRegion._(
          44, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_AP_NORTHEAST_3');
  static const CloudWatchRegion CLOUD_WATCH_REGION_US_EAST_2 =
      CloudWatchRegion._(
          45, _omitEnumNames ? '' : 'CLOUD_WATCH_REGION_US_EAST_2');

  static const $core.List<CloudWatchRegion> values = <CloudWatchRegion>[
    CLOUD_WATCH_REGION_AP_SOUTHEAST_3,
    CLOUD_WATCH_REGION_US_ISOF_SOUTH_1,
    CLOUD_WATCH_REGION_US_ISO_EAST_1,
    CLOUD_WATCH_REGION_EU_WEST_3,
    CLOUD_WATCH_REGION_US_ISO_WEST_1,
    CLOUD_WATCH_REGION_US_GOV_EAST_1,
    CLOUD_WATCH_REGION_AP_NORTHEAST_1,
    CLOUD_WATCH_REGION_IL_CENTRAL_1,
    CLOUD_WATCH_REGION_EU_NORTH_1,
    CLOUD_WATCH_REGION_CN_NORTH_1,
    CLOUD_WATCH_REGION_MX_CENTRAL_1,
    CLOUD_WATCH_REGION_AP_SOUTHEAST_5,
    CLOUD_WATCH_REGION_EUSC_DE_EAST_1,
    CLOUD_WATCH_REGION_US_ISOF_EAST_1,
    CLOUD_WATCH_REGION_AP_NORTHEAST_2,
    CLOUD_WATCH_REGION_US_WEST_1,
    CLOUD_WATCH_REGION_AP_SOUTH_1,
    CLOUD_WATCH_REGION_US_EAST_1,
    CLOUD_WATCH_REGION_US_ISOB_WEST_1,
    CLOUD_WATCH_REGION_CA_WEST_1,
    CLOUD_WATCH_REGION_AP_SOUTHEAST_2,
    CLOUD_WATCH_REGION_EU_SOUTH_1,
    CLOUD_WATCH_REGION_AP_EAST_2,
    CLOUD_WATCH_REGION_EU_WEST_2,
    CLOUD_WATCH_REGION_US_WEST_2,
    CLOUD_WATCH_REGION_ME_CENTRAL_1,
    CLOUD_WATCH_REGION_AP_SOUTHEAST_7,
    CLOUD_WATCH_REGION_US_ISOB_EAST_1,
    CLOUD_WATCH_REGION_CN_NORTHWEST_1,
    CLOUD_WATCH_REGION_EU_SOUTH_2,
    CLOUD_WATCH_REGION_AP_EAST_1,
    CLOUD_WATCH_REGION_EU_CENTRAL_1,
    CLOUD_WATCH_REGION_AP_SOUTHEAST_4,
    CLOUD_WATCH_REGION_CA_CENTRAL_1,
    CLOUD_WATCH_REGION_SA_EAST_1,
    CLOUD_WATCH_REGION_AP_SOUTH_2,
    CLOUD_WATCH_REGION_US_GOV_WEST_1,
    CLOUD_WATCH_REGION_EU_CENTRAL_2,
    CLOUD_WATCH_REGION_AP_SOUTHEAST_1,
    CLOUD_WATCH_REGION_AF_SOUTH_1,
    CLOUD_WATCH_REGION_ME_SOUTH_1,
    CLOUD_WATCH_REGION_EU_WEST_1,
    CLOUD_WATCH_REGION_EU_ISOE_WEST_1,
    CLOUD_WATCH_REGION_AP_SOUTHEAST_6,
    CLOUD_WATCH_REGION_AP_NORTHEAST_3,
    CLOUD_WATCH_REGION_US_EAST_2,
  ];

  static final $core.List<CloudWatchRegion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 45);
  static CloudWatchRegion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CloudWatchRegion._(super.value, super.name);
}

class ComparisonOperator extends $pb.ProtobufEnum {
  static const ComparisonOperator
      COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD = ComparisonOperator._(
          0,
          _omitEnumNames
              ? ''
              : 'COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD');
  static const ComparisonOperator COMPARISON_OPERATOR_GREATERTHANTHRESHOLD =
      ComparisonOperator._(
          1, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_GREATERTHANTHRESHOLD');
  static const ComparisonOperator COMPARISON_OPERATOR_LESSTHANTHRESHOLD =
      ComparisonOperator._(
          2, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_LESSTHANTHRESHOLD');
  static const ComparisonOperator
      COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD = ComparisonOperator._(
          3,
          _omitEnumNames
              ? ''
              : 'COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD');

  static const $core.List<ComparisonOperator> values = <ComparisonOperator>[
    COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD,
    COMPARISON_OPERATOR_GREATERTHANTHRESHOLD,
    COMPARISON_OPERATOR_LESSTHANTHRESHOLD,
    COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD,
  ];

  static final $core.List<ComparisonOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ComparisonOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComparisonOperator._(super.value, super.name);
}

class HealthCheckRegion extends $pb.ProtobufEnum {
  static const HealthCheckRegion HEALTH_CHECK_REGION_AP_NORTHEAST_1 =
      HealthCheckRegion._(
          0, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_AP_NORTHEAST_1');
  static const HealthCheckRegion HEALTH_CHECK_REGION_US_WEST_1 =
      HealthCheckRegion._(
          1, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_US_WEST_1');
  static const HealthCheckRegion HEALTH_CHECK_REGION_US_EAST_1 =
      HealthCheckRegion._(
          2, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_US_EAST_1');
  static const HealthCheckRegion HEALTH_CHECK_REGION_AP_SOUTHEAST_2 =
      HealthCheckRegion._(
          3, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_AP_SOUTHEAST_2');
  static const HealthCheckRegion HEALTH_CHECK_REGION_US_WEST_2 =
      HealthCheckRegion._(
          4, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_US_WEST_2');
  static const HealthCheckRegion HEALTH_CHECK_REGION_SA_EAST_1 =
      HealthCheckRegion._(
          5, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_SA_EAST_1');
  static const HealthCheckRegion HEALTH_CHECK_REGION_AP_SOUTHEAST_1 =
      HealthCheckRegion._(
          6, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_AP_SOUTHEAST_1');
  static const HealthCheckRegion HEALTH_CHECK_REGION_EU_WEST_1 =
      HealthCheckRegion._(
          7, _omitEnumNames ? '' : 'HEALTH_CHECK_REGION_EU_WEST_1');

  static const $core.List<HealthCheckRegion> values = <HealthCheckRegion>[
    HEALTH_CHECK_REGION_AP_NORTHEAST_1,
    HEALTH_CHECK_REGION_US_WEST_1,
    HEALTH_CHECK_REGION_US_EAST_1,
    HEALTH_CHECK_REGION_AP_SOUTHEAST_2,
    HEALTH_CHECK_REGION_US_WEST_2,
    HEALTH_CHECK_REGION_SA_EAST_1,
    HEALTH_CHECK_REGION_AP_SOUTHEAST_1,
    HEALTH_CHECK_REGION_EU_WEST_1,
  ];

  static final $core.List<HealthCheckRegion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static HealthCheckRegion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HealthCheckRegion._(super.value, super.name);
}

class HealthCheckType extends $pb.ProtobufEnum {
  static const HealthCheckType HEALTH_CHECK_TYPE_HTTP =
      HealthCheckType._(0, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_HTTP');
  static const HealthCheckType HEALTH_CHECK_TYPE_TCP =
      HealthCheckType._(1, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_TCP');
  static const HealthCheckType HEALTH_CHECK_TYPE_RECOVERY_CONTROL =
      HealthCheckType._(
          2, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_RECOVERY_CONTROL');
  static const HealthCheckType HEALTH_CHECK_TYPE_CALCULATED = HealthCheckType._(
      3, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_CALCULATED');
  static const HealthCheckType HEALTH_CHECK_TYPE_CLOUDWATCH_METRIC =
      HealthCheckType._(
          4, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_CLOUDWATCH_METRIC');
  static const HealthCheckType HEALTH_CHECK_TYPE_HTTPS_STR_MATCH =
      HealthCheckType._(
          5, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_HTTPS_STR_MATCH');
  static const HealthCheckType HEALTH_CHECK_TYPE_HTTP_STR_MATCH =
      HealthCheckType._(
          6, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_HTTP_STR_MATCH');
  static const HealthCheckType HEALTH_CHECK_TYPE_HTTPS =
      HealthCheckType._(7, _omitEnumNames ? '' : 'HEALTH_CHECK_TYPE_HTTPS');

  static const $core.List<HealthCheckType> values = <HealthCheckType>[
    HEALTH_CHECK_TYPE_HTTP,
    HEALTH_CHECK_TYPE_TCP,
    HEALTH_CHECK_TYPE_RECOVERY_CONTROL,
    HEALTH_CHECK_TYPE_CALCULATED,
    HEALTH_CHECK_TYPE_CLOUDWATCH_METRIC,
    HEALTH_CHECK_TYPE_HTTPS_STR_MATCH,
    HEALTH_CHECK_TYPE_HTTP_STR_MATCH,
    HEALTH_CHECK_TYPE_HTTPS,
  ];

  static final $core.List<HealthCheckType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static HealthCheckType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HealthCheckType._(super.value, super.name);
}

class HostedZoneLimitType extends $pb.ProtobufEnum {
  static const HostedZoneLimitType HOSTED_ZONE_LIMIT_TYPE_MAX_RRSETS_BY_ZONE =
      HostedZoneLimitType._(
          0, _omitEnumNames ? '' : 'HOSTED_ZONE_LIMIT_TYPE_MAX_RRSETS_BY_ZONE');
  static const HostedZoneLimitType
      HOSTED_ZONE_LIMIT_TYPE_MAX_VPCS_ASSOCIATED_BY_ZONE =
      HostedZoneLimitType._(
          1,
          _omitEnumNames
              ? ''
              : 'HOSTED_ZONE_LIMIT_TYPE_MAX_VPCS_ASSOCIATED_BY_ZONE');

  static const $core.List<HostedZoneLimitType> values = <HostedZoneLimitType>[
    HOSTED_ZONE_LIMIT_TYPE_MAX_RRSETS_BY_ZONE,
    HOSTED_ZONE_LIMIT_TYPE_MAX_VPCS_ASSOCIATED_BY_ZONE,
  ];

  static final $core.List<HostedZoneLimitType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static HostedZoneLimitType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HostedZoneLimitType._(super.value, super.name);
}

class HostedZoneType extends $pb.ProtobufEnum {
  static const HostedZoneType HOSTED_ZONE_TYPE_PRIVATE_HOSTED_ZONE =
      HostedZoneType._(
          0, _omitEnumNames ? '' : 'HOSTED_ZONE_TYPE_PRIVATE_HOSTED_ZONE');

  static const $core.List<HostedZoneType> values = <HostedZoneType>[
    HOSTED_ZONE_TYPE_PRIVATE_HOSTED_ZONE,
  ];

  static final $core.List<HostedZoneType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static HostedZoneType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HostedZoneType._(super.value, super.name);
}

class InsufficientDataHealthStatus extends $pb.ProtobufEnum {
  static const InsufficientDataHealthStatus
      INSUFFICIENT_DATA_HEALTH_STATUS_HEALTHY = InsufficientDataHealthStatus._(
          0, _omitEnumNames ? '' : 'INSUFFICIENT_DATA_HEALTH_STATUS_HEALTHY');
  static const InsufficientDataHealthStatus
      INSUFFICIENT_DATA_HEALTH_STATUS_LASTKNOWNSTATUS =
      InsufficientDataHealthStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'INSUFFICIENT_DATA_HEALTH_STATUS_LASTKNOWNSTATUS');
  static const InsufficientDataHealthStatus
      INSUFFICIENT_DATA_HEALTH_STATUS_UNHEALTHY =
      InsufficientDataHealthStatus._(
          2, _omitEnumNames ? '' : 'INSUFFICIENT_DATA_HEALTH_STATUS_UNHEALTHY');

  static const $core.List<InsufficientDataHealthStatus> values =
      <InsufficientDataHealthStatus>[
    INSUFFICIENT_DATA_HEALTH_STATUS_HEALTHY,
    INSUFFICIENT_DATA_HEALTH_STATUS_LASTKNOWNSTATUS,
    INSUFFICIENT_DATA_HEALTH_STATUS_UNHEALTHY,
  ];

  static final $core.List<InsufficientDataHealthStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static InsufficientDataHealthStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InsufficientDataHealthStatus._(super.value, super.name);
}

class RRType extends $pb.ProtobufEnum {
  static const RRType R_R_TYPE_PTR =
      RRType._(0, _omitEnumNames ? '' : 'R_R_TYPE_PTR');
  static const RRType R_R_TYPE_CAA =
      RRType._(1, _omitEnumNames ? '' : 'R_R_TYPE_CAA');
  static const RRType R_R_TYPE_SOA =
      RRType._(2, _omitEnumNames ? '' : 'R_R_TYPE_SOA');
  static const RRType R_R_TYPE_SRV =
      RRType._(3, _omitEnumNames ? '' : 'R_R_TYPE_SRV');
  static const RRType R_R_TYPE_CNAME =
      RRType._(4, _omitEnumNames ? '' : 'R_R_TYPE_CNAME');
  static const RRType R_R_TYPE_TLSA =
      RRType._(5, _omitEnumNames ? '' : 'R_R_TYPE_TLSA');
  static const RRType R_R_TYPE_A =
      RRType._(6, _omitEnumNames ? '' : 'R_R_TYPE_A');
  static const RRType R_R_TYPE_MX =
      RRType._(7, _omitEnumNames ? '' : 'R_R_TYPE_MX');
  static const RRType R_R_TYPE_SSHFP =
      RRType._(8, _omitEnumNames ? '' : 'R_R_TYPE_SSHFP');
  static const RRType R_R_TYPE_DS =
      RRType._(9, _omitEnumNames ? '' : 'R_R_TYPE_DS');
  static const RRType R_R_TYPE_AAAA =
      RRType._(10, _omitEnumNames ? '' : 'R_R_TYPE_AAAA');
  static const RRType R_R_TYPE_NAPTR =
      RRType._(11, _omitEnumNames ? '' : 'R_R_TYPE_NAPTR');
  static const RRType R_R_TYPE_NS =
      RRType._(12, _omitEnumNames ? '' : 'R_R_TYPE_NS');
  static const RRType R_R_TYPE_HTTPS =
      RRType._(13, _omitEnumNames ? '' : 'R_R_TYPE_HTTPS');
  static const RRType R_R_TYPE_SPF =
      RRType._(14, _omitEnumNames ? '' : 'R_R_TYPE_SPF');
  static const RRType R_R_TYPE_SVCB =
      RRType._(15, _omitEnumNames ? '' : 'R_R_TYPE_SVCB');
  static const RRType R_R_TYPE_TXT =
      RRType._(16, _omitEnumNames ? '' : 'R_R_TYPE_TXT');

  static const $core.List<RRType> values = <RRType>[
    R_R_TYPE_PTR,
    R_R_TYPE_CAA,
    R_R_TYPE_SOA,
    R_R_TYPE_SRV,
    R_R_TYPE_CNAME,
    R_R_TYPE_TLSA,
    R_R_TYPE_A,
    R_R_TYPE_MX,
    R_R_TYPE_SSHFP,
    R_R_TYPE_DS,
    R_R_TYPE_AAAA,
    R_R_TYPE_NAPTR,
    R_R_TYPE_NS,
    R_R_TYPE_HTTPS,
    R_R_TYPE_SPF,
    R_R_TYPE_SVCB,
    R_R_TYPE_TXT,
  ];

  static final $core.List<RRType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 16);
  static RRType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RRType._(super.value, super.name);
}

class ResettableElementName extends $pb.ProtobufEnum {
  static const ResettableElementName RESETTABLE_ELEMENT_NAME_REGIONS =
      ResettableElementName._(
          0, _omitEnumNames ? '' : 'RESETTABLE_ELEMENT_NAME_REGIONS');
  static const ResettableElementName RESETTABLE_ELEMENT_NAME_CHILDHEALTHCHECKS =
      ResettableElementName._(
          1, _omitEnumNames ? '' : 'RESETTABLE_ELEMENT_NAME_CHILDHEALTHCHECKS');
  static const ResettableElementName
      RESETTABLE_ELEMENT_NAME_FULLYQUALIFIEDDOMAINNAME =
      ResettableElementName._(
          2,
          _omitEnumNames
              ? ''
              : 'RESETTABLE_ELEMENT_NAME_FULLYQUALIFIEDDOMAINNAME');
  static const ResettableElementName RESETTABLE_ELEMENT_NAME_RESOURCEPATH =
      ResettableElementName._(
          3, _omitEnumNames ? '' : 'RESETTABLE_ELEMENT_NAME_RESOURCEPATH');

  static const $core.List<ResettableElementName> values =
      <ResettableElementName>[
    RESETTABLE_ELEMENT_NAME_REGIONS,
    RESETTABLE_ELEMENT_NAME_CHILDHEALTHCHECKS,
    RESETTABLE_ELEMENT_NAME_FULLYQUALIFIEDDOMAINNAME,
    RESETTABLE_ELEMENT_NAME_RESOURCEPATH,
  ];

  static final $core.List<ResettableElementName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ResettableElementName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResettableElementName._(super.value, super.name);
}

class ResourceRecordSetFailover extends $pb.ProtobufEnum {
  static const ResourceRecordSetFailover
      RESOURCE_RECORD_SET_FAILOVER_SECONDARY = ResourceRecordSetFailover._(
          0, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_FAILOVER_SECONDARY');
  static const ResourceRecordSetFailover RESOURCE_RECORD_SET_FAILOVER_PRIMARY =
      ResourceRecordSetFailover._(
          1, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_FAILOVER_PRIMARY');

  static const $core.List<ResourceRecordSetFailover> values =
      <ResourceRecordSetFailover>[
    RESOURCE_RECORD_SET_FAILOVER_SECONDARY,
    RESOURCE_RECORD_SET_FAILOVER_PRIMARY,
  ];

  static final $core.List<ResourceRecordSetFailover?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ResourceRecordSetFailover? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceRecordSetFailover._(super.value, super.name);
}

class ResourceRecordSetRegion extends $pb.ProtobufEnum {
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_3 = ResourceRecordSetRegion._(
          0, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_3');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_WEST_3 =
      ResourceRecordSetRegion._(
          1, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_WEST_3');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_US_GOV_EAST_1 = ResourceRecordSetRegion._(
          2, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_US_GOV_EAST_1');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_1 = ResourceRecordSetRegion._(
          3, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_IL_CENTRAL_1 =
      ResourceRecordSetRegion._(
          4, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_IL_CENTRAL_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_NORTH_1 =
      ResourceRecordSetRegion._(
          5, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_NORTH_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_CN_NORTH_1 =
      ResourceRecordSetRegion._(
          6, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_CN_NORTH_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_MX_CENTRAL_1 =
      ResourceRecordSetRegion._(
          7, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_MX_CENTRAL_1');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_5 = ResourceRecordSetRegion._(
          8, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_5');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_EUSC_DE_EAST_1 = ResourceRecordSetRegion._(
          9, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EUSC_DE_EAST_1');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_2 = ResourceRecordSetRegion._(10,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_2');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_US_WEST_1 =
      ResourceRecordSetRegion._(
          11, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_US_WEST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_AP_SOUTH_1 =
      ResourceRecordSetRegion._(
          12, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTH_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_US_EAST_1 =
      ResourceRecordSetRegion._(
          13, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_US_EAST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_CA_WEST_1 =
      ResourceRecordSetRegion._(
          14, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_CA_WEST_1');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_2 = ResourceRecordSetRegion._(15,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_2');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_SOUTH_1 =
      ResourceRecordSetRegion._(
          16, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_SOUTH_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_AP_EAST_2 =
      ResourceRecordSetRegion._(
          17, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_EAST_2');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_WEST_2 =
      ResourceRecordSetRegion._(
          18, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_WEST_2');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_US_WEST_2 =
      ResourceRecordSetRegion._(
          19, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_US_WEST_2');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_ME_CENTRAL_1 =
      ResourceRecordSetRegion._(
          20, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_ME_CENTRAL_1');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_7 = ResourceRecordSetRegion._(21,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_7');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_CN_NORTHWEST_1 = ResourceRecordSetRegion._(22,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_CN_NORTHWEST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_SOUTH_2 =
      ResourceRecordSetRegion._(
          23, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_SOUTH_2');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_AP_EAST_1 =
      ResourceRecordSetRegion._(
          24, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_EAST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_CENTRAL_1 =
      ResourceRecordSetRegion._(
          25, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_CENTRAL_1');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_4 = ResourceRecordSetRegion._(26,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_4');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_CA_CENTRAL_1 =
      ResourceRecordSetRegion._(
          27, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_CA_CENTRAL_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_SA_EAST_1 =
      ResourceRecordSetRegion._(
          28, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_SA_EAST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_AP_SOUTH_2 =
      ResourceRecordSetRegion._(
          29, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTH_2');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_US_GOV_WEST_1 = ResourceRecordSetRegion._(
          30, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_US_GOV_WEST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_CENTRAL_2 =
      ResourceRecordSetRegion._(
          31, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_CENTRAL_2');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_1 = ResourceRecordSetRegion._(32,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_AF_SOUTH_1 =
      ResourceRecordSetRegion._(
          33, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AF_SOUTH_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_ME_SOUTH_1 =
      ResourceRecordSetRegion._(
          34, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_ME_SOUTH_1');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_EU_WEST_1 =
      ResourceRecordSetRegion._(
          35, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_EU_WEST_1');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_6 = ResourceRecordSetRegion._(36,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_6');
  static const ResourceRecordSetRegion
      RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_3 = ResourceRecordSetRegion._(37,
          _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_3');
  static const ResourceRecordSetRegion RESOURCE_RECORD_SET_REGION_US_EAST_2 =
      ResourceRecordSetRegion._(
          38, _omitEnumNames ? '' : 'RESOURCE_RECORD_SET_REGION_US_EAST_2');

  static const $core.List<ResourceRecordSetRegion> values =
      <ResourceRecordSetRegion>[
    RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_3,
    RESOURCE_RECORD_SET_REGION_EU_WEST_3,
    RESOURCE_RECORD_SET_REGION_US_GOV_EAST_1,
    RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_1,
    RESOURCE_RECORD_SET_REGION_IL_CENTRAL_1,
    RESOURCE_RECORD_SET_REGION_EU_NORTH_1,
    RESOURCE_RECORD_SET_REGION_CN_NORTH_1,
    RESOURCE_RECORD_SET_REGION_MX_CENTRAL_1,
    RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_5,
    RESOURCE_RECORD_SET_REGION_EUSC_DE_EAST_1,
    RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_2,
    RESOURCE_RECORD_SET_REGION_US_WEST_1,
    RESOURCE_RECORD_SET_REGION_AP_SOUTH_1,
    RESOURCE_RECORD_SET_REGION_US_EAST_1,
    RESOURCE_RECORD_SET_REGION_CA_WEST_1,
    RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_2,
    RESOURCE_RECORD_SET_REGION_EU_SOUTH_1,
    RESOURCE_RECORD_SET_REGION_AP_EAST_2,
    RESOURCE_RECORD_SET_REGION_EU_WEST_2,
    RESOURCE_RECORD_SET_REGION_US_WEST_2,
    RESOURCE_RECORD_SET_REGION_ME_CENTRAL_1,
    RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_7,
    RESOURCE_RECORD_SET_REGION_CN_NORTHWEST_1,
    RESOURCE_RECORD_SET_REGION_EU_SOUTH_2,
    RESOURCE_RECORD_SET_REGION_AP_EAST_1,
    RESOURCE_RECORD_SET_REGION_EU_CENTRAL_1,
    RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_4,
    RESOURCE_RECORD_SET_REGION_CA_CENTRAL_1,
    RESOURCE_RECORD_SET_REGION_SA_EAST_1,
    RESOURCE_RECORD_SET_REGION_AP_SOUTH_2,
    RESOURCE_RECORD_SET_REGION_US_GOV_WEST_1,
    RESOURCE_RECORD_SET_REGION_EU_CENTRAL_2,
    RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_1,
    RESOURCE_RECORD_SET_REGION_AF_SOUTH_1,
    RESOURCE_RECORD_SET_REGION_ME_SOUTH_1,
    RESOURCE_RECORD_SET_REGION_EU_WEST_1,
    RESOURCE_RECORD_SET_REGION_AP_SOUTHEAST_6,
    RESOURCE_RECORD_SET_REGION_AP_NORTHEAST_3,
    RESOURCE_RECORD_SET_REGION_US_EAST_2,
  ];

  static final $core.List<ResourceRecordSetRegion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 38);
  static ResourceRecordSetRegion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceRecordSetRegion._(super.value, super.name);
}

class ReusableDelegationSetLimitType extends $pb.ProtobufEnum {
  static const ReusableDelegationSetLimitType
      REUSABLE_DELEGATION_SET_LIMIT_TYPE_MAX_ZONES_BY_REUSABLE_DELEGATION_SET =
      ReusableDelegationSetLimitType._(
          0,
          _omitEnumNames
              ? ''
              : 'REUSABLE_DELEGATION_SET_LIMIT_TYPE_MAX_ZONES_BY_REUSABLE_DELEGATION_SET');

  static const $core.List<ReusableDelegationSetLimitType> values =
      <ReusableDelegationSetLimitType>[
    REUSABLE_DELEGATION_SET_LIMIT_TYPE_MAX_ZONES_BY_REUSABLE_DELEGATION_SET,
  ];

  static final $core.List<ReusableDelegationSetLimitType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ReusableDelegationSetLimitType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReusableDelegationSetLimitType._(super.value, super.name);
}

class Statistic extends $pb.ProtobufEnum {
  static const Statistic STATISTIC_SUM =
      Statistic._(0, _omitEnumNames ? '' : 'STATISTIC_SUM');
  static const Statistic STATISTIC_SAMPLECOUNT =
      Statistic._(1, _omitEnumNames ? '' : 'STATISTIC_SAMPLECOUNT');
  static const Statistic STATISTIC_AVERAGE =
      Statistic._(2, _omitEnumNames ? '' : 'STATISTIC_AVERAGE');
  static const Statistic STATISTIC_MAXIMUM =
      Statistic._(3, _omitEnumNames ? '' : 'STATISTIC_MAXIMUM');
  static const Statistic STATISTIC_MINIMUM =
      Statistic._(4, _omitEnumNames ? '' : 'STATISTIC_MINIMUM');

  static const $core.List<Statistic> values = <Statistic>[
    STATISTIC_SUM,
    STATISTIC_SAMPLECOUNT,
    STATISTIC_AVERAGE,
    STATISTIC_MAXIMUM,
    STATISTIC_MINIMUM,
  ];

  static final $core.List<Statistic?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static Statistic? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Statistic._(super.value, super.name);
}

class TagResourceType extends $pb.ProtobufEnum {
  static const TagResourceType TAG_RESOURCE_TYPE_HEALTHCHECK =
      TagResourceType._(
          0, _omitEnumNames ? '' : 'TAG_RESOURCE_TYPE_HEALTHCHECK');
  static const TagResourceType TAG_RESOURCE_TYPE_HOSTEDZONE = TagResourceType._(
      1, _omitEnumNames ? '' : 'TAG_RESOURCE_TYPE_HOSTEDZONE');

  static const $core.List<TagResourceType> values = <TagResourceType>[
    TAG_RESOURCE_TYPE_HEALTHCHECK,
    TAG_RESOURCE_TYPE_HOSTEDZONE,
  ];

  static final $core.List<TagResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static TagResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TagResourceType._(super.value, super.name);
}

class VPCRegion extends $pb.ProtobufEnum {
  static const VPCRegion V_P_C_REGION_AP_SOUTHEAST_3 =
      VPCRegion._(0, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTHEAST_3');
  static const VPCRegion V_P_C_REGION_US_ISOF_SOUTH_1 =
      VPCRegion._(1, _omitEnumNames ? '' : 'V_P_C_REGION_US_ISOF_SOUTH_1');
  static const VPCRegion V_P_C_REGION_US_ISO_EAST_1 =
      VPCRegion._(2, _omitEnumNames ? '' : 'V_P_C_REGION_US_ISO_EAST_1');
  static const VPCRegion V_P_C_REGION_EU_WEST_3 =
      VPCRegion._(3, _omitEnumNames ? '' : 'V_P_C_REGION_EU_WEST_3');
  static const VPCRegion V_P_C_REGION_US_ISO_WEST_1 =
      VPCRegion._(4, _omitEnumNames ? '' : 'V_P_C_REGION_US_ISO_WEST_1');
  static const VPCRegion V_P_C_REGION_US_GOV_EAST_1 =
      VPCRegion._(5, _omitEnumNames ? '' : 'V_P_C_REGION_US_GOV_EAST_1');
  static const VPCRegion V_P_C_REGION_AP_NORTHEAST_1 =
      VPCRegion._(6, _omitEnumNames ? '' : 'V_P_C_REGION_AP_NORTHEAST_1');
  static const VPCRegion V_P_C_REGION_IL_CENTRAL_1 =
      VPCRegion._(7, _omitEnumNames ? '' : 'V_P_C_REGION_IL_CENTRAL_1');
  static const VPCRegion V_P_C_REGION_EU_NORTH_1 =
      VPCRegion._(8, _omitEnumNames ? '' : 'V_P_C_REGION_EU_NORTH_1');
  static const VPCRegion V_P_C_REGION_CN_NORTH_1 =
      VPCRegion._(9, _omitEnumNames ? '' : 'V_P_C_REGION_CN_NORTH_1');
  static const VPCRegion V_P_C_REGION_MX_CENTRAL_1 =
      VPCRegion._(10, _omitEnumNames ? '' : 'V_P_C_REGION_MX_CENTRAL_1');
  static const VPCRegion V_P_C_REGION_AP_SOUTHEAST_5 =
      VPCRegion._(11, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTHEAST_5');
  static const VPCRegion V_P_C_REGION_EUSC_DE_EAST_1 =
      VPCRegion._(12, _omitEnumNames ? '' : 'V_P_C_REGION_EUSC_DE_EAST_1');
  static const VPCRegion V_P_C_REGION_US_ISOF_EAST_1 =
      VPCRegion._(13, _omitEnumNames ? '' : 'V_P_C_REGION_US_ISOF_EAST_1');
  static const VPCRegion V_P_C_REGION_AP_NORTHEAST_2 =
      VPCRegion._(14, _omitEnumNames ? '' : 'V_P_C_REGION_AP_NORTHEAST_2');
  static const VPCRegion V_P_C_REGION_US_WEST_1 =
      VPCRegion._(15, _omitEnumNames ? '' : 'V_P_C_REGION_US_WEST_1');
  static const VPCRegion V_P_C_REGION_AP_SOUTH_1 =
      VPCRegion._(16, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTH_1');
  static const VPCRegion V_P_C_REGION_US_EAST_1 =
      VPCRegion._(17, _omitEnumNames ? '' : 'V_P_C_REGION_US_EAST_1');
  static const VPCRegion V_P_C_REGION_US_ISOB_WEST_1 =
      VPCRegion._(18, _omitEnumNames ? '' : 'V_P_C_REGION_US_ISOB_WEST_1');
  static const VPCRegion V_P_C_REGION_CA_WEST_1 =
      VPCRegion._(19, _omitEnumNames ? '' : 'V_P_C_REGION_CA_WEST_1');
  static const VPCRegion V_P_C_REGION_AP_SOUTHEAST_2 =
      VPCRegion._(20, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTHEAST_2');
  static const VPCRegion V_P_C_REGION_EU_SOUTH_1 =
      VPCRegion._(21, _omitEnumNames ? '' : 'V_P_C_REGION_EU_SOUTH_1');
  static const VPCRegion V_P_C_REGION_AP_EAST_2 =
      VPCRegion._(22, _omitEnumNames ? '' : 'V_P_C_REGION_AP_EAST_2');
  static const VPCRegion V_P_C_REGION_EU_WEST_2 =
      VPCRegion._(23, _omitEnumNames ? '' : 'V_P_C_REGION_EU_WEST_2');
  static const VPCRegion V_P_C_REGION_US_WEST_2 =
      VPCRegion._(24, _omitEnumNames ? '' : 'V_P_C_REGION_US_WEST_2');
  static const VPCRegion V_P_C_REGION_ME_CENTRAL_1 =
      VPCRegion._(25, _omitEnumNames ? '' : 'V_P_C_REGION_ME_CENTRAL_1');
  static const VPCRegion V_P_C_REGION_AP_SOUTHEAST_7 =
      VPCRegion._(26, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTHEAST_7');
  static const VPCRegion V_P_C_REGION_US_ISOB_EAST_1 =
      VPCRegion._(27, _omitEnumNames ? '' : 'V_P_C_REGION_US_ISOB_EAST_1');
  static const VPCRegion V_P_C_REGION_CN_NORTHWEST_1 =
      VPCRegion._(28, _omitEnumNames ? '' : 'V_P_C_REGION_CN_NORTHWEST_1');
  static const VPCRegion V_P_C_REGION_EU_SOUTH_2 =
      VPCRegion._(29, _omitEnumNames ? '' : 'V_P_C_REGION_EU_SOUTH_2');
  static const VPCRegion V_P_C_REGION_AP_EAST_1 =
      VPCRegion._(30, _omitEnumNames ? '' : 'V_P_C_REGION_AP_EAST_1');
  static const VPCRegion V_P_C_REGION_EU_CENTRAL_1 =
      VPCRegion._(31, _omitEnumNames ? '' : 'V_P_C_REGION_EU_CENTRAL_1');
  static const VPCRegion V_P_C_REGION_AP_SOUTHEAST_4 =
      VPCRegion._(32, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTHEAST_4');
  static const VPCRegion V_P_C_REGION_CA_CENTRAL_1 =
      VPCRegion._(33, _omitEnumNames ? '' : 'V_P_C_REGION_CA_CENTRAL_1');
  static const VPCRegion V_P_C_REGION_SA_EAST_1 =
      VPCRegion._(34, _omitEnumNames ? '' : 'V_P_C_REGION_SA_EAST_1');
  static const VPCRegion V_P_C_REGION_AP_SOUTH_2 =
      VPCRegion._(35, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTH_2');
  static const VPCRegion V_P_C_REGION_US_GOV_WEST_1 =
      VPCRegion._(36, _omitEnumNames ? '' : 'V_P_C_REGION_US_GOV_WEST_1');
  static const VPCRegion V_P_C_REGION_EU_CENTRAL_2 =
      VPCRegion._(37, _omitEnumNames ? '' : 'V_P_C_REGION_EU_CENTRAL_2');
  static const VPCRegion V_P_C_REGION_AP_SOUTHEAST_1 =
      VPCRegion._(38, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTHEAST_1');
  static const VPCRegion V_P_C_REGION_AF_SOUTH_1 =
      VPCRegion._(39, _omitEnumNames ? '' : 'V_P_C_REGION_AF_SOUTH_1');
  static const VPCRegion V_P_C_REGION_ME_SOUTH_1 =
      VPCRegion._(40, _omitEnumNames ? '' : 'V_P_C_REGION_ME_SOUTH_1');
  static const VPCRegion V_P_C_REGION_EU_WEST_1 =
      VPCRegion._(41, _omitEnumNames ? '' : 'V_P_C_REGION_EU_WEST_1');
  static const VPCRegion V_P_C_REGION_EU_ISOE_WEST_1 =
      VPCRegion._(42, _omitEnumNames ? '' : 'V_P_C_REGION_EU_ISOE_WEST_1');
  static const VPCRegion V_P_C_REGION_AP_SOUTHEAST_6 =
      VPCRegion._(43, _omitEnumNames ? '' : 'V_P_C_REGION_AP_SOUTHEAST_6');
  static const VPCRegion V_P_C_REGION_AP_NORTHEAST_3 =
      VPCRegion._(44, _omitEnumNames ? '' : 'V_P_C_REGION_AP_NORTHEAST_3');
  static const VPCRegion V_P_C_REGION_US_EAST_2 =
      VPCRegion._(45, _omitEnumNames ? '' : 'V_P_C_REGION_US_EAST_2');

  static const $core.List<VPCRegion> values = <VPCRegion>[
    V_P_C_REGION_AP_SOUTHEAST_3,
    V_P_C_REGION_US_ISOF_SOUTH_1,
    V_P_C_REGION_US_ISO_EAST_1,
    V_P_C_REGION_EU_WEST_3,
    V_P_C_REGION_US_ISO_WEST_1,
    V_P_C_REGION_US_GOV_EAST_1,
    V_P_C_REGION_AP_NORTHEAST_1,
    V_P_C_REGION_IL_CENTRAL_1,
    V_P_C_REGION_EU_NORTH_1,
    V_P_C_REGION_CN_NORTH_1,
    V_P_C_REGION_MX_CENTRAL_1,
    V_P_C_REGION_AP_SOUTHEAST_5,
    V_P_C_REGION_EUSC_DE_EAST_1,
    V_P_C_REGION_US_ISOF_EAST_1,
    V_P_C_REGION_AP_NORTHEAST_2,
    V_P_C_REGION_US_WEST_1,
    V_P_C_REGION_AP_SOUTH_1,
    V_P_C_REGION_US_EAST_1,
    V_P_C_REGION_US_ISOB_WEST_1,
    V_P_C_REGION_CA_WEST_1,
    V_P_C_REGION_AP_SOUTHEAST_2,
    V_P_C_REGION_EU_SOUTH_1,
    V_P_C_REGION_AP_EAST_2,
    V_P_C_REGION_EU_WEST_2,
    V_P_C_REGION_US_WEST_2,
    V_P_C_REGION_ME_CENTRAL_1,
    V_P_C_REGION_AP_SOUTHEAST_7,
    V_P_C_REGION_US_ISOB_EAST_1,
    V_P_C_REGION_CN_NORTHWEST_1,
    V_P_C_REGION_EU_SOUTH_2,
    V_P_C_REGION_AP_EAST_1,
    V_P_C_REGION_EU_CENTRAL_1,
    V_P_C_REGION_AP_SOUTHEAST_4,
    V_P_C_REGION_CA_CENTRAL_1,
    V_P_C_REGION_SA_EAST_1,
    V_P_C_REGION_AP_SOUTH_2,
    V_P_C_REGION_US_GOV_WEST_1,
    V_P_C_REGION_EU_CENTRAL_2,
    V_P_C_REGION_AP_SOUTHEAST_1,
    V_P_C_REGION_AF_SOUTH_1,
    V_P_C_REGION_ME_SOUTH_1,
    V_P_C_REGION_EU_WEST_1,
    V_P_C_REGION_EU_ISOE_WEST_1,
    V_P_C_REGION_AP_SOUTHEAST_6,
    V_P_C_REGION_AP_NORTHEAST_3,
    V_P_C_REGION_US_EAST_2,
  ];

  static final $core.List<VPCRegion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 45);
  static VPCRegion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VPCRegion._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
