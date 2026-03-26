// This is a generated file - do not edit.
//
// Generated from cloudtrail.proto.

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

@$core.Deprecated('Use billingModeDescriptor instead')
const BillingMode$json = {
  '1': 'BillingMode',
  '2': [
    {'1': 'BILLING_MODE_FIXED_RETENTION_PRICING', '2': 0},
    {'1': 'BILLING_MODE_EXTENDABLE_RETENTION_PRICING', '2': 1},
  ],
};

/// Descriptor for `BillingMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List billingModeDescriptor = $convert.base64Decode(
    'CgtCaWxsaW5nTW9kZRIoCiRCSUxMSU5HX01PREVfRklYRURfUkVURU5USU9OX1BSSUNJTkcQAB'
    'ItCilCSUxMSU5HX01PREVfRVhURU5EQUJMRV9SRVRFTlRJT05fUFJJQ0lORxAB');

@$core.Deprecated('Use dashboardStatusDescriptor instead')
const DashboardStatus$json = {
  '1': 'DashboardStatus',
  '2': [
    {'1': 'DASHBOARD_STATUS_UPDATING', '2': 0},
    {'1': 'DASHBOARD_STATUS_UPDATED', '2': 1},
    {'1': 'DASHBOARD_STATUS_DELETING', '2': 2},
    {'1': 'DASHBOARD_STATUS_CREATING', '2': 3},
    {'1': 'DASHBOARD_STATUS_CREATED', '2': 4},
  ],
};

/// Descriptor for `DashboardStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dashboardStatusDescriptor = $convert.base64Decode(
    'Cg9EYXNoYm9hcmRTdGF0dXMSHQoZREFTSEJPQVJEX1NUQVRVU19VUERBVElORxAAEhwKGERBU0'
    'hCT0FSRF9TVEFUVVNfVVBEQVRFRBABEh0KGURBU0hCT0FSRF9TVEFUVVNfREVMRVRJTkcQAhId'
    'ChlEQVNIQk9BUkRfU1RBVFVTX0NSRUFUSU5HEAMSHAoYREFTSEJPQVJEX1NUQVRVU19DUkVBVE'
    'VEEAQ=');

@$core.Deprecated('Use dashboardTypeDescriptor instead')
const DashboardType$json = {
  '1': 'DashboardType',
  '2': [
    {'1': 'DASHBOARD_TYPE_CUSTOM', '2': 0},
    {'1': 'DASHBOARD_TYPE_MANAGED', '2': 1},
  ],
};

/// Descriptor for `DashboardType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dashboardTypeDescriptor = $convert.base64Decode(
    'Cg1EYXNoYm9hcmRUeXBlEhkKFURBU0hCT0FSRF9UWVBFX0NVU1RPTRAAEhoKFkRBU0hCT0FSRF'
    '9UWVBFX01BTkFHRUQQAQ==');

@$core.Deprecated('Use deliveryStatusDescriptor instead')
const DeliveryStatus$json = {
  '1': 'DeliveryStatus',
  '2': [
    {'1': 'DELIVERY_STATUS_PENDING', '2': 0},
    {'1': 'DELIVERY_STATUS_ACCESS_DENIED', '2': 1},
    {'1': 'DELIVERY_STATUS_UNKNOWN', '2': 2},
    {'1': 'DELIVERY_STATUS_ACCESS_DENIED_SIGNING_FILE', '2': 3},
    {'1': 'DELIVERY_STATUS_SUCCESS', '2': 4},
    {'1': 'DELIVERY_STATUS_RESOURCE_NOT_FOUND', '2': 5},
    {'1': 'DELIVERY_STATUS_FAILED_SIGNING_FILE', '2': 6},
    {'1': 'DELIVERY_STATUS_CANCELLED', '2': 7},
    {'1': 'DELIVERY_STATUS_FAILED', '2': 8},
  ],
};

/// Descriptor for `DeliveryStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deliveryStatusDescriptor = $convert.base64Decode(
    'Cg5EZWxpdmVyeVN0YXR1cxIbChdERUxJVkVSWV9TVEFUVVNfUEVORElORxAAEiEKHURFTElWRV'
    'JZX1NUQVRVU19BQ0NFU1NfREVOSUVEEAESGwoXREVMSVZFUllfU1RBVFVTX1VOS05PV04QAhIu'
    'CipERUxJVkVSWV9TVEFUVVNfQUNDRVNTX0RFTklFRF9TSUdOSU5HX0ZJTEUQAxIbChdERUxJVk'
    'VSWV9TVEFUVVNfU1VDQ0VTUxAEEiYKIkRFTElWRVJZX1NUQVRVU19SRVNPVVJDRV9OT1RfRk9V'
    'TkQQBRInCiNERUxJVkVSWV9TVEFUVVNfRkFJTEVEX1NJR05JTkdfRklMRRAGEh0KGURFTElWRV'
    'JZX1NUQVRVU19DQU5DRUxMRUQQBxIaChZERUxJVkVSWV9TVEFUVVNfRkFJTEVEEAg=');

@$core.Deprecated('Use destinationTypeDescriptor instead')
const DestinationType$json = {
  '1': 'DestinationType',
  '2': [
    {'1': 'DESTINATION_TYPE_EVENT_DATA_STORE', '2': 0},
    {'1': 'DESTINATION_TYPE_AWS_SERVICE', '2': 1},
  ],
};

/// Descriptor for `DestinationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List destinationTypeDescriptor = $convert.base64Decode(
    'Cg9EZXN0aW5hdGlvblR5cGUSJQohREVTVElOQVRJT05fVFlQRV9FVkVOVF9EQVRBX1NUT1JFEA'
    'ASIAocREVTVElOQVRJT05fVFlQRV9BV1NfU0VSVklDRRAB');

@$core.Deprecated('Use eventCategoryDescriptor instead')
const EventCategory$json = {
  '1': 'EventCategory',
  '2': [
    {'1': 'EVENT_CATEGORY_INSIGHT', '2': 0},
  ],
};

/// Descriptor for `EventCategory`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventCategoryDescriptor = $convert.base64Decode(
    'Cg1FdmVudENhdGVnb3J5EhoKFkVWRU5UX0NBVEVHT1JZX0lOU0lHSFQQAA==');

@$core.Deprecated('Use eventCategoryAggregationDescriptor instead')
const EventCategoryAggregation$json = {
  '1': 'EventCategoryAggregation',
  '2': [
    {'1': 'EVENT_CATEGORY_AGGREGATION_DATA', '2': 0},
  ],
};

/// Descriptor for `EventCategoryAggregation`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventCategoryAggregationDescriptor =
    $convert.base64Decode(
        'ChhFdmVudENhdGVnb3J5QWdncmVnYXRpb24SIwofRVZFTlRfQ0FURUdPUllfQUdHUkVHQVRJT0'
        '5fREFUQRAA');

@$core.Deprecated('Use eventDataStoreStatusDescriptor instead')
const EventDataStoreStatus$json = {
  '1': 'EventDataStoreStatus',
  '2': [
    {'1': 'EVENT_DATA_STORE_STATUS_PENDING_DELETION', '2': 0},
    {'1': 'EVENT_DATA_STORE_STATUS_STOPPED_INGESTION', '2': 1},
    {'1': 'EVENT_DATA_STORE_STATUS_STOPPING_INGESTION', '2': 2},
    {'1': 'EVENT_DATA_STORE_STATUS_STARTING_INGESTION', '2': 3},
    {'1': 'EVENT_DATA_STORE_STATUS_ENABLED', '2': 4},
    {'1': 'EVENT_DATA_STORE_STATUS_CREATED', '2': 5},
  ],
};

/// Descriptor for `EventDataStoreStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventDataStoreStatusDescriptor = $convert.base64Decode(
    'ChRFdmVudERhdGFTdG9yZVN0YXR1cxIsCihFVkVOVF9EQVRBX1NUT1JFX1NUQVRVU19QRU5ESU'
    '5HX0RFTEVUSU9OEAASLQopRVZFTlRfREFUQV9TVE9SRV9TVEFUVVNfU1RPUFBFRF9JTkdFU1RJ'
    'T04QARIuCipFVkVOVF9EQVRBX1NUT1JFX1NUQVRVU19TVE9QUElOR19JTkdFU1RJT04QAhIuCi'
    'pFVkVOVF9EQVRBX1NUT1JFX1NUQVRVU19TVEFSVElOR19JTkdFU1RJT04QAxIjCh9FVkVOVF9E'
    'QVRBX1NUT1JFX1NUQVRVU19FTkFCTEVEEAQSIwofRVZFTlRfREFUQV9TVE9SRV9TVEFUVVNfQ1'
    'JFQVRFRBAF');

@$core.Deprecated('Use federationStatusDescriptor instead')
const FederationStatus$json = {
  '1': 'FederationStatus',
  '2': [
    {'1': 'FEDERATION_STATUS_DISABLED', '2': 0},
    {'1': 'FEDERATION_STATUS_ENABLING', '2': 1},
    {'1': 'FEDERATION_STATUS_DISABLING', '2': 2},
    {'1': 'FEDERATION_STATUS_ENABLED', '2': 3},
  ],
};

/// Descriptor for `FederationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List federationStatusDescriptor = $convert.base64Decode(
    'ChBGZWRlcmF0aW9uU3RhdHVzEh4KGkZFREVSQVRJT05fU1RBVFVTX0RJU0FCTEVEEAASHgoaRk'
    'VERVJBVElPTl9TVEFUVVNfRU5BQkxJTkcQARIfChtGRURFUkFUSU9OX1NUQVRVU19ESVNBQkxJ'
    'TkcQAhIdChlGRURFUkFUSU9OX1NUQVRVU19FTkFCTEVEEAM=');

@$core.Deprecated('Use importFailureStatusDescriptor instead')
const ImportFailureStatus$json = {
  '1': 'ImportFailureStatus',
  '2': [
    {'1': 'IMPORT_FAILURE_STATUS_RETRY', '2': 0},
    {'1': 'IMPORT_FAILURE_STATUS_SUCCEEDED', '2': 1},
    {'1': 'IMPORT_FAILURE_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `ImportFailureStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importFailureStatusDescriptor = $convert.base64Decode(
    'ChNJbXBvcnRGYWlsdXJlU3RhdHVzEh8KG0lNUE9SVF9GQUlMVVJFX1NUQVRVU19SRVRSWRAAEi'
    'MKH0lNUE9SVF9GQUlMVVJFX1NUQVRVU19TVUNDRUVERUQQARIgChxJTVBPUlRfRkFJTFVSRV9T'
    'VEFUVVNfRkFJTEVEEAI=');

@$core.Deprecated('Use importStatusDescriptor instead')
const ImportStatus$json = {
  '1': 'ImportStatus',
  '2': [
    {'1': 'IMPORT_STATUS_STOPPED', '2': 0},
    {'1': 'IMPORT_STATUS_INITIALIZING', '2': 1},
    {'1': 'IMPORT_STATUS_IN_PROGRESS', '2': 2},
    {'1': 'IMPORT_STATUS_COMPLETED', '2': 3},
    {'1': 'IMPORT_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `ImportStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importStatusDescriptor = $convert.base64Decode(
    'CgxJbXBvcnRTdGF0dXMSGQoVSU1QT1JUX1NUQVRVU19TVE9QUEVEEAASHgoaSU1QT1JUX1NUQV'
    'RVU19JTklUSUFMSVpJTkcQARIdChlJTVBPUlRfU1RBVFVTX0lOX1BST0dSRVNTEAISGwoXSU1Q'
    'T1JUX1NUQVRVU19DT01QTEVURUQQAxIYChRJTVBPUlRfU1RBVFVTX0ZBSUxFRBAE');

@$core.Deprecated('Use insightTypeDescriptor instead')
const InsightType$json = {
  '1': 'InsightType',
  '2': [
    {'1': 'INSIGHT_TYPE_APICALLRATEINSIGHT', '2': 0},
    {'1': 'INSIGHT_TYPE_APIERRORRATEINSIGHT', '2': 1},
  ],
};

/// Descriptor for `InsightType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List insightTypeDescriptor = $convert.base64Decode(
    'CgtJbnNpZ2h0VHlwZRIjCh9JTlNJR0hUX1RZUEVfQVBJQ0FMTFJBVEVJTlNJR0hUEAASJAogSU'
    '5TSUdIVF9UWVBFX0FQSUVSUk9SUkFURUlOU0lHSFQQAQ==');

@$core.Deprecated('Use insightsMetricDataTypeDescriptor instead')
const InsightsMetricDataType$json = {
  '1': 'InsightsMetricDataType',
  '2': [
    {'1': 'INSIGHTS_METRIC_DATA_TYPE_FILL_WITH_ZEROS', '2': 0},
    {'1': 'INSIGHTS_METRIC_DATA_TYPE_NON_ZERO_DATA', '2': 1},
  ],
};

/// Descriptor for `InsightsMetricDataType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List insightsMetricDataTypeDescriptor = $convert.base64Decode(
    'ChZJbnNpZ2h0c01ldHJpY0RhdGFUeXBlEi0KKUlOU0lHSFRTX01FVFJJQ19EQVRBX1RZUEVfRk'
    'lMTF9XSVRIX1pFUk9TEAASKwonSU5TSUdIVFNfTUVUUklDX0RBVEFfVFlQRV9OT05fWkVST19E'
    'QVRBEAE=');

@$core.Deprecated('Use listInsightsDataDimensionKeyDescriptor instead')
const ListInsightsDataDimensionKey$json = {
  '1': 'ListInsightsDataDimensionKey',
  '2': [
    {'1': 'LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_NAME', '2': 0},
    {'1': 'LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_SOURCE', '2': 1},
    {'1': 'LIST_INSIGHTS_DATA_DIMENSION_KEY_EVENT_ID', '2': 2},
  ],
};

/// Descriptor for `ListInsightsDataDimensionKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List listInsightsDataDimensionKeyDescriptor = $convert.base64Decode(
    'ChxMaXN0SW5zaWdodHNEYXRhRGltZW5zaW9uS2V5Ei8KK0xJU1RfSU5TSUdIVFNfREFUQV9ESU'
    '1FTlNJT05fS0VZX0VWRU5UX05BTUUQABIxCi1MSVNUX0lOU0lHSFRTX0RBVEFfRElNRU5TSU9O'
    'X0tFWV9FVkVOVF9TT1VSQ0UQARItCilMSVNUX0lOU0lHSFRTX0RBVEFfRElNRU5TSU9OX0tFWV'
    '9FVkVOVF9JRBAC');

@$core.Deprecated('Use listInsightsDataTypeDescriptor instead')
const ListInsightsDataType$json = {
  '1': 'ListInsightsDataType',
  '2': [
    {'1': 'LIST_INSIGHTS_DATA_TYPE_INSIGHTS_EVENTS', '2': 0},
  ],
};

/// Descriptor for `ListInsightsDataType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List listInsightsDataTypeDescriptor = $convert.base64Decode(
    'ChRMaXN0SW5zaWdodHNEYXRhVHlwZRIrCidMSVNUX0lOU0lHSFRTX0RBVEFfVFlQRV9JTlNJR0'
    'hUU19FVkVOVFMQAA==');

@$core.Deprecated('Use lookupAttributeKeyDescriptor instead')
const LookupAttributeKey$json = {
  '1': 'LookupAttributeKey',
  '2': [
    {'1': 'LOOKUP_ATTRIBUTE_KEY_EVENT_NAME', '2': 0},
    {'1': 'LOOKUP_ATTRIBUTE_KEY_USERNAME', '2': 1},
    {'1': 'LOOKUP_ATTRIBUTE_KEY_EVENT_SOURCE', '2': 2},
    {'1': 'LOOKUP_ATTRIBUTE_KEY_EVENT_ID', '2': 3},
    {'1': 'LOOKUP_ATTRIBUTE_KEY_READ_ONLY', '2': 4},
    {'1': 'LOOKUP_ATTRIBUTE_KEY_RESOURCE_TYPE', '2': 5},
    {'1': 'LOOKUP_ATTRIBUTE_KEY_RESOURCE_NAME', '2': 6},
    {'1': 'LOOKUP_ATTRIBUTE_KEY_ACCESS_KEY_ID', '2': 7},
  ],
};

/// Descriptor for `LookupAttributeKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List lookupAttributeKeyDescriptor = $convert.base64Decode(
    'ChJMb29rdXBBdHRyaWJ1dGVLZXkSIwofTE9PS1VQX0FUVFJJQlVURV9LRVlfRVZFTlRfTkFNRR'
    'AAEiEKHUxPT0tVUF9BVFRSSUJVVEVfS0VZX1VTRVJOQU1FEAESJQohTE9PS1VQX0FUVFJJQlVU'
    'RV9LRVlfRVZFTlRfU09VUkNFEAISIQodTE9PS1VQX0FUVFJJQlVURV9LRVlfRVZFTlRfSUQQAx'
    'IiCh5MT09LVVBfQVRUUklCVVRFX0tFWV9SRUFEX09OTFkQBBImCiJMT09LVVBfQVRUUklCVVRF'
    'X0tFWV9SRVNPVVJDRV9UWVBFEAUSJgoiTE9PS1VQX0FUVFJJQlVURV9LRVlfUkVTT1VSQ0VfTk'
    'FNRRAGEiYKIkxPT0tVUF9BVFRSSUJVVEVfS0VZX0FDQ0VTU19LRVlfSUQQBw==');

@$core.Deprecated('Use maxEventSizeDescriptor instead')
const MaxEventSize$json = {
  '1': 'MaxEventSize',
  '2': [
    {'1': 'MAX_EVENT_SIZE_LARGE', '2': 0},
    {'1': 'MAX_EVENT_SIZE_STANDARD', '2': 1},
  ],
};

/// Descriptor for `MaxEventSize`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List maxEventSizeDescriptor = $convert.base64Decode(
    'CgxNYXhFdmVudFNpemUSGAoUTUFYX0VWRU5UX1NJWkVfTEFSR0UQABIbChdNQVhfRVZFTlRfU0'
    'laRV9TVEFOREFSRBAB');

@$core.Deprecated('Use queryStatusDescriptor instead')
const QueryStatus$json = {
  '1': 'QueryStatus',
  '2': [
    {'1': 'QUERY_STATUS_FINISHED', '2': 0},
    {'1': 'QUERY_STATUS_QUEUED', '2': 1},
    {'1': 'QUERY_STATUS_RUNNING', '2': 2},
    {'1': 'QUERY_STATUS_TIMED_OUT', '2': 3},
    {'1': 'QUERY_STATUS_CANCELLED', '2': 4},
    {'1': 'QUERY_STATUS_FAILED', '2': 5},
  ],
};

/// Descriptor for `QueryStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryStatusDescriptor = $convert.base64Decode(
    'CgtRdWVyeVN0YXR1cxIZChVRVUVSWV9TVEFUVVNfRklOSVNIRUQQABIXChNRVUVSWV9TVEFUVV'
    'NfUVVFVUVEEAESGAoUUVVFUllfU1RBVFVTX1JVTk5JTkcQAhIaChZRVUVSWV9TVEFUVVNfVElN'
    'RURfT1VUEAMSGgoWUVVFUllfU1RBVFVTX0NBTkNFTExFRBAEEhcKE1FVRVJZX1NUQVRVU19GQU'
    'lMRUQQBQ==');

@$core.Deprecated('Use readWriteTypeDescriptor instead')
const ReadWriteType$json = {
  '1': 'ReadWriteType',
  '2': [
    {'1': 'READ_WRITE_TYPE_READONLY', '2': 0},
    {'1': 'READ_WRITE_TYPE_ALL', '2': 1},
    {'1': 'READ_WRITE_TYPE_WRITEONLY', '2': 2},
  ],
};

/// Descriptor for `ReadWriteType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List readWriteTypeDescriptor = $convert.base64Decode(
    'Cg1SZWFkV3JpdGVUeXBlEhwKGFJFQURfV1JJVEVfVFlQRV9SRUFET05MWRAAEhcKE1JFQURfV1'
    'JJVEVfVFlQRV9BTEwQARIdChlSRUFEX1dSSVRFX1RZUEVfV1JJVEVPTkxZEAI=');

@$core.Deprecated('Use refreshScheduleFrequencyUnitDescriptor instead')
const RefreshScheduleFrequencyUnit$json = {
  '1': 'RefreshScheduleFrequencyUnit',
  '2': [
    {'1': 'REFRESH_SCHEDULE_FREQUENCY_UNIT_DAYS', '2': 0},
    {'1': 'REFRESH_SCHEDULE_FREQUENCY_UNIT_HOURS', '2': 1},
  ],
};

/// Descriptor for `RefreshScheduleFrequencyUnit`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List refreshScheduleFrequencyUnitDescriptor =
    $convert.base64Decode(
        'ChxSZWZyZXNoU2NoZWR1bGVGcmVxdWVuY3lVbml0EigKJFJFRlJFU0hfU0NIRURVTEVfRlJFUV'
        'VFTkNZX1VOSVRfREFZUxAAEikKJVJFRlJFU0hfU0NIRURVTEVfRlJFUVVFTkNZX1VOSVRfSE9V'
        'UlMQAQ==');

@$core.Deprecated('Use refreshScheduleStatusDescriptor instead')
const RefreshScheduleStatus$json = {
  '1': 'RefreshScheduleStatus',
  '2': [
    {'1': 'REFRESH_SCHEDULE_STATUS_DISABLED', '2': 0},
    {'1': 'REFRESH_SCHEDULE_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `RefreshScheduleStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List refreshScheduleStatusDescriptor = $convert.base64Decode(
    'ChVSZWZyZXNoU2NoZWR1bGVTdGF0dXMSJAogUkVGUkVTSF9TQ0hFRFVMRV9TVEFUVVNfRElTQU'
    'JMRUQQABIjCh9SRUZSRVNIX1NDSEVEVUxFX1NUQVRVU19FTkFCTEVEEAE=');

@$core.Deprecated('Use sourceEventCategoryDescriptor instead')
const SourceEventCategory$json = {
  '1': 'SourceEventCategory',
  '2': [
    {'1': 'SOURCE_EVENT_CATEGORY_DATA', '2': 0},
    {'1': 'SOURCE_EVENT_CATEGORY_MANAGEMENT', '2': 1},
  ],
};

/// Descriptor for `SourceEventCategory`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sourceEventCategoryDescriptor = $convert.base64Decode(
    'ChNTb3VyY2VFdmVudENhdGVnb3J5Eh4KGlNPVVJDRV9FVkVOVF9DQVRFR09SWV9EQVRBEAASJA'
    'ogU09VUkNFX0VWRU5UX0NBVEVHT1JZX01BTkFHRU1FTlQQAQ==');

@$core.Deprecated('Use templateDescriptor instead')
const Template$json = {
  '1': 'Template',
  '2': [
    {'1': 'TEMPLATE_API_ACTIVITY', '2': 0},
    {'1': 'TEMPLATE_RESOURCE_ACCESS', '2': 1},
    {'1': 'TEMPLATE_USER_ACTIONS', '2': 2},
  ],
};

/// Descriptor for `Template`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List templateDescriptor = $convert.base64Decode(
    'CghUZW1wbGF0ZRIZChVURU1QTEFURV9BUElfQUNUSVZJVFkQABIcChhURU1QTEFURV9SRVNPVV'
    'JDRV9BQ0NFU1MQARIZChVURU1QTEFURV9VU0VSX0FDVElPTlMQAg==');

@$core.Deprecated('Use typeDescriptor instead')
const Type$json = {
  '1': 'Type',
  '2': [
    {'1': 'TYPE_TAGCONTEXT', '2': 0},
    {'1': 'TYPE_REQUESTCONTEXT', '2': 1},
  ],
};

/// Descriptor for `Type`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List typeDescriptor = $convert.base64Decode(
    'CgRUeXBlEhMKD1RZUEVfVEFHQ09OVEVYVBAAEhcKE1RZUEVfUkVRVUVTVENPTlRFWFQQAQ==');

@$core.Deprecated('Use accessDeniedExceptionDescriptor instead')
const AccessDeniedException$json = {
  '1': 'AccessDeniedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedExceptionDescriptor = $convert.base64Decode(
    'ChVBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use accountHasOngoingImportExceptionDescriptor instead')
const AccountHasOngoingImportException$json = {
  '1': 'AccountHasOngoingImportException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccountHasOngoingImportException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountHasOngoingImportExceptionDescriptor =
    $convert.base64Decode(
        'CiBBY2NvdW50SGFzT25nb2luZ0ltcG9ydEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use accountNotFoundExceptionDescriptor instead')
const AccountNotFoundException$json = {
  '1': 'AccountNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccountNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChhBY2NvdW50Tm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use accountNotRegisteredExceptionDescriptor instead')
const AccountNotRegisteredException$json = {
  '1': 'AccountNotRegisteredException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccountNotRegisteredException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountNotRegisteredExceptionDescriptor =
    $convert.base64Decode(
        'Ch1BY2NvdW50Tm90UmVnaXN0ZXJlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use accountRegisteredExceptionDescriptor instead')
const AccountRegisteredException$json = {
  '1': 'AccountRegisteredException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccountRegisteredException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountRegisteredExceptionDescriptor =
    $convert.base64Decode(
        'ChpBY2NvdW50UmVnaXN0ZXJlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use addTagsRequestDescriptor instead')
const AddTagsRequest$json = {
  '1': 'AddTagsRequest',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
  ],
};

/// Descriptor for `AddTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addTagsRequestDescriptor = $convert.base64Decode(
    'Cg5BZGRUYWdzUmVxdWVzdBIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcmVzb3VyY2VpZBIvCg'
    'h0YWdzbGlzdBippJjtASADKAsyDy5jbG91ZHRyYWlsLlRhZ1IIdGFnc2xpc3Q=');

@$core.Deprecated('Use addTagsResponseDescriptor instead')
const AddTagsResponse$json = {
  '1': 'AddTagsResponse',
};

/// Descriptor for `AddTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addTagsResponseDescriptor =
    $convert.base64Decode('Cg9BZGRUYWdzUmVzcG9uc2U=');

@$core.Deprecated('Use advancedEventSelectorDescriptor instead')
const AdvancedEventSelector$json = {
  '1': 'AdvancedEventSelector',
  '2': [
    {
      '1': 'fieldselectors',
      '3': 445033784,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedFieldSelector',
      '10': 'fieldselectors'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `AdvancedEventSelector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List advancedEventSelectorDescriptor = $convert.base64Decode(
    'ChVBZHZhbmNlZEV2ZW50U2VsZWN0b3ISTQoOZmllbGRzZWxlY3RvcnMYuNqa1AEgAygLMiEuY2'
    'xvdWR0cmFpbC5BZHZhbmNlZEZpZWxkU2VsZWN0b3JSDmZpZWxkc2VsZWN0b3JzEhUKBG5hbWUY'
    'h+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use advancedFieldSelectorDescriptor instead')
const AdvancedFieldSelector$json = {
  '1': 'AdvancedFieldSelector',
  '2': [
    {'1': 'endswith', '3': 505385124, '4': 3, '5': 9, '10': 'endswith'},
    {'1': 'equals', '3': 513367477, '4': 3, '5': 9, '10': 'equals'},
    {'1': 'field', '3': 263732488, '4': 1, '5': 9, '10': 'field'},
    {'1': 'notendswith', '3': 252128485, '4': 3, '5': 9, '10': 'notendswith'},
    {'1': 'notequals', '3': 424280088, '4': 3, '5': 9, '10': 'notequals'},
    {
      '1': 'notstartswith',
      '3': 305277948,
      '4': 3,
      '5': 9,
      '10': 'notstartswith'
    },
    {'1': 'startswith', '3': 468546557, '4': 3, '5': 9, '10': 'startswith'},
  ],
};

/// Descriptor for `AdvancedFieldSelector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List advancedFieldSelectorDescriptor = $convert.base64Decode(
    'ChVBZHZhbmNlZEZpZWxkU2VsZWN0b3ISHgoIZW5kc3dpdGgYpKH+8AEgAygJUghlbmRzd2l0aB'
    'IaCgZlcXVhbHMYtbvl9AEgAygJUgZlcXVhbHMSFwoFZmllbGQYiPrgfSABKAlSBWZpZWxkEiMK'
    'C25vdGVuZHN3aXRoGOXZnHggAygJUgtub3RlbmRzd2l0aBIgCglub3RlcXVhbHMYmICoygEgAy'
    'gJUglub3RlcXVhbHMSKAoNbm90c3RhcnRzd2l0aBj818iRASADKAlSDW5vdHN0YXJ0c3dpdGgS'
    'IgoKc3RhcnRzd2l0aBj957XfASADKAlSCnN0YXJ0c3dpdGg=');

@$core.Deprecated('Use aggregationConfigurationDescriptor instead')
const AggregationConfiguration$json = {
  '1': 'AggregationConfiguration',
  '2': [
    {
      '1': 'eventcategory',
      '3': 164668724,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.EventCategoryAggregation',
      '10': 'eventcategory'
    },
    {
      '1': 'templates',
      '3': 63634313,
      '4': 3,
      '5': 14,
      '6': '.cloudtrail.Template',
      '10': 'templates'
    },
  ],
};

/// Descriptor for `AggregationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aggregationConfigurationDescriptor = $convert.base64Decode(
    'ChhBZ2dyZWdhdGlvbkNvbmZpZ3VyYXRpb24STQoNZXZlbnRjYXRlZ29yeRi0ysJOIAEoDjIkLm'
    'Nsb3VkdHJhaWwuRXZlbnRDYXRlZ29yeUFnZ3JlZ2F0aW9uUg1ldmVudGNhdGVnb3J5EjUKCXRl'
    'bXBsYXRlcxiJ96seIAMoDjIULmNsb3VkdHJhaWwuVGVtcGxhdGVSCXRlbXBsYXRlcw==');

@$core.Deprecated('Use cancelQueryRequestDescriptor instead')
const CancelQueryRequest$json = {
  '1': 'CancelQueryRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `CancelQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelQueryRequestDescriptor = $convert.base64Decode(
    'ChJDYW5jZWxRdWVyeVJlcXVlc3QSKQoOZXZlbnRkYXRhc3RvcmUYwdudQSABKAlSDmV2ZW50ZG'
    'F0YXN0b3JlEkYKHGV2ZW50ZGF0YXN0b3Jlb3duZXJhY2NvdW50aWQYsN3w4AEgASgJUhxldmVu'
    'dGRhdGFzdG9yZW93bmVyYWNjb3VudGlkEhsKB3F1ZXJ5aWQY7/DmNCABKAlSB3F1ZXJ5aWQ=');

@$core.Deprecated('Use cancelQueryResponseDescriptor instead')
const CancelQueryResponse$json = {
  '1': 'CancelQueryResponse',
  '2': [
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
    {
      '1': 'querystatus',
      '3': 367016406,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.QueryStatus',
      '10': 'querystatus'
    },
  ],
};

/// Descriptor for `CancelQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelQueryResponseDescriptor = $convert.base64Decode(
    'ChNDYW5jZWxRdWVyeVJlc3BvbnNlEkYKHGV2ZW50ZGF0YXN0b3Jlb3duZXJhY2NvdW50aWQYsN'
    '3w4AEgASgJUhxldmVudGRhdGFzdG9yZW93bmVyYWNjb3VudGlkEhsKB3F1ZXJ5aWQY7/DmNCAB'
    'KAlSB3F1ZXJ5aWQSPQoLcXVlcnlzdGF0dXMY1vOArwEgASgOMhcuY2xvdWR0cmFpbC5RdWVyeV'
    'N0YXR1c1ILcXVlcnlzdGF0dXM=');

@$core.Deprecated(
    'Use cannotDelegateManagementAccountExceptionDescriptor instead')
const CannotDelegateManagementAccountException$json = {
  '1': 'CannotDelegateManagementAccountException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CannotDelegateManagementAccountException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cannotDelegateManagementAccountExceptionDescriptor =
    $convert.base64Decode(
        'CihDYW5ub3REZWxlZ2F0ZU1hbmFnZW1lbnRBY2NvdW50RXhjZXB0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use channelDescriptor instead')
const Channel$json = {
  '1': 'Channel',
  '2': [
    {'1': 'channelarn', '3': 99276476, '4': 1, '5': 9, '10': 'channelarn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `Channel`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List channelDescriptor = $convert.base64Decode(
    'CgdDaGFubmVsEiEKCmNoYW5uZWxhcm4YvK2rLyABKAlSCmNoYW5uZWxhcm4SFQoEbmFtZRiH5o'
    'F/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use channelARNInvalidExceptionDescriptor instead')
const ChannelARNInvalidException$json = {
  '1': 'ChannelARNInvalidException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ChannelARNInvalidException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List channelARNInvalidExceptionDescriptor =
    $convert.base64Decode(
        'ChpDaGFubmVsQVJOSW52YWxpZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use channelAlreadyExistsExceptionDescriptor instead')
const ChannelAlreadyExistsException$json = {
  '1': 'ChannelAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ChannelAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List channelAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'Ch1DaGFubmVsQWxyZWFkeUV4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use channelExistsForEDSExceptionDescriptor instead')
const ChannelExistsForEDSException$json = {
  '1': 'ChannelExistsForEDSException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ChannelExistsForEDSException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List channelExistsForEDSExceptionDescriptor =
    $convert.base64Decode(
        'ChxDaGFubmVsRXhpc3RzRm9yRURTRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use channelMaxLimitExceededExceptionDescriptor instead')
const ChannelMaxLimitExceededException$json = {
  '1': 'ChannelMaxLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ChannelMaxLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List channelMaxLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiBDaGFubmVsTWF4TGltaXRFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use channelNotFoundExceptionDescriptor instead')
const ChannelNotFoundException$json = {
  '1': 'ChannelNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ChannelNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List channelNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChhDaGFubmVsTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use cloudTrailARNInvalidExceptionDescriptor instead')
const CloudTrailARNInvalidException$json = {
  '1': 'CloudTrailARNInvalidException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudTrailARNInvalidException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudTrailARNInvalidExceptionDescriptor =
    $convert.base64Decode(
        'Ch1DbG91ZFRyYWlsQVJOSW52YWxpZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use cloudTrailAccessNotEnabledExceptionDescriptor instead')
const CloudTrailAccessNotEnabledException$json = {
  '1': 'CloudTrailAccessNotEnabledException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudTrailAccessNotEnabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudTrailAccessNotEnabledExceptionDescriptor =
    $convert.base64Decode(
        'CiNDbG91ZFRyYWlsQWNjZXNzTm90RW5hYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core
    .Deprecated('Use cloudTrailInvalidClientTokenIdExceptionDescriptor instead')
const CloudTrailInvalidClientTokenIdException$json = {
  '1': 'CloudTrailInvalidClientTokenIdException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudTrailInvalidClientTokenIdException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudTrailInvalidClientTokenIdExceptionDescriptor =
    $convert.base64Decode(
        'CidDbG91ZFRyYWlsSW52YWxpZENsaWVudFRva2VuSWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated(
    'Use cloudWatchLogsDeliveryUnavailableExceptionDescriptor instead')
const CloudWatchLogsDeliveryUnavailableException$json = {
  '1': 'CloudWatchLogsDeliveryUnavailableException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudWatchLogsDeliveryUnavailableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    cloudWatchLogsDeliveryUnavailableExceptionDescriptor =
    $convert.base64Decode(
        'CipDbG91ZFdhdGNoTG9nc0RlbGl2ZXJ5VW5hdmFpbGFibGVFeGNlcHRpb24SGwoHbWVzc2FnZR'
        'iFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use concurrentModificationExceptionDescriptor instead')
const ConcurrentModificationException$json = {
  '1': 'ConcurrentModificationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConcurrentModificationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrentModificationExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Db25jdXJyZW50TW9kaWZpY2F0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use conflictExceptionDescriptor instead')
const ConflictException$json = {
  '1': 'ConflictException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictExceptionDescriptor = $convert.base64Decode(
    'ChFDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use contextKeySelectorDescriptor instead')
const ContextKeySelector$json = {
  '1': 'ContextKeySelector',
  '2': [
    {'1': 'equals', '3': 513367477, '4': 3, '5': 9, '10': 'equals'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.Type',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ContextKeySelector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextKeySelectorDescriptor = $convert.base64Decode(
    'ChJDb250ZXh0S2V5U2VsZWN0b3ISGgoGZXF1YWxzGLW75fQBIAMoCVIGZXF1YWxzEigKBHR5cG'
    'UY7qDXigEgASgOMhAuY2xvdWR0cmFpbC5UeXBlUgR0eXBl');

@$core.Deprecated('Use createChannelRequestDescriptor instead')
const CreateChannelRequest$json = {
  '1': 'CreateChannelRequest',
  '2': [
    {
      '1': 'destinations',
      '3': 404385861,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Destination',
      '10': 'destinations'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateChannelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createChannelRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVDaGFubmVsUmVxdWVzdBI/CgxkZXN0aW5hdGlvbnMYxeDpwAEgAygLMhcuY2xvdW'
    'R0cmFpbC5EZXN0aW5hdGlvblIMZGVzdGluYXRpb25zEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUS'
    'GQoGc291cmNlGPnHig8gASgJUgZzb3VyY2USJwoEdGFncxjBwfa1ASADKAsyDy5jbG91ZHRyYW'
    'lsLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createChannelResponseDescriptor instead')
const CreateChannelResponse$json = {
  '1': 'CreateChannelResponse',
  '2': [
    {'1': 'channelarn', '3': 99276476, '4': 1, '5': 9, '10': 'channelarn'},
    {
      '1': 'destinations',
      '3': 404385861,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Destination',
      '10': 'destinations'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateChannelResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createChannelResponseDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVDaGFubmVsUmVzcG9uc2USIQoKY2hhbm5lbGFybhi8rasvIAEoCVIKY2hhbm5lbG'
    'FybhI/CgxkZXN0aW5hdGlvbnMYxeDpwAEgAygLMhcuY2xvdWR0cmFpbC5EZXN0aW5hdGlvblIM'
    'ZGVzdGluYXRpb25zEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGQoGc291cmNlGPnHig8gASgJUg'
    'Zzb3VyY2USJwoEdGFncxjBwfa1ASADKAsyDy5jbG91ZHRyYWlsLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createDashboardRequestDescriptor instead')
const CreateDashboardRequest$json = {
  '1': 'CreateDashboardRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'refreshschedule',
      '3': 261773338,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.RefreshSchedule',
      '10': 'refreshschedule'
    },
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'widgets',
      '3': 501826147,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.RequestWidget',
      '10': 'widgets'
    },
  ],
};

/// Descriptor for `CreateDashboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDashboardRequestDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVEYXNoYm9hcmRSZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSSAoPcmVmcm'
    'VzaHNjaGVkdWxlGJqw6XwgASgLMhsuY2xvdWR0cmFpbC5SZWZyZXNoU2NoZWR1bGVSD3JlZnJl'
    'c2hzY2hlZHVsZRIvCgh0YWdzbGlzdBippJjtASADKAsyDy5jbG91ZHRyYWlsLlRhZ1IIdGFnc2'
    'xpc3QSRgocdGVybWluYXRpb25wcm90ZWN0aW9uZW5hYmxlZBjc89mzASABKAhSHHRlcm1pbmF0'
    'aW9ucHJvdGVjdGlvbmVuYWJsZWQSNwoHd2lkZ2V0cxjjhKXvASADKAsyGS5jbG91ZHRyYWlsLl'
    'JlcXVlc3RXaWRnZXRSB3dpZGdldHM=');

@$core.Deprecated('Use createDashboardResponseDescriptor instead')
const CreateDashboardResponse$json = {
  '1': 'CreateDashboardResponse',
  '2': [
    {'1': 'dashboardarn', '3': 108051951, '4': 1, '5': 9, '10': 'dashboardarn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'refreshschedule',
      '3': 261773338,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.RefreshSchedule',
      '10': 'refreshschedule'
    },
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DashboardType',
      '10': 'type'
    },
    {
      '1': 'widgets',
      '3': 501826147,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Widget',
      '10': 'widgets'
    },
  ],
};

/// Descriptor for `CreateDashboardResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDashboardResponseDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVEYXNoYm9hcmRSZXNwb25zZRIlCgxkYXNoYm9hcmRhcm4Y7/vCMyABKAlSDGRhc2'
    'hib2FyZGFybhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkgKD3JlZnJlc2hzY2hlZHVsZRiasOl8'
    'IAEoCzIbLmNsb3VkdHJhaWwuUmVmcmVzaFNjaGVkdWxlUg9yZWZyZXNoc2NoZWR1bGUSLwoIdG'
    'Fnc2xpc3QYqaSY7QEgAygLMg8uY2xvdWR0cmFpbC5UYWdSCHRhZ3NsaXN0EkYKHHRlcm1pbmF0'
    'aW9ucHJvdGVjdGlvbmVuYWJsZWQY3PPZswEgASgIUhx0ZXJtaW5hdGlvbnByb3RlY3Rpb25lbm'
    'FibGVkEjEKBHR5cGUY7qDXigEgASgOMhkuY2xvdWR0cmFpbC5EYXNoYm9hcmRUeXBlUgR0eXBl'
    'EjAKB3dpZGdldHMY44Sl7wEgAygLMhIuY2xvdWR0cmFpbC5XaWRnZXRSB3dpZGdldHM=');

@$core.Deprecated('Use createEventDataStoreRequestDescriptor instead')
const CreateEventDataStoreRequest$json = {
  '1': 'CreateEventDataStoreRequest',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.BillingMode',
      '10': 'billingmode'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'multiregionenabled',
      '3': 20620620,
      '4': 1,
      '5': 8,
      '10': 'multiregionenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'organizationenabled',
      '3': 480171176,
      '4': 1,
      '5': 8,
      '10': 'organizationenabled'
    },
    {
      '1': 'retentionperiod',
      '3': 196383721,
      '4': 1,
      '5': 5,
      '10': 'retentionperiod'
    },
    {
      '1': 'startingestion',
      '3': 168263832,
      '4': 1,
      '5': 8,
      '10': 'startingestion'
    },
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
  ],
};

/// Descriptor for `CreateEventDataStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEventDataStoreRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVFdmVudERhdGFTdG9yZVJlcXVlc3QSXAoWYWR2YW5jZWRldmVudHNlbGVjdG9ycx'
    'iytsgRIAMoCzIhLmNsb3VkdHJhaWwuQWR2YW5jZWRFdmVudFNlbGVjdG9yUhZhZHZhbmNlZGV2'
    'ZW50c2VsZWN0b3JzEjwKC2JpbGxpbmdtb2RlGMC06FcgASgOMhcuY2xvdWR0cmFpbC5CaWxsaW'
    '5nTW9kZVILYmlsbGluZ21vZGUSHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEjEKEm11'
    'bHRpcmVnaW9uZW5hYmxlZBjMyuoJIAEoCFISbXVsdGlyZWdpb25lbmFibGVkEhUKBG5hbWUYh+'
    'aBfyABKAlSBG5hbWUSNAoTb3JnYW5pemF0aW9uZW5hYmxlZBioqfvkASABKAhSE29yZ2FuaXph'
    'dGlvbmVuYWJsZWQSKwoPcmV0ZW50aW9ucGVyaW9kGOmn0l0gASgFUg9yZXRlbnRpb25wZXJpb2'
    'QSKQoOc3RhcnRpbmdlc3Rpb24YmIGeUCABKAhSDnN0YXJ0aW5nZXN0aW9uEi8KCHRhZ3NsaXN0'
    'GKmkmO0BIAMoCzIPLmNsb3VkdHJhaWwuVGFnUgh0YWdzbGlzdBJGChx0ZXJtaW5hdGlvbnByb3'
    'RlY3Rpb25lbmFibGVkGNzz2bMBIAEoCFIcdGVybWluYXRpb25wcm90ZWN0aW9uZW5hYmxlZA==');

@$core.Deprecated('Use createEventDataStoreResponseDescriptor instead')
const CreateEventDataStoreResponse$json = {
  '1': 'CreateEventDataStoreResponse',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'multiregionenabled',
      '3': 20620620,
      '4': 1,
      '5': 8,
      '10': 'multiregionenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'organizationenabled',
      '3': 480171176,
      '4': 1,
      '5': 8,
      '10': 'organizationenabled'
    },
    {
      '1': 'retentionperiod',
      '3': 196383721,
      '4': 1,
      '5': 5,
      '10': 'retentionperiod'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.EventDataStoreStatus',
      '10': 'status'
    },
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `CreateEventDataStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEventDataStoreResponseDescriptor = $convert.base64Decode(
    'ChxDcmVhdGVFdmVudERhdGFTdG9yZVJlc3BvbnNlElwKFmFkdmFuY2VkZXZlbnRzZWxlY3Rvcn'
    'MYsrbIESADKAsyIS5jbG91ZHRyYWlsLkFkdmFuY2VkRXZlbnRTZWxlY3RvclIWYWR2YW5jZWRl'
    'dmVudHNlbGVjdG9ycxI8CgtiaWxsaW5nbW9kZRjAtOhXIAEoDjIXLmNsb3VkdHJhaWwuQmlsbG'
    'luZ01vZGVSC2JpbGxpbmdtb2RlEi4KEGNyZWF0ZWR0aW1lc3RhbXAY+tvPnwEgASgJUhBjcmVh'
    'dGVkdGltZXN0YW1wEjAKEWV2ZW50ZGF0YXN0b3JlYXJuGOirl54BIAEoCVIRZXZlbnRkYXRhc3'
    'RvcmVhcm4SHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEjEKEm11bHRpcmVnaW9uZW5h'
    'YmxlZBjMyuoJIAEoCFISbXVsdGlyZWdpb25lbmFibGVkEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'USNAoTb3JnYW5pemF0aW9uZW5hYmxlZBioqfvkASABKAhSE29yZ2FuaXphdGlvbmVuYWJsZWQS'
    'KwoPcmV0ZW50aW9ucGVyaW9kGOmn0l0gASgFUg9yZXRlbnRpb25wZXJpb2QSOwoGc3RhdHVzGJ'
    'Dk+wIgASgOMiAuY2xvdWR0cmFpbC5FdmVudERhdGFTdG9yZVN0YXR1c1IGc3RhdHVzEi8KCHRh'
    'Z3NsaXN0GKmkmO0BIAMoCzIPLmNsb3VkdHJhaWwuVGFnUgh0YWdzbGlzdBJGChx0ZXJtaW5hdG'
    'lvbnByb3RlY3Rpb25lbmFibGVkGNzz2bMBIAEoCFIcdGVybWluYXRpb25wcm90ZWN0aW9uZW5h'
    'YmxlZBItChB1cGRhdGVkdGltZXN0YW1wGIHjkxUgASgJUhB1cGRhdGVkdGltZXN0YW1w');

@$core.Deprecated('Use createTrailRequestDescriptor instead')
const CreateTrailRequest$json = {
  '1': 'CreateTrailRequest',
  '2': [
    {
      '1': 'cloudwatchlogsloggrouparn',
      '3': 519613355,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsloggrouparn'
    },
    {
      '1': 'cloudwatchlogsrolearn',
      '3': 55454690,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsrolearn'
    },
    {
      '1': 'enablelogfilevalidation',
      '3': 80236930,
      '4': 1,
      '5': 8,
      '10': 'enablelogfilevalidation'
    },
    {
      '1': 'includeglobalserviceevents',
      '3': 212227451,
      '4': 1,
      '5': 8,
      '10': 'includeglobalserviceevents'
    },
    {
      '1': 'ismultiregiontrail',
      '3': 468658571,
      '4': 1,
      '5': 8,
      '10': 'ismultiregiontrail'
    },
    {
      '1': 'isorganizationtrail',
      '3': 145256127,
      '4': 1,
      '5': 8,
      '10': 'isorganizationtrail'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 's3bucketname', '3': 320495427, '4': 1, '5': 9, '10': 's3bucketname'},
    {'1': 's3keyprefix', '3': 206015359, '4': 1, '5': 9, '10': 's3keyprefix'},
    {'1': 'snstopicname', '3': 415454800, '4': 1, '5': 9, '10': 'snstopicname'},
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
  ],
};

/// Descriptor for `CreateTrailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrailRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVUcmFpbFJlcXVlc3QSQAoZY2xvdWR3YXRjaGxvZ3Nsb2dncm91cGFybhir1+L3AS'
    'ABKAlSGWNsb3Vkd2F0Y2hsb2dzbG9nZ3JvdXBhcm4SNwoVY2xvdWR3YXRjaGxvZ3Nyb2xlYXJu'
    'GOLXuBogASgJUhVjbG91ZHdhdGNobG9nc3JvbGVhcm4SOwoXZW5hYmxlbG9nZmlsZXZhbGlkYX'
    'Rpb24YgqOhJiABKAhSF2VuYWJsZWxvZ2ZpbGV2YWxpZGF0aW9uEkEKGmluY2x1ZGVnbG9iYWxz'
    'ZXJ2aWNlZXZlbnRzGPuqmWUgASgIUhppbmNsdWRlZ2xvYmFsc2VydmljZWV2ZW50cxIyChJpc2'
    '11bHRpcmVnaW9udHJhaWwYi9O83wEgASgIUhJpc211bHRpcmVnaW9udHJhaWwSMwoTaXNvcmdh'
    'bml6YXRpb250cmFpbBi/3aFFIAEoCFITaXNvcmdhbml6YXRpb250cmFpbBIdCghrbXNrZXlpZB'
    'iNyZcWIAEoCVIIa21za2V5aWQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRImCgxzM2J1Y2tldG5h'
    'bWUYw77pmAEgASgJUgxzM2J1Y2tldG5hbWUSIwoLczNrZXlwcmVmaXgY/5aeYiABKAlSC3Mza2'
    'V5cHJlZml4EiYKDHNuc3RvcGljbmFtZRjQrI3GASABKAlSDHNuc3RvcGljbmFtZRIvCgh0YWdz'
    'bGlzdBippJjtASADKAsyDy5jbG91ZHRyYWlsLlRhZ1IIdGFnc2xpc3Q=');

@$core.Deprecated('Use createTrailResponseDescriptor instead')
const CreateTrailResponse$json = {
  '1': 'CreateTrailResponse',
  '2': [
    {
      '1': 'cloudwatchlogsloggrouparn',
      '3': 519613355,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsloggrouparn'
    },
    {
      '1': 'cloudwatchlogsrolearn',
      '3': 55454690,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsrolearn'
    },
    {
      '1': 'includeglobalserviceevents',
      '3': 212227451,
      '4': 1,
      '5': 8,
      '10': 'includeglobalserviceevents'
    },
    {
      '1': 'ismultiregiontrail',
      '3': 468658571,
      '4': 1,
      '5': 8,
      '10': 'ismultiregiontrail'
    },
    {
      '1': 'isorganizationtrail',
      '3': 145256127,
      '4': 1,
      '5': 8,
      '10': 'isorganizationtrail'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'logfilevalidationenabled',
      '3': 35904346,
      '4': 1,
      '5': 8,
      '10': 'logfilevalidationenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 's3bucketname', '3': 320495427, '4': 1, '5': 9, '10': 's3bucketname'},
    {'1': 's3keyprefix', '3': 206015359, '4': 1, '5': 9, '10': 's3keyprefix'},
    {'1': 'snstopicarn', '3': 380025580, '4': 1, '5': 9, '10': 'snstopicarn'},
    {'1': 'snstopicname', '3': 415454800, '4': 1, '5': 9, '10': 'snstopicname'},
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `CreateTrailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTrailResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVUcmFpbFJlc3BvbnNlEkAKGWNsb3Vkd2F0Y2hsb2dzbG9nZ3JvdXBhcm4Yq9fi9w'
    'EgASgJUhljbG91ZHdhdGNobG9nc2xvZ2dyb3VwYXJuEjcKFWNsb3Vkd2F0Y2hsb2dzcm9sZWFy'
    'bhji17gaIAEoCVIVY2xvdWR3YXRjaGxvZ3Nyb2xlYXJuEkEKGmluY2x1ZGVnbG9iYWxzZXJ2aW'
    'NlZXZlbnRzGPuqmWUgASgIUhppbmNsdWRlZ2xvYmFsc2VydmljZWV2ZW50cxIyChJpc211bHRp'
    'cmVnaW9udHJhaWwYi9O83wEgASgIUhJpc211bHRpcmVnaW9udHJhaWwSMwoTaXNvcmdhbml6YX'
    'Rpb250cmFpbBi/3aFFIAEoCFITaXNvcmdhbml6YXRpb250cmFpbBIdCghrbXNrZXlpZBiNyZcW'
    'IAEoCVIIa21za2V5aWQSPQoYbG9nZmlsZXZhbGlkYXRpb25lbmFibGVkGNq2jxEgASgIUhhsb2'
    'dmaWxldmFsaWRhdGlvbmVuYWJsZWQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRImCgxzM2J1Y2tl'
    'dG5hbWUYw77pmAEgASgJUgxzM2J1Y2tldG5hbWUSIwoLczNrZXlwcmVmaXgY/5aeYiABKAlSC3'
    'Mza2V5cHJlZml4EiQKC3Nuc3RvcGljYXJuGOz1mrUBIAEoCVILc25zdG9waWNhcm4SJgoMc25z'
    'dG9waWNuYW1lGNCsjcYBIAEoCVIMc25zdG9waWNuYW1lEh0KCHRyYWlsYXJuGPeK8BIgASgJUg'
    'h0cmFpbGFybg==');

@$core.Deprecated('Use dashboardDetailDescriptor instead')
const DashboardDetail$json = {
  '1': 'DashboardDetail',
  '2': [
    {'1': 'dashboardarn', '3': 108051951, '4': 1, '5': 9, '10': 'dashboardarn'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DashboardType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `DashboardDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dashboardDetailDescriptor = $convert.base64Decode(
    'Cg9EYXNoYm9hcmREZXRhaWwSJQoMZGFzaGJvYXJkYXJuGO/7wjMgASgJUgxkYXNoYm9hcmRhcm'
    '4SMQoEdHlwZRjuoNeKASABKA4yGS5jbG91ZHRyYWlsLkRhc2hib2FyZFR5cGVSBHR5cGU=');

@$core.Deprecated('Use dataResourceDescriptor instead')
const DataResource$json = {
  '1': 'DataResource',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `DataResource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataResourceDescriptor = $convert.base64Decode(
    'CgxEYXRhUmVzb3VyY2USFgoEdHlwZRjuoNeKASABKAlSBHR5cGUSGQoGdmFsdWVzGNzEtGogAy'
    'gJUgZ2YWx1ZXM=');

@$core.Deprecated(
    'Use delegatedAdminAccountLimitExceededExceptionDescriptor instead')
const DelegatedAdminAccountLimitExceededException$json = {
  '1': 'DelegatedAdminAccountLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DelegatedAdminAccountLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    delegatedAdminAccountLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CitEZWxlZ2F0ZWRBZG1pbkFjY291bnRMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2'
        'UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use deleteChannelRequestDescriptor instead')
const DeleteChannelRequest$json = {
  '1': 'DeleteChannelRequest',
  '2': [
    {'1': 'channel', '3': 90016325, '4': 1, '5': 9, '10': 'channel'},
  ],
};

/// Descriptor for `DeleteChannelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteChannelRequestDescriptor =
    $convert.base64Decode(
        'ChREZWxldGVDaGFubmVsUmVxdWVzdBIbCgdjaGFubmVsGMWU9iogASgJUgdjaGFubmVs');

@$core.Deprecated('Use deleteChannelResponseDescriptor instead')
const DeleteChannelResponse$json = {
  '1': 'DeleteChannelResponse',
};

/// Descriptor for `DeleteChannelResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteChannelResponseDescriptor =
    $convert.base64Decode('ChVEZWxldGVDaGFubmVsUmVzcG9uc2U=');

@$core.Deprecated('Use deleteDashboardRequestDescriptor instead')
const DeleteDashboardRequest$json = {
  '1': 'DeleteDashboardRequest',
  '2': [
    {'1': 'dashboardid', '3': 430615703, '4': 1, '5': 9, '10': 'dashboardid'},
  ],
};

/// Descriptor for `DeleteDashboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDashboardRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVEYXNoYm9hcmRSZXF1ZXN0EiQKC2Rhc2hib2FyZGlkGJfZqs0BIAEoCVILZGFzaG'
        'JvYXJkaWQ=');

@$core.Deprecated('Use deleteDashboardResponseDescriptor instead')
const DeleteDashboardResponse$json = {
  '1': 'DeleteDashboardResponse',
};

/// Descriptor for `DeleteDashboardResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDashboardResponseDescriptor =
    $convert.base64Decode('ChdEZWxldGVEYXNoYm9hcmRSZXNwb25zZQ==');

@$core.Deprecated('Use deleteEventDataStoreRequestDescriptor instead')
const DeleteEventDataStoreRequest$json = {
  '1': 'DeleteEventDataStoreRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
  ],
};

/// Descriptor for `DeleteEventDataStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEventDataStoreRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVFdmVudERhdGFTdG9yZVJlcXVlc3QSKQoOZXZlbnRkYXRhc3RvcmUYwdudQSABKA'
        'lSDmV2ZW50ZGF0YXN0b3Jl');

@$core.Deprecated('Use deleteEventDataStoreResponseDescriptor instead')
const DeleteEventDataStoreResponse$json = {
  '1': 'DeleteEventDataStoreResponse',
};

/// Descriptor for `DeleteEventDataStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEventDataStoreResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVFdmVudERhdGFTdG9yZVJlc3BvbnNl');

@$core.Deprecated('Use deleteResourcePolicyRequestDescriptor instead')
const DeleteResourcePolicyRequest$json = {
  '1': 'DeleteResourcePolicyRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `DeleteResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVSZXNvdXJjZVBvbGljeVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUg'
        'tyZXNvdXJjZWFybg==');

@$core.Deprecated('Use deleteResourcePolicyResponseDescriptor instead')
const DeleteResourcePolicyResponse$json = {
  '1': 'DeleteResourcePolicyResponse',
};

/// Descriptor for `DeleteResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVSZXNvdXJjZVBvbGljeVJlc3BvbnNl');

@$core.Deprecated('Use deleteTrailRequestDescriptor instead')
const DeleteTrailRequest$json = {
  '1': 'DeleteTrailRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteTrailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTrailRequestDescriptor =
    $convert.base64Decode(
        'ChJEZWxldGVUcmFpbFJlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use deleteTrailResponseDescriptor instead')
const DeleteTrailResponse$json = {
  '1': 'DeleteTrailResponse',
};

/// Descriptor for `DeleteTrailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTrailResponseDescriptor =
    $convert.base64Decode('ChNEZWxldGVUcmFpbFJlc3BvbnNl');

@$core.Deprecated(
    'Use deregisterOrganizationDelegatedAdminRequestDescriptor instead')
const DeregisterOrganizationDelegatedAdminRequest$json = {
  '1': 'DeregisterOrganizationDelegatedAdminRequest',
  '2': [
    {
      '1': 'delegatedadminaccountid',
      '3': 104507462,
      '4': 1,
      '5': 9,
      '10': 'delegatedadminaccountid'
    },
  ],
};

/// Descriptor for `DeregisterOrganizationDelegatedAdminRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterOrganizationDelegatedAdminRequestDescriptor =
    $convert.base64Decode(
        'CitEZXJlZ2lzdGVyT3JnYW5pemF0aW9uRGVsZWdhdGVkQWRtaW5SZXF1ZXN0EjsKF2RlbGVnYX'
        'RlZGFkbWluYWNjb3VudGlkGMbQ6jEgASgJUhdkZWxlZ2F0ZWRhZG1pbmFjY291bnRpZA==');

@$core.Deprecated(
    'Use deregisterOrganizationDelegatedAdminResponseDescriptor instead')
const DeregisterOrganizationDelegatedAdminResponse$json = {
  '1': 'DeregisterOrganizationDelegatedAdminResponse',
};

/// Descriptor for `DeregisterOrganizationDelegatedAdminResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterOrganizationDelegatedAdminResponseDescriptor =
    $convert.base64Decode(
        'CixEZXJlZ2lzdGVyT3JnYW5pemF0aW9uRGVsZWdhdGVkQWRtaW5SZXNwb25zZQ==');

@$core.Deprecated('Use describeQueryRequestDescriptor instead')
const DescribeQueryRequest$json = {
  '1': 'DescribeQueryRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {'1': 'queryalias', '3': 512762270, '4': 1, '5': 9, '10': 'queryalias'},
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
    {'1': 'refreshid', '3': 57407610, '4': 1, '5': 9, '10': 'refreshid'},
  ],
};

/// Descriptor for `DescribeQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeQueryRequestDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZVF1ZXJ5UmVxdWVzdBIpCg5ldmVudGRhdGFzdG9yZRjB251BIAEoCVIOZXZlbn'
    'RkYXRhc3RvcmUSRgocZXZlbnRkYXRhc3RvcmVvd25lcmFjY291bnRpZBiw3fDgASABKAlSHGV2'
    'ZW50ZGF0YXN0b3Jlb3duZXJhY2NvdW50aWQSIgoKcXVlcnlhbGlhcxiew8D0ASABKAlSCnF1ZX'
    'J5YWxpYXMSGwoHcXVlcnlpZBjv8OY0IAEoCVIHcXVlcnlpZBIfCglyZWZyZXNoaWQY+vCvGyAB'
    'KAlSCXJlZnJlc2hpZA==');

@$core.Deprecated('Use describeQueryResponseDescriptor instead')
const DescribeQueryResponse$json = {
  '1': 'DescribeQueryResponse',
  '2': [
    {
      '1': 'deliverys3uri',
      '3': 230884460,
      '4': 1,
      '5': 9,
      '10': 'deliverys3uri'
    },
    {
      '1': 'deliverystatus',
      '3': 483265504,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DeliveryStatus',
      '10': 'deliverystatus'
    },
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {'1': 'prompt', '3': 263974748, '4': 1, '5': 9, '10': 'prompt'},
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
    {
      '1': 'querystatistics',
      '3': 260794841,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.QueryStatisticsForDescribeQuery',
      '10': 'querystatistics'
    },
    {
      '1': 'querystatus',
      '3': 367016406,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.QueryStatus',
      '10': 'querystatus'
    },
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
  ],
};

/// Descriptor for `DescribeQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeQueryResponseDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZVF1ZXJ5UmVzcG9uc2USJwoNZGVsaXZlcnlzM3VyaRjsiIxuIAEoCVINZGVsaX'
    'ZlcnlzM3VyaRJGCg5kZWxpdmVyeXN0YXR1cxjgl7jmASABKA4yGi5jbG91ZHRyYWlsLkRlbGl2'
    'ZXJ5U3RhdHVzUg5kZWxpdmVyeXN0YXR1cxImCgxlcnJvcm1lc3NhZ2UYqYqr9wEgASgJUgxlcn'
    'Jvcm1lc3NhZ2USRgocZXZlbnRkYXRhc3RvcmVvd25lcmFjY291bnRpZBiw3fDgASABKAlSHGV2'
    'ZW50ZGF0YXN0b3Jlb3duZXJhY2NvdW50aWQSGQoGcHJvbXB0GNze730gASgJUgZwcm9tcHQSGw'
    'oHcXVlcnlpZBjv8OY0IAEoCVIHcXVlcnlpZBJYCg9xdWVyeXN0YXRpc3RpY3MY2dOtfCABKAsy'
    'Ky5jbG91ZHRyYWlsLlF1ZXJ5U3RhdGlzdGljc0ZvckRlc2NyaWJlUXVlcnlSD3F1ZXJ5c3RhdG'
    'lzdGljcxI9CgtxdWVyeXN0YXR1cxjW84CvASABKA4yFy5jbG91ZHRyYWlsLlF1ZXJ5U3RhdHVz'
    'UgtxdWVyeXN0YXR1cxIkCgtxdWVyeXN0cmluZxjnyu/PASABKAlSC3F1ZXJ5c3RyaW5n');

@$core.Deprecated('Use describeTrailsRequestDescriptor instead')
const DescribeTrailsRequest$json = {
  '1': 'DescribeTrailsRequest',
  '2': [
    {
      '1': 'includeshadowtrails',
      '3': 200547819,
      '4': 1,
      '5': 8,
      '10': 'includeshadowtrails'
    },
    {
      '1': 'trailnamelist',
      '3': 529562769,
      '4': 3,
      '5': 9,
      '10': 'trailnamelist'
    },
  ],
};

/// Descriptor for `DescribeTrailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTrailsRequestDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZVRyYWlsc1JlcXVlc3QSMwoTaW5jbHVkZXNoYWRvd3RyYWlscxjru9BfIAEoCF'
    'ITaW5jbHVkZXNoYWRvd3RyYWlscxIoCg10cmFpbG5hbWVsaXN0GJH5wfwBIAMoCVINdHJhaWxu'
    'YW1lbGlzdA==');

@$core.Deprecated('Use describeTrailsResponseDescriptor instead')
const DescribeTrailsResponse$json = {
  '1': 'DescribeTrailsResponse',
  '2': [
    {
      '1': 'traillist',
      '3': 508517484,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Trail',
      '10': 'traillist'
    },
  ],
};

/// Descriptor for `DescribeTrailsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTrailsResponseDescriptor =
    $convert.base64Decode(
        'ChZEZXNjcmliZVRyYWlsc1Jlc3BvbnNlEjMKCXRyYWlsbGlzdBjsuL3yASADKAsyES5jbG91ZH'
        'RyYWlsLlRyYWlsUgl0cmFpbGxpc3Q=');

@$core.Deprecated('Use destinationDescriptor instead')
const Destination$json = {
  '1': 'Destination',
  '2': [
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DestinationType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `Destination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List destinationDescriptor = $convert.base64Decode(
    'CgtEZXN0aW5hdGlvbhIeCghsb2NhdGlvbhjHm4LeASABKAlSCGxvY2F0aW9uEjMKBHR5cGUY7q'
    'DXigEgASgOMhsuY2xvdWR0cmFpbC5EZXN0aW5hdGlvblR5cGVSBHR5cGU=');

@$core.Deprecated('Use disableFederationRequestDescriptor instead')
const DisableFederationRequest$json = {
  '1': 'DisableFederationRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
  ],
};

/// Descriptor for `DisableFederationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableFederationRequestDescriptor =
    $convert.base64Decode(
        'ChhEaXNhYmxlRmVkZXJhdGlvblJlcXVlc3QSKQoOZXZlbnRkYXRhc3RvcmUYwdudQSABKAlSDm'
        'V2ZW50ZGF0YXN0b3Jl');

@$core.Deprecated('Use disableFederationResponseDescriptor instead')
const DisableFederationResponse$json = {
  '1': 'DisableFederationResponse',
  '2': [
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'federationstatus',
      '3': 146235383,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.FederationStatus',
      '10': 'federationstatus'
    },
  ],
};

/// Descriptor for `DisableFederationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableFederationResponseDescriptor = $convert.base64Decode(
    'ChlEaXNhYmxlRmVkZXJhdGlvblJlc3BvbnNlEjAKEWV2ZW50ZGF0YXN0b3JlYXJuGOirl54BIA'
    'EoCVIRZXZlbnRkYXRhc3RvcmVhcm4SSwoQZmVkZXJhdGlvbnN0YXR1cxj3v91FIAEoDjIcLmNs'
    'b3VkdHJhaWwuRmVkZXJhdGlvblN0YXR1c1IQZmVkZXJhdGlvbnN0YXR1cw==');

@$core.Deprecated('Use enableFederationRequestDescriptor instead')
const EnableFederationRequest$json = {
  '1': 'EnableFederationRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {
      '1': 'federationrolearn',
      '3': 504464364,
      '4': 1,
      '5': 9,
      '10': 'federationrolearn'
    },
  ],
};

/// Descriptor for `EnableFederationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableFederationRequestDescriptor = $convert.base64Decode(
    'ChdFbmFibGVGZWRlcmF0aW9uUmVxdWVzdBIpCg5ldmVudGRhdGFzdG9yZRjB251BIAEoCVIOZX'
    'ZlbnRkYXRhc3RvcmUSMAoRZmVkZXJhdGlvbnJvbGVhcm4Y7IfG8AEgASgJUhFmZWRlcmF0aW9u'
    'cm9sZWFybg==');

@$core.Deprecated('Use enableFederationResponseDescriptor instead')
const EnableFederationResponse$json = {
  '1': 'EnableFederationResponse',
  '2': [
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'federationrolearn',
      '3': 504464364,
      '4': 1,
      '5': 9,
      '10': 'federationrolearn'
    },
    {
      '1': 'federationstatus',
      '3': 146235383,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.FederationStatus',
      '10': 'federationstatus'
    },
  ],
};

/// Descriptor for `EnableFederationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableFederationResponseDescriptor = $convert.base64Decode(
    'ChhFbmFibGVGZWRlcmF0aW9uUmVzcG9uc2USMAoRZXZlbnRkYXRhc3RvcmVhcm4Y6KuXngEgAS'
    'gJUhFldmVudGRhdGFzdG9yZWFybhIwChFmZWRlcmF0aW9ucm9sZWFybhjsh8bwASABKAlSEWZl'
    'ZGVyYXRpb25yb2xlYXJuEksKEGZlZGVyYXRpb25zdGF0dXMY97/dRSABKA4yHC5jbG91ZHRyYW'
    'lsLkZlZGVyYXRpb25TdGF0dXNSEGZlZGVyYXRpb25zdGF0dXM=');

@$core.Deprecated('Use eventDescriptor instead')
const Event$json = {
  '1': 'Event',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {
      '1': 'cloudtrailevent',
      '3': 344297255,
      '4': 1,
      '5': 9,
      '10': 'cloudtrailevent'
    },
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 9, '10': 'eventid'},
    {'1': 'eventname', '3': 264746781, '4': 1, '5': 9, '10': 'eventname'},
    {'1': 'eventsource', '3': 37841339, '4': 1, '5': 9, '10': 'eventsource'},
    {'1': 'eventtime', '3': 222361647, '4': 1, '5': 9, '10': 'eventtime'},
    {'1': 'readonly', '3': 129039352, '4': 1, '5': 9, '10': 'readonly'},
    {
      '1': 'resources',
      '3': 358918291,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Resource',
      '10': 'resources'
    },
    {'1': 'username', '3': 470340826, '4': 1, '5': 9, '10': 'username'},
  ],
};

/// Descriptor for `Event`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDescriptor = $convert.base64Decode(
    'CgVFdmVudBIkCgthY2Nlc3NrZXlpZBigt7fYASABKAlSC2FjY2Vzc2tleWlkEiwKD2Nsb3VkdH'
    'JhaWxldmVudBinnpakASABKAlSD2Nsb3VkdHJhaWxldmVudBIcCgdldmVudGlkGNOW3bMBIAEo'
    'CVIHZXZlbnRpZBIfCglldmVudG5hbWUYne6efiABKAlSCWV2ZW50bmFtZRIjCgtldmVudHNvdX'
    'JjZRi704USIAEoCVILZXZlbnRzb3VyY2USHwoJZXZlbnR0aW1lGK/wg2ogASgJUglldmVudHRp'
    'bWUSHQoIcmVhZG9ubHkY+PfDPSABKAlSCHJlYWRvbmx5EjYKCXJlc291cmNlcxiT0ZKrASADKA'
    'syFC5jbG91ZHRyYWlsLlJlc291cmNlUglyZXNvdXJjZXMSHgoIdXNlcm5hbWUY2qmj4AEgASgJ'
    'Ugh1c2VybmFtZQ==');

@$core.Deprecated('Use eventDataStoreDescriptor instead')
const EventDataStore$json = {
  '1': 'EventDataStore',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'multiregionenabled',
      '3': 20620620,
      '4': 1,
      '5': 8,
      '10': 'multiregionenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'organizationenabled',
      '3': 480171176,
      '4': 1,
      '5': 8,
      '10': 'organizationenabled'
    },
    {
      '1': 'retentionperiod',
      '3': 196383721,
      '4': 1,
      '5': 5,
      '10': 'retentionperiod'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.EventDataStoreStatus',
      '10': 'status'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `EventDataStore`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDataStoreDescriptor = $convert.base64Decode(
    'Cg5FdmVudERhdGFTdG9yZRJcChZhZHZhbmNlZGV2ZW50c2VsZWN0b3JzGLK2yBEgAygLMiEuY2'
    'xvdWR0cmFpbC5BZHZhbmNlZEV2ZW50U2VsZWN0b3JSFmFkdmFuY2VkZXZlbnRzZWxlY3RvcnMS'
    'LgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEGNyZWF0ZWR0aW1lc3RhbXASMAoRZXZlbn'
    'RkYXRhc3RvcmVhcm4Y6KuXngEgASgJUhFldmVudGRhdGFzdG9yZWFybhIxChJtdWx0aXJlZ2lv'
    'bmVuYWJsZWQYzMrqCSABKAhSEm11bHRpcmVnaW9uZW5hYmxlZBIVCgRuYW1lGIfmgX8gASgJUg'
    'RuYW1lEjQKE29yZ2FuaXphdGlvbmVuYWJsZWQYqKn75AEgASgIUhNvcmdhbml6YXRpb25lbmFi'
    'bGVkEisKD3JldGVudGlvbnBlcmlvZBjpp9JdIAEoBVIPcmV0ZW50aW9ucGVyaW9kEjsKBnN0YX'
    'R1cxiQ5PsCIAEoDjIgLmNsb3VkdHJhaWwuRXZlbnREYXRhU3RvcmVTdGF0dXNSBnN0YXR1cxJG'
    'Chx0ZXJtaW5hdGlvbnByb3RlY3Rpb25lbmFibGVkGNzz2bMBIAEoCFIcdGVybWluYXRpb25wcm'
    '90ZWN0aW9uZW5hYmxlZBItChB1cGRhdGVkdGltZXN0YW1wGIHjkxUgASgJUhB1cGRhdGVkdGlt'
    'ZXN0YW1w');

@$core.Deprecated('Use eventDataStoreARNInvalidExceptionDescriptor instead')
const EventDataStoreARNInvalidException$json = {
  '1': 'EventDataStoreARNInvalidException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EventDataStoreARNInvalidException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDataStoreARNInvalidExceptionDescriptor =
    $convert.base64Decode(
        'CiFFdmVudERhdGFTdG9yZUFSTkludmFsaWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use eventDataStoreAlreadyExistsExceptionDescriptor instead')
const EventDataStoreAlreadyExistsException$json = {
  '1': 'EventDataStoreAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EventDataStoreAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDataStoreAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'CiRFdmVudERhdGFTdG9yZUFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated(
    'Use eventDataStoreFederationEnabledExceptionDescriptor instead')
const EventDataStoreFederationEnabledException$json = {
  '1': 'EventDataStoreFederationEnabledException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EventDataStoreFederationEnabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDataStoreFederationEnabledExceptionDescriptor =
    $convert.base64Decode(
        'CihFdmVudERhdGFTdG9yZUZlZGVyYXRpb25FbmFibGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core
    .Deprecated('Use eventDataStoreHasOngoingImportExceptionDescriptor instead')
const EventDataStoreHasOngoingImportException$json = {
  '1': 'EventDataStoreHasOngoingImportException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EventDataStoreHasOngoingImportException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDataStoreHasOngoingImportExceptionDescriptor =
    $convert.base64Decode(
        'CidFdmVudERhdGFTdG9yZUhhc09uZ29pbmdJbXBvcnRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core
    .Deprecated('Use eventDataStoreMaxLimitExceededExceptionDescriptor instead')
const EventDataStoreMaxLimitExceededException$json = {
  '1': 'EventDataStoreMaxLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EventDataStoreMaxLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDataStoreMaxLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CidFdmVudERhdGFTdG9yZU1heExpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use eventDataStoreNotFoundExceptionDescriptor instead')
const EventDataStoreNotFoundException$json = {
  '1': 'EventDataStoreNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EventDataStoreNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDataStoreNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ch9FdmVudERhdGFTdG9yZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated(
    'Use eventDataStoreTerminationProtectedExceptionDescriptor instead')
const EventDataStoreTerminationProtectedException$json = {
  '1': 'EventDataStoreTerminationProtectedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EventDataStoreTerminationProtectedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    eventDataStoreTerminationProtectedExceptionDescriptor =
    $convert.base64Decode(
        'CitFdmVudERhdGFTdG9yZVRlcm1pbmF0aW9uUHJvdGVjdGVkRXhjZXB0aW9uEhsKB21lc3NhZ2'
        'UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use eventSelectorDescriptor instead')
const EventSelector$json = {
  '1': 'EventSelector',
  '2': [
    {
      '1': 'dataresources',
      '3': 126123155,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.DataResource',
      '10': 'dataresources'
    },
    {
      '1': 'excludemanagementeventsources',
      '3': 225676985,
      '4': 3,
      '5': 9,
      '10': 'excludemanagementeventsources'
    },
    {
      '1': 'includemanagementevents',
      '3': 215824550,
      '4': 1,
      '5': 8,
      '10': 'includemanagementevents'
    },
    {
      '1': 'readwritetype',
      '3': 296653333,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ReadWriteType',
      '10': 'readwritetype'
    },
  ],
};

/// Descriptor for `EventSelector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventSelectorDescriptor = $convert.base64Decode(
    'Cg1FdmVudFNlbGVjdG9yEkEKDWRhdGFyZXNvdXJjZXMYk/mRPCADKAsyGC5jbG91ZHRyYWlsLk'
    'RhdGFSZXNvdXJjZVINZGF0YXJlc291cmNlcxJHCh1leGNsdWRlbWFuYWdlbWVudGV2ZW50c291'
    'cmNlcxi5nc5rIAMoCVIdZXhjbHVkZW1hbmFnZW1lbnRldmVudHNvdXJjZXMSOwoXaW5jbHVkZW'
    '1hbmFnZW1lbnRldmVudHMYpvH0ZiABKAhSF2luY2x1ZGVtYW5hZ2VtZW50ZXZlbnRzEkMKDXJl'
    'YWR3cml0ZXR5cGUYlaS6jQEgASgOMhkuY2xvdWR0cmFpbC5SZWFkV3JpdGVUeXBlUg1yZWFkd3'
    'JpdGV0eXBl');

@$core.Deprecated('Use generateQueryRequestDescriptor instead')
const GenerateQueryRequest$json = {
  '1': 'GenerateQueryRequest',
  '2': [
    {
      '1': 'eventdatastores',
      '3': 152154314,
      '4': 3,
      '5': 9,
      '10': 'eventdatastores'
    },
    {'1': 'prompt', '3': 263974748, '4': 1, '5': 9, '10': 'prompt'},
  ],
};

/// Descriptor for `GenerateQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateQueryRequestDescriptor = $convert.base64Decode(
    'ChRHZW5lcmF0ZVF1ZXJ5UmVxdWVzdBIrCg9ldmVudGRhdGFzdG9yZXMYyuHGSCADKAlSD2V2ZW'
    '50ZGF0YXN0b3JlcxIZCgZwcm9tcHQY3N7vfSABKAlSBnByb21wdA==');

@$core.Deprecated('Use generateQueryResponseDescriptor instead')
const GenerateQueryResponse$json = {
  '1': 'GenerateQueryResponse',
  '2': [
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {'1': 'queryalias', '3': 512762270, '4': 1, '5': 9, '10': 'queryalias'},
    {
      '1': 'querystatement',
      '3': 340852217,
      '4': 1,
      '5': 9,
      '10': 'querystatement'
    },
  ],
};

/// Descriptor for `GenerateQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateQueryResponseDescriptor = $convert.base64Decode(
    'ChVHZW5lcmF0ZVF1ZXJ5UmVzcG9uc2USRgocZXZlbnRkYXRhc3RvcmVvd25lcmFjY291bnRpZB'
    'iw3fDgASABKAlSHGV2ZW50ZGF0YXN0b3Jlb3duZXJhY2NvdW50aWQSIgoKcXVlcnlhbGlhcxie'
    'w8D0ASABKAlSCnF1ZXJ5YWxpYXMSKgoOcXVlcnlzdGF0ZW1lbnQY+fvDogEgASgJUg5xdWVyeX'
    'N0YXRlbWVudA==');

@$core.Deprecated('Use generateResponseExceptionDescriptor instead')
const GenerateResponseException$json = {
  '1': 'GenerateResponseException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `GenerateResponseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateResponseExceptionDescriptor =
    $convert.base64Decode(
        'ChlHZW5lcmF0ZVJlc3BvbnNlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use getChannelRequestDescriptor instead')
const GetChannelRequest$json = {
  '1': 'GetChannelRequest',
  '2': [
    {'1': 'channel', '3': 90016325, '4': 1, '5': 9, '10': 'channel'},
  ],
};

/// Descriptor for `GetChannelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChannelRequestDescriptor = $convert.base64Decode(
    'ChFHZXRDaGFubmVsUmVxdWVzdBIbCgdjaGFubmVsGMWU9iogASgJUgdjaGFubmVs');

@$core.Deprecated('Use getChannelResponseDescriptor instead')
const GetChannelResponse$json = {
  '1': 'GetChannelResponse',
  '2': [
    {'1': 'channelarn', '3': 99276476, '4': 1, '5': 9, '10': 'channelarn'},
    {
      '1': 'destinations',
      '3': 404385861,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Destination',
      '10': 'destinations'
    },
    {
      '1': 'ingestionstatus',
      '3': 188310656,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.IngestionStatus',
      '10': 'ingestionstatus'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {
      '1': 'sourceconfig',
      '3': 415255383,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.SourceConfig',
      '10': 'sourceconfig'
    },
  ],
};

/// Descriptor for `GetChannelResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChannelResponseDescriptor = $convert.base64Decode(
    'ChJHZXRDaGFubmVsUmVzcG9uc2USIQoKY2hhbm5lbGFybhi8rasvIAEoCVIKY2hhbm5lbGFybh'
    'I/CgxkZXN0aW5hdGlvbnMYxeDpwAEgAygLMhcuY2xvdWR0cmFpbC5EZXN0aW5hdGlvblIMZGVz'
    'dGluYXRpb25zEkgKD2luZ2VzdGlvbnN0YXR1cxiAyeVZIAEoCzIbLmNsb3VkdHJhaWwuSW5nZX'
    'N0aW9uU3RhdHVzUg9pbmdlc3Rpb25zdGF0dXMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIZCgZz'
    'b3VyY2UY+ceKDyABKAlSBnNvdXJjZRJACgxzb3VyY2Vjb25maWcY15aBxgEgASgLMhguY2xvdW'
    'R0cmFpbC5Tb3VyY2VDb25maWdSDHNvdXJjZWNvbmZpZw==');

@$core.Deprecated('Use getDashboardRequestDescriptor instead')
const GetDashboardRequest$json = {
  '1': 'GetDashboardRequest',
  '2': [
    {'1': 'dashboardid', '3': 430615703, '4': 1, '5': 9, '10': 'dashboardid'},
  ],
};

/// Descriptor for `GetDashboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDashboardRequestDescriptor = $convert.base64Decode(
    'ChNHZXREYXNoYm9hcmRSZXF1ZXN0EiQKC2Rhc2hib2FyZGlkGJfZqs0BIAEoCVILZGFzaGJvYX'
    'JkaWQ=');

@$core.Deprecated('Use getDashboardResponseDescriptor instead')
const GetDashboardResponse$json = {
  '1': 'GetDashboardResponse',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'dashboardarn', '3': 108051951, '4': 1, '5': 9, '10': 'dashboardarn'},
    {
      '1': 'lastrefreshfailurereason',
      '3': 493889247,
      '4': 1,
      '5': 9,
      '10': 'lastrefreshfailurereason'
    },
    {
      '1': 'lastrefreshid',
      '3': 272889110,
      '4': 1,
      '5': 9,
      '10': 'lastrefreshid'
    },
    {
      '1': 'refreshschedule',
      '3': 261773338,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.RefreshSchedule',
      '10': 'refreshschedule'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DashboardStatus',
      '10': 'status'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DashboardType',
      '10': 'type'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
    {
      '1': 'widgets',
      '3': 501826147,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Widget',
      '10': 'widgets'
    },
  ],
};

/// Descriptor for `GetDashboardResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDashboardResponseDescriptor = $convert.base64Decode(
    'ChRHZXREYXNoYm9hcmRSZXNwb25zZRIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCVIQY3'
    'JlYXRlZHRpbWVzdGFtcBIlCgxkYXNoYm9hcmRhcm4Y7/vCMyABKAlSDGRhc2hib2FyZGFybhI+'
    'ChhsYXN0cmVmcmVzaGZhaWx1cmVyZWFzb24Y383A6wEgASgJUhhsYXN0cmVmcmVzaGZhaWx1cm'
    'VyZWFzb24SKAoNbGFzdHJlZnJlc2hpZBiW6o+CASABKAlSDWxhc3RyZWZyZXNoaWQSSAoPcmVm'
    'cmVzaHNjaGVkdWxlGJqw6XwgASgLMhsuY2xvdWR0cmFpbC5SZWZyZXNoU2NoZWR1bGVSD3JlZn'
    'Jlc2hzY2hlZHVsZRI2CgZzdGF0dXMYkOT7AiABKA4yGy5jbG91ZHRyYWlsLkRhc2hib2FyZFN0'
    'YXR1c1IGc3RhdHVzEkYKHHRlcm1pbmF0aW9ucHJvdGVjdGlvbmVuYWJsZWQY3PPZswEgASgIUh'
    'x0ZXJtaW5hdGlvbnByb3RlY3Rpb25lbmFibGVkEjEKBHR5cGUY7qDXigEgASgOMhkuY2xvdWR0'
    'cmFpbC5EYXNoYm9hcmRUeXBlUgR0eXBlEi0KEHVwZGF0ZWR0aW1lc3RhbXAYgeOTFSABKAlSEH'
    'VwZGF0ZWR0aW1lc3RhbXASMAoHd2lkZ2V0cxjjhKXvASADKAsyEi5jbG91ZHRyYWlsLldpZGdl'
    'dFIHd2lkZ2V0cw==');

@$core.Deprecated('Use getEventConfigurationRequestDescriptor instead')
const GetEventConfigurationRequest$json = {
  '1': 'GetEventConfigurationRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {'1': 'trailname', '3': 507774985, '4': 1, '5': 9, '10': 'trailname'},
  ],
};

/// Descriptor for `GetEventConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEventConfigurationRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRFdmVudENvbmZpZ3VyYXRpb25SZXF1ZXN0EikKDmV2ZW50ZGF0YXN0b3JlGMHbnUEgAS'
        'gJUg5ldmVudGRhdGFzdG9yZRIgCgl0cmFpbG5hbWUYiZCQ8gEgASgJUgl0cmFpbG5hbWU=');

@$core.Deprecated('Use getEventConfigurationResponseDescriptor instead')
const GetEventConfigurationResponse$json = {
  '1': 'GetEventConfigurationResponse',
  '2': [
    {
      '1': 'aggregationconfigurations',
      '3': 481530463,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AggregationConfiguration',
      '10': 'aggregationconfigurations'
    },
    {
      '1': 'contextkeyselectors',
      '3': 342040300,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.ContextKeySelector',
      '10': 'contextkeyselectors'
    },
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'maxeventsize',
      '3': 517627763,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.MaxEventSize',
      '10': 'maxeventsize'
    },
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `GetEventConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEventConfigurationResponseDescriptor = $convert.base64Decode(
    'Ch1HZXRFdmVudENvbmZpZ3VyYXRpb25SZXNwb25zZRJmChlhZ2dyZWdhdGlvbmNvbmZpZ3VyYX'
    'Rpb25zGN+kzuUBIAMoCzIkLmNsb3VkdHJhaWwuQWdncmVnYXRpb25Db25maWd1cmF0aW9uUhlh'
    'Z2dyZWdhdGlvbmNvbmZpZ3VyYXRpb25zElQKE2NvbnRleHRrZXlzZWxlY3RvcnMY7L2MowEgAy'
    'gLMh4uY2xvdWR0cmFpbC5Db250ZXh0S2V5U2VsZWN0b3JSE2NvbnRleHRrZXlzZWxlY3RvcnMS'
    'MAoRZXZlbnRkYXRhc3RvcmVhcm4Y6KuXngEgASgJUhFldmVudGRhdGFzdG9yZWFybhJACgxtYX'
    'hldmVudHNpemUY877p9gEgASgOMhguY2xvdWR0cmFpbC5NYXhFdmVudFNpemVSDG1heGV2ZW50'
    'c2l6ZRIdCgh0cmFpbGFybhj3ivASIAEoCVIIdHJhaWxhcm4=');

@$core.Deprecated('Use getEventDataStoreRequestDescriptor instead')
const GetEventDataStoreRequest$json = {
  '1': 'GetEventDataStoreRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
  ],
};

/// Descriptor for `GetEventDataStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEventDataStoreRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRFdmVudERhdGFTdG9yZVJlcXVlc3QSKQoOZXZlbnRkYXRhc3RvcmUYwdudQSABKAlSDm'
        'V2ZW50ZGF0YXN0b3Jl');

@$core.Deprecated('Use getEventDataStoreResponseDescriptor instead')
const GetEventDataStoreResponse$json = {
  '1': 'GetEventDataStoreResponse',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'federationrolearn',
      '3': 504464364,
      '4': 1,
      '5': 9,
      '10': 'federationrolearn'
    },
    {
      '1': 'federationstatus',
      '3': 146235383,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.FederationStatus',
      '10': 'federationstatus'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'multiregionenabled',
      '3': 20620620,
      '4': 1,
      '5': 8,
      '10': 'multiregionenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'organizationenabled',
      '3': 480171176,
      '4': 1,
      '5': 8,
      '10': 'organizationenabled'
    },
    {
      '1': 'partitionkeys',
      '3': 200562986,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.PartitionKey',
      '10': 'partitionkeys'
    },
    {
      '1': 'retentionperiod',
      '3': 196383721,
      '4': 1,
      '5': 5,
      '10': 'retentionperiod'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.EventDataStoreStatus',
      '10': 'status'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `GetEventDataStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEventDataStoreResponseDescriptor = $convert.base64Decode(
    'ChlHZXRFdmVudERhdGFTdG9yZVJlc3BvbnNlElwKFmFkdmFuY2VkZXZlbnRzZWxlY3RvcnMYsr'
    'bIESADKAsyIS5jbG91ZHRyYWlsLkFkdmFuY2VkRXZlbnRTZWxlY3RvclIWYWR2YW5jZWRldmVu'
    'dHNlbGVjdG9ycxI8CgtiaWxsaW5nbW9kZRjAtOhXIAEoDjIXLmNsb3VkdHJhaWwuQmlsbGluZ0'
    '1vZGVSC2JpbGxpbmdtb2RlEi4KEGNyZWF0ZWR0aW1lc3RhbXAY+tvPnwEgASgJUhBjcmVhdGVk'
    'dGltZXN0YW1wEjAKEWV2ZW50ZGF0YXN0b3JlYXJuGOirl54BIAEoCVIRZXZlbnRkYXRhc3Rvcm'
    'Vhcm4SMAoRZmVkZXJhdGlvbnJvbGVhcm4Y7IfG8AEgASgJUhFmZWRlcmF0aW9ucm9sZWFybhJL'
    'ChBmZWRlcmF0aW9uc3RhdHVzGPe/3UUgASgOMhwuY2xvdWR0cmFpbC5GZWRlcmF0aW9uU3RhdH'
    'VzUhBmZWRlcmF0aW9uc3RhdHVzEh0KCGttc2tleWlkGI3JlxYgASgJUghrbXNrZXlpZBIxChJt'
    'dWx0aXJlZ2lvbmVuYWJsZWQYzMrqCSABKAhSEm11bHRpcmVnaW9uZW5hYmxlZBIVCgRuYW1lGI'
    'fmgX8gASgJUgRuYW1lEjQKE29yZ2FuaXphdGlvbmVuYWJsZWQYqKn75AEgASgIUhNvcmdhbml6'
    'YXRpb25lbmFibGVkEkEKDXBhcnRpdGlvbmtleXMYqrLRXyADKAsyGC5jbG91ZHRyYWlsLlBhcn'
    'RpdGlvbktleVINcGFydGl0aW9ua2V5cxIrCg9yZXRlbnRpb25wZXJpb2QY6afSXSABKAVSD3Jl'
    'dGVudGlvbnBlcmlvZBI7CgZzdGF0dXMYkOT7AiABKA4yIC5jbG91ZHRyYWlsLkV2ZW50RGF0YV'
    'N0b3JlU3RhdHVzUgZzdGF0dXMSRgocdGVybWluYXRpb25wcm90ZWN0aW9uZW5hYmxlZBjc89mz'
    'ASABKAhSHHRlcm1pbmF0aW9ucHJvdGVjdGlvbmVuYWJsZWQSLQoQdXBkYXRlZHRpbWVzdGFtcB'
    'iB45MVIAEoCVIQdXBkYXRlZHRpbWVzdGFtcA==');

@$core.Deprecated('Use getEventSelectorsRequestDescriptor instead')
const GetEventSelectorsRequest$json = {
  '1': 'GetEventSelectorsRequest',
  '2': [
    {'1': 'trailname', '3': 507774985, '4': 1, '5': 9, '10': 'trailname'},
  ],
};

/// Descriptor for `GetEventSelectorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEventSelectorsRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRFdmVudFNlbGVjdG9yc1JlcXVlc3QSIAoJdHJhaWxuYW1lGImQkPIBIAEoCVIJdHJhaW'
        'xuYW1l');

@$core.Deprecated('Use getEventSelectorsResponseDescriptor instead')
const GetEventSelectorsResponse$json = {
  '1': 'GetEventSelectorsResponse',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'eventselectors',
      '3': 344575328,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.EventSelector',
      '10': 'eventselectors'
    },
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `GetEventSelectorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEventSelectorsResponseDescriptor = $convert.base64Decode(
    'ChlHZXRFdmVudFNlbGVjdG9yc1Jlc3BvbnNlElwKFmFkdmFuY2VkZXZlbnRzZWxlY3RvcnMYsr'
    'bIESADKAsyIS5jbG91ZHRyYWlsLkFkdmFuY2VkRXZlbnRTZWxlY3RvclIWYWR2YW5jZWRldmVu'
    'dHNlbGVjdG9ycxJFCg5ldmVudHNlbGVjdG9ycxjgmqekASADKAsyGS5jbG91ZHRyYWlsLkV2ZW'
    '50U2VsZWN0b3JSDmV2ZW50c2VsZWN0b3JzEh0KCHRyYWlsYXJuGPeK8BIgASgJUgh0cmFpbGFy'
    'bg==');

@$core.Deprecated('Use getImportRequestDescriptor instead')
const GetImportRequest$json = {
  '1': 'GetImportRequest',
  '2': [
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
  ],
};

/// Descriptor for `GetImportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImportRequestDescriptor = $convert.base64Decode(
    'ChBHZXRJbXBvcnRSZXF1ZXN0Eh4KCGltcG9ydGlkGNqUrMgBIAEoCVIIaW1wb3J0aWQ=');

@$core.Deprecated('Use getImportResponseDescriptor instead')
const GetImportResponse$json = {
  '1': 'GetImportResponse',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'destinations', '3': 404385861, '4': 3, '5': 9, '10': 'destinations'},
    {'1': 'endeventtime', '3': 260441984, '4': 1, '5': 9, '10': 'endeventtime'},
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importsource',
      '3': 41128754,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.ImportSource',
      '10': 'importsource'
    },
    {
      '1': 'importstatistics',
      '3': 48175528,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.ImportStatistics',
      '10': 'importstatistics'
    },
    {
      '1': 'importstatus',
      '3': 129077631,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'starteventtime',
      '3': 107272573,
      '4': 1,
      '5': 9,
      '10': 'starteventtime'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `GetImportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImportResponseDescriptor = $convert.base64Decode(
    'ChFHZXRJbXBvcnRSZXNwb25zZRIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCVIQY3JlYX'
    'RlZHRpbWVzdGFtcBImCgxkZXN0aW5hdGlvbnMYxeDpwAEgAygJUgxkZXN0aW5hdGlvbnMSJQoM'
    'ZW5kZXZlbnR0aW1lGICPmHwgASgJUgxlbmRldmVudHRpbWUSHgoIaW1wb3J0aWQY2pSsyAEgAS'
    'gJUghpbXBvcnRpZBI/CgxpbXBvcnRzb3VyY2UYsqbOEyABKAsyGC5jbG91ZHRyYWlsLkltcG9y'
    'dFNvdXJjZVIMaW1wb3J0c291cmNlEksKEGltcG9ydHN0YXRpc3RpY3MYqLP8FiABKAsyHC5jbG'
    '91ZHRyYWlsLkltcG9ydFN0YXRpc3RpY3NSEGltcG9ydHN0YXRpc3RpY3MSPwoMaW1wb3J0c3Rh'
    'dHVzGP+ixj0gASgOMhguY2xvdWR0cmFpbC5JbXBvcnRTdGF0dXNSDGltcG9ydHN0YXR1cxIpCg'
    '5zdGFydGV2ZW50dGltZRj9spMzIAEoCVIOc3RhcnRldmVudHRpbWUSLQoQdXBkYXRlZHRpbWVz'
    'dGFtcBiB45MVIAEoCVIQdXBkYXRlZHRpbWVzdGFtcA==');

@$core.Deprecated('Use getInsightSelectorsRequestDescriptor instead')
const GetInsightSelectorsRequest$json = {
  '1': 'GetInsightSelectorsRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {'1': 'trailname', '3': 507774985, '4': 1, '5': 9, '10': 'trailname'},
  ],
};

/// Descriptor for `GetInsightSelectorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInsightSelectorsRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRJbnNpZ2h0U2VsZWN0b3JzUmVxdWVzdBIpCg5ldmVudGRhdGFzdG9yZRjB251BIAEoCV'
        'IOZXZlbnRkYXRhc3RvcmUSIAoJdHJhaWxuYW1lGImQkPIBIAEoCVIJdHJhaWxuYW1l');

@$core.Deprecated('Use getInsightSelectorsResponseDescriptor instead')
const GetInsightSelectorsResponse$json = {
  '1': 'GetInsightSelectorsResponse',
  '2': [
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'insightselectors',
      '3': 49628748,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.InsightSelector',
      '10': 'insightselectors'
    },
    {
      '1': 'insightsdestination',
      '3': 396776649,
      '4': 1,
      '5': 9,
      '10': 'insightsdestination'
    },
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `GetInsightSelectorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInsightSelectorsResponseDescriptor = $convert.base64Decode(
    'ChtHZXRJbnNpZ2h0U2VsZWN0b3JzUmVzcG9uc2USMAoRZXZlbnRkYXRhc3RvcmVhcm4Y6KuXng'
    'EgASgJUhFldmVudGRhdGFzdG9yZWFybhJKChBpbnNpZ2h0c2VsZWN0b3JzGMyM1RcgAygLMhsu'
    'Y2xvdWR0cmFpbC5JbnNpZ2h0U2VsZWN0b3JSEGluc2lnaHRzZWxlY3RvcnMSNAoTaW5zaWdodH'
    'NkZXN0aW5hdGlvbhjJqZm9ASABKAlSE2luc2lnaHRzZGVzdGluYXRpb24SHQoIdHJhaWxhcm4Y'
    '94rwEiABKAlSCHRyYWlsYXJu');

@$core.Deprecated('Use getQueryResultsRequestDescriptor instead')
const GetQueryResultsRequest$json = {
  '1': 'GetQueryResultsRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {
      '1': 'maxqueryresults',
      '3': 120660448,
      '4': 1,
      '5': 5,
      '10': 'maxqueryresults'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `GetQueryResultsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryResultsRequestDescriptor = $convert.base64Decode(
    'ChZHZXRRdWVyeVJlc3VsdHNSZXF1ZXN0EikKDmV2ZW50ZGF0YXN0b3JlGMHbnUEgASgJUg5ldm'
    'VudGRhdGFzdG9yZRJGChxldmVudGRhdGFzdG9yZW93bmVyYWNjb3VudGlkGLDd8OABIAEoCVIc'
    'ZXZlbnRkYXRhc3RvcmVvd25lcmFjY291bnRpZBIrCg9tYXhxdWVyeXJlc3VsdHMY4MPEOSABKA'
    'VSD21heHF1ZXJ5cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIbCgdx'
    'dWVyeWlkGO/w5jQgASgJUgdxdWVyeWlk');

@$core.Deprecated('Use getQueryResultsResponseDescriptor instead')
const GetQueryResultsResponse$json = {
  '1': 'GetQueryResultsResponse',
  '2': [
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queryresultrows',
      '3': 240264704,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.GetQueryResultsResponse.QueryresultrowsEntry',
      '10': 'queryresultrows'
    },
    {
      '1': 'querystatistics',
      '3': 260794841,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.QueryStatistics',
      '10': 'querystatistics'
    },
    {
      '1': 'querystatus',
      '3': 367016406,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.QueryStatus',
      '10': 'querystatus'
    },
  ],
  '3': [GetQueryResultsResponse_QueryresultrowsEntry$json],
};

@$core.Deprecated('Use getQueryResultsResponseDescriptor instead')
const GetQueryResultsResponse_QueryresultrowsEntry$json = {
  '1': 'QueryresultrowsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetQueryResultsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryResultsResponseDescriptor = $convert.base64Decode(
    'ChdHZXRRdWVyeVJlc3VsdHNSZXNwb25zZRImCgxlcnJvcm1lc3NhZ2UYqYqr9wEgASgJUgxlcn'
    'Jvcm1lc3NhZ2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SZQoPcXVlcnlyZXN1'
    'bHRyb3dzGIDMyHIgAygLMjguY2xvdWR0cmFpbC5HZXRRdWVyeVJlc3VsdHNSZXNwb25zZS5RdW'
    'VyeXJlc3VsdHJvd3NFbnRyeVIPcXVlcnlyZXN1bHRyb3dzEkgKD3F1ZXJ5c3RhdGlzdGljcxjZ'
    '0618IAEoCzIbLmNsb3VkdHJhaWwuUXVlcnlTdGF0aXN0aWNzUg9xdWVyeXN0YXRpc3RpY3MSPQ'
    'oLcXVlcnlzdGF0dXMY1vOArwEgASgOMhcuY2xvdWR0cmFpbC5RdWVyeVN0YXR1c1ILcXVlcnlz'
    'dGF0dXMaQgoUUXVlcnlyZXN1bHRyb3dzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdW'
    'UYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getResourcePolicyRequestDescriptor instead')
const GetResourcePolicyRequest$json = {
  '1': 'GetResourcePolicyRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZX'
        'NvdXJjZWFybg==');

@$core.Deprecated('Use getResourcePolicyResponseDescriptor instead')
const GetResourcePolicyResponse$json = {
  '1': 'GetResourcePolicyResponse',
  '2': [
    {
      '1': 'delegatedadminresourcepolicy',
      '3': 510966132,
      '4': 1,
      '5': 9,
      '10': 'delegatedadminresourcepolicy'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
  ],
};

/// Descriptor for `GetResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyResponseDescriptor = $convert.base64Decode(
    'ChlHZXRSZXNvdXJjZVBvbGljeVJlc3BvbnNlEkYKHGRlbGVnYXRlZGFkbWlucmVzb3VyY2Vwb2'
    'xpY3kY9PLS8wEgASgJUhxkZWxlZ2F0ZWRhZG1pbnJlc291cmNlcG9saWN5EiQKC3Jlc291cmNl'
    'YXJuGK342a0BIAEoCVILcmVzb3VyY2Vhcm4SKQoOcmVzb3VyY2Vwb2xpY3kYsJTBByABKAlSDn'
    'Jlc291cmNlcG9saWN5');

@$core.Deprecated('Use getTrailRequestDescriptor instead')
const GetTrailRequest$json = {
  '1': 'GetTrailRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetTrailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrailRequestDescriptor = $convert
    .base64Decode('Cg9HZXRUcmFpbFJlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use getTrailResponseDescriptor instead')
const GetTrailResponse$json = {
  '1': 'GetTrailResponse',
  '2': [
    {
      '1': 'trail',
      '3': 470969828,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.Trail',
      '10': 'trail'
    },
  ],
};

/// Descriptor for `GetTrailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrailResponseDescriptor = $convert.base64Decode(
    'ChBHZXRUcmFpbFJlc3BvbnNlEisKBXRyYWlsGOTbyeABIAEoCzIRLmNsb3VkdHJhaWwuVHJhaW'
    'xSBXRyYWls');

@$core.Deprecated('Use getTrailStatusRequestDescriptor instead')
const GetTrailStatusRequest$json = {
  '1': 'GetTrailStatusRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetTrailStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrailStatusRequestDescriptor =
    $convert.base64Decode(
        'ChVHZXRUcmFpbFN0YXR1c1JlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use getTrailStatusResponseDescriptor instead')
const GetTrailStatusResponse$json = {
  '1': 'GetTrailStatusResponse',
  '2': [
    {'1': 'islogging', '3': 61005519, '4': 1, '5': 8, '10': 'islogging'},
    {
      '1': 'latestcloudwatchlogsdeliveryerror',
      '3': 238605214,
      '4': 1,
      '5': 9,
      '10': 'latestcloudwatchlogsdeliveryerror'
    },
    {
      '1': 'latestcloudwatchlogsdeliverytime',
      '3': 517138001,
      '4': 1,
      '5': 9,
      '10': 'latestcloudwatchlogsdeliverytime'
    },
    {
      '1': 'latestdeliveryattemptsucceeded',
      '3': 432951833,
      '4': 1,
      '5': 9,
      '10': 'latestdeliveryattemptsucceeded'
    },
    {
      '1': 'latestdeliveryattempttime',
      '3': 109095435,
      '4': 1,
      '5': 9,
      '10': 'latestdeliveryattempttime'
    },
    {
      '1': 'latestdeliveryerror',
      '3': 119953649,
      '4': 1,
      '5': 9,
      '10': 'latestdeliveryerror'
    },
    {
      '1': 'latestdeliverytime',
      '3': 320897248,
      '4': 1,
      '5': 9,
      '10': 'latestdeliverytime'
    },
    {
      '1': 'latestdigestdeliveryerror',
      '3': 514209249,
      '4': 1,
      '5': 9,
      '10': 'latestdigestdeliveryerror'
    },
    {
      '1': 'latestdigestdeliverytime',
      '3': 469527024,
      '4': 1,
      '5': 9,
      '10': 'latestdigestdeliverytime'
    },
    {
      '1': 'latestnotificationattemptsucceeded',
      '3': 3764412,
      '4': 1,
      '5': 9,
      '10': 'latestnotificationattemptsucceeded'
    },
    {
      '1': 'latestnotificationattempttime',
      '3': 104168880,
      '4': 1,
      '5': 9,
      '10': 'latestnotificationattempttime'
    },
    {
      '1': 'latestnotificationerror',
      '3': 261077338,
      '4': 1,
      '5': 9,
      '10': 'latestnotificationerror'
    },
    {
      '1': 'latestnotificationtime',
      '3': 2082317,
      '4': 1,
      '5': 9,
      '10': 'latestnotificationtime'
    },
    {
      '1': 'startloggingtime',
      '3': 351696170,
      '4': 1,
      '5': 9,
      '10': 'startloggingtime'
    },
    {
      '1': 'stoploggingtime',
      '3': 469874232,
      '4': 1,
      '5': 9,
      '10': 'stoploggingtime'
    },
    {
      '1': 'timeloggingstarted',
      '3': 18088239,
      '4': 1,
      '5': 9,
      '10': 'timeloggingstarted'
    },
    {
      '1': 'timeloggingstopped',
      '3': 272406495,
      '4': 1,
      '5': 9,
      '10': 'timeloggingstopped'
    },
  ],
};

/// Descriptor for `GetTrailStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrailStatusResponseDescriptor = $convert.base64Decode(
    'ChZHZXRUcmFpbFN0YXR1c1Jlc3BvbnNlEh8KCWlzbG9nZ2luZxjPvYsdIAEoCFIJaXNsb2dnaW'
    '5nEk8KIWxhdGVzdGNsb3Vkd2F0Y2hsb2dzZGVsaXZlcnllcnJvchiep+NxIAEoCVIhbGF0ZXN0'
    'Y2xvdWR3YXRjaGxvZ3NkZWxpdmVyeWVycm9yEk4KIGxhdGVzdGNsb3Vkd2F0Y2hsb2dzZGVsaX'
    'Zlcnl0aW1lGNHMy/YBIAEoCVIgbGF0ZXN0Y2xvdWR3YXRjaGxvZ3NkZWxpdmVyeXRpbWUSSgoe'
    'bGF0ZXN0ZGVsaXZlcnlhdHRlbXB0c3VjY2VlZGVkGJmkuc4BIAEoCVIebGF0ZXN0ZGVsaXZlcn'
    'lhdHRlbXB0c3VjY2VlZGVkEj8KGWxhdGVzdGRlbGl2ZXJ5YXR0ZW1wdHRpbWUYi9SCNCABKAlS'
    'GWxhdGVzdGRlbGl2ZXJ5YXR0ZW1wdHRpbWUSMwoTbGF0ZXN0ZGVsaXZlcnllcnJvchjxsZk5IA'
    'EoCVITbGF0ZXN0ZGVsaXZlcnllcnJvchIyChJsYXRlc3RkZWxpdmVyeXRpbWUY4IGCmQEgASgJ'
    'UhJsYXRlc3RkZWxpdmVyeXRpbWUSQAoZbGF0ZXN0ZGlnZXN0ZGVsaXZlcnllcnJvchjh65j1AS'
    'ABKAlSGWxhdGVzdGRpZ2VzdGRlbGl2ZXJ5ZXJyb3ISPgoYbGF0ZXN0ZGlnZXN0ZGVsaXZlcnl0'
    'aW1lGPDT8d8BIAEoCVIYbGF0ZXN0ZGlnZXN0ZGVsaXZlcnl0aW1lElEKImxhdGVzdG5vdGlmaW'
    'NhdGlvbmF0dGVtcHRzdWNjZWVkZWQYvOHlASABKAlSImxhdGVzdG5vdGlmaWNhdGlvbmF0dGVt'
    'cHRzdWNjZWVkZWQSRwodbGF0ZXN0bm90aWZpY2F0aW9uYXR0ZW1wdHRpbWUYsPvVMSABKAlSHW'
    'xhdGVzdG5vdGlmaWNhdGlvbmF0dGVtcHR0aW1lEjsKF2xhdGVzdG5vdGlmaWNhdGlvbmVycm9y'
    'GNryvnwgASgJUhdsYXRlc3Rub3RpZmljYXRpb25lcnJvchI4ChZsYXRlc3Rub3RpZmljYXRpb2'
    '50aW1lGI2MfyABKAlSFmxhdGVzdG5vdGlmaWNhdGlvbnRpbWUSLgoQc3RhcnRsb2dnaW5ndGlt'
    'ZRiq6tmnASABKAlSEHN0YXJ0bG9nZ2luZ3RpbWUSLAoPc3RvcGxvZ2dpbmd0aW1lGLjshuABIA'
    'EoCVIPc3RvcGxvZ2dpbmd0aW1lEjEKEnRpbWVsb2dnaW5nc3RhcnRlZBivgtAIIAEoCVISdGlt'
    'ZWxvZ2dpbmdzdGFydGVkEjIKEnRpbWVsb2dnaW5nc3RvcHBlZBjfr/KBASABKAlSEnRpbWVsb2'
    'dnaW5nc3RvcHBlZA==');

@$core.Deprecated('Use importFailureListItemDescriptor instead')
const ImportFailureListItem$json = {
  '1': 'ImportFailureListItem',
  '2': [
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'errortype', '3': 398848954, '4': 1, '5': 9, '10': 'errortype'},
    {
      '1': 'lastupdatedtime',
      '3': 177046166,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtime'
    },
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ImportFailureStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ImportFailureListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importFailureListItemDescriptor = $convert.base64Decode(
    'ChVJbXBvcnRGYWlsdXJlTGlzdEl0ZW0SJgoMZXJyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3'
    'JtZXNzYWdlEiAKCWVycm9ydHlwZRi655e+ASABKAlSCWVycm9ydHlwZRIrCg9sYXN0dXBkYXRl'
    'ZHRpbWUYloW2VCABKAlSD2xhc3R1cGRhdGVkdGltZRIeCghsb2NhdGlvbhjHm4LeASABKAlSCG'
    'xvY2F0aW9uEjoKBnN0YXR1cxiQ5PsCIAEoDjIfLmNsb3VkdHJhaWwuSW1wb3J0RmFpbHVyZVN0'
    'YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use importNotFoundExceptionDescriptor instead')
const ImportNotFoundException$json = {
  '1': 'ImportNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ImportNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbXBvcnROb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use importSourceDescriptor instead')
const ImportSource$json = {
  '1': 'ImportSource',
  '2': [
    {
      '1': 's3',
      '3': 100795332,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.S3ImportSource',
      '10': 's3'
    },
  ],
};

/// Descriptor for `ImportSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importSourceDescriptor = $convert.base64Decode(
    'CgxJbXBvcnRTb3VyY2USLQoCczMYxIeIMCABKAsyGi5jbG91ZHRyYWlsLlMzSW1wb3J0U291cm'
    'NlUgJzMw==');

@$core.Deprecated('Use importStatisticsDescriptor instead')
const ImportStatistics$json = {
  '1': 'ImportStatistics',
  '2': [
    {
      '1': 'eventscompleted',
      '3': 159166494,
      '4': 1,
      '5': 3,
      '10': 'eventscompleted'
    },
    {
      '1': 'failedentries',
      '3': 86416685,
      '4': 1,
      '5': 3,
      '10': 'failedentries'
    },
    {
      '1': 'filescompleted',
      '3': 300203834,
      '4': 1,
      '5': 3,
      '10': 'filescompleted'
    },
    {
      '1': 'prefixescompleted',
      '3': 433011203,
      '4': 1,
      '5': 3,
      '10': 'prefixescompleted'
    },
    {
      '1': 'prefixesfound',
      '3': 244228528,
      '4': 1,
      '5': 3,
      '10': 'prefixesfound'
    },
  ],
};

/// Descriptor for `ImportStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importStatisticsDescriptor = $convert.base64Decode(
    'ChBJbXBvcnRTdGF0aXN0aWNzEisKD2V2ZW50c2NvbXBsZXRlZBie4PJLIAEoA1IPZXZlbnRzY2'
    '9tcGxldGVkEicKDWZhaWxlZGVudHJpZXMYrbqaKSABKANSDWZhaWxlZGVudHJpZXMSKgoOZmls'
    'ZXNjb21wbGV0ZWQYuv6SjwEgASgDUg5maWxlc2NvbXBsZXRlZBIwChFwcmVmaXhlc2NvbXBsZX'
    'RlZBiD9LzOASABKANSEXByZWZpeGVzY29tcGxldGVkEicKDXByZWZpeGVzZm91bmQYsMO6dCAB'
    'KANSDXByZWZpeGVzZm91bmQ=');

@$core.Deprecated('Use importsListItemDescriptor instead')
const ImportsListItem$json = {
  '1': 'ImportsListItem',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'destinations', '3': 404385861, '4': 3, '5': 9, '10': 'destinations'},
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importstatus',
      '3': 129077631,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `ImportsListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importsListItemDescriptor = $convert.base64Decode(
    'Cg9JbXBvcnRzTGlzdEl0ZW0SLgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEGNyZWF0ZW'
    'R0aW1lc3RhbXASJgoMZGVzdGluYXRpb25zGMXg6cABIAMoCVIMZGVzdGluYXRpb25zEh4KCGlt'
    'cG9ydGlkGNqUrMgBIAEoCVIIaW1wb3J0aWQSPwoMaW1wb3J0c3RhdHVzGP+ixj0gASgOMhguY2'
    'xvdWR0cmFpbC5JbXBvcnRTdGF0dXNSDGltcG9ydHN0YXR1cxItChB1cGRhdGVkdGltZXN0YW1w'
    'GIHjkxUgASgJUhB1cGRhdGVkdGltZXN0YW1w');

@$core.Deprecated('Use inactiveEventDataStoreExceptionDescriptor instead')
const InactiveEventDataStoreException$json = {
  '1': 'InactiveEventDataStoreException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InactiveEventDataStoreException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inactiveEventDataStoreExceptionDescriptor =
    $convert.base64Decode(
        'Ch9JbmFjdGl2ZUV2ZW50RGF0YVN0b3JlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use inactiveQueryExceptionDescriptor instead')
const InactiveQueryException$json = {
  '1': 'InactiveQueryException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InactiveQueryException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inactiveQueryExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbmFjdGl2ZVF1ZXJ5RXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use ingestionStatusDescriptor instead')
const IngestionStatus$json = {
  '1': 'IngestionStatus',
  '2': [
    {
      '1': 'latestingestionattempteventid',
      '3': 204299739,
      '4': 1,
      '5': 9,
      '10': 'latestingestionattempteventid'
    },
    {
      '1': 'latestingestionattempttime',
      '3': 313195389,
      '4': 1,
      '5': 9,
      '10': 'latestingestionattempttime'
    },
    {
      '1': 'latestingestionerrorcode',
      '3': 415423856,
      '4': 1,
      '5': 9,
      '10': 'latestingestionerrorcode'
    },
    {
      '1': 'latestingestionsuccesseventid',
      '3': 31034217,
      '4': 1,
      '5': 9,
      '10': 'latestingestionsuccesseventid'
    },
    {
      '1': 'latestingestionsuccesstime',
      '3': 184241615,
      '4': 1,
      '5': 9,
      '10': 'latestingestionsuccesstime'
    },
  ],
};

/// Descriptor for `IngestionStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ingestionStatusDescriptor = $convert.base64Decode(
    'Cg9Jbmdlc3Rpb25TdGF0dXMSRwodbGF0ZXN0aW5nZXN0aW9uYXR0ZW1wdGV2ZW50aWQY27u1YS'
    'ABKAlSHWxhdGVzdGluZ2VzdGlvbmF0dGVtcHRldmVudGlkEkIKGmxhdGVzdGluZ2VzdGlvbmF0'
    'dGVtcHR0aW1lGP32q5UBIAEoCVIabGF0ZXN0aW5nZXN0aW9uYXR0ZW1wdHRpbWUSPgoYbGF0ZX'
    'N0aW5nZXN0aW9uZXJyb3Jjb2RlGPC6i8YBIAEoCVIYbGF0ZXN0aW5nZXN0aW9uZXJyb3Jjb2Rl'
    'EkcKHWxhdGVzdGluZ2VzdGlvbnN1Y2Nlc3NldmVudGlkGOmW5g4gASgJUh1sYXRlc3Rpbmdlc3'
    'Rpb25zdWNjZXNzZXZlbnRpZBJBChpsYXRlc3Rpbmdlc3Rpb25zdWNjZXNzdGltZRjPm+1XIAEo'
    'CVIabGF0ZXN0aW5nZXN0aW9uc3VjY2Vzc3RpbWU=');

@$core.Deprecated('Use insightNotEnabledExceptionDescriptor instead')
const InsightNotEnabledException$json = {
  '1': 'InsightNotEnabledException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InsightNotEnabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insightNotEnabledExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnNpZ2h0Tm90RW5hYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use insightSelectorDescriptor instead')
const InsightSelector$json = {
  '1': 'InsightSelector',
  '2': [
    {
      '1': 'eventcategories',
      '3': 3676820,
      '4': 3,
      '5': 14,
      '6': '.cloudtrail.SourceEventCategory',
      '10': 'eventcategories'
    },
    {
      '1': 'insighttype',
      '3': 530375860,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.InsightType',
      '10': 'insighttype'
    },
  ],
};

/// Descriptor for `InsightSelector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insightSelectorDescriptor = $convert.base64Decode(
    'Cg9JbnNpZ2h0U2VsZWN0b3ISTAoPZXZlbnRjYXRlZ29yaWVzGJS14AEgAygOMh8uY2xvdWR0cm'
    'FpbC5Tb3VyY2VFdmVudENhdGVnb3J5Ug9ldmVudGNhdGVnb3JpZXMSPQoLaW5zaWdodHR5cGUY'
    'tMnz/AEgASgOMhcuY2xvdWR0cmFpbC5JbnNpZ2h0VHlwZVILaW5zaWdodHR5cGU=');

@$core.Deprecated(
    'Use insufficientDependencyServiceAccessPermissionExceptionDescriptor instead')
const InsufficientDependencyServiceAccessPermissionException$json = {
  '1': 'InsufficientDependencyServiceAccessPermissionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InsufficientDependencyServiceAccessPermissionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    insufficientDependencyServiceAccessPermissionExceptionDescriptor =
    $convert.base64Decode(
        'CjZJbnN1ZmZpY2llbnREZXBlbmRlbmN5U2VydmljZUFjY2Vzc1Blcm1pc3Npb25FeGNlcHRpb2'
        '4SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use insufficientEncryptionPolicyExceptionDescriptor instead')
const InsufficientEncryptionPolicyException$json = {
  '1': 'InsufficientEncryptionPolicyException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InsufficientEncryptionPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insufficientEncryptionPolicyExceptionDescriptor =
    $convert.base64Decode(
        'CiVJbnN1ZmZpY2llbnRFbmNyeXB0aW9uUG9saWN5RXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cC'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use insufficientIAMAccessPermissionExceptionDescriptor instead')
const InsufficientIAMAccessPermissionException$json = {
  '1': 'InsufficientIAMAccessPermissionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InsufficientIAMAccessPermissionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insufficientIAMAccessPermissionExceptionDescriptor =
    $convert.base64Decode(
        'CihJbnN1ZmZpY2llbnRJQU1BY2Nlc3NQZXJtaXNzaW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use insufficientS3BucketPolicyExceptionDescriptor instead')
const InsufficientS3BucketPolicyException$json = {
  '1': 'InsufficientS3BucketPolicyException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InsufficientS3BucketPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insufficientS3BucketPolicyExceptionDescriptor =
    $convert.base64Decode(
        'CiNJbnN1ZmZpY2llbnRTM0J1Y2tldFBvbGljeUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use insufficientSnsTopicPolicyExceptionDescriptor instead')
const InsufficientSnsTopicPolicyException$json = {
  '1': 'InsufficientSnsTopicPolicyException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InsufficientSnsTopicPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insufficientSnsTopicPolicyExceptionDescriptor =
    $convert.base64Decode(
        'CiNJbnN1ZmZpY2llbnRTbnNUb3BpY1BvbGljeUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use invalidCloudWatchLogsLogGroupArnExceptionDescriptor instead')
const InvalidCloudWatchLogsLogGroupArnException$json = {
  '1': 'InvalidCloudWatchLogsLogGroupArnException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidCloudWatchLogsLogGroupArnException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    invalidCloudWatchLogsLogGroupArnExceptionDescriptor = $convert.base64Decode(
        'CilJbnZhbGlkQ2xvdWRXYXRjaExvZ3NMb2dHcm91cEFybkV4Y2VwdGlvbhIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidCloudWatchLogsRoleArnExceptionDescriptor instead')
const InvalidCloudWatchLogsRoleArnException$json = {
  '1': 'InvalidCloudWatchLogsRoleArnException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidCloudWatchLogsRoleArnException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidCloudWatchLogsRoleArnExceptionDescriptor =
    $convert.base64Decode(
        'CiVJbnZhbGlkQ2xvdWRXYXRjaExvZ3NSb2xlQXJuRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cC'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidDateRangeExceptionDescriptor instead')
const InvalidDateRangeException$json = {
  '1': 'InvalidDateRangeException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDateRangeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDateRangeExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkRGF0ZVJhbmdlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidEventCategoryExceptionDescriptor instead')
const InvalidEventCategoryException$json = {
  '1': 'InvalidEventCategoryException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEventCategoryException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEventCategoryExceptionDescriptor =
    $convert.base64Decode(
        'Ch1JbnZhbGlkRXZlbnRDYXRlZ29yeUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core
    .Deprecated('Use invalidEventDataStoreCategoryExceptionDescriptor instead')
const InvalidEventDataStoreCategoryException$json = {
  '1': 'InvalidEventDataStoreCategoryException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEventDataStoreCategoryException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEventDataStoreCategoryExceptionDescriptor =
    $convert.base64Decode(
        'CiZJbnZhbGlkRXZlbnREYXRhU3RvcmVDYXRlZ29yeUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3'
        'AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidEventDataStoreStatusExceptionDescriptor instead')
const InvalidEventDataStoreStatusException$json = {
  '1': 'InvalidEventDataStoreStatusException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEventDataStoreStatusException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEventDataStoreStatusExceptionDescriptor =
    $convert.base64Decode(
        'CiRJbnZhbGlkRXZlbnREYXRhU3RvcmVTdGF0dXNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidEventSelectorsExceptionDescriptor instead')
const InvalidEventSelectorsException$json = {
  '1': 'InvalidEventSelectorsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEventSelectorsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEventSelectorsExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkRXZlbnRTZWxlY3RvcnNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidHomeRegionExceptionDescriptor instead')
const InvalidHomeRegionException$json = {
  '1': 'InvalidHomeRegionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidHomeRegionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidHomeRegionExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkSG9tZVJlZ2lvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidImportSourceExceptionDescriptor instead')
const InvalidImportSourceException$json = {
  '1': 'InvalidImportSourceException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidImportSourceException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidImportSourceExceptionDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkSW1wb3J0U291cmNlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidInsightSelectorsExceptionDescriptor instead')
const InvalidInsightSelectorsException$json = {
  '1': 'InvalidInsightSelectorsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInsightSelectorsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInsightSelectorsExceptionDescriptor =
    $convert.base64Decode(
        'CiBJbnZhbGlkSW5zaWdodFNlbGVjdG9yc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use invalidKmsKeyIdExceptionDescriptor instead')
const InvalidKmsKeyIdException$json = {
  '1': 'InvalidKmsKeyIdException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidKmsKeyIdException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidKmsKeyIdExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkS21zS2V5SWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidLookupAttributesExceptionDescriptor instead')
const InvalidLookupAttributesException$json = {
  '1': 'InvalidLookupAttributesException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidLookupAttributesException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidLookupAttributesExceptionDescriptor =
    $convert.base64Decode(
        'CiBJbnZhbGlkTG9va3VwQXR0cmlidXRlc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use invalidMaxResultsExceptionDescriptor instead')
const InvalidMaxResultsException$json = {
  '1': 'InvalidMaxResultsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidMaxResultsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidMaxResultsExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkTWF4UmVzdWx0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidNextTokenExceptionDescriptor instead')
const InvalidNextTokenException$json = {
  '1': 'InvalidNextTokenException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidNextTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidNextTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkTmV4dFRva2VuRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidParameterCombinationExceptionDescriptor instead')
const InvalidParameterCombinationException$json = {
  '1': 'InvalidParameterCombinationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterCombinationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterCombinationExceptionDescriptor =
    $convert.base64Decode(
        'CiRJbnZhbGlkUGFyYW1ldGVyQ29tYmluYXRpb25FeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidParameterExceptionDescriptor instead')
const InvalidParameterException$json = {
  '1': 'InvalidParameterException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkUGFyYW1ldGVyRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidQueryStatementExceptionDescriptor instead')
const InvalidQueryStatementException$json = {
  '1': 'InvalidQueryStatementException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidQueryStatementException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidQueryStatementExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkUXVlcnlTdGF0ZW1lbnRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidQueryStatusExceptionDescriptor instead')
const InvalidQueryStatusException$json = {
  '1': 'InvalidQueryStatusException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidQueryStatusException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidQueryStatusExceptionDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkUXVlcnlTdGF0dXNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use invalidS3BucketNameExceptionDescriptor instead')
const InvalidS3BucketNameException$json = {
  '1': 'InvalidS3BucketNameException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidS3BucketNameException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidS3BucketNameExceptionDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkUzNCdWNrZXROYW1lRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidS3PrefixExceptionDescriptor instead')
const InvalidS3PrefixException$json = {
  '1': 'InvalidS3PrefixException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidS3PrefixException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidS3PrefixExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkUzNQcmVmaXhFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidSnsTopicNameExceptionDescriptor instead')
const InvalidSnsTopicNameException$json = {
  '1': 'InvalidSnsTopicNameException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSnsTopicNameException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSnsTopicNameExceptionDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkU25zVG9waWNOYW1lRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidSourceExceptionDescriptor instead')
const InvalidSourceException$json = {
  '1': 'InvalidSourceException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSourceException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSourceExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbnZhbGlkU291cmNlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidTagParameterExceptionDescriptor instead')
const InvalidTagParameterException$json = {
  '1': 'InvalidTagParameterException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTagParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTagParameterExceptionDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkVGFnUGFyYW1ldGVyRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidTimeRangeExceptionDescriptor instead')
const InvalidTimeRangeException$json = {
  '1': 'InvalidTimeRangeException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTimeRangeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTimeRangeExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkVGltZVJhbmdlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidTokenExceptionDescriptor instead')
const InvalidTokenException$json = {
  '1': 'InvalidTokenException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTokenExceptionDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkVG9rZW5FeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidTrailNameExceptionDescriptor instead')
const InvalidTrailNameException$json = {
  '1': 'InvalidTrailNameException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTrailNameException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTrailNameExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkVHJhaWxOYW1lRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use kmsExceptionDescriptor instead')
const KmsException$json = {
  '1': 'KmsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsExceptionDescriptor = $convert.base64Decode(
    'CgxLbXNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use kmsKeyDisabledExceptionDescriptor instead')
const KmsKeyDisabledException$json = {
  '1': 'KmsKeyDisabledException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsKeyDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsKeyDisabledExceptionDescriptor =
    $convert.base64Decode(
        'ChdLbXNLZXlEaXNhYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use kmsKeyNotFoundExceptionDescriptor instead')
const KmsKeyNotFoundException$json = {
  '1': 'KmsKeyNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsKeyNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsKeyNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChdLbXNLZXlOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use listChannelsRequestDescriptor instead')
const ListChannelsRequest$json = {
  '1': 'ListChannelsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListChannelsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listChannelsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0Q2hhbm5lbHNSZXF1ZXN0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bH'
    'RzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listChannelsResponseDescriptor instead')
const ListChannelsResponse$json = {
  '1': 'ListChannelsResponse',
  '2': [
    {
      '1': 'channels',
      '3': 155227286,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Channel',
      '10': 'channels'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListChannelsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listChannelsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0Q2hhbm5lbHNSZXNwb25zZRIyCghjaGFubmVscxiWqYJKIAMoCzITLmNsb3VkdHJhaW'
    'wuQ2hhbm5lbFIIY2hhbm5lbHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listDashboardsRequestDescriptor instead')
const ListDashboardsRequest$json = {
  '1': 'ListDashboardsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DashboardType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ListDashboardsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDashboardsRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0RGFzaGJvYXJkc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3'
    'VsdHMSIgoKbmFtZXByZWZpeBib87ysASABKAlSCm5hbWVwcmVmaXgSHwoJbmV4dHRva2VuGP6E'
    'umcgASgJUgluZXh0dG9rZW4SMQoEdHlwZRjuoNeKASABKA4yGS5jbG91ZHRyYWlsLkRhc2hib2'
    'FyZFR5cGVSBHR5cGU=');

@$core.Deprecated('Use listDashboardsResponseDescriptor instead')
const ListDashboardsResponse$json = {
  '1': 'ListDashboardsResponse',
  '2': [
    {
      '1': 'dashboards',
      '3': 261493913,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.DashboardDetail',
      '10': 'dashboards'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDashboardsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDashboardsResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0RGFzaGJvYXJkc1Jlc3BvbnNlEj4KCmRhc2hib2FyZHMYmanYfCADKAsyGy5jbG91ZH'
    'RyYWlsLkRhc2hib2FyZERldGFpbFIKZGFzaGJvYXJkcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlS'
    'CW5leHR0b2tlbg==');

@$core.Deprecated('Use listEventDataStoresRequestDescriptor instead')
const ListEventDataStoresRequest$json = {
  '1': 'ListEventDataStoresRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEventDataStoresRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventDataStoresRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0RXZlbnREYXRhU3RvcmVzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
        'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listEventDataStoresResponseDescriptor instead')
const ListEventDataStoresResponse$json = {
  '1': 'ListEventDataStoresResponse',
  '2': [
    {
      '1': 'eventdatastores',
      '3': 152154314,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.EventDataStore',
      '10': 'eventdatastores'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEventDataStoresResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventDataStoresResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0RXZlbnREYXRhU3RvcmVzUmVzcG9uc2USRwoPZXZlbnRkYXRhc3RvcmVzGMrhxkggAy'
        'gLMhouY2xvdWR0cmFpbC5FdmVudERhdGFTdG9yZVIPZXZlbnRkYXRhc3RvcmVzEh8KCW5leHR0'
        'b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listImportFailuresRequestDescriptor instead')
const ListImportFailuresRequest$json = {
  '1': 'ListImportFailuresRequest',
  '2': [
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListImportFailuresRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportFailuresRequestDescriptor = $convert.base64Decode(
    'ChlMaXN0SW1wb3J0RmFpbHVyZXNSZXF1ZXN0Eh4KCGltcG9ydGlkGNqUrMgBIAEoCVIIaW1wb3'
    'J0aWQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6E'
    'umcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listImportFailuresResponseDescriptor instead')
const ListImportFailuresResponse$json = {
  '1': 'ListImportFailuresResponse',
  '2': [
    {
      '1': 'failures',
      '3': 335467271,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.ImportFailureListItem',
      '10': 'failures'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListImportFailuresResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportFailuresResponseDescriptor =
    $convert.base64Decode(
        'ChpMaXN0SW1wb3J0RmFpbHVyZXNSZXNwb25zZRJBCghmYWlsdXJlcxiHpvufASADKAsyIS5jbG'
        '91ZHRyYWlsLkltcG9ydEZhaWx1cmVMaXN0SXRlbVIIZmFpbHVyZXMSHwoJbmV4dHRva2VuGP6E'
        'umcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listImportsRequestDescriptor instead')
const ListImportsRequest$json = {
  '1': 'ListImportsRequest',
  '2': [
    {'1': 'destination', '3': 457443680, '4': 1, '5': 9, '10': 'destination'},
    {
      '1': 'importstatus',
      '3': 129077631,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ImportStatus',
      '10': 'importstatus'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListImportsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportsRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0SW1wb3J0c1JlcXVlc3QSJAoLZGVzdGluYXRpb24Y4JKQ2gEgASgJUgtkZXN0aW5hdG'
    'lvbhI/CgxpbXBvcnRzdGF0dXMY/6LGPSABKA4yGC5jbG91ZHRyYWlsLkltcG9ydFN0YXR1c1IM'
    'aW1wb3J0c3RhdHVzEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh8KCW5leH'
    'R0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listImportsResponseDescriptor instead')
const ListImportsResponse$json = {
  '1': 'ListImportsResponse',
  '2': [
    {
      '1': 'imports',
      '3': 497717894,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.ImportsListItem',
      '10': 'imports'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListImportsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportsResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0SW1wb3J0c1Jlc3BvbnNlEjkKB2ltcG9ydHMYhqWq7QEgAygLMhsuY2xvdWR0cmFpbC'
    '5JbXBvcnRzTGlzdEl0ZW1SB2ltcG9ydHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9r'
    'ZW4=');

@$core.Deprecated('Use listInsightsDataRequestDescriptor instead')
const ListInsightsDataRequest$json = {
  '1': 'ListInsightsDataRequest',
  '2': [
    {
      '1': 'datatype',
      '3': 67988590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ListInsightsDataType',
      '10': 'datatype'
    },
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.ListInsightsDataRequest.DimensionsEntry',
      '10': 'dimensions'
    },
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {
      '1': 'insightsource',
      '3': 438381967,
      '4': 1,
      '5': 9,
      '10': 'insightsource'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
  ],
  '3': [ListInsightsDataRequest_DimensionsEntry$json],
};

@$core.Deprecated('Use listInsightsDataRequestDescriptor instead')
const ListInsightsDataRequest_DimensionsEntry$json = {
  '1': 'DimensionsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListInsightsDataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInsightsDataRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0SW5zaWdodHNEYXRhUmVxdWVzdBI/CghkYXRhdHlwZRju2LUgIAEoDjIgLmNsb3VkdH'
    'JhaWwuTGlzdEluc2lnaHRzRGF0YVR5cGVSCGRhdGF0eXBlElcKCmRpbWVuc2lvbnMY0Zvf3AEg'
    'AygLMjMuY2xvdWR0cmFpbC5MaXN0SW5zaWdodHNEYXRhUmVxdWVzdC5EaW1lbnNpb25zRW50cn'
    'lSCmRpbWVuc2lvbnMSGwoHZW5kdGltZRjM77weIAEoCVIHZW5kdGltZRIoCg1pbnNpZ2h0c291'
    'cmNlGI/bhNEBIAEoCVINaW5zaWdodHNvdXJjZRIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
    'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIgCglzdGFydHRpbWUY'
    '77TlsAEgASgJUglzdGFydHRpbWUaPQoPRGltZW5zaW9uc0VudHJ5EhAKA2tleRgBIAEoCVIDa2'
    'V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use listInsightsDataResponseDescriptor instead')
const ListInsightsDataResponse$json = {
  '1': 'ListInsightsDataResponse',
  '2': [
    {
      '1': 'events',
      '3': 3416229,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Event',
      '10': 'events'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListInsightsDataResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInsightsDataResponseDescriptor =
    $convert.base64Decode(
        'ChhMaXN0SW5zaWdodHNEYXRhUmVzcG9uc2USLAoGZXZlbnRzGKXB0AEgAygLMhEuY2xvdWR0cm'
        'FpbC5FdmVudFIGZXZlbnRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listInsightsMetricDataRequestDescriptor instead')
const ListInsightsMetricDataRequest$json = {
  '1': 'ListInsightsMetricDataRequest',
  '2': [
    {
      '1': 'datatype',
      '3': 67988590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.InsightsMetricDataType',
      '10': 'datatype'
    },
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'eventname', '3': 264746781, '4': 1, '5': 9, '10': 'eventname'},
    {'1': 'eventsource', '3': 37841339, '4': 1, '5': 9, '10': 'eventsource'},
    {
      '1': 'insighttype',
      '3': 530375860,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.InsightType',
      '10': 'insighttype'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {'1': 'trailname', '3': 507774985, '4': 1, '5': 9, '10': 'trailname'},
  ],
};

/// Descriptor for `ListInsightsMetricDataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInsightsMetricDataRequestDescriptor = $convert.base64Decode(
    'Ch1MaXN0SW5zaWdodHNNZXRyaWNEYXRhUmVxdWVzdBJBCghkYXRhdHlwZRju2LUgIAEoDjIiLm'
    'Nsb3VkdHJhaWwuSW5zaWdodHNNZXRyaWNEYXRhVHlwZVIIZGF0YXR5cGUSGwoHZW5kdGltZRjM'
    '77weIAEoCVIHZW5kdGltZRIfCgllcnJvcmNvZGUYmdbDECABKAlSCWVycm9yY29kZRIfCglldm'
    'VudG5hbWUYne6efiABKAlSCWV2ZW50bmFtZRIjCgtldmVudHNvdXJjZRi704USIAEoCVILZXZl'
    'bnRzb3VyY2USPQoLaW5zaWdodHR5cGUYtMnz/AEgASgOMhcuY2xvdWR0cmFpbC5JbnNpZ2h0VH'
    'lwZVILaW5zaWdodHR5cGUSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJ'
    'bmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SGQoGcGVyaW9kGKWIkjkgASgFUgZwZXJpb2'
    'QSIAoJc3RhcnR0aW1lGO+05bABIAEoCVIJc3RhcnR0aW1lEiAKCXRyYWlsbmFtZRiJkJDyASAB'
    'KAlSCXRyYWlsbmFtZQ==');

@$core.Deprecated('Use listInsightsMetricDataResponseDescriptor instead')
const ListInsightsMetricDataResponse$json = {
  '1': 'ListInsightsMetricDataResponse',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'eventname', '3': 264746781, '4': 1, '5': 9, '10': 'eventname'},
    {'1': 'eventsource', '3': 37841339, '4': 1, '5': 9, '10': 'eventsource'},
    {
      '1': 'insighttype',
      '3': 530375860,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.InsightType',
      '10': 'insighttype'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'timestamps', '3': 213534737, '4': 3, '5': 9, '10': 'timestamps'},
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 1, '10': 'values'},
  ],
};

/// Descriptor for `ListInsightsMetricDataResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInsightsMetricDataResponseDescriptor = $convert.base64Decode(
    'Ch5MaXN0SW5zaWdodHNNZXRyaWNEYXRhUmVzcG9uc2USHwoJZXJyb3Jjb2RlGJnWwxAgASgJUg'
    'llcnJvcmNvZGUSHwoJZXZlbnRuYW1lGJ3unn4gASgJUglldmVudG5hbWUSIwoLZXZlbnRzb3Vy'
    'Y2UYu9OFEiABKAlSC2V2ZW50c291cmNlEj0KC2luc2lnaHR0eXBlGLTJ8/wBIAEoDjIXLmNsb3'
    'VkdHJhaWwuSW5zaWdodFR5cGVSC2luc2lnaHR0eXBlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJ'
    'bmV4dHRva2VuEiEKCnRpbWVzdGFtcHMYkZDpZSADKAlSCnRpbWVzdGFtcHMSHQoIdHJhaWxhcm'
    '4Y94rwEiABKAlSCHRyYWlsYXJuEhkKBnZhbHVlcxjcxLRqIAMoAVIGdmFsdWVz');

@$core.Deprecated('Use listPublicKeysRequestDescriptor instead')
const ListPublicKeysRequest$json = {
  '1': 'ListPublicKeysRequest',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
  ],
};

/// Descriptor for `ListPublicKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPublicKeysRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0UHVibGljS2V5c1JlcXVlc3QSGwoHZW5kdGltZRjM77weIAEoCVIHZW5kdGltZRIfCg'
    'luZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIgCglzdGFydHRpbWUY77TlsAEgASgJUglz'
    'dGFydHRpbWU=');

@$core.Deprecated('Use listPublicKeysResponseDescriptor instead')
const ListPublicKeysResponse$json = {
  '1': 'ListPublicKeysResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'publickeylist',
      '3': 129933128,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.PublicKey',
      '10': 'publickeylist'
    },
  ],
};

/// Descriptor for `ListPublicKeysResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPublicKeysResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0UHVibGljS2V5c1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2'
    'VuEj4KDXB1YmxpY2tleWxpc3QYyL76PSADKAsyFS5jbG91ZHRyYWlsLlB1YmxpY0tleVINcHVi'
    'bGlja2V5bGlzdA==');

@$core.Deprecated('Use listQueriesRequestDescriptor instead')
const ListQueriesRequest$json = {
  '1': 'ListQueriesRequest',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'querystatus',
      '3': 367016406,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.QueryStatus',
      '10': 'querystatus'
    },
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
  ],
};

/// Descriptor for `ListQueriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueriesRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0UXVlcmllc1JlcXVlc3QSGwoHZW5kdGltZRjM77weIAEoCVIHZW5kdGltZRIpCg5ldm'
    'VudGRhdGFzdG9yZRjB251BIAEoCVIOZXZlbnRkYXRhc3RvcmUSIgoKbWF4cmVzdWx0cxiyqJuD'
    'ASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SPQoLcX'
    'VlcnlzdGF0dXMY1vOArwEgASgOMhcuY2xvdWR0cmFpbC5RdWVyeVN0YXR1c1ILcXVlcnlzdGF0'
    'dXMSIAoJc3RhcnR0aW1lGO+05bABIAEoCVIJc3RhcnR0aW1l');

@$core.Deprecated('Use listQueriesResponseDescriptor instead')
const ListQueriesResponse$json = {
  '1': 'ListQueriesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queries',
      '3': 305659620,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Query',
      '10': 'queries'
    },
  ],
};

/// Descriptor for `ListQueriesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueriesResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0UXVlcmllc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEi'
    '8KB3F1ZXJpZXMY5P3fkQEgAygLMhEuY2xvdWR0cmFpbC5RdWVyeVIHcXVlcmllcw==');

@$core.Deprecated('Use listTagsRequestDescriptor instead')
const ListTagsRequest$json = {
  '1': 'ListTagsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourceidlist',
      '3': 272650681,
      '4': 3,
      '5': 9,
      '10': 'resourceidlist'
    },
  ],
};

/// Descriptor for `ListTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsRequestDescriptor = $convert.base64Decode(
    'Cg9MaXN0VGFnc1JlcXVlc3QSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SKgoOcm'
    'Vzb3VyY2VpZGxpc3QYuaOBggEgAygJUg5yZXNvdXJjZWlkbGlzdA==');

@$core.Deprecated('Use listTagsResponseDescriptor instead')
const ListTagsResponse$json = {
  '1': 'ListTagsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcetaglist',
      '3': 301974828,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.ResourceTag',
      '10': 'resourcetaglist'
    },
  ],
};

/// Descriptor for `ListTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsResponseDescriptor = $convert.base64Decode(
    'ChBMaXN0VGFnc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEkUKD3'
    'Jlc291cmNldGFnbGlzdBisiv+PASADKAsyFy5jbG91ZHRyYWlsLlJlc291cmNlVGFnUg9yZXNv'
    'dXJjZXRhZ2xpc3Q=');

@$core.Deprecated('Use listTrailsRequestDescriptor instead')
const ListTrailsRequest$json = {
  '1': 'ListTrailsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListTrailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrailsRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0VHJhaWxzUmVxdWVzdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listTrailsResponseDescriptor instead')
const ListTrailsResponse$json = {
  '1': 'ListTrailsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'trails',
      '3': 16939441,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.TrailInfo',
      '10': 'trails'
    },
  ],
};

/// Descriptor for `ListTrailsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTrailsResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0VHJhaWxzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SMA'
    'oGdHJhaWxzGLHziQggAygLMhUuY2xvdWR0cmFpbC5UcmFpbEluZm9SBnRyYWlscw==');

@$core.Deprecated('Use lookupAttributeDescriptor instead')
const LookupAttribute$json = {
  '1': 'LookupAttribute',
  '2': [
    {
      '1': 'attributekey',
      '3': 104472119,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.LookupAttributeKey',
      '10': 'attributekey'
    },
    {
      '1': 'attributevalue',
      '3': 96769221,
      '4': 1,
      '5': 9,
      '10': 'attributevalue'
    },
  ],
};

/// Descriptor for `LookupAttribute`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lookupAttributeDescriptor = $convert.base64Decode(
    'Cg9Mb29rdXBBdHRyaWJ1dGUSRQoMYXR0cmlidXRla2V5GLe86DEgASgOMh4uY2xvdWR0cmFpbC'
    '5Mb29rdXBBdHRyaWJ1dGVLZXlSDGF0dHJpYnV0ZWtleRIpCg5hdHRyaWJ1dGV2YWx1ZRjFqZIu'
    'IAEoCVIOYXR0cmlidXRldmFsdWU=');

@$core.Deprecated('Use lookupEventsRequestDescriptor instead')
const LookupEventsRequest$json = {
  '1': 'LookupEventsRequest',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {
      '1': 'eventcategory',
      '3': 164668724,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.EventCategory',
      '10': 'eventcategory'
    },
    {
      '1': 'lookupattributes',
      '3': 162162567,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.LookupAttribute',
      '10': 'lookupattributes'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
  ],
};

/// Descriptor for `LookupEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lookupEventsRequestDescriptor = $convert.base64Decode(
    'ChNMb29rdXBFdmVudHNSZXF1ZXN0EhsKB2VuZHRpbWUYzO+8HiABKAlSB2VuZHRpbWUSQgoNZX'
    'ZlbnRjYXRlZ29yeRi0ysJOIAEoDjIZLmNsb3VkdHJhaWwuRXZlbnRDYXRlZ29yeVINZXZlbnRj'
    'YXRlZ29yeRJKChBsb29rdXBhdHRyaWJ1dGVzGIfPqU0gAygLMhsuY2xvdWR0cmFpbC5Mb29rdX'
    'BBdHRyaWJ1dGVSEGxvb2t1cGF0dHJpYnV0ZXMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1h'
    'eHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIAoJc3RhcnR0aW1lGO'
    '+05bABIAEoCVIJc3RhcnR0aW1l');

@$core.Deprecated('Use lookupEventsResponseDescriptor instead')
const LookupEventsResponse$json = {
  '1': 'LookupEventsResponse',
  '2': [
    {
      '1': 'events',
      '3': 3416229,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Event',
      '10': 'events'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `LookupEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lookupEventsResponseDescriptor = $convert.base64Decode(
    'ChRMb29rdXBFdmVudHNSZXNwb25zZRIsCgZldmVudHMYpcHQASADKAsyES5jbG91ZHRyYWlsLk'
    'V2ZW50UgZldmVudHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use maxConcurrentQueriesExceptionDescriptor instead')
const MaxConcurrentQueriesException$json = {
  '1': 'MaxConcurrentQueriesException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MaxConcurrentQueriesException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maxConcurrentQueriesExceptionDescriptor =
    $convert.base64Decode(
        'Ch1NYXhDb25jdXJyZW50UXVlcmllc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core
    .Deprecated('Use maximumNumberOfTrailsExceededExceptionDescriptor instead')
const MaximumNumberOfTrailsExceededException$json = {
  '1': 'MaximumNumberOfTrailsExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MaximumNumberOfTrailsExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maximumNumberOfTrailsExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiZNYXhpbXVtTnVtYmVyT2ZUcmFpbHNFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3'
        'AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use noManagementAccountSLRExistsExceptionDescriptor instead')
const NoManagementAccountSLRExistsException$json = {
  '1': 'NoManagementAccountSLRExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoManagementAccountSLRExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noManagementAccountSLRExistsExceptionDescriptor =
    $convert.base64Decode(
        'CiVOb01hbmFnZW1lbnRBY2NvdW50U0xSRXhpc3RzRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cC'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use notOrganizationManagementAccountExceptionDescriptor instead')
const NotOrganizationManagementAccountException$json = {
  '1': 'NotOrganizationManagementAccountException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NotOrganizationManagementAccountException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    notOrganizationManagementAccountExceptionDescriptor = $convert.base64Decode(
        'CilOb3RPcmdhbml6YXRpb25NYW5hZ2VtZW50QWNjb3VudEV4Y2VwdGlvbhIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use notOrganizationMasterAccountExceptionDescriptor instead')
const NotOrganizationMasterAccountException$json = {
  '1': 'NotOrganizationMasterAccountException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NotOrganizationMasterAccountException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notOrganizationMasterAccountExceptionDescriptor =
    $convert.base64Decode(
        'CiVOb3RPcmdhbml6YXRpb25NYXN0ZXJBY2NvdW50RXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cC'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use operationNotPermittedExceptionDescriptor instead')
const OperationNotPermittedException$json = {
  '1': 'OperationNotPermittedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OperationNotPermittedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List operationNotPermittedExceptionDescriptor =
    $convert.base64Decode(
        'Ch5PcGVyYXRpb25Ob3RQZXJtaXR0ZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated(
    'Use organizationNotInAllFeaturesModeExceptionDescriptor instead')
const OrganizationNotInAllFeaturesModeException$json = {
  '1': 'OrganizationNotInAllFeaturesModeException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OrganizationNotInAllFeaturesModeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    organizationNotInAllFeaturesModeExceptionDescriptor = $convert.base64Decode(
        'CilPcmdhbml6YXRpb25Ob3RJbkFsbEZlYXR1cmVzTW9kZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use organizationsNotInUseExceptionDescriptor instead')
const OrganizationsNotInUseException$json = {
  '1': 'OrganizationsNotInUseException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OrganizationsNotInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List organizationsNotInUseExceptionDescriptor =
    $convert.base64Decode(
        'Ch5Pcmdhbml6YXRpb25zTm90SW5Vc2VFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use partitionKeyDescriptor instead')
const PartitionKey$json = {
  '1': 'PartitionKey',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `PartitionKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List partitionKeyDescriptor = $convert.base64Decode(
    'CgxQYXJ0aXRpb25LZXkSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIWCgR0eXBlGO6g14oBIAEoCV'
    'IEdHlwZQ==');

@$core.Deprecated('Use publicKeyDescriptor instead')
const PublicKey$json = {
  '1': 'PublicKey',
  '2': [
    {'1': 'fingerprint', '3': 502547484, '4': 1, '5': 9, '10': 'fingerprint'},
    {
      '1': 'validityendtime',
      '3': 192996550,
      '4': 1,
      '5': 9,
      '10': 'validityendtime'
    },
    {
      '1': 'validitystarttime',
      '3': 288335101,
      '4': 1,
      '5': 9,
      '10': 'validitystarttime'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 12, '10': 'value'},
  ],
};

/// Descriptor for `PublicKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyDescriptor = $convert.base64Decode(
    'CglQdWJsaWNLZXkSJAoLZmluZ2VycHJpbnQYnIjR7wEgASgJUgtmaW5nZXJwcmludBIrCg92YW'
    'xpZGl0eWVuZHRpbWUYxsmDXCABKAlSD3ZhbGlkaXR5ZW5kdGltZRIwChF2YWxpZGl0eXN0YXJ0'
    'dGltZRj9yb6JASABKAlSEXZhbGlkaXR5c3RhcnR0aW1lEhgKBXZhbHVlGOvyn4oBIAEoDFIFdm'
    'FsdWU=');

@$core.Deprecated('Use putEventConfigurationRequestDescriptor instead')
const PutEventConfigurationRequest$json = {
  '1': 'PutEventConfigurationRequest',
  '2': [
    {
      '1': 'aggregationconfigurations',
      '3': 481530463,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AggregationConfiguration',
      '10': 'aggregationconfigurations'
    },
    {
      '1': 'contextkeyselectors',
      '3': 342040300,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.ContextKeySelector',
      '10': 'contextkeyselectors'
    },
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {
      '1': 'maxeventsize',
      '3': 517627763,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.MaxEventSize',
      '10': 'maxeventsize'
    },
    {'1': 'trailname', '3': 507774985, '4': 1, '5': 9, '10': 'trailname'},
  ],
};

/// Descriptor for `PutEventConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventConfigurationRequestDescriptor = $convert.base64Decode(
    'ChxQdXRFdmVudENvbmZpZ3VyYXRpb25SZXF1ZXN0EmYKGWFnZ3JlZ2F0aW9uY29uZmlndXJhdG'
    'lvbnMY36TO5QEgAygLMiQuY2xvdWR0cmFpbC5BZ2dyZWdhdGlvbkNvbmZpZ3VyYXRpb25SGWFn'
    'Z3JlZ2F0aW9uY29uZmlndXJhdGlvbnMSVAoTY29udGV4dGtleXNlbGVjdG9ycxjsvYyjASADKA'
    'syHi5jbG91ZHRyYWlsLkNvbnRleHRLZXlTZWxlY3RvclITY29udGV4dGtleXNlbGVjdG9ycxIp'
    'Cg5ldmVudGRhdGFzdG9yZRjB251BIAEoCVIOZXZlbnRkYXRhc3RvcmUSQAoMbWF4ZXZlbnRzaX'
    'plGPO+6fYBIAEoDjIYLmNsb3VkdHJhaWwuTWF4RXZlbnRTaXplUgxtYXhldmVudHNpemUSIAoJ'
    'dHJhaWxuYW1lGImQkPIBIAEoCVIJdHJhaWxuYW1l');

@$core.Deprecated('Use putEventConfigurationResponseDescriptor instead')
const PutEventConfigurationResponse$json = {
  '1': 'PutEventConfigurationResponse',
  '2': [
    {
      '1': 'aggregationconfigurations',
      '3': 481530463,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AggregationConfiguration',
      '10': 'aggregationconfigurations'
    },
    {
      '1': 'contextkeyselectors',
      '3': 342040300,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.ContextKeySelector',
      '10': 'contextkeyselectors'
    },
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'maxeventsize',
      '3': 517627763,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.MaxEventSize',
      '10': 'maxeventsize'
    },
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `PutEventConfigurationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventConfigurationResponseDescriptor = $convert.base64Decode(
    'Ch1QdXRFdmVudENvbmZpZ3VyYXRpb25SZXNwb25zZRJmChlhZ2dyZWdhdGlvbmNvbmZpZ3VyYX'
    'Rpb25zGN+kzuUBIAMoCzIkLmNsb3VkdHJhaWwuQWdncmVnYXRpb25Db25maWd1cmF0aW9uUhlh'
    'Z2dyZWdhdGlvbmNvbmZpZ3VyYXRpb25zElQKE2NvbnRleHRrZXlzZWxlY3RvcnMY7L2MowEgAy'
    'gLMh4uY2xvdWR0cmFpbC5Db250ZXh0S2V5U2VsZWN0b3JSE2NvbnRleHRrZXlzZWxlY3RvcnMS'
    'MAoRZXZlbnRkYXRhc3RvcmVhcm4Y6KuXngEgASgJUhFldmVudGRhdGFzdG9yZWFybhJACgxtYX'
    'hldmVudHNpemUY877p9gEgASgOMhguY2xvdWR0cmFpbC5NYXhFdmVudFNpemVSDG1heGV2ZW50'
    'c2l6ZRIdCgh0cmFpbGFybhj3ivASIAEoCVIIdHJhaWxhcm4=');

@$core.Deprecated('Use putEventSelectorsRequestDescriptor instead')
const PutEventSelectorsRequest$json = {
  '1': 'PutEventSelectorsRequest',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'eventselectors',
      '3': 344575328,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.EventSelector',
      '10': 'eventselectors'
    },
    {'1': 'trailname', '3': 507774985, '4': 1, '5': 9, '10': 'trailname'},
  ],
};

/// Descriptor for `PutEventSelectorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventSelectorsRequestDescriptor = $convert.base64Decode(
    'ChhQdXRFdmVudFNlbGVjdG9yc1JlcXVlc3QSXAoWYWR2YW5jZWRldmVudHNlbGVjdG9ycxiyts'
    'gRIAMoCzIhLmNsb3VkdHJhaWwuQWR2YW5jZWRFdmVudFNlbGVjdG9yUhZhZHZhbmNlZGV2ZW50'
    'c2VsZWN0b3JzEkUKDmV2ZW50c2VsZWN0b3JzGOCap6QBIAMoCzIZLmNsb3VkdHJhaWwuRXZlbn'
    'RTZWxlY3RvclIOZXZlbnRzZWxlY3RvcnMSIAoJdHJhaWxuYW1lGImQkPIBIAEoCVIJdHJhaWxu'
    'YW1l');

@$core.Deprecated('Use putEventSelectorsResponseDescriptor instead')
const PutEventSelectorsResponse$json = {
  '1': 'PutEventSelectorsResponse',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'eventselectors',
      '3': 344575328,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.EventSelector',
      '10': 'eventselectors'
    },
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `PutEventSelectorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEventSelectorsResponseDescriptor = $convert.base64Decode(
    'ChlQdXRFdmVudFNlbGVjdG9yc1Jlc3BvbnNlElwKFmFkdmFuY2VkZXZlbnRzZWxlY3RvcnMYsr'
    'bIESADKAsyIS5jbG91ZHRyYWlsLkFkdmFuY2VkRXZlbnRTZWxlY3RvclIWYWR2YW5jZWRldmVu'
    'dHNlbGVjdG9ycxJFCg5ldmVudHNlbGVjdG9ycxjgmqekASADKAsyGS5jbG91ZHRyYWlsLkV2ZW'
    '50U2VsZWN0b3JSDmV2ZW50c2VsZWN0b3JzEh0KCHRyYWlsYXJuGPeK8BIgASgJUgh0cmFpbGFy'
    'bg==');

@$core.Deprecated('Use putInsightSelectorsRequestDescriptor instead')
const PutInsightSelectorsRequest$json = {
  '1': 'PutInsightSelectorsRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {
      '1': 'insightselectors',
      '3': 49628748,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.InsightSelector',
      '10': 'insightselectors'
    },
    {
      '1': 'insightsdestination',
      '3': 396776649,
      '4': 1,
      '5': 9,
      '10': 'insightsdestination'
    },
    {'1': 'trailname', '3': 507774985, '4': 1, '5': 9, '10': 'trailname'},
  ],
};

/// Descriptor for `PutInsightSelectorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putInsightSelectorsRequestDescriptor = $convert.base64Decode(
    'ChpQdXRJbnNpZ2h0U2VsZWN0b3JzUmVxdWVzdBIpCg5ldmVudGRhdGFzdG9yZRjB251BIAEoCV'
    'IOZXZlbnRkYXRhc3RvcmUSSgoQaW5zaWdodHNlbGVjdG9ycxjMjNUXIAMoCzIbLmNsb3VkdHJh'
    'aWwuSW5zaWdodFNlbGVjdG9yUhBpbnNpZ2h0c2VsZWN0b3JzEjQKE2luc2lnaHRzZGVzdGluYX'
    'Rpb24YyamZvQEgASgJUhNpbnNpZ2h0c2Rlc3RpbmF0aW9uEiAKCXRyYWlsbmFtZRiJkJDyASAB'
    'KAlSCXRyYWlsbmFtZQ==');

@$core.Deprecated('Use putInsightSelectorsResponseDescriptor instead')
const PutInsightSelectorsResponse$json = {
  '1': 'PutInsightSelectorsResponse',
  '2': [
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'insightselectors',
      '3': 49628748,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.InsightSelector',
      '10': 'insightselectors'
    },
    {
      '1': 'insightsdestination',
      '3': 396776649,
      '4': 1,
      '5': 9,
      '10': 'insightsdestination'
    },
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `PutInsightSelectorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putInsightSelectorsResponseDescriptor = $convert.base64Decode(
    'ChtQdXRJbnNpZ2h0U2VsZWN0b3JzUmVzcG9uc2USMAoRZXZlbnRkYXRhc3RvcmVhcm4Y6KuXng'
    'EgASgJUhFldmVudGRhdGFzdG9yZWFybhJKChBpbnNpZ2h0c2VsZWN0b3JzGMyM1RcgAygLMhsu'
    'Y2xvdWR0cmFpbC5JbnNpZ2h0U2VsZWN0b3JSEGluc2lnaHRzZWxlY3RvcnMSNAoTaW5zaWdodH'
    'NkZXN0aW5hdGlvbhjJqZm9ASABKAlSE2luc2lnaHRzZGVzdGluYXRpb24SHQoIdHJhaWxhcm4Y'
    '94rwEiABKAlSCHRyYWlsYXJu');

@$core.Deprecated('Use putResourcePolicyRequestDescriptor instead')
const PutResourcePolicyRequest$json = {
  '1': 'PutResourcePolicyRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
  ],
};

/// Descriptor for `PutResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChhQdXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZX'
        'NvdXJjZWFybhIpCg5yZXNvdXJjZXBvbGljeRiwlMEHIAEoCVIOcmVzb3VyY2Vwb2xpY3k=');

@$core.Deprecated('Use putResourcePolicyResponseDescriptor instead')
const PutResourcePolicyResponse$json = {
  '1': 'PutResourcePolicyResponse',
  '2': [
    {
      '1': 'delegatedadminresourcepolicy',
      '3': 510966132,
      '4': 1,
      '5': 9,
      '10': 'delegatedadminresourcepolicy'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'resourcepolicy',
      '3': 15747632,
      '4': 1,
      '5': 9,
      '10': 'resourcepolicy'
    },
  ],
};

/// Descriptor for `PutResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyResponseDescriptor = $convert.base64Decode(
    'ChlQdXRSZXNvdXJjZVBvbGljeVJlc3BvbnNlEkYKHGRlbGVnYXRlZGFkbWlucmVzb3VyY2Vwb2'
    'xpY3kY9PLS8wEgASgJUhxkZWxlZ2F0ZWRhZG1pbnJlc291cmNlcG9saWN5EiQKC3Jlc291cmNl'
    'YXJuGK342a0BIAEoCVILcmVzb3VyY2Vhcm4SKQoOcmVzb3VyY2Vwb2xpY3kYsJTBByABKAlSDn'
    'Jlc291cmNlcG9saWN5');

@$core.Deprecated('Use queryDescriptor instead')
const Query$json = {
  '1': 'Query',
  '2': [
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
    {
      '1': 'querystatus',
      '3': 367016406,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.QueryStatus',
      '10': 'querystatus'
    },
  ],
};

/// Descriptor for `Query`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryDescriptor = $convert.base64Decode(
    'CgVRdWVyeRIlCgxjcmVhdGlvbnRpbWUY5s+qMSABKAlSDGNyZWF0aW9udGltZRIbCgdxdWVyeW'
    'lkGO/w5jQgASgJUgdxdWVyeWlkEj0KC3F1ZXJ5c3RhdHVzGNbzgK8BIAEoDjIXLmNsb3VkdHJh'
    'aWwuUXVlcnlTdGF0dXNSC3F1ZXJ5c3RhdHVz');

@$core.Deprecated('Use queryIdNotFoundExceptionDescriptor instead')
const QueryIdNotFoundException$json = {
  '1': 'QueryIdNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueryIdNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryIdNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChhRdWVyeUlkTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use queryStatisticsDescriptor instead')
const QueryStatistics$json = {
  '1': 'QueryStatistics',
  '2': [
    {'1': 'bytesscanned', '3': 186950329, '4': 1, '5': 3, '10': 'bytesscanned'},
    {'1': 'resultscount', '3': 530057669, '4': 1, '5': 5, '10': 'resultscount'},
    {
      '1': 'totalresultscount',
      '3': 167126221,
      '4': 1,
      '5': 5,
      '10': 'totalresultscount'
    },
  ],
};

/// Descriptor for `QueryStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStatisticsDescriptor = $convert.base64Decode(
    'Cg9RdWVyeVN0YXRpc3RpY3MSJQoMYnl0ZXNzY2FubmVkGLnFklkgASgDUgxieXRlc3NjYW5uZW'
    'QSJgoMcmVzdWx0c2NvdW50GMWT4PwBIAEoBVIMcmVzdWx0c2NvdW50Ei8KEXRvdGFscmVzdWx0'
    'c2NvdW50GM3J2E8gASgFUhF0b3RhbHJlc3VsdHNjb3VudA==');

@$core.Deprecated('Use queryStatisticsForDescribeQueryDescriptor instead')
const QueryStatisticsForDescribeQuery$json = {
  '1': 'QueryStatisticsForDescribeQuery',
  '2': [
    {'1': 'bytesscanned', '3': 186950329, '4': 1, '5': 3, '10': 'bytesscanned'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'eventsmatched',
      '3': 19501829,
      '4': 1,
      '5': 3,
      '10': 'eventsmatched'
    },
    {
      '1': 'eventsscanned',
      '3': 475241657,
      '4': 1,
      '5': 3,
      '10': 'eventsscanned'
    },
    {
      '1': 'executiontimeinmillis',
      '3': 447596178,
      '4': 1,
      '5': 5,
      '10': 'executiontimeinmillis'
    },
  ],
};

/// Descriptor for `QueryStatisticsForDescribeQuery`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStatisticsForDescribeQueryDescriptor = $convert.base64Decode(
    'Ch9RdWVyeVN0YXRpc3RpY3NGb3JEZXNjcmliZVF1ZXJ5EiUKDGJ5dGVzc2Nhbm5lZBi5xZJZIA'
    'EoA1IMYnl0ZXNzY2FubmVkEiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3JlYXRpb250aW1l'
    'EicKDWV2ZW50c21hdGNoZWQYhaamCSABKANSDWV2ZW50c21hdGNoZWQSKAoNZXZlbnRzc2Nhbm'
    '5lZBi5uc7iASABKANSDWV2ZW50c3NjYW5uZWQSOAoVZXhlY3V0aW9udGltZWlubWlsbGlzGJKN'
    't9UBIAEoBVIVZXhlY3V0aW9udGltZWlubWlsbGlz');

@$core.Deprecated('Use refreshScheduleDescriptor instead')
const RefreshSchedule$json = {
  '1': 'RefreshSchedule',
  '2': [
    {
      '1': 'frequency',
      '3': 227673762,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.RefreshScheduleFrequency',
      '10': 'frequency'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.RefreshScheduleStatus',
      '10': 'status'
    },
    {'1': 'timeofday', '3': 77605358, '4': 1, '5': 9, '10': 'timeofday'},
  ],
};

/// Descriptor for `RefreshSchedule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refreshScheduleDescriptor = $convert.base64Decode(
    'Cg9SZWZyZXNoU2NoZWR1bGUSRQoJZnJlcXVlbmN5GKKNyGwgASgLMiQuY2xvdWR0cmFpbC5SZW'
    'ZyZXNoU2NoZWR1bGVGcmVxdWVuY3lSCWZyZXF1ZW5jeRI8CgZzdGF0dXMYkOT7AiABKA4yIS5j'
    'bG91ZHRyYWlsLlJlZnJlc2hTY2hlZHVsZVN0YXR1c1IGc3RhdHVzEh8KCXRpbWVvZmRheRju04'
    'AlIAEoCVIJdGltZW9mZGF5');

@$core.Deprecated('Use refreshScheduleFrequencyDescriptor instead')
const RefreshScheduleFrequency$json = {
  '1': 'RefreshScheduleFrequency',
  '2': [
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.RefreshScheduleFrequencyUnit',
      '10': 'unit'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 5, '10': 'value'},
  ],
};

/// Descriptor for `RefreshScheduleFrequency`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refreshScheduleFrequencyDescriptor = $convert.base64Decode(
    'ChhSZWZyZXNoU2NoZWR1bGVGcmVxdWVuY3kSPwoEdW5pdBiozIVHIAEoDjIoLmNsb3VkdHJhaW'
    'wuUmVmcmVzaFNjaGVkdWxlRnJlcXVlbmN5VW5pdFIEdW5pdBIYCgV2YWx1ZRjr8p+KASABKAVS'
    'BXZhbHVl');

@$core.Deprecated(
    'Use registerOrganizationDelegatedAdminRequestDescriptor instead')
const RegisterOrganizationDelegatedAdminRequest$json = {
  '1': 'RegisterOrganizationDelegatedAdminRequest',
  '2': [
    {
      '1': 'memberaccountid',
      '3': 374644988,
      '4': 1,
      '5': 9,
      '10': 'memberaccountid'
    },
  ],
};

/// Descriptor for `RegisterOrganizationDelegatedAdminRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    registerOrganizationDelegatedAdminRequestDescriptor = $convert.base64Decode(
        'CilSZWdpc3Rlck9yZ2FuaXphdGlvbkRlbGVnYXRlZEFkbWluUmVxdWVzdBIsCg9tZW1iZXJhY2'
        'NvdW50aWQY/MHSsgEgASgJUg9tZW1iZXJhY2NvdW50aWQ=');

@$core.Deprecated(
    'Use registerOrganizationDelegatedAdminResponseDescriptor instead')
const RegisterOrganizationDelegatedAdminResponse$json = {
  '1': 'RegisterOrganizationDelegatedAdminResponse',
};

/// Descriptor for `RegisterOrganizationDelegatedAdminResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    registerOrganizationDelegatedAdminResponseDescriptor =
    $convert.base64Decode(
        'CipSZWdpc3Rlck9yZ2FuaXphdGlvbkRlbGVnYXRlZEFkbWluUmVzcG9uc2U=');

@$core.Deprecated('Use removeTagsRequestDescriptor instead')
const RemoveTagsRequest$json = {
  '1': 'RemoveTagsRequest',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
  ],
};

/// Descriptor for `RemoveTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTagsRequestDescriptor = $convert.base64Decode(
    'ChFSZW1vdmVUYWdzUmVxdWVzdBIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcmVzb3VyY2VpZB'
    'IvCgh0YWdzbGlzdBippJjtASADKAsyDy5jbG91ZHRyYWlsLlRhZ1IIdGFnc2xpc3Q=');

@$core.Deprecated('Use removeTagsResponseDescriptor instead')
const RemoveTagsResponse$json = {
  '1': 'RemoveTagsResponse',
};

/// Descriptor for `RemoveTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTagsResponseDescriptor =
    $convert.base64Decode('ChJSZW1vdmVUYWdzUmVzcG9uc2U=');

@$core.Deprecated('Use requestWidgetDescriptor instead')
const RequestWidget$json = {
  '1': 'RequestWidget',
  '2': [
    {
      '1': 'queryparameters',
      '3': 76317660,
      '4': 3,
      '5': 9,
      '10': 'queryparameters'
    },
    {
      '1': 'querystatement',
      '3': 340852217,
      '4': 1,
      '5': 9,
      '10': 'querystatement'
    },
    {
      '1': 'viewproperties',
      '3': 381974398,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.RequestWidget.ViewpropertiesEntry',
      '10': 'viewproperties'
    },
  ],
  '3': [RequestWidget_ViewpropertiesEntry$json],
};

@$core.Deprecated('Use requestWidgetDescriptor instead')
const RequestWidget_ViewpropertiesEntry$json = {
  '1': 'ViewpropertiesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `RequestWidget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestWidgetDescriptor = $convert.base64Decode(
    'Cg1SZXF1ZXN0V2lkZ2V0EisKD3F1ZXJ5cGFyYW1ldGVycxjch7IkIAMoCVIPcXVlcnlwYXJhbW'
    'V0ZXJzEioKDnF1ZXJ5c3RhdGVtZW50GPn7w6IBIAEoCVIOcXVlcnlzdGF0ZW1lbnQSWQoOdmll'
    'd3Byb3BlcnRpZXMY/u6RtgEgAygLMi0uY2xvdWR0cmFpbC5SZXF1ZXN0V2lkZ2V0LlZpZXdwcm'
    '9wZXJ0aWVzRW50cnlSDnZpZXdwcm9wZXJ0aWVzGkEKE1ZpZXdwcm9wZXJ0aWVzRW50cnkSEAoD'
    'a2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use resourceDescriptor instead')
const Resource$json = {
  '1': 'Resource',
  '2': [
    {'1': 'resourcename', '3': 269834071, '4': 1, '5': 9, '10': 'resourcename'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `Resource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDescriptor = $convert.base64Decode(
    'CghSZXNvdXJjZRImCgxyZXNvdXJjZW5hbWUY167VgAEgASgJUgxyZXNvdXJjZW5hbWUSJgoMcm'
    'Vzb3VyY2V0eXBlGN6+2I8BIAEoCVIMcmVzb3VyY2V0eXBl');

@$core.Deprecated('Use resourceARNNotValidExceptionDescriptor instead')
const ResourceARNNotValidException$json = {
  '1': 'ResourceARNNotValidException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceARNNotValidException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceARNNotValidExceptionDescriptor =
    $convert.base64Decode(
        'ChxSZXNvdXJjZUFSTk5vdFZhbGlkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use resourcePolicyNotFoundExceptionDescriptor instead')
const ResourcePolicyNotFoundException$json = {
  '1': 'ResourcePolicyNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourcePolicyNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcePolicyNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ch9SZXNvdXJjZVBvbGljeU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use resourcePolicyNotValidExceptionDescriptor instead')
const ResourcePolicyNotValidException$json = {
  '1': 'ResourcePolicyNotValidException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourcePolicyNotValidException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcePolicyNotValidExceptionDescriptor =
    $convert.base64Decode(
        'Ch9SZXNvdXJjZVBvbGljeU5vdFZhbGlkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use resourceTagDescriptor instead')
const ResourceTag$json = {
  '1': 'ResourceTag',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'tagslist',
      '3': 497422889,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Tag',
      '10': 'tagslist'
    },
  ],
};

/// Descriptor for `ResourceTag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceTagDescriptor = $convert.base64Decode(
    'CgtSZXNvdXJjZVRhZxIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcmVzb3VyY2VpZBIvCgh0YW'
    'dzbGlzdBippJjtASADKAsyDy5jbG91ZHRyYWlsLlRhZ1IIdGFnc2xpc3Q=');

@$core.Deprecated('Use resourceTypeNotSupportedExceptionDescriptor instead')
const ResourceTypeNotSupportedException$json = {
  '1': 'ResourceTypeNotSupportedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceTypeNotSupportedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceTypeNotSupportedExceptionDescriptor =
    $convert.base64Decode(
        'CiFSZXNvdXJjZVR5cGVOb3RTdXBwb3J0ZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use restoreEventDataStoreRequestDescriptor instead')
const RestoreEventDataStoreRequest$json = {
  '1': 'RestoreEventDataStoreRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
  ],
};

/// Descriptor for `RestoreEventDataStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreEventDataStoreRequestDescriptor =
    $convert.base64Decode(
        'ChxSZXN0b3JlRXZlbnREYXRhU3RvcmVSZXF1ZXN0EikKDmV2ZW50ZGF0YXN0b3JlGMHbnUEgAS'
        'gJUg5ldmVudGRhdGFzdG9yZQ==');

@$core.Deprecated('Use restoreEventDataStoreResponseDescriptor instead')
const RestoreEventDataStoreResponse$json = {
  '1': 'RestoreEventDataStoreResponse',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'multiregionenabled',
      '3': 20620620,
      '4': 1,
      '5': 8,
      '10': 'multiregionenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'organizationenabled',
      '3': 480171176,
      '4': 1,
      '5': 8,
      '10': 'organizationenabled'
    },
    {
      '1': 'retentionperiod',
      '3': 196383721,
      '4': 1,
      '5': 5,
      '10': 'retentionperiod'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.EventDataStoreStatus',
      '10': 'status'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `RestoreEventDataStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreEventDataStoreResponseDescriptor = $convert.base64Decode(
    'Ch1SZXN0b3JlRXZlbnREYXRhU3RvcmVSZXNwb25zZRJcChZhZHZhbmNlZGV2ZW50c2VsZWN0b3'
    'JzGLK2yBEgAygLMiEuY2xvdWR0cmFpbC5BZHZhbmNlZEV2ZW50U2VsZWN0b3JSFmFkdmFuY2Vk'
    'ZXZlbnRzZWxlY3RvcnMSPAoLYmlsbGluZ21vZGUYwLToVyABKA4yFy5jbG91ZHRyYWlsLkJpbG'
    'xpbmdNb2RlUgtiaWxsaW5nbW9kZRIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCVIQY3Jl'
    'YXRlZHRpbWVzdGFtcBIwChFldmVudGRhdGFzdG9yZWFybhjoq5eeASABKAlSEWV2ZW50ZGF0YX'
    'N0b3JlYXJuEh0KCGttc2tleWlkGI3JlxYgASgJUghrbXNrZXlpZBIxChJtdWx0aXJlZ2lvbmVu'
    'YWJsZWQYzMrqCSABKAhSEm11bHRpcmVnaW9uZW5hYmxlZBIVCgRuYW1lGIfmgX8gASgJUgRuYW'
    '1lEjQKE29yZ2FuaXphdGlvbmVuYWJsZWQYqKn75AEgASgIUhNvcmdhbml6YXRpb25lbmFibGVk'
    'EisKD3JldGVudGlvbnBlcmlvZBjpp9JdIAEoBVIPcmV0ZW50aW9ucGVyaW9kEjsKBnN0YXR1cx'
    'iQ5PsCIAEoDjIgLmNsb3VkdHJhaWwuRXZlbnREYXRhU3RvcmVTdGF0dXNSBnN0YXR1cxJGChx0'
    'ZXJtaW5hdGlvbnByb3RlY3Rpb25lbmFibGVkGNzz2bMBIAEoCFIcdGVybWluYXRpb25wcm90ZW'
    'N0aW9uZW5hYmxlZBItChB1cGRhdGVkdGltZXN0YW1wGIHjkxUgASgJUhB1cGRhdGVkdGltZXN0'
    'YW1w');

@$core.Deprecated('Use s3BucketDoesNotExistExceptionDescriptor instead')
const S3BucketDoesNotExistException$json = {
  '1': 'S3BucketDoesNotExistException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `S3BucketDoesNotExistException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3BucketDoesNotExistExceptionDescriptor =
    $convert.base64Decode(
        'Ch1TM0J1Y2tldERvZXNOb3RFeGlzdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use s3ImportSourceDescriptor instead')
const S3ImportSource$json = {
  '1': 'S3ImportSource',
  '2': [
    {
      '1': 's3bucketaccessrolearn',
      '3': 184258167,
      '4': 1,
      '5': 9,
      '10': 's3bucketaccessrolearn'
    },
    {
      '1': 's3bucketregion',
      '3': 223825074,
      '4': 1,
      '5': 9,
      '10': 's3bucketregion'
    },
    {
      '1': 's3locationuri',
      '3': 115913071,
      '4': 1,
      '5': 9,
      '10': 's3locationuri'
    },
  ],
};

/// Descriptor for `S3ImportSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3ImportSourceDescriptor = $convert.base64Decode(
    'Cg5TM0ltcG9ydFNvdXJjZRI3ChVzM2J1Y2tldGFjY2Vzc3JvbGVhcm4Y95zuVyABKAlSFXMzYn'
    'Vja2V0YWNjZXNzcm9sZWFybhIpCg5zM2J1Y2tldHJlZ2lvbhiymd1qIAEoCVIOczNidWNrZXRy'
    'ZWdpb24SJwoNczNsb2NhdGlvbnVyaRjv4qI3IAEoCVINczNsb2NhdGlvbnVyaQ==');

@$core.Deprecated('Use searchSampleQueriesRequestDescriptor instead')
const SearchSampleQueriesRequest$json = {
  '1': 'SearchSampleQueriesRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'searchphrase', '3': 243779077, '4': 1, '5': 9, '10': 'searchphrase'},
  ],
};

/// Descriptor for `SearchSampleQueriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchSampleQueriesRequestDescriptor =
    $convert.base64Decode(
        'ChpTZWFyY2hTYW1wbGVRdWVyaWVzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
        'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIlCgxzZWFyY2hwaHJh'
        'c2UYhYyfdCABKAlSDHNlYXJjaHBocmFzZQ==');

@$core.Deprecated('Use searchSampleQueriesResponseDescriptor instead')
const SearchSampleQueriesResponse$json = {
  '1': 'SearchSampleQueriesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'searchresults',
      '3': 29101186,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.SearchSampleQueriesSearchResult',
      '10': 'searchresults'
    },
  ],
};

/// Descriptor for `SearchSampleQueriesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchSampleQueriesResponseDescriptor =
    $convert.base64Decode(
        'ChtTZWFyY2hTYW1wbGVRdWVyaWVzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SVAoNc2VhcmNocmVzdWx0cxiCmfANIAMoCzIrLmNsb3VkdHJhaWwuU2VhcmNoU2Ft'
        'cGxlUXVlcmllc1NlYXJjaFJlc3VsdFINc2VhcmNocmVzdWx0cw==');

@$core.Deprecated('Use searchSampleQueriesSearchResultDescriptor instead')
const SearchSampleQueriesSearchResult$json = {
  '1': 'SearchSampleQueriesSearchResult',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'relevance', '3': 287814281, '4': 1, '5': 2, '10': 'relevance'},
    {'1': 'sql', '3': 18818720, '4': 1, '5': 9, '10': 'sql'},
  ],
};

/// Descriptor for `SearchSampleQueriesSearchResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchSampleQueriesSearchResultDescriptor =
    $convert.base64Decode(
        'Ch9TZWFyY2hTYW1wbGVRdWVyaWVzU2VhcmNoUmVzdWx0EiMKC2Rlc2NyaXB0aW9uGIr0+TYgAS'
        'gJUgtkZXNjcmlwdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiAKCXJlbGV2YW5jZRiJ5Z6J'
        'ASABKAJSCXJlbGV2YW5jZRITCgNzcWwYoM38CCABKAlSA3NxbA==');

@$core.Deprecated('Use serviceQuotaExceededExceptionDescriptor instead')
const ServiceQuotaExceededException$json = {
  '1': 'ServiceQuotaExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceQuotaExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceQuotaExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1TZXJ2aWNlUXVvdGFFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use sourceConfigDescriptor instead')
const SourceConfig$json = {
  '1': 'SourceConfig',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'applytoallregions',
      '3': 512541819,
      '4': 1,
      '5': 8,
      '10': 'applytoallregions'
    },
  ],
};

/// Descriptor for `SourceConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sourceConfigDescriptor = $convert.base64Decode(
    'CgxTb3VyY2VDb25maWcSXAoWYWR2YW5jZWRldmVudHNlbGVjdG9ycxiytsgRIAMoCzIhLmNsb3'
    'VkdHJhaWwuQWR2YW5jZWRFdmVudFNlbGVjdG9yUhZhZHZhbmNlZGV2ZW50c2VsZWN0b3JzEjAK'
    'EWFwcGx5dG9hbGxyZWdpb25zGPuIs/QBIAEoCFIRYXBwbHl0b2FsbHJlZ2lvbnM=');

@$core.Deprecated('Use startDashboardRefreshRequestDescriptor instead')
const StartDashboardRefreshRequest$json = {
  '1': 'StartDashboardRefreshRequest',
  '2': [
    {'1': 'dashboardid', '3': 430615703, '4': 1, '5': 9, '10': 'dashboardid'},
    {
      '1': 'queryparametervalues',
      '3': 398732929,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.StartDashboardRefreshRequest.QueryparametervaluesEntry',
      '10': 'queryparametervalues'
    },
  ],
  '3': [StartDashboardRefreshRequest_QueryparametervaluesEntry$json],
};

@$core.Deprecated('Use startDashboardRefreshRequestDescriptor instead')
const StartDashboardRefreshRequest_QueryparametervaluesEntry$json = {
  '1': 'QueryparametervaluesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `StartDashboardRefreshRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startDashboardRefreshRequestDescriptor = $convert.base64Decode(
    'ChxTdGFydERhc2hib2FyZFJlZnJlc2hSZXF1ZXN0EiQKC2Rhc2hib2FyZGlkGJfZqs0BIAEoCV'
    'ILZGFzaGJvYXJkaWQSegoUcXVlcnlwYXJhbWV0ZXJ2YWx1ZXMYgd2QvgEgAygLMkIuY2xvdWR0'
    'cmFpbC5TdGFydERhc2hib2FyZFJlZnJlc2hSZXF1ZXN0LlF1ZXJ5cGFyYW1ldGVydmFsdWVzRW'
    '50cnlSFHF1ZXJ5cGFyYW1ldGVydmFsdWVzGkcKGVF1ZXJ5cGFyYW1ldGVydmFsdWVzRW50cnkS'
    'EAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use startDashboardRefreshResponseDescriptor instead')
const StartDashboardRefreshResponse$json = {
  '1': 'StartDashboardRefreshResponse',
  '2': [
    {'1': 'refreshid', '3': 57407610, '4': 1, '5': 9, '10': 'refreshid'},
  ],
};

/// Descriptor for `StartDashboardRefreshResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startDashboardRefreshResponseDescriptor =
    $convert.base64Decode(
        'Ch1TdGFydERhc2hib2FyZFJlZnJlc2hSZXNwb25zZRIfCglyZWZyZXNoaWQY+vCvGyABKAlSCX'
        'JlZnJlc2hpZA==');

@$core.Deprecated('Use startEventDataStoreIngestionRequestDescriptor instead')
const StartEventDataStoreIngestionRequest$json = {
  '1': 'StartEventDataStoreIngestionRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
  ],
};

/// Descriptor for `StartEventDataStoreIngestionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startEventDataStoreIngestionRequestDescriptor =
    $convert.base64Decode(
        'CiNTdGFydEV2ZW50RGF0YVN0b3JlSW5nZXN0aW9uUmVxdWVzdBIpCg5ldmVudGRhdGFzdG9yZR'
        'jB251BIAEoCVIOZXZlbnRkYXRhc3RvcmU=');

@$core.Deprecated('Use startEventDataStoreIngestionResponseDescriptor instead')
const StartEventDataStoreIngestionResponse$json = {
  '1': 'StartEventDataStoreIngestionResponse',
};

/// Descriptor for `StartEventDataStoreIngestionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startEventDataStoreIngestionResponseDescriptor =
    $convert
        .base64Decode('CiRTdGFydEV2ZW50RGF0YVN0b3JlSW5nZXN0aW9uUmVzcG9uc2U=');

@$core.Deprecated('Use startImportRequestDescriptor instead')
const StartImportRequest$json = {
  '1': 'StartImportRequest',
  '2': [
    {'1': 'destinations', '3': 404385861, '4': 3, '5': 9, '10': 'destinations'},
    {'1': 'endeventtime', '3': 260441984, '4': 1, '5': 9, '10': 'endeventtime'},
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importsource',
      '3': 41128754,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.ImportSource',
      '10': 'importsource'
    },
    {
      '1': 'starteventtime',
      '3': 107272573,
      '4': 1,
      '5': 9,
      '10': 'starteventtime'
    },
  ],
};

/// Descriptor for `StartImportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startImportRequestDescriptor = $convert.base64Decode(
    'ChJTdGFydEltcG9ydFJlcXVlc3QSJgoMZGVzdGluYXRpb25zGMXg6cABIAMoCVIMZGVzdGluYX'
    'Rpb25zEiUKDGVuZGV2ZW50dGltZRiAj5h8IAEoCVIMZW5kZXZlbnR0aW1lEh4KCGltcG9ydGlk'
    'GNqUrMgBIAEoCVIIaW1wb3J0aWQSPwoMaW1wb3J0c291cmNlGLKmzhMgASgLMhguY2xvdWR0cm'
    'FpbC5JbXBvcnRTb3VyY2VSDGltcG9ydHNvdXJjZRIpCg5zdGFydGV2ZW50dGltZRj9spMzIAEo'
    'CVIOc3RhcnRldmVudHRpbWU=');

@$core.Deprecated('Use startImportResponseDescriptor instead')
const StartImportResponse$json = {
  '1': 'StartImportResponse',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'destinations', '3': 404385861, '4': 3, '5': 9, '10': 'destinations'},
    {'1': 'endeventtime', '3': 260441984, '4': 1, '5': 9, '10': 'endeventtime'},
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importsource',
      '3': 41128754,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.ImportSource',
      '10': 'importsource'
    },
    {
      '1': 'importstatus',
      '3': 129077631,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'starteventtime',
      '3': 107272573,
      '4': 1,
      '5': 9,
      '10': 'starteventtime'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `StartImportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startImportResponseDescriptor = $convert.base64Decode(
    'ChNTdGFydEltcG9ydFJlc3BvbnNlEi4KEGNyZWF0ZWR0aW1lc3RhbXAY+tvPnwEgASgJUhBjcm'
    'VhdGVkdGltZXN0YW1wEiYKDGRlc3RpbmF0aW9ucxjF4OnAASADKAlSDGRlc3RpbmF0aW9ucxIl'
    'CgxlbmRldmVudHRpbWUYgI+YfCABKAlSDGVuZGV2ZW50dGltZRIeCghpbXBvcnRpZBjalKzIAS'
    'ABKAlSCGltcG9ydGlkEj8KDGltcG9ydHNvdXJjZRiyps4TIAEoCzIYLmNsb3VkdHJhaWwuSW1w'
    'b3J0U291cmNlUgxpbXBvcnRzb3VyY2USPwoMaW1wb3J0c3RhdHVzGP+ixj0gASgOMhguY2xvdW'
    'R0cmFpbC5JbXBvcnRTdGF0dXNSDGltcG9ydHN0YXR1cxIpCg5zdGFydGV2ZW50dGltZRj9spMz'
    'IAEoCVIOc3RhcnRldmVudHRpbWUSLQoQdXBkYXRlZHRpbWVzdGFtcBiB45MVIAEoCVIQdXBkYX'
    'RlZHRpbWVzdGFtcA==');

@$core.Deprecated('Use startLoggingRequestDescriptor instead')
const StartLoggingRequest$json = {
  '1': 'StartLoggingRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `StartLoggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startLoggingRequestDescriptor =
    $convert.base64Decode(
        'ChNTdGFydExvZ2dpbmdSZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use startLoggingResponseDescriptor instead')
const StartLoggingResponse$json = {
  '1': 'StartLoggingResponse',
};

/// Descriptor for `StartLoggingResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startLoggingResponseDescriptor =
    $convert.base64Decode('ChRTdGFydExvZ2dpbmdSZXNwb25zZQ==');

@$core.Deprecated('Use startQueryRequestDescriptor instead')
const StartQueryRequest$json = {
  '1': 'StartQueryRequest',
  '2': [
    {
      '1': 'deliverys3uri',
      '3': 230884460,
      '4': 1,
      '5': 9,
      '10': 'deliverys3uri'
    },
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {'1': 'queryalias', '3': 512762270, '4': 1, '5': 9, '10': 'queryalias'},
    {
      '1': 'queryparameters',
      '3': 76317660,
      '4': 3,
      '5': 9,
      '10': 'queryparameters'
    },
    {
      '1': 'querystatement',
      '3': 340852217,
      '4': 1,
      '5': 9,
      '10': 'querystatement'
    },
  ],
};

/// Descriptor for `StartQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startQueryRequestDescriptor = $convert.base64Decode(
    'ChFTdGFydFF1ZXJ5UmVxdWVzdBInCg1kZWxpdmVyeXMzdXJpGOyIjG4gASgJUg1kZWxpdmVyeX'
    'MzdXJpEkYKHGV2ZW50ZGF0YXN0b3Jlb3duZXJhY2NvdW50aWQYsN3w4AEgASgJUhxldmVudGRh'
    'dGFzdG9yZW93bmVyYWNjb3VudGlkEiIKCnF1ZXJ5YWxpYXMYnsPA9AEgASgJUgpxdWVyeWFsaW'
    'FzEisKD3F1ZXJ5cGFyYW1ldGVycxjch7IkIAMoCVIPcXVlcnlwYXJhbWV0ZXJzEioKDnF1ZXJ5'
    'c3RhdGVtZW50GPn7w6IBIAEoCVIOcXVlcnlzdGF0ZW1lbnQ=');

@$core.Deprecated('Use startQueryResponseDescriptor instead')
const StartQueryResponse$json = {
  '1': 'StartQueryResponse',
  '2': [
    {
      '1': 'eventdatastoreowneraccountid',
      '3': 471609008,
      '4': 1,
      '5': 9,
      '10': 'eventdatastoreowneraccountid'
    },
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `StartQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startQueryResponseDescriptor = $convert.base64Decode(
    'ChJTdGFydFF1ZXJ5UmVzcG9uc2USRgocZXZlbnRkYXRhc3RvcmVvd25lcmFjY291bnRpZBiw3f'
    'DgASABKAlSHGV2ZW50ZGF0YXN0b3Jlb3duZXJhY2NvdW50aWQSGwoHcXVlcnlpZBjv8OY0IAEo'
    'CVIHcXVlcnlpZA==');

@$core.Deprecated('Use stopEventDataStoreIngestionRequestDescriptor instead')
const StopEventDataStoreIngestionRequest$json = {
  '1': 'StopEventDataStoreIngestionRequest',
  '2': [
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
  ],
};

/// Descriptor for `StopEventDataStoreIngestionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopEventDataStoreIngestionRequestDescriptor =
    $convert.base64Decode(
        'CiJTdG9wRXZlbnREYXRhU3RvcmVJbmdlc3Rpb25SZXF1ZXN0EikKDmV2ZW50ZGF0YXN0b3JlGM'
        'HbnUEgASgJUg5ldmVudGRhdGFzdG9yZQ==');

@$core.Deprecated('Use stopEventDataStoreIngestionResponseDescriptor instead')
const StopEventDataStoreIngestionResponse$json = {
  '1': 'StopEventDataStoreIngestionResponse',
};

/// Descriptor for `StopEventDataStoreIngestionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopEventDataStoreIngestionResponseDescriptor =
    $convert
        .base64Decode('CiNTdG9wRXZlbnREYXRhU3RvcmVJbmdlc3Rpb25SZXNwb25zZQ==');

@$core.Deprecated('Use stopImportRequestDescriptor instead')
const StopImportRequest$json = {
  '1': 'StopImportRequest',
  '2': [
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
  ],
};

/// Descriptor for `StopImportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopImportRequestDescriptor = $convert.base64Decode(
    'ChFTdG9wSW1wb3J0UmVxdWVzdBIeCghpbXBvcnRpZBjalKzIASABKAlSCGltcG9ydGlk');

@$core.Deprecated('Use stopImportResponseDescriptor instead')
const StopImportResponse$json = {
  '1': 'StopImportResponse',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'destinations', '3': 404385861, '4': 3, '5': 9, '10': 'destinations'},
    {'1': 'endeventtime', '3': 260441984, '4': 1, '5': 9, '10': 'endeventtime'},
    {'1': 'importid', '3': 420153946, '4': 1, '5': 9, '10': 'importid'},
    {
      '1': 'importsource',
      '3': 41128754,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.ImportSource',
      '10': 'importsource'
    },
    {
      '1': 'importstatistics',
      '3': 48175528,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.ImportStatistics',
      '10': 'importstatistics'
    },
    {
      '1': 'importstatus',
      '3': 129077631,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.ImportStatus',
      '10': 'importstatus'
    },
    {
      '1': 'starteventtime',
      '3': 107272573,
      '4': 1,
      '5': 9,
      '10': 'starteventtime'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `StopImportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopImportResponseDescriptor = $convert.base64Decode(
    'ChJTdG9wSW1wb3J0UmVzcG9uc2USLgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEGNyZW'
    'F0ZWR0aW1lc3RhbXASJgoMZGVzdGluYXRpb25zGMXg6cABIAMoCVIMZGVzdGluYXRpb25zEiUK'
    'DGVuZGV2ZW50dGltZRiAj5h8IAEoCVIMZW5kZXZlbnR0aW1lEh4KCGltcG9ydGlkGNqUrMgBIA'
    'EoCVIIaW1wb3J0aWQSPwoMaW1wb3J0c291cmNlGLKmzhMgASgLMhguY2xvdWR0cmFpbC5JbXBv'
    'cnRTb3VyY2VSDGltcG9ydHNvdXJjZRJLChBpbXBvcnRzdGF0aXN0aWNzGKiz/BYgASgLMhwuY2'
    'xvdWR0cmFpbC5JbXBvcnRTdGF0aXN0aWNzUhBpbXBvcnRzdGF0aXN0aWNzEj8KDGltcG9ydHN0'
    'YXR1cxj/osY9IAEoDjIYLmNsb3VkdHJhaWwuSW1wb3J0U3RhdHVzUgxpbXBvcnRzdGF0dXMSKQ'
    'oOc3RhcnRldmVudHRpbWUY/bKTMyABKAlSDnN0YXJ0ZXZlbnR0aW1lEi0KEHVwZGF0ZWR0aW1l'
    'c3RhbXAYgeOTFSABKAlSEHVwZGF0ZWR0aW1lc3RhbXA=');

@$core.Deprecated('Use stopLoggingRequestDescriptor instead')
const StopLoggingRequest$json = {
  '1': 'StopLoggingRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `StopLoggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopLoggingRequestDescriptor =
    $convert.base64Decode(
        'ChJTdG9wTG9nZ2luZ1JlcXVlc3QSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use stopLoggingResponseDescriptor instead')
const StopLoggingResponse$json = {
  '1': 'StopLoggingResponse',
};

/// Descriptor for `StopLoggingResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopLoggingResponseDescriptor =
    $convert.base64Decode('ChNTdG9wTG9nZ2luZ1Jlc3BvbnNl');

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

@$core.Deprecated('Use tagsLimitExceededExceptionDescriptor instead')
const TagsLimitExceededException$json = {
  '1': 'TagsLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TagsLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagsLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChpUYWdzTGltaXRFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use throttlingExceptionDescriptor instead')
const ThrottlingException$json = {
  '1': 'ThrottlingException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttlingExceptionDescriptor =
    $convert.base64Decode(
        'ChNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use trailDescriptor instead')
const Trail$json = {
  '1': 'Trail',
  '2': [
    {
      '1': 'cloudwatchlogsloggrouparn',
      '3': 519613355,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsloggrouparn'
    },
    {
      '1': 'cloudwatchlogsrolearn',
      '3': 55454690,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsrolearn'
    },
    {
      '1': 'hascustomeventselectors',
      '3': 79311547,
      '4': 1,
      '5': 8,
      '10': 'hascustomeventselectors'
    },
    {
      '1': 'hasinsightselectors',
      '3': 146532978,
      '4': 1,
      '5': 8,
      '10': 'hasinsightselectors'
    },
    {'1': 'homeregion', '3': 300797531, '4': 1, '5': 9, '10': 'homeregion'},
    {
      '1': 'includeglobalserviceevents',
      '3': 212227451,
      '4': 1,
      '5': 8,
      '10': 'includeglobalserviceevents'
    },
    {
      '1': 'ismultiregiontrail',
      '3': 468658571,
      '4': 1,
      '5': 8,
      '10': 'ismultiregiontrail'
    },
    {
      '1': 'isorganizationtrail',
      '3': 145256127,
      '4': 1,
      '5': 8,
      '10': 'isorganizationtrail'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'logfilevalidationenabled',
      '3': 35904346,
      '4': 1,
      '5': 8,
      '10': 'logfilevalidationenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 's3bucketname', '3': 320495427, '4': 1, '5': 9, '10': 's3bucketname'},
    {'1': 's3keyprefix', '3': 206015359, '4': 1, '5': 9, '10': 's3keyprefix'},
    {'1': 'snstopicarn', '3': 380025580, '4': 1, '5': 9, '10': 'snstopicarn'},
    {'1': 'snstopicname', '3': 415454800, '4': 1, '5': 9, '10': 'snstopicname'},
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `Trail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trailDescriptor = $convert.base64Decode(
    'CgVUcmFpbBJAChljbG91ZHdhdGNobG9nc2xvZ2dyb3VwYXJuGKvX4vcBIAEoCVIZY2xvdWR3YX'
    'RjaGxvZ3Nsb2dncm91cGFybhI3ChVjbG91ZHdhdGNobG9nc3JvbGVhcm4Y4te4GiABKAlSFWNs'
    'b3Vkd2F0Y2hsb2dzcm9sZWFybhI7ChdoYXNjdXN0b21ldmVudHNlbGVjdG9ycxi75eglIAEoCF'
    'IXaGFzY3VzdG9tZXZlbnRzZWxlY3RvcnMSMwoTaGFzaW5zaWdodHNlbGVjdG9ycxjy1O9FIAEo'
    'CFITaGFzaW5zaWdodHNlbGVjdG9ycxIiCgpob21lcmVnaW9uGNuct48BIAEoCVIKaG9tZXJlZ2'
    'lvbhJBChppbmNsdWRlZ2xvYmFsc2VydmljZWV2ZW50cxj7qpllIAEoCFIaaW5jbHVkZWdsb2Jh'
    'bHNlcnZpY2VldmVudHMSMgoSaXNtdWx0aXJlZ2lvbnRyYWlsGIvTvN8BIAEoCFISaXNtdWx0aX'
    'JlZ2lvbnRyYWlsEjMKE2lzb3JnYW5pemF0aW9udHJhaWwYv92hRSABKAhSE2lzb3JnYW5pemF0'
    'aW9udHJhaWwSHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEj0KGGxvZ2ZpbGV2YWxpZG'
    'F0aW9uZW5hYmxlZBjato8RIAEoCFIYbG9nZmlsZXZhbGlkYXRpb25lbmFibGVkEhUKBG5hbWUY'
    'h+aBfyABKAlSBG5hbWUSJgoMczNidWNrZXRuYW1lGMO+6ZgBIAEoCVIMczNidWNrZXRuYW1lEi'
    'MKC3Mza2V5cHJlZml4GP+WnmIgASgJUgtzM2tleXByZWZpeBIkCgtzbnN0b3BpY2Fybhjs9Zq1'
    'ASABKAlSC3Nuc3RvcGljYXJuEiYKDHNuc3RvcGljbmFtZRjQrI3GASABKAlSDHNuc3RvcGljbm'
    'FtZRIdCgh0cmFpbGFybhj3ivASIAEoCVIIdHJhaWxhcm4=');

@$core.Deprecated('Use trailAlreadyExistsExceptionDescriptor instead')
const TrailAlreadyExistsException$json = {
  '1': 'TrailAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrailAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trailAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChtUcmFpbEFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use trailInfoDescriptor instead')
const TrailInfo$json = {
  '1': 'TrailInfo',
  '2': [
    {'1': 'homeregion', '3': 300797531, '4': 1, '5': 9, '10': 'homeregion'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `TrailInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trailInfoDescriptor = $convert.base64Decode(
    'CglUcmFpbEluZm8SIgoKaG9tZXJlZ2lvbhjbnLePASABKAlSCmhvbWVyZWdpb24SFQoEbmFtZR'
    'iH5oF/IAEoCVIEbmFtZRIdCgh0cmFpbGFybhj3ivASIAEoCVIIdHJhaWxhcm4=');

@$core.Deprecated('Use trailNotFoundExceptionDescriptor instead')
const TrailNotFoundException$json = {
  '1': 'TrailNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrailNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trailNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChZUcmFpbE5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use trailNotProvidedExceptionDescriptor instead')
const TrailNotProvidedException$json = {
  '1': 'TrailNotProvidedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TrailNotProvidedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trailNotProvidedExceptionDescriptor =
    $convert.base64Decode(
        'ChlUcmFpbE5vdFByb3ZpZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use unsupportedOperationExceptionDescriptor instead')
const UnsupportedOperationException$json = {
  '1': 'UnsupportedOperationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedOperationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedOperationExceptionDescriptor =
    $convert.base64Decode(
        'Ch1VbnN1cHBvcnRlZE9wZXJhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use updateChannelRequestDescriptor instead')
const UpdateChannelRequest$json = {
  '1': 'UpdateChannelRequest',
  '2': [
    {'1': 'channel', '3': 90016325, '4': 1, '5': 9, '10': 'channel'},
    {
      '1': 'destinations',
      '3': 404385861,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Destination',
      '10': 'destinations'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateChannelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateChannelRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVDaGFubmVsUmVxdWVzdBIbCgdjaGFubmVsGMWU9iogASgJUgdjaGFubmVsEj8KDG'
    'Rlc3RpbmF0aW9ucxjF4OnAASADKAsyFy5jbG91ZHRyYWlsLkRlc3RpbmF0aW9uUgxkZXN0aW5h'
    'dGlvbnMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use updateChannelResponseDescriptor instead')
const UpdateChannelResponse$json = {
  '1': 'UpdateChannelResponse',
  '2': [
    {'1': 'channelarn', '3': 99276476, '4': 1, '5': 9, '10': 'channelarn'},
    {
      '1': 'destinations',
      '3': 404385861,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Destination',
      '10': 'destinations'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `UpdateChannelResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateChannelResponseDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVDaGFubmVsUmVzcG9uc2USIQoKY2hhbm5lbGFybhi8rasvIAEoCVIKY2hhbm5lbG'
    'FybhI/CgxkZXN0aW5hdGlvbnMYxeDpwAEgAygLMhcuY2xvdWR0cmFpbC5EZXN0aW5hdGlvblIM'
    'ZGVzdGluYXRpb25zEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGQoGc291cmNlGPnHig8gASgJUg'
    'Zzb3VyY2U=');

@$core.Deprecated('Use updateDashboardRequestDescriptor instead')
const UpdateDashboardRequest$json = {
  '1': 'UpdateDashboardRequest',
  '2': [
    {'1': 'dashboardid', '3': 430615703, '4': 1, '5': 9, '10': 'dashboardid'},
    {
      '1': 'refreshschedule',
      '3': 261773338,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.RefreshSchedule',
      '10': 'refreshschedule'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'widgets',
      '3': 501826147,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.RequestWidget',
      '10': 'widgets'
    },
  ],
};

/// Descriptor for `UpdateDashboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDashboardRequestDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVEYXNoYm9hcmRSZXF1ZXN0EiQKC2Rhc2hib2FyZGlkGJfZqs0BIAEoCVILZGFzaG'
    'JvYXJkaWQSSAoPcmVmcmVzaHNjaGVkdWxlGJqw6XwgASgLMhsuY2xvdWR0cmFpbC5SZWZyZXNo'
    'U2NoZWR1bGVSD3JlZnJlc2hzY2hlZHVsZRJGChx0ZXJtaW5hdGlvbnByb3RlY3Rpb25lbmFibG'
    'VkGNzz2bMBIAEoCFIcdGVybWluYXRpb25wcm90ZWN0aW9uZW5hYmxlZBI3Cgd3aWRnZXRzGOOE'
    'pe8BIAMoCzIZLmNsb3VkdHJhaWwuUmVxdWVzdFdpZGdldFIHd2lkZ2V0cw==');

@$core.Deprecated('Use updateDashboardResponseDescriptor instead')
const UpdateDashboardResponse$json = {
  '1': 'UpdateDashboardResponse',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'dashboardarn', '3': 108051951, '4': 1, '5': 9, '10': 'dashboardarn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'refreshschedule',
      '3': 261773338,
      '4': 1,
      '5': 11,
      '6': '.cloudtrail.RefreshSchedule',
      '10': 'refreshschedule'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.DashboardType',
      '10': 'type'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
    {
      '1': 'widgets',
      '3': 501826147,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Widget',
      '10': 'widgets'
    },
  ],
};

/// Descriptor for `UpdateDashboardResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDashboardResponseDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVEYXNoYm9hcmRSZXNwb25zZRIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCV'
    'IQY3JlYXRlZHRpbWVzdGFtcBIlCgxkYXNoYm9hcmRhcm4Y7/vCMyABKAlSDGRhc2hib2FyZGFy'
    'bhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkgKD3JlZnJlc2hzY2hlZHVsZRiasOl8IAEoCzIbLm'
    'Nsb3VkdHJhaWwuUmVmcmVzaFNjaGVkdWxlUg9yZWZyZXNoc2NoZWR1bGUSRgocdGVybWluYXRp'
    'b25wcm90ZWN0aW9uZW5hYmxlZBjc89mzASABKAhSHHRlcm1pbmF0aW9ucHJvdGVjdGlvbmVuYW'
    'JsZWQSMQoEdHlwZRjuoNeKASABKA4yGS5jbG91ZHRyYWlsLkRhc2hib2FyZFR5cGVSBHR5cGUS'
    'LQoQdXBkYXRlZHRpbWVzdGFtcBiB45MVIAEoCVIQdXBkYXRlZHRpbWVzdGFtcBIwCgd3aWRnZX'
    'RzGOOEpe8BIAMoCzISLmNsb3VkdHJhaWwuV2lkZ2V0Ugd3aWRnZXRz');

@$core.Deprecated('Use updateEventDataStoreRequestDescriptor instead')
const UpdateEventDataStoreRequest$json = {
  '1': 'UpdateEventDataStoreRequest',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'eventdatastore',
      '3': 136801729,
      '4': 1,
      '5': 9,
      '10': 'eventdatastore'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'multiregionenabled',
      '3': 20620620,
      '4': 1,
      '5': 8,
      '10': 'multiregionenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'organizationenabled',
      '3': 480171176,
      '4': 1,
      '5': 8,
      '10': 'organizationenabled'
    },
    {
      '1': 'retentionperiod',
      '3': 196383721,
      '4': 1,
      '5': 5,
      '10': 'retentionperiod'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
  ],
};

/// Descriptor for `UpdateEventDataStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateEventDataStoreRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVFdmVudERhdGFTdG9yZVJlcXVlc3QSXAoWYWR2YW5jZWRldmVudHNlbGVjdG9ycx'
    'iytsgRIAMoCzIhLmNsb3VkdHJhaWwuQWR2YW5jZWRFdmVudFNlbGVjdG9yUhZhZHZhbmNlZGV2'
    'ZW50c2VsZWN0b3JzEjwKC2JpbGxpbmdtb2RlGMC06FcgASgOMhcuY2xvdWR0cmFpbC5CaWxsaW'
    '5nTW9kZVILYmlsbGluZ21vZGUSKQoOZXZlbnRkYXRhc3RvcmUYwdudQSABKAlSDmV2ZW50ZGF0'
    'YXN0b3JlEh0KCGttc2tleWlkGI3JlxYgASgJUghrbXNrZXlpZBIxChJtdWx0aXJlZ2lvbmVuYW'
    'JsZWQYzMrqCSABKAhSEm11bHRpcmVnaW9uZW5hYmxlZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l'
    'EjQKE29yZ2FuaXphdGlvbmVuYWJsZWQYqKn75AEgASgIUhNvcmdhbml6YXRpb25lbmFibGVkEi'
    'sKD3JldGVudGlvbnBlcmlvZBjpp9JdIAEoBVIPcmV0ZW50aW9ucGVyaW9kEkYKHHRlcm1pbmF0'
    'aW9ucHJvdGVjdGlvbmVuYWJsZWQY3PPZswEgASgIUhx0ZXJtaW5hdGlvbnByb3RlY3Rpb25lbm'
    'FibGVk');

@$core.Deprecated('Use updateEventDataStoreResponseDescriptor instead')
const UpdateEventDataStoreResponse$json = {
  '1': 'UpdateEventDataStoreResponse',
  '2': [
    {
      '1': 'advancedeventselectors',
      '3': 36838194,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.AdvancedEventSelector',
      '10': 'advancedeventselectors'
    },
    {
      '1': 'billingmode',
      '3': 184162880,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.BillingMode',
      '10': 'billingmode'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'eventdatastorearn',
      '3': 331732456,
      '4': 1,
      '5': 9,
      '10': 'eventdatastorearn'
    },
    {
      '1': 'federationrolearn',
      '3': 504464364,
      '4': 1,
      '5': 9,
      '10': 'federationrolearn'
    },
    {
      '1': 'federationstatus',
      '3': 146235383,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.FederationStatus',
      '10': 'federationstatus'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'multiregionenabled',
      '3': 20620620,
      '4': 1,
      '5': 8,
      '10': 'multiregionenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'organizationenabled',
      '3': 480171176,
      '4': 1,
      '5': 8,
      '10': 'organizationenabled'
    },
    {
      '1': 'retentionperiod',
      '3': 196383721,
      '4': 1,
      '5': 5,
      '10': 'retentionperiod'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.cloudtrail.EventDataStoreStatus',
      '10': 'status'
    },
    {
      '1': 'terminationprotectionenabled',
      '3': 376863196,
      '4': 1,
      '5': 8,
      '10': 'terminationprotectionenabled'
    },
    {
      '1': 'updatedtimestamp',
      '3': 44364161,
      '4': 1,
      '5': 9,
      '10': 'updatedtimestamp'
    },
  ],
};

/// Descriptor for `UpdateEventDataStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateEventDataStoreResponseDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVFdmVudERhdGFTdG9yZVJlc3BvbnNlElwKFmFkdmFuY2VkZXZlbnRzZWxlY3Rvcn'
    'MYsrbIESADKAsyIS5jbG91ZHRyYWlsLkFkdmFuY2VkRXZlbnRTZWxlY3RvclIWYWR2YW5jZWRl'
    'dmVudHNlbGVjdG9ycxI8CgtiaWxsaW5nbW9kZRjAtOhXIAEoDjIXLmNsb3VkdHJhaWwuQmlsbG'
    'luZ01vZGVSC2JpbGxpbmdtb2RlEi4KEGNyZWF0ZWR0aW1lc3RhbXAY+tvPnwEgASgJUhBjcmVh'
    'dGVkdGltZXN0YW1wEjAKEWV2ZW50ZGF0YXN0b3JlYXJuGOirl54BIAEoCVIRZXZlbnRkYXRhc3'
    'RvcmVhcm4SMAoRZmVkZXJhdGlvbnJvbGVhcm4Y7IfG8AEgASgJUhFmZWRlcmF0aW9ucm9sZWFy'
    'bhJLChBmZWRlcmF0aW9uc3RhdHVzGPe/3UUgASgOMhwuY2xvdWR0cmFpbC5GZWRlcmF0aW9uU3'
    'RhdHVzUhBmZWRlcmF0aW9uc3RhdHVzEh0KCGttc2tleWlkGI3JlxYgASgJUghrbXNrZXlpZBIx'
    'ChJtdWx0aXJlZ2lvbmVuYWJsZWQYzMrqCSABKAhSEm11bHRpcmVnaW9uZW5hYmxlZBIVCgRuYW'
    '1lGIfmgX8gASgJUgRuYW1lEjQKE29yZ2FuaXphdGlvbmVuYWJsZWQYqKn75AEgASgIUhNvcmdh'
    'bml6YXRpb25lbmFibGVkEisKD3JldGVudGlvbnBlcmlvZBjpp9JdIAEoBVIPcmV0ZW50aW9ucG'
    'VyaW9kEjsKBnN0YXR1cxiQ5PsCIAEoDjIgLmNsb3VkdHJhaWwuRXZlbnREYXRhU3RvcmVTdGF0'
    'dXNSBnN0YXR1cxJGChx0ZXJtaW5hdGlvbnByb3RlY3Rpb25lbmFibGVkGNzz2bMBIAEoCFIcdG'
    'VybWluYXRpb25wcm90ZWN0aW9uZW5hYmxlZBItChB1cGRhdGVkdGltZXN0YW1wGIHjkxUgASgJ'
    'UhB1cGRhdGVkdGltZXN0YW1w');

@$core.Deprecated('Use updateTrailRequestDescriptor instead')
const UpdateTrailRequest$json = {
  '1': 'UpdateTrailRequest',
  '2': [
    {
      '1': 'cloudwatchlogsloggrouparn',
      '3': 519613355,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsloggrouparn'
    },
    {
      '1': 'cloudwatchlogsrolearn',
      '3': 55454690,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsrolearn'
    },
    {
      '1': 'enablelogfilevalidation',
      '3': 80236930,
      '4': 1,
      '5': 8,
      '10': 'enablelogfilevalidation'
    },
    {
      '1': 'includeglobalserviceevents',
      '3': 212227451,
      '4': 1,
      '5': 8,
      '10': 'includeglobalserviceevents'
    },
    {
      '1': 'ismultiregiontrail',
      '3': 468658571,
      '4': 1,
      '5': 8,
      '10': 'ismultiregiontrail'
    },
    {
      '1': 'isorganizationtrail',
      '3': 145256127,
      '4': 1,
      '5': 8,
      '10': 'isorganizationtrail'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 's3bucketname', '3': 320495427, '4': 1, '5': 9, '10': 's3bucketname'},
    {'1': 's3keyprefix', '3': 206015359, '4': 1, '5': 9, '10': 's3keyprefix'},
    {'1': 'snstopicname', '3': 415454800, '4': 1, '5': 9, '10': 'snstopicname'},
  ],
};

/// Descriptor for `UpdateTrailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrailRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVUcmFpbFJlcXVlc3QSQAoZY2xvdWR3YXRjaGxvZ3Nsb2dncm91cGFybhir1+L3AS'
    'ABKAlSGWNsb3Vkd2F0Y2hsb2dzbG9nZ3JvdXBhcm4SNwoVY2xvdWR3YXRjaGxvZ3Nyb2xlYXJu'
    'GOLXuBogASgJUhVjbG91ZHdhdGNobG9nc3JvbGVhcm4SOwoXZW5hYmxlbG9nZmlsZXZhbGlkYX'
    'Rpb24YgqOhJiABKAhSF2VuYWJsZWxvZ2ZpbGV2YWxpZGF0aW9uEkEKGmluY2x1ZGVnbG9iYWxz'
    'ZXJ2aWNlZXZlbnRzGPuqmWUgASgIUhppbmNsdWRlZ2xvYmFsc2VydmljZWV2ZW50cxIyChJpc2'
    '11bHRpcmVnaW9udHJhaWwYi9O83wEgASgIUhJpc211bHRpcmVnaW9udHJhaWwSMwoTaXNvcmdh'
    'bml6YXRpb250cmFpbBi/3aFFIAEoCFITaXNvcmdhbml6YXRpb250cmFpbBIdCghrbXNrZXlpZB'
    'iNyZcWIAEoCVIIa21za2V5aWQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRImCgxzM2J1Y2tldG5h'
    'bWUYw77pmAEgASgJUgxzM2J1Y2tldG5hbWUSIwoLczNrZXlwcmVmaXgY/5aeYiABKAlSC3Mza2'
    'V5cHJlZml4EiYKDHNuc3RvcGljbmFtZRjQrI3GASABKAlSDHNuc3RvcGljbmFtZQ==');

@$core.Deprecated('Use updateTrailResponseDescriptor instead')
const UpdateTrailResponse$json = {
  '1': 'UpdateTrailResponse',
  '2': [
    {
      '1': 'cloudwatchlogsloggrouparn',
      '3': 519613355,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsloggrouparn'
    },
    {
      '1': 'cloudwatchlogsrolearn',
      '3': 55454690,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchlogsrolearn'
    },
    {
      '1': 'includeglobalserviceevents',
      '3': 212227451,
      '4': 1,
      '5': 8,
      '10': 'includeglobalserviceevents'
    },
    {
      '1': 'ismultiregiontrail',
      '3': 468658571,
      '4': 1,
      '5': 8,
      '10': 'ismultiregiontrail'
    },
    {
      '1': 'isorganizationtrail',
      '3': 145256127,
      '4': 1,
      '5': 8,
      '10': 'isorganizationtrail'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'logfilevalidationenabled',
      '3': 35904346,
      '4': 1,
      '5': 8,
      '10': 'logfilevalidationenabled'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 's3bucketname', '3': 320495427, '4': 1, '5': 9, '10': 's3bucketname'},
    {'1': 's3keyprefix', '3': 206015359, '4': 1, '5': 9, '10': 's3keyprefix'},
    {'1': 'snstopicarn', '3': 380025580, '4': 1, '5': 9, '10': 'snstopicarn'},
    {'1': 'snstopicname', '3': 415454800, '4': 1, '5': 9, '10': 'snstopicname'},
    {'1': 'trailarn', '3': 39585143, '4': 1, '5': 9, '10': 'trailarn'},
  ],
};

/// Descriptor for `UpdateTrailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTrailResponseDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVUcmFpbFJlc3BvbnNlEkAKGWNsb3Vkd2F0Y2hsb2dzbG9nZ3JvdXBhcm4Yq9fi9w'
    'EgASgJUhljbG91ZHdhdGNobG9nc2xvZ2dyb3VwYXJuEjcKFWNsb3Vkd2F0Y2hsb2dzcm9sZWFy'
    'bhji17gaIAEoCVIVY2xvdWR3YXRjaGxvZ3Nyb2xlYXJuEkEKGmluY2x1ZGVnbG9iYWxzZXJ2aW'
    'NlZXZlbnRzGPuqmWUgASgIUhppbmNsdWRlZ2xvYmFsc2VydmljZWV2ZW50cxIyChJpc211bHRp'
    'cmVnaW9udHJhaWwYi9O83wEgASgIUhJpc211bHRpcmVnaW9udHJhaWwSMwoTaXNvcmdhbml6YX'
    'Rpb250cmFpbBi/3aFFIAEoCFITaXNvcmdhbml6YXRpb250cmFpbBIdCghrbXNrZXlpZBiNyZcW'
    'IAEoCVIIa21za2V5aWQSPQoYbG9nZmlsZXZhbGlkYXRpb25lbmFibGVkGNq2jxEgASgIUhhsb2'
    'dmaWxldmFsaWRhdGlvbmVuYWJsZWQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRImCgxzM2J1Y2tl'
    'dG5hbWUYw77pmAEgASgJUgxzM2J1Y2tldG5hbWUSIwoLczNrZXlwcmVmaXgY/5aeYiABKAlSC3'
    'Mza2V5cHJlZml4EiQKC3Nuc3RvcGljYXJuGOz1mrUBIAEoCVILc25zdG9waWNhcm4SJgoMc25z'
    'dG9waWNuYW1lGNCsjcYBIAEoCVIMc25zdG9waWNuYW1lEh0KCHRyYWlsYXJuGPeK8BIgASgJUg'
    'h0cmFpbGFybg==');

@$core.Deprecated('Use widgetDescriptor instead')
const Widget$json = {
  '1': 'Widget',
  '2': [
    {'1': 'queryalias', '3': 512762270, '4': 1, '5': 9, '10': 'queryalias'},
    {
      '1': 'queryparameters',
      '3': 76317660,
      '4': 3,
      '5': 9,
      '10': 'queryparameters'
    },
    {
      '1': 'querystatement',
      '3': 340852217,
      '4': 1,
      '5': 9,
      '10': 'querystatement'
    },
    {
      '1': 'viewproperties',
      '3': 381974398,
      '4': 3,
      '5': 11,
      '6': '.cloudtrail.Widget.ViewpropertiesEntry',
      '10': 'viewproperties'
    },
  ],
  '3': [Widget_ViewpropertiesEntry$json],
};

@$core.Deprecated('Use widgetDescriptor instead')
const Widget_ViewpropertiesEntry$json = {
  '1': 'ViewpropertiesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Widget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List widgetDescriptor = $convert.base64Decode(
    'CgZXaWRnZXQSIgoKcXVlcnlhbGlhcxiew8D0ASABKAlSCnF1ZXJ5YWxpYXMSKwoPcXVlcnlwYX'
    'JhbWV0ZXJzGNyHsiQgAygJUg9xdWVyeXBhcmFtZXRlcnMSKgoOcXVlcnlzdGF0ZW1lbnQY+fvD'
    'ogEgASgJUg5xdWVyeXN0YXRlbWVudBJSCg52aWV3cHJvcGVydGllcxj+7pG2ASADKAsyJi5jbG'
    '91ZHRyYWlsLldpZGdldC5WaWV3cHJvcGVydGllc0VudHJ5Ug52aWV3cHJvcGVydGllcxpBChNW'
    'aWV3cHJvcGVydGllc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YW'
    'x1ZToCOAE=');
