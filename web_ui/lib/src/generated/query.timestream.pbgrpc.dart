// This is a generated file - do not edit.
//
// Generated from query.timestream.proto.

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
import 'query.timestream.pb.dart' as $0;

export 'query.timestream.pb.dart';

/// TimestreamQueryService provides query.timestream API operations.
@$pb.GrpcServiceName('query.timestream.TimestreamQueryService')
class TimestreamQueryServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  TimestreamQueryServiceClient(super.channel,
      {super.options, super.interceptors});

  /// Cancels a query that has been issued. Cancellation is provided only if the query has not completed running before the cancellation request was issued. Because cancellation is an idempotent operatio...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CancelQueryResponse> cancelQuery(
    $0.CancelQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelQuery, request, options: options);
  }

  /// Create a scheduled query that will be run on your behalf at the configured schedule. Timestream assumes the execution role provided as part of the ScheduledQueryExecutionRoleArn parameter to run th...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateScheduledQueryResponse> createScheduledQuery(
    $0.CreateScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createScheduledQuery, request, options: options);
  }

  /// Deletes a given scheduled query. This is an irreversible operation.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> deleteScheduledQuery(
    $0.DeleteScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteScheduledQuery, request, options: options);
  }

  /// Describes the settings for your account that include the query pricing model and the configured maximum TCUs the service can use for your query workload. You're charged only for the duration of com...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeAccountSettingsResponse>
      describeAccountSettings(
    $0.DescribeAccountSettingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAccountSettings, request,
        options: options);
  }

  /// DescribeEndpoints returns a list of available endpoints to make Timestream API calls against. This API is available through both Write and Query. Because the Timestream SDKs are designed to transpa...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeEndpointsResponse> describeEndpoints(
    $0.DescribeEndpointsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeEndpoints, request, options: options);
  }

  /// Provides detailed information about a scheduled query.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeScheduledQueryResponse>
      describeScheduledQuery(
    $0.DescribeScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeScheduledQuery, request,
        options: options);
  }

  /// You can use this API to run a scheduled query manually. If you enabled QueryInsights, this API also returns insights and metrics related to the query that you executed as part of an Amazon SNS noti...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> executeScheduledQuery(
    $0.ExecuteScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$executeScheduledQuery, request, options: options);
  }

  /// Gets a list of all scheduled queries in the caller's Amazon account and Region. ListScheduledQueries is eventually consistent.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListScheduledQueriesResponse> listScheduledQueries(
    $0.ListScheduledQueriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listScheduledQueries, request, options: options);
  }

  /// List all tags on a Timestream query resource.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// A synchronous operation that allows you to submit a query with parameters to be stored by Timestream for later running. Timestream only supports using this operation with ValidateOnly set to true.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PrepareQueryResponse> prepareQuery(
    $0.PrepareQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$prepareQuery, request, options: options);
  }

  /// Query is a synchronous operation that enables you to run a query against your Amazon Timestream data. If you enabled QueryInsights, this API also returns insights and metrics related to the query t...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.QueryResponse> query(
    $0.QueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$query, request, options: options);
  }

  /// Associate a set of tags with a Timestream resource. You can then activate these user-defined tags so that they appear on the Billing and Cost Management console for cost allocation tracking.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Removes the association of tags from a Timestream query resource.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Transitions your account to use TCUs for query pricing and modifies the maximum query compute units that you've configured. If you reduce the value of MaxQueryTCU to a desired configuration, the ne...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateAccountSettingsResponse> updateAccountSettings(
    $0.UpdateAccountSettingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAccountSettings, request, options: options);
  }

  /// Update a scheduled query.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> updateScheduledQuery(
    $0.UpdateScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateScheduledQuery, request, options: options);
  }

  // method descriptors

  static final _$cancelQuery =
      $grpc.ClientMethod<$0.CancelQueryRequest, $0.CancelQueryResponse>(
          '/query.timestream.TimestreamQueryService/CancelQuery',
          ($0.CancelQueryRequest value) => value.writeToBuffer(),
          $0.CancelQueryResponse.fromBuffer);
  static final _$createScheduledQuery = $grpc.ClientMethod<
          $0.CreateScheduledQueryRequest, $0.CreateScheduledQueryResponse>(
      '/query.timestream.TimestreamQueryService/CreateScheduledQuery',
      ($0.CreateScheduledQueryRequest value) => value.writeToBuffer(),
      $0.CreateScheduledQueryResponse.fromBuffer);
  static final _$deleteScheduledQuery =
      $grpc.ClientMethod<$0.DeleteScheduledQueryRequest, $1.Empty>(
          '/query.timestream.TimestreamQueryService/DeleteScheduledQuery',
          ($0.DeleteScheduledQueryRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$describeAccountSettings = $grpc.ClientMethod<
          $0.DescribeAccountSettingsRequest,
          $0.DescribeAccountSettingsResponse>(
      '/query.timestream.TimestreamQueryService/DescribeAccountSettings',
      ($0.DescribeAccountSettingsRequest value) => value.writeToBuffer(),
      $0.DescribeAccountSettingsResponse.fromBuffer);
  static final _$describeEndpoints = $grpc.ClientMethod<
          $0.DescribeEndpointsRequest, $0.DescribeEndpointsResponse>(
      '/query.timestream.TimestreamQueryService/DescribeEndpoints',
      ($0.DescribeEndpointsRequest value) => value.writeToBuffer(),
      $0.DescribeEndpointsResponse.fromBuffer);
  static final _$describeScheduledQuery = $grpc.ClientMethod<
          $0.DescribeScheduledQueryRequest, $0.DescribeScheduledQueryResponse>(
      '/query.timestream.TimestreamQueryService/DescribeScheduledQuery',
      ($0.DescribeScheduledQueryRequest value) => value.writeToBuffer(),
      $0.DescribeScheduledQueryResponse.fromBuffer);
  static final _$executeScheduledQuery =
      $grpc.ClientMethod<$0.ExecuteScheduledQueryRequest, $1.Empty>(
          '/query.timestream.TimestreamQueryService/ExecuteScheduledQuery',
          ($0.ExecuteScheduledQueryRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$listScheduledQueries = $grpc.ClientMethod<
          $0.ListScheduledQueriesRequest, $0.ListScheduledQueriesResponse>(
      '/query.timestream.TimestreamQueryService/ListScheduledQueries',
      ($0.ListScheduledQueriesRequest value) => value.writeToBuffer(),
      $0.ListScheduledQueriesResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/query.timestream.TimestreamQueryService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$prepareQuery =
      $grpc.ClientMethod<$0.PrepareQueryRequest, $0.PrepareQueryResponse>(
          '/query.timestream.TimestreamQueryService/PrepareQuery',
          ($0.PrepareQueryRequest value) => value.writeToBuffer(),
          $0.PrepareQueryResponse.fromBuffer);
  static final _$query = $grpc.ClientMethod<$0.QueryRequest, $0.QueryResponse>(
      '/query.timestream.TimestreamQueryService/Query',
      ($0.QueryRequest value) => value.writeToBuffer(),
      $0.QueryResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/query.timestream.TimestreamQueryService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/query.timestream.TimestreamQueryService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateAccountSettings = $grpc.ClientMethod<
          $0.UpdateAccountSettingsRequest, $0.UpdateAccountSettingsResponse>(
      '/query.timestream.TimestreamQueryService/UpdateAccountSettings',
      ($0.UpdateAccountSettingsRequest value) => value.writeToBuffer(),
      $0.UpdateAccountSettingsResponse.fromBuffer);
  static final _$updateScheduledQuery =
      $grpc.ClientMethod<$0.UpdateScheduledQueryRequest, $1.Empty>(
          '/query.timestream.TimestreamQueryService/UpdateScheduledQuery',
          ($0.UpdateScheduledQueryRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
}

@$pb.GrpcServiceName('query.timestream.TimestreamQueryService')
abstract class TimestreamQueryServiceBase extends $grpc.Service {
  $core.String get $name => 'query.timestream.TimestreamQueryService';

  TimestreamQueryServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.CancelQueryRequest, $0.CancelQueryResponse>(
            'CancelQuery',
            cancelQuery_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CancelQueryRequest.fromBuffer(value),
            ($0.CancelQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateScheduledQueryRequest,
            $0.CreateScheduledQueryResponse>(
        'CreateScheduledQuery',
        createScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateScheduledQueryRequest.fromBuffer(value),
        ($0.CreateScheduledQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteScheduledQueryRequest, $1.Empty>(
        'DeleteScheduledQuery',
        deleteScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteScheduledQueryRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAccountSettingsRequest,
            $0.DescribeAccountSettingsResponse>(
        'DescribeAccountSettings',
        describeAccountSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAccountSettingsRequest.fromBuffer(value),
        ($0.DescribeAccountSettingsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeEndpointsRequest,
            $0.DescribeEndpointsResponse>(
        'DescribeEndpoints',
        describeEndpoints_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeEndpointsRequest.fromBuffer(value),
        ($0.DescribeEndpointsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeScheduledQueryRequest,
            $0.DescribeScheduledQueryResponse>(
        'DescribeScheduledQuery',
        describeScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeScheduledQueryRequest.fromBuffer(value),
        ($0.DescribeScheduledQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ExecuteScheduledQueryRequest, $1.Empty>(
        'ExecuteScheduledQuery',
        executeScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ExecuteScheduledQueryRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListScheduledQueriesRequest,
            $0.ListScheduledQueriesResponse>(
        'ListScheduledQueries',
        listScheduledQueries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListScheduledQueriesRequest.fromBuffer(value),
        ($0.ListScheduledQueriesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PrepareQueryRequest, $0.PrepareQueryResponse>(
            'PrepareQuery',
            prepareQuery_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PrepareQueryRequest.fromBuffer(value),
            ($0.PrepareQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.QueryRequest, $0.QueryResponse>(
        'Query',
        query_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.QueryRequest.fromBuffer(value),
        ($0.QueryResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
            'TagResource',
            tagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TagResourceRequest.fromBuffer(value),
            ($0.TagResourceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceRequest.fromBuffer(value),
            ($0.UntagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAccountSettingsRequest,
            $0.UpdateAccountSettingsResponse>(
        'UpdateAccountSettings',
        updateAccountSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAccountSettingsRequest.fromBuffer(value),
        ($0.UpdateAccountSettingsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateScheduledQueryRequest, $1.Empty>(
        'UpdateScheduledQuery',
        updateScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateScheduledQueryRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
  }

  $async.Future<$0.CancelQueryResponse> cancelQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CancelQueryRequest> $request) async {
    return cancelQuery($call, await $request);
  }

  $async.Future<$0.CancelQueryResponse> cancelQuery(
      $grpc.ServiceCall call, $0.CancelQueryRequest request);

  $async.Future<$0.CreateScheduledQueryResponse> createScheduledQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateScheduledQueryRequest> $request) async {
    return createScheduledQuery($call, await $request);
  }

  $async.Future<$0.CreateScheduledQueryResponse> createScheduledQuery(
      $grpc.ServiceCall call, $0.CreateScheduledQueryRequest request);

  $async.Future<$1.Empty> deleteScheduledQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteScheduledQueryRequest> $request) async {
    return deleteScheduledQuery($call, await $request);
  }

  $async.Future<$1.Empty> deleteScheduledQuery(
      $grpc.ServiceCall call, $0.DeleteScheduledQueryRequest request);

  $async.Future<$0.DescribeAccountSettingsResponse> describeAccountSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAccountSettingsRequest> $request) async {
    return describeAccountSettings($call, await $request);
  }

  $async.Future<$0.DescribeAccountSettingsResponse> describeAccountSettings(
      $grpc.ServiceCall call, $0.DescribeAccountSettingsRequest request);

  $async.Future<$0.DescribeEndpointsResponse> describeEndpoints_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeEndpointsRequest> $request) async {
    return describeEndpoints($call, await $request);
  }

  $async.Future<$0.DescribeEndpointsResponse> describeEndpoints(
      $grpc.ServiceCall call, $0.DescribeEndpointsRequest request);

  $async.Future<$0.DescribeScheduledQueryResponse> describeScheduledQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeScheduledQueryRequest> $request) async {
    return describeScheduledQuery($call, await $request);
  }

  $async.Future<$0.DescribeScheduledQueryResponse> describeScheduledQuery(
      $grpc.ServiceCall call, $0.DescribeScheduledQueryRequest request);

  $async.Future<$1.Empty> executeScheduledQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ExecuteScheduledQueryRequest> $request) async {
    return executeScheduledQuery($call, await $request);
  }

  $async.Future<$1.Empty> executeScheduledQuery(
      $grpc.ServiceCall call, $0.ExecuteScheduledQueryRequest request);

  $async.Future<$0.ListScheduledQueriesResponse> listScheduledQueries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListScheduledQueriesRequest> $request) async {
    return listScheduledQueries($call, await $request);
  }

  $async.Future<$0.ListScheduledQueriesResponse> listScheduledQueries(
      $grpc.ServiceCall call, $0.ListScheduledQueriesRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.PrepareQueryResponse> prepareQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PrepareQueryRequest> $request) async {
    return prepareQuery($call, await $request);
  }

  $async.Future<$0.PrepareQueryResponse> prepareQuery(
      $grpc.ServiceCall call, $0.PrepareQueryRequest request);

  $async.Future<$0.QueryResponse> query_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.QueryRequest> $request) async {
    return query($call, await $request);
  }

  $async.Future<$0.QueryResponse> query(
      $grpc.ServiceCall call, $0.QueryRequest request);

  $async.Future<$0.TagResourceResponse> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceResponse> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.UntagResourceResponse> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceResponse> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.UpdateAccountSettingsResponse> updateAccountSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateAccountSettingsRequest> $request) async {
    return updateAccountSettings($call, await $request);
  }

  $async.Future<$0.UpdateAccountSettingsResponse> updateAccountSettings(
      $grpc.ServiceCall call, $0.UpdateAccountSettingsRequest request);

  $async.Future<$1.Empty> updateScheduledQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateScheduledQueryRequest> $request) async {
    return updateScheduledQuery($call, await $request);
  }

  $async.Future<$1.Empty> updateScheduledQuery(
      $grpc.ServiceCall call, $0.UpdateScheduledQueryRequest request);
}
