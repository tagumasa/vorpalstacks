// This is a generated file - do not edit.
//
// Generated from events.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ApiDestinationHttpMethod extends $pb.ProtobufEnum {
  static const ApiDestinationHttpMethod API_DESTINATION_HTTP_METHOD_OPTIONS =
      ApiDestinationHttpMethod._(
          0, _omitEnumNames ? '' : 'API_DESTINATION_HTTP_METHOD_OPTIONS');
  static const ApiDestinationHttpMethod API_DESTINATION_HTTP_METHOD_PATCH =
      ApiDestinationHttpMethod._(
          1, _omitEnumNames ? '' : 'API_DESTINATION_HTTP_METHOD_PATCH');
  static const ApiDestinationHttpMethod API_DESTINATION_HTTP_METHOD_POST =
      ApiDestinationHttpMethod._(
          2, _omitEnumNames ? '' : 'API_DESTINATION_HTTP_METHOD_POST');
  static const ApiDestinationHttpMethod API_DESTINATION_HTTP_METHOD_HEAD =
      ApiDestinationHttpMethod._(
          3, _omitEnumNames ? '' : 'API_DESTINATION_HTTP_METHOD_HEAD');
  static const ApiDestinationHttpMethod API_DESTINATION_HTTP_METHOD_GET =
      ApiDestinationHttpMethod._(
          4, _omitEnumNames ? '' : 'API_DESTINATION_HTTP_METHOD_GET');
  static const ApiDestinationHttpMethod API_DESTINATION_HTTP_METHOD_DELETE =
      ApiDestinationHttpMethod._(
          5, _omitEnumNames ? '' : 'API_DESTINATION_HTTP_METHOD_DELETE');
  static const ApiDestinationHttpMethod API_DESTINATION_HTTP_METHOD_PUT =
      ApiDestinationHttpMethod._(
          6, _omitEnumNames ? '' : 'API_DESTINATION_HTTP_METHOD_PUT');

  static const $core.List<ApiDestinationHttpMethod> values =
      <ApiDestinationHttpMethod>[
    API_DESTINATION_HTTP_METHOD_OPTIONS,
    API_DESTINATION_HTTP_METHOD_PATCH,
    API_DESTINATION_HTTP_METHOD_POST,
    API_DESTINATION_HTTP_METHOD_HEAD,
    API_DESTINATION_HTTP_METHOD_GET,
    API_DESTINATION_HTTP_METHOD_DELETE,
    API_DESTINATION_HTTP_METHOD_PUT,
  ];

  static final $core.List<ApiDestinationHttpMethod?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static ApiDestinationHttpMethod? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ApiDestinationHttpMethod._(super.value, super.name);
}

class ApiDestinationState extends $pb.ProtobufEnum {
  static const ApiDestinationState API_DESTINATION_STATE_ACTIVE =
      ApiDestinationState._(
          0, _omitEnumNames ? '' : 'API_DESTINATION_STATE_ACTIVE');
  static const ApiDestinationState API_DESTINATION_STATE_INACTIVE =
      ApiDestinationState._(
          1, _omitEnumNames ? '' : 'API_DESTINATION_STATE_INACTIVE');

  static const $core.List<ApiDestinationState> values = <ApiDestinationState>[
    API_DESTINATION_STATE_ACTIVE,
    API_DESTINATION_STATE_INACTIVE,
  ];

  static final $core.List<ApiDestinationState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ApiDestinationState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ApiDestinationState._(super.value, super.name);
}

class ArchiveState extends $pb.ProtobufEnum {
  static const ArchiveState ARCHIVE_STATE_UPDATING =
      ArchiveState._(0, _omitEnumNames ? '' : 'ARCHIVE_STATE_UPDATING');
  static const ArchiveState ARCHIVE_STATE_UPDATE_FAILED =
      ArchiveState._(1, _omitEnumNames ? '' : 'ARCHIVE_STATE_UPDATE_FAILED');
  static const ArchiveState ARCHIVE_STATE_DISABLED =
      ArchiveState._(2, _omitEnumNames ? '' : 'ARCHIVE_STATE_DISABLED');
  static const ArchiveState ARCHIVE_STATE_CREATING =
      ArchiveState._(3, _omitEnumNames ? '' : 'ARCHIVE_STATE_CREATING');
  static const ArchiveState ARCHIVE_STATE_ENABLED =
      ArchiveState._(4, _omitEnumNames ? '' : 'ARCHIVE_STATE_ENABLED');
  static const ArchiveState ARCHIVE_STATE_CREATE_FAILED =
      ArchiveState._(5, _omitEnumNames ? '' : 'ARCHIVE_STATE_CREATE_FAILED');

  static const $core.List<ArchiveState> values = <ArchiveState>[
    ARCHIVE_STATE_UPDATING,
    ARCHIVE_STATE_UPDATE_FAILED,
    ARCHIVE_STATE_DISABLED,
    ARCHIVE_STATE_CREATING,
    ARCHIVE_STATE_ENABLED,
    ARCHIVE_STATE_CREATE_FAILED,
  ];

  static final $core.List<ArchiveState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ArchiveState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ArchiveState._(super.value, super.name);
}

class AssignPublicIp extends $pb.ProtobufEnum {
  static const AssignPublicIp ASSIGN_PUBLIC_IP_DISABLED =
      AssignPublicIp._(0, _omitEnumNames ? '' : 'ASSIGN_PUBLIC_IP_DISABLED');
  static const AssignPublicIp ASSIGN_PUBLIC_IP_ENABLED =
      AssignPublicIp._(1, _omitEnumNames ? '' : 'ASSIGN_PUBLIC_IP_ENABLED');

  static const $core.List<AssignPublicIp> values = <AssignPublicIp>[
    ASSIGN_PUBLIC_IP_DISABLED,
    ASSIGN_PUBLIC_IP_ENABLED,
  ];

  static final $core.List<AssignPublicIp?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AssignPublicIp? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssignPublicIp._(super.value, super.name);
}

class ConnectionAuthorizationType extends $pb.ProtobufEnum {
  static const ConnectionAuthorizationType
      CONNECTION_AUTHORIZATION_TYPE_OAUTH_CLIENT_CREDENTIALS =
      ConnectionAuthorizationType._(
          0,
          _omitEnumNames
              ? ''
              : 'CONNECTION_AUTHORIZATION_TYPE_OAUTH_CLIENT_CREDENTIALS');
  static const ConnectionAuthorizationType CONNECTION_AUTHORIZATION_TYPE_BASIC =
      ConnectionAuthorizationType._(
          1, _omitEnumNames ? '' : 'CONNECTION_AUTHORIZATION_TYPE_BASIC');
  static const ConnectionAuthorizationType
      CONNECTION_AUTHORIZATION_TYPE_API_KEY = ConnectionAuthorizationType._(
          2, _omitEnumNames ? '' : 'CONNECTION_AUTHORIZATION_TYPE_API_KEY');

  static const $core.List<ConnectionAuthorizationType> values =
      <ConnectionAuthorizationType>[
    CONNECTION_AUTHORIZATION_TYPE_OAUTH_CLIENT_CREDENTIALS,
    CONNECTION_AUTHORIZATION_TYPE_BASIC,
    CONNECTION_AUTHORIZATION_TYPE_API_KEY,
  ];

  static final $core.List<ConnectionAuthorizationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ConnectionAuthorizationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionAuthorizationType._(super.value, super.name);
}

class ConnectionOAuthHttpMethod extends $pb.ProtobufEnum {
  static const ConnectionOAuthHttpMethod CONNECTION_O_AUTH_HTTP_METHOD_POST =
      ConnectionOAuthHttpMethod._(
          0, _omitEnumNames ? '' : 'CONNECTION_O_AUTH_HTTP_METHOD_POST');
  static const ConnectionOAuthHttpMethod CONNECTION_O_AUTH_HTTP_METHOD_GET =
      ConnectionOAuthHttpMethod._(
          1, _omitEnumNames ? '' : 'CONNECTION_O_AUTH_HTTP_METHOD_GET');
  static const ConnectionOAuthHttpMethod CONNECTION_O_AUTH_HTTP_METHOD_PUT =
      ConnectionOAuthHttpMethod._(
          2, _omitEnumNames ? '' : 'CONNECTION_O_AUTH_HTTP_METHOD_PUT');

  static const $core.List<ConnectionOAuthHttpMethod> values =
      <ConnectionOAuthHttpMethod>[
    CONNECTION_O_AUTH_HTTP_METHOD_POST,
    CONNECTION_O_AUTH_HTTP_METHOD_GET,
    CONNECTION_O_AUTH_HTTP_METHOD_PUT,
  ];

  static final $core.List<ConnectionOAuthHttpMethod?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ConnectionOAuthHttpMethod? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionOAuthHttpMethod._(super.value, super.name);
}

class ConnectionState extends $pb.ProtobufEnum {
  static const ConnectionState CONNECTION_STATE_DEAUTHORIZED =
      ConnectionState._(
          0, _omitEnumNames ? '' : 'CONNECTION_STATE_DEAUTHORIZED');
  static const ConnectionState CONNECTION_STATE_UPDATING =
      ConnectionState._(1, _omitEnumNames ? '' : 'CONNECTION_STATE_UPDATING');
  static const ConnectionState CONNECTION_STATE_AUTHORIZING = ConnectionState._(
      2, _omitEnumNames ? '' : 'CONNECTION_STATE_AUTHORIZING');
  static const ConnectionState CONNECTION_STATE_AUTHORIZED =
      ConnectionState._(3, _omitEnumNames ? '' : 'CONNECTION_STATE_AUTHORIZED');
  static const ConnectionState CONNECTION_STATE_DEAUTHORIZING =
      ConnectionState._(
          4, _omitEnumNames ? '' : 'CONNECTION_STATE_DEAUTHORIZING');
  static const ConnectionState CONNECTION_STATE_DELETING =
      ConnectionState._(5, _omitEnumNames ? '' : 'CONNECTION_STATE_DELETING');
  static const ConnectionState CONNECTION_STATE_CREATING =
      ConnectionState._(6, _omitEnumNames ? '' : 'CONNECTION_STATE_CREATING');

  static const $core.List<ConnectionState> values = <ConnectionState>[
    CONNECTION_STATE_DEAUTHORIZED,
    CONNECTION_STATE_UPDATING,
    CONNECTION_STATE_AUTHORIZING,
    CONNECTION_STATE_AUTHORIZED,
    CONNECTION_STATE_DEAUTHORIZING,
    CONNECTION_STATE_DELETING,
    CONNECTION_STATE_CREATING,
  ];

  static final $core.List<ConnectionState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static ConnectionState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionState._(super.value, super.name);
}

class EventSourceState extends $pb.ProtobufEnum {
  static const EventSourceState EVENT_SOURCE_STATE_PENDING =
      EventSourceState._(0, _omitEnumNames ? '' : 'EVENT_SOURCE_STATE_PENDING');
  static const EventSourceState EVENT_SOURCE_STATE_ACTIVE =
      EventSourceState._(1, _omitEnumNames ? '' : 'EVENT_SOURCE_STATE_ACTIVE');
  static const EventSourceState EVENT_SOURCE_STATE_DELETED =
      EventSourceState._(2, _omitEnumNames ? '' : 'EVENT_SOURCE_STATE_DELETED');

  static const $core.List<EventSourceState> values = <EventSourceState>[
    EVENT_SOURCE_STATE_PENDING,
    EVENT_SOURCE_STATE_ACTIVE,
    EVENT_SOURCE_STATE_DELETED,
  ];

  static final $core.List<EventSourceState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EventSourceState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventSourceState._(super.value, super.name);
}

class LaunchType extends $pb.ProtobufEnum {
  static const LaunchType LAUNCH_TYPE_FARGATE =
      LaunchType._(0, _omitEnumNames ? '' : 'LAUNCH_TYPE_FARGATE');
  static const LaunchType LAUNCH_TYPE_EXTERNAL =
      LaunchType._(1, _omitEnumNames ? '' : 'LAUNCH_TYPE_EXTERNAL');
  static const LaunchType LAUNCH_TYPE_EC2 =
      LaunchType._(2, _omitEnumNames ? '' : 'LAUNCH_TYPE_EC2');

  static const $core.List<LaunchType> values = <LaunchType>[
    LAUNCH_TYPE_FARGATE,
    LAUNCH_TYPE_EXTERNAL,
    LAUNCH_TYPE_EC2,
  ];

  static final $core.List<LaunchType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static LaunchType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LaunchType._(super.value, super.name);
}

class PlacementConstraintType extends $pb.ProtobufEnum {
  static const PlacementConstraintType
      PLACEMENT_CONSTRAINT_TYPE_DISTINCT_INSTANCE = PlacementConstraintType._(0,
          _omitEnumNames ? '' : 'PLACEMENT_CONSTRAINT_TYPE_DISTINCT_INSTANCE');
  static const PlacementConstraintType PLACEMENT_CONSTRAINT_TYPE_MEMBER_OF =
      PlacementConstraintType._(
          1, _omitEnumNames ? '' : 'PLACEMENT_CONSTRAINT_TYPE_MEMBER_OF');

  static const $core.List<PlacementConstraintType> values =
      <PlacementConstraintType>[
    PLACEMENT_CONSTRAINT_TYPE_DISTINCT_INSTANCE,
    PLACEMENT_CONSTRAINT_TYPE_MEMBER_OF,
  ];

  static final $core.List<PlacementConstraintType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PlacementConstraintType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PlacementConstraintType._(super.value, super.name);
}

class PlacementStrategyType extends $pb.ProtobufEnum {
  static const PlacementStrategyType PLACEMENT_STRATEGY_TYPE_SPREAD =
      PlacementStrategyType._(
          0, _omitEnumNames ? '' : 'PLACEMENT_STRATEGY_TYPE_SPREAD');
  static const PlacementStrategyType PLACEMENT_STRATEGY_TYPE_RANDOM =
      PlacementStrategyType._(
          1, _omitEnumNames ? '' : 'PLACEMENT_STRATEGY_TYPE_RANDOM');
  static const PlacementStrategyType PLACEMENT_STRATEGY_TYPE_BINPACK =
      PlacementStrategyType._(
          2, _omitEnumNames ? '' : 'PLACEMENT_STRATEGY_TYPE_BINPACK');

  static const $core.List<PlacementStrategyType> values =
      <PlacementStrategyType>[
    PLACEMENT_STRATEGY_TYPE_SPREAD,
    PLACEMENT_STRATEGY_TYPE_RANDOM,
    PLACEMENT_STRATEGY_TYPE_BINPACK,
  ];

  static final $core.List<PlacementStrategyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static PlacementStrategyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PlacementStrategyType._(super.value, super.name);
}

class PropagateTags extends $pb.ProtobufEnum {
  static const PropagateTags PROPAGATE_TAGS_TASK_DEFINITION = PropagateTags._(
      0, _omitEnumNames ? '' : 'PROPAGATE_TAGS_TASK_DEFINITION');

  static const $core.List<PropagateTags> values = <PropagateTags>[
    PROPAGATE_TAGS_TASK_DEFINITION,
  ];

  static final $core.List<PropagateTags?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static PropagateTags? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PropagateTags._(super.value, super.name);
}

class ReplayState extends $pb.ProtobufEnum {
  static const ReplayState REPLAY_STATE_STARTING =
      ReplayState._(0, _omitEnumNames ? '' : 'REPLAY_STATE_STARTING');
  static const ReplayState REPLAY_STATE_RUNNING =
      ReplayState._(1, _omitEnumNames ? '' : 'REPLAY_STATE_RUNNING');
  static const ReplayState REPLAY_STATE_CANCELLED =
      ReplayState._(2, _omitEnumNames ? '' : 'REPLAY_STATE_CANCELLED');
  static const ReplayState REPLAY_STATE_CANCELLING =
      ReplayState._(3, _omitEnumNames ? '' : 'REPLAY_STATE_CANCELLING');
  static const ReplayState REPLAY_STATE_COMPLETED =
      ReplayState._(4, _omitEnumNames ? '' : 'REPLAY_STATE_COMPLETED');
  static const ReplayState REPLAY_STATE_FAILED =
      ReplayState._(5, _omitEnumNames ? '' : 'REPLAY_STATE_FAILED');

  static const $core.List<ReplayState> values = <ReplayState>[
    REPLAY_STATE_STARTING,
    REPLAY_STATE_RUNNING,
    REPLAY_STATE_CANCELLED,
    REPLAY_STATE_CANCELLING,
    REPLAY_STATE_COMPLETED,
    REPLAY_STATE_FAILED,
  ];

  static final $core.List<ReplayState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ReplayState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReplayState._(super.value, super.name);
}

class RuleState extends $pb.ProtobufEnum {
  static const RuleState RULE_STATE_DISABLED =
      RuleState._(0, _omitEnumNames ? '' : 'RULE_STATE_DISABLED');
  static const RuleState RULE_STATE_ENABLED =
      RuleState._(1, _omitEnumNames ? '' : 'RULE_STATE_ENABLED');

  static const $core.List<RuleState> values = <RuleState>[
    RULE_STATE_DISABLED,
    RULE_STATE_ENABLED,
  ];

  static final $core.List<RuleState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RuleState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RuleState._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
