// This is a generated file - do not edit.
//
// Generated from scheduler.proto.

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

import 'scheduler.pb.dart' as $0;

export 'scheduler.pb.dart';

/// SchedulerService provides scheduler API operations.
@$pb.GrpcServiceName('scheduler.SchedulerService')
class SchedulerServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SchedulerServiceClient(super.channel, {super.options, super.interceptors});

  /// Lists the tags associated with the Scheduler resource.
  /// HTTP: GET /tags/{ResourceArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListTagsForResourceOutput> listTagsForResource(
    $0.ListTagsForResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Assigns one or more tags (key-value pairs) to the specified EventBridge Scheduler resource. You can only assign tags to schedule groups.
  /// HTTP: POST /tags/{ResourceArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.TagResourceOutput> tagResource(
    $0.TagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Removes one or more tags from the specified EventBridge Scheduler schedule group.
  /// HTTP: DELETE /tags/{ResourceArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UntagResourceOutput> untagResource(
    $0.UntagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Returns a paginated list of your EventBridge Scheduler schedules.
  /// HTTP: GET /schedules
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListSchedulesOutput> listSchedules(
    $0.ListSchedulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSchedules, request, options: options);
  }

  /// Creates the specified schedule.
  /// HTTP: POST /schedules/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateScheduleOutput> createSchedule(
    $0.CreateScheduleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSchedule, request, options: options);
  }

  /// Retrieves the specified schedule.
  /// HTTP: GET /schedules/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetScheduleOutput> getSchedule(
    $0.GetScheduleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSchedule, request, options: options);
  }

  /// Updates the specified schedule. When you call UpdateSchedule, EventBridge Scheduler uses all values, including empty values, specified in the request and overrides the existing schedule. This is by...
  /// HTTP: PUT /schedules/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateScheduleOutput> updateSchedule(
    $0.UpdateScheduleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSchedule, request, options: options);
  }

  /// Deletes the specified schedule.
  /// HTTP: DELETE /schedules/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteScheduleOutput> deleteSchedule(
    $0.DeleteScheduleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSchedule, request, options: options);
  }

  /// Returns a paginated list of your schedule groups.
  /// HTTP: GET /schedule-groups
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListScheduleGroupsOutput> listScheduleGroups(
    $0.ListScheduleGroupsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listScheduleGroups, request, options: options);
  }

  /// Creates the specified schedule group.
  /// HTTP: POST /schedule-groups/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateScheduleGroupOutput> createScheduleGroup(
    $0.CreateScheduleGroupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createScheduleGroup, request, options: options);
  }

  /// Retrieves the specified schedule group.
  /// HTTP: GET /schedule-groups/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetScheduleGroupOutput> getScheduleGroup(
    $0.GetScheduleGroupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getScheduleGroup, request, options: options);
  }

  /// Deletes the specified schedule group. Deleting a schedule group results in EventBridge Scheduler deleting all schedules associated with the group. When you delete a group, it remains in a DELETING ...
  /// HTTP: DELETE /schedule-groups/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteScheduleGroupOutput> deleteScheduleGroup(
    $0.DeleteScheduleGroupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteScheduleGroup, request, options: options);
  }

  // method descriptors

  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceInput, $0.ListTagsForResourceOutput>(
      '/scheduler.SchedulerService/ListTagsForResource',
      ($0.ListTagsForResourceInput value) => value.writeToBuffer(),
      $0.ListTagsForResourceOutput.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceInput, $0.TagResourceOutput>(
          '/scheduler.SchedulerService/TagResource',
          ($0.TagResourceInput value) => value.writeToBuffer(),
          $0.TagResourceOutput.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
          '/scheduler.SchedulerService/UntagResource',
          ($0.UntagResourceInput value) => value.writeToBuffer(),
          $0.UntagResourceOutput.fromBuffer);
  static final _$listSchedules =
      $grpc.ClientMethod<$0.ListSchedulesInput, $0.ListSchedulesOutput>(
          '/scheduler.SchedulerService/ListSchedules',
          ($0.ListSchedulesInput value) => value.writeToBuffer(),
          $0.ListSchedulesOutput.fromBuffer);
  static final _$createSchedule =
      $grpc.ClientMethod<$0.CreateScheduleInput, $0.CreateScheduleOutput>(
          '/scheduler.SchedulerService/CreateSchedule',
          ($0.CreateScheduleInput value) => value.writeToBuffer(),
          $0.CreateScheduleOutput.fromBuffer);
  static final _$getSchedule =
      $grpc.ClientMethod<$0.GetScheduleInput, $0.GetScheduleOutput>(
          '/scheduler.SchedulerService/GetSchedule',
          ($0.GetScheduleInput value) => value.writeToBuffer(),
          $0.GetScheduleOutput.fromBuffer);
  static final _$updateSchedule =
      $grpc.ClientMethod<$0.UpdateScheduleInput, $0.UpdateScheduleOutput>(
          '/scheduler.SchedulerService/UpdateSchedule',
          ($0.UpdateScheduleInput value) => value.writeToBuffer(),
          $0.UpdateScheduleOutput.fromBuffer);
  static final _$deleteSchedule =
      $grpc.ClientMethod<$0.DeleteScheduleInput, $0.DeleteScheduleOutput>(
          '/scheduler.SchedulerService/DeleteSchedule',
          ($0.DeleteScheduleInput value) => value.writeToBuffer(),
          $0.DeleteScheduleOutput.fromBuffer);
  static final _$listScheduleGroups = $grpc.ClientMethod<
          $0.ListScheduleGroupsInput, $0.ListScheduleGroupsOutput>(
      '/scheduler.SchedulerService/ListScheduleGroups',
      ($0.ListScheduleGroupsInput value) => value.writeToBuffer(),
      $0.ListScheduleGroupsOutput.fromBuffer);
  static final _$createScheduleGroup = $grpc.ClientMethod<
          $0.CreateScheduleGroupInput, $0.CreateScheduleGroupOutput>(
      '/scheduler.SchedulerService/CreateScheduleGroup',
      ($0.CreateScheduleGroupInput value) => value.writeToBuffer(),
      $0.CreateScheduleGroupOutput.fromBuffer);
  static final _$getScheduleGroup =
      $grpc.ClientMethod<$0.GetScheduleGroupInput, $0.GetScheduleGroupOutput>(
          '/scheduler.SchedulerService/GetScheduleGroup',
          ($0.GetScheduleGroupInput value) => value.writeToBuffer(),
          $0.GetScheduleGroupOutput.fromBuffer);
  static final _$deleteScheduleGroup = $grpc.ClientMethod<
          $0.DeleteScheduleGroupInput, $0.DeleteScheduleGroupOutput>(
      '/scheduler.SchedulerService/DeleteScheduleGroup',
      ($0.DeleteScheduleGroupInput value) => value.writeToBuffer(),
      $0.DeleteScheduleGroupOutput.fromBuffer);
}

@$pb.GrpcServiceName('scheduler.SchedulerService')
abstract class SchedulerServiceBase extends $grpc.Service {
  $core.String get $name => 'scheduler.SchedulerService';

  SchedulerServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceInput,
            $0.ListTagsForResourceOutput>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceInput.fromBuffer(value),
        ($0.ListTagsForResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceInput, $0.TagResourceOutput>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagResourceInput.fromBuffer(value),
        ($0.TagResourceOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceInput.fromBuffer(value),
            ($0.UntagResourceOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListSchedulesInput, $0.ListSchedulesOutput>(
            'ListSchedules',
            listSchedules_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListSchedulesInput.fromBuffer(value),
            ($0.ListSchedulesOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateScheduleInput, $0.CreateScheduleOutput>(
            'CreateSchedule',
            createSchedule_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateScheduleInput.fromBuffer(value),
            ($0.CreateScheduleOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetScheduleInput, $0.GetScheduleOutput>(
        'GetSchedule',
        getSchedule_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetScheduleInput.fromBuffer(value),
        ($0.GetScheduleOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateScheduleInput, $0.UpdateScheduleOutput>(
            'UpdateSchedule',
            updateSchedule_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateScheduleInput.fromBuffer(value),
            ($0.UpdateScheduleOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteScheduleInput, $0.DeleteScheduleOutput>(
            'DeleteSchedule',
            deleteSchedule_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteScheduleInput.fromBuffer(value),
            ($0.DeleteScheduleOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListScheduleGroupsInput,
            $0.ListScheduleGroupsOutput>(
        'ListScheduleGroups',
        listScheduleGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListScheduleGroupsInput.fromBuffer(value),
        ($0.ListScheduleGroupsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateScheduleGroupInput,
            $0.CreateScheduleGroupOutput>(
        'CreateScheduleGroup',
        createScheduleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateScheduleGroupInput.fromBuffer(value),
        ($0.CreateScheduleGroupOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetScheduleGroupInput,
            $0.GetScheduleGroupOutput>(
        'GetScheduleGroup',
        getScheduleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetScheduleGroupInput.fromBuffer(value),
        ($0.GetScheduleGroupOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteScheduleGroupInput,
            $0.DeleteScheduleGroupOutput>(
        'DeleteScheduleGroup',
        deleteScheduleGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteScheduleGroupInput.fromBuffer(value),
        ($0.DeleteScheduleGroupOutput value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceInput> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceInput request);

  $async.Future<$0.TagResourceOutput> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceInput> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceOutput> tagResource(
      $grpc.ServiceCall call, $0.TagResourceInput request);

  $async.Future<$0.UntagResourceOutput> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceInput> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceOutput> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceInput request);

  $async.Future<$0.ListSchedulesOutput> listSchedules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSchedulesInput> $request) async {
    return listSchedules($call, await $request);
  }

  $async.Future<$0.ListSchedulesOutput> listSchedules(
      $grpc.ServiceCall call, $0.ListSchedulesInput request);

  $async.Future<$0.CreateScheduleOutput> createSchedule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateScheduleInput> $request) async {
    return createSchedule($call, await $request);
  }

  $async.Future<$0.CreateScheduleOutput> createSchedule(
      $grpc.ServiceCall call, $0.CreateScheduleInput request);

  $async.Future<$0.GetScheduleOutput> getSchedule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetScheduleInput> $request) async {
    return getSchedule($call, await $request);
  }

  $async.Future<$0.GetScheduleOutput> getSchedule(
      $grpc.ServiceCall call, $0.GetScheduleInput request);

  $async.Future<$0.UpdateScheduleOutput> updateSchedule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateScheduleInput> $request) async {
    return updateSchedule($call, await $request);
  }

  $async.Future<$0.UpdateScheduleOutput> updateSchedule(
      $grpc.ServiceCall call, $0.UpdateScheduleInput request);

  $async.Future<$0.DeleteScheduleOutput> deleteSchedule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteScheduleInput> $request) async {
    return deleteSchedule($call, await $request);
  }

  $async.Future<$0.DeleteScheduleOutput> deleteSchedule(
      $grpc.ServiceCall call, $0.DeleteScheduleInput request);

  $async.Future<$0.ListScheduleGroupsOutput> listScheduleGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListScheduleGroupsInput> $request) async {
    return listScheduleGroups($call, await $request);
  }

  $async.Future<$0.ListScheduleGroupsOutput> listScheduleGroups(
      $grpc.ServiceCall call, $0.ListScheduleGroupsInput request);

  $async.Future<$0.CreateScheduleGroupOutput> createScheduleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateScheduleGroupInput> $request) async {
    return createScheduleGroup($call, await $request);
  }

  $async.Future<$0.CreateScheduleGroupOutput> createScheduleGroup(
      $grpc.ServiceCall call, $0.CreateScheduleGroupInput request);

  $async.Future<$0.GetScheduleGroupOutput> getScheduleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetScheduleGroupInput> $request) async {
    return getScheduleGroup($call, await $request);
  }

  $async.Future<$0.GetScheduleGroupOutput> getScheduleGroup(
      $grpc.ServiceCall call, $0.GetScheduleGroupInput request);

  $async.Future<$0.DeleteScheduleGroupOutput> deleteScheduleGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteScheduleGroupInput> $request) async {
    return deleteScheduleGroup($call, await $request);
  }

  $async.Future<$0.DeleteScheduleGroupOutput> deleteScheduleGroup(
      $grpc.ServiceCall call, $0.DeleteScheduleGroupInput request);
}
