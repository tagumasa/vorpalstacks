// This is a generated file - do not edit.
//
// Generated from dynamodb.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ApproximateCreationDateTimePrecision extends $pb.ProtobufEnum {
  static const ApproximateCreationDateTimePrecision
      APPROXIMATE_CREATION_DATE_TIME_PRECISION_MILLISECOND =
      ApproximateCreationDateTimePrecision._(
          0,
          _omitEnumNames
              ? ''
              : 'APPROXIMATE_CREATION_DATE_TIME_PRECISION_MILLISECOND');
  static const ApproximateCreationDateTimePrecision
      APPROXIMATE_CREATION_DATE_TIME_PRECISION_MICROSECOND =
      ApproximateCreationDateTimePrecision._(
          1,
          _omitEnumNames
              ? ''
              : 'APPROXIMATE_CREATION_DATE_TIME_PRECISION_MICROSECOND');

  static const $core.List<ApproximateCreationDateTimePrecision> values =
      <ApproximateCreationDateTimePrecision>[
    APPROXIMATE_CREATION_DATE_TIME_PRECISION_MILLISECOND,
    APPROXIMATE_CREATION_DATE_TIME_PRECISION_MICROSECOND,
  ];

  static final $core.List<ApproximateCreationDateTimePrecision?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ApproximateCreationDateTimePrecision? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ApproximateCreationDateTimePrecision._(super.value, super.name);
}

class AttributeAction extends $pb.ProtobufEnum {
  static const AttributeAction ATTRIBUTE_ACTION_ADD =
      AttributeAction._(0, _omitEnumNames ? '' : 'ATTRIBUTE_ACTION_ADD');
  static const AttributeAction ATTRIBUTE_ACTION_DELETE =
      AttributeAction._(1, _omitEnumNames ? '' : 'ATTRIBUTE_ACTION_DELETE');
  static const AttributeAction ATTRIBUTE_ACTION_PUT =
      AttributeAction._(2, _omitEnumNames ? '' : 'ATTRIBUTE_ACTION_PUT');

  static const $core.List<AttributeAction> values = <AttributeAction>[
    ATTRIBUTE_ACTION_ADD,
    ATTRIBUTE_ACTION_DELETE,
    ATTRIBUTE_ACTION_PUT,
  ];

  static final $core.List<AttributeAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AttributeAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AttributeAction._(super.value, super.name);
}

class BackupStatus extends $pb.ProtobufEnum {
  static const BackupStatus BACKUP_STATUS_AVAILABLE =
      BackupStatus._(0, _omitEnumNames ? '' : 'BACKUP_STATUS_AVAILABLE');
  static const BackupStatus BACKUP_STATUS_CREATING =
      BackupStatus._(1, _omitEnumNames ? '' : 'BACKUP_STATUS_CREATING');
  static const BackupStatus BACKUP_STATUS_DELETED =
      BackupStatus._(2, _omitEnumNames ? '' : 'BACKUP_STATUS_DELETED');

  static const $core.List<BackupStatus> values = <BackupStatus>[
    BACKUP_STATUS_AVAILABLE,
    BACKUP_STATUS_CREATING,
    BACKUP_STATUS_DELETED,
  ];

  static final $core.List<BackupStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static BackupStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BackupStatus._(super.value, super.name);
}

class BackupType extends $pb.ProtobufEnum {
  static const BackupType BACKUP_TYPE_SYSTEM =
      BackupType._(0, _omitEnumNames ? '' : 'BACKUP_TYPE_SYSTEM');
  static const BackupType BACKUP_TYPE_USER =
      BackupType._(1, _omitEnumNames ? '' : 'BACKUP_TYPE_USER');
  static const BackupType BACKUP_TYPE_AWS_BACKUP =
      BackupType._(2, _omitEnumNames ? '' : 'BACKUP_TYPE_AWS_BACKUP');

  static const $core.List<BackupType> values = <BackupType>[
    BACKUP_TYPE_SYSTEM,
    BACKUP_TYPE_USER,
    BACKUP_TYPE_AWS_BACKUP,
  ];

  static final $core.List<BackupType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static BackupType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BackupType._(super.value, super.name);
}

class BackupTypeFilter extends $pb.ProtobufEnum {
  static const BackupTypeFilter BACKUP_TYPE_FILTER_SYSTEM =
      BackupTypeFilter._(0, _omitEnumNames ? '' : 'BACKUP_TYPE_FILTER_SYSTEM');
  static const BackupTypeFilter BACKUP_TYPE_FILTER_USER =
      BackupTypeFilter._(1, _omitEnumNames ? '' : 'BACKUP_TYPE_FILTER_USER');
  static const BackupTypeFilter BACKUP_TYPE_FILTER_AWS_BACKUP =
      BackupTypeFilter._(
          2, _omitEnumNames ? '' : 'BACKUP_TYPE_FILTER_AWS_BACKUP');
  static const BackupTypeFilter BACKUP_TYPE_FILTER_ALL =
      BackupTypeFilter._(3, _omitEnumNames ? '' : 'BACKUP_TYPE_FILTER_ALL');

  static const $core.List<BackupTypeFilter> values = <BackupTypeFilter>[
    BACKUP_TYPE_FILTER_SYSTEM,
    BACKUP_TYPE_FILTER_USER,
    BACKUP_TYPE_FILTER_AWS_BACKUP,
    BACKUP_TYPE_FILTER_ALL,
  ];

  static final $core.List<BackupTypeFilter?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static BackupTypeFilter? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BackupTypeFilter._(super.value, super.name);
}

class BatchStatementErrorCodeEnum extends $pb.ProtobufEnum {
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_REQUESTLIMITEXCEEDED =
      BatchStatementErrorCodeEnum._(
          0,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_REQUESTLIMITEXCEEDED');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_ITEMCOLLECTIONSIZELIMITEXCEEDED =
      BatchStatementErrorCodeEnum._(
          1,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_ITEMCOLLECTIONSIZELIMITEXCEEDED');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_INTERNALSERVERERROR =
      BatchStatementErrorCodeEnum._(
          2,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_INTERNALSERVERERROR');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_RESOURCENOTFOUND =
      BatchStatementErrorCodeEnum._(
          3,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_RESOURCENOTFOUND');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_ACCESSDENIED =
      BatchStatementErrorCodeEnum._(4,
          _omitEnumNames ? '' : 'BATCH_STATEMENT_ERROR_CODE_ENUM_ACCESSDENIED');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_DUPLICATEITEM =
      BatchStatementErrorCodeEnum._(
          5,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_DUPLICATEITEM');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_VALIDATIONERROR =
      BatchStatementErrorCodeEnum._(
          6,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_VALIDATIONERROR');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_CONDITIONALCHECKFAILED =
      BatchStatementErrorCodeEnum._(
          7,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_CONDITIONALCHECKFAILED');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_THROTTLINGERROR =
      BatchStatementErrorCodeEnum._(
          8,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_THROTTLINGERROR');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_PROVISIONEDTHROUGHPUTEXCEEDED =
      BatchStatementErrorCodeEnum._(
          9,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_PROVISIONEDTHROUGHPUTEXCEEDED');
  static const BatchStatementErrorCodeEnum
      BATCH_STATEMENT_ERROR_CODE_ENUM_TRANSACTIONCONFLICT =
      BatchStatementErrorCodeEnum._(
          10,
          _omitEnumNames
              ? ''
              : 'BATCH_STATEMENT_ERROR_CODE_ENUM_TRANSACTIONCONFLICT');

  static const $core.List<BatchStatementErrorCodeEnum> values =
      <BatchStatementErrorCodeEnum>[
    BATCH_STATEMENT_ERROR_CODE_ENUM_REQUESTLIMITEXCEEDED,
    BATCH_STATEMENT_ERROR_CODE_ENUM_ITEMCOLLECTIONSIZELIMITEXCEEDED,
    BATCH_STATEMENT_ERROR_CODE_ENUM_INTERNALSERVERERROR,
    BATCH_STATEMENT_ERROR_CODE_ENUM_RESOURCENOTFOUND,
    BATCH_STATEMENT_ERROR_CODE_ENUM_ACCESSDENIED,
    BATCH_STATEMENT_ERROR_CODE_ENUM_DUPLICATEITEM,
    BATCH_STATEMENT_ERROR_CODE_ENUM_VALIDATIONERROR,
    BATCH_STATEMENT_ERROR_CODE_ENUM_CONDITIONALCHECKFAILED,
    BATCH_STATEMENT_ERROR_CODE_ENUM_THROTTLINGERROR,
    BATCH_STATEMENT_ERROR_CODE_ENUM_PROVISIONEDTHROUGHPUTEXCEEDED,
    BATCH_STATEMENT_ERROR_CODE_ENUM_TRANSACTIONCONFLICT,
  ];

  static final $core.List<BatchStatementErrorCodeEnum?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 10);
  static BatchStatementErrorCodeEnum? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BatchStatementErrorCodeEnum._(super.value, super.name);
}

class BillingMode extends $pb.ProtobufEnum {
  static const BillingMode BILLING_MODE_PAY_PER_REQUEST =
      BillingMode._(0, _omitEnumNames ? '' : 'BILLING_MODE_PAY_PER_REQUEST');
  static const BillingMode BILLING_MODE_PROVISIONED =
      BillingMode._(1, _omitEnumNames ? '' : 'BILLING_MODE_PROVISIONED');

  static const $core.List<BillingMode> values = <BillingMode>[
    BILLING_MODE_PAY_PER_REQUEST,
    BILLING_MODE_PROVISIONED,
  ];

  static final $core.List<BillingMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BillingMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BillingMode._(super.value, super.name);
}

class ComparisonOperator extends $pb.ProtobufEnum {
  static const ComparisonOperator COMPARISON_OPERATOR_GE =
      ComparisonOperator._(0, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_GE');
  static const ComparisonOperator COMPARISON_OPERATOR_NOT_NULL =
      ComparisonOperator._(
          1, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_NOT_NULL');
  static const ComparisonOperator COMPARISON_OPERATOR_EQ =
      ComparisonOperator._(2, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_EQ');
  static const ComparisonOperator COMPARISON_OPERATOR_IN =
      ComparisonOperator._(3, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_IN');
  static const ComparisonOperator COMPARISON_OPERATOR_LT =
      ComparisonOperator._(4, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_LT');
  static const ComparisonOperator COMPARISON_OPERATOR_CONTAINS =
      ComparisonOperator._(
          5, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_CONTAINS');
  static const ComparisonOperator COMPARISON_OPERATOR_GT =
      ComparisonOperator._(6, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_GT');
  static const ComparisonOperator COMPARISON_OPERATOR_NOT_CONTAINS =
      ComparisonOperator._(
          7, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_NOT_CONTAINS');
  static const ComparisonOperator COMPARISON_OPERATOR_BETWEEN =
      ComparisonOperator._(
          8, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_BETWEEN');
  static const ComparisonOperator COMPARISON_OPERATOR_NULL =
      ComparisonOperator._(9, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_NULL');
  static const ComparisonOperator COMPARISON_OPERATOR_BEGINS_WITH =
      ComparisonOperator._(
          10, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_BEGINS_WITH');
  static const ComparisonOperator COMPARISON_OPERATOR_LE =
      ComparisonOperator._(11, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_LE');
  static const ComparisonOperator COMPARISON_OPERATOR_NE =
      ComparisonOperator._(12, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_NE');

  static const $core.List<ComparisonOperator> values = <ComparisonOperator>[
    COMPARISON_OPERATOR_GE,
    COMPARISON_OPERATOR_NOT_NULL,
    COMPARISON_OPERATOR_EQ,
    COMPARISON_OPERATOR_IN,
    COMPARISON_OPERATOR_LT,
    COMPARISON_OPERATOR_CONTAINS,
    COMPARISON_OPERATOR_GT,
    COMPARISON_OPERATOR_NOT_CONTAINS,
    COMPARISON_OPERATOR_BETWEEN,
    COMPARISON_OPERATOR_NULL,
    COMPARISON_OPERATOR_BEGINS_WITH,
    COMPARISON_OPERATOR_LE,
    COMPARISON_OPERATOR_NE,
  ];

  static final $core.List<ComparisonOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 12);
  static ComparisonOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComparisonOperator._(super.value, super.name);
}

class ConditionalOperator extends $pb.ProtobufEnum {
  static const ConditionalOperator CONDITIONAL_OPERATOR_AND =
      ConditionalOperator._(
          0, _omitEnumNames ? '' : 'CONDITIONAL_OPERATOR_AND');
  static const ConditionalOperator CONDITIONAL_OPERATOR_OR =
      ConditionalOperator._(1, _omitEnumNames ? '' : 'CONDITIONAL_OPERATOR_OR');

  static const $core.List<ConditionalOperator> values = <ConditionalOperator>[
    CONDITIONAL_OPERATOR_AND,
    CONDITIONAL_OPERATOR_OR,
  ];

  static final $core.List<ConditionalOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ConditionalOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConditionalOperator._(super.value, super.name);
}

class ContinuousBackupsStatus extends $pb.ProtobufEnum {
  static const ContinuousBackupsStatus CONTINUOUS_BACKUPS_STATUS_DISABLED =
      ContinuousBackupsStatus._(
          0, _omitEnumNames ? '' : 'CONTINUOUS_BACKUPS_STATUS_DISABLED');
  static const ContinuousBackupsStatus CONTINUOUS_BACKUPS_STATUS_ENABLED =
      ContinuousBackupsStatus._(
          1, _omitEnumNames ? '' : 'CONTINUOUS_BACKUPS_STATUS_ENABLED');

  static const $core.List<ContinuousBackupsStatus> values =
      <ContinuousBackupsStatus>[
    CONTINUOUS_BACKUPS_STATUS_DISABLED,
    CONTINUOUS_BACKUPS_STATUS_ENABLED,
  ];

  static final $core.List<ContinuousBackupsStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ContinuousBackupsStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContinuousBackupsStatus._(super.value, super.name);
}

class ContributorInsightsAction extends $pb.ProtobufEnum {
  static const ContributorInsightsAction CONTRIBUTOR_INSIGHTS_ACTION_ENABLE =
      ContributorInsightsAction._(
          0, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_ACTION_ENABLE');
  static const ContributorInsightsAction CONTRIBUTOR_INSIGHTS_ACTION_DISABLE =
      ContributorInsightsAction._(
          1, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_ACTION_DISABLE');

  static const $core.List<ContributorInsightsAction> values =
      <ContributorInsightsAction>[
    CONTRIBUTOR_INSIGHTS_ACTION_ENABLE,
    CONTRIBUTOR_INSIGHTS_ACTION_DISABLE,
  ];

  static final $core.List<ContributorInsightsAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ContributorInsightsAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContributorInsightsAction._(super.value, super.name);
}

class ContributorInsightsMode extends $pb.ProtobufEnum {
  static const ContributorInsightsMode
      CONTRIBUTOR_INSIGHTS_MODE_THROTTLED_KEYS = ContributorInsightsMode._(
          0, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_MODE_THROTTLED_KEYS');
  static const ContributorInsightsMode
      CONTRIBUTOR_INSIGHTS_MODE_ACCESSED_AND_THROTTLED_KEYS =
      ContributorInsightsMode._(
          1,
          _omitEnumNames
              ? ''
              : 'CONTRIBUTOR_INSIGHTS_MODE_ACCESSED_AND_THROTTLED_KEYS');

  static const $core.List<ContributorInsightsMode> values =
      <ContributorInsightsMode>[
    CONTRIBUTOR_INSIGHTS_MODE_THROTTLED_KEYS,
    CONTRIBUTOR_INSIGHTS_MODE_ACCESSED_AND_THROTTLED_KEYS,
  ];

  static final $core.List<ContributorInsightsMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ContributorInsightsMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContributorInsightsMode._(super.value, super.name);
}

class ContributorInsightsStatus extends $pb.ProtobufEnum {
  static const ContributorInsightsStatus CONTRIBUTOR_INSIGHTS_STATUS_DISABLED =
      ContributorInsightsStatus._(
          0, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_STATUS_DISABLED');
  static const ContributorInsightsStatus CONTRIBUTOR_INSIGHTS_STATUS_ENABLING =
      ContributorInsightsStatus._(
          1, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_STATUS_ENABLING');
  static const ContributorInsightsStatus CONTRIBUTOR_INSIGHTS_STATUS_DISABLING =
      ContributorInsightsStatus._(
          2, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_STATUS_DISABLING');
  static const ContributorInsightsStatus CONTRIBUTOR_INSIGHTS_STATUS_ENABLED =
      ContributorInsightsStatus._(
          3, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_STATUS_ENABLED');
  static const ContributorInsightsStatus CONTRIBUTOR_INSIGHTS_STATUS_FAILED =
      ContributorInsightsStatus._(
          4, _omitEnumNames ? '' : 'CONTRIBUTOR_INSIGHTS_STATUS_FAILED');

  static const $core.List<ContributorInsightsStatus> values =
      <ContributorInsightsStatus>[
    CONTRIBUTOR_INSIGHTS_STATUS_DISABLED,
    CONTRIBUTOR_INSIGHTS_STATUS_ENABLING,
    CONTRIBUTOR_INSIGHTS_STATUS_DISABLING,
    CONTRIBUTOR_INSIGHTS_STATUS_ENABLED,
    CONTRIBUTOR_INSIGHTS_STATUS_FAILED,
  ];

  static final $core.List<ContributorInsightsStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ContributorInsightsStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContributorInsightsStatus._(super.value, super.name);
}

class DestinationStatus extends $pb.ProtobufEnum {
  static const DestinationStatus DESTINATION_STATUS_UPDATING =
      DestinationStatus._(
          0, _omitEnumNames ? '' : 'DESTINATION_STATUS_UPDATING');
  static const DestinationStatus DESTINATION_STATUS_DISABLED =
      DestinationStatus._(
          1, _omitEnumNames ? '' : 'DESTINATION_STATUS_DISABLED');
  static const DestinationStatus DESTINATION_STATUS_ENABLING =
      DestinationStatus._(
          2, _omitEnumNames ? '' : 'DESTINATION_STATUS_ENABLING');
  static const DestinationStatus DESTINATION_STATUS_ENABLE_FAILED =
      DestinationStatus._(
          3, _omitEnumNames ? '' : 'DESTINATION_STATUS_ENABLE_FAILED');
  static const DestinationStatus DESTINATION_STATUS_ACTIVE =
      DestinationStatus._(4, _omitEnumNames ? '' : 'DESTINATION_STATUS_ACTIVE');
  static const DestinationStatus DESTINATION_STATUS_DISABLING =
      DestinationStatus._(
          5, _omitEnumNames ? '' : 'DESTINATION_STATUS_DISABLING');

  static const $core.List<DestinationStatus> values = <DestinationStatus>[
    DESTINATION_STATUS_UPDATING,
    DESTINATION_STATUS_DISABLED,
    DESTINATION_STATUS_ENABLING,
    DESTINATION_STATUS_ENABLE_FAILED,
    DESTINATION_STATUS_ACTIVE,
    DESTINATION_STATUS_DISABLING,
  ];

  static final $core.List<DestinationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static DestinationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DestinationStatus._(super.value, super.name);
}

class ExportFormat extends $pb.ProtobufEnum {
  static const ExportFormat EXPORT_FORMAT_DYNAMODB_JSON =
      ExportFormat._(0, _omitEnumNames ? '' : 'EXPORT_FORMAT_DYNAMODB_JSON');
  static const ExportFormat EXPORT_FORMAT_ION =
      ExportFormat._(1, _omitEnumNames ? '' : 'EXPORT_FORMAT_ION');

  static const $core.List<ExportFormat> values = <ExportFormat>[
    EXPORT_FORMAT_DYNAMODB_JSON,
    EXPORT_FORMAT_ION,
  ];

  static final $core.List<ExportFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExportFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExportFormat._(super.value, super.name);
}

class ExportStatus extends $pb.ProtobufEnum {
  static const ExportStatus EXPORT_STATUS_IN_PROGRESS =
      ExportStatus._(0, _omitEnumNames ? '' : 'EXPORT_STATUS_IN_PROGRESS');
  static const ExportStatus EXPORT_STATUS_COMPLETED =
      ExportStatus._(1, _omitEnumNames ? '' : 'EXPORT_STATUS_COMPLETED');
  static const ExportStatus EXPORT_STATUS_FAILED =
      ExportStatus._(2, _omitEnumNames ? '' : 'EXPORT_STATUS_FAILED');

  static const $core.List<ExportStatus> values = <ExportStatus>[
    EXPORT_STATUS_IN_PROGRESS,
    EXPORT_STATUS_COMPLETED,
    EXPORT_STATUS_FAILED,
  ];

  static final $core.List<ExportStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ExportStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExportStatus._(super.value, super.name);
}

class ExportType extends $pb.ProtobufEnum {
  static const ExportType EXPORT_TYPE_FULL_EXPORT =
      ExportType._(0, _omitEnumNames ? '' : 'EXPORT_TYPE_FULL_EXPORT');
  static const ExportType EXPORT_TYPE_INCREMENTAL_EXPORT =
      ExportType._(1, _omitEnumNames ? '' : 'EXPORT_TYPE_INCREMENTAL_EXPORT');

  static const $core.List<ExportType> values = <ExportType>[
    EXPORT_TYPE_FULL_EXPORT,
    EXPORT_TYPE_INCREMENTAL_EXPORT,
  ];

  static final $core.List<ExportType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExportType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExportType._(super.value, super.name);
}

class ExportViewType extends $pb.ProtobufEnum {
  static const ExportViewType EXPORT_VIEW_TYPE_NEW_IMAGE =
      ExportViewType._(0, _omitEnumNames ? '' : 'EXPORT_VIEW_TYPE_NEW_IMAGE');
  static const ExportViewType EXPORT_VIEW_TYPE_NEW_AND_OLD_IMAGES =
      ExportViewType._(
          1, _omitEnumNames ? '' : 'EXPORT_VIEW_TYPE_NEW_AND_OLD_IMAGES');

  static const $core.List<ExportViewType> values = <ExportViewType>[
    EXPORT_VIEW_TYPE_NEW_IMAGE,
    EXPORT_VIEW_TYPE_NEW_AND_OLD_IMAGES,
  ];

  static final $core.List<ExportViewType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExportViewType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExportViewType._(super.value, super.name);
}

class GlobalTableSettingsReplicationMode extends $pb.ProtobufEnum {
  static const GlobalTableSettingsReplicationMode
      GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_DISABLED =
      GlobalTableSettingsReplicationMode._(
          0,
          _omitEnumNames
              ? ''
              : 'GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_DISABLED');
  static const GlobalTableSettingsReplicationMode
      GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED_WITH_OVERRIDES =
      GlobalTableSettingsReplicationMode._(
          1,
          _omitEnumNames
              ? ''
              : 'GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED_WITH_OVERRIDES');
  static const GlobalTableSettingsReplicationMode
      GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED =
      GlobalTableSettingsReplicationMode._(
          2,
          _omitEnumNames
              ? ''
              : 'GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED');

  static const $core.List<GlobalTableSettingsReplicationMode> values =
      <GlobalTableSettingsReplicationMode>[
    GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_DISABLED,
    GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED_WITH_OVERRIDES,
    GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED,
  ];

  static final $core.List<GlobalTableSettingsReplicationMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static GlobalTableSettingsReplicationMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GlobalTableSettingsReplicationMode._(super.value, super.name);
}

class GlobalTableStatus extends $pb.ProtobufEnum {
  static const GlobalTableStatus GLOBAL_TABLE_STATUS_UPDATING =
      GlobalTableStatus._(
          0, _omitEnumNames ? '' : 'GLOBAL_TABLE_STATUS_UPDATING');
  static const GlobalTableStatus GLOBAL_TABLE_STATUS_ACTIVE =
      GlobalTableStatus._(
          1, _omitEnumNames ? '' : 'GLOBAL_TABLE_STATUS_ACTIVE');
  static const GlobalTableStatus GLOBAL_TABLE_STATUS_DELETING =
      GlobalTableStatus._(
          2, _omitEnumNames ? '' : 'GLOBAL_TABLE_STATUS_DELETING');
  static const GlobalTableStatus GLOBAL_TABLE_STATUS_CREATING =
      GlobalTableStatus._(
          3, _omitEnumNames ? '' : 'GLOBAL_TABLE_STATUS_CREATING');

  static const $core.List<GlobalTableStatus> values = <GlobalTableStatus>[
    GLOBAL_TABLE_STATUS_UPDATING,
    GLOBAL_TABLE_STATUS_ACTIVE,
    GLOBAL_TABLE_STATUS_DELETING,
    GLOBAL_TABLE_STATUS_CREATING,
  ];

  static final $core.List<GlobalTableStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static GlobalTableStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GlobalTableStatus._(super.value, super.name);
}

class ImportStatus extends $pb.ProtobufEnum {
  static const ImportStatus IMPORT_STATUS_IN_PROGRESS =
      ImportStatus._(0, _omitEnumNames ? '' : 'IMPORT_STATUS_IN_PROGRESS');
  static const ImportStatus IMPORT_STATUS_CANCELLED =
      ImportStatus._(1, _omitEnumNames ? '' : 'IMPORT_STATUS_CANCELLED');
  static const ImportStatus IMPORT_STATUS_CANCELLING =
      ImportStatus._(2, _omitEnumNames ? '' : 'IMPORT_STATUS_CANCELLING');
  static const ImportStatus IMPORT_STATUS_COMPLETED =
      ImportStatus._(3, _omitEnumNames ? '' : 'IMPORT_STATUS_COMPLETED');
  static const ImportStatus IMPORT_STATUS_FAILED =
      ImportStatus._(4, _omitEnumNames ? '' : 'IMPORT_STATUS_FAILED');

  static const $core.List<ImportStatus> values = <ImportStatus>[
    IMPORT_STATUS_IN_PROGRESS,
    IMPORT_STATUS_CANCELLED,
    IMPORT_STATUS_CANCELLING,
    IMPORT_STATUS_COMPLETED,
    IMPORT_STATUS_FAILED,
  ];

  static final $core.List<ImportStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ImportStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportStatus._(super.value, super.name);
}

class IndexStatus extends $pb.ProtobufEnum {
  static const IndexStatus INDEX_STATUS_UPDATING =
      IndexStatus._(0, _omitEnumNames ? '' : 'INDEX_STATUS_UPDATING');
  static const IndexStatus INDEX_STATUS_ACTIVE =
      IndexStatus._(1, _omitEnumNames ? '' : 'INDEX_STATUS_ACTIVE');
  static const IndexStatus INDEX_STATUS_DELETING =
      IndexStatus._(2, _omitEnumNames ? '' : 'INDEX_STATUS_DELETING');
  static const IndexStatus INDEX_STATUS_CREATING =
      IndexStatus._(3, _omitEnumNames ? '' : 'INDEX_STATUS_CREATING');

  static const $core.List<IndexStatus> values = <IndexStatus>[
    INDEX_STATUS_UPDATING,
    INDEX_STATUS_ACTIVE,
    INDEX_STATUS_DELETING,
    INDEX_STATUS_CREATING,
  ];

  static final $core.List<IndexStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static IndexStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IndexStatus._(super.value, super.name);
}

class InputCompressionType extends $pb.ProtobufEnum {
  static const InputCompressionType INPUT_COMPRESSION_TYPE_NONE =
      InputCompressionType._(
          0, _omitEnumNames ? '' : 'INPUT_COMPRESSION_TYPE_NONE');
  static const InputCompressionType INPUT_COMPRESSION_TYPE_GZIP =
      InputCompressionType._(
          1, _omitEnumNames ? '' : 'INPUT_COMPRESSION_TYPE_GZIP');
  static const InputCompressionType INPUT_COMPRESSION_TYPE_ZSTD =
      InputCompressionType._(
          2, _omitEnumNames ? '' : 'INPUT_COMPRESSION_TYPE_ZSTD');

  static const $core.List<InputCompressionType> values = <InputCompressionType>[
    INPUT_COMPRESSION_TYPE_NONE,
    INPUT_COMPRESSION_TYPE_GZIP,
    INPUT_COMPRESSION_TYPE_ZSTD,
  ];

  static final $core.List<InputCompressionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static InputCompressionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InputCompressionType._(super.value, super.name);
}

class InputFormat extends $pb.ProtobufEnum {
  static const InputFormat INPUT_FORMAT_DYNAMODB_JSON =
      InputFormat._(0, _omitEnumNames ? '' : 'INPUT_FORMAT_DYNAMODB_JSON');
  static const InputFormat INPUT_FORMAT_CSV =
      InputFormat._(1, _omitEnumNames ? '' : 'INPUT_FORMAT_CSV');
  static const InputFormat INPUT_FORMAT_ION =
      InputFormat._(2, _omitEnumNames ? '' : 'INPUT_FORMAT_ION');

  static const $core.List<InputFormat> values = <InputFormat>[
    INPUT_FORMAT_DYNAMODB_JSON,
    INPUT_FORMAT_CSV,
    INPUT_FORMAT_ION,
  ];

  static final $core.List<InputFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static InputFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InputFormat._(super.value, super.name);
}

class KeyType extends $pb.ProtobufEnum {
  static const KeyType KEY_TYPE_HASH =
      KeyType._(0, _omitEnumNames ? '' : 'KEY_TYPE_HASH');
  static const KeyType KEY_TYPE_RANGE =
      KeyType._(1, _omitEnumNames ? '' : 'KEY_TYPE_RANGE');

  static const $core.List<KeyType> values = <KeyType>[
    KEY_TYPE_HASH,
    KEY_TYPE_RANGE,
  ];

  static final $core.List<KeyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static KeyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyType._(super.value, super.name);
}

class MultiRegionConsistency extends $pb.ProtobufEnum {
  static const MultiRegionConsistency MULTI_REGION_CONSISTENCY_EVENTUAL =
      MultiRegionConsistency._(
          0, _omitEnumNames ? '' : 'MULTI_REGION_CONSISTENCY_EVENTUAL');
  static const MultiRegionConsistency MULTI_REGION_CONSISTENCY_STRONG =
      MultiRegionConsistency._(
          1, _omitEnumNames ? '' : 'MULTI_REGION_CONSISTENCY_STRONG');

  static const $core.List<MultiRegionConsistency> values =
      <MultiRegionConsistency>[
    MULTI_REGION_CONSISTENCY_EVENTUAL,
    MULTI_REGION_CONSISTENCY_STRONG,
  ];

  static final $core.List<MultiRegionConsistency?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MultiRegionConsistency? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MultiRegionConsistency._(super.value, super.name);
}

class PointInTimeRecoveryStatus extends $pb.ProtobufEnum {
  static const PointInTimeRecoveryStatus
      POINT_IN_TIME_RECOVERY_STATUS_DISABLED = PointInTimeRecoveryStatus._(
          0, _omitEnumNames ? '' : 'POINT_IN_TIME_RECOVERY_STATUS_DISABLED');
  static const PointInTimeRecoveryStatus POINT_IN_TIME_RECOVERY_STATUS_ENABLED =
      PointInTimeRecoveryStatus._(
          1, _omitEnumNames ? '' : 'POINT_IN_TIME_RECOVERY_STATUS_ENABLED');

  static const $core.List<PointInTimeRecoveryStatus> values =
      <PointInTimeRecoveryStatus>[
    POINT_IN_TIME_RECOVERY_STATUS_DISABLED,
    POINT_IN_TIME_RECOVERY_STATUS_ENABLED,
  ];

  static final $core.List<PointInTimeRecoveryStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PointInTimeRecoveryStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PointInTimeRecoveryStatus._(super.value, super.name);
}

class ProjectionType extends $pb.ProtobufEnum {
  static const ProjectionType PROJECTION_TYPE_KEYS_ONLY =
      ProjectionType._(0, _omitEnumNames ? '' : 'PROJECTION_TYPE_KEYS_ONLY');
  static const ProjectionType PROJECTION_TYPE_INCLUDE =
      ProjectionType._(1, _omitEnumNames ? '' : 'PROJECTION_TYPE_INCLUDE');
  static const ProjectionType PROJECTION_TYPE_ALL =
      ProjectionType._(2, _omitEnumNames ? '' : 'PROJECTION_TYPE_ALL');

  static const $core.List<ProjectionType> values = <ProjectionType>[
    PROJECTION_TYPE_KEYS_ONLY,
    PROJECTION_TYPE_INCLUDE,
    PROJECTION_TYPE_ALL,
  ];

  static final $core.List<ProjectionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ProjectionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ProjectionType._(super.value, super.name);
}

class ReplicaStatus extends $pb.ProtobufEnum {
  static const ReplicaStatus REPLICA_STATUS_UPDATING =
      ReplicaStatus._(0, _omitEnumNames ? '' : 'REPLICA_STATUS_UPDATING');
  static const ReplicaStatus REPLICA_STATUS_ARCHIVING =
      ReplicaStatus._(1, _omitEnumNames ? '' : 'REPLICA_STATUS_ARCHIVING');
  static const ReplicaStatus REPLICA_STATUS_REGION_DISABLED = ReplicaStatus._(
      2, _omitEnumNames ? '' : 'REPLICA_STATUS_REGION_DISABLED');
  static const ReplicaStatus
      REPLICA_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS = ReplicaStatus._(
          3,
          _omitEnumNames
              ? ''
              : 'REPLICA_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS');
  static const ReplicaStatus REPLICA_STATUS_ARCHIVED =
      ReplicaStatus._(4, _omitEnumNames ? '' : 'REPLICA_STATUS_ARCHIVED');
  static const ReplicaStatus REPLICA_STATUS_ACTIVE =
      ReplicaStatus._(5, _omitEnumNames ? '' : 'REPLICA_STATUS_ACTIVE');
  static const ReplicaStatus REPLICA_STATUS_DELETING =
      ReplicaStatus._(6, _omitEnumNames ? '' : 'REPLICA_STATUS_DELETING');
  static const ReplicaStatus REPLICA_STATUS_CREATION_FAILED = ReplicaStatus._(
      7, _omitEnumNames ? '' : 'REPLICA_STATUS_CREATION_FAILED');
  static const ReplicaStatus REPLICA_STATUS_CREATING =
      ReplicaStatus._(8, _omitEnumNames ? '' : 'REPLICA_STATUS_CREATING');
  static const ReplicaStatus REPLICA_STATUS_REPLICATION_NOT_AUTHORIZED =
      ReplicaStatus._(
          9, _omitEnumNames ? '' : 'REPLICA_STATUS_REPLICATION_NOT_AUTHORIZED');

  static const $core.List<ReplicaStatus> values = <ReplicaStatus>[
    REPLICA_STATUS_UPDATING,
    REPLICA_STATUS_ARCHIVING,
    REPLICA_STATUS_REGION_DISABLED,
    REPLICA_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS,
    REPLICA_STATUS_ARCHIVED,
    REPLICA_STATUS_ACTIVE,
    REPLICA_STATUS_DELETING,
    REPLICA_STATUS_CREATION_FAILED,
    REPLICA_STATUS_CREATING,
    REPLICA_STATUS_REPLICATION_NOT_AUTHORIZED,
  ];

  static final $core.List<ReplicaStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 9);
  static ReplicaStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReplicaStatus._(super.value, super.name);
}

class ReturnConsumedCapacity extends $pb.ProtobufEnum {
  static const ReturnConsumedCapacity RETURN_CONSUMED_CAPACITY_NONE =
      ReturnConsumedCapacity._(
          0, _omitEnumNames ? '' : 'RETURN_CONSUMED_CAPACITY_NONE');
  static const ReturnConsumedCapacity RETURN_CONSUMED_CAPACITY_INDEXES =
      ReturnConsumedCapacity._(
          1, _omitEnumNames ? '' : 'RETURN_CONSUMED_CAPACITY_INDEXES');
  static const ReturnConsumedCapacity RETURN_CONSUMED_CAPACITY_TOTAL =
      ReturnConsumedCapacity._(
          2, _omitEnumNames ? '' : 'RETURN_CONSUMED_CAPACITY_TOTAL');

  static const $core.List<ReturnConsumedCapacity> values =
      <ReturnConsumedCapacity>[
    RETURN_CONSUMED_CAPACITY_NONE,
    RETURN_CONSUMED_CAPACITY_INDEXES,
    RETURN_CONSUMED_CAPACITY_TOTAL,
  ];

  static final $core.List<ReturnConsumedCapacity?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ReturnConsumedCapacity? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReturnConsumedCapacity._(super.value, super.name);
}

class ReturnItemCollectionMetrics extends $pb.ProtobufEnum {
  static const ReturnItemCollectionMetrics RETURN_ITEM_COLLECTION_METRICS_NONE =
      ReturnItemCollectionMetrics._(
          0, _omitEnumNames ? '' : 'RETURN_ITEM_COLLECTION_METRICS_NONE');
  static const ReturnItemCollectionMetrics RETURN_ITEM_COLLECTION_METRICS_SIZE =
      ReturnItemCollectionMetrics._(
          1, _omitEnumNames ? '' : 'RETURN_ITEM_COLLECTION_METRICS_SIZE');

  static const $core.List<ReturnItemCollectionMetrics> values =
      <ReturnItemCollectionMetrics>[
    RETURN_ITEM_COLLECTION_METRICS_NONE,
    RETURN_ITEM_COLLECTION_METRICS_SIZE,
  ];

  static final $core.List<ReturnItemCollectionMetrics?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ReturnItemCollectionMetrics? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReturnItemCollectionMetrics._(super.value, super.name);
}

class ReturnValue extends $pb.ProtobufEnum {
  static const ReturnValue RETURN_VALUE_UPDATED_NEW =
      ReturnValue._(0, _omitEnumNames ? '' : 'RETURN_VALUE_UPDATED_NEW');
  static const ReturnValue RETURN_VALUE_NONE =
      ReturnValue._(1, _omitEnumNames ? '' : 'RETURN_VALUE_NONE');
  static const ReturnValue RETURN_VALUE_UPDATED_OLD =
      ReturnValue._(2, _omitEnumNames ? '' : 'RETURN_VALUE_UPDATED_OLD');
  static const ReturnValue RETURN_VALUE_ALL_OLD =
      ReturnValue._(3, _omitEnumNames ? '' : 'RETURN_VALUE_ALL_OLD');
  static const ReturnValue RETURN_VALUE_ALL_NEW =
      ReturnValue._(4, _omitEnumNames ? '' : 'RETURN_VALUE_ALL_NEW');

  static const $core.List<ReturnValue> values = <ReturnValue>[
    RETURN_VALUE_UPDATED_NEW,
    RETURN_VALUE_NONE,
    RETURN_VALUE_UPDATED_OLD,
    RETURN_VALUE_ALL_OLD,
    RETURN_VALUE_ALL_NEW,
  ];

  static final $core.List<ReturnValue?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ReturnValue? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReturnValue._(super.value, super.name);
}

class ReturnValuesOnConditionCheckFailure extends $pb.ProtobufEnum {
  static const ReturnValuesOnConditionCheckFailure
      RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_NONE =
      ReturnValuesOnConditionCheckFailure._(
          0,
          _omitEnumNames
              ? ''
              : 'RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_NONE');
  static const ReturnValuesOnConditionCheckFailure
      RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_ALL_OLD =
      ReturnValuesOnConditionCheckFailure._(
          1,
          _omitEnumNames
              ? ''
              : 'RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_ALL_OLD');

  static const $core.List<ReturnValuesOnConditionCheckFailure> values =
      <ReturnValuesOnConditionCheckFailure>[
    RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_NONE,
    RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_ALL_OLD,
  ];

  static final $core.List<ReturnValuesOnConditionCheckFailure?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ReturnValuesOnConditionCheckFailure? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReturnValuesOnConditionCheckFailure._(super.value, super.name);
}

class S3SseAlgorithm extends $pb.ProtobufEnum {
  static const S3SseAlgorithm S3_SSE_ALGORITHM_AES256 =
      S3SseAlgorithm._(0, _omitEnumNames ? '' : 'S3_SSE_ALGORITHM_AES256');
  static const S3SseAlgorithm S3_SSE_ALGORITHM_KMS =
      S3SseAlgorithm._(1, _omitEnumNames ? '' : 'S3_SSE_ALGORITHM_KMS');

  static const $core.List<S3SseAlgorithm> values = <S3SseAlgorithm>[
    S3_SSE_ALGORITHM_AES256,
    S3_SSE_ALGORITHM_KMS,
  ];

  static final $core.List<S3SseAlgorithm?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static S3SseAlgorithm? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const S3SseAlgorithm._(super.value, super.name);
}

class SSEStatus extends $pb.ProtobufEnum {
  static const SSEStatus S_S_E_STATUS_UPDATING =
      SSEStatus._(0, _omitEnumNames ? '' : 'S_S_E_STATUS_UPDATING');
  static const SSEStatus S_S_E_STATUS_DISABLED =
      SSEStatus._(1, _omitEnumNames ? '' : 'S_S_E_STATUS_DISABLED');
  static const SSEStatus S_S_E_STATUS_ENABLING =
      SSEStatus._(2, _omitEnumNames ? '' : 'S_S_E_STATUS_ENABLING');
  static const SSEStatus S_S_E_STATUS_DISABLING =
      SSEStatus._(3, _omitEnumNames ? '' : 'S_S_E_STATUS_DISABLING');
  static const SSEStatus S_S_E_STATUS_ENABLED =
      SSEStatus._(4, _omitEnumNames ? '' : 'S_S_E_STATUS_ENABLED');

  static const $core.List<SSEStatus> values = <SSEStatus>[
    S_S_E_STATUS_UPDATING,
    S_S_E_STATUS_DISABLED,
    S_S_E_STATUS_ENABLING,
    S_S_E_STATUS_DISABLING,
    S_S_E_STATUS_ENABLED,
  ];

  static final $core.List<SSEStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static SSEStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SSEStatus._(super.value, super.name);
}

class SSEType extends $pb.ProtobufEnum {
  static const SSEType S_S_E_TYPE_AES256 =
      SSEType._(0, _omitEnumNames ? '' : 'S_S_E_TYPE_AES256');
  static const SSEType S_S_E_TYPE_KMS =
      SSEType._(1, _omitEnumNames ? '' : 'S_S_E_TYPE_KMS');

  static const $core.List<SSEType> values = <SSEType>[
    S_S_E_TYPE_AES256,
    S_S_E_TYPE_KMS,
  ];

  static final $core.List<SSEType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SSEType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SSEType._(super.value, super.name);
}

class ScalarAttributeType extends $pb.ProtobufEnum {
  static const ScalarAttributeType SCALAR_ATTRIBUTE_TYPE_B =
      ScalarAttributeType._(0, _omitEnumNames ? '' : 'SCALAR_ATTRIBUTE_TYPE_B');
  static const ScalarAttributeType SCALAR_ATTRIBUTE_TYPE_S =
      ScalarAttributeType._(1, _omitEnumNames ? '' : 'SCALAR_ATTRIBUTE_TYPE_S');
  static const ScalarAttributeType SCALAR_ATTRIBUTE_TYPE_N =
      ScalarAttributeType._(2, _omitEnumNames ? '' : 'SCALAR_ATTRIBUTE_TYPE_N');

  static const $core.List<ScalarAttributeType> values = <ScalarAttributeType>[
    SCALAR_ATTRIBUTE_TYPE_B,
    SCALAR_ATTRIBUTE_TYPE_S,
    SCALAR_ATTRIBUTE_TYPE_N,
  ];

  static final $core.List<ScalarAttributeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ScalarAttributeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScalarAttributeType._(super.value, super.name);
}

class Select extends $pb.ProtobufEnum {
  static const Select SELECT_ALL_ATTRIBUTES =
      Select._(0, _omitEnumNames ? '' : 'SELECT_ALL_ATTRIBUTES');
  static const Select SELECT_COUNT =
      Select._(1, _omitEnumNames ? '' : 'SELECT_COUNT');
  static const Select SELECT_ALL_PROJECTED_ATTRIBUTES =
      Select._(2, _omitEnumNames ? '' : 'SELECT_ALL_PROJECTED_ATTRIBUTES');
  static const Select SELECT_SPECIFIC_ATTRIBUTES =
      Select._(3, _omitEnumNames ? '' : 'SELECT_SPECIFIC_ATTRIBUTES');

  static const $core.List<Select> values = <Select>[
    SELECT_ALL_ATTRIBUTES,
    SELECT_COUNT,
    SELECT_ALL_PROJECTED_ATTRIBUTES,
    SELECT_SPECIFIC_ATTRIBUTES,
  ];

  static final $core.List<Select?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static Select? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Select._(super.value, super.name);
}

class StreamViewType extends $pb.ProtobufEnum {
  static const StreamViewType STREAM_VIEW_TYPE_KEYS_ONLY =
      StreamViewType._(0, _omitEnumNames ? '' : 'STREAM_VIEW_TYPE_KEYS_ONLY');
  static const StreamViewType STREAM_VIEW_TYPE_NEW_IMAGE =
      StreamViewType._(1, _omitEnumNames ? '' : 'STREAM_VIEW_TYPE_NEW_IMAGE');
  static const StreamViewType STREAM_VIEW_TYPE_NEW_AND_OLD_IMAGES =
      StreamViewType._(
          2, _omitEnumNames ? '' : 'STREAM_VIEW_TYPE_NEW_AND_OLD_IMAGES');
  static const StreamViewType STREAM_VIEW_TYPE_OLD_IMAGE =
      StreamViewType._(3, _omitEnumNames ? '' : 'STREAM_VIEW_TYPE_OLD_IMAGE');

  static const $core.List<StreamViewType> values = <StreamViewType>[
    STREAM_VIEW_TYPE_KEYS_ONLY,
    STREAM_VIEW_TYPE_NEW_IMAGE,
    STREAM_VIEW_TYPE_NEW_AND_OLD_IMAGES,
    STREAM_VIEW_TYPE_OLD_IMAGE,
  ];

  static final $core.List<StreamViewType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static StreamViewType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StreamViewType._(super.value, super.name);
}

class TableClass extends $pb.ProtobufEnum {
  static const TableClass TABLE_CLASS_STANDARD_INFREQUENT_ACCESS = TableClass._(
      0, _omitEnumNames ? '' : 'TABLE_CLASS_STANDARD_INFREQUENT_ACCESS');
  static const TableClass TABLE_CLASS_STANDARD =
      TableClass._(1, _omitEnumNames ? '' : 'TABLE_CLASS_STANDARD');

  static const $core.List<TableClass> values = <TableClass>[
    TABLE_CLASS_STANDARD_INFREQUENT_ACCESS,
    TABLE_CLASS_STANDARD,
  ];

  static final $core.List<TableClass?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static TableClass? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TableClass._(super.value, super.name);
}

class TableStatus extends $pb.ProtobufEnum {
  static const TableStatus TABLE_STATUS_UPDATING =
      TableStatus._(0, _omitEnumNames ? '' : 'TABLE_STATUS_UPDATING');
  static const TableStatus TABLE_STATUS_ARCHIVING =
      TableStatus._(1, _omitEnumNames ? '' : 'TABLE_STATUS_ARCHIVING');
  static const TableStatus TABLE_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS =
      TableStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'TABLE_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS');
  static const TableStatus TABLE_STATUS_ARCHIVED =
      TableStatus._(3, _omitEnumNames ? '' : 'TABLE_STATUS_ARCHIVED');
  static const TableStatus TABLE_STATUS_ACTIVE =
      TableStatus._(4, _omitEnumNames ? '' : 'TABLE_STATUS_ACTIVE');
  static const TableStatus TABLE_STATUS_DELETING =
      TableStatus._(5, _omitEnumNames ? '' : 'TABLE_STATUS_DELETING');
  static const TableStatus TABLE_STATUS_CREATING =
      TableStatus._(6, _omitEnumNames ? '' : 'TABLE_STATUS_CREATING');
  static const TableStatus TABLE_STATUS_REPLICATION_NOT_AUTHORIZED =
      TableStatus._(
          7, _omitEnumNames ? '' : 'TABLE_STATUS_REPLICATION_NOT_AUTHORIZED');

  static const $core.List<TableStatus> values = <TableStatus>[
    TABLE_STATUS_UPDATING,
    TABLE_STATUS_ARCHIVING,
    TABLE_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS,
    TABLE_STATUS_ARCHIVED,
    TABLE_STATUS_ACTIVE,
    TABLE_STATUS_DELETING,
    TABLE_STATUS_CREATING,
    TABLE_STATUS_REPLICATION_NOT_AUTHORIZED,
  ];

  static final $core.List<TableStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static TableStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TableStatus._(super.value, super.name);
}

class TimeToLiveStatus extends $pb.ProtobufEnum {
  static const TimeToLiveStatus TIME_TO_LIVE_STATUS_DISABLED =
      TimeToLiveStatus._(
          0, _omitEnumNames ? '' : 'TIME_TO_LIVE_STATUS_DISABLED');
  static const TimeToLiveStatus TIME_TO_LIVE_STATUS_ENABLING =
      TimeToLiveStatus._(
          1, _omitEnumNames ? '' : 'TIME_TO_LIVE_STATUS_ENABLING');
  static const TimeToLiveStatus TIME_TO_LIVE_STATUS_DISABLING =
      TimeToLiveStatus._(
          2, _omitEnumNames ? '' : 'TIME_TO_LIVE_STATUS_DISABLING');
  static const TimeToLiveStatus TIME_TO_LIVE_STATUS_ENABLED =
      TimeToLiveStatus._(
          3, _omitEnumNames ? '' : 'TIME_TO_LIVE_STATUS_ENABLED');

  static const $core.List<TimeToLiveStatus> values = <TimeToLiveStatus>[
    TIME_TO_LIVE_STATUS_DISABLED,
    TIME_TO_LIVE_STATUS_ENABLING,
    TIME_TO_LIVE_STATUS_DISABLING,
    TIME_TO_LIVE_STATUS_ENABLED,
  ];

  static final $core.List<TimeToLiveStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static TimeToLiveStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TimeToLiveStatus._(super.value, super.name);
}

class WitnessStatus extends $pb.ProtobufEnum {
  static const WitnessStatus WITNESS_STATUS_ACTIVE =
      WitnessStatus._(0, _omitEnumNames ? '' : 'WITNESS_STATUS_ACTIVE');
  static const WitnessStatus WITNESS_STATUS_DELETING =
      WitnessStatus._(1, _omitEnumNames ? '' : 'WITNESS_STATUS_DELETING');
  static const WitnessStatus WITNESS_STATUS_CREATING =
      WitnessStatus._(2, _omitEnumNames ? '' : 'WITNESS_STATUS_CREATING');

  static const $core.List<WitnessStatus> values = <WitnessStatus>[
    WITNESS_STATUS_ACTIVE,
    WITNESS_STATUS_DELETING,
    WITNESS_STATUS_CREATING,
  ];

  static final $core.List<WitnessStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static WitnessStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WitnessStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
