// This is a generated file - do not edit.
//
// Generated from logs.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ActionStatus extends $pb.ProtobufEnum {
  static const ActionStatus ACTION_STATUS_IN_PROGRESS =
      ActionStatus._(0, _omitEnumNames ? '' : 'ACTION_STATUS_IN_PROGRESS');
  static const ActionStatus ACTION_STATUS_COMPLETE =
      ActionStatus._(1, _omitEnumNames ? '' : 'ACTION_STATUS_COMPLETE');
  static const ActionStatus ACTION_STATUS_CLIENT_ERROR =
      ActionStatus._(2, _omitEnumNames ? '' : 'ACTION_STATUS_CLIENT_ERROR');
  static const ActionStatus ACTION_STATUS_FAILED =
      ActionStatus._(3, _omitEnumNames ? '' : 'ACTION_STATUS_FAILED');

  static const $core.List<ActionStatus> values = <ActionStatus>[
    ACTION_STATUS_IN_PROGRESS,
    ACTION_STATUS_COMPLETE,
    ACTION_STATUS_CLIENT_ERROR,
    ACTION_STATUS_FAILED,
  ];

  static final $core.List<ActionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ActionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ActionStatus._(super.value, super.name);
}

class AnomalyDetectorStatus extends $pb.ProtobufEnum {
  static const AnomalyDetectorStatus ANOMALY_DETECTOR_STATUS_ANALYZING =
      AnomalyDetectorStatus._(
          0, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_STATUS_ANALYZING');
  static const AnomalyDetectorStatus ANOMALY_DETECTOR_STATUS_INITIALIZING =
      AnomalyDetectorStatus._(
          1, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_STATUS_INITIALIZING');
  static const AnomalyDetectorStatus ANOMALY_DETECTOR_STATUS_TRAINING =
      AnomalyDetectorStatus._(
          2, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_STATUS_TRAINING');
  static const AnomalyDetectorStatus ANOMALY_DETECTOR_STATUS_DELETED =
      AnomalyDetectorStatus._(
          3, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_STATUS_DELETED');
  static const AnomalyDetectorStatus ANOMALY_DETECTOR_STATUS_FAILED =
      AnomalyDetectorStatus._(
          4, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_STATUS_FAILED');
  static const AnomalyDetectorStatus ANOMALY_DETECTOR_STATUS_PAUSED =
      AnomalyDetectorStatus._(
          5, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_STATUS_PAUSED');

  static const $core.List<AnomalyDetectorStatus> values =
      <AnomalyDetectorStatus>[
    ANOMALY_DETECTOR_STATUS_ANALYZING,
    ANOMALY_DETECTOR_STATUS_INITIALIZING,
    ANOMALY_DETECTOR_STATUS_TRAINING,
    ANOMALY_DETECTOR_STATUS_DELETED,
    ANOMALY_DETECTOR_STATUS_FAILED,
    ANOMALY_DETECTOR_STATUS_PAUSED,
  ];

  static final $core.List<AnomalyDetectorStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static AnomalyDetectorStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AnomalyDetectorStatus._(super.value, super.name);
}

class DataProtectionStatus extends $pb.ProtobufEnum {
  static const DataProtectionStatus DATA_PROTECTION_STATUS_DISABLED =
      DataProtectionStatus._(
          0, _omitEnumNames ? '' : 'DATA_PROTECTION_STATUS_DISABLED');
  static const DataProtectionStatus DATA_PROTECTION_STATUS_ACTIVATED =
      DataProtectionStatus._(
          1, _omitEnumNames ? '' : 'DATA_PROTECTION_STATUS_ACTIVATED');
  static const DataProtectionStatus DATA_PROTECTION_STATUS_ARCHIVED =
      DataProtectionStatus._(
          2, _omitEnumNames ? '' : 'DATA_PROTECTION_STATUS_ARCHIVED');
  static const DataProtectionStatus DATA_PROTECTION_STATUS_DELETED =
      DataProtectionStatus._(
          3, _omitEnumNames ? '' : 'DATA_PROTECTION_STATUS_DELETED');

  static const $core.List<DataProtectionStatus> values = <DataProtectionStatus>[
    DATA_PROTECTION_STATUS_DISABLED,
    DATA_PROTECTION_STATUS_ACTIVATED,
    DATA_PROTECTION_STATUS_ARCHIVED,
    DATA_PROTECTION_STATUS_DELETED,
  ];

  static final $core.List<DataProtectionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static DataProtectionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataProtectionStatus._(super.value, super.name);
}

class DeliveryDestinationType extends $pb.ProtobufEnum {
  static const DeliveryDestinationType DELIVERY_DESTINATION_TYPE_XRAY =
      DeliveryDestinationType._(
          0, _omitEnumNames ? '' : 'DELIVERY_DESTINATION_TYPE_XRAY');
  static const DeliveryDestinationType DELIVERY_DESTINATION_TYPE_S3 =
      DeliveryDestinationType._(
          1, _omitEnumNames ? '' : 'DELIVERY_DESTINATION_TYPE_S3');
  static const DeliveryDestinationType DELIVERY_DESTINATION_TYPE_CWL =
      DeliveryDestinationType._(
          2, _omitEnumNames ? '' : 'DELIVERY_DESTINATION_TYPE_CWL');
  static const DeliveryDestinationType DELIVERY_DESTINATION_TYPE_FH =
      DeliveryDestinationType._(
          3, _omitEnumNames ? '' : 'DELIVERY_DESTINATION_TYPE_FH');

  static const $core.List<DeliveryDestinationType> values =
      <DeliveryDestinationType>[
    DELIVERY_DESTINATION_TYPE_XRAY,
    DELIVERY_DESTINATION_TYPE_S3,
    DELIVERY_DESTINATION_TYPE_CWL,
    DELIVERY_DESTINATION_TYPE_FH,
  ];

  static final $core.List<DeliveryDestinationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static DeliveryDestinationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeliveryDestinationType._(super.value, super.name);
}

class Distribution extends $pb.ProtobufEnum {
  static const Distribution DISTRIBUTION_BYLOGSTREAM =
      Distribution._(0, _omitEnumNames ? '' : 'DISTRIBUTION_BYLOGSTREAM');
  static const Distribution DISTRIBUTION_RANDOM =
      Distribution._(1, _omitEnumNames ? '' : 'DISTRIBUTION_RANDOM');

  static const $core.List<Distribution> values = <Distribution>[
    DISTRIBUTION_BYLOGSTREAM,
    DISTRIBUTION_RANDOM,
  ];

  static final $core.List<Distribution?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static Distribution? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Distribution._(super.value, super.name);
}

class EntityRejectionErrorType extends $pb.ProtobufEnum {
  static const EntityRejectionErrorType
      ENTITY_REJECTION_ERROR_TYPE_ENTITY_SIZE_TOO_LARGE =
      EntityRejectionErrorType._(
          0,
          _omitEnumNames
              ? ''
              : 'ENTITY_REJECTION_ERROR_TYPE_ENTITY_SIZE_TOO_LARGE');
  static const EntityRejectionErrorType
      ENTITY_REJECTION_ERROR_TYPE_INVALID_KEY_ATTRIBUTE =
      EntityRejectionErrorType._(
          1,
          _omitEnumNames
              ? ''
              : 'ENTITY_REJECTION_ERROR_TYPE_INVALID_KEY_ATTRIBUTE');
  static const EntityRejectionErrorType
      ENTITY_REJECTION_ERROR_TYPE_UNSUPPORTED_LOG_GROUP_TYPE =
      EntityRejectionErrorType._(
          2,
          _omitEnumNames
              ? ''
              : 'ENTITY_REJECTION_ERROR_TYPE_UNSUPPORTED_LOG_GROUP_TYPE');
  static const EntityRejectionErrorType
      ENTITY_REJECTION_ERROR_TYPE_MISSING_REQUIRED_FIELDS =
      EntityRejectionErrorType._(
          3,
          _omitEnumNames
              ? ''
              : 'ENTITY_REJECTION_ERROR_TYPE_MISSING_REQUIRED_FIELDS');
  static const EntityRejectionErrorType
      ENTITY_REJECTION_ERROR_TYPE_INVALID_ATTRIBUTES =
      EntityRejectionErrorType._(
          4,
          _omitEnumNames
              ? ''
              : 'ENTITY_REJECTION_ERROR_TYPE_INVALID_ATTRIBUTES');
  static const EntityRejectionErrorType
      ENTITY_REJECTION_ERROR_TYPE_INVALID_TYPE_VALUE =
      EntityRejectionErrorType._(
          5,
          _omitEnumNames
              ? ''
              : 'ENTITY_REJECTION_ERROR_TYPE_INVALID_TYPE_VALUE');
  static const EntityRejectionErrorType
      ENTITY_REJECTION_ERROR_TYPE_INVALID_ENTITY = EntityRejectionErrorType._(6,
          _omitEnumNames ? '' : 'ENTITY_REJECTION_ERROR_TYPE_INVALID_ENTITY');

  static const $core.List<EntityRejectionErrorType> values =
      <EntityRejectionErrorType>[
    ENTITY_REJECTION_ERROR_TYPE_ENTITY_SIZE_TOO_LARGE,
    ENTITY_REJECTION_ERROR_TYPE_INVALID_KEY_ATTRIBUTE,
    ENTITY_REJECTION_ERROR_TYPE_UNSUPPORTED_LOG_GROUP_TYPE,
    ENTITY_REJECTION_ERROR_TYPE_MISSING_REQUIRED_FIELDS,
    ENTITY_REJECTION_ERROR_TYPE_INVALID_ATTRIBUTES,
    ENTITY_REJECTION_ERROR_TYPE_INVALID_TYPE_VALUE,
    ENTITY_REJECTION_ERROR_TYPE_INVALID_ENTITY,
  ];

  static final $core.List<EntityRejectionErrorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static EntityRejectionErrorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EntityRejectionErrorType._(super.value, super.name);
}

class EvaluationFrequency extends $pb.ProtobufEnum {
  static const EvaluationFrequency EVALUATION_FREQUENCY_ONE_MIN =
      EvaluationFrequency._(
          0, _omitEnumNames ? '' : 'EVALUATION_FREQUENCY_ONE_MIN');
  static const EvaluationFrequency EVALUATION_FREQUENCY_FIVE_MIN =
      EvaluationFrequency._(
          1, _omitEnumNames ? '' : 'EVALUATION_FREQUENCY_FIVE_MIN');
  static const EvaluationFrequency EVALUATION_FREQUENCY_ONE_HOUR =
      EvaluationFrequency._(
          2, _omitEnumNames ? '' : 'EVALUATION_FREQUENCY_ONE_HOUR');
  static const EvaluationFrequency EVALUATION_FREQUENCY_FIFTEEN_MIN =
      EvaluationFrequency._(
          3, _omitEnumNames ? '' : 'EVALUATION_FREQUENCY_FIFTEEN_MIN');
  static const EvaluationFrequency EVALUATION_FREQUENCY_THIRTY_MIN =
      EvaluationFrequency._(
          4, _omitEnumNames ? '' : 'EVALUATION_FREQUENCY_THIRTY_MIN');
  static const EvaluationFrequency EVALUATION_FREQUENCY_TEN_MIN =
      EvaluationFrequency._(
          5, _omitEnumNames ? '' : 'EVALUATION_FREQUENCY_TEN_MIN');

  static const $core.List<EvaluationFrequency> values = <EvaluationFrequency>[
    EVALUATION_FREQUENCY_ONE_MIN,
    EVALUATION_FREQUENCY_FIVE_MIN,
    EVALUATION_FREQUENCY_ONE_HOUR,
    EVALUATION_FREQUENCY_FIFTEEN_MIN,
    EVALUATION_FREQUENCY_THIRTY_MIN,
    EVALUATION_FREQUENCY_TEN_MIN,
  ];

  static final $core.List<EvaluationFrequency?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static EvaluationFrequency? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EvaluationFrequency._(super.value, super.name);
}

class EventSource extends $pb.ProtobufEnum {
  static const EventSource EVENT_SOURCE_CLOUD_TRAIL =
      EventSource._(0, _omitEnumNames ? '' : 'EVENT_SOURCE_CLOUD_TRAIL');
  static const EventSource EVENT_SOURCE_AWSWAF =
      EventSource._(1, _omitEnumNames ? '' : 'EVENT_SOURCE_AWSWAF');
  static const EventSource EVENT_SOURCE_VPC_FLOW =
      EventSource._(2, _omitEnumNames ? '' : 'EVENT_SOURCE_VPC_FLOW');
  static const EventSource EVENT_SOURCE_ROUTE53_RESOLVER =
      EventSource._(3, _omitEnumNames ? '' : 'EVENT_SOURCE_ROUTE53_RESOLVER');
  static const EventSource EVENT_SOURCE_EKS_AUDIT =
      EventSource._(4, _omitEnumNames ? '' : 'EVENT_SOURCE_EKS_AUDIT');

  static const $core.List<EventSource> values = <EventSource>[
    EVENT_SOURCE_CLOUD_TRAIL,
    EVENT_SOURCE_AWSWAF,
    EVENT_SOURCE_VPC_FLOW,
    EVENT_SOURCE_ROUTE53_RESOLVER,
    EVENT_SOURCE_EKS_AUDIT,
  ];

  static final $core.List<EventSource?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static EventSource? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventSource._(super.value, super.name);
}

class ExecutionStatus extends $pb.ProtobufEnum {
  static const ExecutionStatus EXECUTION_STATUS_FAILED =
      ExecutionStatus._(0, _omitEnumNames ? '' : 'EXECUTION_STATUS_FAILED');
  static const ExecutionStatus EXECUTION_STATUS_RUNNING =
      ExecutionStatus._(1, _omitEnumNames ? '' : 'EXECUTION_STATUS_RUNNING');
  static const ExecutionStatus EXECUTION_STATUS_COMPLETE =
      ExecutionStatus._(2, _omitEnumNames ? '' : 'EXECUTION_STATUS_COMPLETE');
  static const ExecutionStatus EXECUTION_STATUS_TIMEOUT =
      ExecutionStatus._(3, _omitEnumNames ? '' : 'EXECUTION_STATUS_TIMEOUT');
  static const ExecutionStatus EXECUTION_STATUS_INVALIDQUERY =
      ExecutionStatus._(
          4, _omitEnumNames ? '' : 'EXECUTION_STATUS_INVALIDQUERY');

  static const $core.List<ExecutionStatus> values = <ExecutionStatus>[
    EXECUTION_STATUS_FAILED,
    EXECUTION_STATUS_RUNNING,
    EXECUTION_STATUS_COMPLETE,
    EXECUTION_STATUS_TIMEOUT,
    EXECUTION_STATUS_INVALIDQUERY,
  ];

  static final $core.List<ExecutionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ExecutionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutionStatus._(super.value, super.name);
}

class ExportTaskStatusCode extends $pb.ProtobufEnum {
  static const ExportTaskStatusCode EXPORT_TASK_STATUS_CODE_PENDING =
      ExportTaskStatusCode._(
          0, _omitEnumNames ? '' : 'EXPORT_TASK_STATUS_CODE_PENDING');
  static const ExportTaskStatusCode EXPORT_TASK_STATUS_CODE_RUNNING =
      ExportTaskStatusCode._(
          1, _omitEnumNames ? '' : 'EXPORT_TASK_STATUS_CODE_RUNNING');
  static const ExportTaskStatusCode EXPORT_TASK_STATUS_CODE_CANCELLED =
      ExportTaskStatusCode._(
          2, _omitEnumNames ? '' : 'EXPORT_TASK_STATUS_CODE_CANCELLED');
  static const ExportTaskStatusCode EXPORT_TASK_STATUS_CODE_PENDING_CANCEL =
      ExportTaskStatusCode._(
          3, _omitEnumNames ? '' : 'EXPORT_TASK_STATUS_CODE_PENDING_CANCEL');
  static const ExportTaskStatusCode EXPORT_TASK_STATUS_CODE_COMPLETED =
      ExportTaskStatusCode._(
          4, _omitEnumNames ? '' : 'EXPORT_TASK_STATUS_CODE_COMPLETED');
  static const ExportTaskStatusCode EXPORT_TASK_STATUS_CODE_FAILED =
      ExportTaskStatusCode._(
          5, _omitEnumNames ? '' : 'EXPORT_TASK_STATUS_CODE_FAILED');

  static const $core.List<ExportTaskStatusCode> values = <ExportTaskStatusCode>[
    EXPORT_TASK_STATUS_CODE_PENDING,
    EXPORT_TASK_STATUS_CODE_RUNNING,
    EXPORT_TASK_STATUS_CODE_CANCELLED,
    EXPORT_TASK_STATUS_CODE_PENDING_CANCEL,
    EXPORT_TASK_STATUS_CODE_COMPLETED,
    EXPORT_TASK_STATUS_CODE_FAILED,
  ];

  static final $core.List<ExportTaskStatusCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ExportTaskStatusCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExportTaskStatusCode._(super.value, super.name);
}

class FlattenedElement extends $pb.ProtobufEnum {
  static const FlattenedElement FLATTENED_ELEMENT_LAST =
      FlattenedElement._(0, _omitEnumNames ? '' : 'FLATTENED_ELEMENT_LAST');
  static const FlattenedElement FLATTENED_ELEMENT_FIRST =
      FlattenedElement._(1, _omitEnumNames ? '' : 'FLATTENED_ELEMENT_FIRST');

  static const $core.List<FlattenedElement> values = <FlattenedElement>[
    FLATTENED_ELEMENT_LAST,
    FLATTENED_ELEMENT_FIRST,
  ];

  static final $core.List<FlattenedElement?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FlattenedElement? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FlattenedElement._(super.value, super.name);
}

class ImportStatus extends $pb.ProtobufEnum {
  static const ImportStatus IMPORT_STATUS_IN_PROGRESS =
      ImportStatus._(0, _omitEnumNames ? '' : 'IMPORT_STATUS_IN_PROGRESS');
  static const ImportStatus IMPORT_STATUS_CANCELLED =
      ImportStatus._(1, _omitEnumNames ? '' : 'IMPORT_STATUS_CANCELLED');
  static const ImportStatus IMPORT_STATUS_COMPLETED =
      ImportStatus._(2, _omitEnumNames ? '' : 'IMPORT_STATUS_COMPLETED');
  static const ImportStatus IMPORT_STATUS_FAILED =
      ImportStatus._(3, _omitEnumNames ? '' : 'IMPORT_STATUS_FAILED');

  static const $core.List<ImportStatus> values = <ImportStatus>[
    IMPORT_STATUS_IN_PROGRESS,
    IMPORT_STATUS_CANCELLED,
    IMPORT_STATUS_COMPLETED,
    IMPORT_STATUS_FAILED,
  ];

  static final $core.List<ImportStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ImportStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportStatus._(super.value, super.name);
}

class IndexSource extends $pb.ProtobufEnum {
  static const IndexSource INDEX_SOURCE_ACCOUNT =
      IndexSource._(0, _omitEnumNames ? '' : 'INDEX_SOURCE_ACCOUNT');
  static const IndexSource INDEX_SOURCE_LOG_GROUP =
      IndexSource._(1, _omitEnumNames ? '' : 'INDEX_SOURCE_LOG_GROUP');

  static const $core.List<IndexSource> values = <IndexSource>[
    INDEX_SOURCE_ACCOUNT,
    INDEX_SOURCE_LOG_GROUP,
  ];

  static final $core.List<IndexSource?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IndexSource? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IndexSource._(super.value, super.name);
}

class IndexType extends $pb.ProtobufEnum {
  static const IndexType INDEX_TYPE_FACET =
      IndexType._(0, _omitEnumNames ? '' : 'INDEX_TYPE_FACET');
  static const IndexType INDEX_TYPE_FIELD_INDEX =
      IndexType._(1, _omitEnumNames ? '' : 'INDEX_TYPE_FIELD_INDEX');

  static const $core.List<IndexType> values = <IndexType>[
    INDEX_TYPE_FACET,
    INDEX_TYPE_FIELD_INDEX,
  ];

  static final $core.List<IndexType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IndexType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IndexType._(super.value, super.name);
}

class InheritedProperty extends $pb.ProtobufEnum {
  static const InheritedProperty INHERITED_PROPERTY_ACCOUNT_DATA_PROTECTION =
      InheritedProperty._(0,
          _omitEnumNames ? '' : 'INHERITED_PROPERTY_ACCOUNT_DATA_PROTECTION');

  static const $core.List<InheritedProperty> values = <InheritedProperty>[
    INHERITED_PROPERTY_ACCOUNT_DATA_PROTECTION,
  ];

  static final $core.List<InheritedProperty?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static InheritedProperty? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InheritedProperty._(super.value, super.name);
}

class IntegrationStatus extends $pb.ProtobufEnum {
  static const IntegrationStatus INTEGRATION_STATUS_PROVISIONING =
      IntegrationStatus._(
          0, _omitEnumNames ? '' : 'INTEGRATION_STATUS_PROVISIONING');
  static const IntegrationStatus INTEGRATION_STATUS_ACTIVE =
      IntegrationStatus._(1, _omitEnumNames ? '' : 'INTEGRATION_STATUS_ACTIVE');
  static const IntegrationStatus INTEGRATION_STATUS_FAILED =
      IntegrationStatus._(2, _omitEnumNames ? '' : 'INTEGRATION_STATUS_FAILED');

  static const $core.List<IntegrationStatus> values = <IntegrationStatus>[
    INTEGRATION_STATUS_PROVISIONING,
    INTEGRATION_STATUS_ACTIVE,
    INTEGRATION_STATUS_FAILED,
  ];

  static final $core.List<IntegrationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static IntegrationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IntegrationStatus._(super.value, super.name);
}

class IntegrationType extends $pb.ProtobufEnum {
  static const IntegrationType INTEGRATION_TYPE_OPENSEARCH =
      IntegrationType._(0, _omitEnumNames ? '' : 'INTEGRATION_TYPE_OPENSEARCH');

  static const $core.List<IntegrationType> values = <IntegrationType>[
    INTEGRATION_TYPE_OPENSEARCH,
  ];

  static final $core.List<IntegrationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static IntegrationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IntegrationType._(super.value, super.name);
}

class ListAggregateLogGroupSummariesGroupBy extends $pb.ProtobufEnum {
  static const ListAggregateLogGroupSummariesGroupBy
      LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_AND_TYPE =
      ListAggregateLogGroupSummariesGroupBy._(
          0,
          _omitEnumNames
              ? ''
              : 'LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_AND_TYPE');
  static const ListAggregateLogGroupSummariesGroupBy
      LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_TYPE_AND_FORMAT =
      ListAggregateLogGroupSummariesGroupBy._(
          1,
          _omitEnumNames
              ? ''
              : 'LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_TYPE_AND_FORMAT');

  static const $core.List<ListAggregateLogGroupSummariesGroupBy> values =
      <ListAggregateLogGroupSummariesGroupBy>[
    LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_AND_TYPE,
    LIST_AGGREGATE_LOG_GROUP_SUMMARIES_GROUP_BY_DATA_SOURCE_NAME_TYPE_AND_FORMAT,
  ];

  static final $core.List<ListAggregateLogGroupSummariesGroupBy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ListAggregateLogGroupSummariesGroupBy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ListAggregateLogGroupSummariesGroupBy._(super.value, super.name);
}

class LogGroupClass extends $pb.ProtobufEnum {
  static const LogGroupClass LOG_GROUP_CLASS_STANDARD =
      LogGroupClass._(0, _omitEnumNames ? '' : 'LOG_GROUP_CLASS_STANDARD');
  static const LogGroupClass LOG_GROUP_CLASS_DELIVERY =
      LogGroupClass._(1, _omitEnumNames ? '' : 'LOG_GROUP_CLASS_DELIVERY');
  static const LogGroupClass LOG_GROUP_CLASS_INFREQUENT_ACCESS =
      LogGroupClass._(
          2, _omitEnumNames ? '' : 'LOG_GROUP_CLASS_INFREQUENT_ACCESS');

  static const $core.List<LogGroupClass> values = <LogGroupClass>[
    LOG_GROUP_CLASS_STANDARD,
    LOG_GROUP_CLASS_DELIVERY,
    LOG_GROUP_CLASS_INFREQUENT_ACCESS,
  ];

  static final $core.List<LogGroupClass?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static LogGroupClass? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogGroupClass._(super.value, super.name);
}

class OCSFVersion extends $pb.ProtobufEnum {
  static const OCSFVersion O_C_S_F_VERSION_V1_5 =
      OCSFVersion._(0, _omitEnumNames ? '' : 'O_C_S_F_VERSION_V1_5');
  static const OCSFVersion O_C_S_F_VERSION_V1_1 =
      OCSFVersion._(1, _omitEnumNames ? '' : 'O_C_S_F_VERSION_V1_1');

  static const $core.List<OCSFVersion> values = <OCSFVersion>[
    O_C_S_F_VERSION_V1_5,
    O_C_S_F_VERSION_V1_1,
  ];

  static final $core.List<OCSFVersion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static OCSFVersion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OCSFVersion._(super.value, super.name);
}

class OpenSearchResourceStatusType extends $pb.ProtobufEnum {
  static const OpenSearchResourceStatusType
      OPEN_SEARCH_RESOURCE_STATUS_TYPE_NOT_FOUND =
      OpenSearchResourceStatusType._(0,
          _omitEnumNames ? '' : 'OPEN_SEARCH_RESOURCE_STATUS_TYPE_NOT_FOUND');
  static const OpenSearchResourceStatusType
      OPEN_SEARCH_RESOURCE_STATUS_TYPE_ACTIVE = OpenSearchResourceStatusType._(
          1, _omitEnumNames ? '' : 'OPEN_SEARCH_RESOURCE_STATUS_TYPE_ACTIVE');
  static const OpenSearchResourceStatusType
      OPEN_SEARCH_RESOURCE_STATUS_TYPE_ERROR = OpenSearchResourceStatusType._(
          2, _omitEnumNames ? '' : 'OPEN_SEARCH_RESOURCE_STATUS_TYPE_ERROR');

  static const $core.List<OpenSearchResourceStatusType> values =
      <OpenSearchResourceStatusType>[
    OPEN_SEARCH_RESOURCE_STATUS_TYPE_NOT_FOUND,
    OPEN_SEARCH_RESOURCE_STATUS_TYPE_ACTIVE,
    OPEN_SEARCH_RESOURCE_STATUS_TYPE_ERROR,
  ];

  static final $core.List<OpenSearchResourceStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static OpenSearchResourceStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpenSearchResourceStatusType._(super.value, super.name);
}

class OrderBy extends $pb.ProtobufEnum {
  static const OrderBy ORDER_BY_LOGSTREAMNAME =
      OrderBy._(0, _omitEnumNames ? '' : 'ORDER_BY_LOGSTREAMNAME');
  static const OrderBy ORDER_BY_LASTEVENTTIME =
      OrderBy._(1, _omitEnumNames ? '' : 'ORDER_BY_LASTEVENTTIME');

  static const $core.List<OrderBy> values = <OrderBy>[
    ORDER_BY_LOGSTREAMNAME,
    ORDER_BY_LASTEVENTTIME,
  ];

  static final $core.List<OrderBy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static OrderBy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OrderBy._(super.value, super.name);
}

class OutputFormat extends $pb.ProtobufEnum {
  static const OutputFormat OUTPUT_FORMAT_RAW =
      OutputFormat._(0, _omitEnumNames ? '' : 'OUTPUT_FORMAT_RAW');
  static const OutputFormat OUTPUT_FORMAT_JSON =
      OutputFormat._(1, _omitEnumNames ? '' : 'OUTPUT_FORMAT_JSON');
  static const OutputFormat OUTPUT_FORMAT_PARQUET =
      OutputFormat._(2, _omitEnumNames ? '' : 'OUTPUT_FORMAT_PARQUET');
  static const OutputFormat OUTPUT_FORMAT_PLAIN =
      OutputFormat._(3, _omitEnumNames ? '' : 'OUTPUT_FORMAT_PLAIN');
  static const OutputFormat OUTPUT_FORMAT_W3C =
      OutputFormat._(4, _omitEnumNames ? '' : 'OUTPUT_FORMAT_W3C');

  static const $core.List<OutputFormat> values = <OutputFormat>[
    OUTPUT_FORMAT_RAW,
    OUTPUT_FORMAT_JSON,
    OUTPUT_FORMAT_PARQUET,
    OUTPUT_FORMAT_PLAIN,
    OUTPUT_FORMAT_W3C,
  ];

  static final $core.List<OutputFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static OutputFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OutputFormat._(super.value, super.name);
}

class PolicyScope extends $pb.ProtobufEnum {
  static const PolicyScope POLICY_SCOPE_ACCOUNT =
      PolicyScope._(0, _omitEnumNames ? '' : 'POLICY_SCOPE_ACCOUNT');
  static const PolicyScope POLICY_SCOPE_RESOURCE =
      PolicyScope._(1, _omitEnumNames ? '' : 'POLICY_SCOPE_RESOURCE');

  static const $core.List<PolicyScope> values = <PolicyScope>[
    POLICY_SCOPE_ACCOUNT,
    POLICY_SCOPE_RESOURCE,
  ];

  static final $core.List<PolicyScope?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PolicyScope? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PolicyScope._(super.value, super.name);
}

class PolicyType extends $pb.ProtobufEnum {
  static const PolicyType POLICY_TYPE_FIELD_INDEX_POLICY =
      PolicyType._(0, _omitEnumNames ? '' : 'POLICY_TYPE_FIELD_INDEX_POLICY');
  static const PolicyType POLICY_TYPE_DATA_PROTECTION_POLICY = PolicyType._(
      1, _omitEnumNames ? '' : 'POLICY_TYPE_DATA_PROTECTION_POLICY');
  static const PolicyType POLICY_TYPE_TRANSFORMER_POLICY =
      PolicyType._(2, _omitEnumNames ? '' : 'POLICY_TYPE_TRANSFORMER_POLICY');
  static const PolicyType POLICY_TYPE_METRIC_EXTRACTION_POLICY = PolicyType._(
      3, _omitEnumNames ? '' : 'POLICY_TYPE_METRIC_EXTRACTION_POLICY');
  static const PolicyType POLICY_TYPE_SUBSCRIPTION_FILTER_POLICY = PolicyType._(
      4, _omitEnumNames ? '' : 'POLICY_TYPE_SUBSCRIPTION_FILTER_POLICY');

  static const $core.List<PolicyType> values = <PolicyType>[
    POLICY_TYPE_FIELD_INDEX_POLICY,
    POLICY_TYPE_DATA_PROTECTION_POLICY,
    POLICY_TYPE_TRANSFORMER_POLICY,
    POLICY_TYPE_METRIC_EXTRACTION_POLICY,
    POLICY_TYPE_SUBSCRIPTION_FILTER_POLICY,
  ];

  static final $core.List<PolicyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static PolicyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PolicyType._(super.value, super.name);
}

class QueryLanguage extends $pb.ProtobufEnum {
  static const QueryLanguage QUERY_LANGUAGE_CWLI =
      QueryLanguage._(0, _omitEnumNames ? '' : 'QUERY_LANGUAGE_CWLI');
  static const QueryLanguage QUERY_LANGUAGE_SQL =
      QueryLanguage._(1, _omitEnumNames ? '' : 'QUERY_LANGUAGE_SQL');
  static const QueryLanguage QUERY_LANGUAGE_PPL =
      QueryLanguage._(2, _omitEnumNames ? '' : 'QUERY_LANGUAGE_PPL');

  static const $core.List<QueryLanguage> values = <QueryLanguage>[
    QUERY_LANGUAGE_CWLI,
    QUERY_LANGUAGE_SQL,
    QUERY_LANGUAGE_PPL,
  ];

  static final $core.List<QueryLanguage?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static QueryLanguage? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryLanguage._(super.value, super.name);
}

class QueryStatus extends $pb.ProtobufEnum {
  static const QueryStatus QUERY_STATUS_FAILED =
      QueryStatus._(0, _omitEnumNames ? '' : 'QUERY_STATUS_FAILED');
  static const QueryStatus QUERY_STATUS_CANCELLED =
      QueryStatus._(1, _omitEnumNames ? '' : 'QUERY_STATUS_CANCELLED');
  static const QueryStatus QUERY_STATUS_RUNNING =
      QueryStatus._(2, _omitEnumNames ? '' : 'QUERY_STATUS_RUNNING');
  static const QueryStatus QUERY_STATUS_COMPLETE =
      QueryStatus._(3, _omitEnumNames ? '' : 'QUERY_STATUS_COMPLETE');
  static const QueryStatus QUERY_STATUS_TIMEOUT =
      QueryStatus._(4, _omitEnumNames ? '' : 'QUERY_STATUS_TIMEOUT');
  static const QueryStatus QUERY_STATUS_SCHEDULED =
      QueryStatus._(5, _omitEnumNames ? '' : 'QUERY_STATUS_SCHEDULED');
  static const QueryStatus QUERY_STATUS_UNKNOWN =
      QueryStatus._(6, _omitEnumNames ? '' : 'QUERY_STATUS_UNKNOWN');

  static const $core.List<QueryStatus> values = <QueryStatus>[
    QUERY_STATUS_FAILED,
    QUERY_STATUS_CANCELLED,
    QUERY_STATUS_RUNNING,
    QUERY_STATUS_COMPLETE,
    QUERY_STATUS_TIMEOUT,
    QUERY_STATUS_SCHEDULED,
    QUERY_STATUS_UNKNOWN,
  ];

  static final $core.List<QueryStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static QueryStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryStatus._(super.value, super.name);
}

class S3TableIntegrationSourceStatus extends $pb.ProtobufEnum {
  static const S3TableIntegrationSourceStatus
      S3_TABLE_INTEGRATION_SOURCE_STATUS_UNHEALTHY =
      S3TableIntegrationSourceStatus._(0,
          _omitEnumNames ? '' : 'S3_TABLE_INTEGRATION_SOURCE_STATUS_UNHEALTHY');
  static const S3TableIntegrationSourceStatus
      S3_TABLE_INTEGRATION_SOURCE_STATUS_DATA_SOURCE_DELETE_IN_PROGRESS =
      S3TableIntegrationSourceStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'S3_TABLE_INTEGRATION_SOURCE_STATUS_DATA_SOURCE_DELETE_IN_PROGRESS');
  static const S3TableIntegrationSourceStatus
      S3_TABLE_INTEGRATION_SOURCE_STATUS_ACTIVE =
      S3TableIntegrationSourceStatus._(
          2, _omitEnumNames ? '' : 'S3_TABLE_INTEGRATION_SOURCE_STATUS_ACTIVE');
  static const S3TableIntegrationSourceStatus
      S3_TABLE_INTEGRATION_SOURCE_STATUS_FAILED =
      S3TableIntegrationSourceStatus._(
          3, _omitEnumNames ? '' : 'S3_TABLE_INTEGRATION_SOURCE_STATUS_FAILED');

  static const $core.List<S3TableIntegrationSourceStatus> values =
      <S3TableIntegrationSourceStatus>[
    S3_TABLE_INTEGRATION_SOURCE_STATUS_UNHEALTHY,
    S3_TABLE_INTEGRATION_SOURCE_STATUS_DATA_SOURCE_DELETE_IN_PROGRESS,
    S3_TABLE_INTEGRATION_SOURCE_STATUS_ACTIVE,
    S3_TABLE_INTEGRATION_SOURCE_STATUS_FAILED,
  ];

  static final $core.List<S3TableIntegrationSourceStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static S3TableIntegrationSourceStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const S3TableIntegrationSourceStatus._(super.value, super.name);
}

class ScheduledQueryDestinationType extends $pb.ProtobufEnum {
  static const ScheduledQueryDestinationType
      SCHEDULED_QUERY_DESTINATION_TYPE_S3 = ScheduledQueryDestinationType._(
          0, _omitEnumNames ? '' : 'SCHEDULED_QUERY_DESTINATION_TYPE_S3');

  static const $core.List<ScheduledQueryDestinationType> values =
      <ScheduledQueryDestinationType>[
    SCHEDULED_QUERY_DESTINATION_TYPE_S3,
  ];

  static final $core.List<ScheduledQueryDestinationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ScheduledQueryDestinationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScheduledQueryDestinationType._(super.value, super.name);
}

class ScheduledQueryState extends $pb.ProtobufEnum {
  static const ScheduledQueryState SCHEDULED_QUERY_STATE_DISABLED =
      ScheduledQueryState._(
          0, _omitEnumNames ? '' : 'SCHEDULED_QUERY_STATE_DISABLED');
  static const ScheduledQueryState SCHEDULED_QUERY_STATE_ENABLED =
      ScheduledQueryState._(
          1, _omitEnumNames ? '' : 'SCHEDULED_QUERY_STATE_ENABLED');

  static const $core.List<ScheduledQueryState> values = <ScheduledQueryState>[
    SCHEDULED_QUERY_STATE_DISABLED,
    SCHEDULED_QUERY_STATE_ENABLED,
  ];

  static final $core.List<ScheduledQueryState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ScheduledQueryState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScheduledQueryState._(super.value, super.name);
}

class Scope extends $pb.ProtobufEnum {
  static const Scope SCOPE_ALL = Scope._(0, _omitEnumNames ? '' : 'SCOPE_ALL');

  static const $core.List<Scope> values = <Scope>[
    SCOPE_ALL,
  ];

  static final $core.List<Scope?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static Scope? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Scope._(super.value, super.name);
}

class StandardUnit extends $pb.ProtobufEnum {
  static const StandardUnit STANDARD_UNIT_KILOBITSSECOND =
      StandardUnit._(0, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBITSSECOND');
  static const StandardUnit STANDARD_UNIT_COUNTSECOND =
      StandardUnit._(1, _omitEnumNames ? '' : 'STANDARD_UNIT_COUNTSECOND');
  static const StandardUnit STANDARD_UNIT_MEGABITSSECOND =
      StandardUnit._(2, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABITSSECOND');
  static const StandardUnit STANDARD_UNIT_GIGABITS =
      StandardUnit._(3, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABITS');
  static const StandardUnit STANDARD_UNIT_NONE =
      StandardUnit._(4, _omitEnumNames ? '' : 'STANDARD_UNIT_NONE');
  static const StandardUnit STANDARD_UNIT_MILLISECONDS =
      StandardUnit._(5, _omitEnumNames ? '' : 'STANDARD_UNIT_MILLISECONDS');
  static const StandardUnit STANDARD_UNIT_BYTES =
      StandardUnit._(6, _omitEnumNames ? '' : 'STANDARD_UNIT_BYTES');
  static const StandardUnit STANDARD_UNIT_GIGABYTES =
      StandardUnit._(7, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABYTES');
  static const StandardUnit STANDARD_UNIT_TERABITSSECOND =
      StandardUnit._(8, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABITSSECOND');
  static const StandardUnit STANDARD_UNIT_BITS =
      StandardUnit._(9, _omitEnumNames ? '' : 'STANDARD_UNIT_BITS');
  static const StandardUnit STANDARD_UNIT_GIGABITSSECOND =
      StandardUnit._(10, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABITSSECOND');
  static const StandardUnit STANDARD_UNIT_MEGABYTESSECOND =
      StandardUnit._(11, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABYTESSECOND');
  static const StandardUnit STANDARD_UNIT_MEGABYTES =
      StandardUnit._(12, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABYTES');
  static const StandardUnit STANDARD_UNIT_BITSSECOND =
      StandardUnit._(13, _omitEnumNames ? '' : 'STANDARD_UNIT_BITSSECOND');
  static const StandardUnit STANDARD_UNIT_GIGABYTESSECOND =
      StandardUnit._(14, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABYTESSECOND');
  static const StandardUnit STANDARD_UNIT_MEGABITS =
      StandardUnit._(15, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABITS');
  static const StandardUnit STANDARD_UNIT_SECONDS =
      StandardUnit._(16, _omitEnumNames ? '' : 'STANDARD_UNIT_SECONDS');
  static const StandardUnit STANDARD_UNIT_TERABYTESSECOND =
      StandardUnit._(17, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABYTESSECOND');
  static const StandardUnit STANDARD_UNIT_PERCENT =
      StandardUnit._(18, _omitEnumNames ? '' : 'STANDARD_UNIT_PERCENT');
  static const StandardUnit STANDARD_UNIT_COUNT =
      StandardUnit._(19, _omitEnumNames ? '' : 'STANDARD_UNIT_COUNT');
  static const StandardUnit STANDARD_UNIT_KILOBITS =
      StandardUnit._(20, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBITS');
  static const StandardUnit STANDARD_UNIT_KILOBYTES =
      StandardUnit._(21, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBYTES');
  static const StandardUnit STANDARD_UNIT_TERABYTES =
      StandardUnit._(22, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABYTES');
  static const StandardUnit STANDARD_UNIT_MICROSECONDS =
      StandardUnit._(23, _omitEnumNames ? '' : 'STANDARD_UNIT_MICROSECONDS');
  static const StandardUnit STANDARD_UNIT_KILOBYTESSECOND =
      StandardUnit._(24, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBYTESSECOND');
  static const StandardUnit STANDARD_UNIT_BYTESSECOND =
      StandardUnit._(25, _omitEnumNames ? '' : 'STANDARD_UNIT_BYTESSECOND');
  static const StandardUnit STANDARD_UNIT_TERABITS =
      StandardUnit._(26, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABITS');

  static const $core.List<StandardUnit> values = <StandardUnit>[
    STANDARD_UNIT_KILOBITSSECOND,
    STANDARD_UNIT_COUNTSECOND,
    STANDARD_UNIT_MEGABITSSECOND,
    STANDARD_UNIT_GIGABITS,
    STANDARD_UNIT_NONE,
    STANDARD_UNIT_MILLISECONDS,
    STANDARD_UNIT_BYTES,
    STANDARD_UNIT_GIGABYTES,
    STANDARD_UNIT_TERABITSSECOND,
    STANDARD_UNIT_BITS,
    STANDARD_UNIT_GIGABITSSECOND,
    STANDARD_UNIT_MEGABYTESSECOND,
    STANDARD_UNIT_MEGABYTES,
    STANDARD_UNIT_BITSSECOND,
    STANDARD_UNIT_GIGABYTESSECOND,
    STANDARD_UNIT_MEGABITS,
    STANDARD_UNIT_SECONDS,
    STANDARD_UNIT_TERABYTESSECOND,
    STANDARD_UNIT_PERCENT,
    STANDARD_UNIT_COUNT,
    STANDARD_UNIT_KILOBITS,
    STANDARD_UNIT_KILOBYTES,
    STANDARD_UNIT_TERABYTES,
    STANDARD_UNIT_MICROSECONDS,
    STANDARD_UNIT_KILOBYTESSECOND,
    STANDARD_UNIT_BYTESSECOND,
    STANDARD_UNIT_TERABITS,
  ];

  static final $core.List<StandardUnit?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 26);
  static StandardUnit? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StandardUnit._(super.value, super.name);
}

class State extends $pb.ProtobufEnum {
  static const State STATE_ACTIVE =
      State._(0, _omitEnumNames ? '' : 'STATE_ACTIVE');
  static const State STATE_SUPPRESSED =
      State._(1, _omitEnumNames ? '' : 'STATE_SUPPRESSED');
  static const State STATE_BASELINE =
      State._(2, _omitEnumNames ? '' : 'STATE_BASELINE');

  static const $core.List<State> values = <State>[
    STATE_ACTIVE,
    STATE_SUPPRESSED,
    STATE_BASELINE,
  ];

  static final $core.List<State?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static State? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const State._(super.value, super.name);
}

class SuppressionState extends $pb.ProtobufEnum {
  static const SuppressionState SUPPRESSION_STATE_UNSUPPRESSED =
      SuppressionState._(
          0, _omitEnumNames ? '' : 'SUPPRESSION_STATE_UNSUPPRESSED');
  static const SuppressionState SUPPRESSION_STATE_SUPPRESSED =
      SuppressionState._(
          1, _omitEnumNames ? '' : 'SUPPRESSION_STATE_SUPPRESSED');

  static const $core.List<SuppressionState> values = <SuppressionState>[
    SUPPRESSION_STATE_UNSUPPRESSED,
    SUPPRESSION_STATE_SUPPRESSED,
  ];

  static final $core.List<SuppressionState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SuppressionState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SuppressionState._(super.value, super.name);
}

class SuppressionType extends $pb.ProtobufEnum {
  static const SuppressionType SUPPRESSION_TYPE_INFINITE =
      SuppressionType._(0, _omitEnumNames ? '' : 'SUPPRESSION_TYPE_INFINITE');
  static const SuppressionType SUPPRESSION_TYPE_LIMITED =
      SuppressionType._(1, _omitEnumNames ? '' : 'SUPPRESSION_TYPE_LIMITED');

  static const $core.List<SuppressionType> values = <SuppressionType>[
    SUPPRESSION_TYPE_INFINITE,
    SUPPRESSION_TYPE_LIMITED,
  ];

  static final $core.List<SuppressionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SuppressionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SuppressionType._(super.value, super.name);
}

class SuppressionUnit extends $pb.ProtobufEnum {
  static const SuppressionUnit SUPPRESSION_UNIT_MINUTES =
      SuppressionUnit._(0, _omitEnumNames ? '' : 'SUPPRESSION_UNIT_MINUTES');
  static const SuppressionUnit SUPPRESSION_UNIT_SECONDS =
      SuppressionUnit._(1, _omitEnumNames ? '' : 'SUPPRESSION_UNIT_SECONDS');
  static const SuppressionUnit SUPPRESSION_UNIT_HOURS =
      SuppressionUnit._(2, _omitEnumNames ? '' : 'SUPPRESSION_UNIT_HOURS');

  static const $core.List<SuppressionUnit> values = <SuppressionUnit>[
    SUPPRESSION_UNIT_MINUTES,
    SUPPRESSION_UNIT_SECONDS,
    SUPPRESSION_UNIT_HOURS,
  ];

  static final $core.List<SuppressionUnit?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SuppressionUnit? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SuppressionUnit._(super.value, super.name);
}

class Type extends $pb.ProtobufEnum {
  static const Type TYPE_INTEGER =
      Type._(0, _omitEnumNames ? '' : 'TYPE_INTEGER');
  static const Type TYPE_STRING =
      Type._(1, _omitEnumNames ? '' : 'TYPE_STRING');
  static const Type TYPE_BOOLEAN =
      Type._(2, _omitEnumNames ? '' : 'TYPE_BOOLEAN');
  static const Type TYPE_DOUBLE =
      Type._(3, _omitEnumNames ? '' : 'TYPE_DOUBLE');

  static const $core.List<Type> values = <Type>[
    TYPE_INTEGER,
    TYPE_STRING,
    TYPE_BOOLEAN,
    TYPE_DOUBLE,
  ];

  static final $core.List<Type?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static Type? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Type._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
