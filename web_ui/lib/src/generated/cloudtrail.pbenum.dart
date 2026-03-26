// This is a generated file - do not edit.
//
// Generated from cloudtrail.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class BillingMode extends $pb.ProtobufEnum {
  static const BillingMode BILLING_MODE_FIXED_RETENTION_PRICING = BillingMode._(
      0, _omitEnumNames ? '' : 'BILLING_MODE_FIXED_RETENTION_PRICING');
  static const BillingMode BILLING_MODE_EXTENDABLE_RETENTION_PRICING =
      BillingMode._(
          1, _omitEnumNames ? '' : 'BILLING_MODE_EXTENDABLE_RETENTION_PRICING');

  static const $core.List<BillingMode> values = <BillingMode>[
    BILLING_MODE_FIXED_RETENTION_PRICING,
    BILLING_MODE_EXTENDABLE_RETENTION_PRICING,
  ];

  static final $core.List<BillingMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BillingMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BillingMode._(super.value, super.name);
}

class DashboardStatus extends $pb.ProtobufEnum {
  static const DashboardStatus DASHBOARD_STATUS_UPDATING =
      DashboardStatus._(0, _omitEnumNames ? '' : 'DASHBOARD_STATUS_UPDATING');
  static const DashboardStatus DASHBOARD_STATUS_UPDATED =
      DashboardStatus._(1, _omitEnumNames ? '' : 'DASHBOARD_STATUS_UPDATED');
  static const DashboardStatus DASHBOARD_STATUS_DELETING =
      DashboardStatus._(2, _omitEnumNames ? '' : 'DASHBOARD_STATUS_DELETING');
  static const DashboardStatus DASHBOARD_STATUS_CREATING =
      DashboardStatus._(3, _omitEnumNames ? '' : 'DASHBOARD_STATUS_CREATING');
  static const DashboardStatus DASHBOARD_STATUS_CREATED =
      DashboardStatus._(4, _omitEnumNames ? '' : 'DASHBOARD_STATUS_CREATED');

  static const $core.List<DashboardStatus> values = <DashboardStatus>[
    DASHBOARD_STATUS_UPDATING,
    DASHBOARD_STATUS_UPDATED,
    DASHBOARD_STATUS_DELETING,
    DASHBOARD_STATUS_CREATING,
    DASHBOARD_STATUS_CREATED,
  ];

  static final $core.List<DashboardStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static DashboardStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DashboardStatus._(super.value, super.name);
}

class DashboardType extends $pb.ProtobufEnum {
  static const DashboardType DASHBOARD_TYPE_CUSTOM =
      DashboardType._(0, _omitEnumNames ? '' : 'DASHBOARD_TYPE_CUSTOM');
  static const DashboardType DASHBOARD_TYPE_MANAGED =
      DashboardType._(1, _omitEnumNames ? '' : 'DASHBOARD_TYPE_MANAGED');

  static const $core.List<DashboardType> values = <DashboardType>[
    DASHBOARD_TYPE_CUSTOM,
    DASHBOARD_TYPE_MANAGED,
  ];

  static final $core.List<DashboardType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DashboardType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DashboardType._(super.value, super.name);
}

class DeliveryStatus extends $pb.ProtobufEnum {
  static const DeliveryStatus DELIVERY_STATUS_PENDING =
      DeliveryStatus._(0, _omitEnumNames ? '' : 'DELIVERY_STATUS_PENDING');
  static const DeliveryStatus DELIVERY_STATUS_ACCESS_DENIED = DeliveryStatus._(
      1, _omitEnumNames ? '' : 'DELIVERY_STATUS_ACCESS_DENIED');
  static const DeliveryStatus DELIVERY_STATUS_UNKNOWN =
      DeliveryStatus._(2, _omitEnumNames ? '' : 'DELIVERY_STATUS_UNKNOWN');
  static const DeliveryStatus DELIVERY_STATUS_ACCESS_DENIED_SIGNING_FILE =
      DeliveryStatus._(3,
          _omitEnumNames ? '' : 'DELIVERY_STATUS_ACCESS_DENIED_SIGNING_FILE');
  static const DeliveryStatus DELIVERY_STATUS_SUCCESS =
      DeliveryStatus._(4, _omitEnumNames ? '' : 'DELIVERY_STATUS_SUCCESS');
  static const DeliveryStatus DELIVERY_STATUS_RESOURCE_NOT_FOUND =
      DeliveryStatus._(
          5, _omitEnumNames ? '' : 'DELIVERY_STATUS_RESOURCE_NOT_FOUND');
  static const DeliveryStatus DELIVERY_STATUS_FAILED_SIGNING_FILE =
      DeliveryStatus._(
          6, _omitEnumNames ? '' : 'DELIVERY_STATUS_FAILED_SIGNING_FILE');
  static const DeliveryStatus DELIVERY_STATUS_CANCELLED =
      DeliveryStatus._(7, _omitEnumNames ? '' : 'DELIVERY_STATUS_CANCELLED');
  static const DeliveryStatus DELIVERY_STATUS_FAILED =
      DeliveryStatus._(8, _omitEnumNames ? '' : 'DELIVERY_STATUS_FAILED');

  static const $core.List<DeliveryStatus> values = <DeliveryStatus>[
    DELIVERY_STATUS_PENDING,
    DELIVERY_STATUS_ACCESS_DENIED,
    DELIVERY_STATUS_UNKNOWN,
    DELIVERY_STATUS_ACCESS_DENIED_SIGNING_FILE,
    DELIVERY_STATUS_SUCCESS,
    DELIVERY_STATUS_RESOURCE_NOT_FOUND,
    DELIVERY_STATUS_FAILED_SIGNING_FILE,
    DELIVERY_STATUS_CANCELLED,
    DELIVERY_STATUS_FAILED,
  ];

  static final $core.List<DeliveryStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static DeliveryStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeliveryStatus._(super.value, super.name);
}

class DestinationType extends $pb.ProtobufEnum {
  static const DestinationType DESTINATION_TYPE_EVENT_DATA_STORE =
      DestinationType._(
          0, _omitEnumNames ? '' : 'DESTINATION_TYPE_EVENT_DATA_STORE');
  static const DestinationType DESTINATION_TYPE_AWS_SERVICE = DestinationType._(
      1, _omitEnumNames ? '' : 'DESTINATION_TYPE_AWS_SERVICE');

  static const $core.List<DestinationType> values = <DestinationType>[
    DESTINATION_TYPE_EVENT_DATA_STORE,
    DESTINATION_TYPE_AWS_SERVICE,
  ];

  static final $core.List<DestinationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DestinationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DestinationType._(super.value, super.name);
}

class EventCategory extends $pb.ProtobufEnum {
  static const EventCategory EVENT_CATEGORY_INSIGHT =
      EventCategory._(0, _omitEnumNames ? '' : 'EVENT_CATEGORY_INSIGHT');

  static const $core.List<EventCategory> values = <EventCategory>[
    EVENT_CATEGORY_INSIGHT,
  ];

  static final $core.List<EventCategory?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static EventCategory? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventCategory._(super.value, super.name);
}

class EventCategoryAggregation extends $pb.ProtobufEnum {
  static const EventCategoryAggregation EVENT_CATEGORY_AGGREGATION_DATA =
      EventCategoryAggregation._(
          0, _omitEnumNames ? '' : 'EVENT_CATEGORY_AGGREGATION_DATA');

  static const $core.List<EventCategoryAggregation> values =
      <EventCategoryAggregation>[
    EVENT_CATEGORY_AGGREGATION_DATA,
  ];

  static final $core.List<EventCategoryAggregation?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static EventCategoryAggregation? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventCategoryAggregation._(super.value, super.name);
}

class EventDataStoreStatus extends $pb.ProtobufEnum {
  static const EventDataStoreStatus EVENT_DATA_STORE_STATUS_PENDING_DELETION =
      EventDataStoreStatus._(
          0, _omitEnumNames ? '' : 'EVENT_DATA_STORE_STATUS_PENDING_DELETION');
  static const EventDataStoreStatus EVENT_DATA_STORE_STATUS_STOPPED_INGESTION =
      EventDataStoreStatus._(
          1, _omitEnumNames ? '' : 'EVENT_DATA_STORE_STATUS_STOPPED_INGESTION');
  static const EventDataStoreStatus EVENT_DATA_STORE_STATUS_STOPPING_INGESTION =
      EventDataStoreStatus._(2,
          _omitEnumNames ? '' : 'EVENT_DATA_STORE_STATUS_STOPPING_INGESTION');
  static const EventDataStoreStatus EVENT_DATA_STORE_STATUS_STARTING_INGESTION =
      EventDataStoreStatus._(3,
          _omitEnumNames ? '' : 'EVENT_DATA_STORE_STATUS_STARTING_INGESTION');
  static const EventDataStoreStatus EVENT_DATA_STORE_STATUS_ENABLED =
      EventDataStoreStatus._(
          4, _omitEnumNames ? '' : 'EVENT_DATA_STORE_STATUS_ENABLED');
  static const EventDataStoreStatus EVENT_DATA_STORE_STATUS_CREATED =
      EventDataStoreStatus._(
          5, _omitEnumNames ? '' : 'EVENT_DATA_STORE_STATUS_CREATED');

  static const $core.List<EventDataStoreStatus> values = <EventDataStoreStatus>[
    EVENT_DATA_STORE_STATUS_PENDING_DELETION,
    EVENT_DATA_STORE_STATUS_STOPPED_INGESTION,
    EVENT_DATA_STORE_STATUS_STOPPING_INGESTION,
    EVENT_DATA_STORE_STATUS_STARTING_INGESTION,
    EVENT_DATA_STORE_STATUS_ENABLED,
    EVENT_DATA_STORE_STATUS_CREATED,
  ];

  static final $core.List<EventDataStoreStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static EventDataStoreStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventDataStoreStatus._(super.value, super.name);
}

class FederationStatus extends $pb.ProtobufEnum {
  static const FederationStatus FEDERATION_STATUS_DISABLED =
      FederationStatus._(0, _omitEnumNames ? '' : 'FEDERATION_STATUS_DISABLED');
  static const FederationStatus FEDERATION_STATUS_ENABLING =
      FederationStatus._(1, _omitEnumNames ? '' : 'FEDERATION_STATUS_ENABLING');
  static const FederationStatus FEDERATION_STATUS_DISABLING =
      FederationStatus._(
          2, _omitEnumNames ? '' : 'FEDERATION_STATUS_DISABLING');
  static const FederationStatus FEDERATION_STATUS_ENABLED =
      FederationStatus._(3, _omitEnumNames ? '' : 'FEDERATION_STATUS_ENABLED');

  static const $core.List<FederationStatus> values = <FederationStatus>[
    FEDERATION_STATUS_DISABLED,
    FEDERATION_STATUS_ENABLING,
    FEDERATION_STATUS_DISABLING,
    FEDERATION_STATUS_ENABLED,
  ];

  static final $core.List<FederationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static FederationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FederationStatus._(super.value, super.name);
}

class ImportFailureStatus extends $pb.ProtobufEnum {
  static const ImportFailureStatus IMPORT_FAILURE_STATUS_RETRY =
      ImportFailureStatus._(
          0, _omitEnumNames ? '' : 'IMPORT_FAILURE_STATUS_RETRY');
  static const ImportFailureStatus IMPORT_FAILURE_STATUS_SUCCEEDED =
      ImportFailureStatus._(
          1, _omitEnumNames ? '' : 'IMPORT_FAILURE_STATUS_SUCCEEDED');
  static const ImportFailureStatus IMPORT_FAILURE_STATUS_FAILED =
      ImportFailureStatus._(
          2, _omitEnumNames ? '' : 'IMPORT_FAILURE_STATUS_FAILED');

  static const $core.List<ImportFailureStatus> values = <ImportFailureStatus>[
    IMPORT_FAILURE_STATUS_RETRY,
    IMPORT_FAILURE_STATUS_SUCCEEDED,
    IMPORT_FAILURE_STATUS_FAILED,
  ];

  static final $core.List<ImportFailureStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ImportFailureStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportFailureStatus._(super.value, super.name);
}

class ImportStatus extends $pb.ProtobufEnum {
  static const ImportStatus IMPORT_STATUS_STOPPED =
      ImportStatus._(0, _omitEnumNames ? '' : 'IMPORT_STATUS_STOPPED');
  static const ImportStatus IMPORT_STATUS_INITIALIZING =
      ImportStatus._(1, _omitEnumNames ? '' : 'IMPORT_STATUS_INITIALIZING');
  static const ImportStatus IMPORT_STATUS_IN_PROGRESS =
      ImportStatus._(2, _omitEnumNames ? '' : 'IMPORT_STATUS_IN_PROGRESS');
  static const ImportStatus IMPORT_STATUS_COMPLETED =
      ImportStatus._(3, _omitEnumNames ? '' : 'IMPORT_STATUS_COMPLETED');
  static const ImportStatus IMPORT_STATUS_FAILED =
      ImportStatus._(4, _omitEnumNames ? '' : 'IMPORT_STATUS_FAILED');

  static const $core.List<ImportStatus> values = <ImportStatus>[
    IMPORT_STATUS_STOPPED,
    IMPORT_STATUS_INITIALIZING,
    IMPORT_STATUS_IN_PROGRESS,
    IMPORT_STATUS_COMPLETED,
    IMPORT_STATUS_FAILED,
  ];

  static final $core.List<ImportStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ImportStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportStatus._(super.value, super.name);
}

class InsightType extends $pb.ProtobufEnum {
  static const InsightType INSIGHT_TYPE_APICALLRATEINSIGHT =
      InsightType._(0, _omitEnumNames ? '' : 'INSIGHT_TYPE_APICALLRATEINSIGHT');
  static const InsightType INSIGHT_TYPE_APIERRORRATEINSIGHT = InsightType._(
      1, _omitEnumNames ? '' : 'INSIGHT_TYPE_APIERRORRATEINSIGHT');

  static const $core.List<InsightType> values = <InsightType>[
    INSIGHT_TYPE_APICALLRATEINSIGHT,
    INSIGHT_TYPE_APIERRORRATEINSIGHT,
  ];

  static final $core.List<InsightType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InsightType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InsightType._(super.value, super.name);
}

class InsightsMetricDataType extends $pb.ProtobufEnum {
  static const InsightsMetricDataType
      INSIGHTS_METRIC_DATA_TYPE_FILL_WITH_ZEROS = InsightsMetricDataType._(
          0, _omitEnumNames ? '' : 'INSIGHTS_METRIC_DATA_TYPE_FILL_WITH_ZEROS');
  static const InsightsMetricDataType INSIGHTS_METRIC_DATA_TYPE_NON_ZERO_DATA =
      InsightsMetricDataType._(
          1, _omitEnumNames ? '' : 'INSIGHTS_METRIC_DATA_TYPE_NON_ZERO_DATA');

  static const $core.List<InsightsMetricDataType> values =
      <InsightsMetricDataType>[
    INSIGHTS_METRIC_DATA_TYPE_FILL_WITH_ZEROS,
    INSIGHTS_METRIC_DATA_TYPE_NON_ZERO_DATA,
  ];

  static final $core.List<InsightsMetricDataType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InsightsMetricDataType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InsightsMetricDataType._(super.value, super.name);
}

class ListInsightsDataDimensionKey extends $pb.ProtobufEnum {
  static const ListInsightsDataDimensionKey
      LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_NAME =
      ListInsightsDataDimensionKey._(0,
          _omitEnumNames ? '' : 'LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_NAME');
  static const ListInsightsDataDimensionKey
      LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_SOURCE =
      ListInsightsDataDimensionKey._(
          1,
          _omitEnumNames
              ? ''
              : 'LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_SOURCE');
  static const ListInsightsDataDimensionKey
      LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_ID =
      ListInsightsDataDimensionKey._(
          2, _omitEnumNames ? '' : 'LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_ID');

  static const $core.List<ListInsightsDataDimensionKey> values =
      <ListInsightsDataDimensionKey>[
    LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_NAME,
    LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_SOURCE,
    LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_ID,
  ];

  static final $core.List<ListInsightsDataDimensionKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ListInsightsDataDimensionKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ListInsightsDataDimensionKey._(super.value, super.name);
}

class ListInsightsDataType extends $pb.ProtobufEnum {
  static const ListInsightsDataType LIST_INSIGHTS_DATA_TYPE_INSIGHTS_EVENTS =
      ListInsightsDataType._(
          0, _omitEnumNames ? '' : 'LIST_INSIGHTS_DATA_TYPE_INSIGHTS_EVENTS');

  static const $core.List<ListInsightsDataType> values = <ListInsightsDataType>[
    LIST_INSIGHTS_DATA_TYPE_INSIGHTS_EVENTS,
  ];

  static final $core.List<ListInsightsDataType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ListInsightsDataType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ListInsightsDataType._(super.value, super.name);
}

class LookupAttributeKey extends $pb.ProtobufEnum {
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_EVENT_NAME =
      LookupAttributeKey._(
          0, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_EVENT_NAME');
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_USERNAME =
      LookupAttributeKey._(
          1, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_USERNAME');
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_EVENT_SOURCE =
      LookupAttributeKey._(
          2, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_EVENT_SOURCE');
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_EVENT_ID =
      LookupAttributeKey._(
          3, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_EVENT_ID');
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_READ_ONLY =
      LookupAttributeKey._(
          4, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_READ_ONLY');
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_RESOURCE_TYPE =
      LookupAttributeKey._(
          5, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_RESOURCE_TYPE');
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_RESOURCE_NAME =
      LookupAttributeKey._(
          6, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_RESOURCE_NAME');
  static const LookupAttributeKey LOOKUP_ATTRIBUTE_KEY_ACCESS_KEY_ID =
      LookupAttributeKey._(
          7, _omitEnumNames ? '' : 'LOOKUP_ATTRIBUTE_KEY_ACCESS_KEY_ID');

  static const $core.List<LookupAttributeKey> values = <LookupAttributeKey>[
    LOOKUP_ATTRIBUTE_KEY_EVENT_NAME,
    LOOKUP_ATTRIBUTE_KEY_USERNAME,
    LOOKUP_ATTRIBUTE_KEY_EVENT_SOURCE,
    LOOKUP_ATTRIBUTE_KEY_EVENT_ID,
    LOOKUP_ATTRIBUTE_KEY_READ_ONLY,
    LOOKUP_ATTRIBUTE_KEY_RESOURCE_TYPE,
    LOOKUP_ATTRIBUTE_KEY_RESOURCE_NAME,
    LOOKUP_ATTRIBUTE_KEY_ACCESS_KEY_ID,
  ];

  static final $core.List<LookupAttributeKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static LookupAttributeKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LookupAttributeKey._(super.value, super.name);
}

class MaxEventSize extends $pb.ProtobufEnum {
  static const MaxEventSize MAX_EVENT_SIZE_LARGE =
      MaxEventSize._(0, _omitEnumNames ? '' : 'MAX_EVENT_SIZE_LARGE');
  static const MaxEventSize MAX_EVENT_SIZE_STANDARD =
      MaxEventSize._(1, _omitEnumNames ? '' : 'MAX_EVENT_SIZE_STANDARD');

  static const $core.List<MaxEventSize> values = <MaxEventSize>[
    MAX_EVENT_SIZE_LARGE,
    MAX_EVENT_SIZE_STANDARD,
  ];

  static final $core.List<MaxEventSize?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MaxEventSize? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MaxEventSize._(super.value, super.name);
}

class QueryStatus extends $pb.ProtobufEnum {
  static const QueryStatus QUERY_STATUS_FINISHED =
      QueryStatus._(0, _omitEnumNames ? '' : 'QUERY_STATUS_FINISHED');
  static const QueryStatus QUERY_STATUS_QUEUED =
      QueryStatus._(1, _omitEnumNames ? '' : 'QUERY_STATUS_QUEUED');
  static const QueryStatus QUERY_STATUS_RUNNING =
      QueryStatus._(2, _omitEnumNames ? '' : 'QUERY_STATUS_RUNNING');
  static const QueryStatus QUERY_STATUS_TIMED_OUT =
      QueryStatus._(3, _omitEnumNames ? '' : 'QUERY_STATUS_TIMED_OUT');
  static const QueryStatus QUERY_STATUS_CANCELLED =
      QueryStatus._(4, _omitEnumNames ? '' : 'QUERY_STATUS_CANCELLED');
  static const QueryStatus QUERY_STATUS_FAILED =
      QueryStatus._(5, _omitEnumNames ? '' : 'QUERY_STATUS_FAILED');

  static const $core.List<QueryStatus> values = <QueryStatus>[
    QUERY_STATUS_FINISHED,
    QUERY_STATUS_QUEUED,
    QUERY_STATUS_RUNNING,
    QUERY_STATUS_TIMED_OUT,
    QUERY_STATUS_CANCELLED,
    QUERY_STATUS_FAILED,
  ];

  static final $core.List<QueryStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static QueryStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryStatus._(super.value, super.name);
}

class ReadWriteType extends $pb.ProtobufEnum {
  static const ReadWriteType READ_WRITE_TYPE_READONLY =
      ReadWriteType._(0, _omitEnumNames ? '' : 'READ_WRITE_TYPE_READONLY');
  static const ReadWriteType READ_WRITE_TYPE_ALL =
      ReadWriteType._(1, _omitEnumNames ? '' : 'READ_WRITE_TYPE_ALL');
  static const ReadWriteType READ_WRITE_TYPE_WRITEONLY =
      ReadWriteType._(2, _omitEnumNames ? '' : 'READ_WRITE_TYPE_WRITEONLY');

  static const $core.List<ReadWriteType> values = <ReadWriteType>[
    READ_WRITE_TYPE_READONLY,
    READ_WRITE_TYPE_ALL,
    READ_WRITE_TYPE_WRITEONLY,
  ];

  static final $core.List<ReadWriteType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ReadWriteType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReadWriteType._(super.value, super.name);
}

class RefreshScheduleFrequencyUnit extends $pb.ProtobufEnum {
  static const RefreshScheduleFrequencyUnit
      REFRESH_SCHEDULE_FREQUENCY_UNIT_DAYS = RefreshScheduleFrequencyUnit._(
          0, _omitEnumNames ? '' : 'REFRESH_SCHEDULE_FREQUENCY_UNIT_DAYS');
  static const RefreshScheduleFrequencyUnit
      REFRESH_SCHEDULE_FREQUENCY_UNIT_HOURS = RefreshScheduleFrequencyUnit._(
          1, _omitEnumNames ? '' : 'REFRESH_SCHEDULE_FREQUENCY_UNIT_HOURS');

  static const $core.List<RefreshScheduleFrequencyUnit> values =
      <RefreshScheduleFrequencyUnit>[
    REFRESH_SCHEDULE_FREQUENCY_UNIT_DAYS,
    REFRESH_SCHEDULE_FREQUENCY_UNIT_HOURS,
  ];

  static final $core.List<RefreshScheduleFrequencyUnit?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RefreshScheduleFrequencyUnit? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RefreshScheduleFrequencyUnit._(super.value, super.name);
}

class RefreshScheduleStatus extends $pb.ProtobufEnum {
  static const RefreshScheduleStatus REFRESH_SCHEDULE_STATUS_DISABLED =
      RefreshScheduleStatus._(
          0, _omitEnumNames ? '' : 'REFRESH_SCHEDULE_STATUS_DISABLED');
  static const RefreshScheduleStatus REFRESH_SCHEDULE_STATUS_ENABLED =
      RefreshScheduleStatus._(
          1, _omitEnumNames ? '' : 'REFRESH_SCHEDULE_STATUS_ENABLED');

  static const $core.List<RefreshScheduleStatus> values =
      <RefreshScheduleStatus>[
    REFRESH_SCHEDULE_STATUS_DISABLED,
    REFRESH_SCHEDULE_STATUS_ENABLED,
  ];

  static final $core.List<RefreshScheduleStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RefreshScheduleStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RefreshScheduleStatus._(super.value, super.name);
}

class SourceEventCategory extends $pb.ProtobufEnum {
  static const SourceEventCategory SOURCE_EVENT_CATEGORY_DATA =
      SourceEventCategory._(
          0, _omitEnumNames ? '' : 'SOURCE_EVENT_CATEGORY_DATA');
  static const SourceEventCategory SOURCE_EVENT_CATEGORY_MANAGEMENT =
      SourceEventCategory._(
          1, _omitEnumNames ? '' : 'SOURCE_EVENT_CATEGORY_MANAGEMENT');

  static const $core.List<SourceEventCategory> values = <SourceEventCategory>[
    SOURCE_EVENT_CATEGORY_DATA,
    SOURCE_EVENT_CATEGORY_MANAGEMENT,
  ];

  static final $core.List<SourceEventCategory?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SourceEventCategory? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SourceEventCategory._(super.value, super.name);
}

class Template extends $pb.ProtobufEnum {
  static const Template TEMPLATE_API_ACTIVITY =
      Template._(0, _omitEnumNames ? '' : 'TEMPLATE_API_ACTIVITY');
  static const Template TEMPLATE_RESOURCE_ACCESS =
      Template._(1, _omitEnumNames ? '' : 'TEMPLATE_RESOURCE_ACCESS');
  static const Template TEMPLATE_USER_ACTIONS =
      Template._(2, _omitEnumNames ? '' : 'TEMPLATE_USER_ACTIONS');

  static const $core.List<Template> values = <Template>[
    TEMPLATE_API_ACTIVITY,
    TEMPLATE_RESOURCE_ACCESS,
    TEMPLATE_USER_ACTIONS,
  ];

  static final $core.List<Template?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static Template? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Template._(super.value, super.name);
}

class Type extends $pb.ProtobufEnum {
  static const Type TYPE_TAGCONTEXT =
      Type._(0, _omitEnumNames ? '' : 'TYPE_TAGCONTEXT');
  static const Type TYPE_REQUESTCONTEXT =
      Type._(1, _omitEnumNames ? '' : 'TYPE_REQUESTCONTEXT');

  static const $core.List<Type> values = <Type>[
    TYPE_TAGCONTEXT,
    TYPE_REQUESTCONTEXT,
  ];

  static final $core.List<Type?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static Type? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Type._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
