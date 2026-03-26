// This is a generated file - do not edit.
//
// Generated from states.proto.

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

import 'states.pb.dart' as $0;

export 'states.pb.dart';

/// SFNService provides states API operations.
@$pb.GrpcServiceName('states.SFNService')
class SFNServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SFNServiceClient(super.channel, {super.options, super.interceptors});

  /// Creates an activity. An activity is a task that you write in any programming language and host on any machine that has access to Step Functions. Activities must poll Step Functions using the GetAct...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateActivityOutput> createActivity(
    $0.CreateActivityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createActivity, request, options: options);
  }

  /// Creates a state machine. A state machine consists of a collection of states that can do work (Task states), determine to which states to transition next (Choice states), stop an execution with an e...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateStateMachineOutput> createStateMachine(
    $0.CreateStateMachineInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createStateMachine, request, options: options);
  }

  /// Creates an alias for a state machine that points to one or two versions of the same state machine. You can set your application to call StartExecution with an alias and update the version the alias...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateStateMachineAliasOutput>
      createStateMachineAlias(
    $0.CreateStateMachineAliasInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createStateMachineAlias, request,
        options: options);
  }

  /// Deletes an activity.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteActivityOutput> deleteActivity(
    $0.DeleteActivityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteActivity, request, options: options);
  }

  /// Deletes a state machine. This is an asynchronous operation. It sets the state machine's status to DELETING and begins the deletion process. A state machine is deleted only when all its executions a...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteStateMachineOutput> deleteStateMachine(
    $0.DeleteStateMachineInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteStateMachine, request, options: options);
  }

  /// Deletes a state machine alias. After you delete a state machine alias, you can't use it to start executions. When you delete a state machine alias, Step Functions doesn't delete the state machine v...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteStateMachineAliasOutput>
      deleteStateMachineAlias(
    $0.DeleteStateMachineAliasInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteStateMachineAlias, request,
        options: options);
  }

  /// Deletes a state machine version. After you delete a version, you can't call StartExecution using that version's ARN or use the version with a state machine alias. Deleting a state machine version w...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteStateMachineVersionOutput>
      deleteStateMachineVersion(
    $0.DeleteStateMachineVersionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteStateMachineVersion, request,
        options: options);
  }

  /// Describes an activity. This operation is eventually consistent. The results are best effort and may not reflect very recent updates and changes.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeActivityOutput> describeActivity(
    $0.DescribeActivityInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeActivity, request, options: options);
  }

  /// Provides information about a state machine execution, such as the state machine associated with the execution, the execution input and output, and relevant execution metadata. If you've redriven an...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeExecutionOutput> describeExecution(
    $0.DescribeExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeExecution, request, options: options);
  }

  /// Provides information about a Map Run's configuration, progress, and results. If you've redriven a Map Run, this API action also returns information about the redrives of that Map Run. For more info...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeMapRunOutput> describeMapRun(
    $0.DescribeMapRunInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMapRun, request, options: options);
  }

  /// Provides information about a state machine's definition, its IAM role Amazon Resource Name (ARN), and configuration. A qualified state machine ARN can either refer to a Distributed Map state define...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeStateMachineOutput> describeStateMachine(
    $0.DescribeStateMachineInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeStateMachine, request, options: options);
  }

  /// Returns details about a state machine alias. Related operations: CreateStateMachineAlias ListStateMachineAliases UpdateStateMachineAlias DeleteStateMachineAlias
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeStateMachineAliasOutput>
      describeStateMachineAlias(
    $0.DescribeStateMachineAliasInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeStateMachineAlias, request,
        options: options);
  }

  /// Provides information about a state machine's definition, its execution role ARN, and configuration. If a Map Run dispatched the execution, this action returns the Map Run Amazon Resource Name (ARN)...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeStateMachineForExecutionOutput>
      describeStateMachineForExecution(
    $0.DescribeStateMachineForExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeStateMachineForExecution, request,
        options: options);
  }

  /// Used by workers to retrieve a task (with the specified activity ARN) which has been scheduled for execution by a running state machine. This initiates a long poll, where the service holds the HTTP ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetActivityTaskOutput> getActivityTask(
    $0.GetActivityTaskInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getActivityTask, request, options: options);
  }

  /// Returns the history of the specified execution as a list of events. By default, the results are returned in ascending order of the timeStamp of the events. Use the reverseOrder parameter to get the...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetExecutionHistoryOutput> getExecutionHistory(
    $0.GetExecutionHistoryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getExecutionHistory, request, options: options);
  }

  /// Lists the existing activities. If nextToken is returned, there are more results available. The value of nextToken is a unique pagination token for each page. Make the call again using the returned ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListActivitiesOutput> listActivities(
    $0.ListActivitiesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listActivities, request, options: options);
  }

  /// Lists all executions of a state machine or a Map Run. You can list all executions related to a state machine by specifying a state machine Amazon Resource Name (ARN), or those related to a Map Run ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListExecutionsOutput> listExecutions(
    $0.ListExecutionsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listExecutions, request, options: options);
  }

  /// Lists all Map Runs that were started by a given state machine execution. Use this API action to obtain Map Run ARNs, and then call DescribeMapRun to obtain more information, if needed.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListMapRunsOutput> listMapRuns(
    $0.ListMapRunsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMapRuns, request, options: options);
  }

  /// Lists aliases for a specified state machine ARN. Results are sorted by time, with the most recently created aliases listed first. To list aliases that reference a state machine version, you can spe...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListStateMachineAliasesOutput>
      listStateMachineAliases(
    $0.ListStateMachineAliasesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listStateMachineAliases, request,
        options: options);
  }

  /// Lists the existing state machines. If nextToken is returned, there are more results available. The value of nextToken is a unique pagination token for each page. Make the call again using the retur...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListStateMachinesOutput> listStateMachines(
    $0.ListStateMachinesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listStateMachines, request, options: options);
  }

  /// Lists versions for the specified state machine Amazon Resource Name (ARN). The results are sorted in descending order of the version creation time. If nextToken is returned, there are more results ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListStateMachineVersionsOutput>
      listStateMachineVersions(
    $0.ListStateMachineVersionsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listStateMachineVersions, request,
        options: options);
  }

  /// List tags for a given resource. Tags may only contain Unicode letters, digits, white space, or these symbols: _ . : / = + - @.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListTagsForResourceOutput> listTagsForResource(
    $0.ListTagsForResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Creates a version from the current revision of a state machine. Use versions to create immutable snapshots of your state machine. You can start executions from versions either directly or with an a...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PublishStateMachineVersionOutput>
      publishStateMachineVersion(
    $0.PublishStateMachineVersionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$publishStateMachineVersion, request,
        options: options);
  }

  /// Restarts unsuccessful executions of Standard workflows that didn't complete successfully in the last 14 days. These include failed, aborted, or timed out executions. When you redrive an execution, ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.RedriveExecutionOutput> redriveExecution(
    $0.RedriveExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$redriveExecution, request, options: options);
  }

  /// Used by activity workers, Task states using the callback pattern, and optionally Task states using the job run pattern to report that the task identified by the taskToken failed. For an execution w...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.SendTaskFailureOutput> sendTaskFailure(
    $0.SendTaskFailureInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendTaskFailure, request, options: options);
  }

  /// Used by activity workers and Task states using the callback pattern, and optionally Task states using the job run pattern to report to Step Functions that the task represented by the specified task...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.SendTaskHeartbeatOutput> sendTaskHeartbeat(
    $0.SendTaskHeartbeatInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendTaskHeartbeat, request, options: options);
  }

  /// Used by activity workers, Task states using the callback pattern, and optionally Task states using the job run pattern to report that the task identified by the taskToken completed successfully.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.SendTaskSuccessOutput> sendTaskSuccess(
    $0.SendTaskSuccessInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendTaskSuccess, request, options: options);
  }

  /// Starts a state machine execution. A qualified state machine ARN can either refer to a Distributed Map state defined within a state machine, a version ARN, or an alias ARN. The following are some ex...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.StartExecutionOutput> startExecution(
    $0.StartExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startExecution, request, options: options);
  }

  /// Starts a Synchronous Express state machine execution. StartSyncExecution is not available for STANDARD workflows. StartSyncExecution will return a 200 OK response, even if your execution fails, bec...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.StartSyncExecutionOutput> startSyncExecution(
    $0.StartSyncExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startSyncExecution, request, options: options);
  }

  /// Stops an execution. This API action is not supported by EXPRESS state machines. For an execution with encryption enabled, Step Functions will encrypt the error and cause fields using the KMS key fo...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.StopExecutionOutput> stopExecution(
    $0.StopExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopExecution, request, options: options);
  }

  /// Add a tag to a Step Functions resource. An array of key-value pairs. For more information, see Using Cost Allocation Tags in the Amazon Web Services Billing and Cost Management User Guide, and Cont...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.TagResourceOutput> tagResource(
    $0.TagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Accepts the definition of a single state and executes it. You can test a state without creating a state machine or updating an existing state machine. Using this API, you can test the following: A ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.TestStateOutput> testState(
    $0.TestStateInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testState, request, options: options);
  }

  /// Remove a tag from a Step Functions resource
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UntagResourceOutput> untagResource(
    $0.UntagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates an in-progress Map Run's configuration to include changes to the settings that control maximum concurrency and Map Run failure.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateMapRunOutput> updateMapRun(
    $0.UpdateMapRunInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMapRun, request, options: options);
  }

  /// Updates an existing state machine by modifying its definition, roleArn, loggingConfiguration, or EncryptionConfiguration. Running executions will continue to use the previous definition and roleArn...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateStateMachineOutput> updateStateMachine(
    $0.UpdateStateMachineInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateStateMachine, request, options: options);
  }

  /// Updates the configuration of an existing state machine alias by modifying its description or routingConfiguration. You must specify at least one of the description or routingConfiguration parameter...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateStateMachineAliasOutput>
      updateStateMachineAlias(
    $0.UpdateStateMachineAliasInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateStateMachineAlias, request,
        options: options);
  }

  /// Validates the syntax of a state machine definition specified in Amazon States Language (ASL), a JSON-based, structured language. You can validate that a state machine definition is correct without ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ValidateStateMachineDefinitionOutput>
      validateStateMachineDefinition(
    $0.ValidateStateMachineDefinitionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$validateStateMachineDefinition, request,
        options: options);
  }

  // method descriptors

  static final _$createActivity =
      $grpc.ClientMethod<$0.CreateActivityInput, $0.CreateActivityOutput>(
          '/states.SFNService/CreateActivity',
          ($0.CreateActivityInput value) => value.writeToBuffer(),
          $0.CreateActivityOutput.fromBuffer);
  static final _$createStateMachine = $grpc.ClientMethod<
          $0.CreateStateMachineInput, $0.CreateStateMachineOutput>(
      '/states.SFNService/CreateStateMachine',
      ($0.CreateStateMachineInput value) => value.writeToBuffer(),
      $0.CreateStateMachineOutput.fromBuffer);
  static final _$createStateMachineAlias = $grpc.ClientMethod<
          $0.CreateStateMachineAliasInput, $0.CreateStateMachineAliasOutput>(
      '/states.SFNService/CreateStateMachineAlias',
      ($0.CreateStateMachineAliasInput value) => value.writeToBuffer(),
      $0.CreateStateMachineAliasOutput.fromBuffer);
  static final _$deleteActivity =
      $grpc.ClientMethod<$0.DeleteActivityInput, $0.DeleteActivityOutput>(
          '/states.SFNService/DeleteActivity',
          ($0.DeleteActivityInput value) => value.writeToBuffer(),
          $0.DeleteActivityOutput.fromBuffer);
  static final _$deleteStateMachine = $grpc.ClientMethod<
          $0.DeleteStateMachineInput, $0.DeleteStateMachineOutput>(
      '/states.SFNService/DeleteStateMachine',
      ($0.DeleteStateMachineInput value) => value.writeToBuffer(),
      $0.DeleteStateMachineOutput.fromBuffer);
  static final _$deleteStateMachineAlias = $grpc.ClientMethod<
          $0.DeleteStateMachineAliasInput, $0.DeleteStateMachineAliasOutput>(
      '/states.SFNService/DeleteStateMachineAlias',
      ($0.DeleteStateMachineAliasInput value) => value.writeToBuffer(),
      $0.DeleteStateMachineAliasOutput.fromBuffer);
  static final _$deleteStateMachineVersion = $grpc.ClientMethod<
          $0.DeleteStateMachineVersionInput,
          $0.DeleteStateMachineVersionOutput>(
      '/states.SFNService/DeleteStateMachineVersion',
      ($0.DeleteStateMachineVersionInput value) => value.writeToBuffer(),
      $0.DeleteStateMachineVersionOutput.fromBuffer);
  static final _$describeActivity =
      $grpc.ClientMethod<$0.DescribeActivityInput, $0.DescribeActivityOutput>(
          '/states.SFNService/DescribeActivity',
          ($0.DescribeActivityInput value) => value.writeToBuffer(),
          $0.DescribeActivityOutput.fromBuffer);
  static final _$describeExecution =
      $grpc.ClientMethod<$0.DescribeExecutionInput, $0.DescribeExecutionOutput>(
          '/states.SFNService/DescribeExecution',
          ($0.DescribeExecutionInput value) => value.writeToBuffer(),
          $0.DescribeExecutionOutput.fromBuffer);
  static final _$describeMapRun =
      $grpc.ClientMethod<$0.DescribeMapRunInput, $0.DescribeMapRunOutput>(
          '/states.SFNService/DescribeMapRun',
          ($0.DescribeMapRunInput value) => value.writeToBuffer(),
          $0.DescribeMapRunOutput.fromBuffer);
  static final _$describeStateMachine = $grpc.ClientMethod<
          $0.DescribeStateMachineInput, $0.DescribeStateMachineOutput>(
      '/states.SFNService/DescribeStateMachine',
      ($0.DescribeStateMachineInput value) => value.writeToBuffer(),
      $0.DescribeStateMachineOutput.fromBuffer);
  static final _$describeStateMachineAlias = $grpc.ClientMethod<
          $0.DescribeStateMachineAliasInput,
          $0.DescribeStateMachineAliasOutput>(
      '/states.SFNService/DescribeStateMachineAlias',
      ($0.DescribeStateMachineAliasInput value) => value.writeToBuffer(),
      $0.DescribeStateMachineAliasOutput.fromBuffer);
  static final _$describeStateMachineForExecution = $grpc.ClientMethod<
          $0.DescribeStateMachineForExecutionInput,
          $0.DescribeStateMachineForExecutionOutput>(
      '/states.SFNService/DescribeStateMachineForExecution',
      ($0.DescribeStateMachineForExecutionInput value) => value.writeToBuffer(),
      $0.DescribeStateMachineForExecutionOutput.fromBuffer);
  static final _$getActivityTask =
      $grpc.ClientMethod<$0.GetActivityTaskInput, $0.GetActivityTaskOutput>(
          '/states.SFNService/GetActivityTask',
          ($0.GetActivityTaskInput value) => value.writeToBuffer(),
          $0.GetActivityTaskOutput.fromBuffer);
  static final _$getExecutionHistory = $grpc.ClientMethod<
          $0.GetExecutionHistoryInput, $0.GetExecutionHistoryOutput>(
      '/states.SFNService/GetExecutionHistory',
      ($0.GetExecutionHistoryInput value) => value.writeToBuffer(),
      $0.GetExecutionHistoryOutput.fromBuffer);
  static final _$listActivities =
      $grpc.ClientMethod<$0.ListActivitiesInput, $0.ListActivitiesOutput>(
          '/states.SFNService/ListActivities',
          ($0.ListActivitiesInput value) => value.writeToBuffer(),
          $0.ListActivitiesOutput.fromBuffer);
  static final _$listExecutions =
      $grpc.ClientMethod<$0.ListExecutionsInput, $0.ListExecutionsOutput>(
          '/states.SFNService/ListExecutions',
          ($0.ListExecutionsInput value) => value.writeToBuffer(),
          $0.ListExecutionsOutput.fromBuffer);
  static final _$listMapRuns =
      $grpc.ClientMethod<$0.ListMapRunsInput, $0.ListMapRunsOutput>(
          '/states.SFNService/ListMapRuns',
          ($0.ListMapRunsInput value) => value.writeToBuffer(),
          $0.ListMapRunsOutput.fromBuffer);
  static final _$listStateMachineAliases = $grpc.ClientMethod<
          $0.ListStateMachineAliasesInput, $0.ListStateMachineAliasesOutput>(
      '/states.SFNService/ListStateMachineAliases',
      ($0.ListStateMachineAliasesInput value) => value.writeToBuffer(),
      $0.ListStateMachineAliasesOutput.fromBuffer);
  static final _$listStateMachines =
      $grpc.ClientMethod<$0.ListStateMachinesInput, $0.ListStateMachinesOutput>(
          '/states.SFNService/ListStateMachines',
          ($0.ListStateMachinesInput value) => value.writeToBuffer(),
          $0.ListStateMachinesOutput.fromBuffer);
  static final _$listStateMachineVersions = $grpc.ClientMethod<
          $0.ListStateMachineVersionsInput, $0.ListStateMachineVersionsOutput>(
      '/states.SFNService/ListStateMachineVersions',
      ($0.ListStateMachineVersionsInput value) => value.writeToBuffer(),
      $0.ListStateMachineVersionsOutput.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceInput, $0.ListTagsForResourceOutput>(
      '/states.SFNService/ListTagsForResource',
      ($0.ListTagsForResourceInput value) => value.writeToBuffer(),
      $0.ListTagsForResourceOutput.fromBuffer);
  static final _$publishStateMachineVersion = $grpc.ClientMethod<
          $0.PublishStateMachineVersionInput,
          $0.PublishStateMachineVersionOutput>(
      '/states.SFNService/PublishStateMachineVersion',
      ($0.PublishStateMachineVersionInput value) => value.writeToBuffer(),
      $0.PublishStateMachineVersionOutput.fromBuffer);
  static final _$redriveExecution =
      $grpc.ClientMethod<$0.RedriveExecutionInput, $0.RedriveExecutionOutput>(
          '/states.SFNService/RedriveExecution',
          ($0.RedriveExecutionInput value) => value.writeToBuffer(),
          $0.RedriveExecutionOutput.fromBuffer);
  static final _$sendTaskFailure =
      $grpc.ClientMethod<$0.SendTaskFailureInput, $0.SendTaskFailureOutput>(
          '/states.SFNService/SendTaskFailure',
          ($0.SendTaskFailureInput value) => value.writeToBuffer(),
          $0.SendTaskFailureOutput.fromBuffer);
  static final _$sendTaskHeartbeat =
      $grpc.ClientMethod<$0.SendTaskHeartbeatInput, $0.SendTaskHeartbeatOutput>(
          '/states.SFNService/SendTaskHeartbeat',
          ($0.SendTaskHeartbeatInput value) => value.writeToBuffer(),
          $0.SendTaskHeartbeatOutput.fromBuffer);
  static final _$sendTaskSuccess =
      $grpc.ClientMethod<$0.SendTaskSuccessInput, $0.SendTaskSuccessOutput>(
          '/states.SFNService/SendTaskSuccess',
          ($0.SendTaskSuccessInput value) => value.writeToBuffer(),
          $0.SendTaskSuccessOutput.fromBuffer);
  static final _$startExecution =
      $grpc.ClientMethod<$0.StartExecutionInput, $0.StartExecutionOutput>(
          '/states.SFNService/StartExecution',
          ($0.StartExecutionInput value) => value.writeToBuffer(),
          $0.StartExecutionOutput.fromBuffer);
  static final _$startSyncExecution = $grpc.ClientMethod<
          $0.StartSyncExecutionInput, $0.StartSyncExecutionOutput>(
      '/states.SFNService/StartSyncExecution',
      ($0.StartSyncExecutionInput value) => value.writeToBuffer(),
      $0.StartSyncExecutionOutput.fromBuffer);
  static final _$stopExecution =
      $grpc.ClientMethod<$0.StopExecutionInput, $0.StopExecutionOutput>(
          '/states.SFNService/StopExecution',
          ($0.StopExecutionInput value) => value.writeToBuffer(),
          $0.StopExecutionOutput.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceInput, $0.TagResourceOutput>(
          '/states.SFNService/TagResource',
          ($0.TagResourceInput value) => value.writeToBuffer(),
          $0.TagResourceOutput.fromBuffer);
  static final _$testState =
      $grpc.ClientMethod<$0.TestStateInput, $0.TestStateOutput>(
          '/states.SFNService/TestState',
          ($0.TestStateInput value) => value.writeToBuffer(),
          $0.TestStateOutput.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
          '/states.SFNService/UntagResource',
          ($0.UntagResourceInput value) => value.writeToBuffer(),
          $0.UntagResourceOutput.fromBuffer);
  static final _$updateMapRun =
      $grpc.ClientMethod<$0.UpdateMapRunInput, $0.UpdateMapRunOutput>(
          '/states.SFNService/UpdateMapRun',
          ($0.UpdateMapRunInput value) => value.writeToBuffer(),
          $0.UpdateMapRunOutput.fromBuffer);
  static final _$updateStateMachine = $grpc.ClientMethod<
          $0.UpdateStateMachineInput, $0.UpdateStateMachineOutput>(
      '/states.SFNService/UpdateStateMachine',
      ($0.UpdateStateMachineInput value) => value.writeToBuffer(),
      $0.UpdateStateMachineOutput.fromBuffer);
  static final _$updateStateMachineAlias = $grpc.ClientMethod<
          $0.UpdateStateMachineAliasInput, $0.UpdateStateMachineAliasOutput>(
      '/states.SFNService/UpdateStateMachineAlias',
      ($0.UpdateStateMachineAliasInput value) => value.writeToBuffer(),
      $0.UpdateStateMachineAliasOutput.fromBuffer);
  static final _$validateStateMachineDefinition = $grpc.ClientMethod<
          $0.ValidateStateMachineDefinitionInput,
          $0.ValidateStateMachineDefinitionOutput>(
      '/states.SFNService/ValidateStateMachineDefinition',
      ($0.ValidateStateMachineDefinitionInput value) => value.writeToBuffer(),
      $0.ValidateStateMachineDefinitionOutput.fromBuffer);
}

@$pb.GrpcServiceName('states.SFNService')
abstract class SFNServiceBase extends $grpc.Service {
  $core.String get $name => 'states.SFNService';

  SFNServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.CreateActivityInput, $0.CreateActivityOutput>(
            'CreateActivity',
            createActivity_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateActivityInput.fromBuffer(value),
            ($0.CreateActivityOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateStateMachineInput,
            $0.CreateStateMachineOutput>(
        'CreateStateMachine',
        createStateMachine_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateStateMachineInput.fromBuffer(value),
        ($0.CreateStateMachineOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateStateMachineAliasInput,
            $0.CreateStateMachineAliasOutput>(
        'CreateStateMachineAlias',
        createStateMachineAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateStateMachineAliasInput.fromBuffer(value),
        ($0.CreateStateMachineAliasOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteActivityInput, $0.DeleteActivityOutput>(
            'DeleteActivity',
            deleteActivity_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteActivityInput.fromBuffer(value),
            ($0.DeleteActivityOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteStateMachineInput,
            $0.DeleteStateMachineOutput>(
        'DeleteStateMachine',
        deleteStateMachine_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteStateMachineInput.fromBuffer(value),
        ($0.DeleteStateMachineOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteStateMachineAliasInput,
            $0.DeleteStateMachineAliasOutput>(
        'DeleteStateMachineAlias',
        deleteStateMachineAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteStateMachineAliasInput.fromBuffer(value),
        ($0.DeleteStateMachineAliasOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteStateMachineVersionInput,
            $0.DeleteStateMachineVersionOutput>(
        'DeleteStateMachineVersion',
        deleteStateMachineVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteStateMachineVersionInput.fromBuffer(value),
        ($0.DeleteStateMachineVersionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeActivityInput,
            $0.DescribeActivityOutput>(
        'DescribeActivity',
        describeActivity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeActivityInput.fromBuffer(value),
        ($0.DescribeActivityOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeExecutionInput,
            $0.DescribeExecutionOutput>(
        'DescribeExecution',
        describeExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeExecutionInput.fromBuffer(value),
        ($0.DescribeExecutionOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeMapRunInput, $0.DescribeMapRunOutput>(
            'DescribeMapRun',
            describeMapRun_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeMapRunInput.fromBuffer(value),
            ($0.DescribeMapRunOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeStateMachineInput,
            $0.DescribeStateMachineOutput>(
        'DescribeStateMachine',
        describeStateMachine_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeStateMachineInput.fromBuffer(value),
        ($0.DescribeStateMachineOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeStateMachineAliasInput,
            $0.DescribeStateMachineAliasOutput>(
        'DescribeStateMachineAlias',
        describeStateMachineAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeStateMachineAliasInput.fromBuffer(value),
        ($0.DescribeStateMachineAliasOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeStateMachineForExecutionInput,
            $0.DescribeStateMachineForExecutionOutput>(
        'DescribeStateMachineForExecution',
        describeStateMachineForExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeStateMachineForExecutionInput.fromBuffer(value),
        ($0.DescribeStateMachineForExecutionOutput value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetActivityTaskInput, $0.GetActivityTaskOutput>(
            'GetActivityTask',
            getActivityTask_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetActivityTaskInput.fromBuffer(value),
            ($0.GetActivityTaskOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetExecutionHistoryInput,
            $0.GetExecutionHistoryOutput>(
        'GetExecutionHistory',
        getExecutionHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetExecutionHistoryInput.fromBuffer(value),
        ($0.GetExecutionHistoryOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListActivitiesInput, $0.ListActivitiesOutput>(
            'ListActivities',
            listActivities_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListActivitiesInput.fromBuffer(value),
            ($0.ListActivitiesOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListExecutionsInput, $0.ListExecutionsOutput>(
            'ListExecutions',
            listExecutions_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListExecutionsInput.fromBuffer(value),
            ($0.ListExecutionsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMapRunsInput, $0.ListMapRunsOutput>(
        'ListMapRuns',
        listMapRuns_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListMapRunsInput.fromBuffer(value),
        ($0.ListMapRunsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListStateMachineAliasesInput,
            $0.ListStateMachineAliasesOutput>(
        'ListStateMachineAliases',
        listStateMachineAliases_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListStateMachineAliasesInput.fromBuffer(value),
        ($0.ListStateMachineAliasesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListStateMachinesInput,
            $0.ListStateMachinesOutput>(
        'ListStateMachines',
        listStateMachines_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListStateMachinesInput.fromBuffer(value),
        ($0.ListStateMachinesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListStateMachineVersionsInput,
            $0.ListStateMachineVersionsOutput>(
        'ListStateMachineVersions',
        listStateMachineVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListStateMachineVersionsInput.fromBuffer(value),
        ($0.ListStateMachineVersionsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceInput,
            $0.ListTagsForResourceOutput>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceInput.fromBuffer(value),
        ($0.ListTagsForResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PublishStateMachineVersionInput,
            $0.PublishStateMachineVersionOutput>(
        'PublishStateMachineVersion',
        publishStateMachineVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PublishStateMachineVersionInput.fromBuffer(value),
        ($0.PublishStateMachineVersionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RedriveExecutionInput,
            $0.RedriveExecutionOutput>(
        'RedriveExecution',
        redriveExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RedriveExecutionInput.fromBuffer(value),
        ($0.RedriveExecutionOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SendTaskFailureInput, $0.SendTaskFailureOutput>(
            'SendTaskFailure',
            sendTaskFailure_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SendTaskFailureInput.fromBuffer(value),
            ($0.SendTaskFailureOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendTaskHeartbeatInput,
            $0.SendTaskHeartbeatOutput>(
        'SendTaskHeartbeat',
        sendTaskHeartbeat_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendTaskHeartbeatInput.fromBuffer(value),
        ($0.SendTaskHeartbeatOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SendTaskSuccessInput, $0.SendTaskSuccessOutput>(
            'SendTaskSuccess',
            sendTaskSuccess_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SendTaskSuccessInput.fromBuffer(value),
            ($0.SendTaskSuccessOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartExecutionInput, $0.StartExecutionOutput>(
            'StartExecution',
            startExecution_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartExecutionInput.fromBuffer(value),
            ($0.StartExecutionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartSyncExecutionInput,
            $0.StartSyncExecutionOutput>(
        'StartSyncExecution',
        startSyncExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartSyncExecutionInput.fromBuffer(value),
        ($0.StartSyncExecutionOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StopExecutionInput, $0.StopExecutionOutput>(
            'StopExecution',
            stopExecution_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StopExecutionInput.fromBuffer(value),
            ($0.StopExecutionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceInput, $0.TagResourceOutput>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagResourceInput.fromBuffer(value),
        ($0.TagResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestStateInput, $0.TestStateOutput>(
        'TestState',
        testState_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TestStateInput.fromBuffer(value),
        ($0.TestStateOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceInput.fromBuffer(value),
            ($0.UntagResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateMapRunInput, $0.UpdateMapRunOutput>(
        'UpdateMapRun',
        updateMapRun_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateMapRunInput.fromBuffer(value),
        ($0.UpdateMapRunOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateStateMachineInput,
            $0.UpdateStateMachineOutput>(
        'UpdateStateMachine',
        updateStateMachine_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateStateMachineInput.fromBuffer(value),
        ($0.UpdateStateMachineOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateStateMachineAliasInput,
            $0.UpdateStateMachineAliasOutput>(
        'UpdateStateMachineAlias',
        updateStateMachineAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateStateMachineAliasInput.fromBuffer(value),
        ($0.UpdateStateMachineAliasOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ValidateStateMachineDefinitionInput,
            $0.ValidateStateMachineDefinitionOutput>(
        'ValidateStateMachineDefinition',
        validateStateMachineDefinition_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ValidateStateMachineDefinitionInput.fromBuffer(value),
        ($0.ValidateStateMachineDefinitionOutput value) =>
            value.writeToBuffer()));
  }

  $async.Future<$0.CreateActivityOutput> createActivity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateActivityInput> $request) async {
    return createActivity($call, await $request);
  }

  $async.Future<$0.CreateActivityOutput> createActivity(
      $grpc.ServiceCall call, $0.CreateActivityInput request);

  $async.Future<$0.CreateStateMachineOutput> createStateMachine_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateStateMachineInput> $request) async {
    return createStateMachine($call, await $request);
  }

  $async.Future<$0.CreateStateMachineOutput> createStateMachine(
      $grpc.ServiceCall call, $0.CreateStateMachineInput request);

  $async.Future<$0.CreateStateMachineAliasOutput> createStateMachineAlias_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateStateMachineAliasInput> $request) async {
    return createStateMachineAlias($call, await $request);
  }

  $async.Future<$0.CreateStateMachineAliasOutput> createStateMachineAlias(
      $grpc.ServiceCall call, $0.CreateStateMachineAliasInput request);

  $async.Future<$0.DeleteActivityOutput> deleteActivity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteActivityInput> $request) async {
    return deleteActivity($call, await $request);
  }

  $async.Future<$0.DeleteActivityOutput> deleteActivity(
      $grpc.ServiceCall call, $0.DeleteActivityInput request);

  $async.Future<$0.DeleteStateMachineOutput> deleteStateMachine_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteStateMachineInput> $request) async {
    return deleteStateMachine($call, await $request);
  }

  $async.Future<$0.DeleteStateMachineOutput> deleteStateMachine(
      $grpc.ServiceCall call, $0.DeleteStateMachineInput request);

  $async.Future<$0.DeleteStateMachineAliasOutput> deleteStateMachineAlias_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteStateMachineAliasInput> $request) async {
    return deleteStateMachineAlias($call, await $request);
  }

  $async.Future<$0.DeleteStateMachineAliasOutput> deleteStateMachineAlias(
      $grpc.ServiceCall call, $0.DeleteStateMachineAliasInput request);

  $async.Future<$0.DeleteStateMachineVersionOutput>
      deleteStateMachineVersion_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteStateMachineVersionInput> $request) async {
    return deleteStateMachineVersion($call, await $request);
  }

  $async.Future<$0.DeleteStateMachineVersionOutput> deleteStateMachineVersion(
      $grpc.ServiceCall call, $0.DeleteStateMachineVersionInput request);

  $async.Future<$0.DescribeActivityOutput> describeActivity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeActivityInput> $request) async {
    return describeActivity($call, await $request);
  }

  $async.Future<$0.DescribeActivityOutput> describeActivity(
      $grpc.ServiceCall call, $0.DescribeActivityInput request);

  $async.Future<$0.DescribeExecutionOutput> describeExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeExecutionInput> $request) async {
    return describeExecution($call, await $request);
  }

  $async.Future<$0.DescribeExecutionOutput> describeExecution(
      $grpc.ServiceCall call, $0.DescribeExecutionInput request);

  $async.Future<$0.DescribeMapRunOutput> describeMapRun_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeMapRunInput> $request) async {
    return describeMapRun($call, await $request);
  }

  $async.Future<$0.DescribeMapRunOutput> describeMapRun(
      $grpc.ServiceCall call, $0.DescribeMapRunInput request);

  $async.Future<$0.DescribeStateMachineOutput> describeStateMachine_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeStateMachineInput> $request) async {
    return describeStateMachine($call, await $request);
  }

  $async.Future<$0.DescribeStateMachineOutput> describeStateMachine(
      $grpc.ServiceCall call, $0.DescribeStateMachineInput request);

  $async.Future<$0.DescribeStateMachineAliasOutput>
      describeStateMachineAlias_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeStateMachineAliasInput> $request) async {
    return describeStateMachineAlias($call, await $request);
  }

  $async.Future<$0.DescribeStateMachineAliasOutput> describeStateMachineAlias(
      $grpc.ServiceCall call, $0.DescribeStateMachineAliasInput request);

  $async.Future<$0.DescribeStateMachineForExecutionOutput>
      describeStateMachineForExecution_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeStateMachineForExecutionInput>
              $request) async {
    return describeStateMachineForExecution($call, await $request);
  }

  $async.Future<$0.DescribeStateMachineForExecutionOutput>
      describeStateMachineForExecution($grpc.ServiceCall call,
          $0.DescribeStateMachineForExecutionInput request);

  $async.Future<$0.GetActivityTaskOutput> getActivityTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetActivityTaskInput> $request) async {
    return getActivityTask($call, await $request);
  }

  $async.Future<$0.GetActivityTaskOutput> getActivityTask(
      $grpc.ServiceCall call, $0.GetActivityTaskInput request);

  $async.Future<$0.GetExecutionHistoryOutput> getExecutionHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetExecutionHistoryInput> $request) async {
    return getExecutionHistory($call, await $request);
  }

  $async.Future<$0.GetExecutionHistoryOutput> getExecutionHistory(
      $grpc.ServiceCall call, $0.GetExecutionHistoryInput request);

  $async.Future<$0.ListActivitiesOutput> listActivities_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListActivitiesInput> $request) async {
    return listActivities($call, await $request);
  }

  $async.Future<$0.ListActivitiesOutput> listActivities(
      $grpc.ServiceCall call, $0.ListActivitiesInput request);

  $async.Future<$0.ListExecutionsOutput> listExecutions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListExecutionsInput> $request) async {
    return listExecutions($call, await $request);
  }

  $async.Future<$0.ListExecutionsOutput> listExecutions(
      $grpc.ServiceCall call, $0.ListExecutionsInput request);

  $async.Future<$0.ListMapRunsOutput> listMapRuns_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListMapRunsInput> $request) async {
    return listMapRuns($call, await $request);
  }

  $async.Future<$0.ListMapRunsOutput> listMapRuns(
      $grpc.ServiceCall call, $0.ListMapRunsInput request);

  $async.Future<$0.ListStateMachineAliasesOutput> listStateMachineAliases_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListStateMachineAliasesInput> $request) async {
    return listStateMachineAliases($call, await $request);
  }

  $async.Future<$0.ListStateMachineAliasesOutput> listStateMachineAliases(
      $grpc.ServiceCall call, $0.ListStateMachineAliasesInput request);

  $async.Future<$0.ListStateMachinesOutput> listStateMachines_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListStateMachinesInput> $request) async {
    return listStateMachines($call, await $request);
  }

  $async.Future<$0.ListStateMachinesOutput> listStateMachines(
      $grpc.ServiceCall call, $0.ListStateMachinesInput request);

  $async.Future<$0.ListStateMachineVersionsOutput> listStateMachineVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListStateMachineVersionsInput> $request) async {
    return listStateMachineVersions($call, await $request);
  }

  $async.Future<$0.ListStateMachineVersionsOutput> listStateMachineVersions(
      $grpc.ServiceCall call, $0.ListStateMachineVersionsInput request);

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceInput> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceInput request);

  $async.Future<$0.PublishStateMachineVersionOutput>
      publishStateMachineVersion_Pre($grpc.ServiceCall $call,
          $async.Future<$0.PublishStateMachineVersionInput> $request) async {
    return publishStateMachineVersion($call, await $request);
  }

  $async.Future<$0.PublishStateMachineVersionOutput> publishStateMachineVersion(
      $grpc.ServiceCall call, $0.PublishStateMachineVersionInput request);

  $async.Future<$0.RedriveExecutionOutput> redriveExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RedriveExecutionInput> $request) async {
    return redriveExecution($call, await $request);
  }

  $async.Future<$0.RedriveExecutionOutput> redriveExecution(
      $grpc.ServiceCall call, $0.RedriveExecutionInput request);

  $async.Future<$0.SendTaskFailureOutput> sendTaskFailure_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendTaskFailureInput> $request) async {
    return sendTaskFailure($call, await $request);
  }

  $async.Future<$0.SendTaskFailureOutput> sendTaskFailure(
      $grpc.ServiceCall call, $0.SendTaskFailureInput request);

  $async.Future<$0.SendTaskHeartbeatOutput> sendTaskHeartbeat_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendTaskHeartbeatInput> $request) async {
    return sendTaskHeartbeat($call, await $request);
  }

  $async.Future<$0.SendTaskHeartbeatOutput> sendTaskHeartbeat(
      $grpc.ServiceCall call, $0.SendTaskHeartbeatInput request);

  $async.Future<$0.SendTaskSuccessOutput> sendTaskSuccess_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendTaskSuccessInput> $request) async {
    return sendTaskSuccess($call, await $request);
  }

  $async.Future<$0.SendTaskSuccessOutput> sendTaskSuccess(
      $grpc.ServiceCall call, $0.SendTaskSuccessInput request);

  $async.Future<$0.StartExecutionOutput> startExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartExecutionInput> $request) async {
    return startExecution($call, await $request);
  }

  $async.Future<$0.StartExecutionOutput> startExecution(
      $grpc.ServiceCall call, $0.StartExecutionInput request);

  $async.Future<$0.StartSyncExecutionOutput> startSyncExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartSyncExecutionInput> $request) async {
    return startSyncExecution($call, await $request);
  }

  $async.Future<$0.StartSyncExecutionOutput> startSyncExecution(
      $grpc.ServiceCall call, $0.StartSyncExecutionInput request);

  $async.Future<$0.StopExecutionOutput> stopExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StopExecutionInput> $request) async {
    return stopExecution($call, await $request);
  }

  $async.Future<$0.StopExecutionOutput> stopExecution(
      $grpc.ServiceCall call, $0.StopExecutionInput request);

  $async.Future<$0.TagResourceOutput> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceInput> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceOutput> tagResource(
      $grpc.ServiceCall call, $0.TagResourceInput request);

  $async.Future<$0.TestStateOutput> testState_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TestStateInput> $request) async {
    return testState($call, await $request);
  }

  $async.Future<$0.TestStateOutput> testState(
      $grpc.ServiceCall call, $0.TestStateInput request);

  $async.Future<$0.UntagResourceOutput> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceInput> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceOutput> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceInput request);

  $async.Future<$0.UpdateMapRunOutput> updateMapRun_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateMapRunInput> $request) async {
    return updateMapRun($call, await $request);
  }

  $async.Future<$0.UpdateMapRunOutput> updateMapRun(
      $grpc.ServiceCall call, $0.UpdateMapRunInput request);

  $async.Future<$0.UpdateStateMachineOutput> updateStateMachine_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateStateMachineInput> $request) async {
    return updateStateMachine($call, await $request);
  }

  $async.Future<$0.UpdateStateMachineOutput> updateStateMachine(
      $grpc.ServiceCall call, $0.UpdateStateMachineInput request);

  $async.Future<$0.UpdateStateMachineAliasOutput> updateStateMachineAlias_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateStateMachineAliasInput> $request) async {
    return updateStateMachineAlias($call, await $request);
  }

  $async.Future<$0.UpdateStateMachineAliasOutput> updateStateMachineAlias(
      $grpc.ServiceCall call, $0.UpdateStateMachineAliasInput request);

  $async.Future<$0.ValidateStateMachineDefinitionOutput>
      validateStateMachineDefinition_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ValidateStateMachineDefinitionInput>
              $request) async {
    return validateStateMachineDefinition($call, await $request);
  }

  $async.Future<$0.ValidateStateMachineDefinitionOutput>
      validateStateMachineDefinition($grpc.ServiceCall call,
          $0.ValidateStateMachineDefinitionInput request);
}
