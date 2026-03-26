// This is a generated file - do not edit.
//
// Generated from dynamodb.proto.

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

@$core.Deprecated('Use approximateCreationDateTimePrecisionDescriptor instead')
const ApproximateCreationDateTimePrecision$json = {
  '1': 'ApproximateCreationDateTimePrecision',
  '2': [
    {'1': 'APPROXIMATE_CREATION_DATE_TIME_PRECISION_MILLISECOND', '2': 0},
    {'1': 'APPROXIMATE_CREATION_DATE_TIME_PRECISION_MICROSECOND', '2': 1},
  ],
};

/// Descriptor for `ApproximateCreationDateTimePrecision`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List approximateCreationDateTimePrecisionDescriptor =
    $convert.base64Decode(
        'CiRBcHByb3hpbWF0ZUNyZWF0aW9uRGF0ZVRpbWVQcmVjaXNpb24SOAo0QVBQUk9YSU1BVEVfQ1'
        'JFQVRJT05fREFURV9USU1FX1BSRUNJU0lPTl9NSUxMSVNFQ09ORBAAEjgKNEFQUFJPWElNQVRF'
        'X0NSRUFUSU9OX0RBVEVfVElNRV9QUkVDSVNJT05fTUlDUk9TRUNPTkQQAQ==');

@$core.Deprecated('Use attributeActionDescriptor instead')
const AttributeAction$json = {
  '1': 'AttributeAction',
  '2': [
    {'1': 'ATTRIBUTE_ACTION_ADD', '2': 0},
    {'1': 'ATTRIBUTE_ACTION_DELETE', '2': 1},
    {'1': 'ATTRIBUTE_ACTION_PUT', '2': 2},
  ],
};

/// Descriptor for `AttributeAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List attributeActionDescriptor = $convert.base64Decode(
    'Cg9BdHRyaWJ1dGVBY3Rpb24SGAoUQVRUUklCVVRFX0FDVElPTl9BREQQABIbChdBVFRSSUJVVE'
    'VfQUNUSU9OX0RFTEVURRABEhgKFEFUVFJJQlVURV9BQ1RJT05fUFVUEAI=');

@$core.Deprecated('Use backupStatusDescriptor instead')
const BackupStatus$json = {
  '1': 'BackupStatus',
  '2': [
    {'1': 'BACKUP_STATUS_AVAILABLE', '2': 0},
    {'1': 'BACKUP_STATUS_CREATING', '2': 1},
    {'1': 'BACKUP_STATUS_DELETED', '2': 2},
  ],
};

/// Descriptor for `BackupStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List backupStatusDescriptor = $convert.base64Decode(
    'CgxCYWNrdXBTdGF0dXMSGwoXQkFDS1VQX1NUQVRVU19BVkFJTEFCTEUQABIaChZCQUNLVVBfU1'
    'RBVFVTX0NSRUFUSU5HEAESGQoVQkFDS1VQX1NUQVRVU19ERUxFVEVEEAI=');

@$core.Deprecated('Use backupTypeDescriptor instead')
const BackupType$json = {
  '1': 'BackupType',
  '2': [
    {'1': 'BACKUP_TYPE_SYSTEM', '2': 0},
    {'1': 'BACKUP_TYPE_USER', '2': 1},
    {'1': 'BACKUP_TYPE_AWS_BACKUP', '2': 2},
  ],
};

/// Descriptor for `BackupType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List backupTypeDescriptor = $convert.base64Decode(
    'CgpCYWNrdXBUeXBlEhYKEkJBQ0tVUF9UWVBFX1NZU1RFTRAAEhQKEEJBQ0tVUF9UWVBFX1VTRV'
    'IQARIaChZCQUNLVVBfVFlQRV9BV1NfQkFDS1VQEAI=');

@$core.Deprecated('Use backupTypeFilterDescriptor instead')
const BackupTypeFilter$json = {
  '1': 'BackupTypeFilter',
  '2': [
    {'1': 'BACKUP_TYPE_FILTER_SYSTEM', '2': 0},
    {'1': 'BACKUP_TYPE_FILTER_USER', '2': 1},
    {'1': 'BACKUP_TYPE_FILTER_AWS_BACKUP', '2': 2},
    {'1': 'BACKUP_TYPE_FILTER_ALL', '2': 3},
  ],
};

/// Descriptor for `BackupTypeFilter`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List backupTypeFilterDescriptor = $convert.base64Decode(
    'ChBCYWNrdXBUeXBlRmlsdGVyEh0KGUJBQ0tVUF9UWVBFX0ZJTFRFUl9TWVNURU0QABIbChdCQU'
    'NLVVBfVFlQRV9GSUxURVJfVVNFUhABEiEKHUJBQ0tVUF9UWVBFX0ZJTFRFUl9BV1NfQkFDS1VQ'
    'EAISGgoWQkFDS1VQX1RZUEVfRklMVEVSX0FMTBAD');

@$core.Deprecated('Use batchStatementErrorCodeEnumDescriptor instead')
const BatchStatementErrorCodeEnum$json = {
  '1': 'BatchStatementErrorCodeEnum',
  '2': [
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_REQUESTLIMITEXCEEDED', '2': 0},
    {
      '1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_ITEMCOLLECTIONSIZELIMITEXCEEDED',
      '2': 1
    },
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_INTERNALSERVERERROR', '2': 2},
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_RESOURCENOTFOUND', '2': 3},
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_ACCESSDENIED', '2': 4},
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_DUPLICATEITEM', '2': 5},
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_VALIDATIONERROR', '2': 6},
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_CONDITIONALCHECKFAILED', '2': 7},
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_THROTTLINGERROR', '2': 8},
    {
      '1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_PROVISIONEDTHROUGHPUTEXCEEDED',
      '2': 9
    },
    {'1': 'BATCH_STATEMENT_ERROR_CODE_ENUM_TRANSACTIONCONFLICT', '2': 10},
  ],
};

/// Descriptor for `BatchStatementErrorCodeEnum`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List batchStatementErrorCodeEnumDescriptor = $convert.base64Decode(
    'ChtCYXRjaFN0YXRlbWVudEVycm9yQ29kZUVudW0SOAo0QkFUQ0hfU1RBVEVNRU5UX0VSUk9SX0'
    'NPREVfRU5VTV9SRVFVRVNUTElNSVRFWENFRURFRBAAEkMKP0JBVENIX1NUQVRFTUVOVF9FUlJP'
    'Ul9DT0RFX0VOVU1fSVRFTUNPTExFQ1RJT05TSVpFTElNSVRFWENFRURFRBABEjcKM0JBVENIX1'
    'NUQVRFTUVOVF9FUlJPUl9DT0RFX0VOVU1fSU5URVJOQUxTRVJWRVJFUlJPUhACEjQKMEJBVENI'
    'X1NUQVRFTUVOVF9FUlJPUl9DT0RFX0VOVU1fUkVTT1VSQ0VOT1RGT1VORBADEjAKLEJBVENIX1'
    'NUQVRFTUVOVF9FUlJPUl9DT0RFX0VOVU1fQUNDRVNTREVOSUVEEAQSMQotQkFUQ0hfU1RBVEVN'
    'RU5UX0VSUk9SX0NPREVfRU5VTV9EVVBMSUNBVEVJVEVNEAUSMwovQkFUQ0hfU1RBVEVNRU5UX0'
    'VSUk9SX0NPREVfRU5VTV9WQUxJREFUSU9ORVJST1IQBhI6CjZCQVRDSF9TVEFURU1FTlRfRVJS'
    'T1JfQ09ERV9FTlVNX0NPTkRJVElPTkFMQ0hFQ0tGQUlMRUQQBxIzCi9CQVRDSF9TVEFURU1FTl'
    'RfRVJST1JfQ09ERV9FTlVNX1RIUk9UVExJTkdFUlJPUhAIEkEKPUJBVENIX1NUQVRFTUVOVF9F'
    'UlJPUl9DT0RFX0VOVU1fUFJPVklTSU9ORURUSFJPVUdIUFVURVhDRUVERUQQCRI3CjNCQVRDSF'
    '9TVEFURU1FTlRfRVJST1JfQ09ERV9FTlVNX1RSQU5TQUNUSU9OQ09ORkxJQ1QQCg==');

@$core.Deprecated('Use billingModeDescriptor instead')
const BillingMode$json = {
  '1': 'BillingMode',
  '2': [
    {'1': 'BILLING_MODE_PAY_PER_REQUEST', '2': 0},
    {'1': 'BILLING_MODE_PROVISIONED', '2': 1},
  ],
};

/// Descriptor for `BillingMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List billingModeDescriptor = $convert.base64Decode(
    'CgtCaWxsaW5nTW9kZRIgChxCSUxMSU5HX01PREVfUEFZX1BFUl9SRVFVRVNUEAASHAoYQklMTE'
    'lOR19NT0RFX1BST1ZJU0lPTkVEEAE=');

@$core.Deprecated('Use comparisonOperatorDescriptor instead')
const ComparisonOperator$json = {
  '1': 'ComparisonOperator',
  '2': [
    {'1': 'COMPARISON_OPERATOR_GE', '2': 0},
    {'1': 'COMPARISON_OPERATOR_NOT_NULL', '2': 1},
    {'1': 'COMPARISON_OPERATOR_EQ', '2': 2},
    {'1': 'COMPARISON_OPERATOR_IN', '2': 3},
    {'1': 'COMPARISON_OPERATOR_LT', '2': 4},
    {'1': 'COMPARISON_OPERATOR_CONTAINS', '2': 5},
    {'1': 'COMPARISON_OPERATOR_GT', '2': 6},
    {'1': 'COMPARISON_OPERATOR_NOT_CONTAINS', '2': 7},
    {'1': 'COMPARISON_OPERATOR_BETWEEN', '2': 8},
    {'1': 'COMPARISON_OPERATOR_NULL', '2': 9},
    {'1': 'COMPARISON_OPERATOR_BEGINS_WITH', '2': 10},
    {'1': 'COMPARISON_OPERATOR_LE', '2': 11},
    {'1': 'COMPARISON_OPERATOR_NE', '2': 12},
  ],
};

/// Descriptor for `ComparisonOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List comparisonOperatorDescriptor = $convert.base64Decode(
    'ChJDb21wYXJpc29uT3BlcmF0b3ISGgoWQ09NUEFSSVNPTl9PUEVSQVRPUl9HRRAAEiAKHENPTV'
    'BBUklTT05fT1BFUkFUT1JfTk9UX05VTEwQARIaChZDT01QQVJJU09OX09QRVJBVE9SX0VREAIS'
    'GgoWQ09NUEFSSVNPTl9PUEVSQVRPUl9JThADEhoKFkNPTVBBUklTT05fT1BFUkFUT1JfTFQQBB'
    'IgChxDT01QQVJJU09OX09QRVJBVE9SX0NPTlRBSU5TEAUSGgoWQ09NUEFSSVNPTl9PUEVSQVRP'
    'Ul9HVBAGEiQKIENPTVBBUklTT05fT1BFUkFUT1JfTk9UX0NPTlRBSU5TEAcSHwobQ09NUEFSSV'
    'NPTl9PUEVSQVRPUl9CRVRXRUVOEAgSHAoYQ09NUEFSSVNPTl9PUEVSQVRPUl9OVUxMEAkSIwof'
    'Q09NUEFSSVNPTl9PUEVSQVRPUl9CRUdJTlNfV0lUSBAKEhoKFkNPTVBBUklTT05fT1BFUkFUT1'
    'JfTEUQCxIaChZDT01QQVJJU09OX09QRVJBVE9SX05FEAw=');

@$core.Deprecated('Use conditionalOperatorDescriptor instead')
const ConditionalOperator$json = {
  '1': 'ConditionalOperator',
  '2': [
    {'1': 'CONDITIONAL_OPERATOR_AND', '2': 0},
    {'1': 'CONDITIONAL_OPERATOR_OR', '2': 1},
  ],
};

/// Descriptor for `ConditionalOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List conditionalOperatorDescriptor = $convert.base64Decode(
    'ChNDb25kaXRpb25hbE9wZXJhdG9yEhwKGENPTkRJVElPTkFMX09QRVJBVE9SX0FORBAAEhsKF0'
    'NPTkRJVElPTkFMX09QRVJBVE9SX09SEAE=');

@$core.Deprecated('Use continuousBackupsStatusDescriptor instead')
const ContinuousBackupsStatus$json = {
  '1': 'ContinuousBackupsStatus',
  '2': [
    {'1': 'CONTINUOUS_BACKUPS_STATUS_DISABLED', '2': 0},
    {'1': 'CONTINUOUS_BACKUPS_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ContinuousBackupsStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List continuousBackupsStatusDescriptor =
    $convert.base64Decode(
        'ChdDb250aW51b3VzQmFja3Vwc1N0YXR1cxImCiJDT05USU5VT1VTX0JBQ0tVUFNfU1RBVFVTX0'
        'RJU0FCTEVEEAASJQohQ09OVElOVU9VU19CQUNLVVBTX1NUQVRVU19FTkFCTEVEEAE=');

@$core.Deprecated('Use contributorInsightsActionDescriptor instead')
const ContributorInsightsAction$json = {
  '1': 'ContributorInsightsAction',
  '2': [
    {'1': 'CONTRIBUTOR_INSIGHTS_ACTION_ENABLE', '2': 0},
    {'1': 'CONTRIBUTOR_INSIGHTS_ACTION_DISABLE', '2': 1},
  ],
};

/// Descriptor for `ContributorInsightsAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contributorInsightsActionDescriptor =
    $convert.base64Decode(
        'ChlDb250cmlidXRvckluc2lnaHRzQWN0aW9uEiYKIkNPTlRSSUJVVE9SX0lOU0lHSFRTX0FDVE'
        'lPTl9FTkFCTEUQABInCiNDT05UUklCVVRPUl9JTlNJR0hUU19BQ1RJT05fRElTQUJMRRAB');

@$core.Deprecated('Use contributorInsightsModeDescriptor instead')
const ContributorInsightsMode$json = {
  '1': 'ContributorInsightsMode',
  '2': [
    {'1': 'CONTRIBUTOR_INSIGHTS_MODE_THROTTLED_KEYS', '2': 0},
    {'1': 'CONTRIBUTOR_INSIGHTS_MODE_ACCESSED_AND_THROTTLED_KEYS', '2': 1},
  ],
};

/// Descriptor for `ContributorInsightsMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contributorInsightsModeDescriptor = $convert.base64Decode(
    'ChdDb250cmlidXRvckluc2lnaHRzTW9kZRIsCihDT05UUklCVVRPUl9JTlNJR0hUU19NT0RFX1'
    'RIUk9UVExFRF9LRVlTEAASOQo1Q09OVFJJQlVUT1JfSU5TSUdIVFNfTU9ERV9BQ0NFU1NFRF9B'
    'TkRfVEhST1RUTEVEX0tFWVMQAQ==');

@$core.Deprecated('Use contributorInsightsStatusDescriptor instead')
const ContributorInsightsStatus$json = {
  '1': 'ContributorInsightsStatus',
  '2': [
    {'1': 'CONTRIBUTOR_INSIGHTS_STATUS_DISABLED', '2': 0},
    {'1': 'CONTRIBUTOR_INSIGHTS_STATUS_ENABLING', '2': 1},
    {'1': 'CONTRIBUTOR_INSIGHTS_STATUS_DISABLING', '2': 2},
    {'1': 'CONTRIBUTOR_INSIGHTS_STATUS_ENABLED', '2': 3},
    {'1': 'CONTRIBUTOR_INSIGHTS_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `ContributorInsightsStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contributorInsightsStatusDescriptor = $convert.base64Decode(
    'ChlDb250cmlidXRvckluc2lnaHRzU3RhdHVzEigKJENPTlRSSUJVVE9SX0lOU0lHSFRTX1NUQV'
    'RVU19ESVNBQkxFRBAAEigKJENPTlRSSUJVVE9SX0lOU0lHSFRTX1NUQVRVU19FTkFCTElORxAB'
    'EikKJUNPTlRSSUJVVE9SX0lOU0lHSFRTX1NUQVRVU19ESVNBQkxJTkcQAhInCiNDT05UUklCVV'
    'RPUl9JTlNJR0hUU19TVEFUVVNfRU5BQkxFRBADEiYKIkNPTlRSSUJVVE9SX0lOU0lHSFRTX1NU'
    'QVRVU19GQUlMRUQQBA==');

@$core.Deprecated('Use destinationStatusDescriptor instead')
const DestinationStatus$json = {
  '1': 'DestinationStatus',
  '2': [
    {'1': 'DESTINATION_STATUS_UPDATING', '2': 0},
    {'1': 'DESTINATION_STATUS_DISABLED', '2': 1},
    {'1': 'DESTINATION_STATUS_ENABLING', '2': 2},
    {'1': 'DESTINATION_STATUS_ENABLE_FAILED', '2': 3},
    {'1': 'DESTINATION_STATUS_ACTIVE', '2': 4},
    {'1': 'DESTINATION_STATUS_DISABLING', '2': 5},
  ],
};

/// Descriptor for `DestinationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List destinationStatusDescriptor = $convert.base64Decode(
    'ChFEZXN0aW5hdGlvblN0YXR1cxIfChtERVNUSU5BVElPTl9TVEFUVVNfVVBEQVRJTkcQABIfCh'
    'tERVNUSU5BVElPTl9TVEFUVVNfRElTQUJMRUQQARIfChtERVNUSU5BVElPTl9TVEFUVVNfRU5B'
    'QkxJTkcQAhIkCiBERVNUSU5BVElPTl9TVEFUVVNfRU5BQkxFX0ZBSUxFRBADEh0KGURFU1RJTk'
    'FUSU9OX1NUQVRVU19BQ1RJVkUQBBIgChxERVNUSU5BVElPTl9TVEFUVVNfRElTQUJMSU5HEAU=');

@$core.Deprecated('Use exportFormatDescriptor instead')
const ExportFormat$json = {
  '1': 'ExportFormat',
  '2': [
    {'1': 'EXPORT_FORMAT_DYNAMODB_JSON', '2': 0},
    {'1': 'EXPORT_FORMAT_ION', '2': 1},
  ],
};

/// Descriptor for `ExportFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List exportFormatDescriptor = $convert.base64Decode(
    'CgxFeHBvcnRGb3JtYXQSHwobRVhQT1JUX0ZPUk1BVF9EWU5BTU9EQl9KU09OEAASFQoRRVhQT1'
    'JUX0ZPUk1BVF9JT04QAQ==');

@$core.Deprecated('Use exportStatusDescriptor instead')
const ExportStatus$json = {
  '1': 'ExportStatus',
  '2': [
    {'1': 'EXPORT_STATUS_IN_PROGRESS', '2': 0},
    {'1': 'EXPORT_STATUS_COMPLETED', '2': 1},
    {'1': 'EXPORT_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `ExportStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List exportStatusDescriptor = $convert.base64Decode(
    'CgxFeHBvcnRTdGF0dXMSHQoZRVhQT1JUX1NUQVRVU19JTl9QUk9HUkVTUxAAEhsKF0VYUE9SVF'
    '9TVEFUVVNfQ09NUExFVEVEEAESGAoURVhQT1JUX1NUQVRVU19GQUlMRUQQAg==');

@$core.Deprecated('Use exportTypeDescriptor instead')
const ExportType$json = {
  '1': 'ExportType',
  '2': [
    {'1': 'EXPORT_TYPE_FULL_EXPORT', '2': 0},
    {'1': 'EXPORT_TYPE_INCREMENTAL_EXPORT', '2': 1},
  ],
};

/// Descriptor for `ExportType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List exportTypeDescriptor = $convert.base64Decode(
    'CgpFeHBvcnRUeXBlEhsKF0VYUE9SVF9UWVBFX0ZVTExfRVhQT1JUEAASIgoeRVhQT1JUX1RZUE'
    'VfSU5DUkVNRU5UQUxfRVhQT1JUEAE=');

@$core.Deprecated('Use exportViewTypeDescriptor instead')
const ExportViewType$json = {
  '1': 'ExportViewType',
  '2': [
    {'1': 'EXPORT_VIEW_TYPE_NEW_IMAGE', '2': 0},
    {'1': 'EXPORT_VIEW_TYPE_NEW_AND_OLD_IMAGES', '2': 1},
  ],
};

/// Descriptor for `ExportViewType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List exportViewTypeDescriptor = $convert.base64Decode(
    'Cg5FeHBvcnRWaWV3VHlwZRIeChpFWFBPUlRfVklFV19UWVBFX05FV19JTUFHRRAAEicKI0VYUE'
    '9SVF9WSUVXX1RZUEVfTkVXX0FORF9PTERfSU1BR0VTEAE=');

@$core.Deprecated('Use globalTableSettingsReplicationModeDescriptor instead')
const GlobalTableSettingsReplicationMode$json = {
  '1': 'GlobalTableSettingsReplicationMode',
  '2': [
    {'1': 'GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_DISABLED', '2': 0},
    {
      '1': 'GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED_WITH_OVERRIDES',
      '2': 1
    },
    {'1': 'GLOBAL_TABLE_SETTINGS_REPLICATION_MODE_ENABLED', '2': 2},
  ],
};

/// Descriptor for `GlobalTableSettingsReplicationMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List globalTableSettingsReplicationModeDescriptor =
    $convert.base64Decode(
        'CiJHbG9iYWxUYWJsZVNldHRpbmdzUmVwbGljYXRpb25Nb2RlEjMKL0dMT0JBTF9UQUJMRV9TRV'
        'RUSU5HU19SRVBMSUNBVElPTl9NT0RFX0RJU0FCTEVEEAASQQo9R0xPQkFMX1RBQkxFX1NFVFRJ'
        'TkdTX1JFUExJQ0FUSU9OX01PREVfRU5BQkxFRF9XSVRIX09WRVJSSURFUxABEjIKLkdMT0JBTF'
        '9UQUJMRV9TRVRUSU5HU19SRVBMSUNBVElPTl9NT0RFX0VOQUJMRUQQAg==');

@$core.Deprecated('Use globalTableStatusDescriptor instead')
const GlobalTableStatus$json = {
  '1': 'GlobalTableStatus',
  '2': [
    {'1': 'GLOBAL_TABLE_STATUS_UPDATING', '2': 0},
    {'1': 'GLOBAL_TABLE_STATUS_ACTIVE', '2': 1},
    {'1': 'GLOBAL_TABLE_STATUS_DELETING', '2': 2},
    {'1': 'GLOBAL_TABLE_STATUS_CREATING', '2': 3},
  ],
};

/// Descriptor for `GlobalTableStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List globalTableStatusDescriptor = $convert.base64Decode(
    'ChFHbG9iYWxUYWJsZVN0YXR1cxIgChxHTE9CQUxfVEFCTEVfU1RBVFVTX1VQREFUSU5HEAASHg'
    'oaR0xPQkFMX1RBQkxFX1NUQVRVU19BQ1RJVkUQARIgChxHTE9CQUxfVEFCTEVfU1RBVFVTX0RF'
    'TEVUSU5HEAISIAocR0xPQkFMX1RBQkxFX1NUQVRVU19DUkVBVElORxAD');

@$core.Deprecated('Use importStatusDescriptor instead')
const ImportStatus$json = {
  '1': 'ImportStatus',
  '2': [
    {'1': 'IMPORT_STATUS_IN_PROGRESS', '2': 0},
    {'1': 'IMPORT_STATUS_CANCELLED', '2': 1},
    {'1': 'IMPORT_STATUS_CANCELLING', '2': 2},
    {'1': 'IMPORT_STATUS_COMPLETED', '2': 3},
    {'1': 'IMPORT_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `ImportStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importStatusDescriptor = $convert.base64Decode(
    'CgxJbXBvcnRTdGF0dXMSHQoZSU1QT1JUX1NUQVRVU19JTl9QUk9HUkVTUxAAEhsKF0lNUE9SVF'
    '9TVEFUVVNfQ0FOQ0VMTEVEEAESHAoYSU1QT1JUX1NUQVRVU19DQU5DRUxMSU5HEAISGwoXSU1Q'
    'T1JUX1NUQVRVU19DT01QTEVURUQQAxIYChRJTVBPUlRfU1RBVFVTX0ZBSUxFRBAE');

@$core.Deprecated('Use indexStatusDescriptor instead')
const IndexStatus$json = {
  '1': 'IndexStatus',
  '2': [
    {'1': 'INDEX_STATUS_UPDATING', '2': 0},
    {'1': 'INDEX_STATUS_ACTIVE', '2': 1},
    {'1': 'INDEX_STATUS_DELETING', '2': 2},
    {'1': 'INDEX_STATUS_CREATING', '2': 3},
  ],
};

/// Descriptor for `IndexStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List indexStatusDescriptor = $convert.base64Decode(
    'CgtJbmRleFN0YXR1cxIZChVJTkRFWF9TVEFUVVNfVVBEQVRJTkcQABIXChNJTkRFWF9TVEFUVV'
    'NfQUNUSVZFEAESGQoVSU5ERVhfU1RBVFVTX0RFTEVUSU5HEAISGQoVSU5ERVhfU1RBVFVTX0NS'
    'RUFUSU5HEAM=');

@$core.Deprecated('Use inputCompressionTypeDescriptor instead')
const InputCompressionType$json = {
  '1': 'InputCompressionType',
  '2': [
    {'1': 'INPUT_COMPRESSION_TYPE_NONE', '2': 0},
    {'1': 'INPUT_COMPRESSION_TYPE_GZIP', '2': 1},
    {'1': 'INPUT_COMPRESSION_TYPE_ZSTD', '2': 2},
  ],
};

/// Descriptor for `InputCompressionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inputCompressionTypeDescriptor = $convert.base64Decode(
    'ChRJbnB1dENvbXByZXNzaW9uVHlwZRIfChtJTlBVVF9DT01QUkVTU0lPTl9UWVBFX05PTkUQAB'
    'IfChtJTlBVVF9DT01QUkVTU0lPTl9UWVBFX0daSVAQARIfChtJTlBVVF9DT01QUkVTU0lPTl9U'
    'WVBFX1pTVEQQAg==');

@$core.Deprecated('Use inputFormatDescriptor instead')
const InputFormat$json = {
  '1': 'InputFormat',
  '2': [
    {'1': 'INPUT_FORMAT_DYNAMODB_JSON', '2': 0},
    {'1': 'INPUT_FORMAT_CSV', '2': 1},
    {'1': 'INPUT_FORMAT_ION', '2': 2},
  ],
};

/// Descriptor for `InputFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inputFormatDescriptor = $convert.base64Decode(
    'CgtJbnB1dEZvcm1hdBIeChpJTlBVVF9GT1JNQVRfRFlOQU1PREJfSlNPThAAEhQKEElOUFVUX0'
    'ZPUk1BVF9DU1YQARIUChBJTlBVVF9GT1JNQVRfSU9OEAI=');

@$core.Deprecated('Use keyTypeDescriptor instead')
const KeyType$json = {
  '1': 'KeyType',
  '2': [
    {'1': 'KEY_TYPE_HASH', '2': 0},
    {'1': 'KEY_TYPE_RANGE', '2': 1},
  ],
};

/// Descriptor for `KeyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyTypeDescriptor = $convert.base64Decode(
    'CgdLZXlUeXBlEhEKDUtFWV9UWVBFX0hBU0gQABISCg5LRVlfVFlQRV9SQU5HRRAB');

@$core.Deprecated('Use multiRegionConsistencyDescriptor instead')
const MultiRegionConsistency$json = {
  '1': 'MultiRegionConsistency',
  '2': [
    {'1': 'MULTI_REGION_CONSISTENCY_EVENTUAL', '2': 0},
    {'1': 'MULTI_REGION_CONSISTENCY_STRONG', '2': 1},
  ],
};

/// Descriptor for `MultiRegionConsistency`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List multiRegionConsistencyDescriptor =
    $convert.base64Decode(
        'ChZNdWx0aVJlZ2lvbkNvbnNpc3RlbmN5EiUKIU1VTFRJX1JFR0lPTl9DT05TSVNURU5DWV9FVk'
        'VOVFVBTBAAEiMKH01VTFRJX1JFR0lPTl9DT05TSVNURU5DWV9TVFJPTkcQAQ==');

@$core.Deprecated('Use pointInTimeRecoveryStatusDescriptor instead')
const PointInTimeRecoveryStatus$json = {
  '1': 'PointInTimeRecoveryStatus',
  '2': [
    {'1': 'POINT_IN_TIME_RECOVERY_STATUS_DISABLED', '2': 0},
    {'1': 'POINT_IN_TIME_RECOVERY_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `PointInTimeRecoveryStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List pointInTimeRecoveryStatusDescriptor = $convert.base64Decode(
    'ChlQb2ludEluVGltZVJlY292ZXJ5U3RhdHVzEioKJlBPSU5UX0lOX1RJTUVfUkVDT1ZFUllfU1'
    'RBVFVTX0RJU0FCTEVEEAASKQolUE9JTlRfSU5fVElNRV9SRUNPVkVSWV9TVEFUVVNfRU5BQkxF'
    'RBAB');

@$core.Deprecated('Use projectionTypeDescriptor instead')
const ProjectionType$json = {
  '1': 'ProjectionType',
  '2': [
    {'1': 'PROJECTION_TYPE_KEYS_ONLY', '2': 0},
    {'1': 'PROJECTION_TYPE_INCLUDE', '2': 1},
    {'1': 'PROJECTION_TYPE_ALL', '2': 2},
  ],
};

/// Descriptor for `ProjectionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List projectionTypeDescriptor = $convert.base64Decode(
    'Cg5Qcm9qZWN0aW9uVHlwZRIdChlQUk9KRUNUSU9OX1RZUEVfS0VZU19PTkxZEAASGwoXUFJPSk'
    'VDVElPTl9UWVBFX0lOQ0xVREUQARIXChNQUk9KRUNUSU9OX1RZUEVfQUxMEAI=');

@$core.Deprecated('Use replicaStatusDescriptor instead')
const ReplicaStatus$json = {
  '1': 'ReplicaStatus',
  '2': [
    {'1': 'REPLICA_STATUS_UPDATING', '2': 0},
    {'1': 'REPLICA_STATUS_ARCHIVING', '2': 1},
    {'1': 'REPLICA_STATUS_REGION_DISABLED', '2': 2},
    {'1': 'REPLICA_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS', '2': 3},
    {'1': 'REPLICA_STATUS_ARCHIVED', '2': 4},
    {'1': 'REPLICA_STATUS_ACTIVE', '2': 5},
    {'1': 'REPLICA_STATUS_DELETING', '2': 6},
    {'1': 'REPLICA_STATUS_CREATION_FAILED', '2': 7},
    {'1': 'REPLICA_STATUS_CREATING', '2': 8},
    {'1': 'REPLICA_STATUS_REPLICATION_NOT_AUTHORIZED', '2': 9},
  ],
};

/// Descriptor for `ReplicaStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List replicaStatusDescriptor = $convert.base64Decode(
    'Cg1SZXBsaWNhU3RhdHVzEhsKF1JFUExJQ0FfU1RBVFVTX1VQREFUSU5HEAASHAoYUkVQTElDQV'
    '9TVEFUVVNfQVJDSElWSU5HEAESIgoeUkVQTElDQV9TVEFUVVNfUkVHSU9OX0RJU0FCTEVEEAIS'
    'NgoyUkVQTElDQV9TVEFUVVNfSU5BQ0NFU1NJQkxFX0VOQ1JZUFRJT05fQ1JFREVOVElBTFMQAx'
    'IbChdSRVBMSUNBX1NUQVRVU19BUkNISVZFRBAEEhkKFVJFUExJQ0FfU1RBVFVTX0FDVElWRRAF'
    'EhsKF1JFUExJQ0FfU1RBVFVTX0RFTEVUSU5HEAYSIgoeUkVQTElDQV9TVEFUVVNfQ1JFQVRJT0'
    '5fRkFJTEVEEAcSGwoXUkVQTElDQV9TVEFUVVNfQ1JFQVRJTkcQCBItCilSRVBMSUNBX1NUQVRV'
    'U19SRVBMSUNBVElPTl9OT1RfQVVUSE9SSVpFRBAJ');

@$core.Deprecated('Use returnConsumedCapacityDescriptor instead')
const ReturnConsumedCapacity$json = {
  '1': 'ReturnConsumedCapacity',
  '2': [
    {'1': 'RETURN_CONSUMED_CAPACITY_NONE', '2': 0},
    {'1': 'RETURN_CONSUMED_CAPACITY_INDEXES', '2': 1},
    {'1': 'RETURN_CONSUMED_CAPACITY_TOTAL', '2': 2},
  ],
};

/// Descriptor for `ReturnConsumedCapacity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List returnConsumedCapacityDescriptor = $convert.base64Decode(
    'ChZSZXR1cm5Db25zdW1lZENhcGFjaXR5EiEKHVJFVFVSTl9DT05TVU1FRF9DQVBBQ0lUWV9OT0'
    '5FEAASJAogUkVUVVJOX0NPTlNVTUVEX0NBUEFDSVRZX0lOREVYRVMQARIiCh5SRVRVUk5fQ09O'
    'U1VNRURfQ0FQQUNJVFlfVE9UQUwQAg==');

@$core.Deprecated('Use returnItemCollectionMetricsDescriptor instead')
const ReturnItemCollectionMetrics$json = {
  '1': 'ReturnItemCollectionMetrics',
  '2': [
    {'1': 'RETURN_ITEM_COLLECTION_METRICS_NONE', '2': 0},
    {'1': 'RETURN_ITEM_COLLECTION_METRICS_SIZE', '2': 1},
  ],
};

/// Descriptor for `ReturnItemCollectionMetrics`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List returnItemCollectionMetricsDescriptor =
    $convert.base64Decode(
        'ChtSZXR1cm5JdGVtQ29sbGVjdGlvbk1ldHJpY3MSJwojUkVUVVJOX0lURU1fQ09MTEVDVElPTl'
        '9NRVRSSUNTX05PTkUQABInCiNSRVRVUk5fSVRFTV9DT0xMRUNUSU9OX01FVFJJQ1NfU0laRRAB');

@$core.Deprecated('Use returnValueDescriptor instead')
const ReturnValue$json = {
  '1': 'ReturnValue',
  '2': [
    {'1': 'RETURN_VALUE_UPDATED_NEW', '2': 0},
    {'1': 'RETURN_VALUE_NONE', '2': 1},
    {'1': 'RETURN_VALUE_UPDATED_OLD', '2': 2},
    {'1': 'RETURN_VALUE_ALL_OLD', '2': 3},
    {'1': 'RETURN_VALUE_ALL_NEW', '2': 4},
  ],
};

/// Descriptor for `ReturnValue`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List returnValueDescriptor = $convert.base64Decode(
    'CgtSZXR1cm5WYWx1ZRIcChhSRVRVUk5fVkFMVUVfVVBEQVRFRF9ORVcQABIVChFSRVRVUk5fVk'
    'FMVUVfTk9ORRABEhwKGFJFVFVSTl9WQUxVRV9VUERBVEVEX09MRBACEhgKFFJFVFVSTl9WQUxV'
    'RV9BTExfT0xEEAMSGAoUUkVUVVJOX1ZBTFVFX0FMTF9ORVcQBA==');

@$core.Deprecated('Use returnValuesOnConditionCheckFailureDescriptor instead')
const ReturnValuesOnConditionCheckFailure$json = {
  '1': 'ReturnValuesOnConditionCheckFailure',
  '2': [
    {'1': 'RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_NONE', '2': 0},
    {'1': 'RETURN_VALUES_ON_CONDITION_CHECK_FAILURE_ALL_OLD', '2': 1},
  ],
};

/// Descriptor for `ReturnValuesOnConditionCheckFailure`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List returnValuesOnConditionCheckFailureDescriptor =
    $convert.base64Decode(
        'CiNSZXR1cm5WYWx1ZXNPbkNvbmRpdGlvbkNoZWNrRmFpbHVyZRIxCi1SRVRVUk5fVkFMVUVTX0'
        '9OX0NPTkRJVElPTl9DSEVDS19GQUlMVVJFX05PTkUQABI0CjBSRVRVUk5fVkFMVUVTX09OX0NP'
        'TkRJVElPTl9DSEVDS19GQUlMVVJFX0FMTF9PTEQQAQ==');

@$core.Deprecated('Use s3SseAlgorithmDescriptor instead')
const S3SseAlgorithm$json = {
  '1': 'S3SseAlgorithm',
  '2': [
    {'1': 'S3_SSE_ALGORITHM_AES256', '2': 0},
    {'1': 'S3_SSE_ALGORITHM_KMS', '2': 1},
  ],
};

/// Descriptor for `S3SseAlgorithm`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List s3SseAlgorithmDescriptor = $convert.base64Decode(
    'Cg5TM1NzZUFsZ29yaXRobRIbChdTM19TU0VfQUxHT1JJVEhNX0FFUzI1NhAAEhgKFFMzX1NTRV'
    '9BTEdPUklUSE1fS01TEAE=');

@$core.Deprecated('Use sSEStatusDescriptor instead')
const SSEStatus$json = {
  '1': 'SSEStatus',
  '2': [
    {'1': 'S_S_E_STATUS_UPDATING', '2': 0},
    {'1': 'S_S_E_STATUS_DISABLED', '2': 1},
    {'1': 'S_S_E_STATUS_ENABLING', '2': 2},
    {'1': 'S_S_E_STATUS_DISABLING', '2': 3},
    {'1': 'S_S_E_STATUS_ENABLED', '2': 4},
  ],
};

/// Descriptor for `SSEStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sSEStatusDescriptor = $convert.base64Decode(
    'CglTU0VTdGF0dXMSGQoVU19TX0VfU1RBVFVTX1VQREFUSU5HEAASGQoVU19TX0VfU1RBVFVTX0'
    'RJU0FCTEVEEAESGQoVU19TX0VfU1RBVFVTX0VOQUJMSU5HEAISGgoWU19TX0VfU1RBVFVTX0RJ'
    'U0FCTElORxADEhgKFFNfU19FX1NUQVRVU19FTkFCTEVEEAQ=');

@$core.Deprecated('Use sSETypeDescriptor instead')
const SSEType$json = {
  '1': 'SSEType',
  '2': [
    {'1': 'S_S_E_TYPE_AES256', '2': 0},
    {'1': 'S_S_E_TYPE_KMS', '2': 1},
  ],
};

/// Descriptor for `SSEType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sSETypeDescriptor = $convert.base64Decode(
    'CgdTU0VUeXBlEhUKEVNfU19FX1RZUEVfQUVTMjU2EAASEgoOU19TX0VfVFlQRV9LTVMQAQ==');

@$core.Deprecated('Use scalarAttributeTypeDescriptor instead')
const ScalarAttributeType$json = {
  '1': 'ScalarAttributeType',
  '2': [
    {'1': 'SCALAR_ATTRIBUTE_TYPE_B', '2': 0},
    {'1': 'SCALAR_ATTRIBUTE_TYPE_S', '2': 1},
    {'1': 'SCALAR_ATTRIBUTE_TYPE_N', '2': 2},
  ],
};

/// Descriptor for `ScalarAttributeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scalarAttributeTypeDescriptor = $convert.base64Decode(
    'ChNTY2FsYXJBdHRyaWJ1dGVUeXBlEhsKF1NDQUxBUl9BVFRSSUJVVEVfVFlQRV9CEAASGwoXU0'
    'NBTEFSX0FUVFJJQlVURV9UWVBFX1MQARIbChdTQ0FMQVJfQVRUUklCVVRFX1RZUEVfThAC');

@$core.Deprecated('Use selectDescriptor instead')
const Select$json = {
  '1': 'Select',
  '2': [
    {'1': 'SELECT_ALL_ATTRIBUTES', '2': 0},
    {'1': 'SELECT_COUNT', '2': 1},
    {'1': 'SELECT_ALL_PROJECTED_ATTRIBUTES', '2': 2},
    {'1': 'SELECT_SPECIFIC_ATTRIBUTES', '2': 3},
  ],
};

/// Descriptor for `Select`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List selectDescriptor = $convert.base64Decode(
    'CgZTZWxlY3QSGQoVU0VMRUNUX0FMTF9BVFRSSUJVVEVTEAASEAoMU0VMRUNUX0NPVU5UEAESIw'
    'ofU0VMRUNUX0FMTF9QUk9KRUNURURfQVRUUklCVVRFUxACEh4KGlNFTEVDVF9TUEVDSUZJQ19B'
    'VFRSSUJVVEVTEAM=');

@$core.Deprecated('Use streamViewTypeDescriptor instead')
const StreamViewType$json = {
  '1': 'StreamViewType',
  '2': [
    {'1': 'STREAM_VIEW_TYPE_KEYS_ONLY', '2': 0},
    {'1': 'STREAM_VIEW_TYPE_NEW_IMAGE', '2': 1},
    {'1': 'STREAM_VIEW_TYPE_NEW_AND_OLD_IMAGES', '2': 2},
    {'1': 'STREAM_VIEW_TYPE_OLD_IMAGE', '2': 3},
  ],
};

/// Descriptor for `StreamViewType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List streamViewTypeDescriptor = $convert.base64Decode(
    'Cg5TdHJlYW1WaWV3VHlwZRIeChpTVFJFQU1fVklFV19UWVBFX0tFWVNfT05MWRAAEh4KGlNUUk'
    'VBTV9WSUVXX1RZUEVfTkVXX0lNQUdFEAESJwojU1RSRUFNX1ZJRVdfVFlQRV9ORVdfQU5EX09M'
    'RF9JTUFHRVMQAhIeChpTVFJFQU1fVklFV19UWVBFX09MRF9JTUFHRRAD');

@$core.Deprecated('Use tableClassDescriptor instead')
const TableClass$json = {
  '1': 'TableClass',
  '2': [
    {'1': 'TABLE_CLASS_STANDARD_INFREQUENT_ACCESS', '2': 0},
    {'1': 'TABLE_CLASS_STANDARD', '2': 1},
  ],
};

/// Descriptor for `TableClass`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tableClassDescriptor = $convert.base64Decode(
    'CgpUYWJsZUNsYXNzEioKJlRBQkxFX0NMQVNTX1NUQU5EQVJEX0lORlJFUVVFTlRfQUNDRVNTEA'
    'ASGAoUVEFCTEVfQ0xBU1NfU1RBTkRBUkQQAQ==');

@$core.Deprecated('Use tableStatusDescriptor instead')
const TableStatus$json = {
  '1': 'TableStatus',
  '2': [
    {'1': 'TABLE_STATUS_UPDATING', '2': 0},
    {'1': 'TABLE_STATUS_ARCHIVING', '2': 1},
    {'1': 'TABLE_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS', '2': 2},
    {'1': 'TABLE_STATUS_ARCHIVED', '2': 3},
    {'1': 'TABLE_STATUS_ACTIVE', '2': 4},
    {'1': 'TABLE_STATUS_DELETING', '2': 5},
    {'1': 'TABLE_STATUS_CREATING', '2': 6},
    {'1': 'TABLE_STATUS_REPLICATION_NOT_AUTHORIZED', '2': 7},
  ],
};

/// Descriptor for `TableStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tableStatusDescriptor = $convert.base64Decode(
    'CgtUYWJsZVN0YXR1cxIZChVUQUJMRV9TVEFUVVNfVVBEQVRJTkcQABIaChZUQUJMRV9TVEFUVV'
    'NfQVJDSElWSU5HEAESNAowVEFCTEVfU1RBVFVTX0lOQUNDRVNTSUJMRV9FTkNSWVBUSU9OX0NS'
    'RURFTlRJQUxTEAISGQoVVEFCTEVfU1RBVFVTX0FSQ0hJVkVEEAMSFwoTVEFCTEVfU1RBVFVTX0'
    'FDVElWRRAEEhkKFVRBQkxFX1NUQVRVU19ERUxFVElORxAFEhkKFVRBQkxFX1NUQVRVU19DUkVB'
    'VElORxAGEisKJ1RBQkxFX1NUQVRVU19SRVBMSUNBVElPTl9OT1RfQVVUSE9SSVpFRBAH');

@$core.Deprecated('Use timeToLiveStatusDescriptor instead')
const TimeToLiveStatus$json = {
  '1': 'TimeToLiveStatus',
  '2': [
    {'1': 'TIME_TO_LIVE_STATUS_DISABLED', '2': 0},
    {'1': 'TIME_TO_LIVE_STATUS_ENABLING', '2': 1},
    {'1': 'TIME_TO_LIVE_STATUS_DISABLING', '2': 2},
    {'1': 'TIME_TO_LIVE_STATUS_ENABLED', '2': 3},
  ],
};

/// Descriptor for `TimeToLiveStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List timeToLiveStatusDescriptor = $convert.base64Decode(
    'ChBUaW1lVG9MaXZlU3RhdHVzEiAKHFRJTUVfVE9fTElWRV9TVEFUVVNfRElTQUJMRUQQABIgCh'
    'xUSU1FX1RPX0xJVkVfU1RBVFVTX0VOQUJMSU5HEAESIQodVElNRV9UT19MSVZFX1NUQVRVU19E'
    'SVNBQkxJTkcQAhIfChtUSU1FX1RPX0xJVkVfU1RBVFVTX0VOQUJMRUQQAw==');

@$core.Deprecated('Use witnessStatusDescriptor instead')
const WitnessStatus$json = {
  '1': 'WitnessStatus',
  '2': [
    {'1': 'WITNESS_STATUS_ACTIVE', '2': 0},
    {'1': 'WITNESS_STATUS_DELETING', '2': 1},
    {'1': 'WITNESS_STATUS_CREATING', '2': 2},
  ],
};

/// Descriptor for `WitnessStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List witnessStatusDescriptor = $convert.base64Decode(
    'Cg1XaXRuZXNzU3RhdHVzEhkKFVdJVE5FU1NfU1RBVFVTX0FDVElWRRAAEhsKF1dJVE5FU1NfU1'
    'RBVFVTX0RFTEVUSU5HEAESGwoXV0lUTkVTU19TVEFUVVNfQ1JFQVRJTkcQAg==');

@$core.Deprecated('Use archivalSummaryDescriptor instead')
const ArchivalSummary$json = {
  '1': 'ArchivalSummary',
  '2': [
    {
      '1': 'archivalbackuparn',
      '3': 95197293,
      '4': 1,
      '5': 9,
      '10': 'archivalbackuparn'
    },
    {
      '1': 'archivaldatetime',
      '3': 424460035,
      '4': 1,
      '5': 9,
      '10': 'archivaldatetime'
    },
    {
      '1': 'archivalreason',
      '3': 132420188,
      '4': 1,
      '5': 9,
      '10': 'archivalreason'
    },
  ],
};

/// Descriptor for `ArchivalSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List archivalSummaryDescriptor = $convert.base64Decode(
    'Cg9BcmNoaXZhbFN1bW1hcnkSLwoRYXJjaGl2YWxiYWNrdXBhcm4Y7bCyLSABKAlSEWFyY2hpdm'
    'FsYmFja3VwYXJuEi4KEGFyY2hpdmFsZGF0ZXRpbWUYg/6yygEgASgJUhBhcmNoaXZhbGRhdGV0'
    'aW1lEikKDmFyY2hpdmFscmVhc29uGNykkj8gASgJUg5hcmNoaXZhbHJlYXNvbg==');

@$core.Deprecated('Use attributeDefinitionDescriptor instead')
const AttributeDefinition$json = {
  '1': 'AttributeDefinition',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'attributetype',
      '3': 54612120,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ScalarAttributeType',
      '10': 'attributetype'
    },
  ],
};

/// Descriptor for `AttributeDefinition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attributeDefinitionDescriptor = $convert.base64Decode(
    'ChNBdHRyaWJ1dGVEZWZpbml0aW9uEigKDWF0dHJpYnV0ZW5hbWUYrZWYqAEgASgJUg1hdHRyaW'
    'J1dGVuYW1lEkYKDWF0dHJpYnV0ZXR5cGUYmKGFGiABKA4yHS5keW5hbW9kYi5TY2FsYXJBdHRy'
    'aWJ1dGVUeXBlUg1hdHRyaWJ1dGV0eXBl');

@$core.Deprecated('Use attributeValueDescriptor instead')
const AttributeValue$json = {
  '1': 'AttributeValue',
  '2': [
    {'1': 'b', '3': 118225798, '4': 1, '5': 12, '10': 'b'},
    {'1': 'bool', '3': 307144766, '4': 1, '5': 8, '10': 'bool'},
    {'1': 'bs', '3': 232616419, '4': 3, '5': 12, '10': 'bs'},
    {
      '1': 'l',
      '3': 151781036,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'l'
    },
    {
      '1': 'm',
      '3': 135003417,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeValue.MEntry',
      '10': 'm'
    },
    {'1': 'n', '3': 185336274, '4': 1, '5': 9, '10': 'n'},
    {'1': 'ns', '3': 98983847, '4': 3, '5': 9, '10': 'ns'},
    {'1': 'null', '3': 426761765, '4': 1, '5': 8, '10': 'null'},
    {'1': 's', '3': 369890083, '4': 1, '5': 9, '10': 's'},
    {'1': 'ss', '3': 100834020, '4': 3, '5': 9, '10': 'ss'},
  ],
  '3': [AttributeValue_MEntry$json],
};

@$core.Deprecated('Use attributeValueDescriptor instead')
const AttributeValue_MEntry$json = {
  '1': 'MEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `AttributeValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attributeValueDescriptor = $convert.base64Decode(
    'Cg5BdHRyaWJ1dGVWYWx1ZRIPCgFiGIb3rzggASgMUgFiEhYKBGJvb2wYvtC6kgEgASgIUgRib2'
    '9sEhEKAmJzGOPj9W4gAygMUgJicxIpCgFsGKz9r0ggAygLMhguZHluYW1vZGIuQXR0cmlidXRl'
    'VmFsdWVSAWwSMAoBbRiZ+q9AIAMoCzIfLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlLk1FbnRyeV'
    'IBbRIPCgFuGNKDsFggASgJUgFuEhEKAm5zGKe/mS8gAygJUgJucxIWCgRudWxsGKW8v8sBIAEo'
    'CFIEbnVsbBIQCgFzGKOmsLABIAEoCVIBcxIRCgJzcxjktYowIAMoCVICc3MaTgoGTUVudHJ5Eh'
    'AKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFs'
    'dWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use attributeValueUpdateDescriptor instead')
const AttributeValueUpdate$json = {
  '1': 'AttributeValueUpdate',
  '2': [
    {
      '1': 'action',
      '3': 175614240,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.AttributeAction',
      '10': 'action'
    },
    {
      '1': 'value',
      '3': 289929579,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
};

/// Descriptor for `AttributeValueUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attributeValueUpdateDescriptor = $convert.base64Decode(
    'ChRBdHRyaWJ1dGVWYWx1ZVVwZGF0ZRI0CgZhY3Rpb24YoNLeUyABKA4yGS5keW5hbW9kYi5BdH'
    'RyaWJ1dGVBY3Rpb25SBmFjdGlvbhIyCgV2YWx1ZRjr8p+KASABKAsyGC5keW5hbW9kYi5BdHRy'
    'aWJ1dGVWYWx1ZVIFdmFsdWU=');

@$core.Deprecated('Use autoScalingPolicyDescriptionDescriptor instead')
const AutoScalingPolicyDescription$json = {
  '1': 'AutoScalingPolicyDescription',
  '2': [
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'targettrackingscalingpolicyconfiguration',
      '3': 474521147,
      '4': 1,
      '5': 11,
      '6':
          '.dynamodb.AutoScalingTargetTrackingScalingPolicyConfigurationDescription',
      '10': 'targettrackingscalingpolicyconfiguration'
    },
  ],
};

/// Descriptor for `AutoScalingPolicyDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoScalingPolicyDescriptionDescriptor = $convert.base64Decode(
    'ChxBdXRvU2NhbGluZ1BvbGljeURlc2NyaXB0aW9uEiEKCnBvbGljeW5hbWUYvfWHfyABKAlSCn'
    'BvbGljeW5hbWUSqAEKKHRhcmdldHRyYWNraW5nc2NhbGluZ3BvbGljeWNvbmZpZ3VyYXRpb24Y'
    'u7yi4gEgASgLMkguZHluYW1vZGIuQXV0b1NjYWxpbmdUYXJnZXRUcmFja2luZ1NjYWxpbmdQb2'
    'xpY3lDb25maWd1cmF0aW9uRGVzY3JpcHRpb25SKHRhcmdldHRyYWNraW5nc2NhbGluZ3BvbGlj'
    'eWNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use autoScalingPolicyUpdateDescriptor instead')
const AutoScalingPolicyUpdate$json = {
  '1': 'AutoScalingPolicyUpdate',
  '2': [
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'targettrackingscalingpolicyconfiguration',
      '3': 474521147,
      '4': 1,
      '5': 11,
      '6':
          '.dynamodb.AutoScalingTargetTrackingScalingPolicyConfigurationUpdate',
      '10': 'targettrackingscalingpolicyconfiguration'
    },
  ],
};

/// Descriptor for `AutoScalingPolicyUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoScalingPolicyUpdateDescriptor = $convert.base64Decode(
    'ChdBdXRvU2NhbGluZ1BvbGljeVVwZGF0ZRIhCgpwb2xpY3luYW1lGL31h38gASgJUgpwb2xpY3'
    'luYW1lEqMBCih0YXJnZXR0cmFja2luZ3NjYWxpbmdwb2xpY3ljb25maWd1cmF0aW9uGLu8ouIB'
    'IAEoCzJDLmR5bmFtb2RiLkF1dG9TY2FsaW5nVGFyZ2V0VHJhY2tpbmdTY2FsaW5nUG9saWN5Q2'
    '9uZmlndXJhdGlvblVwZGF0ZVIodGFyZ2V0dHJhY2tpbmdzY2FsaW5ncG9saWN5Y29uZmlndXJh'
    'dGlvbg==');

@$core.Deprecated('Use autoScalingSettingsDescriptionDescriptor instead')
const AutoScalingSettingsDescription$json = {
  '1': 'AutoScalingSettingsDescription',
  '2': [
    {
      '1': 'autoscalingdisabled',
      '3': 183835736,
      '4': 1,
      '5': 8,
      '10': 'autoscalingdisabled'
    },
    {
      '1': 'autoscalingrolearn',
      '3': 348085703,
      '4': 1,
      '5': 9,
      '10': 'autoscalingrolearn'
    },
    {'1': 'maximumunits', '3': 501867481, '4': 1, '5': 3, '10': 'maximumunits'},
    {'1': 'minimumunits', '3': 190148659, '4': 1, '5': 3, '10': 'minimumunits'},
    {
      '1': 'scalingpolicies',
      '3': 289494257,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AutoScalingPolicyDescription',
      '10': 'scalingpolicies'
    },
  ],
};

/// Descriptor for `AutoScalingSettingsDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoScalingSettingsDescriptionDescriptor = $convert.base64Decode(
    'Ch5BdXRvU2NhbGluZ1NldHRpbmdzRGVzY3JpcHRpb24SMwoTYXV0b3NjYWxpbmdkaXNhYmxlZB'
    'jYuNRXIAEoCFITYXV0b3NjYWxpbmdkaXNhYmxlZBIyChJhdXRvc2NhbGluZ3JvbGVhcm4Yx7v9'
    'pQEgASgJUhJhdXRvc2NhbGluZ3JvbGVhcm4SJgoMbWF4aW11bXVuaXRzGNnHp+8BIAEoA1IMbW'
    'F4aW11bXVuaXRzEiUKDG1pbmltdW11bml0cxiz4NVaIAEoA1IMbWluaW11bXVuaXRzElQKD3Nj'
    'YWxpbmdwb2xpY2llcxjxqYWKASADKAsyJi5keW5hbW9kYi5BdXRvU2NhbGluZ1BvbGljeURlc2'
    'NyaXB0aW9uUg9zY2FsaW5ncG9saWNpZXM=');

@$core.Deprecated('Use autoScalingSettingsUpdateDescriptor instead')
const AutoScalingSettingsUpdate$json = {
  '1': 'AutoScalingSettingsUpdate',
  '2': [
    {
      '1': 'autoscalingdisabled',
      '3': 183835736,
      '4': 1,
      '5': 8,
      '10': 'autoscalingdisabled'
    },
    {
      '1': 'autoscalingrolearn',
      '3': 348085703,
      '4': 1,
      '5': 9,
      '10': 'autoscalingrolearn'
    },
    {'1': 'maximumunits', '3': 501867481, '4': 1, '5': 3, '10': 'maximumunits'},
    {'1': 'minimumunits', '3': 190148659, '4': 1, '5': 3, '10': 'minimumunits'},
    {
      '1': 'scalingpolicyupdate',
      '3': 189139846,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingPolicyUpdate',
      '10': 'scalingpolicyupdate'
    },
  ],
};

/// Descriptor for `AutoScalingSettingsUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List autoScalingSettingsUpdateDescriptor = $convert.base64Decode(
    'ChlBdXRvU2NhbGluZ1NldHRpbmdzVXBkYXRlEjMKE2F1dG9zY2FsaW5nZGlzYWJsZWQY2LjUVy'
    'ABKAhSE2F1dG9zY2FsaW5nZGlzYWJsZWQSMgoSYXV0b3NjYWxpbmdyb2xlYXJuGMe7/aUBIAEo'
    'CVISYXV0b3NjYWxpbmdyb2xlYXJuEiYKDG1heGltdW11bml0cxjZx6fvASABKANSDG1heGltdW'
    '11bml0cxIlCgxtaW5pbXVtdW5pdHMYs+DVWiABKANSDG1pbmltdW11bml0cxJWChNzY2FsaW5n'
    'cG9saWN5dXBkYXRlGIaXmFogASgLMiEuZHluYW1vZGIuQXV0b1NjYWxpbmdQb2xpY3lVcGRhdG'
    'VSE3NjYWxpbmdwb2xpY3l1cGRhdGU=');

@$core.Deprecated(
    'Use autoScalingTargetTrackingScalingPolicyConfigurationDescriptionDescriptor instead')
const AutoScalingTargetTrackingScalingPolicyConfigurationDescription$json = {
  '1': 'AutoScalingTargetTrackingScalingPolicyConfigurationDescription',
  '2': [
    {
      '1': 'disablescalein',
      '3': 464986377,
      '4': 1,
      '5': 8,
      '10': 'disablescalein'
    },
    {
      '1': 'scaleincooldown',
      '3': 45925876,
      '4': 1,
      '5': 5,
      '10': 'scaleincooldown'
    },
    {
      '1': 'scaleoutcooldown',
      '3': 316995887,
      '4': 1,
      '5': 5,
      '10': 'scaleoutcooldown'
    },
    {'1': 'targetvalue', '3': 118247738, '4': 1, '5': 1, '10': 'targetvalue'},
  ],
};

/// Descriptor for `AutoScalingTargetTrackingScalingPolicyConfigurationDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    autoScalingTargetTrackingScalingPolicyConfigurationDescriptionDescriptor =
    $convert.base64Decode(
        'Cj5BdXRvU2NhbGluZ1RhcmdldFRyYWNraW5nU2NhbGluZ1BvbGljeUNvbmZpZ3VyYXRpb25EZX'
        'NjcmlwdGlvbhIqCg5kaXNhYmxlc2NhbGVpbhiJwtzdASABKAhSDmRpc2FibGVzY2FsZWluEisK'
        'D3NjYWxlaW5jb29sZG93bhj0i/MVIAEoBVIPc2NhbGVpbmNvb2xkb3duEi4KEHNjYWxlb3V0Y2'
        '9vbGRvd24Yr/KTlwEgASgFUhBzY2FsZW91dGNvb2xkb3duEiMKC3RhcmdldHZhbHVlGLqisTgg'
        'ASgBUgt0YXJnZXR2YWx1ZQ==');

@$core.Deprecated(
    'Use autoScalingTargetTrackingScalingPolicyConfigurationUpdateDescriptor instead')
const AutoScalingTargetTrackingScalingPolicyConfigurationUpdate$json = {
  '1': 'AutoScalingTargetTrackingScalingPolicyConfigurationUpdate',
  '2': [
    {
      '1': 'disablescalein',
      '3': 464986377,
      '4': 1,
      '5': 8,
      '10': 'disablescalein'
    },
    {
      '1': 'scaleincooldown',
      '3': 45925876,
      '4': 1,
      '5': 5,
      '10': 'scaleincooldown'
    },
    {
      '1': 'scaleoutcooldown',
      '3': 316995887,
      '4': 1,
      '5': 5,
      '10': 'scaleoutcooldown'
    },
    {'1': 'targetvalue', '3': 118247738, '4': 1, '5': 1, '10': 'targetvalue'},
  ],
};

/// Descriptor for `AutoScalingTargetTrackingScalingPolicyConfigurationUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    autoScalingTargetTrackingScalingPolicyConfigurationUpdateDescriptor =
    $convert.base64Decode(
        'CjlBdXRvU2NhbGluZ1RhcmdldFRyYWNraW5nU2NhbGluZ1BvbGljeUNvbmZpZ3VyYXRpb25VcG'
        'RhdGUSKgoOZGlzYWJsZXNjYWxlaW4YicLc3QEgASgIUg5kaXNhYmxlc2NhbGVpbhIrCg9zY2Fs'
        'ZWluY29vbGRvd24Y9IvzFSABKAVSD3NjYWxlaW5jb29sZG93bhIuChBzY2FsZW91dGNvb2xkb3'
        'duGK/yk5cBIAEoBVIQc2NhbGVvdXRjb29sZG93bhIjCgt0YXJnZXR2YWx1ZRi6orE4IAEoAVIL'
        'dGFyZ2V0dmFsdWU=');

@$core.Deprecated('Use backupDescriptionDescriptor instead')
const BackupDescription$json = {
  '1': 'BackupDescription',
  '2': [
    {
      '1': 'backupdetails',
      '3': 383357432,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.BackupDetails',
      '10': 'backupdetails'
    },
    {
      '1': 'sourcetabledetails',
      '3': 507937979,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SourceTableDetails',
      '10': 'sourcetabledetails'
    },
    {
      '1': 'sourcetablefeaturedetails',
      '3': 21598361,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SourceTableFeatureDetails',
      '10': 'sourcetablefeaturedetails'
    },
  ],
};

/// Descriptor for `BackupDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List backupDescriptionDescriptor = $convert.base64Decode(
    'ChFCYWNrdXBEZXNjcmlwdGlvbhJBCg1iYWNrdXBkZXRhaWxzGPij5rYBIAEoCzIXLmR5bmFtb2'
    'RiLkJhY2t1cERldGFpbHNSDWJhY2t1cGRldGFpbHMSUAoSc291cmNldGFibGVkZXRhaWxzGLuJ'
    'mvIBIAEoCzIcLmR5bmFtb2RiLlNvdXJjZVRhYmxlRGV0YWlsc1ISc291cmNldGFibGVkZXRhaW'
    'xzEmQKGXNvdXJjZXRhYmxlZmVhdHVyZWRldGFpbHMYmaGmCiABKAsyIy5keW5hbW9kYi5Tb3Vy'
    'Y2VUYWJsZUZlYXR1cmVEZXRhaWxzUhlzb3VyY2V0YWJsZWZlYXR1cmVkZXRhaWxz');

@$core.Deprecated('Use backupDetailsDescriptor instead')
const BackupDetails$json = {
  '1': 'BackupDetails',
  '2': [
    {'1': 'backuparn', '3': 370874339, '4': 1, '5': 9, '10': 'backuparn'},
    {
      '1': 'backupcreationdatetime',
      '3': 368022476,
      '4': 1,
      '5': 9,
      '10': 'backupcreationdatetime'
    },
    {
      '1': 'backupexpirydatetime',
      '3': 471291762,
      '4': 1,
      '5': 9,
      '10': 'backupexpirydatetime'
    },
    {'1': 'backupname', '3': 467693789, '4': 1, '5': 9, '10': 'backupname'},
    {
      '1': 'backupsizebytes',
      '3': 147336318,
      '4': 1,
      '5': 3,
      '10': 'backupsizebytes'
    },
    {
      '1': 'backupstatus',
      '3': 382505546,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BackupStatus',
      '10': 'backupstatus'
    },
    {
      '1': 'backuptype',
      '3': 134973992,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BackupType',
      '10': 'backuptype'
    },
  ],
};

/// Descriptor for `BackupDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List backupDetailsDescriptor = $convert.base64Decode(
    'Cg1CYWNrdXBEZXRhaWxzEiAKCWJhY2t1cGFybhjjr+ywASABKAlSCWJhY2t1cGFybhI6ChZiYW'
    'NrdXBjcmVhdGlvbmRhdGV0aW1lGMynvq8BIAEoCVIWYmFja3VwY3JlYXRpb25kYXRldGltZRI2'
    'ChRiYWNrdXBleHBpcnlkYXRldGltZRjyrt3gASABKAlSFGJhY2t1cGV4cGlyeWRhdGV0aW1lEi'
    'IKCmJhY2t1cG5hbWUY3eGB3wEgASgJUgpiYWNrdXBuYW1lEisKD2JhY2t1cHNpemVieXRlcxj+'
    '2KBGIAEoA1IPYmFja3Vwc2l6ZWJ5dGVzEj4KDGJhY2t1cHN0YXR1cxjKpLK2ASABKA4yFi5keW'
    '5hbW9kYi5CYWNrdXBTdGF0dXNSDGJhY2t1cHN0YXR1cxI3CgpiYWNrdXB0eXBlGKiUrkAgASgO'
    'MhQuZHluYW1vZGIuQmFja3VwVHlwZVIKYmFja3VwdHlwZQ==');

@$core.Deprecated('Use backupInUseExceptionDescriptor instead')
const BackupInUseException$json = {
  '1': 'BackupInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BackupInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List backupInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChRCYWNrdXBJblVzZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use backupNotFoundExceptionDescriptor instead')
const BackupNotFoundException$json = {
  '1': 'BackupNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BackupNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List backupNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChdCYWNrdXBOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use backupSummaryDescriptor instead')
const BackupSummary$json = {
  '1': 'BackupSummary',
  '2': [
    {'1': 'backuparn', '3': 370874339, '4': 1, '5': 9, '10': 'backuparn'},
    {
      '1': 'backupcreationdatetime',
      '3': 368022476,
      '4': 1,
      '5': 9,
      '10': 'backupcreationdatetime'
    },
    {
      '1': 'backupexpirydatetime',
      '3': 471291762,
      '4': 1,
      '5': 9,
      '10': 'backupexpirydatetime'
    },
    {'1': 'backupname', '3': 467693789, '4': 1, '5': 9, '10': 'backupname'},
    {
      '1': 'backupsizebytes',
      '3': 147336318,
      '4': 1,
      '5': 3,
      '10': 'backupsizebytes'
    },
    {
      '1': 'backupstatus',
      '3': 382505546,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BackupStatus',
      '10': 'backupstatus'
    },
    {
      '1': 'backuptype',
      '3': 134973992,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BackupType',
      '10': 'backuptype'
    },
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {'1': 'tableid', '3': 449893011, '4': 1, '5': 9, '10': 'tableid'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `BackupSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List backupSummaryDescriptor = $convert.base64Decode(
    'Cg1CYWNrdXBTdW1tYXJ5EiAKCWJhY2t1cGFybhjjr+ywASABKAlSCWJhY2t1cGFybhI6ChZiYW'
    'NrdXBjcmVhdGlvbmRhdGV0aW1lGMynvq8BIAEoCVIWYmFja3VwY3JlYXRpb25kYXRldGltZRI2'
    'ChRiYWNrdXBleHBpcnlkYXRldGltZRjyrt3gASABKAlSFGJhY2t1cGV4cGlyeWRhdGV0aW1lEi'
    'IKCmJhY2t1cG5hbWUY3eGB3wEgASgJUgpiYWNrdXBuYW1lEisKD2JhY2t1cHNpemVieXRlcxj+'
    '2KBGIAEoA1IPYmFja3Vwc2l6ZWJ5dGVzEj4KDGJhY2t1cHN0YXR1cxjKpLK2ASABKA4yFi5keW'
    '5hbW9kYi5CYWNrdXBTdGF0dXNSDGJhY2t1cHN0YXR1cxI3CgpiYWNrdXB0eXBlGKiUrkAgASgO'
    'MhQuZHluYW1vZGIuQmFja3VwVHlwZVIKYmFja3VwdHlwZRIeCgh0YWJsZWFybhjjgOvNASABKA'
    'lSCHRhYmxlYXJuEhwKB3RhYmxlaWQYk6XD1gEgASgJUgd0YWJsZWlkEiAKCXRhYmxlbmFtZRjd'
    '5NqBASABKAlSCXRhYmxlbmFtZQ==');

@$core.Deprecated('Use batchExecuteStatementInputDescriptor instead')
const BatchExecuteStatementInput$json = {
  '1': 'BatchExecuteStatementInput',
  '2': [
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'statements',
      '3': 488352288,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchStatementRequest',
      '10': 'statements'
    },
  ],
};

/// Descriptor for `BatchExecuteStatementInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchExecuteStatementInputDescriptor = $convert.base64Decode(
    'ChpCYXRjaEV4ZWN1dGVTdGF0ZW1lbnRJbnB1dBJbChZyZXR1cm5jb25zdW1lZGNhcGFjaXR5GP'
    '7n4RQgASgOMiAuZHluYW1vZGIuUmV0dXJuQ29uc3VtZWRDYXBhY2l0eVIWcmV0dXJuY29uc3Vt'
    'ZWRjYXBhY2l0eRJDCgpzdGF0ZW1lbnRzGKDU7ugBIAMoCzIfLmR5bmFtb2RiLkJhdGNoU3RhdG'
    'VtZW50UmVxdWVzdFIKc3RhdGVtZW50cw==');

@$core.Deprecated('Use batchExecuteStatementOutputDescriptor instead')
const BatchExecuteStatementOutput$json = {
  '1': 'BatchExecuteStatementOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'responses',
      '3': 114852856,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchStatementResponse',
      '10': 'responses'
    },
  ],
};

/// Descriptor for `BatchExecuteStatementOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchExecuteStatementOutputDescriptor = $convert.base64Decode(
    'ChtCYXRjaEV4ZWN1dGVTdGF0ZW1lbnRPdXRwdXQSSgoQY29uc3VtZWRjYXBhY2l0eRisqqHWAS'
    'ADKAsyGi5keW5hbW9kYi5Db25zdW1lZENhcGFjaXR5UhBjb25zdW1lZGNhcGFjaXR5EkEKCXJl'
    'c3BvbnNlcxj4h+I2IAMoCzIgLmR5bmFtb2RiLkJhdGNoU3RhdGVtZW50UmVzcG9uc2VSCXJlc3'
    'BvbnNlcw==');

@$core.Deprecated('Use batchGetItemInputDescriptor instead')
const BatchGetItemInput$json = {
  '1': 'BatchGetItemInput',
  '2': [
    {
      '1': 'requestitems',
      '3': 247720687,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchGetItemInput.RequestitemsEntry',
      '10': 'requestitems'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
  ],
  '3': [BatchGetItemInput_RequestitemsEntry$json],
};

@$core.Deprecated('Use batchGetItemInputDescriptor instead')
const BatchGetItemInput_RequestitemsEntry$json = {
  '1': 'RequestitemsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.KeysAndAttributes',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `BatchGetItemInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetItemInputDescriptor = $convert.base64Decode(
    'ChFCYXRjaEdldEl0ZW1JbnB1dBJUCgxyZXF1ZXN0aXRlbXMY79WPdiADKAsyLS5keW5hbW9kYi'
    '5CYXRjaEdldEl0ZW1JbnB1dC5SZXF1ZXN0aXRlbXNFbnRyeVIMcmVxdWVzdGl0ZW1zElsKFnJl'
    'dHVybmNvbnN1bWVkY2FwYWNpdHkY/ufhFCABKA4yIC5keW5hbW9kYi5SZXR1cm5Db25zdW1lZE'
    'NhcGFjaXR5UhZyZXR1cm5jb25zdW1lZGNhcGFjaXR5GlwKEVJlcXVlc3RpdGVtc0VudHJ5EhAK'
    'A2tleRgBIAEoCVIDa2V5EjEKBXZhbHVlGAIgASgLMhsuZHluYW1vZGIuS2V5c0FuZEF0dHJpYn'
    'V0ZXNSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use batchGetItemOutputDescriptor instead')
const BatchGetItemOutput$json = {
  '1': 'BatchGetItemOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'responses',
      '3': 114852856,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchGetItemOutput.ResponsesEntry',
      '10': 'responses'
    },
    {
      '1': 'unprocessedkeys',
      '3': 123999275,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchGetItemOutput.UnprocessedkeysEntry',
      '10': 'unprocessedkeys'
    },
  ],
  '3': [
    BatchGetItemOutput_ResponsesEntry$json,
    BatchGetItemOutput_UnprocessedkeysEntry$json
  ],
};

@$core.Deprecated('Use batchGetItemOutputDescriptor instead')
const BatchGetItemOutput_ResponsesEntry$json = {
  '1': 'ResponsesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use batchGetItemOutputDescriptor instead')
const BatchGetItemOutput_UnprocessedkeysEntry$json = {
  '1': 'UnprocessedkeysEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.KeysAndAttributes',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `BatchGetItemOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetItemOutputDescriptor = $convert.base64Decode(
    'ChJCYXRjaEdldEl0ZW1PdXRwdXQSSgoQY29uc3VtZWRjYXBhY2l0eRisqqHWASADKAsyGi5keW'
    '5hbW9kYi5Db25zdW1lZENhcGFjaXR5UhBjb25zdW1lZGNhcGFjaXR5EkwKCXJlc3BvbnNlcxj4'
    'h+I2IAMoCzIrLmR5bmFtb2RiLkJhdGNoR2V0SXRlbU91dHB1dC5SZXNwb25zZXNFbnRyeVIJcm'
    'VzcG9uc2VzEl4KD3VucHJvY2Vzc2Vka2V5cxirqJA7IAMoCzIxLmR5bmFtb2RiLkJhdGNoR2V0'
    'SXRlbU91dHB1dC5VbnByb2Nlc3NlZGtleXNFbnRyeVIPdW5wcm9jZXNzZWRrZXlzGjwKDlJlc3'
    'BvbnNlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEa'
    'XwoUVW5wcm9jZXNzZWRrZXlzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSMQoFdmFsdWUYAiABKA'
    'syGy5keW5hbW9kYi5LZXlzQW5kQXR0cmlidXRlc1IFdmFsdWU6AjgB');

@$core.Deprecated('Use batchStatementErrorDescriptor instead')
const BatchStatementError$json = {
  '1': 'BatchStatementError',
  '2': [
    {
      '1': 'code',
      '3': 425572629,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BatchStatementErrorCodeEnum',
      '10': 'code'
    },
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchStatementError.ItemEntry',
      '10': 'item'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
  '3': [BatchStatementError_ItemEntry$json],
};

@$core.Deprecated('Use batchStatementErrorDescriptor instead')
const BatchStatementError_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `BatchStatementError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchStatementErrorDescriptor = $convert.base64Decode(
    'ChNCYXRjaFN0YXRlbWVudEVycm9yEj0KBGNvZGUYlfL2ygEgASgOMiUuZHluYW1vZGIuQmF0Y2'
    'hTdGF0ZW1lbnRFcnJvckNvZGVFbnVtUgRjb2RlEj8KBGl0ZW0Yh4CS+wEgAygLMicuZHluYW1v'
    'ZGIuQmF0Y2hTdGF0ZW1lbnRFcnJvci5JdGVtRW50cnlSBGl0ZW0SGwoHbWVzc2FnZRiFs7twIA'
    'EoCVIHbWVzc2FnZRpRCglJdGVtRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSLgoFdmFsdWUYAiAB'
    'KAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdWU6AjgB');

@$core.Deprecated('Use batchStatementRequestDescriptor instead')
const BatchStatementRequest$json = {
  '1': 'BatchStatementRequest',
  '2': [
    {
      '1': 'consistentread',
      '3': 531556994,
      '4': 1,
      '5': 8,
      '10': 'consistentread'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'parameters'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'statement', '3': 248790199, '4': 1, '5': 9, '10': 'statement'},
  ],
};

/// Descriptor for `BatchStatementRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchStatementRequestDescriptor = $convert.base64Decode(
    'ChVCYXRjaFN0YXRlbWVudFJlcXVlc3QSKgoOY29uc2lzdGVudHJlYWQYgtW7/QEgASgIUg5jb2'
    '5zaXN0ZW50cmVhZBI8CgpwYXJhbWV0ZXJzGPqn/usBIAMoCzIYLmR5bmFtb2RiLkF0dHJpYnV0'
    'ZVZhbHVlUgpwYXJhbWV0ZXJzEoIBCiNyZXR1cm52YWx1ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbH'
    'VyZRjgl4ECIAEoDjItLmR5bmFtb2RiLlJldHVyblZhbHVlc09uQ29uZGl0aW9uQ2hlY2tGYWls'
    'dXJlUiNyZXR1cm52YWx1ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZRIfCglzdGF0ZW1lbnQYt/'
    'nQdiABKAlSCXN0YXRlbWVudA==');

@$core.Deprecated('Use batchStatementResponseDescriptor instead')
const BatchStatementResponse$json = {
  '1': 'BatchStatementResponse',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.BatchStatementError',
      '10': 'error'
    },
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchStatementResponse.ItemEntry',
      '10': 'item'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [BatchStatementResponse_ItemEntry$json],
};

@$core.Deprecated('Use batchStatementResponseDescriptor instead')
const BatchStatementResponse_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `BatchStatementResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchStatementResponseDescriptor = $convert.base64Decode(
    'ChZCYXRjaFN0YXRlbWVudFJlc3BvbnNlEjcKBWVycm9yGPK5tpwBIAEoCzIdLmR5bmFtb2RiLk'
    'JhdGNoU3RhdGVtZW50RXJyb3JSBWVycm9yEkIKBGl0ZW0Yh4CS+wEgAygLMiouZHluYW1vZGIu'
    'QmF0Y2hTdGF0ZW1lbnRSZXNwb25zZS5JdGVtRW50cnlSBGl0ZW0SIAoJdGFibGVuYW1lGN3k2o'
    'EBIAEoCVIJdGFibGVuYW1lGlEKCUl0ZW1FbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1'
    'ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use batchWriteItemInputDescriptor instead')
const BatchWriteItemInput$json = {
  '1': 'BatchWriteItemInput',
  '2': [
    {
      '1': 'requestitems',
      '3': 247720687,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchWriteItemInput.RequestitemsEntry',
      '10': 'requestitems'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'returnitemcollectionmetrics',
      '3': 255507354,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnItemCollectionMetrics',
      '10': 'returnitemcollectionmetrics'
    },
  ],
  '3': [BatchWriteItemInput_RequestitemsEntry$json],
};

@$core.Deprecated('Use batchWriteItemInputDescriptor instead')
const BatchWriteItemInput_RequestitemsEntry$json = {
  '1': 'RequestitemsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `BatchWriteItemInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchWriteItemInputDescriptor = $convert.base64Decode(
    'ChNCYXRjaFdyaXRlSXRlbUlucHV0ElYKDHJlcXVlc3RpdGVtcxjv1Y92IAMoCzIvLmR5bmFtb2'
    'RiLkJhdGNoV3JpdGVJdGVtSW5wdXQuUmVxdWVzdGl0ZW1zRW50cnlSDHJlcXVlc3RpdGVtcxJb'
    'ChZyZXR1cm5jb25zdW1lZGNhcGFjaXR5GP7n4RQgASgOMiAuZHluYW1vZGIuUmV0dXJuQ29uc3'
    'VtZWRDYXBhY2l0eVIWcmV0dXJuY29uc3VtZWRjYXBhY2l0eRJqChtyZXR1cm5pdGVtY29sbGVj'
    'dGlvbm1ldHJpY3MYmvfqeSABKA4yJS5keW5hbW9kYi5SZXR1cm5JdGVtQ29sbGVjdGlvbk1ldH'
    'JpY3NSG3JldHVybml0ZW1jb2xsZWN0aW9ubWV0cmljcxo/ChFSZXF1ZXN0aXRlbXNFbnRyeRIQ'
    'CgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use batchWriteItemOutputDescriptor instead')
const BatchWriteItemOutput$json = {
  '1': 'BatchWriteItemOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'itemcollectionmetrics',
      '3': 185317452,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchWriteItemOutput.ItemcollectionmetricsEntry',
      '10': 'itemcollectionmetrics'
    },
    {
      '1': 'unprocesseditems',
      '3': 173269603,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BatchWriteItemOutput.UnprocesseditemsEntry',
      '10': 'unprocesseditems'
    },
  ],
  '3': [
    BatchWriteItemOutput_ItemcollectionmetricsEntry$json,
    BatchWriteItemOutput_UnprocesseditemsEntry$json
  ],
};

@$core.Deprecated('Use batchWriteItemOutputDescriptor instead')
const BatchWriteItemOutput_ItemcollectionmetricsEntry$json = {
  '1': 'ItemcollectionmetricsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use batchWriteItemOutputDescriptor instead')
const BatchWriteItemOutput_UnprocesseditemsEntry$json = {
  '1': 'UnprocesseditemsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `BatchWriteItemOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchWriteItemOutputDescriptor = $convert.base64Decode(
    'ChRCYXRjaFdyaXRlSXRlbU91dHB1dBJKChBjb25zdW1lZGNhcGFjaXR5GKyqodYBIAMoCzIaLm'
    'R5bmFtb2RiLkNvbnN1bWVkQ2FwYWNpdHlSEGNvbnN1bWVkY2FwYWNpdHkScgoVaXRlbWNvbGxl'
    'Y3Rpb25tZXRyaWNzGMzwrlggAygLMjkuZHluYW1vZGIuQmF0Y2hXcml0ZUl0ZW1PdXRwdXQuSX'
    'RlbWNvbGxlY3Rpb25tZXRyaWNzRW50cnlSFWl0ZW1jb2xsZWN0aW9ubWV0cmljcxJjChB1bnBy'
    'b2Nlc3NlZGl0ZW1zGOPEz1IgAygLMjQuZHluYW1vZGIuQmF0Y2hXcml0ZUl0ZW1PdXRwdXQuVW'
    '5wcm9jZXNzZWRpdGVtc0VudHJ5UhB1bnByb2Nlc3NlZGl0ZW1zGkgKGkl0ZW1jb2xsZWN0aW9u'
    'bWV0cmljc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOA'
    'EaQwoVVW5wcm9jZXNzZWRpdGVtc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIg'
    'ASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use billingModeSummaryDescriptor instead')
const BillingModeSummary$json = {
  '1': 'BillingModeSummary',
  '2': [
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'lastupdatetopayperrequestdatetime',
      '3': 40315649,
      '4': 1,
      '5': 9,
      '10': 'lastupdatetopayperrequestdatetime'
    },
  ],
};

/// Descriptor for `BillingModeSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List billingModeSummaryDescriptor = $convert.base64Decode(
    'ChJCaWxsaW5nTW9kZVN1bW1hcnkSOgoLYmlsbGluZ21vZGUYwLToVyABKA4yFS5keW5hbW9kYi'
    '5CaWxsaW5nTW9kZVILYmlsbGluZ21vZGUSTwohbGFzdHVwZGF0ZXRvcGF5cGVycmVxdWVzdGRh'
    'dGV0aW1lGIHWnBMgASgJUiFsYXN0dXBkYXRldG9wYXlwZXJyZXF1ZXN0ZGF0ZXRpbWU=');

@$core.Deprecated('Use cancellationReasonDescriptor instead')
const CancellationReason$json = {
  '1': 'CancellationReason',
  '2': [
    {'1': 'code', '3': 425572629, '4': 1, '5': 9, '10': 'code'},
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.CancellationReason.ItemEntry',
      '10': 'item'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
  '3': [CancellationReason_ItemEntry$json],
};

@$core.Deprecated('Use cancellationReasonDescriptor instead')
const CancellationReason_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `CancellationReason`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancellationReasonDescriptor = $convert.base64Decode(
    'ChJDYW5jZWxsYXRpb25SZWFzb24SFgoEY29kZRiV8vbKASABKAlSBGNvZGUSPgoEaXRlbRiHgJ'
    'L7ASADKAsyJi5keW5hbW9kYi5DYW5jZWxsYXRpb25SZWFzb24uSXRlbUVudHJ5UgRpdGVtEhsK'
    'B21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2UaUQoJSXRlbUVudHJ5EhAKA2tleRgBIAEoCVIDa2'
    'V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use capacityDescriptor instead')
const Capacity$json = {
  '1': 'Capacity',
  '2': [
    {
      '1': 'capacityunits',
      '3': 521922929,
      '4': 1,
      '5': 1,
      '10': 'capacityunits'
    },
    {
      '1': 'readcapacityunits',
      '3': 43945489,
      '4': 1,
      '5': 1,
      '10': 'readcapacityunits'
    },
    {
      '1': 'writecapacityunits',
      '3': 26938032,
      '4': 1,
      '5': 1,
      '10': 'writecapacityunits'
    },
  ],
};

/// Descriptor for `Capacity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityDescriptor = $convert.base64Decode(
    'CghDYXBhY2l0eRIoCg1jYXBhY2l0eXVuaXRzGPHS7/gBIAEoAVINY2FwYWNpdHl1bml0cxIvCh'
    'FyZWFkY2FwYWNpdHl1bml0cxiRnPoUIAEoAVIRcmVhZGNhcGFjaXR5dW5pdHMSMQoSd3JpdGVj'
    'YXBhY2l0eXVuaXRzGLCV7AwgASgBUhJ3cml0ZWNhcGFjaXR5dW5pdHM=');

@$core.Deprecated('Use conditionDescriptor instead')
const Condition$json = {
  '1': 'Condition',
  '2': [
    {
      '1': 'attributevaluelist',
      '3': 157424013,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'attributevaluelist'
    },
    {
      '1': 'comparisonoperator',
      '3': 7059861,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ComparisonOperator',
      '10': 'comparisonoperator'
    },
  ],
};

/// Descriptor for `Condition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conditionDescriptor = $convert.base64Decode(
    'CglDb25kaXRpb24SSwoSYXR0cmlidXRldmFsdWVsaXN0GI2ziEsgAygLMhguZHluYW1vZGIuQX'
    'R0cmlidXRlVmFsdWVSEmF0dHJpYnV0ZXZhbHVlbGlzdBJPChJjb21wYXJpc29ub3BlcmF0b3IY'
    'lfOuAyABKA4yHC5keW5hbW9kYi5Db21wYXJpc29uT3BlcmF0b3JSEmNvbXBhcmlzb25vcGVyYX'
    'Rvcg==');

@$core.Deprecated('Use conditionCheckDescriptor instead')
const ConditionCheck$json = {
  '1': 'ConditionCheck',
  '2': [
    {
      '1': 'conditionexpression',
      '3': 409657405,
      '4': 1,
      '5': 9,
      '10': 'conditionexpression'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConditionCheck.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConditionCheck.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConditionCheck.KeyEntry',
      '10': 'key'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [
    ConditionCheck_ExpressionattributenamesEntry$json,
    ConditionCheck_ExpressionattributevaluesEntry$json,
    ConditionCheck_KeyEntry$json
  ],
};

@$core.Deprecated('Use conditionCheckDescriptor instead')
const ConditionCheck_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use conditionCheckDescriptor instead')
const ConditionCheck_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use conditionCheckDescriptor instead')
const ConditionCheck_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ConditionCheck`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conditionCheckDescriptor = $convert.base64Decode(
    'Cg5Db25kaXRpb25DaGVjaxI0ChNjb25kaXRpb25leHByZXNzaW9uGL3Aq8MBIAEoCVITY29uZG'
    'l0aW9uZXhwcmVzc2lvbhJ1ChhleHByZXNzaW9uYXR0cmlidXRlbmFtZXMY/JjRRyADKAsyNi5k'
    'eW5hbW9kYi5Db25kaXRpb25DaGVjay5FeHByZXNzaW9uYXR0cmlidXRlbmFtZXNFbnRyeVIYZX'
    'hwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzEnkKGWV4cHJlc3Npb25hdHRyaWJ1dGV2YWx1ZXMY2Jyg'
    '5wEgAygLMjcuZHluYW1vZGIuQ29uZGl0aW9uQ2hlY2suRXhwcmVzc2lvbmF0dHJpYnV0ZXZhbH'
    'Vlc0VudHJ5UhlleHByZXNzaW9uYXR0cmlidXRldmFsdWVzEjYKA2tleRiNkutoIAMoCzIhLmR5'
    'bmFtb2RiLkNvbmRpdGlvbkNoZWNrLktleUVudHJ5UgNrZXkSggEKI3JldHVybnZhbHVlc29uY2'
    '9uZGl0aW9uY2hlY2tmYWlsdXJlGOCXgQIgASgOMi0uZHluYW1vZGIuUmV0dXJuVmFsdWVzT25D'
    'b25kaXRpb25DaGVja0ZhaWx1cmVSI3JldHVybnZhbHVlc29uY29uZGl0aW9uY2hlY2tmYWlsdX'
    'JlEiAKCXRhYmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZRpLCh1FeHByZXNzaW9uYXR0cmli'
    'dXRlbmFtZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6Aj'
    'gBGmYKHkV4cHJlc3Npb25hdHRyaWJ1dGV2YWx1ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIu'
    'CgV2YWx1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAEaUAoIS2'
    'V5RW50cnkSEAoDa2V5GAEgASgJUgNrZXkSLgoFdmFsdWUYAiABKAsyGC5keW5hbW9kYi5BdHRy'
    'aWJ1dGVWYWx1ZVIFdmFsdWU6AjgB');

@$core.Deprecated('Use conditionalCheckFailedExceptionDescriptor instead')
const ConditionalCheckFailedException$json = {
  '1': 'ConditionalCheckFailedException',
  '2': [
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConditionalCheckFailedException.ItemEntry',
      '10': 'item'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
  '3': [ConditionalCheckFailedException_ItemEntry$json],
};

@$core.Deprecated('Use conditionalCheckFailedExceptionDescriptor instead')
const ConditionalCheckFailedException_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ConditionalCheckFailedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conditionalCheckFailedExceptionDescriptor = $convert.base64Decode(
    'Ch9Db25kaXRpb25hbENoZWNrRmFpbGVkRXhjZXB0aW9uEksKBGl0ZW0Yh4CS+wEgAygLMjMuZH'
    'luYW1vZGIuQ29uZGl0aW9uYWxDaGVja0ZhaWxlZEV4Y2VwdGlvbi5JdGVtRW50cnlSBGl0ZW0S'
    'GwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZRpRCglJdGVtRW50cnkSEAoDa2V5GAEgASgJUg'
    'NrZXkSLgoFdmFsdWUYAiABKAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdWU6AjgB');

@$core.Deprecated('Use consumedCapacityDescriptor instead')
const ConsumedCapacity$json = {
  '1': 'ConsumedCapacity',
  '2': [
    {
      '1': 'capacityunits',
      '3': 521922929,
      '4': 1,
      '5': 1,
      '10': 'capacityunits'
    },
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity.GlobalsecondaryindexesEntry',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'localsecondaryindexes',
      '3': 362339959,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity.LocalsecondaryindexesEntry',
      '10': 'localsecondaryindexes'
    },
    {
      '1': 'readcapacityunits',
      '3': 43945489,
      '4': 1,
      '5': 1,
      '10': 'readcapacityunits'
    },
    {
      '1': 'table',
      '3': 386722688,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Capacity',
      '10': 'table'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'writecapacityunits',
      '3': 26938032,
      '4': 1,
      '5': 1,
      '10': 'writecapacityunits'
    },
  ],
  '3': [
    ConsumedCapacity_GlobalsecondaryindexesEntry$json,
    ConsumedCapacity_LocalsecondaryindexesEntry$json
  ],
};

@$core.Deprecated('Use consumedCapacityDescriptor instead')
const ConsumedCapacity_GlobalsecondaryindexesEntry$json = {
  '1': 'GlobalsecondaryindexesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Capacity',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use consumedCapacityDescriptor instead')
const ConsumedCapacity_LocalsecondaryindexesEntry$json = {
  '1': 'LocalsecondaryindexesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Capacity',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ConsumedCapacity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List consumedCapacityDescriptor = $convert.base64Decode(
    'ChBDb25zdW1lZENhcGFjaXR5EigKDWNhcGFjaXR5dW5pdHMY8dLv+AEgASgBUg1jYXBhY2l0eX'
    'VuaXRzEnIKFmdsb2JhbHNlY29uZGFyeWluZGV4ZXMYqfqMwwEgAygLMjYuZHluYW1vZGIuQ29u'
    'c3VtZWRDYXBhY2l0eS5HbG9iYWxzZWNvbmRhcnlpbmRleGVzRW50cnlSFmdsb2JhbHNlY29uZG'
    'FyeWluZGV4ZXMSbwoVbG9jYWxzZWNvbmRhcnlpbmRleGVzGPe846wBIAMoCzI1LmR5bmFtb2Ri'
    'LkNvbnN1bWVkQ2FwYWNpdHkuTG9jYWxzZWNvbmRhcnlpbmRleGVzRW50cnlSFWxvY2Fsc2Vjb2'
    '5kYXJ5aW5kZXhlcxIvChFyZWFkY2FwYWNpdHl1bml0cxiRnPoUIAEoAVIRcmVhZGNhcGFjaXR5'
    'dW5pdHMSLAoFdGFibGUYgNezuAEgASgLMhIuZHluYW1vZGIuQ2FwYWNpdHlSBXRhYmxlEiAKCX'
    'RhYmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZRIxChJ3cml0ZWNhcGFjaXR5dW5pdHMYsJXs'
    'DCABKAFSEndyaXRlY2FwYWNpdHl1bml0cxpdChtHbG9iYWxzZWNvbmRhcnlpbmRleGVzRW50cn'
    'kSEAoDa2V5GAEgASgJUgNrZXkSKAoFdmFsdWUYAiABKAsyEi5keW5hbW9kYi5DYXBhY2l0eVIF'
    'dmFsdWU6AjgBGlwKGkxvY2Fsc2Vjb25kYXJ5aW5kZXhlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2'
    'V5EigKBXZhbHVlGAIgASgLMhIuZHluYW1vZGIuQ2FwYWNpdHlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use continuousBackupsDescriptionDescriptor instead')
const ContinuousBackupsDescription$json = {
  '1': 'ContinuousBackupsDescription',
  '2': [
    {
      '1': 'continuousbackupsstatus',
      '3': 501975808,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContinuousBackupsStatus',
      '10': 'continuousbackupsstatus'
    },
    {
      '1': 'pointintimerecoverydescription',
      '3': 319300211,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.PointInTimeRecoveryDescription',
      '10': 'pointintimerecoverydescription'
    },
  ],
};

/// Descriptor for `ContinuousBackupsDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousBackupsDescriptionDescriptor = $convert.base64Decode(
    'ChxDb250aW51b3VzQmFja3Vwc0Rlc2NyaXB0aW9uEl8KF2NvbnRpbnVvdXNiYWNrdXBzc3RhdH'
    'VzGICWru8BIAEoDjIhLmR5bmFtb2RiLkNvbnRpbnVvdXNCYWNrdXBzU3RhdHVzUhdjb250aW51'
    'b3VzYmFja3Vwc3N0YXR1cxJ0Ch5wb2ludGludGltZXJlY292ZXJ5ZGVzY3JpcHRpb24Y88SgmA'
    'EgASgLMiguZHluYW1vZGIuUG9pbnRJblRpbWVSZWNvdmVyeURlc2NyaXB0aW9uUh5wb2ludGlu'
    'dGltZXJlY292ZXJ5ZGVzY3JpcHRpb24=');

@$core.Deprecated('Use continuousBackupsUnavailableExceptionDescriptor instead')
const ContinuousBackupsUnavailableException$json = {
  '1': 'ContinuousBackupsUnavailableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ContinuousBackupsUnavailableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuousBackupsUnavailableExceptionDescriptor =
    $convert.base64Decode(
        'CiVDb250aW51b3VzQmFja3Vwc1VuYXZhaWxhYmxlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use contributorInsightsSummaryDescriptor instead')
const ContributorInsightsSummary$json = {
  '1': 'ContributorInsightsSummary',
  '2': [
    {
      '1': 'contributorinsightsmode',
      '3': 86700161,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsMode',
      '10': 'contributorinsightsmode'
    },
    {
      '1': 'contributorinsightsstatus',
      '3': 363347282,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsStatus',
      '10': 'contributorinsightsstatus'
    },
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `ContributorInsightsSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contributorInsightsSummaryDescriptor = $convert.base64Decode(
    'ChpDb250cmlidXRvckluc2lnaHRzU3VtbWFyeRJeChdjb250cmlidXRvcmluc2lnaHRzbW9kZR'
    'iB4aspIAEoDjIhLmR5bmFtb2RiLkNvbnRyaWJ1dG9ySW5zaWdodHNNb2RlUhdjb250cmlidXRv'
    'cmluc2lnaHRzbW9kZRJlChljb250cmlidXRvcmluc2lnaHRzc3RhdHVzGNL6oK0BIAEoDjIjLm'
    'R5bmFtb2RiLkNvbnRyaWJ1dG9ySW5zaWdodHNTdGF0dXNSGWNvbnRyaWJ1dG9yaW5zaWdodHNz'
    'dGF0dXMSHwoJaW5kZXhuYW1lGJHV6zAgASgJUglpbmRleG5hbWUSIAoJdGFibGVuYW1lGN3k2o'
    'EBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use createBackupInputDescriptor instead')
const CreateBackupInput$json = {
  '1': 'CreateBackupInput',
  '2': [
    {'1': 'backupname', '3': 467693789, '4': 1, '5': 9, '10': 'backupname'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `CreateBackupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBackupInputDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVCYWNrdXBJbnB1dBIiCgpiYWNrdXBuYW1lGN3hgd8BIAEoCVIKYmFja3VwbmFtZR'
    'IgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW5hbWU=');

@$core.Deprecated('Use createBackupOutputDescriptor instead')
const CreateBackupOutput$json = {
  '1': 'CreateBackupOutput',
  '2': [
    {
      '1': 'backupdetails',
      '3': 383357432,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.BackupDetails',
      '10': 'backupdetails'
    },
  ],
};

/// Descriptor for `CreateBackupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBackupOutputDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVCYWNrdXBPdXRwdXQSQQoNYmFja3VwZGV0YWlscxj4o+a2ASABKAsyFy5keW5hbW'
    '9kYi5CYWNrdXBEZXRhaWxzUg1iYWNrdXBkZXRhaWxz');

@$core.Deprecated('Use createGlobalSecondaryIndexActionDescriptor instead')
const CreateGlobalSecondaryIndexAction$json = {
  '1': 'CreateGlobalSecondaryIndexAction',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'projection',
      '3': 105045921,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Projection',
      '10': 'projection'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.WarmThroughput',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `CreateGlobalSecondaryIndexAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGlobalSecondaryIndexActionDescriptor = $convert.base64Decode(
    'CiBDcmVhdGVHbG9iYWxTZWNvbmRhcnlJbmRleEFjdGlvbhIfCglpbmRleG5hbWUYkdXrMCABKA'
    'lSCWluZGV4bmFtZRI8CglrZXlzY2hlbWEY6M/diwEgAygLMhouZHluYW1vZGIuS2V5U2NoZW1h'
    'RWxlbWVudFIJa2V5c2NoZW1hElAKEm9uZGVtYW5kdGhyb3VnaHB1dBiC3trlASABKAsyHC5keW'
    '5hbW9kYi5PbkRlbWFuZFRocm91Z2hwdXRSEm9uZGVtYW5kdGhyb3VnaHB1dBI3Cgpwcm9qZWN0'
    'aW9uGKG/izIgASgLMhQuZHluYW1vZGIuUHJvamVjdGlvblIKcHJvamVjdGlvbhJXChVwcm92aX'
    'Npb25lZHRocm91Z2hwdXQYjKNrIAEoCzIfLmR5bmFtb2RiLlByb3Zpc2lvbmVkVGhyb3VnaHB1'
    'dFIVcHJvdmlzaW9uZWR0aHJvdWdocHV0EkQKDndhcm10aHJvdWdocHV0GIPeyIoBIAEoCzIYLm'
    'R5bmFtb2RiLldhcm1UaHJvdWdocHV0Ug53YXJtdGhyb3VnaHB1dA==');

@$core.Deprecated('Use createGlobalTableInputDescriptor instead')
const CreateGlobalTableInput$json = {
  '1': 'CreateGlobalTableInput',
  '2': [
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
    {
      '1': 'replicationgroup',
      '3': 190970785,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Replica',
      '10': 'replicationgroup'
    },
  ],
};

/// Descriptor for `CreateGlobalTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGlobalTableInputDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVHbG9iYWxUYWJsZUlucHV0EiwKD2dsb2JhbHRhYmxlbmFtZRiqpqeHASABKAlSD2'
    'dsb2JhbHRhYmxlbmFtZRJAChByZXBsaWNhdGlvbmdyb3VwGKH3h1sgAygLMhEuZHluYW1vZGIu'
    'UmVwbGljYVIQcmVwbGljYXRpb25ncm91cA==');

@$core.Deprecated('Use createGlobalTableOutputDescriptor instead')
const CreateGlobalTableOutput$json = {
  '1': 'CreateGlobalTableOutput',
  '2': [
    {
      '1': 'globaltabledescription',
      '3': 342116405,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.GlobalTableDescription',
      '10': 'globaltabledescription'
    },
  ],
};

/// Descriptor for `CreateGlobalTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGlobalTableOutputDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVHbG9iYWxUYWJsZU91dHB1dBJcChZnbG9iYWx0YWJsZWRlc2NyaXB0aW9uGLWQka'
    'MBIAEoCzIgLmR5bmFtb2RiLkdsb2JhbFRhYmxlRGVzY3JpcHRpb25SFmdsb2JhbHRhYmxlZGVz'
    'Y3JpcHRpb24=');

@$core.Deprecated(
    'Use createGlobalTableWitnessGroupMemberActionDescriptor instead')
const CreateGlobalTableWitnessGroupMemberAction$json = {
  '1': 'CreateGlobalTableWitnessGroupMemberAction',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
  ],
};

/// Descriptor for `CreateGlobalTableWitnessGroupMemberAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createGlobalTableWitnessGroupMemberActionDescriptor = $convert.base64Decode(
        'CilDcmVhdGVHbG9iYWxUYWJsZVdpdG5lc3NHcm91cE1lbWJlckFjdGlvbhIhCgpyZWdpb25uYW'
        '1lGL+buTUgASgJUgpyZWdpb25uYW1l');

@$core.Deprecated('Use createReplicaActionDescriptor instead')
const CreateReplicaAction$json = {
  '1': 'CreateReplicaAction',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
  ],
};

/// Descriptor for `CreateReplicaAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createReplicaActionDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVSZXBsaWNhQWN0aW9uEiEKCnJlZ2lvbm5hbWUYv5u5NSABKAlSCnJlZ2lvbm5hbW'
    'U=');

@$core.Deprecated('Use createReplicationGroupMemberActionDescriptor instead')
const CreateReplicationGroupMemberAction$json = {
  '1': 'CreateReplicationGroupMemberAction',
  '2': [
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaGlobalSecondaryIndex',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'kmsmasterkeyid',
      '3': 521459443,
      '4': 1,
      '5': 9,
      '10': 'kmsmasterkeyid'
    },
    {
      '1': 'ondemandthroughputoverride',
      '3': 317165234,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughputOverride',
      '10': 'ondemandthroughputoverride'
    },
    {
      '1': 'provisionedthroughputoverride',
      '3': 413332116,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughputOverride',
      '10': 'provisionedthroughputoverride'
    },
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'tableclassoverride',
      '3': 415569842,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableClass',
      '10': 'tableclassoverride'
    },
  ],
};

/// Descriptor for `CreateReplicationGroupMemberAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createReplicationGroupMemberActionDescriptor = $convert.base64Decode(
    'CiJDcmVhdGVSZXBsaWNhdGlvbkdyb3VwTWVtYmVyQWN0aW9uEmEKFmdsb2JhbHNlY29uZGFyeW'
    'luZGV4ZXMYqfqMwwEgAygLMiUuZHluYW1vZGIuUmVwbGljYUdsb2JhbFNlY29uZGFyeUluZGV4'
    'UhZnbG9iYWxzZWNvbmRhcnlpbmRleGVzEioKDmttc21hc3RlcmtleWlkGPOt0/gBIAEoCVIOa2'
    '1zbWFzdGVya2V5aWQSaAoab25kZW1hbmR0aHJvdWdocHV0b3ZlcnJpZGUYsp2elwEgASgLMiQu'
    'ZHluYW1vZGIuT25EZW1hbmRUaHJvdWdocHV0T3ZlcnJpZGVSGm9uZGVtYW5kdGhyb3VnaHB1dG'
    '92ZXJyaWRlEnEKHXByb3Zpc2lvbmVkdGhyb3VnaHB1dG92ZXJyaWRlGJTli8UBIAEoCzInLmR5'
    'bmFtb2RiLlByb3Zpc2lvbmVkVGhyb3VnaHB1dE92ZXJyaWRlUh1wcm92aXNpb25lZHRocm91Z2'
    'hwdXRvdmVycmlkZRIhCgpyZWdpb25uYW1lGL+buTUgASgJUgpyZWdpb25uYW1lEkgKEnRhYmxl'
    'Y2xhc3NvdmVycmlkZRiyr5TGASABKA4yFC5keW5hbW9kYi5UYWJsZUNsYXNzUhJ0YWJsZWNsYX'
    'Nzb3ZlcnJpZGU=');

@$core.Deprecated('Use createTableInputDescriptor instead')
const CreateTableInput$json = {
  '1': 'CreateTableInput',
  '2': [
    {
      '1': 'attributedefinitions',
      '3': 414687108,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeDefinition',
      '10': 'attributedefinitions'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'deletionprotectionenabled',
      '3': 259418450,
      '4': 1,
      '5': 8,
      '10': 'deletionprotectionenabled'
    },
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndex',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'globaltablesettingsreplicationmode',
      '3': 10446577,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.GlobalTableSettingsReplicationMode',
      '10': 'globaltablesettingsreplicationmode'
    },
    {
      '1': 'globaltablesourcearn',
      '3': 443614787,
      '4': 1,
      '5': 9,
      '10': 'globaltablesourcearn'
    },
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'localsecondaryindexes',
      '3': 362339959,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.LocalSecondaryIndex',
      '10': 'localsecondaryindexes'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
    {
      '1': 'ssespecification',
      '3': 31692444,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SSESpecification',
      '10': 'ssespecification'
    },
    {
      '1': 'streamspecification',
      '3': 403922627,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.StreamSpecification',
      '10': 'streamspecification'
    },
    {
      '1': 'tableclass',
      '3': 342890498,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableClass',
      '10': 'tableclass'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Tag',
      '10': 'tags'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.WarmThroughput',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `CreateTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTableInputDescriptor = $convert.base64Decode(
    'ChBDcmVhdGVUYWJsZUlucHV0ElUKFGF0dHJpYnV0ZWRlZmluaXRpb25zGIS/3sUBIAMoCzIdLm'
    'R5bmFtb2RiLkF0dHJpYnV0ZURlZmluaXRpb25SFGF0dHJpYnV0ZWRlZmluaXRpb25zEjoKC2Jp'
    'bGxpbmdtb2RlGMC06FcgASgOMhUuZHluYW1vZGIuQmlsbGluZ01vZGVSC2JpbGxpbmdtb2RlEj'
    '8KGWRlbGV0aW9ucHJvdGVjdGlvbmVuYWJsZWQY0tLZeyABKAhSGWRlbGV0aW9ucHJvdGVjdGlv'
    'bmVuYWJsZWQSWgoWZ2xvYmFsc2Vjb25kYXJ5aW5kZXhlcxip+ozDASADKAsyHi5keW5hbW9kYi'
    '5HbG9iYWxTZWNvbmRhcnlJbmRleFIWZ2xvYmFsc2Vjb25kYXJ5aW5kZXhlcxJ/CiJnbG9iYWx0'
    'YWJsZXNldHRpbmdzcmVwbGljYXRpb25tb2RlGPHN/QQgASgOMiwuZHluYW1vZGIuR2xvYmFsVG'
    'FibGVTZXR0aW5nc1JlcGxpY2F0aW9uTW9kZVIiZ2xvYmFsdGFibGVzZXR0aW5nc3JlcGxpY2F0'
    'aW9ubW9kZRI2ChRnbG9iYWx0YWJsZXNvdXJjZWFybhjDjMTTASABKAlSFGdsb2JhbHRhYmxlc2'
    '91cmNlYXJuEjwKCWtleXNjaGVtYRjoz92LASADKAsyGi5keW5hbW9kYi5LZXlTY2hlbWFFbGVt'
    'ZW50UglrZXlzY2hlbWESVwoVbG9jYWxzZWNvbmRhcnlpbmRleGVzGPe846wBIAMoCzIdLmR5bm'
    'Ftb2RiLkxvY2FsU2Vjb25kYXJ5SW5kZXhSFWxvY2Fsc2Vjb25kYXJ5aW5kZXhlcxJQChJvbmRl'
    'bWFuZHRocm91Z2hwdXQYgt7a5QEgASgLMhwuZHluYW1vZGIuT25EZW1hbmRUaHJvdWdocHV0Uh'
    'JvbmRlbWFuZHRocm91Z2hwdXQSVwoVcHJvdmlzaW9uZWR0aHJvdWdocHV0GIyjayABKAsyHy5k'
    'eW5hbW9kYi5Qcm92aXNpb25lZFRocm91Z2hwdXRSFXByb3Zpc2lvbmVkdGhyb3VnaHB1dBIpCg'
    '5yZXNvdXJjZXBvbGljeRiwlMEHIAEoCVIOcmVzb3VyY2Vwb2xpY3kSSQoQc3Nlc3BlY2lmaWNh'
    'dGlvbhicrY4PIAEoCzIaLmR5bmFtb2RiLlNTRVNwZWNpZmljYXRpb25SEHNzZXNwZWNpZmljYX'
    'Rpb24SUwoTc3RyZWFtc3BlY2lmaWNhdGlvbhjDvc3AASABKAsyHS5keW5hbW9kYi5TdHJlYW1T'
    'cGVjaWZpY2F0aW9uUhNzdHJlYW1zcGVjaWZpY2F0aW9uEjgKCnRhYmxlY2xhc3MYgrDAowEgAS'
    'gOMhQuZHluYW1vZGIuVGFibGVDbGFzc1IKdGFibGVjbGFzcxIgCgl0YWJsZW5hbWUY3eTagQEg'
    'ASgJUgl0YWJsZW5hbWUSJQoEdGFncxjBwfa1ASADKAsyDS5keW5hbW9kYi5UYWdSBHRhZ3MSRA'
    'oOd2FybXRocm91Z2hwdXQYg97IigEgASgLMhguZHluYW1vZGIuV2FybVRocm91Z2hwdXRSDndh'
    'cm10aHJvdWdocHV0');

@$core.Deprecated('Use createTableOutputDescriptor instead')
const CreateTableOutput$json = {
  '1': 'CreateTableOutput',
  '2': [
    {
      '1': 'tabledescription',
      '3': 19962388,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableDescription',
      '10': 'tabledescription'
    },
  ],
};

/// Descriptor for `CreateTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTableOutputDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVUYWJsZU91dHB1dBJJChB0YWJsZWRlc2NyaXB0aW9uGJS0wgkgASgLMhouZHluYW'
    '1vZGIuVGFibGVEZXNjcmlwdGlvblIQdGFibGVkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use csvOptionsDescriptor instead')
const CsvOptions$json = {
  '1': 'CsvOptions',
  '2': [
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {'1': 'headerlist', '3': 424018665, '4': 3, '5': 9, '10': 'headerlist'},
  ],
};

/// Descriptor for `CsvOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List csvOptionsDescriptor = $convert.base64Decode(
    'CgpDc3ZPcHRpb25zEiAKCWRlbGltaXRlchib2YiQASABKAlSCWRlbGltaXRlchIiCgpoZWFkZX'
    'JsaXN0GOmFmMoBIAMoCVIKaGVhZGVybGlzdA==');

@$core.Deprecated('Use deleteDescriptor instead')
const Delete$json = {
  '1': 'Delete',
  '2': [
    {
      '1': 'conditionexpression',
      '3': 409657405,
      '4': 1,
      '5': 9,
      '10': 'conditionexpression'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Delete.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Delete.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Delete.KeyEntry',
      '10': 'key'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [
    Delete_ExpressionattributenamesEntry$json,
    Delete_ExpressionattributevaluesEntry$json,
    Delete_KeyEntry$json
  ],
};

@$core.Deprecated('Use deleteDescriptor instead')
const Delete_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use deleteDescriptor instead')
const Delete_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use deleteDescriptor instead')
const Delete_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `Delete`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDescriptor = $convert.base64Decode(
    'CgZEZWxldGUSNAoTY29uZGl0aW9uZXhwcmVzc2lvbhi9wKvDASABKAlSE2NvbmRpdGlvbmV4cH'
    'Jlc3Npb24SbQoYZXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMi4uZHluYW1vZGIu'
    'RGVsZXRlLkV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9uYXR0cmlidX'
    'RlbmFtZXMScQoZZXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlcxjYnKDnASADKAsyLy5keW5hbW9k'
    'Yi5EZWxldGUuRXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5UhlleHByZXNzaW9uYXR0cm'
    'lidXRldmFsdWVzEi4KA2tleRiNkutoIAMoCzIZLmR5bmFtb2RiLkRlbGV0ZS5LZXlFbnRyeVID'
    'a2V5EoIBCiNyZXR1cm52YWx1ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZRjgl4ECIAEoDjItLm'
    'R5bmFtb2RiLlJldHVyblZhbHVlc09uQ29uZGl0aW9uQ2hlY2tGYWlsdXJlUiNyZXR1cm52YWx1'
    'ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZRIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW'
    '5hbWUaSwodRXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkS'
    'FAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpmCh5FeHByZXNzaW9uYXR0cmlidXRldmFsdWVzRW'
    '50cnkSEAoDa2V5GAEgASgJUgNrZXkSLgoFdmFsdWUYAiABKAsyGC5keW5hbW9kYi5BdHRyaWJ1'
    'dGVWYWx1ZVIFdmFsdWU6AjgBGlAKCEtleUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbH'
    'VlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use deleteBackupInputDescriptor instead')
const DeleteBackupInput$json = {
  '1': 'DeleteBackupInput',
  '2': [
    {'1': 'backuparn', '3': 370874339, '4': 1, '5': 9, '10': 'backuparn'},
  ],
};

/// Descriptor for `DeleteBackupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBackupInputDescriptor = $convert.base64Decode(
    'ChFEZWxldGVCYWNrdXBJbnB1dBIgCgliYWNrdXBhcm4Y46/ssAEgASgJUgliYWNrdXBhcm4=');

@$core.Deprecated('Use deleteBackupOutputDescriptor instead')
const DeleteBackupOutput$json = {
  '1': 'DeleteBackupOutput',
  '2': [
    {
      '1': 'backupdescription',
      '3': 294936980,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.BackupDescription',
      '10': 'backupdescription'
    },
  ],
};

/// Descriptor for `DeleteBackupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBackupOutputDescriptor = $convert.base64Decode(
    'ChJEZWxldGVCYWNrdXBPdXRwdXQSTQoRYmFja3VwZGVzY3JpcHRpb24YlMPRjAEgASgLMhsuZH'
    'luYW1vZGIuQmFja3VwRGVzY3JpcHRpb25SEWJhY2t1cGRlc2NyaXB0aW9u');

@$core.Deprecated('Use deleteGlobalSecondaryIndexActionDescriptor instead')
const DeleteGlobalSecondaryIndexAction$json = {
  '1': 'DeleteGlobalSecondaryIndexAction',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
  ],
};

/// Descriptor for `DeleteGlobalSecondaryIndexAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteGlobalSecondaryIndexActionDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVHbG9iYWxTZWNvbmRhcnlJbmRleEFjdGlvbhIfCglpbmRleG5hbWUYkdXrMCABKA'
        'lSCWluZGV4bmFtZQ==');

@$core.Deprecated(
    'Use deleteGlobalTableWitnessGroupMemberActionDescriptor instead')
const DeleteGlobalTableWitnessGroupMemberAction$json = {
  '1': 'DeleteGlobalTableWitnessGroupMemberAction',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
  ],
};

/// Descriptor for `DeleteGlobalTableWitnessGroupMemberAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteGlobalTableWitnessGroupMemberActionDescriptor = $convert.base64Decode(
        'CilEZWxldGVHbG9iYWxUYWJsZVdpdG5lc3NHcm91cE1lbWJlckFjdGlvbhIhCgpyZWdpb25uYW'
        '1lGL+buTUgASgJUgpyZWdpb25uYW1l');

@$core.Deprecated('Use deleteItemInputDescriptor instead')
const DeleteItemInput$json = {
  '1': 'DeleteItemInput',
  '2': [
    {
      '1': 'conditionexpression',
      '3': 409657405,
      '4': 1,
      '5': 9,
      '10': 'conditionexpression'
    },
    {
      '1': 'conditionaloperator',
      '3': 172066260,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ConditionalOperator',
      '10': 'conditionaloperator'
    },
    {
      '1': 'expected',
      '3': 106557946,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.DeleteItemInput.ExpectedEntry',
      '10': 'expected'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.DeleteItemInput.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.DeleteItemInput.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.DeleteItemInput.KeyEntry',
      '10': 'key'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'returnitemcollectionmetrics',
      '3': 255507354,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnItemCollectionMetrics',
      '10': 'returnitemcollectionmetrics'
    },
    {
      '1': 'returnvalues',
      '3': 402960198,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValue',
      '10': 'returnvalues'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [
    DeleteItemInput_ExpectedEntry$json,
    DeleteItemInput_ExpressionattributenamesEntry$json,
    DeleteItemInput_ExpressionattributevaluesEntry$json,
    DeleteItemInput_KeyEntry$json
  ],
};

@$core.Deprecated('Use deleteItemInputDescriptor instead')
const DeleteItemInput_ExpectedEntry$json = {
  '1': 'ExpectedEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ExpectedAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use deleteItemInputDescriptor instead')
const DeleteItemInput_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use deleteItemInputDescriptor instead')
const DeleteItemInput_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use deleteItemInputDescriptor instead')
const DeleteItemInput_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `DeleteItemInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteItemInputDescriptor = $convert.base64Decode(
    'Cg9EZWxldGVJdGVtSW5wdXQSNAoTY29uZGl0aW9uZXhwcmVzc2lvbhi9wKvDASABKAlSE2Nvbm'
    'RpdGlvbmV4cHJlc3Npb24SUgoTY29uZGl0aW9uYWxvcGVyYXRvchjUi4ZSIAEoDjIdLmR5bmFt'
    'b2RiLkNvbmRpdGlvbmFsT3BlcmF0b3JSE2NvbmRpdGlvbmFsb3BlcmF0b3ISRgoIZXhwZWN0ZW'
    'QY+uPnMiADKAsyJy5keW5hbW9kYi5EZWxldGVJdGVtSW5wdXQuRXhwZWN0ZWRFbnRyeVIIZXhw'
    'ZWN0ZWQSdgoYZXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMjcuZHluYW1vZGIuRG'
    'VsZXRlSXRlbUlucHV0LkV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9u'
    'YXR0cmlidXRlbmFtZXMSegoZZXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlcxjYnKDnASADKAsyOC'
    '5keW5hbW9kYi5EZWxldGVJdGVtSW5wdXQuRXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5'
    'UhlleHByZXNzaW9uYXR0cmlidXRldmFsdWVzEjcKA2tleRiNkutoIAMoCzIiLmR5bmFtb2RiLk'
    'RlbGV0ZUl0ZW1JbnB1dC5LZXlFbnRyeVIDa2V5ElsKFnJldHVybmNvbnN1bWVkY2FwYWNpdHkY'
    '/ufhFCABKA4yIC5keW5hbW9kYi5SZXR1cm5Db25zdW1lZENhcGFjaXR5UhZyZXR1cm5jb25zdW'
    '1lZGNhcGFjaXR5EmoKG3JldHVybml0ZW1jb2xsZWN0aW9ubWV0cmljcxia9+p5IAEoDjIlLmR5'
    'bmFtb2RiLlJldHVybkl0ZW1Db2xsZWN0aW9uTWV0cmljc1IbcmV0dXJuaXRlbWNvbGxlY3Rpb2'
    '5tZXRyaWNzEj0KDHJldHVybnZhbHVlcxjG3pLAASABKA4yFS5keW5hbW9kYi5SZXR1cm5WYWx1'
    'ZVIMcmV0dXJudmFsdWVzEoIBCiNyZXR1cm52YWx1ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZR'
    'jgl4ECIAEoDjItLmR5bmFtb2RiLlJldHVyblZhbHVlc09uQ29uZGl0aW9uQ2hlY2tGYWlsdXJl'
    'UiNyZXR1cm52YWx1ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZRIgCgl0YWJsZW5hbWUY3eTagQ'
    'EgASgJUgl0YWJsZW5hbWUaXQoNRXhwZWN0ZWRFbnRyeRIQCgNrZXkYASABKAlSA2tleRI2CgV2'
    'YWx1ZRgCIAEoCzIgLmR5bmFtb2RiLkV4cGVjdGVkQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AR'
    'pLCh1FeHByZXNzaW9uYXR0cmlidXRlbmFtZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2'
    'YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGmYKHkV4cHJlc3Npb25hdHRyaWJ1dGV2YWx1ZXNFbnRyeR'
    'IQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZh'
    'bHVlUgV2YWx1ZToCOAEaUAoIS2V5RW50cnkSEAoDa2V5GAEgASgJUgNrZXkSLgoFdmFsdWUYAi'
    'ABKAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdWU6AjgB');

@$core.Deprecated('Use deleteItemOutputDescriptor instead')
const DeleteItemOutput$json = {
  '1': 'DeleteItemOutput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.DeleteItemOutput.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'itemcollectionmetrics',
      '3': 185317452,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ItemCollectionMetrics',
      '10': 'itemcollectionmetrics'
    },
  ],
  '3': [DeleteItemOutput_AttributesEntry$json],
};

@$core.Deprecated('Use deleteItemOutputDescriptor instead')
const DeleteItemOutput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `DeleteItemOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteItemOutputDescriptor = $convert.base64Decode(
    'ChBEZWxldGVJdGVtT3V0cHV0Ek0KCmF0dHJpYnV0ZXMYtan7YyADKAsyKi5keW5hbW9kYi5EZW'
    'xldGVJdGVtT3V0cHV0LkF0dHJpYnV0ZXNFbnRyeVIKYXR0cmlidXRlcxJKChBjb25zdW1lZGNh'
    'cGFjaXR5GKyqodYBIAEoCzIaLmR5bmFtb2RiLkNvbnN1bWVkQ2FwYWNpdHlSEGNvbnN1bWVkY2'
    'FwYWNpdHkSWAoVaXRlbWNvbGxlY3Rpb25tZXRyaWNzGMzwrlggASgLMh8uZHluYW1vZGIuSXRl'
    'bUNvbGxlY3Rpb25NZXRyaWNzUhVpdGVtY29sbGVjdGlvbm1ldHJpY3MaVwoPQXR0cmlidXRlc0'
    'VudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmli'
    'dXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use deleteReplicaActionDescriptor instead')
const DeleteReplicaAction$json = {
  '1': 'DeleteReplicaAction',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
  ],
};

/// Descriptor for `DeleteReplicaAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteReplicaActionDescriptor = $convert.base64Decode(
    'ChNEZWxldGVSZXBsaWNhQWN0aW9uEiEKCnJlZ2lvbm5hbWUYv5u5NSABKAlSCnJlZ2lvbm5hbW'
    'U=');

@$core.Deprecated('Use deleteReplicationGroupMemberActionDescriptor instead')
const DeleteReplicationGroupMemberAction$json = {
  '1': 'DeleteReplicationGroupMemberAction',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
  ],
};

/// Descriptor for `DeleteReplicationGroupMemberAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteReplicationGroupMemberActionDescriptor =
    $convert.base64Decode(
        'CiJEZWxldGVSZXBsaWNhdGlvbkdyb3VwTWVtYmVyQWN0aW9uEiEKCnJlZ2lvbm5hbWUYv5u5NS'
        'ABKAlSCnJlZ2lvbm5hbWU=');

@$core.Deprecated('Use deleteRequestDescriptor instead')
const DeleteRequest$json = {
  '1': 'DeleteRequest',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.DeleteRequest.KeyEntry',
      '10': 'key'
    },
  ],
  '3': [DeleteRequest_KeyEntry$json],
};

@$core.Deprecated('Use deleteRequestDescriptor instead')
const DeleteRequest_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `DeleteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRequestDescriptor = $convert.base64Decode(
    'Cg1EZWxldGVSZXF1ZXN0EjUKA2tleRiNkutoIAMoCzIgLmR5bmFtb2RiLkRlbGV0ZVJlcXVlc3'
    'QuS2V5RW50cnlSA2tleRpQCghLZXlFbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1ZRgC'
    'IAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use deleteResourcePolicyInputDescriptor instead')
const DeleteResourcePolicyInput$json = {
  '1': 'DeleteResourcePolicyInput',
  '2': [
    {
      '1': 'expectedrevisionid',
      '3': 456463970,
      '4': 1,
      '5': 9,
      '10': 'expectedrevisionid'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `DeleteResourcePolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyInputDescriptor = $convert.base64Decode(
    'ChlEZWxldGVSZXNvdXJjZVBvbGljeUlucHV0EjIKEmV4cGVjdGVkcmV2aXNpb25pZBjirNTZAS'
    'ABKAlSEmV4cGVjdGVkcmV2aXNpb25pZBIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc291'
    'cmNlYXJu');

@$core.Deprecated('Use deleteResourcePolicyOutputDescriptor instead')
const DeleteResourcePolicyOutput$json = {
  '1': 'DeleteResourcePolicyOutput',
  '2': [
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `DeleteResourcePolicyOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyOutputDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVSZXNvdXJjZVBvbGljeU91dHB1dBIiCgpyZXZpc2lvbmlkGIajnu4BIAEoCVIKcm'
        'V2aXNpb25pZA==');

@$core.Deprecated('Use deleteTableInputDescriptor instead')
const DeleteTableInput$json = {
  '1': 'DeleteTableInput',
  '2': [
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DeleteTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTableInputDescriptor = $convert.base64Decode(
    'ChBEZWxldGVUYWJsZUlucHV0EiAKCXRhYmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZQ==');

@$core.Deprecated('Use deleteTableOutputDescriptor instead')
const DeleteTableOutput$json = {
  '1': 'DeleteTableOutput',
  '2': [
    {
      '1': 'tabledescription',
      '3': 19962388,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableDescription',
      '10': 'tabledescription'
    },
  ],
};

/// Descriptor for `DeleteTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTableOutputDescriptor = $convert.base64Decode(
    'ChFEZWxldGVUYWJsZU91dHB1dBJJChB0YWJsZWRlc2NyaXB0aW9uGJS0wgkgASgLMhouZHluYW'
    '1vZGIuVGFibGVEZXNjcmlwdGlvblIQdGFibGVkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use describeBackupInputDescriptor instead')
const DescribeBackupInput$json = {
  '1': 'DescribeBackupInput',
  '2': [
    {'1': 'backuparn', '3': 370874339, '4': 1, '5': 9, '10': 'backuparn'},
  ],
};

/// Descriptor for `DescribeBackupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeBackupInputDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZUJhY2t1cElucHV0EiAKCWJhY2t1cGFybhjjr+ywASABKAlSCWJhY2t1cGFybg'
    '==');

@$core.Deprecated('Use describeBackupOutputDescriptor instead')
const DescribeBackupOutput$json = {
  '1': 'DescribeBackupOutput',
  '2': [
    {
      '1': 'backupdescription',
      '3': 294936980,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.BackupDescription',
      '10': 'backupdescription'
    },
  ],
};

/// Descriptor for `DescribeBackupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeBackupOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZUJhY2t1cE91dHB1dBJNChFiYWNrdXBkZXNjcmlwdGlvbhiUw9GMASABKAsyGy'
    '5keW5hbW9kYi5CYWNrdXBEZXNjcmlwdGlvblIRYmFja3VwZGVzY3JpcHRpb24=');

@$core.Deprecated('Use describeContinuousBackupsInputDescriptor instead')
const DescribeContinuousBackupsInput$json = {
  '1': 'DescribeContinuousBackupsInput',
  '2': [
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeContinuousBackupsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeContinuousBackupsInputDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZUNvbnRpbnVvdXNCYWNrdXBzSW5wdXQSIAoJdGFibGVuYW1lGN3k2oEBIAEoCV'
        'IJdGFibGVuYW1l');

@$core.Deprecated('Use describeContinuousBackupsOutputDescriptor instead')
const DescribeContinuousBackupsOutput$json = {
  '1': 'DescribeContinuousBackupsOutput',
  '2': [
    {
      '1': 'continuousbackupsdescription',
      '3': 446030650,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ContinuousBackupsDescription',
      '10': 'continuousbackupsdescription'
    },
  ],
};

/// Descriptor for `DescribeContinuousBackupsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeContinuousBackupsOutputDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZUNvbnRpbnVvdXNCYWNrdXBzT3V0cHV0Em4KHGNvbnRpbnVvdXNiYWNrdXBzZG'
        'VzY3JpcHRpb24YusbX1AEgASgLMiYuZHluYW1vZGIuQ29udGludW91c0JhY2t1cHNEZXNjcmlw'
        'dGlvblIcY29udGludW91c2JhY2t1cHNkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use describeContributorInsightsInputDescriptor instead')
const DescribeContributorInsightsInput$json = {
  '1': 'DescribeContributorInsightsInput',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeContributorInsightsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeContributorInsightsInputDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZUNvbnRyaWJ1dG9ySW5zaWdodHNJbnB1dBIfCglpbmRleG5hbWUYkdXrMCABKA'
        'lSCWluZGV4bmFtZRIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW5hbWU=');

@$core.Deprecated('Use describeContributorInsightsOutputDescriptor instead')
const DescribeContributorInsightsOutput$json = {
  '1': 'DescribeContributorInsightsOutput',
  '2': [
    {
      '1': 'contributorinsightsmode',
      '3': 86700161,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsMode',
      '10': 'contributorinsightsmode'
    },
    {
      '1': 'contributorinsightsrulelist',
      '3': 140475206,
      '4': 3,
      '5': 9,
      '10': 'contributorinsightsrulelist'
    },
    {
      '1': 'contributorinsightsstatus',
      '3': 363347282,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsStatus',
      '10': 'contributorinsightsstatus'
    },
    {
      '1': 'failureexception',
      '3': 284668899,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.FailureException',
      '10': 'failureexception'
    },
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'lastupdatedatetime',
      '3': 452274318,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedatetime'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeContributorInsightsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeContributorInsightsOutputDescriptor = $convert.base64Decode(
    'CiFEZXNjcmliZUNvbnRyaWJ1dG9ySW5zaWdodHNPdXRwdXQSXgoXY29udHJpYnV0b3JpbnNpZ2'
    'h0c21vZGUYgeGrKSABKA4yIS5keW5hbW9kYi5Db250cmlidXRvckluc2lnaHRzTW9kZVIXY29u'
    'dHJpYnV0b3JpbnNpZ2h0c21vZGUSQwobY29udHJpYnV0b3JpbnNpZ2h0c3J1bGVsaXN0GMb2/U'
    'IgAygJUhtjb250cmlidXRvcmluc2lnaHRzcnVsZWxpc3QSZQoZY29udHJpYnV0b3JpbnNpZ2h0'
    'c3N0YXR1cxjS+qCtASABKA4yIy5keW5hbW9kYi5Db250cmlidXRvckluc2lnaHRzU3RhdHVzUh'
    'ljb250cmlidXRvcmluc2lnaHRzc3RhdHVzEkoKEGZhaWx1cmVleGNlcHRpb24Y4+fehwEgASgL'
    'MhouZHluYW1vZGIuRmFpbHVyZUV4Y2VwdGlvblIQZmFpbHVyZWV4Y2VwdGlvbhIfCglpbmRleG'
    '5hbWUYkdXrMCABKAlSCWluZGV4bmFtZRIyChJsYXN0dXBkYXRlZGF0ZXRpbWUYjtHU1wEgASgJ'
    'UhJsYXN0dXBkYXRlZGF0ZXRpbWUSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use describeEndpointsRequestDescriptor instead')
const DescribeEndpointsRequest$json = {
  '1': 'DescribeEndpointsRequest',
};

/// Descriptor for `DescribeEndpointsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEndpointsRequestDescriptor =
    $convert.base64Decode('ChhEZXNjcmliZUVuZHBvaW50c1JlcXVlc3Q=');

@$core.Deprecated('Use describeEndpointsResponseDescriptor instead')
const DescribeEndpointsResponse$json = {
  '1': 'DescribeEndpointsResponse',
  '2': [
    {
      '1': 'endpoints',
      '3': 16210494,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Endpoint',
      '10': 'endpoints'
    },
  ],
};

/// Descriptor for `DescribeEndpointsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEndpointsResponseDescriptor =
    $convert.base64Decode(
        'ChlEZXNjcmliZUVuZHBvaW50c1Jlc3BvbnNlEjMKCWVuZHBvaW50cxi+tN0HIAMoCzISLmR5bm'
        'Ftb2RiLkVuZHBvaW50UgllbmRwb2ludHM=');

@$core.Deprecated('Use describeExportInputDescriptor instead')
const DescribeExportInput$json = {
  '1': 'DescribeExportInput',
  '2': [
    {'1': 'exportarn', '3': 3661287, '4': 1, '5': 9, '10': 'exportarn'},
  ],
};

/// Descriptor for `DescribeExportInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeExportInputDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZUV4cG9ydElucHV0Eh8KCWV4cG9ydGFybhjnu98BIAEoCVIJZXhwb3J0YXJu');

@$core.Deprecated('Use describeExportOutputDescriptor instead')
const DescribeExportOutput$json = {
  '1': 'DescribeExportOutput',
  '2': [
    {
      '1': 'exportdescription',
      '3': 14399944,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ExportDescription',
      '10': 'exportdescription'
    },
  ],
};

/// Descriptor for `DescribeExportOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeExportOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZUV4cG9ydE91dHB1dBJMChFleHBvcnRkZXNjcmlwdGlvbhjI8+4GIAEoCzIbLm'
    'R5bmFtb2RiLkV4cG9ydERlc2NyaXB0aW9uUhFleHBvcnRkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use describeGlobalTableInputDescriptor instead')
const DescribeGlobalTableInput$json = {
  '1': 'DescribeGlobalTableInput',
  '2': [
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
  ],
};

/// Descriptor for `DescribeGlobalTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeGlobalTableInputDescriptor =
    $convert.base64Decode(
        'ChhEZXNjcmliZUdsb2JhbFRhYmxlSW5wdXQSLAoPZ2xvYmFsdGFibGVuYW1lGKqmp4cBIAEoCV'
        'IPZ2xvYmFsdGFibGVuYW1l');

@$core.Deprecated('Use describeGlobalTableOutputDescriptor instead')
const DescribeGlobalTableOutput$json = {
  '1': 'DescribeGlobalTableOutput',
  '2': [
    {
      '1': 'globaltabledescription',
      '3': 342116405,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.GlobalTableDescription',
      '10': 'globaltabledescription'
    },
  ],
};

/// Descriptor for `DescribeGlobalTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeGlobalTableOutputDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZUdsb2JhbFRhYmxlT3V0cHV0ElwKFmdsb2JhbHRhYmxlZGVzY3JpcHRpb24YtZ'
    'CRowEgASgLMiAuZHluYW1vZGIuR2xvYmFsVGFibGVEZXNjcmlwdGlvblIWZ2xvYmFsdGFibGVk'
    'ZXNjcmlwdGlvbg==');

@$core.Deprecated('Use describeGlobalTableSettingsInputDescriptor instead')
const DescribeGlobalTableSettingsInput$json = {
  '1': 'DescribeGlobalTableSettingsInput',
  '2': [
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
  ],
};

/// Descriptor for `DescribeGlobalTableSettingsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeGlobalTableSettingsInputDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZUdsb2JhbFRhYmxlU2V0dGluZ3NJbnB1dBIsCg9nbG9iYWx0YWJsZW5hbWUYqq'
        'anhwEgASgJUg9nbG9iYWx0YWJsZW5hbWU=');

@$core.Deprecated('Use describeGlobalTableSettingsOutputDescriptor instead')
const DescribeGlobalTableSettingsOutput$json = {
  '1': 'DescribeGlobalTableSettingsOutput',
  '2': [
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
    {
      '1': 'replicasettings',
      '3': 288882107,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaSettingsDescription',
      '10': 'replicasettings'
    },
  ],
};

/// Descriptor for `DescribeGlobalTableSettingsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeGlobalTableSettingsOutputDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZUdsb2JhbFRhYmxlU2V0dGluZ3NPdXRwdXQSLAoPZ2xvYmFsdGFibGVuYW1lGK'
        'qmp4cBIAEoCVIPZ2xvYmFsdGFibGVuYW1lElIKD3JlcGxpY2FzZXR0aW5ncxi7+9+JASADKAsy'
        'JC5keW5hbW9kYi5SZXBsaWNhU2V0dGluZ3NEZXNjcmlwdGlvblIPcmVwbGljYXNldHRpbmdz');

@$core.Deprecated('Use describeImportInputDescriptor instead')
const DescribeImportInput$json = {
  '1': 'DescribeImportInput',
  '2': [
    {'1': 'importarn', '3': 444379628, '4': 1, '5': 9, '10': 'importarn'},
  ],
};

/// Descriptor for `DescribeImportInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeImportInputDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZUltcG9ydElucHV0EiAKCWltcG9ydGFybhjs4/LTASABKAlSCWltcG9ydGFybg'
    '==');

@$core.Deprecated('Use describeImportOutputDescriptor instead')
const DescribeImportOutput$json = {
  '1': 'DescribeImportOutput',
  '2': [
    {
      '1': 'importtabledescription',
      '3': 407746467,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ImportTableDescription',
      '10': 'importtabledescription'
    },
  ],
};

/// Descriptor for `DescribeImportOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeImportOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZUltcG9ydE91dHB1dBJcChZpbXBvcnR0YWJsZWRlc2NyaXB0aW9uGKPvtsIBIA'
    'EoCzIgLmR5bmFtb2RiLkltcG9ydFRhYmxlRGVzY3JpcHRpb25SFmltcG9ydHRhYmxlZGVzY3Jp'
    'cHRpb24=');

@$core.Deprecated(
    'Use describeKinesisStreamingDestinationInputDescriptor instead')
const DescribeKinesisStreamingDestinationInput$json = {
  '1': 'DescribeKinesisStreamingDestinationInput',
  '2': [
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeKinesisStreamingDestinationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeKinesisStreamingDestinationInputDescriptor =
    $convert.base64Decode(
        'CihEZXNjcmliZUtpbmVzaXNTdHJlYW1pbmdEZXN0aW5hdGlvbklucHV0EiAKCXRhYmxlbmFtZR'
        'jd5NqBASABKAlSCXRhYmxlbmFtZQ==');

@$core.Deprecated(
    'Use describeKinesisStreamingDestinationOutputDescriptor instead')
const DescribeKinesisStreamingDestinationOutput$json = {
  '1': 'DescribeKinesisStreamingDestinationOutput',
  '2': [
    {
      '1': 'kinesisdatastreamdestinations',
      '3': 4645129,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KinesisDataStreamDestination',
      '10': 'kinesisdatastreamdestinations'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeKinesisStreamingDestinationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeKinesisStreamingDestinationOutputDescriptor = $convert.base64Decode(
        'CilEZXNjcmliZUtpbmVzaXNTdHJlYW1pbmdEZXN0aW5hdGlvbk91dHB1dBJvCh1raW5lc2lzZG'
        'F0YXN0cmVhbWRlc3RpbmF0aW9ucxiJwpsCIAMoCzImLmR5bmFtb2RiLktpbmVzaXNEYXRhU3Ry'
        'ZWFtRGVzdGluYXRpb25SHWtpbmVzaXNkYXRhc3RyZWFtZGVzdGluYXRpb25zEiAKCXRhYmxlbm'
        'FtZRjd5NqBASABKAlSCXRhYmxlbmFtZQ==');

@$core.Deprecated('Use describeLimitsInputDescriptor instead')
const DescribeLimitsInput$json = {
  '1': 'DescribeLimitsInput',
};

/// Descriptor for `DescribeLimitsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLimitsInputDescriptor =
    $convert.base64Decode('ChNEZXNjcmliZUxpbWl0c0lucHV0');

@$core.Deprecated('Use describeLimitsOutputDescriptor instead')
const DescribeLimitsOutput$json = {
  '1': 'DescribeLimitsOutput',
  '2': [
    {
      '1': 'accountmaxreadcapacityunits',
      '3': 318006242,
      '4': 1,
      '5': 3,
      '10': 'accountmaxreadcapacityunits'
    },
    {
      '1': 'accountmaxwritecapacityunits',
      '3': 526094897,
      '4': 1,
      '5': 3,
      '10': 'accountmaxwritecapacityunits'
    },
    {
      '1': 'tablemaxreadcapacityunits',
      '3': 165993307,
      '4': 1,
      '5': 3,
      '10': 'tablemaxreadcapacityunits'
    },
    {
      '1': 'tablemaxwritecapacityunits',
      '3': 444127818,
      '4': 1,
      '5': 3,
      '10': 'tablemaxwritecapacityunits'
    },
  ],
};

/// Descriptor for `DescribeLimitsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLimitsOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZUxpbWl0c091dHB1dBJEChthY2NvdW50bWF4cmVhZGNhcGFjaXR5dW5pdHMY4s'
    'fRlwEgASgDUhthY2NvdW50bWF4cmVhZGNhcGFjaXR5dW5pdHMSRgocYWNjb3VudG1heHdyaXRl'
    'Y2FwYWNpdHl1bml0cxixpO76ASABKANSHGFjY291bnRtYXh3cml0ZWNhcGFjaXR5dW5pdHMSPw'
    'oZdGFibGVtYXhyZWFkY2FwYWNpdHl1bml0cxjbtpNPIAEoA1IZdGFibGVtYXhyZWFkY2FwYWNp'
    'dHl1bml0cxJCChp0YWJsZW1heHdyaXRlY2FwYWNpdHl1bml0cxjKtOPTASABKANSGnRhYmxlbW'
    'F4d3JpdGVjYXBhY2l0eXVuaXRz');

@$core.Deprecated('Use describeTableInputDescriptor instead')
const DescribeTableInput$json = {
  '1': 'DescribeTableInput',
  '2': [
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTableInputDescriptor = $convert.base64Decode(
    'ChJEZXNjcmliZVRhYmxlSW5wdXQSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use describeTableOutputDescriptor instead')
const DescribeTableOutput$json = {
  '1': 'DescribeTableOutput',
  '2': [
    {
      '1': 'table',
      '3': 386722688,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableDescription',
      '10': 'table'
    },
  ],
};

/// Descriptor for `DescribeTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTableOutputDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZVRhYmxlT3V0cHV0EjQKBXRhYmxlGIDXs7gBIAEoCzIaLmR5bmFtb2RiLlRhYm'
    'xlRGVzY3JpcHRpb25SBXRhYmxl');

@$core.Deprecated('Use describeTableReplicaAutoScalingInputDescriptor instead')
const DescribeTableReplicaAutoScalingInput$json = {
  '1': 'DescribeTableReplicaAutoScalingInput',
  '2': [
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeTableReplicaAutoScalingInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTableReplicaAutoScalingInputDescriptor =
    $convert.base64Decode(
        'CiREZXNjcmliZVRhYmxlUmVwbGljYUF1dG9TY2FsaW5nSW5wdXQSIAoJdGFibGVuYW1lGN3k2o'
        'EBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use describeTableReplicaAutoScalingOutputDescriptor instead')
const DescribeTableReplicaAutoScalingOutput$json = {
  '1': 'DescribeTableReplicaAutoScalingOutput',
  '2': [
    {
      '1': 'tableautoscalingdescription',
      '3': 490422806,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableAutoScalingDescription',
      '10': 'tableautoscalingdescription'
    },
  ],
};

/// Descriptor for `DescribeTableReplicaAutoScalingOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTableReplicaAutoScalingOutputDescriptor =
    $convert.base64Decode(
        'CiVEZXNjcmliZVRhYmxlUmVwbGljYUF1dG9TY2FsaW5nT3V0cHV0EmsKG3RhYmxlYXV0b3NjYW'
        'xpbmdkZXNjcmlwdGlvbhiWhO3pASABKAsyJS5keW5hbW9kYi5UYWJsZUF1dG9TY2FsaW5nRGVz'
        'Y3JpcHRpb25SG3RhYmxlYXV0b3NjYWxpbmdkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use describeTimeToLiveInputDescriptor instead')
const DescribeTimeToLiveInput$json = {
  '1': 'DescribeTimeToLiveInput',
  '2': [
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeTimeToLiveInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTimeToLiveInputDescriptor =
    $convert.base64Decode(
        'ChdEZXNjcmliZVRpbWVUb0xpdmVJbnB1dBIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW'
        '5hbWU=');

@$core.Deprecated('Use describeTimeToLiveOutputDescriptor instead')
const DescribeTimeToLiveOutput$json = {
  '1': 'DescribeTimeToLiveOutput',
  '2': [
    {
      '1': 'timetolivedescription',
      '3': 367590876,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TimeToLiveDescription',
      '10': 'timetolivedescription'
    },
  ],
};

/// Descriptor for `DescribeTimeToLiveOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTimeToLiveOutputDescriptor = $convert.base64Decode(
    'ChhEZXNjcmliZVRpbWVUb0xpdmVPdXRwdXQSWQoVdGltZXRvbGl2ZWRlc2NyaXB0aW9uGNz7o6'
    '8BIAEoCzIfLmR5bmFtb2RiLlRpbWVUb0xpdmVEZXNjcmlwdGlvblIVdGltZXRvbGl2ZWRlc2Ny'
    'aXB0aW9u');

@$core.Deprecated('Use duplicateItemExceptionDescriptor instead')
const DuplicateItemException$json = {
  '1': 'DuplicateItemException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DuplicateItemException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List duplicateItemExceptionDescriptor =
    $convert.base64Decode(
        'ChZEdXBsaWNhdGVJdGVtRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use enableKinesisStreamingConfigurationDescriptor instead')
const EnableKinesisStreamingConfiguration$json = {
  '1': 'EnableKinesisStreamingConfiguration',
  '2': [
    {
      '1': 'approximatecreationdatetimeprecision',
      '3': 392293352,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ApproximateCreationDateTimePrecision',
      '10': 'approximatecreationdatetimeprecision'
    },
  ],
};

/// Descriptor for `EnableKinesisStreamingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableKinesisStreamingConfigurationDescriptor =
    $convert.base64Decode(
        'CiNFbmFibGVLaW5lc2lzU3RyZWFtaW5nQ29uZmlndXJhdGlvbhKGAQokYXBwcm94aW1hdGVjcm'
        'VhdGlvbmRhdGV0aW1lcHJlY2lzaW9uGOjXh7sBIAEoDjIuLmR5bmFtb2RiLkFwcHJveGltYXRl'
        'Q3JlYXRpb25EYXRlVGltZVByZWNpc2lvblIkYXBwcm94aW1hdGVjcmVhdGlvbmRhdGV0aW1lcH'
        'JlY2lzaW9u');

@$core.Deprecated('Use endpointDescriptor instead')
const Endpoint$json = {
  '1': 'Endpoint',
  '2': [
    {'1': 'address', '3': 268787956, '4': 1, '5': 9, '10': 'address'},
    {
      '1': 'cacheperiodinminutes',
      '3': 71526449,
      '4': 1,
      '5': 3,
      '10': 'cacheperiodinminutes'
    },
  ],
};

/// Descriptor for `Endpoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endpointDescriptor = $convert.base64Decode(
    'CghFbmRwb2ludBIcCgdhZGRyZXNzGPTBlYABIAEoCVIHYWRkcmVzcxI1ChRjYWNoZXBlcmlvZG'
    'lubWludXRlcxix0I0iIAEoA1IUY2FjaGVwZXJpb2Rpbm1pbnV0ZXM=');

@$core.Deprecated('Use executeStatementInputDescriptor instead')
const ExecuteStatementInput$json = {
  '1': 'ExecuteStatementInput',
  '2': [
    {
      '1': 'consistentread',
      '3': 531556994,
      '4': 1,
      '5': 8,
      '10': 'consistentread'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'parameters'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'statement', '3': 248790199, '4': 1, '5': 9, '10': 'statement'},
  ],
};

/// Descriptor for `ExecuteStatementInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executeStatementInputDescriptor = $convert.base64Decode(
    'ChVFeGVjdXRlU3RhdGVtZW50SW5wdXQSKgoOY29uc2lzdGVudHJlYWQYgtW7/QEgASgIUg5jb2'
    '5zaXN0ZW50cmVhZBIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0Eh8KCW5leHR0b2tlbhj+hLpn'
    'IAEoCVIJbmV4dHRva2VuEjwKCnBhcmFtZXRlcnMY+qf+6wEgAygLMhguZHluYW1vZGIuQXR0cm'
    'lidXRlVmFsdWVSCnBhcmFtZXRlcnMSWwoWcmV0dXJuY29uc3VtZWRjYXBhY2l0eRj+5+EUIAEo'
    'DjIgLmR5bmFtb2RiLlJldHVybkNvbnN1bWVkQ2FwYWNpdHlSFnJldHVybmNvbnN1bWVkY2FwYW'
    'NpdHkSggEKI3JldHVybnZhbHVlc29uY29uZGl0aW9uY2hlY2tmYWlsdXJlGOCXgQIgASgOMi0u'
    'ZHluYW1vZGIuUmV0dXJuVmFsdWVzT25Db25kaXRpb25DaGVja0ZhaWx1cmVSI3JldHVybnZhbH'
    'Vlc29uY29uZGl0aW9uY2hlY2tmYWlsdXJlEh8KCXN0YXRlbWVudBi3+dB2IAEoCVIJc3RhdGVt'
    'ZW50');

@$core.Deprecated('Use executeStatementOutputDescriptor instead')
const ExecuteStatementOutput$json = {
  '1': 'ExecuteStatementOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ExecuteStatementOutput.ItemsEntry',
      '10': 'items'
    },
    {
      '1': 'lastevaluatedkey',
      '3': 54319830,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ExecuteStatementOutput.LastevaluatedkeyEntry',
      '10': 'lastevaluatedkey'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
  '3': [
    ExecuteStatementOutput_ItemsEntry$json,
    ExecuteStatementOutput_LastevaluatedkeyEntry$json
  ],
};

@$core.Deprecated('Use executeStatementOutputDescriptor instead')
const ExecuteStatementOutput_ItemsEntry$json = {
  '1': 'ItemsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use executeStatementOutputDescriptor instead')
const ExecuteStatementOutput_LastevaluatedkeyEntry$json = {
  '1': 'LastevaluatedkeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ExecuteStatementOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executeStatementOutputDescriptor = $convert.base64Decode(
    'ChZFeGVjdXRlU3RhdGVtZW50T3V0cHV0EkoKEGNvbnN1bWVkY2FwYWNpdHkYrKqh1gEgASgLMh'
    'ouZHluYW1vZGIuQ29uc3VtZWRDYXBhY2l0eVIQY29uc3VtZWRjYXBhY2l0eRJECgVpdGVtcxiw'
    '8NgBIAMoCzIrLmR5bmFtb2RiLkV4ZWN1dGVTdGF0ZW1lbnRPdXRwdXQuSXRlbXNFbnRyeVIFaX'
    'RlbXMSZQoQbGFzdGV2YWx1YXRlZGtleRjWtfMZIAMoCzI2LmR5bmFtb2RiLkV4ZWN1dGVTdGF0'
    'ZW1lbnRPdXRwdXQuTGFzdGV2YWx1YXRlZGtleUVudHJ5UhBsYXN0ZXZhbHVhdGVka2V5Eh8KCW'
    '5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuGlIKCkl0ZW1zRW50cnkSEAoDa2V5GAEgASgJ'
    'UgNrZXkSLgoFdmFsdWUYAiABKAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdWU6Aj'
    'gBGl0KFUxhc3RldmFsdWF0ZWRrZXlFbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1ZRgC'
    'IAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use executeTransactionInputDescriptor instead')
const ExecuteTransactionInput$json = {
  '1': 'ExecuteTransactionInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'transactstatements',
      '3': 446038782,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ParameterizedStatement',
      '10': 'transactstatements'
    },
  ],
};

/// Descriptor for `ExecuteTransactionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executeTransactionInputDescriptor = $convert.base64Decode(
    'ChdFeGVjdXRlVHJhbnNhY3Rpb25JbnB1dBIyChJjbGllbnRyZXF1ZXN0dG9rZW4Y8e+i2QEgAS'
    'gJUhJjbGllbnRyZXF1ZXN0dG9rZW4SWwoWcmV0dXJuY29uc3VtZWRjYXBhY2l0eRj+5+EUIAEo'
    'DjIgLmR5bmFtb2RiLlJldHVybkNvbnN1bWVkQ2FwYWNpdHlSFnJldHVybmNvbnN1bWVkY2FwYW'
    'NpdHkSVAoSdHJhbnNhY3RzdGF0ZW1lbnRzGP6F2NQBIAMoCzIgLmR5bmFtb2RiLlBhcmFtZXRl'
    'cml6ZWRTdGF0ZW1lbnRSEnRyYW5zYWN0c3RhdGVtZW50cw==');

@$core.Deprecated('Use executeTransactionOutputDescriptor instead')
const ExecuteTransactionOutput$json = {
  '1': 'ExecuteTransactionOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'responses',
      '3': 114852856,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ItemResponse',
      '10': 'responses'
    },
  ],
};

/// Descriptor for `ExecuteTransactionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executeTransactionOutputDescriptor = $convert.base64Decode(
    'ChhFeGVjdXRlVHJhbnNhY3Rpb25PdXRwdXQSSgoQY29uc3VtZWRjYXBhY2l0eRisqqHWASADKA'
    'syGi5keW5hbW9kYi5Db25zdW1lZENhcGFjaXR5UhBjb25zdW1lZGNhcGFjaXR5EjcKCXJlc3Bv'
    'bnNlcxj4h+I2IAMoCzIWLmR5bmFtb2RiLkl0ZW1SZXNwb25zZVIJcmVzcG9uc2Vz');

@$core.Deprecated('Use expectedAttributeValueDescriptor instead')
const ExpectedAttributeValue$json = {
  '1': 'ExpectedAttributeValue',
  '2': [
    {
      '1': 'attributevaluelist',
      '3': 157424013,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'attributevaluelist'
    },
    {
      '1': 'comparisonoperator',
      '3': 7059861,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ComparisonOperator',
      '10': 'comparisonoperator'
    },
    {'1': 'exists', '3': 265084382, '4': 1, '5': 8, '10': 'exists'},
    {
      '1': 'value',
      '3': 289929579,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
};

/// Descriptor for `ExpectedAttributeValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expectedAttributeValueDescriptor = $convert.base64Decode(
    'ChZFeHBlY3RlZEF0dHJpYnV0ZVZhbHVlEksKEmF0dHJpYnV0ZXZhbHVlbGlzdBiNs4hLIAMoCz'
    'IYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUhJhdHRyaWJ1dGV2YWx1ZWxpc3QSTwoSY29tcGFy'
    'aXNvbm9wZXJhdG9yGJXzrgMgASgOMhwuZHluYW1vZGIuQ29tcGFyaXNvbk9wZXJhdG9yUhJjb2'
    '1wYXJpc29ub3BlcmF0b3ISGQoGZXhpc3RzGN67s34gASgIUgZleGlzdHMSMgoFdmFsdWUY6/Kf'
    'igEgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVl');

@$core.Deprecated('Use exportConflictExceptionDescriptor instead')
const ExportConflictException$json = {
  '1': 'ExportConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExportConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportConflictExceptionDescriptor =
    $convert.base64Decode(
        'ChdFeHBvcnRDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use exportDescriptionDescriptor instead')
const ExportDescription$json = {
  '1': 'ExportDescription',
  '2': [
    {
      '1': 'billedsizebytes',
      '3': 443315368,
      '4': 1,
      '5': 3,
      '10': 'billedsizebytes'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'exportarn', '3': 3661287, '4': 1, '5': 9, '10': 'exportarn'},
    {
      '1': 'exportformat',
      '3': 300799837,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportFormat',
      '10': 'exportformat'
    },
    {
      '1': 'exportmanifest',
      '3': 196735863,
      '4': 1,
      '5': 9,
      '10': 'exportmanifest'
    },
    {
      '1': 'exportstatus',
      '3': 459702918,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportStatus',
      '10': 'exportstatus'
    },
    {'1': 'exporttime', '3': 28335083, '4': 1, '5': 9, '10': 'exporttime'},
    {
      '1': 'exporttype',
      '3': 189377484,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportType',
      '10': 'exporttype'
    },
    {'1': 'failurecode', '3': 84707897, '4': 1, '5': 9, '10': 'failurecode'},
    {
      '1': 'failuremessage',
      '3': 353556937,
      '4': 1,
      '5': 9,
      '10': 'failuremessage'
    },
    {
      '1': 'incrementalexportspecification',
      '3': 20792295,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.IncrementalExportSpecification',
      '10': 'incrementalexportspecification'
    },
    {'1': 'itemcount', '3': 26280022, '4': 1, '5': 3, '10': 'itemcount'},
    {'1': 's3bucket', '3': 114031434, '4': 1, '5': 9, '10': 's3bucket'},
    {
      '1': 's3bucketowner',
      '3': 351576129,
      '4': 1,
      '5': 9,
      '10': 's3bucketowner'
    },
    {'1': 's3prefix', '3': 21529336, '4': 1, '5': 9, '10': 's3prefix'},
    {
      '1': 's3ssealgorithm',
      '3': 213032706,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.S3SseAlgorithm',
      '10': 's3ssealgorithm'
    },
    {'1': 's3ssekmskeyid', '3': 9386110, '4': 1, '5': 9, '10': 's3ssekmskeyid'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {'1': 'tableid', '3': 449893011, '4': 1, '5': 9, '10': 'tableid'},
  ],
};

/// Descriptor for `ExportDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportDescriptionDescriptor = $convert.base64Decode(
    'ChFFeHBvcnREZXNjcmlwdGlvbhIsCg9iaWxsZWRzaXplYnl0ZXMYqOmx0wEgASgDUg9iaWxsZW'
    'RzaXplYnl0ZXMSIwoLY2xpZW50dG9rZW4YzPu7QSABKAlSC2NsaWVudHRva2VuEhsKB2VuZHRp'
    'bWUYzO+8HiABKAlSB2VuZHRpbWUSHwoJZXhwb3J0YXJuGOe73wEgASgJUglleHBvcnRhcm4SPg'
    'oMZXhwb3J0Zm9ybWF0GN2ut48BIAEoDjIWLmR5bmFtb2RiLkV4cG9ydEZvcm1hdFIMZXhwb3J0'
    'Zm9ybWF0EikKDmV4cG9ydG1hbmlmZXN0GPfm510gASgJUg5leHBvcnRtYW5pZmVzdBI+CgxleH'
    'BvcnRzdGF0dXMYhoWa2wEgASgOMhYuZHluYW1vZGIuRXhwb3J0U3RhdHVzUgxleHBvcnRzdGF0'
    'dXMSIQoKZXhwb3J0dGltZRjrt8ENIAEoCVIKZXhwb3J0dGltZRI3CgpleHBvcnR0eXBlGMzXpl'
    'ogASgOMhQuZHluYW1vZGIuRXhwb3J0VHlwZVIKZXhwb3J0dHlwZRIjCgtmYWlsdXJlY29kZRi5'
    'lLIoIAEoCVILZmFpbHVyZWNvZGUSKgoOZmFpbHVyZW1lc3NhZ2UYybPLqAEgASgJUg5mYWlsdX'
    'JlbWVzc2FnZRJzCh5pbmNyZW1lbnRhbGV4cG9ydHNwZWNpZmljYXRpb24Y54f1CSABKAsyKC5k'
    'eW5hbW9kYi5JbmNyZW1lbnRhbEV4cG9ydFNwZWNpZmljYXRpb25SHmluY3JlbWVudGFsZXhwb3'
    'J0c3BlY2lmaWNhdGlvbhIfCglpdGVtY291bnQY1oDEDCABKANSCWl0ZW1jb3VudBIdCghzM2J1'
    'Y2tldBjK9q82IAEoCVIIczNidWNrZXQSKAoNczNidWNrZXRvd25lchjBwNKnASABKAlSDXMzYn'
    'Vja2V0b3duZXISHQoIczNwcmVmaXgY+IWiCiABKAlSCHMzcHJlZml4EkMKDnMzc3NlYWxnb3Jp'
    'dGhtGIK+ymUgASgOMhguZHluYW1vZGIuUzNTc2VBbGdvcml0aG1SDnMzc3NlYWxnb3JpdGhtEi'
    'cKDXMzc3Nla21za2V5aWQY/vC8BCABKAlSDXMzc3Nla21za2V5aWQSIAoJc3RhcnR0aW1lGO+0'
    '5bABIAEoCVIJc3RhcnR0aW1lEh4KCHRhYmxlYXJuGOOA680BIAEoCVIIdGFibGVhcm4SHAoHdG'
    'FibGVpZBiTpcPWASABKAlSB3RhYmxlaWQ=');

@$core.Deprecated('Use exportNotFoundExceptionDescriptor instead')
const ExportNotFoundException$json = {
  '1': 'ExportNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExportNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChdFeHBvcnROb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use exportSummaryDescriptor instead')
const ExportSummary$json = {
  '1': 'ExportSummary',
  '2': [
    {'1': 'exportarn', '3': 3661287, '4': 1, '5': 9, '10': 'exportarn'},
    {
      '1': 'exportstatus',
      '3': 459702918,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportStatus',
      '10': 'exportstatus'
    },
    {
      '1': 'exporttype',
      '3': 189377484,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportType',
      '10': 'exporttype'
    },
  ],
};

/// Descriptor for `ExportSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportSummaryDescriptor = $convert.base64Decode(
    'Cg1FeHBvcnRTdW1tYXJ5Eh8KCWV4cG9ydGFybhjnu98BIAEoCVIJZXhwb3J0YXJuEj4KDGV4cG'
    '9ydHN0YXR1cxiGhZrbASABKA4yFi5keW5hbW9kYi5FeHBvcnRTdGF0dXNSDGV4cG9ydHN0YXR1'
    'cxI3CgpleHBvcnR0eXBlGMzXplogASgOMhQuZHluYW1vZGIuRXhwb3J0VHlwZVIKZXhwb3J0dH'
    'lwZQ==');

@$core.Deprecated('Use exportTableToPointInTimeInputDescriptor instead')
const ExportTableToPointInTimeInput$json = {
  '1': 'ExportTableToPointInTimeInput',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'exportformat',
      '3': 300799837,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportFormat',
      '10': 'exportformat'
    },
    {'1': 'exporttime', '3': 28335083, '4': 1, '5': 9, '10': 'exporttime'},
    {
      '1': 'exporttype',
      '3': 189377484,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportType',
      '10': 'exporttype'
    },
    {
      '1': 'incrementalexportspecification',
      '3': 20792295,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.IncrementalExportSpecification',
      '10': 'incrementalexportspecification'
    },
    {'1': 's3bucket', '3': 114031434, '4': 1, '5': 9, '10': 's3bucket'},
    {
      '1': 's3bucketowner',
      '3': 351576129,
      '4': 1,
      '5': 9,
      '10': 's3bucketowner'
    },
    {'1': 's3prefix', '3': 21529336, '4': 1, '5': 9, '10': 's3prefix'},
    {
      '1': 's3ssealgorithm',
      '3': 213032706,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.S3SseAlgorithm',
      '10': 's3ssealgorithm'
    },
    {'1': 's3ssekmskeyid', '3': 9386110, '4': 1, '5': 9, '10': 's3ssekmskeyid'},
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
  ],
};

/// Descriptor for `ExportTableToPointInTimeInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportTableToPointInTimeInputDescriptor = $convert.base64Decode(
    'Ch1FeHBvcnRUYWJsZVRvUG9pbnRJblRpbWVJbnB1dBIjCgtjbGllbnR0b2tlbhjM+7tBIAEoCV'
    'ILY2xpZW50dG9rZW4SPgoMZXhwb3J0Zm9ybWF0GN2ut48BIAEoDjIWLmR5bmFtb2RiLkV4cG9y'
    'dEZvcm1hdFIMZXhwb3J0Zm9ybWF0EiEKCmV4cG9ydHRpbWUY67fBDSABKAlSCmV4cG9ydHRpbW'
    'USNwoKZXhwb3J0dHlwZRjM16ZaIAEoDjIULmR5bmFtb2RiLkV4cG9ydFR5cGVSCmV4cG9ydHR5'
    'cGUScwoeaW5jcmVtZW50YWxleHBvcnRzcGVjaWZpY2F0aW9uGOeH9QkgASgLMiguZHluYW1vZG'
    'IuSW5jcmVtZW50YWxFeHBvcnRTcGVjaWZpY2F0aW9uUh5pbmNyZW1lbnRhbGV4cG9ydHNwZWNp'
    'ZmljYXRpb24SHQoIczNidWNrZXQYyvavNiABKAlSCHMzYnVja2V0EigKDXMzYnVja2V0b3duZX'
    'IYwcDSpwEgASgJUg1zM2J1Y2tldG93bmVyEh0KCHMzcHJlZml4GPiFogogASgJUghzM3ByZWZp'
    'eBJDCg5zM3NzZWFsZ29yaXRobRiCvsplIAEoDjIYLmR5bmFtb2RiLlMzU3NlQWxnb3JpdGhtUg'
    '5zM3NzZWFsZ29yaXRobRInCg1zM3NzZWttc2tleWlkGP7wvAQgASgJUg1zM3NzZWttc2tleWlk'
    'Eh4KCHRhYmxlYXJuGOOA680BIAEoCVIIdGFibGVhcm4=');

@$core.Deprecated('Use exportTableToPointInTimeOutputDescriptor instead')
const ExportTableToPointInTimeOutput$json = {
  '1': 'ExportTableToPointInTimeOutput',
  '2': [
    {
      '1': 'exportdescription',
      '3': 14399944,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ExportDescription',
      '10': 'exportdescription'
    },
  ],
};

/// Descriptor for `ExportTableToPointInTimeOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportTableToPointInTimeOutputDescriptor =
    $convert.base64Decode(
        'Ch5FeHBvcnRUYWJsZVRvUG9pbnRJblRpbWVPdXRwdXQSTAoRZXhwb3J0ZGVzY3JpcHRpb24YyP'
        'PuBiABKAsyGy5keW5hbW9kYi5FeHBvcnREZXNjcmlwdGlvblIRZXhwb3J0ZGVzY3JpcHRpb24=');

@$core.Deprecated('Use failureExceptionDescriptor instead')
const FailureException$json = {
  '1': 'FailureException',
  '2': [
    {
      '1': 'exceptiondescription',
      '3': 75755113,
      '4': 1,
      '5': 9,
      '10': 'exceptiondescription'
    },
    {
      '1': 'exceptionname',
      '3': 449933270,
      '4': 1,
      '5': 9,
      '10': 'exceptionname'
    },
  ],
};

/// Descriptor for `FailureException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List failureExceptionDescriptor = $convert.base64Decode(
    'ChBGYWlsdXJlRXhjZXB0aW9uEjUKFGV4Y2VwdGlvbmRlc2NyaXB0aW9uGOncjyQgASgJUhRleG'
    'NlcHRpb25kZXNjcmlwdGlvbhIoCg1leGNlcHRpb25uYW1lGNbfxdYBIAEoCVINZXhjZXB0aW9u'
    'bmFtZQ==');

@$core.Deprecated('Use getDescriptor instead')
const Get$json = {
  '1': 'Get',
  '2': [
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Get.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Get.KeyEntry',
      '10': 'key'
    },
    {
      '1': 'projectionexpression',
      '3': 150730243,
      '4': 1,
      '5': 9,
      '10': 'projectionexpression'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [Get_ExpressionattributenamesEntry$json, Get_KeyEntry$json],
};

@$core.Deprecated('Use getDescriptor instead')
const Get_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use getDescriptor instead')
const Get_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `Get`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDescriptor = $convert.base64Decode(
    'CgNHZXQSagoYZXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMisuZHluYW1vZGIuR2'
    'V0LkV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9uYXR0cmlidXRlbmFt'
    'ZXMSKwoDa2V5GI2S62ggAygLMhYuZHluYW1vZGIuR2V0LktleUVudHJ5UgNrZXkSNQoUcHJvam'
    'VjdGlvbmV4cHJlc3Npb24Yg+zvRyABKAlSFHByb2plY3Rpb25leHByZXNzaW9uEiAKCXRhYmxl'
    'bmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZRpLCh1FeHByZXNzaW9uYXR0cmlidXRlbmFtZXNFbn'
    'RyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGlAKCEtleUVu'
    'dHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidX'
    'RlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getItemInputDescriptor instead')
const GetItemInput$json = {
  '1': 'GetItemInput',
  '2': [
    {
      '1': 'attributestoget',
      '3': 311382592,
      '4': 3,
      '5': 9,
      '10': 'attributestoget'
    },
    {
      '1': 'consistentread',
      '3': 531556994,
      '4': 1,
      '5': 8,
      '10': 'consistentread'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GetItemInput.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GetItemInput.KeyEntry',
      '10': 'key'
    },
    {
      '1': 'projectionexpression',
      '3': 150730243,
      '4': 1,
      '5': 9,
      '10': 'projectionexpression'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [
    GetItemInput_ExpressionattributenamesEntry$json,
    GetItemInput_KeyEntry$json
  ],
};

@$core.Deprecated('Use getItemInputDescriptor instead')
const GetItemInput_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use getItemInputDescriptor instead')
const GetItemInput_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `GetItemInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getItemInputDescriptor = $convert.base64Decode(
    'CgxHZXRJdGVtSW5wdXQSLAoPYXR0cmlidXRlc3RvZ2V0GMCkvZQBIAMoCVIPYXR0cmlidXRlc3'
    'RvZ2V0EioKDmNvbnNpc3RlbnRyZWFkGILVu/0BIAEoCFIOY29uc2lzdGVudHJlYWQScwoYZXhw'
    'cmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMjQuZHluYW1vZGIuR2V0SXRlbUlucHV0Lk'
    'V4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9uYXR0cmlidXRlbmFtZXMS'
    'NAoDa2V5GI2S62ggAygLMh8uZHluYW1vZGIuR2V0SXRlbUlucHV0LktleUVudHJ5UgNrZXkSNQ'
    'oUcHJvamVjdGlvbmV4cHJlc3Npb24Yg+zvRyABKAlSFHByb2plY3Rpb25leHByZXNzaW9uElsK'
    'FnJldHVybmNvbnN1bWVkY2FwYWNpdHkY/ufhFCABKA4yIC5keW5hbW9kYi5SZXR1cm5Db25zdW'
    '1lZENhcGFjaXR5UhZyZXR1cm5jb25zdW1lZGNhcGFjaXR5EiAKCXRhYmxlbmFtZRjd5NqBASAB'
    'KAlSCXRhYmxlbmFtZRpLCh1FeHByZXNzaW9uYXR0cmlidXRlbmFtZXNFbnRyeRIQCgNrZXkYAS'
    'ABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGlAKCEtleUVudHJ5EhAKA2tleRgB'
    'IAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbH'
    'VlOgI4AQ==');

@$core.Deprecated('Use getItemOutputDescriptor instead')
const GetItemOutput$json = {
  '1': 'GetItemOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GetItemOutput.ItemEntry',
      '10': 'item'
    },
  ],
  '3': [GetItemOutput_ItemEntry$json],
};

@$core.Deprecated('Use getItemOutputDescriptor instead')
const GetItemOutput_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `GetItemOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getItemOutputDescriptor = $convert.base64Decode(
    'Cg1HZXRJdGVtT3V0cHV0EkoKEGNvbnN1bWVkY2FwYWNpdHkYrKqh1gEgASgLMhouZHluYW1vZG'
    'IuQ29uc3VtZWRDYXBhY2l0eVIQY29uc3VtZWRjYXBhY2l0eRI5CgRpdGVtGIeAkvsBIAMoCzIh'
    'LmR5bmFtb2RiLkdldEl0ZW1PdXRwdXQuSXRlbUVudHJ5UgRpdGVtGlEKCUl0ZW1FbnRyeRIQCg'
    'NrZXkYASABKAlSA2tleRIuCgV2YWx1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVl'
    'UgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getResourcePolicyInputDescriptor instead')
const GetResourcePolicyInput$json = {
  '1': 'GetResourcePolicyInput',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetResourcePolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyInputDescriptor =
    $convert.base64Decode(
        'ChZHZXRSZXNvdXJjZVBvbGljeUlucHV0EiQKC3Jlc291cmNlYXJuGK342a0BIAEoCVILcmVzb3'
        'VyY2Vhcm4=');

@$core.Deprecated('Use getResourcePolicyOutputDescriptor instead')
const GetResourcePolicyOutput$json = {
  '1': 'GetResourcePolicyOutput',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `GetResourcePolicyOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyOutputDescriptor =
    $convert.base64Decode(
        'ChdHZXRSZXNvdXJjZVBvbGljeU91dHB1dBIaCgZwb2xpY3kYoO/w4AEgASgJUgZwb2xpY3kSIg'
        'oKcmV2aXNpb25pZBiGo57uASABKAlSCnJldmlzaW9uaWQ=');

@$core.Deprecated('Use globalSecondaryIndexDescriptor instead')
const GlobalSecondaryIndex$json = {
  '1': 'GlobalSecondaryIndex',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'projection',
      '3': 105045921,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Projection',
      '10': 'projection'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.WarmThroughput',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `GlobalSecondaryIndex`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalSecondaryIndexDescriptor = $convert.base64Decode(
    'ChRHbG9iYWxTZWNvbmRhcnlJbmRleBIfCglpbmRleG5hbWUYkdXrMCABKAlSCWluZGV4bmFtZR'
    'I8CglrZXlzY2hlbWEY6M/diwEgAygLMhouZHluYW1vZGIuS2V5U2NoZW1hRWxlbWVudFIJa2V5'
    'c2NoZW1hElAKEm9uZGVtYW5kdGhyb3VnaHB1dBiC3trlASABKAsyHC5keW5hbW9kYi5PbkRlbW'
    'FuZFRocm91Z2hwdXRSEm9uZGVtYW5kdGhyb3VnaHB1dBI3Cgpwcm9qZWN0aW9uGKG/izIgASgL'
    'MhQuZHluYW1vZGIuUHJvamVjdGlvblIKcHJvamVjdGlvbhJXChVwcm92aXNpb25lZHRocm91Z2'
    'hwdXQYjKNrIAEoCzIfLmR5bmFtb2RiLlByb3Zpc2lvbmVkVGhyb3VnaHB1dFIVcHJvdmlzaW9u'
    'ZWR0aHJvdWdocHV0EkQKDndhcm10aHJvdWdocHV0GIPeyIoBIAEoCzIYLmR5bmFtb2RiLldhcm'
    '1UaHJvdWdocHV0Ug53YXJtdGhyb3VnaHB1dA==');

@$core.Deprecated('Use globalSecondaryIndexAutoScalingUpdateDescriptor instead')
const GlobalSecondaryIndexAutoScalingUpdate$json = {
  '1': 'GlobalSecondaryIndexAutoScalingUpdate',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'provisionedwritecapacityautoscalingupdate',
      '3': 316618294,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'provisionedwritecapacityautoscalingupdate'
    },
  ],
};

/// Descriptor for `GlobalSecondaryIndexAutoScalingUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalSecondaryIndexAutoScalingUpdateDescriptor =
    $convert.base64Decode(
        'CiVHbG9iYWxTZWNvbmRhcnlJbmRleEF1dG9TY2FsaW5nVXBkYXRlEh8KCWluZGV4bmFtZRiR1e'
        'swIAEoCVIJaW5kZXhuYW1lEoUBCilwcm92aXNpb25lZHdyaXRlY2FwYWNpdHlhdXRvc2NhbGlu'
        'Z3VwZGF0ZRi27PyWASABKAsyIy5keW5hbW9kYi5BdXRvU2NhbGluZ1NldHRpbmdzVXBkYXRlUi'
        'lwcm92aXNpb25lZHdyaXRlY2FwYWNpdHlhdXRvc2NhbGluZ3VwZGF0ZQ==');

@$core.Deprecated('Use globalSecondaryIndexDescriptionDescriptor instead')
const GlobalSecondaryIndexDescription$json = {
  '1': 'GlobalSecondaryIndexDescription',
  '2': [
    {'1': 'backfilling', '3': 251413370, '4': 1, '5': 8, '10': 'backfilling'},
    {'1': 'indexarn', '3': 374335615, '4': 1, '5': 9, '10': 'indexarn'},
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'indexsizebytes',
      '3': 395738346,
      '4': 1,
      '5': 3,
      '10': 'indexsizebytes'
    },
    {
      '1': 'indexstatus',
      '3': 364436830,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.IndexStatus',
      '10': 'indexstatus'
    },
    {'1': 'itemcount', '3': 26280022, '4': 1, '5': 3, '10': 'itemcount'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'projection',
      '3': 105045921,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Projection',
      '10': 'projection'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughputDescription',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndexWarmThroughputDescription',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `GlobalSecondaryIndexDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalSecondaryIndexDescriptionDescriptor = $convert.base64Decode(
    'Ch9HbG9iYWxTZWNvbmRhcnlJbmRleERlc2NyaXB0aW9uEiMKC2JhY2tmaWxsaW5nGPqG8XcgAS'
    'gIUgtiYWNrZmlsbGluZxIeCghpbmRleGFybhj/0L+yASABKAlSCGluZGV4YXJuEh8KCWluZGV4'
    'bmFtZRiR1eswIAEoCVIJaW5kZXhuYW1lEioKDmluZGV4c2l6ZWJ5dGVzGOr52bwBIAEoA1IOaW'
    '5kZXhzaXplYnl0ZXMSOwoLaW5kZXhzdGF0dXMY3rrjrQEgASgOMhUuZHluYW1vZGIuSW5kZXhT'
    'dGF0dXNSC2luZGV4c3RhdHVzEh8KCWl0ZW1jb3VudBjWgMQMIAEoA1IJaXRlbWNvdW50EjwKCW'
    'tleXNjaGVtYRjoz92LASADKAsyGi5keW5hbW9kYi5LZXlTY2hlbWFFbGVtZW50UglrZXlzY2hl'
    'bWESUAoSb25kZW1hbmR0aHJvdWdocHV0GILe2uUBIAEoCzIcLmR5bmFtb2RiLk9uRGVtYW5kVG'
    'hyb3VnaHB1dFISb25kZW1hbmR0aHJvdWdocHV0EjcKCnByb2plY3Rpb24Yob+LMiABKAsyFC5k'
    'eW5hbW9kYi5Qcm9qZWN0aW9uUgpwcm9qZWN0aW9uEmIKFXByb3Zpc2lvbmVkdGhyb3VnaHB1dB'
    'iMo2sgASgLMiouZHluYW1vZGIuUHJvdmlzaW9uZWRUaHJvdWdocHV0RGVzY3JpcHRpb25SFXBy'
    'b3Zpc2lvbmVkdGhyb3VnaHB1dBJjCg53YXJtdGhyb3VnaHB1dBiD3siKASABKAsyNy5keW5hbW'
    '9kYi5HbG9iYWxTZWNvbmRhcnlJbmRleFdhcm1UaHJvdWdocHV0RGVzY3JpcHRpb25SDndhcm10'
    'aHJvdWdocHV0');

@$core.Deprecated('Use globalSecondaryIndexInfoDescriptor instead')
const GlobalSecondaryIndexInfo$json = {
  '1': 'GlobalSecondaryIndexInfo',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'projection',
      '3': 105045921,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Projection',
      '10': 'projection'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
  ],
};

/// Descriptor for `GlobalSecondaryIndexInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalSecondaryIndexInfoDescriptor = $convert.base64Decode(
    'ChhHbG9iYWxTZWNvbmRhcnlJbmRleEluZm8SHwoJaW5kZXhuYW1lGJHV6zAgASgJUglpbmRleG'
    '5hbWUSPAoJa2V5c2NoZW1hGOjP3YsBIAMoCzIaLmR5bmFtb2RiLktleVNjaGVtYUVsZW1lbnRS'
    'CWtleXNjaGVtYRJQChJvbmRlbWFuZHRocm91Z2hwdXQYgt7a5QEgASgLMhwuZHluYW1vZGIuT2'
    '5EZW1hbmRUaHJvdWdocHV0UhJvbmRlbWFuZHRocm91Z2hwdXQSNwoKcHJvamVjdGlvbhihv4sy'
    'IAEoCzIULmR5bmFtb2RiLlByb2plY3Rpb25SCnByb2plY3Rpb24SVwoVcHJvdmlzaW9uZWR0aH'
    'JvdWdocHV0GIyjayABKAsyHy5keW5hbW9kYi5Qcm92aXNpb25lZFRocm91Z2hwdXRSFXByb3Zp'
    'c2lvbmVkdGhyb3VnaHB1dA==');

@$core.Deprecated('Use globalSecondaryIndexUpdateDescriptor instead')
const GlobalSecondaryIndexUpdate$json = {
  '1': 'GlobalSecondaryIndexUpdate',
  '2': [
    {
      '1': 'create',
      '3': 420340862,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.CreateGlobalSecondaryIndexAction',
      '10': 'create'
    },
    {
      '1': 'delete',
      '3': 395831915,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.DeleteGlobalSecondaryIndexAction',
      '10': 'delete'
    },
    {
      '1': 'update',
      '3': 237178517,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.UpdateGlobalSecondaryIndexAction',
      '10': 'update'
    },
  ],
};

/// Descriptor for `GlobalSecondaryIndexUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalSecondaryIndexUpdateDescriptor = $convert.base64Decode(
    'ChpHbG9iYWxTZWNvbmRhcnlJbmRleFVwZGF0ZRJGCgZjcmVhdGUY/si3yAEgASgLMiouZHluYW'
    '1vZGIuQ3JlYXRlR2xvYmFsU2Vjb25kYXJ5SW5kZXhBY3Rpb25SBmNyZWF0ZRJGCgZkZWxldGUY'
    '69TfvAEgASgLMiouZHluYW1vZGIuRGVsZXRlR2xvYmFsU2Vjb25kYXJ5SW5kZXhBY3Rpb25SBm'
    'RlbGV0ZRJFCgZ1cGRhdGUYlZ2McSABKAsyKi5keW5hbW9kYi5VcGRhdGVHbG9iYWxTZWNvbmRh'
    'cnlJbmRleEFjdGlvblIGdXBkYXRl');

@$core.Deprecated(
    'Use globalSecondaryIndexWarmThroughputDescriptionDescriptor instead')
const GlobalSecondaryIndexWarmThroughputDescription$json = {
  '1': 'GlobalSecondaryIndexWarmThroughputDescription',
  '2': [
    {
      '1': 'readunitspersecond',
      '3': 11400732,
      '4': 1,
      '5': 3,
      '10': 'readunitspersecond'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.IndexStatus',
      '10': 'status'
    },
    {
      '1': 'writeunitspersecond',
      '3': 339770127,
      '4': 1,
      '5': 3,
      '10': 'writeunitspersecond'
    },
  ],
};

/// Descriptor for `GlobalSecondaryIndexWarmThroughputDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    globalSecondaryIndexWarmThroughputDescriptionDescriptor =
    $convert.base64Decode(
        'Ci1HbG9iYWxTZWNvbmRhcnlJbmRleFdhcm1UaHJvdWdocHV0RGVzY3JpcHRpb24SMQoScmVhZH'
        'VuaXRzcGVyc2Vjb25kGJzstwUgASgDUhJyZWFkdW5pdHNwZXJzZWNvbmQSMAoGc3RhdHVzGJDk'
        '+wIgASgOMhUuZHluYW1vZGIuSW5kZXhTdGF0dXNSBnN0YXR1cxI0ChN3cml0ZXVuaXRzcGVyc2'
        'Vjb25kGI/2gaIBIAEoA1ITd3JpdGV1bml0c3BlcnNlY29uZA==');

@$core.Deprecated('Use globalTableDescriptor instead')
const GlobalTable$json = {
  '1': 'GlobalTable',
  '2': [
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
    {
      '1': 'replicationgroup',
      '3': 190970785,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Replica',
      '10': 'replicationgroup'
    },
  ],
};

/// Descriptor for `GlobalTable`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalTableDescriptor = $convert.base64Decode(
    'CgtHbG9iYWxUYWJsZRIsCg9nbG9iYWx0YWJsZW5hbWUYqqanhwEgASgJUg9nbG9iYWx0YWJsZW'
    '5hbWUSQAoQcmVwbGljYXRpb25ncm91cBih94dbIAMoCzIRLmR5bmFtb2RiLlJlcGxpY2FSEHJl'
    'cGxpY2F0aW9uZ3JvdXA=');

@$core.Deprecated('Use globalTableAlreadyExistsExceptionDescriptor instead')
const GlobalTableAlreadyExistsException$json = {
  '1': 'GlobalTableAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `GlobalTableAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalTableAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'CiFHbG9iYWxUYWJsZUFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use globalTableDescriptionDescriptor instead')
const GlobalTableDescription$json = {
  '1': 'GlobalTableDescription',
  '2': [
    {
      '1': 'creationdatetime',
      '3': 48904698,
      '4': 1,
      '5': 9,
      '10': 'creationdatetime'
    },
    {
      '1': 'globaltablearn',
      '3': 262302830,
      '4': 1,
      '5': 9,
      '10': 'globaltablearn'
    },
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
    {
      '1': 'globaltablestatus',
      '3': 8833293,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.GlobalTableStatus',
      '10': 'globaltablestatus'
    },
    {
      '1': 'replicationgroup',
      '3': 190970785,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaDescription',
      '10': 'replicationgroup'
    },
  ],
};

/// Descriptor for `GlobalTableDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalTableDescriptionDescriptor = $convert.base64Decode(
    'ChZHbG9iYWxUYWJsZURlc2NyaXB0aW9uEi0KEGNyZWF0aW9uZGF0ZXRpbWUY+vOoFyABKAlSEG'
    'NyZWF0aW9uZGF0ZXRpbWUSKQoOZ2xvYmFsdGFibGVhcm4Y7tiJfSABKAlSDmdsb2JhbHRhYmxl'
    'YXJuEiwKD2dsb2JhbHRhYmxlbmFtZRiqpqeHASABKAlSD2dsb2JhbHRhYmxlbmFtZRJMChFnbG'
    '9iYWx0YWJsZXN0YXR1cxiNkpsEIAEoDjIbLmR5bmFtb2RiLkdsb2JhbFRhYmxlU3RhdHVzUhFn'
    'bG9iYWx0YWJsZXN0YXR1cxJLChByZXBsaWNhdGlvbmdyb3VwGKH3h1sgAygLMhwuZHluYW1vZG'
    'IuUmVwbGljYURlc2NyaXB0aW9uUhByZXBsaWNhdGlvbmdyb3Vw');

@$core.Deprecated(
    'Use globalTableGlobalSecondaryIndexSettingsUpdateDescriptor instead')
const GlobalTableGlobalSecondaryIndexSettingsUpdate$json = {
  '1': 'GlobalTableGlobalSecondaryIndexSettingsUpdate',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'provisionedwritecapacityautoscalingsettingsupdate',
      '3': 302140761,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'provisionedwritecapacityautoscalingsettingsupdate'
    },
    {
      '1': 'provisionedwritecapacityunits',
      '3': 225881684,
      '4': 1,
      '5': 3,
      '10': 'provisionedwritecapacityunits'
    },
  ],
};

/// Descriptor for `GlobalTableGlobalSecondaryIndexSettingsUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    globalTableGlobalSecondaryIndexSettingsUpdateDescriptor =
    $convert.base64Decode(
        'Ci1HbG9iYWxUYWJsZUdsb2JhbFNlY29uZGFyeUluZGV4U2V0dGluZ3NVcGRhdGUSHwoJaW5kZX'
        'huYW1lGJHV6zAgASgJUglpbmRleG5hbWUSlQEKMXByb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eWF1'
        'dG9zY2FsaW5nc2V0dGluZ3N1cGRhdGUY2ZqJkAEgASgLMiMuZHluYW1vZGIuQXV0b1NjYWxpbm'
        'dTZXR0aW5nc1VwZGF0ZVIxcHJvdmlzaW9uZWR3cml0ZWNhcGFjaXR5YXV0b3NjYWxpbmdzZXR0'
        'aW5nc3VwZGF0ZRJHCh1wcm92aXNpb25lZHdyaXRlY2FwYWNpdHl1bml0cxjU3NprIAEoA1IdcH'
        'JvdmlzaW9uZWR3cml0ZWNhcGFjaXR5dW5pdHM=');

@$core.Deprecated('Use globalTableNotFoundExceptionDescriptor instead')
const GlobalTableNotFoundException$json = {
  '1': 'GlobalTableNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `GlobalTableNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalTableNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChxHbG9iYWxUYWJsZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use globalTableWitnessDescriptionDescriptor instead')
const GlobalTableWitnessDescription$json = {
  '1': 'GlobalTableWitnessDescription',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'witnessstatus',
      '3': 303352887,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.WitnessStatus',
      '10': 'witnessstatus'
    },
  ],
};

/// Descriptor for `GlobalTableWitnessDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalTableWitnessDescriptionDescriptor =
    $convert.base64Decode(
        'Ch1HbG9iYWxUYWJsZVdpdG5lc3NEZXNjcmlwdGlvbhIhCgpyZWdpb25uYW1lGL+buTUgASgJUg'
        'pyZWdpb25uYW1lEkEKDXdpdG5lc3NzdGF0dXMYt5jTkAEgASgOMhcuZHluYW1vZGIuV2l0bmVz'
        'c1N0YXR1c1INd2l0bmVzc3N0YXR1cw==');

@$core.Deprecated('Use globalTableWitnessGroupUpdateDescriptor instead')
const GlobalTableWitnessGroupUpdate$json = {
  '1': 'GlobalTableWitnessGroupUpdate',
  '2': [
    {
      '1': 'create',
      '3': 420340862,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.CreateGlobalTableWitnessGroupMemberAction',
      '10': 'create'
    },
    {
      '1': 'delete',
      '3': 395831915,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.DeleteGlobalTableWitnessGroupMemberAction',
      '10': 'delete'
    },
  ],
};

/// Descriptor for `GlobalTableWitnessGroupUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List globalTableWitnessGroupUpdateDescriptor = $convert.base64Decode(
    'Ch1HbG9iYWxUYWJsZVdpdG5lc3NHcm91cFVwZGF0ZRJPCgZjcmVhdGUY/si3yAEgASgLMjMuZH'
    'luYW1vZGIuQ3JlYXRlR2xvYmFsVGFibGVXaXRuZXNzR3JvdXBNZW1iZXJBY3Rpb25SBmNyZWF0'
    'ZRJPCgZkZWxldGUY69TfvAEgASgLMjMuZHluYW1vZGIuRGVsZXRlR2xvYmFsVGFibGVXaXRuZX'
    'NzR3JvdXBNZW1iZXJBY3Rpb25SBmRlbGV0ZQ==');

@$core.Deprecated('Use idempotentParameterMismatchExceptionDescriptor instead')
const IdempotentParameterMismatchException$json = {
  '1': 'IdempotentParameterMismatchException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IdempotentParameterMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List idempotentParameterMismatchExceptionDescriptor =
    $convert.base64Decode(
        'CiRJZGVtcG90ZW50UGFyYW1ldGVyTWlzbWF0Y2hFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use importConflictExceptionDescriptor instead')
const ImportConflictException$json = {
  '1': 'ImportConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ImportConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importConflictExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbXBvcnRDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use importNotFoundExceptionDescriptor instead')
const ImportNotFoundException$json = {
  '1': 'ImportNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ImportNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbXBvcnROb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use importSummaryDescriptor instead')
const ImportSummary$json = {
  '1': 'ImportSummary',
  '2': [
    {
      '1': 'cloudwatchloggrouparn',
      '3': 171042008,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchloggrouparn'
    },
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'importarn', '3': 444379628, '4': 1, '5': 9, '10': 'importarn'},
    {
      '1': 'importstatus',
      '3': 129077631,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'inputformat',
      '3': 405664101,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.InputFormat',
      '10': 'inputformat'
    },
    {
      '1': 's3bucketsource',
      '3': 202310037,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.S3BucketSource',
      '10': 's3bucketsource'
    },
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
  ],
};

/// Descriptor for `ImportSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importSummaryDescriptor = $convert.base64Decode(
    'Cg1JbXBvcnRTdW1tYXJ5EjcKFWNsb3Vkd2F0Y2hsb2dncm91cGFybhjYycdRIAEoCVIVY2xvdW'
    'R3YXRjaGxvZ2dyb3VwYXJuEhsKB2VuZHRpbWUYzO+8HiABKAlSB2VuZHRpbWUSIAoJaW1wb3J0'
    'YXJuGOzj8tMBIAEoCVIJaW1wb3J0YXJuEj0KDGltcG9ydHN0YXR1cxj/osY9IAEoDjIWLmR5bm'
    'Ftb2RiLkltcG9ydFN0YXR1c1IMaW1wb3J0c3RhdHVzEjsKC2lucHV0Zm9ybWF0GOXit8EBIAEo'
    'DjIVLmR5bmFtb2RiLklucHV0Rm9ybWF0UgtpbnB1dGZvcm1hdBJDCg5zM2J1Y2tldHNvdXJjZR'
    'iVg7xgIAEoCzIYLmR5bmFtb2RiLlMzQnVja2V0U291cmNlUg5zM2J1Y2tldHNvdXJjZRIgCglz'
    'dGFydHRpbWUY77TlsAEgASgJUglzdGFydHRpbWUSHgoIdGFibGVhcm4Y44DrzQEgASgJUgh0YW'
    'JsZWFybg==');

@$core.Deprecated('Use importTableDescriptionDescriptor instead')
const ImportTableDescription$json = {
  '1': 'ImportTableDescription',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'cloudwatchloggrouparn',
      '3': 171042008,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchloggrouparn'
    },
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'errorcount', '3': 311137001, '4': 1, '5': 3, '10': 'errorcount'},
    {'1': 'failurecode', '3': 84707897, '4': 1, '5': 9, '10': 'failurecode'},
    {
      '1': 'failuremessage',
      '3': 353556937,
      '4': 1,
      '5': 9,
      '10': 'failuremessage'
    },
    {'1': 'importarn', '3': 444379628, '4': 1, '5': 9, '10': 'importarn'},
    {
      '1': 'importstatus',
      '3': 129077631,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'importeditemcount',
      '3': 202622198,
      '4': 1,
      '5': 3,
      '10': 'importeditemcount'
    },
    {
      '1': 'inputcompressiontype',
      '3': 392699396,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.InputCompressionType',
      '10': 'inputcompressiontype'
    },
    {
      '1': 'inputformat',
      '3': 405664101,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.InputFormat',
      '10': 'inputformat'
    },
    {
      '1': 'inputformatoptions',
      '3': 249201403,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.InputFormatOptions',
      '10': 'inputformatoptions'
    },
    {
      '1': 'processeditemcount',
      '3': 340163668,
      '4': 1,
      '5': 3,
      '10': 'processeditemcount'
    },
    {
      '1': 'processedsizebytes',
      '3': 72287562,
      '4': 1,
      '5': 3,
      '10': 'processedsizebytes'
    },
    {
      '1': 's3bucketsource',
      '3': 202310037,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.S3BucketSource',
      '10': 's3bucketsource'
    },
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {
      '1': 'tablecreationparameters',
      '3': 487419839,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableCreationParameters',
      '10': 'tablecreationparameters'
    },
    {'1': 'tableid', '3': 449893011, '4': 1, '5': 9, '10': 'tableid'},
  ],
};

/// Descriptor for `ImportTableDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importTableDescriptionDescriptor = $convert.base64Decode(
    'ChZJbXBvcnRUYWJsZURlc2NyaXB0aW9uEiMKC2NsaWVudHRva2VuGMz7u0EgASgJUgtjbGllbn'
    'R0b2tlbhI3ChVjbG91ZHdhdGNobG9nZ3JvdXBhcm4Y2MnHUSABKAlSFWNsb3Vkd2F0Y2hsb2dn'
    'cm91cGFybhIbCgdlbmR0aW1lGMzvvB4gASgJUgdlbmR0aW1lEiIKCmVycm9yY291bnQY6aWulA'
    'EgASgDUgplcnJvcmNvdW50EiMKC2ZhaWx1cmVjb2RlGLmUsiggASgJUgtmYWlsdXJlY29kZRIq'
    'Cg5mYWlsdXJlbWVzc2FnZRjJs8uoASABKAlSDmZhaWx1cmVtZXNzYWdlEiAKCWltcG9ydGFybh'
    'js4/LTASABKAlSCWltcG9ydGFybhI9CgxpbXBvcnRzdGF0dXMY/6LGPSABKA4yFi5keW5hbW9k'
    'Yi5JbXBvcnRTdGF0dXNSDGltcG9ydHN0YXR1cxIvChFpbXBvcnRlZGl0ZW1jb3VudBj2ic9gIA'
    'EoA1IRaW1wb3J0ZWRpdGVtY291bnQSVgoUaW5wdXRjb21wcmVzc2lvbnR5cGUYhLyguwEgASgO'
    'Mh4uZHluYW1vZGIuSW5wdXRDb21wcmVzc2lvblR5cGVSFGlucHV0Y29tcHJlc3Npb250eXBlEj'
    'sKC2lucHV0Zm9ybWF0GOXit8EBIAEoDjIVLmR5bmFtb2RiLklucHV0Rm9ybWF0UgtpbnB1dGZv'
    'cm1hdBJPChJpbnB1dGZvcm1hdG9wdGlvbnMY+4XqdiABKAsyHC5keW5hbW9kYi5JbnB1dEZvcm'
    '1hdE9wdGlvbnNSEmlucHV0Zm9ybWF0b3B0aW9ucxIyChJwcm9jZXNzZWRpdGVtY291bnQY1PiZ'
    'ogEgASgDUhJwcm9jZXNzZWRpdGVtY291bnQSMQoScHJvY2Vzc2Vkc2l6ZWJ5dGVzGMqKvCIgAS'
    'gDUhJwcm9jZXNzZWRzaXplYnl0ZXMSQwoOczNidWNrZXRzb3VyY2UYlYO8YCABKAsyGC5keW5h'
    'bW9kYi5TM0J1Y2tldFNvdXJjZVIOczNidWNrZXRzb3VyY2USIAoJc3RhcnR0aW1lGO+05bABIA'
    'EoCVIJc3RhcnR0aW1lEh4KCHRhYmxlYXJuGOOA680BIAEoCVIIdGFibGVhcm4SXwoXdGFibGVj'
    'cmVhdGlvbnBhcmFtZXRlcnMYv9+16AEgASgLMiEuZHluYW1vZGIuVGFibGVDcmVhdGlvblBhcm'
    'FtZXRlcnNSF3RhYmxlY3JlYXRpb25wYXJhbWV0ZXJzEhwKB3RhYmxlaWQYk6XD1gEgASgJUgd0'
    'YWJsZWlk');

@$core.Deprecated('Use importTableInputDescriptor instead')
const ImportTableInput$json = {
  '1': 'ImportTableInput',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'inputcompressiontype',
      '3': 392699396,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.InputCompressionType',
      '10': 'inputcompressiontype'
    },
    {
      '1': 'inputformat',
      '3': 405664101,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.InputFormat',
      '10': 'inputformat'
    },
    {
      '1': 'inputformatoptions',
      '3': 249201403,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.InputFormatOptions',
      '10': 'inputformatoptions'
    },
    {
      '1': 's3bucketsource',
      '3': 202310037,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.S3BucketSource',
      '10': 's3bucketsource'
    },
    {
      '1': 'tablecreationparameters',
      '3': 487419839,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableCreationParameters',
      '10': 'tablecreationparameters'
    },
  ],
};

/// Descriptor for `ImportTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importTableInputDescriptor = $convert.base64Decode(
    'ChBJbXBvcnRUYWJsZUlucHV0EiMKC2NsaWVudHRva2VuGMz7u0EgASgJUgtjbGllbnR0b2tlbh'
    'JWChRpbnB1dGNvbXByZXNzaW9udHlwZRiEvKC7ASABKA4yHi5keW5hbW9kYi5JbnB1dENvbXBy'
    'ZXNzaW9uVHlwZVIUaW5wdXRjb21wcmVzc2lvbnR5cGUSOwoLaW5wdXRmb3JtYXQY5eK3wQEgAS'
    'gOMhUuZHluYW1vZGIuSW5wdXRGb3JtYXRSC2lucHV0Zm9ybWF0Ek8KEmlucHV0Zm9ybWF0b3B0'
    'aW9ucxj7hep2IAEoCzIcLmR5bmFtb2RiLklucHV0Rm9ybWF0T3B0aW9uc1ISaW5wdXRmb3JtYX'
    'RvcHRpb25zEkMKDnMzYnVja2V0c291cmNlGJWDvGAgASgLMhguZHluYW1vZGIuUzNCdWNrZXRT'
    'b3VyY2VSDnMzYnVja2V0c291cmNlEl8KF3RhYmxlY3JlYXRpb25wYXJhbWV0ZXJzGL/ftegBIA'
    'EoCzIhLmR5bmFtb2RiLlRhYmxlQ3JlYXRpb25QYXJhbWV0ZXJzUhd0YWJsZWNyZWF0aW9ucGFy'
    'YW1ldGVycw==');

@$core.Deprecated('Use importTableOutputDescriptor instead')
const ImportTableOutput$json = {
  '1': 'ImportTableOutput',
  '2': [
    {
      '1': 'importtabledescription',
      '3': 407746467,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ImportTableDescription',
      '10': 'importtabledescription'
    },
  ],
};

/// Descriptor for `ImportTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importTableOutputDescriptor = $convert.base64Decode(
    'ChFJbXBvcnRUYWJsZU91dHB1dBJcChZpbXBvcnR0YWJsZWRlc2NyaXB0aW9uGKPvtsIBIAEoCz'
    'IgLmR5bmFtb2RiLkltcG9ydFRhYmxlRGVzY3JpcHRpb25SFmltcG9ydHRhYmxlZGVzY3JpcHRp'
    'b24=');

@$core.Deprecated('Use incrementalExportSpecificationDescriptor instead')
const IncrementalExportSpecification$json = {
  '1': 'IncrementalExportSpecification',
  '2': [
    {
      '1': 'exportfromtime',
      '3': 466248415,
      '4': 1,
      '5': 9,
      '10': 'exportfromtime'
    },
    {'1': 'exporttotime', '3': 242369368, '4': 1, '5': 9, '10': 'exporttotime'},
    {
      '1': 'exportviewtype',
      '3': 26765639,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ExportViewType',
      '10': 'exportviewtype'
    },
  ],
};

/// Descriptor for `IncrementalExportSpecification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List incrementalExportSpecificationDescriptor =
    $convert.base64Decode(
        'Ch5JbmNyZW1lbnRhbEV4cG9ydFNwZWNpZmljYXRpb24SKgoOZXhwb3J0ZnJvbXRpbWUY38Wp3g'
        'EgASgJUg5leHBvcnRmcm9tdGltZRIlCgxleHBvcnR0b3RpbWUY2IbJcyABKAlSDGV4cG9ydHRv'
        'dGltZRJDCg5leHBvcnR2aWV3dHlwZRjH0uEMIAEoDjIYLmR5bmFtb2RiLkV4cG9ydFZpZXdUeX'
        'BlUg5leHBvcnR2aWV3dHlwZQ==');

@$core.Deprecated('Use indexNotFoundExceptionDescriptor instead')
const IndexNotFoundException$json = {
  '1': 'IndexNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IndexNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List indexNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbmRleE5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use inputFormatOptionsDescriptor instead')
const InputFormatOptions$json = {
  '1': 'InputFormatOptions',
  '2': [
    {
      '1': 'csv',
      '3': 450056448,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.CsvOptions',
      '10': 'csv'
    },
  ],
};

/// Descriptor for `InputFormatOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inputFormatOptionsDescriptor = $convert.base64Decode(
    'ChJJbnB1dEZvcm1hdE9wdGlvbnMSKgoDY3N2GICizdYBIAEoCzIULmR5bmFtb2RiLkNzdk9wdG'
    'lvbnNSA2Nzdg==');

@$core.Deprecated('Use internalServerErrorDescriptor instead')
const InternalServerError$json = {
  '1': 'InternalServerError',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServerError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServerErrorDescriptor =
    $convert.base64Decode(
        'ChNJbnRlcm5hbFNlcnZlckVycm9yEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidEndpointExceptionDescriptor instead')
const InvalidEndpointException$json = {
  '1': 'InvalidEndpointException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEndpointException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEndpointExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkRW5kcG9pbnRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidExportTimeExceptionDescriptor instead')
const InvalidExportTimeException$json = {
  '1': 'InvalidExportTimeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidExportTimeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidExportTimeExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkRXhwb3J0VGltZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidRestoreTimeExceptionDescriptor instead')
const InvalidRestoreTimeException$json = {
  '1': 'InvalidRestoreTimeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidRestoreTimeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRestoreTimeExceptionDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkUmVzdG9yZVRpbWVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use itemCollectionMetricsDescriptor instead')
const ItemCollectionMetrics$json = {
  '1': 'ItemCollectionMetrics',
  '2': [
    {
      '1': 'itemcollectionkey',
      '3': 6645926,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ItemCollectionMetrics.ItemcollectionkeyEntry',
      '10': 'itemcollectionkey'
    },
    {
      '1': 'sizeestimaterangegb',
      '3': 271707895,
      '4': 3,
      '5': 1,
      '10': 'sizeestimaterangegb'
    },
  ],
  '3': [ItemCollectionMetrics_ItemcollectionkeyEntry$json],
};

@$core.Deprecated('Use itemCollectionMetricsDescriptor instead')
const ItemCollectionMetrics_ItemcollectionkeyEntry$json = {
  '1': 'ItemcollectionkeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ItemCollectionMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List itemCollectionMetricsDescriptor = $convert.base64Decode(
    'ChVJdGVtQ29sbGVjdGlvbk1ldHJpY3MSZwoRaXRlbWNvbGxlY3Rpb25rZXkYptGVAyADKAsyNi'
    '5keW5hbW9kYi5JdGVtQ29sbGVjdGlvbk1ldHJpY3MuSXRlbWNvbGxlY3Rpb25rZXlFbnRyeVIR'
    'aXRlbWNvbGxlY3Rpb25rZXkSNAoTc2l6ZWVzdGltYXRlcmFuZ2VnYhj33ceBASADKAFSE3Npem'
    'Vlc3RpbWF0ZXJhbmdlZ2IaXgoWSXRlbWNvbGxlY3Rpb25rZXlFbnRyeRIQCgNrZXkYASABKAlS'
    'A2tleRIuCgV2YWx1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOA'
    'E=');

@$core.Deprecated(
    'Use itemCollectionSizeLimitExceededExceptionDescriptor instead')
const ItemCollectionSizeLimitExceededException$json = {
  '1': 'ItemCollectionSizeLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ItemCollectionSizeLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List itemCollectionSizeLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CihJdGVtQ29sbGVjdGlvblNpemVMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5Z'
        'HIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use itemResponseDescriptor instead')
const ItemResponse$json = {
  '1': 'ItemResponse',
  '2': [
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ItemResponse.ItemEntry',
      '10': 'item'
    },
  ],
  '3': [ItemResponse_ItemEntry$json],
};

@$core.Deprecated('Use itemResponseDescriptor instead')
const ItemResponse_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ItemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List itemResponseDescriptor = $convert.base64Decode(
    'CgxJdGVtUmVzcG9uc2USOAoEaXRlbRiHgJL7ASADKAsyIC5keW5hbW9kYi5JdGVtUmVzcG9uc2'
    'UuSXRlbUVudHJ5UgRpdGVtGlEKCUl0ZW1FbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1'
    'ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use keySchemaElementDescriptor instead')
const KeySchemaElement$json = {
  '1': 'KeySchemaElement',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'keytype',
      '3': 5029221,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.KeyType',
      '10': 'keytype'
    },
  ],
};

/// Descriptor for `KeySchemaElement`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keySchemaElementDescriptor = $convert.base64Decode(
    'ChBLZXlTY2hlbWFFbGVtZW50EigKDWF0dHJpYnV0ZW5hbWUYrZWYqAEgASgJUg1hdHRyaWJ1dG'
    'VuYW1lEi4KB2tleXR5cGUY5fqyAiABKA4yES5keW5hbW9kYi5LZXlUeXBlUgdrZXl0eXBl');

@$core.Deprecated('Use keysAndAttributesDescriptor instead')
const KeysAndAttributes$json = {
  '1': 'KeysAndAttributes',
  '2': [
    {
      '1': 'attributestoget',
      '3': 311382592,
      '4': 3,
      '5': 9,
      '10': 'attributestoget'
    },
    {
      '1': 'consistentread',
      '3': 531556994,
      '4': 1,
      '5': 8,
      '10': 'consistentread'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeysAndAttributes.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'keys',
      '3': 2831086,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeysAndAttributes.KeysEntry',
      '10': 'keys'
    },
    {
      '1': 'projectionexpression',
      '3': 150730243,
      '4': 1,
      '5': 9,
      '10': 'projectionexpression'
    },
  ],
  '3': [
    KeysAndAttributes_ExpressionattributenamesEntry$json,
    KeysAndAttributes_KeysEntry$json
  ],
};

@$core.Deprecated('Use keysAndAttributesDescriptor instead')
const KeysAndAttributes_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use keysAndAttributesDescriptor instead')
const KeysAndAttributes_KeysEntry$json = {
  '1': 'KeysEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `KeysAndAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keysAndAttributesDescriptor = $convert.base64Decode(
    'ChFLZXlzQW5kQXR0cmlidXRlcxIsCg9hdHRyaWJ1dGVzdG9nZXQYwKS9lAEgAygJUg9hdHRyaW'
    'J1dGVzdG9nZXQSKgoOY29uc2lzdGVudHJlYWQYgtW7/QEgASgIUg5jb25zaXN0ZW50cmVhZBJ4'
    'ChhleHByZXNzaW9uYXR0cmlidXRlbmFtZXMY/JjRRyADKAsyOS5keW5hbW9kYi5LZXlzQW5kQX'
    'R0cmlidXRlcy5FeHByZXNzaW9uYXR0cmlidXRlbmFtZXNFbnRyeVIYZXhwcmVzc2lvbmF0dHJp'
    'YnV0ZW5hbWVzEjwKBGtleXMY7uWsASADKAsyJS5keW5hbW9kYi5LZXlzQW5kQXR0cmlidXRlcy'
    '5LZXlzRW50cnlSBGtleXMSNQoUcHJvamVjdGlvbmV4cHJlc3Npb24Yg+zvRyABKAlSFHByb2pl'
    'Y3Rpb25leHByZXNzaW9uGksKHUV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5EhAKA2tleR'
    'gBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaUQoJS2V5c0VudHJ5EhAKA2tl'
    'eRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBX'
    'ZhbHVlOgI4AQ==');

@$core.Deprecated('Use kinesisDataStreamDestinationDescriptor instead')
const KinesisDataStreamDestination$json = {
  '1': 'KinesisDataStreamDestination',
  '2': [
    {
      '1': 'approximatecreationdatetimeprecision',
      '3': 392293352,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ApproximateCreationDateTimePrecision',
      '10': 'approximatecreationdatetimeprecision'
    },
    {
      '1': 'destinationstatus',
      '3': 381248234,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.DestinationStatus',
      '10': 'destinationstatus'
    },
    {
      '1': 'destinationstatusdescription',
      '3': 499573086,
      '4': 1,
      '5': 9,
      '10': 'destinationstatusdescription'
    },
    {'1': 'streamarn', '3': 513423709, '4': 1, '5': 9, '10': 'streamarn'},
  ],
};

/// Descriptor for `KinesisDataStreamDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kinesisDataStreamDestinationDescriptor = $convert.base64Decode(
    'ChxLaW5lc2lzRGF0YVN0cmVhbURlc3RpbmF0aW9uEoYBCiRhcHByb3hpbWF0ZWNyZWF0aW9uZG'
    'F0ZXRpbWVwcmVjaXNpb24Y6NeHuwEgASgOMi4uZHluYW1vZGIuQXBwcm94aW1hdGVDcmVhdGlv'
    'bkRhdGVUaW1lUHJlY2lzaW9uUiRhcHByb3hpbWF0ZWNyZWF0aW9uZGF0ZXRpbWVwcmVjaXNpb2'
    '4STQoRZGVzdGluYXRpb25zdGF0dXMY6sXltQEgASgOMhsuZHluYW1vZGIuRGVzdGluYXRpb25T'
    'dGF0dXNSEWRlc3RpbmF0aW9uc3RhdHVzEkYKHGRlc3RpbmF0aW9uc3RhdHVzZGVzY3JpcHRpb2'
    '4Y3sKb7gEgASgJUhxkZXN0aW5hdGlvbnN0YXR1c2Rlc2NyaXB0aW9uEiAKCXN0cmVhbWFybhjd'
    '8uj0ASABKAlSCXN0cmVhbWFybg==');

@$core.Deprecated('Use kinesisStreamingDestinationInputDescriptor instead')
const KinesisStreamingDestinationInput$json = {
  '1': 'KinesisStreamingDestinationInput',
  '2': [
    {
      '1': 'enablekinesisstreamingconfiguration',
      '3': 300307843,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.EnableKinesisStreamingConfiguration',
      '10': 'enablekinesisstreamingconfiguration'
    },
    {'1': 'streamarn', '3': 513423709, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `KinesisStreamingDestinationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kinesisStreamingDestinationInputDescriptor = $convert.base64Decode(
    'CiBLaW5lc2lzU3RyZWFtaW5nRGVzdGluYXRpb25JbnB1dBKDAQojZW5hYmxla2luZXNpc3N0cm'
    'VhbWluZ2NvbmZpZ3VyYXRpb24Yg6uZjwEgASgLMi0uZHluYW1vZGIuRW5hYmxlS2luZXNpc1N0'
    'cmVhbWluZ0NvbmZpZ3VyYXRpb25SI2VuYWJsZWtpbmVzaXNzdHJlYW1pbmdjb25maWd1cmF0aW'
    '9uEiAKCXN0cmVhbWFybhjd8uj0ASABKAlSCXN0cmVhbWFybhIgCgl0YWJsZW5hbWUY3eTagQEg'
    'ASgJUgl0YWJsZW5hbWU=');

@$core.Deprecated('Use kinesisStreamingDestinationOutputDescriptor instead')
const KinesisStreamingDestinationOutput$json = {
  '1': 'KinesisStreamingDestinationOutput',
  '2': [
    {
      '1': 'destinationstatus',
      '3': 381248234,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.DestinationStatus',
      '10': 'destinationstatus'
    },
    {
      '1': 'enablekinesisstreamingconfiguration',
      '3': 300307843,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.EnableKinesisStreamingConfiguration',
      '10': 'enablekinesisstreamingconfiguration'
    },
    {'1': 'streamarn', '3': 513423709, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `KinesisStreamingDestinationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kinesisStreamingDestinationOutputDescriptor = $convert.base64Decode(
    'CiFLaW5lc2lzU3RyZWFtaW5nRGVzdGluYXRpb25PdXRwdXQSTQoRZGVzdGluYXRpb25zdGF0dX'
    'MY6sXltQEgASgOMhsuZHluYW1vZGIuRGVzdGluYXRpb25TdGF0dXNSEWRlc3RpbmF0aW9uc3Rh'
    'dHVzEoMBCiNlbmFibGVraW5lc2lzc3RyZWFtaW5nY29uZmlndXJhdGlvbhiDq5mPASABKAsyLS'
    '5keW5hbW9kYi5FbmFibGVLaW5lc2lzU3RyZWFtaW5nQ29uZmlndXJhdGlvblIjZW5hYmxla2lu'
    'ZXNpc3N0cmVhbWluZ2NvbmZpZ3VyYXRpb24SIAoJc3RyZWFtYXJuGN3y6PQBIAEoCVIJc3RyZW'
    'FtYXJuEiAKCXRhYmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZQ==');

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

@$core.Deprecated('Use listBackupsInputDescriptor instead')
const ListBackupsInput$json = {
  '1': 'ListBackupsInput',
  '2': [
    {
      '1': 'backuptype',
      '3': 134973992,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BackupTypeFilter',
      '10': 'backuptype'
    },
    {
      '1': 'exclusivestartbackuparn',
      '3': 381696959,
      '4': 1,
      '5': 9,
      '10': 'exclusivestartbackuparn'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'timerangelowerbound',
      '3': 506236543,
      '4': 1,
      '5': 9,
      '10': 'timerangelowerbound'
    },
    {
      '1': 'timerangeupperbound',
      '3': 420889434,
      '4': 1,
      '5': 9,
      '10': 'timerangeupperbound'
    },
  ],
};

/// Descriptor for `ListBackupsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBackupsInputDescriptor = $convert.base64Decode(
    'ChBMaXN0QmFja3Vwc0lucHV0Ej0KCmJhY2t1cHR5cGUYqJSuQCABKA4yGi5keW5hbW9kYi5CYW'
    'NrdXBUeXBlRmlsdGVyUgpiYWNrdXB0eXBlEjwKF2V4Y2x1c2l2ZXN0YXJ0YmFja3VwYXJuGL/3'
    'gLYBIAEoCVIXZXhjbHVzaXZlc3RhcnRiYWNrdXBhcm4SGAoFbGltaXQY1ZXZxAEgASgFUgVsaW'
    '1pdBIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW5hbWUSNAoTdGltZXJhbmdlbG93ZXJi'
    'b3VuZBj/nLLxASABKAlSE3RpbWVyYW5nZWxvd2VyYm91bmQSNAoTdGltZXJhbmdldXBwZXJib3'
    'VuZBjahtnIASABKAlSE3RpbWVyYW5nZXVwcGVyYm91bmQ=');

@$core.Deprecated('Use listBackupsOutputDescriptor instead')
const ListBackupsOutput$json = {
  '1': 'ListBackupsOutput',
  '2': [
    {
      '1': 'backupsummaries',
      '3': 375274198,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.BackupSummary',
      '10': 'backupsummaries'
    },
    {
      '1': 'lastevaluatedbackuparn',
      '3': 254744240,
      '4': 1,
      '5': 9,
      '10': 'lastevaluatedbackuparn'
    },
  ],
};

/// Descriptor for `ListBackupsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBackupsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0QmFja3Vwc091dHB1dBJFCg9iYWNrdXBzdW1tYXJpZXMY1vX4sgEgAygLMhcuZHluYW'
    '1vZGIuQmFja3VwU3VtbWFyeVIPYmFja3Vwc3VtbWFyaWVzEjkKFmxhc3RldmFsdWF0ZWRiYWNr'
    'dXBhcm4YsK28eSABKAlSFmxhc3RldmFsdWF0ZWRiYWNrdXBhcm4=');

@$core.Deprecated('Use listContributorInsightsInputDescriptor instead')
const ListContributorInsightsInput$json = {
  '1': 'ListContributorInsightsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `ListContributorInsightsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContributorInsightsInputDescriptor =
    $convert.base64Decode(
        'ChxMaXN0Q29udHJpYnV0b3JJbnNpZ2h0c0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUg'
        'ptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiAKCXRhYmxlbmFt'
        'ZRjd5NqBASABKAlSCXRhYmxlbmFtZQ==');

@$core.Deprecated('Use listContributorInsightsOutputDescriptor instead')
const ListContributorInsightsOutput$json = {
  '1': 'ListContributorInsightsOutput',
  '2': [
    {
      '1': 'contributorinsightssummaries',
      '3': 228077470,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ContributorInsightsSummary',
      '10': 'contributorinsightssummaries'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListContributorInsightsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContributorInsightsOutputDescriptor = $convert.base64Decode(
    'Ch1MaXN0Q29udHJpYnV0b3JJbnNpZ2h0c091dHB1dBJrChxjb250cmlidXRvcmluc2lnaHRzc3'
    'VtbWFyaWVzGJ7f4GwgAygLMiQuZHluYW1vZGIuQ29udHJpYnV0b3JJbnNpZ2h0c1N1bW1hcnlS'
    'HGNvbnRyaWJ1dG9yaW5zaWdodHNzdW1tYXJpZXMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
    'h0dG9rZW4=');

@$core.Deprecated('Use listExportsInputDescriptor instead')
const ListExportsInput$json = {
  '1': 'ListExportsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
  ],
};

/// Descriptor for `ListExportsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExportsInputDescriptor = $convert.base64Decode(
    'ChBMaXN0RXhwb3J0c0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh'
    '8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEh4KCHRhYmxlYXJuGOOA680BIAEoCVII'
    'dGFibGVhcm4=');

@$core.Deprecated('Use listExportsOutputDescriptor instead')
const ListExportsOutput$json = {
  '1': 'ListExportsOutput',
  '2': [
    {
      '1': 'exportsummaries',
      '3': 378369890,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ExportSummary',
      '10': 'exportsummaries'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListExportsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExportsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0RXhwb3J0c091dHB1dBJFCg9leHBvcnRzdW1tYXJpZXMY4u61tAEgAygLMhcuZHluYW'
    '1vZGIuRXhwb3J0U3VtbWFyeVIPZXhwb3J0c3VtbWFyaWVzEh8KCW5leHR0b2tlbhj+hLpnIAEo'
    'CVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listGlobalTablesInputDescriptor instead')
const ListGlobalTablesInput$json = {
  '1': 'ListGlobalTablesInput',
  '2': [
    {
      '1': 'exclusivestartglobaltablename',
      '3': 89274086,
      '4': 1,
      '5': 9,
      '10': 'exclusivestartglobaltablename'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
  ],
};

/// Descriptor for `ListGlobalTablesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGlobalTablesInputDescriptor = $convert.base64Decode(
    'ChVMaXN0R2xvYmFsVGFibGVzSW5wdXQSRwodZXhjbHVzaXZlc3RhcnRnbG9iYWx0YWJsZW5hbW'
    'UY5u3IKiABKAlSHWV4Y2x1c2l2ZXN0YXJ0Z2xvYmFsdGFibGVuYW1lEhgKBWxpbWl0GNWV2cQB'
    'IAEoBVIFbGltaXQSIQoKcmVnaW9ubmFtZRi/m7k1IAEoCVIKcmVnaW9ubmFtZQ==');

@$core.Deprecated('Use listGlobalTablesOutputDescriptor instead')
const ListGlobalTablesOutput$json = {
  '1': 'ListGlobalTablesOutput',
  '2': [
    {
      '1': 'globaltables',
      '3': 213108516,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalTable',
      '10': 'globaltables'
    },
    {
      '1': 'lastevaluatedglobaltablename',
      '3': 157875353,
      '4': 1,
      '5': 9,
      '10': 'lastevaluatedglobaltablename'
    },
  ],
};

/// Descriptor for `ListGlobalTablesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGlobalTablesOutputDescriptor = $convert.base64Decode(
    'ChZMaXN0R2xvYmFsVGFibGVzT3V0cHV0EjwKDGdsb2JhbHRhYmxlcxikjs9lIAMoCzIVLmR5bm'
    'Ftb2RiLkdsb2JhbFRhYmxlUgxnbG9iYWx0YWJsZXMSRQocbGFzdGV2YWx1YXRlZGdsb2JhbHRh'
    'YmxlbmFtZRiZ+aNLIAEoCVIcbGFzdGV2YWx1YXRlZGdsb2JhbHRhYmxlbmFtZQ==');

@$core.Deprecated('Use listImportsInputDescriptor instead')
const ListImportsInput$json = {
  '1': 'ListImportsInput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
  ],
};

/// Descriptor for `ListImportsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportsInputDescriptor = $convert.base64Decode(
    'ChBMaXN0SW1wb3J0c0lucHV0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEh4KCH'
    'BhZ2VzaXplGLiTgtEBIAEoBVIIcGFnZXNpemUSHgoIdGFibGVhcm4Y44DrzQEgASgJUgh0YWJs'
    'ZWFybg==');

@$core.Deprecated('Use listImportsOutputDescriptor instead')
const ListImportsOutput$json = {
  '1': 'ListImportsOutput',
  '2': [
    {
      '1': 'importsummarylist',
      '3': 19020549,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ImportSummary',
      '10': 'importsummarylist'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListImportsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0SW1wb3J0c091dHB1dBJIChFpbXBvcnRzdW1tYXJ5bGlzdBiF9ogJIAMoCzIXLmR5bm'
    'Ftb2RiLkltcG9ydFN1bW1hcnlSEWltcG9ydHN1bW1hcnlsaXN0Eh8KCW5leHR0b2tlbhj+hLpn'
    'IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listTablesInputDescriptor instead')
const ListTablesInput$json = {
  '1': 'ListTablesInput',
  '2': [
    {
      '1': 'exclusivestarttablename',
      '3': 4554457,
      '4': 1,
      '5': 9,
      '10': 'exclusivestarttablename'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `ListTablesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTablesInputDescriptor = $convert.base64Decode(
    'Cg9MaXN0VGFibGVzSW5wdXQSOwoXZXhjbHVzaXZlc3RhcnR0YWJsZW5hbWUY2f2VAiABKAlSF2'
    'V4Y2x1c2l2ZXN0YXJ0dGFibGVuYW1lEhgKBWxpbWl0GNWV2cQBIAEoBVIFbGltaXQ=');

@$core.Deprecated('Use listTablesOutputDescriptor instead')
const ListTablesOutput$json = {
  '1': 'ListTablesOutput',
  '2': [
    {
      '1': 'lastevaluatedtablename',
      '3': 436880490,
      '4': 1,
      '5': 9,
      '10': 'lastevaluatedtablename'
    },
    {'1': 'tablenames', '3': 354058238, '4': 3, '5': 9, '10': 'tablenames'},
  ],
};

/// Descriptor for `ListTablesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTablesOutputDescriptor = $convert.base64Decode(
    'ChBMaXN0VGFibGVzT3V0cHV0EjoKFmxhc3RldmFsdWF0ZWR0YWJsZW5hbWUY6oip0AEgASgJUh'
    'ZsYXN0ZXZhbHVhdGVkdGFibGVuYW1lEiIKCnRhYmxlbmFtZXMY/v/pqAEgAygJUgp0YWJsZW5h'
    'bWVz');

@$core.Deprecated('Use listTagsOfResourceInputDescriptor instead')
const ListTagsOfResourceInput$json = {
  '1': 'ListTagsOfResourceInput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsOfResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsOfResourceInputDescriptor =
    $convert.base64Decode(
        'ChdMaXN0VGFnc09mUmVzb3VyY2VJbnB1dBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2'
        'tlbhIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use listTagsOfResourceOutputDescriptor instead')
const ListTagsOfResourceOutput$json = {
  '1': 'ListTagsOfResourceOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsOfResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsOfResourceOutputDescriptor =
    $convert.base64Decode(
        'ChhMaXN0VGFnc09mUmVzb3VyY2VPdXRwdXQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
        '9rZW4SJQoEdGFncxjBwfa1ASADKAsyDS5keW5hbW9kYi5UYWdSBHRhZ3M=');

@$core.Deprecated('Use localSecondaryIndexDescriptor instead')
const LocalSecondaryIndex$json = {
  '1': 'LocalSecondaryIndex',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'projection',
      '3': 105045921,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Projection',
      '10': 'projection'
    },
  ],
};

/// Descriptor for `LocalSecondaryIndex`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List localSecondaryIndexDescriptor = $convert.base64Decode(
    'ChNMb2NhbFNlY29uZGFyeUluZGV4Eh8KCWluZGV4bmFtZRiR1eswIAEoCVIJaW5kZXhuYW1lEj'
    'wKCWtleXNjaGVtYRjoz92LASADKAsyGi5keW5hbW9kYi5LZXlTY2hlbWFFbGVtZW50UglrZXlz'
    'Y2hlbWESNwoKcHJvamVjdGlvbhihv4syIAEoCzIULmR5bmFtb2RiLlByb2plY3Rpb25SCnByb2'
    'plY3Rpb24=');

@$core.Deprecated('Use localSecondaryIndexDescriptionDescriptor instead')
const LocalSecondaryIndexDescription$json = {
  '1': 'LocalSecondaryIndexDescription',
  '2': [
    {'1': 'indexarn', '3': 374335615, '4': 1, '5': 9, '10': 'indexarn'},
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'indexsizebytes',
      '3': 395738346,
      '4': 1,
      '5': 3,
      '10': 'indexsizebytes'
    },
    {'1': 'itemcount', '3': 26280022, '4': 1, '5': 3, '10': 'itemcount'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'projection',
      '3': 105045921,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Projection',
      '10': 'projection'
    },
  ],
};

/// Descriptor for `LocalSecondaryIndexDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List localSecondaryIndexDescriptionDescriptor = $convert.base64Decode(
    'Ch5Mb2NhbFNlY29uZGFyeUluZGV4RGVzY3JpcHRpb24SHgoIaW5kZXhhcm4Y/9C/sgEgASgJUg'
    'hpbmRleGFybhIfCglpbmRleG5hbWUYkdXrMCABKAlSCWluZGV4bmFtZRIqCg5pbmRleHNpemVi'
    'eXRlcxjq+dm8ASABKANSDmluZGV4c2l6ZWJ5dGVzEh8KCWl0ZW1jb3VudBjWgMQMIAEoA1IJaX'
    'RlbWNvdW50EjwKCWtleXNjaGVtYRjoz92LASADKAsyGi5keW5hbW9kYi5LZXlTY2hlbWFFbGVt'
    'ZW50UglrZXlzY2hlbWESNwoKcHJvamVjdGlvbhihv4syIAEoCzIULmR5bmFtb2RiLlByb2plY3'
    'Rpb25SCnByb2plY3Rpb24=');

@$core.Deprecated('Use localSecondaryIndexInfoDescriptor instead')
const LocalSecondaryIndexInfo$json = {
  '1': 'LocalSecondaryIndexInfo',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'projection',
      '3': 105045921,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Projection',
      '10': 'projection'
    },
  ],
};

/// Descriptor for `LocalSecondaryIndexInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List localSecondaryIndexInfoDescriptor = $convert.base64Decode(
    'ChdMb2NhbFNlY29uZGFyeUluZGV4SW5mbxIfCglpbmRleG5hbWUYkdXrMCABKAlSCWluZGV4bm'
    'FtZRI8CglrZXlzY2hlbWEY6M/diwEgAygLMhouZHluYW1vZGIuS2V5U2NoZW1hRWxlbWVudFIJ'
    'a2V5c2NoZW1hEjcKCnByb2plY3Rpb24Yob+LMiABKAsyFC5keW5hbW9kYi5Qcm9qZWN0aW9uUg'
    'pwcm9qZWN0aW9u');

@$core.Deprecated('Use onDemandThroughputDescriptor instead')
const OnDemandThroughput$json = {
  '1': 'OnDemandThroughput',
  '2': [
    {
      '1': 'maxreadrequestunits',
      '3': 361996322,
      '4': 1,
      '5': 3,
      '10': 'maxreadrequestunits'
    },
    {
      '1': 'maxwriterequestunits',
      '3': 99152185,
      '4': 1,
      '5': 3,
      '10': 'maxwriterequestunits'
    },
  ],
};

/// Descriptor for `OnDemandThroughput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List onDemandThroughputDescriptor = $convert.base64Decode(
    'ChJPbkRlbWFuZFRocm91Z2hwdXQSNAoTbWF4cmVhZHJlcXVlc3R1bml0cxiiwM6sASABKANSE2'
    '1heHJlYWRyZXF1ZXN0dW5pdHMSNQoUbWF4d3JpdGVyZXF1ZXN0dW5pdHMYueKjLyABKANSFG1h'
    'eHdyaXRlcmVxdWVzdHVuaXRz');

@$core.Deprecated('Use onDemandThroughputOverrideDescriptor instead')
const OnDemandThroughputOverride$json = {
  '1': 'OnDemandThroughputOverride',
  '2': [
    {
      '1': 'maxreadrequestunits',
      '3': 361996322,
      '4': 1,
      '5': 3,
      '10': 'maxreadrequestunits'
    },
  ],
};

/// Descriptor for `OnDemandThroughputOverride`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List onDemandThroughputOverrideDescriptor =
    $convert.base64Decode(
        'ChpPbkRlbWFuZFRocm91Z2hwdXRPdmVycmlkZRI0ChNtYXhyZWFkcmVxdWVzdHVuaXRzGKLAzq'
        'wBIAEoA1ITbWF4cmVhZHJlcXVlc3R1bml0cw==');

@$core.Deprecated('Use parameterizedStatementDescriptor instead')
const ParameterizedStatement$json = {
  '1': 'ParameterizedStatement',
  '2': [
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'parameters'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'statement', '3': 248790199, '4': 1, '5': 9, '10': 'statement'},
  ],
};

/// Descriptor for `ParameterizedStatement`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterizedStatementDescriptor = $convert.base64Decode(
    'ChZQYXJhbWV0ZXJpemVkU3RhdGVtZW50EjwKCnBhcmFtZXRlcnMY+qf+6wEgAygLMhguZHluYW'
    '1vZGIuQXR0cmlidXRlVmFsdWVSCnBhcmFtZXRlcnMSggEKI3JldHVybnZhbHVlc29uY29uZGl0'
    'aW9uY2hlY2tmYWlsdXJlGOCXgQIgASgOMi0uZHluYW1vZGIuUmV0dXJuVmFsdWVzT25Db25kaX'
    'Rpb25DaGVja0ZhaWx1cmVSI3JldHVybnZhbHVlc29uY29uZGl0aW9uY2hlY2tmYWlsdXJlEh8K'
    'CXN0YXRlbWVudBi3+dB2IAEoCVIJc3RhdGVtZW50');

@$core.Deprecated('Use pointInTimeRecoveryDescriptionDescriptor instead')
const PointInTimeRecoveryDescription$json = {
  '1': 'PointInTimeRecoveryDescription',
  '2': [
    {
      '1': 'earliestrestorabledatetime',
      '3': 440478879,
      '4': 1,
      '5': 9,
      '10': 'earliestrestorabledatetime'
    },
    {
      '1': 'latestrestorabledatetime',
      '3': 131030757,
      '4': 1,
      '5': 9,
      '10': 'latestrestorabledatetime'
    },
    {
      '1': 'pointintimerecoverystatus',
      '3': 510838075,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.PointInTimeRecoveryStatus',
      '10': 'pointintimerecoverystatus'
    },
    {
      '1': 'recoveryperiodindays',
      '3': 67453592,
      '4': 1,
      '5': 5,
      '10': 'recoveryperiodindays'
    },
  ],
};

/// Descriptor for `PointInTimeRecoveryDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pointInTimeRecoveryDescriptionDescriptor = $convert.base64Decode(
    'Ch5Qb2ludEluVGltZVJlY292ZXJ5RGVzY3JpcHRpb24SQgoaZWFybGllc3RyZXN0b3JhYmxlZG'
    'F0ZXRpbWUYn9mE0gEgASgJUhplYXJsaWVzdHJlc3RvcmFibGVkYXRldGltZRI9ChhsYXRlc3Ry'
    'ZXN0b3JhYmxlZGF0ZXRpbWUY5b29PiABKAlSGGxhdGVzdHJlc3RvcmFibGVkYXRldGltZRJlCh'
    'lwb2ludGludGltZXJlY292ZXJ5c3RhdHVzGLuKy/MBIAEoDjIjLmR5bmFtb2RiLlBvaW50SW5U'
    'aW1lUmVjb3ZlcnlTdGF0dXNSGXBvaW50aW50aW1lcmVjb3ZlcnlzdGF0dXMSNQoUcmVjb3Zlcn'
    'lwZXJpb2RpbmRheXMYmIWVICABKAVSFHJlY292ZXJ5cGVyaW9kaW5kYXlz');

@$core.Deprecated('Use pointInTimeRecoverySpecificationDescriptor instead')
const PointInTimeRecoverySpecification$json = {
  '1': 'PointInTimeRecoverySpecification',
  '2': [
    {
      '1': 'pointintimerecoveryenabled',
      '3': 294176170,
      '4': 1,
      '5': 8,
      '10': 'pointintimerecoveryenabled'
    },
    {
      '1': 'recoveryperiodindays',
      '3': 67453592,
      '4': 1,
      '5': 5,
      '10': 'recoveryperiodindays'
    },
  ],
};

/// Descriptor for `PointInTimeRecoverySpecification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pointInTimeRecoverySpecificationDescriptor =
    $convert.base64Decode(
        'CiBQb2ludEluVGltZVJlY292ZXJ5U3BlY2lmaWNhdGlvbhJCChpwb2ludGludGltZXJlY292ZX'
        'J5ZW5hYmxlZBiqi6OMASABKAhSGnBvaW50aW50aW1lcmVjb3ZlcnllbmFibGVkEjUKFHJlY292'
        'ZXJ5cGVyaW9kaW5kYXlzGJiFlSAgASgFUhRyZWNvdmVyeXBlcmlvZGluZGF5cw==');

@$core
    .Deprecated('Use pointInTimeRecoveryUnavailableExceptionDescriptor instead')
const PointInTimeRecoveryUnavailableException$json = {
  '1': 'PointInTimeRecoveryUnavailableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PointInTimeRecoveryUnavailableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pointInTimeRecoveryUnavailableExceptionDescriptor =
    $convert.base64Decode(
        'CidQb2ludEluVGltZVJlY292ZXJ5VW5hdmFpbGFibGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkc'
        'gnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use policyNotFoundExceptionDescriptor instead')
const PolicyNotFoundException$json = {
  '1': 'PolicyNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PolicyNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChdQb2xpY3lOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use projectionDescriptor instead')
const Projection$json = {
  '1': 'Projection',
  '2': [
    {
      '1': 'nonkeyattributes',
      '3': 312245447,
      '4': 3,
      '5': 9,
      '10': 'nonkeyattributes'
    },
    {
      '1': 'projectiontype',
      '3': 120720617,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ProjectionType',
      '10': 'projectiontype'
    },
  ],
};

/// Descriptor for `Projection`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List projectionDescriptor = $convert.base64Decode(
    'CgpQcm9qZWN0aW9uEi4KEG5vbmtleWF0dHJpYnV0ZXMYx/nxlAEgAygJUhBub25rZXlhdHRyaW'
    'J1dGVzEkMKDnByb2plY3Rpb250eXBlGOmZyDkgASgOMhguZHluYW1vZGIuUHJvamVjdGlvblR5'
    'cGVSDnByb2plY3Rpb250eXBl');

@$core.Deprecated('Use provisionedThroughputDescriptor instead')
const ProvisionedThroughput$json = {
  '1': 'ProvisionedThroughput',
  '2': [
    {
      '1': 'readcapacityunits',
      '3': 43945489,
      '4': 1,
      '5': 3,
      '10': 'readcapacityunits'
    },
    {
      '1': 'writecapacityunits',
      '3': 26938032,
      '4': 1,
      '5': 3,
      '10': 'writecapacityunits'
    },
  ],
};

/// Descriptor for `ProvisionedThroughput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedThroughputDescriptor = $convert.base64Decode(
    'ChVQcm92aXNpb25lZFRocm91Z2hwdXQSLwoRcmVhZGNhcGFjaXR5dW5pdHMYkZz6FCABKANSEX'
    'JlYWRjYXBhY2l0eXVuaXRzEjEKEndyaXRlY2FwYWNpdHl1bml0cxiwlewMIAEoA1ISd3JpdGVj'
    'YXBhY2l0eXVuaXRz');

@$core.Deprecated('Use provisionedThroughputDescriptionDescriptor instead')
const ProvisionedThroughputDescription$json = {
  '1': 'ProvisionedThroughputDescription',
  '2': [
    {
      '1': 'lastdecreasedatetime',
      '3': 494947041,
      '4': 1,
      '5': 9,
      '10': 'lastdecreasedatetime'
    },
    {
      '1': 'lastincreasedatetime',
      '3': 490688181,
      '4': 1,
      '5': 9,
      '10': 'lastincreasedatetime'
    },
    {
      '1': 'numberofdecreasestoday',
      '3': 34313772,
      '4': 1,
      '5': 3,
      '10': 'numberofdecreasestoday'
    },
    {
      '1': 'readcapacityunits',
      '3': 43945489,
      '4': 1,
      '5': 3,
      '10': 'readcapacityunits'
    },
    {
      '1': 'writecapacityunits',
      '3': 26938032,
      '4': 1,
      '5': 3,
      '10': 'writecapacityunits'
    },
  ],
};

/// Descriptor for `ProvisionedThroughputDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedThroughputDescriptionDescriptor = $convert.base64Decode(
    'CiBQcm92aXNpb25lZFRocm91Z2hwdXREZXNjcmlwdGlvbhI2ChRsYXN0ZGVjcmVhc2VkYXRldG'
    'ltZRjhlYHsASABKAlSFGxhc3RkZWNyZWFzZWRhdGV0aW1lEjYKFGxhc3RpbmNyZWFzZWRhdGV0'
    'aW1lGLWd/ekBIAEoCVIUbGFzdGluY3JlYXNlZGF0ZXRpbWUSOQoWbnVtYmVyb2ZkZWNyZWFzZX'
    'N0b2RheRisrK4QIAEoA1IWbnVtYmVyb2ZkZWNyZWFzZXN0b2RheRIvChFyZWFkY2FwYWNpdHl1'
    'bml0cxiRnPoUIAEoA1IRcmVhZGNhcGFjaXR5dW5pdHMSMQoSd3JpdGVjYXBhY2l0eXVuaXRzGL'
    'CV7AwgASgDUhJ3cml0ZWNhcGFjaXR5dW5pdHM=');

@$core
    .Deprecated('Use provisionedThroughputExceededExceptionDescriptor instead')
const ProvisionedThroughputExceededException$json = {
  '1': 'ProvisionedThroughputExceededException',
  '2': [
    {
      '1': 'throttlingreasons',
      '3': 284436852,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ThrottlingReason',
      '10': 'throttlingreasons'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ProvisionedThroughputExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedThroughputExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiZQcm92aXNpb25lZFRocm91Z2hwdXRFeGNlZWRlZEV4Y2VwdGlvbhJMChF0aHJvdHRsaW5ncm'
        'Vhc29ucxj00tCHASADKAsyGi5keW5hbW9kYi5UaHJvdHRsaW5nUmVhc29uUhF0aHJvdHRsaW5n'
        'cmVhc29ucxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use provisionedThroughputOverrideDescriptor instead')
const ProvisionedThroughputOverride$json = {
  '1': 'ProvisionedThroughputOverride',
  '2': [
    {
      '1': 'readcapacityunits',
      '3': 43945489,
      '4': 1,
      '5': 3,
      '10': 'readcapacityunits'
    },
  ],
};

/// Descriptor for `ProvisionedThroughputOverride`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedThroughputOverrideDescriptor =
    $convert.base64Decode(
        'Ch1Qcm92aXNpb25lZFRocm91Z2hwdXRPdmVycmlkZRIvChFyZWFkY2FwYWNpdHl1bml0cxiRnP'
        'oUIAEoA1IRcmVhZGNhcGFjaXR5dW5pdHM=');

@$core.Deprecated('Use putDescriptor instead')
const Put$json = {
  '1': 'Put',
  '2': [
    {
      '1': 'conditionexpression',
      '3': 409657405,
      '4': 1,
      '5': 9,
      '10': 'conditionexpression'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Put.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Put.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Put.ItemEntry',
      '10': 'item'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [
    Put_ExpressionattributenamesEntry$json,
    Put_ExpressionattributevaluesEntry$json,
    Put_ItemEntry$json
  ],
};

@$core.Deprecated('Use putDescriptor instead')
const Put_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putDescriptor instead')
const Put_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putDescriptor instead')
const Put_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `Put`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDescriptor = $convert.base64Decode(
    'CgNQdXQSNAoTY29uZGl0aW9uZXhwcmVzc2lvbhi9wKvDASABKAlSE2NvbmRpdGlvbmV4cHJlc3'
    'Npb24SagoYZXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMisuZHluYW1vZGIuUHV0'
    'LkV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9uYXR0cmlidXRlbmFtZX'
    'MSbgoZZXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlcxjYnKDnASADKAsyLC5keW5hbW9kYi5QdXQu'
    'RXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5UhlleHByZXNzaW9uYXR0cmlidXRldmFsdW'
    'VzEi8KBGl0ZW0Yh4CS+wEgAygLMhcuZHluYW1vZGIuUHV0Lkl0ZW1FbnRyeVIEaXRlbRKCAQoj'
    'cmV0dXJudmFsdWVzb25jb25kaXRpb25jaGVja2ZhaWx1cmUY4JeBAiABKA4yLS5keW5hbW9kYi'
    '5SZXR1cm5WYWx1ZXNPbkNvbmRpdGlvbkNoZWNrRmFpbHVyZVIjcmV0dXJudmFsdWVzb25jb25k'
    'aXRpb25jaGVja2ZhaWx1cmUSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1lGksKHU'
    'V4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVl'
    'GAIgASgJUgV2YWx1ZToCOAEaZgoeRXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5EhAKA2'
    'tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVS'
    'BXZhbHVlOgI4ARpRCglJdGVtRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSLgoFdmFsdWUYAiABKA'
    'syGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdWU6AjgB');

@$core.Deprecated('Use putItemInputDescriptor instead')
const PutItemInput$json = {
  '1': 'PutItemInput',
  '2': [
    {
      '1': 'conditionexpression',
      '3': 409657405,
      '4': 1,
      '5': 9,
      '10': 'conditionexpression'
    },
    {
      '1': 'conditionaloperator',
      '3': 172066260,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ConditionalOperator',
      '10': 'conditionaloperator'
    },
    {
      '1': 'expected',
      '3': 106557946,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.PutItemInput.ExpectedEntry',
      '10': 'expected'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.PutItemInput.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.PutItemInput.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.PutItemInput.ItemEntry',
      '10': 'item'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'returnitemcollectionmetrics',
      '3': 255507354,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnItemCollectionMetrics',
      '10': 'returnitemcollectionmetrics'
    },
    {
      '1': 'returnvalues',
      '3': 402960198,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValue',
      '10': 'returnvalues'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [
    PutItemInput_ExpectedEntry$json,
    PutItemInput_ExpressionattributenamesEntry$json,
    PutItemInput_ExpressionattributevaluesEntry$json,
    PutItemInput_ItemEntry$json
  ],
};

@$core.Deprecated('Use putItemInputDescriptor instead')
const PutItemInput_ExpectedEntry$json = {
  '1': 'ExpectedEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ExpectedAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putItemInputDescriptor instead')
const PutItemInput_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putItemInputDescriptor instead')
const PutItemInput_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use putItemInputDescriptor instead')
const PutItemInput_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `PutItemInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putItemInputDescriptor = $convert.base64Decode(
    'CgxQdXRJdGVtSW5wdXQSNAoTY29uZGl0aW9uZXhwcmVzc2lvbhi9wKvDASABKAlSE2NvbmRpdG'
    'lvbmV4cHJlc3Npb24SUgoTY29uZGl0aW9uYWxvcGVyYXRvchjUi4ZSIAEoDjIdLmR5bmFtb2Ri'
    'LkNvbmRpdGlvbmFsT3BlcmF0b3JSE2NvbmRpdGlvbmFsb3BlcmF0b3ISQwoIZXhwZWN0ZWQY+u'
    'PnMiADKAsyJC5keW5hbW9kYi5QdXRJdGVtSW5wdXQuRXhwZWN0ZWRFbnRyeVIIZXhwZWN0ZWQS'
    'cwoYZXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMjQuZHluYW1vZGIuUHV0SXRlbU'
    'lucHV0LkV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9uYXR0cmlidXRl'
    'bmFtZXMSdwoZZXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlcxjYnKDnASADKAsyNS5keW5hbW9kYi'
    '5QdXRJdGVtSW5wdXQuRXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5UhlleHByZXNzaW9u'
    'YXR0cmlidXRldmFsdWVzEjgKBGl0ZW0Yh4CS+wEgAygLMiAuZHluYW1vZGIuUHV0SXRlbUlucH'
    'V0Lkl0ZW1FbnRyeVIEaXRlbRJbChZyZXR1cm5jb25zdW1lZGNhcGFjaXR5GP7n4RQgASgOMiAu'
    'ZHluYW1vZGIuUmV0dXJuQ29uc3VtZWRDYXBhY2l0eVIWcmV0dXJuY29uc3VtZWRjYXBhY2l0eR'
    'JqChtyZXR1cm5pdGVtY29sbGVjdGlvbm1ldHJpY3MYmvfqeSABKA4yJS5keW5hbW9kYi5SZXR1'
    'cm5JdGVtQ29sbGVjdGlvbk1ldHJpY3NSG3JldHVybml0ZW1jb2xsZWN0aW9ubWV0cmljcxI9Cg'
    'xyZXR1cm52YWx1ZXMYxt6SwAEgASgOMhUuZHluYW1vZGIuUmV0dXJuVmFsdWVSDHJldHVybnZh'
    'bHVlcxKCAQojcmV0dXJudmFsdWVzb25jb25kaXRpb25jaGVja2ZhaWx1cmUY4JeBAiABKA4yLS'
    '5keW5hbW9kYi5SZXR1cm5WYWx1ZXNPbkNvbmRpdGlvbkNoZWNrRmFpbHVyZVIjcmV0dXJudmFs'
    'dWVzb25jb25kaXRpb25jaGVja2ZhaWx1cmUSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibG'
    'VuYW1lGl0KDUV4cGVjdGVkRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSNgoFdmFsdWUYAiABKAsy'
    'IC5keW5hbW9kYi5FeHBlY3RlZEF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAEaSwodRXhwcmVzc2'
    'lvbmF0dHJpYnV0ZW5hbWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4ARpmCh5FeHByZXNzaW9uYXR0cmlidXRldmFsdWVzRW50cnkSEAoDa2V5GAEgAS'
    'gJUgNrZXkSLgoFdmFsdWUYAiABKAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdWU6'
    'AjgBGlEKCUl0ZW1FbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1ZRgCIAEoCzIYLmR5bm'
    'Ftb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use putItemOutputDescriptor instead')
const PutItemOutput$json = {
  '1': 'PutItemOutput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.PutItemOutput.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'itemcollectionmetrics',
      '3': 185317452,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ItemCollectionMetrics',
      '10': 'itemcollectionmetrics'
    },
  ],
  '3': [PutItemOutput_AttributesEntry$json],
};

@$core.Deprecated('Use putItemOutputDescriptor instead')
const PutItemOutput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `PutItemOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putItemOutputDescriptor = $convert.base64Decode(
    'Cg1QdXRJdGVtT3V0cHV0EkoKCmF0dHJpYnV0ZXMYtan7YyADKAsyJy5keW5hbW9kYi5QdXRJdG'
    'VtT3V0cHV0LkF0dHJpYnV0ZXNFbnRyeVIKYXR0cmlidXRlcxJKChBjb25zdW1lZGNhcGFjaXR5'
    'GKyqodYBIAEoCzIaLmR5bmFtb2RiLkNvbnN1bWVkQ2FwYWNpdHlSEGNvbnN1bWVkY2FwYWNpdH'
    'kSWAoVaXRlbWNvbGxlY3Rpb25tZXRyaWNzGMzwrlggASgLMh8uZHluYW1vZGIuSXRlbUNvbGxl'
    'Y3Rpb25NZXRyaWNzUhVpdGVtY29sbGVjdGlvbm1ldHJpY3MaVwoPQXR0cmlidXRlc0VudHJ5Eh'
    'AKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFs'
    'dWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use putRequestDescriptor instead')
const PutRequest$json = {
  '1': 'PutRequest',
  '2': [
    {
      '1': 'item',
      '3': 526680071,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.PutRequest.ItemEntry',
      '10': 'item'
    },
  ],
  '3': [PutRequest_ItemEntry$json],
};

@$core.Deprecated('Use putRequestDescriptor instead')
const PutRequest_ItemEntry$json = {
  '1': 'ItemEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `PutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRequestDescriptor = $convert.base64Decode(
    'CgpQdXRSZXF1ZXN0EjYKBGl0ZW0Yh4CS+wEgAygLMh4uZHluYW1vZGIuUHV0UmVxdWVzdC5JdG'
    'VtRW50cnlSBGl0ZW0aUQoJSXRlbUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIg'
    'ASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use putResourcePolicyInputDescriptor instead')
const PutResourcePolicyInput$json = {
  '1': 'PutResourcePolicyInput',
  '2': [
    {
      '1': 'confirmremoveselfresourceaccess',
      '3': 88224484,
      '4': 1,
      '5': 8,
      '10': 'confirmremoveselfresourceaccess'
    },
    {
      '1': 'expectedrevisionid',
      '3': 456463970,
      '4': 1,
      '5': 9,
      '10': 'expectedrevisionid'
    },
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `PutResourcePolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyInputDescriptor = $convert.base64Decode(
    'ChZQdXRSZXNvdXJjZVBvbGljeUlucHV0EksKH2NvbmZpcm1yZW1vdmVzZWxmcmVzb3VyY2VhY2'
    'Nlc3MY5OWIKiABKAhSH2NvbmZpcm1yZW1vdmVzZWxmcmVzb3VyY2VhY2Nlc3MSMgoSZXhwZWN0'
    'ZWRyZXZpc2lvbmlkGOKs1NkBIAEoCVISZXhwZWN0ZWRyZXZpc2lvbmlkEhoKBnBvbGljeRig7/'
    'DgASABKAlSBnBvbGljeRIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use putResourcePolicyOutputDescriptor instead')
const PutResourcePolicyOutput$json = {
  '1': 'PutResourcePolicyOutput',
  '2': [
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `PutResourcePolicyOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyOutputDescriptor =
    $convert.base64Decode(
        'ChdQdXRSZXNvdXJjZVBvbGljeU91dHB1dBIiCgpyZXZpc2lvbmlkGIajnu4BIAEoCVIKcmV2aX'
        'Npb25pZA==');

@$core.Deprecated('Use queryInputDescriptor instead')
const QueryInput$json = {
  '1': 'QueryInput',
  '2': [
    {
      '1': 'attributestoget',
      '3': 311382592,
      '4': 3,
      '5': 9,
      '10': 'attributestoget'
    },
    {
      '1': 'conditionaloperator',
      '3': 172066260,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ConditionalOperator',
      '10': 'conditionaloperator'
    },
    {
      '1': 'consistentread',
      '3': 531556994,
      '4': 1,
      '5': 8,
      '10': 'consistentread'
    },
    {
      '1': 'exclusivestartkey',
      '3': 348137297,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.QueryInput.ExclusivestartkeyEntry',
      '10': 'exclusivestartkey'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.QueryInput.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.QueryInput.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'filterexpression',
      '3': 68019610,
      '4': 1,
      '5': 9,
      '10': 'filterexpression'
    },
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'keyconditionexpression',
      '3': 218579280,
      '4': 1,
      '5': 9,
      '10': 'keyconditionexpression'
    },
    {
      '1': 'keyconditions',
      '3': 307034879,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.QueryInput.KeyconditionsEntry',
      '10': 'keyconditions'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'projectionexpression',
      '3': 150730243,
      '4': 1,
      '5': 9,
      '10': 'projectionexpression'
    },
    {
      '1': 'queryfilter',
      '3': 107880942,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.QueryInput.QueryfilterEntry',
      '10': 'queryfilter'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'scanindexforward',
      '3': 360066954,
      '4': 1,
      '5': 8,
      '10': 'scanindexforward'
    },
    {
      '1': 'select',
      '3': 512305998,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.Select',
      '10': 'select'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
  '3': [
    QueryInput_ExclusivestartkeyEntry$json,
    QueryInput_ExpressionattributenamesEntry$json,
    QueryInput_ExpressionattributevaluesEntry$json,
    QueryInput_KeyconditionsEntry$json,
    QueryInput_QueryfilterEntry$json
  ],
};

@$core.Deprecated('Use queryInputDescriptor instead')
const QueryInput_ExclusivestartkeyEntry$json = {
  '1': 'ExclusivestartkeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use queryInputDescriptor instead')
const QueryInput_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use queryInputDescriptor instead')
const QueryInput_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use queryInputDescriptor instead')
const QueryInput_KeyconditionsEntry$json = {
  '1': 'KeyconditionsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Condition',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use queryInputDescriptor instead')
const QueryInput_QueryfilterEntry$json = {
  '1': 'QueryfilterEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Condition',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `QueryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryInputDescriptor = $convert.base64Decode(
    'CgpRdWVyeUlucHV0EiwKD2F0dHJpYnV0ZXN0b2dldBjApL2UASADKAlSD2F0dHJpYnV0ZXN0b2'
    'dldBJSChNjb25kaXRpb25hbG9wZXJhdG9yGNSLhlIgASgOMh0uZHluYW1vZGIuQ29uZGl0aW9u'
    'YWxPcGVyYXRvclITY29uZGl0aW9uYWxvcGVyYXRvchIqCg5jb25zaXN0ZW50cmVhZBiC1bv9AS'
    'ABKAhSDmNvbnNpc3RlbnRyZWFkEl0KEWV4Y2x1c2l2ZXN0YXJ0a2V5GNHOgKYBIAMoCzIrLmR5'
    'bmFtb2RiLlF1ZXJ5SW5wdXQuRXhjbHVzaXZlc3RhcnRrZXlFbnRyeVIRZXhjbHVzaXZlc3Rhcn'
    'RrZXkScQoYZXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMjIuZHluYW1vZGIuUXVl'
    'cnlJbnB1dC5FeHByZXNzaW9uYXR0cmlidXRlbmFtZXNFbnRyeVIYZXhwcmVzc2lvbmF0dHJpYn'
    'V0ZW5hbWVzEnUKGWV4cHJlc3Npb25hdHRyaWJ1dGV2YWx1ZXMY2Jyg5wEgAygLMjMuZHluYW1v'
    'ZGIuUXVlcnlJbnB1dC5FeHByZXNzaW9uYXR0cmlidXRldmFsdWVzRW50cnlSGWV4cHJlc3Npb2'
    '5hdHRyaWJ1dGV2YWx1ZXMSLQoQZmlsdGVyZXhwcmVzc2lvbhiay7cgIAEoCVIQZmlsdGVyZXhw'
    'cmVzc2lvbhIfCglpbmRleG5hbWUYkdXrMCABKAlSCWluZGV4bmFtZRI5ChZrZXljb25kaXRpb2'
    '5leHByZXNzaW9uGNCCnWggASgJUhZrZXljb25kaXRpb25leHByZXNzaW9uElEKDWtleWNvbmRp'
    'dGlvbnMY//WzkgEgAygLMicuZHluYW1vZGIuUXVlcnlJbnB1dC5LZXljb25kaXRpb25zRW50cn'
    'lSDWtleWNvbmRpdGlvbnMSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBI1ChRwcm9qZWN0aW9u'
    'ZXhwcmVzc2lvbhiD7O9HIAEoCVIUcHJvamVjdGlvbmV4cHJlc3Npb24SSgoLcXVlcnlmaWx0ZX'
    'IY7sO4MyADKAsyJS5keW5hbW9kYi5RdWVyeUlucHV0LlF1ZXJ5ZmlsdGVyRW50cnlSC3F1ZXJ5'
    'ZmlsdGVyElsKFnJldHVybmNvbnN1bWVkY2FwYWNpdHkY/ufhFCABKA4yIC5keW5hbW9kYi5SZX'
    'R1cm5Db25zdW1lZENhcGFjaXR5UhZyZXR1cm5jb25zdW1lZGNhcGFjaXR5Ei4KEHNjYW5pbmRl'
    'eGZvcndhcmQYit/YqwEgASgIUhBzY2FuaW5kZXhmb3J3YXJkEiwKBnNlbGVjdBjO1qT0ASABKA'
    '4yEC5keW5hbW9kYi5TZWxlY3RSBnNlbGVjdBIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJs'
    'ZW5hbWUaXgoWRXhjbHVzaXZlc3RhcnRrZXlFbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YW'
    'x1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAEaSwodRXhwcmVz'
    'c2lvbmF0dHJpYnV0ZW5hbWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKA'
    'lSBXZhbHVlOgI4ARpmCh5FeHByZXNzaW9uYXR0cmlidXRldmFsdWVzRW50cnkSEAoDa2V5GAEg'
    'ASgJUgNrZXkSLgoFdmFsdWUYAiABKAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdW'
    'U6AjgBGlUKEktleWNvbmRpdGlvbnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIpCgV2YWx1ZRgC'
    'IAEoCzITLmR5bmFtb2RiLkNvbmRpdGlvblIFdmFsdWU6AjgBGlMKEFF1ZXJ5ZmlsdGVyRW50cn'
    'kSEAoDa2V5GAEgASgJUgNrZXkSKQoFdmFsdWUYAiABKAsyEy5keW5hbW9kYi5Db25kaXRpb25S'
    'BXZhbHVlOgI4AQ==');

@$core.Deprecated('Use queryOutputDescriptor instead')
const QueryOutput$json = {
  '1': 'QueryOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {'1': 'count', '3': 31963285, '4': 1, '5': 5, '10': 'count'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.QueryOutput.ItemsEntry',
      '10': 'items'
    },
    {
      '1': 'lastevaluatedkey',
      '3': 54319830,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.QueryOutput.LastevaluatedkeyEntry',
      '10': 'lastevaluatedkey'
    },
    {'1': 'scannedcount', '3': 531161315, '4': 1, '5': 5, '10': 'scannedcount'},
  ],
  '3': [QueryOutput_ItemsEntry$json, QueryOutput_LastevaluatedkeyEntry$json],
};

@$core.Deprecated('Use queryOutputDescriptor instead')
const QueryOutput_ItemsEntry$json = {
  '1': 'ItemsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use queryOutputDescriptor instead')
const QueryOutput_LastevaluatedkeyEntry$json = {
  '1': 'LastevaluatedkeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `QueryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryOutputDescriptor = $convert.base64Decode(
    'CgtRdWVyeU91dHB1dBJKChBjb25zdW1lZGNhcGFjaXR5GKyqodYBIAEoCzIaLmR5bmFtb2RiLk'
    'NvbnN1bWVkQ2FwYWNpdHlSEGNvbnN1bWVkY2FwYWNpdHkSFwoFY291bnQYlfGeDyABKAVSBWNv'
    'dW50EjkKBWl0ZW1zGLDw2AEgAygLMiAuZHluYW1vZGIuUXVlcnlPdXRwdXQuSXRlbXNFbnRyeV'
    'IFaXRlbXMSWgoQbGFzdGV2YWx1YXRlZGtleRjWtfMZIAMoCzIrLmR5bmFtb2RiLlF1ZXJ5T3V0'
    'cHV0Lkxhc3RldmFsdWF0ZWRrZXlFbnRyeVIQbGFzdGV2YWx1YXRlZGtleRImCgxzY2FubmVkY2'
    '91bnQY48Gj/QEgASgFUgxzY2FubmVkY291bnQaUgoKSXRlbXNFbnRyeRIQCgNrZXkYASABKAlS'
    'A2tleRIuCgV2YWx1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOA'
    'EaXQoVTGFzdGV2YWx1YXRlZGtleUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIg'
    'ASgLMhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use replicaDescriptor instead')
const Replica$json = {
  '1': 'Replica',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
  ],
};

/// Descriptor for `Replica`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaDescriptor = $convert.base64Decode(
    'CgdSZXBsaWNhEiEKCnJlZ2lvbm5hbWUYv5u5NSABKAlSCnJlZ2lvbm5hbWU=');

@$core.Deprecated('Use replicaAlreadyExistsExceptionDescriptor instead')
const ReplicaAlreadyExistsException$json = {
  '1': 'ReplicaAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ReplicaAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'Ch1SZXBsaWNhQWxyZWFkeUV4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use replicaAutoScalingDescriptionDescriptor instead')
const ReplicaAutoScalingDescription$json = {
  '1': 'ReplicaAutoScalingDescription',
  '2': [
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaGlobalSecondaryIndexAutoScalingDescription',
      '10': 'globalsecondaryindexes'
    },
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'replicaprovisionedreadcapacityautoscalingsettings',
      '3': 125131885,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'replicaprovisionedreadcapacityautoscalingsettings'
    },
    {
      '1': 'replicaprovisionedwritecapacityautoscalingsettings',
      '3': 293841796,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'replicaprovisionedwritecapacityautoscalingsettings'
    },
    {
      '1': 'replicastatus',
      '3': 466739730,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReplicaStatus',
      '10': 'replicastatus'
    },
  ],
};

/// Descriptor for `ReplicaAutoScalingDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaAutoScalingDescriptionDescriptor = $convert.base64Decode(
    'Ch1SZXBsaWNhQXV0b1NjYWxpbmdEZXNjcmlwdGlvbhJ3ChZnbG9iYWxzZWNvbmRhcnlpbmRleG'
    'VzGKn6jMMBIAMoCzI7LmR5bmFtb2RiLlJlcGxpY2FHbG9iYWxTZWNvbmRhcnlJbmRleEF1dG9T'
    'Y2FsaW5nRGVzY3JpcHRpb25SFmdsb2JhbHNlY29uZGFyeWluZGV4ZXMSIQoKcmVnaW9ubmFtZR'
    'i/m7k1IAEoCVIKcmVnaW9ubmFtZRKZAQoxcmVwbGljYXByb3Zpc2lvbmVkcmVhZGNhcGFjaXR5'
    'YXV0b3NjYWxpbmdzZXR0aW5ncxjtuNU7IAEoCzIoLmR5bmFtb2RiLkF1dG9TY2FsaW5nU2V0dG'
    'luZ3NEZXNjcmlwdGlvblIxcmVwbGljYXByb3Zpc2lvbmVkcmVhZGNhcGFjaXR5YXV0b3NjYWxp'
    'bmdzZXR0aW5ncxKcAQoycmVwbGljYXByb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eWF1dG9zY2FsaW'
    '5nc2V0dGluZ3MYhNeOjAEgASgLMiguZHluYW1vZGIuQXV0b1NjYWxpbmdTZXR0aW5nc0Rlc2Ny'
    'aXB0aW9uUjJyZXBsaWNhcHJvdmlzaW9uZWR3cml0ZWNhcGFjaXR5YXV0b3NjYWxpbmdzZXR0aW'
    '5ncxJBCg1yZXBsaWNhc3RhdHVzGJLEx94BIAEoDjIXLmR5bmFtb2RiLlJlcGxpY2FTdGF0dXNS'
    'DXJlcGxpY2FzdGF0dXM=');

@$core.Deprecated('Use replicaAutoScalingUpdateDescriptor instead')
const ReplicaAutoScalingUpdate$json = {
  '1': 'ReplicaAutoScalingUpdate',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'replicaglobalsecondaryindexupdates',
      '3': 16441201,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaGlobalSecondaryIndexAutoScalingUpdate',
      '10': 'replicaglobalsecondaryindexupdates'
    },
    {
      '1': 'replicaprovisionedreadcapacityautoscalingupdate',
      '3': 262279073,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'replicaprovisionedreadcapacityautoscalingupdate'
    },
  ],
};

/// Descriptor for `ReplicaAutoScalingUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaAutoScalingUpdateDescriptor = $convert.base64Decode(
    'ChhSZXBsaWNhQXV0b1NjYWxpbmdVcGRhdGUSIQoKcmVnaW9ubmFtZRi/m7k1IAEoCVIKcmVnaW'
    '9ubmFtZRKJAQoicmVwbGljYWdsb2JhbHNlY29uZGFyeWluZGV4dXBkYXRlcxjxvusHIAMoCzI2'
    'LmR5bmFtb2RiLlJlcGxpY2FHbG9iYWxTZWNvbmRhcnlJbmRleEF1dG9TY2FsaW5nVXBkYXRlUi'
    'JyZXBsaWNhZ2xvYmFsc2Vjb25kYXJ5aW5kZXh1cGRhdGVzEpABCi9yZXBsaWNhcHJvdmlzaW9u'
    'ZWRyZWFkY2FwYWNpdHlhdXRvc2NhbGluZ3VwZGF0ZRihn4h9IAEoCzIjLmR5bmFtb2RiLkF1dG'
    '9TY2FsaW5nU2V0dGluZ3NVcGRhdGVSL3JlcGxpY2Fwcm92aXNpb25lZHJlYWRjYXBhY2l0eWF1'
    'dG9zY2FsaW5ndXBkYXRl');

@$core.Deprecated('Use replicaDescriptionDescriptor instead')
const ReplicaDescription$json = {
  '1': 'ReplicaDescription',
  '2': [
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaGlobalSecondaryIndexDescription',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'globaltablesettingsreplicationmode',
      '3': 10446577,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.GlobalTableSettingsReplicationMode',
      '10': 'globaltablesettingsreplicationmode'
    },
    {
      '1': 'kmsmasterkeyid',
      '3': 521459443,
      '4': 1,
      '5': 9,
      '10': 'kmsmasterkeyid'
    },
    {
      '1': 'ondemandthroughputoverride',
      '3': 317165234,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughputOverride',
      '10': 'ondemandthroughputoverride'
    },
    {
      '1': 'provisionedthroughputoverride',
      '3': 413332116,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughputOverride',
      '10': 'provisionedthroughputoverride'
    },
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'replicainaccessibledatetime',
      '3': 20060608,
      '4': 1,
      '5': 9,
      '10': 'replicainaccessibledatetime'
    },
    {
      '1': 'replicastatus',
      '3': 466739730,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReplicaStatus',
      '10': 'replicastatus'
    },
    {
      '1': 'replicastatusdescription',
      '3': 416683638,
      '4': 1,
      '5': 9,
      '10': 'replicastatusdescription'
    },
    {
      '1': 'replicastatuspercentprogress',
      '3': 456382130,
      '4': 1,
      '5': 9,
      '10': 'replicastatuspercentprogress'
    },
    {
      '1': 'replicatableclasssummary',
      '3': 519483614,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableClassSummary',
      '10': 'replicatableclasssummary'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableWarmThroughputDescription',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `ReplicaDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaDescriptionDescriptor = $convert.base64Decode(
    'ChJSZXBsaWNhRGVzY3JpcHRpb24SbAoWZ2xvYmFsc2Vjb25kYXJ5aW5kZXhlcxip+ozDASADKA'
    'syMC5keW5hbW9kYi5SZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhEZXNjcmlwdGlvblIWZ2xv'
    'YmFsc2Vjb25kYXJ5aW5kZXhlcxJ/CiJnbG9iYWx0YWJsZXNldHRpbmdzcmVwbGljYXRpb25tb2'
    'RlGPHN/QQgASgOMiwuZHluYW1vZGIuR2xvYmFsVGFibGVTZXR0aW5nc1JlcGxpY2F0aW9uTW9k'
    'ZVIiZ2xvYmFsdGFibGVzZXR0aW5nc3JlcGxpY2F0aW9ubW9kZRIqCg5rbXNtYXN0ZXJrZXlpZB'
    'jzrdP4ASABKAlSDmttc21hc3RlcmtleWlkEmgKGm9uZGVtYW5kdGhyb3VnaHB1dG92ZXJyaWRl'
    'GLKdnpcBIAEoCzIkLmR5bmFtb2RiLk9uRGVtYW5kVGhyb3VnaHB1dE92ZXJyaWRlUhpvbmRlbW'
    'FuZHRocm91Z2hwdXRvdmVycmlkZRJxCh1wcm92aXNpb25lZHRocm91Z2hwdXRvdmVycmlkZRiU'
    '5YvFASABKAsyJy5keW5hbW9kYi5Qcm92aXNpb25lZFRocm91Z2hwdXRPdmVycmlkZVIdcHJvdm'
    'lzaW9uZWR0aHJvdWdocHV0b3ZlcnJpZGUSIQoKcmVnaW9ubmFtZRi/m7k1IAEoCVIKcmVnaW9u'
    'bmFtZRJDChtyZXBsaWNhaW5hY2Nlc3NpYmxlZGF0ZXRpbWUYwLPICSABKAlSG3JlcGxpY2Fpbm'
    'FjY2Vzc2libGVkYXRldGltZRJBCg1yZXBsaWNhc3RhdHVzGJLEx94BIAEoDjIXLmR5bmFtb2Ri'
    'LlJlcGxpY2FTdGF0dXNSDXJlcGxpY2FzdGF0dXMSPgoYcmVwbGljYXN0YXR1c2Rlc2NyaXB0aW'
    '9uGPas2MYBIAEoCVIYcmVwbGljYXN0YXR1c2Rlc2NyaXB0aW9uEkYKHHJlcGxpY2FzdGF0dXNw'
    'ZXJjZW50cHJvZ3Jlc3MYsq3P2QEgASgJUhxyZXBsaWNhc3RhdHVzcGVyY2VudHByb2dyZXNzEl'
    'sKGHJlcGxpY2F0YWJsZWNsYXNzc3VtbWFyeRje4dr3ASABKAsyGy5keW5hbW9kYi5UYWJsZUNs'
    'YXNzU3VtbWFyeVIYcmVwbGljYXRhYmxlY2xhc3NzdW1tYXJ5ElQKDndhcm10aHJvdWdocHV0GI'
    'PeyIoBIAEoCzIoLmR5bmFtb2RiLlRhYmxlV2FybVRocm91Z2hwdXREZXNjcmlwdGlvblIOd2Fy'
    'bXRocm91Z2hwdXQ=');

@$core.Deprecated('Use replicaGlobalSecondaryIndexDescriptor instead')
const ReplicaGlobalSecondaryIndex$json = {
  '1': 'ReplicaGlobalSecondaryIndex',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'ondemandthroughputoverride',
      '3': 317165234,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughputOverride',
      '10': 'ondemandthroughputoverride'
    },
    {
      '1': 'provisionedthroughputoverride',
      '3': 413332116,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughputOverride',
      '10': 'provisionedthroughputoverride'
    },
  ],
};

/// Descriptor for `ReplicaGlobalSecondaryIndex`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaGlobalSecondaryIndexDescriptor = $convert.base64Decode(
    'ChtSZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXgSHwoJaW5kZXhuYW1lGJHV6zAgASgJUglpbm'
    'RleG5hbWUSaAoab25kZW1hbmR0aHJvdWdocHV0b3ZlcnJpZGUYsp2elwEgASgLMiQuZHluYW1v'
    'ZGIuT25EZW1hbmRUaHJvdWdocHV0T3ZlcnJpZGVSGm9uZGVtYW5kdGhyb3VnaHB1dG92ZXJyaW'
    'RlEnEKHXByb3Zpc2lvbmVkdGhyb3VnaHB1dG92ZXJyaWRlGJTli8UBIAEoCzInLmR5bmFtb2Ri'
    'LlByb3Zpc2lvbmVkVGhyb3VnaHB1dE92ZXJyaWRlUh1wcm92aXNpb25lZHRocm91Z2hwdXRvdm'
    'VycmlkZQ==');

@$core.Deprecated(
    'Use replicaGlobalSecondaryIndexAutoScalingDescriptionDescriptor instead')
const ReplicaGlobalSecondaryIndexAutoScalingDescription$json = {
  '1': 'ReplicaGlobalSecondaryIndexAutoScalingDescription',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'indexstatus',
      '3': 364436830,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.IndexStatus',
      '10': 'indexstatus'
    },
    {
      '1': 'provisionedreadcapacityautoscalingsettings',
      '3': 85617667,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'provisionedreadcapacityautoscalingsettings'
    },
    {
      '1': 'provisionedwritecapacityautoscalingsettings',
      '3': 168409114,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'provisionedwritecapacityautoscalingsettings'
    },
  ],
};

/// Descriptor for `ReplicaGlobalSecondaryIndexAutoScalingDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    replicaGlobalSecondaryIndexAutoScalingDescriptionDescriptor =
    $convert.base64Decode(
        'CjFSZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhBdXRvU2NhbGluZ0Rlc2NyaXB0aW9uEh8KCW'
        'luZGV4bmFtZRiR1eswIAEoCVIJaW5kZXhuYW1lEjsKC2luZGV4c3RhdHVzGN66460BIAEoDjIV'
        'LmR5bmFtb2RiLkluZGV4U3RhdHVzUgtpbmRleHN0YXR1cxKLAQoqcHJvdmlzaW9uZWRyZWFkY2'
        'FwYWNpdHlhdXRvc2NhbGluZ3NldHRpbmdzGIPY6SggASgLMiguZHluYW1vZGIuQXV0b1NjYWxp'
        'bmdTZXR0aW5nc0Rlc2NyaXB0aW9uUipwcm92aXNpb25lZHJlYWRjYXBhY2l0eWF1dG9zY2FsaW'
        '5nc2V0dGluZ3MSjQEKK3Byb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eWF1dG9zY2FsaW5nc2V0dGlu'
        'Z3MYmvCmUCABKAsyKC5keW5hbW9kYi5BdXRvU2NhbGluZ1NldHRpbmdzRGVzY3JpcHRpb25SK3'
        'Byb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eWF1dG9zY2FsaW5nc2V0dGluZ3M=');

@$core.Deprecated(
    'Use replicaGlobalSecondaryIndexAutoScalingUpdateDescriptor instead')
const ReplicaGlobalSecondaryIndexAutoScalingUpdate$json = {
  '1': 'ReplicaGlobalSecondaryIndexAutoScalingUpdate',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'provisionedreadcapacityautoscalingupdate',
      '3': 22287551,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'provisionedreadcapacityautoscalingupdate'
    },
  ],
};

/// Descriptor for `ReplicaGlobalSecondaryIndexAutoScalingUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    replicaGlobalSecondaryIndexAutoScalingUpdateDescriptor =
    $convert.base64Decode(
        'CixSZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhBdXRvU2NhbGluZ1VwZGF0ZRIfCglpbmRleG'
        '5hbWUYkdXrMCABKAlSCWluZGV4bmFtZRKCAQoocHJvdmlzaW9uZWRyZWFkY2FwYWNpdHlhdXRv'
        'c2NhbGluZ3VwZGF0ZRi/qdAKIAEoCzIjLmR5bmFtb2RiLkF1dG9TY2FsaW5nU2V0dGluZ3NVcG'
        'RhdGVSKHByb3Zpc2lvbmVkcmVhZGNhcGFjaXR5YXV0b3NjYWxpbmd1cGRhdGU=');

@$core
    .Deprecated('Use replicaGlobalSecondaryIndexDescriptionDescriptor instead')
const ReplicaGlobalSecondaryIndexDescription$json = {
  '1': 'ReplicaGlobalSecondaryIndexDescription',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'ondemandthroughputoverride',
      '3': 317165234,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughputOverride',
      '10': 'ondemandthroughputoverride'
    },
    {
      '1': 'provisionedthroughputoverride',
      '3': 413332116,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughputOverride',
      '10': 'provisionedthroughputoverride'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndexWarmThroughputDescription',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `ReplicaGlobalSecondaryIndexDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaGlobalSecondaryIndexDescriptionDescriptor = $convert.base64Decode(
    'CiZSZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhEZXNjcmlwdGlvbhIfCglpbmRleG5hbWUYkd'
    'XrMCABKAlSCWluZGV4bmFtZRJoChpvbmRlbWFuZHRocm91Z2hwdXRvdmVycmlkZRiynZ6XASAB'
    'KAsyJC5keW5hbW9kYi5PbkRlbWFuZFRocm91Z2hwdXRPdmVycmlkZVIab25kZW1hbmR0aHJvdW'
    'docHV0b3ZlcnJpZGUScQodcHJvdmlzaW9uZWR0aHJvdWdocHV0b3ZlcnJpZGUYlOWLxQEgASgL'
    'MicuZHluYW1vZGIuUHJvdmlzaW9uZWRUaHJvdWdocHV0T3ZlcnJpZGVSHXByb3Zpc2lvbmVkdG'
    'hyb3VnaHB1dG92ZXJyaWRlEmMKDndhcm10aHJvdWdocHV0GIPeyIoBIAEoCzI3LmR5bmFtb2Ri'
    'Lkdsb2JhbFNlY29uZGFyeUluZGV4V2FybVRocm91Z2hwdXREZXNjcmlwdGlvblIOd2FybXRocm'
    '91Z2hwdXQ=');

@$core.Deprecated(
    'Use replicaGlobalSecondaryIndexSettingsDescriptionDescriptor instead')
const ReplicaGlobalSecondaryIndexSettingsDescription$json = {
  '1': 'ReplicaGlobalSecondaryIndexSettingsDescription',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'indexstatus',
      '3': 364436830,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.IndexStatus',
      '10': 'indexstatus'
    },
    {
      '1': 'provisionedreadcapacityautoscalingsettings',
      '3': 85617667,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'provisionedreadcapacityautoscalingsettings'
    },
    {
      '1': 'provisionedreadcapacityunits',
      '3': 350750021,
      '4': 1,
      '5': 3,
      '10': 'provisionedreadcapacityunits'
    },
    {
      '1': 'provisionedwritecapacityautoscalingsettings',
      '3': 168409114,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'provisionedwritecapacityautoscalingsettings'
    },
    {
      '1': 'provisionedwritecapacityunits',
      '3': 225881684,
      '4': 1,
      '5': 3,
      '10': 'provisionedwritecapacityunits'
    },
  ],
};

/// Descriptor for `ReplicaGlobalSecondaryIndexSettingsDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaGlobalSecondaryIndexSettingsDescriptionDescriptor = $convert.base64Decode(
    'Ci5SZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhTZXR0aW5nc0Rlc2NyaXB0aW9uEh8KCWluZG'
    'V4bmFtZRiR1eswIAEoCVIJaW5kZXhuYW1lEjsKC2luZGV4c3RhdHVzGN66460BIAEoDjIVLmR5'
    'bmFtb2RiLkluZGV4U3RhdHVzUgtpbmRleHN0YXR1cxKLAQoqcHJvdmlzaW9uZWRyZWFkY2FwYW'
    'NpdHlhdXRvc2NhbGluZ3NldHRpbmdzGIPY6SggASgLMiguZHluYW1vZGIuQXV0b1NjYWxpbmdT'
    'ZXR0aW5nc0Rlc2NyaXB0aW9uUipwcm92aXNpb25lZHJlYWRjYXBhY2l0eWF1dG9zY2FsaW5nc2'
    'V0dGluZ3MSRgoccHJvdmlzaW9uZWRyZWFkY2FwYWNpdHl1bml0cxjFiqCnASABKANSHHByb3Zp'
    'c2lvbmVkcmVhZGNhcGFjaXR5dW5pdHMSjQEKK3Byb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eWF1dG'
    '9zY2FsaW5nc2V0dGluZ3MYmvCmUCABKAsyKC5keW5hbW9kYi5BdXRvU2NhbGluZ1NldHRpbmdz'
    'RGVzY3JpcHRpb25SK3Byb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eWF1dG9zY2FsaW5nc2V0dGluZ3'
    'MSRwodcHJvdmlzaW9uZWR3cml0ZWNhcGFjaXR5dW5pdHMY1NzaayABKANSHXByb3Zpc2lvbmVk'
    'd3JpdGVjYXBhY2l0eXVuaXRz');

@$core.Deprecated(
    'Use replicaGlobalSecondaryIndexSettingsUpdateDescriptor instead')
const ReplicaGlobalSecondaryIndexSettingsUpdate$json = {
  '1': 'ReplicaGlobalSecondaryIndexSettingsUpdate',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'provisionedreadcapacityautoscalingsettingsupdate',
      '3': 57419900,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'provisionedreadcapacityautoscalingsettingsupdate'
    },
    {
      '1': 'provisionedreadcapacityunits',
      '3': 350750021,
      '4': 1,
      '5': 3,
      '10': 'provisionedreadcapacityunits'
    },
  ],
};

/// Descriptor for `ReplicaGlobalSecondaryIndexSettingsUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    replicaGlobalSecondaryIndexSettingsUpdateDescriptor = $convert.base64Decode(
        'CilSZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhTZXR0aW5nc1VwZGF0ZRIfCglpbmRleG5hbW'
        'UYkdXrMCABKAlSCWluZGV4bmFtZRKSAQowcHJvdmlzaW9uZWRyZWFkY2FwYWNpdHlhdXRvc2Nh'
        'bGluZ3NldHRpbmdzdXBkYXRlGPzQsBsgASgLMiMuZHluYW1vZGIuQXV0b1NjYWxpbmdTZXR0aW'
        '5nc1VwZGF0ZVIwcHJvdmlzaW9uZWRyZWFkY2FwYWNpdHlhdXRvc2NhbGluZ3NldHRpbmdzdXBk'
        'YXRlEkYKHHByb3Zpc2lvbmVkcmVhZGNhcGFjaXR5dW5pdHMYxYqgpwEgASgDUhxwcm92aXNpb2'
        '5lZHJlYWRjYXBhY2l0eXVuaXRz');

@$core.Deprecated('Use replicaNotFoundExceptionDescriptor instead')
const ReplicaNotFoundException$json = {
  '1': 'ReplicaNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ReplicaNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChhSZXBsaWNhTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use replicaSettingsDescriptionDescriptor instead')
const ReplicaSettingsDescription$json = {
  '1': 'ReplicaSettingsDescription',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'replicabillingmodesummary',
      '3': 456333796,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.BillingModeSummary',
      '10': 'replicabillingmodesummary'
    },
    {
      '1': 'replicaglobalsecondaryindexsettings',
      '3': 521234416,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaGlobalSecondaryIndexSettingsDescription',
      '10': 'replicaglobalsecondaryindexsettings'
    },
    {
      '1': 'replicaprovisionedreadcapacityautoscalingsettings',
      '3': 125131885,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'replicaprovisionedreadcapacityautoscalingsettings'
    },
    {
      '1': 'replicaprovisionedreadcapacityunits',
      '3': 82081083,
      '4': 1,
      '5': 3,
      '10': 'replicaprovisionedreadcapacityunits'
    },
    {
      '1': 'replicaprovisionedwritecapacityautoscalingsettings',
      '3': 293841796,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsDescription',
      '10': 'replicaprovisionedwritecapacityautoscalingsettings'
    },
    {
      '1': 'replicaprovisionedwritecapacityunits',
      '3': 356738858,
      '4': 1,
      '5': 3,
      '10': 'replicaprovisionedwritecapacityunits'
    },
    {
      '1': 'replicastatus',
      '3': 466739730,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReplicaStatus',
      '10': 'replicastatus'
    },
    {
      '1': 'replicatableclasssummary',
      '3': 519483614,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableClassSummary',
      '10': 'replicatableclasssummary'
    },
  ],
};

/// Descriptor for `ReplicaSettingsDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaSettingsDescriptionDescriptor = $convert.base64Decode(
    'ChpSZXBsaWNhU2V0dGluZ3NEZXNjcmlwdGlvbhIhCgpyZWdpb25uYW1lGL+buTUgASgJUgpyZW'
    'dpb25uYW1lEl4KGXJlcGxpY2FiaWxsaW5nbW9kZXN1bW1hcnkY5LPM2QEgASgLMhwuZHluYW1v'
    'ZGIuQmlsbGluZ01vZGVTdW1tYXJ5UhlyZXBsaWNhYmlsbGluZ21vZGVzdW1tYXJ5Eo4BCiNyZX'
    'BsaWNhZ2xvYmFsc2Vjb25kYXJ5aW5kZXhzZXR0aW5ncxjwz8X4ASADKAsyOC5keW5hbW9kYi5S'
    'ZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhTZXR0aW5nc0Rlc2NyaXB0aW9uUiNyZXBsaWNhZ2'
    'xvYmFsc2Vjb25kYXJ5aW5kZXhzZXR0aW5ncxKZAQoxcmVwbGljYXByb3Zpc2lvbmVkcmVhZGNh'
    'cGFjaXR5YXV0b3NjYWxpbmdzZXR0aW5ncxjtuNU7IAEoCzIoLmR5bmFtb2RiLkF1dG9TY2FsaW'
    '5nU2V0dGluZ3NEZXNjcmlwdGlvblIxcmVwbGljYXByb3Zpc2lvbmVkcmVhZGNhcGFjaXR5YXV0'
    'b3NjYWxpbmdzZXR0aW5ncxJTCiNyZXBsaWNhcHJvdmlzaW9uZWRyZWFkY2FwYWNpdHl1bml0cx'
    'i76pEnIAEoA1IjcmVwbGljYXByb3Zpc2lvbmVkcmVhZGNhcGFjaXR5dW5pdHMSnAEKMnJlcGxp'
    'Y2Fwcm92aXNpb25lZHdyaXRlY2FwYWNpdHlhdXRvc2NhbGluZ3NldHRpbmdzGITXjowBIAEoCz'
    'IoLmR5bmFtb2RiLkF1dG9TY2FsaW5nU2V0dGluZ3NEZXNjcmlwdGlvblIycmVwbGljYXByb3Zp'
    'c2lvbmVkd3JpdGVjYXBhY2l0eWF1dG9zY2FsaW5nc2V0dGluZ3MSVgokcmVwbGljYXByb3Zpc2'
    'lvbmVkd3JpdGVjYXBhY2l0eXVuaXRzGKrOjaoBIAEoA1IkcmVwbGljYXByb3Zpc2lvbmVkd3Jp'
    'dGVjYXBhY2l0eXVuaXRzEkEKDXJlcGxpY2FzdGF0dXMYksTH3gEgASgOMhcuZHluYW1vZGIuUm'
    'VwbGljYVN0YXR1c1INcmVwbGljYXN0YXR1cxJbChhyZXBsaWNhdGFibGVjbGFzc3N1bW1hcnkY'
    '3uHa9wEgASgLMhsuZHluYW1vZGIuVGFibGVDbGFzc1N1bW1hcnlSGHJlcGxpY2F0YWJsZWNsYX'
    'Nzc3VtbWFyeQ==');

@$core.Deprecated('Use replicaSettingsUpdateDescriptor instead')
const ReplicaSettingsUpdate$json = {
  '1': 'ReplicaSettingsUpdate',
  '2': [
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'replicaglobalsecondaryindexsettingsupdate',
      '3': 116195935,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaGlobalSecondaryIndexSettingsUpdate',
      '10': 'replicaglobalsecondaryindexsettingsupdate'
    },
    {
      '1': 'replicaprovisionedreadcapacityautoscalingsettingsupdate',
      '3': 245262702,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'replicaprovisionedreadcapacityautoscalingsettingsupdate'
    },
    {
      '1': 'replicaprovisionedreadcapacityunits',
      '3': 82081083,
      '4': 1,
      '5': 3,
      '10': 'replicaprovisionedreadcapacityunits'
    },
    {
      '1': 'replicatableclass',
      '3': 248679204,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableClass',
      '10': 'replicatableclass'
    },
  ],
};

/// Descriptor for `ReplicaSettingsUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaSettingsUpdateDescriptor = $convert.base64Decode(
    'ChVSZXBsaWNhU2V0dGluZ3NVcGRhdGUSIQoKcmVnaW9ubmFtZRi/m7k1IAEoCVIKcmVnaW9ubm'
    'FtZRKUAQopcmVwbGljYWdsb2JhbHNlY29uZGFyeWluZGV4c2V0dGluZ3N1cGRhdGUY34S0NyAD'
    'KAsyMy5keW5hbW9kYi5SZXBsaWNhR2xvYmFsU2Vjb25kYXJ5SW5kZXhTZXR0aW5nc1VwZGF0ZV'
    'IpcmVwbGljYWdsb2JhbHNlY29uZGFyeWluZGV4c2V0dGluZ3N1cGRhdGUSoAEKN3JlcGxpY2Fw'
    'cm92aXNpb25lZHJlYWRjYXBhY2l0eWF1dG9zY2FsaW5nc2V0dGluZ3N1cGRhdGUY7tL5dCABKA'
    'syIy5keW5hbW9kYi5BdXRvU2NhbGluZ1NldHRpbmdzVXBkYXRlUjdyZXBsaWNhcHJvdmlzaW9u'
    'ZWRyZWFkY2FwYWNpdHlhdXRvc2NhbGluZ3NldHRpbmdzdXBkYXRlElMKI3JlcGxpY2Fwcm92aX'
    'Npb25lZHJlYWRjYXBhY2l0eXVuaXRzGLvqkScgASgDUiNyZXBsaWNhcHJvdmlzaW9uZWRyZWFk'
    'Y2FwYWNpdHl1bml0cxJFChFyZXBsaWNhdGFibGVjbGFzcxiklsp2IAEoDjIULmR5bmFtb2RiLl'
    'RhYmxlQ2xhc3NSEXJlcGxpY2F0YWJsZWNsYXNz');

@$core.Deprecated('Use replicaUpdateDescriptor instead')
const ReplicaUpdate$json = {
  '1': 'ReplicaUpdate',
  '2': [
    {
      '1': 'create',
      '3': 420340862,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.CreateReplicaAction',
      '10': 'create'
    },
    {
      '1': 'delete',
      '3': 395831915,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.DeleteReplicaAction',
      '10': 'delete'
    },
  ],
};

/// Descriptor for `ReplicaUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaUpdateDescriptor = $convert.base64Decode(
    'Cg1SZXBsaWNhVXBkYXRlEjkKBmNyZWF0ZRj+yLfIASABKAsyHS5keW5hbW9kYi5DcmVhdGVSZX'
    'BsaWNhQWN0aW9uUgZjcmVhdGUSOQoGZGVsZXRlGOvU37wBIAEoCzIdLmR5bmFtb2RiLkRlbGV0'
    'ZVJlcGxpY2FBY3Rpb25SBmRlbGV0ZQ==');

@$core.Deprecated('Use replicatedWriteConflictExceptionDescriptor instead')
const ReplicatedWriteConflictException$json = {
  '1': 'ReplicatedWriteConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ReplicatedWriteConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicatedWriteConflictExceptionDescriptor =
    $convert.base64Decode(
        'CiBSZXBsaWNhdGVkV3JpdGVDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use replicationGroupUpdateDescriptor instead')
const ReplicationGroupUpdate$json = {
  '1': 'ReplicationGroupUpdate',
  '2': [
    {
      '1': 'create',
      '3': 420340862,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.CreateReplicationGroupMemberAction',
      '10': 'create'
    },
    {
      '1': 'delete',
      '3': 395831915,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.DeleteReplicationGroupMemberAction',
      '10': 'delete'
    },
    {
      '1': 'update',
      '3': 237178517,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.UpdateReplicationGroupMemberAction',
      '10': 'update'
    },
  ],
};

/// Descriptor for `ReplicationGroupUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationGroupUpdateDescriptor = $convert.base64Decode(
    'ChZSZXBsaWNhdGlvbkdyb3VwVXBkYXRlEkgKBmNyZWF0ZRj+yLfIASABKAsyLC5keW5hbW9kYi'
    '5DcmVhdGVSZXBsaWNhdGlvbkdyb3VwTWVtYmVyQWN0aW9uUgZjcmVhdGUSSAoGZGVsZXRlGOvU'
    '37wBIAEoCzIsLmR5bmFtb2RiLkRlbGV0ZVJlcGxpY2F0aW9uR3JvdXBNZW1iZXJBY3Rpb25SBm'
    'RlbGV0ZRJHCgZ1cGRhdGUYlZ2McSABKAsyLC5keW5hbW9kYi5VcGRhdGVSZXBsaWNhdGlvbkdy'
    'b3VwTWVtYmVyQWN0aW9uUgZ1cGRhdGU=');

@$core.Deprecated('Use requestLimitExceededDescriptor instead')
const RequestLimitExceeded$json = {
  '1': 'RequestLimitExceeded',
  '2': [
    {
      '1': 'throttlingreasons',
      '3': 284436852,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ThrottlingReason',
      '10': 'throttlingreasons'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RequestLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestLimitExceededDescriptor = $convert.base64Decode(
    'ChRSZXF1ZXN0TGltaXRFeGNlZWRlZBJMChF0aHJvdHRsaW5ncmVhc29ucxj00tCHASADKAsyGi'
    '5keW5hbW9kYi5UaHJvdHRsaW5nUmVhc29uUhF0aHJvdHRsaW5ncmVhc29ucxIbCgdtZXNzYWdl'
    'GOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use resourceInUseExceptionDescriptor instead')
const ResourceInUseException$json = {
  '1': 'ResourceInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChZSZXNvdXJjZUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

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

@$core.Deprecated('Use restoreSummaryDescriptor instead')
const RestoreSummary$json = {
  '1': 'RestoreSummary',
  '2': [
    {
      '1': 'restoredatetime',
      '3': 176772301,
      '4': 1,
      '5': 9,
      '10': 'restoredatetime'
    },
    {
      '1': 'restoreinprogress',
      '3': 508499598,
      '4': 1,
      '5': 8,
      '10': 'restoreinprogress'
    },
    {
      '1': 'sourcebackuparn',
      '3': 16805842,
      '4': 1,
      '5': 9,
      '10': 'sourcebackuparn'
    },
    {
      '1': 'sourcetablearn',
      '3': 124215796,
      '4': 1,
      '5': 9,
      '10': 'sourcetablearn'
    },
  ],
};

/// Descriptor for `RestoreSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreSummaryDescriptor = $convert.base64Decode(
    'Cg5SZXN0b3JlU3VtbWFyeRIrCg9yZXN0b3JlZGF0ZXRpbWUYzamlVCABKAlSD3Jlc3RvcmVkYX'
    'RldGltZRIwChFyZXN0b3JlaW5wcm9ncmVzcxiOrbzyASABKAhSEXJlc3RvcmVpbnByb2dyZXNz'
    'EisKD3NvdXJjZWJhY2t1cGFybhjS34EIIAEoCVIPc291cmNlYmFja3VwYXJuEikKDnNvdXJjZX'
    'RhYmxlYXJuGPTDnTsgASgJUg5zb3VyY2V0YWJsZWFybg==');

@$core.Deprecated('Use restoreTableFromBackupInputDescriptor instead')
const RestoreTableFromBackupInput$json = {
  '1': 'RestoreTableFromBackupInput',
  '2': [
    {'1': 'backuparn', '3': 370874339, '4': 1, '5': 9, '10': 'backuparn'},
    {
      '1': 'billingmodeoverride',
      '3': 357784560,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'billingmodeoverride'
    },
    {
      '1': 'globalsecondaryindexoverride',
      '3': 369844021,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndex',
      '10': 'globalsecondaryindexoverride'
    },
    {
      '1': 'localsecondaryindexoverride',
      '3': 431784607,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.LocalSecondaryIndex',
      '10': 'localsecondaryindexoverride'
    },
    {
      '1': 'ondemandthroughputoverride',
      '3': 317165234,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughputoverride'
    },
    {
      '1': 'provisionedthroughputoverride',
      '3': 413332116,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughputoverride'
    },
    {
      '1': 'ssespecificationoverride',
      '3': 421570884,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SSESpecification',
      '10': 'ssespecificationoverride'
    },
    {
      '1': 'targettablename',
      '3': 298767720,
      '4': 1,
      '5': 9,
      '10': 'targettablename'
    },
  ],
};

/// Descriptor for `RestoreTableFromBackupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreTableFromBackupInputDescriptor = $convert.base64Decode(
    'ChtSZXN0b3JlVGFibGVGcm9tQmFja3VwSW5wdXQSIAoJYmFja3VwYXJuGOOv7LABIAEoCVIJYm'
    'Fja3VwYXJuEksKE2JpbGxpbmdtb2Rlb3ZlcnJpZGUY8LfNqgEgASgOMhUuZHluYW1vZGIuQmls'
    'bGluZ01vZGVSE2JpbGxpbmdtb2Rlb3ZlcnJpZGUSZgocZ2xvYmFsc2Vjb25kYXJ5aW5kZXhvdm'
    'VycmlkZRi1vq2wASADKAsyHi5keW5hbW9kYi5HbG9iYWxTZWNvbmRhcnlJbmRleFIcZ2xvYmFs'
    'c2Vjb25kYXJ5aW5kZXhvdmVycmlkZRJjChtsb2NhbHNlY29uZGFyeWluZGV4b3ZlcnJpZGUYn4'
    'XyzQEgAygLMh0uZHluYW1vZGIuTG9jYWxTZWNvbmRhcnlJbmRleFIbbG9jYWxzZWNvbmRhcnlp'
    'bmRleG92ZXJyaWRlEmAKGm9uZGVtYW5kdGhyb3VnaHB1dG92ZXJyaWRlGLKdnpcBIAEoCzIcLm'
    'R5bmFtb2RiLk9uRGVtYW5kVGhyb3VnaHB1dFIab25kZW1hbmR0aHJvdWdocHV0b3ZlcnJpZGUS'
    'aQodcHJvdmlzaW9uZWR0aHJvdWdocHV0b3ZlcnJpZGUYlOWLxQEgASgLMh8uZHluYW1vZGIuUH'
    'JvdmlzaW9uZWRUaHJvdWdocHV0Uh1wcm92aXNpb25lZHRocm91Z2hwdXRvdmVycmlkZRJaChhz'
    'c2VzcGVjaWZpY2F0aW9ub3ZlcnJpZGUYxNKCyQEgASgLMhouZHluYW1vZGIuU1NFU3BlY2lmaW'
    'NhdGlvblIYc3Nlc3BlY2lmaWNhdGlvbm92ZXJyaWRlEiwKD3RhcmdldHRhYmxlbmFtZRjoqruO'
    'ASABKAlSD3RhcmdldHRhYmxlbmFtZQ==');

@$core.Deprecated('Use restoreTableFromBackupOutputDescriptor instead')
const RestoreTableFromBackupOutput$json = {
  '1': 'RestoreTableFromBackupOutput',
  '2': [
    {
      '1': 'tabledescription',
      '3': 19962388,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableDescription',
      '10': 'tabledescription'
    },
  ],
};

/// Descriptor for `RestoreTableFromBackupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreTableFromBackupOutputDescriptor =
    $convert.base64Decode(
        'ChxSZXN0b3JlVGFibGVGcm9tQmFja3VwT3V0cHV0EkkKEHRhYmxlZGVzY3JpcHRpb24YlLTCCS'
        'ABKAsyGi5keW5hbW9kYi5UYWJsZURlc2NyaXB0aW9uUhB0YWJsZWRlc2NyaXB0aW9u');

@$core.Deprecated('Use restoreTableToPointInTimeInputDescriptor instead')
const RestoreTableToPointInTimeInput$json = {
  '1': 'RestoreTableToPointInTimeInput',
  '2': [
    {
      '1': 'billingmodeoverride',
      '3': 357784560,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'billingmodeoverride'
    },
    {
      '1': 'globalsecondaryindexoverride',
      '3': 369844021,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndex',
      '10': 'globalsecondaryindexoverride'
    },
    {
      '1': 'localsecondaryindexoverride',
      '3': 431784607,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.LocalSecondaryIndex',
      '10': 'localsecondaryindexoverride'
    },
    {
      '1': 'ondemandthroughputoverride',
      '3': 317165234,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughputoverride'
    },
    {
      '1': 'provisionedthroughputoverride',
      '3': 413332116,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughputoverride'
    },
    {
      '1': 'restoredatetime',
      '3': 176772301,
      '4': 1,
      '5': 9,
      '10': 'restoredatetime'
    },
    {
      '1': 'ssespecificationoverride',
      '3': 421570884,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SSESpecification',
      '10': 'ssespecificationoverride'
    },
    {
      '1': 'sourcetablearn',
      '3': 124215796,
      '4': 1,
      '5': 9,
      '10': 'sourcetablearn'
    },
    {
      '1': 'sourcetablename',
      '3': 165368952,
      '4': 1,
      '5': 9,
      '10': 'sourcetablename'
    },
    {
      '1': 'targettablename',
      '3': 298767720,
      '4': 1,
      '5': 9,
      '10': 'targettablename'
    },
    {
      '1': 'uselatestrestorabletime',
      '3': 434512618,
      '4': 1,
      '5': 8,
      '10': 'uselatestrestorabletime'
    },
  ],
};

/// Descriptor for `RestoreTableToPointInTimeInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreTableToPointInTimeInputDescriptor = $convert.base64Decode(
    'Ch5SZXN0b3JlVGFibGVUb1BvaW50SW5UaW1lSW5wdXQSSwoTYmlsbGluZ21vZGVvdmVycmlkZR'
    'jwt82qASABKA4yFS5keW5hbW9kYi5CaWxsaW5nTW9kZVITYmlsbGluZ21vZGVvdmVycmlkZRJm'
    'ChxnbG9iYWxzZWNvbmRhcnlpbmRleG92ZXJyaWRlGLW+rbABIAMoCzIeLmR5bmFtb2RiLkdsb2'
    'JhbFNlY29uZGFyeUluZGV4UhxnbG9iYWxzZWNvbmRhcnlpbmRleG92ZXJyaWRlEmMKG2xvY2Fs'
    'c2Vjb25kYXJ5aW5kZXhvdmVycmlkZRifhfLNASADKAsyHS5keW5hbW9kYi5Mb2NhbFNlY29uZG'
    'FyeUluZGV4Uhtsb2NhbHNlY29uZGFyeWluZGV4b3ZlcnJpZGUSYAoab25kZW1hbmR0aHJvdWdo'
    'cHV0b3ZlcnJpZGUYsp2elwEgASgLMhwuZHluYW1vZGIuT25EZW1hbmRUaHJvdWdocHV0Uhpvbm'
    'RlbWFuZHRocm91Z2hwdXRvdmVycmlkZRJpCh1wcm92aXNpb25lZHRocm91Z2hwdXRvdmVycmlk'
    'ZRiU5YvFASABKAsyHy5keW5hbW9kYi5Qcm92aXNpb25lZFRocm91Z2hwdXRSHXByb3Zpc2lvbm'
    'VkdGhyb3VnaHB1dG92ZXJyaWRlEisKD3Jlc3RvcmVkYXRldGltZRjNqaVUIAEoCVIPcmVzdG9y'
    'ZWRhdGV0aW1lEloKGHNzZXNwZWNpZmljYXRpb25vdmVycmlkZRjE0oLJASABKAsyGi5keW5hbW'
    '9kYi5TU0VTcGVjaWZpY2F0aW9uUhhzc2VzcGVjaWZpY2F0aW9ub3ZlcnJpZGUSKQoOc291cmNl'
    'dGFibGVhcm4Y9MOdOyABKAlSDnNvdXJjZXRhYmxlYXJuEisKD3NvdXJjZXRhYmxlbmFtZRj4qO'
    '1OIAEoCVIPc291cmNldGFibGVuYW1lEiwKD3RhcmdldHRhYmxlbmFtZRjoqruOASABKAlSD3Rh'
    'cmdldHRhYmxlbmFtZRI8Chd1c2VsYXRlc3RyZXN0b3JhYmxldGltZRjqxZjPASABKAhSF3VzZW'
    'xhdGVzdHJlc3RvcmFibGV0aW1l');

@$core.Deprecated('Use restoreTableToPointInTimeOutputDescriptor instead')
const RestoreTableToPointInTimeOutput$json = {
  '1': 'RestoreTableToPointInTimeOutput',
  '2': [
    {
      '1': 'tabledescription',
      '3': 19962388,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableDescription',
      '10': 'tabledescription'
    },
  ],
};

/// Descriptor for `RestoreTableToPointInTimeOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreTableToPointInTimeOutputDescriptor =
    $convert.base64Decode(
        'Ch9SZXN0b3JlVGFibGVUb1BvaW50SW5UaW1lT3V0cHV0EkkKEHRhYmxlZGVzY3JpcHRpb24YlL'
        'TCCSABKAsyGi5keW5hbW9kYi5UYWJsZURlc2NyaXB0aW9uUhB0YWJsZWRlc2NyaXB0aW9u');

@$core.Deprecated('Use s3BucketSourceDescriptor instead')
const S3BucketSource$json = {
  '1': 'S3BucketSource',
  '2': [
    {'1': 's3bucket', '3': 114031434, '4': 1, '5': 9, '10': 's3bucket'},
    {
      '1': 's3bucketowner',
      '3': 351576129,
      '4': 1,
      '5': 9,
      '10': 's3bucketowner'
    },
    {'1': 's3keyprefix', '3': 206015359, '4': 1, '5': 9, '10': 's3keyprefix'},
  ],
};

/// Descriptor for `S3BucketSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3BucketSourceDescriptor = $convert.base64Decode(
    'Cg5TM0J1Y2tldFNvdXJjZRIdCghzM2J1Y2tldBjK9q82IAEoCVIIczNidWNrZXQSKAoNczNidW'
    'NrZXRvd25lchjBwNKnASABKAlSDXMzYnVja2V0b3duZXISIwoLczNrZXlwcmVmaXgY/5aeYiAB'
    'KAlSC3Mza2V5cHJlZml4');

@$core.Deprecated('Use sSEDescriptionDescriptor instead')
const SSEDescription$json = {
  '1': 'SSEDescription',
  '2': [
    {
      '1': 'inaccessibleencryptiondatetime',
      '3': 478188797,
      '4': 1,
      '5': 9,
      '10': 'inaccessibleencryptiondatetime'
    },
    {
      '1': 'kmsmasterkeyarn',
      '3': 281937987,
      '4': 1,
      '5': 9,
      '10': 'kmsmasterkeyarn'
    },
    {
      '1': 'ssetype',
      '3': 431846435,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.SSEType',
      '10': 'ssetype'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.SSEStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `SSEDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sSEDescriptionDescriptor = $convert.base64Decode(
    'Cg5TU0VEZXNjcmlwdGlvbhJKCh5pbmFjY2Vzc2libGVlbmNyeXB0aW9uZGF0ZXRpbWUY/amC5A'
    'EgASgJUh5pbmFjY2Vzc2libGVlbmNyeXB0aW9uZGF0ZXRpbWUSLAoPa21zbWFzdGVya2V5YXJu'
    'GMOQuIYBIAEoCVIPa21zbWFzdGVya2V5YXJuEi8KB3NzZXR5cGUYo+j1zQEgASgOMhEuZHluYW'
    '1vZGIuU1NFVHlwZVIHc3NldHlwZRIuCgZzdGF0dXMYkOT7AiABKA4yEy5keW5hbW9kYi5TU0VT'
    'dGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use sSESpecificationDescriptor instead')
const SSESpecification$json = {
  '1': 'SSESpecification',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'kmsmasterkeyid',
      '3': 521459443,
      '4': 1,
      '5': 9,
      '10': 'kmsmasterkeyid'
    },
    {
      '1': 'ssetype',
      '3': 431846435,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.SSEType',
      '10': 'ssetype'
    },
  ],
};

/// Descriptor for `SSESpecification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sSESpecificationDescriptor = $convert.base64Decode(
    'ChBTU0VTcGVjaWZpY2F0aW9uEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVkEioKDmttc2'
    '1hc3RlcmtleWlkGPOt0/gBIAEoCVIOa21zbWFzdGVya2V5aWQSLwoHc3NldHlwZRij6PXNASAB'
    'KA4yES5keW5hbW9kYi5TU0VUeXBlUgdzc2V0eXBl');

@$core.Deprecated('Use scanInputDescriptor instead')
const ScanInput$json = {
  '1': 'ScanInput',
  '2': [
    {
      '1': 'attributestoget',
      '3': 311382592,
      '4': 3,
      '5': 9,
      '10': 'attributestoget'
    },
    {
      '1': 'conditionaloperator',
      '3': 172066260,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ConditionalOperator',
      '10': 'conditionaloperator'
    },
    {
      '1': 'consistentread',
      '3': 531556994,
      '4': 1,
      '5': 8,
      '10': 'consistentread'
    },
    {
      '1': 'exclusivestartkey',
      '3': 348137297,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ScanInput.ExclusivestartkeyEntry',
      '10': 'exclusivestartkey'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ScanInput.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ScanInput.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'filterexpression',
      '3': 68019610,
      '4': 1,
      '5': 9,
      '10': 'filterexpression'
    },
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'projectionexpression',
      '3': 150730243,
      '4': 1,
      '5': 9,
      '10': 'projectionexpression'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'scanfilter',
      '3': 272885755,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ScanInput.ScanfilterEntry',
      '10': 'scanfilter'
    },
    {'1': 'segment', '3': 279654279, '4': 1, '5': 5, '10': 'segment'},
    {
      '1': 'select',
      '3': 512305998,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.Select',
      '10': 'select'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'totalsegments',
      '3': 149136904,
      '4': 1,
      '5': 5,
      '10': 'totalsegments'
    },
  ],
  '3': [
    ScanInput_ExclusivestartkeyEntry$json,
    ScanInput_ExpressionattributenamesEntry$json,
    ScanInput_ExpressionattributevaluesEntry$json,
    ScanInput_ScanfilterEntry$json
  ],
};

@$core.Deprecated('Use scanInputDescriptor instead')
const ScanInput_ExclusivestartkeyEntry$json = {
  '1': 'ExclusivestartkeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use scanInputDescriptor instead')
const ScanInput_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use scanInputDescriptor instead')
const ScanInput_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use scanInputDescriptor instead')
const ScanInput_ScanfilterEntry$json = {
  '1': 'ScanfilterEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Condition',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ScanInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scanInputDescriptor = $convert.base64Decode(
    'CglTY2FuSW5wdXQSLAoPYXR0cmlidXRlc3RvZ2V0GMCkvZQBIAMoCVIPYXR0cmlidXRlc3RvZ2'
    'V0ElIKE2NvbmRpdGlvbmFsb3BlcmF0b3IY1IuGUiABKA4yHS5keW5hbW9kYi5Db25kaXRpb25h'
    'bE9wZXJhdG9yUhNjb25kaXRpb25hbG9wZXJhdG9yEioKDmNvbnNpc3RlbnRyZWFkGILVu/0BIA'
    'EoCFIOY29uc2lzdGVudHJlYWQSXAoRZXhjbHVzaXZlc3RhcnRrZXkY0c6ApgEgAygLMiouZHlu'
    'YW1vZGIuU2NhbklucHV0LkV4Y2x1c2l2ZXN0YXJ0a2V5RW50cnlSEWV4Y2x1c2l2ZXN0YXJ0a2'
    'V5EnAKGGV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lcxj8mNFHIAMoCzIxLmR5bmFtb2RiLlNjYW5J'
    'bnB1dC5FeHByZXNzaW9uYXR0cmlidXRlbmFtZXNFbnRyeVIYZXhwcmVzc2lvbmF0dHJpYnV0ZW'
    '5hbWVzEnQKGWV4cHJlc3Npb25hdHRyaWJ1dGV2YWx1ZXMY2Jyg5wEgAygLMjIuZHluYW1vZGIu'
    'U2NhbklucHV0LkV4cHJlc3Npb25hdHRyaWJ1dGV2YWx1ZXNFbnRyeVIZZXhwcmVzc2lvbmF0dH'
    'JpYnV0ZXZhbHVlcxItChBmaWx0ZXJleHByZXNzaW9uGJrLtyAgASgJUhBmaWx0ZXJleHByZXNz'
    'aW9uEh8KCWluZGV4bmFtZRiR1eswIAEoCVIJaW5kZXhuYW1lEhgKBWxpbWl0GNWV2cQBIAEoBV'
    'IFbGltaXQSNQoUcHJvamVjdGlvbmV4cHJlc3Npb24Yg+zvRyABKAlSFHByb2plY3Rpb25leHBy'
    'ZXNzaW9uElsKFnJldHVybmNvbnN1bWVkY2FwYWNpdHkY/ufhFCABKA4yIC5keW5hbW9kYi5SZX'
    'R1cm5Db25zdW1lZENhcGFjaXR5UhZyZXR1cm5jb25zdW1lZGNhcGFjaXR5EkcKCnNjYW5maWx0'
    'ZXIY+8+PggEgAygLMiMuZHluYW1vZGIuU2NhbklucHV0LlNjYW5maWx0ZXJFbnRyeVIKc2Nhbm'
    'ZpbHRlchIcCgdzZWdtZW50GIffrIUBIAEoBVIHc2VnbWVudBIsCgZzZWxlY3QYztak9AEgASgO'
    'MhAuZHluYW1vZGIuU2VsZWN0UgZzZWxlY3QSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibG'
    'VuYW1lEicKDXRvdGFsc2VnbWVudHMYiMyORyABKAVSDXRvdGFsc2VnbWVudHMaXgoWRXhjbHVz'
    'aXZlc3RhcnRrZXlFbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1ZRgCIAEoCzIYLmR5bm'
    'Ftb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAEaSwodRXhwcmVzc2lvbmF0dHJpYnV0ZW5h'
    'bWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpmCh'
    '5FeHByZXNzaW9uYXR0cmlidXRldmFsdWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSLgoFdmFs'
    'dWUYAiABKAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZVIFdmFsdWU6AjgBGlIKD1NjYW5maW'
    'x0ZXJFbnRyeRIQCgNrZXkYASABKAlSA2tleRIpCgV2YWx1ZRgCIAEoCzITLmR5bmFtb2RiLkNv'
    'bmRpdGlvblIFdmFsdWU6AjgB');

@$core.Deprecated('Use scanOutputDescriptor instead')
const ScanOutput$json = {
  '1': 'ScanOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {'1': 'count', '3': 31963285, '4': 1, '5': 5, '10': 'count'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ScanOutput.ItemsEntry',
      '10': 'items'
    },
    {
      '1': 'lastevaluatedkey',
      '3': 54319830,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ScanOutput.LastevaluatedkeyEntry',
      '10': 'lastevaluatedkey'
    },
    {'1': 'scannedcount', '3': 531161315, '4': 1, '5': 5, '10': 'scannedcount'},
  ],
  '3': [ScanOutput_ItemsEntry$json, ScanOutput_LastevaluatedkeyEntry$json],
};

@$core.Deprecated('Use scanOutputDescriptor instead')
const ScanOutput_ItemsEntry$json = {
  '1': 'ItemsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use scanOutputDescriptor instead')
const ScanOutput_LastevaluatedkeyEntry$json = {
  '1': 'LastevaluatedkeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `ScanOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scanOutputDescriptor = $convert.base64Decode(
    'CgpTY2FuT3V0cHV0EkoKEGNvbnN1bWVkY2FwYWNpdHkYrKqh1gEgASgLMhouZHluYW1vZGIuQ2'
    '9uc3VtZWRDYXBhY2l0eVIQY29uc3VtZWRjYXBhY2l0eRIXCgVjb3VudBiV8Z4PIAEoBVIFY291'
    'bnQSOAoFaXRlbXMYsPDYASADKAsyHy5keW5hbW9kYi5TY2FuT3V0cHV0Lkl0ZW1zRW50cnlSBW'
    'l0ZW1zElkKEGxhc3RldmFsdWF0ZWRrZXkY1rXzGSADKAsyKi5keW5hbW9kYi5TY2FuT3V0cHV0'
    'Lkxhc3RldmFsdWF0ZWRrZXlFbnRyeVIQbGFzdGV2YWx1YXRlZGtleRImCgxzY2FubmVkY291bn'
    'QY48Gj/QEgASgFUgxzY2FubmVkY291bnQaUgoKSXRlbXNFbnRyeRIQCgNrZXkYASABKAlSA2tl'
    'eRIuCgV2YWx1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAEaXQ'
    'oVTGFzdGV2YWx1YXRlZGtleUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgL'
    'MhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use sourceTableDetailsDescriptor instead')
const SourceTableDetails$json = {
  '1': 'SourceTableDetails',
  '2': [
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'billingmode'
    },
    {'1': 'itemcount', '3': 26280022, '4': 1, '5': 3, '10': 'itemcount'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {
      '1': 'tablecreationdatetime',
      '3': 152706380,
      '4': 1,
      '5': 9,
      '10': 'tablecreationdatetime'
    },
    {'1': 'tableid', '3': 449893011, '4': 1, '5': 9, '10': 'tableid'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'tablesizebytes',
      '3': 220631294,
      '4': 1,
      '5': 3,
      '10': 'tablesizebytes'
    },
  ],
};

/// Descriptor for `SourceTableDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sourceTableDetailsDescriptor = $convert.base64Decode(
    'ChJTb3VyY2VUYWJsZURldGFpbHMSOgoLYmlsbGluZ21vZGUYwLToVyABKA4yFS5keW5hbW9kYi'
    '5CaWxsaW5nTW9kZVILYmlsbGluZ21vZGUSHwoJaXRlbWNvdW50GNaAxAwgASgDUglpdGVtY291'
    'bnQSPAoJa2V5c2NoZW1hGOjP3YsBIAMoCzIaLmR5bmFtb2RiLktleVNjaGVtYUVsZW1lbnRSCW'
    'tleXNjaGVtYRJQChJvbmRlbWFuZHRocm91Z2hwdXQYgt7a5QEgASgLMhwuZHluYW1vZGIuT25E'
    'ZW1hbmRUaHJvdWdocHV0UhJvbmRlbWFuZHRocm91Z2hwdXQSVwoVcHJvdmlzaW9uZWR0aHJvdW'
    'docHV0GIyjayABKAsyHy5keW5hbW9kYi5Qcm92aXNpb25lZFRocm91Z2hwdXRSFXByb3Zpc2lv'
    'bmVkdGhyb3VnaHB1dBIeCgh0YWJsZWFybhjjgOvNASABKAlSCHRhYmxlYXJuEjcKFXRhYmxlY3'
    'JlYXRpb25kYXRldGltZRjMuuhIIAEoCVIVdGFibGVjcmVhdGlvbmRhdGV0aW1lEhwKB3RhYmxl'
    'aWQYk6XD1gEgASgJUgd0YWJsZWlkEiAKCXRhYmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZR'
    'IpCg50YWJsZXNpemVieXRlcxj+oZppIAEoA1IOdGFibGVzaXplYnl0ZXM=');

@$core.Deprecated('Use sourceTableFeatureDetailsDescriptor instead')
const SourceTableFeatureDetails$json = {
  '1': 'SourceTableFeatureDetails',
  '2': [
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndexInfo',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'localsecondaryindexes',
      '3': 362339959,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.LocalSecondaryIndexInfo',
      '10': 'localsecondaryindexes'
    },
    {
      '1': 'ssedescription',
      '3': 350068773,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SSEDescription',
      '10': 'ssedescription'
    },
    {
      '1': 'streamdescription',
      '3': 363737034,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.StreamSpecification',
      '10': 'streamdescription'
    },
    {
      '1': 'timetolivedescription',
      '3': 367590876,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TimeToLiveDescription',
      '10': 'timetolivedescription'
    },
  ],
};

/// Descriptor for `SourceTableFeatureDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sourceTableFeatureDetailsDescriptor = $convert.base64Decode(
    'ChlTb3VyY2VUYWJsZUZlYXR1cmVEZXRhaWxzEl4KFmdsb2JhbHNlY29uZGFyeWluZGV4ZXMYqf'
    'qMwwEgAygLMiIuZHluYW1vZGIuR2xvYmFsU2Vjb25kYXJ5SW5kZXhJbmZvUhZnbG9iYWxzZWNv'
    'bmRhcnlpbmRleGVzElsKFWxvY2Fsc2Vjb25kYXJ5aW5kZXhlcxj3vOOsASADKAsyIS5keW5hbW'
    '9kYi5Mb2NhbFNlY29uZGFyeUluZGV4SW5mb1IVbG9jYWxzZWNvbmRhcnlpbmRleGVzEkQKDnNz'
    'ZWRlc2NyaXB0aW9uGKXA9qYBIAEoCzIYLmR5bmFtb2RiLlNTRURlc2NyaXB0aW9uUg5zc2VkZX'
    'NjcmlwdGlvbhJPChFzdHJlYW1kZXNjcmlwdGlvbhjK37itASABKAsyHS5keW5hbW9kYi5TdHJl'
    'YW1TcGVjaWZpY2F0aW9uUhFzdHJlYW1kZXNjcmlwdGlvbhJZChV0aW1ldG9saXZlZGVzY3JpcH'
    'Rpb24Y3PujrwEgASgLMh8uZHluYW1vZGIuVGltZVRvTGl2ZURlc2NyaXB0aW9uUhV0aW1ldG9s'
    'aXZlZGVzY3JpcHRpb24=');

@$core.Deprecated('Use streamSpecificationDescriptor instead')
const StreamSpecification$json = {
  '1': 'StreamSpecification',
  '2': [
    {
      '1': 'streamenabled',
      '3': 266707711,
      '4': 1,
      '5': 8,
      '10': 'streamenabled'
    },
    {
      '1': 'streamviewtype',
      '3': 380488241,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.StreamViewType',
      '10': 'streamviewtype'
    },
  ],
};

/// Descriptor for `StreamSpecification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamSpecificationDescriptor = $convert.base64Decode(
    'ChNTdHJlYW1TcGVjaWZpY2F0aW9uEicKDXN0cmVhbWVuYWJsZWQY/8WWfyABKAhSDXN0cmVhbW'
    'VuYWJsZWQSRAoOc3RyZWFtdmlld3R5cGUYsZS3tQEgASgOMhguZHluYW1vZGIuU3RyZWFtVmll'
    'd1R5cGVSDnN0cmVhbXZpZXd0eXBl');

@$core.Deprecated('Use tableAlreadyExistsExceptionDescriptor instead')
const TableAlreadyExistsException$json = {
  '1': 'TableAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TableAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChtUYWJsZUFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use tableAutoScalingDescriptionDescriptor instead')
const TableAutoScalingDescription$json = {
  '1': 'TableAutoScalingDescription',
  '2': [
    {
      '1': 'replicas',
      '3': 306066781,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaAutoScalingDescription',
      '10': 'replicas'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'tablestatus',
      '3': 207908810,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableStatus',
      '10': 'tablestatus'
    },
  ],
};

/// Descriptor for `TableAutoScalingDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableAutoScalingDescriptionDescriptor = $convert.base64Decode(
    'ChtUYWJsZUF1dG9TY2FsaW5nRGVzY3JpcHRpb24SRwoIcmVwbGljYXMY3er4kQEgAygLMicuZH'
    'luYW1vZGIuUmVwbGljYUF1dG9TY2FsaW5nRGVzY3JpcHRpb25SCHJlcGxpY2FzEiAKCXRhYmxl'
    'bmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZRI6Cgt0YWJsZXN0YXR1cxjK35FjIAEoDjIVLmR5bm'
    'Ftb2RiLlRhYmxlU3RhdHVzUgt0YWJsZXN0YXR1cw==');

@$core.Deprecated('Use tableClassSummaryDescriptor instead')
const TableClassSummary$json = {
  '1': 'TableClassSummary',
  '2': [
    {
      '1': 'lastupdatedatetime',
      '3': 452274318,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedatetime'
    },
    {
      '1': 'tableclass',
      '3': 342890498,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableClass',
      '10': 'tableclass'
    },
  ],
};

/// Descriptor for `TableClassSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableClassSummaryDescriptor = $convert.base64Decode(
    'ChFUYWJsZUNsYXNzU3VtbWFyeRIyChJsYXN0dXBkYXRlZGF0ZXRpbWUYjtHU1wEgASgJUhJsYX'
    'N0dXBkYXRlZGF0ZXRpbWUSOAoKdGFibGVjbGFzcxiCsMCjASABKA4yFC5keW5hbW9kYi5UYWJs'
    'ZUNsYXNzUgp0YWJsZWNsYXNz');

@$core.Deprecated('Use tableCreationParametersDescriptor instead')
const TableCreationParameters$json = {
  '1': 'TableCreationParameters',
  '2': [
    {
      '1': 'attributedefinitions',
      '3': 414687108,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeDefinition',
      '10': 'attributedefinitions'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndex',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'ssespecification',
      '3': 31692444,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SSESpecification',
      '10': 'ssespecification'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `TableCreationParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableCreationParametersDescriptor = $convert.base64Decode(
    'ChdUYWJsZUNyZWF0aW9uUGFyYW1ldGVycxJVChRhdHRyaWJ1dGVkZWZpbml0aW9ucxiEv97FAS'
    'ADKAsyHS5keW5hbW9kYi5BdHRyaWJ1dGVEZWZpbml0aW9uUhRhdHRyaWJ1dGVkZWZpbml0aW9u'
    'cxI6CgtiaWxsaW5nbW9kZRjAtOhXIAEoDjIVLmR5bmFtb2RiLkJpbGxpbmdNb2RlUgtiaWxsaW'
    '5nbW9kZRJaChZnbG9iYWxzZWNvbmRhcnlpbmRleGVzGKn6jMMBIAMoCzIeLmR5bmFtb2RiLkds'
    'b2JhbFNlY29uZGFyeUluZGV4UhZnbG9iYWxzZWNvbmRhcnlpbmRleGVzEjwKCWtleXNjaGVtYR'
    'joz92LASADKAsyGi5keW5hbW9kYi5LZXlTY2hlbWFFbGVtZW50UglrZXlzY2hlbWESUAoSb25k'
    'ZW1hbmR0aHJvdWdocHV0GILe2uUBIAEoCzIcLmR5bmFtb2RiLk9uRGVtYW5kVGhyb3VnaHB1dF'
    'ISb25kZW1hbmR0aHJvdWdocHV0ElcKFXByb3Zpc2lvbmVkdGhyb3VnaHB1dBiMo2sgASgLMh8u'
    'ZHluYW1vZGIuUHJvdmlzaW9uZWRUaHJvdWdocHV0UhVwcm92aXNpb25lZHRocm91Z2hwdXQSSQ'
    'oQc3Nlc3BlY2lmaWNhdGlvbhicrY4PIAEoCzIaLmR5bmFtb2RiLlNTRVNwZWNpZmljYXRpb25S'
    'EHNzZXNwZWNpZmljYXRpb24SIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use tableDescriptionDescriptor instead')
const TableDescription$json = {
  '1': 'TableDescription',
  '2': [
    {
      '1': 'archivalsummary',
      '3': 52039658,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ArchivalSummary',
      '10': 'archivalsummary'
    },
    {
      '1': 'attributedefinitions',
      '3': 414687108,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeDefinition',
      '10': 'attributedefinitions'
    },
    {
      '1': 'billingmodesummary',
      '3': 163529882,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.BillingModeSummary',
      '10': 'billingmodesummary'
    },
    {
      '1': 'creationdatetime',
      '3': 48904698,
      '4': 1,
      '5': 9,
      '10': 'creationdatetime'
    },
    {
      '1': 'deletionprotectionenabled',
      '3': 259418450,
      '4': 1,
      '5': 8,
      '10': 'deletionprotectionenabled'
    },
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndexDescription',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'globaltablesettingsreplicationmode',
      '3': 10446577,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.GlobalTableSettingsReplicationMode',
      '10': 'globaltablesettingsreplicationmode'
    },
    {
      '1': 'globaltableversion',
      '3': 68234287,
      '4': 1,
      '5': 9,
      '10': 'globaltableversion'
    },
    {
      '1': 'globaltablewitnesses',
      '3': 4521286,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalTableWitnessDescription',
      '10': 'globaltablewitnesses'
    },
    {'1': 'itemcount', '3': 26280022, '4': 1, '5': 3, '10': 'itemcount'},
    {
      '1': 'keyschema',
      '3': 293038056,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.KeySchemaElement',
      '10': 'keyschema'
    },
    {
      '1': 'lateststreamarn',
      '3': 207365682,
      '4': 1,
      '5': 9,
      '10': 'lateststreamarn'
    },
    {
      '1': 'lateststreamlabel',
      '3': 328475293,
      '4': 1,
      '5': 9,
      '10': 'lateststreamlabel'
    },
    {
      '1': 'localsecondaryindexes',
      '3': 362339959,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.LocalSecondaryIndexDescription',
      '10': 'localsecondaryindexes'
    },
    {
      '1': 'multiregionconsistency',
      '3': 446019131,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.MultiRegionConsistency',
      '10': 'multiregionconsistency'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughputDescription',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'replicas',
      '3': 306066781,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaDescription',
      '10': 'replicas'
    },
    {
      '1': 'restoresummary',
      '3': 330529648,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.RestoreSummary',
      '10': 'restoresummary'
    },
    {
      '1': 'ssedescription',
      '3': 350068773,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SSEDescription',
      '10': 'ssedescription'
    },
    {
      '1': 'streamspecification',
      '3': 403922627,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.StreamSpecification',
      '10': 'streamspecification'
    },
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {
      '1': 'tableclasssummary',
      '3': 4371552,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableClassSummary',
      '10': 'tableclasssummary'
    },
    {'1': 'tableid', '3': 449893011, '4': 1, '5': 9, '10': 'tableid'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'tablesizebytes',
      '3': 220631294,
      '4': 1,
      '5': 3,
      '10': 'tablesizebytes'
    },
    {
      '1': 'tablestatus',
      '3': 207908810,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableStatus',
      '10': 'tablestatus'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableWarmThroughputDescription',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `TableDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableDescriptionDescriptor = $convert.base64Decode(
    'ChBUYWJsZURlc2NyaXB0aW9uEkYKD2FyY2hpdmFsc3VtbWFyeRjqn+gYIAEoCzIZLmR5bmFtb2'
    'RiLkFyY2hpdmFsU3VtbWFyeVIPYXJjaGl2YWxzdW1tYXJ5ElUKFGF0dHJpYnV0ZWRlZmluaXRp'
    'b25zGIS/3sUBIAMoCzIdLmR5bmFtb2RiLkF0dHJpYnV0ZURlZmluaXRpb25SFGF0dHJpYnV0ZW'
    'RlZmluaXRpb25zEk8KEmJpbGxpbmdtb2Rlc3VtbWFyeRiaif1NIAEoCzIcLmR5bmFtb2RiLkJp'
    'bGxpbmdNb2RlU3VtbWFyeVISYmlsbGluZ21vZGVzdW1tYXJ5Ei0KEGNyZWF0aW9uZGF0ZXRpbW'
    'UY+vOoFyABKAlSEGNyZWF0aW9uZGF0ZXRpbWUSPwoZZGVsZXRpb25wcm90ZWN0aW9uZW5hYmxl'
    'ZBjS0tl7IAEoCFIZZGVsZXRpb25wcm90ZWN0aW9uZW5hYmxlZBJlChZnbG9iYWxzZWNvbmRhcn'
    'lpbmRleGVzGKn6jMMBIAMoCzIpLmR5bmFtb2RiLkdsb2JhbFNlY29uZGFyeUluZGV4RGVzY3Jp'
    'cHRpb25SFmdsb2JhbHNlY29uZGFyeWluZGV4ZXMSfwoiZ2xvYmFsdGFibGVzZXR0aW5nc3JlcG'
    'xpY2F0aW9ubW9kZRjxzf0EIAEoDjIsLmR5bmFtb2RiLkdsb2JhbFRhYmxlU2V0dGluZ3NSZXBs'
    'aWNhdGlvbk1vZGVSImdsb2JhbHRhYmxlc2V0dGluZ3NyZXBsaWNhdGlvbm1vZGUSMQoSZ2xvYm'
    'FsdGFibGV2ZXJzaW9uGK/YxCAgASgJUhJnbG9iYWx0YWJsZXZlcnNpb24SXgoUZ2xvYmFsdGFi'
    'bGV3aXRuZXNzZXMYxvqTAiADKAsyJy5keW5hbW9kYi5HbG9iYWxUYWJsZVdpdG5lc3NEZXNjcm'
    'lwdGlvblIUZ2xvYmFsdGFibGV3aXRuZXNzZXMSHwoJaXRlbWNvdW50GNaAxAwgASgDUglpdGVt'
    'Y291bnQSPAoJa2V5c2NoZW1hGOjP3YsBIAMoCzIaLmR5bmFtb2RiLktleVNjaGVtYUVsZW1lbn'
    'RSCWtleXNjaGVtYRIrCg9sYXRlc3RzdHJlYW1hcm4YsszwYiABKAlSD2xhdGVzdHN0cmVhbWFy'
    'bhIwChFsYXRlc3RzdHJlYW1sYWJlbBidxdCcASABKAlSEWxhdGVzdHN0cmVhbWxhYmVsEmIKFW'
    'xvY2Fsc2Vjb25kYXJ5aW5kZXhlcxj3vOOsASADKAsyKC5keW5hbW9kYi5Mb2NhbFNlY29uZGFy'
    'eUluZGV4RGVzY3JpcHRpb25SFWxvY2Fsc2Vjb25kYXJ5aW5kZXhlcxJcChZtdWx0aXJlZ2lvbm'
    'NvbnNpc3RlbmN5GLvs1tQBIAEoDjIgLmR5bmFtb2RiLk11bHRpUmVnaW9uQ29uc2lzdGVuY3lS'
    'Fm11bHRpcmVnaW9uY29uc2lzdGVuY3kSUAoSb25kZW1hbmR0aHJvdWdocHV0GILe2uUBIAEoCz'
    'IcLmR5bmFtb2RiLk9uRGVtYW5kVGhyb3VnaHB1dFISb25kZW1hbmR0aHJvdWdocHV0EmIKFXBy'
    'b3Zpc2lvbmVkdGhyb3VnaHB1dBiMo2sgASgLMiouZHluYW1vZGIuUHJvdmlzaW9uZWRUaHJvdW'
    'docHV0RGVzY3JpcHRpb25SFXByb3Zpc2lvbmVkdGhyb3VnaHB1dBI8CghyZXBsaWNhcxjd6viR'
    'ASADKAsyHC5keW5hbW9kYi5SZXBsaWNhRGVzY3JpcHRpb25SCHJlcGxpY2FzEkQKDnJlc3Rvcm'
    'VzdW1tYXJ5GPD2zZ0BIAEoCzIYLmR5bmFtb2RiLlJlc3RvcmVTdW1tYXJ5Ug5yZXN0b3Jlc3Vt'
    'bWFyeRJECg5zc2VkZXNjcmlwdGlvbhilwPamASABKAsyGC5keW5hbW9kYi5TU0VEZXNjcmlwdG'
    'lvblIOc3NlZGVzY3JpcHRpb24SUwoTc3RyZWFtc3BlY2lmaWNhdGlvbhjDvc3AASABKAsyHS5k'
    'eW5hbW9kYi5TdHJlYW1TcGVjaWZpY2F0aW9uUhNzdHJlYW1zcGVjaWZpY2F0aW9uEh4KCHRhYm'
    'xlYXJuGOOA680BIAEoCVIIdGFibGVhcm4STAoRdGFibGVjbGFzc3N1bW1hcnkY4OiKAiABKAsy'
    'Gy5keW5hbW9kYi5UYWJsZUNsYXNzU3VtbWFyeVIRdGFibGVjbGFzc3N1bW1hcnkSHAoHdGFibG'
    'VpZBiTpcPWASABKAlSB3RhYmxlaWQSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1l'
    'EikKDnRhYmxlc2l6ZWJ5dGVzGP6hmmkgASgDUg50YWJsZXNpemVieXRlcxI6Cgt0YWJsZXN0YX'
    'R1cxjK35FjIAEoDjIVLmR5bmFtb2RiLlRhYmxlU3RhdHVzUgt0YWJsZXN0YXR1cxJUCg53YXJt'
    'dGhyb3VnaHB1dBiD3siKASABKAsyKC5keW5hbW9kYi5UYWJsZVdhcm1UaHJvdWdocHV0RGVzY3'
    'JpcHRpb25SDndhcm10aHJvdWdocHV0');

@$core.Deprecated('Use tableInUseExceptionDescriptor instead')
const TableInUseException$json = {
  '1': 'TableInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TableInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChNUYWJsZUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tableNotFoundExceptionDescriptor instead')
const TableNotFoundException$json = {
  '1': 'TableNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TableNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChZUYWJsZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use tableWarmThroughputDescriptionDescriptor instead')
const TableWarmThroughputDescription$json = {
  '1': 'TableWarmThroughputDescription',
  '2': [
    {
      '1': 'readunitspersecond',
      '3': 11400732,
      '4': 1,
      '5': 3,
      '10': 'readunitspersecond'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableStatus',
      '10': 'status'
    },
    {
      '1': 'writeunitspersecond',
      '3': 339770127,
      '4': 1,
      '5': 3,
      '10': 'writeunitspersecond'
    },
  ],
};

/// Descriptor for `TableWarmThroughputDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableWarmThroughputDescriptionDescriptor =
    $convert.base64Decode(
        'Ch5UYWJsZVdhcm1UaHJvdWdocHV0RGVzY3JpcHRpb24SMQoScmVhZHVuaXRzcGVyc2Vjb25kGJ'
        'zstwUgASgDUhJyZWFkdW5pdHNwZXJzZWNvbmQSMAoGc3RhdHVzGJDk+wIgASgOMhUuZHluYW1v'
        'ZGIuVGFibGVTdGF0dXNSBnN0YXR1cxI0ChN3cml0ZXVuaXRzcGVyc2Vjb25kGI/2gaIBIAEoA1'
        'ITd3JpdGV1bml0c3BlcnNlY29uZA==');

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

@$core.Deprecated('Use tagResourceInputDescriptor instead')
const TagResourceInput$json = {
  '1': 'TagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceInputDescriptor = $convert.base64Decode(
    'ChBUYWdSZXNvdXJjZUlucHV0EiQKC3Jlc291cmNlYXJuGK342a0BIAEoCVILcmVzb3VyY2Vhcm'
    '4SJQoEdGFncxjBwfa1ASADKAsyDS5keW5hbW9kYi5UYWdSBHRhZ3M=');

@$core.Deprecated('Use throttlingExceptionDescriptor instead')
const ThrottlingException$json = {
  '1': 'ThrottlingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'throttlingreasons',
      '3': 170616084,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ThrottlingReason',
      '10': 'throttlingreasons'
    },
  ],
};

/// Descriptor for `ThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttlingExceptionDescriptor = $convert.base64Decode(
    'ChNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2USSwoRdG'
    'hyb3R0bGluZ3JlYXNvbnMYlMqtUSADKAsyGi5keW5hbW9kYi5UaHJvdHRsaW5nUmVhc29uUhF0'
    'aHJvdHRsaW5ncmVhc29ucw==');

@$core.Deprecated('Use throttlingReasonDescriptor instead')
const ThrottlingReason$json = {
  '1': 'ThrottlingReason',
  '2': [
    {'1': 'reason', '3': 413359642, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
  ],
};

/// Descriptor for `ThrottlingReason`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttlingReasonDescriptor = $convert.base64Decode(
    'ChBUaHJvdHRsaW5nUmVhc29uEhoKBnJlYXNvbhiavI3FASABKAlSBnJlYXNvbhIdCghyZXNvdX'
    'JjZRj2//1OIAEoCVIIcmVzb3VyY2U=');

@$core.Deprecated('Use timeToLiveDescriptionDescriptor instead')
const TimeToLiveDescription$json = {
  '1': 'TimeToLiveDescription',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'timetolivestatus',
      '3': 279467554,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TimeToLiveStatus',
      '10': 'timetolivestatus'
    },
  ],
};

/// Descriptor for `TimeToLiveDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timeToLiveDescriptionDescriptor = $convert.base64Decode(
    'ChVUaW1lVG9MaXZlRGVzY3JpcHRpb24SKAoNYXR0cmlidXRlbmFtZRitlZioASABKAlSDWF0dH'
    'JpYnV0ZW5hbWUSSgoQdGltZXRvbGl2ZXN0YXR1cxiirKGFASABKA4yGi5keW5hbW9kYi5UaW1l'
    'VG9MaXZlU3RhdHVzUhB0aW1ldG9saXZlc3RhdHVz');

@$core.Deprecated('Use timeToLiveSpecificationDescriptor instead')
const TimeToLiveSpecification$json = {
  '1': 'TimeToLiveSpecification',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
  ],
};

/// Descriptor for `TimeToLiveSpecification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timeToLiveSpecificationDescriptor =
    $convert.base64Decode(
        'ChdUaW1lVG9MaXZlU3BlY2lmaWNhdGlvbhIoCg1hdHRyaWJ1dGVuYW1lGK2VmKgBIAEoCVINYX'
        'R0cmlidXRlbmFtZRIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZA==');

@$core.Deprecated('Use transactGetItemDescriptor instead')
const TransactGetItem$json = {
  '1': 'TransactGetItem',
  '2': [
    {
      '1': 'get',
      '3': 379010808,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Get',
      '10': 'get'
    },
  ],
};

/// Descriptor for `TransactGetItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactGetItemDescriptor = $convert.base64Decode(
    'Cg9UcmFuc2FjdEdldEl0ZW0SIwoDZ2V0GPj93LQBIAEoCzINLmR5bmFtb2RiLkdldFIDZ2V0');

@$core.Deprecated('Use transactGetItemsInputDescriptor instead')
const TransactGetItemsInput$json = {
  '1': 'TransactGetItemsInput',
  '2': [
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'transactitems',
      '3': 506245290,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.TransactGetItem',
      '10': 'transactitems'
    },
  ],
};

/// Descriptor for `TransactGetItemsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactGetItemsInputDescriptor = $convert.base64Decode(
    'ChVUcmFuc2FjdEdldEl0ZW1zSW5wdXQSWwoWcmV0dXJuY29uc3VtZWRjYXBhY2l0eRj+5+EUIA'
    'EoDjIgLmR5bmFtb2RiLlJldHVybkNvbnN1bWVkQ2FwYWNpdHlSFnJldHVybmNvbnN1bWVkY2Fw'
    'YWNpdHkSQwoNdHJhbnNhY3RpdGVtcxiq4bLxASADKAsyGS5keW5hbW9kYi5UcmFuc2FjdEdldE'
    'l0ZW1SDXRyYW5zYWN0aXRlbXM=');

@$core.Deprecated('Use transactGetItemsOutputDescriptor instead')
const TransactGetItemsOutput$json = {
  '1': 'TransactGetItemsOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'responses',
      '3': 114852856,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ItemResponse',
      '10': 'responses'
    },
  ],
};

/// Descriptor for `TransactGetItemsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactGetItemsOutputDescriptor = $convert.base64Decode(
    'ChZUcmFuc2FjdEdldEl0ZW1zT3V0cHV0EkoKEGNvbnN1bWVkY2FwYWNpdHkYrKqh1gEgAygLMh'
    'ouZHluYW1vZGIuQ29uc3VtZWRDYXBhY2l0eVIQY29uc3VtZWRjYXBhY2l0eRI3CglyZXNwb25z'
    'ZXMY+IfiNiADKAsyFi5keW5hbW9kYi5JdGVtUmVzcG9uc2VSCXJlc3BvbnNlcw==');

@$core.Deprecated('Use transactWriteItemDescriptor instead')
const TransactWriteItem$json = {
  '1': 'TransactWriteItem',
  '2': [
    {
      '1': 'conditioncheck',
      '3': 167629711,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConditionCheck',
      '10': 'conditioncheck'
    },
    {
      '1': 'delete',
      '3': 395831915,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Delete',
      '10': 'delete'
    },
    {
      '1': 'put',
      '3': 242822543,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Put',
      '10': 'put'
    },
    {
      '1': 'update',
      '3': 237178517,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.Update',
      '10': 'update'
    },
  ],
};

/// Descriptor for `TransactWriteItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactWriteItemDescriptor = $convert.base64Decode(
    'ChFUcmFuc2FjdFdyaXRlSXRlbRJDCg5jb25kaXRpb25jaGVjaxiPp/dPIAEoCzIYLmR5bmFtb2'
    'RiLkNvbmRpdGlvbkNoZWNrUg5jb25kaXRpb25jaGVjaxIsCgZkZWxldGUY69TfvAEgASgLMhAu'
    'ZHluYW1vZGIuRGVsZXRlUgZkZWxldGUSIgoDcHV0GI/b5HMgASgLMg0uZHluYW1vZGIuUHV0Ug'
    'NwdXQSKwoGdXBkYXRlGJWdjHEgASgLMhAuZHluYW1vZGIuVXBkYXRlUgZ1cGRhdGU=');

@$core.Deprecated('Use transactWriteItemsInputDescriptor instead')
const TransactWriteItemsInput$json = {
  '1': 'TransactWriteItemsInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'returnitemcollectionmetrics',
      '3': 255507354,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnItemCollectionMetrics',
      '10': 'returnitemcollectionmetrics'
    },
    {
      '1': 'transactitems',
      '3': 506245290,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.TransactWriteItem',
      '10': 'transactitems'
    },
  ],
};

/// Descriptor for `TransactWriteItemsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactWriteItemsInputDescriptor = $convert.base64Decode(
    'ChdUcmFuc2FjdFdyaXRlSXRlbXNJbnB1dBIyChJjbGllbnRyZXF1ZXN0dG9rZW4Y8e+i2QEgAS'
    'gJUhJjbGllbnRyZXF1ZXN0dG9rZW4SWwoWcmV0dXJuY29uc3VtZWRjYXBhY2l0eRj+5+EUIAEo'
    'DjIgLmR5bmFtb2RiLlJldHVybkNvbnN1bWVkQ2FwYWNpdHlSFnJldHVybmNvbnN1bWVkY2FwYW'
    'NpdHkSagobcmV0dXJuaXRlbWNvbGxlY3Rpb25tZXRyaWNzGJr36nkgASgOMiUuZHluYW1vZGIu'
    'UmV0dXJuSXRlbUNvbGxlY3Rpb25NZXRyaWNzUhtyZXR1cm5pdGVtY29sbGVjdGlvbm1ldHJpY3'
    'MSRQoNdHJhbnNhY3RpdGVtcxiq4bLxASADKAsyGy5keW5hbW9kYi5UcmFuc2FjdFdyaXRlSXRl'
    'bVINdHJhbnNhY3RpdGVtcw==');

@$core.Deprecated('Use transactWriteItemsOutputDescriptor instead')
const TransactWriteItemsOutput$json = {
  '1': 'TransactWriteItemsOutput',
  '2': [
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'itemcollectionmetrics',
      '3': 185317452,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.TransactWriteItemsOutput.ItemcollectionmetricsEntry',
      '10': 'itemcollectionmetrics'
    },
  ],
  '3': [TransactWriteItemsOutput_ItemcollectionmetricsEntry$json],
};

@$core.Deprecated('Use transactWriteItemsOutputDescriptor instead')
const TransactWriteItemsOutput_ItemcollectionmetricsEntry$json = {
  '1': 'ItemcollectionmetricsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TransactWriteItemsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactWriteItemsOutputDescriptor = $convert.base64Decode(
    'ChhUcmFuc2FjdFdyaXRlSXRlbXNPdXRwdXQSSgoQY29uc3VtZWRjYXBhY2l0eRisqqHWASADKA'
    'syGi5keW5hbW9kYi5Db25zdW1lZENhcGFjaXR5UhBjb25zdW1lZGNhcGFjaXR5EnYKFWl0ZW1j'
    'b2xsZWN0aW9ubWV0cmljcxjM8K5YIAMoCzI9LmR5bmFtb2RiLlRyYW5zYWN0V3JpdGVJdGVtc0'
    '91dHB1dC5JdGVtY29sbGVjdGlvbm1ldHJpY3NFbnRyeVIVaXRlbWNvbGxlY3Rpb25tZXRyaWNz'
    'GkgKGkl0ZW1jb2xsZWN0aW9ubWV0cmljc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbH'
    'VlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use transactionCanceledExceptionDescriptor instead')
const TransactionCanceledException$json = {
  '1': 'TransactionCanceledException',
  '2': [
    {
      '1': 'cancellationreasons',
      '3': 346602284,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.CancellationReason',
      '10': 'cancellationreasons'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TransactionCanceledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionCanceledExceptionDescriptor =
    $convert.base64Decode(
        'ChxUcmFuc2FjdGlvbkNhbmNlbGVkRXhjZXB0aW9uElIKE2NhbmNlbGxhdGlvbnJlYXNvbnMYrP'
        'aipQEgAygLMhwuZHluYW1vZGIuQ2FuY2VsbGF0aW9uUmVhc29uUhNjYW5jZWxsYXRpb25yZWFz'
        'b25zEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use transactionConflictExceptionDescriptor instead')
const TransactionConflictException$json = {
  '1': 'TransactionConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TransactionConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionConflictExceptionDescriptor =
    $convert.base64Decode(
        'ChxUcmFuc2FjdGlvbkNvbmZsaWN0RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use transactionInProgressExceptionDescriptor instead')
const TransactionInProgressException$json = {
  '1': 'TransactionInProgressException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TransactionInProgressException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionInProgressExceptionDescriptor =
    $convert.base64Decode(
        'Ch5UcmFuc2FjdGlvbkluUHJvZ3Jlc3NFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use untagResourceInputDescriptor instead')
const UntagResourceInput$json = {
  '1': 'UntagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceInputDescriptor = $convert.base64Decode(
    'ChJVbnRhZ1Jlc291cmNlSW5wdXQSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJjZW'
    'FybhIcCgd0YWdrZXlzGPzD85gBIAMoCVIHdGFna2V5cw==');

@$core.Deprecated('Use updateDescriptor instead')
const Update$json = {
  '1': 'Update',
  '2': [
    {
      '1': 'conditionexpression',
      '3': 409657405,
      '4': 1,
      '5': 9,
      '10': 'conditionexpression'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Update.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Update.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.Update.KeyEntry',
      '10': 'key'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'updateexpression',
      '3': 389263879,
      '4': 1,
      '5': 9,
      '10': 'updateexpression'
    },
  ],
  '3': [
    Update_ExpressionattributenamesEntry$json,
    Update_ExpressionattributevaluesEntry$json,
    Update_KeyEntry$json
  ],
};

@$core.Deprecated('Use updateDescriptor instead')
const Update_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateDescriptor instead')
const Update_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateDescriptor instead')
const Update_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `Update`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDescriptor = $convert.base64Decode(
    'CgZVcGRhdGUSNAoTY29uZGl0aW9uZXhwcmVzc2lvbhi9wKvDASABKAlSE2NvbmRpdGlvbmV4cH'
    'Jlc3Npb24SbQoYZXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMi4uZHluYW1vZGIu'
    'VXBkYXRlLkV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9uYXR0cmlidX'
    'RlbmFtZXMScQoZZXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlcxjYnKDnASADKAsyLy5keW5hbW9k'
    'Yi5VcGRhdGUuRXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5UhlleHByZXNzaW9uYXR0cm'
    'lidXRldmFsdWVzEi4KA2tleRiNkutoIAMoCzIZLmR5bmFtb2RiLlVwZGF0ZS5LZXlFbnRyeVID'
    'a2V5EoIBCiNyZXR1cm52YWx1ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZRjgl4ECIAEoDjItLm'
    'R5bmFtb2RiLlJldHVyblZhbHVlc09uQ29uZGl0aW9uQ2hlY2tGYWlsdXJlUiNyZXR1cm52YWx1'
    'ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZRIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW'
    '5hbWUSLgoQdXBkYXRlZXhwcmVzc2lvbhiH5M65ASABKAlSEHVwZGF0ZWV4cHJlc3Npb24aSwod'
    'RXhwcmVzc2lvbmF0dHJpYnV0ZW5hbWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdW'
    'UYAiABKAlSBXZhbHVlOgI4ARpmCh5FeHByZXNzaW9uYXR0cmlidXRldmFsdWVzRW50cnkSEAoD'
    'a2V5GAEgASgJUgNrZXkSLgoFdmFsdWUYAiABKAsyGC5keW5hbW9kYi5BdHRyaWJ1dGVWYWx1ZV'
    'IFdmFsdWU6AjgBGlAKCEtleUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgL'
    'MhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use updateContinuousBackupsInputDescriptor instead')
const UpdateContinuousBackupsInput$json = {
  '1': 'UpdateContinuousBackupsInput',
  '2': [
    {
      '1': 'pointintimerecoveryspecification',
      '3': 395539950,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.PointInTimeRecoverySpecification',
      '10': 'pointintimerecoveryspecification'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `UpdateContinuousBackupsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContinuousBackupsInputDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVDb250aW51b3VzQmFja3Vwc0lucHV0EnoKIHBvaW50aW50aW1lcmVjb3ZlcnlzcG'
    'VjaWZpY2F0aW9uGO7rzbwBIAEoCzIqLmR5bmFtb2RiLlBvaW50SW5UaW1lUmVjb3ZlcnlTcGVj'
    'aWZpY2F0aW9uUiBwb2ludGludGltZXJlY292ZXJ5c3BlY2lmaWNhdGlvbhIgCgl0YWJsZW5hbW'
    'UY3eTagQEgASgJUgl0YWJsZW5hbWU=');

@$core.Deprecated('Use updateContinuousBackupsOutputDescriptor instead')
const UpdateContinuousBackupsOutput$json = {
  '1': 'UpdateContinuousBackupsOutput',
  '2': [
    {
      '1': 'continuousbackupsdescription',
      '3': 446030650,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ContinuousBackupsDescription',
      '10': 'continuousbackupsdescription'
    },
  ],
};

/// Descriptor for `UpdateContinuousBackupsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContinuousBackupsOutputDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVDb250aW51b3VzQmFja3Vwc091dHB1dBJuChxjb250aW51b3VzYmFja3Vwc2Rlc2'
        'NyaXB0aW9uGLrG19QBIAEoCzImLmR5bmFtb2RiLkNvbnRpbnVvdXNCYWNrdXBzRGVzY3JpcHRp'
        'b25SHGNvbnRpbnVvdXNiYWNrdXBzZGVzY3JpcHRpb24=');

@$core.Deprecated('Use updateContributorInsightsInputDescriptor instead')
const UpdateContributorInsightsInput$json = {
  '1': 'UpdateContributorInsightsInput',
  '2': [
    {
      '1': 'contributorinsightsaction',
      '3': 81909182,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsAction',
      '10': 'contributorinsightsaction'
    },
    {
      '1': 'contributorinsightsmode',
      '3': 86700161,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsMode',
      '10': 'contributorinsightsmode'
    },
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `UpdateContributorInsightsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContributorInsightsInputDescriptor = $convert.base64Decode(
    'Ch5VcGRhdGVDb250cmlidXRvckluc2lnaHRzSW5wdXQSZAoZY29udHJpYnV0b3JpbnNpZ2h0c2'
    'FjdGlvbhi+q4cnIAEoDjIjLmR5bmFtb2RiLkNvbnRyaWJ1dG9ySW5zaWdodHNBY3Rpb25SGWNv'
    'bnRyaWJ1dG9yaW5zaWdodHNhY3Rpb24SXgoXY29udHJpYnV0b3JpbnNpZ2h0c21vZGUYgeGrKS'
    'ABKA4yIS5keW5hbW9kYi5Db250cmlidXRvckluc2lnaHRzTW9kZVIXY29udHJpYnV0b3JpbnNp'
    'Z2h0c21vZGUSHwoJaW5kZXhuYW1lGJHV6zAgASgJUglpbmRleG5hbWUSIAoJdGFibGVuYW1lGN'
    '3k2oEBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use updateContributorInsightsOutputDescriptor instead')
const UpdateContributorInsightsOutput$json = {
  '1': 'UpdateContributorInsightsOutput',
  '2': [
    {
      '1': 'contributorinsightsmode',
      '3': 86700161,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsMode',
      '10': 'contributorinsightsmode'
    },
    {
      '1': 'contributorinsightsstatus',
      '3': 363347282,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ContributorInsightsStatus',
      '10': 'contributorinsightsstatus'
    },
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `UpdateContributorInsightsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContributorInsightsOutputDescriptor = $convert.base64Decode(
    'Ch9VcGRhdGVDb250cmlidXRvckluc2lnaHRzT3V0cHV0El4KF2NvbnRyaWJ1dG9yaW5zaWdodH'
    'Ntb2RlGIHhqykgASgOMiEuZHluYW1vZGIuQ29udHJpYnV0b3JJbnNpZ2h0c01vZGVSF2NvbnRy'
    'aWJ1dG9yaW5zaWdodHNtb2RlEmUKGWNvbnRyaWJ1dG9yaW5zaWdodHNzdGF0dXMY0vqgrQEgAS'
    'gOMiMuZHluYW1vZGIuQ29udHJpYnV0b3JJbnNpZ2h0c1N0YXR1c1IZY29udHJpYnV0b3JpbnNp'
    'Z2h0c3N0YXR1cxIfCglpbmRleG5hbWUYkdXrMCABKAlSCWluZGV4bmFtZRIgCgl0YWJsZW5hbW'
    'UY3eTagQEgASgJUgl0YWJsZW5hbWU=');

@$core.Deprecated('Use updateGlobalSecondaryIndexActionDescriptor instead')
const UpdateGlobalSecondaryIndexAction$json = {
  '1': 'UpdateGlobalSecondaryIndexAction',
  '2': [
    {'1': 'indexname', '3': 102427281, '4': 1, '5': 9, '10': 'indexname'},
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.WarmThroughput',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `UpdateGlobalSecondaryIndexAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGlobalSecondaryIndexActionDescriptor = $convert.base64Decode(
    'CiBVcGRhdGVHbG9iYWxTZWNvbmRhcnlJbmRleEFjdGlvbhIfCglpbmRleG5hbWUYkdXrMCABKA'
    'lSCWluZGV4bmFtZRJQChJvbmRlbWFuZHRocm91Z2hwdXQYgt7a5QEgASgLMhwuZHluYW1vZGIu'
    'T25EZW1hbmRUaHJvdWdocHV0UhJvbmRlbWFuZHRocm91Z2hwdXQSVwoVcHJvdmlzaW9uZWR0aH'
    'JvdWdocHV0GIyjayABKAsyHy5keW5hbW9kYi5Qcm92aXNpb25lZFRocm91Z2hwdXRSFXByb3Zp'
    'c2lvbmVkdGhyb3VnaHB1dBJECg53YXJtdGhyb3VnaHB1dBiD3siKASABKAsyGC5keW5hbW9kYi'
    '5XYXJtVGhyb3VnaHB1dFIOd2FybXRocm91Z2hwdXQ=');

@$core.Deprecated('Use updateGlobalTableInputDescriptor instead')
const UpdateGlobalTableInput$json = {
  '1': 'UpdateGlobalTableInput',
  '2': [
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
    {
      '1': 'replicaupdates',
      '3': 260731936,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaUpdate',
      '10': 'replicaupdates'
    },
  ],
};

/// Descriptor for `UpdateGlobalTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGlobalTableInputDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVHbG9iYWxUYWJsZUlucHV0EiwKD2dsb2JhbHRhYmxlbmFtZRiqpqeHASABKAlSD2'
    'dsb2JhbHRhYmxlbmFtZRJCCg5yZXBsaWNhdXBkYXRlcxig6Kl8IAMoCzIXLmR5bmFtb2RiLlJl'
    'cGxpY2FVcGRhdGVSDnJlcGxpY2F1cGRhdGVz');

@$core.Deprecated('Use updateGlobalTableOutputDescriptor instead')
const UpdateGlobalTableOutput$json = {
  '1': 'UpdateGlobalTableOutput',
  '2': [
    {
      '1': 'globaltabledescription',
      '3': 342116405,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.GlobalTableDescription',
      '10': 'globaltabledescription'
    },
  ],
};

/// Descriptor for `UpdateGlobalTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGlobalTableOutputDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVHbG9iYWxUYWJsZU91dHB1dBJcChZnbG9iYWx0YWJsZWRlc2NyaXB0aW9uGLWQka'
    'MBIAEoCzIgLmR5bmFtb2RiLkdsb2JhbFRhYmxlRGVzY3JpcHRpb25SFmdsb2JhbHRhYmxlZGVz'
    'Y3JpcHRpb24=');

@$core.Deprecated('Use updateGlobalTableSettingsInputDescriptor instead')
const UpdateGlobalTableSettingsInput$json = {
  '1': 'UpdateGlobalTableSettingsInput',
  '2': [
    {
      '1': 'globaltablebillingmode',
      '3': 285868859,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'globaltablebillingmode'
    },
    {
      '1': 'globaltableglobalsecondaryindexsettingsupdate',
      '3': 368411116,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalTableGlobalSecondaryIndexSettingsUpdate',
      '10': 'globaltableglobalsecondaryindexsettingsupdate'
    },
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
    {
      '1': 'globaltableprovisionedwritecapacityautoscalingsettingsupdate',
      '3': 435219182,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'globaltableprovisionedwritecapacityautoscalingsettingsupdate'
    },
    {
      '1': 'globaltableprovisionedwritecapacityunits',
      '3': 127172283,
      '4': 1,
      '5': 3,
      '10': 'globaltableprovisionedwritecapacityunits'
    },
    {
      '1': 'replicasettingsupdate',
      '3': 463498612,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaSettingsUpdate',
      '10': 'replicasettingsupdate'
    },
  ],
};

/// Descriptor for `UpdateGlobalTableSettingsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGlobalTableSettingsInputDescriptor = $convert.base64Decode(
    'Ch5VcGRhdGVHbG9iYWxUYWJsZVNldHRpbmdzSW5wdXQSUQoWZ2xvYmFsdGFibGViaWxsaW5nbW'
    '9kZRi7hqiIASABKA4yFS5keW5hbW9kYi5CaWxsaW5nTW9kZVIWZ2xvYmFsdGFibGViaWxsaW5n'
    'bW9kZRKhAQotZ2xvYmFsdGFibGVnbG9iYWxzZWNvbmRhcnlpbmRleHNldHRpbmdzdXBkYXRlGO'
    'yD1q8BIAMoCzI3LmR5bmFtb2RiLkdsb2JhbFRhYmxlR2xvYmFsU2Vjb25kYXJ5SW5kZXhTZXR0'
    'aW5nc1VwZGF0ZVItZ2xvYmFsdGFibGVnbG9iYWxzZWNvbmRhcnlpbmRleHNldHRpbmdzdXBkYX'
    'RlEiwKD2dsb2JhbHRhYmxlbmFtZRiqpqeHASABKAlSD2dsb2JhbHRhYmxlbmFtZRKrAQo8Z2xv'
    'YmFsdGFibGVwcm92aXNpb25lZHdyaXRlY2FwYWNpdHlhdXRvc2NhbGluZ3NldHRpbmdzdXBkYX'
    'RlGO7Vw88BIAEoCzIjLmR5bmFtb2RiLkF1dG9TY2FsaW5nU2V0dGluZ3NVcGRhdGVSPGdsb2Jh'
    'bHRhYmxlcHJvdmlzaW9uZWR3cml0ZWNhcGFjaXR5YXV0b3NjYWxpbmdzZXR0aW5nc3VwZGF0ZR'
    'JdCihnbG9iYWx0YWJsZXByb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eXVuaXRzGLv90TwgASgDUihn'
    'bG9iYWx0YWJsZXByb3Zpc2lvbmVkd3JpdGVjYXBhY2l0eXVuaXRzElkKFXJlcGxpY2FzZXR0aW'
    '5nc3VwZGF0ZRj02oHdASADKAsyHy5keW5hbW9kYi5SZXBsaWNhU2V0dGluZ3NVcGRhdGVSFXJl'
    'cGxpY2FzZXR0aW5nc3VwZGF0ZQ==');

@$core.Deprecated('Use updateGlobalTableSettingsOutputDescriptor instead')
const UpdateGlobalTableSettingsOutput$json = {
  '1': 'UpdateGlobalTableSettingsOutput',
  '2': [
    {
      '1': 'globaltablename',
      '3': 283759402,
      '4': 1,
      '5': 9,
      '10': 'globaltablename'
    },
    {
      '1': 'replicasettings',
      '3': 288882107,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaSettingsDescription',
      '10': 'replicasettings'
    },
  ],
};

/// Descriptor for `UpdateGlobalTableSettingsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateGlobalTableSettingsOutputDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVHbG9iYWxUYWJsZVNldHRpbmdzT3V0cHV0EiwKD2dsb2JhbHRhYmxlbmFtZRiqpq'
        'eHASABKAlSD2dsb2JhbHRhYmxlbmFtZRJSCg9yZXBsaWNhc2V0dGluZ3MYu/vfiQEgAygLMiQu'
        'ZHluYW1vZGIuUmVwbGljYVNldHRpbmdzRGVzY3JpcHRpb25SD3JlcGxpY2FzZXR0aW5ncw==');

@$core.Deprecated('Use updateItemInputDescriptor instead')
const UpdateItemInput$json = {
  '1': 'UpdateItemInput',
  '2': [
    {
      '1': 'attributeupdates',
      '3': 32403512,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.UpdateItemInput.AttributeupdatesEntry',
      '10': 'attributeupdates'
    },
    {
      '1': 'conditionexpression',
      '3': 409657405,
      '4': 1,
      '5': 9,
      '10': 'conditionexpression'
    },
    {
      '1': 'conditionaloperator',
      '3': 172066260,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ConditionalOperator',
      '10': 'conditionaloperator'
    },
    {
      '1': 'expected',
      '3': 106557946,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.UpdateItemInput.ExpectedEntry',
      '10': 'expected'
    },
    {
      '1': 'expressionattributenames',
      '3': 150228092,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.UpdateItemInput.ExpressionattributenamesEntry',
      '10': 'expressionattributenames'
    },
    {
      '1': 'expressionattributevalues',
      '3': 484970072,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.UpdateItemInput.ExpressionattributevaluesEntry',
      '10': 'expressionattributevalues'
    },
    {
      '1': 'key',
      '3': 219859213,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.UpdateItemInput.KeyEntry',
      '10': 'key'
    },
    {
      '1': 'returnconsumedcapacity',
      '3': 43545598,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnConsumedCapacity',
      '10': 'returnconsumedcapacity'
    },
    {
      '1': 'returnitemcollectionmetrics',
      '3': 255507354,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnItemCollectionMetrics',
      '10': 'returnitemcollectionmetrics'
    },
    {
      '1': 'returnvalues',
      '3': 402960198,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValue',
      '10': 'returnvalues'
    },
    {
      '1': 'returnvaluesonconditioncheckfailure',
      '3': 4213728,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ReturnValuesOnConditionCheckFailure',
      '10': 'returnvaluesonconditioncheckfailure'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'updateexpression',
      '3': 389263879,
      '4': 1,
      '5': 9,
      '10': 'updateexpression'
    },
  ],
  '3': [
    UpdateItemInput_AttributeupdatesEntry$json,
    UpdateItemInput_ExpectedEntry$json,
    UpdateItemInput_ExpressionattributenamesEntry$json,
    UpdateItemInput_ExpressionattributevaluesEntry$json,
    UpdateItemInput_KeyEntry$json
  ],
};

@$core.Deprecated('Use updateItemInputDescriptor instead')
const UpdateItemInput_AttributeupdatesEntry$json = {
  '1': 'AttributeupdatesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValueUpdate',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateItemInputDescriptor instead')
const UpdateItemInput_ExpectedEntry$json = {
  '1': 'ExpectedEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ExpectedAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateItemInputDescriptor instead')
const UpdateItemInput_ExpressionattributenamesEntry$json = {
  '1': 'ExpressionattributenamesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateItemInputDescriptor instead')
const UpdateItemInput_ExpressionattributevaluesEntry$json = {
  '1': 'ExpressionattributevaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateItemInputDescriptor instead')
const UpdateItemInput_KeyEntry$json = {
  '1': 'KeyEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateItemInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateItemInputDescriptor = $convert.base64Decode(
    'Cg9VcGRhdGVJdGVtSW5wdXQSXgoQYXR0cmlidXRldXBkYXRlcxi44LkPIAMoCzIvLmR5bmFtb2'
    'RiLlVwZGF0ZUl0ZW1JbnB1dC5BdHRyaWJ1dGV1cGRhdGVzRW50cnlSEGF0dHJpYnV0ZXVwZGF0'
    'ZXMSNAoTY29uZGl0aW9uZXhwcmVzc2lvbhi9wKvDASABKAlSE2NvbmRpdGlvbmV4cHJlc3Npb2'
    '4SUgoTY29uZGl0aW9uYWxvcGVyYXRvchjUi4ZSIAEoDjIdLmR5bmFtb2RiLkNvbmRpdGlvbmFs'
    'T3BlcmF0b3JSE2NvbmRpdGlvbmFsb3BlcmF0b3ISRgoIZXhwZWN0ZWQY+uPnMiADKAsyJy5keW'
    '5hbW9kYi5VcGRhdGVJdGVtSW5wdXQuRXhwZWN0ZWRFbnRyeVIIZXhwZWN0ZWQSdgoYZXhwcmVz'
    'c2lvbmF0dHJpYnV0ZW5hbWVzGPyY0UcgAygLMjcuZHluYW1vZGIuVXBkYXRlSXRlbUlucHV0Lk'
    'V4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5UhhleHByZXNzaW9uYXR0cmlidXRlbmFtZXMS'
    'egoZZXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlcxjYnKDnASADKAsyOC5keW5hbW9kYi5VcGRhdG'
    'VJdGVtSW5wdXQuRXhwcmVzc2lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5UhlleHByZXNzaW9uYXR0'
    'cmlidXRldmFsdWVzEjcKA2tleRiNkutoIAMoCzIiLmR5bmFtb2RiLlVwZGF0ZUl0ZW1JbnB1dC'
    '5LZXlFbnRyeVIDa2V5ElsKFnJldHVybmNvbnN1bWVkY2FwYWNpdHkY/ufhFCABKA4yIC5keW5h'
    'bW9kYi5SZXR1cm5Db25zdW1lZENhcGFjaXR5UhZyZXR1cm5jb25zdW1lZGNhcGFjaXR5EmoKG3'
    'JldHVybml0ZW1jb2xsZWN0aW9ubWV0cmljcxia9+p5IAEoDjIlLmR5bmFtb2RiLlJldHVybkl0'
    'ZW1Db2xsZWN0aW9uTWV0cmljc1IbcmV0dXJuaXRlbWNvbGxlY3Rpb25tZXRyaWNzEj0KDHJldH'
    'VybnZhbHVlcxjG3pLAASABKA4yFS5keW5hbW9kYi5SZXR1cm5WYWx1ZVIMcmV0dXJudmFsdWVz'
    'EoIBCiNyZXR1cm52YWx1ZXNvbmNvbmRpdGlvbmNoZWNrZmFpbHVyZRjgl4ECIAEoDjItLmR5bm'
    'Ftb2RiLlJldHVyblZhbHVlc09uQ29uZGl0aW9uQ2hlY2tGYWlsdXJlUiNyZXR1cm52YWx1ZXNv'
    'bmNvbmRpdGlvbmNoZWNrZmFpbHVyZRIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW5hbW'
    'USLgoQdXBkYXRlZXhwcmVzc2lvbhiH5M65ASABKAlSEHVwZGF0ZWV4cHJlc3Npb24aYwoVQXR0'
    'cmlidXRldXBkYXRlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EjQKBXZhbHVlGAIgASgLMh4uZH'
    'luYW1vZGIuQXR0cmlidXRlVmFsdWVVcGRhdGVSBXZhbHVlOgI4ARpdCg1FeHBlY3RlZEVudHJ5'
    'EhAKA2tleRgBIAEoCVIDa2V5EjYKBXZhbHVlGAIgASgLMiAuZHluYW1vZGIuRXhwZWN0ZWRBdH'
    'RyaWJ1dGVWYWx1ZVIFdmFsdWU6AjgBGksKHUV4cHJlc3Npb25hdHRyaWJ1dGVuYW1lc0VudHJ5'
    'EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaZgoeRXhwcmVzc2'
    'lvbmF0dHJpYnV0ZXZhbHVlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgL'
    'MhguZHluYW1vZGIuQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4ARpQCghLZXlFbnRyeRIQCgNrZX'
    'kYASABKAlSA2tleRIuCgV2YWx1ZRgCIAEoCzIYLmR5bmFtb2RiLkF0dHJpYnV0ZVZhbHVlUgV2'
    'YWx1ZToCOAE=');

@$core.Deprecated('Use updateItemOutputDescriptor instead')
const UpdateItemOutput$json = {
  '1': 'UpdateItemOutput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.UpdateItemOutput.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'consumedcapacity',
      '3': 449336620,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ConsumedCapacity',
      '10': 'consumedcapacity'
    },
    {
      '1': 'itemcollectionmetrics',
      '3': 185317452,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ItemCollectionMetrics',
      '10': 'itemcollectionmetrics'
    },
  ],
  '3': [UpdateItemOutput_AttributesEntry$json],
};

@$core.Deprecated('Use updateItemOutputDescriptor instead')
const UpdateItemOutput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateItemOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateItemOutputDescriptor = $convert.base64Decode(
    'ChBVcGRhdGVJdGVtT3V0cHV0Ek0KCmF0dHJpYnV0ZXMYtan7YyADKAsyKi5keW5hbW9kYi5VcG'
    'RhdGVJdGVtT3V0cHV0LkF0dHJpYnV0ZXNFbnRyeVIKYXR0cmlidXRlcxJKChBjb25zdW1lZGNh'
    'cGFjaXR5GKyqodYBIAEoCzIaLmR5bmFtb2RiLkNvbnN1bWVkQ2FwYWNpdHlSEGNvbnN1bWVkY2'
    'FwYWNpdHkSWAoVaXRlbWNvbGxlY3Rpb25tZXRyaWNzGMzwrlggASgLMh8uZHluYW1vZGIuSXRl'
    'bUNvbGxlY3Rpb25NZXRyaWNzUhVpdGVtY29sbGVjdGlvbm1ldHJpY3MaVwoPQXR0cmlidXRlc0'
    'VudHJ5EhAKA2tleRgBIAEoCVIDa2V5Ei4KBXZhbHVlGAIgASgLMhguZHluYW1vZGIuQXR0cmli'
    'dXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use updateKinesisStreamingConfigurationDescriptor instead')
const UpdateKinesisStreamingConfiguration$json = {
  '1': 'UpdateKinesisStreamingConfiguration',
  '2': [
    {
      '1': 'approximatecreationdatetimeprecision',
      '3': 392293352,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.ApproximateCreationDateTimePrecision',
      '10': 'approximatecreationdatetimeprecision'
    },
  ],
};

/// Descriptor for `UpdateKinesisStreamingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKinesisStreamingConfigurationDescriptor =
    $convert.base64Decode(
        'CiNVcGRhdGVLaW5lc2lzU3RyZWFtaW5nQ29uZmlndXJhdGlvbhKGAQokYXBwcm94aW1hdGVjcm'
        'VhdGlvbmRhdGV0aW1lcHJlY2lzaW9uGOjXh7sBIAEoDjIuLmR5bmFtb2RiLkFwcHJveGltYXRl'
        'Q3JlYXRpb25EYXRlVGltZVByZWNpc2lvblIkYXBwcm94aW1hdGVjcmVhdGlvbmRhdGV0aW1lcH'
        'JlY2lzaW9u');

@$core
    .Deprecated('Use updateKinesisStreamingDestinationInputDescriptor instead')
const UpdateKinesisStreamingDestinationInput$json = {
  '1': 'UpdateKinesisStreamingDestinationInput',
  '2': [
    {'1': 'streamarn', '3': 513423709, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'updatekinesisstreamingconfiguration',
      '3': 239134845,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.UpdateKinesisStreamingConfiguration',
      '10': 'updatekinesisstreamingconfiguration'
    },
  ],
};

/// Descriptor for `UpdateKinesisStreamingDestinationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKinesisStreamingDestinationInputDescriptor =
    $convert.base64Decode(
        'CiZVcGRhdGVLaW5lc2lzU3RyZWFtaW5nRGVzdGluYXRpb25JbnB1dBIgCglzdHJlYW1hcm4Y3f'
        'Lo9AEgASgJUglzdHJlYW1hcm4SIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1lEoIB'
        'CiN1cGRhdGVraW5lc2lzc3RyZWFtaW5nY29uZmlndXJhdGlvbhj90INyIAEoCzItLmR5bmFtb2'
        'RiLlVwZGF0ZUtpbmVzaXNTdHJlYW1pbmdDb25maWd1cmF0aW9uUiN1cGRhdGVraW5lc2lzc3Ry'
        'ZWFtaW5nY29uZmlndXJhdGlvbg==');

@$core
    .Deprecated('Use updateKinesisStreamingDestinationOutputDescriptor instead')
const UpdateKinesisStreamingDestinationOutput$json = {
  '1': 'UpdateKinesisStreamingDestinationOutput',
  '2': [
    {
      '1': 'destinationstatus',
      '3': 381248234,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.DestinationStatus',
      '10': 'destinationstatus'
    },
    {'1': 'streamarn', '3': 513423709, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'updatekinesisstreamingconfiguration',
      '3': 239134845,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.UpdateKinesisStreamingConfiguration',
      '10': 'updatekinesisstreamingconfiguration'
    },
  ],
};

/// Descriptor for `UpdateKinesisStreamingDestinationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKinesisStreamingDestinationOutputDescriptor =
    $convert.base64Decode(
        'CidVcGRhdGVLaW5lc2lzU3RyZWFtaW5nRGVzdGluYXRpb25PdXRwdXQSTQoRZGVzdGluYXRpb2'
        '5zdGF0dXMY6sXltQEgASgOMhsuZHluYW1vZGIuRGVzdGluYXRpb25TdGF0dXNSEWRlc3RpbmF0'
        'aW9uc3RhdHVzEiAKCXN0cmVhbWFybhjd8uj0ASABKAlSCXN0cmVhbWFybhIgCgl0YWJsZW5hbW'
        'UY3eTagQEgASgJUgl0YWJsZW5hbWUSggEKI3VwZGF0ZWtpbmVzaXNzdHJlYW1pbmdjb25maWd1'
        'cmF0aW9uGP3Qg3IgASgLMi0uZHluYW1vZGIuVXBkYXRlS2luZXNpc1N0cmVhbWluZ0NvbmZpZ3'
        'VyYXRpb25SI3VwZGF0ZWtpbmVzaXNzdHJlYW1pbmdjb25maWd1cmF0aW9u');

@$core.Deprecated('Use updateReplicationGroupMemberActionDescriptor instead')
const UpdateReplicationGroupMemberAction$json = {
  '1': 'UpdateReplicationGroupMemberAction',
  '2': [
    {
      '1': 'globalsecondaryindexes',
      '3': 409156905,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaGlobalSecondaryIndex',
      '10': 'globalsecondaryindexes'
    },
    {
      '1': 'kmsmasterkeyid',
      '3': 521459443,
      '4': 1,
      '5': 9,
      '10': 'kmsmasterkeyid'
    },
    {
      '1': 'ondemandthroughputoverride',
      '3': 317165234,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughputOverride',
      '10': 'ondemandthroughputoverride'
    },
    {
      '1': 'provisionedthroughputoverride',
      '3': 413332116,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughputOverride',
      '10': 'provisionedthroughputoverride'
    },
    {'1': 'regionname', '3': 112086463, '4': 1, '5': 9, '10': 'regionname'},
    {
      '1': 'tableclassoverride',
      '3': 415569842,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableClass',
      '10': 'tableclassoverride'
    },
  ],
};

/// Descriptor for `UpdateReplicationGroupMemberAction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateReplicationGroupMemberActionDescriptor = $convert.base64Decode(
    'CiJVcGRhdGVSZXBsaWNhdGlvbkdyb3VwTWVtYmVyQWN0aW9uEmEKFmdsb2JhbHNlY29uZGFyeW'
    'luZGV4ZXMYqfqMwwEgAygLMiUuZHluYW1vZGIuUmVwbGljYUdsb2JhbFNlY29uZGFyeUluZGV4'
    'UhZnbG9iYWxzZWNvbmRhcnlpbmRleGVzEioKDmttc21hc3RlcmtleWlkGPOt0/gBIAEoCVIOa2'
    '1zbWFzdGVya2V5aWQSaAoab25kZW1hbmR0aHJvdWdocHV0b3ZlcnJpZGUYsp2elwEgASgLMiQu'
    'ZHluYW1vZGIuT25EZW1hbmRUaHJvdWdocHV0T3ZlcnJpZGVSGm9uZGVtYW5kdGhyb3VnaHB1dG'
    '92ZXJyaWRlEnEKHXByb3Zpc2lvbmVkdGhyb3VnaHB1dG92ZXJyaWRlGJTli8UBIAEoCzInLmR5'
    'bmFtb2RiLlByb3Zpc2lvbmVkVGhyb3VnaHB1dE92ZXJyaWRlUh1wcm92aXNpb25lZHRocm91Z2'
    'hwdXRvdmVycmlkZRIhCgpyZWdpb25uYW1lGL+buTUgASgJUgpyZWdpb25uYW1lEkgKEnRhYmxl'
    'Y2xhc3NvdmVycmlkZRiyr5TGASABKA4yFC5keW5hbW9kYi5UYWJsZUNsYXNzUhJ0YWJsZWNsYX'
    'Nzb3ZlcnJpZGU=');

@$core.Deprecated('Use updateTableInputDescriptor instead')
const UpdateTableInput$json = {
  '1': 'UpdateTableInput',
  '2': [
    {
      '1': 'attributedefinitions',
      '3': 414687108,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.AttributeDefinition',
      '10': 'attributedefinitions'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'deletionprotectionenabled',
      '3': 259418450,
      '4': 1,
      '5': 8,
      '10': 'deletionprotectionenabled'
    },
    {
      '1': 'globalsecondaryindexupdates',
      '3': 265760923,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndexUpdate',
      '10': 'globalsecondaryindexupdates'
    },
    {
      '1': 'globaltablesettingsreplicationmode',
      '3': 10446577,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.GlobalTableSettingsReplicationMode',
      '10': 'globaltablesettingsreplicationmode'
    },
    {
      '1': 'globaltablewitnessupdates',
      '3': 269201458,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalTableWitnessGroupUpdate',
      '10': 'globaltablewitnessupdates'
    },
    {
      '1': 'multiregionconsistency',
      '3': 446019131,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.MultiRegionConsistency',
      '10': 'multiregionconsistency'
    },
    {
      '1': 'ondemandthroughput',
      '3': 481734402,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.OnDemandThroughput',
      '10': 'ondemandthroughput'
    },
    {
      '1': 'provisionedthroughput',
      '3': 1757580,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.ProvisionedThroughput',
      '10': 'provisionedthroughput'
    },
    {
      '1': 'replicaupdates',
      '3': 260731936,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicationGroupUpdate',
      '10': 'replicaupdates'
    },
    {
      '1': 'ssespecification',
      '3': 31692444,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.SSESpecification',
      '10': 'ssespecification'
    },
    {
      '1': 'streamspecification',
      '3': 403922627,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.StreamSpecification',
      '10': 'streamspecification'
    },
    {
      '1': 'tableclass',
      '3': 342890498,
      '4': 1,
      '5': 14,
      '6': '.dynamodb.TableClass',
      '10': 'tableclass'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.WarmThroughput',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `UpdateTableInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTableInputDescriptor = $convert.base64Decode(
    'ChBVcGRhdGVUYWJsZUlucHV0ElUKFGF0dHJpYnV0ZWRlZmluaXRpb25zGIS/3sUBIAMoCzIdLm'
    'R5bmFtb2RiLkF0dHJpYnV0ZURlZmluaXRpb25SFGF0dHJpYnV0ZWRlZmluaXRpb25zEjoKC2Jp'
    'bGxpbmdtb2RlGMC06FcgASgOMhUuZHluYW1vZGIuQmlsbGluZ01vZGVSC2JpbGxpbmdtb2RlEj'
    '8KGWRlbGV0aW9ucHJvdGVjdGlvbmVuYWJsZWQY0tLZeyABKAhSGWRlbGV0aW9ucHJvdGVjdGlv'
    'bmVuYWJsZWQSaQobZ2xvYmFsc2Vjb25kYXJ5aW5kZXh1cGRhdGVzGJvh3H4gAygLMiQuZHluYW'
    '1vZGIuR2xvYmFsU2Vjb25kYXJ5SW5kZXhVcGRhdGVSG2dsb2JhbHNlY29uZGFyeWluZGV4dXBk'
    'YXRlcxJ/CiJnbG9iYWx0YWJsZXNldHRpbmdzcmVwbGljYXRpb25tb2RlGPHN/QQgASgOMiwuZH'
    'luYW1vZGIuR2xvYmFsVGFibGVTZXR0aW5nc1JlcGxpY2F0aW9uTW9kZVIiZ2xvYmFsdGFibGVz'
    'ZXR0aW5nc3JlcGxpY2F0aW9ubW9kZRJpChlnbG9iYWx0YWJsZXdpdG5lc3N1cGRhdGVzGLLgro'
    'ABIAMoCzInLmR5bmFtb2RiLkdsb2JhbFRhYmxlV2l0bmVzc0dyb3VwVXBkYXRlUhlnbG9iYWx0'
    'YWJsZXdpdG5lc3N1cGRhdGVzElwKFm11bHRpcmVnaW9uY29uc2lzdGVuY3kYu+zW1AEgASgOMi'
    'AuZHluYW1vZGIuTXVsdGlSZWdpb25Db25zaXN0ZW5jeVIWbXVsdGlyZWdpb25jb25zaXN0ZW5j'
    'eRJQChJvbmRlbWFuZHRocm91Z2hwdXQYgt7a5QEgASgLMhwuZHluYW1vZGIuT25EZW1hbmRUaH'
    'JvdWdocHV0UhJvbmRlbWFuZHRocm91Z2hwdXQSVwoVcHJvdmlzaW9uZWR0aHJvdWdocHV0GIyj'
    'ayABKAsyHy5keW5hbW9kYi5Qcm92aXNpb25lZFRocm91Z2hwdXRSFXByb3Zpc2lvbmVkdGhyb3'
    'VnaHB1dBJLCg5yZXBsaWNhdXBkYXRlcxig6Kl8IAMoCzIgLmR5bmFtb2RiLlJlcGxpY2F0aW9u'
    'R3JvdXBVcGRhdGVSDnJlcGxpY2F1cGRhdGVzEkkKEHNzZXNwZWNpZmljYXRpb24YnK2ODyABKA'
    'syGi5keW5hbW9kYi5TU0VTcGVjaWZpY2F0aW9uUhBzc2VzcGVjaWZpY2F0aW9uElMKE3N0cmVh'
    'bXNwZWNpZmljYXRpb24Yw73NwAEgASgLMh0uZHluYW1vZGIuU3RyZWFtU3BlY2lmaWNhdGlvbl'
    'ITc3RyZWFtc3BlY2lmaWNhdGlvbhI4Cgp0YWJsZWNsYXNzGIKwwKMBIAEoDjIULmR5bmFtb2Ri'
    'LlRhYmxlQ2xhc3NSCnRhYmxlY2xhc3MSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW'
    '1lEkQKDndhcm10aHJvdWdocHV0GIPeyIoBIAEoCzIYLmR5bmFtb2RiLldhcm1UaHJvdWdocHV0'
    'Ug53YXJtdGhyb3VnaHB1dA==');

@$core.Deprecated('Use updateTableOutputDescriptor instead')
const UpdateTableOutput$json = {
  '1': 'UpdateTableOutput',
  '2': [
    {
      '1': 'tabledescription',
      '3': 19962388,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableDescription',
      '10': 'tabledescription'
    },
  ],
};

/// Descriptor for `UpdateTableOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTableOutputDescriptor = $convert.base64Decode(
    'ChFVcGRhdGVUYWJsZU91dHB1dBJJChB0YWJsZWRlc2NyaXB0aW9uGJS0wgkgASgLMhouZHluYW'
    '1vZGIuVGFibGVEZXNjcmlwdGlvblIQdGFibGVkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use updateTableReplicaAutoScalingInputDescriptor instead')
const UpdateTableReplicaAutoScalingInput$json = {
  '1': 'UpdateTableReplicaAutoScalingInput',
  '2': [
    {
      '1': 'globalsecondaryindexupdates',
      '3': 265760923,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.GlobalSecondaryIndexAutoScalingUpdate',
      '10': 'globalsecondaryindexupdates'
    },
    {
      '1': 'provisionedwritecapacityautoscalingupdate',
      '3': 316618294,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.AutoScalingSettingsUpdate',
      '10': 'provisionedwritecapacityautoscalingupdate'
    },
    {
      '1': 'replicaupdates',
      '3': 260731936,
      '4': 3,
      '5': 11,
      '6': '.dynamodb.ReplicaAutoScalingUpdate',
      '10': 'replicaupdates'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `UpdateTableReplicaAutoScalingInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTableReplicaAutoScalingInputDescriptor = $convert.base64Decode(
    'CiJVcGRhdGVUYWJsZVJlcGxpY2FBdXRvU2NhbGluZ0lucHV0EnQKG2dsb2JhbHNlY29uZGFyeW'
    'luZGV4dXBkYXRlcxib4dx+IAMoCzIvLmR5bmFtb2RiLkdsb2JhbFNlY29uZGFyeUluZGV4QXV0'
    'b1NjYWxpbmdVcGRhdGVSG2dsb2JhbHNlY29uZGFyeWluZGV4dXBkYXRlcxKFAQopcHJvdmlzaW'
    '9uZWR3cml0ZWNhcGFjaXR5YXV0b3NjYWxpbmd1cGRhdGUYtuz8lgEgASgLMiMuZHluYW1vZGIu'
    'QXV0b1NjYWxpbmdTZXR0aW5nc1VwZGF0ZVIpcHJvdmlzaW9uZWR3cml0ZWNhcGFjaXR5YXV0b3'
    'NjYWxpbmd1cGRhdGUSTQoOcmVwbGljYXVwZGF0ZXMYoOipfCADKAsyIi5keW5hbW9kYi5SZXBs'
    'aWNhQXV0b1NjYWxpbmdVcGRhdGVSDnJlcGxpY2F1cGRhdGVzEiAKCXRhYmxlbmFtZRjd5NqBAS'
    'ABKAlSCXRhYmxlbmFtZQ==');

@$core.Deprecated('Use updateTableReplicaAutoScalingOutputDescriptor instead')
const UpdateTableReplicaAutoScalingOutput$json = {
  '1': 'UpdateTableReplicaAutoScalingOutput',
  '2': [
    {
      '1': 'tableautoscalingdescription',
      '3': 490422806,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TableAutoScalingDescription',
      '10': 'tableautoscalingdescription'
    },
  ],
};

/// Descriptor for `UpdateTableReplicaAutoScalingOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTableReplicaAutoScalingOutputDescriptor =
    $convert.base64Decode(
        'CiNVcGRhdGVUYWJsZVJlcGxpY2FBdXRvU2NhbGluZ091dHB1dBJrCht0YWJsZWF1dG9zY2FsaW'
        '5nZGVzY3JpcHRpb24YloTt6QEgASgLMiUuZHluYW1vZGIuVGFibGVBdXRvU2NhbGluZ0Rlc2Ny'
        'aXB0aW9uUht0YWJsZWF1dG9zY2FsaW5nZGVzY3JpcHRpb24=');

@$core.Deprecated('Use updateTimeToLiveInputDescriptor instead')
const UpdateTimeToLiveInput$json = {
  '1': 'UpdateTimeToLiveInput',
  '2': [
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'timetolivespecification',
      '3': 315882193,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TimeToLiveSpecification',
      '10': 'timetolivespecification'
    },
  ],
};

/// Descriptor for `UpdateTimeToLiveInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTimeToLiveInputDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVUaW1lVG9MaXZlSW5wdXQSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW'
    '1lEl8KF3RpbWV0b2xpdmVzcGVjaWZpY2F0aW9uGNH1z5YBIAEoCzIhLmR5bmFtb2RiLlRpbWVU'
    'b0xpdmVTcGVjaWZpY2F0aW9uUhd0aW1ldG9saXZlc3BlY2lmaWNhdGlvbg==');

@$core.Deprecated('Use updateTimeToLiveOutputDescriptor instead')
const UpdateTimeToLiveOutput$json = {
  '1': 'UpdateTimeToLiveOutput',
  '2': [
    {
      '1': 'timetolivespecification',
      '3': 315882193,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.TimeToLiveSpecification',
      '10': 'timetolivespecification'
    },
  ],
};

/// Descriptor for `UpdateTimeToLiveOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTimeToLiveOutputDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVUaW1lVG9MaXZlT3V0cHV0El8KF3RpbWV0b2xpdmVzcGVjaWZpY2F0aW9uGNH1z5'
    'YBIAEoCzIhLmR5bmFtb2RiLlRpbWVUb0xpdmVTcGVjaWZpY2F0aW9uUhd0aW1ldG9saXZlc3Bl'
    'Y2lmaWNhdGlvbg==');

@$core.Deprecated('Use warmThroughputDescriptor instead')
const WarmThroughput$json = {
  '1': 'WarmThroughput',
  '2': [
    {
      '1': 'readunitspersecond',
      '3': 11400732,
      '4': 1,
      '5': 3,
      '10': 'readunitspersecond'
    },
    {
      '1': 'writeunitspersecond',
      '3': 339770127,
      '4': 1,
      '5': 3,
      '10': 'writeunitspersecond'
    },
  ],
};

/// Descriptor for `WarmThroughput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List warmThroughputDescriptor = $convert.base64Decode(
    'Cg5XYXJtVGhyb3VnaHB1dBIxChJyZWFkdW5pdHNwZXJzZWNvbmQYnOy3BSABKANSEnJlYWR1bm'
    'l0c3BlcnNlY29uZBI0ChN3cml0ZXVuaXRzcGVyc2Vjb25kGI/2gaIBIAEoA1ITd3JpdGV1bml0'
    'c3BlcnNlY29uZA==');

@$core.Deprecated('Use writeRequestDescriptor instead')
const WriteRequest$json = {
  '1': 'WriteRequest',
  '2': [
    {
      '1': 'deleterequest',
      '3': 477809064,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.DeleteRequest',
      '10': 'deleterequest'
    },
    {
      '1': 'putrequest',
      '3': 443252836,
      '4': 1,
      '5': 11,
      '6': '.dynamodb.PutRequest',
      '10': 'putrequest'
    },
  ],
};

/// Descriptor for `WriteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List writeRequestDescriptor = $convert.base64Decode(
    'CgxXcml0ZVJlcXVlc3QSQQoNZGVsZXRlcmVxdWVzdBiok+vjASABKAsyFy5keW5hbW9kYi5EZW'
    'xldGVSZXF1ZXN0Ug1kZWxldGVyZXF1ZXN0EjgKCnB1dHJlcXVlc3QY5ICu0wEgASgLMhQuZHlu'
    'YW1vZGIuUHV0UmVxdWVzdFIKcHV0cmVxdWVzdA==');
