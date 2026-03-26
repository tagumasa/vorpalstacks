// This is a generated file - do not edit.
//
// Generated from lambda.proto.

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

@$core.Deprecated('Use applicationLogLevelDescriptor instead')
const ApplicationLogLevel$json = {
  '1': 'ApplicationLogLevel',
  '2': [
    {'1': 'APPLICATION_LOG_LEVEL_WARN', '2': 0},
    {'1': 'APPLICATION_LOG_LEVEL_FATAL', '2': 1},
    {'1': 'APPLICATION_LOG_LEVEL_ERROR', '2': 2},
    {'1': 'APPLICATION_LOG_LEVEL_TRACE', '2': 3},
    {'1': 'APPLICATION_LOG_LEVEL_INFO', '2': 4},
    {'1': 'APPLICATION_LOG_LEVEL_DEBUG', '2': 5},
  ],
};

/// Descriptor for `ApplicationLogLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List applicationLogLevelDescriptor = $convert.base64Decode(
    'ChNBcHBsaWNhdGlvbkxvZ0xldmVsEh4KGkFQUExJQ0FUSU9OX0xPR19MRVZFTF9XQVJOEAASHw'
    'obQVBQTElDQVRJT05fTE9HX0xFVkVMX0ZBVEFMEAESHwobQVBQTElDQVRJT05fTE9HX0xFVkVM'
    'X0VSUk9SEAISHwobQVBQTElDQVRJT05fTE9HX0xFVkVMX1RSQUNFEAMSHgoaQVBQTElDQVRJT0'
    '5fTE9HX0xFVkVMX0lORk8QBBIfChtBUFBMSUNBVElPTl9MT0dfTEVWRUxfREVCVUcQBQ==');

@$core.Deprecated('Use architectureDescriptor instead')
const Architecture$json = {
  '1': 'Architecture',
  '2': [
    {'1': 'ARCHITECTURE_ARM64', '2': 0},
    {'1': 'ARCHITECTURE_X86_64', '2': 1},
  ],
};

/// Descriptor for `Architecture`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List architectureDescriptor = $convert.base64Decode(
    'CgxBcmNoaXRlY3R1cmUSFgoSQVJDSElURUNUVVJFX0FSTTY0EAASFwoTQVJDSElURUNUVVJFX1'
    'g4Nl82NBAB');

@$core.Deprecated('Use capacityProviderPredefinedMetricTypeDescriptor instead')
const CapacityProviderPredefinedMetricType$json = {
  '1': 'CapacityProviderPredefinedMetricType',
  '2': [
    {
      '1':
          'CAPACITY_PROVIDER_PREDEFINED_METRIC_TYPE_LAMBDACAPACITYPROVIDERAVERAGECPUUTILIZATION',
      '2': 0
    },
  ],
};

/// Descriptor for `CapacityProviderPredefinedMetricType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List capacityProviderPredefinedMetricTypeDescriptor =
    $convert.base64Decode(
        'CiRDYXBhY2l0eVByb3ZpZGVyUHJlZGVmaW5lZE1ldHJpY1R5cGUSWApUQ0FQQUNJVFlfUFJPVk'
        'lERVJfUFJFREVGSU5FRF9NRVRSSUNfVFlQRV9MQU1CREFDQVBBQ0lUWVBST1ZJREVSQVZFUkFH'
        'RUNQVVVUSUxJWkFUSU9OEAA=');

@$core.Deprecated('Use capacityProviderScalingModeDescriptor instead')
const CapacityProviderScalingMode$json = {
  '1': 'CapacityProviderScalingMode',
  '2': [
    {'1': 'CAPACITY_PROVIDER_SCALING_MODE_MANUAL', '2': 0},
    {'1': 'CAPACITY_PROVIDER_SCALING_MODE_AUTO', '2': 1},
  ],
};

/// Descriptor for `CapacityProviderScalingMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List capacityProviderScalingModeDescriptor =
    $convert.base64Decode(
        'ChtDYXBhY2l0eVByb3ZpZGVyU2NhbGluZ01vZGUSKQolQ0FQQUNJVFlfUFJPVklERVJfU0NBTE'
        'lOR19NT0RFX01BTlVBTBAAEicKI0NBUEFDSVRZX1BST1ZJREVSX1NDQUxJTkdfTU9ERV9BVVRP'
        'EAE=');

@$core.Deprecated('Use capacityProviderStateDescriptor instead')
const CapacityProviderState$json = {
  '1': 'CapacityProviderState',
  '2': [
    {'1': 'CAPACITY_PROVIDER_STATE_ACTIVE', '2': 0},
    {'1': 'CAPACITY_PROVIDER_STATE_FAILED', '2': 1},
    {'1': 'CAPACITY_PROVIDER_STATE_DELETING', '2': 2},
    {'1': 'CAPACITY_PROVIDER_STATE_PENDING', '2': 3},
  ],
};

/// Descriptor for `CapacityProviderState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List capacityProviderStateDescriptor = $convert.base64Decode(
    'ChVDYXBhY2l0eVByb3ZpZGVyU3RhdGUSIgoeQ0FQQUNJVFlfUFJPVklERVJfU1RBVEVfQUNUSV'
    'ZFEAASIgoeQ0FQQUNJVFlfUFJPVklERVJfU1RBVEVfRkFJTEVEEAESJAogQ0FQQUNJVFlfUFJP'
    'VklERVJfU1RBVEVfREVMRVRJTkcQAhIjCh9DQVBBQ0lUWV9QUk9WSURFUl9TVEFURV9QRU5ESU'
    '5HEAM=');

@$core.Deprecated('Use codeSigningPolicyDescriptor instead')
const CodeSigningPolicy$json = {
  '1': 'CodeSigningPolicy',
  '2': [
    {'1': 'CODE_SIGNING_POLICY_WARN', '2': 0},
    {'1': 'CODE_SIGNING_POLICY_ENFORCE', '2': 1},
  ],
};

/// Descriptor for `CodeSigningPolicy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List codeSigningPolicyDescriptor = $convert.base64Decode(
    'ChFDb2RlU2lnbmluZ1BvbGljeRIcChhDT0RFX1NJR05JTkdfUE9MSUNZX1dBUk4QABIfChtDT0'
    'RFX1NJR05JTkdfUE9MSUNZX0VORk9SQ0UQAQ==');

@$core.Deprecated('Use endPointTypeDescriptor instead')
const EndPointType$json = {
  '1': 'EndPointType',
  '2': [
    {'1': 'END_POINT_TYPE_KAFKA_BOOTSTRAP_SERVERS', '2': 0},
  ],
};

/// Descriptor for `EndPointType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List endPointTypeDescriptor = $convert.base64Decode(
    'CgxFbmRQb2ludFR5cGUSKgomRU5EX1BPSU5UX1RZUEVfS0FGS0FfQk9PVFNUUkFQX1NFUlZFUl'
    'MQAA==');

@$core.Deprecated('Use eventSourceMappingMetricDescriptor instead')
const EventSourceMappingMetric$json = {
  '1': 'EventSourceMappingMetric',
  '2': [
    {'1': 'EVENT_SOURCE_MAPPING_METRIC_EVENTCOUNT', '2': 0},
    {'1': 'EVENT_SOURCE_MAPPING_METRIC_KAFKAMETRICS', '2': 1},
    {'1': 'EVENT_SOURCE_MAPPING_METRIC_ERRORCOUNT', '2': 2},
  ],
};

/// Descriptor for `EventSourceMappingMetric`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventSourceMappingMetricDescriptor = $convert.base64Decode(
    'ChhFdmVudFNvdXJjZU1hcHBpbmdNZXRyaWMSKgomRVZFTlRfU09VUkNFX01BUFBJTkdfTUVUUk'
    'lDX0VWRU5UQ09VTlQQABIsCihFVkVOVF9TT1VSQ0VfTUFQUElOR19NRVRSSUNfS0FGS0FNRVRS'
    'SUNTEAESKgomRVZFTlRfU09VUkNFX01BUFBJTkdfTUVUUklDX0VSUk9SQ09VTlQQAg==');

@$core.Deprecated('Use eventSourceMappingSystemLogLevelDescriptor instead')
const EventSourceMappingSystemLogLevel$json = {
  '1': 'EventSourceMappingSystemLogLevel',
  '2': [
    {'1': 'EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_WARN', '2': 0},
    {'1': 'EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_INFO', '2': 1},
    {'1': 'EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_DEBUG', '2': 2},
  ],
};

/// Descriptor for `EventSourceMappingSystemLogLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventSourceMappingSystemLogLevelDescriptor =
    $convert.base64Decode(
        'CiBFdmVudFNvdXJjZU1hcHBpbmdTeXN0ZW1Mb2dMZXZlbBIuCipFVkVOVF9TT1VSQ0VfTUFQUE'
        'lOR19TWVNURU1fTE9HX0xFVkVMX1dBUk4QABIuCipFVkVOVF9TT1VSQ0VfTUFQUElOR19TWVNU'
        'RU1fTE9HX0xFVkVMX0lORk8QARIvCitFVkVOVF9TT1VSQ0VfTUFQUElOR19TWVNURU1fTE9HX0'
        'xFVkVMX0RFQlVHEAI=');

@$core.Deprecated('Use eventSourcePositionDescriptor instead')
const EventSourcePosition$json = {
  '1': 'EventSourcePosition',
  '2': [
    {'1': 'EVENT_SOURCE_POSITION_AT_TIMESTAMP', '2': 0},
    {'1': 'EVENT_SOURCE_POSITION_TRIM_HORIZON', '2': 1},
    {'1': 'EVENT_SOURCE_POSITION_LATEST', '2': 2},
  ],
};

/// Descriptor for `EventSourcePosition`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventSourcePositionDescriptor = $convert.base64Decode(
    'ChNFdmVudFNvdXJjZVBvc2l0aW9uEiYKIkVWRU5UX1NPVVJDRV9QT1NJVElPTl9BVF9USU1FU1'
    'RBTVAQABImCiJFVkVOVF9TT1VSQ0VfUE9TSVRJT05fVFJJTV9IT1JJWk9OEAESIAocRVZFTlRf'
    'U09VUkNFX1BPU0lUSU9OX0xBVEVTVBAC');

@$core.Deprecated('Use eventTypeDescriptor instead')
const EventType$json = {
  '1': 'EventType',
  '2': [
    {'1': 'EVENT_TYPE_CHAINEDINVOKESTARTED', '2': 0},
    {'1': 'EVENT_TYPE_INVOCATIONCOMPLETED', '2': 1},
    {'1': 'EVENT_TYPE_EXECUTIONTIMEDOUT', '2': 2},
    {'1': 'EVENT_TYPE_CALLBACKSUCCEEDED', '2': 3},
    {'1': 'EVENT_TYPE_EXECUTIONSTOPPED', '2': 4},
    {'1': 'EVENT_TYPE_WAITCANCELLED', '2': 5},
    {'1': 'EVENT_TYPE_STEPFAILED', '2': 6},
    {'1': 'EVENT_TYPE_EXECUTIONFAILED', '2': 7},
    {'1': 'EVENT_TYPE_CHAINEDINVOKESUCCEEDED', '2': 8},
    {'1': 'EVENT_TYPE_CALLBACKSTARTED', '2': 9},
    {'1': 'EVENT_TYPE_CONTEXTFAILED', '2': 10},
    {'1': 'EVENT_TYPE_WAITSTARTED', '2': 11},
    {'1': 'EVENT_TYPE_EXECUTIONSTARTED', '2': 12},
    {'1': 'EVENT_TYPE_CALLBACKTIMEDOUT', '2': 13},
    {'1': 'EVENT_TYPE_CHAINEDINVOKETIMEDOUT', '2': 14},
    {'1': 'EVENT_TYPE_STEPSTARTED', '2': 15},
    {'1': 'EVENT_TYPE_CALLBACKFAILED', '2': 16},
    {'1': 'EVENT_TYPE_CONTEXTSTARTED', '2': 17},
    {'1': 'EVENT_TYPE_CHAINEDINVOKEFAILED', '2': 18},
    {'1': 'EVENT_TYPE_EXECUTIONSUCCEEDED', '2': 19},
    {'1': 'EVENT_TYPE_CHAINEDINVOKESTOPPED', '2': 20},
    {'1': 'EVENT_TYPE_STEPSUCCEEDED', '2': 21},
    {'1': 'EVENT_TYPE_WAITSUCCEEDED', '2': 22},
    {'1': 'EVENT_TYPE_CONTEXTSUCCEEDED', '2': 23},
  ],
};

/// Descriptor for `EventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventTypeDescriptor = $convert.base64Decode(
    'CglFdmVudFR5cGUSIwofRVZFTlRfVFlQRV9DSEFJTkVESU5WT0tFU1RBUlRFRBAAEiIKHkVWRU'
    '5UX1RZUEVfSU5WT0NBVElPTkNPTVBMRVRFRBABEiAKHEVWRU5UX1RZUEVfRVhFQ1VUSU9OVElN'
    'RURPVVQQAhIgChxFVkVOVF9UWVBFX0NBTExCQUNLU1VDQ0VFREVEEAMSHwobRVZFTlRfVFlQRV'
    '9FWEVDVVRJT05TVE9QUEVEEAQSHAoYRVZFTlRfVFlQRV9XQUlUQ0FOQ0VMTEVEEAUSGQoVRVZF'
    'TlRfVFlQRV9TVEVQRkFJTEVEEAYSHgoaRVZFTlRfVFlQRV9FWEVDVVRJT05GQUlMRUQQBxIlCi'
    'FFVkVOVF9UWVBFX0NIQUlORURJTlZPS0VTVUNDRUVERUQQCBIeChpFVkVOVF9UWVBFX0NBTExC'
    'QUNLU1RBUlRFRBAJEhwKGEVWRU5UX1RZUEVfQ09OVEVYVEZBSUxFRBAKEhoKFkVWRU5UX1RZUE'
    'VfV0FJVFNUQVJURUQQCxIfChtFVkVOVF9UWVBFX0VYRUNVVElPTlNUQVJURUQQDBIfChtFVkVO'
    'VF9UWVBFX0NBTExCQUNLVElNRURPVVQQDRIkCiBFVkVOVF9UWVBFX0NIQUlORURJTlZPS0VUSU'
    '1FRE9VVBAOEhoKFkVWRU5UX1RZUEVfU1RFUFNUQVJURUQQDxIdChlFVkVOVF9UWVBFX0NBTExC'
    'QUNLRkFJTEVEEBASHQoZRVZFTlRfVFlQRV9DT05URVhUU1RBUlRFRBAREiIKHkVWRU5UX1RZUE'
    'VfQ0hBSU5FRElOVk9LRUZBSUxFRBASEiEKHUVWRU5UX1RZUEVfRVhFQ1VUSU9OU1VDQ0VFREVE'
    'EBMSIwofRVZFTlRfVFlQRV9DSEFJTkVESU5WT0tFU1RPUFBFRBAUEhwKGEVWRU5UX1RZUEVfU1'
    'RFUFNVQ0NFRURFRBAVEhwKGEVWRU5UX1RZUEVfV0FJVFNVQ0NFRURFRBAWEh8KG0VWRU5UX1RZ'
    'UEVfQ09OVEVYVFNVQ0NFRURFRBAX');

@$core.Deprecated('Use executionStatusDescriptor instead')
const ExecutionStatus$json = {
  '1': 'ExecutionStatus',
  '2': [
    {'1': 'EXECUTION_STATUS_RUNNING', '2': 0},
    {'1': 'EXECUTION_STATUS_TIMED_OUT', '2': 1},
    {'1': 'EXECUTION_STATUS_STOPPED', '2': 2},
    {'1': 'EXECUTION_STATUS_SUCCEEDED', '2': 3},
    {'1': 'EXECUTION_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `ExecutionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executionStatusDescriptor = $convert.base64Decode(
    'Cg9FeGVjdXRpb25TdGF0dXMSHAoYRVhFQ1VUSU9OX1NUQVRVU19SVU5OSU5HEAASHgoaRVhFQ1'
    'VUSU9OX1NUQVRVU19USU1FRF9PVVQQARIcChhFWEVDVVRJT05fU1RBVFVTX1NUT1BQRUQQAhIe'
    'ChpFWEVDVVRJT05fU1RBVFVTX1NVQ0NFRURFRBADEhsKF0VYRUNVVElPTl9TVEFUVVNfRkFJTE'
    'VEEAQ=');

@$core.Deprecated('Use fullDocumentDescriptor instead')
const FullDocument$json = {
  '1': 'FullDocument',
  '2': [
    {'1': 'FULL_DOCUMENT_UPDATELOOKUP', '2': 0},
    {'1': 'FULL_DOCUMENT_DEFAULT', '2': 1},
  ],
};

/// Descriptor for `FullDocument`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List fullDocumentDescriptor = $convert.base64Decode(
    'CgxGdWxsRG9jdW1lbnQSHgoaRlVMTF9ET0NVTUVOVF9VUERBVEVMT09LVVAQABIZChVGVUxMX0'
    'RPQ1VNRU5UX0RFRkFVTFQQAQ==');

@$core.Deprecated('Use functionResponseTypeDescriptor instead')
const FunctionResponseType$json = {
  '1': 'FunctionResponseType',
  '2': [
    {'1': 'FUNCTION_RESPONSE_TYPE_REPORTBATCHITEMFAILURES', '2': 0},
  ],
};

/// Descriptor for `FunctionResponseType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List functionResponseTypeDescriptor = $convert.base64Decode(
    'ChRGdW5jdGlvblJlc3BvbnNlVHlwZRIyCi5GVU5DVElPTl9SRVNQT05TRV9UWVBFX1JFUE9SVE'
    'JBVENISVRFTUZBSUxVUkVTEAA=');

@$core.Deprecated('Use functionUrlAuthTypeDescriptor instead')
const FunctionUrlAuthType$json = {
  '1': 'FunctionUrlAuthType',
  '2': [
    {'1': 'FUNCTION_URL_AUTH_TYPE_NONE', '2': 0},
    {'1': 'FUNCTION_URL_AUTH_TYPE_AWS_IAM', '2': 1},
  ],
};

/// Descriptor for `FunctionUrlAuthType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List functionUrlAuthTypeDescriptor = $convert.base64Decode(
    'ChNGdW5jdGlvblVybEF1dGhUeXBlEh8KG0ZVTkNUSU9OX1VSTF9BVVRIX1RZUEVfTk9ORRAAEi'
    'IKHkZVTkNUSU9OX1VSTF9BVVRIX1RZUEVfQVdTX0lBTRAB');

@$core.Deprecated('Use functionVersionDescriptor instead')
const FunctionVersion$json = {
  '1': 'FunctionVersion',
  '2': [
    {'1': 'FUNCTION_VERSION_ALL', '2': 0},
  ],
};

/// Descriptor for `FunctionVersion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List functionVersionDescriptor = $convert.base64Decode(
    'Cg9GdW5jdGlvblZlcnNpb24SGAoURlVOQ1RJT05fVkVSU0lPTl9BTEwQAA==');

@$core.Deprecated('Use functionVersionLatestPublishedDescriptor instead')
const FunctionVersionLatestPublished$json = {
  '1': 'FunctionVersionLatestPublished',
  '2': [
    {'1': 'FUNCTION_VERSION_LATEST_PUBLISHED_LATEST_PUBLISHED', '2': 0},
  ],
};

/// Descriptor for `FunctionVersionLatestPublished`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List functionVersionLatestPublishedDescriptor =
    $convert.base64Decode(
        'Ch5GdW5jdGlvblZlcnNpb25MYXRlc3RQdWJsaXNoZWQSNgoyRlVOQ1RJT05fVkVSU0lPTl9MQV'
        'RFU1RfUFVCTElTSEVEX0xBVEVTVF9QVUJMSVNIRUQQAA==');

@$core.Deprecated('Use invocationTypeDescriptor instead')
const InvocationType$json = {
  '1': 'InvocationType',
  '2': [
    {'1': 'INVOCATION_TYPE_EVENT', '2': 0},
    {'1': 'INVOCATION_TYPE_REQUESTRESPONSE', '2': 1},
    {'1': 'INVOCATION_TYPE_DRYRUN', '2': 2},
  ],
};

/// Descriptor for `InvocationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List invocationTypeDescriptor = $convert.base64Decode(
    'Cg5JbnZvY2F0aW9uVHlwZRIZChVJTlZPQ0FUSU9OX1RZUEVfRVZFTlQQABIjCh9JTlZPQ0FUSU'
    '9OX1RZUEVfUkVRVUVTVFJFU1BPTlNFEAESGgoWSU5WT0NBVElPTl9UWVBFX0RSWVJVThAC');

@$core.Deprecated('Use invokeModeDescriptor instead')
const InvokeMode$json = {
  '1': 'InvokeMode',
  '2': [
    {'1': 'INVOKE_MODE_RESPONSE_STREAM', '2': 0},
    {'1': 'INVOKE_MODE_BUFFERED', '2': 1},
  ],
};

/// Descriptor for `InvokeMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List invokeModeDescriptor = $convert.base64Decode(
    'CgpJbnZva2VNb2RlEh8KG0lOVk9LRV9NT0RFX1JFU1BPTlNFX1NUUkVBTRAAEhgKFElOVk9LRV'
    '9NT0RFX0JVRkZFUkVEEAE=');

@$core.Deprecated('Use kafkaSchemaRegistryAuthTypeDescriptor instead')
const KafkaSchemaRegistryAuthType$json = {
  '1': 'KafkaSchemaRegistryAuthType',
  '2': [
    {
      '1': 'KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_CLIENT_CERTIFICATE_TLS_AUTH',
      '2': 0
    },
    {'1': 'KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_SERVER_ROOT_CA_CERTIFICATE', '2': 1},
    {'1': 'KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_BASIC_AUTH', '2': 2},
  ],
};

/// Descriptor for `KafkaSchemaRegistryAuthType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List kafkaSchemaRegistryAuthTypeDescriptor = $convert.base64Decode(
    'ChtLYWZrYVNjaGVtYVJlZ2lzdHJ5QXV0aFR5cGUSPwo7S0FGS0FfU0NIRU1BX1JFR0lTVFJZX0'
    'FVVEhfVFlQRV9DTElFTlRfQ0VSVElGSUNBVEVfVExTX0FVVEgQABI+CjpLQUZLQV9TQ0hFTUFf'
    'UkVHSVNUUllfQVVUSF9UWVBFX1NFUlZFUl9ST09UX0NBX0NFUlRJRklDQVRFEAESLgoqS0FGS0'
    'FfU0NIRU1BX1JFR0lTVFJZX0FVVEhfVFlQRV9CQVNJQ19BVVRIEAI=');

@$core.Deprecated('Use kafkaSchemaValidationAttributeDescriptor instead')
const KafkaSchemaValidationAttribute$json = {
  '1': 'KafkaSchemaValidationAttribute',
  '2': [
    {'1': 'KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_VALUE', '2': 0},
    {'1': 'KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_KEY', '2': 1},
  ],
};

/// Descriptor for `KafkaSchemaValidationAttribute`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List kafkaSchemaValidationAttributeDescriptor =
    $convert.base64Decode(
        'Ch5LYWZrYVNjaGVtYVZhbGlkYXRpb25BdHRyaWJ1dGUSKwonS0FGS0FfU0NIRU1BX1ZBTElEQV'
        'RJT05fQVRUUklCVVRFX1ZBTFVFEAASKQolS0FGS0FfU0NIRU1BX1ZBTElEQVRJT05fQVRUUklC'
        'VVRFX0tFWRAB');

@$core.Deprecated('Use lastUpdateStatusDescriptor instead')
const LastUpdateStatus$json = {
  '1': 'LastUpdateStatus',
  '2': [
    {'1': 'LAST_UPDATE_STATUS_SUCCESSFUL', '2': 0},
    {'1': 'LAST_UPDATE_STATUS_FAILED', '2': 1},
    {'1': 'LAST_UPDATE_STATUS_INPROGRESS', '2': 2},
  ],
};

/// Descriptor for `LastUpdateStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List lastUpdateStatusDescriptor = $convert.base64Decode(
    'ChBMYXN0VXBkYXRlU3RhdHVzEiEKHUxBU1RfVVBEQVRFX1NUQVRVU19TVUNDRVNTRlVMEAASHQ'
    'oZTEFTVF9VUERBVEVfU1RBVFVTX0ZBSUxFRBABEiEKHUxBU1RfVVBEQVRFX1NUQVRVU19JTlBS'
    'T0dSRVNTEAI=');

@$core.Deprecated('Use lastUpdateStatusReasonCodeDescriptor instead')
const LastUpdateStatusReasonCode$json = {
  '1': 'LastUpdateStatusReasonCode',
  '2': [
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_ENILIMITEXCEEDED', '2': 0},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_SUBNETOUTOFIPADDRESSES', '2': 1},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_KMSKEYNOTFOUND', '2': 2},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITTIMEOUT', '2': 3},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDZIPFILEEXCEPTION', '2': 4},
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL',
      '2': 5
    },
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR',
      '2': 6
    },
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_EC2REQUESTLIMITEXCEEDED', '2': 7},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDIMAGE', '2': 8},
    {
      '1':
          'LAST_UPDATE_STATUS_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED',
      '2': 9
    },
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED',
      '2': 10
    },
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDRUNTIME', '2': 11},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_DISABLEDKMSKEY', '2': 12},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_VCPULIMITEXCEEDED', '2': 13},
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS',
      '2': 14
    },
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROR', '2': 15},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDCONFIGURATION', '2': 16},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDSUBNET', '2': 17},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDSTATEKMSKEY', '2': 18},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_IMAGEACCESSDENIED', '2': 19},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_IMAGEDELETED', '2': 20},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTTIMEOUT', '2': 21},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTFAILURE', '2': 22},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_EFSIOERROR', '2': 23},
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR',
      '2': 24
    },
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDSECURITYGROUP', '2': 25},
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT',
      '2': 26
    },
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTCONNECTIVITYERROR', '2': 27},
    {
      '1':
          'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY',
      '2': 28
    },
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTROLEPERMISSIONS',
      '2': 29
    },
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_KMSKEYACCESSDENIED', '2': 30},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INTERNALERROR', '2': 31},
    {'1': 'LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTCAPACITY', '2': 32},
    {
      '1': 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED',
      '2': 33
    },
  ],
};

/// Descriptor for `LastUpdateStatusReasonCode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List lastUpdateStatusReasonCodeDescriptor = $convert.base64Decode(
    'ChpMYXN0VXBkYXRlU3RhdHVzUmVhc29uQ29kZRIzCi9MQVNUX1VQREFURV9TVEFUVVNfUkVBU0'
    '9OX0NPREVfRU5JTElNSVRFWENFRURFRBAAEjkKNUxBU1RfVVBEQVRFX1NUQVRVU19SRUFTT05f'
    'Q09ERV9TVUJORVRPVVRPRklQQUREUkVTU0VTEAESMQotTEFTVF9VUERBVEVfU1RBVFVTX1JFQV'
    'NPTl9DT0RFX0tNU0tFWU5PVEZPVU5EEAISOwo3TEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9D'
    'T0RFX0ZVTkNUSU9ORVJST1JJTklUVElNRU9VVBADEjoKNkxBU1RfVVBEQVRFX1NUQVRVU19SRU'
    'FTT05fQ09ERV9JTlZBTElEWklQRklMRUVYQ0VQVElPThAEEkMKP0xBU1RfVVBEQVRFX1NUQVRV'
    'U19SRUFTT05fQ09ERV9ESVNBTExPV0VEQllWUENFTkNSWVBUSU9OQ09OVFJPTBAFEkAKPExBU1'
    'RfVVBEQVRFX1NUQVRVU19SRUFTT05fQ09ERV9GVU5DVElPTkVSUk9SUlVOVElNRUlOSVRFUlJP'
    'UhAGEjoKNkxBU1RfVVBEQVRFX1NUQVRVU19SRUFTT05fQ09ERV9FQzJSRVFVRVNUTElNSVRFWE'
    'NFRURFRBAHEi8KK0xBU1RfVVBEQVRFX1NUQVRVU19SRUFTT05fQ09ERV9JTlZBTElESU1BR0UQ'
    'CBJHCkNMQVNUX1VQREFURV9TVEFUVVNfUkVBU09OX0NPREVfQ0FQQUNJVFlQUk9WSURFUlNDQU'
    'xJTkdMSU1JVEVYQ0VFREVEEAkSQAo8TEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0ZV'
    'TkNUSU9ORVJST1JQRVJNSVNTSU9OREVOSUVEEAoSMQotTEFTVF9VUERBVEVfU1RBVFVTX1JFQV'
    'NPTl9DT0RFX0lOVkFMSURSVU5USU1FEAsSMQotTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9D'
    'T0RFX0RJU0FCTEVES01TS0VZEAwSNAowTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX1'
    'ZDUFVMSU1JVEVYQ0VFREVEEA0SQQo9TEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0ZV'
    'TkNUSU9ORVJST1JUT09NQU5ZRVhURU5TSU9OUxAOEjAKLExBU1RfVVBEQVRFX1NUQVRVU19SRU'
    'FTT05fQ09ERV9GVU5DVElPTkVSUk9SEA8SNwozTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9D'
    'T0RFX0lOVkFMSURDT05GSUdVUkFUSU9OEBASMAosTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl'
    '9DT0RFX0lOVkFMSURTVUJORVQQERI1CjFMQVNUX1VQREFURV9TVEFUVVNfUkVBU09OX0NPREVf'
    'SU5WQUxJRFNUQVRFS01TS0VZEBISNAowTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0'
    'lNQUdFQUNDRVNTREVOSUVEEBMSLworTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0lN'
    'QUdFREVMRVRFRBAUEjIKLkxBU1RfVVBEQVRFX1NUQVRVU19SRUFTT05fQ09ERV9FRlNNT1VOVF'
    'RJTUVPVVQQFRIyCi5MQVNUX1VQREFURV9TVEFUVVNfUkVBU09OX0NPREVfRUZTTU9VTlRGQUlM'
    'VVJFEBYSLQopTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0VGU0lPRVJST1IQFxJCCj'
    '5MQVNUX1VQREFURV9TVEFUVVNfUkVBU09OX0NPREVfRlVOQ1RJT05FUlJPUkVYVEVOU0lPTklO'
    'SVRFUlJPUhAYEjcKM0xBU1RfVVBEQVRFX1NUQVRVU19SRUFTT05fQ09ERV9JTlZBTElEU0VDVV'
    'JJVFlHUk9VUBAZEkEKPUxBU1RfVVBEQVRFX1NUQVRVU19SRUFTT05fQ09ERV9GVU5DVElPTkVS'
    'Uk9SSU5WQUxJREVOVFJZUE9JTlQQGhI8CjhMQVNUX1VQREFURV9TVEFUVVNfUkVBU09OX0NPRE'
    'VfRUZTTU9VTlRDT05ORUNUSVZJVFlFUlJPUhAbEkcKQ0xBU1RfVVBEQVRFX1NUQVRVU19SRUFT'
    'T05fQ09ERV9GVU5DVElPTkVSUk9SSU5WQUxJRFdPUktJTkdESVJFQ1RPUlkQHBI+CjpMQVNUX1'
    'VQREFURV9TVEFUVVNfUkVBU09OX0NPREVfSU5TVUZGSUNJRU5UUk9MRVBFUk1JU1NJT05TEB0S'
    'NQoxTEFTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0tNU0tFWUFDQ0VTU0RFTklFRBAeEj'
    'AKLExBU1RfVVBEQVRFX1NUQVRVU19SRUFTT05fQ09ERV9JTlRFUk5BTEVSUk9SEB8SNwozTEFT'
    'VF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0lOU1VGRklDSUVOVENBUEFDSVRZECASRQpBTE'
    'FTVF9VUERBVEVfU1RBVFVTX1JFQVNPTl9DT0RFX0ZVTkNUSU9ORVJST1JJTklUUkVTT1VSQ0VF'
    'WEhBVVNURUQQIQ==');

@$core.Deprecated('Use logFormatDescriptor instead')
const LogFormat$json = {
  '1': 'LogFormat',
  '2': [
    {'1': 'LOG_FORMAT_TEXT', '2': 0},
    {'1': 'LOG_FORMAT_JSON', '2': 1},
  ],
};

/// Descriptor for `LogFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List logFormatDescriptor = $convert.base64Decode(
    'CglMb2dGb3JtYXQSEwoPTE9HX0ZPUk1BVF9URVhUEAASEwoPTE9HX0ZPUk1BVF9KU09OEAE=');

@$core.Deprecated('Use logTypeDescriptor instead')
const LogType$json = {
  '1': 'LogType',
  '2': [
    {'1': 'LOG_TYPE_NONE', '2': 0},
    {'1': 'LOG_TYPE_TAIL', '2': 1},
  ],
};

/// Descriptor for `LogType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List logTypeDescriptor = $convert.base64Decode(
    'CgdMb2dUeXBlEhEKDUxPR19UWVBFX05PTkUQABIRCg1MT0dfVFlQRV9UQUlMEAE=');

@$core.Deprecated('Use operationActionDescriptor instead')
const OperationAction$json = {
  '1': 'OperationAction',
  '2': [
    {'1': 'OPERATION_ACTION_SUCCEED', '2': 0},
    {'1': 'OPERATION_ACTION_RETRY', '2': 1},
    {'1': 'OPERATION_ACTION_START', '2': 2},
    {'1': 'OPERATION_ACTION_CANCEL', '2': 3},
    {'1': 'OPERATION_ACTION_FAIL', '2': 4},
  ],
};

/// Descriptor for `OperationAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List operationActionDescriptor = $convert.base64Decode(
    'Cg9PcGVyYXRpb25BY3Rpb24SHAoYT1BFUkFUSU9OX0FDVElPTl9TVUNDRUVEEAASGgoWT1BFUk'
    'FUSU9OX0FDVElPTl9SRVRSWRABEhoKFk9QRVJBVElPTl9BQ1RJT05fU1RBUlQQAhIbChdPUEVS'
    'QVRJT05fQUNUSU9OX0NBTkNFTBADEhkKFU9QRVJBVElPTl9BQ1RJT05fRkFJTBAE');

@$core.Deprecated('Use operationStatusDescriptor instead')
const OperationStatus$json = {
  '1': 'OperationStatus',
  '2': [
    {'1': 'OPERATION_STATUS_PENDING', '2': 0},
    {'1': 'OPERATION_STATUS_TIMED_OUT', '2': 1},
    {'1': 'OPERATION_STATUS_STOPPED', '2': 2},
    {'1': 'OPERATION_STATUS_STARTED', '2': 3},
    {'1': 'OPERATION_STATUS_READY', '2': 4},
    {'1': 'OPERATION_STATUS_SUCCEEDED', '2': 5},
    {'1': 'OPERATION_STATUS_CANCELLED', '2': 6},
    {'1': 'OPERATION_STATUS_FAILED', '2': 7},
  ],
};

/// Descriptor for `OperationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List operationStatusDescriptor = $convert.base64Decode(
    'Cg9PcGVyYXRpb25TdGF0dXMSHAoYT1BFUkFUSU9OX1NUQVRVU19QRU5ESU5HEAASHgoaT1BFUk'
    'FUSU9OX1NUQVRVU19USU1FRF9PVVQQARIcChhPUEVSQVRJT05fU1RBVFVTX1NUT1BQRUQQAhIc'
    'ChhPUEVSQVRJT05fU1RBVFVTX1NUQVJURUQQAxIaChZPUEVSQVRJT05fU1RBVFVTX1JFQURZEA'
    'QSHgoaT1BFUkFUSU9OX1NUQVRVU19TVUNDRUVERUQQBRIeChpPUEVSQVRJT05fU1RBVFVTX0NB'
    'TkNFTExFRBAGEhsKF09QRVJBVElPTl9TVEFUVVNfRkFJTEVEEAc=');

@$core.Deprecated('Use operationTypeDescriptor instead')
const OperationType$json = {
  '1': 'OperationType',
  '2': [
    {'1': 'OPERATION_TYPE_CALLBACK', '2': 0},
    {'1': 'OPERATION_TYPE_STEP', '2': 1},
    {'1': 'OPERATION_TYPE_CHAINED_INVOKE', '2': 2},
    {'1': 'OPERATION_TYPE_EXECUTION', '2': 3},
    {'1': 'OPERATION_TYPE_WAIT', '2': 4},
    {'1': 'OPERATION_TYPE_CONTEXT', '2': 5},
  ],
};

/// Descriptor for `OperationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List operationTypeDescriptor = $convert.base64Decode(
    'Cg1PcGVyYXRpb25UeXBlEhsKF09QRVJBVElPTl9UWVBFX0NBTExCQUNLEAASFwoTT1BFUkFUSU'
    '9OX1RZUEVfU1RFUBABEiEKHU9QRVJBVElPTl9UWVBFX0NIQUlORURfSU5WT0tFEAISHAoYT1BF'
    'UkFUSU9OX1RZUEVfRVhFQ1VUSU9OEAMSFwoTT1BFUkFUSU9OX1RZUEVfV0FJVBAEEhoKFk9QRV'
    'JBVElPTl9UWVBFX0NPTlRFWFQQBQ==');

@$core.Deprecated('Use packageTypeDescriptor instead')
const PackageType$json = {
  '1': 'PackageType',
  '2': [
    {'1': 'PACKAGE_TYPE_IMAGE', '2': 0},
    {'1': 'PACKAGE_TYPE_ZIP', '2': 1},
  ],
};

/// Descriptor for `PackageType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List packageTypeDescriptor = $convert.base64Decode(
    'CgtQYWNrYWdlVHlwZRIWChJQQUNLQUdFX1RZUEVfSU1BR0UQABIUChBQQUNLQUdFX1RZUEVfWk'
    'lQEAE=');

@$core.Deprecated('Use provisionedConcurrencyStatusEnumDescriptor instead')
const ProvisionedConcurrencyStatusEnum$json = {
  '1': 'ProvisionedConcurrencyStatusEnum',
  '2': [
    {'1': 'PROVISIONED_CONCURRENCY_STATUS_ENUM_IN_PROGRESS', '2': 0},
    {'1': 'PROVISIONED_CONCURRENCY_STATUS_ENUM_READY', '2': 1},
    {'1': 'PROVISIONED_CONCURRENCY_STATUS_ENUM_FAILED', '2': 2},
  ],
};

/// Descriptor for `ProvisionedConcurrencyStatusEnum`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List provisionedConcurrencyStatusEnumDescriptor =
    $convert.base64Decode(
        'CiBQcm92aXNpb25lZENvbmN1cnJlbmN5U3RhdHVzRW51bRIzCi9QUk9WSVNJT05FRF9DT05DVV'
        'JSRU5DWV9TVEFUVVNfRU5VTV9JTl9QUk9HUkVTUxAAEi0KKVBST1ZJU0lPTkVEX0NPTkNVUlJF'
        'TkNZX1NUQVRVU19FTlVNX1JFQURZEAESLgoqUFJPVklTSU9ORURfQ09OQ1VSUkVOQ1lfU1RBVF'
        'VTX0VOVU1fRkFJTEVEEAI=');

@$core.Deprecated('Use recursiveLoopDescriptor instead')
const RecursiveLoop$json = {
  '1': 'RecursiveLoop',
  '2': [
    {'1': 'RECURSIVE_LOOP_ALLOW', '2': 0},
    {'1': 'RECURSIVE_LOOP_TERMINATE', '2': 1},
  ],
};

/// Descriptor for `RecursiveLoop`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recursiveLoopDescriptor = $convert.base64Decode(
    'Cg1SZWN1cnNpdmVMb29wEhgKFFJFQ1VSU0lWRV9MT09QX0FMTE9XEAASHAoYUkVDVVJTSVZFX0'
    'xPT1BfVEVSTUlOQVRFEAE=');

@$core.Deprecated('Use responseStreamingInvocationTypeDescriptor instead')
const ResponseStreamingInvocationType$json = {
  '1': 'ResponseStreamingInvocationType',
  '2': [
    {'1': 'RESPONSE_STREAMING_INVOCATION_TYPE_REQUESTRESPONSE', '2': 0},
    {'1': 'RESPONSE_STREAMING_INVOCATION_TYPE_DRYRUN', '2': 1},
  ],
};

/// Descriptor for `ResponseStreamingInvocationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List responseStreamingInvocationTypeDescriptor =
    $convert.base64Decode(
        'Ch9SZXNwb25zZVN0cmVhbWluZ0ludm9jYXRpb25UeXBlEjYKMlJFU1BPTlNFX1NUUkVBTUlOR1'
        '9JTlZPQ0FUSU9OX1RZUEVfUkVRVUVTVFJFU1BPTlNFEAASLQopUkVTUE9OU0VfU1RSRUFNSU5H'
        'X0lOVk9DQVRJT05fVFlQRV9EUllSVU4QAQ==');

@$core.Deprecated('Use runtimeDescriptor instead')
const Runtime$json = {
  '1': 'Runtime',
  '2': [
    {'1': 'RUNTIME_DOTNETCORE31', '2': 0},
    {'1': 'RUNTIME_NODEJS16X', '2': 1},
    {'1': 'RUNTIME_RUBY25', '2': 2},
    {'1': 'RUNTIME_NODEJS24X', '2': 3},
    {'1': 'RUNTIME_PYTHON39', '2': 4},
    {'1': 'RUNTIME_DOTNET6', '2': 5},
    {'1': 'RUNTIME_RUBY32', '2': 6},
    {'1': 'RUNTIME_GO1X', '2': 7},
    {'1': 'RUNTIME_PYTHON36', '2': 8},
    {'1': 'RUNTIME_PROVIDEDAL2023', '2': 9},
    {'1': 'RUNTIME_NODEJS22X', '2': 10},
    {'1': 'RUNTIME_PYTHON312', '2': 11},
    {'1': 'RUNTIME_DOTNETCORE21', '2': 12},
    {'1': 'RUNTIME_NODEJS810', '2': 13},
    {'1': 'RUNTIME_RUBY27', '2': 14},
    {'1': 'RUNTIME_JAVA21', '2': 15},
    {'1': 'RUNTIME_NODEJS14X', '2': 16},
    {'1': 'RUNTIME_PROVIDED', '2': 17},
    {'1': 'RUNTIME_NODEJS610', '2': 18},
    {'1': 'RUNTIME_PYTHON310', '2': 19},
    {'1': 'RUNTIME_PYTHON38', '2': 20},
    {'1': 'RUNTIME_JAVA17', '2': 21},
    {'1': 'RUNTIME_NODEJS43', '2': 22},
    {'1': 'RUNTIME_NODEJS', '2': 23},
    {'1': 'RUNTIME_NODEJS20X', '2': 24},
    {'1': 'RUNTIME_PYTHON313', '2': 25},
    {'1': 'RUNTIME_DOTNETCORE20', '2': 26},
    {'1': 'RUNTIME_NODEJS43EDGE', '2': 27},
    {'1': 'RUNTIME_JAVA11', '2': 28},
    {'1': 'RUNTIME_NODEJS12X', '2': 29},
    {'1': 'RUNTIME_RUBY33', '2': 30},
    {'1': 'RUNTIME_JAVA8', '2': 31},
    {'1': 'RUNTIME_DOTNETCORE10', '2': 32},
    {'1': 'RUNTIME_PYTHON37', '2': 33},
    {'1': 'RUNTIME_PYTHON311', '2': 34},
    {'1': 'RUNTIME_JAVA8AL2', '2': 35},
    {'1': 'RUNTIME_RUBY34', '2': 36},
    {'1': 'RUNTIME_PROVIDEDAL2', '2': 37},
    {'1': 'RUNTIME_NODEJS10X', '2': 38},
    {'1': 'RUNTIME_JAVA25', '2': 39},
    {'1': 'RUNTIME_PYTHON27', '2': 40},
    {'1': 'RUNTIME_DOTNET8', '2': 41},
    {'1': 'RUNTIME_NODEJS18X', '2': 42},
    {'1': 'RUNTIME_DOTNET10', '2': 43},
    {'1': 'RUNTIME_PYTHON314', '2': 44},
  ],
};

/// Descriptor for `Runtime`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List runtimeDescriptor = $convert.base64Decode(
    'CgdSdW50aW1lEhgKFFJVTlRJTUVfRE9UTkVUQ09SRTMxEAASFQoRUlVOVElNRV9OT0RFSlMxNl'
    'gQARISCg5SVU5USU1FX1JVQlkyNRACEhUKEVJVTlRJTUVfTk9ERUpTMjRYEAMSFAoQUlVOVElN'
    'RV9QWVRIT04zORAEEhMKD1JVTlRJTUVfRE9UTkVUNhAFEhIKDlJVTlRJTUVfUlVCWTMyEAYSEA'
    'oMUlVOVElNRV9HTzFYEAcSFAoQUlVOVElNRV9QWVRIT04zNhAIEhoKFlJVTlRJTUVfUFJPVklE'
    'RURBTDIwMjMQCRIVChFSVU5USU1FX05PREVKUzIyWBAKEhUKEVJVTlRJTUVfUFlUSE9OMzEyEA'
    'sSGAoUUlVOVElNRV9ET1RORVRDT1JFMjEQDBIVChFSVU5USU1FX05PREVKUzgxMBANEhIKDlJV'
    'TlRJTUVfUlVCWTI3EA4SEgoOUlVOVElNRV9KQVZBMjEQDxIVChFSVU5USU1FX05PREVKUzE0WB'
    'AQEhQKEFJVTlRJTUVfUFJPVklERUQQERIVChFSVU5USU1FX05PREVKUzYxMBASEhUKEVJVTlRJ'
    'TUVfUFlUSE9OMzEwEBMSFAoQUlVOVElNRV9QWVRIT04zOBAUEhIKDlJVTlRJTUVfSkFWQTE3EB'
    'USFAoQUlVOVElNRV9OT0RFSlM0MxAWEhIKDlJVTlRJTUVfTk9ERUpTEBcSFQoRUlVOVElNRV9O'
    'T0RFSlMyMFgQGBIVChFSVU5USU1FX1BZVEhPTjMxMxAZEhgKFFJVTlRJTUVfRE9UTkVUQ09SRT'
    'IwEBoSGAoUUlVOVElNRV9OT0RFSlM0M0VER0UQGxISCg5SVU5USU1FX0pBVkExMRAcEhUKEVJV'
    'TlRJTUVfTk9ERUpTMTJYEB0SEgoOUlVOVElNRV9SVUJZMzMQHhIRCg1SVU5USU1FX0pBVkE4EB'
    '8SGAoUUlVOVElNRV9ET1RORVRDT1JFMTAQIBIUChBSVU5USU1FX1BZVEhPTjM3ECESFQoRUlVO'
    'VElNRV9QWVRIT04zMTEQIhIUChBSVU5USU1FX0pBVkE4QUwyECMSEgoOUlVOVElNRV9SVUJZMz'
    'QQJBIXChNSVU5USU1FX1BST1ZJREVEQUwyECUSFQoRUlVOVElNRV9OT0RFSlMxMFgQJhISCg5S'
    'VU5USU1FX0pBVkEyNRAnEhQKEFJVTlRJTUVfUFlUSE9OMjcQKBITCg9SVU5USU1FX0RPVE5FVD'
    'gQKRIVChFSVU5USU1FX05PREVKUzE4WBAqEhQKEFJVTlRJTUVfRE9UTkVUMTAQKxIVChFSVU5U'
    'SU1FX1BZVEhPTjMxNBAs');

@$core.Deprecated('Use schemaRegistryEventRecordFormatDescriptor instead')
const SchemaRegistryEventRecordFormat$json = {
  '1': 'SchemaRegistryEventRecordFormat',
  '2': [
    {'1': 'SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_JSON', '2': 0},
    {'1': 'SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_SOURCE', '2': 1},
  ],
};

/// Descriptor for `SchemaRegistryEventRecordFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List schemaRegistryEventRecordFormatDescriptor =
    $convert.base64Decode(
        'Ch9TY2hlbWFSZWdpc3RyeUV2ZW50UmVjb3JkRm9ybWF0EiwKKFNDSEVNQV9SRUdJU1RSWV9FVk'
        'VOVF9SRUNPUkRfRk9STUFUX0pTT04QABIuCipTQ0hFTUFfUkVHSVNUUllfRVZFTlRfUkVDT1JE'
        'X0ZPUk1BVF9TT1VSQ0UQAQ==');

@$core.Deprecated('Use snapStartApplyOnDescriptor instead')
const SnapStartApplyOn$json = {
  '1': 'SnapStartApplyOn',
  '2': [
    {'1': 'SNAP_START_APPLY_ON_NONE', '2': 0},
    {'1': 'SNAP_START_APPLY_ON_PUBLISHEDVERSIONS', '2': 1},
  ],
};

/// Descriptor for `SnapStartApplyOn`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List snapStartApplyOnDescriptor = $convert.base64Decode(
    'ChBTbmFwU3RhcnRBcHBseU9uEhwKGFNOQVBfU1RBUlRfQVBQTFlfT05fTk9ORRAAEikKJVNOQV'
    'BfU1RBUlRfQVBQTFlfT05fUFVCTElTSEVEVkVSU0lPTlMQAQ==');

@$core.Deprecated('Use snapStartOptimizationStatusDescriptor instead')
const SnapStartOptimizationStatus$json = {
  '1': 'SnapStartOptimizationStatus',
  '2': [
    {'1': 'SNAP_START_OPTIMIZATION_STATUS_ON', '2': 0},
    {'1': 'SNAP_START_OPTIMIZATION_STATUS_OFF', '2': 1},
  ],
};

/// Descriptor for `SnapStartOptimizationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List snapStartOptimizationStatusDescriptor =
    $convert.base64Decode(
        'ChtTbmFwU3RhcnRPcHRpbWl6YXRpb25TdGF0dXMSJQohU05BUF9TVEFSVF9PUFRJTUlaQVRJT0'
        '5fU1RBVFVTX09OEAASJgoiU05BUF9TVEFSVF9PUFRJTUlaQVRJT05fU1RBVFVTX09GRhAB');

@$core.Deprecated('Use sourceAccessTypeDescriptor instead')
const SourceAccessType$json = {
  '1': 'SourceAccessType',
  '2': [
    {'1': 'SOURCE_ACCESS_TYPE_VPC_SUBNET', '2': 0},
    {'1': 'SOURCE_ACCESS_TYPE_VPC_SECURITY_GROUP', '2': 1},
    {'1': 'SOURCE_ACCESS_TYPE_CLIENT_CERTIFICATE_TLS_AUTH', '2': 2},
    {'1': 'SOURCE_ACCESS_TYPE_SASL_SCRAM_256_AUTH', '2': 3},
    {'1': 'SOURCE_ACCESS_TYPE_VIRTUAL_HOST', '2': 4},
    {'1': 'SOURCE_ACCESS_TYPE_SASL_SCRAM_512_AUTH', '2': 5},
    {'1': 'SOURCE_ACCESS_TYPE_SERVER_ROOT_CA_CERTIFICATE', '2': 6},
    {'1': 'SOURCE_ACCESS_TYPE_BASIC_AUTH', '2': 7},
  ],
};

/// Descriptor for `SourceAccessType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sourceAccessTypeDescriptor = $convert.base64Decode(
    'ChBTb3VyY2VBY2Nlc3NUeXBlEiEKHVNPVVJDRV9BQ0NFU1NfVFlQRV9WUENfU1VCTkVUEAASKQ'
    'olU09VUkNFX0FDQ0VTU19UWVBFX1ZQQ19TRUNVUklUWV9HUk9VUBABEjIKLlNPVVJDRV9BQ0NF'
    'U1NfVFlQRV9DTElFTlRfQ0VSVElGSUNBVEVfVExTX0FVVEgQAhIqCiZTT1VSQ0VfQUNDRVNTX1'
    'RZUEVfU0FTTF9TQ1JBTV8yNTZfQVVUSBADEiMKH1NPVVJDRV9BQ0NFU1NfVFlQRV9WSVJUVUFM'
    'X0hPU1QQBBIqCiZTT1VSQ0VfQUNDRVNTX1RZUEVfU0FTTF9TQ1JBTV81MTJfQVVUSBAFEjEKLV'
    'NPVVJDRV9BQ0NFU1NfVFlQRV9TRVJWRVJfUk9PVF9DQV9DRVJUSUZJQ0FURRAGEiEKHVNPVVJD'
    'RV9BQ0NFU1NfVFlQRV9CQVNJQ19BVVRIEAc=');

@$core.Deprecated('Use stateDescriptor instead')
const State$json = {
  '1': 'State',
  '2': [
    {'1': 'STATE_ACTIVE', '2': 0},
    {'1': 'STATE_DEACTIVATED', '2': 1},
    {'1': 'STATE_FAILED', '2': 2},
    {'1': 'STATE_ACTIVENONINVOCABLE', '2': 3},
    {'1': 'STATE_INACTIVE', '2': 4},
    {'1': 'STATE_DELETING', '2': 5},
    {'1': 'STATE_DEACTIVATING', '2': 6},
    {'1': 'STATE_PENDING', '2': 7},
  ],
};

/// Descriptor for `State`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stateDescriptor = $convert.base64Decode(
    'CgVTdGF0ZRIQCgxTVEFURV9BQ1RJVkUQABIVChFTVEFURV9ERUFDVElWQVRFRBABEhAKDFNUQV'
    'RFX0ZBSUxFRBACEhwKGFNUQVRFX0FDVElWRU5PTklOVk9DQUJMRRADEhIKDlNUQVRFX0lOQUNU'
    'SVZFEAQSEgoOU1RBVEVfREVMRVRJTkcQBRIWChJTVEFURV9ERUFDVElWQVRJTkcQBhIRCg1TVE'
    'FURV9QRU5ESU5HEAc=');

@$core.Deprecated('Use stateReasonCodeDescriptor instead')
const StateReasonCode$json = {
  '1': 'StateReasonCode',
  '2': [
    {'1': 'STATE_REASON_CODE_ENILIMITEXCEEDED', '2': 0},
    {'1': 'STATE_REASON_CODE_SUBNETOUTOFIPADDRESSES', '2': 1},
    {'1': 'STATE_REASON_CODE_IDLE', '2': 2},
    {'1': 'STATE_REASON_CODE_KMSKEYNOTFOUND', '2': 3},
    {'1': 'STATE_REASON_CODE_FUNCTIONERRORINITTIMEOUT', '2': 4},
    {'1': 'STATE_REASON_CODE_INVALIDZIPFILEEXCEPTION', '2': 5},
    {'1': 'STATE_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL', '2': 6},
    {'1': 'STATE_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR', '2': 7},
    {'1': 'STATE_REASON_CODE_RESTORING', '2': 8},
    {'1': 'STATE_REASON_CODE_EC2REQUESTLIMITEXCEEDED', '2': 9},
    {'1': 'STATE_REASON_CODE_INVALIDIMAGE', '2': 10},
    {'1': 'STATE_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED', '2': 11},
    {'1': 'STATE_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED', '2': 12},
    {'1': 'STATE_REASON_CODE_INVALIDRUNTIME', '2': 13},
    {'1': 'STATE_REASON_CODE_DISABLEDKMSKEY', '2': 14},
    {'1': 'STATE_REASON_CODE_VCPULIMITEXCEEDED', '2': 15},
    {'1': 'STATE_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS', '2': 16},
    {'1': 'STATE_REASON_CODE_FUNCTIONERROR', '2': 17},
    {'1': 'STATE_REASON_CODE_INVALIDCONFIGURATION', '2': 18},
    {'1': 'STATE_REASON_CODE_INVALIDSUBNET', '2': 19},
    {'1': 'STATE_REASON_CODE_INVALIDSTATEKMSKEY', '2': 20},
    {'1': 'STATE_REASON_CODE_IMAGEACCESSDENIED', '2': 21},
    {'1': 'STATE_REASON_CODE_IMAGEDELETED', '2': 22},
    {'1': 'STATE_REASON_CODE_EFSMOUNTTIMEOUT', '2': 23},
    {'1': 'STATE_REASON_CODE_EFSMOUNTFAILURE', '2': 24},
    {'1': 'STATE_REASON_CODE_EFSIOERROR', '2': 25},
    {'1': 'STATE_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR', '2': 26},
    {'1': 'STATE_REASON_CODE_CREATING', '2': 27},
    {'1': 'STATE_REASON_CODE_INVALIDSECURITYGROUP', '2': 28},
    {'1': 'STATE_REASON_CODE_DRAININGDURABLEEXECUTIONS', '2': 29},
    {'1': 'STATE_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT', '2': 30},
    {'1': 'STATE_REASON_CODE_EFSMOUNTCONNECTIVITYERROR', '2': 31},
    {'1': 'STATE_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY', '2': 32},
    {'1': 'STATE_REASON_CODE_INSUFFICIENTROLEPERMISSIONS', '2': 33},
    {'1': 'STATE_REASON_CODE_KMSKEYACCESSDENIED', '2': 34},
    {'1': 'STATE_REASON_CODE_INTERNALERROR', '2': 35},
    {'1': 'STATE_REASON_CODE_INSUFFICIENTCAPACITY', '2': 36},
    {'1': 'STATE_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED', '2': 37},
  ],
};

/// Descriptor for `StateReasonCode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stateReasonCodeDescriptor = $convert.base64Decode(
    'Cg9TdGF0ZVJlYXNvbkNvZGUSJgoiU1RBVEVfUkVBU09OX0NPREVfRU5JTElNSVRFWENFRURFRB'
    'AAEiwKKFNUQVRFX1JFQVNPTl9DT0RFX1NVQk5FVE9VVE9GSVBBRERSRVNTRVMQARIaChZTVEFU'
    'RV9SRUFTT05fQ09ERV9JRExFEAISJAogU1RBVEVfUkVBU09OX0NPREVfS01TS0VZTk9URk9VTk'
    'QQAxIuCipTVEFURV9SRUFTT05fQ09ERV9GVU5DVElPTkVSUk9SSU5JVFRJTUVPVVQQBBItCilT'
    'VEFURV9SRUFTT05fQ09ERV9JTlZBTElEWklQRklMRUVYQ0VQVElPThAFEjYKMlNUQVRFX1JFQV'
    'NPTl9DT0RFX0RJU0FMTE9XRURCWVZQQ0VOQ1JZUFRJT05DT05UUk9MEAYSMwovU1RBVEVfUkVB'
    'U09OX0NPREVfRlVOQ1RJT05FUlJPUlJVTlRJTUVJTklURVJST1IQBxIfChtTVEFURV9SRUFTT0'
    '5fQ09ERV9SRVNUT1JJTkcQCBItCilTVEFURV9SRUFTT05fQ09ERV9FQzJSRVFVRVNUTElNSVRF'
    'WENFRURFRBAJEiIKHlNUQVRFX1JFQVNPTl9DT0RFX0lOVkFMSURJTUFHRRAKEjoKNlNUQVRFX1'
    'JFQVNPTl9DT0RFX0NBUEFDSVRZUFJPVklERVJTQ0FMSU5HTElNSVRFWENFRURFRBALEjMKL1NU'
    'QVRFX1JFQVNPTl9DT0RFX0ZVTkNUSU9ORVJST1JQRVJNSVNTSU9OREVOSUVEEAwSJAogU1RBVE'
    'VfUkVBU09OX0NPREVfSU5WQUxJRFJVTlRJTUUQDRIkCiBTVEFURV9SRUFTT05fQ09ERV9ESVNB'
    'QkxFREtNU0tFWRAOEicKI1NUQVRFX1JFQVNPTl9DT0RFX1ZDUFVMSU1JVEVYQ0VFREVEEA8SNA'
    'owU1RBVEVfUkVBU09OX0NPREVfRlVOQ1RJT05FUlJPUlRPT01BTllFWFRFTlNJT05TEBASIwof'
    'U1RBVEVfUkVBU09OX0NPREVfRlVOQ1RJT05FUlJPUhAREioKJlNUQVRFX1JFQVNPTl9DT0RFX0'
    'lOVkFMSURDT05GSUdVUkFUSU9OEBISIwofU1RBVEVfUkVBU09OX0NPREVfSU5WQUxJRFNVQk5F'
    'VBATEigKJFNUQVRFX1JFQVNPTl9DT0RFX0lOVkFMSURTVEFURUtNU0tFWRAUEicKI1NUQVRFX1'
    'JFQVNPTl9DT0RFX0lNQUdFQUNDRVNTREVOSUVEEBUSIgoeU1RBVEVfUkVBU09OX0NPREVfSU1B'
    'R0VERUxFVEVEEBYSJQohU1RBVEVfUkVBU09OX0NPREVfRUZTTU9VTlRUSU1FT1VUEBcSJQohU1'
    'RBVEVfUkVBU09OX0NPREVfRUZTTU9VTlRGQUlMVVJFEBgSIAocU1RBVEVfUkVBU09OX0NPREVf'
    'RUZTSU9FUlJPUhAZEjUKMVNUQVRFX1JFQVNPTl9DT0RFX0ZVTkNUSU9ORVJST1JFWFRFTlNJT0'
    '5JTklURVJST1IQGhIeChpTVEFURV9SRUFTT05fQ09ERV9DUkVBVElORxAbEioKJlNUQVRFX1JF'
    'QVNPTl9DT0RFX0lOVkFMSURTRUNVUklUWUdST1VQEBwSLworU1RBVEVfUkVBU09OX0NPREVfRF'
    'JBSU5JTkdEVVJBQkxFRVhFQ1VUSU9OUxAdEjQKMFNUQVRFX1JFQVNPTl9DT0RFX0ZVTkNUSU9O'
    'RVJST1JJTlZBTElERU5UUllQT0lOVBAeEi8KK1NUQVRFX1JFQVNPTl9DT0RFX0VGU01PVU5UQ0'
    '9OTkVDVElWSVRZRVJST1IQHxI6CjZTVEFURV9SRUFTT05fQ09ERV9GVU5DVElPTkVSUk9SSU5W'
    'QUxJRFdPUktJTkdESVJFQ1RPUlkQIBIxCi1TVEFURV9SRUFTT05fQ09ERV9JTlNVRkZJQ0lFTl'
    'RST0xFUEVSTUlTU0lPTlMQIRIoCiRTVEFURV9SRUFTT05fQ09ERV9LTVNLRVlBQ0NFU1NERU5J'
    'RUQQIhIjCh9TVEFURV9SRUFTT05fQ09ERV9JTlRFUk5BTEVSUk9SECMSKgomU1RBVEVfUkVBU0'
    '9OX0NPREVfSU5TVUZGSUNJRU5UQ0FQQUNJVFkQJBI4CjRTVEFURV9SRUFTT05fQ09ERV9GVU5D'
    'VElPTkVSUk9SSU5JVFJFU09VUkNFRVhIQVVTVEVEECU=');

@$core.Deprecated('Use systemLogLevelDescriptor instead')
const SystemLogLevel$json = {
  '1': 'SystemLogLevel',
  '2': [
    {'1': 'SYSTEM_LOG_LEVEL_WARN', '2': 0},
    {'1': 'SYSTEM_LOG_LEVEL_INFO', '2': 1},
    {'1': 'SYSTEM_LOG_LEVEL_DEBUG', '2': 2},
  ],
};

/// Descriptor for `SystemLogLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List systemLogLevelDescriptor = $convert.base64Decode(
    'Cg5TeXN0ZW1Mb2dMZXZlbBIZChVTWVNURU1fTE9HX0xFVkVMX1dBUk4QABIZChVTWVNURU1fTE'
    '9HX0xFVkVMX0lORk8QARIaChZTWVNURU1fTE9HX0xFVkVMX0RFQlVHEAI=');

@$core.Deprecated('Use tenantIsolationModeDescriptor instead')
const TenantIsolationMode$json = {
  '1': 'TenantIsolationMode',
  '2': [
    {'1': 'TENANT_ISOLATION_MODE_PER_TENANT', '2': 0},
  ],
};

/// Descriptor for `TenantIsolationMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tenantIsolationModeDescriptor = $convert.base64Decode(
    'ChNUZW5hbnRJc29sYXRpb25Nb2RlEiQKIFRFTkFOVF9JU09MQVRJT05fTU9ERV9QRVJfVEVOQU'
    '5UEAA=');

@$core.Deprecated('Use throttleReasonDescriptor instead')
const ThrottleReason$json = {
  '1': 'ThrottleReason',
  '2': [
    {'1': 'THROTTLE_REASON_CONCURRENTSNAPSHOTCREATELIMITEXCEEDED', '2': 0},
    {'1': 'THROTTLE_REASON_FUNCTIONINVOCATIONRATELIMITEXCEEDED', '2': 1},
    {'1': 'THROTTLE_REASON_CONCURRENTINVOCATIONLIMITEXCEEDED', '2': 2},
    {'1': 'THROTTLE_REASON_CALLERRATELIMITEXCEEDED', '2': 3},
    {
      '1': 'THROTTLE_REASON_RESERVEDFUNCTIONCONCURRENTINVOCATIONLIMITEXCEEDED',
      '2': 4
    },
    {
      '1': 'THROTTLE_REASON_RESERVEDFUNCTIONINVOCATIONRATELIMITEXCEEDED',
      '2': 5
    },
  ],
};

/// Descriptor for `ThrottleReason`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List throttleReasonDescriptor = $convert.base64Decode(
    'Cg5UaHJvdHRsZVJlYXNvbhI5CjVUSFJPVFRMRV9SRUFTT05fQ09OQ1VSUkVOVFNOQVBTSE9UQ1'
    'JFQVRFTElNSVRFWENFRURFRBAAEjcKM1RIUk9UVExFX1JFQVNPTl9GVU5DVElPTklOVk9DQVRJ'
    'T05SQVRFTElNSVRFWENFRURFRBABEjUKMVRIUk9UVExFX1JFQVNPTl9DT05DVVJSRU5USU5WT0'
    'NBVElPTkxJTUlURVhDRUVERUQQAhIrCidUSFJPVFRMRV9SRUFTT05fQ0FMTEVSUkFURUxJTUlU'
    'RVhDRUVERUQQAxJFCkFUSFJPVFRMRV9SRUFTT05fUkVTRVJWRURGVU5DVElPTkNPTkNVUlJFTl'
    'RJTlZPQ0FUSU9OTElNSVRFWENFRURFRBAEEj8KO1RIUk9UVExFX1JFQVNPTl9SRVNFUlZFREZV'
    'TkNUSU9OSU5WT0NBVElPTlJBVEVMSU1JVEVYQ0VFREVEEAU=');

@$core.Deprecated('Use tracingModeDescriptor instead')
const TracingMode$json = {
  '1': 'TracingMode',
  '2': [
    {'1': 'TRACING_MODE_ACTIVE', '2': 0},
    {'1': 'TRACING_MODE_PASSTHROUGH', '2': 1},
  ],
};

/// Descriptor for `TracingMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tracingModeDescriptor = $convert.base64Decode(
    'CgtUcmFjaW5nTW9kZRIXChNUUkFDSU5HX01PREVfQUNUSVZFEAASHAoYVFJBQ0lOR19NT0RFX1'
    'BBU1NUSFJPVUdIEAE=');

@$core.Deprecated('Use updateRuntimeOnDescriptor instead')
const UpdateRuntimeOn$json = {
  '1': 'UpdateRuntimeOn',
  '2': [
    {'1': 'UPDATE_RUNTIME_ON_MANUAL', '2': 0},
    {'1': 'UPDATE_RUNTIME_ON_AUTO', '2': 1},
    {'1': 'UPDATE_RUNTIME_ON_FUNCTIONUPDATE', '2': 2},
  ],
};

/// Descriptor for `UpdateRuntimeOn`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List updateRuntimeOnDescriptor = $convert.base64Decode(
    'Cg9VcGRhdGVSdW50aW1lT24SHAoYVVBEQVRFX1JVTlRJTUVfT05fTUFOVUFMEAASGgoWVVBEQV'
    'RFX1JVTlRJTUVfT05fQVVUTxABEiQKIFVQREFURV9SVU5USU1FX09OX0ZVTkNUSU9OVVBEQVRF'
    'EAI=');

@$core.Deprecated('Use accountLimitDescriptor instead')
const AccountLimit$json = {
  '1': 'AccountLimit',
  '2': [
    {
      '1': 'codesizeunzipped',
      '3': 508806555,
      '4': 1,
      '5': 3,
      '10': 'codesizeunzipped'
    },
    {
      '1': 'codesizezipped',
      '3': 26010534,
      '4': 1,
      '5': 3,
      '10': 'codesizezipped'
    },
    {
      '1': 'concurrentexecutions',
      '3': 333066772,
      '4': 1,
      '5': 5,
      '10': 'concurrentexecutions'
    },
    {
      '1': 'totalcodesize',
      '3': 202795702,
      '4': 1,
      '5': 3,
      '10': 'totalcodesize'
    },
    {
      '1': 'unreservedconcurrentexecutions',
      '3': 6163497,
      '4': 1,
      '5': 5,
      '10': 'unreservedconcurrentexecutions'
    },
  ],
};

/// Descriptor for `AccountLimit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountLimitDescriptor = $convert.base64Decode(
    'CgxBY2NvdW50TGltaXQSLgoQY29kZXNpemV1bnppcHBlZBibi8/yASABKANSEGNvZGVzaXpldW'
    '56aXBwZWQSKQoOY29kZXNpemV6aXBwZWQYpsezDCABKANSDmNvZGVzaXplemlwcGVkEjYKFGNv'
    'bmN1cnJlbnRleGVjdXRpb25zGJTk6J4BIAEoBVIUY29uY3VycmVudGV4ZWN1dGlvbnMSJwoNdG'
    '90YWxjb2Rlc2l6ZRi21dlgIAEoA1INdG90YWxjb2Rlc2l6ZRJJCh51bnJlc2VydmVkY29uY3Vy'
    'cmVudGV4ZWN1dGlvbnMYqZj4AiABKAVSHnVucmVzZXJ2ZWRjb25jdXJyZW50ZXhlY3V0aW9ucw'
    '==');

@$core.Deprecated('Use accountUsageDescriptor instead')
const AccountUsage$json = {
  '1': 'AccountUsage',
  '2': [
    {
      '1': 'functioncount',
      '3': 204403569,
      '4': 1,
      '5': 3,
      '10': 'functioncount'
    },
    {
      '1': 'totalcodesize',
      '3': 202795702,
      '4': 1,
      '5': 3,
      '10': 'totalcodesize'
    },
  ],
};

/// Descriptor for `AccountUsage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountUsageDescriptor = $convert.base64Decode(
    'CgxBY2NvdW50VXNhZ2USJwoNZnVuY3Rpb25jb3VudBjx5rthIAEoA1INZnVuY3Rpb25jb3VudB'
    'InCg10b3RhbGNvZGVzaXplGLbV2WAgASgDUg10b3RhbGNvZGVzaXpl');

@$core.Deprecated('Use addLayerVersionPermissionRequestDescriptor instead')
const AddLayerVersionPermissionRequest$json = {
  '1': 'AddLayerVersionPermissionRequest',
  '2': [
    {'1': 'action', '3': 175614240, '4': 1, '5': 9, '10': 'action'},
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
    {
      '1': 'organizationid',
      '3': 311006120,
      '4': 1,
      '5': 9,
      '10': 'organizationid'
    },
    {'1': 'principal', '3': 361640138, '4': 1, '5': 9, '10': 'principal'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'statementid', '3': 169602348, '4': 1, '5': 9, '10': 'statementid'},
    {
      '1': 'versionnumber',
      '3': 209346775,
      '4': 1,
      '5': 3,
      '10': 'versionnumber'
    },
  ],
};

/// Descriptor for `AddLayerVersionPermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addLayerVersionPermissionRequestDescriptor = $convert.base64Decode(
    'CiBBZGRMYXllclZlcnNpb25QZXJtaXNzaW9uUmVxdWVzdBIZCgZhY3Rpb24YoNLeUyABKAlSBm'
    'FjdGlvbhIgCglsYXllcm5hbWUYlqqY7QEgASgJUglsYXllcm5hbWUSKgoOb3JnYW5pemF0aW9u'
    'aWQYqKemlAEgASgJUg5vcmdhbml6YXRpb25pZBIgCglwcmluY2lwYWwYyuG4rAEgASgJUglwcm'
    'luY2lwYWwSIgoKcmV2aXNpb25pZBiGo57uASABKAlSCnJldmlzaW9uaWQSIwoLc3RhdGVtZW50'
    'aWQYrNrvUCABKAlSC3N0YXRlbWVudGlkEicKDXZlcnNpb25udW1iZXIY18HpYyABKANSDXZlcn'
    'Npb25udW1iZXI=');

@$core.Deprecated('Use addLayerVersionPermissionResponseDescriptor instead')
const AddLayerVersionPermissionResponse$json = {
  '1': 'AddLayerVersionPermissionResponse',
  '2': [
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'statement', '3': 248790199, '4': 1, '5': 9, '10': 'statement'},
  ],
};

/// Descriptor for `AddLayerVersionPermissionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addLayerVersionPermissionResponseDescriptor =
    $convert.base64Decode(
        'CiFBZGRMYXllclZlcnNpb25QZXJtaXNzaW9uUmVzcG9uc2USIgoKcmV2aXNpb25pZBiGo57uAS'
        'ABKAlSCnJldmlzaW9uaWQSHwoJc3RhdGVtZW50GLf50HYgASgJUglzdGF0ZW1lbnQ=');

@$core.Deprecated('Use addPermissionRequestDescriptor instead')
const AddPermissionRequest$json = {
  '1': 'AddPermissionRequest',
  '2': [
    {'1': 'action', '3': 175614240, '4': 1, '5': 9, '10': 'action'},
    {
      '1': 'eventsourcetoken',
      '3': 178305168,
      '4': 1,
      '5': 9,
      '10': 'eventsourcetoken'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'functionurlauthtype',
      '3': 417744441,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionUrlAuthType',
      '10': 'functionurlauthtype'
    },
    {
      '1': 'invokedviafunctionurl',
      '3': 425988239,
      '4': 1,
      '5': 8,
      '10': 'invokedviafunctionurl'
    },
    {'1': 'principal', '3': 361640138, '4': 1, '5': 9, '10': 'principal'},
    {
      '1': 'principalorgid',
      '3': 133988161,
      '4': 1,
      '5': 9,
      '10': 'principalorgid'
    },
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {
      '1': 'sourceaccount',
      '3': 118748296,
      '4': 1,
      '5': 9,
      '10': 'sourceaccount'
    },
    {'1': 'sourcearn', '3': 439903072, '4': 1, '5': 9, '10': 'sourcearn'},
    {'1': 'statementid', '3': 169602348, '4': 1, '5': 9, '10': 'statementid'},
  ],
};

/// Descriptor for `AddPermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addPermissionRequestDescriptor = $convert.base64Decode(
    'ChRBZGRQZXJtaXNzaW9uUmVxdWVzdBIZCgZhY3Rpb24YoNLeUyABKAlSBmFjdGlvbhItChBldm'
    'VudHNvdXJjZXRva2VuGJDxglUgASgJUhBldmVudHNvdXJjZXRva2VuEiYKDGZ1bmN0aW9ubmFt'
    'ZRijiL/fASABKAlSDGZ1bmN0aW9ubmFtZRJRChNmdW5jdGlvbnVybGF1dGh0eXBlGLmMmccBIA'
    'EoDjIbLmxhbWJkYS5GdW5jdGlvblVybEF1dGhUeXBlUhNmdW5jdGlvbnVybGF1dGh0eXBlEjgK'
    'FWludm9rZWR2aWFmdW5jdGlvbnVybBiPoZDLASABKAhSFWludm9rZWR2aWFmdW5jdGlvbnVybB'
    'IgCglwcmluY2lwYWwYyuG4rAEgASgJUglwcmluY2lwYWwSKQoOcHJpbmNpcGFsb3JnaWQYwf7x'
    'PyABKAlSDnByaW5jaXBhbG9yZ2lkEiAKCXF1YWxpZmllchjgtZH7ASABKAlSCXF1YWxpZmllch'
    'IiCgpyZXZpc2lvbmlkGIajnu4BIAEoCVIKcmV2aXNpb25pZBInCg1zb3VyY2VhY2NvdW50GIjp'
    'zzggASgJUg1zb3VyY2VhY2NvdW50EiAKCXNvdXJjZWFybhjgxuHRASABKAlSCXNvdXJjZWFybh'
    'IjCgtzdGF0ZW1lbnRpZBis2u9QIAEoCVILc3RhdGVtZW50aWQ=');

@$core.Deprecated('Use addPermissionResponseDescriptor instead')
const AddPermissionResponse$json = {
  '1': 'AddPermissionResponse',
  '2': [
    {'1': 'statement', '3': 248790199, '4': 1, '5': 9, '10': 'statement'},
  ],
};

/// Descriptor for `AddPermissionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addPermissionResponseDescriptor = $convert.base64Decode(
    'ChVBZGRQZXJtaXNzaW9uUmVzcG9uc2USHwoJc3RhdGVtZW50GLf50HYgASgJUglzdGF0ZW1lbn'
    'Q=');

@$core.Deprecated('Use aliasConfigurationDescriptor instead')
const AliasConfiguration$json = {
  '1': 'AliasConfiguration',
  '2': [
    {'1': 'aliasarn', '3': 461101595, '4': 1, '5': 9, '10': 'aliasarn'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'functionversion',
      '3': 365780244,
      '4': 1,
      '5': 9,
      '10': 'functionversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {
      '1': 'routingconfig',
      '3': 477346140,
      '4': 1,
      '5': 11,
      '6': '.lambda.AliasRoutingConfiguration',
      '10': 'routingconfig'
    },
  ],
};

/// Descriptor for `AliasConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aliasConfigurationDescriptor = $convert.base64Decode(
    'ChJBbGlhc0NvbmZpZ3VyYXRpb24SHgoIYWxpYXNhcm4Ym7Tv2wEgASgJUghhbGlhc2FybhIjCg'
    'tkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SLAoPZnVuY3Rpb252ZXJzaW9uGJS6'
    'ta4BIAEoCVIPZnVuY3Rpb252ZXJzaW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSIgoKcmV2aX'
    'Npb25pZBiGo57uASABKAlSCnJldmlzaW9uaWQSSwoNcm91dGluZ2NvbmZpZxjc8s7jASABKAsy'
    'IS5sYW1iZGEuQWxpYXNSb3V0aW5nQ29uZmlndXJhdGlvblINcm91dGluZ2NvbmZpZw==');

@$core.Deprecated('Use aliasRoutingConfigurationDescriptor instead')
const AliasRoutingConfiguration$json = {
  '1': 'AliasRoutingConfiguration',
  '2': [
    {
      '1': 'additionalversionweights',
      '3': 426806330,
      '4': 3,
      '5': 11,
      '6': '.lambda.AliasRoutingConfiguration.AdditionalversionweightsEntry',
      '10': 'additionalversionweights'
    },
  ],
  '3': [AliasRoutingConfiguration_AdditionalversionweightsEntry$json],
};

@$core.Deprecated('Use aliasRoutingConfigurationDescriptor instead')
const AliasRoutingConfiguration_AdditionalversionweightsEntry$json = {
  '1': 'AdditionalversionweightsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AliasRoutingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aliasRoutingConfigurationDescriptor = $convert.base64Decode(
    'ChlBbGlhc1JvdXRpbmdDb25maWd1cmF0aW9uEn8KGGFkZGl0aW9uYWx2ZXJzaW9ud2VpZ2h0cx'
    'i6mMLLASADKAsyPy5sYW1iZGEuQWxpYXNSb3V0aW5nQ29uZmlndXJhdGlvbi5BZGRpdGlvbmFs'
    'dmVyc2lvbndlaWdodHNFbnRyeVIYYWRkaXRpb25hbHZlcnNpb253ZWlnaHRzGksKHUFkZGl0aW'
    '9uYWx2ZXJzaW9ud2VpZ2h0c0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgB'
    'UgV2YWx1ZToCOAE=');

@$core.Deprecated('Use allowedPublishersDescriptor instead')
const AllowedPublishers$json = {
  '1': 'AllowedPublishers',
  '2': [
    {
      '1': 'signingprofileversionarns',
      '3': 187946552,
      '4': 3,
      '5': 9,
      '10': 'signingprofileversionarns'
    },
  ],
};

/// Descriptor for `AllowedPublishers`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List allowedPublishersDescriptor = $convert.base64Decode(
    'ChFBbGxvd2VkUHVibGlzaGVycxI/ChlzaWduaW5ncHJvZmlsZXZlcnNpb25hcm5zGLisz1kgAy'
    'gJUhlzaWduaW5ncHJvZmlsZXZlcnNpb25hcm5z');

@$core.Deprecated('Use amazonManagedKafkaEventSourceConfigDescriptor instead')
const AmazonManagedKafkaEventSourceConfig$json = {
  '1': 'AmazonManagedKafkaEventSourceConfig',
  '2': [
    {
      '1': 'consumergroupid',
      '3': 222398686,
      '4': 1,
      '5': 9,
      '10': 'consumergroupid'
    },
    {
      '1': 'schemaregistryconfig',
      '3': 482259768,
      '4': 1,
      '5': 11,
      '6': '.lambda.KafkaSchemaRegistryConfig',
      '10': 'schemaregistryconfig'
    },
  ],
};

/// Descriptor for `AmazonManagedKafkaEventSourceConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List amazonManagedKafkaEventSourceConfigDescriptor =
    $convert.base64Decode(
        'CiNBbWF6b25NYW5hZ2VkS2Fma2FFdmVudFNvdXJjZUNvbmZpZxIrCg9jb25zdW1lcmdyb3VwaW'
        'QY3pGGaiABKAlSD2NvbnN1bWVyZ3JvdXBpZBJZChRzY2hlbWFyZWdpc3RyeWNvbmZpZxi45vrl'
        'ASABKAsyIS5sYW1iZGEuS2Fma2FTY2hlbWFSZWdpc3RyeUNvbmZpZ1IUc2NoZW1hcmVnaXN0cn'
        'ljb25maWc=');

@$core.Deprecated('Use callbackDetailsDescriptor instead')
const CallbackDetails$json = {
  '1': 'CallbackDetails',
  '2': [
    {'1': 'callbackid', '3': 70101916, '4': 1, '5': 9, '10': 'callbackid'},
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '10': 'error'
    },
    {'1': 'result', '3': 273346629, '4': 1, '5': 9, '10': 'result'},
  ],
};

/// Descriptor for `CallbackDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callbackDetailsDescriptor = $convert.base64Decode(
    'Cg9DYWxsYmFja0RldGFpbHMSIQoKY2FsbGJhY2tpZBic17YhIAEoCVIKY2FsbGJhY2tpZBItCg'
    'VlcnJvchjyubacASABKAsyEy5sYW1iZGEuRXJyb3JPYmplY3RSBWVycm9yEhoKBnJlc3VsdBjF'
    '4KuCASABKAlSBnJlc3VsdA==');

@$core.Deprecated('Use callbackFailedDetailsDescriptor instead')
const CallbackFailedDetails$json = {
  '1': 'CallbackFailedDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `CallbackFailedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callbackFailedDetailsDescriptor = $convert.base64Decode(
    'ChVDYWxsYmFja0ZhaWxlZERldGFpbHMSLAoFZXJyb3IY8rm2nAEgASgLMhIubGFtYmRhLkV2ZW'
    '50RXJyb3JSBWVycm9y');

@$core.Deprecated('Use callbackOptionsDescriptor instead')
const CallbackOptions$json = {
  '1': 'CallbackOptions',
  '2': [
    {
      '1': 'heartbeattimeoutseconds',
      '3': 532445120,
      '4': 1,
      '5': 5,
      '10': 'heartbeattimeoutseconds'
    },
    {
      '1': 'timeoutseconds',
      '3': 336148022,
      '4': 1,
      '5': 5,
      '10': 'timeoutseconds'
    },
  ],
};

/// Descriptor for `CallbackOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callbackOptionsDescriptor = $convert.base64Decode(
    'Cg9DYWxsYmFja09wdGlvbnMSPAoXaGVhcnRiZWF0dGltZW91dHNlY29uZHMYwO/x/QEgASgFUh'
    'doZWFydGJlYXR0aW1lb3V0c2Vjb25kcxIqCg50aW1lb3V0c2Vjb25kcxi27KSgASABKAVSDnRp'
    'bWVvdXRzZWNvbmRz');

@$core.Deprecated('Use callbackStartedDetailsDescriptor instead')
const CallbackStartedDetails$json = {
  '1': 'CallbackStartedDetails',
  '2': [
    {'1': 'callbackid', '3': 70101916, '4': 1, '5': 9, '10': 'callbackid'},
    {
      '1': 'heartbeattimeout',
      '3': 35823659,
      '4': 1,
      '5': 5,
      '10': 'heartbeattimeout'
    },
    {'1': 'timeout', '3': 47808041, '4': 1, '5': 5, '10': 'timeout'},
  ],
};

/// Descriptor for `CallbackStartedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callbackStartedDetailsDescriptor = $convert.base64Decode(
    'ChZDYWxsYmFja1N0YXJ0ZWREZXRhaWxzEiEKCmNhbGxiYWNraWQYnNe2ISABKAlSCmNhbGxiYW'
    'NraWQSLQoQaGVhcnRiZWF0dGltZW91dBirwIoRIAEoBVIQaGVhcnRiZWF0dGltZW91dBIbCgd0'
    'aW1lb3V0GKn85RYgASgFUgd0aW1lb3V0');

@$core.Deprecated('Use callbackSucceededDetailsDescriptor instead')
const CallbackSucceededDetails$json = {
  '1': 'CallbackSucceededDetails',
  '2': [
    {
      '1': 'result',
      '3': 273346629,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventResult',
      '10': 'result'
    },
  ],
};

/// Descriptor for `CallbackSucceededDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callbackSucceededDetailsDescriptor =
    $convert.base64Decode(
        'ChhDYWxsYmFja1N1Y2NlZWRlZERldGFpbHMSLwoGcmVzdWx0GMXgq4IBIAEoCzITLmxhbWJkYS'
        '5FdmVudFJlc3VsdFIGcmVzdWx0');

@$core.Deprecated('Use callbackTimedOutDetailsDescriptor instead')
const CallbackTimedOutDetails$json = {
  '1': 'CallbackTimedOutDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `CallbackTimedOutDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callbackTimedOutDetailsDescriptor =
    $convert.base64Decode(
        'ChdDYWxsYmFja1RpbWVkT3V0RGV0YWlscxIsCgVlcnJvchjyubacASABKAsyEi5sYW1iZGEuRX'
        'ZlbnRFcnJvclIFZXJyb3I=');

@$core.Deprecated('Use callbackTimeoutExceptionDescriptor instead')
const CallbackTimeoutException$json = {
  '1': 'CallbackTimeoutException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `CallbackTimeoutException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List callbackTimeoutExceptionDescriptor =
    $convert.base64Decode(
        'ChhDYWxsYmFja1RpbWVvdXRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use capacityProviderDescriptor instead')
const CapacityProvider$json = {
  '1': 'CapacityProvider',
  '2': [
    {
      '1': 'capacityproviderarn',
      '3': 109580344,
      '4': 1,
      '5': 9,
      '10': 'capacityproviderarn'
    },
    {
      '1': 'capacityproviderscalingconfig',
      '3': 499413336,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderScalingConfig',
      '10': 'capacityproviderscalingconfig'
    },
    {
      '1': 'instancerequirements',
      '3': 520528435,
      '4': 1,
      '5': 11,
      '6': '.lambda.InstanceRequirements',
      '10': 'instancerequirements'
    },
    {'1': 'kmskeyarn', '3': 110041649, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'permissionsconfig',
      '3': 249189794,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderPermissionsConfig',
      '10': 'permissionsconfig'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.lambda.CapacityProviderState',
      '10': 'state'
    },
    {
      '1': 'vpcconfig',
      '3': 194980743,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderVpcConfig',
      '10': 'vpcconfig'
    },
  ],
};

/// Descriptor for `CapacityProvider`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderDescriptor = $convert.base64Decode(
    'ChBDYXBhY2l0eVByb3ZpZGVyEjMKE2NhcGFjaXR5cHJvdmlkZXJhcm4YuKCgNCABKAlSE2NhcG'
    'FjaXR5cHJvdmlkZXJhcm4SbwodY2FwYWNpdHlwcm92aWRlcnNjYWxpbmdjb25maWcY2OKR7gEg'
    'ASgLMiUubGFtYmRhLkNhcGFjaXR5UHJvdmlkZXJTY2FsaW5nQ29uZmlnUh1jYXBhY2l0eXByb3'
    'ZpZGVyc2NhbGluZ2NvbmZpZxJUChRpbnN0YW5jZXJlcXVpcmVtZW50cxizxJr4ASABKAsyHC5s'
    'YW1iZGEuSW5zdGFuY2VSZXF1aXJlbWVudHNSFGluc3RhbmNlcmVxdWlyZW1lbnRzEh8KCWttc2'
    'tleWFybhixtLw0IAEoCVIJa21za2V5YXJuEiYKDGxhc3Rtb2RpZmllZBinnPzOASABKAlSDGxh'
    'c3Rtb2RpZmllZBJaChFwZXJtaXNzaW9uc2NvbmZpZxiiq+l2IAEoCzIpLmxhbWJkYS5DYXBhY2'
    'l0eVByb3ZpZGVyUGVybWlzc2lvbnNDb25maWdSEXBlcm1pc3Npb25zY29uZmlnEjcKBXN0YXRl'
    'GJfJsu8BIAEoDjIdLmxhbWJkYS5DYXBhY2l0eVByb3ZpZGVyU3RhdGVSBXN0YXRlEkIKCXZwY2'
    'NvbmZpZxiH1/xcIAEoCzIhLmxhbWJkYS5DYXBhY2l0eVByb3ZpZGVyVnBjQ29uZmlnUgl2cGNj'
    'b25maWc=');

@$core.Deprecated('Use capacityProviderConfigDescriptor instead')
const CapacityProviderConfig$json = {
  '1': 'CapacityProviderConfig',
  '2': [
    {
      '1': 'lambdamanagedinstancescapacityproviderconfig',
      '3': 297153183,
      '4': 1,
      '5': 11,
      '6': '.lambda.LambdaManagedInstancesCapacityProviderConfig',
      '10': 'lambdamanagedinstancescapacityproviderconfig'
    },
  ],
};

/// Descriptor for `CapacityProviderConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderConfigDescriptor = $convert.base64Decode(
    'ChZDYXBhY2l0eVByb3ZpZGVyQ29uZmlnEpwBCixsYW1iZGFtYW5hZ2VkaW5zdGFuY2VzY2FwYW'
    'NpdHlwcm92aWRlcmNvbmZpZxif5diNASABKAsyNC5sYW1iZGEuTGFtYmRhTWFuYWdlZEluc3Rh'
    'bmNlc0NhcGFjaXR5UHJvdmlkZXJDb25maWdSLGxhbWJkYW1hbmFnZWRpbnN0YW5jZXNjYXBhY2'
    'l0eXByb3ZpZGVyY29uZmln');

@$core
    .Deprecated('Use capacityProviderLimitExceededExceptionDescriptor instead')
const CapacityProviderLimitExceededException$json = {
  '1': 'CapacityProviderLimitExceededException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CapacityProviderLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiZDYXBhY2l0eVByb3ZpZGVyTGltaXRFeGNlZWRlZEV4Y2VwdGlvbhIWCgR0eXBlGO6g14oBIA'
        'EoCVIEdHlwZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use capacityProviderPermissionsConfigDescriptor instead')
const CapacityProviderPermissionsConfig$json = {
  '1': 'CapacityProviderPermissionsConfig',
  '2': [
    {
      '1': 'capacityprovideroperatorrolearn',
      '3': 320252472,
      '4': 1,
      '5': 9,
      '10': 'capacityprovideroperatorrolearn'
    },
  ],
};

/// Descriptor for `CapacityProviderPermissionsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderPermissionsConfigDescriptor =
    $convert.base64Decode(
        'CiFDYXBhY2l0eVByb3ZpZGVyUGVybWlzc2lvbnNDb25maWcSTAofY2FwYWNpdHlwcm92aWRlcm'
        '9wZXJhdG9ycm9sZWFybhi41NqYASABKAlSH2NhcGFjaXR5cHJvdmlkZXJvcGVyYXRvcnJvbGVh'
        'cm4=');

@$core.Deprecated('Use capacityProviderScalingConfigDescriptor instead')
const CapacityProviderScalingConfig$json = {
  '1': 'CapacityProviderScalingConfig',
  '2': [
    {'1': 'maxvcpucount', '3': 336271977, '4': 1, '5': 5, '10': 'maxvcpucount'},
    {
      '1': 'scalingmode',
      '3': 210356138,
      '4': 1,
      '5': 14,
      '6': '.lambda.CapacityProviderScalingMode',
      '10': 'scalingmode'
    },
    {
      '1': 'scalingpolicies',
      '3': 289494257,
      '4': 3,
      '5': 11,
      '6': '.lambda.TargetTrackingScalingPolicy',
      '10': 'scalingpolicies'
    },
  ],
};

/// Descriptor for `CapacityProviderScalingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderScalingConfigDescriptor = $convert.base64Decode(
    'Ch1DYXBhY2l0eVByb3ZpZGVyU2NhbGluZ0NvbmZpZxImCgxtYXh2Y3B1Y291bnQY6bSsoAEgAS'
    'gFUgxtYXh2Y3B1Y291bnQSSAoLc2NhbGluZ21vZGUYqo+nZCABKA4yIy5sYW1iZGEuQ2FwYWNp'
    'dHlQcm92aWRlclNjYWxpbmdNb2RlUgtzY2FsaW5nbW9kZRJRCg9zY2FsaW5ncG9saWNpZXMY8a'
    'mFigEgAygLMiMubGFtYmRhLlRhcmdldFRyYWNraW5nU2NhbGluZ1BvbGljeVIPc2NhbGluZ3Bv'
    'bGljaWVz');

@$core.Deprecated('Use capacityProviderVpcConfigDescriptor instead')
const CapacityProviderVpcConfig$json = {
  '1': 'CapacityProviderVpcConfig',
  '2': [
    {
      '1': 'securitygroupids',
      '3': 13337861,
      '4': 3,
      '5': 9,
      '10': 'securitygroupids'
    },
    {'1': 'subnetids', '3': 266219411, '4': 3, '5': 9, '10': 'subnetids'},
  ],
};

/// Descriptor for `CapacityProviderVpcConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderVpcConfigDescriptor =
    $convert.base64Decode(
        'ChlDYXBhY2l0eVByb3ZpZGVyVnBjQ29uZmlnEi0KEHNlY3VyaXR5Z3JvdXBpZHMYhYquBiADKA'
        'lSEHNlY3VyaXR5Z3JvdXBpZHMSHwoJc3VibmV0aWRzGJPf+H4gAygJUglzdWJuZXRpZHM=');

@$core.Deprecated('Use chainedInvokeDetailsDescriptor instead')
const ChainedInvokeDetails$json = {
  '1': 'ChainedInvokeDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '10': 'error'
    },
    {'1': 'result', '3': 273346629, '4': 1, '5': 9, '10': 'result'},
  ],
};

/// Descriptor for `ChainedInvokeDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chainedInvokeDetailsDescriptor = $convert.base64Decode(
    'ChRDaGFpbmVkSW52b2tlRGV0YWlscxItCgVlcnJvchjyubacASABKAsyEy5sYW1iZGEuRXJyb3'
    'JPYmplY3RSBWVycm9yEhoKBnJlc3VsdBjF4KuCASABKAlSBnJlc3VsdA==');

@$core.Deprecated('Use chainedInvokeFailedDetailsDescriptor instead')
const ChainedInvokeFailedDetails$json = {
  '1': 'ChainedInvokeFailedDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `ChainedInvokeFailedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chainedInvokeFailedDetailsDescriptor =
    $convert.base64Decode(
        'ChpDaGFpbmVkSW52b2tlRmFpbGVkRGV0YWlscxIsCgVlcnJvchjyubacASABKAsyEi5sYW1iZG'
        'EuRXZlbnRFcnJvclIFZXJyb3I=');

@$core.Deprecated('Use chainedInvokeOptionsDescriptor instead')
const ChainedInvokeOptions$json = {
  '1': 'ChainedInvokeOptions',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
  ],
};

/// Descriptor for `ChainedInvokeOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chainedInvokeOptionsDescriptor = $convert.base64Decode(
    'ChRDaGFpbmVkSW52b2tlT3B0aW9ucxImCgxmdW5jdGlvbm5hbWUYo4i/3wEgASgJUgxmdW5jdG'
    'lvbm5hbWUSHQoIdGVuYW50aWQYgfjvZCABKAlSCHRlbmFudGlk');

@$core.Deprecated('Use chainedInvokeStartedDetailsDescriptor instead')
const ChainedInvokeStartedDetails$json = {
  '1': 'ChainedInvokeStartedDetails',
  '2': [
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {
      '1': 'executedversion',
      '3': 457672809,
      '4': 1,
      '5': 9,
      '10': 'executedversion'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'input',
      '3': 529785116,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventInput',
      '10': 'input'
    },
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
  ],
};

/// Descriptor for `ChainedInvokeStartedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chainedInvokeStartedDetailsDescriptor = $convert.base64Decode(
    'ChtDaGFpbmVkSW52b2tlU3RhcnRlZERldGFpbHMSNAoTZHVyYWJsZWV4ZWN1dGlvbmFybhiKzb'
    'eAASABKAlSE2R1cmFibGVleGVjdXRpb25hcm4SLAoPZXhlY3V0ZWR2ZXJzaW9uGOmQntoBIAEo'
    'CVIPZXhlY3V0ZWR2ZXJzaW9uEiYKDGZ1bmN0aW9ubmFtZRijiL/fASABKAlSDGZ1bmN0aW9ubm'
    'FtZRIsCgVpbnB1dBicws/8ASABKAsyEi5sYW1iZGEuRXZlbnRJbnB1dFIFaW5wdXQSHQoIdGVu'
    'YW50aWQYgfjvZCABKAlSCHRlbmFudGlk');

@$core.Deprecated('Use chainedInvokeStoppedDetailsDescriptor instead')
const ChainedInvokeStoppedDetails$json = {
  '1': 'ChainedInvokeStoppedDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `ChainedInvokeStoppedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chainedInvokeStoppedDetailsDescriptor =
    $convert.base64Decode(
        'ChtDaGFpbmVkSW52b2tlU3RvcHBlZERldGFpbHMSLAoFZXJyb3IY8rm2nAEgASgLMhIubGFtYm'
        'RhLkV2ZW50RXJyb3JSBWVycm9y');

@$core.Deprecated('Use chainedInvokeSucceededDetailsDescriptor instead')
const ChainedInvokeSucceededDetails$json = {
  '1': 'ChainedInvokeSucceededDetails',
  '2': [
    {
      '1': 'result',
      '3': 273346629,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventResult',
      '10': 'result'
    },
  ],
};

/// Descriptor for `ChainedInvokeSucceededDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chainedInvokeSucceededDetailsDescriptor =
    $convert.base64Decode(
        'Ch1DaGFpbmVkSW52b2tlU3VjY2VlZGVkRGV0YWlscxIvCgZyZXN1bHQYxeCrggEgASgLMhMubG'
        'FtYmRhLkV2ZW50UmVzdWx0UgZyZXN1bHQ=');

@$core.Deprecated('Use chainedInvokeTimedOutDetailsDescriptor instead')
const ChainedInvokeTimedOutDetails$json = {
  '1': 'ChainedInvokeTimedOutDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `ChainedInvokeTimedOutDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chainedInvokeTimedOutDetailsDescriptor =
    $convert.base64Decode(
        'ChxDaGFpbmVkSW52b2tlVGltZWRPdXREZXRhaWxzEiwKBWVycm9yGPK5tpwBIAEoCzISLmxhbW'
        'JkYS5FdmVudEVycm9yUgVlcnJvcg==');

@$core.Deprecated('Use checkpointDurableExecutionRequestDescriptor instead')
const CheckpointDurableExecutionRequest$json = {
  '1': 'CheckpointDurableExecutionRequest',
  '2': [
    {
      '1': 'checkpointtoken',
      '3': 27028053,
      '4': 1,
      '5': 9,
      '10': 'checkpointtoken'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {
      '1': 'updates',
      '3': 137393574,
      '4': 3,
      '5': 11,
      '6': '.lambda.OperationUpdate',
      '10': 'updates'
    },
  ],
};

/// Descriptor for `CheckpointDurableExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkpointDurableExecutionRequestDescriptor = $convert.base64Decode(
    'CiFDaGVja3BvaW50RHVyYWJsZUV4ZWN1dGlvblJlcXVlc3QSKwoPY2hlY2twb2ludHRva2VuGN'
    'XU8QwgASgJUg9jaGVja3BvaW50dG9rZW4SIwoLY2xpZW50dG9rZW4YzPu7QSABKAlSC2NsaWVu'
    'dHRva2VuEjQKE2R1cmFibGVleGVjdXRpb25hcm4Yis23gAEgASgJUhNkdXJhYmxlZXhlY3V0aW'
    '9uYXJuEjQKB3VwZGF0ZXMYpuvBQSADKAsyFy5sYW1iZGEuT3BlcmF0aW9uVXBkYXRlUgd1cGRh'
    'dGVz');

@$core.Deprecated('Use checkpointDurableExecutionResponseDescriptor instead')
const CheckpointDurableExecutionResponse$json = {
  '1': 'CheckpointDurableExecutionResponse',
  '2': [
    {
      '1': 'checkpointtoken',
      '3': 27028053,
      '4': 1,
      '5': 9,
      '10': 'checkpointtoken'
    },
    {
      '1': 'newexecutionstate',
      '3': 298388571,
      '4': 1,
      '5': 11,
      '6': '.lambda.CheckpointUpdatedExecutionState',
      '10': 'newexecutionstate'
    },
  ],
};

/// Descriptor for `CheckpointDurableExecutionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkpointDurableExecutionResponseDescriptor =
    $convert.base64Decode(
        'CiJDaGVja3BvaW50RHVyYWJsZUV4ZWN1dGlvblJlc3BvbnNlEisKD2NoZWNrcG9pbnR0b2tlbh'
        'jV1PEMIAEoCVIPY2hlY2twb2ludHRva2VuElkKEW5ld2V4ZWN1dGlvbnN0YXRlGNuYpI4BIAEo'
        'CzInLmxhbWJkYS5DaGVja3BvaW50VXBkYXRlZEV4ZWN1dGlvblN0YXRlUhFuZXdleGVjdXRpb2'
        '5zdGF0ZQ==');

@$core.Deprecated('Use checkpointUpdatedExecutionStateDescriptor instead')
const CheckpointUpdatedExecutionState$json = {
  '1': 'CheckpointUpdatedExecutionState',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {
      '1': 'operations',
      '3': 126776656,
      '4': 3,
      '5': 11,
      '6': '.lambda.Operation',
      '10': 'operations'
    },
  ],
};

/// Descriptor for `CheckpointUpdatedExecutionState`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkpointUpdatedExecutionStateDescriptor =
    $convert.base64Decode(
        'Ch9DaGVja3BvaW50VXBkYXRlZEV4ZWN1dGlvblN0YXRlEiIKCm5leHRtYXJrZXIYo4Gu/QEgAS'
        'gJUgpuZXh0bWFya2VyEjQKCm9wZXJhdGlvbnMY0Oq5PCADKAsyES5sYW1iZGEuT3BlcmF0aW9u'
        'UgpvcGVyYXRpb25z');

@$core.Deprecated('Use codeSigningConfigDescriptor instead')
const CodeSigningConfig$json = {
  '1': 'CodeSigningConfig',
  '2': [
    {
      '1': 'allowedpublishers',
      '3': 217933299,
      '4': 1,
      '5': 11,
      '6': '.lambda.AllowedPublishers',
      '10': 'allowedpublishers'
    },
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
    {
      '1': 'codesigningconfigid',
      '3': 77491229,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigid'
    },
    {
      '1': 'codesigningpolicies',
      '3': 254459202,
      '4': 1,
      '5': 11,
      '6': '.lambda.CodeSigningPolicies',
      '10': 'codesigningpolicies'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
  ],
};

/// Descriptor for `CodeSigningConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeSigningConfigDescriptor = $convert.base64Decode(
    'ChFDb2RlU2lnbmluZ0NvbmZpZxJKChFhbGxvd2VkcHVibGlzaGVycxjzy/VnIAEoCzIZLmxhbW'
    'JkYS5BbGxvd2VkUHVibGlzaGVyc1IRYWxsb3dlZHB1Ymxpc2hlcnMSNgoUY29kZXNpZ25pbmdj'
    'b25maWdhcm4Ywfz38AEgASgJUhRjb2Rlc2lnbmluZ2NvbmZpZ2FybhIzChNjb2Rlc2lnbmluZ2'
    'NvbmZpZ2lkGJ3Y+SQgASgJUhNjb2Rlc2lnbmluZ2NvbmZpZ2lkElAKE2NvZGVzaWduaW5ncG9s'
    'aWNpZXMYwvqqeSABKAsyGy5sYW1iZGEuQ29kZVNpZ25pbmdQb2xpY2llc1ITY29kZXNpZ25pbm'
    'dwb2xpY2llcxIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SJgoMbGFzdG1v'
    'ZGlmaWVkGKec/M4BIAEoCVIMbGFzdG1vZGlmaWVk');

@$core.Deprecated('Use codeSigningConfigNotFoundExceptionDescriptor instead')
const CodeSigningConfigNotFoundException$json = {
  '1': 'CodeSigningConfigNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `CodeSigningConfigNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeSigningConfigNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CiJDb2RlU2lnbmluZ0NvbmZpZ05vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKA'
        'lSB21lc3NhZ2USFgoEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use codeSigningPoliciesDescriptor instead')
const CodeSigningPolicies$json = {
  '1': 'CodeSigningPolicies',
  '2': [
    {
      '1': 'untrustedartifactondeployment',
      '3': 67707554,
      '4': 1,
      '5': 14,
      '6': '.lambda.CodeSigningPolicy',
      '10': 'untrustedartifactondeployment'
    },
  ],
};

/// Descriptor for `CodeSigningPolicies`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeSigningPoliciesDescriptor = $convert.base64Decode(
    'ChNDb2RlU2lnbmluZ1BvbGljaWVzEmIKHXVudHJ1c3RlZGFydGlmYWN0b25kZXBsb3ltZW50GK'
    'LFpCAgASgOMhkubGFtYmRhLkNvZGVTaWduaW5nUG9saWN5Uh11bnRydXN0ZWRhcnRpZmFjdG9u'
    'ZGVwbG95bWVudA==');

@$core.Deprecated('Use codeStorageExceededExceptionDescriptor instead')
const CodeStorageExceededException$json = {
  '1': 'CodeStorageExceededException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CodeStorageExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeStorageExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChxDb2RlU3RvcmFnZUV4Y2VlZGVkRXhjZXB0aW9uEhYKBHR5cGUY7qDXigEgASgJUgR0eXBlEh'
        'sKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use codeVerificationFailedExceptionDescriptor instead')
const CodeVerificationFailedException$json = {
  '1': 'CodeVerificationFailedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `CodeVerificationFailedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List codeVerificationFailedExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Db2RlVmVyaWZpY2F0aW9uRmFpbGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2USFgoEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use concurrencyDescriptor instead')
const Concurrency$json = {
  '1': 'Concurrency',
  '2': [
    {
      '1': 'reservedconcurrentexecutions',
      '3': 40149212,
      '4': 1,
      '5': 5,
      '10': 'reservedconcurrentexecutions'
    },
  ],
};

/// Descriptor for `Concurrency`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrencyDescriptor = $convert.base64Decode(
    'CgtDb25jdXJyZW5jeRJFChxyZXNlcnZlZGNvbmN1cnJlbnRleGVjdXRpb25zGNzBkhMgASgFUh'
    'xyZXNlcnZlZGNvbmN1cnJlbnRleGVjdXRpb25z');

@$core.Deprecated('Use contextDetailsDescriptor instead')
const ContextDetails$json = {
  '1': 'ContextDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '10': 'error'
    },
    {
      '1': 'replaychildren',
      '3': 146095102,
      '4': 1,
      '5': 8,
      '10': 'replaychildren'
    },
    {'1': 'result', '3': 273346629, '4': 1, '5': 9, '10': 'result'},
  ],
};

/// Descriptor for `ContextDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextDetailsDescriptor = $convert.base64Decode(
    'Cg5Db250ZXh0RGV0YWlscxItCgVlcnJvchjyubacASABKAsyEy5sYW1iZGEuRXJyb3JPYmplY3'
    'RSBWVycm9yEikKDnJlcGxheWNoaWxkcmVuGP731EUgASgIUg5yZXBsYXljaGlsZHJlbhIaCgZy'
    'ZXN1bHQYxeCrggEgASgJUgZyZXN1bHQ=');

@$core.Deprecated('Use contextFailedDetailsDescriptor instead')
const ContextFailedDetails$json = {
  '1': 'ContextFailedDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `ContextFailedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextFailedDetailsDescriptor = $convert.base64Decode(
    'ChRDb250ZXh0RmFpbGVkRGV0YWlscxIsCgVlcnJvchjyubacASABKAsyEi5sYW1iZGEuRXZlbn'
    'RFcnJvclIFZXJyb3I=');

@$core.Deprecated('Use contextOptionsDescriptor instead')
const ContextOptions$json = {
  '1': 'ContextOptions',
  '2': [
    {
      '1': 'replaychildren',
      '3': 146095102,
      '4': 1,
      '5': 8,
      '10': 'replaychildren'
    },
  ],
};

/// Descriptor for `ContextOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextOptionsDescriptor = $convert.base64Decode(
    'Cg5Db250ZXh0T3B0aW9ucxIpCg5yZXBsYXljaGlsZHJlbhj+99RFIAEoCFIOcmVwbGF5Y2hpbG'
    'RyZW4=');

@$core.Deprecated('Use contextStartedDetailsDescriptor instead')
const ContextStartedDetails$json = {
  '1': 'ContextStartedDetails',
};

/// Descriptor for `ContextStartedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextStartedDetailsDescriptor =
    $convert.base64Decode('ChVDb250ZXh0U3RhcnRlZERldGFpbHM=');

@$core.Deprecated('Use contextSucceededDetailsDescriptor instead')
const ContextSucceededDetails$json = {
  '1': 'ContextSucceededDetails',
  '2': [
    {
      '1': 'result',
      '3': 273346629,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventResult',
      '10': 'result'
    },
  ],
};

/// Descriptor for `ContextSucceededDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextSucceededDetailsDescriptor =
    $convert.base64Decode(
        'ChdDb250ZXh0U3VjY2VlZGVkRGV0YWlscxIvCgZyZXN1bHQYxeCrggEgASgLMhMubGFtYmRhLk'
        'V2ZW50UmVzdWx0UgZyZXN1bHQ=');

@$core.Deprecated('Use corsDescriptor instead')
const Cors$json = {
  '1': 'Cors',
  '2': [
    {
      '1': 'allowcredentials',
      '3': 487229205,
      '4': 1,
      '5': 8,
      '10': 'allowcredentials'
    },
    {'1': 'allowheaders', '3': 511432073, '4': 3, '5': 9, '10': 'allowheaders'},
    {'1': 'allowmethods', '3': 441007253, '4': 3, '5': 9, '10': 'allowmethods'},
    {'1': 'alloworigins', '3': 343821862, '4': 3, '5': 9, '10': 'alloworigins'},
    {
      '1': 'exposeheaders',
      '3': 290364554,
      '4': 3,
      '5': 9,
      '10': 'exposeheaders'
    },
    {'1': 'maxage', '3': 243937897, '4': 1, '5': 5, '10': 'maxage'},
  ],
};

/// Descriptor for `Cors`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List corsDescriptor = $convert.base64Decode(
    'CgRDb3JzEi4KEGFsbG93Y3JlZGVudGlhbHMYlY6q6AEgASgIUhBhbGxvd2NyZWRlbnRpYWxzEi'
    'YKDGFsbG93aGVhZGVycxiJq+/zASADKAlSDGFsbG93aGVhZGVycxImCgxhbGxvd21ldGhvZHMY'
    'lfmk0gEgAygJUgxhbGxvd21ldGhvZHMSJgoMYWxsb3dvcmlnaW5zGKac+aMBIAMoCVIMYWxsb3'
    'dvcmlnaW5zEigKDWV4cG9zZWhlYWRlcnMYirm6igEgAygJUg1leHBvc2VoZWFkZXJzEhkKBm1h'
    'eGFnZRjp5Kh0IAEoBVIGbWF4YWdl');

@$core.Deprecated('Use createAliasRequestDescriptor instead')
const CreateAliasRequest$json = {
  '1': 'CreateAliasRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'functionversion',
      '3': 365780244,
      '4': 1,
      '5': 9,
      '10': 'functionversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'routingconfig',
      '3': 477346140,
      '4': 1,
      '5': 11,
      '6': '.lambda.AliasRoutingConfiguration',
      '10': 'routingconfig'
    },
  ],
};

/// Descriptor for `CreateAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAliasRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVBbGlhc1JlcXVlc3QSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW'
    '9uEiYKDGZ1bmN0aW9ubmFtZRijiL/fASABKAlSDGZ1bmN0aW9ubmFtZRIsCg9mdW5jdGlvbnZl'
    'cnNpb24YlLq1rgEgASgJUg9mdW5jdGlvbnZlcnNpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZR'
    'JLCg1yb3V0aW5nY29uZmlnGNzyzuMBIAEoCzIhLmxhbWJkYS5BbGlhc1JvdXRpbmdDb25maWd1'
    'cmF0aW9uUg1yb3V0aW5nY29uZmln');

@$core.Deprecated('Use createCapacityProviderRequestDescriptor instead')
const CreateCapacityProviderRequest$json = {
  '1': 'CreateCapacityProviderRequest',
  '2': [
    {
      '1': 'capacityprovidername',
      '3': 466230132,
      '4': 1,
      '5': 9,
      '10': 'capacityprovidername'
    },
    {
      '1': 'capacityproviderscalingconfig',
      '3': 499413336,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderScalingConfig',
      '10': 'capacityproviderscalingconfig'
    },
    {
      '1': 'instancerequirements',
      '3': 520528435,
      '4': 1,
      '5': 11,
      '6': '.lambda.InstanceRequirements',
      '10': 'instancerequirements'
    },
    {'1': 'kmskeyarn', '3': 110041649, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {
      '1': 'permissionsconfig',
      '3': 249189794,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderPermissionsConfig',
      '10': 'permissionsconfig'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.lambda.CreateCapacityProviderRequest.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'vpcconfig',
      '3': 194980743,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderVpcConfig',
      '10': 'vpcconfig'
    },
  ],
  '3': [CreateCapacityProviderRequest_TagsEntry$json],
};

@$core.Deprecated('Use createCapacityProviderRequestDescriptor instead')
const CreateCapacityProviderRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateCapacityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCapacityProviderRequestDescriptor = $convert.base64Decode(
    'Ch1DcmVhdGVDYXBhY2l0eVByb3ZpZGVyUmVxdWVzdBI2ChRjYXBhY2l0eXByb3ZpZGVybmFtZR'
    'j0tqjeASABKAlSFGNhcGFjaXR5cHJvdmlkZXJuYW1lEm8KHWNhcGFjaXR5cHJvdmlkZXJzY2Fs'
    'aW5nY29uZmlnGNjike4BIAEoCzIlLmxhbWJkYS5DYXBhY2l0eVByb3ZpZGVyU2NhbGluZ0Nvbm'
    'ZpZ1IdY2FwYWNpdHlwcm92aWRlcnNjYWxpbmdjb25maWcSVAoUaW5zdGFuY2VyZXF1aXJlbWVu'
    'dHMYs8Sa+AEgASgLMhwubGFtYmRhLkluc3RhbmNlUmVxdWlyZW1lbnRzUhRpbnN0YW5jZXJlcX'
    'VpcmVtZW50cxIfCglrbXNrZXlhcm4YsbS8NCABKAlSCWttc2tleWFybhJaChFwZXJtaXNzaW9u'
    'c2NvbmZpZxiiq+l2IAEoCzIpLmxhbWJkYS5DYXBhY2l0eVByb3ZpZGVyUGVybWlzc2lvbnNDb2'
    '5maWdSEXBlcm1pc3Npb25zY29uZmlnEkcKBHRhZ3MYwcH2tQEgAygLMi8ubGFtYmRhLkNyZWF0'
    'ZUNhcGFjaXR5UHJvdmlkZXJSZXF1ZXN0LlRhZ3NFbnRyeVIEdGFncxJCCgl2cGNjb25maWcYh9'
    'f8XCABKAsyIS5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlclZwY0NvbmZpZ1IJdnBjY29uZmlnGjcK'
    'CVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use createCapacityProviderResponseDescriptor instead')
const CreateCapacityProviderResponse$json = {
  '1': 'CreateCapacityProviderResponse',
  '2': [
    {
      '1': 'capacityprovider',
      '3': 448074961,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProvider',
      '10': 'capacityprovider'
    },
  ],
};

/// Descriptor for `CreateCapacityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCapacityProviderResponseDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVDYXBhY2l0eVByb3ZpZGVyUmVzcG9uc2USSAoQY2FwYWNpdHlwcm92aWRlchjRqd'
        'TVASABKAsyGC5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlclIQY2FwYWNpdHlwcm92aWRlcg==');

@$core.Deprecated('Use createCodeSigningConfigRequestDescriptor instead')
const CreateCodeSigningConfigRequest$json = {
  '1': 'CreateCodeSigningConfigRequest',
  '2': [
    {
      '1': 'allowedpublishers',
      '3': 217933299,
      '4': 1,
      '5': 11,
      '6': '.lambda.AllowedPublishers',
      '10': 'allowedpublishers'
    },
    {
      '1': 'codesigningpolicies',
      '3': 254459202,
      '4': 1,
      '5': 11,
      '6': '.lambda.CodeSigningPolicies',
      '10': 'codesigningpolicies'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.lambda.CreateCodeSigningConfigRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [CreateCodeSigningConfigRequest_TagsEntry$json],
};

@$core.Deprecated('Use createCodeSigningConfigRequestDescriptor instead')
const CreateCodeSigningConfigRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCodeSigningConfigRequestDescriptor = $convert.base64Decode(
    'Ch5DcmVhdGVDb2RlU2lnbmluZ0NvbmZpZ1JlcXVlc3QSSgoRYWxsb3dlZHB1Ymxpc2hlcnMY88'
    'v1ZyABKAsyGS5sYW1iZGEuQWxsb3dlZFB1Ymxpc2hlcnNSEWFsbG93ZWRwdWJsaXNoZXJzElAK'
    'E2NvZGVzaWduaW5ncG9saWNpZXMYwvqqeSABKAsyGy5sYW1iZGEuQ29kZVNpZ25pbmdQb2xpY2'
    'llc1ITY29kZXNpZ25pbmdwb2xpY2llcxIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3Jp'
    'cHRpb24SSAoEdGFncxjBwfa1ASADKAsyMC5sYW1iZGEuQ3JlYXRlQ29kZVNpZ25pbmdDb25maW'
    'dSZXF1ZXN0LlRhZ3NFbnRyeVIEdGFncxo3CglUYWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkS'
    'FAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createCodeSigningConfigResponseDescriptor instead')
const CreateCodeSigningConfigResponse$json = {
  '1': 'CreateCodeSigningConfigResponse',
  '2': [
    {
      '1': 'codesigningconfig',
      '3': 130458490,
      '4': 1,
      '5': 11,
      '6': '.lambda.CodeSigningConfig',
      '10': 'codesigningconfig'
    },
  ],
};

/// Descriptor for `CreateCodeSigningConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCodeSigningConfigResponseDescriptor =
    $convert.base64Decode(
        'Ch9DcmVhdGVDb2RlU2lnbmluZ0NvbmZpZ1Jlc3BvbnNlEkoKEWNvZGVzaWduaW5nY29uZmlnGP'
        'rGmj4gASgLMhkubGFtYmRhLkNvZGVTaWduaW5nQ29uZmlnUhFjb2Rlc2lnbmluZ2NvbmZpZw==');

@$core.Deprecated('Use createEventSourceMappingRequestDescriptor instead')
const CreateEventSourceMappingRequest$json = {
  '1': 'CreateEventSourceMappingRequest',
  '2': [
    {
      '1': 'amazonmanagedkafkaeventsourceconfig',
      '3': 60136380,
      '4': 1,
      '5': 11,
      '6': '.lambda.AmazonManagedKafkaEventSourceConfig',
      '10': 'amazonmanagedkafkaeventsourceconfig'
    },
    {'1': 'batchsize', '3': 318039259, '4': 1, '5': 5, '10': 'batchsize'},
    {
      '1': 'bisectbatchonfunctionerror',
      '3': 276143707,
      '4': 1,
      '5': 8,
      '10': 'bisectbatchonfunctionerror'
    },
    {
      '1': 'destinationconfig',
      '3': 184834158,
      '4': 1,
      '5': 11,
      '6': '.lambda.DestinationConfig',
      '10': 'destinationconfig'
    },
    {
      '1': 'documentdbeventsourceconfig',
      '3': 173060622,
      '4': 1,
      '5': 11,
      '6': '.lambda.DocumentDBEventSourceConfig',
      '10': 'documentdbeventsourceconfig'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'filtercriteria',
      '3': 439219323,
      '4': 1,
      '5': 11,
      '6': '.lambda.FilterCriteria',
      '10': 'filtercriteria'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'functionresponsetypes',
      '3': 382292260,
      '4': 3,
      '5': 14,
      '6': '.lambda.FunctionResponseType',
      '10': 'functionresponsetypes'
    },
    {'1': 'kmskeyarn', '3': 117627377, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {
      '1': 'loggingconfig',
      '3': 424359625,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventSourceMappingLoggingConfig',
      '10': 'loggingconfig'
    },
    {
      '1': 'maximumbatchingwindowinseconds',
      '3': 346663320,
      '4': 1,
      '5': 5,
      '10': 'maximumbatchingwindowinseconds'
    },
    {
      '1': 'maximumrecordageinseconds',
      '3': 102344982,
      '4': 1,
      '5': 5,
      '10': 'maximumrecordageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
    {
      '1': 'metricsconfig',
      '3': 412971857,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventSourceMappingMetricsConfig',
      '10': 'metricsconfig'
    },
    {
      '1': 'parallelizationfactor',
      '3': 400337694,
      '4': 1,
      '5': 5,
      '10': 'parallelizationfactor'
    },
    {
      '1': 'provisionedpollerconfig',
      '3': 275676602,
      '4': 1,
      '5': 11,
      '6': '.lambda.ProvisionedPollerConfig',
      '10': 'provisionedpollerconfig'
    },
    {'1': 'queues', '3': 222519730, '4': 3, '5': 9, '10': 'queues'},
    {
      '1': 'scalingconfig',
      '3': 392871661,
      '4': 1,
      '5': 11,
      '6': '.lambda.ScalingConfig',
      '10': 'scalingconfig'
    },
    {
      '1': 'selfmanagedeventsource',
      '3': 283601786,
      '4': 1,
      '5': 11,
      '6': '.lambda.SelfManagedEventSource',
      '10': 'selfmanagedeventsource'
    },
    {
      '1': 'selfmanagedkafkaeventsourceconfig',
      '3': 322222578,
      '4': 1,
      '5': 11,
      '6': '.lambda.SelfManagedKafkaEventSourceConfig',
      '10': 'selfmanagedkafkaeventsourceconfig'
    },
    {
      '1': 'sourceaccessconfigurations',
      '3': 371593554,
      '4': 3,
      '5': 11,
      '6': '.lambda.SourceAccessConfiguration',
      '10': 'sourceaccessconfigurations'
    },
    {
      '1': 'startingposition',
      '3': 428771919,
      '4': 1,
      '5': 14,
      '6': '.lambda.EventSourcePosition',
      '10': 'startingposition'
    },
    {
      '1': 'startingpositiontimestamp',
      '3': 144323607,
      '4': 1,
      '5': 9,
      '10': 'startingpositiontimestamp'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.lambda.CreateEventSourceMappingRequest.TagsEntry',
      '10': 'tags'
    },
    {'1': 'topics', '3': 219850038, '4': 3, '5': 9, '10': 'topics'},
    {
      '1': 'tumblingwindowinseconds',
      '3': 372493124,
      '4': 1,
      '5': 5,
      '10': 'tumblingwindowinseconds'
    },
  ],
  '3': [CreateEventSourceMappingRequest_TagsEntry$json],
};

@$core.Deprecated('Use createEventSourceMappingRequestDescriptor instead')
const CreateEventSourceMappingRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateEventSourceMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEventSourceMappingRequestDescriptor = $convert.base64Decode(
    'Ch9DcmVhdGVFdmVudFNvdXJjZU1hcHBpbmdSZXF1ZXN0EoABCiNhbWF6b25tYW5hZ2Vka2Fma2'
    'FldmVudHNvdXJjZWNvbmZpZxi8t9YcIAEoCzIrLmxhbWJkYS5BbWF6b25NYW5hZ2VkS2Fma2FF'
    'dmVudFNvdXJjZUNvbmZpZ1IjYW1hem9ubWFuYWdlZGthZmthZXZlbnRzb3VyY2Vjb25maWcSIA'
    'oJYmF0Y2hzaXplGNvJ05cBIAEoBVIJYmF0Y2hzaXplEkIKGmJpc2VjdGJhdGNob25mdW5jdGlv'
    'bmVycm9yGNu81oMBIAEoCFIaYmlzZWN0YmF0Y2hvbmZ1bmN0aW9uZXJyb3ISSgoRZGVzdGluYX'
    'Rpb25jb25maWcY7rCRWCABKAsyGS5sYW1iZGEuRGVzdGluYXRpb25Db25maWdSEWRlc3RpbmF0'
    'aW9uY29uZmlnEmgKG2RvY3VtZW50ZGJldmVudHNvdXJjZWNvbmZpZxiO5MJSIAEoCzIjLmxhbW'
    'JkYS5Eb2N1bWVudERCRXZlbnRTb3VyY2VDb25maWdSG2RvY3VtZW50ZGJldmVudHNvdXJjZWNv'
    'bmZpZxIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIqCg5ldmVudHNvdXJjZWFybhjGyo'
    'qSASABKAlSDmV2ZW50c291cmNlYXJuEkIKDmZpbHRlcmNyaXRlcmlhGPvot9EBIAEoCzIWLmxh'
    'bWJkYS5GaWx0ZXJDcml0ZXJpYVIOZmlsdGVyY3JpdGVyaWESJgoMZnVuY3Rpb25uYW1lGKOIv9'
    '8BIAEoCVIMZnVuY3Rpb25uYW1lElYKFWZ1bmN0aW9ucmVzcG9uc2V0eXBlcxikoqW2ASADKA4y'
    'HC5sYW1iZGEuRnVuY3Rpb25SZXNwb25zZVR5cGVSFWZ1bmN0aW9ucmVzcG9uc2V0eXBlcxIfCg'
    'lrbXNrZXlhcm4Y8bOLOCABKAlSCWttc2tleWFybhJRCg1sb2dnaW5nY29uZmlnGMntrMoBIAEo'
    'CzInLmxhbWJkYS5FdmVudFNvdXJjZU1hcHBpbmdMb2dnaW5nQ29uZmlnUg1sb2dnaW5nY29uZm'
    'lnEkoKHm1heGltdW1iYXRjaGluZ3dpbmRvd2luc2Vjb25kcxiY06alASABKAVSHm1heGltdW1i'
    'YXRjaGluZ3dpbmRvd2luc2Vjb25kcxI/ChltYXhpbXVtcmVjb3JkYWdlaW5zZWNvbmRzGJbS5j'
    'AgASgFUhltYXhpbXVtcmVjb3JkYWdlaW5zZWNvbmRzEjUKFG1heGltdW1yZXRyeWF0dGVtcHRz'
    'GMCouTUgASgFUhRtYXhpbXVtcmV0cnlhdHRlbXB0cxJRCg1tZXRyaWNzY29uZmlnGNHm9cQBIA'
    'EoCzInLmxhbWJkYS5FdmVudFNvdXJjZU1hcHBpbmdNZXRyaWNzQ29uZmlnUg1tZXRyaWNzY29u'
    'ZmlnEjgKFXBhcmFsbGVsaXphdGlvbmZhY3Rvchie1vK+ASABKAVSFXBhcmFsbGVsaXphdGlvbm'
    'ZhY3RvchJdChdwcm92aXNpb25lZHBvbGxlcmNvbmZpZxi6+7mDASABKAsyHy5sYW1iZGEuUHJv'
    'dmlzaW9uZWRQb2xsZXJDb25maWdSF3Byb3Zpc2lvbmVkcG9sbGVyY29uZmlnEhkKBnF1ZXVlcx'
    'iyw41qIAMoCVIGcXVldWVzEj8KDXNjYWxpbmdjb25maWcY7f2quwEgASgLMhUubGFtYmRhLlNj'
    'YWxpbmdDb25maWdSDXNjYWxpbmdjb25maWcSWgoWc2VsZm1hbmFnZWRldmVudHNvdXJjZRj61p'
    '2HASABKAsyHi5sYW1iZGEuU2VsZk1hbmFnZWRFdmVudFNvdXJjZVIWc2VsZm1hbmFnZWRldmVu'
    'dHNvdXJjZRJ7CiFzZWxmbWFuYWdlZGthZmthZXZlbnRzb3VyY2Vjb25maWcY8vPSmQEgASgLMi'
    'kubGFtYmRhLlNlbGZNYW5hZ2VkS2Fma2FFdmVudFNvdXJjZUNvbmZpZ1Ihc2VsZm1hbmFnZWRr'
    'YWZrYWV2ZW50c291cmNlY29uZmlnEmUKGnNvdXJjZWFjY2Vzc2NvbmZpZ3VyYXRpb25zGNKimL'
    'EBIAMoCzIhLmxhbWJkYS5Tb3VyY2VBY2Nlc3NDb25maWd1cmF0aW9uUhpzb3VyY2VhY2Nlc3Nj'
    'b25maWd1cmF0aW9ucxJLChBzdGFydGluZ3Bvc2l0aW9uGM+UuswBIAEoDjIbLmxhbWJkYS5Fdm'
    'VudFNvdXJjZVBvc2l0aW9uUhBzdGFydGluZ3Bvc2l0aW9uEj8KGXN0YXJ0aW5ncG9zaXRpb250'
    'aW1lc3RhbXAYl+joRCABKAlSGXN0YXJ0aW5ncG9zaXRpb250aW1lc3RhbXASSQoEdGFncxjBwf'
    'a1ASADKAsyMS5sYW1iZGEuQ3JlYXRlRXZlbnRTb3VyY2VNYXBwaW5nUmVxdWVzdC5UYWdzRW50'
    'cnlSBHRhZ3MSGQoGdG9waWNzGLbK6mggAygJUgZ0b3BpY3MSPAoXdHVtYmxpbmd3aW5kb3dpbn'
    'NlY29uZHMYxJbPsQEgASgFUhd0dW1ibGluZ3dpbmRvd2luc2Vjb25kcxo3CglUYWdzRW50cnkS'
    'EAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createFunctionRequestDescriptor instead')
const CreateFunctionRequest$json = {
  '1': 'CreateFunctionRequest',
  '2': [
    {
      '1': 'architectures',
      '3': 530490948,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'architectures'
    },
    {
      '1': 'capacityproviderconfig',
      '3': 52030623,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderConfig',
      '10': 'capacityproviderconfig'
    },
    {
      '1': 'code',
      '3': 425572629,
      '4': 1,
      '5': 11,
      '6': '.lambda.FunctionCode',
      '10': 'code'
    },
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
    {
      '1': 'deadletterconfig',
      '3': 79786642,
      '4': 1,
      '5': 11,
      '6': '.lambda.DeadLetterConfig',
      '10': 'deadletterconfig'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'durableconfig',
      '3': 206326279,
      '4': 1,
      '5': 11,
      '6': '.lambda.DurableConfig',
      '10': 'durableconfig'
    },
    {
      '1': 'environment',
      '3': 119823003,
      '4': 1,
      '5': 11,
      '6': '.lambda.Environment',
      '10': 'environment'
    },
    {
      '1': 'ephemeralstorage',
      '3': 365965382,
      '4': 1,
      '5': 11,
      '6': '.lambda.EphemeralStorage',
      '10': 'ephemeralstorage'
    },
    {
      '1': 'filesystemconfigs',
      '3': 490453750,
      '4': 3,
      '5': 11,
      '6': '.lambda.FileSystemConfig',
      '10': 'filesystemconfigs'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'handler', '3': 81160724, '4': 1, '5': 9, '10': 'handler'},
    {
      '1': 'imageconfig',
      '3': 281970485,
      '4': 1,
      '5': 11,
      '6': '.lambda.ImageConfig',
      '10': 'imageconfig'
    },
    {'1': 'kmskeyarn', '3': 117627377, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {'1': 'layers', '3': 478144896, '4': 3, '5': 9, '10': 'layers'},
    {
      '1': 'loggingconfig',
      '3': 424359625,
      '4': 1,
      '5': 11,
      '6': '.lambda.LoggingConfig',
      '10': 'loggingconfig'
    },
    {'1': 'memorysize', '3': 55523120, '4': 1, '5': 5, '10': 'memorysize'},
    {
      '1': 'packagetype',
      '3': 517524132,
      '4': 1,
      '5': 14,
      '6': '.lambda.PackageType',
      '10': 'packagetype'
    },
    {'1': 'publish', '3': 207759785, '4': 1, '5': 8, '10': 'publish'},
    {
      '1': 'publishto',
      '3': 524127682,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionVersionLatestPublished',
      '10': 'publishto'
    },
    {'1': 'role', '3': 271285818, '4': 1, '5': 9, '10': 'role'},
    {
      '1': 'runtime',
      '3': 359311308,
      '4': 1,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'runtime'
    },
    {
      '1': 'snapstart',
      '3': 283273032,
      '4': 1,
      '5': 11,
      '6': '.lambda.SnapStart',
      '10': 'snapstart'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.lambda.CreateFunctionRequest.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'tenancyconfig',
      '3': 215700986,
      '4': 1,
      '5': 11,
      '6': '.lambda.TenancyConfig',
      '10': 'tenancyconfig'
    },
    {'1': 'timeout', '3': 47808041, '4': 1, '5': 5, '10': 'timeout'},
    {
      '1': 'tracingconfig',
      '3': 19554860,
      '4': 1,
      '5': 11,
      '6': '.lambda.TracingConfig',
      '10': 'tracingconfig'
    },
    {
      '1': 'vpcconfig',
      '3': 194980743,
      '4': 1,
      '5': 11,
      '6': '.lambda.VpcConfig',
      '10': 'vpcconfig'
    },
  ],
  '3': [CreateFunctionRequest_TagsEntry$json],
};

@$core.Deprecated('Use createFunctionRequestDescriptor instead')
const CreateFunctionRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFunctionRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVGdW5jdGlvblJlcXVlc3QSPgoNYXJjaGl0ZWN0dXJlcxjEzPr8ASADKA4yFC5sYW'
    '1iZGEuQXJjaGl0ZWN0dXJlUg1hcmNoaXRlY3R1cmVzElkKFmNhcGFjaXR5cHJvdmlkZXJjb25m'
    'aWcYn9nnGCABKAsyHi5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlckNvbmZpZ1IWY2FwYWNpdHlwcm'
    '92aWRlcmNvbmZpZxIsCgRjb2RlGJXy9soBIAEoCzIULmxhbWJkYS5GdW5jdGlvbkNvZGVSBGNv'
    'ZGUSNgoUY29kZXNpZ25pbmdjb25maWdhcm4Ywfz38AEgASgJUhRjb2Rlc2lnbmluZ2NvbmZpZ2'
    'FybhJHChBkZWFkbGV0dGVyY29uZmlnGJLlhSYgASgLMhgubGFtYmRhLkRlYWRMZXR0ZXJDb25m'
    'aWdSEGRlYWRsZXR0ZXJjb25maWcSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW'
    '9uEj4KDWR1cmFibGVjb25maWcYh5SxYiABKAsyFS5sYW1iZGEuRHVyYWJsZUNvbmZpZ1INZHVy'
    'YWJsZWNvbmZpZxI4CgtlbnZpcm9ubWVudBibtZE5IAEoCzITLmxhbWJkYS5FbnZpcm9ubWVudF'
    'ILZW52aXJvbm1lbnQSSAoQZXBoZW1lcmFsc3RvcmFnZRjG4MCuASABKAsyGC5sYW1iZGEuRXBo'
    'ZW1lcmFsU3RvcmFnZVIQZXBoZW1lcmFsc3RvcmFnZRJKChFmaWxlc3lzdGVtY29uZmlncxj29e'
    '7pASADKAsyGC5sYW1iZGEuRmlsZVN5c3RlbUNvbmZpZ1IRZmlsZXN5c3RlbWNvbmZpZ3MSJgoM'
    'ZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb25uYW1lEhsKB2hhbmRsZXIYlNTZJiABKA'
    'lSB2hhbmRsZXISOQoLaW1hZ2Vjb25maWcYtY66hgEgASgLMhMubGFtYmRhLkltYWdlQ29uZmln'
    'UgtpbWFnZWNvbmZpZxIfCglrbXNrZXlhcm4Y8bOLOCABKAlSCWttc2tleWFybhIaCgZsYXllcn'
    'MYgNP/4wEgAygJUgZsYXllcnMSPwoNbG9nZ2luZ2NvbmZpZxjJ7azKASABKAsyFS5sYW1iZGEu'
    'TG9nZ2luZ0NvbmZpZ1INbG9nZ2luZ2NvbmZpZxIhCgptZW1vcnlzaXplGLDuvBogASgFUgptZW'
    '1vcnlzaXplEjkKC3BhY2thZ2V0eXBlGKSV4/YBIAEoDjITLmxhbWJkYS5QYWNrYWdlVHlwZVIL'
    'cGFja2FnZXR5cGUSGwoHcHVibGlzaBip04hjIAEoCFIHcHVibGlzaBJICglwdWJsaXNodG8Ywp'
    'v2+QEgASgOMiYubGFtYmRhLkZ1bmN0aW9uVmVyc2lvbkxhdGVzdFB1Ymxpc2hlZFIJcHVibGlz'
    'aHRvEhYKBHJvbGUYuvytgQEgASgJUgRyb2xlEi0KB3J1bnRpbWUYzM+qqwEgASgOMg8ubGFtYm'
    'RhLlJ1bnRpbWVSB3J1bnRpbWUSMwoJc25hcHN0YXJ0GMjOiYcBIAEoCzIRLmxhbWJkYS5TbmFw'
    'U3RhcnRSCXNuYXBzdGFydBI/CgR0YWdzGMHB9rUBIAMoCzInLmxhbWJkYS5DcmVhdGVGdW5jdG'
    'lvblJlcXVlc3QuVGFnc0VudHJ5UgR0YWdzEj4KDXRlbmFuY3ljb25maWcY+qvtZiABKAsyFS5s'
    'YW1iZGEuVGVuYW5jeUNvbmZpZ1INdGVuYW5jeWNvbmZpZxIbCgd0aW1lb3V0GKn85RYgASgFUg'
    'd0aW1lb3V0Ej4KDXRyYWNpbmdjb25maWcYrMSpCSABKAsyFS5sYW1iZGEuVHJhY2luZ0NvbmZp'
    'Z1INdHJhY2luZ2NvbmZpZxIyCgl2cGNjb25maWcYh9f8XCABKAsyES5sYW1iZGEuVnBjQ29uZm'
    'lnUgl2cGNjb25maWcaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIg'
    'ASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createFunctionUrlConfigRequestDescriptor instead')
const CreateFunctionUrlConfigRequest$json = {
  '1': 'CreateFunctionUrlConfigRequest',
  '2': [
    {
      '1': 'authtype',
      '3': 477704248,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionUrlAuthType',
      '10': 'authtype'
    },
    {
      '1': 'cors',
      '3': 260753653,
      '4': 1,
      '5': 11,
      '6': '.lambda.Cors',
      '10': 'cors'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'invokemode',
      '3': 414956667,
      '4': 1,
      '5': 14,
      '6': '.lambda.InvokeMode',
      '10': 'invokemode'
    },
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `CreateFunctionUrlConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFunctionUrlConfigRequestDescriptor = $convert.base64Decode(
    'Ch5DcmVhdGVGdW5jdGlvblVybENvbmZpZ1JlcXVlc3QSOwoIYXV0aHR5cGUYuODk4wEgASgOMh'
    'subGFtYmRhLkZ1bmN0aW9uVXJsQXV0aFR5cGVSCGF1dGh0eXBlEiMKBGNvcnMY9ZGrfCABKAsy'
    'DC5sYW1iZGEuQ29yc1IEY29ycxImCgxmdW5jdGlvbm5hbWUYo4i/3wEgASgJUgxmdW5jdGlvbm'
    '5hbWUSNgoKaW52b2tlbW9kZRj7+O7FASABKA4yEi5sYW1iZGEuSW52b2tlTW9kZVIKaW52b2tl'
    'bW9kZRIgCglxdWFsaWZpZXIY4LWR+wEgASgJUglxdWFsaWZpZXI=');

@$core.Deprecated('Use createFunctionUrlConfigResponseDescriptor instead')
const CreateFunctionUrlConfigResponse$json = {
  '1': 'CreateFunctionUrlConfigResponse',
  '2': [
    {
      '1': 'authtype',
      '3': 477704248,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionUrlAuthType',
      '10': 'authtype'
    },
    {
      '1': 'cors',
      '3': 260753653,
      '4': 1,
      '5': 11,
      '6': '.lambda.Cors',
      '10': 'cors'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'functionurl', '3': 449381947, '4': 1, '5': 9, '10': 'functionurl'},
    {
      '1': 'invokemode',
      '3': 414956667,
      '4': 1,
      '5': 14,
      '6': '.lambda.InvokeMode',
      '10': 'invokemode'
    },
  ],
};

/// Descriptor for `CreateFunctionUrlConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFunctionUrlConfigResponseDescriptor = $convert.base64Decode(
    'Ch9DcmVhdGVGdW5jdGlvblVybENvbmZpZ1Jlc3BvbnNlEjsKCGF1dGh0eXBlGLjg5OMBIAEoDj'
    'IbLmxhbWJkYS5GdW5jdGlvblVybEF1dGhUeXBlUghhdXRodHlwZRIjCgRjb3JzGPWRq3wgASgL'
    'MgwubGFtYmRhLkNvcnNSBGNvcnMSJQoMY3JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlvbn'
    'RpbWUSJAoLZnVuY3Rpb25hcm4Y8cmOtgEgASgJUgtmdW5jdGlvbmFybhIkCgtmdW5jdGlvbnVy'
    'bBi7jKTWASABKAlSC2Z1bmN0aW9udXJsEjYKCmludm9rZW1vZGUY+/juxQEgASgOMhIubGFtYm'
    'RhLkludm9rZU1vZGVSCmludm9rZW1vZGU=');

@$core.Deprecated('Use deadLetterConfigDescriptor instead')
const DeadLetterConfig$json = {
  '1': 'DeadLetterConfig',
  '2': [
    {'1': 'targetarn', '3': 217664144, '4': 1, '5': 9, '10': 'targetarn'},
  ],
};

/// Descriptor for `DeadLetterConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deadLetterConfigDescriptor = $convert.base64Decode(
    'ChBEZWFkTGV0dGVyQ29uZmlnEh8KCXRhcmdldGFybhiQleVnIAEoCVIJdGFyZ2V0YXJu');

@$core.Deprecated('Use deleteAliasRequestDescriptor instead')
const DeleteAliasRequest$json = {
  '1': 'DeleteAliasRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAliasRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVBbGlhc1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb2'
    '5uYW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deleteCapacityProviderRequestDescriptor instead')
const DeleteCapacityProviderRequest$json = {
  '1': 'DeleteCapacityProviderRequest',
  '2': [
    {
      '1': 'capacityprovidername',
      '3': 466230132,
      '4': 1,
      '5': 9,
      '10': 'capacityprovidername'
    },
  ],
};

/// Descriptor for `DeleteCapacityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCapacityProviderRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVDYXBhY2l0eVByb3ZpZGVyUmVxdWVzdBI2ChRjYXBhY2l0eXByb3ZpZGVybmFtZR'
        'j0tqjeASABKAlSFGNhcGFjaXR5cHJvdmlkZXJuYW1l');

@$core.Deprecated('Use deleteCapacityProviderResponseDescriptor instead')
const DeleteCapacityProviderResponse$json = {
  '1': 'DeleteCapacityProviderResponse',
  '2': [
    {
      '1': 'capacityprovider',
      '3': 448074961,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProvider',
      '10': 'capacityprovider'
    },
  ],
};

/// Descriptor for `DeleteCapacityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCapacityProviderResponseDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVDYXBhY2l0eVByb3ZpZGVyUmVzcG9uc2USSAoQY2FwYWNpdHlwcm92aWRlchjRqd'
        'TVASABKAsyGC5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlclIQY2FwYWNpdHlwcm92aWRlcg==');

@$core.Deprecated('Use deleteCodeSigningConfigRequestDescriptor instead')
const DeleteCodeSigningConfigRequest$json = {
  '1': 'DeleteCodeSigningConfigRequest',
  '2': [
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
  ],
};

/// Descriptor for `DeleteCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCodeSigningConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVDb2RlU2lnbmluZ0NvbmZpZ1JlcXVlc3QSNgoUY29kZXNpZ25pbmdjb25maWdhcm'
        '4Ywfz38AEgASgJUhRjb2Rlc2lnbmluZ2NvbmZpZ2Fybg==');

@$core.Deprecated('Use deleteCodeSigningConfigResponseDescriptor instead')
const DeleteCodeSigningConfigResponse$json = {
  '1': 'DeleteCodeSigningConfigResponse',
};

/// Descriptor for `DeleteCodeSigningConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCodeSigningConfigResponseDescriptor =
    $convert.base64Decode('Ch9EZWxldGVDb2RlU2lnbmluZ0NvbmZpZ1Jlc3BvbnNl');

@$core.Deprecated('Use deleteEventSourceMappingRequestDescriptor instead')
const DeleteEventSourceMappingRequest$json = {
  '1': 'DeleteEventSourceMappingRequest',
  '2': [
    {'1': 'uuid', '3': 91981875, '4': 1, '5': 9, '10': 'uuid'},
  ],
};

/// Descriptor for `DeleteEventSourceMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEventSourceMappingRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZWxldGVFdmVudFNvdXJjZU1hcHBpbmdSZXF1ZXN0EhUKBHV1aWQYs5DuKyABKAlSBHV1aW'
        'Q=');

@$core
    .Deprecated('Use deleteFunctionCodeSigningConfigRequestDescriptor instead')
const DeleteFunctionCodeSigningConfigRequest$json = {
  '1': 'DeleteFunctionCodeSigningConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `DeleteFunctionCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFunctionCodeSigningConfigRequestDescriptor =
    $convert.base64Decode(
        'CiZEZWxldGVGdW5jdGlvbkNvZGVTaWduaW5nQ29uZmlnUmVxdWVzdBImCgxmdW5jdGlvbm5hbW'
        'UYo4i/3wEgASgJUgxmdW5jdGlvbm5hbWU=');

@$core.Deprecated('Use deleteFunctionConcurrencyRequestDescriptor instead')
const DeleteFunctionConcurrencyRequest$json = {
  '1': 'DeleteFunctionConcurrencyRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `DeleteFunctionConcurrencyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFunctionConcurrencyRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVGdW5jdGlvbkNvbmN1cnJlbmN5UmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4i/3w'
        'EgASgJUgxmdW5jdGlvbm5hbWU=');

@$core
    .Deprecated('Use deleteFunctionEventInvokeConfigRequestDescriptor instead')
const DeleteFunctionEventInvokeConfigRequest$json = {
  '1': 'DeleteFunctionEventInvokeConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `DeleteFunctionEventInvokeConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFunctionEventInvokeConfigRequestDescriptor =
    $convert.base64Decode(
        'CiZEZWxldGVGdW5jdGlvbkV2ZW50SW52b2tlQ29uZmlnUmVxdWVzdBImCgxmdW5jdGlvbm5hbW'
        'UYo4i/3wEgASgJUgxmdW5jdGlvbm5hbWUSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlm'
        'aWVy');

@$core.Deprecated('Use deleteFunctionRequestDescriptor instead')
const DeleteFunctionRequest$json = {
  '1': 'DeleteFunctionRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `DeleteFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFunctionRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVGdW5jdGlvblJlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3'
    'Rpb25uYW1lEiAKCXF1YWxpZmllchjgtZH7ASABKAlSCXF1YWxpZmllcg==');

@$core.Deprecated('Use deleteFunctionResponseDescriptor instead')
const DeleteFunctionResponse$json = {
  '1': 'DeleteFunctionResponse',
  '2': [
    {'1': 'statuscode', '3': 303830783, '4': 1, '5': 5, '10': 'statuscode'},
  ],
};

/// Descriptor for `DeleteFunctionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFunctionResponseDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVGdW5jdGlvblJlc3BvbnNlEiIKCnN0YXR1c2NvZGUY/63wkAEgASgFUgpzdGF0dX'
        'Njb2Rl');

@$core.Deprecated('Use deleteFunctionUrlConfigRequestDescriptor instead')
const DeleteFunctionUrlConfigRequest$json = {
  '1': 'DeleteFunctionUrlConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `DeleteFunctionUrlConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteFunctionUrlConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVGdW5jdGlvblVybENvbmZpZ1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIA'
        'EoCVIMZnVuY3Rpb25uYW1lEiAKCXF1YWxpZmllchjgtZH7ASABKAlSCXF1YWxpZmllcg==');

@$core.Deprecated('Use deleteLayerVersionRequestDescriptor instead')
const DeleteLayerVersionRequest$json = {
  '1': 'DeleteLayerVersionRequest',
  '2': [
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
    {
      '1': 'versionnumber',
      '3': 209346775,
      '4': 1,
      '5': 3,
      '10': 'versionnumber'
    },
  ],
};

/// Descriptor for `DeleteLayerVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteLayerVersionRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVMYXllclZlcnNpb25SZXF1ZXN0EiAKCWxheWVybmFtZRiWqpjtASABKAlSCWxheW'
        'VybmFtZRInCg12ZXJzaW9ubnVtYmVyGNfB6WMgASgDUg12ZXJzaW9ubnVtYmVy');

@$core.Deprecated(
    'Use deleteProvisionedConcurrencyConfigRequestDescriptor instead')
const DeleteProvisionedConcurrencyConfigRequest$json = {
  '1': 'DeleteProvisionedConcurrencyConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `DeleteProvisionedConcurrencyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteProvisionedConcurrencyConfigRequestDescriptor = $convert.base64Decode(
        'CilEZWxldGVQcm92aXNpb25lZENvbmN1cnJlbmN5Q29uZmlnUmVxdWVzdBImCgxmdW5jdGlvbm'
        '5hbWUYo4i/3wEgASgJUgxmdW5jdGlvbm5hbWUSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVh'
        'bGlmaWVy');

@$core.Deprecated('Use destinationConfigDescriptor instead')
const DestinationConfig$json = {
  '1': 'DestinationConfig',
  '2': [
    {
      '1': 'onfailure',
      '3': 424696739,
      '4': 1,
      '5': 11,
      '6': '.lambda.OnFailure',
      '10': 'onfailure'
    },
    {
      '1': 'onsuccess',
      '3': 332738196,
      '4': 1,
      '5': 11,
      '6': '.lambda.OnSuccess',
      '10': 'onsuccess'
    },
  ],
};

/// Descriptor for `DestinationConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List destinationConfigDescriptor = $convert.base64Decode(
    'ChFEZXN0aW5hdGlvbkNvbmZpZxIzCglvbmZhaWx1cmUYo7fBygEgASgLMhEubGFtYmRhLk9uRm'
    'FpbHVyZVIJb25mYWlsdXJlEjMKCW9uc3VjY2VzcxiU3dSeASABKAsyES5sYW1iZGEuT25TdWNj'
    'ZXNzUglvbnN1Y2Nlc3M=');

@$core.Deprecated('Use documentDBEventSourceConfigDescriptor instead')
const DocumentDBEventSourceConfig$json = {
  '1': 'DocumentDBEventSourceConfig',
  '2': [
    {
      '1': 'collectionname',
      '3': 350862923,
      '4': 1,
      '5': 9,
      '10': 'collectionname'
    },
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {
      '1': 'fulldocument',
      '3': 450918030,
      '4': 1,
      '5': 14,
      '6': '.lambda.FullDocument',
      '10': 'fulldocument'
    },
  ],
};

/// Descriptor for `DocumentDBEventSourceConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentDBEventSourceConfigDescriptor = $convert.base64Decode(
    'ChtEb2N1bWVudERCRXZlbnRTb3VyY2VDb25maWcSKgoOY29sbGVjdGlvbm5hbWUYy/ympwEgAS'
    'gJUg5jb2xsZWN0aW9ubmFtZRIlCgxkYXRhYmFzZW5hbWUY3LLZKiABKAlSDGRhdGFiYXNlbmFt'
    'ZRI8CgxmdWxsZG9jdW1lbnQYju2B1wEgASgOMhQubGFtYmRhLkZ1bGxEb2N1bWVudFIMZnVsbG'
    'RvY3VtZW50');

@$core.Deprecated('Use durableConfigDescriptor instead')
const DurableConfig$json = {
  '1': 'DurableConfig',
  '2': [
    {
      '1': 'executiontimeout',
      '3': 502184601,
      '4': 1,
      '5': 5,
      '10': 'executiontimeout'
    },
    {
      '1': 'retentionperiodindays',
      '3': 324977873,
      '4': 1,
      '5': 5,
      '10': 'retentionperiodindays'
    },
  ],
};

/// Descriptor for `DurableConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List durableConfigDescriptor = $convert.base64Decode(
    'Cg1EdXJhYmxlQ29uZmlnEi4KEGV4ZWN1dGlvbnRpbWVvdXQYmfW67wEgASgFUhBleGVjdXRpb2'
    '50aW1lb3V0EjgKFXJldGVudGlvbnBlcmlvZGluZGF5cxjRifuaASABKAVSFXJldGVudGlvbnBl'
    'cmlvZGluZGF5cw==');

@$core
    .Deprecated('Use durableExecutionAlreadyStartedExceptionDescriptor instead')
const DurableExecutionAlreadyStartedException$json = {
  '1': 'DurableExecutionAlreadyStartedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `DurableExecutionAlreadyStartedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List durableExecutionAlreadyStartedExceptionDescriptor =
    $convert.base64Decode(
        'CidEdXJhYmxlRXhlY3V0aW9uQWxyZWFkeVN0YXJ0ZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZRIWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use eC2AccessDeniedExceptionDescriptor instead')
const EC2AccessDeniedException$json = {
  '1': 'EC2AccessDeniedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `EC2AccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eC2AccessDeniedExceptionDescriptor =
    $convert.base64Decode(
        'ChhFQzJBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use eC2ThrottledExceptionDescriptor instead')
const EC2ThrottledException$json = {
  '1': 'EC2ThrottledException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `EC2ThrottledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eC2ThrottledExceptionDescriptor = $convert.base64Decode(
    'ChVFQzJUaHJvdHRsZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZRIWCg'
    'R0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use eC2UnexpectedExceptionDescriptor instead')
const EC2UnexpectedException$json = {
  '1': 'EC2UnexpectedException',
  '2': [
    {'1': 'ec2errorcode', '3': 35455333, '4': 1, '5': 9, '10': 'ec2errorcode'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `EC2UnexpectedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eC2UnexpectedExceptionDescriptor = $convert.base64Decode(
    'ChZFQzJVbmV4cGVjdGVkRXhjZXB0aW9uEiUKDGVjMmVycm9yY29kZRjlgvQQIAEoCVIMZWMyZX'
    'Jyb3Jjb2RlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2USFgoEdHlwZRjuoNeKASABKAlS'
    'BHR5cGU=');

@$core.Deprecated('Use eFSIOExceptionDescriptor instead')
const EFSIOException$json = {
  '1': 'EFSIOException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `EFSIOException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eFSIOExceptionDescriptor = $convert.base64Decode(
    'Cg5FRlNJT0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdlEhYKBHR5cGUY7q'
    'DXigEgASgJUgR0eXBl');

@$core.Deprecated('Use eFSMountConnectivityExceptionDescriptor instead')
const EFSMountConnectivityException$json = {
  '1': 'EFSMountConnectivityException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `EFSMountConnectivityException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eFSMountConnectivityExceptionDescriptor =
    $convert.base64Decode(
        'Ch1FRlNNb3VudENvbm5lY3Rpdml0eUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdlEhYKBHR5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use eFSMountFailureExceptionDescriptor instead')
const EFSMountFailureException$json = {
  '1': 'EFSMountFailureException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `EFSMountFailureException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eFSMountFailureExceptionDescriptor =
    $convert.base64Decode(
        'ChhFRlNNb3VudEZhaWx1cmVFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use eFSMountTimeoutExceptionDescriptor instead')
const EFSMountTimeoutException$json = {
  '1': 'EFSMountTimeoutException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `EFSMountTimeoutException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eFSMountTimeoutExceptionDescriptor =
    $convert.base64Decode(
        'ChhFRlNNb3VudFRpbWVvdXRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use eNILimitReachedExceptionDescriptor instead')
const ENILimitReachedException$json = {
  '1': 'ENILimitReachedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `ENILimitReachedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eNILimitReachedExceptionDescriptor =
    $convert.base64Decode(
        'ChhFTklMaW1pdFJlYWNoZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use environmentDescriptor instead')
const Environment$json = {
  '1': 'Environment',
  '2': [
    {
      '1': 'variables',
      '3': 429322339,
      '4': 3,
      '5': 11,
      '6': '.lambda.Environment.VariablesEntry',
      '10': 'variables'
    },
  ],
  '3': [Environment_VariablesEntry$json],
};

@$core.Deprecated('Use environmentDescriptor instead')
const Environment_VariablesEntry$json = {
  '1': 'VariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Environment`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List environmentDescriptor = $convert.base64Decode(
    'CgtFbnZpcm9ubWVudBJECgl2YXJpYWJsZXMY4+DbzAEgAygLMiIubGFtYmRhLkVudmlyb25tZW'
    '50LlZhcmlhYmxlc0VudHJ5Ugl2YXJpYWJsZXMaPAoOVmFyaWFibGVzRW50cnkSEAoDa2V5GAEg'
    'ASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use environmentErrorDescriptor instead')
const EnvironmentError$json = {
  '1': 'EnvironmentError',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EnvironmentError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List environmentErrorDescriptor = $convert.base64Decode(
    'ChBFbnZpcm9ubWVudEVycm9yEh8KCWVycm9yY29kZRiZ1sMQIAEoCVIJZXJyb3Jjb2RlEhsKB2'
    '1lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use environmentResponseDescriptor instead')
const EnvironmentResponse$json = {
  '1': 'EnvironmentResponse',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EnvironmentError',
      '10': 'error'
    },
    {
      '1': 'variables',
      '3': 429322339,
      '4': 3,
      '5': 11,
      '6': '.lambda.EnvironmentResponse.VariablesEntry',
      '10': 'variables'
    },
  ],
  '3': [EnvironmentResponse_VariablesEntry$json],
};

@$core.Deprecated('Use environmentResponseDescriptor instead')
const EnvironmentResponse_VariablesEntry$json = {
  '1': 'VariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `EnvironmentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List environmentResponseDescriptor = $convert.base64Decode(
    'ChNFbnZpcm9ubWVudFJlc3BvbnNlEjIKBWVycm9yGPK5tpwBIAEoCzIYLmxhbWJkYS5FbnZpcm'
    '9ubWVudEVycm9yUgVlcnJvchJMCgl2YXJpYWJsZXMY4+DbzAEgAygLMioubGFtYmRhLkVudmly'
    'b25tZW50UmVzcG9uc2UuVmFyaWFibGVzRW50cnlSCXZhcmlhYmxlcxo8Cg5WYXJpYWJsZXNFbn'
    'RyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use ephemeralStorageDescriptor instead')
const EphemeralStorage$json = {
  '1': 'EphemeralStorage',
  '2': [
    {'1': 'size', '3': 105352829, '4': 1, '5': 5, '10': 'size'},
  ],
};

/// Descriptor for `EphemeralStorage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ephemeralStorageDescriptor = $convert
    .base64Decode('ChBFcGhlbWVyYWxTdG9yYWdlEhUKBHNpemUY/ZyeMiABKAVSBHNpemU=');

@$core.Deprecated('Use errorObjectDescriptor instead')
const ErrorObject$json = {
  '1': 'ErrorObject',
  '2': [
    {'1': 'errordata', '3': 164186722, '4': 1, '5': 9, '10': 'errordata'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'errortype', '3': 398848954, '4': 1, '5': 9, '10': 'errortype'},
    {'1': 'stacktrace', '3': 306550173, '4': 3, '5': 9, '10': 'stacktrace'},
  ],
};

/// Descriptor for `ErrorObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorObjectDescriptor = $convert.base64Decode(
    'CgtFcnJvck9iamVjdBIfCgllcnJvcmRhdGEY4pSlTiABKAlSCWVycm9yZGF0YRImCgxlcnJvcm'
    '1lc3NhZ2UYqYqr9wEgASgJUgxlcnJvcm1lc3NhZ2USIAoJZXJyb3J0eXBlGLrnl74BIAEoCVIJ'
    'ZXJyb3J0eXBlEiIKCnN0YWNrdHJhY2UYnauWkgEgAygJUgpzdGFja3RyYWNl');

@$core.Deprecated('Use eventDescriptor instead')
const Event$json = {
  '1': 'Event',
  '2': [
    {
      '1': 'callbackfaileddetails',
      '3': 428025264,
      '4': 1,
      '5': 11,
      '6': '.lambda.CallbackFailedDetails',
      '10': 'callbackfaileddetails'
    },
    {
      '1': 'callbackstarteddetails',
      '3': 99932848,
      '4': 1,
      '5': 11,
      '6': '.lambda.CallbackStartedDetails',
      '10': 'callbackstarteddetails'
    },
    {
      '1': 'callbacksucceededdetails',
      '3': 335863628,
      '4': 1,
      '5': 11,
      '6': '.lambda.CallbackSucceededDetails',
      '10': 'callbacksucceededdetails'
    },
    {
      '1': 'callbacktimedoutdetails',
      '3': 57602236,
      '4': 1,
      '5': 11,
      '6': '.lambda.CallbackTimedOutDetails',
      '10': 'callbacktimedoutdetails'
    },
    {
      '1': 'chainedinvokefaileddetails',
      '3': 470187921,
      '4': 1,
      '5': 11,
      '6': '.lambda.ChainedInvokeFailedDetails',
      '10': 'chainedinvokefaileddetails'
    },
    {
      '1': 'chainedinvokestarteddetails',
      '3': 211938659,
      '4': 1,
      '5': 11,
      '6': '.lambda.ChainedInvokeStartedDetails',
      '10': 'chainedinvokestarteddetails'
    },
    {
      '1': 'chainedinvokestoppeddetails',
      '3': 38526227,
      '4': 1,
      '5': 11,
      '6': '.lambda.ChainedInvokeStoppedDetails',
      '10': 'chainedinvokestoppeddetails'
    },
    {
      '1': 'chainedinvokesucceededdetails',
      '3': 519370675,
      '4': 1,
      '5': 11,
      '6': '.lambda.ChainedInvokeSucceededDetails',
      '10': 'chainedinvokesucceededdetails'
    },
    {
      '1': 'chainedinvoketimedoutdetails',
      '3': 245892105,
      '4': 1,
      '5': 11,
      '6': '.lambda.ChainedInvokeTimedOutDetails',
      '10': 'chainedinvoketimedoutdetails'
    },
    {
      '1': 'contextfaileddetails',
      '3': 449331754,
      '4': 1,
      '5': 11,
      '6': '.lambda.ContextFailedDetails',
      '10': 'contextfaileddetails'
    },
    {
      '1': 'contextstarteddetails',
      '3': 340244010,
      '4': 1,
      '5': 11,
      '6': '.lambda.ContextStartedDetails',
      '10': 'contextstarteddetails'
    },
    {
      '1': 'contextsucceededdetails',
      '3': 278760906,
      '4': 1,
      '5': 11,
      '6': '.lambda.ContextSucceededDetails',
      '10': 'contextsucceededdetails'
    },
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 5, '10': 'eventid'},
    {
      '1': 'eventtimestamp',
      '3': 184687758,
      '4': 1,
      '5': 9,
      '10': 'eventtimestamp'
    },
    {
      '1': 'eventtype',
      '3': 468897896,
      '4': 1,
      '5': 14,
      '6': '.lambda.EventType',
      '10': 'eventtype'
    },
    {
      '1': 'executionfaileddetails',
      '3': 177664947,
      '4': 1,
      '5': 11,
      '6': '.lambda.ExecutionFailedDetails',
      '10': 'executionfaileddetails'
    },
    {
      '1': 'executionstarteddetails',
      '3': 465319829,
      '4': 1,
      '5': 11,
      '6': '.lambda.ExecutionStartedDetails',
      '10': 'executionstarteddetails'
    },
    {
      '1': 'executionstoppeddetails',
      '3': 210839205,
      '4': 1,
      '5': 11,
      '6': '.lambda.ExecutionStoppedDetails',
      '10': 'executionstoppeddetails'
    },
    {
      '1': 'executionsucceededdetails',
      '3': 497432545,
      '4': 1,
      '5': 11,
      '6': '.lambda.ExecutionSucceededDetails',
      '10': 'executionsucceededdetails'
    },
    {
      '1': 'executiontimedoutdetails',
      '3': 364692067,
      '4': 1,
      '5': 11,
      '6': '.lambda.ExecutionTimedOutDetails',
      '10': 'executiontimedoutdetails'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'invocationcompleteddetails',
      '3': 202146139,
      '4': 1,
      '5': 11,
      '6': '.lambda.InvocationCompletedDetails',
      '10': 'invocationcompleteddetails'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'parentid', '3': 182711433, '4': 1, '5': 9, '10': 'parentid'},
    {
      '1': 'stepfaileddetails',
      '3': 93829193,
      '4': 1,
      '5': 11,
      '6': '.lambda.StepFailedDetails',
      '10': 'stepfaileddetails'
    },
    {
      '1': 'stepstarteddetails',
      '3': 509042731,
      '4': 1,
      '5': 11,
      '6': '.lambda.StepStartedDetails',
      '10': 'stepstarteddetails'
    },
    {
      '1': 'stepsucceededdetails',
      '3': 171047003,
      '4': 1,
      '5': 11,
      '6': '.lambda.StepSucceededDetails',
      '10': 'stepsucceededdetails'
    },
    {'1': 'subtype', '3': 152988350, '4': 1, '5': 9, '10': 'subtype'},
    {
      '1': 'waitcancelleddetails',
      '3': 203669828,
      '4': 1,
      '5': 11,
      '6': '.lambda.WaitCancelledDetails',
      '10': 'waitcancelleddetails'
    },
    {
      '1': 'waitstarteddetails',
      '3': 272898878,
      '4': 1,
      '5': 11,
      '6': '.lambda.WaitStartedDetails',
      '10': 'waitstarteddetails'
    },
    {
      '1': 'waitsucceededdetails',
      '3': 273317614,
      '4': 1,
      '5': 11,
      '6': '.lambda.WaitSucceededDetails',
      '10': 'waitsucceededdetails'
    },
  ],
};

/// Descriptor for `Event`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDescriptor = $convert.base64Decode(
    'CgVFdmVudBJXChVjYWxsYmFja2ZhaWxlZGRldGFpbHMYsMuMzAEgASgLMh0ubGFtYmRhLkNhbG'
    'xiYWNrRmFpbGVkRGV0YWlsc1IVY2FsbGJhY2tmYWlsZWRkZXRhaWxzElkKFmNhbGxiYWNrc3Rh'
    'cnRlZGRldGFpbHMYsLXTLyABKAsyHi5sYW1iZGEuQ2FsbGJhY2tTdGFydGVkRGV0YWlsc1IWY2'
    'FsbGJhY2tzdGFydGVkZGV0YWlscxJgChhjYWxsYmFja3N1Y2NlZWRlZGRldGFpbHMYzL6ToAEg'
    'ASgLMiAubGFtYmRhLkNhbGxiYWNrU3VjY2VlZGVkRGV0YWlsc1IYY2FsbGJhY2tzdWNjZWVkZW'
    'RkZXRhaWxzElwKF2NhbGxiYWNrdGltZWRvdXRkZXRhaWxzGLzhuxsgASgLMh8ubGFtYmRhLkNh'
    'bGxiYWNrVGltZWRPdXREZXRhaWxzUhdjYWxsYmFja3RpbWVkb3V0ZGV0YWlscxJmChpjaGFpbm'
    'VkaW52b2tlZmFpbGVkZGV0YWlscxiR/5ngASABKAsyIi5sYW1iZGEuQ2hhaW5lZEludm9rZUZh'
    'aWxlZERldGFpbHNSGmNoYWluZWRpbnZva2VmYWlsZWRkZXRhaWxzEmgKG2NoYWluZWRpbnZva2'
    'VzdGFydGVkZGV0YWlscxjj2odlIAEoCzIjLmxhbWJkYS5DaGFpbmVkSW52b2tlU3RhcnRlZERl'
    'dGFpbHNSG2NoYWluZWRpbnZva2VzdGFydGVkZGV0YWlscxJoChtjaGFpbmVkaW52b2tlc3RvcH'
    'BlZGRldGFpbHMYk7qvEiABKAsyIy5sYW1iZGEuQ2hhaW5lZEludm9rZVN0b3BwZWREZXRhaWxz'
    'UhtjaGFpbmVkaW52b2tlc3RvcHBlZGRldGFpbHMSbwodY2hhaW5lZGludm9rZXN1Y2NlZWRlZG'
    'RldGFpbHMYs+/T9wEgASgLMiUubGFtYmRhLkNoYWluZWRJbnZva2VTdWNjZWVkZWREZXRhaWxz'
    'Uh1jaGFpbmVkaW52b2tlc3VjY2VlZGVkZGV0YWlscxJrChxjaGFpbmVkaW52b2tldGltZWRvdX'
    'RkZXRhaWxzGImIoHUgASgLMiQubGFtYmRhLkNoYWluZWRJbnZva2VUaW1lZE91dERldGFpbHNS'
    'HGNoYWluZWRpbnZva2V0aW1lZG91dGRldGFpbHMSVAoUY29udGV4dGZhaWxlZGRldGFpbHMYqo'
    'Sh1gEgASgLMhwubGFtYmRhLkNvbnRleHRGYWlsZWREZXRhaWxzUhRjb250ZXh0ZmFpbGVkZGV0'
    'YWlscxJXChVjb250ZXh0c3RhcnRlZGRldGFpbHMYquyeogEgASgLMh0ubGFtYmRhLkNvbnRleH'
    'RTdGFydGVkRGV0YWlsc1IVY29udGV4dHN0YXJ0ZWRkZXRhaWxzEl0KF2NvbnRleHRzdWNjZWVk'
    'ZWRkZXRhaWxzGMqb9oQBIAEoCzIfLmxhbWJkYS5Db250ZXh0U3VjY2VlZGVkRGV0YWlsc1IXY2'
    '9udGV4dHN1Y2NlZWRlZGRldGFpbHMSHAoHZXZlbnRpZBjTlt2zASABKAVSB2V2ZW50aWQSKQoO'
    'ZXZlbnR0aW1lc3RhbXAYjrmIWCABKAlSDmV2ZW50dGltZXN0YW1wEjMKCWV2ZW50dHlwZRjooM'
    'vfASABKA4yES5sYW1iZGEuRXZlbnRUeXBlUglldmVudHR5cGUSWQoWZXhlY3V0aW9uZmFpbGVk'
    'ZGV0YWlscxiz59tUIAEoCzIeLmxhbWJkYS5FeGVjdXRpb25GYWlsZWREZXRhaWxzUhZleGVjdX'
    'Rpb25mYWlsZWRkZXRhaWxzEl0KF2V4ZWN1dGlvbnN0YXJ0ZWRkZXRhaWxzGJXv8N0BIAEoCzIf'
    'LmxhbWJkYS5FeGVjdXRpb25TdGFydGVkRGV0YWlsc1IXZXhlY3V0aW9uc3RhcnRlZGRldGFpbH'
    'MSXAoXZXhlY3V0aW9uc3RvcHBlZGRldGFpbHMYpc3EZCABKAsyHy5sYW1iZGEuRXhlY3V0aW9u'
    'U3RvcHBlZERldGFpbHNSF2V4ZWN1dGlvbnN0b3BwZWRkZXRhaWxzEmMKGWV4ZWN1dGlvbnN1Y2'
    'NlZWRlZGRldGFpbHMY4e+Y7QEgASgLMiEubGFtYmRhLkV4ZWN1dGlvblN1Y2NlZWRlZERldGFp'
    'bHNSGWV4ZWN1dGlvbnN1Y2NlZWRlZGRldGFpbHMSYAoYZXhlY3V0aW9udGltZWRvdXRkZXRhaW'
    'xzGOOE860BIAEoCzIgLmxhbWJkYS5FeGVjdXRpb25UaW1lZE91dERldGFpbHNSGGV4ZWN1dGlv'
    'bnRpbWVkb3V0ZGV0YWlscxISCgJpZBiB8qK3ASABKAlSAmlkEmUKGmludm9jYXRpb25jb21wbG'
    'V0ZWRkZXRhaWxzGNuCsmAgASgLMiIubGFtYmRhLkludm9jYXRpb25Db21wbGV0ZWREZXRhaWxz'
    'UhppbnZvY2F0aW9uY29tcGxldGVkZGV0YWlscxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEh0KCH'
    'BhcmVudGlkGInpj1cgASgJUghwYXJlbnRpZBJKChFzdGVwZmFpbGVkZGV0YWlscxjJ8N4sIAEo'
    'CzIZLmxhbWJkYS5TdGVwRmFpbGVkRGV0YWlsc1IRc3RlcGZhaWxlZGRldGFpbHMSTgoSc3RlcH'
    'N0YXJ0ZWRkZXRhaWxzGKvA3fIBIAEoCzIaLmxhbWJkYS5TdGVwU3RhcnRlZERldGFpbHNSEnN0'
    'ZXBzdGFydGVkZGV0YWlscxJTChRzdGVwc3VjY2VlZGVkZGV0YWlscxjb8MdRIAEoCzIcLmxhbW'
    'JkYS5TdGVwU3VjY2VlZGVkRGV0YWlsc1IUc3RlcHN1Y2NlZWRlZGRldGFpbHMSGwoHc3VidHlw'
    'ZRi+1flIIAEoCVIHc3VidHlwZRJTChR3YWl0Y2FuY2VsbGVkZGV0YWlscxjEgo9hIAEoCzIcLm'
    'xhbWJkYS5XYWl0Q2FuY2VsbGVkRGV0YWlsc1IUd2FpdGNhbmNlbGxlZGRldGFpbHMSTgoSd2Fp'
    'dHN0YXJ0ZWRkZXRhaWxzGL62kIIBIAEoCzIaLmxhbWJkYS5XYWl0U3RhcnRlZERldGFpbHNSEn'
    'dhaXRzdGFydGVkZGV0YWlscxJUChR3YWl0c3VjY2VlZGVkZGV0YWlscxju/amCASABKAsyHC5s'
    'YW1iZGEuV2FpdFN1Y2NlZWRlZERldGFpbHNSFHdhaXRzdWNjZWVkZWRkZXRhaWxz');

@$core.Deprecated('Use eventErrorDescriptor instead')
const EventError$json = {
  '1': 'EventError',
  '2': [
    {
      '1': 'payload',
      '3': 6526790,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '10': 'payload'
    },
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `EventError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventErrorDescriptor = $convert.base64Decode(
    'CgpFdmVudEVycm9yEjAKB3BheWxvYWQYxq6OAyABKAsyEy5sYW1iZGEuRXJyb3JPYmplY3RSB3'
    'BheWxvYWQSHwoJdHJ1bmNhdGVkGMrv2EggASgIUgl0cnVuY2F0ZWQ=');

@$core.Deprecated('Use eventInputDescriptor instead')
const EventInput$json = {
  '1': 'EventInput',
  '2': [
    {'1': 'payload', '3': 6526790, '4': 1, '5': 9, '10': 'payload'},
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `EventInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventInputDescriptor = $convert.base64Decode(
    'CgpFdmVudElucHV0EhsKB3BheWxvYWQYxq6OAyABKAlSB3BheWxvYWQSHwoJdHJ1bmNhdGVkGM'
    'rv2EggASgIUgl0cnVuY2F0ZWQ=');

@$core.Deprecated('Use eventResultDescriptor instead')
const EventResult$json = {
  '1': 'EventResult',
  '2': [
    {'1': 'payload', '3': 6526790, '4': 1, '5': 9, '10': 'payload'},
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `EventResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventResultDescriptor = $convert.base64Decode(
    'CgtFdmVudFJlc3VsdBIbCgdwYXlsb2FkGMaujgMgASgJUgdwYXlsb2FkEh8KCXRydW5jYXRlZB'
    'jK79hIIAEoCFIJdHJ1bmNhdGVk');

@$core.Deprecated('Use eventSourceMappingConfigurationDescriptor instead')
const EventSourceMappingConfiguration$json = {
  '1': 'EventSourceMappingConfiguration',
  '2': [
    {
      '1': 'amazonmanagedkafkaeventsourceconfig',
      '3': 60136380,
      '4': 1,
      '5': 11,
      '6': '.lambda.AmazonManagedKafkaEventSourceConfig',
      '10': 'amazonmanagedkafkaeventsourceconfig'
    },
    {'1': 'batchsize', '3': 318039259, '4': 1, '5': 5, '10': 'batchsize'},
    {
      '1': 'bisectbatchonfunctionerror',
      '3': 276143707,
      '4': 1,
      '5': 8,
      '10': 'bisectbatchonfunctionerror'
    },
    {
      '1': 'destinationconfig',
      '3': 184834158,
      '4': 1,
      '5': 11,
      '6': '.lambda.DestinationConfig',
      '10': 'destinationconfig'
    },
    {
      '1': 'documentdbeventsourceconfig',
      '3': 173060622,
      '4': 1,
      '5': 11,
      '6': '.lambda.DocumentDBEventSourceConfig',
      '10': 'documentdbeventsourceconfig'
    },
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {
      '1': 'eventsourcemappingarn',
      '3': 176331828,
      '4': 1,
      '5': 9,
      '10': 'eventsourcemappingarn'
    },
    {
      '1': 'filtercriteria',
      '3': 439219323,
      '4': 1,
      '5': 11,
      '6': '.lambda.FilterCriteria',
      '10': 'filtercriteria'
    },
    {
      '1': 'filtercriteriaerror',
      '3': 355388477,
      '4': 1,
      '5': 11,
      '6': '.lambda.FilterCriteriaError',
      '10': 'filtercriteriaerror'
    },
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {
      '1': 'functionresponsetypes',
      '3': 382292260,
      '4': 3,
      '5': 14,
      '6': '.lambda.FunctionResponseType',
      '10': 'functionresponsetypes'
    },
    {'1': 'kmskeyarn', '3': 117627377, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'lastprocessingresult',
      '3': 278539700,
      '4': 1,
      '5': 9,
      '10': 'lastprocessingresult'
    },
    {
      '1': 'loggingconfig',
      '3': 424359625,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventSourceMappingLoggingConfig',
      '10': 'loggingconfig'
    },
    {
      '1': 'maximumbatchingwindowinseconds',
      '3': 346663320,
      '4': 1,
      '5': 5,
      '10': 'maximumbatchingwindowinseconds'
    },
    {
      '1': 'maximumrecordageinseconds',
      '3': 102344982,
      '4': 1,
      '5': 5,
      '10': 'maximumrecordageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
    {
      '1': 'metricsconfig',
      '3': 412971857,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventSourceMappingMetricsConfig',
      '10': 'metricsconfig'
    },
    {
      '1': 'parallelizationfactor',
      '3': 400337694,
      '4': 1,
      '5': 5,
      '10': 'parallelizationfactor'
    },
    {
      '1': 'provisionedpollerconfig',
      '3': 275676602,
      '4': 1,
      '5': 11,
      '6': '.lambda.ProvisionedPollerConfig',
      '10': 'provisionedpollerconfig'
    },
    {'1': 'queues', '3': 222519730, '4': 3, '5': 9, '10': 'queues'},
    {
      '1': 'scalingconfig',
      '3': 392871661,
      '4': 1,
      '5': 11,
      '6': '.lambda.ScalingConfig',
      '10': 'scalingconfig'
    },
    {
      '1': 'selfmanagedeventsource',
      '3': 283601786,
      '4': 1,
      '5': 11,
      '6': '.lambda.SelfManagedEventSource',
      '10': 'selfmanagedeventsource'
    },
    {
      '1': 'selfmanagedkafkaeventsourceconfig',
      '3': 322222578,
      '4': 1,
      '5': 11,
      '6': '.lambda.SelfManagedKafkaEventSourceConfig',
      '10': 'selfmanagedkafkaeventsourceconfig'
    },
    {
      '1': 'sourceaccessconfigurations',
      '3': 371593554,
      '4': 3,
      '5': 11,
      '6': '.lambda.SourceAccessConfiguration',
      '10': 'sourceaccessconfigurations'
    },
    {
      '1': 'startingposition',
      '3': 428771919,
      '4': 1,
      '5': 14,
      '6': '.lambda.EventSourcePosition',
      '10': 'startingposition'
    },
    {
      '1': 'startingpositiontimestamp',
      '3': 144323607,
      '4': 1,
      '5': 9,
      '10': 'startingpositiontimestamp'
    },
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {
      '1': 'statetransitionreason',
      '3': 79099714,
      '4': 1,
      '5': 9,
      '10': 'statetransitionreason'
    },
    {'1': 'topics', '3': 219850038, '4': 3, '5': 9, '10': 'topics'},
    {
      '1': 'tumblingwindowinseconds',
      '3': 372493124,
      '4': 1,
      '5': 5,
      '10': 'tumblingwindowinseconds'
    },
    {'1': 'uuid', '3': 91981875, '4': 1, '5': 9, '10': 'uuid'},
  ],
};

/// Descriptor for `EventSourceMappingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventSourceMappingConfigurationDescriptor = $convert.base64Decode(
    'Ch9FdmVudFNvdXJjZU1hcHBpbmdDb25maWd1cmF0aW9uEoABCiNhbWF6b25tYW5hZ2Vka2Fma2'
    'FldmVudHNvdXJjZWNvbmZpZxi8t9YcIAEoCzIrLmxhbWJkYS5BbWF6b25NYW5hZ2VkS2Fma2FF'
    'dmVudFNvdXJjZUNvbmZpZ1IjYW1hem9ubWFuYWdlZGthZmthZXZlbnRzb3VyY2Vjb25maWcSIA'
    'oJYmF0Y2hzaXplGNvJ05cBIAEoBVIJYmF0Y2hzaXplEkIKGmJpc2VjdGJhdGNob25mdW5jdGlv'
    'bmVycm9yGNu81oMBIAEoCFIaYmlzZWN0YmF0Y2hvbmZ1bmN0aW9uZXJyb3ISSgoRZGVzdGluYX'
    'Rpb25jb25maWcY7rCRWCABKAsyGS5sYW1iZGEuRGVzdGluYXRpb25Db25maWdSEWRlc3RpbmF0'
    'aW9uY29uZmlnEmgKG2RvY3VtZW50ZGJldmVudHNvdXJjZWNvbmZpZxiO5MJSIAEoCzIjLmxhbW'
    'JkYS5Eb2N1bWVudERCRXZlbnRTb3VyY2VDb25maWdSG2RvY3VtZW50ZGJldmVudHNvdXJjZWNv'
    'bmZpZxIqCg5ldmVudHNvdXJjZWFybhjGyoqSASABKAlSDmV2ZW50c291cmNlYXJuEjcKFWV2ZW'
    '50c291cmNlbWFwcGluZ2Fybhi0uIpUIAEoCVIVZXZlbnRzb3VyY2VtYXBwaW5nYXJuEkIKDmZp'
    'bHRlcmNyaXRlcmlhGPvot9EBIAEoCzIWLmxhbWJkYS5GaWx0ZXJDcml0ZXJpYVIOZmlsdGVyY3'
    'JpdGVyaWESUQoTZmlsdGVyY3JpdGVyaWFlcnJvchi9mLupASABKAsyGy5sYW1iZGEuRmlsdGVy'
    'Q3JpdGVyaWFFcnJvclITZmlsdGVyY3JpdGVyaWFlcnJvchIkCgtmdW5jdGlvbmFybhjxyY62AS'
    'ABKAlSC2Z1bmN0aW9uYXJuElYKFWZ1bmN0aW9ucmVzcG9uc2V0eXBlcxikoqW2ASADKA4yHC5s'
    'YW1iZGEuRnVuY3Rpb25SZXNwb25zZVR5cGVSFWZ1bmN0aW9ucmVzcG9uc2V0eXBlcxIfCglrbX'
    'NrZXlhcm4Y8bOLOCABKAlSCWttc2tleWFybhImCgxsYXN0bW9kaWZpZWQYp5z8zgEgASgJUgxs'
    'YXN0bW9kaWZpZWQSNgoUbGFzdHByb2Nlc3NpbmdyZXN1bHQYtNvohAEgASgJUhRsYXN0cHJvY2'
    'Vzc2luZ3Jlc3VsdBJRCg1sb2dnaW5nY29uZmlnGMntrMoBIAEoCzInLmxhbWJkYS5FdmVudFNv'
    'dXJjZU1hcHBpbmdMb2dnaW5nQ29uZmlnUg1sb2dnaW5nY29uZmlnEkoKHm1heGltdW1iYXRjaG'
    'luZ3dpbmRvd2luc2Vjb25kcxiY06alASABKAVSHm1heGltdW1iYXRjaGluZ3dpbmRvd2luc2Vj'
    'b25kcxI/ChltYXhpbXVtcmVjb3JkYWdlaW5zZWNvbmRzGJbS5jAgASgFUhltYXhpbXVtcmVjb3'
    'JkYWdlaW5zZWNvbmRzEjUKFG1heGltdW1yZXRyeWF0dGVtcHRzGMCouTUgASgFUhRtYXhpbXVt'
    'cmV0cnlhdHRlbXB0cxJRCg1tZXRyaWNzY29uZmlnGNHm9cQBIAEoCzInLmxhbWJkYS5FdmVudF'
    'NvdXJjZU1hcHBpbmdNZXRyaWNzQ29uZmlnUg1tZXRyaWNzY29uZmlnEjgKFXBhcmFsbGVsaXph'
    'dGlvbmZhY3Rvchie1vK+ASABKAVSFXBhcmFsbGVsaXphdGlvbmZhY3RvchJdChdwcm92aXNpb2'
    '5lZHBvbGxlcmNvbmZpZxi6+7mDASABKAsyHy5sYW1iZGEuUHJvdmlzaW9uZWRQb2xsZXJDb25m'
    'aWdSF3Byb3Zpc2lvbmVkcG9sbGVyY29uZmlnEhkKBnF1ZXVlcxiyw41qIAMoCVIGcXVldWVzEj'
    '8KDXNjYWxpbmdjb25maWcY7f2quwEgASgLMhUubGFtYmRhLlNjYWxpbmdDb25maWdSDXNjYWxp'
    'bmdjb25maWcSWgoWc2VsZm1hbmFnZWRldmVudHNvdXJjZRj61p2HASABKAsyHi5sYW1iZGEuU2'
    'VsZk1hbmFnZWRFdmVudFNvdXJjZVIWc2VsZm1hbmFnZWRldmVudHNvdXJjZRJ7CiFzZWxmbWFu'
    'YWdlZGthZmthZXZlbnRzb3VyY2Vjb25maWcY8vPSmQEgASgLMikubGFtYmRhLlNlbGZNYW5hZ2'
    'VkS2Fma2FFdmVudFNvdXJjZUNvbmZpZ1Ihc2VsZm1hbmFnZWRrYWZrYWV2ZW50c291cmNlY29u'
    'ZmlnEmUKGnNvdXJjZWFjY2Vzc2NvbmZpZ3VyYXRpb25zGNKimLEBIAMoCzIhLmxhbWJkYS5Tb3'
    'VyY2VBY2Nlc3NDb25maWd1cmF0aW9uUhpzb3VyY2VhY2Nlc3Njb25maWd1cmF0aW9ucxJLChBz'
    'dGFydGluZ3Bvc2l0aW9uGM+UuswBIAEoDjIbLmxhbWJkYS5FdmVudFNvdXJjZVBvc2l0aW9uUh'
    'BzdGFydGluZ3Bvc2l0aW9uEj8KGXN0YXJ0aW5ncG9zaXRpb250aW1lc3RhbXAYl+joRCABKAlS'
    'GXN0YXJ0aW5ncG9zaXRpb250aW1lc3RhbXASGAoFc3RhdGUYl8my7wEgASgJUgVzdGF0ZRI3Ch'
    'VzdGF0ZXRyYW5zaXRpb25yZWFzb24Ywu7bJSABKAlSFXN0YXRldHJhbnNpdGlvbnJlYXNvbhIZ'
    'CgZ0b3BpY3MYtsrqaCADKAlSBnRvcGljcxI8Chd0dW1ibGluZ3dpbmRvd2luc2Vjb25kcxjEls'
    '+xASABKAVSF3R1bWJsaW5nd2luZG93aW5zZWNvbmRzEhUKBHV1aWQYs5DuKyABKAlSBHV1aWQ=');

@$core.Deprecated('Use eventSourceMappingLoggingConfigDescriptor instead')
const EventSourceMappingLoggingConfig$json = {
  '1': 'EventSourceMappingLoggingConfig',
  '2': [
    {
      '1': 'systemloglevel',
      '3': 530478525,
      '4': 1,
      '5': 14,
      '6': '.lambda.EventSourceMappingSystemLogLevel',
      '10': 'systemloglevel'
    },
  ],
};

/// Descriptor for `EventSourceMappingLoggingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventSourceMappingLoggingConfigDescriptor =
    $convert.base64Decode(
        'Ch9FdmVudFNvdXJjZU1hcHBpbmdMb2dnaW5nQ29uZmlnElQKDnN5c3RlbWxvZ2xldmVsGL3r+f'
        'wBIAEoDjIoLmxhbWJkYS5FdmVudFNvdXJjZU1hcHBpbmdTeXN0ZW1Mb2dMZXZlbFIOc3lzdGVt'
        'bG9nbGV2ZWw=');

@$core.Deprecated('Use eventSourceMappingMetricsConfigDescriptor instead')
const EventSourceMappingMetricsConfig$json = {
  '1': 'EventSourceMappingMetricsConfig',
  '2': [
    {
      '1': 'metrics',
      '3': 436365847,
      '4': 3,
      '5': 14,
      '6': '.lambda.EventSourceMappingMetric',
      '10': 'metrics'
    },
  ],
};

/// Descriptor for `EventSourceMappingMetricsConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventSourceMappingMetricsConfigDescriptor =
    $convert.base64Decode(
        'Ch9FdmVudFNvdXJjZU1hcHBpbmdNZXRyaWNzQ29uZmlnEj4KB21ldHJpY3MYl9SJ0AEgAygOMi'
        'AubGFtYmRhLkV2ZW50U291cmNlTWFwcGluZ01ldHJpY1IHbWV0cmljcw==');

@$core.Deprecated('Use executionDescriptor instead')
const Execution$json = {
  '1': 'Execution',
  '2': [
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {
      '1': 'durableexecutionname',
      '3': 251828526,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionname'
    },
    {'1': 'endtimestamp', '3': 340967267, '4': 1, '5': 9, '10': 'endtimestamp'},
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {
      '1': 'starttimestamp',
      '3': 392954318,
      '4': 1,
      '5': 9,
      '10': 'starttimestamp'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.lambda.ExecutionStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `Execution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionDescriptor = $convert.base64Decode(
    'CglFeGVjdXRpb24SNAoTZHVyYWJsZWV4ZWN1dGlvbmFybhiKzbeAASABKAlSE2R1cmFibGVleG'
    'VjdXRpb25hcm4SNQoUZHVyYWJsZWV4ZWN1dGlvbm5hbWUYrrKKeCABKAlSFGR1cmFibGVleGVj'
    'dXRpb25uYW1lEiYKDGVuZHRpbWVzdGFtcBjj/sqiASABKAlSDGVuZHRpbWVzdGFtcBIkCgtmdW'
    '5jdGlvbmFybhjxyY62ASABKAlSC2Z1bmN0aW9uYXJuEioKDnN0YXJ0dGltZXN0YW1wGM6DsLsB'
    'IAEoCVIOc3RhcnR0aW1lc3RhbXASMgoGc3RhdHVzGJDk+wIgASgOMhcubGFtYmRhLkV4ZWN1dG'
    'lvblN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use executionDetailsDescriptor instead')
const ExecutionDetails$json = {
  '1': 'ExecutionDetails',
  '2': [
    {'1': 'inputpayload', '3': 228187532, '4': 1, '5': 9, '10': 'inputpayload'},
  ],
};

/// Descriptor for `ExecutionDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionDetailsDescriptor = $convert.base64Decode(
    'ChBFeGVjdXRpb25EZXRhaWxzEiUKDGlucHV0cGF5bG9hZBiMu+dsIAEoCVIMaW5wdXRwYXlsb2'
    'Fk');

@$core.Deprecated('Use executionFailedDetailsDescriptor instead')
const ExecutionFailedDetails$json = {
  '1': 'ExecutionFailedDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `ExecutionFailedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionFailedDetailsDescriptor =
    $convert.base64Decode(
        'ChZFeGVjdXRpb25GYWlsZWREZXRhaWxzEiwKBWVycm9yGPK5tpwBIAEoCzISLmxhbWJkYS5Fdm'
        'VudEVycm9yUgVlcnJvcg==');

@$core.Deprecated('Use executionStartedDetailsDescriptor instead')
const ExecutionStartedDetails$json = {
  '1': 'ExecutionStartedDetails',
  '2': [
    {
      '1': 'executiontimeout',
      '3': 502184601,
      '4': 1,
      '5': 5,
      '10': 'executiontimeout'
    },
    {
      '1': 'input',
      '3': 529785116,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventInput',
      '10': 'input'
    },
  ],
};

/// Descriptor for `ExecutionStartedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionStartedDetailsDescriptor = $convert.base64Decode(
    'ChdFeGVjdXRpb25TdGFydGVkRGV0YWlscxIuChBleGVjdXRpb250aW1lb3V0GJn1uu8BIAEoBV'
    'IQZXhlY3V0aW9udGltZW91dBIsCgVpbnB1dBicws/8ASABKAsyEi5sYW1iZGEuRXZlbnRJbnB1'
    'dFIFaW5wdXQ=');

@$core.Deprecated('Use executionStoppedDetailsDescriptor instead')
const ExecutionStoppedDetails$json = {
  '1': 'ExecutionStoppedDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `ExecutionStoppedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionStoppedDetailsDescriptor =
    $convert.base64Decode(
        'ChdFeGVjdXRpb25TdG9wcGVkRGV0YWlscxIsCgVlcnJvchjyubacASABKAsyEi5sYW1iZGEuRX'
        'ZlbnRFcnJvclIFZXJyb3I=');

@$core.Deprecated('Use executionSucceededDetailsDescriptor instead')
const ExecutionSucceededDetails$json = {
  '1': 'ExecutionSucceededDetails',
  '2': [
    {
      '1': 'result',
      '3': 273346629,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventResult',
      '10': 'result'
    },
  ],
};

/// Descriptor for `ExecutionSucceededDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionSucceededDetailsDescriptor =
    $convert.base64Decode(
        'ChlFeGVjdXRpb25TdWNjZWVkZWREZXRhaWxzEi8KBnJlc3VsdBjF4KuCASABKAsyEy5sYW1iZG'
        'EuRXZlbnRSZXN1bHRSBnJlc3VsdA==');

@$core.Deprecated('Use executionTimedOutDetailsDescriptor instead')
const ExecutionTimedOutDetails$json = {
  '1': 'ExecutionTimedOutDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `ExecutionTimedOutDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionTimedOutDetailsDescriptor =
    $convert.base64Decode(
        'ChhFeGVjdXRpb25UaW1lZE91dERldGFpbHMSLAoFZXJyb3IY8rm2nAEgASgLMhIubGFtYmRhLk'
        'V2ZW50RXJyb3JSBWVycm9y');

@$core.Deprecated('Use fileSystemConfigDescriptor instead')
const FileSystemConfig$json = {
  '1': 'FileSystemConfig',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'localmountpath',
      '3': 69699719,
      '4': 1,
      '5': 9,
      '10': 'localmountpath'
    },
  ],
};

/// Descriptor for `FileSystemConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fileSystemConfigDescriptor = $convert.base64Decode(
    'ChBGaWxlU3lzdGVtQ29uZmlnEhQKA2Fybhidm+2/ASABKAlSA2FybhIpCg5sb2NhbG1vdW50cG'
    'F0aBiHkZ4hIAEoCVIObG9jYWxtb3VudHBhdGg=');

@$core.Deprecated('Use filterDescriptor instead')
const Filter$json = {
  '1': 'Filter',
  '2': [
    {'1': 'pattern', '3': 345055434, '4': 1, '5': 9, '10': 'pattern'},
  ],
};

/// Descriptor for `Filter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterDescriptor = $convert
    .base64Decode('CgZGaWx0ZXISHAoHcGF0dGVybhjKwcSkASABKAlSB3BhdHRlcm4=');

@$core.Deprecated('Use filterCriteriaDescriptor instead')
const FilterCriteria$json = {
  '1': 'FilterCriteria',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.lambda.Filter',
      '10': 'filters'
    },
  ],
};

/// Descriptor for `FilterCriteria`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterCriteriaDescriptor = $convert.base64Decode(
    'Cg5GaWx0ZXJDcml0ZXJpYRIrCgdmaWx0ZXJzGO3N6lkgAygLMg4ubGFtYmRhLkZpbHRlclIHZm'
    'lsdGVycw==');

@$core.Deprecated('Use filterCriteriaErrorDescriptor instead')
const FilterCriteriaError$json = {
  '1': 'FilterCriteriaError',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FilterCriteriaError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterCriteriaErrorDescriptor = $convert.base64Decode(
    'ChNGaWx0ZXJDcml0ZXJpYUVycm9yEh8KCWVycm9yY29kZRiZ1sMQIAEoCVIJZXJyb3Jjb2RlEh'
    'sKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use functionCodeDescriptor instead')
const FunctionCode$json = {
  '1': 'FunctionCode',
  '2': [
    {'1': 'imageuri', '3': 412238461, '4': 1, '5': 9, '10': 'imageuri'},
    {'1': 's3bucket', '3': 114031434, '4': 1, '5': 9, '10': 's3bucket'},
    {'1': 's3key', '3': 490298907, '4': 1, '5': 9, '10': 's3key'},
    {
      '1': 's3objectversion',
      '3': 194809669,
      '4': 1,
      '5': 9,
      '10': 's3objectversion'
    },
    {
      '1': 'sourcekmskeyarn',
      '3': 203651164,
      '4': 1,
      '5': 9,
      '10': 'sourcekmskeyarn'
    },
    {'1': 'zipfile', '3': 2519299, '4': 1, '5': 12, '10': 'zipfile'},
  ],
};

/// Descriptor for `FunctionCode`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionCodeDescriptor = $convert.base64Decode(
    'CgxGdW5jdGlvbkNvZGUSHgoIaW1hZ2V1cmkY/YTJxAEgASgJUghpbWFnZXVyaRIdCghzM2J1Y2'
    'tldBjK9q82IAEoCVIIczNidWNrZXQSGAoFczNrZXkYm7zl6QEgASgJUgVzM2tleRIrCg9zM29i'
    'amVjdHZlcnNpb24YxZ7yXCABKAlSD3Mzb2JqZWN0dmVyc2lvbhIrCg9zb3VyY2VrbXNrZXlhcm'
    '4Y3PCNYSABKAlSD3NvdXJjZWttc2tleWFybhIbCgd6aXBmaWxlGIPimQEgASgMUgd6aXBmaWxl');

@$core.Deprecated('Use functionCodeLocationDescriptor instead')
const FunctionCodeLocation$json = {
  '1': 'FunctionCodeLocation',
  '2': [
    {'1': 'imageuri', '3': 412238461, '4': 1, '5': 9, '10': 'imageuri'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'repositorytype',
      '3': 224453310,
      '4': 1,
      '5': 9,
      '10': 'repositorytype'
    },
    {
      '1': 'resolvedimageuri',
      '3': 456188253,
      '4': 1,
      '5': 9,
      '10': 'resolvedimageuri'
    },
    {
      '1': 'sourcekmskeyarn',
      '3': 203651164,
      '4': 1,
      '5': 9,
      '10': 'sourcekmskeyarn'
    },
  ],
};

/// Descriptor for `FunctionCodeLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionCodeLocationDescriptor = $convert.base64Decode(
    'ChRGdW5jdGlvbkNvZGVMb2NhdGlvbhIeCghpbWFnZXVyaRj9hMnEASABKAlSCGltYWdldXJpEh'
    '4KCGxvY2F0aW9uGMebgt4BIAEoCVIIbG9jYXRpb24SKQoOcmVwb3NpdG9yeXR5cGUYvsWDayAB'
    'KAlSDnJlcG9zaXRvcnl0eXBlEi4KEHJlc29sdmVkaW1hZ2V1cmkY3cLD2QEgASgJUhByZXNvbH'
    'ZlZGltYWdldXJpEisKD3NvdXJjZWttc2tleWFybhjc8I1hIAEoCVIPc291cmNla21za2V5YXJu');

@$core.Deprecated('Use functionConfigurationDescriptor instead')
const FunctionConfiguration$json = {
  '1': 'FunctionConfiguration',
  '2': [
    {
      '1': 'architectures',
      '3': 530490948,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'architectures'
    },
    {
      '1': 'capacityproviderconfig',
      '3': 52030623,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderConfig',
      '10': 'capacityproviderconfig'
    },
    {'1': 'codesha256', '3': 46450860, '4': 1, '5': 9, '10': 'codesha256'},
    {'1': 'codesize', '3': 74450158, '4': 1, '5': 3, '10': 'codesize'},
    {'1': 'configsha256', '3': 145714121, '4': 1, '5': 9, '10': 'configsha256'},
    {
      '1': 'deadletterconfig',
      '3': 79786642,
      '4': 1,
      '5': 11,
      '6': '.lambda.DeadLetterConfig',
      '10': 'deadletterconfig'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'durableconfig',
      '3': 206326279,
      '4': 1,
      '5': 11,
      '6': '.lambda.DurableConfig',
      '10': 'durableconfig'
    },
    {
      '1': 'environment',
      '3': 119823003,
      '4': 1,
      '5': 11,
      '6': '.lambda.EnvironmentResponse',
      '10': 'environment'
    },
    {
      '1': 'ephemeralstorage',
      '3': 365965382,
      '4': 1,
      '5': 11,
      '6': '.lambda.EphemeralStorage',
      '10': 'ephemeralstorage'
    },
    {
      '1': 'filesystemconfigs',
      '3': 490453750,
      '4': 3,
      '5': 11,
      '6': '.lambda.FileSystemConfig',
      '10': 'filesystemconfigs'
    },
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'handler', '3': 81160724, '4': 1, '5': 9, '10': 'handler'},
    {
      '1': 'imageconfigresponse',
      '3': 319434068,
      '4': 1,
      '5': 11,
      '6': '.lambda.ImageConfigResponse',
      '10': 'imageconfigresponse'
    },
    {'1': 'kmskeyarn', '3': 117627377, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'lastupdatestatus',
      '3': 269946843,
      '4': 1,
      '5': 14,
      '6': '.lambda.LastUpdateStatus',
      '10': 'lastupdatestatus'
    },
    {
      '1': 'lastupdatestatusreason',
      '3': 165558159,
      '4': 1,
      '5': 9,
      '10': 'lastupdatestatusreason'
    },
    {
      '1': 'lastupdatestatusreasoncode',
      '3': 275687348,
      '4': 1,
      '5': 14,
      '6': '.lambda.LastUpdateStatusReasonCode',
      '10': 'lastupdatestatusreasoncode'
    },
    {
      '1': 'layers',
      '3': 478144896,
      '4': 3,
      '5': 11,
      '6': '.lambda.Layer',
      '10': 'layers'
    },
    {
      '1': 'loggingconfig',
      '3': 424359625,
      '4': 1,
      '5': 11,
      '6': '.lambda.LoggingConfig',
      '10': 'loggingconfig'
    },
    {'1': 'masterarn', '3': 74605927, '4': 1, '5': 9, '10': 'masterarn'},
    {'1': 'memorysize', '3': 55523120, '4': 1, '5': 5, '10': 'memorysize'},
    {
      '1': 'packagetype',
      '3': 517524132,
      '4': 1,
      '5': 14,
      '6': '.lambda.PackageType',
      '10': 'packagetype'
    },
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'role', '3': 271285818, '4': 1, '5': 9, '10': 'role'},
    {
      '1': 'runtime',
      '3': 359311308,
      '4': 1,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'runtime'
    },
    {
      '1': 'runtimeversionconfig',
      '3': 486723720,
      '4': 1,
      '5': 11,
      '6': '.lambda.RuntimeVersionConfig',
      '10': 'runtimeversionconfig'
    },
    {
      '1': 'signingjobarn',
      '3': 343397691,
      '4': 1,
      '5': 9,
      '10': 'signingjobarn'
    },
    {
      '1': 'signingprofileversionarn',
      '3': 432885567,
      '4': 1,
      '5': 9,
      '10': 'signingprofileversionarn'
    },
    {
      '1': 'snapstart',
      '3': 283273032,
      '4': 1,
      '5': 11,
      '6': '.lambda.SnapStartResponse',
      '10': 'snapstart'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.lambda.State',
      '10': 'state'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
    {
      '1': 'statereasoncode',
      '3': 319263936,
      '4': 1,
      '5': 14,
      '6': '.lambda.StateReasonCode',
      '10': 'statereasoncode'
    },
    {
      '1': 'tenancyconfig',
      '3': 215700986,
      '4': 1,
      '5': 11,
      '6': '.lambda.TenancyConfig',
      '10': 'tenancyconfig'
    },
    {'1': 'timeout', '3': 47808041, '4': 1, '5': 5, '10': 'timeout'},
    {
      '1': 'tracingconfig',
      '3': 19554860,
      '4': 1,
      '5': 11,
      '6': '.lambda.TracingConfigResponse',
      '10': 'tracingconfig'
    },
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
    {
      '1': 'vpcconfig',
      '3': 194980743,
      '4': 1,
      '5': 11,
      '6': '.lambda.VpcConfigResponse',
      '10': 'vpcconfig'
    },
  ],
};

/// Descriptor for `FunctionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionConfigurationDescriptor = $convert.base64Decode(
    'ChVGdW5jdGlvbkNvbmZpZ3VyYXRpb24SPgoNYXJjaGl0ZWN0dXJlcxjEzPr8ASADKA4yFC5sYW'
    '1iZGEuQXJjaGl0ZWN0dXJlUg1hcmNoaXRlY3R1cmVzElkKFmNhcGFjaXR5cHJvdmlkZXJjb25m'
    'aWcYn9nnGCABKAsyHi5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlckNvbmZpZ1IWY2FwYWNpdHlwcm'
    '92aWRlcmNvbmZpZxIhCgpjb2Rlc2hhMjU2GKyRkxYgASgJUgpjb2Rlc2hhMjU2Eh0KCGNvZGVz'
    'aXplGO6JwCMgASgDUghjb2Rlc2l6ZRIlCgxjb25maWdzaGEyNTYYyde9RSABKAlSDGNvbmZpZ3'
    'NoYTI1NhJHChBkZWFkbGV0dGVyY29uZmlnGJLlhSYgASgLMhgubGFtYmRhLkRlYWRMZXR0ZXJD'
    'b25maWdSEGRlYWRsZXR0ZXJjb25maWcSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaX'
    'B0aW9uEj4KDWR1cmFibGVjb25maWcYh5SxYiABKAsyFS5sYW1iZGEuRHVyYWJsZUNvbmZpZ1IN'
    'ZHVyYWJsZWNvbmZpZxJACgtlbnZpcm9ubWVudBibtZE5IAEoCzIbLmxhbWJkYS5FbnZpcm9ubW'
    'VudFJlc3BvbnNlUgtlbnZpcm9ubWVudBJIChBlcGhlbWVyYWxzdG9yYWdlGMbgwK4BIAEoCzIY'
    'LmxhbWJkYS5FcGhlbWVyYWxTdG9yYWdlUhBlcGhlbWVyYWxzdG9yYWdlEkoKEWZpbGVzeXN0ZW'
    '1jb25maWdzGPb17ukBIAMoCzIYLmxhbWJkYS5GaWxlU3lzdGVtQ29uZmlnUhFmaWxlc3lzdGVt'
    'Y29uZmlncxIkCgtmdW5jdGlvbmFybhjxyY62ASABKAlSC2Z1bmN0aW9uYXJuEiYKDGZ1bmN0aW'
    '9ubmFtZRijiL/fASABKAlSDGZ1bmN0aW9ubmFtZRIbCgdoYW5kbGVyGJTU2SYgASgJUgdoYW5k'
    'bGVyElEKE2ltYWdlY29uZmlncmVzcG9uc2UY1NqomAEgASgLMhsubGFtYmRhLkltYWdlQ29uZm'
    'lnUmVzcG9uc2VSE2ltYWdlY29uZmlncmVzcG9uc2USHwoJa21za2V5YXJuGPGzizggASgJUglr'
    'bXNrZXlhcm4SJgoMbGFzdG1vZGlmaWVkGKec/M4BIAEoCVIMbGFzdG1vZGlmaWVkEkgKEGxhc3'
    'R1cGRhdGVzdGF0dXMY25/cgAEgASgOMhgubGFtYmRhLkxhc3RVcGRhdGVTdGF0dXNSEGxhc3R1'
    'cGRhdGVzdGF0dXMSOQoWbGFzdHVwZGF0ZXN0YXR1c3JlYXNvbhiP7/hOIAEoCVIWbGFzdHVwZG'
    'F0ZXN0YXR1c3JlYXNvbhJmChpsYXN0dXBkYXRlc3RhdHVzcmVhc29uY29kZRi0z7qDASABKA4y'
    'Ii5sYW1iZGEuTGFzdFVwZGF0ZVN0YXR1c1JlYXNvbkNvZGVSGmxhc3R1cGRhdGVzdGF0dXNyZW'
    'Fzb25jb2RlEikKBmxheWVycxiA0//jASADKAsyDS5sYW1iZGEuTGF5ZXJSBmxheWVycxI/Cg1s'
    'b2dnaW5nY29uZmlnGMntrMoBIAEoCzIVLmxhbWJkYS5Mb2dnaW5nQ29uZmlnUg1sb2dnaW5nY2'
    '9uZmlnEh8KCW1hc3RlcmFybhjnyskjIAEoCVIJbWFzdGVyYXJuEiEKCm1lbW9yeXNpemUYsO68'
    'GiABKAVSCm1lbW9yeXNpemUSOQoLcGFja2FnZXR5cGUYpJXj9gEgASgOMhMubGFtYmRhLlBhY2'
    'thZ2VUeXBlUgtwYWNrYWdldHlwZRIiCgpyZXZpc2lvbmlkGIajnu4BIAEoCVIKcmV2aXNpb25p'
    'ZBIWCgRyb2xlGLr8rYEBIAEoCVIEcm9sZRItCgdydW50aW1lGMzPqqsBIAEoDjIPLmxhbWJkYS'
    '5SdW50aW1lUgdydW50aW1lElQKFHJ1bnRpbWV2ZXJzaW9uY29uZmlnGIihi+gBIAEoCzIcLmxh'
    'bWJkYS5SdW50aW1lVmVyc2lvbkNvbmZpZ1IUcnVudGltZXZlcnNpb25jb25maWcSKAoNc2lnbm'
    'luZ2pvYmFybhi7qt+jASABKAlSDXNpZ25pbmdqb2Jhcm4SPgoYc2lnbmluZ3Byb2ZpbGV2ZXJz'
    'aW9uYXJuGL+etc4BIAEoCVIYc2lnbmluZ3Byb2ZpbGV2ZXJzaW9uYXJuEjsKCXNuYXBzdGFydB'
    'jIzomHASABKAsyGS5sYW1iZGEuU25hcFN0YXJ0UmVzcG9uc2VSCXNuYXBzdGFydBInCgVzdGF0'
    'ZRiXybLvASABKA4yDS5sYW1iZGEuU3RhdGVSBXN0YXRlEiQKC3N0YXRlcmVhc29uGPPVrbMBIA'
    'EoCVILc3RhdGVyZWFzb24SRQoPc3RhdGVyZWFzb25jb2RlGMCpnpgBIAEoDjIXLmxhbWJkYS5T'
    'dGF0ZVJlYXNvbkNvZGVSD3N0YXRlcmVhc29uY29kZRI+Cg10ZW5hbmN5Y29uZmlnGPqr7WYgAS'
    'gLMhUubGFtYmRhLlRlbmFuY3lDb25maWdSDXRlbmFuY3ljb25maWcSGwoHdGltZW91dBip/OUW'
    'IAEoBVIHdGltZW91dBJGCg10cmFjaW5nY29uZmlnGKzEqQkgASgLMh0ubGFtYmRhLlRyYWNpbm'
    'dDb25maWdSZXNwb25zZVINdHJhY2luZ2NvbmZpZxIcCgd2ZXJzaW9uGLiqt+4BIAEoCVIHdmVy'
    'c2lvbhI6Cgl2cGNjb25maWcYh9f8XCABKAsyGS5sYW1iZGEuVnBjQ29uZmlnUmVzcG9uc2VSCX'
    'ZwY2NvbmZpZw==');

@$core.Deprecated('Use functionEventInvokeConfigDescriptor instead')
const FunctionEventInvokeConfig$json = {
  '1': 'FunctionEventInvokeConfig',
  '2': [
    {
      '1': 'destinationconfig',
      '3': 184834158,
      '4': 1,
      '5': 11,
      '6': '.lambda.DestinationConfig',
      '10': 'destinationconfig'
    },
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'maximumeventageinseconds',
      '3': 393041563,
      '4': 1,
      '5': 5,
      '10': 'maximumeventageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
  ],
};

/// Descriptor for `FunctionEventInvokeConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionEventInvokeConfigDescriptor = $convert.base64Decode(
    'ChlGdW5jdGlvbkV2ZW50SW52b2tlQ29uZmlnEkoKEWRlc3RpbmF0aW9uY29uZmlnGO6wkVggAS'
    'gLMhkubGFtYmRhLkRlc3RpbmF0aW9uQ29uZmlnUhFkZXN0aW5hdGlvbmNvbmZpZxIkCgtmdW5j'
    'dGlvbmFybhjxyY62ASABKAlSC2Z1bmN0aW9uYXJuEiYKDGxhc3Rtb2RpZmllZBinnPzOASABKA'
    'lSDGxhc3Rtb2RpZmllZBI+ChhtYXhpbXVtZXZlbnRhZ2VpbnNlY29uZHMYm621uwEgASgFUhht'
    'YXhpbXVtZXZlbnRhZ2VpbnNlY29uZHMSNQoUbWF4aW11bXJldHJ5YXR0ZW1wdHMYwKi5NSABKA'
    'VSFG1heGltdW1yZXRyeWF0dGVtcHRz');

@$core.Deprecated('Use functionScalingConfigDescriptor instead')
const FunctionScalingConfig$json = {
  '1': 'FunctionScalingConfig',
  '2': [
    {
      '1': 'maxexecutionenvironments',
      '3': 337598368,
      '4': 1,
      '5': 5,
      '10': 'maxexecutionenvironments'
    },
    {
      '1': 'minexecutionenvironments',
      '3': 123880642,
      '4': 1,
      '5': 5,
      '10': 'minexecutionenvironments'
    },
  ],
};

/// Descriptor for `FunctionScalingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionScalingConfigDescriptor = $convert.base64Decode(
    'ChVGdW5jdGlvblNjYWxpbmdDb25maWcSPgoYbWF4ZXhlY3V0aW9uZW52aXJvbm1lbnRzGKCv/a'
    'ABIAEoBVIYbWF4ZXhlY3V0aW9uZW52aXJvbm1lbnRzEj0KGG1pbmV4ZWN1dGlvbmVudmlyb25t'
    'ZW50cxjCiYk7IAEoBVIYbWluZXhlY3V0aW9uZW52aXJvbm1lbnRz');

@$core.Deprecated('Use functionUrlConfigDescriptor instead')
const FunctionUrlConfig$json = {
  '1': 'FunctionUrlConfig',
  '2': [
    {
      '1': 'authtype',
      '3': 477704248,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionUrlAuthType',
      '10': 'authtype'
    },
    {
      '1': 'cors',
      '3': 260753653,
      '4': 1,
      '5': 11,
      '6': '.lambda.Cors',
      '10': 'cors'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'functionurl', '3': 449381947, '4': 1, '5': 9, '10': 'functionurl'},
    {
      '1': 'invokemode',
      '3': 414956667,
      '4': 1,
      '5': 14,
      '6': '.lambda.InvokeMode',
      '10': 'invokemode'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `FunctionUrlConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionUrlConfigDescriptor = $convert.base64Decode(
    'ChFGdW5jdGlvblVybENvbmZpZxI7CghhdXRodHlwZRi44OTjASABKA4yGy5sYW1iZGEuRnVuY3'
    'Rpb25VcmxBdXRoVHlwZVIIYXV0aHR5cGUSIwoEY29ycxj1kat8IAEoCzIMLmxhbWJkYS5Db3Jz'
    'UgRjb3JzEiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3JlYXRpb250aW1lEiQKC2Z1bmN0aW'
    '9uYXJuGPHJjrYBIAEoCVILZnVuY3Rpb25hcm4SJAoLZnVuY3Rpb251cmwYu4yk1gEgASgJUgtm'
    'dW5jdGlvbnVybBI2CgppbnZva2Vtb2RlGPv47sUBIAEoDjISLmxhbWJkYS5JbnZva2VNb2RlUg'
    'ppbnZva2Vtb2RlEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRp'
    'bWU=');

@$core.Deprecated(
    'Use functionVersionsByCapacityProviderListItemDescriptor instead')
const FunctionVersionsByCapacityProviderListItem$json = {
  '1': 'FunctionVersionsByCapacityProviderListItem',
  '2': [
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.lambda.State',
      '10': 'state'
    },
  ],
};

/// Descriptor for `FunctionVersionsByCapacityProviderListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    functionVersionsByCapacityProviderListItemDescriptor =
    $convert.base64Decode(
        'CipGdW5jdGlvblZlcnNpb25zQnlDYXBhY2l0eVByb3ZpZGVyTGlzdEl0ZW0SJAoLZnVuY3Rpb2'
        '5hcm4Y8cmOtgEgASgJUgtmdW5jdGlvbmFybhInCgVzdGF0ZRiXybLvASABKA4yDS5sYW1iZGEu'
        'U3RhdGVSBXN0YXRl');

@$core.Deprecated(
    'Use functionVersionsPerCapacityProviderLimitExceededExceptionDescriptor instead')
const FunctionVersionsPerCapacityProviderLimitExceededException$json = {
  '1': 'FunctionVersionsPerCapacityProviderLimitExceededException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FunctionVersionsPerCapacityProviderLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    functionVersionsPerCapacityProviderLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CjlGdW5jdGlvblZlcnNpb25zUGVyQ2FwYWNpdHlQcm92aWRlckxpbWl0RXhjZWVkZWRFeGNlcH'
        'Rpb24SFgoEdHlwZRjuoNeKASABKAlSBHR5cGUSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2Fn'
        'ZQ==');

@$core.Deprecated('Use getAccountSettingsRequestDescriptor instead')
const GetAccountSettingsRequest$json = {
  '1': 'GetAccountSettingsRequest',
};

/// Descriptor for `GetAccountSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountSettingsRequestDescriptor =
    $convert.base64Decode('ChlHZXRBY2NvdW50U2V0dGluZ3NSZXF1ZXN0');

@$core.Deprecated('Use getAccountSettingsResponseDescriptor instead')
const GetAccountSettingsResponse$json = {
  '1': 'GetAccountSettingsResponse',
  '2': [
    {
      '1': 'accountlimit',
      '3': 433716352,
      '4': 1,
      '5': 11,
      '6': '.lambda.AccountLimit',
      '10': 'accountlimit'
    },
    {
      '1': 'accountusage',
      '3': 508300372,
      '4': 1,
      '5': 11,
      '6': '.lambda.AccountUsage',
      '10': 'accountusage'
    },
  ],
};

/// Descriptor for `GetAccountSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountSettingsResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRBY2NvdW50U2V0dGluZ3NSZXNwb25zZRI8CgxhY2NvdW50bGltaXQYgPnnzgEgASgLMh'
        'QubGFtYmRhLkFjY291bnRMaW1pdFIMYWNjb3VudGxpbWl0EjwKDGFjY291bnR1c2FnZRjUmLDy'
        'ASABKAsyFC5sYW1iZGEuQWNjb3VudFVzYWdlUgxhY2NvdW50dXNhZ2U=');

@$core.Deprecated('Use getAliasRequestDescriptor instead')
const GetAliasRequest$json = {
  '1': 'GetAliasRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAliasRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRBbGlhc1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb25uYW'
    '1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use getCapacityProviderRequestDescriptor instead')
const GetCapacityProviderRequest$json = {
  '1': 'GetCapacityProviderRequest',
  '2': [
    {
      '1': 'capacityprovidername',
      '3': 466230132,
      '4': 1,
      '5': 9,
      '10': 'capacityprovidername'
    },
  ],
};

/// Descriptor for `GetCapacityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCapacityProviderRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRDYXBhY2l0eVByb3ZpZGVyUmVxdWVzdBI2ChRjYXBhY2l0eXByb3ZpZGVybmFtZRj0tq'
        'jeASABKAlSFGNhcGFjaXR5cHJvdmlkZXJuYW1l');

@$core.Deprecated('Use getCapacityProviderResponseDescriptor instead')
const GetCapacityProviderResponse$json = {
  '1': 'GetCapacityProviderResponse',
  '2': [
    {
      '1': 'capacityprovider',
      '3': 448074961,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProvider',
      '10': 'capacityprovider'
    },
  ],
};

/// Descriptor for `GetCapacityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCapacityProviderResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRDYXBhY2l0eVByb3ZpZGVyUmVzcG9uc2USSAoQY2FwYWNpdHlwcm92aWRlchjRqdTVAS'
        'ABKAsyGC5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlclIQY2FwYWNpdHlwcm92aWRlcg==');

@$core.Deprecated('Use getCodeSigningConfigRequestDescriptor instead')
const GetCodeSigningConfigRequest$json = {
  '1': 'GetCodeSigningConfigRequest',
  '2': [
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
  ],
};

/// Descriptor for `GetCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCodeSigningConfigRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRDb2RlU2lnbmluZ0NvbmZpZ1JlcXVlc3QSNgoUY29kZXNpZ25pbmdjb25maWdhcm4Ywf'
        'z38AEgASgJUhRjb2Rlc2lnbmluZ2NvbmZpZ2Fybg==');

@$core.Deprecated('Use getCodeSigningConfigResponseDescriptor instead')
const GetCodeSigningConfigResponse$json = {
  '1': 'GetCodeSigningConfigResponse',
  '2': [
    {
      '1': 'codesigningconfig',
      '3': 130458490,
      '4': 1,
      '5': 11,
      '6': '.lambda.CodeSigningConfig',
      '10': 'codesigningconfig'
    },
  ],
};

/// Descriptor for `GetCodeSigningConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCodeSigningConfigResponseDescriptor =
    $convert.base64Decode(
        'ChxHZXRDb2RlU2lnbmluZ0NvbmZpZ1Jlc3BvbnNlEkoKEWNvZGVzaWduaW5nY29uZmlnGPrGmj'
        '4gASgLMhkubGFtYmRhLkNvZGVTaWduaW5nQ29uZmlnUhFjb2Rlc2lnbmluZ2NvbmZpZw==');

@$core.Deprecated('Use getDurableExecutionHistoryRequestDescriptor instead')
const GetDurableExecutionHistoryRequest$json = {
  '1': 'GetDurableExecutionHistoryRequest',
  '2': [
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {
      '1': 'includeexecutiondata',
      '3': 194359928,
      '4': 1,
      '5': 8,
      '10': 'includeexecutiondata'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'reverseorder', '3': 427445336, '4': 1, '5': 8, '10': 'reverseorder'},
  ],
};

/// Descriptor for `GetDurableExecutionHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDurableExecutionHistoryRequestDescriptor = $convert.base64Decode(
    'CiFHZXREdXJhYmxlRXhlY3V0aW9uSGlzdG9yeVJlcXVlc3QSNAoTZHVyYWJsZWV4ZWN1dGlvbm'
    'FybhiKzbeAASABKAlSE2R1cmFibGVleGVjdXRpb25hcm4SNQoUaW5jbHVkZWV4ZWN1dGlvbmRh'
    'dGEY+OTWXCABKAhSFGluY2x1ZGVleGVjdXRpb25kYXRhEhkKBm1hcmtlchi43c0qIAEoCVIGbW'
    'Fya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSJgoMcmV2ZXJzZW9yZGVyGNiY'
    '6csBIAEoCFIMcmV2ZXJzZW9yZGVy');

@$core.Deprecated('Use getDurableExecutionHistoryResponseDescriptor instead')
const GetDurableExecutionHistoryResponse$json = {
  '1': 'GetDurableExecutionHistoryResponse',
  '2': [
    {
      '1': 'events',
      '3': 3416229,
      '4': 3,
      '5': 11,
      '6': '.lambda.Event',
      '10': 'events'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `GetDurableExecutionHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDurableExecutionHistoryResponseDescriptor =
    $convert.base64Decode(
        'CiJHZXREdXJhYmxlRXhlY3V0aW9uSGlzdG9yeVJlc3BvbnNlEigKBmV2ZW50cxilwdABIAMoCz'
        'INLmxhbWJkYS5FdmVudFIGZXZlbnRzEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFy'
        'a2Vy');

@$core.Deprecated('Use getDurableExecutionRequestDescriptor instead')
const GetDurableExecutionRequest$json = {
  '1': 'GetDurableExecutionRequest',
  '2': [
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
  ],
};

/// Descriptor for `GetDurableExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDurableExecutionRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXREdXJhYmxlRXhlY3V0aW9uUmVxdWVzdBI0ChNkdXJhYmxlZXhlY3V0aW9uYXJuGIrNt4'
        'ABIAEoCVITZHVyYWJsZWV4ZWN1dGlvbmFybg==');

@$core.Deprecated('Use getDurableExecutionResponseDescriptor instead')
const GetDurableExecutionResponse$json = {
  '1': 'GetDurableExecutionResponse',
  '2': [
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {
      '1': 'durableexecutionname',
      '3': 251828526,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionname'
    },
    {'1': 'endtimestamp', '3': 340967267, '4': 1, '5': 9, '10': 'endtimestamp'},
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '10': 'error'
    },
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'inputpayload', '3': 228187532, '4': 1, '5': 9, '10': 'inputpayload'},
    {'1': 'result', '3': 273346629, '4': 1, '5': 9, '10': 'result'},
    {
      '1': 'starttimestamp',
      '3': 392954318,
      '4': 1,
      '5': 9,
      '10': 'starttimestamp'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.lambda.ExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'traceheader',
      '3': 303628800,
      '4': 1,
      '5': 11,
      '6': '.lambda.TraceHeader',
      '10': 'traceheader'
    },
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `GetDurableExecutionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDurableExecutionResponseDescriptor = $convert.base64Decode(
    'ChtHZXREdXJhYmxlRXhlY3V0aW9uUmVzcG9uc2USNAoTZHVyYWJsZWV4ZWN1dGlvbmFybhiKzb'
    'eAASABKAlSE2R1cmFibGVleGVjdXRpb25hcm4SNQoUZHVyYWJsZWV4ZWN1dGlvbm5hbWUYrrKK'
    'eCABKAlSFGR1cmFibGVleGVjdXRpb25uYW1lEiYKDGVuZHRpbWVzdGFtcBjj/sqiASABKAlSDG'
    'VuZHRpbWVzdGFtcBItCgVlcnJvchjyubacASABKAsyEy5sYW1iZGEuRXJyb3JPYmplY3RSBWVy'
    'cm9yEiQKC2Z1bmN0aW9uYXJuGPHJjrYBIAEoCVILZnVuY3Rpb25hcm4SJQoMaW5wdXRwYXlsb2'
    'FkGIy752wgASgJUgxpbnB1dHBheWxvYWQSGgoGcmVzdWx0GMXgq4IBIAEoCVIGcmVzdWx0EioK'
    'DnN0YXJ0dGltZXN0YW1wGM6DsLsBIAEoCVIOc3RhcnR0aW1lc3RhbXASMgoGc3RhdHVzGJDk+w'
    'IgASgOMhcubGFtYmRhLkV4ZWN1dGlvblN0YXR1c1IGc3RhdHVzEjkKC3RyYWNlaGVhZGVyGICE'
    '5JABIAEoCzITLmxhbWJkYS5UcmFjZUhlYWRlclILdHJhY2VoZWFkZXISHAoHdmVyc2lvbhi4qr'
    'fuASABKAlSB3ZlcnNpb24=');

@$core.Deprecated('Use getDurableExecutionStateRequestDescriptor instead')
const GetDurableExecutionStateRequest$json = {
  '1': 'GetDurableExecutionStateRequest',
  '2': [
    {
      '1': 'checkpointtoken',
      '3': 27028053,
      '4': 1,
      '5': 9,
      '10': 'checkpointtoken'
    },
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `GetDurableExecutionStateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDurableExecutionStateRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXREdXJhYmxlRXhlY3V0aW9uU3RhdGVSZXF1ZXN0EisKD2NoZWNrcG9pbnR0b2tlbhjV1P'
        'EMIAEoCVIPY2hlY2twb2ludHRva2VuEjQKE2R1cmFibGVleGVjdXRpb25hcm4Yis23gAEgASgJ'
        'UhNkdXJhYmxlZXhlY3V0aW9uYXJuEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG1heG'
        'l0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use getDurableExecutionStateResponseDescriptor instead')
const GetDurableExecutionStateResponse$json = {
  '1': 'GetDurableExecutionStateResponse',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {
      '1': 'operations',
      '3': 126776656,
      '4': 3,
      '5': 11,
      '6': '.lambda.Operation',
      '10': 'operations'
    },
  ],
};

/// Descriptor for `GetDurableExecutionStateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDurableExecutionStateResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXREdXJhYmxlRXhlY3V0aW9uU3RhdGVSZXNwb25zZRIiCgpuZXh0bWFya2VyGKOBrv0BIA'
        'EoCVIKbmV4dG1hcmtlchI0CgpvcGVyYXRpb25zGNDquTwgAygLMhEubGFtYmRhLk9wZXJhdGlv'
        'blIKb3BlcmF0aW9ucw==');

@$core.Deprecated('Use getEventSourceMappingRequestDescriptor instead')
const GetEventSourceMappingRequest$json = {
  '1': 'GetEventSourceMappingRequest',
  '2': [
    {'1': 'uuid', '3': 91981875, '4': 1, '5': 9, '10': 'uuid'},
  ],
};

/// Descriptor for `GetEventSourceMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEventSourceMappingRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRFdmVudFNvdXJjZU1hcHBpbmdSZXF1ZXN0EhUKBHV1aWQYs5DuKyABKAlSBHV1aWQ=');

@$core.Deprecated('Use getFunctionCodeSigningConfigRequestDescriptor instead')
const GetFunctionCodeSigningConfigRequest$json = {
  '1': 'GetFunctionCodeSigningConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `GetFunctionCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionCodeSigningConfigRequestDescriptor =
    $convert.base64Decode(
        'CiNHZXRGdW5jdGlvbkNvZGVTaWduaW5nQ29uZmlnUmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4'
        'i/3wEgASgJUgxmdW5jdGlvbm5hbWU=');

@$core.Deprecated('Use getFunctionCodeSigningConfigResponseDescriptor instead')
const GetFunctionCodeSigningConfigResponse$json = {
  '1': 'GetFunctionCodeSigningConfigResponse',
  '2': [
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `GetFunctionCodeSigningConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionCodeSigningConfigResponseDescriptor =
    $convert.base64Decode(
        'CiRHZXRGdW5jdGlvbkNvZGVTaWduaW5nQ29uZmlnUmVzcG9uc2USNgoUY29kZXNpZ25pbmdjb2'
        '5maWdhcm4Ywfz38AEgASgJUhRjb2Rlc2lnbmluZ2NvbmZpZ2FybhImCgxmdW5jdGlvbm5hbWUY'
        'o4i/3wEgASgJUgxmdW5jdGlvbm5hbWU=');

@$core.Deprecated('Use getFunctionConcurrencyRequestDescriptor instead')
const GetFunctionConcurrencyRequest$json = {
  '1': 'GetFunctionConcurrencyRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `GetFunctionConcurrencyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionConcurrencyRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRGdW5jdGlvbkNvbmN1cnJlbmN5UmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4i/3wEgAS'
        'gJUgxmdW5jdGlvbm5hbWU=');

@$core.Deprecated('Use getFunctionConcurrencyResponseDescriptor instead')
const GetFunctionConcurrencyResponse$json = {
  '1': 'GetFunctionConcurrencyResponse',
  '2': [
    {
      '1': 'reservedconcurrentexecutions',
      '3': 40149212,
      '4': 1,
      '5': 5,
      '10': 'reservedconcurrentexecutions'
    },
  ],
};

/// Descriptor for `GetFunctionConcurrencyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionConcurrencyResponseDescriptor =
    $convert.base64Decode(
        'Ch5HZXRGdW5jdGlvbkNvbmN1cnJlbmN5UmVzcG9uc2USRQoccmVzZXJ2ZWRjb25jdXJyZW50ZX'
        'hlY3V0aW9ucxjcwZITIAEoBVIccmVzZXJ2ZWRjb25jdXJyZW50ZXhlY3V0aW9ucw==');

@$core.Deprecated('Use getFunctionConfigurationRequestDescriptor instead')
const GetFunctionConfigurationRequest$json = {
  '1': 'GetFunctionConfigurationRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetFunctionConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionConfigurationRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRGdW5jdGlvbkNvbmZpZ3VyYXRpb25SZXF1ZXN0EiYKDGZ1bmN0aW9ubmFtZRijiL/fAS'
        'ABKAlSDGZ1bmN0aW9ubmFtZRIgCglxdWFsaWZpZXIY4LWR+wEgASgJUglxdWFsaWZpZXI=');

@$core.Deprecated('Use getFunctionEventInvokeConfigRequestDescriptor instead')
const GetFunctionEventInvokeConfigRequest$json = {
  '1': 'GetFunctionEventInvokeConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetFunctionEventInvokeConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionEventInvokeConfigRequestDescriptor =
    $convert.base64Decode(
        'CiNHZXRGdW5jdGlvbkV2ZW50SW52b2tlQ29uZmlnUmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4'
        'i/3wEgASgJUgxmdW5jdGlvbm5hbWUSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlmaWVy');

@$core.Deprecated('Use getFunctionRecursionConfigRequestDescriptor instead')
const GetFunctionRecursionConfigRequest$json = {
  '1': 'GetFunctionRecursionConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `GetFunctionRecursionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionRecursionConfigRequestDescriptor =
    $convert.base64Decode(
        'CiFHZXRGdW5jdGlvblJlY3Vyc2lvbkNvbmZpZ1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv9'
        '8BIAEoCVIMZnVuY3Rpb25uYW1l');

@$core.Deprecated('Use getFunctionRecursionConfigResponseDescriptor instead')
const GetFunctionRecursionConfigResponse$json = {
  '1': 'GetFunctionRecursionConfigResponse',
  '2': [
    {
      '1': 'recursiveloop',
      '3': 2821758,
      '4': 1,
      '5': 14,
      '6': '.lambda.RecursiveLoop',
      '10': 'recursiveloop'
    },
  ],
};

/// Descriptor for `GetFunctionRecursionConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionRecursionConfigResponseDescriptor =
    $convert.base64Decode(
        'CiJHZXRGdW5jdGlvblJlY3Vyc2lvbkNvbmZpZ1Jlc3BvbnNlEj4KDXJlY3Vyc2l2ZWxvb3AY/p'
        'ysASABKA4yFS5sYW1iZGEuUmVjdXJzaXZlTG9vcFINcmVjdXJzaXZlbG9vcA==');

@$core.Deprecated('Use getFunctionRequestDescriptor instead')
const GetFunctionRequest$json = {
  '1': 'GetFunctionRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionRequestDescriptor = $convert.base64Decode(
    'ChJHZXRGdW5jdGlvblJlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb2'
    '5uYW1lEiAKCXF1YWxpZmllchjgtZH7ASABKAlSCXF1YWxpZmllcg==');

@$core.Deprecated('Use getFunctionResponseDescriptor instead')
const GetFunctionResponse$json = {
  '1': 'GetFunctionResponse',
  '2': [
    {
      '1': 'code',
      '3': 425572629,
      '4': 1,
      '5': 11,
      '6': '.lambda.FunctionCodeLocation',
      '10': 'code'
    },
    {
      '1': 'concurrency',
      '3': 492620473,
      '4': 1,
      '5': 11,
      '6': '.lambda.Concurrency',
      '10': 'concurrency'
    },
    {
      '1': 'configuration',
      '3': 442426458,
      '4': 1,
      '5': 11,
      '6': '.lambda.FunctionConfiguration',
      '10': 'configuration'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.lambda.GetFunctionResponse.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'tagserror',
      '3': 86984231,
      '4': 1,
      '5': 11,
      '6': '.lambda.TagsError',
      '10': 'tagserror'
    },
  ],
  '3': [GetFunctionResponse_TagsEntry$json],
};

@$core.Deprecated('Use getFunctionResponseDescriptor instead')
const GetFunctionResponse_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetFunctionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionResponseDescriptor = $convert.base64Decode(
    'ChNHZXRGdW5jdGlvblJlc3BvbnNlEjQKBGNvZGUYlfL2ygEgASgLMhwubGFtYmRhLkZ1bmN0aW'
    '9uQ29kZUxvY2F0aW9uUgRjb2RlEjkKC2NvbmN1cnJlbmN5GLmV8+oBIAEoCzITLmxhbWJkYS5D'
    'b25jdXJyZW5jeVILY29uY3VycmVuY3kSRwoNY29uZmlndXJhdGlvbhjayPvSASABKAsyHS5sYW'
    '1iZGEuRnVuY3Rpb25Db25maWd1cmF0aW9uUg1jb25maWd1cmF0aW9uEj0KBHRhZ3MYwcH2tQEg'
    'AygLMiUubGFtYmRhLkdldEZ1bmN0aW9uUmVzcG9uc2UuVGFnc0VudHJ5UgR0YWdzEjIKCXRhZ3'
    'NlcnJvchinjL0pIAEoCzIRLmxhbWJkYS5UYWdzRXJyb3JSCXRhZ3NlcnJvcho3CglUYWdzRW50'
    'cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getFunctionScalingConfigRequestDescriptor instead')
const GetFunctionScalingConfigRequest$json = {
  '1': 'GetFunctionScalingConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetFunctionScalingConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionScalingConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRGdW5jdGlvblNjYWxpbmdDb25maWdSZXF1ZXN0EiYKDGZ1bmN0aW9ubmFtZRijiL/fAS'
        'ABKAlSDGZ1bmN0aW9ubmFtZRIgCglxdWFsaWZpZXIY4LWR+wEgASgJUglxdWFsaWZpZXI=');

@$core.Deprecated('Use getFunctionScalingConfigResponseDescriptor instead')
const GetFunctionScalingConfigResponse$json = {
  '1': 'GetFunctionScalingConfigResponse',
  '2': [
    {
      '1': 'appliedfunctionscalingconfig',
      '3': 440367614,
      '4': 1,
      '5': 11,
      '6': '.lambda.FunctionScalingConfig',
      '10': 'appliedfunctionscalingconfig'
    },
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {
      '1': 'requestedfunctionscalingconfig',
      '3': 190466507,
      '4': 1,
      '5': 11,
      '6': '.lambda.FunctionScalingConfig',
      '10': 'requestedfunctionscalingconfig'
    },
  ],
};

/// Descriptor for `GetFunctionScalingConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionScalingConfigResponseDescriptor = $convert.base64Decode(
    'CiBHZXRGdW5jdGlvblNjYWxpbmdDb25maWdSZXNwb25zZRJlChxhcHBsaWVkZnVuY3Rpb25zY2'
    'FsaW5nY29uZmlnGP7z/dEBIAEoCzIdLmxhbWJkYS5GdW5jdGlvblNjYWxpbmdDb25maWdSHGFw'
    'cGxpZWRmdW5jdGlvbnNjYWxpbmdjb25maWcSJAoLZnVuY3Rpb25hcm4Y8cmOtgEgASgJUgtmdW'
    '5jdGlvbmFybhJoCh5yZXF1ZXN0ZWRmdW5jdGlvbnNjYWxpbmdjb25maWcYy5PpWiABKAsyHS5s'
    'YW1iZGEuRnVuY3Rpb25TY2FsaW5nQ29uZmlnUh5yZXF1ZXN0ZWRmdW5jdGlvbnNjYWxpbmdjb2'
    '5maWc=');

@$core.Deprecated('Use getFunctionUrlConfigRequestDescriptor instead')
const GetFunctionUrlConfigRequest$json = {
  '1': 'GetFunctionUrlConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetFunctionUrlConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionUrlConfigRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRGdW5jdGlvblVybENvbmZpZ1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCV'
        'IMZnVuY3Rpb25uYW1lEiAKCXF1YWxpZmllchjgtZH7ASABKAlSCXF1YWxpZmllcg==');

@$core.Deprecated('Use getFunctionUrlConfigResponseDescriptor instead')
const GetFunctionUrlConfigResponse$json = {
  '1': 'GetFunctionUrlConfigResponse',
  '2': [
    {
      '1': 'authtype',
      '3': 477704248,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionUrlAuthType',
      '10': 'authtype'
    },
    {
      '1': 'cors',
      '3': 260753653,
      '4': 1,
      '5': 11,
      '6': '.lambda.Cors',
      '10': 'cors'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'functionurl', '3': 449381947, '4': 1, '5': 9, '10': 'functionurl'},
    {
      '1': 'invokemode',
      '3': 414956667,
      '4': 1,
      '5': 14,
      '6': '.lambda.InvokeMode',
      '10': 'invokemode'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `GetFunctionUrlConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFunctionUrlConfigResponseDescriptor = $convert.base64Decode(
    'ChxHZXRGdW5jdGlvblVybENvbmZpZ1Jlc3BvbnNlEjsKCGF1dGh0eXBlGLjg5OMBIAEoDjIbLm'
    'xhbWJkYS5GdW5jdGlvblVybEF1dGhUeXBlUghhdXRodHlwZRIjCgRjb3JzGPWRq3wgASgLMgwu'
    'bGFtYmRhLkNvcnNSBGNvcnMSJQoMY3JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlvbnRpbW'
    'USJAoLZnVuY3Rpb25hcm4Y8cmOtgEgASgJUgtmdW5jdGlvbmFybhIkCgtmdW5jdGlvbnVybBi7'
    'jKTWASABKAlSC2Z1bmN0aW9udXJsEjYKCmludm9rZW1vZGUY+/juxQEgASgOMhIubGFtYmRhLk'
    'ludm9rZU1vZGVSCmludm9rZW1vZGUSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFz'
    'dG1vZGlmaWVkdGltZQ==');

@$core.Deprecated('Use getLayerVersionByArnRequestDescriptor instead')
const GetLayerVersionByArnRequest$json = {
  '1': 'GetLayerVersionByArnRequest',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `GetLayerVersionByArnRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLayerVersionByArnRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRMYXllclZlcnNpb25CeUFyblJlcXVlc3QSFAoDYXJuGJ2b7b8BIAEoCVIDYXJu');

@$core.Deprecated('Use getLayerVersionPolicyRequestDescriptor instead')
const GetLayerVersionPolicyRequest$json = {
  '1': 'GetLayerVersionPolicyRequest',
  '2': [
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
    {
      '1': 'versionnumber',
      '3': 209346775,
      '4': 1,
      '5': 3,
      '10': 'versionnumber'
    },
  ],
};

/// Descriptor for `GetLayerVersionPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLayerVersionPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRMYXllclZlcnNpb25Qb2xpY3lSZXF1ZXN0EiAKCWxheWVybmFtZRiWqpjtASABKAlSCW'
        'xheWVybmFtZRInCg12ZXJzaW9ubnVtYmVyGNfB6WMgASgDUg12ZXJzaW9ubnVtYmVy');

@$core.Deprecated('Use getLayerVersionPolicyResponseDescriptor instead')
const GetLayerVersionPolicyResponse$json = {
  '1': 'GetLayerVersionPolicyResponse',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `GetLayerVersionPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLayerVersionPolicyResponseDescriptor =
    $convert.base64Decode(
        'Ch1HZXRMYXllclZlcnNpb25Qb2xpY3lSZXNwb25zZRIaCgZwb2xpY3kYoO/w4AEgASgJUgZwb2'
        'xpY3kSIgoKcmV2aXNpb25pZBiGo57uASABKAlSCnJldmlzaW9uaWQ=');

@$core.Deprecated('Use getLayerVersionRequestDescriptor instead')
const GetLayerVersionRequest$json = {
  '1': 'GetLayerVersionRequest',
  '2': [
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
    {
      '1': 'versionnumber',
      '3': 209346775,
      '4': 1,
      '5': 3,
      '10': 'versionnumber'
    },
  ],
};

/// Descriptor for `GetLayerVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLayerVersionRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRMYXllclZlcnNpb25SZXF1ZXN0EiAKCWxheWVybmFtZRiWqpjtASABKAlSCWxheWVybm'
        'FtZRInCg12ZXJzaW9ubnVtYmVyGNfB6WMgASgDUg12ZXJzaW9ubnVtYmVy');

@$core.Deprecated('Use getLayerVersionResponseDescriptor instead')
const GetLayerVersionResponse$json = {
  '1': 'GetLayerVersionResponse',
  '2': [
    {
      '1': 'compatiblearchitectures',
      '3': 20400170,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'compatiblearchitectures'
    },
    {
      '1': 'compatibleruntimes',
      '3': 300943187,
      '4': 3,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'compatibleruntimes'
    },
    {
      '1': 'content',
      '3': 23568227,
      '4': 1,
      '5': 11,
      '6': '.lambda.LayerVersionContentOutput',
      '10': 'content'
    },
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'layerarn', '3': 319652978, '4': 1, '5': 9, '10': 'layerarn'},
    {
      '1': 'layerversionarn',
      '3': 174478342,
      '4': 1,
      '5': 9,
      '10': 'layerversionarn'
    },
    {'1': 'licenseinfo', '3': 133903897, '4': 1, '5': 9, '10': 'licenseinfo'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `GetLayerVersionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getLayerVersionResponseDescriptor = $convert.base64Decode(
    'ChdHZXRMYXllclZlcnNpb25SZXNwb25zZRJRChdjb21wYXRpYmxlYXJjaGl0ZWN0dXJlcxiqkN'
    '0JIAMoDjIULmxhbWJkYS5BcmNoaXRlY3R1cmVSF2NvbXBhdGlibGVhcmNoaXRlY3R1cmVzEkMK'
    'EmNvbXBhdGlibGVydW50aW1lcxjTjsCPASADKA4yDy5sYW1iZGEuUnVudGltZVISY29tcGF0aW'
    'JsZXJ1bnRpbWVzEj4KB2NvbnRlbnQY476eCyABKAsyIS5sYW1iZGEuTGF5ZXJWZXJzaW9uQ29u'
    'dGVudE91dHB1dFIHY29udGVudBIkCgtjcmVhdGVkZGF0ZRiwsOfGASABKAlSC2NyZWF0ZWRkYX'
    'RlEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIeCghsYXllcmFybhjyiLaY'
    'ASABKAlSCGxheWVyYXJuEisKD2xheWVydmVyc2lvbmFybhiGqJlTIAEoCVIPbGF5ZXJ2ZXJzaW'
    '9uYXJuEiMKC2xpY2Vuc2VpbmZvGJns7D8gASgJUgtsaWNlbnNlaW5mbxIcCgd2ZXJzaW9uGLiq'
    't+4BIAEoA1IHdmVyc2lvbg==');

@$core.Deprecated('Use getPolicyRequestDescriptor instead')
const GetPolicyRequest$json = {
  '1': 'GetPolicyRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPolicyRequestDescriptor = $convert.base64Decode(
    'ChBHZXRQb2xpY3lSZXF1ZXN0EiYKDGZ1bmN0aW9ubmFtZRijiL/fASABKAlSDGZ1bmN0aW9ubm'
    'FtZRIgCglxdWFsaWZpZXIY4LWR+wEgASgJUglxdWFsaWZpZXI=');

@$core.Deprecated('Use getPolicyResponseDescriptor instead')
const GetPolicyResponse$json = {
  '1': 'GetPolicyResponse',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `GetPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPolicyResponseDescriptor = $convert.base64Decode(
    'ChFHZXRQb2xpY3lSZXNwb25zZRIaCgZwb2xpY3kYoO/w4AEgASgJUgZwb2xpY3kSIgoKcmV2aX'
    'Npb25pZBiGo57uASABKAlSCnJldmlzaW9uaWQ=');

@$core
    .Deprecated('Use getProvisionedConcurrencyConfigRequestDescriptor instead')
const GetProvisionedConcurrencyConfigRequest$json = {
  '1': 'GetProvisionedConcurrencyConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetProvisionedConcurrencyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getProvisionedConcurrencyConfigRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXRQcm92aXNpb25lZENvbmN1cnJlbmN5Q29uZmlnUmVxdWVzdBImCgxmdW5jdGlvbm5hbW'
        'UYo4i/3wEgASgJUgxmdW5jdGlvbm5hbWUSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlm'
        'aWVy');

@$core
    .Deprecated('Use getProvisionedConcurrencyConfigResponseDescriptor instead')
const GetProvisionedConcurrencyConfigResponse$json = {
  '1': 'GetProvisionedConcurrencyConfigResponse',
  '2': [
    {
      '1': 'allocatedprovisionedconcurrentexecutions',
      '3': 468402411,
      '4': 1,
      '5': 5,
      '10': 'allocatedprovisionedconcurrentexecutions'
    },
    {
      '1': 'availableprovisionedconcurrentexecutions',
      '3': 32168809,
      '4': 1,
      '5': 5,
      '10': 'availableprovisionedconcurrentexecutions'
    },
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'requestedprovisionedconcurrentexecutions',
      '3': 58431158,
      '4': 1,
      '5': 5,
      '10': 'requestedprovisionedconcurrentexecutions'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.lambda.ProvisionedConcurrencyStatusEnum',
      '10': 'status'
    },
    {'1': 'statusreason', '3': 139234172, '4': 1, '5': 9, '10': 'statusreason'},
  ],
};

/// Descriptor for `GetProvisionedConcurrencyConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getProvisionedConcurrencyConfigResponseDescriptor = $convert.base64Decode(
    'CidHZXRQcm92aXNpb25lZENvbmN1cnJlbmN5Q29uZmlnUmVzcG9uc2USXgooYWxsb2NhdGVkcH'
    'JvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxjrga3fASABKAVSKGFsbG9jYXRlZHByb3Zp'
    'c2lvbmVkY29uY3VycmVudGV4ZWN1dGlvbnMSXQooYXZhaWxhYmxlcHJvdmlzaW9uZWRjb25jdX'
    'JyZW50ZXhlY3V0aW9ucxjptqsPIAEoBVIoYXZhaWxhYmxlcHJvdmlzaW9uZWRjb25jdXJyZW50'
    'ZXhlY3V0aW9ucxImCgxsYXN0bW9kaWZpZWQYp5z8zgEgASgJUgxsYXN0bW9kaWZpZWQSXQoocm'
    'VxdWVzdGVkcHJvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxi2re4bIAEoBVIocmVxdWVz'
    'dGVkcHJvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxJDCgZzdGF0dXMYkOT7AiABKA4yKC'
    '5sYW1iZGEuUHJvdmlzaW9uZWRDb25jdXJyZW5jeVN0YXR1c0VudW1SBnN0YXR1cxIlCgxzdGF0'
    'dXNyZWFzb24Y/JayQiABKAlSDHN0YXR1c3JlYXNvbg==');

@$core.Deprecated('Use getRuntimeManagementConfigRequestDescriptor instead')
const GetRuntimeManagementConfigRequest$json = {
  '1': 'GetRuntimeManagementConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `GetRuntimeManagementConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRuntimeManagementConfigRequestDescriptor =
    $convert.base64Decode(
        'CiFHZXRSdW50aW1lTWFuYWdlbWVudENvbmZpZ1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv9'
        '8BIAEoCVIMZnVuY3Rpb25uYW1lEiAKCXF1YWxpZmllchjgtZH7ASABKAlSCXF1YWxpZmllcg==');

@$core.Deprecated('Use getRuntimeManagementConfigResponseDescriptor instead')
const GetRuntimeManagementConfigResponse$json = {
  '1': 'GetRuntimeManagementConfigResponse',
  '2': [
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {
      '1': 'runtimeversionarn',
      '3': 532500669,
      '4': 1,
      '5': 9,
      '10': 'runtimeversionarn'
    },
    {
      '1': 'updateruntimeon',
      '3': 285197810,
      '4': 1,
      '5': 14,
      '6': '.lambda.UpdateRuntimeOn',
      '10': 'updateruntimeon'
    },
  ],
};

/// Descriptor for `GetRuntimeManagementConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRuntimeManagementConfigResponseDescriptor =
    $convert.base64Decode(
        'CiJHZXRSdW50aW1lTWFuYWdlbWVudENvbmZpZ1Jlc3BvbnNlEiQKC2Z1bmN0aW9uYXJuGPHJjr'
        'YBIAEoCVILZnVuY3Rpb25hcm4SMAoRcnVudGltZXZlcnNpb25hcm4YvaH1/QEgASgJUhFydW50'
        'aW1ldmVyc2lvbmFybhJFCg91cGRhdGVydW50aW1lb24Y8ov/hwEgASgOMhcubGFtYmRhLlVwZG'
        'F0ZVJ1bnRpbWVPblIPdXBkYXRlcnVudGltZW9u');

@$core.Deprecated('Use imageConfigDescriptor instead')
const ImageConfig$json = {
  '1': 'ImageConfig',
  '2': [
    {'1': 'command', '3': 108826451, '4': 3, '5': 9, '10': 'command'},
    {'1': 'entrypoint', '3': 73718972, '4': 3, '5': 9, '10': 'entrypoint'},
    {
      '1': 'workingdirectory',
      '3': 478970252,
      '4': 1,
      '5': 9,
      '10': 'workingdirectory'
    },
  ],
};

/// Descriptor for `ImageConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List imageConfigDescriptor = $convert.base64Decode(
    'CgtJbWFnZUNvbmZpZxIbCgdjb21tYW5kGNOe8jMgAygJUgdjb21tYW5kEiEKCmVudHJ5cG9pbn'
    'QYvLmTIyADKAlSCmVudHJ5cG9pbnQSLgoQd29ya2luZ2RpcmVjdG9yeRiMg7LkASABKAlSEHdv'
    'cmtpbmdkaXJlY3Rvcnk=');

@$core.Deprecated('Use imageConfigErrorDescriptor instead')
const ImageConfigError$json = {
  '1': 'ImageConfigError',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ImageConfigError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List imageConfigErrorDescriptor = $convert.base64Decode(
    'ChBJbWFnZUNvbmZpZ0Vycm9yEh8KCWVycm9yY29kZRiZ1sMQIAEoCVIJZXJyb3Jjb2RlEhsKB2'
    '1lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use imageConfigResponseDescriptor instead')
const ImageConfigResponse$json = {
  '1': 'ImageConfigResponse',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ImageConfigError',
      '10': 'error'
    },
    {
      '1': 'imageconfig',
      '3': 281970485,
      '4': 1,
      '5': 11,
      '6': '.lambda.ImageConfig',
      '10': 'imageconfig'
    },
  ],
};

/// Descriptor for `ImageConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List imageConfigResponseDescriptor = $convert.base64Decode(
    'ChNJbWFnZUNvbmZpZ1Jlc3BvbnNlEjIKBWVycm9yGPK5tpwBIAEoCzIYLmxhbWJkYS5JbWFnZU'
    'NvbmZpZ0Vycm9yUgVlcnJvchI5CgtpbWFnZWNvbmZpZxi1jrqGASABKAsyEy5sYW1iZGEuSW1h'
    'Z2VDb25maWdSC2ltYWdlY29uZmln');

@$core.Deprecated('Use instanceRequirementsDescriptor instead')
const InstanceRequirements$json = {
  '1': 'InstanceRequirements',
  '2': [
    {
      '1': 'allowedinstancetypes',
      '3': 276234658,
      '4': 3,
      '5': 9,
      '10': 'allowedinstancetypes'
    },
    {
      '1': 'architectures',
      '3': 530490948,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'architectures'
    },
    {
      '1': 'excludedinstancetypes',
      '3': 249351054,
      '4': 3,
      '5': 9,
      '10': 'excludedinstancetypes'
    },
  ],
};

/// Descriptor for `InstanceRequirements`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceRequirementsDescriptor = $convert.base64Decode(
    'ChRJbnN0YW5jZVJlcXVpcmVtZW50cxI2ChRhbGxvd2VkaW5zdGFuY2V0eXBlcxiig9yDASADKA'
    'lSFGFsbG93ZWRpbnN0YW5jZXR5cGVzEj4KDWFyY2hpdGVjdHVyZXMYxMz6/AEgAygOMhQubGFt'
    'YmRhLkFyY2hpdGVjdHVyZVINYXJjaGl0ZWN0dXJlcxI3ChVleGNsdWRlZGluc3RhbmNldHlwZX'
    'MYjpfzdiADKAlSFWV4Y2x1ZGVkaW5zdGFuY2V0eXBlcw==');

@$core.Deprecated('Use invalidCodeSignatureExceptionDescriptor instead')
const InvalidCodeSignatureException$json = {
  '1': 'InvalidCodeSignatureException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `InvalidCodeSignatureException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidCodeSignatureExceptionDescriptor =
    $convert.base64Decode(
        'Ch1JbnZhbGlkQ29kZVNpZ25hdHVyZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdlEhYKBHR5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use invalidParameterValueExceptionDescriptor instead')
const InvalidParameterValueException$json = {
  '1': 'InvalidParameterValueException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterValueException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterValueExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkUGFyYW1ldGVyVmFsdWVFeGNlcHRpb24SFgoEdHlwZRjuoNeKASABKAlSBHR5cG'
        'USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidRequestContentExceptionDescriptor instead')
const InvalidRequestContentException$json = {
  '1': 'InvalidRequestContentException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidRequestContentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRequestContentExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkUmVxdWVzdENvbnRlbnRFeGNlcHRpb24SFgoEdHlwZRjuoNeKASABKAlSBHR5cG'
        'USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidRuntimeExceptionDescriptor instead')
const InvalidRuntimeException$json = {
  '1': 'InvalidRuntimeException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `InvalidRuntimeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRuntimeExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnZhbGlkUnVudGltZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdlEh'
        'YKBHR5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use invalidSecurityGroupIDExceptionDescriptor instead')
const InvalidSecurityGroupIDException$json = {
  '1': 'InvalidSecurityGroupIDException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `InvalidSecurityGroupIDException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSecurityGroupIDExceptionDescriptor =
    $convert.base64Decode(
        'Ch9JbnZhbGlkU2VjdXJpdHlHcm91cElERXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2USFgoEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use invalidSubnetIDExceptionDescriptor instead')
const InvalidSubnetIDException$json = {
  '1': 'InvalidSubnetIDException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `InvalidSubnetIDException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSubnetIDExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkU3VibmV0SURFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use invalidZipFileExceptionDescriptor instead')
const InvalidZipFileException$json = {
  '1': 'InvalidZipFileException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `InvalidZipFileException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidZipFileExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnZhbGlkWmlwRmlsZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdlEh'
        'YKBHR5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use invocationCompletedDetailsDescriptor instead')
const InvocationCompletedDetails$json = {
  '1': 'InvocationCompletedDetails',
  '2': [
    {'1': 'endtimestamp', '3': 340967267, '4': 1, '5': 9, '10': 'endtimestamp'},
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
    {'1': 'requestid', '3': 396701568, '4': 1, '5': 9, '10': 'requestid'},
    {
      '1': 'starttimestamp',
      '3': 392954318,
      '4': 1,
      '5': 9,
      '10': 'starttimestamp'
    },
  ],
};

/// Descriptor for `InvocationCompletedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invocationCompletedDetailsDescriptor = $convert.base64Decode(
    'ChpJbnZvY2F0aW9uQ29tcGxldGVkRGV0YWlscxImCgxlbmR0aW1lc3RhbXAY4/7KogEgASgJUg'
    'xlbmR0aW1lc3RhbXASLAoFZXJyb3IY8rm2nAEgASgLMhIubGFtYmRhLkV2ZW50RXJyb3JSBWVy'
    'cm9yEiAKCXJlcXVlc3RpZBiA35S9ASABKAlSCXJlcXVlc3RpZBIqCg5zdGFydHRpbWVzdGFtcB'
    'jOg7C7ASABKAlSDnN0YXJ0dGltZXN0YW1w');

@$core.Deprecated('Use invocationRequestDescriptor instead')
const InvocationRequest$json = {
  '1': 'InvocationRequest',
  '2': [
    {
      '1': 'clientcontext',
      '3': 354405962,
      '4': 1,
      '5': 9,
      '10': 'clientcontext'
    },
    {
      '1': 'durableexecutionname',
      '3': 251828526,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionname'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'invocationtype',
      '3': 492784524,
      '4': 1,
      '5': 14,
      '6': '.lambda.InvocationType',
      '10': 'invocationtype'
    },
    {
      '1': 'logtype',
      '3': 116703930,
      '4': 1,
      '5': 14,
      '6': '.lambda.LogType',
      '10': 'logtype'
    },
    {'1': 'payload', '3': 6526790, '4': 1, '5': 12, '8': {}, '10': 'payload'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
  ],
};

/// Descriptor for `InvocationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invocationRequestDescriptor = $convert.base64Decode(
    'ChFJbnZvY2F0aW9uUmVxdWVzdBIoCg1jbGllbnRjb250ZXh0GMqc/6gBIAEoCVINY2xpZW50Y2'
    '9udGV4dBI1ChRkdXJhYmxlZXhlY3V0aW9ubmFtZRiusop4IAEoCVIUZHVyYWJsZWV4ZWN1dGlv'
    'bm5hbWUSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb25uYW1lEkIKDmludm9jYX'
    'Rpb250eXBlGIyX/eoBIAEoDjIWLmxhbWJkYS5JbnZvY2F0aW9uVHlwZVIOaW52b2NhdGlvbnR5'
    'cGUSLAoHbG9ndHlwZRi6hdM3IAEoDjIPLmxhbWJkYS5Mb2dUeXBlUgdsb2d0eXBlEiEKB3BheW'
    'xvYWQYxq6OAyABKAxCBIi1GAFSB3BheWxvYWQSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVh'
    'bGlmaWVyEh0KCHRlbmFudGlkGIH472QgASgJUgh0ZW5hbnRpZA==');

@$core.Deprecated('Use invocationResponseDescriptor instead')
const InvocationResponse$json = {
  '1': 'InvocationResponse',
  '2': [
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {
      '1': 'executedversion',
      '3': 457672809,
      '4': 1,
      '5': 9,
      '10': 'executedversion'
    },
    {
      '1': 'functionerror',
      '3': 246179982,
      '4': 1,
      '5': 9,
      '10': 'functionerror'
    },
    {'1': 'logresult', '3': 367982401, '4': 1, '5': 9, '10': 'logresult'},
    {'1': 'payload', '3': 6526790, '4': 1, '5': 12, '8': {}, '10': 'payload'},
    {'1': 'statuscode', '3': 303830783, '4': 1, '5': 5, '10': 'statuscode'},
  ],
};

/// Descriptor for `InvocationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invocationResponseDescriptor = $convert.base64Decode(
    'ChJJbnZvY2F0aW9uUmVzcG9uc2USNAoTZHVyYWJsZWV4ZWN1dGlvbmFybhiKzbeAASABKAlSE2'
    'R1cmFibGVleGVjdXRpb25hcm4SLAoPZXhlY3V0ZWR2ZXJzaW9uGOmQntoBIAEoCVIPZXhlY3V0'
    'ZWR2ZXJzaW9uEicKDWZ1bmN0aW9uZXJyb3IYjtGxdSABKAlSDWZ1bmN0aW9uZXJyb3ISIAoJbG'
    '9ncmVzdWx0GMHuu68BIAEoCVIJbG9ncmVzdWx0EiEKB3BheWxvYWQYxq6OAyABKAxCBIi1GAFS'
    'B3BheWxvYWQSIgoKc3RhdHVzY29kZRj/rfCQASABKAVSCnN0YXR1c2NvZGU=');

@$core.Deprecated('Use invokeAsyncRequestDescriptor instead')
const InvokeAsyncRequest$json = {
  '1': 'InvokeAsyncRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'invokeargs',
      '3': 88134693,
      '4': 1,
      '5': 12,
      '8': {},
      '10': 'invokeargs'
    },
  ],
};

/// Descriptor for `InvokeAsyncRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invokeAsyncRequestDescriptor = $convert.base64Decode(
    'ChJJbnZva2VBc3luY1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb2'
    '5uYW1lEicKCmludm9rZWFyZ3MYpaiDKiABKAxCBIi1GAFSCmludm9rZWFyZ3M=');

@$core.Deprecated('Use invokeAsyncResponseDescriptor instead')
const InvokeAsyncResponse$json = {
  '1': 'InvokeAsyncResponse',
  '2': [
    {'1': 'status', '3': 6222352, '4': 1, '5': 5, '10': 'status'},
  ],
};

/// Descriptor for `InvokeAsyncResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invokeAsyncResponseDescriptor =
    $convert.base64Decode(
        'ChNJbnZva2VBc3luY1Jlc3BvbnNlEhkKBnN0YXR1cxiQ5PsCIAEoBVIGc3RhdHVz');

@$core.Deprecated('Use invokeResponseStreamUpdateDescriptor instead')
const InvokeResponseStreamUpdate$json = {
  '1': 'InvokeResponseStreamUpdate',
  '2': [
    {'1': 'payload', '3': 6526790, '4': 1, '5': 12, '10': 'payload'},
  ],
};

/// Descriptor for `InvokeResponseStreamUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invokeResponseStreamUpdateDescriptor =
    $convert.base64Decode(
        'ChpJbnZva2VSZXNwb25zZVN0cmVhbVVwZGF0ZRIbCgdwYXlsb2FkGMaujgMgASgMUgdwYXlsb2'
        'Fk');

@$core.Deprecated('Use invokeWithResponseStreamCompleteEventDescriptor instead')
const InvokeWithResponseStreamCompleteEvent$json = {
  '1': 'InvokeWithResponseStreamCompleteEvent',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errordetails', '3': 369268426, '4': 1, '5': 9, '10': 'errordetails'},
    {'1': 'logresult', '3': 367982401, '4': 1, '5': 9, '10': 'logresult'},
  ],
};

/// Descriptor for `InvokeWithResponseStreamCompleteEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invokeWithResponseStreamCompleteEventDescriptor =
    $convert.base64Decode(
        'CiVJbnZva2VXaXRoUmVzcG9uc2VTdHJlYW1Db21wbGV0ZUV2ZW50Eh8KCWVycm9yY29kZRiZ1s'
        'MQIAEoCVIJZXJyb3Jjb2RlEiYKDGVycm9yZGV0YWlscxjKrYqwASABKAlSDGVycm9yZGV0YWls'
        'cxIgCglsb2dyZXN1bHQYwe67rwEgASgJUglsb2dyZXN1bHQ=');

@$core.Deprecated('Use invokeWithResponseStreamRequestDescriptor instead')
const InvokeWithResponseStreamRequest$json = {
  '1': 'InvokeWithResponseStreamRequest',
  '2': [
    {
      '1': 'clientcontext',
      '3': 354405962,
      '4': 1,
      '5': 9,
      '10': 'clientcontext'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'invocationtype',
      '3': 492784524,
      '4': 1,
      '5': 14,
      '6': '.lambda.ResponseStreamingInvocationType',
      '10': 'invocationtype'
    },
    {
      '1': 'logtype',
      '3': 116703930,
      '4': 1,
      '5': 14,
      '6': '.lambda.LogType',
      '10': 'logtype'
    },
    {'1': 'payload', '3': 6526790, '4': 1, '5': 12, '8': {}, '10': 'payload'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
  ],
};

/// Descriptor for `InvokeWithResponseStreamRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invokeWithResponseStreamRequestDescriptor = $convert.base64Decode(
    'Ch9JbnZva2VXaXRoUmVzcG9uc2VTdHJlYW1SZXF1ZXN0EigKDWNsaWVudGNvbnRleHQYypz/qA'
    'EgASgJUg1jbGllbnRjb250ZXh0EiYKDGZ1bmN0aW9ubmFtZRijiL/fASABKAlSDGZ1bmN0aW9u'
    'bmFtZRJTCg5pbnZvY2F0aW9udHlwZRiMl/3qASABKA4yJy5sYW1iZGEuUmVzcG9uc2VTdHJlYW'
    '1pbmdJbnZvY2F0aW9uVHlwZVIOaW52b2NhdGlvbnR5cGUSLAoHbG9ndHlwZRi6hdM3IAEoDjIP'
    'LmxhbWJkYS5Mb2dUeXBlUgdsb2d0eXBlEiEKB3BheWxvYWQYxq6OAyABKAxCBIi1GAFSB3BheW'
    'xvYWQSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlmaWVyEh0KCHRlbmFudGlkGIH472Qg'
    'ASgJUgh0ZW5hbnRpZA==');

@$core.Deprecated('Use invokeWithResponseStreamResponseDescriptor instead')
const InvokeWithResponseStreamResponse$json = {
  '1': 'InvokeWithResponseStreamResponse',
  '2': [
    {
      '1': 'eventstream',
      '3': 26857888,
      '4': 1,
      '5': 11,
      '6': '.lambda.InvokeWithResponseStreamResponseEvent',
      '8': {},
      '10': 'eventstream'
    },
    {
      '1': 'executedversion',
      '3': 457672809,
      '4': 1,
      '5': 9,
      '10': 'executedversion'
    },
    {
      '1': 'responsestreamcontenttype',
      '3': 65265032,
      '4': 1,
      '5': 9,
      '10': 'responsestreamcontenttype'
    },
    {'1': 'statuscode', '3': 303830783, '4': 1, '5': 5, '10': 'statuscode'},
  ],
};

/// Descriptor for `InvokeWithResponseStreamResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invokeWithResponseStreamResponseDescriptor = $convert.base64Decode(
    'CiBJbnZva2VXaXRoUmVzcG9uc2VTdHJlYW1SZXNwb25zZRJYCgtldmVudHN0cmVhbRigo+cMIA'
    'EoCzItLmxhbWJkYS5JbnZva2VXaXRoUmVzcG9uc2VTdHJlYW1SZXNwb25zZUV2ZW50QgSItRgB'
    'UgtldmVudHN0cmVhbRIsCg9leGVjdXRlZHZlcnNpb24Y6ZCe2gEgASgJUg9leGVjdXRlZHZlcn'
    'Npb24SPwoZcmVzcG9uc2VzdHJlYW1jb250ZW50dHlwZRiIu48fIAEoCVIZcmVzcG9uc2VzdHJl'
    'YW1jb250ZW50dHlwZRIiCgpzdGF0dXNjb2RlGP+t8JABIAEoBVIKc3RhdHVzY29kZQ==');

@$core.Deprecated('Use invokeWithResponseStreamResponseEventDescriptor instead')
const InvokeWithResponseStreamResponseEvent$json = {
  '1': 'InvokeWithResponseStreamResponseEvent',
  '2': [
    {
      '1': 'invokecomplete',
      '3': 272618555,
      '4': 1,
      '5': 11,
      '6': '.lambda.InvokeWithResponseStreamCompleteEvent',
      '10': 'invokecomplete'
    },
    {
      '1': 'payloadchunk',
      '3': 230139107,
      '4': 1,
      '5': 11,
      '6': '.lambda.InvokeResponseStreamUpdate',
      '10': 'payloadchunk'
    },
  ],
};

/// Descriptor for `InvokeWithResponseStreamResponseEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invokeWithResponseStreamResponseEventDescriptor =
    $convert.base64Decode(
        'CiVJbnZva2VXaXRoUmVzcG9uc2VTdHJlYW1SZXNwb25zZUV2ZW50ElkKDmludm9rZWNvbXBsZX'
        'RlGLuo/4EBIAEoCzItLmxhbWJkYS5JbnZva2VXaXRoUmVzcG9uc2VTdHJlYW1Db21wbGV0ZUV2'
        'ZW50Ug5pbnZva2Vjb21wbGV0ZRJJCgxwYXlsb2FkY2h1bmsY48nebSABKAsyIi5sYW1iZGEuSW'
        '52b2tlUmVzcG9uc2VTdHJlYW1VcGRhdGVSDHBheWxvYWRjaHVuaw==');

@$core.Deprecated('Use kMSAccessDeniedExceptionDescriptor instead')
const KMSAccessDeniedException$json = {
  '1': 'KMSAccessDeniedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `KMSAccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSAccessDeniedExceptionDescriptor =
    $convert.base64Decode(
        'ChhLTVNBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use kMSDisabledExceptionDescriptor instead')
const KMSDisabledException$json = {
  '1': 'KMSDisabledException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `KMSDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSDisabledExceptionDescriptor = $convert.base64Decode(
    'ChRLTVNEaXNhYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdlEhYKBH'
    'R5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use kMSInvalidStateExceptionDescriptor instead')
const KMSInvalidStateException$json = {
  '1': 'KMSInvalidStateException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `KMSInvalidStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSInvalidStateExceptionDescriptor =
    $convert.base64Decode(
        'ChhLTVNJbnZhbGlkU3RhdGVFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use kMSNotFoundExceptionDescriptor instead')
const KMSNotFoundException$json = {
  '1': 'KMSNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `KMSNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSNotFoundExceptionDescriptor = $convert.base64Decode(
    'ChRLTVNOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdlEhYKBH'
    'R5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use kafkaSchemaRegistryAccessConfigDescriptor instead')
const KafkaSchemaRegistryAccessConfig$json = {
  '1': 'KafkaSchemaRegistryAccessConfig',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.lambda.KafkaSchemaRegistryAuthType',
      '10': 'type'
    },
    {'1': 'uri', '3': 443116318, '4': 1, '5': 9, '10': 'uri'},
  ],
};

/// Descriptor for `KafkaSchemaRegistryAccessConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kafkaSchemaRegistryAccessConfigDescriptor =
    $convert.base64Decode(
        'Ch9LYWZrYVNjaGVtYVJlZ2lzdHJ5QWNjZXNzQ29uZmlnEjsKBHR5cGUY7qDXigEgASgOMiMubG'
        'FtYmRhLkthZmthU2NoZW1hUmVnaXN0cnlBdXRoVHlwZVIEdHlwZRIUCgN1cmkYntal0wEgASgJ'
        'UgN1cmk=');

@$core.Deprecated('Use kafkaSchemaRegistryConfigDescriptor instead')
const KafkaSchemaRegistryConfig$json = {
  '1': 'KafkaSchemaRegistryConfig',
  '2': [
    {
      '1': 'accessconfigs',
      '3': 217772895,
      '4': 3,
      '5': 11,
      '6': '.lambda.KafkaSchemaRegistryAccessConfig',
      '10': 'accessconfigs'
    },
    {
      '1': 'eventrecordformat',
      '3': 212276390,
      '4': 1,
      '5': 14,
      '6': '.lambda.SchemaRegistryEventRecordFormat',
      '10': 'eventrecordformat'
    },
    {
      '1': 'schemaregistryuri',
      '3': 227288142,
      '4': 1,
      '5': 9,
      '10': 'schemaregistryuri'
    },
    {
      '1': 'schemavalidationconfigs',
      '3': 355660573,
      '4': 3,
      '5': 11,
      '6': '.lambda.KafkaSchemaValidationConfig',
      '10': 'schemavalidationconfigs'
    },
  ],
};

/// Descriptor for `KafkaSchemaRegistryConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kafkaSchemaRegistryConfigDescriptor = $convert.base64Decode(
    'ChlLYWZrYVNjaGVtYVJlZ2lzdHJ5Q29uZmlnElAKDWFjY2Vzc2NvbmZpZ3MY3+brZyADKAsyJy'
    '5sYW1iZGEuS2Fma2FTY2hlbWFSZWdpc3RyeUFjY2Vzc0NvbmZpZ1INYWNjZXNzY29uZmlncxJY'
    'ChFldmVudHJlY29yZGZvcm1hdBimqZxlIAEoDjInLmxhbWJkYS5TY2hlbWFSZWdpc3RyeUV2ZW'
    '50UmVjb3JkRm9ybWF0UhFldmVudHJlY29yZGZvcm1hdBIvChFzY2hlbWFyZWdpc3RyeXVyaRjO'
    'yLBsIAEoCVIRc2NoZW1hcmVnaXN0cnl1cmkSYQoXc2NoZW1hdmFsaWRhdGlvbmNvbmZpZ3MYne'
    'bLqQEgAygLMiMubGFtYmRhLkthZmthU2NoZW1hVmFsaWRhdGlvbkNvbmZpZ1IXc2NoZW1hdmFs'
    'aWRhdGlvbmNvbmZpZ3M=');

@$core.Deprecated('Use kafkaSchemaValidationConfigDescriptor instead')
const KafkaSchemaValidationConfig$json = {
  '1': 'KafkaSchemaValidationConfig',
  '2': [
    {
      '1': 'attribute',
      '3': 49977488,
      '4': 1,
      '5': 14,
      '6': '.lambda.KafkaSchemaValidationAttribute',
      '10': 'attribute'
    },
  ],
};

/// Descriptor for `KafkaSchemaValidationConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kafkaSchemaValidationConfigDescriptor =
    $convert.base64Decode(
        'ChtLYWZrYVNjaGVtYVZhbGlkYXRpb25Db25maWcSRwoJYXR0cmlidXRlGJCx6hcgASgOMiYubG'
        'FtYmRhLkthZmthU2NoZW1hVmFsaWRhdGlvbkF0dHJpYnV0ZVIJYXR0cmlidXRl');

@$core.Deprecated(
    'Use lambdaManagedInstancesCapacityProviderConfigDescriptor instead')
const LambdaManagedInstancesCapacityProviderConfig$json = {
  '1': 'LambdaManagedInstancesCapacityProviderConfig',
  '2': [
    {
      '1': 'capacityproviderarn',
      '3': 109580344,
      '4': 1,
      '5': 9,
      '10': 'capacityproviderarn'
    },
    {
      '1': 'executionenvironmentmemorygibpervcpu',
      '3': 227590513,
      '4': 1,
      '5': 1,
      '10': 'executionenvironmentmemorygibpervcpu'
    },
    {
      '1': 'perexecutionenvironmentmaxconcurrency',
      '3': 120967319,
      '4': 1,
      '5': 5,
      '10': 'perexecutionenvironmentmaxconcurrency'
    },
  ],
};

/// Descriptor for `LambdaManagedInstancesCapacityProviderConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    lambdaManagedInstancesCapacityProviderConfigDescriptor =
    $convert.base64Decode(
        'CixMYW1iZGFNYW5hZ2VkSW5zdGFuY2VzQ2FwYWNpdHlQcm92aWRlckNvbmZpZxIzChNjYXBhY2'
        'l0eXByb3ZpZGVyYXJuGLigoDQgASgJUhNjYXBhY2l0eXByb3ZpZGVyYXJuElUKJGV4ZWN1dGlv'
        'bmVudmlyb25tZW50bWVtb3J5Z2licGVydmNwdRjxgsNsIAEoAVIkZXhlY3V0aW9uZW52aXJvbm'
        '1lbnRtZW1vcnlnaWJwZXJ2Y3B1ElcKJXBlcmV4ZWN1dGlvbmVudmlyb25tZW50bWF4Y29uY3Vy'
        'cmVuY3kYl6HXOSABKAVSJXBlcmV4ZWN1dGlvbmVudmlyb25tZW50bWF4Y29uY3VycmVuY3k=');

@$core.Deprecated('Use layerDescriptor instead')
const Layer$json = {
  '1': 'Layer',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'codesize', '3': 74450158, '4': 1, '5': 3, '10': 'codesize'},
    {
      '1': 'signingjobarn',
      '3': 343397691,
      '4': 1,
      '5': 9,
      '10': 'signingjobarn'
    },
    {
      '1': 'signingprofileversionarn',
      '3': 432885567,
      '4': 1,
      '5': 9,
      '10': 'signingprofileversionarn'
    },
  ],
};

/// Descriptor for `Layer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List layerDescriptor = $convert.base64Decode(
    'CgVMYXllchIUCgNhcm4YnZvtvwEgASgJUgNhcm4SHQoIY29kZXNpemUY7onAIyABKANSCGNvZG'
    'VzaXplEigKDXNpZ25pbmdqb2Jhcm4Yu6rfowEgASgJUg1zaWduaW5nam9iYXJuEj4KGHNpZ25p'
    'bmdwcm9maWxldmVyc2lvbmFybhi/nrXOASABKAlSGHNpZ25pbmdwcm9maWxldmVyc2lvbmFybg'
    '==');

@$core.Deprecated('Use layerVersionContentInputDescriptor instead')
const LayerVersionContentInput$json = {
  '1': 'LayerVersionContentInput',
  '2': [
    {'1': 's3bucket', '3': 114031434, '4': 1, '5': 9, '10': 's3bucket'},
    {'1': 's3key', '3': 490298907, '4': 1, '5': 9, '10': 's3key'},
    {
      '1': 's3objectversion',
      '3': 194809669,
      '4': 1,
      '5': 9,
      '10': 's3objectversion'
    },
    {'1': 'zipfile', '3': 2519299, '4': 1, '5': 12, '10': 'zipfile'},
  ],
};

/// Descriptor for `LayerVersionContentInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List layerVersionContentInputDescriptor = $convert.base64Decode(
    'ChhMYXllclZlcnNpb25Db250ZW50SW5wdXQSHQoIczNidWNrZXQYyvavNiABKAlSCHMzYnVja2'
    'V0EhgKBXMza2V5GJu85ekBIAEoCVIFczNrZXkSKwoPczNvYmplY3R2ZXJzaW9uGMWe8lwgASgJ'
    'Ug9zM29iamVjdHZlcnNpb24SGwoHemlwZmlsZRiD4pkBIAEoDFIHemlwZmlsZQ==');

@$core.Deprecated('Use layerVersionContentOutputDescriptor instead')
const LayerVersionContentOutput$json = {
  '1': 'LayerVersionContentOutput',
  '2': [
    {'1': 'codesha256', '3': 46450860, '4': 1, '5': 9, '10': 'codesha256'},
    {'1': 'codesize', '3': 74450158, '4': 1, '5': 3, '10': 'codesize'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'signingjobarn',
      '3': 343397691,
      '4': 1,
      '5': 9,
      '10': 'signingjobarn'
    },
    {
      '1': 'signingprofileversionarn',
      '3': 432885567,
      '4': 1,
      '5': 9,
      '10': 'signingprofileversionarn'
    },
  ],
};

/// Descriptor for `LayerVersionContentOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List layerVersionContentOutputDescriptor = $convert.base64Decode(
    'ChlMYXllclZlcnNpb25Db250ZW50T3V0cHV0EiEKCmNvZGVzaGEyNTYYrJGTFiABKAlSCmNvZG'
    'VzaGEyNTYSHQoIY29kZXNpemUY7onAIyABKANSCGNvZGVzaXplEh4KCGxvY2F0aW9uGMebgt4B'
    'IAEoCVIIbG9jYXRpb24SKAoNc2lnbmluZ2pvYmFybhi7qt+jASABKAlSDXNpZ25pbmdqb2Jhcm'
    '4SPgoYc2lnbmluZ3Byb2ZpbGV2ZXJzaW9uYXJuGL+etc4BIAEoCVIYc2lnbmluZ3Byb2ZpbGV2'
    'ZXJzaW9uYXJu');

@$core.Deprecated('Use layerVersionsListItemDescriptor instead')
const LayerVersionsListItem$json = {
  '1': 'LayerVersionsListItem',
  '2': [
    {
      '1': 'compatiblearchitectures',
      '3': 20400170,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'compatiblearchitectures'
    },
    {
      '1': 'compatibleruntimes',
      '3': 300943187,
      '4': 3,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'compatibleruntimes'
    },
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'layerversionarn',
      '3': 174478342,
      '4': 1,
      '5': 9,
      '10': 'layerversionarn'
    },
    {'1': 'licenseinfo', '3': 133903897, '4': 1, '5': 9, '10': 'licenseinfo'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `LayerVersionsListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List layerVersionsListItemDescriptor = $convert.base64Decode(
    'ChVMYXllclZlcnNpb25zTGlzdEl0ZW0SUQoXY29tcGF0aWJsZWFyY2hpdGVjdHVyZXMYqpDdCS'
    'ADKA4yFC5sYW1iZGEuQXJjaGl0ZWN0dXJlUhdjb21wYXRpYmxlYXJjaGl0ZWN0dXJlcxJDChJj'
    'b21wYXRpYmxlcnVudGltZXMY047AjwEgAygOMg8ubGFtYmRhLlJ1bnRpbWVSEmNvbXBhdGlibG'
    'VydW50aW1lcxIkCgtjcmVhdGVkZGF0ZRiwsOfGASABKAlSC2NyZWF0ZWRkYXRlEiMKC2Rlc2Ny'
    'aXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIrCg9sYXllcnZlcnNpb25hcm4YhqiZUyABKA'
    'lSD2xheWVydmVyc2lvbmFybhIjCgtsaWNlbnNlaW5mbxiZ7Ow/IAEoCVILbGljZW5zZWluZm8S'
    'HAoHdmVyc2lvbhi4qrfuASABKANSB3ZlcnNpb24=');

@$core.Deprecated('Use layersListItemDescriptor instead')
const LayersListItem$json = {
  '1': 'LayersListItem',
  '2': [
    {
      '1': 'latestmatchingversion',
      '3': 389365332,
      '4': 1,
      '5': 11,
      '6': '.lambda.LayerVersionsListItem',
      '10': 'latestmatchingversion'
    },
    {'1': 'layerarn', '3': 319652978, '4': 1, '5': 9, '10': 'layerarn'},
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
  ],
};

/// Descriptor for `LayersListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List layersListItemDescriptor = $convert.base64Decode(
    'Cg5MYXllcnNMaXN0SXRlbRJXChVsYXRlc3RtYXRjaGluZ3ZlcnNpb24Y1PzUuQEgASgLMh0ubG'
    'FtYmRhLkxheWVyVmVyc2lvbnNMaXN0SXRlbVIVbGF0ZXN0bWF0Y2hpbmd2ZXJzaW9uEh4KCGxh'
    'eWVyYXJuGPKItpgBIAEoCVIIbGF5ZXJhcm4SIAoJbGF5ZXJuYW1lGJaqmO0BIAEoCVIJbGF5ZX'
    'JuYW1l');

@$core.Deprecated('Use listAliasesRequestDescriptor instead')
const ListAliasesRequest$json = {
  '1': 'ListAliasesRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'functionversion',
      '3': 365780244,
      '4': 1,
      '5': 9,
      '10': 'functionversion'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListAliasesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAliasesRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0QWxpYXNlc1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb2'
    '5uYW1lEiwKD2Z1bmN0aW9udmVyc2lvbhiUurWuASABKAlSD2Z1bmN0aW9udmVyc2lvbhIZCgZt'
    'YXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated('Use listAliasesResponseDescriptor instead')
const ListAliasesResponse$json = {
  '1': 'ListAliasesResponse',
  '2': [
    {
      '1': 'aliases',
      '3': 476693696,
      '4': 3,
      '5': 11,
      '6': '.lambda.AliasConfiguration',
      '10': 'aliases'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListAliasesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAliasesResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0QWxpYXNlc1Jlc3BvbnNlEjgKB2FsaWFzZXMYwImn4wEgAygLMhoubGFtYmRhLkFsaW'
    'FzQ29uZmlndXJhdGlvblIHYWxpYXNlcxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1h'
    'cmtlcg==');

@$core.Deprecated('Use listCapacityProvidersRequestDescriptor instead')
const ListCapacityProvidersRequest$json = {
  '1': 'ListCapacityProvidersRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.lambda.CapacityProviderState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `ListCapacityProvidersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCapacityProvidersRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0Q2FwYWNpdHlQcm92aWRlcnNSZXF1ZXN0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2'
        'VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSNwoFc3RhdGUYl8my7wEgASgOMh0u'
        'bGFtYmRhLkNhcGFjaXR5UHJvdmlkZXJTdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use listCapacityProvidersResponseDescriptor instead')
const ListCapacityProvidersResponse$json = {
  '1': 'ListCapacityProvidersResponse',
  '2': [
    {
      '1': 'capacityproviders',
      '3': 235895962,
      '4': 3,
      '5': 11,
      '6': '.lambda.CapacityProvider',
      '10': 'capacityproviders'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListCapacityProvidersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCapacityProvidersResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q2FwYWNpdHlQcm92aWRlcnNSZXNwb25zZRJJChFjYXBhY2l0eXByb3ZpZGVycxia+b'
        '1wIAMoCzIYLmxhbWJkYS5DYXBhY2l0eVByb3ZpZGVyUhFjYXBhY2l0eXByb3ZpZGVycxIiCgpu'
        'ZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use listCodeSigningConfigsRequestDescriptor instead')
const ListCodeSigningConfigsRequest$json = {
  '1': 'ListCodeSigningConfigsRequest',
  '2': [
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListCodeSigningConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCodeSigningConfigsRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q29kZVNpZ25pbmdDb25maWdzUmVxdWVzdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcm'
        'tlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated('Use listCodeSigningConfigsResponseDescriptor instead')
const ListCodeSigningConfigsResponse$json = {
  '1': 'ListCodeSigningConfigsResponse',
  '2': [
    {
      '1': 'codesigningconfigs',
      '3': 129149119,
      '4': 3,
      '5': 11,
      '6': '.lambda.CodeSigningConfig',
      '10': 'codesigningconfigs'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListCodeSigningConfigsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCodeSigningConfigsResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0Q29kZVNpZ25pbmdDb25maWdzUmVzcG9uc2USTAoSY29kZXNpZ25pbmdjb25maWdzGL'
        '/Ryj0gAygLMhkubGFtYmRhLkNvZGVTaWduaW5nQ29uZmlnUhJjb2Rlc2lnbmluZ2NvbmZpZ3MS'
        'IgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leHRtYXJrZXI=');

@$core
    .Deprecated('Use listDurableExecutionsByFunctionRequestDescriptor instead')
const ListDurableExecutionsByFunctionRequest$json = {
  '1': 'ListDurableExecutionsByFunctionRequest',
  '2': [
    {
      '1': 'durableexecutionname',
      '3': 251828526,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionname'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
    {'1': 'reverseorder', '3': 427445336, '4': 1, '5': 8, '10': 'reverseorder'},
    {'1': 'startedafter', '3': 284744661, '4': 1, '5': 9, '10': 'startedafter'},
    {
      '1': 'startedbefore',
      '3': 145036858,
      '4': 1,
      '5': 9,
      '10': 'startedbefore'
    },
    {
      '1': 'statuses',
      '3': 374056024,
      '4': 3,
      '5': 14,
      '6': '.lambda.ExecutionStatus',
      '10': 'statuses'
    },
  ],
};

/// Descriptor for `ListDurableExecutionsByFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDurableExecutionsByFunctionRequestDescriptor = $convert.base64Decode(
    'CiZMaXN0RHVyYWJsZUV4ZWN1dGlvbnNCeUZ1bmN0aW9uUmVxdWVzdBI1ChRkdXJhYmxlZXhlY3'
    'V0aW9ubmFtZRiusop4IAEoCVIUZHVyYWJsZWV4ZWN1dGlvbm5hbWUSJgoMZnVuY3Rpb25uYW1l'
    'GKOIv98BIAEoCVIMZnVuY3Rpb25uYW1lEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEh4KCG'
    '1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXMSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVh'
    'bGlmaWVyEiYKDHJldmVyc2VvcmRlchjYmOnLASABKAhSDHJldmVyc2VvcmRlchImCgxzdGFydG'
    'VkYWZ0ZXIY1bfjhwEgASgJUgxzdGFydGVkYWZ0ZXISJwoNc3RhcnRlZGJlZm9yZRi6rJRFIAEo'
    'CVINc3RhcnRlZGJlZm9yZRI3CghzdGF0dXNlcxjYyK6yASADKA4yFy5sYW1iZGEuRXhlY3V0aW'
    '9uU3RhdHVzUghzdGF0dXNlcw==');

@$core
    .Deprecated('Use listDurableExecutionsByFunctionResponseDescriptor instead')
const ListDurableExecutionsByFunctionResponse$json = {
  '1': 'ListDurableExecutionsByFunctionResponse',
  '2': [
    {
      '1': 'durableexecutions',
      '3': 65150760,
      '4': 3,
      '5': 11,
      '6': '.lambda.Execution',
      '10': 'durableexecutions'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListDurableExecutionsByFunctionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDurableExecutionsByFunctionResponseDescriptor =
    $convert.base64Decode(
        'CidMaXN0RHVyYWJsZUV4ZWN1dGlvbnNCeUZ1bmN0aW9uUmVzcG9uc2USQgoRZHVyYWJsZWV4ZW'
        'N1dGlvbnMYqL6IHyADKAsyES5sYW1iZGEuRXhlY3V0aW9uUhFkdXJhYmxlZXhlY3V0aW9ucxIi'
        'CgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use listEventSourceMappingsRequestDescriptor instead')
const ListEventSourceMappingsRequest$json = {
  '1': 'ListEventSourceMappingsRequest',
  '2': [
    {
      '1': 'eventsourcearn',
      '3': 306357574,
      '4': 1,
      '5': 9,
      '10': 'eventsourcearn'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListEventSourceMappingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventSourceMappingsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0RXZlbnRTb3VyY2VNYXBwaW5nc1JlcXVlc3QSKgoOZXZlbnRzb3VyY2Vhcm4YxsqKkg'
        'EgASgJUg5ldmVudHNvdXJjZWFybhImCgxmdW5jdGlvbm5hbWUYo4i/3wEgASgJUgxmdW5jdGlv'
        'bm5hbWUSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMYlNba8QEgASgFUg'
        'htYXhpdGVtcw==');

@$core.Deprecated('Use listEventSourceMappingsResponseDescriptor instead')
const ListEventSourceMappingsResponse$json = {
  '1': 'ListEventSourceMappingsResponse',
  '2': [
    {
      '1': 'eventsourcemappings',
      '3': 349608254,
      '4': 3,
      '5': 11,
      '6': '.lambda.EventSourceMappingConfiguration',
      '10': 'eventsourcemappings'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListEventSourceMappingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEventSourceMappingsResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0RXZlbnRTb3VyY2VNYXBwaW5nc1Jlc3BvbnNlEl0KE2V2ZW50c291cmNlbWFwcGluZ3'
        'MYvrLapgEgAygLMicubGFtYmRhLkV2ZW50U291cmNlTWFwcGluZ0NvbmZpZ3VyYXRpb25SE2V2'
        'ZW50c291cmNlbWFwcGluZ3MSIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leHRtYXJrZXI=');

@$core.Deprecated('Use listFunctionEventInvokeConfigsRequestDescriptor instead')
const ListFunctionEventInvokeConfigsRequest$json = {
  '1': 'ListFunctionEventInvokeConfigsRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListFunctionEventInvokeConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionEventInvokeConfigsRequestDescriptor =
    $convert.base64Decode(
        'CiVMaXN0RnVuY3Rpb25FdmVudEludm9rZUNvbmZpZ3NSZXF1ZXN0EiYKDGZ1bmN0aW9ubmFtZR'
        'ijiL/fASABKAlSDGZ1bmN0aW9ubmFtZRIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCght'
        'YXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core
    .Deprecated('Use listFunctionEventInvokeConfigsResponseDescriptor instead')
const ListFunctionEventInvokeConfigsResponse$json = {
  '1': 'ListFunctionEventInvokeConfigsResponse',
  '2': [
    {
      '1': 'functioneventinvokeconfigs',
      '3': 154516615,
      '4': 3,
      '5': 11,
      '6': '.lambda.FunctionEventInvokeConfig',
      '10': 'functioneventinvokeconfigs'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListFunctionEventInvokeConfigsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionEventInvokeConfigsResponseDescriptor =
    $convert.base64Decode(
        'CiZMaXN0RnVuY3Rpb25FdmVudEludm9rZUNvbmZpZ3NSZXNwb25zZRJkChpmdW5jdGlvbmV2ZW'
        '50aW52b2tlY29uZmlncxiH+dZJIAMoCzIhLmxhbWJkYS5GdW5jdGlvbkV2ZW50SW52b2tlQ29u'
        'ZmlnUhpmdW5jdGlvbmV2ZW50aW52b2tlY29uZmlncxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCV'
        'IKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use listFunctionUrlConfigsRequestDescriptor instead')
const ListFunctionUrlConfigsRequest$json = {
  '1': 'ListFunctionUrlConfigsRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListFunctionUrlConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionUrlConfigsRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0RnVuY3Rpb25VcmxDb25maWdzUmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4i/3wEgAS'
        'gJUgxmdW5jdGlvbm5hbWUSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMY'
        'lNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listFunctionUrlConfigsResponseDescriptor instead')
const ListFunctionUrlConfigsResponse$json = {
  '1': 'ListFunctionUrlConfigsResponse',
  '2': [
    {
      '1': 'functionurlconfigs',
      '3': 300063654,
      '4': 3,
      '5': 11,
      '6': '.lambda.FunctionUrlConfig',
      '10': 'functionurlconfigs'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListFunctionUrlConfigsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionUrlConfigsResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0RnVuY3Rpb25VcmxDb25maWdzUmVzcG9uc2USTQoSZnVuY3Rpb251cmxjb25maWdzGK'
        'a3io8BIAMoCzIZLmxhbWJkYS5GdW5jdGlvblVybENvbmZpZ1ISZnVuY3Rpb251cmxjb25maWdz'
        'EiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2Vy');

@$core.Deprecated(
    'Use listFunctionVersionsByCapacityProviderRequestDescriptor instead')
const ListFunctionVersionsByCapacityProviderRequest$json = {
  '1': 'ListFunctionVersionsByCapacityProviderRequest',
  '2': [
    {
      '1': 'capacityprovidername',
      '3': 466230132,
      '4': 1,
      '5': 9,
      '10': 'capacityprovidername'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListFunctionVersionsByCapacityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listFunctionVersionsByCapacityProviderRequestDescriptor =
    $convert.base64Decode(
        'Ci1MaXN0RnVuY3Rpb25WZXJzaW9uc0J5Q2FwYWNpdHlQcm92aWRlclJlcXVlc3QSNgoUY2FwYW'
        'NpdHlwcm92aWRlcm5hbWUY9Lao3gEgASgJUhRjYXBhY2l0eXByb3ZpZGVybmFtZRIZCgZtYXJr'
        'ZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated(
    'Use listFunctionVersionsByCapacityProviderResponseDescriptor instead')
const ListFunctionVersionsByCapacityProviderResponse$json = {
  '1': 'ListFunctionVersionsByCapacityProviderResponse',
  '2': [
    {
      '1': 'capacityproviderarn',
      '3': 109580344,
      '4': 1,
      '5': 9,
      '10': 'capacityproviderarn'
    },
    {
      '1': 'functionversions',
      '3': 306839073,
      '4': 3,
      '5': 11,
      '6': '.lambda.FunctionVersionsByCapacityProviderListItem',
      '10': 'functionversions'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListFunctionVersionsByCapacityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listFunctionVersionsByCapacityProviderResponseDescriptor =
    $convert.base64Decode(
        'Ci5MaXN0RnVuY3Rpb25WZXJzaW9uc0J5Q2FwYWNpdHlQcm92aWRlclJlc3BvbnNlEjMKE2NhcG'
        'FjaXR5cHJvdmlkZXJhcm4YuKCgNCABKAlSE2NhcGFjaXR5cHJvdmlkZXJhcm4SYgoQZnVuY3Rp'
        'b252ZXJzaW9ucxih/KeSASADKAsyMi5sYW1iZGEuRnVuY3Rpb25WZXJzaW9uc0J5Q2FwYWNpdH'
        'lQcm92aWRlckxpc3RJdGVtUhBmdW5jdGlvbnZlcnNpb25zEiIKCm5leHRtYXJrZXIYo4Gu/QEg'
        'ASgJUgpuZXh0bWFya2Vy');

@$core
    .Deprecated('Use listFunctionsByCodeSigningConfigRequestDescriptor instead')
const ListFunctionsByCodeSigningConfigRequest$json = {
  '1': 'ListFunctionsByCodeSigningConfigRequest',
  '2': [
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListFunctionsByCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionsByCodeSigningConfigRequestDescriptor =
    $convert.base64Decode(
        'CidMaXN0RnVuY3Rpb25zQnlDb2RlU2lnbmluZ0NvbmZpZ1JlcXVlc3QSNgoUY29kZXNpZ25pbm'
        'djb25maWdhcm4Ywfz38AEgASgJUhRjb2Rlc2lnbmluZ2NvbmZpZ2FybhIZCgZtYXJrZXIYuN3N'
        'KiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated(
    'Use listFunctionsByCodeSigningConfigResponseDescriptor instead')
const ListFunctionsByCodeSigningConfigResponse$json = {
  '1': 'ListFunctionsByCodeSigningConfigResponse',
  '2': [
    {'1': 'functionarns', '3': 419166778, '4': 3, '5': 9, '10': 'functionarns'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListFunctionsByCodeSigningConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionsByCodeSigningConfigResponseDescriptor =
    $convert.base64Decode(
        'CihMaXN0RnVuY3Rpb25zQnlDb2RlU2lnbmluZ0NvbmZpZ1Jlc3BvbnNlEiYKDGZ1bmN0aW9uYX'
        'Jucxi69O/HASADKAlSDGZ1bmN0aW9uYXJucxIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4'
        'dG1hcmtlcg==');

@$core.Deprecated('Use listFunctionsRequestDescriptor instead')
const ListFunctionsRequest$json = {
  '1': 'ListFunctionsRequest',
  '2': [
    {
      '1': 'functionversion',
      '3': 365780244,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionVersion',
      '10': 'functionversion'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'masterregion', '3': 114009692, '4': 1, '5': 9, '10': 'masterregion'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListFunctionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0RnVuY3Rpb25zUmVxdWVzdBJFCg9mdW5jdGlvbnZlcnNpb24YlLq1rgEgASgOMhcubG'
    'FtYmRhLkZ1bmN0aW9uVmVyc2lvblIPZnVuY3Rpb252ZXJzaW9uEhkKBm1hcmtlchi43c0qIAEo'
    'CVIGbWFya2VyEiUKDG1hc3RlcnJlZ2lvbhjczK42IAEoCVIMbWFzdGVycmVnaW9uEh4KCG1heG'
    'l0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listFunctionsResponseDescriptor instead')
const ListFunctionsResponse$json = {
  '1': 'ListFunctionsResponse',
  '2': [
    {
      '1': 'functions',
      '3': 165123823,
      '4': 3,
      '5': 11,
      '6': '.lambda.FunctionConfiguration',
      '10': 'functions'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListFunctionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFunctionsResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0RnVuY3Rpb25zUmVzcG9uc2USPgoJZnVuY3Rpb25zGO+t3k4gAygLMh0ubGFtYmRhLk'
    'Z1bmN0aW9uQ29uZmlndXJhdGlvblIJZnVuY3Rpb25zEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJ'
    'UgpuZXh0bWFya2Vy');

@$core.Deprecated('Use listLayerVersionsRequestDescriptor instead')
const ListLayerVersionsRequest$json = {
  '1': 'ListLayerVersionsRequest',
  '2': [
    {
      '1': 'compatiblearchitecture',
      '3': 27235489,
      '4': 1,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'compatiblearchitecture'
    },
    {
      '1': 'compatibleruntime',
      '3': 238958294,
      '4': 1,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'compatibleruntime'
    },
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListLayerVersionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLayerVersionsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0TGF5ZXJWZXJzaW9uc1JlcXVlc3QSTwoWY29tcGF0aWJsZWFyY2hpdGVjdHVyZRihqf'
    '4MIAEoDjIULmxhbWJkYS5BcmNoaXRlY3R1cmVSFmNvbXBhdGlibGVhcmNoaXRlY3R1cmUSQAoR'
    'Y29tcGF0aWJsZXJ1bnRpbWUY1u34cSABKA4yDy5sYW1iZGEuUnVudGltZVIRY29tcGF0aWJsZX'
    'J1bnRpbWUSIAoJbGF5ZXJuYW1lGJaqmO0BIAEoCVIJbGF5ZXJuYW1lEhkKBm1hcmtlchi43c0q'
    'IAEoCVIGbWFya2VyEh4KCG1heGl0ZW1zGJTW2vEBIAEoBVIIbWF4aXRlbXM=');

@$core.Deprecated('Use listLayerVersionsResponseDescriptor instead')
const ListLayerVersionsResponse$json = {
  '1': 'ListLayerVersionsResponse',
  '2': [
    {
      '1': 'layerversions',
      '3': 409515564,
      '4': 3,
      '5': 11,
      '6': '.lambda.LayerVersionsListItem',
      '10': 'layerversions'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListLayerVersionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLayerVersionsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0TGF5ZXJWZXJzaW9uc1Jlc3BvbnNlEkcKDWxheWVydmVyc2lvbnMYrOyiwwEgAygLMh'
    '0ubGFtYmRhLkxheWVyVmVyc2lvbnNMaXN0SXRlbVINbGF5ZXJ2ZXJzaW9ucxIiCgpuZXh0bWFy'
    'a2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlcg==');

@$core.Deprecated('Use listLayersRequestDescriptor instead')
const ListLayersRequest$json = {
  '1': 'ListLayersRequest',
  '2': [
    {
      '1': 'compatiblearchitecture',
      '3': 27235489,
      '4': 1,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'compatiblearchitecture'
    },
    {
      '1': 'compatibleruntime',
      '3': 238958294,
      '4': 1,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'compatibleruntime'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListLayersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLayersRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0TGF5ZXJzUmVxdWVzdBJPChZjb21wYXRpYmxlYXJjaGl0ZWN0dXJlGKGp/gwgASgOMh'
    'QubGFtYmRhLkFyY2hpdGVjdHVyZVIWY29tcGF0aWJsZWFyY2hpdGVjdHVyZRJAChFjb21wYXRp'
    'YmxlcnVudGltZRjW7fhxIAEoDjIPLmxhbWJkYS5SdW50aW1lUhFjb21wYXRpYmxlcnVudGltZR'
    'IZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIeCghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0'
    'ZW1z');

@$core.Deprecated('Use listLayersResponseDescriptor instead')
const ListLayersResponse$json = {
  '1': 'ListLayersResponse',
  '2': [
    {
      '1': 'layers',
      '3': 478144896,
      '4': 3,
      '5': 11,
      '6': '.lambda.LayersListItem',
      '10': 'layers'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
  ],
};

/// Descriptor for `ListLayersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listLayersResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0TGF5ZXJzUmVzcG9uc2USMgoGbGF5ZXJzGIDT/+MBIAMoCzIWLmxhbWJkYS5MYXllcn'
    'NMaXN0SXRlbVIGbGF5ZXJzEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2Vy');

@$core.Deprecated(
    'Use listProvisionedConcurrencyConfigsRequestDescriptor instead')
const ListProvisionedConcurrencyConfigsRequest$json = {
  '1': 'ListProvisionedConcurrencyConfigsRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListProvisionedConcurrencyConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listProvisionedConcurrencyConfigsRequestDescriptor =
    $convert.base64Decode(
        'CihMaXN0UHJvdmlzaW9uZWRDb25jdXJyZW5jeUNvbmZpZ3NSZXF1ZXN0EiYKDGZ1bmN0aW9ubm'
        'FtZRijiL/fASABKAlSDGZ1bmN0aW9ubmFtZRIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlchIe'
        'CghtYXhpdGVtcxiU1trxASABKAVSCG1heGl0ZW1z');

@$core.Deprecated(
    'Use listProvisionedConcurrencyConfigsResponseDescriptor instead')
const ListProvisionedConcurrencyConfigsResponse$json = {
  '1': 'ListProvisionedConcurrencyConfigsResponse',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {
      '1': 'provisionedconcurrencyconfigs',
      '3': 289785780,
      '4': 3,
      '5': 11,
      '6': '.lambda.ProvisionedConcurrencyConfigListItem',
      '10': 'provisionedconcurrencyconfigs'
    },
  ],
};

/// Descriptor for `ListProvisionedConcurrencyConfigsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listProvisionedConcurrencyConfigsResponseDescriptor = $convert.base64Decode(
        'CilMaXN0UHJvdmlzaW9uZWRDb25jdXJyZW5jeUNvbmZpZ3NSZXNwb25zZRIiCgpuZXh0bWFya2'
        'VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchJ2Ch1wcm92aXNpb25lZGNvbmN1cnJlbmN5Y29uZmln'
        'cxi0j5eKASADKAsyLC5sYW1iZGEuUHJvdmlzaW9uZWRDb25jdXJyZW5jeUNvbmZpZ0xpc3RJdG'
        'VtUh1wcm92aXNpb25lZGNvbmN1cnJlbmN5Y29uZmlncw==');

@$core.Deprecated('Use listTagsRequestDescriptor instead')
const ListTagsRequest$json = {
  '1': 'ListTagsRequest',
  '2': [
    {'1': 'resource', '3': 61921302, '4': 1, '5': 9, '10': 'resource'},
  ],
};

/// Descriptor for `ListTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsRequestDescriptor = $convert.base64Decode(
    'Cg9MaXN0VGFnc1JlcXVlc3QSHQoIcmVzb3VyY2UYlrDDHSABKAlSCHJlc291cmNl');

@$core.Deprecated('Use listTagsResponseDescriptor instead')
const ListTagsResponse$json = {
  '1': 'ListTagsResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.lambda.ListTagsResponse.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [ListTagsResponse_TagsEntry$json],
};

@$core.Deprecated('Use listTagsResponseDescriptor instead')
const ListTagsResponse_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsResponseDescriptor = $convert.base64Decode(
    'ChBMaXN0VGFnc1Jlc3BvbnNlEjoKBHRhZ3MYwcH2tQEgAygLMiIubGFtYmRhLkxpc3RUYWdzUm'
    'VzcG9uc2UuVGFnc0VudHJ5UgR0YWdzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIU'
    'CgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use listVersionsByFunctionRequestDescriptor instead')
const ListVersionsByFunctionRequest$json = {
  '1': 'ListVersionsByFunctionRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxitems', '3': 506899220, '4': 1, '5': 5, '10': 'maxitems'},
  ],
};

/// Descriptor for `ListVersionsByFunctionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVersionsByFunctionRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0VmVyc2lvbnNCeUZ1bmN0aW9uUmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4i/3wEgAS'
        'gJUgxmdW5jdGlvbm5hbWUSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXISHgoIbWF4aXRlbXMY'
        'lNba8QEgASgFUghtYXhpdGVtcw==');

@$core.Deprecated('Use listVersionsByFunctionResponseDescriptor instead')
const ListVersionsByFunctionResponse$json = {
  '1': 'ListVersionsByFunctionResponse',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {
      '1': 'versions',
      '3': 252099085,
      '4': 3,
      '5': 11,
      '6': '.lambda.FunctionConfiguration',
      '10': 'versions'
    },
  ],
};

/// Descriptor for `ListVersionsByFunctionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVersionsByFunctionResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0VmVyc2lvbnNCeUZ1bmN0aW9uUmVzcG9uc2USIgoKbmV4dG1hcmtlchijga79ASABKA'
        'lSCm5leHRtYXJrZXISPAoIdmVyc2lvbnMYjfSaeCADKAsyHS5sYW1iZGEuRnVuY3Rpb25Db25m'
        'aWd1cmF0aW9uUgh2ZXJzaW9ucw==');

@$core.Deprecated('Use loggingConfigDescriptor instead')
const LoggingConfig$json = {
  '1': 'LoggingConfig',
  '2': [
    {
      '1': 'applicationloglevel',
      '3': 282824602,
      '4': 1,
      '5': 14,
      '6': '.lambda.ApplicationLogLevel',
      '10': 'applicationloglevel'
    },
    {
      '1': 'logformat',
      '3': 198463503,
      '4': 1,
      '5': 14,
      '6': '.lambda.LogFormat',
      '10': 'logformat'
    },
    {'1': 'loggroup', '3': 148580073, '4': 1, '5': 9, '10': 'loggroup'},
    {
      '1': 'systemloglevel',
      '3': 530478525,
      '4': 1,
      '5': 14,
      '6': '.lambda.SystemLogLevel',
      '10': 'systemloglevel'
    },
  ],
};

/// Descriptor for `LoggingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loggingConfigDescriptor = $convert.base64Decode(
    'Cg1Mb2dnaW5nQ29uZmlnElEKE2FwcGxpY2F0aW9ubG9nbGV2ZWwYmp/uhgEgASgOMhsubGFtYm'
    'RhLkFwcGxpY2F0aW9uTG9nTGV2ZWxSE2FwcGxpY2F0aW9ubG9nbGV2ZWwSMgoJbG9nZm9ybWF0'
    'GI+g0V4gASgOMhEubGFtYmRhLkxvZ0Zvcm1hdFIJbG9nZm9ybWF0Eh0KCGxvZ2dyb3VwGOnN7E'
    'YgASgJUghsb2dncm91cBJCCg5zeXN0ZW1sb2dsZXZlbBi96/n8ASABKA4yFi5sYW1iZGEuU3lz'
    'dGVtTG9nTGV2ZWxSDnN5c3RlbWxvZ2xldmVs');

@$core.Deprecated('Use noPublishedVersionExceptionDescriptor instead')
const NoPublishedVersionException$json = {
  '1': 'NoPublishedVersionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `NoPublishedVersionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noPublishedVersionExceptionDescriptor =
    $convert.base64Decode(
        'ChtOb1B1Ymxpc2hlZFZlcnNpb25FeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZRIWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use onFailureDescriptor instead')
const OnFailure$json = {
  '1': 'OnFailure',
  '2': [
    {'1': 'destination', '3': 457443680, '4': 1, '5': 9, '10': 'destination'},
  ],
};

/// Descriptor for `OnFailure`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List onFailureDescriptor = $convert.base64Decode(
    'CglPbkZhaWx1cmUSJAoLZGVzdGluYXRpb24Y4JKQ2gEgASgJUgtkZXN0aW5hdGlvbg==');

@$core.Deprecated('Use onSuccessDescriptor instead')
const OnSuccess$json = {
  '1': 'OnSuccess',
  '2': [
    {'1': 'destination', '3': 457443680, '4': 1, '5': 9, '10': 'destination'},
  ],
};

/// Descriptor for `OnSuccess`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List onSuccessDescriptor = $convert.base64Decode(
    'CglPblN1Y2Nlc3MSJAoLZGVzdGluYXRpb24Y4JKQ2gEgASgJUgtkZXN0aW5hdGlvbg==');

@$core.Deprecated('Use operationDescriptor instead')
const Operation$json = {
  '1': 'Operation',
  '2': [
    {
      '1': 'callbackdetails',
      '3': 167924733,
      '4': 1,
      '5': 11,
      '6': '.lambda.CallbackDetails',
      '10': 'callbackdetails'
    },
    {
      '1': 'chainedinvokedetails',
      '3': 26439088,
      '4': 1,
      '5': 11,
      '6': '.lambda.ChainedInvokeDetails',
      '10': 'chainedinvokedetails'
    },
    {
      '1': 'contextdetails',
      '3': 12172411,
      '4': 1,
      '5': 11,
      '6': '.lambda.ContextDetails',
      '10': 'contextdetails'
    },
    {'1': 'endtimestamp', '3': 340967267, '4': 1, '5': 9, '10': 'endtimestamp'},
    {
      '1': 'executiondetails',
      '3': 308036630,
      '4': 1,
      '5': 11,
      '6': '.lambda.ExecutionDetails',
      '10': 'executiondetails'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'parentid', '3': 182711433, '4': 1, '5': 9, '10': 'parentid'},
    {
      '1': 'starttimestamp',
      '3': 392954318,
      '4': 1,
      '5': 9,
      '10': 'starttimestamp'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.lambda.OperationStatus',
      '10': 'status'
    },
    {
      '1': 'stepdetails',
      '3': 448807816,
      '4': 1,
      '5': 11,
      '6': '.lambda.StepDetails',
      '10': 'stepdetails'
    },
    {'1': 'subtype', '3': 152988350, '4': 1, '5': 9, '10': 'subtype'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.lambda.OperationType',
      '10': 'type'
    },
    {
      '1': 'waitdetails',
      '3': 37054967,
      '4': 1,
      '5': 11,
      '6': '.lambda.WaitDetails',
      '10': 'waitdetails'
    },
  ],
};

/// Descriptor for `Operation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List operationDescriptor = $convert.base64Decode(
    'CglPcGVyYXRpb24SRAoPY2FsbGJhY2tkZXRhaWxzGP2niVAgASgLMhcubGFtYmRhLkNhbGxiYW'
    'NrRGV0YWlsc1IPY2FsbGJhY2tkZXRhaWxzElMKFGNoYWluZWRpbnZva2VkZXRhaWxzGLDbzQwg'
    'ASgLMhwubGFtYmRhLkNoYWluZWRJbnZva2VEZXRhaWxzUhRjaGFpbmVkaW52b2tlZGV0YWlscx'
    'JBCg5jb250ZXh0ZGV0YWlscxj7+OYFIAEoCzIWLmxhbWJkYS5Db250ZXh0RGV0YWlsc1IOY29u'
    'dGV4dGRldGFpbHMSJgoMZW5kdGltZXN0YW1wGOP+yqIBIAEoCVIMZW5kdGltZXN0YW1wEkgKEG'
    'V4ZWN1dGlvbmRldGFpbHMYlojxkgEgASgLMhgubGFtYmRhLkV4ZWN1dGlvbkRldGFpbHNSEGV4'
    'ZWN1dGlvbmRldGFpbHMSEgoCaWQYgfKitwEgASgJUgJpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW'
    '1lEh0KCHBhcmVudGlkGInpj1cgASgJUghwYXJlbnRpZBIqCg5zdGFydHRpbWVzdGFtcBjOg7C7'
    'ASABKAlSDnN0YXJ0dGltZXN0YW1wEjIKBnN0YXR1cxiQ5PsCIAEoDjIXLmxhbWJkYS5PcGVyYX'
    'Rpb25TdGF0dXNSBnN0YXR1cxI5CgtzdGVwZGV0YWlscxiIh4HWASABKAsyEy5sYW1iZGEuU3Rl'
    'cERldGFpbHNSC3N0ZXBkZXRhaWxzEhsKB3N1YnR5cGUYvtX5SCABKAlSB3N1YnR5cGUSLQoEdH'
    'lwZRjuoNeKASABKA4yFS5sYW1iZGEuT3BlcmF0aW9uVHlwZVIEdHlwZRI4Cgt3YWl0ZGV0YWls'
    'cxj309URIAEoCzITLmxhbWJkYS5XYWl0RGV0YWlsc1ILd2FpdGRldGFpbHM=');

@$core.Deprecated('Use operationUpdateDescriptor instead')
const OperationUpdate$json = {
  '1': 'OperationUpdate',
  '2': [
    {
      '1': 'action',
      '3': 175614240,
      '4': 1,
      '5': 14,
      '6': '.lambda.OperationAction',
      '10': 'action'
    },
    {
      '1': 'callbackoptions',
      '3': 490425321,
      '4': 1,
      '5': 11,
      '6': '.lambda.CallbackOptions',
      '10': 'callbackoptions'
    },
    {
      '1': 'chainedinvokeoptions',
      '3': 241465076,
      '4': 1,
      '5': 11,
      '6': '.lambda.ChainedInvokeOptions',
      '10': 'chainedinvokeoptions'
    },
    {
      '1': 'contextoptions',
      '3': 181976067,
      '4': 1,
      '5': 11,
      '6': '.lambda.ContextOptions',
      '10': 'contextoptions'
    },
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '10': 'error'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'parentid', '3': 182711433, '4': 1, '5': 9, '10': 'parentid'},
    {'1': 'payload', '3': 6526790, '4': 1, '5': 9, '10': 'payload'},
    {
      '1': 'stepoptions',
      '3': 160655340,
      '4': 1,
      '5': 11,
      '6': '.lambda.StepOptions',
      '10': 'stepoptions'
    },
    {'1': 'subtype', '3': 152988350, '4': 1, '5': 9, '10': 'subtype'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.lambda.OperationType',
      '10': 'type'
    },
    {
      '1': 'waitoptions',
      '3': 240436735,
      '4': 1,
      '5': 11,
      '6': '.lambda.WaitOptions',
      '10': 'waitoptions'
    },
  ],
};

/// Descriptor for `OperationUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List operationUpdateDescriptor = $convert.base64Decode(
    'Cg9PcGVyYXRpb25VcGRhdGUSMgoGYWN0aW9uGKDS3lMgASgOMhcubGFtYmRhLk9wZXJhdGlvbk'
    'FjdGlvblIGYWN0aW9uEkUKD2NhbGxiYWNrb3B0aW9ucxjpl+3pASABKAsyFy5sYW1iZGEuQ2Fs'
    'bGJhY2tPcHRpb25zUg9jYWxsYmFja29wdGlvbnMSUwoUY2hhaW5lZGludm9rZW9wdGlvbnMY9O'
    '2RcyABKAsyHC5sYW1iZGEuQ2hhaW5lZEludm9rZU9wdGlvbnNSFGNoYWluZWRpbnZva2VvcHRp'
    'b25zEkEKDmNvbnRleHRvcHRpb25zGIP44lYgASgLMhYubGFtYmRhLkNvbnRleHRPcHRpb25zUg'
    '5jb250ZXh0b3B0aW9ucxItCgVlcnJvchjyubacASABKAsyEy5sYW1iZGEuRXJyb3JPYmplY3RS'
    'BWVycm9yEhIKAmlkGIHyorcBIAEoCVICaWQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIdCghwYX'
    'JlbnRpZBiJ6Y9XIAEoCVIIcGFyZW50aWQSGwoHcGF5bG9hZBjGro4DIAEoCVIHcGF5bG9hZBI4'
    'CgtzdGVwb3B0aW9ucxjsz81MIAEoCzITLmxhbWJkYS5TdGVwT3B0aW9uc1ILc3RlcG9wdGlvbn'
    'MSGwoHc3VidHlwZRi+1flIIAEoCVIHc3VidHlwZRItCgR0eXBlGO6g14oBIAEoDjIVLmxhbWJk'
    'YS5PcGVyYXRpb25UeXBlUgR0eXBlEjgKC3dhaXRvcHRpb25zGP+L03IgASgLMhMubGFtYmRhLl'
    'dhaXRPcHRpb25zUgt3YWl0b3B0aW9ucw==');

@$core.Deprecated('Use policyLengthExceededExceptionDescriptor instead')
const PolicyLengthExceededException$json = {
  '1': 'PolicyLengthExceededException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PolicyLengthExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyLengthExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1Qb2xpY3lMZW5ndGhFeGNlZWRlZEV4Y2VwdGlvbhIWCgR0eXBlGO6g14oBIAEoCVIEdHlwZR'
        'IbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use preconditionFailedExceptionDescriptor instead')
const PreconditionFailedException$json = {
  '1': 'PreconditionFailedException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PreconditionFailedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List preconditionFailedExceptionDescriptor =
    $convert.base64Decode(
        'ChtQcmVjb25kaXRpb25GYWlsZWRFeGNlcHRpb24SFgoEdHlwZRjuoNeKASABKAlSBHR5cGUSGw'
        'oHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use provisionedConcurrencyConfigListItemDescriptor instead')
const ProvisionedConcurrencyConfigListItem$json = {
  '1': 'ProvisionedConcurrencyConfigListItem',
  '2': [
    {
      '1': 'allocatedprovisionedconcurrentexecutions',
      '3': 468402411,
      '4': 1,
      '5': 5,
      '10': 'allocatedprovisionedconcurrentexecutions'
    },
    {
      '1': 'availableprovisionedconcurrentexecutions',
      '3': 32168809,
      '4': 1,
      '5': 5,
      '10': 'availableprovisionedconcurrentexecutions'
    },
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'requestedprovisionedconcurrentexecutions',
      '3': 58431158,
      '4': 1,
      '5': 5,
      '10': 'requestedprovisionedconcurrentexecutions'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.lambda.ProvisionedConcurrencyStatusEnum',
      '10': 'status'
    },
    {'1': 'statusreason', '3': 139234172, '4': 1, '5': 9, '10': 'statusreason'},
  ],
};

/// Descriptor for `ProvisionedConcurrencyConfigListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedConcurrencyConfigListItemDescriptor = $convert.base64Decode(
    'CiRQcm92aXNpb25lZENvbmN1cnJlbmN5Q29uZmlnTGlzdEl0ZW0SXgooYWxsb2NhdGVkcHJvdm'
    'lzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxjrga3fASABKAVSKGFsbG9jYXRlZHByb3Zpc2lv'
    'bmVkY29uY3VycmVudGV4ZWN1dGlvbnMSXQooYXZhaWxhYmxlcHJvdmlzaW9uZWRjb25jdXJyZW'
    '50ZXhlY3V0aW9ucxjptqsPIAEoBVIoYXZhaWxhYmxlcHJvdmlzaW9uZWRjb25jdXJyZW50ZXhl'
    'Y3V0aW9ucxIkCgtmdW5jdGlvbmFybhjxyY62ASABKAlSC2Z1bmN0aW9uYXJuEiYKDGxhc3Rtb2'
    'RpZmllZBinnPzOASABKAlSDGxhc3Rtb2RpZmllZBJdCihyZXF1ZXN0ZWRwcm92aXNpb25lZGNv'
    'bmN1cnJlbnRleGVjdXRpb25zGLat7hsgASgFUihyZXF1ZXN0ZWRwcm92aXNpb25lZGNvbmN1cn'
    'JlbnRleGVjdXRpb25zEkMKBnN0YXR1cxiQ5PsCIAEoDjIoLmxhbWJkYS5Qcm92aXNpb25lZENv'
    'bmN1cnJlbmN5U3RhdHVzRW51bVIGc3RhdHVzEiUKDHN0YXR1c3JlYXNvbhj8lrJCIAEoCVIMc3'
    'RhdHVzcmVhc29u');

@$core.Deprecated(
    'Use provisionedConcurrencyConfigNotFoundExceptionDescriptor instead')
const ProvisionedConcurrencyConfigNotFoundException$json = {
  '1': 'ProvisionedConcurrencyConfigNotFoundException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ProvisionedConcurrencyConfigNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    provisionedConcurrencyConfigNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ci1Qcm92aXNpb25lZENvbmN1cnJlbmN5Q29uZmlnTm90Rm91bmRFeGNlcHRpb24SFgoEdHlwZR'
        'juoNeKASABKAlSBHR5cGUSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use provisionedPollerConfigDescriptor instead')
const ProvisionedPollerConfig$json = {
  '1': 'ProvisionedPollerConfig',
  '2': [
    {
      '1': 'maximumpollers',
      '3': 170898825,
      '4': 1,
      '5': 5,
      '10': 'maximumpollers'
    },
    {
      '1': 'minimumpollers',
      '3': 302938719,
      '4': 1,
      '5': 5,
      '10': 'minimumpollers'
    },
    {
      '1': 'pollergroupname',
      '3': 316867254,
      '4': 1,
      '5': 9,
      '10': 'pollergroupname'
    },
  ],
};

/// Descriptor for `ProvisionedPollerConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedPollerConfigDescriptor = $convert.base64Decode(
    'ChdQcm92aXNpb25lZFBvbGxlckNvbmZpZxIpCg5tYXhpbXVtcG9sbGVycxiJ675RIAEoBVIObW'
    'F4aW11bXBvbGxlcnMSKgoObWluaW11bXBvbGxlcnMY3/S5kAEgASgFUg5taW5pbXVtcG9sbGVy'
    'cxIsCg9wb2xsZXJncm91cG5hbWUYtoWMlwEgASgJUg9wb2xsZXJncm91cG5hbWU=');

@$core.Deprecated('Use publishLayerVersionRequestDescriptor instead')
const PublishLayerVersionRequest$json = {
  '1': 'PublishLayerVersionRequest',
  '2': [
    {
      '1': 'compatiblearchitectures',
      '3': 20400170,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'compatiblearchitectures'
    },
    {
      '1': 'compatibleruntimes',
      '3': 300943187,
      '4': 3,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'compatibleruntimes'
    },
    {
      '1': 'content',
      '3': 23568227,
      '4': 1,
      '5': 11,
      '6': '.lambda.LayerVersionContentInput',
      '10': 'content'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
    {'1': 'licenseinfo', '3': 133903897, '4': 1, '5': 9, '10': 'licenseinfo'},
  ],
};

/// Descriptor for `PublishLayerVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishLayerVersionRequestDescriptor = $convert.base64Decode(
    'ChpQdWJsaXNoTGF5ZXJWZXJzaW9uUmVxdWVzdBJRChdjb21wYXRpYmxlYXJjaGl0ZWN0dXJlcx'
    'iqkN0JIAMoDjIULmxhbWJkYS5BcmNoaXRlY3R1cmVSF2NvbXBhdGlibGVhcmNoaXRlY3R1cmVz'
    'EkMKEmNvbXBhdGlibGVydW50aW1lcxjTjsCPASADKA4yDy5sYW1iZGEuUnVudGltZVISY29tcG'
    'F0aWJsZXJ1bnRpbWVzEj0KB2NvbnRlbnQY476eCyABKAsyIC5sYW1iZGEuTGF5ZXJWZXJzaW9u'
    'Q29udGVudElucHV0Ugdjb250ZW50EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdG'
    'lvbhIgCglsYXllcm5hbWUYlqqY7QEgASgJUglsYXllcm5hbWUSIwoLbGljZW5zZWluZm8Ymezs'
    'PyABKAlSC2xpY2Vuc2VpbmZv');

@$core.Deprecated('Use publishLayerVersionResponseDescriptor instead')
const PublishLayerVersionResponse$json = {
  '1': 'PublishLayerVersionResponse',
  '2': [
    {
      '1': 'compatiblearchitectures',
      '3': 20400170,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'compatiblearchitectures'
    },
    {
      '1': 'compatibleruntimes',
      '3': 300943187,
      '4': 3,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'compatibleruntimes'
    },
    {
      '1': 'content',
      '3': 23568227,
      '4': 1,
      '5': 11,
      '6': '.lambda.LayerVersionContentOutput',
      '10': 'content'
    },
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'layerarn', '3': 319652978, '4': 1, '5': 9, '10': 'layerarn'},
    {
      '1': 'layerversionarn',
      '3': 174478342,
      '4': 1,
      '5': 9,
      '10': 'layerversionarn'
    },
    {'1': 'licenseinfo', '3': 133903897, '4': 1, '5': 9, '10': 'licenseinfo'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `PublishLayerVersionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishLayerVersionResponseDescriptor = $convert.base64Decode(
    'ChtQdWJsaXNoTGF5ZXJWZXJzaW9uUmVzcG9uc2USUQoXY29tcGF0aWJsZWFyY2hpdGVjdHVyZX'
    'MYqpDdCSADKA4yFC5sYW1iZGEuQXJjaGl0ZWN0dXJlUhdjb21wYXRpYmxlYXJjaGl0ZWN0dXJl'
    'cxJDChJjb21wYXRpYmxlcnVudGltZXMY047AjwEgAygOMg8ubGFtYmRhLlJ1bnRpbWVSEmNvbX'
    'BhdGlibGVydW50aW1lcxI+Cgdjb250ZW50GOO+ngsgASgLMiEubGFtYmRhLkxheWVyVmVyc2lv'
    'bkNvbnRlbnRPdXRwdXRSB2NvbnRlbnQSJAoLY3JlYXRlZGRhdGUYsLDnxgEgASgJUgtjcmVhdG'
    'VkZGF0ZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SHgoIbGF5ZXJhcm4Y'
    '8oi2mAEgASgJUghsYXllcmFybhIrCg9sYXllcnZlcnNpb25hcm4YhqiZUyABKAlSD2xheWVydm'
    'Vyc2lvbmFybhIjCgtsaWNlbnNlaW5mbxiZ7Ow/IAEoCVILbGljZW5zZWluZm8SHAoHdmVyc2lv'
    'bhi4qrfuASABKANSB3ZlcnNpb24=');

@$core.Deprecated('Use publishVersionRequestDescriptor instead')
const PublishVersionRequest$json = {
  '1': 'PublishVersionRequest',
  '2': [
    {'1': 'codesha256', '3': 46450860, '4': 1, '5': 9, '10': 'codesha256'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'publishto',
      '3': 524127682,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionVersionLatestPublished',
      '10': 'publishto'
    },
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
  ],
};

/// Descriptor for `PublishVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishVersionRequestDescriptor = $convert.base64Decode(
    'ChVQdWJsaXNoVmVyc2lvblJlcXVlc3QSIQoKY29kZXNoYTI1NhiskZMWIAEoCVIKY29kZXNoYT'
    'I1NhIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SJgoMZnVuY3Rpb25uYW1l'
    'GKOIv98BIAEoCVIMZnVuY3Rpb25uYW1lEkgKCXB1Ymxpc2h0bxjCm/b5ASABKA4yJi5sYW1iZG'
    'EuRnVuY3Rpb25WZXJzaW9uTGF0ZXN0UHVibGlzaGVkUglwdWJsaXNodG8SIgoKcmV2aXNpb25p'
    'ZBiGo57uASABKAlSCnJldmlzaW9uaWQ=');

@$core.Deprecated('Use putFunctionCodeSigningConfigRequestDescriptor instead')
const PutFunctionCodeSigningConfigRequest$json = {
  '1': 'PutFunctionCodeSigningConfigRequest',
  '2': [
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `PutFunctionCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionCodeSigningConfigRequestDescriptor =
    $convert.base64Decode(
        'CiNQdXRGdW5jdGlvbkNvZGVTaWduaW5nQ29uZmlnUmVxdWVzdBI2ChRjb2Rlc2lnbmluZ2Nvbm'
        'ZpZ2FybhjB/PfwASABKAlSFGNvZGVzaWduaW5nY29uZmlnYXJuEiYKDGZ1bmN0aW9ubmFtZRij'
        'iL/fASABKAlSDGZ1bmN0aW9ubmFtZQ==');

@$core.Deprecated('Use putFunctionCodeSigningConfigResponseDescriptor instead')
const PutFunctionCodeSigningConfigResponse$json = {
  '1': 'PutFunctionCodeSigningConfigResponse',
  '2': [
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
  ],
};

/// Descriptor for `PutFunctionCodeSigningConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionCodeSigningConfigResponseDescriptor =
    $convert.base64Decode(
        'CiRQdXRGdW5jdGlvbkNvZGVTaWduaW5nQ29uZmlnUmVzcG9uc2USNgoUY29kZXNpZ25pbmdjb2'
        '5maWdhcm4Ywfz38AEgASgJUhRjb2Rlc2lnbmluZ2NvbmZpZ2FybhImCgxmdW5jdGlvbm5hbWUY'
        'o4i/3wEgASgJUgxmdW5jdGlvbm5hbWU=');

@$core.Deprecated('Use putFunctionConcurrencyRequestDescriptor instead')
const PutFunctionConcurrencyRequest$json = {
  '1': 'PutFunctionConcurrencyRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'reservedconcurrentexecutions',
      '3': 40149212,
      '4': 1,
      '5': 5,
      '10': 'reservedconcurrentexecutions'
    },
  ],
};

/// Descriptor for `PutFunctionConcurrencyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionConcurrencyRequestDescriptor =
    $convert.base64Decode(
        'Ch1QdXRGdW5jdGlvbkNvbmN1cnJlbmN5UmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4i/3wEgAS'
        'gJUgxmdW5jdGlvbm5hbWUSRQoccmVzZXJ2ZWRjb25jdXJyZW50ZXhlY3V0aW9ucxjcwZITIAEo'
        'BVIccmVzZXJ2ZWRjb25jdXJyZW50ZXhlY3V0aW9ucw==');

@$core.Deprecated('Use putFunctionEventInvokeConfigRequestDescriptor instead')
const PutFunctionEventInvokeConfigRequest$json = {
  '1': 'PutFunctionEventInvokeConfigRequest',
  '2': [
    {
      '1': 'destinationconfig',
      '3': 184834158,
      '4': 1,
      '5': 11,
      '6': '.lambda.DestinationConfig',
      '10': 'destinationconfig'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'maximumeventageinseconds',
      '3': 393041563,
      '4': 1,
      '5': 5,
      '10': 'maximumeventageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `PutFunctionEventInvokeConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionEventInvokeConfigRequestDescriptor = $convert.base64Decode(
    'CiNQdXRGdW5jdGlvbkV2ZW50SW52b2tlQ29uZmlnUmVxdWVzdBJKChFkZXN0aW5hdGlvbmNvbm'
    'ZpZxjusJFYIAEoCzIZLmxhbWJkYS5EZXN0aW5hdGlvbkNvbmZpZ1IRZGVzdGluYXRpb25jb25m'
    'aWcSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb25uYW1lEj4KGG1heGltdW1ldm'
    'VudGFnZWluc2Vjb25kcxibrbW7ASABKAVSGG1heGltdW1ldmVudGFnZWluc2Vjb25kcxI1ChRt'
    'YXhpbXVtcmV0cnlhdHRlbXB0cxjAqLk1IAEoBVIUbWF4aW11bXJldHJ5YXR0ZW1wdHMSIAoJcX'
    'VhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlmaWVy');

@$core.Deprecated('Use putFunctionRecursionConfigRequestDescriptor instead')
const PutFunctionRecursionConfigRequest$json = {
  '1': 'PutFunctionRecursionConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'recursiveloop',
      '3': 2821758,
      '4': 1,
      '5': 14,
      '6': '.lambda.RecursiveLoop',
      '10': 'recursiveloop'
    },
  ],
};

/// Descriptor for `PutFunctionRecursionConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionRecursionConfigRequestDescriptor =
    $convert.base64Decode(
        'CiFQdXRGdW5jdGlvblJlY3Vyc2lvbkNvbmZpZ1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv9'
        '8BIAEoCVIMZnVuY3Rpb25uYW1lEj4KDXJlY3Vyc2l2ZWxvb3AY/pysASABKA4yFS5sYW1iZGEu'
        'UmVjdXJzaXZlTG9vcFINcmVjdXJzaXZlbG9vcA==');

@$core.Deprecated('Use putFunctionRecursionConfigResponseDescriptor instead')
const PutFunctionRecursionConfigResponse$json = {
  '1': 'PutFunctionRecursionConfigResponse',
  '2': [
    {
      '1': 'recursiveloop',
      '3': 2821758,
      '4': 1,
      '5': 14,
      '6': '.lambda.RecursiveLoop',
      '10': 'recursiveloop'
    },
  ],
};

/// Descriptor for `PutFunctionRecursionConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionRecursionConfigResponseDescriptor =
    $convert.base64Decode(
        'CiJQdXRGdW5jdGlvblJlY3Vyc2lvbkNvbmZpZ1Jlc3BvbnNlEj4KDXJlY3Vyc2l2ZWxvb3AY/p'
        'ysASABKA4yFS5sYW1iZGEuUmVjdXJzaXZlTG9vcFINcmVjdXJzaXZlbG9vcA==');

@$core.Deprecated('Use putFunctionScalingConfigRequestDescriptor instead')
const PutFunctionScalingConfigRequest$json = {
  '1': 'PutFunctionScalingConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'functionscalingconfig',
      '3': 125660113,
      '4': 1,
      '5': 11,
      '6': '.lambda.FunctionScalingConfig',
      '10': 'functionscalingconfig'
    },
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `PutFunctionScalingConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionScalingConfigRequestDescriptor =
    $convert.base64Decode(
        'Ch9QdXRGdW5jdGlvblNjYWxpbmdDb25maWdSZXF1ZXN0EiYKDGZ1bmN0aW9ubmFtZRijiL/fAS'
        'ABKAlSDGZ1bmN0aW9ubmFtZRJWChVmdW5jdGlvbnNjYWxpbmdjb25maWcY0df1OyABKAsyHS5s'
        'YW1iZGEuRnVuY3Rpb25TY2FsaW5nQ29uZmlnUhVmdW5jdGlvbnNjYWxpbmdjb25maWcSIAoJcX'
        'VhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlmaWVy');

@$core.Deprecated('Use putFunctionScalingConfigResponseDescriptor instead')
const PutFunctionScalingConfigResponse$json = {
  '1': 'PutFunctionScalingConfigResponse',
  '2': [
    {
      '1': 'functionstate',
      '3': 63089739,
      '4': 1,
      '5': 14,
      '6': '.lambda.State',
      '10': 'functionstate'
    },
  ],
};

/// Descriptor for `PutFunctionScalingConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putFunctionScalingConfigResponseDescriptor =
    $convert.base64Decode(
        'CiBQdXRGdW5jdGlvblNjYWxpbmdDb25maWdSZXNwb25zZRI2Cg1mdW5jdGlvbnN0YXRlGMvYih'
        '4gASgOMg0ubGFtYmRhLlN0YXRlUg1mdW5jdGlvbnN0YXRl');

@$core
    .Deprecated('Use putProvisionedConcurrencyConfigRequestDescriptor instead')
const PutProvisionedConcurrencyConfigRequest$json = {
  '1': 'PutProvisionedConcurrencyConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'provisionedconcurrentexecutions',
      '3': 184243952,
      '4': 1,
      '5': 5,
      '10': 'provisionedconcurrentexecutions'
    },
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `PutProvisionedConcurrencyConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putProvisionedConcurrencyConfigRequestDescriptor =
    $convert.base64Decode(
        'CiZQdXRQcm92aXNpb25lZENvbmN1cnJlbmN5Q29uZmlnUmVxdWVzdBImCgxmdW5jdGlvbm5hbW'
        'UYo4i/3wEgASgJUgxmdW5jdGlvbm5hbWUSSwofcHJvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0'
        'aW9ucxjwre1XIAEoBVIfcHJvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxIgCglxdWFsaW'
        'ZpZXIY4LWR+wEgASgJUglxdWFsaWZpZXI=');

@$core
    .Deprecated('Use putProvisionedConcurrencyConfigResponseDescriptor instead')
const PutProvisionedConcurrencyConfigResponse$json = {
  '1': 'PutProvisionedConcurrencyConfigResponse',
  '2': [
    {
      '1': 'allocatedprovisionedconcurrentexecutions',
      '3': 468402411,
      '4': 1,
      '5': 5,
      '10': 'allocatedprovisionedconcurrentexecutions'
    },
    {
      '1': 'availableprovisionedconcurrentexecutions',
      '3': 32168809,
      '4': 1,
      '5': 5,
      '10': 'availableprovisionedconcurrentexecutions'
    },
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'requestedprovisionedconcurrentexecutions',
      '3': 58431158,
      '4': 1,
      '5': 5,
      '10': 'requestedprovisionedconcurrentexecutions'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.lambda.ProvisionedConcurrencyStatusEnum',
      '10': 'status'
    },
    {'1': 'statusreason', '3': 139234172, '4': 1, '5': 9, '10': 'statusreason'},
  ],
};

/// Descriptor for `PutProvisionedConcurrencyConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putProvisionedConcurrencyConfigResponseDescriptor = $convert.base64Decode(
    'CidQdXRQcm92aXNpb25lZENvbmN1cnJlbmN5Q29uZmlnUmVzcG9uc2USXgooYWxsb2NhdGVkcH'
    'JvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxjrga3fASABKAVSKGFsbG9jYXRlZHByb3Zp'
    'c2lvbmVkY29uY3VycmVudGV4ZWN1dGlvbnMSXQooYXZhaWxhYmxlcHJvdmlzaW9uZWRjb25jdX'
    'JyZW50ZXhlY3V0aW9ucxjptqsPIAEoBVIoYXZhaWxhYmxlcHJvdmlzaW9uZWRjb25jdXJyZW50'
    'ZXhlY3V0aW9ucxImCgxsYXN0bW9kaWZpZWQYp5z8zgEgASgJUgxsYXN0bW9kaWZpZWQSXQoocm'
    'VxdWVzdGVkcHJvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxi2re4bIAEoBVIocmVxdWVz'
    'dGVkcHJvdmlzaW9uZWRjb25jdXJyZW50ZXhlY3V0aW9ucxJDCgZzdGF0dXMYkOT7AiABKA4yKC'
    '5sYW1iZGEuUHJvdmlzaW9uZWRDb25jdXJyZW5jeVN0YXR1c0VudW1SBnN0YXR1cxIlCgxzdGF0'
    'dXNyZWFzb24Y/JayQiABKAlSDHN0YXR1c3JlYXNvbg==');

@$core.Deprecated('Use putRuntimeManagementConfigRequestDescriptor instead')
const PutRuntimeManagementConfigRequest$json = {
  '1': 'PutRuntimeManagementConfigRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
    {
      '1': 'runtimeversionarn',
      '3': 532500669,
      '4': 1,
      '5': 9,
      '10': 'runtimeversionarn'
    },
    {
      '1': 'updateruntimeon',
      '3': 285197810,
      '4': 1,
      '5': 14,
      '6': '.lambda.UpdateRuntimeOn',
      '10': 'updateruntimeon'
    },
  ],
};

/// Descriptor for `PutRuntimeManagementConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRuntimeManagementConfigRequestDescriptor = $convert.base64Decode(
    'CiFQdXRSdW50aW1lTWFuYWdlbWVudENvbmZpZ1JlcXVlc3QSJgoMZnVuY3Rpb25uYW1lGKOIv9'
    '8BIAEoCVIMZnVuY3Rpb25uYW1lEiAKCXF1YWxpZmllchjgtZH7ASABKAlSCXF1YWxpZmllchIw'
    'ChFydW50aW1ldmVyc2lvbmFybhi9ofX9ASABKAlSEXJ1bnRpbWV2ZXJzaW9uYXJuEkUKD3VwZG'
    'F0ZXJ1bnRpbWVvbhjyi/+HASABKA4yFy5sYW1iZGEuVXBkYXRlUnVudGltZU9uUg91cGRhdGVy'
    'dW50aW1lb24=');

@$core.Deprecated('Use putRuntimeManagementConfigResponseDescriptor instead')
const PutRuntimeManagementConfigResponse$json = {
  '1': 'PutRuntimeManagementConfigResponse',
  '2': [
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {
      '1': 'runtimeversionarn',
      '3': 532500669,
      '4': 1,
      '5': 9,
      '10': 'runtimeversionarn'
    },
    {
      '1': 'updateruntimeon',
      '3': 285197810,
      '4': 1,
      '5': 14,
      '6': '.lambda.UpdateRuntimeOn',
      '10': 'updateruntimeon'
    },
  ],
};

/// Descriptor for `PutRuntimeManagementConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRuntimeManagementConfigResponseDescriptor =
    $convert.base64Decode(
        'CiJQdXRSdW50aW1lTWFuYWdlbWVudENvbmZpZ1Jlc3BvbnNlEiQKC2Z1bmN0aW9uYXJuGPHJjr'
        'YBIAEoCVILZnVuY3Rpb25hcm4SMAoRcnVudGltZXZlcnNpb25hcm4YvaH1/QEgASgJUhFydW50'
        'aW1ldmVyc2lvbmFybhJFCg91cGRhdGVydW50aW1lb24Y8ov/hwEgASgOMhcubGFtYmRhLlVwZG'
        'F0ZVJ1bnRpbWVPblIPdXBkYXRlcnVudGltZW9u');

@$core.Deprecated('Use recursiveInvocationExceptionDescriptor instead')
const RecursiveInvocationException$json = {
  '1': 'RecursiveInvocationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `RecursiveInvocationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recursiveInvocationExceptionDescriptor =
    $convert.base64Decode(
        'ChxSZWN1cnNpdmVJbnZvY2F0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2USFgoEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use removeLayerVersionPermissionRequestDescriptor instead')
const RemoveLayerVersionPermissionRequest$json = {
  '1': 'RemoveLayerVersionPermissionRequest',
  '2': [
    {'1': 'layername', '3': 497423638, '4': 1, '5': 9, '10': 'layername'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'statementid', '3': 169602348, '4': 1, '5': 9, '10': 'statementid'},
    {
      '1': 'versionnumber',
      '3': 209346775,
      '4': 1,
      '5': 3,
      '10': 'versionnumber'
    },
  ],
};

/// Descriptor for `RemoveLayerVersionPermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeLayerVersionPermissionRequestDescriptor =
    $convert.base64Decode(
        'CiNSZW1vdmVMYXllclZlcnNpb25QZXJtaXNzaW9uUmVxdWVzdBIgCglsYXllcm5hbWUYlqqY7Q'
        'EgASgJUglsYXllcm5hbWUSIgoKcmV2aXNpb25pZBiGo57uASABKAlSCnJldmlzaW9uaWQSIwoL'
        'c3RhdGVtZW50aWQYrNrvUCABKAlSC3N0YXRlbWVudGlkEicKDXZlcnNpb25udW1iZXIY18HpYy'
        'ABKANSDXZlcnNpb25udW1iZXI=');

@$core.Deprecated('Use removePermissionRequestDescriptor instead')
const RemovePermissionRequest$json = {
  '1': 'RemovePermissionRequest',
  '2': [
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'statementid', '3': 169602348, '4': 1, '5': 9, '10': 'statementid'},
  ],
};

/// Descriptor for `RemovePermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removePermissionRequestDescriptor = $convert.base64Decode(
    'ChdSZW1vdmVQZXJtaXNzaW9uUmVxdWVzdBImCgxmdW5jdGlvbm5hbWUYo4i/3wEgASgJUgxmdW'
    '5jdGlvbm5hbWUSIAoJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlmaWVyEiIKCnJldmlzaW9u'
    'aWQYhqOe7gEgASgJUgpyZXZpc2lvbmlkEiMKC3N0YXRlbWVudGlkGKza71AgASgJUgtzdGF0ZW'
    '1lbnRpZA==');

@$core.Deprecated('Use requestTooLargeExceptionDescriptor instead')
const RequestTooLargeException$json = {
  '1': 'RequestTooLargeException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RequestTooLargeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestTooLargeExceptionDescriptor =
    $convert.base64Decode(
        'ChhSZXF1ZXN0VG9vTGFyZ2VFeGNlcHRpb24SFgoEdHlwZRjuoNeKASABKAlSBHR5cGUSGwoHbW'
        'Vzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use resourceConflictExceptionDescriptor instead')
const ResourceConflictException$json = {
  '1': 'ResourceConflictException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceConflictExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZUNvbmZsaWN0RXhjZXB0aW9uEhYKBHR5cGUY7qDXigEgASgJUgR0eXBlEhsKB2'
        '1lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use resourceInUseExceptionDescriptor instead')
const ResourceInUseException$json = {
  '1': 'ResourceInUseException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `ResourceInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChZSZXNvdXJjZUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2USFg'
        'oEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'USFgoEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use resourceNotReadyExceptionDescriptor instead')
const ResourceNotReadyException$json = {
  '1': 'ResourceNotReadyException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotReadyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotReadyExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdFJlYWR5RXhjZXB0aW9uEhYKBHR5cGUY7qDXigEgASgJUgR0eXBlEhsKB2'
        '1lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use retryDetailsDescriptor instead')
const RetryDetails$json = {
  '1': 'RetryDetails',
  '2': [
    {
      '1': 'currentattempt',
      '3': 340355832,
      '4': 1,
      '5': 5,
      '10': 'currentattempt'
    },
    {
      '1': 'nextattemptdelayseconds',
      '3': 502780186,
      '4': 1,
      '5': 5,
      '10': 'nextattemptdelayseconds'
    },
  ],
};

/// Descriptor for `RetryDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List retryDetailsDescriptor = $convert.base64Decode(
    'CgxSZXRyeURldGFpbHMSKgoOY3VycmVudGF0dGVtcHQY+NWlogEgASgFUg5jdXJyZW50YXR0ZW'
    '1wdBI8ChduZXh0YXR0ZW1wdGRlbGF5c2Vjb25kcxiaot/vASABKAVSF25leHRhdHRlbXB0ZGVs'
    'YXlzZWNvbmRz');

@$core.Deprecated('Use runtimeVersionConfigDescriptor instead')
const RuntimeVersionConfig$json = {
  '1': 'RuntimeVersionConfig',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.RuntimeVersionError',
      '10': 'error'
    },
    {
      '1': 'runtimeversionarn',
      '3': 532500669,
      '4': 1,
      '5': 9,
      '10': 'runtimeversionarn'
    },
  ],
};

/// Descriptor for `RuntimeVersionConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List runtimeVersionConfigDescriptor = $convert.base64Decode(
    'ChRSdW50aW1lVmVyc2lvbkNvbmZpZxI1CgVlcnJvchjyubacASABKAsyGy5sYW1iZGEuUnVudG'
    'ltZVZlcnNpb25FcnJvclIFZXJyb3ISMAoRcnVudGltZXZlcnNpb25hcm4YvaH1/QEgASgJUhFy'
    'dW50aW1ldmVyc2lvbmFybg==');

@$core.Deprecated('Use runtimeVersionErrorDescriptor instead')
const RuntimeVersionError$json = {
  '1': 'RuntimeVersionError',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RuntimeVersionError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List runtimeVersionErrorDescriptor = $convert.base64Decode(
    'ChNSdW50aW1lVmVyc2lvbkVycm9yEh8KCWVycm9yY29kZRiZ1sMQIAEoCVIJZXJyb3Jjb2RlEh'
    'sKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use scalingConfigDescriptor instead')
const ScalingConfig$json = {
  '1': 'ScalingConfig',
  '2': [
    {
      '1': 'maximumconcurrency',
      '3': 257783285,
      '4': 1,
      '5': 5,
      '10': 'maximumconcurrency'
    },
  ],
};

/// Descriptor for `ScalingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scalingConfigDescriptor = $convert.base64Decode(
    'Cg1TY2FsaW5nQ29uZmlnEjEKEm1heGltdW1jb25jdXJyZW5jeRj16/V6IAEoBVISbWF4aW11bW'
    'NvbmN1cnJlbmN5');

@$core.Deprecated('Use selfManagedEventSourceDescriptor instead')
const SelfManagedEventSource$json = {
  '1': 'SelfManagedEventSource',
  '2': [
    {
      '1': 'endpoints',
      '3': 16210494,
      '4': 3,
      '5': 11,
      '6': '.lambda.SelfManagedEventSource.EndpointsEntry',
      '10': 'endpoints'
    },
  ],
  '3': [SelfManagedEventSource_EndpointsEntry$json],
};

@$core.Deprecated('Use selfManagedEventSourceDescriptor instead')
const SelfManagedEventSource_EndpointsEntry$json = {
  '1': 'EndpointsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SelfManagedEventSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selfManagedEventSourceDescriptor = $convert.base64Decode(
    'ChZTZWxmTWFuYWdlZEV2ZW50U291cmNlEk4KCWVuZHBvaW50cxi+tN0HIAMoCzItLmxhbWJkYS'
    '5TZWxmTWFuYWdlZEV2ZW50U291cmNlLkVuZHBvaW50c0VudHJ5UgllbmRwb2ludHMaPAoORW5k'
    'cG9pbnRzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ'
    '==');

@$core.Deprecated('Use selfManagedKafkaEventSourceConfigDescriptor instead')
const SelfManagedKafkaEventSourceConfig$json = {
  '1': 'SelfManagedKafkaEventSourceConfig',
  '2': [
    {
      '1': 'consumergroupid',
      '3': 222398686,
      '4': 1,
      '5': 9,
      '10': 'consumergroupid'
    },
    {
      '1': 'schemaregistryconfig',
      '3': 482259768,
      '4': 1,
      '5': 11,
      '6': '.lambda.KafkaSchemaRegistryConfig',
      '10': 'schemaregistryconfig'
    },
  ],
};

/// Descriptor for `SelfManagedKafkaEventSourceConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selfManagedKafkaEventSourceConfigDescriptor =
    $convert.base64Decode(
        'CiFTZWxmTWFuYWdlZEthZmthRXZlbnRTb3VyY2VDb25maWcSKwoPY29uc3VtZXJncm91cGlkGN'
        '6RhmogASgJUg9jb25zdW1lcmdyb3VwaWQSWQoUc2NoZW1hcmVnaXN0cnljb25maWcYuOb65QEg'
        'ASgLMiEubGFtYmRhLkthZmthU2NoZW1hUmVnaXN0cnlDb25maWdSFHNjaGVtYXJlZ2lzdHJ5Y2'
        '9uZmln');

@$core.Deprecated(
    'Use sendDurableExecutionCallbackFailureRequestDescriptor instead')
const SendDurableExecutionCallbackFailureRequest$json = {
  '1': 'SendDurableExecutionCallbackFailureRequest',
  '2': [
    {'1': 'callbackid', '3': 70101916, '4': 1, '5': 9, '10': 'callbackid'},
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '8': {},
      '10': 'error'
    },
  ],
};

/// Descriptor for `SendDurableExecutionCallbackFailureRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    sendDurableExecutionCallbackFailureRequestDescriptor =
    $convert.base64Decode(
        'CipTZW5kRHVyYWJsZUV4ZWN1dGlvbkNhbGxiYWNrRmFpbHVyZVJlcXVlc3QSIQoKY2FsbGJhY2'
        'tpZBic17YhIAEoCVIKY2FsbGJhY2tpZBIzCgVlcnJvchjyubacASABKAsyEy5sYW1iZGEuRXJy'
        'b3JPYmplY3RCBIi1GAFSBWVycm9y');

@$core.Deprecated(
    'Use sendDurableExecutionCallbackFailureResponseDescriptor instead')
const SendDurableExecutionCallbackFailureResponse$json = {
  '1': 'SendDurableExecutionCallbackFailureResponse',
};

/// Descriptor for `SendDurableExecutionCallbackFailureResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    sendDurableExecutionCallbackFailureResponseDescriptor =
    $convert.base64Decode(
        'CitTZW5kRHVyYWJsZUV4ZWN1dGlvbkNhbGxiYWNrRmFpbHVyZVJlc3BvbnNl');

@$core.Deprecated(
    'Use sendDurableExecutionCallbackHeartbeatRequestDescriptor instead')
const SendDurableExecutionCallbackHeartbeatRequest$json = {
  '1': 'SendDurableExecutionCallbackHeartbeatRequest',
  '2': [
    {'1': 'callbackid', '3': 70101916, '4': 1, '5': 9, '10': 'callbackid'},
  ],
};

/// Descriptor for `SendDurableExecutionCallbackHeartbeatRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    sendDurableExecutionCallbackHeartbeatRequestDescriptor =
    $convert.base64Decode(
        'CixTZW5kRHVyYWJsZUV4ZWN1dGlvbkNhbGxiYWNrSGVhcnRiZWF0UmVxdWVzdBIhCgpjYWxsYm'
        'Fja2lkGJzXtiEgASgJUgpjYWxsYmFja2lk');

@$core.Deprecated(
    'Use sendDurableExecutionCallbackHeartbeatResponseDescriptor instead')
const SendDurableExecutionCallbackHeartbeatResponse$json = {
  '1': 'SendDurableExecutionCallbackHeartbeatResponse',
};

/// Descriptor for `SendDurableExecutionCallbackHeartbeatResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    sendDurableExecutionCallbackHeartbeatResponseDescriptor =
    $convert.base64Decode(
        'Ci1TZW5kRHVyYWJsZUV4ZWN1dGlvbkNhbGxiYWNrSGVhcnRiZWF0UmVzcG9uc2U=');

@$core.Deprecated(
    'Use sendDurableExecutionCallbackSuccessRequestDescriptor instead')
const SendDurableExecutionCallbackSuccessRequest$json = {
  '1': 'SendDurableExecutionCallbackSuccessRequest',
  '2': [
    {'1': 'callbackid', '3': 70101916, '4': 1, '5': 9, '10': 'callbackid'},
    {'1': 'result', '3': 273346629, '4': 1, '5': 12, '8': {}, '10': 'result'},
  ],
};

/// Descriptor for `SendDurableExecutionCallbackSuccessRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    sendDurableExecutionCallbackSuccessRequestDescriptor =
    $convert.base64Decode(
        'CipTZW5kRHVyYWJsZUV4ZWN1dGlvbkNhbGxiYWNrU3VjY2Vzc1JlcXVlc3QSIQoKY2FsbGJhY2'
        'tpZBic17YhIAEoCVIKY2FsbGJhY2tpZBIgCgZyZXN1bHQYxeCrggEgASgMQgSItRgBUgZyZXN1'
        'bHQ=');

@$core.Deprecated(
    'Use sendDurableExecutionCallbackSuccessResponseDescriptor instead')
const SendDurableExecutionCallbackSuccessResponse$json = {
  '1': 'SendDurableExecutionCallbackSuccessResponse',
};

/// Descriptor for `SendDurableExecutionCallbackSuccessResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    sendDurableExecutionCallbackSuccessResponseDescriptor =
    $convert.base64Decode(
        'CitTZW5kRHVyYWJsZUV4ZWN1dGlvbkNhbGxiYWNrU3VjY2Vzc1Jlc3BvbnNl');

@$core.Deprecated(
    'Use serializedRequestEntityTooLargeExceptionDescriptor instead')
const SerializedRequestEntityTooLargeException$json = {
  '1': 'SerializedRequestEntityTooLargeException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SerializedRequestEntityTooLargeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serializedRequestEntityTooLargeExceptionDescriptor =
    $convert.base64Decode(
        'CihTZXJpYWxpemVkUmVxdWVzdEVudGl0eVRvb0xhcmdlRXhjZXB0aW9uEhYKBHR5cGUY7qDXig'
        'EgASgJUgR0eXBlEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use serviceExceptionDescriptor instead')
const ServiceException$json = {
  '1': 'ServiceException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `ServiceException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceExceptionDescriptor = $convert.base64Decode(
    'ChBTZXJ2aWNlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2USFgoEdHlwZR'
    'juoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use snapStartDescriptor instead')
const SnapStart$json = {
  '1': 'SnapStart',
  '2': [
    {
      '1': 'applyon',
      '3': 531091671,
      '4': 1,
      '5': 14,
      '6': '.lambda.SnapStartApplyOn',
      '10': 'applyon'
    },
  ],
};

/// Descriptor for `SnapStart`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List snapStartDescriptor = $convert.base64Decode(
    'CglTbmFwU3RhcnQSNgoHYXBwbHlvbhjXoZ/9ASABKA4yGC5sYW1iZGEuU25hcFN0YXJ0QXBwbH'
    'lPblIHYXBwbHlvbg==');

@$core.Deprecated('Use snapStartExceptionDescriptor instead')
const SnapStartException$json = {
  '1': 'SnapStartException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `SnapStartException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List snapStartExceptionDescriptor = $convert.base64Decode(
    'ChJTbmFwU3RhcnRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZRIWCgR0eX'
    'BlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use snapStartNotReadyExceptionDescriptor instead')
const SnapStartNotReadyException$json = {
  '1': 'SnapStartNotReadyException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `SnapStartNotReadyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List snapStartNotReadyExceptionDescriptor =
    $convert.base64Decode(
        'ChpTbmFwU3RhcnROb3RSZWFkeUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dlEhYKBHR5cGUY7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use snapStartResponseDescriptor instead')
const SnapStartResponse$json = {
  '1': 'SnapStartResponse',
  '2': [
    {
      '1': 'applyon',
      '3': 531091671,
      '4': 1,
      '5': 14,
      '6': '.lambda.SnapStartApplyOn',
      '10': 'applyon'
    },
    {
      '1': 'optimizationstatus',
      '3': 155803311,
      '4': 1,
      '5': 14,
      '6': '.lambda.SnapStartOptimizationStatus',
      '10': 'optimizationstatus'
    },
  ],
};

/// Descriptor for `SnapStartResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List snapStartResponseDescriptor = $convert.base64Decode(
    'ChFTbmFwU3RhcnRSZXNwb25zZRI2CgdhcHBseW9uGNehn/0BIAEoDjIYLmxhbWJkYS5TbmFwU3'
    'RhcnRBcHBseU9uUgdhcHBseW9uElYKEm9wdGltaXphdGlvbnN0YXR1cxivvaVKIAEoDjIjLmxh'
    'bWJkYS5TbmFwU3RhcnRPcHRpbWl6YXRpb25TdGF0dXNSEm9wdGltaXphdGlvbnN0YXR1cw==');

@$core.Deprecated('Use snapStartTimeoutExceptionDescriptor instead')
const SnapStartTimeoutException$json = {
  '1': 'SnapStartTimeoutException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `SnapStartTimeoutException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List snapStartTimeoutExceptionDescriptor =
    $convert.base64Decode(
        'ChlTbmFwU3RhcnRUaW1lb3V0RXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'USFgoEdHlwZRjuoNeKASABKAlSBHR5cGU=');

@$core.Deprecated('Use sourceAccessConfigurationDescriptor instead')
const SourceAccessConfiguration$json = {
  '1': 'SourceAccessConfiguration',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.lambda.SourceAccessType',
      '10': 'type'
    },
    {'1': 'uri', '3': 443116318, '4': 1, '5': 9, '10': 'uri'},
  ],
};

/// Descriptor for `SourceAccessConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sourceAccessConfigurationDescriptor =
    $convert.base64Decode(
        'ChlTb3VyY2VBY2Nlc3NDb25maWd1cmF0aW9uEjAKBHR5cGUY7qDXigEgASgOMhgubGFtYmRhLl'
        'NvdXJjZUFjY2Vzc1R5cGVSBHR5cGUSFAoDdXJpGJ7WpdMBIAEoCVIDdXJp');

@$core.Deprecated('Use stepDetailsDescriptor instead')
const StepDetails$json = {
  '1': 'StepDetails',
  '2': [
    {'1': 'attempt', '3': 105696463, '4': 1, '5': 5, '10': 'attempt'},
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '10': 'error'
    },
    {
      '1': 'nextattempttimestamp',
      '3': 217497488,
      '4': 1,
      '5': 9,
      '10': 'nextattempttimestamp'
    },
    {'1': 'result', '3': 273346629, '4': 1, '5': 9, '10': 'result'},
  ],
};

/// Descriptor for `StepDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stepDetailsDescriptor = $convert.base64Decode(
    'CgtTdGVwRGV0YWlscxIbCgdhdHRlbXB0GM+ZszIgASgFUgdhdHRlbXB0Ei0KBWVycm9yGPK5tp'
    'wBIAEoCzITLmxhbWJkYS5FcnJvck9iamVjdFIFZXJyb3ISNQoUbmV4dGF0dGVtcHR0aW1lc3Rh'
    'bXAYkP/aZyABKAlSFG5leHRhdHRlbXB0dGltZXN0YW1wEhoKBnJlc3VsdBjF4KuCASABKAlSBn'
    'Jlc3VsdA==');

@$core.Deprecated('Use stepFailedDetailsDescriptor instead')
const StepFailedDetails$json = {
  '1': 'StepFailedDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
    {
      '1': 'retrydetails',
      '3': 405480610,
      '4': 1,
      '5': 11,
      '6': '.lambda.RetryDetails',
      '10': 'retrydetails'
    },
  ],
};

/// Descriptor for `StepFailedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stepFailedDetailsDescriptor = $convert.base64Decode(
    'ChFTdGVwRmFpbGVkRGV0YWlscxIsCgVlcnJvchjyubacASABKAsyEi5sYW1iZGEuRXZlbnRFcn'
    'JvclIFZXJyb3ISPAoMcmV0cnlkZXRhaWxzGKLJrMEBIAEoCzIULmxhbWJkYS5SZXRyeURldGFp'
    'bHNSDHJldHJ5ZGV0YWlscw==');

@$core.Deprecated('Use stepOptionsDescriptor instead')
const StepOptions$json = {
  '1': 'StepOptions',
  '2': [
    {
      '1': 'nextattemptdelayseconds',
      '3': 502780186,
      '4': 1,
      '5': 5,
      '10': 'nextattemptdelayseconds'
    },
  ],
};

/// Descriptor for `StepOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stepOptionsDescriptor = $convert.base64Decode(
    'CgtTdGVwT3B0aW9ucxI8ChduZXh0YXR0ZW1wdGRlbGF5c2Vjb25kcxiaot/vASABKAVSF25leH'
    'RhdHRlbXB0ZGVsYXlzZWNvbmRz');

@$core.Deprecated('Use stepStartedDetailsDescriptor instead')
const StepStartedDetails$json = {
  '1': 'StepStartedDetails',
};

/// Descriptor for `StepStartedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stepStartedDetailsDescriptor =
    $convert.base64Decode('ChJTdGVwU3RhcnRlZERldGFpbHM=');

@$core.Deprecated('Use stepSucceededDetailsDescriptor instead')
const StepSucceededDetails$json = {
  '1': 'StepSucceededDetails',
  '2': [
    {
      '1': 'result',
      '3': 273346629,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventResult',
      '10': 'result'
    },
    {
      '1': 'retrydetails',
      '3': 405480610,
      '4': 1,
      '5': 11,
      '6': '.lambda.RetryDetails',
      '10': 'retrydetails'
    },
  ],
};

/// Descriptor for `StepSucceededDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stepSucceededDetailsDescriptor = $convert.base64Decode(
    'ChRTdGVwU3VjY2VlZGVkRGV0YWlscxIvCgZyZXN1bHQYxeCrggEgASgLMhMubGFtYmRhLkV2ZW'
    '50UmVzdWx0UgZyZXN1bHQSPAoMcmV0cnlkZXRhaWxzGKLJrMEBIAEoCzIULmxhbWJkYS5SZXRy'
    'eURldGFpbHNSDHJldHJ5ZGV0YWlscw==');

@$core.Deprecated('Use stopDurableExecutionRequestDescriptor instead')
const StopDurableExecutionRequest$json = {
  '1': 'StopDurableExecutionRequest',
  '2': [
    {
      '1': 'durableexecutionarn',
      '3': 269346442,
      '4': 1,
      '5': 9,
      '10': 'durableexecutionarn'
    },
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.ErrorObject',
      '8': {},
      '10': 'error'
    },
  ],
};

/// Descriptor for `StopDurableExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopDurableExecutionRequestDescriptor =
    $convert.base64Decode(
        'ChtTdG9wRHVyYWJsZUV4ZWN1dGlvblJlcXVlc3QSNAoTZHVyYWJsZWV4ZWN1dGlvbmFybhiKzb'
        'eAASABKAlSE2R1cmFibGVleGVjdXRpb25hcm4SMwoFZXJyb3IY8rm2nAEgASgLMhMubGFtYmRh'
        'LkVycm9yT2JqZWN0QgSItRgBUgVlcnJvcg==');

@$core.Deprecated('Use stopDurableExecutionResponseDescriptor instead')
const StopDurableExecutionResponse$json = {
  '1': 'StopDurableExecutionResponse',
  '2': [
    {
      '1': 'stoptimestamp',
      '3': 418051076,
      '4': 1,
      '5': 9,
      '10': 'stoptimestamp'
    },
  ],
};

/// Descriptor for `StopDurableExecutionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopDurableExecutionResponseDescriptor =
    $convert.base64Decode(
        'ChxTdG9wRHVyYWJsZUV4ZWN1dGlvblJlc3BvbnNlEigKDXN0b3B0aW1lc3RhbXAYhOirxwEgAS'
        'gJUg1zdG9wdGltZXN0YW1w');

@$core.Deprecated('Use subnetIPAddressLimitReachedExceptionDescriptor instead')
const SubnetIPAddressLimitReachedException$json = {
  '1': 'SubnetIPAddressLimitReachedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `SubnetIPAddressLimitReachedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subnetIPAddressLimitReachedExceptionDescriptor =
    $convert.base64Decode(
        'CiRTdWJuZXRJUEFkZHJlc3NMaW1pdFJlYWNoZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZRIWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resource', '3': 61921302, '4': 1, '5': 9, '10': 'resource'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.lambda.TagResourceRequest.TagsEntry',
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
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSHQoIcmVzb3VyY2UYlrDDHSABKAlSCHJlc291cmNlEjwKBH'
    'RhZ3MYwcH2tQEgAygLMiQubGFtYmRhLlRhZ1Jlc291cmNlUmVxdWVzdC5UYWdzRW50cnlSBHRh'
    'Z3MaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZT'
    'oCOAE=');

@$core.Deprecated('Use tagsErrorDescriptor instead')
const TagsError$json = {
  '1': 'TagsError',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TagsError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagsErrorDescriptor = $convert.base64Decode(
    'CglUYWdzRXJyb3ISHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcnJvcmNvZGUSGwoHbWVzc2FnZR'
    'iFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use targetTrackingScalingPolicyDescriptor instead')
const TargetTrackingScalingPolicy$json = {
  '1': 'TargetTrackingScalingPolicy',
  '2': [
    {
      '1': 'predefinedmetrictype',
      '3': 334491508,
      '4': 1,
      '5': 14,
      '6': '.lambda.CapacityProviderPredefinedMetricType',
      '10': 'predefinedmetrictype'
    },
    {'1': 'targetvalue', '3': 118247738, '4': 1, '5': 1, '10': 'targetvalue'},
  ],
};

/// Descriptor for `TargetTrackingScalingPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetTrackingScalingPolicyDescriptor = $convert.base64Decode(
    'ChtUYXJnZXRUcmFja2luZ1NjYWxpbmdQb2xpY3kSZAoUcHJlZGVmaW5lZG1ldHJpY3R5cGUY9N'
    '6/nwEgASgOMiwubGFtYmRhLkNhcGFjaXR5UHJvdmlkZXJQcmVkZWZpbmVkTWV0cmljVHlwZVIU'
    'cHJlZGVmaW5lZG1ldHJpY3R5cGUSIwoLdGFyZ2V0dmFsdWUYuqKxOCABKAFSC3RhcmdldHZhbH'
    'Vl');

@$core.Deprecated('Use tenancyConfigDescriptor instead')
const TenancyConfig$json = {
  '1': 'TenancyConfig',
  '2': [
    {
      '1': 'tenantisolationmode',
      '3': 380929873,
      '4': 1,
      '5': 14,
      '6': '.lambda.TenantIsolationMode',
      '10': 'tenantisolationmode'
    },
  ],
};

/// Descriptor for `TenancyConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tenancyConfigDescriptor = $convert.base64Decode(
    'Cg1UZW5hbmN5Q29uZmlnElEKE3RlbmFudGlzb2xhdGlvbm1vZGUY0Y7StQEgASgOMhsubGFtYm'
    'RhLlRlbmFudElzb2xhdGlvbk1vZGVSE3RlbmFudGlzb2xhdGlvbm1vZGU=');

@$core.Deprecated('Use tooManyRequestsExceptionDescriptor instead')
const TooManyRequestsException$json = {
  '1': 'TooManyRequestsException',
  '2': [
    {
      '1': 'reason',
      '3': 20005178,
      '4': 1,
      '5': 14,
      '6': '.lambda.ThrottleReason',
      '10': 'reason'
    },
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'retryafterseconds',
      '3': 436039555,
      '4': 1,
      '5': 9,
      '10': 'retryafterseconds'
    },
  ],
};

/// Descriptor for `TooManyRequestsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyRequestsExceptionDescriptor = $convert.base64Decode(
    'ChhUb29NYW55UmVxdWVzdHNFeGNlcHRpb24SMQoGcmVhc29uGLqCxQkgASgOMhYubGFtYmRhLl'
    'Rocm90dGxlUmVhc29uUgZyZWFzb24SFgoEdHlwZRjuoNeKASABKAlSBHR5cGUSGwoHbWVzc2Fn'
    'ZRjlkcgnIAEoCVIHbWVzc2FnZRIwChFyZXRyeWFmdGVyc2Vjb25kcxiD3/XPASABKAlSEXJldH'
    'J5YWZ0ZXJzZWNvbmRz');

@$core.Deprecated('Use traceHeaderDescriptor instead')
const TraceHeader$json = {
  '1': 'TraceHeader',
  '2': [
    {'1': 'xamzntraceid', '3': 182818336, '4': 1, '5': 9, '10': 'xamzntraceid'},
  ],
};

/// Descriptor for `TraceHeader`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List traceHeaderDescriptor = $convert.base64Decode(
    'CgtUcmFjZUhlYWRlchIlCgx4YW16bnRyYWNlaWQYoKyWVyABKAlSDHhhbXpudHJhY2VpZA==');

@$core.Deprecated('Use tracingConfigDescriptor instead')
const TracingConfig$json = {
  '1': 'TracingConfig',
  '2': [
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.lambda.TracingMode',
      '10': 'mode'
    },
  ],
};

/// Descriptor for `TracingConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tracingConfigDescriptor = $convert.base64Decode(
    'Cg1UcmFjaW5nQ29uZmlnEisKBG1vZGUYs+65mgEgASgOMhMubGFtYmRhLlRyYWNpbmdNb2RlUg'
    'Rtb2Rl');

@$core.Deprecated('Use tracingConfigResponseDescriptor instead')
const TracingConfigResponse$json = {
  '1': 'TracingConfigResponse',
  '2': [
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.lambda.TracingMode',
      '10': 'mode'
    },
  ],
};

/// Descriptor for `TracingConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tracingConfigResponseDescriptor = $convert.base64Decode(
    'ChVUcmFjaW5nQ29uZmlnUmVzcG9uc2USKwoEbW9kZRiz7rmaASABKA4yEy5sYW1iZGEuVHJhY2'
    'luZ01vZGVSBG1vZGU=');

@$core.Deprecated('Use unsupportedMediaTypeExceptionDescriptor instead')
const UnsupportedMediaTypeException$json = {
  '1': 'UnsupportedMediaTypeException',
  '2': [
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedMediaTypeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedMediaTypeExceptionDescriptor =
    $convert.base64Decode(
        'Ch1VbnN1cHBvcnRlZE1lZGlhVHlwZUV4Y2VwdGlvbhIWCgR0eXBlGO6g14oBIAEoCVIEdHlwZR'
        'IbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resource', '3': 61921302, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIdCghyZXNvdXJjZRiWsMMdIAEoCVIIcmVzb3VyY2USHA'
    'oHdGFna2V5cxj8w/OYASADKAlSB3RhZ2tleXM=');

@$core.Deprecated('Use updateAliasRequestDescriptor instead')
const UpdateAliasRequest$json = {
  '1': 'UpdateAliasRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'functionversion',
      '3': 365780244,
      '4': 1,
      '5': 9,
      '10': 'functionversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {
      '1': 'routingconfig',
      '3': 477346140,
      '4': 1,
      '5': 11,
      '6': '.lambda.AliasRoutingConfiguration',
      '10': 'routingconfig'
    },
  ],
};

/// Descriptor for `UpdateAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAliasRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVBbGlhc1JlcXVlc3QSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW'
    '9uEiYKDGZ1bmN0aW9ubmFtZRijiL/fASABKAlSDGZ1bmN0aW9ubmFtZRIsCg9mdW5jdGlvbnZl'
    'cnNpb24YlLq1rgEgASgJUg9mdW5jdGlvbnZlcnNpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZR'
    'IiCgpyZXZpc2lvbmlkGIajnu4BIAEoCVIKcmV2aXNpb25pZBJLCg1yb3V0aW5nY29uZmlnGNzy'
    'zuMBIAEoCzIhLmxhbWJkYS5BbGlhc1JvdXRpbmdDb25maWd1cmF0aW9uUg1yb3V0aW5nY29uZm'
    'ln');

@$core.Deprecated('Use updateCapacityProviderRequestDescriptor instead')
const UpdateCapacityProviderRequest$json = {
  '1': 'UpdateCapacityProviderRequest',
  '2': [
    {
      '1': 'capacityprovidername',
      '3': 466230132,
      '4': 1,
      '5': 9,
      '10': 'capacityprovidername'
    },
    {
      '1': 'capacityproviderscalingconfig',
      '3': 499413336,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderScalingConfig',
      '10': 'capacityproviderscalingconfig'
    },
  ],
};

/// Descriptor for `UpdateCapacityProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCapacityProviderRequestDescriptor = $convert.base64Decode(
    'Ch1VcGRhdGVDYXBhY2l0eVByb3ZpZGVyUmVxdWVzdBI2ChRjYXBhY2l0eXByb3ZpZGVybmFtZR'
    'j0tqjeASABKAlSFGNhcGFjaXR5cHJvdmlkZXJuYW1lEm8KHWNhcGFjaXR5cHJvdmlkZXJzY2Fs'
    'aW5nY29uZmlnGNjike4BIAEoCzIlLmxhbWJkYS5DYXBhY2l0eVByb3ZpZGVyU2NhbGluZ0Nvbm'
    'ZpZ1IdY2FwYWNpdHlwcm92aWRlcnNjYWxpbmdjb25maWc=');

@$core.Deprecated('Use updateCapacityProviderResponseDescriptor instead')
const UpdateCapacityProviderResponse$json = {
  '1': 'UpdateCapacityProviderResponse',
  '2': [
    {
      '1': 'capacityprovider',
      '3': 448074961,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProvider',
      '10': 'capacityprovider'
    },
  ],
};

/// Descriptor for `UpdateCapacityProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCapacityProviderResponseDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVDYXBhY2l0eVByb3ZpZGVyUmVzcG9uc2USSAoQY2FwYWNpdHlwcm92aWRlchjRqd'
        'TVASABKAsyGC5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlclIQY2FwYWNpdHlwcm92aWRlcg==');

@$core.Deprecated('Use updateCodeSigningConfigRequestDescriptor instead')
const UpdateCodeSigningConfigRequest$json = {
  '1': 'UpdateCodeSigningConfigRequest',
  '2': [
    {
      '1': 'allowedpublishers',
      '3': 217933299,
      '4': 1,
      '5': 11,
      '6': '.lambda.AllowedPublishers',
      '10': 'allowedpublishers'
    },
    {
      '1': 'codesigningconfigarn',
      '3': 505282113,
      '4': 1,
      '5': 9,
      '10': 'codesigningconfigarn'
    },
    {
      '1': 'codesigningpolicies',
      '3': 254459202,
      '4': 1,
      '5': 11,
      '6': '.lambda.CodeSigningPolicies',
      '10': 'codesigningpolicies'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `UpdateCodeSigningConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCodeSigningConfigRequestDescriptor = $convert.base64Decode(
    'Ch5VcGRhdGVDb2RlU2lnbmluZ0NvbmZpZ1JlcXVlc3QSSgoRYWxsb3dlZHB1Ymxpc2hlcnMY88'
    'v1ZyABKAsyGS5sYW1iZGEuQWxsb3dlZFB1Ymxpc2hlcnNSEWFsbG93ZWRwdWJsaXNoZXJzEjYK'
    'FGNvZGVzaWduaW5nY29uZmlnYXJuGMH89/ABIAEoCVIUY29kZXNpZ25pbmdjb25maWdhcm4SUA'
    'oTY29kZXNpZ25pbmdwb2xpY2llcxjC+qp5IAEoCzIbLmxhbWJkYS5Db2RlU2lnbmluZ1BvbGlj'
    'aWVzUhNjb2Rlc2lnbmluZ3BvbGljaWVzEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcm'
    'lwdGlvbg==');

@$core.Deprecated('Use updateCodeSigningConfigResponseDescriptor instead')
const UpdateCodeSigningConfigResponse$json = {
  '1': 'UpdateCodeSigningConfigResponse',
  '2': [
    {
      '1': 'codesigningconfig',
      '3': 130458490,
      '4': 1,
      '5': 11,
      '6': '.lambda.CodeSigningConfig',
      '10': 'codesigningconfig'
    },
  ],
};

/// Descriptor for `UpdateCodeSigningConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCodeSigningConfigResponseDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVDb2RlU2lnbmluZ0NvbmZpZ1Jlc3BvbnNlEkoKEWNvZGVzaWduaW5nY29uZmlnGP'
        'rGmj4gASgLMhkubGFtYmRhLkNvZGVTaWduaW5nQ29uZmlnUhFjb2Rlc2lnbmluZ2NvbmZpZw==');

@$core.Deprecated('Use updateEventSourceMappingRequestDescriptor instead')
const UpdateEventSourceMappingRequest$json = {
  '1': 'UpdateEventSourceMappingRequest',
  '2': [
    {
      '1': 'amazonmanagedkafkaeventsourceconfig',
      '3': 60136380,
      '4': 1,
      '5': 11,
      '6': '.lambda.AmazonManagedKafkaEventSourceConfig',
      '10': 'amazonmanagedkafkaeventsourceconfig'
    },
    {'1': 'batchsize', '3': 318039259, '4': 1, '5': 5, '10': 'batchsize'},
    {
      '1': 'bisectbatchonfunctionerror',
      '3': 276143707,
      '4': 1,
      '5': 8,
      '10': 'bisectbatchonfunctionerror'
    },
    {
      '1': 'destinationconfig',
      '3': 184834158,
      '4': 1,
      '5': 11,
      '6': '.lambda.DestinationConfig',
      '10': 'destinationconfig'
    },
    {
      '1': 'documentdbeventsourceconfig',
      '3': 173060622,
      '4': 1,
      '5': 11,
      '6': '.lambda.DocumentDBEventSourceConfig',
      '10': 'documentdbeventsourceconfig'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'filtercriteria',
      '3': 439219323,
      '4': 1,
      '5': 11,
      '6': '.lambda.FilterCriteria',
      '10': 'filtercriteria'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'functionresponsetypes',
      '3': 382292260,
      '4': 3,
      '5': 14,
      '6': '.lambda.FunctionResponseType',
      '10': 'functionresponsetypes'
    },
    {'1': 'kmskeyarn', '3': 117627377, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {
      '1': 'loggingconfig',
      '3': 424359625,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventSourceMappingLoggingConfig',
      '10': 'loggingconfig'
    },
    {
      '1': 'maximumbatchingwindowinseconds',
      '3': 346663320,
      '4': 1,
      '5': 5,
      '10': 'maximumbatchingwindowinseconds'
    },
    {
      '1': 'maximumrecordageinseconds',
      '3': 102344982,
      '4': 1,
      '5': 5,
      '10': 'maximumrecordageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
    {
      '1': 'metricsconfig',
      '3': 412971857,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventSourceMappingMetricsConfig',
      '10': 'metricsconfig'
    },
    {
      '1': 'parallelizationfactor',
      '3': 400337694,
      '4': 1,
      '5': 5,
      '10': 'parallelizationfactor'
    },
    {
      '1': 'provisionedpollerconfig',
      '3': 275676602,
      '4': 1,
      '5': 11,
      '6': '.lambda.ProvisionedPollerConfig',
      '10': 'provisionedpollerconfig'
    },
    {
      '1': 'scalingconfig',
      '3': 392871661,
      '4': 1,
      '5': 11,
      '6': '.lambda.ScalingConfig',
      '10': 'scalingconfig'
    },
    {
      '1': 'selfmanagedkafkaeventsourceconfig',
      '3': 322222578,
      '4': 1,
      '5': 11,
      '6': '.lambda.SelfManagedKafkaEventSourceConfig',
      '10': 'selfmanagedkafkaeventsourceconfig'
    },
    {
      '1': 'sourceaccessconfigurations',
      '3': 371593554,
      '4': 3,
      '5': 11,
      '6': '.lambda.SourceAccessConfiguration',
      '10': 'sourceaccessconfigurations'
    },
    {
      '1': 'tumblingwindowinseconds',
      '3': 372493124,
      '4': 1,
      '5': 5,
      '10': 'tumblingwindowinseconds'
    },
    {'1': 'uuid', '3': 91981875, '4': 1, '5': 9, '10': 'uuid'},
  ],
};

/// Descriptor for `UpdateEventSourceMappingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateEventSourceMappingRequestDescriptor = $convert.base64Decode(
    'Ch9VcGRhdGVFdmVudFNvdXJjZU1hcHBpbmdSZXF1ZXN0EoABCiNhbWF6b25tYW5hZ2Vka2Fma2'
    'FldmVudHNvdXJjZWNvbmZpZxi8t9YcIAEoCzIrLmxhbWJkYS5BbWF6b25NYW5hZ2VkS2Fma2FF'
    'dmVudFNvdXJjZUNvbmZpZ1IjYW1hem9ubWFuYWdlZGthZmthZXZlbnRzb3VyY2Vjb25maWcSIA'
    'oJYmF0Y2hzaXplGNvJ05cBIAEoBVIJYmF0Y2hzaXplEkIKGmJpc2VjdGJhdGNob25mdW5jdGlv'
    'bmVycm9yGNu81oMBIAEoCFIaYmlzZWN0YmF0Y2hvbmZ1bmN0aW9uZXJyb3ISSgoRZGVzdGluYX'
    'Rpb25jb25maWcY7rCRWCABKAsyGS5sYW1iZGEuRGVzdGluYXRpb25Db25maWdSEWRlc3RpbmF0'
    'aW9uY29uZmlnEmgKG2RvY3VtZW50ZGJldmVudHNvdXJjZWNvbmZpZxiO5MJSIAEoCzIjLmxhbW'
    'JkYS5Eb2N1bWVudERCRXZlbnRTb3VyY2VDb25maWdSG2RvY3VtZW50ZGJldmVudHNvdXJjZWNv'
    'bmZpZxIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBJCCg5maWx0ZXJjcml0ZXJpYRj76L'
    'fRASABKAsyFi5sYW1iZGEuRmlsdGVyQ3JpdGVyaWFSDmZpbHRlcmNyaXRlcmlhEiYKDGZ1bmN0'
    'aW9ubmFtZRijiL/fASABKAlSDGZ1bmN0aW9ubmFtZRJWChVmdW5jdGlvbnJlc3BvbnNldHlwZX'
    'MYpKKltgEgAygOMhwubGFtYmRhLkZ1bmN0aW9uUmVzcG9uc2VUeXBlUhVmdW5jdGlvbnJlc3Bv'
    'bnNldHlwZXMSHwoJa21za2V5YXJuGPGzizggASgJUglrbXNrZXlhcm4SUQoNbG9nZ2luZ2Nvbm'
    'ZpZxjJ7azKASABKAsyJy5sYW1iZGEuRXZlbnRTb3VyY2VNYXBwaW5nTG9nZ2luZ0NvbmZpZ1IN'
    'bG9nZ2luZ2NvbmZpZxJKCh5tYXhpbXVtYmF0Y2hpbmd3aW5kb3dpbnNlY29uZHMYmNOmpQEgAS'
    'gFUh5tYXhpbXVtYmF0Y2hpbmd3aW5kb3dpbnNlY29uZHMSPwoZbWF4aW11bXJlY29yZGFnZWlu'
    'c2Vjb25kcxiW0uYwIAEoBVIZbWF4aW11bXJlY29yZGFnZWluc2Vjb25kcxI1ChRtYXhpbXVtcm'
    'V0cnlhdHRlbXB0cxjAqLk1IAEoBVIUbWF4aW11bXJldHJ5YXR0ZW1wdHMSUQoNbWV0cmljc2Nv'
    'bmZpZxjR5vXEASABKAsyJy5sYW1iZGEuRXZlbnRTb3VyY2VNYXBwaW5nTWV0cmljc0NvbmZpZ1'
    'INbWV0cmljc2NvbmZpZxI4ChVwYXJhbGxlbGl6YXRpb25mYWN0b3IYntbyvgEgASgFUhVwYXJh'
    'bGxlbGl6YXRpb25mYWN0b3ISXQoXcHJvdmlzaW9uZWRwb2xsZXJjb25maWcYuvu5gwEgASgLMh'
    '8ubGFtYmRhLlByb3Zpc2lvbmVkUG9sbGVyQ29uZmlnUhdwcm92aXNpb25lZHBvbGxlcmNvbmZp'
    'ZxI/Cg1zY2FsaW5nY29uZmlnGO39qrsBIAEoCzIVLmxhbWJkYS5TY2FsaW5nQ29uZmlnUg1zY2'
    'FsaW5nY29uZmlnEnsKIXNlbGZtYW5hZ2Vka2Fma2FldmVudHNvdXJjZWNvbmZpZxjy89KZASAB'
    'KAsyKS5sYW1iZGEuU2VsZk1hbmFnZWRLYWZrYUV2ZW50U291cmNlQ29uZmlnUiFzZWxmbWFuYW'
    'dlZGthZmthZXZlbnRzb3VyY2Vjb25maWcSZQoac291cmNlYWNjZXNzY29uZmlndXJhdGlvbnMY'
    '0qKYsQEgAygLMiEubGFtYmRhLlNvdXJjZUFjY2Vzc0NvbmZpZ3VyYXRpb25SGnNvdXJjZWFjY2'
    'Vzc2NvbmZpZ3VyYXRpb25zEjwKF3R1bWJsaW5nd2luZG93aW5zZWNvbmRzGMSWz7EBIAEoBVIX'
    'dHVtYmxpbmd3aW5kb3dpbnNlY29uZHMSFQoEdXVpZBizkO4rIAEoCVIEdXVpZA==');

@$core.Deprecated('Use updateFunctionCodeRequestDescriptor instead')
const UpdateFunctionCodeRequest$json = {
  '1': 'UpdateFunctionCodeRequest',
  '2': [
    {
      '1': 'architectures',
      '3': 530490948,
      '4': 3,
      '5': 14,
      '6': '.lambda.Architecture',
      '10': 'architectures'
    },
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'imageuri', '3': 412238461, '4': 1, '5': 9, '10': 'imageuri'},
    {'1': 'publish', '3': 207759785, '4': 1, '5': 8, '10': 'publish'},
    {
      '1': 'publishto',
      '3': 524127682,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionVersionLatestPublished',
      '10': 'publishto'
    },
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 's3bucket', '3': 114031434, '4': 1, '5': 9, '10': 's3bucket'},
    {'1': 's3key', '3': 490298907, '4': 1, '5': 9, '10': 's3key'},
    {
      '1': 's3objectversion',
      '3': 194809669,
      '4': 1,
      '5': 9,
      '10': 's3objectversion'
    },
    {
      '1': 'sourcekmskeyarn',
      '3': 203651164,
      '4': 1,
      '5': 9,
      '10': 'sourcekmskeyarn'
    },
    {'1': 'zipfile', '3': 2519299, '4': 1, '5': 12, '10': 'zipfile'},
  ],
};

/// Descriptor for `UpdateFunctionCodeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFunctionCodeRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVGdW5jdGlvbkNvZGVSZXF1ZXN0Ej4KDWFyY2hpdGVjdHVyZXMYxMz6/AEgAygOMh'
    'QubGFtYmRhLkFyY2hpdGVjdHVyZVINYXJjaGl0ZWN0dXJlcxIZCgZkcnlydW4YuN/7KyABKAhS'
    'BmRyeXJ1bhImCgxmdW5jdGlvbm5hbWUYo4i/3wEgASgJUgxmdW5jdGlvbm5hbWUSHgoIaW1hZ2'
    'V1cmkY/YTJxAEgASgJUghpbWFnZXVyaRIbCgdwdWJsaXNoGKnTiGMgASgIUgdwdWJsaXNoEkgK'
    'CXB1Ymxpc2h0bxjCm/b5ASABKA4yJi5sYW1iZGEuRnVuY3Rpb25WZXJzaW9uTGF0ZXN0UHVibG'
    'lzaGVkUglwdWJsaXNodG8SIgoKcmV2aXNpb25pZBiGo57uASABKAlSCnJldmlzaW9uaWQSHQoI'
    'czNidWNrZXQYyvavNiABKAlSCHMzYnVja2V0EhgKBXMza2V5GJu85ekBIAEoCVIFczNrZXkSKw'
    'oPczNvYmplY3R2ZXJzaW9uGMWe8lwgASgJUg9zM29iamVjdHZlcnNpb24SKwoPc291cmNla21z'
    'a2V5YXJuGNzwjWEgASgJUg9zb3VyY2VrbXNrZXlhcm4SGwoHemlwZmlsZRiD4pkBIAEoDFIHem'
    'lwZmlsZQ==');

@$core.Deprecated('Use updateFunctionConfigurationRequestDescriptor instead')
const UpdateFunctionConfigurationRequest$json = {
  '1': 'UpdateFunctionConfigurationRequest',
  '2': [
    {
      '1': 'capacityproviderconfig',
      '3': 52030623,
      '4': 1,
      '5': 11,
      '6': '.lambda.CapacityProviderConfig',
      '10': 'capacityproviderconfig'
    },
    {
      '1': 'deadletterconfig',
      '3': 79786642,
      '4': 1,
      '5': 11,
      '6': '.lambda.DeadLetterConfig',
      '10': 'deadletterconfig'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'durableconfig',
      '3': 206326279,
      '4': 1,
      '5': 11,
      '6': '.lambda.DurableConfig',
      '10': 'durableconfig'
    },
    {
      '1': 'environment',
      '3': 119823003,
      '4': 1,
      '5': 11,
      '6': '.lambda.Environment',
      '10': 'environment'
    },
    {
      '1': 'ephemeralstorage',
      '3': 365965382,
      '4': 1,
      '5': 11,
      '6': '.lambda.EphemeralStorage',
      '10': 'ephemeralstorage'
    },
    {
      '1': 'filesystemconfigs',
      '3': 490453750,
      '4': 3,
      '5': 11,
      '6': '.lambda.FileSystemConfig',
      '10': 'filesystemconfigs'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {'1': 'handler', '3': 81160724, '4': 1, '5': 9, '10': 'handler'},
    {
      '1': 'imageconfig',
      '3': 281970485,
      '4': 1,
      '5': 11,
      '6': '.lambda.ImageConfig',
      '10': 'imageconfig'
    },
    {'1': 'kmskeyarn', '3': 117627377, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {'1': 'layers', '3': 478144896, '4': 3, '5': 9, '10': 'layers'},
    {
      '1': 'loggingconfig',
      '3': 424359625,
      '4': 1,
      '5': 11,
      '6': '.lambda.LoggingConfig',
      '10': 'loggingconfig'
    },
    {'1': 'memorysize', '3': 55523120, '4': 1, '5': 5, '10': 'memorysize'},
    {'1': 'revisionid', '3': 499618182, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'role', '3': 271285818, '4': 1, '5': 9, '10': 'role'},
    {
      '1': 'runtime',
      '3': 359311308,
      '4': 1,
      '5': 14,
      '6': '.lambda.Runtime',
      '10': 'runtime'
    },
    {
      '1': 'snapstart',
      '3': 283273032,
      '4': 1,
      '5': 11,
      '6': '.lambda.SnapStart',
      '10': 'snapstart'
    },
    {'1': 'timeout', '3': 47808041, '4': 1, '5': 5, '10': 'timeout'},
    {
      '1': 'tracingconfig',
      '3': 19554860,
      '4': 1,
      '5': 11,
      '6': '.lambda.TracingConfig',
      '10': 'tracingconfig'
    },
    {
      '1': 'vpcconfig',
      '3': 194980743,
      '4': 1,
      '5': 11,
      '6': '.lambda.VpcConfig',
      '10': 'vpcconfig'
    },
  ],
};

/// Descriptor for `UpdateFunctionConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFunctionConfigurationRequestDescriptor = $convert.base64Decode(
    'CiJVcGRhdGVGdW5jdGlvbkNvbmZpZ3VyYXRpb25SZXF1ZXN0ElkKFmNhcGFjaXR5cHJvdmlkZX'
    'Jjb25maWcYn9nnGCABKAsyHi5sYW1iZGEuQ2FwYWNpdHlQcm92aWRlckNvbmZpZ1IWY2FwYWNp'
    'dHlwcm92aWRlcmNvbmZpZxJHChBkZWFkbGV0dGVyY29uZmlnGJLlhSYgASgLMhgubGFtYmRhLk'
    'RlYWRMZXR0ZXJDb25maWdSEGRlYWRsZXR0ZXJjb25maWcSIwoLZGVzY3JpcHRpb24YivT5NiAB'
    'KAlSC2Rlc2NyaXB0aW9uEj4KDWR1cmFibGVjb25maWcYh5SxYiABKAsyFS5sYW1iZGEuRHVyYW'
    'JsZUNvbmZpZ1INZHVyYWJsZWNvbmZpZxI4CgtlbnZpcm9ubWVudBibtZE5IAEoCzITLmxhbWJk'
    'YS5FbnZpcm9ubWVudFILZW52aXJvbm1lbnQSSAoQZXBoZW1lcmFsc3RvcmFnZRjG4MCuASABKA'
    'syGC5sYW1iZGEuRXBoZW1lcmFsU3RvcmFnZVIQZXBoZW1lcmFsc3RvcmFnZRJKChFmaWxlc3lz'
    'dGVtY29uZmlncxj29e7pASADKAsyGC5sYW1iZGEuRmlsZVN5c3RlbUNvbmZpZ1IRZmlsZXN5c3'
    'RlbWNvbmZpZ3MSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb25uYW1lEhsKB2hh'
    'bmRsZXIYlNTZJiABKAlSB2hhbmRsZXISOQoLaW1hZ2Vjb25maWcYtY66hgEgASgLMhMubGFtYm'
    'RhLkltYWdlQ29uZmlnUgtpbWFnZWNvbmZpZxIfCglrbXNrZXlhcm4Y8bOLOCABKAlSCWttc2tl'
    'eWFybhIaCgZsYXllcnMYgNP/4wEgAygJUgZsYXllcnMSPwoNbG9nZ2luZ2NvbmZpZxjJ7azKAS'
    'ABKAsyFS5sYW1iZGEuTG9nZ2luZ0NvbmZpZ1INbG9nZ2luZ2NvbmZpZxIhCgptZW1vcnlzaXpl'
    'GLDuvBogASgFUgptZW1vcnlzaXplEiIKCnJldmlzaW9uaWQYhqOe7gEgASgJUgpyZXZpc2lvbm'
    'lkEhYKBHJvbGUYuvytgQEgASgJUgRyb2xlEi0KB3J1bnRpbWUYzM+qqwEgASgOMg8ubGFtYmRh'
    'LlJ1bnRpbWVSB3J1bnRpbWUSMwoJc25hcHN0YXJ0GMjOiYcBIAEoCzIRLmxhbWJkYS5TbmFwU3'
    'RhcnRSCXNuYXBzdGFydBIbCgd0aW1lb3V0GKn85RYgASgFUgd0aW1lb3V0Ej4KDXRyYWNpbmdj'
    'b25maWcYrMSpCSABKAsyFS5sYW1iZGEuVHJhY2luZ0NvbmZpZ1INdHJhY2luZ2NvbmZpZxIyCg'
    'l2cGNjb25maWcYh9f8XCABKAsyES5sYW1iZGEuVnBjQ29uZmlnUgl2cGNjb25maWc=');

@$core
    .Deprecated('Use updateFunctionEventInvokeConfigRequestDescriptor instead')
const UpdateFunctionEventInvokeConfigRequest$json = {
  '1': 'UpdateFunctionEventInvokeConfigRequest',
  '2': [
    {
      '1': 'destinationconfig',
      '3': 184834158,
      '4': 1,
      '5': 11,
      '6': '.lambda.DestinationConfig',
      '10': 'destinationconfig'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'maximumeventageinseconds',
      '3': 393041563,
      '4': 1,
      '5': 5,
      '10': 'maximumeventageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `UpdateFunctionEventInvokeConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFunctionEventInvokeConfigRequestDescriptor =
    $convert.base64Decode(
        'CiZVcGRhdGVGdW5jdGlvbkV2ZW50SW52b2tlQ29uZmlnUmVxdWVzdBJKChFkZXN0aW5hdGlvbm'
        'NvbmZpZxjusJFYIAEoCzIZLmxhbWJkYS5EZXN0aW5hdGlvbkNvbmZpZ1IRZGVzdGluYXRpb25j'
        'b25maWcSJgoMZnVuY3Rpb25uYW1lGKOIv98BIAEoCVIMZnVuY3Rpb25uYW1lEj4KGG1heGltdW'
        '1ldmVudGFnZWluc2Vjb25kcxibrbW7ASABKAVSGG1heGltdW1ldmVudGFnZWluc2Vjb25kcxI1'
        'ChRtYXhpbXVtcmV0cnlhdHRlbXB0cxjAqLk1IAEoBVIUbWF4aW11bXJldHJ5YXR0ZW1wdHMSIA'
        'oJcXVhbGlmaWVyGOC1kfsBIAEoCVIJcXVhbGlmaWVy');

@$core.Deprecated('Use updateFunctionUrlConfigRequestDescriptor instead')
const UpdateFunctionUrlConfigRequest$json = {
  '1': 'UpdateFunctionUrlConfigRequest',
  '2': [
    {
      '1': 'authtype',
      '3': 477704248,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionUrlAuthType',
      '10': 'authtype'
    },
    {
      '1': 'cors',
      '3': 260753653,
      '4': 1,
      '5': 11,
      '6': '.lambda.Cors',
      '10': 'cors'
    },
    {'1': 'functionname', '3': 468698147, '4': 1, '5': 9, '10': 'functionname'},
    {
      '1': 'invokemode',
      '3': 414956667,
      '4': 1,
      '5': 14,
      '6': '.lambda.InvokeMode',
      '10': 'invokemode'
    },
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `UpdateFunctionUrlConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFunctionUrlConfigRequestDescriptor = $convert.base64Decode(
    'Ch5VcGRhdGVGdW5jdGlvblVybENvbmZpZ1JlcXVlc3QSOwoIYXV0aHR5cGUYuODk4wEgASgOMh'
    'subGFtYmRhLkZ1bmN0aW9uVXJsQXV0aFR5cGVSCGF1dGh0eXBlEiMKBGNvcnMY9ZGrfCABKAsy'
    'DC5sYW1iZGEuQ29yc1IEY29ycxImCgxmdW5jdGlvbm5hbWUYo4i/3wEgASgJUgxmdW5jdGlvbm'
    '5hbWUSNgoKaW52b2tlbW9kZRj7+O7FASABKA4yEi5sYW1iZGEuSW52b2tlTW9kZVIKaW52b2tl'
    'bW9kZRIgCglxdWFsaWZpZXIY4LWR+wEgASgJUglxdWFsaWZpZXI=');

@$core.Deprecated('Use updateFunctionUrlConfigResponseDescriptor instead')
const UpdateFunctionUrlConfigResponse$json = {
  '1': 'UpdateFunctionUrlConfigResponse',
  '2': [
    {
      '1': 'authtype',
      '3': 477704248,
      '4': 1,
      '5': 14,
      '6': '.lambda.FunctionUrlAuthType',
      '10': 'authtype'
    },
    {
      '1': 'cors',
      '3': 260753653,
      '4': 1,
      '5': 11,
      '6': '.lambda.Cors',
      '10': 'cors'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'functionarn', '3': 381920497, '4': 1, '5': 9, '10': 'functionarn'},
    {'1': 'functionurl', '3': 449381947, '4': 1, '5': 9, '10': 'functionurl'},
    {
      '1': 'invokemode',
      '3': 414956667,
      '4': 1,
      '5': 14,
      '6': '.lambda.InvokeMode',
      '10': 'invokemode'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
  ],
};

/// Descriptor for `UpdateFunctionUrlConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateFunctionUrlConfigResponseDescriptor = $convert.base64Decode(
    'Ch9VcGRhdGVGdW5jdGlvblVybENvbmZpZ1Jlc3BvbnNlEjsKCGF1dGh0eXBlGLjg5OMBIAEoDj'
    'IbLmxhbWJkYS5GdW5jdGlvblVybEF1dGhUeXBlUghhdXRodHlwZRIjCgRjb3JzGPWRq3wgASgL'
    'MgwubGFtYmRhLkNvcnNSBGNvcnMSJQoMY3JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlvbn'
    'RpbWUSJAoLZnVuY3Rpb25hcm4Y8cmOtgEgASgJUgtmdW5jdGlvbmFybhIkCgtmdW5jdGlvbnVy'
    'bBi7jKTWASABKAlSC2Z1bmN0aW9udXJsEjYKCmludm9rZW1vZGUY+/juxQEgASgOMhIubGFtYm'
    'RhLkludm9rZU1vZGVSCmludm9rZW1vZGUSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQ'
    'bGFzdG1vZGlmaWVkdGltZQ==');

@$core.Deprecated('Use vpcConfigDescriptor instead')
const VpcConfig$json = {
  '1': 'VpcConfig',
  '2': [
    {
      '1': 'ipv6allowedfordualstack',
      '3': 138147060,
      '4': 1,
      '5': 8,
      '10': 'ipv6allowedfordualstack'
    },
    {
      '1': 'securitygroupids',
      '3': 13337861,
      '4': 3,
      '5': 9,
      '10': 'securitygroupids'
    },
    {'1': 'subnetids', '3': 266219411, '4': 3, '5': 9, '10': 'subnetids'},
  ],
};

/// Descriptor for `VpcConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcConfigDescriptor = $convert.base64Decode(
    'CglWcGNDb25maWcSOwoXaXB2NmFsbG93ZWRmb3JkdWFsc3RhY2sY9OnvQSABKAhSF2lwdjZhbG'
    'xvd2VkZm9yZHVhbHN0YWNrEi0KEHNlY3VyaXR5Z3JvdXBpZHMYhYquBiADKAlSEHNlY3VyaXR5'
    'Z3JvdXBpZHMSHwoJc3VibmV0aWRzGJPf+H4gAygJUglzdWJuZXRpZHM=');

@$core.Deprecated('Use vpcConfigResponseDescriptor instead')
const VpcConfigResponse$json = {
  '1': 'VpcConfigResponse',
  '2': [
    {
      '1': 'ipv6allowedfordualstack',
      '3': 138147060,
      '4': 1,
      '5': 8,
      '10': 'ipv6allowedfordualstack'
    },
    {
      '1': 'securitygroupids',
      '3': 13337861,
      '4': 3,
      '5': 9,
      '10': 'securitygroupids'
    },
    {'1': 'subnetids', '3': 266219411, '4': 3, '5': 9, '10': 'subnetids'},
    {'1': 'vpcid', '3': 412355958, '4': 1, '5': 9, '10': 'vpcid'},
  ],
};

/// Descriptor for `VpcConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vpcConfigResponseDescriptor = $convert.base64Decode(
    'ChFWcGNDb25maWdSZXNwb25zZRI7ChdpcHY2YWxsb3dlZGZvcmR1YWxzdGFjaxj06e9BIAEoCF'
    'IXaXB2NmFsbG93ZWRmb3JkdWFsc3RhY2sSLQoQc2VjdXJpdHlncm91cGlkcxiFiq4GIAMoCVIQ'
    'c2VjdXJpdHlncm91cGlkcxIfCglzdWJuZXRpZHMYk9/4fiADKAlSCXN1Ym5ldGlkcxIYCgV2cG'
    'NpZBj2mtDEASABKAlSBXZwY2lk');

@$core.Deprecated('Use waitCancelledDetailsDescriptor instead')
const WaitCancelledDetails$json = {
  '1': 'WaitCancelledDetails',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.lambda.EventError',
      '10': 'error'
    },
  ],
};

/// Descriptor for `WaitCancelledDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List waitCancelledDetailsDescriptor = $convert.base64Decode(
    'ChRXYWl0Q2FuY2VsbGVkRGV0YWlscxIsCgVlcnJvchjyubacASABKAsyEi5sYW1iZGEuRXZlbn'
    'RFcnJvclIFZXJyb3I=');

@$core.Deprecated('Use waitDetailsDescriptor instead')
const WaitDetails$json = {
  '1': 'WaitDetails',
  '2': [
    {
      '1': 'scheduledendtimestamp',
      '3': 325003612,
      '4': 1,
      '5': 9,
      '10': 'scheduledendtimestamp'
    },
  ],
};

/// Descriptor for `WaitDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List waitDetailsDescriptor = $convert.base64Decode(
    'CgtXYWl0RGV0YWlscxI4ChVzY2hlZHVsZWRlbmR0aW1lc3RhbXAY3NL8mgEgASgJUhVzY2hlZH'
    'VsZWRlbmR0aW1lc3RhbXA=');

@$core.Deprecated('Use waitOptionsDescriptor instead')
const WaitOptions$json = {
  '1': 'WaitOptions',
  '2': [
    {'1': 'waitseconds', '3': 305214910, '4': 1, '5': 5, '10': 'waitseconds'},
  ],
};

/// Descriptor for `WaitOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List waitOptionsDescriptor = $convert.base64Decode(
    'CgtXYWl0T3B0aW9ucxIkCgt3YWl0c2Vjb25kcxi+68SRASABKAVSC3dhaXRzZWNvbmRz');

@$core.Deprecated('Use waitStartedDetailsDescriptor instead')
const WaitStartedDetails$json = {
  '1': 'WaitStartedDetails',
  '2': [
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {
      '1': 'scheduledendtimestamp',
      '3': 325003612,
      '4': 1,
      '5': 9,
      '10': 'scheduledendtimestamp'
    },
  ],
};

/// Descriptor for `WaitStartedDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List waitStartedDetailsDescriptor = $convert.base64Decode(
    'ChJXYWl0U3RhcnRlZERldGFpbHMSHgoIZHVyYXRpb24YrpKdpgEgASgFUghkdXJhdGlvbhI4Ch'
    'VzY2hlZHVsZWRlbmR0aW1lc3RhbXAY3NL8mgEgASgJUhVzY2hlZHVsZWRlbmR0aW1lc3RhbXA=');

@$core.Deprecated('Use waitSucceededDetailsDescriptor instead')
const WaitSucceededDetails$json = {
  '1': 'WaitSucceededDetails',
  '2': [
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
  ],
};

/// Descriptor for `WaitSucceededDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List waitSucceededDetailsDescriptor = $convert.base64Decode(
    'ChRXYWl0U3VjY2VlZGVkRGV0YWlscxIeCghkdXJhdGlvbhiukp2mASABKAVSCGR1cmF0aW9u');
