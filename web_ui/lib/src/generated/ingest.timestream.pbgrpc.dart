// This is a generated file - do not edit.
//
// Generated from ingest.timestream.proto.

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
import 'ingest.timestream.pb.dart' as $0;

export 'ingest.timestream.pb.dart';

/// TimestreamWriteService provides ingest.timestream API operations.
@$pb.GrpcServiceName('ingest.timestream.TimestreamWriteService')
class TimestreamWriteServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  TimestreamWriteServiceClient(super.channel,
      {super.options, super.interceptors});

  /// Creates a new Timestream batch load task. A batch load task processes data from a CSV source in an S3 location and writes to a Timestream table. A mapping from source to target is defined in a batc...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateBatchLoadTaskResponse> createBatchLoadTask(
    $0.CreateBatchLoadTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createBatchLoadTask, request, options: options);
  }

  /// Creates a new Timestream database. If the KMS key is not specified, the database will be encrypted with a Timestream managed KMS key located in your account. For more information, see Amazon Web Se...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateDatabaseResponse> createDatabase(
    $0.CreateDatabaseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDatabase, request, options: options);
  }

  /// Adds a new table to an existing database in your account. In an Amazon Web Services account, table names must be at least unique within each Region if they are in the same database. You might have ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateTableResponse> createTable(
    $0.CreateTableRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTable, request, options: options);
  }

  /// Deletes a given Timestream database. This is an irreversible operation. After a database is deleted, the time-series data from its tables cannot be recovered. All tables in the database must be del...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> deleteDatabase(
    $0.DeleteDatabaseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDatabase, request, options: options);
  }

  /// Deletes a given Timestream table. This is an irreversible operation. After a Timestream database table is deleted, the time-series data stored in the table cannot be recovered. Due to the nature of...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> deleteTable(
    $0.DeleteTableRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTable, request, options: options);
  }

  /// Returns information about the batch load task, including configurations, mappings, progress, and other details. Service quotas apply. See code sample for details.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeBatchLoadTaskResponse> describeBatchLoadTask(
    $0.DescribeBatchLoadTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeBatchLoadTask, request, options: options);
  }

  /// Returns information about the database, including the database name, time that the database was created, and the total number of tables found within the database. Service quotas apply. See code sam...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeDatabaseResponse> describeDatabase(
    $0.DescribeDatabaseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeDatabase, request, options: options);
  }

  /// Returns a list of available endpoints to make Timestream API calls against. This API operation is available through both the Write and Query APIs. Because the Timestream SDKs are designed to transp...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeEndpointsResponse> describeEndpoints(
    $0.DescribeEndpointsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeEndpoints, request, options: options);
  }

  /// Returns information about the table, including the table name, database name, retention duration of the memory store and the magnetic store. Service quotas apply. See code sample for details.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeTableResponse> describeTable(
    $0.DescribeTableRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeTable, request, options: options);
  }

  /// Provides a list of batch load tasks, along with the name, status, when the task is resumable until, and other details. See code sample for details.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListBatchLoadTasksResponse> listBatchLoadTasks(
    $0.ListBatchLoadTasksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listBatchLoadTasks, request, options: options);
  }

  /// Returns a list of your Timestream databases. Service quotas apply. See code sample for details.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListDatabasesResponse> listDatabases(
    $0.ListDatabasesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDatabases, request, options: options);
  }

  /// Provides a list of tables, along with the name, status, and retention properties of each table. See code sample for details.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListTablesResponse> listTables(
    $0.ListTablesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTables, request, options: options);
  }

  /// Lists all tags on a Timestream resource.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  ///
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ResumeBatchLoadTaskResponse> resumeBatchLoadTask(
    $0.ResumeBatchLoadTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resumeBatchLoadTask, request, options: options);
  }

  /// Associates a set of tags with a Timestream resource. You can then activate these user-defined tags so that they appear on the Billing and Cost Management console for cost allocation tracking.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Removes the association of tags from a Timestream resource.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Modifies the KMS key for an existing database. While updating the database, you must specify the database name and the identifier of the new KMS key to be used (KmsKeyId). If there are any concurre...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateDatabaseResponse> updateDatabase(
    $0.UpdateDatabaseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDatabase, request, options: options);
  }

  /// Modifies the retention duration of the memory store and magnetic store for your Timestream table. Note that the change in retention duration takes effect immediately. For example, if the retention ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateTableResponse> updateTable(
    $0.UpdateTableRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTable, request, options: options);
  }

  /// Enables you to write your time-series data into Timestream. You can specify a single data point or a batch of data points to be inserted into the system. Timestream offers you a flexible schema tha...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.WriteRecordsResponse> writeRecords(
    $0.WriteRecordsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$writeRecords, request, options: options);
  }

  // method descriptors

  static final _$createBatchLoadTask = $grpc.ClientMethod<
          $0.CreateBatchLoadTaskRequest, $0.CreateBatchLoadTaskResponse>(
      '/ingest.timestream.TimestreamWriteService/CreateBatchLoadTask',
      ($0.CreateBatchLoadTaskRequest value) => value.writeToBuffer(),
      $0.CreateBatchLoadTaskResponse.fromBuffer);
  static final _$createDatabase =
      $grpc.ClientMethod<$0.CreateDatabaseRequest, $0.CreateDatabaseResponse>(
          '/ingest.timestream.TimestreamWriteService/CreateDatabase',
          ($0.CreateDatabaseRequest value) => value.writeToBuffer(),
          $0.CreateDatabaseResponse.fromBuffer);
  static final _$createTable =
      $grpc.ClientMethod<$0.CreateTableRequest, $0.CreateTableResponse>(
          '/ingest.timestream.TimestreamWriteService/CreateTable',
          ($0.CreateTableRequest value) => value.writeToBuffer(),
          $0.CreateTableResponse.fromBuffer);
  static final _$deleteDatabase =
      $grpc.ClientMethod<$0.DeleteDatabaseRequest, $1.Empty>(
          '/ingest.timestream.TimestreamWriteService/DeleteDatabase',
          ($0.DeleteDatabaseRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteTable =
      $grpc.ClientMethod<$0.DeleteTableRequest, $1.Empty>(
          '/ingest.timestream.TimestreamWriteService/DeleteTable',
          ($0.DeleteTableRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$describeBatchLoadTask = $grpc.ClientMethod<
          $0.DescribeBatchLoadTaskRequest, $0.DescribeBatchLoadTaskResponse>(
      '/ingest.timestream.TimestreamWriteService/DescribeBatchLoadTask',
      ($0.DescribeBatchLoadTaskRequest value) => value.writeToBuffer(),
      $0.DescribeBatchLoadTaskResponse.fromBuffer);
  static final _$describeDatabase = $grpc.ClientMethod<
          $0.DescribeDatabaseRequest, $0.DescribeDatabaseResponse>(
      '/ingest.timestream.TimestreamWriteService/DescribeDatabase',
      ($0.DescribeDatabaseRequest value) => value.writeToBuffer(),
      $0.DescribeDatabaseResponse.fromBuffer);
  static final _$describeEndpoints = $grpc.ClientMethod<
          $0.DescribeEndpointsRequest, $0.DescribeEndpointsResponse>(
      '/ingest.timestream.TimestreamWriteService/DescribeEndpoints',
      ($0.DescribeEndpointsRequest value) => value.writeToBuffer(),
      $0.DescribeEndpointsResponse.fromBuffer);
  static final _$describeTable =
      $grpc.ClientMethod<$0.DescribeTableRequest, $0.DescribeTableResponse>(
          '/ingest.timestream.TimestreamWriteService/DescribeTable',
          ($0.DescribeTableRequest value) => value.writeToBuffer(),
          $0.DescribeTableResponse.fromBuffer);
  static final _$listBatchLoadTasks = $grpc.ClientMethod<
          $0.ListBatchLoadTasksRequest, $0.ListBatchLoadTasksResponse>(
      '/ingest.timestream.TimestreamWriteService/ListBatchLoadTasks',
      ($0.ListBatchLoadTasksRequest value) => value.writeToBuffer(),
      $0.ListBatchLoadTasksResponse.fromBuffer);
  static final _$listDatabases =
      $grpc.ClientMethod<$0.ListDatabasesRequest, $0.ListDatabasesResponse>(
          '/ingest.timestream.TimestreamWriteService/ListDatabases',
          ($0.ListDatabasesRequest value) => value.writeToBuffer(),
          $0.ListDatabasesResponse.fromBuffer);
  static final _$listTables =
      $grpc.ClientMethod<$0.ListTablesRequest, $0.ListTablesResponse>(
          '/ingest.timestream.TimestreamWriteService/ListTables',
          ($0.ListTablesRequest value) => value.writeToBuffer(),
          $0.ListTablesResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/ingest.timestream.TimestreamWriteService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$resumeBatchLoadTask = $grpc.ClientMethod<
          $0.ResumeBatchLoadTaskRequest, $0.ResumeBatchLoadTaskResponse>(
      '/ingest.timestream.TimestreamWriteService/ResumeBatchLoadTask',
      ($0.ResumeBatchLoadTaskRequest value) => value.writeToBuffer(),
      $0.ResumeBatchLoadTaskResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/ingest.timestream.TimestreamWriteService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/ingest.timestream.TimestreamWriteService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateDatabase =
      $grpc.ClientMethod<$0.UpdateDatabaseRequest, $0.UpdateDatabaseResponse>(
          '/ingest.timestream.TimestreamWriteService/UpdateDatabase',
          ($0.UpdateDatabaseRequest value) => value.writeToBuffer(),
          $0.UpdateDatabaseResponse.fromBuffer);
  static final _$updateTable =
      $grpc.ClientMethod<$0.UpdateTableRequest, $0.UpdateTableResponse>(
          '/ingest.timestream.TimestreamWriteService/UpdateTable',
          ($0.UpdateTableRequest value) => value.writeToBuffer(),
          $0.UpdateTableResponse.fromBuffer);
  static final _$writeRecords =
      $grpc.ClientMethod<$0.WriteRecordsRequest, $0.WriteRecordsResponse>(
          '/ingest.timestream.TimestreamWriteService/WriteRecords',
          ($0.WriteRecordsRequest value) => value.writeToBuffer(),
          $0.WriteRecordsResponse.fromBuffer);
}

@$pb.GrpcServiceName('ingest.timestream.TimestreamWriteService')
abstract class TimestreamWriteServiceBase extends $grpc.Service {
  $core.String get $name => 'ingest.timestream.TimestreamWriteService';

  TimestreamWriteServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateBatchLoadTaskRequest,
            $0.CreateBatchLoadTaskResponse>(
        'CreateBatchLoadTask',
        createBatchLoadTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateBatchLoadTaskRequest.fromBuffer(value),
        ($0.CreateBatchLoadTaskResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDatabaseRequest,
            $0.CreateDatabaseResponse>(
        'CreateDatabase',
        createDatabase_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDatabaseRequest.fromBuffer(value),
        ($0.CreateDatabaseResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateTableRequest, $0.CreateTableResponse>(
            'CreateTable',
            createTable_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateTableRequest.fromBuffer(value),
            ($0.CreateTableResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDatabaseRequest, $1.Empty>(
        'DeleteDatabase',
        deleteDatabase_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDatabaseRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTableRequest, $1.Empty>(
        'DeleteTable',
        deleteTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteTableRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeBatchLoadTaskRequest,
            $0.DescribeBatchLoadTaskResponse>(
        'DescribeBatchLoadTask',
        describeBatchLoadTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeBatchLoadTaskRequest.fromBuffer(value),
        ($0.DescribeBatchLoadTaskResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeDatabaseRequest,
            $0.DescribeDatabaseResponse>(
        'DescribeDatabase',
        describeDatabase_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeDatabaseRequest.fromBuffer(value),
        ($0.DescribeDatabaseResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeEndpointsRequest,
            $0.DescribeEndpointsResponse>(
        'DescribeEndpoints',
        describeEndpoints_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeEndpointsRequest.fromBuffer(value),
        ($0.DescribeEndpointsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeTableRequest, $0.DescribeTableResponse>(
            'DescribeTable',
            describeTable_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeTableRequest.fromBuffer(value),
            ($0.DescribeTableResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListBatchLoadTasksRequest,
            $0.ListBatchLoadTasksResponse>(
        'ListBatchLoadTasks',
        listBatchLoadTasks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListBatchLoadTasksRequest.fromBuffer(value),
        ($0.ListBatchLoadTasksResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListDatabasesRequest, $0.ListDatabasesResponse>(
            'ListDatabases',
            listDatabases_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDatabasesRequest.fromBuffer(value),
            ($0.ListDatabasesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTablesRequest, $0.ListTablesResponse>(
        'ListTables',
        listTables_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListTablesRequest.fromBuffer(value),
        ($0.ListTablesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResumeBatchLoadTaskRequest,
            $0.ResumeBatchLoadTaskResponse>(
        'ResumeBatchLoadTask',
        resumeBatchLoadTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResumeBatchLoadTaskRequest.fromBuffer(value),
        ($0.ResumeBatchLoadTaskResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.UpdateDatabaseRequest,
            $0.UpdateDatabaseResponse>(
        'UpdateDatabase',
        updateDatabase_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDatabaseRequest.fromBuffer(value),
        ($0.UpdateDatabaseResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateTableRequest, $0.UpdateTableResponse>(
            'UpdateTable',
            updateTable_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateTableRequest.fromBuffer(value),
            ($0.UpdateTableResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.WriteRecordsRequest, $0.WriteRecordsResponse>(
            'WriteRecords',
            writeRecords_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.WriteRecordsRequest.fromBuffer(value),
            ($0.WriteRecordsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.CreateBatchLoadTaskResponse> createBatchLoadTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateBatchLoadTaskRequest> $request) async {
    return createBatchLoadTask($call, await $request);
  }

  $async.Future<$0.CreateBatchLoadTaskResponse> createBatchLoadTask(
      $grpc.ServiceCall call, $0.CreateBatchLoadTaskRequest request);

  $async.Future<$0.CreateDatabaseResponse> createDatabase_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDatabaseRequest> $request) async {
    return createDatabase($call, await $request);
  }

  $async.Future<$0.CreateDatabaseResponse> createDatabase(
      $grpc.ServiceCall call, $0.CreateDatabaseRequest request);

  $async.Future<$0.CreateTableResponse> createTable_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateTableRequest> $request) async {
    return createTable($call, await $request);
  }

  $async.Future<$0.CreateTableResponse> createTable(
      $grpc.ServiceCall call, $0.CreateTableRequest request);

  $async.Future<$1.Empty> deleteDatabase_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDatabaseRequest> $request) async {
    return deleteDatabase($call, await $request);
  }

  $async.Future<$1.Empty> deleteDatabase(
      $grpc.ServiceCall call, $0.DeleteDatabaseRequest request);

  $async.Future<$1.Empty> deleteTable_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteTableRequest> $request) async {
    return deleteTable($call, await $request);
  }

  $async.Future<$1.Empty> deleteTable(
      $grpc.ServiceCall call, $0.DeleteTableRequest request);

  $async.Future<$0.DescribeBatchLoadTaskResponse> describeBatchLoadTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeBatchLoadTaskRequest> $request) async {
    return describeBatchLoadTask($call, await $request);
  }

  $async.Future<$0.DescribeBatchLoadTaskResponse> describeBatchLoadTask(
      $grpc.ServiceCall call, $0.DescribeBatchLoadTaskRequest request);

  $async.Future<$0.DescribeDatabaseResponse> describeDatabase_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeDatabaseRequest> $request) async {
    return describeDatabase($call, await $request);
  }

  $async.Future<$0.DescribeDatabaseResponse> describeDatabase(
      $grpc.ServiceCall call, $0.DescribeDatabaseRequest request);

  $async.Future<$0.DescribeEndpointsResponse> describeEndpoints_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeEndpointsRequest> $request) async {
    return describeEndpoints($call, await $request);
  }

  $async.Future<$0.DescribeEndpointsResponse> describeEndpoints(
      $grpc.ServiceCall call, $0.DescribeEndpointsRequest request);

  $async.Future<$0.DescribeTableResponse> describeTable_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeTableRequest> $request) async {
    return describeTable($call, await $request);
  }

  $async.Future<$0.DescribeTableResponse> describeTable(
      $grpc.ServiceCall call, $0.DescribeTableRequest request);

  $async.Future<$0.ListBatchLoadTasksResponse> listBatchLoadTasks_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListBatchLoadTasksRequest> $request) async {
    return listBatchLoadTasks($call, await $request);
  }

  $async.Future<$0.ListBatchLoadTasksResponse> listBatchLoadTasks(
      $grpc.ServiceCall call, $0.ListBatchLoadTasksRequest request);

  $async.Future<$0.ListDatabasesResponse> listDatabases_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDatabasesRequest> $request) async {
    return listDatabases($call, await $request);
  }

  $async.Future<$0.ListDatabasesResponse> listDatabases(
      $grpc.ServiceCall call, $0.ListDatabasesRequest request);

  $async.Future<$0.ListTablesResponse> listTables_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTablesRequest> $request) async {
    return listTables($call, await $request);
  }

  $async.Future<$0.ListTablesResponse> listTables(
      $grpc.ServiceCall call, $0.ListTablesRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ResumeBatchLoadTaskResponse> resumeBatchLoadTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ResumeBatchLoadTaskRequest> $request) async {
    return resumeBatchLoadTask($call, await $request);
  }

  $async.Future<$0.ResumeBatchLoadTaskResponse> resumeBatchLoadTask(
      $grpc.ServiceCall call, $0.ResumeBatchLoadTaskRequest request);

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

  $async.Future<$0.UpdateDatabaseResponse> updateDatabase_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDatabaseRequest> $request) async {
    return updateDatabase($call, await $request);
  }

  $async.Future<$0.UpdateDatabaseResponse> updateDatabase(
      $grpc.ServiceCall call, $0.UpdateDatabaseRequest request);

  $async.Future<$0.UpdateTableResponse> updateTable_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateTableRequest> $request) async {
    return updateTable($call, await $request);
  }

  $async.Future<$0.UpdateTableResponse> updateTable(
      $grpc.ServiceCall call, $0.UpdateTableRequest request);

  $async.Future<$0.WriteRecordsResponse> writeRecords_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.WriteRecordsRequest> $request) async {
    return writeRecords($call, await $request);
  }

  $async.Future<$0.WriteRecordsResponse> writeRecords(
      $grpc.ServiceCall call, $0.WriteRecordsRequest request);
}
