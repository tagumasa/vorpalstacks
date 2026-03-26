// This is a generated file - do not edit.
//
// Generated from events.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'common.pb.dart' as $1;
import 'events.pb.dart' as $0;

export 'events.pb.dart';

/// CloudWatchEventsService provides events API operations.
@$pb.GrpcServiceName('events.CloudWatchEventsService')
class CloudWatchEventsServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CloudWatchEventsServiceClient(super.channel,
      {super.options, super.interceptors});

  /// Activates a partner event source that has been deactivated. Once activated, your matching event bus will start receiving events from the event source.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> activateEventSource(
    $0.ActivateEventSourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$activateEventSource, request, options: options);
  }

  /// Cancels the specified replay.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelReplayResponse> cancelReplay(
    $0.CancelReplayRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelReplay, request, options: options);
  }

  /// Creates an API destination, which is an HTTP invocation endpoint configured as a target for events.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateApiDestinationResponse> createApiDestination(
    $0.CreateApiDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createApiDestination, request, options: options);
  }

  /// Creates an archive of events with the specified settings. When you create an archive, incoming events might not immediately start being sent to the archive. Allow a short period of time for changes...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateArchiveResponse> createArchive(
    $0.CreateArchiveRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createArchive, request, options: options);
  }

  /// Creates a connection. A connection defines the authorization type and credentials to use for authorization with an API destination HTTP endpoint.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateConnectionResponse> createConnection(
    $0.CreateConnectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createConnection, request, options: options);
  }

  /// Creates a new event bus within your account. This can be a custom event bus which you can use to receive events from your custom applications and services, or it can be a partner event bus which ca...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateEventBusResponse> createEventBus(
    $0.CreateEventBusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createEventBus, request, options: options);
  }

  /// Called by an SaaS partner to create a partner event source. This operation is not used by Amazon Web Services customers. Each partner event source can be used by one Amazon Web Services account to ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreatePartnerEventSourceResponse>
      createPartnerEventSource(
    $0.CreatePartnerEventSourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPartnerEventSource, request,
        options: options);
  }

  /// You can use this operation to temporarily stop receiving events from the specified partner event source. The matching event bus is not deleted. When you deactivate a partner event source, the sourc...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deactivateEventSource(
    $0.DeactivateEventSourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deactivateEventSource, request, options: options);
  }

  /// Removes all authorization parameters from the connection. This lets you remove the secret from the connection so you can reuse it without having to create a new connection.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeauthorizeConnectionResponse> deauthorizeConnection(
    $0.DeauthorizeConnectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deauthorizeConnection, request, options: options);
  }

  /// Deletes the specified API destination.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteApiDestinationResponse> deleteApiDestination(
    $0.DeleteApiDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteApiDestination, request, options: options);
  }

  /// Deletes the specified archive.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteArchiveResponse> deleteArchive(
    $0.DeleteArchiveRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteArchive, request, options: options);
  }

  /// Deletes a connection.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteConnectionResponse> deleteConnection(
    $0.DeleteConnectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteConnection, request, options: options);
  }

  /// Deletes the specified custom event bus or partner event bus. All rules associated with this event bus need to be deleted. You can't delete your account's default event bus.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteEventBus(
    $0.DeleteEventBusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteEventBus, request, options: options);
  }

  /// This operation is used by SaaS partners to delete a partner event source. This operation is not used by Amazon Web Services customers. When you delete an event source, the status of the correspondi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deletePartnerEventSource(
    $0.DeletePartnerEventSourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePartnerEventSource, request,
        options: options);
  }

  /// Deletes the specified rule. Before you can delete the rule, you must remove all targets, using RemoveTargets. When you delete a rule, incoming events might continue to match to the deleted rule. Al...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteRule(
    $0.DeleteRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRule, request, options: options);
  }

  /// Retrieves details about an API destination.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeApiDestinationResponse>
      describeApiDestination(
    $0.DescribeApiDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeApiDestination, request,
        options: options);
  }

  /// Retrieves details about an archive.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeArchiveResponse> describeArchive(
    $0.DescribeArchiveRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeArchive, request, options: options);
  }

  /// Retrieves details about a connection.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeConnectionResponse> describeConnection(
    $0.DescribeConnectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeConnection, request, options: options);
  }

  /// Displays details about an event bus in your account. This can include the external Amazon Web Services accounts that are permitted to write events to your default event bus, and the associated poli...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeEventBusResponse> describeEventBus(
    $0.DescribeEventBusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeEventBus, request, options: options);
  }

  /// This operation lists details about a partner event source that is shared with your account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeEventSourceResponse> describeEventSource(
    $0.DescribeEventSourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeEventSource, request, options: options);
  }

  /// An SaaS partner can use this operation to list details about a partner event source that they have created. Amazon Web Services customers do not use this operation. Instead, Amazon Web Services cus...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribePartnerEventSourceResponse>
      describePartnerEventSource(
    $0.DescribePartnerEventSourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describePartnerEventSource, request,
        options: options);
  }

  /// Retrieves details about a replay. Use DescribeReplay to determine the progress of a running replay. A replay processes events to replay based on the time in the event, and replays them using 1 minu...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeReplayResponse> describeReplay(
    $0.DescribeReplayRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeReplay, request, options: options);
  }

  /// Describes the specified rule. DescribeRule does not list the targets of a rule. To see the targets associated with a rule, use ListTargetsByRule.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeRuleResponse> describeRule(
    $0.DescribeRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeRule, request, options: options);
  }

  /// Disables the specified rule. A disabled rule won't match any events, and won't self-trigger if it has a schedule expression. When you disable a rule, incoming events might continue to match to the ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> disableRule(
    $0.DisableRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableRule, request, options: options);
  }

  /// Enables the specified rule. If the rule does not exist, the operation fails. When you enable a rule, incoming events might not immediately start matching to a newly enabled rule. Allow a short peri...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> enableRule(
    $0.EnableRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableRule, request, options: options);
  }

  /// Retrieves a list of API destination in the account in the current Region.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListApiDestinationsResponse> listApiDestinations(
    $0.ListApiDestinationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listApiDestinations, request, options: options);
  }

  /// Lists your archives. You can either list all the archives or you can provide a prefix to match to the archive names. Filter parameters are exclusive.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListArchivesResponse> listArchives(
    $0.ListArchivesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listArchives, request, options: options);
  }

  /// Retrieves a list of connections from the account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListConnectionsResponse> listConnections(
    $0.ListConnectionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConnections, request, options: options);
  }

  /// Lists all the event buses in your account, including the default event bus, custom event buses, and partner event buses.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListEventBusesResponse> listEventBuses(
    $0.ListEventBusesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEventBuses, request, options: options);
  }

  /// You can use this to see all the partner event sources that have been shared with your Amazon Web Services account. For more information about partner event sources, see CreateEventBus.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListEventSourcesResponse> listEventSources(
    $0.ListEventSourcesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEventSources, request, options: options);
  }

  /// An SaaS partner can use this operation to display the Amazon Web Services account ID that a particular partner event source name is associated with. This operation is not used by Amazon Web Service...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListPartnerEventSourceAccountsResponse>
      listPartnerEventSourceAccounts(
    $0.ListPartnerEventSourceAccountsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPartnerEventSourceAccounts, request,
        options: options);
  }

  /// An SaaS partner can use this operation to list all the partner event source names that they have created. This operation is not used by Amazon Web Services customers.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListPartnerEventSourcesResponse>
      listPartnerEventSources(
    $0.ListPartnerEventSourcesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPartnerEventSources, request,
        options: options);
  }

  /// Lists your replays. You can either list all the replays or you can provide a prefix to match to the replay names. Filter parameters are exclusive.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListReplaysResponse> listReplays(
    $0.ListReplaysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listReplays, request, options: options);
  }

  /// Lists the rules for the specified target. You can see which of the rules in Amazon EventBridge can invoke a specific target in your account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRuleNamesByTargetResponse> listRuleNamesByTarget(
    $0.ListRuleNamesByTargetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRuleNamesByTarget, request, options: options);
  }

  /// Lists your Amazon EventBridge rules. You can either list all the rules or you can provide a prefix to match to the rule names. ListRules does not list the targets of a rule. To see the targets asso...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListRulesResponse> listRules(
    $0.ListRulesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRules, request, options: options);
  }

  /// Displays the tags associated with an EventBridge resource. In EventBridge, rules and event buses can be tagged.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Lists the targets assigned to the specified rule.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTargetsByRuleResponse> listTargetsByRule(
    $0.ListTargetsByRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTargetsByRule, request, options: options);
  }

  /// Sends custom events to Amazon EventBridge so that they can be matched to rules.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutEventsResponse> putEvents(
    $0.PutEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putEvents, request, options: options);
  }

  /// This is used by SaaS partners to write events to a customer's partner event bus. Amazon Web Services customers do not use this operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutPartnerEventsResponse> putPartnerEvents(
    $0.PutPartnerEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putPartnerEvents, request, options: options);
  }

  /// Running PutPermission permits the specified Amazon Web Services account or Amazon Web Services organization to put events to the specified event bus. Amazon EventBridge (CloudWatch Events) rules in...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putPermission(
    $0.PutPermissionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putPermission, request, options: options);
  }

  /// Creates or updates the specified rule. Rules are enabled by default, or based on value of the state. You can disable a rule using DisableRule. A single rule watches for events from a single event b...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutRuleResponse> putRule(
    $0.PutRuleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putRule, request, options: options);
  }

  /// Adds the specified targets to the specified rule, or updates the targets if they are already associated with the rule. Targets are the resources that are invoked when a rule is triggered. You can c...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutTargetsResponse> putTargets(
    $0.PutTargetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putTargets, request, options: options);
  }

  /// Revokes the permission of another Amazon Web Services account to be able to put events to the specified event bus. Specify the account to revoke by the StatementId value that you associated with th...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> removePermission(
    $0.RemovePermissionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removePermission, request, options: options);
  }

  /// Removes the specified targets from the specified rule. When the rule is triggered, those targets are no longer be invoked. When you remove a target, when the associated rule triggers, removed targe...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RemoveTargetsResponse> removeTargets(
    $0.RemoveTargetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeTargets, request, options: options);
  }

  /// Starts the specified replay. Events are not necessarily replayed in the exact same order that they were added to the archive. A replay processes events to replay based on the time in the event, and...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartReplayResponse> startReplay(
    $0.StartReplayRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startReplay, request, options: options);
  }

  /// Assigns one or more tags (key-value pairs) to the specified EventBridge resource. Tags can help you organize and categorize your resources. You can also use them to scope user permissions by granti...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Tests whether the specified event pattern matches the provided event. Most services in Amazon Web Services treat : or / as the same character in Amazon Resource Names (ARNs). However, EventBridge u...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TestEventPatternResponse> testEventPattern(
    $0.TestEventPatternRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testEventPattern, request, options: options);
  }

  /// Removes one or more tags from the specified EventBridge resource. In Amazon EventBridge (CloudWatch Events), rules and event buses can be tagged.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates an API destination.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateApiDestinationResponse> updateApiDestination(
    $0.UpdateApiDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateApiDestination, request, options: options);
  }

  /// Updates the specified archive.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateArchiveResponse> updateArchive(
    $0.UpdateArchiveRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateArchive, request, options: options);
  }

  /// Updates settings for a connection.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateConnectionResponse> updateConnection(
    $0.UpdateConnectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateConnection, request, options: options);
  }

  // method descriptors

  static final _$activateEventSource =
      $grpc.ClientMethod<$0.ActivateEventSourceRequest, $1.Empty>(
          '/events.CloudWatchEventsService/ActivateEventSource',
          ($0.ActivateEventSourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$cancelReplay =
      $grpc.ClientMethod<$0.CancelReplayRequest, $0.CancelReplayResponse>(
          '/events.CloudWatchEventsService/CancelReplay',
          ($0.CancelReplayRequest value) => value.writeToBuffer(),
          $0.CancelReplayResponse.fromBuffer);
  static final _$createApiDestination = $grpc.ClientMethod<
          $0.CreateApiDestinationRequest, $0.CreateApiDestinationResponse>(
      '/events.CloudWatchEventsService/CreateApiDestination',
      ($0.CreateApiDestinationRequest value) => value.writeToBuffer(),
      $0.CreateApiDestinationResponse.fromBuffer);
  static final _$createArchive =
      $grpc.ClientMethod<$0.CreateArchiveRequest, $0.CreateArchiveResponse>(
          '/events.CloudWatchEventsService/CreateArchive',
          ($0.CreateArchiveRequest value) => value.writeToBuffer(),
          $0.CreateArchiveResponse.fromBuffer);
  static final _$createConnection = $grpc.ClientMethod<
          $0.CreateConnectionRequest, $0.CreateConnectionResponse>(
      '/events.CloudWatchEventsService/CreateConnection',
      ($0.CreateConnectionRequest value) => value.writeToBuffer(),
      $0.CreateConnectionResponse.fromBuffer);
  static final _$createEventBus =
      $grpc.ClientMethod<$0.CreateEventBusRequest, $0.CreateEventBusResponse>(
          '/events.CloudWatchEventsService/CreateEventBus',
          ($0.CreateEventBusRequest value) => value.writeToBuffer(),
          $0.CreateEventBusResponse.fromBuffer);
  static final _$createPartnerEventSource = $grpc.ClientMethod<
          $0.CreatePartnerEventSourceRequest,
          $0.CreatePartnerEventSourceResponse>(
      '/events.CloudWatchEventsService/CreatePartnerEventSource',
      ($0.CreatePartnerEventSourceRequest value) => value.writeToBuffer(),
      $0.CreatePartnerEventSourceResponse.fromBuffer);
  static final _$deactivateEventSource =
      $grpc.ClientMethod<$0.DeactivateEventSourceRequest, $1.Empty>(
          '/events.CloudWatchEventsService/DeactivateEventSource',
          ($0.DeactivateEventSourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deauthorizeConnection = $grpc.ClientMethod<
          $0.DeauthorizeConnectionRequest, $0.DeauthorizeConnectionResponse>(
      '/events.CloudWatchEventsService/DeauthorizeConnection',
      ($0.DeauthorizeConnectionRequest value) => value.writeToBuffer(),
      $0.DeauthorizeConnectionResponse.fromBuffer);
  static final _$deleteApiDestination = $grpc.ClientMethod<
          $0.DeleteApiDestinationRequest, $0.DeleteApiDestinationResponse>(
      '/events.CloudWatchEventsService/DeleteApiDestination',
      ($0.DeleteApiDestinationRequest value) => value.writeToBuffer(),
      $0.DeleteApiDestinationResponse.fromBuffer);
  static final _$deleteArchive =
      $grpc.ClientMethod<$0.DeleteArchiveRequest, $0.DeleteArchiveResponse>(
          '/events.CloudWatchEventsService/DeleteArchive',
          ($0.DeleteArchiveRequest value) => value.writeToBuffer(),
          $0.DeleteArchiveResponse.fromBuffer);
  static final _$deleteConnection = $grpc.ClientMethod<
          $0.DeleteConnectionRequest, $0.DeleteConnectionResponse>(
      '/events.CloudWatchEventsService/DeleteConnection',
      ($0.DeleteConnectionRequest value) => value.writeToBuffer(),
      $0.DeleteConnectionResponse.fromBuffer);
  static final _$deleteEventBus =
      $grpc.ClientMethod<$0.DeleteEventBusRequest, $1.Empty>(
          '/events.CloudWatchEventsService/DeleteEventBus',
          ($0.DeleteEventBusRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deletePartnerEventSource =
      $grpc.ClientMethod<$0.DeletePartnerEventSourceRequest, $1.Empty>(
          '/events.CloudWatchEventsService/DeletePartnerEventSource',
          ($0.DeletePartnerEventSourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRule =
      $grpc.ClientMethod<$0.DeleteRuleRequest, $1.Empty>(
          '/events.CloudWatchEventsService/DeleteRule',
          ($0.DeleteRuleRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$describeApiDestination = $grpc.ClientMethod<
          $0.DescribeApiDestinationRequest, $0.DescribeApiDestinationResponse>(
      '/events.CloudWatchEventsService/DescribeApiDestination',
      ($0.DescribeApiDestinationRequest value) => value.writeToBuffer(),
      $0.DescribeApiDestinationResponse.fromBuffer);
  static final _$describeArchive =
      $grpc.ClientMethod<$0.DescribeArchiveRequest, $0.DescribeArchiveResponse>(
          '/events.CloudWatchEventsService/DescribeArchive',
          ($0.DescribeArchiveRequest value) => value.writeToBuffer(),
          $0.DescribeArchiveResponse.fromBuffer);
  static final _$describeConnection = $grpc.ClientMethod<
          $0.DescribeConnectionRequest, $0.DescribeConnectionResponse>(
      '/events.CloudWatchEventsService/DescribeConnection',
      ($0.DescribeConnectionRequest value) => value.writeToBuffer(),
      $0.DescribeConnectionResponse.fromBuffer);
  static final _$describeEventBus = $grpc.ClientMethod<
          $0.DescribeEventBusRequest, $0.DescribeEventBusResponse>(
      '/events.CloudWatchEventsService/DescribeEventBus',
      ($0.DescribeEventBusRequest value) => value.writeToBuffer(),
      $0.DescribeEventBusResponse.fromBuffer);
  static final _$describeEventSource = $grpc.ClientMethod<
          $0.DescribeEventSourceRequest, $0.DescribeEventSourceResponse>(
      '/events.CloudWatchEventsService/DescribeEventSource',
      ($0.DescribeEventSourceRequest value) => value.writeToBuffer(),
      $0.DescribeEventSourceResponse.fromBuffer);
  static final _$describePartnerEventSource = $grpc.ClientMethod<
          $0.DescribePartnerEventSourceRequest,
          $0.DescribePartnerEventSourceResponse>(
      '/events.CloudWatchEventsService/DescribePartnerEventSource',
      ($0.DescribePartnerEventSourceRequest value) => value.writeToBuffer(),
      $0.DescribePartnerEventSourceResponse.fromBuffer);
  static final _$describeReplay =
      $grpc.ClientMethod<$0.DescribeReplayRequest, $0.DescribeReplayResponse>(
          '/events.CloudWatchEventsService/DescribeReplay',
          ($0.DescribeReplayRequest value) => value.writeToBuffer(),
          $0.DescribeReplayResponse.fromBuffer);
  static final _$describeRule =
      $grpc.ClientMethod<$0.DescribeRuleRequest, $0.DescribeRuleResponse>(
          '/events.CloudWatchEventsService/DescribeRule',
          ($0.DescribeRuleRequest value) => value.writeToBuffer(),
          $0.DescribeRuleResponse.fromBuffer);
  static final _$disableRule =
      $grpc.ClientMethod<$0.DisableRuleRequest, $1.Empty>(
          '/events.CloudWatchEventsService/DisableRule',
          ($0.DisableRuleRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$enableRule =
      $grpc.ClientMethod<$0.EnableRuleRequest, $1.Empty>(
          '/events.CloudWatchEventsService/EnableRule',
          ($0.EnableRuleRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$listApiDestinations = $grpc.ClientMethod<
          $0.ListApiDestinationsRequest, $0.ListApiDestinationsResponse>(
      '/events.CloudWatchEventsService/ListApiDestinations',
      ($0.ListApiDestinationsRequest value) => value.writeToBuffer(),
      $0.ListApiDestinationsResponse.fromBuffer);
  static final _$listArchives =
      $grpc.ClientMethod<$0.ListArchivesRequest, $0.ListArchivesResponse>(
          '/events.CloudWatchEventsService/ListArchives',
          ($0.ListArchivesRequest value) => value.writeToBuffer(),
          $0.ListArchivesResponse.fromBuffer);
  static final _$listConnections =
      $grpc.ClientMethod<$0.ListConnectionsRequest, $0.ListConnectionsResponse>(
          '/events.CloudWatchEventsService/ListConnections',
          ($0.ListConnectionsRequest value) => value.writeToBuffer(),
          $0.ListConnectionsResponse.fromBuffer);
  static final _$listEventBuses =
      $grpc.ClientMethod<$0.ListEventBusesRequest, $0.ListEventBusesResponse>(
          '/events.CloudWatchEventsService/ListEventBuses',
          ($0.ListEventBusesRequest value) => value.writeToBuffer(),
          $0.ListEventBusesResponse.fromBuffer);
  static final _$listEventSources = $grpc.ClientMethod<
          $0.ListEventSourcesRequest, $0.ListEventSourcesResponse>(
      '/events.CloudWatchEventsService/ListEventSources',
      ($0.ListEventSourcesRequest value) => value.writeToBuffer(),
      $0.ListEventSourcesResponse.fromBuffer);
  static final _$listPartnerEventSourceAccounts = $grpc.ClientMethod<
          $0.ListPartnerEventSourceAccountsRequest,
          $0.ListPartnerEventSourceAccountsResponse>(
      '/events.CloudWatchEventsService/ListPartnerEventSourceAccounts',
      ($0.ListPartnerEventSourceAccountsRequest value) => value.writeToBuffer(),
      $0.ListPartnerEventSourceAccountsResponse.fromBuffer);
  static final _$listPartnerEventSources = $grpc.ClientMethod<
          $0.ListPartnerEventSourcesRequest,
          $0.ListPartnerEventSourcesResponse>(
      '/events.CloudWatchEventsService/ListPartnerEventSources',
      ($0.ListPartnerEventSourcesRequest value) => value.writeToBuffer(),
      $0.ListPartnerEventSourcesResponse.fromBuffer);
  static final _$listReplays =
      $grpc.ClientMethod<$0.ListReplaysRequest, $0.ListReplaysResponse>(
          '/events.CloudWatchEventsService/ListReplays',
          ($0.ListReplaysRequest value) => value.writeToBuffer(),
          $0.ListReplaysResponse.fromBuffer);
  static final _$listRuleNamesByTarget = $grpc.ClientMethod<
          $0.ListRuleNamesByTargetRequest, $0.ListRuleNamesByTargetResponse>(
      '/events.CloudWatchEventsService/ListRuleNamesByTarget',
      ($0.ListRuleNamesByTargetRequest value) => value.writeToBuffer(),
      $0.ListRuleNamesByTargetResponse.fromBuffer);
  static final _$listRules =
      $grpc.ClientMethod<$0.ListRulesRequest, $0.ListRulesResponse>(
          '/events.CloudWatchEventsService/ListRules',
          ($0.ListRulesRequest value) => value.writeToBuffer(),
          $0.ListRulesResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/events.CloudWatchEventsService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listTargetsByRule = $grpc.ClientMethod<
          $0.ListTargetsByRuleRequest, $0.ListTargetsByRuleResponse>(
      '/events.CloudWatchEventsService/ListTargetsByRule',
      ($0.ListTargetsByRuleRequest value) => value.writeToBuffer(),
      $0.ListTargetsByRuleResponse.fromBuffer);
  static final _$putEvents =
      $grpc.ClientMethod<$0.PutEventsRequest, $0.PutEventsResponse>(
          '/events.CloudWatchEventsService/PutEvents',
          ($0.PutEventsRequest value) => value.writeToBuffer(),
          $0.PutEventsResponse.fromBuffer);
  static final _$putPartnerEvents = $grpc.ClientMethod<
          $0.PutPartnerEventsRequest, $0.PutPartnerEventsResponse>(
      '/events.CloudWatchEventsService/PutPartnerEvents',
      ($0.PutPartnerEventsRequest value) => value.writeToBuffer(),
      $0.PutPartnerEventsResponse.fromBuffer);
  static final _$putPermission =
      $grpc.ClientMethod<$0.PutPermissionRequest, $1.Empty>(
          '/events.CloudWatchEventsService/PutPermission',
          ($0.PutPermissionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putRule =
      $grpc.ClientMethod<$0.PutRuleRequest, $0.PutRuleResponse>(
          '/events.CloudWatchEventsService/PutRule',
          ($0.PutRuleRequest value) => value.writeToBuffer(),
          $0.PutRuleResponse.fromBuffer);
  static final _$putTargets =
      $grpc.ClientMethod<$0.PutTargetsRequest, $0.PutTargetsResponse>(
          '/events.CloudWatchEventsService/PutTargets',
          ($0.PutTargetsRequest value) => value.writeToBuffer(),
          $0.PutTargetsResponse.fromBuffer);
  static final _$removePermission =
      $grpc.ClientMethod<$0.RemovePermissionRequest, $1.Empty>(
          '/events.CloudWatchEventsService/RemovePermission',
          ($0.RemovePermissionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$removeTargets =
      $grpc.ClientMethod<$0.RemoveTargetsRequest, $0.RemoveTargetsResponse>(
          '/events.CloudWatchEventsService/RemoveTargets',
          ($0.RemoveTargetsRequest value) => value.writeToBuffer(),
          $0.RemoveTargetsResponse.fromBuffer);
  static final _$startReplay =
      $grpc.ClientMethod<$0.StartReplayRequest, $0.StartReplayResponse>(
          '/events.CloudWatchEventsService/StartReplay',
          ($0.StartReplayRequest value) => value.writeToBuffer(),
          $0.StartReplayResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/events.CloudWatchEventsService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$testEventPattern = $grpc.ClientMethod<
          $0.TestEventPatternRequest, $0.TestEventPatternResponse>(
      '/events.CloudWatchEventsService/TestEventPattern',
      ($0.TestEventPatternRequest value) => value.writeToBuffer(),
      $0.TestEventPatternResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/events.CloudWatchEventsService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateApiDestination = $grpc.ClientMethod<
          $0.UpdateApiDestinationRequest, $0.UpdateApiDestinationResponse>(
      '/events.CloudWatchEventsService/UpdateApiDestination',
      ($0.UpdateApiDestinationRequest value) => value.writeToBuffer(),
      $0.UpdateApiDestinationResponse.fromBuffer);
  static final _$updateArchive =
      $grpc.ClientMethod<$0.UpdateArchiveRequest, $0.UpdateArchiveResponse>(
          '/events.CloudWatchEventsService/UpdateArchive',
          ($0.UpdateArchiveRequest value) => value.writeToBuffer(),
          $0.UpdateArchiveResponse.fromBuffer);
  static final _$updateConnection = $grpc.ClientMethod<
          $0.UpdateConnectionRequest, $0.UpdateConnectionResponse>(
      '/events.CloudWatchEventsService/UpdateConnection',
      ($0.UpdateConnectionRequest value) => value.writeToBuffer(),
      $0.UpdateConnectionResponse.fromBuffer);
}

@$pb.GrpcServiceName('events.CloudWatchEventsService')
abstract class CloudWatchEventsServiceBase extends $grpc.Service {
  $core.String get $name => 'events.CloudWatchEventsService';

  CloudWatchEventsServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ActivateEventSourceRequest, $1.Empty>(
        'ActivateEventSource',
        activateEventSource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ActivateEventSourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CancelReplayRequest, $0.CancelReplayResponse>(
            'CancelReplay',
            cancelReplay_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CancelReplayRequest.fromBuffer(value),
            ($0.CancelReplayResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateApiDestinationRequest,
            $0.CreateApiDestinationResponse>(
        'CreateApiDestination',
        createApiDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateApiDestinationRequest.fromBuffer(value),
        ($0.CreateApiDestinationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateArchiveRequest, $0.CreateArchiveResponse>(
            'CreateArchive',
            createArchive_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateArchiveRequest.fromBuffer(value),
            ($0.CreateArchiveResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateConnectionRequest,
            $0.CreateConnectionResponse>(
        'CreateConnection',
        createConnection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateConnectionRequest.fromBuffer(value),
        ($0.CreateConnectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateEventBusRequest,
            $0.CreateEventBusResponse>(
        'CreateEventBus',
        createEventBus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateEventBusRequest.fromBuffer(value),
        ($0.CreateEventBusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePartnerEventSourceRequest,
            $0.CreatePartnerEventSourceResponse>(
        'CreatePartnerEventSource',
        createPartnerEventSource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePartnerEventSourceRequest.fromBuffer(value),
        ($0.CreatePartnerEventSourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeactivateEventSourceRequest, $1.Empty>(
        'DeactivateEventSource',
        deactivateEventSource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeactivateEventSourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeauthorizeConnectionRequest,
            $0.DeauthorizeConnectionResponse>(
        'DeauthorizeConnection',
        deauthorizeConnection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeauthorizeConnectionRequest.fromBuffer(value),
        ($0.DeauthorizeConnectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteApiDestinationRequest,
            $0.DeleteApiDestinationResponse>(
        'DeleteApiDestination',
        deleteApiDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteApiDestinationRequest.fromBuffer(value),
        ($0.DeleteApiDestinationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteArchiveRequest, $0.DeleteArchiveResponse>(
            'DeleteArchive',
            deleteArchive_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteArchiveRequest.fromBuffer(value),
            ($0.DeleteArchiveResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteConnectionRequest,
            $0.DeleteConnectionResponse>(
        'DeleteConnection',
        deleteConnection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteConnectionRequest.fromBuffer(value),
        ($0.DeleteConnectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteEventBusRequest, $1.Empty>(
        'DeleteEventBus',
        deleteEventBus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteEventBusRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeletePartnerEventSourceRequest, $1.Empty>(
            'DeletePartnerEventSource',
            deletePartnerEventSource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeletePartnerEventSourceRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRuleRequest, $1.Empty>(
        'DeleteRule',
        deleteRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteRuleRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeApiDestinationRequest,
            $0.DescribeApiDestinationResponse>(
        'DescribeApiDestination',
        describeApiDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeApiDestinationRequest.fromBuffer(value),
        ($0.DescribeApiDestinationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeArchiveRequest,
            $0.DescribeArchiveResponse>(
        'DescribeArchive',
        describeArchive_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeArchiveRequest.fromBuffer(value),
        ($0.DescribeArchiveResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeConnectionRequest,
            $0.DescribeConnectionResponse>(
        'DescribeConnection',
        describeConnection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeConnectionRequest.fromBuffer(value),
        ($0.DescribeConnectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeEventBusRequest,
            $0.DescribeEventBusResponse>(
        'DescribeEventBus',
        describeEventBus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeEventBusRequest.fromBuffer(value),
        ($0.DescribeEventBusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeEventSourceRequest,
            $0.DescribeEventSourceResponse>(
        'DescribeEventSource',
        describeEventSource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeEventSourceRequest.fromBuffer(value),
        ($0.DescribeEventSourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribePartnerEventSourceRequest,
            $0.DescribePartnerEventSourceResponse>(
        'DescribePartnerEventSource',
        describePartnerEventSource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribePartnerEventSourceRequest.fromBuffer(value),
        ($0.DescribePartnerEventSourceResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeReplayRequest,
            $0.DescribeReplayResponse>(
        'DescribeReplay',
        describeReplay_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeReplayRequest.fromBuffer(value),
        ($0.DescribeReplayResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeRuleRequest, $0.DescribeRuleResponse>(
            'DescribeRule',
            describeRule_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeRuleRequest.fromBuffer(value),
            ($0.DescribeRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableRuleRequest, $1.Empty>(
        'DisableRule',
        disableRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableRuleRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableRuleRequest, $1.Empty>(
        'EnableRule',
        enableRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.EnableRuleRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListApiDestinationsRequest,
            $0.ListApiDestinationsResponse>(
        'ListApiDestinations',
        listApiDestinations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListApiDestinationsRequest.fromBuffer(value),
        ($0.ListApiDestinationsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListArchivesRequest, $0.ListArchivesResponse>(
            'ListArchives',
            listArchives_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListArchivesRequest.fromBuffer(value),
            ($0.ListArchivesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListConnectionsRequest,
            $0.ListConnectionsResponse>(
        'ListConnections',
        listConnections_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListConnectionsRequest.fromBuffer(value),
        ($0.ListConnectionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEventBusesRequest,
            $0.ListEventBusesResponse>(
        'ListEventBuses',
        listEventBuses_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEventBusesRequest.fromBuffer(value),
        ($0.ListEventBusesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEventSourcesRequest,
            $0.ListEventSourcesResponse>(
        'ListEventSources',
        listEventSources_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEventSourcesRequest.fromBuffer(value),
        ($0.ListEventSourcesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPartnerEventSourceAccountsRequest,
            $0.ListPartnerEventSourceAccountsResponse>(
        'ListPartnerEventSourceAccounts',
        listPartnerEventSourceAccounts_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPartnerEventSourceAccountsRequest.fromBuffer(value),
        ($0.ListPartnerEventSourceAccountsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPartnerEventSourcesRequest,
            $0.ListPartnerEventSourcesResponse>(
        'ListPartnerEventSources',
        listPartnerEventSources_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPartnerEventSourcesRequest.fromBuffer(value),
        ($0.ListPartnerEventSourcesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListReplaysRequest, $0.ListReplaysResponse>(
            'ListReplays',
            listReplays_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListReplaysRequest.fromBuffer(value),
            ($0.ListReplaysResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRuleNamesByTargetRequest,
            $0.ListRuleNamesByTargetResponse>(
        'ListRuleNamesByTarget',
        listRuleNamesByTarget_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRuleNamesByTargetRequest.fromBuffer(value),
        ($0.ListRuleNamesByTargetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRulesRequest, $0.ListRulesResponse>(
        'ListRules',
        listRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListRulesRequest.fromBuffer(value),
        ($0.ListRulesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTargetsByRuleRequest,
            $0.ListTargetsByRuleResponse>(
        'ListTargetsByRule',
        listTargetsByRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTargetsByRuleRequest.fromBuffer(value),
        ($0.ListTargetsByRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutEventsRequest, $0.PutEventsResponse>(
        'PutEvents',
        putEvents_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutEventsRequest.fromBuffer(value),
        ($0.PutEventsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutPartnerEventsRequest,
            $0.PutPartnerEventsResponse>(
        'PutPartnerEvents',
        putPartnerEvents_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutPartnerEventsRequest.fromBuffer(value),
        ($0.PutPartnerEventsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutPermissionRequest, $1.Empty>(
        'PutPermission',
        putPermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutPermissionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutRuleRequest, $0.PutRuleResponse>(
        'PutRule',
        putRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutRuleRequest.fromBuffer(value),
        ($0.PutRuleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutTargetsRequest, $0.PutTargetsResponse>(
        'PutTargets',
        putTargets_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutTargetsRequest.fromBuffer(value),
        ($0.PutTargetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemovePermissionRequest, $1.Empty>(
        'RemovePermission',
        removePermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemovePermissionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RemoveTargetsRequest, $0.RemoveTargetsResponse>(
            'RemoveTargets',
            removeTargets_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RemoveTargetsRequest.fromBuffer(value),
            ($0.RemoveTargetsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartReplayRequest, $0.StartReplayResponse>(
            'StartReplay',
            startReplay_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartReplayRequest.fromBuffer(value),
            ($0.StartReplayResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
            'TagResource',
            tagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TagResourceRequest.fromBuffer(value),
            ($0.TagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestEventPatternRequest,
            $0.TestEventPatternResponse>(
        'TestEventPattern',
        testEventPattern_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TestEventPatternRequest.fromBuffer(value),
        ($0.TestEventPatternResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceRequest.fromBuffer(value),
            ($0.UntagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateApiDestinationRequest,
            $0.UpdateApiDestinationResponse>(
        'UpdateApiDestination',
        updateApiDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateApiDestinationRequest.fromBuffer(value),
        ($0.UpdateApiDestinationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateArchiveRequest, $0.UpdateArchiveResponse>(
            'UpdateArchive',
            updateArchive_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateArchiveRequest.fromBuffer(value),
            ($0.UpdateArchiveResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateConnectionRequest,
            $0.UpdateConnectionResponse>(
        'UpdateConnection',
        updateConnection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateConnectionRequest.fromBuffer(value),
        ($0.UpdateConnectionResponse value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> activateEventSource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ActivateEventSourceRequest> $request) async {
    return activateEventSource($call, await $request);
  }

  $async.Future<$1.Empty> activateEventSource(
      $grpc.ServiceCall call, $0.ActivateEventSourceRequest request);

  $async.Future<$0.CancelReplayResponse> cancelReplay_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelReplayRequest> $request) async {
    return cancelReplay($call, await $request);
  }

  $async.Future<$0.CancelReplayResponse> cancelReplay(
      $grpc.ServiceCall call, $0.CancelReplayRequest request);

  $async.Future<$0.CreateApiDestinationResponse> createApiDestination_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateApiDestinationRequest> $request) async {
    return createApiDestination($call, await $request);
  }

  $async.Future<$0.CreateApiDestinationResponse> createApiDestination(
      $grpc.ServiceCall call, $0.CreateApiDestinationRequest request);

  $async.Future<$0.CreateArchiveResponse> createArchive_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateArchiveRequest> $request) async {
    return createArchive($call, await $request);
  }

  $async.Future<$0.CreateArchiveResponse> createArchive(
      $grpc.ServiceCall call, $0.CreateArchiveRequest request);

  $async.Future<$0.CreateConnectionResponse> createConnection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateConnectionRequest> $request) async {
    return createConnection($call, await $request);
  }

  $async.Future<$0.CreateConnectionResponse> createConnection(
      $grpc.ServiceCall call, $0.CreateConnectionRequest request);

  $async.Future<$0.CreateEventBusResponse> createEventBus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateEventBusRequest> $request) async {
    return createEventBus($call, await $request);
  }

  $async.Future<$0.CreateEventBusResponse> createEventBus(
      $grpc.ServiceCall call, $0.CreateEventBusRequest request);

  $async.Future<$0.CreatePartnerEventSourceResponse>
      createPartnerEventSource_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreatePartnerEventSourceRequest> $request) async {
    return createPartnerEventSource($call, await $request);
  }

  $async.Future<$0.CreatePartnerEventSourceResponse> createPartnerEventSource(
      $grpc.ServiceCall call, $0.CreatePartnerEventSourceRequest request);

  $async.Future<$1.Empty> deactivateEventSource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeactivateEventSourceRequest> $request) async {
    return deactivateEventSource($call, await $request);
  }

  $async.Future<$1.Empty> deactivateEventSource(
      $grpc.ServiceCall call, $0.DeactivateEventSourceRequest request);

  $async.Future<$0.DeauthorizeConnectionResponse> deauthorizeConnection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeauthorizeConnectionRequest> $request) async {
    return deauthorizeConnection($call, await $request);
  }

  $async.Future<$0.DeauthorizeConnectionResponse> deauthorizeConnection(
      $grpc.ServiceCall call, $0.DeauthorizeConnectionRequest request);

  $async.Future<$0.DeleteApiDestinationResponse> deleteApiDestination_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteApiDestinationRequest> $request) async {
    return deleteApiDestination($call, await $request);
  }

  $async.Future<$0.DeleteApiDestinationResponse> deleteApiDestination(
      $grpc.ServiceCall call, $0.DeleteApiDestinationRequest request);

  $async.Future<$0.DeleteArchiveResponse> deleteArchive_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteArchiveRequest> $request) async {
    return deleteArchive($call, await $request);
  }

  $async.Future<$0.DeleteArchiveResponse> deleteArchive(
      $grpc.ServiceCall call, $0.DeleteArchiveRequest request);

  $async.Future<$0.DeleteConnectionResponse> deleteConnection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteConnectionRequest> $request) async {
    return deleteConnection($call, await $request);
  }

  $async.Future<$0.DeleteConnectionResponse> deleteConnection(
      $grpc.ServiceCall call, $0.DeleteConnectionRequest request);

  $async.Future<$1.Empty> deleteEventBus_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteEventBusRequest> $request) async {
    return deleteEventBus($call, await $request);
  }

  $async.Future<$1.Empty> deleteEventBus(
      $grpc.ServiceCall call, $0.DeleteEventBusRequest request);

  $async.Future<$1.Empty> deletePartnerEventSource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeletePartnerEventSourceRequest> $request) async {
    return deletePartnerEventSource($call, await $request);
  }

  $async.Future<$1.Empty> deletePartnerEventSource(
      $grpc.ServiceCall call, $0.DeletePartnerEventSourceRequest request);

  $async.Future<$1.Empty> deleteRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRuleRequest> $request) async {
    return deleteRule($call, await $request);
  }

  $async.Future<$1.Empty> deleteRule(
      $grpc.ServiceCall call, $0.DeleteRuleRequest request);

  $async.Future<$0.DescribeApiDestinationResponse> describeApiDestination_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeApiDestinationRequest> $request) async {
    return describeApiDestination($call, await $request);
  }

  $async.Future<$0.DescribeApiDestinationResponse> describeApiDestination(
      $grpc.ServiceCall call, $0.DescribeApiDestinationRequest request);

  $async.Future<$0.DescribeArchiveResponse> describeArchive_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeArchiveRequest> $request) async {
    return describeArchive($call, await $request);
  }

  $async.Future<$0.DescribeArchiveResponse> describeArchive(
      $grpc.ServiceCall call, $0.DescribeArchiveRequest request);

  $async.Future<$0.DescribeConnectionResponse> describeConnection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeConnectionRequest> $request) async {
    return describeConnection($call, await $request);
  }

  $async.Future<$0.DescribeConnectionResponse> describeConnection(
      $grpc.ServiceCall call, $0.DescribeConnectionRequest request);

  $async.Future<$0.DescribeEventBusResponse> describeEventBus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeEventBusRequest> $request) async {
    return describeEventBus($call, await $request);
  }

  $async.Future<$0.DescribeEventBusResponse> describeEventBus(
      $grpc.ServiceCall call, $0.DescribeEventBusRequest request);

  $async.Future<$0.DescribeEventSourceResponse> describeEventSource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeEventSourceRequest> $request) async {
    return describeEventSource($call, await $request);
  }

  $async.Future<$0.DescribeEventSourceResponse> describeEventSource(
      $grpc.ServiceCall call, $0.DescribeEventSourceRequest request);

  $async.Future<$0.DescribePartnerEventSourceResponse>
      describePartnerEventSource_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribePartnerEventSourceRequest> $request) async {
    return describePartnerEventSource($call, await $request);
  }

  $async.Future<$0.DescribePartnerEventSourceResponse>
      describePartnerEventSource(
          $grpc.ServiceCall call, $0.DescribePartnerEventSourceRequest request);

  $async.Future<$0.DescribeReplayResponse> describeReplay_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeReplayRequest> $request) async {
    return describeReplay($call, await $request);
  }

  $async.Future<$0.DescribeReplayResponse> describeReplay(
      $grpc.ServiceCall call, $0.DescribeReplayRequest request);

  $async.Future<$0.DescribeRuleResponse> describeRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeRuleRequest> $request) async {
    return describeRule($call, await $request);
  }

  $async.Future<$0.DescribeRuleResponse> describeRule(
      $grpc.ServiceCall call, $0.DescribeRuleRequest request);

  $async.Future<$1.Empty> disableRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DisableRuleRequest> $request) async {
    return disableRule($call, await $request);
  }

  $async.Future<$1.Empty> disableRule(
      $grpc.ServiceCall call, $0.DisableRuleRequest request);

  $async.Future<$1.Empty> enableRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EnableRuleRequest> $request) async {
    return enableRule($call, await $request);
  }

  $async.Future<$1.Empty> enableRule(
      $grpc.ServiceCall call, $0.EnableRuleRequest request);

  $async.Future<$0.ListApiDestinationsResponse> listApiDestinations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListApiDestinationsRequest> $request) async {
    return listApiDestinations($call, await $request);
  }

  $async.Future<$0.ListApiDestinationsResponse> listApiDestinations(
      $grpc.ServiceCall call, $0.ListApiDestinationsRequest request);

  $async.Future<$0.ListArchivesResponse> listArchives_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListArchivesRequest> $request) async {
    return listArchives($call, await $request);
  }

  $async.Future<$0.ListArchivesResponse> listArchives(
      $grpc.ServiceCall call, $0.ListArchivesRequest request);

  $async.Future<$0.ListConnectionsResponse> listConnections_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConnectionsRequest> $request) async {
    return listConnections($call, await $request);
  }

  $async.Future<$0.ListConnectionsResponse> listConnections(
      $grpc.ServiceCall call, $0.ListConnectionsRequest request);

  $async.Future<$0.ListEventBusesResponse> listEventBuses_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEventBusesRequest> $request) async {
    return listEventBuses($call, await $request);
  }

  $async.Future<$0.ListEventBusesResponse> listEventBuses(
      $grpc.ServiceCall call, $0.ListEventBusesRequest request);

  $async.Future<$0.ListEventSourcesResponse> listEventSources_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEventSourcesRequest> $request) async {
    return listEventSources($call, await $request);
  }

  $async.Future<$0.ListEventSourcesResponse> listEventSources(
      $grpc.ServiceCall call, $0.ListEventSourcesRequest request);

  $async.Future<$0.ListPartnerEventSourceAccountsResponse>
      listPartnerEventSourceAccounts_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListPartnerEventSourceAccountsRequest>
              $request) async {
    return listPartnerEventSourceAccounts($call, await $request);
  }

  $async.Future<$0.ListPartnerEventSourceAccountsResponse>
      listPartnerEventSourceAccounts($grpc.ServiceCall call,
          $0.ListPartnerEventSourceAccountsRequest request);

  $async.Future<$0.ListPartnerEventSourcesResponse> listPartnerEventSources_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPartnerEventSourcesRequest> $request) async {
    return listPartnerEventSources($call, await $request);
  }

  $async.Future<$0.ListPartnerEventSourcesResponse> listPartnerEventSources(
      $grpc.ServiceCall call, $0.ListPartnerEventSourcesRequest request);

  $async.Future<$0.ListReplaysResponse> listReplays_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListReplaysRequest> $request) async {
    return listReplays($call, await $request);
  }

  $async.Future<$0.ListReplaysResponse> listReplays(
      $grpc.ServiceCall call, $0.ListReplaysRequest request);

  $async.Future<$0.ListRuleNamesByTargetResponse> listRuleNamesByTarget_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRuleNamesByTargetRequest> $request) async {
    return listRuleNamesByTarget($call, await $request);
  }

  $async.Future<$0.ListRuleNamesByTargetResponse> listRuleNamesByTarget(
      $grpc.ServiceCall call, $0.ListRuleNamesByTargetRequest request);

  $async.Future<$0.ListRulesResponse> listRules_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListRulesRequest> $request) async {
    return listRules($call, await $request);
  }

  $async.Future<$0.ListRulesResponse> listRules(
      $grpc.ServiceCall call, $0.ListRulesRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListTargetsByRuleResponse> listTargetsByRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTargetsByRuleRequest> $request) async {
    return listTargetsByRule($call, await $request);
  }

  $async.Future<$0.ListTargetsByRuleResponse> listTargetsByRule(
      $grpc.ServiceCall call, $0.ListTargetsByRuleRequest request);

  $async.Future<$0.PutEventsResponse> putEvents_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutEventsRequest> $request) async {
    return putEvents($call, await $request);
  }

  $async.Future<$0.PutEventsResponse> putEvents(
      $grpc.ServiceCall call, $0.PutEventsRequest request);

  $async.Future<$0.PutPartnerEventsResponse> putPartnerEvents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutPartnerEventsRequest> $request) async {
    return putPartnerEvents($call, await $request);
  }

  $async.Future<$0.PutPartnerEventsResponse> putPartnerEvents(
      $grpc.ServiceCall call, $0.PutPartnerEventsRequest request);

  $async.Future<$1.Empty> putPermission_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutPermissionRequest> $request) async {
    return putPermission($call, await $request);
  }

  $async.Future<$1.Empty> putPermission(
      $grpc.ServiceCall call, $0.PutPermissionRequest request);

  $async.Future<$0.PutRuleResponse> putRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutRuleRequest> $request) async {
    return putRule($call, await $request);
  }

  $async.Future<$0.PutRuleResponse> putRule(
      $grpc.ServiceCall call, $0.PutRuleRequest request);

  $async.Future<$0.PutTargetsResponse> putTargets_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutTargetsRequest> $request) async {
    return putTargets($call, await $request);
  }

  $async.Future<$0.PutTargetsResponse> putTargets(
      $grpc.ServiceCall call, $0.PutTargetsRequest request);

  $async.Future<$1.Empty> removePermission_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemovePermissionRequest> $request) async {
    return removePermission($call, await $request);
  }

  $async.Future<$1.Empty> removePermission(
      $grpc.ServiceCall call, $0.RemovePermissionRequest request);

  $async.Future<$0.RemoveTargetsResponse> removeTargets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RemoveTargetsRequest> $request) async {
    return removeTargets($call, await $request);
  }

  $async.Future<$0.RemoveTargetsResponse> removeTargets(
      $grpc.ServiceCall call, $0.RemoveTargetsRequest request);

  $async.Future<$0.StartReplayResponse> startReplay_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StartReplayRequest> $request) async {
    return startReplay($call, await $request);
  }

  $async.Future<$0.StartReplayResponse> startReplay(
      $grpc.ServiceCall call, $0.StartReplayRequest request);

  $async.Future<$0.TagResourceResponse> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceResponse> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.TestEventPatternResponse> testEventPattern_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestEventPatternRequest> $request) async {
    return testEventPattern($call, await $request);
  }

  $async.Future<$0.TestEventPatternResponse> testEventPattern(
      $grpc.ServiceCall call, $0.TestEventPatternRequest request);

  $async.Future<$0.UntagResourceResponse> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceResponse> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.UpdateApiDestinationResponse> updateApiDestination_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateApiDestinationRequest> $request) async {
    return updateApiDestination($call, await $request);
  }

  $async.Future<$0.UpdateApiDestinationResponse> updateApiDestination(
      $grpc.ServiceCall call, $0.UpdateApiDestinationRequest request);

  $async.Future<$0.UpdateArchiveResponse> updateArchive_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateArchiveRequest> $request) async {
    return updateArchive($call, await $request);
  }

  $async.Future<$0.UpdateArchiveResponse> updateArchive(
      $grpc.ServiceCall call, $0.UpdateArchiveRequest request);

  $async.Future<$0.UpdateConnectionResponse> updateConnection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateConnectionRequest> $request) async {
    return updateConnection($call, await $request);
  }

  $async.Future<$0.UpdateConnectionResponse> updateConnection(
      $grpc.ServiceCall call, $0.UpdateConnectionRequest request);
}
