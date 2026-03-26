// This is a generated file - do not edit.
//
// Generated from lambda.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ApplicationLogLevel extends $pb.ProtobufEnum {
  static const ApplicationLogLevel APPLICATION_LOG_LEVEL_WARN =
      ApplicationLogLevel._(
          0, _omitEnumNames ? '' : 'APPLICATION_LOG_LEVEL_WARN');
  static const ApplicationLogLevel APPLICATION_LOG_LEVEL_FATAL =
      ApplicationLogLevel._(
          1, _omitEnumNames ? '' : 'APPLICATION_LOG_LEVEL_FATAL');
  static const ApplicationLogLevel APPLICATION_LOG_LEVEL_ERROR =
      ApplicationLogLevel._(
          2, _omitEnumNames ? '' : 'APPLICATION_LOG_LEVEL_ERROR');
  static const ApplicationLogLevel APPLICATION_LOG_LEVEL_TRACE =
      ApplicationLogLevel._(
          3, _omitEnumNames ? '' : 'APPLICATION_LOG_LEVEL_TRACE');
  static const ApplicationLogLevel APPLICATION_LOG_LEVEL_INFO =
      ApplicationLogLevel._(
          4, _omitEnumNames ? '' : 'APPLICATION_LOG_LEVEL_INFO');
  static const ApplicationLogLevel APPLICATION_LOG_LEVEL_DEBUG =
      ApplicationLogLevel._(
          5, _omitEnumNames ? '' : 'APPLICATION_LOG_LEVEL_DEBUG');

  static const $core.List<ApplicationLogLevel> values = <ApplicationLogLevel>[
    APPLICATION_LOG_LEVEL_WARN,
    APPLICATION_LOG_LEVEL_FATAL,
    APPLICATION_LOG_LEVEL_ERROR,
    APPLICATION_LOG_LEVEL_TRACE,
    APPLICATION_LOG_LEVEL_INFO,
    APPLICATION_LOG_LEVEL_DEBUG,
  ];

  static final $core.List<ApplicationLogLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ApplicationLogLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ApplicationLogLevel._(super.value, super.name);
}

class Architecture extends $pb.ProtobufEnum {
  static const Architecture ARCHITECTURE_ARM64 =
      Architecture._(0, _omitEnumNames ? '' : 'ARCHITECTURE_ARM64');
  static const Architecture ARCHITECTURE_X86_64 =
      Architecture._(1, _omitEnumNames ? '' : 'ARCHITECTURE_X86_64');

  static const $core.List<Architecture> values = <Architecture>[
    ARCHITECTURE_ARM64,
    ARCHITECTURE_X86_64,
  ];

  static final $core.List<Architecture?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static Architecture? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Architecture._(super.value, super.name);
}

class CapacityProviderPredefinedMetricType extends $pb.ProtobufEnum {
  static const CapacityProviderPredefinedMetricType
      CAPACITY_PROVIDER_PREDEFINED_METRIC_TYPE_LAMBDACAPACITYPROVIDERAVERAGECPUUTILIZATION =
      CapacityProviderPredefinedMetricType._(
          0,
          _omitEnumNames
              ? ''
              : 'CAPACITY_PROVIDER_PREDEFINED_METRIC_TYPE_LAMBDACAPACITYPROVIDERAVERAGECPUUTILIZATION');

  static const $core.List<CapacityProviderPredefinedMetricType> values =
      <CapacityProviderPredefinedMetricType>[
    CAPACITY_PROVIDER_PREDEFINED_METRIC_TYPE_LAMBDACAPACITYPROVIDERAVERAGECPUUTILIZATION,
  ];

  static final $core.List<CapacityProviderPredefinedMetricType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static CapacityProviderPredefinedMetricType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CapacityProviderPredefinedMetricType._(super.value, super.name);
}

class CapacityProviderScalingMode extends $pb.ProtobufEnum {
  static const CapacityProviderScalingMode
      CAPACITY_PROVIDER_SCALING_MODE_MANUAL = CapacityProviderScalingMode._(
          0, _omitEnumNames ? '' : 'CAPACITY_PROVIDER_SCALING_MODE_MANUAL');
  static const CapacityProviderScalingMode CAPACITY_PROVIDER_SCALING_MODE_AUTO =
      CapacityProviderScalingMode._(
          1, _omitEnumNames ? '' : 'CAPACITY_PROVIDER_SCALING_MODE_AUTO');

  static const $core.List<CapacityProviderScalingMode> values =
      <CapacityProviderScalingMode>[
    CAPACITY_PROVIDER_SCALING_MODE_MANUAL,
    CAPACITY_PROVIDER_SCALING_MODE_AUTO,
  ];

  static final $core.List<CapacityProviderScalingMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CapacityProviderScalingMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CapacityProviderScalingMode._(super.value, super.name);
}

class CapacityProviderState extends $pb.ProtobufEnum {
  static const CapacityProviderState CAPACITY_PROVIDER_STATE_ACTIVE =
      CapacityProviderState._(
          0, _omitEnumNames ? '' : 'CAPACITY_PROVIDER_STATE_ACTIVE');
  static const CapacityProviderState CAPACITY_PROVIDER_STATE_FAILED =
      CapacityProviderState._(
          1, _omitEnumNames ? '' : 'CAPACITY_PROVIDER_STATE_FAILED');
  static const CapacityProviderState CAPACITY_PROVIDER_STATE_DELETING =
      CapacityProviderState._(
          2, _omitEnumNames ? '' : 'CAPACITY_PROVIDER_STATE_DELETING');
  static const CapacityProviderState CAPACITY_PROVIDER_STATE_PENDING =
      CapacityProviderState._(
          3, _omitEnumNames ? '' : 'CAPACITY_PROVIDER_STATE_PENDING');

  static const $core.List<CapacityProviderState> values =
      <CapacityProviderState>[
    CAPACITY_PROVIDER_STATE_ACTIVE,
    CAPACITY_PROVIDER_STATE_FAILED,
    CAPACITY_PROVIDER_STATE_DELETING,
    CAPACITY_PROVIDER_STATE_PENDING,
  ];

  static final $core.List<CapacityProviderState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static CapacityProviderState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CapacityProviderState._(super.value, super.name);
}

class CodeSigningPolicy extends $pb.ProtobufEnum {
  static const CodeSigningPolicy CODE_SIGNING_POLICY_WARN =
      CodeSigningPolicy._(0, _omitEnumNames ? '' : 'CODE_SIGNING_POLICY_WARN');
  static const CodeSigningPolicy CODE_SIGNING_POLICY_ENFORCE =
      CodeSigningPolicy._(
          1, _omitEnumNames ? '' : 'CODE_SIGNING_POLICY_ENFORCE');

  static const $core.List<CodeSigningPolicy> values = <CodeSigningPolicy>[
    CODE_SIGNING_POLICY_WARN,
    CODE_SIGNING_POLICY_ENFORCE,
  ];

  static final $core.List<CodeSigningPolicy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CodeSigningPolicy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CodeSigningPolicy._(super.value, super.name);
}

class EndPointType extends $pb.ProtobufEnum {
  static const EndPointType END_POINT_TYPE_KAFKA_BOOTSTRAP_SERVERS =
      EndPointType._(
          0, _omitEnumNames ? '' : 'END_POINT_TYPE_KAFKA_BOOTSTRAP_SERVERS');

  static const $core.List<EndPointType> values = <EndPointType>[
    END_POINT_TYPE_KAFKA_BOOTSTRAP_SERVERS,
  ];

  static final $core.List<EndPointType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static EndPointType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EndPointType._(super.value, super.name);
}

class EventSourceMappingMetric extends $pb.ProtobufEnum {
  static const EventSourceMappingMetric EVENT_SOURCE_MAPPING_METRIC_EVENTCOUNT =
      EventSourceMappingMetric._(
          0, _omitEnumNames ? '' : 'EVENT_SOURCE_MAPPING_METRIC_EVENTCOUNT');
  static const EventSourceMappingMetric
      EVENT_SOURCE_MAPPING_METRIC_KAFKAMETRICS = EventSourceMappingMetric._(
          1, _omitEnumNames ? '' : 'EVENT_SOURCE_MAPPING_METRIC_KAFKAMETRICS');
  static const EventSourceMappingMetric EVENT_SOURCE_MAPPING_METRIC_ERRORCOUNT =
      EventSourceMappingMetric._(
          2, _omitEnumNames ? '' : 'EVENT_SOURCE_MAPPING_METRIC_ERRORCOUNT');

  static const $core.List<EventSourceMappingMetric> values =
      <EventSourceMappingMetric>[
    EVENT_SOURCE_MAPPING_METRIC_EVENTCOUNT,
    EVENT_SOURCE_MAPPING_METRIC_KAFKAMETRICS,
    EVENT_SOURCE_MAPPING_METRIC_ERRORCOUNT,
  ];

  static final $core.List<EventSourceMappingMetric?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EventSourceMappingMetric? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventSourceMappingMetric._(super.value, super.name);
}

class EventSourceMappingSystemLogLevel extends $pb.ProtobufEnum {
  static const EventSourceMappingSystemLogLevel
      EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_WARN =
      EventSourceMappingSystemLogLevel._(0,
          _omitEnumNames ? '' : 'EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_WARN');
  static const EventSourceMappingSystemLogLevel
      EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_INFO =
      EventSourceMappingSystemLogLevel._(1,
          _omitEnumNames ? '' : 'EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_INFO');
  static const EventSourceMappingSystemLogLevel
      EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_DEBUG =
      EventSourceMappingSystemLogLevel._(2,
          _omitEnumNames ? '' : 'EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_DEBUG');

  static const $core.List<EventSourceMappingSystemLogLevel> values =
      <EventSourceMappingSystemLogLevel>[
    EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_WARN,
    EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_INFO,
    EVENT_SOURCE_MAPPING_SYSTEM_LOG_LEVEL_DEBUG,
  ];

  static final $core.List<EventSourceMappingSystemLogLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EventSourceMappingSystemLogLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventSourceMappingSystemLogLevel._(super.value, super.name);
}

class EventSourcePosition extends $pb.ProtobufEnum {
  static const EventSourcePosition EVENT_SOURCE_POSITION_AT_TIMESTAMP =
      EventSourcePosition._(
          0, _omitEnumNames ? '' : 'EVENT_SOURCE_POSITION_AT_TIMESTAMP');
  static const EventSourcePosition EVENT_SOURCE_POSITION_TRIM_HORIZON =
      EventSourcePosition._(
          1, _omitEnumNames ? '' : 'EVENT_SOURCE_POSITION_TRIM_HORIZON');
  static const EventSourcePosition EVENT_SOURCE_POSITION_LATEST =
      EventSourcePosition._(
          2, _omitEnumNames ? '' : 'EVENT_SOURCE_POSITION_LATEST');

  static const $core.List<EventSourcePosition> values = <EventSourcePosition>[
    EVENT_SOURCE_POSITION_AT_TIMESTAMP,
    EVENT_SOURCE_POSITION_TRIM_HORIZON,
    EVENT_SOURCE_POSITION_LATEST,
  ];

  static final $core.List<EventSourcePosition?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EventSourcePosition? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventSourcePosition._(super.value, super.name);
}

class EventType extends $pb.ProtobufEnum {
  static const EventType EVENT_TYPE_CHAINEDINVOKESTARTED =
      EventType._(0, _omitEnumNames ? '' : 'EVENT_TYPE_CHAINEDINVOKESTARTED');
  static const EventType EVENT_TYPE_INVOCATIONCOMPLETED =
      EventType._(1, _omitEnumNames ? '' : 'EVENT_TYPE_INVOCATIONCOMPLETED');
  static const EventType EVENT_TYPE_EXECUTIONTIMEDOUT =
      EventType._(2, _omitEnumNames ? '' : 'EVENT_TYPE_EXECUTIONTIMEDOUT');
  static const EventType EVENT_TYPE_CALLBACKSUCCEEDED =
      EventType._(3, _omitEnumNames ? '' : 'EVENT_TYPE_CALLBACKSUCCEEDED');
  static const EventType EVENT_TYPE_EXECUTIONSTOPPED =
      EventType._(4, _omitEnumNames ? '' : 'EVENT_TYPE_EXECUTIONSTOPPED');
  static const EventType EVENT_TYPE_WAITCANCELLED =
      EventType._(5, _omitEnumNames ? '' : 'EVENT_TYPE_WAITCANCELLED');
  static const EventType EVENT_TYPE_STEPFAILED =
      EventType._(6, _omitEnumNames ? '' : 'EVENT_TYPE_STEPFAILED');
  static const EventType EVENT_TYPE_EXECUTIONFAILED =
      EventType._(7, _omitEnumNames ? '' : 'EVENT_TYPE_EXECUTIONFAILED');
  static const EventType EVENT_TYPE_CHAINEDINVOKESUCCEEDED =
      EventType._(8, _omitEnumNames ? '' : 'EVENT_TYPE_CHAINEDINVOKESUCCEEDED');
  static const EventType EVENT_TYPE_CALLBACKSTARTED =
      EventType._(9, _omitEnumNames ? '' : 'EVENT_TYPE_CALLBACKSTARTED');
  static const EventType EVENT_TYPE_CONTEXTFAILED =
      EventType._(10, _omitEnumNames ? '' : 'EVENT_TYPE_CONTEXTFAILED');
  static const EventType EVENT_TYPE_WAITSTARTED =
      EventType._(11, _omitEnumNames ? '' : 'EVENT_TYPE_WAITSTARTED');
  static const EventType EVENT_TYPE_EXECUTIONSTARTED =
      EventType._(12, _omitEnumNames ? '' : 'EVENT_TYPE_EXECUTIONSTARTED');
  static const EventType EVENT_TYPE_CALLBACKTIMEDOUT =
      EventType._(13, _omitEnumNames ? '' : 'EVENT_TYPE_CALLBACKTIMEDOUT');
  static const EventType EVENT_TYPE_CHAINEDINVOKETIMEDOUT =
      EventType._(14, _omitEnumNames ? '' : 'EVENT_TYPE_CHAINEDINVOKETIMEDOUT');
  static const EventType EVENT_TYPE_STEPSTARTED =
      EventType._(15, _omitEnumNames ? '' : 'EVENT_TYPE_STEPSTARTED');
  static const EventType EVENT_TYPE_CALLBACKFAILED =
      EventType._(16, _omitEnumNames ? '' : 'EVENT_TYPE_CALLBACKFAILED');
  static const EventType EVENT_TYPE_CONTEXTSTARTED =
      EventType._(17, _omitEnumNames ? '' : 'EVENT_TYPE_CONTEXTSTARTED');
  static const EventType EVENT_TYPE_CHAINEDINVOKEFAILED =
      EventType._(18, _omitEnumNames ? '' : 'EVENT_TYPE_CHAINEDINVOKEFAILED');
  static const EventType EVENT_TYPE_EXECUTIONSUCCEEDED =
      EventType._(19, _omitEnumNames ? '' : 'EVENT_TYPE_EXECUTIONSUCCEEDED');
  static const EventType EVENT_TYPE_CHAINEDINVOKESTOPPED =
      EventType._(20, _omitEnumNames ? '' : 'EVENT_TYPE_CHAINEDINVOKESTOPPED');
  static const EventType EVENT_TYPE_STEPSUCCEEDED =
      EventType._(21, _omitEnumNames ? '' : 'EVENT_TYPE_STEPSUCCEEDED');
  static const EventType EVENT_TYPE_WAITSUCCEEDED =
      EventType._(22, _omitEnumNames ? '' : 'EVENT_TYPE_WAITSUCCEEDED');
  static const EventType EVENT_TYPE_CONTEXTSUCCEEDED =
      EventType._(23, _omitEnumNames ? '' : 'EVENT_TYPE_CONTEXTSUCCEEDED');

  static const $core.List<EventType> values = <EventType>[
    EVENT_TYPE_CHAINEDINVOKESTARTED,
    EVENT_TYPE_INVOCATIONCOMPLETED,
    EVENT_TYPE_EXECUTIONTIMEDOUT,
    EVENT_TYPE_CALLBACKSUCCEEDED,
    EVENT_TYPE_EXECUTIONSTOPPED,
    EVENT_TYPE_WAITCANCELLED,
    EVENT_TYPE_STEPFAILED,
    EVENT_TYPE_EXECUTIONFAILED,
    EVENT_TYPE_CHAINEDINVOKESUCCEEDED,
    EVENT_TYPE_CALLBACKSTARTED,
    EVENT_TYPE_CONTEXTFAILED,
    EVENT_TYPE_WAITSTARTED,
    EVENT_TYPE_EXECUTIONSTARTED,
    EVENT_TYPE_CALLBACKTIMEDOUT,
    EVENT_TYPE_CHAINEDINVOKETIMEDOUT,
    EVENT_TYPE_STEPSTARTED,
    EVENT_TYPE_CALLBACKFAILED,
    EVENT_TYPE_CONTEXTSTARTED,
    EVENT_TYPE_CHAINEDINVOKEFAILED,
    EVENT_TYPE_EXECUTIONSUCCEEDED,
    EVENT_TYPE_CHAINEDINVOKESTOPPED,
    EVENT_TYPE_STEPSUCCEEDED,
    EVENT_TYPE_WAITSUCCEEDED,
    EVENT_TYPE_CONTEXTSUCCEEDED,
  ];

  static final $core.List<EventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 23);
  static EventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventType._(super.value, super.name);
}

class ExecutionStatus extends $pb.ProtobufEnum {
  static const ExecutionStatus EXECUTION_STATUS_RUNNING =
      ExecutionStatus._(0, _omitEnumNames ? '' : 'EXECUTION_STATUS_RUNNING');
  static const ExecutionStatus EXECUTION_STATUS_TIMED_OUT =
      ExecutionStatus._(1, _omitEnumNames ? '' : 'EXECUTION_STATUS_TIMED_OUT');
  static const ExecutionStatus EXECUTION_STATUS_STOPPED =
      ExecutionStatus._(2, _omitEnumNames ? '' : 'EXECUTION_STATUS_STOPPED');
  static const ExecutionStatus EXECUTION_STATUS_SUCCEEDED =
      ExecutionStatus._(3, _omitEnumNames ? '' : 'EXECUTION_STATUS_SUCCEEDED');
  static const ExecutionStatus EXECUTION_STATUS_FAILED =
      ExecutionStatus._(4, _omitEnumNames ? '' : 'EXECUTION_STATUS_FAILED');

  static const $core.List<ExecutionStatus> values = <ExecutionStatus>[
    EXECUTION_STATUS_RUNNING,
    EXECUTION_STATUS_TIMED_OUT,
    EXECUTION_STATUS_STOPPED,
    EXECUTION_STATUS_SUCCEEDED,
    EXECUTION_STATUS_FAILED,
  ];

  static final $core.List<ExecutionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ExecutionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutionStatus._(super.value, super.name);
}

class FullDocument extends $pb.ProtobufEnum {
  static const FullDocument FULL_DOCUMENT_UPDATELOOKUP =
      FullDocument._(0, _omitEnumNames ? '' : 'FULL_DOCUMENT_UPDATELOOKUP');
  static const FullDocument FULL_DOCUMENT_DEFAULT =
      FullDocument._(1, _omitEnumNames ? '' : 'FULL_DOCUMENT_DEFAULT');

  static const $core.List<FullDocument> values = <FullDocument>[
    FULL_DOCUMENT_UPDATELOOKUP,
    FULL_DOCUMENT_DEFAULT,
  ];

  static final $core.List<FullDocument?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FullDocument? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FullDocument._(super.value, super.name);
}

class FunctionResponseType extends $pb.ProtobufEnum {
  static const FunctionResponseType
      FUNCTION_RESPONSE_TYPE_REPORTBATCHITEMFAILURES = FunctionResponseType._(
          0,
          _omitEnumNames
              ? ''
              : 'FUNCTION_RESPONSE_TYPE_REPORTBATCHITEMFAILURES');

  static const $core.List<FunctionResponseType> values = <FunctionResponseType>[
    FUNCTION_RESPONSE_TYPE_REPORTBATCHITEMFAILURES,
  ];

  static final $core.List<FunctionResponseType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static FunctionResponseType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FunctionResponseType._(super.value, super.name);
}

class FunctionUrlAuthType extends $pb.ProtobufEnum {
  static const FunctionUrlAuthType FUNCTION_URL_AUTH_TYPE_NONE =
      FunctionUrlAuthType._(
          0, _omitEnumNames ? '' : 'FUNCTION_URL_AUTH_TYPE_NONE');
  static const FunctionUrlAuthType FUNCTION_URL_AUTH_TYPE_AWS_IAM =
      FunctionUrlAuthType._(
          1, _omitEnumNames ? '' : 'FUNCTION_URL_AUTH_TYPE_AWS_IAM');

  static const $core.List<FunctionUrlAuthType> values = <FunctionUrlAuthType>[
    FUNCTION_URL_AUTH_TYPE_NONE,
    FUNCTION_URL_AUTH_TYPE_AWS_IAM,
  ];

  static final $core.List<FunctionUrlAuthType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FunctionUrlAuthType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FunctionUrlAuthType._(super.value, super.name);
}

class FunctionVersion extends $pb.ProtobufEnum {
  static const FunctionVersion FUNCTION_VERSION_ALL =
      FunctionVersion._(0, _omitEnumNames ? '' : 'FUNCTION_VERSION_ALL');

  static const $core.List<FunctionVersion> values = <FunctionVersion>[
    FUNCTION_VERSION_ALL,
  ];

  static final $core.List<FunctionVersion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static FunctionVersion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FunctionVersion._(super.value, super.name);
}

class FunctionVersionLatestPublished extends $pb.ProtobufEnum {
  static const FunctionVersionLatestPublished
      FUNCTION_VERSION_LATEST_PUBLISHED_LATEST_PUBLISHED =
      FunctionVersionLatestPublished._(
          0,
          _omitEnumNames
              ? ''
              : 'FUNCTION_VERSION_LATEST_PUBLISHED_LATEST_PUBLISHED');

  static const $core.List<FunctionVersionLatestPublished> values =
      <FunctionVersionLatestPublished>[
    FUNCTION_VERSION_LATEST_PUBLISHED_LATEST_PUBLISHED,
  ];

  static final $core.List<FunctionVersionLatestPublished?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static FunctionVersionLatestPublished? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FunctionVersionLatestPublished._(super.value, super.name);
}

class InvocationType extends $pb.ProtobufEnum {
  static const InvocationType INVOCATION_TYPE_EVENT =
      InvocationType._(0, _omitEnumNames ? '' : 'INVOCATION_TYPE_EVENT');
  static const InvocationType INVOCATION_TYPE_REQUESTRESPONSE =
      InvocationType._(
          1, _omitEnumNames ? '' : 'INVOCATION_TYPE_REQUESTRESPONSE');
  static const InvocationType INVOCATION_TYPE_DRYRUN =
      InvocationType._(2, _omitEnumNames ? '' : 'INVOCATION_TYPE_DRYRUN');

  static const $core.List<InvocationType> values = <InvocationType>[
    INVOCATION_TYPE_EVENT,
    INVOCATION_TYPE_REQUESTRESPONSE,
    INVOCATION_TYPE_DRYRUN,
  ];

  static final $core.List<InvocationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static InvocationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InvocationType._(super.value, super.name);
}

class InvokeMode extends $pb.ProtobufEnum {
  static const InvokeMode INVOKE_MODE_RESPONSE_STREAM =
      InvokeMode._(0, _omitEnumNames ? '' : 'INVOKE_MODE_RESPONSE_STREAM');
  static const InvokeMode INVOKE_MODE_BUFFERED =
      InvokeMode._(1, _omitEnumNames ? '' : 'INVOKE_MODE_BUFFERED');

  static const $core.List<InvokeMode> values = <InvokeMode>[
    INVOKE_MODE_RESPONSE_STREAM,
    INVOKE_MODE_BUFFERED,
  ];

  static final $core.List<InvokeMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InvokeMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InvokeMode._(super.value, super.name);
}

class KafkaSchemaRegistryAuthType extends $pb.ProtobufEnum {
  static const KafkaSchemaRegistryAuthType
      KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_CLIENT_CERTIFICATE_TLS_AUTH =
      KafkaSchemaRegistryAuthType._(
          0,
          _omitEnumNames
              ? ''
              : 'KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_CLIENT_CERTIFICATE_TLS_AUTH');
  static const KafkaSchemaRegistryAuthType
      KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_SERVER_ROOT_CA_CERTIFICATE =
      KafkaSchemaRegistryAuthType._(
          1,
          _omitEnumNames
              ? ''
              : 'KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_SERVER_ROOT_CA_CERTIFICATE');
  static const KafkaSchemaRegistryAuthType
      KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_BASIC_AUTH =
      KafkaSchemaRegistryAuthType._(2,
          _omitEnumNames ? '' : 'KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_BASIC_AUTH');

  static const $core.List<KafkaSchemaRegistryAuthType> values =
      <KafkaSchemaRegistryAuthType>[
    KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_CLIENT_CERTIFICATE_TLS_AUTH,
    KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_SERVER_ROOT_CA_CERTIFICATE,
    KAFKA_SCHEMA_REGISTRY_AUTH_TYPE_BASIC_AUTH,
  ];

  static final $core.List<KafkaSchemaRegistryAuthType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static KafkaSchemaRegistryAuthType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KafkaSchemaRegistryAuthType._(super.value, super.name);
}

class KafkaSchemaValidationAttribute extends $pb.ProtobufEnum {
  static const KafkaSchemaValidationAttribute
      KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_VALUE =
      KafkaSchemaValidationAttribute._(
          0, _omitEnumNames ? '' : 'KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_VALUE');
  static const KafkaSchemaValidationAttribute
      KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_KEY = KafkaSchemaValidationAttribute._(
          1, _omitEnumNames ? '' : 'KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_KEY');

  static const $core.List<KafkaSchemaValidationAttribute> values =
      <KafkaSchemaValidationAttribute>[
    KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_VALUE,
    KAFKA_SCHEMA_VALIDATION_ATTRIBUTE_KEY,
  ];

  static final $core.List<KafkaSchemaValidationAttribute?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static KafkaSchemaValidationAttribute? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KafkaSchemaValidationAttribute._(super.value, super.name);
}

class LastUpdateStatus extends $pb.ProtobufEnum {
  static const LastUpdateStatus LAST_UPDATE_STATUS_SUCCESSFUL =
      LastUpdateStatus._(
          0, _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_SUCCESSFUL');
  static const LastUpdateStatus LAST_UPDATE_STATUS_FAILED =
      LastUpdateStatus._(1, _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_FAILED');
  static const LastUpdateStatus LAST_UPDATE_STATUS_INPROGRESS =
      LastUpdateStatus._(
          2, _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_INPROGRESS');

  static const $core.List<LastUpdateStatus> values = <LastUpdateStatus>[
    LAST_UPDATE_STATUS_SUCCESSFUL,
    LAST_UPDATE_STATUS_FAILED,
    LAST_UPDATE_STATUS_INPROGRESS,
  ];

  static final $core.List<LastUpdateStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static LastUpdateStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LastUpdateStatus._(super.value, super.name);
}

class LastUpdateStatusReasonCode extends $pb.ProtobufEnum {
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_ENILIMITEXCEEDED =
      LastUpdateStatusReasonCode._(
          0,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_ENILIMITEXCEEDED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_SUBNETOUTOFIPADDRESSES =
      LastUpdateStatusReasonCode._(
          1,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_SUBNETOUTOFIPADDRESSES');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_KMSKEYNOTFOUND =
      LastUpdateStatusReasonCode._(
          2,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_KMSKEYNOTFOUND');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITTIMEOUT =
      LastUpdateStatusReasonCode._(
          3,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITTIMEOUT');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INVALIDZIPFILEEXCEPTION =
      LastUpdateStatusReasonCode._(
          4,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDZIPFILEEXCEPTION');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL =
      LastUpdateStatusReasonCode._(
          5,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR =
      LastUpdateStatusReasonCode._(
          6,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_EC2REQUESTLIMITEXCEEDED =
      LastUpdateStatusReasonCode._(
          7,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_EC2REQUESTLIMITEXCEEDED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INVALIDIMAGE =
      LastUpdateStatusReasonCode._(8,
          _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDIMAGE');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED =
      LastUpdateStatusReasonCode._(
          9,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED =
      LastUpdateStatusReasonCode._(
          10,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INVALIDRUNTIME =
      LastUpdateStatusReasonCode._(
          11,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDRUNTIME');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_DISABLEDKMSKEY =
      LastUpdateStatusReasonCode._(
          12,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_DISABLEDKMSKEY');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_VCPULIMITEXCEEDED =
      LastUpdateStatusReasonCode._(
          13,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_VCPULIMITEXCEEDED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS =
      LastUpdateStatusReasonCode._(
          14,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROR =
      LastUpdateStatusReasonCode._(15,
          _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROR');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INVALIDCONFIGURATION =
      LastUpdateStatusReasonCode._(
          16,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDCONFIGURATION');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INVALIDSUBNET =
      LastUpdateStatusReasonCode._(17,
          _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDSUBNET');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INVALIDSTATEKMSKEY =
      LastUpdateStatusReasonCode._(
          18,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDSTATEKMSKEY');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_IMAGEACCESSDENIED =
      LastUpdateStatusReasonCode._(
          19,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_IMAGEACCESSDENIED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_IMAGEDELETED =
      LastUpdateStatusReasonCode._(20,
          _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_REASON_CODE_IMAGEDELETED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTTIMEOUT =
      LastUpdateStatusReasonCode._(
          21,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTTIMEOUT');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTFAILURE =
      LastUpdateStatusReasonCode._(
          22,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTFAILURE');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_EFSIOERROR = LastUpdateStatusReasonCode._(
          23,
          _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_REASON_CODE_EFSIOERROR');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR =
      LastUpdateStatusReasonCode._(
          24,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INVALIDSECURITYGROUP =
      LastUpdateStatusReasonCode._(
          25,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_INVALIDSECURITYGROUP');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT =
      LastUpdateStatusReasonCode._(
          26,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTCONNECTIVITYERROR =
      LastUpdateStatusReasonCode._(
          27,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTCONNECTIVITYERROR');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY =
      LastUpdateStatusReasonCode._(
          28,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTROLEPERMISSIONS =
      LastUpdateStatusReasonCode._(
          29,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTROLEPERMISSIONS');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_KMSKEYACCESSDENIED =
      LastUpdateStatusReasonCode._(
          30,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_KMSKEYACCESSDENIED');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INTERNALERROR =
      LastUpdateStatusReasonCode._(31,
          _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_REASON_CODE_INTERNALERROR');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTCAPACITY =
      LastUpdateStatusReasonCode._(
          32,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTCAPACITY');
  static const LastUpdateStatusReasonCode
      LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED =
      LastUpdateStatusReasonCode._(
          33,
          _omitEnumNames
              ? ''
              : 'LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED');

  static const $core.List<LastUpdateStatusReasonCode> values =
      <LastUpdateStatusReasonCode>[
    LAST_UPDATE_STATUS_REASON_CODE_ENILIMITEXCEEDED,
    LAST_UPDATE_STATUS_REASON_CODE_SUBNETOUTOFIPADDRESSES,
    LAST_UPDATE_STATUS_REASON_CODE_KMSKEYNOTFOUND,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITTIMEOUT,
    LAST_UPDATE_STATUS_REASON_CODE_INVALIDZIPFILEEXCEPTION,
    LAST_UPDATE_STATUS_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR,
    LAST_UPDATE_STATUS_REASON_CODE_EC2REQUESTLIMITEXCEEDED,
    LAST_UPDATE_STATUS_REASON_CODE_INVALIDIMAGE,
    LAST_UPDATE_STATUS_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED,
    LAST_UPDATE_STATUS_REASON_CODE_INVALIDRUNTIME,
    LAST_UPDATE_STATUS_REASON_CODE_DISABLEDKMSKEY,
    LAST_UPDATE_STATUS_REASON_CODE_VCPULIMITEXCEEDED,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROR,
    LAST_UPDATE_STATUS_REASON_CODE_INVALIDCONFIGURATION,
    LAST_UPDATE_STATUS_REASON_CODE_INVALIDSUBNET,
    LAST_UPDATE_STATUS_REASON_CODE_INVALIDSTATEKMSKEY,
    LAST_UPDATE_STATUS_REASON_CODE_IMAGEACCESSDENIED,
    LAST_UPDATE_STATUS_REASON_CODE_IMAGEDELETED,
    LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTTIMEOUT,
    LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTFAILURE,
    LAST_UPDATE_STATUS_REASON_CODE_EFSIOERROR,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR,
    LAST_UPDATE_STATUS_REASON_CODE_INVALIDSECURITYGROUP,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT,
    LAST_UPDATE_STATUS_REASON_CODE_EFSMOUNTCONNECTIVITYERROR,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY,
    LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTROLEPERMISSIONS,
    LAST_UPDATE_STATUS_REASON_CODE_KMSKEYACCESSDENIED,
    LAST_UPDATE_STATUS_REASON_CODE_INTERNALERROR,
    LAST_UPDATE_STATUS_REASON_CODE_INSUFFICIENTCAPACITY,
    LAST_UPDATE_STATUS_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED,
  ];

  static final $core.List<LastUpdateStatusReasonCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 33);
  static LastUpdateStatusReasonCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LastUpdateStatusReasonCode._(super.value, super.name);
}

class LogFormat extends $pb.ProtobufEnum {
  static const LogFormat LOG_FORMAT_TEXT =
      LogFormat._(0, _omitEnumNames ? '' : 'LOG_FORMAT_TEXT');
  static const LogFormat LOG_FORMAT_JSON =
      LogFormat._(1, _omitEnumNames ? '' : 'LOG_FORMAT_JSON');

  static const $core.List<LogFormat> values = <LogFormat>[
    LOG_FORMAT_TEXT,
    LOG_FORMAT_JSON,
  ];

  static final $core.List<LogFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static LogFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogFormat._(super.value, super.name);
}

class LogType extends $pb.ProtobufEnum {
  static const LogType LOG_TYPE_NONE =
      LogType._(0, _omitEnumNames ? '' : 'LOG_TYPE_NONE');
  static const LogType LOG_TYPE_TAIL =
      LogType._(1, _omitEnumNames ? '' : 'LOG_TYPE_TAIL');

  static const $core.List<LogType> values = <LogType>[
    LOG_TYPE_NONE,
    LOG_TYPE_TAIL,
  ];

  static final $core.List<LogType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static LogType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogType._(super.value, super.name);
}

class OperationAction extends $pb.ProtobufEnum {
  static const OperationAction OPERATION_ACTION_SUCCEED =
      OperationAction._(0, _omitEnumNames ? '' : 'OPERATION_ACTION_SUCCEED');
  static const OperationAction OPERATION_ACTION_RETRY =
      OperationAction._(1, _omitEnumNames ? '' : 'OPERATION_ACTION_RETRY');
  static const OperationAction OPERATION_ACTION_START =
      OperationAction._(2, _omitEnumNames ? '' : 'OPERATION_ACTION_START');
  static const OperationAction OPERATION_ACTION_CANCEL =
      OperationAction._(3, _omitEnumNames ? '' : 'OPERATION_ACTION_CANCEL');
  static const OperationAction OPERATION_ACTION_FAIL =
      OperationAction._(4, _omitEnumNames ? '' : 'OPERATION_ACTION_FAIL');

  static const $core.List<OperationAction> values = <OperationAction>[
    OPERATION_ACTION_SUCCEED,
    OPERATION_ACTION_RETRY,
    OPERATION_ACTION_START,
    OPERATION_ACTION_CANCEL,
    OPERATION_ACTION_FAIL,
  ];

  static final $core.List<OperationAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static OperationAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OperationAction._(super.value, super.name);
}

class OperationStatus extends $pb.ProtobufEnum {
  static const OperationStatus OPERATION_STATUS_PENDING =
      OperationStatus._(0, _omitEnumNames ? '' : 'OPERATION_STATUS_PENDING');
  static const OperationStatus OPERATION_STATUS_TIMED_OUT =
      OperationStatus._(1, _omitEnumNames ? '' : 'OPERATION_STATUS_TIMED_OUT');
  static const OperationStatus OPERATION_STATUS_STOPPED =
      OperationStatus._(2, _omitEnumNames ? '' : 'OPERATION_STATUS_STOPPED');
  static const OperationStatus OPERATION_STATUS_STARTED =
      OperationStatus._(3, _omitEnumNames ? '' : 'OPERATION_STATUS_STARTED');
  static const OperationStatus OPERATION_STATUS_READY =
      OperationStatus._(4, _omitEnumNames ? '' : 'OPERATION_STATUS_READY');
  static const OperationStatus OPERATION_STATUS_SUCCEEDED =
      OperationStatus._(5, _omitEnumNames ? '' : 'OPERATION_STATUS_SUCCEEDED');
  static const OperationStatus OPERATION_STATUS_CANCELLED =
      OperationStatus._(6, _omitEnumNames ? '' : 'OPERATION_STATUS_CANCELLED');
  static const OperationStatus OPERATION_STATUS_FAILED =
      OperationStatus._(7, _omitEnumNames ? '' : 'OPERATION_STATUS_FAILED');

  static const $core.List<OperationStatus> values = <OperationStatus>[
    OPERATION_STATUS_PENDING,
    OPERATION_STATUS_TIMED_OUT,
    OPERATION_STATUS_STOPPED,
    OPERATION_STATUS_STARTED,
    OPERATION_STATUS_READY,
    OPERATION_STATUS_SUCCEEDED,
    OPERATION_STATUS_CANCELLED,
    OPERATION_STATUS_FAILED,
  ];

  static final $core.List<OperationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static OperationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OperationStatus._(super.value, super.name);
}

class OperationType extends $pb.ProtobufEnum {
  static const OperationType OPERATION_TYPE_CALLBACK =
      OperationType._(0, _omitEnumNames ? '' : 'OPERATION_TYPE_CALLBACK');
  static const OperationType OPERATION_TYPE_STEP =
      OperationType._(1, _omitEnumNames ? '' : 'OPERATION_TYPE_STEP');
  static const OperationType OPERATION_TYPE_CHAINED_INVOKE =
      OperationType._(2, _omitEnumNames ? '' : 'OPERATION_TYPE_CHAINED_INVOKE');
  static const OperationType OPERATION_TYPE_EXECUTION =
      OperationType._(3, _omitEnumNames ? '' : 'OPERATION_TYPE_EXECUTION');
  static const OperationType OPERATION_TYPE_WAIT =
      OperationType._(4, _omitEnumNames ? '' : 'OPERATION_TYPE_WAIT');
  static const OperationType OPERATION_TYPE_CONTEXT =
      OperationType._(5, _omitEnumNames ? '' : 'OPERATION_TYPE_CONTEXT');

  static const $core.List<OperationType> values = <OperationType>[
    OPERATION_TYPE_CALLBACK,
    OPERATION_TYPE_STEP,
    OPERATION_TYPE_CHAINED_INVOKE,
    OPERATION_TYPE_EXECUTION,
    OPERATION_TYPE_WAIT,
    OPERATION_TYPE_CONTEXT,
  ];

  static final $core.List<OperationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static OperationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OperationType._(super.value, super.name);
}

class PackageType extends $pb.ProtobufEnum {
  static const PackageType PACKAGE_TYPE_IMAGE =
      PackageType._(0, _omitEnumNames ? '' : 'PACKAGE_TYPE_IMAGE');
  static const PackageType PACKAGE_TYPE_ZIP =
      PackageType._(1, _omitEnumNames ? '' : 'PACKAGE_TYPE_ZIP');

  static const $core.List<PackageType> values = <PackageType>[
    PACKAGE_TYPE_IMAGE,
    PACKAGE_TYPE_ZIP,
  ];

  static final $core.List<PackageType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PackageType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PackageType._(super.value, super.name);
}

class ProvisionedConcurrencyStatusEnum extends $pb.ProtobufEnum {
  static const ProvisionedConcurrencyStatusEnum
      PROVISIONED_CONCURRENCY_STATUS_ENUM_IN_PROGRESS =
      ProvisionedConcurrencyStatusEnum._(
          0,
          _omitEnumNames
              ? ''
              : 'PROVISIONED_CONCURRENCY_STATUS_ENUM_IN_PROGRESS');
  static const ProvisionedConcurrencyStatusEnum
      PROVISIONED_CONCURRENCY_STATUS_ENUM_READY =
      ProvisionedConcurrencyStatusEnum._(
          1, _omitEnumNames ? '' : 'PROVISIONED_CONCURRENCY_STATUS_ENUM_READY');
  static const ProvisionedConcurrencyStatusEnum
      PROVISIONED_CONCURRENCY_STATUS_ENUM_FAILED =
      ProvisionedConcurrencyStatusEnum._(2,
          _omitEnumNames ? '' : 'PROVISIONED_CONCURRENCY_STATUS_ENUM_FAILED');

  static const $core.List<ProvisionedConcurrencyStatusEnum> values =
      <ProvisionedConcurrencyStatusEnum>[
    PROVISIONED_CONCURRENCY_STATUS_ENUM_IN_PROGRESS,
    PROVISIONED_CONCURRENCY_STATUS_ENUM_READY,
    PROVISIONED_CONCURRENCY_STATUS_ENUM_FAILED,
  ];

  static final $core.List<ProvisionedConcurrencyStatusEnum?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ProvisionedConcurrencyStatusEnum? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ProvisionedConcurrencyStatusEnum._(super.value, super.name);
}

class RecursiveLoop extends $pb.ProtobufEnum {
  static const RecursiveLoop RECURSIVE_LOOP_ALLOW =
      RecursiveLoop._(0, _omitEnumNames ? '' : 'RECURSIVE_LOOP_ALLOW');
  static const RecursiveLoop RECURSIVE_LOOP_TERMINATE =
      RecursiveLoop._(1, _omitEnumNames ? '' : 'RECURSIVE_LOOP_TERMINATE');

  static const $core.List<RecursiveLoop> values = <RecursiveLoop>[
    RECURSIVE_LOOP_ALLOW,
    RECURSIVE_LOOP_TERMINATE,
  ];

  static final $core.List<RecursiveLoop?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RecursiveLoop? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecursiveLoop._(super.value, super.name);
}

class ResponseStreamingInvocationType extends $pb.ProtobufEnum {
  static const ResponseStreamingInvocationType
      RESPONSE_STREAMING_INVOCATION_TYPE_REQUESTRESPONSE =
      ResponseStreamingInvocationType._(
          0,
          _omitEnumNames
              ? ''
              : 'RESPONSE_STREAMING_INVOCATION_TYPE_REQUESTRESPONSE');
  static const ResponseStreamingInvocationType
      RESPONSE_STREAMING_INVOCATION_TYPE_DRYRUN =
      ResponseStreamingInvocationType._(
          1, _omitEnumNames ? '' : 'RESPONSE_STREAMING_INVOCATION_TYPE_DRYRUN');

  static const $core.List<ResponseStreamingInvocationType> values =
      <ResponseStreamingInvocationType>[
    RESPONSE_STREAMING_INVOCATION_TYPE_REQUESTRESPONSE,
    RESPONSE_STREAMING_INVOCATION_TYPE_DRYRUN,
  ];

  static final $core.List<ResponseStreamingInvocationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ResponseStreamingInvocationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResponseStreamingInvocationType._(super.value, super.name);
}

class Runtime extends $pb.ProtobufEnum {
  static const Runtime RUNTIME_DOTNETCORE31 =
      Runtime._(0, _omitEnumNames ? '' : 'RUNTIME_DOTNETCORE31');
  static const Runtime RUNTIME_NODEJS16X =
      Runtime._(1, _omitEnumNames ? '' : 'RUNTIME_NODEJS16X');
  static const Runtime RUNTIME_RUBY25 =
      Runtime._(2, _omitEnumNames ? '' : 'RUNTIME_RUBY25');
  static const Runtime RUNTIME_NODEJS24X =
      Runtime._(3, _omitEnumNames ? '' : 'RUNTIME_NODEJS24X');
  static const Runtime RUNTIME_PYTHON39 =
      Runtime._(4, _omitEnumNames ? '' : 'RUNTIME_PYTHON39');
  static const Runtime RUNTIME_DOTNET6 =
      Runtime._(5, _omitEnumNames ? '' : 'RUNTIME_DOTNET6');
  static const Runtime RUNTIME_RUBY32 =
      Runtime._(6, _omitEnumNames ? '' : 'RUNTIME_RUBY32');
  static const Runtime RUNTIME_GO1X =
      Runtime._(7, _omitEnumNames ? '' : 'RUNTIME_GO1X');
  static const Runtime RUNTIME_PYTHON36 =
      Runtime._(8, _omitEnumNames ? '' : 'RUNTIME_PYTHON36');
  static const Runtime RUNTIME_PROVIDEDAL2023 =
      Runtime._(9, _omitEnumNames ? '' : 'RUNTIME_PROVIDEDAL2023');
  static const Runtime RUNTIME_NODEJS22X =
      Runtime._(10, _omitEnumNames ? '' : 'RUNTIME_NODEJS22X');
  static const Runtime RUNTIME_PYTHON312 =
      Runtime._(11, _omitEnumNames ? '' : 'RUNTIME_PYTHON312');
  static const Runtime RUNTIME_DOTNETCORE21 =
      Runtime._(12, _omitEnumNames ? '' : 'RUNTIME_DOTNETCORE21');
  static const Runtime RUNTIME_NODEJS810 =
      Runtime._(13, _omitEnumNames ? '' : 'RUNTIME_NODEJS810');
  static const Runtime RUNTIME_RUBY27 =
      Runtime._(14, _omitEnumNames ? '' : 'RUNTIME_RUBY27');
  static const Runtime RUNTIME_JAVA21 =
      Runtime._(15, _omitEnumNames ? '' : 'RUNTIME_JAVA21');
  static const Runtime RUNTIME_NODEJS14X =
      Runtime._(16, _omitEnumNames ? '' : 'RUNTIME_NODEJS14X');
  static const Runtime RUNTIME_PROVIDED =
      Runtime._(17, _omitEnumNames ? '' : 'RUNTIME_PROVIDED');
  static const Runtime RUNTIME_NODEJS610 =
      Runtime._(18, _omitEnumNames ? '' : 'RUNTIME_NODEJS610');
  static const Runtime RUNTIME_PYTHON310 =
      Runtime._(19, _omitEnumNames ? '' : 'RUNTIME_PYTHON310');
  static const Runtime RUNTIME_PYTHON38 =
      Runtime._(20, _omitEnumNames ? '' : 'RUNTIME_PYTHON38');
  static const Runtime RUNTIME_JAVA17 =
      Runtime._(21, _omitEnumNames ? '' : 'RUNTIME_JAVA17');
  static const Runtime RUNTIME_NODEJS43 =
      Runtime._(22, _omitEnumNames ? '' : 'RUNTIME_NODEJS43');
  static const Runtime RUNTIME_NODEJS =
      Runtime._(23, _omitEnumNames ? '' : 'RUNTIME_NODEJS');
  static const Runtime RUNTIME_NODEJS20X =
      Runtime._(24, _omitEnumNames ? '' : 'RUNTIME_NODEJS20X');
  static const Runtime RUNTIME_PYTHON313 =
      Runtime._(25, _omitEnumNames ? '' : 'RUNTIME_PYTHON313');
  static const Runtime RUNTIME_DOTNETCORE20 =
      Runtime._(26, _omitEnumNames ? '' : 'RUNTIME_DOTNETCORE20');
  static const Runtime RUNTIME_NODEJS43EDGE =
      Runtime._(27, _omitEnumNames ? '' : 'RUNTIME_NODEJS43EDGE');
  static const Runtime RUNTIME_JAVA11 =
      Runtime._(28, _omitEnumNames ? '' : 'RUNTIME_JAVA11');
  static const Runtime RUNTIME_NODEJS12X =
      Runtime._(29, _omitEnumNames ? '' : 'RUNTIME_NODEJS12X');
  static const Runtime RUNTIME_RUBY33 =
      Runtime._(30, _omitEnumNames ? '' : 'RUNTIME_RUBY33');
  static const Runtime RUNTIME_JAVA8 =
      Runtime._(31, _omitEnumNames ? '' : 'RUNTIME_JAVA8');
  static const Runtime RUNTIME_DOTNETCORE10 =
      Runtime._(32, _omitEnumNames ? '' : 'RUNTIME_DOTNETCORE10');
  static const Runtime RUNTIME_PYTHON37 =
      Runtime._(33, _omitEnumNames ? '' : 'RUNTIME_PYTHON37');
  static const Runtime RUNTIME_PYTHON311 =
      Runtime._(34, _omitEnumNames ? '' : 'RUNTIME_PYTHON311');
  static const Runtime RUNTIME_JAVA8AL2 =
      Runtime._(35, _omitEnumNames ? '' : 'RUNTIME_JAVA8AL2');
  static const Runtime RUNTIME_RUBY34 =
      Runtime._(36, _omitEnumNames ? '' : 'RUNTIME_RUBY34');
  static const Runtime RUNTIME_PROVIDEDAL2 =
      Runtime._(37, _omitEnumNames ? '' : 'RUNTIME_PROVIDEDAL2');
  static const Runtime RUNTIME_NODEJS10X =
      Runtime._(38, _omitEnumNames ? '' : 'RUNTIME_NODEJS10X');
  static const Runtime RUNTIME_JAVA25 =
      Runtime._(39, _omitEnumNames ? '' : 'RUNTIME_JAVA25');
  static const Runtime RUNTIME_PYTHON27 =
      Runtime._(40, _omitEnumNames ? '' : 'RUNTIME_PYTHON27');
  static const Runtime RUNTIME_DOTNET8 =
      Runtime._(41, _omitEnumNames ? '' : 'RUNTIME_DOTNET8');
  static const Runtime RUNTIME_NODEJS18X =
      Runtime._(42, _omitEnumNames ? '' : 'RUNTIME_NODEJS18X');
  static const Runtime RUNTIME_DOTNET10 =
      Runtime._(43, _omitEnumNames ? '' : 'RUNTIME_DOTNET10');
  static const Runtime RUNTIME_PYTHON314 =
      Runtime._(44, _omitEnumNames ? '' : 'RUNTIME_PYTHON314');

  static const $core.List<Runtime> values = <Runtime>[
    RUNTIME_DOTNETCORE31,
    RUNTIME_NODEJS16X,
    RUNTIME_RUBY25,
    RUNTIME_NODEJS24X,
    RUNTIME_PYTHON39,
    RUNTIME_DOTNET6,
    RUNTIME_RUBY32,
    RUNTIME_GO1X,
    RUNTIME_PYTHON36,
    RUNTIME_PROVIDEDAL2023,
    RUNTIME_NODEJS22X,
    RUNTIME_PYTHON312,
    RUNTIME_DOTNETCORE21,
    RUNTIME_NODEJS810,
    RUNTIME_RUBY27,
    RUNTIME_JAVA21,
    RUNTIME_NODEJS14X,
    RUNTIME_PROVIDED,
    RUNTIME_NODEJS610,
    RUNTIME_PYTHON310,
    RUNTIME_PYTHON38,
    RUNTIME_JAVA17,
    RUNTIME_NODEJS43,
    RUNTIME_NODEJS,
    RUNTIME_NODEJS20X,
    RUNTIME_PYTHON313,
    RUNTIME_DOTNETCORE20,
    RUNTIME_NODEJS43EDGE,
    RUNTIME_JAVA11,
    RUNTIME_NODEJS12X,
    RUNTIME_RUBY33,
    RUNTIME_JAVA8,
    RUNTIME_DOTNETCORE10,
    RUNTIME_PYTHON37,
    RUNTIME_PYTHON311,
    RUNTIME_JAVA8AL2,
    RUNTIME_RUBY34,
    RUNTIME_PROVIDEDAL2,
    RUNTIME_NODEJS10X,
    RUNTIME_JAVA25,
    RUNTIME_PYTHON27,
    RUNTIME_DOTNET8,
    RUNTIME_NODEJS18X,
    RUNTIME_DOTNET10,
    RUNTIME_PYTHON314,
  ];

  static final $core.List<Runtime?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 44);
  static Runtime? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Runtime._(super.value, super.name);
}

class SchemaRegistryEventRecordFormat extends $pb.ProtobufEnum {
  static const SchemaRegistryEventRecordFormat
      SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_JSON =
      SchemaRegistryEventRecordFormat._(
          0, _omitEnumNames ? '' : 'SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_JSON');
  static const SchemaRegistryEventRecordFormat
      SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_SOURCE =
      SchemaRegistryEventRecordFormat._(1,
          _omitEnumNames ? '' : 'SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_SOURCE');

  static const $core.List<SchemaRegistryEventRecordFormat> values =
      <SchemaRegistryEventRecordFormat>[
    SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_JSON,
    SCHEMA_REGISTRY_EVENT_RECORD_FORMAT_SOURCE,
  ];

  static final $core.List<SchemaRegistryEventRecordFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SchemaRegistryEventRecordFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SchemaRegistryEventRecordFormat._(super.value, super.name);
}

class SnapStartApplyOn extends $pb.ProtobufEnum {
  static const SnapStartApplyOn SNAP_START_APPLY_ON_NONE =
      SnapStartApplyOn._(0, _omitEnumNames ? '' : 'SNAP_START_APPLY_ON_NONE');
  static const SnapStartApplyOn SNAP_START_APPLY_ON_PUBLISHEDVERSIONS =
      SnapStartApplyOn._(
          1, _omitEnumNames ? '' : 'SNAP_START_APPLY_ON_PUBLISHEDVERSIONS');

  static const $core.List<SnapStartApplyOn> values = <SnapStartApplyOn>[
    SNAP_START_APPLY_ON_NONE,
    SNAP_START_APPLY_ON_PUBLISHEDVERSIONS,
  ];

  static final $core.List<SnapStartApplyOn?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SnapStartApplyOn? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SnapStartApplyOn._(super.value, super.name);
}

class SnapStartOptimizationStatus extends $pb.ProtobufEnum {
  static const SnapStartOptimizationStatus SNAP_START_OPTIMIZATION_STATUS_ON =
      SnapStartOptimizationStatus._(
          0, _omitEnumNames ? '' : 'SNAP_START_OPTIMIZATION_STATUS_ON');
  static const SnapStartOptimizationStatus SNAP_START_OPTIMIZATION_STATUS_OFF =
      SnapStartOptimizationStatus._(
          1, _omitEnumNames ? '' : 'SNAP_START_OPTIMIZATION_STATUS_OFF');

  static const $core.List<SnapStartOptimizationStatus> values =
      <SnapStartOptimizationStatus>[
    SNAP_START_OPTIMIZATION_STATUS_ON,
    SNAP_START_OPTIMIZATION_STATUS_OFF,
  ];

  static final $core.List<SnapStartOptimizationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SnapStartOptimizationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SnapStartOptimizationStatus._(super.value, super.name);
}

class SourceAccessType extends $pb.ProtobufEnum {
  static const SourceAccessType SOURCE_ACCESS_TYPE_VPC_SUBNET =
      SourceAccessType._(
          0, _omitEnumNames ? '' : 'SOURCE_ACCESS_TYPE_VPC_SUBNET');
  static const SourceAccessType SOURCE_ACCESS_TYPE_VPC_SECURITY_GROUP =
      SourceAccessType._(
          1, _omitEnumNames ? '' : 'SOURCE_ACCESS_TYPE_VPC_SECURITY_GROUP');
  static const SourceAccessType SOURCE_ACCESS_TYPE_CLIENT_CERTIFICATE_TLS_AUTH =
      SourceAccessType._(
          2,
          _omitEnumNames
              ? ''
              : 'SOURCE_ACCESS_TYPE_CLIENT_CERTIFICATE_TLS_AUTH');
  static const SourceAccessType SOURCE_ACCESS_TYPE_SASL_SCRAM_256_AUTH =
      SourceAccessType._(
          3, _omitEnumNames ? '' : 'SOURCE_ACCESS_TYPE_SASL_SCRAM_256_AUTH');
  static const SourceAccessType SOURCE_ACCESS_TYPE_VIRTUAL_HOST =
      SourceAccessType._(
          4, _omitEnumNames ? '' : 'SOURCE_ACCESS_TYPE_VIRTUAL_HOST');
  static const SourceAccessType SOURCE_ACCESS_TYPE_SASL_SCRAM_512_AUTH =
      SourceAccessType._(
          5, _omitEnumNames ? '' : 'SOURCE_ACCESS_TYPE_SASL_SCRAM_512_AUTH');
  static const SourceAccessType SOURCE_ACCESS_TYPE_SERVER_ROOT_CA_CERTIFICATE =
      SourceAccessType._(
          6,
          _omitEnumNames
              ? ''
              : 'SOURCE_ACCESS_TYPE_SERVER_ROOT_CA_CERTIFICATE');
  static const SourceAccessType SOURCE_ACCESS_TYPE_BASIC_AUTH =
      SourceAccessType._(
          7, _omitEnumNames ? '' : 'SOURCE_ACCESS_TYPE_BASIC_AUTH');

  static const $core.List<SourceAccessType> values = <SourceAccessType>[
    SOURCE_ACCESS_TYPE_VPC_SUBNET,
    SOURCE_ACCESS_TYPE_VPC_SECURITY_GROUP,
    SOURCE_ACCESS_TYPE_CLIENT_CERTIFICATE_TLS_AUTH,
    SOURCE_ACCESS_TYPE_SASL_SCRAM_256_AUTH,
    SOURCE_ACCESS_TYPE_VIRTUAL_HOST,
    SOURCE_ACCESS_TYPE_SASL_SCRAM_512_AUTH,
    SOURCE_ACCESS_TYPE_SERVER_ROOT_CA_CERTIFICATE,
    SOURCE_ACCESS_TYPE_BASIC_AUTH,
  ];

  static final $core.List<SourceAccessType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static SourceAccessType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SourceAccessType._(super.value, super.name);
}

class State extends $pb.ProtobufEnum {
  static const State STATE_ACTIVE =
      State._(0, _omitEnumNames ? '' : 'STATE_ACTIVE');
  static const State STATE_DEACTIVATED =
      State._(1, _omitEnumNames ? '' : 'STATE_DEACTIVATED');
  static const State STATE_FAILED =
      State._(2, _omitEnumNames ? '' : 'STATE_FAILED');
  static const State STATE_ACTIVENONINVOCABLE =
      State._(3, _omitEnumNames ? '' : 'STATE_ACTIVENONINVOCABLE');
  static const State STATE_INACTIVE =
      State._(4, _omitEnumNames ? '' : 'STATE_INACTIVE');
  static const State STATE_DELETING =
      State._(5, _omitEnumNames ? '' : 'STATE_DELETING');
  static const State STATE_DEACTIVATING =
      State._(6, _omitEnumNames ? '' : 'STATE_DEACTIVATING');
  static const State STATE_PENDING =
      State._(7, _omitEnumNames ? '' : 'STATE_PENDING');

  static const $core.List<State> values = <State>[
    STATE_ACTIVE,
    STATE_DEACTIVATED,
    STATE_FAILED,
    STATE_ACTIVENONINVOCABLE,
    STATE_INACTIVE,
    STATE_DELETING,
    STATE_DEACTIVATING,
    STATE_PENDING,
  ];

  static final $core.List<State?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static State? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const State._(super.value, super.name);
}

class StateReasonCode extends $pb.ProtobufEnum {
  static const StateReasonCode STATE_REASON_CODE_ENILIMITEXCEEDED =
      StateReasonCode._(
          0, _omitEnumNames ? '' : 'STATE_REASON_CODE_ENILIMITEXCEEDED');
  static const StateReasonCode STATE_REASON_CODE_SUBNETOUTOFIPADDRESSES =
      StateReasonCode._(
          1, _omitEnumNames ? '' : 'STATE_REASON_CODE_SUBNETOUTOFIPADDRESSES');
  static const StateReasonCode STATE_REASON_CODE_IDLE =
      StateReasonCode._(2, _omitEnumNames ? '' : 'STATE_REASON_CODE_IDLE');
  static const StateReasonCode STATE_REASON_CODE_KMSKEYNOTFOUND =
      StateReasonCode._(
          3, _omitEnumNames ? '' : 'STATE_REASON_CODE_KMSKEYNOTFOUND');
  static const StateReasonCode STATE_REASON_CODE_FUNCTIONERRORINITTIMEOUT =
      StateReasonCode._(4,
          _omitEnumNames ? '' : 'STATE_REASON_CODE_FUNCTIONERRORINITTIMEOUT');
  static const StateReasonCode STATE_REASON_CODE_INVALIDZIPFILEEXCEPTION =
      StateReasonCode._(
          5, _omitEnumNames ? '' : 'STATE_REASON_CODE_INVALIDZIPFILEEXCEPTION');
  static const StateReasonCode
      STATE_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL = StateReasonCode._(
          6,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL');
  static const StateReasonCode STATE_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR =
      StateReasonCode._(
          7,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR');
  static const StateReasonCode STATE_REASON_CODE_RESTORING =
      StateReasonCode._(8, _omitEnumNames ? '' : 'STATE_REASON_CODE_RESTORING');
  static const StateReasonCode STATE_REASON_CODE_EC2REQUESTLIMITEXCEEDED =
      StateReasonCode._(
          9, _omitEnumNames ? '' : 'STATE_REASON_CODE_EC2REQUESTLIMITEXCEEDED');
  static const StateReasonCode STATE_REASON_CODE_INVALIDIMAGE =
      StateReasonCode._(
          10, _omitEnumNames ? '' : 'STATE_REASON_CODE_INVALIDIMAGE');
  static const StateReasonCode
      STATE_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED =
      StateReasonCode._(
          11,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED');
  static const StateReasonCode STATE_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED =
      StateReasonCode._(
          12,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED');
  static const StateReasonCode STATE_REASON_CODE_INVALIDRUNTIME =
      StateReasonCode._(
          13, _omitEnumNames ? '' : 'STATE_REASON_CODE_INVALIDRUNTIME');
  static const StateReasonCode STATE_REASON_CODE_DISABLEDKMSKEY =
      StateReasonCode._(
          14, _omitEnumNames ? '' : 'STATE_REASON_CODE_DISABLEDKMSKEY');
  static const StateReasonCode STATE_REASON_CODE_VCPULIMITEXCEEDED =
      StateReasonCode._(
          15, _omitEnumNames ? '' : 'STATE_REASON_CODE_VCPULIMITEXCEEDED');
  static const StateReasonCode
      STATE_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS = StateReasonCode._(
          16,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS');
  static const StateReasonCode STATE_REASON_CODE_FUNCTIONERROR =
      StateReasonCode._(
          17, _omitEnumNames ? '' : 'STATE_REASON_CODE_FUNCTIONERROR');
  static const StateReasonCode STATE_REASON_CODE_INVALIDCONFIGURATION =
      StateReasonCode._(
          18, _omitEnumNames ? '' : 'STATE_REASON_CODE_INVALIDCONFIGURATION');
  static const StateReasonCode STATE_REASON_CODE_INVALIDSUBNET =
      StateReasonCode._(
          19, _omitEnumNames ? '' : 'STATE_REASON_CODE_INVALIDSUBNET');
  static const StateReasonCode STATE_REASON_CODE_INVALIDSTATEKMSKEY =
      StateReasonCode._(
          20, _omitEnumNames ? '' : 'STATE_REASON_CODE_INVALIDSTATEKMSKEY');
  static const StateReasonCode STATE_REASON_CODE_IMAGEACCESSDENIED =
      StateReasonCode._(
          21, _omitEnumNames ? '' : 'STATE_REASON_CODE_IMAGEACCESSDENIED');
  static const StateReasonCode STATE_REASON_CODE_IMAGEDELETED =
      StateReasonCode._(
          22, _omitEnumNames ? '' : 'STATE_REASON_CODE_IMAGEDELETED');
  static const StateReasonCode STATE_REASON_CODE_EFSMOUNTTIMEOUT =
      StateReasonCode._(
          23, _omitEnumNames ? '' : 'STATE_REASON_CODE_EFSMOUNTTIMEOUT');
  static const StateReasonCode STATE_REASON_CODE_EFSMOUNTFAILURE =
      StateReasonCode._(
          24, _omitEnumNames ? '' : 'STATE_REASON_CODE_EFSMOUNTFAILURE');
  static const StateReasonCode STATE_REASON_CODE_EFSIOERROR = StateReasonCode._(
      25, _omitEnumNames ? '' : 'STATE_REASON_CODE_EFSIOERROR');
  static const StateReasonCode
      STATE_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR = StateReasonCode._(
          26,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR');
  static const StateReasonCode STATE_REASON_CODE_CREATING =
      StateReasonCode._(27, _omitEnumNames ? '' : 'STATE_REASON_CODE_CREATING');
  static const StateReasonCode STATE_REASON_CODE_INVALIDSECURITYGROUP =
      StateReasonCode._(
          28, _omitEnumNames ? '' : 'STATE_REASON_CODE_INVALIDSECURITYGROUP');
  static const StateReasonCode STATE_REASON_CODE_DRAININGDURABLEEXECUTIONS =
      StateReasonCode._(29,
          _omitEnumNames ? '' : 'STATE_REASON_CODE_DRAININGDURABLEEXECUTIONS');
  static const StateReasonCode
      STATE_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT = StateReasonCode._(
          30,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT');
  static const StateReasonCode STATE_REASON_CODE_EFSMOUNTCONNECTIVITYERROR =
      StateReasonCode._(31,
          _omitEnumNames ? '' : 'STATE_REASON_CODE_EFSMOUNTCONNECTIVITYERROR');
  static const StateReasonCode
      STATE_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY =
      StateReasonCode._(
          32,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY');
  static const StateReasonCode STATE_REASON_CODE_INSUFFICIENTROLEPERMISSIONS =
      StateReasonCode._(
          33,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_INSUFFICIENTROLEPERMISSIONS');
  static const StateReasonCode STATE_REASON_CODE_KMSKEYACCESSDENIED =
      StateReasonCode._(
          34, _omitEnumNames ? '' : 'STATE_REASON_CODE_KMSKEYACCESSDENIED');
  static const StateReasonCode STATE_REASON_CODE_INTERNALERROR =
      StateReasonCode._(
          35, _omitEnumNames ? '' : 'STATE_REASON_CODE_INTERNALERROR');
  static const StateReasonCode STATE_REASON_CODE_INSUFFICIENTCAPACITY =
      StateReasonCode._(
          36, _omitEnumNames ? '' : 'STATE_REASON_CODE_INSUFFICIENTCAPACITY');
  static const StateReasonCode
      STATE_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED = StateReasonCode._(
          37,
          _omitEnumNames
              ? ''
              : 'STATE_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED');

  static const $core.List<StateReasonCode> values = <StateReasonCode>[
    STATE_REASON_CODE_ENILIMITEXCEEDED,
    STATE_REASON_CODE_SUBNETOUTOFIPADDRESSES,
    STATE_REASON_CODE_IDLE,
    STATE_REASON_CODE_KMSKEYNOTFOUND,
    STATE_REASON_CODE_FUNCTIONERRORINITTIMEOUT,
    STATE_REASON_CODE_INVALIDZIPFILEEXCEPTION,
    STATE_REASON_CODE_DISALLOWEDBYVPCENCRYPTIONCONTROL,
    STATE_REASON_CODE_FUNCTIONERRORRUNTIMEINITERROR,
    STATE_REASON_CODE_RESTORING,
    STATE_REASON_CODE_EC2REQUESTLIMITEXCEEDED,
    STATE_REASON_CODE_INVALIDIMAGE,
    STATE_REASON_CODE_CAPACITYPROVIDERSCALINGLIMITEXCEEDED,
    STATE_REASON_CODE_FUNCTIONERRORPERMISSIONDENIED,
    STATE_REASON_CODE_INVALIDRUNTIME,
    STATE_REASON_CODE_DISABLEDKMSKEY,
    STATE_REASON_CODE_VCPULIMITEXCEEDED,
    STATE_REASON_CODE_FUNCTIONERRORTOOMANYEXTENSIONS,
    STATE_REASON_CODE_FUNCTIONERROR,
    STATE_REASON_CODE_INVALIDCONFIGURATION,
    STATE_REASON_CODE_INVALIDSUBNET,
    STATE_REASON_CODE_INVALIDSTATEKMSKEY,
    STATE_REASON_CODE_IMAGEACCESSDENIED,
    STATE_REASON_CODE_IMAGEDELETED,
    STATE_REASON_CODE_EFSMOUNTTIMEOUT,
    STATE_REASON_CODE_EFSMOUNTFAILURE,
    STATE_REASON_CODE_EFSIOERROR,
    STATE_REASON_CODE_FUNCTIONERROREXTENSIONINITERROR,
    STATE_REASON_CODE_CREATING,
    STATE_REASON_CODE_INVALIDSECURITYGROUP,
    STATE_REASON_CODE_DRAININGDURABLEEXECUTIONS,
    STATE_REASON_CODE_FUNCTIONERRORINVALIDENTRYPOINT,
    STATE_REASON_CODE_EFSMOUNTCONNECTIVITYERROR,
    STATE_REASON_CODE_FUNCTIONERRORINVALIDWORKINGDIRECTORY,
    STATE_REASON_CODE_INSUFFICIENTROLEPERMISSIONS,
    STATE_REASON_CODE_KMSKEYACCESSDENIED,
    STATE_REASON_CODE_INTERNALERROR,
    STATE_REASON_CODE_INSUFFICIENTCAPACITY,
    STATE_REASON_CODE_FUNCTIONERRORINITRESOURCEEXHAUSTED,
  ];

  static final $core.List<StateReasonCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 37);
  static StateReasonCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StateReasonCode._(super.value, super.name);
}

class SystemLogLevel extends $pb.ProtobufEnum {
  static const SystemLogLevel SYSTEM_LOG_LEVEL_WARN =
      SystemLogLevel._(0, _omitEnumNames ? '' : 'SYSTEM_LOG_LEVEL_WARN');
  static const SystemLogLevel SYSTEM_LOG_LEVEL_INFO =
      SystemLogLevel._(1, _omitEnumNames ? '' : 'SYSTEM_LOG_LEVEL_INFO');
  static const SystemLogLevel SYSTEM_LOG_LEVEL_DEBUG =
      SystemLogLevel._(2, _omitEnumNames ? '' : 'SYSTEM_LOG_LEVEL_DEBUG');

  static const $core.List<SystemLogLevel> values = <SystemLogLevel>[
    SYSTEM_LOG_LEVEL_WARN,
    SYSTEM_LOG_LEVEL_INFO,
    SYSTEM_LOG_LEVEL_DEBUG,
  ];

  static final $core.List<SystemLogLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SystemLogLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SystemLogLevel._(super.value, super.name);
}

class TenantIsolationMode extends $pb.ProtobufEnum {
  static const TenantIsolationMode TENANT_ISOLATION_MODE_PER_TENANT =
      TenantIsolationMode._(
          0, _omitEnumNames ? '' : 'TENANT_ISOLATION_MODE_PER_TENANT');

  static const $core.List<TenantIsolationMode> values = <TenantIsolationMode>[
    TENANT_ISOLATION_MODE_PER_TENANT,
  ];

  static final $core.List<TenantIsolationMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static TenantIsolationMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TenantIsolationMode._(super.value, super.name);
}

class ThrottleReason extends $pb.ProtobufEnum {
  static const ThrottleReason
      THROTTLE_REASON_CONCURRENTSNAPSHOTCREATELIMITEXCEEDED = ThrottleReason._(
          0,
          _omitEnumNames
              ? ''
              : 'THROTTLE_REASON_CONCURRENTSNAPSHOTCREATELIMITEXCEEDED');
  static const ThrottleReason
      THROTTLE_REASON_FUNCTIONINVOCATIONRATELIMITEXCEEDED = ThrottleReason._(
          1,
          _omitEnumNames
              ? ''
              : 'THROTTLE_REASON_FUNCTIONINVOCATIONRATELIMITEXCEEDED');
  static const ThrottleReason
      THROTTLE_REASON_CONCURRENTINVOCATIONLIMITEXCEEDED = ThrottleReason._(
          2,
          _omitEnumNames
              ? ''
              : 'THROTTLE_REASON_CONCURRENTINVOCATIONLIMITEXCEEDED');
  static const ThrottleReason THROTTLE_REASON_CALLERRATELIMITEXCEEDED =
      ThrottleReason._(
          3, _omitEnumNames ? '' : 'THROTTLE_REASON_CALLERRATELIMITEXCEEDED');
  static const ThrottleReason
      THROTTLE_REASON_RESERVEDFUNCTIONCONCURRENTINVOCATIONLIMITEXCEEDED =
      ThrottleReason._(
          4,
          _omitEnumNames
              ? ''
              : 'THROTTLE_REASON_RESERVEDFUNCTIONCONCURRENTINVOCATIONLIMITEXCEEDED');
  static const ThrottleReason
      THROTTLE_REASON_RESERVEDFUNCTIONINVOCATIONRATELIMITEXCEEDED =
      ThrottleReason._(
          5,
          _omitEnumNames
              ? ''
              : 'THROTTLE_REASON_RESERVEDFUNCTIONINVOCATIONRATELIMITEXCEEDED');

  static const $core.List<ThrottleReason> values = <ThrottleReason>[
    THROTTLE_REASON_CONCURRENTSNAPSHOTCREATELIMITEXCEEDED,
    THROTTLE_REASON_FUNCTIONINVOCATIONRATELIMITEXCEEDED,
    THROTTLE_REASON_CONCURRENTINVOCATIONLIMITEXCEEDED,
    THROTTLE_REASON_CALLERRATELIMITEXCEEDED,
    THROTTLE_REASON_RESERVEDFUNCTIONCONCURRENTINVOCATIONLIMITEXCEEDED,
    THROTTLE_REASON_RESERVEDFUNCTIONINVOCATIONRATELIMITEXCEEDED,
  ];

  static final $core.List<ThrottleReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ThrottleReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ThrottleReason._(super.value, super.name);
}

class TracingMode extends $pb.ProtobufEnum {
  static const TracingMode TRACING_MODE_ACTIVE =
      TracingMode._(0, _omitEnumNames ? '' : 'TRACING_MODE_ACTIVE');
  static const TracingMode TRACING_MODE_PASSTHROUGH =
      TracingMode._(1, _omitEnumNames ? '' : 'TRACING_MODE_PASSTHROUGH');

  static const $core.List<TracingMode> values = <TracingMode>[
    TRACING_MODE_ACTIVE,
    TRACING_MODE_PASSTHROUGH,
  ];

  static final $core.List<TracingMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static TracingMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TracingMode._(super.value, super.name);
}

class UpdateRuntimeOn extends $pb.ProtobufEnum {
  static const UpdateRuntimeOn UPDATE_RUNTIME_ON_MANUAL =
      UpdateRuntimeOn._(0, _omitEnumNames ? '' : 'UPDATE_RUNTIME_ON_MANUAL');
  static const UpdateRuntimeOn UPDATE_RUNTIME_ON_AUTO =
      UpdateRuntimeOn._(1, _omitEnumNames ? '' : 'UPDATE_RUNTIME_ON_AUTO');
  static const UpdateRuntimeOn UPDATE_RUNTIME_ON_FUNCTIONUPDATE =
      UpdateRuntimeOn._(
          2, _omitEnumNames ? '' : 'UPDATE_RUNTIME_ON_FUNCTIONUPDATE');

  static const $core.List<UpdateRuntimeOn> values = <UpdateRuntimeOn>[
    UPDATE_RUNTIME_ON_MANUAL,
    UPDATE_RUNTIME_ON_AUTO,
    UPDATE_RUNTIME_ON_FUNCTIONUPDATE,
  ];

  static final $core.List<UpdateRuntimeOn?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static UpdateRuntimeOn? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const UpdateRuntimeOn._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
