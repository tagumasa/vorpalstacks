// This is a generated file - do not edit.
//
// Generated from logs.proto.

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

@$core.Deprecated('Use actionStatusDescriptor instead')
const ActionStatus$json = {
  '1': 'ActionStatus',
  '2': [
    {'1': 'ACTION_STATUS_IN_PROGRESS', '2': 0},
    {'1': 'ACTION_STATUS_COMPLETE', '2': 1},
    {'1': 'ACTION_STATUS_CLIENT_ERROR', '2': 2},
    {'1': 'ACTION_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `ActionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List actionStatusDescriptor = $convert.base64Decode(
    'CgxBY3Rpb25TdGF0dXMSHQoZQUNUSU9OX1NUQVRVU19JTl9QUk9HUkVTUxAAEhoKFkFDVElPTl'
    '9TVEFUVVNfQ09NUExFVEUQARIeChpBQ1RJT05fU1RBVFVTX0NMSUVOVF9FUlJPUhACEhgKFEFD'
    'VElPTl9TVEFUVVNfRkFJTEVEEAM=');

@$core.Deprecated('Use anomalyDetectorStatusDescriptor instead')
const AnomalyDetectorStatus$json = {
  '1': 'AnomalyDetectorStatus',
  '2': [
    {'1': 'ANOMALY_DETECTOR_STATUS_ANALYZING', '2': 0},
    {'1': 'ANOMALY_DETECTOR_STATUS_INITIALIZING', '2': 1},
    {'1': 'ANOMALY_DETECTOR_STATUS_TRAINING', '2': 2},
    {'1': 'ANOMALY_DETECTOR_STATUS_DELETED', '2': 3},
    {'1': 'ANOMALY_DETECTOR_STATUS_FAILED', '2': 4},
    {'1': 'ANOMALY_DETECTOR_STATUS_PAUSED', '2': 5},
  ],
};

/// Descriptor for `AnomalyDetectorStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List anomalyDetectorStatusDescriptor = $convert.base64Decode(
    'ChVBbm9tYWx5RGV0ZWN0b3JTdGF0dXMSJQohQU5PTUFMWV9ERVRFQ1RPUl9TVEFUVVNfQU5BTF'
    'laSU5HEAASKAokQU5PTUFMWV9ERVRFQ1RPUl9TVEFUVVNfSU5JVElBTElaSU5HEAESJAogQU5P'
    'TUFMWV9ERVRFQ1RPUl9TVEFUVVNfVFJBSU5JTkcQAhIjCh9BTk9NQUxZX0RFVEVDVE9SX1NUQV'
    'RVU19ERUxFVEVEEAMSIgoeQU5PTUFMWV9ERVRFQ1RPUl9TVEFUVVNfRkFJTEVEEAQSIgoeQU5P'
    'TUFMWV9ERVRFQ1RPUl9TVEFUVVNfUEFVU0VEEAU=');

@$core.Deprecated('Use dataProtectionStatusDescriptor instead')
const DataProtectionStatus$json = {
  '1': 'DataProtectionStatus',
  '2': [
    {'1': 'DATA_PROTECTION_STATUS_DISABLED', '2': 0},
    {'1': 'DATA_PROTECTION_STATUS_ACTIVATED', '2': 1},
    {'1': 'DATA_PROTECTION_STATUS_ARCHIVED', '2': 2},
    {'1': 'DATA_PROTECTION_STATUS_DELETED', '2': 3},
  ],
};

/// Descriptor for `DataProtectionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dataProtectionStatusDescriptor = $convert.base64Decode(
    'ChREYXRhUHJvdGVjdGlvblN0YXR1cxIjCh9EQVRBX1BST1RFQ1RJT05fU1RBVFVTX0RJU0FCTE'
    'VEEAASJAogREFUQV9QUk9URUNUSU9OX1NUQVRVU19BQ1RJVkFURUQQARIjCh9EQVRBX1BST1RF'
    'Q1RJT05fU1RBVFVTX0FSQ0hJVkVEEAISIgoeREFUQV9QUk9URUNUSU9OX1NUQVRVU19ERUxFVE'
    'VEEAM=');

@$core.Deprecated('Use deliveryDestinationTypeDescriptor instead')
const DeliveryDestinationType$json = {
  '1': 'DeliveryDestinationType',
  '2': [
    {'1': 'DELIVERY_DESTINATION_TYPE_XRAY', '2': 0},
    {'1': 'DELIVERY_DESTINATION_TYPE_S3', '2': 1},
    {'1': 'DELIVERY_DESTINATION_TYPE_CWL', '2': 2},
    {'1': 'DELIVERY_DESTINATION_TYPE_FH', '2': 3},
  ],
};

/// Descriptor for `DeliveryDestinationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deliveryDestinationTypeDescriptor = $convert.base64Decode(
    'ChdEZWxpdmVyeURlc3RpbmF0aW9uVHlwZRIiCh5ERUxJVkVSWV9ERVNUSU5BVElPTl9UWVBFX1'
    'hSQVkQABIgChxERUxJVkVSWV9ERVNUSU5BVElPTl9UWVBFX1MzEAESIQodREVMSVZFUllfREVT'
    'VElOQVRJT05fVFlQRV9DV0wQAhIgChxERUxJVkVSWV9ERVNUSU5BVElPTl9UWVBFX0ZIEAM=');

@$core.Deprecated('Use distributionDescriptor instead')
const Distribution$json = {
  '1': 'Distribution',
  '2': [
    {'1': 'DISTRIBUTION_BYLOGSTREAM', '2': 0},
    {'1': 'DISTRIBUTION_RANDOM', '2': 1},
  ],
};

/// Descriptor for `Distribution`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List distributionDescriptor = $convert.base64Decode(
    'CgxEaXN0cmlidXRpb24SHAoYRElTVFJJQlVUSU9OX0JZTE9HU1RSRUFNEAASFwoTRElTVFJJQl'
    'VUSU9OX1JBTkRPTRAB');

@$core.Deprecated('Use entityRejectionErrorTypeDescriptor instead')
const EntityRejectionErrorType$json = {
  '1': 'EntityRejectionErrorType',
  '2': [
    {'1': 'ENTITY_REJECTION_ERROR_TYPE_ENTITY_SIZE_TOO_LARGE', '2': 0},
    {'1': 'ENTITY_REJECTION_ERROR_TYPE_INVALID_KEY_ATTRIBUTE', '2': 1},
    {'1': 'ENTITY_REJECTION_ERROR_TYPE_UNSUPPORTED_LOG_GROUP_TYPE', '2': 2},
    {'1': 'ENTITY_REJECTION_ERROR_TYPE_MISSING_REQUIRED_FIELDS', '2': 3},
    {'1': 'ENTITY_REJECTION_ERROR_TYPE_INVALID_ATTRIBUTES', '2': 4},
    {'1': 'ENTITY_REJECTION_ERROR_TYPE_INVALID_TYPE_VALUE', '2': 5},
    {'1': 'ENTITY_REJECTION_ERROR_TYPE_INVALID_ENTITY', '2': 6},
  ],
};

/// Descriptor for `EntityRejectionErrorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List entityRejectionErrorTypeDescriptor = $convert.base64Decode(
    'ChhFbnRpdHlSZWplY3Rpb25FcnJvclR5cGUSNQoxRU5USVRZX1JFSkVDVElPTl9FUlJPUl9UWV'
    'BFX0VOVElUWV9TSVpFX1RPT19MQVJHRRAAEjUKMUVOVElUWV9SRUpFQ1RJT05fRVJST1JfVFlQ'
    'RV9JTlZBTElEX0tFWV9BVFRSSUJVVEUQARI6CjZFTlRJVFlfUkVKRUNUSU9OX0VSUk9SX1RZUE'
    'VfVU5TVVBQT1JURURfTE9HX0dST1VQX1RZUEUQAhI3CjNFTlRJVFlfUkVKRUNUSU9OX0VSUk9S'
    'X1RZUEVfTUlTU0lOR19SRVFVSVJFRF9GSUVMRFMQAxIyCi5FTlRJVFlfUkVKRUNUSU9OX0VSUk'
    '9SX1RZUEVfSU5WQUxJRF9BVFRSSUJVVEVTEAQSMgouRU5USVRZX1JFSkVDVElPTl9FUlJPUl9U'
    'WVBFX0lOVkFMSURfVFlQRV9WQUxVRRAFEi4KKkVOVElUWV9SRUpFQ1RJT05fRVJST1JfVFlQRV'
    '9JTlZBTElEX0VOVElUWRAG');

@$core.Deprecated('Use evaluationFrequencyDescriptor instead')
const EvaluationFrequency$json = {
  '1': 'EvaluationFrequency',
  '2': [
    {'1': 'EVALUATION_FREQUENCY_ONE_MIN', '2': 0},
    {'1': 'EVALUATION_FREQUENCY_FIVE_MIN', '2': 1},
    {'1': 'EVALUATION_FREQUENCY_ONE_HOUR', '2': 2},
    {'1': 'EVALUATION_FREQUENCY_FIFTEEN_MIN', '2': 3},
    {'1': 'EVALUATION_FREQUENCY_THIRTY_MIN', '2': 4},
    {'1': 'EVALUATION_FREQUENCY_TEN_MIN', '2': 5},
  ],
};

/// Descriptor for `EvaluationFrequency`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List evaluationFrequencyDescriptor = $convert.base64Decode(
    'ChNFdmFsdWF0aW9uRnJlcXVlbmN5EiAKHEVWQUxVQVRJT05fRlJFUVVFTkNZX09ORV9NSU4QAB'
    'IhCh1FVkFMVUFUSU9OX0ZSRVFVRU5DWV9GSVZFX01JThABEiEKHUVWQUxVQVRJT05fRlJFUVVF'
    'TkNZX09ORV9IT1VSEAISJAogRVZBTFVBVElPTl9GUkVRVUVOQ1lfRklGVEVFTl9NSU4QAxIjCh'
    '9FVkFMVUFUSU9OX0ZSRVFVRU5DWV9USElSVFlfTUlOEAQSIAocRVZBTFVBVElPTl9GUkVRVUVO'
    'Q1lfVEVOX01JThAF');

@$core.Deprecated('Use eventSourceDescriptor instead')
const EventSource$json = {
  '1': 'EventSource',
  '2': [
    {'1': 'EVENT_SOURCE_CLOUD_TRAIL', '2': 0},
    {'1': 'EVENT_SOURCE_AWSWAF', '2': 1},
    {'1': 'EVENT_SOURCE_VPC_FLOW', '2': 2},
    {'1': 'EVENT_SOURCE_ROUTE53_RESOLVER', '2': 3},
    {'1': 'EVENT_SOURCE_EKS_AUDIT', '2': 4},
  ],
};

/// Descriptor for `EventSource`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventSourceDescriptor = $convert.base64Decode(
    'CgtFdmVudFNvdXJjZRIcChhFVkVOVF9TT1VSQ0VfQ0xPVURfVFJBSUwQABIXChNFVkVOVF9TT1'
    'VSQ0VfQVdTV0FGEAESGQoVRVZFTlRfU09VUkNFX1ZQQ19GTE9XEAISIQodRVZFTlRfU09VUkNF'
    'X1JPVVRFNTNfUkVTT0xWRVIQAxIaChZFVkVOVF9TT1VSQ0VfRUtTX0FVRElUEAQ=');

@$core.Deprecated('Use executionStatusDescriptor instead')
const ExecutionStatus$json = {
  '1': 'ExecutionStatus',
  '2': [
    {'1': 'EXECUTION_STATUS_FAILED', '2': 0},
    {'1': 'EXECUTION_STATUS_RUNNING', '2': 1},
    {'1': 'EXECUTION_STATUS_COMPLETE', '2': 2},
    {'1': 'EXECUTION_STATUS_TIMEOUT', '2': 3},
    {'1': 'EXECUTION_STATUS_INVALIDQUERY', '2': 4},
  ],
};

/// Descriptor for `ExecutionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executionStatusDescriptor = $convert.base64Decode(
    'Cg9FeGVjdXRpb25TdGF0dXMSGwoXRVhFQ1VUSU9OX1NUQVRVU19GQUlMRUQQABIcChhFWEVDVV'
    'RJT05fU1RBVFVTX1JVTk5JTkcQARIdChlFWEVDVVRJT05fU1RBVFVTX0NPTVBMRVRFEAISHAoY'
    'RVhFQ1VUSU9OX1NUQVRVU19USU1FT1VUEAMSIQodRVhFQ1VUSU9OX1NUQVRVU19JTlZBTElEUV'
    'VFUlkQBA==');

@$core.Deprecated('Use exportTaskStatusCodeDescriptor instead')
const ExportTaskStatusCode$json = {
  '1': 'ExportTaskStatusCode',
  '2': [
    {'1': 'EXPORT_TASK_STATUS_CODE_PENDING', '2': 0},
    {'1': 'EXPORT_TASK_STATUS_CODE_RUNNING', '2': 1},
    {'1': 'EXPORT_TASK_STATUS_CODE_CANCELLED', '2': 2},
    {'1': 'EXPORT_TASK_STATUS_CODE_PENDING_CANCEL', '2': 3},
    {'1': 'EXPORT_TASK_STATUS_CODE_COMPLETED', '2': 4},
    {'1': 'EXPORT_TASK_STATUS_CODE_FAILED', '2': 5},
  ],
};

/// Descriptor for `ExportTaskStatusCode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List exportTaskStatusCodeDescriptor = $convert.base64Decode(
    'ChRFeHBvcnRUYXNrU3RhdHVzQ29kZRIjCh9FWFBPUlRfVEFTS19TVEFUVVNfQ09ERV9QRU5ESU'
    '5HEAASIwofRVhQT1JUX1RBU0tfU1RBVFVTX0NPREVfUlVOTklORxABEiUKIUVYUE9SVF9UQVNL'
    'X1NUQVRVU19DT0RFX0NBTkNFTExFRBACEioKJkVYUE9SVF9UQVNLX1NUQVRVU19DT0RFX1BFTk'
    'RJTkdfQ0FOQ0VMEAMSJQohRVhQT1JUX1RBU0tfU1RBVFVTX0NPREVfQ09NUExFVEVEEAQSIgoe'
    'RVhQT1JUX1RBU0tfU1RBVFVTX0NPREVfRkFJTEVEEAU=');

@$core.Deprecated('Use flattenedElementDescriptor instead')
const FlattenedElement$json = {
  '1': 'FlattenedElement',
  '2': [
    {'1': 'FLATTENED_ELEMENT_LAST', '2': 0},
    {'1': 'FLATTENED_ELEMENT_FIRST', '2': 1},
  ],
};

/// Descriptor for `FlattenedElement`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List flattenedElementDescriptor = $convert.base64Decode(
    'ChBGbGF0dGVuZWRFbGVtZW50EhoKFkZMQVRURU5FRF9FTEVNRU5UX0xBU1QQABIbChdGTEFUVE'
    'VORURfRUxFTUVOVF9GSVJTVBAB');

@$core.Deprecated('Use importStatusDescriptor instead')
const ImportStatus$json = {
  '1': 'ImportStatus',
  '2': [
    {'1': 'IMPORT_STATUS_IN_PROGRESS', '2': 0},
    {'1': 'IMPORT_STATUS_CANCELLED', '2': 1},
    {'1': 'IMPORT_STATUS_COMPLETED', '2': 2},
    {'1': 'IMPORT_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `ImportStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importStatusDescriptor = $convert.base64Decode(
    'CgxJbXBvcnRTdGF0dXMSHQoZSU1QT1JUX1NUQVRVU19JTl9QUk9HUkVTUxAAEhsKF0lNUE9SVF'
    '9TVEFUVVNfQ0FOQ0VMTEVEEAESGwoXSU1QT1JUX1NUQVRVU19DT01QTEVURUQQAhIYChRJTVBP'
    'UlRfU1RBVFVTX0ZBSUxFRBAD');

@$core.Deprecated('Use indexSourceDescriptor instead')
const IndexSource$json = {
  '1': 'IndexSource',
  '2': [
    {'1': 'INDEX_SOURCE_ACCOUNT', '2': 0},
    {'1': 'INDEX_SOURCE_LOG_GROUP', '2': 1},
  ],
};

/// Descriptor for `IndexSource`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List indexSourceDescriptor = $convert.base64Decode(
    'CgtJbmRleFNvdXJjZRIYChRJTkRFWF9TT1VSQ0VfQUNDT1VOVBAAEhoKFklOREVYX1NPVVJDRV'
    '9MT0dfR1JPVVAQAQ==');

@$core.Deprecated('Use indexTypeDescriptor instead')
const IndexType$json = {
  '1': 'IndexType',
  '2': [
    {'1': 'INDEX_TYPE_FACET', '2': 0},
    {'1': 'INDEX_TYPE_FIELD_INDEX', '2': 1},
  ],
};

/// Descriptor for `IndexType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List indexTypeDescriptor = $convert.base64Decode(
    'CglJbmRleFR5cGUSFAoQSU5ERVhfVFlQRV9GQUNFVBAAEhoKFklOREVYX1RZUEVfRklFTERfSU'
    '5ERVgQAQ==');

@$core.Deprecated('Use inheritedPropertyDescriptor instead')
const InheritedProperty$json = {
  '1': 'InheritedProperty',
  '2': [
    {'1': 'INHERITED_PROPERTY_ACCOUNT_DATA_PROTECTION', '2': 0},
  ],
};

/// Descriptor for `InheritedProperty`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inheritedPropertyDescriptor = $convert.base64Decode(
    'ChFJbmhlcml0ZWRQcm9wZXJ0eRIuCipJTkhFUklURURfUFJPUEVSVFlfQUNDT1VOVF9EQVRBX1'
    'BST1RFQ1RJT04QAA==');

@$core.Deprecated('Use integrationStatusDescriptor instead')
const IntegrationStatus$json = {
  '1': 'IntegrationStatus',
  '2': [
    {'1': 'INTEGRATION_STATUS_PROVISIONING', '2': 0},
    {'1': 'INTEGRATION_STATUS_ACTIVE', '2': 1},
    {'1': 'INTEGRATION_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `IntegrationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List integrationStatusDescriptor = $convert.base64Decode(
    'ChFJbnRlZ3JhdGlvblN0YXR1cxIjCh9JTlRFR1JBVElPTl9TVEFUVVNfUFJPVklTSU9OSU5HEA'
    'ASHQoZSU5URUdSQVRJT05fU1RBVFVTX0FDVElWRRABEh0KGUlOVEVHUkFUSU9OX1NUQVRVU19G'
    'QUlMRUQQAg==');

@$core.Deprecated('Use integrationTypeDescriptor instead')
const IntegrationType$json = {
  '1': 'IntegrationType',
  '2': [
    {'1': 'INTEGRATION_TYPE_OPENSEARCH', '2': 0},
  ],
};

/// Descriptor for `IntegrationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List integrationTypeDescriptor = $convert.base64Decode(
    'Cg9JbnRlZ3JhdGlvblR5cGUSHwobSU5URUdSQVRJT05fVFlQRV9PUEVOU0VBUkNIEAA=');

@$core.Deprecated('Use listAggregateLogGroupSummariesGroupByDescriptor instead')
const ListAggregateLogGroupSummariesGroupBy$json = {
  '1': 'ListAggregateLogGroupSummariesGroupBy',
  '2': [
    {
      '1':
          'LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_AND_TYPE',
      '2': 0
    },
    {
      '1':
          'LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_TYPE_AND_FORMAT',
      '2': 1
    },
  ],
};

/// Descriptor for `ListAggregateLogGroupSummariesGroupBy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List listAggregateLogGroupSummariesGroupByDescriptor =
    $convert.base64Decode(
        'CiVMaXN0QWdncmVnYXRlTG9nR3JvdXBTdW1tYXJpZXNHcm91cEJ5EkkKRUxJU1RfQUdHUkVHQV'
        'RFX0xPR19HUk9VUF9TVU1NQVJJRVNfR1JPVVBfQllfREFUQV9TT1VSQ0VfTkFNRV9BTkRfVFlQ'
        'RRAAElAKTExJU1RfQUdHUkVHQVRFX0xPR19HUk9VUF9TVU1NQVJJRVNfR1JPVVBfQllfREFUQV'
        '9TT1VSQ0VfTkFNRV9UWVBFX0FORF9GT1JNQVQQAQ==');

@$core.Deprecated('Use logGroupClassDescriptor instead')
const LogGroupClass$json = {
  '1': 'LogGroupClass',
  '2': [
    {'1': 'LOG_GROUP_CLASS_STANDARD', '2': 0},
    {'1': 'LOG_GROUP_CLASS_DELIVERY', '2': 1},
    {'1': 'LOG_GROUP_CLASS_INFREQUENT_ACCESS', '2': 2},
  ],
};

/// Descriptor for `LogGroupClass`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List logGroupClassDescriptor = $convert.base64Decode(
    'Cg1Mb2dHcm91cENsYXNzEhwKGExPR19HUk9VUF9DTEFTU19TVEFOREFSRBAAEhwKGExPR19HUk'
    '9VUF9DTEFTU19ERUxJVkVSWRABEiUKIUxPR19HUk9VUF9DTEFTU19JTkZSRVFVRU5UX0FDQ0VT'
    'UxAC');

@$core.Deprecated('Use oCSFVersionDescriptor instead')
const OCSFVersion$json = {
  '1': 'OCSFVersion',
  '2': [
    {'1': 'O_C_S_F_VERSION_V1_5', '2': 0},
    {'1': 'O_C_S_F_VERSION_V1_1', '2': 1},
  ],
};

/// Descriptor for `OCSFVersion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List oCSFVersionDescriptor = $convert.base64Decode(
    'CgtPQ1NGVmVyc2lvbhIYChRPX0NfU19GX1ZFUlNJT05fVjFfNRAAEhgKFE9fQ19TX0ZfVkVSU0'
    'lPTl9WMV8xEAE=');

@$core.Deprecated('Use openSearchResourceStatusTypeDescriptor instead')
const OpenSearchResourceStatusType$json = {
  '1': 'OpenSearchResourceStatusType',
  '2': [
    {'1': 'OPEN_SEARCH_RESOURCE_STATUS_TYPE_NOT_FOUND', '2': 0},
    {'1': 'OPEN_SEARCH_RESOURCE_STATUS_TYPE_ACTIVE', '2': 1},
    {'1': 'OPEN_SEARCH_RESOURCE_STATUS_TYPE_ERROR', '2': 2},
  ],
};

/// Descriptor for `OpenSearchResourceStatusType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List openSearchResourceStatusTypeDescriptor = $convert.base64Decode(
    'ChxPcGVuU2VhcmNoUmVzb3VyY2VTdGF0dXNUeXBlEi4KKk9QRU5fU0VBUkNIX1JFU09VUkNFX1'
    'NUQVRVU19UWVBFX05PVF9GT1VORBAAEisKJ09QRU5fU0VBUkNIX1JFU09VUkNFX1NUQVRVU19U'
    'WVBFX0FDVElWRRABEioKJk9QRU5fU0VBUkNIX1JFU09VUkNFX1NUQVRVU19UWVBFX0VSUk9SEA'
    'I=');

@$core.Deprecated('Use orderByDescriptor instead')
const OrderBy$json = {
  '1': 'OrderBy',
  '2': [
    {'1': 'ORDER_BY_LOGSTREAMNAME', '2': 0},
    {'1': 'ORDER_BY_LASTEVENTTIME', '2': 1},
  ],
};

/// Descriptor for `OrderBy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List orderByDescriptor = $convert.base64Decode(
    'CgdPcmRlckJ5EhoKFk9SREVSX0JZX0xPR1NUUkVBTU5BTUUQABIaChZPUkRFUl9CWV9MQVNURV'
    'ZFTlRUSU1FEAE=');

@$core.Deprecated('Use outputFormatDescriptor instead')
const OutputFormat$json = {
  '1': 'OutputFormat',
  '2': [
    {'1': 'OUTPUT_FORMAT_RAW', '2': 0},
    {'1': 'OUTPUT_FORMAT_JSON', '2': 1},
    {'1': 'OUTPUT_FORMAT_PARQUET', '2': 2},
    {'1': 'OUTPUT_FORMAT_PLAIN', '2': 3},
    {'1': 'OUTPUT_FORMAT_W3C', '2': 4},
  ],
};

/// Descriptor for `OutputFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List outputFormatDescriptor = $convert.base64Decode(
    'CgxPdXRwdXRGb3JtYXQSFQoRT1VUUFVUX0ZPUk1BVF9SQVcQABIWChJPVVRQVVRfRk9STUFUX0'
    'pTT04QARIZChVPVVRQVVRfRk9STUFUX1BBUlFVRVQQAhIXChNPVVRQVVRfRk9STUFUX1BMQUlO'
    'EAMSFQoRT1VUUFVUX0ZPUk1BVF9XM0MQBA==');

@$core.Deprecated('Use policyScopeDescriptor instead')
const PolicyScope$json = {
  '1': 'PolicyScope',
  '2': [
    {'1': 'POLICY_SCOPE_ACCOUNT', '2': 0},
    {'1': 'POLICY_SCOPE_RESOURCE', '2': 1},
  ],
};

/// Descriptor for `PolicyScope`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyScopeDescriptor = $convert.base64Decode(
    'CgtQb2xpY3lTY29wZRIYChRQT0xJQ1lfU0NPUEVfQUNDT1VOVBAAEhkKFVBPTElDWV9TQ09QRV'
    '9SRVNPVVJDRRAB');

@$core.Deprecated('Use policyTypeDescriptor instead')
const PolicyType$json = {
  '1': 'PolicyType',
  '2': [
    {'1': 'POLICY_TYPE_FIELD_INDEX_POLICY', '2': 0},
    {'1': 'POLICY_TYPE_DATA_PROTECTION_POLICY', '2': 1},
    {'1': 'POLICY_TYPE_TRANSFORMER_POLICY', '2': 2},
    {'1': 'POLICY_TYPE_METRIC_EXTRACTION_POLICY', '2': 3},
    {'1': 'POLICY_TYPE_SUBSCRIPTION_FILTER_POLICY', '2': 4},
  ],
};

/// Descriptor for `PolicyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List policyTypeDescriptor = $convert.base64Decode(
    'CgpQb2xpY3lUeXBlEiIKHlBPTElDWV9UWVBFX0ZJRUxEX0lOREVYX1BPTElDWRAAEiYKIlBPTE'
    'lDWV9UWVBFX0RBVEFfUFJPVEVDVElPTl9QT0xJQ1kQARIiCh5QT0xJQ1lfVFlQRV9UUkFOU0ZP'
    'Uk1FUl9QT0xJQ1kQAhIoCiRQT0xJQ1lfVFlQRV9NRVRSSUNfRVhUUkFDVElPTl9QT0xJQ1kQAx'
    'IqCiZQT0xJQ1lfVFlQRV9TVUJTQ1JJUFRJT05fRklMVEVSX1BPTElDWRAE');

@$core.Deprecated('Use queryLanguageDescriptor instead')
const QueryLanguage$json = {
  '1': 'QueryLanguage',
  '2': [
    {'1': 'QUERY_LANGUAGE_CWLI', '2': 0},
    {'1': 'QUERY_LANGUAGE_SQL', '2': 1},
    {'1': 'QUERY_LANGUAGE_PPL', '2': 2},
  ],
};

/// Descriptor for `QueryLanguage`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryLanguageDescriptor = $convert.base64Decode(
    'Cg1RdWVyeUxhbmd1YWdlEhcKE1FVRVJZX0xBTkdVQUdFX0NXTEkQABIWChJRVUVSWV9MQU5HVU'
    'FHRV9TUUwQARIWChJRVUVSWV9MQU5HVUFHRV9QUEwQAg==');

@$core.Deprecated('Use queryStatusDescriptor instead')
const QueryStatus$json = {
  '1': 'QueryStatus',
  '2': [
    {'1': 'QUERY_STATUS_FAILED', '2': 0},
    {'1': 'QUERY_STATUS_CANCELLED', '2': 1},
    {'1': 'QUERY_STATUS_RUNNING', '2': 2},
    {'1': 'QUERY_STATUS_COMPLETE', '2': 3},
    {'1': 'QUERY_STATUS_TIMEOUT', '2': 4},
    {'1': 'QUERY_STATUS_SCHEDULED', '2': 5},
    {'1': 'QUERY_STATUS_UNKNOWN', '2': 6},
  ],
};

/// Descriptor for `QueryStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryStatusDescriptor = $convert.base64Decode(
    'CgtRdWVyeVN0YXR1cxIXChNRVUVSWV9TVEFUVVNfRkFJTEVEEAASGgoWUVVFUllfU1RBVFVTX0'
    'NBTkNFTExFRBABEhgKFFFVRVJZX1NUQVRVU19SVU5OSU5HEAISGQoVUVVFUllfU1RBVFVTX0NP'
    'TVBMRVRFEAMSGAoUUVVFUllfU1RBVFVTX1RJTUVPVVQQBBIaChZRVUVSWV9TVEFUVVNfU0NIRU'
    'RVTEVEEAUSGAoUUVVFUllfU1RBVFVTX1VOS05PV04QBg==');

@$core.Deprecated('Use s3TableIntegrationSourceStatusDescriptor instead')
const S3TableIntegrationSourceStatus$json = {
  '1': 'S3TableIntegrationSourceStatus',
  '2': [
    {'1': 'S3_TABLE_INTEGRATION_SOURCE_STATUS_UNHEALTHY', '2': 0},
    {
      '1': 'S3_TABLE_INTEGRATION_SOURCE_STATUS_DATA_SOURCE_DELETE_IN_PROGRESS',
      '2': 1
    },
    {'1': 'S3_TABLE_INTEGRATION_SOURCE_STATUS_ACTIVE', '2': 2},
    {'1': 'S3_TABLE_INTEGRATION_SOURCE_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `S3TableIntegrationSourceStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List s3TableIntegrationSourceStatusDescriptor = $convert.base64Decode(
    'Ch5TM1RhYmxlSW50ZWdyYXRpb25Tb3VyY2VTdGF0dXMSMAosUzNfVEFCTEVfSU5URUdSQVRJT0'
    '5fU09VUkNFX1NUQVRVU19VTkhFQUxUSFkQABJFCkFTM19UQUJMRV9JTlRFR1JBVElPTl9TT1VS'
    'Q0VfU1RBVFVTX0RBVEFfU09VUkNFX0RFTEVURV9JTl9QUk9HUkVTUxABEi0KKVMzX1RBQkxFX0'
    'lOVEVHUkFUSU9OX1NPVVJDRV9TVEFUVVNfQUNUSVZFEAISLQopUzNfVEFCTEVfSU5URUdSQVRJ'
    'T05fU09VUkNFX1NUQVRVU19GQUlMRUQQAw==');

@$core.Deprecated('Use scheduledQueryDestinationTypeDescriptor instead')
const ScheduledQueryDestinationType$json = {
  '1': 'ScheduledQueryDestinationType',
  '2': [
    {'1': 'SCHEDULED_QUERY_DESTINATION_TYPE_S3', '2': 0},
  ],
};

/// Descriptor for `ScheduledQueryDestinationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scheduledQueryDestinationTypeDescriptor =
    $convert.base64Decode(
        'Ch1TY2hlZHVsZWRRdWVyeURlc3RpbmF0aW9uVHlwZRInCiNTQ0hFRFVMRURfUVVFUllfREVTVE'
        'lOQVRJT05fVFlQRV9TMxAA');

@$core.Deprecated('Use scheduledQueryStateDescriptor instead')
const ScheduledQueryState$json = {
  '1': 'ScheduledQueryState',
  '2': [
    {'1': 'SCHEDULED_QUERY_STATE_DISABLED', '2': 0},
    {'1': 'SCHEDULED_QUERY_STATE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ScheduledQueryState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scheduledQueryStateDescriptor = $convert.base64Decode(
    'ChNTY2hlZHVsZWRRdWVyeVN0YXRlEiIKHlNDSEVEVUxFRF9RVUVSWV9TVEFURV9ESVNBQkxFRB'
    'AAEiEKHVNDSEVEVUxFRF9RVUVSWV9TVEFURV9FTkFCTEVEEAE=');

@$core.Deprecated('Use scopeDescriptor instead')
const Scope$json = {
  '1': 'Scope',
  '2': [
    {'1': 'SCOPE_ALL', '2': 0},
  ],
};

/// Descriptor for `Scope`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scopeDescriptor =
    $convert.base64Decode('CgVTY29wZRINCglTQ09QRV9BTEwQAA==');

@$core.Deprecated('Use standardUnitDescriptor instead')
const StandardUnit$json = {
  '1': 'StandardUnit',
  '2': [
    {'1': 'STANDARD_UNIT_KILOBITSSECOND', '2': 0},
    {'1': 'STANDARD_UNIT_COUNTSECOND', '2': 1},
    {'1': 'STANDARD_UNIT_MEGABITSSECOND', '2': 2},
    {'1': 'STANDARD_UNIT_GIGABITS', '2': 3},
    {'1': 'STANDARD_UNIT_NONE', '2': 4},
    {'1': 'STANDARD_UNIT_MILLISECONDS', '2': 5},
    {'1': 'STANDARD_UNIT_BYTES', '2': 6},
    {'1': 'STANDARD_UNIT_GIGABYTES', '2': 7},
    {'1': 'STANDARD_UNIT_TERABITSSECOND', '2': 8},
    {'1': 'STANDARD_UNIT_BITS', '2': 9},
    {'1': 'STANDARD_UNIT_GIGABITSSECOND', '2': 10},
    {'1': 'STANDARD_UNIT_MEGABYTESSECOND', '2': 11},
    {'1': 'STANDARD_UNIT_MEGABYTES', '2': 12},
    {'1': 'STANDARD_UNIT_BITSSECOND', '2': 13},
    {'1': 'STANDARD_UNIT_GIGABYTESSECOND', '2': 14},
    {'1': 'STANDARD_UNIT_MEGABITS', '2': 15},
    {'1': 'STANDARD_UNIT_SECONDS', '2': 16},
    {'1': 'STANDARD_UNIT_TERABYTESSECOND', '2': 17},
    {'1': 'STANDARD_UNIT_PERCENT', '2': 18},
    {'1': 'STANDARD_UNIT_COUNT', '2': 19},
    {'1': 'STANDARD_UNIT_KILOBITS', '2': 20},
    {'1': 'STANDARD_UNIT_KILOBYTES', '2': 21},
    {'1': 'STANDARD_UNIT_TERABYTES', '2': 22},
    {'1': 'STANDARD_UNIT_MICROSECONDS', '2': 23},
    {'1': 'STANDARD_UNIT_KILOBYTESSECOND', '2': 24},
    {'1': 'STANDARD_UNIT_BYTESSECOND', '2': 25},
    {'1': 'STANDARD_UNIT_TERABITS', '2': 26},
  ],
};

/// Descriptor for `StandardUnit`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List standardUnitDescriptor = $convert.base64Decode(
    'CgxTdGFuZGFyZFVuaXQSIAocU1RBTkRBUkRfVU5JVF9LSUxPQklUU1NFQ09ORBAAEh0KGVNUQU'
    '5EQVJEX1VOSVRfQ09VTlRTRUNPTkQQARIgChxTVEFOREFSRF9VTklUX01FR0FCSVRTU0VDT05E'
    'EAISGgoWU1RBTkRBUkRfVU5JVF9HSUdBQklUUxADEhYKElNUQU5EQVJEX1VOSVRfTk9ORRAEEh'
    '4KGlNUQU5EQVJEX1VOSVRfTUlMTElTRUNPTkRTEAUSFwoTU1RBTkRBUkRfVU5JVF9CWVRFUxAG'
    'EhsKF1NUQU5EQVJEX1VOSVRfR0lHQUJZVEVTEAcSIAocU1RBTkRBUkRfVU5JVF9URVJBQklUU1'
    'NFQ09ORBAIEhYKElNUQU5EQVJEX1VOSVRfQklUUxAJEiAKHFNUQU5EQVJEX1VOSVRfR0lHQUJJ'
    'VFNTRUNPTkQQChIhCh1TVEFOREFSRF9VTklUX01FR0FCWVRFU1NFQ09ORBALEhsKF1NUQU5EQV'
    'JEX1VOSVRfTUVHQUJZVEVTEAwSHAoYU1RBTkRBUkRfVU5JVF9CSVRTU0VDT05EEA0SIQodU1RB'
    'TkRBUkRfVU5JVF9HSUdBQllURVNTRUNPTkQQDhIaChZTVEFOREFSRF9VTklUX01FR0FCSVRTEA'
    '8SGQoVU1RBTkRBUkRfVU5JVF9TRUNPTkRTEBASIQodU1RBTkRBUkRfVU5JVF9URVJBQllURVNT'
    'RUNPTkQQERIZChVTVEFOREFSRF9VTklUX1BFUkNFTlQQEhIXChNTVEFOREFSRF9VTklUX0NPVU'
    '5UEBMSGgoWU1RBTkRBUkRfVU5JVF9LSUxPQklUUxAUEhsKF1NUQU5EQVJEX1VOSVRfS0lMT0JZ'
    'VEVTEBUSGwoXU1RBTkRBUkRfVU5JVF9URVJBQllURVMQFhIeChpTVEFOREFSRF9VTklUX01JQ1'
    'JPU0VDT05EUxAXEiEKHVNUQU5EQVJEX1VOSVRfS0lMT0JZVEVTU0VDT05EEBgSHQoZU1RBTkRB'
    'UkRfVU5JVF9CWVRFU1NFQ09ORBAZEhoKFlNUQU5EQVJEX1VOSVRfVEVSQUJJVFMQGg==');

@$core.Deprecated('Use stateDescriptor instead')
const State$json = {
  '1': 'State',
  '2': [
    {'1': 'STATE_ACTIVE', '2': 0},
    {'1': 'STATE_SUPPRESSED', '2': 1},
    {'1': 'STATE_BASELINE', '2': 2},
  ],
};

/// Descriptor for `State`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stateDescriptor = $convert.base64Decode(
    'CgVTdGF0ZRIQCgxTVEFURV9BQ1RJVkUQABIUChBTVEFURV9TVVBQUkVTU0VEEAESEgoOU1RBVE'
    'VfQkFTRUxJTkUQAg==');

@$core.Deprecated('Use suppressionStateDescriptor instead')
const SuppressionState$json = {
  '1': 'SuppressionState',
  '2': [
    {'1': 'SUPPRESSION_STATE_UNSUPPRESSED', '2': 0},
    {'1': 'SUPPRESSION_STATE_SUPPRESSED', '2': 1},
  ],
};

/// Descriptor for `SuppressionState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List suppressionStateDescriptor = $convert.base64Decode(
    'ChBTdXBwcmVzc2lvblN0YXRlEiIKHlNVUFBSRVNTSU9OX1NUQVRFX1VOU1VQUFJFU1NFRBAAEi'
    'AKHFNVUFBSRVNTSU9OX1NUQVRFX1NVUFBSRVNTRUQQAQ==');

@$core.Deprecated('Use suppressionTypeDescriptor instead')
const SuppressionType$json = {
  '1': 'SuppressionType',
  '2': [
    {'1': 'SUPPRESSION_TYPE_INFINITE', '2': 0},
    {'1': 'SUPPRESSION_TYPE_LIMITED', '2': 1},
  ],
};

/// Descriptor for `SuppressionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List suppressionTypeDescriptor = $convert.base64Decode(
    'Cg9TdXBwcmVzc2lvblR5cGUSHQoZU1VQUFJFU1NJT05fVFlQRV9JTkZJTklURRAAEhwKGFNVUF'
    'BSRVNTSU9OX1RZUEVfTElNSVRFRBAB');

@$core.Deprecated('Use suppressionUnitDescriptor instead')
const SuppressionUnit$json = {
  '1': 'SuppressionUnit',
  '2': [
    {'1': 'SUPPRESSION_UNIT_MINUTES', '2': 0},
    {'1': 'SUPPRESSION_UNIT_SECONDS', '2': 1},
    {'1': 'SUPPRESSION_UNIT_HOURS', '2': 2},
  ],
};

/// Descriptor for `SuppressionUnit`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List suppressionUnitDescriptor = $convert.base64Decode(
    'Cg9TdXBwcmVzc2lvblVuaXQSHAoYU1VQUFJFU1NJT05fVU5JVF9NSU5VVEVTEAASHAoYU1VQUF'
    'JFU1NJT05fVU5JVF9TRUNPTkRTEAESGgoWU1VQUFJFU1NJT05fVU5JVF9IT1VSUxAC');

@$core.Deprecated('Use typeDescriptor instead')
const Type$json = {
  '1': 'Type',
  '2': [
    {'1': 'TYPE_INTEGER', '2': 0},
    {'1': 'TYPE_STRING', '2': 1},
    {'1': 'TYPE_BOOLEAN', '2': 2},
    {'1': 'TYPE_DOUBLE', '2': 3},
  ],
};

/// Descriptor for `Type`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List typeDescriptor = $convert.base64Decode(
    'CgRUeXBlEhAKDFRZUEVfSU5URUdFUhAAEg8KC1RZUEVfU1RSSU5HEAESEAoMVFlQRV9CT09MRU'
    'FOEAISDwoLVFlQRV9ET1VCTEUQAw==');

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

@$core.Deprecated('Use accountPolicyDescriptor instead')
const AccountPolicy$json = {
  '1': 'AccountPolicy',
  '2': [
    {'1': 'accountid', '3': 516110962, '4': 1, '5': 9, '10': 'accountid'},
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'policytype',
      '3': 319277736,
      '4': 1,
      '5': 14,
      '6': '.logs.PolicyType',
      '10': 'policytype'
    },
    {
      '1': 'scope',
      '3': 506131436,
      '4': 1,
      '5': 14,
      '6': '.logs.Scope',
      '10': 'scope'
    },
    {
      '1': 'selectioncriteria',
      '3': 145052429,
      '4': 1,
      '5': 9,
      '10': 'selectioncriteria'
    },
  ],
};

/// Descriptor for `AccountPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountPolicyDescriptor = $convert.base64Decode(
    'Cg1BY2NvdW50UG9saWN5EiAKCWFjY291bnRpZBjy9Iz2ASABKAlSCWFjY291bnRpZBIsCg9sYX'
    'N0dXBkYXRlZHRpbWUY9reVuQEgASgDUg9sYXN0dXBkYXRlZHRpbWUSKQoOcG9saWN5ZG9jdW1l'
    'bnQY6+n2VCABKAlSDnBvbGljeWRvY3VtZW50EiEKCnBvbGljeW5hbWUY3eLyNiABKAlSCnBvbG'
    'ljeW5hbWUSNAoKcG9saWN5dHlwZRiolZ+YASABKA4yEC5sb2dzLlBvbGljeVR5cGVSCnBvbGlj'
    'eXR5cGUSJQoFc2NvcGUY7Oer8QEgASgOMgsubG9ncy5TY29wZVIFc2NvcGUSLwoRc2VsZWN0aW'
    '9uY3JpdGVyaWEYjaaVRSABKAlSEXNlbGVjdGlvbmNyaXRlcmlh');

@$core.Deprecated('Use addKeyEntryDescriptor instead')
const AddKeyEntry$json = {
  '1': 'AddKeyEntry',
  '2': [
    {'1': 'key', '3': 135645293, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'overwriteifexists',
      '3': 230880030,
      '4': 1,
      '5': 8,
      '10': 'overwriteifexists'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `AddKeyEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addKeyEntryDescriptor = $convert.base64Decode(
    'CgtBZGRLZXlFbnRyeRITCgNrZXkY7ZDXQCABKAlSA2tleRIvChFvdmVyd3JpdGVpZmV4aXN0cx'
    'ie5otuIAEoCFIRb3ZlcndyaXRlaWZleGlzdHMSFwoFdmFsdWUYy6f7EiABKAlSBXZhbHVl');

@$core.Deprecated('Use addKeysDescriptor instead')
const AddKeys$json = {
  '1': 'AddKeys',
  '2': [
    {
      '1': 'entries',
      '3': 257458932,
      '4': 3,
      '5': 11,
      '6': '.logs.AddKeyEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `AddKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addKeysDescriptor = $convert.base64Decode(
    'CgdBZGRLZXlzEi4KB2VudHJpZXMY9IXieiADKAsyES5sb2dzLkFkZEtleUVudHJ5UgdlbnRyaW'
    'Vz');

@$core.Deprecated('Use aggregateLogGroupSummaryDescriptor instead')
const AggregateLogGroupSummary$json = {
  '1': 'AggregateLogGroupSummary',
  '2': [
    {
      '1': 'groupingidentifiers',
      '3': 363642629,
      '4': 3,
      '5': 11,
      '6': '.logs.GroupingIdentifier',
      '10': 'groupingidentifiers'
    },
    {
      '1': 'loggroupcount',
      '3': 176522920,
      '4': 1,
      '5': 5,
      '10': 'loggroupcount'
    },
  ],
};

/// Descriptor for `AggregateLogGroupSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aggregateLogGroupSummaryDescriptor = $convert.base64Decode(
    'ChhBZ2dyZWdhdGVMb2dHcm91cFN1bW1hcnkSTgoTZ3JvdXBpbmdpZGVudGlmaWVycxiF/rKtAS'
    'ADKAsyGC5sb2dzLkdyb3VwaW5nSWRlbnRpZmllclITZ3JvdXBpbmdpZGVudGlmaWVycxInCg1s'
    'b2dncm91cGNvdW50GKiNllQgASgFUg1sb2dncm91cGNvdW50');

@$core.Deprecated('Use anomalyDescriptor instead')
const Anomaly$json = {
  '1': 'Anomaly',
  '2': [
    {'1': 'active', '3': 427137328, '4': 1, '5': 8, '10': 'active'},
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
    {'1': 'anomalyid', '3': 201902912, '4': 1, '5': 9, '10': 'anomalyid'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'firstseen', '3': 101793055, '4': 1, '5': 3, '10': 'firstseen'},
    {
      '1': 'histogram',
      '3': 297617058,
      '4': 3,
      '5': 11,
      '6': '.logs.Anomaly.HistogramEntry',
      '10': 'histogram'
    },
    {
      '1': 'ispatternlevelsuppression',
      '3': 247622013,
      '4': 1,
      '5': 8,
      '10': 'ispatternlevelsuppression'
    },
    {'1': 'lastseen', '3': 414117031, '4': 1, '5': 3, '10': 'lastseen'},
    {
      '1': 'loggrouparnlist',
      '3': 374867736,
      '4': 3,
      '5': 9,
      '10': 'loggrouparnlist'
    },
    {
      '1': 'logsamples',
      '3': 85300633,
      '4': 3,
      '5': 11,
      '6': '.logs.LogEvent',
      '10': 'logsamples'
    },
    {'1': 'patternid', '3': 292391181, '4': 1, '5': 9, '10': 'patternid'},
    {'1': 'patternregex', '3': 249573581, '4': 1, '5': 9, '10': 'patternregex'},
    {
      '1': 'patternstring',
      '3': 484736421,
      '4': 1,
      '5': 9,
      '10': 'patternstring'
    },
    {
      '1': 'patterntokens',
      '3': 434268088,
      '4': 3,
      '5': 11,
      '6': '.logs.PatternToken',
      '10': 'patterntokens'
    },
    {'1': 'priority', '3': 350544650, '4': 1, '5': 9, '10': 'priority'},
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.State',
      '10': 'state'
    },
    {'1': 'suppressed', '3': 141710296, '4': 1, '5': 8, '10': 'suppressed'},
    {
      '1': 'suppresseddate',
      '3': 533956404,
      '4': 1,
      '5': 3,
      '10': 'suppresseddate'
    },
    {
      '1': 'suppresseduntil',
      '3': 437063138,
      '4': 1,
      '5': 3,
      '10': 'suppresseduntil'
    },
  ],
  '3': [Anomaly_HistogramEntry$json],
};

@$core.Deprecated('Use anomalyDescriptor instead')
const Anomaly_HistogramEntry$json = {
  '1': 'HistogramEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 3, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Anomaly`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anomalyDescriptor = $convert.base64Decode(
    'CgdBbm9tYWx5EhoKBmFjdGl2ZRiwstbLASABKAhSBmFjdGl2ZRIyChJhbm9tYWx5ZGV0ZWN0b3'
    'Jhcm4YrM/z1AEgASgJUhJhbm9tYWx5ZGV0ZWN0b3Jhcm4SHwoJYW5vbWFseWlkGMCWo2AgASgJ'
    'Uglhbm9tYWx5aWQSJAoLZGVzY3JpcHRpb24Y6va8owEgASgJUgtkZXNjcmlwdGlvbhIfCglmaX'
    'JzdHNlZW4Yn/rEMCABKANSCWZpcnN0c2VlbhI+CgloaXN0b2dyYW0Yoo31jQEgAygLMhwubG9n'
    'cy5Bbm9tYWx5Lkhpc3RvZ3JhbUVudHJ5UgloaXN0b2dyYW0SPwoZaXNwYXR0ZXJubGV2ZWxzdX'
    'BwcmVzc2lvbhj90ol2IAEoCFIZaXNwYXR0ZXJubGV2ZWxzdXBwcmVzc2lvbhIeCghsYXN0c2Vl'
    'bhin2bvFASABKANSCGxhc3RzZWVuEiwKD2xvZ2dyb3VwYXJubGlzdBiYjuCyASADKAlSD2xvZ2'
    'dyb3VwYXJubGlzdBIxCgpsb2dzYW1wbGVzGJmr1iggAygLMg4ubG9ncy5Mb2dFdmVudFIKbG9n'
    'c2FtcGxlcxIgCglwYXR0ZXJuaWQYjZK2iwEgASgJUglwYXR0ZXJuaWQSJQoMcGF0dGVybnJlZ2'
    'V4GM3hgHcgASgJUgxwYXR0ZXJucmVnZXgSKAoNcGF0dGVybnN0cmluZxil+5HnASABKAlSDXBh'
    'dHRlcm5zdHJpbmcSPAoNcGF0dGVybnRva2Vucxi4z4nPASADKAsyEi5sb2dzLlBhdHRlcm5Ub2'
    'tlblINcGF0dGVybnRva2VucxIeCghwcmlvcml0eRiKxpOnASABKAlSCHByaW9yaXR5EiUKBXN0'
    'YXRlGPflxMEBIAEoDjILLmxvZ3MuU3RhdGVSBXN0YXRlEiEKCnN1cHByZXNzZWQY2KfJQyABKA'
    'hSCnN1cHByZXNzZWQSKgoOc3VwcHJlc3NlZGRhdGUYtI7O/gEgASgDUg5zdXBwcmVzc2VkZGF0'
    'ZRIsCg9zdXBwcmVzc2VkdW50aWwY4pu00AEgASgDUg9zdXBwcmVzc2VkdW50aWwaPAoOSGlzdG'
    '9ncmFtRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKANSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use anomalyDetectorDescriptor instead')
const AnomalyDetector$json = {
  '1': 'AnomalyDetector',
  '2': [
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
    {
      '1': 'anomalydetectorstatus',
      '3': 458778431,
      '4': 1,
      '5': 14,
      '6': '.logs.AnomalyDetectorStatus',
      '10': 'anomalydetectorstatus'
    },
    {
      '1': 'anomalyvisibilitytime',
      '3': 287987260,
      '4': 1,
      '5': 3,
      '10': 'anomalyvisibilitytime'
    },
    {
      '1': 'creationtimestamp',
      '3': 206588645,
      '4': 1,
      '5': 3,
      '10': 'creationtimestamp'
    },
    {'1': 'detectorname', '3': 114651981, '4': 1, '5': 9, '10': 'detectorname'},
    {
      '1': 'evaluationfrequency',
      '3': 533700360,
      '4': 1,
      '5': 14,
      '6': '.logs.EvaluationFrequency',
      '10': 'evaluationfrequency'
    },
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'kmskeyid', '3': 510698477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'lastmodifiedtimestamp',
      '3': 40019279,
      '4': 1,
      '5': 3,
      '10': 'lastmodifiedtimestamp'
    },
    {
      '1': 'loggrouparnlist',
      '3': 374867736,
      '4': 3,
      '5': 9,
      '10': 'loggrouparnlist'
    },
  ],
};

/// Descriptor for `AnomalyDetector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anomalyDetectorDescriptor = $convert.base64Decode(
    'Cg9Bbm9tYWx5RGV0ZWN0b3ISMgoSYW5vbWFseWRldGVjdG9yYXJuGKzP89QBIAEoCVISYW5vbW'
    'FseWRldGVjdG9yYXJuElUKFWFub21hbHlkZXRlY3RvcnN0YXR1cxi/zuHaASABKA4yGy5sb2dz'
    'LkFub21hbHlEZXRlY3RvclN0YXR1c1IVYW5vbWFseWRldGVjdG9yc3RhdHVzEjgKFWFub21hbH'
    'l2aXNpYmlsaXR5dGltZRi8rKmJASABKANSFWFub21hbHl2aXNpYmlsaXR5dGltZRIvChFjcmVh'
    'dGlvbnRpbWVzdGFtcBjllcFiIAEoA1IRY3JlYXRpb250aW1lc3RhbXASJQoMZGV0ZWN0b3JuYW'
    '1lGM3m1TYgASgJUgxkZXRlY3Rvcm5hbWUSTwoTZXZhbHVhdGlvbmZyZXF1ZW5jeRiIvr7+ASAB'
    'KA4yGS5sb2dzLkV2YWx1YXRpb25GcmVxdWVuY3lSE2V2YWx1YXRpb25mcmVxdWVuY3kSJwoNZm'
    'lsdGVycGF0dGVybhiYh4pFIAEoCVINZmlsdGVycGF0dGVybhIeCghrbXNrZXlpZBjtx8LzASAB'
    'KAlSCGttc2tleWlkEjcKFWxhc3Rtb2RpZmllZHRpbWVzdGFtcBjPyooTIAEoA1IVbGFzdG1vZG'
    'lmaWVkdGltZXN0YW1wEiwKD2xvZ2dyb3VwYXJubGlzdBiYjuCyASADKAlSD2xvZ2dyb3VwYXJu'
    'bGlzdA==');

@$core.Deprecated('Use associateKmsKeyRequestDescriptor instead')
const AssociateKmsKeyRequest$json = {
  '1': 'AssociateKmsKeyRequest',
  '2': [
    {'1': 'kmskeyid', '3': 510698477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'resourceidentifier',
      '3': 309427407,
      '4': 1,
      '5': 9,
      '10': 'resourceidentifier'
    },
  ],
};

/// Descriptor for `AssociateKmsKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateKmsKeyRequestDescriptor = $convert.base64Decode(
    'ChZBc3NvY2lhdGVLbXNLZXlSZXF1ZXN0Eh4KCGttc2tleWlkGO3HwvMBIAEoCVIIa21za2V5aW'
    'QSJQoMbG9nZ3JvdXBuYW1lGMy/1C0gASgJUgxsb2dncm91cG5hbWUSMgoScmVzb3VyY2VpZGVu'
    'dGlmaWVyGM/5xZMBIAEoCVIScmVzb3VyY2VpZGVudGlmaWVy');

@$core.Deprecated(
    'Use associateSourceToS3TableIntegrationRequestDescriptor instead')
const AssociateSourceToS3TableIntegrationRequest$json = {
  '1': 'AssociateSourceToS3TableIntegrationRequest',
  '2': [
    {
      '1': 'datasource',
      '3': 345762713,
      '4': 1,
      '5': 11,
      '6': '.logs.DataSource',
      '10': 'datasource'
    },
    {
      '1': 'integrationarn',
      '3': 432021733,
      '4': 1,
      '5': 9,
      '10': 'integrationarn'
    },
  ],
};

/// Descriptor for `AssociateSourceToS3TableIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    associateSourceToS3TableIntegrationRequestDescriptor =
    $convert.base64Decode(
        'CipBc3NvY2lhdGVTb3VyY2VUb1MzVGFibGVJbnRlZ3JhdGlvblJlcXVlc3QSNAoKZGF0YXNvdX'
        'JjZRiZ1++kASABKAsyEC5sb2dzLkRhdGFTb3VyY2VSCmRhdGFzb3VyY2USKgoOaW50ZWdyYXRp'
        'b25hcm4Y5cGAzgEgASgJUg5pbnRlZ3JhdGlvbmFybg==');

@$core.Deprecated(
    'Use associateSourceToS3TableIntegrationResponseDescriptor instead')
const AssociateSourceToS3TableIntegrationResponse$json = {
  '1': 'AssociateSourceToS3TableIntegrationResponse',
  '2': [
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `AssociateSourceToS3TableIntegrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    associateSourceToS3TableIntegrationResponseDescriptor =
    $convert.base64Decode(
        'CitBc3NvY2lhdGVTb3VyY2VUb1MzVGFibGVJbnRlZ3JhdGlvblJlc3BvbnNlEiEKCmlkZW50aW'
        'ZpZXIYv9CWRSABKAlSCmlkZW50aWZpZXI=');

@$core.Deprecated('Use cSVDescriptor instead')
const CSV$json = {
  '1': 'CSV',
  '2': [
    {'1': 'columns', '3': 179509821, '4': 3, '5': 9, '10': 'columns'},
    {'1': 'delimiter', '3': 202749435, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'quotecharacter',
      '3': 92052223,
      '4': 1,
      '5': 9,
      '10': 'quotecharacter'
    },
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `CSV`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cSVDescriptor = $convert.base64Decode(
    'CgNDU1YSGwoHY29sdW1ucxi9tMxVIAMoCVIHY29sdW1ucxIfCglkZWxpbWl0ZXIY++vWYCABKA'
    'lSCWRlbGltaXRlchIpCg5xdW90ZWNoYXJhY3Rlchj/tfIrIAEoCVIOcXVvdGVjaGFyYWN0ZXIS'
    'GgoGc291cmNlGNnTvN4BIAEoCVIGc291cmNl');

@$core.Deprecated('Use cancelExportTaskRequestDescriptor instead')
const CancelExportTaskRequest$json = {
  '1': 'CancelExportTaskRequest',
  '2': [
    {'1': 'taskid', '3': 216769858, '4': 1, '5': 9, '10': 'taskid'},
  ],
};

/// Descriptor for `CancelExportTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelExportTaskRequestDescriptor =
    $convert.base64Decode(
        'ChdDYW5jZWxFeHBvcnRUYXNrUmVxdWVzdBIZCgZ0YXNraWQYwsquZyABKAlSBnRhc2tpZA==');

@$core.Deprecated('Use cancelImportTaskRequestDescriptor instead')
const CancelImportTaskRequest$json = {
  '1': 'CancelImportTaskRequest',
  '2': [
    {'1': 'importid', '3': 513429114, '4': 1, '5': 9, '10': 'importid'},
  ],
};

/// Descriptor for `CancelImportTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelImportTaskRequestDescriptor =
    $convert.base64Decode(
        'ChdDYW5jZWxJbXBvcnRUYXNrUmVxdWVzdBIeCghpbXBvcnRpZBj6nOn0ASABKAlSCGltcG9ydG'
        'lk');

@$core.Deprecated('Use cancelImportTaskResponseDescriptor instead')
const CancelImportTaskResponse$json = {
  '1': 'CancelImportTaskResponse',
  '2': [
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {'1': 'importid', '3': 513429114, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importstatistics',
      '3': 60366280,
      '4': 1,
      '5': 11,
      '6': '.logs.ImportStatistics',
      '10': 'importstatistics'
    },
    {
      '1': 'importstatus',
      '3': 31427999,
      '4': 1,
      '5': 14,
      '6': '.logs.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
  ],
};

/// Descriptor for `CancelImportTaskResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelImportTaskResponseDescriptor = $convert.base64Decode(
    'ChhDYW5jZWxJbXBvcnRUYXNrUmVzcG9uc2USJQoMY3JlYXRpb250aW1lGIbFxBkgASgDUgxjcm'
    'VhdGlvbnRpbWUSHgoIaW1wb3J0aWQY+pzp9AEgASgJUghpbXBvcnRpZBJFChBpbXBvcnRzdGF0'
    'aXN0aWNzGMi75BwgASgLMhYubG9ncy5JbXBvcnRTdGF0aXN0aWNzUhBpbXBvcnRzdGF0aXN0aW'
    'NzEjkKDGltcG9ydHN0YXR1cxifm/4OIAEoDjISLmxvZ3MuSW1wb3J0U3RhdHVzUgxpbXBvcnRz'
    'dGF0dXMSLAoPbGFzdHVwZGF0ZWR0aW1lGPa3lbkBIAEoA1IPbGFzdHVwZGF0ZWR0aW1l');

@$core.Deprecated('Use configurationTemplateDescriptor instead')
const ConfigurationTemplate$json = {
  '1': 'ConfigurationTemplate',
  '2': [
    {
      '1': 'allowedactionforallowvendedlogsdeliveryforresource',
      '3': 179630824,
      '4': 1,
      '5': 9,
      '10': 'allowedactionforallowvendedlogsdeliveryforresource'
    },
    {
      '1': 'allowedfielddelimiters',
      '3': 266366692,
      '4': 3,
      '5': 9,
      '10': 'allowedfielddelimiters'
    },
    {
      '1': 'allowedfields',
      '3': 425965483,
      '4': 3,
      '5': 11,
      '6': '.logs.RecordField',
      '10': 'allowedfields'
    },
    {
      '1': 'allowedoutputformats',
      '3': 395034331,
      '4': 3,
      '5': 14,
      '6': '.logs.OutputFormat',
      '10': 'allowedoutputformats'
    },
    {
      '1': 'allowedsuffixpathfields',
      '3': 529615989,
      '4': 3,
      '5': 9,
      '10': 'allowedsuffixpathfields'
    },
    {
      '1': 'defaultdeliveryconfigvalues',
      '3': 235569423,
      '4': 1,
      '5': 11,
      '6': '.logs.ConfigurationTemplateDeliveryConfigValues',
      '10': 'defaultdeliveryconfigvalues'
    },
    {
      '1': 'deliverydestinationtype',
      '3': 33253176,
      '4': 1,
      '5': 14,
      '6': '.logs.DeliveryDestinationType',
      '10': 'deliverydestinationtype'
    },
    {'1': 'logtype', '3': 257838938, '4': 1, '5': 9, '10': 'logtype'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
    {'1': 'service', '3': 383770213, '4': 1, '5': 9, '10': 'service'},
  ],
};

/// Descriptor for `ConfigurationTemplate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List configurationTemplateDescriptor = $convert.base64Decode(
    'ChVDb25maWd1cmF0aW9uVGVtcGxhdGUScQoyYWxsb3dlZGFjdGlvbmZvcmFsbG93dmVuZGVkbG'
    '9nc2RlbGl2ZXJ5Zm9ycmVzb3VyY2UY6OXTVSABKAlSMmFsbG93ZWRhY3Rpb25mb3JhbGxvd3Zl'
    'bmRlZGxvZ3NkZWxpdmVyeWZvcnJlc291cmNlEjkKFmFsbG93ZWRmaWVsZGRlbGltaXRlcnMY5N'
    '2BfyADKAlSFmFsbG93ZWRmaWVsZGRlbGltaXRlcnMSOwoNYWxsb3dlZGZpZWxkcxir747LASAD'
    'KAsyES5sb2dzLlJlY29yZEZpZWxkUg1hbGxvd2VkZmllbGRzEkoKFGFsbG93ZWRvdXRwdXRmb3'
    'JtYXRzGNv9rrwBIAMoDjISLmxvZ3MuT3V0cHV0Rm9ybWF0UhRhbGxvd2Vkb3V0cHV0Zm9ybWF0'
    'cxI8ChdhbGxvd2Vkc3VmZml4cGF0aGZpZWxkcxj1mMX8ASADKAlSF2FsbG93ZWRzdWZmaXhwYX'
    'RoZmllbGRzEnQKG2RlZmF1bHRkZWxpdmVyeWNvbmZpZ3ZhbHVlcxiPgqpwIAEoCzIvLmxvZ3Mu'
    'Q29uZmlndXJhdGlvblRlbXBsYXRlRGVsaXZlcnlDb25maWdWYWx1ZXNSG2RlZmF1bHRkZWxpdm'
    'VyeWNvbmZpZ3ZhbHVlcxJaChdkZWxpdmVyeWRlc3RpbmF0aW9udHlwZRi4zu0PIAEoDjIdLmxv'
    'Z3MuRGVsaXZlcnlEZXN0aW5hdGlvblR5cGVSF2RlbGl2ZXJ5ZGVzdGluYXRpb250eXBlEhsKB2'
    'xvZ3R5cGUY2p75eiABKAlSB2xvZ3R5cGUSJQoMcmVzb3VyY2V0eXBlGP6V0AMgASgJUgxyZXNv'
    'dXJjZXR5cGUSHAoHc2VydmljZRjlvP+2ASABKAlSB3NlcnZpY2U=');

@$core.Deprecated(
    'Use configurationTemplateDeliveryConfigValuesDescriptor instead')
const ConfigurationTemplateDeliveryConfigValues$json = {
  '1': 'ConfigurationTemplateDeliveryConfigValues',
  '2': [
    {
      '1': 'fielddelimiter',
      '3': 433214845,
      '4': 1,
      '5': 9,
      '10': 'fielddelimiter'
    },
    {'1': 'recordfields', '3': 309625178, '4': 3, '5': 9, '10': 'recordfields'},
    {
      '1': 's3deliveryconfiguration',
      '3': 95221108,
      '4': 1,
      '5': 11,
      '6': '.logs.S3DeliveryConfiguration',
      '10': 's3deliveryconfiguration'
    },
  ],
};

/// Descriptor for `ConfigurationTemplateDeliveryConfigValues`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    configurationTemplateDeliveryConfigValuesDescriptor = $convert.base64Decode(
        'CilDb25maWd1cmF0aW9uVGVtcGxhdGVEZWxpdmVyeUNvbmZpZ1ZhbHVlcxIqCg5maWVsZGRlbG'
        'ltaXRlchj9qsnOASABKAlSDmZpZWxkZGVsaW1pdGVyEiYKDHJlY29yZGZpZWxkcxjagtKTASAD'
        'KAlSDHJlY29yZGZpZWxkcxJaChdzM2RlbGl2ZXJ5Y29uZmlndXJhdGlvbhj06rMtIAEoCzIdLm'
        'xvZ3MuUzNEZWxpdmVyeUNvbmZpZ3VyYXRpb25SF3MzZGVsaXZlcnljb25maWd1cmF0aW9u');

@$core.Deprecated('Use conflictExceptionDescriptor instead')
const ConflictException$json = {
  '1': 'ConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictExceptionDescriptor = $convert.base64Decode(
    'ChFDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use copyValueDescriptor instead')
const CopyValue$json = {
  '1': 'CopyValue',
  '2': [
    {
      '1': 'entries',
      '3': 257458932,
      '4': 3,
      '5': 11,
      '6': '.logs.CopyValueEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `CopyValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyValueDescriptor = $convert.base64Decode(
    'CglDb3B5VmFsdWUSMQoHZW50cmllcxj0heJ6IAMoCzIULmxvZ3MuQ29weVZhbHVlRW50cnlSB2'
    'VudHJpZXM=');

@$core.Deprecated('Use copyValueEntryDescriptor instead')
const CopyValueEntry$json = {
  '1': 'CopyValueEntry',
  '2': [
    {
      '1': 'overwriteifexists',
      '3': 230880030,
      '4': 1,
      '5': 8,
      '10': 'overwriteifexists'
    },
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
    {'1': 'target', '3': 308316233, '4': 1, '5': 9, '10': 'target'},
  ],
};

/// Descriptor for `CopyValueEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyValueEntryDescriptor = $convert.base64Decode(
    'Cg5Db3B5VmFsdWVFbnRyeRIvChFvdmVyd3JpdGVpZmV4aXN0cxie5otuIAEoCFIRb3ZlcndyaX'
    'RlaWZleGlzdHMSGgoGc291cmNlGNnTvN4BIAEoCVIGc291cmNlEhoKBnRhcmdldBjJkIKTASAB'
    'KAlSBnRhcmdldA==');

@$core.Deprecated('Use createDeliveryRequestDescriptor instead')
const CreateDeliveryRequest$json = {
  '1': 'CreateDeliveryRequest',
  '2': [
    {
      '1': 'deliverydestinationarn',
      '3': 156800339,
      '4': 1,
      '5': 9,
      '10': 'deliverydestinationarn'
    },
    {
      '1': 'deliverysourcename',
      '3': 277065836,
      '4': 1,
      '5': 9,
      '10': 'deliverysourcename'
    },
    {
      '1': 'fielddelimiter',
      '3': 433214845,
      '4': 1,
      '5': 9,
      '10': 'fielddelimiter'
    },
    {'1': 'recordfields', '3': 309625178, '4': 3, '5': 9, '10': 'recordfields'},
    {
      '1': 's3deliveryconfiguration',
      '3': 95221108,
      '4': 1,
      '5': 11,
      '6': '.logs.S3DeliveryConfiguration',
      '10': 's3deliveryconfiguration'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.CreateDeliveryRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [CreateDeliveryRequest_TagsEntry$json],
};

@$core.Deprecated('Use createDeliveryRequestDescriptor instead')
const CreateDeliveryRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateDeliveryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDeliveryRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVEZWxpdmVyeVJlcXVlc3QSOQoWZGVsaXZlcnlkZXN0aW5hdGlvbmFybhjTquJKIA'
    'EoCVIWZGVsaXZlcnlkZXN0aW5hdGlvbmFybhIyChJkZWxpdmVyeXNvdXJjZW5hbWUY7OCOhAEg'
    'ASgJUhJkZWxpdmVyeXNvdXJjZW5hbWUSKgoOZmllbGRkZWxpbWl0ZXIY/arJzgEgASgJUg5maW'
    'VsZGRlbGltaXRlchImCgxyZWNvcmRmaWVsZHMY2oLSkwEgAygJUgxyZWNvcmRmaWVsZHMSWgoX'
    'czNkZWxpdmVyeWNvbmZpZ3VyYXRpb24Y9OqzLSABKAsyHS5sb2dzLlMzRGVsaXZlcnlDb25maW'
    'd1cmF0aW9uUhdzM2RlbGl2ZXJ5Y29uZmlndXJhdGlvbhI9CgR0YWdzGKHX26ABIAMoCzIlLmxv'
    'Z3MuQ3JlYXRlRGVsaXZlcnlSZXF1ZXN0LlRhZ3NFbnRyeVIEdGFncxo3CglUYWdzRW50cnkSEA'
    'oDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createDeliveryResponseDescriptor instead')
const CreateDeliveryResponse$json = {
  '1': 'CreateDeliveryResponse',
  '2': [
    {
      '1': 'delivery',
      '3': 432355926,
      '4': 1,
      '5': 11,
      '6': '.logs.Delivery',
      '10': 'delivery'
    },
  ],
};

/// Descriptor for `CreateDeliveryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDeliveryResponseDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVEZWxpdmVyeVJlc3BvbnNlEi4KCGRlbGl2ZXJ5GNb0lM4BIAEoCzIOLmxvZ3MuRG'
        'VsaXZlcnlSCGRlbGl2ZXJ5');

@$core.Deprecated('Use createExportTaskRequestDescriptor instead')
const CreateExportTaskRequest$json = {
  '1': 'CreateExportTaskRequest',
  '2': [
    {'1': 'destination', '3': 316564672, '4': 1, '5': 9, '10': 'destination'},
    {
      '1': 'destinationprefix',
      '3': 172629044,
      '4': 1,
      '5': 9,
      '10': 'destinationprefix'
    },
    {'1': 'from', '3': 365789302, '4': 1, '5': 3, '10': 'from'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'logstreamnameprefix',
      '3': 24068871,
      '4': 1,
      '5': 9,
      '10': 'logstreamnameprefix'
    },
    {'1': 'taskname', '3': 82438536, '4': 1, '5': 9, '10': 'taskname'},
    {'1': 'to', '3': 38094885, '4': 1, '5': 3, '10': 'to'},
  ],
};

/// Descriptor for `CreateExportTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createExportTaskRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVFeHBvcnRUYXNrUmVxdWVzdBIkCgtkZXN0aW5hdGlvbhjAyfmWASABKAlSC2Rlc3'
    'RpbmF0aW9uEi8KEWRlc3RpbmF0aW9ucHJlZml4GLS4qFIgASgJUhFkZXN0aW5hdGlvbnByZWZp'
    'eBIWCgRmcm9tGPaAtq4BIAEoA1IEZnJvbRIlCgxsb2dncm91cG5hbWUYzL/ULSABKAlSDGxvZ2'
    'dyb3VwbmFtZRIzChNsb2dzdHJlYW1uYW1lcHJlZml4GIeGvQsgASgJUhNsb2dzdHJlYW1uYW1l'
    'cHJlZml4Eh0KCHRhc2tuYW1lGIjTpycgASgJUgh0YXNrbmFtZRIRCgJ0bxilkJUSIAEoA1ICdG'
    '8=');

@$core.Deprecated('Use createExportTaskResponseDescriptor instead')
const CreateExportTaskResponse$json = {
  '1': 'CreateExportTaskResponse',
  '2': [
    {'1': 'taskid', '3': 216769858, '4': 1, '5': 9, '10': 'taskid'},
  ],
};

/// Descriptor for `CreateExportTaskResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createExportTaskResponseDescriptor =
    $convert.base64Decode(
        'ChhDcmVhdGVFeHBvcnRUYXNrUmVzcG9uc2USGQoGdGFza2lkGMLKrmcgASgJUgZ0YXNraWQ=');

@$core.Deprecated('Use createImportTaskRequestDescriptor instead')
const CreateImportTaskRequest$json = {
  '1': 'CreateImportTaskRequest',
  '2': [
    {
      '1': 'importfilter',
      '3': 170561015,
      '4': 1,
      '5': 11,
      '6': '.logs.ImportFilter',
      '10': 'importfilter'
    },
    {
      '1': 'importrolearn',
      '3': 453628468,
      '4': 1,
      '5': 9,
      '10': 'importrolearn'
    },
    {
      '1': 'importsourcearn',
      '3': 161570329,
      '4': 1,
      '5': 9,
      '10': 'importsourcearn'
    },
  ],
};

/// Descriptor for `CreateImportTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createImportTaskRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVJbXBvcnRUYXNrUmVxdWVzdBI5CgxpbXBvcnRmaWx0ZXIY95uqUSABKAsyEi5sb2'
    'dzLkltcG9ydEZpbHRlclIMaW1wb3J0ZmlsdGVyEigKDWltcG9ydHJvbGVhcm4YtKSn2AEgASgJ'
    'Ug1pbXBvcnRyb2xlYXJuEisKD2ltcG9ydHNvdXJjZWFybhiZvIVNIAEoCVIPaW1wb3J0c291cm'
    'NlYXJu');

@$core.Deprecated('Use createImportTaskResponseDescriptor instead')
const CreateImportTaskResponse$json = {
  '1': 'CreateImportTaskResponse',
  '2': [
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'importdestinationarn',
      '3': 94477180,
      '4': 1,
      '5': 9,
      '10': 'importdestinationarn'
    },
    {'1': 'importid', '3': 513429114, '4': 1, '5': 9, '10': 'importid'},
  ],
};

/// Descriptor for `CreateImportTaskResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createImportTaskResponseDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVJbXBvcnRUYXNrUmVzcG9uc2USJQoMY3JlYXRpb250aW1lGIbFxBkgASgDUgxjcm'
    'VhdGlvbnRpbWUSNQoUaW1wb3J0ZGVzdGluYXRpb25hcm4Y/LaGLSABKAlSFGltcG9ydGRlc3Rp'
    'bmF0aW9uYXJuEh4KCGltcG9ydGlkGPqc6fQBIAEoCVIIaW1wb3J0aWQ=');

@$core.Deprecated('Use createLogAnomalyDetectorRequestDescriptor instead')
const CreateLogAnomalyDetectorRequest$json = {
  '1': 'CreateLogAnomalyDetectorRequest',
  '2': [
    {
      '1': 'anomalyvisibilitytime',
      '3': 287987260,
      '4': 1,
      '5': 3,
      '10': 'anomalyvisibilitytime'
    },
    {'1': 'detectorname', '3': 114651981, '4': 1, '5': 9, '10': 'detectorname'},
    {
      '1': 'evaluationfrequency',
      '3': 533700360,
      '4': 1,
      '5': 14,
      '6': '.logs.EvaluationFrequency',
      '10': 'evaluationfrequency'
    },
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'kmskeyid', '3': 510698477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'loggrouparnlist',
      '3': 374867736,
      '4': 3,
      '5': 9,
      '10': 'loggrouparnlist'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.CreateLogAnomalyDetectorRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [CreateLogAnomalyDetectorRequest_TagsEntry$json],
};

@$core.Deprecated('Use createLogAnomalyDetectorRequestDescriptor instead')
const CreateLogAnomalyDetectorRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateLogAnomalyDetectorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createLogAnomalyDetectorRequestDescriptor = $convert.base64Decode(
    'Ch9DcmVhdGVMb2dBbm9tYWx5RGV0ZWN0b3JSZXF1ZXN0EjgKFWFub21hbHl2aXNpYmlsaXR5dG'
    'ltZRi8rKmJASABKANSFWFub21hbHl2aXNpYmlsaXR5dGltZRIlCgxkZXRlY3Rvcm5hbWUYzebV'
    'NiABKAlSDGRldGVjdG9ybmFtZRJPChNldmFsdWF0aW9uZnJlcXVlbmN5GIi+vv4BIAEoDjIZLm'
    'xvZ3MuRXZhbHVhdGlvbkZyZXF1ZW5jeVITZXZhbHVhdGlvbmZyZXF1ZW5jeRInCg1maWx0ZXJw'
    'YXR0ZXJuGJiHikUgASgJUg1maWx0ZXJwYXR0ZXJuEh4KCGttc2tleWlkGO3HwvMBIAEoCVIIa2'
    '1za2V5aWQSLAoPbG9nZ3JvdXBhcm5saXN0GJiO4LIBIAMoCVIPbG9nZ3JvdXBhcm5saXN0EkcK'
    'BHRhZ3MYodfboAEgAygLMi8ubG9ncy5DcmVhdGVMb2dBbm9tYWx5RGV0ZWN0b3JSZXF1ZXN0Ll'
    'RhZ3NFbnRyeVIEdGFncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUY'
    'AiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createLogAnomalyDetectorResponseDescriptor instead')
const CreateLogAnomalyDetectorResponse$json = {
  '1': 'CreateLogAnomalyDetectorResponse',
  '2': [
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
  ],
};

/// Descriptor for `CreateLogAnomalyDetectorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createLogAnomalyDetectorResponseDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVMb2dBbm9tYWx5RGV0ZWN0b3JSZXNwb25zZRIyChJhbm9tYWx5ZGV0ZWN0b3Jhcm'
        '4YrM/z1AEgASgJUhJhbm9tYWx5ZGV0ZWN0b3Jhcm4=');

@$core.Deprecated('Use createLogGroupRequestDescriptor instead')
const CreateLogGroupRequest$json = {
  '1': 'CreateLogGroupRequest',
  '2': [
    {
      '1': 'deletionprotectionenabled',
      '3': 475522738,
      '4': 1,
      '5': 8,
      '10': 'deletionprotectionenabled'
    },
    {'1': 'kmskeyid', '3': 510698477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'loggroupclass',
      '3': 518605953,
      '4': 1,
      '5': 14,
      '6': '.logs.LogGroupClass',
      '10': 'loggroupclass'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.CreateLogGroupRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [CreateLogGroupRequest_TagsEntry$json],
};

@$core.Deprecated('Use createLogGroupRequestDescriptor instead')
const CreateLogGroupRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateLogGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createLogGroupRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVMb2dHcm91cFJlcXVlc3QSQAoZZGVsZXRpb25wcm90ZWN0aW9uZW5hYmxlZBiyzd'
    '/iASABKAhSGWRlbGV0aW9ucHJvdGVjdGlvbmVuYWJsZWQSHgoIa21za2V5aWQY7cfC8wEgASgJ'
    'UghrbXNrZXlpZBI9Cg1sb2dncm91cGNsYXNzGIGZpfcBIAEoDjITLmxvZ3MuTG9nR3JvdXBDbG'
    'Fzc1INbG9nZ3JvdXBjbGFzcxIlCgxsb2dncm91cG5hbWUYzL/ULSABKAlSDGxvZ2dyb3VwbmFt'
    'ZRI9CgR0YWdzGKHX26ABIAMoCzIlLmxvZ3MuQ3JlYXRlTG9nR3JvdXBSZXF1ZXN0LlRhZ3NFbn'
    'RyeVIEdGFncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createLogStreamRequestDescriptor instead')
const CreateLogStreamRequest$json = {
  '1': 'CreateLogStreamRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
  ],
};

/// Descriptor for `CreateLogStreamRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createLogStreamRequestDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVMb2dTdHJlYW1SZXF1ZXN0EiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3'
        'JvdXBuYW1lEigKDWxvZ3N0cmVhbW5hbWUYo/fu0AEgASgJUg1sb2dzdHJlYW1uYW1l');

@$core.Deprecated('Use createScheduledQueryRequestDescriptor instead')
const CreateScheduledQueryRequest$json = {
  '1': 'CreateScheduledQueryRequest',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'destinationconfiguration',
      '3': 145518016,
      '4': 1,
      '5': 11,
      '6': '.logs.DestinationConfiguration',
      '10': 'destinationconfiguration'
    },
    {
      '1': 'executionrolearn',
      '3': 230613553,
      '4': 1,
      '5': 9,
      '10': 'executionrolearn'
    },
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'scheduleendtime',
      '3': 111645113,
      '4': 1,
      '5': 3,
      '10': 'scheduleendtime'
    },
    {
      '1': 'scheduleexpression',
      '3': 287794975,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'schedulestarttime',
      '3': 464194170,
      '4': 1,
      '5': 3,
      '10': 'schedulestarttime'
    },
    {
      '1': 'starttimeoffset',
      '3': 308525274,
      '4': 1,
      '5': 3,
      '10': 'starttimeoffset'
    },
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryState',
      '10': 'state'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.CreateScheduledQueryRequest.TagsEntry',
      '10': 'tags'
    },
    {'1': 'timezone', '3': 190615331, '4': 1, '5': 9, '10': 'timezone'},
  ],
  '3': [CreateScheduledQueryRequest_TagsEntry$json],
};

@$core.Deprecated('Use createScheduledQueryRequestDescriptor instead')
const CreateScheduledQueryRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduledQueryRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVTY2hlZHVsZWRRdWVyeVJlcXVlc3QSJAoLZGVzY3JpcHRpb24Y6va8owEgASgJUg'
    'tkZXNjcmlwdGlvbhJdChhkZXN0aW5hdGlvbmNvbmZpZ3VyYXRpb24YwNuxRSABKAsyHi5sb2dz'
    'LkRlc3RpbmF0aW9uQ29uZmlndXJhdGlvblIYZGVzdGluYXRpb25jb25maWd1cmF0aW9uEi0KEG'
    'V4ZWN1dGlvbnJvbGVhcm4YscT7bSABKAlSEGV4ZWN1dGlvbnJvbGVhcm4SNAoTbG9nZ3JvdXBp'
    'ZGVudGlmaWVycxj52rjDASADKAlSE2xvZ2dyb3VwaWRlbnRpZmllcnMSFQoEbmFtZRjn++ZpIA'
    'EoCVIEbmFtZRI9Cg1xdWVyeWxhbmd1YWdlGPav96MBIAEoDjITLmxvZ3MuUXVlcnlMYW5ndWFn'
    'ZVINcXVlcnlsYW5ndWFnZRIkCgtxdWVyeXN0cmluZxiHgZ34ASABKAlSC3F1ZXJ5c3RyaW5nEi'
    'sKD3NjaGVkdWxlZW5kdGltZRi5o541IAEoA1IPc2NoZWR1bGVlbmR0aW1lEjIKEnNjaGVkdWxl'
    'ZXhwcmVzc2lvbhifzp2JASABKAlSEnNjaGVkdWxlZXhwcmVzc2lvbhIwChFzY2hlZHVsZXN0YX'
    'J0dGltZRj6lKzdASABKANSEXNjaGVkdWxlc3RhcnR0aW1lEiwKD3N0YXJ0dGltZW9mZnNldBja'
    '8Y6TASABKANSD3N0YXJ0dGltZW9mZnNldBIzCgVzdGF0ZRj35cTBASABKA4yGS5sb2dzLlNjaG'
    'VkdWxlZFF1ZXJ5U3RhdGVSBXN0YXRlEkMKBHRhZ3MYodfboAEgAygLMisubG9ncy5DcmVhdGVT'
    'Y2hlZHVsZWRRdWVyeVJlcXVlc3QuVGFnc0VudHJ5UgR0YWdzEh0KCHRpbWV6b25lGKOe8logAS'
    'gJUgh0aW1lem9uZRo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiAB'
    'KAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createScheduledQueryResponseDescriptor instead')
const CreateScheduledQueryResponse$json = {
  '1': 'CreateScheduledQueryResponse',
  '2': [
    {
      '1': 'scheduledqueryarn',
      '3': 240292916,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `CreateScheduledQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduledQueryResponseDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVTY2hlZHVsZWRRdWVyeVJlc3BvbnNlEi8KEXNjaGVkdWxlZHF1ZXJ5YXJuGLSoyn'
        'IgASgJUhFzY2hlZHVsZWRxdWVyeWFybhIzCgVzdGF0ZRj35cTBASABKA4yGS5sb2dzLlNjaGVk'
        'dWxlZFF1ZXJ5U3RhdGVSBXN0YXRl');

@$core.Deprecated('Use dataAlreadyAcceptedExceptionDescriptor instead')
const DataAlreadyAcceptedException$json = {
  '1': 'DataAlreadyAcceptedException',
  '2': [
    {
      '1': 'expectedsequencetoken',
      '3': 419198802,
      '4': 1,
      '5': 9,
      '10': 'expectedsequencetoken'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DataAlreadyAcceptedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataAlreadyAcceptedExceptionDescriptor =
    $convert.base64Decode(
        'ChxEYXRhQWxyZWFkeUFjY2VwdGVkRXhjZXB0aW9uEjgKFWV4cGVjdGVkc2VxdWVuY2V0b2tlbh'
        'jS7vHHASABKAlSFWV4cGVjdGVkc2VxdWVuY2V0b2tlbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdt'
        'ZXNzYWdl');

@$core.Deprecated('Use dataSourceDescriptor instead')
const DataSource$json = {
  '1': 'DataSource',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'type', '3': 287830350, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `DataSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataSourceDescriptor = $convert.base64Decode(
    'CgpEYXRhU291cmNlEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSFgoEdHlwZRjO4p+JASABKAlSBH'
    'R5cGU=');

@$core.Deprecated('Use dataSourceFilterDescriptor instead')
const DataSourceFilter$json = {
  '1': 'DataSourceFilter',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'type', '3': 287830350, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `DataSourceFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataSourceFilterDescriptor = $convert.base64Decode(
    'ChBEYXRhU291cmNlRmlsdGVyEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSFgoEdHlwZRjO4p+JAS'
    'ABKAlSBHR5cGU=');

@$core.Deprecated('Use dateTimeConverterDescriptor instead')
const DateTimeConverter$json = {
  '1': 'DateTimeConverter',
  '2': [
    {'1': 'locale', '3': 186372248, '4': 1, '5': 9, '10': 'locale'},
    {
      '1': 'matchpatterns',
      '3': 310725878,
      '4': 3,
      '5': 9,
      '10': 'matchpatterns'
    },
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
    {
      '1': 'sourcetimezone',
      '3': 154904028,
      '4': 1,
      '5': 9,
      '10': 'sourcetimezone'
    },
    {'1': 'target', '3': 308316233, '4': 1, '5': 9, '10': 'target'},
    {'1': 'targetformat', '3': 416820428, '4': 1, '5': 9, '10': 'targetformat'},
    {
      '1': 'targettimezone',
      '3': 310709132,
      '4': 1,
      '5': 9,
      '10': 'targettimezone'
    },
  ],
};

/// Descriptor for `DateTimeConverter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dateTimeConverterDescriptor = $convert.base64Decode(
    'ChFEYXRlVGltZUNvbnZlcnRlchIZCgZsb2NhbGUYmKHvWCABKAlSBmxvY2FsZRIoCg1tYXRjaH'
    'BhdHRlcm5zGPaZlZQBIAMoCVINbWF0Y2hwYXR0ZXJucxIaCgZzb3VyY2UY2dO83gEgASgJUgZz'
    'b3VyY2USKQoOc291cmNldGltZXpvbmUY3MvuSSABKAlSDnNvdXJjZXRpbWV6b25lEhoKBnRhcm'
    'dldBjJkIKTASABKAlSBnRhcmdldBImCgx0YXJnZXRmb3JtYXQYzNngxgEgASgJUgx0YXJnZXRm'
    'b3JtYXQSKgoOdGFyZ2V0dGltZXpvbmUYjJeUlAEgASgJUg50YXJnZXR0aW1lem9uZQ==');

@$core.Deprecated('Use deleteAccountPolicyRequestDescriptor instead')
const DeleteAccountPolicyRequest$json = {
  '1': 'DeleteAccountPolicyRequest',
  '2': [
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'policytype',
      '3': 319277736,
      '4': 1,
      '5': 14,
      '6': '.logs.PolicyType',
      '10': 'policytype'
    },
  ],
};

/// Descriptor for `DeleteAccountPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAccountPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVBY2NvdW50UG9saWN5UmVxdWVzdBIhCgpwb2xpY3luYW1lGN3i8jYgASgJUgpwb2'
        'xpY3luYW1lEjQKCnBvbGljeXR5cGUYqJWfmAEgASgOMhAubG9ncy5Qb2xpY3lUeXBlUgpwb2xp'
        'Y3l0eXBl');

@$core.Deprecated('Use deleteDataProtectionPolicyRequestDescriptor instead')
const DeleteDataProtectionPolicyRequest$json = {
  '1': 'DeleteDataProtectionPolicyRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
  ],
};

/// Descriptor for `DeleteDataProtectionPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDataProtectionPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiFEZWxldGVEYXRhUHJvdGVjdGlvblBvbGljeVJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaW'
        'VyGNzPpd8BIAEoCVISbG9nZ3JvdXBpZGVudGlmaWVy');

@$core
    .Deprecated('Use deleteDeliveryDestinationPolicyRequestDescriptor instead')
const DeleteDeliveryDestinationPolicyRequest$json = {
  '1': 'DeleteDeliveryDestinationPolicyRequest',
  '2': [
    {
      '1': 'deliverydestinationname',
      '3': 331358541,
      '4': 1,
      '5': 9,
      '10': 'deliverydestinationname'
    },
  ],
};

/// Descriptor for `DeleteDeliveryDestinationPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDeliveryDestinationPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiZEZWxldGVEZWxpdmVyeURlc3RpbmF0aW9uUG9saWN5UmVxdWVzdBI8ChdkZWxpdmVyeWRlc3'
        'RpbmF0aW9ubmFtZRjNwoCeASABKAlSF2RlbGl2ZXJ5ZGVzdGluYXRpb25uYW1l');

@$core.Deprecated('Use deleteDeliveryDestinationRequestDescriptor instead')
const DeleteDeliveryDestinationRequest$json = {
  '1': 'DeleteDeliveryDestinationRequest',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteDeliveryDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDeliveryDestinationRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVEZWxpdmVyeURlc3RpbmF0aW9uUmVxdWVzdBIVCgRuYW1lGOf75mkgASgJUgRuYW'
        '1l');

@$core.Deprecated('Use deleteDeliveryRequestDescriptor instead')
const DeleteDeliveryRequest$json = {
  '1': 'DeleteDeliveryRequest',
  '2': [
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteDeliveryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDeliveryRequestDescriptor =
    $convert.base64Decode(
        'ChVEZWxldGVEZWxpdmVyeVJlcXVlc3QSEgoCaWQY4dXhuQEgASgJUgJpZA==');

@$core.Deprecated('Use deleteDeliverySourceRequestDescriptor instead')
const DeleteDeliverySourceRequest$json = {
  '1': 'DeleteDeliverySourceRequest',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteDeliverySourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDeliverySourceRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVEZWxpdmVyeVNvdXJjZVJlcXVlc3QSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZQ==');

@$core.Deprecated('Use deleteDestinationRequestDescriptor instead')
const DeleteDestinationRequest$json = {
  '1': 'DeleteDestinationRequest',
  '2': [
    {
      '1': 'destinationname',
      '3': 284844189,
      '4': 1,
      '5': 9,
      '10': 'destinationname'
    },
  ],
};

/// Descriptor for `DeleteDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDestinationRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVEZXN0aW5hdGlvblJlcXVlc3QSLAoPZGVzdGluYXRpb25uYW1lGJ3B6YcBIAEoCV'
        'IPZGVzdGluYXRpb25uYW1l');

@$core.Deprecated('Use deleteIndexPolicyRequestDescriptor instead')
const DeleteIndexPolicyRequest$json = {
  '1': 'DeleteIndexPolicyRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
  ],
};

/// Descriptor for `DeleteIndexPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIndexPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVJbmRleFBvbGljeVJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIA'
        'EoCVISbG9nZ3JvdXBpZGVudGlmaWVy');

@$core.Deprecated('Use deleteIndexPolicyResponseDescriptor instead')
const DeleteIndexPolicyResponse$json = {
  '1': 'DeleteIndexPolicyResponse',
};

/// Descriptor for `DeleteIndexPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIndexPolicyResponseDescriptor =
    $convert.base64Decode('ChlEZWxldGVJbmRleFBvbGljeVJlc3BvbnNl');

@$core.Deprecated('Use deleteIntegrationRequestDescriptor instead')
const DeleteIntegrationRequest$json = {
  '1': 'DeleteIntegrationRequest',
  '2': [
    {'1': 'force', '3': 430540933, '4': 1, '5': 8, '10': 'force'},
    {
      '1': 'integrationname',
      '3': 183183535,
      '4': 1,
      '5': 9,
      '10': 'integrationname'
    },
  ],
};

/// Descriptor for `DeleteIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIntegrationRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVJbnRlZ3JhdGlvblJlcXVlc3QSGAoFZm9yY2UYhZGmzQEgASgIUgVmb3JjZRIrCg'
        '9pbnRlZ3JhdGlvbm5hbWUYr9GsVyABKAlSD2ludGVncmF0aW9ubmFtZQ==');

@$core.Deprecated('Use deleteIntegrationResponseDescriptor instead')
const DeleteIntegrationResponse$json = {
  '1': 'DeleteIntegrationResponse',
};

/// Descriptor for `DeleteIntegrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIntegrationResponseDescriptor =
    $convert.base64Decode('ChlEZWxldGVJbnRlZ3JhdGlvblJlc3BvbnNl');

@$core.Deprecated('Use deleteKeysDescriptor instead')
const DeleteKeys$json = {
  '1': 'DeleteKeys',
  '2': [
    {'1': 'withkeys', '3': 161392106, '4': 3, '5': 9, '10': 'withkeys'},
  ],
};

/// Descriptor for `DeleteKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteKeysDescriptor = $convert.base64Decode(
    'CgpEZWxldGVLZXlzEh0KCHdpdGhrZXlzGOrL+kwgAygJUgh3aXRoa2V5cw==');

@$core.Deprecated('Use deleteLogAnomalyDetectorRequestDescriptor instead')
const DeleteLogAnomalyDetectorRequest$json = {
  '1': 'DeleteLogAnomalyDetectorRequest',
  '2': [
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
  ],
};

/// Descriptor for `DeleteLogAnomalyDetectorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteLogAnomalyDetectorRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVMb2dBbm9tYWx5RGV0ZWN0b3JSZXF1ZXN0EjIKEmFub21hbHlkZXRlY3RvcmFybh'
        'isz/PUASABKAlSEmFub21hbHlkZXRlY3RvcmFybg==');

@$core.Deprecated('Use deleteLogGroupRequestDescriptor instead')
const DeleteLogGroupRequest$json = {
  '1': 'DeleteLogGroupRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
  ],
};

/// Descriptor for `DeleteLogGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteLogGroupRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVMb2dHcm91cFJlcXVlc3QSJQoMbG9nZ3JvdXBuYW1lGMy/1C0gASgJUgxsb2dncm'
    '91cG5hbWU=');

@$core.Deprecated('Use deleteLogStreamRequestDescriptor instead')
const DeleteLogStreamRequest$json = {
  '1': 'DeleteLogStreamRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
  ],
};

/// Descriptor for `DeleteLogStreamRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteLogStreamRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVMb2dTdHJlYW1SZXF1ZXN0EiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3'
        'JvdXBuYW1lEigKDWxvZ3N0cmVhbW5hbWUYo/fu0AEgASgJUg1sb2dzdHJlYW1uYW1l');

@$core.Deprecated('Use deleteMetricFilterRequestDescriptor instead')
const DeleteMetricFilterRequest$json = {
  '1': 'DeleteMetricFilterRequest',
  '2': [
    {'1': 'filtername', '3': 395125013, '4': 1, '5': 9, '10': 'filtername'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
  ],
};

/// Descriptor for `DeleteMetricFilterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMetricFilterRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVNZXRyaWNGaWx0ZXJSZXF1ZXN0EiIKCmZpbHRlcm5hbWUYlcK0vAEgASgJUgpmaW'
        'x0ZXJuYW1lEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW1l');

@$core.Deprecated('Use deleteQueryDefinitionRequestDescriptor instead')
const DeleteQueryDefinitionRequest$json = {
  '1': 'DeleteQueryDefinitionRequest',
  '2': [
    {
      '1': 'querydefinitionid',
      '3': 178455620,
      '4': 1,
      '5': 9,
      '10': 'querydefinitionid'
    },
  ],
};

/// Descriptor for `DeleteQueryDefinitionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteQueryDefinitionRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVRdWVyeURlZmluaXRpb25SZXF1ZXN0Ei8KEXF1ZXJ5ZGVmaW5pdGlvbmlkGMSIjF'
        'UgASgJUhFxdWVyeWRlZmluaXRpb25pZA==');

@$core.Deprecated('Use deleteQueryDefinitionResponseDescriptor instead')
const DeleteQueryDefinitionResponse$json = {
  '1': 'DeleteQueryDefinitionResponse',
  '2': [
    {'1': 'success', '3': 442482449, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `DeleteQueryDefinitionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteQueryDefinitionResponseDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVRdWVyeURlZmluaXRpb25SZXNwb25zZRIcCgdzdWNjZXNzGJH+/tIBIAEoCFIHc3'
        'VjY2Vzcw==');

@$core.Deprecated('Use deleteResourcePolicyRequestDescriptor instead')
const DeleteResourcePolicyRequest$json = {
  '1': 'DeleteResourcePolicyRequest',
  '2': [
    {
      '1': 'expectedrevisionid',
      '3': 469613378,
      '4': 1,
      '5': 9,
      '10': 'expectedrevisionid'
    },
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `DeleteResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVSZXNvdXJjZVBvbGljeVJlcXVlc3QSMgoSZXhwZWN0ZWRyZXZpc2lvbmlkGML29t'
        '8BIAEoCVISZXhwZWN0ZWRyZXZpc2lvbmlkEiEKCnBvbGljeW5hbWUY3eLyNiABKAlSCnBvbGlj'
        'eW5hbWUSIwoLcmVzb3VyY2Vhcm4YzcyqICABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use deleteRetentionPolicyRequestDescriptor instead')
const DeleteRetentionPolicyRequest$json = {
  '1': 'DeleteRetentionPolicyRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
  ],
};

/// Descriptor for `DeleteRetentionPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteRetentionPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVSZXRlbnRpb25Qb2xpY3lSZXF1ZXN0EiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCV'
        'IMbG9nZ3JvdXBuYW1l');

@$core.Deprecated('Use deleteScheduledQueryRequestDescriptor instead')
const DeleteScheduledQueryRequest$json = {
  '1': 'DeleteScheduledQueryRequest',
  '2': [
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `DeleteScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteScheduledQueryRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVTY2hlZHVsZWRRdWVyeVJlcXVlc3QSIQoKaWRlbnRpZmllchi/0JZFIAEoCVIKaW'
        'RlbnRpZmllcg==');

@$core.Deprecated('Use deleteScheduledQueryResponseDescriptor instead')
const DeleteScheduledQueryResponse$json = {
  '1': 'DeleteScheduledQueryResponse',
};

/// Descriptor for `DeleteScheduledQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteScheduledQueryResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVTY2hlZHVsZWRRdWVyeVJlc3BvbnNl');

@$core.Deprecated('Use deleteSubscriptionFilterRequestDescriptor instead')
const DeleteSubscriptionFilterRequest$json = {
  '1': 'DeleteSubscriptionFilterRequest',
  '2': [
    {'1': 'filtername', '3': 395125013, '4': 1, '5': 9, '10': 'filtername'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
  ],
};

/// Descriptor for `DeleteSubscriptionFilterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSubscriptionFilterRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVTdWJzY3JpcHRpb25GaWx0ZXJSZXF1ZXN0EiIKCmZpbHRlcm5hbWUYlcK0vAEgAS'
        'gJUgpmaWx0ZXJuYW1lEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW1l');

@$core.Deprecated('Use deleteTransformerRequestDescriptor instead')
const DeleteTransformerRequest$json = {
  '1': 'DeleteTransformerRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
  ],
};

/// Descriptor for `DeleteTransformerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTransformerRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVUcmFuc2Zvcm1lclJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIA'
        'EoCVISbG9nZ3JvdXBpZGVudGlmaWVy');

@$core.Deprecated('Use deliveryDescriptor instead')
const Delivery$json = {
  '1': 'Delivery',
  '2': [
    {'1': 'arn', '3': 359604989, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'deliverydestinationarn',
      '3': 156800339,
      '4': 1,
      '5': 9,
      '10': 'deliverydestinationarn'
    },
    {
      '1': 'deliverydestinationtype',
      '3': 33253176,
      '4': 1,
      '5': 14,
      '6': '.logs.DeliveryDestinationType',
      '10': 'deliverydestinationtype'
    },
    {
      '1': 'deliverysourcename',
      '3': 277065836,
      '4': 1,
      '5': 9,
      '10': 'deliverysourcename'
    },
    {
      '1': 'fielddelimiter',
      '3': 433214845,
      '4': 1,
      '5': 9,
      '10': 'fielddelimiter'
    },
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'recordfields', '3': 309625178, '4': 3, '5': 9, '10': 'recordfields'},
    {
      '1': 's3deliveryconfiguration',
      '3': 95221108,
      '4': 1,
      '5': 11,
      '6': '.logs.S3DeliveryConfiguration',
      '10': 's3deliveryconfiguration'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.Delivery.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [Delivery_TagsEntry$json],
};

@$core.Deprecated('Use deliveryDescriptor instead')
const Delivery_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Delivery`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deliveryDescriptor = $convert.base64Decode(
    'CghEZWxpdmVyeRIUCgNhcm4Y/cW8qwEgASgJUgNhcm4SOQoWZGVsaXZlcnlkZXN0aW5hdGlvbm'
    'FybhjTquJKIAEoCVIWZGVsaXZlcnlkZXN0aW5hdGlvbmFybhJaChdkZWxpdmVyeWRlc3RpbmF0'
    'aW9udHlwZRi4zu0PIAEoDjIdLmxvZ3MuRGVsaXZlcnlEZXN0aW5hdGlvblR5cGVSF2RlbGl2ZX'
    'J5ZGVzdGluYXRpb250eXBlEjIKEmRlbGl2ZXJ5c291cmNlbmFtZRjs4I6EASABKAlSEmRlbGl2'
    'ZXJ5c291cmNlbmFtZRIqCg5maWVsZGRlbGltaXRlchj9qsnOASABKAlSDmZpZWxkZGVsaW1pdG'
    'VyEhIKAmlkGOHV4bkBIAEoCVICaWQSJgoMcmVjb3JkZmllbGRzGNqC0pMBIAMoCVIMcmVjb3Jk'
    'ZmllbGRzEloKF3MzZGVsaXZlcnljb25maWd1cmF0aW9uGPTqsy0gASgLMh0ubG9ncy5TM0RlbG'
    'l2ZXJ5Q29uZmlndXJhdGlvblIXczNkZWxpdmVyeWNvbmZpZ3VyYXRpb24SMAoEdGFncxih19ug'
    'ASADKAsyGC5sb2dzLkRlbGl2ZXJ5LlRhZ3NFbnRyeVIEdGFncxo3CglUYWdzRW50cnkSEAoDa2'
    'V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use deliveryDestinationDescriptor instead')
const DeliveryDestination$json = {
  '1': 'DeliveryDestination',
  '2': [
    {'1': 'arn', '3': 359604989, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'deliverydestinationconfiguration',
      '3': 256504432,
      '4': 1,
      '5': 11,
      '6': '.logs.DeliveryDestinationConfiguration',
      '10': 'deliverydestinationconfiguration'
    },
    {
      '1': 'deliverydestinationtype',
      '3': 33253176,
      '4': 1,
      '5': 14,
      '6': '.logs.DeliveryDestinationType',
      '10': 'deliverydestinationtype'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputformat',
      '3': 217347480,
      '4': 1,
      '5': 14,
      '6': '.logs.OutputFormat',
      '10': 'outputformat'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.DeliveryDestination.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [DeliveryDestination_TagsEntry$json],
};

@$core.Deprecated('Use deliveryDestinationDescriptor instead')
const DeliveryDestination_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DeliveryDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deliveryDestinationDescriptor = $convert.base64Decode(
    'ChNEZWxpdmVyeURlc3RpbmF0aW9uEhQKA2Fybhj9xbyrASABKAlSA2FybhJ1CiBkZWxpdmVyeW'
    'Rlc3RpbmF0aW9uY29uZmlndXJhdGlvbhjw5Kd6IAEoCzImLmxvZ3MuRGVsaXZlcnlEZXN0aW5h'
    'dGlvbkNvbmZpZ3VyYXRpb25SIGRlbGl2ZXJ5ZGVzdGluYXRpb25jb25maWd1cmF0aW9uEloKF2'
    'RlbGl2ZXJ5ZGVzdGluYXRpb250eXBlGLjO7Q8gASgOMh0ubG9ncy5EZWxpdmVyeURlc3RpbmF0'
    'aW9uVHlwZVIXZGVsaXZlcnlkZXN0aW5hdGlvbnR5cGUSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZR'
    'I5CgxvdXRwdXRmb3JtYXQYmOvRZyABKA4yEi5sb2dzLk91dHB1dEZvcm1hdFIMb3V0cHV0Zm9y'
    'bWF0EjsKBHRhZ3MYodfboAEgAygLMiMubG9ncy5EZWxpdmVyeURlc3RpbmF0aW9uLlRhZ3NFbn'
    'RyeVIEdGFncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4AQ==');

@$core.Deprecated('Use deliveryDestinationConfigurationDescriptor instead')
const DeliveryDestinationConfiguration$json = {
  '1': 'DeliveryDestinationConfiguration',
  '2': [
    {
      '1': 'destinationresourcearn',
      '3': 430428943,
      '4': 1,
      '5': 9,
      '10': 'destinationresourcearn'
    },
  ],
};

/// Descriptor for `DeliveryDestinationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deliveryDestinationConfigurationDescriptor =
    $convert.base64Decode(
        'CiBEZWxpdmVyeURlc3RpbmF0aW9uQ29uZmlndXJhdGlvbhI6ChZkZXN0aW5hdGlvbnJlc291cm'
        'NlYXJuGI+mn80BIAEoCVIWZGVzdGluYXRpb25yZXNvdXJjZWFybg==');

@$core.Deprecated('Use deliverySourceDescriptor instead')
const DeliverySource$json = {
  '1': 'DeliverySource',
  '2': [
    {'1': 'arn', '3': 359604989, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'logtype', '3': 257838938, '4': 1, '5': 9, '10': 'logtype'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'resourcearns', '3': 465810734, '4': 3, '5': 9, '10': 'resourcearns'},
    {'1': 'service', '3': 383770213, '4': 1, '5': 9, '10': 'service'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.DeliverySource.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [DeliverySource_TagsEntry$json],
};

@$core.Deprecated('Use deliverySourceDescriptor instead')
const DeliverySource_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DeliverySource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deliverySourceDescriptor = $convert.base64Decode(
    'Cg5EZWxpdmVyeVNvdXJjZRIUCgNhcm4Y/cW8qwEgASgJUgNhcm4SGwoHbG9ndHlwZRjanvl6IA'
    'EoCVIHbG9ndHlwZRIVCgRuYW1lGOf75mkgASgJUgRuYW1lEiYKDHJlc291cmNlYXJucxiu6o7e'
    'ASADKAlSDHJlc291cmNlYXJucxIcCgdzZXJ2aWNlGOW8/7YBIAEoCVIHc2VydmljZRI2CgR0YW'
    'dzGKHX26ABIAMoCzIeLmxvZ3MuRGVsaXZlcnlTb3VyY2UuVGFnc0VudHJ5UgR0YWdzGjcKCVRh'
    'Z3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use describeAccountPoliciesRequestDescriptor instead')
const DescribeAccountPoliciesRequest$json = {
  '1': 'DescribeAccountPoliciesRequest',
  '2': [
    {
      '1': 'accountidentifiers',
      '3': 349304053,
      '4': 3,
      '5': 9,
      '10': 'accountidentifiers'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'policytype',
      '3': 319277736,
      '4': 1,
      '5': 14,
      '6': '.logs.PolicyType',
      '10': 'policytype'
    },
  ],
};

/// Descriptor for `DescribeAccountPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAccountPoliciesRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZUFjY291bnRQb2xpY2llc1JlcXVlc3QSMgoSYWNjb3VudGlkZW50aWZpZXJzGP'
        'Xpx6YBIAMoCVISYWNjb3VudGlkZW50aWZpZXJzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4'
        'dHRva2VuEiEKCnBvbGljeW5hbWUY3eLyNiABKAlSCnBvbGljeW5hbWUSNAoKcG9saWN5dHlwZR'
        'iolZ+YASABKA4yEC5sb2dzLlBvbGljeVR5cGVSCnBvbGljeXR5cGU=');

@$core.Deprecated('Use describeAccountPoliciesResponseDescriptor instead')
const DescribeAccountPoliciesResponse$json = {
  '1': 'DescribeAccountPoliciesResponse',
  '2': [
    {
      '1': 'accountpolicies',
      '3': 43946431,
      '4': 3,
      '5': 11,
      '6': '.logs.AccountPolicy',
      '10': 'accountpolicies'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAccountPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAccountPoliciesResponseDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZUFjY291bnRQb2xpY2llc1Jlc3BvbnNlEkAKD2FjY291bnRwb2xpY2llcxi/o/'
        'oUIAMoCzITLmxvZ3MuQWNjb3VudFBvbGljeVIPYWNjb3VudHBvbGljaWVzEh8KCW5leHR0b2tl'
        'bhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeConfigurationTemplatesRequestDescriptor instead')
const DescribeConfigurationTemplatesRequest$json = {
  '1': 'DescribeConfigurationTemplatesRequest',
  '2': [
    {
      '1': 'deliverydestinationtypes',
      '3': 46370829,
      '4': 3,
      '5': 14,
      '6': '.logs.DeliveryDestinationType',
      '10': 'deliverydestinationtypes'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'logtypes', '3': 460758815, '4': 3, '5': 9, '10': 'logtypes'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcetypes',
      '3': 78421387,
      '4': 3,
      '5': 9,
      '10': 'resourcetypes'
    },
    {'1': 'service', '3': 383770213, '4': 1, '5': 9, '10': 'service'},
  ],
};

/// Descriptor for `DescribeConfigurationTemplatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeConfigurationTemplatesRequestDescriptor = $convert.base64Decode(
    'CiVEZXNjcmliZUNvbmZpZ3VyYXRpb25UZW1wbGF0ZXNSZXF1ZXN0ElwKGGRlbGl2ZXJ5ZGVzdG'
    'luYXRpb250eXBlcxiNoI4WIAMoDjIdLmxvZ3MuRGVsaXZlcnlEZXN0aW5hdGlvblR5cGVSGGRl'
    'bGl2ZXJ5ZGVzdGluYXRpb250eXBlcxIYCgVsaW1pdBi1suuWASABKAVSBWxpbWl0Eh4KCGxvZ3'
    'R5cGVzGJ++2tsBIAMoCVIIbG9ndHlwZXMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9r'
    'ZW4SJwoNcmVzb3VyY2V0eXBlcxiLu7IlIAMoCVINcmVzb3VyY2V0eXBlcxIcCgdzZXJ2aWNlGO'
    'W8/7YBIAEoCVIHc2VydmljZQ==');

@$core
    .Deprecated('Use describeConfigurationTemplatesResponseDescriptor instead')
const DescribeConfigurationTemplatesResponse$json = {
  '1': 'DescribeConfigurationTemplatesResponse',
  '2': [
    {
      '1': 'configurationtemplates',
      '3': 91645093,
      '4': 3,
      '5': 11,
      '6': '.logs.ConfigurationTemplate',
      '10': 'configurationtemplates'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeConfigurationTemplatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeConfigurationTemplatesResponseDescriptor =
    $convert.base64Decode(
        'CiZEZXNjcmliZUNvbmZpZ3VyYXRpb25UZW1wbGF0ZXNSZXNwb25zZRJWChZjb25maWd1cmF0aW'
        '9udGVtcGxhdGVzGKXJ2SsgAygLMhsubG9ncy5Db25maWd1cmF0aW9uVGVtcGxhdGVSFmNvbmZp'
        'Z3VyYXRpb250ZW1wbGF0ZXMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeDeliveriesRequestDescriptor instead')
const DescribeDeliveriesRequest$json = {
  '1': 'DescribeDeliveriesRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDeliveriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDeliveriesRequestDescriptor =
    $convert.base64Decode(
        'ChlEZXNjcmliZURlbGl2ZXJpZXNSZXF1ZXN0EhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSHw'
        'oJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeDeliveriesResponseDescriptor instead')
const DescribeDeliveriesResponse$json = {
  '1': 'DescribeDeliveriesResponse',
  '2': [
    {
      '1': 'deliveries',
      '3': 154341066,
      '4': 3,
      '5': 11,
      '6': '.logs.Delivery',
      '10': 'deliveries'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDeliveriesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDeliveriesResponseDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZURlbGl2ZXJpZXNSZXNwb25zZRIxCgpkZWxpdmVyaWVzGMqdzEkgAygLMg4ubG'
        '9ncy5EZWxpdmVyeVIKZGVsaXZlcmllcxIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tl'
        'bg==');

@$core.Deprecated('Use describeDeliveryDestinationsRequestDescriptor instead')
const DescribeDeliveryDestinationsRequest$json = {
  '1': 'DescribeDeliveryDestinationsRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDeliveryDestinationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDeliveryDestinationsRequestDescriptor =
    $convert.base64Decode(
        'CiNEZXNjcmliZURlbGl2ZXJ5RGVzdGluYXRpb25zUmVxdWVzdBIYCgVsaW1pdBi1suuWASABKA'
        'VSBWxpbWl0Eh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeDeliveryDestinationsResponseDescriptor instead')
const DescribeDeliveryDestinationsResponse$json = {
  '1': 'DescribeDeliveryDestinationsResponse',
  '2': [
    {
      '1': 'deliverydestinations',
      '3': 236883221,
      '4': 3,
      '5': 11,
      '6': '.logs.DeliveryDestination',
      '10': 'deliverydestinations'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDeliveryDestinationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDeliveryDestinationsResponseDescriptor =
    $convert.base64Decode(
        'CiREZXNjcmliZURlbGl2ZXJ5RGVzdGluYXRpb25zUmVzcG9uc2USUAoUZGVsaXZlcnlkZXN0aW'
        '5hdGlvbnMYlZr6cCADKAsyGS5sb2dzLkRlbGl2ZXJ5RGVzdGluYXRpb25SFGRlbGl2ZXJ5ZGVz'
        'dGluYXRpb25zEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeDeliverySourcesRequestDescriptor instead')
const DescribeDeliverySourcesRequest$json = {
  '1': 'DescribeDeliverySourcesRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDeliverySourcesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDeliverySourcesRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZURlbGl2ZXJ5U291cmNlc1JlcXVlc3QSGAoFbGltaXQYtbLrlgEgASgFUgVsaW'
        '1pdBIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeDeliverySourcesResponseDescriptor instead')
const DescribeDeliverySourcesResponse$json = {
  '1': 'DescribeDeliverySourcesResponse',
  '2': [
    {
      '1': 'deliverysources',
      '3': 5432066,
      '4': 3,
      '5': 11,
      '6': '.logs.DeliverySource',
      '10': 'deliverysources'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDeliverySourcesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDeliverySourcesResponseDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZURlbGl2ZXJ5U291cmNlc1Jlc3BvbnNlEkEKD2RlbGl2ZXJ5c291cmNlcxiCxs'
        'sCIAMoCzIULmxvZ3MuRGVsaXZlcnlTb3VyY2VSD2RlbGl2ZXJ5c291cmNlcxIfCgluZXh0dG9r'
        'ZW4YnvOdNyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeDestinationsRequestDescriptor instead')
const DescribeDestinationsRequest$json = {
  '1': 'DescribeDestinationsRequest',
  '2': [
    {
      '1': 'destinationnameprefix',
      '3': 160941145,
      '4': 1,
      '5': 9,
      '10': 'destinationnameprefix'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDestinationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDestinationsRequestDescriptor =
    $convert.base64Decode(
        'ChtEZXNjcmliZURlc3RpbmF0aW9uc1JlcXVlc3QSNwoVZGVzdGluYXRpb25uYW1lcHJlZml4GN'
        'mI30wgASgJUhVkZXN0aW5hdGlvbm5hbWVwcmVmaXgSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1p'
        'dBIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeDestinationsResponseDescriptor instead')
const DescribeDestinationsResponse$json = {
  '1': 'DescribeDestinationsResponse',
  '2': [
    {
      '1': 'destinations',
      '3': 1617189,
      '4': 3,
      '5': 11,
      '6': '.logs.Destination',
      '10': 'destinations'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDestinationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDestinationsResponseDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZURlc3RpbmF0aW9uc1Jlc3BvbnNlEjcKDGRlc3RpbmF0aW9ucxil2mIgAygLMh'
        'EubG9ncy5EZXN0aW5hdGlvblIMZGVzdGluYXRpb25zEh8KCW5leHR0b2tlbhie8503IAEoCVIJ'
        'bmV4dHRva2Vu');

@$core.Deprecated('Use describeExportTasksRequestDescriptor instead')
const DescribeExportTasksRequest$json = {
  '1': 'DescribeExportTasksRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statuscode',
      '3': 299352223,
      '4': 1,
      '5': 14,
      '6': '.logs.ExportTaskStatusCode',
      '10': 'statuscode'
    },
    {'1': 'taskid', '3': 216769858, '4': 1, '5': 9, '10': 'taskid'},
  ],
};

/// Descriptor for `DescribeExportTasksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeExportTasksRequestDescriptor = $convert.base64Decode(
    'ChpEZXNjcmliZUV4cG9ydFRhc2tzUmVxdWVzdBIYCgVsaW1pdBi1suuWASABKAVSBWxpbWl0Eh'
    '8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2VuEj4KCnN0YXR1c2NvZGUYn4HfjgEgASgO'
    'MhoubG9ncy5FeHBvcnRUYXNrU3RhdHVzQ29kZVIKc3RhdHVzY29kZRIZCgZ0YXNraWQYwsquZy'
    'ABKAlSBnRhc2tpZA==');

@$core.Deprecated('Use describeExportTasksResponseDescriptor instead')
const DescribeExportTasksResponse$json = {
  '1': 'DescribeExportTasksResponse',
  '2': [
    {
      '1': 'exporttasks',
      '3': 40181908,
      '4': 3,
      '5': 11,
      '6': '.logs.ExportTask',
      '10': 'exporttasks'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeExportTasksResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeExportTasksResponseDescriptor =
    $convert.base64Decode(
        'ChtEZXNjcmliZUV4cG9ydFRhc2tzUmVzcG9uc2USNQoLZXhwb3J0dGFza3MYlMGUEyADKAsyEC'
        '5sb2dzLkV4cG9ydFRhc2tSC2V4cG9ydHRhc2tzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4'
        'dHRva2Vu');

@$core.Deprecated('Use describeFieldIndexesRequestDescriptor instead')
const DescribeFieldIndexesRequest$json = {
  '1': 'DescribeFieldIndexesRequest',
  '2': [
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeFieldIndexesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeFieldIndexesRequestDescriptor =
    $convert.base64Decode(
        'ChtEZXNjcmliZUZpZWxkSW5kZXhlc1JlcXVlc3QSNAoTbG9nZ3JvdXBpZGVudGlmaWVycxj52r'
        'jDASADKAlSE2xvZ2dyb3VwaWRlbnRpZmllcnMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0'
        'dG9rZW4=');

@$core.Deprecated('Use describeFieldIndexesResponseDescriptor instead')
const DescribeFieldIndexesResponse$json = {
  '1': 'DescribeFieldIndexesResponse',
  '2': [
    {
      '1': 'fieldindexes',
      '3': 138171666,
      '4': 3,
      '5': 11,
      '6': '.logs.FieldIndex',
      '10': 'fieldindexes'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeFieldIndexesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeFieldIndexesResponseDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZUZpZWxkSW5kZXhlc1Jlc3BvbnNlEjcKDGZpZWxkaW5kZXhlcxiSqvFBIAMoCz'
        'IQLmxvZ3MuRmllbGRJbmRleFIMZmllbGRpbmRleGVzEh8KCW5leHR0b2tlbhie8503IAEoCVIJ'
        'bmV4dHRva2Vu');

@$core.Deprecated('Use describeImportTaskBatchesRequestDescriptor instead')
const DescribeImportTaskBatchesRequest$json = {
  '1': 'DescribeImportTaskBatchesRequest',
  '2': [
    {
      '1': 'batchimportstatus',
      '3': 117307821,
      '4': 3,
      '5': 14,
      '6': '.logs.ImportStatus',
      '10': 'batchimportstatus'
    },
    {'1': 'importid', '3': 513429114, '4': 1, '5': 9, '10': 'importid'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeImportTaskBatchesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeImportTaskBatchesRequestDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZUltcG9ydFRhc2tCYXRjaGVzUmVxdWVzdBJDChFiYXRjaGltcG9ydHN0YXR1cx'
        'it8/c3IAMoDjISLmxvZ3MuSW1wb3J0U3RhdHVzUhFiYXRjaGltcG9ydHN0YXR1cxIeCghpbXBv'
        'cnRpZBj6nOn0ASABKAlSCGltcG9ydGlkEhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSHwoJbm'
        'V4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeImportTaskBatchesResponseDescriptor instead')
const DescribeImportTaskBatchesResponse$json = {
  '1': 'DescribeImportTaskBatchesResponse',
  '2': [
    {
      '1': 'importbatches',
      '3': 506952441,
      '4': 3,
      '5': 11,
      '6': '.logs.ImportBatch',
      '10': 'importbatches'
    },
    {'1': 'importid', '3': 513429114, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importsourcearn',
      '3': 161570329,
      '4': 1,
      '5': 9,
      '10': 'importsourcearn'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeImportTaskBatchesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeImportTaskBatchesResponseDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZUltcG9ydFRhc2tCYXRjaGVzUmVzcG9uc2USOwoNaW1wb3J0YmF0Y2hlcxj59d'
        '3xASADKAsyES5sb2dzLkltcG9ydEJhdGNoUg1pbXBvcnRiYXRjaGVzEh4KCGltcG9ydGlkGPqc'
        '6fQBIAEoCVIIaW1wb3J0aWQSKwoPaW1wb3J0c291cmNlYXJuGJm8hU0gASgJUg9pbXBvcnRzb3'
        'VyY2Vhcm4SHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeImportTasksRequestDescriptor instead')
const DescribeImportTasksRequest$json = {
  '1': 'DescribeImportTasksRequest',
  '2': [
    {'1': 'importid', '3': 513429114, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importsourcearn',
      '3': 161570329,
      '4': 1,
      '5': 9,
      '10': 'importsourcearn'
    },
    {
      '1': 'importstatus',
      '3': 31427999,
      '4': 1,
      '5': 14,
      '6': '.logs.ImportStatus',
      '10': 'importstatus'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeImportTasksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeImportTasksRequestDescriptor = $convert.base64Decode(
    'ChpEZXNjcmliZUltcG9ydFRhc2tzUmVxdWVzdBIeCghpbXBvcnRpZBj6nOn0ASABKAlSCGltcG'
    '9ydGlkEisKD2ltcG9ydHNvdXJjZWFybhiZvIVNIAEoCVIPaW1wb3J0c291cmNlYXJuEjkKDGlt'
    'cG9ydHN0YXR1cxifm/4OIAEoDjISLmxvZ3MuSW1wb3J0U3RhdHVzUgxpbXBvcnRzdGF0dXMSGA'
    'oFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tl'
    'bg==');

@$core.Deprecated('Use describeImportTasksResponseDescriptor instead')
const DescribeImportTasksResponse$json = {
  '1': 'DescribeImportTasksResponse',
  '2': [
    {
      '1': 'imports',
      '3': 218216166,
      '4': 3,
      '5': 11,
      '6': '.logs.Import',
      '10': 'imports'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeImportTasksResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeImportTasksResponseDescriptor =
    $convert.base64Decode(
        'ChtEZXNjcmliZUltcG9ydFRhc2tzUmVzcG9uc2USKQoHaW1wb3J0cxjm7YZoIAMoCzIMLmxvZ3'
        'MuSW1wb3J0UgdpbXBvcnRzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeIndexPoliciesRequestDescriptor instead')
const DescribeIndexPoliciesRequest$json = {
  '1': 'DescribeIndexPoliciesRequest',
  '2': [
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeIndexPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeIndexPoliciesRequestDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZUluZGV4UG9saWNpZXNSZXF1ZXN0EjQKE2xvZ2dyb3VwaWRlbnRpZmllcnMY+d'
        'q4wwEgAygJUhNsb2dncm91cGlkZW50aWZpZXJzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4'
        'dHRva2Vu');

@$core.Deprecated('Use describeIndexPoliciesResponseDescriptor instead')
const DescribeIndexPoliciesResponse$json = {
  '1': 'DescribeIndexPoliciesResponse',
  '2': [
    {
      '1': 'indexpolicies',
      '3': 463340482,
      '4': 3,
      '5': 11,
      '6': '.logs.IndexPolicy',
      '10': 'indexpolicies'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeIndexPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeIndexPoliciesResponseDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZUluZGV4UG9saWNpZXNSZXNwb25zZRI7Cg1pbmRleHBvbGljaWVzGMKH+NwBIA'
        'MoCzIRLmxvZ3MuSW5kZXhQb2xpY3lSDWluZGV4cG9saWNpZXMSHwoJbmV4dHRva2VuGJ7znTcg'
        'ASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeLogGroupsRequestDescriptor instead')
const DescribeLogGroupsRequest$json = {
  '1': 'DescribeLogGroupsRequest',
  '2': [
    {
      '1': 'accountidentifiers',
      '3': 349304053,
      '4': 3,
      '5': 9,
      '10': 'accountidentifiers'
    },
    {
      '1': 'includelinkedaccounts',
      '3': 56034131,
      '4': 1,
      '5': 8,
      '10': 'includelinkedaccounts'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'loggroupclass',
      '3': 518605953,
      '4': 1,
      '5': 14,
      '6': '.logs.LogGroupClass',
      '10': 'loggroupclass'
    },
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {
      '1': 'loggroupnamepattern',
      '3': 253299540,
      '4': 1,
      '5': 9,
      '10': 'loggroupnamepattern'
    },
    {
      '1': 'loggroupnameprefix',
      '3': 255310944,
      '4': 1,
      '5': 9,
      '10': 'loggroupnameprefix'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeLogGroupsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLogGroupsRequestDescriptor = $convert.base64Decode(
    'ChhEZXNjcmliZUxvZ0dyb3Vwc1JlcXVlc3QSMgoSYWNjb3VudGlkZW50aWZpZXJzGPXpx6YBIA'
    'MoCVISYWNjb3VudGlkZW50aWZpZXJzEjcKFWluY2x1ZGVsaW5rZWRhY2NvdW50cxjThtwaIAEo'
    'CFIVaW5jbHVkZWxpbmtlZGFjY291bnRzEhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSPQoNbG'
    '9nZ3JvdXBjbGFzcxiBmaX3ASABKA4yEy5sb2dzLkxvZ0dyb3VwQ2xhc3NSDWxvZ2dyb3VwY2xh'
    'c3MSNAoTbG9nZ3JvdXBpZGVudGlmaWVycxj52rjDASADKAlSE2xvZ2dyb3VwaWRlbnRpZmllcn'
    'MSMwoTbG9nZ3JvdXBuYW1lcGF0dGVybhjUluR4IAEoCVITbG9nZ3JvdXBuYW1lcGF0dGVybhIx'
    'ChJsb2dncm91cG5hbWVwcmVmaXgY4PjeeSABKAlSEmxvZ2dyb3VwbmFtZXByZWZpeBIfCgluZX'
    'h0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeLogGroupsResponseDescriptor instead')
const DescribeLogGroupsResponse$json = {
  '1': 'DescribeLogGroupsResponse',
  '2': [
    {
      '1': 'loggroups',
      '3': 507247778,
      '4': 3,
      '5': 11,
      '6': '.logs.LogGroup',
      '10': 'loggroups'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeLogGroupsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLogGroupsResponseDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZUxvZ0dyb3Vwc1Jlc3BvbnNlEjAKCWxvZ2dyb3Vwcxii+e/xASADKAsyDi5sb2'
    'dzLkxvZ0dyb3VwUglsb2dncm91cHMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeLogStreamsRequestDescriptor instead')
const DescribeLogStreamsRequest$json = {
  '1': 'DescribeLogStreamsRequest',
  '2': [
    {'1': 'descending', '3': 378787644, '4': 1, '5': 8, '10': 'descending'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'logstreamnameprefix',
      '3': 24068871,
      '4': 1,
      '5': 9,
      '10': 'logstreamnameprefix'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'orderby',
      '3': 63062531,
      '4': 1,
      '5': 14,
      '6': '.logs.OrderBy',
      '10': 'orderby'
    },
  ],
};

/// Descriptor for `DescribeLogStreamsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLogStreamsRequestDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZUxvZ1N0cmVhbXNSZXF1ZXN0EiIKCmRlc2NlbmRpbmcYvK7PtAEgASgIUgpkZX'
    'NjZW5kaW5nEhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSMgoSbG9nZ3JvdXBpZGVudGlmaWVy'
    'GNzPpd8BIAEoCVISbG9nZ3JvdXBpZGVudGlmaWVyEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCV'
    'IMbG9nZ3JvdXBuYW1lEjMKE2xvZ3N0cmVhbW5hbWVwcmVmaXgYh4a9CyABKAlSE2xvZ3N0cmVh'
    'bW5hbWVwcmVmaXgSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4SKgoHb3JkZXJieR'
    'iDhIkeIAEoDjINLmxvZ3MuT3JkZXJCeVIHb3JkZXJieQ==');

@$core.Deprecated('Use describeLogStreamsResponseDescriptor instead')
const DescribeLogStreamsResponse$json = {
  '1': 'DescribeLogStreamsResponse',
  '2': [
    {
      '1': 'logstreams',
      '3': 201936239,
      '4': 3,
      '5': 11,
      '6': '.logs.LogStream',
      '10': 'logstreams'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeLogStreamsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLogStreamsResponseDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZUxvZ1N0cmVhbXNSZXNwb25zZRIyCgpsb2dzdHJlYW1zGO+apWAgAygLMg8ubG'
        '9ncy5Mb2dTdHJlYW1SCmxvZ3N0cmVhbXMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9r'
        'ZW4=');

@$core.Deprecated('Use describeMetricFiltersRequestDescriptor instead')
const DescribeMetricFiltersRequest$json = {
  '1': 'DescribeMetricFiltersRequest',
  '2': [
    {
      '1': 'filternameprefix',
      '3': 46988721,
      '4': 1,
      '5': 9,
      '10': 'filternameprefix'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'metricname', '3': 204020635, '4': 1, '5': 9, '10': 'metricname'},
    {
      '1': 'metricnamespace',
      '3': 315894261,
      '4': 1,
      '5': 9,
      '10': 'metricnamespace'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeMetricFiltersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMetricFiltersRequestDescriptor = $convert.base64Decode(
    'ChxEZXNjcmliZU1ldHJpY0ZpbHRlcnNSZXF1ZXN0Ei0KEGZpbHRlcm5hbWVwcmVmaXgYsfuzFi'
    'ABKAlSEGZpbHRlcm5hbWVwcmVmaXgSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIlCgxsb2dn'
    'cm91cG5hbWUYzL/ULSABKAlSDGxvZ2dyb3VwbmFtZRIhCgptZXRyaWNuYW1lGJu3pGEgASgJUg'
    'ptZXRyaWNuYW1lEiwKD21ldHJpY25hbWVzcGFjZRj109CWASABKAlSD21ldHJpY25hbWVzcGFj'
    'ZRIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeMetricFiltersResponseDescriptor instead')
const DescribeMetricFiltersResponse$json = {
  '1': 'DescribeMetricFiltersResponse',
  '2': [
    {
      '1': 'metricfilters',
      '3': 449173825,
      '4': 3,
      '5': 11,
      '6': '.logs.MetricFilter',
      '10': 'metricfilters'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeMetricFiltersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMetricFiltersResponseDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZU1ldHJpY0ZpbHRlcnNSZXNwb25zZRI8Cg1tZXRyaWNmaWx0ZXJzGMGyl9YBIA'
        'MoCzISLmxvZ3MuTWV0cmljRmlsdGVyUg1tZXRyaWNmaWx0ZXJzEh8KCW5leHR0b2tlbhie8503'
        'IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeQueriesRequestDescriptor instead')
const DescribeQueriesRequest$json = {
  '1': 'DescribeQueriesRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `DescribeQueriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeQueriesRequestDescriptor = $convert.base64Decode(
    'ChZEZXNjcmliZVF1ZXJpZXNSZXF1ZXN0EiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3'
    'JvdXBuYW1lEiIKCm1heHJlc3VsdHMY0tzn3QEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tl'
    'bhie8503IAEoCVIJbmV4dHRva2VuEj0KDXF1ZXJ5bGFuZ3VhZ2UY9q/3owEgASgOMhMubG9ncy'
    '5RdWVyeUxhbmd1YWdlUg1xdWVyeWxhbmd1YWdlEi0KBnN0YXR1cxjw763SASABKA4yES5sb2dz'
    'LlF1ZXJ5U3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use describeQueriesResponseDescriptor instead')
const DescribeQueriesResponse$json = {
  '1': 'DescribeQueriesResponse',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queries',
      '3': 26054724,
      '4': 3,
      '5': 11,
      '6': '.logs.QueryInfo',
      '10': 'queries'
    },
  ],
};

/// Descriptor for `DescribeQueriesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeQueriesResponseDescriptor =
    $convert.base64Decode(
        'ChdEZXNjcmliZVF1ZXJpZXNSZXNwb25zZRIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2'
        'tlbhIsCgdxdWVyaWVzGMSgtgwgAygLMg8ubG9ncy5RdWVyeUluZm9SB3F1ZXJpZXM=');

@$core.Deprecated('Use describeQueryDefinitionsRequestDescriptor instead')
const DescribeQueryDefinitionsRequest$json = {
  '1': 'DescribeQueryDefinitionsRequest',
  '2': [
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'querydefinitionnameprefix',
      '3': 28618146,
      '4': 1,
      '5': 9,
      '10': 'querydefinitionnameprefix'
    },
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
  ],
};

/// Descriptor for `DescribeQueryDefinitionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeQueryDefinitionsRequestDescriptor = $convert.base64Decode(
    'Ch9EZXNjcmliZVF1ZXJ5RGVmaW5pdGlvbnNSZXF1ZXN0EiIKCm1heHJlc3VsdHMY0tzn3QEgAS'
    'gFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2VuEj8KGXF1ZXJ5'
    'ZGVmaW5pdGlvbm5hbWVwcmVmaXgYotvSDSABKAlSGXF1ZXJ5ZGVmaW5pdGlvbm5hbWVwcmVmaX'
    'gSPQoNcXVlcnlsYW5ndWFnZRj2r/ejASABKA4yEy5sb2dzLlF1ZXJ5TGFuZ3VhZ2VSDXF1ZXJ5'
    'bGFuZ3VhZ2U=');

@$core.Deprecated('Use describeQueryDefinitionsResponseDescriptor instead')
const DescribeQueryDefinitionsResponse$json = {
  '1': 'DescribeQueryDefinitionsResponse',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'querydefinitions',
      '3': 195658648,
      '4': 3,
      '5': 11,
      '6': '.logs.QueryDefinition',
      '10': 'querydefinitions'
    },
  ],
};

/// Descriptor for `DescribeQueryDefinitionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeQueryDefinitionsResponseDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZVF1ZXJ5RGVmaW5pdGlvbnNSZXNwb25zZRIfCgluZXh0dG9rZW4YnvOdNyABKA'
        'lSCW5leHR0b2tlbhJEChBxdWVyeWRlZmluaXRpb25zGJiHpl0gAygLMhUubG9ncy5RdWVyeURl'
        'ZmluaXRpb25SEHF1ZXJ5ZGVmaW5pdGlvbnM=');

@$core.Deprecated('Use describeResourcePoliciesRequestDescriptor instead')
const DescribeResourcePoliciesRequest$json = {
  '1': 'DescribeResourcePoliciesRequest',
  '2': [
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'policyscope',
      '3': 288841470,
      '4': 1,
      '5': 14,
      '6': '.logs.PolicyScope',
      '10': 'policyscope'
    },
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `DescribeResourcePoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeResourcePoliciesRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZVJlc291cmNlUG9saWNpZXNSZXF1ZXN0EhgKBWxpbWl0GLWy65YBIAEoBVIFbG'
        'ltaXQSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4SNwoLcG9saWN5c2NvcGUY/r3d'
        'iQEgASgOMhEubG9ncy5Qb2xpY3lTY29wZVILcG9saWN5c2NvcGUSIwoLcmVzb3VyY2Vhcm4Yzc'
        'yqICABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use describeResourcePoliciesResponseDescriptor instead')
const DescribeResourcePoliciesResponse$json = {
  '1': 'DescribeResourcePoliciesResponse',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcepolicies',
      '3': 220332616,
      '4': 3,
      '5': 11,
      '6': '.logs.ResourcePolicy',
      '10': 'resourcepolicies'
    },
  ],
};

/// Descriptor for `DescribeResourcePoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeResourcePoliciesResponseDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZVJlc291cmNlUG9saWNpZXNSZXNwb25zZRIfCgluZXh0dG9rZW4YnvOdNyABKA'
        'lSCW5leHR0b2tlbhJDChByZXNvdXJjZXBvbGljaWVzGMiEiGkgAygLMhQubG9ncy5SZXNvdXJj'
        'ZVBvbGljeVIQcmVzb3VyY2Vwb2xpY2llcw==');

@$core.Deprecated('Use describeSubscriptionFiltersRequestDescriptor instead')
const DescribeSubscriptionFiltersRequest$json = {
  '1': 'DescribeSubscriptionFiltersRequest',
  '2': [
    {
      '1': 'filternameprefix',
      '3': 46988721,
      '4': 1,
      '5': 9,
      '10': 'filternameprefix'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeSubscriptionFiltersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeSubscriptionFiltersRequestDescriptor =
    $convert.base64Decode(
        'CiJEZXNjcmliZVN1YnNjcmlwdGlvbkZpbHRlcnNSZXF1ZXN0Ei0KEGZpbHRlcm5hbWVwcmVmaX'
        'gYsfuzFiABKAlSEGZpbHRlcm5hbWVwcmVmaXgSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIl'
        'Cgxsb2dncm91cG5hbWUYzL/ULSABKAlSDGxvZ2dyb3VwbmFtZRIfCgluZXh0dG9rZW4YnvOdNy'
        'ABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeSubscriptionFiltersResponseDescriptor instead')
const DescribeSubscriptionFiltersResponse$json = {
  '1': 'DescribeSubscriptionFiltersResponse',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'subscriptionfilters',
      '3': 453910816,
      '4': 3,
      '5': 11,
      '6': '.logs.SubscriptionFilter',
      '10': 'subscriptionfilters'
    },
  ],
};

/// Descriptor for `DescribeSubscriptionFiltersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeSubscriptionFiltersResponseDescriptor =
    $convert.base64Decode(
        'CiNEZXNjcmliZVN1YnNjcmlwdGlvbkZpbHRlcnNSZXNwb25zZRIfCgluZXh0dG9rZW4YnvOdNy'
        'ABKAlSCW5leHR0b2tlbhJOChNzdWJzY3JpcHRpb25maWx0ZXJzGKDCuNgBIAMoCzIYLmxvZ3Mu'
        'U3Vic2NyaXB0aW9uRmlsdGVyUhNzdWJzY3JpcHRpb25maWx0ZXJz');

@$core.Deprecated('Use destinationDescriptor instead')
const Destination$json = {
  '1': 'Destination',
  '2': [
    {'1': 'accesspolicy', '3': 336061518, '4': 1, '5': 9, '10': 'accesspolicy'},
    {'1': 'arn', '3': 359604989, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'destinationname',
      '3': 284844189,
      '4': 1,
      '5': 9,
      '10': 'destinationname'
    },
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {'1': 'targetarn', '3': 367964720, '4': 1, '5': 9, '10': 'targetarn'},
  ],
};

/// Descriptor for `Destination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List destinationDescriptor = $convert.base64Decode(
    'CgtEZXN0aW5hdGlvbhImCgxhY2Nlc3Nwb2xpY3kYzsifoAEgASgJUgxhY2Nlc3Nwb2xpY3kSFA'
    'oDYXJuGP3FvKsBIAEoCVIDYXJuEiUKDGNyZWF0aW9udGltZRiGxcQZIAEoA1IMY3JlYXRpb250'
    'aW1lEiwKD2Rlc3RpbmF0aW9ubmFtZRidwemHASABKAlSD2Rlc3RpbmF0aW9ubmFtZRIbCgdyb2'
    'xlYXJuGKGXiVEgASgJUgdyb2xlYXJuEiAKCXRhcmdldGFybhiw5LqvASABKAlSCXRhcmdldGFy'
    'bg==');

@$core.Deprecated('Use destinationConfigurationDescriptor instead')
const DestinationConfiguration$json = {
  '1': 'DestinationConfiguration',
  '2': [
    {
      '1': 's3configuration',
      '3': 506726172,
      '4': 1,
      '5': 11,
      '6': '.logs.S3Configuration',
      '10': 's3configuration'
    },
  ],
};

/// Descriptor for `DestinationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List destinationConfigurationDescriptor =
    $convert.base64Decode(
        'ChhEZXN0aW5hdGlvbkNvbmZpZ3VyYXRpb24SQwoPczNjb25maWd1cmF0aW9uGJyO0PEBIAEoCz'
        'IVLmxvZ3MuUzNDb25maWd1cmF0aW9uUg9zM2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use disassociateKmsKeyRequestDescriptor instead')
const DisassociateKmsKeyRequest$json = {
  '1': 'DisassociateKmsKeyRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'resourceidentifier',
      '3': 309427407,
      '4': 1,
      '5': 9,
      '10': 'resourceidentifier'
    },
  ],
};

/// Descriptor for `DisassociateKmsKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disassociateKmsKeyRequestDescriptor = $convert.base64Decode(
    'ChlEaXNhc3NvY2lhdGVLbXNLZXlSZXF1ZXN0EiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG'
    '9nZ3JvdXBuYW1lEjIKEnJlc291cmNlaWRlbnRpZmllchjP+cWTASABKAlSEnJlc291cmNlaWRl'
    'bnRpZmllcg==');

@$core.Deprecated(
    'Use disassociateSourceFromS3TableIntegrationRequestDescriptor instead')
const DisassociateSourceFromS3TableIntegrationRequest$json = {
  '1': 'DisassociateSourceFromS3TableIntegrationRequest',
  '2': [
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `DisassociateSourceFromS3TableIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    disassociateSourceFromS3TableIntegrationRequestDescriptor =
    $convert.base64Decode(
        'Ci9EaXNhc3NvY2lhdGVTb3VyY2VGcm9tUzNUYWJsZUludGVncmF0aW9uUmVxdWVzdBIhCgppZG'
        'VudGlmaWVyGL/QlkUgASgJUgppZGVudGlmaWVy');

@$core.Deprecated(
    'Use disassociateSourceFromS3TableIntegrationResponseDescriptor instead')
const DisassociateSourceFromS3TableIntegrationResponse$json = {
  '1': 'DisassociateSourceFromS3TableIntegrationResponse',
  '2': [
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `DisassociateSourceFromS3TableIntegrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    disassociateSourceFromS3TableIntegrationResponseDescriptor =
    $convert.base64Decode(
        'CjBEaXNhc3NvY2lhdGVTb3VyY2VGcm9tUzNUYWJsZUludGVncmF0aW9uUmVzcG9uc2USIQoKaW'
        'RlbnRpZmllchi/0JZFIAEoCVIKaWRlbnRpZmllcg==');

@$core.Deprecated('Use entityDescriptor instead')
const Entity$json = {
  '1': 'Entity',
  '2': [
    {
      '1': 'attributes',
      '3': 33545109,
      '4': 3,
      '5': 11,
      '6': '.logs.Entity.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'keyattributes',
      '3': 299230926,
      '4': 3,
      '5': 11,
      '6': '.logs.Entity.KeyattributesEntry',
      '10': 'keyattributes'
    },
  ],
  '3': [Entity_AttributesEntry$json, Entity_KeyattributesEntry$json],
};

@$core.Deprecated('Use entityDescriptor instead')
const Entity_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use entityDescriptor instead')
const Entity_KeyattributesEntry$json = {
  '1': 'KeyattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Entity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityDescriptor = $convert.base64Decode(
    'CgZFbnRpdHkSPwoKYXR0cmlidXRlcxiVt/8PIAMoCzIcLmxvZ3MuRW50aXR5LkF0dHJpYnV0ZX'
    'NFbnRyeVIKYXR0cmlidXRlcxJJCg1rZXlhdHRyaWJ1dGVzGM7N144BIAMoCzIfLmxvZ3MuRW50'
    'aXR5LktleWF0dHJpYnV0ZXNFbnRyeVINa2V5YXR0cmlidXRlcxo9Cg9BdHRyaWJ1dGVzRW50cn'
    'kSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpAChJLZXlhdHRy'
    'aWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ'
    '==');

@$core.Deprecated('Use exportTaskDescriptor instead')
const ExportTask$json = {
  '1': 'ExportTask',
  '2': [
    {'1': 'destination', '3': 316564672, '4': 1, '5': 9, '10': 'destination'},
    {
      '1': 'destinationprefix',
      '3': 172629044,
      '4': 1,
      '5': 9,
      '10': 'destinationprefix'
    },
    {
      '1': 'executioninfo',
      '3': 110255734,
      '4': 1,
      '5': 11,
      '6': '.logs.ExportTaskExecutionInfo',
      '10': 'executioninfo'
    },
    {'1': 'from', '3': 365789302, '4': 1, '5': 3, '10': 'from'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.ExportTaskStatus',
      '10': 'status'
    },
    {'1': 'taskid', '3': 216769858, '4': 1, '5': 9, '10': 'taskid'},
    {'1': 'taskname', '3': 82438536, '4': 1, '5': 9, '10': 'taskname'},
    {'1': 'to', '3': 38094885, '4': 1, '5': 3, '10': 'to'},
  ],
};

/// Descriptor for `ExportTask`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportTaskDescriptor = $convert.base64Decode(
    'CgpFeHBvcnRUYXNrEiQKC2Rlc3RpbmF0aW9uGMDJ+ZYBIAEoCVILZGVzdGluYXRpb24SLwoRZG'
    'VzdGluYXRpb25wcmVmaXgYtLioUiABKAlSEWRlc3RpbmF0aW9ucHJlZml4EkYKDWV4ZWN1dGlv'
    'bmluZm8Y9rzJNCABKAsyHS5sb2dzLkV4cG9ydFRhc2tFeGVjdXRpb25JbmZvUg1leGVjdXRpb2'
    '5pbmZvEhYKBGZyb20Y9oC2rgEgASgDUgRmcm9tEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIM'
    'bG9nZ3JvdXBuYW1lEjIKBnN0YXR1cxjw763SASABKAsyFi5sb2dzLkV4cG9ydFRhc2tTdGF0dX'
    'NSBnN0YXR1cxIZCgZ0YXNraWQYwsquZyABKAlSBnRhc2tpZBIdCgh0YXNrbmFtZRiI06cnIAEo'
    'CVIIdGFza25hbWUSEQoCdG8YpZCVEiABKANSAnRv');

@$core.Deprecated('Use exportTaskExecutionInfoDescriptor instead')
const ExportTaskExecutionInfo$json = {
  '1': 'ExportTaskExecutionInfo',
  '2': [
    {
      '1': 'completiontime',
      '3': 45320599,
      '4': 1,
      '5': 3,
      '10': 'completiontime'
    },
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
  ],
};

/// Descriptor for `ExportTaskExecutionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportTaskExecutionInfoDescriptor = $convert.base64Decode(
    'ChdFeHBvcnRUYXNrRXhlY3V0aW9uSW5mbxIpCg5jb21wbGV0aW9udGltZRiXk84VIAEoA1IOY2'
    '9tcGxldGlvbnRpbWUSJQoMY3JlYXRpb250aW1lGIbFxBkgASgDUgxjcmVhdGlvbnRpbWU=');

@$core.Deprecated('Use exportTaskStatusDescriptor instead')
const ExportTaskStatus$json = {
  '1': 'ExportTaskStatus',
  '2': [
    {
      '1': 'code',
      '3': 422669557,
      '4': 1,
      '5': 14,
      '6': '.logs.ExportTaskStatusCode',
      '10': 'code'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExportTaskStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportTaskStatusDescriptor = $convert.base64Decode(
    'ChBFeHBvcnRUYXNrU3RhdHVzEjIKBGNvZGUY9dnFyQEgASgOMhoubG9ncy5FeHBvcnRUYXNrU3'
    'RhdHVzQ29kZVIEY29kZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use fieldIndexDescriptor instead')
const FieldIndex$json = {
  '1': 'FieldIndex',
  '2': [
    {
      '1': 'fieldindexname',
      '3': 52826287,
      '4': 1,
      '5': 9,
      '10': 'fieldindexname'
    },
    {
      '1': 'firsteventtime',
      '3': 236445651,
      '4': 1,
      '5': 3,
      '10': 'firsteventtime'
    },
    {
      '1': 'lasteventtime',
      '3': 54472587,
      '4': 1,
      '5': 3,
      '10': 'lasteventtime'
    },
    {'1': 'lastscantime', '3': 464250390, '4': 1, '5': 3, '10': 'lastscantime'},
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.logs.IndexType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `FieldIndex`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldIndexDescriptor = $convert.base64Decode(
    'CgpGaWVsZEluZGV4EikKDmZpZWxkaW5kZXhuYW1lGK+hmBkgASgJUg5maWVsZGluZGV4bmFtZR'
    'IpCg5maXJzdGV2ZW50dGltZRjTv99wIAEoA1IOZmlyc3RldmVudHRpbWUSJwoNbGFzdGV2ZW50'
    'dGltZRiL3/wZIAEoA1INbGFzdGV2ZW50dGltZRImCgxsYXN0c2NhbnRpbWUYlsyv3QEgASgDUg'
    'xsYXN0c2NhbnRpbWUSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIAEoCVISbG9nZ3JvdXBp'
    'ZGVudGlmaWVyEicKBHR5cGUYzuKfiQEgASgOMg8ubG9ncy5JbmRleFR5cGVSBHR5cGU=');

@$core.Deprecated('Use fieldsDataDescriptor instead')
const FieldsData$json = {
  '1': 'FieldsData',
  '2': [
    {'1': 'data', '3': 410182310, '4': 1, '5': 12, '10': 'data'},
  ],
};

/// Descriptor for `FieldsData`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fieldsDataDescriptor =
    $convert.base64Decode('CgpGaWVsZHNEYXRhEhYKBGRhdGEYpsXLwwEgASgMUgRkYXRh');

@$core.Deprecated('Use filterLogEventsRequestDescriptor instead')
const FilterLogEventsRequest$json = {
  '1': 'FilterLogEventsRequest',
  '2': [
    {'1': 'endtime', '3': 329679852, '4': 1, '5': 3, '10': 'endtime'},
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'interleaved', '3': 405276841, '4': 1, '5': 8, '10': 'interleaved'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'logstreamnameprefix',
      '3': 24068871,
      '4': 1,
      '5': 9,
      '10': 'logstreamnameprefix'
    },
    {
      '1': 'logstreamnames',
      '3': 178825732,
      '4': 3,
      '5': 9,
      '10': 'logstreamnames'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'starttime', '3': 178154767, '4': 1, '5': 3, '10': 'starttime'},
    {'1': 'unmask', '3': 363382499, '4': 1, '5': 8, '10': 'unmask'},
  ],
};

/// Descriptor for `FilterLogEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterLogEventsRequestDescriptor = $convert.base64Decode(
    'ChZGaWx0ZXJMb2dFdmVudHNSZXF1ZXN0EhwKB2VuZHRpbWUY7IeanQEgASgDUgdlbmR0aW1lEi'
    'cKDWZpbHRlcnBhdHRlcm4YmIeKRSABKAlSDWZpbHRlcnBhdHRlcm4SJAoLaW50ZXJsZWF2ZWQY'
    'qZGgwQEgASgIUgtpbnRlcmxlYXZlZBIYCgVsaW1pdBi1suuWASABKAVSBWxpbWl0EjIKEmxvZ2'
    'dyb3VwaWRlbnRpZmllchjcz6XfASABKAlSEmxvZ2dyb3VwaWRlbnRpZmllchIlCgxsb2dncm91'
    'cG5hbWUYzL/ULSABKAlSDGxvZ2dyb3VwbmFtZRIzChNsb2dzdHJlYW1uYW1lcHJlZml4GIeGvQ'
    'sgASgJUhNsb2dzdHJlYW1uYW1lcHJlZml4EikKDmxvZ3N0cmVhbW5hbWVzGITUolUgAygJUg5s'
    'b2dzdHJlYW1uYW1lcxIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tlbhIfCglzdGFydH'
    'RpbWUYj9r5VCABKANSCXN0YXJ0dGltZRIaCgZ1bm1hc2sY442jrQEgASgIUgZ1bm1hc2s=');

@$core.Deprecated('Use filterLogEventsResponseDescriptor instead')
const FilterLogEventsResponse$json = {
  '1': 'FilterLogEventsResponse',
  '2': [
    {
      '1': 'events',
      '3': 316203909,
      '4': 3,
      '5': 11,
      '6': '.logs.FilteredLogEvent',
      '10': 'events'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'searchedlogstreams',
      '3': 67617196,
      '4': 3,
      '5': 11,
      '6': '.logs.SearchedLogStream',
      '10': 'searchedlogstreams'
    },
  ],
};

/// Descriptor for `FilterLogEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterLogEventsResponseDescriptor = $convert.base64Decode(
    'ChdGaWx0ZXJMb2dFdmVudHNSZXNwb25zZRIyCgZldmVudHMYhcfjlgEgAygLMhYubG9ncy5GaW'
    'x0ZXJlZExvZ0V2ZW50UgZldmVudHMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4S'
    'SgoSc2VhcmNoZWRsb2dzdHJlYW1zGKyDnyAgAygLMhcubG9ncy5TZWFyY2hlZExvZ1N0cmVhbV'
    'ISc2VhcmNoZWRsb2dzdHJlYW1z');

@$core.Deprecated('Use filteredLogEventDescriptor instead')
const FilteredLogEvent$json = {
  '1': 'FilteredLogEvent',
  '2': [
    {'1': 'eventid', '3': 255267571, '4': 1, '5': 9, '10': 'eventid'},
    {
      '1': 'ingestiontime',
      '3': 179367957,
      '4': 1,
      '5': 3,
      '10': 'ingestiontime'
    },
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'timestamp', '3': 310629668, '4': 1, '5': 3, '10': 'timestamp'},
  ],
};

/// Descriptor for `FilteredLogEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filteredLogEventDescriptor = $convert.base64Decode(
    'ChBGaWx0ZXJlZExvZ0V2ZW50EhsKB2V2ZW50aWQY86XceSABKAlSB2V2ZW50aWQSJwoNaW5nZX'
    'N0aW9udGltZRiV4MNVIAEoA1INaW5nZXN0aW9udGltZRIoCg1sb2dzdHJlYW1uYW1lGKP37tAB'
    'IAEoCVINbG9nc3RyZWFtbmFtZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdlEiAKCXRpbW'
    'VzdGFtcBikqo+UASABKANSCXRpbWVzdGFtcA==');

@$core.Deprecated('Use getDataProtectionPolicyRequestDescriptor instead')
const GetDataProtectionPolicyRequest$json = {
  '1': 'GetDataProtectionPolicyRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
  ],
};

/// Descriptor for `GetDataProtectionPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataProtectionPolicyRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXREYXRhUHJvdGVjdGlvblBvbGljeVJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGN'
        'zPpd8BIAEoCVISbG9nZ3JvdXBpZGVudGlmaWVy');

@$core.Deprecated('Use getDataProtectionPolicyResponseDescriptor instead')
const GetDataProtectionPolicyResponse$json = {
  '1': 'GetDataProtectionPolicyResponse',
  '2': [
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
  ],
};

/// Descriptor for `GetDataProtectionPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataProtectionPolicyResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXREYXRhUHJvdGVjdGlvblBvbGljeVJlc3BvbnNlEiwKD2xhc3R1cGRhdGVkdGltZRj2t5'
        'W5ASABKANSD2xhc3R1cGRhdGVkdGltZRIyChJsb2dncm91cGlkZW50aWZpZXIY3M+l3wEgASgJ'
        'UhJsb2dncm91cGlkZW50aWZpZXISKQoOcG9saWN5ZG9jdW1lbnQY6+n2VCABKAlSDnBvbGljeW'
        'RvY3VtZW50');

@$core.Deprecated('Use getDeliveryDestinationPolicyRequestDescriptor instead')
const GetDeliveryDestinationPolicyRequest$json = {
  '1': 'GetDeliveryDestinationPolicyRequest',
  '2': [
    {
      '1': 'deliverydestinationname',
      '3': 331358541,
      '4': 1,
      '5': 9,
      '10': 'deliverydestinationname'
    },
  ],
};

/// Descriptor for `GetDeliveryDestinationPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliveryDestinationPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiNHZXREZWxpdmVyeURlc3RpbmF0aW9uUG9saWN5UmVxdWVzdBI8ChdkZWxpdmVyeWRlc3Rpbm'
        'F0aW9ubmFtZRjNwoCeASABKAlSF2RlbGl2ZXJ5ZGVzdGluYXRpb25uYW1l');

@$core.Deprecated('Use getDeliveryDestinationPolicyResponseDescriptor instead')
const GetDeliveryDestinationPolicyResponse$json = {
  '1': 'GetDeliveryDestinationPolicyResponse',
  '2': [
    {
      '1': 'policy',
      '3': 247528064,
      '4': 1,
      '5': 11,
      '6': '.logs.Policy',
      '10': 'policy'
    },
  ],
};

/// Descriptor for `GetDeliveryDestinationPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliveryDestinationPolicyResponseDescriptor =
    $convert.base64Decode(
        'CiRHZXREZWxpdmVyeURlc3RpbmF0aW9uUG9saWN5UmVzcG9uc2USJwoGcG9saWN5GID1g3YgAS'
        'gLMgwubG9ncy5Qb2xpY3lSBnBvbGljeQ==');

@$core.Deprecated('Use getDeliveryDestinationRequestDescriptor instead')
const GetDeliveryDestinationRequest$json = {
  '1': 'GetDeliveryDestinationRequest',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetDeliveryDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliveryDestinationRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXREZWxpdmVyeURlc3RpbmF0aW9uUmVxdWVzdBIVCgRuYW1lGOf75mkgASgJUgRuYW1l');

@$core.Deprecated('Use getDeliveryDestinationResponseDescriptor instead')
const GetDeliveryDestinationResponse$json = {
  '1': 'GetDeliveryDestinationResponse',
  '2': [
    {
      '1': 'deliverydestination',
      '3': 103332720,
      '4': 1,
      '5': 11,
      '6': '.logs.DeliveryDestination',
      '10': 'deliverydestination'
    },
  ],
};

/// Descriptor for `GetDeliveryDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliveryDestinationResponseDescriptor =
    $convert.base64Decode(
        'Ch5HZXREZWxpdmVyeURlc3RpbmF0aW9uUmVzcG9uc2USTgoTZGVsaXZlcnlkZXN0aW5hdGlvbh'
        'jw9qIxIAEoCzIZLmxvZ3MuRGVsaXZlcnlEZXN0aW5hdGlvblITZGVsaXZlcnlkZXN0aW5hdGlv'
        'bg==');

@$core.Deprecated('Use getDeliveryRequestDescriptor instead')
const GetDeliveryRequest$json = {
  '1': 'GetDeliveryRequest',
  '2': [
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetDeliveryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliveryRequestDescriptor = $convert
    .base64Decode('ChJHZXREZWxpdmVyeVJlcXVlc3QSEgoCaWQY4dXhuQEgASgJUgJpZA==');

@$core.Deprecated('Use getDeliveryResponseDescriptor instead')
const GetDeliveryResponse$json = {
  '1': 'GetDeliveryResponse',
  '2': [
    {
      '1': 'delivery',
      '3': 432355926,
      '4': 1,
      '5': 11,
      '6': '.logs.Delivery',
      '10': 'delivery'
    },
  ],
};

/// Descriptor for `GetDeliveryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliveryResponseDescriptor = $convert.base64Decode(
    'ChNHZXREZWxpdmVyeVJlc3BvbnNlEi4KCGRlbGl2ZXJ5GNb0lM4BIAEoCzIOLmxvZ3MuRGVsaX'
    'ZlcnlSCGRlbGl2ZXJ5');

@$core.Deprecated('Use getDeliverySourceRequestDescriptor instead')
const GetDeliverySourceRequest$json = {
  '1': 'GetDeliverySourceRequest',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetDeliverySourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliverySourceRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXREZWxpdmVyeVNvdXJjZVJlcXVlc3QSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZQ==');

@$core.Deprecated('Use getDeliverySourceResponseDescriptor instead')
const GetDeliverySourceResponse$json = {
  '1': 'GetDeliverySourceResponse',
  '2': [
    {
      '1': 'deliverysource',
      '3': 1553897,
      '4': 1,
      '5': 11,
      '6': '.logs.DeliverySource',
      '10': 'deliverysource'
    },
  ],
};

/// Descriptor for `GetDeliverySourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliverySourceResponseDescriptor =
    $convert.base64Decode(
        'ChlHZXREZWxpdmVyeVNvdXJjZVJlc3BvbnNlEj4KDmRlbGl2ZXJ5c291cmNlGOnrXiABKAsyFC'
        '5sb2dzLkRlbGl2ZXJ5U291cmNlUg5kZWxpdmVyeXNvdXJjZQ==');

@$core.Deprecated('Use getIntegrationRequestDescriptor instead')
const GetIntegrationRequest$json = {
  '1': 'GetIntegrationRequest',
  '2': [
    {
      '1': 'integrationname',
      '3': 183183535,
      '4': 1,
      '5': 9,
      '10': 'integrationname'
    },
  ],
};

/// Descriptor for `GetIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIntegrationRequestDescriptor = $convert.base64Decode(
    'ChVHZXRJbnRlZ3JhdGlvblJlcXVlc3QSKwoPaW50ZWdyYXRpb25uYW1lGK/RrFcgASgJUg9pbn'
    'RlZ3JhdGlvbm5hbWU=');

@$core.Deprecated('Use getIntegrationResponseDescriptor instead')
const GetIntegrationResponse$json = {
  '1': 'GetIntegrationResponse',
  '2': [
    {
      '1': 'integrationdetails',
      '3': 475519038,
      '4': 1,
      '5': 11,
      '6': '.logs.IntegrationDetails',
      '10': 'integrationdetails'
    },
    {
      '1': 'integrationname',
      '3': 183183535,
      '4': 1,
      '5': 9,
      '10': 'integrationname'
    },
    {
      '1': 'integrationstatus',
      '3': 12110360,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationStatus',
      '10': 'integrationstatus'
    },
    {
      '1': 'integrationtype',
      '3': 307087270,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationType',
      '10': 'integrationtype'
    },
  ],
};

/// Descriptor for `GetIntegrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIntegrationResponseDescriptor = $convert.base64Decode(
    'ChZHZXRJbnRlZ3JhdGlvblJlc3BvbnNlEkwKEmludGVncmF0aW9uZGV0YWlscxi+sN/iASABKA'
    'syGC5sb2dzLkludGVncmF0aW9uRGV0YWlsc1ISaW50ZWdyYXRpb25kZXRhaWxzEisKD2ludGVn'
    'cmF0aW9ubmFtZRiv0axXIAEoCVIPaW50ZWdyYXRpb25uYW1lEkgKEWludGVncmF0aW9uc3RhdH'
    'VzGJiU4wUgASgOMhcubG9ncy5JbnRlZ3JhdGlvblN0YXR1c1IRaW50ZWdyYXRpb25zdGF0dXMS'
    'QwoPaW50ZWdyYXRpb250eXBlGKaPt5IBIAEoDjIVLmxvZ3MuSW50ZWdyYXRpb25UeXBlUg9pbn'
    'RlZ3JhdGlvbnR5cGU=');

@$core.Deprecated('Use getLogAnomalyDetectorRequestDescriptor instead')
const GetLogAnomalyDetectorRequest$json = {
  '1': 'GetLogAnomalyDetectorRequest',
  '2': [
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
  ],
};

/// Descriptor for `GetLogAnomalyDetectorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogAnomalyDetectorRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRMb2dBbm9tYWx5RGV0ZWN0b3JSZXF1ZXN0EjIKEmFub21hbHlkZXRlY3RvcmFybhisz/'
        'PUASABKAlSEmFub21hbHlkZXRlY3RvcmFybg==');

@$core.Deprecated('Use getLogAnomalyDetectorResponseDescriptor instead')
const GetLogAnomalyDetectorResponse$json = {
  '1': 'GetLogAnomalyDetectorResponse',
  '2': [
    {
      '1': 'anomalydetectorstatus',
      '3': 458778431,
      '4': 1,
      '5': 14,
      '6': '.logs.AnomalyDetectorStatus',
      '10': 'anomalydetectorstatus'
    },
    {
      '1': 'anomalyvisibilitytime',
      '3': 287987260,
      '4': 1,
      '5': 3,
      '10': 'anomalyvisibilitytime'
    },
    {
      '1': 'creationtimestamp',
      '3': 206588645,
      '4': 1,
      '5': 3,
      '10': 'creationtimestamp'
    },
    {'1': 'detectorname', '3': 114651981, '4': 1, '5': 9, '10': 'detectorname'},
    {
      '1': 'evaluationfrequency',
      '3': 533700360,
      '4': 1,
      '5': 14,
      '6': '.logs.EvaluationFrequency',
      '10': 'evaluationfrequency'
    },
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'kmskeyid', '3': 510698477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'lastmodifiedtimestamp',
      '3': 40019279,
      '4': 1,
      '5': 3,
      '10': 'lastmodifiedtimestamp'
    },
    {
      '1': 'loggrouparnlist',
      '3': 374867736,
      '4': 3,
      '5': 9,
      '10': 'loggrouparnlist'
    },
  ],
};

/// Descriptor for `GetLogAnomalyDetectorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogAnomalyDetectorResponseDescriptor = $convert.base64Decode(
    'Ch1HZXRMb2dBbm9tYWx5RGV0ZWN0b3JSZXNwb25zZRJVChVhbm9tYWx5ZGV0ZWN0b3JzdGF0dX'
    'MYv87h2gEgASgOMhsubG9ncy5Bbm9tYWx5RGV0ZWN0b3JTdGF0dXNSFWFub21hbHlkZXRlY3Rv'
    'cnN0YXR1cxI4ChVhbm9tYWx5dmlzaWJpbGl0eXRpbWUYvKypiQEgASgDUhVhbm9tYWx5dmlzaW'
    'JpbGl0eXRpbWUSLwoRY3JlYXRpb250aW1lc3RhbXAY5ZXBYiABKANSEWNyZWF0aW9udGltZXN0'
    'YW1wEiUKDGRldGVjdG9ybmFtZRjN5tU2IAEoCVIMZGV0ZWN0b3JuYW1lEk8KE2V2YWx1YXRpb2'
    '5mcmVxdWVuY3kYiL6+/gEgASgOMhkubG9ncy5FdmFsdWF0aW9uRnJlcXVlbmN5UhNldmFsdWF0'
    'aW9uZnJlcXVlbmN5EicKDWZpbHRlcnBhdHRlcm4YmIeKRSABKAlSDWZpbHRlcnBhdHRlcm4SHg'
    'oIa21za2V5aWQY7cfC8wEgASgJUghrbXNrZXlpZBI3ChVsYXN0bW9kaWZpZWR0aW1lc3RhbXAY'
    'z8qKEyABKANSFWxhc3Rtb2RpZmllZHRpbWVzdGFtcBIsCg9sb2dncm91cGFybmxpc3QYmI7gsg'
    'EgAygJUg9sb2dncm91cGFybmxpc3Q=');

@$core.Deprecated('Use getLogEventsRequestDescriptor instead')
const GetLogEventsRequest$json = {
  '1': 'GetLogEventsRequest',
  '2': [
    {'1': 'endtime', '3': 329679852, '4': 1, '5': 3, '10': 'endtime'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'startfromhead',
      '3': 472305066,
      '4': 1,
      '5': 8,
      '10': 'startfromhead'
    },
    {'1': 'starttime', '3': 178154767, '4': 1, '5': 3, '10': 'starttime'},
    {'1': 'unmask', '3': 363382499, '4': 1, '5': 8, '10': 'unmask'},
  ],
};

/// Descriptor for `GetLogEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogEventsRequestDescriptor = $convert.base64Decode(
    'ChNHZXRMb2dFdmVudHNSZXF1ZXN0EhwKB2VuZHRpbWUY7IeanQEgASgDUgdlbmR0aW1lEhgKBW'
    'xpbWl0GLWy65YBIAEoBVIFbGltaXQSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIAEoCVIS'
    'bG9nZ3JvdXBpZGVudGlmaWVyEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW'
    '1lEigKDWxvZ3N0cmVhbW5hbWUYo/fu0AEgASgJUg1sb2dzdHJlYW1uYW1lEh8KCW5leHR0b2tl'
    'bhie8503IAEoCVIJbmV4dHRva2VuEigKDXN0YXJ0ZnJvbWhlYWQYqpub4QEgASgIUg1zdGFydG'
    'Zyb21oZWFkEh8KCXN0YXJ0dGltZRiP2vlUIAEoA1IJc3RhcnR0aW1lEhoKBnVubWFzaxjjjaOt'
    'ASABKAhSBnVubWFzaw==');

@$core.Deprecated('Use getLogEventsResponseDescriptor instead')
const GetLogEventsResponse$json = {
  '1': 'GetLogEventsResponse',
  '2': [
    {
      '1': 'events',
      '3': 316203909,
      '4': 3,
      '5': 11,
      '6': '.logs.OutputLogEvent',
      '10': 'events'
    },
    {
      '1': 'nextbackwardtoken',
      '3': 265738167,
      '4': 1,
      '5': 9,
      '10': 'nextbackwardtoken'
    },
    {
      '1': 'nextforwardtoken',
      '3': 262037315,
      '4': 1,
      '5': 9,
      '10': 'nextforwardtoken'
    },
  ],
};

/// Descriptor for `GetLogEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogEventsResponseDescriptor = $convert.base64Decode(
    'ChRHZXRMb2dFdmVudHNSZXNwb25zZRIwCgZldmVudHMYhcfjlgEgAygLMhQubG9ncy5PdXRwdX'
    'RMb2dFdmVudFIGZXZlbnRzEi8KEW5leHRiYWNrd2FyZHRva2VuGLev234gASgJUhFuZXh0YmFj'
    'a3dhcmR0b2tlbhItChBuZXh0Zm9yd2FyZHRva2VuGMO++XwgASgJUhBuZXh0Zm9yd2FyZHRva2'
    'Vu');

@$core.Deprecated('Use getLogFieldsRequestDescriptor instead')
const GetLogFieldsRequest$json = {
  '1': 'GetLogFieldsRequest',
  '2': [
    {
      '1': 'datasourcename',
      '3': 231923996,
      '4': 1,
      '5': 9,
      '10': 'datasourcename'
    },
    {
      '1': 'datasourcetype',
      '3': 635377,
      '4': 1,
      '5': 9,
      '10': 'datasourcetype'
    },
  ],
};

/// Descriptor for `GetLogFieldsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogFieldsRequestDescriptor = $convert.base64Decode(
    'ChNHZXRMb2dGaWVsZHNSZXF1ZXN0EikKDmRhdGFzb3VyY2VuYW1lGJzCy24gASgJUg5kYXRhc2'
    '91cmNlbmFtZRIoCg5kYXRhc291cmNldHlwZRjx4yYgASgJUg5kYXRhc291cmNldHlwZQ==');

@$core.Deprecated('Use getLogFieldsResponseDescriptor instead')
const GetLogFieldsResponse$json = {
  '1': 'GetLogFieldsResponse',
  '2': [
    {
      '1': 'logfields',
      '3': 229969497,
      '4': 3,
      '5': 11,
      '6': '.logs.LogFieldsListItem',
      '10': 'logfields'
    },
  ],
};

/// Descriptor for `GetLogFieldsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogFieldsResponseDescriptor = $convert.base64Decode(
    'ChRHZXRMb2dGaWVsZHNSZXNwb25zZRI4Cglsb2dmaWVsZHMY2ZzUbSADKAsyFy5sb2dzLkxvZ0'
    'ZpZWxkc0xpc3RJdGVtUglsb2dmaWVsZHM=');

@$core.Deprecated('Use getLogGroupFieldsRequestDescriptor instead')
const GetLogGroupFieldsRequest$json = {
  '1': 'GetLogGroupFieldsRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'time', '3': 490511333, '4': 1, '5': 3, '10': 'time'},
  ],
};

/// Descriptor for `GetLogGroupFieldsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogGroupFieldsRequestDescriptor = $convert.base64Decode(
    'ChhHZXRMb2dHcm91cEZpZWxkc1JlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIA'
    'EoCVISbG9nZ3JvdXBpZGVudGlmaWVyEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3Jv'
    'dXBuYW1lEhYKBHRpbWUY5bfy6QEgASgDUgR0aW1l');

@$core.Deprecated('Use getLogGroupFieldsResponseDescriptor instead')
const GetLogGroupFieldsResponse$json = {
  '1': 'GetLogGroupFieldsResponse',
  '2': [
    {
      '1': 'loggroupfields',
      '3': 79081558,
      '4': 3,
      '5': 11,
      '6': '.logs.LogGroupField',
      '10': 'loggroupfields'
    },
  ],
};

/// Descriptor for `GetLogGroupFieldsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogGroupFieldsResponseDescriptor =
    $convert.base64Decode(
        'ChlHZXRMb2dHcm91cEZpZWxkc1Jlc3BvbnNlEj4KDmxvZ2dyb3VwZmllbGRzGNbg2iUgAygLMh'
        'MubG9ncy5Mb2dHcm91cEZpZWxkUg5sb2dncm91cGZpZWxkcw==');

@$core.Deprecated('Use getLogObjectRequestDescriptor instead')
const GetLogObjectRequest$json = {
  '1': 'GetLogObjectRequest',
  '2': [
    {
      '1': 'logobjectpointer',
      '3': 460396600,
      '4': 1,
      '5': 9,
      '10': 'logobjectpointer'
    },
    {'1': 'unmask', '3': 363382499, '4': 1, '5': 8, '10': 'unmask'},
  ],
};

/// Descriptor for `GetLogObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogObjectRequestDescriptor = $convert.base64Decode(
    'ChNHZXRMb2dPYmplY3RSZXF1ZXN0Ei4KEGxvZ29iamVjdHBvaW50ZXIYuLDE2wEgASgJUhBsb2'
    'dvYmplY3Rwb2ludGVyEhoKBnVubWFzaxjjjaOtASABKAhSBnVubWFzaw==');

@$core.Deprecated('Use getLogObjectResponseDescriptor instead')
const GetLogObjectResponse$json = {
  '1': 'GetLogObjectResponse',
  '2': [
    {
      '1': 'fieldstream',
      '3': 531106040,
      '4': 1,
      '5': 11,
      '6': '.logs.GetLogObjectResponseStream',
      '10': 'fieldstream'
    },
  ],
};

/// Descriptor for `GetLogObjectResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogObjectResponseDescriptor = $convert.base64Decode(
    'ChRHZXRMb2dPYmplY3RSZXNwb25zZRJGCgtmaWVsZHN0cmVhbRj4kaD9ASABKAsyIC5sb2dzLk'
    'dldExvZ09iamVjdFJlc3BvbnNlU3RyZWFtUgtmaWVsZHN0cmVhbQ==');

@$core.Deprecated('Use getLogObjectResponseStreamDescriptor instead')
const GetLogObjectResponseStream$json = {
  '1': 'GetLogObjectResponseStream',
  '2': [
    {
      '1': 'internalstreamingexception',
      '3': 431472842,
      '4': 1,
      '5': 11,
      '6': '.logs.InternalStreamingException',
      '10': 'internalstreamingexception'
    },
    {
      '1': 'fields',
      '3': 104883581,
      '4': 1,
      '5': 11,
      '6': '.logs.FieldsData',
      '10': 'fields'
    },
  ],
};

/// Descriptor for `GetLogObjectResponseStream`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogObjectResponseStreamDescriptor = $convert.base64Decode(
    'ChpHZXRMb2dPYmplY3RSZXNwb25zZVN0cmVhbRJkChppbnRlcm5hbHN0cmVhbWluZ2V4Y2VwdG'
    'lvbhjKgd/NASABKAsyIC5sb2dzLkludGVybmFsU3RyZWFtaW5nRXhjZXB0aW9uUhppbnRlcm5h'
    'bHN0cmVhbWluZ2V4Y2VwdGlvbhIrCgZmaWVsZHMY/cqBMiABKAsyEC5sb2dzLkZpZWxkc0RhdG'
    'FSBmZpZWxkcw==');

@$core.Deprecated('Use getLogRecordRequestDescriptor instead')
const GetLogRecordRequest$json = {
  '1': 'GetLogRecordRequest',
  '2': [
    {
      '1': 'logrecordpointer',
      '3': 456941138,
      '4': 1,
      '5': 9,
      '10': 'logrecordpointer'
    },
    {'1': 'unmask', '3': 363382499, '4': 1, '5': 8, '10': 'unmask'},
  ],
};

/// Descriptor for `GetLogRecordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogRecordRequestDescriptor = $convert.base64Decode(
    'ChNHZXRMb2dSZWNvcmRSZXF1ZXN0Ei4KEGxvZ3JlY29yZHBvaW50ZXIY0rzx2QEgASgJUhBsb2'
    'dyZWNvcmRwb2ludGVyEhoKBnVubWFzaxjjjaOtASABKAhSBnVubWFzaw==');

@$core.Deprecated('Use getLogRecordResponseDescriptor instead')
const GetLogRecordResponse$json = {
  '1': 'GetLogRecordResponse',
  '2': [
    {
      '1': 'logrecord',
      '3': 456252785,
      '4': 3,
      '5': 11,
      '6': '.logs.GetLogRecordResponse.LogrecordEntry',
      '10': 'logrecord'
    },
  ],
  '3': [GetLogRecordResponse_LogrecordEntry$json],
};

@$core.Deprecated('Use getLogRecordResponseDescriptor instead')
const GetLogRecordResponse_LogrecordEntry$json = {
  '1': 'LogrecordEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetLogRecordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLogRecordResponseDescriptor = $convert.base64Decode(
    'ChRHZXRMb2dSZWNvcmRSZXNwb25zZRJLCglsb2dyZWNvcmQY8brH2QEgAygLMikubG9ncy5HZX'
    'RMb2dSZWNvcmRSZXNwb25zZS5Mb2dyZWNvcmRFbnRyeVIJbG9ncmVjb3JkGjwKDkxvZ3JlY29y'
    'ZEVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getQueryResultsRequestDescriptor instead')
const GetQueryResultsRequest$json = {
  '1': 'GetQueryResultsRequest',
  '2': [
    {'1': 'queryid', '3': 336975759, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `GetQueryResultsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryResultsRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRRdWVyeVJlc3VsdHNSZXF1ZXN0EhwKB3F1ZXJ5aWQYj6/XoAEgASgJUgdxdWVyeWlk');

@$core.Deprecated('Use getQueryResultsResponseDescriptor instead')
const GetQueryResultsResponse$json = {
  '1': 'GetQueryResultsResponse',
  '2': [
    {
      '1': 'encryptionkey',
      '3': 138098284,
      '4': 1,
      '5': 9,
      '10': 'encryptionkey'
    },
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {
      '1': 'results',
      '3': 206523126,
      '4': 3,
      '5': 11,
      '6': '.logs.ResultField',
      '10': 'results'
    },
    {
      '1': 'statistics',
      '3': 222129163,
      '4': 1,
      '5': 11,
      '6': '.logs.QueryStatistics',
      '10': 'statistics'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetQueryResultsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryResultsResponseDescriptor = $convert.base64Decode(
    'ChdHZXRRdWVyeVJlc3VsdHNSZXNwb25zZRInCg1lbmNyeXB0aW9ua2V5GOzs7EEgASgJUg1lbm'
    'NyeXB0aW9ua2V5Ej0KDXF1ZXJ5bGFuZ3VhZ2UY9q/3owEgASgOMhMubG9ncy5RdWVyeUxhbmd1'
    'YWdlUg1xdWVyeWxhbmd1YWdlEi4KB3Jlc3VsdHMY9pW9YiADKAsyES5sb2dzLlJlc3VsdEZpZW'
    'xkUgdyZXN1bHRzEjgKCnN0YXRpc3RpY3MYi9j1aSABKAsyFS5sb2dzLlF1ZXJ5U3RhdGlzdGlj'
    'c1IKc3RhdGlzdGljcxItCgZzdGF0dXMY8O+t0gEgASgOMhEubG9ncy5RdWVyeVN0YXR1c1IGc3'
    'RhdHVz');

@$core.Deprecated('Use getScheduledQueryHistoryRequestDescriptor instead')
const GetScheduledQueryHistoryRequest$json = {
  '1': 'GetScheduledQueryHistoryRequest',
  '2': [
    {'1': 'endtime', '3': 329679852, '4': 1, '5': 3, '10': 'endtime'},
    {
      '1': 'executionstatuses',
      '3': 457739688,
      '4': 3,
      '5': 14,
      '6': '.logs.ExecutionStatus',
      '10': 'executionstatuses'
    },
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'starttime', '3': 178154767, '4': 1, '5': 3, '10': 'starttime'},
  ],
};

/// Descriptor for `GetScheduledQueryHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduledQueryHistoryRequestDescriptor = $convert.base64Decode(
    'Ch9HZXRTY2hlZHVsZWRRdWVyeUhpc3RvcnlSZXF1ZXN0EhwKB2VuZHRpbWUY7IeanQEgASgDUg'
    'dlbmR0aW1lEkcKEWV4ZWN1dGlvbnN0YXR1c2VzGKibotoBIAMoDjIVLmxvZ3MuRXhlY3V0aW9u'
    'U3RhdHVzUhFleGVjdXRpb25zdGF0dXNlcxIhCgppZGVudGlmaWVyGL/QlkUgASgJUgppZGVudG'
    'lmaWVyEiIKCm1heHJlc3VsdHMY0tzn3QEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhie'
    '8503IAEoCVIJbmV4dHRva2VuEh8KCXN0YXJ0dGltZRiP2vlUIAEoA1IJc3RhcnR0aW1l');

@$core.Deprecated('Use getScheduledQueryHistoryResponseDescriptor instead')
const GetScheduledQueryHistoryResponse$json = {
  '1': 'GetScheduledQueryHistoryResponse',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'scheduledqueryarn',
      '3': 240292916,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
    {
      '1': 'triggerhistory',
      '3': 6841438,
      '4': 3,
      '5': 11,
      '6': '.logs.TriggerHistoryRecord',
      '10': 'triggerhistory'
    },
  ],
};

/// Descriptor for `GetScheduledQueryHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduledQueryHistoryResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXRTY2hlZHVsZWRRdWVyeUhpc3RvcnlSZXNwb25zZRIVCgRuYW1lGOf75mkgASgJUgRuYW'
        '1lEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2VuEi8KEXNjaGVkdWxlZHF1ZXJ5YXJu'
        'GLSoynIgASgJUhFzY2hlZHVsZWRxdWVyeWFybhJFCg50cmlnZ2VyaGlzdG9yeRjeyKEDIAMoCz'
        'IaLmxvZ3MuVHJpZ2dlckhpc3RvcnlSZWNvcmRSDnRyaWdnZXJoaXN0b3J5');

@$core.Deprecated('Use getScheduledQueryRequestDescriptor instead')
const GetScheduledQueryRequest$json = {
  '1': 'GetScheduledQueryRequest',
  '2': [
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
  ],
};

/// Descriptor for `GetScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduledQueryRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRTY2hlZHVsZWRRdWVyeVJlcXVlc3QSIQoKaWRlbnRpZmllchi/0JZFIAEoCVIKaWRlbn'
        'RpZmllcg==');

@$core.Deprecated('Use getScheduledQueryResponseDescriptor instead')
const GetScheduledQueryResponse$json = {
  '1': 'GetScheduledQueryResponse',
  '2': [
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'destinationconfiguration',
      '3': 145518016,
      '4': 1,
      '5': 11,
      '6': '.logs.DestinationConfiguration',
      '10': 'destinationconfiguration'
    },
    {
      '1': 'executionrolearn',
      '3': 230613553,
      '4': 1,
      '5': 9,
      '10': 'executionrolearn'
    },
    {
      '1': 'lastexecutionstatus',
      '3': 7800596,
      '4': 1,
      '5': 14,
      '6': '.logs.ExecutionStatus',
      '10': 'lastexecutionstatus'
    },
    {
      '1': 'lasttriggeredtime',
      '3': 397057656,
      '4': 1,
      '5': 3,
      '10': 'lasttriggeredtime'
    },
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'scheduleendtime',
      '3': 111645113,
      '4': 1,
      '5': 3,
      '10': 'scheduleendtime'
    },
    {
      '1': 'scheduleexpression',
      '3': 287794975,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'schedulestarttime',
      '3': 464194170,
      '4': 1,
      '5': 3,
      '10': 'schedulestarttime'
    },
    {
      '1': 'scheduledqueryarn',
      '3': 240292916,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
    {
      '1': 'starttimeoffset',
      '3': 308525274,
      '4': 1,
      '5': 3,
      '10': 'starttimeoffset'
    },
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryState',
      '10': 'state'
    },
    {'1': 'timezone', '3': 190615331, '4': 1, '5': 9, '10': 'timezone'},
  ],
};

/// Descriptor for `GetScheduledQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduledQueryResponseDescriptor = $convert.base64Decode(
    'ChlHZXRTY2hlZHVsZWRRdWVyeVJlc3BvbnNlEiUKDGNyZWF0aW9udGltZRiGxcQZIAEoA1IMY3'
    'JlYXRpb250aW1lEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SXQoYZGVz'
    'dGluYXRpb25jb25maWd1cmF0aW9uGMDbsUUgASgLMh4ubG9ncy5EZXN0aW5hdGlvbkNvbmZpZ3'
    'VyYXRpb25SGGRlc3RpbmF0aW9uY29uZmlndXJhdGlvbhItChBleGVjdXRpb25yb2xlYXJuGLHE'
    '+20gASgJUhBleGVjdXRpb25yb2xlYXJuEkoKE2xhc3RleGVjdXRpb25zdGF0dXMYlI7cAyABKA'
    '4yFS5sb2dzLkV4ZWN1dGlvblN0YXR1c1ITbGFzdGV4ZWN1dGlvbnN0YXR1cxIwChFsYXN0dHJp'
    'Z2dlcmVkdGltZRj4vKq9ASABKANSEWxhc3R0cmlnZ2VyZWR0aW1lEiwKD2xhc3R1cGRhdGVkdG'
    'ltZRj2t5W5ASABKANSD2xhc3R1cGRhdGVkdGltZRI0ChNsb2dncm91cGlkZW50aWZpZXJzGPna'
    'uMMBIAMoCVITbG9nZ3JvdXBpZGVudGlmaWVycxIVCgRuYW1lGOf75mkgASgJUgRuYW1lEj0KDX'
    'F1ZXJ5bGFuZ3VhZ2UY9q/3owEgASgOMhMubG9ncy5RdWVyeUxhbmd1YWdlUg1xdWVyeWxhbmd1'
    'YWdlEiQKC3F1ZXJ5c3RyaW5nGIeBnfgBIAEoCVILcXVlcnlzdHJpbmcSKwoPc2NoZWR1bGVlbm'
    'R0aW1lGLmjnjUgASgDUg9zY2hlZHVsZWVuZHRpbWUSMgoSc2NoZWR1bGVleHByZXNzaW9uGJ/O'
    'nYkBIAEoCVISc2NoZWR1bGVleHByZXNzaW9uEjAKEXNjaGVkdWxlc3RhcnR0aW1lGPqUrN0BIA'
    'EoA1IRc2NoZWR1bGVzdGFydHRpbWUSLwoRc2NoZWR1bGVkcXVlcnlhcm4YtKjKciABKAlSEXNj'
    'aGVkdWxlZHF1ZXJ5YXJuEiwKD3N0YXJ0dGltZW9mZnNldBja8Y6TASABKANSD3N0YXJ0dGltZW'
    '9mZnNldBIzCgVzdGF0ZRj35cTBASABKA4yGS5sb2dzLlNjaGVkdWxlZFF1ZXJ5U3RhdGVSBXN0'
    'YXRlEh0KCHRpbWV6b25lGKOe8logASgJUgh0aW1lem9uZQ==');

@$core.Deprecated('Use getTransformerRequestDescriptor instead')
const GetTransformerRequest$json = {
  '1': 'GetTransformerRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
  ],
};

/// Descriptor for `GetTransformerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTransformerRequestDescriptor = $convert.base64Decode(
    'ChVHZXRUcmFuc2Zvcm1lclJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIAEoCV'
    'ISbG9nZ3JvdXBpZGVudGlmaWVy');

@$core.Deprecated('Use getTransformerResponseDescriptor instead')
const GetTransformerResponse$json = {
  '1': 'GetTransformerResponse',
  '2': [
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'lastmodifiedtime',
      '3': 374093504,
      '4': 1,
      '5': 3,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'transformerconfig',
      '3': 384836439,
      '4': 3,
      '5': 11,
      '6': '.logs.Processor',
      '10': 'transformerconfig'
    },
  ],
};

/// Descriptor for `GetTransformerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTransformerResponseDescriptor = $convert.base64Decode(
    'ChZHZXRUcmFuc2Zvcm1lclJlc3BvbnNlEiUKDGNyZWF0aW9udGltZRiGxcQZIAEoA1IMY3JlYX'
    'Rpb250aW1lEi4KEGxhc3Rtb2RpZmllZHRpbWUYwO2wsgEgASgDUhBsYXN0bW9kaWZpZWR0aW1l'
    'EjIKEmxvZ2dyb3VwaWRlbnRpZmllchjcz6XfASABKAlSEmxvZ2dyb3VwaWRlbnRpZmllchJBCh'
    'F0cmFuc2Zvcm1lcmNvbmZpZxjXxsC3ASADKAsyDy5sb2dzLlByb2Nlc3NvclIRdHJhbnNmb3Jt'
    'ZXJjb25maWc=');

@$core.Deprecated('Use grokDescriptor instead')
const Grok$json = {
  '1': 'Grok',
  '2': [
    {'1': 'match', '3': 505425815, '4': 1, '5': 9, '10': 'match'},
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `Grok`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grokDescriptor = $convert.base64Decode(
    'CgRHcm9rEhgKBW1hdGNoGJffgPEBIAEoCVIFbWF0Y2gSGgoGc291cmNlGNnTvN4BIAEoCVIGc2'
    '91cmNl');

@$core.Deprecated('Use groupingIdentifierDescriptor instead')
const GroupingIdentifier$json = {
  '1': 'GroupingIdentifier',
  '2': [
    {'1': 'key', '3': 135645293, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `GroupingIdentifier`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List groupingIdentifierDescriptor = $convert.base64Decode(
    'ChJHcm91cGluZ0lkZW50aWZpZXISEwoDa2V5GO2Q10AgASgJUgNrZXkSFwoFdmFsdWUYy6f7Ei'
    'ABKAlSBXZhbHVl');

@$core.Deprecated('Use importDescriptor instead')
const Import$json = {
  '1': 'Import',
  '2': [
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {'1': 'errormessage', '3': 136873289, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'importdestinationarn',
      '3': 94477180,
      '4': 1,
      '5': 9,
      '10': 'importdestinationarn'
    },
    {
      '1': 'importfilter',
      '3': 170561015,
      '4': 1,
      '5': 11,
      '6': '.logs.ImportFilter',
      '10': 'importfilter'
    },
    {'1': 'importid', '3': 513429114, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importsourcearn',
      '3': 161570329,
      '4': 1,
      '5': 9,
      '10': 'importsourcearn'
    },
    {
      '1': 'importstatistics',
      '3': 60366280,
      '4': 1,
      '5': 11,
      '6': '.logs.ImportStatistics',
      '10': 'importstatistics'
    },
    {
      '1': 'importstatus',
      '3': 31427999,
      '4': 1,
      '5': 14,
      '6': '.logs.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
  ],
};

/// Descriptor for `Import`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importDescriptor = $convert.base64Decode(
    'CgZJbXBvcnQSJQoMY3JlYXRpb250aW1lGIbFxBkgASgDUgxjcmVhdGlvbnRpbWUSJQoMZXJyb3'
    'JtZXNzYWdlGMmKokEgASgJUgxlcnJvcm1lc3NhZ2USNQoUaW1wb3J0ZGVzdGluYXRpb25hcm4Y'
    '/LaGLSABKAlSFGltcG9ydGRlc3RpbmF0aW9uYXJuEjkKDGltcG9ydGZpbHRlchj3m6pRIAEoCz'
    'ISLmxvZ3MuSW1wb3J0RmlsdGVyUgxpbXBvcnRmaWx0ZXISHgoIaW1wb3J0aWQY+pzp9AEgASgJ'
    'UghpbXBvcnRpZBIrCg9pbXBvcnRzb3VyY2Vhcm4YmbyFTSABKAlSD2ltcG9ydHNvdXJjZWFybh'
    'JFChBpbXBvcnRzdGF0aXN0aWNzGMi75BwgASgLMhYubG9ncy5JbXBvcnRTdGF0aXN0aWNzUhBp'
    'bXBvcnRzdGF0aXN0aWNzEjkKDGltcG9ydHN0YXR1cxifm/4OIAEoDjISLmxvZ3MuSW1wb3J0U3'
    'RhdHVzUgxpbXBvcnRzdGF0dXMSLAoPbGFzdHVwZGF0ZWR0aW1lGPa3lbkBIAEoA1IPbGFzdHVw'
    'ZGF0ZWR0aW1l');

@$core.Deprecated('Use importBatchDescriptor instead')
const ImportBatch$json = {
  '1': 'ImportBatch',
  '2': [
    {'1': 'batchid', '3': 194363415, '4': 1, '5': 9, '10': 'batchid'},
    {'1': 'errormessage', '3': 136873289, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.logs.ImportStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ImportBatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importBatchDescriptor = $convert.base64Decode(
    'CgtJbXBvcnRCYXRjaBIbCgdiYXRjaGlkGJeA11wgASgJUgdiYXRjaGlkEiUKDGVycm9ybWVzc2'
    'FnZRjJiqJBIAEoCVIMZXJyb3JtZXNzYWdlEi4KBnN0YXR1cxjw763SASABKA4yEi5sb2dzLklt'
    'cG9ydFN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use importFilterDescriptor instead')
const ImportFilter$json = {
  '1': 'ImportFilter',
  '2': [
    {'1': 'endeventtime', '3': 345758432, '4': 1, '5': 3, '10': 'endeventtime'},
    {
      '1': 'starteventtime',
      '3': 406085597,
      '4': 1,
      '5': 3,
      '10': 'starteventtime'
    },
  ],
};

/// Descriptor for `ImportFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importFilterDescriptor = $convert.base64Decode(
    'CgxJbXBvcnRGaWx0ZXISJgoMZW5kZXZlbnR0aW1lGOC176QBIAEoA1IMZW5kZXZlbnR0aW1lEi'
    'oKDnN0YXJ0ZXZlbnR0aW1lGN2/0cEBIAEoA1IOc3RhcnRldmVudHRpbWU=');

@$core.Deprecated('Use importStatisticsDescriptor instead')
const ImportStatistics$json = {
  '1': 'ImportStatistics',
  '2': [
    {
      '1': 'bytesimported',
      '3': 216355553,
      '4': 1,
      '5': 3,
      '10': 'bytesimported'
    },
  ],
};

/// Descriptor for `ImportStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importStatisticsDescriptor = $convert.base64Decode(
    'ChBJbXBvcnRTdGF0aXN0aWNzEicKDWJ5dGVzaW1wb3J0ZWQY4aWVZyABKANSDWJ5dGVzaW1wb3'
    'J0ZWQ=');

@$core.Deprecated('Use indexPolicyDescriptor instead')
const IndexPolicy$json = {
  '1': 'IndexPolicy',
  '2': [
    {
      '1': 'lastupdatetime',
      '3': 189280082,
      '4': 1,
      '5': 3,
      '10': 'lastupdatetime'
    },
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'source',
      '3': 466561497,
      '4': 1,
      '5': 14,
      '6': '.logs.IndexSource',
      '10': 'source'
    },
  ],
};

/// Descriptor for `IndexPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List indexPolicyDescriptor = $convert.base64Decode(
    'CgtJbmRleFBvbGljeRIpCg5sYXN0dXBkYXRldGltZRjS3qBaIAEoA1IObGFzdHVwZGF0ZXRpbW'
    'USMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIAEoCVISbG9nZ3JvdXBpZGVudGlmaWVyEikK'
    'DnBvbGljeWRvY3VtZW50GOvp9lQgASgJUg5wb2xpY3lkb2N1bWVudBIhCgpwb2xpY3luYW1lGN'
    '3i8jYgASgJUgpwb2xpY3luYW1lEi0KBnNvdXJjZRjZ07zeASABKA4yES5sb2dzLkluZGV4U291'
    'cmNlUgZzb3VyY2U=');

@$core.Deprecated('Use inputLogEventDescriptor instead')
const InputLogEvent$json = {
  '1': 'InputLogEvent',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'timestamp', '3': 310629668, '4': 1, '5': 3, '10': 'timestamp'},
  ],
};

/// Descriptor for `InputLogEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inputLogEventDescriptor = $convert.base64Decode(
    'Cg1JbnB1dExvZ0V2ZW50EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2USIAoJdGltZXN0YW'
    '1wGKSqj5QBIAEoA1IJdGltZXN0YW1w');

@$core.Deprecated('Use integrationDetailsDescriptor instead')
const IntegrationDetails$json = {
  '1': 'IntegrationDetails',
  '2': [
    {
      '1': 'opensearchintegrationdetails',
      '3': 425301806,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchIntegrationDetails',
      '10': 'opensearchintegrationdetails'
    },
  ],
};

/// Descriptor for `IntegrationDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List integrationDetailsDescriptor = $convert.base64Decode(
    'ChJJbnRlZ3JhdGlvbkRldGFpbHMSagocb3BlbnNlYXJjaGludGVncmF0aW9uZGV0YWlscxiuru'
    'bKASABKAsyIi5sb2dzLk9wZW5TZWFyY2hJbnRlZ3JhdGlvbkRldGFpbHNSHG9wZW5zZWFyY2hp'
    'bnRlZ3JhdGlvbmRldGFpbHM=');

@$core.Deprecated('Use integrationSummaryDescriptor instead')
const IntegrationSummary$json = {
  '1': 'IntegrationSummary',
  '2': [
    {
      '1': 'integrationname',
      '3': 183183535,
      '4': 1,
      '5': 9,
      '10': 'integrationname'
    },
    {
      '1': 'integrationstatus',
      '3': 12110360,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationStatus',
      '10': 'integrationstatus'
    },
    {
      '1': 'integrationtype',
      '3': 307087270,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationType',
      '10': 'integrationtype'
    },
  ],
};

/// Descriptor for `IntegrationSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List integrationSummaryDescriptor = $convert.base64Decode(
    'ChJJbnRlZ3JhdGlvblN1bW1hcnkSKwoPaW50ZWdyYXRpb25uYW1lGK/RrFcgASgJUg9pbnRlZ3'
    'JhdGlvbm5hbWUSSAoRaW50ZWdyYXRpb25zdGF0dXMYmJTjBSABKA4yFy5sb2dzLkludGVncmF0'
    'aW9uU3RhdHVzUhFpbnRlZ3JhdGlvbnN0YXR1cxJDCg9pbnRlZ3JhdGlvbnR5cGUYpo+3kgEgAS'
    'gOMhUubG9ncy5JbnRlZ3JhdGlvblR5cGVSD2ludGVncmF0aW9udHlwZQ==');

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

@$core.Deprecated('Use internalStreamingExceptionDescriptor instead')
const InternalStreamingException$json = {
  '1': 'InternalStreamingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalStreamingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalStreamingExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnRlcm5hbFN0cmVhbWluZ0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidOperationExceptionDescriptor instead')
const InvalidOperationException$json = {
  '1': 'InvalidOperationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOperationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOperationExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkT3BlcmF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidParameterExceptionDescriptor instead')
const InvalidParameterException$json = {
  '1': 'InvalidParameterException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkUGFyYW1ldGVyRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidSequenceTokenExceptionDescriptor instead')
const InvalidSequenceTokenException$json = {
  '1': 'InvalidSequenceTokenException',
  '2': [
    {
      '1': 'expectedsequencetoken',
      '3': 419198802,
      '4': 1,
      '5': 9,
      '10': 'expectedsequencetoken'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSequenceTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSequenceTokenExceptionDescriptor =
    $convert.base64Decode(
        'Ch1JbnZhbGlkU2VxdWVuY2VUb2tlbkV4Y2VwdGlvbhI4ChVleHBlY3RlZHNlcXVlbmNldG9rZW'
        '4Y0u7xxwEgASgJUhVleHBlY3RlZHNlcXVlbmNldG9rZW4SGwoHbWVzc2FnZRjlkcgnIAEoCVIH'
        'bWVzc2FnZQ==');

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

@$core.Deprecated('Use listAggregateLogGroupSummariesRequestDescriptor instead')
const ListAggregateLogGroupSummariesRequest$json = {
  '1': 'ListAggregateLogGroupSummariesRequest',
  '2': [
    {
      '1': 'accountidentifiers',
      '3': 349304053,
      '4': 3,
      '5': 9,
      '10': 'accountidentifiers'
    },
    {
      '1': 'datasources',
      '3': 477389554,
      '4': 3,
      '5': 11,
      '6': '.logs.DataSourceFilter',
      '10': 'datasources'
    },
    {
      '1': 'groupby',
      '3': 372746826,
      '4': 1,
      '5': 14,
      '6': '.logs.ListAggregateLogGroupSummariesGroupBy',
      '10': 'groupby'
    },
    {
      '1': 'includelinkedaccounts',
      '3': 56034131,
      '4': 1,
      '5': 8,
      '10': 'includelinkedaccounts'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'loggroupclass',
      '3': 518605953,
      '4': 1,
      '5': 14,
      '6': '.logs.LogGroupClass',
      '10': 'loggroupclass'
    },
    {
      '1': 'loggroupnamepattern',
      '3': 253299540,
      '4': 1,
      '5': 9,
      '10': 'loggroupnamepattern'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAggregateLogGroupSummariesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAggregateLogGroupSummariesRequestDescriptor = $convert.base64Decode(
    'CiVMaXN0QWdncmVnYXRlTG9nR3JvdXBTdW1tYXJpZXNSZXF1ZXN0EjIKEmFjY291bnRpZGVudG'
    'lmaWVycxj16cemASADKAlSEmFjY291bnRpZGVudGlmaWVycxI8CgtkYXRhc291cmNlcxjyxdHj'
    'ASADKAsyFi5sb2dzLkRhdGFTb3VyY2VGaWx0ZXJSC2RhdGFzb3VyY2VzEkkKB2dyb3VwYnkYyt'
    'TesQEgASgOMisubG9ncy5MaXN0QWdncmVnYXRlTG9nR3JvdXBTdW1tYXJpZXNHcm91cEJ5Ugdn'
    'cm91cGJ5EjcKFWluY2x1ZGVsaW5rZWRhY2NvdW50cxjThtwaIAEoCFIVaW5jbHVkZWxpbmtlZG'
    'FjY291bnRzEhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSPQoNbG9nZ3JvdXBjbGFzcxiBmaX3'
    'ASABKA4yEy5sb2dzLkxvZ0dyb3VwQ2xhc3NSDWxvZ2dyb3VwY2xhc3MSMwoTbG9nZ3JvdXBuYW'
    '1lcGF0dGVybhjUluR4IAEoCVITbG9nZ3JvdXBuYW1lcGF0dGVybhIfCgluZXh0dG9rZW4YnvOd'
    'NyABKAlSCW5leHR0b2tlbg==');

@$core
    .Deprecated('Use listAggregateLogGroupSummariesResponseDescriptor instead')
const ListAggregateLogGroupSummariesResponse$json = {
  '1': 'ListAggregateLogGroupSummariesResponse',
  '2': [
    {
      '1': 'aggregateloggroupsummaries',
      '3': 255519548,
      '4': 3,
      '5': 11,
      '6': '.logs.AggregateLogGroupSummary',
      '10': 'aggregateloggroupsummaries'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAggregateLogGroupSummariesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAggregateLogGroupSummariesResponseDescriptor =
    $convert.base64Decode(
        'CiZMaXN0QWdncmVnYXRlTG9nR3JvdXBTdW1tYXJpZXNSZXNwb25zZRJhChphZ2dyZWdhdGVsb2'
        'dncm91cHN1bW1hcmllcxi81ut5IAMoCzIeLmxvZ3MuQWdncmVnYXRlTG9nR3JvdXBTdW1tYXJ5'
        'UhphZ2dyZWdhdGVsb2dncm91cHN1bW1hcmllcxIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leH'
        'R0b2tlbg==');

@$core.Deprecated('Use listAnomaliesRequestDescriptor instead')
const ListAnomaliesRequest$json = {
  '1': 'ListAnomaliesRequest',
  '2': [
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'suppressionstate',
      '3': 124822782,
      '4': 1,
      '5': 14,
      '6': '.logs.SuppressionState',
      '10': 'suppressionstate'
    },
  ],
};

/// Descriptor for `ListAnomaliesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAnomaliesRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0QW5vbWFsaWVzUmVxdWVzdBIyChJhbm9tYWx5ZGV0ZWN0b3Jhcm4YrM/z1AEgASgJUh'
    'Jhbm9tYWx5ZGV0ZWN0b3Jhcm4SGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBIfCgluZXh0dG9r'
    'ZW4YnvOdNyABKAlSCW5leHR0b2tlbhJFChBzdXBwcmVzc2lvbnN0YXRlGP7JwjsgASgOMhYubG'
    '9ncy5TdXBwcmVzc2lvblN0YXRlUhBzdXBwcmVzc2lvbnN0YXRl');

@$core.Deprecated('Use listAnomaliesResponseDescriptor instead')
const ListAnomaliesResponse$json = {
  '1': 'ListAnomaliesResponse',
  '2': [
    {
      '1': 'anomalies',
      '3': 242998947,
      '4': 3,
      '5': 11,
      '6': '.logs.Anomaly',
      '10': 'anomalies'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAnomaliesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAnomaliesResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0QW5vbWFsaWVzUmVzcG9uc2USLgoJYW5vbWFsaWVzGKO973MgAygLMg0ubG9ncy5Bbm'
    '9tYWx5Uglhbm9tYWxpZXMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listIntegrationsRequestDescriptor instead')
const ListIntegrationsRequest$json = {
  '1': 'ListIntegrationsRequest',
  '2': [
    {
      '1': 'integrationnameprefix',
      '3': 275225987,
      '4': 1,
      '5': 9,
      '10': 'integrationnameprefix'
    },
    {
      '1': 'integrationstatus',
      '3': 12110360,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationStatus',
      '10': 'integrationstatus'
    },
    {
      '1': 'integrationtype',
      '3': 307087270,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationType',
      '10': 'integrationtype'
    },
  ],
};

/// Descriptor for `ListIntegrationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIntegrationsRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0SW50ZWdyYXRpb25zUmVxdWVzdBI4ChVpbnRlZ3JhdGlvbm5hbWVwcmVmaXgYg7uegw'
    'EgASgJUhVpbnRlZ3JhdGlvbm5hbWVwcmVmaXgSSAoRaW50ZWdyYXRpb25zdGF0dXMYmJTjBSAB'
    'KA4yFy5sb2dzLkludGVncmF0aW9uU3RhdHVzUhFpbnRlZ3JhdGlvbnN0YXR1cxJDCg9pbnRlZ3'
    'JhdGlvbnR5cGUYpo+3kgEgASgOMhUubG9ncy5JbnRlZ3JhdGlvblR5cGVSD2ludGVncmF0aW9u'
    'dHlwZQ==');

@$core.Deprecated('Use listIntegrationsResponseDescriptor instead')
const ListIntegrationsResponse$json = {
  '1': 'ListIntegrationsResponse',
  '2': [
    {
      '1': 'integrationsummaries',
      '3': 72020988,
      '4': 3,
      '5': 11,
      '6': '.logs.IntegrationSummary',
      '10': 'integrationsummaries'
    },
  ],
};

/// Descriptor for `ListIntegrationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIntegrationsResponseDescriptor =
    $convert.base64Decode(
        'ChhMaXN0SW50ZWdyYXRpb25zUmVzcG9uc2USTwoUaW50ZWdyYXRpb25zdW1tYXJpZXMY/OerIi'
        'ADKAsyGC5sb2dzLkludGVncmF0aW9uU3VtbWFyeVIUaW50ZWdyYXRpb25zdW1tYXJpZXM=');

@$core.Deprecated('Use listLogAnomalyDetectorsRequestDescriptor instead')
const ListLogAnomalyDetectorsRequest$json = {
  '1': 'ListLogAnomalyDetectorsRequest',
  '2': [
    {
      '1': 'filterloggrouparn',
      '3': 72011610,
      '4': 1,
      '5': 9,
      '10': 'filterloggrouparn'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListLogAnomalyDetectorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLogAnomalyDetectorsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0TG9nQW5vbWFseURldGVjdG9yc1JlcXVlc3QSLwoRZmlsdGVybG9nZ3JvdXBhcm4Y2p'
        '6rIiABKAlSEWZpbHRlcmxvZ2dyb3VwYXJuEhgKBWxpbWl0GLWy65YBIAEoBVIFbGltaXQSHwoJ'
        'bmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listLogAnomalyDetectorsResponseDescriptor instead')
const ListLogAnomalyDetectorsResponse$json = {
  '1': 'ListLogAnomalyDetectorsResponse',
  '2': [
    {
      '1': 'anomalydetectors',
      '3': 379268934,
      '4': 3,
      '5': 11,
      '6': '.logs.AnomalyDetector',
      '10': 'anomalydetectors'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListLogAnomalyDetectorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLogAnomalyDetectorsResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0TG9nQW5vbWFseURldGVjdG9yc1Jlc3BvbnNlEkUKEGFub21hbHlkZXRlY3RvcnMYxt'
        '7stAEgAygLMhUubG9ncy5Bbm9tYWx5RGV0ZWN0b3JSEGFub21hbHlkZXRlY3RvcnMSHwoJbmV4'
        'dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listLogGroupsForQueryRequestDescriptor instead')
const ListLogGroupsForQueryRequest$json = {
  '1': 'ListLogGroupsForQueryRequest',
  '2': [
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'queryid', '3': 336975759, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `ListLogGroupsForQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLogGroupsForQueryRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0TG9nR3JvdXBzRm9yUXVlcnlSZXF1ZXN0EiIKCm1heHJlc3VsdHMY0tzn3QEgASgFUg'
        'ptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2VuEhwKB3F1ZXJ5aWQY'
        'j6/XoAEgASgJUgdxdWVyeWlk');

@$core.Deprecated('Use listLogGroupsForQueryResponseDescriptor instead')
const ListLogGroupsForQueryResponse$json = {
  '1': 'ListLogGroupsForQueryResponse',
  '2': [
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListLogGroupsForQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLogGroupsForQueryResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0TG9nR3JvdXBzRm9yUXVlcnlSZXNwb25zZRI0ChNsb2dncm91cGlkZW50aWZpZXJzGP'
        'nauMMBIAMoCVITbG9nZ3JvdXBpZGVudGlmaWVycxIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5l'
        'eHR0b2tlbg==');

@$core.Deprecated('Use listLogGroupsRequestDescriptor instead')
const ListLogGroupsRequest$json = {
  '1': 'ListLogGroupsRequest',
  '2': [
    {
      '1': 'accountidentifiers',
      '3': 349304053,
      '4': 3,
      '5': 9,
      '10': 'accountidentifiers'
    },
    {
      '1': 'datasources',
      '3': 477389554,
      '4': 3,
      '5': 11,
      '6': '.logs.DataSourceFilter',
      '10': 'datasources'
    },
    {
      '1': 'fieldindexnames',
      '3': 300714984,
      '4': 3,
      '5': 9,
      '10': 'fieldindexnames'
    },
    {
      '1': 'includelinkedaccounts',
      '3': 56034131,
      '4': 1,
      '5': 8,
      '10': 'includelinkedaccounts'
    },
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'loggroupclass',
      '3': 518605953,
      '4': 1,
      '5': 14,
      '6': '.logs.LogGroupClass',
      '10': 'loggroupclass'
    },
    {
      '1': 'loggroupnamepattern',
      '3': 253299540,
      '4': 1,
      '5': 9,
      '10': 'loggroupnamepattern'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListLogGroupsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLogGroupsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0TG9nR3JvdXBzUmVxdWVzdBIyChJhY2NvdW50aWRlbnRpZmllcnMY9enHpgEgAygJUh'
    'JhY2NvdW50aWRlbnRpZmllcnMSPAoLZGF0YXNvdXJjZXMY8sXR4wEgAygLMhYubG9ncy5EYXRh'
    'U291cmNlRmlsdGVyUgtkYXRhc291cmNlcxIsCg9maWVsZGluZGV4bmFtZXMY6JeyjwEgAygJUg'
    '9maWVsZGluZGV4bmFtZXMSNwoVaW5jbHVkZWxpbmtlZGFjY291bnRzGNOG3BogASgIUhVpbmNs'
    'dWRlbGlua2VkYWNjb3VudHMSGAoFbGltaXQYtbLrlgEgASgFUgVsaW1pdBI9Cg1sb2dncm91cG'
    'NsYXNzGIGZpfcBIAEoDjITLmxvZ3MuTG9nR3JvdXBDbGFzc1INbG9nZ3JvdXBjbGFzcxIzChNs'
    'b2dncm91cG5hbWVwYXR0ZXJuGNSW5HggASgJUhNsb2dncm91cG5hbWVwYXR0ZXJuEh8KCW5leH'
    'R0b2tlbhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listLogGroupsResponseDescriptor instead')
const ListLogGroupsResponse$json = {
  '1': 'ListLogGroupsResponse',
  '2': [
    {
      '1': 'loggroups',
      '3': 507247778,
      '4': 3,
      '5': 11,
      '6': '.logs.LogGroupSummary',
      '10': 'loggroups'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListLogGroupsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLogGroupsResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0TG9nR3JvdXBzUmVzcG9uc2USNwoJbG9nZ3JvdXBzGKL57/EBIAMoCzIVLmxvZ3MuTG'
    '9nR3JvdXBTdW1tYXJ5Uglsb2dncm91cHMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9r'
    'ZW4=');

@$core.Deprecated('Use listScheduledQueriesRequestDescriptor instead')
const ListScheduledQueriesRequest$json = {
  '1': 'ListScheduledQueriesRequest',
  '2': [
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `ListScheduledQueriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listScheduledQueriesRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0U2NoZWR1bGVkUXVlcmllc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxjS3OfdASABKAVSCm'
        '1heHJlc3VsdHMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4SMwoFc3RhdGUY9+XE'
        'wQEgASgOMhkubG9ncy5TY2hlZHVsZWRRdWVyeVN0YXRlUgVzdGF0ZQ==');

@$core.Deprecated('Use listScheduledQueriesResponseDescriptor instead')
const ListScheduledQueriesResponse$json = {
  '1': 'ListScheduledQueriesResponse',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'scheduledqueries',
      '3': 275803785,
      '4': 3,
      '5': 11,
      '6': '.logs.ScheduledQuerySummary',
      '10': 'scheduledqueries'
    },
  ],
};

/// Descriptor for `ListScheduledQueriesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listScheduledQueriesResponseDescriptor =
    $convert.base64Decode(
        'ChxMaXN0U2NoZWR1bGVkUXVlcmllc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhie8503IAEoCVIJbm'
        'V4dHRva2VuEksKEHNjaGVkdWxlZHF1ZXJpZXMYid3BgwEgAygLMhsubG9ncy5TY2hlZHVsZWRR'
        'dWVyeVN1bW1hcnlSEHNjaGVkdWxlZHF1ZXJpZXM=');

@$core
    .Deprecated('Use listSourcesForS3TableIntegrationRequestDescriptor instead')
const ListSourcesForS3TableIntegrationRequest$json = {
  '1': 'ListSourcesForS3TableIntegrationRequest',
  '2': [
    {
      '1': 'integrationarn',
      '3': 432021733,
      '4': 1,
      '5': 9,
      '10': 'integrationarn'
    },
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListSourcesForS3TableIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSourcesForS3TableIntegrationRequestDescriptor =
    $convert.base64Decode(
        'CidMaXN0U291cmNlc0ZvclMzVGFibGVJbnRlZ3JhdGlvblJlcXVlc3QSKgoOaW50ZWdyYXRpb2'
        '5hcm4Y5cGAzgEgASgJUg5pbnRlZ3JhdGlvbmFybhIiCgptYXhyZXN1bHRzGNLc590BIAEoBVIK'
        'bWF4cmVzdWx0cxIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated(
    'Use listSourcesForS3TableIntegrationResponseDescriptor instead')
const ListSourcesForS3TableIntegrationResponse$json = {
  '1': 'ListSourcesForS3TableIntegrationResponse',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'sources',
      '3': 303994930,
      '4': 3,
      '5': 11,
      '6': '.logs.S3TableIntegrationSource',
      '10': 'sources'
    },
  ],
};

/// Descriptor for `ListSourcesForS3TableIntegrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSourcesForS3TableIntegrationResponseDescriptor =
    $convert.base64Decode(
        'CihMaXN0U291cmNlc0ZvclMzVGFibGVJbnRlZ3JhdGlvblJlc3BvbnNlEh8KCW5leHR0b2tlbh'
        'ie8503IAEoCVIJbmV4dHRva2VuEjwKB3NvdXJjZXMYsrD6kAEgAygLMh4ubG9ncy5TM1RhYmxl'
        'SW50ZWdyYXRpb25Tb3VyY2VSB3NvdXJjZXM=');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIjCgtyZXNvdXJjZWFybhjNzKogIAEoCVILcm'
        'Vzb3VyY2Vhcm4=');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.ListTagsForResourceResponse.TagsEntry',
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
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USQwoEdGFncxih19ugASADKAsyKy5sb2dzLk'
        'xpc3RUYWdzRm9yUmVzb3VyY2VSZXNwb25zZS5UYWdzRW50cnlSBHRhZ3MaNwoJVGFnc0VudHJ5'
        'EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use listTagsLogGroupRequestDescriptor instead')
const ListTagsLogGroupRequest$json = {
  '1': 'ListTagsLogGroupRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
  ],
};

/// Descriptor for `ListTagsLogGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsLogGroupRequestDescriptor =
    $convert.base64Decode(
        'ChdMaXN0VGFnc0xvZ0dyb3VwUmVxdWVzdBIlCgxsb2dncm91cG5hbWUYzL/ULSABKAlSDGxvZ2'
        'dyb3VwbmFtZQ==');

@$core.Deprecated('Use listTagsLogGroupResponseDescriptor instead')
const ListTagsLogGroupResponse$json = {
  '1': 'ListTagsLogGroupResponse',
  '2': [
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.ListTagsLogGroupResponse.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [ListTagsLogGroupResponse_TagsEntry$json],
};

@$core.Deprecated('Use listTagsLogGroupResponseDescriptor instead')
const ListTagsLogGroupResponse_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListTagsLogGroupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsLogGroupResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0VGFnc0xvZ0dyb3VwUmVzcG9uc2USQAoEdGFncxih19ugASADKAsyKC5sb2dzLkxpc3'
    'RUYWdzTG9nR3JvdXBSZXNwb25zZS5UYWdzRW50cnlSBHRhZ3MaNwoJVGFnc0VudHJ5EhAKA2tl'
    'eRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use listToMapDescriptor instead')
const ListToMap$json = {
  '1': 'ListToMap',
  '2': [
    {'1': 'flatten', '3': 100024266, '4': 1, '5': 8, '10': 'flatten'},
    {
      '1': 'flattenedelement',
      '3': 9195721,
      '4': 1,
      '5': 14,
      '6': '.logs.FlattenedElement',
      '10': 'flattenedelement'
    },
    {'1': 'key', '3': 135645293, '4': 1, '5': 9, '10': 'key'},
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
    {'1': 'target', '3': 308316233, '4': 1, '5': 9, '10': 'target'},
    {'1': 'valuekey', '3': 260470114, '4': 1, '5': 9, '10': 'valuekey'},
  ],
};

/// Descriptor for `ListToMap`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listToMapDescriptor = $convert.base64Decode(
    'CglMaXN0VG9NYXASGwoHZmxhdHRlbhjK/9gvIAEoCFIHZmxhdHRlbhJFChBmbGF0dGVuZWRlbG'
    'VtZW50GMmhsQQgASgOMhYubG9ncy5GbGF0dGVuZWRFbGVtZW50UhBmbGF0dGVuZWRlbGVtZW50'
    'EhMKA2tleRjtkNdAIAEoCVIDa2V5EhoKBnNvdXJjZRjZ07zeASABKAlSBnNvdXJjZRIaCgZ0YX'
    'JnZXQYyZCCkwEgASgJUgZ0YXJnZXQSHQoIdmFsdWVrZXkY4uqZfCABKAlSCHZhbHVla2V5');

@$core.Deprecated('Use liveTailSessionLogEventDescriptor instead')
const LiveTailSessionLogEvent$json = {
  '1': 'LiveTailSessionLogEvent',
  '2': [
    {
      '1': 'ingestiontime',
      '3': 179367957,
      '4': 1,
      '5': 3,
      '10': 'ingestiontime'
    },
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'timestamp', '3': 310629668, '4': 1, '5': 3, '10': 'timestamp'},
  ],
};

/// Descriptor for `LiveTailSessionLogEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List liveTailSessionLogEventDescriptor = $convert.base64Decode(
    'ChdMaXZlVGFpbFNlc3Npb25Mb2dFdmVudBInCg1pbmdlc3Rpb250aW1lGJXgw1UgASgDUg1pbm'
    'dlc3Rpb250aW1lEjIKEmxvZ2dyb3VwaWRlbnRpZmllchjcz6XfASABKAlSEmxvZ2dyb3VwaWRl'
    'bnRpZmllchIoCg1sb2dzdHJlYW1uYW1lGKP37tABIAEoCVINbG9nc3RyZWFtbmFtZRIbCgdtZX'
    'NzYWdlGOWRyCcgASgJUgdtZXNzYWdlEiAKCXRpbWVzdGFtcBikqo+UASABKANSCXRpbWVzdGFt'
    'cA==');

@$core.Deprecated('Use liveTailSessionMetadataDescriptor instead')
const LiveTailSessionMetadata$json = {
  '1': 'LiveTailSessionMetadata',
  '2': [
    {'1': 'sampled', '3': 187095290, '4': 1, '5': 8, '10': 'sampled'},
  ],
};

/// Descriptor for `LiveTailSessionMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List liveTailSessionMetadataDescriptor =
    $convert.base64Decode(
        'ChdMaXZlVGFpbFNlc3Npb25NZXRhZGF0YRIbCgdzYW1wbGVkGPqxm1kgASgIUgdzYW1wbGVk');

@$core.Deprecated('Use liveTailSessionStartDescriptor instead')
const LiveTailSessionStart$json = {
  '1': 'LiveTailSessionStart',
  '2': [
    {
      '1': 'logeventfilterpattern',
      '3': 81051802,
      '4': 1,
      '5': 9,
      '10': 'logeventfilterpattern'
    },
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {
      '1': 'logstreamnameprefixes',
      '3': 109414303,
      '4': 3,
      '5': 9,
      '10': 'logstreamnameprefixes'
    },
    {
      '1': 'logstreamnames',
      '3': 178825732,
      '4': 3,
      '5': 9,
      '10': 'logstreamnames'
    },
    {'1': 'requestid', '3': 376827552, '4': 1, '5': 9, '10': 'requestid'},
    {'1': 'sessionid', '3': 505783131, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `LiveTailSessionStart`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List liveTailSessionStartDescriptor = $convert.base64Decode(
    'ChRMaXZlVGFpbFNlc3Npb25TdGFydBI3ChVsb2dldmVudGZpbHRlcnBhdHRlcm4YmoHTJiABKA'
    'lSFWxvZ2V2ZW50ZmlsdGVycGF0dGVybhI0ChNsb2dncm91cGlkZW50aWZpZXJzGPnauMMBIAMo'
    'CVITbG9nZ3JvdXBpZGVudGlmaWVycxI3ChVsb2dzdHJlYW1uYW1lcHJlZml4ZXMYn4+WNCADKA'
    'lSFWxvZ3N0cmVhbW5hbWVwcmVmaXhlcxIpCg5sb2dzdHJlYW1uYW1lcxiE1KJVIAMoCVIObG9n'
    'c3RyZWFtbmFtZXMSIAoJcmVxdWVzdGlkGKDd17MBIAEoCVIJcmVxdWVzdGlkEiAKCXNlc3Npb2'
    '5pZBjbxpbxASABKAlSCXNlc3Npb25pZA==');

@$core.Deprecated('Use liveTailSessionUpdateDescriptor instead')
const LiveTailSessionUpdate$json = {
  '1': 'LiveTailSessionUpdate',
  '2': [
    {
      '1': 'sessionmetadata',
      '3': 7268203,
      '4': 1,
      '5': 11,
      '6': '.logs.LiveTailSessionMetadata',
      '10': 'sessionmetadata'
    },
    {
      '1': 'sessionresults',
      '3': 292392456,
      '4': 3,
      '5': 11,
      '6': '.logs.LiveTailSessionLogEvent',
      '10': 'sessionresults'
    },
  ],
};

/// Descriptor for `LiveTailSessionUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List liveTailSessionUpdateDescriptor = $convert.base64Decode(
    'ChVMaXZlVGFpbFNlc3Npb25VcGRhdGUSSgoPc2Vzc2lvbm1ldGFkYXRhGOvOuwMgASgLMh0ubG'
    '9ncy5MaXZlVGFpbFNlc3Npb25NZXRhZGF0YVIPc2Vzc2lvbm1ldGFkYXRhEkkKDnNlc3Npb25y'
    'ZXN1bHRzGIictosBIAMoCzIdLmxvZ3MuTGl2ZVRhaWxTZXNzaW9uTG9nRXZlbnRSDnNlc3Npb2'
    '5yZXN1bHRz');

@$core.Deprecated('Use logEventDescriptor instead')
const LogEvent$json = {
  '1': 'LogEvent',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'timestamp', '3': 310629668, '4': 1, '5': 3, '10': 'timestamp'},
  ],
};

/// Descriptor for `LogEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logEventDescriptor = $convert.base64Decode(
    'CghMb2dFdmVudBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdlEiAKCXRpbWVzdGFtcBikqo'
    '+UASABKANSCXRpbWVzdGFtcA==');

@$core.Deprecated('Use logFieldTypeDescriptor instead')
const LogFieldType$json = {
  '1': 'LogFieldType',
  '2': [
    {
      '1': 'element',
      '3': 256719864,
      '4': 1,
      '5': 11,
      '6': '.logs.LogFieldType',
      '10': 'element'
    },
    {
      '1': 'fields',
      '3': 104883581,
      '4': 3,
      '5': 11,
      '6': '.logs.LogFieldsListItem',
      '10': 'fields'
    },
    {'1': 'type', '3': 287830350, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `LogFieldType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logFieldTypeDescriptor = $convert.base64Decode(
    'CgxMb2dGaWVsZFR5cGUSLwoHZWxlbWVudBj497R6IAEoCzISLmxvZ3MuTG9nRmllbGRUeXBlUg'
    'dlbGVtZW50EjIKBmZpZWxkcxj9yoEyIAMoCzIXLmxvZ3MuTG9nRmllbGRzTGlzdEl0ZW1SBmZp'
    'ZWxkcxIWCgR0eXBlGM7in4kBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use logFieldsListItemDescriptor instead')
const LogFieldsListItem$json = {
  '1': 'LogFieldsListItem',
  '2': [
    {'1': 'logfieldname', '3': 243412097, '4': 1, '5': 9, '10': 'logfieldname'},
    {
      '1': 'logfieldtype',
      '3': 491059716,
      '4': 1,
      '5': 11,
      '6': '.logs.LogFieldType',
      '10': 'logfieldtype'
    },
  ],
};

/// Descriptor for `LogFieldsListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logFieldsListItemDescriptor = $convert.base64Decode(
    'ChFMb2dGaWVsZHNMaXN0SXRlbRIlCgxsb2dmaWVsZG5hbWUYgdmIdCABKAlSDGxvZ2ZpZWxkbm'
    'FtZRI6Cgxsb2dmaWVsZHR5cGUYhPST6gEgASgLMhIubG9ncy5Mb2dGaWVsZFR5cGVSDGxvZ2Zp'
    'ZWxkdHlwZQ==');

@$core.Deprecated('Use logGroupDescriptor instead')
const LogGroup$json = {
  '1': 'LogGroup',
  '2': [
    {'1': 'arn', '3': 359604989, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'bearertokenauthenticationenabled',
      '3': 513390155,
      '4': 1,
      '5': 8,
      '10': 'bearertokenauthenticationenabled'
    },
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'dataprotectionstatus',
      '3': 24294469,
      '4': 1,
      '5': 14,
      '6': '.logs.DataProtectionStatus',
      '10': 'dataprotectionstatus'
    },
    {
      '1': 'deletionprotectionenabled',
      '3': 475522738,
      '4': 1,
      '5': 8,
      '10': 'deletionprotectionenabled'
    },
    {
      '1': 'inheritedproperties',
      '3': 178880393,
      '4': 3,
      '5': 14,
      '6': '.logs.InheritedProperty',
      '10': 'inheritedproperties'
    },
    {'1': 'kmskeyid', '3': 510698477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'loggrouparn', '3': 6742512, '4': 1, '5': 9, '10': 'loggrouparn'},
    {
      '1': 'loggroupclass',
      '3': 518605953,
      '4': 1,
      '5': 14,
      '6': '.logs.LogGroupClass',
      '10': 'loggroupclass'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'metricfiltercount',
      '3': 347145747,
      '4': 1,
      '5': 5,
      '10': 'metricfiltercount'
    },
    {
      '1': 'retentionindays',
      '3': 258337482,
      '4': 1,
      '5': 5,
      '10': 'retentionindays'
    },
    {'1': 'storedbytes', '3': 260235622, '4': 1, '5': 3, '10': 'storedbytes'},
  ],
};

/// Descriptor for `LogGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logGroupDescriptor = $convert.base64Decode(
    'CghMb2dHcm91cBIUCgNhcm4Y/cW8qwEgASgJUgNhcm4STgogYmVhcmVydG9rZW5hdXRoZW50aW'
    'NhdGlvbmVuYWJsZWQYy+zm9AEgASgIUiBiZWFyZXJ0b2tlbmF1dGhlbnRpY2F0aW9uZW5hYmxl'
    'ZBIlCgxjcmVhdGlvbnRpbWUYhsXEGSABKANSDGNyZWF0aW9udGltZRJRChRkYXRhcHJvdGVjdG'
    'lvbnN0YXR1cxjF6MoLIAEoDjIaLmxvZ3MuRGF0YVByb3RlY3Rpb25TdGF0dXNSFGRhdGFwcm90'
    'ZWN0aW9uc3RhdHVzEkAKGWRlbGV0aW9ucHJvdGVjdGlvbmVuYWJsZWQYss3f4gEgASgIUhlkZW'
    'xldGlvbnByb3RlY3Rpb25lbmFibGVkEkwKE2luaGVyaXRlZHByb3BlcnRpZXMYif+lVSADKA4y'
    'Fy5sb2dzLkluaGVyaXRlZFByb3BlcnR5UhNpbmhlcml0ZWRwcm9wZXJ0aWVzEh4KCGttc2tleW'
    'lkGO3HwvMBIAEoCVIIa21za2V5aWQSIwoLbG9nZ3JvdXBhcm4Y8MObAyABKAlSC2xvZ2dyb3Vw'
    'YXJuEj0KDWxvZ2dyb3VwY2xhc3MYgZml9wEgASgOMhMubG9ncy5Mb2dHcm91cENsYXNzUg1sb2'
    'dncm91cGNsYXNzEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW1lEjAKEW1l'
    'dHJpY2ZpbHRlcmNvdW50GJOMxKUBIAEoBVIRbWV0cmljZmlsdGVyY291bnQSKwoPcmV0ZW50aW'
    '9uaW5kYXlzGMrVl3sgASgFUg9yZXRlbnRpb25pbmRheXMSIwoLc3RvcmVkYnl0ZXMY5sKLfCAB'
    'KANSC3N0b3JlZGJ5dGVz');

@$core.Deprecated('Use logGroupFieldDescriptor instead')
const LogGroupField$json = {
  '1': 'LogGroupField',
  '2': [
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'percent', '3': 368704091, '4': 1, '5': 5, '10': 'percent'},
  ],
};

/// Descriptor for `LogGroupField`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logGroupFieldDescriptor = $convert.base64Decode(
    'Cg1Mb2dHcm91cEZpZWxkEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSHAoHcGVyY2VudBjb9OevAS'
    'ABKAVSB3BlcmNlbnQ=');

@$core.Deprecated('Use logGroupSummaryDescriptor instead')
const LogGroupSummary$json = {
  '1': 'LogGroupSummary',
  '2': [
    {'1': 'loggrouparn', '3': 6742512, '4': 1, '5': 9, '10': 'loggrouparn'},
    {
      '1': 'loggroupclass',
      '3': 518605953,
      '4': 1,
      '5': 14,
      '6': '.logs.LogGroupClass',
      '10': 'loggroupclass'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
  ],
};

/// Descriptor for `LogGroupSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logGroupSummaryDescriptor = $convert.base64Decode(
    'Cg9Mb2dHcm91cFN1bW1hcnkSIwoLbG9nZ3JvdXBhcm4Y8MObAyABKAlSC2xvZ2dyb3VwYXJuEj'
    '0KDWxvZ2dyb3VwY2xhc3MYgZml9wEgASgOMhMubG9ncy5Mb2dHcm91cENsYXNzUg1sb2dncm91'
    'cGNsYXNzEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW1l');

@$core.Deprecated('Use logStreamDescriptor instead')
const LogStream$json = {
  '1': 'LogStream',
  '2': [
    {'1': 'arn', '3': 359604989, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'firsteventtimestamp',
      '3': 354762842,
      '4': 1,
      '5': 3,
      '10': 'firsteventtimestamp'
    },
    {
      '1': 'lasteventtimestamp',
      '3': 338474610,
      '4': 1,
      '5': 3,
      '10': 'lasteventtimestamp'
    },
    {
      '1': 'lastingestiontime',
      '3': 56929529,
      '4': 1,
      '5': 3,
      '10': 'lastingestiontime'
    },
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
    {'1': 'storedbytes', '3': 260235622, '4': 1, '5': 3, '10': 'storedbytes'},
    {
      '1': 'uploadsequencetoken',
      '3': 139889111,
      '4': 1,
      '5': 9,
      '10': 'uploadsequencetoken'
    },
  ],
};

/// Descriptor for `LogStream`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logStreamDescriptor = $convert.base64Decode(
    'CglMb2dTdHJlYW0SFAoDYXJuGP3FvKsBIAEoCVIDYXJuEiUKDGNyZWF0aW9udGltZRiGxcQZIA'
    'EoA1IMY3JlYXRpb250aW1lEjQKE2ZpcnN0ZXZlbnR0aW1lc3RhbXAY2oCVqQEgASgDUhNmaXJz'
    'dGV2ZW50dGltZXN0YW1wEjIKEmxhc3RldmVudHRpbWVzdGFtcBjy7LKhASABKANSEmxhc3Rldm'
    'VudHRpbWVzdGFtcBIvChFsYXN0aW5nZXN0aW9udGltZRj52ZIbIAEoA1IRbGFzdGluZ2VzdGlv'
    'bnRpbWUSKAoNbG9nc3RyZWFtbmFtZRij9+7QASABKAlSDWxvZ3N0cmVhbW5hbWUSIwoLc3Rvcm'
    'VkYnl0ZXMY5sKLfCABKANSC3N0b3JlZGJ5dGVzEjMKE3VwbG9hZHNlcXVlbmNldG9rZW4Y15Pa'
    'QiABKAlSE3VwbG9hZHNlcXVlbmNldG9rZW4=');

@$core.Deprecated('Use lowerCaseStringDescriptor instead')
const LowerCaseString$json = {
  '1': 'LowerCaseString',
  '2': [
    {'1': 'withkeys', '3': 161392106, '4': 3, '5': 9, '10': 'withkeys'},
  ],
};

/// Descriptor for `LowerCaseString`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lowerCaseStringDescriptor = $convert.base64Decode(
    'Cg9Mb3dlckNhc2VTdHJpbmcSHQoId2l0aGtleXMY6sv6TCADKAlSCHdpdGhrZXlz');

@$core.Deprecated('Use malformedQueryExceptionDescriptor instead')
const MalformedQueryException$json = {
  '1': 'MalformedQueryException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'querycompileerror',
      '3': 116576135,
      '4': 1,
      '5': 11,
      '6': '.logs.QueryCompileError',
      '10': 'querycompileerror'
    },
  ],
};

/// Descriptor for `MalformedQueryException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List malformedQueryExceptionDescriptor = $convert.base64Decode(
    'ChdNYWxmb3JtZWRRdWVyeUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdlEk'
    'gKEXF1ZXJ5Y29tcGlsZWVycm9yGIefyzcgASgLMhcubG9ncy5RdWVyeUNvbXBpbGVFcnJvclIR'
    'cXVlcnljb21waWxlZXJyb3I=');

@$core.Deprecated('Use metricFilterDescriptor instead')
const MetricFilter$json = {
  '1': 'MetricFilter',
  '2': [
    {
      '1': 'applyontransformedlogs',
      '3': 99775525,
      '4': 1,
      '5': 8,
      '10': 'applyontransformedlogs'
    },
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'emitsystemfielddimensions',
      '3': 437431321,
      '4': 3,
      '5': 9,
      '10': 'emitsystemfielddimensions'
    },
    {
      '1': 'fieldselectioncriteria',
      '3': 303984807,
      '4': 1,
      '5': 9,
      '10': 'fieldselectioncriteria'
    },
    {'1': 'filtername', '3': 395125013, '4': 1, '5': 9, '10': 'filtername'},
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'metrictransformations',
      '3': 169353806,
      '4': 3,
      '5': 11,
      '6': '.logs.MetricTransformation',
      '10': 'metrictransformations'
    },
  ],
};

/// Descriptor for `MetricFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricFilterDescriptor = $convert.base64Decode(
    'CgxNZXRyaWNGaWx0ZXISOQoWYXBwbHlvbnRyYW5zZm9ybWVkbG9ncxil6MkvIAEoCFIWYXBwbH'
    'lvbnRyYW5zZm9ybWVkbG9ncxIlCgxjcmVhdGlvbnRpbWUYhsXEGSABKANSDGNyZWF0aW9udGlt'
    'ZRJAChllbWl0c3lzdGVtZmllbGRkaW1lbnNpb25zGJnYytABIAMoCVIZZW1pdHN5c3RlbWZpZW'
    'xkZGltZW5zaW9ucxI6ChZmaWVsZHNlbGVjdGlvbmNyaXRlcmlhGKfh+ZABIAEoCVIWZmllbGRz'
    'ZWxlY3Rpb25jcml0ZXJpYRIiCgpmaWx0ZXJuYW1lGJXCtLwBIAEoCVIKZmlsdGVybmFtZRInCg'
    '1maWx0ZXJwYXR0ZXJuGJiHikUgASgJUg1maWx0ZXJwYXR0ZXJuEiUKDGxvZ2dyb3VwbmFtZRjM'
    'v9QtIAEoCVIMbG9nZ3JvdXBuYW1lElMKFW1ldHJpY3RyYW5zZm9ybWF0aW9ucxjOxOBQIAMoCz'
    'IaLmxvZ3MuTWV0cmljVHJhbnNmb3JtYXRpb25SFW1ldHJpY3RyYW5zZm9ybWF0aW9ucw==');

@$core.Deprecated('Use metricFilterMatchRecordDescriptor instead')
const MetricFilterMatchRecord$json = {
  '1': 'MetricFilterMatchRecord',
  '2': [
    {'1': 'eventmessage', '3': 299743039, '4': 1, '5': 9, '10': 'eventmessage'},
    {'1': 'eventnumber', '3': 220470463, '4': 1, '5': 3, '10': 'eventnumber'},
    {
      '1': 'extractedvalues',
      '3': 341078182,
      '4': 3,
      '5': 11,
      '6': '.logs.MetricFilterMatchRecord.ExtractedvaluesEntry',
      '10': 'extractedvalues'
    },
  ],
  '3': [MetricFilterMatchRecord_ExtractedvaluesEntry$json],
};

@$core.Deprecated('Use metricFilterMatchRecordDescriptor instead')
const MetricFilterMatchRecord_ExtractedvaluesEntry$json = {
  '1': 'ExtractedvaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `MetricFilterMatchRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricFilterMatchRecordDescriptor = $convert.base64Decode(
    'ChdNZXRyaWNGaWx0ZXJNYXRjaFJlY29yZBImCgxldmVudG1lc3NhZ2UYv+72jgEgASgJUgxldm'
    'VudG1lc3NhZ2USIwoLZXZlbnRudW1iZXIYv7mQaSABKANSC2V2ZW50bnVtYmVyEmAKD2V4dHJh'
    'Y3RlZHZhbHVlcxim4dGiASADKAsyMi5sb2dzLk1ldHJpY0ZpbHRlck1hdGNoUmVjb3JkLkV4dH'
    'JhY3RlZHZhbHVlc0VudHJ5Ug9leHRyYWN0ZWR2YWx1ZXMaQgoURXh0cmFjdGVkdmFsdWVzRW50'
    'cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use metricTransformationDescriptor instead')
const MetricTransformation$json = {
  '1': 'MetricTransformation',
  '2': [
    {'1': 'defaultvalue', '3': 403858624, '4': 1, '5': 1, '10': 'defaultvalue'},
    {
      '1': 'dimensions',
      '3': 116965553,
      '4': 3,
      '5': 11,
      '6': '.logs.MetricTransformation.DimensionsEntry',
      '10': 'dimensions'
    },
    {'1': 'metricname', '3': 204020635, '4': 1, '5': 9, '10': 'metricname'},
    {
      '1': 'metricnamespace',
      '3': 315894261,
      '4': 1,
      '5': 9,
      '10': 'metricnamespace'
    },
    {'1': 'metricvalue', '3': 51084559, '4': 1, '5': 9, '10': 'metricvalue'},
    {
      '1': 'unit',
      '3': 146086408,
      '4': 1,
      '5': 14,
      '6': '.logs.StandardUnit',
      '10': 'unit'
    },
  ],
  '3': [MetricTransformation_DimensionsEntry$json],
};

@$core.Deprecated('Use metricTransformationDescriptor instead')
const MetricTransformation_DimensionsEntry$json = {
  '1': 'DimensionsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `MetricTransformation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricTransformationDescriptor = $convert.base64Decode(
    'ChRNZXRyaWNUcmFuc2Zvcm1hdGlvbhImCgxkZWZhdWx0dmFsdWUYwMnJwAEgASgBUgxkZWZhdW'
    'x0dmFsdWUSTQoKZGltZW5zaW9ucxixgeM3IAMoCzIqLmxvZ3MuTWV0cmljVHJhbnNmb3JtYXRp'
    'b24uRGltZW5zaW9uc0VudHJ5UgpkaW1lbnNpb25zEiEKCm1ldHJpY25hbWUYm7ekYSABKAlSCm'
    '1ldHJpY25hbWUSLAoPbWV0cmljbmFtZXNwYWNlGPXT0JYBIAEoCVIPbWV0cmljbmFtZXNwYWNl'
    'EiMKC21ldHJpY3ZhbHVlGI/6rRggASgJUgttZXRyaWN2YWx1ZRIpCgR1bml0GIi01EUgASgOMh'
    'IubG9ncy5TdGFuZGFyZFVuaXRSBHVuaXQaPQoPRGltZW5zaW9uc0VudHJ5EhAKA2tleRgBIAEo'
    'CVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use moveKeyEntryDescriptor instead')
const MoveKeyEntry$json = {
  '1': 'MoveKeyEntry',
  '2': [
    {
      '1': 'overwriteifexists',
      '3': 230880030,
      '4': 1,
      '5': 8,
      '10': 'overwriteifexists'
    },
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
    {'1': 'target', '3': 308316233, '4': 1, '5': 9, '10': 'target'},
  ],
};

/// Descriptor for `MoveKeyEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List moveKeyEntryDescriptor = $convert.base64Decode(
    'CgxNb3ZlS2V5RW50cnkSLwoRb3ZlcndyaXRlaWZleGlzdHMYnuaLbiABKAhSEW92ZXJ3cml0ZW'
    'lmZXhpc3RzEhoKBnNvdXJjZRjZ07zeASABKAlSBnNvdXJjZRIaCgZ0YXJnZXQYyZCCkwEgASgJ'
    'UgZ0YXJnZXQ=');

@$core.Deprecated('Use moveKeysDescriptor instead')
const MoveKeys$json = {
  '1': 'MoveKeys',
  '2': [
    {
      '1': 'entries',
      '3': 257458932,
      '4': 3,
      '5': 11,
      '6': '.logs.MoveKeyEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `MoveKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List moveKeysDescriptor = $convert.base64Decode(
    'CghNb3ZlS2V5cxIvCgdlbnRyaWVzGPSF4nogAygLMhIubG9ncy5Nb3ZlS2V5RW50cnlSB2VudH'
    'JpZXM=');

@$core.Deprecated('Use openSearchApplicationDescriptor instead')
const OpenSearchApplication$json = {
  '1': 'OpenSearchApplication',
  '2': [
    {
      '1': 'applicationarn',
      '3': 230367005,
      '4': 1,
      '5': 9,
      '10': 'applicationarn'
    },
    {
      '1': 'applicationendpoint',
      '3': 379513629,
      '4': 1,
      '5': 9,
      '10': 'applicationendpoint'
    },
    {
      '1': 'applicationid',
      '3': 531796353,
      '4': 1,
      '5': 9,
      '10': 'applicationid'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `OpenSearchApplication`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchApplicationDescriptor = $convert.base64Decode(
    'ChVPcGVuU2VhcmNoQXBwbGljYXRpb24SKQoOYXBwbGljYXRpb25hcm4Ynb7sbSABKAlSDmFwcG'
    'xpY2F0aW9uYXJuEjQKE2FwcGxpY2F0aW9uZW5kcG9pbnQYndb7tAEgASgJUhNhcHBsaWNhdGlv'
    'bmVuZHBvaW50EigKDWFwcGxpY2F0aW9uaWQYgaPK/QEgASgJUg1hcHBsaWNhdGlvbmlkEjoKBn'
    'N0YXR1cxjw763SASABKAsyHi5sb2dzLk9wZW5TZWFyY2hSZXNvdXJjZVN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use openSearchCollectionDescriptor instead')
const OpenSearchCollection$json = {
  '1': 'OpenSearchCollection',
  '2': [
    {
      '1': 'collectionarn',
      '3': 461567321,
      '4': 1,
      '5': 9,
      '10': 'collectionarn'
    },
    {
      '1': 'collectionendpoint',
      '3': 354351753,
      '4': 1,
      '5': 9,
      '10': 'collectionendpoint'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `OpenSearchCollection`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchCollectionDescriptor = $convert.base64Decode(
    'ChRPcGVuU2VhcmNoQ29sbGVjdGlvbhIoCg1jb2xsZWN0aW9uYXJuGNnqi9wBIAEoCVINY29sbG'
    'VjdGlvbmFybhIyChJjb2xsZWN0aW9uZW5kcG9pbnQYifX7qAEgASgJUhJjb2xsZWN0aW9uZW5k'
    'cG9pbnQSOgoGc3RhdHVzGPDvrdIBIAEoCzIeLmxvZ3MuT3BlblNlYXJjaFJlc291cmNlU3RhdH'
    'VzUgZzdGF0dXM=');

@$core.Deprecated('Use openSearchDataAccessPolicyDescriptor instead')
const OpenSearchDataAccessPolicy$json = {
  '1': 'OpenSearchDataAccessPolicy',
  '2': [
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `OpenSearchDataAccessPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchDataAccessPolicyDescriptor =
    $convert.base64Decode(
        'ChpPcGVuU2VhcmNoRGF0YUFjY2Vzc1BvbGljeRIhCgpwb2xpY3luYW1lGN3i8jYgASgJUgpwb2'
        'xpY3luYW1lEjoKBnN0YXR1cxjw763SASABKAsyHi5sb2dzLk9wZW5TZWFyY2hSZXNvdXJjZVN0'
        'YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use openSearchDataSourceDescriptor instead')
const OpenSearchDataSource$json = {
  '1': 'OpenSearchDataSource',
  '2': [
    {
      '1': 'datasourcename',
      '3': 231923996,
      '4': 1,
      '5': 9,
      '10': 'datasourcename'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `OpenSearchDataSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchDataSourceDescriptor = $convert.base64Decode(
    'ChRPcGVuU2VhcmNoRGF0YVNvdXJjZRIpCg5kYXRhc291cmNlbmFtZRicwstuIAEoCVIOZGF0YX'
    'NvdXJjZW5hbWUSOgoGc3RhdHVzGPDvrdIBIAEoCzIeLmxvZ3MuT3BlblNlYXJjaFJlc291cmNl'
    'U3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use openSearchEncryptionPolicyDescriptor instead')
const OpenSearchEncryptionPolicy$json = {
  '1': 'OpenSearchEncryptionPolicy',
  '2': [
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `OpenSearchEncryptionPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchEncryptionPolicyDescriptor =
    $convert.base64Decode(
        'ChpPcGVuU2VhcmNoRW5jcnlwdGlvblBvbGljeRIhCgpwb2xpY3luYW1lGN3i8jYgASgJUgpwb2'
        'xpY3luYW1lEjoKBnN0YXR1cxjw763SASABKAsyHi5sb2dzLk9wZW5TZWFyY2hSZXNvdXJjZVN0'
        'YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use openSearchIntegrationDetailsDescriptor instead')
const OpenSearchIntegrationDetails$json = {
  '1': 'OpenSearchIntegrationDetails',
  '2': [
    {
      '1': 'accesspolicy',
      '3': 336061518,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchDataAccessPolicy',
      '10': 'accesspolicy'
    },
    {
      '1': 'application',
      '3': 524788294,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchApplication',
      '10': 'application'
    },
    {
      '1': 'collection',
      '3': 531465362,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchCollection',
      '10': 'collection'
    },
    {
      '1': 'datasource',
      '3': 345762713,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchDataSource',
      '10': 'datasource'
    },
    {
      '1': 'encryptionpolicy',
      '3': 70395159,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchEncryptionPolicy',
      '10': 'encryptionpolicy'
    },
    {
      '1': 'lifecyclepolicy',
      '3': 210416482,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchLifecyclePolicy',
      '10': 'lifecyclepolicy'
    },
    {
      '1': 'networkpolicy',
      '3': 200333472,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchNetworkPolicy',
      '10': 'networkpolicy'
    },
    {
      '1': 'workspace',
      '3': 254637495,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchWorkspace',
      '10': 'workspace'
    },
  ],
};

/// Descriptor for `OpenSearchIntegrationDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchIntegrationDetailsDescriptor = $convert.base64Decode(
    'ChxPcGVuU2VhcmNoSW50ZWdyYXRpb25EZXRhaWxzEkgKDGFjY2Vzc3BvbGljeRjOyJ+gASABKA'
    'syIC5sb2dzLk9wZW5TZWFyY2hEYXRhQWNjZXNzUG9saWN5UgxhY2Nlc3Nwb2xpY3kSQQoLYXBw'
    'bGljYXRpb24YxsSe+gEgASgLMhsubG9ncy5PcGVuU2VhcmNoQXBwbGljYXRpb25SC2FwcGxpY2'
    'F0aW9uEj4KCmNvbGxlY3Rpb24Ykom2/QEgASgLMhoubG9ncy5PcGVuU2VhcmNoQ29sbGVjdGlv'
    'blIKY29sbGVjdGlvbhI+CgpkYXRhc291cmNlGJnX76QBIAEoCzIaLmxvZ3MuT3BlblNlYXJjaE'
    'RhdGFTb3VyY2VSCmRhdGFzb3VyY2USTwoQZW5jcnlwdGlvbnBvbGljeRiXysghIAEoCzIgLmxv'
    'Z3MuT3BlblNlYXJjaEVuY3J5cHRpb25Qb2xpY3lSEGVuY3J5cHRpb25wb2xpY3kSTAoPbGlmZW'
    'N5Y2xlcG9saWN5GOLmqmQgASgLMh8ubG9ncy5PcGVuU2VhcmNoTGlmZWN5Y2xlUG9saWN5Ug9s'
    'aWZlY3ljbGVwb2xpY3kSRgoNbmV0d29ya3BvbGljeRigscNfIAEoCzIdLmxvZ3MuT3BlblNlYX'
    'JjaE5ldHdvcmtQb2xpY3lSDW5ldHdvcmtwb2xpY3kSOgoJd29ya3NwYWNlGLfrtXkgASgLMhku'
    'bG9ncy5PcGVuU2VhcmNoV29ya3NwYWNlUgl3b3Jrc3BhY2U=');

@$core.Deprecated('Use openSearchLifecyclePolicyDescriptor instead')
const OpenSearchLifecyclePolicy$json = {
  '1': 'OpenSearchLifecyclePolicy',
  '2': [
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `OpenSearchLifecyclePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchLifecyclePolicyDescriptor = $convert.base64Decode(
    'ChlPcGVuU2VhcmNoTGlmZWN5Y2xlUG9saWN5EiEKCnBvbGljeW5hbWUY3eLyNiABKAlSCnBvbG'
    'ljeW5hbWUSOgoGc3RhdHVzGPDvrdIBIAEoCzIeLmxvZ3MuT3BlblNlYXJjaFJlc291cmNlU3Rh'
    'dHVzUgZzdGF0dXM=');

@$core.Deprecated('Use openSearchNetworkPolicyDescriptor instead')
const OpenSearchNetworkPolicy$json = {
  '1': 'OpenSearchNetworkPolicy',
  '2': [
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `OpenSearchNetworkPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchNetworkPolicyDescriptor = $convert.base64Decode(
    'ChdPcGVuU2VhcmNoTmV0d29ya1BvbGljeRIhCgpwb2xpY3luYW1lGN3i8jYgASgJUgpwb2xpY3'
    'luYW1lEjoKBnN0YXR1cxjw763SASABKAsyHi5sb2dzLk9wZW5TZWFyY2hSZXNvdXJjZVN0YXR1'
    'c1IGc3RhdHVz');

@$core.Deprecated('Use openSearchResourceConfigDescriptor instead')
const OpenSearchResourceConfig$json = {
  '1': 'OpenSearchResourceConfig',
  '2': [
    {
      '1': 'applicationarn',
      '3': 230367005,
      '4': 1,
      '5': 9,
      '10': 'applicationarn'
    },
    {
      '1': 'dashboardviewerprincipals',
      '3': 157196451,
      '4': 3,
      '5': 9,
      '10': 'dashboardviewerprincipals'
    },
    {
      '1': 'datasourcerolearn',
      '3': 205136432,
      '4': 1,
      '5': 9,
      '10': 'datasourcerolearn'
    },
    {'1': 'kmskeyarn', '3': 341492497, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {
      '1': 'retentiondays',
      '3': 536209775,
      '4': 1,
      '5': 5,
      '10': 'retentiondays'
    },
  ],
};

/// Descriptor for `OpenSearchResourceConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchResourceConfigDescriptor = $convert.base64Decode(
    'ChhPcGVuU2VhcmNoUmVzb3VyY2VDb25maWcSKQoOYXBwbGljYXRpb25hcm4Ynb7sbSABKAlSDm'
    'FwcGxpY2F0aW9uYXJuEj8KGWRhc2hib2FyZHZpZXdlcnByaW5jaXBhbHMYo8H6SiADKAlSGWRh'
    'c2hib2FyZHZpZXdlcnByaW5jaXBhbHMSLwoRZGF0YXNvdXJjZXJvbGVhcm4YsMToYSABKAlSEW'
    'RhdGFzb3VyY2Vyb2xlYXJuEiAKCWttc2tleWFybhiRhuuiASABKAlSCWttc2tleWFybhIoCg1y'
    'ZXRlbnRpb25kYXlzGO/S1/8BIAEoBVINcmV0ZW50aW9uZGF5cw==');

@$core.Deprecated('Use openSearchResourceStatusDescriptor instead')
const OpenSearchResourceStatus$json = {
  '1': 'OpenSearchResourceStatus',
  '2': [
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.logs.OpenSearchResourceStatusType',
      '10': 'status'
    },
    {
      '1': 'statusmessage',
      '3': 474462255,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
  ],
};

/// Descriptor for `OpenSearchResourceStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchResourceStatusDescriptor = $convert.base64Decode(
    'ChhPcGVuU2VhcmNoUmVzb3VyY2VTdGF0dXMSPgoGc3RhdHVzGPDvrdIBIAEoDjIiLmxvZ3MuT3'
    'BlblNlYXJjaFJlc291cmNlU3RhdHVzVHlwZVIGc3RhdHVzEigKDXN0YXR1c21lc3NhZ2UYr/Ce'
    '4gEgASgJUg1zdGF0dXNtZXNzYWdl');

@$core.Deprecated('Use openSearchWorkspaceDescriptor instead')
const OpenSearchWorkspace$json = {
  '1': 'OpenSearchWorkspace',
  '2': [
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceStatus',
      '10': 'status'
    },
    {'1': 'workspaceid', '3': 98455084, '4': 1, '5': 9, '10': 'workspaceid'},
  ],
};

/// Descriptor for `OpenSearchWorkspace`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List openSearchWorkspaceDescriptor = $convert.base64Decode(
    'ChNPcGVuU2VhcmNoV29ya3NwYWNlEjoKBnN0YXR1cxjw763SASABKAsyHi5sb2dzLk9wZW5TZW'
    'FyY2hSZXNvdXJjZVN0YXR1c1IGc3RhdHVzEiMKC3dvcmtzcGFjZWlkGKyc+S4gASgJUgt3b3Jr'
    'c3BhY2VpZA==');

@$core.Deprecated('Use operationAbortedExceptionDescriptor instead')
const OperationAbortedException$json = {
  '1': 'OperationAbortedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OperationAbortedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List operationAbortedExceptionDescriptor =
    $convert.base64Decode(
        'ChlPcGVyYXRpb25BYm9ydGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use outputLogEventDescriptor instead')
const OutputLogEvent$json = {
  '1': 'OutputLogEvent',
  '2': [
    {
      '1': 'ingestiontime',
      '3': 179367957,
      '4': 1,
      '5': 3,
      '10': 'ingestiontime'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'timestamp', '3': 310629668, '4': 1, '5': 3, '10': 'timestamp'},
  ],
};

/// Descriptor for `OutputLogEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List outputLogEventDescriptor = $convert.base64Decode(
    'Cg5PdXRwdXRMb2dFdmVudBInCg1pbmdlc3Rpb250aW1lGJXgw1UgASgDUg1pbmdlc3Rpb250aW'
    '1lEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2USIAoJdGltZXN0YW1wGKSqj5QBIAEoA1IJ'
    'dGltZXN0YW1w');

@$core.Deprecated('Use parseCloudfrontDescriptor instead')
const ParseCloudfront$json = {
  '1': 'ParseCloudfront',
  '2': [
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParseCloudfront`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parseCloudfrontDescriptor = $convert.base64Decode(
    'Cg9QYXJzZUNsb3VkZnJvbnQSGgoGc291cmNlGNnTvN4BIAEoCVIGc291cmNl');

@$core.Deprecated('Use parseJSONDescriptor instead')
const ParseJSON$json = {
  '1': 'ParseJSON',
  '2': [
    {'1': 'destination', '3': 316564672, '4': 1, '5': 9, '10': 'destination'},
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParseJSON`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parseJSONDescriptor = $convert.base64Decode(
    'CglQYXJzZUpTT04SJAoLZGVzdGluYXRpb24YwMn5lgEgASgJUgtkZXN0aW5hdGlvbhIaCgZzb3'
    'VyY2UY2dO83gEgASgJUgZzb3VyY2U=');

@$core.Deprecated('Use parseKeyValueDescriptor instead')
const ParseKeyValue$json = {
  '1': 'ParseKeyValue',
  '2': [
    {'1': 'destination', '3': 316564672, '4': 1, '5': 9, '10': 'destination'},
    {
      '1': 'fielddelimiter',
      '3': 433214845,
      '4': 1,
      '5': 9,
      '10': 'fielddelimiter'
    },
    {'1': 'keyprefix', '3': 318004009, '4': 1, '5': 9, '10': 'keyprefix'},
    {
      '1': 'keyvaluedelimiter',
      '3': 186408507,
      '4': 1,
      '5': 9,
      '10': 'keyvaluedelimiter'
    },
    {
      '1': 'nonmatchvalue',
      '3': 229964203,
      '4': 1,
      '5': 9,
      '10': 'nonmatchvalue'
    },
    {
      '1': 'overwriteifexists',
      '3': 230880030,
      '4': 1,
      '5': 8,
      '10': 'overwriteifexists'
    },
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParseKeyValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parseKeyValueDescriptor = $convert.base64Decode(
    'Cg1QYXJzZUtleVZhbHVlEiQKC2Rlc3RpbmF0aW9uGMDJ+ZYBIAEoCVILZGVzdGluYXRpb24SKg'
    'oOZmllbGRkZWxpbWl0ZXIY/arJzgEgASgJUg5maWVsZGRlbGltaXRlchIgCglrZXlwcmVmaXgY'
    'qbbRlwEgASgJUglrZXlwcmVmaXgSLwoRa2V5dmFsdWVkZWxpbWl0ZXIYu7zxWCABKAlSEWtleX'
    'ZhbHVlZGVsaW1pdGVyEicKDW5vbm1hdGNodmFsdWUYq/PTbSABKAlSDW5vbm1hdGNodmFsdWUS'
    'LwoRb3ZlcndyaXRlaWZleGlzdHMYnuaLbiABKAhSEW92ZXJ3cml0ZWlmZXhpc3RzEhoKBnNvdX'
    'JjZRjZ07zeASABKAlSBnNvdXJjZQ==');

@$core.Deprecated('Use parsePostgresDescriptor instead')
const ParsePostgres$json = {
  '1': 'ParsePostgres',
  '2': [
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParsePostgres`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parsePostgresDescriptor = $convert.base64Decode(
    'Cg1QYXJzZVBvc3RncmVzEhoKBnNvdXJjZRjZ07zeASABKAlSBnNvdXJjZQ==');

@$core.Deprecated('Use parseRoute53Descriptor instead')
const ParseRoute53$json = {
  '1': 'ParseRoute53',
  '2': [
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParseRoute53`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parseRoute53Descriptor = $convert
    .base64Decode('CgxQYXJzZVJvdXRlNTMSGgoGc291cmNlGNnTvN4BIAEoCVIGc291cmNl');

@$core.Deprecated('Use parseToOCSFDescriptor instead')
const ParseToOCSF$json = {
  '1': 'ParseToOCSF',
  '2': [
    {
      '1': 'eventsource',
      '3': 260249947,
      '4': 1,
      '5': 14,
      '6': '.logs.EventSource',
      '10': 'eventsource'
    },
    {
      '1': 'mappingversion',
      '3': 29454072,
      '4': 1,
      '5': 9,
      '10': 'mappingversion'
    },
    {
      '1': 'ocsfversion',
      '3': 3840275,
      '4': 1,
      '5': 14,
      '6': '.logs.OCSFVersion',
      '10': 'ocsfversion'
    },
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParseToOCSF`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parseToOCSFDescriptor = $convert.base64Decode(
    'CgtQYXJzZVRvT0NTRhI2CgtldmVudHNvdXJjZRjbsox8IAEoDjIRLmxvZ3MuRXZlbnRTb3VyY2'
    'VSC2V2ZW50c291cmNlEikKDm1hcHBpbmd2ZXJzaW9uGPjdhQ4gASgJUg5tYXBwaW5ndmVyc2lv'
    'bhI2CgtvY3NmdmVyc2lvbhiTsuoBIAEoDjIRLmxvZ3MuT0NTRlZlcnNpb25SC29jc2Z2ZXJzaW'
    '9uEhoKBnNvdXJjZRjZ07zeASABKAlSBnNvdXJjZQ==');

@$core.Deprecated('Use parseVPCDescriptor instead')
const ParseVPC$json = {
  '1': 'ParseVPC',
  '2': [
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParseVPC`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parseVPCDescriptor = $convert
    .base64Decode('CghQYXJzZVZQQxIaCgZzb3VyY2UY2dO83gEgASgJUgZzb3VyY2U=');

@$core.Deprecated('Use parseWAFDescriptor instead')
const ParseWAF$json = {
  '1': 'ParseWAF',
  '2': [
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `ParseWAF`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parseWAFDescriptor = $convert
    .base64Decode('CghQYXJzZVdBRhIaCgZzb3VyY2UY2dO83gEgASgJUgZzb3VyY2U=');

@$core.Deprecated('Use patternTokenDescriptor instead')
const PatternToken$json = {
  '1': 'PatternToken',
  '2': [
    {
      '1': 'dynamictokenposition',
      '3': 9081241,
      '4': 1,
      '5': 5,
      '10': 'dynamictokenposition'
    },
    {
      '1': 'enumerations',
      '3': 14171326,
      '4': 3,
      '5': 11,
      '6': '.logs.PatternToken.EnumerationsEntry',
      '10': 'enumerations'
    },
    {
      '1': 'inferredtokenname',
      '3': 57553257,
      '4': 1,
      '5': 9,
      '10': 'inferredtokenname'
    },
    {'1': 'isdynamic', '3': 17462495, '4': 1, '5': 8, '10': 'isdynamic'},
    {'1': 'tokenstring', '3': 180629740, '4': 1, '5': 9, '10': 'tokenstring'},
  ],
  '3': [PatternToken_EnumerationsEntry$json],
};

@$core.Deprecated('Use patternTokenDescriptor instead')
const PatternToken_EnumerationsEntry$json = {
  '1': 'EnumerationsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 3, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PatternToken`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patternTokenDescriptor = $convert.base64Decode(
    'CgxQYXR0ZXJuVG9rZW4SNQoUZHluYW1pY3Rva2VucG9zaXRpb24YmaOqBCABKAVSFGR5bmFtaW'
    'N0b2tlbnBvc2l0aW9uEksKDGVudW1lcmF0aW9ucxi++eAGIAMoCzIkLmxvZ3MuUGF0dGVyblRv'
    'a2VuLkVudW1lcmF0aW9uc0VudHJ5UgxlbnVtZXJhdGlvbnMSLwoRaW5mZXJyZWR0b2tlbm5hbW'
    'UY6eK4GyABKAlSEWluZmVycmVkdG9rZW5uYW1lEh8KCWlzZHluYW1pYxjf6akIIAEoCFIJaXNk'
    'eW5hbWljEiMKC3Rva2Vuc3RyaW5nGOzhkFYgASgJUgt0b2tlbnN0cmluZxo/ChFFbnVtZXJhdG'
    'lvbnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoA1IFdmFsdWU6AjgB');

@$core.Deprecated('Use policyDescriptor instead')
const Policy$json = {
  '1': 'Policy',
  '2': [
    {
      '1': 'deliverydestinationpolicy',
      '3': 413241666,
      '4': 1,
      '5': 9,
      '10': 'deliverydestinationpolicy'
    },
  ],
};

/// Descriptor for `Policy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyDescriptor = $convert.base64Decode(
    'CgZQb2xpY3kSQAoZZGVsaXZlcnlkZXN0aW5hdGlvbnBvbGljeRjCoobFASABKAlSGWRlbGl2ZX'
    'J5ZGVzdGluYXRpb25wb2xpY3k=');

@$core.Deprecated('Use processorDescriptor instead')
const Processor$json = {
  '1': 'Processor',
  '2': [
    {
      '1': 'addkeys',
      '3': 300392065,
      '4': 1,
      '5': 11,
      '6': '.logs.AddKeys',
      '10': 'addkeys'
    },
    {
      '1': 'copyvalue',
      '3': 431370302,
      '4': 1,
      '5': 11,
      '6': '.logs.CopyValue',
      '10': 'copyvalue'
    },
    {
      '1': 'csv',
      '3': 407316064,
      '4': 1,
      '5': 11,
      '6': '.logs.CSV',
      '10': 'csv'
    },
    {
      '1': 'datetimeconverter',
      '3': 519051029,
      '4': 1,
      '5': 11,
      '6': '.logs.DateTimeConverter',
      '10': 'datetimeconverter'
    },
    {
      '1': 'deletekeys',
      '3': 217334743,
      '4': 1,
      '5': 11,
      '6': '.logs.DeleteKeys',
      '10': 'deletekeys'
    },
    {
      '1': 'grok',
      '3': 8553927,
      '4': 1,
      '5': 11,
      '6': '.logs.Grok',
      '10': 'grok'
    },
    {
      '1': 'listtomap',
      '3': 329538039,
      '4': 1,
      '5': 11,
      '6': '.logs.ListToMap',
      '10': 'listtomap'
    },
    {
      '1': 'lowercasestring',
      '3': 442295000,
      '4': 1,
      '5': 11,
      '6': '.logs.LowerCaseString',
      '10': 'lowercasestring'
    },
    {
      '1': 'movekeys',
      '3': 293478497,
      '4': 1,
      '5': 11,
      '6': '.logs.MoveKeys',
      '10': 'movekeys'
    },
    {
      '1': 'parsecloudfront',
      '3': 439788377,
      '4': 1,
      '5': 11,
      '6': '.logs.ParseCloudfront',
      '10': 'parsecloudfront'
    },
    {
      '1': 'parsejson',
      '3': 93591791,
      '4': 1,
      '5': 11,
      '6': '.logs.ParseJSON',
      '10': 'parsejson'
    },
    {
      '1': 'parsekeyvalue',
      '3': 415670841,
      '4': 1,
      '5': 11,
      '6': '.logs.ParseKeyValue',
      '10': 'parsekeyvalue'
    },
    {
      '1': 'parsepostgres',
      '3': 349269892,
      '4': 1,
      '5': 11,
      '6': '.logs.ParsePostgres',
      '10': 'parsepostgres'
    },
    {
      '1': 'parseroute53',
      '3': 448861006,
      '4': 1,
      '5': 11,
      '6': '.logs.ParseRoute53',
      '10': 'parseroute53'
    },
    {
      '1': 'parsetoocsf',
      '3': 443672911,
      '4': 1,
      '5': 11,
      '6': '.logs.ParseToOCSF',
      '10': 'parsetoocsf'
    },
    {
      '1': 'parsevpc',
      '3': 346687996,
      '4': 1,
      '5': 11,
      '6': '.logs.ParseVPC',
      '10': 'parsevpc'
    },
    {
      '1': 'parsewaf',
      '3': 342220243,
      '4': 1,
      '5': 11,
      '6': '.logs.ParseWAF',
      '10': 'parsewaf'
    },
    {
      '1': 'renamekeys',
      '3': 11638236,
      '4': 1,
      '5': 11,
      '6': '.logs.RenameKeys',
      '10': 'renamekeys'
    },
    {
      '1': 'splitstring',
      '3': 417544327,
      '4': 1,
      '5': 11,
      '6': '.logs.SplitString',
      '10': 'splitstring'
    },
    {
      '1': 'substitutestring',
      '3': 130738861,
      '4': 1,
      '5': 11,
      '6': '.logs.SubstituteString',
      '10': 'substitutestring'
    },
    {
      '1': 'trimstring',
      '3': 18323845,
      '4': 1,
      '5': 11,
      '6': '.logs.TrimString',
      '10': 'trimstring'
    },
    {
      '1': 'typeconverter',
      '3': 487701320,
      '4': 1,
      '5': 11,
      '6': '.logs.TypeConverter',
      '10': 'typeconverter'
    },
    {
      '1': 'uppercasestring',
      '3': 96994447,
      '4': 1,
      '5': 11,
      '6': '.logs.UpperCaseString',
      '10': 'uppercasestring'
    },
  ],
};

/// Descriptor for `Processor`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List processorDescriptor = $convert.base64Decode(
    'CglQcm9jZXNzb3ISKwoHYWRka2V5cxiBvZ6PASABKAsyDS5sb2dzLkFkZEtleXNSB2FkZGtleX'
    'MSMQoJY29weXZhbHVlGL7g2M0BIAEoCzIPLmxvZ3MuQ29weVZhbHVlUgljb3B5dmFsdWUSHwoD'
    'Y3N2GODMnMIBIAEoCzIJLmxvZ3MuQ1NWUgNjc3YSSQoRZGF0ZXRpbWVjb252ZXJ0ZXIYla7A9w'
    'EgASgLMhcubG9ncy5EYXRlVGltZUNvbnZlcnRlclIRZGF0ZXRpbWVjb252ZXJ0ZXISMwoKZGVs'
    'ZXRla2V5cxjXh9FnIAEoCzIQLmxvZ3MuRGVsZXRlS2V5c1IKZGVsZXRla2V5cxIhCgRncm9rGM'
    'eLigQgASgLMgoubG9ncy5Hcm9rUgRncm9rEjEKCWxpc3R0b21hcBj3s5GdASABKAsyDy5sb2dz'
    'Lkxpc3RUb01hcFIJbGlzdHRvbWFwEkMKD2xvd2VyY2FzZXN0cmluZxjYxfPSASABKAsyFS5sb2'
    'dzLkxvd2VyQ2FzZVN0cmluZ1IPbG93ZXJjYXNlc3RyaW5nEi4KCG1vdmVrZXlzGOHA+IsBIAEo'
    'CzIOLmxvZ3MuTW92ZUtleXNSCG1vdmVrZXlzEkMKD3BhcnNlY2xvdWRmcm9udBjZxtrRASABKA'
    'syFS5sb2dzLlBhcnNlQ2xvdWRmcm9udFIPcGFyc2VjbG91ZGZyb250EjAKCXBhcnNlanNvbhjv'
    'sdAsIAEoCzIPLmxvZ3MuUGFyc2VKU09OUglwYXJzZWpzb24SPQoNcGFyc2VrZXl2YWx1ZRi5xJ'
    'rGASABKAsyEy5sb2dzLlBhcnNlS2V5VmFsdWVSDXBhcnNla2V5dmFsdWUSPQoNcGFyc2Vwb3N0'
    'Z3JlcxiE38WmASABKAsyEy5sb2dzLlBhcnNlUG9zdGdyZXNSDXBhcnNlcG9zdGdyZXMSOgoMcG'
    'Fyc2Vyb3V0ZTUzGM6mhNYBIAEoCzISLmxvZ3MuUGFyc2VSb3V0ZTUzUgxwYXJzZXJvdXRlNTMS'
    'NwoLcGFyc2V0b29jc2YYz9LH0wEgASgLMhEubG9ncy5QYXJzZVRvT0NTRlILcGFyc2V0b29jc2'
    'YSLgoIcGFyc2V2cGMY/JOopQEgASgLMg4ubG9ncy5QYXJzZVZQQ1IIcGFyc2V2cGMSLgoIcGFy'
    'c2V3YWYY07uXowEgASgLMg4ubG9ncy5QYXJzZVdBRlIIcGFyc2V3YWYSMwoKcmVuYW1la2V5cx'
    'jcq8YFIAEoCzIQLmxvZ3MuUmVuYW1lS2V5c1IKcmVuYW1la2V5cxI3CgtzcGxpdHN0cmluZxiH'
    '8YzHASABKAsyES5sb2dzLlNwbGl0U3RyaW5nUgtzcGxpdHN0cmluZxJFChBzdWJzdGl0dXRlc3'
    'RyaW5nGK3Vqz4gASgLMhYubG9ncy5TdWJzdGl0dXRlU3RyaW5nUhBzdWJzdGl0dXRlc3RyaW5n'
    'EjMKCnRyaW1zdHJpbmcYhbPeCCABKAsyEC5sb2dzLlRyaW1TdHJpbmdSCnRyaW1zdHJpbmcSPQ'
    'oNdHlwZWNvbnZlcnRlchjI9sboASABKAsyEy5sb2dzLlR5cGVDb252ZXJ0ZXJSDXR5cGVjb252'
    'ZXJ0ZXISQgoPdXBwZXJjYXNlc3RyaW5nGI+JoC4gASgLMhUubG9ncy5VcHBlckNhc2VTdHJpbm'
    'dSD3VwcGVyY2FzZXN0cmluZw==');

@$core.Deprecated('Use putAccountPolicyRequestDescriptor instead')
const PutAccountPolicyRequest$json = {
  '1': 'PutAccountPolicyRequest',
  '2': [
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'policytype',
      '3': 319277736,
      '4': 1,
      '5': 14,
      '6': '.logs.PolicyType',
      '10': 'policytype'
    },
    {
      '1': 'scope',
      '3': 506131436,
      '4': 1,
      '5': 14,
      '6': '.logs.Scope',
      '10': 'scope'
    },
    {
      '1': 'selectioncriteria',
      '3': 145052429,
      '4': 1,
      '5': 9,
      '10': 'selectioncriteria'
    },
  ],
};

/// Descriptor for `PutAccountPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountPolicyRequestDescriptor = $convert.base64Decode(
    'ChdQdXRBY2NvdW50UG9saWN5UmVxdWVzdBIpCg5wb2xpY3lkb2N1bWVudBjr6fZUIAEoCVIOcG'
    '9saWN5ZG9jdW1lbnQSIQoKcG9saWN5bmFtZRjd4vI2IAEoCVIKcG9saWN5bmFtZRI0Cgpwb2xp'
    'Y3l0eXBlGKiVn5gBIAEoDjIQLmxvZ3MuUG9saWN5VHlwZVIKcG9saWN5dHlwZRIlCgVzY29wZR'
    'js56vxASABKA4yCy5sb2dzLlNjb3BlUgVzY29wZRIvChFzZWxlY3Rpb25jcml0ZXJpYRiNppVF'
    'IAEoCVIRc2VsZWN0aW9uY3JpdGVyaWE=');

@$core.Deprecated('Use putAccountPolicyResponseDescriptor instead')
const PutAccountPolicyResponse$json = {
  '1': 'PutAccountPolicyResponse',
  '2': [
    {
      '1': 'accountpolicy',
      '3': 186824791,
      '4': 1,
      '5': 11,
      '6': '.logs.AccountPolicy',
      '10': 'accountpolicy'
    },
  ],
};

/// Descriptor for `PutAccountPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountPolicyResponseDescriptor =
    $convert.base64Decode(
        'ChhQdXRBY2NvdW50UG9saWN5UmVzcG9uc2USPAoNYWNjb3VudHBvbGljeRjX8IpZIAEoCzITLm'
        'xvZ3MuQWNjb3VudFBvbGljeVINYWNjb3VudHBvbGljeQ==');

@$core.Deprecated('Use putBearerTokenAuthenticationRequestDescriptor instead')
const PutBearerTokenAuthenticationRequest$json = {
  '1': 'PutBearerTokenAuthenticationRequest',
  '2': [
    {
      '1': 'bearertokenauthenticationenabled',
      '3': 513390155,
      '4': 1,
      '5': 8,
      '10': 'bearertokenauthenticationenabled'
    },
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
  ],
};

/// Descriptor for `PutBearerTokenAuthenticationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBearerTokenAuthenticationRequestDescriptor =
    $convert.base64Decode(
        'CiNQdXRCZWFyZXJUb2tlbkF1dGhlbnRpY2F0aW9uUmVxdWVzdBJOCiBiZWFyZXJ0b2tlbmF1dG'
        'hlbnRpY2F0aW9uZW5hYmxlZBjL7Ob0ASABKAhSIGJlYXJlcnRva2VuYXV0aGVudGljYXRpb25l'
        'bmFibGVkEjIKEmxvZ2dyb3VwaWRlbnRpZmllchjcz6XfASABKAlSEmxvZ2dyb3VwaWRlbnRpZm'
        'llcg==');

@$core.Deprecated('Use putDataProtectionPolicyRequestDescriptor instead')
const PutDataProtectionPolicyRequest$json = {
  '1': 'PutDataProtectionPolicyRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
  ],
};

/// Descriptor for `PutDataProtectionPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDataProtectionPolicyRequestDescriptor =
    $convert.base64Decode(
        'Ch5QdXREYXRhUHJvdGVjdGlvblBvbGljeVJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGN'
        'zPpd8BIAEoCVISbG9nZ3JvdXBpZGVudGlmaWVyEikKDnBvbGljeWRvY3VtZW50GOvp9lQgASgJ'
        'Ug5wb2xpY3lkb2N1bWVudA==');

@$core.Deprecated('Use putDataProtectionPolicyResponseDescriptor instead')
const PutDataProtectionPolicyResponse$json = {
  '1': 'PutDataProtectionPolicyResponse',
  '2': [
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
  ],
};

/// Descriptor for `PutDataProtectionPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDataProtectionPolicyResponseDescriptor =
    $convert.base64Decode(
        'Ch9QdXREYXRhUHJvdGVjdGlvblBvbGljeVJlc3BvbnNlEiwKD2xhc3R1cGRhdGVkdGltZRj2t5'
        'W5ASABKANSD2xhc3R1cGRhdGVkdGltZRIyChJsb2dncm91cGlkZW50aWZpZXIY3M+l3wEgASgJ'
        'UhJsb2dncm91cGlkZW50aWZpZXISKQoOcG9saWN5ZG9jdW1lbnQY6+n2VCABKAlSDnBvbGljeW'
        'RvY3VtZW50');

@$core.Deprecated('Use putDeliveryDestinationPolicyRequestDescriptor instead')
const PutDeliveryDestinationPolicyRequest$json = {
  '1': 'PutDeliveryDestinationPolicyRequest',
  '2': [
    {
      '1': 'deliverydestinationname',
      '3': 331358541,
      '4': 1,
      '5': 9,
      '10': 'deliverydestinationname'
    },
    {
      '1': 'deliverydestinationpolicy',
      '3': 413241666,
      '4': 1,
      '5': 9,
      '10': 'deliverydestinationpolicy'
    },
  ],
};

/// Descriptor for `PutDeliveryDestinationPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliveryDestinationPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiNQdXREZWxpdmVyeURlc3RpbmF0aW9uUG9saWN5UmVxdWVzdBI8ChdkZWxpdmVyeWRlc3Rpbm'
        'F0aW9ubmFtZRjNwoCeASABKAlSF2RlbGl2ZXJ5ZGVzdGluYXRpb25uYW1lEkAKGWRlbGl2ZXJ5'
        'ZGVzdGluYXRpb25wb2xpY3kYwqKGxQEgASgJUhlkZWxpdmVyeWRlc3RpbmF0aW9ucG9saWN5');

@$core.Deprecated('Use putDeliveryDestinationPolicyResponseDescriptor instead')
const PutDeliveryDestinationPolicyResponse$json = {
  '1': 'PutDeliveryDestinationPolicyResponse',
  '2': [
    {
      '1': 'policy',
      '3': 247528064,
      '4': 1,
      '5': 11,
      '6': '.logs.Policy',
      '10': 'policy'
    },
  ],
};

/// Descriptor for `PutDeliveryDestinationPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliveryDestinationPolicyResponseDescriptor =
    $convert.base64Decode(
        'CiRQdXREZWxpdmVyeURlc3RpbmF0aW9uUG9saWN5UmVzcG9uc2USJwoGcG9saWN5GID1g3YgAS'
        'gLMgwubG9ncy5Qb2xpY3lSBnBvbGljeQ==');

@$core.Deprecated('Use putDeliveryDestinationRequestDescriptor instead')
const PutDeliveryDestinationRequest$json = {
  '1': 'PutDeliveryDestinationRequest',
  '2': [
    {
      '1': 'deliverydestinationconfiguration',
      '3': 256504432,
      '4': 1,
      '5': 11,
      '6': '.logs.DeliveryDestinationConfiguration',
      '10': 'deliverydestinationconfiguration'
    },
    {
      '1': 'deliverydestinationtype',
      '3': 33253176,
      '4': 1,
      '5': 14,
      '6': '.logs.DeliveryDestinationType',
      '10': 'deliverydestinationtype'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputformat',
      '3': 217347480,
      '4': 1,
      '5': 14,
      '6': '.logs.OutputFormat',
      '10': 'outputformat'
    },
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.PutDeliveryDestinationRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [PutDeliveryDestinationRequest_TagsEntry$json],
};

@$core.Deprecated('Use putDeliveryDestinationRequestDescriptor instead')
const PutDeliveryDestinationRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutDeliveryDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliveryDestinationRequestDescriptor = $convert.base64Decode(
    'Ch1QdXREZWxpdmVyeURlc3RpbmF0aW9uUmVxdWVzdBJ1CiBkZWxpdmVyeWRlc3RpbmF0aW9uY2'
    '9uZmlndXJhdGlvbhjw5Kd6IAEoCzImLmxvZ3MuRGVsaXZlcnlEZXN0aW5hdGlvbkNvbmZpZ3Vy'
    'YXRpb25SIGRlbGl2ZXJ5ZGVzdGluYXRpb25jb25maWd1cmF0aW9uEloKF2RlbGl2ZXJ5ZGVzdG'
    'luYXRpb250eXBlGLjO7Q8gASgOMh0ubG9ncy5EZWxpdmVyeURlc3RpbmF0aW9uVHlwZVIXZGVs'
    'aXZlcnlkZXN0aW5hdGlvbnR5cGUSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZRI5CgxvdXRwdXRmb3'
    'JtYXQYmOvRZyABKA4yEi5sb2dzLk91dHB1dEZvcm1hdFIMb3V0cHV0Zm9ybWF0EkUKBHRhZ3MY'
    'odfboAEgAygLMi0ubG9ncy5QdXREZWxpdmVyeURlc3RpbmF0aW9uUmVxdWVzdC5UYWdzRW50cn'
    'lSBHRhZ3MaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2'
    'YWx1ZToCOAE=');

@$core.Deprecated('Use putDeliveryDestinationResponseDescriptor instead')
const PutDeliveryDestinationResponse$json = {
  '1': 'PutDeliveryDestinationResponse',
  '2': [
    {
      '1': 'deliverydestination',
      '3': 103332720,
      '4': 1,
      '5': 11,
      '6': '.logs.DeliveryDestination',
      '10': 'deliverydestination'
    },
  ],
};

/// Descriptor for `PutDeliveryDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliveryDestinationResponseDescriptor =
    $convert.base64Decode(
        'Ch5QdXREZWxpdmVyeURlc3RpbmF0aW9uUmVzcG9uc2USTgoTZGVsaXZlcnlkZXN0aW5hdGlvbh'
        'jw9qIxIAEoCzIZLmxvZ3MuRGVsaXZlcnlEZXN0aW5hdGlvblITZGVsaXZlcnlkZXN0aW5hdGlv'
        'bg==');

@$core.Deprecated('Use putDeliverySourceRequestDescriptor instead')
const PutDeliverySourceRequest$json = {
  '1': 'PutDeliverySourceRequest',
  '2': [
    {'1': 'logtype', '3': 257838938, '4': 1, '5': 9, '10': 'logtype'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.PutDeliverySourceRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [PutDeliverySourceRequest_TagsEntry$json],
};

@$core.Deprecated('Use putDeliverySourceRequestDescriptor instead')
const PutDeliverySourceRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutDeliverySourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliverySourceRequestDescriptor = $convert.base64Decode(
    'ChhQdXREZWxpdmVyeVNvdXJjZVJlcXVlc3QSGwoHbG9ndHlwZRjanvl6IAEoCVIHbG9ndHlwZR'
    'IVCgRuYW1lGOf75mkgASgJUgRuYW1lEiMKC3Jlc291cmNlYXJuGM3MqiAgASgJUgtyZXNvdXJj'
    'ZWFybhJACgR0YWdzGKHX26ABIAMoCzIoLmxvZ3MuUHV0RGVsaXZlcnlTb3VyY2VSZXF1ZXN0Ll'
    'RhZ3NFbnRyeVIEdGFncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUY'
    'AiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use putDeliverySourceResponseDescriptor instead')
const PutDeliverySourceResponse$json = {
  '1': 'PutDeliverySourceResponse',
  '2': [
    {
      '1': 'deliverysource',
      '3': 1553897,
      '4': 1,
      '5': 11,
      '6': '.logs.DeliverySource',
      '10': 'deliverysource'
    },
  ],
};

/// Descriptor for `PutDeliverySourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliverySourceResponseDescriptor =
    $convert.base64Decode(
        'ChlQdXREZWxpdmVyeVNvdXJjZVJlc3BvbnNlEj4KDmRlbGl2ZXJ5c291cmNlGOnrXiABKAsyFC'
        '5sb2dzLkRlbGl2ZXJ5U291cmNlUg5kZWxpdmVyeXNvdXJjZQ==');

@$core.Deprecated('Use putDestinationPolicyRequestDescriptor instead')
const PutDestinationPolicyRequest$json = {
  '1': 'PutDestinationPolicyRequest',
  '2': [
    {'1': 'accesspolicy', '3': 336061518, '4': 1, '5': 9, '10': 'accesspolicy'},
    {
      '1': 'destinationname',
      '3': 284844189,
      '4': 1,
      '5': 9,
      '10': 'destinationname'
    },
    {'1': 'forceupdate', '3': 21245222, '4': 1, '5': 8, '10': 'forceupdate'},
  ],
};

/// Descriptor for `PutDestinationPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDestinationPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChtQdXREZXN0aW5hdGlvblBvbGljeVJlcXVlc3QSJgoMYWNjZXNzcG9saWN5GM7In6ABIAEoCV'
        'IMYWNjZXNzcG9saWN5EiwKD2Rlc3RpbmF0aW9ubmFtZRidwemHASABKAlSD2Rlc3RpbmF0aW9u'
        'bmFtZRIjCgtmb3JjZXVwZGF0ZRim2pAKIAEoCFILZm9yY2V1cGRhdGU=');

@$core.Deprecated('Use putDestinationRequestDescriptor instead')
const PutDestinationRequest$json = {
  '1': 'PutDestinationRequest',
  '2': [
    {
      '1': 'destinationname',
      '3': 284844189,
      '4': 1,
      '5': 9,
      '10': 'destinationname'
    },
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.PutDestinationRequest.TagsEntry',
      '10': 'tags'
    },
    {'1': 'targetarn', '3': 367964720, '4': 1, '5': 9, '10': 'targetarn'},
  ],
  '3': [PutDestinationRequest_TagsEntry$json],
};

@$core.Deprecated('Use putDestinationRequestDescriptor instead')
const PutDestinationRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDestinationRequestDescriptor = $convert.base64Decode(
    'ChVQdXREZXN0aW5hdGlvblJlcXVlc3QSLAoPZGVzdGluYXRpb25uYW1lGJ3B6YcBIAEoCVIPZG'
    'VzdGluYXRpb25uYW1lEhsKB3JvbGVhcm4YoZeJUSABKAlSB3JvbGVhcm4SPQoEdGFncxih19ug'
    'ASADKAsyJS5sb2dzLlB1dERlc3RpbmF0aW9uUmVxdWVzdC5UYWdzRW50cnlSBHRhZ3MSIAoJdG'
    'FyZ2V0YXJuGLDkuq8BIAEoCVIJdGFyZ2V0YXJuGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlS'
    'A2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use putDestinationResponseDescriptor instead')
const PutDestinationResponse$json = {
  '1': 'PutDestinationResponse',
  '2': [
    {
      '1': 'destination',
      '3': 316564672,
      '4': 1,
      '5': 11,
      '6': '.logs.Destination',
      '10': 'destination'
    },
  ],
};

/// Descriptor for `PutDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDestinationResponseDescriptor =
    $convert.base64Decode(
        'ChZQdXREZXN0aW5hdGlvblJlc3BvbnNlEjcKC2Rlc3RpbmF0aW9uGMDJ+ZYBIAEoCzIRLmxvZ3'
        'MuRGVzdGluYXRpb25SC2Rlc3RpbmF0aW9u');

@$core.Deprecated('Use putIndexPolicyRequestDescriptor instead')
const PutIndexPolicyRequest$json = {
  '1': 'PutIndexPolicyRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
  ],
};

/// Descriptor for `PutIndexPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putIndexPolicyRequestDescriptor = $convert.base64Decode(
    'ChVQdXRJbmRleFBvbGljeVJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIAEoCV'
    'ISbG9nZ3JvdXBpZGVudGlmaWVyEikKDnBvbGljeWRvY3VtZW50GOvp9lQgASgJUg5wb2xpY3lk'
    'b2N1bWVudA==');

@$core.Deprecated('Use putIndexPolicyResponseDescriptor instead')
const PutIndexPolicyResponse$json = {
  '1': 'PutIndexPolicyResponse',
  '2': [
    {
      '1': 'indexpolicy',
      '3': 373363214,
      '4': 1,
      '5': 11,
      '6': '.logs.IndexPolicy',
      '10': 'indexpolicy'
    },
  ],
};

/// Descriptor for `PutIndexPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putIndexPolicyResponseDescriptor =
    $convert.base64Decode(
        'ChZQdXRJbmRleFBvbGljeVJlc3BvbnNlEjcKC2luZGV4cG9saWN5GI6khLIBIAEoCzIRLmxvZ3'
        'MuSW5kZXhQb2xpY3lSC2luZGV4cG9saWN5');

@$core.Deprecated('Use putIntegrationRequestDescriptor instead')
const PutIntegrationRequest$json = {
  '1': 'PutIntegrationRequest',
  '2': [
    {
      '1': 'integrationname',
      '3': 183183535,
      '4': 1,
      '5': 9,
      '10': 'integrationname'
    },
    {
      '1': 'integrationtype',
      '3': 307087270,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationType',
      '10': 'integrationtype'
    },
    {
      '1': 'resourceconfig',
      '3': 462395352,
      '4': 1,
      '5': 11,
      '6': '.logs.ResourceConfig',
      '10': 'resourceconfig'
    },
  ],
};

/// Descriptor for `PutIntegrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putIntegrationRequestDescriptor = $convert.base64Decode(
    'ChVQdXRJbnRlZ3JhdGlvblJlcXVlc3QSKwoPaW50ZWdyYXRpb25uYW1lGK/RrFcgASgJUg9pbn'
    'RlZ3JhdGlvbm5hbWUSQwoPaW50ZWdyYXRpb250eXBlGKaPt5IBIAEoDjIVLmxvZ3MuSW50ZWdy'
    'YXRpb25UeXBlUg9pbnRlZ3JhdGlvbnR5cGUSQAoOcmVzb3VyY2Vjb25maWcY2K++3AEgASgLMh'
    'QubG9ncy5SZXNvdXJjZUNvbmZpZ1IOcmVzb3VyY2Vjb25maWc=');

@$core.Deprecated('Use putIntegrationResponseDescriptor instead')
const PutIntegrationResponse$json = {
  '1': 'PutIntegrationResponse',
  '2': [
    {
      '1': 'integrationname',
      '3': 183183535,
      '4': 1,
      '5': 9,
      '10': 'integrationname'
    },
    {
      '1': 'integrationstatus',
      '3': 12110360,
      '4': 1,
      '5': 14,
      '6': '.logs.IntegrationStatus',
      '10': 'integrationstatus'
    },
  ],
};

/// Descriptor for `PutIntegrationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putIntegrationResponseDescriptor = $convert.base64Decode(
    'ChZQdXRJbnRlZ3JhdGlvblJlc3BvbnNlEisKD2ludGVncmF0aW9ubmFtZRiv0axXIAEoCVIPaW'
    '50ZWdyYXRpb25uYW1lEkgKEWludGVncmF0aW9uc3RhdHVzGJiU4wUgASgOMhcubG9ncy5JbnRl'
    'Z3JhdGlvblN0YXR1c1IRaW50ZWdyYXRpb25zdGF0dXM=');

@$core.Deprecated('Use putLogEventsRequestDescriptor instead')
const PutLogEventsRequest$json = {
  '1': 'PutLogEventsRequest',
  '2': [
    {
      '1': 'entity',
      '3': 322958811,
      '4': 1,
      '5': 11,
      '6': '.logs.Entity',
      '10': 'entity'
    },
    {
      '1': 'logevents',
      '3': 220433545,
      '4': 3,
      '5': 11,
      '6': '.logs.InputLogEvent',
      '10': 'logevents'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
    {'1': 'sequencetoken', '3': 2309662, '4': 1, '5': 9, '10': 'sequencetoken'},
  ],
};

/// Descriptor for `PutLogEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putLogEventsRequestDescriptor = $convert.base64Decode(
    'ChNQdXRMb2dFdmVudHNSZXF1ZXN0EigKBmVudGl0eRjb6/+ZASABKAsyDC5sb2dzLkVudGl0eV'
    'IGZW50aXR5EjQKCWxvZ2V2ZW50cxiJmY5pIAMoCzITLmxvZ3MuSW5wdXRMb2dFdmVudFIJbG9n'
    'ZXZlbnRzEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW1lEigKDWxvZ3N0cm'
    'VhbW5hbWUYo/fu0AEgASgJUg1sb2dzdHJlYW1uYW1lEicKDXNlcXVlbmNldG9rZW4YnvyMASAB'
    'KAlSDXNlcXVlbmNldG9rZW4=');

@$core.Deprecated('Use putLogEventsResponseDescriptor instead')
const PutLogEventsResponse$json = {
  '1': 'PutLogEventsResponse',
  '2': [
    {
      '1': 'nextsequencetoken',
      '3': 211087843,
      '4': 1,
      '5': 9,
      '10': 'nextsequencetoken'
    },
    {
      '1': 'rejectedentityinfo',
      '3': 4125469,
      '4': 1,
      '5': 11,
      '6': '.logs.RejectedEntityInfo',
      '10': 'rejectedentityinfo'
    },
    {
      '1': 'rejectedlogeventsinfo',
      '3': 417437083,
      '4': 1,
      '5': 11,
      '6': '.logs.RejectedLogEventsInfo',
      '10': 'rejectedlogeventsinfo'
    },
  ],
};

/// Descriptor for `PutLogEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putLogEventsResponseDescriptor = $convert.base64Decode(
    'ChRQdXRMb2dFdmVudHNSZXNwb25zZRIvChFuZXh0c2VxdWVuY2V0b2tlbhjj49NkIAEoCVIRbm'
    'V4dHNlcXVlbmNldG9rZW4SSwoScmVqZWN0ZWRlbnRpdHlpbmZvGJ3m+wEgASgLMhgubG9ncy5S'
    'ZWplY3RlZEVudGl0eUluZm9SEnJlamVjdGVkZW50aXR5aW5mbxJVChVyZWplY3RlZGxvZ2V2ZW'
    '50c2luZm8Ym6uGxwEgASgLMhsubG9ncy5SZWplY3RlZExvZ0V2ZW50c0luZm9SFXJlamVjdGVk'
    'bG9nZXZlbnRzaW5mbw==');

@$core.Deprecated('Use putLogGroupDeletionProtectionRequestDescriptor instead')
const PutLogGroupDeletionProtectionRequest$json = {
  '1': 'PutLogGroupDeletionProtectionRequest',
  '2': [
    {
      '1': 'deletionprotectionenabled',
      '3': 475522738,
      '4': 1,
      '5': 8,
      '10': 'deletionprotectionenabled'
    },
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
  ],
};

/// Descriptor for `PutLogGroupDeletionProtectionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putLogGroupDeletionProtectionRequestDescriptor =
    $convert.base64Decode(
        'CiRQdXRMb2dHcm91cERlbGV0aW9uUHJvdGVjdGlvblJlcXVlc3QSQAoZZGVsZXRpb25wcm90ZW'
        'N0aW9uZW5hYmxlZBiyzd/iASABKAhSGWRlbGV0aW9ucHJvdGVjdGlvbmVuYWJsZWQSMgoSbG9n'
        'Z3JvdXBpZGVudGlmaWVyGNzPpd8BIAEoCVISbG9nZ3JvdXBpZGVudGlmaWVy');

@$core.Deprecated('Use putMetricFilterRequestDescriptor instead')
const PutMetricFilterRequest$json = {
  '1': 'PutMetricFilterRequest',
  '2': [
    {
      '1': 'applyontransformedlogs',
      '3': 99775525,
      '4': 1,
      '5': 8,
      '10': 'applyontransformedlogs'
    },
    {
      '1': 'emitsystemfielddimensions',
      '3': 437431321,
      '4': 3,
      '5': 9,
      '10': 'emitsystemfielddimensions'
    },
    {
      '1': 'fieldselectioncriteria',
      '3': 303984807,
      '4': 1,
      '5': 9,
      '10': 'fieldselectioncriteria'
    },
    {'1': 'filtername', '3': 395125013, '4': 1, '5': 9, '10': 'filtername'},
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'metrictransformations',
      '3': 169353806,
      '4': 3,
      '5': 11,
      '6': '.logs.MetricTransformation',
      '10': 'metrictransformations'
    },
  ],
};

/// Descriptor for `PutMetricFilterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putMetricFilterRequestDescriptor = $convert.base64Decode(
    'ChZQdXRNZXRyaWNGaWx0ZXJSZXF1ZXN0EjkKFmFwcGx5b250cmFuc2Zvcm1lZGxvZ3MYpejJLy'
    'ABKAhSFmFwcGx5b250cmFuc2Zvcm1lZGxvZ3MSQAoZZW1pdHN5c3RlbWZpZWxkZGltZW5zaW9u'
    'cxiZ2MrQASADKAlSGWVtaXRzeXN0ZW1maWVsZGRpbWVuc2lvbnMSOgoWZmllbGRzZWxlY3Rpb2'
    '5jcml0ZXJpYRin4fmQASABKAlSFmZpZWxkc2VsZWN0aW9uY3JpdGVyaWESIgoKZmlsdGVybmFt'
    'ZRiVwrS8ASABKAlSCmZpbHRlcm5hbWUSJwoNZmlsdGVycGF0dGVybhiYh4pFIAEoCVINZmlsdG'
    'VycGF0dGVybhIlCgxsb2dncm91cG5hbWUYzL/ULSABKAlSDGxvZ2dyb3VwbmFtZRJTChVtZXRy'
    'aWN0cmFuc2Zvcm1hdGlvbnMYzsTgUCADKAsyGi5sb2dzLk1ldHJpY1RyYW5zZm9ybWF0aW9uUh'
    'VtZXRyaWN0cmFuc2Zvcm1hdGlvbnM=');

@$core.Deprecated('Use putQueryDefinitionRequestDescriptor instead')
const PutQueryDefinitionRequest$json = {
  '1': 'PutQueryDefinitionRequest',
  '2': [
    {'1': 'clienttoken', '3': 272531820, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'loggroupnames',
      '3': 337702569,
      '4': 3,
      '5': 9,
      '10': 'loggroupnames'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'querydefinitionid',
      '3': 178455620,
      '4': 1,
      '5': 9,
      '10': 'querydefinitionid'
    },
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
  ],
};

/// Descriptor for `PutQueryDefinitionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putQueryDefinitionRequestDescriptor = $convert.base64Decode(
    'ChlQdXRRdWVyeURlZmluaXRpb25SZXF1ZXN0EiQKC2NsaWVudHRva2VuGOyC+oEBIAEoCVILY2'
    'xpZW50dG9rZW4SKAoNbG9nZ3JvdXBuYW1lcxip3YOhASADKAlSDWxvZ2dyb3VwbmFtZXMSFQoE'
    'bmFtZRjn++ZpIAEoCVIEbmFtZRIvChFxdWVyeWRlZmluaXRpb25pZBjEiIxVIAEoCVIRcXVlcn'
    'lkZWZpbml0aW9uaWQSPQoNcXVlcnlsYW5ndWFnZRj2r/ejASABKA4yEy5sb2dzLlF1ZXJ5TGFu'
    'Z3VhZ2VSDXF1ZXJ5bGFuZ3VhZ2USJAoLcXVlcnlzdHJpbmcYh4Gd+AEgASgJUgtxdWVyeXN0cm'
    'luZw==');

@$core.Deprecated('Use putQueryDefinitionResponseDescriptor instead')
const PutQueryDefinitionResponse$json = {
  '1': 'PutQueryDefinitionResponse',
  '2': [
    {
      '1': 'querydefinitionid',
      '3': 178455620,
      '4': 1,
      '5': 9,
      '10': 'querydefinitionid'
    },
  ],
};

/// Descriptor for `PutQueryDefinitionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putQueryDefinitionResponseDescriptor =
    $convert.base64Decode(
        'ChpQdXRRdWVyeURlZmluaXRpb25SZXNwb25zZRIvChFxdWVyeWRlZmluaXRpb25pZBjEiIxVIA'
        'EoCVIRcXVlcnlkZWZpbml0aW9uaWQ=');

@$core.Deprecated('Use putResourcePolicyRequestDescriptor instead')
const PutResourcePolicyRequest$json = {
  '1': 'PutResourcePolicyRequest',
  '2': [
    {
      '1': 'expectedrevisionid',
      '3': 469613378,
      '4': 1,
      '5': 9,
      '10': 'expectedrevisionid'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `PutResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyRequestDescriptor = $convert.base64Decode(
    'ChhQdXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSMgoSZXhwZWN0ZWRyZXZpc2lvbmlkGML29t8BIA'
    'EoCVISZXhwZWN0ZWRyZXZpc2lvbmlkEikKDnBvbGljeWRvY3VtZW50GOvp9lQgASgJUg5wb2xp'
    'Y3lkb2N1bWVudBIhCgpwb2xpY3luYW1lGN3i8jYgASgJUgpwb2xpY3luYW1lEiMKC3Jlc291cm'
    'NlYXJuGM3MqiAgASgJUgtyZXNvdXJjZWFybg==');

@$core.Deprecated('Use putResourcePolicyResponseDescriptor instead')
const PutResourcePolicyResponse$json = {
  '1': 'PutResourcePolicyResponse',
  '2': [
    {
      '1': 'resourcepolicy',
      '3': 305429200,
      '4': 1,
      '5': 11,
      '6': '.logs.ResourcePolicy',
      '10': 'resourcepolicy'
    },
    {'1': 'revisionid', '3': 369170086, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `PutResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyResponseDescriptor = $convert.base64Decode(
    'ChlQdXRSZXNvdXJjZVBvbGljeVJlc3BvbnNlEkAKDnJlc291cmNlcG9saWN5GND10ZEBIAEoCz'
    'IULmxvZ3MuUmVzb3VyY2VQb2xpY3lSDnJlc291cmNlcG9saWN5EiIKCnJldmlzaW9uaWQYpq2E'
    'sAEgASgJUgpyZXZpc2lvbmlk');

@$core.Deprecated('Use putRetentionPolicyRequestDescriptor instead')
const PutRetentionPolicyRequest$json = {
  '1': 'PutRetentionPolicyRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'retentionindays',
      '3': 258337482,
      '4': 1,
      '5': 5,
      '10': 'retentionindays'
    },
  ],
};

/// Descriptor for `PutRetentionPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRetentionPolicyRequestDescriptor = $convert.base64Decode(
    'ChlQdXRSZXRlbnRpb25Qb2xpY3lSZXF1ZXN0EiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG'
    '9nZ3JvdXBuYW1lEisKD3JldGVudGlvbmluZGF5cxjK1Zd7IAEoBVIPcmV0ZW50aW9uaW5kYXlz');

@$core.Deprecated('Use putSubscriptionFilterRequestDescriptor instead')
const PutSubscriptionFilterRequest$json = {
  '1': 'PutSubscriptionFilterRequest',
  '2': [
    {
      '1': 'applyontransformedlogs',
      '3': 99775525,
      '4': 1,
      '5': 8,
      '10': 'applyontransformedlogs'
    },
    {
      '1': 'destinationarn',
      '3': 427601315,
      '4': 1,
      '5': 9,
      '10': 'destinationarn'
    },
    {
      '1': 'distribution',
      '3': 345526572,
      '4': 1,
      '5': 14,
      '6': '.logs.Distribution',
      '10': 'distribution'
    },
    {
      '1': 'emitsystemfields',
      '3': 392618203,
      '4': 3,
      '5': 9,
      '10': 'emitsystemfields'
    },
    {
      '1': 'fieldselectioncriteria',
      '3': 303984807,
      '4': 1,
      '5': 9,
      '10': 'fieldselectioncriteria'
    },
    {'1': 'filtername', '3': 395125013, '4': 1, '5': 9, '10': 'filtername'},
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
  ],
};

/// Descriptor for `PutSubscriptionFilterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putSubscriptionFilterRequestDescriptor = $convert.base64Decode(
    'ChxQdXRTdWJzY3JpcHRpb25GaWx0ZXJSZXF1ZXN0EjkKFmFwcGx5b250cmFuc2Zvcm1lZGxvZ3'
    'MYpejJLyABKAhSFmFwcGx5b250cmFuc2Zvcm1lZGxvZ3MSKgoOZGVzdGluYXRpb25hcm4Yo9vy'
    'ywEgASgJUg5kZXN0aW5hdGlvbmFybhI6CgxkaXN0cmlidXRpb24YrKLhpAEgASgOMhIubG9ncy'
    '5EaXN0cmlidXRpb25SDGRpc3RyaWJ1dGlvbhIuChBlbWl0c3lzdGVtZmllbGRzGNvBm7sBIAMo'
    'CVIQZW1pdHN5c3RlbWZpZWxkcxI6ChZmaWVsZHNlbGVjdGlvbmNyaXRlcmlhGKfh+ZABIAEoCV'
    'IWZmllbGRzZWxlY3Rpb25jcml0ZXJpYRIiCgpmaWx0ZXJuYW1lGJXCtLwBIAEoCVIKZmlsdGVy'
    'bmFtZRInCg1maWx0ZXJwYXR0ZXJuGJiHikUgASgJUg1maWx0ZXJwYXR0ZXJuEiUKDGxvZ2dyb3'
    'VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW1lEhsKB3JvbGVhcm4YoZeJUSABKAlSB3JvbGVh'
    'cm4=');

@$core.Deprecated('Use putTransformerRequestDescriptor instead')
const PutTransformerRequest$json = {
  '1': 'PutTransformerRequest',
  '2': [
    {
      '1': 'loggroupidentifier',
      '3': 468281308,
      '4': 1,
      '5': 9,
      '10': 'loggroupidentifier'
    },
    {
      '1': 'transformerconfig',
      '3': 384836439,
      '4': 3,
      '5': 11,
      '6': '.logs.Processor',
      '10': 'transformerconfig'
    },
  ],
};

/// Descriptor for `PutTransformerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putTransformerRequestDescriptor = $convert.base64Decode(
    'ChVQdXRUcmFuc2Zvcm1lclJlcXVlc3QSMgoSbG9nZ3JvdXBpZGVudGlmaWVyGNzPpd8BIAEoCV'
    'ISbG9nZ3JvdXBpZGVudGlmaWVyEkEKEXRyYW5zZm9ybWVyY29uZmlnGNfGwLcBIAMoCzIPLmxv'
    'Z3MuUHJvY2Vzc29yUhF0cmFuc2Zvcm1lcmNvbmZpZw==');

@$core.Deprecated('Use queryCompileErrorDescriptor instead')
const QueryCompileError$json = {
  '1': 'QueryCompileError',
  '2': [
    {
      '1': 'location',
      '3': 200649127,
      '4': 1,
      '5': 11,
      '6': '.logs.QueryCompileErrorLocation',
      '10': 'location'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueryCompileError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryCompileErrorDescriptor = $convert.base64Decode(
    'ChFRdWVyeUNvbXBpbGVFcnJvchI+Cghsb2NhdGlvbhin09ZfIAEoCzIfLmxvZ3MuUXVlcnlDb2'
    '1waWxlRXJyb3JMb2NhdGlvblIIbG9jYXRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2Fn'
    'ZQ==');

@$core.Deprecated('Use queryCompileErrorLocationDescriptor instead')
const QueryCompileErrorLocation$json = {
  '1': 'QueryCompileErrorLocation',
  '2': [
    {
      '1': 'endcharoffset',
      '3': 14912634,
      '4': 1,
      '5': 5,
      '10': 'endcharoffset'
    },
    {
      '1': 'startcharoffset',
      '3': 391795325,
      '4': 1,
      '5': 5,
      '10': 'startcharoffset'
    },
  ],
};

/// Descriptor for `QueryCompileErrorLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryCompileErrorLocationDescriptor = $convert.base64Decode(
    'ChlRdWVyeUNvbXBpbGVFcnJvckxvY2F0aW9uEicKDWVuZGNoYXJvZmZzZXQY+piOByABKAVSDW'
    'VuZGNoYXJvZmZzZXQSLAoPc3RhcnRjaGFyb2Zmc2V0GP2k6boBIAEoBVIPc3RhcnRjaGFyb2Zm'
    'c2V0');

@$core.Deprecated('Use queryDefinitionDescriptor instead')
const QueryDefinition$json = {
  '1': 'QueryDefinition',
  '2': [
    {'1': 'lastmodified', '3': 410744903, '4': 1, '5': 3, '10': 'lastmodified'},
    {
      '1': 'loggroupnames',
      '3': 337702569,
      '4': 3,
      '5': 9,
      '10': 'loggroupnames'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'querydefinitionid',
      '3': 178455620,
      '4': 1,
      '5': 9,
      '10': 'querydefinitionid'
    },
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
  ],
};

/// Descriptor for `QueryDefinition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryDefinitionDescriptor = $convert.base64Decode(
    'Cg9RdWVyeURlZmluaXRpb24SJgoMbGFzdG1vZGlmaWVkGMfw7cMBIAEoA1IMbGFzdG1vZGlmaW'
    'VkEigKDWxvZ2dyb3VwbmFtZXMYqd2DoQEgAygJUg1sb2dncm91cG5hbWVzEhUKBG5hbWUY5/vm'
    'aSABKAlSBG5hbWUSLwoRcXVlcnlkZWZpbml0aW9uaWQYxIiMVSABKAlSEXF1ZXJ5ZGVmaW5pdG'
    'lvbmlkEj0KDXF1ZXJ5bGFuZ3VhZ2UY9q/3owEgASgOMhMubG9ncy5RdWVyeUxhbmd1YWdlUg1x'
    'dWVyeWxhbmd1YWdlEiQKC3F1ZXJ5c3RyaW5nGIeBnfgBIAEoCVILcXVlcnlzdHJpbmc=');

@$core.Deprecated('Use queryInfoDescriptor instead')
const QueryInfo$json = {
  '1': 'QueryInfo',
  '2': [
    {'1': 'createtime', '3': 297700189, '4': 1, '5': 3, '10': 'createtime'},
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'queryid', '3': 336975759, '4': 1, '5': 9, '10': 'queryid'},
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `QueryInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryInfoDescriptor = $convert.base64Decode(
    'CglRdWVyeUluZm8SIgoKY3JlYXRldGltZRjdlvqNASABKANSCmNyZWF0ZXRpbWUSJQoMbG9nZ3'
    'JvdXBuYW1lGMy/1C0gASgJUgxsb2dncm91cG5hbWUSHAoHcXVlcnlpZBiPr9egASABKAlSB3F1'
    'ZXJ5aWQSPQoNcXVlcnlsYW5ndWFnZRj2r/ejASABKA4yEy5sb2dzLlF1ZXJ5TGFuZ3VhZ2VSDX'
    'F1ZXJ5bGFuZ3VhZ2USJAoLcXVlcnlzdHJpbmcYh4Gd+AEgASgJUgtxdWVyeXN0cmluZxItCgZz'
    'dGF0dXMY8O+t0gEgASgOMhEubG9ncy5RdWVyeVN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use queryStatisticsDescriptor instead')
const QueryStatistics$json = {
  '1': 'QueryStatistics',
  '2': [
    {'1': 'bytesscanned', '3': 533301977, '4': 1, '5': 1, '10': 'bytesscanned'},
    {
      '1': 'estimatedbytesskipped',
      '3': 243901805,
      '4': 1,
      '5': 1,
      '10': 'estimatedbytesskipped'
    },
    {
      '1': 'estimatedrecordsskipped',
      '3': 495275708,
      '4': 1,
      '5': 1,
      '10': 'estimatedrecordsskipped'
    },
    {
      '1': 'loggroupsscanned',
      '3': 188426568,
      '4': 1,
      '5': 1,
      '10': 'loggroupsscanned'
    },
    {
      '1': 'recordsmatched',
      '3': 338750832,
      '4': 1,
      '5': 1,
      '10': 'recordsmatched'
    },
    {
      '1': 'recordsscanned',
      '3': 431708156,
      '4': 1,
      '5': 1,
      '10': 'recordsscanned'
    },
  ],
};

/// Descriptor for `QueryStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStatisticsDescriptor = $convert.base64Decode(
    'Cg9RdWVyeVN0YXRpc3RpY3MSJgoMYnl0ZXNzY2FubmVkGNmVpv4BIAEoAVIMYnl0ZXNzY2Fubm'
    'VkEjcKFWVzdGltYXRlZGJ5dGVzc2tpcHBlZBjtyqZ0IAEoAVIVZXN0aW1hdGVkYnl0ZXNza2lw'
    'cGVkEjwKF2VzdGltYXRlZHJlY29yZHNza2lwcGVkGLydlewBIAEoAVIXZXN0aW1hdGVkcmVjb3'
    'Jkc3NraXBwZWQSLQoQbG9nZ3JvdXBzc2Nhbm5lZBjI0uxZIAEoAVIQbG9nZ3JvdXBzc2Nhbm5l'
    'ZBIqCg5yZWNvcmRzbWF0Y2hlZBjw2sOhASABKAFSDnJlY29yZHNtYXRjaGVkEioKDnJlY29yZH'
    'NzY2FubmVkGPyv7c0BIAEoAVIOcmVjb3Jkc3NjYW5uZWQ=');

@$core.Deprecated('Use recordFieldDescriptor instead')
const RecordField$json = {
  '1': 'RecordField',
  '2': [
    {'1': 'mandatory', '3': 264682169, '4': 1, '5': 8, '10': 'mandatory'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `RecordField`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordFieldDescriptor = $convert.base64Decode(
    'CgtSZWNvcmRGaWVsZBIfCgltYW5kYXRvcnkYufWafiABKAhSCW1hbmRhdG9yeRIVCgRuYW1lGO'
    'f75mkgASgJUgRuYW1l');

@$core.Deprecated('Use rejectedEntityInfoDescriptor instead')
const RejectedEntityInfo$json = {
  '1': 'RejectedEntityInfo',
  '2': [
    {
      '1': 'errortype',
      '3': 183231834,
      '4': 1,
      '5': 14,
      '6': '.logs.EntityRejectionErrorType',
      '10': 'errortype'
    },
  ],
};

/// Descriptor for `RejectedEntityInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rejectedEntityInfoDescriptor = $convert.base64Decode(
    'ChJSZWplY3RlZEVudGl0eUluZm8SPwoJZXJyb3J0eXBlGNrKr1cgASgOMh4ubG9ncy5FbnRpdH'
    'lSZWplY3Rpb25FcnJvclR5cGVSCWVycm9ydHlwZQ==');

@$core.Deprecated('Use rejectedLogEventsInfoDescriptor instead')
const RejectedLogEventsInfo$json = {
  '1': 'RejectedLogEventsInfo',
  '2': [
    {
      '1': 'expiredlogeventendindex',
      '3': 299438860,
      '4': 1,
      '5': 5,
      '10': 'expiredlogeventendindex'
    },
    {
      '1': 'toonewlogeventstartindex',
      '3': 397251384,
      '4': 1,
      '5': 5,
      '10': 'toonewlogeventstartindex'
    },
    {
      '1': 'toooldlogeventendindex',
      '3': 96023564,
      '4': 1,
      '5': 5,
      '10': 'toooldlogeventendindex'
    },
  ],
};

/// Descriptor for `RejectedLogEventsInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rejectedLogEventsInfoDescriptor = $convert.base64Decode(
    'ChVSZWplY3RlZExvZ0V2ZW50c0luZm8SPAoXZXhwaXJlZGxvZ2V2ZW50ZW5kaW5kZXgYjKbkjg'
    'EgASgFUhdleHBpcmVkbG9nZXZlbnRlbmRpbmRleBI+Chh0b29uZXdsb2dldmVudHN0YXJ0aW5k'
    'ZXgYuKa2vQEgASgFUhh0b29uZXdsb2dldmVudHN0YXJ0aW5kZXgSOQoWdG9vb2xkbG9nZXZlbn'
    'RlbmRpbmRleBiM6OQtIAEoBVIWdG9vb2xkbG9nZXZlbnRlbmRpbmRleA==');

@$core.Deprecated('Use renameKeyEntryDescriptor instead')
const RenameKeyEntry$json = {
  '1': 'RenameKeyEntry',
  '2': [
    {'1': 'key', '3': 135645293, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'overwriteifexists',
      '3': 230880030,
      '4': 1,
      '5': 8,
      '10': 'overwriteifexists'
    },
    {'1': 'renameto', '3': 370293599, '4': 1, '5': 9, '10': 'renameto'},
  ],
};

/// Descriptor for `RenameKeyEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renameKeyEntryDescriptor = $convert.base64Decode(
    'Cg5SZW5hbWVLZXlFbnRyeRITCgNrZXkY7ZDXQCABKAlSA2tleRIvChFvdmVyd3JpdGVpZmV4aX'
    'N0cxie5otuIAEoCFIRb3ZlcndyaXRlaWZleGlzdHMSHgoIcmVuYW1ldG8Y3/bIsAEgASgJUghy'
    'ZW5hbWV0bw==');

@$core.Deprecated('Use renameKeysDescriptor instead')
const RenameKeys$json = {
  '1': 'RenameKeys',
  '2': [
    {
      '1': 'entries',
      '3': 257458932,
      '4': 3,
      '5': 11,
      '6': '.logs.RenameKeyEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `RenameKeys`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renameKeysDescriptor = $convert.base64Decode(
    'CgpSZW5hbWVLZXlzEjEKB2VudHJpZXMY9IXieiADKAsyFC5sb2dzLlJlbmFtZUtleUVudHJ5Ug'
    'dlbnRyaWVz');

@$core.Deprecated('Use resourceAlreadyExistsExceptionDescriptor instead')
const ResourceAlreadyExistsException$json = {
  '1': 'ResourceAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'Ch5SZXNvdXJjZUFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use resourceConfigDescriptor instead')
const ResourceConfig$json = {
  '1': 'ResourceConfig',
  '2': [
    {
      '1': 'opensearchresourceconfig',
      '3': 365109608,
      '4': 1,
      '5': 11,
      '6': '.logs.OpenSearchResourceConfig',
      '10': 'opensearchresourceconfig'
    },
  ],
};

/// Descriptor for `ResourceConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceConfigDescriptor = $convert.base64Decode(
    'Cg5SZXNvdXJjZUNvbmZpZxJeChhvcGVuc2VhcmNocmVzb3VyY2Vjb25maWcY6MKMrgEgASgLMh'
    '4ubG9ncy5PcGVuU2VhcmNoUmVzb3VyY2VDb25maWdSGG9wZW5zZWFyY2hyZXNvdXJjZWNvbmZp'
    'Zw==');

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

@$core.Deprecated('Use resourcePolicyDescriptor instead')
const ResourcePolicy$json = {
  '1': 'ResourcePolicy',
  '2': [
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'policydocument',
      '3': 178107627,
      '4': 1,
      '5': 9,
      '10': 'policydocument'
    },
    {'1': 'policyname', '3': 115126621, '4': 1, '5': 9, '10': 'policyname'},
    {
      '1': 'policyscope',
      '3': 288841470,
      '4': 1,
      '5': 14,
      '6': '.logs.PolicyScope',
      '10': 'policyscope'
    },
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'revisionid', '3': 369170086, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `ResourcePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcePolicyDescriptor = $convert.base64Decode(
    'Cg5SZXNvdXJjZVBvbGljeRIsCg9sYXN0dXBkYXRlZHRpbWUY9reVuQEgASgDUg9sYXN0dXBkYX'
    'RlZHRpbWUSKQoOcG9saWN5ZG9jdW1lbnQY6+n2VCABKAlSDnBvbGljeWRvY3VtZW50EiEKCnBv'
    'bGljeW5hbWUY3eLyNiABKAlSCnBvbGljeW5hbWUSNwoLcG9saWN5c2NvcGUY/r3diQEgASgOMh'
    'EubG9ncy5Qb2xpY3lTY29wZVILcG9saWN5c2NvcGUSIwoLcmVzb3VyY2Vhcm4YzcyqICABKAlS'
    'C3Jlc291cmNlYXJuEiIKCnJldmlzaW9uaWQYpq2EsAEgASgJUgpyZXZpc2lvbmlk');

@$core.Deprecated('Use resultFieldDescriptor instead')
const ResultField$json = {
  '1': 'ResultField',
  '2': [
    {'1': 'field', '3': 125985384, '4': 1, '5': 9, '10': 'field'},
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ResultField`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultFieldDescriptor = $convert.base64Decode(
    'CgtSZXN1bHRGaWVsZBIXCgVmaWVsZBjoxIk8IAEoCVIFZmllbGQSFwoFdmFsdWUYy6f7EiABKA'
    'lSBXZhbHVl');

@$core.Deprecated('Use s3ConfigurationDescriptor instead')
const S3Configuration$json = {
  '1': 'S3Configuration',
  '2': [
    {
      '1': 'destinationidentifier',
      '3': 399670053,
      '4': 1,
      '5': 9,
      '10': 'destinationidentifier'
    },
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
  ],
};

/// Descriptor for `S3Configuration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3ConfigurationDescriptor = $convert.base64Decode(
    'Cg9TM0NvbmZpZ3VyYXRpb24SOAoVZGVzdGluYXRpb25pZGVudGlmaWVyGKX2yb4BIAEoCVIVZG'
    'VzdGluYXRpb25pZGVudGlmaWVyEhsKB3JvbGVhcm4YoZeJUSABKAlSB3JvbGVhcm4=');

@$core.Deprecated('Use s3DeliveryConfigurationDescriptor instead')
const S3DeliveryConfiguration$json = {
  '1': 'S3DeliveryConfiguration',
  '2': [
    {
      '1': 'enablehivecompatiblepath',
      '3': 83012438,
      '4': 1,
      '5': 8,
      '10': 'enablehivecompatiblepath'
    },
    {'1': 'suffixpath', '3': 249314700, '4': 1, '5': 9, '10': 'suffixpath'},
  ],
};

/// Descriptor for `S3DeliveryConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3DeliveryConfigurationDescriptor = $convert.base64Decode(
    'ChdTM0RlbGl2ZXJ5Q29uZmlndXJhdGlvbhI9ChhlbmFibGVoaXZlY29tcGF0aWJsZXBhdGgY1t'
    'bKJyABKAhSGGVuYWJsZWhpdmVjb21wYXRpYmxlcGF0aBIhCgpzdWZmaXhwYXRoGIz78HYgASgJ'
    'UgpzdWZmaXhwYXRo');

@$core.Deprecated('Use s3TableIntegrationSourceDescriptor instead')
const S3TableIntegrationSource$json = {
  '1': 'S3TableIntegrationSource',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 462845754,
      '4': 1,
      '5': 3,
      '10': 'createdtimestamp'
    },
    {
      '1': 'datasource',
      '3': 345762713,
      '4': 1,
      '5': 11,
      '6': '.logs.DataSource',
      '10': 'datasource'
    },
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.logs.S3TableIntegrationSourceStatus',
      '10': 'status'
    },
    {'1': 'statusreason', '3': 352592412, '4': 1, '5': 9, '10': 'statusreason'},
  ],
};

/// Descriptor for `S3TableIntegrationSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3TableIntegrationSourceDescriptor = $convert.base64Decode(
    'ChhTM1RhYmxlSW50ZWdyYXRpb25Tb3VyY2USLgoQY3JlYXRlZHRpbWVzdGFtcBi67tncASABKA'
    'NSEGNyZWF0ZWR0aW1lc3RhbXASNAoKZGF0YXNvdXJjZRiZ1++kASABKAsyEC5sb2dzLkRhdGFT'
    'b3VyY2VSCmRhdGFzb3VyY2USIQoKaWRlbnRpZmllchi/0JZFIAEoCVIKaWRlbnRpZmllchJACg'
    'ZzdGF0dXMY8O+t0gEgASgOMiQubG9ncy5TM1RhYmxlSW50ZWdyYXRpb25Tb3VyY2VTdGF0dXNS'
    'BnN0YXR1cxImCgxzdGF0dXNyZWFzb24YnMSQqAEgASgJUgxzdGF0dXNyZWFzb24=');

@$core.Deprecated('Use scheduledQueryDestinationDescriptor instead')
const ScheduledQueryDestination$json = {
  '1': 'ScheduledQueryDestination',
  '2': [
    {
      '1': 'destinationidentifier',
      '3': 399670053,
      '4': 1,
      '5': 9,
      '10': 'destinationidentifier'
    },
    {
      '1': 'destinationtype',
      '3': 488995304,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryDestinationType',
      '10': 'destinationtype'
    },
    {'1': 'errormessage', '3': 136873289, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'processedidentifier',
      '3': 194381665,
      '4': 1,
      '5': 9,
      '10': 'processedidentifier'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.logs.ActionStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ScheduledQueryDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledQueryDestinationDescriptor = $convert.base64Decode(
    'ChlTY2hlZHVsZWRRdWVyeURlc3RpbmF0aW9uEjgKFWRlc3RpbmF0aW9uaWRlbnRpZmllchil9s'
    'm+ASABKAlSFWRlc3RpbmF0aW9uaWRlbnRpZmllchJRCg9kZXN0aW5hdGlvbnR5cGUY6POV6QEg'
    'ASgOMiMubG9ncy5TY2hlZHVsZWRRdWVyeURlc3RpbmF0aW9uVHlwZVIPZGVzdGluYXRpb250eX'
    'BlEiUKDGVycm9ybWVzc2FnZRjJiqJBIAEoCVIMZXJyb3JtZXNzYWdlEjMKE3Byb2Nlc3NlZGlk'
    'ZW50aWZpZXIY4Y7YXCABKAlSE3Byb2Nlc3NlZGlkZW50aWZpZXISLgoGc3RhdHVzGPDvrdIBIA'
    'EoDjISLmxvZ3MuQWN0aW9uU3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use scheduledQuerySummaryDescriptor instead')
const ScheduledQuerySummary$json = {
  '1': 'ScheduledQuerySummary',
  '2': [
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'destinationconfiguration',
      '3': 145518016,
      '4': 1,
      '5': 11,
      '6': '.logs.DestinationConfiguration',
      '10': 'destinationconfiguration'
    },
    {
      '1': 'lastexecutionstatus',
      '3': 7800596,
      '4': 1,
      '5': 14,
      '6': '.logs.ExecutionStatus',
      '10': 'lastexecutionstatus'
    },
    {
      '1': 'lasttriggeredtime',
      '3': 397057656,
      '4': 1,
      '5': 3,
      '10': 'lasttriggeredtime'
    },
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'scheduleexpression',
      '3': 287794975,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduledqueryarn',
      '3': 240292916,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryState',
      '10': 'state'
    },
    {'1': 'timezone', '3': 190615331, '4': 1, '5': 9, '10': 'timezone'},
  ],
};

/// Descriptor for `ScheduledQuerySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledQuerySummaryDescriptor = $convert.base64Decode(
    'ChVTY2hlZHVsZWRRdWVyeVN1bW1hcnkSJQoMY3JlYXRpb250aW1lGIbFxBkgASgDUgxjcmVhdG'
    'lvbnRpbWUSXQoYZGVzdGluYXRpb25jb25maWd1cmF0aW9uGMDbsUUgASgLMh4ubG9ncy5EZXN0'
    'aW5hdGlvbkNvbmZpZ3VyYXRpb25SGGRlc3RpbmF0aW9uY29uZmlndXJhdGlvbhJKChNsYXN0ZX'
    'hlY3V0aW9uc3RhdHVzGJSO3AMgASgOMhUubG9ncy5FeGVjdXRpb25TdGF0dXNSE2xhc3RleGVj'
    'dXRpb25zdGF0dXMSMAoRbGFzdHRyaWdnZXJlZHRpbWUY+LyqvQEgASgDUhFsYXN0dHJpZ2dlcm'
    'VkdGltZRIsCg9sYXN0dXBkYXRlZHRpbWUY9reVuQEgASgDUg9sYXN0dXBkYXRlZHRpbWUSFQoE'
    'bmFtZRjn++ZpIAEoCVIEbmFtZRIyChJzY2hlZHVsZWV4cHJlc3Npb24Yn86diQEgASgJUhJzY2'
    'hlZHVsZWV4cHJlc3Npb24SLwoRc2NoZWR1bGVkcXVlcnlhcm4YtKjKciABKAlSEXNjaGVkdWxl'
    'ZHF1ZXJ5YXJuEjMKBXN0YXRlGPflxMEBIAEoDjIZLmxvZ3MuU2NoZWR1bGVkUXVlcnlTdGF0ZV'
    'IFc3RhdGUSHQoIdGltZXpvbmUYo57yWiABKAlSCHRpbWV6b25l');

@$core.Deprecated('Use searchedLogStreamDescriptor instead')
const SearchedLogStream$json = {
  '1': 'SearchedLogStream',
  '2': [
    {
      '1': 'logstreamname',
      '3': 438025123,
      '4': 1,
      '5': 9,
      '10': 'logstreamname'
    },
    {
      '1': 'searchedcompletely',
      '3': 449344361,
      '4': 1,
      '5': 8,
      '10': 'searchedcompletely'
    },
  ],
};

/// Descriptor for `SearchedLogStream`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchedLogStreamDescriptor = $convert.base64Decode(
    'ChFTZWFyY2hlZExvZ1N0cmVhbRIoCg1sb2dzdHJlYW1uYW1lGKP37tABIAEoCVINbG9nc3RyZW'
    'FtbmFtZRIyChJzZWFyY2hlZGNvbXBsZXRlbHkY6eah1gEgASgIUhJzZWFyY2hlZGNvbXBsZXRl'
    'bHk=');

@$core.Deprecated('Use serviceQuotaExceededExceptionDescriptor instead')
const ServiceQuotaExceededException$json = {
  '1': 'ServiceQuotaExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceQuotaExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceQuotaExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1TZXJ2aWNlUXVvdGFFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use serviceUnavailableExceptionDescriptor instead')
const ServiceUnavailableException$json = {
  '1': 'ServiceUnavailableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceUnavailableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceUnavailableExceptionDescriptor =
    $convert.base64Decode(
        'ChtTZXJ2aWNlVW5hdmFpbGFibGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use sessionStreamingExceptionDescriptor instead')
const SessionStreamingException$json = {
  '1': 'SessionStreamingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SessionStreamingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionStreamingExceptionDescriptor =
    $convert.base64Decode(
        'ChlTZXNzaW9uU3RyZWFtaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use sessionTimeoutExceptionDescriptor instead')
const SessionTimeoutException$json = {
  '1': 'SessionTimeoutException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SessionTimeoutException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionTimeoutExceptionDescriptor =
    $convert.base64Decode(
        'ChdTZXNzaW9uVGltZW91dEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use splitStringDescriptor instead')
const SplitString$json = {
  '1': 'SplitString',
  '2': [
    {
      '1': 'entries',
      '3': 257458932,
      '4': 3,
      '5': 11,
      '6': '.logs.SplitStringEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `SplitString`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List splitStringDescriptor = $convert.base64Decode(
    'CgtTcGxpdFN0cmluZxIzCgdlbnRyaWVzGPSF4nogAygLMhYubG9ncy5TcGxpdFN0cmluZ0VudH'
    'J5UgdlbnRyaWVz');

@$core.Deprecated('Use splitStringEntryDescriptor instead')
const SplitStringEntry$json = {
  '1': 'SplitStringEntry',
  '2': [
    {'1': 'delimiter', '3': 202749435, '4': 1, '5': 9, '10': 'delimiter'},
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `SplitStringEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List splitStringEntryDescriptor = $convert.base64Decode(
    'ChBTcGxpdFN0cmluZ0VudHJ5Eh8KCWRlbGltaXRlchj769ZgIAEoCVIJZGVsaW1pdGVyEhoKBn'
    'NvdXJjZRjZ07zeASABKAlSBnNvdXJjZQ==');

@$core.Deprecated('Use startLiveTailRequestDescriptor instead')
const StartLiveTailRequest$json = {
  '1': 'StartLiveTailRequest',
  '2': [
    {
      '1': 'logeventfilterpattern',
      '3': 81051802,
      '4': 1,
      '5': 9,
      '10': 'logeventfilterpattern'
    },
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {
      '1': 'logstreamnameprefixes',
      '3': 109414303,
      '4': 3,
      '5': 9,
      '10': 'logstreamnameprefixes'
    },
    {
      '1': 'logstreamnames',
      '3': 178825732,
      '4': 3,
      '5': 9,
      '10': 'logstreamnames'
    },
  ],
};

/// Descriptor for `StartLiveTailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startLiveTailRequestDescriptor = $convert.base64Decode(
    'ChRTdGFydExpdmVUYWlsUmVxdWVzdBI3ChVsb2dldmVudGZpbHRlcnBhdHRlcm4YmoHTJiABKA'
    'lSFWxvZ2V2ZW50ZmlsdGVycGF0dGVybhI0ChNsb2dncm91cGlkZW50aWZpZXJzGPnauMMBIAMo'
    'CVITbG9nZ3JvdXBpZGVudGlmaWVycxI3ChVsb2dzdHJlYW1uYW1lcHJlZml4ZXMYn4+WNCADKA'
    'lSFWxvZ3N0cmVhbW5hbWVwcmVmaXhlcxIpCg5sb2dzdHJlYW1uYW1lcxiE1KJVIAMoCVIObG9n'
    'c3RyZWFtbmFtZXM=');

@$core.Deprecated('Use startLiveTailResponseDescriptor instead')
const StartLiveTailResponse$json = {
  '1': 'StartLiveTailResponse',
  '2': [
    {
      '1': 'responsestream',
      '3': 238600735,
      '4': 1,
      '5': 11,
      '6': '.logs.StartLiveTailResponseStream',
      '10': 'responsestream'
    },
  ],
};

/// Descriptor for `StartLiveTailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startLiveTailResponseDescriptor = $convert.base64Decode(
    'ChVTdGFydExpdmVUYWlsUmVzcG9uc2USTAoOcmVzcG9uc2VzdHJlYW0Yn4TjcSABKAsyIS5sb2'
    'dzLlN0YXJ0TGl2ZVRhaWxSZXNwb25zZVN0cmVhbVIOcmVzcG9uc2VzdHJlYW0=');

@$core.Deprecated('Use startLiveTailResponseStreamDescriptor instead')
const StartLiveTailResponseStream$json = {
  '1': 'StartLiveTailResponseStream',
  '2': [
    {
      '1': 'sessionstreamingexception',
      '3': 433896323,
      '4': 1,
      '5': 11,
      '6': '.logs.SessionStreamingException',
      '10': 'sessionstreamingexception'
    },
    {
      '1': 'sessiontimeoutexception',
      '3': 457664772,
      '4': 1,
      '5': 11,
      '6': '.logs.SessionTimeoutException',
      '10': 'sessiontimeoutexception'
    },
    {
      '1': 'sessionstart',
      '3': 464332638,
      '4': 1,
      '5': 11,
      '6': '.logs.LiveTailSessionStart',
      '10': 'sessionstart'
    },
    {
      '1': 'sessionupdate',
      '3': 498561863,
      '4': 1,
      '5': 11,
      '6': '.logs.LiveTailSessionUpdate',
      '10': 'sessionupdate'
    },
  ],
};

/// Descriptor for `StartLiveTailResponseStream`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startLiveTailResponseStreamDescriptor = $convert.base64Decode(
    'ChtTdGFydExpdmVUYWlsUmVzcG9uc2VTdHJlYW0SYQoZc2Vzc2lvbnN0cmVhbWluZ2V4Y2VwdG'
    'lvbhiD9/LOASABKAsyHy5sb2dzLlNlc3Npb25TdHJlYW1pbmdFeGNlcHRpb25SGXNlc3Npb25z'
    'dHJlYW1pbmdleGNlcHRpb24SWwoXc2Vzc2lvbnRpbWVvdXRleGNlcHRpb24YhNKd2gEgASgLMh'
    '0ubG9ncy5TZXNzaW9uVGltZW91dEV4Y2VwdGlvblIXc2Vzc2lvbnRpbWVvdXRleGNlcHRpb24S'
    'QgoMc2Vzc2lvbnN0YXJ0GN7OtN0BIAEoCzIaLmxvZ3MuTGl2ZVRhaWxTZXNzaW9uU3RhcnRSDH'
    'Nlc3Npb25zdGFydBJFCg1zZXNzaW9udXBkYXRlGMfm3e0BIAEoCzIbLmxvZ3MuTGl2ZVRhaWxT'
    'ZXNzaW9uVXBkYXRlUg1zZXNzaW9udXBkYXRl');

@$core.Deprecated('Use startQueryRequestDescriptor instead')
const StartQueryRequest$json = {
  '1': 'StartQueryRequest',
  '2': [
    {'1': 'endtime', '3': 329679852, '4': 1, '5': 3, '10': 'endtime'},
    {'1': 'limit', '3': 316332341, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'loggroupnames',
      '3': 337702569,
      '4': 3,
      '5': 9,
      '10': 'loggroupnames'
    },
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
    {'1': 'starttime', '3': 178154767, '4': 1, '5': 3, '10': 'starttime'},
  ],
};

/// Descriptor for `StartQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startQueryRequestDescriptor = $convert.base64Decode(
    'ChFTdGFydFF1ZXJ5UmVxdWVzdBIcCgdlbmR0aW1lGOyHmp0BIAEoA1IHZW5kdGltZRIYCgVsaW'
    '1pdBi1suuWASABKAVSBWxpbWl0EjQKE2xvZ2dyb3VwaWRlbnRpZmllcnMY+dq4wwEgAygJUhNs'
    'b2dncm91cGlkZW50aWZpZXJzEiUKDGxvZ2dyb3VwbmFtZRjMv9QtIAEoCVIMbG9nZ3JvdXBuYW'
    '1lEigKDWxvZ2dyb3VwbmFtZXMYqd2DoQEgAygJUg1sb2dncm91cG5hbWVzEj0KDXF1ZXJ5bGFu'
    'Z3VhZ2UY9q/3owEgASgOMhMubG9ncy5RdWVyeUxhbmd1YWdlUg1xdWVyeWxhbmd1YWdlEiQKC3'
    'F1ZXJ5c3RyaW5nGIeBnfgBIAEoCVILcXVlcnlzdHJpbmcSHwoJc3RhcnR0aW1lGI/a+VQgASgD'
    'UglzdGFydHRpbWU=');

@$core.Deprecated('Use startQueryResponseDescriptor instead')
const StartQueryResponse$json = {
  '1': 'StartQueryResponse',
  '2': [
    {'1': 'queryid', '3': 336975759, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `StartQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startQueryResponseDescriptor =
    $convert.base64Decode(
        'ChJTdGFydFF1ZXJ5UmVzcG9uc2USHAoHcXVlcnlpZBiPr9egASABKAlSB3F1ZXJ5aWQ=');

@$core.Deprecated('Use stopQueryRequestDescriptor instead')
const StopQueryRequest$json = {
  '1': 'StopQueryRequest',
  '2': [
    {'1': 'queryid', '3': 336975759, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `StopQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopQueryRequestDescriptor = $convert.base64Decode(
    'ChBTdG9wUXVlcnlSZXF1ZXN0EhwKB3F1ZXJ5aWQYj6/XoAEgASgJUgdxdWVyeWlk');

@$core.Deprecated('Use stopQueryResponseDescriptor instead')
const StopQueryResponse$json = {
  '1': 'StopQueryResponse',
  '2': [
    {'1': 'success', '3': 442482449, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `StopQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopQueryResponseDescriptor = $convert.base64Decode(
    'ChFTdG9wUXVlcnlSZXNwb25zZRIcCgdzdWNjZXNzGJH+/tIBIAEoCFIHc3VjY2Vzcw==');

@$core.Deprecated('Use subscriptionFilterDescriptor instead')
const SubscriptionFilter$json = {
  '1': 'SubscriptionFilter',
  '2': [
    {
      '1': 'applyontransformedlogs',
      '3': 99775525,
      '4': 1,
      '5': 8,
      '10': 'applyontransformedlogs'
    },
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {
      '1': 'destinationarn',
      '3': 427601315,
      '4': 1,
      '5': 9,
      '10': 'destinationarn'
    },
    {
      '1': 'distribution',
      '3': 345526572,
      '4': 1,
      '5': 14,
      '6': '.logs.Distribution',
      '10': 'distribution'
    },
    {
      '1': 'emitsystemfields',
      '3': 392618203,
      '4': 3,
      '5': 9,
      '10': 'emitsystemfields'
    },
    {
      '1': 'fieldselectioncriteria',
      '3': 303984807,
      '4': 1,
      '5': 9,
      '10': 'fieldselectioncriteria'
    },
    {'1': 'filtername', '3': 395125013, '4': 1, '5': 9, '10': 'filtername'},
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
  ],
};

/// Descriptor for `SubscriptionFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscriptionFilterDescriptor = $convert.base64Decode(
    'ChJTdWJzY3JpcHRpb25GaWx0ZXISOQoWYXBwbHlvbnRyYW5zZm9ybWVkbG9ncxil6MkvIAEoCF'
    'IWYXBwbHlvbnRyYW5zZm9ybWVkbG9ncxIlCgxjcmVhdGlvbnRpbWUYhsXEGSABKANSDGNyZWF0'
    'aW9udGltZRIqCg5kZXN0aW5hdGlvbmFybhij2/LLASABKAlSDmRlc3RpbmF0aW9uYXJuEjoKDG'
    'Rpc3RyaWJ1dGlvbhisouGkASABKA4yEi5sb2dzLkRpc3RyaWJ1dGlvblIMZGlzdHJpYnV0aW9u'
    'Ei4KEGVtaXRzeXN0ZW1maWVsZHMY28GbuwEgAygJUhBlbWl0c3lzdGVtZmllbGRzEjoKFmZpZW'
    'xkc2VsZWN0aW9uY3JpdGVyaWEYp+H5kAEgASgJUhZmaWVsZHNlbGVjdGlvbmNyaXRlcmlhEiIK'
    'CmZpbHRlcm5hbWUYlcK0vAEgASgJUgpmaWx0ZXJuYW1lEicKDWZpbHRlcnBhdHRlcm4YmIeKRS'
    'ABKAlSDWZpbHRlcnBhdHRlcm4SJQoMbG9nZ3JvdXBuYW1lGMy/1C0gASgJUgxsb2dncm91cG5h'
    'bWUSGwoHcm9sZWFybhihl4lRIAEoCVIHcm9sZWFybg==');

@$core.Deprecated('Use substituteStringDescriptor instead')
const SubstituteString$json = {
  '1': 'SubstituteString',
  '2': [
    {
      '1': 'entries',
      '3': 257458932,
      '4': 3,
      '5': 11,
      '6': '.logs.SubstituteStringEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `SubstituteString`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List substituteStringDescriptor = $convert.base64Decode(
    'ChBTdWJzdGl0dXRlU3RyaW5nEjgKB2VudHJpZXMY9IXieiADKAsyGy5sb2dzLlN1YnN0aXR1dG'
    'VTdHJpbmdFbnRyeVIHZW50cmllcw==');

@$core.Deprecated('Use substituteStringEntryDescriptor instead')
const SubstituteStringEntry$json = {
  '1': 'SubstituteStringEntry',
  '2': [
    {'1': 'from', '3': 365789302, '4': 1, '5': 9, '10': 'from'},
    {'1': 'source', '3': 466561497, '4': 1, '5': 9, '10': 'source'},
    {'1': 'to', '3': 38094885, '4': 1, '5': 9, '10': 'to'},
  ],
};

/// Descriptor for `SubstituteStringEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List substituteStringEntryDescriptor = $convert.base64Decode(
    'ChVTdWJzdGl0dXRlU3RyaW5nRW50cnkSFgoEZnJvbRj2gLauASABKAlSBGZyb20SGgoGc291cm'
    'NlGNnTvN4BIAEoCVIGc291cmNlEhEKAnRvGKWQlRIgASgJUgJ0bw==');

@$core.Deprecated('Use suppressionPeriodDescriptor instead')
const SuppressionPeriod$json = {
  '1': 'SuppressionPeriod',
  '2': [
    {
      '1': 'suppressionunit',
      '3': 423868455,
      '4': 1,
      '5': 14,
      '6': '.logs.SuppressionUnit',
      '10': 'suppressionunit'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 5, '10': 'value'},
  ],
};

/// Descriptor for `SuppressionPeriod`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionPeriodDescriptor = $convert.base64Decode(
    'ChFTdXBwcmVzc2lvblBlcmlvZBJDCg9zdXBwcmVzc2lvbnVuaXQYp/COygEgASgOMhUubG9ncy'
    '5TdXBwcmVzc2lvblVuaXRSD3N1cHByZXNzaW9udW5pdBIXCgV2YWx1ZRjLp/sSIAEoBVIFdmFs'
    'dWU=');

@$core.Deprecated('Use tagLogGroupRequestDescriptor instead')
const TagLogGroupRequest$json = {
  '1': 'TagLogGroupRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.TagLogGroupRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [TagLogGroupRequest_TagsEntry$json],
};

@$core.Deprecated('Use tagLogGroupRequestDescriptor instead')
const TagLogGroupRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TagLogGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagLogGroupRequestDescriptor = $convert.base64Decode(
    'ChJUYWdMb2dHcm91cFJlcXVlc3QSJQoMbG9nZ3JvdXBuYW1lGMy/1C0gASgJUgxsb2dncm91cG'
    '5hbWUSOgoEdGFncxih19ugASADKAsyIi5sb2dzLlRhZ0xvZ0dyb3VwUmVxdWVzdC5UYWdzRW50'
    'cnlSBHRhZ3MaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUg'
    'V2YWx1ZToCOAE=');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.logs.TagResourceRequest.TagsEntry',
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
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSIwoLcmVzb3VyY2Vhcm4YzcyqICABKAlSC3Jlc291cmNlYX'
    'JuEjoKBHRhZ3MYodfboAEgAygLMiIubG9ncy5UYWdSZXNvdXJjZVJlcXVlc3QuVGFnc0VudHJ5'
    'UgR0YWdzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdm'
    'FsdWU6AjgB');

@$core.Deprecated('Use testMetricFilterRequestDescriptor instead')
const TestMetricFilterRequest$json = {
  '1': 'TestMetricFilterRequest',
  '2': [
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
    {
      '1': 'logeventmessages',
      '3': 32656852,
      '4': 3,
      '5': 9,
      '10': 'logeventmessages'
    },
  ],
};

/// Descriptor for `TestMetricFilterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testMetricFilterRequestDescriptor = $convert.base64Decode(
    'ChdUZXN0TWV0cmljRmlsdGVyUmVxdWVzdBInCg1maWx0ZXJwYXR0ZXJuGJiHikUgASgJUg1maW'
    'x0ZXJwYXR0ZXJuEi0KEGxvZ2V2ZW50bWVzc2FnZXMY1JvJDyADKAlSEGxvZ2V2ZW50bWVzc2Fn'
    'ZXM=');

@$core.Deprecated('Use testMetricFilterResponseDescriptor instead')
const TestMetricFilterResponse$json = {
  '1': 'TestMetricFilterResponse',
  '2': [
    {
      '1': 'matches',
      '3': 351545999,
      '4': 3,
      '5': 11,
      '6': '.logs.MetricFilterMatchRecord',
      '10': 'matches'
    },
  ],
};

/// Descriptor for `TestMetricFilterResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testMetricFilterResponseDescriptor =
    $convert.base64Decode(
        'ChhUZXN0TWV0cmljRmlsdGVyUmVzcG9uc2USOwoHbWF0Y2hlcxiP1dCnASADKAsyHS5sb2dzLk'
        '1ldHJpY0ZpbHRlck1hdGNoUmVjb3JkUgdtYXRjaGVz');

@$core.Deprecated('Use testTransformerRequestDescriptor instead')
const TestTransformerRequest$json = {
  '1': 'TestTransformerRequest',
  '2': [
    {
      '1': 'logeventmessages',
      '3': 32656852,
      '4': 3,
      '5': 9,
      '10': 'logeventmessages'
    },
    {
      '1': 'transformerconfig',
      '3': 384836439,
      '4': 3,
      '5': 11,
      '6': '.logs.Processor',
      '10': 'transformerconfig'
    },
  ],
};

/// Descriptor for `TestTransformerRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testTransformerRequestDescriptor = $convert.base64Decode(
    'ChZUZXN0VHJhbnNmb3JtZXJSZXF1ZXN0Ei0KEGxvZ2V2ZW50bWVzc2FnZXMY1JvJDyADKAlSEG'
    'xvZ2V2ZW50bWVzc2FnZXMSQQoRdHJhbnNmb3JtZXJjb25maWcY18bAtwEgAygLMg8ubG9ncy5Q'
    'cm9jZXNzb3JSEXRyYW5zZm9ybWVyY29uZmln');

@$core.Deprecated('Use testTransformerResponseDescriptor instead')
const TestTransformerResponse$json = {
  '1': 'TestTransformerResponse',
  '2': [
    {
      '1': 'transformedlogs',
      '3': 292768762,
      '4': 3,
      '5': 11,
      '6': '.logs.TransformedLogRecord',
      '10': 'transformedlogs'
    },
  ],
};

/// Descriptor for `TestTransformerResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testTransformerResponseDescriptor =
    $convert.base64Decode(
        'ChdUZXN0VHJhbnNmb3JtZXJSZXNwb25zZRJICg90cmFuc2Zvcm1lZGxvZ3MY+pfNiwEgAygLMh'
        'oubG9ncy5UcmFuc2Zvcm1lZExvZ1JlY29yZFIPdHJhbnNmb3JtZWRsb2dz');

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

@$core.Deprecated('Use tooManyTagsExceptionDescriptor instead')
const TooManyTagsException$json = {
  '1': 'TooManyTagsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'resourcename', '3': 17776375, '4': 1, '5': 9, '10': 'resourcename'},
  ],
};

/// Descriptor for `TooManyTagsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyTagsExceptionDescriptor = $convert.base64Decode(
    'ChRUb29NYW55VGFnc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdlEiUKDH'
    'Jlc291cmNlbmFtZRj3/bwIIAEoCVIMcmVzb3VyY2VuYW1l');

@$core.Deprecated('Use transformedLogRecordDescriptor instead')
const TransformedLogRecord$json = {
  '1': 'TransformedLogRecord',
  '2': [
    {'1': 'eventmessage', '3': 299743039, '4': 1, '5': 9, '10': 'eventmessage'},
    {'1': 'eventnumber', '3': 220470463, '4': 1, '5': 3, '10': 'eventnumber'},
    {
      '1': 'transformedeventmessage',
      '3': 209045014,
      '4': 1,
      '5': 9,
      '10': 'transformedeventmessage'
    },
  ],
};

/// Descriptor for `TransformedLogRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transformedLogRecordDescriptor = $convert.base64Decode(
    'ChRUcmFuc2Zvcm1lZExvZ1JlY29yZBImCgxldmVudG1lc3NhZ2UYv+72jgEgASgJUgxldmVudG'
    '1lc3NhZ2USIwoLZXZlbnRudW1iZXIYv7mQaSABKANSC2V2ZW50bnVtYmVyEjsKF3RyYW5zZm9y'
    'bWVkZXZlbnRtZXNzYWdlGJaM12MgASgJUhd0cmFuc2Zvcm1lZGV2ZW50bWVzc2FnZQ==');

@$core.Deprecated('Use triggerHistoryRecordDescriptor instead')
const TriggerHistoryRecord$json = {
  '1': 'TriggerHistoryRecord',
  '2': [
    {
      '1': 'destinations',
      '3': 1617189,
      '4': 3,
      '5': 11,
      '6': '.logs.ScheduledQueryDestination',
      '10': 'destinations'
    },
    {'1': 'errormessage', '3': 136873289, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'executionstatus',
      '3': 6216448,
      '4': 1,
      '5': 14,
      '6': '.logs.ExecutionStatus',
      '10': 'executionstatus'
    },
    {'1': 'queryid', '3': 336975759, '4': 1, '5': 9, '10': 'queryid'},
    {
      '1': 'triggeredtimestamp',
      '3': 257556475,
      '4': 1,
      '5': 3,
      '10': 'triggeredtimestamp'
    },
  ],
};

/// Descriptor for `TriggerHistoryRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List triggerHistoryRecordDescriptor = $convert.base64Decode(
    'ChRUcmlnZ2VySGlzdG9yeVJlY29yZBJFCgxkZXN0aW5hdGlvbnMYpdpiIAMoCzIfLmxvZ3MuU2'
    'NoZWR1bGVkUXVlcnlEZXN0aW5hdGlvblIMZGVzdGluYXRpb25zEiUKDGVycm9ybWVzc2FnZRjJ'
    'iqJBIAEoCVIMZXJyb3JtZXNzYWdlEkIKD2V4ZWN1dGlvbnN0YXR1cxiAtvsCIAEoDjIVLmxvZ3'
    'MuRXhlY3V0aW9uU3RhdHVzUg9leGVjdXRpb25zdGF0dXMSHAoHcXVlcnlpZBiPr9egASABKAlS'
    'B3F1ZXJ5aWQSMQoSdHJpZ2dlcmVkdGltZXN0YW1wGPv/53ogASgDUhJ0cmlnZ2VyZWR0aW1lc3'
    'RhbXA=');

@$core.Deprecated('Use trimStringDescriptor instead')
const TrimString$json = {
  '1': 'TrimString',
  '2': [
    {'1': 'withkeys', '3': 161392106, '4': 3, '5': 9, '10': 'withkeys'},
  ],
};

/// Descriptor for `TrimString`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trimStringDescriptor = $convert.base64Decode(
    'CgpUcmltU3RyaW5nEh0KCHdpdGhrZXlzGOrL+kwgAygJUgh3aXRoa2V5cw==');

@$core.Deprecated('Use typeConverterDescriptor instead')
const TypeConverter$json = {
  '1': 'TypeConverter',
  '2': [
    {
      '1': 'entries',
      '3': 257458932,
      '4': 3,
      '5': 11,
      '6': '.logs.TypeConverterEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `TypeConverter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List typeConverterDescriptor = $convert.base64Decode(
    'Cg1UeXBlQ29udmVydGVyEjUKB2VudHJpZXMY9IXieiADKAsyGC5sb2dzLlR5cGVDb252ZXJ0ZX'
    'JFbnRyeVIHZW50cmllcw==');

@$core.Deprecated('Use typeConverterEntryDescriptor instead')
const TypeConverterEntry$json = {
  '1': 'TypeConverterEntry',
  '2': [
    {'1': 'key', '3': 135645293, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.logs.Type',
      '10': 'type'
    },
  ],
};

/// Descriptor for `TypeConverterEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List typeConverterEntryDescriptor = $convert.base64Decode(
    'ChJUeXBlQ29udmVydGVyRW50cnkSEwoDa2V5GO2Q10AgASgJUgNrZXkSIgoEdHlwZRjO4p+JAS'
    'ABKA4yCi5sb2dzLlR5cGVSBHR5cGU=');

@$core.Deprecated('Use unrecognizedClientExceptionDescriptor instead')
const UnrecognizedClientException$json = {
  '1': 'UnrecognizedClientException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnrecognizedClientException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unrecognizedClientExceptionDescriptor =
    $convert.base64Decode(
        'ChtVbnJlY29nbml6ZWRDbGllbnRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use untagLogGroupRequestDescriptor instead')
const UntagLogGroupRequest$json = {
  '1': 'UntagLogGroupRequest',
  '2': [
    {'1': 'loggroupname', '3': 95756236, '4': 1, '5': 9, '10': 'loggroupname'},
    {'1': 'tags', '3': 337046433, '4': 3, '5': 9, '10': 'tags'},
  ],
};

/// Descriptor for `UntagLogGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagLogGroupRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ0xvZ0dyb3VwUmVxdWVzdBIlCgxsb2dncm91cG5hbWUYzL/ULSABKAlSDGxvZ2dyb3'
    'VwbmFtZRIWCgR0YWdzGKHX26ABIAMoCVIEdGFncw==');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 78811036, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIjCgtyZXNvdXJjZWFybhjNzKogIAEoCVILcmVzb3VyY2'
    'Vhcm4SGwoHdGFna2V5cxicn8olIAMoCVIHdGFna2V5cw==');

@$core.Deprecated('Use updateAnomalyRequestDescriptor instead')
const UpdateAnomalyRequest$json = {
  '1': 'UpdateAnomalyRequest',
  '2': [
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
    {'1': 'anomalyid', '3': 201902912, '4': 1, '5': 9, '10': 'anomalyid'},
    {'1': 'baseline', '3': 200859203, '4': 1, '5': 8, '10': 'baseline'},
    {'1': 'patternid', '3': 292391181, '4': 1, '5': 9, '10': 'patternid'},
    {
      '1': 'suppressionperiod',
      '3': 217163962,
      '4': 1,
      '5': 11,
      '6': '.logs.SuppressionPeriod',
      '10': 'suppressionperiod'
    },
    {
      '1': 'suppressiontype',
      '3': 322214769,
      '4': 1,
      '5': 14,
      '6': '.logs.SuppressionType',
      '10': 'suppressiontype'
    },
  ],
};

/// Descriptor for `UpdateAnomalyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAnomalyRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVBbm9tYWx5UmVxdWVzdBIyChJhbm9tYWx5ZGV0ZWN0b3Jhcm4YrM/z1AEgASgJUh'
    'Jhbm9tYWx5ZGV0ZWN0b3Jhcm4SHwoJYW5vbWFseWlkGMCWo2AgASgJUglhbm9tYWx5aWQSHQoI'
    'YmFzZWxpbmUYw7zjXyABKAhSCGJhc2VsaW5lEiAKCXBhdHRlcm5pZBiNkraLASABKAlSCXBhdH'
    'Rlcm5pZBJIChFzdXBwcmVzc2lvbnBlcmlvZBi60cZnIAEoCzIXLmxvZ3MuU3VwcHJlc3Npb25Q'
    'ZXJpb2RSEXN1cHByZXNzaW9ucGVyaW9kEkMKD3N1cHByZXNzaW9udHlwZRjxttKZASABKA4yFS'
    '5sb2dzLlN1cHByZXNzaW9uVHlwZVIPc3VwcHJlc3Npb250eXBl');

@$core.Deprecated('Use updateDeliveryConfigurationRequestDescriptor instead')
const UpdateDeliveryConfigurationRequest$json = {
  '1': 'UpdateDeliveryConfigurationRequest',
  '2': [
    {
      '1': 'fielddelimiter',
      '3': 433214845,
      '4': 1,
      '5': 9,
      '10': 'fielddelimiter'
    },
    {'1': 'id', '3': 389573345, '4': 1, '5': 9, '10': 'id'},
    {'1': 'recordfields', '3': 309625178, '4': 3, '5': 9, '10': 'recordfields'},
    {
      '1': 's3deliveryconfiguration',
      '3': 95221108,
      '4': 1,
      '5': 11,
      '6': '.logs.S3DeliveryConfiguration',
      '10': 's3deliveryconfiguration'
    },
  ],
};

/// Descriptor for `UpdateDeliveryConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDeliveryConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiJVcGRhdGVEZWxpdmVyeUNvbmZpZ3VyYXRpb25SZXF1ZXN0EioKDmZpZWxkZGVsaW1pdGVyGP'
        '2qyc4BIAEoCVIOZmllbGRkZWxpbWl0ZXISEgoCaWQY4dXhuQEgASgJUgJpZBImCgxyZWNvcmRm'
        'aWVsZHMY2oLSkwEgAygJUgxyZWNvcmRmaWVsZHMSWgoXczNkZWxpdmVyeWNvbmZpZ3VyYXRpb2'
        '4Y9OqzLSABKAsyHS5sb2dzLlMzRGVsaXZlcnlDb25maWd1cmF0aW9uUhdzM2RlbGl2ZXJ5Y29u'
        'ZmlndXJhdGlvbg==');

@$core.Deprecated('Use updateDeliveryConfigurationResponseDescriptor instead')
const UpdateDeliveryConfigurationResponse$json = {
  '1': 'UpdateDeliveryConfigurationResponse',
};

/// Descriptor for `UpdateDeliveryConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDeliveryConfigurationResponseDescriptor =
    $convert
        .base64Decode('CiNVcGRhdGVEZWxpdmVyeUNvbmZpZ3VyYXRpb25SZXNwb25zZQ==');

@$core.Deprecated('Use updateLogAnomalyDetectorRequestDescriptor instead')
const UpdateLogAnomalyDetectorRequest$json = {
  '1': 'UpdateLogAnomalyDetectorRequest',
  '2': [
    {
      '1': 'anomalydetectorarn',
      '3': 446490540,
      '4': 1,
      '5': 9,
      '10': 'anomalydetectorarn'
    },
    {
      '1': 'anomalyvisibilitytime',
      '3': 287987260,
      '4': 1,
      '5': 3,
      '10': 'anomalyvisibilitytime'
    },
    {'1': 'enabled', '3': 49525663, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'evaluationfrequency',
      '3': 533700360,
      '4': 1,
      '5': 14,
      '6': '.logs.EvaluationFrequency',
      '10': 'evaluationfrequency'
    },
    {
      '1': 'filterpattern',
      '3': 144868248,
      '4': 1,
      '5': 9,
      '10': 'filterpattern'
    },
  ],
};

/// Descriptor for `UpdateLogAnomalyDetectorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateLogAnomalyDetectorRequestDescriptor = $convert.base64Decode(
    'Ch9VcGRhdGVMb2dBbm9tYWx5RGV0ZWN0b3JSZXF1ZXN0EjIKEmFub21hbHlkZXRlY3RvcmFybh'
    'isz/PUASABKAlSEmFub21hbHlkZXRlY3RvcmFybhI4ChVhbm9tYWx5dmlzaWJpbGl0eXRpbWUY'
    'vKypiQEgASgDUhVhbm9tYWx5dmlzaWJpbGl0eXRpbWUSGwoHZW5hYmxlZBif584XIAEoCFIHZW'
    '5hYmxlZBJPChNldmFsdWF0aW9uZnJlcXVlbmN5GIi+vv4BIAEoDjIZLmxvZ3MuRXZhbHVhdGlv'
    'bkZyZXF1ZW5jeVITZXZhbHVhdGlvbmZyZXF1ZW5jeRInCg1maWx0ZXJwYXR0ZXJuGJiHikUgAS'
    'gJUg1maWx0ZXJwYXR0ZXJu');

@$core.Deprecated('Use updateScheduledQueryRequestDescriptor instead')
const UpdateScheduledQueryRequest$json = {
  '1': 'UpdateScheduledQueryRequest',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'destinationconfiguration',
      '3': 145518016,
      '4': 1,
      '5': 11,
      '6': '.logs.DestinationConfiguration',
      '10': 'destinationconfiguration'
    },
    {
      '1': 'executionrolearn',
      '3': 230613553,
      '4': 1,
      '5': 9,
      '10': 'executionrolearn'
    },
    {'1': 'identifier', '3': 145074239, '4': 1, '5': 9, '10': 'identifier'},
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'scheduleendtime',
      '3': 111645113,
      '4': 1,
      '5': 3,
      '10': 'scheduleendtime'
    },
    {
      '1': 'scheduleexpression',
      '3': 287794975,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'schedulestarttime',
      '3': 464194170,
      '4': 1,
      '5': 3,
      '10': 'schedulestarttime'
    },
    {
      '1': 'starttimeoffset',
      '3': 308525274,
      '4': 1,
      '5': 3,
      '10': 'starttimeoffset'
    },
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryState',
      '10': 'state'
    },
    {'1': 'timezone', '3': 190615331, '4': 1, '5': 9, '10': 'timezone'},
  ],
};

/// Descriptor for `UpdateScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateScheduledQueryRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVTY2hlZHVsZWRRdWVyeVJlcXVlc3QSJAoLZGVzY3JpcHRpb24Y6va8owEgASgJUg'
    'tkZXNjcmlwdGlvbhJdChhkZXN0aW5hdGlvbmNvbmZpZ3VyYXRpb24YwNuxRSABKAsyHi5sb2dz'
    'LkRlc3RpbmF0aW9uQ29uZmlndXJhdGlvblIYZGVzdGluYXRpb25jb25maWd1cmF0aW9uEi0KEG'
    'V4ZWN1dGlvbnJvbGVhcm4YscT7bSABKAlSEGV4ZWN1dGlvbnJvbGVhcm4SIQoKaWRlbnRpZmll'
    'chi/0JZFIAEoCVIKaWRlbnRpZmllchI0ChNsb2dncm91cGlkZW50aWZpZXJzGPnauMMBIAMoCV'
    'ITbG9nZ3JvdXBpZGVudGlmaWVycxI9Cg1xdWVyeWxhbmd1YWdlGPav96MBIAEoDjITLmxvZ3Mu'
    'UXVlcnlMYW5ndWFnZVINcXVlcnlsYW5ndWFnZRIkCgtxdWVyeXN0cmluZxiHgZ34ASABKAlSC3'
    'F1ZXJ5c3RyaW5nEisKD3NjaGVkdWxlZW5kdGltZRi5o541IAEoA1IPc2NoZWR1bGVlbmR0aW1l'
    'EjIKEnNjaGVkdWxlZXhwcmVzc2lvbhifzp2JASABKAlSEnNjaGVkdWxlZXhwcmVzc2lvbhIwCh'
    'FzY2hlZHVsZXN0YXJ0dGltZRj6lKzdASABKANSEXNjaGVkdWxlc3RhcnR0aW1lEiwKD3N0YXJ0'
    'dGltZW9mZnNldBja8Y6TASABKANSD3N0YXJ0dGltZW9mZnNldBIzCgVzdGF0ZRj35cTBASABKA'
    '4yGS5sb2dzLlNjaGVkdWxlZFF1ZXJ5U3RhdGVSBXN0YXRlEh0KCHRpbWV6b25lGKOe8logASgJ'
    'Ugh0aW1lem9uZQ==');

@$core.Deprecated('Use updateScheduledQueryResponseDescriptor instead')
const UpdateScheduledQueryResponse$json = {
  '1': 'UpdateScheduledQueryResponse',
  '2': [
    {'1': 'creationtime', '3': 53551750, '4': 1, '5': 3, '10': 'creationtime'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'destinationconfiguration',
      '3': 145518016,
      '4': 1,
      '5': 11,
      '6': '.logs.DestinationConfiguration',
      '10': 'destinationconfiguration'
    },
    {
      '1': 'executionrolearn',
      '3': 230613553,
      '4': 1,
      '5': 9,
      '10': 'executionrolearn'
    },
    {
      '1': 'lastexecutionstatus',
      '3': 7800596,
      '4': 1,
      '5': 14,
      '6': '.logs.ExecutionStatus',
      '10': 'lastexecutionstatus'
    },
    {
      '1': 'lasttriggeredtime',
      '3': 397057656,
      '4': 1,
      '5': 3,
      '10': 'lasttriggeredtime'
    },
    {
      '1': 'lastupdatedtime',
      '3': 388324342,
      '4': 1,
      '5': 3,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'loggroupidentifiers',
      '3': 409873785,
      '4': 3,
      '5': 9,
      '10': 'loggroupidentifiers'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'querylanguage',
      '3': 343791606,
      '4': 1,
      '5': 14,
      '6': '.logs.QueryLanguage',
      '10': 'querylanguage'
    },
    {'1': 'querystring', '3': 520568967, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'scheduleendtime',
      '3': 111645113,
      '4': 1,
      '5': 3,
      '10': 'scheduleendtime'
    },
    {
      '1': 'scheduleexpression',
      '3': 287794975,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'schedulestarttime',
      '3': 464194170,
      '4': 1,
      '5': 3,
      '10': 'schedulestarttime'
    },
    {
      '1': 'scheduledqueryarn',
      '3': 240292916,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
    {
      '1': 'starttimeoffset',
      '3': 308525274,
      '4': 1,
      '5': 3,
      '10': 'starttimeoffset'
    },
    {
      '1': 'state',
      '3': 405877495,
      '4': 1,
      '5': 14,
      '6': '.logs.ScheduledQueryState',
      '10': 'state'
    },
    {'1': 'timezone', '3': 190615331, '4': 1, '5': 9, '10': 'timezone'},
  ],
};

/// Descriptor for `UpdateScheduledQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateScheduledQueryResponseDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVTY2hlZHVsZWRRdWVyeVJlc3BvbnNlEiUKDGNyZWF0aW9udGltZRiGxcQZIAEoA1'
    'IMY3JlYXRpb250aW1lEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SXQoY'
    'ZGVzdGluYXRpb25jb25maWd1cmF0aW9uGMDbsUUgASgLMh4ubG9ncy5EZXN0aW5hdGlvbkNvbm'
    'ZpZ3VyYXRpb25SGGRlc3RpbmF0aW9uY29uZmlndXJhdGlvbhItChBleGVjdXRpb25yb2xlYXJu'
    'GLHE+20gASgJUhBleGVjdXRpb25yb2xlYXJuEkoKE2xhc3RleGVjdXRpb25zdGF0dXMYlI7cAy'
    'ABKA4yFS5sb2dzLkV4ZWN1dGlvblN0YXR1c1ITbGFzdGV4ZWN1dGlvbnN0YXR1cxIwChFsYXN0'
    'dHJpZ2dlcmVkdGltZRj4vKq9ASABKANSEWxhc3R0cmlnZ2VyZWR0aW1lEiwKD2xhc3R1cGRhdG'
    'VkdGltZRj2t5W5ASABKANSD2xhc3R1cGRhdGVkdGltZRI0ChNsb2dncm91cGlkZW50aWZpZXJz'
    'GPnauMMBIAMoCVITbG9nZ3JvdXBpZGVudGlmaWVycxIVCgRuYW1lGOf75mkgASgJUgRuYW1lEj'
    '0KDXF1ZXJ5bGFuZ3VhZ2UY9q/3owEgASgOMhMubG9ncy5RdWVyeUxhbmd1YWdlUg1xdWVyeWxh'
    'bmd1YWdlEiQKC3F1ZXJ5c3RyaW5nGIeBnfgBIAEoCVILcXVlcnlzdHJpbmcSKwoPc2NoZWR1bG'
    'VlbmR0aW1lGLmjnjUgASgDUg9zY2hlZHVsZWVuZHRpbWUSMgoSc2NoZWR1bGVleHByZXNzaW9u'
    'GJ/OnYkBIAEoCVISc2NoZWR1bGVleHByZXNzaW9uEjAKEXNjaGVkdWxlc3RhcnR0aW1lGPqUrN'
    '0BIAEoA1IRc2NoZWR1bGVzdGFydHRpbWUSLwoRc2NoZWR1bGVkcXVlcnlhcm4YtKjKciABKAlS'
    'EXNjaGVkdWxlZHF1ZXJ5YXJuEiwKD3N0YXJ0dGltZW9mZnNldBja8Y6TASABKANSD3N0YXJ0dG'
    'ltZW9mZnNldBIzCgVzdGF0ZRj35cTBASABKA4yGS5sb2dzLlNjaGVkdWxlZFF1ZXJ5U3RhdGVS'
    'BXN0YXRlEh0KCHRpbWV6b25lGKOe8logASgJUgh0aW1lem9uZQ==');

@$core.Deprecated('Use upperCaseStringDescriptor instead')
const UpperCaseString$json = {
  '1': 'UpperCaseString',
  '2': [
    {'1': 'withkeys', '3': 161392106, '4': 3, '5': 9, '10': 'withkeys'},
  ],
};

/// Descriptor for `UpperCaseString`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List upperCaseStringDescriptor = $convert.base64Decode(
    'Cg9VcHBlckNhc2VTdHJpbmcSHQoId2l0aGtleXMY6sv6TCADKAlSCHdpdGhrZXlz');

@$core.Deprecated('Use validationExceptionDescriptor instead')
const ValidationException$json = {
  '1': 'ValidationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationExceptionDescriptor =
    $convert.base64Decode(
        'ChNWYWxpZGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');
